# TR.A1 mutation transcript (iteration 71)

All Go mutants below were applied one at a time. Each landed by a changed anchor and SHA-256,
compiled with `GOTOOLCHAIN=go1.25.6 go build ./...`, was exercised with a scoped `-run`, and had a
package inverse arm using `-skip` return rc=0. Restores used `cp` and reproduced the recorded before
SHA exactly. The common restored SHA values were:

- `host/store/store.go`: `eb54d3d8c179aadd29d273bcc8021a3666c9569a324da9da672352c029422a64`
- `host/transitionreg/transitionreg.go`: `674cc001185fd4d05f84e4f426e1ba32477590f3a224767104fb9750c007a86d`
- `host/transitionreg/codec.go`: `11ad144ec17e6291cecdf2540f8c74031cd9b7c48353705116867d168aa5ba49`

| mutation | landed mutant SHA | build | result and failing assertion | inverse |
|---|---|---:|---|---:|
| MUT-CAS-BLIND | `86810ee7d6fb9877c22904e4db424f30ab66b508b763585ea0221b97345e1d70` | 0 | KILLED — `TestCompareAndSetRegistryHead/stale_expected_conflicts`: conflict nil, err nil | 0 |
| MUT-CAS-DANGLING | `d9e08227e1ae9036ee7c0c665c8ab7665379ced75d4bd3374b2da3af469a4bb8` | 0 | KILLED — `TestCompareAndSetRegistryHead/dangling_next_refused`: accepted dangling next | 0 |
| MUT-CAS-EPOCH-HEAD | `06e31f44dbb482a3bdb7edf3030045f4494e68b7be78ce65e73460ee8962f9c8` | 0 | KILLED — `TestCompareAndSetRegistryHead/epoch_registry_isolation`: epoch head changed | 0 |
| MUT-ID-ACCEPT-EMPTY | `b105bf3d754e1437066bb139d1d54a5278a97f26aa0420208684a4be019c6512` | 0 | KILLED — `TestDescriptorValidationRefusals/id_grammar`: invalid value accepted | 0 |
| MUT-ID-ZERO-FN | `9b3e399f2118f98ea96bee051676c70a258f4353cc51373c2add0513ea5f02ac` | 0 | KILLED — `TestDescriptorValidationRefusals/zero_transition_fn` | 0 |
| MUT-ID-ZERO-INTERP | `e7d3f79a670687db4a2c6d25f4ded1037a8ac135eda44594fd228935d0b9adde` | 0 | KILLED — `TestDescriptorValidationRefusals/zero_interpreter` | 0 |
| MUT-SCHEMA-ANY-BYTES | `de903086dab9b94b2fa9bccdba5d22cd8e1563a0a647d3f4935671176b4a7f17` | 0 | KILLED — `TestDescriptorValidationRefusals/schema_not_an_object` | 0 |
| MUT-SCHEMA-NO-LIMIT (raw) | `0d3bc38897b146812983b0456b053a9392b50ee37384e595e5942b2ed38271bc` | 0 | KILLED after fixture repair — `.../schema_raw_over_262144` | 0 |
| MUT-SCHEMA-NO-LIMIT (canonical) | `4447d4559d67f8faedc6d35a41095a7b52c18f5c14ff046522d7dcdde804615a` | 0 | KILLED — `.../schema_canonical_over_65536` | 0 |
| MUT-CODEC-NO-KEYSORT | `a95e5536fe963adf6cf564327cb8c6b5befef56233985c48349ecff6789febd2` | 0 | KILLED — `TestCodecGoldenRoundTrip`: literal empty golden byte mismatch | 0 |
| MUT-CODEC-INDENT | `fbabd632100f57640434abc8cecda8e6286585421f8edc91fdb6da8bd4fb5603` | 0 | KILLED — `TestCodecGoldenRoundTrip`: literal empty golden byte mismatch | 0 |
| MUT-NEG-EPOCH-OK | `f8c525453090be964e78117125c74c3cf76d7c0e989ad9cd08e7c9ad6ff23c38` | 0 | KILLED — `.../negative_semantics_epoch` | 0 |
| MUT-NEG-COST-OK | `4737762427c2566612d20c9d9987b5795dc4a8a934dfcb5d9a96ef545a1fb16a` | 0 | KILLED — `.../negative_cost` | 0 |

`MUT-CAS-EPOCH-HEAD` control: `TestRegistryHeadRoundTrip` returned rc=0 on the landed mutant. Only
the isolation subtest fired, confirming the controller's measured correction.

