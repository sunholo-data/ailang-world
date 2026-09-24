# M4 live timed walk — raw transcript (w-prove-1-0-phase-a)

Run 2026-09-24 against a scratch store (`.walk-scratch/walk.db`, untracked, never committed) in
worktree `.wt-world-iter182` at `ee03ab4` + the M4 working-tree changes. Every block below is
the command (prefixed `$ `) followed by its stdout+stderr **unedited**, then `[rc=N]`. Commands
were run through a small logging helper (`.walk-scratch/r.sh`: `cd` to the worktree, `export
AILANG_BIN=$HOME/.pinned-ailang/ailang`, `unset AILANG_REGISTRY_API_KEY`, then `bash -c` the
command). Lines starting `---` or `#` are the executor's annotations, written into the log at
the moment of the run.

## Deviations from the plan's commands (plan §5 M4), measured

- **D-1 — `--addr` needs a scheme and is not a `serve` flag.** The plan writes
  `--addr 127.0.0.1:7654`; the client builds `base + path`, so a scheme-less value fails with
  `first path segment in URL cannot contain colon`. Actual: `--addr http://127.0.0.1:7654`
  (the default is `http://127.0.0.1:7644`). The plan's serve line also puts `--addr` before
  `serve`; `main.go` refuses that (`--addr is a client flag and is not valid for 'serve'`).
  Actual serve line: `ailang-worldd serve --db … --bind 127.0.0.1:7654 --ailang-bin …`. All
  other verb/flag names (`health`, `head`, `log range --from`, `object get <ref> --payload`,
  `commit --file --session`, `session mint --db --episode --grant --ttl --out`) match the plan.
- **D-2 — `AILANG_REGISTRY_API_KEY` is exported in the rig shell.** `ailang-worldd`'s startup
  refusal (w-self-mod-vertical Decision 4) fires ahead of flag parsing for every verb. Per the
  refusal message's own instruction, it is unset for every `ailang-worldd` process (value never
  printed). Not a transport block: the refusal is the binary's intended hygiene check.
- **D-3 — mint cannot run while `serve` holds the db.** `session mint --db` opens the store as
  a writer; the single-writer lock refuses it while `serve` is up. The plan's order (serve →
  mint → commit) is therefore impossible as written; actual order: serve (pin check) → stop →
  mint → restart serve → pin re-check → commit. The tty fence itself was satisfied by the
  plan's `pty.spawn` wrapper on the first try (PF-2 did NOT fire). Note `pty.spawn` returns
  rc=0 even when the child fails — the `TOKEN-OK` check is the real assertion.
- **D-4 — the daemon prints no drain line.** A clean D7 drain writes nothing to stderr, so the
  second serve was wrapped (`wait $p; echo "serve exited rc=$?"`) to record its exit status.
- **D-5 — `gen_fixtures.py --out` default (fixed in M4).** The default was `"commits"` relative
  to CWD, so the plan's G2 command (no `--out`, run from the repo root) failed with
  `check: expected 4 fixtures, found 0`. The default is now the `commits/` directory next to
  the script. Proof below (E-fix).

## E-fix — gen_fixtures.py `--out` default

Before the fix (from the repo root): `check: expected 4 fixtures, found 0` / `check: 1 error(s)` / rc=1.
After:

```
$ python3 design_docs/verification/w-1-0-value-demonstration/gen_fixtures.py --check   # from repo root, no --out
check: OK — 4 fixture(s), schema + content-hashes + chain + interpreter asserted against /Users/voightkampff/.pinned-ailang/ailang
[rc=0]
$ python3 design_docs/verification/w-1-0-value-demonstration/gen_fixtures.py --out .walk-scratch/regen
wrote .walk-scratch/regen/q1-iter171-index-row.json
wrote .walk-scratch/regen/q2-iter154-unfenced-pi.json
wrote .walk-scratch/regen/q3-iter155-faithfulness-proof.json
wrote .walk-scratch/regen/q4-iter181-ci-bench-401.json
self-check: 0 error(s); interpreter pin sha256:1a67b01468584...
[rc=0]
$ diff -r .walk-scratch/regen design_docs/verification/w-1-0-value-demonstration/commits && echo BYTE-IDENTICAL
BYTE-IDENTICAL
[rc=0]
$ git diff --stat -- design_docs/verification/w-1-0-value-demonstration/commits
[rc=0]
```

## Setup (untimed; wall-clock recorded as context)

