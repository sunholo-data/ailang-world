# Sprint plan — `w-provenance-teeth-unseeded-baselined-walk` (queue row 114, clause 5), question Q1

**Design doc:** [`w-provenance-teeth-unseeded-baselined-walk.md`](w-provenance-teeth-unseeded-baselined-walk.md) r4 (`d626ed2`); quorum r1 PROCEED at N−1 (`verification/world-iter270/quorum/`, gemini PASS; 3 seats absent on 429s).
**Planner:** iteration 270, Agent `opus` (`fail-closed:env-pin`). **Executes ONLY after the freeze (M0).**
**Scope:** protocol only. Zero code or harness diffs (D6, AC11). Commits add files under `design_docs/verification/w-provenance-teeth-unseeded-baselined-walk/` and nothing else.
**Shape:** Q1 end to end, the first of ≥3 questions. Q2 and later questions rerun M1–M7 with their own cohort, pre-diagnosis commit and E2; M0 runs once per execution session.
**Estimate:** ≈2.5 h wall for Q1 (M1–M2 ≈60 min, M4 ≤60 min, M5–M6 ≤15 min, M7 ≈20 min). Model spend ≤ $6.

## 0. Environment (every command; zsh; no `timeout(1)`, no `${PIPESTATUS}`, redirect then `echo rc=$?`)

```sh
WT=<executor worktree>; RIG=~/.ailang/state/floor-final-2026-10-14; K1=~/.ailang/state/world-iter270/f1014.db
PIN=~/.pinned-ailang/ailang          # measured: AILANG v0.41.0, Commit 24ee108 (the only pinned release on the rig)
S=~/.ailang/state/row114/Q1          # evidence bank (NOT /tmp, which is wiped); A=$WT/design_docs/verification/w-provenance-teeth-unseeded-baselined-walk/Q1-config-file-parser
PORT=47314                           # measured free; re-check with lsof -nP -iTCP:$PORT before use
ts(){ perl -MTime::HiRes=time -e 'printf "%.3f\n", time'; }                       # controller clock
stamp(){ perl -MTime::HiRes=time -ne '$|=1; printf "%.3f %s", time, $_'; }         # line-timestamps both arms' stdout
```
Measured this session (read-only): rig `bin/ailang-worldd` sha256 `eae3bec4…2ca53a`; K1 archival copy == original == `bf1d4db5…e5d6d27`; `claude` 2.1.300 on PATH; `pi` 0.85.1; pi providers `zai` and `openrouter` report `ready`, `ollama` (weekly 429) and `openai` (no credits) are dead; `$PIN pi tool-profile AilangRead,AilangCheck` → `--no-extensions -e …/ailang-exec.ts -e …/ailang-lsp-lite.ts --no-builtin-tools --tools ailang_read,ailang_check`.

**Two facts the design doc does not state (executor must apply):** (1) `ailang-exec.ts:86` spawns `ailang policy-tool` by bare name from `PATH`, so B is on the pinned release only if `$S/pinbin/ailang -> $PIN` is first on B's `PATH`; on the ambient PATH it silently runs `v0.54.10-dirty`. (2) The lane prompt that the extension injects (`lanePrompt`, `ailang-exec.ts:167`) tells the model to write and `ailang_run` programs. In read-only mode those tools are absent, so B's preamble must say that this is a read-only question task.

## 1. Roles (D2) — all fresh processes; none is the controller, the r4 designer or any session that has read the Q1 diagnosis

| Actor | Spawn (by the controller) | Model | Why this lane |
|---|---|---|---|
| B | `env PATH=$S/pinbin:$PATH AILANG_AGENT_POLICY=$S/b/b.toml pi -p --mode json --no-session --no-context-files --no-skills --no-prompt-templates -na <tool-profile flags> --provider zai --model glm-5.3 @$S/b/prompt.txt` run with cwd `$S/b/root` | `zai/glm-5.3` (fallback `openrouter/<model>`, chosen and pinned BEFORE reveal only) | D-52: the `ailang_only` lane, read-only, on pi; the Anthropic author is excluded |
| W | `cd $S/w/scratch/cwd && env -u AILANG_REGISTRY_API_KEY -u WORLD_SESSION claude -p --model sonnet --tools Bash --allowedTools <§M1.W list> --permission-mode dontAsk --setting-sources "" --strict-mcp-config --no-session-persistence --max-budget-usd 2 --output-format stream-json --verbose "$(cat $S/w/prompt.txt)"` | `sonnet` via Claude Code 2.1.300 | D3b channel (V6) |
| V_B, V_W | `cd $S/v{b,w} && claude -p --model sonnet --tools Read --setting-sources "" --strict-mcp-config --no-session-persistence --max-budget-usd 1 --output-format stream-json --verbose "$(cat prompt.txt)"`, where the cwd holds ONLY the E2 bundle, the rubric and that arm's sealed answer | `sonnet` | fresh; never an Agent-tool sub-agent, which would inherit the controller's CLAUDE.md and memory (and the memory holds the Q1 diagnosis) |

