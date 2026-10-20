# Evaluator r1 - iteration 289, row 156 (V25 + M5e) - 1b61eba9 vs 0ffbb9c9

SCORE 90/100 - PASS. Correctness 28/30, test quality 21/25, plan/design fidelity 18/20, evidence honesty 13/15, hygiene 10/10. BLOCKING: none.

## Checks
A. V25 re-measure: PASS. `python3 -I <script> <bin> <scratch>` for both scripts on both binaries; `diff` against the banked .out.txt is empty for all four (v7r_v0521, v7r_v05415, v7rb_v0521, v7rb_v05415). shasum of ~/.pinned-ailang-tools/v0.54.15/ailang = d821dad0...ac043 = host/pins/ailang-releases.json tool binary_sha256 (darwin). v0.52.1 = 0dd70a1d..., matches binaries.sha256.txt. Re-confirmed: 52.1 writes all symlink/hard-link aliases; 54.15 refuses with host=none; without the protected set every 54.15 symlink arm is written; multi-link refusal fires with EMPTY fs_deny_write (v7rb last row); `**/*.test.ts` refuses depth 0/1/2 on 54.15 (depth 1 only on 52.1); `tests/**` does not match cli/tests/x.py.
B. Design amendment: supported by the evidence. Direction unchanged: 4.3.2 and D-156-4 say "neither built nor dropped" until D-WORLD-86; D-156-4 stays "Answered A (D-WORLD-82)", which matches the ledger row (I read it). R-156-10 is mostly accurate (hard-linked files unwritable through the tools) but see NB1/NB2.
C. Plan/design/JSON consistency: PASS. Table summed with awk: 24 rows, 5407 LOC, 670 min; held M5c+M5d 500 LOC/60 min; 4907/610 if dropped - matches the totals line and the JSON estimate. JSON milestone ids/statuses/depends_on match the table (M5e planned dep M5b; M5c/M5d held:D-WORLD-86; M6-pre deps M5b,M1d). D-156-1..4 answered A in plan, JSON scope and design. SHAs re-derived with `git log --merges --oneline 0ffbb9c9 | grep '#n '`: #312=02bd1410, #314=78f2860c, #319=826f2a0b, #322=c9ab49b7, #325=0adeb703, #326=96d49010, #328=7fb99dae - all match.
D. M5e: real pin (WORLD_TOOL_AILANG_BIN=...v0.54.15) all 14 subtests PASS (14.5 s); pin unset -> discovered SKIP before any fixture.
E. Gates: go vet ./host/daemon/ rc 0; gofmt -l empty; go test -race -count=1 ./host/daemon/ ok 286 s (pin unset, so the real-tool test skipped there, as designed).
F. Hygiene: diff is exactly the 12 allowlisted files (docs/evidence/JSON + host/daemon/workspace_protected_test.go, test-only); check_no_personal_email.sh passes; no address-shaped added lines. Scratch deleted; mutated files restored and `shasum -c` OK after every drill.

## Mutant table (own mutants, one at a time, real v0.54.15 unless noted)
| Mutant | M5e result | Other tests |
|---|---|---|
| (i) retired v0.52.1 + temp ToolBinaryRelease hunk | RED: 5 symlink arms + hardlink-write + hardlink-edit + no-protected-hardlink; GREEN: 5 no-protected-symlink + singlelink (exactly as predicted) | - |
| (ii-a) workspace.go RenderEpisodePolicy(sandbox, nil) | RED 5 symlink arms (hardlink arms stay green) | - |
| (ii-b) AilangToolConfig.Protected: nil (policy still has set) | SURVIVES (PASS) | killed by existing TestWorkspaceProtectedPolicyRequests |
| (ii-c) un-rebased protected (rel ignored) | SURVIVES (episode root == sandbox, rebase is identity) | killed by existing TestEpisodePolicyProtectedRebased/* and TestWorkspaceProtectedPolicyRequests |
| (ii-d) protected[:1] (drops one pattern) | RED symlink-dir-* (3 arms; canonical order puts *_test.go first) | - |
| (iii) skip creating lnk / alias.txt / hl.txt in symlink-dir-new-write, symlink-file-edit, hardlink-edit | each RED on the mechanism assertion ("resolves to"/"fs_deny_write" or "more than one hard link" missing; tool returned ok:true or "no such file") plus the setup-shape assertions | - |
| (iii) same neuter in no-protected-symlink-file-write, no-protected-hardlink | RED (bytes / ok / shape) | - |

(iv) Vacuity: symlink arms require BOTH "resolves to" and "fs_deny_write", so a lexical or generic refusal cannot satisfy them; the no-protected controls prove the aliases really alias. Sole killer for wiring mutants (ii-a, ii-d) = the symlink arms; hardlink-write/edit are NOT sole killers of anything about protection wiring (they stay green under ii-a) because the refusal is unconditional.

## Non-blocking findings
NB1. R-156-10 says the refusal happens "because World's episode policy always has a non-empty fs_deny_write". V25 itself shows the refusal is unconditional (empty fs_deny_write also refuses). The "because" implies a condition the evidence refutes. Repro: `python3 -I design_docs/verification/world-iter289/v7r_doublestar.py <54.15 bin> <scratch>` last row. Fix: drop the clause.
NB2. R-156-10's examples (pnpm stores, uv caches) are unmeasured; also only the policy-tool path was measured, not `ailang-run` FS effects (the README does say this for R-156-3).
NB3. D-WORLD-86 is cited as OPEN in design, plan and JSON but is not in design_docs/world-mission.md (grep count 0; last id is D-WORLD-85). The plan says "being filed"; the controller must file it with the PR or the hold points at nothing.
NB4. Protected hardlink-write/edit arms duplicate no-protected-hardlink in what they can distinguish; World's wiring for hard links is untested by M5e (acknowledged in the plan: refusal is upstream and unconditional). The symlink arms do not assert which pattern matched (tests/** vs *_test.go); ii-d is still caught by outcome.
NB5. M5e is real-pin-only and SKIPs in CI (no pin-free companion, row 160); the plan states this plainly, so the regression guard only bites on controller/pin-move runs.

## Blocking
None.
