package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sunholo-data/ailang-world/host/childenv"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/pkgproj"
)

// EffectRegistryPublish is the one effect in this repository that can perform
// an irreversible public write. The registry is immutable and the publisher
// cannot recall a version, so every refusal below is one-way: a false REFUSAL
// costs a retry, a false ALLOW costs a permanent public artifact.
const EffectRegistryPublish = "Registry.Publish"

// PublishCost is the frozen cost of one live publish attempt (Decision 3).
// The unit is one irreversible HTTP publish attempt, charged BEFORE dispatch,
// and an attended grant starts at budget 1 — so success, definite failure and
// indeterminate dispatch all consume the same single unit.
const PublishCost = 1

// PublishPayloadSchema tags the canonical fixed-field-order request payload.
const PublishPayloadSchema = "world/registry-publish-request/v1"

// PublishResultSchema tags the canonical result object.
const PublishResultSchema = "world/registry-publish-result/v1"

// ApprovedValidatorOrigin is the ONLY validator origin the production
// constructor accepts. Measured, not assumed: it is v0.30.0's compiled
// DefaultValidatorURL at e37b370:cmd/ailang/pkg_info.go:14.
const ApprovedValidatorOrigin = "https://registry.ailang.sunholo.com"

// forbiddenPublisherFlags are publisher flags this handler refuses to pass on.
// --allow-dotted-tool-names downgrades v0.30.0's naming gate from error to
// warning (e37b370:cmd/ailang/pkg_publish.go:25-27) — i.e. it makes the
// registry accept a package the registry's own validator would otherwise
// reject, permanently.
var forbiddenPublisherFlags = []string{"--allow-dotted-tool-names"}

// Publish outcome classes. These are the receipt-facing statuses of Decision
// 3's definite/ambiguous table, and the split between the first three and the
// fourth is the whole point: the first three are RESOLVED, the fourth is not.
const (
	PublishStatusSucceeded            = "succeeded"
	PublishStatusFailedBeforeDispatch = "failed-before-dispatch"
	PublishStatusFailed               = "failed"
	PublishStatusIndeterminate        = "indeterminate"
)

// PublishScope builds the frozen scope grammar of Decision 3:
//
//	registry:<registry-origin>/package:<vendor>/<name>/version:<version>
//
// The version is part of the scope, so an attended grant for 0.1.0 cannot
// authorize 0.1.1 — see MUT-SM-GRANT-SCOPE.
func PublishScope(registryOrigin, vendor, name, version string) string {
	return "registry:" + registryOrigin + "/package:" + vendor + "/" + name + "/version:" + version
}

// PublishIdentity is the caller-supplied publication identity. Every field
// lands in the canonical payload and is therefore bound by the EffectRecord's
// content-addressed RequestRef.
type PublishIdentity struct {
	Vendor          string
	Name            string
	Version         string
	RegistryOrigin  string
	ManifestRef     hashref.HashRef
	ApprovalRef     hashref.HashRef
	Exports         []string
	Effects         []string
	CompilerVersion string
	CompilerSHA256  string
}

// PublishHashes are the three package digests recomputed immediately before
// dispatch by host/pkgproj.
type PublishHashes struct {
	TarballSHA256 string
	ContentHash   string
	InterfaceHash string
}

// PublishApproval is the attended one-shot stamp. SM.B2a consumes only its
// identity and hash fields; SM.B2b binds ApprovalRef to the durable
// single-use claim.
type PublishApproval struct {
	ApprovalRef    hashref.HashRef
	Vendor         string
	Name           string
	Version        string
	RegistryOrigin string
	Hashes         PublishHashes
}

// publishPayloadWire is the canonical request payload. THE FIELD ORDER IS
// FROZEN and is exactly Decision 3's list; encoding/json emits struct fields
// in declaration order, so this declaration IS the wire order. It is pinned by
// TestPublishPayloadFieldOrderIsFrozen, which reads the key order back out of
// the encoded bytes rather than trusting this comment.
type publishPayloadWire struct {
	Schema          string   `json:"schema"`
	Vendor          string   `json:"vendor"`
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	RegistryOrigin  string   `json:"registryOrigin"`
	ManifestRef     string   `json:"manifestRef"`
	ApprovalRef     string   `json:"approvalRef"`
	TarballSHA256   string   `json:"tarballSHA256"`
	ContentHash     string   `json:"contentHash"`
	InterfaceHash   string   `json:"interfaceHash"`
	Exports         []string `json:"exports"`
	Effects         []string `json:"effects"`
	CompilerVersion string   `json:"compilerVersion"`
	CompilerSHA256  string   `json:"compilerSHA256"`
}