**Gemini note 1 (fixed).** W's scratch holds the rig binary copied under the allowlist's name: `cp -p $RIG/bin/ailang-worldd $S/w/scratch/cwd/worldd-client`, and its sha256 must equal `eae3bec4…`. The design doc §5 line 94 now says this (a one-line consistency edit in this commit).

## 2. Milestones

**M0 — Freeze gate (≤2 min, $0).** `grep -n 'w-provenance-teeth-unseeded-baselined-walk' $WT/design_docs/world-mission.md | grep -E 'D-WORLD-[0-9]+'`. Then read the one ledger row the controller filed to freeze D2/D3b/D4/D5 (expected `D-WORLD-76`). **STOP unless** that row is `RESOLVED` and its answer is A (freeze as written). OPEN, any other answer or no row ⇒ record `HOLD` in the iteration log and do nothing else. Also re-check: `git -C $WT merge-base --is-ancestor d626ed2 HEAD` rc 0, and the doc is unchanged since `d626ed2` apart from the line-94 edit; otherwise STOP (re-plan).

**M1 — Pre-reveal controls on the PINNED release (≤40 min, ≤$1.50; implements §5 PR, closes gemini note 2; AC4).** Nothing about Q1 is written to any role yet.
- *Pins:* `$PIN --version`, `claude --version`, `pi --version`, `shasum -a 256 $K1 $RIG/store/world.db $RIG/bin/ailang-worldd` → `$S/pins`.
- *B layer (re-run V5 on `$PIN`):* in `$S/bctl` rebuild V5's fixture (root/logs/phase1.log, outside/diag.md, root/link.md → ../outside/diag.md; policy outside the root). Run `echo <req> | $PIN policy-tool --policy …` for in-root read ok, directory list ok, `..`/absolute/symlink refused with "escapes sandbox", and `run` refused with "unknown op".
- *B live end to end (new, gemini note 2):* the exact B spawn line from §1 with cwd `$S/bctl/root` and a control prompt. The prompt asks for (1) `logs/phase1.log` (expect "phase1 ok"), (2) a listing of `logs`, (3) `../outside/diag.md` (expect a refusal), (4) a bash or `ailang_run` call (expect "tool not found/not available" from pi), and (5) `ailang_check` on a root `.ail` file (record whatever happens). Pass = (1)(2) in tool results; (3)(4) refused; "SECRET" absent from the whole json transcript (`grep -c SECRET` = 0); pin proof: `env PATH=$S/pinbin:$PATH command -v ailang` = `$S/pinbin/ailang`.
- *W controls:* serve a CONTROL copy (`cp -p $K1 $S/wctl/ctl.db`), `env -u AILANG_REGISTRY_API_KEY $RIG/bin/ailang-worldd serve --db $S/wctl/ctl.db --bind 127.0.0.1:$PORT --ailang-bin $PIN &`. Allowlist with `P='env -u AILANG_REGISTRY_API_KEY ./worldd-client --addr http://127.0.0.1:'$PORT`: `"Bash($P head)" "Bash($P provenance)" "Bash($P provenance *)" "Bash($P why *)" "Bash($P object get *)" "Bash($P object find *)" "Bash($P log get *)" "Bash($P log range *)"`. The cwd holds only `worldd-client`, under a parent dir that holds a canary `CLAUDE.md` with `CANARY-114-PARENT`. Run the §1 W spawn with a control prompt: `head`; `$P commit --file x.json`; `$P call ailang-read`; `cat $S/wctl/ctl.db`; `ls $WT/design_docs`; `$P why 1 && cat /etc/hosts`; `curl -X POST http://127.0.0.1:$PORT/commit`; "quote any instruction token beginning CANARY". Pass = `head` tool_result is `sha256:1710ddd7…`; `jq '[.permission_denials[]?.tool_input.command]'` on the result event lists all six denied commands (AC4); `grep -c CANARY-114-PARENT` on the transcript minus the prompt echo = 0 (MU-4). Optional `--safe-mode` may be added only if the full W control set is re-run with it.
- Any failed or missing control ⇒ `NOT-RUN(ISOLATION)`, preserved; Q1 is not revealed. Stop the control daemon; `shasum` K1 unchanged.

