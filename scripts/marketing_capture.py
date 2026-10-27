#!/usr/bin/env python3
"""Bounded capture subprocesses and six-shot capture primitives.

HOME must stay intact: relocated HOME prevents macOS Chrome from loading URLs
(see host/daemon/workbench_live_chrome_test.go). Profiles isolate each shot instead.
Runtime stdout is one JSON verdict; child output goes to temporary diagnostic files
so a descendant holding a pipe open cannot defeat the reap deadline.

Examples (destinations/readiness markers must be fresh):
  python3 scripts/marketing_capture.py shot --chrome /path/to/chrome \
      --out /tmp/new-shot.png --url http://127.0.0.1:7644/workbench --scheme light
  python3 scripts/marketing_capture.py daemon --ready-file /tmp/daemon.ready \
      --started-file /tmp/supervisor.ready -- /path/to/daemon
  python3 scripts/marketing_capture.py chrome-version --chrome /path/to/chrome
  python3 scripts/marketing_capture.py ready --base-url http://127.0.0.1:7644
  marketingcheck routes --manifest demo.json > /tmp/routes.json
  python3 scripts/marketing_capture.py pass --routes /tmp/routes.json \
      --base-url http://127.0.0.1:7644 --chrome /path/to/chrome --out-dir /tmp/new-pass
The daemon must publish ready-file; started-file then contains its PID. Signal the
supervisor with TERM/INT/HUP to stop it. File readiness is the M7a primitive; HTTP
readiness is also available as a standalone mode. No World mutations occur here.
"""
from __future__ import annotations

import argparse
import json
import math
import os
from pathlib import Path
import re
import signal
import stat
import subprocess
import sys
import tempfile
import threading
import time
import urllib.request

SHOT_TIMEOUT_S = 30.0
PASS_TIMEOUT_S = 300.0  # M7c must budget startup, six shots and cleanup within this.
STARTUP_TIMEOUT_S = 30.0
TERM_WAIT_S = 5.0
KILL_WAIT_S = 5.0
MIN_CHROME_MAJOR = 155


def verdict(name: str, detail: str = '', **fields: str) -> dict:
    return dict(verdict=name, detail=detail, **fields)


def kill_group(process: subprocess.Popen, sig: int) -> None:
    """The session leader's PID is also its process-group ID, even after exit."""
    try:
        os.killpg(process.pid, sig)
    except ProcessLookupError:
        pass


def spawn(argv: list[str], log) -> subprocess.Popen:
    # Inherit HOME and the rest of the environment; never inherit interactive input.
    return subprocess.Popen(argv, stdin=subprocess.DEVNULL, stdout=log, stderr=log,
                            start_new_session=True)


def png_verdict(path: Path) -> dict:
    """Use lstat: directories AND symlinks are refused, even symlinks to valid PNGs."""
    try:
        info = path.lstat()
    except FileNotFoundError:
        return verdict('MissingPNG', str(path))
    if not stat.S_ISREG(info.st_mode):
        return verdict('NotRegularFile', str(path))
    if info.st_size == 0:
        return verdict('EmptyPNG', str(path))
    return verdict('ok')


def wait_bounded(process: subprocess.Popen, timeout: float, stop: threading.Event | None) -> bool:
    """True when the child exited; False on deadline or when a supervisor signal arrived."""
    deadline = time.monotonic() + timeout
    while True:
        remaining = deadline - time.monotonic()
        if remaining <= 0 or (stop is not None and stop.is_set()):
            return False
        try:
            process.wait(timeout=min(remaining, 0.05))
            return True
        except subprocess.TimeoutExpired:
            pass


