// Command world-publish is the ATTENDED operator entrypoint for the one
// irreversible effect in this repository: publishing world/core to the public
// AILANG registry.
//
// ---------------------------------------------------------------------------
// WHY THIS PROGRAM EXISTS, AND WHY IT IS A SEPARATE PROGRAM
//
// Until this milestone, the ABSENCE of a caller WAS the safety property.
// host/broker/registry_publish.go's production constructor demands https and
// refuses loopback, and every caller in the tree was an httptest server reached
// through the UNEXPORTED loopback constructor — so nothing here could publish
// and no headless process could trip it. Measured at 6d1dce0: zero Publish or
// Approve identifiers anywhere under cmd/ (control: 27 func declarations), and
// exactly one command in the repository.
//
// That made the attended runbook's Stage B a sequence of PROSE. On 2026-08-10
// the attended session could not execute it — not for want of a human decision,
// but because the code that performs it did not exist. This program is that
// code, and its central deliverable is a FENCE, not a feature.
//
// It is deliberately NOT a verb on ailang-worldd. That daemon is what CI, the
// bench and the autonomous loop already execute; a publish verb there would put
// every existing invocation one flag away from the irreversible path. A
// separate program keeps the blast radius at "a human typed a different program
// name".
//
// THE DEFAULT OUTCOME OF EVERY PATH BELOW IS STOP (exit 3).
//
// ---------------------------------------------------------------------------
// THE FOUR VERBS
//
//	packet     recompute the projection's identity and compare it with the
//	           committed golden. READ-ONLY, headless-permitted.
//	approve    mint the one-shot attended approval. ATTENDED.
//	publish    spend it, exactly once. ATTENDED. IRREVERSIBLE under --live.
//	reconcile  report indeterminate publish receipts, and (with --probe) resolve
//	           them with read-only GETs. READ-ONLY, headless-permitted.
//
// Minting and spending are TWO invocations on purpose (Decision D4). One
// command that minted-then-spent would collapse two human acts into one
// keystroke and make the single-use property invisible at the operator's
// surface. Two invocations mean the operator can review the minted ref, and the
// durable claim is spent by a command that could not have created it.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/pkgproj"
	"github.com/sunholo-data/ailang-world/host/store"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, environment{
		getenv: os.Getenv,
		probe:  probeControllingTerminal,
	}))
}

// environment is everything the command reads from OUTSIDE its flags. It is a
// struct so the fence rows can be driven in-process with an injected probe and
// an injected getenv, rather than by hoping a subprocess is in the right state.
//
// There is deliberately no field here that supplies a credential. The registry
// API key is read from a FILE named by --credential-file and from nowhere else;
// broker.AssertNoAmbientRegistryCredential (called inside the production
// constructor) refuses to build a handler at all if the variable is set in this
// process's environment.
type environment struct {
	getenv func(string) string
	probe  func() ttyProbe
}

// options is the complete flag surface.
type options struct {
	store          string
	packageDir     string
	golden         string
	registryOrigin string
	publisher      string
	credentialFile string
	approvalRef    string
	episode        string
	requester      string
	decidedBy      string
	now            int64
	expires        int64
	dryRun         bool
	live           bool
	probe          bool
}

// flagNames is the FROZEN flag surface, stated once and INDEPENDENTLY of the
// FlagSet so the two can be compared by exact set equality.
//
// THERE IS NO FLAG THAT SETS THE VALIDATOR ORIGIN, and there must never be one.
// broker.ApprovedValidatorOrigin is a compiled constant; a --validator-origin
// flag, however convenient for testing, would make the production constructor's
// one unwidenable guarantee configurable from the command line. AC24(b)
// compares this list with the FlagSet as an exact set, so ADDING a flag reds
// even if it is never passed.
var flagNames = []string{
	"approval-ref",
	"credential-file",
	"decided-by",
	"dry-run",
	"episode",
	"expires",
	"golden",
	"live",
	"now",
	"package-dir",
	"probe",
	"publisher",
	"registry-origin",
	"requester",
	"store",
}

const (
	defaultPackageDir = "packages/world-core"
	defaultGolden     = "scripts/world_package_ready_packet.golden.json"
)

