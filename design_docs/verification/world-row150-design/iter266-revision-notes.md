# Row 150 revision verification — iteration 266

Base: c21f30310e04941afd3ccfb6b8748177db7b06c8. Documentation only.
Both quorum artifacts were parsed, including every present reviewer's strongest
objection, catch and proposed fix, and each controller-in-session result.

Executed, bounded checks:

- `go test -count=1 -timeout 45s ./host/coordinator -run '^TestDecisionPacketCodec$'`: PASS (0.202s).
- `go test -count=1 -timeout 45s ./host/daemon -run '^TestWorkbench(GraphWiring|DecisionsPane|LiveCursorCrossesGap)$'`: PASS (0.493s); recorder-based, no browser.
- `bash scripts/tracked_binary_gate.sh --allow-png 'website/static/img/v1/' --png-file-cap 1048576 --png-dir-cap 8388608`: PASS, 0 binary blobs / 3516 tracked files. This inspects HEAD, not proposed images.
- `$HOME/.pinned-ailang/ailang --version`: v0.41.0, commit 24ee108.
- `$HOME/.pinned-ailang/ailang test --format json world/logepoch.ail`: rc=0, success=true, 10 passed, 0 failed/skipped/vacuous skips.
- `$HOME/.pinned-ailang/ailang ai-check -timeout 5s world/logepoch.ail`: rc=0, check.passed=true; verified identities sameRef and servesEntry; errors/counterexamples/skipped=0.
- `$HOME/.pinned-ailang/ailang test --format json world/contracts.ail`: BASE FINDING, rc=1, success=false, 7 passed, 4 no_generator/vacuous skips. Do not use this as passing TESTED evidence. This is not a claim that the repository's full verify gate fails.
- `$HOME/.pinned-ailang/ailang ai-check -timeout 5s world/contracts.ail`: rc=0, five identities verified, no skips/errors/counterexamples. Its imports make it unsuitable for the producer's single-source staging without extra work.

AILANG subprocesses above were supervised with Python subprocess.run(timeout=40).
Protected-path OBSERVATION: 12 tracked paths (world/*.ail, go.mod, go.sum,
tools/launchd/*, charter); aggregate HEAD/working SHA256 both
`05852e6f08c102b42ab5952c4928a75779815512000a319de518d64e2a0f6202`; newline counts
HEAD=7363, working=7363. Compared git-show HEAD bytes with current bytes;
none was written.

No socket-denied result was encountered. No product Chrome run was attempted here.
The controller's Chrome 155 instrument results are in iter266-chrome-exit-probes.md;
those results do not establish product capture, dark forcing, or six-shot equality.

## PENDING controller measurements

Run these after the proposed tooling is implemented on a socket/browser-capable host.
`ADDR` is the scratch capture daemon's loopback URL; `PACKET` and `ENTRY` come from
its seed manifest, never a guessed route. All these commands fail on timeout.

Raw HTML repeatability, including request-time strings (2 seconds apart):

```bash
python3 - "$ADDR" "$PACKET" "$ENTRY" <<'PY'
import sys, time, urllib.request
base, packet, entry = sys.argv[1:]
paths = ['/workbench', '/workbench?object='+packet,
         '/workbench?from='+entry+'&entry='+entry]
for path in paths:
    a = urllib.request.urlopen(base+path, timeout=10).read()
    time.sleep(2)
    b = urllib.request.urlopen(base+path, timeout=10).read()
    assert a == b, 'request-time HTML difference: '+path
    assert b'<script' not in a, 'still mode emitted a script: '+path
    assert b'Live updates: off. Reload to refresh.' in a
print('PASS: exact HTML equality, static footer, no script')
PY
```

Chrome 155 dark-media probe (an instrument control, not a product image). This
preserves HOME, disconnects stdin, kills the process group on a 30-second deadline,
and checks the DOM's measured CSS/media values instead of trusting a flag name:

```bash
python3 - <<'PY'
import os, pathlib, signal, subprocess, tempfile
chrome = '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'
with tempfile.TemporaryDirectory(prefix='row150-dark-') as directory:
    root = pathlib.Path(directory)
    page = root/'probe.html'
    page.write_text('''<!doctype html><style>:root{--bg:#f7f6f2}
@media(prefers-color-scheme:dark){:root{--bg:#1f1e1b}}</style>
<body><script>document.body.textContent =
matchMedia('(prefers-color-scheme: dark)').matches + '|' +
getComputedStyle(document.documentElement).getPropertyValue('--bg').trim();
</script></body>''')
    for mode in ('light', 'dark'):
        args = [chrome, '--headless=new', '--no-first-run',
                '--no-default-browser-check', '--disable-background-networking',
                '--user-data-dir='+str(root/mode), '--dump-dom']
        if mode == 'dark': args += ['--force-dark-mode']
        args += [page.as_uri()]
        p = subprocess.Popen(args, stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
                             stderr=subprocess.PIPE, start_new_session=True)
        try: stdout, stderr = p.communicate(timeout=30)
        except subprocess.TimeoutExpired:
            os.killpg(p.pid, signal.SIGKILL)
            p.communicate(timeout=5)
            raise SystemExit('FAIL: Chrome dark probe timeout')
        want = b'true|#1f1e1b' if mode == 'dark' else b'false|#f7f6f2'
        assert p.returncode == 0 and want in stdout, (mode, p.returncode, stdout)
        print('PASS:', mode, want.decode())
PY
```

Product PNG completion/content/byte equality (the proposed script's supervisor
uses the same deadline and kill semantics; its verifier is Go image/png):

```bash
scripts/capture_marketing.sh --db "$DEMO_DB" --manifest "$DEMO_MANIFEST" \
  --proof-key-file "$DEMO_PROOF_KEY" --ailang-bin "$HOME/.pinned-ailang/ailang" \
  --chrome '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' \
  --out "$CAPTURE_OUT" --verify-twice
```

This command is a FUTURE acceptance command, not an existing script or a reported
pass. It must exercise all six routes in two isolated sets of profiles. Required
raw Chrome argv per shot:
`--headless=new --no-first-run --no-default-browser-check --disable-background-networking --user-data-dir=<unique-absolute-profile> --screenshot=<new-absolute-png> --window-size=W,H --force-device-scale-factor=1 --lang=en-US <loopback-shot-URL>`;
light shots add `--blink-settings=preferredColorScheme=1` and the dark shot `--blink-settings=preferredColorScheme=0` (controller V23 superseded `--force-dark-mode`, see iter266-colour-scheme-probes.md). Neither JS-disabling flag is used.
