MISSION-ROLE: evaluator
You are the INDEPENDENT record evaluator for mission-world iteration 218 — a bookkeeping-only block report (no product work exists in this iteration). Generator ≠ judge: the controller (pi:ollama/glm-5.3:cloud) authored this record; you are a different model in a fresh context. Judge ONLY the iteration-218 RECORD (docs), never any product code.

Your working directory is a read-only evaluation worktree checked out at commit ba9fdc6611b9b5e3d9ee53eaf39c820795a1ed85 (the record candidate, immutable). Do NOT edit tracked files, do NOT commit, do NOT push. Run any read-only commands you need (git, grep, scripts) to verify the record's claims FIRST-PARTY — a claim you did not re-derive is a claim you cannot bank.

THE RECORD'S CLAIMS TO VERIFY (verify each by running commands yourself):
1. Block report completeness: `design_docs/world-mission-log.md` entry `## 218` states the clause map (clauses 1/2/3/7 MET; 4/5/6 UNMET), names each blocked row with its exact blocker (row 108 → D-WORLD-49 OPEN, default B; row 114 → follows 108 per D-WORLD-46 = A; row 93 → needs 108 landed), cites the charter's standing rule (d) (BLOCKED MEANS STOP, NOT SIDE WORK), and picks NOTHING.
2. Cross-check against the live charter: run `scripts/mission_decisions.sh --check --file design_docs/world-mission.md` (expect: valid, 36 rows) and `scripts/mission_decisions.sh --open --file design_docs/world-mission.md` (expect: exactly D-WORLD-49 OPEN). D-WORLD-46/45/47 rows are RESOLVED with the answers the record cites.
3. STATUS rotation: `/usr/bin/grep -c "^## STATUS 2026" design_docs/world-mission.md` == 3; the 218 stamp is present; `(iteration 215)` moved to `design_docs/world-mission-status-archive.md` (grep it there: ≥ 1).
4. Index: `design_docs/world-mission-index.md` has a `| 218 |` row dated 2026-10-02.
5. Dashboard: `design_docs/world-mission-dashboard.md` is ≤ 40 lines and reflects the blocked state (D-WORLD-49 default B, no fourth round, goal unmoved).
6. No overstatement: the record must NOT claim product acceptance, must NOT claim PR #173 or #166 were merged or reviewed as product, must NOT claim a fourth row-108 round ran, must NOT claim any clause moved.
7. Diff scope: `git show --stat ba9fdc6611b9b5e3d9ee53eaf39c820795a1ed85` touches ONLY doc files (dashboard, log, index, charter, status-archive). ANY production code in the record diff is a HARD FAIL.
8. Commit-message scan: `git log --format=%B -1 ba9fdc6611b9b5e3d9ee53eaf39c820795a1ed85` contains no closing keyword followed by #N (fix/fixes/fixed/close/closes/closed/resolve/resolves/resolved + #number).

SCORING — report a score /100 and a verdict PASS or FAIL:
- HARD FAIL (verdict FAIL, score < 60): a false claim in the block report; any production-code file in the record diff; the 215 stamp lost (absent from BOTH charter and archive); the record claims product acceptance or a row-108 judge verdict.
- Deduct for: a missing fixed section (the log entry should carry Kind, Pick and why, Gate 0/1, Premise, Designer, Planner, Executor, Independent evaluator, Controller gates, Gate 3b, Routing evidence, Record, Ruled out, Retro, Progress, Next); a blocker named without its decision id; a clause map disagreeing with the charter's bar; an index row or dashboard missing or stale.
- PASS = score ≥ 60, zero HARD FAILs, zero BLOCKING findings.

DELIVERABLE: write your report to `record-eval-218-r1.md` IN THIS WORKTREE's root (untracked file — not a repo edit) with: first line `VERDICT: PASS` or `VERDICT: FAIL`, then `SCORE: N/100`, then findings as bullets each marked BLOCKING or NON-BLOCKING, then the commands you ran with their key outputs, then your model self-identification. FINALLY, print as your last line of output exactly: `RECORD-VERDICT: <PASS|FAIL> <score>/100 <one-line summary>`. If a prerequisite fails (missing commit, unreadable charter), report the exact error and stop — that is a transport failure, never a PASS or FAIL.
