# Sprint plan — `VL.B` (queue item 9 `w-verify-binary-lockfile`)

**Milestone**: `VL.B` — *the shim accept-arms assert the FULL gate outcome, because job 2 now has a solver*
**Status**: PLANNED · one small milestone · no split needed
**Item**: 9 `w-verify-binary-lockfile` (charter-tracked; **no design-doc file** — the spec is
`design_docs/world-mission.md` STATUS iter-68 + OD rows `9/OD-10`, `9/OD-11`)
**Base**: `dev` @ `541258e`, clean
**Unblocked by**: `9/OD-11` **RATIFIED** (Mark, verbatim: *"Yes you can install z3 on cicd"*).
This widens `9/OD-10` clause (a) (item 9 scoped to zero `.github/` edits) **for the Z3 install only**.
**Executor**: `codex:gpt-5.6-sol`, sandboxed worktree, **NO git write permission**.
**THE CONTROLLER MAKES ALL COMMITS.** The executor never runs `git commit`, `git add`, `git push`,
`git checkout`, `git stash`, `git restore`, or `gh pr`. Restores are `cp` from a backup — see §7.

---

## 0. Planner's verification of the handoff — what I reproduced, corrected, and refuted

Every controller measurement was re-run first-party on this rig. Verdicts:

| # | Controller premise | Verdict | My command / result |
|---|---|---|---|
| 1 | Z3 installed in job 1 only | **CONFIRMED**, count corrected | `grep -rin z3 .github/` = **11** hits (controller said 7 — that was case-sensitive `z3`); `sed -n '92,174p' .github/workflows/ci.yml \| grep -ci z3` = **0**. Job 2 has no solver. |
| 2 | binary SHELLS OUT to z3, honours `AILANG_Z3_PATH` | **CONFIRMED**, path list corrected | `otool -L /tmp/ailang-v0300/ailang \| grep -ci z3` = **0** (no libz3). `strings` shows `AILANG_Z3_PATH`, `/opt/homebrew/bin/z3`, `/usr/local/bin/z3`, `/snap/bin/z3` **and `/usr/bin/z3`** (controller omitted `/usr/bin/z3`). Also found the string **`AILANG_Z3_PATH set to %q but file not found`** — so a bogus `AILANG_Z3_PATH` is a *named refusal*, **not** a silent fallback to PATH. That is what makes the control faithful. |
| 3 | PATH-hiding is not a valid control; `AILANG_Z3_PATH=<nonexistent>` is | **CONFIRMED** | With no `AILANG_Z3_PATH` at all the gate still passes (hardcoded `/opt/homebrew/bin/z3` search) → `✓ verify gate PASSED`. With `AILANG_Z3_PATH=/tmp/w_nosuch/z3` the gate reds. No fallback. |
| 4 | the two-arm defect | **CONFIRMED** (see §1) | verbatim, both arms |
| 5 | "a full `verify_ail.sh` run takes **4 seconds**" | **REFINED** | Measured `2.944 s` wall for the passing arm and `0.680 s` for the Z3-absent (fail-fast) arm. The controller's 4 s is the right order of magnitude; the number that actually matters is the **package** delta, §5. |
| 6 | `WORLD_PKG_AILANG_BIN` is the real pin in `shimEnv` | **CONFIRMED** | `ail_binary_gate_test.go:86` — `"WORLD_PKG_AILANG_BIN": pinned`. |
| 7 | `requirePinned` fails loudly, never skips | **CONFIRMED** | `:38-49`, two `t.Fatal` paths. Preserved untouched by this plan. |

### Three things the controller did NOT say, all measured, all change the plan

**(i) `AILANG_BIN` must be ABSOLUTE, or you get a FALSE "shim did not delegate".**
My first reproduction of Arm A used `AILANG_BIN=scripts/testdata/ailang_version_shim.sh` (relative)
and got `rc=1` with `could not parse ai-check JSON = 1` — the exact `noDelegateMarker` signature,
i.e. it *looks* like a delegation failure and is really a cwd failure (leg 1 runs inside
`( cd "$base" … )`, `verify_ail.sh:174`). The Go tests are already correct (`shimEnv` uses
`filepath.Join(repoRoot, …)`). **Every shell reproduction in this plan therefore uses `$PWD/…`.**
Executor: if you ever see `could not parse ai-check JSON` while the shim is committed and 0755,
check for a relative `AILANG_BIN` **before** concluding anything about delegation.

**(ii) `TestFixtureDiscrimination`'s leg-3 assertion is CONDITIONALLY VACUOUS in job 2 today, and
the Z3 install de-vacuifies it for free.** `:458` reads
`if strings.Contains(out, "── Leg 3") && !strings.Contains(out, "wrong compiler version")`.
Measured, same fixture, both arms:

| arm | rc | `── Leg 3` | `wrong compiler version` |
|---|---|---|---|
| Z3 present | 1 | **1** | **1** |
| Z3 absent (= CI job 2 today) | 1 | **0** | **0** |

So in CI the guarded body has **never executed**. Once job 2 has Z3, `── Leg 3` is always reached
and the condition can be made unconditional. This is a second, unclaimed win of `9/OD-11` and it is
in-scope: same file, same theme (an assertion weakened to survive a Z3-less lane). **AC5 + `MUT-VLB-FIXDISC`.**

**(iii) `verify_go.sh` runs the suite TWICE** (`:108` plain, `:112-116` `-race -timeout 8m`), so the
package cost added to job 2 is **doubled**. See §5 — it is still trivially affordable.

### One controller sketch item I am REFUTING: (d), the "two pins must agree" guard

The controller proposed duplicating `Z3_VER`/`Z3_SHA` into job 2 and then adding a test asserting
the two literals agree. **Do not detect a drift you can decline to create.** GitHub Actions supports
a workflow-level `env:` block that every step of every job inherits, so the pin can be declared
**once**. Measured at base: `grep -c 'Z3_VER=' ci.yml` = **1**, `grep -c 'Z3_VER:' ci.yml` = **0** —
there is no workflow-level block today, and adding one deletes the drift class instead of guarding it.
Consequently:

- there is **no** `MUT-…-SHA-DRIFT` mutation in this plan, because after the hoist a wrong sha is
  wrong in *both* jobs and reds `sha256sum -c` loudly in both. There is nothing left to disagree.
- what survives is a *different* class — someone re-introduces a job-local literal, or job 2's
  install line gets edited away. That class is real and is guarded by **AC6** with two mutations
  (`MUT-VLB-PIN-DUP`, `MUT-VLB-PIN-DROP`).

**SCOPE NOTE the controller may veto before routing.** The hoist edits job 1's *existing, green*
step (deleting its two literal assignments). `9/OD-11` ratifies "install z3 on cicd"; it does not
literally speak to refactoring job 1. Blast radius if the hoist is wrong: `set -u` makes it
`Z3_VER: unbound variable` → **job 1 red on the PR, loud, immediate, not silent**. I judge that
acceptable and the single-pin payoff worth it. **Fallback if vetoed**: duplicate both literals into
job 2 verbatim and change AC6's counts from `Z3_VER: == 1 / Z3_VER= == 0` to
`Z3_VER= == 2 / the two values are byte-equal`; everything else in this plan is unchanged.

