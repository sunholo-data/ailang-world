---
title: Install
sidebar_position: 1
description: Build ailang-worldd and world-publish from source with Go, and install the two pinned AILANG binaries World runs.
---

# Install

World has no packaged release. You build it from source and point it at two pinned AILANG
binaries.

## Requirements

- **Go 1.26.6** (the version in `go.mod` and CI). The build fetches the Go module
  `github.com/sunholo-data/ailang` at `v0.47.2`, which provides the MCP and A2A wire
  handlers.
- **git**, **curl**, **python3** (the runbooks use it to build commit JSON and parse SSE).
- macOS or Linux. The tool-confinement facts were measured on macOS (APFS); CI runs on Linux.

## Build the two programs

```bash
git clone https://github.com/sunholo-data/ailang-world
cd ailang-world
go build -o /tmp/ailang-worldd ./cmd/ailang-worldd
go build -o /tmp/world-publish ./cmd/world-publish
```

| Binary | What it is |
|---|---|
| `ailang-worldd` | The daemon (`serve`), its bounded client verbs, and `session mint`/`revoke` |
| `world-publish` | The attended operator entrypoint: transition-registry publish and the `world/core` package publish |

`world-publish` is a separate program on purpose: nothing CI or the daemon runs is one flag
away from an irreversible write.

## The two pinned AILANG binaries

World uses two AILANG releases for two different jobs. Both are checked by version at startup.

| Name in the docs | Release | Used for | Passed as |
|---|---|---|---|
| `PIN` | **AILANG v0.41.0** (commit `24ee108`) | The interpreter that runs every transition (plan and finish phases), type-checks transition sources at publish time, and builds/verifies `world/core` | `serve --ailang-bin`, `world-publish transitions --ailang-bin` |
| `TOOL` | **AILANG v0.51.0** (commit `b99dd25`) | The binary the coding tools run inside the episode worktree (`policy-tool` and `run --policy`) | `serve --tool-ailang-bin` |

**Why two pins.**

- `PIN` is the **replay pin**. Its content hash is written into every log entry and every
  registry descriptor, so it only moves deliberately. v0.41.0 is also the compiler the
  `world/core` golden ready-packet names (`frozenCompilerVersion` in `cmd/world-publish`),
  and CI verifies the `.ail` sources against it.
- `TOOL` is a separate, archived, hash-verified binary chosen for its **confinement
  behaviour**. The tool confinement (sandbox, `.git` protection including the macOS `.GIT`
  case-fold fix that landed in v0.50.1, `--chdir`, `--args-json`) and the four `run`
  outcome shapes were measured on exactly v0.51.0. The daemon refuses any other release for
  `--tool-ailang-bin`: a different binary changes the measured facts, and World never widens to
  an unmeasured one (residual R-SE-8).

Install them at the paths the runbooks use. Release assets are named
`<os>.<arch>.ailang.tar.gz`, each with a `.sha256` file (CI uses `linux.x64`; the tool binary
was measured with `darwin.arm64`). For example, on Apple silicon:

```bash
mkdir -p ~/.pinned-ailang ~/.pinned-ailang-tools/v0.51.0
cd "$(mktemp -d)"

base=https://github.com/sunholo-data/ailang/releases/download/v0.41.0
curl -fsSL -o pin.tar.gz "$base/darwin.arm64.ailang.tar.gz"
curl -fsSL -o pin.tar.gz.sha256 "$base/darwin.arm64.ailang.tar.gz.sha256"
echo "$(cut -d' ' -f1 pin.tar.gz.sha256)  pin.tar.gz" | shasum -a 256 -c -
tar -xzf pin.tar.gz -C ~/.pinned-ailang ailang

base=https://github.com/sunholo-data/ailang/releases/download/v0.51.0
curl -fsSL -o tool.tar.gz "$base/darwin.arm64.ailang.tar.gz"
curl -fsSL -o tool.tar.gz.sha256 "$base/darwin.arm64.ailang.tar.gz.sha256"
echo "$(cut -d' ' -f1 tool.tar.gz.sha256)  tool.tar.gz" | shasum -a 256 -c -
tar -xzf tool.tar.gz -C ~/.pinned-ailang-tools/v0.51.0 ailang

export PIN=$HOME/.pinned-ailang/ailang
export TOOL=$HOME/.pinned-ailang-tools/v0.51.0/ailang
$PIN --version && $TOOL --version
```

The first line of each `--version` must read `AILANG v0.41.0` and `AILANG v0.51.0`.

:::warning Keep the pins out of `/tmp`
macOS clears `/tmp` on reboot. Older runbooks used `/tmp/ailang-v0300/ailang` (v0.30.0, the
pin before 2026-09-21); that path and version are retired. Use the durable paths above.
:::

:::tip Never use a development build
Run World against released, tagged binaries only — never a `-dirty` build and never whatever
`ailang` happens to be on `PATH`.
:::

## Optional: an examples corpus

The `examples-search` tool needs an AILANG examples corpus. It is not built into the tool
binary. `$TOOL examples download` fills `~/.ailang/examples`, which `serve` uses by default
when it exists. See [Coding tools](coding-tools.md).

## Before you run anything

Unset the registry credential in the shell that runs `ailang-worldd`:

```bash
unset AILANG_REGISTRY_API_KEY
```

With it set, every `ailang-worldd` command, client verbs included, refuses to start (exit 2).
The registry is immutable, so an ambient key would hand unrecallable publish authority to every
process World starts. Keep the key in a mode-`0600` file outside the working tree instead.

Next: [Your first world](first-world.md).