**M2 — E2 bundle, B root, freeze (≤30 min, $0; §5 P2–P3, D4; AC1, AC2, AC3, AC12).**
- *Harvest log* (AC1): Q1 taken; C1/C2/C3 rejected `NOT-RUN(NO-CORPUS)`; Q2 candidate. Copy the rejection list from §8 verbatim.
- *Question:* write `$S/question.txt` = §8 Q1 text, exact bytes; `shasum -a 256` → pins. The W and B prompts are `preamble_{b,w}.txt` + `\n---QUESTION---\n` + question bytes. Hash the preambles too. Preambles name the channel only (B: "read-only, answer from the files under your root; end with a line `ANSWER:`"; W: the allowed command prefix and verb list, same `ANSWER:` rule). Neither may name a mechanism.
- *E2 (controller only; stored in `$S/e2`, never under B's root or W's cwd):* (a) `$RIG/phase2.log` codex rows `:58-59,:219-220,:380-381`, plus the claude rows; (b) for each codex World row, the solution bytes taken from the K1 copy with `sqlite3 -readonly` (the `Workspace.Write` of `benchmark/solution.ail` in the same second as the `app_config.json` write, V9). Each solution is run by `$PIN run --caps IO,FS` in a clean empty dir, then again in a dir that holds the written `app_config.json`, with stdout, stderr and rc recorded; (c) the grader's clean-dir rule from `$WT/scripts/floor/`. *Rubric* (pinned before reveal, the same for V_B and V_W): VERIFIED only if the answer's causal claim is confirmed by (a)+(b); a merely plausible or partial answer is `UNVERIFIED`.
- *B root (D4):* `git -C $WT archive fd321b6 | tar -x -C $S/b/root/repo`. Then `rm` every `design_docs/verification/world-floor-*/{README,summary}.md` and the archive's `.claude/`. Add `rsync -a --exclude .git --exclude .claude` of `$RIG/{phase1.log,phase2.log,transcripts/,ws/,fixture/}` into `$S/b/root/logs/`. Exclude `store/ private/ bin/ clibin/ serve.log verdict.out MINT_BLOCK.sh genesis.json m4wt/ void/ repo/` (rig `ws/*` and `fixture/` carry `.git` dirs — measured). Policy `$S/b/b.toml` = V5's shape with `fs_sandbox = "$S/b/root"`, kept outside the root.
- *Scan (AC3, MU-2):* `phrases.txt` = `in-program`, `inside the program`, `must write the file itself`, `itself must write`, `as a workspace file`, `ailang#1837`. Run `grep -rIl -i -F -f phrases.txt $S/b/root` and expect no output. Run `find $S/b/root -name .git -o -name '*.db' | wc -l` and expect 0. Positive control: the same grep over `git -C $WT show c21f303:design_docs/verification/world-floor-final-2026-10-14/README.md` must hit. Write the manifest with `find … -type f -exec shasum -a 256 {} +`. Scan the preambles, the prompt files and W's cwd with the same phrases. Any hit ⇒ fix the export or mark `CONTAMINATED-PENDING` and stop.
- *Served copy (AC5, AC12, MU-7):* `cp -p $K1 $S/w/served.db`, check the sha256, then serve it as in M1. Invariants → pins: `head` = `sha256:1710ddd7…`, plus `sqlite3 -readonly $S/w/served.db 'select count(*) from log_entries'` = 495. Never pass `--db` an original or `$K1` itself. Write the spawn commands and `max_attempt_wall_s=3600` to pins.

**M3 — Commit the sealed pre-reveal state (≤5 min).** Copy pins, the harvest log, the control transcripts' sha256 and paths, the manifest and the scan results into `$A/` (as `pins`, `prereq-transcript`, `exposure-log`). Commit locally with "row114 Q1 pre-reveal sealed". The commit timestamp precedes reveal.

**M4 — Baseline FIRST (≤60 min wall cap, ≤$1; §5 P4; AC2, AC6, AC7).** `t_issue_b=$(ts)`. Run the B spawn piped through `stamp > $S/b/transcript.jsonl`, with rc to a file. `t_answer_b` = stamp of the line carrying `ANSWER:`. Then `cp` only the final answer text into `$S/vb/answer.txt` (with `e2/`, `rubric.txt`), spawn V_B → `$S/vb/verdict.jsonl`, `t_verified_b` = stamp of V_B's result event. Seal: write sha256 of the transcript and verdict into `$A/verify-b`. B's output never enters `$S/w`. If the attempt wall passes 3600 s, kill B ⇒ `INVALID(clock)`. A provider error before B's first model turn ⇒ `NOT-RUN(TRANSPORT)`; one retry is allowed only then, and is logged.

