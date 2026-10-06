# Attended publish runbook — `world/core`

**This runbook stops at readiness by default.** Steps 1–3 are automated, run in CI, and perform no
public write. Steps 4–8 are **attended**: they are never run headless and never run in CI.

The asymmetry that shapes every refusal below: a public publish is **immutable**
(`registry-validator` 409s on an existing version), so a wrong publish permanently consumes
`world/core@0.1.1`, while a refusal costs a human five minutes.

**`world/core@0.1.1` — PUBLISHED 2026-10-02T07:58:43Z (Mark, attended), record
`60b46d4b…3b24` (full value in the evidence folder below), from the fresh store
`~/.ailang/world/world-0.1.1.db` (`D-WORLD-45` = A).** The served metadata matched the reviewed
golden on all four digests and the byte length, and the registry validator verified 16/16 contracts;
it is banked at `design_docs/verification/world-attended-2026-10-02-publish-0.1.1/`. The next
release needs a new candidate; the text below describes the 0.1.1 procedure as run.
`world/core@0.1.0` is already published (2026-09-21) and immutable; its served record is
preserved verbatim at `host/broker/testdata/metadata_world_core_0.1.0.json` and is never
rewritten. `D-WORLD-41` (attended, 2026-09-28) chose 0.1.1 as a proof-hardening release with
no exported signature change. The loop prepares and rehearses it; **only a human runs the
irreversible publish** (Stage B).

---

## Stage A — readiness (automated, no public write)

### 1. Build the projection

```bash
./scripts/build_world_package.sh
```

Deterministic, allowlisted projection of the four frozen modules into `packages/world-core/`.

### 2. Run the readiness gate

```bash
./scripts/verify_world_package.sh
```

Nine steps, each asserted to perform non-zero work. It ends by writing the **ready packet** and
comparing it byte-for-byte against `scripts/world_package_ready_packet.golden.json`. Green means
the artifact's identity is exactly what was reviewed.

This gate needs the **pinned v0.41.0** compiler, which is *not* the binary every other leg uses — the
replay goldens are still v0.30.0-scoped and that pin deliberately did NOT move. Point
it with `WORLD_PKG_AILANG_BIN`; there is deliberately **no silent fallback**, so a wrong binary fails
loudly rather than skipping.

### 3. Run the full repo gate

```bash
./scripts/verify_ail.sh   # its final leg invokes verify_world_package.sh
./scripts/verify_go.sh
```

`verify_go.sh` prints **two `WARNING: DATA RACE` lines on a healthy run** — that is its own
known-positive control, and it FATALs if they are absent. Do not "fix" them.

**Stage A ends here.** Everything below requires a human at a terminal.

---

## Stage B — attended publish

Every command in this stage is refused unless a human is present. The refusal is **structural, not
advisory**: `world-publish` opens `/dev/tty` and reads the typed phrase from the controlling
terminal: from stdin when stdin *is* that terminal, from the opened `/dev/tty` otherwise, never from
a redirected stdin. A CI runner and an autonomous agent have no controlling terminal, so `/dev/tty`
does not open and the command refuses; a pipe or a redirect can only feed a stdin that is not read.
(Until 2026-10-06 a stdin that was not the very file `/dev/tty` was refused outright, which stopped
IDE-pane operators too; the `< /dev/tty` that `self_mod_publish.sh` adds is still harmless.)

`--dry-run` is a rehearsal that passes every one of those fences and then makes no request at all,
so an operator can walk the exact keystrokes before the real one.

**Exit codes:** `0` done · `1` failed (or indeterminate — see step 8) · `2` usage ·
`3` **STOP**, printed as `STOP fence=<name>` on stderr, meaning *nothing happened and nothing will*.

### 4. Review the ready packet

Read `scripts/world_package_ready_packet.golden.json`. It carries `package`, `version`, `exports`,
`effects`, `tarballSHA256`, `contentHash`, `interfaceHash`, `interfaceHashV2`, `tarballBytes` and
`compilerVersion`.
`compilerSHA256` is provenance about the *machine*, not the package, and is deliberately kept out
of the byte-compared golden.

**What the readiness gate actually guarantees.** Step 7 proves the three hashes *agree* across the
local recomputation, the pkgproj projection and the manifest. Step 9 proves the canonical ready
packet equals the committed golden **byte-for-byte** (`cmp -s`, and it prints a diff on failure).
Those two facts are the verification.

