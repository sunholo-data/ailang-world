---
description: Print the command that explains one AILANG World log entry or MCP call (ailang-worldd why). It only prints.
argument-hint: <entry|mcp:id>
allowed-tools: []
---

Print this line for `$ARGUMENTS`: a log index, `head`, a `sha256:` entry ref, or an
`mcp:<id>` call id (ask for one if it is empty). Print only; do not run it.

```sh
ailang-worldd --addr 127.0.0.1:47330 why <entry> --json
```

Say that `127.0.0.1:47330` is the rig's daemon; use the address the operator's `serve --bind`
names.
