# Evaluator report, World iteration 304 (docs-only landing of 485d42f0 + 00f7baf9)

**VERDICT: PASS (82/100) for landing these docs-only commits.** The parent design doc stays quorum-BLOCKED (r5, r6), so nothing may execute from it. See item 4.

## 1. Designer accuracy (485d42f0)
- V29 re-run. `Google Chrome --version` printed `Google Chrome 156.0.8078.47`. `ls -la` shows mtime `Oct 20 22:47`. `git diff --quiet HEAD -- scripts/marketing_capture.py` gave rc 0. `38ac4167`, `98f48b20` and `c73c7904` all resolve to the cited subjects. I did not re-run the shot/PNG decode. That is the unverified part of V29; its claimed hashes (`7e96e882…`, `d1240c97…`) and timings are unconfirmed.
- V30 re-run. The cited text is present at HEAD: `evidence.ail:4`, `:30`, `:32`; `workbench_packet_evidence.go:265-267`; `demo_seed.go:274/275/1017`; `world_demo_test.go:221`; `28ff4db3` and `bcd2d7cd`. The mapping is real and the V30 reading is correct.
- V31 re-run. `merge-base --is-ancestor` for `38db793`, `9319e6b` and `b8b6a3f` against `c21f3031` gave rc 0, rc 0, rc 0. Note the doc's own control claims "HEAD 039bfa77 → rc 1". I ran it with the sha as the descendant and got rc 0, because I tested the SHA against HEAD, not HEAD against the base. This is a harmless wording ambiguity in the control, NON-BLOCKING.
- "Verbatim". The M7 fallback sentence and the Acceptance sentence match Kimi r5 `proposed_fix` character for character in substance. I checked both against `quorum-r5.json`. The designer did NOT apply Kimi's second suggestion verbatim: V29 is annotated by naming the committed supervisor instead of "never committed". That is correct, because the file is committed.
- Loosening. The Chrome rule changes the contract from "fixed Chrome 155" to "same `--version` string before A and after B, major ≥155, recorded". It is bounded and keeps raw-byte equality between passes of one run, so it is not pixel similarity. The M7 fallback is tightened (M7 fails, escalate to Mark). The only slack is that there is no upper major bound, which is acceptable because the string is recorded. NON-BLOCKING.
- Diff scope: the doc diff is exactly header, annotations, V29-V31, the §2 rule, the roster sentence, the M7 row, the acceptance sentence and the r5 paragraph. That matches the commit message.
- **NON-BLOCKING defect INTRODUCED by this iteration:** the §2 heading "(measured on 155 and 156)" (line 217) overstates the evidence. On 156 only a synthetic `data:` page was measured (V29 says so itself). Kimi r6 is right. Fix: "system Chrome, major ≥155".
- **NON-BLOCKING, introduced by omission:** the `Status:` line still says "REVISED r4 (iteration 266)", and `Base:` still says `c21f303` while V29-V31 are measured at `039bfa77`. GLM and Gemini r6 both flag this.

## 2. Planner accuracy (00f7baf9)
All spot-checks pass (12 of 12):
- `refuse` is at `png.go:22`.
- `wantCode` is at `png_test.go:23`.
- The gate `args[0] != "check"` is at `main.go:13`.
- The usage `Fprintf` format is at `main.go:36`.
- `OpenFile(... O_EXCL)` is at `main.go:73`.
- `cli/unknown` with `[]string{"docs"}` expecting 2 is at `png_test.go:800`.
- `TestProductionGoSurface` allow-list contains `scripts/marketingcheck/main.go` at `host/store/context_roots_test.go:205`. The new file `observation_test.go` is a `_test.go`, so no allow-list edit is needed, as the plan says.
- The `python3` mandate is at `supervisor_test.go:28-30`.
- The `go build` 60 s pattern is at `png_test.go:734-741`.
- go.mod is `go 1.26.6`, so `slices` is available.
- The C11 Python one-liner exists verbatim in the plan at about line 120. Run with `python3 -I -c`, it takes 0.03 s and emits 12 records, which matches the planner's P7 count. The 30 s bound is therefore ≥1000×. The 60 s build and 120 s test bounds are comfortable.
- The Python `sorted(Path)` order can differ from Go string order, for example `a.sh` vs `a/b`. The plan's parity test re-sorts, and `CompareObservations` sorts itself, so this is harmless.
- Mutants, reasoned through (not prototyped):
  - Digest guard to `if false` gives nil on `same-line-count-change`, so `wantCode` reds. OK.
  - Line guard to `if false` gives nil on `digest-equal-line-change`, so it reds. OK.
  - Set guard to `if false` with a path-keyed loop gives `observation-digest` on `deleted`, so `wantCode` reds on the wrong code. OK.
  - NON-BLOCKING ambiguity: if the executor writes a position-indexed loop, the third mutant panics instead of reaching the assertion. The plan does say "zero-value record", which implies a map lookup. An executor note saying "look up by path in a map" would remove the doubt.
