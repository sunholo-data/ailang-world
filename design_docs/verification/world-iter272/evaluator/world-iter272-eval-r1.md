# Evaluator report r1 - row 150 M3a (f80502f)

Score: 82/100 - PASS (no blocking findings)

## Gates (all rc=0, top-level --- PASS)
- TestDemoProofHTTPProducerRoundTrip, TestDemoProofSourceHashRefusal (cmd/ailang-worldd): PASS
- TestProducerEndToEndMintsEnvelopeValidatorAccepts (host/evidence), TestWorkbenchGradeProvenFromSubjectBoundEnvelope (host/daemon): PASS
- ^TestDemoProof (3 tests incl. TestDemoProofAdapterGuards) PASS; ^TestDemoSeed (7 tests) PASS; go vet rc=0; gofmt -l empty.
- Diff scope: exactly demo_proof.go and demo_proof_test.go (741 insertions). No production caller of demoProve/adapters.

## Design findings
- Real v0.41.0 producer, real daemon via httptest over REST, fresh crypto/rand key, store handle closed before daemon starts. PROVEN for exact subject; different subject refused as subject_mismatch (named negative). Missing AILANG_BIN is t.Fatal, and requireBinaryRelease pins v0.41.0 (FAIL not SKIP).
- reader/writer parameters on demoProve = the natural NewProducer shape (evidence.ObjectReader/ObjectWriter), unexported, no global/hook/env seam. Not a new injection seam in the plan's sense.
- UncertainError: transport error, 503 CommitUncertain, and 200-with-wrong-head all -> store.UncertainError; single POST each (posts==1 asserted); plain refusals (409) not uncertain. No retry loop in code.
- Checker pinned from health: version line, ref parse, archive resolve, bytes re-hashed. All four have negative controls.

## Mutation table (production demo_proof.go; run: TestDemoProof* ; mutant -> rc -> reds at)
Killed (rc=1, at an assertion naming the guard): reader-status (AdapterGuards/reader-status-refused), reader-404, reader-hash, reader-size, reader-iface, reader-meta-semantic (RoundTrip:207), reader-nopayload (red via nil-pointer PANIC at test line 328 - non-assertion; confirmed), writer-nohead, writer-noworld, writer-noentry, writer-interp, writer-transport-plain, writer-503-plain, writer-status, writer-headmatch, writer-revision (RoundTrip:247 carrying-entry assertion), writer-token, writer-obj-payload, writer-obj-iface, writer-obj-sem, checker-version, checker-ref, checker-resolve, checker-hash, checker-compiler, checker-version-field, bound-obj, bound-exec, bound-max, bound-whole, cfg-readtimeout, cfg-writetimeout, prove-noreq.
Survivors (rc=0): reader-malformed (ignore JSON unmarshal err; hash guard downstream masks), reader-meta-prov, writer-503-class (any 503 -> uncertain), writer-epoch (+1 epoch accepted), writer-parent (wrong parent log head accepted), checker-readfile (ReadFile error guard; path empty case masked), bound-http (13s HTTP bound), busy-reader/busy-writer (BusyTimeout 5s), bound-whole-nocancel (100x whole bound), prove-checker-err (masked by NewProducer config check).

Plan mutations as worded:
- M3a-source (flip a source byte in forwarding read arm of the RoundTrip test): rc=1, RoundTrip line 199 "source ... bytes hash to sha256:562e..." calls=0 -> reds at the producer hash guard as required.
- M3a-no-rest (forwarding writer stores via real store.PutObject): rc=1, but at RoundTrip:199 (SQLite writer lock held by daemon: "only one writer may be active"), NOT at the carrying-entry assertion. Reproduced the executor's claim. The plan wording is unattainable as written (daemon holds writer lock); executor's variant (commit with foreign transitionRef) is the right intent, and my independent writer-revision mutant independently reaches line 247 ("object exists without a commit"). Accepted as satisfying intent.

## Non-blocking findings
1. demoProofHTTPTimeout (10s) is used only by tests; production adapters take whatever *client the caller supplies, so the 10s < 12s ordering is not enforced by M3a code (bound-http survives). Enforce in M3b wiring or set it in a constructor and test it.
2. BusyTimeout()==0 on reader/writer is unasserted (busy-* survive); add a one-line assertion.
3. "Assert no competing writable store" is structural only (store closed before daemon start); no assertion.
4. Survivors writer-epoch/writer-parent/reader-malformed/reader-meta-prov/503-class have no control; epoch/parent are copied from head and daemon does not reject mutations.
5. reader-no-payload-refused reds via panic (nil deref in test line 328 when guard removed) rather than an assertion; still red.
6. Whole-90s bound is only exercised by shrinking it, no test that 90s is the configured value.

Worktree restored (git status clean).
