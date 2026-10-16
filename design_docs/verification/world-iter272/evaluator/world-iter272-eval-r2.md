# Evaluator report r2 - row 150 M3a (0dd8f60, diff f50ec08..0dd8f60)

Score: 90/100 - PASS (no blocking findings)

## Gates (rc=0 each, top-level --- PASS)
^TestDemoProof (3 tests), same with -race (ok), ^TestDemoSeed (7), the two plan base gates, go vet rc=0, gofmt -l empty. demo_proof.go diff f80502f..0dd8f60 = 0 lines (byte-identical). Milestone diff still exactly the two allowed files.

## r1 survivors vs new tests (restored after each mutant)
- writer-parent: KILLED, RoundTrip:253 carrying entry prevEntryHash != observed head log head
- writer-epoch: KILLED, RoundTrip:256 semanticsEpoch 3 != head entry epoch 2
- writer-503-class: KILLED, AdapterGuards/writer-503-other-class-refused
- reader-malformed: KILLED, AdapterGuards/reader-malformed-body-refused (:359)
- reader-meta-prov: KILLED, AdapterGuards/reader-valid-object-positive (:406)
- busy-reader / busy-writer (5s): KILLED, busy-reader-zero / busy-writer-zero
- reader-nopayload: still red, now by recover() -> t.Fatalf naming the guard and the panic (:344), not a raw panic. Acceptable.

## Vacuity / over-fit check
- Read-back and inventory go over real REST (client against httptest daemon; store closed before daemon start).
- prevWorld/prevEntry are read via proofHeadEntry AFTER the source commit and BEFORE demoProve, so they are the pre-commit head, not derived from the result. Assertions compare entry index to Revision+1, prev hash to the observed LogHead, epoch to the observed head entry. Genesis epoch 1 and the observed values differ from the mutated ones (killed with concrete numbers), so not vacuous. Mild over-fit: the carrying entry is selected by transitionRef+interpreter (last match), fine.

## Remaining non-blocking (retained from r1)
1. demoProofHTTPTimeout (10s) is used only by tests; production adapters rely on the caller's client timeout (bound-http still survives). Enforce in M3b wiring.
2. "No competing writable store" is structural only (store closed before daemon start), no assertion.
3. Configured 90s whole bound / cancel not asserted (bound-whole-nocancel survives); checker-readfile and prove-checker-err guards still uncontrolled (masked downstream).
4. Plan's M3a-no-rest as worded still reds at an earlier assertion (writer lock); executor's variant plus the writer-revision mutant (RoundTrip carrying-entry assertion) satisfy the intent.

Worktree restored; git status clean.
