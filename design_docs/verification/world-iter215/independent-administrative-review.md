# Administrative recovery report — row127 park disposition (iteration 215)

**Disposition:** PASS — park `w-a2a-operator-log-coverage` (row 127) as `needs-human-review` is correct.
**Score:** 85 / 100
**Scope:** Administrative review only. No product implementation exists, was evaluated, or is accepted.

## Basis

Two quorum rounds returned BLOCKED. All three present reviewers rejected with independent,
non-overlapping blocking objections recorded in `evaluator-r2-evidence.json`:

1. **oc-glm-5-3** — card-route constant `GET /.well-known/agent.json` is assertion, not
   verification; load-bearing in the log line, the `fmt.Fprintf` snippet, the byte-exact
   oracle, and the mutation table. No verification-log row reads the route registration.
2. **oc-kimi-k3** — the Success gate is unsatisfiable: it requires full `go test ./...` +
   `./scripts/verify_ail.sh` green while V16 records that exact pristine baseline RED in
   capsule/pkgproj/replay/verifygate, and "an existing red cannot be waved through." No
   terminating comparator exists.
3. **gemini-3-1-pro** — the negative controls assert zero lines for "invalid JSON/params" and
   "unlisted skill", which deterministically emit resolver/registry lines (A2A resolver error,
   registry admission error); these contradict the coverage additions.

One reviewer was absent (claude-sonnet-5@claude-p, quota). No common fix emerged: GLM proposes
a baseline-relative pass rule; Kimi accepts a relative rule only via remediation **or explicit
written quorum acceptance**; Gemini's negative-control fix is orthogonal. No verbatim common fix
selects a new gate policy.

## Controller-measured evidence does not clear the gate

The additional controller measurement
`AILANG_BIN=~/.pinned-ailang/ailang go test ./... -p 1 -count=1 -timeout=300s` → rc0 changes the
execution shape (serial, `-p 1`, `-count=1`, 300s timeout). It is **not** the default
`verify_go` comparator against which V16's pristine reds were recorded, and it is neither a
unique cause nor a quorum acceptance. It therefore resolves none of the three blocking
objections and must not be read as a baseline comparator.

## Blocking findings (unresolved, for the next gate)

- **B1 (card-route constant):** no source/wire verification of the registered card path and
  accepted HTTP method; `GET /.well-known/agent.json` remains unproven.
- **B2 (unsatisfiable success gate):** a terminating pass condition requires either remediation
  of the four pre-existing V16 reds, or written quorum acceptance of a fingerprint comparator
  — neither is presently recorded.
- **B3 (negative controls):** zero-line assertions conflict with resolver/registry admission
  lines; controls must be re-specified.

## Required next gate

A new gate policy must be ratified (human or unattended quorum) before this row may be routed to
sprint. Minimum contents: (1) a verification row reading the card-route registration; (2) a
satisfiable Success criterion — remediation row for the four baseline reds **or** explicit
written quorum acceptance of a baseline-relative fingerprint comparator with the exact
comparator command attached; (3) corrected negative controls relocating invalid-JSON and
unlisted-skill cases to positive 1-line expectations. The controller's serial rc0 run is not a
substitute for any of these.

Product implementation was **not** evaluated and is **not** accepted by this report.