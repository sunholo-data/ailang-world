# Attended publish runbook — `world/core`

**This runbook stops at readiness by default.** Steps 1–3 are automated, run in CI, and perform no
public write. Steps 4–7 are **attended**: they are blocked on `8/OD-1` (Mark's one-shot approval for
the irreversible first publish) and are never run headless or in CI.

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

**Stage A ends here.** Everything below requires a human.

---

## Stage B — attended publish (blocked on `8/OD-1`)

### 4. Review the ready packet

Read `scripts/world_package_ready_packet.golden.json`. It carries `package`, `version`, `exports`,
`effects`, `tarballSHA256`, `contentHash`, `interfaceHash`, `tarballBytes` and `compilerVersion`.
`compilerSHA256` is provenance about the *machine*, not the package, and is deliberately kept out
of the byte-compared golden.

Confirm all three digests against the gate's own output before continuing. The approval you are
about to mint binds **bytes, not a name**.

### 5. Mint the one-shot approval

The attended stamp is an `ApprovalRequestV1` → `ApprovalDecisionV1` pair. `Session.Invoke` traverses
`payload.approvalRef` → the landed decision → its request → the canonical scope, and refuses
**before** the credential is loaded and **before** any POST. Single use is enforced by
`approval_claims`' PRIMARY KEY — durably, not by in-memory budget — so the stamp cannot be spent
twice even across a process restart.

### 6. Invoke the publish effect exactly once

Invoke the brokered `Registry.Publish` effect. Then read the outcome:

- **resolved success / resolved failure** → done; no reconciliation.
- **typed indeterminate** → the attempt may or may not have landed publicly. Go to step 7. **Do not
  retry.** A retry here is the double-publish this whole design exists to prevent.

### 7. Reconcile an indeterminate attempt (read-only)

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

## What is NOT enforced

`world/` is **convention enforcement on World, not registry namespace ownership.** The registry
validator defers namespace auth (`cmd/registry-validator/main.go:177`, "accept all publishers for
now") and authorizes with a single shared `REGISTRY_API_KEY` that is never checked against the vendor
prefix. Routed upstream as **`sunholo-data/ailang#633`**; tracked here as `8/OD-2`. Until that lands,
World enforces the vendor in its own gate and claims nothing more.
