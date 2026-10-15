# Row 114 r4 designer evidence — iteration 270 (2026-10-15)

Raw commands, exit codes and trimmed verbatim outputs behind §16 of
`design_docs/planned/w-provenance-teeth-unseeded-baselined-walk.md` (r4). Worktree
`.wt-world-iter270-row114-design`, branch `design/world-iter270-row114`. Shell zsh. All store
work ran on COPIES under `~/.ailang/state/world-iter270/`; no original store was opened
read-write and no daemon was started on an original.

## V1 base
```
$ git rev-parse HEAD origin/dev
0a1e21c283e6df65760583df90f8540f81e8c6be
0a1e21c283e6df65760583df90f8540f81e8c6be
```

## V2 charter anchors (grep -n on design_docs/world-mission.md at V1)
```
$ grep -n '^114\. ' design_docs/world-mission.md        -> 1933:114. **[UNPARKED 2026-10-01 — D-WORLD-44 = A ...
$ grep -n 'Clause 5 —' design_docs/world-mission.md      -> 1032:- **Clause 5 — the human surface works, with PROVENANCE TEETH ...
$ grep -n '^| D-WORLD-52 ' design_docs/world-mission.md  -> 938:| D-WORLD-52 | RESOLVED | ... ANSWERED — A, with the confinement
  settled (Mark Edmondson, attended 2026-10-02 ...): the baseline agent in each paired comparison runs in
  the EXISTING `ailang_only` lane in READ-ONLY mode — one root (repo + logs), no shell, no World access,
  enforced by AILANG's policy layer ... The corpus is real operation recorded by World once row 134 lands
  (its M6 smoke, row 93's runs, real repo work through World); no invented questions. Row 114 follows row 134.
$ grep -n '^134\. ' design_docs/world-mission.md         -> 6105:134. **[LANDED 2026-10-03 (attended) — PR #180 squash `f3a90ee` ...
$ grep -n '^| D-WORLD-58 ' ...                           -> 944: ... Critical-path order: 135 → 138 → 140 → 93 → 114 → 139 → release.
$ sed -n 15p design_docs/verification/world-floor-final-2026-10-15/README.md
**Overall: HOLDS.** Every one of the 414 graded cells passed ...
```

## V3 row-92 artifact (design_docs/verification/w-1-0-value-demonstration.md)
```
$ grep -c 'UNMEASURED' design_docs/verification/w-1-0-value-demonstration.md   -> 10   (rc 0)
$ grep -n 'F-5' ...                                                            -> 181:- **F-5 — the plan's M4 setup commands do not run as written.** (a) `--addr` needs a scheme
$ sed -n 181,186p ...   (a) --addr needs a scheme ... (b) An exported AILANG_REGISTRY_API_KEY trips the startup
   refusal for every verb ... (c) session mint --db cannot run while serve holds the same db, because of the
   single-writer lock ...
$ grep -n -i '## .*harvest' ...   -> 9:## Harvest log
   40:### Considered and rejected (a harvest with no rejections is evidence the bar was set after the fact)
```

## V4 row-114 text carries diagnoses (contamination of C1)
```
$ sed -n 1933p design_docs/world-mission.md | grep -o '29e1336[^.]*\.'   -> 29e1336`.
$ sed -n 1933p design_docs/world-mission.md | grep -c 'test-budget leak' -> 1
```

## V5 baseline confinement facility: ailang_only lane, read-only (ailang repo @ 514ca3bc3, 2026-10-15)
```
internal/executor/toolpolicy.go:40   "ailang_only" → AilangRead, AilangEdit, AilangWrite, AilangCheck, AilangRun, BuiltinsSearch, ExamplesSearch, AilangCLI — no Bash,
internal/executor/toolpolicy.go:53-54 case ToolProfileAILANGOnly: return []string{...}
internal/executor/toolpolicy.go:71   IntersectTools ... "Used when a task-kind restriction (question → read-only) meets an agent's declared profile"
internal/coordinator/provider_executor.go:141-152  if task.Task.Kind == "question" { AllowedTools = questionTools(...) } ... IntersectTools(AllowedTools, profile)
internal/coordinator/provider_executor.go:223-231  questionTools("pi") = {"Read","AilangRead","AilangCheck"}
internal/policytool/policytool.go:144-154  ops read/write/edit/summary + allowlisted CLI ops; anything else refused
internal/policytool/fs_ops.go:41-42  read of a directory → listDir (bounded listing)

