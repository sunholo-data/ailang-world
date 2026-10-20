# Evaluator r2 — iteration 287, row 156 M5a, commit 3852dbc2

VERDICT: PASS — score 91/100
Breakdown: correctness/design fidelity 29/30; test quality and non-vacuity 24/25; scope/standards 20/20; gates 10/10; mutation evidence honesty 8/10.

Delta from r1 (6350cb35): `git diff 6350cb35 3852dbc2` is test-only, +10 lines in TestEpisodePolicyProtectedGolden/escapes (requires `\u007F`, refuses raw rune < 0x20 or == 0x7f after stripping "\n").

## Gates (WORLD_TOOL_AILANG_BIN unset, pinned AILANG_BIN)
- go vet ./... rc 0.
- broker 6 named tests: 6/6 PASS, rc 0 (unmutated tree, so the new assertions pass on a correct renderer, not vacuously failing).
- daemon 3 named tests: PASS (ok).
- -race suites not re-run (UNMEASURED by me).

## Mutants (tomlString, one at a time, cp restore, git diff empty after each)
| Mutant | Result |
|---|---|
| `case r < 0x20:` | KILLED, Golden/escapes :751 "DEL is not escaped" |
| `case r < 0x7f && unicode.IsControl(r):` (DEL raw) | KILLED, Golden/escapes (DEL assertion) |
| `case r < 0x20 \|\| C1 range:` (DEL raw) | KILLED, :751 |
| `case r > 0 && IsControl(r):` (NUL raw) | KILLED, :744 invalid character in string literal |
| `case r < 0x7f:` literally (as the controller phrased it) | KILLED, but by TestRenderRunVariantKeySetsPerClass :111 (it escapes all printable ASCII), not by /escapes; a mutated-away-DEL variant is the meaningful one and is killed above |
| lowercase `\u%04x` | KILLED :751 (the test now pins uppercase hex; semantically equivalent in TOML, harmless over-pinning) |
Note: the plain "no-op" control mutant was not applied (identical text), so the unmutated-tree PASS above is the vacuity control.

## Carry-forward
- N1: CLOSED. DEL and C0 raw-control mutants now go RED on assertions.
- N2: still stands as a documentation caveat only (JSON decoder is a sound proxy for the emitted subset; the new raw-rune scan now covers its main blind spot, raw DEL/C0). Non-blocking.
- N3: NON-BLOCKING, unreachable through a profile. `ParseProtectedPatterns` (host/broker/exec_protected.go:75) refuses `!utf8.ValidString` entries, and also refuses any control rune (:89) and backslash, so a parsed pattern is always valid UTF-8 before rebasing/rendering; NFC normalization preserves validity. Reachable only by a direct caller of RenderEpisodePolicy bypassing the parser. Carry to M5b as the controller proposes (M5b must call the parser before render). UNMEASURED: M5b wiring itself.
- N4: still stands as informational; chain sound.

No blocking findings. Eval worktree clean; scratch removed.
