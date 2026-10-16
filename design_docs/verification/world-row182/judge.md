# Row 182 judge r1 — design/row182-episode-outcome-capture (PR #280)
Score 68/100 — FAIL (2 blocking, both small). Go pkg: ok (65s). Python: 194 tests OK (skipped=5); no ERROR seen (1 run).

## AC check
- AC1 PASS: outcome.go:147 canonical Marshal, :155 object, :213-225 commit; outcome_test.go:143-195 (byte-equal, find, entry writtenBy harness:floor, entry 1).
- AC2 PASS w/ nit: refusals outcome.go:86-146; outcome_test.go:42-92. Range check allows to==from-1 (outcome.go:124, tested :80), a deliberate deviation from literal "from<=to". Size-cap test is vacuous (N1).
- AC3 PASS: 401 -> exit 2 (test :214-219), uncertain -> 3 via seam (:227). The real 503 / 200-mismatch mapping is covered by row-150 TestDemoProofAdapterGuards (X6/X7 killed there).
- AC4 PASS: demo_proof.go:185-252, demo_seed.go:264-272. Error strings "proof writer: ..." are identical via op="proof"; Op "proof commit" identical. Row-150 tests unmodified and green.
- AC5 FAIL: no test calls run_cell; the hook (run.py:746) is untested.
- AC6 PASS: outcome_test.go:187-194 (why head, why <ref>). `provenance --episode` not exercised.
- AC7: nothing in the code blocks it (verb + hook + knob present; session passed as a file path, no token on argv). But B2 would produce record errors on non-ASCII output in the smoke.

## BLOCKING
B1 AC5 untested. run.py:746-747 (hook). Mutations P4 (hook ignores client/token), P5 (hook deleted), P6 (record before classify), P8 (knob ignored) ALL SURVIVE the full suite; MU-3 and MU-5 are undetectable. Fix: add a run_cell test (fake Ctx, run_task.run_agent, grade) or extract the hook into a function, asserting: knob off -> no outcome_record; client None -> none; recorder error -> row['class'] unchanged and outcome_record_errors counts it; recorder sees row['class'] already set (MU-5); both arms.
B2 Byte/char bound mismatch. run.py row builds grade_stdout_head = stdout[:400] (CHARACTERS); outcome.go:134 refuses len() > 400 BYTES. 400 x 'é' = 800 bytes -> exit 1 -> {'error'}, counted in outcome_record_errors, which a FINAL treats as a harness defect. Any task printing non-ASCII (or U+FFFD, 3 bytes) trips it. Fix: count runes in Go (utf8.RuneCountInString) or truncate to 400 encoded bytes on a rune boundary in outcome_payload; add a test with 400 multibyte chars.

## NON-BLOCKING
N1 Size test (outcome_test.go:89) uses 16K+1 spaces, refused by JSON decode not by the cap; X3/X3b survive. Use a structurally valid >16 KiB payload (20 argv items x 1000 B).
N2 P1 (drop returncode check) survives; unreachable in practice.
N3 `to == from-1` accepted; document in design D2.
N4 Not verified: world.check_task_run / canary over a log containing harness:floor entries. Reading run_task.py:282, the next task's log_from = head+1 is read after the record, so ranges are not shifted. Confirm in the AC7 real-daemon smoke.
N5 Secrets OK: session file path only; env strips AILANG_REGISTRY_API_KEY / WORLD_SESSION (P2 killed).

## Mutations
| id | where | result | by |
|---|---|---|---|
| MU-1 world entries nil accepted | outcome.go:121 | KILLED | TestOutcomeParseRefusals (panic) |
| MU-1b range check off | outcome.go:124 | KILLED | TestOutcomeParseRefusals |
| MU-2 no canonicalization | outcome.go:147 | KILLED | TestOutcomeCanonical |
| MU-3 swallow failure into class | run.py hook | NOT TESTABLE | no run_cell test |
| MU-4 writtenBy demo | outcome.go:31 | KILLED | RoundTrip, UncertainExit |
| MU-5 record before grading (P6 form) | run.py:745-746 | SURVIVED | - |
| X1 allow unknown fields | outcome.go:90 | KILLED | ParseRefusals |
| X2 skip trailing-data check | outcome.go:95 | SURVIVED (equivalent, :101 also refuses) | - |
| X3/X3b no size cap / cap x100 | outcome.go:86 | SURVIVED | test uses spaces |
| X4 shell with entries ok | outcome.go:117 | KILLED | ParseRefusals |
| X5 bad arm accepted | outcome.go:129 | KILLED | ParseRefusals, Refusals |
| X6 503 CommitUncertain unmapped | demo_proof.go commitOneObject | KILLED | TestDemoProofAdapterGuards |
| X7 selectedHead check dropped | commitOneObject | KILLED | TestDemoProofAdapterGuards |
| X8 interpreter "" branch inverted | commitOneObject | KILLED | Guards, RoundTrip, Refusals |
| X9 pinned-interpreter check removed | commitOneObject | KILLED | Guards |
| X10 exit 3 -> 2 | outcome.go:220 | KILLED | UncertainExit |
| X11 head length check off | outcome.go:134 | KILLED | ParseRefusals |
| X12 returns entry 0 | commitOneObject | KILLED | RoundTrip |
| P1 drop returncode check | run.py:~806 | SURVIVED | - |
| P2 no env filter | run.py record_outcome | KILLED | test_success_parses_ref_and_entry |
| P3 spawn error raises | run.py:~805 | KILLED | test_failure_is_an_error_never_a_raise |
| P4 hook ignores client/token | run.py:746 | SURVIVED | - |
| P5 hook removed | run.py:747 | SURVIVED | - |
| P6 hook before classify_row | run.py:745 | SURVIVED | - |
| P7 inputFiles unsorted | run.py outcome_payload | KILLED | test_input_files_are_named |
| P8 knob ignored | run.py:746 | SURVIVED | - |
Worktree restored after each mutation (git status clean).

## r2 (head 0da335e)
Score 91/100 — PASS. No blocking findings. Go pkg ok (50s); Python 201 tests OK (skipped=5), no ERROR. Worktree clean after mutations.
- B1 fixed: hook extracted as maybe_record_outcome (run.py ~812), called after classify_row; Hook tests cover both arms, knob off/absent, no client, no token, recorder error and raise, class untouched, errors counted.
- B2 fixed: _head_bytes cuts to 400 UTF-8 bytes on a char boundary (surrogate escapes -> U+FFFD first); tested in Python and against the real Go validator.
- N1/N2/N3 fixed.
Re-run mutations (all KILLED): P1 (Record.test_failure_is_an_error_never_a_raise), P4 client / P4b token / P8 knob (Hook.test_off_or_no_client_or_no_token), P5 hook call removed and P6 hook before classify (Hook.test_run_cell_calls_the_hook_after_classify), MU-3 no except / MU-3b record overwrites class (Hook.test_failure_never_changes_class), MU-5 = P6, B2 uncut head (HeadBytes x2 + real-validator test), X3 no cap and X3b cap x100 (TestOutcomeParseRefusals).
Non-blocking: P5/P6 are killed only by a source-text inspection of run_cell (index of two strings), not by executing run_cell; a behavioural run_cell test would be stronger. Not exercised: provenance --episode (AC6) and real-daemon canary over outcome entries (defer to AC7 smoke).
No new defect found in the fixes (ctx.cfg/ctx.L exist on Ctx; tool_bin(ctx.cfg) same as before).
