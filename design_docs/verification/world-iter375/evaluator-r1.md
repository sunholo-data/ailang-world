# Evaluator r1 — iter 375, row 201 A1b sprint 1 (048537c3)
Verdict: PASS 91/100. Sub-scores: correctness 28/30, tests+mutation 26/30, security 20/20 (properties confirmed), standards 9/10, scope 8/10.
Mutants: 100 run (plan section 4 set + 25 of mine); 87 killed by the named test; 13 survived the named test;
after running the whole package, 5 survivors are equivalent/unreachable, 1 message-only, 1 propagation, 2 real gaps.
Real survivors: PAT-UPPER2 (uppercase allowed in non-first name chars; approver.go:25; only first char tested), LK-NAMEPREFIX (Approver() by name LIKE-prefix; approvers.go:135; only key prefix tested).
Equivalent/unreachable: RE-FORMAT, EN-INSERT-FORMAT, RV-N1-GT, CL-BEHIND-GE, CL-DEFAULT. Message-only: EN-NAMEQ. Propagation: LS-ROWSERR. Harness artefact: E4-GE-b (wrong test targeted).
BLOCKING: none. host/broker TestExecProbeNet6 flake in go test ./... is unrelated (package untouched).
Each commit builds, vets and passes host tests. Stub non-vacuity: M1b, M2a, M2b, M2c headline tests fail. Fixtures NUL-free, SHA-pinned (append to p384.pub fails the pin test).
