# Iteration222 designer evidence

Author: fresh native GPT6.1Sol fallback, FLAGGED per controller D48 directive; preferred GLM5.3 Unknown model. Scope: D44=A ONE protocol revision after row108 (LANDED221); no code, harness, walk, baseline, credit, commits/push or external messages.

Readiness: revision complete for fresh full author-excluding quorum; execution BLOCKED. No existing method-enforcement facility/configuration/control has been measured; authority D2–D5 remains unfrozen. The revision corrects World-only verified wall timing, removes the bundle alternative and excludes diagnosis content regardless of ancestry, requires attributable pre-reveal method/bypass controls, and replaces shorthand logs.

Authoritative reads: design-doc-creator SKILL.md and design_doc_structure.md; mission-control gate-2-pick.md/gate-3-route.md; World CLAUDE.md/coding-standards.md, charter clause5/D44/row108/row114, DESIGN thesis and boundaries. No scaffold helper invoked. Inherited doc and quorum artifacts read directly from named git objects.

## Exact measured premises

### V1

```sh
git rev-parse HEAD origin/dev
```

Exit: 0

```text
03ea641884327045ac59754d9fc25ed046066d7d
03ea641884327045ac59754d9fc25ed046066d7d
```

### V2

```sh
sed -n '715,745p' host/daemon/daemon.go
```

Exit: 0

```text
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", d.handleHealth)
	mux.HandleFunc("GET /v1/head", d.handleHead)
	mux.HandleFunc("GET /v1/worlds/{ref}", d.handleWorld)
	mux.HandleFunc("GET /v1/objects/{ref}", d.handleObject)
	mux.HandleFunc("GET /v1/objects/by-semantic-id/{name...}", d.handleObjectsBySemanticID)
	mux.HandleFunc("GET /v1/log/{index}", d.handleLogEntry)
	mux.HandleFunc("GET /v1/log", d.handleLogRange)
	mux.HandleFunc("GET /v1/registry/{name...}", d.handleRegistry)
	mux.HandleFunc("POST /v1/commit", d.handleCommit)
	mux.HandleFunc("GET /v1/receipts/{id}", d.handleReceipt)
	mux.HandleFunc("GET /workbench", d.handleWorkbench)
	// The two A2A projection routes (w-a2a-session-projection P6.B-A2A-CARD)
	// are ADDITIVE: the frozen /v1/ table above is untouched, and the routes
	// are NOT in isProtected — the projection handler resolves the session
	// itself so /a2a/ can answer in JSON-RPC form (B5).
	mux.HandleFunc("GET /.well-known/agent.json", d.projection.AgentCard)
	mux.HandleFunc("POST /a2a/", d.projection.A2A)
	mux.HandleFunc("POST /mcp/", d.projection.MCP)
	return NewSessionMiddleware(d.resolver, d.credentialBudget, d.writeInternalError).Wrap(d.isProtected, mux)
}

// isProtected reports whether a request must carry a valid session credential
// (w-session-authority D6). This sprint it is ONLY POST /v1/commit — the
// daemon's sole mutation. The eight GET routes pass through unauthenticated as
// declared residual R1 (full /v1/* read enforcement is a follow-up queue row; a
// config flip here later is a rewrite-free change).
func (d *Daemon) isProtected(r *http.Request) bool {
	return r.Method == http.MethodPost && r.URL.Path == "/v1/commit"
}

```

### V3

```sh
sed -n '41,79p' host/daemon/middleware.go
```

Exit: 0

```text
func NewSessionMiddleware(resolver authority.Resolver, budget time.Duration, fail func(http.ResponseWriter, *http.Request, error)) *SessionMiddleware {
	return &SessionMiddleware{resolver: resolver, budget: budget, fail: fail}
}

// Wrap returns a handler that enforces the session boundary on requests for
// which protected(r) is true; all other requests pass through to mux untouched
// (residual R1 — the read routes pass unauthenticated this sprint).
func (m *SessionMiddleware) Wrap(protected func(*http.Request) bool, mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !protected(r) {
			mux.ServeHTTP(w, r)
			return
		}
		header := r.Header.Get("Authorization")
		lctx, cancel := context.WithTimeout(r.Context(), m.budget)
		out, err := m.resolver.ResolveContext(lctx, header, time.Now().Unix())
		timeout := err != nil && timedOut(lctx, err)
		cancel()
		if err != nil {
			// A lookup that did not finish is not a denial: never 401.
			if timeout {
				writeAPIError(w, "Timeout", fmt.Sprintf("credential lookup deadline (%s) exceeded", m.budget), http.StatusServiceUnavailable)
				return
			}
			m.fail(w, r, err)
			return
		}
		if out.Denied != nil {
			writeSessionDenial(w, *out.Denied)
			return
		}
		ctx := authority.WithBinding(r.Context(), out.Success)
		mux.ServeHTTP(w, r.WithContext(ctx))
	})
}

// writeSessionDenial is THE ONE session-denial -> HTTP mapping: Wrap uses it
// for /v1/commit, and the daemon injects it into host/projection's Config.Deny
// so a card-route denial is byte-identical to the middleware's (B1) — the
```

