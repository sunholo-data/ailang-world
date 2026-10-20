---
description: Print the ailang-worldd episode start line for this repo. It only prints; the human runs the line at a terminal and types y.
argument-hint: <episode>
allowed-tools: Bash(git rev-parse:*), Bash(ls:*)
---

Print the `episode start` line for episode `$ARGUMENTS` on the current repo. **Print only:
never run `ailang-worldd`, never pipe or feed it input, and never mint a session.**

1. The episode is `$ARGUMENTS`. It must match `^[a-z0-9][a-z0-9-]{0,63}$`; if it is empty or
   does not, ask for one and stop.
2. The repo: run `git rev-parse --show-toplevel`.
3. The profile: run `ls ~/.ailang-world-exec/profiles/`. The active profiles are the
   `<project>.json` files; a `*.proposed.json` is not active. Pick the project that names this
   repo (usually its directory name), or the only one if there is just one.
4. If an active profile fits, print:

   ```sh
   ailang-worldd episode start <episode> --repo <toplevel> --branch world/<episode> --profile <project> --addr 127.0.0.1:47330
   ```

   Add `--task <file>` if the human has a task file for the agent.

   If none fits, print the line that proposes one instead, and say that the human reviews the
   proposal and runs the `mv` it prints before starting:

   ```sh
   ailang-worldd profile init --repo <toplevel> --project <name>
   ```

5. Tell the human: **run the line yourself at a terminal and type `y`** at the y/N. It mints
   the session and prints the agent's launch line; start the agent from that line. An agent
   cannot run this step (there is no terminal for it).