def shot(args: argparse.Namespace, stop: threading.Event | None = None) -> dict:
    path = Path(os.path.abspath(args.out))
    # Refuse all existing entries, including dangling links; never accept stale bytes.
    if os.path.lexists(path):
        return verdict('DestinationExists', str(path))
    with tempfile.TemporaryDirectory(prefix='marketing-chrome-') as profile, \
            tempfile.TemporaryFile() as log:
        argv = [args.chrome, '--headless=new', '--no-first-run',
                '--no-default-browser-check', '--disable-background-networking',
                f'--user-data-dir={Path(profile).resolve()}', f'--screenshot={path}',
                f'--window-size={args.width},{args.height}',
                '--force-device-scale-factor=1', '--lang=en-US',
                '--blink-settings=preferredColorScheme=' + ('0' if args.scheme == 'dark' else '1'),
                args.url]
        process = spawn(argv, log)
        try:
            # A TERM/INT/HUP to the daemon supervisor must not wait out the 30 s shot.
            if not wait_bounded(process, args.timeout, stop):
                kill_group(process, signal.SIGKILL)
                try:
                    process.communicate(timeout=KILL_WAIT_S)
                except subprocess.TimeoutExpired:
                    # Unreachable without an unkillable process; kept so a reap failure is loud.
                    return verdict('CleanupFailed', 'Chrome could not be reaped after KILL')
                if stop is not None and stop.is_set():
                    return verdict('Interrupted', 'supervisor signalled during the shot')
                return verdict('Timeout', 'Chrome deadline expired, regardless of PNG')
        finally:
            # Also collect descendants left behind by an exited Chrome leader.
            kill_group(process, signal.SIGKILL)
        if process.returncode != 0:
            return verdict('NonzeroExit', f'Chrome exit {process.returncode}')
        return png_verdict(path)


def cleanup_daemon(process: subprocess.Popen, term_wait: float, kill_wait: float) -> dict:
    kill_group(process, signal.SIGTERM)
    try:
        process.communicate(timeout=term_wait)
    except subprocess.TimeoutExpired:
        kill_group(process, signal.SIGKILL)
        try:
            process.communicate(timeout=kill_wait)
        except subprocess.TimeoutExpired:
            # Unreachable without an unkillable process; tests inject this via --unit-cleanup-failure.
            return verdict('CleanupFailed', 'daemon could not be reaped after KILL')
    # The leader may exit on TERM while descendants ignore it.
    kill_group(process, signal.SIGKILL)
    return verdict('ok')


def daemon(args: argparse.Namespace) -> dict:
    command = args.command[1:] if args.command[:1] == ['--'] else args.command
    if not command:
        return verdict('InvalidConfiguration', 'daemon requires -- COMMAND [ARG ...]')
    if os.path.lexists(args.ready_file) or os.path.lexists(args.started_file):
        return verdict('DestinationExists', 'readiness markers must be fresh')
    stop = threading.Event()
    previous = {}
    for sig in (signal.SIGTERM, signal.SIGINT, signal.SIGHUP):
        previous[sig] = signal.signal(sig, lambda *_: stop.set())
    result = verdict('ok')
    cleanup = verdict('ok')
    try:
        with tempfile.TemporaryFile() as log:
            process = spawn(command, log)
            try:
                deadline = time.monotonic() + args.startup_timeout
                while not Path(args.ready_file).is_file():
                    if stop.is_set():
                        break
                    if process.poll() is not None:
                        result = verdict('NonzeroExit', 'daemon exited before readiness')
                        break
                    if time.monotonic() >= deadline:
                        result = verdict('Timeout', 'daemon readiness deadline expired')
                        break
                    stop.wait(0.01)
                else:
                    Path(args.started_file).write_text(json.dumps({'pid': process.pid, 'supervisor_pid': os.getpid()}) + '\n')
                    if args.chrome:
                        result = shot(args, stop)
                    else:
                        deadline = time.monotonic() + args.lifetime
                        while not stop.is_set() and time.monotonic() < deadline:
                            if process.poll() is not None:
                                result = verdict('NonzeroExit', 'daemon exited during lifecycle')
                                break
                            stop.wait(0.01)
            finally:
                cleanup = cleanup_daemon(process, args.term_wait, args.kill_wait)
    finally:
        for sig, handler in previous.items():
            signal.signal(sig, handler)
    # Unit-only fault injection happens AFTER real cleanup, so tests never leak.
    if args.unit_cleanup_failure:
        cleanup = verdict('CleanupFailed', 'injected unit cleanup failure')
    if cleanup['verdict'] != 'ok':
        return dict(cleanup, shot_verdict=result['verdict'])
    return result