---

## 1. The defect, measured in both arms at base (this is the spine)

Fixture = the committed shim driving the **real** `scripts/verify_ail.sh`, i.e. exactly what
`host/verifygate`'s accept-arms do. `R=/Users/voightkampff/dev/sunholo-data/ailang-world`.

```
env AILANG_BIN=$R/scripts/testdata/ailang_version_shim.sh \
    AILANG_SHIM_VERSION_LINE="AILANG v0.33.0" \
    AILANG_SHIM_DELEGATE=/tmp/ailang-v0300/ailang \
    WORLD_PKG_AILANG_BIN=/tmp/ailang-v0300/ailang \
    AILANG_Z3_PATH=<A or B> ./scripts/verify_ail.sh
```

| marker | Arm A `AILANG_Z3_PATH=/opt/homebrew/bin/z3` | Arm B `AILANG_Z3_PATH=/tmp/w_nosuch/z3` (**= CI job 2 today**) |
|---|---|---|
| `rc` | **0** | **1** |
| `AILANG_BIN refused` | 0 | 0 |
| `── Leg 1` | 1 | 1 |
| `could not parse ai-check JSON` | 0 | 0 |
| **`verify gate PASSED`** | **1** | **0** |

Arm B's failure line: `✗ world/contracts.ail: required identity (contracts, isValidNextWorld)
MISSING from verify.results[] (vanished silently, V20)`.

**The three markers `requireProceeded` asserts (`:121-132`) are IDENTICAL in both arms.** Only
`verify gate PASSED` separates them. And the package-level consequence, measured at base:

```
$ AILANG_BIN=/tmp/ailang-v0300/ailang AILANG_Z3_PATH=/tmp/w_nosuch/z3 go test ./host/verifygate/ -count=1
ok   github.com/sunholo-data/ailang-world/host/verifygate   5.197s        ← rc=0
$ AILANG_BIN=/tmp/ailang-v0300/ailang AILANG_Z3_PATH=/opt/homebrew/bin/z3 go test ./host/verifygate/ -count=1
ok   github.com/sunholo-data/ailang-world/host/verifygate  14.482s        ← rc=0
```

**The accept-arms are GREEN in the lane where the gate they drive is RED.** That is iter-68's spine
turned inward — a green proving the tree passes where you ran it, not where it must — and `VL.B`
closes it. The "before" arm needs no source archaeology: **it is a live, reproducible `rc=0`**.

Solver discrimination, from `ai-check` directly (this is what the new probe reads):

```
$ AILANG_Z3_PATH=/opt/homebrew/bin/z3 /tmp/ailang-v0300/ailang ai-check -timeout 5s world/contracts.ail
  → verify: {"available": true,  "verified": 1, ... "results":[{"function":"isValidNextWorld","status":"verified"}]}
$ AILANG_Z3_PATH=/tmp/w_nosuch/z3  /tmp/ailang-v0300/ailang ai-check -timeout 5s world/contracts.ail
  → verify: {"available": false, "verified": 0, ... "results":[]}
```

---

## 2. Scope — exactly two files

| file | change |
|---|---|
| `.github/workflows/ci.yml` | hoist the Z3 pin to workflow-level `env:`; add the identical install step to job 2 `go-verify` |
| `host/verifygate/ail_binary_gate_test.go` | assert `verify gate PASSED`; rewrite the now-false comment; strip ambient `AILANG_Z3_PATH` in `runGate`; add the self-controlled solver probe; make the leg-3 assertion unconditional; add the ci.yml pin-structure guard |

**No other file may change.** AC0 enforces this.

Base sha256 (record before editing; used for the restore protocol):

```
27450846847ffb4a74b8ad34b23b5030bac5ed9b5e02d9efa41d45e82343bb9b  host/verifygate/ail_binary_gate_test.go
7677ccb60a2771d0e3acf09dd3b759441fd643a63d44d0b3d0d1c6cbf3edf26c  .github/workflows/ci.yml
```

---

## 3. Prescriptive edits

Set `R=<repo root>` and `cd "$R"` for everything below.

### E1 — `.github/workflows/ci.yml`: hoist the pin

Line 7 is the blank line between `  pull_request:` (line 6) and `jobs:` (line 8).
**Insert before `jobs:`**, so the file reads:

```yaml
  pull_request:

# The z3 pin is declared ONCE, at workflow scope, because BOTH jobs must install the SAME solver.
# Job 1 runs `ai-check` directly; job 2 runs it THROUGH host/verifygate's shim arms, which drive
# scripts/verify_ail.sh end to end. Two job-local literals would let the lanes drift silently —
# the same "two lanes resolve different things" defect 9/CF-A-2 was. Ratified: 9/OD-11.
env:
  Z3_VER: 4.16.0
  Z3_SHA: 7288c49a5bd6dbafd7b0b0d1f65956b91672da24b08f09242919af159be3418e

jobs:
```

### E2 — `.github/workflows/ci.yml`: job 1 reads the hoisted pin

In the step at line 34, **delete exactly these two lines** (currently 41 and 43):

```
          Z3_VER=4.16.0
          Z3_SHA=7288c49a5bd6dbafd7b0b0d1f65956b91672da24b08f09242919af159be3418e
```

Nothing else in that step changes: `Z3_DIR="z3-${Z3_VER}-x64-glibc-2.39"` and
`echo "${Z3_SHA}  z3.zip" | sha256sum -c -` now resolve from the workflow env.

### E3 — `.github/workflows/ci.yml`: job 2 gets the identical step

Insert **after** job 2's `echo "AILANG_BIN=$HOME/.local/bin/ailang" >> "$GITHUB_ENV"` (line 138)
and **before** `- name: go build + test gate` (line 140):

```yaml
      - name: Install Z3 ${{ env.Z3_VER }} (host/verifygate drives the REAL .ail gate under go test)
        run: |
          set -euo pipefail
          # 9/OD-11, RATIFIED: "Yes you can install z3 on cicd".
          # host/verifygate's shim arms run scripts/verify_ail.sh end to end. Without a solver
          # `ai-check` reports verify.available=false and every contract vanishes from
          # verify.results[] (V20/V27), so the gate reds at leg 1 and the accept-arms could only
          # assert markers emitted BEFORE Z3 matters — they were green in a lane where the gate
          # they drive was red. Same version and sha as job 1, from the workflow-level pin above.
          Z3_DIR="z3-${Z3_VER}-x64-glibc-2.39"
          curl -fsSL -o z3.zip "https://github.com/Z3Prover/z3/releases/download/z3-${Z3_VER}/${Z3_DIR}.zip"
          echo "${Z3_SHA}  z3.zip" | sha256sum -c -
          unzip -q z3.zip
          sudo install -m 0755 "${Z3_DIR}/bin/z3" /usr/local/bin/z3
          z3 --version