```
$ date -u +%Y-%m-%dT%H:%M:%SZ   # setup start
2026-09-24T09:34:47Z
[rc=0]
$ go build -o .walk-scratch/ailang-worldd ./cmd/ailang-worldd
[rc=0]
$ rm -f .walk-scratch/walk.db*; ls .walk-scratch/walk.db 2>&1
ls: .walk-scratch/walk.db: No such file or directory
[rc=1]
$ .walk-scratch/ailang-worldd --addr 127.0.0.1:7654 serve --db .walk-scratch/walk.db --bind 127.0.0.1:7654 --ailang-bin "$AILANG_BIN"   # plan literal
ailang-worldd: broker: AILANG_REGISTRY_API_KEY is set in the process environment: move it to a mode-0600 file outside the working tree and unset it (see design_docs/planned/w-self-mod-vertical.md Decision 4)
[rc=2]
$ .walk-scratch/ailang-worldd serve --db .walk-scratch/walk.db --bind 127.0.0.1:7654 --ailang-bin "$AILANG_BIN" > .walk-scratch/serve.log 2>&1 &   # actual (background)
serve pid 97813

--- serve.log of the attempt above (the daemon exited at startup):
ailang-worldd: broker: AILANG_REGISTRY_API_KEY is set in the process environment: move it to a mode-0600 file outside the working tree and unset it (see design_docs/planned/w-self-mod-vertical.md Decision 4)
--- DEVIATION D-2: the rig shell exports AILANG_REGISTRY_API_KEY (value never printed); ailang-worldd's startup refusal (w-self-mod-vertical Decision 4) fires for EVERY verb. Per the refusal message, it is unset for every ailang-worldd process below (`unset AILANG_REGISTRY_API_KEY` in the helper / `env -u` for serve).
$ env -u AILANG_REGISTRY_API_KEY .walk-scratch/ailang-worldd serve --db .walk-scratch/walk.db --bind 127.0.0.1:7654 --ailang-bin "$AILANG_BIN" > .walk-scratch/serve.log 2>&1 &
serve pid 98061
$ deadline=$(( $(date +%s) + 15 )); until .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 health >/dev/null 2>&1; do [ $(date +%s) -ge $deadline ] && { echo HEALTH-TIMEOUT; break; }; sleep 0.2; done; echo polled
polled
[rc=0]
$ .walk-scratch/ailang-worldd --addr 127.0.0.1:7654 health   # plan literal (no scheme)
ailang-worldd: build GET 127.0.0.1:7654/v1/health: parse "127.0.0.1:7654/v1/health": first path segment in URL cannot contain colon
[rc=1]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 health
{"status":"ok","daemon_version":"0.1.0","db_path":".walk-scratch/walk.db","interpreter_ref":"sha256:1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f","interpreter_version":"AILANG v0.41.0\nCommit: 24ee108\nFull:   24ee1088776e21cd06a3781ed18e77f40be06db3\nBuilt:  2026-09-21T15:50:02Z\n\nThe AI-First Programming Language\nCopyright (c) 2025-2026\n"}
[rc=0]
$ cat .walk-scratch/serve.log
ailang-worldd listening on http://127.0.0.1:7654
[rc=0]
$ .walk-scratch/ailang-worldd --addr 127.0.0.1:7654 serve --db .walk-scratch/walk-unused.db --bind 127.0.0.1:7655 --ailang-bin "$AILANG_BIN"   # plan literal shape, key unset: --addr before serve is refused (DEVIATION D-1)
ailang-worldd: --addr is a client flag and is not valid for 'serve'; use --bind to choose the listen address
[rc=1]
$ ls .walk-scratch/walk-unused.db 2>&1
ls: .walk-scratch/walk-unused.db: No such file or directory
[rc=1]
$ test "$(.walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 health | python3 -c 'import json,sys; print(json.load(sys.stdin)["interpreter_ref"])')" = "sha256:$(shasum -a 256 "$AILANG_BIN" | awk '{print $1}')" && echo REPLAY-PIN-OK
REPLAY-PIN-OK
[rc=0]
$ printf 'y\n' | python3 -c "import pty,sys; pty.spawn(sys.argv[1:])" .walk-scratch/ailang-worldd session mint --db .walk-scratch/walk.db --episode w-1-0-value-demo --grant fs.read=/tmp:1 --ttl 3600 --out .walk-scratch/token
y
Confirm mint for episode w-1-0-value-demo (1 grant(s), expiry +3600s) with a session credential? [y/N] ailang-worldd session mint: cannot open world store ".walk-scratch/walk.db": store: another process already holds the writer lock for "/Users/voightkampff/dev/sunholo-data/.wt-world-iter182/.walk-scratch/walk.db" (lock file "/Users/voightkampff/dev/sunholo-data/.wt-world-iter182/.walk-scratch/walk.db.writer.lock"); only one writer may be active per database
[rc=0]
$ test "$(wc -c < .walk-scratch/token | tr -d " ")" = "64" && echo TOKEN-OK; stat -f "%Sp" .walk-scratch/token
bash: .walk-scratch/token: No such file or directory
stat: .walk-scratch/token: stat: No such file or directory
[rc=1]
--- DEVIATION D-3: mint opens --db as a WRITER; the store's single-writer lock refuses it while serve holds the same db. (pty.spawn returned rc=0 — its exit status is not the child's; the TOKEN check is the real fence.) Order changed: stop serve -> mint -> restart serve.
$ kill -TERM $(cat .walk-scratch/serve.pid)
$ pid=$(cat .walk-scratch/serve.pid); deadline=$(( $(date +%s) + 15 )); while kill -0 $pid 2>/dev/null; do [ $(date +%s) -ge $deadline ] && { echo STOP-TIMEOUT; break; }; sleep 0.2; done; cat .walk-scratch/serve.log
ailang-worldd listening on http://127.0.0.1:7654
[rc=0]
$ printf 'y\n' | python3 -c "import pty,sys; pty.spawn(sys.argv[1:])" .walk-scratch/ailang-worldd session mint --db .walk-scratch/walk.db --episode w-1-0-value-demo --grant fs.read=/tmp:1 --ttl 3600 --out .walk-scratch/token
y
Confirm mint for episode w-1-0-value-demo (1 grant(s), expiry +3600s) with a session credential? [y/N] minted session credential for episode w-1-0-value-demo: 1 grant(s), expires epoch 1790246119, written to .walk-scratch/token
ailang-worldd session mint: credential_id=0f3ad65fa1c48a3fc4aed9cf55b4ca20bc551241ef3dd4e3f15f64b493067070 (episode w-1-0-value-demo, 1 grant(s), expires 1790246119)
[rc=0]
$ test "$(wc -c < .walk-scratch/token | tr -d " ")" = "64" && echo TOKEN-OK; stat -f "%Sp" .walk-scratch/token
TOKEN-OK
-rw-------
[rc=0]
--- note: the first serve exited on SIGTERM with no stderr line (a clean D7 drain prints nothing). The restart below wraps serve so its exit status is appended to serve.log.
$ nohup bash -c 'env -u AILANG_REGISTRY_API_KEY .walk-scratch/ailang-worldd serve --db .walk-scratch/walk.db --bind 127.0.0.1:7654 --ailang-bin "$AILANG_BIN" & p=$!; echo $p > .walk-scratch/serve.pid; wait $p; echo "serve exited rc=$?"' > .walk-scratch/serve.log 2>&1 &
$ deadline=$(( $(date +%s) + 15 )); until .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 health >/dev/null 2>&1; do [ $(date +%s) -ge $deadline ] && { echo HEALTH-TIMEOUT; break; }; sleep 0.2; done; echo polled; cat .walk-scratch/serve.log
polled
ailang-worldd listening on http://127.0.0.1:7654
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 health
{"status":"ok","daemon_version":"0.1.0","db_path":".walk-scratch/walk.db","interpreter_ref":"sha256:1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f","interpreter_version":"AILANG v0.41.0\nCommit: 24ee108\nFull:   24ee1088776e21cd06a3781ed18e77f40be06db3\nBuilt:  2026-09-21T15:50:02Z\n\nThe AI-First Programming Language\nCopyright (c) 2025-2026\n"}
[rc=0]
$ test "$(.walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 health | python3 -c 'import json,sys; print(json.load(sys.stdin)["interpreter_ref"])')" = "sha256:$(shasum -a 256 "$AILANG_BIN" | awk '{print $1}')" && echo REPLAY-PIN-OK
REPLAY-PIN-OK
[rc=0]
$ for f in design_docs/verification/w-1-0-value-demonstration/commits/q*.json; do echo "== $f"; .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 commit --file "$f" --session "$(cat .walk-scratch/token)"; echo "[commit rc=$?]"; done
== design_docs/verification/w-1-0-value-demonstration/commits/q1-iter171-index-row.json
{"selectedHead":"sha256:a3ae98dcd018dc0bcd18f27f78f8d53bfa4f9ade3aa2521bb4c3858f457e75eb"}
[commit rc=0]
== design_docs/verification/w-1-0-value-demonstration/commits/q2-iter154-unfenced-pi.json
{"selectedHead":"sha256:22567594c4c250297512b0b4a2f2a5c598b23aa436ec1d5e6ef6aff6af15e0ee"}
[commit rc=0]
== design_docs/verification/w-1-0-value-demonstration/commits/q3-iter155-faithfulness-proof.json
{"selectedHead":"sha256:749fc1728b737e8764d5ce09a70ed51a7693ab184f60c6190bccde1026534457"}
[commit rc=0]
== design_docs/verification/w-1-0-value-demonstration/commits/q4-iter181-ci-bench-401.json
{"selectedHead":"sha256:ef2d72ebf335c265f47efcaf29a1f469349e593700d9d1f9e0a6fa4064f135aa"}
[commit rc=0]
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 head
sha256:ef2d72ebf335c265f47efcaf29a1f469349e593700d9d1f9e0a6fa4064f135aa
[rc=0]
$ date -u +%Y-%m-%dT%H:%M:%SZ   # setup end
2026-09-24T09:35:32Z
[rc=0]
```

