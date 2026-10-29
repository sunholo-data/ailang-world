# w-surface-index-follow (queue row 220): the app follows merges and running episodes — r0

**Status:** r0, written in World as episode `w220d`, Codex, 2026-10-29. Design only; not quorum-reviewed. Roadmap: [w-all-agent-work-through-world](w-all-agent-work-through-world.md) §2 row 4 and §3 step 3, beside row 188's N2d. Rulings needed: §6 R1–R7, especially the local-versus-remote merge clock (R7). Friction logged: **46** in [the dogfood log](../verification/world-dogfood/friction.md). No implementation is blocked by this writing session. Host timing, remote synchronization and row 218's eventual record reader remain **verify** in §1. The operator files decisions as D-WORLD-NN; this document does not resolve them.

## §0 What is wrong today

The index is a command's output. `ailang-worldd surface-index` reads one repo, one charter and one revision. It defaults to `dev` and 200 merges. It resolves the revision, reads first-parent merges, and runs merge-base once per merge (`cmd/ailang-worldd/surface_index.go:492`). It builds and validates a pure projection, then writes a temporary file beside the output and renames it (`cmd/ailang-worldd/surface_index.go:531`). It refuses an output inside the source repo.

With `--rig --addr`, the command reads episode rounds and `trial/prs.json`. It fetches reports and resolves diffs and read origins. This is M0c, not a source of running states (`cmd/ailang-worldd/surface_index.go:306`; `host/surface/episodes.go:374`). AttachRig attaches only episodes whose PR already names a merged step. Other episodes become unattached IDs. The whole build has a five-minute budget; a report GET can take 30 seconds (`cmd/ailang-worldd/surface_index.go:110`; `cmd/ailang-worldd/surface_index.go:118`). Repeating that path on a timer cannot promise a one-minute result.

The daemon is less static than the problem statement might suggest. It rereads the configured file on every request. It bounds it at 4 MiB, unmarshals it and validates it (`host/daemon/surface_index.go:49`). A rename is visible on the next request without restarting. No configured file or an absent file is a positive empty state. Invalid bytes produce 503. There is no last-good cache today. The daemon uses json.Unmarshal here; the CLI's stricter Decode is a different path.

The home service merely forwards an optional absolute `--surface-index` path to serve (`cmd/ailang-worldd/service.go:409`). It does not rebuild it. A PR can therefore remain absent indefinitely.

Live sessions are already read from the credentials store, read-only. Only unexpired rows are listed, with a refusal to count a potentially incomplete 500-row result (`host/daemon/surface_page.go:176`). The model adds those sessions separately (`host/daemon/surface_model.go:81`). The page renders a Running section, but a session is not a timeline state (`host/surfacepage/app/app.js:489`). A valid state ID is a commit hash, and its kind must be shared or step (`host/surface/validate.go:57`). A running episode cannot satisfy that contract.

N2c already refreshes on committed World log entries. It fetches the model and selected stack before replacing the render (`host/surfacepage/app/app.js:806`). A git merge alone makes no such entry. Rebuilding the file alone would still leave an idle open app stale.

**Exit criterion:** a merged PR appears in the app within a minute with no manual command; a running episode appears as an open state. §3 and R7 define which source makes the merge observable.

## §1 Premise verification log

Measured means read in this episode's worktree. It does not mean a host experiment ran.

