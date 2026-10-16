# Row 139 M1 judge r1 (PR #282, 5713f79) -- score 90/100 PASS

Verification: go vet clean; go test ./cmd/... ./host/... all ok (AILANG_BIN v0.41.0); verify_ail.sh PASS (16 identities, 40 + 444 tests).
v0.54.11 verify_ail.sh: runs through World package steps 1-8 (tolerant parser accepts v0.54.11 dry-run), refuses at step 9 "wrong compiler version: AILANG v0.54.11" (behaviour as intended: leg 3 pinned).

1 Digest fidelity: all 8 digests in pins/ailang-releases.json equal origin/dev pins.go (4 rows) and verify_world_package.sh (v0.41.0 binary x2); embedded copy byte-equal (cmp) and tested. No invented digest. Commits (24ee108, c68ded4) were in repo (24ee108 in verify script; c68ded4 not found on origin/dev in a quick grep -- unverified, non-blocking).
2 Behaviour: ToolBinaryRelease derives to "AILANG v0.52.1"; proof/frozen v0.41.0; help text content unchanged (const->var).
3 AC1: exhaustive over git ls-files, Go tokenised (strings only, raw strings counted), shell/yaml/py/json line-scanned, whole-line # comments dropped. Escapes: (a) \b means "ailang_v0.41.0" (underscore-prefixed) is not seen; (b) concatenation/Sprintf-built literals; (c) unprefixed "0.41.0" outside the scanned pattern and ranges like ">=0.30.0" (fences.go:395, packages/*/ailang.toml); (d) file types outside the production set (*.ail, *.toml, *.js, *.ts). All non-blocking.
4 D6: both fixtures real; regex requires ": " after label so "Interface hash (v2):" not matched (fixture contains it; tests pass); negative controls present.
5 Mutations: see table in report.
