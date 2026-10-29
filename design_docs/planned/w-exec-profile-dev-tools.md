# w-exec-profile-dev-tools (queue row 221): the agent runs its own browser and build checks — r0

**Status:** r0, written in World as episode `w221d`, Codex, 2026-10-29, against dev `bcb07394`. Design only; not quorum-reviewed. Roadmap: [w-all-agent-work-through-world](w-all-agent-work-through-world.md) §2 row 5 and §3 step 4, with row 199. Rulings needed: §6 R1–R7. D-140-4 stays settled; D-WORLD-84 supplies the long-run lane. The operator files the rulings. Friction logged: 50–51 in [the dogfood log](../verification/world-dogfood/friction.md). Writing is complete. Chrome, cargo, HOME compatibility, inner sandbox compatibility and host-service enforcement cannot be exercised through this episode's grep/ls profile; those measurements remain **verify**, with explicit milestone gates.

## §0 What is wrong today

The agent can change a browser test but cannot establish that it ran. `TestSurfaceLiveChrome` skips when `WORLD_CHROME` is absent (`host/daemon/surface_live_chrome_test.go:22`). Its test server uses `httptest.NewServer` in the same file at line 42. Supplying a Chrome path alone will not fix the denied bind.

Exec reads `world/exec-profile/v1`. Its strict top-level keys are `profile, project, root, path, env, read_roots, caches, timeout_ms, commands, probe, protected` (`host/broker/exec_profile.go:137`). Commands resolve argv[0] through the profile path (`host/broker/exec_profile.go:489`). The handler verifies its recorded hash before each call (`host/broker/handlers_exec.go:270`). There is no declared secondary browser executable identity for a Go test that launches Chrome.

