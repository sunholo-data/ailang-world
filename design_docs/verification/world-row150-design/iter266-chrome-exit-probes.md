# iter-266 controller Chrome capture probes (first-party, 2026-10-14 ~22:58 CEST)
Chrome: Google Chrome 155.0.8059.41  at /Applications/Google Chrome.app; base flags: --headless=new --no-first-run --no-default-browser-check --disable-background-networking --user-data-dir=<isolated> ; HOME preserved; stdin </dev/null; 40 s deadline.
- --screenshot=<abs path> --window-size=800,400 file://p.html : rc=0 in 1-2 s, file written, B2/B3 isolated profiles byte-identical sha256 prefix a486c6ea7b48d460; chunks IHDR/IDAT/IEND only, 0 trailing bytes (compatible with tracked_binary_gate --allow-png).
- + --virtual-time-budget=2000 : rc=0, file written, same bytes.
- + --blink-settings=scriptEnabled=false : rc=0 in ~1 s and NO FILE WRITTEN (4 of 4 runs: run1, run2, D, F). SILENT FALSE GREEN.
- + --disable-javascript : NOT effective (js.html control: script ran in both J1 control and J2; identical PNG f23b481c50223af9, red background = script executed).
- Contrast iter-246 (Chrome 154.0.8037.98, scriptEnabled=false): PNG written, then process timed out at 30 s.
