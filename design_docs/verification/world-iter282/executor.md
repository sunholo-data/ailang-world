M1c and M1d completed in about 10 minutes.

| Milestone | Files changed | Changed LOC: prod/test |
|---|---|---:|
| M1c | `host/broker/exec_protected.go`, `host/broker/exec_protected_test.go` | 140 / 110 |
| M1d | `host/broker/exec_protected_git_test.go` | 0 / 196 |

LOC counts include additions and deletions. Cumulative snapshots were written and byte-verified in `.snap/M1c/` and `.snap/M1d/`. All 25 Git rows passed from both CWDs.

| Gate | M1c rc | M1d rc |
|---|---:|---:|
| `gofmt -l host/broker` — empty output | 0 | 0 |
| `go vet ./host/broker/` | 0 | 0 |
| `go test -count=1 ./host/broker/ -run '^TestProtected' -timeout 300s` | 0 | 0 |
| `go test -race -count=1 ./host/broker/ -run '^TestProtected' -timeout 600s` | 0 | 0 |

Mutation results below group separately executed killers. `Parse` means `TestProtectedCandidateParseFailsClosed`; `Git` means `TestProtectedGitCandidateEquivalence`.

| Mutation | Killer | Observed FAIL | Reverted SHA match |
|---|---|---|---|
| MUT-WIDEN-SET-STALE | `TestProtectedFoldToASCIIClosed/closed-set` | Yes | Yes |
| MUT-CAP-SILENT | `Parse/over-cap-ignored` | Yes | Yes |
| MUT-NO-IGNORED-BASELINE | `TestProtectedIgnoredBaseline/standing` | Yes | Yes |
| MUT-PARSE-MALFORMED-OPEN | `Parse/missing-nul` | Yes | Yes |
| MUT-PARSE-STATUS-OPEN | `Parse/unknown-status` | Yes | Yes |
| MUT-PARSE-FLAG-OPEN | `Parse/bad-flag` | Yes | Yes |
| MUT-PARSE-FIELDS-OPEN | `Parse/bad-fields` | Yes | Yes |
| MUT-ASCII-ONLY-SUPERSET | `Git/{long-s,kelvin,greek-question-mark}/root-CWD` | All yes | All yes |
| MUT-QMARK-KEPT | `Git/{qmark-utf8,class-utf8}/root-CWD` | All yes | All yes |
| MUT-NO-MIRROR-FILTER | `Git/non-hit-sibling/root-CWD` | Yes | Yes |
| MUT-IGNORED-MATCHING | `Git/ignored-nm-c-txt/root-CWD` | Yes | Yes |
| MUT-NO-TOP | `Git/root/subdirectory-CWD` | Yes | Yes |
| MUT-NO-BARE-D | `Git/directory-replaced-by-{file,symlink}/subdirectory-CWD` | Both yes | Both yes |
| MUT-PATHSPEC-NO-ICASE | `Git/pkg-A_TEST.GO/root-CWD` | Yes | Yes |

All 18 trials returned rc 1 at the intended assertion, then rc 0 after restoration. Production SHA-256 before/after each trial:

`1155b863ca6b7e80155c5013add76a62bc29453390f2600827deccddfe7bd0ce`

Implementation addition: a pure ignored-command constant enables MUT-IGNORED-MATCHING to alter production bytes and exercise the actual Git invocation without production I/O.

UNMEASURED: external guard before/after hash census, CI execution, optional pinned AILANG arm, and unfiltered broker suite (prohibited by PF6).

No git writes were performed against the mission worktree; fixture Git writes were confined to `t.TempDir()`. No dependencies were added. `.scratch-exec282/` is deleted.