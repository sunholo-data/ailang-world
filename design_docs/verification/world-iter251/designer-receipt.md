# Designer receipt — iteration 251, queue row 26 (`w-bounded-z3-report-producer`)

- **Doc**: `design_docs/planned/w-bounded-z3-report-producer.md` (474 lines, written by the
  scaffold script `~/.claude/skills/design-doc-creator/scripts/create_planned_doc.sh
  w-bounded-z3-report-producer` run from this worktree, then filled; no scaffold section left
  as a placeholder — verified by grep for `[placeholder]`/`[Feature Name]`/`TODO`: 0 hits).
- **Verification Log rows in the doc**: **33** (V1–V33, §10 of the doc). Every negative-existence
  claim carries a same-call known-positive control; every AILANG-language claim is live-checked
  against the pinned binary (`~/.pinned-ailang/ailang`, v0.41.0); the doc's own sweep note maps
  the negative claims to their rows.
- **Evidence artifacts committed beside this receipt** (all under
  `design_docs/verification/world-iter251/`):
  - `gen_aicheck_fixtures.sh` — the reproducible fixture generator (scratch `.ail` modules in a
    `mktemp` dir outside the repo, so no new `.ail` enters a `verify_ail.sh`-swept tree; the
    generator is the same shape the doc's test plan commits at
    `host/evidence/testdata/aicheck/regen.sh`).
  - `fixtures/` — seven real pinned-binary `ai-check` outputs with `.rc`/`.stderr` sidecars:
    `verifiedpass` (rc 0, verified), `refuted` (rc 1, counterexample+model), `encerror` (rc 1,
    errors=1), `mixed` (rc 1, verified+error), `skippedshape` (rc **0**, skipped=1),
    `loaderror` (rc 1, check failed), `nocontracts` (rc **0**, results=[]).

## Commands that surprised me

1. **The silent-failure class changed shape between the pinned releases.** The inherited claim
   ("a Z3 encoding error exits 0 silently", parent V10, measured on v0.30.0) does NOT reproduce
   on the pinned v0.41.0: my encoding-error fixture (string interpolation in `ensures`) exits
   **rc 1**. Instead, a DIFFERENT rc-0 silent class is live: the **skipped** shape
   (`verify.skipped=1, verified=0`, exit 0) — a clean exit code over a contract that was never
   verified. The design keeps the JSON-authoritative discipline and adds `skipped>0` as an
   explicit refusal, citing the live measurement rather than the inherited v0.30.0 claim (doc
   V25, AC27).
2. **The pinned binary needs a writable `TMPDIR` for Z3 solver scratch files.** With the
   variable stripped, the run fails with `solver error: creating temp file: … operation not
   permitted`; with `TMPDIR` redirected into the scratch root the run succeeds and cleans up.
   This is now a measured design fact (doc V26, §D3): the producer's child env redirects
   `TMPDIR` into the removed-on-cleanup scratch root.
3. **A contract-free module is a "clean" run that proves nothing** — rc 0, `check.passed=true`,
   every counter zero, `results: []` (doc V27). This is the anti-vacuity guard's fixture, and it
   is not constructible by hand-waving: it had to be generated to be believed.
4. **Line-number drift in the parent's citations**: `db.SetMaxOpenConns(1)` is at
   `host/store/store.go:382` today, not `:298` as the parent's AC16 cites (doc V9 records the
   drift; the doc cites the measured line). Several of my own first-draft line citations were
   also wrong for the same reason and were re-verified one by one before finalizing
   (validator.go's `ObjectReader` is at :50-53; the report codec's sorted-unique branch at
   report_codec.go:86-97; `beginDurable`'s pool-wait comment at durable.go:48-51).
5. **The rig's default `grep` is ugrep and direct `/tmp` writes are fence-blocked** — recorded
   commands use `/usr/bin/grep` per the directive, and the fixture generator writes its scratch
   through `mktemp "${TMPDIR:-/tmp}/…"` (the per-session TMPDIR is writable) with the committed
   outputs landing inside the worktree.

## Design-time defect caught by the doc's own mutation discipline

The first-draft D4 had TWO vacuity guards with one sentinel (a zero-`results` branch plus a
zero-`verified` branch, both `ErrVacuousCheck`). The mutation drill showed the first branch
would be unkillable — the second intercepts the same fixture with the same sentinel, exactly
the "two guards, one observable" redundancy the mission has rejected before. Fixed in the doc:
ONE merged anti-vacuity guard with a detail that distinguishes its two input shapes, and M34
re-anchored to it (under the mutant the `nocontracts` fixture refuses with the WRONG sentinel,
so the named test reds on sentinel identity).

## Deviations from the directive

None material. Two notes for the record:

- The directive's fixture rule was applied by committing the generator + its real outputs under
  `design_docs/verification/world-iter251/` (evidence home) and by specifying in the test plan
  that the implementation commits the same shapes under `host/evidence/testdata/aicheck/` —
  outside every `verify_ail.sh` sweep, on the existing `host/replay/testdata/transition_fixture.ail`
  precedent (doc V28). No `.ail` file was added anywhere the gate sweeps.
- The estimate is honest arithmetic, not the carried figure: the queue row carries 1.0 d from
  item 17's cut and explicitly prices the round-12 fixes as "unpriced, pending design"; the doc
  prices them (~2.0 d total, §Timeline) and says so in those words.

## Scope discipline

The doc was created and the evidence generated; nothing was implemented, no quorum was run, no
git add/commit/branch/stash was performed. The working tree carries only this document, the
evidence directory, and this receipt (V1 of the doc's log).