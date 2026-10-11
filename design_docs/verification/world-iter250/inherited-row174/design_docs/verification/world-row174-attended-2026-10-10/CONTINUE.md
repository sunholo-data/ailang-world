# Attended benchmark continuation

**BLOCKED before step 1:** M1 final normal Go gate is red, round3 FAIL68.
The separate exec-timeout fixture follow-up needs attended authorization and
passing normal gates/new review before publication, mint or performance runs.


No benchmark inference, publish, or session mint has run in this investigation. The scratch
store contains bootstrap state only. Its binaries and worktrees are prepared.

1. Finish local candidate checks and fresh-context review before using these steps.
2. User runs `bash attended-mint.sh` in a real terminal (absolute script path below).
   The script publishes the **original parent** manifest from the detached baseline
   checkout and asks for two scratch-session confirmations. It never types an answer.
3. Agent runs all four fresh baseline cells with unchanged `floor_config.json`:

```sh
cd /private/tmp/ailang-world-row174-baseline-20261010
python3 scripts/floor/run.py iterate --iteration 6 --knob row174-descriptor-baseline \
  --side baseline --change 'fresh original manifest baseline on current CLI versions' \
  --root /Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig \
  --addr 127.0.0.1:7646
```

Before inference, capture the exact eight-tool live metadata plus manifest, binary,
CLI and config provenance. Bank the complete iter-6 directory and resulting ledger
row in the candidate branch, commit that ledger, and only then proceed. Check the
ledger head first: if another session has advanced it, allocate the next unused
iteration rather than overwriting a run. Abort on existing capacity rules; a partial
run is not a performance pass. Do not reuse cells.

4. After the baseline, user runs `attended-publish-candidate.sh` in a terminal. It
   publishes the candidate manifest with its own typed confirmation. Sessions can
   remain valid within their original TTL/budgets; expiry requires a new typed mint.
5. Bank candidate product/served metadata provenance; run the next unused iteration,
   all four cells fresh, `--side world --knob row174-tool-descriptions`. Commit all
   evidence and the ledger. Then run one unchanged-candidate all-four-cell repeat
   at the next unused iteration. Keep the task set, cheap models, N, grader, teaching,
   thresholds, deadline and grant set unchanged. No `--reuse-from`.
6. Evaluate every attempt against the plan: Codex overhead <=20% tuning target,
   Claude <=25%, deltas >=-2%, no harness faults. Report noise and failed attempts.
   Missing or failing evidence keeps row 174 open. A later FINAL needs a new committed
   preregistration, normal sealing protocol and human-typed FINAL sessions after landing.

Initial human command:

```sh
bash /private/tmp/ailang-world-row174-20261010/design_docs/verification/world-row174-attended-2026-10-10/attended-mint.sh
```

Authority source: `docs/QUICKSTART.md` section 0 says these two commands read a line
“you type” and “Run them yourself”. This cannot be replaced by an agent-created PTY
or scripted answers. Sonnet is unnecessary for the permitted fresh-context judge;
the benchmark's fixed Haiku/Luna model lanes still need capacity.