| # | premise | how checked | result |
|---|---|---|---|
| P1 | The charter assigns automatic rebuilds, open states and several repos to 220. | World profiled grep of `220. ` in `design_docs/world-mission.md`, row at line 6287; roadmap §2 and §3 | measured; item and exit criterion agree |
| P2 | Git history builds two states per merge from one revision. | `cmd/ailang-worldd/surface_index.go:492`; `host/surface/surface.go:255` | measured; first-parent merges, second-parent step, merge-base |
| P3 | Rig attachment is for already merged PR steps. | `host/surface/episodes.go:374`; `cmd/ailang-worldd/surface_index.go:306` | measured; unmatched episodes are not open states |
| P4 | Current timing is too loose for a minute. | `cmd/ailang-worldd/surface_index.go:110`; `cmd/ailang-worldd/surface_index.go:118` | measured constants: five minutes and 30 seconds |
| P5 | File replacement is atomic; daemon reads every request. | `cmd/ailang-worldd/surface_index.go:531`; `host/daemon/surface_index.go:49` | measured; validation failure refuses, rather than retaining last good |
| P6 | Existing tests cover flags, rig bounds, report errors and atomic writing. | Read `cmd/ailang-worldd/surface_index_test.go`, especially lines 17, 125, 167 and 301 | measured tests, not executed; tests explicitly do not run the git build |
| P7 | The pure schema is single-source and commit-only. | `host/surface/surface.go:66`; `host/surface/validate.go:57` | measured; no repo namespace or open kind |
| P8 | Service forwards a file path, not source configuration. | `cmd/ailang-worldd/service.go:296`; `cmd/ailang-worldd/service.go:409` | measured |
| P9 | Preset starts already record repo, base, branch, task, agent and model. | `cmd/ailang-worldd/episode_preset.go:205`; `host/store/episode_starts.go:18` | measured; row 217's design claim that nothing records starts is historical |
| P10 | A start record survives credential changes and has a read-only lookup. | `host/store/episode_starts.go:16`; `host/store/episode_starts.go:68` | measured; immutable mint history, not evidence that a process is alive |
| P11 | Sessions already supply a running list, not timeline states. | `host/daemon/surface_page.go:176`; `host/surfacepage/model.go:50` | measured; unexpired credential means session live |
| P12 | The browser follows World commits, not index generations. | `host/surfacepage/app/app.js:806`; `host/surfacepage/app/app.js:875` | measured; independent refresh trigger needed |
| P13 | Finish will record commit and PR before revoking sessions. | [row 218](w-episode-finish.md) §3.8; charter D-WORLD-100, read by grep | quoted design and resolved operator-route ruling; actual finish reader/API **verify** when 218 lands |
| P14 | Waiting/open work is a later data milestone, beside N2d. | [sprint plan](w-world-os-surface-sprint-plan.md), Data milestones after S1; roadmap §3 step 3 | measured text; it does not authorize inventing approval waits |
| P15 | Four configured repos can rebuild within 20 seconds on the home host. | Requires fixture benchmark and attended load check, unavailable through this grep/ls profile | **verify**; proposed supported bound, not a measured performance claim |
| P16 | A GitHub merge reaches the configured local branch within the exit budget. | No remote or host sync command available in this session | **verify**; no such synchronizer found in the builder; R7 must settle the clock |

## §2 Conflict surface

| row | overlapping files or contract | boundary for 220 |
|---|---|---|
| 188, N2d | `host/surface/surface.go`, `episodes.go`, `validate.go`; `host/daemon/surface_index.go`, `surface_model.go`, `surface_page.go`; `host/surfacepage/model.go`, `layout.go`, stack code and `app/app.js` | Agree the v2 data contract first. 220 owns source following, generation status and open episode facts. N2d consumes those facts and owns the richer waiting/stack presentation. Ship the minimal open glyph/card and refresh wiring together so the exit run is visible. Do not run competing schema edits. |
| 217 | `cmd/ailang-worldd/episode.go`, `episode_preset.go`; `host/store/episode_starts.go` | Read the existing immutable start record. Add only a bounded read query if needed. No start schema change, launch work or new agent call. Older sessions without starts remain explicit unknowns. |
| 218 | Future finish storage/reader and episode-report finish block; `host/surface/episodes.go` and daemon model integration | Consume the operator-authenticated finish record after its transaction. A finish notification is an optional wake-up, not a correctness dependency. Do not implement commit, push, PR creation, checks or revocation here. |
| 219 | `cmd/ailang-worldd/service.go`, its plist argv and record; serve configuration | Add one follow-config argument to the existing service installation path. Keep supervision and upgrades in 219. No second scheduled service and no launchd timer. |

This episode changes only this design, the roadmap's row-220 link, and the friction log. The file list above is a future implementation boundary.

## §3 The design

### Sources and configuration

Propose `serve --surface-follow ~/.ailang-world/surface-follow.json`. It is mutually exclusive with `--surface-index`. The manual single-repo command and v1 serving mode keep their current contract.

One operator-owned JSON file configures one surface. Paths below are examples; the loader expands a leading tilde itself.

