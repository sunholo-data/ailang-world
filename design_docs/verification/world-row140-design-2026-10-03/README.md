# Row 140 design evidence — `Workspace.Exec` (2026-10-03)

Measurements behind `design_docs/planned/w-workspace-exec-toolchain-effect.md` §3 (V1–V35).

- **Rig:** macOS 26.6.2 arm64; node v26.0.0; go1.26.6 (module-cache toolchain); Python 3.12.13 (uv store).
- **Sandbox:** `@anthropic-ai/sandbox-runtime` 0.0.71, the fleet's install under `ailang/tools/pi-extensions/sandbox/node_modules`, used read-only. 0.0.78 was installed into scratch for V32.
- **Linux:** not measured (V33).

## Layout
- `probes/` holds the exact scripts that were run. `lib.sh` defines `run` (print the command, its output and `[rc=N]`) and needs `S` set to a scratch directory.
- `out/` holds the outputs. They were redacted after the run:
  - the scratch dir is shown as `$S` and the home dir as `$HOME`;
  - `node …/sandbox-runtime/dist/cli.js` is shown as `SRT` (or `SRT78`);
  - srt's per-run proxy credential is `<redacted>`.

## Projects and data
- **Projects:** every project was an APFS clone (`cp -c`) in scratch, never the real checkout:
  - the ailang compiler (`go.mod`, `go.sum`, `internal/`, `cmd/`);
  - TwilightGame (without `.git`, `dist`, `public`, `output`, `.claude`, `incoming-art`);
  - `sunholo-platform/cli` with its `.venv`.
- **The decoy secret** is the scratch file `~/.srt-row140-decoy`, containing `DECOY-NOT-A-SECRET`. No real secret was read.

## Probe → V-row map

| Probe | V-rows |
|---|---|
| `p1_fs` | V3, V7, V8, V11 |
| `p2_net` | V9, V10, V11, V12 |
| `p3_paths` | V4, V5, V6 |
| `p4_go` | V16–V19 |
| `p4_go_net` / `p4_go_tls` | V20 |
| `p4_ts` | V22, V8b |
| `p4_ts_fence` | V8, V8b, V23 |
| `p4_py` / `p4_py2` | V24, V25, V26 |
| `p5_overhead` | V15 |
| `p6_inject` | V28 |
| `p7_seed` | V21 |
| `p8_proc` | V11, V13, V14 |
| `p9_v078` | V32 |
| `p10_full` / `p10b_baseline` | V27 |
| `p11_deps` | V29 |
| `p12_topology` | V35 |
| `p13_quorum_r1` | V36, V37 (quorum round-1 premises) |
| `p14_quorum_r2` | V38, V39 (quorum round-2 premises) |
| `p15_quorum_r3` / `p15b_archive` / `p15c_tramp` | V40, V41, V42 (quorum round-3 premises) |
| `src_grep` | V4 (source) |

## Quorum
Quorum artifacts (`quorum/*.json`) are written by `ailang design-quorum --artifact-dir` for each round.
