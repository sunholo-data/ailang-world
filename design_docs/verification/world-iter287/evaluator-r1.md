# Evaluator r1 — iteration 287, row 156 M5a (commit 6350cb35 on 0adeb703)

VERDICT: PASS  — score 88/100

Breakdown: correctness/design fidelity 29/30; test quality and non-vacuity 22/25 (-3: escape guard has surviving mutants); scope/standards discipline 20/20; gates 10/10; mutation evidence honesty 7/10 (executor table verified accurate; one "independent decoder" claim overstated).
(My breakdown is approximate against the sprint-evaluator categories.)

## Gates (host, WORLD_TOOL_AILANG_BIN unset, AILANG_BIN=pinned)
- go vet ./... : rc 0. gofmt -l host/ : clean.
- broker focused (6 named tests, -v): rc 0, 6/6 PASS (Golden, AllVariants, NoProtectedIsByteIdentical, ThroughHandler, DesignKeySet, KeySetsPerClass).
- daemon focused (3 named tests): rc 0, 3/3 PASS.
- Race suites not re-run (controller measured green); UNMEASURED by me.
- All production callers of RenderEpisodePolicy/RenderRunVariant updated (daemon passes nil, by design). No World symlink/hard-link pre-check present (correct for M5a).

## Executor mutants re-verified (one mutant at a time, restored with cp, git diff empty after each)
| Mutant | Result | Failing assertion |
|---|---|---|
| NO-POLICY-APPEND (drop append of spec.protected) | KILLED by Golden | handlers_ailang_run_test.go:717 "protected policy differs from independent golden" |
| VARIANT-MISSED all (drop spec.protected in RenderRunVariant) | KILLED by AllVariants/{Declassify,Net,Env} | :765 |
| VARIANT-MISSED Env only (skip when mode=trusted_host) | KILLED, AllVariants/Env only | :765 |
| VARIANT-MISSED Net only (skip when netAllow set) | KILLED, AllVariants/Net only | :765 |
| POLICY-QUOTE (remove `"` escape case) | KILLED by Golden/quote | :726 and :728 |
| POLICY-EMPTY-BYTES (append ".zz/**" when no protection) | KILLED by NoProtectedIsByteIdentical (base,Declassify,Net,Env + episode) | :797, :811 |
| ALIAS (deny := episodeDenyWrite, append into shared backing) | KILLED by NoProtectedIsByteIdentical/base (+others) | :801 backing array mutated |
| HANDLER-DROP (ctor copy removed; and separately runVariantFor line removed) | both KILLED by ThroughHandler/true | :886 |
| VERIFY-DROP (VerifyRunCaps passes nil) | KILLED by ThroughHandler/true | :899 |
All seven executor claims hold.

## My diff-anchored mutants
| Mutant | Result | Note |
|---|---|---|
| ctor copy removed (h.protected = cfg.Protected) | KILLED (ThroughHandler :886) | caller-mutation check works |
| runVariantFor copy removed (spec.protected = h.protected) | SURVIVED - EQUIVALENT | renderPolicy copies into a fresh deny slice, nothing retains spec.protected |
| first EMPTY-BYTES form (`protected != nil && len==0`) | survived - EQUIVALENT | unreachable: append([]string(nil), {}...) is nil |
| backslash escape removed | KILLED Golden/escapes :744 | |
| `\\` emitted as single backslash | KILLED Golden/escapes :744 | |
| IsControl -> r < 0x20 (DEL and C1 raw) | SURVIVED | see finding N1 |
| IsControl -> r < 0x7f (C1 only raw) | SURVIVED | C1 raw is legal TOML; harmless |
| `\u%04X` -> lowercase hex | SURVIVED - EQUIVALENT | TOML accepts both; test only probes \u0000 |
| escape all non-ASCII | KILLED (Golden, base golden contains U+017F) | |
| protected placed before built-in denies | KILLED Golden :717, AllVariants :765 | |
| trustedHostGitDeny moved after episodeDenyWrite | KILLED (KeySetsPerClass :128, NoProtected/Env :797) | |
| RenderRunVariant base path passes nil | KILLED AllVariants/base :765 | |
| RenderEpisodePolicy forces spec.protected=nil | KILLED Golden :717 | |
| daemon caller passes non-nil list | KILLED (TestWorkspacePolicyRenderedPerEpisode, ModuleRootIsTheSandbox, PackageCache...) | nil-by-design pinned |
| remove AilangToolConfig.Protected field | compile error (not a kill; irrelevant, covered by HANDLER-DROP) | |

## Findings
N1 NON-BLOCKING: escape test cannot distinguish raw U+007F (DEL) from `\u007F`. Repro: change `case unicode.IsControl(r):` to `case r < 0x20:` in tomlString; all M5a tests pass. Cause: protectedPolicyDeny decodes with encoding/json, which accepts raw DEL, whereas TOML basic strings forbid it; the only literal check is `back\\slash\u0000`. Fix: also assert Contains `\u007F` in the escapes subtest. Production code is correct today.

N2 NON-BLOCKING: the JSON decoder is a sound proxy only for the subset the renderer emits (\\, \", \uXXXX, raw printable/non-ASCII). It is not a TOML validator: it accepts raw DEL, \/ \b \f, and lenient invalid UTF-8. The comment states the subset honestly, but "independently decodes" overstates it for control-character validity (see N1). Byte-level goldens (base/quote) are genuinely independent of the renderer: the legacy deny list is a hand-written literal and variants are string substitutions; they match production byte-for-byte.

N3 NON-BLOCKING: tomlString ranges over runes, so invalid UTF-8 in a pattern is silently rewritten to U+FFFD rather than refused (renderPolicy's root check rejects ReplacementChar, protected patterns get no equivalent check). Repro: RenderEpisodePolicy(root, []string{"a\xffb"}) renders "a�b" bytes with nil error, so the rendered deny no longer matches the input. Patterns come from fileguard (M5b wiring) so likely unreachable; M5b should validate or reject. Not an M5a plan requirement.

N4 NON-BLOCKING: AllVariants asserts relative to the production-rendered unprotected variant; absolute pinning of the unprotected variants is supplied by NoProtectedIsByteIdentical goldens, so the chain is sound.

Vacuity: with the production change reverted the tests do not compile (new signature) and Golden/AllVariants/ThroughHandler would fail on assertions; every guard has a RED-on-assertion mutant above except the items listed as equivalent or N1.

Plan scope: all M5a items present (RenderEpisodePolicy/RenderRunVariant signatures, Config/handler plumbing, runVariantFor, VerifyRunCaps, TOML-safe encoding, byte preservation for nil/[], all caller sites, three named tests). D-156-4=A pre-check correctly absent.

UNMEASURED: broker and daemon full -race suites (controller figures only); AILANG's real TOML parser behavior on raw DEL; whether fileguard grammar can admit control characters.
