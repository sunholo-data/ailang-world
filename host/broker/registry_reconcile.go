package broker

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Probe-then-resolve reconciliation (SM.C).
//
// This file resolves ONE indeterminate Registry.Publish attempt by READING the
// public bucket. It never publishes, never POSTs, and never talks to the
// validator service. Its single network verb is GET.
//
// The asymmetry that shapes every branch here: resolving `not-published` is
// what RE-AUTHORIZES an irreversible POST, so a false absence costs a permanent
// duplicate public artifact, while a false `probe-unavailable` costs a human
// five minutes. Every refusal below is therefore one-way on purpose.
// ---------------------------------------------------------------------------

// Reconciliation receipt states. These are the four outcomes a reconciliation
// pass may report, and only the first three are resolutions: `probe-unavailable`
// is an explicit refusal to decide that requires a human.
const (
	ReconcileSucceededReconciled = "succeeded-reconciled"
	ReconcileConflict            = "conflict"
	ReconcileNotPublished        = "not-published"
	ReconcileProbeUnavailable    = "probe-unavailable"
)

// The three-state per-sample classifier's verdicts. UNINFORMATIVE is not a
// third kind of answer about the registry — it is the statement that this
// sample carries NO information about the registry, because the instrument that
// produced it was not shown to be working in the same pass.
const (
	ProbeAbsent        = "absent"
	ProbePresent       = "present"
	ProbeUninformative = "uninformative"
)

// The default same-pass known-positive control: a package/version MEASURED to
// exist in the read-only bucket (design doc row V-N arm 1; re-measured
// 2026-08-08, 200 with 1289 bytes of well-formed JSON).
const (
	DefaultProbeControlVendor  = "sunholo"
	DefaultProbeControlName    = "auth"
	DefaultProbeControlVersion = "0.4.1"
)

// DefaultAbsentSamplesRequired is how many absent-with-a-firing-control samples
// must be collected before `not-published` may be resolved.
const DefaultAbsentSamplesRequired = 3

// ReconcileConfig pins one reconciliation pass. There is deliberately no field
// that accepts a ready-made URL: every URL this package fetches is built by
// metadataObjectURL from RegistryOrigin, so no caller can point the control at
// a different key-space than the target.
type ReconcileConfig struct {
	// RegistryOrigin is $AILANG_REGISTRY: the READ-ONLY public bucket base
	// (e.g. https://storage.googleapis.com/ailang-registry). It is NOT
	// $AILANG_REGISTRY_VALIDATOR and must never be failed over to it.
	RegistryOrigin string

	// Vendor, Name and Version identify the package version whose existence is
	// in question — the absence TARGET.
	Vendor  string
	Name    string
	Version string

	// Expected are the three digests the durable publish request bound. A
	// present document must match all three or the pass resolves `conflict`.
	Expected PublishHashes

	// ControlVendor/Name/Version select the same-pass known-positive control.
	// Left entirely zero they default to sunholo/auth@0.4.1. The control is
	// always a metadata object under the SAME origin, never an index route.
	ControlVendor  string
	ControlName    string
	ControlVersion string

	// ObservedPublishStatus is the HTTP status the ambiguous publish attempt
	// observed, or 0 if none was seen. It is RECORDED and never DECIDES: a 409
	// is evidence only after the metadata probe, because someone else may have
	// won the immutable version with different bytes.
	ObservedPublishStatus int

	// AbsentSamplesRequired is the bounded absence window (default 3).
	AbsentSamplesRequired int
	// MaxAttempts bounds the TOTAL number of probe passes, informative or not.
	// Uninformative samples do not count toward AbsentSamplesRequired but they
	// do consume an attempt, so a permanently broken instrument terminates at
	// `probe-unavailable` instead of looping forever.
	MaxAttempts int
	// SampleInterval is the delay between attempts.
	SampleInterval time.Duration
	// RequestTimeout bounds one GET.
	RequestTimeout time.Duration
	// MaxBodyBytes bounds one response body. A body over the bound is
	// TRUNCATED, which makes it fail to parse, which makes the sample
	// uninformative — fail-closed by construction.
	MaxBodyBytes int64
	// Client is the HTTP client. Tests inject one; production may leave it nil.
	Client *http.Client
}

