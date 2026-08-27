# w-toolchain-pin-normalizer-accepts-malformed-gotoolchain — one normalizer serves two grammars, so the strict one inherits the lax one's tolerance: a GOTOOLCHAIN value the Go runtime refuses outright passes the pin gate

**Status**: PLANNED
**Date**: 2026-08-28
**Queue item**: 45, `w-toolchain-pin-normalizer-accepts-malformed-gotoolchain` (clause-2, queued
iter-128; folds in two controller-measured findings from iteration 134)
**Estimated**: ~0.15 day — ONE file, `host/verifygate/toolchain_pin_gate_test.go` (~+55/−20):
split `normalizeToolchainPin` into a per-convention canonicalizer and a strict toolchain-name
validator; flip Test B's two `t.Errorf` instrument floors to `t.Fatalf`; add one
comparison-direction pin on run.sh's pinned-OK guard. Zero `.ail` touched; single-gate (Go)
change; `ci.yml` and `run.sh` ship byte-untouched (they appear below only as mutation venues).
**Designer**: `claude-fable-5` (design-doc-creator, iteration 134)
**Revision**: round-1 quorum applied 2026-08-28 (3/3 present, three REJECTs — see §Quorum
round 1); this is the protocol-mandated single revision, one re-quorum follows.
**Toolchain boundary**: every VERIFIED-BY-ME command below ran first-party in this checkout
(`/Users/voightkampff/dev/sunholo-data/ailang-world`, zsh→bash-safe quoting, darwin/arm64,
macOS 26.5.0) at `07668e1`, tree clean (`git status --porcelain` = 0 lines), 2026-08-28, with
`go version` = `go1.26.6 darwin/arm64`. Rows marked INHERITED were measured by the controller
this iteration (attribution: "controller, iteration 134") in a detached probe worktree at
`07668e1` with every mutation proven LANDED before its result was read and restored to
porcelain-0; the load-bearing non-mutation arms are re-derived here first-party. The
`go/version.IsValid` table ran in a throwaway `/tmp/vercheck` module (nothing entered the
tree). The one network read (Go's own toolchain documentation) is dated and quoted verbatim.

> **Thesis:** `normalizeToolchainPin` answers one question — *what version does this pin
> mean?* — by blindly prepending `go` to whatever lacks it, and that is the correct reading
> for exactly ONE of the three conventions it serves. `actions/setup-go`'s `go-version:` and
> go.mod's `go` directive legitimately omit the prefix; a `GOTOOLCHAIN` env value without it
> is not a laxer spelling of the same pin but a value the Go runtime itself refuses with
> `go: invalid GOTOOLCHAIN` (P3) — and the gate's auto-correction certifies exactly the byte
> string that would kill every `go` command in the job (P4). One normalizer serving two
> grammars silently imports the laxer grammar's tolerance into the stricter one. The fix is
> not a special case bolted onto the shared function — it is to stop sharing: each convention
> gets the validator its own grammar defines, with `GOTOOLCHAIN` validity sourced from Go's
> documentation (P6) and its name-shape checked by `go/version.IsValid`, the stdlib function
> this file already imports and already trusts for KNOWN_BAD tokens (P9, P17) — a first,
> attributing filter for the SYNTAX class only (availability is a different failure class,
> P19), which STOPS the test on failure so the kept agreement-loop and floor comparisons
> never grade a non-pin operand (P21), while those kept comparisons pin any well-formed
> name to the exact go.mod floor string behind it. Two smaller
> instances of the same one-namespace-two-meanings shape ride along, both controller-measured:
> Test B grades instrument failure `Errorf` where Test A and the Z3 precedent grade the same
> class `Fatalf` (P10, P11, P12), and Test B counts the pinned-OK guard's SITES while being
> blind to its comparison DIRECTION (P13, P14).

## The finding in one paragraph

`normalizeToolchainPin` (`host/verifygate/toolchain_pin_gate_test.go:13`) trims, strips
quotes, and prepends `go` to any value not already starting with it; it is called from three
conventions — workflow keyed lines (both `GOTOOLCHAIN` and `go-version`), go.mod `go`
directives, and run.sh's `PINNED=` assignment (P1, P2). For `go-version: '1.26.6'` and
`go 1.26.6` the prefix-less form is the convention's native spelling (P8); for
`GOTOOLCHAIN` it is malformed — `GOTOOLCHAIN=1.26.6 go version` exits 1 with
`go: invalid GOTOOLCHAIN "1.26.6"` while the control arm `GOTOOLCHAIN=go1.26.6 go version`
exits 0 (P3) — yet the controller measured that mutating `ci.yml:21` to exactly that value
leaves `TestGoToolchainPinsAgreeAndMatchJobList` green (P4): the gate whose one job is pin
hygiene certifies a workflow in which every subsequent `go` invocation would refuse to run.
The row names `:21`, but the extractor collects TWO `GOTOOLCHAIN` sites (`:21`, `:102`) and
TWO `go-version` sites (`:28`, `:109`) — one pair per enumerated job (P7) — so the repair must
validate every value the extractor collects, never a line number. A value that fails
validation stops Test A at ONE attributed message — under this doc's round-1 draft
(`Errorf`) the same malformed value produced THREE, two of them asserting a disagreement
between pins when one operand was not a pin at all (P21; quorum round 1, R3). Two adjacent
one-namespace defects fold in: Test B's two `instrument failure:` floors use `t.Errorf`
where Test A's identical floors and ALL seven of the Z3 precedent's use `t.Fatal/Fatalf`
(P10, P11), so a broken instrument there limps through ~10 downstream assertions instead of
stopping at the cause (P12); and Test B's `saw_pinned_ok` site-count assertion (≥3) is
direction-blind — the controller measured that flipping run.sh:140's `-eq 0` to `-ne 0`
preserves the count at 3 and leaves the whole 18-arm test file green (P13, P14).

## Premises

`VERIFIED-BY-ME` = run by the designer at `07668e1`, output observed first-hand this session.
`INHERITED` = measured by the controller, iteration 134, in a detached probe worktree at
`07668e1` (mutation LANDED-proof asserted before each result was read; porcelain 0 after
restore), cited with attribution as instructed; non-mutation arms re-derived here where
load-bearing.

| # | Claim | Command (verbatim) | Observed | Status |
|---|---|---|---|---|
| P1 | The shared normalizer has exactly three call-site conventions | `grep -rn 'normalizeToolchainPin' --include='*.go' .` with same-call control `grep -rn 'moduleGoFloor' --include='*.go' .` | 4 hits: definition `:13`, `pinValues` `:30` (serves BOTH `GOTOOLCHAIN` and `go-version` keys), `moduleGoFloor` `:45` (go.mod `go ` directives), `:244` (run.sh `PINNED=`); control: 4 hits (`:36,:148,:238,:272`) — pattern live | VERIFIED-BY-ME |
| P2 | The normalizer blindly prepends `go` to whatever lacks it | `sed -n '13,23p' host/verifygate/toolchain_pin_gate_test.go` | trim → quote-strip → `if value != "" && !strings.HasPrefix(value, "go") { value = "go" + value }` — no grammar check of any kind | VERIFIED-BY-ME |
| P3 | A prefix-less GOTOOLCHAIN is invalid to the Go runtime itself; control arm differs | probe loop over both arms (see P5 for the full loop): `GOTOOLCHAIN="1.26.6" go version`; `GOTOOLCHAIN="go1.26.6" go version` | `1.26.6 -> rc=1: go: invalid GOTOOLCHAIN "1.26.6"`; `go1.26.6 -> rc=0: go version go1.26.6 darwin/arm64` — the two arms differ | VERIFIED-BY-ME (re-derives controller measurement A, iteration 134, identically) |
| P4 | The gate is blind to the malformation TODAY | controller, iteration 134, probe worktree at `07668e1`: baseline `go test ./host/verifygate/ -run TestGoToolchainPinsAgreeAndMatchJobList -count=1 -v` → rc=0 `--- PASS`; mutate `.github/workflows/ci.yml:21` `GOTOOLCHAIN: go1.26.6` → `GOTOOLCHAIN: 1.26.6` (LANDED: old literal 2→1, new literal 0→1, line-21 content asserted); re-run | mutant **rc=0 `--- PASS`**; restored, `git status --porcelain` = 0 | INHERITED (controller, iteration 134) — this doc's headline defect; becomes RED via mutation M1 below |
| P5 | The measured GOTOOLCHAIN value grammar on the pinned-major toolchain, including the validity-vs-availability split | `for v in auto local path go1.26.6 go1.26.6+auto go1.26.6+path 1.26.6 min1.26.6; do printf '%s -> ' "$v"; GOTOOLCHAIN="$v" go version >/tmp/gt.out 2>&1; printf 'rc=%s: %s\n' "$?" "$(head -1 /tmp/gt.out)"; done` | `auto -> rc=0`; `local -> rc=0: go version go1.26.4` (the bundled toolchain — proof `local` is a MODE, not a pin); `path -> rc=1: go: cannot find "go1.26.6" in PATH` (valid value, unavailable channel — NOT the invalid-value error class); `go1.26.6 -> rc=0`; `go1.26.6+auto -> rc=0`; `go1.26.6+path -> rc=1: cannot find in PATH` (same class as `path`); `1.26.6 -> rc=1: go: invalid GOTOOLCHAIN "1.26.6"`; `min1.26.6 -> rc=1: go: invalid GOTOOLCHAIN "min1.26.6"` | VERIFIED-BY-ME — two distinct failure classes: *invalid value* vs *valid value, toolchain unavailable*; only the first is this row's defect |
| P6 | Go's own documentation of the grammar (the row said "check `go help toolchain`"; that topic does not exist on go1.26.6 — recording what I found instead) | `go help toolchain` → `go help toolchain: unknown help topic. Run 'go help'.`; `go help` topic list contains no toolchain topic (goauth/gopath/goproxy present as adjacent controls); `go help environment 2>&1 \| grep -n -A6 'GOTOOLCHAIN'` → `:83 GOTOOLCHAIN / Controls which Go toolchain is used. See https://go.dev/doc/toolchain.`; fetched that page 2026-08-28 | verbatim: *"It consults the `GOTOOLCHAIN` setting, which takes the form `<name>`, `<name>+auto`, or `<name>+path`. `GOTOOLCHAIN=auto` is shorthand for `GOTOOLCHAIN=local+auto`; similarly, `GOTOOLCHAIN=path` is shorthand for `GOTOOLCHAIN=local+path`. The `<name>` sets the default Go toolchain: `local` indicates the bundled Go toolchain … and otherwise `<name>` must be a specific Go toolchain name, such as `go1.21.0`."* and *"The standard Go toolchains are named `goV` where V is a Go version denoting a beta release, release candidate, or release."* | VERIFIED-BY-ME — the grammar the validator below implements, matching P5's measured arms shape-for-shape |
| P7 | TWO `GOTOOLCHAIN` sites and TWO `go-version` sites exist — one pair per enumerated job; the row's `:21` names half the surface | `grep -n 'GOTOOLCHAIN\|go-version' .github/workflows/ci.yml` | `21:      GOTOOLCHAIN: go1.26.6` / `28:          go-version: '1.26.6'` / `102:      GOTOOLCHAIN: go1.26.6` / `109:          go-version: '1.26.6'` | VERIFIED-BY-ME (matches controller, iteration 134) |
| P8 | The prefix-less form IS the native spelling for the two lax conventions | `grep -n '^go ' go.mod`; `sed -n '19,29p' .github/workflows/ci.yml` | `3:go 1.26.6` (the go.mod directive grammar carries no `go` prefix in its value); the setup-go block reads `uses: actions/setup-go@v5` / `with:` / `go-version: '1.26.6'` — setup-go's documented version-spec form, currently green under Test A (P15) | VERIFIED-BY-ME — prepending is CORRECT for these two; the defect is sharing that correction with GOTOOLCHAIN |
| P9 | `go/version.IsValid` accepts exactly the bare toolchain-name form and rejects every non-pin shape | `/tmp/vercheck` module, `go run .` over 11 shapes | `IsValid("go1.26.6")=true`, `("go1.26")=true`, `("go1.21rc1")=true`, `("go1.99.99")=true`; `("1.26.6")=false`, `("go1.26.6+auto")=false`, `("go1.26.6+path")=false`, `("auto")=false`, `("local")=false`, `("path")=false`, `("goauto")=false` | VERIFIED-BY-ME — the stdlib discriminator for `<name>`; note `go1.99.99`=true: validity ≠ availability (Declared residual 2) |
| P10 | Exactly two `instrument failure` floors in the pin-gate file are non-fatal, both in Test B | `grep -n 'instrument failure' host/verifygate/toolchain_pin_gate_test.go` | 11 hits; `:199` (`does not contain known-positive control`) and `:209` (`exact shebang count`) are `t.Errorf`; the other 9 (`:78,:128,:274,:284,:288,:293,:312,:345,:365`) are `t.Fatalf` — Test A's SAME control-loop floor at `:78` is fatal | VERIFIED-BY-ME (corroborates controller measurement B, iteration 134) |
| P11 | The Z3 precedent Test B's doc-comment claims to mirror grades this class fatal, without exception — verified myself as the row demanded | `grep -n 'instrument failure' host/verifygate/ail_binary_gate_test.go`; `grep -c 't.Fatalf' host/verifygate/ail_binary_gate_test.go`; `grep -c 't.Errorf' host/verifygate/ail_binary_gate_test.go` | 8 hits: `:193,:268,:272,:282,:679` `t.Fatalf`, `:559,:580` `t.Fatal`, `:264` a comment line — 7 call sites, ALL Fatal-class, zero Errorf among them; file-wide 33 `t.Fatalf` vs 10 `t.Errorf` | VERIFIED-BY-ME — the precedent claim is TRUE; Test B departs from it at exactly its two floors |
| P12 | Under `Errorf`, a broken instrument limps through ~10 further assertions instead of stopping at the cause | controller, iteration 134 (measurement B); mechanism corroborated by code read: after `:199`/`:209`, Test B continues through the assignment-count, non-empty, floor-membership, PINNED-equality, exec-bit, site-count, and guard-message assertions (`:212–:268`) | a missing control token cascades into the downstream assertion noise rather than one attributed refusal | INHERITED (controller, iteration 134); corroborated VERIFIED-BY-ME by read |
| P13 | The pinned-OK guard: 3 `saw_pinned_ok` sites; the `-eq 0` guard literal appears exactly once, its flipped form zero times (same-call control pair) | `grep -n 'saw_pinned_ok' design_docs/verification/w-race-gate-blindspot/run.sh`; `grep -c 'saw_pinned_ok" -eq 0' …/run.sh`; `grep -c 'saw_pinned_ok" -ne 0' …/run.sh` | sites `31:saw_pinned_ok=0` (declaration), `87:… [ "$tc" = "$PINNED" ] && saw_pinned_ok=1 ;;` (OK set), `140:if [ "$saw_pinned_ok" -eq 0 ]; then` (guard); literal counts: `-eq 0` form = **1**, `-ne 0` form = **0** — the 1 is the same-call known-positive control for the 0 | VERIFIED-BY-ME |
| P14 | Flipping the guard's direction survives every committed assertion TODAY | controller, iteration 134, probe worktree: mutate run.sh:140 `-eq 0` → `-ne 0` (LANDED: site count before=3, after=3, new literal present); `go test ./host/verifygate/ -run TestMiscompileInstrumentProbesPinnedToolchain -count=1` | **rc=0 `ok`** on the mutant — invisible to all 18 named arms; loud only at RUNTIME (the flipped guard refuses whenever the pin DID report OK, i.e. on every healthy gated run); restored, porcelain 0 | **RE-DERIVED FIRST-PARTY, round-2 R3** (controller, iteration 134, second independent run this session — see P23 for the verbatim command set and observed output) — since row 44 landed, that runtime refusal is a CI red on the NEXT push; the static gap is attribution and pre-push detection (honesty note in Decision) |
| P15 | Baseline: both pin-gate tests green at `07668e1`, tree clean | `go test ./host/verifygate/ -run 'TestGoToolchainPinsAgreeAndMatchJobList\|TestMiscompileInstrumentProbesPinnedToolchain' -count=1 -v`; `git log -1 --format='%h %s'`; `git status --porcelain \| wc -l` | both `=== RUN` + `--- PASS`, `ok … 0.279s`; `07668e1 record(world) iter 133: …`; `0` | VERIFIED-BY-ME |
| P16 | `PINNED=` carries a toolchain NAME by convention, and run.sh compares it byte-wise against toolchain names at runtime | `sed -n '25,31p' design_docs/verification/w-race-gate-blindspot/run.sh` (P13's grep supplies `:87`) | `PINNED="go1.26.6"   # the toolchain go.mod pins; …` — already prefixed; `:87` `[ "$tc" = "$PINNED" ] && saw_pinned_ok=1` compares against `$tc` values drawn from the `goX.Y.Z`-named lists. So a prefix-less `PINNED="1.26.6"` is auto-corrected to static green by P2's function while `:87` could never set the flag — the static gate would certify a script whose own `:140` floor then refuses every run | VERIFIED-BY-ME (mutant arm itself is design-shipped as M5, not measured today) |
| P17 | `go/version` is already imported by this file — the strict validator adds zero dependencies | `sed -n '1,11p' host/verifygate/toolchain_pin_gate_test.go` | import block contains `"go/version"` (used by `TestReproModuleFloorStaysBelowKnownBadToolchains` at `:273,:292,:296,:299`) | VERIFIED-BY-ME |
| P18 | *(amended round 1, R2 — the prior scope was `*.sh` only, which excluded the Go gate machinery; widened to the whole repo)* No competing static validator of the pin TEXT exists anywhere in the repo; every other occurrence of the token is prose, message text, a fence over a different file, dynamic-by-design use, or a runtime observation of the ACTIVE toolchain | `grep -rn 'GOTOOLCHAIN' --include='*.go' .` with same-call control `grep -rn 'normalizeToolchainPin' --include='*.go' .`; whole-repo census `grep -rn 'GOTOOLCHAIN' . --exclude-dir=.git`; prior `*.sh` scope retained | Go hits: exactly **8**, all in `host/verifygate/toolchain_pin_gate_test.go`, per-hit disposition: `:64` doc-comment prose; `:100`/`:101` the job-enumeration Errorf message + its two `pinValues` calls; `:104` the extraction itself; `:108` the keyed-line count floor; `:120` the agreement loop's `toolchain pins disagree` message; `:314`/`:315` `TestCanaryDeclaresPositiveArmOnly`'s `strings.Count(src, "GOTOOLCHAIN") != 0` fence over `canaryPath` — a token fence over a DIFFERENT file (`host/store/toolchain_canary_test.go`), not a value validator. Control: 4 hits (`:13,:30,:45,:244`) — grep live. Census beyond `.go`: code surfaces are `ci.yml` ×2 (P7's pin sites — the SUBJECT), `scripts/verify_go.sh` ×1, `run.sh:69` ×1 (`GOTOOLCHAIN="$tc" go build …`, the probe's intended dynamic use); every remaining hit is `design_docs/**`/sprint-JSON prose. Also in conflict scope and previously unmentioned: `verify_go.sh:214–222` deny-lists `go env GOVERSION` (`# observe; never assign`) — a RUNTIME observation of what ran, whose `:222` message text merely SUGGESTS `GOTOOLCHAIN=go1.25.6`; it validates no pin text. Conclusion at its true strength: **no competing static validator of the pin TEXT exists; `:314` is a token fence over a different file and `verify_go.sh` observes the ACTIVE toolchain, so neither conflicts** | VERIFIED-BY-ME (re-derived first-party this pass; matches controller, iteration 134) |
| P19 | Name-shaped-but-unusable GOTOOLCHAIN values fail in the AVAILABILITY class, never the invalid-value class — `IsValid` discriminates syntax, not availability; and `+auto` on the installed floor name RUNS at rc=0 | controller, iteration 134: `for v in go1.26 go1.26.6 go1.99.99 go1.21rc1 go1.26.0 go1 go1.26.6+auto; do … GOTOOLCHAIN="$v" go version …; done`; split-shape corroborated first-party with network OFF: same loop under `GOPROXY=off` over `go1.26 go1 go1.99.99 go1.26.6 go1.26.6+auto 1.26.6`; `IsValid` re-checked in throwaway `/tmp/vercheck2` | controller: `go1.26`, `go1`, `go1.99.99`, `go1.21rc1` → rc=1 as `go: downloading <name> (darwin/arm64)` — the availability class, NOT `go: invalid GOTOOLCHAIN`; `go1.26.6`, `go1.26.0` → rc=0; **`go1.26.6+auto` → rc=0 and runs go1.26.6** — so banning mode words/`+` suffixes is pin-stability POLICY, not a validity claim. First-party corroboration (no network): `go1.26`/`go1`/`go1.99.99` → rc=1 `go: downloading <name>`; `go1.26.6`/`go1.26.6+auto` → rc=0 `go version go1.26.6`; same-call invalid-class control `1.26.6` → rc=1 `go: invalid GOTOOLCHAIN "1.26.6"`. `IsValid("go1")=true` (newly covered; controls `("go1.26")=true`, `("1.26.6")=false` match P9). Only non-name-shaped values (P3/P5's `1.26.6`, `min1.26.6`) ever produce the invalid-value error | INHERITED (controller, iteration 134); split shape corroborated VERIFIED-BY-ME under `GOPROXY=off` |
| P20 | This design leaves the `:314` canary fence's subject byte-untouched: the `GOTOOLCHAIN` token count in `host/store/toolchain_canary_test.go` is 0 today and no milestone edits that file | `grep -c 'GOTOOLCHAIN' host/store/toolchain_canary_test.go` with same-call control `grep -c 'func Test' host/store/toolchain_canary_test.go` | count = **0** (the fence wants exactly 0 — green today); control = **1**, proving the file was read; M1–M3 edit only `host/verifygate/toolchain_pin_gate_test.go`, so the fence and its subject keep their bytes | VERIFIED-BY-ME |
| P21 | Under the round-1 draft (validator arms as `t.Errorf`), ONE malformed value yields THREE messages, two misattributing — the E0 shape reproduced one call later | code read of the kept comparisons, first-party this pass: `sed -n '104,152p' host/verifygate/toolchain_pin_gate_test.go` | with `GOTOOLCHAIN: 1.26.6` at one site, `keyedValues` returns RAW values: `goToolchains=["1.26.6","go1.26.6"]`, `goVersions=["go1.26.6","go1.26.6"]`, `allPins[0]="1.26.6"` → `:120` fires `toolchain pins disagree: GOTOOLCHAIN=[1.26.6 go1.26.6] go-version=[…]` AND `:150` fires `go.mod floor="go1.26.6" disagrees with ci.yml toolchain pin="1.26.6"` — plus the validator's own message = 3 messages for one defect, 2 asserting a disagreement between things one of which is not a pin | VERIFIED-BY-ME code read (mechanism confirmed by controller, iteration 134); discharged by the shipped `Fatalf` form, absence asserted by M9 |
| P22 | The round-2 R1 hypothesis — that `version.IsValid` rejects custom toolchain names the runtime accepts — is REFUTED on the shapes the reviewer named, and the two agree on every measured shape; the policy reframing is applied anyway, because a finite sample cannot license a universal claim | `for v in go1.26.6-corp go1.26.6-devel go1.26.6_x devel; do printf '%-18s ' "$v"; GOTOOLCHAIN="$v" go version >/tmp/gt2.out 2>&1; printf 'rc=%s :: %s\n' "$?" "$(head -1 /tmp/gt2.out)"; done` with BOTH controls in the same session (`go1.26.6` known-good, `1.26.6` known-invalid), and `go/version.IsValid` over the identical four shapes | runtime: `go1.26.6-corp -> rc=1 :: go: downloading go1.26.6-corp (darwin/arm64)`; `go1.26.6-devel -> rc=1 :: go: downloading …`; `go1.26.6_x -> rc=1 :: go: invalid GOTOOLCHAIN "go1.26.6_x"`; `devel -> rc=1 :: go: invalid GOTOOLCHAIN "devel"`; controls `go1.26.6 -> rc=0 :: go version go1.26.6 darwin/arm64` and `1.26.6 -> rc=1 :: go: invalid GOTOOLCHAIN "1.26.6"`. `IsValid`: `("go1.26.6-corp")=true`, `("go1.26.6-devel")=true`, `("go1.26.6_x")=false`, `("devel")=false`, `("go1.26.6")=true`. **The `goV-suffix` custom form is accepted by BOTH** — it lands in the availability class, exactly like `go1.26`; the two shapes on which `IsValid` says false are the two the runtime calls `invalid GOTOOLCHAIN`. 4/4 agreement, plus P9/P19's eleven | VERIFIED-BY-ME (controller, iteration 134) — the reviewer's mechanism is refuted; its EPISTEMIC point stands and its fix is applied verbatim, because 15 agreeing shapes are still a sample and the shipped message must not assert what only a sample supports. `go doc go/version.IsValid` documents it as *"IsValid reports whether the version x is valid"* — a statement about Go VERSIONS, never about GOTOOLCHAIN, which is itself the reason not to quote it as a runtime-equivalence oracle |
| P23 | Round-2 R3's demanded re-derivation, run verbatim as the reviewer prescribed it | detached probe worktree at `07668e1`, `git status --porcelain` = 0 at entry; `shasum -a 256` of `run.sh` recorded BEFORE; `sed -i '' '140s/-eq 0/-ne 0/'`; LANDED proof = the guard-literal pair AND the site count; `go test ./host/verifygate/ -run TestMiscompileInstrumentProbesPinnedToolchain -count=1 -v`; `git checkout --` restore; sha256 re-compared; porcelain re-read | before `-eq 0`=**1** `-ne 0`=**0** sites=**3**; after `-eq 0`=**0** `-ne 0`=**1** sites=**3** → **MUTATION LANDED** (the guard literal moved, the site count did not — which is the defect itself); `=== RUN TestMiscompileInstrumentProbesPinnedToolchain` / `--- PASS` / `ok … 0.278s`, **mutant rc=0**; restore: **sha256 IDENTICAL** (`b80109aa…d0bb`), porcelain **0** | VERIFIED-BY-ME (controller, iteration 134) — this is the second independent run of measurement C in this session; the round-1 run and this one agree in every count |

## Options considered

**E0 — special-case the shared function: "don't prepend when the key is GOTOOLCHAIN."**
*Cost:* ~3 lines. *What it actually does:* inverts the bug without defining validity —
exactly what this item's brief forbids. A malformed `1.26.6` would then red, but as
`toolchain pins disagree: GOTOOLCHAIN=[1.26.6] go-version=[go1.26.6]` — a message asserting
the pins *disagree* when the true defect is that one of them *is not a value at all* (P3).
Misattribution is this mission's named failure mode (row 44 landed on it one week ago), and
the special case imports no grammar: the next convention added to `pinValues` inherits the
same coin-flip. Killed on arrival.
*Round-1 addendum (R3):* the first draft of Option A reproduced E0's shape one call
later — a validator whose arms were `t.Errorf` let execution continue into the very same
agreement/floor comparisons, producing three messages for one malformed value (P21). E0's
kill is honest only because the shipped Option A now actually differs from it on exactly
this axis: the validator STOPS (`t.Fatalf`), so a value that is not a pin never becomes
the operand of a comparison that would report a "disagreement" about it.

**A — split the normalizer per convention; validate GOTOOLCHAIN against Go's own grammar
(CHOSEN).** One canonicalizer for the conventions whose native form omits the prefix
(setup-go `go-version:`, go.mod `go` directives — P8), one strict validator for the
conventions that carry a toolchain NAME (`GOTOOLCHAIN` keyed lines, run.sh `PINNED=` — P7,
P16), with the name shape checked by `go/version.IsValid` (P9, P17) and the full setting
grammar cited from Go's documentation (P6). Two distinct FATAL red messages — one defect,
one message, the test stops rather than grading comparisons over a non-pin operand (P21;
R3): *invalid to Go itself* (P5's `rc=1 invalid` class) vs *valid to Go but a selection
mode, not a pin* (`auto`/`local`/`path`/`+auto`/`+path` — values under which the resolved
toolchain can move without the file changing; the ban is pin-stability policy, not a
validity claim — P19 measures `go1.26.6+auto` at rc=0). Fold in the two measured siblings:
Test B's floors go fatal (P10–P12), and the pinned-OK guard gets a direction pin (P13,
P14). *Cost:* ~0.15 d, one file. *Catches:* M1–M9 below. *Misses:* Declared residuals 1–6.

**B — parse ci.yml with a real YAML parser (or adopt actionlint) so pins are extracted
structurally.** *Cost:* a new dependency or tool for a repo whose gate philosophy is
line-scans with hand-maintained constants as the bar — Test A's own doc-comment declares
line-parsing and no-actionlint as accepted residuals (read at P10's file). Structural parsing
would not have caught THIS defect anyway: the value `1.26.6` parses fine as YAML; the
malformation is in Go's value grammar, which is exactly what Option A checks. Strictly larger
than the item, orthogonal to its cause. Deferred Scope.

**C — rely on runtime: an invalid GOTOOLCHAIN kills every `go` command in the job, so CI
reds anyway.** TRUE, and stated honestly (P3: rc=1 on the first `go` invocation). But the red
arrives as whatever `go` command runs first in the job — attribution lands on an innocent
step, not the pin line — and the *valid-but-not-a-pin* class (`GOTOOLCHAIN: auto`) has no
runtime red AT ALL: the job runs green today while the pin silently becomes a floating
selection mode whose resolved toolchain moves when go.mod moves. A pin gate that outsources
pin validity to a crash elsewhere reports `success` over its own subject — one merge after
row 44 taught precisely that lesson about a step reporting success over a refusal. Killed.

## Decision

**Option A.** Stop sharing one normalizer across two grammars; make each convention's
validator state that convention's own rule, sourced from the convention's own documentation.

Why this is not a package (S3): it is test-layer hardening of an existing host gate — no new
kernel surface, no new world/ or host/ production code; the S3 question does not arise beyond
this sentence.

All edits in `host/verifygate/toolchain_pin_gate_test.go`. Sketch (verbatim intent; the
executor may adjust identifiers, not semantics):

```go
func stripPinQuotes(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 && ((value[0] == '\'' && value[len(value)-1] == '\'') ||
		(value[0] == '"' && value[len(value)-1] == '"')) {
		value = value[1 : len(value)-1]
	}
	return value
}

// canonicalizeVersionPin serves ONLY the conventions whose native spelling omits the
// "go" prefix: actions/setup-go `go-version:` values and go.mod `go` directives
// (design doc P8). Prepending is correct here and nowhere else.
func canonicalizeVersionPin(value string) string {
	value = stripPinQuotes(value)
	if value != "" && !strings.HasPrefix(value, "go") {
		value = "go" + value
	}
	return value
}

// requireToolchainNamePin serves the conventions that carry a Go toolchain NAME:
// ci.yml GOTOOLCHAIN pins and run.sh's PINNED=. Per https://go.dev/doc/toolchain
// (where `go help environment` sends GOTOOLCHAIN readers), the setting's grammar is
// <name>, <name>+auto, <name>+path, or the shorthands auto/local/path — but only the
// bare <name> form PINS. The mode words and +suffix forms are VALID to Go — measured:
// GOTOOLCHAIN=go1.26.6+auto is rc=0 and runs go1.26.6 (design doc P19) — so rejecting
// them is a pin-stability POLICY choice, not a validity claim: they are selection
// channels under which the resolved toolchain can move without this file changing.
// The second arm is this REPOSITORY'S PIN POLICY, not a claim of equivalence with
// the runtime's own name grammar (round-2 R1; custom `goV-suffix` names measured in
// design doc P22). Both arms are Fatalf, and neither is an `instrument failure:` floor — the
// instrument is fine; the INPUT is not a pin. Stopping is still correct, because
// every comparison downstream (the agreement loop, the go.mod-floor check) takes
// this value as its operand: grading agreement over a non-pin is the E0
// misattribution reproduced one call later (design doc P21). One defect, one
// attributed message. Nothing is auto-corrected, because a value Go itself refuses
// (`go: invalid GOTOOLCHAIN "1.26.6"`, measured) must never be repaired into a
// value it would accept.
func requireToolchainNamePin(t *testing.T, source, key, raw string) string {
	t.Helper()
	value := stripPinQuotes(raw)
	switch {
	case value == "auto" || value == "local" || value == "path" || strings.Contains(value, "+"):
		t.Fatalf("%s: %s=%q is a toolchain-selection mode, not a pin; only a bare toolchain name (e.g. go1.26.6) pins", source, key, value)
	case !version.IsValid(value):
		t.Fatalf("%s: %s=%q is not an allowed standard Go toolchain pin; this repository requires a bare standard toolchain version accepted by go/version.IsValid (for example go1.26.6)", source, key, value)
	}
	return value
}
```

**What the shape check does and does not claim (round-1 R1; RESTATED round-2 R1, the
reviewer's own words).** `version.IsValid` enforces this repository's standard `goV` pin
policy; P3 separately verifies that the motivating value `1.26.6` is refused by the Go
runtime. This design does not claim `version.IsValid` recognizes every runtime-accepted
custom toolchain name. Within the measured sample the two agree (P22), but the check is a
POLICY, not a runtime-equivalence proof. It is NOT a discriminator for
availability: `go1.26`, `go1`, and `go1.99.99` all pass it, and all fail only as
`go: downloading <name>`, the availability class (P19). This design never leans on
`IsValid` alone. It is the first, *attributing* filter; behind it the design KEEPS the two
existing comparisons that pin the value to the exact go.mod floor string — the agreement
loop (`toolchain pins disagree`, `:120`) and the floor comparison
(`go.mod floor=%q disagrees…`, `:150`) — so a name-shaped non-floor value like `go1.26`
still reds, and reds correctly attributed: as a genuine disagreement between well-formed
pins, which is exactly what it is. What `IsValid` alone would admit in a repo WITHOUT a
floor to compare against is Declared residual 6.

Call-site disposition, one per P1 site:

- `pinValues` (:25) loses its normalize call and becomes `keyedValues` returning
  quote-stripped RAW values. In `TestGoToolchainPinsAgreeAndMatchJobList`, the
  `GOTOOLCHAIN` slice maps through `requireToolchainNamePin(t, "ci.yml", "GOTOOLCHAIN", raw)`
  — **per collected value**, so both of P7's sites (and any future third) are covered with
  no line numbers anywhere; the `go-version` slice maps through `canonicalizeVersionPin`.
  The keyed-count floors, agreement loop, setup-go count, floor comparison, and job
  enumeration keep their current semantics over the resulting slices (a validated
  GOTOOLCHAIN value participates in agreement exactly as before, since for well-formed input
  both old and new paths yield the identical string — base-green is preserved, AC7). A value
  that FAILS validation stops the test right there (`Fatalf`): the agreement loop, floor
  comparison, and job enumeration never grade it, so one malformed value gets ONE attributed
  message, not P21's three (M9 asserts the message set).
- `moduleGoFloor` (:45) → `canonicalizeVersionPin` (go.mod's native prefix-less grammar, P8).
- Test B `:244` → `requireToolchainNamePin(t, scriptPath, "PINNED", pinnedAssignments[0])`
  (P16: PINNED is compared byte-wise against toolchain names at run.sh:87; a prefix-less
  value must red statically, not be repaired into a string the script never uses).
- `normalizeToolchainPin` is **deleted**. The name with two meanings dies; its two
  successors each carry one grammar and say which (AC7 asserts zero remaining references
  with a same-call control).

**Finding B fold-in:** Test B's two floors at `:199` and `:209` change `t.Errorf` →
`t.Fatalf`, token-for-token otherwise — mirroring Test A's `:78`/`:128` and the Z3
precedent's 7-of-7 (P10, P11). An instrument that cannot find its own controls stops at the
cause instead of grading ~10 assertions with a broken ruler (P12).

**Finding C fold-in:** after Test B's existing `saw_pinned_ok` site-count check, add a
direction pin over comment-stripped CODE lines (the house pattern already shipped in
`TestMiscompileInstrumentStepIsGatedInCI`'s executable-use scan):

```go
	// Direction pin (row 45, finding C): the site COUNT above cannot see the guard's
	// comparison direction — flipping `-eq 0` to `-ne 0` preserves all 3 sites and
	// every committed assertion (controller, iteration 134). Count comment-stripped
	// CODE lines only: the guard and the prose explaining it must not compete for one
	// namespace (iteration-133 lesson) — a comment may quote the literal freely.
	guardLines := 0
	for _, line := range lines {
		code := line
		if idx := strings.Index(code, "#"); idx >= 0 {
			code = code[:idx]
		}
		if strings.Contains(code, `[ "$saw_pinned_ok" -eq 0 ]`) {
			guardLines++
		}
	}
	if guardLines != 1 {
		t.Errorf("%s: executable pinned-OK guard-line count=%d, want exactly 1 — the floor must test ABSENCE of the OK flag (`-eq 0`); a flipped or duplicated guard is a different instrument", scriptPath, guardLines)
	}
```

The `!= 1` form is two-sided: 0 (flipped/rewritten/deleted guard) and ≥2 (duplicated guard)
both red, so the pin cannot pass vacuously on its own null case (S6). Honesty note carried
from P14: since row 44 gated the step, a flipped guard is ALSO loud at runtime on the next
push; this pin's value is pre-push detection and correct attribution, and it is the only
COMMITTED assertion that sees the direction.

Coupled-code disposition: `run.sh` and `ci.yml` ship byte-untouched (mutation venues only).
`TestReproModuleFloorStaysBelowKnownBadToolchains` — behavior unchanged (its
`moduleGoFloor` inputs are `go `-directive values, P8's lax convention, and its KNOWN_BAD
validation already uses `version.IsValid` directly). Row 50's `shellAssignmentValues` —
byte-untouched. `host/store/toolchain_canary_test.go` — byte-untouched, so
`TestCanaryDeclaresPositiveArmOnly`'s `:314` token fence keeps its subject at
`GOTOOLCHAIN` count 0 (P20). `scripts/verify_go.sh:214–222` — untouched and
non-conflicting: its deny-list observes `go env GOVERSION`, the ACTIVE toolchain at
runtime, never the pin text (P18). Row 44's wiring test at the file end — byte-untouched; this item edits
`:13–:34`, `:104`-region, `:199`, `:209`, `:244`, and inserts inside Test B — disjoint from
the file-end hunks, as row 44's own conflict note anticipated.

## Milestones

**M1 — split the normalizer; validate names strictly (~0.08 d).**
`stripPinQuotes` + `canonicalizeVersionPin` + `requireToolchainNamePin` land;
`pinValues`→`keyedValues`; the three call-site dispositions above; `normalizeToolchainPin`
deleted. Whole-package green on the unmutated tree before commit (AC7); mutations M1–M5 and the
M9 attribution pair rehearsed per the house recipe (landed-proof → observe → restore by
sha256 → porcelain 0).

**M2 — Test B's instrument floors go fatal (~0.03 d).**
Two tokens: `:199` and `:209` `t.Errorf` → `t.Fatalf`. Mutations M7–M8 rehearsed.

**M3 — pinned-OK guard direction pin (~0.04 d).**
The `guardLines` block lands in Test B. Mutation M6 rehearsed, plus its comment-quote green
control (AC6).

Each milestone is independently committable and leaves the whole gate green; order is free
except M3 reads Test B's `lines`, which both M1 and M2 leave untouched.

## Acceptance criteria

Counts below are 0/1 predicates or extractor-derived cardinalities, never transcribed line
numbers; every zero carries a same-call known-positive control. "Probe worktree" = a
detached worktree at the landing commit, mutations proven LANDED before results are read,
restored to porcelain 0 (house recipe).

- **AC1 — the validators exist and RUN (run-existence form).**
  `go test ./host/verifygate/ -run 'TestGoToolchainPinsAgreeAndMatchJobList|TestMiscompileInstrumentProbesPinnedToolchain' -count=1 -v`
  → rc=0 with one `=== RUN` and one `--- PASS` per test; paired nonsense control
  `-run '^TestNoSuchPinGateZZZ$'` prints `[no tests to run]`. **Base:** both PASS today
  (P15) — AC1's teeth are M1–M9, not this green.
- **AC2 — every collected GOTOOLCHAIN site rejects the prefix-less form (property over the
  extractor's set, not line numbers).** Derive the site set at execution time:
  `grep -n 'GOTOOLCHAIN:' .github/workflows/ci.yml` must yield exactly as many lines as the
  enumerated-job cardinality (same-call control: non-zero, currently the two P7 sites). For
  EACH site in turn: strip the value's `go` prefix, assert LANDED, run Test A scoped `-v` →
  rc≠0 with the not-an-allowed-pin message carrying the malformed value as the ONLY failure
  message (`Fatalf`: the agreement loop and floor comparison never run for that arm —
  M9(a) asserts the absences); restore between arms. **Base:** the first-site arm is measured GREEN today (P4) — the
  headline flip this item exists to reverse.
- **AC3 — a valid-but-not-a-pin value reds with the NOT-A-PIN class as the ONLY
  message.** Probe worktree: one collected site's value → `auto`; Test A scoped `-v` →
  rc≠0; the output contains the selection-mode message EXACTLY once, ZERO occurrences of
  `toolchain pins disagree`, and ZERO occurrences of `disagrees with ci.yml toolchain
  pin` — exact counts, no hedge (round-1 R3), with M9(b)'s same-session control proving
  both absent strings ARE observable when a genuine disagreement fires. **Base
  (code-derived from P2, not measured):** today `auto` normalizes to `goauto` and reds
  only as `pins disagree` — red for the wrong reason; after M1 one defect gets one
  attributed message.
- **AC4 — the zero-extraction floor stays reachable across the extraction split.** Probe
  worktree: every `GOTOOLCHAIN:` key → `GO_TOOLCHAIN:` (LANDED: `grep -c 'GOTOOLCHAIN:'` →
  0 with same-call control `grep -c 'GO_TOOLCHAIN:'` = the old cardinality); Test A scoped →
  rc≠0 via the keyed-line-count floor (`count=0, want` job cardinality). This is the input
  shape that makes the validator loop vacuous — the floor, not the loop, must red it (S6).
- **AC5 — instrument floors are FATAL: one cause, one message.** Probe worktree: rename
  run.sh's `PINNED="` assignment prefix so the control token vanishes (LANDED:
  `grep -c 'PINNED=' run.sh` → 0, same-call control `grep -c 'KNOWN_BAD=' run.sh` ≥ 1);
  Test B scoped `-v` → rc≠0 whose output contains the `does not contain known-positive
  control` floor EXACTLY once and ZERO of Test B's downstream assertion messages
  (assignment-count, PINNED-equality, exec-bit, guard-message). **Base:** controller,
  iteration 134 (P12): today the Errorf floor limps into the downstream cascade.
- **AC6 — the guard-direction pin sees the flip, and comments stay free.** Probe worktree,
  two arms: (a) run.sh's single executable `-eq 0` guard line → `-ne 0` (LANDED per P13's
  count pair inverting: `-eq` form 1→0, `-ne` form 0→1); Test B scoped → rc≠0 naming the
  guard-line count. (b) Green control for the iteration-133 namespace lesson: restore, then
  append a COMMENT line quoting `[ "$saw_pinned_ok" -eq 0 ]` verbatim; Test B scoped →
  rc=0. Restore. **Base:** arm (a) is measured rc=0 today (P14).
- **AC7 — hygiene, base-green preservation, and the name is dead.**
  `gofmt -l host/verifygate/` → empty; `go vet ./host/verifygate/` → rc=0;
  `go test ./host/verifygate/ -count=1` → rc=0 on the unmutated tree (base-green: the
  well-formed P7/P8 values pass both new paths byte-for-byte as before);
  `grep -c 'normalizeToolchainPin' host/verifygate/toolchain_pin_gate_test.go` → 0 with
  same-call control `grep -c 'canonicalizeVersionPin'` ≥ 2 (definition + a call site)
  proving the grep is live; `git status --porcelain` → 0 lines after all rehearsals.

## Named RED mutations

**Needle discipline (round-2 carve-out).** Where a row below says *not-an-allowed-pin
message* or *selection-mode message*, the executor must grep the **shipped literal** from the
Decision's sketch — `is not an allowed standard Go toolchain pin` and `is a toolchain-selection
mode, not a pin` respectively — never this prose. Round 2 rewrote the first message's text
(R1's verbatim fix), and a mutation arm keyed on a description rather than on the shipped bytes
is how a correct landing gets redded.

Venue: `scripts/verify_go.sh` runs `go test ./... -count=1` inside the gated `go-verify`
job (row 44's P15, re-usable unchanged), so every arm below is a CI-gate red. House recipe
per arm: prove LANDED by count/content before reading the result; restore byte-identical by
sha256; porcelain 0 after. Cross-checked against the Decision: every arm targets a construct
the design SHIPS (no arm references `normalizeToolchainPin`, which the design deletes).

| # | Mutation | File | What it neuters | Predicted result | Landed-proof before reading the result |
|---|---|---|---|---|---|
| M1 | First collected `GOTOOLCHAIN` site (today `:21`): `go1.26.6` → `1.26.6` | `ci.yml` | the headline validation — the arm measured GREEN today (P4) | Test A RED via `Fatalf`: the not-an-allowed-pin message carrying `"1.26.6"` is the ONLY failure message (M9(a) asserts the absences); the Go-runtime control for the same value is P3's rc=1 | `go1.26.6` count on GOTOOLCHAIN-keyed lines 2→1; bare `1.26.6` present on a `GOTOOLCHAIN:` line |
| M2 | SAME flip at the other collected site (today `:102`) | `ci.yml` | per-site coverage — proves the design validates the extractor's SET, not the row's one line number | Test A RED, same message class as its only message, naming the second site's value | as M1, on the other keyed line |
| M3 | One collected site: `go1.26.6` → `auto` | `ci.yml` | the pin/mode distinction — a value Go ACCEPTS (P5 rc=0) that is not a pin | Test A RED via the selection-mode arm as the ONLY message (AC3's exact counts) | `grep -c 'GOTOOLCHAIN: auto'` 0→1 |
| M4 | Every `GOTOOLCHAIN:` key → `GO_TOOLCHAIN:` | `ci.yml` | the extraction itself — the validator loop's vacuous-input shape (S6 floor arm) | Test A RED via the keyed-count floor: count=0 vs enumerated-job cardinality | `grep -c 'GOTOOLCHAIN:'` →0 with `grep -c 'GO_TOOLCHAIN:'` = old cardinality in the same call |
| M5 | `PINNED="go1.26.6"` → `PINNED="1.26.6"` | `run.sh` | the PINNED-site strict validation | Test B RED: not-an-allowed-pin message for key `PINNED`. TODAY this mutant is statically GREEN by auto-correction (code-derived, P2+P16) while run.sh:87's byte compare could never set `saw_pinned_ok` — the static/runtime mispolarity this arm closes | `grep -c 'PINNED="1.26.6"'` 0→1 |
| M6 | run.sh:140 `-eq 0` → `-ne 0` | `run.sh` | the guard's comparison direction — measured invisible to all committed assertions today (P14) | Test B RED via the direction pin: `guard-line count=0, want exactly 1` | P13's count pair inverts: `-eq` form 1→0, `-ne` form 0→1 |
| M7 | `PINNED="` → `PINHOLE="` (assignment prefix renamed; control token gone) | `run.sh` | fatality of the control floor (finding B) | Test B RED via the `:199`-class floor, now `Fatalf`: `-v` output carries the floor message EXACTLY once and ZERO downstream Test B assertion messages — the sharpened prediction that distinguishes the shipped `Fatalf` from today's measured cascade (P12) | `grep -c 'PINNED=' run.sh` →0 with same-call control `grep -c 'KNOWN_BAD=' run.sh` ≥1 |
| M8 | Append a second `#!/usr/bin/env bash` line at end of file | `run.sh` | fatality of the shebang floor (`:209`-class); its input shape is trivially producible — this is how | Test B RED via the shebang floor, now `Fatalf`: `exact shebang count=2, want 1`, single failure message, zero downstream | exact-match shebang line count 1→2 |
| M9 | Attribution pair, same probe session (round-1 R3): **(a)** M1's flip — `go1.26.6` → `1.26.6` at the first collected site; **(b)** restore, then `go1.26.6` → `go1.25.6` at the same site — a VALID name that genuinely disagrees with the other three pins and the go.mod floor | `ci.yml` | single-attribution itself: an absence assertion is vacuous unless the absent strings are observable — arm (b) is the same-session known-positive control | **(a)** Test A RED, `-v` contains the not-an-allowed-pin message EXACTLY once, `toolchain pins disagree` count = **0**, `disagrees with ci.yml toolchain pin` count = **0**; **(b)** Test A RED, `-v` CONTAINS `toolchain pins disagree` AND the floor-disagree message — both strings CAN fire (the validator passes `go1.25.6`), so arm (a)'s zeros are measured absences, not dead greps | (a) as M1; (b) `go1.25.6` count on the keyed line 0→1, `go1.26.6` count on GOTOOLCHAIN-keyed lines 2→1 |

Green control for all arms: the unmutated post-sprint tree passes AC1–AC7, and AC6(b) is the
named must-stay-green arm (a comment quoting the guard literal). Every floor this design
adds or hardens has an arm: keyed-count/vacuous-extraction → M4; control floor → M7; shebang
floor → M8; single-attribution of the `Fatalf` validator → M9, whose (b) arm doubles as the
known-positive control for AC3's and M9(a)'s absence greps; guard pin's own two-sided form
→ M6 (zero side) with the ≥2 side exercised by M8's recipe applied to the guard line if the
evaluator wants a tenth arm (not required — the `!= 1` predicate is one construct).

## Deferred Scope

- **Structural YAML extraction / actionlint adoption (Option B)** — a repo-wide gate
  philosophy change, not a row-45 repair; the line-scan residual class is already declared
  in Test A's doc-comment and inherited unchanged here. One line, one reason: the
  malformation this row measures is in Go's value grammar, which line-scans see fine.
- **Availability probing of a valid-shaped pin** — `go1.99.99` passes `version.IsValid`
  (P9) but downloads nothing (row 44's M6 measured the SKIP shape). Runtime layers (setup-go
  install, the gated run.sh coverage floor) own availability; a static availability check
  means network in a unit test. Declared residual 2 names the gap.
- **Strict validation for `verify_go.sh`'s GOTOOLCHAIN surface** — `:214–222` is a RUNTIME
  deny-list over `go env GOVERSION` (the ACTIVE toolchain; `# observe; never assign`), and
  `:222` is message text that merely SUGGESTS `GOTOOLCHAIN=go1.25.6`; it validates no pin
  text, so it neither conflicts with nor substitutes for this gate (P18, amended round 1).
  Teaching the gate to read shell `export` pins generally is row-50-adjacent surface
  (assignment parsing) and stays out.
- **Generalizing the direction-pin pattern to run.sh's other floors** — the same
  site-count-vs-direction blindness plausibly holds for `saw_bad`/`saw_good`/coverage
  guards; measuring which flips survive is a fresh controller measurement, not an assumption
  this doc may transcribe. One row, if the measurement lands.

## Declared residuals

1. **Still a static text scan.** Everything Test A's doc-comment declares survives
   unchanged: a step-level `if:` that never runs, YAML anchors/flow-style hiding keyed
   lines (a keyed-count RED, not a silent pass), values computed at workflow runtime, and a
   `GOTOOLCHAIN` exported in a script rather than keyed in YAML are all invisible. This
   item narrows the VALUE grammar, not the extraction physics.
2. **Validity ≠ availability.** `GOTOOLCHAIN: go1.99.99` passes the shape gate (P9) and
   dies only at runtime/download. The gate now refuses to bless malformed bytes; it still
   cannot promise the well-formed bytes resolve.
3. **The direction pin is a byte pin on one literal.** A semantically equivalent rewrite of
   the guard (`[ 0 -eq "$saw_pinned_ok" ]`, `test`, arithmetic `(( ))`) reds it FALSELY — a
   loud false red demanding a deliberate pin update, the accepted failure polarity — and a
   logically inverted guard rebuilt in different bytes elsewhere is unseen. Runtime loudness
   (P14's note: the gated step reds on the next push) backstops the unseen case.
4. **The mode-word arm bans a capability, not just a defect.** If this repo ever
   legitimately wants `GOTOOLCHAIN: local+path` in CI, the pin gate will red it and the
   change needs its own row — by design: under a pin regime, "the resolved toolchain can
   move without a diff" IS the defect, but the residual is that the gate cannot tell a
   future ratified policy change from a mistake.
5. **Fatality trades breadth for attribution.** After M2, a run with BOTH a missing control
   and (say) a flipped guard reports only the first; the second surfaces on the re-run
   after the first repair. The same now holds for the name validator's `Fatalf` arms
   (round-1 R3): a workflow with two malformed pins reports the first collected one per
   run. The Z3 precedent (P11) and Test A (P10) already accepted this trade for
   instrument-class floors; this doc extends it to those floors AND to the validator's
   input-defect arms, and says so.
6. **The shape check's teeth against name-shaped drift are floor-dependent.**
   `version.IsValid` admits every well-formed toolchain name — `go1.26`, `go1`,
   `go1.99.99` (P9, P19). HERE those red anyway, correctly attributed as genuine
   disagreements, because the kept agreement loop and floor comparison pin the value to
   the exact go.mod floor string. In a repo (or a future convention) without a floor to
   compare against, `IsValid` alone would bless any name-shaped value, and only
   runtime/download would object — in the availability class, with attribution on an
   innocent step. Porting `requireToolchainNamePin` without also porting a floor
   comparison ships a syntax gate, not a pin gate.

7. **Custom `goV-suffix` toolchain names are intentionally unsupported as pins here.**
   `go1.26.6-corp` is accepted by both the runtime (it proceeds to download) and by
   `version.IsValid` (P22), so nothing in this gate would reject it on shape — it is the
   go.mod-floor comparison that refuses it, as a disagreement. That is the correct
   attribution for this repository, which pins one exact released toolchain; a repository
   that genuinely ships a custom toolchain would need the floor and the pin to name it in
   the same edit, and this design does not attempt to serve that case. Declared, per
   round-2 R1's "either test/document custom-name behavior or declare custom toolchain
   names intentionally unsupported" — measured rather than assumed.

## Quorum round 1

Full-strength review: 3/3 present, `absent_reviewers` EMPTY. Verdicts: R1 REJECT, R2
REJECT, R3 REJECT — blocked. This is the protocol-mandated single revision; exactly one
re-quorum follows, and this section is what it should re-read the revision against. Every
count stated in the applied fixes was re-derived first-party this pass, not transcribed.

### R1 — gpt5-6-sol — REJECT — **half conceded, half corrected by measurement; both halves applied**

*Objection (verbatim core):* "`go/version.IsValid` validates Go version syntax, not the
stricter executable-toolchain-name grammar. The document's own P9 shows it accepts
`go1.26`, while P5 never tests whether `GOTOOLCHAIN=go1.26` is accepted by the Go
command. Therefore the load-bearing claim that `version.IsValid` is the discriminator for
`<name>` is unsupported and likely false; the gate can still misclassify a value the
runtime rejects as valid."

- **Conceded plainly:** P5 never tested a name-shaped-but-unusable value, and `IsValid`
  accepts `go1.26`, `go1`, and `go1.99.99` (P9; `IsValid("go1")=true` newly measured this
  pass, P19). The coverage gap was real.
- **Corrected by measurement (controller, iteration 134; split shape corroborated
  first-party under `GOPROXY=off`, P19):** the runtime does NOT reject those values in
  the invalid-value class. Every name-shaped value fails as `go: downloading <name>` —
  the AVAILABILITY class — and only non-name-shaped values (`1.26.6`, `min1.26.6`)
  produce `go: invalid GOTOOLCHAIN`. So `IsValid` IS an exact discriminator for the class
  this row is about (syntax) and is NOT one for availability. The Decision now claims
  exactly the first and disclaims the second ("What the shape check does and does not
  claim").
- **The real gap the objection exposes, applied:** the doc's prose leaned on `IsValid`
  alone while the design already keeps two stronger comparisons — the agreement loop
  (`:120`) and the floor comparison (`:150`) — that pin the value to the exact go.mod
  floor string and catch `go1.26` regardless, correctly attributed as a genuine
  disagreement. The Decision now states this layering explicitly, and new Declared
  residual 6 names what `IsValid` alone would admit in a repo without a floor.
- **Mode-word note applied:** `GOTOOLCHAIN=go1.26.6+auto` is rc=0 and RUNS go1.26.6
  (P19), so the `auto`/`local`/`path`/`+` ban is a pin-stability POLICY choice, not a
  validity claim — the validator's comment now says so in those words, and no sentence in
  this doc calls those values invalid.

### R2 — gemini-3-1-pro — REJECT — **applied verbatim: search space widened to the whole repo, per-hit dispositions printed**

*Objection (verbatim core):* "Premise 18 claims the absence of a conflicting static
GOTOOLCHAIN validator is 'measured, not assumed', but the provided command restricts the
search entirely to shell scripts. Since the existing static gate machinery is written in
Go, excluding `.go` files from the search space renders the premise unverified and the
conflict-surface analysis incomplete."

- **P18 amended in place:** `grep -rn 'GOTOOLCHAIN' --include='*.go' .` → exactly 8 hits,
  all in `toolchain_pin_gate_test.go`, each dispositioned (doc-comment prose, enumeration
  message + calls, extraction, keyed-count floor, agreement loop, and the `:314`/`:315`
  canary token fence — a fence over a DIFFERENT file, not a value validator); same-call
  known-positive control `normalizeToolchainPin` → 4 hits; plus a whole-repo census (code
  surfaces beyond `.go`: `ci.yml` ×2 = the subject, `verify_go.sh` ×1, `run.sh` ×1;
  everything else design-doc/sprint-plan prose).
- **Conflict surface completed:** `verify_go.sh:214–222`'s runtime deny-list over
  `go env GOVERSION` now appears in P18, the Coupled-code disposition, and the Deferred
  Scope bullet — it observes the ACTIVE toolchain, never the pin text.
- **Conclusion restated at its true strength, in the reviewer-supplied words:** no
  competing static validator of the pin TEXT exists; `:314` is a token fence over a
  different file and `verify_go.sh` observes the ACTIVE toolchain, so neither conflicts.
- **Self-check against `:314` as instructed:** new P20 — the `GOTOOLCHAIN` token count in
  `host/store/toolchain_canary_test.go` is **0** (the fence wants exactly 0; same-call
  control `func Test` = 1 proves the file was read), and every milestone edits only
  `toolchain_pin_gate_test.go`, so the fence's subject stays byte-untouched.

### R3 — oc-glm-5-2 — REJECT — **applied in the reviewer's terms: the validator STOPS, the message set is ONE, M9 asserts the absences**

*Objection (verbatim core):* "The chosen `requireToolchainNamePin` uses `t.Errorf`,
allowing execution to continue into the unchanged agreement loop, which then emits a
'pins disagree' message for the same malformed GOTOOLCHAIN value — the exact
misattribution the doc identifies as a killable offense in Option E0. AC3's hedge ('does
NOT attribute the defect solely to pin disagreement') silently concedes both messages
appear but never reconciles this with E0's reasoning… the new validator it introduces
reproduces the same Errorf-into-cascade shape for value-validation failures."

- **The reviewer's frame adopted as the rule:** a value that is not a pin makes every
  downstream comparison meaningless, so the validator must STOP rather than grade. Both
  validator arms are now `t.Fatalf` (Decision sketch), with the classification the doc
  owed: this is NOT an `instrument failure:` floor — the instrument is fine, the INPUT is
  not a pin — but stopping is still correct because the agreement loop and floor
  comparison take the value as their operand (the sketch's comment and P21 say exactly
  this).
- **The three-message mechanism is now a premise, not a hedge:** P21 derives it
  first-party from the kept code (raw slices → `:120` disagree + `:150` floor mismatch +
  the validator's own message = 3 messages for one defect), confirming the controller's
  iteration-134 mechanism read.
- **Message set named; absence asserted by a named RED mutation:** one malformed value →
  ONE attributed message (AC2, AC3, M1–M3 all updated to exact single-message
  predictions). New M9: arm (a) asserts `toolchain pins disagree` count 0 and
  floor-mismatch count 0 under the malformed value; arm (b), same probe session, is the
  known-positive control — a valid-but-genuinely-disagreeing `go1.25.6` makes both
  strings fire, proving the zeros are measured absences, not dead greps.
- **Consistency repairs:** AC3's hedge deleted (exact counts, no "solely"); E0's kill
  reasoning carries a round-1 addendum acknowledging the draft reproduced E0's shape one
  call later and stating why the shipped `Fatalf` form does not; Declared residual 5
  extended to name the breadth-for-attribution trade the validator's arms now share.

## Quorum round 2 — and the carve-out that closed it

Round 2 ran at **full strength**: `absent_reviewers` **EMPTY**, 3 of 3 present, caps
pre-raised to `$0.20`/reviewer so the self-selecting budget trap could not fire.
Verdict **blocked**, 2 reject / 1 pass — `gemini-3-1-pro` flipped to **PASS**, its round-2
objection non-blocking (the `#`-stripping heuristic is "acceptable … True parsing of shell
syntax properly belongs in Option B (Deferred Scope)" — its own words, and already where
this doc put it).

Closed under the **narrow-refinement carve-out**, ratified for this mission at iteration 44,
so the first-use gate does not apply. Both conditions were checked and both hold: every
remaining blocking objection carries a concrete **reviewer-authored `proposed_fix`**, and
neither disputes the design DIRECTION — R1 is an over-claim in message TEXT and prose, R3 is
a verification-completeness ask. Each fix is applied in the reviewer's own words; nothing was
overridden and nothing was argued past.

### R1 — gpt5-6-sol — REJECT — **applied verbatim; its mechanism separately measured and REFUTED, which does not weaken the fix**

*Objection:* "The load-bearing claim that `go/version.IsValid` is an exact discriminator for
values the Go runtime accepts as toolchain names is unverified. P9 and P19 test only a finite
sample, while the Decision universally claims that every rejected value would produce
`go: invalid GOTOOLCHAIN`. Go toolchain names may include non-standard/custom naming forms
beyond the standard `goV` syntax recognized by `version.IsValid`, so the validator can reject
a runtime-accepted name while falsely reporting that the runtime itself refuses it."

*Applied, verbatim from `proposed_fix`:* the second arm's message is now
`%s: %s=%q is not an allowed standard Go toolchain pin; this repository requires a bare
standard toolchain version accepted by go/version.IsValid (for example go1.26.6)`, and the
universal prose is replaced by the reviewer's sentence — *"`version.IsValid` enforces this
repository's standard `goV` pin policy; P3 separately verifies that the motivating value
`1.26.6` is refused by the Go runtime. This design does not claim `version.IsValid`
recognizes every runtime-accepted custom toolchain name."* The premise row the fix asks for
is **P22**, and residual **7** takes the fix's second branch ("declare custom toolchain names
intentionally unsupported").

*And the mechanism it names was measured before anything was applied, per ghost discipline —
it is **REFUTED** on the shapes the reviewer itself named.* `go1.26.6-corp` and
`go1.26.6-devel` are accepted by BOTH the runtime (download / availability class) and
`version.IsValid` (**true** for both); the two shapes where `IsValid` says false —
`go1.26.6_x`, `devel` — are exactly the two the runtime calls `invalid GOTOOLCHAIN`. Four of
four agree, on top of P9/P19's eleven. So no such divergence exists on this rig today. **The
fix is applied anyway, and this is the honest reason:** fifteen agreeing shapes are still a
sample, and the shipped `Fatalf` text was asserting something about *the Go runtime* that only
*this repository's policy* can support. The reviewer's epistemic point survives the refutation
of its example, and it is the epistemic point the message text was violating.

### R2 — gemini-3-1-pro — PASS — non-blocking; catch already in Deferred Scope

*Catch:* "Test B's `#` stripping is lexically naive and vulnerable to quoted `#` characters
earlier on the same line." *Its own fix:* "the naive strip is acceptable as a heuristic. True
parsing of shell syntax properly belongs in Option B (Deferred Scope)." Recorded here so the
next reader does not re-derive it; no change made, because the reviewer asked for none.

### R3 — oc-glm-5-2 — REJECT — **applied verbatim: P14 re-derived first-party, second independent run**

*Objection:* "P14's headline mutation claim — that flipping `run.sh:140`'s `-eq 0` to `-ne 0`
leaves all 18 test arms green today — is marked INHERITED … yet it is the load-bearing premise
for Finding C (the guard-direction pin), M6, and AC6. … An unverified load-bearing mutation
premise is grounds to reject."

*Applied, verbatim from `proposed_fix`:* the mutation was re-run in the detached probe worktree
at `07668e1` — LANDED proved via the count pair, `go test ./host/verifygate/ -run
TestMiscompileInstrumentProbesPinnedToolchain -count=1 -v` observed, restore verified by
sha256, porcelain asserted 0 — and recorded as **P23**, with P14's status changed from
`INHERITED` to a re-derivation pointer. **One honest amendment to the attribution the reviewer
assumed:** the reviewer asked the *designer* to re-run it; the re-run is the **controller's**,
performed in this iteration's own session. Both runs agree in every count (guard literal
1→0, flipped form 0→1, site count unmoved at 3, mutant `rc=0 --- PASS`, sha256 identical on
restore). The row says whose hands ran it rather than borrowing the designer's label — the
premise is first-party to this iteration, which is what the objection was actually about.


## Related Documents

- `design_docs/implemented/w-setup-go-pin-unguarded.md` — Test A's origin (the pin
  agreement gate whose normalizer this item splits).
- `design_docs/implemented/w-miscompile-instrument-inert-in-ci.md` — row 44, landed one
  merge ago (`46add2c`): same file, disjoint hunks (its own conflict note anticipated this
  item); source of the gated-step fact behind P14's honesty note and the house
  mutation/landed-proof recipe reused verbatim.
- `design_docs/implemented/w-race-gate-blindspot.md` — the instrument whose `PINNED=` and
  pinned-OK guard this item's Test B changes bind more tightly.
- `design_docs/coding-standards.md` — S6 (a gate must fail loudly on its own null case:
  AC4, AC5, the two-sided guard pin) and the S3 non-package answer in Decision.
- `design_docs/world-mission.md` — queue row 45 (the item, verbatim intent).
- https://go.dev/doc/toolchain — Go's own GOTOOLCHAIN grammar (P6), the external source the
  validator's comment cites.