// publishPayloadFields is the frozen field order stated once, independently of
// the struct, so the two can be compared. A single source would make the
// ordering test a tautology.
var publishPayloadFields = []string{
	"schema", "vendor", "name", "version", "registryOrigin",
	"manifestRef", "approvalRef", "tarballSHA256", "contentHash", "interfaceHash",
	"exports", "effects", "compilerVersion", "compilerSHA256",
}

// EncodePublishPayload is the single canonical publish payload codec. Like
// EncodeRecord it uses a positional literal, so adding a field to the wire
// struct without deciding its position fails to compile.
func EncodePublishPayload(id PublishIdentity, hashes PublishHashes) []byte {
	payload, err := json.Marshal(publishPayloadWire{
		PublishPayloadSchema, id.Vendor, id.Name, id.Version, id.RegistryOrigin,
		id.ManifestRef.String(), id.ApprovalRef.String(),
		hashes.TarballSHA256, hashes.ContentHash, hashes.InterfaceHash,
		exportedList(id.Exports), exportedList(id.Effects),
		id.CompilerVersion, id.CompilerSHA256,
	})
	if err != nil {
		panic("broker: fixed publish payload cannot fail JSON encoding: " + err.Error())
	}
	return payload
}

// exportedList normalizes nil to an empty list so an absent and an empty
// export set encode identically, and a payload never carries JSON null where
// the registry expects an array.
func exportedList(items []string) []string {
	if items == nil {
		return []string{}
	}
	return append([]string(nil), items...)
}

// DecodePublishPayload decodes and validates the canonical payload. Unknown
// fields are rejected: a payload the handler cannot fully account for must not
// authorize an irreversible write.
func DecodePublishPayload(payload []byte) (PublishIdentity, PublishHashes, error) {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	var wire publishPayloadWire
	if err := dec.Decode(&wire); err != nil {
		return PublishIdentity{}, PublishHashes{}, fmt.Errorf("broker: decode publish payload: %w", err)
	}
	if wire.Schema != PublishPayloadSchema {
		return PublishIdentity{}, PublishHashes{}, fmt.Errorf(
			"broker: publish payload schema %q, want %q", wire.Schema, PublishPayloadSchema)
	}
	manifestRef, err := parseRequiredRef("manifestRef", wire.ManifestRef)
	if err != nil {
		return PublishIdentity{}, PublishHashes{}, err
	}
	approvalRef, err := parseRequiredRef("approvalRef", wire.ApprovalRef)
	if err != nil {
		return PublishIdentity{}, PublishHashes{}, err
	}
	return PublishIdentity{
			Vendor: wire.Vendor, Name: wire.Name, Version: wire.Version,
			RegistryOrigin: wire.RegistryOrigin,
			ManifestRef:    manifestRef, ApprovalRef: approvalRef,
			Exports: wire.Exports, Effects: wire.Effects,
			CompilerVersion: wire.CompilerVersion, CompilerSHA256: wire.CompilerSHA256,
		}, PublishHashes{
			TarballSHA256: wire.TarballSHA256,
			ContentHash:   wire.ContentHash,
			InterfaceHash: wire.InterfaceHash,
		}, nil
}

// publishApprovalRef extracts the approval reference the durable single-use
// claim is keyed on, without interpreting anything else in the payload. It is
// the broker's only read of a publish payload; the handler owns the rest.
func publishApprovalRef(payload []byte) (hashref.HashRef, error) {
	id, _, err := DecodePublishPayload(payload)
	if err != nil {
		return hashref.HashRef{}, err
	}
	return id.ApprovalRef, nil
}

