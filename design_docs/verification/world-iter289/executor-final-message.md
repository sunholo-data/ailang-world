Implemented M5e in the single allowed file: 123 added lines, covering all 14 arms. No Git writes performed. Scratch directory deleted.

`git status --porcelain`:

```text
 M host/daemon/workspace_protected_test.go
```

`git diff --stat`:

```text
 host/daemon/workspace_protected_test.go | 123 ++++++++++++++++++++++++++++++++
 1 file changed, 123 insertions(+)
```

All results below are **sandbox observations only**; controller reruns remain required.

| Gate | rc | Observed verdict |
|---|---:|---|
| G1 real pin, final restored file | 0 | `--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool (14.57s)`; 14 subtests PASS, none SKIP |
| G2 pin unset | 0 | `--- SKIP: TestWorkspaceWriteProtectedAliasesPolicyTool (0.00s)` |
| G3 vet | 0 | PASS; empty output |
| G4 gofmt | 0 | PASS; empty output |
| G5 existing lexical test | 0 | `--- PASS: TestWorkspaceWriteProtectedPolicyTool (1.80s)` |

Final G1 subtest lines:

```text
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/symlink-dir-new-write (1.04s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/symlink-dir-write (1.03s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/symlink-dir-edit (1.04s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/symlink-file-write (1.02s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/symlink-file-edit (1.04s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/hardlink-write (1.02s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/hardlink-edit (1.03s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/no-protected-symlink-dir-new-write (1.02s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/no-protected-symlink-dir-write (1.04s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/no-protected-symlink-dir-edit (1.04s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/no-protected-symlink-file-write (1.04s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/no-protected-symlink-file-edit (1.03s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/no-protected-hardlink (1.04s)
--- PASS: TestWorkspaceWriteProtectedAliasesPolicyTool/no-protected-singlelink (1.15s)
```

For the drill table, **S** means `symlink-dir-new-write`, `symlink-dir-write`, `symlink-dir-edit`, `symlink-file-write`, `symlink-file-edit`; **H** means `hardlink-write`, `hardlink-edit`.

| Mutant | Applied; diff shown? | Build rc | Killer rc; failed subtests | Failing assertion text |
|---|---|---:|---|---|
| D1 MUT-NO-PROTECTED-IN-POLICY | Yes | 0 | 1; S | `refusal mechanism missing "resolves to"` and `"fs_deny_write"` |
| D2 omit `lnk` | Yes | 0 | 1; `symlink-dir-new-write` | Same two mechanism assertions |
| D2 omit `alias.txt` | Yes | 0 | 1; `symlink-file-edit` | Same two mechanism assertions, despite generic `ok:false` |
| D2 omit `hl.txt` | Yes | 0 | 1; `hardlink-edit` | `refusal mechanism missing "more than one hard link"` |
| D3 MUT-TOOL-REGRESSION | Yes | 0 | 1; S, H, `no-protected-hardlink` | Missing symlink/hard-link mechanism text |

D1 and D3 additionally reported `host target bytes = "changed\n", want "seed\n"` and, for the new target, `new protected target must remain absent: <nil>`. D3’s five unprotected symlink controls and single-link control passed.

Every mutation was restored using `cp`, with matching SHA-256 and green restoration checks. Final hashes:

```text
test:         32af0c658b8caf38ed994fb97cb1b52246085e0282b5adec8a9828f5c0b802c2
workspace.go: 23099ff0669766d360d26ec2ca6034a1b6e29f86033bce05ed059750e5276baa
```

Finding: the plan’s isolated D3 selector lacks a subtest-leading `^`, so it also selects `no-protected-hardlink`. The literal selector returned rc 1 with eight failures; an additional anchored selector returned rc 1 with exactly seven protected failures and six passing controls. No implementation deviation; race tests were not run.

`UNINFORMATIVE UNDER SANDBOX`: the initial cleanup command was rejected with `rm -f style commands are not permitted. Use a safer approach`. Path-validated cleanup subsequently succeeded; `.tmp-iter289` is absent.