`MUT-SCHEMA-NO-LIMIT/raw` initially survived. The fixture used a huge string and therefore also
hit the independent canonical-size guard. It was repaired to call `canonicalSchema` with
262,144 spaces plus `{}`: the canonical form is two bytes, so only the raw guard can reject it.
The repaired mutant then failed the named subtest.

`MUT-DELETE-TR-A-TEST` renamed `TestCodecGoldenRoundTrip`: mutant test SHA
`5e34896a3bc8235c09822fd7b786e3b282a9bea09ca544c09b6412f28fca3668`, build rc=0, AC1 count
became 2 and AC10 count became 0; both activated gates returned rc=1. The inverse returned rc=0 and
restore reproduced `90f5a9377cb305b846aef38df902ddc876505fe7a97a88f29d17ca2ad3a825f0`.
Renaming `TestCompareAndSetRegistryHead` gave mutant SHA
`42f098ff1a33c076124124af57cc06156229fbbee7354108a7194ec0da875a45`, build rc=0, AC4 count 1
and gate rc=1; inverse rc=0; restore reproduced
`06d315fa534dc0a8a430ce2c2cd83e6faa32b32c78d31160bd262f17cd34f5fb`.

`MUT-AIL-EMPTY-MODULE` landed as a fifth `world/*.ail` file with SHA
`cdbc91b73f9b3e78c7009e1740439aa303cd29952a62dc76b5d4c89458e94646`. The underlying script
returned rc=0 and printed 12 modules; repaired AC9 observed `modules11=0` and returned rc=1. The
file was deleted and the world file count restored to four.

## Six arms deferred by the executor, then run by the controller

The executor did not execute these six T4-assigned arms and said so, claiming no kill. The
controller ran all six. Every mutant was asserted LANDED (differing sha256 over a single-match
anchor) and BUILDING (`GOTOOLCHAIN=go1.25.6 go build ./...` rc=0) before its result was read;
each kill arm was `-run`-scoped, each had an inverse `-skip` arm requiring rc=0, and every
restore came from a `cp` backup and was verified byte-identical.

| Mutation | Landed (sha) | Builds | Result | Failing test |
|---|---|---|---|---|
| `MUT-ID-NO-LENGTH-BOUND` | `674cc001…` → `31488eb8…` | rc=0 | **KILLED** | `TestDescriptorValidationRefusals/segment_too_long` |
| `MUT-CJSON-DUP-KEY-OK` | `11ad144e…` → `8887f146…` | rc=0 | **KILLED** | `…/duplicate_schema_key_nested` |
| `MUT-CJSON-NO-NUMBER-BOUND` | `11ad144e…` → `071d6791…` | rc=0 | **KILLED** | `…/number_coefficient_overflow` |
| `MUT-CODEC-NUMBER-RAW` | `11ad144e…` → `76f129d9…` | rc=0 | **KILLED** | `TestCodecGoldenRoundTrip` |
| `MUT-CJSON-SURROGATE-OK` | `11ad144e…` → `707b1ce8…` | rc=0 | **SURVIVED → now KILLED** | `…/lone_surrogate` |
| `MUT-CJSON-UNKNOWN-KEY-OK` | `11ad144e…` → `083410ce…` | rc=0 | **SURVIVED → now KILLED** | `…/unknown_revision_key` |

**The two survivals were genuine, and they share one cause.** Both mutants landed and built, so
neither was instrument failure. The refusal table asserted only *that* an error was returned,
never *which branch* returned it — and `DecodeRevision` performs a canonical re-encode comparison
(`codec.go`), so every guard has a second refuser standing behind it. Neutering the named guard
left the input still refused, by the backstop, and the message-agnostic assertion stayed green.
On the clean tree the named guards *do* fire, which is exactly why this was invisible.

Repaired in `c0e72de`: all 16 refusal cases now pin their own measured message. Re-run against the
strengthened test, both survivors KILL, and the failure text names the masking mechanism —
`input schema is not canonical` and `revision is not canonical for its typed schema`.

The masking is specific to inputs that are *silently transformed* during decode→re-encode
(a surrogate becomes U+FFFD; an unknown key is dropped). Duplicate-key, ID-length and
digit-overflow violations are outright rejections with nothing to mask them, which is why they
killed on the first arm. The evaluator swept for a third masked case independently and found none.

T5–T8/TR.A2 mutations are outside TR.A1 by controller direction.