### V4

```sh
rg -l 'ailang-worldd|/v1/commit|worldd commit' scripts tools docs
```

Exit: 0

```text
docs/QUICKSTART.md
scripts/verify_go.sh
```

### V6a

```sh
rg -c UNMEASURED design_docs/verification/w-1-0-value-demonstration.md
```

Exit: 0

```text
10
```

### V6b

```sh
sed -n '143,158p' design_docs/verification/w-1-0-value-demonstration.md
```

Exit: 0

```text
|---|---|---|---|---|---|
| Q1 `iter171-index-row` | `world-mission.md` row-86 bullet and the v0.30.0 `mission`-command bullet (~lines 82–101); log `## 172`: grep the index for the row, control-grep a neighbour, run the prescribed regenerator, positive-control another command group | UNMEASURED (record states method, not cost) | 20 | surviving `ailang.v0.30.0` binary and `git log -S` over `world-mission-index.md` | PASS |
| Q2 `iter154-unfenced-pi` | `world-mission.md` row 78; log `## 154` containment note: read the recipe's two-flag invocation block, read the runner's actual invocation, compare | UNMEASURED (record states method, not cost) | 13 | live fleet `scripts/mission_pi_run.sh` via `gh api` and `gh issue view 1043` | PASS |
| Q3 `iter155-faithfulness-proof` | `world-mission.md` row 81; log `## 155`: trust three green `git commit`s, get a green `shasum -c`, catch the problem only by reading `git show --stat` per commit, rebuild, hit the zsh word-splitting trap on retry | UNMEASURED (record states method, not cost) | 18 | in-git charter row 81 and log `## 155` text, and `git show --stat` of the rebuilt row-59 split | PASS |
| Q4 `iter181-ci-bench-401` | `world-mission.md` STATUS 2026-09-24 (iter-181): local gates all green, remote CI red, read the CI failure, reproduce, judge round 2 | UNMEASURED (record states method, not cost) | 10 | git `a036062` diff of `host/daemon/bench_test.go` and the in-git STATUS text | PASS |

**Reading this table honestly.** All four taken questions pass: each was verified, and each took
at most 20 s against the 300 s bar. That meets clause 5's "≥3". There is one caveat, recorded
rather than hidden. M2 wrote the incident objects from the same mission record that later
served as verification: the answers were diagnosed first and then committed. So the walk times
measure **retrieval of a recorded diagnosis through World's provenance chain**, not a fresh
diagnosis. The baselines are the original diagnoses, and they are UNMEASURED, so no speed-up
ratio is claimed. The walk found each incident by a linear scan (F-1), which is cheap at 4
entries and says nothing about mission scale.

