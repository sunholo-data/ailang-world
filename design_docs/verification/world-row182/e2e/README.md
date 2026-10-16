# Row 182 AC7: end-to-end on a real daemon (attended 2026-10-16)

Fresh rig `~/.ailang/state/row182-e2e` built by `final-prepare` from the PR head; Mark published se-tools and
minted one session (`fl-cx-r1`). Tuning-ledger iteration 15 (`world-floor-m4-2026-10-06/iter-15/`): codex
(`gpt-6.1-sol`) only, paired, `config_file_parser` at the OLD wording (corpus `76a5aef`), config `config.json`.

**Capture.** Both arms FAILed for real (`logic_error`). All four graded cells (2 canaries + the task, both arms)
were recorded: `outcome_record_errors` 0. The World task's outcome names entries 6–10 and is entry 11; the
shell outcome is entry 12. Each next task's `log_from` began after the previous outcome entry. The harness's own
provenance check still passed (`prov=world`). On the served copy: `provenance --episode fl-cx-r1` → entries 2–10
(outcome entries are not coordinator entries); `why head` explains an outcome entry ("not a coordinator
invocation … world/task-outcome/v1", 2/2 links verified); `object find world/task-outcome/v1` lists all four.

**Walk.** A fresh Claude Code 2.1.300 `sonnet` walk, same channel as row 114 Q1 (World client verbs only, served
copy, `--setting-sources ""`, scratch cwd), was asked: "Why did codex fail `config_file_parser` in both its World
arm and its shell arm in this run?" It answered in **20.8 s** with 8 allowlisted calls (`why head` → the outcome
entry → `object get` of the outcomes), naming the cause for both arms. A fresh `sonnet` verifier, against an E2
covering both arms (grader rows, the grading rule, both final solutions re-executed clean vs with the agent's
`app_config.json`), returned **VERIFIED**. Envelope issue → verified: **32.8 s**. Cost $0.11. Store head and
original sha256 unchanged; 0 denials, 0 non-allowlisted calls.

**Scope.** This is a capability check, not a clause-5 qualifying walk: the question was posed by the controller
about a run made for the check, not harvested from operation (row 114 E1). For contrast, row 114 Q1 asked the same
kind of question of a store without outcome objects: 55 calls, 140 s, "could not determine", UNVERIFIED.
