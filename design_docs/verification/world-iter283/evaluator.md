# Eval iter-283 r1 — row 156 M2a (ebc6f0f) + M2b (eb46e92)
VERDICT PASS, 90/100. Correctness 38/40, tests/non-vacuity 25/30, scope 15/15, quality 12/15.
Gates green: named broker tests, TestExecStartupRefusalTable (51 rows), go vet, go test -race broker (164s). Changed files gofmt-clean. Worktree restored (git status empty).
Blocking: none.
Non-blocking:
1. D5: M2b "protected-malformed-value" expects substring `protected`, but the temp-dir path in the error contains the test name "...Tableprotected-malformed-value...", so the substring is satisfied even if the field is renamed (mutant `profileErr("prot", ...)` SURVIVED at the daemon row). The broker test (Field == "protected") does catch it, so net coverage holds; the daemon row's field-naming check is vacuous for that row. Fix: match on `profile: protected:` or avoid the id in the temp path.
2. Error text doubles the prefix: "exec profile: protected: protected[0] ...".
3. The dispatches==0 / summaries==0 / listener-less checks cannot fail under any realistic mutant (no profile is parsed on refusal, so nothing can dispatch); the archive/exec-sandbox stat is the check that actually catches a late refusal (D1). Probe marker/exec-probe stats are belt-and-braces.
4. M3 (drop element-type guard) is killed only by null-element (via an incidental entry-length error); the number element is caught by decodeField anyway. Guard is thin but covered.
(a) A refusal after probe: not constructible with the grammar error because the probe needs a parsed profile; a deferred-refusal mutant in StartExec was killed by the archive/exec-sandbox assertion. (b) Compatibility subtest uses literal v1 plus leading/trailing whitespace; TrimSpace-hash mutant killed, so raw bytes are pinned. (c) With M2a production hunk reverted: Protected/ProtectedDigest fields missing so the test file does not compile; semantically the unknown-key refusal would also fail list/wrong-type rows. Non-vacuous.
Mutations (restored by cp + shasum/cmp, green after):
| id | edit | test | result |
| MUT-PROFILE-NO-CHECK | ParseProtectedPatterns replaced by passthrough | UnknownFormRefused | RED exec_profile_test.go:295 "protected[0] names refused src/*.ts: err = <nil>" |
| MUT-PROFILE-TYPE (array guard) | raw[0]!='[' -> false | Protected/wrong-type | RED :267 (null) |
| elem guard | entry[0]!='"' -> false | wrong-type | RED :267 null-element |
| wrap dropped | return err unwrapped | UnknownFormRefused | RED :295 (not *ExecProfileError) |
| digest unset | ProtectedDigest="" | /list | RED :248 |
| MUT-PROFILE-COMPAT | hash TrimSpace(data) | /compatibility | RED :285 raw-byte digest |
| no canonicalization | p.Protected=patterns | /list | RED :248 |
| M8 equivalent (len>=0) | no-op | /empty | SURVIVED (equivalent mutant: protectedDigest(empty)=="") |
| MUT-STARTUP-PROTECTED-IGNORED | `if err != nil` -> `if false` for grammar error | Table/protected- | RED :482 does not name protected[0] (4 rows; probe ran) |
| D1 delayed refusal | StartExec defers profile error past archive+probe | Table/protected- | RED :489 reached sandbox construction (all 4 rows) |
| D4 | drop QUICKSTART protected-slash-glob row | Table | RED :435 census 50 vs 51 |
| D5 | malformed field name "prot" | Table/protected-malformed | SURVIVED (see NB1) |

## Round 2
Tree: 9bbcbd7 (detached). `git diff eb46e92 FETCH_HEAD` = docs/QUICKSTART.md one line only (nothing missed); no Go change.
Mutant `profileErr("prot", ...)` vs `-run TestExecStartupRefusalTable/protected-malformed-value`: now RED at exec_startup_test.go:482 (refusal does not name `protected: must be an array of strings`). Restored (git status empty); full TestExecStartupRefusalTable green (20.3s).
NB1 CLOSED. NB2 (doubled prefix), NB3 (unfireable dispatch/summary checks), NB4 (thin element guard) remain OPEN, non-blocking, unchanged.
VERDICT PASS, r2 score 92/100 (tests 27/30).