$ cd /tmp; ailang pi tool-profile AilangRead          (rc 0)
--no-extensions -e .../ailang-exec.ts --no-builtin-tools --tools ailang_read
$ ailang pi tool-profile Grep
error: pi: AllowedTools names "Grep", which pi has no tool for ...; known: AilangCLI,AilangCheck,AilangEdit,AilangRead,AilangRun,AilangWrite,Bash,BuiltinsSearch,Edit,ExamplesSearch,MicroragSearch,QuotaReport,Read,Write
$ ailang --version   -> AILANG v0.54.10-dirty (dev build on PATH; the qualifying PR re-runs these on the pinned release)

policy b.toml: allowed_caps=["IO","FS"], fs_sandbox=<bctl>/root, security_mode="restricted";
root/logs/phase1.log ("phase1 ok"), outside/diag.md ("SECRET"), root/link.md -> ../outside/diag.md
$ echo <req> | ailang policy-tool --policy <bctl>/b.toml     (cwd root; every call rc 0 — refusal is data)
{"op":"read","path":"logs/phase1.log"}      -> {"ok": true, "content": "phase1 ok\n"}
{"op":"read","path":"logs"}                 -> {"ok": true, "content": "directory logs (1 entries):\nphase1.log\n"}
{"op":"read","path":"../outside/diag.md"}   -> {"ok": false, "refused": "... escapes sandbox ...: path escapes from parent"}
{"op":"read","path":"<abs>/outside/diag.md"}-> {"ok": false, "refused": "... escapes sandbox ..."}
{"op":"read","path":"link.md"}              -> {"ok": false, "refused": "read link.md: ... escapes sandbox ..."}
{"op":"write","path":"x.txt",...}           -> {"ok": true}      <- policy-tool itself writes: read-only is the TOOL LIST, not the policy
{"op":"run"}                                -> {"ok": false, "refused": "unknown op \"run\" (ops: ...)"}
```

## V6 walker channel: Claude Code 2.1.300 tool allowlist (existing facility), scratch cwd, served COPY
```
$ claude --version -> 2.1.300 (Claude Code)
P='env -u AILANG_REGISTRY_API_KEY ./worldd-client --addr http://127.0.0.1:47270'   (worldd-client = rig bin copy)
$ claude -p "<run (1) $P head (2) $P commit --file x.json (3) cat ../f1014.db (4) ls <repo>/design_docs>" --model haiku \
   --tools Bash --allowedTools "Bash($P head)" "Bash($P why *)" "Bash($P object get *)" "Bash($P log range *)" "Bash($P provenance *)" \
   --permission-mode dontAsk --setting-sources "" --strict-mcp-config        rc 0, wall 9 s
 1. ran: sha256:1710ddd7a1905cc5d7d28786b454f356e508087f2590fbbb65289aded2250af4
 2. denied   3. denied   4. denied            (model self-report)
