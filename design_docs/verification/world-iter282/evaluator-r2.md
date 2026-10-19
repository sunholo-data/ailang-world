# Round 2 evaluation, row 156 (tree detached at ba544c0)

Score: 88/100, PASS.

## Delta reviewed (e4810dc..ba544c0)
- host/broker/exec_protected_test.go: only the three over-cap rows of TestProtectedCandidateParseFailsClosed changed. They now use well-formed NUL-terminated streams (`? a\0` x cap/4+1, `a\0` x cap/2+1, `H a\0` x cap/4+1). No other test changed, no production code changed.
- Design doc: the §4.5 measured-cost sentence now says the 165-195 ms figure is V21's clean shared clone only and cites V23 (406-504 ms warm and 2114 ms cold ignored leg on ailang). R-156-9 added: one-character widened patterns render as `**/*`, correct because the mirror filters, bounded by the 8 MiB cap which gives `error`, never `intact`.

## (a) Cap mutants against TestProtectedCandidateParseFailsClosed
| Mutant | Result |
|---|---|
| cap check removed | KILLED: over-cap-status, over-cap-ignored, over-cap-flags all FAIL (round 1: survived) |
| cap bound doubled (`> 2*max`) | KILLED: same three |
| over-cap truncated to the cap and accepted | KILLED: same three |

Unmutated baseline: ok. exec_protected.go sha e8602fdc... before and after (restored byte-identically).

## (b) Gates
| Gate | rc | Result |
|---|---|---|
| gofmt -l host/broker | 0 | empty |
| go vet ./... | 0 | |
| go test -count=1 -run '^TestProtected' -timeout 300s -v | 0 | 155 PASS, 1 SKIP (TestProtectedMatchFixturePinned, needs WORLD_TOOL_AILANG_BIN), 0 FAIL |
| same with -race | 0 | ok |
| jq -e . on sprint JSON | 0 | |

## (c) Re-score
| Category | Round 1 | Round 2 |
|---|---|---|
| Acceptance criteria | 22/30 | 27/30 |
| Test quality and non-vacuity | 10/20 | 18/20 |
| Design fidelity | 19/20 | 19/20 |
| Plan and JSON | 9/10 | 9/10 |
| Code quality and standards | 8/10 | 8/10 |
| Documentation (cost sentence fixed, R-156-9 added) | n/a | 7/10 |
| Total | 68 | 88 |

Round-1 BLOCKING 1 is resolved. NB2 and NB3 are resolved. Remaining NON-BLOCKING, unchanged:
- NB4 (opinion, no reproduction): the parser returns head "" with a nil error when the status stream has no `# branch.oid`. M6's head != base comparison fails safe.
- NB5: the controller chose the 2500 ms phase bound and 3200 ms construction bound itself, rather than applying oc-glm-5-3's split-budget wording. This is disclosed and deferred to V24.

Not measured: CI, the pinned-binary arm, a Linux git run, the unfiltered broker suite.
Tree: `git status --porcelain` empty, scratch deleted.

EVALUATION_RESULT: pass
EVALUATION_SCORE: 88/100