func newFlagSet(opts *options) *flag.FlagSet {
	fs := flag.NewFlagSet("world-publish", flag.ContinueOnError)
	fs.StringVar(&opts.store, "store", "", "path to the world database this publish is recorded in")
	fs.StringVar(&opts.packageDir, "package-dir", defaultPackageDir, "the projected package directory")
	fs.StringVar(&opts.golden, "golden", defaultGolden, "the committed ready-packet golden")
	fs.StringVar(&opts.registryOrigin, "registry-origin", "", "the read-only public bucket origin")
	fs.StringVar(&opts.publisher, "publisher", "", "path to the pinned released ailang binary")
	fs.StringVar(&opts.credentialFile, "credential-file", "",
		"file outside the working tree holding the registry API key")
	fs.StringVar(&opts.approvalRef, "approval-ref", "",
		"the ApprovalDecisionV1 hashref minted by `world-publish approve`")
	fs.StringVar(&opts.episode, "episode", "attended-publish", "durable episode ID")
	fs.StringVar(&opts.requester, "requester", "", "who is requesting the approval")
	fs.StringVar(&opts.decidedBy, "decided-by", "", "who is granting the approval")
	fs.Int64Var(&opts.now, "now", 0, "logical time of this act")
	fs.Int64Var(&opts.expires, "expires", 0, "logical time the approval expires at")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "rehearse the publish: every fence, no request")
	fs.BoolVar(&opts.live, "live", false, "PERFORM THE IRREVERSIBLE PUBLIC WRITE")
	fs.BoolVar(&opts.probe, "probe", false, "reconcile: issue the read-only metadata GETs")
	return fs
}

const usageText = `world-publish — the attended entrypoint for an IRREVERSIBLE public write

  world-publish packet     [--package-dir D] [--golden G]
  world-publish approve    --store S [--registry-origin O] --now N --expires E
  world-publish publish    --store S --registry-origin O --publisher P \
                           --credential-file C --approval-ref R --now N --expires E (--live | --dry-run)
  world-publish reconcile  --store S [--registry-origin O] [--probe]

Exit codes: 0 done · 1 failed · 2 usage · 3 STOP (a fence refused; nothing happened)
`

func run(args []string, in io.Reader, out, errw io.Writer, env environment) int {
	if len(args) == 0 {
		fmt.Fprint(errw, usageText)
		return exitUsage
	}
	verb := args[0]
	var opts options
	fs := newFlagSet(&opts)
	fs.SetOutput(errw)
	if err := fs.Parse(args[1:]); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(errw, "world-publish: unexpected argument %q\n", fs.Arg(0))
		return exitUsage
	}
	switch verb {
	case "packet":
		return runPacket(opts, out, errw)
	case "approve":
		return runApprove(opts, in, out, errw, env)
	case "publish":
		return runPublish(opts, in, out, errw, env)
	case "reconcile":
		return runReconcile(opts, out, errw)
	default:
		fmt.Fprintf(errw, "world-publish: unknown verb %q\n\n%s", verb, usageText)
		return exitUsage
	}
}

// ---------------------------------------------------------------------------
// packet — read-only, headless-permitted
// ---------------------------------------------------------------------------

func runPacket(opts options, out, errw io.Writer) int {
	if serr := refuseLiveOnReadOnlyVerb("packet", opts.live); serr != nil {
		return report(errw, serr)
	}
	packet, serr := requireUndriftedPacket(opts.packageDir, opts.golden)
	if serr != nil {
		return report(errw, serr)
	}
	printPacket(out, packet)
	return exitOK
}

func printPacket(out io.Writer, packet pkgproj.ReadyPacket) {
	fmt.Fprintf(out, "ready packet for %s@%s (compiler %s)\n",
		packet.Package, packet.Version, packet.CompilerVersion)
	for _, name := range pkgproj.ReadyPacketFields {
		fmt.Fprintf(out, "  %-16s %s\n", name, packet.Field(name))
	}
	fmt.Fprintln(out, "  recomputed from disk and EQUAL to the committed golden in every field")
}

// ---------------------------------------------------------------------------
// approve — attended
// ---------------------------------------------------------------------------

