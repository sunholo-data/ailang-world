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
