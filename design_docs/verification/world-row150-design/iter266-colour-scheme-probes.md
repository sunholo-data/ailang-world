# iter-266 controller colour-scheme probes (Chrome 155.0.8059.41, macOS AppleInterfaceStyle=Dark)
Base flags as iter266-chrome-exit-probes.md; probe page reports matchMedia(prefers-color-scheme: dark)|--bg.
- no flag (--dump-dom): true|#1f1e1b  -> headless Chrome FOLLOWS THE OS APPEARANCE; an unflagged "light" shot is host-dependent.
- --force-dark-mode: true|#1f1e1b (indistinguishable from the OS default on a Dark host; not shown to force anything).
- --blink-settings=preferredColorScheme=0: true|#1f1e1b (dark)
- --blink-settings=preferredColorScheme=1: false|#f7f6f2 (LIGHT, overrides the Dark OS)
- --screenshot with preferredColorScheme=0 / =1: rc=0 and PNG WRITTEN in both (unlike scriptEnabled=false); hashes differ (262b2f8b1b12 vs 135f95e75120); the =0 PNG equals the --force-dark-mode PNG byte-for-byte.
Untested: behaviour on a host whose OS appearance is Light (the explicit =1/=0 setting is the host-independent choice).