func runApprove(opts options, in io.Reader, out, errw io.Writer, env environment) int {
	if serr := requireStorePath(opts.store); serr != nil {
		return report(errw, serr)
	}
	packet, serr := requireUndriftedPacket(opts.packageDir, opts.golden)
	if serr != nil {
		return report(errw, serr)
	}
	plan, err := attendedPlan(opts, packet)
	if err != nil {
		fmt.Fprintln(errw, "world-publish: "+err.Error())
		return exitError
	}
	fmt.Fprintf(out, "about to mint a ONE-SHOT approval for:\n  %s\n", plan.ApprovalScope())
	// Minting authority for an irreversible act is itself an attended act: a
	// headless loop that could mint its own stamp would have defeated the whole
	// approval layer, which is single-use but NOT headless/attended.
	if serr := requireAttendedOperator(in, out, env.getenv, env.probe()); serr != nil {
		return report(errw, serr)
	}
	db, err := store.Open(opts.store)
	if err != nil {
		return report(errw, storeOpenFailed(opts.store, err))
	}
	defer func() { _ = db.Close() }()

	ref, err := broker.MintAttendedApproval(db, plan)
	if err != nil {
		fmt.Fprintln(errw, "world-publish: mint approval: "+err.Error())
		return exitError
	}
	fmt.Fprintf(out, "minted approval %s\n", ref.String())
	fmt.Fprintf(out, "spend it exactly once with:\n"+
		"  go run ./cmd/world-publish publish --live --approval-ref %s ...\n", ref.String())
	return exitOK
}

// ---------------------------------------------------------------------------
// publish — attended, and the only irreversible path in this repository
// ---------------------------------------------------------------------------

func runPublish(opts options, in io.Reader, out, errw io.Writer, env environment) int {
	if serr := requireExactlyOneMode(opts.dryRun, opts.live); serr != nil {
		return report(errw, serr)
	}
	if serr := requireStorePath(opts.store); serr != nil {
		return report(errw, serr)
	}
	// The store is opened HERE, above the fence, deliberately: it is a local
	// SQLite file, it carries no secret and it reaches no network, and an
	// operator who mistyped the path should learn that before being asked to
	// type an irreversible confirmation phrase.
	db, err := store.Open(opts.store)
	if err != nil {
		return report(errw, storeOpenFailed(opts.store, err))
	}
	defer func() { _ = db.Close() }()

	var approvalRef hashref.HashRef
	if opts.live {
		var serr *stopError
		if approvalRef, serr = requireApprovalRef(opts.approvalRef); serr != nil {
			return report(errw, serr)
		}
		if serr := requireCredentialFile(opts.credentialFile); serr != nil {
			return report(errw, serr)
		}
	}
	packet, serr := requireUndriftedPacket(opts.packageDir, opts.golden)
	if serr != nil {
		return report(errw, serr)
	}
	plan, planErr := attendedPlan(opts, packet)
	if planErr != nil {
		fmt.Fprintln(errw, "world-publish: "+planErr.Error())
		return exitError
	}
	credential := &deferredFileCredential{path: opts.credentialFile}
	cfg, cfgErr := publishHandlerConfig(opts, plan, credential)
	if cfgErr != nil {
		fmt.Fprintln(errw, "world-publish: "+cfgErr.Error())
		return exitError
	}
	// ======================================================================
	// THE FENCE. Everything above this line is local, reversible and cheap.
	// Everything below it can reach the network.
	//
	// AC22 asserts over the AST that this call and the
	// broker.NewRegistryPublishHandler call below each appear EXACTLY ONCE in
	// this function and that this one's statement index is strictly lower.
	// MUT-D0-FENCE-ORDER swaps the two. The mutant BUILDS (rc=0) and AC22 reds,
	// naming the ordering defect structurally — which is why the guard is not a
	// comment.
	//
	// CORRECTED (controller iter-67, measured — the original comment here, and
	// the same sentence in wiring_test.go and the plan, claimed AC22 was the
	// ONLY killer and that "every AC21 row still passes"). It is not: under the
	// mutant 6 of 15 AC21 rows ALSO red — R-CI, R-TTY-OPEN, R-TTY-CHARDEV,
	// R-TTY-SAMEFILE, R-PHRASE-EOF, R-PHRASE. The cause is worth knowing: those
	// rows' fixture supplies a LOOPBACK registry origin as a defensive
	// baseline, so with the constructor hoisted above the fence its own
	// loopback/ambient-credential refusal fires first and the row observes
	// "STOP fence=handler" instead of its documented line. AC22 still earns its
	// place — it names the defect precisely, and it would red even against a
	// fixture using a valid https origin, where every AC21 row WOULD pass. But
	// it is not the unique killer, and a non-vacuity claim written down without
	// being run as literally described is the exact class this mission keeps
	// closing.
	// ======================================================================
	if serr := requireAttendedOperator(in, out, env.getenv, env.probe()); serr != nil {
		return report(errw, serr)
	}
	if opts.dryRun {
		// --dry-run is the attended REHEARSAL: every fence above, then a full
		// statement of what --live would do, and no request of any kind. It is
		// not the pinned compiler's `publish --dry-run`; that is step 7 of
		// scripts/verify_world_package.sh and it needs no attended operator.
		printRehearsal(out, plan, cfg)
		return exitOK
	}
	handler, handlerErr := broker.NewRegistryPublishHandler(cfg)
	if handlerErr != nil {
		// Every refusal reachable here is INHERITED — loopback/https/wildcard/
		// userinfo origin refusals, the ambient-credential refusal, the
		// missing-approval and missing-credential-provider refusals. They are
		// reported, never restated.
		return report(errw, &stopError{Fence: fenceHandler, Reason: "refused",
			Detail: handlerErr.Error()})
	}
	// Pre-flight the credential file with the LANDED provider constructor, so a
	// typo costs a refusal rather than the one-shot approval: Session.Invoke
	// consumes the durable claim BEFORE the handler reads the secret.
	if _, err := credential.provider(); err != nil {
		return report(errw, &stopError{Fence: fenceCredential, Reason: "unusable",
			Detail: err.Error()})
	}

	fmt.Fprintf(out, "publishing %s@%s to %s\n", plan.Identity.Vendor+"/"+plan.Identity.Name,
		plan.Identity.Version, plan.Identity.RegistryOrigin)
	result, invokeErr := broker.InvokeAttendedPublish(context.Background(), db, handler, plan, approvalRef)
	return reportPublishResult(out, errw, result, invokeErr)
}

