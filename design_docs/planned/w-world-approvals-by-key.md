# w-world-approvals-by-key: approve with an enrolled key, in open formats, with the typed `y` as root

- **Status:** planned (attended design, 2026-10-21). Rulings R1–R5 are answered: **D-WORLD-89**, all at their defaults (Mark, attended 2026-10-21: "all defaults").
- **Queue row:** 201, `w-world-approvals-by-key` (attended). Sequenced after row 188's style port, because it builds on the Tauri surface's *needs you* panel (D-WORLD-88).
- **Asked by Mark** (attended, 2026-10-21):
  - "when we say mints stayed typed, is that the cli we are doing today? why is that the only option - what is a better UX for this that is possible?"
  - On 1Password: "making 1password a dependency isnt so good. how do we get something that works cross platform, or give a range of options? we can keep y as a fallback".

## §1 Problem (measured)

- **Today's mint.** Every episode starts with the operator pasting a mint script into a terminal. `ailang-worldd episode start` (row 189) prints the episode, worktree, every grant, the exec profile digest and the TTL. It then reads `y` from **`/dev/tty`**: `cmd/ailang-worldd/session.go:23`, `:43`, after `cmd/world-publish/tty.go`. That is decision D1 of [`w-session-authority.md`](w-session-authority.md), §D1.
- **Why `/dev/tty`.** It is the one channel shown to be out of an agent's reach. Agents, including the loop, run without a controlling terminal, and the operator never simulates one. Every other channel is reachable by an agent: a file, an HTTP request, a button in a page a script can click.
- **D-WORLD-87 R3** asked whether approval could also be a gesture, and ruled "no" by default, because no gesture had been shown to be unforgeable.
- **The cost, measured over row 189 and row 188 (19–21 Oct).** Twelve episodes (w189a–e and w188a–g) each needed Mark at a terminal, pasting `! sh …/MINT_BLOCK_N.sh` and typing `y`. That is the single remaining reason World needs a terminal, so it cannot be operated from the surface, a phone or another machine. The surface's *needs you* has no measured source, because nothing asks for approval anywhere but the TTY.

## §2 The idea

An approval is **a signature over the exact request, made by a key the operator enrolled**. The daemon accepts two open formats, so no vendor is required, and keeps the typed `y` as the fallback and the only way to enrol or revoke a key.

- **Proof of a human.** A key that cannot be used without the person (a biometric, a security-key touch, a passphrase) proves the human as well as a TTY does. An agent with full file access still cannot produce the signature.
- **Proof of what was approved.** The signature covers a hash of the exact request, so an approval cannot be moved to a bigger one.

### §2.1 Formats (what the daemon verifies)

