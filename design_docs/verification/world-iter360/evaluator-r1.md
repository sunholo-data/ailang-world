# Iteration 360 — independent evaluator (Agent `sonnet`), row 255 M1 at fa5b9183

Spawned foreground by the controller against an isolated `git archive fa5b9183` copy
(`.judge-iter360`, outside the sprint worktree, because the A5 mutation runner was editing the
worktree in place at the time). The directive gave questions and files, not the controller's
numbers. 13 tool calls, 566 s, 54,802 tokens. The report is reproduced verbatim below.

---

**Score: 88/100. PASS.**

Correctness 27/30, test quality 18/20, design conformance 18/20, scope discipline 14/15, evidence 11/15. I could not read the A5 mutation matrix (not in my copy), so I did not judge it.

I did not open the design doc beyond the row text, so Q8 is a code-versus-claims comparison only. `go test` on `host/...` and `cmd/ailang-worldd` passed everywhere except `host/pins`, which failed on `git ls-files` because my copy has no `.git`.

**Q1. Surfaces and tests: all carry the line, none pinned vacuously.**
- Stored record: `broker.go:298-315` sets `Failure: reason` on the failed record. `TestExecFailureClassesReachTheRecord` and `TestFailureLineMarksTheStoreFirst` pin it. Dropping the field turned both red.
- MCP/A2A result: `effectful.go:199` copies `failed.Reason` into the world block's `reason` field. `TestEffectFailureReasonReachesTheWorldBlockAndReplays` pins it. Dropping the copy turned it red.
- `why`: `why.go:912-914` prints `no result: <line>`. `TestWhyShowsTheFailureLine` pins it. Reverting the line turned it red.
- Log: `mcp.go:135` and `projection.go:379` call `logEffectFailures`. Dropping the MCP call turned `TestMCPEffectFailureLineReturnedAndLoggedOnce` red.
- `call` CLI: `TestCallPrintsTheFailureLine` went red when I dropped the line from `mcpclient.go`.

**Q2. Secrecy and bounding: no leak path found.**
- I checked every `execFailure` and `HandlerFailure` construction site in `handlers_exec.go` and `failure.go`. The error text never reaches the line.
- Details come from fixed words, the profile root, cache names and the command id. The profile root must be a clean relative path (`exec_profile.go:196`) and cache names are validated. Only the base name of any path outside the disclosed set (the node/argv0 paths already returned by a successful exec) is named.
- `seedWhat` uses the error's type and exit code only, so seed paths and child output do not reach the line.
- `driftWhat` and `settingsWhat` map to fixed phrases; an OS error text becomes "unreadable".
- `FailureLine` replaces invalid UTF-8, control characters and Cf (bidi) characters, collapses whitespace and cuts at a rune boundary to 200 bytes. An unknown class gives "".
- `ValidFailureLine` checks that the line is a fixed point of the rule. `DecodeRecord` enforces it through `FailureConsistent` (`record.go:84`), and `putRecord` enforces it on write (`broker.go:411`).

**Q3. Compatibility.**
- The `failure` field is `omitempty`, so records without a line keep their bytes. `TestRecordFailureGoldenBytes` and the existing golden tests are green.
- Old records decode, because the field is optional.
- The mark runs before `putRecord` (`broker.go:306-311`). A failed mark returns `OutcomeWriteError` and writes no record, with the debit standing.
- `MarkFailureLines` is only called from `invoke`. The only lined-record writer is `putRecord`, after the mark. Marking also fails on a read-only handle (tested). A store at 6 opens under the new binary. A binary that knows only version 5 refuses a version-6 store at open.

**Q4. Failure, refusal, denial.**
- Refusals still go through `execRefusal` and are unchanged.
- `Execute` only wraps its existing errors, so no call changed category.
- The time-budget class is named `time-budget-spent` with the text "the call's time was spent", which does not collide with the capability-budget denial.

