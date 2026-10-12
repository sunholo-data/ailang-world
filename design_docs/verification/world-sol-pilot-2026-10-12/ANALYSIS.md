# Codex World overhead: transcript analysis (2026-10-12, attended, zero inference spend)

Inputs: pilot iteration 11 (23 tasks, N=1, rows + 46 transcripts, `analyze.py` → `analysis.json`) and the
2026-10-09 FINAL codex rows (23 tasks × N=3). Read-only; no model calls.

## Finding: the overhead is a FIXED per-task cost, not agent behaviour

| | pilot (N=1) | FINAL (N=3) |
|---|---|---|
| median absolute World − shell wall | **+5.9 s** | **+6.6 s** (IQR +4.4..+9.2) |
| median shell wall | 30.4 s | 28.5 s |
| tasks with shell < 30 s: median o / abs | +39% / +7.8 s | +35% / +7.6 s |
| tasks with shell ≥ 30 s: median o / abs | −1% / −1.3 s | +14% / +4.6 s |
| tool calls, World vs shell (median) | 3–4 vs 4 (ratio 0.75) | 3 vs 4 |
| World tool errors | 0 | — |
| output tokens World/shell | 0.9× | — |
| input tokens World/shell | ~1.7× | 263k vs 126k |

- World does **less** work: fewer tool calls, no tool errors, fewer output tokens. It still pays a roughly
  constant +6–8 s per task.
- The bar is a **ratio**, so a constant cost is tiny on long tasks and large on short ones. The 11 tasks above
  +25% are almost exactly the short ones. The median lands wherever the 12th task's shell time falls, which is
  why the pilot (+20.0%) and the FINAL (+26.6%) differ. The fixed cost is about the same in both, so **this data
  does not show that row 174 or the client upgrade reduced it**.
- **Input tokens are quantized:** World − shell deltas cluster at about +44k, +88–96k and +133–138k, i.e.
  ~22k extra tokens on every model request. Both arms get the same teaching prompt (only the short World tools
  block differs), and the MCP tool list is ~10 KB (~2.5k tokens). So the extra ~22k per request comes from
  codex's own request construction in the World configuration (MCP-only, shell disabled), not from our prompt.
- Time and tokens do not track each other per task (`prompt_injection` +2k tokens but +15 s;
  `state_machine_vending` +1k, −6 s). So some or all of the fixed time is **per-call latency** (MCP round-trip
  → broker → capsule spawn), which these banked artefacts cannot separate: codex `--json` events carry no
  timestamps, the daemon logs no per-call latency, and the store is deliberately timestamp-free.

## What would move the bar

Passing needs the median o ≤ +25% at a median shell time around 29 s, i.e. a fixed cost ≤ ~7 s. Today it is
6–8 s, right on the edge. Halving it passes comfortably. Two zero-spend probes would split the cost:

1. **Daemon-side latency:** drive the World MCP tools directly with a script that replays a typical
   read → write → run sequence against the rig daemon (no LLM), timing each call end to end. This measures
   the broker/capsule share of the +6–8 s.
2. **Codex request size:** point codex at a local mock model endpoint (no provider call) in both arm
   configurations and capture the request bodies, to find what the ~22k tokens per request are
   (tool schemas, sandbox/approval instructions, or MCP metadata).

Either result names a product row (cut capsule/round-trip latency, or slim what codex injects) before a
new FINAL is pre-registered.

## Probes (2026-10-12, attended, zero inference spend): cause found

Scripts in `probes/`. Neither probe made a provider call. Probe 2 points codex at a local mock Responses
endpoint on 127.0.0.1. Raw request bodies are kept local, not committed, because they carry codex
client metadata.

**Probe 1, daemon latency** (`probe1.py`, `probes/probe1-result.json`, N=10 each): World MCP calls take
read 171 ms, write 169 ms, check 283 ms, run 196 ms (initialize/tools-list < 1 ms), against 18 ms natively for
`ailang check`/`run`. A typical 3–4 call World task therefore spends **< 1 s** in the daemon. **The daemon is
not the main cost.**

**Probe 2, what codex sends** (codex-cli 0.162.0, exact pilot argv plus a provider override):
- The first request is *smaller* in World (62,839 bytes) than in shell (67,081 bytes). The only difference
  is that `exec_command`/`write_stdin` are removed. There is no static per-request bloat.
- **codex 0.162 runs tools through a JavaScript `exec` "code mode" tool. In the World configuration all 8 World
  tools are *deferred*:** their schemas are absent from the request and appear only in the runtime `ALL_TOOLS`
  list (`mcp__world__ailang_read`, …), which the model must inspect via `exec` before it can call them.
  The shell arm's primary tool, `exec_command`, is declared in full in the `exec` description.
  Confirmed by a scripted mock that answered request 0 with an `exec` listing `ALL_TOOLS` (request 1 shows the
  8 World tools only there).
- These discovery `exec` calls **do not appear in `codex exec --json` output** (the probe's own `exec`
  call produced no item), which is why the transcripts show *fewer* tool calls for World.
- One full World request is about 63 KB of codex scaffolding plus the ~25k-token teaching prompt, roughly
  40–45k tokens. That matches the **+44k quantum** in the World − shell input-token deltas: each quantum is
  one extra hidden model round trip. At 2–3 s of model latency each, 1–3 extra round trips per task
  accounts for the **+6–8 s fixed cost**. This is strongly supported, not yet proven by a timed run.
- `--disable code_mode_host` and `--disable tool_search_always_defer_mcp_tools` (alone or together) leave the
  request byte-identical. The codex binary exposes no per-server config key for this. **It cannot be switched
  off from configuration.**

## Implication

The codex World overhead is a **codex tool-exposure tax** (deferred MCP tools → hidden discovery round
trips), not a World daemon or product cost. Row 174's descriptor work could not have fixed it, which matches
the data. Candidate fixes:
1. **Tuning knob (pre-FINAL, existing `world_tools_wording` knob family):** state in the World tools
   block the exact `tools.mcp__world__*` call names and argument shapes, so the model can call them
   without a discovery round trip. Equalizes with shell, where codex inlines `exec_command`. Needs one
   measured smoke iteration (inference spend; Mark's call).
2. **Upstream/fleet:** ask for a codex setting to inline (not defer) a named MCP server's tools; route as a
   fleet note, not a local workaround.
