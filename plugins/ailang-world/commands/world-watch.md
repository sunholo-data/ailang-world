---
description: Print the commands that watch a running AILANG World episode (episode report --follow and provenance). It only prints.
argument-hint: <episode>
allowed-tools: []
---

Print these two lines for episode `$ARGUMENTS` (ask for an episode if it is empty). Print
only; do not run them.

```sh
ailang-worldd episode report <episode> --follow --addr 127.0.0.1:47330
ailang-worldd --addr 127.0.0.1:47330 provenance --episode <episode>
```

Say that `127.0.0.1:47330` is the rig's daemon; use the address the operator's `serve --bind`
names. `episode report --follow` re-prints the summary as the episode moves; `provenance`
lists its entries, and `/world-why <entry>` explains one.