// ProbeSample is one probe pass: the control and the target, fetched together,
// and the verdict the classifier reached about the pair.
type ProbeSample struct {
	Verdict       string
	ControlURL    string
	ControlStatus int
	TargetURL     string
	TargetStatus  int
	Reason        string
}

func (s ProbeSample) String() string {
	return fmt.Sprintf("verdict=%s control=%d target=%d reason=%q",
		s.Verdict, s.ControlStatus, s.TargetStatus, s.Reason)
}

// ReconcileReceipt is the verbatim result of one reconciliation pass.
type ReconcileReceipt struct {
	State                 string
	RegistryOrigin        string
	Vendor                string
	Name                  string
	Version               string
	TargetURL             string
	ControlURL            string
	ObservedPublishStatus int
	Samples               []ProbeSample
	AbsentSamples         int
	UninformativeSamples  int
	// Served is the metadata document a PRESENT sample returned, or nil.
	Served *RegistryMetadata
	Detail string
}

func (r ReconcileReceipt) String() string {
	return fmt.Sprintf(
		"state=%s package=%s/%s@%s absent=%d uninformative=%d samples=%d observedPublishStatus=%d detail=%q",
		r.State, r.Vendor, r.Name, r.Version,
		r.AbsentSamples, r.UninformativeSamples, len(r.Samples), r.ObservedPublishStatus, r.Detail)
}

// RegistryMetadata is the subset of the served bucket document reconciliation
// compares. The key names are MEASURED from the live object, not guessed
// (V-N arm 1; re-measured 2026-08-08).
type RegistryMetadata struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	TarballHash   string `json:"tarball_hash"`
	ContentHash   string `json:"content_hash"`
	InterfaceHash string `json:"interface_hash"`
}

// metadataObjectURL builds the ONE URL shape reconciliation is allowed to
// fetch:
//
//	<registryOrigin>/packages/<vendor>/<name>/<version>/metadata.json
//
// STRUCTURAL INVARIANT, and the reason this function exists at all: BOTH the
// absence TARGET and the same-pass known-positive CONTROL are built here, from
// the SAME registryOrigin, differing ONLY in vendor/name/version.
//
// That is not tidiness. MEASURED 2026-08-08: the validator origin
// https://registry.ailang.sunholo.com answers 200 with 35 KB of well-formed
// JSON at /api/packages, and 404 at /packages/{vendor}/{name}/{version}/
// metadata.json. So if the origin is ever misconfigured to the validator AND
// the control is implemented as "fetch the registry index" — the natural choice
// there, and the one a failover or a copy-paste produces — then the control
// returns 200-with-JSON while the target returns 404. The control FIRES, the
// sample is believed ABSENT, the bounded window resolves `not-published`, and
// an irreversible POST is re-authorized. That is precisely the double-publish
// this whole design exists to prevent.
//
// A control that does not travel the target's own key-space proves nothing
// about the target's key-space. TestAbsenceControlSharesTargetKeySpace pins it.
func metadataObjectURL(registryOrigin, vendor, name, version string) (string, error) {
	for _, segment := range []struct{ what, value string }{
		{"vendor", vendor}, {"name", name}, {"version", version},
	} {
		// U1 — an empty segment would collapse the key space.
		if segment.value == "" {
			return "", &PublishRefusalError{Why: "metadata object " + segment.what + " is empty"}
		}
		// U2 — "." and ".." escape the key space while surviving PathEscape.
		if segment.value == "." || segment.value == ".." {
			return "", &PublishRefusalError{Why: fmt.Sprintf(
				"metadata object %s %q is a path traversal segment", segment.what, segment.value)}
		}
		// U3 — anything needing escaping is not a single safe path segment.
		if url.PathEscape(segment.value) != segment.value {
			return "", &PublishRefusalError{Why: fmt.Sprintf(
				"metadata object %s %q is not a single safe path segment", segment.what, segment.value)}
		}
	}
	return registryOrigin + "/packages/" + vendor + "/" + name + "/" + version + "/metadata.json", nil
}

// ReconcileRegistryPublish is the PRODUCTION entry point. It accepts only a
// non-loopback https bucket origin; there is no flag that widens it, exactly as
// with NewRegistryPublishHandler.
func ReconcileRegistryPublish(ctx context.Context, cfg ReconcileConfig) (ReconcileReceipt, error) {
	return reconcileRegistryPublish(ctx, cfg, false)
}