func reportPublishResult(out, errw io.Writer, result broker.AttendedPublishResult, err error) int {
	switch result.Status {
	case broker.AttendedPublishSucceeded:
		fmt.Fprintf(out, "PUBLISHED. record %s\n", result.RecordRef.String())
		return exitOK
	case broker.AttendedPublishIndeterminate:
		// The one outcome that must not be retried. The request body may already
		// have reached a registry that cannot be asked to take it back.
		fmt.Fprintf(errw, "INDETERMINATE. episode %s ordinal %d invocation %s\n",
			result.EpisodeID, result.Ordinal, result.InvocationID)
		fmt.Fprintln(errw, "  DO NOT RETRY. Resolve it read-only:")
		fmt.Fprintln(errw, "  go run ./cmd/world-publish reconcile --store <db> --registry-origin <bucket> --probe")
		return exitError
	default:
		fmt.Fprintf(errw, "FAILED. %v\n", err)
		return exitError
	}
}

func printRehearsal(out io.Writer, plan broker.AttendedPublishPlan, cfg broker.RegistryPublishConfig) {
	fmt.Fprintln(out, "REHEARSAL — no request of any kind was made.")
	fmt.Fprintf(out, "  package        %s/%s@%s\n", plan.Identity.Vendor, plan.Identity.Name, plan.Identity.Version)
	fmt.Fprintf(out, "  registry       %s\n", plan.Identity.RegistryOrigin)
	fmt.Fprintf(out, "  validator      %s\n", cfg.ValidatorOrigin)
	fmt.Fprintf(out, "  publisher      %s\n", cfg.PublisherPath)
	fmt.Fprintf(out, "  capability     %s\n", plan.PublishScope())
	fmt.Fprintf(out, "  approval scope %s\n", plan.ApprovalScope())
	fmt.Fprintf(out, "  tarball        %s\n", plan.Hashes.TarballSHA256)
	fmt.Fprintf(out, "  content        %s\n", plan.Hashes.ContentHash)
	fmt.Fprintf(out, "  interface      %s\n", plan.Hashes.InterfaceHash)
}

// ---------------------------------------------------------------------------
// reconcile — read-only, headless-permitted
// ---------------------------------------------------------------------------

