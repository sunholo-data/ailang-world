Ground truth (E2) for one question about trial episode w185 (ailang-world, 2026-10-16), produced by the controller from git and the archival store after the episode; neither answering arm produced or saw it.
- row-185-text.txt: queue row 185 as the charter states it (the task the agent was given): degraded "recorded on the session, shown by /v1/health / session list".
- session_new-comment.diff: the first part of the committed diff of cmd/ailang-worldd/session_new.go, including the header comment the change added.
- store-schema.txt: the session_credentials table in the episode's store (no degraded column), and host/store/store.go at the base commit: the store's schema version constant and the errors a store of another version raises.