## Q1 (timed)

```
# Q1 — Why did iteration 171 have a full log entry and no index row?
$ T0=$(date +%s); echo $T0 > .walk-scratch/q1.t0; date -u +%Y-%m-%dT%H:%M:%SZ   # START clock: "Why did iteration 171 have a full log entry and no index row?"
2026-09-24T09:35:37Z
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 log range --from 0
{"items":[{"header":{"entryIndex":0,"semanticsEpoch":1,"transitionFn":"sha256:5e670bae3ff5816cd7f6faed983dd61e2942d5b0268a3c66a53a9ea6bc1d3f1d","interpreter":"sha256:1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f","prevEntryHash":"sha256:ee2957a18e7806d9854f199f6533effdd3ccab4b04cf9126cc054818cfe26900","writtenBy":"w-prove-1-0-phase-a"},"entryHash":"sha256:cd8180864257c0510a64981eed25e913a35f278f81d9a9247d6258fd69a44ee5","transitionRef":"sha256:c07048332f5cd09d5d632a1059941f56ba1aa2914846f9be1c5228d6979d6206"},{"header":{"entryIndex":1,"semanticsEpoch":1,"transitionFn":"sha256:5e670bae3ff5816cd7f6faed983dd61e2942d5b0268a3c66a53a9ea6bc1d3f1d","interpreter":"sha256:1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f","prevEntryHash":"sha256:cd8180864257c0510a64981eed25e913a35f278f81d9a9247d6258fd69a44ee5","writtenBy":"w-prove-1-0-phase-a"},"entryHash":"sha256:8549b68c5dd3074c9ed9ebf5fc908fc0ee8e7f7da7ab7d22c164e8a88e55a446","transitionRef":"sha256:290a22fdcc72a313cb22245ff95bbda95d83625a884954fb47d15579bf9703ea"},{"header":{"entryIndex":2,"semanticsEpoch":1,"transitionFn":"sha256:5e670bae3ff5816cd7f6faed983dd61e2942d5b0268a3c66a53a9ea6bc1d3f1d","interpreter":"sha256:1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f","prevEntryHash":"sha256:8549b68c5dd3074c9ed9ebf5fc908fc0ee8e7f7da7ab7d22c164e8a88e55a446","writtenBy":"w-prove-1-0-phase-a"},"entryHash":"sha256:af33ffe45e6a0dc099020337ebbce65f762564684d719161df2e49008c8d9e28","transitionRef":"sha256:cacf22c211a74d8c28c94c4dafecf2d2f78e4ee41854b37290889fe7944b40dd"},{"header":{"entryIndex":3,"semanticsEpoch":1,"transitionFn":"sha256:5e670bae3ff5816cd7f6faed983dd61e2942d5b0268a3c66a53a9ea6bc1d3f1d","interpreter":"sha256:1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f","prevEntryHash":"sha256:af33ffe45e6a0dc099020337ebbce65f762564684d719161df2e49008c8d9e28","writtenBy":"w-prove-1-0-phase-a"},"entryHash":"sha256:23a8b9bb7e62b98c11543a9fbb66031fd791e458d033584a6e7055cfce29070f","transitionRef":"sha256:722751941829c9165764b76d78d435f83ea7ef3aeadd573c7b29d4fdaf629d2b"}]}
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:c07048332f5cd09d5d632a1059941f56ba1aa2914846f9be1c5228d6979d6206
{"hash":"sha256:c07048332f5cd09d5d632a1059941f56ba1aa2914846f9be1c5228d6979d6206","interfaceHash":"sha256:c604ef41412830eae9cead714250e42c65accb445d92888387c25da902fd6a42","semanticId":"world/mission/incident/iter171-index-row","provenance":"w-prove-1-0-phase-a"}
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:c07048332f5cd09d5d632a1059941f56ba1aa2914846f9be1c5228d6979d6206 --payload | python3 -c 'import json,sys,base64; print(base64.b64decode(json.load(sys.stdin)["payload"]).decode())'
{"question":"Why did iteration 171 have a full log entry and no index row?","answer":"The pinned v0.30.0 binary lacks the mission command group: `ailang.v0.30.0 mission rotate-log world --keep 31` prints `Error: unknown command 'mission'`; the index hole is visible in git history of world-mission-index.md (the row exists today because it was added manually after the iter-172 rule).","diagnosed_at":"2026-09-08","sources":["sha256:006e5c6db28078565852a45ddd2f4b4977ef887e74ff1870a31c3b44c76493f1","sha256:b419a0de9c12ba52db42be7ec3be9c88c3e2d1cd61521ef278e164cd78c42ef0","sha256:28f698b35c67e8399e305d313ec177b6a5a52f484e69ff629e2448d10ab76c32"]}
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:006e5c6db28078565852a45ddd2f4b4977ef887e74ff1870a31c3b44c76493f1 --payload | python3 -c 'import json,sys,base64; o=json.load(sys.stdin); print(o["semanticId"]); print(base64.b64decode(o["payload"]).decode())'
world/mission/evidence/iter171-index-row/0
{"kind":"binary","ref":"$HOME/.pinned-ailang/ailang.v0.30.0","check":"ailang.v0.30.0 mission rotate-log world --keep 31 => Error: unknown command 'mission'","excerpt":"The pinned v0.30.0 binary has no `ailang mission` command at all: `~/.pinned-ailang/ailang mission rotate-log world --keep 31` prints **`Error: unknown command 'mission'`**."}
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:b419a0de9c12ba52db42be7ec3be9c88c3e2d1cd61521ef278e164cd78c42ef0 --payload | python3 -c 'import json,sys,base64; o=json.load(sys.stdin); print(o["semanticId"]); print(base64.b64decode(o["payload"]).decode())'
world/mission/evidence/iter171-index-row/1
{"kind":"index","ref":"world-mission-index.md (git history)","check":"grep -c '^| 171 |' world-mission-index.md => 1 NOW (the manual repair, not iter-171's Gate 4)","excerpt":"iteration **171** had a full log entry and **no index row** (`grep -c '^| 171 |'` -> **0**, control `170` -> **1**)."}
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:28f698b35c67e8399e305d313ec177b6a5a52f484e69ff629e2448d10ab76c32 --payload | python3 -c 'import json,sys,base64; o=json.load(sys.stdin); print(o["semanticId"]); print(base64.b64decode(o["payload"]).decode())'
world/mission/evidence/iter171-index-row/2
{"kind":"record","ref":"world-mission.md ~lines 82-101; log ## 172","check":"read the `mission`-command pin-version-gap bullet","excerpt":"the failure is SILENT in the direction that matters: Gate 2's first instruction is to grep the index before picking, so a missing row does not produce an error."}
[rc=0]
# assembled answer (incident.answer): The pinned v0.30.0 binary lacks the mission command group ... the index hole is visible in git history of world-mission-index.md
$ "$HOME/.pinned-ailang/ailang.v0.30.0" mission rotate-log world --keep 31 2>&1 | grep -F "unknown command 'mission'"
Error: unknown command 'mission'
[rc=0]
$ "$HOME/.pinned-ailang/ailang.v0.30.0" --version 2>&1 | head -2   # positive control: this IS the v0.30.0 binary
2026/09/24 11:35:50 Observatory: 586MB (DB=581MB WAL=5MB) — running retention cleanup
2026/09/24 11:35:50 Observatory cleanup: deleted spans=0 summaries=0 metrics=0 chat=0 tools=0 (size 586MB → 581MB)
[rc=0]
$ git log --oneline -S'| 171 |' -- design_docs/world-mission-index.md | tail -3
0d31638 docs(mission): iteration 172 record — row 72 landed, the count is a reading, and a byte-level fixture lesson
[rc=0]
$ "$HOME/.pinned-ailang/ailang.v0.30.0" --version 2>/dev/null | grep -m1 -i 'v0\.'   # positive control: this IS the v0.30.0 binary
AILANG v0.30.0
[rc=0]
$ git show --stat --oneline 0d31638 | head -8; git show 0d31638 -- design_docs/world-mission-index.md | grep -E '^\+\| 17[12] \|' | cut -c1-120
0d31638 docs(mission): iteration 172 record — row 72 landed, the count is a reading, and a byte-level fixture lesson
 .../w-queue-closed-count-census-sprint-plan.md     |  0
 .../w-queue-closed-count-census.md                 |  0
 design_docs/world-mission-dashboard.md             | 69 +++++++++-------------
 design_docs/world-mission-index.md                 | 13 ++++
 design_docs/world-mission-log.md                   | 42 +++++++++++++
 design_docs/world-mission-status-archive.md        |  2 +
 design_docs/world-mission.md                       | 32 +++++++++-
+| 172 | 2026-09-08 | row 72 LANDED: `scripts/queue_census.sh` — the loop's own closed-row count becomes a READING wit
+| 171 | 2026-09-08 | row 71 LANDED: critical mission state never lives where the OS wipes; all four roles on pi lanes, 
[rc=0]
$ git show 0d31638^:design_docs/world-mission-index.md | grep -c '^| 171 |'; git show 0d31638^:design_docs/world-mission-index.md | grep -c '^| 170 |'   # the hole: parent of the repair has 0 rows for 171, control 170 -> 1
0
1
[rc=0]
$ T1=$(date +%s); T0=$(cat .walk-scratch/q1.t0); date -u +%Y-%m-%dT%H:%M:%SZ; echo "duration_s=$((T1-T0))"   # STOP clock
2026-09-24T09:35:57Z
duration_s=20
[rc=0]
```

