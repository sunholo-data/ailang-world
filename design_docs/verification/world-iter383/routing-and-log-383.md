
## 383 — 2026-11-11 — row 201 A3 child design r2 (`w-world-approvals-by-key-a3.md`): quorum blocked twice, controller carve-out, independent audit VALID 84; D-WORLD-114 filed on its two amendments to the parent; row 201 A2 sprint 4 LANDED (PR #552, merge `88ff1565`) once the CI runner apt outage cleared on attempt 8 [PRODUCT]

**Pick:** carry-over first: land PR #552 by iteration 382's predicate. While the outage persisted, ci-health's rule for a provider-side red applied: do not fix-forward or revert, and pick routable work that needs no landing gate. That was row 201 A3's child design, which the charter's Next names after Mark's attended A2 exit run (attended, not loop work). The outage cleared on attempt 8, and #552 landed in the same fire.
- **Re-measured at pick** (`origin/dev` `27ce163e`, record 381; PR #552's branch carries record 382, read as Gate-1 input): no `a3` doc, branch, PR or worktree; `mc_quorum_gate2 w-world-approvals-by-key-a3` rc 1.
- **Ledger:** valid (100 rows before, 101 after). OPEN: D-WORLD-105; D-WORLD-114 filed this iteration.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 not routable (unchanged).
- **Orphans:** none. Iteration 382's fire ended normally; its fyi arrived 3 minutes before this fire.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #524, of 16 comments. The watermark file is absent, so all comments were read.
- Self-notices: rc 2, an instrument failure and not a verdict; no `-last-seen` file exists for 524 or 449 (same as 371–382).
- Inbox: nothing unread addressed to `mission-world` (`-inbox mission-world` empty).
- Skill drift: none. SKILL.md and all 13 resource files at the resolved symlink target match the driver clone's `origin/dev`.
- Weekly issue sweep: not owed.

**Gate 1 / 3b — PR #552 (head `0140bd7f`):** `rerun --failed` attempts 2–8 this fire (16:06–17:00Z). The host job was red at `apt-get install -y -qq bubblewrap socat ripgrep` (exit 100) on attempts 2–7, with 6–7 `Temporary failure resolving 'azure.archive.ubuntu.com'` lines each, before any Go test ran. **Attempt 8** (17:00Z): the apt step passed and every later step ran; the host job was green, so all **3/3** checks passed on `0140bd7f` (`commits/<sha>/check-runs` checks=3, all success). `ailang-code verify` and `go race` were green throughout.
- **Landed:** `gh pr merge 552 --merge --match-head-commit 0140bd7f…` → merge **`88ff1565`** at 17:11:12Z (no closing keywords in the PR text). Push CI on the merge (run 43065674075): 3/3 success (run 43065674075, 17:26Z).
- **Controls:**
  - the GitHub status API: `none` / All Systems Operational, no unresolved incidents;
  - from the rig, `azure.archive.ubuntu.com` resolves (a traffic-manager CNAME to a Sweden Central load balancer);
  - the 382 `dev` dispatch control already failed at the same step.
- **Reading:** the runner-side Ubuntu mirror is broken, which makes this a provider-side red. ci-health: re-run, record the diagnosis, do not fix-forward, do not revert.
- **Not done:** pointing the apt step at another mirror would have been a fix-forward of a provider outage; it cleared on its own about 3 hours after it was first seen.

**Design** (`w-world-approvals-by-key-a3.md`, a child of the attended parent). Resolver: `recipe claude:claude-opus-5-5 declared:provider-pin`; spawned as Agent `opus`, following the 379 precedent. The rotation file reads `claude:claude-opus-5-5`; codex was unusable at preflight (rc 75). Worktree `.wt-world-iter383-a3`, branch `design/world-iter383-row201-a3`, stacked on #552's head.
- **r0** (12.2 min, 798 lines). Split:
  - **A3a**, the web flow, in three sprints:
    - credentials schema 7 `surface_approvals`;
    - one shared judge;
    - the CLI minter `approve <id> --await`;
    - `GET /surface/requests.json`, `GET /surface/requests/{id}`, the card page and `POST /surface/approve/{id}`, with Host, Origin, `Sec-Fetch-Site` and JSON checks;
    - *Needs you*, the decision list and an opt-in notification.
  - **A3b**, the Tauri notification and Touch ID: constraints only, because no Tauri app exists.
  - The designer measured that minting is CLI-side (`session_new.go:404`) and that the daemon never keeps the credentials writer (`host/daemon/credentials.go:50-53`). **The controller re-read both, plus the absence of `desktop/`, first-party.** It also probed Chrome's POST headers and the cross-origin preflight (V17).
- **Quorum round 1 BLOCKED 3/3** ($0.385). Author Anthropic benched; seated oc-glm-5-3, oc-kimi-k3, gpt6-1-sol; absent none. The controller checked each objection against r0, and all three were real:
  - R3 blocked SSHSIG behind a tailnet origin;
  - the `crypto.subtle` premise was cited as measured but was not;
  - the card omitted signed fields.
- **r1** (designer resumed by SendMessage, 4.7 min, 928 lines):
  - a surface origin set with `serve --surface-origin` (RT13, CF1);
  - V23–V26 measured in Chrome 157 at `http://localhost` under `surfaceCSP` (digest, the 90-byte message and `MessageHash` byte-equal to Go on two golden requests; `Notification` behaviour);
  - a `FIELDS` table rendering every `MintRequest` key path, with PG3.
- **Quorum round 2 BLOCKED 3/3** ($0.426, same seats). Every objection carried a reviewer-authored `proposed_fix` and none disputed the direction:
  - glm: the missing-path refusal assumes no `omitempty`;
  - kimi: two amendments to the attended parent have no record;
  - sol: the card's GET and POST have no network deadline.
- **r2 by controller narrow-refinement carve-out:** the fixes applied in the reviewers' own text.
  - glm's branch (a) holds, measured: `omitempty` count 0 (V28).
  - V29 (`env.after`) and V30 (`DecidedMintRequests("")`) measured.
  - kimi's *Amendment posture*, V31 (the parent's R1–R5 fix no mint location or A3 vehicle) and V32 (base-drift trigger), plus the §10 mirror and the §15 checklist.
  - sol's 10-second AbortController paragraph and XA5.
