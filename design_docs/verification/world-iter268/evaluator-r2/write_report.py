from pathlib import Path
import json,time,hashlib,subprocess,shutil
W=Path('/Users/voightkampff/dev/sunholo-data/.eval-world-iter268-r2');E=W/'design_docs/verification/world-iter268/evaluator-r2';start=json.loads((E/'start.json').read_text());post=json.loads((E/'post.json').read_text());focused=json.loads((E/'focused-gates.json').read_text());attacks=json.loads((E/'attacks.json').read_text());restored=json.loads((E/'restoration.json').read_text());assert restored==start['files']
ci=json.loads((E/'ci-observed.json').read_text());full=Path.home()/'.ailang/state/world-iter268/controller/candidate-7ba691b';full_status={}
for n in ['candidate-ail','candidate-go']:
 p=full/(n+'.json')
 if p.exists():
  shutil.copyfile(p,E/(n+'.json'));full_status[n]=json.loads(p.read_text())
 else: full_status[n]={'status':'IN_PROGRESS at observation; no terminal JSON','observed_epoch':time.time(),'log':str(full/(n+'.log'))}
for n in ['candidate-ail.log']:
 if (full/n).exists():shutil.copyfile(full/n,E/n)
# Complete source-refusal inventory, distinguished from proven mutation coverage.
inventory={
 'runDemo':{'mutated':['unknown subcommand','positional arguments','required flags','90s overall budget (SURVIVES)'],'existing_helpers_exercised':['parseVerbFlags help/flags','resolveSession absent token'],'propagation':['demoSeedDomain error'],'success_boundary':'explicit incomplete exit1/no manifest, authenticated test'},
 'demoGet':{'mutated':['HTTP read error','absent required GET']},
 'demoPreflight':{'mutated':['existing manifest','DB/health identity','pin read','pin hash','archive unavailable','archive Hash','archive Size','archive/health Version equality','v0.41.0 version prefix','nonempty world/head/log','registry object integrity','registry wire decode','epoch nomination'],'propagation_not_per_site_mutated':['demoGet health','demoGrants','headState read','log page read','demoGet registry','object read'],'unprovoked_operational':['filepath.Abs failure for db','filepath.Abs failure for health db']},
 'demoGrants':{'mutated':['session denial','episode mismatch','live sufficient grants','subtract persisted spend','world-row scan threshold'],'diagnostic_only_mutation':['OpenReadOnly failure message (wrong refusal assertion; not a bypass kill)'],'unprovoked_operational':['OpenReadOnly actual error path is exercised by readonly-db','ResolveContext SQLite query error','EffectSpend SQLite query error','ScanUnreadableWorlds SQLite query error'],'defensive_contract':'nil success/denial pair not injected; resolver promises either denial or success'},
 'demoSeedDomain':{'mutated':['source read','CAS mismatch','commit transport uncertainty','HTTP status','response selectedHead/JSON guard','bearer'],'propagation_not_per_site_mutated':['preflight error','headState read'],'nonrepresentable_supported_input':'json.Marshal of fixed strings/int64/slices/bytes struct cannot fail through supported input; no arbitrary interface/float/channel field'},
 'supporting_diff':{'mutated':['frozen interpreter header','semantics epoch header','carrying transition reference','new documentation artifact command'],'observed_survivors':['StateRoot derivation changed to hash(o.Hash) only: authenticated test remains green','entryHash changed to hash(append(encoded,0)): authenticated test remains green','90s command budget changed to24h: all four TestDemoSeed top-level tests remain green'],'main_dispatch_help':'own-M2a production revert kills each three named M2a acceptance tests and extra commit-failure test; earlier mint regression stays green','scope':'No blanket mutation completeness, no whole-repo test coverage, no direct Linux certification.'}}