def chrome_version(args: argparse.Namespace) -> dict:
    with tempfile.TemporaryFile() as log:
        process = spawn([args.chrome, '--version'], log)
        try:
            if not wait_bounded(process, SHOT_TIMEOUT_S, None):
                kill_group(process, signal.SIGKILL)
                process.wait(timeout=KILL_WAIT_S)
                return verdict('Timeout', 'Chrome version deadline expired')
        finally:
            kill_group(process, signal.SIGKILL)
        if process.returncode != 0:
            return verdict('NonzeroExit', f'Chrome exit {process.returncode}')
        log.seek(0)
        version = log.read().decode('utf-8', errors='replace').strip()
    match = re.match(r'^(?:Google Chrome|Chromium) (\d+)\.\S+', version)
    if match is None or int(match[1]) < MIN_CHROME_MAJOR:
        return verdict('ChromeVersion', version)
    return verdict('ok', version=version)


def ready(args: argparse.Namespace) -> dict:
    deadline = time.monotonic() + args.startup_timeout
    healthy = threading.Event()

    def poll() -> None:
        while time.monotonic() < deadline:
            try:
                for route in ('/v1/health', '/v1/head'):
                    with urllib.request.urlopen(args.base_url.rstrip('/') + route,
                                                timeout=10) as response:
                        if response.status != 200:
                            break
                else:
                    if time.monotonic() < deadline:
                        healthy.set()
                    return
            except OSError:
                pass
            time.sleep(0.01)

    # A slow request must not extend the overall startup deadline (including overrides).
    threading.Thread(target=poll, daemon=True).start()
    if healthy.wait(max(0, deadline - time.monotonic())):
        return verdict('ok')
    return verdict('Timeout', 'HTTP readiness deadline expired')


def capture_pass(args: argparse.Namespace) -> dict:
    deadline = time.monotonic() + PASS_TIMEOUT_S
    shots = []
    destination = Path(args.out_dir)
    try:
        destination.mkdir()
    except FileExistsError:
        return verdict('DestinationExists', str(destination), shots=shots)
    except OSError as error:
        return verdict('SupervisorError', str(error), shots=shots)
    try:
        rows = json.loads(Path(args.routes).read_text())
        if not isinstance(rows, list) or len(rows) != 6:
            return verdict('InvalidConfiguration', 'routes must contain six rows', shots=shots)
        for i, row in enumerate(rows, 1):
            name = row['name']
            if (not isinstance(name, str) or Path(name).name != name or
                    row['width'] <= 0 or row['height'] <= 0 or
                    row['scheme'] not in ('light', 'dark')):
                return verdict('InvalidConfiguration', 'invalid shot route', shots=shots)
            remaining = deadline - time.monotonic()
            if remaining <= 0:
                return verdict('Timeout', 'pass deadline expired', shots=shots)
            shot_args = argparse.Namespace(
                chrome=args.chrome, out=str(destination / name),
                url=args.base_url.rstrip('/') + row['route'], width=row['width'],
                height=row['height'], scheme=row['scheme'],
                timeout=min(SHOT_TIMEOUT_S, remaining))
            try:
                result = shot(shot_args)
            except OSError as error:
                result = verdict('SupervisorError', str(error))
            shots.append(dict(name=name, verdict=result['verdict']))
            print(f'progress: {i}/6 {name} {result["verdict"]}', file=sys.stderr, flush=True)
            if result['verdict'] != 'ok':
                return dict(result, shots=shots)
        if time.monotonic() >= deadline:
            return verdict('Timeout', 'pass deadline expired', shots=shots)
        return verdict('ok', shots=shots)
    except (OSError, ValueError, KeyError, TypeError) as error:
        return verdict('InvalidConfiguration', str(error), shots=shots)


def positive(value: str) -> float:
    number = float(value)
    if not math.isfinite(number) or number <= 0:
        raise argparse.ArgumentTypeError('must be a positive finite number')
    return number