## Q2 (timed)

```
# Q2 — Why did three designer runs execute unfenced while the rulebook said they were sandboxed?
$ date +%s > .walk-scratch/q2.t0; date -u +%Y-%m-%dT%H:%M:%SZ   # START clock: "Why did three designer runs execute unfenced while the rulebook said they were sandboxed?"
2026-09-24T09:36:03Z
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 log range --from 0 | python3 -c 'import json,sys; [print(e["header"]["entryIndex"], e["transitionRef"]) for e in json.load(sys.stdin)["items"]]'
0 sha256:c07048332f5cd09d5d632a1059941f56ba1aa2914846f9be1c5228d6979d6206
1 sha256:290a22fdcc72a313cb22245ff95bbda95d83625a884954fb47d15579bf9703ea
2 sha256:cacf22c211a74d8c28c94c4dafecf2d2f78e4ee41854b37290889fe7944b40dd
3 sha256:722751941829c9165764b76d78d435f83ea7ef3aeadd573c7b29d4fdaf629d2b
[rc=0]
$ for r in $(.walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 log range --from 0 | python3 -c 'import json,sys; [print(e["transitionRef"]) for e in json.load(sys.stdin)["items"]]'); do .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get $r; done
{"hash":"sha256:c07048332f5cd09d5d632a1059941f56ba1aa2914846f9be1c5228d6979d6206","interfaceHash":"sha256:c604ef41412830eae9cead714250e42c65accb445d92888387c25da902fd6a42","semanticId":"world/mission/incident/iter171-index-row","provenance":"w-prove-1-0-phase-a"}
{"hash":"sha256:290a22fdcc72a313cb22245ff95bbda95d83625a884954fb47d15579bf9703ea","interfaceHash":"sha256:258c0f772883db38aa9996e8709369b2d1d728286f6edfbc43fd9cad792c6ff3","semanticId":"world/mission/incident/iter154-unfenced-pi","provenance":"w-prove-1-0-phase-a"}
{"hash":"sha256:cacf22c211a74d8c28c94c4dafecf2d2f78e4ee41854b37290889fe7944b40dd","interfaceHash":"sha256:cca7a1ba281879dbf7cfc070cdbc6fa69d69339c5c05d15bebe73136434b0546","semanticId":"world/mission/incident/iter155-faithfulness-proof","provenance":"w-prove-1-0-phase-a"}
{"hash":"sha256:722751941829c9165764b76d78d435f83ea7ef3aeadd573c7b29d4fdaf629d2b","interfaceHash":"sha256:8dc8c03857020a11a4c39f31616c82de40a553548b1cfc4a6b4bb54d2537294f","semanticId":"world/mission/incident/iter181-ci-bench-401","provenance":"w-prove-1-0-phase-a"}
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:290a22fdcc72a313cb22245ff95bbda95d83625a884954fb47d15579bf9703ea --payload | python3 -c 'import json,sys,base64; print(base64.b64decode(json.load(sys.stdin)["payload"]).decode())'
{"question":"Why did three designer runs execute unfenced while the rulebook said they were sandboxed?","answer":"The fleet's scripts/mission_pi_run.sh invokes `pi --mode json --no-session --model \"$MODEL\" < \"$DIRECTIVE\"` with no -e flag and no PI_FENCE_ROOT anywhere in the file, so the shared skill's two sandbox extensions are never wired; issue sunholo-data/ailang#1043 is OPEN.","diagnosed_at":"2026-09-04","sources":["sha256:eb21c6282318d49279b77b15c09e3432ec41817ef45d65d2d83a3a68465db643","sha256:360bde0a5262e459636f283e4c87fd8efbd423ddcef10020fce37cad6c847fdf","sha256:d5ff45ac5536c4794650f0b8d5a9406cde0fea85bf6f27f21d88f1fa808976e2"]}
[rc=0]
$ for h in $(.walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:290a22fdcc72a313cb22245ff95bbda95d83625a884954fb47d15579bf9703ea --payload | python3 -c 'import json,sys,base64; [print(s) for s in json.loads(base64.b64decode(json.load(sys.stdin)["payload"]))["sources"]]'); do echo "-- $h"; .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get $h --payload | python3 -c 'import json,sys,base64; o=json.load(sys.stdin); print(o["semanticId"]); print(base64.b64decode(o["payload"]).decode())'; done
-- sha256:eb21c6282318d49279b77b15c09e3432ec41817ef45d65d2d83a3a68465db643
world/mission/evidence/iter154-unfenced-pi/0
{"kind":"script","ref":"scripts/mission_pi_run.sh (fleet, live via gh api)","check":"grep -n 'pi --mode|PI_FENCE_ROOT| -e ' => invocation matches, ZERO -e/PI_FENCE_ROOT hits","excerpt":"measured first-party at `scripts/mission_pi_run.sh:155`: the invocation is `pi --mode json --no-session --model \"$MODEL\" < \"$DIRECTIVE\"`, with **no `-e` flag at all** and no `PI_FENCE_ROOT`."}
-- sha256:360bde0a5262e459636f283e4c87fd8efbd423ddcef10020fce37cad6c847fdf
world/mission/evidence/iter154-unfenced-pi/1
{"kind":"recipe","ref":"shared-skill pi recipe (fleet-owned)","check":"recipe invocation block mandates -e sandbox/index.ts and -e worktree-fence.ts","excerpt":"The recipe's invocation block is explicit — `-e \"$REPO/tools/pi-extensions/sandbox/index.ts\" -e \"$REPO/tools/pi-extensions/worktree-fence.ts\"`."}
-- sha256:d5ff45ac5536c4794650f0b8d5a9406cde0fea85bf6f27f21d88f1fa808976e2
world/mission/evidence/iter154-unfenced-pi/2
{"kind":"issue","ref":"sunholo-data/ailang#1043","check":"gh issue view 1043 --repo sunholo-data/ailang --json state --jq .state => OPEN","excerpt":"the code half is filed upstream as [`ailang#1043`](https://github.com/sunholo-data/ailang/issues/1043)."}
[rc=0]
# assembled answer (incident.answer): mission_pi_run.sh invokes pi with no -e flag and no PI_FENCE_ROOT; the sandbox extensions are never wired; ailang#1043 OPEN
$ /opt/homebrew/bin/gh api repos/sunholo-data/ailang/contents/scripts/mission_pi_run.sh --jq .content | base64 -d > .walk-scratch/mission_pi_run.sh; wc -l < .walk-scratch/mission_pi_run.sh; echo known-present:$(grep -c 'pi --mode' .walk-scratch/mission_pi_run.sh) PI_FENCE_ROOT:$(grep -c 'PI_FENCE_ROOT' .walk-scratch/mission_pi_run.sh) space-e-space:$(grep -c ' -e ' .walk-scratch/mission_pi_run.sh); grep -n 'pi --mode\|PI_FENCE_ROOT\| -e ' .walk-scratch/mission_pi_run.sh
     279
