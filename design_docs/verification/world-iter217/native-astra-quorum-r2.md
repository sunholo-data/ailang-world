# Native Astra design quorum — iteration217 round2, r4

**REJECT immutable r4 — one narrow syntax/type defect; zero direction objections. Exact measured refinement accepted.** Design review only; implementation acceptance NOT RUN, product/admin score none. Transport: native Agent `gpt-6-astra`, author `gpt-6.1-sol`, **judge-independence: same-vendor** (different model, separate judge context) under D-WORLD-48. This seat does not rewrite CLI synthesis or override another present reviewer.

Exact snapshot `077027ed1d74f62151947cacacb62e2a689fbead`; canonical design SHA256 `8a596295fc540afd4d568896f0d1a4aeef248182037bc695bc50c25904a31086`. Read the full236-line r4, all freeze choices and carried objections against the design-doc premise/conflict/axiom gates and sprint-evaluator rubric. The latter's implementation score is not used to manufacture a product verdict at design stage.

## Strongest objection and its disposition

The strongest r3 objection was real: allowedDescriptors interprets any snapshot failure as empty success when its precheck saw no head. This reviewer missed that interaction in r1. R4 now requires typed confirmed absence, preserves wrapped genuine errors, repairs the shared helper before MCP reuse, and specifies non-vacuous list plus internally observed Invoke tests. That is a complete semantic design-level answer. One narrow Go-form defect remains in the immutable text, detailed below.

**Catch:** In the proposed repaired helper, restore the old catch-all while retaining the typed reader. The new list test must turn RED because the response becomes empty SSE success. Invoke must also assert the classified error/operator evidence because the frozen generic envelope masks the difference between genuine failure and ordinary not-allowed refusal. Test true absence, absent-then-success, present-then-absence, wrapped typed absence and a genuine same-text error.

**Proposed fix:** Correct the declared errors.As target as specified below; no direction revision is required. Implement D108-12 before MCP reuse and preserve its exact controls and mutation gates through planning. Treat A2A unchanged as unchanged normal wire contract except for the explicitly specified raced-failure correction. Prove the reader checks cancellation after the head lookup before creating a typed absence.

Combined r1 was correctly BLOCKED despite my prior PASS. R4 responds to a verified source defect, not a requirement for extra reviewer independence. Its D108-12 is a **narrow concrete classification repair**, not a new attended direction decision.

## Why D108-12 closes the blocker at design level

Independently re-read `projection.go:375–407`, `transitionreg.go:70–130`, `bind.go:116–128`, and `projection_test.go:575–670`. The existing helper swallows every rerr on absent precheck; ReadSnapshot has separate context/head/object/integrity/codec failures; NewRequest wraps with `%w`. Existing tests exercise true absence, present failure and absent→success but miss absent→failure. A same-call `rg` scan confirms the proposed RegistryHeadAbsentError name is unallocated, with the existing untyped absence and other error types as controls. Read the controller's stored repro source/log showing200/skills[] instead of503; that execution is attributed to the controller, not claimed as mine.

R4 requires the type only at successful head lookup reporting absence with live ctx, and empty success only when BOTH the precheck was absent and errors.As identifies that type. Therefore:

| Precheck / snapshot | Required result |
|---|---|
| Precheck fails | Propagate without snapshot |
| Absent / confirmed typed absence | Empty allowed set |
| Absent / successful fresh snapshot | Use fresh populated snapshot |
| Absent / genuine store, integrity or context error | Propagate genuine failure |
| Present / typed absence or other failure | Existing unavailable behavior |

The exact old catch-all mutation is non-vacuously killed on **list** response shape. The Invoker path acknowledges that generic upstream errors can look identical on the wire, requiring internal classified-error/sink evidence and reached-stage counters. Same-text genuine errors kill string matching; a real StoreReader raced-in-head/GetObject failure leg ties the helper test to production classification. Typed `%w` propagation, zero task mint/Dispatch and real-store receipt absence are specified. No A2A wire taxonomy or coordinator law is invented.

## Full design/freeze re-review

