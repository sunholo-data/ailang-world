# Row 174 attended investigation — 2026-10-10

Base: `5ca111547fdd80480a8d42578fb59942ee8a0b63` (origin/dev).
Branch: `attended/row174-codex-overhead`.

The 2026-10-09 FINAL remains FAIL. Row 174 remains open until fresh unchanged-model
smoke evidence demonstrates the required margin; this investigation is not a new FINAL.

## Historical evidence

`profile_final.py` independently hashes all 138 Codex transcripts, requires 69 distinct
task/trial rows per arm, and reports tool counts, serialized result sizes and cached input.
`final-profile.json` banks its output. Input totals are not request counts or timing data.
The transcripts have no per-call timestamps. Duplicate MCP content is upstream protocol
behavior; no local serializer fork is proposed.

## Current local daemon measurements

`daemon-probe.log` is an uncached real-binary daemon test using the pinned v0.41.0
interpreter and v0.52.1 tool binary. The reproducible probe source is
`daemon_probe_test.go.txt`. To reproduce, copy it to
`host/daemon/zz_row174_probe_test.go`, run the command below, then remove the temporary file.

```sh
AILANG_BIN=$HOME/.pinned-ailang/ailang \
WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.52.1/ailang \
go test ./host/daemon -run '^TestRow174AttendedProbe$' -count=1 -v
```

It exercises six calls each of read, write and run, with real scratch-store commits and
effects; failures are fatal. Recorded read calls 135.6–172.4 ms, writes 136.5–147.2 ms,
runs 174.1–500.3 ms. Whole-binary read+SHA256 costs 39.7–44.7 ms (interpreter) and
41.2–71.1 ms (tool). The 16,991-byte tools/list wire response has NINE tools, including
workspace-exec: it is not the floor's eight-tool response. These single-client local
measurements do not establish historical daemon timings or explain the complete wall gap.

## Routing and pending benchmark

User confirms Sonnet judge unavailable until API reset. D-WORLD-48 explicitly permits
GPT-6.1 Sol review in a fresh separate context, FLAG `same-model-fresh-context`.
Actual runtime model/token attestation is unavailable through native collaboration tools.
This is an attended development session, not a fabricated numbered mission iteration.

Scratch rig prepared without inference or authority mint at:
`/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig`.
Scratch port: 7646. User must publish the changed manifest and mint sessions themselves,
as required by `docs/QUICKSTART.md` §0. No session credential is banked in evidence.

Run all four smoke cells fresh: do not use `--reuse-from`. The existing floor config
digest does not include the package descriptions. Separately bank the product commit,
manifest SHA256, daemon binary SHA256 and served eight-tool metadata SHA256 before running.
Preserve task set, models, grader, N, thresholds, grants and teaching. Codex CLI is now
0.162.0 (historical FINAL 0.159.2); record this provenance and avoid cross-version
causal claims. No benchmark has been run in this attended session yet.

## Candidate and remaining release path

The local candidate compacted only the eight floor tool descriptions. Canonical
metadata is 13,947 -> 10,301 bytes; description prose 7,430 -> 3,794 bytes. The
served exact-eight SSE response is 14,054 -> 10,368 bytes. All non-description
fields and all 28 tracked package sources are unchanged. Scoped controls and
14 isolated mutation arms passed their required green/red expectations.

All three fresh-context GPT-6.1 Sol reviews are FAIL, 68/100. Round1 rejected
the full Go timeout failure; round2 rejected the first repair's finite-fit guard.
The revised row166 fixture has current plain10/race5, eight expected-red mutants
and exact-restored controls, plus independent judge proof. Its normal/plain test
passed at `84d3f5a`, but the complete final Go gate failed the separate unchanged
`TestExecHandlerTimeoutKillsTheGroup` (800 ms, missing partial stdout), rc1;
race NOT_REACHED. All started product hashes remained exact. The untouched parent
passed that second test 20/20 in isolation; no unique cause or waiver is established.
M1 remains unapproved; M2 has not run. All prior red and green records are retained.

A concrete one-fixture follow-up is proposed at
`design_docs/planned/w-exec-timeout-stimulus-fixture.md`. The user is being asked to
authorize one bounded follow-up beyond the shared three-round review limit.
No further repair, fourth review, merge, publish or mint occurs before that answer.

The running World driver pin is `939beb6fd1deef15296fc55f455e9e68e049659b`.
Its Codex environment-forwarder matches the fixed core implementation, and the
actual control suite returned 18 PASS / 0 FAIL / 0 SKIP. This establishes deployment
of that specific regression repair; it does not establish row174 latency efficacy.
See `codex-driver-fix-deployed.json`.

At origin/dev `5ca1115`, clauses 1/2/3/6/7 are MET; 4/5 remain UNMET. The release
critical path is row174 -> new row93 FINAL -> row114 provenance walk -> row139 attended
interpreter-epoch upgrade, with row156 protected-test-path completion before tagging. D69 showcase
prerequisites run alongside this path; row150 has the recorded fleet PNG gate.
No release percentage or ETA follows from the clause count. Historical FINAL
Codex overhead is +26.57598%, against the +25% limit; pass-rate delta -1/69 meets
its -2% limit. Claude passes. The failed FINAL remains immutable.
