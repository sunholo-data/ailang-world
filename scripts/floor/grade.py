#!/usr/bin/env python3
"""Row 93 floor harness — the grader port (design §4.2, D-NF-1 = A; V9, V10, V37; AC1.2, AC1.3).

A port of the ailang eval harness's AILANG grading path at the pinned corpus commit:
``validateSolution`` -> ``runAILANGSolution`` -> ``AILANGRunner.Run`` (internal/eval_harness/
agent_validation.go, runner.go) and ``GradeStdout``/``CompareOutput`` (runner.go).

What is ported exactly:
  * the run argv: ``run --entry main --quiet --relax-modules [--caps C] [--ai-stub]
    [--net-allow-http --net-allow-localhost] benchmark/solution.ail [-- cli_args]``;
  * a FRESH workspace holding ``benchmark/solution.ail`` with ``input_files`` seeded at its
    root, the child's cwd = that workspace, stdin piped when the spec has one (else /dev/null);
  * the 10 s timeout (process group killed), the 1 MiB output caps, ``runtime_ok = exit 0``,
    ``compile_ok`` false only on a failed exit whose stderr names a parse/type/syntax error;
  * ``CompareOutput``: TrimSpace -> exact -> jsonEqual -> normalizedEqual, with Go's rules
    (Go's ``unicode.IsSpace`` trim set, RE2's ASCII ``\\s``/``\\d``/``\\b``, ``json.Unmarshal``
    into ``interface{}`` — every number a float64, bool never equal to a number — and
    ``strconv.FormatFloat(f, 'f', -1, 64)`` for decimal canonicalisation).

What deliberately differs (and why):
  * the binary is an EXPLICIT absolute path (``$TOOL``, the v0.51.0 tool binary the World arm
    uses), never ``ailang`` from PATH (V21: PATH is a -dirty dev build). Its sha256 is recorded
    on every grade;
  * ``--stdlib-path`` is never passed: the Go harness adds ``--stdlib-path <cwd>/std`` when its
    cwd holds a stdlib (V37 — grading against the ailang repo's SOURCE stdlib). The workspace
    here is a fresh temp dir, asserted to hold no ``std/``, so the binary's built-in stdlib is
    the one graded against;
  * a missing, empty or placeholder-only solution never reaches the binary; it is reported as
    ``solution_state`` (classified a logic_error, never a harness fault — §4.2).
Only the default grading mode exists in the gate tier (V7); ``quine``/``prefix_line`` specs are
refused rather than mis-graded.
"""
from __future__ import annotations

import decimal
import hashlib
import http.server
import json
import math
import os
import re
import shutil
import signal
import subprocess
import tempfile
import threading
import time

GRADER_VERSION = 'floor-grade/1 (port of ailang runner.go+agent_validation.go @ corpus commit)'
DEFAULT_TOOL_BIN = os.path.expanduser('~/.pinned-ailang-tools/v0.51.0/ailang')
TIMEOUT_S = 10.0
MAX_OUTPUT = 1 * 1024 * 1024  # runner.go MaxOutputSize
# agent_runner_multi.go: the placeholder seeded into benchmark/solution.ail before the agent runs
GO_PLACEHOLDER = ('module benchmark/solution\n'
                  '// DO NOT CHANGE THE MODULE DECLARATION ABOVE!\n'
                  '// TODO: Add your solution code below\n'
                  '\n')
COMPILE_MARKERS = ('parse error', 'type error', 'syntax error')
TIMEOUT_STDERR = 'execution timed out'


class GradeSetupError(Exception):
    """The grader could not be set up (a harness_setup fault, never a model outcome)."""


# ---------------------------------------------------------------- CompareOutput port

# Go strings.TrimSpace trims unicode.IsSpace: Latin-1 '\t','\n','\v','\f','\r',' ',U+0085,
# U+00A0, plus the Unicode White_Space set above Latin-1. (Python's str.strip differs: it also
# strips U+001C..U+001F and does not match Go exactly.)
_GO_SPACE = ('\t\n\v\f\r \x85\xa0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008'
             '\u2009\u200a\u2028\u2029\u202f\u205f\u3000')


def go_trim_space(s: str) -> str:
    return s.strip(_GO_SPACE)


class _NotJSON(Exception):
    pass


def _reject_constant(name):
    raise _NotJSON(name)  # Go's encoding/json rejects NaN/Infinity/-Infinity


def _go_float(s: str) -> float:
    f = float(s)
    if math.isinf(f):
        raise _NotJSON(s)  # Go: number out of float64 range is an Unmarshal error
    return f