known-present:1 PI_FENCE_ROOT:0 space-e-space:0
165:    pi --mode json --no-session --model "$MODEL" < "$DIRECTIVE" 2>"$ERR" |
[rc=0]
$ /opt/homebrew/bin/gh issue view 1043 --repo sunholo-data/ailang --json state --jq .state
OPEN
[rc=0]
$ echo "duration_s=$(( $(date +%s) - $(cat .walk-scratch/q2.t0) ))"; date -u +%Y-%m-%dT%H:%M:%SZ   # STOP clock
duration_s=13
2026-09-24T09:36:16Z
[rc=0]
```

## Q3 (timed)

```
# Q3 — Why did a green `shasum -c` faithfulness proof pass over a destroyed commit split?
$ date +%s > .walk-scratch/q3.t0; date -u +%Y-%m-%dT%H:%M:%SZ   # START clock: "Why did a green shasum -c faithfulness proof pass over a destroyed commit split?"
2026-09-24T09:36:21Z
[rc=0]
$ for r in $(.walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 log range --from 0 | python3 -c 'import json,sys; [print(e["transitionRef"]) for e in json.load(sys.stdin)["items"]]'); do .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get $r; done
{"hash":"sha256:c07048332f5cd09d5d632a1059941f56ba1aa2914846f9be1c5228d6979d6206","interfaceHash":"sha256:c604ef41412830eae9cead714250e42c65accb445d92888387c25da902fd6a42","semanticId":"world/mission/incident/iter171-index-row","provenance":"w-prove-1-0-phase-a"}
{"hash":"sha256:290a22fdcc72a313cb22245ff95bbda95d83625a884954fb47d15579bf9703ea","interfaceHash":"sha256:258c0f772883db38aa9996e8709369b2d1d728286f6edfbc43fd9cad792c6ff3","semanticId":"world/mission/incident/iter154-unfenced-pi","provenance":"w-prove-1-0-phase-a"}
{"hash":"sha256:cacf22c211a74d8c28c94c4dafecf2d2f78e4ee41854b37290889fe7944b40dd","interfaceHash":"sha256:cca7a1ba281879dbf7cfc070cdbc6fa69d69339c5c05d15bebe73136434b0546","semanticId":"world/mission/incident/iter155-faithfulness-proof","provenance":"w-prove-1-0-phase-a"}
{"hash":"sha256:722751941829c9165764b76d78d435f83ea7ef3aeadd573c7b29d4fdaf629d2b","interfaceHash":"sha256:8dc8c03857020a11a4c39f31616c82de40a553548b1cfc4a6b4bb54d2537294f","semanticId":"world/mission/incident/iter181-ci-bench-401","provenance":"w-prove-1-0-phase-a"}
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:cacf22c211a74d8c28c94c4dafecf2d2f78e4ee41854b37290889fe7944b40dd --payload | python3 -c 'import json,sys,base64; print(base64.b64decode(json.load(sys.stdin)["payload"]).decode())'
{"question":"Why did a green `shasum -c` faithfulness proof pass over a destroyed commit split?","answer":"The prescribed manifest covers the **final tree**, identical whether the split is correct or collapsed; `git add design_docs host` staged disk state, MS1 swallowed MS3's files, MS3 committed nothing, and `shasum -c` returned OK on every file.","diagnosed_at":"2026-09-05","sources":["sha256:be8bc3f149f18c4c03074d259332a08f298ebec4cddfa2189944507aec5ba301","sha256:854dd4c651ead8e1ac1fc1ac5fbf25dcca674feb60623b74b56cae7b07226650","sha256:a8c2a061b960c5a5137726515499104807e89291d91eaacc9e3a7048d51c60e0"]}
[rc=0]
$ for h in $(.walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:cacf22c211a74d8c28c94c4dafecf2d2f78e4ee41854b37290889fe7944b40dd --payload | python3 -c 'import json,sys,base64; [print(s) for s in json.loads(base64.b64decode(json.load(sys.stdin)["payload"]))["sources"]]'); do echo "-- $h"; .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get $h --payload | python3 -c 'import json,sys,base64; o=json.load(sys.stdin); print(o["semanticId"]); print(base64.b64decode(o["payload"]).decode())'; done
-- sha256:be8bc3f149f18c4c03074d259332a08f298ebec4cddfa2189944507aec5ba301
world/mission/evidence/iter155-faithfulness-proof/0
{"kind":"manifest","ref":"world-mission.md row 81","check":"grep -n \"manifest is over the FINAL TREE\" world-mission.md","excerpt":"**That manifest is over the FINAL TREE.** Bisectability is a property of the SPLIT, and the final tree is identical whether the split is correct or whether commit 1 swallowed everything and commits 2-3 are empty."}
-- sha256:854dd4c651ead8e1ac1fc1ac5fbf25dcca674feb60623b74b56cae7b07226650
world/mission/evidence/iter155-faithfulness-proof/1
{"kind":"log","ref":"world-mission-log.md ## 155","check":"git show --stat per commit; rebuilt row-59 split in git history","excerpt":"**My first commit reconstruction was wrong and I rebuilt it.** `git add design_docs host` staged whatever was on disk rather than the snapshot's named files, so MS1's commit swallowed MS3's two doc files."}
-- sha256:a8c2a061b960c5a5137726515499104807e89291d91eaacc9e3a7048d51c60e0
world/mission/evidence/iter155-faithfulness-proof/2
{"kind":"retry","ref":"world-mission-log.md ## 155","check":"rebuild boundary gate; zsh word-splitting trap","excerpt":"zsh does not word-split unquoted expansions, so all three paths arrived as ONE argument, `cp` failed, and the boundary gates still printed `vet=0 verifygate=0`."}
[rc=0]
# assembled answer (incident.answer): the prescribed manifest covers the FINAL TREE (identical for correct or collapsed split); git add design_docs host staged disk state; MS1 swallowed MS3's files; MS3 committed nothing; shasum -c OK on every file
$ grep -n "manifest is over the FINAL TREE" design_docs/world-mission.md | cut -c1-200
5931:81. **w-snapshot-reconstruction-faithfulness-proof-is-blind-to-the-split-it-sits-next-to** · clause-2 · **THE SHARED SKILL PRESCRIBES A sha256 FAITHFULNESS PROOF FOR MULTI-MILESTONE COMMIT RECO
[rc=0]
$ grep -n "commit reconstruction was wrong and I rebuilt it" design_docs/world-mission-log.md | cut -c1-200
1688:- **My first commit reconstruction was wrong and I rebuilt it.** `git add design_docs host` staged
[rc=0]
$ git log --oneline --since=2026-09-04 --until=2026-09-07 --grep='row 59' --all | head -10
c6456dc record(world) iter 155: row 59 lands and row 50 closes with it — the planner, not the quorum, found the doc committing its own thesis four more times
d353ef1 docs(standards): load-bearing criteria need a mutation, not a grep
5d84209 test(verifygate): add the fixture-shape gate, with its own red arms
fb8bc29 feat(verifygate): move the toolchain pins into a data-only fixture run.sh executes
5928453 docs(design): correct three measured defects in the row-59 design before its sprint
[rc=0]
$ git log --oneline --reverse 5928453..c6456dc | head -12
fb8bc29 feat(verifygate): move the toolchain pins into a data-only fixture run.sh executes
5d84209 test(verifygate): add the fixture-shape gate, with its own red arms
d353ef1 docs(standards): load-bearing criteria need a mutation, not a grep
c6456dc record(world) iter 155: row 59 lands and row 50 closes with it — the planner, not the quorum, found the doc committing its own thesis four more times
[rc=0]
$ for c in fb8bc29 5d84209 d353ef1; do git show --stat --format='%h %ad %s' --date=short $c | cat; echo; done
fb8bc29 2026-09-05 feat(verifygate): move the toolchain pins into a data-only fixture run.sh executes

 .../verification/w-race-gate-blindspot/run.sh      |  40 ++-
 .../w-race-gate-blindspot/toolchain_pins.conf      |   7 +
 host/verifygate/toolchain_pin_gate_test.go         | 268 +++++++++++++++++++--
 3 files changed, 294 insertions(+), 21 deletions(-)

