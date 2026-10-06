# Sprint plan — w-mcp-effects-unrecorded-is-untyped (row 136, iteration 237)

**Planner-Lane:** claude-opus-5-5 (sprint-planner role, iteration 237)
**Authority:** `D-WORLD-66` = A (Mark, attended 2026-10-06): ship option C; option D is out of
scope (row 155). The design's round-2 fixes were applied verbatim at `0ce6e43`.
**Design:** [w-mcp-effects-unrecorded-is-untyped.md](w-mcp-effects-unrecorded-is-untyped.md).
It is binding. §5 below lists every place where this plan's measurements disagree with it.
**Base:** branch `sprint/row136-mcp-typed-errors`, HEAD `0ce6e43` on `origin/dev` `1bf36bd`.
`git diff 3d965f6 HEAD -- host/projection host/coordinator host/daemon host/store` is empty, so
every **code** anchor in the design still holds. Docs moved (§2).
**Prototype:** M1 and M2 were prototyped in throwaway worktrees, gated, and mutation-tested.
Both worktrees are deleted. The diffs are kept under
`~/.ailang/state/mission-world-iter237-evidence/row136-planner/`:
- `proto_m1.diff`: the M1 slice as first built.
- `proto_m2.diff`: M2 on top of M1.
- `proto_final.diff`: M1+M2 plus the §5 fixes. 31 files, +355/−102.

That directory also holds every log cited below. **The executor starts from `proto_final.diff`.**
It splits the diff into the two commits this plan defines and re-runs every gate. It is not a
skeleton to rewrite.