**Q5. Replay.**
- `invokeReplay` sets `Reason: rec.Failure`.
- The replay test compares the replayed bytes to the committed output exactly.
- The world block's `reason` is `omitempty`, so old committed outputs replay unchanged.
- Existing replay tests in `host/replay` and `host/coordinator` pass.

**Q6. Mutations (smallest change that should break each test).**

| Mutation | Test | Result |
|---|---|---|
| `why.go`: drop the failure line | `TestWhyShowsTheFailureLine` | red |
| `effectful.go`: drop the reason copy | `TestEffectFailureReasonReachesTheWorldBlockAndReplays` | red |
| `broker.go`: drop `Failure:` from the record | `TestExecFailureClassesReachTheRecord`, `TestFailureLineMarksTheStoreFirst` | red |
| `broker.go`: skip the `MarkFailureLines` call | `TestFailureLineMarksTheStoreFirst` | red |
| `mcp.go`: drop the log call | `TestMCPEffectFailureLineReturnedAndLoggedOnce` | red |
| `record.go`: drop the decode check | `Failure\|Record` pattern | red |
| `mcpclient.go`: drop the line from `call` | `TestCallPrintsTheFailureLine` | red |
| `failure.go`: cut at 200 instead of 197 | `TestFailureLineBoundsAndSanitises` | red |
| `failure.go`: remove the Cf sanitiser | `TestFailureLineBoundsAndSanitises` | red |
| `failure.go`: raise the keep-whole threshold from 200 to 250 | `TestFailureLineBoundsAndSanitises` | green (see NB1) |

I restored every file and confirmed each with a byte compare.

**Q7. Scope.** The diff touches `docs/QUICKSTART.md` plus the broker, coordinator, projection, store and `cmd/ailang-worldd` files. It touches no `tools/launchd`, skill or `.ail` file, and no row 256/257 code.

**Q8.** I found no claim in the code that the design lacks, and no design claim the code fails to do.

**Findings (all NON-BLOCKING)**
- **NB1.** `host/broker/failure.go:~96`: the bound test does not cover lines of 201 to 250 bytes. With the keep-whole threshold raised to `FailureLineMax+50`, `TestFailureLineBoundsAndSanitises` stays green because it only checks a 300-rune input. A line that long is not reachable in production, since every detail is short. Reproduce with `sed 's/len(line) <= FailureLineMax/len(line) <= FailureLineMax+50/'`, then `go test -run '^TestFailureLineBoundsAndSanitises$' ./host/broker/`. A test with a detail of about 210 bytes would close it.
- **NB2.** `handlers_exec.go:~268`: the misrouted class puts `req.Effect@req.Scope` verbatim into the line. The text is bounded and sanitised, and these are effect and scope names, not paths, so I could not construct a leak (UNVERIFIED for hostile registry names).
- **NB3.** A failed schema mark leaves an intent with no outcome (`OutcomeWriteError`, debit stands). I accept this as designed, but it is only as safe as the later reconciliation of that intent.
- **NB4.** `TestReleaseLiteralsExplainedByManifest` fails in my copy only because there is no `.git`. It is not a sprint defect.

**PASS**

---

## Controller follow-through on the judge's findings

- **NB1 reproduced** at fa5b9183 in the judge copy: the `FailureLineMax+50` mutant → `TestFailureLineBoundsAndSanitises` rc 0 (survives); restored, `cmp` identical, test rc 0. The judge's correction of its own aside stands: the line `len(line) <= FailureLineMax` is at `failure.go:84`, not ~96. Logged, not filed (D-WORLD-109).
- **NB1 understates one thing the controller found separately** (A5 run 1, `controller-M1/summary.txt`): the same 200-byte cap makes `TestExecFailureClassesReachTheRecord/exec-stack-drift_node` red whenever the temp path is long, because its expected string is the uncut line. So the cap is reachable in a test fixture, not just "not reachable in production". Logged, not filed.
- NB2, NB3: accepted as stated (NB3 is the r3 design choice, §3.7 and Claim F). NB4: environment of the copy, not the sprint.
