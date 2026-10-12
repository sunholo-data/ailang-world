# Evaluator report, iteration 254, row 105 (8390c6b vs edb1b63)
SCORE 88/100, PASS (no blocking findings).
Rubric: AC coverage 27/30; test non-vacuity 21/25; honesty/security 18/20; gates/scope 13/15; deviations/process 9/10.

Gates (AILANG_BIN pinned v0.41.0): go vet ./... clean; go test ./host/daemon ./cmd/ailang-worldd green; go test -race ./host/daemon -run 'Grade|ProofKey|Workbench' green. Scope: host/evidence, store, workbench, scripts, .github, world unchanged; workbench_test.go untouched.

Deviations: (1) `class="verdict-` / `aria-label="test verdict` are correct: render.go:168 stylesheet always has .verdict-*; markup is render.go:208. Injecting HasVerdict into the PROVEN view reds 3 tests, so the narrowed marker still kills M10. ACCEPTED. (2) `(*proofKeyFile)[:0]` is needed because dropping the field leaves an unused var (BUILDFAIL); mutant reds TestServeProofFlags. ACCEPTED. (3) Excluding subject_mismatch is right: it is folded to bound-elsewhere by the scan loop, not by the map. ACCEPTED. (4)/(5) not material to correctness.

Honesty: sentences follow validator order (MAC before subject, then tool, proof, identities), so "subject-bound ... authenticates under the daemon's key" sentences are only reachable after MAC and subject checks pass. Unauthenticated sentence claims nothing about binding. Fallback is an invariant sentence. OK.
Security: key opened, Stat on open fd, regular only, LimitReader 33, exact 32, errors name path/sizes not bytes, no logging in file, no context.Background in production, half configs fail closed at StageConfig.

Mutations (own, from production diff; command: go test ./host/daemon ./cmd/ailang-worldd -run 'Grade|ProofKey|Workbench|ServeProof'):
A remove unconfigured early return: RED Unconfigured
B fetch +1 -> +0: RED ScanTruncation
C truncated == -> >: RED ScanTruncation
D drop candidates[:100] slice: GREEN  (survivor)
E key >32 bound shifted: RED ProofKeyWiring
F short-key check removed: RED ProofKeyWiring
G LimitReader +1 -> +0: RED ProofKeyWiring
H non-regular check bypassed: RED ProofKeyWiring
I TrimRight -> TrimSpace version line: RED DerivationsAgree + others
J whole verbatim instead of first line: RED
K Resolve Proven ignored: RED MapsEvery, ValidatorRefusal
L wrong subject passed to ValidateProof: RED (many)
M subject_mismatch treated as failure: RED BoundToOther
N firstFailure -> last failure wins: GREEN (survivor)
O truncation qualifier dropped: RED
P res.Err swallowed: RED ValidationReadError
Q scan error swallowed: RED ScanFailure
R fallback sentence changed: RED
S map tool_mismatch -> proof_failed: RED
T interpreter-hash check removed: RED
U identities-without-key check removed: RED
V Compiler hash omitted: RED
W ObjectReadTimeout 0: RED
X errGradeNoVerdict -> continue: RED MapsEvery
Y workbench error classification bypass: RED (3 sanitized tests)
AA identities flag not wired: RED ServeProofFlags
AB2 empty-identity flag check off: RED ServeProofFlags
AC configureGrader not called: RED
AD context.Background in scan: RED
AE2 label PROVEN -> TESTED: RED
AE3 HasVerdict on PROVEN view: RED (3)
Z/Z2 ProofKeyFile removed from Config: BUILDFAIL; M17a form [:0]: RED
AF unbounded io.ReadAll(f): GREEN (survivor; only observable with a huge regular file)
(First attempt AB matched a different flag's `v == ""`; redone as AB2.)

NON-BLOCKING findings
1. Mutant D survives: the test for the scan bound never makes the 101st candidate the subject-bound valid one, so the "first 100 only" claim in the rendered qualifier is not pinned (a valid envelope past the bound could render PROVEN). Repro: delete `candidates = candidates[:gradeCandidateScanLimit]`, tests stay green. Fix: add a 101st valid envelope sorted last; expect the bound-elsewhere + truncation sentence.
2. Mutant N survives: "first failure wins" precedence is untested (sentence is true either way).
3. Mutant AF: bounded key read not behaviorally pinned (equivalent except for memory on huge regular files).
No blocking findings. Worktree restored; git status --porcelain empty.