5d84209 2026-09-05 test(verifygate): add the fixture-shape gate, with its own red arms

 host/verifygate/toolchain_pin_gate_test.go | 58 ++++++++++++++++++++++++++++++
 1 file changed, 58 insertions(+)

d353ef1 2026-09-05 docs(standards): load-bearing criteria need a mutation, not a grep

 design_docs/coding-standards.md                              |  8 ++++++++
 ...w-shell-assignment-parser-drops-an-indented-assignment.md | 12 ++++++++++--
 2 files changed, 18 insertions(+), 2 deletions(-)

[rc=0]
$ echo "duration_s=$(( $(date +%s) - $(cat .walk-scratch/q3.t0) ))"; date -u +%Y-%m-%dT%H:%M:%SZ   # STOP clock
duration_s=18
2026-09-24T09:36:39Z
[rc=0]
```

## Q4 (timed)

```
# Q4 — Why did PR #141's remote CI go red on `BenchmarkRESTCommit` when every local gate was green?
$ date +%s > .walk-scratch/q4.t0; date -u +%Y-%m-%dT%H:%M:%SZ   # START clock: "Why did PR #141s remote CI go red on BenchmarkRESTCommit when every local gate was green?"
2026-09-24T09:36:45Z
[rc=0]
$ for r in $(.walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 log range --from 0 | python3 -c 'import json,sys; [print(e["transitionRef"]) for e in json.load(sys.stdin)["items"]]'); do .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get $r; done
{"hash":"sha256:c07048332f5cd09d5d632a1059941f56ba1aa2914846f9be1c5228d6979d6206","interfaceHash":"sha256:c604ef41412830eae9cead714250e42c65accb445d92888387c25da902fd6a42","semanticId":"world/mission/incident/iter171-index-row","provenance":"w-prove-1-0-phase-a"}
{"hash":"sha256:290a22fdcc72a313cb22245ff95bbda95d83625a884954fb47d15579bf9703ea","interfaceHash":"sha256:258c0f772883db38aa9996e8709369b2d1d728286f6edfbc43fd9cad792c6ff3","semanticId":"world/mission/incident/iter154-unfenced-pi","provenance":"w-prove-1-0-phase-a"}
{"hash":"sha256:cacf22c211a74d8c28c94c4dafecf2d2f78e4ee41854b37290889fe7944b40dd","interfaceHash":"sha256:cca7a1ba281879dbf7cfc070cdbc6fa69d69339c5c05d15bebe73136434b0546","semanticId":"world/mission/incident/iter155-faithfulness-proof","provenance":"w-prove-1-0-phase-a"}
{"hash":"sha256:722751941829c9165764b76d78d435f83ea7ef3aeadd573c7b29d4fdaf629d2b","interfaceHash":"sha256:8dc8c03857020a11a4c39f31616c82de40a553548b1cfc4a6b4bb54d2537294f","semanticId":"world/mission/incident/iter181-ci-bench-401","provenance":"w-prove-1-0-phase-a"}
[rc=0]
$ .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:722751941829c9165764b76d78d435f83ea7ef3aeadd573c7b29d4fdaf629d2b --payload | python3 -c 'import json,sys,base64; print(base64.b64decode(json.load(sys.stdin)["payload"]).decode())'
{"question":"Why did PR #141's remote CI go red on `BenchmarkRESTCommit` when every local gate was green?","answer":"a036062 is an ancestor of base and its diff adds `req.Header.Set(\"Authorization\", auth)` to bench_test.go — the bench predated the session gate and POSTed to the now-protected /v1/commit unauthenticated, 401 (the middleware working as designed).","diagnosed_at":"2026-09-24","sources":["sha256:43bc14de993204f80bcb62beaba93fd6f5bf7940108ad82c29ebfb5fa0f35b03","sha256:483111f4c441085a9cd9ede52f21218884d7cde6a01be7d2cbf55e1d203c43a3"]}
[rc=0]
$ for h in $(.walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get sha256:722751941829c9165764b76d78d435f83ea7ef3aeadd573c7b29d4fdaf629d2b --payload | python3 -c 'import json,sys,base64; [print(s) for s in json.loads(base64.b64decode(json.load(sys.stdin)["payload"]))["sources"]]'); do echo "-- $h"; .walk-scratch/ailang-worldd --addr http://127.0.0.1:7654 object get $h --payload | python3 -c 'import json,sys,base64; o=json.load(sys.stdin); print(o["semanticId"]); print(base64.b64decode(o["payload"]).decode())'; done
-- sha256:43bc14de993204f80bcb62beaba93fd6f5bf7940108ad82c29ebfb5fa0f35b03
world/mission/evidence/iter181-ci-bench-401/0
{"kind":"commit","ref":"git a036062","check":"git show a036062 -- host/daemon/bench_test.go | grep -c 'req.Header.Set(\"Authorization\", auth)' >= 1","excerpt":"`a036062` is an ancestor of base and its diff adds `req.Header.Set(\"Authorization\", auth)` to `bench_test.go` — the bench predated the session gate and 401'd."}
-- sha256:483111f4c441085a9cd9ede52f21218884d7cde6a01be7d2cbf55e1d203c43a3
world/mission/evidence/iter181-ci-bench-401/1
{"kind":"status","ref":"world-mission.md STATUS 2026-09-24 (iter-181)","check":"grep \"bench predates the sprint\" world-mission.md","excerpt":"remote CI on PR #141 then caught what every local gate missed — `BenchmarkRESTCommit` POSTing to the now-protected `/v1/commit` unauthenticated (401: the middleware working as designed; the bench predates the sprint and was not in the plan's local gate list)."}
[rc=0]
# assembled answer (incident.answer): the bench predated the session gate and POSTed to the now-protected /v1/commit unauthenticated -> 401; a036062 adds the Authorization header to bench_test.go
$ git merge-base --is-ancestor a036062 HEAD && echo ANCESTOR; git show a036062 -- host/daemon/bench_test.go | grep -c 'req.Header.Set("Authorization", auth)'
ANCESTOR
1
[rc=0]
$ git show a036062 -- host/daemon/bench_test.go | grep -n '^[-+].*Authorization'
91:+		req.Header.Set("Authorization", auth)
[rc=0]
$ grep -n "bench predates the sprint" design_docs/world-mission.md | cut -c1-220
928:## STATUS 2026-09-24 (iteration 181) — **ROW 39 LANDED: THE PARKED, TWICE-JUDGED SPRINT IS MERGED AS `a036062` WITH REMOTE CI GREEN — THE FIRST PI-CONTROLLER LAND IN THIS REPO, ON THE ATTENDED `D-WORLD-36` HARNES
[rc=0]
$ echo "duration_s=$(( $(date +%s) - $(cat .walk-scratch/q4.t0) ))"; date -u +%Y-%m-%dT%H:%M:%SZ   # STOP clock
duration_s=10
2026-09-24T09:36:54Z
[rc=0]
```

## Daemon stop

```
$ date -u +%Y-%m-%dT%H:%M:%SZ; kill -TERM $(cat .walk-scratch/serve.pid); echo "SIGTERM -> pid $(cat .walk-scratch/serve.pid)"
2026-09-24T09:37:03Z
SIGTERM -> pid 98897
[rc=0]
$ deadline=$(( $(date +%s) + 20 )); until grep -q "serve exited" .walk-scratch/serve.log; do [ $(date +%s) -ge $deadline ] && { echo STOP-TIMEOUT; break; }; sleep 0.2; done; cat .walk-scratch/serve.log
ailang-worldd listening on http://127.0.0.1:7654
serve exited rc=0
[rc=0]
$ lsof -nP -iTCP:7654 -sTCP:LISTEN; echo "listeners-on-7654 rc=$? (1 = none)"
listeners-on-7654 rc=1 (1 = none)
[rc=0]
```