(E/'refusal-support-inventory.json').write_text(json.dumps(inventory,indent=2))
score={'tests_pass':0,'lint_clean':10,'M2a_acceptance':30,'code_quality':10,'documentation':10,'design_fidelity':10}
report={'sprint_id':'w-marketing-capture','iteration':268,'milestone':'M2a','evaluation_round':2,'result':'FAIL','score':sum(score.values()),'blocking_count':1,'head':start['head'],'base':'39066a46c1ac13c994bacb162b8d4b754d5ec083','previous_candidate':'67ee2efe5b2fdca3dcab4eaeda4849a8caea1c8c','judge':{'model':'gpt-6.1-sol','role':'independent native Agent, fresh context; did not author implementation or r1 verdict','judge_independence':'same-model-fresh-context FLAG under D-WORLD-48/current Gate3','fallback':'native Sonnet Unknown model; Anthropic/OpenRouter over ration per routing evidence, not re-probed','actual_tokens':'unknown'},'platform':'darwin/arm64; Go1.26.6; v0.41.0 capsule and v0.52.1 tool pin unless explicitly unset toolpin arm','blocking':[{'id':'R2-B1','file':'cmd/ailang-worldd/demo_test.go','line':394,'finding':'New pin-mismatch refusal fixture requires WORLD_TOOL_AILANG_BIN; ordinary CI installs AILANG_BIN only. Exact repaired head fails with tool variable absent.','command':['go','test','./cmd/ailang-worldd','-run','^TestDemoSeedPreflightRefusals$/^pin-mismatch$','-count=1','-timeout=90s','-v'],'rc':1,'evidence':'evaluator-r2/ci-env-unset-toolpin.log','fix':'Create byte-distinct mismatching pin test data in t.TempDir; hash preflight rejects it before execution. Preserve exact refusal/inventory assertions and test both environment arms; do not skip or edit CI gates.'}],
 'r1_dispositions':{'R1-B1':'CLOSED: direct registry import removed; unchanged dependency allowlist/recording writer tests PASS; local wire schema equivalent to registry.Decode; decode mutant kill','R1-B2':'CLOSED: every prior QUICKSTART byte preserved; section12 uses private DEMO_BIN consistently; section11 and decision-packet doc regressions PASS; wiring mutant kill','R1-B3':'CLOSED: direct os/exec launcher removed; existing read-only manifest path reuses startup bounded scrubbed probe; unchanged nine-site driven/scrub test PASS; all archive metadata refusal mutants compile and die'},'rubric':score,'conditional_categories':{'shared_compilation':'not triggered','performance':'not triggered'},'focused_gates':focused,'post_commands':post,'attacks':attacks,'initial_mutation_assertion_mismatch':'commit-response first expected no retries but actual named kill is commit refusal boundary; rerun same compiling mutant confirms actual assertion. One unique mutant only. Initial manifest replacement-count0 was an instrument refusal before writing, excluded.','mutation_summary':{'distinct_assertion_killed_mutants':34,'of_which_diagnostic_only':1,'distinct_timeout_survivor':1,'additional_supporting_survivors':2,'build_failures_counted_as_kills':0,'old_29_executor_mutants_used_as_current_own_proof':False},'refusal_support_inventory':inventory,'own_diff_nonvacuity':{'build_rc':0,'log':'evaluator-r2/M2a-own-revert.log','production_revert':'remove demo.go/demo_seed.go and restore only main.go to M2a base, keeping tests/current docs and earlier M1a/M1b source','M2a_named_tests':'AuthenticatedCommits, PreflightRefusals, HelpPreparation all FAIL; extra CommitFailures also FAIL','prior_regressions':'SessionMint_RequiresControllingTerminal stays PASS; REST decision-packet prohibition is earlier-work regression protection, not M2a new behavior credit'},'full_signals':full_status,'ci_observation':ci,'restoration':restored,'protected_paths_observation':'44 tracked protected paths SHA256/newline observations match head and base; no writes/restores to them','scope_limits':['Only M2a five implementation files judged; M1a/M1b prior work; M2b onward intentionally pending','r1-r3 design quorum remain BLOCKED; r4 narrow-refinement disposition inherited, no new quorum credit','No full seed, release, image, endpoint, Mint, package publication or attended run certified','Unprovoked OS/SQLite errors not all independently forced; no declaration of cross-platform unreachability','Archive sidecar is existing startup trust/diagnostic material, not a new signature or independent version execution; crosschecks establish requested local-byte/health agreement only','Controller both-pin full Go result cannot discharge ordinary-CI fixture failure; no duplicate full suite run'], 'elapsed_seconds_since_evidence_start':round(time.time()-start['start_epoch'],3),'tokens_actual':'unknown'}