```json
{
  "schema": "world/surface-follow-config/v1",
  "repos": [
    {"id": "world", "repo": "~/dev/ailang-world",
     "charter": "design_docs/world-mission.md", "branch": "dev", "limit": 200},
    {"id": "tools", "repo": "~/dev/tools",
     "charter": "design_docs/mission.md", "branch": "main", "limit": 200}
  ]
}
```

IDs are unique stable slugs. Repo paths are canonical main checkouts; duplicate canonical repos are refused. Charter paths are relative to that repo, with no escape. Branches name local branches under refs/heads, not arbitrary revision expressions. Read the charter blob at the pinned branch head. A working-copy edit must not silently change the projection. An absent charter refuses that repo's first build.

Keep config outside every workspace root, mode 0600, with an operator-owned directory. Refuse unknown keys, symlinks into workspaces and writable-by-agent locations. Canonical path checks must include parent directories. Changes take effect on explicit service restart; no agent-writable preset can enroll a repo. Service install gains `--surface-follow`, forwards its absolute path, and records it through the existing plist/record mechanism. No credentials, remote URL or source filesystem path is sent to the browser.

### Trigger and clock

A background scheduler inside serve probes each configured branch every 15 seconds. It starts immediately at startup. Polling handles merges made outside World, packed refs and restarts. Git hooks require installation and miss other writers. Filesystem watches are platform-specific and can miss packed-ref replacement. Neither is the correctness trigger.

Resolve the branch once per probe. A changed head invalidates that repo. This includes rewinds and force updates, not just new merges. Rebuild the bounded window from the new pinned hash. Read log and charter using that hash, not a second read of the moving branch. Cache merge-bases by immutable parent pair. Coalesce changes into one pending latest-head rebuild per repo.

A successful row-218 finish transaction may wake the episode overlay immediately. It is not a merge signal: opening a PR is not merging it. Reconcile starts, sessions and finish records every 15 seconds even when git is unchanged. No new agent-side call or heartbeat is required.

The supported local-source budget is 15 seconds to detect, 20 seconds to rebuild and publish, 15 seconds for an already open visible app to refetch: at most 50 seconds. The remaining ten seconds are margin. A poll already in flight does not defer the next probe by another whole interval. Tests use a fake clock and controlled process seam; no minute-long sleeps.

The page polls the model every 15 seconds while visible, independently of AG-UI. A generation token avoids re-rendering unchanged content. Resume from a hidden tab triggers an immediate fetch. N2c's log refresh remains useful but is not needed for a merge. A failed stream cannot disable the model poll. Fetch the model and selected stack from the same generation; if generations differ, retry once and preserve the previous coherent view. Keep pan, zoom, focus and address. Old deep links outside the configured window return a framed missing-selection reason.

**Remote boundary:** this process never fetches. A GitHub merge is invisible until some authorized synchronizer updates the local branch. The 50-second bound starts when that local ref advances. This is narrower than an unconditional reading of the charter exit criterion. R7 asks the operator to rule on it. Do not report the remote exit as passed without measuring the full remote-to-app interval and naming the synchronizer. Its delay is additional, not hidden inside “fresh.” The app says “local branch checked”; remote freshness is unknown.

### Execution and publication

Run the scheduler in serve, with git subprocesses behind a narrow read-only seam. Extract reusable collection from the CLI into a host package; keep host/surface pure. Do not recursively invoke the installed CLI, launch a shell, or run one child builder per HTTP request. Service supervision restarts serve; it does not own collection.

Allow only revision resolution, bounded log, merge-base and charter-object reads in the critical path. Disable optional locks, hooks, pagers and external helpers through the existing scrubbed worktree execution boundary. Never run fetch, checkout, reset, worktree add, config, commit or push. Never write source refs, index or working files. Output goes only to daemon-owned cache outside all source repos and workspace roots. Instrument the seam so tests assert the command list and unchanged refs/worktree.

Build an immutable candidate, validate it, enforce byte bounds, then serialize it. Write a same-directory temporary file, sync and rename, as the existing writer does. Only after success swap one in-memory snapshot pointer containing data, generation and source statuses. Every request captures one pointer. Failed candidates never replace accepted data. Restart loads the last validated cache as stale until sources are checked. Cache writes and status reads are projections, not new World log entries.

