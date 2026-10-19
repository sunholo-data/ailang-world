# Independent evidence evaluation — round 1

EVALUATION_RESULT: pass
EVALUATION_SCORE: 83/100
EVALUATION_ROUND: 1

Judged SHA: `d8e17680b642b34736ce1e57c74c73aa54c07d60`.
Comparison base: `a0c73053f10092fd511ad22fcf6b95d6a3f896aa`.
Actual judge: **gpt-6.1-sol / OpenAI, fallback FLAGGED** under the dispatched D-WORLD-48 exception. Dispatcher reported configured MiniMax containment failure and Sonnet/Opus weekly-quota 429s; this judge did not independently re-probe those providers. Fresh-context independent judgment, no prior judge evidence or fixtures reused. All execution and report writes used `/Users/voightkampff/dev/sunholo-data/.eval-world-iter281-fallback`, detached at the exact SHA. Initial required instruction reads occurred in the initially assigned tree before the controller switched trees; no fixtures/evidence or writes there were used.

**PASS applies only to the evidence/refutation deliverable. M1c remains incomplete and its equivalence premise is false on the measured host. This does not approve a replacement algorithm, implementation, merge, or dependent milestones.** Full tests, race gates, lint/vet and CI are **UNMEASURED by this judge**; PR309 CI belongs to the controller.

## Findings