def _go_json(s: str):
    try:
        return json.loads(s, parse_int=_go_float, parse_float=_go_float, parse_constant=_reject_constant)
    except (ValueError, _NotJSON, RecursionError):
        raise _NotJSON(s)


def _go_deep_equal(a, b) -> bool:
    """reflect.DeepEqual over json.Unmarshal's interface{} values (types must match)."""
    if type(a) is not type(b):
        return False
    if isinstance(a, dict):
        return a.keys() == b.keys() and all(_go_deep_equal(a[k], b[k]) for k in a)
    if isinstance(a, list):
        return len(a) == len(b) and all(_go_deep_equal(x, y) for x, y in zip(a, b))
    return a == b


def json_equal(a: str, b: str) -> bool:
    try:
        va = _go_json(a)
        vb = _go_json(b)
    except _NotJSON:
        return False
    return _go_deep_equal(va, vb)


_RE_TRUE = re.compile(r'\bTrue\b', re.ASCII)
_RE_FALSE = re.compile(r'\bFalse\b', re.ASCII)
_RE_DECIMAL = re.compile(r'-?\d+\.\d+', re.ASCII)
# RE2's \s is [\t\n\f\r ] — no \v (Python's ASCII \s would include it)
_RE_PUNCT_SPACE = re.compile(r'[\t\n\f\r ]*([,\[\](){}:])[\t\n\f\r ]*')


def canon_decimal(s: str) -> str:
    """strconv.FormatFloat(ParseFloat(s), 'f', -1, 64); unparseable (overflow) -> unchanged."""
    f = float(s)
    if math.isinf(f):
        return s  # ParseFloat returns ErrRange; Go keeps the literal
    out = format(decimal.Decimal(repr(f)), 'f')
    if '.' in out:
        out = out.rstrip('0').rstrip('.')
    return out


def normalize_line(s: str) -> str:
    s = _RE_TRUE.sub('true', s)
    s = _RE_FALSE.sub('false', s)
    s = _RE_DECIMAL.sub(lambda m: canon_decimal(m.group(0)), s)
    s = _RE_PUNCT_SPACE.sub(r'\1', s)
    return go_trim_space(s)


def normalized_equal(expected: str, actual: str) -> bool:
    el = expected.split('\n')
    al = actual.split('\n')
    if len(el) != len(al):
        return False
    return all(normalize_line(e) == normalize_line(a) for e, a in zip(el, al))


def compare_output(expected: str, actual: str) -> bool:
    expected = go_trim_space(expected)
    actual = go_trim_space(actual)
    if expected == actual:
        return True
    if json_equal(expected, actual):
        return True
    return normalized_equal(expected, actual)


# ---------------------------------------------------------------- the loopback HTTP mock (V8)

def _go_canonical_header(k: str) -> str:
    """net/textproto.CanonicalMIMEHeaderKey for valid tokens."""
    if not re.fullmatch(r"[!#$%&'*+\-.^_`|~0-9A-Za-z]+", k):
        return k
    return '-'.join(p[:1].upper() + p[1:].lower() for p in k.split('-'))


def _go_json_encode(v) -> bytes:
    """json.NewEncoder(w).Encode(v): sorted map keys, compact, HTML-escaped, trailing newline."""
    s = json.dumps(v, sort_keys=True, separators=(',', ':'), ensure_ascii=False)
    s = s.replace('<', '\\u003c').replace('>', '\\u003e').replace('&', '\\u0026')
    s = s.replace('\u2028', '\\u2028').replace('\u2029', '\\u2029')
    return (s + '\n').encode()


class _MockHandler(http.server.BaseHTTPRequestHandler):
    protocol_version = 'HTTP/1.1'

    def _any(self):
        n = int(self.headers.get('Content-Length') or 0)
        body = self.rfile.read(n) if n > 0 else b''
        parsed = None
        if body:
            try:
                parsed = json.loads(body)
            except ValueError:
                parsed = body.decode('utf-8', errors='replace')
        headers = {}
        for k, v in self.headers.items():
            ck = _go_canonical_header(k)
            if ck != 'Host':  # net/http promotes Host to r.Host and drops it from r.Header
                headers[ck] = v  # last wins, as the Go mock
        out = _go_json_encode({'headers': headers, 'json': parsed, 'method': self.command,
                               'url': self.path})
        self.server.hits.append({'method': self.command, 'path': self.path, 'headers': headers,
                                 'json': parsed})
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.send_header('Content-Length', str(len(out)))
        self.end_headers()
        self.wfile.write(out)

    do_GET = do_POST = do_PUT = do_DELETE = do_PATCH = do_HEAD = do_OPTIONS = _any

    def log_message(self, *a):  # quiet
        pass