(W/'design_docs/verification/world-iter268/evaluator-r2.json').write_text(json.dumps(report,indent=2)+'\n')
rows='\n'.join('| '+x['name']+' | '+str(x['rc'])+' | '+str(x['seconds'])+' |' for x in focused)
md=f'''# Independent evaluator round2 — World268 / row150 / M2a

**FAIL — 70/100; 1 blocking finding.** Subject `{start['head']}`, base `39066a46c1ac13c994bacb162b8d4b754d5ec083`, previous candidate `67ee2efe5b2fdca3dcab4eaeda4849a8caea1c8c`. The three r1 blockers are closed; the CI fixture dependency below blocks this repaired head.

Independent native Agent judge `gpt-6.1-sol`, fresh context, neither author nor prior judge. **judge-independence: same-model-fresh-context FLAG**, permitted by D-WORLD-48/current Gate3 after native Sonnet Unknown model and Anthropic/OpenRouter ration limits documented in routing.json. Tokens actually used: unknown. Local measurements are darwin/arm64, Go1.26.6, released v0.41.0 capsule/proof pin and separate v0.52.1 tool pin. No Linux certification is inferred.

## Blocking finding

**R2-B1 — `cmd/ailang-worldd/demo_test.go:394`: ordinary CI cannot satisfy the new pin-mismatch fixture.** The test requires `WORLD_TOOL_AILANG_BIN` and calls Fatal if absent. Independently inspected `.github/workflows/ci.yml`: the Go job installs/exports `AILANG_BIN` only and has no WORLD_TOOL_AILANG_BIN provisioning. Independently ran on pristine exact7ba, removing only that tool variable beyond the standard four telemetry removals:

```sh
go test ./cmd/ailang-worldd -run '^TestDemoSeedPreflightRefusals$/^pin-mismatch$' -count=1 -timeout=90s -v
```

**rc1,2.740s**, exact assertion `demo_test.go:394: required tool pin absent`. Raw evidence: [ci-env-unset-toolpin.log](evaluator-r2/ci-env-unset-toolpin.log). This is a newly introduced test/CI portability failure, not evidence that the CLI's hash rejection is wrong. A full local PASS with both pins set cannot override it. Repair the fixture with byte-distinct temporary pin data, which this hash refusal reaches before any execution; preserve refusal/inventory assertions and verify both environment arms. No skip or workflow modification is needed.

## Prior findings and authority target

- **R1-B1 CLOSED:** forbidden direct host/registry import removed. Existing dependency allowlist and recording-writer tests pass unchanged. The anonymous JSON wire shape preserves the former Decode behavior; fresh decode mutant compiles and dies on wrong-refusal assertion.
- **R1-B2 CLOSED:** all 75,634 prior QUICKSTART bytes are preserved as an exact prefix. New section12 uses its private DEMO_BIN export/build/mint/serve/seed consistently and transfers variables to the serving terminal. Both old §11 binding and `TestQuickstartDecisionPacketSection` pass unchanged; the latter independently rc0,1.175s. New wiring mutation fails its named assertion.
- **R1-B3 CLOSED:** seed's direct subprocess is removed, not exempted from the audit. Existing nine-site driven credential-scrub audit passes. `archive.New(matchedDB).ReadManifest(pinRef)` only reads; startup archive already probes with timeout and childenv.Scrubbed. Fresh missing/hash/size/version-equality/release-prefix mutants compile and fail intended assertions. No credential probe was necessary because the offending launcher no longer exists.

Read-only clarification is sound within its stated scope: health DB identity first, bounded authority resolution binds the opaque credential to episode demo, EffectSpend subtracts prior charges, broker predicates check required live grants, and ScanUnreadableWorlds(empty-afterRef,1).Scanned rejects a real stored unselected World even with empty head/log. Deferred read-only Close happens before REST commits. No writable store, Mint or new endpoint is added. Every public write retains bearer and observed-head CAS. Bootstrap session/registry rows remain allowed; no all-object emptiness claim. Manifest checks bind Hash/Size to local hash-matched bytes and Version to health plus supported release; they reuse the existing startup sidecar trust, not new signature/authenticity authority.

## Own focused gates and mutations

All eight focused commands passed with both pins set; top-level names were discovered, no skips:

| Command/test group | rc | Seconds |
|---|---:|---:|
{rows}

Exact argv are in [focused-gates.json](evaluator-r2/focused-gates.json). Gates use explicit 90s/120s Go bounds and150s outer bounds; children remove only GOOGLE_CLOUD_PROJECT, OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_RESOURCE_ATTRIBUTES, OTLP_GOOGLE_CLOUD_PROJECT. The single-variable CI experiment is separately named above.

Fresh exact-head drills: **34 distinct compiling assertion kills**, including one diagnostic-only wrong-refusal mutation of readonly-open (not a bypass proof). Each attack records landed SHA256, build rc0, named assertion and byte-exact restore. The initial commit-response expected-assertion label was wrong: actual assertion is `commit refusal boundary`; the same mutant was rebuilt and redriven, counted once. A manifest replacement-count0 instrument refusal happened before any write and is excluded. No build-red, panic, timeout or missing test counts as a kill; the old executor29/r1judge20 are not claimed as my current proof.

**Own-M2a diff non-vacuity:** removed the two new production files and restored only main.go's production delta to the M2a base, leaving current tests/docs and M1a/M1b intact. Build rc0. All three named M2a acceptance tests fail; extra commit-failure test fails too. Earlier `TestSessionMint_RequiresControllingTerminal` remains PASS, correctly classed as regression protection; earlier REST packet prohibition is likewise not new M2a behavior credit. [Own-revert log](evaluator-r2/M2a-own-revert.log).

Diff enumeration is exhaustive by path:181 base→head entries (five implementation, two plan clarifications,174 verification/routing),107 previous-candidate→head entries. Full patches/path inventory are banked; no omitted implementation file is assumed from controller orientation. Refusal/support inventory explicitly lists unprovoked filepath.Abs, SQLite resolver/spend/scan and per-call read-error propagation, defensive resolver nil result, and nonrepresentable fixed-primitive JSON marshal errors. These are limitations, not blanket unreachable claims or current independent kills.

Three measured **nonblocking coverage survivors**: replacing90s overall budget with24h leaves all four TestDemoSeed top-level tests PASS; replacing StateRoot derivation with hash(o.Hash) and replacing entryHash with hash(append(encoded,0)) each leaves authenticated test PASS after build0. These identify unpinned supporting behavior; they do not establish violation of a normative canonical hash formula. No suite-wide/sole-killer claim is made for them. Queue them under3n rather than widening this repair. Details/raw logs in JSON and refusal-support inventory.

After every attack, source was restored from external byte copies. Final restored authenticated/preflight/help/commit-failure tests pass, formatting/vet/scoped code-doc diffcheck pass. All five file hashes equal reviewed HEAD.44 protected tracked paths retain head/base SHA256 and newline equality; none was written/restored. Raw-log whitespace remains unaltered.

## Full signals, rubric and limits

Exact-head PR267 CI observation is banked: AILANG job SUCCESS and Go job IN_PROGRESS at observation (no pending result called PASS). Controller exact7ba AIL gate rc0 is copied; fullGo status at final snapshot is in JSON. I did not duplicate full suites. The independently reproduced ordinary-CI test failure already makes Tests Pass0/HARD FAIL, regardless of a both-pin local result. Historical base green and old67ee fullGo red remain scoped to those heads.

| Category | Score | Basis |
|---|---:|---|
| Tests pass | 0/20 | HARD FAIL: new test fails ordinary-CI environment |
| Lint clean | 10/10 | affected-package vet, formatting and scoped diffcheck green |
| M2a acceptance | 30/30 | three narrow behavior criteria independently exercised with both pins; CI defect separately blocks |
| Code quality | 10/15 | capped-5 new functions >50lines; no >800line file or new TODO/HACK/FIXME |
| Documentation | 10/15 | no changelog0/5; runnable generator/runbook5/5; honest incremental status5/5 |
| Design fidelity | 10/10 | repaired architecture and named read-only clarification match narrow M2a |
| Conditional compiler/performance | N/A | neither category triggered |
| **Total** | **70/100 FAIL** | **1blocking** |

Only M2a is judged. M1a/M1b were already landed; M2b+ remain intentionally pending, without false incomplete-sprint deductions. Quorum r1–r3 stay BLOCKED; r4 inherits iteration266 narrow-refinement disposition, not a new quorum PASS. No whole seed/demo, proof stage, images, capture, release, publication, attended mint or promotion credit. No git writes, inbox calls, production publish/mint, confirmations or subagents were used. All own commands are terminal and source restored.

EVALUATION_RESULT: fail
EVALUATION_SCORE: 70/100
EVALUATION_ROUND: 2
EVALUATION_REPORT_PATH: design_docs/verification/world-iter268/evaluator-r2.json
FEEDBACK_SUMMARY: Replace tool-pin environment dependence in pin-mismatch fixture; preserve all assertions and verify ordinary-CI unset arm.
'''
(W/'design_docs/verification/world-iter268/evaluator-r2.md').write_text(md)
print('REPORT TERMINAL FAIL70/100 1blocking',report['elapsed_seconds_since_evidence_start'])
