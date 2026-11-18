# w-world-address-bar-verbs-v4 — start and finish an episode from the World app

- **Status** Planned r0. Designed in iteration 406, 2026-11-18, queue row 243, parent phase V4.
- **Parent** [w-world-address-bar-verbs](w-world-address-bar-verbs.md), with [V3](w-world-address-bar-verbs-v3.md) and [its S6 handoff](w-world-address-bar-verbs-v3s6-sprint-plan.md#appendix-c--after-this-sprint-s6s-close-the-v3-exit).
- **Authority** D-WORLD-89 R1–R5, D-WORLD-94, D-WORLD-107 = A, D-WORLD-109 = A, D-WORLD-115 = A, D-WORLD-118 = A (W02). Decisions requested in §10.
- **Designer lane** `codex:gpt-6.2-sol`.
- **Measurement base** `41331b4d880e4ea1a4c19a2f1e65a74204580bc4`, branch `design/world-iter406-row243-v4`, origin/dev base as assigned (W01).
- **Instruments** Git read commands, rg, nl, sed, wc, read-only macOS and tool-discovery probes (§12).
- **Exit** Mark starts one confined headless episode, reads its report, and finishes it to a PR in World, approving each step with his enrolled key. The real run is attended.

## 0. Scope and ordering

Start with create-and-launch, then implement finish. A confined browser launch establishes the user identity, signed action, durable claim and restart boundaries that finish also needs. Finish adds local and remote effects, a second prompt, and recovery from uncertain GitHub responses.

- Phone workflow draft followed for order, pinned review, durable outcomes and one headless Sol. Its mobile page redesign and merge action are separate work.
- Human-agent surface draft followed for a durable link, short progress/outcome and visible blockers. Multi-user handoff and the episode-page redesign are separate work.
- V5 owns attended exit-run readiness. V4 supplies the executable flow, its fake runbook and the attended hook. Source briefs and parent scope in W30–W31.

### 0.1 Amendments

| ID | Document and location | Proposed amendment and reason |
|---|---|---|
| AM1 | Parent §0.3 V4 and §12 Q4 | User LaunchAgent runner, create-and-launch first, calibrated sprint plan in §8. |
| AM2 | Parent §2.1, §10 steps 1–5 | Surface start uses inline bounded task text. Finish first authorizes local preparation, then uses the real post-commit confirmation. Headless is the exit path. |
| AM3 | V3 §0.4, §2.7 F1, RV4 and S6 Appendix C | Only a verified runner invocation can cross F1 without a terminal. The public CLI retains the F1 behavior, including with the V3 flag. |
| AM4 | Row 201 parent §2.4 and A3 §§2.1, 2.5–2.6 | Add a separate operator-start proof and local consumer alongside the A3 CLI minter. A3 card, request schema and judge stay dedicated to mint requests. |
| AM5 | Parent RK5 and row 219 service design | Require and measure runner-specific PATH, GitHub and Git authentication. Extend typed service lifecycle for an opt-in companion. Do not ask the serving job to inherit interactive credentials. |
| AM6 | V3 §2.6 and §2.7 | A V3 card linked to a runner action has a liveness/generation guard and a durable waiting owner. The finish owner still supplies and verifies the exact payload, consistent with A4’s reserved scope rationale. |
| AM7 | QUICKSTART service, episode start and §20 (W27) | Add service companion setup, surface start and two-step surface finish. Existing terminal entry stays available. Mark’s end-to-end walkthrough is handed to V5. |

## 1. Measurements and premises

| ID | Parent or handoff premise | HEAD measurement and design consequence |
|---|---|---|
| P1 | Parent Q4 says the daemon spawns nothing | No direct process launch was found in the daemon package, but the daemon brokers workspace execution on the host. Q4's distinction is between brokered tools and operator process supervision (W03). |
| P2 | Parent §0.3 suggests a small wrapper and interactive hand-off | The service supplies one GUI LaunchAgent for `serve`; the plugin's New World agent action does a typed mint and can select an unconfined fallback. Neither is a durable operator runner (W04, W08). |
| P3 | RK5 describes the risk as missing `gh` auth in a launchd daemon | The installed job on the Studio is a user LaunchAgent, and its actual PATH cannot find `gh`. An explicit Homebrew PATH finds it, but its clean-env auth probe did not establish usable credentials. Git push auth is a separate preflight (W05–W07, W20). |
| P4 | Parent §10 expects the start request prompt with a launch preset | `--request` refuses launch presets, launch and headless. A3's signed mint and the preset mint with a start record enter different functions and transactions (W09–W12). |
| P5 | V3 §0.4 / Appendix C call for running the CLI and relaxing F1 in surface mode | F1 is still real, now at lines 89–94. Dropping F1 or adding the flag alone would let the application run checks and commit before the first key, and would allow reuse of a finish approval from a recorded dispatch failure. Use a private checked entry and a separate preparation proof (W16–W18). |
| P6 | Appendix C describes the V3 card and waiter as the finish hand-off | The approval row records a confirmation, and the waiter rechecks in process. There is no worker claim, checkpoint for the in-progress commit or dispatch, or exact post-timeout PR identity check in that path. Preserve its signature law, add durable ownership and recovery (W15, W17, W19–W21). |
| P7 | Parent §10 step 2 offers herdr or row 206 as launch destinations | The installed CLI has pane control, but the checked-in plugin's action takes a terminal and assumes herdr context. The row 206 charter entry supplies roadmap placement, not a callable adapter for V4 (W02, W08, W22). |
| P8 | V3 says 10 minutes / at most 610 s and the bar answers `terminal` | Both are substantiated, with distinct 600 s approval and 10 s observation grace. The public CLI calls F1 before effect I/O, and the bar is a GET classification route. A durable waiter and a guarded POST are separate additions (W15, W16, W18, W23). |
| P9 | Appendix C leaves the choice of per-check approval and notification open | Unsandboxed checks are refused on fresh V3 surface finish before the check loop. Needs you already discovers V3 requests, but the notification projection selects mint/decision items. V4 preserves those choices and checks all runner finish modes (W16, W23). |
| P10 | Parent §0.3 estimates 1–2 sprints / ≈2,000 lines | Separate admission, credential transaction, process ownership, effect reconciliation and tests are needed. V3's S5 alone added 1,398 changed lines in the inherited budget, with F1 kept (W24). |
| P11 | Parent §10 names `--task <file>` for the first bar run | The CLI reads a host file. A browser-supplied pathname would give the page an operator-file read primitive. The surface supplies inline text, and the rig derives and pins its own task copy (W09, W12, W13). |
| P12 | The first slice is one Sol launch | The code selects a host agent family and a preset model. Sol is an operator preset to validate on the Studio, with no assumed hard-coded family/model mapping (W09, W13, W25). |
| P13 | V3 handoff treats dispatch and outcome recording as one uninterrupted command | The finish record route rejects expired scopes, and independently sets the used and episode heads. A delayed receipt after confirmed remote success needs a scoped observed-writeback path (W28). |

## 2. Q4 — runner

### 2.1 Options

| Option | Identity, environment and lifecycle | Learns of work | Reuse and blast radius |
|---|---|---|---|
| A. Daemon spawns | Same UID as today's GUI job. The fixed PATH and demonstrated auth gap would still need explicit configuration. Service start/restart/upgrade would own browser, world writer and worker lifetime (W04–W07). | In-process after durable confirmation; still needs a claim and effect journal. | Reuse daemon approvals and launch seams (W03, W14, W15). A malformed dispatch or hung child would share the World writer's process, shutdown and environment. |
| **B. Separate user LaunchAgent** | Opt-in through attended `service`. Same UID, independent job, explicit trusted HOME/PATH and per-tool configuration. Launchd restarts it, and paired upgrades drain and fail closed on mismatch. | Poll actions in the local credentials database; heartbeat, claims and result checkpoints use bounded transactions. Browser sees the daemon's read model. | Reuse service plist/controller, status/doctor, credential locking, launch and finish seams (W04, W10, W14–W15, W17, W20). Its UID and GitHub credentials still have Mark's reach. Process separation limits accidental coupling to the serving daemon. |
| C. herdr plugin owns runner | User UID with herdr's service and pane environment, which requires measurement of HOME, PATH, Git and `gh`. Plugin/server startup or restart would own execution and availability, and CLI/plugin compatibility would join upgrade. | Would need a persistent plugin watcher and the same verified durable claims. Popup and pane tracking supply UI context only (W08, W22). | Reuse pane creation/run and metadata. The existing popup mints interactively and offers fallback; widening that entry point would couple headless work to ambient pane configuration (W08). |
| D. PIN a manual user process | Same terminal login credentials at start; operator manually supervises and restarts it, and binary skew stays possible. | Could poll the same table. | Useful for a synthetic smoke fixture. A job that requires a terminal to keep the product available misses the requested service behavior. |

Choose **B**. Studio evidence establishes that the GUI service installation model is available, and the inherited PATH is inadequate. Adding Homebrew to PATH still leaves authentication unverified. Gate finish on E2 and launch on E3. Using B with a separate job preserves Q4's process boundary and makes herdr optional for the headless exit.

### 2.2 Lifecycle and configuration

- NEW `service runner install|start|stop|status|cancel|uninstall` and NEW internal `operator-runner` mode, with no command string or verb argv accepted at its entry. Install is optional, typed and user-scoped. Start uses the previously approved immutable configuration. Stop, replacement and uninstall are attended.
- NEW separate label, 0700 private configuration and spool, 0600 state files and redacted rotation. Reuse `serviceEnv.confirm`, `renderServicePlist`, `systemServiceController` and the build/status precedents (W04). Record the canonical daemon store pair, login UID, allowed preset names, absolute executable paths and hashes, owner-validated HOME, tools dirs and allowed backend. Refuse root, a writable-by-others configuration component, a symlink or changed executable, mismatched stores/build and non-loopback daemon.
- Service constructs PATH from approved directories and the system floor, HOME from the approved login directory, and child-specific allowlists. Set Git/SSH/GitHub to non-interactive. The Git seam strips ambient GIT_* (W14), so enforce the reviewed SSH command and HTTPS credential helper policy in trusted Git configuration/arguments rather than ambient shell overrides. If an SSH socket is configured, require a live user-owned path and use the allowed socket only for Git. HTTPS helpers, Keychain and agent login are measured under the runner environment. Do not persist source tokens in a plist, runner config or logs. Agent provider auth uses its approved, existing user configuration or an explicitly reviewed credential source. Browser input supplies no environment variables, SSH config, process name or command.
- Daemon and runner share a pinned build and store paths. On upgrade, refuse new work, drain active children or preserve their uncertain checkpoints, then replace/restart both jobs through the attended service seam. On failure, start only a compatible pair. A pre-v10 binary refuses v10 (W24). Roll back by stopping the companion and running a schema-compatible binary, with a future forward migration if the table must be retired.
- NEW heartbeat record in the same private store carries protocol/build, lease generation, readiness per verb and a last-seen time. Poll every 2 s with bounded reads; after 10 s stale, refuse new surface actions. Heartbeat is liveness only. A missed heartbeat cannot release execution ownership or prove a child died. Hold one kernel-backed runner lock per served credentials store and one effect-exclusion lock per episode; terminal start/launch/finish also honor the same per-episode lock for managed episodes and release it while waiting only if the durable phase prohibits conflicting writes.
- Precompute a small NEW read-only catalogue from configured presets and completed/eligible episodes at boot and on a fixed timer, independent of HTTP requests. Use bounded existing read/health and Git seams. Publish only a redacted, complete, short-lived immutable generation through the credentials store. A browser command never triggers a pre-key subprocess. Snapshot age, claimed runner exclusivity, credential status and software identities are checked again at execution. Serve no arbitrary paths or live shell environment.
- `doctor` and service status gain a NEW companion diagnostic for absent, stale, schema/build conflict, version ready and auth blockers. Per-verb readiness separates agent runtime from GitHub. Output names a remedy and reports no token, signed artifact or task.

## 3. Durable actions and cryptographic admission

### 3.1 Persistence, claim and proof

Recommend a **new credentials schema 10** and NEW `operator_actions` table. V3's `verb_requests` is single-use confirmation storage with its own episode supersession, pending cap and approval outcome. Extending it would entangle two lifetimes and allow a terminal finish to supersede a start job (W15). An independent queue file avoids a credentials migration but cannot atomically couple start mint, proof consumption and its durable receipt. The credentials DB is the owner. NEW companion tables `operator_runner_state` (single heartbeat and catalogue generation) and `operator_action_steps` (bounded per-action effect receipts) are in the same migration. New symbols, states, fields, routes, schemas, limits and tests in this design are proposed unless cited.

- NEW action id `oa-` plus 32 lowercase hex digits, challenge of 32 random bytes, signed request expiry 600 s, exact canonical resolved-intent bytes and digest, phase generation, owner nonce, linked V3 confirmation id if any, enrolled approver identity and artifact, checkpoints and redacted outcome. NEW browser operation id is a random 128-bit value retained across retries. Use a review digest over semantic inputs before assigning an action id or challenge and a separate signature digest over the final envelope. Unique kind/episode/review digest for active actions, unique operation id with exact review-digest equality, and one active action per episode. Different bytes using an existing operation id conflict. Identical independent taps on the same frozen intent link the same active or completed action. Expired and cancelled intents require an explicit fresh review; the dedup keys and permanent effect receipts survive pruning and are never repurposed.
- Signed schema NEW `world/operator-action/v1`; message NEW `world-operator-action/v1`, NUL, 32-byte canonical digest, 32-byte challenge, big-endian 8-byte expiry; NEW SSHSIG namespace `world-operator-action@sunholo`. Canonical action includes schema, action id, kind, episode and served store identity, command words and admitted options, frozen resolution, and displayed review bytes or their hash reconstructed by both sides. Closed typed object, deterministic field order, reject duplicate/unknown keys, trailing JSON, invalid UTF-8, NUL and over-limit input before canonicalization. Compare using canonical bytes. NEW limit 128 KiB resolved action and a separate 8 KiB UTF-8 task cap.
- Reuse existing key lookup, SSHSIG, WebAuthn assertion primitives and the card ceremony structure (`host/authority/judge_verb.go:30,44,63,82,107`; `host/daemon/surface_confirm.go:97–207`, W15, W26). NEW typed verified-action proof can originate only in these judges. Preserve UV, exact RP/origin, namespace, enrollment/revocation checks and monotonic passkey use. Signature is over action identity, kind and the exact review, never over the raw page payload alone.
- Daemon verifies, then atomically changes requested → confirmed and advances approver use/counter with CAS on id, digest, generation and unexpired state. A duplicate confirmed gesture returns the existing action link using a fresh read, with zero second counter update and zero second claim. Runner derives its own bytes from frozen values and revalidation, recomputes the signed message, judges the stored artifact against the current key and compares the approver identity before any effect. Recheck expiry and revocation again at each new effect boundary. If checks exhaust the preparation authorization before commit, stop with their receipts. Reauthorization can approve a successor with identical worktree and check inputs and explicitly carried receipts. No second check run is implicit. The transaction that mints start rechecks proof, approver and owner with a fresh clock. A counter already advanced at confirmation is the expected observation, never advanced a second time.
- Claim by transactional compare-and-swap from confirmed to claimed, with matching current runner generation, same digest, fresh approval and no owned active action for the episode; exactly one affected row is required. All subsequent writes CAS the action/step version and owner. Two runners are also excluded by the kernel lock. Recovery takes that lock and verifies that any recorded child or reservation is quiescent before changing owner; a lost lock or lease never authorizes a second child. A live matching child is observed, an unprovable child is uncertain.
- Daemon routes create and confirm requests, read stable status and expire untouched work logically. Runner claims, verifies, dispatches, renews liveness, checkpoints and reconciles. Browser tab close is a read disconnect. No cancellation on disconnect and no unsigned cancel/deny. NEW attended service cancel marks an unstarted action cancelled by CAS; an in-flight effect first goes through bounded stop/reconciliation. Cancellation never deletes a worktree, commit or PR. Once a prep/launch dispatch is accepted, its signed expiry limits **new** effects, not the execution of an already-authorized bounded child or observation/writeback of its outcome.

### 3.2 States and restart

| Durable state | Writer / comparison | Next state or observation |
|---|---|---|
| Requested | Daemon insert only after grammar, capacity, snapshot and same-origin guards | Review URL. No effect. Expired at the signature deadline; only explicit retry creates another request. |
| Confirmed | Daemon signed CAS and passkey use transaction | One live runner can claim; cancellation or expiry may win before a claim. |
| Claimed | Runner CAS after acquiring locks, once per generation | Verify bytes, key, time and inputs before moving to running. Restart can return an untouched claim to confirmed only after proving no effect was reserved. |
| Running | Runner owns action and per-effect reservation | A reserved step permits one attempted effect. Effects with positive receipts advance. A known refusal with no side effect fails; missing or conflicting evidence enters uncertain. |
| Awaiting PR confirmation | Runner has prep checkpoint and links one V3 row | V3 card records key proof; runner recovers checkpoint and proves own V3 bytes, or expires/cancels. Durable prep is not erased. |
| Succeeded | Runner stores authoritative receipt and stable user link | No automatic new execution, including after reconnect. Launch action succeeds on observed start; completion updates the episode report/status separately. |
| Failed | Runner records confirmed no effect for the failed step, or an observed child exit with its durable effects retained | Display the completed steps and explicit recovery. No automatic re-run on page refresh. |
| Uncertain | Owner lost, timeout or conflicting evidence after reservation | Read-only reconciliation can recover a matching receipt. No blind replay. Preserve the lockout across restarts when evidence is insufficient, with operator handoff. |
| Expired / cancelled | Logical expiry or attended CAS with runner boundary checks | Preserve effects and status. Allow read-only reconciliation of work already reserved. Fresh approval required to take any permitted new effect. |

Daemon restart preserves the database. An unavailable daemon pauses new work and authoritative report reads; a child already admitted may finish and its receipt is stored locally, then reconciled when the daemon returns. Local checkpoints use finite contexts, independent of the page, and release SQLite writers before network/process waits. Pending caps are NEW 16 total and one active action per episode across starts and finishes. Terminal statuses may be pruned after 30 days if operation-id and effect-identity tombstones and uncertain reservations are retained. Daemon never turns an uncertain action into a new claim.

### 3.3 Side-effect identities and replay rules

| Effect | Intent and reservation before issuing | Reconcile before further writes |
|---|---|---|
| Worktree / branch | Episode, canonical repo, pinned base commit, expected new branch and unique path. Reserve under the episode lock, demand initial absence, use `worktree.Add` (`host/worktree/worktree.go:63`, W11/W20). | Validate Git registration, physical path, branch, repo identity and starting SHA from the recorded reservation. Accept a matching worktree only with exclusive provenance. Unattributable or partial creation is uncertain; never remove someone else's tree. |
| Mint and start record | NEW atomic transaction keyed by action id and episode. Proof and fresh clock, claim generation, unique start record, repeated capability equality, session issuance and credential id receipt. Reuse the session/start insertion path (W11). | Query atomic receipt, start and credential. If all absent after an established rollback and proof remains valid, retry the mint in-place. If committed, never mint again for this action. Verify token file hash against credential and task/MCP files against pinned digests. If token output was lost, revoke the orphan and stop with a recoverable outcome; no plaintext token in the journal. No launch with an expired/revoked credential. |
| Task, token, MCP config | Use private staging / exclusive no-replace install to canonical recorded paths. Reuse `writeTokenFile` and `writeMCPConfig` (W11–W12). | Match ownership, mode and exact hashes or approved metadata. No overwrites. Revoke only the credential this action minted if files were unrecoverable before spawn; retain the worktree and start evidence. |
| Headless or interactive launch | Reserve one immutable run id per action. NEW launch receipt records executable identity, process birth identity, run id and observations. Reservation is durable **before** invoking the launcher. Reuse the distinct launcher lock and child supervisor (W13–W14). Existing stale PID lock files need positive quiescence and identity evidence before reconciliation can remove a lock this action owns (W32). | Positive start from the supervisor or its matching receipt adopts the run. Supervisor proof of non-spawn may release explicit retry with a fresh approval. A crash between reservation and receipt is uncertain if a live or prior launch cannot be ruled out, even with no PID. Never launch twice to find out. |
| Check / outcomes | Use approved sandboxed checks, order and input digests; reserve each check with an action/check id. Reuse finish check and world outcome seams (W16). | Use existing content-addressed outcome evidence if it matches the checkpoint. A lost process/partial output with possible effects blocks replay until read-only reconciliation proves outcome. Never rerun an unsandboxed check. |
| Local Git commit | Bind prior HEAD, exact approved input snapshot, staged tree identity, ordered paths, message/trailer and operator committer. Reserve before the bounded commit (W16, W20). | Inspect HEAD, parent, tree, message/trailer and accepted committer; adopt only the exact matching descendant under the lock. Proven no mutation and unchanged input may explicitly resume. Diverged, partially staged or unprovable output is uncertain. Never reset on recovery. |
| Finish body and PR prompt | Freeze exact body bytes and digest in private storage. Persist commit, scope, argv, prompt bytes, id, challenge and expiry **before** putting a V3 confirmation row. | After restart, regenerate only from durable canonical bytes/check results. Match V3 row byte-for-byte against its reserved id and validate the key/decision. An absent row may be inserted with that same id if still fresh and no competitor; superseded or expired status requires a new review. |
| Git push | Reserve approved remote URL, exact SHA, full branch ref and prerequisite remote-head observation. Reuse non-force `worktree.Push` (W20). | Read exact remote ref. Matching SHA is success. A different or missing ref after an uncertain response is no proof of non-execution; require positive no-send evidence or stop for attended reconciliation. Never force or use a different remote. |
| Label creation | Reserve repository and label only if the V3 scope explicitly allows it and live lookup still requires it (W17, W20). | Exact repo/name lookup confirms it. A conflict with different repository/label data blocks. After a lost write response, missing label leaves uncertainty absent positive no-send evidence. |
| GitHub PR | Reserve repository, exact head repo/ref/OID, base ref, title/body digests, label and draft flag. Use V3 evidence anchor plus a typed creation call (W17, W20). | Query by exact head repository/ref and base, across open/closed/merged states. Verify head OID, title/body, draft and label, and a unique PR identity. Adopt one verified matching outcome; show mismatches or multiple candidates for attended review. Empty/stale/failed reads after an uncertain write never allow automatic create. |
| Finish record, detach, revoke | Reserve record keyed by episode, commit and approval ref. V3 record operation uses operator proof; `finishRevoke` names episode credentials (W17, W21). Detach exact finished worktree only. | Read back record and exact PR/approval fields before write replay. On conflict read resolved record; never create again. Observe HEAD detached at exact SHA, and read each relevant credential's revoked state. Retry cleanup only for that pinned worktree and credential set. Resumed credentials added later require a separate operation. |

## 4. Operator paths

### 4.1 Bar and card

The bar uses the shared `verbs.Split`, `Lookup` and `ParseInterleaved` seams (W23). Add a NEW pure admitted subset joined to the CLI's flag definitions with parity tests. Supported starts are `episode start <ep> --preset <name> --launch --headless` and, under D-WORLD-119 A when E4 passes, `episode start <ep> --preset <name> --launch` for a pinned interactive preset. Explicit episode is required. Episode IDs use the existing 1–64 ASCII grammar and preset names its existing 51-character cap (W29). Action IDs use the NEW exact lowercase hexadecimal grammar in §3.1. Task entry is a bounded field in the review panel. The bar never reads a caller-supplied file. Supported finish is `episode finish <ep> --preset <name>`, optionally with the existing `--confirm-on-surface` spelling and, only under D-WORLD-120 A, `--ready`. Numeric, path, grant, remote, title, check waiver, dry-run and interactive overrides are outside this initial surface subset; show the CLI usage/refusal and keep the input. Enforce single occurrences even when equivalent spellings or interleaving are accepted. A launch-mode boolean must be true when explicit. Reject explicit false, empty preset/episode and padded spellings. Reuse the trusted per-episode local defaults to build canonical arguments. A page cannot type `y` to approve.

- NEW `GET /surface/verb.json` outcome `review` in the existing route (`host/daemon/daemon.go:1154`, W23) selects only the closed, supported tuple and returns short-lived reference to an available catalogue generation. Current report, why, address and approve behavior keep their existing routes (W23). When the companion is absent or incompatible, return `terminal` through the current panel pattern with the precise reason and the compatible attended CLI instruction (W23).
- NEW `POST /surface/operator-actions` creates an inert request, with strict JSON, operation id, reviewed tuple, optional bounded task and catalogue generation. No subprocess. NEW `GET /surface/operator-actions/{id}` supplies the stable HTML card and NEW `GET /surface/operator-actions/{id}/status.json` its projection. NEW `POST /surface/operator-actions/{id}/confirm` accepts only exact key ceremony data after guard checks, then returns the action link and status. Route components are whole anchored typed ids; no redirects from request data.
- Two explicit steps in the page. Enter opens the approved command's review panel. The request opens the signed review card. It shows the exact canonical human review and the current pending/running/blocked/outcome state. Browser refresh reads status and restores focus/links. Action status includes timestamps, episode/report destination, and the V3 confirm card or validated GitHub PR URL as applicable. Strict same-origin path and exact id allowlists. A GitHub result link requires parsed HTTPS and the pinned repo and PR number, no credentials/query/fragment. Keep task and proof data out of routine lists and logs.
- The bar refuses `world-publish`, `approver enroll|revoke`, `session mint|revoke` and `service`. Existing non-operator CLI commands and unrelated addresses retain their classifications (W23, W29). `world-publish` is a separate binary and falls outside those CLI word sets at HEAD. Add an explicit V4 refusal before navigation for it. No generic POST run route and no raw argv-to-shell dispatch.

### 4.2 Start, one headless Sol

Option 1 is drive the existing A3 mint request and run CLI launch afterwards. That request signs base mint fields, and today's provision path chooses request consume or a start-record transaction, separately (W10–W12). Option 2 is extend A3's schema and all compatible clients. Choose option 3, a separate start action proof using A3/V3 ceremony and authority primitives, plus the atomic start receipt. This avoids accepting an operator-start artifact via any A3 CLI approval route.

1. An attended-installed allowed preset supplies repository, remote, base ref/OID, fixed new workspace and branch, preset digest, interpreter/profile identity, task text and digest, resolved agent/model, launch mode, cap, effective grants/checks/budgets/TTL and visible committer. Render the preset start prompt along the existing fields (`episode_preset.go:240–268`, W12) plus task text and launch effect. Catalogue reads only approved presets. User inputs episode id and task text, with no repository pathname or provider tokens accepted.
2. After the key gesture, factor the effect-free portion of `runPresetEpisodeStart` into a shared resolution seam (W12). Recheck every pinned input, daemon properties and session paths before mint. A changed preset/base/profile/tool/task causes a stale response and fresh review, even when it would have lowered authority. Run the attested fixed plan, never re-read mutable user input as the action to execute.
3. Use `provisionEpisode`'s bounded worktree and file mechanisms and the existing start/session insertion parts (W11). NEW verified start proof type, held inside the runner, is a separate authorization arm from A3 request consume. Within the tentative mint transaction compare actual grants/binding/TTL/start/digest against the proof and compare-and-swap mint receipt. New storage proof can't be constructed by setting `r.start`, by an internal flag, or by supplying a plain digest. Adopt credential and file receipts per §3.3. Existing A3 typed fallback and request shapes remain available at the terminal (W09–W11).
4. Reuse `launchEpisode`, the pinned task copy and the headless `agentrun.Run` path (W13–W14). Tie one supervised process to the action/run id. Use its bounded wall-time and scrubbed environment. After mint, validate matching current session, positive TTL, action proof freshness, same pinned profile and binary before the durable spawn reservation. If the runner dies after mint but before reserve, recover one valid token and the frozen start, with no second mint. If expired during the gap, stop and show the credential expiry and attended recovery. If reserved before death, reconcile per §3.3.
5. First slice ends when fake Chrome can approve one reviewed headless start, close the page and find one durable start, one worktree, one credential and one observed confined launch. The product demonstration uses Mark's Sol preset only in the attended run. Other providers, provisioning a remote repo and the mobile form are outside this slice.

### 4.3 Finish and F1

Options are broad tty removal, a service pseudo-tty, and a narrow internal capability. Choose the narrow capability. A pseudo-tty satisfies a transport predicate while adding no key proof, and broad removal would allow V3's pre-approval local effects from CI (W16–W18).

1. NEW surface prep action signs the reviewed source report/base, episode, preset and exact worktree input manifest: expected branch/HEAD, pinned remote and target, changed path digests including untracked file contents, checks, check provenance, protected/unattributed paths, planned draft readiness and local commit intent. Read-only preview comes from the runner's periodic snapshot. Stale manifest or active agent causes a refusal. Reject symlinks or special files in dirty inputs for this slice and cap preview costs. Pin or filter Git hooks, helpers, filters and config to the approved execution boundary so a staged path cannot pick a program from browser or episode data. Detect unexpected check-generated changes before signing the PR card.
2. Runner verifies prep proof before opening the NEW private internal entry of `runEpisodeFinishWith` (`episode_finish.go:63`, W16). F1 remains first on all public invocations, including the flag, and typed mode still reads the terminal. Internal entry requires a claimed, current, locally verified finish-prep capability and a surface-only mode. No public flag or env bypass. Supply safe structured progress output and no tty reader. Do the sandboxed checks and commit, applying §3.3. Probe `seam.Prepare` / Git auth as the runner **before** mutation in V4, even on resume. No unsupported flags and no unsandboxed checks in any runner mode.
3. At the V3 prompt, issue the real canonical `VerbRequest` from actual commit/body and canonical argv using the existing helper (W15–W18). Persist its reserved id/challenge/expiry and exact bytes before issuance. Suspend on the durable `awaiting PR confirmation` state. Daemon still owns V3 key verification and world decision/credential consume. When a V3 row is linked to an action, a NEW callback verifies that the action is live, in the awaiting phase and generation, under the same confirmation/cancellation arbitration. Check again during credentials consume. Standalone V3 cards have no linked action check. If approval fails to land the card is retriable at its V3 boundary before expiry. Runner restarts from its own checkpoint, never takes expected bytes from the row being verified.
4. Use `pollVerbRequest`, `authority.RecheckVerb` and `readBackFinishDecision` (W15/W18) with the runner's own bytes. Verify generation, approval timestamps, exact V3 id and key, decision scope and expiry; recheck F6 and all pinned dispatch data immediately before external effects. Hold the episode mutex or atomic reservation across local revalidation and write. Do not import reusable typed or another action's V3 approval from `finishReusableApproval` (W19). Keep the 610 s observation bound, with actual key authorization ending at the approval deadline. NEW write admission uses the earlier applicable deadline, never grace to authorize a new remote write. A confirmed V3 action can dispatch after the preparation proof expires only for the exact durable commit, as authorized by the independent V3 proof. An in-flight remote request can still be reconciled after expiry.
5. Reuse typed Git/GitHub operations through the dispatch seam, adding checkpoints and exact read-only verification for the runner. A V3 CLI operation or external actor that has already created a mismatched PR blocks this action. Reconcile before any retry, require a fresh key to authorize a later permitted write, and leave uncertain outcomes locked out when prior execution cannot be excluded. Finish retains the commit for expired or declined remote approval. Expose one result including PR URL, record and credential cleanup status.

7. The current record route validates approval expiry and writes the used head before the episode head (W28). Runner recovery therefore needs NEW private `record-observed` admission on the operator-proof route. Use the existing decision scope read and a NEW durable verified dispatch receipt binding action, exact PR, proof and start times. The daemon checks the recorded outbound reservation began under a live key and validates the same receipt against scope/episode. Its NEW store helper uses one durable transaction for both registry heads, following the store primitives in W28. The handler is reachable only with an operator proof plus matching runner receipt. It allows post-expiry recording of already-issued effects, then reuses exact-record/credential read-back. It may repair a prior matching used-head-only state by CAS if no newer episode head intervened. A mismatching head blocks with the PR link and attended reconciliation.
6. Mark's attended exit exercises start, launch, `episode report <ep>`, finish prep, V3 PR approval and the PR link from one browser. The review includes explicit check outcomes and protected changes. V4 opens a PR; merge, changing ready state on an existing PR and GitHub review actions are separate authority.

### 4.4 Interactive hand-off

Choose herdr over a row 206 integration for this phase, conditional on D-WORLD-119 A and E4. Row 222 is the confined interactive launcher predecessor, while row 206 is governed roadmap work (W02, W13, W22). Do not call the plugin's New World agent popup (W08). Reuse pane metadata and the installed confined `episode launch` interactive path (`episode_launch.go:300–322,517`, W13). Reserve one hand-off for the signed run id and a preconfigured session/workspace. The adapter uses fixed command form plus validated opaque action id, no task/auth/config in shell commands. It records pane id and hands a single-use private ticket to the confined launcher inside a TTY. Launcher validates ticket, exact action/start/scope, key and live credential before spawning. Reserve before creating a pane or calling pane run. Search by trusted namespace/workspace/run metadata and receipt; a lost pane creation reply or missing metadata after creation blocks another hand-off. The browser shows the pane destination or precise blocker. Missing executable, workspace, socket, version or confinement stops before mint when known. Never open an unconfined fallback. Under B or failed E4 return a terminal hand-off instruction after reporting the blocker; headless exit still works.

## 5. Branch traces and mutation plan

The following are implementation gates. Test names are NEW. A table entry gives the precise branch,
operator result and assertion that must fail under the listed mutation. No test here has been run
against implementation. Each test supplies a positive matching case in the same harness with a real
fixture key and counting fake effects. Positive-control fixture assertions are independent of
production allowlist constants. A red test with zero discovered tests or a skipped Chrome leg does
not count. Refusals assert both their public classification and absence of new effect calls.

### 5.1 Bar, schema and request admission

| Input / failure | Decision and visible result | Test and mutant killed |
|---|---|---|
| Existing address, report, why and approve | Original navigate/read/result paths | `TestOperatorOtherVerbParity` kills routing every GET into the action branch. |
| Valid supported start | Read-only review panel, then one pending card on explicit POST | `TestOperatorBarStartReview` kills inserting from Enter/GET. Assert GET leaves action count zero. |
| Valid supported finish | Read-only finish intent preview, then one pending prep card | `TestOperatorBarFinishReview` kills handing the raw line to execution. Assert zero calls before proof. |
| Case mismatch in verb, flag or id | Refused, fixed grammar shown | `TestOperatorGrammarCase` kills case-folding. |
| Extra positional or missing episode, preset or flag value | Usage refusal | `TestOperatorGrammarArity` kills each arity/value condition separately. |
| Prefix or substring resembling a supported root or flag | Refused or original address classification, no action | `TestOperatorGrammarExactWords` kills token prefix and substring matches separately. |
| Quoted/equals/interleaved accepted flags | Same canonical plan; preserve CLI-equivalent tokens | `TestOperatorGrammarCanonical` kills loss of value semantics on the shared parse path. |
| Duplicate flag, even mixed spellings | Usage refusal | `TestOperatorGrammarDuplicates` kills per-flag uniqueness. |
| Unsupported start/finish option, explicit false or standalone launch | No review and correct terminal guidance | `TestOperatorGrammarOptions` adds one disallowed member at a time, including dry-run, grants, branch, task path, credentials/store/daemon and waiver overrides. |
| Denied authority mode | Terminal guidance, including interactive under 119 B or ready under 120 B | `TestOperatorGrammarPolicy` kills each answer-dependent gate separately. |
| Enrol/revoke/mint/service/publish and arbitrary shell | Refused or unchanged read/navigation classification, no action | `TestOperatorNeverDispatch` inserts each forbidden family into the execution allowlist. |
| Episode/preset/action id with prefix/suffix junk | Refused | `TestOperatorIdentifierAnchors` separately removes leading and trailing anchors. |
| Multiline id or case mismatch | Refused | `TestOperatorIdentifierFlags` adds multiline and ignore-case flags separately. |
| Identifier just outside allowed length | Refused | `TestOperatorIdentifierLength` widens lower/upper limit independently. |
| Empty or over-cap request/body/task | Bad request / too large with a field-specific fix | `TestOperatorRequestLimits` disables each zero/max gate and checks both boundaries. |
| Non-object, non-array argv, non-string element or type-confused status | Unreadable/refused, no link or action | `TestOperatorWireTypes` removes each shape check independently. |
| Unknown/duplicate keys, trailing JSON, invalid byte or NUL | Bad request | `TestOperatorStrictJSON` allows each invalid class separately. |
| Shell metacharacters inside literal task | Signed and displayed literally; never become an argument fragment or HTML | `TestOperatorLiteralTask` kills shell interpolation and unescaped rendering separately. |
| Bad Host | 421 before reading body or admitting | `TestOperatorMisdirected` bypasses `approvalRoute` and mutates Host exactness. |
| Wrong scheme/host/port, missing or null Origin | 403 | `TestOperatorCrossOrigin` removes exact Origin and substitutes substring, positive cases loopback and configured. |
| Cross-site, same-site or malformed fetch metadata | 403 | `TestOperatorFetchSite` relaxes strict fetch-site membership. |
| Non-JSON content type / body over read cap | 415 / 413 | `TestOperatorMediaType` bypasses guard and `TestOperatorBodyBound` removes cap. |
| OPTIONS or other wrong method | Method refusal with no CORS grant | `TestOperatorUnsupportedMethod` registers side-effect handler for GET/OPTIONS. |
| Runner absent, paused, stale, wrong owner, wrong build/store, no compatible schema | Read-only terminal fallback naming reason, direct POST refuses | `TestOperatorRunnerUnavailable` removes each readiness predicate independently and checks direct POST. |
| Catalogue expired, identity changed, unavailable, missing preset or busy storage | Review stale/unavailable or recoverable store error | `TestOperatorCatalogueAdmission` skips each freshness, generation, membership and error guard. |
| Active action of another kind, 17th pending or reused op id with different intent | Conflict or queue full, existing work visible | `TestOperatorAdmissionBoundaries` drops unique episode constraint, cap or digest match in separate mutants. |
| Same operation id, or two taps with distinct ids on identical intent | Existing link and one action | `TestOperatorDuplicateTap` removes operation id unique index and separately disables intent dedup. |
| Lost POST response, refresh or daemon restart | Retry reads the same action; no double insert | `TestOperatorReconnect` clears durable dedup at restart. |
| GET timeout, POST timeout or schema-invalid reply | Keep the form, offer safe reconnect; no alternate execution | `TestOperatorPanelNetwork` navigates on GET failure or changes operation id on retry. |
| Superseded Enter/fetch response | Stale response dropped without page navigation | `TestOperatorPanelLatest` removes request generation/abort checking. |
| Crafted status/card/PR link | No link unless exact route or pinned HTTPS PR destination | `TestOperatorLinkAllowlist` independently removes path anchors, adds regex flags, permits suffix/substring, enables case-folding, widens id arity/length, accepts object href, omits host/port/scheme/repo/userinfo/query/fragment checks. |

### 5.2 Key confirmation, status and claim

| Input / failure | Decision and visible result | Test and mutant killed |
|---|---|---|
| Good passkey and good enrolled SSHSIG | Confirmed with correct approver; one CAS | `TestOperatorConfirmBothKeys` skips proof or discards recorded approver. |
| Plain y, no proof, malformed proof, bad challenge | Refused, pending stays inert | `TestOperatorConfirmProof` admits each invalid variant independently. |
| Forged, unknown, wrong SSH namespace or WebAuthn UV/RP/origin | Refused | `TestOperatorConfirmAuthenticity` bypasses each check independently. |
| Revoked key or replayed/unchanged nonzero counter | Refused | `TestOperatorConfirmRevocationCounter` omits revocation or counter comparison. |
| Proof from another action, episode, kind, store, argv, preset or task | Refused | `TestOperatorConfirmBinding` drops each covered field or allows a V3/A3 artifact. |
| Expiry equal to now, in past or outside issue/lifetime bounds | Expired | `TestOperatorConfirmExpiry` changes strict equality boundary or drops max/issued bound. |
| Approved row replay / lost confirm response | Already confirmed, one counter advancement | `TestOperatorConfirmReplay` allows second CAS or repeats counter use. |
| Pending row missing, read-only DB busy, counter CAS lost, writer lock or daemon death | Missing, retry or confirmed result only when durable | `TestOperatorConfirmCrash` split-writes state and counter, or maps store error to success. |
| Timeout, cancel wins CAS or runner offline between confirm and claim | Expired, cancelled or pending until fresh runner | `TestOperatorClaimBeforeExpiry` ignores each of expiry, cancellation and readiness. |
| Two runners or stale generation | Exactly one claim and reservation | `TestOperatorClaimRace` drops exclusive lock and separately removes transactional CAS/current-owner comparison. |
| Forged DB row / mismatch after claim | Refused, no effects | `TestOperatorRunnerRechecks` trusts row digest, stored approver or daemon status without proof. |
| Process restart before reserve | Recheck and recover untouched claim | `TestOperatorRecoverUntouched` marks a pristine claim orphaned or skips proof on recovery. |
| Process restart after reserve / live child / unknown child | Adopt matching child, or uncertain lockout | `TestOperatorRecoverReserved` reruns after heartbeat timeout or ignores birth identity. |
| Status poll offline, read error or invalid state | Last known state with timestamp and explicit disconnect | `TestOperatorStatusDisconnected` turns stale into success or failed into empty. |
| Pruning and cap boundary | Only terminal metadata ages out, tombstones still dedup | `TestOperatorPrune` removes tombstone or prunes uncertain rows. |
| Cancel after in-flight checkpoint | Reconcile before marking cancelled, no cleanup of effects | `TestOperatorCancelInflight` releases reservation or calls delete. |

### 5.3 Start and launch

| Input / failure | Decision and visible result | Test and mutant killed |
|---|---|---|
| Ready configured headless preset and valid key | One matching start/worktree/mint/run, stable link | `TestOperatorStartOnce` launches before receipt or ignores confinement. |
| Preset/base/profile/task/daemon changed after approve | Fresh review needed, no effect | `TestOperatorStartChanged` drops one digest/generation/identity at a time. |
| Existing episode, branch, worktree or token/task/MCP path | Conflict, no adoption or overwrite | `TestOperatorStartCollisions` relaxes each absence or canonical path predicate. |
| Worktree add returns error or partial outcome | Failed if absent and unambiguous, otherwise uncertain | `TestOperatorWorktreeWindow` assumes successful add from a path alone. |
| Signed grants differ from effective grants/TTL or revoked approver at mint | Blocked, no credential issued | `TestOperatorMintBinding` removes exact grant/TTL, proof or fresh clock comparison. |
| Mint busy/error before transaction commit | Previous receipt or bounded retry after verified rollback | `TestOperatorMintRollback` leaks start or session when transaction aborts. |
| Crash after mint, before token write/checkpoint | Find mint receipt, revoke only orphaned credential, no remint | `TestOperatorMintReceiptWindow` writes receipt after mint commit or generates a second credential. |
| Task/token/MCP file mismatched, linked or unreadable | Blocked with repair path; never overwritten | `TestOperatorFilesReconcile` trusts mode/digest/path alone or opens with truncation. |
| Session expires, revoked or changed before launch | Expired or blocked, no start | `TestOperatorLaunchSession` removes expiry/revocation/current-session match separately. |
| Child exec missing, permission denied, spawn refused | Failed with action and known no spawn | `TestOperatorSpawnRefused` reports success after spawn error. |
| Crash just before or after launch, live child, PID reuse | Use durable run/birth proof, else uncertain | `TestOperatorLaunchCrashWindow` reserves after spawn or accepts PID without birth/run identity. |
| Nonzero exit, timeout, output failure | Report stopped/failed/capped with scrubbed evidence | `TestOperatorLaunchExit` conflates child stop with admission, drops wall cap or logs fake canary. |
| Browser closes at confirmation; daemon restarts at mint/launch | Same status link and one start, no session duplication | `TestOperatorStartBrowser` ties child lifetime to request context or repeats mint after restart. |
| Interactive selected without herdr/socket/version/workspace/confinement | Named blocker before mint; terminal alternative | `TestOperatorHerdrUnavailable` skips each preflight or calls plain fallback. |
| herdr ready, exact action and TTY | One pane/launcher with metadata, scoped ticket | `TestOperatorHerdrHandoff` omits ticket/TTY check or starts another agent in an existing pane. |
| Pane reply lost, stale ticket, no metadata, or socket goes down | Reconcile matching run, otherwise uncertain/block | `TestOperatorHerdrWindow` reruns pane creation on transport error or trusts an old ticket. |

### 5.4 Finish, F1 and remote reconciliation

V3's C0–C15, CF, CX and XC branches still run through the existing fixtures (W24, W26).
The runner's new branches below are additional.

| Input / failure | Decision and visible result | Test and mutant killed |
|---|---|---|
| Public interactive typed finish | Existing tty prompt and typed y | `TestOperatorFinishTTYParity` removes prompt or routes typed path through private capability. |
| Public no-tty finish, with/without flag or caller-supplied bypass | F1, before local/remote effects | `TestOperatorFinishF1` drops F1 when flag is set or uses env/flag marker. |
| Private command with missing/wrong prep proof | Refused before checks/commit | `TestOperatorFinishNoProof` opens private entry without typed proof. |
| Valid prep, completed launch, attested clean/pinned manifest | Run pinned checks, reserve commit and return V3 confirmation | `TestOperatorFinishPrep` commits before proof or skips durable commit reservation. |
| No start, degraded or empty tree, running agent, sessions changed/expired or read failure | Block with episode/check/session blocker | `TestOperatorFinishEligibility` drops each report, tree, stop or session predicate separately. |
| Stale preset, worktree HEAD, untracked file, path race, remote/base or profile | Re-review; no execution on changed inputs | `TestOperatorFinishStale` drops each matching predicate or runs diff after commit. |
| Unsupported protected/unattributed change or untrusted Git hook/helper | Refusal naming operator input required | `TestOperatorFinishWorktreeFence` trusts repo-controlled hooks/config, ignores path guard or suppresses visible blocker. |
| Pinned unsandboxed check, including resume | Terminal-only check blocker before new effect | `TestOperatorFinishUnsandboxed` moves guard into fresh-only path or after first check. |
| Missing gh / expired login / Git SSH or HTTPS failure | Finish unavailable before V4 effects | `TestOperatorFinishAuth` calls `seam.Prepare` after commit or uses gh result as proof of Git auth. |
| Check fails, aborts, times out or has unexpected mutations | Block before commit with bounded output | `TestOperatorFinishChecks` ignores each outcome, reruns unknown check or commits changed input. |
| Worktree commit definite failure | Fail with retained evidence and no prompt | `TestOperatorCommitFailed` proceeds after commit error. |
| Commit reply lost, partial stage, HEAD advanced differently | Adopt exact proof or show uncertain | `TestOperatorCommitReconcile` matches trailer alone, removes parent/tree/message check, or resets on conflict. |
| Prompt too large, entropy/store busy, V3 cap | No PR dispatch, commit retained and seen | `TestOperatorFinishCardWriteRefusal` synthesizes default approval or loses commit receipt. |
| Daemon dies during V3 store boundary / runner dies before V3 row insert | Recover one linked request from frozen bytes, confirm only within lifetime | `TestOperatorFinishCardRestart` rotates id/challenge on retry or reconstructs expected bytes from served row. |
| Row missing, superseded, expired, read failed or key revoked | Visible specific blocker, commit retained | `TestOperatorFinishCardRefusal` skips each refused state, revocation or actual signed deadline. |
| Typed or another action's reusable approval | Require linked V3 confirmation | `TestOperatorFinishRejectReuse` runs existing reuse branch for runner. |
| Real V3 proof against a different prompt or decision | Refusal; zero push/PR | `TestOperatorFinishOwnBytes` trusts daemon decision or replaces either digest/challenge/expiry with row data. |
| Changed HEAD/remote/body/title / branch target conflict after PR approval | F6/stale with commit kept; fresh review | `TestOperatorFinishPreDispatchFence` drops each final comparison. |
| Remote head equals approved SHA / valid ancestor | Skip push when equal; only approved non-force push otherwise | `TestOperatorPushIdentity` pushes again on equality, uses HEAD instead of SHA or adds force. |
| Remote diverged, lookup unavailable or uncertain push | Block or reconcile exact ref, with no blind retry | `TestOperatorPushUncertain` treats missing ref after timeout as permission to push again. |
| Label exists, missing, creation refused or reply lost | Verify pinned label identity/scope before creating or declaring observed | `TestOperatorLabelReconcile` creates before proof, repeats on unknown outcome or uses repo from response. |
| One matching PR exists | Adopt exact result without create | `TestOperatorReadExistingPR` ignores verified result and creates anyway. |
| Mismatching, duplicate or closed/merged PR evidence | No additional create; report the observed mismatch/status | `TestOperatorPRIdentity` drops each repo/head/OID/base/draft/title/body/label check, filters away closed and merged, or selects the first of several. |
| PR creation definite refusal before send | Failed, expose explicit fresh authorization for permitted retry | `TestOperatorPRNoSend` retries without a new key. |
| GitHub timeout, lost response, stale or empty reconciliation | Show checking/uncertain, no second create | `TestOperatorPRResponseLost` turns an empty search into automatic retry. |
| Browser closes during submission; daemon restarts before writeback | Same URL and one finish record after verified read-back | `TestOperatorFinishBrowser` duplicates create/record or loses the stable action link. |
| V3 card references expired/cancelled or replaced runner action | Refuse to consume a linked card; standalone V3 continues | `TestOperatorFinishLinkedCard` skips each linked action generation/phase/expiry check or applies it to every V3 row. |
| Prep expires after successful sandboxed checks | Commit remains absent, receipts available for explicit reauthorization | `TestOperatorFinishPrepExpiry` uses initial admission time to authorize commit or automatically reruns checks. |
| Record arrives after approval expiry with exact receipt, or forged receipt | Record only a previously authorized effect; preserve valid PR link | `TestOperatorFinishObservedRecord` makes expiry reopen dispatch, drops initiation clock/scope/receipt equality or accepts missing proof. |
| Store failure between used and episode heads, or stale previous head | Recover exact record or stop on conflict | `TestOperatorFinishRecordCAS` writes heads in separate transactions or overwrites a later head. |
| Finish record write/receipt failed after PR; detach/revoke retries | PR link and incomplete cleanup; exact scoped reconciliation | `TestOperatorFinishCleanup` creates again, detaches changed HEAD or revokes a newly resumed credential. |

### 5.5 Service admission, versioning and lifecycle

| Input / failure | Decision and visible result | Test and mutant killed |
|---|---|---|
| Install invoked without typed approval / stale approved plist | Refuse install or start, retain existing state | `TestOperatorServiceFence` removes typed y, accepts edited plist or registers service operations from the bar. |
| Root, mismatched UID, writable or linked config/spool, changed binary | Fail before starting companion | `TestOperatorServiceOwnership` removes each owner/mode/path/hash check separately. |
| Bad config, unknown fields, absent binary or unknown runner launch option | Explicit invalid config/setup result | `TestOperatorServiceConfig` opens the schema or accepts raw command at internal mode. |
| Synthetic restart and mismatched store/build | Reenter only compatible pair, report restart pending on mismatch | `TestOperatorServicePair` accepts wrong store, ignores either binary hash or opens work prematurely. |
| Daemon down, launchctl boot failure, runner start timeout or schema skew | Bounded diagnostic and no queued execution | `TestOperatorServiceStartup` suppresses each failure, starts unpaired or waits indefinitely. |
| Upgrade with active child, bad swapped binary or failed paired restart | Drain or stop with receipt, rollback only to schema-compatible pair | `TestOperatorServiceUpgrade` skips drain, erases reservation, or starts incompatible reader. |
| Correct synthetic installation and stop/uninstall | Only companion lifecycle changes, retain action evidence | `TestOperatorServiceWriteScope` changes base service path or deletes persistent state. |
| Runner permission/environment/credential fields forged by browser | Refused, no environment changes | `TestOperatorServiceEnvironment` allows browser-sourced vars or logs fake secrets. |
| Fresh schema, upgrade from v9, legacy reader and migration failure | Parity with existing rows; incompatible reader refuses; failed tx rolls back | `TestOperatorSchema10` omits a migration, accepts newer store or partially commits DDL. |

### 5.6 Judge-405 pins and package gates

S1 takes all three small pins. Same harness includes a positive valid V3 row and valid decision href.

| Test | One mutation and assertion |
|---|---|
| `TestNeedsYouConfirmationEmbeddedNewline` | Add `m` to the id regex in `confirmLineHTML` (`app.js:456–458`, W23). A valid-looking line followed by LF and junk must render no link. Also mutate each anchor and `i` separately. |
| `TestNeedsYouConfirmationsNonArray` | Drop `Array.isArray` (`app.js:434`, W23). Object/string/null fixtures render the other queue items without throwing or linking the malformed entry. |
| `TestVerbDecidePrefixBoundary` | Replace decision `startsWith` with `includes` (`app.js:507`, W23). A foreign origin/path containing the prefix must refuse. Also remove the colon delimiter and check that altered prefix is refused. |

For each mutation, run the **whole affected package** with fresh caches as needed and record
KILLED/SURVIVED with the named failing assertion. Repeat schema/claim/restart tests under race with
bounded clocks and subprocesses. Use fake `gh`, fake remotes, fake model/agent and launcher,
synthetic HOME and seeded test keys only. No real enrolment, mint, agent launch, publish or merge in
the loop. Fake process exit, delayed/reordered writes and corrupt receipts explicitly. Check all
real-path guarantees as observations of sha256 and line count (or `absent`) before/after under
synthetic HOME only. Never seed the user's real credentials, plist, log or state path.

Live Chrome virtual-authenticator coverage gates each page change. Extend the existing S6 fixture
seams (W24/W26). S1 pins queue rendering with one live smoke. S8 runs start from bar through
request/key/reconnect with fake launch, invalid payload/refusal and phone viewport (390 × 844),
focus and back link. S10 runs the V3 card as a second gesture on one stable action URL, with expiry.
S11 runs a lost remote response and reconnect with fake GitHub, plus exact PR URL. When an
environment denies loopback or browser access treat it as UNINFORMATIVE UNDER SANDBOX and schedule
the gate on the controlled acceptance host. Include the executed-verbatim fake Quickstart at each
new usage surface, plus the repository verification gates (W01).

## 6. Threats and placement

| Threat | Control and remaining boundary | Pin |
|---|---|---|
| Forgery, revoked enrollment, wrong key type/counter | Independent runner recheck against the current enrolled key; signed ceremony with UV; counter recorded once by daemon and checked atomically | ConfirmAuthenticity, ConfirmRevocationCounter, RunnerRechecks |
| Replay, stale timestamp, cross-episode/action/argv/preset/store | Domain-separated canonical message, signed identity and random challenge, absolute expiry and CAS; V3 proof linked to one exact frozen action | ConfirmBinding, ConfirmExpiry, ConfirmReplay, FinishOwnBytes |
| Page tries to act without a key / typed y | Only inert request writes; bounded same-origin POST; guard before effect lock or helper; action state is not authority | BarStartReview, ConfirmProof, FinishNoProof |
| CSRF, origin reflection, DNS rebinding, cross-site display | `approvalRoute` allowlist (`surface_approve.go:58–85`) plus explicit exact Origin/fetch site/JSON guards (`surface_approve_guard.go:41–64`, W26). No CORS or state-changing GET. Inherit CSP on cards (`surface_page.go:19–23`, W26). | Misdirected, CrossOrigin, FetchSite, MediaType, UnsupportedMethod |
| Runner tricked into different argv / shell code | Closed command variants and immutable resolved typed plan, fixed pipe of argument vectors. Interactive wrapper takes only a validated id/ticket. | GrammarOptions, LiteralTask, ConfirmBinding, HerdrHandoff |
| TOCTOU after approval | Per-episode exclusion, compare immutable bytes and live authorities before each new effect, effect-specific identity and conservative recovery. A same-UID operator remains within the trusted local boundary. | StartChanged, FinishStale, FinishPreDispatchFence, ClaimRace |
| Malicious child or arbitrary subprocess | Companion user UID has Mark's reach. Keep secrets out of checked-out repo, keep existing sandbox and model confinement, per-tool env allowlist, noninteractive auth and owner-validated absolute binaries. | StartOnce, FinishAuth, FinishWorktreeFence |
| Token leak by logs, argv or durable status | Token in protected existing file paths only. Scrub child output, do not expose artifacts or raw environments in regular logs. Use private checkpoints capped and mode-restricted. | LaunchExit, FilesReconcile, StatusDisconnected |
| Malicious PR/progress link | Exact type, route and repo validation; escape all variable text; never redirect from untrusted payload | LinkAllowlist, WireTypes |
| Crash or double runner as replay vector | Owner generation and lock, durable mutable reservations, read-only reconciliation after uncertainty | ClaimRace, RecoverReserved, MintReceiptWindow, PRResponseLost |
| Same-user compromise or compromised signed binary | UID/file/key roots are trusted; the runner cannot protect this account from code already acting as it. Build/config pinning and attended installation reduce accidental scope. | RunnerUnavailable, FinishWorktreeFence |

Placement follows coding standards S1–S3 (W01). NEW approval canonicalization and validators belong
in `host/approval` and adapters in `host/authority`; the persistence/transactions belong in
`host/store`; admission/transport in `host/daemon`; runner and UX seams in the CLI and surface.
These are OS identity, signature implementation, SQLite or browser transport boundaries. A registry
package cannot own an OS process credential, its SQLite transaction or a browser ceremony. No
`.ail`, kernel, effect catalog, driver, pipeline or publishing change is expected.

## 7. Files and collision surface

- `cmd/ailang-worldd` service, episode preset/start/launch/finish files and focused runner files.
  Reuse the bounded package seams, keep daemon operator dispatch subprocess-free.
- `host/{approval,authority,store,verbs,daemon,surfacepage}` and existing narrowly shared process
  helper files if a planned receipt needs them. Keep changes partitioned by package responsibility.
- `docs/QUICKSTART.md` operator setup/status, start and finish usage, fake examples and an attended
  hook. Update only the affected call paths and their runbook checks.
- PR **#604**, attended `world/mobile-needs-you`, open/unmerged per operator brief, lists
  `host/surfacepage/app/{app.js,app.css,index.html}`. Its hunks are unverified. Keep S1 test-only
  against HEAD; develop the action card as separate assets and fixtures. Touch the shared app shell
  in S8 only after #604 is resolved and the controller has remeasured the overlap, or after the
  controller explicitly sequences V4 first. No design dependency on #604 content. Avoid CSS or
  index-shell changes in V4 unless porting the approved hooks proves necessary.
- The optional herdr release work stays in `cmd` for the adapter and CLI hooks. Any plugin edit must
  fit the same 1,400 total limit including integration paths, with no edits to shared launchd tools.

## 8. Sprint split and gates

Raw estimates are planning budgets, not measured implementation. Calibrated changed lines below
use 1.35× rounded up. The V3 S5 exception hit 1.54× (W24), so the parallel planning cap is 900
raw, which leaves margin under 1,400 even at that ratio. Count **insertions + deletions** by
`git diff --numstat <sprint-base> -- cmd host docs app`, also inspect the unfiltered sum and
never hide work in another path. Re-measure staged patches and cut between milestones before
1,400. Include a working S7 fake example and help wherever a new public entry lands.

| Sprint | Milestones, raw budget | Raw → calibrated | Waits for |
|---|---|---:|---|
| S1 | Judge-405 three independent pins 210; pure admitted CLI-surface subset and parity 410 | 620 → 837 | None. Test-only shell changes and unregistered parser. |
| S2 | Action canonical envelope / fixtures 360; typed SSHSIG/passkey proof and recheck against fakes 480 | 840 → 1,134 | None. Pure law/adapters against synthetic fixtures. No migrations or served routes. |
| S3 | Schema 10 and upgrade parity 330; receipt, CAS/expiry/dedup and restart tests 560 | 890 → 1,202 | 119 A/B. |
| S4 | Typed runner install and isolated process/status/upgrade seam 520; heartbeat and inert catalogue fixture 360 | 880 → 1,188 | 119 A/B, E1/E3. |
| S5 | Shared pinned-start read seam and semantic parity 390; atomic minted-start proof, token recovery and cases 510 | 900 → 1,215 | 119 A/B, S3–S4. |
| S6 | Confined headless launch and process receipt 460; launch/restart, expiry and scrub testing 430 | 890 → 1,202 | 119 A/B, E3, S5. |
| S7 | Inert request and guarded action confirmation routes 540; status/Needs you projection against fixtures 330 | 870 → 1,175 | 119 A/B, S3–S6. |
| S8 | Standalone card and app bar bindings 490; Chrome fake Sol start, viewport/reconnect and usage 370 | 860 → 1,161 | 119 A/B, S7, E5 and #604 sequencing. **First browser slice.** |
| S9 | Verified private F1 entry and auth readiness 350; frozen manifest, sandboxed checks and local commit journal 550 | 900 → 1,215 | 119 A/B and 120 A/B, E2, S8. |
| S10 | Durable V3-confirm hand-off/read-back and failure windows 500; exact review/link and Chrome approval 370 | 870 → 1,175 | 120 A/B, S9. |
| S11 | Exact remote receipt reconciliation, record/revoke 540; uncertain-dispatch and full start-to-PR pins 360 | 900 → 1,215 | 120 A/B, E2/E5, S10. **Headless exit.** |
| S12 | Fixed herdr adapter/ticket 420; fake restart/missing-pane, tty confinement and documentation 280 | 700 → 945 | 119 A, E4, S8. Schedule after S11 when 120 = A/B, or after S8 when 120 = C. |

Total 10,120 raw / 13,664 calibrated including optional S12. Without S12, 9,420 raw / 12,719
calibrated. Parent ≈2,000 estimated the choice and a wrapper before measuring. The additional
work is a read-only catalogue, new action proof and transactional mint/start, execution ownership,
local and remote write-ahead receipts, process and credential recovery, guarded UI and test matrix.
W24 gives measured surface-area and V3 calibration, and S5/S8/S11 are explicit fallback split
points if a compiled milestone breaks the budget. Split tests or an independently useful vertical
milestone before activation; never cut an authority or recovery gate to make the count.

## 9. Residuals and explicit scope choices

- **RV1 / Ctrl-C and RV5 / denial.** The public V3 CLI retains its existing behavior (W18/W27).
  V4 uses restartable pending actions. An attended cancel expires the linked dispatch capability
  and blocks V4-linked confirmation consumption. The pre-existing V3 two-store orphan window can
  still leave a world decision with no usable row. No unsigned browser cancel. Do not extend
  V3's row outcome enum for this.
- **RV2 / two stores.** Existing V3 orphan decision window persists (W26). Runner still requires
  row-bound evidence plus decision read-back; a standalone world decision cannot claim an action.
- **RV3 / per-check unsandboxed confirmation.** Out. It would add a separate program/cwd/env
  authorization, operator display and recovery protocol. Surface finish refuses and names each
  unsandboxed check. E5 checks that the attended exit preset can use this lane.
- **RV4 / F1.** Private verified surface path only, per §4.3. Public CLI surface waiting keeps F1.
- **RV6 / GET visibility and RV8 / provenance.** Share current Host policy, minimize lists to short
  redacted summaries, show exact signed data only at its card, link decision and episode report.
  Do not make status depend on reconstructing the request from live input.
- **RV7 / schema.** Compatible binaries and restart are operationally paired. Stage a stopped
  companion while handling a previous schema; current daemon is the schema upgrader (W24).
- **RV9 / size.** Measure each sprint under real numstat, budget rollovers per §8.
- **Judge-402 note 2 / S3 §1.4 F1.** Defer V3's day-expired same-episode supersession/prune ordering
  and permissive superseded decision-ref DDL. V4's action DDL pins strict transition/artifact
  invariants. Route V3 changes separately because v9/v10-upgraded stores must match fresh stores and
  V3's ledger trace stays reviewable (W15/W27).
- **S6-D9** finite V4 fake clocks and bounded acceptance harnesses. Preserve S6's existing rig
  identity and route any frozen fixture change through its owning plan (W27).
- **S6-D3 / notification.** Out. Show pending work in Needs you and the stable episode/action link.
  Notification dedup and delivery across devices are a separate mechanism, with no effect on
  authorization. A new notification is unnecessary for the addressed operator's attended run.
- **Merge.** Separate future ledger decision; V4 adds no merge, mark-ready or review authority.
  V4 supplies the exact PR destination and observed outcome for that later work.

## 10. Ledger asks for the controller

> | D-WORLD-119 | OPEN | **Row 243 V4: may the World app start and launch a confined episode through a user-level runner after an enrolled-key approval?** The runner is an opt-in service companion running as the operator. The key approves the exact task, preset, repository and base, grants and lifetime, and launch mode. The runner rechecks the proof and pinned inputs before creating the worktree, minting a session and launching the agent. **A (recommended):** yes for headless launch, and for a later confined herdr hand-off once the same proof and a no-fallback rig test pass. The headless path lands first and is sufficient for the exit run. **B:** yes for headless only. Interactive launches stay attended in their terminal path. **C:** no. Surface start and launch stay terminal-only. **Default while unanswered:** build S1 and S2 against fixtures only. Leave the schema, service installation, browser writes, mint and launch paths unchanged. S3–S8 wait for A or B. S12 waits for A. The loop uses fakes and test keys only. | Filed iteration 406 (unattended loop). Evidence: V4 design §2–§5, §8 and §12. |

> | D-WORLD-120 | OPEN | **Row 243 V4: may the World app run finish without a terminal using a separate enrolled-key approval for local preparation and the V3 approval for the exact pull request?** The preparation approval binds the pinned checks, worktree input and local commit intent. The operator then signs the actual post-commit prompt before the runner may push or open the PR. The runner verifies both proofs and reconciles uncertain GitHub results before any further write. **A (recommended):** yes, for draft or ready pull requests. Unsandboxed checks and unsupported overrides stay at a terminal. **B:** yes for draft pull requests only. Ready pull requests stay at a terminal. **C:** no. Keep F1 for all finish invocations. **Default while unanswered:** keep the public F1 terminal guard and do not build or enable the private runner finish path. S9–S11 wait for A or B and the delivered A/B start path under D-WORLD-119. D-WORLD-118 continues to govern the terminal-launched V3 confirmation. No option authorizes merge, publish or service/key administration from the app. | Filed iteration 406 (unattended loop). Evidence: V4 design §3–§5, §8 and §12. |

If 119 = C, leave only S1 and the dormant, policy-independent portions of S2. Do not plan S3–S12
for this design. If 120 = C, S1–S8 can deliver start under 119 A/B and S12 can follow under 119 A.
A route activated by S8 must already have earned every applicable start gate. A finish route
activated by S11 must already have earned both finish gates.

## 11. Named planner experiments

The Studio observations in W05–W07 and W22 informed Q4. E1–E5 below remain **UNVERIFIED**.
E1 was denied in this sandbox. E2's first-stage PATH result is measured, but a stored-token error
with restricted egress cannot establish the real job's credential parity. E3 and E4 require
supervisor/process operations outside this designer assignment. E5 requires Mark's chosen preset.
Run the relevant experiment before planning its gated milestone. A sandbox denial is
**UNINFORMATIVE UNDER SANDBOX** and requires a capable controlled host, regardless of exit code.

### E1 — running Studio job UID and identity (S4)

Run read-only on the Studio. Filter before logging the job data.

```sh
job="gui/$(/usr/bin/id -u)/com.sunholo.ailang-world.worldd"
job_pid=$(/bin/launchctl print "$job" | /usr/bin/awk '$1=="pid" && $2=="=" {print $3; exit}')
/usr/bin/id -u
test -n "$job_pid" && /bin/ps -p "$job_pid" -o uid=,comm=
```

Pass if the running executable is the expected installed World binary and UID matches the GUI
user. Missing/mismatched UID is a deployment blocker if the probe ran successfully with
permission. No elevated identity is needed for the companion. W06 recorded `ps` denied here.

### E2 — noninteractive GitHub and Git identity (S9/S11)

Set `V4_REMOTE` to the exact operator-approved GitHub remote and `V4_REPO` to its `owner/name`.
Set `V4_RUNNER_PATH` to the proposed approved PATH and, only for an approved SSH route, set
`V4_SSH_SOCKET` to its measured socket. These are read-only rig probes. Use fake endpoints for
loop test execution. Do not log returned credential data or raw environment values.

```sh
/usr/bin/env -i HOME="$HOME" PATH="$V4_RUNNER_PATH" \
  GH_PROMPT_DISABLED=1 GIT_TERMINAL_PROMPT=0 V4_REPO="$V4_REPO" \
  V4_REMOTE="$V4_REMOTE" SSH_AUTH_SOCK="${V4_SSH_SOCKET:-}" \
  GIT_SSH_COMMAND='/usr/bin/ssh -o BatchMode=yes' /bin/sh -c '
    command -v gh && command -v git || exit 2
    gh auth status --hostname github.com >/dev/null 2>&1 || exit 3
    gh label list -R "$V4_REPO" --search ailang-world --json name >/dev/null 2>&1 || exit 4
    gh repo view "$V4_REPO" --json viewerPermission \
      --jq ".viewerPermission" 2>/dev/null || exit 5
    git -c core.sshCommand="/usr/bin/ssh -o BatchMode=yes" \
      -c credential.interactive=never ls-remote --symref "$V4_REMOTE" HEAD \
      >/dev/null 2>&1 || exit 6
    printf "%s\n" rig-auth-read-ok
  '
```

Pass for readiness planning if both tools resolve, repo access and remote read pass without a
prompt, and GitHub reports WRITE, MAINTAIN or ADMIN for the configured account. Exit codes
2–6 identify the failing step. This is read-level evidence for the exact remotes and a permission
signal. Actual write permission and branch rules remain observed outcomes of Mark's attended
finish; the probe performs no remote write.

### E3 — sanitized GUI LaunchAgent and restart on a fake child (S4/S6)

Run on the Studio under the Mark UID, with **all** test files scoped to a synthetic home and a
unique synthetic label. This experiment installs and stops only its fake job. Confirm no live
World label is touched. The fake child is a one-shot shell observation; no model is launched.

```sh
set -eu
V4_SYNTH_HOME=$(mktemp -d)
V4_UID=$(/usr/bin/id -u)
V4_LABEL="org.example.v4-fixture.$V4_UID.$$"
V4_JOB="gui/$V4_UID/$V4_LABEL"
V4_PLIST="$V4_SYNTH_HOME/$V4_LABEL.plist"
trap '/bin/launchctl bootout "$V4_JOB" >/dev/null 2>&1 || true' EXIT
obs() {
  for file in "$V4_SYNTH_HOME/receipt" "$V4_SYNTH_HOME/.ailang-world/service.json" \
              "$V4_SYNTH_HOME/Library/LaunchAgents/com.sunholo.ailang-world.worldd.plist"; do
    if [ -f "$file" ]; then /usr/bin/shasum -a 256 "$file"; /usr/bin/wc -l "$file"
    else printf 'absent %s\n' "$file"; fi
  done
}
obs
cat >"$V4_SYNTH_HOME/record" <<'SH'
#!/bin/sh
set -eu
printf '%s %s %s %s\n' "$(/usr/bin/id -u)" "$$" "$HOME" "$PATH" >>"$HOME/receipt"
SH
/bin/chmod 700 "$V4_SYNTH_HOME/record"
cat >"$V4_PLIST" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0"><dict>
<key>Label</key><string>$V4_LABEL</string>
<key>ProgramArguments</key><array><string>/bin/sh</string><string>$V4_SYNTH_HOME/record</string></array>
<key>EnvironmentVariables</key><dict>
<key>HOME</key><string>$V4_SYNTH_HOME</string>
<key>PATH</key><string>/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
</dict></dict></plist>
PLIST
/bin/launchctl bootstrap "gui/$V4_UID" "$V4_PLIST"
for n in 1 2; do
  /bin/launchctl kickstart "$V4_JOB"
  limit=$(( $(/bin/date +%s) + 10 ))
  until [ -r "$V4_SYNTH_HOME/receipt" ] &&
    [ "$(/usr/bin/wc -l <"$V4_SYNTH_HOME/receipt")" -ge "$n" ]; do
    [ "$(/bin/date +%s)" -lt "$limit" ] || exit 7
    /bin/sleep 1
  done
done
/usr/bin/awk -v id="$V4_UID" -v h="$V4_SYNTH_HOME" '
  $1!=id || $3!=h || $4!="/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin" { bad=1 }
  !seen[$2]++ { identities++ }
  END { exit !(NR==2 && identities==2 && !bad) }
' "$V4_SYNTH_HOME/receipt"
obs
```

Pass if the fake job runs twice with two process identities and the explicit HOME/PATH/UID, and
the two synthetic World files are observed unchanged or absent. This establishes service
environment control and restart reentry. S6 then pins the crash windows in §5.3 with a separately
supervised fake child implementing the proposed receipt. E3 is not evidence that real Sol auth
or launch works.

### E4 — herdr hand-off contract (optional S12)

W22 measured CLI help on the Studio. For the live socket portion use an **already running synthetic
herdr session** on the controlled host, one known fake workspace and fake pane. Start the synthetic
server only under a synthetic HOME, with separately verified socket and configuration paths. Set
`V4_HERDR` to the absolute binary, `V4_SYNTH_SESSION` to that unique fixture session, and
`V4_SYNTH_PANE` to its opaque id read from the snapshot. Run these against the fake pane only.

```sh
"$V4_HERDR" --session "$V4_SYNTH_SESSION" api snapshot
"$V4_HERDR" --session "$V4_SYNTH_SESSION" pane run "$V4_SYNTH_PANE" \
  '/usr/bin/printf "%s\n" world-v4-fixture'
"$V4_HERDR" --session "$V4_SYNTH_SESSION" pane read "$V4_SYNTH_PANE"
```

Pass if the expected marker is observed in that one pane and the before/after snapshot names the
same synthetic workspace/pane and no second launch. Check isolated config/log paths only by
sha + line-count observations under the synthetic HOME. Simulate a lost response and require the
planned adapter's fake `TestOperatorHerdrWindow` to stop or adopt from the exact receipt. If the
daemon socket cannot give an unambiguous outcome, the state is uncertain. If the CLI/mock cannot
pin the workspace, pane and confined launcher, ship headless and keep the interactive route closed.

### E5 — actual preset eligibility (before planning S8/S11)

Ask Mark to point to the approved Sol preset with `V4_PRESET_FILE` on the rig. This is a read-only
probe against the chosen file. Print counts and mode only.

```sh
/usr/bin/shasum -a 256 "$V4_PRESET_FILE"
/usr/bin/wc -l "$V4_PRESET_FILE"
jq -er '{
  agent,
  model_present: ((.model // "") | length > 0),
  mode: (.launch.mode // "default"),
  checks: ((.checks // []) | length),
  unsandboxed: ([.checks[]? | select(type == "object") |
                  select(.unsandboxed == true)] | length)
}' "$V4_PRESET_FILE"
```

Pass for preset review if it resolves to Mark's approved host/model, its grants and environment
are suitable under the runner, effective headless mode and its cap validate, and effective resolved
finish checks are all sandboxed (W32). The default interactive preset mode can be overridden by
the explicit headless form, with the result visible and signed (W32). Run the existing `preset check` path against fakes in planning for
resolution/parity (W25); the JSON probe alone establishes only the file's explicit settings.
A real model invocation, if needed to validate provider auth, belongs to the attended exit.

## 12. Verification Log

Commands were run from the assigned worktree unless noted. File names in command columns are
relative to that known root. Outputs are trimmed to inspected signals. `F` below denotes
`cmd/ailang-worldd`, `D` denotes `host/daemon`, and `P` denotes
`integrations/herdr/ailang-world`. Expand these aliases as paths when reproducing a command.
No Go build, Chrome, new daemon, real mint/agent/enrolment, GitHub write or git write was run.
The read-only `gh` probe output was filtered before logging.

### W01

**Commands.** `git rev-parse HEAD && git branch --show-current && git status --short`; `git log --oneline -7`; `cat CLAUDE.md design_docs/coding-standards.md`; `sed -n '27,48p;610,641p' design_docs/DESIGN.md`

**Observed.** `41331b4d880e4ea1a4c19a2f1e65a74204580bc4`; `design/world-iter406-row243-v4`; initial status empty, positive base/branch. HEAD merge #603, V3 S6 merge `0077a7b0`. CLAUDE verify gate and frozen boundaries; standards S1–S8; DESIGN immutable kernel and host boundary.

### W02

**Commands.** `rg -n '^\| D-WORLD-(94|95|98|101|106|107|109|115|118)' design_docs/world-mission.md`; required ledger rows via `rg -n '^\| *\*?\*?D-WORLD-(89|94|107|109|115|118|119|120)\b' design_docs/world-mission.md`; `rg -n '^206\.|^219\.|^222\.' design_docs/world-mission.md`

**Observed.** Lines 1008 R1 mints and Human.Approve/key, R2 formats, R3 key admin, R4 publish, R5 order; 1013–1014 roadmap and placement; 1017 service placement; 1020 preset decisions; 1026 bar; 1028 fake-only and attended; 1034 A4; 1037 finish A. Queue lines 6294/6307/6310. Exact probe `rg -n '^\| D-WORLD-(118|119|120) \|' design_docs/world-mission.md` yielded only :1037 D-WORLD-118, the positive control.

### W03

**Commands.** `rg -n 'exec\.Command(Context)?\(|func \(w \*workspaceTools\) execHandler' host/{daemon,broker} --glob '*.go' --glob '!**/*_test.go'`; `rg -n 'execHandler|shellHandler|func NewExec|exec.CommandContext' D/workspace.go host/broker/handlers.go`

**Observed.** No direct exec call in `host/daemon` in this scan, positive controls `workspace.go:460,478` and `handlers.go:159 cmd := exec.CommandContext(…)`.

### W04

**Commands.** `nl -ba F/service.go | sed -n '20,55p;85,94p;166,198p;234,254p;419,446p'`; `nl -ba F/doctor_service.go | sed -n '1,95p'`; service sources `rg -n 'LaunchAgents|gui/|launchctl|func ' F/service*.go F/doctor_service.go`

**Observed.** `service.go:29` label and one `serve` argv at 435; paths 48–51; `gui/<uid>` 85–88; typed `confirm` 167–179; `renderServicePlist` 186–196, bootstrap/build check 234–253. `service_launchctl.go:49,59,63,67,71` controller; `service_upgrade.go:125,243` upgrade and rollback; best-effort doctor 12–48.

### W05

**Commands.** `/bin/hostname -s`; `/usr/bin/id -un`; `/usr/bin/id -u`; `command -v gh git herdr`; filtered `/bin/launchctl print gui/$(id -u)` and `/bin/launchctl print gui/$(id -u)/com.sunholo.ailang-world.worldd`

**Observed.** `Voights-Mac-Studio`, `voightkampff`, uid 501. Both launchctl calls rc 0; GUI line lists World; job `state = running`, `pid = 49644`, installed binary `~/.ailang-world/bin/ailang-worldd`. Shell tools `/opt/homebrew/bin/{gh,git}`, `~/.local/bin/herdr`.

### W06

**Commands.** For keys `Label ProgramArguments.0 ProgramArguments.1 EnvironmentVariables.PATH RunAtLoad KeepAlive`, `/usr/bin/plutil -extract "$k" raw -o - "$HOME/Library/LaunchAgents/com.sunholo.ailang-world.worldd.plist"`; `/bin/ps -p 49644 -o uid=,comm=`

**Observed.** Expected label and binary, `serve`, PATH `/usr/bin:/bin:/usr/sbin:/sbin`, RunAtLoad/KeepAlive true. `ps: operation not permitted` **UNINFORMATIVE UNDER SANDBOX**, E1.

### W07

**Commands.** `/usr/bin/env -i HOME="$HOME" PATH=/usr/bin:/bin:/usr/sbin:/sbin /bin/sh -c 'for x in git gh ssh herdr; do command -v "$x" || printf "%s: missing\\n" "$x"; done'`; for system and Homebrew-prefixed PATH, `/usr/bin/env -i HOME="$HOME" PATH="$runner_path" GH_PROMPT_DISABLED=1 /bin/sh -c 'gh auth status --hostname github.com'`, filtered

**Observed.** System PATH finds `/usr/bin/git`, `/usr/bin/ssh`; misses gh and herdr. Auth probe rc 127 with system PATH, rc 1 with Homebrew PATH; absolute gh invocation printed filtered failed login/stored token invalid markers. Authentication parity is **UNVERIFIED** E2; command output is not proof of invalid rig credentials under unrestricted egress.

### W08

**Commands.** `rg -n 'exec|Socket|pane (create|split)|Command|preset|fallback|episode start' P/bin/world-agent-popup`; `rg -n 'world_herdr|stty|HERDR_' P/bin/world-agent-popup`; `rg -n 'herdr|plugin' P/herdr-plugin.toml`; `rg -n 'func |StateDir' P/lib/world.sh`

**Observed.** Popup plain fallback 75–93/140; attended start 107–125; plugin state 120,146,163; shell launcher 155, pane run/metadata 156–158; manifest min herdr 0.9.3 at 4. `lib/world.sh:50–68,125–130` environment and pane assumptions.

### W09

**Commands.** `nl -ba F/episode.go | sed -n '190,377p'`

**Observed.** Flags 195–217; shared parse 218; `--request` with resume/launch/headless or non-`se-tools` preset refuses 240–250; launch preset branch 263–265; request route 357–358 and typed provision 364–370.

### W10

**Commands.** `nl -ba F/mint_request.go | sed -n '23,155p'`; `rg -n 'awaitSurfaceApproval|mcpAfterToken|MintSession|consume' F/approve*.go`; `rg -n '^func |^type ' host/approval/request.go`

**Observed.** `buildMintRequest` 28 and fields 56–62, `fromRequest` 68, `mcpAfterToken` 92, put 129–140, `renderRequest` at 156 in `rg`; approve await call 324, request consume setup 394–396. Existing `MintRequest` at request.go:33 and `SignedMessage` :211.

### W11

**Commands.** `nl -ba F/session_new.go | sed -n '405,581p'`; `rg -n '^(func|type) .*|MintForRequest' host/authority/mint.go`; `nl -ba host/store/credentials.go | sed -n '445,503p'`; `rg -n 'ConsumeMint|func ' host/store/mint_requests.go`; `rg -n 'writeTokenFile' F/session_new.go`

**Observed.** `provisionEpisode` 405, signed fence 409–416, worktree 494–502, request/start selection 519/524, token/files 529–544, rollback 546–574, `writeTokenFile` 663. `MintStarted` and `MintForRequest` authority/mint.go:50/112. Store `MintSessionWithStart` :446 and session + `episode_starts` transaction :474–495; request consume separate at mint_requests.go:309.

### W12

**Commands.** `nl -ba F/episode_preset.go | sed -n '119,195p;218,308p'`; `rg -n 'ReadFile|sha256|resolvePreset|func .*Preset' F/episode_preset.go F/preset*.go`; `rg -n 'func writeMCPConfig|func runEpisodeStart' F/episode.go`

**Observed.** Validate binary/mode/output 119–134; branch checks 147–171; copy path 173 and task read :208; `req.startPoint` and start fields 239–243; prompt 245–268; health 269–274; task write via `O_EXCL`, MCP config 276–299; provision/launch 301–305. `writeMCPConfig` episode.go:540.

### W13

**Commands.** `rg -n 'func |\.Launch|Flock|agentrun.Run|launchLock|Credential|readFile|Through|pinned' F/episode_launch.go`; `rg -n 'type EpisodeStart|func .*EpisodeStart|type .*aunch|func .*aunch' host/store/episode* F/episode_launch*.go`

**Observed.** `launchMode` 44, `launchBinary` 151, launch 258, pinned bytes 292–300, TTY branch 300–322, task digest 324–327, lock 383, credential reads 424/430/454, interactive/headless 517/519, launch receipt appended 560–562. `EpisodeStart` :20, `EpisodeLaunch` :92, `AppendEpisodeLaunch` :112 in episode_starts.go.

### W14

**Commands.** `nl -ba host/agentrun/agentrun.go | sed -n '22,99p'`; `rg -n 'Run\(|process|kill|Start\(' host/agentrun/agentrun.go`; `rg -n 'func .*|[Tt]railer|environ' host/worktree/{commit.go,worktree.go,push.go}`

**Observed.** Headless `Run` 26 and group 49–51; interactive `RunInteractive` 32 and TTY 33–36; bounded context 44, `cmd.Start` 68, `procbound.Wait` 92. `worktree.Env` :49 uses scrubber, `Add` :63, `Output` :83, `runGit` :135.

### W15

**Commands.** `rg -n 'func |verb_requests|credentialsSchemaVersion' host/store/*{verb,credential}*.go`; `nl -ba host/store/verb_requests.go | sed -n '23,53p;129,150p'`; `rg -n 'func |^type ' host/approval/verb.go host/authority/judge_verb.go`

**Observed.** `credentials.go:45` is 9. DDL 24–39 only approved/superseded outcomes, positive control. `PutVerbRequest` :121, supersedes before prune :136–140, cap :143–147, `ConsumeVerbRequest` :239 and update :281. `VerbRequest` verb.go:51 and canonical :120. Judge SSH/passkey/recheck judge_verb.go:63/82/107.

### W16

**Commands.** `nl -ba F/episode_finish.go | sed -n '63,140p;360,495p'`; `rg -n '^func |openTerminal|Prepare\(|Unsandboxed|finishStep|Dispatch\(' F/episode_finish.go`

**Observed.** Flags 67–76, terminal fence 89–94, dry-run refusal 95–97; fresh unsandboxed pre-check 263; checks at 279–307, commit 361, body 389, prepare 400, typed path 432–439, remote/head/body checks 445–451. `finishSandboxCheck` :650.

### W17

**Commands.** Same numbered finish slices as W16; `rg -n 'func .*FinishDispatch|gh|ls-remote|func .*Prepare|func .*Dispatch' F/episode_finish_dispatch.go host/worktree`

**Observed.** Scope populated 404–408; reuse precedes surface branch 418–423; V3 confirmation 422–430, dispatch 453, record 460–469, detach/revoke 470–475, PR URL 494. `Prepare` :51 uses `exec.LookPath("gh")` :53 and label query :57.

### W18

**Commands.** `rg -n 'func |verbLandGrace|RecheckVerb|approved' F/episode_finish_surface*.go`; `sed -n '215,240p' F/episode_finish_surface.go`; `nl -ba host/approval/verb.go | sed -n '25,120p'`

**Observed.** `VerbRequestLifetime` :25 ten minutes, only finish at :31–32. `verbLandGrace` :28 ten seconds. `canonicalFinishArgv` :53, `finishVerbRequest` :75, put :168, `awaitVerbConfirmation` :188, polling :214/239, `readBackFinishDecision` :252. Poll binds own digest/challenge/expiry.

### W19

**Commands.** `rg -n 'func |^const ' F/episode_finish_resume.go`; `nl -ba F/episode_finish_dispatch.go | sed -n '100,240p'`

**Observed.** `finishMode` :19, `finishResumeBody` :79, `finishReusableApproval` :118. Current `ViewPR` :100 reads number/url/draft/state keyed by branch and only accepts OPEN, positive control for missing exact head/base/title/body identity. Timeout classification :155–163, push/label/view/create :164–197.

### W20

**Commands.** `rg -n 'remote|label|GitHub' F/episode_finish_dispatch.go`; `rg -n 'func .*|[Tt]railer|environ' host/worktree/{commit.go,worktree.go,push.go}`

**Observed.** `RemoteHead` checks exact ref via `ls-remote` :67–68, `Ancestor` :82, `Push` :91, `CreateLabel` :95, create PR :123–138; repo parse :141–152. `CommitWithCommitter` commit.go:17, `Push` push.go:11.

### W21

**Commands.** `rg -n 'func .*record|Phase|Operator|time|approval' D/episode_finish.go`; `rg -n 'func |^const ' F/episode_finish_resume.go`

**Observed.** `consumeOperatorProof` :92, phase/proof :154–160, typed approve 166–178, record compares scope :192 and rejects consumed done approval :211; `finishDetach` :133 and `finishRevoke` :148.

### W22

**Commands.** `"$HOME/.local/bin/herdr" --version`; `herdr --help`, `pane --help`, `pane split --help`, `pane run --help`, `tab create --help`, `api --help`, `session --help`, `workspace create --help` (absolute binary)

**Observed.** `herdr 0.9.3`; split accepts pane/direction/cwd/env, tab accepts workspace/cwd/env, run takes `<PANE_ID> <COMMAND>...`, API snapshot, named sessions, config override. Help is not a socket/launch experiment.

### W23

**Commands.** `rg -n 'func (Lookup|Split|ParseInterleaved)|case "(launch|start|finish)"' host/verbs/*.go`; `nl -ba D/surface_verb.go | sed -n '12,74p'`; `rg -n '/surface/(verb.json|confirm|requests.json|approve)' D/daemon.go`; `sed -n '453,485p;494,511p' host/surfacepage/app/app.js`; `rg -n 'confirmLineHTML|Array.isArray|surface/decide|verb.json|terminal' host/surfacepage/app/app.js`

**Observed.** `split.go:35`, `lookup.go:44–48,67,90–96,129`, `flags.go:14`; GET route 1154, `VerbSchema` :14 and handler :38; app array guard 434, confirm regex 456–458, notifications 473–482 mint and decisions with no confirmation leg, verb outcomes :496/507–508, fetch :526.

### W24

**Commands.** `for c in c251132f 48b8b6b6 0077a7b0; do git diff --numstat "$c^1" "$c" -- cmd host docs app | awk '{i+=$1; d+=$2} END{print i,d,i+d}'; done`; `rg -n 'git diff --numstat|1,398|1.35' design_docs/planned/w-world-address-bar-verbs-v3{,s5-sprint-plan,s6-sprint-plan}.md`; `wc -l` over service/preset/launch/finish/store/confirm files; `rg -n 'UpgradeCredentials|CredentialsSchemaVersion' F D --glob '!**/*_test.go'`

**Observed.** V2b 1186+69=1255; V3 S5 1387+11=1398; S6 418+7=425. V3 r2 reports 1.35× baseline, S5 plan :207 explicitly 910→1398 (1.54×). `service.go` 737, preset 314, launch 678, finish 748, surface finish 283, dispatch 199, credentials 623, verb requests 293, confirm 215 and approve 311 lines. `D/credentials.go:59–77` read-version and serve upgrade; `host/store/credentials.go:344–345` refusal of newer schemas.

### W25

**Commands.** `rg -n 'ReadFile|sha256|resolvePreset|func .*Preset' F/episode_preset.go F/preset*.go`; `rg -n 'Checks .*json|Unsandboxed .*json|workspaceEpisodePattern' F/{preset,episode}.go host/store/episode_starts.go`; `rg -n '\bSol\b|gpt-6\.2-sol|func launchBinary|Model .*json' F P --glob '!**/*_test.go'`

**Observed.** `defaultPresetsDir` preset.go:114, strict preset JSON :121, `loadLaunchPreset` :327, named read :349, `presetDaemonCheck` :502, CLI owner :515. `checks` preset.go:72/85 and `StartCheck.Unsandboxed` episode_starts.go:63. Sol-specific names absent; positive control `episode_launch.go:151 func launchBinary` and `preset.go:67 Model`.

### W26

**Commands.** `nl -ba D/surface_approve.go | sed -n '55,89p'`; `rg -n 'Origin|sameOriginFetch|jsonMediaType|func .*Approval|handleSurfaceRequests' D/surface_approve*.go`; `rg -n 'func |humanMu|commitBudget|Origin|sameOriginFetch|jsonMediaType|RecordOperatorApproval|ConsumeVerbRequest' D/surface_confirm.go`; `rg -n 'surfaceCSP' D/surface_page.go`; `rg -n 'func .*Test.*Chrome' D/surface_*chrome* F/*chrome*`

**Observed.** `approvalOrigins` :58, wrapper :76–85, POST Origin/fetch/JSON :218–226; confirm bounded body :119 and writer :129–131, world first :199 and consume :210. CSP :19–23. `confirm_surface_chrome_test.go:52,97` real route virtual passkey fixtures.

### W27

**Commands.** `sed -n '479,600p;852,891p' design_docs/planned/w-world-address-bar-verbs-v3.md`; `sed -n '722,749p' design_docs/planned/w-world-address-bar-verbs-v3s6-sprint-plan.md`; `rg -n '^##.*(ervice|pisode|inish)|^###.*(ervice|pisode|inish)' docs/QUICKSTART.md`

**Observed.** V3 §2.7 and RV1–RV9, S6 Appendix C explicitly carries prune/DDL/RV1, no notifications and fake-clock risk. Quickstart service :23, episode start :232, finish V3 :1915–1925.

### W28

**Commands.** `nl -ba D/episode_finish.go | sed -n '179,270p;265,282p'`; `rg -n 'func .*SetRegistryHead|func .*GetRegistryHead|func .*Tx|func .*Commit.*Registry' host/store --glob '!**/*_test.go'`

**Observed.** `finishDecisionScope` :243, rejects expired at :276. Record :186–195 checks scope, :206–235 checks/sets used head then episode head. `store.go:802/832/884/937` head primitives, `durable.go:52/73` durable begin/finish.

### W29

**Commands.** `nl -ba host/verbs/lookup.go | sed -n '10,122p'`; `rg -n '^var workspaceEpisodePattern|^const maxPresetNameLen|^func .*RecordOperatorApproval|^func \(c \*Credentials\) RevokeSession' F host/{store,broker,verbs} --glob '!**/*_test.go'`; `nl -ba host/worktree/commit.go | sed -n '1,96p'`

**Observed.** Lookup word sets :36–39 and address fallback :119; no `world-publish` entry, with positive known-never `service` :37. `doctor.go:461` grammar, `preset.go:150` cap 51. `approve.go:722` record helper, `credentials.go:557` credential-specific revoke. Commit resets index :34, literal staging :37–40, disables hooks :44.

### W30

**Commands.** `sed -n '52,120p;207,294p;604,666p' design_docs/planned/w-world-address-bar-verbs.md`; `sed -n '1,39p;98,144p' design_docs/planned/w-world-address-bar-verbs-v3.md`; `rg -n '^### 2\.[0-9]' design_docs/planned/w-world-approvals-by-key-a3.md`.

**Observed.** Parent §0.3 V4 says 1–2 sprints and ≈2,000, its §10 names a host task file and steps 1–5, RK5 asks to measure GitHub auth, and Q4 prefers a runner outside the daemon. V3 §0.4 carries F1, bar start/finish and herdr/open check scope. A3 §2.1 distinguishes surface and minter, §2.5 is the common judge, §2.6 is CLI await.

### W31

**Commands.** `cat /tmp/w406-designer/inputs/w-phone-episode-workflow.md /tmp/w406-designer/inputs/w-human-agent-surface.md`.

**Observed.** Phone draft assigns the runner design to row 243 V4, asks for a measured runner and one headless Sol, separates PR creation and merge, and describes #604. Human-agent draft asks for stable context, visible blockers and evidence-backed outcome after disconnect. Local briefs supply product intent; their authority remains unratified.

### W32

**Commands.** `nl -ba F/episode_launch.go | sed -n '44,73p;151,187p;437,467p'`; `nl -ba F/preset.go | sed -n '632,675p'`; `"$HOME/.local/bin/herdr" pane read --help`; `command -v jq`.

**Observed.** `launchMode` defaults to interactive and an explicit headless flag selects the headless branch (:44–68). Launch uses a PID file with `O_EXCL` (:442–447); preflight refuses an existing PID file regardless of liveness (:168–184). Preset check resolves first then lists resolved sandbox/unsandboxed status (:632–661). herdr read takes `<PANE_ID>`. jq is `/usr/bin/jq`.

### Reproduction notes

W02's positive ID check is on ledger rows
rather than references elsewhere in the charter. Filtered runtime probes record command results,
not results inferred from the local interactive environment. Any runtime test or mutation
described elsewhere in this design is a planner/build gate and has no KILLED/SURVIVED verdict yet.
