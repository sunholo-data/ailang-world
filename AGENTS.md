# Working in AILANG World

Read [CLAUDE.md](CLAUDE.md) for the project charter, coding standards, and verification commands.

## Attended World commands in Herdr

When a World command needs a person to answer `Proceed? [y/N]` and Herdr is available, open a visible tab in the operator’s World workspace. The agent prepares and runs the command. The operator reviews the change and types the answer in that tab.

1. Run `herdr workspace list` on the machine hosting World. Find the `ailang-world` workspace and use its returned ID.
2. Run `herdr tab create --workspace "$workspace_id" --cwd "$repo" --label "Deploy World" --focus`. Choose a label that names the operation. The response includes the tab number and `root_pane.pane_id`.
3. Start the requested command in that pane. For a service upgrade, `candidate` is the absolute path to the verified binary on that machine.

   ```sh
   herdr pane run "$pane_id" "$HOME/.ailang-world/bin/ailang-worldd" service upgrade --binary "$candidate"
   herdr pane read "$pane_id" --source recent-unwrapped --lines 12
   ```

4. Once the prompt appears, tell the operator which workspace and tab to open, what they are approving, and which key confirms it. Leave the prompt for their own keystroke. If World presents a passkey link, give them the link for that step.
5. Watch the pane for success, cancellation, or rollback. After an upgrade, check `ailang-worldd service status --json` with the installed binary, compare the running digest with the candidate, and open the live surface.

If Herdr is unavailable, give the operator the exact command to run in their terminal and verify the outcome afterward.