func runReconcile(opts options, out, errw io.Writer) int {
	if serr := refuseLiveOnReadOnlyVerb("reconcile", opts.live); serr != nil {
		return report(errw, serr)
	}
	if serr := requireStorePath(opts.store); serr != nil {
		return report(errw, serr)
	}
	packet, serr := requireUndriftedPacket(opts.packageDir, opts.golden)
	if serr != nil {
		return report(errw, serr)
	}
	db, err := store.Open(opts.store)
	if err != nil {
		return report(errw, storeOpenFailed(opts.store, err))
	}
	defer func() { _ = db.Close() }()

	pending, err := db.PendingEffectIntents(reconcileScanLimit)
	if err != nil {
		fmt.Fprintln(errw, "world-publish: scan pending effect intents: "+err.Error())
		return exitError
	}
	outstanding := 0
	for _, intent := range pending {
		if intent.Intent.Effect != broker.EffectRegistryPublish {
			continue
		}
		outstanding++
		fmt.Fprintf(out, "INDETERMINATE %s\n  scope %s\n", intent.InvocationID, intent.Intent.Scope)
	}
	fmt.Fprintf(out, "%d indeterminate %s receipt(s) in %s\n",
		outstanding, broker.EffectRegistryPublish, opts.store)
	if !opts.probe {
		fmt.Fprintln(out, "no probe issued (--probe requests the read-only metadata GETs)")
		return exitOK
	}
	receipt, err := broker.ReconcileRegistryPublish(context.Background(), reconcileConfigFor(opts, packet))
	if err != nil {
		fmt.Fprintln(errw, "world-publish: reconcile: "+err.Error())
		return exitError
	}
	fmt.Fprintln(out, receipt.String())
	return exitOK
}

// reconcileScanLimit bounds one pending-intent page. It is a constant rather
// than a flag because a flag here would be a way to make the report LOOK empty.
const reconcileScanLimit = 200

// reconcileConfigFor is pure, so the arguments handed to the production
// reconciler can be asserted without issuing a single GET.
//
// The expected digests come from the READY PACKET, not from the durable intent:
// reconciliation asks "does the public bucket serve the bytes we reviewed", and
// answering it from the same place the publish was authorized from is what
// makes `conflict` distinguishable from `succeeded-reconciled`.
//
// There is deliberately no field here for a control package or a ready-made
// URL. host/broker builds both the target and its same-pass known-positive
// control from ONE origin, differing only in vendor/name/version, so no caller
// can point the control at a different key-space than the target.
func reconcileConfigFor(opts options, packet pkgproj.ReadyPacket) broker.ReconcileConfig {
	vendor, name := splitPackage(packet.Package)
	return broker.ReconcileConfig{
		RegistryOrigin: opts.registryOrigin,
		Vendor:         vendor,
		Name:           name,
		Version:        packet.Version,
		Expected: broker.PublishHashes{
			TarballSHA256: packet.TarballSHA256,
			ContentHash:   packet.ContentHash,
			InterfaceHash: packet.InterfaceHash,
		},
	}
}

// ---------------------------------------------------------------------------
// the plan and the handler config
// ---------------------------------------------------------------------------

// attendedPlan builds the ONE plan the mint and the spend must agree on. It is
// derived entirely from the READY PACKET, so `approve` and `publish` cannot
// disagree about what is being published without one of them failing
// R-PACKET-DRIFT first.
func attendedPlan(opts options, packet pkgproj.ReadyPacket) (broker.AttendedPublishPlan, error) {
	vendor, name := splitPackage(packet.Package)
	if vendor == "" || name == "" {
		return broker.AttendedPublishPlan{}, fmt.Errorf(
			"ready packet names package %q, which is not <vendor>/<name>", packet.Package)
	}
	manifestRef, err := manifestReference(opts.packageDir)
	if err != nil {
		return broker.AttendedPublishPlan{}, err
	}
	compilerSHA, err := publisherDigest(opts.publisher)
	if err != nil {
		return broker.AttendedPublishPlan{}, err
	}
	return broker.AttendedPublishPlan{
		Identity: broker.PublishIdentity{
			Vendor:          vendor,
			Name:            name,
			Version:         packet.Version,
			RegistryOrigin:  opts.registryOrigin,
			ManifestRef:     manifestRef,
			Exports:         packet.Exports,
			Effects:         packet.Effects,
			CompilerVersion: packet.CompilerVersion,
			CompilerSHA256:  compilerSHA,
		},
		Hashes: broker.PublishHashes{
			TarballSHA256: packet.TarballSHA256,
			ContentHash:   packet.ContentHash,
			InterfaceHash: packet.InterfaceHash,
		},
		Requester:   opts.requester,
		DecidedBy:   opts.decidedBy,
		EpisodeID:   opts.episode,
		RequestedAt: opts.now,
		DecidedAt:   opts.now,
		PublishAt:   opts.now,
		ExpiresAt:   opts.expires,
	}, nil
}