```

### E4 — `ail_binary_gate_test.go`: strip ambient `AILANG_Z3_PATH` in `runGate`

At `:56-58` the `blocked` map has three lines. **Add `AILANG_Z3_PATH` to it**:

```go
		blocked := map[string]bool{
			"AILANG_BIN": true, "WORLD_PKG_AILANG_BIN": true,
			"AILANG_SHIM_VERSION_LINE": true, "AILANG_SHIM_DELEGATE": true,
			// Ambient AILANG_Z3_PATH is stripped so every arm resolves the solver the way the
			// LANE does (the binary's own search: PATH + /usr/local/bin, /usr/bin, /snap/bin,
			// /opt/homebrew/bin). `blocked` filters os.Environ() only — an explicit entry in the
			// `env` map below still reaches the child, which is how a Z3-absent control is armed
			// deterministically instead of relying on os/exec's dedup order.
			"AILANG_Z3_PATH": true,
		}
```

### E5 — `ail_binary_gate_test.go`: add the marker and assert it

Add to the `const` block at `:101-109`:

```go
	// The gate's terminal success line (verify_ail.sh:304). Unambiguous: the only other "gate
	// PASSED" line is `world package gate PASSED`, which does not contain this substring.
	passedMarker = "verify gate PASSED"
```

**Replace the entire comment at `:111-120` and add the fourth assertion.** The old comment's
central claim — *"It deliberately does NOT assert `verify gate PASSED`"* — is now **false**, and a
comment that contradicts its code is a defect this mission has already paid for. New text:

```go
// requireProceeded asserts the version block ACCEPTED, the real gate ran on the real delegate, and
// THE GATE PASSED.
//
// The fourth assertion is what 9/OD-11 bought. It was absent because these tests run inside
// `go test ./...`, i.e. CI's go-verify job, which installed no Z3 — and without a solver `ai-check`
// reports verify.available=false and every contract vanishes from verify.results[] (V20/V27), so
// the gate reds at leg 1 for a reason with nothing to do with the version predicate. The three
// remaining markers were re-aimed at that lane and are all emitted before Z3 matters. Measured, the
// cost of that portability was exact and total: with AILANG_Z3_PATH pointed at a missing file, a
// PASSING gate (rc=0) and a FAILING gate (rc=1) produce IDENTICAL values for all three — refused=0,
// leg1=1, noDelegate=0 — so the accept-arms were green in the lane where the gate they drive was
// red. `verify gate PASSED` is the only marker that separates them.
//
// Mark ratified the install ("Yes you can install z3 on cicd", 9/OD-11); ci.yml now installs the
// same pinned Z3 in BOTH jobs, so the assertion is reachable everywhere this package runs. If a
// lane loses its solver, TestSolverAvailableInThisLane says so by name before these arms confuse
// anyone.
func requireProceeded(t *testing.T, label, out string) {
	t.Helper()
	if strings.Contains(out, refusalMarker) {
		t.Fatalf("%s: version block refused a release token\n%s", label, out)
	}
	if !strings.Contains(out, leg1Marker) {
		t.Fatalf("%s: gate never reached leg 1 — it did not proceed past the version block\n%s", label, out)
	}
	if strings.Contains(out, noDelegateMarker) {
		t.Fatalf("%s: shim did not delegate — ai-check produced no parseable output\n%s", label, out)
	}
	if !strings.Contains(out, passedMarker) {
		t.Fatalf("%s: gate proceeded but did NOT pass — %q absent. If this lane has no Z3 the whole "+
			"gate reds at leg 1; see TestSolverAvailableInThisLane and the Z3 install steps in "+
			".github/workflows/ci.yml\n%s", label, passedMarker, out)
	}
}
```

Also update the comment at `:99-100` (`Markers the accept-arms assert on. All three are emitted
BEFORE anything Z3-dependent…`) — it now says three where there are four, and the portability
rationale no longer holds. Replace its first sentence with:

```go
// Markers the accept-arms assert on. The first three are emitted BEFORE anything Z3-dependent; the
// fourth (passedMarker) is the gate's terminal success line and needs a solver in the lane.
```

### E6 — `ail_binary_gate_test.go`: the self-controlled solver probe

Add this test (and add `encoding/json` to the import block). It is **two-armed**, so it proves in
every CI run that the instrument discriminates — no source mutation required for that claim.

```go
// TestSolverAvailableInThisLane is the self-naming instrument check for 9/OD-11. Every accept-arm
// above now asserts `verify gate PASSED`, which is unreachable without a solver; if a lane loses
// Z3 this test says exactly that, by name, instead of four arms failing with a leg-1 contract error.
//
// It is SELF-CONTROLLED: arm B points AILANG_Z3_PATH at a path that does not exist and requires
// verify.available to be FALSE. That is not decoration — it is the anti-vacuity proof, and unlike a
// source mutation it runs on every CI run forever. Measured on the rig: the binary does not link
// libz3 (`otool -L` names none) and does not fall back to PATH when AILANG_Z3_PATH is set and
// missing (it carries the string "AILANG_Z3_PATH set to %q but file not found"), so arm B is a
// faithful model of a solverless runner. PATH manipulation is NOT a valid control here: removing
// the toolchain directory also removes `go`, and the gate then reds for an unrelated reason.
func TestSolverAvailableInThisLane(t *testing.T) {
	requirePinned(t)

	type aiCheck struct {
		Check struct {
			Passed bool `json:"passed"`
		} `json:"check"`
		Verify struct {
			Available bool `json:"available"`
			Verified  int  `json:"verified"`
		} `json:"verify"`
	}

	probe := func(t *testing.T, z3 string) aiCheck {
		t.Helper()
		cmd := exec.Command(pinned, "ai-check", "-timeout", "5s", "world/contracts.ail")
		cmd.Dir = repoRoot
		env := make([]string, 0, len(os.Environ())+1)
		for _, item := range os.Environ() {
			if strings.SplitN(item, "=", 2)[0] != "AILANG_Z3_PATH" {
				env = append(env, item)
			}
		}
		if z3 != "" {
			env = append(env, "AILANG_Z3_PATH="+z3)
		}
		cmd.Env = env
		// Exit status is ADVISORY here (V10/V20) exactly as in verify_ail.sh; the JSON is
		// authoritative. Only an unparseable body is instrument failure.
		raw, _ := cmd.Output()
		start := bytes.IndexByte(raw, '{')
		if start < 0 {
			t.Fatalf("instrument failure: ai-check produced no JSON object (%d bytes): %q", len(raw), raw)
		}
		var parsed aiCheck
		if err := json.Unmarshal(raw[start:], &parsed); err != nil {
			t.Fatalf("instrument failure: ai-check JSON did not parse: %v\n%s", err, raw)
		}
		return parsed
	}

	// Arm A — this lane, as the gate will see it.
	armA := probe(t, "")
	// Known-positive control FIRST: if check.passed is false the run is broken and an
	// available=false verdict would be attributed to the wrong cause.
	if !armA.Check.Passed {
		t.Fatalf("instrument failure: ai-check on world/contracts.ail did not pass its CHECK phase; "+
			"the solver verdict below would be meaningless (%+v)", armA)
	}
	if !armA.Verify.Available {
		t.Fatalf("NO SMT SOLVER IN THIS LANE: ai-check reports verify.available=false, so every " +
			"contract vanishes from verify.results[] (V20/V27) and scripts/verify_ail.sh cannot " +
			"reach `verify gate PASSED`. Install the pinned Z3 (see the workflow-level Z3_VER/Z3_SHA " +
			"and the two `Install Z3` steps in .github/workflows/ci.yml), or locally: " +
			"brew install z3 / apt install z3, or set AILANG_Z3_PATH.")
	}
	if armA.Verify.Verified < 1 {
		t.Fatalf("solver is available but verified %d contracts in world/contracts.ail, want >=1 (%+v)",
			armA.Verify.Verified, armA)
	}

	// Arm B — the negative control. Assembled, never a rig-absolute literal (TestNoRigAbsolutePaths).
	missing := filepath.Join(t.TempDir(), "nosuch", "z3")
	armB := probe(t, missing)
	if armB.Verify.Available {
		t.Fatalf("instrument is NOT discriminating: verify.available stayed true with AILANG_Z3_PATH "+
			"pointed at %q, so arm A's green proves nothing about this lane", missing)
	}
	if armA.Verify.Available == armB.Verify.Available {
		t.Fatalf("arms do not differ (A=%v B=%v) — the probe is not reading the solver",
			armA.Verify.Available, armB.Verify.Available)
	}
}
```

### E7 — `ail_binary_gate_test.go`: make the leg-3 assertion unconditional

`TestFixtureDiscrimination` `:458-460`. Replace:

```go
	if strings.Contains(out, "── Leg 3") && !strings.Contains(out, "wrong compiler version") {
		t.Fatalf("fixture discrimination: reached leg 3 but not via the compiler-version check\n%s", out)
	}
