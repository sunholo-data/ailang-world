# ailang-world — the AILANG World plugin for Claude Code

Two skills and three commands for working with AILANG World:

- **Skill `ailang-world`** (the agent side): loads inside a World episode, when the agent's
  tools are World's `mcp__world__*` tools. How to work so the record is complete and gradeable.
- **Skill `ailang-world-operator`** (the human side): starting, watching, grading and finishing
  an episode on a repo, and publishing.
- **`/world-start <episode>`**, **`/world-watch <episode>`**, **`/world-why <entry>`**: they
  print `ailang-worldd` lines and never run them.

There is **no MCP server** in this plugin (D-WORLD-85, D-189-4 = A). An agent gets World's tools
only from the launch line `ailang-worldd episode start` prints, and minting that session is a
human act at a terminal. The plugin holds no token and no address book: the session token
lives only in `~/.ailang/world/sessions/<episode>.session` (0600) and the per-launch MCP config.

## Install

```sh
claude plugin marketplace add sunholo-data/ailang-world
claude plugin install ailang-world@ailang-world-marketplace
```

`marketplace add` also takes a git URL or the path of a local checkout of this repo.

## Kept honest

Every `ailang-worldd` and `world-publish` line in these files is checked against the CLI's own
help in CI (`TestPluginCommandsMatchTheCLIHelp` in `cmd/ailang-worldd` and
`cmd/world-publish`): an unknown verb or flag fails the build.