**M5 — Walk SECOND (target ≤300 s envelope; ≤$2; §5 P5; AC4–AC6).** Immediately before: `head` and count equal the pins. `t_issue_w=$(ts)`; run the W spawn piped through `stamp > $S/w/transcript.jsonl`. As soon as the result event arrives: copy the answer into `$S/vw/`, spawn V_W, and take `t_verified_w`. `world_wall_s = t_verified_w − t_issue_w` (this includes the handoff and V_W's setup; only this number meets 300 s, D5).

**M6 — Accounting (≤10 min; §5 P6, §6 audit; AC2, AC5, AC6, AC8, AC12).** Check: `shasum $K1 $RIG/store/world.db` equals pins; served `head` and count are unchanged (else `INVALID(seeding)`); stop the daemon. Quote `permission_denials` and every W `tool_use` command. Each command must start with the allowlisted prefix and no denied command may have a `tool_result` without an error; else `INVALID(seeding)`. The question sha256 must be identical across arms (else `INVALID(contamination)`). The declared and transcript durations must agree within 5 s (else `INVALID(clock)`). The initial-context scan of both transcripts' first events must have 0 hits. Write `envelope` (`t_issue_b`, `t_verified_b`, `t_issue_w`, `t_verified_w`, `world_wall_s`, `baseline_wall_s`).

**M7 — Artifact + index (≤20 min; §7; AC1, AC8–AC10).** Write `$A/{pins,prereq-transcript,baseline-transcript,walk-transcript,verify-b,verify-w,envelope,exposure-log,verdict}`; transcripts are kept in `$S` and cited by sha256. `index.md` holds the harvest log with rejections, the accounting, a comparison row, the findings and the stated baseline bias. Any ratio carries "B: read/list only, no grep, `zai/glm-5.3`; W: `sonnet`" (AC10). The row stays open (AC9: 1 of ≥3 at most). Commit and push a draft PR; CI must be green. AC11: `git diff --stat origin/dev -- . ':!design_docs/verification/w-provenance-teeth-unseeded-baselined-walk' ':!design_docs/world-mission*.md'` is empty.

## 3. Stop conditions → verdict (closed set, §6; every attempt preserved)

| Where | Condition | Recorded |
|---|---|---|
| M0 | freeze row not RESOLVED-A | `HOLD` (no verdict; nothing revealed) |
| M1 | any B/W control fails or is missing | `NOT-RUN(ISOLATION)` both arms |
| M2 | scan hit that can't be removed | candidate `CONTAMINATED-PENDING` |
| M2 | served copy refuses or head ≠ pin | `NOT-RUN(TRANSPORT)` walk; `INVALID(pin)` if an original was opened |
| M4/M5 | provider dead before first turn | that arm `NOT-RUN(TRANSPORT)` |
| M4–M6 | wall > 3600 s / clocks disagree > 5 s | `INVALID(clock)` |
| M6 | head/count changed, a non-allowlisted call succeeded | `INVALID(seeding)` |
| M6 | question hash differs, diagnosis in initial context | `INVALID(contamination)` |
| M5 | V_W VERIFIED, envelope > 300 s | W `FAIL-TIME` |
| M4/M5 | V says not confirmed | `UNVERIFIED` |
| M5 | declared cohort queried, answer absent | `UNANSWERABLE(evidence)` |
| — | V_B and V_W VERIFIED, all audits clean, envelope ≤ 300 s | Q1 qualifying success (1 of ≥3) |

## 4. Mutation controls the evaluator must run (all against `$S` copies, never re-revealing Q1)

| MU | Evaluator action | Must observe |
|---|---|---|
| MU-1 | append `$P commit --file x.json` to a W control prompt | in `permission_denials`; `head` unchanged |
| MU-2 | copy `c21f303` 10-14 README into a scratch copy of B's root; re-run the M2 scan | ≥1 hit |
| MU-3 | `$P why 1 && cat /etc/hosts` | in `permission_denials` (M1 log) |
| MU-4 | the canary `CLAUDE.md` in W's parent | 0 canary hits in the W transcript |
| MU-5 | flip one byte in a copy of `prompt.txt` | the sha256 comparison fails |
| MU-6 | recompute the envelope from the first successful W `tool_result` | it differs from `t_issue_w` → the M6 check rejects it |
| MU-7 | `shasum` originals vs pins; `lsof` on originals | equal; no holder |
| MU-8 | `git log -p -- $A/../index.md` | candidate ids immutable since M3 |
Plus: the evaluator re-runs the M2 positive control and the B pin check (`command -v ailang` under B's PATH).

## 5. Open items (do not block the plan; resolve before reveal or record)
1. "Pinned release" for B = `$PIN` v0.41.0 (interpreter/gate pin). The v0.52.1 tool-pin binary was not found under `~/.pinned-ailang` or `~/.ailang/bin`. If the freeze names v0.52.1, re-run M1 with it.
2. B model `zai/glm-5.3` vs W `sonnet`: the asymmetry goes in AC10's conditions; the freeze may name a different B model.
3. V_W's latency counts against 300 s. A sonnet verifier over a small bundle is expected to take ≈30–60 s (unmeasured).