- **Independent audit** (evaluator sonnet via the Agent tool; generator ≠ judge): **CARVE-OUT VALID, 84/100, zero blocking.**
  - It re-ran V28–V30, read V31 and D-WORLD-88 R1 (the deferral sequences A3's native add-ons and does not contradict R1), and judged kimi's objection a missing record, not a direction dispute.
  - Its consistency findings were applied (`fad08b71`): PG3b/XA5 added to §9's gates and the `-run` pattern (XA5 would otherwise not have run), PG3b renamed, V29/V30 cited.
  - Noted, not changed: glm's third golden case is a test obligation rather than a measured row, and the D-WORLD-114 sentence is a controller addition.
- **Designer diet:** one doc, r0 plus one protocol-mandated revision. Held; r2 cost no designer run.

**Not done this fire, on purpose: the A3a sprint-1 planner.** Sprint 1 builds the CLI minter, which is amendment (b), the thing D-WORLD-114 lets Mark veto, and the remaining slot after a design fire is not enough for a planner, executor and judge round (382's planner alone took 218 min). So the default (A) waits one report cycle before any A3a code.

**Routing evidence:** base=27ce163eb1ae8af55151ea6548591c15a3126835@2026-11-11T16:06:23Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: `claude:claude-opus-5-5` via Agent `opus` (r0 197,650 tok; r1 resume 251,486 tok as reported); r2 by controller carve-out.
- Quorum: `ailang design-quorum --author claude:claude-opus-5-5`, rounds 1+2 $0.811 metered (artifacts `mission-quorum-world/w-world-approvals-by-key-a3-2026-11-11T16-31-03Z.json`, `…T16-47-14Z.json`).
- Planner and executor: not spawned (design-only fire; see above).
- Evaluator: `sonnet` via Agent (44,668 tok), auditing the r2 carve-out.
- The driver noted codex unusable at preflight (probe rc 75).
- Metered total: $0.81 (quorum only).

**Progress:** 1.0 clauses 1–7 MET; the goal is unmoved. A2 is code-complete on `dev`: an operator can enrol a passkey from the printed link and approve a pending request with it, with no World terminal prompt. Mark's real-passkey run (A2 §10, QUICKSTART §16) is now runnable. Row 201 also has a reviewed design for its last human-surface milestone: an operator approves a mint from the surface's *Needs you* card with a passkey or SSH signature.

**Ruled out:**
- **Merging #552 over the red (attempts 2–7):** an item is not LANDED until remote CI passes; it merged only after attempt 8 went 3/3.
- **Fix-forward of the apt step:** ci-health forbids it for a provider-side red.
- **Parking A3 as `needs-human-review` after round 2:**
  - every remaining objection had a concrete fix and none disputed the direction, so the carve-out applied;
  - an independent judge confirmed it;
  - the human question that remains (the two amendments) is asked as D-WORLD-114 with a default, not as a park.
- **Codex for the designer:** its lane was unusable at preflight (rc 75), and the resolver named opus.

**Retro:**
- **A design premise about browser behaviour gets probed at design time, not cited.** r0 measured the POST headers but asserted `crypto.subtle` byte-equality, and round 1 caught it. The designer measured it in five minutes once asked. Next design directive: list every browser-capability premise and require a V-row per premise.
- **Amending an attended parent needs a record in r0, not a disclosure.** r0 labelled the departure *DEPARTED-FROM-PARENT*. Round 1 let it pass; round 2 blocked on it. Next directive: any departure from an attended parent gets an *Amendment posture* paragraph and a ledger ask in r0.
- **The outage lasted about 3 hours** (first seen before 14:53Z, still failing at 16:54Z, clear at 17:00Z). Re-running and doing other work while waiting cost one design fire and no lost work; the same handling applies next time.

**Next:**
1. If this record's own PR is still open, land it first (docs only: the A3 design, record 383).
2. Then the A3a sprint-1 planner from `w-world-approvals-by-key-a3.md` r2 (no re-design). D-WORLD-114 default A applies from iteration 384 if unanswered. Run V32's base-drift check first.
3. Mark's A2 exit run (§10, QUICKSTART §16) is runnable and stays attended.
