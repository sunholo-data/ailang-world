# Mission Dashboard — World (snapshot 2026-09-25, iteration 187)

- **State**: **row 40 card half LANDED** (`w-a2a-session-projection`, clause 6). PR #146 was squash-merged as **`8bb8502`**, and remote CI is **green on the merge** (2/2, Gate 3b SHA-pinned). World now serves a session-scoped A2A agent card at `GET /.well-known/agent.json` and a fail-closed `/a2a/`, both built on row 39's Bearer session boundary.
  - `authority.ResolveContext` bounds session resolution, so a timeout is no longer reported as an unknown credential.
  - The upstream `serveapi/protocol` package is admitted by one package-path allowlist line.
- **What it is not yet**: invocation. Nothing in World invokes a registered transition server-side, and nothing publishes the transition registry in production, so a live daemon serves a zero-skills card. Every authorized `/a2a/` call gets a constant "not available" error. Split to rows **106** (coordinator) and **107** (publisher).
- **Quality**: planner prototype 30/30 mutations killed; executor re-drill 15/15. Judged **PASS 93/100, zero blocking** (sonnet, own worktree). Its 4 surviving mutations were closed test-only, each by a single load-bearing test; r2 judge **98/100**.
- **Design gate**: r1 blocked 3/3; astra's objection was upheld (a context-free resolve on a single-connection store) and became P6.A-CTX. r2 blocked on one surface, closed with the reviewers' fixes applied verbatim.
- **Upstream finding**: `sunholo-data/ailang#885` was closed "Delivered", but `serveapi/protocol` is blob-identical between v0.33.2 and v0.42.0. The MCP child (row **108**) stays blocked; a re-open was requested on the issue and via mission-control.
- **Next**: 27, **93** (clause-4 floor run), 106 (needs a design doc), 107, 99, 100, 96, 97, 103, 104.
- **Parked for Mark**: **nothing — ledger 23 rows, ZERO OPEN.**
- **Cadence/routing**:
  - **Roles:** controller `claude-opus-5-5`; designer pi deepseek (astra skipped, codex over ration); planner pi kimi (Ollama 429 mid-run, continued on OpenRouter); executor opus (end of chain); evaluator sonnet.
  - **Spend:** **$3.18 metered**, $2.63 of it the OpenRouter planner continuation.
  - **Ration:** codex, Ollama and OpenRouter were all over ration by the end of the fire.
- **Maintenance for an attended session**: row 101 — `gofmt -l` is red at base on 2 `host/store` files (the sprint fixed the third), and no CI step checks it.
