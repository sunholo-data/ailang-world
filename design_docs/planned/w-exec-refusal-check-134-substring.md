# w-exec-refusal-check-134-substring: `requireF15` reds any refusal whose temp path contains "134"

- **Status**: Planned (designed iteration 320, r2 after quorum round 1 BLOCKED; see §10). Queue row 234 (the controller filed it this
  iteration; it is not in the charter at the base yet, V19). It is a red on `dev`, so it outranks the queue.
- **Item**: a test-helper defect from row 219 M3 (`44eed437`, w219c, Codex in an attended World
  episode). The helper landed in PR #418, merged 2026-10-29T16:11:00Z (V5). Only test code changes.
- **Estimated**: ~0.1 day. About +75/−8 LOC, all in `host/broker/exec_tools_test.go`.
- **Measurement base**: `785b3bd1` (worktree `sprint/world-iter320-f15-tempdir`, clean, V1).
- **Instruments**: `go1.26.6 darwin/arm64`, `AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0,
  the CI pin, V18). I ran the prototypes and mutants in throwaway copies (`/tmp/w320d/scratch`,
  `/tmp/w320d/full`, made with `git archive HEAD`), never in the worktree.

---

## 1. The defect

`requireF15` (`host/broker/exec_tools_test.go:75`) asserts that a call was refused **before
launch**. Its guard at `:89` is:

```go
if refusal.OK || refusal.Refused == "" || strings.Contains(string(out), "134") || f.h.Spawns() != before {
```

The `"134"` term means "the call launched and exited 134" (SIGABRT, dyld). But it scans the
whole reply, and every F15 refusal quotes the tool's paths. Those paths sit under `t.TempDir()`,
which is `<GOTMPDIR or TMPDIR>/<TestName><random uint32>/<seq>` (V10). So the guard fails
whenever the random suffix, the temp root or the test name contains `134`.

- **CI, two reds on `dev` (V4).** Run 40909036384 (`171ece91`, "go race gate"):
  `TestExecToolCellarProfilePathRefusesByName`, path `/tmp/…RefusesByName3490134876/001/…`.
  Run 40911404613 attempt 2 (`785b3bd1`, "go host build + test gate"):
  `TestExecToolSymlinkInsideWorkspaceRefuses`, path `/tmp/…Refuses596413420/001/…` (`5964`**`134`**`20`).
  Each run has exactly one `--- FAIL`, and both messages start `not a pre-launch refusal:`.
- **Reproduced (V6, V7).** With `TMPDIR=/tmp/w320d/tmp-x134x`, all 6 callers fail (0 PASS, 6 FAIL).
  The "control" dir `tmp-ctl` failed too on its first run, from a natural suffix `4134186054`.
  At `-count=50` in `tmp-ctl` I got 295 PASS and 5 FAIL. Every failing path had `134` in its random
  suffix. That is a **~1.7% per-test flake** (5/300 on darwin), and CI's linux runners hit it twice in two runs (V4). It is not a CI quirk.
- **The term never did its job (V8, V9).** I measured both real ways a launch can exit 134:
  1. **Healthy node, the command aborts.** The reply is `execResult` with `exit_code: 134`
     (`handlers_exec.go:209`, no `omitempty`) and no `refused` key. The guard's
     `Refused == ""` and its spawn counter (`:322`, incremented before `runBounded` at `:323`)
     already catch this. The `"134"` match is redundant.
  2. **Broken node: the launch exits 134 and the `--version` re-check fails** (`:326-337`). The
     reply is `{"ok":false,"refused":"exec node … does not start: dyld: …"}`. That is the same
     envelope a pre-launch refusal uses, and its text contains **no** `134` (`has134=false`). Only
     the spawn counter (0→1) shows it launched. The `"134"` match **misses** it.

  So the term adds only false reds. The real guards are the envelope shape and the spawn counter.

## 2. The decision: a structural predicate; the substring goes

Pull the guard out into a pure test helper, and have `requireF15` call it:

```go
// f15RefusalProblem says why out is not a pre-launch F15 refusal ("" when it is). out must be
// exactly one complete JSON value (trailing whitespace allowed): the two-key envelope. A launch is
// seen by structure, never by text (row 234): a launched reply is execResult, which always
// carries exit_code, so a strict decode into the two-key envelope fails; a post-launch refusal
// (a 134 exit with a broken node, handlers_exec.go:326) has the envelope's shape but moved the
// spawn counter.
func f15RefusalProblem(out []byte, before, after int64) string {
	dec := json.NewDecoder(bytes.NewReader(out))
	dec.DisallowUnknownFields()
	var r struct {
		OK      *bool   `json:"ok"`
		Refused *string `json:"refused"`
	}
	if err := dec.Decode(&r); err != nil { return "not the refusal envelope: " + err.Error() }
	if err := dec.Decode(&json.RawMessage{}); err != io.EOF { // [r2] exactly one JSON value
		return fmt.Sprintf("trailing content after the envelope: %v", err)
	}
	if r.OK == nil || *r.OK { return "ok is not false" }
	if r.Refused == nil || *r.Refused == "" { return "refused is empty" }
	if after != before { return "the call spawned" }
	return ""
}
```

`requireF15` then fails with `not a pre-launch refusal (<problem>): <out> (spawns a -> b)`.
That keeps the CI signature's prefix. The `parts` loop is unchanged.

**[r2] Exactly one JSON value.** `json.Decoder.Decode` reads only the first value. Without the
second `Decode`, `<envelope>{"exit_code":134}` and `<envelope> garbage` are both accepted (V20).
So the strict claim is: `out` is **exactly one complete JSON value**, it is the two-key envelope,
and only whitespace may follow it. The second `Decode` must return `io.EOF`. Concatenated JSON
returns `nil`, and garbage returns a syntax error, so both are rejected. Trailing `" \n\t"` returns
`io.EOF` and is accepted (V20). This adds the `io` and `fmt` imports.

**Why a real exit-134 launch is still caught.** Shape 1 fails the strict decode (`exit_code`,
`argv`, … are unknown fields), and it fails the spawn check. Shape 2 fails the spawn check. Both
are pinned by mutants below, on live fixtures. The new check needs no exit-code text at all,
and no path can satisfy or break it.

Alternatives I rejected:
- **(a) Narrow the token to `"exit status 134"`.** It is still text, and it still misses shape 2,
  whose text has no exit status.
- **(b) Match `"134"` only outside path-like spans.** This is a parser for a free-text field.
- **(c) Delete the term and add nothing.** That would leave the envelope check loose
  (`json.Unmarshal` ignores extra keys).

The strict decode is what turns "this reply came from a launch" into a structural fact.

## 3. Risk

The change touches test code only. **The predicate is stricter than today's guard.** It also
rejects a refusal that carries any extra key, or any non-whitespace content after the envelope. Production emits exactly `{"ok","refused"}`
(`execRefusal`, `handlers_exec.go:96-97`). With the r2 prototype, all 6 callers pass 20/20 with the
root forced through `GOTMPDIR`, both for a `134` root and for a control root (V21). The whole
`host/broker` package passes under a `GOTMPDIR` `134` root, plain and under `-race` (V24).
A future refusal key would red these tests loudly, never silently. That is the intended direction.

## 4. Mutation table

| ID | Mutant (on the fixed helper) | Killer: the assertion that fires (all measured: r1 V12–V13, re-run on the r2 helper V22) |
|---|---|---|
| M1 MUT-PATH-SUBSTRING | today's `strings.Contains(out, "134")` restored as a reject | T1 rows `refusal-temp-root-134` and (r2) `trailing-whitespace`, whose envelope also carries the `134` paths; T2 (`not a pre-launch refusal`) |
| M2 MUT-NO-SPAWN-GUARD | drop `after != before` | T1 row `post-launch-f15`; T3 `broken-node` (`a launch that exited 134 passes requireF15`) |
| M3 MUT-LOOSE-DECODE | drop `DisallowUnknownFields` | T1 row `refusal-plus-exit-code` |
| M4 MUT-OK-IGNORED | drop the `ok` check | T1 row `ok-true` |
| M5 MUT-REFUSED-OPTIONAL | accept `"refused":""` | T1 row `refused-empty` |
| M7 MUT-NO-EOF-CHECK **[r2]** | the second `Decode`'s `err != io.EOF` becomes `false && …` | T1 rows `concatenated-json` and `trailing-garbage` (`problem "", want accepted=false`), and nothing else (V22) |
| M6 MUT-ACCEPT-ALL | predicate returns `""` | T1's 7 reject rows; T3 `broken-node` **and** `command-aborted` (the guard's purpose survives) |
| M0 OLD-HELPER | base `requireF15` (`:75-98`) with T2 added | T2: 10/10 FAIL at `-count=10` (the paired red at base) |

## 5. Tests

- **T1 `TestF15RefusalProblemIsStructural`** (pure table, no fixture). Rows:
  - `refusal-temp-root-134`: an F15 envelope whose paths hold `x134x` and `…1349870625`, with
    spawns 0→0. It is accepted.
  - `post-launch-f15`: the same envelope with spawns 0→1. Rejected.
  - `launched-exit-134`: an `execResult` with `exit_code:134`. Rejected.
  - `refusal-plus-exit-code`, `ok-true`, `refused-empty`: all rejected.
  - **[r2]** `concatenated-json` (the envelope then `{"exit_code":134}`) and `trailing-garbage`
    (the envelope then ` garbage`) are rejected. The control `trailing-whitespace` (the envelope
    then `" \n\t"`) is accepted.
- **T2 `TestExecToolF15RefusalUnderTempRoot134`** (live and deterministic: it forces the temp
  root).
  1. `root := os.MkdirTemp("", "w234-x134x-")`, removed in cleanup.
  2. `t.Setenv("GOTMPDIR", root)`, before the test's first `TempDir`. `testing` reads
     `GOTMPDIR` first and makes one parent per test (V10), so setting `TMPDIR` would not be
     enough where `GOTMPDIR` is set.
  3. Build `newMovingExec`, then check the precondition: `strings.Contains(m.link, "x134x")`, or
     fatal "temp root not honoured". This is the S6 null case: the test cannot pass by never
     building the hazard.
  4. Remove the target and call `requireF15(t, m.f, m.link, m.a, "the target is gone")`.

  The test name itself contains `134`, so the parent dir carries it a second time (V13). The
  test does not rely on that.
- **T3: extend the existing `TestExecToolExit134ChecksNode`** (`:230`). Read `before := f.h.Spawns()` before
  `f.call`, then assert `f15RefusalProblem(out, before, f.h.Spawns()) != ""` in both subtests.
  These are the live shapes 1 and 2 of §1.
- **AC-A [r2, the root is forced through `GOTMPDIR`].** Run
  `GOTMPDIR=/tmp/wNNNd/tmp-x134x go test ./host/broker/ -run '^(<six>|TestF15RefusalProblemIsStructural|TestExecToolF15RefusalUnderTempRoot134|TestExecToolExit134ChecksNode)$' -count=20 -v`,
  then the same with a control root. The assignment overrides any ambient `GOTMPDIR`, and per V10
  it also wins over an ambient `TMPDIR`. Want 0 FAIL in each run and a `--- PASS:` line for every
  one of the 9 names (prototype: 180 PASS, 0 FAIL, 9 names, both roots; V21). r1's
  `TMPDIR`-only form was vacuous under an ambient `GOTMPDIR`: at base,
  `GOTMPDIR=<ctl> TMPDIR=<x134x>` gives 6 PASS where the root should red all six (V23).
- **AC-A-ROOT [r2], the root-honoured assertion.** glm proposed that `grep -c tmp-x134x` over the
  `-v` output equal the PASS count. That cannot be met: passing tests print no paths, so the count
  is **0** for 180 PASS (V21). The assertion that can be run:
  1. Build once with `go test -c -o /tmp/wNNNd/broker.test ./host/broker/`. A prebuilt binary is
     needed because the go tool puts its own work dir under `GOTMPDIR`; a read-only root fails
     the build first, with `go: creating work dir: mkdir …: permission denied` (V23).
  2. From `host/broker/`, run `GOTMPDIR=<root>-ro /tmp/wNNNd/broker.test -test.run '^(<six>)$' -test.v -test.timeout=2m`,
     where `<root>-ro` is mode 555.
  3. Want 6 top-level `--- FAIL`, and `grep -c 'TempDir: mkdir <root>-ro'` equal to the number of
     leaf `--- FAIL` lines. A leaf is every `--- FAIL` except the `TestExecToolPlacementReapplied`
     parent. The prototype gave 10 leaf FAILs and 10 `TempDir: mkdir` lines (V23).
  4. **[r3, quorum r2 fix applied verbatim]** A directly executed test binary has no default overall
     timeout, so the `-test.timeout=2m` deadline is required. A timeout (`panic: test timed out`) is an
     acceptance FAILURE, not an expected read-only-root failure.

  This proves that the variable AC-A sets is the one every caller's `TempDir` reads. It mirrors
  T2's in-test precondition fatal, and T2's `t.Setenv` stays unchanged.
- **AC-B.** M0: T2 is red against the base helper (prototype 10/10, V13).
- **AC-C.** Apply M1–M7 one at a time. Each reds exactly its named killers, and the file is
  restored byte-identical (r2 prototype sha256 `b8b1904a8d76…` before and after, V22).
- **AC-D.** `GOTMPDIR=/tmp/wNNNd/tmp-x134x go test ./host/broker/ -count=1` and the same with
  `-race` are green, with the 134 leg forced through `GOTMPDIR` as in AC-A (r2 prototype: plain
  `ok … 189.342s`, race `ok … 248.741s`; V24). Also the full
  gate: `go build ./... && go test ./...` and `./scripts/verify_ail.sh`, with no `.ail` touched.
  CI's run of the fixed head must show `ok  github.com/sunholo-data/ailang-world/host/broker` in both
  "go race gate" and "go host build + test gate".
- **No CI list edit.** No `tests=` list names a `TestExecTool*` test (V18). They run in both
  `./...` legs.

## 6. Files touched, conflict surface

| File | Change |
|---|---|
| `host/broker/exec_tools_test.go` | `f15RefusalProblem`; `requireF15` uses it (`:75-98`); T1, T2; T3's 4 lines in `TestExecToolExit134ChecksNode`; `bytes`, `fmt`, `io` imports |

`exec_tools_round2_test.go:357` calls `requireF15` and needs no edit (V3). No production file
changes, and nothing changes in `tools/launchd/*`, `.ail` or docs. **S2/S3/S1/S4/S5:** n/a (test
helper only). **S6:** the change removes a guard term that produced false reds and caught nothing, and T2
carries its own null-case precondition. **S7:** n/a (no usage surface). **Row 219 (attended, Codex):** this touches
only row 219 M3's test helper, in one file. Any open row-219 branch that edits `requireF15` or
`TestExecToolExit134ChecksNode` conflicts textually there and nowhere else.

## 7. Sibling sweep (bounded to this repo's Go tests)

**The hazard class:** a *match → fail* (absence) substring assertion, using a bare number or a
≤4-character token, over a reply that embeds `t.TempDir`/`os.MkdirTemp` paths. Presence
assertions (`!Contains → fail`) cannot red from a path. A path can only make them pass
vacuously. That is a different class, recorded as R2.

| Grep (V15–V17) | Hits | Control | Verdict |
|---|---|---|---|
| A: `strings.(Contains\|Count\|Index\|LastIndex\|HasSuffix\|HasPrefix)(…, "<digits>")` | 5 | `exec_tools_test.go:89` ✓ | **:89 is the only absence hit.** The other 4 are presence: `main_test.go:155` and `scripts/marketingcheck/inventory_test.go:124` (`"404"`), `outcome_test.go:226` (`"401"`), `doctor_hosts_test.go:106` (`"17"`) |
| C: the same functions with token `[A-Za-z0-9_.-]{1,4}` | 20 | `:89` ✓ | **Absence hits:** `:89` (hazard). `episode_report_test.go:137` `"a.go"` is safe: its haystack is `printReport(sampleReport())` with no temp path. `registry_reconcile_test.go:1022` `"url"`/`"uri"` is a match-conditioned branch, not an output assertion, and it is safe: it scans reflect field names. The `strings.Count(…,"x")` lines at `handlers_ailang_test.go:197` and `host/daemon/workspace_test.go:132` are safe: they read counter files that the stubs fill with `echo x` lines (`:156,163`; `workspace_test.go:104,109`). **The other 14 are presence** (`F5`, `F6`, `test`, `404`, `401`, `17`, `four`, `ok`, `.`×2, `1s`, `2s`, `cap`, `404`) |
| D: absence or presence with `fmt.Sprint[f]`/`strconv.Itoa` tokens | 12 | n/a (no expected hit; grep A's control shares the shape) | Safe: every token carries literal framing (`"entryIndex":%d`, `busy_timeout%%28%d%%29`, `>entry %d</text>`, …) that no temp path contains |
| E **[r2]**: `bytes.(Contains\|Index)(…, []byte("<digits>"))` | **0** | the same shape with any literal, `[]byte("[^"]+")`, gives 16, so the pattern matches real code (V25) | No bare-number `bytes` hits |
| F **[r2]**: variable-token absence, `strings.(Contains\|Index\|Count)(…, <identifier>)` not preceded by `!` | **170** (and 287 presence) | `:89` is not of this shape. The count shows the grep fires | **Not triaged.** 170 sites need context reads to know the token's length and whether the haystack holds temp paths. That is residual R3, not scope |

**Result [r2, glm's wording]:** there is **one true hit among literal-token `strings.*` absence
checks**, `:89`, which this doc fixes. `bytes.*` bare-number checks: 0. Variable-token checks are
not closed (R3). As a package-wide check, the whole `host/broker` package (home of the F15
refusals and `requireF15`) is green under a `GOTMPDIR` `134` root (V24). Other packages were not
run under a `134` root.

## 8. Residuals

- **R1. A TMPDIR containing a letter token.** On darwin, `$TMPDIR` is
  `/var/folders/<2>/<random alnum>/T/`, so a per-machine letter token in it is possible. No
  absence hit in the sweep has a temp-path haystack, so nothing here is exposed. I note it only
  because the sweep's class rests on it.
- **R2. Presence checks that a temp path can satisfy vacuously.** `doctor_hosts_test.go:106` `"17"`
  is the one whose haystack may quote temp files (`files[name]`). I did not verify this. The
  `404`/`401` checks read httptest error text, whose URLs have random ports. These are weak but
  never red. Candidate row: `w-test-presence-tokens-vacuous`. Out of scope.
- **R3 [r2]. Variable-token absence checks are not triaged.** Grep F finds 170 sites (V25). A site
  is exposed only if its token variable can be short, and a temp path can contain it, and its
  haystack quotes temp paths. Deciding that needs a per-site context read, which this small row
  does not cover. Candidate row: `w-test-variable-token-absence-sweep`. Interim evidence: only
  `host/broker` was run whole under a `134` root (V24), and that covers numeric `134` only, not
  the class.

## 9. Verification Log

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and tools | `cd <wt> && git log --oneline -1 && git status --short`; `go version` | `785b3bd1 Merge pull request #419 …`; no output from status; `go1.26.6 darwin/arm64` |
| V2 | The guard and its origin | `grep -n 134 host/broker/exec_tools_test.go`; `git blame -L89,89` | `:89 … strings.Contains(string(out), "134") …`; `44eed4376 (… 2026-10-29 14:55:18 +0100 89)`; `git log -1 44eed437` = `feat(row 219 M3): exec tools resolved per call (w219c, Codex)` |
| V3 | The six callers | `grep -n 'requireF15(' host/broker/*_test.go` + an awk map to the enclosing `func Test` | 9 call lines: `:124` DeletedTargetRefusesF15; `:142,143,149` FailedProbeRefusesF15; `:200` PlacementReapplied; `:336` CellarProfilePathRefusesByName; `:402,403` StartFailureIsRetriedWithoutChangingIdentity; `round2_test.go:357` SymlinkInsideWorkspaceRefuses |
| V4 | The CI reds | `gh run list --commit <sha>`; `gh run view <id> --log-failed \| grep -E -- '--- FAIL\|not a pre-launch'` | `171ece91` → 40909036384 `failure attempt=1`: race gate, `--- FAIL: TestExecToolCellarProfilePathRefusesByName`, `exec_tools_test.go:336: not a pre-launch refusal: {"ok":false,"refused":"exec /tmp/TestExecToolCellarProfilePathRefusesByName3490134876/001/Cellar/go/…`. `785b3bd1` → 40911404613 `failure attempt=2`: host build+test gate, `--- FAIL: TestExecToolSymlinkInsideWorkspaceRefuses`, `round2_test.go:357: not a pre-launch refusal: {… /tmp/TestExecToolSymlinkInsideWorkspaceRefuses596413420/001/node …`. 2 top-level `--- FAIL` lines across both logs |
| V5 | PR #418 | `gh pr view 418 --json mergedAt,mergeCommit,files --jq …` | `merged 2026-10-29T16:11:00Z`, oid `171ece91`, 19 files, including `exec_tools_test.go`, `_round2_test.go`, `_round3_test.go` |
| V6 | Reproduction | `TMPDIR=/tmp/w320d/<d> go test ./host/broker/ -run '^(<six>)$' -count=1 -v` | `tmp-x134x`: rc=1, PASS 0, FAIL 6 (`/private/tmp/w320d/tmp-x134x/…`). `tmp-ctl`: rc=1, PASS 5, FAIL 1, on `…/TestExecToolDeletedTargetRefusesF154134186054/…` (a natural hit) |
| V7 | Natural flake rate | the same, `tmp-ctl`, `-count=50` | rc=1, PASS 295, FAIL 5. Failing parents: `…ByName1349870625`, `…ByName2138134077`, `…ByName2575513453`, `…F152134511095`, `…nodeworkspace4251347714` |
| V8 | The two real exit-134 shapes | scratch test (`/tmp/w320d/scratch`, the `:230` fixture body, logging keys, spawns and `Contains "134"`), `-run TestW320Probe -v` | broken-node: `spawns 0->1 keys=[ok refused] has134=false`, `refused="exec node …/node does not start: dyld: Library not loaded: libada.3.dylib (--version failed)…"`. command-aborted: `spawns 0->1 keys=[argv duration_ms exit_code limit profile sandbox stderr … timed_out] exit_code=134` |
| V9 | Production shapes | `sed -n 300,345p host/broker/handlers_exec.go`; `grep -n 'type execResult' -A3`; `sed -n 95,97p` | `:322 h.spawns.Add(1)`, `:323 runBounded(…)`, `:326 … code == 134 …`, `:333 return execRefusal(h.tools.nodeFailed(sandbox, why))`; `:209 ExitCode *int \`json:"exit_code"\``; `:97 json.Marshal(map[string]any{"ok": false, "refused": why})` |
| V10 | `t.TempDir` path shape | read `func (c *common) makeTempDir` in `$(go env GOROOT)/src/testing/testing.go` | "a single parent directory for all the temporary directories created by a test"; `pattern := c.Name()` … `os.MkdirTemp(os.Getenv("GOTMPDIR"), pattern)`; `c.tempDirSeq++` |
| V11 | (r1, `TMPDIR` form; superseded for the ACs by V21) Prototype green | `/tmp/w320d/scratch`: fixed helper + T1–T3; `TMPDIR=<d> go test ./host/broker/ -run '^(<six>\|T1\|T2\|TestExecToolExit134ChecksNode)$' -count=20 -v` | `tmp-x134x rc=0 PASS 180 FAIL 0`; `tmp-ctl rc=0 PASS 180 FAIL 0` |
| V12 | Mutants | `/tmp/w320d/mut.py <file> M<n>` (one exact-string replace, `assert count==1`), then `-run '^(T1\|T2\|Exit134)$' -count=1 -v`, then restore | M1: T1 `refusal-temp-root-134`, T2. M2: T1 `post-launch-f15`, Exit134/`broken-node`. M3: T1 `refusal-plus-exit-code`. M4: T1 `ok-true`. M5: T1 `refused-empty`. M6: T1 (5 rows), Exit134/`broken-node` + `command-aborted`. All rc=1. sha256 `0cc5c2118d86` before and after |
| V13 | M0: the base helper reds T2 | the base `requireF15` spliced in from the worktree; `-run '^TestExecToolF15RefusalUnderTempRoot134$' -count=10 -v` | rc=1, FAIL 10, PASS 0; `not a pre-launch refusal: {… /private/tmp/w320d/tmp-ctl/w234-x134x-1772733125/TestExecToolF15RefusalUnderTempRoot1342635915153/001/node …` |
| V14 | (r1, `TMPDIR` form; superseded for AC-D by V24) Whole package, fixed | `/tmp/w320d/full` (full `git archive HEAD` + fixed file), `AILANG_BIN=…/.pinned-ailang/ailang`: `go test -race ./host/broker/ -count=1`; `TMPDIR=/tmp/w320d/tmp-x134x go test ./host/broker/ -count=1` | race rc=0, `ok … 249.882s`; plain with a `134` TMPDIR rc=0, `ok … 189.035s`. (A partial copy without `docs/` and `AILANG_BIN` gave unrelated reds such as `open ../../docs/QUICKSTART.md`, so I discarded it.) |
| V15 | Sweep A | `grep -rnE 'strings\.(Contains\|Count\|Index\|LastIndex\|HasSuffix\|HasPrefix)\([^;]*, *"[0-9]+"\)' --include='*_test.go' . \| grep -v /.claude/ \| wc -l` | `5`, including the control `host/broker/exec_tools_test.go:89` |
| V16 | Sweep C | `grep -rnoE '(!?)strings\.(Contains\|Count\|Index\|LastIndex)\([^()]*(\([^()]*\))?[^()]*, *"[A-Za-z0-9_.-]{1,4}"\)' --include='*_test.go' . \| grep -v /.claude/`; context read with `sed -n` on each absence hit | `20`, including the control `:89`. Absence hits are as in §7 (`episode_report_test.go:131` `printReport(&out, rep)` over `sampleReport()`; `handlers_ailang_test.go:156` `echo x >> "$d/dispatches"`) |
| V17 | Sweep D | `grep -rnE 'strings\.(Contains\|Count\|Index)\([^()]*(\([^()]*\))?[^()]*, *(strconv\.Itoa\|fmt\.Sprint[f]?)\(' --include='*_test.go' .` | 12 hits, all with literal framing (for example `agui_test.go:717` `"entryIndex":%d`, `context_read_test.go:229` `busy_timeout%%28%d%%29`) |
| V18 | CI env and lists; the pin | `grep -nE 'GOTMPDIR\|TMPDIR' .github/workflows/*.yml`; `grep -nE 'TestExecTool\|requireF15' .github/workflows/ci.yml docs/QUICKSTART.md`; `grep -n 'PINNED' ci.yml`; `~/.pinned-ailang/ailang --version` | no output from the first two (control: `ci.yml:245` `tests="TestExecSrtGoTest …` matches `tests=`); `:77 Install PINNED v0.41.0 ailang`; `AILANG v0.41.0` |
| V19 | Row 234 is not in the charter at base | `grep -cE '^\| *234 *\|' design_docs/world-mission.md` | `0`. Control: the same file's `:1021` `\| D-WORLD-102 \| OPEN \|` and STATUS `:1031` (iteration 320) are present |
| V20 | **[r2]** `Decode` reads only the first value (gpt6-1-sol's objection, measured) | `/tmp/w320d/dec/main.go`: the r1 decode (`DisallowUnknownFields`, `*bool`/`*string`), then a second `dec.Decode(&json.RawMessage{})`, over three inputs; `go run .` | `{…}{"exit_code":134}` → `first=<nil> second=<nil>`. `{…} garbage` → `first=<nil> second=invalid character 'g' looking for beginning of value`. `{…} \n\t` → `first=<nil> second=EOF`. So r1 accepted the first two |
| V21 | **[r2]** AC-A with the root forced through `GOTMPDIR` (r2 helper, sha256 `b8b1904a8d76`) | `/tmp/w320d/scratch`: `GOTMPDIR=/tmp/w320d/<d> go test ./host/broker/ -run '^(<six>\|T1\|T2\|TestExecToolExit134ChecksNode)$' -count=20 -v`; count `^--- PASS`, `^--- FAIL`, `grep -c <d>`, unique PASS names | `tmp-x134x`: rc=0, PASS 180, FAIL 0, root-lines **0**, names 9. `tmp-ctl`: rc=0, PASS 180, FAIL 0, root-lines **0**, names 9. Passing `-v` output prints no temp paths, so glm's `grep -c` = PASS form cannot be met |
| V22 | **[r2]** Mutants M1–M7 on the r2 helper | as V12, with `GOTMPDIR=/tmp/w320d/tmp-ctl`; M7 replaces `err != io.EOF {` with `false && err != io.EOF {`. A first M7 form that deleted the block failed to build (`"fmt" imported and not used`), so I replaced it with this one | M1: T1 `refusal-temp-root-134`, `trailing-whitespace`; T2. M2: T1 `post-launch-f15`; Exit134/`broken-node`. M3: `refusal-plus-exit-code`. M4: `ok-true`. M5: `refused-empty`. M6: T1 (`post-launch-f15`, `launched-exit-134`, `refusal-plus-exit-code`, `ok-true`, `refused-empty`, `concatenated-json`, `trailing-garbage`); Exit134 both subtests. **M7: only T1 `concatenated-json`, `trailing-garbage`.** All rc=1. sha256 `b8b1904a8d76` before and after |
| V23 | **[r2]** `GOTMPDIR` dominance, and the root-honoured control | at base (worktree): `GOTMPDIR=<x134x> TMPDIR=<ctl> go test … -run '^(<six>)$' -v` and the reverse. On the r2 prototype: `GOTMPDIR=/tmp/w320d/tmp-x134x-ro go test ./host/broker/ -run '^(<six>)$' -count=1 -v`; then `go test -c -o /tmp/w320d/broker.test ./host/broker/` and, from `host/broker/`, `GOTMPDIR=/tmp/w320d/tmp-x134x-ro /tmp/w320d/broker.test -test.run '^(<six>)$' -test.v` (root mode 555) | base `GOTMPDIR=x134x TMPDIR=ctl`: rc=1, PASS 0, FAIL 6, 9 `tmp-x134x` lines. Base `GOTMPDIR=ctl TMPDIR=x134x`: rc=0, **PASS 6**, FAIL 0, 0 lines (glm's vacuity, confirmed). `go test` with the ro root: `go: creating work dir: mkdir /tmp/w320d/tmp-x134x-ro/go-build…: permission denied`. Prebuilt binary: rc=1, top-level FAIL 6, all FAIL 11 (10 leaves + the `PlacementReapplied` parent), `TempDir: mkdir /tmp/w320d/tmp-x134x-ro` lines **10**; first `testing.go:1426: TempDir: mkdir /tmp/w320d/tmp-x134x-ro/TestExecToolSymlinkInsideWorkspaceRefuses2299717663: permission denied` **[r3, controller re-run with the exact AC-A-ROOT command]** `go test -c -o /tmp/w320/ctl-r3/broker.test ./host/broker/` (rc 0) at `785b3bd1`; from `host/broker/`, `GOTMPDIR=/tmp/w320/ctl-r3/tmp-x134x-ro` (mode 555) `broker.test -test.run '^(<six>)$' -test.v -test.timeout=2m` → exit 1, elapsed 1 s, 6 top-level `--- FAIL`, 11 `--- FAIL` lines in all (10 leaves + the `TestExecToolPlacementReapplied` parent), `grep -c 'TempDir: mkdir <root>'` → 10, `panic: test timed out` → 0; control `=== RUN` → 11. |
| V24 | **[r2]** AC-D with the `134` root forced through `GOTMPDIR` (r2 helper) | `/tmp/w320d/full` (full `git archive HEAD` + the r2 file), `GOTMPDIR=/tmp/w320d/tmp-x134x AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/broker/ -count=1`, then the same with `-race` | plain rc=0, `ok … 189.342s`; race rc=0, `ok … 248.741s` |
| V25 | **[r2]** Sweeps E and F | E: `grep -rnE 'bytes\.(Contains\|Index)\([^;]*, *\[\]byte\("[0-9]+"\)\)' --include='*_test.go' . \| grep -v /.claude/ \| wc -l`; its control, the same with `"[^"]+"`. F: `grep -rnoE '(^\|[^!])strings\.(Contains\|Index\|Count)\([^()]*(\([^()]*\))?[^()]*, *[A-Za-z_][A-Za-z0-9_.]*\)' … \| wc -l`, and the `!`-prefixed twin | E: `0`; control `16`. F: `170` absence-shaped; `287` presence-shaped |

## 10. Quorum log

**Round 1** (artifact `w-exec-refusal-check-134-substring-2026-10-29T17-22-10Z.json`, $0.1694):
**BLOCKED**.

Seats: `gpt6-1-sol` REJECT, `gemini-4-argon` PASS, `oc-glm-5-3` REJECT. The reserve was
`oc-kimi-k3`. The author's vendor (claude) was benched.

- **gpt6-1-sol (REJECT): envelope completeness.** Quoted: "json.Decoder.Decode reads only the
  first JSON value; a valid refusal followed by another JSON object or trailing garbage is
  accepted when before == after, so the 'strict envelope' claim is stronger than the
  implementation." **Measured: correct** (V20). **r2 applies their fix verbatim:**
  - after the first `Decode`, a second `Decode` into a temporary value must return `io.EOF`,
    with a trailing-content error otherwise, and the `io` import (§2);
  - T1 rows `concatenated-json` and `trailing-garbage` (rejected) and the control
    `trailing-whitespace` (accepted) (§5);
  - mutant M7 MUT-NO-EOF-CHECK, killed by exactly those two rows (V22);
  - the claim now reads "exactly one complete JSON value" (§2, §3).
- **gemini-4-argon (PASS).** No objection is carried into r2.
- **oc-glm-5-3 (REJECT): AC hermeticity and sweep scope.** Quoted: "the ACs force the 134 root
  through TMPDIR only, but per your own V10 GOTMPDIR is read first, so an ambient GOTMPDIR
  silently makes AC-A and AC-D's 134 leg vacuous." **Measured: correct.** At base,
  `GOTMPDIR=<ctl> TMPDIR=<x134x>` passes all six (V23). **r2 applies:**
  - AC-A and AC-D's `134` leg now set `GOTMPDIR` (V21, V24), and T2's `t.Setenv` is unchanged.
  - **One part of their fix was measured as not achievable as written, and is replaced by its
    measured equivalent.** The fix said: "a root-honoured assertion that grep -c tmp-x134x on
    the -v output matches the PASS count." Passing tests print no paths, so that count is 0
    against 180 PASS (V21). AC-A-ROOT asserts the same fact in a form that can be run. A
    prebuilt test binary under a read-only `GOTMPDIR` root must give one
    `TempDir: mkdir <root>` line per leaf FAIL (10 = 10), plus 6 top-level FAILs (V23). The go
    tool also writes its work dir under `GOTMPDIR`, so a prebuilt binary is required.
  - On sweep scope, their choice set was "either rewrite … or extend the sweep … before
    claiming the class is closed". r2 does both, as far as is measurable:
    - the §7 result sentence uses their wording, "one true hit among literal-token `strings.*`
      absence checks";
    - the `bytes` bare-number grep is run: 0 hits, control 16 (V25);
    - the variable-token pass is run and counted at 170 hits, but not triaged. It is filed as R3
      instead of being claimed closed.

**Round 2** (artifact `w-exec-refusal-check-134-substring-2026-10-29T17-40-31Z.json`, same seats, $0.1535): synthesis BLOCKED. `gpt6-1-sol` **reject**, surface: AC-A-ROOT hermeticity — the prebuilt `broker.test` run had no `-test.timeout`, an unbounded wait (a new surface, introduced by the r2 revision). `gemini-4-argon` **pass** (duplicate keys inside one object are still accepted by `encoding/json`; production replies come from one `json.Marshal`, so they cannot carry them — non-blocking). `oc-glm-5-3` **pass** (the class is provably closed for `host/broker` only; R3's 170 variable-token sites are untriaged — non-blocking, already residual R3). No absent reviewers.

**Round 3 — narrow-refinement carve-out (controller, iteration 320).** After the one re-quorum, the only remaining blocking objection carries a concrete reviewer-authored `proposed_fix` and does not dispute the design direction (it is a determinism/boundedness fix to one acceptance command). The controller applied the fix verbatim: `-test.timeout=2m` added to AC-A-ROOT step 2; step 4 states a timeout is an acceptance failure; V23 re-run with that exact command (elapsed 1 s, exit 1, 6 top-level and 10 leaf FAILs, 10 `TempDir: mkdir` errors, 0 timeout panics). Routed to the planner without a further quorum round. Reviewer text, verbatim: "Add an explicit deadline to AC-A-ROOT: GOTMPDIR=<root>-ro /tmp/wNNNd/broker.test -test.run '^(<six>)$' -test.v -test.timeout=2m. Update V23 with a rerun using that exact command, recording elapsed time, exit status, the expected six top-level failures and ten TempDir errors, and absence of a timeout panic. State that a timeout is an acceptance failure, not an expected read-only-root failure."
