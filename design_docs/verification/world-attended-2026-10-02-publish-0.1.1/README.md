# Attended publish — `world/core@0.1.1` (2026-10-02)

Run attended by Mark in his own terminal via `tools/attended/self_mod_publish.sh`
(preflight → approve → dry-run → live), per `docs/SELF_MOD_PUBLISH.md`. Store: fresh
`~/.ailang/world/world-0.1.1.db` (`D-WORLD-45` = A); the legacy `world.db` was not touched.

| item | value |
|---|---|
| approval ref (one-shot, spent) | `sha256:deb19c62a01447a16b8684fbd32989610a39e57469330fa226e861a6001a06ff` |
| publish record | `sha256:60b46d4bbf8e28c648a0dad34c0e93c88b286d48c5ae88155962175c36613b24` |
| outcome | `PUBLISHED` (exit 0), registry `published_at` 2026-10-02T07:58:43Z |

Stage A, same session, pinned AILANG v0.41.0 (`24ee108`): `build_world_package.sh` rc=0 with
no tree change; `verify_world_package.sh` 9/9 and ready packet byte-equal to the golden;
`verify_ail.sh` rc=0; `go vet ./...` rc=0; `go test ./...` rc=0 (24 packages);
`world-publish packet` EQUAL in every field.

Independent post-publish check: `metadata_world_core_0.1.1.served.json` (fetched from
`https://storage.googleapis.com/ailang-registry/packages/world/core/0.1.1/metadata.json`)
equals the golden on content_hash, interface_hash, interface_hash_v2, tarball_hash and
tarball_size_bytes 10610; registry validator (AILANG v0.51.0) `contracts_verified` 16/16,
0 skipped. Controls in the same pass: 0.1.0 metadata 200, absent 0.1.9 metadata 404.

Operator note: the first dry-run/live refused (nothing sent) because a pasted placeholder
`WORLD_APPROVAL` was exported; re-exported the saved ref and the rehearsal passed.