- **D108-1 dependency: accept.** Actual post-import World closure remains explicitly unmeasured and mandatory. Re-read prefix allowlist and mount source; compare12 relevant source files and go.mod byte-for-byte with my r1 snapshot: all identical. My r1 independent measurements (seam exactly protocol/hostcall/mcphttp, current daemon54 packages, four old protocol files identical, proxy/sumdb200) remain prior evidence, not newly rerun claims. Landing still requires fresh fetchability and exact resulting graph.
- **D108-2 mapping/schema: accept.** Re-read World grammar/canonicalSchema and released descriptor validation. Escaped underscores make M reversible; >64-byte admission is explicitly partial and fails the entire surface. R4 now specifies bounded per-POST upstream error, no startup poison, and same-handler recovery after registry repair. Normalization retains RawMessage or UseNumber precision, tests9007199254740993, and keeps typed-object bytes. Output schema is unchanged. These are explicit compatibility costs, not silent truncation or discarded constraints.
- **D108-3/4 mount/authority: accept.** Re-read daemon.Handler/isProtected and A2A admission. One resolver with bearer-only authority, no registry reads for denial, one additive POST mount and no new protected middleware path. D108-12 shared-helper repair is prerequisite, in file/milestone/conflict scope. AST mount identity plus same-parse A2A control still kills the methodless mutation that upstream405 hides from network tests. Existing source-guard narrowing must retain A2A's own name/codec prohibitions.
- **D108-5/10 bounds: accept.** Re-read Runner, daemon constants and procbound. One20s parent context clips all sequential callbacks;20s per-call and8 slots do not imply process-capacity identity. Actual callback return, recovery and teardown remain required.3+10+20K arithmetic and32.6/34.6s mutant recipes are coherent. R4 explicitly delays only initial Tools admission so fresh Invoke admission reaches capsule stimulus. Body/CPU/write overhead and noncooperative callbacks remain residuals; no universal20s completion guarantee is asserted.
- **D108-8/9/11 batch/durability: accept.** Re-read released authorize/serveMessages/callTool and coordinator Dispatch. Once-per-POST authorization plus K sequential invokes, whole-POST error/null id, K+1 successful admissions, per-item random IDs and independent commit/receipt accounting match the code. A failure can hide prior committed outputs; retry can duplicate them. That remains the largest accepted product limitation, explicitly documented and tested, with no idempotency invention. Genuine admission failures must now escape the repaired helper before task mint or Dispatch.

All three pending direction freezes accepted **by this seat**; D108-12's semantics also accepted, subject to the one exact syntax/type correction below. Net axiom+6 remains reasonable, no hard A1/A3/A4/A7 negative found. Conflict surface includes the new shared-reader/A2A repair. No additional attended direction ruling is required by this review.

## B1 — exact narrow refinement (final r2 disposition)

After the full semantic review, the controller forwarded Google's r2 finding. It is correct: immutable r4 line80 literally says `errors.As(rerr, &RegistryHeadAbsentError)`, taking an address of a type. That is not a valid Go expression. Treating it as conceptual shorthand would leave the design's executable recipe wrong. **Final verdict on the immutable snapshot is REJECT for this one narrow defect**, not a rejection of the typed-classification direction.

I independently read the controller's `errors-as-snippet-repro.json`, including both sources and raw outputs: original rc1, `RegistryHeadAbsentError (type) is not an expression`; declared target rc0 with `true` for wrapped typed absence and `false` for genuine error. I did not execute this scratch program myself. I also read the amended line in the controller worktree: declare `var target transitionreg.RegistryHeadAbsentError`, then test `!hasHead && errors.As(rerr, &target)`, and define/emit a value error with value-receiver `Error()` so the assignable target representation matches. This is correct; a pointer-emitted variant would require a corresponding pointer target instead, so retaining the explicit value representation matters.

**The exact measured refinement is ACCEPTED.** It is a narrow concrete syntax/type correction preserving the same classifier and all accepted freeze choices, suitable for the supplied Gate2 refinement carveout. It creates no attended direction question. This supplemental acceptance does not rewrite the original CLI or immutable-snapshot verdict, does not declare implementation tests passed, and does not waive the future catch-all/string-match mutation gates. The JSON pins both the original design hash and the observed refined controller-draft hash.

## Nonblocking implementation cautions and limits

1. Pin post-head-read cancellation in the typed reader: a head double returning nil error/ok=false after cancel must yield context error. This implements r4's “ctx still live” requirement; it is not a new direction.
2. “A2A stays unchanged” is overly broad shorthand beside an intentional bug fix. The detailed D108-12 contract controls: normal vocabulary/true absence/race success remain, masked genuine failures change to unavailable. Clarify the shorthand in planning prose.

Previous Astra notes on schema precision, initial-Tools-only timing, and sequential baseline prose are incorporated. Saved sequential baseline is AIL rc0 and Go rc1; no cause or exemption is inferred. Final exact full-green gates remain mandatory.

No source/test implementation, code mutation, git mutation, scratch code, additional agent or inbox action performed. No background work outstanding. Reports are this round's only writes. Proposed tests and mutations were not executed, and neither this design review nor the narrow-refinement acceptance may be quoted as implementation acceptance.