**The gate prints no digests.** Measured by running it: zero full-length `sha256:…` strings in the
whole gate log, and zero truncated ones either. The dry-run's displayed hashes are truncated
*upstream* by the compiler to 17 hex characters — 68 bits — and **are not the verification**. Any
instruction to "compare the digests against the gate's output" is therefore impossible to follow;
this step used to say exactly that, and this paragraph is its repair.

So the eyeball check is **document against artifact, at full length**. These four digests are the
identity of the local reviewed projection. **They are the UNPUBLISHED `world/core@0.1.1`
candidate** (manifest version, `## 0.1.1` CHANGELOG section, publisher version fence and
confirmation phrase moved together in iteration 211 under `D-WORLD-41`). They differ from the
already-published 0.1.0 artifact by design: the tarball gained five proven contracts and the
0.1.1 release notes, while `interfaceHash` and `interfaceHashV2` are unchanged because no
exported signature changed. The loop never publishes; a human runs steps 6–7 below:


| field | digest |
|---|---|
| `contentHash` | `sha256:473517079249f3959b5aba9727f62b40130624c9915612a47736fdc1dae2a780` |
| `interfaceHash` (manifest coverage only) | `sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083` |
| `interfaceHashV2` (exported interface) | `sha256:ifacev2:b25fe03155db0c7bf595cf730295b945d6ac64ec415a1998fae8a693d621e8d8` |
| `tarballSHA256` | `sha256:07d6180986739249bc4e5adaf3eb2bf44ce59821861a86652c31149bc09df4a3` |

They are gated against the golden by `host/runbook`, so this table cannot rot silently: change the
package without reprojecting, or edit one nibble here, and the repository gate reds.

You do not have to compare them by eye. Step 5 does it mechanically, from the same file, and
refuses if any field differs.

### 5. Confirm the projection has not drifted

```bash
go run ./cmd/world-publish packet
```

Recomputes the packet from `packages/world-core/` and compares all nine fields with the committed
golden. Read-only; no terminal required; exits `3` with `STOP fence=packet reason=drift` naming the
first field that differs.

### 6. Set the session variables and build the command

**There is a helper: `./tools/attended/self_mod_publish.sh`.** It does not automate the publish and
cannot satisfy any fence for you — it exists because both failures of the first attended attempt
were ENVIRONMENT rather than decision (`fence=store reason=unopenable`, then
`fence=tty reason=stdin-is-not-the-controlling-terminal`). It sets the variables, runs the
preflight below, and invokes `approve` with stdin connected to `/dev/tty`, which is the fence's
own requirement rather than a way around it — you still type the phrase. It keeps mint, rehearse
and spend as **separate invocations** (`preflight` · `approve` · `dry-run` · `live`) for the
reason stated in step 7. The manual path below remains authoritative.


The store's PARENT DIRECTORY must exist — the store refuses to create it, and on a machine that
has never published there is nothing to create it. Measured 2026-09-21, on the first real attended
attempt: `STOP fence=store reason=unopenable … resolve parent of "…/world/world.db": no such file
or directory`. The refusal is correct and loud; the omission was this runbook's.

**Which store (decision D-WORLD-45, RESOLVED A — fresh store; used for the 0.1.1 publish).** The store that recorded the 0.1.0
publish, `~/.ailang/world/world.db`, is schema `user_version 2`; the current binary requires 4
and has no migration, so it refuses to open or modify it (measured iteration 211 on a byte
copy: `STOP fence=store reason=unopenable … has user_version 2 … binary requires 4; refusing
to modify`). Two options:

- **CHOSEN (D-WORLD-45 = A): a fresh store file** (`world-0.1.1.db`, below). The 0.1.1
  publish record then lives in a new store and is not continuous with the 0.1.0 record.
- **Alternative: migrate the legacy store first** (a future queue row; no migration exists
  today). This would keep one provenance chain.

Either way, the legacy `~/.ailang/world/world.db` must never be modified or deleted: it is the
durable 0.1.0 publish record. The helper's default still points at the legacy file, so export
the store path below before running it.

**Paste the helper's subcommands one at a time.** Measured on the 0.1.1 run: a multi-line paste
that included a placeholder `export WORLD_APPROVAL="sha256:<…>"` made `dry-run` refuse (correctly —
nothing was sent), and interactive zsh does not treat `#` as a comment, so inline notes became stray
words. The helper's own message names the saved ref; export that exact value.

