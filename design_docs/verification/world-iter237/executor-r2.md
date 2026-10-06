# Row 136 executor round 2

Commit 8b4440a on sprint/row136-mcp-typed-errors (not pushed).

- New binary-free test TestMCPPreDispatchRefusalLinesAreLabelled (host/projection/mcp_conformance_test.go), 3 subtests: tools/list registry read (mcp.go:76), tools/list unprojectable descriptor (:81, a revision-2 registry with input schema type array), tools/call admission (:107, head read fails from its 2nd call, since tools/call lists first). Each asserts exactly one line with prefix 'ailang-worldd: mcp refusal: tools/list -: "' or 'tools/call -: "' plus the cause.
- Design Status set to IMPLEMENTED; design and sprint plan git mv'd planned/ to implemented/. Remaining 'planned/w-mcp-effects' references: design_docs/world-mission.md:960 (charter, left alone) and design_docs/world-mission-log.md:1327 (append-only historical log entry, left as written).

Gates: go vet rc=0; go test -race -count=1 ./host/projection/ ./host/daemon/ rc=0 (557 RUN, 0 FAIL); verify_ail.sh rc=0 (16 identities, 40 + 444 named tests).

Mutants (each applied, new test run alone with AILANG_BIN only, restored; mcp.go clean after each): all 7 KILLED
- :76 label mcp->a2a, :76 method tools/list->"mcp tools": KILLED (tools_list_registry_read)
- :81 label, :81 method: KILLED (tools_list_unprojectable)
- :107 label, method ("mcp invoke"), id ("x"): KILLED (tools_call_admission)

R2 DONE 8b4440a
