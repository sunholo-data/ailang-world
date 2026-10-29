# w-episode-finish (queue row 218): finishing an episode is one attended command — r0

**Status:** r0, written attended 2026-10-28 through World itself (episode `h2610280901` on the
trial rig; Mark: *"Write the design doc for row 218 (w-episode-finish) … per
w-all-agent-work-through-world §3 step 3. Cover: the preset's declared checks, the commit under
World provenance, push and PR as one outward effect Mark approves, the row 198 report as the PR
body, recording the PR on the episode, and closing the session. List open rulings as R1…Rn with
defaults. Log friction as you go."*). It is a design, not a build: no code exists yet. It is not
quorum-reviewed. The pick-time designer re-measures every premise marked **verify** in §1 (they
sit on GitHub, in `gh`, or on the Mac, none of which this session can reach) and runs the quorum
as usual. Roadmap: [`w-all-agent-work-through-world`](w-all-agent-work-through-world.md) §2 row
3. The charter row cites §3 step 3, where 218 is *used* (row 188's N2d onward "now run through 217
and 218"). It is *built* at the end of step 2 (`219 → 217 → 218`), after row 217's presets,
which it reads (§2). Rulings it needs: §6 R1–R8. Friction logged while writing: **42**.

## §0 What is wrong today

Row 218's own words: **"FINISHING AN EPISODE IS OPERATOR HANDWORK."** Today an episode ends like
this:

- **Commit, push and PR are by hand.** The rig's `trial.py` grades, records outcomes, banks the
  report and prints a draft-PR body, but "it never pushes and never calls `gh`: the operator
  opens the draft" (its D6). Attended episodes outside `trial.py` do all of it by hand. This
  session is an example: asked to "commit and push a branch, and tell me the PR", it could only
  hand the commit back (friction **38**, owned by this row).
- **The provenance trailer is pasted.** D-WORLD-60 asks for a `World-Provenance:` line in the PR
  and, "from row 138 on, a commit trailer". `trial.py` puts it in the PR body only, and no code
  puts it in a commit.
- **The report and the PR are two artefacts.** Row 198's `episode report` is World's own account
  of the episode, but a PR body is assembled separately from the trial row.
- **Nothing on the episode says how it ended.** The store holds every call, but no record of
  the commit, the PR or the session's end. The credential simply expires.

**Exit criterion (row 218):** an episode goes from start to merged PR with no `trial.py` and no
hand-written PR.

## §1 Premise verification log

Each premise is **measured** (read in this worktree's code on 2026-10-28, or observed in this
episode's calls), **quoted** (from a charter row, design or friction entry, not re-measured), or
**verify** (GitHub, `gh` or macOS behaviour this session can't reach: no network, no git).

| # | premise | source | status |
|---|---|---|---|
| P1 | Row 218's item: run the preset's declared checks, commit the worktree, push and open the PR with the episode report (row 198) as its body as an outward effect Mark approves, record the PR on the episode, close the session. The roadmap's "What doesn't change": every outward effect (push, PR, publish) needs Mark's approval, and a mint stays his act. | `design_docs/world-mission.md` row 218; roadmap §2 | measured |
| P2 | **Launch presets do not exist yet.** `--preset` today names a grant set, and the only one is `se-tools` (`unknown --preset %q (known: se-tools)`). The preset file with repo, profile, grants, TTL, agent and model is row 217's item, not built. **There is nowhere to declare checks today.** | `cmd/ailang-worldd/session_new.go`, `episode.go`; row 217 | measured |
| P3 | `trial.py`'s graders are the nearest thing to declared checks. Each is `{name, argv, cwd, timeout_s}` from the trial spec, run **unsandboxed** in the worktree after the agent stops (D5), in its own process group with a wall cap. Each is recorded as a `world/task-outcome/v1` with `ailang-worldd outcome record` under the episode's own session. Unsandboxed because the exec sandbox denies `/dev/ptmx` and binding a loopback port (friction 3, 4). | `scripts/trial/trial.py` (`run_grader`, `outcome_payload`, `record_outcome`); friction 3–4 | measured / quoted |
| P4 | D-WORLD-60 (Mark, 2026-10-03): a GitHub label `ailang-world` (colour `#e73c17`, with a fixed description), "created in a repo the first time it is needed"; a provenance line in the PR and, from row 138 on, a commit trailer `World-Provenance: store=<db> episode=<ep> entries=<from>-<to>`; and the banked evidence directory linked from the PR. | `design_docs/world-mission.md` D-WORLD-60 | measured |
| P5 | `ailang-worldd provenance` prints that trailer, with `store=` the db **path** from `/v1/health`, from an entry scan of at most `--scan` entries (default 500, **max 5000**). An episode longer than the scan window gets no trailer from it. | `cmd/ailang-worldd/provenance.go` | measured |
| P6 | Row 198's report is `GET /v1/episodes/<ep>/report`, also `episode report <ep> [--json]`. It holds `entries` (first and last, which needs no scan), `total_calls`, `calls_by_tool`, effect decisions, budgets, plan and handler refusals, `reads` (files read), `first_at`/`last_at`, `uncommitted_intents` and `notes`. **It has no list of files written and no markdown form.** It exits 1 for an episode with no committed call. | `host/daemon/episode_report.go`; `cmd/ailang-worldd/episode_report.go` | measured |
| P7 | `trial.py pr-body` refuses without a trailer, and refuses to print anything token-shaped: a bare 64-hex run, or `bearer`. | `trial.py` (`pr_body`) | measured |
| P8 | A broker `Git.Commit` handler exists (`git commit --no-edit -m <payload>`, empty `HOME`, identity `AILANG World <ailang-world@invalid>`), but **`NewGitHandler` has no production caller**: it is unwired. No push or PR effect exists anywhere in the repo. | `host/broker/handlers_git.go`; grep for `NewGitHandler` | measured |
| P9 | `host/worktree` is "World's one bounded, environment-scrubbed git site": `Add`, `Remove` and `Output` through one `runGit`, with GIT_* and registry variables scrubbed and a 30 s timeout. The AC10 census requires every production subprocess site to be driven and to scrub the registry credential. | `host/worktree/worktree.go` | measured |
| P10 | **`ailang-worldd commit` is not git.** It posts a World object to the session-gated `/v1/commit`. Every message this verb prints must say "git commit" where it means git. | `cmd/ailang-worldd/cli.go` (`runCommit`) | measured |
| P11 | The precedent for an irreversible outward effect is `Registry.Publish`. Its scope binds the exact version; its cost is 1, charged before dispatch; its authority is a **landed `ApprovalDecisionV1`** naming an `ApprovalRequestV1`, traversed by content address and **consumed exactly once** (`approval_claims` primary key). The bytes are recomputed before dispatch and must match the approval (`comparePublishHashes`). The credential is loaded only after every refusal; the child environment is built from nothing. Outcomes are classed **definite or ambiguous** ("a false REFUSAL costs a retry, a false ALLOW costs a permanent public artifact"). | `host/broker/registry_publish.go`, `approve.go` | measured |
| P12 | **An approval object is a single-use fence, not an attended one.** `world-publish`'s own fence comment says a headless loop can mint a decision by calling `broker.DecideApproval`. The attended fence is a **controlling terminal**: `/dev/tty` must open, and the typed line is read from it, never from a redirected stdin. `session new`/`episode start` use the same idea (`openTerminal` before any side effect, MUT-FENCE-BYPASS). | `cmd/world-publish/tty.go`; `cmd/ailang-worldd/session_new.go` | measured |
| P13 | `DecideApproval` is called in production only by `world-publish` (a store-direct CLI). The daemon has **no decision route**. A store-direct writer needs the world store's writer lock, so it can't run while the daemon serves (219's P8: `world-publish transitions` needs the daemon stopped). | grep for `DecideApproval`; `w-worldd-home-service` §1 P8 | measured / quoted |
| P14 | `outcome record` commits a `world/task-outcome/v1` through the shared one-object writer (`commitOneObject`) **under the episode's own session credential**, written by `harness:floor`. The agent holds that same credential, so anything committed this way could also have been committed by the agent. | `cmd/ailang-worldd/outcome.go` | measured |
| P15 | `session revoke --episode <ep>` deletes every credential row of the episode in the credentials store, under that store's own lock. It works while the daemon serves, needs no terminal fence, and the daemon refuses the token from its next request. | `cmd/ailang-worldd/session.go` (`sessionRevokeHelp`) | measured |
| P16 | The worktree is `git worktree add -b world/<ep>` from the main checkout's `HEAD` at mint, and nothing records that base (row 228's P10). | `host/worktree.Add`; `episode.go`; [`w-agent-search-and-read-cap`](w-agent-search-and-read-cap.md) §1 | measured |
| P17 | Every worktree-changing effect leaves a trace a finish can read. `Workspace.Write`/`Workspace.Delete` payloads carry the path. `Ailang.Check`, `Ailang.CLI` and `Workspace.Exec` results carry a `worktree_delta` (row 184): git-status basis, gitignored paths out of scope, at most 256 changes listed, and an `error` instead of a failure when the snapshot fails. | `host/broker/worktree_delta.go`, `handlers_ailang.go`, `handlers_exec.go` | measured |
| P18 | Row 219's service runs with `EnvironmentVariables` = `PATH` only, by design: no registry key and no other ambient secret. A daemon-side push would need a GitHub credential it doesn't have. | [`w-worldd-home-service`](w-worldd-home-service.md) §3.1 | quoted |
| P19 | `gh` authenticates from Mark's keyring at his terminal, and pushing over SSH needs his agent socket. Both are present in an interactive shell and absent from a LaunchAgent's minimal environment. | macOS, `gh` | **verify** |
| P20 | GitHub's PR body limit (65,536 characters as commonly documented), and whether each target repo's branch protection requires signed commits, a verified author email, or status checks. | GitHub | **verify** |

**P12–P14 shape the design.** The approval object alone doesn't prove Mark was there (P12).
The obvious write path for records is one the agent's own token can use (P14). And the
decision path that exists needs the daemon stopped (P13). So §3 puts the attended fence at the
terminal, as `world-publish` does, and treats World's records of the finish as **evidence
cross-checked against git and GitHub**, not as the gate (R6).

## §2 Conflict surface

- **Row 217 (presets).** 218 reads a preset's `checks` (§3.2) and its repo, base branch and
  remote. 217 owns the file and its schema, and 218 adds one section to it (R1). Per the roadmap order
  217 lands first. Until it does, finish takes `--check` flags for testing only.
- **Row 198 (report).** 218 adds two things to the report: a `writes` list (P17) and a markdown
  rendering (`episode report --markdown`). Both are additive. The JSON route stays as it is for
  `ledger.py`.
- **Row 201 (approvals by key)** replaces the typed `y` with a signed approval of the same
  request object (§3.6). 218 creates that object so 201 has something to sign.
- **Row 224 (standing allowances)** draws *sessions* without a prompt. An outward effect is not a
  session, and the roadmap keeps it Mark's act every time. Finish never draws from an
  allowance.
- **Row 228** gives the agent base/lag and read-only git. 218 owns every git *write*, and
  friction 38 is 218's.
- **Row 219 (home service)** decides the daemon address that finish talks to. A daemon-side push
  would conflict with 219's no-secrets plist (P18): R3.
- **Row 156 (protected test paths).** A commit that changes a protected path is named on the
  approval prompt and in the PR body (§3.4). 218 doesn't change the protected check.
- **Row 226 (scorecard)** reads the finish record (§3.8) for operator minutes and outcome.
- **Row 204 (GC)** reclaims worktrees. Finish keeps the worktree (R7).
- **Row 183's `trial.py`** stays for paired shell arms and graded trials. Its `pr-body` and the
  operator's hand PR are what finish replaces for World arms.
- **Row 196 (the loop through World)** stays gated. Finish is attended by construction (the
  terminal fence) and can't be called from the loop.

## §3 The design

### §3.1 The verb

```
ailang-worldd episode finish <episode> [--preset <name>] [--title <text>] [--ready] [--dry-run]
```

It runs at Mark's terminal and is resumable: each step reads what the previous run recorded
and skips what is done (§3.9). The order is fixed. Every step before §3.6's `y` is local and
reversible, and nothing leaves the machine before it.

1. **Fence:** a controlling terminal (P12's `requireControllingTerminal`, shared and not copied)
   before any side effect. No terminal means exit 1 and nothing written.
2. **Validate** (§3.2's inputs): the session is unexpired, since finish records under it until §3.8;
   the worktree is a linked git worktree, not degraded; the report's `uncommitted_intents` is 0;
   the branch is `world/<ep>`; and the worktree has changes against its base.
3. **Checks** (§3.2).
4. **Attribution** (§3.3).
5. **Git commit**, local (§3.4).
6. **Body:** render the PR body from the report (§3.5).
7. **Approve:** one prompt binding the exact bytes, answered `y` (§3.6).
8. **Push and open the PR**, one outward effect (§3.7).
9. **Record** the finish on the episode (§3.8).
10. **Close** the session (§3.8).

`--dry-run` does steps 1–6, prints the commit, the body and the approval scope, and then resets
the commit, so its only lasting trace is the check outcomes it recorded.

### §3.2 The preset's declared checks

A preset (row 217) gains one section, using `trial.py`'s grader shape (P3) so the trials and
finish agree:

```json
"checks": [
  {"name": "go-test", "argv": ["go", "test", "./..."], "cwd": ".", "timeout_s": 900},
  {"name": "ail", "argv": ["./scripts/verify_ail.sh"], "cwd": ".", "timeout_s": 600,
   "unsandboxed": true}
]
```

**Where they run is R2.** The default (A) runs each check **through the episode's exec sandbox**,
as a `workspace-exec` call on the episode's own session and profile. The check is then a recorded
`Workspace.Exec` effect like any agent call, and its `worktree_delta` shows what it wrote. A check
that can't run sandboxed (P3: a pty, a loopback port) says `"unsandboxed": true`. It is then named
on the terminal **before it runs** ("runs the agent's code outside the sandbox"), runs as
`trial.py`'s graders do, and needs a `y` of its own. Either way, each check's verdict is recorded
as a `world/task-outcome/v1` (P3, P14), with its entry ref kept for the body.

A failing check stops finish before the commit, unless R2's override is used. Nothing is
committed, and the refusal names the check, its exit code, the outcome ref and the stderr head.

### §3.3 Attribution: is every change World's?

"Commit under World provenance" has to mean more than adding a trailer. Finish compares the
worktree against its base (`git status` and `diff` through `host/worktree.Output`, P9) and puts
every changed path in one class:

| class | evidence |
|---|---|
| **written** | a `Workspace.Write`/`Workspace.Delete` record of this episode names the path (P17) |
| **produced** | a `worktree_delta` of this episode's `Ailang.Check`/`Ailang.CLI`/`Workspace.Exec` lists it |
| **unattributed** | neither: an operator edit, a check's artefact, a delta over its 256-change cap, or a delta that recorded an `error` |

The report gains this as `writes` (§2, row 198), so a reader sees it without running finish. What
happens to **unattributed** paths is R4. By default finish refuses and lists them, and the operator
either removes them or passes `--include-unattributed`. They are then committed and listed in
the body under **"Not made through World"**. They are never silently mixed in.

### §3.4 The git commit

One commit per finish, through `host/worktree` (P9), which gains a `Commit`. It is the one
bounded git site, and the unwired broker handler (P8) is retired (R8). It runs with
`-c core.fsmonitor=false -c core.hooksPath=/dev/null -c commit.gpgsign=<from R5>` and stages
exactly the attributed paths, plus the unattributed ones if R4 allows. The message is:

```
<title: --title, else the episode's task title, else "World episode <ep>">

<one paragraph: agent and model, calls, checks passed, files changed>

World-Provenance: store=<db path from /v1/health> episode=<ep> entries=<first>-<last>
```

The entries come from the report's `entries` (P6), not from `provenance`'s scan (P5), so a
long episode still gets its trailer. Who the author and committer are is R5. Protected paths
(row 156) that changed are listed in the paragraph.

### §3.5 The body is the report

`episode report <ep> --markdown` (new, §2) renders the row 198 report. Finish adds only what the
report can't know: the commit sha, the checks' outcome refs and the label note. The sections:

1. One line: `Episode <ep>: <agent> <model> working through AILANG World only.` and the D-WORLD-60
   label note.
2. **Checks:** `name` → PASS/FAIL, sandboxed or not, outcome ref and entry.
3. **What the agent did:** calls by tool, effects by decision, refusals with their first
   reasons, budgets used/left. These are the report's own numbers.
4. **Files:** written / produced / not made through World (§3.3), and files read (count).
5. `World entries <first>-<last>; walk any of them with ailang-worldd why <entry>.`
6. The trailer, the same line as the commit's.

`trial.py`'s scrub is applied to the whole body (P7). It is held under GitHub's limit (P20), with
truncated sections saying so and naming `episode report <ep>` for the full form.

### §3.6 The approval: one prompt, exact bytes

Finish builds one `world/approval-request/v1` (the landed `Human.Approve` object, P11) whose
scope binds every byte the outward effect will send. The scope is canonical
and fixed-order, like `PublishApprovalScope`:

```
github-pr:remote=<url> head=world/<ep> sha=<commit sha> base=<branch> title-sha256=<…>
  body-sha256=<…> label=ailang-world draft=<bool> create-label=<bool> expires=<unix>
```

The terminal shows it in readable form: the repository, the branch and sha, the base, the title,
the body's first lines and size, the label (and that it will be **created** if missing, a second
public write that D-WORLD-60 allows), and draft or ready. `y` lands an `ApprovalDecisionV1` for
it. Anything else exits 1 with the commit kept, so a later run resumes at this prompt.

Where the decision is written while the daemon serves (P13) is part of R6. Row 201 later lets
Mark sign this same request with a key instead of typing at a terminal.

### §3.7 Push and PR, one outward effect

Dispatch follows `Registry.Publish`'s order (P11): refusals first, then credentials, then the
network. **Who dispatches is R3.** The default (A) is the finish process itself, at Mark's
terminal, with his own `git`/`gh` credentials (P19), so the service never holds a GitHub secret
(P18). Immediately before dispatch it recomputes the sha and the body's sha256 and refuses on
any mismatch with the approval. Then:

1. `git push <remote> <sha>:refs/heads/world/<ep>`, **never forced**. If the remote branch
   exists at a different sha, it refuses unless that sha is an ancestor (a fix-up round, §3.9).
2. If the label is missing and the scope said `create-label=true`: `gh label create
   ailang-world --color e73c17 --description <D-WORLD-60's text>`.
3. `gh pr create --base <base> --head world/<ep> --title <title> --body-file <body>
   --label ailang-world [--draft]`. If `gh pr view world/<ep>` already shows an open PR, the
   existing PR is taken instead (idempotent).

**Outcome classes** (as `Registry.Publish`'s):

| class | meaning | what a re-run does |
|---|---|---|
| `done` | pushed and PR open | nothing; prints the PR |
| `push-failed` | definite: nothing left the machine, or the remote refused | retries from the push; the approval still binds the same bytes |
| `pr-failed` | the push landed, but the PR wasn't created | resumes at step 3; the push is not repeated |
| `ambiguous` | timed out mid-request | queries the remote (`ls-remote`, `gh pr view`) before doing anything |

The approval is consumed once (P11). A retry after `push-failed` or `pr-failed` under the same
bytes reuses it, since the effect never completed, and any change to the bytes needs a new `y`.

### §3.8 Record the PR, then close the session

**Record.** One `world/episode-finish/v1` object, committed while the session still exists:

```json
{"episode": "<ep>", "branch": "world/<ep>", "base": "<sha>", "commit": "<sha>",
 "remote": "<url>", "pr": {"url": "…", "number": <n>, "draft": true},
 "approval": "<decision ref>", "checks": ["<outcome ref>", …],
 "unattributed": ["<path>", …], "report_sha256": "<…>", "trailer": "World-Provenance: …",
 "outcome": "done|push-failed|pr-failed|ambiguous", "finished_at": "<RFC 3339>"}
```

`episode report` shows it (a `finish` block), and the surface index (row 220) can show the PR
on the episode. Because a session-gated write is one the agent could also make (P14), **a reader
trusts this record only when its `commit` exists on `branch` and its `pr.url` names that head**.
`ailang-worldd why` checks both when asked (R6).

**Close.** `session revoke --episode <ep>` semantics (P15): every credential of the episode is
deleted, and the agent's next call is refused as an unknown credential. A running agent pane
sees the context block's "session invalid or expired" line. Finish revokes only after the record
lands. If revoking fails, finish says so, exits 2, and the next run revokes.

### §3.9 Re-runs and fix-up rounds

Each run reads its state from git and the store, never from a local file. If `HEAD` already
carries this episode's trailer, the commit step is done. If a decision exists for the current
bytes, the prompt is done. If the remote has the sha and `gh pr view` shows the PR, dispatch is
done. A **fix-up round** (review comments) is a new episode started on the same branch, or the
same episode while its session lives. Its finish adds a commit on top, which pushes as a
fast-forward to the same PR after a new `y`. Whether a fix-up reuses the episode id is §7 Q1.

## §4 Failure modes (three-part, row 227)

| # | condition | first line → `next:` |
|---|---|---|
| F1 | no controlling terminal | `refusing: no controlling terminal` → run it in a terminal (an agent or CI can't finish an episode) |
| F2 | session expired before finish | `episode <ep>'s session has expired: finish records under it` → `ailang-worldd episode finish <ep> --renew` mints a short finish-only session after its own `y` (§7 Q3) |
| F3 | `uncommitted_intents` > 0 | `an effect of <ep> is still in flight` → wait, or `ailang-worldd episode report <ep> --follow` |
| F4 | a check failed | `check <name> failed (rc N), outcome <ref>` → fix it in a new round, or R2's override |
| F5 | unattributed paths | the list → remove them, or `--include-unattributed` (R4) |
| F6 | bytes changed after `y` | `the commit or body changed since approval` → re-run to approve the new bytes |
| F7 | remote branch diverged | `world/<ep> on <remote> is at <sha>, not an ancestor` → never forced; resolve by hand |
| F8 | `gh` not authenticated | `gh can't reach GitHub as you` → `gh auth login`, then re-run (resumes at dispatch) |
| F9 | ambiguous dispatch | `the push or PR may have happened` → re-run; it queries before acting |
| F10 | revoke failed | `PR recorded, session NOT closed` → re-run, or `ailang-worldd session revoke --episode <ep>` |

## §5 What this design does not do

It doesn't merge: merging stays Mark's on GitHub, and the exit criterion ends there. It doesn't
remove the worktree (R7). It doesn't give the agent any git write: finish is a terminal verb
the agent can't call (F1), so friction 38's answer is "the operator runs `episode finish`", not a
tool. And it never force-pushes.

## §6 Rulings (proposed; recorded as D-WORLD-NN when made)

- **R1, where checks are declared.** **(A, default)** a `checks` list in row 217's preset file,
  in `trial.py`'s grader shape (§3.2), with `--check` flags for tests only. **(B)** exec-profile
  commands marked `finish: true`, so no preset section, but a profile describes what an agent
  *may* run, not what an episode *must* pass.
- **R2, where checks run, and failing ones.** **(A, default)** in the episode's exec sandbox as
  recorded `Workspace.Exec`, with a per-check `unsandboxed: true` that is named and confirmed
  before it runs. A failing check stops finish, and `--allow-failing <name>` (typed name, listed in
  the body as FAILED) is the override. **(B)** all checks unsandboxed, as `trial.py` does, so no
  sandbox gaps, but finish then runs agent-written code on the host with no confirmation step.
- **R3, who dispatches push and PR.** **(A, default)** the finish process at Mark's terminal,
  with his own credentials, after the approval binds the bytes. **(B)** a brokered effect in the
  daemon, `GitHub.OpenPR`, in `Registry.Publish`'s exact shape. It is more uniform, but needs a
  standing GitHub secret in the service that 219 keeps out (P18).
- **R4, unattributed changes.** **(A, default)** refuse and list them, with
  `--include-unattributed` committing them under "Not made through World". **(B)** always
  commit and list, so there is no refusal, but provenance becomes a label rather than a check.
- **R5, commit identity.** **(A, default)** author `<agent> via AILANG World
  <ailang-world@invalid>` (the unwired handler's address, P8), committer Mark's own git identity
  and signing settings. The commit says that Mark committed what World's agent wrote. **(B)**
  Mark as both, with the agent named only in the trailer paragraph. Either is subject to P20's
  branch-protection facts.
- **R6, where the decision and the finish record are written while the daemon serves.**
  **(A, default)** a new **operator route** on the daemon, `POST /v1/episodes/<ep>/finish`,
  that refuses any session credential. It accepts only a request from the CLI carrying an
  operator proof minted at the terminal (row 201's key once it exists; until then, a one-time
  secret the CLI reads from a 0600 file beside the credentials store). The route lands the
  approval decision and the finish record, so the agent's token can write neither. **(B)** the
  session-gated `/v1/commit` (P14), with the records' integrity resting on the cross-check against
  git and GitHub (§3.8) alone. Nothing new is needed, but the agent could plant a look-alike
  record before revoke.
- **R7, the worktree after finish.** **(A, default)** kept until the PR merges or closes, then
  reclaimed by row 204. **(B)** removed at finish, so the disk is freed at once, but a fix-up
  round needs a new worktree.
- **R8, the git write site.** **(A, default)** `host/worktree` gains `Commit` (and `Push` under
  R3 A), and the unwired `Git.Commit` handler is deleted with its tests. **(B)** wire the broker
  handler for the commit only, so two git write sites exist.

## §7 Open questions

1. **Fix-up rounds' episode id.** Is a review round a new episode on the same branch (its own
   report, its own trailer range: two trailers in one PR), or the same episode with a
   renewed session? Default: a new episode, `--branch world/<first-ep>`, since each round's
   provenance is then whole.
2. **Evidence directory (D-WORLD-60 (3)).** Finish has no banked directory to link. Should it
   commit the report JSON beside the change, attach it to the PR, or link the surface's episode
   page (row 220)? Default: the surface link once 220 serves it, and nothing before that.
3. **An expired session at finish (F2).** A finish-only renewal is a mint, so a `y`. Should
   `episode start` mint with TTL headroom for finish instead? Default: renewal, since a TTL
   shouldn't be guessed at start.
4. **Repos whose base is not `dev`.** The base comes from the preset (row 217). With no preset,
   is it the remote's default branch (`gh repo view --json defaultBranchRef`), or refused?
   Default: refuse, naming the preset field.

## §8 Milestones

| M | what | gate |
|---|---|---|
| **M0** | row 198 additions: the report's `writes` (§3.3) and `--markdown` (§3.5) | Go tests over seeded episodes (writes, deltas over cap, delta errors); the body scrub test |
| **M1** | **M0/M1 built (w218a)**: fence, validation, pinned checks (R2), attribution, `host/worktree.Commit` with the trailer, body render, `--dry-run`; protected paths left to M2 | Go tests with a fixture repo and the fake terminal; ACs 1–6 |
| **M2** | the approval scope, the operator route (R6), dispatch through a `git`/`gh` seam (a fake in tests), outcome classes, the finish record, revoke | Go tests with the fake seam for every outcome class; ACs 7–12 |
| **M3** | re-runs and fix-up rounds (§3.9); delete the unwired handler (R8) | Go tests killing finish after each step and re-running; AC13 |
| **M4** | the exit run, attended: one episode from `episode start` to a merged PR, on a real repo | AC14 |

## §9 Acceptance criteria (no model call; an operator or CI can check each)

1. **The fence.** `episode finish ep1 </dev/null` under `setsid` (no controlling terminal) exits 1
   with F1, and the worktree's `git status --porcelain`, `HEAD`, the store head and the
   credentials store are all byte-identical before and after.
2. **Validation.** A degraded episode, an expired session, `uncommitted_intents` > 0 and an
   unchanged worktree each refuse with their F-line, and nothing is run.
3. **Checks recorded.** With two declared checks (one passing, one failing), the report shows a
   `Workspace.Exec` entry per sandboxed check and a `world/task-outcome/v1` per check, and finish
   stops with F4 before any commit (`HEAD` unchanged).
4. **Unsandboxed is confirmed first.** An `unsandboxed: true` check answered `n` never runs. A
   marker file the check would create doesn't exist.
5. **Attribution.** In a fixture where the agent wrote `a.go` through `ailang-write`, an exec
   call produced `gen.txt`, and the operator touched `b.go` by hand, the classes are exactly
   written `{a.go}`, produced `{gen.txt}`, unattributed `{b.go}`. Finish refuses with F5, and with
   `--include-unattributed` the body lists `b.go` under "Not made through World".
6. **The commit.** `git log -1 --format=%B` ends with exactly the trailer, whose `entries`
   equal the report's `entries`. The author and committer are R5's. `git show --stat` lists
   exactly the staged classes. An episode of 6,000 entries still gets its trailer (P5's scan
   limit doesn't apply).
7. **Nothing leaves before `y`.** With the dispatch seam recording calls, answering `n` makes
   zero seam calls, and the commit stays.
8. **Bytes bind.** Changing one byte of the body after `y` (a test hook) refuses with F6, with
   zero seam calls.
9. **One outward effect, in order.** Under `y`, the seam sees exactly push → (label create, if
   scoped) → PR create, with the approved sha, base, title, body and label, and `--draft` unless
   `--ready`.
10. **Outcome classes.** For each of `push-failed`, `pr-failed` and `ambiguous` (seam-injected),
    the record's `outcome` is that class, and the re-run performs only §3.7's listed remaining
    calls: no second push after `pr-failed`, and queries before any action after `ambiguous`.
11. **Record, then revoke.** After `done`, `episode report ep1 --json | jq .finish.pr.url` is the
    PR, and the next MCP call with the episode's token is refused as an unknown credential.
    Under R6 A, a `POST /v1/episodes/ep1/finish` carrying the episode's session token is
    refused.
12. **The body is the report.** The body's calls, refusals and budgets equal `episode report
    ep1 --json`'s. The scrub refuses a seeded bare 64-hex value. The trailer is the body's last
    line.
13. **Resumable.** Killing finish after each of steps 5, 7 and 8 and re-running ends in one
    commit, one push, one PR and one record: never two.
14. **The exit run** (attended, M4). An episode started with `episode start --preset <p>` goes
    to a merged PR, with `trial.py` not invoked and the PR body unedited by hand. The PR carries
    the `ailang-world` label, its head commit carries the trailer, and `why <first entry>`
    walks back from it.

## Sources

- Charter row 218 (`design_docs/world-mission.md` line 6281), and rows 156, 183, 196, 198, 201,
  204, 217, 219, 220, 224, 226, 227 and 228; ledger D-WORLD-60.
- [`w-all-agent-work-through-world`](w-all-agent-work-through-world.md) §2 (row 3, "What doesn't
  change") and §3 (steps 2 and 3); [`w-worldd-home-service`](w-worldd-home-service.md) §1 P8 and
  §3.1; [`w-agent-search-and-read-cap`](w-agent-search-and-read-cap.md) §1 P10 and P17.
- [`friction.md`](../verification/world-dogfood/friction.md) entries 3, 4, 6, 38, and 42 (logged
  while writing this).
- Code: `scripts/trial/trial.py`; `cmd/ailang-worldd/{episode,episode_report,provenance,outcome,session,session_new,cli,main}.go`;
  `cmd/world-publish/tty.go`; `host/broker/{handlers_git,registry_publish,approve,worktree_delta,handlers_ailang,handlers_exec}.go`;
  `host/daemon/episode_report.go`; `host/worktree/worktree.go`.