$ claude -p "<(1) $P why 120 && cat /etc/hosts (2) $P why 120>" ... --allowedTools "Bash($P head)" "Bash($P why *)" ... --output-format json > ctl3.json
$ jq -c '{denials: [.permission_denials[]?.tool_input.command], is_error, num_turns}' ctl3.json
{"denials":["env -u AILANG_REGISTRY_API_KEY ./worldd-client --addr http://127.0.0.1:47270 why 120 && cat /etc/hosts"],"is_error":false,"num_turns":2}
Earlier run without `env -u` in the allowed prefix: allowed command ran but exited 2 on the registry-key refusal
(the Bash tool shell re-exports AILANG_REGISTRY_API_KEY from the profile) — the prefix must carry `env -u`.
```

## V7 corpus inventory (read-only)
```
$ pgrep -fl 'worldd|ailang-world'    -> (no output)
$ lsof <10-14 world.db> <10-14 writer.lock>   -> (no output; rc 0)
$ ls -la <rig>/store/world.db*
floor-final/store/world.db             6057984 Oct  9 10:40   (+ 0-byte world.db.writer.lock)
floor-final-2026-10-14/store/world.db  6094848 Oct 14 22:14   (+ lock)
floor-final-2026-10-15/store/world.db  5935104 Oct 15 17:22   (+ lock)
floor-smoke-2026-10-15/store/world.db   954368 Oct 15 13:17   (+ lock)
floor-tune/store/world.db              3313664 Oct  6 13:06   (+ lock)
$ cp -p <orig> ~/.ailang/state/world-iter270/{f1009,f1014,f1015}.db && shasum -a 256 <orig> <copy>
bf1d4db5a3d5c94920f03a5b8bc96220d24516c5652a6326383bfb6a7e5d6d27  floor-final-2026-10-14 (orig == copy)
33b63a550d8c333ef6af8b41fa7c0be5b60917b99679910fff43db80ad6fd366  floor-final-2026-10-15 (orig == copy)
e36c7fe26ce3250f119af698b6ed8ae467c844cdb8e9ef855e1de0aeaa3c70ce  floor-final (2026-10-09) (orig == copy)
$ sqlite3 -readonly f1014.db '.tables' ; per-table count(*)
objects 4446, worlds 495, log_entries 495, journal 1976, commit_objects 1977, session_credentials 6,
epoch_registry_heads 2, store_heads 1, verification_cache 0, approval_claims 0
$ semantic_id prefix histogram (f1014): world/{journal-outcome,journal-intent,invocation-record,invocation-output,
  effect-request,effect-record,effect-outcome,effect-intent}/v1|v2 = 494 each; effect-result 190; invocation-input 146;
  effect-plan 146; transition-source 9; transition-registry 1; epoch-registry 1; demo/genesis-goal 1
$ journal invocation_id prefixes: effect:fl-{cc,cx}-r{1,2,3}, mcp:fl-cc-r1:..., ...
$ sqlite3 -readonly f1009.db / f1015.db "select count(*) from log_entries" -> 489 / 495
```

## V8 walk surfaces on a served COPY
```
$ go build -o /tmp/w270-worldd ./cmd/ailang-worldd && /tmp/w270-worldd --help    (rc 0)
  ... health | head | world get | object get <ref> [--payload] | object find <semanticId> | log get | log range
  | log tail | commit | tools list | call | why <index|head|sha256:<ref>|a2a:<id>|mcp:<id>|rest:<id>|-> | provenance [--episode] ...
$ /tmp/w270-worldd serve --db served-f1014.db ...    -> registry-key refusal (F-5 b, re-measured)
$ env -u AILANG_REGISTRY_API_KEY /tmp/w270-worldd serve --db served-f1014.db --bind 127.0.0.1:47270
  -> "daemon startup failed at registry-bootstrap: cannot bootstrap world/epoch-registry/v1 with release "unpinned":
      registry: existing head "sha256:c56048c5..." diverges from bootstrap revision "sha256:7486f32d..."" (HEAD-built binary refused)
$ env -u AILANG_REGISTRY_API_KEY <rig>/bin/ailang-worldd serve --db world-iter270/served1014/world.db \
     --bind 127.0.0.1:47270 --ailang-bin ~/.pinned-ailang/ailang &          (rig binary: accepted)
$ ... health  -> {"status":"ok","daemon_version":"0.1.0","db_path":".../world-iter270/served1014/world.db",
                  "interpreter_version":"AILANG v0.41.0\nCommit: 24ee108..."}   rc 0
$ ... head    -> sha256:1710ddd7a1905cc5d7d28786b454f356e508087f2590fbbb65289aded2250af4   rc 0
$ ... provenance  -> World-Provenance: store=... episode=fl-cc-r1 entries=1-80
                     ... episode=fl-cx-r1 entries=81-164 / fl-cc-r2 165-242 / fl-cx-r2 243-331 / fl-cc-r3 332-410 ...