| format | how the operator approves | what holds the key: their choice |
|---|---|---|
| **WebAuthn assertion** (passkeys, FIDO2), user verification required | the surface in any modern browser: Mac, Windows, Linux, iOS and Android, including the tailnet address | iCloud Keychain, Google Password Manager, Windows Hello, 1Password, Bitwarden, a YubiKey or another FIDO2 key |
| **OpenSSH signature** (`SSHSIG`, `ssh-keygen -Y sign`), Ed25519, ECDSA P-256 or `sk-*` | `ailang-worldd approve <request>` in any terminal; the Tauri app | a key file with a passphrase; any SSH agent (Secretive on the Mac's Secure Enclave, 1Password, Bitwarden, KeePassXC); a YubiKey `ed25519-sk` key, which needs a physical touch |
| **typed `y`** at the daemon's `/dev/tty` (unchanged) | the terminal | the person at the keyboard |

- **Optional native convenience.** The Tauri app may offer Touch ID (macOS), Windows Hello or polkit (Linux) through one cross-platform crate, but it **signs in one of the two formats above**, so the daemon learns nothing new.
- **Verification uses the standard library only:** `crypto/ecdsa` (P-256), `crypto/ed25519` and `crypto/sha256`. That needs two small in-repo parsers: WebAuthn authenticator data plus the COSE public key read once at enrolment (a bounded CBOR subset, no attestation, `fmt: none`), and the SSHSIG envelope. There are no new Go dependencies (`TestDaemonDependencyAllowlist`).

### §2.2 What is signed

A **mint request**, and in S2 any `Human.Approve` decision packet (`world/decision-packet/v1`, row 172), carries:
- the **request digest**: SHA-256 over the canonical request bytes (the episode, worktree, every grant with its scope and budget, the TTL, the exec profile digest, and the requesting daemon's build identity);
- a **one-time challenge**: 32 random bytes from the daemon;
- an **expiry**: 10 minutes by default.

The signed message is a fixed-prefix encoding of `("world-approval/v1", request digest, challenge, expiry)`.
- **For WebAuthn,** the message hash is the `challenge`. The relying-party ID is the surface's host name: the tailnet `*.ts.net` name, or `localhost` on the daemon's machine.
- **For SSHSIG,** the namespace is `world-approval@sunholo` and the message is those bytes.

A verified approval is consumed once (replay is refused). It mints exactly the requested session; a different grant set is a different digest.

### §2.3 Enrolment and revocation (the root stays typed)

- `ailang-worldd approver enroll --name "Mark's iPhone" --format webauthn|ssh [--pubkey <file>]` runs at the daemon's terminal and asks `y` at `/dev/tty`. WebAuthn enrolment completes in the browser: the surface shows a one-time enrolment code the operator types into the terminal prompt. **No key can enrol without a typed `y`**, so a stolen phone cannot enrol itself.
- Enrolled keys live in the **credentials store** (row 189's `credentials.db`, never under `--workspace-root`), with name, format, public key, sign counter (WebAuthn), enrolment time and last use.
- `ailang-worldd approver list` and `approver revoke <name>` manage them. Revoking is typed too.

### §2.4 The flow

1. **A request appears.** `episode start` no longer has to be pasted: the operator, or later an agent's proposal, creates a **mint request**. The daemon stores it and puts it under *needs you* on every surface, and the native app can notify.
2. **The surface shows the exact request,** in the brane view's open step: grants as pills, TTL, profile and worktree. Nothing is pre-ticked beyond the request itself.
3. **The operator approves** with whatever is at hand: Face ID on the phone (passkey), Touch ID on the Mac (native, or an SSH agent), a YubiKey touch, `ailang-worldd approve` in a terminal, or the typed `y`.
4. **The daemon verifies** the signature against an enrolled key, the digest, the challenge and the expiry. It then mints the session and records **which key approved it** in the episode's record and the decision list ("You · Minted · 12 h · iPhone passkey").
5. **Refusal is visible.** A wrong key, a stale challenge or a digest mismatch is refused with the reason, under *needs you*.

## §3 What does not change (invariants)

- **The agent never approves.** Agents have no enrolment path (typed only) and no key. An agent may *propose* a mint request in S5 ("change from here"), never sign one.
- **Publishing stays typed-only.** `world-publish` to the public, immutable registry is outside this design: one irreversible public act keeps the strictest channel.
- **The typed `y` keeps working** for every mint, on every machine with a terminal.
- **Read routes stay as open as row 199 leaves them.** Approving is a POST bound to a verified signature, never a cookie or a session.

## §4 Threats considered

| threat | answer |
|---|---|
| An agent clicks Approve in the surface | A click proves nothing; only a signature does, and the agent has no key and no enrolment path. |
| A script asks an SSH agent to sign | The key's holder prompts (biometric or touch). An unexpected prompt naming another program is the operator's cue to refuse. Holders that can bind a key to one app (Secretive, `sk` keys) are recommended. |
| Replay or reuse of an approval | One-time challenge, 10-minute expiry, consumed on use; WebAuthn sign counters checked. |
| Approving a small request, then minting a big one | The signature covers the request digest, and the minted session must equal the request. |
| A stolen phone | Revoke by name at the terminal. The phone's passkey still needs the owner's biometric. Enrolment is typed, so the thief cannot add another device. |
| A phishing page on the tailnet | WebAuthn binds the relying-party ID to the origin. The surface is the only origin enrolled. |

## §5 Milestones (S-numbering inside row 201)

- **A0:** the request object, its canonical bytes, digest, challenge, expiry and consumption, plus mint-by-request, still approved by the typed `y`. Pure Go, testable in World's sandbox.
- **A1:** SSHSIG verification, plus `approver enroll|list|revoke` and `approve`. Graded with real `ssh-keygen` signatures, including from a key-file stand-in for an agent-held key, and the refusal arms.
- **A2:** WebAuthn verification, with enrolment and assertions from recorded authenticator fixtures, plus a real browser check on the tailnet with a passkey.
- **A3:** the surface flow: *needs you*, the request card and Approve, plus a native notification and the optional Touch ID convenience in the Tauri app.
- **A4:** `Human.Approve` decision packets (row 172) approved the same way.

## §6 Rulings for Mark (D-WORLD-89)

- **R1** reverse D-WORLD-87 R3 for **mints and `Human.Approve`**: approval may be a verified signature from an enrolled key, as well as the typed `y` (default **yes**).
- **R2** the formats are **WebAuthn and OpenSSH signatures**; native app authenticators only sign in those formats (default **yes**).
- **R3** **enrolment and revocation stay typed-only** at the daemon's terminal (default **yes**).
- **R4** **publishing stays typed-only** (default **yes**).
- **R5** order: after row 188's style port, before S2's decision packets go live (default **yes**).
