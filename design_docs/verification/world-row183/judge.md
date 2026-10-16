# Row 183 judge r1 — PR #281 sprint/row183-trial-runner
Score 91/100 — PASS. Blocking: none.

Suites: trial 26/26 OK; floor 201 OK (5 skipped); go vet ./cmd/ailang-worldd clean. Floor/cmd/design_docs diff vs origin/dev empty (AC6 holds); only ci.yml +4 lines (new trial step).

D4: trial.py:147-163 (argv = floor build_argv, --tools '' kept, allowedTools = 8 + mcp__world__workspace-exec); codex shell stays disabled via floor argv (test_codex_has_no_tool_allowlist_and_no_shell). Pre-flight 166-180 refuses without workspace-exec before any spawn (193 log_from=head+1 taken pre-spawn; 214 log_to after exit). 45 min deadline :54, prompt verbatim + note :138.
D5: grader runs after spawn returns (:440-443), process group + killpg on timeout (:268-290), worktree hash BEFORE graders (:440), heads byte-bounded (:225), one outcome per grader, record/payload failures counted not raised (:449-459), secrets stripped from grader/worldd env (:232), session by file path (:354). AC2/AC3 ok: pr-body never calls gh, refuses with no trailer or token-shaped text (:476-506).
Naming trial.py vs run.py, and compileOk=runtimeOk=stdoutOk=verdict: neither matters (schema has no split, validator accepts; README documents it). Amend design text D4 name.

Non-blocking:
1. trial.py:247-254/261 run_grader: sha256_file on an executable that is not readable (mode 0111) raises PermissionError out of run (not caught, unlike spawn). Wrap tool hashing in try/except OSError.
2. trial.py:275 communicate buffers unbounded grader output in memory and writes it whole to evidence/graders; cap it.
3. grader cwd check (:128) is lexical; a symlink inside the agent-written worktree can escape. Graders are unsandboxed by design, so low risk; realpath-check optionally.
4. Timeout kill uses killpg but no test spawns a grandchild (M7 survives); add one.
5. Token never reaches argv (config file / env only), so redact_argv is defensive; M11 survives harmlessly.
6. worktree_state runs git in an agent-controlled worktree (config/attributes could run drivers) — same trust class as graders.

Mutations (15): killed M1 allowedTools (test_claude_allows_workspace_exec_plus_the_eight), M2 preflight check (test_preflight_without_workspace_exec_spawns_nothing), M3 log_from no +1 (test_claude_episode_end_to_end), M4 token in record argv (FakeWorldd/validator tests; weak, no dedicated token-in-argv assertion), M5 grader env secrets (test_graders_get_no_secrets), M6 hash after graders (end_to_end), M8 record errors uncounted (test_record_failure_is_counted_not_raised), M9 char-bounded heads (test_byte_bounded_heads), M10 recorder exception propagates (test_record_failure...), M12 pass on timeout (test_grader_wall_cap), M14 mcp config left (end_to_end), M15 log_to before spawn (end_to_end), M16 cwd traversal (test_bad_specs_are_refused). SURVIVED M7 (kill vs killpg), M11 (no argv redaction).
