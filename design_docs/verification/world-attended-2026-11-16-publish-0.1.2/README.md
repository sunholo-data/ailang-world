# Attended publish: `world/core@0.1.2` (2026-11-16)

Mark ran this attended, in their own terminal, using `tools/attended/self_mod_publish.sh` (preflight → approve → dry-run → live) and following `docs/SELF_MOD_PUBLISH.md`. It records D-WORLD-117 = A. The store was a fresh `~/.ailang/world/world-0.1.2.db` (D-WORLD-45 = A). The home service's `world.db` was not touched.

| item | value |
|---|---|
| approval ref (one-shot, spent; the only `approval_claims` row) | `sha256:cd197f324171ee4e992d7948194404a25fa054aa494658e4d1ebf6cf13916cac` |
| publish outcome record (journal seq 6, `effect:attended-publish:0`) | `sha256:329708a3a693836705c6a685b3524ae9d30067ae7c7846ede5af791bd166187c` |
| outcome | `PUBLISHED`; registry `published_at` 2026-11-16T14:52:18Z |
| store journal | 3 intents, 3 outcomes (approve, poll, publish); `reconcile`: 0 indeterminate receipts |

**Stage A, run by the operator before Mark ran the publish** (pinned AILANG v0.41.0):
- `build_world_package.sh` rc=0, with no tree change;
- `verify_world_package.sh` 9/9, with the ready packet byte-equal to the golden;
- `world-publish packet` EQUAL in every field;
- the registry served 404 for 0.1.2 beforehand.

**Independent check after the publish.** `metadata_world_core_0.1.2.served.json.txt` holds the served bytes exactly as fetched from `https://storage.googleapis.com/ailang-registry/packages/world/core/0.1.2/metadata.json` (sha256 `7786d059e2ba4c619fd6ab5cad0ad9eee8f72606b8f332d386b994eae1929f31`).
- **Digests:** `content_hash`, `interface_hash`, `interface_hash_v2`, `tarball_hash` and `tarball_size_bytes` (11030) all equal the golden `scripts/world_package_ready_packet.golden.json`.
- **Validation:** the registry validator (AILANG v0.60.5) reported `contracts_verified` 16/16, 0 skipped.
- **CLI:** `ailang pkg versions world/core` lists 3 versions, with v0.1.2 as latest ("✓ compiles ✓ effects 16/16 contracts"), and `ailang install world/core@0.1.2` downloads 11030 bytes.
- **Control:** 0.1.1 metadata still answers 200.

**Registry defect, filed upstream as sunholo-data/ailang#2373** (the mission-control mirror is #2374). The served 0.1.2 metadata is **not valid JSON**, which is why it is banked with a `.txt` suffix.
- `validation.ailang_version` carries raw newlines (the full multi-line `--version` output).
- `interface_signatures` uses the invalid escape `\,`.

0.1.1's document, validated by v0.51.0, escaped both correctly. The content fields are correct, and the `ailang` CLI reads the document, but strict decoders fail on it. That includes Go's `encoding/json` in World's `host/broker/registry_reconcile.go`. This publish ended PUBLISHED, so no reconciliation was needed. A future indeterminate publish would not reconcile against such a document until the validator is fixed.
