# V25 — design V7's alias arms re-measured on the v0.54.15 tool pin (row 156, iteration 289)

Controller measurement, 2026-10-20, darwin/arm64 (APFS). Linux is not measured.

**Question.** Iteration 288's **Next** was to re-measure V7 (design §3) on the tool binary World now
pins, and keep or drop M5c/M5d (D-156-4 = A, World's own symlink/hard-link pre-check before
`ailang-write`/`ailang-edit`) on the result. Upstream ailang#1873 (filed by World from V7) was closed
as fixed by `db9719bb7` in v0.54.15. That closure is a claim; this is the measurement.

**Binaries.** `binaries.sha256.txt`. The v0.54.15 digest `d821dad0…ac043` equals
`host/pins/ailang-releases.json` `binary_sha256` for the `tool` role, so the subject is the production
tool binary. v0.52.1 is the retired tool pin V7 measured, run here as the control arm.

**Instruments.**

- `v7r_policytool_aliases.py <bin> <scratch>`: V7's fixture and patterns (`tests/**`, `*_test.go`,
  `*.test.ts`, `conftest.py`, after `.github/**`), the three V7 alias arms with write AND edit, plus the
  hard-link arms V7 never ran (D-156-4 names both), plus an unprotected hard link. 13 arms × 2 deny
  lists. Each arm gets a fresh tree; the verdict is the host bytes read back after the call, not the
  tool's JSON alone. Outputs `v7r_v0521.out.txt`, `v7r_v05415.out.txt`.
- `v7r_doublestar.py <bin> <scratch>`: V7's pattern-form arms (`**/*.test.ts`, `tests/**` only) and
  an empty `fs_deny_write` with a single-linked and a hard-linked file. Outputs `v7rb_*.out.txt`.

**Result (both arms printed side by side in the outputs).**

| Arm (protected set present) | v0.52.1 (control) | v0.54.15 (pin) |
|---|---|---|
| lexical `pkg/a_test.go`, `tests/new.py` | refused | refused |
| symlink dir `lnk/new5.py` (`lnk → tests`, new file) | **written** `tests/new5.py` | refused, "resolves to tests/new5.py" |
| symlink dir `lnk/test_x.py` write / edit | **changed** `tests/test_x.py` | refused, "resolves to tests/test_x.py" |
| symlink file `alias.txt → pkg/a_test.go` write / edit | **changed** | refused, "resolves to pkg/a_test.go" |
| hard link `hl.txt` = `pkg/a_test.go` write / edit | **changed** | refused, "more than one hard link" |
| hard link `pkg/hl2.go` = `tests/test_x.py` write / edit | **changed** | refused, "more than one hard link" |
| control `pkg/a.go` | written | written |

Host bytes were unchanged in every v0.54.15 refusal (`host=none`). Without the protected set, every
symlink arm is WRITTEN on v0.54.15 too (control), so the refusal is the resolved-name match, not a
blanket symlink ban. Hard links are set up only in the hard-link arms, so the symlink arms' targets are
single-linked (a first run hard-linked them too and the no-pattern control was refused for the
multi-link reason; that confound was removed before banking).

**Two behaviours beyond the question.**

1. **The multi-link refusal is unconditional.** On v0.54.15 a write to ANY file with more than one
   hard link is refused, with or without a protected set, and with an EMPTY `fs_deny_write`
   (`v7rb_v05415.out.txt`, last row; single-linked control written). v0.52.1 wrote it. World
   episodes always carry a non-empty deny list (`episodeDenyWrite`), so since the row-139 D7 tool
   move (2026-10-17) `ailang-write`/`ailang-edit` refuse every hard-linked file in an episode
   worktree. That is a property of the pin, not of row 156.
2. **`**` is doublestar on v0.54.15.** `**/*.test.ts` refuses depth 0, 1 and 2; v0.52.1 refused depth
   1 only (V7). `tests/**` still does not match `cli/tests/x.py` on either binary. Design §4.1's refusal
   text for `**/x` cites V7's one-level reading, which is now v0.52.1-only.

**Not measured.** The FS effect inside `ailang-run` (R-156-3): upstream's commit adds
`internal/effects/fs_deny_symlink_test.go`, but World did not run an AILANG program here. Linux.