```bash
mkdir -p "$HOME/.ailang/world"
export WORLD_STORE="$HOME/.ailang/world/world-0.1.1.db"
export WORLD_BIN="$(mktemp -d)/world-publish"
export WORLD_REGISTRY="https://storage.googleapis.com/ailang-registry"
export WORLD_CREDENTIAL="$HOME/.config/ailang/registry.key"
export WORLD_COMPILER="$HOME/.pinned-ailang/ailang"   # NOT /tmp — macOS wipes it on boot
```

```bash
go build -o "$WORLD_BIN" ./cmd/world-publish
```

The binary is built into a temp directory on purpose. `go run` is **not** used here, and the reason
is measured: on go1.25.6 `go run` exits `1` for a child that exited `3`, printing `exit status 3` to
stderr instead of propagating it. The STOP contract is an exit code, so the runbook uses a binary
whose code survives.

**`WORLD_COMPILER` was `/tmp/ailang-v0300/ailang` until 2026-09-21 and that path does not
survive a reboot** — macOS wipes `/tmp`, and `host/archive/archive_test.go` already records
`~/.pinned-ailang/ailang` as the durable location for exactly this reason. Measured at pre-flight
on 2026-09-21: absent. The durable pin is the one Stage A's step 9 verifies by exact bytes
(`AILANG v0.41.0`, commit `24ee108`).

`WORLD_CREDENTIAL` must name a mode-`0600` file **outside the working tree**. The API key is never
read from the environment: if `AILANG_REGISTRY_API_KEY` is set in your shell, the production handler
refuses to be constructed at all.

### 7. Mint the one-shot approval, then spend it

**This step cannot be delegated to an agent, by design.** `approve` requires a controlling
terminal (`STOP fence=tty reason=no-controlling-terminal`), and `cmd/world-publish/tty.go` records
why: every weaker candidate — an environment variable, the typed phrase alone, a sentinel file, a
naive isatty check — was rejected by asking *"can THIS loop satisfy it?"*, and each one could be.
A controlling terminal is the one candidate the loop is structurally unable to satisfy, measured
first-party. If you are reading this because an agent could not complete the publish: that is the
fence working, not a defect to route around.


Minting and spending are two separate invocations, deliberately. One command that did both would
collapse two human acts into one keystroke and hide the single-use property at the very surface
where it matters.

```bash
"$WORLD_BIN" approve --store "$WORLD_STORE" --registry-origin "$WORLD_REGISTRY" --now 1 --expires 1000 --requester "$USER" --decided-by "$USER"
```

It prints the scope you are approving, requires the typed confirmation phrase at a real terminal,
and then prints the minted `ApprovalDecisionV1` reference. Export it:

```bash
export WORLD_APPROVAL="sha256:<the ref that was printed>"
```

The stamp is an `ApprovalRequestV1` → `ApprovalDecisionV1` pair. `Session.Invoke` traverses
`payload.approvalRef` → the landed decision → its request → the canonical scope, and refuses
**before** the credential is loaded and **before** any POST. Single use is enforced by
`approval_claims`' PRIMARY KEY — durably, not by in-memory budget — so the stamp cannot be spent
twice even across a process restart. The scope binds `contentHash`, `interfaceHash` and `tarballSHA256` above (not
`interfaceHashV2`: the frozen 0.1.0 scope grammar does not carry it), so an approval minted
for these bytes authorizes no other bytes, and an approval for `0.1.0` cannot authorize `0.1.1`.
The typed confirmation phrase names the version too — for this candidate it is
`publish world/core@0.1.1 irreversibly` — so a stale 0.1.0 approval ref saved by the helper
(`~/.ailang/world/.last_approval_ref`) is refused at the scope check; mint a fresh one.

Rehearse first. This runs every fence and makes no request:

```bash
"$WORLD_BIN" publish --dry-run --store "$WORLD_STORE" --registry-origin "$WORLD_REGISTRY" --publisher "$WORLD_COMPILER" --credential-file "$WORLD_CREDENTIAL" --approval-ref "$WORLD_APPROVAL" --now 2 --expires 1000
```

Then perform the irreversible write, exactly once:

