# TR.C mutation transcript

All commands used `PATH=/opt/homebrew/bin:$PATH`, `GOTOOLCHAIN=go1.25.6`, and broker tests used
`AILANG_BIN=/tmp/ailang-v0300/ailang`. Each target was copied to `/tmp/trc_backup` before mutation,
the before/after hashes differed, `go build ./...` and `go vet ./host/...` were required before the
kill result, and restoration was byte-identical. The recorded inverse rc=0 selected the binding test
and skipped it. This is only the scoped inverse instrument: the prescribed whole-package inverse is
**DEFERRED / UNINFORMATIVE UNDER SANDBOX**, because the package panics on its first `httptest`
listener with `bind: operation not permitted` even with the known base flake skipped. The controller
must run that arm outside the sandbox. M10's meaningful inequality-only inverse was rc=0. Hashes
below are full SHA-256 values except where the repeated baseline is abbreviated in the display.

| Arm | Before → mutant SHA-256 | Build / vet | Killed by (exact observation) | Inverse | Restore |
|---|---|---:|---|---:|---|
| M1 `FOURTH-INVOKE-IN` | `c6fbdfc7178651d9ee80c5fb796f2184c5c45cb62d074b66130caabcdf962f37` → `aef42b2f299657c04a06fedfd4ed1d0314f3979513195ce32a3b7f1b7d6f3e64` | 0 / 0 | `inside_broker_exemption`: `count=4 want 3`, naming `confined.go:18 kind=invoke-call fn=trcMutationInvokeInside` | 0 | match |
| M2 `FOURTH-INVOKE-OUT` | `068fc0e47a1404e53ff8e38c2e5d3d5515fbd088914e90aa8bf34e53c129746c` → `d66c18196262943da57ff6bb26fd86347ada9bcbc743b2f385c998baaff98b54` | 0 / 0 | `outside_broker_is_clean`: `bind.go:16 kind=invoke-call fn=trcMutationInvokeOutside` | 0 | match |
| M3 `MOVE-SITE` | `9c04777a3505d980f5cc30446bf43bf0b231db53656edecc8d5f69b5bbbe4098` → `ee43988ffaa49e63b603cbded450f9ab76b4478fe24afada3e3426bb2f02a964` | 0 / 0 | `inside_broker_exemption`: identity mismatch names `trcMovedPollInvoke`; count stayed 3 | 0 | match |
| M4 `CTOR-LIVE` | `068fc0e47a1404e53ff8e38c2e5d3d5515fbd088914e90aa8bf34e53c129746c` → `5a46bfde79f60c186c19bb8ab19f287a59ba75489d44fcaf8f3ae62703ee957a` | 0 / 0 | `outside_broker_is_clean`: `kind=ctor-live fn=trcMutationLiveConstructor` | 0 | match |
| M5 `CTOR-REPLAY` | `068fc0e47a1404e53ff8e38c2e5d3d5515fbd088914e90aa8bf34e53c129746c` → `10fa863b0ac359c217fa5635a4b1dd859c651e4911ba6e5ba7153c75bc5eb954` | 0 / 0 | `outside_broker_is_clean`: `kind=ctor-replay fn=trcMutationReplayConstructor` | 0 | match |
| M6 `SESSION-TYPE` | `068fc0e47a1404e53ff8e38c2e5d3d5515fbd088914e90aa8bf34e53c129746c` → `424ca6c0cbcc56a2b0524521041dae36be52b6075ee841f4ec42a314739981f5` | 0 / 0 | `outside_broker_is_clean`: `kind=session-type fn=trcMutationSessionType` | 0 | match |
| M7 `ALIAS-IMPORT` | `3deab0c0c76f537f06e0765b68eb2b3fd7fa8bbc98f79b93b4be19af2f120c6f` → `fac004e24e00cf8d3b201b8af730fbb3e892240e39e9856f3b0d7f25450446c3` | 0 / 0 | `outside_broker_is_clean`: aliased `session-type` in `cmd/ailang-worldd/main.go` | 0 | match |
| M8 `DOT-IMPORT` | `77e743839a5e191a2c1d4f95b1965fc687ce1614b6e0a820fdb3187f955593f5` → `5cb76871e7ddf56ed6a09e87f4757fb4699b65e3700e83065ca7c6517d77376b` | 0 / 0 | `outside_broker_is_clean`: `host/daemon/daemon.go kind=dot-import` | 0 | match |
| M9 `RAISE-COUNT` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `a9c4c14f1f16a7e4db9bf53eec59089200c305efa1b83d31a9da1aa9476af24e` | 0 / 0 | `inside_broker_exemption`: `count=3 want 4` | 0 | match |
| M10 `COUNT-INEQUALITY` + M1 | test `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `45103fe82ea071a4c68f53002e520aa5fe4309c630ed8d8151f8f4a7c6597789`; production `c6fbdfc7178651d9ee80c5fb796f2184c5c45cb62d074b66130caabcdf962f37` → `aef42b2f299657c04a06fedfd4ed1d0314f3979513195ce32a3b7f1b7d6f3e64` | 0 / 0 | `inside_broker_exemption`: identity mismatch names fourth site; weakened count stayed green | 0 (M10 alone) | both match |
| M11 `IF-FALSE-INVOKE` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `3dfddb09ab990422515e315a2a602b8d51955f05891e791f17cd3b384de113f6` | 0 / 0 | `detector_controls/POS-invoke-outside`: `findings=0 want 1` | 0 | match |
| M12 `IF-FALSE-CTOR-LIVE` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `1dd895779e247503778b7b9f44de6ed0d040630f0d6435972bd58c758908098a` | 0 / 0 | `POS-ctor-live`: `findings=0 want 1` | 0 | match |
| M13 `IF-FALSE-CTOR-REPLAY` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `3315e85e8541e98df781ef8d69e480ec30af5a460efa215863ca86a7dcba1f25` | 0 / 0 | `POS-ctor-replay`: `findings=0 want 1` | 0 | match |
| M14 `IF-FALSE-SESSION-TYPE` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `895e5a4ab770d76e84f4b8d7c3525866d7e100aa8fcfa5080f5aa586aa94cd4d` | 0 / 0 | both `POS-session-type` and `POS-alias-session`: `findings=0 want 1` | 0 | match |
| M15 `IF-FALSE-DOT-IMPORT` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `c2b5c070b32513efe0bd35384de1444c8299d85a32c90b4634df867b7e85e190` | 0 / 0 | `POS-dot-import`: `findings=0 want 1` | 0 | match |
| M16 `EMPTY-WALK` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `5884959e80119595be450517b5711ad1c10b31f8c7166df698407f270d301b28` | 0 / 0 | `enumeration`: `walked ZERO production .go files` | 0 | match |
| M17 `WALK-SKIP-BROKER` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `e12952a2e457ab42790a00767420d7d8f5897a91e1d8480a042017772a8ca44e` | 0 / 0 | `enumeration`: `walked only 25 production .go files, want at least 30` | 0 | match |
| M18 `WALK-INCLUDE-TESTS` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `0d3615a6fb547aef17f55f46fafb5d9bcb0c8f045fa34f6406d989a9836981d3` | 0 / 0 | `enumeration`: names `.snap/T1/host/broker/invoke_boundary_test.go` | 0 | match |
| M19 `GOLIST-DEAD` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `938b4ca0b0a19a318c8c5f20616c3266f470ea678977ddf6f3a71c62c2100e5b` | 0 / 0 | `enumeration`: `go list enumerated 0 files, want at least 30` | 0 | match |
| M20 `TEXT-SCANNER` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `579e39293b28928a54d23af965a4e3e00c53f024b5c00ee7525f00fc91a6f729` | 0 / 0 | `inside_broker_exemption`: text count `22 want 3` | 0 | match |
| M21 `NEG-CONTROL-DEAD` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `394e963f798136fa6f754fc69905b111bfc2f849e9df91175353b8085751047c` | 0 / 0 | all four `NEG-*`: `findings=1 want 0` | 0 | match |
| M22 `PARSE-SWALLOW` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `fd52c07c07cd925dbf055561190451575901df901c6481e5d1c2a1e6ce8c1c45` | 0 / 0 | `WALK-unparseable`: `unparseable production file was silently accepted` | 0 | match |
| M23 `DELETE-TR-C-TEST` | `b4f1d997dc92ba9887231349b9f3067486a6b0d3e905eb4b7f0aba0f9689b985` → `e3b0acec2f987f99928fe4bf76910eb38270a98bdf9bb5ffe1eef50da2d015f3` | 0 / 0 | activated AC11 inventory `count=1`, wanted 2 | 0 | match; restored count 2 |

## Uncertain rows, corrections, and deviations

M1 used a referenced function containing a compile-time-dead call. Its inverse is satisfiable and
rc=0; no landed co-detector fired. M20 retained the helper signature and changed only its body, so
both build and vet were rc=0. The live text mutant counted 22 occurrences inside `host/broker`
(rather than the prototype's repository-wide 26) and was killed by the exact-count assertion.

The first attempted M7 did not build because only some package qualifiers were aliased; it was not
banked. The recorded M7 uses the one-reference `cmd/ailang-worldd` import and builds/vets. The first
attempted M8 duplicated the broker import and did not build; it was not banked. The recorded M8 uses
`host/daemon` with a dot import plus a symbol-use control and builds/vets. The first M21 form used an
early return and failed vet; it was not banked. The recorded form appends one finding after the real
detector and builds/vets.

The mandated snapshots increase the live `skipped_tests` measurement: T4 saw 48 rather than the
pre-snapshot base 44. The production walk stayed 39 and go-list stayed 38.

Protocol deviation: for M1–M9 and M11–M23, the table's inverse `0` is the scoped skip control, not
the required whole-package inverse. Attempting the whole-package inverse is uninformative in this
sandbox: `TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce` panics while
binding `[::1]:0`. Those 22 full inverse arms are explicitly DEFERRED to the controller's
outside-sandbox sweep; none is silently represented as satisfied. M10's prescribed composed-vs-
alone inverse is satisfied locally.

## Rule-3j inventory from the actual untracked file

The untracked-safe diff and the `t.Fatal`/`t.Fatalf` cut fire. Reading the file yields detector
mechanisms R1 invoke, R2 live constructor, R3 replay constructor, R4 session type, R5 dot import,
R17 alias resolution; exemption refusals R6 exact count and R7 exact identity; enumeration E1 zero,
E2 floor, E3 anchors, E4 test exclusion/non-vacuity, E5 go-list floor/superset; R14 parse errors and
R16 negative controls. The plan names all of those.

The actual file also has infrastructure refusal branches not named as mutation rows in the plan:
walker filesystem/stat/relative-path errors; `runtime.Caller` failure; missing `go`; go-list timeout
or command failure; go-list relative-path failure; repository parser failure; synthetic-source
parser failure; synthetic fixture mkdir/write failures; nested-module fixture identity mismatch; and
the positive-control kind mismatch. These are reported rather than silently folded into the frozen
23-row table.

## Controller sweep, outside the sandbox (iteration 75)

The executor's 23 arms were re-read, and two things it could not do were done here: the substantive
assertion branches its 23 do not name were mutated, and the 22 whole-package inverse arms it
explicitly DEFERRED were run. Every arm below asserts LANDED (sha256 pair) and BUILDS before any
test result is read; for a test-only milestone the compile gate is `go vet ./host/broker`, because
`go build ./...` does not compile `_test.go` files at all.

**The rule-3j cut, and its negative control.** `git diff --no-index /dev/null
host/broker/invoke_boundary_test.go | grep -cE '^\+.*(t\.Fatal|t\.Errorf|t\.Fatalf)'` returns **29**
(control: 354 added lines). The same cut through ordinary `git diff` returns **0** — the
untracked-file trap iteration 74 recorded, reproduced first-party on the file this milestone adds.

### Four branches the plan's 23 do not name — all KILLED

| Arm | Branch | Mutation | Killed by |
|---|---|---|---|
| C1 | `:216` `skippedTests == 0` | `stats.skippedTests++` → `_ = stats` | `enumeration`: `production walk skipped zero _test.go files; exclusion is vacuous` |
| C2 | `:225` required anchors | walker drops `writer_lock_other.go` | `enumeration`: `production walk missed required anchor host/store/writer_lock_other.go` |
| C3 | `:312` control kind | rename the `invoke-call` literal everywhere | `inside_broker_exemption` `:280` — **NOT** `:312`; the mutation moved observer and observed together, so this arm does not reach the branch it was aimed at |
| C4 | `:312` control kind, ISOLATED | relabel at the PRODUCER only (`:54`), leaving both expectation sites intact | `:312` on all six `POS-*` controls (`kind=wrong-kind want invoke-call`), plus `:280` |

C3 → C4 is the finding worth keeping: an arm that rewrites a constant used by BOTH the mechanism and
the assertion cannot fail for the reason it claims (rule 3i). C3 alone would have recorded `:312` as
covered while leaving it unobserved; only mutating at the producer reaches it.

### The 22 deferred whole-package inverse arms

Negative control first: on the UNMUTATED tree, `go test ./host/broker -skip
'TestRegistryDispatchBindingBoundary|TestHandlerTimeoutKillsTheWholeProcessGroup' -count=1` is
**rc=0, ok 35.403s** — so a red inverse would have meant something.

| Arm | Shape | Kill | Whole-pkg inverse |
|---|---|---:|---:|
| P1 | fourth `Invoke` inside `host/broker` (`confined.go`) | rc=1 `exemption count=4 want 3`, naming `confined.go:75 fn=trcP1` | **rc=0** |
| P2 | `Invoke` selector outside (`host/transitionreg/bind.go`) | rc=1 `outside_broker_is_clean` | **rc=0** |
| P3 | `broker.NewSession` outside | rc=1 `outside_broker_is_clean` | **rc=0** |
| P4 | `broker.Session` type outside | rc=1 `outside_broker_is_clean` | **rc=0** |
| M16-shape | walker appends nothing (test-file-only) | rc=1 `walked ZERO production .go files` | **rc=0** |

Two arms were DISCARDED rather than banked, both for the reason this discipline exists: a first
`P1` used `s.Invoke(ctx, req)` and failed `not enough arguments`, and `P5` (dot-import in
`host/daemon`) failed `undefined: broker` — `host/daemon` does not import the broker at all. A
mutant that does not build reds in exactly the direction you predicted, which is the one case a
negative control agrees with you for the wrong reason.

The remaining test-file-only arms (M9, M11–M15, M17–M23, C1–C4) have a whole-package inverse of
rc=0 **by construction**, and the construction is measured rather than assumed: every helper the
gate defines — `parseAndScan`, `enumerateProductionFiles`, `assertEnumeration`,
`goListProductionFiles`, `repositoryRoot`, `bindingFinding`, `walkStats` — has **0** references
outside `invoke_boundary_test.go`, against a firing control of **35** for `Session`. Skipping the
gate therefore executes none of the mutated code. The M16-shape arm above converts that derivation
into one measurement.

**Result: 23 executor arms + 9 controller arms, ZERO survivals.**

### Gates, all outside the sandbox

`verify_go.sh` rc=0 with **0** `FAIL` and exactly 2 healthy `WARNING: DATA RACE` · `verify_ail.sh`
rc=0, totals **4 identities / 11 modules / 14 tests UNMOVED** · `go vet ./host/...` rc=0 · AC5=2
AC6=3 AC7=3 unmoved, **AC11 1 → 2 ACTIVATED** · `AC-INVOKE3` n=3 p=3 (control 90) · `host/broker`
under `-race` **90.896s** against a 92.3s base. Final gate on a `.snap`-free tree, i.e. the tree CI
sees: `walked=39 skipped_tests=45 skipped_nested_modules=2 golist=38`.
