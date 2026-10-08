# Independent design/disposition review — iteration 246

Result: **ACCEPT design discovery and PARK disposition**, with the H2 qualification below. This is neither implementation acceptance nor sprint score. Judge independence: **FLAG same-model-fresh-context**, D-WORLD-48 fallback; no controller/generator reasoning used.

Reviewed base: `94fc5ba4d017505849e2adab1e87c8aae8ad0bdc` (independently read HEAD). Contract: tracked `design_docs/world-mission.md`, row 150, D-WORLD-23/24/48/64; umbrella D8; focused candidate `planned/w-marketing-capture.md`. No candidate or git edits.

## H1: independently justified

Row150 explicitly owns a real packet and every grade. D64 authorizes queue placement and mandates all four release rows; it does not narrow those requirements. D23 permits keeping a tranche from absorbing separately-owned work, with a named OPEN residual owner and a pinned composition condition. It does not license declaring row150 complete after deleting its expressly named packet/every-grade deliverables. D24 explicitly distinguishes shedding a core deliverable and records an attended decision for doing so. Row26 remains an unlanded owner of the shed bounded proof producer. Thus the alternative reduced showcase needs attended scope ratification; retaining full scope requires prerequisite routing/production composition rather than a counterfeit demo proof.

Independent scoped negative search (`host cmd`, non-test Go) finds no `DecisionPacket|decision-packet`; positive control finds `world/approval-request/v1` in broker/approve.go. World/types.ail declares the packet but expressly maps ProofReceipt to CLAIMED. Workbench object handler supplies NewGradeUnavailable (477), decisions use RecentApprovals (577). Non-test producer/consumer search finds proof codecs and NewValidator definition, no external production validator composition. Implemented evidence-boundary doc explicitly labels itself library-only/non-production (364–365). These independently support the discovery premise; absence is scoped to this landed repo, not all external systems.

## H2: operational facts confirmed, mandatory human judgment only partly established

CLI mint is attended/TTY fenced. A supplied valid session is an existing authorized normal path; preserving the fence needs no new human policy decision. Daemon coordinator uses time.Now().Unix (694), while Coordinator.Config already exposes required Now and constructs records with it (coordinator.go 48–79,354; effectful.go 284; plan.go 103–160). Current daemon path cannot guarantee fresh-seed public-byte identity across wall-clock seconds. Tests and source confirm this composition fact, not a ban on logical time: frozen DecisionPacket itself explicitly uses logical createdAt/deadlineAt (types.ail 155–168).

I found no explicit decision requiring human ratification of *every* bounded demo-local logical-clock adapter. A new production clock knob affecting ordinary invocation/expiry authority warrants the proposed freeze; an isolated orchestrator reusing the existing logical Now seam may be an ordinary design choice, subject to normal quorum and authority/expiry tests. Candidate should distinguish these. Do not raise a separate human ask merely to supply an attended session or reuse an existing logical-time seam. H1 independently sustains PARK, so this qualification does not change the disposition.

## Executed positive/negative controls

Two independent bounded invocations: Python subprocess timeout 110s; Go test timeout 90s, count=1; both rc=0. No whole-repo gates claimed.

1. host/evidence: TestMismatchedProofSubjectIsRefused, TestAttackerChosenValidatorCannotMintForHostAuthority, TestValidatorMintIdentitiesAreDistinct; each subject/authority test includes its PROVEN positive control. Also TestOtherwisePerfectReportWithoutMACIsUnauthenticated and WithWrongMACIsUnauthenticated: positive validated reports and forged/absent MAC refusal green.
2. CLI: TestSessionMint_PrintsOnceExactly positive; RequiresControllingTerminal negative; TestSessionNewProvisionsTheWorktreeAndTheSession positive; WithoutATerminalProvisionsNothing negative; ValidatesBeforeTheFence subcases green.
3. Why: TargetFormsResolveToOneEntry, ResolvesMCPInvocationID positive; BrokenLinkIsIntegrityRefusal and ScanBound negative/bounded controls green. Plain-object vs coordinator completeness additionally read in why.go; the candidate correctly demands coordinator flag and named chain checks rather than mere command success.

## Limitations and follow-up

No implementation exists to score. No new mutation, production packet ingress, six-shot product capture, fresh seed comparison, Chrome lifecycle, or all-grade renderer was run. Candidate appropriately marks these pending. Chrome timeout artifact is not successful lifecycle evidence. Normal quorum remains required by D64; this independent review does not ratify H1 or authorize scope deletion. Recommendation: preserve PARK under H1; refine H2 into a concrete alternative architecture and identify the precise new authority semantics, if any, before asking a human to resolve it.