// reconcileLoopback is the TEST-ONLY entry point, unexported for the same
// reason newLoopbackRegistryPublishHandler is: no package outside host/broker
// can aim reconciliation at a non-loopback origin through it.
func reconcileLoopback(ctx context.Context, cfg ReconcileConfig) (ReconcileReceipt, error) {
	return reconcileRegistryPublish(ctx, cfg, true)
}

func reconcileRegistryPublish(
	ctx context.Context,
	cfg ReconcileConfig,
	allowLoopback bool,
) (ReconcileReceipt, error) {
	cfg = cfg.withDefaults()

	// R1 — the origin passes the same refusal set the publish handler applies:
	// https, no wildcard, no userinfo, no query/fragment, no trailing slash,
	// and non-loopback in production.
	if err := validatePublishOrigin("reconcile registry origin", cfg.RegistryOrigin, allowLoopback); err != nil {
		return ReconcileReceipt{}, err
	}
	// R2 — the bucket is not the validator. The design says this must never be
	// pointed at or failed over to the publish service; this is that sentence
	// as a refusal rather than as prose.
	if !allowLoopback && sameOriginHost(cfg.RegistryOrigin, ApprovedValidatorOrigin) {
		return ReconcileReceipt{}, &PublishRefusalError{Why: fmt.Sprintf(
			"reconciliation reads the bucket, not the validator service %s", ApprovedValidatorOrigin)}
	}
	targetURL, err := metadataObjectURL(cfg.RegistryOrigin, cfg.Vendor, cfg.Name, cfg.Version)
	// R3 — an unbuildable target URL is a refusal, never an absence.
	if err != nil {
		return ReconcileReceipt{}, err
	}
	controlURL, err := metadataObjectURL(cfg.RegistryOrigin, cfg.ControlVendor, cfg.ControlName, cfg.ControlVersion)
	// R4 — an unbuildable control URL is a refusal: a pass with no instrument
	// must not run at all rather than run uncontrolled.
	if err != nil {
		return ReconcileReceipt{}, err
	}
	// R5 — the control must not BE the target. If it were, a 404 target would
	// make the control 404 too and every sample would be uninformative; worse,
	// a 200 target would "prove" its own presence.
	if controlURL == targetURL {
		return ReconcileReceipt{}, &PublishRefusalError{Why: fmt.Sprintf(
			"the known-positive control %s is the probe target itself", controlURL)}
	}
	// R6 — an empty expected digest would make a present document compare equal
	// to anything, turning `conflict` into `succeeded-reconciled`.
	if cfg.Expected.TarballSHA256 == "" || cfg.Expected.ContentHash == "" || cfg.Expected.InterfaceHash == "" {
		return ReconcileReceipt{}, &PublishRefusalError{
			Why: "reconciliation requires all three expected digests"}
	}
	// R7 — a window that cannot close is not a bound.
	if cfg.AbsentSamplesRequired < 1 || cfg.MaxAttempts < cfg.AbsentSamplesRequired {
		return ReconcileReceipt{}, &PublishRefusalError{Why: fmt.Sprintf(
			"absence window is unsatisfiable: %d samples required within %d attempts",
			cfg.AbsentSamplesRequired, cfg.MaxAttempts)}
	}

	receipt := ReconcileReceipt{
		State:                 ReconcileProbeUnavailable,
		RegistryOrigin:        cfg.RegistryOrigin,
		Vendor:                cfg.Vendor,
		Name:                  cfg.Name,
		Version:               cfg.Version,
		TargetURL:             targetURL,
		ControlURL:            controlURL,
		ObservedPublishStatus: cfg.ObservedPublishStatus,
	}

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if attempt > 0 && cfg.SampleInterval > 0 {
			select {
			case <-ctx.Done():
				// R8 — a cancelled pass reports the refusal it has reached so
				// far. It never upgrades to a resolution.
				receipt.Detail = "reconciliation cancelled: " + ctx.Err().Error()
				return receipt, nil
			case <-time.After(cfg.SampleInterval):
			}
		}

		// The control and the target are fetched in the SAME pass, in this
		// order, so the instrument is shown to work before its silence is read
		// as an answer.
		control := cfg.fetch(ctx, controlURL)
		target := cfg.fetch(ctx, targetURL)
		verdict, reason := classifyProbeSample(control, target)
		receipt.Samples = append(receipt.Samples, ProbeSample{
			Verdict:       verdict,
			ControlURL:    control.URL,
			ControlStatus: control.Status,
			TargetURL:     target.URL,
			TargetStatus:  target.Status,
			Reason:        reason,
		})

		switch verdict {
		case ProbePresent:
			return resolvePresent(receipt, target.Body, cfg)
		case ProbeAbsent:
			receipt.AbsentSamples++
			if receipt.AbsentSamples >= cfg.AbsentSamplesRequired {
				receipt.State = ReconcileNotPublished
				receipt.Detail = fmt.Sprintf(
					"%d absent samples, each with a firing same-pass control", receipt.AbsentSamples)
				return receipt, nil
			}
		default:
			receipt.UninformativeSamples++
		}
	}

	// R9 — the window closed without enough absent-with-a-firing-control
	// samples. This is the refusal the whole design is built around: an empty
	// result from an instrument that was never shown to work is a claim, not a
	// measurement, and here that claim costs a duplicate immutable publish.
	receipt.State = ReconcileProbeUnavailable
	receipt.Detail = fmt.Sprintf(
		"%d attempts yielded %d absent and %d uninformative samples; %d absent required — human required",
		len(receipt.Samples), receipt.AbsentSamples, receipt.UninformativeSamples, cfg.AbsentSamplesRequired)
	return receipt, nil
}