A successful repo build can replace its partition while another repo retains its last good partition. Publish the resulting whole envelope atomically. Each partition carries its own head, built time, checked time and last failure. The aggregate generation identifies that exact set plus the episode overlay. An older or canceled job cannot publish over a newer head. Recheck after publication if the source moved during collection.

A rebuild in progress shows the previous graph, “updating,” and its build time. A failed rebuild shows the previous graph with “stale,” last successful local check, cause class and next command. Do not equate a healthy daemon with fresh sources. No accepted partition yet means “building” or a per-repo refusal. With no usable data, model/index return 503 with a structured reason; partial usable data returns 200 with explicit degraded status. Static pages show the same status text. A status transition is visible even if the graph did not change.

### Versioned graph and several repos

Propose `world/surface-follow-index/v1`, a new envelope rather than loosening the commit-only v1 validator. It contains repo partitions, each with an unchanged validated `world/surface-index/v1` merge projection, and a separately validated open-episode collection. Model and stack schemas advance together for follow mode. Legacy mode remains v1. All follow routes use the same snapshot reader, including index.json used by the address parser.

Within each repo, keep commit IDs and parent rules unchanged. At the envelope/model boundary qualify every state, brane and link with repo ID. For example, a commit becomes `world:commit:<sha>`, an open episode `world:episode:<ep>`, and a brane `world:row-220`. Never pass a qualified ID to git. PR number alone is not an identity. Two repos with row 220, PR 7 or the same commit hash remain distinct.

Render one surface, grouped by repo in config order. Each repo has its own shared lane, branch label and branes. There is no fabricated shared ancestor or merge edge between repos. Short address aliases are allowed only when unique; ambiguity lists repo-qualified choices. Legacy URLs resolve only when they identify exactly one repo. No separate daemon or separate page installation per repo.

### Running episodes and closure

Read sessions and start records from this daemon's credentials store using a bounded read-only transaction. Deduplicate credentials by episode; any unexpired credential makes its session live. Join each episode to its immutable start record, then map canonical start.repo to configured repo. Show only whitelisted fields: episode ID, start time, base SHA, branch, sanitized task title, agent and model. Never serialize preset bytes, paths, remote URLs, checks or tokens.

An episode with a live session and a matching start becomes an open state. “Running” means “session live,” not a claim that an agent process is alive or making progress. Its stable ID is episode-based, not a fake commit hash. Anchor it to the recorded base when that base is in the retained graph. Otherwise show an explicit base-outside-window marker. Do not synthesize a merge or a landed count. Use the branch's row when present; otherwise use an episode lane. Several episodes on one branch retain separate IDs.

A live legacy session without a start stays in Running with “start record unavailable”; it cannot be assigned a guessed repo. Starts from unconfigured repos are counted separately as “repo not configured.” Failure to read sessions is “running status unavailable,” not zero running. Paginate a new read query if needed; never reuse the current 500-row result as a complete join. If the overlay cap is exceeded, show the cap and unknown remainder.

Once row 218 lands, a trusted successful finish makes the state “PR open” even after revoke. A finish with push-failed, pr-failed or ambiguous outcome shows that outcome; it does not imply a PR exists. Use the operator-authenticated record per D-WORLD-100 R6. A PR URL is displayed only after its scheme/host and repo identity are checked. No GitHub call is made by this follower. Remote closure without a merge remains unknown.

On expiry or revoke without a finish, retain the open state as “session ended; outcome unknown.” Keep such states for seven days after the last session expiry or observed end, subject to the cap. Keep recorded PR-open states until their merge is observed, subject to the same explicit retention bounds. A successful merge in the same repo links the episode to its landed step using finish commit/PR evidence, and removes the open glyph in that same publication. Preserve an episode-ID alias to that landed selection while it is retained. If the merge is already known while the session is still live, show the landed state and live-session badge, not duplicate work.

Rig data is optional enrichment. Keep AttachRig's existing merged-step attribution and sanitization. Cache enrichment by source hashes and report cursor. It runs after core publication and cannot hold the minute gate. If enrichment fails, keep previous facts with their own age, or show unavailable. Do not copy old facts onto a different merge head.