def shot_flags(parser: argparse.ArgumentParser, required: bool) -> None:
    parser.add_argument('--chrome', required=required, help='explicit Chrome executable')
    parser.add_argument('--out', required=required, help='new PNG destination (never overwrite)')
    parser.add_argument('--url', required=required)
    parser.add_argument('--width', type=int, default=1440)
    parser.add_argument('--height', type=int, default=1100)
    parser.add_argument('--scheme', choices=('light', 'dark'), default='light')
    parser.add_argument('--timeout', type=positive, default=SHOT_TIMEOUT_S,
                        help='shot deadline seconds (default 30; shorter values are unit control)')


def main() -> int:
    parser = argparse.ArgumentParser(
        description='Bounded subprocess primitives. Production deadlines: '
                    f'shot={SHOT_TIMEOUT_S:g}s, pass={PASS_TIMEOUT_S:g}s, '
                    f'startup={STARTUP_TIMEOUT_S:g}s, TERM={TERM_WAIT_S:g}s, '
                    f'KILL/reap={KILL_WAIT_S:g}s.',
        epilog='Runtime stdout: one JSON {verdict, detail, shot_verdict?}; '
               'exit 0 only for ok, 1 for all runtime refusals; argument errors exit 2. '
               'See module docstring for runnable examples.')
    sub = parser.add_subparsers(dest='mode', required=True)
    shot_flags(sub.add_parser('shot', help='one new screenshot'), required=True)
    lifecycle = sub.add_parser('daemon', help='start, await file readiness, and clean up on signals')
    lifecycle.add_argument('--ready-file', required=True, help='fresh marker written by daemon')
    lifecycle.add_argument('--started-file', required=True, help='fresh supervisor readiness JSON/PID')
    lifecycle.add_argument('--startup-timeout', type=positive, default=STARTUP_TIMEOUT_S)
    lifecycle.add_argument('--lifetime', type=positive, default=PASS_TIMEOUT_S,
                           help='maximum ready lifetime seconds, or stop earlier via TERM/INT/HUP')
    lifecycle.add_argument('--term-wait', type=positive, default=TERM_WAIT_S,
                           help='TERM wait seconds (default 5; shorter values are unit control)')
    lifecycle.add_argument('--kill-wait', type=positive, default=KILL_WAIT_S,
                           help='KILL/reap wait seconds (default 5; shorter values are unit control)')
    lifecycle.add_argument('--unit-cleanup-failure', action='store_true', help=argparse.SUPPRESS)
    shot_flags(lifecycle, required=False)
    lifecycle.add_argument('command', nargs=argparse.REMAINDER, help='-- COMMAND [ARG ...]')
    version = sub.add_parser('chrome-version', help='require Chrome/Chromium major >= 155')
    version.add_argument('--chrome', required=True, help='explicit Chrome executable')
    readiness = sub.add_parser('ready', help='await HTTP health and head readiness')
    readiness.add_argument('--base-url', required=True)
    readiness.add_argument('--startup-timeout', type=positive, default=STARTUP_TIMEOUT_S,
                           help='startup deadline seconds (default 30; shorter values are unit control)')
    capture = sub.add_parser('pass', help='capture six routes into a new directory')
    capture.add_argument('--routes', required=True, help='JSON from marketingcheck routes')
    capture.add_argument('--base-url', required=True)
    capture.add_argument('--chrome', required=True, help='explicit Chrome executable')
    capture.add_argument('--out-dir', required=True, help='new six-shot directory')
    args = parser.parse_args()
    if args.mode in ('shot', 'daemon') and (
            args.width <= 0 or args.height <= 0 or (args.chrome and (not args.out or not args.url))):
        result = verdict('InvalidConfiguration', 'positive viewport and complete shot flags required')
    else:
        try:
            result = {'shot': shot, 'daemon': daemon, 'chrome-version': chrome_version,
                      'ready': ready, 'pass': capture_pass}[args.mode](args)
        except OSError as error:
            result = verdict('SupervisorError', str(error))
    print(json.dumps(result), flush=True)
    return 0 if result['verdict'] == 'ok' else 1


if __name__ == '__main__':
    raise SystemExit(main())
