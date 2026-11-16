# Sprint plan — w-world-address-bar-verbs-v2, sprint V2a "the seam" (iteration 397, row 243)

> **THE SPRINT — row 243 phase V2, sprint V2a (design §0.4, §2.1, §3, §6 "V2a — the seam", §7, §8 M1–M3):**
> `cmd/ailang-worldd/why.go` moves to a new pure package **`host/why`** (`git`-paired rename) behind a
> **`Getter`** (one GET path → status and body), with **`Run`**, **`ParseTarget`**, **`Head`** and **`Elide`**;
> the CLI's `why` becomes a front (`cmd/ailang-worldd/why_cli.go`) that makes the terminal's two reads
> (`--result`'s file, `-`'s stdin) and `--expect-store`, then calls `why.Run` and writes its `Answer`; a
> **compatibility layer** keeps every other CLI reader of World (`demo`, `demo proof`, `episode report`,
> `log tail`, `provenance`, `outcome`) and all their tests **unedited**; `mcpclient.go`'s `headState` and
> `elide` become one-line delegations to `why.Head`/`why.Elide`. **Behaviour-neutral**: no route, page,
> `/v1`, `.ail` or user-visible change; the bar's `why` stays `terminal`.
> - **Measured: 1,299 changed lines** — `git diff --numstat 18d30a17 HEAD -- cmd host docs` (git's default
>   rename detection, renames on), insertions + deletions, after a commit of each of the three scratch stages
>   (**786 production** — `host/why/why.go` 518 rename-paired, `why_cli.go` 241, `mcpclient.go` 27; **513 test**
>   incl. the 11-line golden; **0 docs**). Cap **1,300**: margin **1**.
> - **The cut is taken (V2a-D2):** with M1+M2 and M3's tests the sprint measures 1,299, so design §2.2
>   (`WhyFlags` & friends in `host/verbs`, the CLI's use of them, VA1) **moves to V2b**, whose `AdmitWhy` is its
>   first caller. That is the design §8 cut ("M3's tests before M3's CLI repointing") taken whole: without the
>   repointing, `WhyFlags` would land as dead code pinned only to literals (Appendix C, with its size).
> - **The proof is behaviour neutrality** (§3.1): every existing `why` test and every caller of a moved name
>   passes **unedited** (the whole `cmd/ailang-worldd` package in M12's acceptance), and a **characterisation
>   pin** (VC1, eleven paths no existing test drove, found by a coverage run at the base) is generated **at the
>   base** and lands first (M0), so the move is proven byte-identical at its own commit.
> - **Eleven premises are false or incomplete and routed minimally** (§1.4), incl. two the design could not see
>   (vet refuses another package's unkeyed struct literal; macOS `patch` leaves a deleted file as 0 bytes
>   without `-E`).
> - **Every line of code, test and golden is a unified diff in Appendix P**, generated from compiled, vetted,
>   gofmt-clean, green scratch stages, expressed as whole-file deletions + creations (no rename headers);
>   extracted **from this file** with §2's `awk`, applied with `patch -p1 -E` to a fresh `git archive 18d30a17`
>   export and compared byte-for-byte with the stage tree (§1.2 C17). The executor applies; it writes nothing.

- **Design**: [w-world-address-bar-verbs-v2.md](w-world-address-bar-verbs-v2.md), **r1** (quorum r1 BLOCKED 3/3 →
  r1 → quorum r2 PROCEED 3/3, no absentee; independent audit PASS 84, 0 blocking; iteration 396; **not
  reopened**). Where this plan and the design differ, this plan was measured at `18d30a17` and wins; each
  difference is a decision (§1.3) or a premise found false (§1.4). The audit's planner inputs (§13.2) 1, 2, 3 and 7
  are resolved in §1.3 (V2a-D3 … V2a-D6); 4–6 are V2b's (Appendix C).
- **Authority**: D-WORLD-107 = A, D-WORLD-109 = A (test stores, temp dirs and test keys only), D-WORLD-89 R1–R5
  untouched (no POST, write or key path). No ledger ask (design Q4).
- **Planner base**: `18d30a17a058e9adad71c42b6078b538374378d7` (= origin/dev, PR #584's merge). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter397-row243-v2a`, branch `sprint/world-iter397-row243-v2a`.
  The design was measured at `200ea5d5`; `git diff --stat 200ea5d5 18d30a17 -- cmd host docs` is empty (§1.2 C2),
  and every design line this plan relies on was re-measured at `18d30a17`.
- **Planner scratch** (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter397/plan-scratch/`
  (`base/` the pristine `git archive` export, a throwaway repo tagged `BASE`; `work/` the tree the stages were cut
  from; `work2/` the stage repo, tags `M0`, `M12`, `M3t`, the tree every mutant cloned; `rt/` the round trip;
  `patches/`; `tools/` (`mkwhy.py`, `mkcli.py`, `mkmcp.py`, `mkprot.py` — the scripted edits of the move —
  `stages.sh`, `mkpatches.sh`, `accept.sh`, `accept_all.sh`, `gates.sh`, `mutants.py`, `mutate.py`, `roundtrip.sh`,
  `assemble.py`, `plan_template.md`); `logs/` (`base_gates/`, `final_gates/`, `accept_<M>.log`, `accept_all.txt`,
  `mutant_ledger.tsv`, `mut/`, `roundtrip.log`, `base_cmd_cover.out`, `m3t_why_cover.out`)).
- **Executor**: as routed by the controller. **One run** (M0, M12, M3t). Each sub-milestone is one `patch -p1 -E`
  (§2) plus its acceptance (§3). **Apply, never retype.**
- **Touched** (§5): new `host/why/{why.go (the move), why_input_test.go (the move), why_protected_test.go,
  run_test.go}`, `cmd/ailang-worldd/{why_cli.go, why_seam_test.go, why_characterise_test.go,
  testdata/why-characterise.golden.jsonl}`; edited `cmd/ailang-worldd/mcpclient.go` (+5 −22) and
  `cmd/ailang-worldd/why_protected_test.go` (+1 −65, the split, V2a-D4); deleted `cmd/ailang-worldd/why.go`,
  `cmd/ailang-worldd/why_input_test.go` (moved). **Outside the design's §7 lists**: `why_seam_test.go` (VW2–VW5
  need cmd's walk rig, §1.4 P9), `why_characterise_test.go` + its golden (VC1, directive retro 4),
  `host/why/why_protected_test.go` as a split rather than a move (V2a-D4).
- **Not touched**: `cli.go` (V2a-D5), `demo_seed.go` (V2a-D3), every other `cmd/ailang-worldd` file and test,
  `host/verbs/*` (moved to V2b with §2.2, V2a-D2), `host/daemon/*` (incl. `surface_verb*.go`), `host/surfacepage/*`,
  `world/*.ail`, `packages/**`, `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, row 201's files, any skill,
  `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`. **No `/v1` route, no new route, no POST, store
  write, process, network listener beyond `httptest` loopback, or `.ail`.**

---

## 0. Binding constraints (each is restated where it governs)

1. **Test stores, temp dirs and test keys only (D-WORLD-109 = A).** Every daemon is `daemon.New` behind `httptest`
   (cmd's `newWalkRig`, `why_test.go:52`) on a `t.TempDir()` store; host/why's tests use in-memory `Getter`s. No
   real store, `~/.ssh`, `/dev/tty`, browser, mint, publish or listener beyond loopback.
2. **Pure core, effects at the boundary (coding standards S1–S3, design §3).** `host/why`'s production file
   imports `context encoding/json errors fmt io net/url sort strconv strings` + `host/broker`,
   `host/coordinator`, `host/hashref` — no `os`, `os/exec`, `net`, `net/http`, `io/fs`, `database/sql`,
   `host/store` (VS2a pins **direct** imports; `go list -deps` shows `host/store` reached transitively through
   `host/broker`/`host/coordinator`, §1.2 C9). Every read goes through the injected `Getter`. The `--result` file
   read, the stdin read and `--expect-store`'s `filepath.Abs` stay in cmd (`why_cli.go`).
3. **Behaviour-neutral for the CLI (design §0.4).** Every existing test passes **unedited**, except the two test
   files that move with the code, as re-cut in V2a-D4 (exactly: `why_input_test.go` moved with 3 changed lines,
   `why_protected_test.go` split). No other pre-existing `_test.go` changes (§1.2 C10).
4. **One mutant per new or changed production branch, anchored to the diff** (§4): the named test and the
   behavioural suite with every static test `-skip`ped; every run appends a ledger line; the table's ids are the
   ledger's ids.
5. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
6. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet`
   green, the `_test.go` compile fence green, its own tests and every earlier one green (§1.2 C15).
7. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `18d30a17` (measured by the planner, 2026-11-16)

| Leg (`tools/gates.sh <tree> <logs> base`) | rc | Time |
|---|---|---|
| `status` | **0** | 0s |
| `verify_ail` | **0** | 48s |
| `gofmt` | **0** | 0s |
| `build` | **0** | 7s |
| `vet` | **0** | 3s |
| `fence` | **0** | 35s |
| `pkgs` | **0** | 237s |
| `race` | **0** | 169s |
| `inv` | **0** | 38s |
| named PASS: 43 | | |
| `verify_go` | **1** | 617s |
| `GATES_DONE 1154s FAILS: verify_go` | | |

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0), `go1.26.6 darwin/arm64`, in
a `git archive 18d30a17` export at `plan-scratch/base` (a throwaway git repo; `diff -r cmd`, `host`, `docs` against
the worktree: identical — `BASE_EQUALS_WORKTREE`; the worktree itself was only read), with §6's script in `base`
mode (the sprint's new tests and `host/why` left out). Load average 8–25 throughout (this planner's stage,
acceptance, coverage and mutant runs alongside). Logs: `plan-scratch/logs/base_gates/`.

- `verify_ail`: `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`.
- `pkgs` rc 0 (`ok host/verbs`, `ok cmd/ailang-worldd 236 s`); `race` rc 0; `inv` rc 0: **43** named `--- PASS`, 0
  SKIP, 0 FAIL, none missing.
- `verify_go` rc 1: **60** packages `ok` (incl. `cmd/ailang-worldd` 384 s); red only in untouched packages, each a
  deadline under load: `host/broker` (`TestExecProtectedCheckTimeout` — the directive's list —
  and `TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce`, 4.96 s), `host/daemon`
  (`TestServeStartsWhenTheProbeSpendsTheStartupClock`, red at V1's base too), `host/verifygate`
  (`TestMissionConfigDefaultFlowIntegration`, the directive's list).
- Whole script: `GATES_DONE 1154s FAILS: verify_go`.

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red untouched
package alone once and does not chase it. **Never start a gate under `nohup`.**

### 1.2 Codebase claims checked (at `18d30a17` unless "work" — the planner's final stage tree `work2` at `M3t`)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base and the tools | `git rev-parse HEAD` (worktree) → `18d30a17a058e9adad71c42b6078b538374378d7`; `git status --porcelain \| wc -l` → `0` (before the two plan files); `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0`; `go version` → `go1.26.6 darwin/arm64`; `/usr/bin/patch --version` → `patch 2.0-12u11-Apple`; `git --version` → `git version 2.54.0` |
| C2 | **No drift since the design's base** (negative, with control) | `git diff --stat 200ea5d5 18d30a17 -- cmd host docs \| tail -1` → *(empty)*; control `… -- design_docs \| tail -1` → `11 files changed, 1204 insertions(+), 124 deletions(-)` |
| C3 | **`why.go`'s three parts** (design W2) | `wc -l cmd/ailang-worldd/why.go` → `1050`; the read surface `:44-260`, the CLI front `:255-422` (`whyHelp :258`, `whyTarget :297`, `parseWhyTarget :304`, `trimOneNewline :342`, `runWhy :349`), the walk `:424-1050` |
| C4 | **Every `why.go` identifier other cmd files use** (audit input 1; `bash -c` loop of `grep -lw` over `cmd/ailang-worldd/*.go` minus `why.go`) | `wireHeader`: demo_seed; `wireEntry`: demo_proof_test, demo_seed, demo_test, logtail, provenance, why_test; `wireObject`: demo_seed, demo_test; `wireWorld`: demo_proof_test, demo_seed, demo_test; `invRecord`: demo_test, provenance; `worldReader`: demo_proof, demo_proof_test, demo_seed, demo_test, episode_report, logtail, outcome_test, provenance; `healthInfo`: cliwalk_e2e_test, demo_test, marketing_chrome_test, why_test; `whyTarget`: **demo_seed** (`:658` `r.locate(whyTarget{kind: "index", index: index}, whyScanDefault)`); `sum(`: **demo_seed `:730`, `:801`, `:829`, `:998`, `:1064`** (+ demo_test); `whyChain`, `entryWire`, `worldWire`, `sumJSON`: demo_seed, demo_test; `whyBudget`: episode_markdown_test, episode_report, provenance; `whyScanDefault`: demo_seed, provenance; `whyScanMax`: provenance; `logPageMax`: logtail; `whyHelp`: keyguard; `runWhy`: demo_seed, main, outcome_test, service_pointers_test, why_test; `contains`: 10 files; `whyLink`: demo_test, why_protected_test; `protectedLine`, `resultProtected`: why_protected_test; `whyInputLine`, `whyKeyFieldMax`, `whyKeyFieldCount`: why_input_test; `errScanExhausted`, `parseWhyTarget`, `trimOneNewline`, `byWorld`, `scanFor`, `failNote`, `printChain`, `worktreeLine`, `resultWorktreeDelta`, `whyProtected`, `whyWorktreePaths`, `frozenRigPort`, `healthCache`: none |
| C5 | **Which `worldReader` methods other files call** (negative for two, with control) | `grep -nE '\.(locate\|getJSON\|effectRecord\|checkedObject\|buildChain\|headIndex\|scanEntries\|entry\|page\|object\|world\|health\|record)\(' cmd/ailang-worldd/*.go \| grep -v '^cmd/ailang-worldd/why'` → callers for all thirteen (e.g. `demo_seed.go:43` getJSON, `:658` locate, `:662` buildChain, `:965` checkedObject; `logtail.go:82` headIndex, `:95` page; `episode_report.go:135` health; `provenance.go` scanEntries); `grep -nE '\.(byWorld\|scanFor)\(' … \| grep -v why.go \| wc -l` → `0`; control `grep -cE '\.(byWorld\|scanFor)\(' cmd/ailang-worldd/why.go` → `12` |
| C6 | **`c.get` is taken** (audit input 3) | `grep -n '^func (c \*client) \(do\|get\)(' cmd/ailang-worldd/cli.go` → `48: func (c *client) do(ctx context.Context, method, path string, body io.Reader) (int, []byte, error) {`, `94: func (c *client) get(path string) (int, string, error) {`; work: `cli.go` unchanged (`git diff --name-only BASE M3t -- cmd/ailang-worldd/cli.go` → empty) |
| C7 | **The moving tests use cmd's rig and `elide`** (audit input 2) | `grep -n 'newWalkRig(t)\|whyRun(t\|mustJSONT' cmd/ailang-worldd/why_protected_test.go` → `52`, `59`, `74` (subtest `e2e-exec`), **`96`, `100`, `107`** (subtest `e2e-tampered`, not in the audit's list), `140`, `156` (`mustJSONT`, defined `why_test.go:97`); `grep -n 'elide(p)' cmd/ailang-worldd/why_input_test.go` → `31`, `58` |
| C8 | **`headState`'s call sites** (audit input 7; design W35/W39) | `grep -n "headState(" cmd/ailang-worldd/*.go \| grep -v _test \| grep -v "func headState"` → 8 lines: `demo_proof.go:204`, `demo_seed.go:98`, `:209`, `:699`, `episode_report.go:161`, `mcpclient.go:305`, `:500`, `why.go:148` (W35's `9` counts the definition line); work: `why.go:148` is `Head(r.Ctx, r.Get)` inside host/why, the seven others unedited |
| C9 | **`host/why` reaches `host/store` only transitively** (audit input 7) | work: `go list -deps ./host/why \| grep 'ailang-world/host/'` → `childenv hashref pkgproj procbound canon store broker archive capsule pins registry transitionreg coordinator why`; direct imports (VS2a) exclude it — `TestWhyImportsNoEffect` PASS, its control finds `host/broker` |
| C10 | **No other pre-existing test file changes** (design §0.4 gate) | work: `git diff --name-status --no-renames BASE M3t -- '*_test.go'` → `A cmd/ailang-worldd/why_characterise_test.go`, `M cmd/ailang-worldd/why_protected_test.go`, `A cmd/ailang-worldd/why_seam_test.go`, `D cmd/ailang-worldd/why_input_test.go`, `A host/why/run_test.go`, `A host/why/why_input_test.go`, `A host/why/why_protected_test.go`; with renames: `{cmd/ailang-worldd => host/why}/why_input_test.go 3 3` |
| C11 | **vet refuses another package's unkeyed literal** (§1.4 P1) | work, with `entryWire = why.EntryWire` aliases: `go vet ./cmd/ailang-worldd/` → `demo_seed.go:804:60: …worldWire struct literal uses unkeyed fields`, `:807:20 …entryWire…`, `demo_test.go:1864:30`, `:1866:28`; with defined types (V2a-D7): `VET_OK`. At the base: `grep -n 'worldWire{\|entryWire{' demo_seed.go demo_test.go` → those four lines (positional) |
| C12 | **The paths no existing test drove** (retro 4) | `go test -coverprofile ./cmd/ailang-worldd/` at the base (`ok … 238.1s coverage: 83.9%`): **59** zero-count blocks in `why.go` (`logs/base_why_uncovered.txt`), among them `runWhy`'s trailing-target arm `:364-369`, `--expect-store`'s health cache `:169`, the `rest:` receipt `:449-460`, `parseWhyTarget`'s `arg == ""` `:316`, malformed `sha256:` `:329`, `--result` + target `:306`, `--result` read error `:310`, the not-JSON result fallback `:478`, `byWorld`'s absent world `:538`, `getJSON`'s cancelled-context `:100`; `mcpclient.go:261` (`headState`'s error arm). §3.1 routes each |
| C13 | **The walk rig is deterministic** (VC1's premise) | VC1 under `WORLD_UPDATE_WHY_GOLDEN=1` at the base, then twice more on fresh rigs without it: `ok` ×2 (every hash in the eleven answers equal; only `t.TempDir()` varies, replaced by `$TMP`) |
| C14 | **The move pairs as a rename** (directive size check) | work2: `git diff -M --numstat BASE M12 -- cmd host docs` → `233 285 {cmd/ailang-worldd => host/why}/why.go`, `3 3 {cmd/ailang-worldd => host/why}/why_input_test.go`; `git config --get diff.renames` → *(unset: git's default, renames on)* |
| C15 | **Every sub-milestone is green when committed** | `tools/accept_all.sh` (an APFS clone of `work2` checked out at each tag; `logs/accept_<M>.log`): `M0 rc=0 ACCEPT_OK M0`; `M3t rc=0 ACCEPT_OK M3t`; `M12 rc=0 ACCEPT_OK M12`; `M3t rc=0 ACCEPT_OK M3t`; `M3t rc=0 ACCEPT_OK M3t`; `M0 rc=0 ACCEPT_OK M0`; `M3t rc=0 ACCEPT_OK M3t`; `M12 rc=0 ACCEPT_OK M12` — **`ACCEPT_OK` for all three** (M3t re-run twice, after VW3 gained its deadline row and then its health-error row; every run `ACCEPT_OK`) |
| C16 | **Every sub-milestone is committed and measured** | `tools/stages.sh` (commits each stage in `work2`, measures `git diff -M --numstat BASE <tag> -- cmd host docs`): `NUMSTAT M0 77`, `M12 1039`, `M3t 1299` (`-- cmd host` the same), `WORK2_EQUALS_WORK` |
| C17 | **The patches reproduce the stage tree** | `tools/roundtrip.sh <this plan>`: a fresh `git archive 18d30a17` export + `patch -p1 -E` (after `--dry-run`) of Appendix P's three patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; no `.orig`/`.rej`; the two deleted files gone), each committed and measured → `CMP_EQUAL M0`; `NUMSTAT M0 77`; `CMP_EQUAL M12`; `NUMSTAT M12 1039`; `CMP_EQUAL M3t`; `NUMSTAT M3t 1299`; `233	285	{cmd/ailang-worldd => host/why}/why.go`; `3	3	{cmd/ailang-worldd => host/why}/why_input_test.go`; `ROUNDTRIP_IDENTICAL_TO_STAGE_TREE` |
| C18 | **Every mutant; the table's ids are the ledger's** | 72 rows (63 `K`, 3 `S`, 2 `E`, 1 `N` non-compiling first form, 1 `Kw` (killed by a behavioural test of the whole package, outside the suite regex), 2 controls); 86 ledger lines; the last line of each row behaves as expected for **72 of 72** — `K` rows killed by the named test **and** by the behavioural suite with every static test (`TestWhyHeadIsOneReader`, `TestWhyConstantsMatchTheCLI`, `TestWhyImportsNoEffect`) `-skip`ped; `S` rows killed only by VW4's constant pin (named, reasons in Appendix C); `E` rows named equivalents (reasons in §4); **0 not as expected**. Table ids = ledger ids (`tools/assemble.py` fills §4 from the ledger's last line per id) |
| C19 | **The gates, on the final tree** | `tools/gates.sh <final-tree> logs/final_gates final` (`final-tree` = `work2` at `M3t`, a git repo with every stage committed, byte-identical to the round-trip tree; load 9–25): `status` 0; `verify_ail` rc 0 (no `.ail` changed; `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`, as the base); `gofmt`, `build`, **`vet`** (§1.4 P1's leg), `fence` rc 0; **`pkgs` rc 0** (`ok host/why`, `ok host/verbs`, `ok cmd/ailang-worldd 466 s` incl. every test of the package unedited; an earlier final run on the previous stage tree had two load reds in untouched host-launcher tests, `TestWorldContextAndStatusline`, `TestHostsLaunchersExecHost` "`signal: killed`", the directive's family — they did not recur); **`race` rc 0**; **`inv` rc 0: 52 named `--- PASS`**, 0 SKIP, 0 FAIL, none missing — incl. `TestWorldBoundaryDependencyAllowlist`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestBoundaryASTWriteGuard`, `TestProductionContextRoots`, `TestProductionGoSurface`, `TestDaemonDependencyAllowlist`, `TestSurfaceVerbIsReadOnly`, `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, every moved and new test (`TestElideBounds`, added after that run's script was fixed, `--- PASS` named separately, `logs/final_gates_elide_named.log`); `verify_go` rc 1: **62** packages `ok` (incl. `host/why`, `host/verbs`, `cmd/ailang-worldd 251 s`), red only `host/broker` (`TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce` — red at the base too —, `TestPublishDenialsPersistAndNeverDispatchWithLivePositiveControl`, `TestWorktreeSnapshotRefusesAConfiguredFilter`) and `host/verifygate` (`TestMissionConfigDefaultFlowIntegration`) — untouched, the directive's load families, **all green re-run alone** (`ok host/broker 1.3s`, `ok host/verifygate 11.4s`; `logs/final_broker_alone.log`, `logs/final_verifygate_alone.log`). **`GATES_DONE 1260s FAILS: verify_go`**, every red an untouched load test green alone; `ALL_DONE.M3t` written |
| C20 | **Every name the patches use exists at HEAD or is new** (396 retro) | at HEAD: `client.do` `cli.go:48`, `newClient` `cli.go:31`, `budgetContext` `cli.go:90`, `maxClientResponseBytes` `cli.go:21`, `exitOK`/`exitUsage` `main.go:74-75`, `exitIntegrity` `mcpclient.go:41`, `parseVerbFlags` `mcpclient.go:414`, `serviceMessage` `service.go:96`, `sameStorePath` `doctor.go:369`, `headState` `mcpclient.go:259`, `elide` `mcpclient.go:651`; the test rig: `newWalkRig` `why_test.go:52`, `invoke` `:109`, `synthInv.input`/`outputText` `:44-48`, `whyRun` `:208`, `writeFileT` `:305`, `tamperProxy` `:313`; `hashref.Parse`, `hashref.SumSHA256`, `broker.*`, `coordinator.*` as `why.go` uses them. **New** (`host/why`): `Getter`, `Reader{Ctx, Get, HealthCache}`, `Head`, `Target{Kind, Index, Text, Result}`, `ParseTarget`, `Options{Scan, JSON}`, `Answer{Stdout, Stderr, Exit}`, `Run`, `Elide`, `ElideAt`, `Sum`, `SumJSON`, `ScanDefault`, `ScanMax`, `LogPageMax`, `MaxResponseBytes`, and the renamed `Header Entry Object World Record Health Chain Link EntryWire WorldWire` with methods `GetJSON Entry Page Object World HeadIndex Health Record ScanEntries Locate EffectRecord ByWorld ScanFor CheckedObject BuildChain`. **New** (cmd): `clientGetter`, `worldReader.why`, the test helpers `recorded`, `recordingGetter`, `replayGetter`, `whyAnswer` |
| C21 | **`TestElideBounds`' expectations are the base's** | the test copied into the base export as `package main` with `Elide(in)` → `elide(in)` (`mcpclient.go:651`): `--- PASS: TestZZElideBoundsAtBase` (`logs/elide_at_base.log`); removed after (`git status --porcelain \| wc -l` → `0`) |

### 1.3 Planner decisions where the design leaves a choice

- **V2a-D1 — Three sub-milestones; the move is one atomic unit.** **M0** (VC1, the characterisation pin, at the
  base), **M12** (design M1 + M2: the move, `why_cli.go`, the delegations, the two test files), **M3t** (design
  M3's tests: VW2–VW5, VS2a, `ParseTarget`/`Head`/`GetJSON` classes). Design M1 cannot be green alone: deleting
  `cmd/ailang-worldd/why.go` without `why_cli.go`'s compatibility layer breaks eight cmd files, and
  `why_input_test.go` must move in the same commit because `whyInputLine` does. M12 is therefore 962 own lines —
  over the ~250–400 guide, which the directive allows for "the rename-paired move" — and pairs as a rename at its
  own commit (C14). M0 first means the base's bytes are pinned **before** anything moves, and stay green through
  M12 and M3t (C15).
- **V2a-D2 — The cut (design §8, directive size rule).** M0+M12+M3t measure **1,299 ≤ 1,300**. Design §2.2
  (`WhyFlags`, `WhyArgs`, `WhyLeadingTarget`, `WhyTrailingTarget`, `WhyScanDefault/Max`, `WhyScanOK`, the CLI's
  use, VA1) does not fit and **moves to V2b** whole. The design's cut ("M3's tests before M3's CLI repointing")
  read literally would land `WhyFlags` with no caller (dead code, VA1 pinning it to typed literals); taking §2.2
  whole is the smaller, fully proven sprint the directive prefers, and V2b's `AdmitWhy` + the CLI repoint give it
  its two callers at once (Appendix C). The bar's `why` is `terminal` either way (design §0.4).
- **V2a-D3 — `sum` and `whyTarget` (audit input 1): both in the compatibility layer; `demo_seed.go` unedited.**
  `func sum(b []byte) string { return why.Sum(b) }` (one hash implementation). `whyTarget` stays a cmd struct with
  its four unexported fields, because `demo_seed.go:658` builds `whyTarget{kind: "index", index: index}`; the
  `locate` delegation converts it to `why.Target{Kind, Index, Text, Result}`. `parseWhyTarget` returns
  `why.Target` directly (its only caller is `runWhy`, C4). Reason: "no other CLI file is edited" stays true and is
  checked (C10, §5).
- **V2a-D4 — The moving tests, re-cut (audit input 2).** (a) `why_input_test.go` **moves** (`{cmd => host/why}`, 3
  3): line 1 `package main` → `package why`; lines 31 and 58 `elide(p)` → `Elide(p)` — the same function under its
  moved, exported name (design §2.1 moves `elide` to `why.Elide`); no assertion, input or expected value changes.
  (b) `why_protected_test.go` **splits**, because `TestWhyProtectedLine` mixes 18 pure table subtests with two
  subtests on cmd's rig (`e2e-exec` `:51-91`, `e2e-tampered` `:92-121`; `newWalkRig`, `whyRun`,
  `assertChainOK`, `tamperProxy` are cmd test helpers): the pure part — `protectedFixture`, the 18-case table
  loop, `assertProtectedBlock`, `TestWhyProtectedJSON` — goes to the new `host/why/why_protected_test.go` with
  `whyLink` → `Link` on its 4 lines and a 7-line copy of `mustJSONT` (`why_test.go:97`, cmd's helper); cmd's
  file keeps `protectedFixture`, `assertProtectedBlock` and a `TestWhyProtectedLine` holding **only the two e2e
  subtests, byte for byte** (its diff is deletions plus one comment line; the unused `bytes` import goes). Test ids
  stay: `cmd … TestWhyProtectedLine/e2e-exec`, `/e2e-tampered`; `host/why … TestWhyProtectedLine/<18 cases>`,
  `TestWhyProtectedJSON`. Each change is a move: the same assertions on the same inputs, run where the code
  under test now lives. VW1 is these two files passing (M12's acceptance names every test in them).
- **V2a-D5 — The getter (audit input 3): a closure, `clientGetter(c *client) why.Getter`, in `why_cli.go`;
  `cli.go` unedited.** `c.get` (`cli.go:94`) is the auth-less string GET and keeps its shape. `headState`'s body is
  `return why.Head(ctx, clientGetter(c))`.
- **V2a-D6 — The call-site census and the import pin (audit input 7).** `headState` has 8 non-test call lines
  (C8): the 7 outside the walk keep calling cmd's one-line `headState`; `why.go:148` becomes `Head(r.Ctx, r.Get)`
  inside host/why — one `/v1/head` reader (VW5). VS2a pins **direct** imports only; `host/store` is reached
  transitively through `host/broker`/`host/coordinator` (C9) — the walk itself performs no store call (every read
  is a `Getter` call). Size is from the real patch (C14, C16), not the W27 rehearsal (§1.4 P7).
- **V2a-D7 — `entryWire`/`worldWire` are defined types in cmd, not aliases** (§1.4 P1): `type entryWire
  why.EntryWire`, `type worldWire why.WorldWire`. Same fields, tags and encoding; vet accepts the four positional
  literals in `demo_seed.go`/`demo_test.go` as literals of a local type. The other eight names are aliases.
- **V2a-D8 — Thirteen delegations, not fifteen** (§1.4 P3): `byWorld` and `scanFor` have no caller outside the
  walk (C5), so cmd does not carry them. Every delegation that exists has a caller and a mutant (§4 `C-D-*`).
- **V2a-D9 — `elideAt` leaves cmd** (§1.4 P4): its only cmd use was `elide`'s body; `why.ElideAt` is the constant,
  pinned by VW4.
- **V2a-D10 — VC1, the characterisation pin.** For the eleven CLI paths C12 found undriven, `TestWhyCharacterisedPaths`
  runs the real CLI (`whyRun`) on `newWalkRig` and compares the eleven `{stdout, stderr, exit}` with
  `testdata/why-characterise.golden.jsonl`, which the same test wrote **at the base** under
  `WORLD_UPDATE_WHY_GOLDEN=1` (the planner copied the test into the base export, generated, ran twice more fresh,
  C13). The test compiles against the base's names only, so it is M0 and is green before the move.
- **V2a-D11 — `Run`'s streams.** `Run` writes into two `strings.Builder`s and a deferred assignment fills
  `Answer.Stdout`/`Stderr` on every return; its text is `runWhy`'s `:397-421` with `*scan`/`*asJSON` → `o.Scan`/
  `o.JSON` and `exitUsage`/`exitIntegrity`/`exitOK` as host/why constants (1/3/0; unexported — VW2 asserts the
  exits through the CLI's own constants). `runWhy` then `fmt.Fprint`s both and returns `Exit`. Before, the chain
  streamed to stdout; now it is written whole — the same bytes (VW2, VC1, `why_test.go`).
- **V2a-D12 — Statuses without `net/http`.** host/why's `statusOK = 200`, `statusNotFound = 404` replace
  `http.StatusOK/StatusNotFound` (design §2.1); mutants `W-STATUS-*`.
- **V2a-D13 — The move is scripted, so the diff stays minimal.** `tools/mkwhy.py` builds `host/why/why.go` from
  the base file by word-boundary renames, the three replaced blocks (imports, constants, the CLI front → `Target`/
  `ParseTarget`/`Options`/`Answer`/`Run`) and three insertions (`Getter`, `Head`, `Elide`), then `gofmt`. No
  comment is rewritten (the moved doc comments keep their lower-case first words) and nothing is reflowed but what
  `gofmt` realigns after a shorter type name (`Entry`'s three fields).

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §2.1 "positional literals in `demo_test.go` keep compiling because field order is unchanged".** They
  compile, but `go vet`'s `composites` check refuses an unkeyed literal of another package's struct, and two of
  the four are in **production** `demo_seed.go` (`:804`, `:807`; `demo_test.go:1864`, `:1866`) — the `vet` gate
  would go red (C11). Routed: V2a-D7. (`go test`'s built-in vet subset omits `composites`, so only `go vet`
  sees it.)
- **P2 — Audit input 2 lists `why_protected_test.go:52,59,74`.** The second e2e subtest (`e2e-tampered`, `:96`,
  `:100`, `:107`) uses the rig too (C7). Routed: both e2e subtests stay in cmd (V2a-D4).
- **P3 — Design §2.1 "each of its 15 methods becomes a one-line delegation".** Two have no caller (C5). Routed:
  V2a-D8.
- **P4 — Design §2.1 "`mcpclient.go`'s become a one-line delegation and a constant alias".** `elideAt` has no cmd
  use besides `elide` (`grep -n elideAt cmd/ailang-worldd/*.go` → `mcpclient.go:595`, `:650`, `:653`, `:656`, and
  `why.go`). Routed: V2a-D9.
- **P5 — Design §2.1 "type aliases … `entryWire = why.EntryWire`, `worldWire = why.WorldWire`".** See P1; routed
  V2a-D7.
- **P6 — Directive "a whole-file DELETION diff … so `patch -p1` alone reproduces the tree".** macOS `/usr/bin/patch`
  (`2.0-12u11-Apple`) applies a deletion hunk but leaves a **0-byte file** (`ls -la` → `0 … why.go`, which `go
  build` then rejects); `-E` (`--remove-empty-files`) removes it. Routed: §2 applies with `patch -p1 -E`. No
  created file is empty, so `-E` removes exactly the two deleted files. Git rename headers are not used (the
  patches carry `--- a/…`/`+++ /dev/null` and `--- /dev/null`/`+++ b/…` only).
- **P7 — Design §8 "M1 ≈ 410 (rename-paired)" / W27 "56 209".** Measured: `233 285` for `why.go` (C14): the
  rehearsal renamed five types; the real patch exports ten types, fifteen methods and four helpers, which touches
  ≈ 117 kept lines on both sides. With `why_cli.go` (241, of which ≈ 160 are the moved front, new to git) M12
  measures 962, and V2a 1,299 against the design's ≈ 970 — routed by the cut (V2a-D2). Pairing holds (> 50 %).
- **P8 — Design §7 "New: `host/why/run_test.go`" holding VW2/VW3.** Their replay bodies must be recorded from the
  walk rig's real routes (`newWalkRig`, `tamperProxy` — cmd test helpers). Routed: VW2–VW5 live in
  `cmd/ailang-worldd/why_seam_test.go`; `host/why/run_test.go` holds the pure tests (VS2a, `ParseTarget`, `Head`/
  `GetJSON` classes).
- **P9 — Design §6 VW3 "one [getter] returning `err` → the CLI's `ailang-worldd why: <err>` line".** True; and a
  getter returning `context.DeadlineExceeded` prints **"not found in the last N entries"** (P4 of the design, kept
  by V2a): VW3 pins that row too, so `Run`'s deadline arm (no CLI test reaches the 60 s budget) has a behavioural
  kill (`W-RUN-DEADLINE`).
- **P10 — Design W35 "`headState` has 9 non-test call lines".** 8 calls + the definition (C8); the audit's "7 + `why.go:148`" is right.
- **P11 — Design §2.1 "`Run` contains `runWhy`'s lines from `r.locate` to the exit choice (`why.go:394-422`)".**
  At `18d30a17` they are `:397-421` (`:394-396` close `--expect-store`); no consequence.

### 1.5 The size

Lines are insertions + deletions: **`git diff --numstat 18d30a17 HEAD -- cmd host docs`**, with git's default rename
detection (renames on; `diff.renames` unset). **Every row is measured** from the scratch stages, each committed in
the throwaway repo `work2`, and re-measured on the round-trip tree (C16, C17).

| Sub-milestone | Production | Test | Docs | Own | Cumulative vs `18d30a17` (after its commit) | `-- cmd host` | Run |
|---|---|---|---|---|---|---|---|
| M0 — VC1 characterisation pin + its golden (at the base) | 0 | 77 | 0 | 77 | **77** | 77 | 1 |
| M12 — the move (`host/why`), `why_cli.go`, `mcpclient.go`; the two test files | 786 | 176 | 0 | 962 | **1,039** | 1,039 | 1 |
| M3t — VW2–VW5 (`why_seam_test.go`), VS2a + pure tests (`run_test.go`) | 0 | 260 | 0 | 260 | **1,299** | 1,299 | 1 |
| **The sprint** | **786** | **513** | **0** | | **1,299 ≤ 1,300** (margin 1) | 1,299 | |

Per file (+/−): `{cmd/ailang-worldd => host/why}/why.go` +233 −285; `cmd/ailang-worldd/why_cli.go` +241;
`mcpclient.go` +5 −22; `{cmd/ailang-worldd => host/why}/why_input_test.go` +3 −3; `host/why/why_protected_test.go`
+104; `cmd/ailang-worldd/why_protected_test.go` +1 −65; `why_characterise_test.go` +66;
`testdata/why-characterise.golden.jsonl` +11; `why_seam_test.go` +160; `host/why/run_test.go` +100. **If a later
re-measure exceeded 1,300**, the next cut is M3t's `TestHeadAndGetJSONClasses` (37 lines; `Head`'s classes are
also driven by VW2's `head` row and `why_test.go`'s `head` target) — never VC1, VW2 or VW3.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four
backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`
and an `@@ -0,0 …` hunk; a **deleted** file's patch has `+++ /dev/null` and an `@@ … +0,0 @@` hunk (M12 deletes
`cmd/ailang-worldd/why.go` and `why_input_test.go` this way and creates their `host/why/` copies whole; git pairs
the rename when the controller commits). Apply it **only** this way, in one `bash` shell (the controller's shell is
zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M0   # <- the sub-milestone, exactly as in the Appendix P heading: M0 M12 M3t
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter397-row243-v2a
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v2a-sprint-plan.md"
T="$HOME/.ailang/state/world-iter397/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
for f in cmd/ailang-worldd/why.go cmd/ailang-worldd/why_input_test.go; do
  if [ "$M" != M0 ] && [ -e "$f" ]; then echo "STOP: $f should be deleted"; exit 1; fi
done
git status --short --untracked-files=all
EOF2
```

- **`-E` is required** (§1.4 P6): without it the two deleted files remain as empty files and `go build` fails.
- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan and run it, in one bash shell,
  `tee`'d to a log (never only through `tail`). **Never `nohup`, never background a test leg; pass a 600000 ms
  timeout to every long Bash call.**
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every sub-milestone's acceptance is the arm of this script named for it, preceded by the same
**[base]** lines (gofmt, `go build ./...`, `go vet` of the touched packages, the `_test.go` compile fence). Each `v`
line must show `--- PASS` for every test it names and **no `--- FAIL`, no `--- SKIP`**; each `p` line must end `ok`.
The script ends `ACCEPT_OK <M>` or prints `ACCEPT_RED <M>: <step>` and exits 1. The planner's run of it on every
stage is §1.2 C15. Extract and run it (`M` and `W` as in §2):

```bash
bash <<'EOF2'
set -uo pipefail
M=M0
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter397-row243-v2a
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v2a-sprint-plan.md"
L="$HOME/.ailang/state/world-iter397/exec-logs"; mkdir -p "$L"
awk '$0 == "````bash accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$L/accept.sh"
/bin/bash "$L/accept.sh" "$W" "$M" 2>&1 | tee "$L/accept_$M.log"
EOF2
```

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M>: exactly plan §3's acceptance for sub-milestone M (one bash shell; rc 0 = ACCEPT_OK).
T=$1; M=$2
cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }   # -v: every named test PASS, none SKIP/FAIL
p() { go test -count=1 "$@"; }
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M: $*"; exit 1; }; }
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
VET="./host/verbs/ ./cmd/ailang-worldd/"; [ "$M" = M0 ] || VET="./host/why/ $VET"
step go vet $VET
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M: FENCE_RED"; exit 1; fi
WHYT='TestWhyCharacterisedPaths|TestWhyTargetFormsResolveToOneEntry|TestWhyResolvesMCPInvocationID|TestWhyBrokenLinkIsIntegrityRefusal|TestWhyScanBound|TestWhyShowsTheWorktreeDelta|TestWhyShowsTheFailureLine|TestLogTailFollowPrintsEachCommitOnce|TestProvenanceTrailer|TestWhyProtectedLine'
case $M in
M0)  step v ./cmd/ailang-worldd/ -run "^($WHYT|TestWhyProtectedJSON|TestWhyInputLine.*)\$"
     step p ./cmd/ailang-worldd/ -run '^(TestWhy|TestDemo|TestLogTail|TestProvenance)' ;;
M12) step v ./host/why/ -run '^(TestWhyInputLineShowsThePathOfALongEdit|TestWhyInputLineLeavesAShortInputUnchanged|TestWhyInputLineKeyFieldBounds|TestWhyInputLineFallsBackToElide|TestWhyProtectedLine|TestWhyProtectedJSON)$'
     step v ./cmd/ailang-worldd/ -run "^($WHYT|TestEveryVerbHasHelp|TestVerbTableMatchesTheCLISwitch)\$"
     step v ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestEveryCommandDirectoryIsAProtectedGroup|TestBoundaryASTWriteGuard|TestBareNetHTTPExemptionIsPerGroup)$'
     step v ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface)$'
     step v ./host/daemon/ -run '^(TestSurfaceVerbIsReadOnly|TestDaemonDependencyAllowlist)$'
     step p ./host/why/
     step p -timeout 30m ./cmd/ailang-worldd/ ;;
M3t) step v ./cmd/ailang-worldd/ -run "^($WHYT|TestWhyRunReplaysTheCLI|TestWhyRunGetterClasses|TestWhyConstantsMatchTheCLI|TestWhyHeadIsOneReader)\$"
     step v ./host/why/ -run '^(TestWhyImportsNoEffect|TestParseTargetForms|TestHeadAndGetJSONClasses|TestElideBounds|TestWhyProtectedLine|TestWhyProtectedJSON)$'
     step p ./host/why/
     step p ./cmd/ailang-worldd/ -run '^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP)' ;;
*) echo "unknown $M"; exit 2 ;;
esac
echo "ACCEPT_OK $M"
````

### M0 — the characterisation pin (VC1) · 77 test · Run 1

**Patch** `M0`: new `cmd/ailang-worldd/why_characterise_test.go` — `TestWhyCharacterisedPaths`: one `newWalkRig`
with one invocation; the CLI (`whyRun`) for eleven argument vectors — no target; `sha256:zz`; an absent
`sha256:000…0`; `--json 1` (the trailing-target arm); `1 2`; `1 --result <f>`; `--result <absent>`;
`rest:nobody` (the daemon's real receipt state: `not-started`); `1 --expect-store <the rig's db_path>` (the health
cache); `- --scan 3` with stdin `x` (the not-JSON result fallback); `9` — each `{Name, Stdout, Stderr, Code}`
(temp dir → `$TMP`) as one JSON line, compared byte for byte with new
`cmd/ailang-worldd/testdata/why-characterise.golden.jsonl` (11 lines, **generated at the base**, V2a-D10).
**Expected cumulative numstat after the controller's commit of M0: 77.**

### M12 — the move and the seam (design M1 + M2) · 786 prod + 176 test · Run 1

**Patch** `M12` (whole-file deletion of `cmd/ailang-worldd/why.go` and `why_input_test.go`, whole-file creation of
the rest): new **`host/why/why.go`** — the base's `why.go` with: a package doc; `package why`; imports without
`flag net/http os time`; constants `ScanDefault`, `ScanMax`, `LogPageMax`, `MaxResponseBytes` (= 8,388,608, the
CLI client's cap; V2b's getter applies it), `statusOK`/`statusNotFound`, `exitOK`/`exitUsage`/`exitIntegrity`;
`Getter`; `Reader{Ctx, Get, HealthCache}` (was `worldReader{ctx, c, healthCache}`) whose fifteen methods are
exported with unchanged bodies except `r.c.do(r.ctx, GET, p, nil)` → `r.Get(r.Ctx, p)` and `headState(r.ctx,
r.c)` → `Head(r.Ctx, r.Get)`; `Head` (`headState`'s body over a `Getter`, same texts); the CLI front replaced by
`Target`, `ParseTarget` (`parseWhyTarget` minus the `--result` and `-` arms, same texts), `Options`, `Answer`,
`Run` (V2a-D11); `sumJSON`/`sum` → `SumJSON`/`Sum`; `elide` → `Elide` at its four calls and `ElideAt` in
`whyInputLine`; `Elide`/`ElideAt` appended (moved from `mcpclient.go`). New **`cmd/ailang-worldd/why_cli.go`** —
`whyHelp` verbatim; `parseWhyTarget` (the two reads, then `why.ParseTarget`); `trimOneNewline`; `runWhy` (the
base's parse, bound and `--expect-store` lines verbatim, then `why.Run(r.why(), t, why.Options{…})` and two
`Fprint`s); the compatibility layer (eight aliases, two defined types V2a-D7, `sumJSON`, `sum`, `contains`,
`worldReader` + `why()`, `clientGetter`, `whyTarget`, thirteen delegations). **`mcpclient.go`**: `headState` →
`return why.Head(ctx, clientGetter(c))`; `const elideAt` removed; `elide` → `return why.Elide(s)`; `+ host/why`
import. **Tests**: `why_input_test.go` moved (V2a-D4a); `why_protected_test.go` split (V2a-D4b).
**Expected cumulative numstat after the controller's commit of M12: 1,039**, and it **must** show
`233 285 {cmd/ailang-worldd => host/why}/why.go` and `3 3 {cmd/ailang-worldd => host/why}/why_input_test.go`
(C14). If the pairing is missing, **stop** (design §8).

### M3t — the seam's tests (design M3's VW2–VW5, VS2a) · 260 test · Run 1

**Patch** `M3t`: new `cmd/ailang-worldd/why_seam_test.go` — `recordingGetter` (the CLI's client, keeping every
`(status, body)`), `replayGetter` (answers from a recording, `unrecorded <path>` otherwise, one answer
replaceable); **VW2** `TestWhyRunReplaysTheCLI`: for `1`, `head --json`, `1` over `tamperProxy` (the input object
forged), `9`, `a2a:nobody:x --scan 1`, the CLI's `{stdout, stderr, exit}` equals `why.Run` over the live recording
getter **and** over its replay, with the expected exit (0/0/3/1/1), the right stream empty, and the CLI's own line
(`all 7 link(s) verified`, `"matchedBy": "head"`, `BROKEN — integrity refusal (exit 3)`, `ailang-worldd why: log
entry 9 not found`, `… not found in the last 1 entries (invocationId)`); **VW3** `TestWhyRunGetterClasses` over a
recording: the input object answering 404 → exit 3 with `✗ input … — ✗ does not resolve`; `/v1/log/1` answering
an error → `ailang-worldd why: boom`; 500 → `… GET /v1/log/1 returned HTTP 500: down`; 404 → `log entry 1 not
found`; `/v1/health` answering an error → `ailang-worldd why: health down` (`Run`'s `BuildChain` error arm, which no test drove at the base, C12 `:407`); `context.DeadlineExceeded` → `not found in the last 500 entries (entry index)` (§1.4 P9); **VW4**
`TestWhyConstantsMatchTheCLI` (literals: 8,388,608 = `why.MaxResponseBytes` = `maxClientResponseBytes`, 500, 5,000,
500, 160); **VW5** `TestWhyHeadIsOneReader` (AST over `mcpclient.go`: `headState` and `elide` are each one `return
why.Head(…)`/`why.Elide(…)`). New `host/why/run_test.go` — **VS2a** `TestWhyImportsNoEffect` (production files'
direct imports exclude `os os/exec net net/http io/fs database/sql host/store`; control: `host/broker` found);
`TestParseTargetForms` (seven forms; four error texts, the base's); `TestHeadAndGetJSONClasses` (`Head` 200 trimmed
/ 404 absent / 500 error text / transport error; `GetJSON` malformed body; a done context reads nothing);
`TestElideBounds` (`Elide` at 160 and 161 bytes, a newline-escaped value, a multi-byte cut — the expected values
are the base's `elide`'s: the same test run against `mcpclient.go`'s `elide` at the base PASSES, C21).
**Expected cumulative numstat after the controller's commit of M3t: 1,299.**

### 3.1 Behaviour neutrality: every path through the moved code, from its entry point (directive retro 4)

Entry points: **CLI** = `run([…, "why", …])` / `runWhy` (`whyRun`, in-process against `newWalkRig`'s loopback
daemon); **CALLER** = another CLI verb reading World through `worldReader`/`headState`/`elide`; **RUN** = `why.Run`
over a `Getter`; **UNIT** = a direct call. "Pinned" = byte-for-byte equal to the base's output.

| Path (moved or changed code) | Input → output | Driven, unedited (existing) | New in V2a |
|---|---|---|---|
| `runWhy` leading target, flags, `--scan` bound, `parseWhyTarget` (cmd, verbatim) | `3`, `head`, `-`, `--result f`, `--scan 0/5001`, `bogus:…` | `why_test.go` `TestWhyTargetFormsResolveToOneEntry`, `TestWhyScanBound`, `TestWhyResolvesMCPInvocationID` | VC1: no target, `--json 1`, `1 2`, `1 --result f`, `--result <absent>`, `sha256:zz` — **pinned** |
| `--expect-store` (cmd, verbatim) → `HealthCache` → `Reader.Health` | mismatch → service message; match → chain without a second health read | `service_pointers_test.go:458-474` (mismatch) | VC1 `--expect-store this store` — **pinned** (`C-HEALTHCACHE` is equivalent: a second read of the same daemon's health renders the same `pin`, §4) |
| `Run` → `Locate` → `Entry`/`HeadIndex`/`Head`/`Object`/`ByWorld`/`ScanFor`/`ScanEntries`/`Record` | every target kind | `why_test.go` (all eight tests), `why_protected_test.go` e2e ×2 | VW2 (CLI = RUN live = RUN replayed, five forms), VC1 `absent sha256` (`ByWorld`'s absent world), `9` |
| `Locate` `rest:` receipt | not resolved → `receipt rest:x is <state>: nothing committed` | none at the base (C12 `:449-460`) | VC1 `rest: unknown receipt` — **pinned**. The *resolved* arm (`ByWorld(rc.ResultRef)`) needs a REST receipt the rig cannot write — **residual RS-REST** (rename-only lines; V2b's VE1 runs a `rest:` receipt through the route) |
| `Locate` result fallback (`world.plan`) | not-JSON bytes → exhausted | `why_test.go:365` (JSON with `world.plan`) | VC1 `- not a result` — **pinned** |
| `Run`'s exits and streams | 0 stdout / 3 stdout / 1 stderr | `why_test.go` (exit 0, 3, 1) | VW2 (all three, stream emptiness), VW3 |
| `Run`'s deadline arm (`DeadlineExceeded` → "not found") | a walk past the 60 s budget | none (no test spends 60 s) | VW3 deadline row |
| `Run`'s `BuildChain` error arm (`ailang-worldd why: <err>`, exit 1) | the health or record read fails after `Locate` | none (C12 `:407`) | VW3 health-error row (added after the M3t coverage run below found it uncovered) |
| `GetJSON` classes: 200 / 404 / other / error / malformed / done ctx | | `why_test.go` (200, 404) | VW3 (404, error, 500), `TestHeadAndGetJSONClasses` (malformed, done ctx) |
| `BuildChain`, `CheckedObject`, `EffectRecord`, `printChain`, `whyInputLine`, `protectedLine`, `worktreeLine` (bodies unchanged but renames) | ✓/✗ chains, `--json` | `why_test.go`, `why_protected_test.go` (both halves), `why_input_test.go`, `demo_test.go` | VW2 tamper row, VW3 404 row |
| CALLER `demo`/`demo proof`/`demo seed` through 13 delegations, `sum`, `sumJSON`, `whyTarget`, `entryWire`/`worldWire` | | `demo_test.go`, `demo_proof_test.go` (unedited; in the M12 whole-package run) | `C-D-*`, `C-SUM*` mutants |
| CALLER `log tail` (`headIndex`, `page`, `logPageMax`, `effectRecord`), `provenance` (`scanEntries`, `whyScan*`, `whyBudget`), `episode report` (`health`, `whyBudget`), `outcome` (`object`, `entry`) | | `TestLogTailFollowPrintsEachCommitOnce`, `TestProvenanceTrailer`, `TestEpisodeReport*`, `outcome_test.go` | `C-D-*` mutants |
| CALLER `headState` (7 call sites), `elide` (5 call sites) | | `demo_*`, `episode_report*`, `mcpclient_test.go`, call tests | VW5 (AST), `M-HEADSTATE`, `M-ELIDE` |
| `Elide`'s bound, count and escaping (moved from `mcpclient.go`) | ≤ 160 bytes unchanged; longer → 160 bytes + `… (N bytes)` | none pinned the bound or the count at the base (R1's `W-ELIDE-AT`/`-SHORT`/`-COUNT` survived every suite) | `TestElideBounds`, its expectations checked against the base's `elide` (C21) |
| `headState`'s transport-error arm (`mcpclient.go:261`, uncovered at the base) | | none | `TestHeadAndGetJSONClasses` (`Head`'s error row: one implementation now) |

**Measured, not assumed:** a coverage run of the final stage (`go test -coverpkg=./host/why/ ./cmd/ailang-worldd/ ./host/why/`, `logs/m3t_why_cover.out`, `cmd 81.2 %` of host/why's statements) left **42** zero-count blocks in `host/why/why.go` (`logs/m3t_why_uncovered.txt`), every one an error or edge arm inside a moved body whose only change is a rename — except `Run`'s `BuildChain` error arm, which VW3's health row then covered. The base had the same arms uncovered (C12).

**Residuals** (no test drives them; each is a rename-only line of moved code): **RS-REST** above; `HeadIndex`'s
"selected head has no world row" and `ScanFor`'s "no head" (C12 `:156`, `:557`: a store whose head names no world
cannot be built through `POST /v1/commit`); `ScanEntries`' forward page cap (`:217`: `provenance` scans forward over
≤ 500 entries in its test). Each is listed in Appendix C.

---

## 4. Mutation table

**Procedure** (`tools/mutate.py <run> <tree> <workers> [ids…]`, rows in `tools/mutants.py`): for each row, an APFS
clone of the stage repo `work2` at `M3t` (`cp -cR`; byte-identical to the round-trip tree, C17), one exact-string
replacement whose anchor occurs exactly once in the file (else `ANCHOR<n>`, not a run; every anchor was checked
unique before the run), then `go test -count=1` of the **named test** and of the **behavioural suite**: `host/why`
whole (for `host/why` rows) and `cmd/ailang-worldd` `-run '^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)'`,
both with **every static test `-skip`ped** — `TestWhyHeadIsOneReader` (AST), `TestWhyConstantsMatchTheCLI`
(literal constants), `TestWhyImportsNoEffect` (imports) — and `TestDemoSeedPreflightRefusals` (73 s, the directive's
load-sensitive list), so **no kill in the suite column rests on a static test** (directive retro 3). Go rows ran
with `GOFLAGS=-vet=off`, four workers, at load 15–25. A build failure would be `NOCOMPILE` and not a mutant (none
occurred: every row compiled). Every run appends one line per row to `plan-scratch/logs/mutant_ledger.tsv` (time,
run, tree, id, file, from, to, named regex, named result, named failures, suite result, suite failures, seconds,
expected); per-row output in `logs/mut/<run>-<id>.log`. **The ids below are the ledger's ids**; the table shows each
row's last line. Expected: `K` killed by the named test and by the suite; `S` killed only by a static pin (the
named test), the suite survives — a constant with no behaviour the suites reach (reason per row in Appendix C);
`E` named equivalent (reason below); `CTL` unmutated control, both green; `N` a first form that did not compile (not a mutant; rewritten as `…b`); `Kw` killed by the named **whole-package** run through a behavioural test outside the suite regex (the failing test is named in the row). **Runs.** **R1**: every row on the stage tree before its last three test additions. **R2**: the controls and every row whose R1 result was not as expected or whose anchor/named test changed, on the **final** stage tree (byte-identical to the round-trip tree): `W-RUN-BUILD-ERR` survived R1 (no test drove `Run`'s `BuildChain` error arm), so VW3 gained its `/v1/health` error row; `W-ELIDE-AT`, `-NL`, `-SHORT`, `-COUNT` survived R1's named `TestWhyInputLine*` (whose expectation is itself `Elide`), so `TestElideBounds` was added with the base's values (C21) and named; `W-HEAD-TRIM` and `W-RUN-SCAN` were killed in R1 only by a suite test, which R2 names; `W-EXIT-*` were re-anchored after `gofmt` realigned their constants; `C-RUN-JSON` did not compile (`asJSON` unused — not a mutant, `N`) and is rewritten as `C-RUN-JSONb`; `C-CONTAINS` survived R1's suite (its callers' tests are outside the suite regex) and R2 runs the **whole** `cmd/ailang-worldd` package as its named test. Every R1 line stays in the ledger. One row was **re-classified after run
R1**: `W-SCANMAX` (`ScanMax` 5000 → 4999) was expected `K` via `TestWhyScanBound` and **survived** both the named
test and the suite — no test at the base or here drives `--scan 5000`/`5001`; it is killed only by VW4's literal
(`whyScanMax != 5000`), so it is `S`, and the CLI's upper `--scan` bound is a named residual (Appendix C). The
ledger line keeps its original `K`.

**Equivalents.** `C-HEALTHCACHE` (`--expect-store`'s cached health not passed to `why.Reader`): `BuildChain` then
reads `/v1/health` a second time from the same daemon, and the only member it renders (`pin`, design §2.1) is the
same — only a request count could see it. `C-READER-CTX` (the reader's context replaced by `context.Background()`):
the CLI's 60 s budget no longer bounds the walk's reads, but each read is still bounded by the client's own
deadline (`cli.go`), and no test can wait 60 s; it is the CLI's P4 deadline arm, pinned at the `Run` level by VW3's
deadline row (`W-RUN-DEADLINE`).

| Id | File | Mutation (from → to) | Named test | Named | Suite (static skipped) | Failing tests | Exp |
|---|---|---|---|---|---|---|---|
| `CTL-W` | `why.go` | `-` → `-` | `.` | SURVIVED | SURVIVED |  | CTL |
| `CTL-C` | `why_cli.go` | `-` → `-` | `^(TestWhy)` | SURVIVED | SURVIVED |  | CTL |
| `W-CTX` | `why.go` | `if err := r.Ctx.Err(); err != nil {⏎ return false, err` → `if err := r.Ctx.Err(); false && err != nil {⏎ return fals…` | `^TestHeadAndGetJSONClasses$` | KILLED | KILLED | TestHeadAndGetJSONClasses | K |
| `W-GET-ERR` | `why.go` | `status, body, err := r.Get(r.Ctx, path)⏎ if err != nil {⏎…` → `status, body, err := r.Get(r.Ctx, path)⏎ if err != nil {⏎…` | `^TestWhyRunGetterClasses$` | KILLED | KILLED | TestWhyRunGetterClasses | K |
| `W-STATUS-OK` | `why.go` | `statusOK = 200` → `statusOK = 201` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestHeadAndGetJSONClasses, TestCallHostFailureProbe, TestCallJSONOutIsByteExact, TestCallStrictRefusal (+31) | K |
| `W-STATUS-404` | `why.go` | `statusNotFound = 404` → `statusNotFound = 410` | `^TestWhyRunGetterClasses$` | KILLED | KILLED | TestHeadAndGetJSONClasses, TestDemoProofAdapterGuards, TestDemoSeedAuthenticatedCommits, TestDemoSeedCommitFailures (+5) | K |
| `W-EXIT-OK` | `why.go` | `exitOK = 0` → `exitOK = 2` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls, TestDemoSeedWhyNamedLinks, TestOutcomeRecordRoundTrip (+11) | K |
| `W-EXIT-USAGE` | `why.go` | `exitUsage = 1` → `exitUsage = 2` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestWhyCharacterisedPaths, TestWhyRunGetterClasses, TestWhyRunReplaysTheCLI, TestWhyScanBound | K |
| `W-EXIT-INTEGRITY` | `why.go` | `exitIntegrity = 3` → `exitIntegrity = 1` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestWhyBrokenLinkIsIntegrityRefusal, TestWhyProtectedLine, TestWhyRunGetterClasses, TestWhyRunReplaysTheCLI | K |
| `W-SCANDEFAULT` | `why.go` | `ScanDefault = 500` → `ScanDefault = 501` | `^TestWhyConstantsMatchTheCLI$` | KILLED | KILLED | TestWhyScanBound | K |
| `W-SCANMAX` | `why.go` | `ScanMax = 5000` → `ScanMax = 4999` | `^TestWhyConstantsMatchTheCLI$` | KILLED | SURVIVED | TestWhyConstantsMatchTheCLI | S |
| `W-LOGPAGE` | `why.go` | `LogPageMax = 500` → `LogPageMax = 499` | `^TestWhyConstantsMatchTheCLI$` | KILLED | SURVIVED | TestWhyConstantsMatchTheCLI | S |
| `W-MAXRESP` | `why.go` | `MaxResponseBytes = 8388608` → `MaxResponseBytes = 8388607` | `^TestWhyConstantsMatchTheCLI$` | KILLED | SURVIVED | TestWhyConstantsMatchTheCLI | S |
| `W-HEAD-PATH` | `why.go` | `get(ctx, "/v1/head")` → `get(ctx, "/v1/heads")` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestCallHostFailureProbe, TestDemoProofAdapterGuards, TestDemoProofHTTPProducerRoundTrip, TestDemoProofNoVerifiedSetRefuses (+18) | K |
| `W-HEAD-ERR` | `why.go` | `status, body, err := get(ctx, "/v1/head")⏎ if err != nil …` → `status, body, err := get(ctx, "/v1/head")⏎ if err != nil …` | `^TestHeadAndGetJSONClasses$` | KILLED | KILLED | TestHeadAndGetJSONClasses | K |
| `W-HEAD-404` | `why.go` | `case statusNotFound:⏎ return "", false, nil⏎ default:⏎ re…` → `case statusNotFound:⏎ return "", false, errors.New("absen…` | `^TestHeadAndGetJSONClasses$` | KILLED | KILLED | TestHeadAndGetJSONClasses, TestCallHostFailureProbe, TestDemoProofAdapterGuards, TestDemoSeedAuthenticatedCommits (+1) | K |
| `W-HEAD-TRIM` | `why.go` | `return strings.TrimSpace(string(body)), true, nil` → `return string(body), true, nil` | `^TestHeadAndGetJSONClasses$` | KILLED | KILLED | TestHeadAndGetJSONClasses | K |
| `W-HEAD-5XX` | `why.go` | `return "", false, fmt.Errorf("GET /v1/head returned HTTP …` → `return "", false, nil` | `^TestHeadAndGetJSONClasses$` | KILLED | KILLED | TestHeadAndGetJSONClasses, TestDemoProofAdapterGuards | K |
| `W-HEADINDEX` | `why.go` | `ref, ok, err := Head(r.Ctx, r.Get)` → `ref, ok, err := Head(r.Ctx, func(ctx context.Context, p s…` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestOutcomeRecordRoundTrip, TestProvenanceTrailer, TestWhyBrokenLinkIsIntegrityRefusal, TestWhyProtectedLine (+6) | K |
| `W-RUN-EXHAUST` | `why.go` | `if errors.Is(err, errScanExhausted) \|\| errors.Is(err, c…` → `if errors.Is(err, context.DeadlineExceeded) {` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestWhyCharacterisedPaths, TestWhyRunReplaysTheCLI, TestWhyScanBound | K |
| `W-RUN-DEADLINE` | `why.go` | `if errors.Is(err, errScanExhausted) \|\| errors.Is(err, c…` → `if errors.Is(err, errScanExhausted) {` | `^TestWhyRunGetterClasses$` | KILLED | KILLED | TestWhyRunGetterClasses | K |
| `W-RUN-STREAMS` | `why.go` | `a.Stdout, a.Stderr = stdout.String(), stderr.String()` → `a.Stdout, a.Stderr = stderr.String(), stdout.String()` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls, TestDemoSeedWhyNamedLinks, TestProvenanceTrailer (+9) | K |
| `W-RUN-SCAN` | `why.go` | `e, matchedBy, err := r.Locate(t, o.Scan)` → `e, matchedBy, err := r.Locate(t, ScanDefault)` | `^TestWhyScanBound$` | KILLED | KILLED | TestWhyScanBound | K |
| `W-RUN-RESULT` | `why.go` | `ch, err := r.BuildChain(e, t.Result)` → `ch, err := r.BuildChain(e, nil)` | `^TestWhyBrokenLinkIsIntegrityRefusal$` | KILLED | KILLED | TestWhyBrokenLinkIsIntegrityRefusal, TestWhyTargetFormsResolveToOneEntry | K |
| `W-RUN-BUILD-ERR` | `why.go` | `ch, err := r.BuildChain(e, t.Result)⏎ if err != nil {` → `ch, err := r.BuildChain(e, t.Result)⏎ if err != nil && fa…` | `^TestWhyRunGetterClasses$` | KILLED | KILLED | TestWhyRunGetterClasses | K |
| `W-RUN-MATCHED` | `why.go` | `ch.MatchedBy = matchedBy⏎` → `ch.MatchedBy = ""⏎` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestWhyBrokenLinkIsIntegrityRefusal, TestWhyCharacterisedPaths, TestWhyRunReplaysTheCLI | K |
| `W-RUN-JSON` | `why.go` | `if o.JSON {` → `if !o.JSON {` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls, TestDemoSeedWhyNamedLinks, TestProvenanceTrailer (+10) | K |
| `W-RUN-INTEGRITY` | `why.go` | `if !ch.OK {⏎ return Answer{Exit: exitIntegrity}` → `if false && !ch.OK {⏎ return Answer{Exit: exitIntegrity}` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestWhyBrokenLinkIsIntegrityRefusal, TestWhyProtectedLine, TestWhyRunGetterClasses, TestWhyRunReplaysTheCLI | K |
| `W-PT-EMPTY` | `why.go` | `case arg == "":⏎ return Target{}, errors.New("a target is…` → `case arg == "" && false:⏎ return Target{}, errors.New("a …` | `^TestParseTargetForms$` | KILLED | KILLED | TestParseTargetForms, TestWhyCharacterisedPaths | K |
| `W-PT-HEAD` | `why.go` | `case arg == "head":⏎ return Target{Kind: "head"}, nil` → `case arg == "HEAD":⏎ return Target{Kind: "head"}, nil` | `^TestParseTargetForms$` | KILLED | KILLED | TestParseTargetForms, TestOutcomeRecordRoundTrip, TestWhyBrokenLinkIsIntegrityRefusal, TestWhyRunReplaysTheCLI (+1) | K |
| `W-PT-REST` | `why.go` | `strings.HasPrefix(arg, "mcp:") \|\| strings.HasPrefix(arg…` → `strings.HasPrefix(arg, "mcp:"):` | `^TestParseTargetForms$` | KILLED | KILLED | TestParseTargetForms, TestWhyCharacterisedPaths | K |
| `W-PT-SHA` | `why.go` | `if _, err := hashref.Parse(arg); err != nil {` → `if _, err := hashref.Parse(arg); false && err != nil {` | `^TestParseTargetForms$` | KILLED | KILLED | TestParseTargetForms, TestWhyCharacterisedPaths | K |
| `W-PT-REF-KIND` | `why.go` | `return Target{Kind: "ref", Text: arg}, nil` → `return Target{Kind: "invocation", Text: arg}, nil` | `^TestParseTargetForms$` | KILLED | KILLED | TestParseTargetForms, TestOutcomeRecordRoundTrip, TestWhyCharacterisedPaths, TestWhyTargetFormsResolveToOneEntry | K |
| `W-PT-NEG` | `why.go` | `if err != nil \|\| n < 0 {⏎ return Target{}, fmt.Errorf(` → `if err != nil {⏎ return Target{}, fmt.Errorf(` | `^TestParseTargetForms$` | KILLED | KILLED | TestParseTargetForms | K |
| `W-PT-INDEX` | `why.go` | `return Target{Kind: "index", Index: n}, nil` → `return Target{Kind: "index"}, nil` | `^TestParseTargetForms$` | KILLED | KILLED | TestParseTargetForms, TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls, TestDemoSeedWhyNamedLinks (+6) | K |
| `W-ELIDE-AT` | `why.go` | `const ElideAt = 160` → `const ElideAt = 161` | `^TestElideBounds$` | KILLED | KILLED | TestElideBounds | K |
| `W-ELIDE-NL` | `why.go` | `one := strings.ReplaceAll(s, "\n", `\n`)` → `one := s` | `^TestElideBounds$` | KILLED | KILLED | TestElideBounds, TestMCPWireShapes | K |
| `W-ELIDE-SHORT` | `why.go` | `if len(one) <= ElideAt {⏎ return one⏎ }⏎ return fmt.Sprintf` → `if len(one) < ElideAt {⏎ return one⏎ }⏎ return fmt.Sprintf` | `^TestElideBounds$` | KILLED | KILLED | TestElideBounds | K |
| `W-ELIDE-COUNT` | `why.go` | `strings.ToValidUTF8(one[:ElideAt], ""), len(s))` → `strings.ToValidUTF8(one[:ElideAt], ""), len(one))` | `^TestElideBounds$` | KILLED | KILLED | TestElideBounds | K |
| `C-GETTER` | `why_cli.go` | `return c.do(ctx, http.MethodGet, path, nil)` → `return c.do(ctx, http.MethodPost, path, nil)` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestDemoProofAdapterGuards, TestDemoProofBounds, TestDemoProofHTTPProducerRoundTrip, TestDemoProofNoVerifiedSetRefuses (+23) | K |
| `C-HEALTHCACHE` | `why_cli.go` | `HealthCache: r.healthCache}` → `HealthCache: nil}` | `^TestWhyCharacterisedPaths$` | SURVIVED | SURVIVED |  | E |
| `C-STDIN` | `why_cli.go` | `if arg == "-" {` → `if arg == "--" {` | `^TestWhyTargetFormsResolveToOneEntry$` | KILLED | KILLED | TestWhyBrokenLinkIsIntegrityRefusal, TestWhyCharacterisedPaths, TestWhyTargetFormsResolveToOneEntry | K |
| `C-RESULT-BOTH` | `why_cli.go` | `if arg != "" {⏎ return why.Target{}, errors.New(` → `if false {⏎ return why.Target{}, errors.New(` | `^TestWhyCharacterisedPaths$` | KILLED | KILLED | TestWhyCharacterisedPaths | K |
| `C-RESULT-READERR` | `why_cli.go` | `b, err := os.ReadFile(resultFile)⏎ if err != nil {` → `b, err := os.ReadFile(resultFile)⏎ if err != nil && false {` | `^TestWhyCharacterisedPaths$` | KILLED | KILLED | TestWhyCharacterisedPaths | K |
| `C-RESULT-KIND` | `why_cli.go` | `return why.Target{Kind: "result", Result: trimOneNewline(…` → `return why.Target{Kind: "index", Result: trimOneNewline(b…` | `^TestWhyTargetFormsResolveToOneEntry$` | KILLED | KILLED | TestWhyTargetFormsResolveToOneEntry | K |
| `C-STDIN-KIND` | `why_cli.go` | `return why.Target{Kind: "result", Result: trimOneNewline(…` → `return why.Target{Kind: "result", Result: b}, nil⏎ }⏎ ret…` | `^TestWhyTargetFormsResolveToOneEntry$` | KILLED | KILLED | TestWhyTargetFormsResolveToOneEntry | K |
| `C-TRIM` | `why_cli.go` | `if n := len(b); n > 0 && b[n-1] == '\n' {` → `if n := len(b); false && n > 0 && b[n-1] == '\n' {` | `^TestWhyTargetFormsResolveToOneEntry$` | KILLED | KILLED | TestWhyTargetFormsResolveToOneEntry | K |
| `C-PARSE` | `why_cli.go` | `return why.ParseTarget(arg)⏎}` → `return why.ParseTarget("head")⏎}` | `^TestWhyTargetFormsResolveToOneEntry$` | KILLED | KILLED | TestWhyCharacterisedPaths, TestWhyResolvesMCPInvocationID, TestWhyRunReplaysTheCLI, TestWhyScanBound (+3) | K |
| `C-RUN-STREAMS` | `why_cli.go` | `fmt.Fprint(stdout, ans.Stdout)⏎ fmt.Fprint(stderr, ans.St…` → `fmt.Fprint(stderr, ans.Stdout)⏎ fmt.Fprint(stdout, ans.St…` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls, TestDemoSeedWhyNamedLinks, TestProvenanceTrailer (+9) | K |
| `C-RUN-EXIT` | `why_cli.go` | `return ans.Exit⏎` → `return exitOK⏎` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestWhyBrokenLinkIsIntegrityRefusal, TestWhyCharacterisedPaths, TestWhyProtectedLine, TestWhyRunReplaysTheCLI (+1) | K |
| `C-RUN-SCAN` | `why_cli.go` | `why.Options{Scan: *scan, JSON: *asJSON}` → `why.Options{Scan: whyScanDefault, JSON: *asJSON}` | `^TestWhyCharacterisedPaths$` | KILLED | KILLED | TestWhyCharacterisedPaths, TestWhyRunReplaysTheCLI, TestWhyScanBound | K |
| `C-RUN-JSON` | `why_cli.go` | `why.Options{Scan: *scan, JSON: *asJSON}` → `why.Options{Scan: *scan, JSON: false}` | `^TestWhyRunReplaysTheCLI$` | NOCOMPILE | NOCOMPILE |  | N |
| `C-RUN-JSONb` | `why_cli.go` | `why.Options{Scan: *scan, JSON: *asJSON}` → `why.Options{Scan: *scan, JSON: *asJSON && false}` | `^TestWhyRunReplaysTheCLI$` | KILLED | KILLED | TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls, TestDemoSeedWhyNamedLinks, TestWhyCharacterisedPaths (+4) | K |
| `C-READER-CTX` | `why_cli.go` | `return why.Reader{Ctx: r.ctx, Get:` → `return why.Reader{Ctx: context.Background(), Get:` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | SURVIVED | SURVIVED |  | E |
| `C-SUM` | `why_cli.go` | `return why.Sum(b) }` → `return why.SumJSON(b) }` | `^TestDemo` | KILLED | KILLED | TestDemoSeedM4bBoundaryControls, TestDemoSeedManifestReferences, TestDemoSeedPacketObjectControls | K |
| `C-SUMJSON` | `why_cli.go` | `return why.SumJSON(v) }` → `return why.Sum(nil) }` | `^TestDemo` | KILLED | KILLED | TestDemoSeedM4bBoundaryControls, TestDemoSeedManifestReferences, TestDemoSeedPacketControls | K |
| `C-CONTAINS` | `why_cli.go` | `if v == x {⏎ return true` → `if v == x && false {⏎ return true` | `.` | KILLED | SURVIVED | TestEpisodeStartProvisionsAndPrintsTheLaunchWithTheTokenInTwoFilesOnly | Kw |
| `C-D-GETJSON` | `why_cli.go` | `{ return r.why().GetJSON(path, v) }` → `{ return false, nil }` | `^TestDemo` | KILLED | KILLED | TestDemoProofBounds, TestDemoProofHTTPProducerRoundTrip, TestDemoProofNoVerifiedSetRefuses, TestDemoProofSourceHashRefusal (+8) | K |
| `C-D-ENTRY` | `why_cli.go` | `{ return r.why().Entry(i) }` → `{ return wireEntry{}, false, nil }` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestDemoProofAdapterGuards, TestDemoProofHTTPProducerRoundTrip, TestDemoProofNoVerifiedSetRefuses, TestDemoProofSourceHashRefusal (+9) | K |
| `C-D-PAGE` | `why_cli.go` | `return r.why().Page(from, n)` → `return nil, nil` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestDemoProofHTTPProducerRoundTrip, TestDemoSeedCommitFailures, TestLogTailFollowPrintsEachCommitOnce | K |
| `C-D-OBJECT` | `why_cli.go` | `{ return r.why().Object(ref) }` → `{ return wireObject{}, false, nil }` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestDemoSeedAuthenticatedCommits, TestDemoSeedCommitFailures, TestDemoSeedM4bBoundaryControls, TestDemoSeedManifestReferences (+1) | K |
| `C-D-WORLD` | `why_cli.go` | `{ return r.why().World(ref) }` → `{ return wireWorld{}, false, nil }` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestDemoProofAdapterGuards, TestDemoProofHTTPProducerRoundTrip, TestDemoProofNoVerifiedSetRefuses, TestDemoProofSourceHashRefusal (+7) | K |
| `C-D-HEADINDEX` | `why_cli.go` | `{ return r.why().HeadIndex() }` → `{ return 0, false, nil }` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestProvenanceTrailer | K |
| `C-D-HEALTH` | `why_cli.go` | `{ return r.why().Health() }` → `{ return healthInfo{}, nil }` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestDemoSeedM4bBoundaryControls, TestDemoSeedManifestReferences, TestProvenanceTrailer, TestWhyExpectStore | K |
| `C-D-RECORD` | `why_cli.go` | `return r.why().Record(e)` → `return nil, wireObject{}, nil` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestDemoSeedManifestReferences | K |
| `C-D-SCANENTRIES` | `why_cli.go` | `return r.why().ScanEntries(from, to, reverse, max, visit)` → `return 0, nil` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestProvenanceTrailer | K |
| `C-D-LOCATE` | `why_cli.go` | `why.Target{Kind: t.kind, Index: t.index, Text: t.text, Re…` → `why.Target{Kind: t.kind, Text: t.text, Result: t.result}` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestDemoSeedM4bBoundaryControls, TestDemoSeedManifestReferences, TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls (+1) | K |
| `C-D-EFFECTRECORD` | `why_cli.go` | `return r.why().EffectRecord(ref)` → `return broker.EffectRecord{}, false` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestLogTailFollowPrintsEachCommitOnce | K |
| `C-D-CHECKEDOBJECT` | `why_cli.go` | `return r.why().CheckedObject(ref, semantic)` → `return wireObject{}, true, ""` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestDemoSeedM4bBoundaryControls, TestDemoSeedManifestReferences, TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls | K |
| `C-D-BUILDCHAIN` | `why_cli.go` | `return r.why().BuildChain(e, result)` → `return whyChain{OK: true}, nil` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestDemoSeedM4bBoundaryControls, TestDemoSeedManifestReferences, TestDemoSeedPacketControls, TestDemoSeedPacketObjectControls (+1) | K |
| `M-HEADSTATE` | `mcpclient.go` | `return why.Head(ctx, clientGetter(c))` → `return "", false, nil` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestCallHostFailureProbe, TestDemoProofAdapterGuards, TestDemoProofHTTPProducerRoundTrip, TestDemoProofNoVerifiedSetRefuses (+9) | K |
| `M-ELIDE` | `mcpclient.go` | `func elide(s string) string { return why.Elide(s) }` → `func elide(s string) string { return s }` | `^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestCall|TestMCP|TestEpisodeReport|TestServicePointers|TestCLIWalk|TestCli)` | KILLED | KILLED | TestMCPWireShapes | K |

---

## 5. Files (final `git status --short --untracked-files=all` after M3t, against the base)

```
 M cmd/ailang-worldd/mcpclient.go
 D cmd/ailang-worldd/why.go
 D cmd/ailang-worldd/why_input_test.go
 M cmd/ailang-worldd/why_protected_test.go
?? cmd/ailang-worldd/testdata/why-characterise.golden.jsonl
?? cmd/ailang-worldd/why_characterise_test.go
?? cmd/ailang-worldd/why_cli.go
?? cmd/ailang-worldd/why_seam_test.go
?? host/why/run_test.go
?? host/why/why.go
?? host/why/why_input_test.go
?? host/why/why_protected_test.go
```
Plus the two plan files (`design_docs/planned/w-world-address-bar-verbs-v2a-sprint-plan.md`,
`design_docs/planned/sprint_w-world-address-bar-verbs-v2a.json`), which the controller commits. After M0 only the two
M0 files; after M12 all but `why_seam_test.go` and `run_test.go`. Nothing under `world/`, `packages/`, `scripts/`,
`tools/`, `host/verbs/`, `host/daemon/`, `host/surfacepage/`, `host/store/`, `host/broker/`.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

Run it in **one bash shell, in the foreground, never under `nohup`** (macOS `/bin/bash` 3.2; every leg records its
own rc and every leg runs — no `set -e`; every leg is `tee`'d), extracting the script from this plan:

```bash
bash <<'EOF2'
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter397-row243-v2a
G="$HOME/.ailang/state/world-iter397/gates"; mkdir -p "$G"
awk '$0 == "````bash gates.sh" {f=1; next} f && $0 == "````" {exit} f' "$W/design_docs/planned/w-world-address-bar-verbs-v2a-sprint-plan.md" > "$G/gates.sh"
/bin/bash "$G/gates.sh" "$W" "$G" final "$(git -C "$W" rev-parse HEAD)"
EOF2
```

It writes `gates/legs.txt`, one log per leg, and **`gates/ALL_DONE.<sha>`** when every leg has finished (the judge
waits on it). The `verify_go` leg reads committed blobs, so the script runs **after** the last commit. The legs:
`status`, `verify_ail` (no `.ail` changed: the base's result, §1.1), `gofmt`, `build`, `vet` (`go vet ./...` — the
leg §1.4 P1 would have turned red), `fence` (`go test -run '^$' ./...`), `pkgs` (`host/why`, `host/verbs`,
`cmd/ailang-worldd` whole), `race` (`-race` on `host/why`, `host/verbs` and cmd's walk/caller tests), `inv` (every
named inventory `-v` + every new and moved test: none missing, none SKIP), `verify_go`. The script, verbatim:

````bash gates.sh
#!/bin/bash
# Row 243 V2a gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter397/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.log, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.log"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
PK="./host/verbs/ ./cmd/ailang-worldd/"
[ "$MODE" = base ] || PK="./host/why/ $PK"
leg pkgs "go test -count=1 -timeout 30m $PK"
RACEW=""; [ "$MODE" = base ] || RACEW="go test -race -count=1 ./host/why/ && "
leg race "$RACEW go test -race -count=1 ./host/verbs/ && go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestEpisodeReport|TestMCP|TestCall|TestVerb)'"
NEW_W="TestWhyInputLineShowsThePathOfALongEdit TestWhyInputLineLeavesAShortInputUnchanged TestWhyInputLineKeyFieldBounds TestWhyInputLineFallsBackToElide TestWhyProtectedLine TestWhyProtectedJSON TestWhyImportsNoEffect TestParseTargetForms TestHeadAndGetJSONClasses TestElideBounds"
NEW_V=""   # host/verbs is untouched in V2a (WhyFlags moved to V2b, plan V2a-D2)
NEW_C="TestWhyCharacterisedPaths TestWhyRunReplaysTheCLI TestWhyRunGetterClasses TestWhyHeadIsOneReader TestWhyConstantsMatchTheCLI"
if [ "$MODE" = base ]; then NEW_W=""; NEW_V=""; NEW_C=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes TestSurfaceAppRoutes TestDaemonCredentialWriterSites TestDaemonConstructorCallerInventory TestDaemonCallerInventoryAliasControl"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_VERBS="TestSplit TestQuoteRoundTrip TestLookup TestEpisodeReportFlags TestParseInterleaved TestApproveTarget"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestVerbRouteReportIsTheCLIs TestWhyTargetFormsResolveToOneEntry TestWhyResolvesMCPInvocationID TestWhyBrokenLinkIsIntegrityRefusal TestWhyScanBound TestWhyShowsTheWorktreeDelta TestWhyShowsTheFailureLine TestLogTailFollowPrintsEachCommitOnce TestProvenanceTrailer"
[ "$MODE" = base ] && INV_CMD="$INV_CMD TestWhyInputLineShowsThePathOfALongEdit TestWhyInputLineLeavesAShortInputUnchanged TestWhyInputLineKeyFieldBounds TestWhyInputLineFallsBackToElide TestWhyProtectedLine TestWhyProtectedJSON"
[ "$MODE" = base ] || INV_CMD="$INV_CMD TestWhyProtectedLine"
leg inv "set -o pipefail
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verbs/ -run '^($(echo $INV_VERBS | tr ' ' '|')$(alt $NEW_V))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$' &&
{ [ -z '$NEW_W' ] || go test -v -count=1 ./host/why/ -run '^($(echo $NEW_W | tr ' ' '|'))\$'; }"
MISSING=""
for t in $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_VG $INV_VERBS $NEW_V $INV_CMD $NEW_C $NEW_W; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || MISSING="$MISSING $t"
done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then FAILS="$FAILS inv_skip_or_fail"; fi
[ -z "$MISSING" ] || { echo "NOT PASSED:$MISSING" | tee -a "$L/legs.txt"; FAILS="$FAILS inv_missing"; }
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/legs.txt"
leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
````

**Expected:** every leg rc 0 except a known-load red, re-run alone once (§1.1's rule); `inv` prints **named PASS**
for every inventory and every new test. The planner's run of this script on the final tree: §1.2 C19.

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, rm, mv, commit, stash, checkout, reset, branch). **`patch -p1 -E` (§2) is
  the only writer** of the tree; no hand edit, no `gofmt -w`, no fuzz, no `WORLD_UPDATE_WHY_GOLDEN=1` (the golden
  ships in M0's patch). A failed `--dry-run` or a red acceptance → **stop and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** M12's
  acceptance runs the whole `cmd/ailang-worldd` package (≈ 2.5–6.5 min at load 8–25); if a call nears 600 s, run
  `accept.sh` again after a timeout rather than backgrounding it.
- **Run**: one run, M0 → M12 → M3t. After each sub-milestone's acceptance, snapshot **cumulatively against the
  base** to `$HOME/.ailang/state/world-iter397/snap/<M>/`: every file added or changed since `18d30a17`
  (repo-relative paths), a **`DELETED`** file listing every repo path removed since the base (one per line), and
  `MANIFEST.sha256`; the two plan files are excluded (they are not under `cmd host docs`). **Never an `rsync` of
  whole directories.** The listing is against the base, not the index, so it is right whether or not the
  controller has committed earlier sub-milestones:
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M0   # the sub-milestone just accepted
  B=18d30a17a058e9adad71c42b6078b538374378d7
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter397-row243-v2a
  D="$HOME/.ailang/state/world-iter397/snap/$M"; rm -rf "$D"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find cmd host docs -type f -print0 2>/dev/null | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED:"; cat "$D/DELETED"; wc -l < "$D/MANIFEST.sha256"
  EOF2
  ```
  Expected: M0 — `STATUS` 2 lines (`A`), `DELETED` empty, MANIFEST 2; M12 — `STATUS` 10 lines, `DELETED` =
  `cmd/ailang-worldd/why.go` and `cmd/ailang-worldd/why_input_test.go`, MANIFEST 8; M3t — `STATUS` 12, the same
  `DELETED`, MANIFEST 10. The controller copies the snapshot over the worktree, `git rm`s each `DELETED` path,
  commits each snapshot as one commit, and checks the numstat against §3's value (`git diff --numstat 18d30a17
  HEAD -- cmd host docs | awk '{s+=$1+$2} END {print s}'`): M0 **77**, M12 **1,039** (with the two `=>` rename
  lines, C14), M3t **1,299**. A different number → stop.
- `export AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests (`accept.sh` sets it). The
  controller's shell is zsh: run each block under `bash <<'EOF2' … EOF2`.
- **Stop conditions**: a dry-run failure; a deleted file still present after M12's patch; any `--- FAIL`, `--- SKIP`
  or `FAIL` line in an acceptance (`ACCEPT_RED`); gofmt output; a numstat mismatch or a missing rename line; any file
  outside §5 changed (`git status --short` after each patch must equal §5's list up to that stage).
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone.
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty`, a real browser profile, or a real
  daemon (§0.1).

---

## 8. Risks

- **R1 — Size margin 1.** The patches are exact and the executor writes nothing else, so the margin is a fact, not a
  forecast; any executor-side growth is a stop (§7), never a hand edit.
- **R2 — Rename pairing** is a size-accounting fact of git's heuristic (C14: `233 285` pairs at 62 % similarity);
  a different git would still produce the same tree (the patches carry no rename headers).
- **R3 — The compatibility layer** (eight aliases, two defined types, thirteen delegations, `sum`, `sumJSON`,
  `contains`, `whyTarget`) is the price of not editing eight CLI files and their tests (design RV4); every shim has
  a mutant (§4).
- **R4 — RS-REST and the two head residuals** (§3.1) are rename-only lines no test reaches; V2b's VE1 drives the
  `rest:` receipt through the route.
- **R5 — Load.** The planner's runs saw load 8–25; every base red was a load timeout in an untouched package (§1.1).

---

## 9. Acceptance criteria (the judge scores against these)

1. **`host/why` exists and is pure at its edge**: the walk reads only through `Getter`; direct imports exclude `os
   os/exec net net/http io/fs database/sql host/store` (VS2a, `TestWhyImportsNoEffect`); the `--result` file read,
   the stdin read and `--expect-store` stay in `cmd/ailang-worldd/why_cli.go`.
2. **The move is a rename**: after M12's commit `git diff --numstat 18d30a17 HEAD -- cmd host docs` shows
   `233 285 {cmd/ailang-worldd => host/why}/why.go` and `3 3 {… => host/why}/why_input_test.go`.
3. **One walk**: the CLI's `runWhy` prints exactly `why.Run`'s `Answer` and returns its `Exit`; `why.Run` over a
   live recording getter and over its replay equals the CLI's bytes for `1`, `head --json`, a tampered input, `9`,
   an exhausted scan (VW2), with exits 0/0/3/1/1 and the other stream empty.
4. **Getter classes**: 404 → absent (✗ "does not resolve" in a chain); error/5xx → the CLI's error line; a deadline
   → the CLI's "not found" (VW3); `Head` and `GetJSON` classes (unit).
5. **One `/v1/head` reader and one `elide`**: cmd's `headState` and `elide` are single `return why.Head(…)` /
   `why.Elide(…)` statements (VW5); the constants are pinned as literals (VW4).
6. **Behaviour neutrality** (the deliverable): VC1's eleven answers equal the golden generated at the base; every
   existing test of `cmd/ailang-worldd` passes **unedited** (M12's whole-package acceptance); §3.1 names every path
   and its driver; residuals are named.
7. **Exactly these pre-existing test files change**, each a move (V2a-D4): `cmd/ailang-worldd/why_input_test.go` →
   `host/why/why_input_test.go` (package line; `elide(p)` → `Elide(p)` at lines 31 and 58), and
   `cmd/ailang-worldd/why_protected_test.go` (its 18 pure subtests and `TestWhyProtectedJSON` moved to
   `host/why/why_protected_test.go` with `whyLink` → `Link`; its two e2e subtests unchanged in place). No other
   `_test.go` is edited; `demo_seed.go` and `cli.go` are not edited (audit inputs 1, 3).
8. **Audit inputs resolved and recorded**: 1 (V2a-D3), 2 (V2a-D4), 3 (V2a-D5), 7 (V2a-D6, C8, C9); 4–6 carried to
   V2b (Appendix C).
9. **Inventories** all PASS (`TestWorldBoundaryDependencyAllowlist`, `TestEveryCommandDirectoryIsAProtectedGroup`,
   `TestBoundaryASTWriteGuard`, `TestProductionContextRoots`, `TestProductionGoSurface`, `TestDaemonDependencyAllowlist`,
   `TestSurfaceVerbIsReadOnly`, `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch`,
   `TestNoParallelWriteForkPackagesOutsideVerifygate`), none edited. No production context root is added (no
   `context.Background`/`TODO` in production code: `why.Run` takes the reader's context).
10. **Size** ≤ 1,300 (`-- cmd host docs`, renames on): 1,299 measured; each sub-milestone's numstat as §3.
11. **Every mutant in §4 behaves as expected**, its id in the ledger; no suite kill rests on a static test.
12. **The round trip is identical** (C17) and the gates (§6) pass, `vet` and `verify_go` after commit included;
    `./scripts/verify_ail.sh` gives the base's result.

**What an operator can do after this sprint:** nothing new — by design. `ailang-worldd why <target>` prints the same
bytes with the same exits; the World app's bar still answers `why` with its CLI line (`terminal`). What changes is
underneath: the walk is one function, `why.Run`, over an injected reader, which V2b's daemon calls in-process.

---

## Appendix C — after this sprint (V2a's close and the V2b handoff)

- **Moved to V2b: design §2.2 whole** (V2a-D2) — `verbs.WhyFlags`/`WhyArgs` (the four flags, today's usage
  strings: `why_cli.go`'s `runWhy` lines), `WhyLeadingTarget` (`runWhy`'s leading-target rule),
  `WhyTrailingTarget` (the `unexpected argument %q` rule), `WhyScanDefault`/`WhyScanMax`/`WhyScanOK`, the CLI's use
  of them, and **VA1** (+ `verbs.WhyScan* == why.Scan*`). **Measured** (a compiled probe on the final stage tree, not shipped: `plan-scratch/v2b-probe`, patch
`plan-scratch/patches/v2b-probe-2.2.patch`): **111** changed lines — `host/verbs/flags.go` +46, its VA1 test +39,
`why_cli.go` +9 −16 (the repoint), `why_seam_test.go` +3 −1 (`verbs.WhyScan* == why.Scan*` in VW4); green on
`host/verbs`, VS1 `TestSurfaceVerbIsReadOnly` (host/verbs stays in its pure set with `strings`), and cmd's `why`
tests incl. VC1 unedited. V2a + §2.2 would be **≈ 1,410 > 1,300**. V2b's budget was ≈ 930 (design §8); with
  this it is ≈ 1,041 (the design's ≈ 930 + the measured 111).
- **Audit inputs 4–6 carry to V2b unchanged**: (4) VS2/T1's pinned `d.*` set must include `whyGetter`, `whyMux`,
  the seven handler method values and any `sync.Once` field; (5) `host/verbs/lookup.go:99`'s stale `why` reason must
  be cleared or replaced on the `result` path, with a test; (6) the 1 MiB cap is on raw bytes — the JSON-encoded body
  can be larger (`\u00XX` escapes); say so in RV2.
- **Residuals carried**: RS-REST (the resolved `rest:` receipt arm — V2b's VE1 runs one through the route);
  `HeadIndex`'s "selected head has no world row" and `ScanFor`'s "no head" (unreachable through `POST /v1/commit`);
  `ScanEntries`' forward page cap (> 500 entries forward); design RV3 (the CLI keeps P4/P5: a deadline prints "not
  found", a read failure is a ✗ — pinned by VW3 as today's behaviour, changed by nobody in V2).
- **V2b handoff**: `why.Run(why.Reader{Ctx: ctx, Get: d.whyGetter(ctx, &tally)}, t, why.Options{…})` is the call;
  `why.MaxResponseBytes` is the getter's cap (VW4 pins it to the CLI's); `why.ParseTarget` is the bar's target
  parse (it refuses `""` and returns the CLI's texts; `-` and `--result` must be refused **before** it, design §2.3);
  `why.Head` is the one `/v1/head` reader. The recording/replay getters in `why_seam_test.go` are the shape VR4's
  transport-parity recording can reuse.
- **Named survivors / equivalents carried**: the three `S` rows are constants no behavioural test reaches, each pinned by VW4's literal —
`W-LOGPAGE` (`LogPageMax` 500 → 499: a page one entry shorter returns the same walk on every store the tests
build, all < 499 entries), `W-MAXRESP` (`MaxResponseBytes` has no reader in V2a; V2b's getter applies it and its
VR5 is the behavioural kill), `W-SCANMAX` (the CLI's `--scan` upper bound: no test at the base or here runs
`--scan 5000`/`5001`; a VC1-style row would cost 2 lines V2a does not have — V2b's VA2 `--scan` out-of-range row
drives it); the two `E` rows (`C-HEALTHCACHE`, `C-READER-CTX`) are named equivalents with their reasons in §4.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M0 — the characterisation pin (VC1), at the base

````diff M0
diff --git a/cmd/ailang-worldd/testdata/why-characterise.golden.jsonl b/cmd/ailang-worldd/testdata/why-characterise.golden.jsonl
new file mode 100644
--- /dev/null
+++ b/cmd/ailang-worldd/testdata/why-characterise.golden.jsonl
@@ -0,0 +1,11 @@
+{"Name":"no target","Stdout":"","Stderr":"ailang-worldd why: a target is required\nusage: ailang-worldd [--addr \u003curl\u003e] why \u003ctarget\u003e [--scan N] [--json] [--expect-store \u003cpath\u003e]\n       ailang-worldd [--addr \u003curl\u003e] why --result \u003cfile\u003e [--scan N] [--json]\n\nWalks one committed invocation back to its whole provenance chain, checking\nevery link (✓ or ✗):\n\n  world    the world ref, recomputed from the entry and the output, and\n           confirmed by GET /v1/worlds/\u003cref\u003e\n  entry    index, writtenBy, prev, transitionFn and the interpreter (with its\n           pin name); the entry hash is recomputed\n  record   invocation id, episode, skill (its hash and its fields agree with\n           the entry)\n  input    your arguments, decoded and elided; when elided, the short\n           top-level fields (path, command, op, …) are printed first\n  plan     the effects planned, or the refusal result\n  effect   per effect record: allowed/failed/denial, budget before→after,\n           request and result sizes; a toolchain call that changed the\n           worktree adds \"worktree: N created, M modified, K deleted\" and\n           its first paths (--json: the whole worktree_delta)\n  output   the committed output; it must equal world.stateRoot\n\nTargets:\n  \u003cindex\u003e | head               a log entry\n  sha256:\u003cref\u003e                 any object (dispatched on its semanticId: a\n                               record, input, output, plan, effect record,\n                               effect request or result), or a world ref\n  a2a:\u003cid\u003e | mcp:\u003cid\u003e | rest:\u003cid\u003e\n                               an invocation id\n  - | --result \u003cfile\u003e          the output of 'call --json-out': its sha256 is\n                               the output ref (the world.plan is the fallback)\n\nA target that is not an index is found by scanning the log backwards from\nhead: --scan N entries (default 500, max 5000) within a 60 s budget; past\neither it prints \"not found in the last N entries\" and exits 1.\n\nExit: 0 every link ✓; 3 any link ✗ (an integrity refusal); 1 not found or\nusage. --json prints the chain as JSON.\n","Code":1}
+{"Name":"malformed sha256","Stdout":"","Stderr":"ailang-worldd why: \"sha256:zz\" is not a content address: hashref: digest is not lowercase hexadecimal: \"zz\"\nusage: ailang-worldd [--addr \u003curl\u003e] why \u003ctarget\u003e [--scan N] [--json] [--expect-store \u003cpath\u003e]\n       ailang-worldd [--addr \u003curl\u003e] why --result \u003cfile\u003e [--scan N] [--json]\n\nWalks one committed invocation back to its whole provenance chain, checking\nevery link (✓ or ✗):\n\n  world    the world ref, recomputed from the entry and the output, and\n           confirmed by GET /v1/worlds/\u003cref\u003e\n  entry    index, writtenBy, prev, transitionFn and the interpreter (with its\n           pin name); the entry hash is recomputed\n  record   invocation id, episode, skill (its hash and its fields agree with\n           the entry)\n  input    your arguments, decoded and elided; when elided, the short\n           top-level fields (path, command, op, …) are printed first\n  plan     the effects planned, or the refusal result\n  effect   per effect record: allowed/failed/denial, budget before→after,\n           request and result sizes; a toolchain call that changed the\n           worktree adds \"worktree: N created, M modified, K deleted\" and\n           its first paths (--json: the whole worktree_delta)\n  output   the committed output; it must equal world.stateRoot\n\nTargets:\n  \u003cindex\u003e | head               a log entry\n  sha256:\u003cref\u003e                 any object (dispatched on its semanticId: a\n                               record, input, output, plan, effect record,\n                               effect request or result), or a world ref\n  a2a:\u003cid\u003e | mcp:\u003cid\u003e | rest:\u003cid\u003e\n                               an invocation id\n  - | --result \u003cfile\u003e          the output of 'call --json-out': its sha256 is\n                               the output ref (the world.plan is the fallback)\n\nA target that is not an index is found by scanning the log backwards from\nhead: --scan N entries (default 500, max 5000) within a 60 s budget; past\neither it prints \"not found in the last N entries\" and exits 1.\n\nExit: 0 every link ✓; 3 any link ✗ (an integrity refusal); 1 not found or\nusage. --json prints the chain as JSON.\n","Code":1}
+{"Name":"absent sha256","Stdout":"","Stderr":"ailang-worldd why: sha256:0000000000000000000000000000000000000000000000000000000000000000 is neither an object nor a world in this store\n","Code":1}
+{"Name":"flag before target","Stdout":"{\n  \"entry\": 1,\n  \"matchedBy\": \"entry index\",\n  \"coordinator\": true,\n  \"links\": [\n    {\n      \"name\": \"world\",\n      \"ok\": true,\n      \"ref\": \"sha256:14e78926e09f5db554f2523bfce82bc25adf42d903060674883127db88f18cb6\",\n      \"detail\": \"revision 1, stateRoot sha256:999ee9bde6bea137f94b3317dd4288c2e71513127ca8ac1fdb1271f65336abd6, logHead sha256:e616b00c1a37c8b11f8c7cc7fda616a4b601125ee173618cb7ff6de21cb26a82 (recomputed)\"\n    },\n    {\n      \"name\": \"entry\",\n      \"ok\": true,\n      \"ref\": \"sha256:e616b00c1a37c8b11f8c7cc7fda616a4b601125ee173618cb7ff6de21cb26a82\",\n      \"detail\": \"#1 writtenBy coordinator:a2a prev sha256:bc5596099b820473274c98efa4076422f0f44544264c92cc6aa6ff281a431a5a transitionFn sha256:02d362e69d9240eecc26afd1e96d1e7056da3e591614ae6d487dc251d9189ff4 interpreter sha256:9666d9e8899447735cb9897b77dbb121754fd4d503609758755bd3fdae3a4b22 (not the daemon's pinned interpreter)\"\n    },\n    {\n      \"name\": \"record\",\n      \"ok\": true,\n      \"ref\": \"sha256:e77abd8a9912c0cb9f44efdefa0d54a05d6b0c58fd7ef83e05abad6675ceb623\",\n      \"detail\": \"world/invocation-record/v2 invocation a2a:ep1:task-1 episode ep1 skill ailang-read\"\n    },\n    {\n      \"name\": \"input\",\n      \"ok\": true,\n      \"ref\": \"sha256:db1faa2de9a84a540eda94cd70b242170715d146a6bf25ade75e36ddd121c29a\",\n      \"detail\": \"{\\\"path\\\":\\\"data.txt\\\"}\"\n    },\n    {\n      \"name\": \"plan\",\n      \"ok\": true,\n      \"ref\": \"sha256:704017c470a4c3bfc5963d16efeb31c0fc5c5cd1c598f4561b76f7ea04749709\",\n      \"detail\": \"1 effect(s): e1 Workspace.Read@worktree cost 1\"\n    },\n    {\n      \"name\": \"effect\",\n      \"ok\": true,\n      \"ref\": \"sha256:127f1e6b16651c3521b5164f1968794ea5bf3c281fb8770dd49ea719caa068da\",\n      \"detail\": \"e1 allowed Workspace.Read@worktree cost 1 budget 20→19 request 47 B result 21 B\"\n    },\n    {\n      \"name\": \"output\",\n      \"ok\": true,\n      \"ref\": \"sha256:999ee9bde6bea137f94b3317dd4288c2e71513127ca8ac1fdb1271f65336abd6\",\n      \"detail\": \"243 B = world.stateRoot\"\n    }\n  ],\n  \"ok\": true\n}\n","Stderr":"","Code":0}
+{"Name":"two targets","Stdout":"","Stderr":"ailang-worldd why: unexpected argument \"2\"\n","Code":1}
+{"Name":"target and --result","Stdout":"","Stderr":"ailang-worldd why: give either a target or --result, not both\nusage: ailang-worldd [--addr \u003curl\u003e] why \u003ctarget\u003e [--scan N] [--json] [--expect-store \u003cpath\u003e]\n       ailang-worldd [--addr \u003curl\u003e] why --result \u003cfile\u003e [--scan N] [--json]\n\nWalks one committed invocation back to its whole provenance chain, checking\nevery link (✓ or ✗):\n\n  world    the world ref, recomputed from the entry and the output, and\n           confirmed by GET /v1/worlds/\u003cref\u003e\n  entry    index, writtenBy, prev, transitionFn and the interpreter (with its\n           pin name); the entry hash is recomputed\n  record   invocation id, episode, skill (its hash and its fields agree with\n           the entry)\n  input    your arguments, decoded and elided; when elided, the short\n           top-level fields (path, command, op, …) are printed first\n  plan     the effects planned, or the refusal result\n  effect   per effect record: allowed/failed/denial, budget before→after,\n           request and result sizes; a toolchain call that changed the\n           worktree adds \"worktree: N created, M modified, K deleted\" and\n           its first paths (--json: the whole worktree_delta)\n  output   the committed output; it must equal world.stateRoot\n\nTargets:\n  \u003cindex\u003e | head               a log entry\n  sha256:\u003cref\u003e                 any object (dispatched on its semanticId: a\n                               record, input, output, plan, effect record,\n                               effect request or result), or a world ref\n  a2a:\u003cid\u003e | mcp:\u003cid\u003e | rest:\u003cid\u003e\n                               an invocation id\n  - | --result \u003cfile\u003e          the output of 'call --json-out': its sha256 is\n                               the output ref (the world.plan is the fallback)\n\nA target that is not an index is found by scanning the log backwards from\nhead: --scan N entries (default 500, max 5000) within a 60 s budget; past\neither it prints \"not found in the last N entries\" and exits 1.\n\nExit: 0 every link ✓; 3 any link ✗ (an integrity refusal); 1 not found or\nusage. --json prints the chain as JSON.\n","Code":1}
+{"Name":"--result absent file","Stdout":"","Stderr":"ailang-worldd why: read --result $TMP/absent: no such file or directory\nusage: ailang-worldd [--addr \u003curl\u003e] why \u003ctarget\u003e [--scan N] [--json] [--expect-store \u003cpath\u003e]\n       ailang-worldd [--addr \u003curl\u003e] why --result \u003cfile\u003e [--scan N] [--json]\n\nWalks one committed invocation back to its whole provenance chain, checking\nevery link (✓ or ✗):\n\n  world    the world ref, recomputed from the entry and the output, and\n           confirmed by GET /v1/worlds/\u003cref\u003e\n  entry    index, writtenBy, prev, transitionFn and the interpreter (with its\n           pin name); the entry hash is recomputed\n  record   invocation id, episode, skill (its hash and its fields agree with\n           the entry)\n  input    your arguments, decoded and elided; when elided, the short\n           top-level fields (path, command, op, …) are printed first\n  plan     the effects planned, or the refusal result\n  effect   per effect record: allowed/failed/denial, budget before→after,\n           request and result sizes; a toolchain call that changed the\n           worktree adds \"worktree: N created, M modified, K deleted\" and\n           its first paths (--json: the whole worktree_delta)\n  output   the committed output; it must equal world.stateRoot\n\nTargets:\n  \u003cindex\u003e | head               a log entry\n  sha256:\u003cref\u003e                 any object (dispatched on its semanticId: a\n                               record, input, output, plan, effect record,\n                               effect request or result), or a world ref\n  a2a:\u003cid\u003e | mcp:\u003cid\u003e | rest:\u003cid\u003e\n                               an invocation id\n  - | --result \u003cfile\u003e          the output of 'call --json-out': its sha256 is\n                               the output ref (the world.plan is the fallback)\n\nA target that is not an index is found by scanning the log backwards from\nhead: --scan N entries (default 500, max 5000) within a 60 s budget; past\neither it prints \"not found in the last N entries\" and exits 1.\n\nExit: 0 every link ✓; 3 any link ✗ (an integrity refusal); 1 not found or\nusage. --json prints the chain as JSON.\n","Code":1}
+{"Name":"rest: unknown receipt","Stdout":"","Stderr":"ailang-worldd why: receipt rest:nobody is not-started: nothing committed\n","Code":1}
+{"Name":"--expect-store this store","Stdout":"why: entry 1 (matched by entry index)\n  ✓ world   sha256:14e78926e09f5db554f2523bfce82bc25adf42d903060674883127db88f18cb6  revision 1, stateRoot sha256:999ee9bde6bea137f94b3317dd4288c2e71513127ca8ac1fdb1271f65336abd6, logHead sha256:e616b00c1a37c8b11f8c7cc7fda616a4b601125ee173618cb7ff6de21cb26a82 (recomputed)\n  ✓ entry   sha256:e616b00c1a37c8b11f8c7cc7fda616a4b601125ee173618cb7ff6de21cb26a82  #1 writtenBy coordinator:a2a prev sha256:bc5596099b820473274c98efa4076422f0f44544264c92cc6aa6ff281a431a5a transitionFn sha256:02d362e69d9240eecc26afd1e96d1e7056da3e591614ae6d487dc251d9189ff4 interpreter sha256:9666d9e8899447735cb9897b77dbb121754fd4d503609758755bd3fdae3a4b22 (not the daemon's pinned interpreter)\n  ✓ record  sha256:e77abd8a9912c0cb9f44efdefa0d54a05d6b0c58fd7ef83e05abad6675ceb623  world/invocation-record/v2 invocation a2a:ep1:task-1 episode ep1 skill ailang-read\n  ✓ input   sha256:db1faa2de9a84a540eda94cd70b242170715d146a6bf25ade75e36ddd121c29a  {\"path\":\"data.txt\"}\n  ✓ plan    sha256:704017c470a4c3bfc5963d16efeb31c0fc5c5cd1c598f4561b76f7ea04749709  1 effect(s): e1 Workspace.Read@worktree cost 1\n  ✓ effect  sha256:127f1e6b16651c3521b5164f1968794ea5bf3c281fb8770dd49ea719caa068da  e1 allowed Workspace.Read@worktree cost 1 budget 20→19 request 47 B result 21 B\n  ✓ output  sha256:999ee9bde6bea137f94b3317dd4288c2e71513127ca8ac1fdb1271f65336abd6  243 B = world.stateRoot\nall 7 link(s) verified\n","Stderr":"","Code":0}
+{"Name":"- not a result","Stdout":"","Stderr":"ailang-worldd why: not found in the last 3 entries (output (sha256 of the result))\n","Code":1}
+{"Name":"index past head","Stdout":"","Stderr":"ailang-worldd why: log entry 9 not found\n","Code":1}
diff --git a/cmd/ailang-worldd/why_characterise_test.go b/cmd/ailang-worldd/why_characterise_test.go
new file mode 100644
--- /dev/null
+++ b/cmd/ailang-worldd/why_characterise_test.go
@@ -0,0 +1,66 @@
+package main
+
+// Row 243 V2a (plan §3.1, VC1): the why paths no other test drove, pinned to the bytes the CLI
+// printed before the walk moved to host/why (the golden: WORLD_UPDATE_WHY_GOLDEN=1 on the base tree).
+
+import (
+	"context"
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+)
+
+const whyGolden = "testdata/why-characterise.golden.jsonl"
+
+func TestWhyCharacterisedPaths(t *testing.T) {
+	r := newWalkRig(t)
+	inv := r.invoke("ep1", "ailang-read", "task-1", "hello\n")
+	h, err := (worldReader{ctx: context.Background(), c: newClient(r.url)}).health()
+	if err != nil {
+		t.Fatal(err)
+	}
+	dir := t.TempDir()
+	result := filepath.Join(dir, "result.json")
+	writeFileT(t, result, inv.outputText+"\n")
+	type run struct {
+		Name, Stdout, Stderr string
+		Code                 int
+	}
+	var got []run
+	for _, c := range []struct {
+		name, stdin string
+		args        []string
+	}{
+		{"no target", "", nil},
+		{"malformed sha256", "", []string{"sha256:zz"}},
+		{"absent sha256", "", []string{"sha256:" + strings.Repeat("0", 64)}},
+		{"flag before target", "", []string{"--json", "1"}},
+		{"two targets", "", []string{"1", "2"}},
+		{"target and --result", "", []string{"1", "--result", result}},
+		{"--result absent file", "", []string{"--result", filepath.Join(dir, "absent")}},
+		{"rest: unknown receipt", "", []string{"rest:nobody"}},
+		{"--expect-store this store", "", []string{"1", "--expect-store", h.DBPath}},
+		{"- not a result", "x", []string{"-", "--scan", "3"}},
+		{"index past head", "", []string{"9"}},
+	} {
+		code, out, errOut := whyRun(t, r.url, c.stdin, c.args...)
+		got = append(got, run{c.name, strings.ReplaceAll(out, dir, "$TMP"), strings.ReplaceAll(errOut, dir, "$TMP"), code})
+	}
+	var b []byte
+	for _, g := range got { // one JSON line per run
+		line, err := json.Marshal(g)
+		if err != nil {
+			t.Fatal(err)
+		}
+		b = append(append(b, line...), '\n')
+	}
+	if os.Getenv("WORLD_UPDATE_WHY_GOLDEN") == "1" {
+		writeFileT(t, whyGolden, string(b))
+	}
+	want, err := os.ReadFile(whyGolden)
+	if err != nil || string(want) != string(b) {
+		t.Fatalf("why's output moved from the base's (%v):\n%s", err, b)
+	}
+}
````

### Patch M12 — the move and the seam (design M1 + M2): `host/why`, `why_cli.go`, `mcpclient.go`, the two test files

````diff M12
diff --git a/cmd/ailang-worldd/mcpclient.go b/cmd/ailang-worldd/mcpclient.go
--- a/cmd/ailang-worldd/mcpclient.go
+++ b/cmd/ailang-worldd/mcpclient.go
@@ -34,6 +34,8 @@ import (
 	"sort"
 	"strings"
 	"time"
+
+	"github.com/sunholo-data/ailang-world/host/why"
 )
 
 // exitIntegrity is the row-138 integrity refusal: a broken provenance link
@@ -257,18 +259,7 @@ func (m mcpClient) post(ctx context.Context, method string, id int, params any)
 // headState reads GET /v1/head: the selected world ref, or present=false on
 // the 404 "no world head has been selected yet".
 func headState(ctx context.Context, c *client) (ref string, present bool, err error) {
-	status, body, err := c.do(ctx, http.MethodGet, "/v1/head", nil)
-	if err != nil {
-		return "", false, err
-	}
-	switch status {
-	case http.StatusOK:
-		return strings.TrimSpace(string(body)), true, nil
-	case http.StatusNotFound:
-		return "", false, nil
-	default:
-		return "", false, fmt.Errorf("GET /v1/head returned HTTP %d: %s", status, strings.TrimSpace(string(body)))
-	}
+	return why.Head(ctx, clientGetter(c))
 }
 
 // worldRevision resolves a world ref to its revision (= the entry index of
@@ -592,8 +583,6 @@ func callArguments(kv, kvJSON []string, whole string, stdin io.Reader) (json.Raw
 	return json.Marshal(obj)
 }
 
-const elideAt = 160
-
 // printCallResult renders a result: each top-level field but world (long
 // values elided with byte counts), then the world block.
 func printCallResult(w io.Writer, text string, fields map[string]json.RawMessage) {
@@ -647,14 +636,8 @@ func printCallResult(w io.Writer, text string, fields map[string]json.RawMessage
 	}
 }
 
-// elide shortens a value to elideAt bytes plus its byte count, on one line.
-func elide(s string) string {
-	one := strings.ReplaceAll(s, "\n", `\n`)
-	if len(one) <= elideAt {
-		return one
-	}
-	return fmt.Sprintf("%s… (%d bytes)", strings.ToValidUTF8(one[:elideAt], ""), len(s))
-}
+// elide shortens a value to why.ElideAt bytes plus its byte count, on one line.
+func elide(s string) string { return why.Elide(s) }
 
 func oneLine(s string, n int) string {
 	s = strings.Join(strings.Fields(s), " ")
diff --git a/cmd/ailang-worldd/why.go b/cmd/ailang-worldd/why.go
deleted file mode 100644
--- a/cmd/ailang-worldd/why.go
+++ /dev/null
@@ -1,1050 +0,0 @@
-package main
-
-// Row 138 M2 (design_docs/planned/w-worldd-developer-cli.md §3.3): `why`
-// walks a result, an object, an invocation id or an entry back to its whole
-// provenance chain, recomputing and checking every link.
-//
-// An MCP result carries no entry index, invocation id or record ref (V12), and
-// /v1/receipts refuses a2a:/mcp: ids (V13), so a target that is not an index is
-// found by scanning the log backwards from head — one log page of 500 at a
-// time and one object GET per entry — bounded by --scan (default 500, max
-// 5000) and a 60 s budget. The reads are unauthenticated GETs (V15); `why`
-// trusts that read surface (R-CLI-2) but verifies every content address it is
-// handed against the bytes it is handed.
-
-import (
-	"context"
-	"encoding/json"
-	"errors"
-	"flag"
-	"fmt"
-	"io"
-	"net/http"
-	"net/url"
-	"os"
-	"sort"
-	"strconv"
-	"strings"
-	"time"
-
-	"github.com/sunholo-data/ailang-world/host/broker"
-	"github.com/sunholo-data/ailang-world/host/coordinator"
-	"github.com/sunholo-data/ailang-world/host/hashref"
-)
-
-const (
-	whyBudget      = 60 * time.Second
-	whyScanDefault = 500
-	whyScanMax     = 5000
-	logPageMax     = 500
-	frozenRigPort  = 47_330 // historical trailer recovery, never a default
-)
-
-// ---------------------------------------------------------------------------
-// the read surface
-// ---------------------------------------------------------------------------
-
-type wireHeader struct {
-	EntryIndex     int64  `json:"entryIndex"`
-	SemanticsEpoch int64  `json:"semanticsEpoch"`
-	TransitionFn   string `json:"transitionFn"`
-	Interpreter    string `json:"interpreter"`
-	PrevEntryHash  string `json:"prevEntryHash"`
-	WrittenBy      string `json:"writtenBy"`
-}
-
-type wireEntry struct {
-	Header        wireHeader `json:"header"`
-	EntryHash     string     `json:"entryHash"`
-	TransitionRef string     `json:"transitionRef"`
-}
-
-type wireObject struct {
-	Hash       string `json:"hash"`
-	SemanticID string `json:"semanticId"`
-	Payload    []byte `json:"payload"` // base64 on the wire
-}
-
-type wireWorld struct {
-	Ref       string `json:"ref"`
-	Revision  int64  `json:"revision"`
-	StateRoot string `json:"stateRoot"`
-	LogHead   string `json:"logHead"`
-}
-
-// invRecord is a world/invocation-record/v1 or /v2 payload.
-type invRecord struct {
-	SemanticID     string   `json:"-"`
-	InvocationID   string   `json:"invocationId"`
-	EpisodeID      string   `json:"episodeId"`
-	SkillID        string   `json:"skillId"`
-	TransitionFn   string   `json:"transitionFn"`
-	Interpreter    string   `json:"interpreter"`
-	SemanticsEpoch int64    `json:"semanticsEpoch"`
-	Input          string   `json:"input"`
-	Output         string   `json:"output"`
-	Plan           string   `json:"plan"`
-	Effects        []string `json:"effects"`
-	payload        []byte
-}
-
-// worldReader is the GET side of the daemon, under one command budget.
-type worldReader struct {
-	ctx         context.Context
-	c           *client
-	healthCache *healthInfo
-}
-
-// getJSON decodes a 200 into v; a 404 is (false, nil); anything else errors.
-func (r worldReader) getJSON(path string, v any) (bool, error) {
-	if err := r.ctx.Err(); err != nil {
-		return false, err
-	}
-	status, body, err := r.c.do(r.ctx, http.MethodGet, path, nil)
-	if err != nil {
-		return false, err
-	}
-	switch status {
-	case http.StatusOK:
-		if err := json.Unmarshal(body, v); err != nil {
-			return false, fmt.Errorf("GET %s: malformed body: %v", path, err)
-		}
-		return true, nil
-	case http.StatusNotFound:
-		return false, nil
-	default:
-		return false, fmt.Errorf("GET %s returned HTTP %d: %s", path, status, strings.TrimSpace(string(body)))
-	}
-}
-
-func (r worldReader) entry(i int64) (wireEntry, bool, error) {
-	var e wireEntry
-	ok, err := r.getJSON("/v1/log/"+strconv.FormatInt(i, 10), &e)
-	return e, ok, err
-}
-
-func (r worldReader) page(from int64, limit int) ([]wireEntry, error) {
-	var p struct {
-		Items []wireEntry `json:"items"`
-	}
-	_, err := r.getJSON(fmt.Sprintf("/v1/log?from=%d&limit=%d", from, limit), &p)
-	return p.Items, err
-}
-
-func (r worldReader) object(ref string) (wireObject, bool, error) {
-	var o wireObject
-	ok, err := r.getJSON("/v1/objects/"+url.PathEscape(ref)+"?payload=true", &o)
-	return o, ok, err
-}
-
-func (r worldReader) world(ref string) (wireWorld, bool, error) {
-	var w wireWorld
-	ok, err := r.getJSON("/v1/worlds/"+url.PathEscape(ref), &w)
-	return w, ok, err
-}
-
-// headIndex is the selected world's revision (= its log head's entry index).
-func (r worldReader) headIndex() (int64, bool, error) {
-	ref, ok, err := headState(r.ctx, r.c)
-	if err != nil || !ok {
-		return 0, ok, err
-	}
-	w, found, err := r.world(ref)
-	if err != nil {
-		return 0, false, err
-	}
-	if !found {
-		return 0, false, fmt.Errorf("the selected head %s has no world row", ref)
-	}
-	return w.Revision, true, nil
-}
-
-type healthInfo struct {
-	DBPath             string `json:"db_path"`
-	InterpreterRef     string `json:"interpreter_ref"`
-	InterpreterVersion string `json:"interpreter_version"`
-}
-
-func (r worldReader) health() (healthInfo, error) {
-	if r.healthCache != nil {
-		return *r.healthCache, nil
-	}
-	var h healthInfo
-	_, err := r.getJSON("/v1/health", &h)
-	return h, err
-}
-
-// record fetches an entry's transitionRef and decodes it when it is a
-// coordinator invocation record; rec is nil for any other entry.
-func (r worldReader) record(e wireEntry) (*invRecord, wireObject, error) {
-	obj, ok, err := r.object(e.TransitionRef)
-	if err != nil || !ok {
-		if err == nil {
-			err = fmt.Errorf("entry %d: transitionRef %s does not resolve", e.Header.EntryIndex, e.TransitionRef)
-		}
-		return nil, obj, err
-	}
-	if obj.SemanticID != coordinator.RecordV1 && obj.SemanticID != coordinator.RecordV2 {
-		return nil, obj, nil
-	}
-	var rec invRecord
-	if err := json.Unmarshal(obj.Payload, &rec); err != nil {
-		return nil, obj, fmt.Errorf("entry %d: undecodable %s: %v", e.Header.EntryIndex, obj.SemanticID, err)
-	}
-	rec.SemanticID, rec.payload = obj.SemanticID, obj.Payload
-	return &rec, obj, nil
-}
-
-// errScanExhausted reports a bounded scan that ran out of entries or budget.
-var errScanExhausted = errors.New("scan exhausted")
-
-// scanEntries visits the entries [from, to] — backwards from `to` when
-// reverse — one page of up to 500 at a time, decoding each entry's record
-// (one object GET per entry). It stops at the first visit that returns
-// true, after max entries, or when the budget ends.
-func (r worldReader) scanEntries(from, to int64, reverse bool, max int,
-	visit func(e wireEntry, rec *invRecord) (bool, error)) (int, error) {
-	seen := 0
-	for lo, hi := from, to; lo <= hi && seen < max; {
-		pFrom := lo
-		if reverse {
-			pFrom = hi - logPageMax + 1
-			if pFrom < lo {
-				pFrom = lo
-			}
-		}
-		n := hi - pFrom + 1
-		if !reverse && n > logPageMax {
-			n = logPageMax
-		}
-		items, err := r.page(pFrom, int(n))
-		if err != nil {
-			return seen, err
-		}
-		if len(items) == 0 {
-			return seen, nil
-		}
-		for k := range items {
-			i := k
-			if reverse {
-				i = len(items) - 1 - k
-			}
-			if seen >= max {
-				break
-			}
-			seen++
-			rec, _, err := r.record(items[i])
-			if err != nil {
-				return seen, err
-			}
-			stop, err := visit(items[i], rec)
-			if err != nil || stop {
-				return seen, err
-			}
-		}
-		if reverse {
-			hi = pFrom - 1
-		} else {
-			lo = pFrom + int64(len(items))
-		}
-	}
-	return seen, nil
-}
-
-// ---------------------------------------------------------------------------
-// why
-// ---------------------------------------------------------------------------
-
-const whyHelp = `usage: ailang-worldd [--addr <url>] why <target> [--scan N] [--json] [--expect-store <path>]
-       ailang-worldd [--addr <url>] why --result <file> [--scan N] [--json]
-
-Walks one committed invocation back to its whole provenance chain, checking
-every link (✓ or ✗):
-
-  world    the world ref, recomputed from the entry and the output, and
-           confirmed by GET /v1/worlds/<ref>
-  entry    index, writtenBy, prev, transitionFn and the interpreter (with its
-           pin name); the entry hash is recomputed
-  record   invocation id, episode, skill (its hash and its fields agree with
-           the entry)
-  input    your arguments, decoded and elided; when elided, the short
-           top-level fields (path, command, op, …) are printed first
-  plan     the effects planned, or the refusal result
-  effect   per effect record: allowed/failed/denial, budget before→after,
-           request and result sizes; a toolchain call that changed the
-           worktree adds "worktree: N created, M modified, K deleted" and
-           its first paths (--json: the whole worktree_delta)
-  output   the committed output; it must equal world.stateRoot
-
-Targets:
-  <index> | head               a log entry
-  sha256:<ref>                 any object (dispatched on its semanticId: a
-                               record, input, output, plan, effect record,
-                               effect request or result), or a world ref
-  a2a:<id> | mcp:<id> | rest:<id>
-                               an invocation id
-  - | --result <file>          the output of 'call --json-out': its sha256 is
-                               the output ref (the world.plan is the fallback)
-
-A target that is not an index is found by scanning the log backwards from
-head: --scan N entries (default 500, max 5000) within a 60 s budget; past
-either it prints "not found in the last N entries" and exits 1.
-
-Exit: 0 every link ✓; 3 any link ✗ (an integrity refusal); 1 not found or
-usage. --json prints the chain as JSON.
-`
-
-type whyTarget struct {
-	kind   string // index, head, ref, invocation, result
-	index  int64
-	text   string
-	result []byte
-}
-
-func parseWhyTarget(arg, resultFile string, stdin io.Reader) (whyTarget, error) {
-	if resultFile != "" {
-		if arg != "" {
-			return whyTarget{}, errors.New("give either a target or --result, not both")
-		}
-		b, err := os.ReadFile(resultFile)
-		if err != nil {
-			return whyTarget{}, fmt.Errorf("read --result %s: %v", resultFile, errors.Unwrap(err))
-		}
-		return whyTarget{kind: "result", result: trimOneNewline(b)}, nil
-	}
-	switch {
-	case arg == "":
-		return whyTarget{}, errors.New("a target is required")
-	case arg == "-":
-		b, err := io.ReadAll(io.LimitReader(stdin, maxClientResponseBytes+1))
-		if err != nil {
-			return whyTarget{}, fmt.Errorf("read stdin: %v", err)
-		}
-		return whyTarget{kind: "result", result: trimOneNewline(b)}, nil
-	case arg == "head":
-		return whyTarget{kind: "head"}, nil
-	case strings.HasPrefix(arg, "a2a:") || strings.HasPrefix(arg, "mcp:") || strings.HasPrefix(arg, "rest:"):
-		return whyTarget{kind: "invocation", text: arg}, nil
-	case strings.HasPrefix(arg, "sha256:"):
-		if _, err := hashref.Parse(arg); err != nil {
-			return whyTarget{}, fmt.Errorf("%q is not a content address: %v", arg, err)
-		}
-		return whyTarget{kind: "ref", text: arg}, nil
-	}
-	n, err := strconv.ParseInt(arg, 10, 64)
-	if err != nil || n < 0 {
-		return whyTarget{}, fmt.Errorf("%q is not an entry index, head, sha256:<ref>, a2a:/mcp:/rest:<id>, or -", arg)
-	}
-	return whyTarget{kind: "index", index: n}, nil
-}
-
-// trimOneNewline drops the one newline `call --json-out` appends.
-func trimOneNewline(b []byte) []byte {
-	if n := len(b); n > 0 && b[n-1] == '\n' {
-		return b[:n-1]
-	}
-	return b
-}
-
-func runWhy(addr string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
-	target := ""
-	if len(args) > 0 && (!strings.HasPrefix(args[0], "-") || args[0] == "-") {
-		target, args = args[0], args[1:]
-	}
-	fs := flag.NewFlagSet("ailang-worldd why", flag.ContinueOnError)
-	fs.SetOutput(stderr)
-	fs.Usage = func() {}
-	expectStore := fs.String("expect-store", "", "refuse unless the daemon serves this store (World-Provenance store= path)")
-	resultFile := fs.String("result", "", "a 'call --json-out' output file")
-	scan := fs.Int("scan", whyScanDefault, "entries to scan back from head (max 5000)")
-	asJSON := fs.Bool("json", false, "print the chain as JSON")
-	if code, done := parseVerbFlags(fs, args, whyHelp, stdout, stderr); done {
-		return code
-	}
-	if fs.NArg() > 0 {
-		if target != "" || fs.NArg() > 1 {
-			fmt.Fprintf(stderr, "ailang-worldd why: unexpected argument %q\n", fs.Arg(fs.NArg()-1))
-			return exitUsage
-		}
-		target = fs.Arg(0)
-	}
-	if *scan < 1 || *scan > whyScanMax {
-		fmt.Fprintf(stderr, "ailang-worldd why: --scan must be 1..%d\n", whyScanMax)
-		return exitUsage
-	}
-	t, err := parseWhyTarget(target, *resultFile, stdin)
-	if err != nil {
-		fmt.Fprintf(stderr, "ailang-worldd why: %v\n%s", err, whyHelp)
-		return exitUsage
-	}
-	ctx, cancel := budgetContext(whyBudget)
-	defer cancel()
-	r := worldReader{ctx: ctx, c: newClient(addr)}
-	if *expectStore != "" {
-		h, err := r.health()
-		if err != nil {
-			fmt.Fprintf(stderr, "ailang-worldd why: %v\n", err)
-			return exitUsage
-		}
-		r.healthCache = &h
-		if !sameStorePath(h.DBPath, *expectStore) {
-			return serviceMessage(stderr,
-				fmt.Sprintf("ailang-worldd why: the daemon at %s serves store %s, not %s", addr, h.DBPath, *expectStore),
-				"entry indexes are per store, so the same index in another store is another entry (a World-Provenance trailer names its store as store=<path>)",
-				fmt.Sprintf("pass the --addr of the daemon serving %s (the frozen rig: start it from daemon-cmdline.txt on 127.0.0.1:%d), or drop --expect-store", *expectStore, frozenRigPort), exitUsage)
-		}
-	}
-	e, matchedBy, err := r.locate(t, *scan)
-	if err != nil {
-		if errors.Is(err, errScanExhausted) || errors.Is(err, context.DeadlineExceeded) {
-			fmt.Fprintf(stderr, "ailang-worldd why: not found in the last %d entries (%s)\n", *scan, matchedBy)
-			return exitUsage
-		}
-		fmt.Fprintf(stderr, "ailang-worldd why: %v\n", err)
-		return exitUsage
-	}
-	ch, err := r.buildChain(e, t.result)
-	if err != nil {
-		fmt.Fprintf(stderr, "ailang-worldd why: %v\n", err)
-		return exitUsage
-	}
-	ch.MatchedBy = matchedBy
-	if *asJSON {
-		b, _ := json.MarshalIndent(ch, "", "  ")
-		fmt.Fprintln(stdout, string(b))
-	} else {
-		printChain(stdout, ch)
-	}
-	if !ch.OK {
-		return exitIntegrity
-	}
-	return exitOK
-}
-
-// locate resolves a target to one log entry. matchedBy names how.
-func (r worldReader) locate(t whyTarget, scan int) (wireEntry, string, error) {
-	byIndex := func(i int64, how string) (wireEntry, string, error) {
-		e, ok, err := r.entry(i)
-		if err != nil {
-			return e, how, err
-		}
-		if !ok {
-			return e, how, fmt.Errorf("log entry %d not found", i)
-		}
-		return e, how, nil
-	}
-	switch t.kind {
-	case "index":
-		return byIndex(t.index, "entry index")
-	case "head":
-		h, ok, err := r.headIndex()
-		if err != nil {
-			return wireEntry{}, "head", err
-		}
-		if !ok {
-			return wireEntry{}, "head", errors.New("no world head has been selected yet")
-		}
-		return byIndex(h, "head")
-	case "invocation":
-		if strings.HasPrefix(t.text, "rest:") {
-			var rc struct {
-				State     string `json:"state"`
-				ResultRef string `json:"resultRef"`
-			}
-			if _, err := r.getJSON("/v1/receipts/"+url.PathEscape(t.text), &rc); err != nil {
-				return wireEntry{}, "rest receipt", err
-			}
-			if rc.State != "resolved" || rc.ResultRef == "" {
-				return wireEntry{}, "rest receipt", fmt.Errorf("receipt %s is %s: nothing committed", t.text, rc.State)
-			}
-			return r.byWorld(rc.ResultRef, "rest receipt → world")
-		}
-		return r.scanFor(scan, "invocationId", func(e wireEntry, rec *invRecord) bool {
-			return rec != nil && rec.InvocationID == t.text
-		})
-	case "result":
-		ref := hashref.SumSHA256(t.result).String()
-		e, how, err := r.scanFor(scan, "output (sha256 of the result)", func(e wireEntry, rec *invRecord) bool {
-			return rec != nil && rec.Output == ref
-		})
-		if err == nil || !errors.Is(err, errScanExhausted) {
-			return e, how, err
-		}
-		var out struct {
-			World struct {
-				Plan string `json:"plan"`
-			} `json:"world"`
-		}
-		if json.Unmarshal(t.result, &out) != nil || out.World.Plan == "" {
-			return e, how, err
-		}
-		return r.scanFor(scan, "world.plan (fallback: the result bytes are not a committed output)", func(e wireEntry, rec *invRecord) bool {
-			return rec != nil && rec.Plan == out.World.Plan
-		})
-	}
-	// A content address: dispatch on what it names.
-	obj, ok, err := r.object(t.text)
-	if err != nil {
-		return wireEntry{}, "object", err
-	}
-	if !ok {
-		return r.byWorld(t.text, "world ref")
-	}
-	ref := t.text
-	switch obj.SemanticID {
-	case coordinator.RecordV1, coordinator.RecordV2:
-		return r.scanFor(scan, "record", func(e wireEntry, _ *invRecord) bool { return e.TransitionRef == ref })
-	case coordinator.OutputV1:
-		return r.scanFor(scan, "output", func(_ wireEntry, rec *invRecord) bool { return rec != nil && rec.Output == ref })
-	case coordinator.InputV1:
-		return r.scanFor(scan, "input", func(_ wireEntry, rec *invRecord) bool { return rec != nil && rec.Input == ref })
-	case coordinator.EffectPlanV1:
-		return r.scanFor(scan, "plan", func(_ wireEntry, rec *invRecord) bool { return rec != nil && rec.Plan == ref })
-	case broker.EffectRecordV1:
-		return r.scanFor(scan, "effect record", func(_ wireEntry, rec *invRecord) bool { return rec != nil && contains(rec.Effects, ref) })
-	case broker.EffectResultV1, broker.EffectRequestV1:
-		return r.scanFor(scan, "effect "+strings.TrimSuffix(strings.TrimPrefix(obj.SemanticID, "world/effect-"), "/v1"),
-			func(_ wireEntry, rec *invRecord) bool {
-				if rec == nil {
-					return false
-				}
-				for _, er := range rec.Effects {
-					if d, ok := r.effectRecord(er); ok && (d.ResultRef.String() == ref || d.RequestRef.String() == ref) {
-						return true
-					}
-				}
-				return false
-			})
-	default:
-		return r.scanFor(scan, "transitionRef ("+obj.SemanticID+")", func(e wireEntry, _ *invRecord) bool { return e.TransitionRef == ref })
-	}
-}
-
-func (r worldReader) effectRecord(ref string) (broker.EffectRecord, bool) {
-	obj, ok, err := r.object(ref)
-	if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 {
-		return broker.EffectRecord{}, false
-	}
-	d, err := broker.DecodeRecord(obj.Payload)
-	return d, err == nil
-}
-
-// byWorld resolves a world ref to the entry its log head names.
-func (r worldReader) byWorld(ref, how string) (wireEntry, string, error) {
-	w, ok, err := r.world(ref)
-	if err != nil {
-		return wireEntry{}, how, err
-	}
-	if !ok {
-		return wireEntry{}, how, fmt.Errorf("%s is neither an object nor a world in this store", ref)
-	}
-	e, ok, err := r.entry(w.Revision)
-	if err != nil {
-		return e, how, err
-	}
-	if !ok || e.EntryHash != w.LogHead {
-		return e, how, fmt.Errorf("world %s names log head %s, which is not entry %d", ref, w.LogHead, w.Revision)
-	}
-	return e, how, nil
-}
-
-// scanFor scans back from head for the first entry match accepts.
-func (r worldReader) scanFor(scan int, how string, match func(wireEntry, *invRecord) bool) (wireEntry, string, error) {
-	h, ok, err := r.headIndex()
-	if err != nil {
-		return wireEntry{}, how, err
-	}
-	if !ok {
-		return wireEntry{}, how, errScanExhausted
-	}
-	var found *wireEntry
-	_, err = r.scanEntries(0, h, true, scan, func(e wireEntry, rec *invRecord) (bool, error) {
-		if match(e, rec) {
-			found = &e
-			return true, nil
-		}
-		return false, nil
-	})
-	if err != nil {
-		return wireEntry{}, how, err
-	}
-	if found == nil {
-		return wireEntry{}, how, errScanExhausted
-	}
-	return *found, how, nil
-}
-
-func contains(xs []string, x string) bool {
-	for _, v := range xs {
-		if v == x {
-			return true
-		}
-	}
-	return false
-}
-
-// ---------------------------------------------------------------------------
-// the chain
-// ---------------------------------------------------------------------------
-
-type whyLink struct {
-	Name   string `json:"name"`
-	OK     bool   `json:"ok"`
-	Ref    string `json:"ref,omitempty"`
-	Detail string `json:"detail"`
-	// Worktree (row 184 D3) is an effect result's worktree_delta: the files
-	// the toolchain call created, modified or deleted in the worktree.
-	Worktree   *broker.WorktreeDelta `json:"worktree,omitempty"`
-	Protected  json.RawMessage       `json:"protected,omitempty"`
-	ExecResult bool                  `json:"-"`
-}
-
-// Keep the complete stored block for JSON; nil alone means the key is absent.
-func resultProtected(payload []byte) json.RawMessage {
-	var res struct {
-		Protected json.RawMessage `json:"protected"`
-	}
-	if json.Unmarshal(payload, &res) != nil {
-		return nil
-	}
-	return res.Protected
-}
-
-// whyProtected mirrors broker's private protectedRecord for text decoding.
-type whyProtected struct {
-	Digest            string   `json:"digest"`
-	Patterns          []string `json:"patterns"`
-	Enforced          string   `json:"enforced"`
-	Base              string   `json:"base"`
-	BaseDigest        string   `json:"base_digest"`
-	BaseFiles         int      `json:"base_files"`
-	BaseIgnoredDigest string   `json:"base_ignored_digest"`
-	BaseIgnoredFiles  int      `json:"base_ignored_files"`
-	IntactBefore      *bool    `json:"intact_before"`
-	IntactAfter       *bool    `json:"intact_after"`
-	Changes           []struct {
-		Path          string `json:"path"`
-		Status        string `json:"status"`
-		Phase         string `json:"phase"`
-		PathTruncated bool   `json:"path_truncated,omitempty"`
-	} `json:"changes"`
-	Truncated bool   `json:"truncated"`
-	Total     int    `json:"total"`
-	Error     string `json:"error"`
-}
-
-func protectedLine(raw json.RawMessage, exec bool) string {
-	if raw == nil {
-		if exec {
-			return "protected: none declared"
-		}
-		return ""
-	}
-	var p *whyProtected
-	if json.Unmarshal(raw, &p) != nil || p == nil {
-		return "protected: not checked — malformed protected block"
-	}
-	var phases []string
-	if p.IntactBefore != nil && !*p.IntactBefore {
-		phases = append(phases, "before")
-	}
-	if p.IntactAfter != nil && !*p.IntactAfter {
-		phases = append(phases, "after")
-	}
-	if len(phases) > 0 {
-		var paths []string
-		for _, c := range p.Changes {
-			if len(paths) == whyWorktreePaths {
-				break
-			}
-			paths = append(paths, fmt.Sprintf("%s (%s)", c.Path, c.Status))
-		}
-		line := "protected: CHANGED " + strings.Join(phases, " and ") + " the call — " + strings.Join(paths, ", ")
-		if p.Total > len(paths) {
-			line += fmt.Sprintf(" (+%d more)", p.Total-len(paths))
-		}
-		if p.Error != "" {
-			line += "; not checked — " + p.Error
-		}
-		return line
-	}
-	if p.IntactBefore == nil || p.IntactAfter == nil || p.Error != "" {
-		reason := p.Error
-		if reason == "" {
-			reason = "no verdict recorded"
-		}
-		return "protected: not checked — " + reason
-	}
-	return fmt.Sprintf("protected: intact (%d files @ %.7s, digest %s, enforced %s)", p.BaseFiles, p.Base, p.BaseDigest, p.Enforced)
-}
-
-// whyWorktreePaths is how many changed paths the text form prints.
-const whyWorktreePaths = 5
-
-// resultWorktreeDelta decodes an effect result's worktree_delta, or nil when
-// the result has none (or is not a JSON object).
-func resultWorktreeDelta(payload []byte) *broker.WorktreeDelta {
-	var res struct {
-		Delta *broker.WorktreeDelta `json:"worktree_delta"`
-	}
-	if json.Unmarshal(payload, &res) != nil {
-		return nil
-	}
-	return res.Delta
-}
-
-// worktreeLine renders a delta as `worktree: N created, M modified, K
-// deleted` plus the first paths (D3), or the snapshot's error or cap.
-func worktreeLine(d *broker.WorktreeDelta) string {
-	switch {
-	case d.Error != "":
-		return "worktree: not recorded — " + d.Error
-	case d.Basis == broker.WorktreeBasisWalkCapped:
-		return fmt.Sprintf("worktree: not recorded — the non-git worktree has more than %d files (basis %s)", d.CappedAt, d.Basis)
-	}
-	// Counts are of the listed changes; a truncated delta says so.
-	var created, modified, deleted int
-	paths := make([]string, 0, whyWorktreePaths)
-	for _, c := range d.Changes {
-		switch c.Op {
-		case "created":
-			created++
-		case "modified":
-			modified++
-		case "deleted":
-			deleted++
-		}
-		if len(paths) < whyWorktreePaths {
-			paths = append(paths, c.Op[:1]+" "+c.Path)
-		}
-	}
-	line := fmt.Sprintf("worktree: %d created, %d modified, %d deleted (basis %s)", created, modified, deleted, d.Basis)
-	if d.Truncated {
-		line += fmt.Sprintf(" — first %d of %d changes listed", len(d.Changes), d.Total)
-	}
-	if len(paths) > 0 {
-		line += ": " + strings.Join(paths, ", ")
-		if more := d.Total - len(paths); more > 0 {
-			line += fmt.Sprintf(" (+%d more)", more)
-		}
-	}
-	return line
-}
-
-type whyChain struct {
-	Entry       int64     `json:"entry"`
-	MatchedBy   string    `json:"matchedBy"`
-	Coordinator bool      `json:"coordinator"`
-	Note        string    `json:"note,omitempty"`
-	Links       []whyLink `json:"links"`
-	OK          bool      `json:"ok"`
-}
-
-func (ch *whyChain) add(name string, ok bool, ref, detail string, args ...any) {
-	ch.Links = append(ch.Links, whyLink{Name: name, OK: ok, Ref: ref, Detail: fmt.Sprintf(detail, args...)})
-	if !ok {
-		ch.OK = false
-	}
-}
-
-// entryWire / worldWire mirror host/coordinator/plan.go's canonical encodings
-// (field order is the encoding order), so the entry hash and the world ref
-// are recomputed, not trusted.
-type entryWire struct {
-	EntryIndex     int64  `json:"entryIndex"`
-	SemanticsEpoch int64  `json:"semanticsEpoch"`
-	TransitionFn   string `json:"transitionFn"`
-	Interpreter    string `json:"interpreter"`
-	PrevEntryHash  string `json:"prevEntryHash"`
-	WrittenBy      string `json:"writtenBy"`
-	TransitionRef  string `json:"transitionRef"`
-}
-
-type worldWire struct {
-	Revision  int64  `json:"revision"`
-	StateRoot string `json:"stateRoot"`
-	LogHead   string `json:"logHead"`
-}
-
-func sumJSON(v any) string {
-	b, _ := json.Marshal(v)
-	return hashref.SumSHA256(b).String()
-}
-
-func sum(b []byte) string { return hashref.SumSHA256(b).String() }
-
-// checkedObject fetches ref and reports whether it resolves, carries the
-// wanted semantic id, and its payload hashes to ref.
-func (r worldReader) checkedObject(ref, semantic string) (wireObject, bool, string) {
-	obj, ok, err := r.object(ref)
-	switch {
-	case err != nil:
-		return obj, false, err.Error()
-	case !ok:
-		return obj, false, "does not resolve"
-	case semantic != "" && obj.SemanticID != semantic:
-		return obj, false, fmt.Sprintf("semanticId %s, want %s", obj.SemanticID, semantic)
-	case sum(obj.Payload) != ref:
-		return obj, false, fmt.Sprintf("payload hashes to %s, not its ref", sum(obj.Payload))
-	}
-	return obj, true, ""
-}
-
-func failNote(ok bool, why string) string {
-	if ok {
-		return ""
-	}
-	return " — ✗ " + why
-}
-
-// buildChain renders and checks one entry's provenance chain. result is the
-// caller's result bytes (why - / --result), or nil.
-func (r worldReader) buildChain(e wireEntry, result []byte) (whyChain, error) {
-	ch := whyChain{Entry: e.Header.EntryIndex, OK: true}
-	h, err := r.health()
-	if err != nil {
-		return ch, err
-	}
-	pin := "not the daemon's pinned interpreter"
-	if e.Header.Interpreter == h.InterpreterRef && h.InterpreterRef != "" {
-		pin = strings.TrimSpace(strings.SplitN(h.InterpreterVersion, "\n", 2)[0]) + ", the daemon's pin"
-	}
-	rec, recObj, err := r.record(e)
-	if err != nil {
-		return ch, err
-	}
-	entryLine := fmt.Sprintf("#%d writtenBy %s prev %s transitionFn %s interpreter %s (%s)",
-		e.Header.EntryIndex, e.Header.WrittenBy, e.Header.PrevEntryHash, e.Header.TransitionFn, e.Header.Interpreter, pin)
-	if rec == nil {
-		ch.Note = fmt.Sprintf("entry %d is not a coordinator invocation (writtenBy %s; its transitionRef is a %s): "+
-			"there is no record chain to walk, so only the entry and its object are shown", e.Header.EntryIndex, e.Header.WrittenBy, recObj.SemanticID)
-		ch.add("entry", true, e.EntryHash, "%s (a client-supplied entry hash: not recomputed)", entryLine)
-		ok := sum(recObj.Payload) == e.TransitionRef
-		ch.add("object", ok, e.TransitionRef, "%s, %d bytes%s", recObj.SemanticID, len(recObj.Payload), failNote(ok, "payload does not hash to transitionRef"))
-		return ch, nil
-	}
-	ch.Coordinator = true
-
-	// world: recomputed from the entry and the output, confirmed by the route.
-	wref := sumJSON(worldWire{Revision: e.Header.EntryIndex, StateRoot: rec.Output, LogHead: e.EntryHash})
-	w, wok, werr := r.world(wref)
-	worldOK := werr == nil && wok && w.Revision == e.Header.EntryIndex && w.StateRoot == rec.Output && w.LogHead == e.EntryHash
-	why := "GET /v1/worlds/<recomputed ref> does not return this world"
-	if werr != nil {
-		why = werr.Error()
-	}
-	ch.add("world", worldOK, wref, "revision %d, stateRoot %s, logHead %s (recomputed)%s", e.Header.EntryIndex, rec.Output, e.EntryHash, failNote(worldOK, why))
-
-	// entry: the hash is recomputed over the canonical header + transitionRef.
-	eh := sumJSON(entryWire{EntryIndex: e.Header.EntryIndex, SemanticsEpoch: e.Header.SemanticsEpoch,
-		TransitionFn: e.Header.TransitionFn, Interpreter: e.Header.Interpreter, PrevEntryHash: e.Header.PrevEntryHash,
-		WrittenBy: e.Header.WrittenBy, TransitionRef: e.TransitionRef})
-	ch.add("entry", eh == e.EntryHash, e.EntryHash, "%s%s", entryLine, failNote(eh == e.EntryHash, "entry hash recomputes to "+eh))
-
-	// record: its hash, and its fields agree with the entry.
-	recOK := sum(rec.payload) == e.TransitionRef && rec.TransitionFn == e.Header.TransitionFn &&
-		rec.Interpreter == e.Header.Interpreter && rec.SemanticsEpoch == e.Header.SemanticsEpoch
-	ch.add("record", recOK, e.TransitionRef, "%s invocation %s episode %s skill %s%s", rec.SemanticID, rec.InvocationID, rec.EpisodeID, rec.SkillID,
-		failNote(recOK, "the record's hash or its transitionFn/interpreter/epoch disagree with the entry"))
-
-	// input.
-	in, inOK, inWhy := r.checkedObject(rec.Input, coordinator.InputV1)
-	ch.add("input", inOK, rec.Input, "%s%s", whyInputLine(in.Payload), failNote(inOK, inWhy))
-
-	// plan and effects (v2 only).
-	var planned []struct {
-		ID     string `json:"id"`
-		Effect string `json:"effect"`
-		Scope  string `json:"scope"`
-		Cost   int64  `json:"cost"`
-	}
-	if rec.SemanticID == coordinator.RecordV2 {
-		po, pOK, pWhy := r.checkedObject(rec.Plan, coordinator.EffectPlanV1)
-		var plan struct {
-			Effects []struct {
-				ID     string `json:"id"`
-				Effect string `json:"effect"`
-				Scope  string `json:"scope"`
-				Cost   int64  `json:"cost"`
-			} `json:"effects"`
-			Result json.RawMessage `json:"result"`
-		}
-		detail := "(undecodable)"
-		if pOK && json.Unmarshal(po.Payload, &plan) == nil {
-			planned = plan.Effects
-			if len(plan.Effects) == 0 {
-				detail = "no effects; result " + elide(string(plan.Result))
-			} else {
-				parts := make([]string, 0, len(plan.Effects))
-				for _, pe := range plan.Effects {
-					parts = append(parts, fmt.Sprintf("%s %s@%s cost %d", pe.ID, pe.Effect, pe.Scope, pe.Cost))
-				}
-				detail = fmt.Sprintf("%d effect(s): %s", len(plan.Effects), strings.Join(parts, "; "))
-			}
-		}
-		ch.add("plan", pOK, rec.Plan, "%s%s", detail, failNote(pOK, pWhy))
-		for i, ref := range rec.Effects {
-			obj, ok, oWhy := r.checkedObject(ref, broker.EffectRecordV1)
-			id := fmt.Sprintf("#%d", i)
-			if i < len(planned) {
-				id = planned[i].ID
-			}
-			if !ok {
-				ch.add("effect", false, ref, "%s — ✗ %s", id, oWhy)
-				continue
-			}
-			d, err := broker.DecodeRecord(obj.Payload)
-			if err != nil {
-				ch.add("effect", false, ref, "%s — ✗ %v", id, err)
-				continue
-			}
-			verdict := "allowed"
-			switch {
-			case !d.Allowed:
-				verdict = "denied:" + d.Denial
-			case d.Failed:
-				verdict = "failed"
-			}
-			okAll := broker.RecordConsistent(d)
-			req, reqOK, _ := r.checkedObject(d.RequestRef.String(), broker.EffectRequestV1)
-			okAll = okAll && reqOK
-			resText := "no result"
-			if d.Failure != "" {
-				// Row 255: the failed record's bounded failure line.
-				resText = "no result: " + d.Failure
-			}
-			var delta *broker.WorktreeDelta
-			var protected json.RawMessage
-			execResult := false
-			if !d.ResultRef.IsZero() {
-				res, resOK, _ := r.checkedObject(d.ResultRef.String(), broker.EffectResultV1)
-				okAll = okAll && resOK
-				resText = fmt.Sprintf("result %d B", len(res.Payload))
-				if resOK {
-					delta = resultWorktreeDelta(res.Payload)
-					protected = resultProtected(res.Payload)
-					execResult = d.Effect == broker.EffectWorkspaceExec
-				}
-			}
-			ch.add("effect", okAll, ref, "%s %s %s@%s cost %d budget %d→%d request %d B %s%s", id, verdict, d.Effect, d.Scope, d.Cost,
-				d.BudgetBefore, d.BudgetAfter, len(req.Payload), resText,
-				failNote(okAll, "the record is inconsistent or its request/result does not resolve to its hash"))
-			ch.Links[len(ch.Links)-1].Worktree = delta
-			ch.Links[len(ch.Links)-1].Protected = protected
-			ch.Links[len(ch.Links)-1].ExecResult = execResult
-		}
-	}
-
-	// output: hashes to record.output, equals world.stateRoot, and its world
-	// block names this record's plan and effect records.
-	out, oOK, oWhy := r.checkedObject(rec.Output, coordinator.OutputV1)
-	if oOK && !(worldOK && w.StateRoot == rec.Output) {
-		oOK, oWhy = false, "the output is not the committed world's stateRoot"
-	}
-	if oOK && rec.SemanticID == coordinator.RecordV2 {
-		var ob struct {
-			World struct {
-				Plan    string `json:"plan"`
-				Effects []struct {
-					Record *string `json:"record"`
-				} `json:"effects"`
-			} `json:"world"`
-		}
-		_ = json.Unmarshal(out.Payload, &ob)
-		var recs []string
-		for _, fx := range ob.World.Effects {
-			if fx.Record != nil {
-				recs = append(recs, *fx.Record)
-			}
-		}
-		if ob.World.Plan != rec.Plan || strings.Join(recs, ",") != strings.Join(rec.Effects, ",") {
-			oOK, oWhy = false, "the output's world block does not name this record's plan and effect records"
-		}
-	}
-	detail := fmt.Sprintf("%d B = world.stateRoot", len(out.Payload))
-	if result != nil {
-		if sum(result) == rec.Output {
-			detail += "; your result is these exact bytes"
-		} else {
-			oOK, oWhy = false, "your result ("+sum(result)+") is not the committed output bytes"
-		}
-	}
-	ch.add("output", oOK, rec.Output, "%s%s", detail, failNote(oOK, oWhy))
-	return ch, nil
-}
-
-func printChain(w io.Writer, ch whyChain) {
-	fmt.Fprintf(w, "why: entry %d (matched by %s)\n", ch.Entry, ch.MatchedBy)
-	if ch.Note != "" {
-		fmt.Fprintf(w, "  note: %s\n", ch.Note)
-	}
-	bad := 0
-	for _, l := range ch.Links {
-		mark := "✓"
-		if !l.OK {
-			mark, bad = "✗", bad+1
-		}
-		ref := ""
-		if l.Ref != "" {
-			ref = l.Ref + "  "
-		}
-		fmt.Fprintf(w, "  %s %-7s %s%s\n", mark, l.Name, ref, l.Detail)
-		if l.Worktree != nil {
-			fmt.Fprintf(w, "            %s\n", worktreeLine(l.Worktree))
-		}
-		if line := protectedLine(l.Protected, l.ExecResult); line != "" {
-			fmt.Fprintf(w, "            %s\n", line)
-		}
-	}
-	if bad == 0 {
-		fmt.Fprintf(w, "all %d link(s) verified\n", len(ch.Links))
-	} else {
-		fmt.Fprintf(w, "%d of %d link(s) BROKEN — integrity refusal (exit 3)\n", bad, len(ch.Links))
-	}
-}
-
-// whyKeyFieldMax bounds one key field's encoded value; whyKeyFieldCount the
-// number of key fields shown.
-const (
-	whyKeyFieldMax   = 120
-	whyKeyFieldCount = 6
-)
-
-// whyInputLine renders an invocation input for `why`. A payload that fits in
-// elideAt is printed whole, unchanged. A longer one is elided — and because
-// canonical JSON orders keys, the field that names WHAT a call touched can
-// sit past the cut (ailang-edit's `path` follows `new_text`; row 114 F-7, a
-// walk missed an edited file this way). So the short top-level fields of an
-// object input — each scalar or small array whose encoding is at most
-// whyKeyFieldMax bytes, in key order, at most whyKeyFieldCount — are printed
-// first, then the elided payload.
-func whyInputLine(payload []byte) string {
-	one := strings.ReplaceAll(string(payload), "\n", `\n`)
-	if len(one) <= elideAt {
-		return one
-	}
-	var obj map[string]json.RawMessage
-	if json.Unmarshal(payload, &obj) != nil || len(obj) == 0 {
-		return elide(string(payload))
-	}
-	keys := make([]string, 0, len(obj))
-	for k := range obj {
-		keys = append(keys, k)
-	}
-	sort.Strings(keys)
-	var fields []string
-	for _, k := range keys {
-		v := obj[k]
-		if len(v) == 0 || len(v) > whyKeyFieldMax || v[0] == '{' {
-			continue
-		}
-		fields = append(fields, k+"="+strings.ReplaceAll(string(v), "\n", `\n`))
-		if len(fields) == whyKeyFieldCount {
-			break
-		}
-	}
-	if len(fields) == 0 {
-		return elide(string(payload))
-	}
-	return strings.Join(fields, " ") + " · " + elide(string(payload))
-}
diff --git a/cmd/ailang-worldd/why_cli.go b/cmd/ailang-worldd/why_cli.go
new file mode 100644
--- /dev/null
+++ b/cmd/ailang-worldd/why_cli.go
@@ -0,0 +1,241 @@
+package main
+
+// Row 243 V2a (design_docs/planned/w-world-address-bar-verbs-v2.md §2.1): why's CLI front
+// (help, flags, the --result and - reads, --expect-store) over host/why's walk, then the
+// names the other CLI files use, kept so that none of them is edited.
+
+import (
+	"context"
+	"errors"
+	"flag"
+	"fmt"
+	"io"
+	"net/http"
+	"os"
+	"strings"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/why"
+)
+
+const (
+	whyBudget      = 60 * time.Second
+	whyScanDefault = why.ScanDefault
+	whyScanMax     = why.ScanMax
+	logPageMax     = why.LogPageMax
+	frozenRigPort  = 47_330 // historical trailer recovery, never a default
+)
+
+const whyHelp = `usage: ailang-worldd [--addr <url>] why <target> [--scan N] [--json] [--expect-store <path>]
+       ailang-worldd [--addr <url>] why --result <file> [--scan N] [--json]
+
+Walks one committed invocation back to its whole provenance chain, checking
+every link (✓ or ✗):
+
+  world    the world ref, recomputed from the entry and the output, and
+           confirmed by GET /v1/worlds/<ref>
+  entry    index, writtenBy, prev, transitionFn and the interpreter (with its
+           pin name); the entry hash is recomputed
+  record   invocation id, episode, skill (its hash and its fields agree with
+           the entry)
+  input    your arguments, decoded and elided; when elided, the short
+           top-level fields (path, command, op, …) are printed first
+  plan     the effects planned, or the refusal result
+  effect   per effect record: allowed/failed/denial, budget before→after,
+           request and result sizes; a toolchain call that changed the
+           worktree adds "worktree: N created, M modified, K deleted" and
+           its first paths (--json: the whole worktree_delta)
+  output   the committed output; it must equal world.stateRoot
+
+Targets:
+  <index> | head               a log entry
+  sha256:<ref>                 any object (dispatched on its semanticId: a
+                               record, input, output, plan, effect record,
+                               effect request or result), or a world ref
+  a2a:<id> | mcp:<id> | rest:<id>
+                               an invocation id
+  - | --result <file>          the output of 'call --json-out': its sha256 is
+                               the output ref (the world.plan is the fallback)
+
+A target that is not an index is found by scanning the log backwards from
+head: --scan N entries (default 500, max 5000) within a 60 s budget; past
+either it prints "not found in the last N entries" and exits 1.
+
+Exit: 0 every link ✓; 3 any link ✗ (an integrity refusal); 1 not found or
+usage. --json prints the chain as JSON.
+`
+
+// parseWhyTarget makes the CLI's two reads (--result's file, -'s stdin) and
+// hands every other target to why.ParseTarget.
+func parseWhyTarget(arg, resultFile string, stdin io.Reader) (why.Target, error) {
+	if resultFile != "" {
+		if arg != "" {
+			return why.Target{}, errors.New("give either a target or --result, not both")
+		}
+		b, err := os.ReadFile(resultFile)
+		if err != nil {
+			return why.Target{}, fmt.Errorf("read --result %s: %v", resultFile, errors.Unwrap(err))
+		}
+		return why.Target{Kind: "result", Result: trimOneNewline(b)}, nil
+	}
+	if arg == "-" {
+		b, err := io.ReadAll(io.LimitReader(stdin, maxClientResponseBytes+1))
+		if err != nil {
+			return why.Target{}, fmt.Errorf("read stdin: %v", err)
+		}
+		return why.Target{Kind: "result", Result: trimOneNewline(b)}, nil
+	}
+	return why.ParseTarget(arg)
+}
+
+// trimOneNewline drops the one newline `call --json-out` appends.
+func trimOneNewline(b []byte) []byte {
+	if n := len(b); n > 0 && b[n-1] == '\n' {
+		return b[:n-1]
+	}
+	return b
+}
+
+func runWhy(addr string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
+	target := ""
+	if len(args) > 0 && (!strings.HasPrefix(args[0], "-") || args[0] == "-") {
+		target, args = args[0], args[1:]
+	}
+	fs := flag.NewFlagSet("ailang-worldd why", flag.ContinueOnError)
+	fs.SetOutput(stderr)
+	fs.Usage = func() {}
+	expectStore := fs.String("expect-store", "", "refuse unless the daemon serves this store (World-Provenance store= path)")
+	resultFile := fs.String("result", "", "a 'call --json-out' output file")
+	scan := fs.Int("scan", whyScanDefault, "entries to scan back from head (max 5000)")
+	asJSON := fs.Bool("json", false, "print the chain as JSON")
+	if code, done := parseVerbFlags(fs, args, whyHelp, stdout, stderr); done {
+		return code
+	}
+	if fs.NArg() > 0 {
+		if target != "" || fs.NArg() > 1 {
+			fmt.Fprintf(stderr, "ailang-worldd why: unexpected argument %q\n", fs.Arg(fs.NArg()-1))
+			return exitUsage
+		}
+		target = fs.Arg(0)
+	}
+	if *scan < 1 || *scan > whyScanMax {
+		fmt.Fprintf(stderr, "ailang-worldd why: --scan must be 1..%d\n", whyScanMax)
+		return exitUsage
+	}
+	t, err := parseWhyTarget(target, *resultFile, stdin)
+	if err != nil {
+		fmt.Fprintf(stderr, "ailang-worldd why: %v\n%s", err, whyHelp)
+		return exitUsage
+	}
+	ctx, cancel := budgetContext(whyBudget)
+	defer cancel()
+	r := worldReader{ctx: ctx, c: newClient(addr)}
+	if *expectStore != "" {
+		h, err := r.health()
+		if err != nil {
+			fmt.Fprintf(stderr, "ailang-worldd why: %v\n", err)
+			return exitUsage
+		}
+		r.healthCache = &h
+		if !sameStorePath(h.DBPath, *expectStore) {
+			return serviceMessage(stderr,
+				fmt.Sprintf("ailang-worldd why: the daemon at %s serves store %s, not %s", addr, h.DBPath, *expectStore),
+				"entry indexes are per store, so the same index in another store is another entry (a World-Provenance trailer names its store as store=<path>)",
+				fmt.Sprintf("pass the --addr of the daemon serving %s (the frozen rig: start it from daemon-cmdline.txt on 127.0.0.1:%d), or drop --expect-store", *expectStore, frozenRigPort), exitUsage)
+		}
+	}
+	ans := why.Run(r.why(), t, why.Options{Scan: *scan, JSON: *asJSON})
+	fmt.Fprint(stdout, ans.Stdout)
+	fmt.Fprint(stderr, ans.Stderr)
+	return ans.Exit
+}
+
+// The compatibility layer: the names other CLI files use, on host/why.
+
+type (
+	wireHeader = why.Header
+	wireEntry  = why.Entry
+	wireObject = why.Object
+	wireWorld  = why.World
+	invRecord  = why.Record
+	healthInfo = why.Health
+	whyChain   = why.Chain
+	whyLink    = why.Link
+)
+
+// entryWire and worldWire are defined (not aliased) because demo_seed.go and
+// demo_test.go build them positionally, which vet refuses for another
+// package's struct; the encoding stays host/why's.
+type (
+	entryWire why.EntryWire
+	worldWire why.WorldWire
+)
+
+func sumJSON(v any) string { return why.SumJSON(v) }
+func sum(b []byte) string  { return why.Sum(b) }
+
+func contains(xs []string, x string) bool {
+	for _, v := range xs {
+		if v == x {
+			return true
+		}
+	}
+	return false
+}
+
+// worldReader is the GET side of the daemon, under one command budget.
+type worldReader struct {
+	ctx         context.Context
+	c           *client
+	healthCache *healthInfo
+}
+
+// why is r as host/why's Reader, reading through the CLI's HTTP client.
+func (r worldReader) why() why.Reader {
+	return why.Reader{Ctx: r.ctx, Get: clientGetter(r.c), HealthCache: r.healthCache}
+}
+
+// clientGetter is c's GET as a why.Getter (c.get, cli.go, is another shape).
+func clientGetter(c *client) why.Getter {
+	return func(ctx context.Context, path string) (int, []byte, error) {
+		return c.do(ctx, http.MethodGet, path, nil)
+	}
+}
+
+type whyTarget struct { // a target as demo_seed builds it (kind index)
+	kind   string // index, head, ref, invocation, result
+	index  int64
+	text   string
+	result []byte
+}
+
+// The reader's methods other CLI files call, each host/why's.
+func (r worldReader) getJSON(path string, v any) (bool, error) { return r.why().GetJSON(path, v) }
+func (r worldReader) entry(i int64) (wireEntry, bool, error)   { return r.why().Entry(i) }
+func (r worldReader) page(from int64, n int) ([]wireEntry, error) {
+	return r.why().Page(from, n)
+}
+func (r worldReader) object(ref string) (wireObject, bool, error) { return r.why().Object(ref) }
+func (r worldReader) world(ref string) (wireWorld, bool, error)   { return r.why().World(ref) }
+func (r worldReader) headIndex() (int64, bool, error)             { return r.why().HeadIndex() }
+func (r worldReader) health() (healthInfo, error)                 { return r.why().Health() }
+func (r worldReader) record(e wireEntry) (*invRecord, wireObject, error) {
+	return r.why().Record(e)
+}
+func (r worldReader) scanEntries(from, to int64, reverse bool, max int,
+	visit func(e wireEntry, rec *invRecord) (bool, error)) (int, error) {
+	return r.why().ScanEntries(from, to, reverse, max, visit)
+}
+func (r worldReader) locate(t whyTarget, scan int) (wireEntry, string, error) {
+	return r.why().Locate(why.Target{Kind: t.kind, Index: t.index, Text: t.text, Result: t.result}, scan)
+}
+func (r worldReader) effectRecord(ref string) (broker.EffectRecord, bool) {
+	return r.why().EffectRecord(ref)
+}
+func (r worldReader) checkedObject(ref, semantic string) (wireObject, bool, string) {
+	return r.why().CheckedObject(ref, semantic)
+}
+func (r worldReader) buildChain(e wireEntry, result []byte) (whyChain, error) {
+	return r.why().BuildChain(e, result)
+}
diff --git a/cmd/ailang-worldd/why_input_test.go b/cmd/ailang-worldd/why_input_test.go
deleted file mode 100644
--- a/cmd/ailang-worldd/why_input_test.go
+++ /dev/null
@@ -1,62 +0,0 @@
-package main
-
-import (
-	"strconv"
-	"strings"
-	"testing"
-)
-
-// Row 114 F-7: `why 198` on an ailang-edit showed `{"new_text":"// Row 185:
-// IsLinked accepts…… (1048 bytes)` — the edited file's path sat past the
-// cut, and a walk concluded the file had not been edited.
-func TestWhyInputLineShowsThePathOfALongEdit(t *testing.T) {
-	payload := []byte(`{"new_text":"` + strings.Repeat(`// Row 185: IsLinked accepts the worktree Add makes\n`, 20) +
-		`","old_text":"}\n","path":"host/worktree/worktree_test.go"}`)
-	got := whyInputLine(payload)
-	if !strings.HasPrefix(got, `old_text="}\n" path="host/worktree/worktree_test.go" · {"new_text":`) {
-		t.Fatalf("key fields not first:\n%s", got)
-	}
-	if !strings.HasSuffix(got, "("+strconv.Itoa(len(payload))+" bytes)") {
-		t.Fatalf("the elided payload and its size must still follow:\n%s", got)
-	}
-	if strings.Contains(got, "new_text=") {
-		t.Fatalf("a long field must not be repeated as a key field:\n%s", got)
-	}
-}
-
-// A payload that fits is printed whole and byte-identical to before (no key
-// fields), so short inputs — most of them — read exactly as they did.
-func TestWhyInputLineLeavesAShortInputUnchanged(t *testing.T) {
-	for _, p := range []string{`{"path":"cmd/ailang-worldd/setup.go"}`, `{"args":["-run","TestX","./cmd"],"command":"go-test"}`, `"x"`} {
-		if got := whyInputLine([]byte(p)); got != elide(p) || got != p {
-			t.Fatalf("short input changed: %q -> %q", p, got)
-		}
-	}
-}
-
-// Small arrays count as key fields (an exec call's args); nested objects and
-// over-long values do not; at most whyKeyFieldCount are shown.
-func TestWhyInputLineKeyFieldBounds(t *testing.T) {
-	long := strings.Repeat("a", whyKeyFieldMax)
-	payload := []byte(`{"a1":1,"a2":2,"a3":3,"a4":4,"a5":5,"a6":6,"a7":7,"args":["x"],"big":"` + long +
-		`","nested":{"k":"v"}}`)
-	got := whyInputLine(payload)
-	head, _, ok := strings.Cut(got, " · ")
-	if !ok || head != "a1=1 a2=2 a3=3 a4=4 a5=5 a6=6" {
-		t.Fatalf("key fields = %q (want the first %d in key order)", head, whyKeyFieldCount)
-	}
-	small := []byte(`{"args":["-run","TestSessionNew","./cmd/ailang-worldd"],"command":"go-test","note":"` + long + `"}`)
-	head, _, _ = strings.Cut(whyInputLine(small), " · ")
-	if head != `args=["-run","TestSessionNew","./cmd/ailang-worldd"] command="go-test"` {
-		t.Fatalf("array/scalar key fields = %q", head)
-	}
-}
-
-// Not an object, or an object with no short field: plain elide, as before.
-func TestWhyInputLineFallsBackToElide(t *testing.T) {
-	for _, p := range []string{`"` + strings.Repeat("z", 300) + `"`, `{"blob":"` + strings.Repeat("z", 300) + `"}`, strings.Repeat("not json ", 40)} {
-		if got := whyInputLine([]byte(p)); got != elide(p) {
-			t.Fatalf("fallback changed: %q", got)
-		}
-	}
-}
diff --git a/cmd/ailang-worldd/why_protected_test.go b/cmd/ailang-worldd/why_protected_test.go
--- a/cmd/ailang-worldd/why_protected_test.go
+++ b/cmd/ailang-worldd/why_protected_test.go
@@ -1,7 +1,6 @@
 package main
 
 import (
-	"bytes"
 	"encoding/json"
 	"fmt"
 	"reflect"
@@ -14,40 +13,8 @@ import (
 
 const protectedFixture = `{"digest":"sha256:patterns","patterns":["*_test.go","tests/**"],"enforced":"dirs-only","base":"123456789abcdef","base_digest":"sha256:tree","base_files":12,"base_ignored_digest":"sha256:ignored","base_ignored_files":3,"intact_before":null,"intact_after":null,"changes":[{"path":"tests/a","status":"untracked","phase":"after","path_truncated":true}],"truncated":true,"total":40,"error":"snapshot failed"}`
 
+// The protected line end to end (its pure cases are host/why's TestWhyProtectedLine).
 func TestWhyProtectedLine(t *testing.T) {
-	for _, tc := range []struct {
-		name, block, want string
-		exec              bool
-	}{
-		{"intact", `{"intact_before":true,"intact_after":true,"base_files":12,"base":"123456789","base_digest":"sha256:tree","enforced":"glob"}`, "protected: intact (12 files @ 1234567, digest sha256:tree, enforced glob)", true},
-		{"changed-before", `{"intact_before":false,"intact_after":true,"changes":[{"path":"a_test.go","status":"modified","phase":"before"}],"total":1}`, "protected: CHANGED before the call — a_test.go (modified)", true},
-		{"changed-after", `{"intact_before":true,"intact_after":false,"changes":[{"path":"tests/a","status":"untracked","phase":"after"}],"total":1}`, "protected: CHANGED after the call — tests/a (untracked)", true},
-		{"changed-both", `{"intact_before":false,"intact_after":false,"changes":[{"path":"a","status":"deleted","phase":"before"},{"path":"b","status":"ignored","phase":"after"}],"total":2}`, "protected: CHANGED before and after the call — a (deleted), b (ignored)", true},
-		{"index-flag", `{"intact_before":false,"changes":[{"path":"a_test.go","status":"index-flag"}],"total":1}`, "protected: CHANGED before the call — a_test.go (index-flag)", true},
-		{"head-moved", `{"intact_after":false,"changes":[{"path":"HEAD","status":"head-moved"}],"total":1}`, "protected: CHANGED after the call — HEAD (head-moved)", true},
-		{"truncated", `{"intact_after":false,"changes":[{"path":"a","status":"modified"},{"path":"b","status":"deleted"},{"path":"c","status":"untracked"},{"path":"d","status":"ignored"},{"path":"e","status":"modified"},{"path":"f","status":"deleted"}],"truncated":true,"total":40}`, "protected: CHANGED after the call — a (modified), b (deleted), c (untracked), d (ignored), e (modified) (+35 more)", true},
-		{"error", protectedFixture, "protected: not checked — snapshot failed", true},
-		{"error-beside-changed", `{"intact_after":false,"changes":[{"path":"a","status":"modified"}],"total":1,"error":"before failed"}`, "protected: CHANGED after the call — a (modified); not checked — before failed", true},
-		{"error-beside-true", `{"intact_before":true,"intact_after":true,"error":"failed"}`, "protected: not checked — failed", true},
-		{"missing-verdict", `{"intact_before":true}`, "protected: not checked — no verdict recorded", true},
-		{"missing-before", `{"intact_after":true,"base_files":12,"base":"123456789","base_digest":"sha256:tree","enforced":"glob"}`, "protected: not checked — no verdict recorded", true},
-		{"malformed", `"broken"`, "protected: not checked — malformed protected block", true},
-		{"malformed-null", `null`, "protected: not checked — malformed protected block", true},
-		{"malformed-field", `{"intact_before":"true"}`, "protected: not checked — malformed protected block", true},
-		{"none-declared", "", "protected: none declared", true},
-		{"non-exec", "", "", false},
-		{"old-record", "", "protected: none declared", true},
-	} {
-		t.Run(tc.name, func(t *testing.T) {
-			payload := []byte(`{"content":"old"}`)
-			if tc.block != "" {
-				payload = []byte(`{"protected":` + tc.block + `}`)
-			}
-			if got := protectedLine(resultProtected(payload), tc.exec); got != tc.want {
-				t.Fatalf("protected line = %q, want %q", got, tc.want)
-			}
-		})
-	}
 	t.Run("e2e-exec", func(t *testing.T) {
 		r := newWalkRig(t) // Every REST fixture owns its bounded context in this closure.
 		for i, tc := range []struct{ effect, result, want string }{
@@ -134,34 +101,3 @@ func assertProtectedBlock(t *testing.T, got, want json.RawMessage) {
 		t.Fatalf("protected JSON = %s, want %s", got, want)
 	}
 }
-
-func TestWhyProtectedJSON(t *testing.T) {
-	raw := resultProtected([]byte(`{"protected":` + protectedFixture + `}`))
-	b := mustJSONT(t, whyLink{Name: "effect", OK: true, Detail: "result", Protected: raw})
-	var link map[string]json.RawMessage
-	if err := json.Unmarshal(b, &link); err != nil {
-		t.Fatal(err)
-	}
-	assertProtectedBlock(t, link["protected"], json.RawMessage(protectedFixture))
-	var block map[string]json.RawMessage
-	if err := json.Unmarshal(link["protected"], &block); err != nil {
-		t.Fatal(err)
-	}
-	for _, key := range []string{"intact_before", "intact_after"} {
-		if string(block[key]) != "null" {
-			t.Fatalf("%s = %s, want null", key, block[key])
-		}
-	}
-	old := whyLink{Name: "effect", OK: true, Ref: "hash", Detail: "result", Protected: resultProtected([]byte(`{"content":"old"}`)), ExecResult: true}
-	b = mustJSONT(t, old)
-	if want := `{"name":"effect","ok":true,"ref":"hash","detail":"result"}`; string(b) != want {
-		t.Fatalf("old link JSON = %s, want %s", b, want)
-	}
-	var decoded whyLink
-	if err := json.Unmarshal(b, &decoded); err != nil || decoded.Protected != nil {
-		t.Fatalf("old record: %+v, %v", decoded, err)
-	}
-	if bytes.Contains(b, []byte(`"protected"`)) {
-		t.Fatalf("old link carries protected: %s", b)
-	}
-}
diff --git a/host/why/why.go b/host/why/why.go
new file mode 100644
--- /dev/null
+++ b/host/why/why.go
@@ -0,0 +1,998 @@
+// Package why is the walk behind `ailang-worldd why`, moved from cmd/ailang-worldd (row 243 V2a, design §2.1):
+// it reads World only through an injected Getter, so the CLI (HTTP) and the daemon (in process) run one walk.
+package why
+
+// Row 138 M2 (design_docs/planned/w-worldd-developer-cli.md §3.3): `why`
+// walks a result, an object, an invocation id or an entry back to its whole
+// provenance chain, recomputing and checking every link.
+//
+// An MCP result carries no entry index, invocation id or record ref (V12), and
+// /v1/receipts refuses a2a:/mcp: ids (V13), so a target that is not an index is
+// found by scanning the log backwards from head — one log page of 500 at a
+// time and one object GET per entry — bounded by --scan (default 500, max
+// 5000) and a 60 s budget. The reads are unauthenticated GETs (V15); `why`
+// trusts that read surface (R-CLI-2) but verifies every content address it is
+// handed against the bytes it is handed.
+
+import (
+	"context"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"net/url"
+	"sort"
+	"strconv"
+	"strings"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+)
+
+const (
+	ScanDefault = 500
+	ScanMax     = 5000
+	LogPageMax  = 500
+	// MaxResponseBytes is the CLI client's cap on one response (maxClientResponseBytes).
+	MaxResponseBytes = 8388608
+	// The two statuses a read distinguishes (host/why imports no net/http).
+	statusOK       = 200
+	statusNotFound = 404
+	exitOK         = 0 // the CLI's exits: every link ✓
+	exitUsage      = 1 // usage, or not found
+	exitIntegrity  = 3 // any link ✗
+)
+
+// ---------------------------------------------------------------------------
+// the read surface
+// ---------------------------------------------------------------------------
+
+type Header struct {
+	EntryIndex     int64  `json:"entryIndex"`
+	SemanticsEpoch int64  `json:"semanticsEpoch"`
+	TransitionFn   string `json:"transitionFn"`
+	Interpreter    string `json:"interpreter"`
+	PrevEntryHash  string `json:"prevEntryHash"`
+	WrittenBy      string `json:"writtenBy"`
+}
+
+type Entry struct {
+	Header        Header `json:"header"`
+	EntryHash     string `json:"entryHash"`
+	TransitionRef string `json:"transitionRef"`
+}
+
+type Object struct {
+	Hash       string `json:"hash"`
+	SemanticID string `json:"semanticId"`
+	Payload    []byte `json:"payload"` // base64 on the wire
+}
+
+type World struct {
+	Ref       string `json:"ref"`
+	Revision  int64  `json:"revision"`
+	StateRoot string `json:"stateRoot"`
+	LogHead   string `json:"logHead"`
+}
+
+// Record is a world/invocation-record/v1 or /v2 payload.
+type Record struct {
+	SemanticID     string   `json:"-"`
+	InvocationID   string   `json:"invocationId"`
+	EpisodeID      string   `json:"episodeId"`
+	SkillID        string   `json:"skillId"`
+	TransitionFn   string   `json:"transitionFn"`
+	Interpreter    string   `json:"interpreter"`
+	SemanticsEpoch int64    `json:"semanticsEpoch"`
+	Input          string   `json:"input"`
+	Output         string   `json:"output"`
+	Plan           string   `json:"plan"`
+	Effects        []string `json:"effects"`
+	payload        []byte
+}
+
+// Getter reads one GET path of the daemon's read surface: its status and body. A 404 is an
+// answer, not an error; err is a transport or bound failure.
+type Getter func(ctx context.Context, path string) (status int, body []byte, err error)
+
+// Reader is the GET side of the daemon, under one command budget.
+type Reader struct {
+	Ctx         context.Context
+	Get         Getter
+	HealthCache *Health
+}
+
+// getJSON decodes a 200 into v; a 404 is (false, nil); anything else errors.
+func (r Reader) GetJSON(path string, v any) (bool, error) {
+	if err := r.Ctx.Err(); err != nil {
+		return false, err
+	}
+	status, body, err := r.Get(r.Ctx, path)
+	if err != nil {
+		return false, err
+	}
+	switch status {
+	case statusOK:
+		if err := json.Unmarshal(body, v); err != nil {
+			return false, fmt.Errorf("GET %s: malformed body: %v", path, err)
+		}
+		return true, nil
+	case statusNotFound:
+		return false, nil
+	default:
+		return false, fmt.Errorf("GET %s returned HTTP %d: %s", path, status, strings.TrimSpace(string(body)))
+	}
+}
+
+func (r Reader) Entry(i int64) (Entry, bool, error) {
+	var e Entry
+	ok, err := r.GetJSON("/v1/log/"+strconv.FormatInt(i, 10), &e)
+	return e, ok, err
+}
+
+func (r Reader) Page(from int64, limit int) ([]Entry, error) {
+	var p struct {
+		Items []Entry `json:"items"`
+	}
+	_, err := r.GetJSON(fmt.Sprintf("/v1/log?from=%d&limit=%d", from, limit), &p)
+	return p.Items, err
+}
+
+func (r Reader) Object(ref string) (Object, bool, error) {
+	var o Object
+	ok, err := r.GetJSON("/v1/objects/"+url.PathEscape(ref)+"?payload=true", &o)
+	return o, ok, err
+}
+
+func (r Reader) World(ref string) (World, bool, error) {
+	var w World
+	ok, err := r.GetJSON("/v1/worlds/"+url.PathEscape(ref), &w)
+	return w, ok, err
+}
+
+// headIndex is the selected world's revision (= its log head's entry index).
+func (r Reader) HeadIndex() (int64, bool, error) {
+	ref, ok, err := Head(r.Ctx, r.Get)
+	if err != nil || !ok {
+		return 0, ok, err
+	}
+	w, found, err := r.World(ref)
+	if err != nil {
+		return 0, false, err
+	}
+	if !found {
+		return 0, false, fmt.Errorf("the selected head %s has no world row", ref)
+	}
+	return w.Revision, true, nil
+}
+
+// Head reads GET /v1/head: the selected world ref, or present=false on
+// the 404 "no world head has been selected yet".
+func Head(ctx context.Context, get Getter) (ref string, present bool, err error) {
+	status, body, err := get(ctx, "/v1/head")
+	if err != nil {
+		return "", false, err
+	}
+	switch status {
+	case statusOK:
+		return strings.TrimSpace(string(body)), true, nil
+	case statusNotFound:
+		return "", false, nil
+	default:
+		return "", false, fmt.Errorf("GET /v1/head returned HTTP %d: %s", status, strings.TrimSpace(string(body)))
+	}
+}
+
+type Health struct {
+	DBPath             string `json:"db_path"`
+	InterpreterRef     string `json:"interpreter_ref"`
+	InterpreterVersion string `json:"interpreter_version"`
+}
+
+func (r Reader) Health() (Health, error) {
+	if r.HealthCache != nil {
+		return *r.HealthCache, nil
+	}
+	var h Health
+	_, err := r.GetJSON("/v1/health", &h)
+	return h, err
+}
+
+// record fetches an entry's transitionRef and decodes it when it is a
+// coordinator invocation record; rec is nil for any other entry.
+func (r Reader) Record(e Entry) (*Record, Object, error) {
+	obj, ok, err := r.Object(e.TransitionRef)
+	if err != nil || !ok {
+		if err == nil {
+			err = fmt.Errorf("entry %d: transitionRef %s does not resolve", e.Header.EntryIndex, e.TransitionRef)
+		}
+		return nil, obj, err
+	}
+	if obj.SemanticID != coordinator.RecordV1 && obj.SemanticID != coordinator.RecordV2 {
+		return nil, obj, nil
+	}
+	var rec Record
+	if err := json.Unmarshal(obj.Payload, &rec); err != nil {
+		return nil, obj, fmt.Errorf("entry %d: undecodable %s: %v", e.Header.EntryIndex, obj.SemanticID, err)
+	}
+	rec.SemanticID, rec.payload = obj.SemanticID, obj.Payload
+	return &rec, obj, nil
+}
+
+// errScanExhausted reports a bounded scan that ran out of entries or budget.
+var errScanExhausted = errors.New("scan exhausted")
+
+// scanEntries visits the entries [from, to] — backwards from `to` when
+// reverse — one page of up to 500 at a time, decoding each entry's record
+// (one object GET per entry). It stops at the first visit that returns
+// true, after max entries, or when the budget ends.
+func (r Reader) ScanEntries(from, to int64, reverse bool, max int,
+	visit func(e Entry, rec *Record) (bool, error)) (int, error) {
+	seen := 0
+	for lo, hi := from, to; lo <= hi && seen < max; {
+		pFrom := lo
+		if reverse {
+			pFrom = hi - LogPageMax + 1
+			if pFrom < lo {
+				pFrom = lo
+			}
+		}
+		n := hi - pFrom + 1
+		if !reverse && n > LogPageMax {
+			n = LogPageMax
+		}
+		items, err := r.Page(pFrom, int(n))
+		if err != nil {
+			return seen, err
+		}
+		if len(items) == 0 {
+			return seen, nil
+		}
+		for k := range items {
+			i := k
+			if reverse {
+				i = len(items) - 1 - k
+			}
+			if seen >= max {
+				break
+			}
+			seen++
+			rec, _, err := r.Record(items[i])
+			if err != nil {
+				return seen, err
+			}
+			stop, err := visit(items[i], rec)
+			if err != nil || stop {
+				return seen, err
+			}
+		}
+		if reverse {
+			hi = pFrom - 1
+		} else {
+			lo = pFrom + int64(len(items))
+		}
+	}
+	return seen, nil
+}
+
+// ---------------------------------------------------------------------------
+// ---------------------------------------------------------------------------
+// why
+// ---------------------------------------------------------------------------
+
+// Target is what a walk starts from. Kind is index, head, ref, invocation or
+// result; a result's bytes are read by the CLI (--result, -), never here.
+type Target struct {
+	Kind   string
+	Index  int64
+	Text   string
+	Result []byte
+}
+
+// ParseTarget parses a target argument other than - and --result (the two
+// reads the CLI makes at its boundary), with the CLI's error texts.
+func ParseTarget(arg string) (Target, error) {
+	switch {
+	case arg == "":
+		return Target{}, errors.New("a target is required")
+	case arg == "head":
+		return Target{Kind: "head"}, nil
+	case strings.HasPrefix(arg, "a2a:") || strings.HasPrefix(arg, "mcp:") || strings.HasPrefix(arg, "rest:"):
+		return Target{Kind: "invocation", Text: arg}, nil
+	case strings.HasPrefix(arg, "sha256:"):
+		if _, err := hashref.Parse(arg); err != nil {
+			return Target{}, fmt.Errorf("%q is not a content address: %v", arg, err)
+		}
+		return Target{Kind: "ref", Text: arg}, nil
+	}
+	n, err := strconv.ParseInt(arg, 10, 64)
+	if err != nil || n < 0 {
+		return Target{}, fmt.Errorf("%q is not an entry index, head, sha256:<ref>, a2a:/mcp:/rest:<id>, or -", arg)
+	}
+	return Target{Kind: "index", Index: n}, nil
+}
+
+// Options are a walk's --scan and --json.
+type Options struct {
+	Scan int
+	JSON bool
+}
+
+// Answer is exactly what the CLI writes: its stdout, its stderr, its exit.
+type Answer struct {
+	Stdout, Stderr string
+	Exit           int
+}
+
+// Run walks t to its chain and renders it as `ailang-worldd why` does: exit 0
+// when every link is ✓, 3 when any is ✗, 1 when t is not found.
+func Run(r Reader, t Target, o Options) (a Answer) {
+	var stdout, stderr strings.Builder
+	defer func() { a.Stdout, a.Stderr = stdout.String(), stderr.String() }()
+	e, matchedBy, err := r.Locate(t, o.Scan)
+	if err != nil {
+		if errors.Is(err, errScanExhausted) || errors.Is(err, context.DeadlineExceeded) {
+			fmt.Fprintf(&stderr, "ailang-worldd why: not found in the last %d entries (%s)\n", o.Scan, matchedBy)
+			return Answer{Exit: exitUsage}
+		}
+		fmt.Fprintf(&stderr, "ailang-worldd why: %v\n", err)
+		return Answer{Exit: exitUsage}
+	}
+	ch, err := r.BuildChain(e, t.Result)
+	if err != nil {
+		fmt.Fprintf(&stderr, "ailang-worldd why: %v\n", err)
+		return Answer{Exit: exitUsage}
+	}
+	ch.MatchedBy = matchedBy
+	if o.JSON {
+		b, _ := json.MarshalIndent(ch, "", "  ")
+		fmt.Fprintln(&stdout, string(b))
+	} else {
+		printChain(&stdout, ch)
+	}
+	if !ch.OK {
+		return Answer{Exit: exitIntegrity}
+	}
+	return Answer{Exit: exitOK}
+}
+
+// locate resolves a target to one log entry. matchedBy names how.
+func (r Reader) Locate(t Target, scan int) (Entry, string, error) {
+	byIndex := func(i int64, how string) (Entry, string, error) {
+		e, ok, err := r.Entry(i)
+		if err != nil {
+			return e, how, err
+		}
+		if !ok {
+			return e, how, fmt.Errorf("log entry %d not found", i)
+		}
+		return e, how, nil
+	}
+	switch t.Kind {
+	case "index":
+		return byIndex(t.Index, "entry index")
+	case "head":
+		h, ok, err := r.HeadIndex()
+		if err != nil {
+			return Entry{}, "head", err
+		}
+		if !ok {
+			return Entry{}, "head", errors.New("no world head has been selected yet")
+		}
+		return byIndex(h, "head")
+	case "invocation":
+		if strings.HasPrefix(t.Text, "rest:") {
+			var rc struct {
+				State     string `json:"state"`
+				ResultRef string `json:"resultRef"`
+			}
+			if _, err := r.GetJSON("/v1/receipts/"+url.PathEscape(t.Text), &rc); err != nil {
+				return Entry{}, "rest receipt", err
+			}
+			if rc.State != "resolved" || rc.ResultRef == "" {
+				return Entry{}, "rest receipt", fmt.Errorf("receipt %s is %s: nothing committed", t.Text, rc.State)
+			}
+			return r.ByWorld(rc.ResultRef, "rest receipt → world")
+		}
+		return r.ScanFor(scan, "invocationId", func(e Entry, rec *Record) bool {
+			return rec != nil && rec.InvocationID == t.Text
+		})
+	case "result":
+		ref := hashref.SumSHA256(t.Result).String()
+		e, how, err := r.ScanFor(scan, "output (sha256 of the result)", func(e Entry, rec *Record) bool {
+			return rec != nil && rec.Output == ref
+		})
+		if err == nil || !errors.Is(err, errScanExhausted) {
+			return e, how, err
+		}
+		var out struct {
+			World struct {
+				Plan string `json:"plan"`
+			} `json:"world"`
+		}
+		if json.Unmarshal(t.Result, &out) != nil || out.World.Plan == "" {
+			return e, how, err
+		}
+		return r.ScanFor(scan, "world.plan (fallback: the result bytes are not a committed output)", func(e Entry, rec *Record) bool {
+			return rec != nil && rec.Plan == out.World.Plan
+		})
+	}
+	// A content address: dispatch on what it names.
+	obj, ok, err := r.Object(t.Text)
+	if err != nil {
+		return Entry{}, "object", err
+	}
+	if !ok {
+		return r.ByWorld(t.Text, "world ref")
+	}
+	ref := t.Text
+	switch obj.SemanticID {
+	case coordinator.RecordV1, coordinator.RecordV2:
+		return r.ScanFor(scan, "record", func(e Entry, _ *Record) bool { return e.TransitionRef == ref })
+	case coordinator.OutputV1:
+		return r.ScanFor(scan, "output", func(_ Entry, rec *Record) bool { return rec != nil && rec.Output == ref })
+	case coordinator.InputV1:
+		return r.ScanFor(scan, "input", func(_ Entry, rec *Record) bool { return rec != nil && rec.Input == ref })
+	case coordinator.EffectPlanV1:
+		return r.ScanFor(scan, "plan", func(_ Entry, rec *Record) bool { return rec != nil && rec.Plan == ref })
+	case broker.EffectRecordV1:
+		return r.ScanFor(scan, "effect record", func(_ Entry, rec *Record) bool { return rec != nil && contains(rec.Effects, ref) })
+	case broker.EffectResultV1, broker.EffectRequestV1:
+		return r.ScanFor(scan, "effect "+strings.TrimSuffix(strings.TrimPrefix(obj.SemanticID, "world/effect-"), "/v1"),
+			func(_ Entry, rec *Record) bool {
+				if rec == nil {
+					return false
+				}
+				for _, er := range rec.Effects {
+					if d, ok := r.EffectRecord(er); ok && (d.ResultRef.String() == ref || d.RequestRef.String() == ref) {
+						return true
+					}
+				}
+				return false
+			})
+	default:
+		return r.ScanFor(scan, "transitionRef ("+obj.SemanticID+")", func(e Entry, _ *Record) bool { return e.TransitionRef == ref })
+	}
+}
+
+func (r Reader) EffectRecord(ref string) (broker.EffectRecord, bool) {
+	obj, ok, err := r.Object(ref)
+	if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 {
+		return broker.EffectRecord{}, false
+	}
+	d, err := broker.DecodeRecord(obj.Payload)
+	return d, err == nil
+}
+
+// byWorld resolves a world ref to the entry its log head names.
+func (r Reader) ByWorld(ref, how string) (Entry, string, error) {
+	w, ok, err := r.World(ref)
+	if err != nil {
+		return Entry{}, how, err
+	}
+	if !ok {
+		return Entry{}, how, fmt.Errorf("%s is neither an object nor a world in this store", ref)
+	}
+	e, ok, err := r.Entry(w.Revision)
+	if err != nil {
+		return e, how, err
+	}
+	if !ok || e.EntryHash != w.LogHead {
+		return e, how, fmt.Errorf("world %s names log head %s, which is not entry %d", ref, w.LogHead, w.Revision)
+	}
+	return e, how, nil
+}
+
+// scanFor scans back from head for the first entry match accepts.
+func (r Reader) ScanFor(scan int, how string, match func(Entry, *Record) bool) (Entry, string, error) {
+	h, ok, err := r.HeadIndex()
+	if err != nil {
+		return Entry{}, how, err
+	}
+	if !ok {
+		return Entry{}, how, errScanExhausted
+	}
+	var found *Entry
+	_, err = r.ScanEntries(0, h, true, scan, func(e Entry, rec *Record) (bool, error) {
+		if match(e, rec) {
+			found = &e
+			return true, nil
+		}
+		return false, nil
+	})
+	if err != nil {
+		return Entry{}, how, err
+	}
+	if found == nil {
+		return Entry{}, how, errScanExhausted
+	}
+	return *found, how, nil
+}
+
+func contains(xs []string, x string) bool {
+	for _, v := range xs {
+		if v == x {
+			return true
+		}
+	}
+	return false
+}
+
+// ---------------------------------------------------------------------------
+// the chain
+// ---------------------------------------------------------------------------
+
+type Link struct {
+	Name   string `json:"name"`
+	OK     bool   `json:"ok"`
+	Ref    string `json:"ref,omitempty"`
+	Detail string `json:"detail"`
+	// Worktree (row 184 D3) is an effect result's worktree_delta: the files
+	// the toolchain call created, modified or deleted in the worktree.
+	Worktree   *broker.WorktreeDelta `json:"worktree,omitempty"`
+	Protected  json.RawMessage       `json:"protected,omitempty"`
+	ExecResult bool                  `json:"-"`
+}
+
+// Keep the complete stored block for JSON; nil alone means the key is absent.
+func resultProtected(payload []byte) json.RawMessage {
+	var res struct {
+		Protected json.RawMessage `json:"protected"`
+	}
+	if json.Unmarshal(payload, &res) != nil {
+		return nil
+	}
+	return res.Protected
+}
+
+// whyProtected mirrors broker's private protectedRecord for text decoding.
+type whyProtected struct {
+	Digest            string   `json:"digest"`
+	Patterns          []string `json:"patterns"`
+	Enforced          string   `json:"enforced"`
+	Base              string   `json:"base"`
+	BaseDigest        string   `json:"base_digest"`
+	BaseFiles         int      `json:"base_files"`
+	BaseIgnoredDigest string   `json:"base_ignored_digest"`
+	BaseIgnoredFiles  int      `json:"base_ignored_files"`
+	IntactBefore      *bool    `json:"intact_before"`
+	IntactAfter       *bool    `json:"intact_after"`
+	Changes           []struct {
+		Path          string `json:"path"`
+		Status        string `json:"status"`
+		Phase         string `json:"phase"`
+		PathTruncated bool   `json:"path_truncated,omitempty"`
+	} `json:"changes"`
+	Truncated bool   `json:"truncated"`
+	Total     int    `json:"total"`
+	Error     string `json:"error"`
+}
+
+func protectedLine(raw json.RawMessage, exec bool) string {
+	if raw == nil {
+		if exec {
+			return "protected: none declared"
+		}
+		return ""
+	}
+	var p *whyProtected
+	if json.Unmarshal(raw, &p) != nil || p == nil {
+		return "protected: not checked — malformed protected block"
+	}
+	var phases []string
+	if p.IntactBefore != nil && !*p.IntactBefore {
+		phases = append(phases, "before")
+	}
+	if p.IntactAfter != nil && !*p.IntactAfter {
+		phases = append(phases, "after")
+	}
+	if len(phases) > 0 {
+		var paths []string
+		for _, c := range p.Changes {
+			if len(paths) == whyWorktreePaths {
+				break
+			}
+			paths = append(paths, fmt.Sprintf("%s (%s)", c.Path, c.Status))
+		}
+		line := "protected: CHANGED " + strings.Join(phases, " and ") + " the call — " + strings.Join(paths, ", ")
+		if p.Total > len(paths) {
+			line += fmt.Sprintf(" (+%d more)", p.Total-len(paths))
+		}
+		if p.Error != "" {
+			line += "; not checked — " + p.Error
+		}
+		return line
+	}
+	if p.IntactBefore == nil || p.IntactAfter == nil || p.Error != "" {
+		reason := p.Error
+		if reason == "" {
+			reason = "no verdict recorded"
+		}
+		return "protected: not checked — " + reason
+	}
+	return fmt.Sprintf("protected: intact (%d files @ %.7s, digest %s, enforced %s)", p.BaseFiles, p.Base, p.BaseDigest, p.Enforced)
+}
+
+// whyWorktreePaths is how many changed paths the text form prints.
+const whyWorktreePaths = 5
+
+// resultWorktreeDelta decodes an effect result's worktree_delta, or nil when
+// the result has none (or is not a JSON object).
+func resultWorktreeDelta(payload []byte) *broker.WorktreeDelta {
+	var res struct {
+		Delta *broker.WorktreeDelta `json:"worktree_delta"`
+	}
+	if json.Unmarshal(payload, &res) != nil {
+		return nil
+	}
+	return res.Delta
+}
+
+// worktreeLine renders a delta as `worktree: N created, M modified, K
+// deleted` plus the first paths (D3), or the snapshot's error or cap.
+func worktreeLine(d *broker.WorktreeDelta) string {
+	switch {
+	case d.Error != "":
+		return "worktree: not recorded — " + d.Error
+	case d.Basis == broker.WorktreeBasisWalkCapped:
+		return fmt.Sprintf("worktree: not recorded — the non-git worktree has more than %d files (basis %s)", d.CappedAt, d.Basis)
+	}
+	// Counts are of the listed changes; a truncated delta says so.
+	var created, modified, deleted int
+	paths := make([]string, 0, whyWorktreePaths)
+	for _, c := range d.Changes {
+		switch c.Op {
+		case "created":
+			created++
+		case "modified":
+			modified++
+		case "deleted":
+			deleted++
+		}
+		if len(paths) < whyWorktreePaths {
+			paths = append(paths, c.Op[:1]+" "+c.Path)
+		}
+	}
+	line := fmt.Sprintf("worktree: %d created, %d modified, %d deleted (basis %s)", created, modified, deleted, d.Basis)
+	if d.Truncated {
+		line += fmt.Sprintf(" — first %d of %d changes listed", len(d.Changes), d.Total)
+	}
+	if len(paths) > 0 {
+		line += ": " + strings.Join(paths, ", ")
+		if more := d.Total - len(paths); more > 0 {
+			line += fmt.Sprintf(" (+%d more)", more)
+		}
+	}
+	return line
+}
+
+type Chain struct {
+	Entry       int64  `json:"entry"`
+	MatchedBy   string `json:"matchedBy"`
+	Coordinator bool   `json:"coordinator"`
+	Note        string `json:"note,omitempty"`
+	Links       []Link `json:"links"`
+	OK          bool   `json:"ok"`
+}
+
+func (ch *Chain) add(name string, ok bool, ref, detail string, args ...any) {
+	ch.Links = append(ch.Links, Link{Name: name, OK: ok, Ref: ref, Detail: fmt.Sprintf(detail, args...)})
+	if !ok {
+		ch.OK = false
+	}
+}
+
+// EntryWire / WorldWire mirror host/coordinator/plan.go's canonical encodings
+// (field order is the encoding order), so the entry hash and the world ref
+// are recomputed, not trusted.
+type EntryWire struct {
+	EntryIndex     int64  `json:"entryIndex"`
+	SemanticsEpoch int64  `json:"semanticsEpoch"`
+	TransitionFn   string `json:"transitionFn"`
+	Interpreter    string `json:"interpreter"`
+	PrevEntryHash  string `json:"prevEntryHash"`
+	WrittenBy      string `json:"writtenBy"`
+	TransitionRef  string `json:"transitionRef"`
+}
+
+type WorldWire struct {
+	Revision  int64  `json:"revision"`
+	StateRoot string `json:"stateRoot"`
+	LogHead   string `json:"logHead"`
+}
+
+func SumJSON(v any) string {
+	b, _ := json.Marshal(v)
+	return hashref.SumSHA256(b).String()
+}
+
+func Sum(b []byte) string { return hashref.SumSHA256(b).String() }
+
+// checkedObject fetches ref and reports whether it resolves, carries the
+// wanted semantic id, and its payload hashes to ref.
+func (r Reader) CheckedObject(ref, semantic string) (Object, bool, string) {
+	obj, ok, err := r.Object(ref)
+	switch {
+	case err != nil:
+		return obj, false, err.Error()
+	case !ok:
+		return obj, false, "does not resolve"
+	case semantic != "" && obj.SemanticID != semantic:
+		return obj, false, fmt.Sprintf("semanticId %s, want %s", obj.SemanticID, semantic)
+	case Sum(obj.Payload) != ref:
+		return obj, false, fmt.Sprintf("payload hashes to %s, not its ref", Sum(obj.Payload))
+	}
+	return obj, true, ""
+}
+
+func failNote(ok bool, why string) string {
+	if ok {
+		return ""
+	}
+	return " — ✗ " + why
+}
+
+// buildChain renders and checks one entry's provenance chain. result is the
+// caller's result bytes (why - / --result), or nil.
+func (r Reader) BuildChain(e Entry, result []byte) (Chain, error) {
+	ch := Chain{Entry: e.Header.EntryIndex, OK: true}
+	h, err := r.Health()
+	if err != nil {
+		return ch, err
+	}
+	pin := "not the daemon's pinned interpreter"
+	if e.Header.Interpreter == h.InterpreterRef && h.InterpreterRef != "" {
+		pin = strings.TrimSpace(strings.SplitN(h.InterpreterVersion, "\n", 2)[0]) + ", the daemon's pin"
+	}
+	rec, recObj, err := r.Record(e)
+	if err != nil {
+		return ch, err
+	}
+	entryLine := fmt.Sprintf("#%d writtenBy %s prev %s transitionFn %s interpreter %s (%s)",
+		e.Header.EntryIndex, e.Header.WrittenBy, e.Header.PrevEntryHash, e.Header.TransitionFn, e.Header.Interpreter, pin)
+	if rec == nil {
+		ch.Note = fmt.Sprintf("entry %d is not a coordinator invocation (writtenBy %s; its transitionRef is a %s): "+
+			"there is no record chain to walk, so only the entry and its object are shown", e.Header.EntryIndex, e.Header.WrittenBy, recObj.SemanticID)
+		ch.add("entry", true, e.EntryHash, "%s (a client-supplied entry hash: not recomputed)", entryLine)
+		ok := Sum(recObj.Payload) == e.TransitionRef
+		ch.add("object", ok, e.TransitionRef, "%s, %d bytes%s", recObj.SemanticID, len(recObj.Payload), failNote(ok, "payload does not hash to transitionRef"))
+		return ch, nil
+	}
+	ch.Coordinator = true
+
+	// world: recomputed from the entry and the output, confirmed by the route.
+	wref := SumJSON(WorldWire{Revision: e.Header.EntryIndex, StateRoot: rec.Output, LogHead: e.EntryHash})
+	w, wok, werr := r.World(wref)
+	worldOK := werr == nil && wok && w.Revision == e.Header.EntryIndex && w.StateRoot == rec.Output && w.LogHead == e.EntryHash
+	why := "GET /v1/worlds/<recomputed ref> does not return this world"
+	if werr != nil {
+		why = werr.Error()
+	}
+	ch.add("world", worldOK, wref, "revision %d, stateRoot %s, logHead %s (recomputed)%s", e.Header.EntryIndex, rec.Output, e.EntryHash, failNote(worldOK, why))
+
+	// entry: the hash is recomputed over the canonical header + transitionRef.
+	eh := SumJSON(EntryWire{EntryIndex: e.Header.EntryIndex, SemanticsEpoch: e.Header.SemanticsEpoch,
+		TransitionFn: e.Header.TransitionFn, Interpreter: e.Header.Interpreter, PrevEntryHash: e.Header.PrevEntryHash,
+		WrittenBy: e.Header.WrittenBy, TransitionRef: e.TransitionRef})
+	ch.add("entry", eh == e.EntryHash, e.EntryHash, "%s%s", entryLine, failNote(eh == e.EntryHash, "entry hash recomputes to "+eh))
+
+	// record: its hash, and its fields agree with the entry.
+	recOK := Sum(rec.payload) == e.TransitionRef && rec.TransitionFn == e.Header.TransitionFn &&
+		rec.Interpreter == e.Header.Interpreter && rec.SemanticsEpoch == e.Header.SemanticsEpoch
+	ch.add("record", recOK, e.TransitionRef, "%s invocation %s episode %s skill %s%s", rec.SemanticID, rec.InvocationID, rec.EpisodeID, rec.SkillID,
+		failNote(recOK, "the record's hash or its transitionFn/interpreter/epoch disagree with the entry"))
+
+	// input.
+	in, inOK, inWhy := r.CheckedObject(rec.Input, coordinator.InputV1)
+	ch.add("input", inOK, rec.Input, "%s%s", whyInputLine(in.Payload), failNote(inOK, inWhy))
+
+	// plan and effects (v2 only).
+	var planned []struct {
+		ID     string `json:"id"`
+		Effect string `json:"effect"`
+		Scope  string `json:"scope"`
+		Cost   int64  `json:"cost"`
+	}
+	if rec.SemanticID == coordinator.RecordV2 {
+		po, pOK, pWhy := r.CheckedObject(rec.Plan, coordinator.EffectPlanV1)
+		var plan struct {
+			Effects []struct {
+				ID     string `json:"id"`
+				Effect string `json:"effect"`
+				Scope  string `json:"scope"`
+				Cost   int64  `json:"cost"`
+			} `json:"effects"`
+			Result json.RawMessage `json:"result"`
+		}
+		detail := "(undecodable)"
+		if pOK && json.Unmarshal(po.Payload, &plan) == nil {
+			planned = plan.Effects
+			if len(plan.Effects) == 0 {
+				detail = "no effects; result " + Elide(string(plan.Result))
+			} else {
+				parts := make([]string, 0, len(plan.Effects))
+				for _, pe := range plan.Effects {
+					parts = append(parts, fmt.Sprintf("%s %s@%s cost %d", pe.ID, pe.Effect, pe.Scope, pe.Cost))
+				}
+				detail = fmt.Sprintf("%d effect(s): %s", len(plan.Effects), strings.Join(parts, "; "))
+			}
+		}
+		ch.add("plan", pOK, rec.Plan, "%s%s", detail, failNote(pOK, pWhy))
+		for i, ref := range rec.Effects {
+			obj, ok, oWhy := r.CheckedObject(ref, broker.EffectRecordV1)
+			id := fmt.Sprintf("#%d", i)
+			if i < len(planned) {
+				id = planned[i].ID
+			}
+			if !ok {
+				ch.add("effect", false, ref, "%s — ✗ %s", id, oWhy)
+				continue
+			}
+			d, err := broker.DecodeRecord(obj.Payload)
+			if err != nil {
+				ch.add("effect", false, ref, "%s — ✗ %v", id, err)
+				continue
+			}
+			verdict := "allowed"
+			switch {
+			case !d.Allowed:
+				verdict = "denied:" + d.Denial
+			case d.Failed:
+				verdict = "failed"
+			}
+			okAll := broker.RecordConsistent(d)
+			req, reqOK, _ := r.CheckedObject(d.RequestRef.String(), broker.EffectRequestV1)
+			okAll = okAll && reqOK
+			resText := "no result"
+			if d.Failure != "" {
+				// Row 255: the failed record's bounded failure line.
+				resText = "no result: " + d.Failure
+			}
+			var delta *broker.WorktreeDelta
+			var protected json.RawMessage
+			execResult := false
+			if !d.ResultRef.IsZero() {
+				res, resOK, _ := r.CheckedObject(d.ResultRef.String(), broker.EffectResultV1)
+				okAll = okAll && resOK
+				resText = fmt.Sprintf("result %d B", len(res.Payload))
+				if resOK {
+					delta = resultWorktreeDelta(res.Payload)
+					protected = resultProtected(res.Payload)
+					execResult = d.Effect == broker.EffectWorkspaceExec
+				}
+			}
+			ch.add("effect", okAll, ref, "%s %s %s@%s cost %d budget %d→%d request %d B %s%s", id, verdict, d.Effect, d.Scope, d.Cost,
+				d.BudgetBefore, d.BudgetAfter, len(req.Payload), resText,
+				failNote(okAll, "the record is inconsistent or its request/result does not resolve to its hash"))
+			ch.Links[len(ch.Links)-1].Worktree = delta
+			ch.Links[len(ch.Links)-1].Protected = protected
+			ch.Links[len(ch.Links)-1].ExecResult = execResult
+		}
+	}
+
+	// output: hashes to record.output, equals world.stateRoot, and its world
+	// block names this record's plan and effect records.
+	out, oOK, oWhy := r.CheckedObject(rec.Output, coordinator.OutputV1)
+	if oOK && !(worldOK && w.StateRoot == rec.Output) {
+		oOK, oWhy = false, "the output is not the committed world's stateRoot"
+	}
+	if oOK && rec.SemanticID == coordinator.RecordV2 {
+		var ob struct {
+			World struct {
+				Plan    string `json:"plan"`
+				Effects []struct {
+					Record *string `json:"record"`
+				} `json:"effects"`
+			} `json:"world"`
+		}
+		_ = json.Unmarshal(out.Payload, &ob)
+		var recs []string
+		for _, fx := range ob.World.Effects {
+			if fx.Record != nil {
+				recs = append(recs, *fx.Record)
+			}
+		}
+		if ob.World.Plan != rec.Plan || strings.Join(recs, ",") != strings.Join(rec.Effects, ",") {
+			oOK, oWhy = false, "the output's world block does not name this record's plan and effect records"
+		}
+	}
+	detail := fmt.Sprintf("%d B = world.stateRoot", len(out.Payload))
+	if result != nil {
+		if Sum(result) == rec.Output {
+			detail += "; your result is these exact bytes"
+		} else {
+			oOK, oWhy = false, "your result ("+Sum(result)+") is not the committed output bytes"
+		}
+	}
+	ch.add("output", oOK, rec.Output, "%s%s", detail, failNote(oOK, oWhy))
+	return ch, nil
+}
+
+func printChain(w io.Writer, ch Chain) {
+	fmt.Fprintf(w, "why: entry %d (matched by %s)\n", ch.Entry, ch.MatchedBy)
+	if ch.Note != "" {
+		fmt.Fprintf(w, "  note: %s\n", ch.Note)
+	}
+	bad := 0
+	for _, l := range ch.Links {
+		mark := "✓"
+		if !l.OK {
+			mark, bad = "✗", bad+1
+		}
+		ref := ""
+		if l.Ref != "" {
+			ref = l.Ref + "  "
+		}
+		fmt.Fprintf(w, "  %s %-7s %s%s\n", mark, l.Name, ref, l.Detail)
+		if l.Worktree != nil {
+			fmt.Fprintf(w, "            %s\n", worktreeLine(l.Worktree))
+		}
+		if line := protectedLine(l.Protected, l.ExecResult); line != "" {
+			fmt.Fprintf(w, "            %s\n", line)
+		}
+	}
+	if bad == 0 {
+		fmt.Fprintf(w, "all %d link(s) verified\n", len(ch.Links))
+	} else {
+		fmt.Fprintf(w, "%d of %d link(s) BROKEN — integrity refusal (exit 3)\n", bad, len(ch.Links))
+	}
+}
+
+// whyKeyFieldMax bounds one key field's encoded value; whyKeyFieldCount the
+// number of key fields shown.
+const (
+	whyKeyFieldMax   = 120
+	whyKeyFieldCount = 6
+)
+
+// whyInputLine renders an invocation input for `why`. A payload that fits in
+// ElideAt is printed whole, unchanged. A longer one is elided — and because
+// canonical JSON orders keys, the field that names WHAT a call touched can
+// sit past the cut (ailang-edit's `path` follows `new_text`; row 114 F-7, a
+// walk missed an edited file this way). So the short top-level fields of an
+// object input — each scalar or small array whose encoding is at most
+// whyKeyFieldMax bytes, in key order, at most whyKeyFieldCount — are printed
+// first, then the elided payload.
+func whyInputLine(payload []byte) string {
+	one := strings.ReplaceAll(string(payload), "\n", `\n`)
+	if len(one) <= ElideAt {
+		return one
+	}
+	var obj map[string]json.RawMessage
+	if json.Unmarshal(payload, &obj) != nil || len(obj) == 0 {
+		return Elide(string(payload))
+	}
+	keys := make([]string, 0, len(obj))
+	for k := range obj {
+		keys = append(keys, k)
+	}
+	sort.Strings(keys)
+	var fields []string
+	for _, k := range keys {
+		v := obj[k]
+		if len(v) == 0 || len(v) > whyKeyFieldMax || v[0] == '{' {
+			continue
+		}
+		fields = append(fields, k+"="+strings.ReplaceAll(string(v), "\n", `\n`))
+		if len(fields) == whyKeyFieldCount {
+			break
+		}
+	}
+	if len(fields) == 0 {
+		return Elide(string(payload))
+	}
+	return strings.Join(fields, " ") + " · " + Elide(string(payload))
+}
+
+// ElideAt is how many bytes of a value Elide keeps.
+const ElideAt = 160
+
+// Elide shortens a value to ElideAt bytes plus its byte count, on one line.
+func Elide(s string) string {
+	one := strings.ReplaceAll(s, "\n", `\n`)
+	if len(one) <= ElideAt {
+		return one
+	}
+	return fmt.Sprintf("%s… (%d bytes)", strings.ToValidUTF8(one[:ElideAt], ""), len(s))
+}
diff --git a/host/why/why_input_test.go b/host/why/why_input_test.go
new file mode 100644
--- /dev/null
+++ b/host/why/why_input_test.go
@@ -0,0 +1,62 @@
+package why
+
+import (
+	"strconv"
+	"strings"
+	"testing"
+)
+
+// Row 114 F-7: `why 198` on an ailang-edit showed `{"new_text":"// Row 185:
+// IsLinked accepts…… (1048 bytes)` — the edited file's path sat past the
+// cut, and a walk concluded the file had not been edited.
+func TestWhyInputLineShowsThePathOfALongEdit(t *testing.T) {
+	payload := []byte(`{"new_text":"` + strings.Repeat(`// Row 185: IsLinked accepts the worktree Add makes\n`, 20) +
+		`","old_text":"}\n","path":"host/worktree/worktree_test.go"}`)
+	got := whyInputLine(payload)
+	if !strings.HasPrefix(got, `old_text="}\n" path="host/worktree/worktree_test.go" · {"new_text":`) {
+		t.Fatalf("key fields not first:\n%s", got)
+	}
+	if !strings.HasSuffix(got, "("+strconv.Itoa(len(payload))+" bytes)") {
+		t.Fatalf("the elided payload and its size must still follow:\n%s", got)
+	}
+	if strings.Contains(got, "new_text=") {
+		t.Fatalf("a long field must not be repeated as a key field:\n%s", got)
+	}
+}
+
+// A payload that fits is printed whole and byte-identical to before (no key
+// fields), so short inputs — most of them — read exactly as they did.
+func TestWhyInputLineLeavesAShortInputUnchanged(t *testing.T) {
+	for _, p := range []string{`{"path":"cmd/ailang-worldd/setup.go"}`, `{"args":["-run","TestX","./cmd"],"command":"go-test"}`, `"x"`} {
+		if got := whyInputLine([]byte(p)); got != Elide(p) || got != p {
+			t.Fatalf("short input changed: %q -> %q", p, got)
+		}
+	}
+}
+
+// Small arrays count as key fields (an exec call's args); nested objects and
+// over-long values do not; at most whyKeyFieldCount are shown.
+func TestWhyInputLineKeyFieldBounds(t *testing.T) {
+	long := strings.Repeat("a", whyKeyFieldMax)
+	payload := []byte(`{"a1":1,"a2":2,"a3":3,"a4":4,"a5":5,"a6":6,"a7":7,"args":["x"],"big":"` + long +
+		`","nested":{"k":"v"}}`)
+	got := whyInputLine(payload)
+	head, _, ok := strings.Cut(got, " · ")
+	if !ok || head != "a1=1 a2=2 a3=3 a4=4 a5=5 a6=6" {
+		t.Fatalf("key fields = %q (want the first %d in key order)", head, whyKeyFieldCount)
+	}
+	small := []byte(`{"args":["-run","TestSessionNew","./cmd/ailang-worldd"],"command":"go-test","note":"` + long + `"}`)
+	head, _, _ = strings.Cut(whyInputLine(small), " · ")
+	if head != `args=["-run","TestSessionNew","./cmd/ailang-worldd"] command="go-test"` {
+		t.Fatalf("array/scalar key fields = %q", head)
+	}
+}
+
+// Not an object, or an object with no short field: plain elide, as before.
+func TestWhyInputLineFallsBackToElide(t *testing.T) {
+	for _, p := range []string{`"` + strings.Repeat("z", 300) + `"`, `{"blob":"` + strings.Repeat("z", 300) + `"}`, strings.Repeat("not json ", 40)} {
+		if got := whyInputLine([]byte(p)); got != Elide(p) {
+			t.Fatalf("fallback changed: %q", got)
+		}
+	}
+}
diff --git a/host/why/why_protected_test.go b/host/why/why_protected_test.go
new file mode 100644
--- /dev/null
+++ b/host/why/why_protected_test.go
@@ -0,0 +1,104 @@
+package why
+
+import (
+	"bytes"
+	"encoding/json"
+	"reflect"
+	"testing"
+)
+
+const protectedFixture = `{"digest":"sha256:patterns","patterns":["*_test.go","tests/**"],"enforced":"dirs-only","base":"123456789abcdef","base_digest":"sha256:tree","base_files":12,"base_ignored_digest":"sha256:ignored","base_ignored_files":3,"intact_before":null,"intact_after":null,"changes":[{"path":"tests/a","status":"untracked","phase":"after","path_truncated":true}],"truncated":true,"total":40,"error":"snapshot failed"}`
+
+func TestWhyProtectedLine(t *testing.T) {
+	for _, tc := range []struct {
+		name, block, want string
+		exec              bool
+	}{
+		{"intact", `{"intact_before":true,"intact_after":true,"base_files":12,"base":"123456789","base_digest":"sha256:tree","enforced":"glob"}`, "protected: intact (12 files @ 1234567, digest sha256:tree, enforced glob)", true},
+		{"changed-before", `{"intact_before":false,"intact_after":true,"changes":[{"path":"a_test.go","status":"modified","phase":"before"}],"total":1}`, "protected: CHANGED before the call — a_test.go (modified)", true},
+		{"changed-after", `{"intact_before":true,"intact_after":false,"changes":[{"path":"tests/a","status":"untracked","phase":"after"}],"total":1}`, "protected: CHANGED after the call — tests/a (untracked)", true},
+		{"changed-both", `{"intact_before":false,"intact_after":false,"changes":[{"path":"a","status":"deleted","phase":"before"},{"path":"b","status":"ignored","phase":"after"}],"total":2}`, "protected: CHANGED before and after the call — a (deleted), b (ignored)", true},
+		{"index-flag", `{"intact_before":false,"changes":[{"path":"a_test.go","status":"index-flag"}],"total":1}`, "protected: CHANGED before the call — a_test.go (index-flag)", true},
+		{"head-moved", `{"intact_after":false,"changes":[{"path":"HEAD","status":"head-moved"}],"total":1}`, "protected: CHANGED after the call — HEAD (head-moved)", true},
+		{"truncated", `{"intact_after":false,"changes":[{"path":"a","status":"modified"},{"path":"b","status":"deleted"},{"path":"c","status":"untracked"},{"path":"d","status":"ignored"},{"path":"e","status":"modified"},{"path":"f","status":"deleted"}],"truncated":true,"total":40}`, "protected: CHANGED after the call — a (modified), b (deleted), c (untracked), d (ignored), e (modified) (+35 more)", true},
+		{"error", protectedFixture, "protected: not checked — snapshot failed", true},
+		{"error-beside-changed", `{"intact_after":false,"changes":[{"path":"a","status":"modified"}],"total":1,"error":"before failed"}`, "protected: CHANGED after the call — a (modified); not checked — before failed", true},
+		{"error-beside-true", `{"intact_before":true,"intact_after":true,"error":"failed"}`, "protected: not checked — failed", true},
+		{"missing-verdict", `{"intact_before":true}`, "protected: not checked — no verdict recorded", true},
+		{"missing-before", `{"intact_after":true,"base_files":12,"base":"123456789","base_digest":"sha256:tree","enforced":"glob"}`, "protected: not checked — no verdict recorded", true},
+		{"malformed", `"broken"`, "protected: not checked — malformed protected block", true},
+		{"malformed-null", `null`, "protected: not checked — malformed protected block", true},
+		{"malformed-field", `{"intact_before":"true"}`, "protected: not checked — malformed protected block", true},
+		{"none-declared", "", "protected: none declared", true},
+		{"non-exec", "", "", false},
+		{"old-record", "", "protected: none declared", true},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			payload := []byte(`{"content":"old"}`)
+			if tc.block != "" {
+				payload = []byte(`{"protected":` + tc.block + `}`)
+			}
+			if got := protectedLine(resultProtected(payload), tc.exec); got != tc.want {
+				t.Fatalf("protected line = %q, want %q", got, tc.want)
+			}
+		})
+	}
+}
+
+func assertProtectedBlock(t *testing.T, got, want json.RawMessage) {
+	t.Helper()
+	var g, w any
+	if len(got) != 0 {
+		if err := json.Unmarshal(got, &g); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if len(want) != 0 {
+		if err := json.Unmarshal(want, &w); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if !reflect.DeepEqual(g, w) {
+		t.Fatalf("protected JSON = %s, want %s", got, want)
+	}
+}
+
+func TestWhyProtectedJSON(t *testing.T) {
+	raw := resultProtected([]byte(`{"protected":` + protectedFixture + `}`))
+	b := mustJSONT(t, Link{Name: "effect", OK: true, Detail: "result", Protected: raw})
+	var link map[string]json.RawMessage
+	if err := json.Unmarshal(b, &link); err != nil {
+		t.Fatal(err)
+	}
+	assertProtectedBlock(t, link["protected"], json.RawMessage(protectedFixture))
+	var block map[string]json.RawMessage
+	if err := json.Unmarshal(link["protected"], &block); err != nil {
+		t.Fatal(err)
+	}
+	for _, key := range []string{"intact_before", "intact_after"} {
+		if string(block[key]) != "null" {
+			t.Fatalf("%s = %s, want null", key, block[key])
+		}
+	}
+	old := Link{Name: "effect", OK: true, Ref: "hash", Detail: "result", Protected: resultProtected([]byte(`{"content":"old"}`)), ExecResult: true}
+	b = mustJSONT(t, old)
+	if want := `{"name":"effect","ok":true,"ref":"hash","detail":"result"}`; string(b) != want {
+		t.Fatalf("old link JSON = %s, want %s", b, want)
+	}
+	var decoded Link
+	if err := json.Unmarshal(b, &decoded); err != nil || decoded.Protected != nil {
+		t.Fatalf("old record: %+v, %v", decoded, err)
+	}
+	if bytes.Contains(b, []byte(`"protected"`)) {
+		t.Fatalf("old link carries protected: %s", b)
+	}
+}
+
+func mustJSONT(t *testing.T, v any) []byte {
+	t.Helper()
+	b, err := json.Marshal(v)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return b
+}
````

### Patch M3t — the seam's tests: VW2–VW5 (`why_seam_test.go`), VS2a and the pure tests (`host/why/run_test.go`)

````diff M3t
diff --git a/cmd/ailang-worldd/why_seam_test.go b/cmd/ailang-worldd/why_seam_test.go
new file mode 100644
--- /dev/null
+++ b/cmd/ailang-worldd/why_seam_test.go
@@ -0,0 +1,160 @@
+package main
+
+// Row 243 V2a (design §6 VW2–VW5): host/why's Run is the CLI's why over any Getter; bodies are recorded from the rig.
+
+import (
+	"context"
+	"errors"
+	"go/ast"
+	"go/parser"
+	"go/token"
+	"strconv"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/why"
+)
+
+type recorded struct {
+	status int
+	body   []byte
+}
+
+// recordingGetter reads addr through the CLI's client and keeps every answer.
+func recordingGetter(addr string, into map[string]recorded) why.Getter {
+	get := clientGetter(newClient(addr))
+	return func(ctx context.Context, path string) (int, []byte, error) {
+		status, body, err := get(ctx, path)
+		if err == nil {
+			into[path] = recorded{status, body}
+		}
+		return status, body, err
+	}
+}
+
+// replayGetter answers from a recording; edit may replace one answer.
+func replayGetter(rec map[string]recorded, edit func(path string, r recorded) (recorded, error)) why.Getter {
+	return func(_ context.Context, path string) (int, []byte, error) {
+		r, ok := rec[path]
+		if !ok {
+			return 0, nil, errors.New("unrecorded " + path)
+		}
+		if edit != nil {
+			var err error
+			if r, err = edit(path, r); err != nil {
+				return 0, nil, err
+			}
+		}
+		return r.status, r.body, nil
+	}
+}
+
+func whyAnswer(get why.Getter, arg string, o why.Options) why.Answer {
+	tg, _ := why.ParseTarget(arg)
+	return why.Run(why.Reader{Ctx: context.Background(), Get: get}, tg, o)
+}
+
+// VW2: Run's exit and streams are the CLI's, live and replayed.
+func TestWhyRunReplaysTheCLI(t *testing.T) {
+	r := newWalkRig(t)
+	inv := r.invoke("ep1", "ailang-read", "task-1", "hello\n")
+	for _, c := range []struct {
+		addr, arg string
+		o         why.Options
+		exit      int
+		want      string
+	}{
+		{r.url, "1", why.Options{Scan: 500}, exitOK, "all 7 link(s) verified\n"},
+		{r.url, "head", why.Options{Scan: 500, JSON: true}, exitOK, `"matchedBy": "head"`},
+		{tamperProxy(t, r.url, inv.input, []byte("forged")), "1", why.Options{Scan: 500}, exitIntegrity, "BROKEN — integrity refusal (exit 3)\n"},
+		{r.url, "9", why.Options{Scan: 500}, exitUsage, "ailang-worldd why: log entry 9 not found\n"},
+		{r.url, "a2a:nobody:x", why.Options{Scan: 1}, exitUsage, "ailang-worldd why: not found in the last 1 entries (invocationId)\n"},
+	} {
+		args := []string{c.arg, "--scan", strconv.Itoa(c.o.Scan)}
+		if c.o.JSON {
+			args = append(args, "--json")
+		}
+		code, out, errOut := whyRun(t, c.addr, "", args...)
+		rec := map[string]recorded{}
+		live := whyAnswer(recordingGetter(c.addr, rec), c.arg, c.o)
+		replay := whyAnswer(replayGetter(rec, nil), c.arg, c.o)
+		cli := why.Answer{Stdout: out, Stderr: errOut, Exit: code}
+		if live != cli || replay != cli || code != c.exit || !strings.Contains(out+errOut, c.want) ||
+			(out == "") != (code == exitUsage) || (errOut == "") != (code != exitUsage) {
+			t.Fatalf("why %v: cli %+v\nlive %+v\nreplay %+v", args, cli, live, replay)
+		}
+	}
+}
+
+// VW3: a 404 is an answer (✗ "does not resolve"); an error or 5xx is the CLI's error line; a deadline is "not found".
+func TestWhyRunGetterClasses(t *testing.T) {
+	r := newWalkRig(t)
+	inv := r.invoke("ep1", "ailang-read", "task-1", "hello\n")
+	rec := map[string]recorded{}
+	if a := whyAnswer(recordingGetter(r.url, rec), "1", why.Options{Scan: 500}); a.Exit != exitOK {
+		t.Fatalf("recording walk: %+v", a)
+	}
+	input := "/v1/objects/" + inv.input + "?payload=true"
+	for _, c := range []struct {
+		path   string
+		answer recorded
+		err    error
+		exit   int
+		want   string
+	}{
+		{input, recorded{404, []byte(`{"error":"NotFound"}`)}, nil, exitIntegrity, "✗ input   " + inv.input + "   — ✗ does not resolve\n"},
+		{"/v1/log/1", recorded{}, errors.New("boom"), exitUsage, "ailang-worldd why: boom\n"},
+		{"/v1/log/1", recorded{500, []byte("down\n")}, nil, exitUsage, "ailang-worldd why: GET /v1/log/1 returned HTTP 500: down\n"},
+		{"/v1/log/1", recorded{404, nil}, nil, exitUsage, "ailang-worldd why: log entry 1 not found\n"},
+		{"/v1/health", recorded{}, errors.New("health down"), exitUsage, "ailang-worldd why: health down\n"},
+		{"/v1/log/1", recorded{}, context.DeadlineExceeded, exitUsage, "ailang-worldd why: not found in the last 500 entries (entry index)\n"},
+	} {
+		a := whyAnswer(replayGetter(rec, func(path string, r recorded) (recorded, error) {
+			if path != c.path {
+				return r, nil
+			}
+			return c.answer, c.err
+		}), "1", why.Options{Scan: 500})
+		if a.Exit != c.exit || !strings.Contains(a.Stdout+a.Stderr, c.want) {
+			t.Fatalf("%s %+v %v: %+v", c.path, c.answer, c.err, a)
+		}
+	}
+}
+
+// VW4: the moved constants, as literals, so a change on both sides still fails.
+func TestWhyConstantsMatchTheCLI(t *testing.T) {
+	if why.MaxResponseBytes != maxClientResponseBytes || maxClientResponseBytes != 8388608 ||
+		whyScanDefault != 500 || whyScanMax != 5000 || logPageMax != 500 || why.ElideAt != 160 {
+		t.Fatalf("constants drifted: %d %d %d %d %d", why.MaxResponseBytes, whyScanDefault, whyScanMax, logPageMax, why.ElideAt)
+	}
+}
+
+// VW5: one /v1/head reader, and one elide: cmd's are single delegations to host/why.
+func TestWhyHeadIsOneReader(t *testing.T) {
+	f, err := parser.ParseFile(token.NewFileSet(), "mcpclient.go", nil, 0)
+	if err != nil {
+		t.Fatal(err)
+	}
+	want := map[string]string{"headState": "Head", "elide": "Elide"}
+	for _, d := range f.Decls {
+		fn, ok := d.(*ast.FuncDecl)
+		if !ok || want[fn.Name.Name] == "" {
+			continue
+		}
+		var sel *ast.SelectorExpr
+		if len(fn.Body.List) == 1 {
+			if ret, ok := fn.Body.List[0].(*ast.ReturnStmt); ok && len(ret.Results) == 1 {
+				if call, ok := ret.Results[0].(*ast.CallExpr); ok {
+					sel, _ = call.Fun.(*ast.SelectorExpr)
+				}
+			}
+		}
+		if sel == nil || sel.Sel.Name != want[fn.Name.Name] || ast.Unparen(sel.X).(*ast.Ident).Name != "why" {
+			t.Fatalf("%s is not one return of why.%s", fn.Name.Name, want[fn.Name.Name])
+		}
+		delete(want, fn.Name.Name)
+	}
+	if len(want) != 0 {
+		t.Fatalf("not found in mcpclient.go: %v", want)
+	}
+}
diff --git a/host/why/run_test.go b/host/why/run_test.go
new file mode 100644
--- /dev/null
+++ b/host/why/run_test.go
@@ -0,0 +1,100 @@
+package why
+
+// Row 243 V2a (design §6 VS2a, §2.1): the walk's pure parts and its imports.
+
+import (
+	"context"
+	"errors"
+	"go/parser"
+	"go/token"
+	"path/filepath"
+	"strconv"
+	"strings"
+	"testing"
+)
+
+// VS2a: host/why's production files import no effect (DIRECT imports only;
+// host/broker and host/coordinator reach host/store transitively).
+func TestWhyImportsNoEffect(t *testing.T) {
+	files, _ := filepath.Glob("*.go")
+	seen := map[string]bool{}
+	for _, f := range files {
+		if strings.HasSuffix(f, "_test.go") {
+			continue
+		}
+		p, err := parser.ParseFile(token.NewFileSet(), f, nil, parser.ImportsOnly)
+		if err != nil {
+			t.Fatal(err)
+		}
+		for _, i := range p.Imports {
+			path, _ := strconv.Unquote(i.Path.Value)
+			seen[path] = true
+		}
+	}
+	if !seen["github.com/sunholo-data/ailang-world/host/broker"] { // the instrument's control
+		t.Fatalf("imports not read: %v", seen)
+	}
+	for _, bad := range []string{"os", "os/exec", "net", "net/http", "io/fs", "database/sql", "github.com/sunholo-data/ailang-world/host/store"} {
+		if seen[bad] {
+			t.Errorf("host/why imports %s", bad)
+		}
+	}
+}
+
+// ParseTarget: each form, and the CLI's error texts (why.go at the base, :316-336).
+func TestParseTargetForms(t *testing.T) {
+	ref := "sha256:" + strings.Repeat("a", 64)
+	for arg, want := range map[string]Target{"head": {Kind: "head"}, "0": {Kind: "index"}, "42": {Kind: "index", Index: 42}, "a2a:e:t": {Kind: "invocation", Text: "a2a:e:t"}, "mcp:x": {Kind: "invocation", Text: "mcp:x"}, "rest:x": {Kind: "invocation", Text: "rest:x"}, ref: {Kind: "ref", Text: ref}} {
+		if got, err := ParseTarget(arg); err != nil || got.Kind != want.Kind || got.Index != want.Index || got.Text != want.Text {
+			t.Errorf("ParseTarget(%q) = %+v, %v", arg, got, err)
+		}
+	}
+	nt := " is not an entry index, head, sha256:<ref>, a2a:/mcp:/rest:<id>, or -"
+	for arg, want := range map[string]string{"": "a target is required", "-1": `"-1"` + nt, "x": `"x"` + nt,
+		"sha256:zz": `"sha256:zz" is not a content address: hashref: digest is not lowercase hexadecimal: "zz"`} {
+		if _, err := ParseTarget(arg); err == nil || err.Error() != want {
+			t.Errorf("ParseTarget(%q) error = %v, want %s", arg, err, want)
+		}
+	}
+}
+
+// Head and GetJSON: 200 is the answer, 404 is absent, any other status or error is an error.
+func TestHeadAndGetJSONClasses(t *testing.T) {
+	answer := func(status int, body string, err error) Getter {
+		return func(context.Context, string) (int, []byte, error) { return status, []byte(body), err }
+	}
+	ctx := context.Background()
+	if ref, ok, err := Head(ctx, answer(200, " sha256:w\n", nil)); ref != "sha256:w" || !ok || err != nil {
+		t.Errorf("200: %q %v %v", ref, ok, err)
+	}
+	if ref, ok, err := Head(ctx, answer(404, "x", nil)); ref != "" || ok || err != nil {
+		t.Errorf("404: %q %v %v", ref, ok, err)
+	}
+	if _, ok, err := Head(ctx, answer(500, "down\n", nil)); ok || err == nil || err.Error() != "GET /v1/head returned HTTP 500: down" {
+		t.Errorf("500: %v %v", ok, err)
+	}
+	if _, ok, err := Head(ctx, answer(0, "", errors.New("boom"))); ok || err == nil || err.Error() != "boom" {
+		t.Errorf("err: %v %v", ok, err)
+	}
+	var v struct{ A int }
+	if ok, err := (Reader{Ctx: ctx, Get: answer(200, "{", nil)}).GetJSON("/v1/x", &v); ok || err == nil || !strings.HasPrefix(err.Error(), "GET /v1/x: malformed body:") {
+		t.Errorf("malformed: %v %v", ok, err)
+	}
+	done, cancel := context.WithCancel(ctx)
+	cancel()
+	if ok, err := (Reader{Ctx: done, Get: answer(200, `{"A":1}`, nil)}).GetJSON("/v1/x", &v); ok || !errors.Is(err, context.Canceled) {
+		t.Errorf("a done context still read: %v %v", ok, err)
+	}
+}
+
+// Elide's bounds, as the base's elide (mcpclient.go:651) answered them (run there, plan §1.2 C21).
+func TestElideBounds(t *testing.T) {
+	a := strings.Repeat("a", 160)
+	for in, want := range map[string]string{a: a, a + "a": a + "… (161 bytes)",
+		strings.Repeat("a\n", 81): strings.Repeat(`a\n`, 53) + "a… (162 bytes)",
+		strings.Repeat("é", 100):  strings.Repeat("é", 80) + "… (200 bytes)"} {
+		if got := Elide(in); got != want {
+			t.Errorf("Elide(%q) = %q, want %q", in, got, want)
+		}
+	}
+}
````