- M7c1 is implementable: two files, about 240 lines, 20 minutes. There is a minor overrun risk for a small model.
- Chrome rule across children: M7c2a uses `MIN_CHROME_MAJOR = 155` (plan line 1147). Its mutant sets it to 154 and the fake-154 case reds. M7c2b runs `chrome-version` before A and after B, refuses on A≠B, and writes `diagnostics.txt`. M7c4 repeats the rule explicitly (design §2 / V29). All match the revised §2.
- JSON validity and agreement: `json.load` OK. The chain is M7b→M7c1→M7c2a→M7c2b→M7c3→M7c4→M8a, matching the markdown. Estimated minutes are 20/28/28/25/25, matching the markdown table. The JSON files lists agree with the markdown (M7c1 main.go and observation_test.go; M7c3 main.go and marketing_inventory_test.go; and so on).
- Consequence for dependent rows: the M8b re-drill row references the M7c1 mutant by name, consistent.

## 3. Base-gate claim
`git diff --stat 039bfa77 00f7baf9 -- '*.go' '*.py' '*.sh'` is EMPTY, so the failing tests are pre-existing at the base by construction.
My runs (sandboxed tree; `uptime` load average 31-34 at the time, far above the planner's 8-15):
- Package run: rc 1. `TimeoutKillsGroup`, `cleanup-overrides-shot` and `shell-hup` all failed. The first two failed with `[Errno 1] Operation not permitted`.
- Single test, run 1: rc 1. `interrupted-shot` failed with the same EPERM.
- Single test, run 2: **ok** (7.1 s).
- Single test, run 3: rc 1. `shell-hup` failed with exit 1 instead of 129 (0.59 s). This reproduces the planner's observation.
Conclusion: `shell-hup` exit 1 reproduces; the planner's "isolated runs green" holds only intermittently. The EPERM variant is probably an artefact of my sandbox (a signal/process-group call being denied) and I cannot separate it from load. NON-BLOCKING for landing a docs-only change. However, the plan's claim that the carried flake is "only `shell-hup`" and "3/3 green" is load-dependent, and the M7c1 executor should be told the supervisor tests are flaky at base (report it, do not fix it, do not count it as a regression). The plan already says to carry it.

## 4. Gate disposition
Rules read in `gate-2-pick.md` (QUORUM-AT-PICK, NARROW-REFINEMENT CARVE-OUT, the absent-reviewer rule, "WHEN A DOC BLOCKS ROUND AFTER ROUND").
- The earlier pick-time PROCEED with GLM and Kimi absent is "a pass with a named hole", not a pass. The closed eyes then rejected in both r5 and r6, which is exactly what the absent-reviewer rule warns about. Proceeding to an executor on that artifact would have been wrong.
- One revision and one re-quorum is the permitted budget; r6 is the re-quorum and it BLOCKED. So the default is `needs-human-review`/park, unless the carve-out or the sequence rule applies.
- Carve-out does NOT apply. Both r6 objections lack a narrow verbatim fix. GLM asks for the authority under which M1a, M5a, M7a and M7b merged ahead of a passing quorum, which needs controller judgment. Kimi asks for new Chrome 156 real-content measurements, which is a measurement and not a text edit. Both also challenge the doc's frame (stale r4 vs as-built).
- The sequence rule fits. Reviewer objections moved from different surfaces (r5: CLAIMED premise, fallback authority, dependency SHAs) onto one surface in r6: the doc is a frozen-at-r4 plan that now documents landed work. Gemini passed in both rounds. Rule (b) says SPLIT, with the surviving objection carried verbatim, and rule (d) says a split is a routing call and NOT `needs-human-review`. So the controller's disposition (no executor; route M7c-M8 into a new, smaller as-built doc anchored at 039bfa77, authored and quorumed next iteration) follows the rules, and is the better choice over parking.
- Caveats, NON-BLOCKING: (a) Gemini never "flipped" to pass; it passed throughout, so the "flip" trigger is only loosely met and the justification should be recorded as "objections localised onto the stale-baseline surface". (b) Rule (e) requires the round count in the report: this doc is past round 4 (r1-r6). (c) The plan commit (M7c1-M7c4) is anchored on a still-BLOCKED parent. When the as-built doc lands, the plan JSON/markdown may need re-pointing. (d) M7c1 is independent of the disputed premises (it only needs the protected-file comparator), but running it is nevertheless a controller call; the conservative no-executor choice is defensible.
- Objections that this iteration INTRODUCED rather than inherited: Kimi's "measured on 155 and 156" heading (introduced); the unchanged stale `Status: REVISED r4` and `Base` lines (introduced by V29-V31 being added at a new HEAD without re-anchoring); the self-contradiction between "changes design_docs only"/"future work" and the annotations (introduced as annotations, though correct). Pre-existing: the `Planner-Lane: codex-ok after quorum/freeze` contradiction with already-merged M1a-M7b, and the V1-V28 measurement base `c21f303`. By rule (c), the pre-existing ones belong in the next doc or a queue row, not another revision of this one.

## Score
Accuracy 90 (the measurements reproduce), planner precision 85, introduced heading overstatement and unrefreshed header -8, base-test flake characterisation -5, unrun mutants -5. Total 82.