```

with:

```go
	// Unconditional since 9/OD-11. Before the Z3 install this was guarded by `Contains(…, "── Leg 3")`
	// because the solverless go-verify job never got past leg 1 — measured: `── Leg 3` count 0 there
	// vs 1 locally, so in CI the guarded body had NEVER executed. With a solver in both jobs the leg
	// is always reached and the guard would only hide a regression.
	if !strings.Contains(out, "── Leg 3") {
		t.Fatalf("fixture discrimination: gate never reached leg 3 — it stopped earlier, so the "+
			"compiler-identity check below was never exercised (a Z3-less lane looks exactly like "+
			"this; see TestSolverAvailableInThisLane)\n%s", out)
	}
	if !strings.Contains(out, "wrong compiler version") {
		t.Fatalf("fixture discrimination: reached leg 3 but not via the compiler-version check\n%s", out)
	}
```

Also update this test's doc comment `:434-440`: the sentence *"The precise `wrong compiler version`
reason lives at leg 3, which is only reachable where Z3 is installed, so it is asserted as a
strengthening when the gate got that far rather than as a portability requirement"* is now false.
Replace from *"The precise"* to the end with: *"The precise `wrong compiler version` reason lives at
leg 3; since 9/OD-11 both CI jobs install Z3, so leg 3 is reachable everywhere this package runs and
the reason is asserted unconditionally."*

### E8 — `ail_binary_gate_test.go`: the ci.yml pin-structure guard

```go
// TestZ3PinDeclaredOnceAndInstalledInBothJobs guards the structure 9/OD-11 bought.
//
// Two claims, and the second is the one with teeth: (1) the Z3 version and sha256 are declared
// EXACTLY ONCE, at workflow scope, so the two jobs cannot install different solvers — a job-local
// re-declaration would shadow the pin for that job only, which is 9/CF-A-2's "two lanes resolve
// different things" all over again; (2) BOTH jobs actually install it. Claim 2 is why this test is
// not decoration: if job 2's install is dropped, every accept-arm's `verify gate PASSED` becomes
// unreachable and 10 tests red with a leg-1 contract error that names no cause. This names the cause.
func TestZ3PinDeclaredOnceAndInstalledInBothJobs(t *testing.T) {
	path := filepath.Join(repoRoot, ".github", "workflows", "ci.yml")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	// Known-positive control: the scan must be able to see strings that ARE present, or a
	// zero-count assertion below could be satisfied by reading the wrong file.
	for _, control := range []string{"ailang-verify:", "go-verify:", "./scripts/verify_go.sh"} {
		if !strings.Contains(src, control) {
			t.Fatalf("instrument failure: %s does not contain known-positive control %q", path, control)
		}
	}
	for _, tc := range []struct {
		needle string
		want   int
		why    string
	}{
		{"Z3_VER:", 1, "the version must be declared exactly once, at workflow scope"},
		{"Z3_SHA:", 1, "the sha256 must be declared exactly once, at workflow scope"},
		{"Z3_VER=", 0, "a job-local Z3_VER= shadows the workflow pin for that job only"},
		{"Z3_SHA=", 0, "a job-local Z3_SHA= shadows the workflow pin for that job only"},
		{`sudo install -m 0755 "${Z3_DIR}/bin/z3" /usr/local/bin/z3`, 2, "BOTH jobs must install the solver"},
	} {
		if got := strings.Count(src, tc.needle); got != tc.want {
			t.Errorf("ci.yml: count(%q)=%d, want %d — %s", tc.needle, got, tc.want, tc.why)
		}
	}
}
```

### Executor warnings (read before typing)

1. **`TestNoRigAbsolutePaths` scans this very file** for the assembled needles `/tmp` + `/ailang`,
   `/Us` + `ers/`, `/home` + `/runner/`. Nothing you add — code **or comment** — may contain those
   substrings. This bit iter-68 twice. `/usr/local/bin/z3` is safe. Never write `/tmp/ailang-v0300`.
2. New imports needed: `encoding/json`. `bytes`, `os`, `os/exec`, `filepath`, `strings` are present.
3. Test files are not built by `go build ./...`. Type-check them with **`go vet ./host/verifygate/`**
   (rc=0 at base, verified).
4. `verify gate PASSED` already occurs **once** in this file at base — inside the comment you are
   replacing (`:113`). Do not treat a count of 1 as "already asserted".
5. Never `git checkout --` anything. Restores are `cp` from your own backup (§7).

---

## 4. Acceptance criteria — command, baseline at base, expected after

`R` = repo root. `PIN=/tmp/ailang-v0300/ailang`. `NOZ3` = a path that does not exist; create it as
`NOZ3="$(mktemp -d)/nosuch/z3"` (do **not** hardcode a `/tmp/…` literal into any source file).
**Every "base" column below was measured first-party by the planner on `541258e` unless marked
`[executor must record]`.** Run the base column FIRST, before any edit (rule 3e).

| AC | Command | **BASE** | **AFTER** | where verifiable |
|---|---|---|---|---|
| **AC0** | `git status --porcelain` then `git diff --name-only` (read-only; controller runs the diff) | clean | **exactly 2** paths: `.github/workflows/ci.yml`, `host/verifygate/ail_binary_gate_test.go` | LOCAL |
| **AC1** *(the defect)* | `AILANG_BIN=$PIN AILANG_Z3_PATH="$NOZ3" go test ./host/verifygate/ -count=1` | **rc=0, `ok … 5.197s`** — the accept-arms pass in a solverless lane | **rc=1**, and the failure names `verify gate PASSED` absent + points at `ci.yml` | LOCAL |
| **AC2** *(no regression)* | `AILANG_BIN=$PIN go test ./host/verifygate/ -count=1` (ambient solver) | rc=0, `ok … 14.482s` | **rc=0** | LOCAL |
| **AC3** *(the probe discriminates)* | `AILANG_BIN=$PIN go test ./host/verifygate/ -run TestSolverAvailableInThisLane -count=1 -v` | test does not exist → `no tests to run`, rc=0 | **rc=0**, `--- PASS: TestSolverAvailableInThisLane` | LOCAL |
| **AC4** *(the probe names the cause)* | `AILANG_BIN=$PIN AILANG_Z3_PATH="$NOZ3" go test ./host/verifygate/ -run TestSolverAvailableInThisLane -count=1` | n/a | **rc=1**, output contains `NO SMT SOLVER IN THIS LANE` and `.github/workflows/ci.yml` | LOCAL |
| **AC5** *(leg 3 unconditional)* | `AILANG_BIN=$PIN AILANG_Z3_PATH="$NOZ3" go test ./host/verifygate/ -run TestFixtureDiscrimination -count=1` | **rc=0** (the `&&` guard makes it vacuous) | **rc=1**, message `gate never reached leg 3` | LOCAL |
| **AC6** *(pin structure)* | `AILANG_BIN=$PIN go test ./host/verifygate/ -run TestZ3PinDeclaredOnceAndInstalledInBothJobs -count=1` | does not exist | **rc=0** | LOCAL |
| **AC7** *(counts, first-party)* | `grep -c 'Z3_VER:' ci.yml; grep -c 'Z3_SHA:' ci.yml; grep -c 'Z3_VER=' ci.yml; grep -c 'Z3_SHA=' ci.yml; grep -c 'sudo install -m 0755' ci.yml` | `0 0 1 1 1` | **`1 1 0 0 2`** | LOCAL |
| **AC8** *(workflow still valid)* | `actionlint .github/workflows/ci.yml` | rc=0 (verified; binary at `/opt/homebrew/bin/actionlint`) | **rc=0** | LOCAL — if `command -v actionlint` is empty in your sandbox, record `actionlint: ABSENT, deferred to controller` and do **not** fake it |
| **AC9** *(type-check)* | `go build ./... && go vet ./host/verifygate/` | rc=0 (verified) | **rc=0** | LOCAL |
| **AC10** *(stale comment gone)* | `grep -c 'deliberately does NOT assert' host/verifygate/ail_binary_gate_test.go` and `grep -c 'installs no Z3' host/verifygate/ail_binary_gate_test.go` | `1` and `1` | **`0` and `0`** | LOCAL |
| **AC11** *(repo gates)* | `AILANG_BIN=$PIN ./scripts/verify_ail.sh` ; `AILANG_BIN=$PIN ./scripts/verify_go.sh` | `verify_ail.sh` rc=0 (verified) · `verify_go.sh` **[executor must record]** | **both rc=0** | LOCAL |
| **AC12** *(job 2 really has a solver)* | GitHub Actions run for the PR head SHA: job `go host build + test gate` step `Install Z3 …` present and **success**, `z3 --version` printing `Z3 version 4.16.0`; job step log for `go build + test gate` contains `ok …/host/verifygate` | job 2 has no such step | present + success | **CI ONLY** |
| **AC13** *(job 1 survived the hoist)* | Same run, job `ailang-code verify gate`: `Install Z3 …` step **success**, all steps green, `failed=0` | green today | green | **CI ONLY** |
| **AC14** *(the claim only CI can settle)* | Same run, job 2 conclusion `success` with `host/verifygate` `ok` — i.e. the accept-arms reached `verify gate PASSED` **on ubuntu-latest**, which is the whole point of the milestone | job 2 green **without** the assertion | job 2 green **with** it | **CI ONLY** |

### Local vs CI — stated explicitly, because this is the item's spine

**LOCAL (macOS, z3 4.16.0 present) settles:** AC0–AC11. These prove the *code* is right — the
assertion fires when it should, is quiet when it should be, and the workflow is well-formed.

**LOCAL CANNOT settle, and no local green may be quoted for it:** AC12, AC13, AC14. Whether
`ubuntu-latest` job 2 actually acquires a solver, and whether the accept-arms then reach
`verify gate PASSED` **there**, is a claim about a runner this rig is not. `AILANG_Z3_PATH="$NOZ3"`
is a *model* of a solverless lane, not the lane. Iter-68's spine — *a green proves the tree passes
where you RAN it, never where it MUST* — is exactly this milestone's residual risk, and the only
instrument that settles it is the CI step log on the PR head SHA, read step-by-step (not a
green tick). **The controller must read that log before merging.**

**Pre-registered CI risks** (declare now, so a red is diagnosed rather than re-run):
- **W1** `unzip` absent on `ubuntu-latest` — job 1 already uses it and is green, so it is present.
- **W2** the workflow-level `env:` not reaching a `run:` body → job 1 reds with
  `Z3_VER: unbound variable`. Loud, immediate, and the fallback in §0 applies.
- **W3** a Z3 4.16.0 asset URL/sha change would red **both** jobs identically — which is the
  hoist working as designed, not a new failure mode.

---

## 5. Wall-clock cost added to job 2, against `timeout-minutes: 25`

Measured on this rig (§1): `host/verifygate` goes **5.197 s → 14.482 s**, i.e. **+9.3 s** per
`go test ./...` pass. `scripts/verify_go.sh` runs the suite **twice** (`:108` plain, `:112-116`
`-race -timeout 8m`), so the delta inside the timed step is **≈ +19 s locally**. Assume
`ubuntu-latest` is 2–3× slower on this workload: **≈ +40–60 s** against a **1500 s** budget — about
**3–4%**. Affordable, and it does **not** approach the ceiling that `4e/OD-4` guards.

The Z3 install itself (~50 MB download + unzip, ~10–20 s judging by job 1) is a **separate step**
and therefore **outside** the `timeout-minutes: 25` on the `verify_go.sh` step entirely. It adds to
job 2's total, not to the guarded budget. **Do not raise `timeout-minutes`.** If the step ever
overruns, route it to `4e/OD-4` — never silently raise the ceiling.

Also note: `-race` and this package are compatible today (base rc=0 via `verify_go.sh` at HEAD); the
gate runs are subprocesses, so the added time is wall-clock, not race-detector overhead.

---

## 6. Mutations — non-vacuity, 7 named

**Universal protocol, applied to EVERY mutation. No step is optional.**

```
# 0. backup (ONCE, before any mutation)
cp host/verifygate/ail_binary_gate_test.go /tmp/vlb_bak_test.go
cp .github/workflows/ci.yml               /tmp/vlb_bak_ci.yml
shasum -a 256 host/verifygate/ail_binary_gate_test.go .github/workflows/ci.yml   # record PRE

