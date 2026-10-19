# M1c designer freshness assessment

Outcome: BLOCKED — capability mismatch in the existing ratified equivalence premise. No new design proposed; no production/test/design files changed.

Base: a0c73053f10092fd511ad22fcf6b95d6a3f896aa. Scope: M1c only, `TestProtectedGitPathspecEquivalence`, real-Git equivalence of §4.2 rendering including §12 r3. Read CLAUDE.md, charter preamble/repo profile, coding standards, DESIGN §§1/14, design-doc-creator SKILL.md, ratified design §4.2/§12, M1c/PF6, current exec_protected.go and full exec_protected_test.go, and existing worktree_delta Git test pattern. Existing design remains the authority; no creation/quorum workflow was launched for this assessment.

## Blocking finding

`matchProtected` uses NFC + Unicode SimpleFold orbit minima. Existing shared fixture rows require `tests/**` to match `teſts/x.py` and `*k_test.go` to match `pkg/K_test.go` (exec_protected_test.go:196–197). Git 2.54.0 `icase` does not provide those Unicode equivalences. Both current renderings miss tracked and modified protected files that the mirror positively matches. M1c requires the common bank and Unicode probes, so an honest full equivalence assertion cannot pass with the current production renderer. Dropping these rows or treating the misses as protection would change the ratified claim. Resolving that semantic/capability mismatch requires a design-direction disposition beyond this bounded test-only milestone; this assessment makes none.

## Command/output evidence

Every command ran with `export PATH=/opt/homebrew/bin:$PATH` and cwd `/Users/voightkampff/dev/sunholo-data/.wt-world-iter281-row156-m1c`. `git rev-parse HEAD` returned the base above; `git --version` returned `git version 2.54.0`.

Real-Git probes: Python `tempfile.TemporaryDirectory(dir=Path('.designer-m1c-scratch').resolve())`; one separate initialized repository per Unicode/control row, fake HOME under the fixture, `GIT_CONFIG_NOSYSTEM=1`, `GIT_CONFIG_GLOBAL=/dev/null`, `GIT_OPTIONAL_LOCKS=0`, `LANG=C`, `LC_ALL=C`. Each repository configured `core.ignorecase=false`, committed the named path, then modified its bytes. Commands were `git -c core.fsmonitor=false ls-files [-- PATHSPEC]` and `git -c core.fsmonitor=false status --porcelain=v1 [-- PATHSPEC]`. All subprocesses checked exit success; none of the empty results below is a failed command.

| Fixture | Pathspec | Unfiltered ls-files / status | Selected ls-files / status |
|---|---|---|---|
| ASCII directory positive control | `:(top,glob,icase)tests/**` | `tests/x.py` / ` M tests/x.py` | `tests/x.py` / ` M tests/x.py` |
| Long-s | same | `"te\305\277ts/x.py"` / ` M "te\305\277ts/x.py"` | empty / empty |
| ASCII k positive control | `:(top,glob,icase)**/*k_test.go` | `pkg/k_test.go` / ` M pkg/k_test.go` | `pkg/k_test.go` / ` M pkg/k_test.go` |
| Kelvin sign | same | `"pkg/\342\204\252_test.go"` / ` M "pkg/\342\204\252_test.go"` | empty / empty |

Separate fixture repos are significant on this APFS host: a preliminary mixed directory fixture aliased `tests` and `teſts` on disk. The table above is the corrected isolated experiment, with unfiltered same-repo positive instrumentation controls plus ASCII counterpart controls. NFC/NFD behavior on this filesystem was not used to establish an additional blocker.

Focused mirror/renderer verification:

```sh
TMPDIR="$PWD/.designer-m1c-scratch" GOTMPDIR="$PWD/.designer-m1c-scratch" go test -count=1 ./host/broker/... -run '^TestProtected(MatchFixture|Renderers)$' -timeout 120s -v
```

Output: PASS; `TestProtectedMatchFixture/long-s`, `/kelvin`, `/nfc` all PASS; `TestProtectedRenderers/git` PASS; package `ok .../host/broker 0.348s`. This demonstrates the semantic mismatch, rather than assuming the mirror expects those Unicode matches. It does not claim M1c already passes.

## r3 remains technically valid

In a disposable repository with tracked modified `root_test.go`, cwd `<repo>/sub`, both `status --porcelain=v1 --untracked-files=all --ignored=matching -- ':(top,glob,icase)**/*_test.go'` and `ls-files --` with that pathspec included the root file. With magic changed to `:(glob,icase)`, both returned empty. Same-scope positive control was the identical command with `top` retained. `ls-files` emitted `../root_test.go` while status emitted the repository-root path; a future test must account for this output-relative spelling without confusing it with match selection.

For initially untracked `newtests` as a regular file, and independently as a symlink to tracked `outside.go`, status with `:(top,glob,icase)newtests` plus `:(top,glob,icase)newtests/**` printed `?? newtests`. Removing only bare D yielded empty for each. These are valid independent positive/negative controls with no deleted descendant to mask the omission. ASCII case control also succeeded: `:(top,glob,icase)**/*_test.go` selected modified `pkg/A_TEST.GO`.

## Boundaries and verification limits

`rg -n 'func TestProtectedGitPathspecEquivalence|func TestProtectedRenderers|func TestProtectedMatchFixture|long-s|kelvin|nfc' host/broker/exec_protected_test.go` found the existing match/renderer tests and Unicode rows, and no M1c function; this absent-test finding is paired with existing-function hits from the same search. The production renderer already emits the exact r3 `top,glob,icase` spelling and paired bare-D/descendant pathspecs; the missing work is real-Git proof, not an unimplemented r3 renderer.

No unfiltered broker suite ran (PF6); no global gate is claimed green. No Git writes occurred in the checkout metadata; all Git init/add/commit/config mutations were confined to disposable worktree-local fixture repositories, which TemporaryDirectory removed. Go scratch was worktree-local. No real HOME state was seeded. No AILANG capability claim was made, no `.ail` changed, no driver/skill changed, and no publication/release/message action occurred.

Final disposition: BLOCKED for the existing Unicode equivalence premise; r3 anchoring and bare-D direction verified fresh. Report only; executor must not silently redefine the common fixture expectations to make M1c green.