// publishHandlerConfig is PURE: it performs no I/O and cannot fail for any
// reason outside its arguments, which is what lets the fence and the
// constructor sit as two adjacent statements in runPublish.
//
// ValidatorOrigin is broker.ApprovedValidatorOrigin, a compiled constant, and
// there is no flag that reaches it.
func publishHandlerConfig(
	opts options, plan broker.AttendedPublishPlan, credential broker.RegistryCredentialProvider,
) (broker.RegistryPublishConfig, error) {
	publisher, err := filepath.Abs(opts.publisher)
	if err != nil {
		return broker.RegistryPublishConfig{}, fmt.Errorf("resolve --publisher: %w", err)
	}
	packageDir, err := filepath.Abs(opts.packageDir)
	if err != nil {
		return broker.RegistryPublishConfig{}, fmt.Errorf("resolve --package-dir: %w", err)
	}
	return broker.RegistryPublishConfig{
		PublisherPath:   publisher,
		PackageDir:      packageDir,
		Manifest:        worldCoreManifest,
		RegistryOrigin:  opts.registryOrigin,
		ValidatorOrigin: broker.ApprovedValidatorOrigin,
		Credential:      credential,
		Approval: broker.PublishApproval{
			ApprovalRef:    mustParseRefOrZero(opts.approvalRef),
			Vendor:         plan.Identity.Vendor,
			Name:           plan.Identity.Name,
			Version:        plan.Identity.Version,
			RegistryOrigin: opts.registryOrigin,
			Hashes:         plan.Hashes,
		},
	}, nil
}

// deferredFileCredential defers broker.NewFileRegistryCredentialProvider to
// first use.
//
// It exists so publishHandlerConfig can be pure. That is not cosmetic: an eager
// provider construction between the fence and the constructor would make the
// two non-adjacent, and MUT-D0-FENCE-ORDER — which must BUILD in order to prove
// anything — could not swap them. It also STRENGTHENS the ordering it is there
// to preserve: the credential file is not even stat'ed until the pre-flight
// below the fence.
type deferredFileCredential struct {
	path  string
	cache broker.RegistryCredentialProvider
}

func (c *deferredFileCredential) provider() (broker.RegistryCredentialProvider, error) {
	if c.cache != nil {
		return c.cache, nil
	}
	root, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("resolve working directory: %w", err)
	}
	provider, err := broker.NewFileRegistryCredentialProvider(c.path, root)
	if err != nil {
		return nil, err
	}
	c.cache = provider
	return provider, nil
}

func (c *deferredFileCredential) Credential() ([]byte, error) {
	provider, err := c.provider()
	if err != nil {
		return nil, err
	}
	return provider.Credential()
}

func splitPackage(pkg string) (vendor, name string) {
	vendor, name, _ = strings.Cut(pkg, "/")
	return vendor, name
}

// manifestReference content-addresses the manifest the approval scope binds. It
// is the ailang.toml on disk, so an approval minted against one manifest cannot
// be spent against another.
func manifestReference(packageDir string) (hashref.HashRef, error) {
	data, err := os.ReadFile(filepath.Join(packageDir, "ailang.toml"))
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("read package manifest: %w", err)
	}
	return hashref.SumSHA256(data), nil
}

// publisherDigest records WHICH binary performed the publish. It is provenance
// carried in the payload; it is not part of the approval scope, so `approve`
// (which needs no --publisher) and `publish` mint and spend the same stamp.
func publisherDigest(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read --publisher: %w", err)
	}
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// mustParseRefOrZero keeps publishHandlerConfig pure. An unparseable ref cannot
// reach it: R-APPROVAL-ABSENT has already refused one for every --live
// invocation, and a --dry-run rehearsal never constructs a handler.
func mustParseRefOrZero(raw string) hashref.HashRef {
	ref, err := hashref.Parse(strings.TrimSpace(raw))
	if err != nil {
		return hashref.HashRef{}
	}
	return ref
}
