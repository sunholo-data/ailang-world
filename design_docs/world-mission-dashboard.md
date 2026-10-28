# World mission dashboard — iteration 313 (2026-10-28)

- **This fire:** row 150 M8a BUILT (PR #400): reproduction docs, captions, release manifest schema and `marketingcheck docs`.
- **What landed:** QUICKSTART §13 (attended capture runbook, all thirteen `capture_marketing.sh` flags, the still-mode preview, the release check, P0/P3 and the open criteria); the website page `getting-started/marketing-capture` (six `useBaseUrl` images); `captions.md`; README Screenshots links; `docs/marketing-capture-manifest.schema.json` (no hash).
- **`marketingcheck docs`:** wiring checks against `Shots`; release mode (`--assets … --capture-manifest … --require-assets`) checks the PNG SHA-256 values against the attended capture manifest.
- **Measured first:** a markdown image to an absent PNG fails the Docusaurus build; a raw `<img>` and `pathname://` drop the `/ailang-world/` base; `useBaseUrl` works. The plan's attended block lacked three required capture flags.
- **Judges:** sonnet r1 PASS 82 → controller fix (the page said all evidence was demo-authored; only CLAIMED is — the judge missed it; plus disclosure gates F1, F6, F7, F8) → sonnet r2 PASS 89, zero blocking.
- **Roles:** executor haiku (codex over ration), evaluator sonnet ×2; no designer or planner needed.
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81): 156 (closed) → 150 → 151 → v1.0.0.
- **Row 150:** 22/23 executor milestones built (M8b left); R1 measured; P0, P3 (controller) and A1 (attended) remain.
- **Next:** the C8′(iii)/P0 store-dependent-geometry amendment (doc r3, designer + quorum) → M8b.
- **Open before P0:** phone C8′(iii) (scrollbar thumb 324–338 px across seeds) and P0's ±2 px for wrapping packet paragraphs.
- **For Mark at A1 (N2):** the release-screenshots page is in the sidebar before promotion, so a docs deploy before A1 shows six broken images. Deploy after A1, or mark the page `unlisted` until then.
- **Carried:** judge non-blocking F2–F5, F9 and N1 (negative arms, parse-leg kill, ungated new sentences); iteration 311/312 carries.
- **Instrument note:** build Docusaurus probes with a private `node_modules` (no `.cache`); a symlinked one shares webpack's cache and fails every docs page.
- **Parked on Mark:** nothing — the ledger has 0 OPEN rows.
- **Cost:** metered $0.00; Anthropic quota only.
- **Record/evidence:** world-mission-log.md entry 313; verification/world-iter313/.