“PR open” is waiting for merge; it does not mean an approval is waiting on Mark. The broader N2d “needs you” view still needs actual approval/mint records. This row supplies open states and honest lifecycle labels, not invented waits.

### Bounds

Support at most four repos in r0 implementation. Default limit is 200 merges per repo. Keep the CLI's 1..2500 per-repo limit, but refuse a follow config whose summed limits exceed 2000: 4000 merged states leave 1000 slots for open states within the aggregate 5000-state ceiling. Cap open episodes at 1000 and total serialized envelope and model each at 4 MiB. Refuse oversized candidates; do not silently truncate a merge chain or auto-lower limit. Tell the operator which limit to lower. Show each repo's retained window.

Use at most two collection workers, one job per repo, with per-git-call deadline 2 seconds and a 20-second critical generation deadline including queueing. At rest there are at most 16 git probes per minute for four repos. Cache hits require no merge-base call. A continuously changing source gets at most four log/charter collections per minute per repo. Cap all git launches, including enrichment, at 1200 per rolling minute across the daemon; core work has priority. Default cold work for four 200-merge repos is about 812 calls (probe, log, charter and merge-bases); larger limits may miss the deadline and must say so. They are bounds, not a performance guarantee.

Limit git subprocess output before parsing; use the 4 MiB aggregate publication cap and a bounded intermediate log buffer. Reject an overlong record explicitly. On a budget miss keep last good, record stale, and retry next poll without overlapping or an unbounded backlog. Retry failures at 15, 30 then 60 seconds; a changed head or explicit restart resets backoff. CPU use is bounded by worker count, deadlines, capped input and calls. There is no model call. A host that cannot meet the default timing fails the gate; do not increase clocks to pass it.

## §4 Failure modes (three-part, row 227)

Commands using the proposed flag are future commands. They must enter help before implementation gates use them. Browser reasons name repo IDs, not expanded source paths.

| what failed | why | the next command |
|---|---|---|
| Follow did not start | Config is absent or invalid; name field and repo ID. No sources were enrolled. | `ailang-worldd serve --help`; correct the named config field, then `ailang-worldd service start` |
| Config is refused | A path resolves inside a workspace or an agent-writable directory. No source read begins. | `ailang-worldd service install --help`; place the config outside workspaces and reinstall with `--surface-follow` |
| Repo cannot rebuild | Local branch or charter blob is missing. Keep its old partition stale. | `git -C <repo> show <branch>:<charter>` |
| Git collection timed out | Call or generation budget expired. Show limit, elapsed time and last successful head. | `ailang-worldd service status`; inspect host load before restarting |
| New PR merge is absent | The checked local head has not changed; remote freshness is unknown. | `git -C <repo> log -1 <branch>`; operator checks the external synchronizer |
| New candidate is refused | Validation failed or 4 MiB/5000-state/call bound was exceeded. Last good remains. | `ailang-worldd service status`; lower the named repo's limit in the config, then `ailang-worldd service stop` and `ailang-worldd service start` |
| Running status is unavailable | Credentials read failed or the bounded query was incomplete. Do not display zero. | `ailang-worldd session list --help`; inspect the named credentials store using the operator CLI |
| Episode has no open glyph | Legacy start metadata is absent or its repo is not configured. Running list states which. | `ailang-worldd episode report <ep> --json` |
| PR evidence is inconsistent | Finish head/repo does not match the indexed step. Do not attach by PR number alone. | `ailang-worldd episode report <ep> --json`; operator checks the recorded branch and commit |
| Cache cannot publish | Temporary write, sync or rename failed. Memory and disk retain last accepted generation. | `ailang-worldd service status`; inspect the service log and cache directory permissions |
| Browser still shows old data | Poll failed or model/stack generations disagreed. Preserve the coherent view and show its age. | `ailang-worldd service status`; reload `/surface` after the named daemon or network fault is fixed |

## §5 What this design does not do

It does not merge PRs, fetch remotes, mint credentials, launch agents, run finish checks, reclaim worktrees or change source repositories. It does not give the service GitHub secrets. It does not infer process liveness from a credential or infer operator approval waits from a PR.