// resolvePresent decides between `succeeded-reconciled` and `conflict` for a
// document the classifier already accepted as present.
func resolvePresent(receipt ReconcileReceipt, body []byte, cfg ReconcileConfig) (ReconcileReceipt, error) {
	var served RegistryMetadata
	// P1 — the classifier already proved this parses, so a failure here means
	// the two decoders disagree. Refuse rather than resolve.
	if err := json.Unmarshal(body, &served); err != nil {
		receipt.State = ReconcileProbeUnavailable
		receipt.Detail = "present document did not decode: " + err.Error()
		return receipt, nil
	}
	receipt.Served = &served
	wantName := cfg.Vendor + "/" + cfg.Name
	// P2 — a document for a DIFFERENT package or version is not evidence about
	// this one. It is a conflict, never a success.
	if served.Name != wantName || served.Version != cfg.Version {
		receipt.State = ReconcileConflict
		receipt.Detail = fmt.Sprintf(
			"served document identifies %s@%s, want %s@%s", served.Name, served.Version, wantName, cfg.Version)
		return receipt, nil
	}
	// P3 — all three digests, compared by the SAME comparator the publish
	// handler uses, so the two cannot drift apart.
	if err := comparePublishHashes("public metadata", PublishHashes{
		TarballSHA256: served.TarballHash,
		ContentHash:   served.ContentHash,
		InterfaceHash: served.InterfaceHash,
	}, cfg.Expected); err != nil {
		receipt.State = ReconcileConflict
		receipt.Detail = err.Error()
		return receipt, nil
	}
	receipt.State = ReconcileSucceededReconciled
	receipt.Detail = "served metadata matches all three expected digests"
	return receipt, nil
}

// probeResult is one bounded, read-only GET.
type probeResult struct {
	URL    string
	Status int
	Body   []byte
	Err    error
}

// classifyProbeSample is the three-state classifier, and it is REFUSAL-SHAPED:
// six of its eight branches decline to call the sample informative.
//
// The control is examined FIRST and completely. Only an instrument that has
// just been shown to return 200 with well-formed JSON, from the target's own
// origin and key-space, earns the right to have its silence about the target
// read as absence.
func classifyProbeSample(control, target probeResult) (string, string) {
	// C1 — the control never answered.
	if control.Err != nil {
		return ProbeUninformative, "same-pass control fetch failed: " + control.Err.Error()
	}
	// C2 — the control answered, but not 200. Captive portal, 403, wrong
	// origin, bucket permission change: all land here.
	if control.Status != http.StatusOK {
		return ProbeUninformative, fmt.Sprintf(
			"same-pass control returned %d, want 200", control.Status)
	}
	// C3 — 200 is not enough: an interception page is a 200. The control body
	// must be well-formed JSON, which is what the real metadata object serves.
	if !wellFormedJSONObject(control.Body) {
		return ProbeUninformative, "same-pass control body is not a well-formed JSON object"
	}
	// C4 — the instrument works but the target request itself failed.
	if target.Err != nil {
		return ProbeUninformative, "target fetch failed: " + target.Err.Error()
	}
	switch target.Status {
	case http.StatusNotFound:
		// C5 — THE ONLY PATH TO `absent`. Absence is believed only on the
		// MEASURED GCS error document (V-N arm 2), decoded as XML rather than
		// string-matched, so a truncated or garbage body cannot impersonate it.
		if isGCSNoSuchKey(target.Body) {
			return ProbeAbsent, "target 404 with the measured GCS NoSuchKey document and a firing control"
		}
		// C6 — a 404 whose body is not that document is an unexplained 404.
		return ProbeUninformative,
			"target 404 body is not the measured GCS NoSuchKey document"
	case http.StatusOK:
		// C7 — a present document must at least be well-formed JSON before the
		// hash comparison is allowed to run on it.
		if !wellFormedJSONObject(target.Body) {
			return ProbeUninformative, "target 200 body is not a well-formed JSON object"
		}
		return ProbePresent, "target 200 with a well-formed metadata document and a firing control"
	default:
		// C8 — 403, 500, 301: none of these is absence.
		return ProbeUninformative, fmt.Sprintf("target returned %d", target.Status)
	}
}

