# Attended publish runbook — `world/core`

**This runbook stops at readiness by default.** Steps 1–3 are automated, run in CI, and perform no
public write. Steps 4–8 are **attended**: they are never run headless and never run in CI.

The asymmetry that shapes every refusal below: a public publish is **immutable**
(`registry-validator` 409s on an existing version), so a wrong publish permanently consumes
`world/core@0.1.0`, while a refusal costs a human five minutes.

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

This gate needs the **pinned v0.30.0** compiler, which is *not* the binary the other legs use. Point
it with `WORLD_PKG_AILANG_BIN`; there is deliberately **no silent fallback**, so a wrong binary fails
loudly rather than skipping.

### 3. Run the full repo gate

```bash
./scripts/verify_ail.sh   # invokes verify_world_package.sh at :224
./scripts/verify_go.sh
```

`verify_go.sh` prints **two `WARNING: DATA RACE` lines on a healthy run** — that is its own
known-positive control, and it FATALs if they are absent. Do not "fix" them.

**Stage A ends here.** Everything below requires a human at a terminal.

---

## Stage B — attended publish

Every command in this stage is refused unless a human is present. The refusal is **structural, not
advisory**: `world-publish` opens `/dev/tty` and requires stdin to be *that same file*. A pipe, a
redirect, a CI runner and an autonomous agent all fail it, and none of them can arrange otherwise.

`--dry-run` is a rehearsal that passes every one of those fences and then makes no request at all,
so an operator can walk the exact keystrokes before the real one.

**Exit codes:** `0` done · `1` failed (or indeterminate — see step 8) · `2` usage ·
`3` **STOP**, printed as `STOP fence=<name>` on stderr, meaning *nothing happened and nothing will*.

### 4. Review the ready packet

Read `scripts/world_package_ready_packet.golden.json`. It carries `package`, `version`, `exports`,
`effects`, `tarballSHA256`, `contentHash`, `interfaceHash`, `tarballBytes` and `compilerVersion`.
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

So the eyeball check is **document against artifact, at full length**. These three digests are the
identity of what you are about to publish permanently:

| field | digest |
|---|---|
| `contentHash` | `sha256:0c8c60616e592dc01891e8bbb59350786f242a2f79a9eb2c587ae8b0ca2e00b9` |
| `interfaceHash` | `sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083` |
| `tarballSHA256` | `sha256:d4ff710e4850cd7009ce37a01c77c7cd21576b6bae176b52b5de76e49d0506f7` |

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

```bash
export WORLD_STORE="$HOME/.ailang/world/world.db"
export WORLD_BIN="$(mktemp -d)/world-publish"
export WORLD_REGISTRY="https://storage.googleapis.com/ailang-registry"
export WORLD_CREDENTIAL="$HOME/.config/ailang/registry.key"
export WORLD_COMPILER="/tmp/ailang-v0300/ailang"
```

```bash
go build -o "$WORLD_BIN" ./cmd/world-publish
```

The binary is built into a temp directory on purpose. `go run` is **not** used here, and the reason
is measured: on go1.25.6 `go run` exits `1` for a child that exited `3`, printing `exit status 3` to
stderr instead of propagating it. The STOP contract is an exit code, so the runbook uses a binary
whose code survives.

`WORLD_CREDENTIAL` must name a mode-`0600` file **outside the working tree**. The API key is never
read from the environment: if `AILANG_REGISTRY_API_KEY` is set in your shell, the production handler
refuses to be constructed at all.

### 7. Mint the one-shot approval, then spend it

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
twice even across a process restart. The scope binds the three digests above, so an approval minted
for these bytes authorizes no other bytes, and an approval for `0.1.0` cannot authorize `0.1.1`.

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
| `succeeded-reconciled` | Served metadata matches all three expected digests | Done — the publish landed |
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
`/dev/tty` to open *and* stdin to be the same file as that device. Measured in this repository's own
loop, `/dev/tty` fails to open with "device not configured" and stdin is a socket.

A naive `isatty` would not be enough, and this is measured too: **`/dev/null` is a character
device**, so a character-device test alone admits `< /dev/null`. The `os.SameFile` comparison is the
repair, and deleting it reopens the hole.

---

## What is NOT enforced

`world/` is **convention enforcement on World, not registry namespace ownership.** The registry
validator defers namespace auth (`cmd/registry-validator/main.go:177`, "accept all publishers for
now") and authorizes with a single shared `REGISTRY_API_KEY` that is never checked against the vendor
prefix. Routed upstream as **`sunholo-data/ailang#633`**; tracked here as `8/OD-2`. Until that lands,
World enforces the vendor in its own gate and claims nothing more.