It does not replace World's event store with git, rebuild all history without bounds, or make live sessions dependent on trial.py. The manual builder remains useful for fixtures and offline inspection. Several daemon stores and cross-machine episode federation are outside this row.

## §6 Rulings (proposed; recorded as D-WORLD-NN when made)

- **R1, trigger. (A, default)** Poll local refs and episode facts every 15 seconds; optional finish wake-up. **(B)** Install git hooks and filesystem watchers. A covers outside writers and restart recovery uniformly; B can reduce latency but still needs reconciliation.
- **R2, execution site. (A, default)** A bounded scheduler inside serve, with read-only git children. **(B)** A supervised service timer invoking the CLI. A shares store reads and atomic publication; B isolates collection but adds process/version coordination and retains the current builder's excessive critical path.
- **R3, schema. (A, default)** Versioned multi-repo envelope containing validated v1 merge partitions and separate open states. **(B)** Add repo and open kinds directly to v1. A preserves existing commit invariants and makes migration explicit; B is smaller initially but breaks old validators and address consumers silently.
- **R4, stale policy. (A, default)** Keep last validated partitions with age and failure status; 503 only when no usable data exists. **(B)** Refuse the whole surface on any repo failure. A keeps other repos and historical inspection useful; B makes freshness failure harder to overlook but hides valid data.
- **R5, episode source. (A, default)** Join live sessions to immutable starts; later consume trusted finish records, with no agent heartbeat. **(B)** Require an agent lifecycle tool for open/close. A works for existing preset starts and cannot be omitted by the agent; B could measure process intent but adds an unreliable agent obligation.
- **R6, several repos and bounds. (A, default)** One grouped surface, four repos, namespaced IDs, 4 MiB and 5000 aggregate states, with explicit window limits. **(B)** One surface per repo. A meets the row's one-surface item; B preserves more existing layout code but fragments navigation.
- **R7, remote clock. (A, default)** Guarantee the minute from advancement of the configured local branch, and name remote synchronization as an external prerequisite. **(B)** Require an independently authorized remote synchronizer and gate its delay within the same minute from GitHub merge. A keeps this row strictly read-only but needs an explicit charter interpretation; B meets the broader remote reading, adds another owner and must budget synchronization before claiming the exit.

## §7 Open questions

1. Does the operator accept R7 A, or must the attended exit start at GitHub's merge time? Until ruled, a local-only run cannot establish the remote claim.
2. What exact bounded finish-record read API will row 218 expose? Use its authenticated storage, not an untrusted workspace JSON file. Re-measure once implemented.
3. Can the home host sustain the 20-second generation budget on four default windows during concurrent episodes? P15 is a gate, not an assumption to waive.
4. Should ended unknown states be retained longer than seven days? Proposed default bounds history; PR-open states also need an explicit “retention limit reached” count when dropped.
5. N2d is named in the roadmap, while the measured sprint plan lists the waiting/open milestone without a dedicated N2d table row. Who owns the richer waiting cards? Default: 188, after the shared schema lands.
6. If several start records claim one PR across fix-up episodes, should the landed card show all rounds? Default: retain episode aliases, keep AttachRig's present attribution until 188 agrees a separate change.

## §8 Milestones

**M0/M1 built (w220a):** D-WORLD-103 all defaults; strict follow envelope with absent overlay, read-only pinned collector/config and bounded sink-only scheduler; fake-clock/git and real-git source-preservation tests. M2 onward remains unbuilt.

| M | what | gate |
|---|---|---|
| **M0** | Re-measure premises; settle R1–R7 and the shared N2d schema boundary | Operator records rulings; Go fixture tests preserve legacy v1 decode/validate and reject malformed envelope namespaces |
| **M1** | Read-only source collector, config, pinned-head collection and bounded scheduler | Go tests with fake clock/git: unchanged probe counts, moving refs, rewind, cancellation, deadlines, call budget and no git writes |
| **M2** | Atomic cache and snapshot reader, per-repo status and service wiring | Go tests for rename/write failure, restart recovery, no partial generation, stale partition retention and plist argv mutual exclusion |
| **M3** | Session/start overlay, ended states and finish integration seam | Go tests for deduplication, expiry, missing starts, unconfigured repo, cap/refusal and finish-before-revoke; operator verifies actual 218 reader |
| **M4** | Namespaced multi-repo model, minimal open cards, independent visible-page polling and generation-safe stack refresh | Go tests for equal PR/row/hash collisions and model/stack generation handling; operator opens two repos, disables AG-UI and observes refresh with navigation retained |
| **M5** | Bounds and loaded-host measurement | Go tests for all caps and last-good preservation; operator times four default cold windows and a changed head under concurrent session load |
| **M6** | Attended exit run | Operator merges a PR and measures app appearance within 60 seconds without surface-index or reload; starts an episode with no agent lifecycle call and sees its open state. Repeat across two repos. Record source clock, synchronizer if needed by R7, elapsed times and evidence |