## Findings
```

### V7a

```sh
sed -n '79,102p' host/store/object_references.go
```

Exit: 0

```text
func (s *Store) ObjectReferences(ctx context.Context, ref hashref.HashRef, after *ObjectReferenceCursor, limit int) ([]ObjectReference, error) {
	if err := s.checkQuarantine(); err != nil {
		return nil, err
	}
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if err := validateObjectReferences(ref, after, limit); err != nil {
		return nil, err
	}
	if !s.referenceIndexesAvailable {
		return nil, &ReferenceIndexUnavailableError{}
	}
	items := make([]ObjectReference, 0, limit)
	queries := [...]string{transitionReferencesSQL, functionReferencesSQL, interpreterReferencesSQL, worldReferencesSQL}
	for kind := ReferenceTransitionRef; kind <= ReferenceStateRoot && len(items) < limit; kind++ {
		if after != nil && kind < after.Kind {
			continue
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		query := queries[kind]
		var afterKey any
```

### V7b

```sh
sed -n '26,48p' host/store/object_commits.go
```

Exit: 0

```text
func (s *Store) ObjectCommits(ctx context.Context, ref hashref.HashRef, afterEntry int64, limit int) ([]int64, error) {
	if err := s.checkQuarantine(); err != nil {
		return nil, err
	}
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	const op = "ObjectCommits"
	if limit < 1 || limit > MaxObjectCommitPage {
		return nil, &InvalidLimitError{Op: op, Limit: limit, Max: MaxObjectCommitPage}
	}
	if _, err := hashref.Parse(ref.String()); err != nil {
		return nil, &InvalidRefError{Op: op, Field: "ref", Text: ref.String(), Err: err}
	}
	if afterEntry < -1 {
		return nil, &InvalidObjectCommitCursorError{AfterEntry: afterEntry}
	}
	rows, err := s.db.QueryContext(ctx, objectCommitsSQL, ref.String(), afterEntry, limit)
	if err != nil {
		return nil, fmt.Errorf("store: object commits query: %w", err)
	}
	entries := make([]int64, 0, limit)
	for rows.Next() {
```

### V9a

```sh
sed -n '29,53p' host/projection/mcp.go
```

Exit: 0

```text
func (h *Handler) MCP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.invokeWait)
	defer cancel()
	h.mcp.ServeHTTP(w, r.WithContext(ctx))
}

type mcpAdapter struct{ h *Handler }

func (a mcpAdapter) ResolveSession(ctx context.Context, r *http.Request) (protocol.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, a.h.credentialWait)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out, err := a.h.resolver.ResolveContext(ctx, r.Header.Get("Authorization"), time.Now().Unix())
	if err != nil {
		return nil, err
	}
	if out.Denied != nil {
		message := msgUnknown
		switch *out.Denied {
		case authority.DenialAbsent:
			message = msgAbsent
		case authority.DenialMalformed:
			message = msgMalformed
```

### V9b

```sh
sed -n '120,134p' host/projection/mcp.go
```

Exit: 0

```text
	if err := ctx.Err(); err != nil {
		return protocol.InvocationResult{}, err
	}
	if a.h.coord == nil {
		return protocol.InvocationResult{}, errors.New("projection: invocation coordinator is unavailable")
	}
	task, err := a.h.mintTask()
	if err != nil {
		return protocol.InvocationResult{}, err
	}
	result, err := a.h.coord.Dispatch(ctx, coordinator.Call{Request: request, EpisodeID: binding.EpisodeID, Grants: binding.Caps, SkillID: id, TaskID: task, Input: input})
	if err != nil {
		a.h.logRefusal("mcp invoke", coordinator.InvocationID(binding.EpisodeID, task), err)
		return protocol.InvocationResult{}, err
	}
```

### V9c

```sh
sed -n '355,367p' host/projection/projection.go
```

Exit: 0

```text
				}
				parsed, err := hashref.Parse(s)
				if err != nil {
					protocol.A2AError(w, req.ID, codeInvalidParams, msgInvalidParams)
					return
				}
				pin = &parsed
			}
			result, err := h.coord.Dispatch(ctx, coordinator.Call{Request: admitted, EpisodeID: out.Success.EpisodeID,
				Grants: out.Success.Caps, SkillID: d.ID, TaskID: params.ID, Input: params.Message.Parts[0].Data, PinnedFn: pin})
			if err != nil {
				// A timed-out invocation must release its transport as well as its
				// capsule. The client can then retry with the same task ID after
```

### V10

```sh
rg -n 'GET.only|GET-only|method.filter|method filter|read.only.*proxy|read-only.*proxy' host cmd scripts docs
```

Exit: 0

```text
host/broker/registry_reconcile_test.go:942:			"a TARGET-only transport failure", got)
host/daemon/daemon_test.go:1203:			t.Fatalf("%s /.well-known/agent.json = (%d, %s), want 405 — the card route is GET-only", method, rec.Code, rec.Body)
```

### V11a

```sh
git merge-base --is-ancestor 29e1336 HEAD
```

Exit: 1

```text
```

### V11b

```sh
git merge-base --is-ancestor 29e1336 b702dad27cdb96c785e189f6afcc43b419985c41
```

Exit: 1

```text
```

### V11c

```sh
git cat-file -t 29e1336
```

Exit: 0

```text
commit
```

## Inherited criticism — verbatim

### r1 gemini-3-1-pro — strongest_objection

Section 7 schema requires `baseline-transcript` and `walk-transcript` to contain the `verify_last_cmd` block and total duration, and Section 5 sets the timer STOP at 'the LAST verification command'. However, D2 and P4-P6 strictly isolate V into a separate, fresh role that executes after B and W declare their answers. Isolated agents B and W cannot log V's subsequent commands in their own transcripts, and enforcing a continuous `first->last timestamp` cross-check across isolated role sessions is impossible without inadvertently measuring the inter-agent idle time.

### r1 gemini-3-1-pro — catch

The timer stop rule and transcript schema were inherited from row 92's single-executor design, which contradicts row 114's introduction of three strictly isolated B/W/V roles.

### r1 gemini-3-1-pro — proposed_fix

In §5 Timer rules, redefine STOP for B and W as 'the moment the answer is declared'. In §7, remove `verify_last_cmd` from `baseline-transcript` and `walk-transcript`. Add a separate `verification_duration_s` block to V's `verification` file, and update D5 to explicitly calculate the <=300s total duration as the sum of the isolated agent's duration plus V's duration for that arm.

### r1 oc-glm-5-3 — strongest_objection

E4/D4 assert retrospective isolation 'by construction', but the specified mechanism — 'a worktree at the... pre-diagnosis pin' (§5 P3) — cannot deliver the claimed property that 'the packet environments exclude everything after it' (E4): a git worktree shares its parent repository's object database and refs, so B, whose packet explicitly grants git (§P3), can read post-pin commits (git show <post-pin-sha>:design_docs/world-mission.md, or the fix commit 29e1336's context) and retrieve the controller diagnosis already recorded at world-mission.md:1891 (P4). Since the only two registered candidates (C1/C2, §8) are retrospective, both executable baselines collapse into retrieval of a pre-written answer — the exact row-92 seeding flaw this doc exists to kill — and the P7 transcript audit cannot distinguish contaminated git archaeology from honest archaeology, so any D5/AC9 speedup ratio against these baselines is silently fabricated.

### r1 oc-glm-5-3 — catch

Two more load-bearing premises are asserted rather than verified: (1) §6.1/P2 require the frozen store copy's sha256 to be unchanged post-walk, but no verification-log row establishes that serve + the E3 existence probes are byte-inert on the sqlite file — row 92's own F-5 stop→mint→restart dance (P10) suggests the server process writes to the db, and if boot/WAL checkpointing mutates the copy then every attempt, including the anticipated UNANSWERABLE(capture) gap findings, self-voids as INVALID(pin). (2) The diagnosis embargo (Design Freeze, P0/MU-9) is an unbounded wait: if the protocol stalls before P6 (e.g. NOT-RUN(TRANSPORT)), a live incident's diagnosis becomes unrecordable indefinitely, with no release-on-abandon path anywhere in §13.

### r1 oc-glm-5-3 — proposed_fix

Replace the B-environment mechanism in §P3/E4 with a history-truncated environment: `git bundle create b.bundle <pin>` (a bundle carries only history reachable from the pin), then `git clone b.bundle b-env && git checkout <pin>`; record b.bundle sha256 in the §7 pins schema, and add preflight leg E4b with a quoted mechanical check: in B's environment `git rev-list --all --count` equals the count at the pin AND `git cat-file -e <sha-of-diagnosis-bearing-commit>` fails (that sha listed in pins). Add verification-log row V10 at the base pin: sha256 the store copy → serve on a scratch port → run the E3 existence probes → stop → re-sha; if bytes changed, redefine §6.1 to audit against the pristine pre-serve copy plus an explicit served-file delta manifest (or open sqlite read-only) so daemon-caused mutation can never be read as INVALID(seeding)/INVALID(pin). Bound the embargo: it lifts on any NOT-RUN(TRANSPORT) or abandoned attempt, at which point the controller writes the diagnosis to the harvest log and the attempt is preserved as NOT-RUN, never re-armed.

### r1 controller — verbatim

No authorization for designer-labelled human decisions or ratified embargo; absent corpus differs from unanswerable complete corpus; completion must require measured baseline and unseeded integrity; no-credential alone is not enforced role confinement; Verification Log lacks controls for absence claims.

### r2 gemini-3-1-pro — strongest_objection

The Verification Log (§16) contains typographical ellipses within the command literals for V6 and V7, proving these shell commands were never executed as written. This violates the premise verification gate, as the logged commands are synthetically generated shorthand rather than copy-pasted instructions from a live terminal session.

### r2 gemini-3-1-pro — catch

§16 V6 uses 'grep -c UNMEASURED …/w-1-0-value-demonstration.md' and V7 uses '…ObjectCommits''. Both contain the literal '…' character.

### r2 gemini-3-1-pro — proposed_fix

Replace the shorthand command literals in V6 and V7 with their exact, fully-qualified executable paths (e.g., 'grep -c UNMEASURED design_docs/verification/w-1-0-value-demonstration.md') and update the result columns with their true, verbatim shell outputs from the live worktree.

### r2 oc-glm-5-3 — strongest_objection

The walk arm's confinement is not enforced for its primary seeding vector, contradicting the doc's own binding definition of CONFINED. §P3-W allows W 'network except loopback scratch port' with no method restriction — the one sanctioned destination carries the full HTTP surface including POST /v1/commit, /v1/receipts, /workbench, /a2a/ (per V2/V3). The doc itself disclaims the only server-side barrier twice: P9 declares isProtected 'a routing fact, NOT a confinement mechanism', and §P3-W demotes the absent session credential to 'belt-and-suspenders, not the mechanism'. The PR deny-controls (§5 PR) test only 'network except loopback scratch port + any write outside scratch' — there is no deny-leg for POST/PUT to the scratch port. So W can physically POST to the served copy through its sanctioned channel, and prevention rests entirely on an unverified assumption that the scratch serve rejects unauthenticated mutation — a claim absent from §16 (V3 verifies only which routes isProtected covers, not what it enforces or whether the scratch-serve launch flags retain it). What remains is post-hoc detection (§6.5, delta audit, MU-7), which voids the attempt after the scarce real Gate-2 incident has been burned — the exact 'capability-absence prose is not isolation' failure §2 forbids, now on the walk arm's own channel.

### r2 oc-glm-5-3 — catch

Two things to verify or fix before re-submission: (1) what isProtected actually enforces (token comparison? env flag? file?) and whether the scratch serve launched with default/pinned flags preserves it — currently unverified in §16; (2) the §P3 bundle alternative is internally unsatisfiable for C1/C2: a bundle 'truncated at the pin' necessarily contains the diagnosis commits (e.g. 29e1336) as ancestors, so 'rev-list --all --count equals pin-time count' and 'cat-file -e <diagnosis-sha> fails' cannot both hold — the permitted alternative silently degenerates for exactly the candidates that carry diagnosis_sha[] in pins.

### r2 oc-glm-5-3 — proposed_fix

Amend §P3-W and §5 PR as follows: (a) W's confinement envelope becomes 'network except loopback scratch port, GET methods only' — enforced at the confinement layer (e.g., method-filtered loopback proxy), not at the server; (b) add a mandatory PR deny-leg for W, run before question reveal: issue a POST to the scratch serve's mutation routes (minimum: POST /v1/commit) from inside W's environment and quote the transport result (expected HTTP 4xx ⇒ rc≠0); failure ⇒ NOT-RUN(ISOLATION); (c) add a §16 verification row quoting the isProtected enforcement code (exact check and lines) plus the scratch-serve launch flags, pinned at P2; (d) reword D3 to: 'prevention for filesystem/process/egress plus method-level denial on the allowed channel; §6.5 post-hoc detection is a backstop only'; (e) restrict the bundle alternative to candidates whose pinned diagnosis_sha[] contains no ancestor of the pin, and state explicitly that C1/C2 (diagnosis commits already ancestors of 9b2cfac) are archive-only pending decontamination grep — eliminating the contradictory proof path.

### r2 controller — verbatim

Reject pending attended decisions D2-D5 and timer correctness: clause5 bounds the World walk to verified answer, but r2 binds 300s to baseline question issue through BOTH arms and orchestration, imposing an extra baseline-time bar. Need separate World issued-to-verified envelope including its own idle; whole experiment elapsed is reported without changing clause5. Preserve NOT-RUN(NO-CORPUS), confinement and unseeded controls. C3 author status still running while baseline now completed82.938396125s; treat as historic snapshot and do not claim paired result. No third author revision this iteration.

## Final document checks

Planner-Lane is `codex-ok`, justified only by protocol/documentation scope; freeze gates still block planning. Designer/controller diagnosis exposure excludes them from blinded B/W roles. `git diff --check` returned exit0 (untracked artifacts require direct whitespace check); direct changed-file scan found zero trailing whitespace. No code tests run for this protocol-only change.