Node runners already have a working declaration: [TwilightGame's profile](../verification/world-row140-m4/twilightgame.json) runs tsc and vitest through node. This row adds availability and confinement gates for Chrome and cargo. It does not invent node testing.

The runtime pin is srt 0.0.78 (`host/broker/exec_sandbox.go:27`). The settings renderer supplies empty allowed and denied domain lists and leaves local binding false (`host/broker/exec_settings.go:118`). Row 145's V6–V9 recorded a successful bind with local binding enabled, but also a successful unauthenticated daemon connection. Current `isProtected` still protects only POST /v1/commit (`host/daemon/daemon.go:1088`). Those historical network measurements are not a fresh experiment here.

Chrome is a second mismatch. The launcher uses headless mode and a fresh user-data directory (`host/daemon/workbench_live_chrome_test.go:68`). Its comment records that relocating HOME on macOS stopped URL loading (`host/daemon/workbench_live_chrome_test.go:77`). Exec always relocates HOME and TMPDIR (`host/broker/handlers_exec.go:454`). The comment is evidence of the earlier measurement, not proof that a proposed confined setup works.

The synchronous cap is 9,000 ms (`host/broker/exec_profile.go:33`), enforced by the parser (`host/broker/exec_profile.go:249`). HandlerCap is 10 s (`host/coordinator/effectful.go:47`). Browser startup, cold cargo builds and full node suites need row 145's asynchronous lane. Protected-path checks can reduce a short command's usable time further.

Finally, tests are executable project code. The real service controller calls /bin/launchctl through toolprobe (`cmd/ailang-worldd/service_launchctl.go:33`). Friction 47 records an operator's outer sandbox wrapper breaking process-group tests. Operator-supplied friction 48 records long-suite timeouts and denied httptest binds. Operator-supplied friction 49 records a sandbox pass followed by a host failure involving /var and /private/var. Entries 48–49 are not in this base's friction log. They are supplied reports, not measurements made here.

## §1 Premise verification log

“Checked” below means source inspected through World, unless a command is named. No browser, cargo or Go test ran.

| premise | how checked | result |
|---|---|---|
| P1 Profile schema and digest are strict. | `host/broker/exec_profile.go:137`; digest assignment at line 156 | checked; digest hashes exact profile bytes, including future fields |
| P2 Program identity is checked before spawn. | `host/broker/exec_profile.go:489`; `host/broker/handlers_exec.go:270` | checked; current resolution is at profile parse, hash verification is per call |
| P3 The short cap is fixed. | `host/broker/exec_profile.go:33`; `host/broker/exec_profile.go:249`; `host/coordinator/effectful.go:47` | checked; 9 s profile maximum, 10 s handler |
| P4 srt is pinned and network is closed. | `host/broker/exec_sandbox.go:27`; `host/broker/exec_settings.go:92`; `host/broker/exec_settings.go:118` | checked; 0.0.78, no domains, no local-binding assignment |
| P5 HOME and TMPDIR belong to World. | `host/broker/exec_profile.go:60`; `host/broker/handlers_exec.go:454` | checked; profile env cannot override World's names |
| P6 Opening local binding also exposes loopback listeners. | [row 145](w-workspace-exec-long-run-lane.md) §3 V6–V9 and §4.9; `host/daemon/daemon.go:1088` | historical measurement and current auth code checked; live network reproduction **verify** |
| P7 Node test commands already exist. | [TwilightGame profile](../verification/world-row140-m4/twilightgame.json) commands | checked; fixed node prefixes for tsc and vitest |
| P8 Chrome skips and retains HOME deliberately. | `host/daemon/surface_live_chrome_test.go:22`; `host/daemon/workbench_live_chrome_test.go:59` and line 77 | checked source; fresh HOME and nested sandbox experiments **verify** |
| P9 Protected patterns must survive added tools. | `host/broker/exec_protected.go:50`, `host/broker/exec_protected.go:236`, `host/broker/exec_protected.go:317`; [row 156](../implemented/w-exec-protected-test-paths.md) §4 | checked parser, matcher and OS-specific rendering; Linux basename patterns remain detect-only |
| P10 Mint binds the selected profile. | `cmd/ailang-worldd/session_new.go:333` and line 379; `cmd/ailang-worldd/main.go:508` | checked load, existing-session digest comparison and serve flag |
| P11 Subprocess additions require the census. | `host/broker/registry_publish_test.go:1137` | checked; driven-site and credential-scrubbing test exists |
| P12 Home service uses a real host controller. | `cmd/ailang-worldd/service_launchctl.go:33`; friction 47 | checked source and logged report; descendant enforcement **verify** |
| P13 Adjacent designs own resolution, long runs and git reads. | Required designs found with World `ls design_docs/planned design_docs/implemented`; row 140 §4.3–§4.5, row 145 §3 and §8, row 219 §3.6/M3, row 228 §3.4/M3 read | checked; proposals distinguished from base code |
| P14 Temp-dir differences can hide a host failure. | Operator's supplied friction 49; current TMPDIR assignment above | supplied report; paired reproduction on this base **verify** |
| P15 Chrome and cargo can run under the proposed fence. | This mint permits only `grep` and `ls`; neither browser nor cargo run is available | **verify**; no performance or compatibility claim |
| P16 Search flags and output handling cost extra work. | World `workspace-exec grep` with `-n` refused; broad read batches exceeded host output | measured this episode; friction 50–51 |

## §2 Conflict surface

| owner | boundary for row 221 |
|---|---|
| [145: long-run lane](w-workspace-exec-long-run-lane.md) | D-WORLD-84 accepted all A. Reuse start/collect, ExecLong grant, long_timeout_ms, snapshots, admission, cancellation and replay. Do not enlarge HandlerCap or add a second runner. |
| 199: daemon read auth, then loopback | Security prerequisite. Row 199 owns authentication and per-command allow_local_binding. 221 consumes that field only after its negative daemon-read test passes. No early global switch. |
| [219 M3](w-worldd-home-service.md) §3.6 | Owns per-call resolution, changed-target probes and clearer launch errors. Base code still resolves commands at parse. 221 extends the same resolver to declared auxiliary tools; it must not fork resolution logic. |
| [228 M3](w-agent-search-and-read-cap.md) §3.4 | workspace-git is a separate read-only broker effect with filter and revision guards. It does not require exposing git administration through exec. Preserve row 156's common git and protected-path code. |
| 225: run wrapper | Owns confinement of adapter-less agent hosts and their daemon/model network. A dev-tool child must not inherit those permissions or credentials. Nested sandbox-exec is not a workable macOS composition. 225 must use brokered tools across its boundary or prove an equivalent single enforced policy before this row claims wrapper support. |

Row 221 can land profile parsing and offline fixtures before 199 or 145. It cannot claim the attended browser exit until both dependencies work. Host-service enforcement is required for every episode exec command, including existing Go commands that can spawn arbitrary descendants.

## §3 The design

### Profile declarations and identities

Keep named commands, fixed argv prefixes and existing argument grammars. No free-form executable, shell text or agent env is added. Propose optional command keys `dev_tool` (enum `chrome-test | cargo | node-test`) and `tools` (a list of auxiliary tool IDs). Propose a top-level `dev_tools` map. Each entry declares `program`, `sha256` and optional `resource_root` with `resource_sha256`. Paths are operator-owned absolute paths outside worktrees and World state. Validate their literal and resolved read roots against the existing fence. Reject unknown keys, missing references and invalid hashes.

A Chrome-test command is still a fixed Go test command. It references a Chrome entry in tools. World supplies WORLD_CHROME only to that command, after verifying the tool. A cargo or node-test command declares its primary executable in argv and references the pinned toolchain entry too. This prevents a pinned cargo launcher from silently selecting an unpinned rustc. It is a declaration of the intended toolchain, not a claim that arbitrary project code cannot invoke another readable executable.

For node, use a fixed prefix such as node plus node_modules/vitest/vitest.mjs plus run. Worktree runner code remains project input. Record its lockfile digest and snapshot identity; do not describe that mutable script as an immutable host tool.

The exact profile bytes, including dev_tools, class, auxiliary IDs, allow_local_binding and row 145's long_timeout_ms, are covered by the mint's profile digest. Changing any requires a newly bound mint. These keys are proposals: today's strict parser rejects them. Before adoption, coordinate their schema spelling with 145 and 199. Old profiles retain their meaning.

Use row 219 M3's per-call resolver for primary and auxiliary executables. Record configured path, resolved path, actual binary hashes, resource identity, profile digest and settings digest. This row's explicit tool hash is an additional constraint: re-pointing a Chrome symlink to different bytes refuses even if 219's probe succeeds. An operator updates the profile and remints to accept an upgrade. Hash and probe failures spawn no project code. Pin browser helper/resources as a read-only installation tree; executable-only hashing is inadequate for a Chrome app bundle. Toolchain seeds must have an operator manifest and digest too.

### Chrome: isolated data, measured HOME, two sandboxes

The operator installs a dedicated Chrome for Testing distribution and pins its real platform executable and resource tree. Do not discover a personal Chrome by PATH or use the operator's default browser installation as an implicit dependency. Do not install or update Chrome during an episode.

Each invocation gets private cache, tmp and user-data directories beneath its episode cache, or beneath row 145's run cache. Mode 0700; no reuse across episodes. Force --headless=new, --no-first-run, --no-default-browser-check, --disable-background-networking and --user-data-dir=<fresh directory>. Do not expose cookies, extensions, login data, keychain material or the operator's browser profile. Agent flags cannot replace the profile directory, load extensions or disable sandboxing. Debugging endpoints, if used, are local and ephemeral.

First measure the existing launcher under the exact rendered World settings, with relocated HOME. Require an actual request and a page assertion, not merely a live Chrome process. If it fails as the source comment predicts, measure a broker-owned Chrome child environment that retains the original HOME string while the operator home remains denied by srt. Only the dedicated browser installation and private cache directories are allowed. The HOME string is not permission to read that home. Other children keep World's relocated HOME. This exception requires R3 and a driven launch seam; an editable project script must not receive authority to widen the read fence.

Prefer relocated HOME if it works. If retaining the string needs access to real user data or a broader home read, refuse Chrome on that platform. The compatibility gate stays open. A dedicated OS account is a possible later isolation choice, not a claim of a working r0 solution.

Keep Chrome's own renderer sandbox enabled inside srt. No --no-sandbox, weaker network isolation, Apple Events or universal Unix sockets. On macOS measure the nested Seatbelt behavior; on Linux measure the namespace/seccomp behavior. The known nested sandbox-exec failure makes this a real risk. If the combination fails, return unsupported with its reason. Neither an unsandboxed browser nor a disabled inner sandbox is an automatic fallback. All browser descendants must be reaped on timeout and cancellation. Any new production launch site joins the existing subprocess census and credential-scrubbing controls.

### Cargo and node: dependencies ready before the episode

Registry fetches stay off. Use cargo test --offline --locked and a pinned Rust toolchain with no rustup auto-install. Provide vendored sources or an operator-seeded per-episode CARGO_HOME. Keep compiler and sysroot resources read-only. Targets, incremental outputs and writable cache copies live inside the worktree/run cache. Cargo build scripts and tests are untrusted project code inside the same fence.

Use pre-provisioned node_modules or an operator-seeded offline dependency cache. Prefer direct node runner prefixes over npm scripts that can hide an install. No npx download, npm registry fetch, git dependency fetch or telemetry permission. If an offline materialization step is declared, its scripts still run confined. Record the lockfile and seed digests. Never make a shared seed writable.

An absent dependency fails explicitly. The operator seeds it outside the episode, reviews the manifest and remints if profile bytes change. This row adds no allowlisted registry domains. A later network design must preserve the distinction between fetching dependencies and allowing arbitrary test code to send network requests.

### Loopback: auth first, opt-in second

Row 199 must first require credentials for /v1/* reads, with episode scope where applicable. Prove that a sandboxed child without a token cannot read log entries, objects or reports. Also prove a wrong-episode token cannot read another episode's scoped routes. Preserve operator CLI operation through its authenticated path.

Only then accept row 199's per-command allow_local_binding, default false and named in the profile digest. Render it per invocation, including row 145's snapshot settings. Strip daemon session tokens and registry credentials from every dev-tool environment. Browser tests use their own test daemon/server, not the operator's live daemon credential.

Be explicit about the residual: srt's measured switch permits connections to other loopback listeners, not just the port the test bound. Daemon authentication protects World records; it does not protect an unrelated unauthenticated localhost service. Declare this residual at mint and permit it only on an operator-approved test host. If that residual is unacceptable, leave the command disabled until stronger network isolation exists. Measure IPv4 and IPv6 reachability and non-loopback binding too; do not call allowLocalBinding a port allowlist or assert it guarantees loopback-only listening.

### Timeouts and long runs

Keep timeout_ms at 9,000 or below. Short version probes, warm targeted cargo tests and individual node files may use workspace-exec only after measured fit, including verification and protected-check headroom. Never infer fit from an unsandboxed timing.

Chrome drills, cold cargo builds and full node suites use row 145's workspace-exec-start and workspace-exec-collect. Declare long_timeout_ms per command beneath its operator ceiling. Use its snapshot, protected proof, terminal result, separate pool and cancellation. Resolve auxiliary paths and rebase cache/data paths for that snapshot. Do not leak live worktree output back as if it were snapshot output. Collect retains the existing output limits and recorded replay. Until the lane lands, report these runs unavailable rather than repeatedly retrying 9 s calls.

### Host-service deny for episodes

Make denial mandatory for the entire episode exec process tree. It is not an optional flag the agent can switch off, and it is not merely omission of launchctl from commands. Project tests can call absolute paths.

Deny execution of /bin/launchctl on macOS and systemctl on Linux, including systemctl --user. Deny access to the host user-service control channels as well: launchd management IPC on macOS and the user systemd/D-Bus sockets on Linux. Remove DBUS_SESSION_BUS_ADDRESS and XDG_RUNTIME_DIR inheritance. A copied controller or direct IPC must not bypass the policy. Ordinary test pipes and private test sockets need separate measurements.

Enforce this in World's single srt confinement boundary. Do not add the friction-47 outer sandbox-exec wrapper. Keep production host service operations outside episode authority, behind their existing attended route. Test service code through fake controllers. The enforcement mechanism for macOS management IPC is **verify** at M0; do not pretend denyRead alone is a general execute/IPC guard. If srt cannot express a sufficient deny, block activation on that platform and resolve the runtime policy extension before rollout.

Run process-group lifecycle tests under the exact final policy, with controls on clean dev. No new host subprocess may bypass the AC10 census. Refuse unsupported settings rather than silently dropping them.

### Temp-dir parity

Keep private TMPDIR; do not make the host's shared temp writable. Record the effective temp root and its canonical realpath in check evidence. On macOS measure both /var and /private/var spellings using the same fixture. Host reruns must include the World temp layout and a normal host t.TempDir layout, with explicit canonical-path assertions.

A single matching private temp layout cannot establish host parity. Friction 49's reported path-alias failure must be reproduced and killed by the paired check. Fix path assumptions in the relevant future implementation task, not by expanding this design-only episode's file scope. An environment mismatch must be visible in the report, even when the sandbox test passes.

## §4 Failure modes (three-part, row 227)

Commands below describe the proposed interface and declared command IDs. They are not available in this writing episode. Operator-only commands are labelled.

| what failed | why | the next command |
|---|---|---|
| chrome-test refused before spawn | missing tool declaration or pinned bytes changed | Operator: `ailang-worldd profile init --help`, then review the dedicated tool entry and mint a replacement session. Never auto-accept the new hash. |
| Chrome started but the drill failed | no request reached the test URL; HOME or inner sandbox compatibility probe failed | `workspace-exec-start {command:"chrome-probe"}` after the operator declares the diagnostic; collect its startup log and HOME-mode result. |
| browser command cannot bind | allow_local_binding is false, or row 199 has not passed its auth gate | Operator: `ailang-worldd doctor --session <session-file>`; complete row 199 before changing the profile. |
| test reached daemon but read was denied | no child token; this is the required security result | `workspace-exec-start {command:"daemon-read-negative"}` to retain the expected denial as test evidence; do not inject a live token. |
| short build timed out | command exceeded its effective short budget | `workspace-exec-start {command:"cargo-test"}` if declared long-capable; otherwise operator remints with that declaration. |
| cargo/node cannot find a dependency | offline source or seeded cache is incomplete | Operator: `ailang-worldd profile init --help`; provision the declared seed using the project's locked dependency procedure, then repeat the same command. No online retry in the child. |
| service-control attempt denied | episode descendants have no host-service authority | `workspace-exec-start {command:"service-tests"}` with fake controllers; genuine host administration stays attended. |
| host rerun disagreed with sandbox | temp spelling or environment differs | `workspace-exec-start {command:"temp-parity"}`; compare its recorded canonical paths with the operator's paired host check. |
| long run is still running | terminal collect has not arrived | `workspace-exec-collect {run_id:"<id>",wait_ms:8000}`; cancellation uses the same tool with cancel:true. |

Error records must name the command, failed pin or setting, platform and whether anything spawned. Compatibility failures include bounded Chrome stderr. Credentials and personal data must not appear in diagnostics.

## §5 What this design does not do

It does not implement code, tests or profile files in this episode. It does not change the charter or resolve rulings.

It adds no general shell, arbitrary env, Docker, dependency network access or personal-browser attachment. It does not reopen D-140-4 or implement row 145 again. It does not implement row 199's auth or claim its switch isolates localhost services. It does not build row 225's wrapper, add git write verbs, deploy services or authorize outward effects. Executable declarations are provenance and admission controls; arbitrary project code still requires confinement. Passing tests with skipped Chrome do not satisfy the exit criterion.

## §6 Rulings (proposed; recorded as D-WORLD-NN when made)

- **R1 — timeouts. (A, default)** Keep the 9 s synchronous cap and use row 145's declared long_timeout_ms lane for browser drills and longer builds. **(B)** Add a separate per-command synchronous cap above 9 s. A preserves D-140-4 and the frozen deadline chain; B conflicts with that ruling and is not adopted by this row.
- **R2 — loopback scope. (A, default)** Depend on row 199: read auth first, then per-command allow_local_binding with the all-listeners residual declared for an approved test host. **(B)** Wait for network namespaces or a port-specific isolation mechanism before any loopback opt-in. A meets the measured need sooner but does not protect unrelated local services; B reduces exposure and delays the exit run.
- **R3 — Chrome. (A, default)** Pin dedicated Chrome for Testing bytes/resources, require fresh private user data and both sandboxes, and gate the narrowly described HOME-string exception on measurement. **(B)** Support Chrome only with relocated HOME, using a separate test account or host if needed. A may preserve macOS compatibility without opening user data; B simplifies environment rules but may leave this host unsupported.
- **R4 — dependency network. (A, default)** No registry fetches; cargo and node use vendored or operator-seeded locked dependencies only. **(B)** Declare allowlisted registry domains per command in the digest. A preserves the current no-network boundary; B improves cold setup but grants project code egress and needs a separate network design.
- **R5 — host services. (A, default)** Mandatory descendant-wide launchctl/systemctl and management-IPC deny in World's confinement, failing closed when unsupported. **(B)** Depend on fake-controller test conventions and an attended grader wrapper. A prevents project code reaching the real service; B is weaker and repeats friction 47's composition failure.
- **R6 — declared identity. (A, default)** Extend named commands with dev_tool/tools and a strict dev_tools map; explicit hashes pin tools across row 219's per-call resolution. **(B)** Rely only on argv[0] and record whichever auxiliary binary the project chooses. A adds profile work and deliberate remints on upgrades; B leaves the browser outside the mint's declared tool identity.
- **R7 — temp parity. (A, default)** Keep private temp roots and require paired canonical/alias fixtures on the host and inside World. **(B)** Accept only the sandbox result with its environment recorded. A costs a second check and exposes friction 49; B can miss the same host failure again.

## §7 Open questions

1. Can dedicated macOS Chrome load the test URL with relocated HOME under srt 0.0.78? If not, does the retained-string variant work while all real home data stays denied?
2. Can Chrome's own sandbox run inside the exact World boundary on both platforms? What denial explains failure? No answer is assumed here.
3. Which srt policy mechanism blocks launchd management IPC without breaking test process groups? Can Linux block direct user-manager IPC as well as systemctl?
4. What are the measured browser and cold-build caps, including snapshot and protected-proof costs? Row 145 owns cap admission.
5. What operator provisioning procedure produces reproducible read-only Chrome and Rust resource manifests, including helpers and symlinks?
6. Will row 199 require authentication for non-/v1 surfaces that could reveal records? Its threat review must inventory those routes before loopback rollout.
7. How will row 225 compose with the broker without nested sandbox-exec or inherited daemon credentials? This row makes no compatibility promise until that check exists.

## §8 Milestones

All test names below are proposed gates unless explicitly called existing. They are not tests written in this episode.

| M | what | gate |
|---|---|---|
| M0 | Measure dedicated Chrome, HOME variants, nested sandbox, host-service denies and temp aliases on macOS/Linux | Operator check under exact srt 0.0.78 World settings; bank requests, stderr and negative controls; unsupported combinations block activation |
| M1 | Strict declarations, digest binding, auxiliary pins and row 219 resolver integration | Go tests TestExecDevToolsProfileStrict, TestExecDevToolsMintDigest and TestExecDevToolsChangedTargetRefuses |
| M2 | Offline cargo/node fixtures, private caches and resource manifests | Go test TestExecDevToolsOfflineFixtures; operator check with empty versus seeded caches and no network |
| M3 | Consume row 199 loopback after auth; strip child credentials | Go tests TestExecDevToolsCannotReadDaemon and TestExecDevToolsWrongEpisodeReadDenied; operator IPv4/IPv6 and external-network check |
| M4 | Browser launch/data isolation, service deny and row 145 lifecycle integration | Go tests TestExecDevToolsChromeIsolation and TestExecDevToolsServiceDeny; existing subprocess census and process-group tests; operator real-Chrome request check |
| M5 | Snapshot proofs, replay, timeouts, temp parity and complete reporting | Go tests TestExecDevToolsLongRunReplay, TestExecDevToolsCancelTree and TestExecDevToolsTempParity; operator paired host rerun |
| M6 | Attended exit run: fresh episode like w188l, agent runs its own Chrome test and operator rerun only confirms it | Operator checks non-skipped TestSurfaceLiveChrome evidence, matching pinned tools/input snapshot, terminal World result and confirming host rerun; include cargo/node checks on their fixtures |

## §9 Acceptance criteria (no model call; an operator or CI can check each)

1. A Go parser table rejects unknown dev-tool keys, missing auxiliary IDs, malformed hashes and invalid resource roots. Existing v1 profiles still parse identically.
2. Changing tool hash, class, loopback flag or long_timeout_ms changes the exact-byte profile digest. A session bound to the old digest cannot execute the new declaration.
3. Re-pointing a tool symlink invokes row 219's resolution checks. Different explicit tool bytes refuse before project spawn. The record of a successful run names actual primary and auxiliary hashes.
4. A Chrome fixture makes a real request and asserts page behavior under World. Missing WORLD_CHROME or a skipped drill fails this acceptance check, even if go test exits zero.
5. Chrome gets a fresh private user-data directory per run. It cannot read a planted operator-profile decoy or sibling episode cache. Both the default HOME mode and any exception have explicit passing/denied-read evidence.
6. The final Chrome invocation has no --no-sandbox or equivalent sandbox-disabling argument. An incompatible inner sandbox yields a named unsupported refusal, never a host fallback.
7. Before loopback opt-in is enabled, a sandboxed no-token request to daemon /v1/log, object and report routes is denied. A valid wrong-episode token cannot read episode-scoped data. Authenticated operator reads still work.
8. With opt-in false, httptest binding fails under the measured fence. With it true, the declared fixture binds and serves. Evidence records reachability to another dummy loopback listener and IPv6 behavior; external registry access remains denied.
9. Cargo offline/locked and direct node runner fixtures pass from declared seeds. Empty-cache fixtures fail without fetching. Seed contents remain unchanged and outputs stay in permitted roots.
10. timeout_ms=9001 still refuses. A deliberately longer fixture completes through row 145 start/collect and carries the snapshot digest, tool identities, protected proof and terminal output.
11. Cancelling or timing out a browser/build long run reaps its descendants. Another episode cannot collect its run ID. Replay returns recorded outputs without starting a browser or compiler.
12. An episode descendant cannot execute the real launchctl or systemctl --user, or reach equivalent host-management IPC through a copied client/direct connection. Fake-controller tests and existing process-group tests pass under the same policy.
13. New subprocess sites pass the existing driven-site/registry-credential census. A sentinel daemon credential and registry credential are absent from browser, cargo, node and their children.
14. Row 156 patterns and integrity reporting remain unchanged, including Linux detect-only basename behavior. A fixture trying to alter a protected test cannot produce a clean intact result.
15. The paired temp fixture exercises /var and /private/var aliases on macOS and canonical temp roots on Linux. World and host reports expose differences; the supplied friction-49 failure must be reproduced before its fix is accepted.
16. In an attended fresh episode, the agent's terminal World record proves TestSurfaceLiveChrome ran without skip on the reviewed input. The operator reruns that test with the same tool pins, then with host temp layout; both confirm it. An operator-only first execution does not meet this criterion.

## Sources

- [Roadmap](w-all-agent-work-through-world.md) §2 row 5, §3 step 4; [charter](../world-mission.md) rows 199, 221 and 225.
- [Row 140](w-workspace-exec-toolchain-effect.md) §4.3–§4.5; [row 145](w-workspace-exec-long-run-lane.md) §3 V6–V9, §4 and D-WORLD-84.
- [Row 156](../implemented/w-exec-protected-test-paths.md) §4; [row 219](w-worldd-home-service.md) §3.6 and M3; [row 228](w-agent-search-and-read-cap.md) §3.4 and M3.
- [Shape and voice model](w-surface-index-follow.md); [TwilightGame declared commands](../verification/world-row140-m4/twilightgame.json).
- [Friction log](../verification/world-dogfood/friction.md) 42, 44 and 47, plus this episode's 50–51. Friction 48–49 are operator-supplied unmerged reports quoted in the task.
- Go sources checked in §1. Compatibility, performance and enforcement experiments remain **verify** until M0; source inspection is not a passing runtime test.
