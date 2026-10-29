# Sprint plan — `w-exec-refusal-check-134-substring` (iteration 320, row 234, red on `dev`)

**Design doc:** [`design_docs/planned/w-exec-refusal-check-134-substring.md`](w-exec-refusal-check-134-substring.md)
(r2 + the quorum-r2 fix applied by the controller as round 3; §10 Quorum log). **Sprint base:**
worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter320-f15`, branch
`sprint/world-iter320-f15-tempdir`, HEAD `5407c160`. Its Go tree equals the design's measurement
base `785b3bd1`: `git diff --stat 785b3bd1 HEAD` lists only 4 `design_docs/` files (V-P1).
**Planner:** mission-control iteration 320 (planner role). **Executor:** claude-haiku, so every edit
is spelled out below. **Copy code from the fenced blocks, never from table cells.**
**Estimate:** about 0.1 day. **Two milestones**, each independently green. **Risk:** low (test code
only). The planner prototyped the whole change, every acceptance command and every mutant in a
scratch copy outside the repo (`/tmp/w320p/full` and `/tmp/w320p/mut`, made with `git archive HEAD`);
the results are in §9 and §10.

**Scope (the only files you may edit or create):**
- `host/broker/exec_tools_test.go` (M1; the only code file)
- `design_docs/verification/world-iter320/executor-report.md` (M2; new)
- `design_docs/verification/world-iter320/executor-*.txt` (M2; new evidence logs)

**No production `.go` file. No `.ail` file. No CI file. No `docs/` file. No `tools/launchd/*`.**
**Do NOT edit `design_docs/world-mission.md`, `design_docs/world-mission-log.md`,
`design_docs/world-mission-status-archive.md` or any `design_docs/world-mission*.md`.** Queue row 234
does not exist in the charter at base (`grep -cE '^\| *234 *\|' design_docs/world-mission.md` → `0`,
V-P9). The **controller** writes row 234's charter/queue text, from the design doc's §1 (the defect)
and §2 (the decision) plus this sprint's executor report. The executor never touches it. Do not move
the design doc from `planned/` to `implemented/`; the controller does that too.

---

## 0. The one-paragraph version

`requireF15` (`host/broker/exec_tools_test.go:75`) asserts a call was refused before launch. Its
guard (`:89`) includes `strings.Contains(string(out), "134")`, meaning "the call launched and exited
134". But every F15 refusal quotes the tool's `t.TempDir()` paths, and those paths contain `134`
whenever the random suffix, the temp root or the test name does: a ~1.7 % per-test flake, red twice on
`dev` CI. The term never caught anything (design §1, V8/V9). The fix replaces the whole guard with a
pure helper `f15RefusalProblem(out, before, after)`: a strict decode of **exactly one** JSON value into
the two-key envelope `{"ok":false,"refused":"…"}`, plus the spawn counter. Tests:
- **T1** `TestF15RefusalProblemIsStructural`: a pure 9-row table over the helper.
- **T2** `TestExecToolF15RefusalUnderTempRoot134`: a live F15 refusal under a forced `GOTMPDIR` root
  holding `x134x` (red at base = M0 / AC-B).
- **T3**: two assertions added to the existing `TestExecToolExit134ChecksNode` (the two real
  exit-134 launch shapes must be rejected by the helper).

Milestones:
- **M1:** T2 red-first at base (AC-B), then the fix, T1, T3; AC-A, AC-A-ROOT; the M1–M7 mutation drill (AC-C).
- **M2:** AC-D (whole `host/broker` under a `134` root, plain and `-race`; full gate; `verify_ail.sh`)
  and the executor report with evidence logs.

---

## 1. Environment (on EVERY command)

```sh
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter320-f15
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang     # reports "AILANG v0.41.0" (measured, V-P2)
unset GOTMPDIR                                    # only the AC commands set it, inline
mkdir -p /tmp/w320x                               # the executor's scratch dir (the design's /tmp/wNNNd)
```

| Pin | Measured this session | Why it matters |
|---|---|---|
| `go` | `go version go1.26.6 darwin/arm64` | The CI pin (`go.mod`) |
| `AILANG_BIN` | `AILANG v0.41.0` (the CI pin, `ci.yml:77`) | Without it daemon e2e tests SKIP and `verify_ail.sh` cannot run |
| `GOTMPDIR` | unset in the ambient shell; `TMPDIR=/var/folders/kr/…/T/` | `testing` reads `GOTMPDIR` first (design V10). An ambient value would make the AC legs vacuous (design V23) |
| `WORLD_EXEC_SRT_NODE_MODULES` | unset | Real-srt tests SKIP locally. Do NOT set it |

**Instrument rule.** Never write `cmd | tail; echo rc=$?`. Redirect to a file, put `echo "rc=$?"` on
the next line, then read the file. A `-run` that matches nothing exits 0, so read every verdict from
counted `--- PASS:` / `--- FAIL:` lines.

**Compile fence.** `go build ./...` does **not** compile `_test.go` files. The compile fence is
`go test -run '^$' ./host/broker/` (or `go vet ./host/broker/`).

**Executor constraints (binding).**
1. **No git write operations of any kind:** no `commit`, `add`, `stash`, `checkout -- file`,
   `restore` or `reset`. The controller commits. `git diff`, `git status` and `git show` are allowed.
2. Edit only the *Scope* files, and only at the sites this plan names. Snapshots and logs go under
   `/tmp/w320x/`, **never** inside the worktree (no `.snap/` dir in the repo).
3. Never edit `tools/launchd/*`. Never touch the main checkout
   `/Users/voightkampff/dev/sunholo-data/ailang-world` or `~/.ailang/state/mission-v1*`. Never edit
   `design_docs/world-mission*.md` (the controller writes row 234). Do not run `ailang messages`.
4. Apply each mutation **alone**, run **only** its named killer command (§M1.6), see the named rows
   red, then restore from the snapshot and check with `shasum -a 256 -c`.
5. Do not add `t.Parallel()` (T2 uses `t.Setenv`, which panics under `t.Parallel`).
6. **Bound every wait**: wrap long runs in `perl -e 'alarm N; exec @ARGV' …`. Never end your turn
   announcing a wait.

---

## 2. Gate list

| ID | Command | Pass |
|---|---|---|
| G1 | `gofmt -l host/broker/` | empty output |
| G2 | `go vet ./host/broker/` | rc 0 |
| G3 | `go test -run '^$' ./host/broker/` | rc 0, `ok … [no tests to run]` (**the compile fence**) |
| G4 | AC-A (§M1.4) | both roots: rc 0, 180 `^--- PASS`, 320 `^    --- PASS`, 0 `--- FAIL`, 9 names |
| G5 | AC-A-ROOT (§M1.5) | rc 1, 6 top-level FAIL, 10 leaf FAIL = 10 `TempDir: mkdir` lines, 0 timeout panics |
| G6 | AC-D (§M2.1) | plain and race `ok … host/broker`; full gate rc 0; `verify_ail.sh` rc 0 |
| G7 | `git diff --name-only` | after M1: exactly `host/broker/exec_tools_test.go`; after M2 also the new files (`git status --porcelain`) |

Copy every `-run` regex from the fenced code blocks. In a markdown table a `\|` renders as `|`.

---

## 3. Pristine-tree baselines (record these BEFORE the first edit)

```sh
shasum -a 256 host/broker/exec_tools_test.go host/broker/exec_tools_round2_test.go host/broker/handlers_exec.go > /tmp/w320x/base.sha
cat /tmp/w320x/base.sha
cp host/broker/exec_tools_test.go /tmp/w320x/exec_tools_test.go.base
git status --porcelain
```

Expected (planner, this session):

| File | sha256 |
|---|---|
| `host/broker/exec_tools_test.go` | `20179c17a1fb9c3e4949d01ecc2c618b6e0a809d2efd81705903f239a589fe63` |
| `host/broker/exec_tools_round2_test.go` | `fe2deb2a078c5c62600776c9c4a83abd63d2ad56c493179cf0b057dcb9c272d7` (must stay unchanged) |
| `host/broker/handlers_exec.go` | `d85457256d78add86ac475fc0eb77bcf8cc373082c27139ba3d9f7fd023b2ad0` (must stay unchanged) |

`git status --porcelain` at base shows only this plan and its JSON (the controller commits them).
If any sha256 differs, **stop** and report: the tree is not the planned base.

### 3.1 Base RED / GREEN

- **Red at base (deterministic):** T2 against the base `requireF15` (M0, AC-B), 10/10.
- **Red at base (flaky, ~1.7 %/test):** the six callers, whenever a temp path holds `134`; 6/6 red under
  a forced `134` root (design V6, V23).
- **Cannot compile at base:** T1 and T3's new assertions (they call `f15RefusalProblem`).
- **Green at base and must stay green:** the six `requireF15` callers under a control root,
  `TestExecToolExit134ChecksNode`, the whole `host/broker` package.

The six callers (design V3; `grep -n 'requireF15(' host/broker/*_test.go`, re-run V-P3):

| Test | Call lines |
|---|---|
| `TestExecToolDeletedTargetRefusesF15` | `exec_tools_test.go:124` |
| `TestExecToolFailedProbeRefusesF15` | `:142`, `:143`, `:149` |
| `TestExecToolPlacementReapplied` (4 subtests `node/workspace` … `argv0/state` call it; `startup` does not) | `:200` |
| `TestExecToolCellarProfilePathRefusesByName` | `:336` |
| `TestExecToolStartFailureIsRetriedWithoutChangingIdentity` | `:402`, `:403` |
| `TestExecToolSymlinkInsideWorkspaceRefuses` | `exec_tools_round2_test.go:357` (no edit needed) |

---

## 4. Milestones

### Milestone M1 — T2 red-first, `f15RefusalProblem`, `requireF15` rewired, T1, T3, AC-A, AC-A-ROOT, mutation drill (+91/−4, one file)

**File:** `host/broker/exec_tools_test.go` only. Do the edits **in this order**.

#### M1.1 Red first: append T2 to the pristine file, run AC-B (= M0)

Append this block at the **end of the file**, after the closing `}` of
`TestExecToolArchiveDriftPrecedesReacceptance` (the last function, base `:445-454`):

```go
// T2 (row 234): the six requireF15 callers' refusals quote t.TempDir paths, so a temp root
// holding "134" must not red a real pre-launch refusal. The root is forced through GOTMPDIR
// (testing reads it first, one parent per test); the precondition makes the hazard real.
func TestExecToolF15RefusalUnderTempRoot134(t *testing.T) {
	root, err := os.MkdirTemp("", "w234-x134x-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(root) })
	t.Setenv("GOTMPDIR", root)
	m := newMovingExec(t, execFixtureOpts{})
	if !strings.Contains(m.link, "x134x") {
		t.Fatalf("temp root not honoured: %s is not under %s", m.link, root)
	}
	if err := os.Remove(m.a); err != nil {
		t.Fatal(err)
	}
	requireF15(t, m.f, m.link, m.a, "the target is gone")
}
```

It compiles at base (it needs only `os`, `strings`, `testing`, `newMovingExec` and `requireF15`).
Then run **AC-B** (design §5: "M0: T2 is red against the base helper"):

```sh
gofmt -l host/broker/                                   # empty
go vet ./host/broker/; echo "rc=$?"                     # rc 0
perl -e 'alarm 300; exec @ARGV' go test ./host/broker/ -run '^TestExecToolF15RefusalUnderTempRoot134$' -count=10 -v > /tmp/w320x/acb-m0.txt 2>&1; echo "rc=$?"
grep -c '^--- FAIL: TestExecToolF15RefusalUnderTempRoot134 ' /tmp/w320x/acb-m0.txt   # want 10
grep -c '^--- PASS' /tmp/w320x/acb-m0.txt                                          # want 0
grep -c 'temp root not honoured' /tmp/w320x/acb-m0.txt                             # want 0 (the precondition held)
grep -m1 'not a pre-launch refusal' /tmp/w320x/acb-m0.txt | cut -c1-240            # path holds w234-x134x-
```

Planner's prototype: rc 1, 10 FAIL, 0 PASS; first line
`exec_tools_test.go:473: not a pre-launch refusal: {"ok":false,"refused":"exec node /private/var/folders/…/T/w234-x134x-2539520704/TestExecToolF15RefusalUnderTempRoot1342643642204/001/node → …` (V-P4).
If T2 PASSes here, **stop**: the hazard was not built.

#### M1.2 The fix: four edits

**(a) Imports** (base `:3-13`). Replace

```go
import (
	"context"
	"encoding/json"
	"errors"
	"os"
```

with

```go
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
```

**(b) The predicate.** Insert this block **directly above** the line `func requireF15(t *testing.T, f *execFixture, parts ...string) string {` (base `:75`), with one blank line kept between the end of `newMovingExec` and the comment:

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
	if err := dec.Decode(&r); err != nil {
		return "not the refusal envelope: " + err.Error()
	}
	if err := dec.Decode(&json.RawMessage{}); err != io.EOF {
		return fmt.Sprintf("trailing content after the envelope: %v", err)
	}
	if r.OK == nil || *r.OK {
		return "ok is not false"
	}
	if r.Refused == nil || *r.Refused == "" {
		return "refused is empty"
	}
	if after != before {
		return "the call spawned"
	}
	return ""
}
```

Keep every `if` on its own lines exactly as shown: the mutants in §M1.6 replace those exact lines.

**(c) `requireF15`'s guard** (base `:82-91`). Replace these 10 lines:

```go
	var refusal struct {
		OK      bool   `json:"ok"`
		Refused string `json:"refused"`
	}
	if err := json.Unmarshal(out, &refusal); err != nil {
		t.Fatal(err)
	}
	if refusal.OK || refusal.Refused == "" || strings.Contains(string(out), "134") || f.h.Spawns() != before {
		t.Fatalf("not a pre-launch refusal: %s (spawns %d -> %d)", out, before, f.h.Spawns())
	}
```

with:

```go
	after := f.h.Spawns()
	if problem := f15RefusalProblem(out, before, after); problem != "" {
		t.Fatalf("not a pre-launch refusal (%s): %s (spawns %d -> %d)", problem, out, before, after)
	}
	var refusal struct {
		Refused string `json:"refused"`
	}
	if err := json.Unmarshal(out, &refusal); err != nil {
		t.Fatal(err)
	}
```

The `parts` loop and `return refusal.Refused` below it stay unchanged. The failure message keeps
the CI signature's prefix `not a pre-launch refusal`.

**(d) T3, in `TestExecToolExit134ChecksNode`** (base `:243-248`). Replace:

```go
			f := newExecFixture(t, execFixtureOpts{node: node})
			out, err := f.call(t, "test", "a")
			if err != nil {
				t.Fatal(err)
			}
			if fail {
```

with:

```go
			f := newExecFixture(t, execFixtureOpts{node: node})
			before := f.h.Spawns()
			out, err := f.call(t, "test", "a")
			if err != nil {
				t.Fatal(err)
			}
			if f15RefusalProblem(out, before, f.h.Spawns()) == "" {
				t.Fatalf("a launch that exited 134 passes requireF15: %s", out)
			}
			if fail {
```

The existing `if fail { … } else if … ExitCode …` checks below stay unchanged.

#### M1.3 T1: append at the end of the file (after T2)

```go
// T1 (row 234): f15RefusalProblem judges structure and the spawn counter, never text. The
// accepted rows quote temp paths holding "134"; every reject row fails on one property.
func TestF15RefusalProblemIsStructural(t *testing.T) {
	env := `{"ok":false,"refused":"exec /tmp/x134x/TestExecToolCellarProfilePathRefusesByName1349870625/001/node does not start: the target is gone"}`
	for _, row := range []struct {
		name          string
		out           string
		before, after int64
		accepted      bool
	}{
		{"refusal-temp-root-134", env, 0, 0, true},
		{"trailing-whitespace", env + " \n\t", 0, 0, true},
		{"post-launch-f15", env, 0, 1, false},
		{"launched-exit-134", `{"exit_code":134,"timed_out":false,"limit":null,"duration_ms":3,"stdout":"","stderr":"/tmp/x134x/node: Abort trap: 6"}`, 0, 1, false},
		{"refusal-plus-exit-code", `{"ok":false,"refused":"exec /tmp/x134x/node does not start","exit_code":134}`, 0, 0, false},
		{"ok-true", `{"ok":true,"refused":"exec /tmp/x134x/node does not start"}`, 0, 0, false},
		{"refused-empty", `{"ok":false,"refused":""}`, 0, 0, false},
		{"concatenated-json", env + `{"exit_code":134}`, 0, 0, false},
		{"trailing-garbage", env + " garbage", 0, 0, false},
	} {
		t.Run(row.name, func(t *testing.T) {
			problem := f15RefusalProblem([]byte(row.out), row.before, row.after)
			if (problem == "") != row.accepted {
				t.Fatalf("problem %q, want accepted=%v: %s", problem, row.accepted, row.out)
			}
		})
	}
}
```

Row map (design §5 T1; the literal JSON is the planner's, D1):

| Row | Input | Spawns | Want | Property it pins |
|---|---|---|---|---|
| `refusal-temp-root-134` | envelope; paths hold `x134x` and `…1349870625` | 0→0 | accepted | no path can red it (M1) |
| `trailing-whitespace` | envelope + `" \n\t"` | 0→0 | accepted | trailing whitespace is allowed (control for M7; also sees M1) |
| `post-launch-f15` | envelope | 0→1 | rejected | the spawn counter (M2) |
| `launched-exit-134` | an `execResult`-shaped reply with `exit_code:134`, no `ok`/`refused` | 0→1 | rejected | (M6 only) |
| `refusal-plus-exit-code` | envelope keys + `"exit_code":134` | 0→0 | rejected | strict decode (M3) |
| `ok-true` | `{"ok":true,"refused":"…"}` | 0→0 | rejected | `ok` must be false (M4) |
| `refused-empty` | `{"ok":false,"refused":""}` | 0→0 | rejected | `refused` non-empty (M5) |
| `concatenated-json` | envelope + `{"exit_code":134}` | 0→0 | rejected | exactly one JSON value (M7) |
| `trailing-garbage` | envelope + `" garbage"` | 0→0 | rejected | exactly one JSON value (M7) |

#### M1.4 Acceptance (M1): gates, then AC-A

```sh
gofmt -l host/broker/                                   # G1: empty
go vet ./host/broker/; echo "rc=$?"                     # G2: rc 0
go test -run '^$' ./host/broker/ > /tmp/w320x/g3.txt 2>&1; echo "rc=$?"; cat /tmp/w320x/g3.txt   # G3: rc 0
git diff --stat                                         # exec_tools_test.go | 95 (+91 -4)
```

**AC-A** (design §5, the root forced through `GOTMPDIR`; `/tmp/wNNNd` = `/tmp/w320x`). The `-run`
regex is the design's `'^(<six>|TestF15RefusalProblemIsStructural|TestExecToolF15RefusalUnderTempRoot134|TestExecToolExit134ChecksNode)$'` with `<six>` expanded:

```sh
RUN9='^(TestExecToolDeletedTargetRefusesF15|TestExecToolFailedProbeRefusesF15|TestExecToolPlacementReapplied|TestExecToolCellarProfilePathRefusesByName|TestExecToolStartFailureIsRetriedWithoutChangingIdentity|TestExecToolSymlinkInsideWorkspaceRefuses|TestF15RefusalProblemIsStructural|TestExecToolF15RefusalUnderTempRoot134|TestExecToolExit134ChecksNode)$'
for d in tmp-x134x tmp-ctl; do
  rm -rf /tmp/w320x/$d; mkdir -p /tmp/w320x/$d
  GOTMPDIR=/tmp/w320x/$d perl -e 'alarm 900; exec @ARGV' go test ./host/broker/ -run "$RUN9" -count=20 -v > /tmp/w320x/aca-$d.txt 2>&1; echo "$d rc=$?"
  f=/tmp/w320x/aca-$d.txt
  echo "$d top-PASS=$(grep -c '^--- PASS' $f) sub-PASS=$(grep -c '^    --- PASS' $f) FAIL=$(grep -c -- '--- FAIL' $f) SKIP=$(grep -c -- '--- SKIP' $f) names=$(grep '^--- PASS' $f | awk '{print $3}' | sort -u | wc -l)"
done
```

Want, **for each root**: rc 0, `top-PASS=180` (9 names × 20), `sub-PASS=320` ((5 + 2 + 9) subtests
× 20), `FAIL=0`, `SKIP=0`, `names=9`. Planner's prototype: exactly these numbers for both roots,
43 s and 39 s (V-P5). (`tmp-ctl`'s path holds no `134`; the `-v` output of passing tests prints no
paths, so do not grep it for the root; AC-A-ROOT proves the root is honoured.)

#### M1.5 AC-A-ROOT (design §5, steps 1–4, `-test.timeout=2m` required)

```sh
perl -e 'alarm 300; exec @ARGV' go test -c -o /tmp/w320x/broker.test ./host/broker/; echo "build rc=$?"   # rc 0
rm -rf /tmp/w320x/tmp-x134x-ro; mkdir -p /tmp/w320x/tmp-x134x-ro && chmod 555 /tmp/w320x/tmp-x134x-ro
SIX='^(TestExecToolDeletedTargetRefusesF15|TestExecToolFailedProbeRefusesF15|TestExecToolPlacementReapplied|TestExecToolCellarProfilePathRefusesByName|TestExecToolStartFailureIsRetriedWithoutChangingIdentity|TestExecToolSymlinkInsideWorkspaceRefuses)$'
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter320-f15/host/broker && s=$(date +%s); GOTMPDIR=/tmp/w320x/tmp-x134x-ro /tmp/w320x/broker.test -test.run "$SIX" -test.v -test.timeout=2m > /tmp/w320x/acaroot.txt 2>&1; echo "rc=$? elapsed=$(( $(date +%s) - s ))s"
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter320-f15
f=/tmp/w320x/acaroot.txt
echo "top=$(grep -c '^--- FAIL' $f) all=$(grep -c -- '--- FAIL' $f) leaf=$(grep -- '--- FAIL' $f | grep -vc -- '--- FAIL: TestExecToolPlacementReapplied ') mkdir=$(grep -c 'TempDir: mkdir /tmp/w320x/tmp-x134x-ro' $f) timeout=$(grep -c 'panic: test timed out' $f) run=$(grep -c '^=== RUN' $f)"
chmod 755 /tmp/w320x/tmp-x134x-ro
```

Want: rc 1, `top=6`, `all=11`, `leaf=10`, `mkdir=10` (**leaf must equal mkdir**), `timeout=0`
(design step 4: a `panic: test timed out` is an acceptance FAILURE, not an expected read-only-root
failure), `run=11` (control: 6 tests + 5 `PlacementReapplied` subtests started). The binary must be
prebuilt: `go test` itself puts its work dir under `GOTMPDIR` and fails first on a mode-555 root.
Planner's prototype: rc 1, elapsed 1 s, `top=6 all=11 leaf=10 mkdir=10 timeout=0 run=11` (V-P6).

#### M1.6 Mutation drill (AC-C, design §4 M1–M7)

Protocol (§6) for each row: snapshot, apply **one** edit with the Edit tool (old line → new line,
exactly as in the code block below), run **only** that mutant's killer line(s), see every listed
`--- FAIL:` row, restore, check the sha. "Control PASS" rows must print `--- PASS:` under the mutant.
Every mutant edits only `f15RefusalProblem`.

**Mutant edits** (each a one-line replacement; copy from here, not from the table):

```text
M1 MUT-PATH-SUBSTRING   old: 	if after != before {
                        new: 	if after != before || strings.Contains(string(out), "134") {
M2 MUT-NO-SPAWN-GUARD   old: 	if after != before {
                        new: 	if false && after != before {
M3 MUT-LOOSE-DECODE     old: 	dec.DisallowUnknownFields()
                        new: (delete the line)
M4 MUT-OK-IGNORED       old: 	if r.OK == nil || *r.OK {
                        new: 	if false && (r.OK == nil || *r.OK) {
M5 MUT-REFUSED-OPTIONAL old: 	if r.Refused == nil || *r.Refused == "" {
                        new: 	if r.Refused == nil {
M6 MUT-ACCEPT-ALL       old: 	dec := json.NewDecoder(bytes.NewReader(out))
                        new: three lines inserted ABOVE it:
                             	if true {
                             		return ""
                             	}
M7 MUT-NO-EOF-CHECK     old: 	if err := dec.Decode(&json.RawMessage{}); err != io.EOF {
                        new: 	if err := dec.Decode(&json.RawMessage{}); false && err != io.EOF {
```

(The leading whitespace is one tab. Every `old` line occurs exactly once in the file after M1.2;
the planner's script asserted `count == 1` for each, V-P7.) Every form compiles: M6 uses `if true`
so nothing is unreachable to the compiler, and M7 keeps `err`, `io` and `fmt` used (the design's
first M7 form, deleting the block, failed with `"fmt" imported and not used`, design V22). A mutant
that does not compile is a **protocol error, never a kill**. Do not run `gofmt` on a mutant.

| ID | Named killer (must FAIL) | In-run control (must PASS) | Planner's prototype (V-P7) |
|---|---|---|---|
| M0 OLD-HELPER | T2 (done in M1.1 as AC-B) | — | 10/10 FAIL |
| M1 MUT-PATH-SUBSTRING | T1 `refusal-temp-root-134`, T1 `trailing-whitespace`; T2 | — | all 3 FAIL. T2's message reads `(the call spawned)` because the mutant shares that branch (D3); the kill is the red |
| M2 MUT-NO-SPAWN-GUARD | T1 `post-launch-f15`; T3 `broken-node` (`a launch that exited 134 passes requireF15`) | T3 `command-aborted` (the strict decode still rejects it) | 2 FAIL, control PASS |
| M3 MUT-LOOSE-DECODE | T1 `refusal-plus-exit-code` | T1 `launched-exit-134` (rejected by the `ok` check) | FAIL, control PASS |
| M4 MUT-OK-IGNORED | T1 `ok-true` | — | FAIL |
| M5 MUT-REFUSED-OPTIONAL | T1 `refused-empty` | — | FAIL |
| M6 MUT-ACCEPT-ALL | T1's 7 reject rows; T3 `broken-node` **and** `command-aborted` | — | 7 + 2 FAIL |
| M7 MUT-NO-EOF-CHECK | T1 `concatenated-json`, T1 `trailing-garbage` (`problem "", want accepted=false`) | — | 2 FAIL |

**Killer commands.** Copy from this block. Each `-run` targets exactly the named rows.

```sh
T1=TestF15RefusalProblemIsStructural; T2=TestExecToolF15RefusalUnderTempRoot134; T3=TestExecToolExit134ChecksNode
mkdir -p /tmp/w320x/tmp-ctl-m
k() { GOTMPDIR=/tmp/w320x/tmp-ctl-m perl -e 'alarm 300; exec @ARGV' go test ./host/broker/ -run "$1" -count=1 -v > /tmp/w320x/mut-$2.txt 2>&1; echo "rc=$?"; grep -- '--- FAIL:\|--- PASS:\|build failed\|imported and not used\|declared and not used' /tmp/w320x/mut-$2.txt; }
# M1: want FAIL T1/refusal-temp-root-134, T1/trailing-whitespace, and T2
k "^$T1\$/^(refusal-temp-root-134|trailing-whitespace)\$" M1a; k "^$T2\$" M1b
# M2: want FAIL T1/post-launch-f15, T3/broken-node; control PASS T3/command-aborted
k "^$T1\$/^post-launch-f15\$" M2a; k "^$T3\$/^(broken-node|command-aborted)\$" M2b
# M3: want FAIL T1/refusal-plus-exit-code; control PASS T1/launched-exit-134
k "^$T1\$/^(refusal-plus-exit-code|launched-exit-134)\$" M3
# M4: want FAIL T1/ok-true
k "^$T1\$/^ok-true\$" M4
# M5: want FAIL T1/refused-empty
k "^$T1\$/^refused-empty\$" M5
# M6: want FAIL the 7 T1 reject rows, and T3/broken-node + T3/command-aborted
k "^$T1\$/^(post-launch-f15|launched-exit-134|refusal-plus-exit-code|ok-true|refused-empty|concatenated-json|trailing-garbage)\$" M6a; k "^$T3\$" M6b
# M7: want FAIL T1/concatenated-json, T1/trailing-garbage
k "^$T1\$/^(concatenated-json|trailing-garbage)\$" M7
```

Each `k` call must show its listed `--- FAIL:` subtest lines (and the parent `--- FAIL:` line), and
rc ≠ 0. A `build failed` / `imported and not used` line means a protocol error: restore, fix the
mutant's form, re-run. After the last mutant, `shasum -a 256 -c /tmp/w320x/fixed.sha` must print `OK`.

---

### Milestone M2 — AC-D (the whole package under a `134` root, plain and race; full gate), executor report

**Files:** `design_docs/verification/world-iter320/executor-report.md` (new) and
`design_docs/verification/world-iter320/executor-*.txt` (new evidence copies). No code change in M2.
Do not run M2 while a mutant is applied (`shasum -a 256 -c /tmp/w320x/fixed.sha` must say `OK` first).

#### M2.1 AC-D (design §5, copied; `/tmp/wNNNd` = `/tmp/w320x`)

```sh
rm -rf /tmp/w320x/tmp-x134x; mkdir -p /tmp/w320x/tmp-x134x
s=$(date +%s); GOTMPDIR=/tmp/w320x/tmp-x134x perl -e 'alarm 1200; exec @ARGV' go test ./host/broker/ -count=1 > /tmp/w320x/acd-plain.txt 2>&1; echo "rc=$? elapsed=$(( $(date +%s) - s ))s"; tail -3 /tmp/w320x/acd-plain.txt
s=$(date +%s); GOTMPDIR=/tmp/w320x/tmp-x134x perl -e 'alarm 1500; exec @ARGV' go test -race ./host/broker/ -count=1 > /tmp/w320x/acd-race.txt 2>&1; echo "rc=$? elapsed=$(( $(date +%s) - s ))s"; tail -3 /tmp/w320x/acd-race.txt; grep -c 'DATA RACE' /tmp/w320x/acd-race.txt
perl -e 'alarm 2400; exec @ARGV' sh -c 'go build ./... && go test ./...' > /tmp/w320x/acd-full.txt 2>&1; echo "rc=$?"; grep -E '^(FAIL|ok|---)' /tmp/w320x/acd-full.txt | grep -v '^ok' | head -20
perl -e 'alarm 900; exec @ARGV' ./scripts/verify_ail.sh > /tmp/w320x/acd-ail.txt 2>&1; echo "rc=$?"; tail -3 /tmp/w320x/acd-ail.txt
```

Want: plain rc 0 and `ok  	github.com/sunholo-data/ailang-world/host/broker`; race rc 0, the same
`ok` line and `DATA RACE` count 0; full gate rc 0 (no `FAIL` line); `verify_ail.sh` rc 0 (no `.ail`
touched). Run these one at a time, never alongside AC-A or a mutant run. Planner's prototype: plain `ok … 221.614s` (rc 0, 222 s), race `ok … 322.370s` (rc 0, 328 s, 0 `DATA RACE`);
the full gate takes about 520 s (V-P8, V-P8b).

**Known unrelated load flakes** (seen by earlier sprints or by the planner's full-gate prototype run,
V-P8b; none is in `exec_tools*`): `TestWorktreeSnapshotRefusesAConfiguredFilter`
(`host/broker/worktree_delta_test.go`), `TestExecProtectedCheckTimeout/after-cancel`,
`TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse/chatty` (`host/broker`),
`TestMarketingSupervisorFakeDaemonCleanup/shell-hup` (`scripts/marketingcheck`), and
`host/verifygate TestMissionConfigDefaultFlowIntegration` timing out. If the ONLY failure is one of these, re-run that test alone with `-count=3` (same
`GOTMPDIR` leg if it was in an AC-D `host/broker` leg). PASS alone = record it as a load flake and move
on. **Any failure in an `exec_tools*` test, or any other failure, is a stop:** report it, do not
patch around it.

**CI (controller reads it, not the executor):** the fixed head's PR CI must show
`ok  	github.com/sunholo-data/ailang-world/host/broker` in both "go race gate" and "go host build + test gate".
No CI list edit is needed: no `tests=` list names a `TestExecTool*` test (design V18; re-checked V-P9).

#### M2.2 Executor report and evidence

Copy the logs, then write the report:

```sh
for f in acb-m0 aca-tmp-x134x aca-tmp-ctl acaroot mut-M1a mut-M1b mut-M2a mut-M2b mut-M3 mut-M4 mut-M5 mut-M6a mut-M6b mut-M7 acd-plain acd-race acd-ail; do
  cp /tmp/w320x/$f.txt design_docs/verification/world-iter320/executor-$f.txt
done
grep -E '^(FAIL|ok)' /tmp/w320x/acd-full.txt > design_docs/verification/world-iter320/executor-acd-full-summary.txt
```

`design_docs/verification/world-iter320/executor-report.md` follows the shape of
`design_docs/verification/world-iter318/executor-report.md`: worktree/branch/base, the base sha256
check, then per milestone a table of every command, rc and counted result (AC-B, G1–G3, AC-A for
both roots, AC-A-ROOT, every mutant with its FAIL/PASS rows and its restore `OK`, AC-D's four legs),
and a **Deviations** section (anything that differed from this plan, with the measurement). State
plainly that `design_docs/world-mission*.md` were not touched.

#### M2.3 Acceptance (M2)

```sh
shasum -a 256 -c /tmp/w320x/fixed.sha                 # OK (no mutant left in)
shasum -a 256 host/broker/exec_tools_round2_test.go host/broker/handlers_exec.go   # equal to §3
git diff --name-only                                   # exactly: host/broker/exec_tools_test.go
git status --porcelain                                 # the M1 file, the two plan files, and the new executor-* files under design_docs/verification/world-iter320/
git status --porcelain | grep -c 'world-mission'       # want 0
```

---

## 5. Controller-owned items

- **Row 234 charter/queue text**: the controller writes it into `design_docs/world-mission.md`
  (queue + STATUS) and the mission log, from design §1/§2 and the executor report. The executor never
  edits `design_docs/world-mission*.md`.
- Moving the design doc and this plan from `planned/` to `implemented/`.
- Committing, the PR, and reading CI: both "go race gate" and "go host build + test gate" must show
  `ok  	github.com/sunholo-data/ailang-world/host/broker` on the fixed head (design AC-D).

## 6. Mutation protocol (every mutant)

Take the snapshot **after** M1.4 (gates and AC-A) is green:

```sh
cp host/broker/exec_tools_test.go /tmp/w320x/exec_tools_test.go.snap
shasum -a 256 host/broker/exec_tools_test.go > /tmp/w320x/fixed.sha
cat /tmp/w320x/fixed.sha        # planner's prototype: 475fb831692647b87b00690c1142ae2aa6efc9aa8e41df236c2ac2682e039be6
#   apply ONE edit from the M1.6 block with the Edit tool
#   run ONLY that mutant's killer line(s) from the M1.6 killer block
cp /tmp/w320x/exec_tools_test.go.snap host/broker/exec_tools_test.go
shasum -a 256 -c /tmp/w320x/fixed.sha        # OK
```

If your fixed file's sha256 equals the planner's `475fb831…`, your M1.2/M1.3 edits are byte-identical
to the prototype. If it differs, the edits are still acceptable when G1–G4 pass, but say so in the
report and `diff` it against the plan's blocks.

## 7. Risks

| Risk | Mitigation |
|---|---|
| The predicate is stricter than today's guard (any extra key or trailing content reds) | Intended (design §3). Production emits exactly `{"ok","refused"}` (`handlers_exec.go:96-97`). The whole package is green under a `134` root (V-P8) |
| T2 sets `GOTMPDIR` with `t.Setenv` | It is restored after the test; `t.Setenv` forbids `t.Parallel`, so do not add it |
| An ambient `GOTMPDIR` makes an AC leg vacuous | §1 `unset GOTMPDIR`; every AC leg assigns it inline; AC-A-ROOT proves the assignment is honoured |
| A mode-555 dir left behind breaks `rm -rf` later | M1.5 ends with `chmod 755` |
| `go test ./...` reds on a known load flake | M2.1's re-run-alone rule. Never patch another test |

## 8. Definition of done

1. M1: AC-B red 10/10 before the fix. After it, G1–G3 green; AC-A both roots 180/320/0/9; AC-A-ROOT
   rc 1, 6/11/10/10, 0 timeout; every M1–M7 named killer red and every control green; the file
   restores `OK`.
2. M2: AC-D's four legs green; executor report and evidence files written; `git diff --name-only` is
   exactly `host/broker/exec_tools_test.go`; no `world-mission*` file touched.

## 9. Deviations from the design (planner-measured; flagged to the controller)

- **D1 (T1 literal content).** The design names T1's rows and their verdicts but gives no literal JSON.
  The plan fixes the literals (§M1.3). `launched-exit-134` uses spawns 0→1 (a real launch moves the
  counter); it is rejected by the strict decode / `ok` check independently of the counter, so only M6
  kills it, as design §4 says. `refusal-temp-root-134` holds both `x134x` and `…1349870625`, as the
  design specifies.
- **D2 (M0 realised as the red-first step).** The design's M0 is "base `requireF15` (`:75-98`) with
  T2 added". The plan does exactly that, in TDD order: T2 is appended to the **pristine** file before
  the fix (M1.1), so no base code has to be spliced back in. Measured: 10/10 FAIL, 0 PASS (V-P4),
  matching design V13.
- **D3 (M1's failure text).** The plan's M1 form ORs the `"134"` reject into the spawn branch, so a
  killed T2 prints `not a pre-launch refusal (the call spawned)`. The kill (named rows red) is what
  counts; the label is the shared branch's. Measured kills: T1 `refusal-temp-root-134`,
  `trailing-whitespace`, and T2 (V-P7).
- **D4 (gofmt form).** Design §2's snippet writes `if … { return … }` on one line, which gofmt
  expands. The plan's block is the gofmt form (G1 empty), same semantics, and drops the inline
  `// [r2]` marker. The design's doc comment is kept verbatim.
- **D5 (`requireF15` body).** The design shows only the predicate. The plan reads the counter once
  (`after := f.h.Spawns()`), calls the predicate **before** the loose `json.Unmarshal` (so a non-JSON
  reply fails with the `not a pre-launch refusal (not the refusal envelope: …)` message, not a bare
  `t.Fatal(err)`), and drops the now-unused `OK` field from the loose struct, which only feeds the
  `parts` loop.
- **D6 (size).** Measured diff `+91/−4` (design estimate `+75/−8`): the predicate is gofmt-expanded
  (+10 lines) and T1/T2 carry doc comments.
- **D7 (T3's own name holds `134`).** `TestExecToolExit134ChecksNode`'s temp paths contain `134`
  through its test name (prototype: `…/TestExecToolExit134ChecksNodebroken-node2674606336/001/node`,
  V-P7). So in this test the base substring would have "caught" shape 2 by accident; design V8's
  `has134=false` was measured in a probe test whose name had no `134`. This does not affect the fix:
  T3 asserts on `f15RefusalProblem`, which reads no text, and M2/M6 show T3 sees the spawn and decode
  checks.
- **D8 (placeholder and scratch dir).** The design's `/tmp/wNNNd` is instantiated as `/tmp/w320x`
  for the executor (the planner used `/tmp/w320p`). AC-A's `<six>` is expanded in the fenced block.
- **D9 (subtest count added to AC-A).** AC-A also asserts `320` `^    --- PASS` lines
  ((5 PlacementReapplied + 2 Exit134 + 9 T1) × 20), measured on both roots, so a silently skipped
  T1 row cannot pass AC-A.

**Nothing in the design was found false.** These held when re-measured: the guard at `:89`;
`requireF15` at `:75-98`; the six callers and 9 call lines (design V3); `handlers_exec.go:96-97`
(`execRefusal`), `:209` (`ExitCode *int` `json:"exit_code"`), `:322` `h.spawns.Add(1)`, `:323`
`runBounded`, `:326` the `code == 134` branch, `:333` the post-launch `execRefusal`; `Spawns() int64`
at `:435`; T2 red at base 10/10; AC-A 180 PASS both roots; AC-A-ROOT 6 / 10 = 10 / 0 timeout;
every M1–M7 kill exactly as design V22.

## 10. Verification notes (planner, this session, worktree `5407c160`)

| # | Claim | Command | Observed |
|---|---|---|---|
| V-P1 | Base: the Go tree equals `785b3bd1` | `git -C <wt> log --oneline -3; git status --porcelain; git diff --stat 785b3bd1 HEAD` | HEAD `5407c160 design(row 234): …`; status empty before the plan files; 4 files, all under `design_docs/` (the design, `quorum/r1.json`, `quorum/r2.json`, `v23-r3-rerun.log`) |
| V-P2 | Toolchain and env | `go version`; `~/.pinned-ailang/ailang --version`; `echo GOTMPDIR=$GOTMPDIR TMPDIR=$TMPDIR` | `go1.26.6 darwin/arm64`; `AILANG v0.41.0`; `GOTMPDIR=` (unset), `TMPDIR=/var/folders/kr/…/T/` |
| V-P3 | Edit sites | `sed -n '3,13p;75,98p;230,257p' host/broker/exec_tools_test.go`; `grep -n 'requireF15(' host/broker/*_test.go`; `sed -n 92,99p;205,212p;318,340p host/broker/handlers_exec.go`; `grep -n 'func (h \*ExecHandler) Spawns'` | imports `:3-13`; old struct + guard `:82-91`; guard `:89`; T3 `f := newExecFixture` `:243`, `f.call` `:244`; the 9 call lines of §3.1; `execRefusal` `:96-97`; `ExitCode` `:209`; spawns `:322`; `Spawns() int64` `:435`; base sha256s in §3 |
| V-P4 | M0 / AC-B (prototype `/tmp/w320p/full`, T2 appended to the pristine file) | `go test ./host/broker/ -run '^TestExecToolF15RefusalUnderTempRoot134$' -count=10 -v` | rc 1, 10 FAIL, 0 PASS; `exec_tools_test.go:473: not a pre-launch refusal: {"ok":false,"refused":"exec node /private/var/folders/…/T/w234-x134x-2539520704/TestExecToolF15RefusalUnderTempRoot1342643642204/001/node …` |
| V-P5 | Fix applied by script (`/tmp/w320p/scripts/apply.py`, exact replaces, `assert count==1`) + AC-A | `gofmt -l`; `go vet`; `go test -run '^$'`; `diff -u … \| diffstat`; AC-A loop of §M1.4 with `/tmp/w320p/<d>` | gofmt empty; vet ok; fence ok; `+91/−4`; fixed sha256 `475fb831…039be6`. `tmp-x134x`: rc 0, top-PASS 180, sub-PASS 320, FAIL 0, SKIP 0, names 9, `ok … 43.092s`. `tmp-ctl`: the same, `ok … 38.961s` |
| V-P6 | AC-A-ROOT | §M1.5 with `/tmp/w320p/broker.test` and `/tmp/w320p/tmp-x134x-ro` (mode 555) | build rc 0; rc 1, elapsed 1 s, `top=6 all=11 leaf=10 mkdir=10 timeout=0 run=11` |
| V-P7 | AC-C mutants (copy `/tmp/w320p/mut`; `/tmp/w320p/scripts/mut.py` one exact replace `assert count==1`; `GOTMPDIR=/tmp/w320p/tmp-ctl-m`; the §M1.6 killer lines, plus a wide `^(T1\|T3)$` run per mutant as a cross-check) | `/tmp/w320p/scripts/muts.sh` | every killer rc 1. M1: T1 `refusal-temp-root-134`, `trailing-whitespace` FAIL; T2 FAIL (`(the call spawned)`). M2: T1 `post-launch-f15` FAIL; T3 `broken-node` FAIL (`a launch that exited 134 passes requireF15`), `command-aborted` PASS. M3: `refusal-plus-exit-code` FAIL, `launched-exit-134` PASS. M4: `ok-true` FAIL. M5: `refused-empty` FAIL. M6: all 7 reject rows FAIL; T3 both FAIL. M7: `concatenated-json`, `trailing-garbage` FAIL (`problem "", want accepted=false`). Wide runs: each mutant reds **only** its named T1/T3 rows (M1–M5, M7 leave T3 green; M6 leaves the 2 accept rows green). No build failure. `shasum -c` OK after every restore |
| V-P8 | AC-D on the prototype (`/tmp/w320p/full`, full `git archive HEAD` + fixed file, `AILANG_BIN` pinned) | `GOTMPDIR=/tmp/w320p/tmp-x134x go test ./host/broker/ -count=1`; the same with `-race` | rc 0 and rc 0: plain `ok  github.com/sunholo-data/ailang-world/host/broker 221.614s` (222 s wall); race `ok … 322.370s` (328 s wall), 0 `DATA RACE` |
| V-P8b | Full gate on the prototype (informational; the executor runs it in the real worktree) | `go build ./... && go test ./...` in `/tmp/w320p/full`, ambient `GOTMPDIR` unset | rc 1 in 519 s: 32 `ok` packages and 3 FAIL, none in `exec_tools*`. (1) `host/broker` `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse/chatty` (arm 6 `accepts 2`, under full-suite load); re-run alone `-count=3`: 3 PASS. (2) `scripts/marketingcheck` `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` (`cleanup exceeded bound: 29.5s`); alone `-count=3`: 3 PASS. (3) `host/pins` `TestReleaseLiteralsExplainedByManifest` (`git ls-files: exit status 128`): the `git archive` copy is not a git repo; in the worktree it PASSes (1/1). So the full gate needs the real worktree, which is where the executor runs it |
| V-P9 | Charter row 234 absent; no CI env/list hit | `grep -cE '^\| *234 *\|' design_docs/world-mission.md`; `grep -nE 'GOTMPDIR\|TMPDIR' .github/workflows/*.yml`; `grep -nE 'TestExecTool\|requireF15' .github/workflows/ci.yml docs/QUICKSTART.md`; `grep -n 'go-race-verify\|go host build' ci.yml` | `0`; no output; no output; `:117 name: go host build + test gate`, `:418 go-race-verify:` |
| V-P10 | The repo worktree was not modified by the prototype | `shasum -a 256 -c /tmp/w320p/base.sha` in the worktree; `git status --porcelain` | all `OK`; only the two plan files |
| V-P11 | The plan's fenced blocks reproduce the prototype byte-for-byte | extract the 9 fenced `go` blocks from this file, apply them to a copy of the base `exec_tools_test.go` in the M1.1–M1.3 order (`/tmp/w320p/rt`); `grep` each M1.6 `old:` line in the result | sha256 `475fb831…039be6`, equal to the prototype; each of the 7 `old:` lines occurs exactly once |
