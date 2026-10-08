# ADMIN packaging delta

First PR236 head32ef208, CI37783874633: AILsuccess, Go failure before product tests: tracked-binary hygiene named control1.png and control2.png. This candidate introduced both; baseline has zero Git-binary blobs. The detector in unchanged scripts/verify_go.sh185–219 rejects every such blob, no image/path exception.

Raw instrument PNGs copied byte-identically to persistent local paths in chrome-artifact-manifest.json, then removed from git. Text probe, exact timeout log and hash manifest remain tracked. Images are instrument-only; both Chrome processes timed out. No successful lifecycle/product capture claim.

The same detector conflicts with literal future tracked releasePNG requirement. Fleet-owned ticket inbox_1791466110794_90c512a9/signatureci:world:tracked-binary-rejects-release-png, row173, blocking150. World made no gate/shared-driver/ignore-rule/dependency edit. ADMIN preservation is not a product workaround. D69 remains independently OPEN.
