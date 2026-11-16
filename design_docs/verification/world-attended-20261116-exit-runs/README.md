# Attended exit runs: 2026-11-16

The operator ran these at the Studio's console, with the controller (Claude Code, in session)
preparing commands and checking each result. All `y` answers, Touch ID checks and browser
actions were the operator's.

The build was `b0fbcd3e` (`ailang-worldd`, sha256 `5f3be220b8bc…`), built from a clean
standalone clone of `origin/dev`, with VCS stamping present and `vcs.modified=false`.

## Row 201: approvals by key (QUICKSTART §15, §16, §17). PASS

The runs used a scratch daemon on `127.0.0.1:7645` with a scratch store under `$SCRATCH/xr`,
and a private `--sessions-dir`. They never touched the home service or its stores.
`scratch-decided-and-approvers.txt` is the final `approve --list --decided` and `approver list`.

| Step | Result |
|---|---|
| §15: enrol `demo` (ssh-ed25519). The fingerprint shown matched `ssh-keygen -l` | enrolled after `y` |
| §15: `--write-message`, `ssh-keygen -Y sign`, `--signature` | `sig-demo` minted, `approved_by=ssh:demo`, no terminal |
| §15: `approver revoke demo` | `y`; listed the 1 session it approved, revoked none |
| §15: `--sign-with` (re-run with a fresh key `demo2`, because a revoked key can never be re-enrolled) | `sig-demo-2` minted, `approved_by=ssh:demo2` |
| §16: enrol `demo-mac` (passkey, origin `http://localhost:7645`). 26-char enrolment code typed. ES256, synced passkey | enrolled after `y` |
| §16: `--webauthn-link`, then the page, Touch ID, Save, then `--webauthn-assertion` | `pk-demo` minted, `approved_by=webauthn:demo-mac` |
| §16: the same assertion file a second time | refused: `already used … by webauthn:demo-mac` |
| §17: `approve --await`, then the app showed **Needs you · 1**, then the card, then **Approve with passkey** | the waiting CLI minted `surface-demo`, `approved by webauthn:demo-mac` |
| §17: **Approve with passkey** on the same card again | refused: card shows `Already used: approved by webauthn:demo-mac` |
| §17: **Decided** in the app | lists all four, newest first (operator screenshot) |
| Revoke `demo-mac`, `demo2` | both revoked |

The opt-in notification (XA4) was not exercised. It stays a per-browser check.

## Row 220: M6 follow on the home service. Merge half PASS

The operator stopped the home service, backed up both stores (`*.bak-68c29960`), and ran
`service install --replace` with `--surface-follow ~/.ailang-world/surface-follow.json`. The
config has schema `world/surface-follow-config/v1`, mode 0600, and one repo: `world`, the
`ailang-world` main checkout on `dev`. After install, `home-service-after.json.txt` shows:

- health build `b0fbcd3e`, supervisor `com.sunholo.ailang-world.worldd`;
- `doctor --quick`: no check failed;
- follow `mode: follow`, repo `world` `fresh`, open-episode reconcile `complete`.

**Merge-to-visible:** the local `dev` was 54 commits behind `origin/dev` (`8338eebb` → `a51684af`).
With the app open at `http://localhost:7644/surface`, the operator ran
`git pull --ff-only`. The app updated with no reload in **about 10 s**, against the 60 s
budget. Remote boundary (R7 = A): the follower never fetches. The clock starts at the
local-ref advance.

**Episode ring:** pending the operator's next real episode.

## Row 219: node upgrade (AC6) not evidenced by this upgrade

Homebrew moved `node` from 26.12.1 to 26.13.0_2 on 2026-11-10. Health then showed
`target=26.13.0_2, accepted=26.12.1, probe=pending` until a node call. No `twilightgame`
episode ran before the reinstall restarted the service, which re-resolved node
(`accepted=26.13.0_2, probe=ok`). AC6's "no restart" leg waits for the next node release.

## Operator-surface findings (follow-ups, not blockers)

1. **The app cannot be copied from.** `host/surfacepage/app/app.css` sets `user-select: none`
   on `body`, so the request card's id, message hash and paths cannot be selected. Operator:
   "it would be good UX to allow copy paste".
2. **The card splits each grant across three rows** (effect, scope, budget): 27 rows for 9
   grants. One `Effect scope:budget` line per grant reads better.
3. **"start record unavailable"** in *Running* for every request-minted episode (`pk-demo`,
   `sig-demo`, `sig-demo-2`, `surface-demo`). The mint-by-request path may not write the
   start record that the row 220 M3 overlay reads.
4. **`service install --replace`'s plist diff is unusable.** It printed the whole plist as
   removed and then re-added, with identical text on both sides, instead of the two added
   `--surface-follow` lines.
5. **The doctor line misleads** during install: "Publish transitions … then commit genesis
   (QUICKSTART §9)" on a store that already has a registry and a world head.
6. **`install --replace` keeps no `.prev` and resets `upgrades: []`.** Rollback relies on the
   operator's backups and `candidates/`.
7. Controller note: interactive zsh treats `#` as an argument unless
   `setopt interactivecomments` is set. Checklists for the operator must not carry inline
   `#` comments.
