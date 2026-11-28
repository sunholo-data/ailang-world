# Row 274 S0, part 1: can Blender run under World's exec sandbox? (attended, 2026-11-28)

**Result: no, not on macOS with srt 0.0.78.** Blender 5.2.1 crashes at startup in every srt
configuration tried, including an open filesystem with every Mach service allowed. srt's macOS
profile allows IOKit user clients only for IOSurface, and Metal needs the GPU's
`AGXDeviceUserClient`. macOS Blender has no non-Metal GPU backend, and it crashes on a null
device instead of falling back. So D-WORLD-123's options A and B (`Workspace.Exec` under srt)
cannot run Blender on the Studio. Only C (a host-mediated effect) or a new option D (a
World-authored sandbox profile that admits the GPU) can.

Host: Mac Studio (`Voights-Mac-Studio.local`, Darwin 25.6.0). Blender `5.2.1 LTS` (hash
`9e2066aef7ef`) at `/Applications/Blender.app`. srt 0.0.78 at
`~/.ailang-world-exec/srt/node_modules` (World's pinned copy). Live `ailang-worldd` from
`~/.ailang-world/bin`. Nothing here touched the live service, its store, its profiles directory or
its presets. All scratch state was under the session scratchpad.

## R1. Outside any sandbox it is fast

```
Blender -b --factory-startup --python-expr 'import bpy; bpy.ops.mesh.primitive_torus_add();
  bpy.ops.export_scene.gltf(filepath="probe.glb"); print("PROBE_OK")'
→ PROBE_OK, real 1.65 s; second run real 0.61 s; probe.glb 83,280 bytes
```

## R2. Through World: the startup probe of a scratch daemon

A scratch `ailang-worldd serve` (own `--db` and `--credentials-db` in a scratch state dir, a
scratch `--workspace-root`, `--bind 127.0.0.1:7699`, World's srt, `--exec-profile` holding one
profile `blender`, with `read_roots: ["/Applications/Blender.app"]`, `timeout_ms: 9000`, and the
probe `Blender --background --factory-startup --python-exit-code 1 --python-expr "import bpy; …"`).

1. With the profile inside the state dir, startup refused it, correctly and clearly: *"the
   confinement stack (srt, node, libnode, profiles, cache seeds) must lie outside the workspace
   root, the state dir and every exec cache"*. The profile was moved.
2. Then:
   ```
   daemon startup failed at config: the exec configuration is refused: exec startup probe failed
   for profile "blender" (1 of 8 arms): arm7-toolchain: probe "…Blender --background …": exit 139;
   stderr "ArchWarn: ARCH_CACHE_LINE_SIZE != Arch_ObtainCacheLineSize() …"
   stdout "Blender 5.2.1 LTS … Writing: <state>/exec-cache/wep-…/blen…"
   ```
   Exit 139 is SIGSEGV. World refused to start rather than serve a broken tool. The error names
   the profile, the arm, the exit code and the stderr, which meets row 227's "clear errors" bar.

## R3. srt alone, without World, at three strictness levels

`node <srt>/dist/cli.js -s <settings> Blender --background --factory-startup --python-exit-code 1
--python-expr "…torus…gltf…PROBE_OK"`, with HOME pointed at a scratch worktree:

| Settings | Result |
|---|---|
| `world`: denyRead `$HOME` and the World state dir, allowRead the worktree and `/Applications/Blender.app`, allowWrite the worktree, no network | rc 1, crash |
| `open`: no denyRead, allowWrite the worktree, no network | rc 1, crash |
| `mach`: as `open`, plus `network.allowMachLookup: ["*"]` | rc 1, crash |
| `--gpu-backend opengl` / `vulkan` | refused: *"Unrecognized GPU backend … expected one of [metal]"* |

The crash is the same every time (`blender.crash.txt`):

```
5   libsystem_platform.dylib  _platform_strstr + 48
6   Blender  blender::gpu::supports_barycentric_whitelist(id<MTLDevice>) + 64
7   Blender  blender::gpu::MTLBackend::metal_is_supported() + 364
8   Blender  blender::GPU_backend_type_selection_detect() + 476
9   Blender  blender::wm_homefile_read_ex(...) + 1812
10  Blender  blender::WM_init(...) + 304
```

## R4. The denial that causes it

`log show --last 2m --predicate 'eventMessage CONTAINS "Sandbox" AND eventMessage CONTAINS
"Blender"'` during the `mach` run:

```
deny(1) iokit-open-user-client AGXDeviceUserClient
deny(1) mach-lookup com.apple.DiskArbitration.diskarbitrationd
deny(1) mach-lookup com.apple.CoreServices.coreservicesd
```

srt 0.0.78's macOS profile (`dist/sandbox/macos-sandbox-utils.js`, "IOKit - specific operations
only") allows `iokit-open` only for `IOSurfaceRootUserClient`, `RootDomainUserClient` and
`IOSurfaceSendRight`. No setting extends it: `sandbox-config.js` exposes `allowMachLookup` but no
IOKit or GPU setting. Without `AGXDeviceUserClient` there is no Metal device. Blender reads the
missing device's name without a null check and segfaults, even in `--background`.

## What this changes

- **D-WORLD-123:** A and B are ruled out on macOS by measurement. Remaining are **C**, a
  host-mediated effect (`Host.Blender`, row 272's pattern), where World runs a fixed command
  outside srt with worktree-only inputs and outputs, and **D**, a World-authored sandbox profile
  that adds the GPU user client to srt's rules. D needs either an upstream srt setting or World
  rendering its own Seatbelt profile.
- **Not filed:** a Blender bug (a crash instead of a fallback when no Metal device exists) and an
  srt feature request (a GPU/IOKit setting). Both are external projects. Filing either is Mark's
  call.
- **S0 part 2** (a real attended episode in `~/dev/blender`) waits for the D-WORLD-123 ruling.
  Under the old plan it would only hit this wall again.