// isGCSNoSuchKey decodes the measured bucket 404 document. It is a DECODE, not
// a strings.Contains: a truncated body must fail, and MUT-SM-XML-AS-ABSENT's
// garbage arm must differ from its NoSuchKey arm.
func isGCSNoSuchKey(body []byte) bool {
	var doc struct {
		XMLName xml.Name `xml:"Error"`
		Code    string   `xml:"Code"`
	}
	if err := xml.Unmarshal(body, &doc); err != nil {
		return false
	}
	return doc.Code == "NoSuchKey"
}

// wellFormedJSONObject is the control's liveness predicate: a JSON object with
// at least one member. An empty body, an HTML interception page, a bare string
// and `{}` all fail it.
func wellFormedJSONObject(body []byte) bool {
	var members map[string]json.RawMessage
	if err := json.Unmarshal(body, &members); err != nil {
		return false
	}
	return len(members) > 0
}

// fetch performs one bounded, read-only GET. The verb is a constant: this
// package has no code path that can issue anything but GET.
func (cfg ReconcileConfig) fetch(ctx context.Context, target string) probeResult {
	result := probeResult{URL: target}
	reqCtx, cancel := context.WithTimeout(ctx, cfg.RequestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, target, nil)
	if err != nil {
		result.Err = err
		return result
	}
	req.Header.Set("Accept", "application/json")
	resp, err := cfg.client().Do(req)
	if err != nil {
		result.Err = err
		return result
	}
	defer func() { _ = resp.Body.Close() }()
	result.Status = resp.StatusCode
	// A body over the bound is truncated, which makes it fail to parse, which
	// makes the sample uninformative. Fail-closed by construction.
	body, err := io.ReadAll(io.LimitReader(resp.Body, cfg.MaxBodyBytes))
	if err != nil {
		result.Err = err
		return result
	}
	result.Body = body
	return result
}

func (cfg ReconcileConfig) client() *http.Client {
	if cfg.Client != nil {
		return cfg.Client
	}
	return &http.Client{Timeout: cfg.RequestTimeout}
}

func (cfg ReconcileConfig) withDefaults() ReconcileConfig {
	if cfg.ControlVendor == "" && cfg.ControlName == "" && cfg.ControlVersion == "" {
		cfg.ControlVendor = DefaultProbeControlVendor
		cfg.ControlName = DefaultProbeControlName
		cfg.ControlVersion = DefaultProbeControlVersion
	}
	if cfg.AbsentSamplesRequired == 0 {
		cfg.AbsentSamplesRequired = DefaultAbsentSamplesRequired
	}
	if cfg.MaxAttempts == 0 {
		cfg.MaxAttempts = 2*cfg.AbsentSamplesRequired + 2
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = 10 * time.Second
	}
	if cfg.MaxBodyBytes <= 0 {
		cfg.MaxBodyBytes = 1 << 20
	}
	return cfg
}

func sameOriginHost(a, b string) bool {
	parsedA, errA := url.Parse(a)
	parsedB, errB := url.Parse(b)
	if errA != nil || errB != nil {
		return false
	}
	return strings.EqualFold(parsedA.Hostname(), parsedB.Hostname())
}