// publishResultWire is the canonical success/definite-failure result object.
// PublisherOutput is redacted before it reaches this struct.
type publishResultWire struct {
	Schema          string `json:"schema"`
	Status          string `json:"status"`
	Vendor          string `json:"vendor"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	RegistryOrigin  string `json:"registryOrigin"`
	ValidatorOrigin string `json:"validatorOrigin"`
	TarballSHA256   string `json:"tarballSHA256"`
	ContentHash     string `json:"contentHash"`
	InterfaceHash   string `json:"interfaceHash"`
	PublisherOutput string `json:"publisherOutput"`
}

// PublishRefusalError is a handler refusal raised BEFORE any subprocess is
// launched. Every instance of it means the dispatch counter is still zero.
type PublishRefusalError struct {
	Why string
}

func (e *PublishRefusalError) Error() string {
	return "broker: Registry.Publish refused: " + e.Why
}

// RegistryPublishConfig pins everything the handler is allowed to do. There is
// deliberately no field that lets a caller add arbitrary environment
// variables: the child environment is constructed here and nowhere else.
type RegistryPublishConfig struct {
	// PublisherPath is the pinned released AILANG binary.
	PublisherPath string
	// PackageDir is the projection directory whose hashes are recomputed.
	PackageDir string
	// Manifest is the exact package manifest the interface hash is computed
	// from. It is cross-checked against the payload's exports and effects.
	Manifest pkgproj.Manifest
	// RegistryOrigin is $AILANG_REGISTRY for the child (the read-only bucket).
	RegistryOrigin string
	// ValidatorOrigin is $AILANG_REGISTRY_VALIDATOR for the child (the
	// publish service).
	ValidatorOrigin string
	// Credential is consulted exactly once per live dispatch and must be nil
	// for a dry-run handler.
	Credential RegistryCredentialProvider
	// Approval is the attended stamp the recomputed hashes must match.
	Approval PublishApproval
	// ExtraArgs are additional publisher flags. forbiddenPublisherFlags are
	// rejected at construction.
	ExtraArgs []string
	// DryRun runs the publisher with --dry-run, with no credential provider
	// and with the credential variable unset in the child.
	DryRun bool

	ExecTimeout    time.Duration
	MaxOutputBytes int64
}

// RegistryPublishHandler performs one brokered publish attempt by launching
// the pinned binary in an explicitly constructed minimal environment.
//
// It contains no authority logic: the Session has already decided effect,
// scope, expiry and budget by the time Execute runs. Everything here is
// defence in depth on an operation that cannot be undone.
type RegistryPublishHandler struct {
	cfg            RegistryPublishConfig
	allowLoopback  bool
	bounds         handlerBounds
	dispatches     atomic.Int64
	credentialLoad atomic.Int64
}

// NewRegistryPublishHandler is the PRODUCTION constructor. It accepts only the
// compiled/approved public validator origin; there is no flag, environment
// variable or configuration value that widens it.
func NewRegistryPublishHandler(cfg RegistryPublishConfig) (*RegistryPublishHandler, error) {
	if cfg.ValidatorOrigin != ApprovedValidatorOrigin {
		return nil, &PublishRefusalError{Why: fmt.Sprintf(
			"production validator origin is %q, got %q", ApprovedValidatorOrigin, cfg.ValidatorOrigin)}
	}
	if err := AssertNoAmbientRegistryCredential(os.Environ()); err != nil {
		return nil, err
	}
	return newRegistryPublishHandler(cfg, false)
}

// newLoopbackRegistryPublishHandler is the TEST-ONLY constructor named by
// Decision 3. It is UNEXPORTED on purpose: no package outside host/broker can
// build a publisher pointed anywhere but ApprovedValidatorOrigin, and this
// constructor additionally REQUIRES both origins to be loopback, so a test
// cannot reach the public registry even by mistake.
func newLoopbackRegistryPublishHandler(cfg RegistryPublishConfig) (*RegistryPublishHandler, error) {
	return newRegistryPublishHandler(cfg, true)
}

func newRegistryPublishHandler(cfg RegistryPublishConfig, allowLoopback bool) (*RegistryPublishHandler, error) {
	if !filepath.IsAbs(cfg.PublisherPath) {
		return nil, &PublishRefusalError{Why: "publisher path must be absolute, got " + cfg.PublisherPath}
	}
	if !filepath.IsAbs(cfg.PackageDir) {
		return nil, &PublishRefusalError{Why: "package directory must be absolute, got " + cfg.PackageDir}
	}
	if err := validatePublishOrigin("registry origin", cfg.RegistryOrigin, allowLoopback); err != nil {
		return nil, err
	}
	if err := validatePublishOrigin("validator origin", cfg.ValidatorOrigin, allowLoopback); err != nil {
		return nil, err
	}
	for _, arg := range cfg.ExtraArgs {
		for _, forbidden := range forbiddenPublisherFlags {
			if arg == forbidden || strings.HasPrefix(arg, forbidden+"=") {
				return nil, &PublishRefusalError{Why: "publisher flag " + forbidden + " is refused"}
			}
		}
		if arg == "--dry-run" {
			return nil, &PublishRefusalError{Why: "--dry-run is selected by RegistryPublishConfig.DryRun, not by an argument"}
		}
	}
	switch {
	case cfg.DryRun && cfg.Credential != nil:
		return nil, &PublishRefusalError{Why: "a dry-run handler must have no credential provider"}
	case !cfg.DryRun && cfg.Credential == nil:
		return nil, &PublishRefusalError{Why: "a live publish handler requires a credential provider"}
	}
	if !cfg.DryRun {
		if cfg.Approval.ApprovalRef.IsZero() {
			return nil, &PublishRefusalError{Why: "a live publish handler requires an attended approval"}
		}
		if cfg.Approval.RegistryOrigin != cfg.RegistryOrigin {
			return nil, &PublishRefusalError{Why: fmt.Sprintf(
				"approval registry origin %q does not match configured %q",
				cfg.Approval.RegistryOrigin, cfg.RegistryOrigin)}
		}
	}
	return &RegistryPublishHandler{
		cfg:           cfg,
		allowLoopback: allowLoopback,
		bounds: handlerBounds{
			execTimeout: cfg.ExecTimeout, maxOutputBytes: cfg.MaxOutputBytes,
		}.normalized(),
	}, nil
}

// Dispatches is the number of publisher subprocess launches this handler has
// performed. AC7 asserts it stays ZERO across every denial and reaches exactly
// one under a valid grant in the same run; a zero counter with no positive
// control proves only that the counter exists.
func (h *RegistryPublishHandler) Dispatches() int64 { return h.dispatches.Load() }

// CredentialLoads is the number of times the credential provider was consulted.
// A refusal must leave it at zero: the secret is not read to decide whether
// the publish is allowed.
func (h *RegistryPublishHandler) CredentialLoads() int64 { return h.credentialLoad.Load() }

// validatePublishOrigin is the refusal set of Decision 3 applied to one
// origin. It is a pure function so each refusal is independently testable
// without launching anything.
func validatePublishOrigin(what, raw string, allowLoopback bool) error {
	if raw == "" {
		return &PublishRefusalError{Why: what + " is empty"}
	}
	if strings.Contains(raw, "*") {
		return &PublishRefusalError{Why: what + " is a wildcard: " + raw}
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return &PublishRefusalError{Why: what + " is unparseable: " + err.Error()}
	}
	if parsed.User != nil {
		// Never echo raw here: userinfo IS a credential.
		return &PublishRefusalError{Why: what + " embeds credentials in the URL"}
	}
	if parsed.Host == "" {
		return &PublishRefusalError{Why: what + " has no host: " + raw}
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return &PublishRefusalError{Why: what + " carries a query or fragment: " + raw}
	}
	if strings.HasSuffix(raw, "/") {
		return &PublishRefusalError{Why: what + " must not end in a slash: " + raw}
	}
	loopback := isLoopbackHost(parsed.Hostname())
	switch {
	case allowLoopback && !loopback:
		return &PublishRefusalError{Why: fmt.Sprintf(
			"the loopback-only test constructor refuses non-loopback %s %q", what, raw)}
	case allowLoopback:
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return &PublishRefusalError{Why: what + " scheme must be http or https on loopback: " + raw}
		}
	case parsed.Scheme != "https":
		return &PublishRefusalError{Why: what + " must be https for a live origin: " + raw}
	case loopback:
		return &PublishRefusalError{Why: "the production constructor refuses a loopback " + what + ": " + raw}
	}
	return nil
}

func isLoopbackHost(host string) bool {
	switch host {
	case "127.0.0.1", "::1", "localhost":
		return true
	}
	return strings.HasPrefix(host, "127.")
}

// Execute performs one publish attempt.
//
// Order is authority-bearing and every step before the dispatch counter is a
// refusal that costs nothing: shape, cost law, scope law, payload identity,
// recomputed hashes against the attended approval, and only then the
// credential and the subprocess.
func (h *RegistryPublishHandler) Execute(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	if req.Effect != EffectRegistryPublish {
		return nil, &PublishRefusalError{Why: fmt.Sprintf(
			"handler does not implement %q", req.Effect)}
	}
	if req.Cost != PublishCost {
		// The cost is the accounting of one irreversible attempt. A zero-cost
		// request would let an attended budget of one authorize unboundedly
		// many public writes — see MUT-SM-COST-ZERO.
		return nil, &PublishRefusalError{Why: fmt.Sprintf(
			"cost is %d, want exactly %d", req.Cost, PublishCost)}
	}
	id, hashes, err := DecodePublishPayload(payload)
	if err != nil {
		return nil, &PublishRefusalError{Why: err.Error()}
	}
	if id.RegistryOrigin != h.cfg.RegistryOrigin {
		return nil, &PublishRefusalError{Why: fmt.Sprintf(
			"payload registry origin %q does not match the configured %q",
			id.RegistryOrigin, h.cfg.RegistryOrigin)}
	}
	if err := validatePublishOrigin("payload registry origin", id.RegistryOrigin, h.allowLoopback); err != nil {
		return nil, err
	}
	wantScope := PublishScope(id.RegistryOrigin, id.Vendor, id.Name, id.Version)
	if req.Scope != wantScope {
		return nil, &PublishRefusalError{Why: fmt.Sprintf(
			"scope %q does not describe the payload package: want %q", req.Scope, wantScope)}
	}
	if !h.cfg.DryRun {
		if err := h.checkApprovalIdentity(id); err != nil {
			return nil, err
		}
	}
	if err := h.checkManifestIdentity(id); err != nil {
		return nil, err
	}

	// Recompute every hash immediately before dispatch. Nothing carried in the
	// payload or the approval is trusted as a hash — they are COMPARED with
	// digests computed from the bytes on disk at this instant.
	recomputed, err := recomputePublishHashes(h.cfg.PackageDir, h.cfg.Manifest)
	if err != nil {
		return nil, &PublishRefusalError{Why: "recompute package hashes: " + err.Error()}
	}
	if err := comparePublishHashes("payload", recomputed, hashes); err != nil {
		return nil, err
	}
	if !h.cfg.DryRun {
		if err := comparePublishHashes("approval stamp", recomputed, h.cfg.Approval.Hashes); err != nil {
			return nil, err
		}
	}

	var secret []byte
	if !h.cfg.DryRun {
		h.credentialLoad.Add(1)
		secret, err = h.cfg.Credential.Credential()
		if err != nil {
			return nil, fmt.Errorf("broker: load registry credential: %w", err)
		}
	}

	home, err := os.MkdirTemp("", "broker-publish-home-*")
	if err != nil {
		return nil, fmt.Errorf("broker: create empty publish HOME: %w", err)
	}
	defer func() { _ = os.RemoveAll(home) }()

	args := append([]string{"publish"}, h.cfg.ExtraArgs...)
	if h.cfg.DryRun {
		args = append(args, "--dry-run")
	}

	env := h.childEnvironment(home, secret)
	// Fail closed on the one property a dry-run must never lose. This is a
	// runtime assertion on the constructed environment, not a reading of the
	// literal that built it, so it survives an edit to that literal.
	if h.cfg.DryRun && childenv.Has(env, RegistryCredentialVariable) {
		return nil, &PublishRefusalError{
			Why: "dry-run child environment carries " + RegistryCredentialVariable}
	}

	h.dispatches.Add(1)
	out, runErr := runBounded(ctx, h.bounds, handlerCommand{
		path: h.cfg.PublisherPath,
		args: args,
		dir:  h.cfg.PackageDir,
		env:  env,
	})
	redacted := redactSecret(string(publisherOutput(out, runErr)), secret)

	status := classifyPublisherResult(redacted, runErr, id)
	switch status {
	case PublishStatusIndeterminate:
		// The typed ambiguity. No record, no outcome, no result object: the
		// request body may already have reached a registry that cannot be
		// asked to take it back. Session.Invoke leaves the durable intent
		// INDETERMINATE and SM.C reconciles it read-only.
		return nil, &IndeterminateEffectError{
			Effect: req.Effect,
			Scope:  req.Scope,
			Detail: redactSecret(publisherDetail(runErr, redacted), secret),
		}
	case PublishStatusSucceeded:
		return h.resultObject(status, id, recomputed, redacted), nil
	default:
		return nil, &PublishDispatchError{
			Status: status,
			Detail: redactSecret(publisherDetail(runErr, redacted), secret),
		}
	}
}

// PublishDispatchError is a DEFINITE failure: the attempt is over, the
// outcome is known, and Session.Invoke records a resolved failed outcome for
// it exactly as it does for every other handler failure.
type PublishDispatchError struct {
	Status string
	Detail string
}

func (e *PublishDispatchError) Error() string {
	return fmt.Sprintf("broker: Registry.Publish %s: %s", e.Status, e.Detail)
}

func (h *RegistryPublishHandler) checkApprovalIdentity(id PublishIdentity) error {
	approval := h.cfg.Approval
	if id.ApprovalRef != approval.ApprovalRef {
		return &PublishRefusalError{Why: fmt.Sprintf(
			"payload approval ref %s is not the attended approval %s",
			id.ApprovalRef.String(), approval.ApprovalRef.String())}
	}
	if id.Vendor != approval.Vendor || id.Name != approval.Name || id.Version != approval.Version {
		return &PublishRefusalError{Why: fmt.Sprintf(
			"payload publishes %s/%s@%s but the approval stamps %s/%s@%s",
			id.Vendor, id.Name, id.Version, approval.Vendor, approval.Name, approval.Version)}
	}
	return nil
}

func (h *RegistryPublishHandler) checkManifestIdentity(id PublishIdentity) error {
	wantName := id.Vendor + "/" + id.Name
	if h.cfg.Manifest.Package.Name != wantName {
		return &PublishRefusalError{Why: fmt.Sprintf(
			"manifest package %q does not match payload package %q",
			h.cfg.Manifest.Package.Name, wantName)}
	}
	if !sameStrings(h.cfg.Manifest.Exports.Modules, id.Exports) {
		return &PublishRefusalError{Why: fmt.Sprintf(
			"manifest exports %v do not match payload exports %v",
			h.cfg.Manifest.Exports.Modules, id.Exports)}
	}
	if !sameStrings(h.cfg.Manifest.Effects.Max, id.Effects) {
		return &PublishRefusalError{Why: fmt.Sprintf(
			"manifest effects %v do not match payload effects %v",
			h.cfg.Manifest.Effects.Max, id.Effects)}
	}
	return nil
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func recomputePublishHashes(packageDir string, manifest pkgproj.Manifest) (PublishHashes, error) {
	content, err := pkgproj.ContentHash(packageDir)
	if err != nil {
		return PublishHashes{}, err
	}
	tarball, err := pkgproj.CreateTarball(packageDir)
	if err != nil {
		return PublishHashes{}, err
	}
	return PublishHashes{
		TarballSHA256: pkgproj.TarballHash(tarball),
		ContentHash:   content,
		InterfaceHash: pkgproj.InterfaceHash(manifest),
	}, nil
}

// comparePublishHashes names BOTH values on mismatch and names which of the
// three arms diverged, because content/interface and tarball divergence have
// opposite remedies (design doc DD-1).
func comparePublishHashes(what string, recomputed, claimed PublishHashes) error {
	for _, arm := range []struct{ name, got, want string }{
		{"tarball", recomputed.TarballSHA256, claimed.TarballSHA256},
		{"content", recomputed.ContentHash, claimed.ContentHash},
		{"interface", recomputed.InterfaceHash, claimed.InterfaceHash},
	} {
		if arm.got != arm.want {
			return &PublishRefusalError{Why: fmt.Sprintf(
				"recomputed %s hash %s does not match the %s's %s",
				arm.name, arm.got, what, arm.want)}
		}
	}
	return nil
}

// childEnvironment is the ONLY place a publisher child's environment is built.
// It starts from nothing rather than from os.Environ(), so a registry variable
// cannot arrive by inheritance; the credential is added for the live call and
// for no other.
func (h *RegistryPublishHandler) childEnvironment(home string, secret []byte) []string {
	env := []string{
		"PATH=/usr/bin:/bin",
		"HOME=" + home,
		"LANG=C",
		"LC_ALL=C",
		"AILANG_REGISTRY=" + h.cfg.RegistryOrigin,
		"AILANG_REGISTRY_VALIDATOR=" + h.cfg.ValidatorOrigin,
	}
	if !h.cfg.DryRun && len(secret) != 0 {
		env = append(env, RegistryCredentialVariable+"="+string(secret))
	}
	return env
}

// classifyPublisherResult maps the pinned publisher's observable behaviour to
// Decision 3's outcome table.
//
// Every marker below was READ from the pinned v0.30.0 source at e37b370,
// cmd/ailang/pkg_publish.go, not invented:
//
//   - "upload failed: " (line 282) is client.Do's transport error. The request
//     body may already have left this process, so it is the AMBIGUOUS class.
//   - "publish blocked: ", "asset validation failed: ", "no ailang.toml found",
//     "failed to create package tarball" (lines 53, 87, 100, 166, 173) all
//     return BEFORE the POST, so they are definite failed-before-dispatch.
//   - success prints "Published <name>@<version>" (line 129).
//
// A handler timeout is ambiguous for the same reason as a transport error: the
// bound expired while the request was in flight.
func classifyPublisherResult(output string, runErr error, id PublishIdentity) string {
	var timeout *HandlerTimeoutError
	if errors.As(runErr, &timeout) {
		return PublishStatusIndeterminate
	}
	if strings.Contains(output, "upload failed: ") {
		return PublishStatusIndeterminate
	}
	for _, marker := range []string{
		"publish blocked: ",
		"asset validation failed: ",
		"no ailang.toml found",
		"failed to create package tarball",
		"smoke runner failed",
	} {
		if strings.Contains(output, marker) {
			return PublishStatusFailedBeforeDispatch
		}
	}
	if runErr != nil {
		return PublishStatusFailed
	}
	// A clean exit is NOT believed on its own. v0.30.0 prints an explicit
	// success line, so silence means the binary is not the one this handler
	// was written against — which is a failure, never a success.
	if !strings.Contains(output, "Published "+id.Vendor+"/"+id.Name+"@"+id.Version) {
		return PublishStatusFailed
	}
	return PublishStatusSucceeded
}

// publisherOutput is the publisher's observed byte stream.
//
// It exists because runBounded returns a NIL slice alongside a
// *HandlerExitError and carries the bytes inside the error instead — so a
// classifier that reads only the returned slice sees an empty stream for every
// non-zero exit, which is every interesting arm. Measured while building this
// milestone: the ambiguous "upload failed:" transport marker was present in
// HandlerExitError.Output and absent from the returned slice, and the arm
// mis-classified as a definite failure until this function was added.
func publisherOutput(out []byte, runErr error) []byte {
	if len(out) != 0 {
		return out
	}
	var exit *HandlerExitError
	if errors.As(runErr, &exit) {
		return exit.Output
	}
	return out
}

func publisherDetail(runErr error, redacted string) string {
	if runErr == nil {
		return redacted
	}
	var exit *HandlerExitError
	if errors.As(runErr, &exit) {
		// HandlerExitError.Error() embeds the UNREDACTED output. Rebuild the
		// message from the redacted stream instead of formatting the error.
		return fmt.Sprintf("publisher exited with %v (output %q)", exit.Err, redacted)
	}
	return fmt.Sprintf("%v (output %q)", runErr, redacted)
}

func (h *RegistryPublishHandler) resultObject(
	status string,
	id PublishIdentity,
	hashes PublishHashes,
	redactedOutput string,
) []byte {
	payload, err := json.Marshal(publishResultWire{
		PublishResultSchema, status, id.Vendor, id.Name, id.Version,
		h.cfg.RegistryOrigin, h.cfg.ValidatorOrigin,
		hashes.TarballSHA256, hashes.ContentHash, hashes.InterfaceHash,
		redactedOutput,
	})
	if err != nil {
		panic("broker: fixed publish result cannot fail JSON encoding: " + err.Error())
	}
	return payload
}
