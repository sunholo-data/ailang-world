# Row 139 M2 independent judge r1 (PR #288, origin/sprint/row139-m2-nominate)

Score: 78/100. FAIL (one BLOCKING, point A).

Gates (all rc checked, output to files): go vet rc0; targeted go test (registry, pins, daemon, transitionreg, world-publish, worldd) rc0; go test ./host/... ./cmd/... rc0; verify_ail.sh rc0 (16 identities, 40 named tests).

## BLOCKING

B1 (point A) cmd/world-publish/epoch.go loadReleaseManifest (and main.go --release-manifest flag): the verb accepts ANY file passing pins.Parse as the trust anchor, in headless --candidate mode. D-WORLD-78's digest condition ("binary digest matches GitHub's .sha256") roots in the reviewed, compiled-in pins/ailang-releases.json. An unattended caller can write a manifest with a row for any release, any digests, plus a verdict echoing them (golden_identical true), and nominate it with no human anywhere. Evidence: the tests themselves do exactly this (epoch_test.go writeTestManifestFor, v1_repro_test.go step 3 nominates headless with a fabricated manifest row); no test or fence refuses it. The compiled manifest is where a human reviews a merged row; the flag removes that review.
Minimal fix: delete the --release-manifest flag from options/flagNames/usage; production uses pins.Current() only. Tests get an unexported package var (e.g. `var releaseManifestOverride *pins.Manifest`, or a `runEpochNominateWith(manifest)` inner func) that is not reachable from argv. Add a test that `--release-manifest` is rejected as an unknown flag and that a verdict whose digests differ from the compiled manifest is refused. (Design D2's verb text lists --release-manifest, but that predates D-WORLD-78 = a; the ruling's digest condition is what governs and the design should be amended in the same PR.) Alternative acceptable: flag allowed only with the attended fence stack.

## Point B (boundary guard): NOT blocking
Guard (host/boundary/allowlist_world_test.go forbiddenImportPrefixes, added SM.A 13315da) is a lexical scan of each group's OWN source files' imports (checkGoGroup walks group.dir, ImportsOnly); the closure is only used for the transport check. Its rationale line reads "all direct registry/cloud imports": keeping the package-publish path off remote registry/cloud/network surfaces. host/registry is the local epoch-registry store object, no network. cmd/world-publish already depends transitively on host/registry through host/transitionreg/epoch.go (on dev since before), and `transitions` already writes registry-adjacent state. So the intent (no network/cloud from world-publish) is satisfied. Letter-wise it is a lexical bypass done silently, with aliases and no guard comment or test. Not blocking; recommended NB: add a comment on the prefix noting host/registry is the local epoch registry reached via transitionreg, or a narrow reviewed exception plus a test asserting world-publish's own files import neither host/registry nor net/http.

## AC verdicts
- AC2: immutable revisions (PutObject new object, prior never modified, test asserts bootstrap object unchanged), CAS (precheck + store CAS, both mutants killed), candidate (last epoch) vs new-epoch (N+1) shapes, refusals (already-nominated, both/none mode, unknown release, absent store/registry, verdict malformations, stale head, running daemon) all covered. PASS.
- AC3: existing store un-nominated refuses with the exact command (daemon.go NotNominatedError, tested incl. unpinned variant); nominated starts; fresh-store bootstrap bytes/hash/provenance pinned by literal test. V1 reproduction test exists (real pin; fake candidate unless WORLD_V1_CANDIDATE_BIN). PASS.
- Verdict schema/admit: unknown fields, trailing data, required booleans via pointers, all_pass consistency, unique check names, RFC3339, release/platform/binary+tarball digest vs manifest, golden_identical only for candidate. PASS.
- New-epoch fence: requireAttendedLocalWrite = identical automation/tty/typed-phrase stack as transitions; phrase names epoch+release+store. PASS.
- No release literals in production Go: host/pins scan green. PASS.

## NON-BLOCKING
1. Verdict is self-attested and unsigned; the verb never hashes any binary. Even with the compiled manifest, a forged verdict echoing manifest digests passes. Inherent to the ratified design until M3's script produces it; startup could also compare the served binary's sha256 to the manifest (not done). Record as a known trust residual.
2. Boundary-guard rationale comment (see B).
3. Real-candidate WORLD_V1_CANDIDATE_BIN run not executed by me (optional); the fake-binary path is green.
4. Daemon package test is ~200 s; fine.

## Mutations (17 run on a scratch copy; killed 17 / survived 0)
golden-check-off KILLED (TestAdmitVerdictReasons); binary-digest-off KILLED; tarball-digest-off KILLED; platform-off KILLED; release-mismatch-off KILLED; all_pass admit off KILLED; all_pass consistency (parse) off KILLED (TestParseVerdictControlAndMalformations); CAS precheck off KILLED (TestNominateRefusals); CAS write skipped KILLED (TestNominateCandidateAndNewEpochShapes); new-epoch phrase gate off KILLED (TestEpochNominateNewEpochIsAttended); startup accepts any release KILLED (TestBootstrapOnExistingHeadAdmitsNominatedReleases); candidate to first epoch KILLED; already-nominated off KILLED; both-modes allowed KILLED; manifest-lookup off KILLED; provenance swapped KILLED; epoch contiguity off KILLED (TestNextIsPure).
Note: no mutant targets the --release-manifest bypass because no test asserts its absence; that gap is B1.

Re-judge after B1 fix: expect PASS ~92.

## r2 (head f4f4884)

Score 92/100. PASS, zero blocking.

Gates (rc checked): go vet rc0; go test ./host/... ./cmd/... rc0 (includes boundary + census + host/pins scan); verify_ail.sh rc0.

Point A closed. No argv path names another manifest (flag deleted, test asserts it is unknown; nothing else in the epoch path reads one; runEpochNominateWith is unexported). Fabricated digests: the verdict must equal the compiled manifest, and --binary must hash to the compiled digest, so a forged verdict plus a fake binary is refused at --binary (TestEpochNominateAgainstTheCompiledManifest). Only a real reviewed-release binary can pass.

--binary bypass analysis (all NON-BLOCKING):
- Symlink: os.Open and exec both follow it; hash is of the target. A symlink to the real pin cannot make a different binary pass the hash; the probe runs the same target.
- TOCTOU hash-then-exec: a swap between them can only change what the version probe runs; the hash already matched real bytes, so nothing is gained. The probe is a consistency check, not a security fence.
- Path resolved differently: a bare name (--binary ailang) is opened relative to cwd for the hash but exec.Command does a PATH lookup. Again the hashed bytes are the real release, so no bypass, but hash and probe could describe different files. Fix: filepath.Abs before both (nit).
- The daemon later serves a binary by --ailang-bin and trusts its version line; the nomination does not bind served bytes to the digest. Residual trust gap, outside M2 (verdict remains self-attested for goldens and checks).
Subprocess census: no new exec site. archive.ProbeVersion reuses probeVersion (scrubbed env via childenv.Scrubbed, bounded timeout); epoch.go has no exec. Boundary comment at forbiddenImportPrefixes addresses point B.

Mutations on the new check (cmd/world-publish tests): version-check off KILLED; --binary-absent off KILLED; whole requireReleaseBinary skipped KILLED; both digest comparisons off KILLED. Dropping only the manifest compare, or only the verdict compare, SURVIVED: equivalent mutants, since AdmitVerdict already forces verdict digest == manifest digest, so got==verdict implies got==manifest. Redundant by transitivity, harmless.