- **Blocking design premise, correctly disclosed:** Source `foldProtectedKey` performs NFC plus Unicode SimpleFold orbit minima; `matchProtected` therefore requires the existing long-s and Kelvin fixture answers true. `gitPathspecs` emits literal top/glob/icase forms, including both bare D and D/**. Independently, real Git 2.54.0 selected ASCII counterparts but missed `teſts/x.py` with `tests/**` and `pkg/K_test.go` with `*k_test.go`, in both status and ls-files. Both patterns are ASCII; restricting pattern bytes cannot repair these counterexamples. The new status/V16/R-156-7/M1c annotations correctly state this and preserve semantics.
- **No high/medium evidence defect found:** Four separate repos avoid APFS alias interactions. Unfiltered index/status controls and exact literal Unicode recovery rule out absent paths, empty instrumentation and display quoting. Narrow mirror and renderer tests pass. Structured recorded result payloads decode consistently, including every base64/NUL-field representation. Changed files are only evidence plus factual design/plan annotations; no production or test implementation was smuggled in.
- **Low, wording:** README ends “No production, tests, design, sprint status, driver or skill files were changed.” As a measurement-script claim this is reasonable, but as an aggregate statement about the judged diff it is inaccurate: the design and plan receive factual annotations. Label it explicitly as the script/executor scope. This does not invalidate the measured counterexample.
- **Low, harness maintainability:** `repro.py` main is longer than 50 lines. Its selected-output classifier treats any non-exact selection as false; stored Unicode outputs are actually empty and independent reproduction confirms emptiness, so this is not an observed validity failure. A future characterization should distinguish unexpected partial output from empty selection explicitly.

## Independent method and measured checks

Every shell used `export PATH=/opt/homebrew/bin:$PATH` and pinned the fallback cwd. No generator run was counted as independent reproduction.

1. Read CLAUDE.md fully first, coding standards, evaluator skill/rubric, relevant charter/DESIGN boundaries, design and companion plan with M1c/r3/PF6, and current matcher/renderer/fixture source. Reviewed the complete base..HEAD change set, including structured audit of the 1,130-line result JSON rather than relying on a truncated diff. Commands: `git rev-parse HEAD`, `git diff --stat BASE..HEAD`, `git diff BASE..HEAD`, focused source/doc reads, and Python JSON/base64 recursive validation. `git diff --check BASE..HEAD` exited 0. BASE denotes the exact base above.
2. Exact authorized narrow test command (rc0; package 0.346s):

```sh
env -u WORLD_TOOL_AILANG_BIN HOME="$PWD/.judge/home" TMPDIR="$PWD/.judge/tmp" GOTMPDIR="$PWD/.judge/tmp" GOCACHE="$PWD/.judge/gocache" GOMODCACHE=/Users/voightkampff/go/pkg/mod GOPROXY=off go test -count=1 ./host/broker/... -run '^TestProtected(MatchFixture|Renderers)$' -timeout 120s -v
```

All selected tests were discovered and PASS, including long-s, kelvin, nfc and renderer/git; no optional pinned-binary arm was selected. HOME, temporary directories and build cache were worktree-local. Existing module cache was read with downloads disabled; no real HOME state was seeded.
3. Own alternate implementation: `python3 .judge/independent.py > .judge/independent.json` (rc0). For each path in a fresh worktree-local repo: `git init -q`; write new bytes; `git hash-object -w --stdin` of old bytes; `git update-index --add --cacheinfo 100644 BLOB PATH`. **No fixture commit was needed:** an index-only tracked AM entry and porcelain v2 replace the generator’s committed/v1 method. Commands use `git -c core.fsmonitor=false -c core.ignorecase=false ls-files -v -z` and `status --porcelain=v2 -z --untracked-files=all --ignored=matching`, first unfiltered, then `--` and independently source-derived specs. Git environment clears inherited GIT_*; HOME/TMPDIR are under `.judge`, global/system configuration is disabled and LC_ALL=C.

| Row | Unfiltered index/status | Selected index/status | Exact literal recovery |
|---|---|---|---|
| tests/x.py | populated / AM | populated / AM | yes |
| teſts/x.py | populated / AM | empty / empty | yes |
| pkg/k_test.go | populated / AM | populated / AM | yes |
| pkg/K_test.go | populated / AM | empty / empty | yes |

Directory specs were `:(top,glob,icase)tests` plus `:(top,glob,icase)tests/**`; basename spec was `:(top,glob,icase)**/*k_test.go`. Exact recovery used `ls-files -v -z -- :(top,literal)PATH`. All commands succeeded with 10-second subprocess limits. Separate fixture dirs were automatically removed.
4. Own instrument-health mutants: feed empty unfiltered results to the independent validator → `unfiltered-index-health` fails on all four rows; suppress selected ASCII results → `ascii-selected-health` fails on both controls; replace Unicode selected empty results with populated unfiltered results → `unicode-empty-selection` fails on both probes. These are deliberately **instrument-health data mutants**, not M1c production-mutant coverage.
5. Additional generator instrument health: copied repro into `.judge/repro-mutant.py`, redirected ROOT/OUT into this evaluation and changed only directory specs to missing/missing/**. `python3 .judge/repro-mutant.py` exited nonzero at `ASCII selected positive control missing: ascii-directory`; no stale result appeared. Isolated copy restored with SHA256 `183961e7b8836a0a604efa47c75d28e1bfe5fbeb612c50155bccc3f4cc610a70` verified equal. No production/design bytes mutated. Baseline independent method and narrow test remained green.

Retained local artifacts: `.judge/independent.py`, `.judge/independent.json`, `.judge/narrow.log`, restored `.judge/repro-mutant.py`. These are evaluator scratch inside the isolated tree, not additions to the judged commit. Fixture Git writes were confined there; checkout index/HEAD were never written.

## Explicit 100-point evidence-only rubric adaptation

The assigned deliverable is a reproducible falsification with honest annotations, not implementation of M1c. Accordingly acceptance/docs/fidelity are judged against that deliverable. Pending product milestones are required honesty, not failed completion paperwork. No performance or shared compiler change triggers a conditional category.

| Category | Score | Evidence and limits |
|---|---:|---|
| Tests / verification (20) | 15 | Named narrow suite and independently controlled counterexample measured; five points withheld because full test/race gates are UNMEASURED. This is not a full-test green. |
| Lint clean (10) | 0 | Lint/vet UNMEASURED locally; diff whitespace check is not lint. |
| Acceptance (30) | 30 | 6/6 evidence criteria: valid mismatch, source-linked expectation/rendering, reproducibility, health controls, honest incomplete scope, preservation/review/replan disposition. No assertion of M1c acceptance. |
| Code quality (15) | 13 | Bounded stdlib harness and structured records; -2 for >50-line main. Large result JSON is generated evidence, not an oversized implementation file. |
| Documentation (15) | 15 | Adapted to runnable reproduction (5), recorded outputs/limits (5), and accurate blocked status annotations (5). No new user feature requires an AIL example or changelog. Low aggregate wording nit is nonblocking. |
| Design fidelity (10) | 10 | Preserves NFC/SimpleFold semantics, defaults and original r3 obligations; adds no replacement algorithm and falsely completes no milestone. |
| **Total** | **83/100** | **PASS, threshold 70; no evidence-deliverable hard fail.** |

No full-suite/lint green has been awarded. An actually failing required CI gate still blocks landing regardless of this evidence score. The existing `TestProtectedGitPathspecEquivalence` is not implemented; root/subdirectory, bare-D replacements, ignored/untracked matrix and required production mutants remain incomplete. The weekly-sweep artifact is peripheral; its historical enumeration claim was not independently remeasured and contributes no score.

## Required disposition and limits

Bank the counterexample and keep M1c incomplete. Next work is bounded designer revision preserving the established Unicode matching semantics, review/quorum, then revised decomposition for affected detection milestones. Coverage must retain r3 anchoring/bare-D controls, ignored/untracked paths, index flags/HEAD and time bounds; none is discharged by this report. No implementation algorithm is proposed or approved here, and a human policy change is not needed merely to preserve the existing semantics while redesigning.

The result is specific to recorded macOS/Git 2.54.0 behavior; it does not prove all Git/platform/locale combinations or external-engine behavior. No unfiltered broker tests ran (PF6), no AIL/driver/skill changed, no real HOME state was seeded, and no external messages, push, merge, release or publication occurred.
