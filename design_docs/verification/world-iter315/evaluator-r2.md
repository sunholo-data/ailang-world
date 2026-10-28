# Evaluator report round 2: iteration 315 fix round (0e975c53)

Score 84/100, PASS. Correctness 25/30, verification 21/25, consistency 17/20, honesty 12/15, scope 9/10. No BLOCKING findings.

## Diff walk (git diff c36a4306 0e975c53)
Six files, exactly the coordinator's list. Not named in the §9 "Judge round 1" paragraph (all harmless):
- the Go header comment of `v66-phonecheck-main.go.txt` was gofmt-reflowed (part of the re-pinned hash);
- V67's banked `desktop-light` line now shows `thumb={-1 -1 -1 -1}` (old span test found grey on the 1440 px image, the run test does not; verdict unchanged);
- the 12th V67 arm (`scattered-pixels-no-thumb`) is the pixel file from MY scratch dir (`/tmp/world-iter315-eval/phn/`, sha c57b01b1... matches my file); its generator is not banked, so the arm cannot be rebuilt from the repo.
Everything else in the diff is named in §9.

## (a) Rebuild and re-run
- sha256 of `v66-phonecheck-main.go.txt` = `31576922bc56...0b4c`, equals the §4.4 pin and V72. Build rc 0.
- 19 positives: output byte-identical to my r1 run AND to banked `v66-phone-positive.txt` (`r3-accept=19 r2-accept=1`).
- My 15 r1 phone arms: only `scattered-pixels-no-thumb` flips ACCEPT -> REFUSE; the other 14 verdicts are unchanged.
- V67 generator rerun with the new checker: 13 candidates, 2 accepted (shift-2/4); with the scattered arm 14 candidates, still 2 accepted. V72 controls: accept ACCEPT, refuse REFUSE.

## (b) Attacks on the contiguous-run fix
| Arm | Result (new) | Old |
|---|---|---|
| thumb split by one off-range pixel at x=150 | REFUSE (longest run 151..333, X0 151 > 8) | ACCEPT |
| split at x=8 | REFUSE (X0 9) | ACCEPT |
| split at x=250 | ACCEPT (longest run 6..249) | ACCEPT |
| 100 px run only above the bottom 20 rows | REFUSE | REFUSE |
| grey bar x 0-300, rows 830-834, no real thumb (page-content grey) | ACCEPT | ACCEPT |
| full-width grey bar | REFUSE (X1 390) | REFUSE |
| two pieces on different rows | REFUSE | REFUSE |
| exact 100 px run (x1-x0=99) / 101 px | REFUSE / ACCEPT | same |
So the fix is correct for the stray-pixel case but introduces a fragility: a genuine thumb that has one out-of-range pixel anywhere in its first qualifying row, left of its midpoint, is now refused. The 19 real images are all clean runs, so it does not bite on evidence. The first-qualifying-row-only X0/X1 behaviour is inherited. Shape-only rule: a grey bar of page content starting at x<=8 and ending <=384 still passes (clip and dims must also hold).

## (c) Rule text
§4.4 phone paragraph, (a1)-(a4), (b), (iii') bounds and C8' text are byte-unchanged. Changed: a qualifier narrowing a prose claim ("whole lines in the two body-text gaps"), the 97/73 readings, "12 of 12" notes, C13 (now names the r5 carve-out as the disposition; a gating label), header/opener/preamble labels, the removal of the unverifiable "iteration 95" parenthetical, the re-pinned hash. No acceptance rule changed. Confirmed.

## Adjudication of r1 findings
- Divergence 1 (run vs span): CLOSED, with the new split-thumb fragility above (NON-BLOCKING, new).
- Divergence 2 ((ii) std-dev stand-in): CARRIED (listed in §9; §4.4 prose still says "the M7b gate").
- Divergence 3 (X1 exclusive): CARRIED (listed; not yet in §4.4).
- NB1 stale text: CLOSED. Opener now r5, C13 names the r5 carve-out, header says Gemini also passed the r4 round, §5 wording, §10 preamble V64-V73 plus quorum-r4, 97/73 beside every 72/96 (§4.4, §8, r3 table). I found no remaining stale live sentence.
- NB2 row 231 not in charter: STILL OPEN. §4.4 and §9 still say "is filed"; `grep -E '^231\.' world-mission.md` is empty. "Carried to row 231" has no destination until the controller files it.
- NB3 identity accept controls: CARRIED.
- NB4 P0 item (1) probe logs, no command: CARRIED.
- NB5 V65 does not isolate four clauses: CARRIED.
- NB6 whole-line slides: wording CLOSED, behaviour CARRIED.
- NB7 per-fix AC and mutation (charter guardrail): STILL OPEN. Not mentioned in §9 or carried.
- NB8 truncated V50/seed-5 hashes in §4.4: STILL OPEN. Not mentioned.

## New findings (all NON-BLOCKING)
1. Split-thumb false refusal (above).
2. The 12th V67 arm is not reproducible from banked sources and its sha line names a /tmp path of the judge's scratch dir; bank the generator (a 3-line addition to `v67-phoneneg-main.go.txt`) or reword.
3. Carry list targets row 231, which does not exist yet and is scoped as "commit the checker"; items such as P0 item (1) command text and the accept-control choice are controller-step hardening that can be done now without the row.