class HTTPMock:
    """A port of eval_harness.StartHTTPMock: 127.0.0.1, any request -> 200 + JSON echo.

    ``port=0`` binds an ephemeral port (a live task-run); a fixed port re-creates the mock a
    banked solution has baked into its source (the agreement re-grade)."""

    def __init__(self, port: int = 0):
        self.server = http.server.ThreadingHTTPServer(('127.0.0.1', port), _MockHandler)
        self.server.hits = []
        self.server.daemon_threads = True
        self.thread = threading.Thread(target=self.server.serve_forever, daemon=True)

    @property
    def url(self) -> str:
        return f'http://127.0.0.1:{self.server.server_address[1]}'

    @property
    def hits(self) -> list:
        return self.server.hits

    def __enter__(self):
        self.thread.start()
        return self

    def __exit__(self, *exc):
        self.server.shutdown()
        self.server.server_close()


# ---------------------------------------------------------------- the run

def sha256_file(path: str) -> str:
    h = hashlib.sha256()
    with open(path, 'rb') as f:
        for chunk in iter(lambda: f.read(1 << 20), b''):
            h.update(chunk)
    return h.hexdigest()


def resolve_tool(tool_bin: str) -> str:
    """The grading binary must be an explicit absolute path to an executable file (P12)."""
    if not tool_bin or not os.path.isabs(tool_bin):
        raise GradeSetupError(f'tool binary must be an absolute path, got {tool_bin!r} '
                              '(PATH `ailang` is never used to grade)')
    if not (os.path.isfile(tool_bin) and os.access(tool_bin, os.X_OK)):
        raise GradeSetupError(f'tool binary is not an executable file: {tool_bin}')
    return tool_bin


def build_argv(tool_bin: str, spec: dict) -> list:
    args = [tool_bin, 'run', '--entry', 'main', '--quiet', '--relax-modules']
    caps = [c for c in (spec.get('caps') or [])]
    if caps:
        args += ['--caps', ','.join(caps)]
    if any(str(c).strip().lower() == 'ai' for c in caps):
        args.append('--ai-stub')
    if spec.get('net_allow_localhost'):
        args += ['--net-allow-http', '--net-allow-localhost']
    args.append('benchmark/solution.ail')
    if spec.get('cli_args'):
        args.append('--')
        args += [str(a) for a in spec['cli_args']]
    return args


def solution_state(solution: bytes | None, placeholder: bytes | None = None) -> str:
    if solution is None:
        return 'missing'
    if len(solution) == 0:
        return 'empty'
    ph = GO_PLACEHOLDER.encode() if placeholder is None else placeholder
    if solution.strip() == ph.strip():
        return 'placeholder'
    return 'present'


def _run_guarded(argv: list, cwd: str, stdin: bytes | None, timeout: float) -> dict:
    t0 = time.monotonic()
    p = subprocess.Popen(argv, cwd=cwd, stdin=subprocess.PIPE if stdin is not None else subprocess.DEVNULL,
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE, start_new_session=True)
    try:
        out, err = p.communicate(input=stdin, timeout=timeout)
        timed_out = False
    except subprocess.TimeoutExpired:
        try:
            os.killpg(p.pid, signal.SIGKILL)
        except ProcessLookupError:
            pass
        out, err = p.communicate()
        timed_out = True
    return {'stdout': out[:MAX_OUTPUT], 'stderr': err[:MAX_OUTPUT], 'rc': p.returncode,
            'timed_out': timed_out, 'elapsed_s': time.monotonic() - t0}