# 1. anchor count — assert it equals the number in the table BEFORE editing.
#    A count that differs is INSTRUMENT FAILURE: stop, re-derive, do not proceed.
# 2. apply the exact edit.
# 3. shasum -a 256 <file>   # record POST; POST != PRE or the edit did not land — instrument failure.
# 4. compile the mutant:  go build ./... && go vet ./host/verifygate/   → rc MUST be 0.
#    A mutant that does not compile is INSTRUMENT FAILURE, not a kill. (Iter-68 banked one of these
#    and had to re-run it — a malformed mutant is instrument failure, never "survived".)
# 5. KILL arm  — the -run-scoped command in the table, with its expected rc.
# 6. INVERSE arm — the -skip-scoped command, expected rc=0. This proves YOUR test is the killer and
#    not a bystander. An inverse arm that reds means the mutant broke something else: instrument
#    failure, re-scope the mutation.
# 7. restore:  cp /tmp/vlb_bak_test.go host/verifygate/ail_binary_gate_test.go   (or the ci.yml one)
#    shasum -a 256 <file>  → MUST equal PRE, byte-identical. NEVER `git checkout --`.
```

Prefer a **parsing neuter** (`if false && …`) or a **string edit**. **Never a deletion.**

Common prefix: `E="AILANG_BIN=$PIN"`. `SKIPSET` (the Z3-sensitive tests, excluded from inverse arms
that run solverless) = `'TestKnownPositiveDelegates|TestReleaseChangeNotice|TestSolverAvailableInThisLane|TestFixtureDiscrimination'`.

| # | name | file | anchor (count) | exact edit | KILL arm → expected | INVERSE arm → expected |
|---|---|---|---|---|---|---|
| M1 | `MUT-VLB-PASSED-MARKER` | test | `passedMarker = "verify gate PASSED"` (1) | → `passedMarker = "verify gate PASSED XX"` | `env $E go test ./host/verifygate/ -run 'TestKnownPositiveDelegates' -count=1` → **rc=1**, msg names `verify gate PASSED XX` | `env $E go test ./host/verifygate/ -skip $SKIPSET -count=1` → **rc=0** |
| M2 | `MUT-VLB-PASSED-NEUTER` | test | `if !strings.Contains(out, passedMarker) {` (1) | → `if false && !strings.Contains(out, passedMarker) {` | `env $E AILANG_Z3_PATH="$NOZ3" go test ./host/verifygate/ -run 'TestKnownPositiveDelegates' -count=1` → **rc=0 on the MUTANT** vs **rc=1 on the ORIGINAL** (AC1). The *difference* is the kill: it proves the new assertion, and nothing else, is what reds a solverless lane. | `env $E AILANG_Z3_PATH="$NOZ3" go test ./host/verifygate/ -skip $SKIPSET -count=1` → **rc=0** (a solverless lane reds no bystander) |
| M3 | `MUT-VLB-PROBE-A` | test | `if !armA.Verify.Available {` (1) | → `if false && !armA.Verify.Available {` | `env $E AILANG_Z3_PATH="$NOZ3" go test ./host/verifygate/ -run TestSolverAvailableInThisLane -count=1` → **rc=0 on the MUTANT** vs **rc=1 on the ORIGINAL** (AC4) | `env $E AILANG_Z3_PATH="$NOZ3" go test ./host/verifygate/ -skip $SKIPSET -count=1` → **rc=0** |
| M4 | `MUT-VLB-PROBE-B` | test | `env = append(env, "AILANG_Z3_PATH="+z3)` (1) | → `env = append(env, "WORLD_VLB_UNUSED="+z3)` (arm B now inherits the working lane) | `env $E go test ./host/verifygate/ -run TestSolverAvailableInThisLane -count=1` → **rc=1**, msg `instrument is NOT discriminating` | `env $E go test ./host/verifygate/ -skip $SKIPSET -count=1` → **rc=0** |
| M5 | `MUT-VLB-FIXDISC` | test | `"── Leg 3"` (1, after E7) | → `"── Leg 3 XX"` | `env $E go test ./host/verifygate/ -run TestFixtureDiscrimination -count=1` → **rc=1**, msg `never reached leg 3` | `env $E go test ./host/verifygate/ -skip $SKIPSET -count=1` → **rc=0** |
| M6 | `MUT-VLB-PIN-DUP` | ci.yml | `Z3_VER=` (0, after E1–E3) | insert one line `          Z3_VER=4.15.0` immediately after `          set -euo pipefail` in **job 2's** Z3 step | `env $E go test ./host/verifygate/ -run TestZ3PinDeclaredOnceAndInstalledInBothJobs -count=1` → **rc=1**, msg `count("Z3_VER=")=1, want 0` | `env $E go test ./host/verifygate/ -skip 'TestZ3PinDeclaredOnceAndInstalledInBothJobs' -count=1` → **rc=0** |
| M7 | `MUT-VLB-PIN-DROP` | ci.yml | `sudo install -m 0755 "${Z3_DIR}/bin/z3" /usr/local/bin/z3` (2, after E1–E3) | in **job 2's** step only, change the target to `/usr/local/bin/z3x` (an EDIT, not a deletion) | same `-run` → **rc=1**, msg `count(…)=1, want 2` | same `-skip` → **rc=0** |

For M6 and M7 the "mutant must compile" step (§6 step 4) is **`actionlint .github/workflows/ci.yml`
rc=0**; if actionlint is absent in the sandbox, record `ABSENT` and note that the mutant is a
one-line, syntactically-local change. Both mutants must still be restored byte-identically.

**Why there is no sha-drift mutation:** see §0. The hoist makes the pin single-valued, so there are
no two values to disagree. What is guarded instead is *re-introducing* a second declaration (M6) and
*losing* an install (M7).

**Known asymmetry, declare it — do not overstate precision.** M2 and M3 are *survival-direction*
mutations: the mutant makes a red go away rather than making a green go red. That is the correct
shape for an assertion whose failure condition is an environment, and the ORIGINAL-vs-MUTANT rc pair
is the evidence. M1, M4, M5, M6, M7 are ordinary kill-direction. Also, carrying `9/CF-VLA-1`
forward: M1's `-run 'TestKnownPositiveDelegates'` is *one* of four call sites that would red
(`TestReleaseChangeNotice` has three more), so the `killed_by` mapping is 1-of-N, not 1:1. Say so in
the report; do not claim unique attribution.

---

## 7. Execution order, and what the executor may NOT do

1. `cd $R`; record `git status --porcelain` (must be clean) and the two base sha256s.
2. Run **every AC's BASE column** and record verbatim rc + observable. **`AC1` base rc=0 is the
   defect** — if it comes back rc=1, STOP and report: the premise is refuted and this plan is wrong.
3. `cp` both files to `/tmp/vlb_bak_*` (§6 step 0).
4. Apply E1 → E8 in order.
5. Run AC0–AC11. All must show the AFTER column.
6. Run M1 → M7 with the full universal protocol. Record for each: anchor count, PRE sha, POST sha,
   compile rc, kill rc + message, inverse rc, restore sha == PRE.
7. Re-run AC2, AC9, AC11 on the restored tree to prove the restores were clean.
8. **STOP. Report.** Hand the controller: the AC table with base/after both filled, the 7-row
   mutation table, and any deviation with a stated reason (a self-reported deviation is better
   evidence than a silent one; the controller will adjudicate it by measurement in both directions).

**The executor MUST NOT:** run any `git` write command; create a branch; open a PR; edit any file
other than the two named; edit `tools/launchd/*`; touch `~/.ailang/state/mission-v1*`; modify
`design_docs/world-mission.md`; raise `timeout-minutes`; or use `git checkout --` to restore.
**The CONTROLLER makes all commits, opens the PR, and reads the CI step log for AC12–AC14.**

---

## 8. Sizing, risk, and what is deliberately NOT in this milestone

**Size**: 2 files. `ci.yml` ≈ +18/−2. `ail_binary_gate_test.go` ≈ +105/−14 (one 4-line assertion,
two rewritten comments, one `blocked` entry, two new tests ≈ 75 lines). Net ≈ **+123/−16** — well
inside this mission's per-milestone range (313–2790 for item 8; `VL.A` was +538/−25). **One small
milestone. No split.**

**Risk**: LOW-MEDIUM, entirely concentrated in AC12–AC14 (the CI-only claims) and W2 (the hoist).
Every failure mode identified is loud, not silent.

**Deliberately NOT here** (do not let scope creep in):
- `9/CF-VLA-2` (`defer os.Remove` not crash-safe in the two temp-mutant-script tests) — unchanged.
  Adjacent, but it is a different defect and would grow the diff.
- `9/CF-VLA-3` (a synthetic `v0.33.0-105` is accepted as a release) — unchanged; it is a predicate
  question, not a solver question.
- Pinning job 1's `releases/latest` — **forbidden**, `9/OD-10` clause (a) is ratified ACCEPT.
- Changing `:64`/`:118`'s v0.30.0 pin — **forbidden**, `9/OD-10` clause (c).
- Making `requirePinned` skippable — **forbidden**; its fail-loud posture is the point.

**On `9/OD-11`'s closure**: this milestone discharges it. On merge the controller should mark the
registry row **CLOSED — implemented by `VL.B`**, and note that item 9 was already COMPLETE, so
`VL.B` is a *post-completion strengthening* of a closed item, not a re-opening. `[NEXT]` after this
remains item 5 `w-mcp-projection`'s transition-registry prerequisite, still absent at HEAD.

---

## Verification Log — executor, 2026-08-11

Pinned compiler confirmed before all gate/test commands: `/tmp/ailang-v0300/ailang version` printed
`AILANG v0.30.0` (rc=0). The no-solver control was a generated nonexistent path ending in
`/nosuch/z3`; PATH was not manipulated. Command output was redirected to `/tmp/vlb_*.log` and each
exit status was read directly after the command, never through a pipe.

### Acceptance criteria

| AC | Result | Base observation | After / restored observation |
|---|---|---|---|
| AC0 | **FAIL** | The two implementation files matched the documented base SHA256, but status already contained `?? design_docs/planned/w-verify-binary-lockfile-vlb-sprint-plan.md`. | `git diff --name-only` lists exactly the two implementation paths; `git status --porcelain` also lists this authorized, still-untracked plan. |
| AC1 | **FAIL** | rc=0, package `ok`, 5.168s (defect reproduced). | rc=0, package `ok`, 15.289s; expected rc=1 is unreachable because E4 strips the command's ambient `AILANG_Z3_PATH` before `runGate` starts the gate. |
| AC2 | **PASS** | rc=0, 14.732s. | rc=0, 13.955s; post-mutation restore rerun rc=0, 13.252s. |
| AC3 | **PASS** | rc=0 with `[no tests to run]`. | rc=0; `--- PASS: TestSolverAvailableInThisLane` (0.04s). |
| AC4 | **FAIL** | n/a. | rc=0, not rc=1: E6 removes ambient `AILANG_Z3_PATH`, then arm A calls `probe(t, "")` and deliberately adds no replacement. Therefore neither required failure marker can appear. |
| AC5 | **FAIL** | rc=0, 0.830s (guard vacuous). | rc=0, 2.722s, not rc=1: E4 strips the ambient no-Z3 control and the installed solver reaches leg 3. |
| AC6 | **PASS** | rc=0 with `[no tests to run]`. | rc=0. |
| AC7 | **PASS** | counts `0 0 1 1 1`. | counts `1 1 0 0 2`. |
| AC8 | **PASS** | `actionlint` present at `/opt/homebrew/bin/actionlint`, rc=0. | rc=0. |
| AC9 | **PASS** | `go build ./...` rc=0; `go vet ./host/verifygate/` rc=0. | both rc=0 before mutations and after restoration. |
| AC10 | **PASS** | stale-comment counts `1 1`. | counts `0 0`. |
| AC11 | **FAIL** | `verify_ail.sh` rc=0; `verify_go.sh` rc=1. | `verify_ail.sh` rc=0; `verify_go.sh` rc=1. The latter stops at its explicit guard: active `go1.26.4` miscompiles `host/store/scan.go`; it requests `GOTOOLCHAIN=go1.25.6`. This was not a sandbox denial. |
| AC12 | **CI-ONLY** | Not verifiable locally. | Not claimed; controller must inspect the Ubuntu PR step log. |
| AC13 | **CI-ONLY** | Not verifiable locally. | Not claimed; controller must inspect the Ubuntu PR step log. |
| AC14 | **CI-ONLY** | Not verifiable locally. | Not claimed; controller must inspect the Ubuntu PR step log. |

### Mutations

All test mutants compiled with `go build ./...` rc=0 and `go vet ./host/verifygate/` rc=0.
Workflow mutants passed `actionlint` rc=0. Every post-SHA differed from its pre-SHA, every inverse
arm returned rc=0, and each restoration returned byte-identically to pre (`0c9e70cc…` for the test
file, `5424cada…` for the workflow).

| Mutation | Result | Anchor count | pre → post SHA prefix | compile rc | kill rc | inverse rc | Observation |
|---|---|---:|---|---:|---:|---:|---|
| M1 PASSED-MARKER | **KILLED** | 1 | `0c9e70cc` → `b65d3276` | 0 | 1 | 0 | Names `verify gate PASSED XX` absent. This is a 1-of-N call-site kill, not unique attribution. |
| M2 PASSED-NEUTER | **SURVIVED** | 1 | `0c9e70cc` → `7f0dbe8b` | 0 | 0 | 0 | Original AC1 is also rc=0 after E4, so the required original-vs-mutant difference does not exist. |
| M3 PROBE-A | **SURVIVED** | 1 | `0c9e70cc` → `2a7f6d66` | 0 | 0 | 0 | Original AC4 is also rc=0 because E6 strips ambient no-Z3 before arm A. |
| M4 PROBE-B | **KILLED** | 1 | `0c9e70cc` → `25dbe05b` | 0 | 1 | Names `instrument is NOT discriminating`. |
| M5 FIXDISC | **KILLED** | specified anchor 2; re-scoped executable anchor 1 | `0c9e70cc` → `4d2aea50` | 0 | 1 | Names `gate never reached leg 3`. The plan's raw-string count was 2 because its mandated comment also contains the literal. |
| M6 PIN-DUP | **KILLED** | 0 | `5424cada` → `29eae413` | 0 | 1 | Names `count("Z3_VER=")=1, want 0`. |
| M7 PIN-DROP | **SURVIVED** | 2 | `5424cada` → `d9388eea` | 0 | 0 | 0 | `/usr/local/bin/z3x` still contains the guard's `/usr/local/bin/z3` substring, so `strings.Count` remains 2. |

### Deviations and findings

1. The plan was pre-existing and untracked, so AC0's demanded clean baseline was impossible without
   an unauthorized destructive/removal action. Both implementation files nevertheless matched the
   documented base hashes exactly.
2. E4/E6 and AC1/AC4/AC5/M2/M3 contradict each other: the edits deliberately remove the ambient
   variable that those acceptance and mutation commands rely on. The prescribed edits were retained;
   failed criteria and survived mutations are reported rather than hidden by weakening isolation.
3. M5's specified raw anchor count is 2, not 1. Mutation was paused and re-scoped to the unique
   executable conditional (count 1) before applying the exact code-string change.
4. M7's edit is not detected because the tested needle is a prefix of the mutant target. The mutant
   is honestly reported SURVIVED.
5. `verify_go.sh` was run exactly with the pinned AILANG binary and returns rc=1 on both base and
   restored trees due to its active-Go-version refusal (`go1.26.4`, wants `go1.25.6`). No existing
   test or timeout was changed to conceal it.

---

## Controller verification log (iteration 69, outside the sandbox)

The executor reported **3 of 7 mutations SURVIVED** (M2, M3, M7) and 5 ACs FAIL. It was right, it
was honest, and its self-reported deviations are what made the repair cheap. All three survivals
trace to **one root**: E6 adds `AILANG_Z3_PATH` to `runGate`'s `blocked` map, so the Z3-absent
control can never be armed from the AMBIENT environment — which is exactly how AC1/AC4/AC5 and
M2/M3 were written. The plan contradicted itself across the file boundary between its edits and its
drills.

**The repair is structural, not a weakening.** An ambient drill proves a mechanism once, on a tree
that no longer exists; a committed control proves it on every CI run forever.

| Repair | What it closes |
|---|---|
| `checkProceeded` split out of `requireProceeded` as a pure, error-returning predicate | the accept contract can now be pointed at output the suite chooses — the only way to exercise it against a FAILING gate |
| `TestAcceptContractRejectsASolverlessGate` (new, committed) | arms `AILANG_Z3_PATH` through the `env` map (which `blocked` does not filter), re-measures that all THREE legacy markers are satisfied by a failing gate, then requires the current contract to REJECT it |
| `TestZ3PinDeclaredOnceAndInstalledInBothJobs` install count made LINE-EXACT | `strings.Count` on a prefix-shaped needle could not see a suffix-shaped mutation |
| declared residual on `TestSolverAvailableInThisLane` arm A | measured-unkillable, named in code rather than assumed away |

### Mutations re-run after repair (full six-part drill, `cp` backups, byte-identical restores)

| Mutation | Before | After | Anchor | sha256 pre → post | Compiles | Kill rc | Inverse rc |
|---|---|---|---:|---|---:|---:|---:|
| M2 `PASSED-NEUTER` (`if false && !strings.Contains(out, passedMarker)`) | SURVIVED | **KILLED** | 1 | `15551c521d6d` → `a0cf51decefe` | rc=0 | 1 | 0 |
| M7a `PIN-REDIRECT` (job 2 installs to `/usr/local/bin/z3x`) | SURVIVED | **KILLED** | 2 | `5424cadadfff` → `d9388eea8564` | n/a (yaml) | 1 | 0 |
| M7b `PIN-DROP` (job 2's install step deleted) | not run | **KILLED** | 2 | landed (`sudo install` 2 → 1) | n/a (yaml) | 1 | — |
| M3 `PROBE-A` | SURVIVED | **SURVIVED, DECLARED** | 1 | — | rc=0 | 0 | 0 |

M3 is unkillable **by construction**: arm A's assertion is a diagnostic that fires only in a
solverless lane, and `probe()` deliberately strips ambient `AILANG_Z3_PATH` so each arm resolves the
solver the way its LANE does. Rule 3j allows this when it is declared in the code and the log; it is
now declared in both. Arm B (M4, KILLED) plus the arms-differ check protect the pair.

**Instrument failure caught, not banked:** M7's first anchor read `0` from a pattern nested inside
`"$( … )"`. The mutation had in fact landed (sha differed, the mutated line was visible, and the
guard's own count moved 2 → 1). Re-derived two ways at the top level: **2**. An unexplained zero on
an anchor assertion is instrument failure, never survival.

### Gates, outside the sandbox, `AILANG_BIN=<pinned v0.30.0>` and `GOTOOLCHAIN=go1.25.6`

| Gate | rc | Observed |
|---|---:|---|
| `./scripts/verify_ail.sh` | 0 | `verify gate PASSED: 4 required identities verified, 14 named tests pass` |
| `./scripts/verify_go.sh` | 0 | 32 `ok`, **0** `FAIL` |
| `go vet ./host/verifygate/` · `go build ./...` | 0 · 0 | — |
| `actionlint .github/workflows/ci.yml` | 0 | — |
| `go test ./host/verifygate/ -v` | 0 | **20/20 PASS**, including both new tests |

**AC11 corrected.** The executor read `verify_go.sh` rc=1 on **both base and mutated trees** and
correctly refused to attribute it to the diff: the script rejects the ambient `go1.26.4` and asks
for `GOTOOLCHAIN=go1.25.6`. That is rule 3e(a) — a gate already red at base measures the repo, not
the change. With the toolchain exported it is rc=0.

**AC1 / AC4 / AC5 are SUPERSEDED, not failed.** They specified the ambient-drill control that E6
makes unreachable. `TestAcceptContractRejectsASolverlessGate` answers the same question as a
committed, always-running assertion.

### What only CI can settle (AC12–AC14) — no local green may be quoted for these

Whether ubuntu-latest job 2 actually acquires a solver, and whether the accept-arms reach
`verify gate PASSED` there, is a claim about a machine this rig is not. It is settled by reading the
job-2 step log on the merge commit, and nowhere else. That is the previous milestone's spine:
**a green proves the tree passes where you RAN it, never where it MUST.**
