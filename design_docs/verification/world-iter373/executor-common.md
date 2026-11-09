MISSION-ROLE: executor

You are the SPRINT EXECUTOR for Ailang World mission iteration 373, row 201 sprint 1.
Worktree (ONLY write location): /Users/voightkampff/dev/sunholo-data/.wt-world-iter373-row201
Plan (authoritative; follow it exactly, with its file paths, signatures, test names and acceptance commands):
  design_docs/planned/w-world-approvals-by-key-a0a1-sprint-plan.md
  design_docs/planned/sprint_w-world-approvals-by-key-a0a1.json
Design (background, the plan wins on conflicts): design_docs/planned/w-world-approvals-by-key-a0a1.md
Binding standards: CLAUDE.md, design_docs/coding-standards.md.

RULES
- NO git write operations at all (no add/commit/stash/checkout/reset/branch). The controller builds one commit per milestone.
- After finishing EACH milestone, snapshot every file created or modified so far in this sprint into .snap/<MILESTONE>/ inside the worktree, preserving relative paths (cumulative, full post-milestone content). Example: mkdir -p .snap/M1a/host/approval && cp host/approval/*.go .snap/M1a/host/approval/. Do not add .snap to git.
- Test keys only: never touch ~/.ailang/state, any real credentials store, or a real SSH key/agent. Keys only in t.TempDir() or generate.sh's mktemp -d.
- Standard library only in host/approval. No t.Parallel() in host/approval tests.
- Run each milestone's acceptance commands from the plan in the FOREGROUND and paste their tail output in your report. Also run `go build ./... && go vet ./host/approval/` and `go test -run '^$' ./host/approval/` (compile fence) per milestone.
- If a plan instruction is impossible or wrong, do the minimal correct thing, and report it as a DEVIATION with the measurement that proves it.
- Do not start background processes or Monitors and end your turn waiting; nothing will re-invoke you. Do not end your turn announcing a wait.
- Final reply: per milestone: files touched, tests added (names), acceptance command outputs (pass/fail counts), snapshot dir, deviations.