## §9 Acceptance criteria (no model call; an operator or CI can check each)

1. An unchanged four-repo fixture over one minute issues at most 16 git probes and no history rebuild after initialization.
2. A local branch merge just after a poll appears in an already visible app within 60 seconds without a manual command or World log commit. Record detection, publication and render times.
3. Advancing a source during collection produces a coherent pinned-head index and schedules its successor. A canceled older build cannot overwrite that successor.
4. The collector's recorded argv contain only admitted read commands. Tests compare source refs, git index and worktree bytes before and after; all remain unchanged.
5. A successful publication returns either the old or new whole generation to concurrent readers. Injected write, sync, rename and validation failures preserve last-good data and expose a stale reason.
6. A newly minted preset episode with a live credential appears as an open state within 30 seconds. Its base, task title, agent and model match the immutable start. No agent-side call is made.
7. Two credentials for one episode produce one glyph. Expiry/revoke without a finish changes its label to outcome unknown. Missing start metadata produces an explicit Running-list reason.
8. A trusted finish changes a live glyph to PR open before revoke removes the credential. A failure outcome never becomes PR open. A matching merge replaces the open glyph with an episode alias to the landed step.
9. Two configured repos with the same row, PR number and commit hash render distinct selections and edges. Ambiguous short addresses fail with qualified choices.
10. Invalid config, duplicate canonical repo, unsafe config path, missing branch and missing charter each return a cause and next command; no source write occurs.
11. Exceeding repo, summed limit, open episode, byte, state or git-call caps rejects the candidate or marks the bounded overlay incomplete explicitly. The page never claims a complete zero-running result after a capped query.
12. Failure of one repo leaves another repo's latest valid partition visible. Its source status and the failed repo's last successful local check remain separately visible.
13. With AG-UI disconnected, visible-page polling still refreshes merges. Hidden-to-visible transition fetches immediately. A mismatched model/stack generation preserves the prior view until a coherent pair arrives.
14. Serializing every follow route emits no tokens, preset bytes, personal email addresses, source filesystem paths or remote credential URLs. Static pages carry the same freshness/refusal meaning.
15. Legacy `--surface-index` mode and the manual builder still accept their existing v1 fixtures. Supplying it together with `--surface-follow` refuses before serve starts.
16. The attended M6 evidence states whether timing began at a local ref update or remote merge. If R7 B is chosen, remote merge to app is at most 60 seconds including sync; local-only timing cannot pass that gate.

## Sources

- [Roadmap](w-all-agent-work-through-world.md) §2 row 4 and §3 step 3.
- [Mission charter](../world-mission.md), row 220; D-WORLD-98 and D-WORLD-100.
- [One-command episode design](w-episode-one-command.md), shape and voice; §3.4 immutable start.
- [Finish design](w-episode-finish.md) §3.8 and R6; proposal versus current code is distinguished above.
- [Home service design](w-worldd-home-service.md), row-220 boundary and service argv.
- [Surface sprint plan](w-world-os-surface-sprint-plan.md), M0a/M0c/M1, N2c and Data milestones after S1.
- Measured implementation: `cmd/ailang-worldd/surface_index.go` and its test; `service.go`, `episode.go`, `episode_preset.go`; `host/store/episode_starts.go`; `host/surface/surface.go`, `episodes.go`, `validate.go`; `host/daemon/surface_index.go`, `surface_model.go`, `surface_page.go`; `host/surfacepage/model.go`, `layout.go`, `app/app.js`.
