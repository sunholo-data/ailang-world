# w-worldd-developer-cli — a developer CLI for World, so operators and agents stop hand-rolling curl

**Status**: PLANNED — attended draft 2026-10-03 (read-only Plan agent; written to file by the attended controller). Queue row 138. **RULED D-WORLD-57** (keep two pins + `setup`; all four CLI groups). **D-CLI-1..5 = recommended defaults (Mark, attended 2026-10-03).** Base `5a5347b63ee6`.
**Clauses**: 2 (operator usability, S7 docs), 5 (the provenance walk becomes one command).
**Direction (ruled)**: `ailang-worldd` subcommands `tools list`, `call`, `why`, `log tail`, `setup`, `doctor`, `session new|list|revoke`. Helpers never bypass an attended fence.
**Absorbs**: every M6 snag encoded in `tools/attended/se_smoke.sh` (V27); row 137 drift items (5) and (6) (§9).
**Estimate**: ≈ 5.25 executor days (M1–M6) + 0.25 d attended (M7), versus the row's ≈ 2 d. Three measured things drive the difference:
- the walk back from a result needs a log scan (V12, V13);
- `setup` needs a new network package to keep the boundary gate honest (V23);
- `session new` adds a subprocess site that the AC10 census must drive (V25).

## 1. Why (measured)
- **P1 — there is no MCP client.** The CLI has no `tools`/`call`/`why`/`setup`/`doctor` verb and no `log tail` (V4). Calling a tool means curl with a bearer header, two Accept types, SSE `data:` parsing and base64 decoding (V9, V14). `se_smoke.sh` does all of that in Python (V27).
- **P2 — `/mcp/` has three response shapes**, which a naive client conflates:
  - success is SSE (V9);
  - a session problem is a `text/plain` 401 (V10);
  - every coordinator failure, including no-world and EffectsUnrecorded, is a 200 `application/json` `-32603 "host callback failed"` envelope (V11, V16; row 136).
- **P3 — provenance can't be walked back from an MCP result.** The result carries `world.plan` and `world.effects[].record`, but no invocation id, entry index or record ref (V12). `/v1/receipts` refuses anything but `rest:` ids (V13).
- **P4 — install is a runbook.** Pins are installed by hand. CI checks only the release's own `.sha256`, which the release itself attests (V21). The rig's interpreter directory holds a stale v0.30.0 tarball (V20).
- **P5 — the session helpers are buggy** (row 137):
  - the revoke usage line fails as written;
  - an unknown id prints "revoked" and exits 0;
  - **revoke or mint on a missing `--db` creates a store** (V5);
  - the top-level usage omits `session` and `commit --session` (V2);
  - mint and revoke need the daemon stopped (V6);
  - there is no `session list`.
- **P6 — the same snags get rediscovered by hand:** the API-key guard (V8), the two different TTY fences (V7), a busy port or held lock (V6), and a missing genesis (V16).

## 2. Verification Log (2026-10-03, base V1, read-only; `B` = the CLI built from V1)