```bash
"$WORLD_BIN" publish --live --store "$WORLD_STORE" --registry-origin "$WORLD_REGISTRY" --publisher "$WORLD_COMPILER" --credential-file "$WORLD_CREDENTIAL" --approval-ref "$WORLD_APPROVAL" --now 2 --expires 1000
```

Read the outcome:

- **`PUBLISHED`** (exit 0) → done; no reconciliation.
- **`FAILED`** (exit 1) → the attempt is over and its outcome is known and recorded.
- **`INDETERMINATE`** (exit 1) → the attempt may or may not have landed publicly. Go to step 8.
  **Do not retry.** A retry here is the double-publish this whole design exists to prevent.

### 8. Reconcile an indeterminate attempt (read-only)

```bash
"$WORLD_BIN" reconcile --store "$WORLD_STORE" --registry-origin "$WORLD_REGISTRY"
```

Lists every durable publish intent with no outcome. It issues no network request; it needs no
terminal, so an autonomous loop may run it. Add `--probe` to issue the read-only metadata `GET`s
that resolve one. `--live` is **refused** by this verb rather than ignored.

Reconciliation reads the **public bucket** and never the validator service; its single network verb
is `GET`. It resolves to exactly one of four states:

| State | Meaning | Next action |
|---|---|---|
| `succeeded-reconciled` | Served metadata matches all four expected digests | Done — the publish landed |
| `conflict` | A document exists but does not match | **Stop.** Someone else won the immutable version, or the bytes differ. Human decision |
| `not-published` | Bounded repeated absence, every sample with a firing same-pass control | A live retry is permitted — **but it requires a NEW attended approval/grant** |
| `probe-unavailable` | The instrument was not shown to be working | **Stop.** Human required. This is a refusal to decide, not a third answer |

`probe-unavailable` is the default, not the exception: absence is believed only on the measured GCS
`NoSuchKey` XML document, decoded as XML rather than string-matched, **and** only when a same-pass
known-positive control (a package MEASURED to exist, fetched from the target's own key-space) returns
`200` with well-formed JSON in the same pass.

That control must travel the target's own key-space. Measured 2026-08-08: the validator origin
answers `200` with 35 KB of JSON at `/api/packages` while returning `404` at
`/packages/{vendor}/{name}/{version}/metadata.json`. So a misconfigured origin plus an index-shaped
control would make the control fire, the target read absent, the window resolve `not-published`, and
an irreversible POST be re-authorized. A control that does not travel the target's key-space proves
nothing about the target's key-space.

---

## The fences, and which one is load-bearing

`world-publish` refuses in fourteen enumerated places. They are not equally strong, and pretending
otherwise is how a fence stack rots:

| Fence | Defeats | Strength |
|---|---|---|
| `mode` | an ambiguous or defaulted intent | the irreversible path is never a default |
| `store` · `approval` · `credential` | an incomplete invocation | refuses before anything is opened |
| `packet` | publishing bytes nobody reviewed | binds the golden by name |
| `ci` | a runner that somehow got a terminal | **a declared TRIPWIRE, not the fence.** `env -u CI` defeats it |
| `confirmation` | a slip of the hand | defeats ACCIDENT. A script can echo the phrase |
| **`tty`** | **automation** | **the load-bearing layer** |

The controlling-terminal check is the one an autonomous process cannot satisfy: it requires
`/dev/tty` to open, and the phrase is read from that terminal. Measured in this repository's own
loop, `/dev/tty` fails to open with "device not configured" and stdin is a socket.

A naive `isatty` would not be enough, and this is measured too: **`/dev/null` is a character
device**, so a character-device test alone would admit `< /dev/null`. The `os.SameFile` comparison
decides where the phrase is read from: stdin only when it is the same file as `/dev/tty`, the
opened `/dev/tty` otherwise. Weakening it so that a redirected stdin is read reopens the hole
(`TestRedirectedStdinIsNeverTheConfirmation` reds).

---

## What is NOT enforced

`world/` is **convention enforcement on World, not registry namespace ownership.** The registry
validator defers namespace auth (`cmd/registry-validator/main.go:177`, "accept all publishers for
now") and authorizes with a single shared `REGISTRY_API_KEY` that is never checked against the vendor
prefix. Routed upstream as **`sunholo-data/ailang#633`**; tracked here as `8/OD-2`. Until that lands,
World enforces the vendor in its own gate and claims nothing more.
