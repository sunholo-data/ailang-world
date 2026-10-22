# World surface, working copy v3: a design fixture (row 188, D-WORLD-88)

This is the Tauri prototype Mark judged on 2026-10-21. After v3 he said "starting to look good … we will iterate more im sure but this is good to build on top of". It is a **design fixture**, like [`../world-os-surface.dc.html`](../world-os-surface.dc.html): the visual and interaction reference that the re-cut S1's N2 (the app shell at `/surface`) and N3 (the shell in `desktop/`) port from. It is not product code, and no test or build reads it.

## How it was made
- **v1:** a same-data spike, Tauri (Claude) against SwiftUI (Codex). D-WORLD-88 chose Tauri.
- **v2:** a port of the approved artifact's own markup, inline styles, layout maths and themes, compared screenshot to screenshot against a local render of the artifact.
- **v3:** timeline navigation (pan, five zoom levels, an overview strip, Now) and a human layer: step titles, the agents' own summaries, brane names, goals and progress, a Today line, and hashes only on hover.

## What is here
- `dist/`: the interface. `index.html`, `app.js` and `app.css` are plain HTML, CSS and JS with no libraries. The fonts are OFL. Their licence texts are in `dist/fonts/`, but the font files themselves are **fetched, not tracked**, because the repo's tracked-binary gate keeps binaries out of git: `./fetch-fonts.sh` downloads them from google/fonts at pinned commits and refuses any sha256 mismatch.
  - `spike-data.json` is **frozen real data** from 2026-10-21. `data.js` is the same data as a script, for when `fetch` of a bundled file is refused.
  - The data's `provenance` block names the source of every field and lists what the records do not hold: ran-on, mutants caught, waiting steps, agent views.
- `src/`, `build.rs`, `Cargo.toml`, `Cargo.lock`, `tauri.conf.json`, `capabilities/`, `icons/`: the Tauri v2 shell, which embeds `dist/`.
- `make-app.sh`: builds `World.app` by hand, with no tauri-cli: `cargo build --release`, then an Info.plist, an icon and an ad-hoc signature.
- **No screenshots or icon** are tracked, for the same reason. To see it, open it; `make-app.sh` uses the default icon when `icons/icon.png` is absent.
- `NOTES.md`: the v3 notes, covering navigation, the human layer, and what is still mechanical and why.

## Open it
- **Browser:** run `./fetch-fonts.sh`, serve `dist/` from any static server on loopback, then open `/`.
- **Native:** run `./fetch-fonts.sh`, then `./make-app.sh` (needs `cargo`; no Xcode), then `open World.app`.

## What N2 changes
N2 embeds `dist/` into `host/surfacepage`, served at `/surface`. It swaps `spike-data.json` for the live routes from N1 (`/surface/model.json`, `/surface/stack.json`, which use the same shapes), self-hosts the fonts under `/surface`'s own CSP, and adds live refresh. Fields the records do not hold stay absent, never invented.
