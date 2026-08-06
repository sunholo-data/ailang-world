# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-06, iteration 55 · dev @ `1761a9c` · CI green both jobs (SHA-addressed)

## In flight

- **Item 8 `w-self-mod-vertical` — `AC12` REPAIRED** (PR #45 → `1761a9c`). The World-boundary gate's
  loopback exception was true of **one** group and a single shared list gave it to **all three**:
  bare `net/http` blank-imported into `host/store` passed the gate green. Now per-group.
- **`[NEXT]` is still `SM.B2a`** — brokered publish handler, de-ambient credential, typed
  indeterminate (~780 LOC). Gated on nothing, and now lands against a gate with teeth.
- **NEW item 10 `w-boundary-gate-tree-mutation`** — the same gate mutates *other packages'*
  production sources in the live tree while `go test ./...` builds them concurrently. Measured on
  pristine `dev`. Consider landing **before** SM.B2a, which lengthens the broker suite and widens
  the window.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite. Unchanged.

## Latest — the guard against network had a hole in two of its three groups

- **Found at the SM.B2a boundary, before the network code it guards exists.** Every protected
  group's mutation was `net/http/httputil`, which the gate *does* catch — the mutation was shaped to
  the check, not the threat. Bare `net/http` → `host/store` **rc=0 PASS**, → `host/replay` **rc=0
  PASS**; both now RED, `cmd/ailang-worldd` still permitted, pristine still green.
- **The exemption was unforced**, not just unnecessary: baseline `net/http` in each closure is
  `host/store` **0**, `host/replay` **0**, `cmd/ailang-worldd` **1**. Only one group ever needed it.
- **Second finding, pre-existing and unrelated to the pick.** That same gate rewrites
  `main.go`/`store.go`/`replay.go` in place; a concurrent reader on pristine `dev` observed all
  three in a **non-compiling** state (5 samples of 90; control **0/200** with the gate idle). It
  red-lit `TestCLIRealSubprocessEpisode` once here. Latent in CI: **0/8** full-suite runs.

## Loop · cost · asks

- launchd `mission-world`; controller `claude-opus-5`. Designer/planner/executor/evaluator **not
  fired** — a controller-sized measured repair; rotation unchanged at `codex:gpt-5.6-sol`.
  Verify profile `ailang-code`; AILANG pinned **v0.30.0**. Issue **#32**.
- **`metered=$0.00`** vs the $5 ceiling — controller-only, every role on a quota bucket.
- **Parked on Mark: NONE.** `8/OD-2` open, non-blocking. FYI not blocking: item 9's human-gated half
  (pin CI job 1 vs keep tracking `latest`).
