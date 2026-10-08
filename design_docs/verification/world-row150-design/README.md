# Row150 design probe (2026-10-08)

Instrument-only Chrome control, not product validation. Two calls used system Chrome 154.0.8037.98 with separate profiles, original HOME, `--headless=new --no-first-run --disable-background-networking --hide-scrollbars --force-device-scale-factor=1 --window-size=1280,900 --user-data-dir=<isolated-profile> --screenshot=<controlN.png> file://<absolute-chrome-probe.html>`. Python `subprocess.run(..., timeout=30)` bounded each. Both wrote PNGs but did not exit by 30s. Profiles were removed; images and exact log retained. PNG-byte equality does not discharge the product capture criterion or timeout/cleanup criterion.

Also measured at base94fc5ba4: pinned `ailang check world/types.ail` exits0 (released v0.41.0 /24ee108); `go test ./cmd/ailang-worldd -run TestWhy -count=1` with110s alarm passes (package0.455s). These support code premises, not future demo implementation.
