# Host changelog

## Unreleased — candidate; acceptance pending

The host adds session-filtered MCP tools at `POST /mcp/`, alongside the agent card
and A2A surface. Each invocation item re-admits against the current registry and
receives a fresh task ID. Older-version batches run sequentially with independent
commits: a later failure can hide earlier committed results, and retrying can repeat
an effect. Encoded MCP names are reversible; A2A IDs remain verbatim.

Credential resolution, registry reads, aggregate requests and transport writes use
the frozen 3/10/20/30-second budgets. The shared callback runner uses 20 seconds and
eight slots; the distinct subprocess reservation limit is also eight. MCP wire and
SSE handling come from upstream module v0.47.2; the archived compiler remains
v0.41.0. Slow bodies and noncooperative callbacks retain the documented residuals.

The resend socket fixture uses deterministic cooperative execution to test durable
commit followed by a lost response and same-task reconciliation. The separate real
capsule execution/replay test remains unchanged. Normal full profiles, independent
product evaluation and required CI remain pending; this entry does not assert a
release or product acceptance.