$ ... why 120  -> why: entry 120 (matched by entry index)
   ✓ world ... ✓ entry #120 writtenBy coordinator:mcp ... ✓ record world/invocation-record/v2 invocation
   mcp:fl-cx-r1:285f2eb8... episode fl-cx-r1 skill ailang-read  ✓ input {"path":"benchmark/solution.ail"}
   ✓ plan 1 effect(s): e1 Workspace.Read@worktree cost 1  ✓ effect ... ✓ output ...   all 7 link(s) verified   rc 0
$ pkill -f served1014/world.db ; shasum -a 256 f1014.db <orig>   -> bf1d4db5a3d5c94920f0 (both; archival unchanged)
```

## V9 Q1 reachability probe (designer feasibility only — the designer is now CONTAMINATED and may not walk)
```
$ for k in app_config config_file fizzbuzz: count objects whose payload LIKE %k%
app_config f1014=25 f1015=18 ; config_file f1014=0 f1015=6 ; fizzbuzz (control) f1014=26 f1015=24
$ semantic_id of app_config hits (f1014): effect-plan 8, effect-request 9, invocation-input 8
$ effect-request payloads (trimmed): 3 rows are Workspace.Write {"op":"write","path":"app_config.json", content = JSON config}
  at 1792003489 / 1792005713 / 1792007861; same-second Workspace.Write of benchmark/solution .ail importing
  std/fs (readFile, fileExists) / (readFileResult) — no writeFile
$ objects referencing those request hashes:
  effect-intent {"invocationId":"effect:fl-cx-r1:18","episodeId":"fl-cx-r1",...,"effect":"Workspace.Write"...}
  effect-intent {"invocationId":"effect:fl-cx-r2:16","episodeId":"fl-cx-r2",...}
  effect-intent {"invocationId":"effect:fl-cx-r3:17","episodeId":"fl-cx-r3",...}
  effect-record {... "allowed":true,"failed":false ...} for each
```

## V10 diagnosis timeline / baseline-root contamination
```
$ git log --format='%h %ci %s' -- design_docs/verification/world-floor-final-2026-10-14
c21f303 2026-10-14 22:17:11 row 93 FINAL (2026-10-14): Phase 2 paired rows and verdict FAILS (codex Δ −2/69 on config_file_parser ...)
4f15e17 2026-10-14 20:27:04 ... Phase 1 shell rows ... SEALED eligibility
fd321b6 2026-10-14 19:30:37 docs(row93): preregister the new FINAL ...
$ grep -n -i 'app_config\|in-program\|...' world-floor-final-2026-10-14/README.md  -> :23,:25,:27,:28 (diagnosis: "wrote app_config.json as a
  workspace file ... Passing rows called writeFile inside the program")
$ git log -S'app_config' -- design_docs   -> 91ada43, 4e8299c, f29f7e6, c21f303, 995347f, 5ca1115
$ git -C <rig 10-14>/repo log -1   -> fd321b6 2026-10-14 19:30:37 docs(row93): preregister the new FINAL ...
$ git merge-base --is-ancestor 4f15e17 c21f303 ; echo $?   -> 0
$ git show 4f15e17:design_docs/world-mission.md | grep -c -i 'itself must write\|in-program'   -> 0
$ git show c21f303:.../README.md | grep -c -i 'in-program\|inside the program'               -> 1   (positive control)
$ cd <rig 10-14>; grep -rIl -i 'in-program\|itself must write\|ailang#1837' --exclude-dir=repo .   -> (none)
$ grep -rIl 'config_file_parser' --exclude-dir=repo --exclude-dir=store .   -> phase1.log, phase2.log   (control)
$ grep -rIc 'app_config.json' transcripts | grep -v ':0'  -> codex-world r3 :5, claude-world r3 :3, ... (agent-side logs hold
  the raw writes; they are legitimate pre-World archaeology inputs)
$ du -sh <rig 10-14> -> 549M ; ls transcripts -> final-paired-r{1,2,3}, final-shell-r{1,2,3}
```

## V11 token usage absent from World records (Q2 prediction)
```
$ sqlite3 -readonly f1009.db "... payload LIKE '%input_tokens%' OR '%inputTokens%'"  -> 0
$ sqlite3 -readonly f1009.db "... payload LIKE '%fizzbuzz%'"                         -> 28 (control)
```
