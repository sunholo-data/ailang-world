# Evaluator report r1 — iteration 382, row 201 A2 sprint 4 (M7), head e07b7bad

## Verdict: PASS — 91/100, no BLOCKING findings

| Category | Score |
|---|---|
| Correctness vs design and plan §9 | 28/30 |
| Tests: non-vacuous, kill mutants, cover operands | 27/30 |
| Security / fail-closed | 19/20 |
| Coding standards | 9/10 |
| Scope and size | 8/10 (see note on the 1,054 figure) |

## Review
- Routes: `GET /surface/passkey` and `/surface/passkey.js` (host/daemon/daemon.go, surface_page.go) call `setSurfaceHeaders` (surfaceCSP unchanged: `git diff` shows no change to the const). They serve the embedded bytes, with no-store, nosniff, no-referrer and the correct content types. Non-GET gets 405 and the sub-paths get 404 (TestSurfaceAppRoutes extended). Frozen v1 table test still passes.
- Page: no fetch, XHR, form, innerHTML or eval. Text is set with textContent. The fragment is parsed strictly (base64url regex, JSON, v==1, op in {create,get}). origin and rp are compared to location.origin and location.hostname and the page returns before showing any button. The `signal` is passed beside `publicKey` on create and get. The timer is 240000 ms, a late result is ignored (`expired`), and the button is disabled while pending. After expiry nothing is saved (the save link is only created after the `expired` and `failure` guards). The enrolment code is built per §2.8 and XB1 compares it to Go's code. The blob: save works under the CSP (the download is observed in XB1).
- XB1 uses the product path: the link is parsed from the CLI's stdout, `runApprover enroll` and `runApprove --webauthn-link` and `--webauthn-assertion` run against temp stores, and the sign count and consumed row are checked. It does not re-encode.
- Scope: `git diff --name-status` matches plan §5 (20 files). Nothing under tools/launchd, world/, packages/, scripts/ or cmd/world-publish. No Go file in host/approval changed. host/approval has no non-stdlib imports (the grep for dotted imports printed nothing). No new production exec site or route group; the Node launch is test-only.
- QUICKSTART §16 is tagged "attended — pending first verbatim run", is consistent with the CLI flags (`--webauthn-link`, `--webauthn-assertion`, `--db`/`--credentials-db`) and with the CLI's next-step output, uses the scratch store and localhost only, and is documentation only (not run).
- Plan §1.2 C22 narrowing: sound. The page bytes reach tests only through the daemon handler (PG1) and the passkey_page_* tests, and my page mutants were all killed by the passkey tests alone.

## Numstat (`git diff --numstat 27ce163e HEAD -- cmd host`)
Total 1,054 (insertions + deletions), under 2,000. Note: the stated unit also counts test and driver code (848 of it). Production is about 206.

## Test runs (own, in the judge worktree, AILANG_BIN pinned)
- XB1, XB2 and PG1 with WORLD_REQUIRE_PASSKEY_BROWSER=1, headless Chrome 157, run twice. Run 1: rc 0, XB2 4.14 s, XB1 6.43 s. Run 2: rc 0, XB2 4.03 s, XB1 3.32 s. Logs: judge/xb-run1.txt and xb-run2.txt.
- `go test -count=1 ./host/surfacepage/... ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/ -timeout 20m`: rc 0 (cmd 186 s, daemon 358 s).
- `go test -race -count=1 ./host/surfacepage/`: rc 0.
- `go vet ./...`: rc 0.
- `go test -run '^$' ./...`: rc 0.
- I did not run verify_go.sh or a whole `go test ./...`, as instructed.

## Mutants (own ledger: /Users/voightkampff/.ailang/state/world-iter382/judge/mutants.tsv, 46 lines; per-run outputs are mut-<id>.txt)
Every run was checked for a non-empty `git diff --stat` and a vet/compile pass before the test result was read. The worktree is clean afterwards (`git status --short` empty).
Test selectors:
- browser = XB1 and XB2
- pg1 = PG1 + TestSurfaceAppRoutes + frozen-table
- wn2a = TestPasskeyEnrollFlagRules + TestApproveByPasskey
- gen = TestPasskeyFixtureGeneratorRefuses

**Controls**
- CTL-LIVE-comment (comment text edit): SURVIVED, as expected (must-live)
- CTL-DIE-timeout (240000 to 239000): KILLED (must-die)

**KILLED (34 of 38 non-control rows)**
- Page checks: J-ver, J-op-create, J-op-get, J-origin, J-rp (each operand separate), J-rp-origin-swap (hostname to host), J-b64-strict, J-json-catch, J-refuse-return1, J-refuse-return2.
- Ceremony guards: J-latecheck, J-expired-set, J-abort, J-signal-get, J-btn-disable, J-clear, J-nosave-after-fail.
- Code and UV: J-code-label, J-code-order, J-uv-get.
- Routes and headers: R-page-route, R-script-route, R-method-page, R-method-script, R-hdr-page, R-ctype, R-subpath.
- WN2a: W-wn2a, W-wn2a-order.
- HTML: H-html-script, H-save-hidden.
- Non-vacuity and carry-forward rows: NV-routes-both-XB, NV-routes-both-PG1, NV-wn2a-CW1, N-nextstep-order, G-argc, G-notdir, G-notempty, G-chrome.

**SURVIVED**
- J-allow-empty (allowCredentials removed).
- J-userhandle (userHandle forced to null).
- J-blob-type (Blob type changed to text/html).
- R-csp (blob: added to the style-src of the surfaceCSP value).
- H-go-enabled (initial `disabled` attribute on the button removed).

The plan names the first two as equivalents. See the findings below for how much weight each deserves.

## Non-vacuity per milestone
- M7e: removing both route registrations fails PG1 and also XB1/XB2 (NV-routes-both-PG1 and NV-routes-both-XB, KILLED).
- M7b: 24 page and HTML mutants are killed by XB1/XB2, so they depend on the page bytes.
- M7a: removing the WN2a check fails the WN2a arm of `TestPasskeyEnrollFlagRules` (NV-wn2a-CW1, KILLED). Swapping the next-step args fails the pinned next-step line (N-nextstep-order, KILLED).
- M7f: the generate.sh refusal mutants (G-argc, G-notdir, G-notempty, G-chrome) each fail TestPasskeyFixtureGeneratorRefuses.

## BLOCKING findings
None.

## Non-blocking findings
1. J-allow-empty survives: no test proves the get ceremony passes the CLI's allow list. The authenticator would find the single resident credential anyway, and the CLI judges the credential ID. This is a plan-named equivalent, but the allow list is a real behaviour (it narrows which credential can sign). A second enrolled credential in XB1 would kill it.
2. R-csp survives: PG1 compares the header to the `surfaceCSP` identifier, not a literal, so a weakened CSP value passes PG1. Reproduce with `sed -i '' "s/style-src 'unsafe-inline'/style-src 'unsafe-inline' blob:/" host/daemon/surface_page.go && go test ./host/daemon -run TestPasskeyPageRoutes`. The CSP is unchanged in this sprint and I did not run a test that pins its value. A literal pin of the CSP on the passkey routes would close this.
3. J-blob-type survives, so the saved file's MIME type is unasserted. This is harmless because the file is saved as bytes.
4. H-go-enabled survives. The initial `disabled` attribute is redundant with the script. This is cosmetic.
5. 1,054 counts test and driver lines. Production is about 206. This is within the limit.