| V | Command | Observed (short) |
|---|---|---|
| V1 | `git log -1 --format='%H %cd' --date=short` | `5a5347b63ee6… 2026-10-03` |
| V2 | `$B` (no args); `sed -n 50,92p cmd/ailang-worldd/main.go` | lists serve, health, head, world/object/log/registry get, object find, log range, `commit --file`; **no `session`, no `commit --session`**; rc 1 (`$B help` rc 0) |
| V3 | `$B serve --help`; `$B session mint --help`; `$B session revoke --help` | serve prints the top-level usage, rc 1 (no per-verb help); mint lists `lifetime … (default 3600) (default 3600)`; revoke lists only `-db`; all rc 1 |
| V4 | `$B tools list` / `call x` / `why x` / `setup` / `doctor`; `$B log tail`; `$B session list` | `unknown command` rc 1; `log: usage: log get <index> \| log range …`; `session: unknown subcommand "list"` |
| V5 | `$B session revoke $(64×a) --db $S/store/world.db`; then `--db $S/store/world.db $(64×a)` against an empty dir | (1) `usage: session revoke …` rc 1. (2) `revoked session credential aaaa…` **rc 0**, and `world.db` + `world.db.writer.lock` **were created** (`store.go:1352`: deleting an absent credential is a no-op) |
| V6 | with a daemon running on the store: `$B session revoke --db …`; a second `serve`; `serve` on a busy port | `store: another process already holds the writer lock …` rc 2; `daemon startup failed at store-open`; `… at listen: cannot bind 127.0.0.1:7699 … address already in use` |
| V7 | `$B session mint … </dev/null`; `cmd/ailang-worldd/session.go:129-134`; `cmd/world-publish/tty.go:91-121` | mint: `refusing: no controlling terminal (open /dev/tty: …)` (it only needs `/dev/tty` to open). world-publish has three refusals: `no-controlling-terminal`, `stdin-not-a-terminal`, `stdin-is-not-the-controlling-terminal` |
| V8 | `AILANG_REGISTRY_API_KEY=x $B health`; `main.go:99-110` | refused rc 2 before flags are parsed; "there is no subcommand that is exempt" |
| V9 | the M6 `call-read.json`; `mcphttp` `wire.go:83-89`, `methods.go:97-100` (v0.47.2) | `event: message\ndata: {…"result":{"content":[{"type":"text","text":"<output>"}],"structuredContent":{…}}}`; `text` = the output bytes |
| V10 | `curl -si … /mcp/` with no / unknown / malformed bearer | `401` `text/plain`: `session credential is absent: …` / `unknown session credential` / `malformed Authorization header …` |
| V11 | `serveapi/protocol/envelope.go:21-51`; `host/projection/mcp.go:126-135`; mcphttp `handler.go:95` | a host error → 200 `application/json` `{"error":{"code":-32603,"message":"host callback failed"}}` (**not SSE**); a wrong Accept → 400 |
| V12 | `mcp.go:135`; `coordinator.go:388-391`; `projection.go:381` | MCP returns only `OutputBytes`; only A2A returns `metadata{invocation_id, world_ref, entry_index}` |
| V13 | `curl -si :7699/v1/receipts/a2a:ep1:abc`; `…/rest:nope` | 400 `receipt id must be rest:<id>`; `{"invocationId":"rest:nope","state":"not-started"}` |
| V14 | `host/coordinator/plan.go:36-47,130-157`; `host/broker/record.go:32-43`; `host/daemon/handlers.go:43-68` | record v2 = v1 + `plan` + `effects[]`; `entryIndex = revision+1`; `stateRoot` = the output hash; `worldRef = sha256(json{revision,stateRoot,logHead})`; object payloads are base64; by-semantic-id pages are capped at 500 |
| V15 | `daemon.go:793-801`; on an empty store `curl :7699/v1/head`, `…/v1/log?from=0` | only `POST /v1/commit` is protected; head → 404 `no world head has been selected yet`; log → `{"items":[]}`; a present head is `text/plain` |
| V16 | `coordinator.go:349-355` | no head → `WorldAbsentError` → reaches MCP as V11's generic envelope |
| V17 | `grep -n 'writeTimeout\|invokeDeadline\|defaultClientTimeout\|readDeadline' host/daemon/daemon.go`; `cli.go:21` | write 30 s, invoke 20 s, read 10 s; CLI client 30 s, response cap 8 MiB; MCP request cap 4 MiB |
| V18 | `gh release view v0.52.1 -R sunholo-data/ailang --json isDraft,publishedAt,assets`; the same for v0.41.0 | v0.52.1 **`isDraft:false`, published 2026-10-03T09:50:52Z**. Per platform: `<plat>.ailang.tar.gz{,.sha256,.sig,.pem}`, plus `SHA256SUMS{,.sig,.pem}`. Downloads 302 to `release-assets.githubusercontent.com` |
| V19 | download + `shasum` + `tar tzvf` for each pin | v0.41.0 tarball `b08f3cde…598e0b` → binary **`1a67b014…5b9f`**; v0.52.1 `576236fe…a7579e` → **`0dd70a1d…a8f5`**; v0.51.0 `784a68c7…63dbdb` → **`e55ff71c…9c7f`** (= the `tool` field in V9) |
| V20 | `ls -la ~/.pinned-ailang ~/.pinned-ailang-tools/*`; `shasum` | the interpreter dir has `ailang` (`1a67b014…`) plus a **stale v0.30.0 tarball** (`ac3174e0…`); the tool dirs are versioned |
| V21 | `.github/workflows/ci.yml:78-108,268-287`; `scripts/verify_world_package.sh:25-37` | CI checks only the release's own `.sha256`; the package gate compiles in the binary digest per platform (Darwin/arm64 `1a67b014…`, Linux/x86_64 `8e7a275d…`) |
| V22 | `grep -rn '"AILANG v0\.' host cmd` (non-test) | `host/daemon/workspace.go:38 ToolBinaryRelease = "AILANG v0.51.0"`; `cmd/world-publish/fences.go:91 frozenCompilerVersion = "AILANG v0.41.0"` |
| V23 | `host/boundary/allowlist_world_test.go:43-58,880-900`; `go list -deps ./host/replay/... \| grep -c host/archive` | cmd's net/http is the "documented loopback-IPC exception"; store, replay and world-publish forbid net/http in their closures; **replay imports archive** |
| V24 | `host/store/context_roots_test.go:17-56,171-214` | named context-root pins; `TestProductionGoSurface` walks host/ and cmd/ |
| V25 | `host/broker/registry_publish_test.go:1075-1120,1180-1200` | the AC10 census enumerates every `exec.Command(Context)` in host/ **and cmd/**; each needs a driver |
| V26 | `cmd/ailang-worldd/session_budget_test.go` | pins the literal revoke context and call |
| V27 | `tools/attended/se_smoke.sh` header + mint/MCP sections | the snag list; 6 grants × `worktree:50`, `--ttl 14400`; MCP parsed as the first `data:` line |
| V28 | `host/daemon/workspace.go:44-52` | episode grammar `^[a-z0-9][a-z0-9-]{0,63}$`; 6 workspace effects; 8 tools |
| V29 | `grep -n 'func (s \*Store) .*Session\|^func OpenReadOnly' host/store/store.go` | only Mint/Resolve/RevokeSession, no list; `OpenReadOnly` takes no writer lock |
| V30 | `website/docs/reference/cli.md:58-66,132-146`; `curl` the site | the docs record the drift; the site answers 200 at `https://www.sunholo.com/ailang-world/` |
| V31 | `find ~/.ailang/examples -name '*.ail' \| wc -l`; `ls ~/.ailang/world` | 339; `world-0.1.1.db`, `world.db` |
| V32 | `curl -s :7699/v1/health` | `{"status":"ok","daemon_version":"0.1.0","db_path":…,"interpreter_ref":"sha256:1a67b014…","interpreter_version":"AILANG v0.41.0…"}`; no tool-binary or workspace fields |

## 3. Design
### 3.1 Shape and placement
- **New files in `cmd/ailang-worldd/`:** `mcpclient.go` (tools, call), `why.go`, `logtail.go`, `setup.go`, `doctor.go`, `session_new.go`, `pins.go`.
- **Two new host packages, each with one job:**
  - **`host/pinfetch`** is the only non-loopback network code;
  - **`host/worktree`** is one bounded, env-scrubbed `git worktree add|remove` site.
- **Nothing in `host/archive`,** because replay imports it (V23).
- **Exit codes** are 0, 1 and 2, plus **3 for an integrity refusal** (a hash mismatch or a broken provenance link).
- **One session resolver,** `resolveSession(flag)`, also adopted by `commit --session`:
  - precedence is `--session` > `WORLD_SESSION` > error;
  - a 64-hex value is the token; anything else is a file (≤ 256 B, warn if group/other can read it);
  - a token given raw on argv gets a warning;
  - **the token never appears in output** (MUT-TOKEN-ECHO).

### 3.2 `tools list` and `call` (a built-in MCP client)
`mcpPost` runs over the existing client (30 s, 8 MiB cap) with both Accept types. It classifies **by status and Content-Type, never by sniffing the body**:

| Wire | Classification | Output |
|---|---|---|
| 200 `text/event-stream` | SSE: multi-line `data:` joined per spec; pick the event whose JSON-RPC `id` matches; a JSON-RPC `error` → tool error | result |
| 200 `application/json` with `-32603` | host failure (V11) | runs the **disambiguation probe** (below) |
| 401 `text/plain` | session denial | the body plus a fix hint |
| 400 | transport | verbatim |

**Disambiguation probe.** It reads `/v1/head` before and after the call:
- no head → "commit genesis";
- the log grew → "a commit landed (entry N) — do not retry; `why N`";
- otherwise it passes the error through. Typed errors from row 136 pass through unchanged.

**`tools list [--json]`** shows name, description and required arguments.

**`call <tool> [--arg k=v]… [--arg-json k=<json>]… | --json '<obj>'|@file|- [--json-out] [--strict]`** prints the result fields (long strings elided with byte counts), then the `world` block (plan, and per-effect `id status record`).
- `--json-out` prints the result text byte-exact, so `sha256(stdout − \n)` is the output ref and can be piped to `why -`.
- `--strict` exits 3 when a committed result has `ok:false`.

### 3.3 `why <target>` (the provenance walk)
**Targets:**
- an entry index, or `head`;
- any `sha256:` object, dispatched on its `semanticId`;
- an `a2a:`/`rest:` invocation id;
- `-` or `--result <file>`: a `call --json-out` result, whose sha256 is the output ref (the plan is the fallback).

**Resolution** scans the log backwards from head, a page of 500 at a time, one object GET per entry. It matches on transitionRef, invocationId, plan, output, `effects[]`, or each effect record's resultRef. It is bounded by `--scan N` (default 500, max 5000) and a 60 s budget; past either it prints "not found in the last N entries" and exits 1.

**Render, each link checked ✓ or ✗:**
1. **world:** the ref is recomputed and confirmed via `/v1/worlds/<ref>`;
2. **entry:** index, writtenBy, prev, transitionFn, and the interpreter with its pin name;
3. **record:** invocation, episode, skill;
4. **input:** decoded and elided;
5. **plan:** effects or the refusal result (record.plan must equal the plan hash);
6. **each effect record:** allowed/failed/denial and budget before→after, plus the request/result sizes;
7. **output:** output must equal `world.stateRoot`.

A ✗ exits 3; `--json` prints the chain. Non-coordinator entries (a REST genesis) render the entry and objects only, and say why.

### 3.4 `log tail [--from N] [--follow] [--interval 1s] [--raw]`
- **Default start:** head − 19.
- **Each line:** index, short hash, writtenBy, and for coordinator records the skill plus effect statuses.
- **`--follow`** polls `/v1/log?from=<next>`. The cursor advances by **items received**, never by page size (MUT-FOLLOW-SKIP).
- **Ctrl-C** uses `serveSignalContext`.
- **A daemon restart** is retried with backoff up to 5 s, and the gap is named.

### 3.5 `setup` (fetch and verify both pins)
**The pin table** is `pins.go`, keyed `(release, GOOS/GOARCH)` → `{asset, tarballSHA256, binarySHA256}`, seeded from V19/V21. Linux/amd64 digests are measured in M3. The tool release comes from **`daemon.ToolBinaryRelease`** (one source of truth, so row 135's bump flows through); the interpreter release is test-bound to `frozenCompilerVersion`, `ci.yml` and `verify_world_package.sh`.

**For each pin:**
1. **Already there?** If the destination hashes to `binarySHA256`, print `✓ present` and make no network request.
2. **Download.** Fetch the release `.sha256` (≤ 1 KiB), then stream the tarball (≤ 128 MiB) into an `O_EXCL` 0600 temp file in the destination directory, hashing as it arrives.
3. **Check the tarball.** The release's `.sha256`, the computed hash and the **compiled-in** tarball digest must all agree. Otherwise exit 3, print all three, and delete the temp file.
4. **Extract.** Take exactly one regular entry named `ailang` (≤ 256 MiB): no links, no `..`, no absolute paths.
5. **Check the binary.** It must hash to the compiled-in `binarySHA256`.
6. **Install.** Only then chmod 0755 and atomically rename it to `~/.pinned-ailang/ailang` or `~/.pinned-ailang-tools/<ver>/ailang`.
7. **Clean up.** Keep no tarball. Write `pin.json` (release, platform, both digests, URL, fetchedAt).
8. **Existing file with a different hash:** refuse unless `--replace`, which keeps the old file as `ailang.prev-<sha8>`.

**Setup never executes a downloaded byte.**

**Network bounds (`host/pinfetch`):**
- https only, on the first URL and on every redirect;
- hosts allowlisted to `github.com`, `release-assets.githubusercontent.com` and `objects.githubusercontent.com`;
- URLs come only from the compiled template (no URL flag);
- `--from-dir` installs offline, with the same verification;
- timeouts 30 s (asset) / 5 min (tarball), at most 10 redirects, GET only.

**Afterwards** it creates the store directory and workspace root (0700). It refuses if the root contains the store, then prints the next attended steps with the resolved paths. An unsupported platform exits 1, naming the reason.

### 3.6 `doctor [--db] [--addr] [--workspace-root] [--online]`
Read-only. Each check prints ✓, ! or ✗ with a one-line fix, and any ✗ exits 1.
- **API key:** ✓ by construction. If the key were set, the process would already have refused (V8), so the guard stays exemption-free.
- **TTY:** the mint fence and the publish fence are checked separately, with the three publish reasons. If only the publish fence fails, it suggests `</dev/tty`.
- **Pins:** present, and sha256 compared to the table (no exec). Stale tarballs are noted.
- **Port:** whether a worldd answers (db_path, interpreter), a foreign listener, or nothing.
- **Store:** the file exists; doctor **never creates** one. The writer lock is held or free, via a new non-mutating `store.WriterLockHeld` (`LOCK_SH|LOCK_NB` on the lock file, opened read-only). A head is present, via `OpenReadOnly`.
- **Examples corpus:** its `.ail` count.
- **Workspace root:** sanity, and that each child is a git worktree with a grammar-valid name.
- **`--online`:** the site and release-asset reachability.

### 3.7 `session new|list|revoke`
- **`session new <episode> --db --workspace-root --repo [--preset se-tools] [--grant …] [--budget 50] [--ttl 3600] [--out] [--branch]`** runs in this order:
  1. validate: the grammar, that `--db` **exists**, and that the root is sane;
  2. **the unchanged mint fence**, `env.openTerminal()`, before any side effect;
  3. a y/N that names the worktree and the grants;
  4. `store.Open`, which fails fast with "stop the daemon first";
  5. `worktree.Add` (detached);
  6. `authority.Mint` through the existing core;
  7. on mint failure, remove the worktree this run created.

  The `se-tools` preset is **derived from `transitions.json` declaredEffects**.
- **`session list --db [--episode] [--json]`** uses a new `store.ListSessions` (≤ 500 rows) on `OpenReadOnly`, so it works with a live daemon. It shows the credential id hash, episode, grants, created/expires and live/expired.
- **`session revoke`** accepts flags and the id in either order. An unknown id → rc 1. A missing `--db` → refused, never created (for mint and list as well). `--episode` revokes every credential of that episode. It still needs the daemon stopped (R-CLI-1).
- **Usage:** the top-level usage lists every verb, and `<verb> --help` / `help <verb>` exit 0. The duplicated default is fixed.

## 4. Premises

| P | Premise | Evidence |
|---|---|---|
| P1 | The CLI lacks every ruled verb; its usage omits `session` and `commit --session` | V2–V4 |
| P2 | `/mcp/` has three shapes, classified by status and Content-Type | V9–V11 |
| P3 | The result text is the committed output bytes, so sha256(text) is the output ref and equals stateRoot | V9, V14 |
| P4 | MCP results carry no index or id and receipts refuse `a2a:`, so `why` must scan | V12, V13 |
| P5 | Reads are unauthenticated; payloads are base64; pages are ≤ 500 | V14, V15 |
| P6 | A no-world call surfaces as the generic -32603 | V11, V16 |
| P7 | The release assets exist, are published, and redirect to a githubusercontent https host | V18 |
| P8 | The digests are measured, and the binary digest equals the log's interpreter ref / the result's `tool` ref | V19, V9 |
| P9 | Today's install trusts only the release's own `.sha256` | V21 |
| P10 | Network code cannot live in host/archive; cmd's net/http is loopback-only | V23 |
| P11 | New exec sites must be driven by AC10, and new context roots must be pinned | V24, V25 |
| P12 | revoke: argument order fails; an unknown id is rc 0; a missing db is created; mint and revoke need the daemon stopped | V5, V6 |
| P13 | The mint fence needs only `/dev/tty`; the publish fence also needs stdin to be the tty | V7 |
| P14 | The API-key guard runs before every verb | V8 |
| P15 | Listing sessions needs a new store read; a read-only open works with a live daemon | V29 |

## 5. Milestones
- **M1 (1 d) — MCP client, `tools list`, `call`, the session resolver.**
  - AC1.1: an httptest table covering SSE single-line, multi-line and id-order; 200 JSON -32603; 401 ×4 hints; 400.
  - AC1.2: `--json-out` bytes are exact and their sha256 is the output ref.
  - AC1.3: a daemon end-to-end run: `tools list` shows 8; `call ailang-read` shows the world block, and its record resolves.
  - AC1.4: with no genesis, rc 1 and `no world head`.
  - AC1.5: a token grep finds 0 matches; `commit --session <file>` works.
  - AC1.6: the context-root pins are updated.
- **M2 (1 d) — `why`, `log tail`.**
  - AC2.1: five target forms resolve to the same entry, with every link ✓.
  - AC2.2: a crafted mismatched output → ✗ and rc 3.
  - AC2.3: `--scan 1` → not found, rc 1.
  - AC2.4: `--follow` prints each of three commits exactly once.
- **M3 (1 d) — `host/pinfetch` + `setup`.**
  - AC3.1: a good release installs, leaves no tarball, and writes `pin.json`.
  - AC3.2: a consistently tampered tarball and `.sha256` → rc 3, and the destination is unchanged.
  - AC3.3: `http://`, an http redirect and a non-allowlisted host are refused, with 0 bytes written.
  - AC3.4: an over-cap download is refused and cleaned up.
  - AC3.5: tar tricks are refused.
  - AC3.6: a second run makes 0 requests.
  - AC3.7: the table covers `daemon.ToolBinaryRelease` and the interpreter release on darwin/arm64 and linux/amd64 (linux digests measured as V-rows).
  - AC3.8: the boundary group for pinfetch is added deliberately (4 → 5), and no store/replay closure contains it.
- **M4 (0.75 d) — `doctor`, `store.WriterLockHeld`.**
  - AC4.1: lock held vs free is reported correctly, and no file is created.
  - AC4.2: all three TTY reasons.
  - AC4.3: a flipped pin byte → ✗ showing both digests.
  - AC4.4: a worldd listener is told apart from a foreign one.
- **M5 (1 d) — `session new|list|revoke`, usage/help, `host/worktree`.**
  - AC5.1: with no TTY → no worktree and no row.
  - AC5.2: answering `n` → neither.
  - AC5.3: a mint failure rolls back the worktree it created.
  - AC5.4: the preset equals the declaredEffects.
  - AC5.5: revoke works in both argument orders; an unknown id → rc 1; a missing db → rc 1 and still absent.
  - AC5.6: `list` works with a live writer.
  - AC5.7: help exits 0 for every verb, and the usage is complete.
  - AC5.8: the AC10 driver for `host/worktree` shows the env is scrubbed.
- **M6 (0.5 d) — docs (S7).** Update `website/docs/reference/cli.md` (and drop its drift notes), `install.md` (setup), `agents/provenance.md` and `agents/connecting.md` (why/call), `operating-the-daemon.md`, and QUICKSTART §9. `se_smoke.sh` uses the CLI. AC6.1: every verb in `cli.md` is bound to `help <verb>`.
- **M7 (attended, 0.25 d).** On a fresh `HOME`, run: `setup` → `doctor` → publish → `session new ep1 --preset se-tools` → `serve` → genesis → `call` → `why -` → `log tail`. Bank the evidence, with no token in it.

## 6. Load-bearing mutations

| Mutant | Killed by |
|---|---|
| MUT-SSE-FIRSTLINE | AC1.1 |
| MUT-ENVELOPE-AS-RESULT | AC1.1 |
| MUT-TOKEN-ECHO | AC1.5 |
| MUT-WHY-NOVERIFY | AC2.2 |
| MUT-FOLLOW-SKIP | AC2.4 |
| MUT-HASH-RELEASE-ONLY | AC3.2 |
| MUT-INSTALL-BEFORE-HASH | AC3.2 |
| MUT-HTTPS-OFF / MUT-REDIRECT-OPEN | AC3.3 |
| MUT-NO-CAP | AC3.4 |
| MUT-TAR-NAME | AC3.5 |
| MUT-DOCTOR-OPENS-STORE | AC4.1 |
| MUT-FENCE-BYPASS | AC5.1 |
| MUT-REVOKE-UNKNOWN-OK | AC5.5 |
| MUT-STORE-CREATE | AC5.5 |
| MUT-PRESET-DRIFT | AC5.4 |
| MUT-WORKTREE-ENV | AC5.8 |

## 7. Decisions — RULED (Mark, attended 2026-10-03): all recommended defaults
- **D-CLI-1:** downloads live in a new `host/pinfetch` package.
- **D-CLI-2:** platforms are darwin/arm64 and linux/amd64 only; others are refused with the reason.
- **D-CLI-3:** `call` exits 0 on a committed refusal; `--strict` gives 3.
- **D-CLI-4:** the store defaults to `~/.ailang/world/world.db` and the workspace root to `~/.ailang/world-ws`.
- **D-CLI-5:** release-signature (`.sig`/`.pem`) verification is out of scope → R-CLI-3.

## 8. Risks and residuals
- **R-CLI-1:** revoke and mint need the daemon stopped; a live revoke needs an admin route (a later row).
- **R-CLI-2:** `why` is O(entries) and trusts the read surface. Owner: row 136's surface tag, or an index route.
- **R-CLI-3:** trust rests on the compiled-in digests plus the release `.sha256`; cosign is unverified.
- **R-CLI-4:** doctor's shared-lock probe has a microsecond window (benign).
- **R-CLI-5:** the row-135 bump must add pin-table rows in the same commit (AC3.7).
- **R-CLI-6:** the macOS first-exec delay applies to freshly installed pins.

## 9. Conflict surface
- **Row 135:** `workspace.go:38`, the ci.yml tool step, `se_smoke.sh` (land 135 M0 first or rebase).
- **Row 136:** MCP error mapping (the classifier passes typed errors through).
- **Row 137:** shares `cli.md`, QUICKSTART and the usage text. **138 takes items 5 and 6 plus the missing-db bug; item 7 and the rest stay with 137.**
- **Row 93:** the floor scripts' pin paths.
- **Gates touched deliberately:** context roots, the boundary allowlist, the AC10 drivers, `session_budget_test.go`, `quickstart_se_test.go`, `store.go` (`ListSessions`, `WriterLockHeld`).

## 10. Non-goals
- changing the frozen `/v1` table, read-route auth, or daemon admin routes;
- automating the human acts of publish and mint;
- a genesis helper;
- cosign;
- packaging, or Windows;
- an A2A client or an MCP stdio proxy;
- executing downloaded binaries in `setup`;
- downloading the examples corpus.