Command prefix for everything below:

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang                        # v0.41.0 GATE pin (not the v0.52.1 tool pin)
export TOOL=$HOME/.pinned-ailang-tools/v0.52.1/ailang                # tool binary, for the se-tools e2e only
E=$HOME/.ailang/state/mission-world-iter237-evidence/row136-executor  # bank logs here; /tmp is wiped
```

## 1. Verification Log (planner, measured this iteration)

| # | Command | Observed |
|---|---|---|
| P1 | Base `0ce6e43`: `go vet ./...` | rc=0 |
| P2 | Base: `AILANG_BIN=… go test -race -count=1 -v ./host/... ./cmd/...` | rc=0. **2605** `=== RUN`, 2551 PASS, 0 FAIL, **54 SKIP**, 27 `ok` |
| P3 | `grep -rn 'WORLD_TOOL_AILANG_BIN\|TestSeTools\|TestCLIAgainstSeTools' .github/ scripts/` | **0 hits.** CI never sets the tool binary. Every `TestSeTools*` test and `TestCLIAgainstSeToolsDaemon` SKIPs in CI and in P2 (§5 D1) |
| P4 | M1 prototype: vet + the P2 command | vet rc=0. Test rc=1: **2606** RUN, 2550 PASS, 1 FAIL, 55 SKIP. The FAIL was `host/broker` `TestExecHandlerTimeoutKillsTheGroup` (`handlers_exec_m2_test.go:398`, a timed-out exec whose `echo started` never reached the partial head). It ran while planner mutants loaded the machine, in a package M1 does not touch. Isolated it passed 3/3. A `host/broker` package re-run at the M1 tree passed (rc=0, 627 RUN, 0 FAIL). **A load flake, not M1** |
| P5 | M1 prototype, with the tool binary: `-run '^(TestSeToolsMCPPostEffectFailureIsLogged\|TestSeToolsA2APostEffectFailureIsMapped)$'` | both `--- PASS` (2 RUN). The line reads `ailang-worldd: mcp refusal: tools/call a2a:ep1:<64hex>: "coordinator: invocation a2a:ep1:… requested effects (records [sha256:…]) …"` |
| P6 | M2 prototype (before the §5 fixes D2–D4): vet + the P2 command | vet rc=0. Test rc=0: **2610** RUN (+5 vs base), 2554 PASS, 0 FAIL, **56 SKIP** (+2: the two new se-tools e2e tests SKIP without `TOOL`) |
| P7 | M2 prototype, with the tool binary: `-race -v ./host/daemon/ ./host/projection/ ./host/coordinator/ ./cmd/...` | rc=0. 932 RUN, 924 PASS, 0 FAIL, 8 SKIP (srt/pty/scale tests only). `TestSeToolsMCPPostEffectFailureIsLogged`, `TestSeToolsSurfaceTaggedInLog`, `TestSeToolsA2APostEffectFailureIsMapped` and `TestCLIAgainstSeToolsDaemon` each print `--- PASS` |
| P8 | M2 prototype **before** the `cliwalk_e2e_test.go:232` edit, with the tool binary | `TestCLIAgainstSeToolsDaemon` **FAIL** at `cliwalk_e2e_test.go:234` `no invocation id in the chain`. The design missed this site (§5 D2), and P6 cannot see it |
| P9 | Final prototype: `AILANG_BIN=… PATH=$HOME/.pinned-ailang:$PATH ./scripts/verify_ail.sh` | rc=0: `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 444 se-tools named tests pass` (no `.ail` changes) |
| P10 | Final prototype, after the §5 companion tests: `go vet ./...`; `-race -v ./host/coordinator/ ./host/projection/` | vet rc=0. rc=0, 355 RUN, 0 FAIL |
| P11 | Mutants: §4 table | 22 mutant runs (`mutants_m1.txt` 4, `mutants_m2.txt` 14, `mutants_ci.txt` 4); every AC mutant KILLED once D1 and D3 are applied |

## 2. Anchors re-measured at `0ce6e43` (design cites `3d965f6`)

Every code anchor is unchanged. Where the design was off or incomplete:

| Design cite | At HEAD | Note |
|---|---|---|
| `coordinator.go:199` (`committed()` compares ids) | `:200` | off by one |
| `mcp-and-a2a.md:204–207` (warning) | `:204–208` | `:::warning` at 204, closing `:::` at 208 |
| `provenance.md:45,50,91` | `:40, :45, :50, :91, :100, :101` | `:40` header JSON, `:100–101` the `why` output. All come from `ailang-worldd call`, which is MCP (`cmd/ailang-worldd/mcpclient.go:247` posts `/mcp/`), so after M2 they read `coordinator:mcp` / `mcp:` |
| `reference/http-api.md` (no line) | `:103, :114` | example `provenance` / `writtenBy` |
| `getting-started/coding-tools.md` (no line) | `:262` | "Tool-call entries have `writtenBy: "coordinator:a2a"`" |
| `website/static/AGENTS.md` (no line) | `:165` | "MCP and A2A calls are both logged with …" |
| (not cited) | `README.md:240` | `why <…|a2a:<id>|->` usage: add `mcp:<id>` (and the missing `rest:<id>`) |
| (not cited) | `cmd/ailang-worldd/cliwalk_e2e_test.go:232` | regex `` `invocation (a2a:\S+)` `` over a `call` (MCP) result. **Must** become `mcp:` in M2 (§5 D2) |
| "callers … tests" of `InvocationID` | `coordinator_test.go:165, 647, 771, 818, 990, 1044`; `effectful_test.go:379`; `inflight_test.go:58`; `projection_test.go:885, 996`; `absence_test.go:165`; `mcp_conformance_test.go:127, 353`; `why_test.go:123` | **13 test sites.** `mcp_conformance_test.go:127, :137 (TrimPrefix), :353` take `SurfaceMCP`. `inflight_test.go:58` takes `call.Surface`. All others are A2A |
| `planInvocation` test caller | `coordinator_test.go:193` (`lawPlan`) | takes `SurfaceA2A` |
| AC1.3 "dispatchError table (~800–848)" | `projection_test.go` `TestA2ADispatch` `:804`; rows `:836–841`; assertion `:847` | |
| AC2.5 `commit_budget_test.go:356` | `TestCommitInvocationReceiptIdentity` `:296`, subtest `foreign_namespaces_are_refused`, loop `:356` | |
| AC2.6 `coordinator_test.go:818`, `:990` | `TestNewRefusesMissingSeamsAndCaps` `:796` (assert `:818`); `TestReconcileRefusesDamagedRecord` `:894` (assert `:990`) | |
| AC2.1 `mcp_conformance_test.go:127,137` | the `assertReceipts` helper `:113`, used by `TestMCPBatchAccounting` `:155` and `TestMCPBatchConformance` `:185`; also `TestMCPAbsentPrecheckSuccessfulInvocation` `:353` | |
| V15 `docs/QUICKSTART.md:573` (pi `/mcp tools`) | `:573` before M1; `:574` after M1 adds a line | expected shift |

## 3. Milestones

Two code commits, each CI-green and mergeable alone, plus one docs commit for M3. Follow the
repo convention: one PR, with the commits kept separate.

### M1 — MCP operator line: correct label, tested refs

**Production (2 files):**
- `host/projection/projection.go:500–501`: `logRefusal(method, id, err)` becomes
  `logRefusal(surface, method, id string, err error)` with format
  `"ailang-worldd: %s refusal: %s %s: %q\n"`. The three A2A call sites pass `"a2a"`: `:334`,
  `:341` and `:374`. **A2A bytes are unchanged.**
- `host/projection/mcp.go`:
  - `:76` and `:81`: `logRefusal("mcp", "tools/list", "-", err)`.
  - `:107`: `logRefusal("mcp", "tools/call", "-", err)`.
  - `:132`: `logRefusal("mcp", "tools/call", coordinator.InvocationID(binding.EpisodeID, task), err)`.

  The cause stays `err` (`%q`, one line).

**Tests:**
- `host/daemon/setools_e2e_test.go`: pull the R14 rebinding out of
  `TestSeToolsA2APostEffectFailureIsMapped` (`:725–746`) into `func (r *seRig) serveR14(errLog io.Writer)`.
  The A2A test calls it with `io.Discard`, so its assertions do not change. Add
  `TestSeToolsMCPPostEffectFailureIsLogged` (**AC1.1**). Its `ErrorLog` is a mutex-guarded writer
  (`syncLog`), because the hostcall goroutine writes it. It asserts:
  1. The wire is exactly `-32603 "host callback failed"`.
  2. The log is exactly one line and starts `ailang-worldd: mcp refusal: tools/call a2a:ep1:` (const `mcpRefusalPrefix`).
  3. **(§5 D3)** The line's id, captured by `` ^ailang-worldd: mcp refusal: tools/call (\S+): ``,
     appears in the cause as `coordinator: invocation <id> requested effects`.
  4. Some `sha256:` ref on the line resolves through `GetObject` to `broker.EffectRecordV1`, with
     `Allowed` and `Effect == Workspace.Read`.
  5. `entryCount` is unchanged.
- `host/projection/mcp_conformance_test.go`: add `TestMCPRefusalLineNamesSurfaceAndCause`
  (**AC1.1a**, §5 D1). It is the binary-free companion that CI actually runs. It uses
  `batchFixture` with `runner.failAt = 1` and `h.errorLog` set to a locked sink. It asserts one
  line starting `ailang-worldd: mcp refusal: tools/call ` + `InvocationID(<ep-a>, s.tasks[0]) + ": \""`
  that carries `batch injected capsule failure`. At M1 the id is the 2-argument `InvocationID`;
  M2 switches it to `SurfaceMCP`.

**Docs (same commit; the glm r2 fix):**
- `docs/QUICKSTART.md:92–93`: add the MCP line after the A2A line:
  `` `ailang-worldd: mcp refusal: tools/call <invocation-id-or->: "<escaped cause>"` (or `tools/list -`) ``.
- `website/docs/guides/operating-the-daemon.md:98–99`: "A2A and MCP refusal causes", with both lines.
- `website/docs/guides/troubleshooting.md:187`: "an `a2a refusal` line (A2A) or an `mcp refusal` line (MCP)".
- `website/docs/reference/mcp-and-a2a.md:174`: append "(an `mcp refusal: tools/call` line over MCP)".

**M1 gate (all must pass; bank every log in `$E`):**
1. `./scripts/verify_ail.sh` with `AILANG_BIN` set and the pin first on `PATH`: rc=0.
2. `go vet ./...`: rc=0.
3. `AILANG_BIN=$AILANG_BIN go test -race -count=1 ./...`: rc=0.
4. **Tool-binary leg (§5 D1; a SKIP is a fail):**
   ```sh
   T="TestSeToolsMCPPostEffectFailureIsLogged TestSeToolsA2APostEffectFailureIsMapped TestSeToolsMCPEndToEnd TestCLIAgainstSeToolsDaemon"
   AILANG_BIN=$AILANG_BIN WORLD_TOOL_AILANG_BIN=$TOOL go test -race -count=1 -v ./host/daemon/ ./cmd/ailang-worldd/ \
     -run "^(${T// /|})\$" > $E/m1_tool.log 2>&1; echo rc=$?
   for t in $T; do grep -q -- "^--- PASS: $t " $E/m1_tool.log || echo "MISSING PASS: $t"; done
   ```
   Use `|`, never `\|`: an RE2 `-run 'A\|B'` matches zero tests.
5. V15 docs sweep: `grep -rnE "a2a refusal|mcp invoke|mcp tools" website docs`. Every remaining
   hit is an A2A line or pi's `/mcp tools` slash command (prototype: 7 hits, 3 of them the
   slash-command lines `connecting.md:79`, `coding-tools.md:225`, `QUICKSTART.md:574`). No hit
   reads `mcp invoke`.

### M2 — Surface tag in the invocation id and the entry header

**Production:**
- `host/coordinator/coordinator.go`:
  - Above `Call` (`:80`), add `type Surface string` with `SurfaceA2A = "a2a"` and `SurfaceMCP = "mcp"`.
  - `Call` gains `Surface Surface` as its first field.
  - `:106–108`: `InvocationID(surface Surface, episodeID, taskID string) string` returns
    `string(surface)+":"+episodeID+":"+taskID`. Update the doc comment.
  - `validateCall` `:122`: the first check refuses any surface other than A2A/MCP, **including
    empty**, as `&InvalidCallError{Field: "surface"}`.
  - `:250`: `InvocationID(call.Surface, …)`.
  - `:312`: `dispatchEffectful(ctx, call.Surface, id, …)`.
  - `:332`: `planInvocation(world, call.Surface, id, …)`.
  - Comment `:224`: "coordinator worlds" (was `"a2a:" worlds`).
- `host/coordinator/effectful.go`:
  - `:248`: `dispatchEffectful(ctx, surface Surface, id, episodeID string, …)`.
  - `:275` and `:308`: `planEffectInvocation(world, surface, id, …)`.
- `host/coordinator/plan.go`:
  - Delete the const `writtenBy` at `:17`. Add `func writtenBy(surface Surface) string { return "coordinator:" + string(surface) }`.
  - `object(by, semanticID, payload)` at `:80–84`.
  - `planInvocation` (`:99`) and `planEffectInvocation` (`:114`) take `surface Surface` after `w`
    and compute `by := writtenBy(surface)`.
  - `planCommit(w, by, id, …)`, with `WrittenBy: by` at `:135`.
  - Run `gofmt`: the const block re-aligns.
- `host/projection/projection.go`: `:363` `coordinator.Call{Surface: coordinator.SurfaceA2A, …}`;
  `:374` `InvocationID(coordinator.SurfaceA2A, …)`.
- `host/projection/mcp.go`: `:130` `coordinator.Call{Surface: coordinator.SurfaceMCP, …}`;
  `:132` `InvocationID(coordinator.SurfaceMCP, …)`.
- `host/daemon/handlers.go:353`: comment `"a2a:"/"mcp:"`.
- `cmd/ailang-worldd/why.go`:
  - Delete `coordinatorAuth` (`:39`, zero references) and re-align the const block.
  - `:8` comment `a2a:/mcp:`.
  - Usage `:276`: `a2a:<id> | mcp:<id> | rest:<id>`.
  - `:317`: also accept `strings.HasPrefix(arg, "mcp:")`.
  - `:327`: the error text `a2a:/mcp:/rest:<id>`.
- `cmd/ailang-worldd/main.go:29, :92`: add `mcp:` to the `why` usage.

**Call-literal sweep, re-measured at HEAD:** `grep -rnE '\bCall\{' host cmd --include='*.go'`
gives 6 hits. Five are composite literals and all five get an explicit `Surface`:
- `projection.go:363` (A2A)
- `mcp.go:130` (MCP)
- `coordinator_test.go:155` (the `rig.call` helper, A2A)
- `workspace_test.go:524` (A2A)
- `commit_budget_test.go:154` (A2A: **add** `Surface`, do not change the expected `DeadlineExceeded`)

`effectful_test.go:640` is a `[]Call{…}` slice built from the helper. The prototype confirms
there are no other literals.

**Test edits:**
- The 13 `InvocationID` test sites in §2.
- `coordinator_test.go:193`: `planInvocation(w, SurfaceA2A, "a2a:e:t", …)`.
- `coordinator_test.go:818`: `InvocationID(SurfaceA2A, "ep", "task") != "a2a:ep:task"` (AC2.6).
- `mcp_conformance_test.go:137`: `TrimPrefix(id, "mcp:ep-a:")`.
- `cliwalk_e2e_test.go:232`: `` `invocation (mcp:\S+)` `` (§5 D2).
- `setools_e2e_test.go`: `mcpRefusalPrefix` changes from `…tools/call a2a:ep1:` to `…tools/call mcp:ep1:` (**AC2.1b**).
- AC1.1a's `InvocationID` takes `coordinator.SurfaceMCP`.
- `why_test.go:108–138`: `invoke` becomes a wrapper over a new
  `invokeAs(surface coordinator.Surface, …)`. `invokeAs` sets `by := "coordinator:" + string(surface)`
  in place of the five `"coordinator:a2a"` literals, and its id is `InvocationID(surface, episode, task)`
  (§5 D4).
- `commit_budget_test.go:356`: add `"mcp:ep:task"` to the foreign-namespace list (**AC2.5**).

**New tests:**
- `coordinator_test.go` `TestDispatchRefuses` (`:343`): add the rows `R1_empty_surface`
  (`c.Surface = ""`) and `R1_unknown_surface` (`c.Surface = "rest"`). Each expects
  `InvalidCallError{Field:"surface"}` and `assertUntouched` (**AC2.2**). The design names a new
  `TestDispatchRefusesUntaggedSurface`. Rows in the existing R1 table reuse its `assertUntouched`
  store check, so the prototype used rows (§5 D5). Either form is acceptable if the mutant dies.
- `coordinator_test.go` `TestPlanWrittenByNamesSurface` (**AC2.3a**, §5 D1). It is pure and runs
  in CI. For each of `SurfaceA2A` and `SurfaceMCP`, through both `planInvocation` and
  `planEffectInvocation`, the header `WrittenBy` and every object's `Provenance` equal
  `coordinator:<surface>`.
- `setools_e2e_test.go` `TestSeToolsSurfaceTaggedInLog` (**AC2.3**): one MCP call, then one A2A
  call (`a2a-tag-1`) on one rig, read through a `seRig.entryTag(n)` helper.
  - Entry `first` must be `coordinator:mcp` with an id starting `mcp:ep1:`.
  - Entry `first+1` must be `coordinator:a2a` with id `a2a:ep1:a2a-tag-1`, and the A2A
    `metadata.invocation_id` must equal that id.
- `why_test.go` `TestWhyResolvesMCPInvocationID` (**AC2.4**): an A2A and an MCP invocation share
  `ep1`/task `x`. Then:
  - `why mcp:ep1:x` resolves to entry 2, with all links ✓ and `writtenBy coordinator:mcp`.
  - `why a2a:ep1:x` resolves to entry 1.
  - `why bogus:ep1:x` exits `exitUsage` with "is not an entry index".

**Docs (S7, same commit).** Every MCP-produced example flips:
- `website/docs/agents/provenance.md`:
  - `:40` and `:100` become `coordinator:mcp`; `:50` and `:101` become `mcp:ep1:…`.
  - `:45`: replace the "does not say which surface" sentence. `writtenBy` names the surface:
    `coordinator:mcp` for MCP `tools/call` (and `ailang-worldd call`), `coordinator:a2a` for A2A.
    Entries written before this tag say `coordinator:a2a` for both.
  - `:91`: "your `mcp:` or `a2a:` invocation id".
- `concepts/sessions-and-episodes.md:16`: `mcp:<episode>:<task-id>` (MCP) or `a2a:…` (A2A).
- `guides/provenance-walks.md:79`: name both forms.
- `reference/cli.md`:
  - `:37` adds `mcp:<id>`.
  - `:254` becomes `` `mcp:<id>` or `a2a:<id>` ``.
  - `:284–286` become `coordinator:mcp` (the walk drives `call`).
- `reference/http-api.md:103, :114`: `coordinator:mcp`.
- `getting-started/coding-tools.md:262`: `coordinator:mcp` (an A2A call writes `coordinator:a2a`).
- `website/static/AGENTS.md:165`: state both tags and both id forms.
- `README.md:240`: the `why` usage (§2).
- Leave unchanged:
  - `mcp-and-a2a.md:144`: a genuine A2A example.
  - The frozen evidence under `design_docs/verification/`.

**M2 gate:** steps 1–3 of M1, then:

4. **Tool-binary leg:** add `TestSeToolsSurfaceTaggedInLog` to `T` in M1 step 4. Every name must
   print `--- PASS`. `TestCLIAgainstSeToolsDaemon` is the guard for §5 D2.
5. M2 docs sweep: `grep -rn 'coordinator:a2a\|a2a:<' website docs README.md`. Every remaining hit
   must name A2A explicitly or describe pre-tag stores. Prototype: 7 hits, all of that kind.
6. `grep -rnE '\bCall\{' host cmd --include='*.go'` must still give exactly 6 hits, with 5
   literals carrying `Surface:`. This is an instrument-health control only; AC2.2 is the
   load-bearing check.

### M3 — Upstream ask; the client-visible typed error stays gated

The **controller** owns filing the `sunholo-data/ailang` issue (area:serveapi, V1–V5 as the
repro, the ask exactly as in design §M3) and the `ailang messages send mission-control` note.
The body is positional; assert the artifact afterwards. The controller also owns the follow-up
charter row with the verbatim Gate-1 predicate. **The executor files nothing.**

**Executor (docs-only commit):** replace `website/docs/reference/mcp-and-a2a.md:204–208` with:

```markdown
:::warning EffectsUnrecorded is only named over A2A
The effect-record refs of an `EffectsUnrecordedError` reach the caller over `/a2a/`. Over
`/mcp/` the caller sees only `host callback failed`: the released MCP handler has no typed host
error yet (upstream ask: <ISSUE-URL>). Until it ships, the operator recovers the refs from the
daemon's stderr. The line `ailang-worldd: mcp refusal: tools/call mcp:<episode>:<task>: "…"`
names the invocation id and, in its quoted cause, every effect record (`records [sha256:…]`).
Check that line and the episode's effect records before retrying. An MCP retry is a new
invocation and runs its effects again.
:::
```

`<ISSUE-URL>` stays a literal placeholder in the executor's commit, and the controller replaces
it. **The M3 commit must not merge while `<ISSUE-URL>` is still in the file.** Gate: `grep -c
'<ISSUE-URL>' website/docs/reference/mcp-and-a2a.md` is `1` before the controller edit and `0`
before merge.

## 4. Acceptance criteria, tests and mutants (prototype results)

Each mutant is applied with `perl -0pi`. A run first checks the mutant actually applied
(`cmp` ≠), then runs the named test(s), then restores the file. Scripts: `$E/../row136-planner/mut.sh`
(with the tool binary) and `mut_ci.sh` (AILANG_BIN only, as CI runs). Results are in
`mutants_m1.txt`, `mutants_m2.txt` and `mutants_ci.txt`.

| AC | Test | Mutant | Result (first assertion that fired) | Runs in CI? |
|---|---|---|---|---|
| AC1.1 | `TestSeToolsMCPPostEffectFailureIsLogged` | **MUT-MCP-LOG-NOCAUSE** (`:132` logs `errors.New("dispatch failed")`) | KILLED. At M1 (before D3): `setools_e2e_test.go:825` "operator line carries no effect-record ref". With D3: the id-in-cause check fires first | **no** (needs `TOOL`) |
| AC1.1 | same | **MUT-MCP-LOG-LABEL** (MCP passes `"a2a"`) | KILLED, `:821` prefix assertion | no |
| AC1.1a | `TestMCPRefusalLineNamesSurfaceAndCause` | MUT-MCP-LOG-LABEL / MUT-MCP-LOG-NOCAUSE | both KILLED, `mcp_conformance_test.go:403` | **yes** |
| AC1.2 | `TestA2ADispatchWire` (`a2a_wire_test.go:193`, `assertRefusalLog` `:494–496`) | **MUT-A2A-LABEL** (A2A `:374` passes `"mcp"`) | KILLED, 3 subtests, `a2a_wire_test.go:212` | yes |
| AC1.3 | `TestA2ADispatch` rows `EffectsUnrecorded*` (`projection_test.go:836–841`) | **MUT-UNRECORDED-ORDER** (the `unrecorded` case moved below `conflict`) | KILLED, `projection_test.go:847` | yes |
| AC2.1 | `TestMCPBatchAccounting` / `TestMCPBatchConformance` / `TestMCPAbsentPrecheckSuccessfulInvocation` | **MUT-MCP-SURFACE** (a) Call + log id → `SurfaceA2A`; (b) Call literal only | both KILLED, `mcp_conformance_test.go:166` "task/receipt identity mismatch" (9 RUN) | yes |
| AC2.1b | AC1.1 prefix, now `mcp:ep1:` | MUT-MCP-SURFACE (a) | KILLED, `:821` | no |
| AC2.1b | same | MUT-MCP-SURFACE **(b)**, Call literal only | **SURVIVED without D3** (the logged id is derived separately from the dispatched one). **KILLED with D3** at `:826` | no |
| AC2.2 | `TestDispatchRefuses/R1_empty_surface` | **MUT-SURFACE-DEFAULT** (`validateCall` defaults `""` to A2A) | KILLED, `coordinator_test.go:402` `err = <nil>` | yes |
| AC2.3 | `TestSeToolsSurfaceTaggedInLog` | **MUT-WRITTENBY-CONST** (`WrittenBy: "coordinator:a2a"`) | KILLED, `setools_e2e_test.go:892` | no |
| AC2.3 (CI) | every CI-visible test in `host/coordinator`, `host/projection`, `host/daemon`, `cmd/ailang-worldd` | MUT-WRITTENBY-CONST | **SURVIVED** (878 RUN, rc=0) without AC2.3a | — |
| AC2.3a | `TestPlanWrittenByNamesSurface` | MUT-WRITTENBY-CONST; also the object variant `Provenance: "coordinator:a2a"` | both KILLED, `coordinator_test.go:208` / `:212` | **yes** |
| AC2.4 | `TestWhyResolvesMCPInvocationID` | **MUT-WHY-PREFIX** (drop `mcp:` from `parseTarget`) | KILLED, `why_test.go:259` | yes |
| AC2.5 | `TestCommitInvocationReceiptIdentity/foreign_namespaces_are_refused` | **MUT-REST-ACCEPTS-MCP** | KILLED, `commit_budget_test.go:358` `status = 200, want 400` | yes |
| AC2.6 | `TestNewRefusesMissingSeamsAndCaps` `:818` + `TestReconcileRefusesDamagedRecord` `:990` | **MUT-A2A-ID-FORMAT** (surface after episode) | KILLED, `:823` and 6 reconcile subtests | yes |
| (D2) | `TestCLIAgainstSeToolsDaemon` | `cliwalk_e2e_test.go:232` left at `a2a:` | KILLED (fails at `:234`) | **no** |
| AC3.1 | none: evidence (issue URL + `ailang messages` id in the iteration log) | n/a, labelled instrument | — | — |

The evaluator re-runs every mutant from a clean tree. No survivors remain once D1 and D3 are
applied.

## 5. Design deviations (measured; none changes direction)

- **D1. The design's gate cannot see its own e2e ACs (S6).** AC1.1, AC2.1b and AC2.3 live in
  `setools_e2e_test.go`. That file SKIPs unless `WORLD_TOOL_AILANG_BIN` is set (`seToolBins`,
  `:57–66`), and CI never sets it (P3). P2 and P6 SKIP them under the design's gate
  (`go test ./...` with `AILANG_BIN`). MUT-WRITTENBY-CONST survived every CI-visible test (§4).
  - **This plan:** add two binary-free companions, AC1.1a and AC2.3a (both prototyped, both kill
    their mutants with AILANG_BIN only), and an explicit tool-binary gate leg with PASS-line
    checks.
  - **Controller:** CI coverage of the se-tools e2e is a separate gap: there is no tool binary in
    CI. Recommend a new row, not work inside 136.
- **D2. The M2 sweep missed `cmd/ailang-worldd/cliwalk_e2e_test.go:232`.** Its regex
  `` `invocation (a2a:\S+)` `` reads a `why` walk of an `ailang-worldd call` result, which is MCP.
  The M2 prototype without the edit fails `TestCLIAgainstSeToolsDaemon` (P8). That test also
  SKIPs in CI (D1), so the gate in this plan's M2 step 4 is the only guard. `README.md:240` was
  missed as well (§2).
- **D3. AC2.1b does not kill every form of MUT-MCP-SURFACE as the design writes it.** `mcp.go:132`
  re-derives the logged id with its own `InvocationID(...)` call, separate from the `Call{Surface}`
  at `:130`. A mutant that changes only the dispatched surface leaves the log line reading `mcp:`
  and AC1.1 stays green (measured: SURVIVED). The fix in this plan is AC1.1 assertion 3: the line's
  id must equal the id the coordinator names in the cause. After that fix the mutant is KILLED.
  AC2.1 (`mcp_conformance`) already killed both forms.
- **D4. AC2.6's "why_test.go:112–122 byte-unchanged" cannot hold literally.** `why_test.go:123`
  calls `InvocationID`, whose signature changes, and AC2.4 needs the same fixture to mint
  `coordinator:mcp`. The A2A fixture still commits byte-identical objects: `invoke` delegates to
  `invokeAs(SurfaceA2A, …)` with the same strings. Every existing why/log test passes unchanged
  (P6/P7), including `TestLogTailFollowPrintsEachCommitOnce`'s `coordinator:a2a ep1 ailang-read`
  check at `:370`. Read AC2.6 as "the old-shape store bytes are unchanged", not "the source text
  is unchanged".
- **D5. AC2.2 is implemented as rows of the existing `TestDispatchRefuses` R1 table**, not as a
  new function, so it reuses `assertUntouched`. The mutant dies either way (§4).
- **D6. Doc anchors moved** (§2). `provenance.md`, `cli.md` and `coding-tools.md` examples come
  from `ailang-worldd call` (MCP), so M2 flips them to `coordinator:mcp`. The design's sweep grep
  finds them, but its prose implied only the "cannot tell" sentence changes.

No premise of V1–V14 was contradicted. The code anchors, the R14 injection, the reachability of
`dispatchError` and the `rest:`-only `/v1/commit` all measured as stated.

## 6. Executor notes

- Run the full `-race` suite with no other heavy job on the machine. P4's `host/broker` flake
  happened under planner load. If an untouched package goes red, re-run it alone 3× before
  attributing it.
- `go build ./...` is not a compile fence for `_test.go` files. Use `go vet ./...`.
- When splitting `proto_final.diff` into M1 and M2, the M1 slice is `proto_m1.diff` plus D3's
  assertion and AC1.1a. At M1 the prefix is `a2a:ep1:` and AC1.1a uses the 2-argument
  `InvocationID`. Gate the M1 tree on its own before building M2 on it.
- Bank all logs under `$E`; `/tmp` is wiped.
