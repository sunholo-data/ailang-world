Updated only the [sprint plan](/Users/voightkampff/dev/sunholo-data/.planner-wt-iter282/design_docs/planned/w-exec-protected-test-paths-sprint-plan.md) and [JSON](/Users/voightkampff/dev/sunholo-data/.planner-wt-iter282/design_docs/planned/sprint_w-exec-protected-test-paths.json).

Changed sections: authority/status, baseline provenance, milestone table/dependencies, M1c/M1d, M4a/M4b checker wiring, M6-pre/V24, M6a/M6b/M6c, and corresponding JSON gates, caps, risks and coverage. JSON schema and all four done milestone objects remain unchanged.

The affected table rows, verbatim:

| ID | Title | Changed LOC total | Active minutes | Dependency |
|---|---|---:|---:|---|
| M1c | Pure candidate renderer and mirror filter | 250 | 30 | M1b |
| M1d | Real-Git candidate and mirror equivalence | 240 | 30 | M1c |
| M2a | Load the additive protected profile key | 200 | 25 | M1d |
| M6-pre | Measure V24 and size phase/construction bounds | 200 | 30 | M5b, M1d |
| M6a | Capture the immutable episode baseline and bounded record shape | 245 | 30 | M6-pre |
| M6b | Attach before/after integrity while preserving the command deadline | 245 | 30 | M6a |
| M6c | Kill tamper, ignored-file, flag and moved-HEAD mutants; pin replay | 245 | 30 | M6b |

**M1c tests:** `TestProtectedFoldToASCIIClosed`, `TestProtectedCandidatePathspecsGolden`, `TestProtectedCandidateParseFailsClosed`, `TestProtectedIgnoredBaseline`, and re-derived `TestProtectedRenderers/git`.

| Mutation | Exact killer |
|---|---|
| MUT-WIDEN-SET-STALE | `TestProtectedFoldToASCIIClosed/closed-set` |
| MUT-CAP-SILENT | `TestProtectedCandidateParseFailsClosed/over-cap-ignored` |
| MUT-NO-IGNORED-BASELINE | `TestProtectedIgnoredBaseline/standing` |
| MUT-PARSE-MALFORMED-OPEN | `TestProtectedCandidateParseFailsClosed/missing-nul` |
| MUT-PARSE-STATUS-OPEN | `TestProtectedCandidateParseFailsClosed/unknown-status` |
| MUT-PARSE-FLAG-OPEN | `TestProtectedCandidateParseFailsClosed/bad-flag` |
| MUT-PARSE-FIELDS-OPEN | `TestProtectedCandidateParseFailsClosed/bad-fields` |

**M1d test:** `TestProtectedGitCandidateEquivalence`, retaining every fixture row plus the complete design bank, including ß/ﬀ controls and both CWDs. All killers below are prefixed with `TestProtectedGitCandidateEquivalence/`.

| Mutation | Killer subtest |
|---|---|
| MUT-ASCII-ONLY-SUPERSET | `long-s/root-CWD`, `kelvin/root-CWD`, `greek-question-mark/root-CWD` |
| MUT-QMARK-KEPT | `qmark-utf8/root-CWD`, `class-utf8/root-CWD` |
| MUT-NO-MIRROR-FILTER | `non-hit-sibling/root-CWD` |
| MUT-IGNORED-MATCHING | `ignored-nm-c-txt/root-CWD` |
| MUT-NO-TOP | `root/subdirectory-CWD` |
| MUT-NO-BARE-D | `directory-replaced-by-file/subdirectory-CWD`, `directory-replaced-by-symlink/subdirectory-CWD` |
| MUT-PATHSPEC-NO-ICASE | `pkg-A_TEST.GO/root-CWD` |

**M6-pre test:** `TestProtectedV24RealWorktree`.

| Mutation | Exact killer |
|---|---|
| MUT-V24-NULL-PASS | `TestProtectedV24RealWorktree/reject-null` |
| MUT-V24-BUDGET-PASS | `TestProtectedV24RealWorktree/reject-over-budget` |
| MUT-V24-NO-IGNORED | `TestProtectedV24RealWorktree/required-legs` |

M6 now specifies three candidate streams through `parseProtectedCandidates`, immutable `base_ignored_digest`/`base_ignored_files`, 2500 ms phase and 3200 ms construction bounds subject to V24 resizing, and the recomputed child reserve. `TestExecProtectedAmbientIgnoredIntact/standing` kills MUT-NO-IGNORED-BASELINE. M6a explicitly cannot start before measured V24 and resized bounds.

**Measured baselines at `451dc3d`:**

| Gate | Result |
|---|---|
| `go vet ./host/broker/` | GREEN, rc 0 |
| `go test -count=1 ./host/broker/ -run '^TestProtected' -timeout 300s` | GREEN, rc 0; package time 0.495 s |
| `jq -e . design_docs/planned/sprint_w-exec-protected-test-paths.json` | PASS, rc 0 |
| `git diff --check` | PASS, rc 0 |

Both Go gates used unset `WORLD_TOOL_AILANG_BIN` and worktree-local TMPDIR/GOTMPDIR. Scratch was removed.

**UNMEASURED:** V24, new-test standalone gates, mutation drills, repository-wide vet, and unfiltered broker/daemon/CLI suites. The broker suite was not run because PF6 identifies legacy real-HOME writes. New tests are **absent at base; zero tests is not a pass**. Socket-touching checks are **UNINFORMATIVE UNDER SANDBOX**. No network, Git writes, fixture repositories, production edits or test edits occurred.