def grade(solution: bytes | None, spec: dict, tool_bin: str = DEFAULT_TOOL_BIN, *,
          timeout: float = TIMEOUT_S, placeholder: bytes | None = None,
          tmp_root: str | None = None) -> dict:
    """Grade one solution against one spec. Returns the ailang result fields plus provenance.

    The caller owns the HTTP mock (it must stay up from the agent's run through grading, §4.2);
    a spec with ``net_allow_localhost`` is graded against whatever URL its source names."""
    tool = resolve_tool(tool_bin)
    mode = spec.get('grading') or ''
    if mode:
        raise GradeSetupError(f'grading mode {mode!r} is not ported (gate tier uses the default only)')
    if spec.get('grade_entrypoint') or spec.get('solution_files'):
        raise GradeSetupError('multi-file workspace grading is not ported')
    expected = spec.get('expected_stdout')
    if expected is None:
        raise GradeSetupError('spec has no expected_stdout')
    result = {'compile_ok': False, 'runtime_ok': False, 'stdout_ok': False, 'stdout': '',
              'stderr': '', 'solution_state': solution_state(solution, placeholder),
              'timed_out': False, 'exit_code': None, 'argv': None,
              'tool_bin': tool, 'tool_sha256': sha256_file(tool), 'grader': GRADER_VERSION,
              'solution_sha256': hashlib.sha256(solution).hexdigest() if solution is not None else None}
    if result['solution_state'] != 'present':
        result['stderr'] = f'solution {result["solution_state"]}'
        return result
    ws = tempfile.mkdtemp(prefix='floor-grade-', dir=tmp_root)
    try:
        if os.path.exists(os.path.join(ws, 'std')):
            raise GradeSetupError('grading workspace holds a std/ (V37)')
        os.makedirs(os.path.join(ws, 'benchmark'))
        with open(os.path.join(ws, 'benchmark', 'solution.ail'), 'wb') as f:
            f.write(solution)
        for name, content in (spec.get('input_files') or {}).items():
            fp = os.path.join(ws, name)
            if not os.path.realpath(fp).startswith(os.path.realpath(ws) + os.sep):
                raise GradeSetupError(f'input file escapes the workspace: {name}')
            os.makedirs(os.path.dirname(fp), exist_ok=True)
            with open(fp, 'wb') as f:
                f.write((content or '').encode())
        argv = build_argv(tool, spec)
        result['argv'] = argv
        # The digest recorded is of the binary ACTUALLY executed (argv[0] as the OS resolves
        # it), so a grader that drifted to PATH `ailang` records a digest != $TOOL (P12).
        executed = argv[0] if os.path.isabs(argv[0]) else shutil.which(argv[0])
        if not executed:
            raise GradeSetupError(f'cannot resolve the executed binary {argv[0]!r}')
        result['executed_bin'] = executed
        result['tool_sha256'] = sha256_file(executed)
        stdin = spec.get('stdin')
        r = _run_guarded(argv, ws, stdin.encode() if stdin else None, timeout)
    finally:
        shutil.rmtree(ws, ignore_errors=True)
    stdout = r['stdout'].decode('utf-8', errors='surrogateescape')
    stderr = r['stderr'].decode('utf-8', errors='replace')
    result.update(stdout=stdout.encode('utf-8', errors='surrogateescape').decode('utf-8', errors='replace'),
                  exit_code=r['rc'], elapsed_s=round(r['elapsed_s'], 3))
    if r['timed_out']:
        result.update(compile_ok=True, runtime_ok=False, timed_out=True, stderr=TIMEOUT_STDERR)
        return result
    compile_ok = True
    if r['rc'] != 0 and any(m in stderr for m in COMPILE_MARKERS):
        compile_ok = False
    runtime_ok = r['rc'] == 0
    result.update(compile_ok=compile_ok, runtime_ok=runtime_ok, stderr=stderr,
                  stdout_ok=bool(runtime_ok and compare_output(expected, stdout)))
    return result


def error_category(g: dict) -> str:
    """eval_harness.CategorizeError over a grade (none/compile/runtime/logic).

    A non-present solution is a logic_error (§4.2: never a harness fault)."""
    if g['stdout_ok']:
        return 'none'
    if g['solution_state'] != 'present':
        return 'logic_error'
    if not g['compile_ok']:
        return 'compile_error'
    if not g['runtime_ok']:
        return 'runtime_error'
    return 'logic_error'


def main(argv=None) -> int:
    import argparse
    import corpus as corpus_mod
    ap = argparse.ArgumentParser(description='Grade one solution file against one pinned task.')
    ap.add_argument('--repo', required=True)
    ap.add_argument('--commit', required=True)
    ap.add_argument('--task', required=True)
    ap.add_argument('--solution', required=True)
    ap.add_argument('--tool-bin', default=DEFAULT_TOOL_BIN)
    a = ap.parse_args(argv)
    c = corpus_mod.load_corpus(repo=a.repo, commit=a.commit)
    task = c.by_id()[a.task]
    sol = open(a.solution, 'rb').read() if os.path.exists(a.solution) else None
    g = grade(sol, task.spec, a.tool_bin)
    g['error_category'] = error_category(g)
    print(json.dumps(g, indent=2, sort_keys=True))
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
