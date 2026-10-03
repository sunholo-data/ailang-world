package broker

// Row 135 (design_docs/planned/w-ailang-run-stdin-argv-caps.md §4.3–§4.6):
// ailang-run's stdin, argv and caps. The run's extra capabilities pass three
// gates — the broker grant (Ailang.RunEnv / Ailang.RunNet, or Ailang.Run for
// a Declassify-only run, D-135-5 = A), the operator's allowlist
// (RunCapsConfig, `serve --run-allow-caps / --run-net-allow /
// --run-net-allow-http`) and the AILANG policy layer, which enforces a
// per-episode, per-cap-set policy VARIANT rendered here and verified by
// AILANG's own `policy-tool summary` before its first run. No confinement is
// added in Go: the variant is the whole grant, and AILANG enforces it.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Run capability names (§4.1).
const (
	RunCapIO         = "IO"
	RunCapFS         = "FS"
	RunCapEnv        = "Env"
	RunCapNet        = "Net"
	RunCapDeclassify = "Declassify"
)

// Payload bounds (§4.1), the twins of run.ail's contracted predicates.
const (
	runMaxStdinBytes = 65536
	runMaxArgv       = 32
	runMaxArgBytes   = 1024
)

// runBaseCaps is row 134's default set: a run without caps gets it, and a
// request for exactly it runs under the episode's base policy.
var runBaseCaps = []string{RunCapFS, RunCapIO}

// runKnownCaps is every name caps may hold; runOperatorCaps the ones beyond
// the base set, each of which the operator must enable.
var (
	runKnownCaps    = []string{RunCapDeclassify, RunCapEnv, RunCapFS, RunCapIO, RunCapNet}
	runOperatorCaps = []string{RunCapDeclassify, RunCapEnv, RunCapNet}
)

// The byte limits a trusted_host variant sets explicitly (V38): trusted_host
// leaves them unbounded unless rendered, so they are rendered at exactly the
// restricted-mode defaults measured on v0.52.1.
var trustedHostLimits = []struct {
	key   string
	value int64
}{
	{"max_source_bytes", 1048576},
	{"max_module_graph_bytes", 16777216},
	{"max_output_bytes", 8388608},
	{"max_fs_transfer_bytes", 8388608},
}

// trustedHostGitDeny is the one `.git` entry a trusted_host variant prepends
// (V37): trusted_host drops AILANG's own `.git` protection, and v0.52.1's
// folded matcher makes this single literal-prefix entry cover every case
// variant, the bare `.git` and a worktree's `.git` pointer file.
const trustedHostGitDeny = ".git/**"

// RunCapsConfig is the operator's run allowlist (§4.6). The zero value admits
// only the base IO/FS runs.
type RunCapsConfig struct {
	// Allow names the extra capabilities a run may request: a subset of
	// {Declassify, Env, Net} (`--run-allow-caps`).
	Allow []string
	// NetAllow holds the port-qualified loopback literals a Net run may reach
	// (`--run-net-allow`, repeatable), rendered verbatim as net_allow.
	NetAllow []string
	// NetAllowHTTP renders net_allow_http = true (`--run-net-allow-http`).
	NetAllowHTTP bool
}

func (c RunCapsConfig) allows(name string) bool { return slices.Contains(c.Allow, name) }

// RunCapsConfigError is a startup refusal of the operator's run allowlist.
type RunCapsConfigError struct{ Why string }

func (e *RunCapsConfigError) Error() string { return "run capabilities: " + e.Why }

// Validate is the startup refusal table of §4.3 gate 2 and §4.6: unknown or
// repeated names, Net without a --run-net-allow, --run-net-allow or
// --run-net-allow-http without Net, and any --run-net-allow entry that is
// not a canonical loopback IP literal with a port 1–65535.
func (c RunCapsConfig) Validate() error {
	seen := map[string]bool{}
	for _, name := range c.Allow {
		if !slices.Contains(runOperatorCaps, name) {
			return &RunCapsConfigError{Why: fmt.Sprintf("--run-allow-caps names %q; it admits only Declassify, Env, Net "+
				"(IO and FS are always available)", name)}
		}
		if seen[name] {
			return &RunCapsConfigError{Why: fmt.Sprintf("--run-allow-caps names %q twice", name)}
		}
		seen[name] = true
	}
	switch {
	case c.allows(RunCapNet) && len(c.NetAllow) == 0:
		return &RunCapsConfigError{Why: "--run-allow-caps Net needs at least one --run-net-allow HOST:PORT"}
	case !c.allows(RunCapNet) && len(c.NetAllow) > 0:
		return &RunCapsConfigError{Why: "--run-net-allow is set but --run-allow-caps does not name Net"}
	case !c.allows(RunCapNet) && c.NetAllowHTTP:
		return &RunCapsConfigError{Why: "--run-net-allow-http is set but --run-allow-caps does not name Net"}
	}
	entries := map[string]bool{}
	for _, entry := range c.NetAllow {
		if err := ValidateRunNetAllow(entry); err != nil {
			return err
		}
		if entries[entry] {
			return &RunCapsConfigError{Why: fmt.Sprintf("--run-net-allow names %q twice", entry)}
		}
		entries[entry] = true
	}
	return nil
}

// ValidateRunNetAllow refuses a --run-net-allow entry that is not exactly
// IPLITERAL:PORT for a loopback IP (127.0.0.0/8 or ::1, bracketed), in its
// canonical spelling, with a port 1–65535. This mirrors v0.52.1's load-time
// refusals (V35: a bare host opens every loopback port; a loopback NAME, a
// private, link-local or unspecified address is never admitted) and narrows
// them to loopback, since public-network Net is a later row (§2).
func ValidateRunNetAllow(entry string) error {
	host, port, err := net.SplitHostPort(entry)
	if err != nil {
		return &RunCapsConfigError{Why: fmt.Sprintf("--run-net-allow %q is not HOST:PORT (a bare host would open every "+
			"loopback port; name one port, e.g. 127.0.0.1:7655)", entry)}
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return &RunCapsConfigError{Why: fmt.Sprintf("--run-net-allow %q names a host by name; only a loopback IP literal "+
			"is admitted (127.0.0.1:PORT or [::1]:PORT)", entry)}
	}
	if !ip.IsLoopback() {
		return &RunCapsConfigError{Why: fmt.Sprintf("--run-net-allow %q is not a loopback address; a run's Net reaches "+
			"loopback only", entry)}
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 || strconv.Itoa(n) != port {
		return &RunCapsConfigError{Why: fmt.Sprintf("--run-net-allow %q has an invalid port (want 1-65535)", entry)}
	}
	if canonical := net.JoinHostPort(ip.String(), port); canonical != entry {
		return &RunCapsConfigError{Why: fmt.Sprintf("--run-net-allow %q is not in canonical form; write %q", entry, canonical)}
	}
	return nil
}

// runPolicySpec is one rendered policy's variable part; everything else is
// row 134's key set.
type runPolicySpec struct {
	mode         string   // security_mode
	caps         []string // allowed_caps, in rendered order
	gitDeny      bool     // prepend trustedHostGitDeny to fs_deny_write
	netAllow     []string // net_allow (rendered only when non-empty)
	netAllowHTTP bool
	limits       bool // render trustedHostLimits
	fsBudget     int
}

// baseRunSpec is row 134's episode policy (RenderEpisodePolicy).
var baseRunSpec = runPolicySpec{mode: "restricted", caps: []string{RunCapIO, RunCapFS}, fsBudget: 1000}

func renderPolicy(root string, spec runPolicySpec) ([]byte, error) {
	if !filepath.IsAbs(root) || filepath.Clean(root) != root {
		return nil, fmt.Errorf("broker: episode policy root %q is not a clean absolute path", root)
	}
	for _, r := range root {
		if r == '"' || r == '\\' || r == unicode.ReplacementChar || !unicode.IsPrint(r) {
			return nil, fmt.Errorf("broker: episode policy root %q has a character TOML would need escaped", root)
		}
	}
	deny := episodeDenyWrite
	if spec.gitDeny {
		deny = append([]string{trustedHostGitDeny}, episodeDenyWrite...)
	}
	var b strings.Builder
	b.WriteString(`security_mode = "` + spec.mode + `"` + "\n")
	b.WriteString("allowed_caps = " + tomlStrings(spec.caps) + "\n")
	b.WriteString(`fs_sandbox = "` + root + `"` + "\n")
	fmt.Fprintf(&b, "timeout_ms = %d\n", episodePolicyTimeoutMS)
	b.WriteString("fs_deny_write = " + tomlStrings(deny) + "\n")
	b.WriteString(`entry = "main"` + "\n")
	if len(spec.netAllow) > 0 {
		b.WriteString("net_allow = " + tomlStrings(spec.netAllow) + "\n")
		if spec.netAllowHTTP {
			b.WriteString("net_allow_http = true\n")
		}
	}
	if spec.limits {
		for _, l := range trustedHostLimits {
			fmt.Fprintf(&b, "%s = %d\n", l.key, l.value)
		}
	}
	fmt.Fprintf(&b, "\n[budgets]\nFS = %d\n", spec.fsBudget)
	return []byte(b.String()), nil
}

// tomlStrings renders a TOML array of strings that need no escaping (cap
// names, the fold-variant deny globs, validated IP:PORT literals).
func tomlStrings(xs []string) string {
	quoted := make([]string, len(xs))
	for i, x := range xs {
		quoted[i] = `"` + x + `"`
	}
	return "[" + strings.Join(quoted, ", ") + "]"
}

// runVariantSpec maps a requested cap set (canonical: sorted, distinct,
// known, not Env with Net) onto its policy (§4.4):
//
//	⊆ {IO, FS, Declassify}  restricted; FS always allowed, at budget 0 when
//	                        not requested (the entry-inside-sandbox check and
//	                        the deny list need FS present, V23/V33)
//	with Net                restricted; net_allow = the operator's literals,
//	                        net_allow_http as configured (V34)
//	with Env                trusted_host; one `.git/**` prepended (V37) and
//	                        the four byte limits explicit (V38)
func runVariantSpec(caps []string, op RunCapsConfig) runPolicySpec {
	spec := runPolicySpec{mode: "restricted", fsBudget: 0}
	effective := append([]string(nil), caps...)
	if slices.Contains(caps, RunCapFS) {
		spec.fsBudget = 1000
	} else {
		effective = append(effective, RunCapFS)
	}
	sort.Strings(effective)
	spec.caps = effective
	if slices.Contains(caps, RunCapNet) {
		spec.netAllow = append([]string(nil), op.NetAllow...)
		spec.netAllowHTTP = op.NetAllowHTTP
	}
	if slices.Contains(caps, RunCapEnv) {
		spec.mode = "trusted_host"
		spec.gitDeny = true
		spec.limits = true
	}
	return spec
}

// RenderRunVariant renders the policy a run with exactly caps executes
// under, for the episode worktree root and the operator's allowlist. caps
// must be canonical (validateRunCaps).
func RenderRunVariant(root string, caps []string, op RunCapsConfig) ([]byte, error) {
	if err := validateRunCaps(caps); err != nil {
		return nil, err
	}
	if slices.Equal(caps, runBaseCaps) {
		return RenderEpisodePolicy(root)
	}
	return renderPolicy(root, runVariantSpec(caps, op))
}

// validateRunCaps requires a canonical, non-empty cap set.
func validateRunCaps(caps []string) error {
	if len(caps) == 0 {
		return fmt.Errorf("caps names no capability")
	}
	for i, c := range caps {
		if !slices.Contains(runKnownCaps, c) {
			return fmt.Errorf("unknown capability %q", c)
		}
		if i > 0 && caps[i-1] >= c {
			return fmt.Errorf("caps %v is not sorted and distinct", caps)
		}
	}
	if slices.Contains(caps, RunCapEnv) && slices.Contains(caps, RunCapNet) {
		return fmt.Errorf("caps names both Env and Net; one run takes at most one of them")
	}
	return nil
}

// runEffectFor is the Go twin of run.ail's contracted runEffect.
func runEffectFor(caps []string) string {
	switch {
	case slices.Contains(caps, RunCapEnv):
		return EffectAilangRunEnv
	case slices.Contains(caps, RunCapNet):
		return EffectAilangRunNet
	default:
		return EffectAilangRun
	}
}

// runPolicyInfo is the `policy` block of a run's result (§4.3): which
// verified policy the run executed under.
type runPolicyInfo struct {
	Digest       string   `json:"digest"`
	SecurityMode string   `json:"security_mode"`
	Caps         []string `json:"caps"`
	NetAllow     []string `json:"net_allow"`
}

// runVariant is one rendered, summary-verified policy file.
type runVariant struct {
	path string
	info runPolicyInfo
}

// variantKey names a cap set in a file name: <ep>.run-Declassify.IO.toml.
func variantKey(caps []string) string { return strings.Join(caps, ".") }

// runVariantFor returns the verified policy for caps, rendering and
// verifying it on first use and caching it with the handler (one handler per
// episode, so per episode per cap set, §4.4).
func (h *AilangToolHandler) runVariantFor(ctx context.Context, caps []string) (runVariant, error) {
	if slices.Equal(caps, runBaseCaps) {
		return h.baseVariant, nil
	}
	key := variantKey(caps)
	h.variantMu.Lock()
	defer h.variantMu.Unlock()
	if v, ok := h.variants[key]; ok {
		return v, nil
	}
	spec := runVariantSpec(caps, h.runCaps)
	data, err := renderPolicy(h.root, spec)
	if err != nil {
		return runVariant{}, err
	}
	path := strings.TrimSuffix(h.policyPath, ".toml") + ".run-" + key + ".toml"
	if err := CheckPolicyOutsideRoot(path, h.root); err != nil {
		return runVariant{}, err
	}
	if err := writeFileAtomic(path, data); err != nil {
		return runVariant{}, fmt.Errorf("broker: write run policy variant: %w", err)
	}
	v, err := h.verifyVariant(ctx, path, data, spec)
	if err != nil {
		return runVariant{}, err
	}
	if h.variants == nil {
		h.variants = map[string]runVariant{}
	}
	h.variants[key] = v
	return v, nil
}

// RunVariantUnverifiedError reports a rendered variant whose policy-tool
// summary does not match what World rendered (or that AILANG refuses at
// load, V35): the run never executes under it.
type RunVariantUnverifiedError struct {
	Path string
	Why  string
}

func (e *RunVariantUnverifiedError) Error() string {
	return fmt.Sprintf("broker: run policy variant %s is not verified by policy-tool summary: %s", e.Path, e.Why)
}

// policySummary is the part of `policy-tool summary` World checks (V28).
type policySummary struct {
	SecurityMode string   `json:"security_mode"`
	PolicyDigest string   `json:"policy_digest"`
	FSSandbox    string   `json:"fs_sandbox"`
	Caps         []string `json:"caps"`
	NetAllow     []string `json:"net_allow"`
	CLI          []string `json:"cli"`
}

func (h *AilangToolHandler) summary(ctx context.Context, policyPath string) (policySummary, []byte, error) {
	stdout, err := h.policyToolWith(ctx, policyPath, []byte(`{"op":"summary"}`))
	if err != nil {
		return policySummary{}, stdout, err
	}
	var resp struct {
		OK      bool          `json:"ok"`
		Summary policySummary `json:"summary"`
	}
	if err := json.Unmarshal(stdout, &resp); err != nil || !resp.OK {
		return policySummary{}, stdout, fmt.Errorf("not an ok summary: %q", stdout)
	}
	return resp.Summary, stdout, nil
}

// verifyVariant is §4.3 gate 3's check: AILANG's own summary of the file
// must report the rendered mode, exactly the rendered caps and net_allow,
// the episode root as fs_sandbox, and the digest of the bytes World wrote.
func (h *AilangToolHandler) verifyVariant(ctx context.Context, path string, data []byte, spec runPolicySpec) (runVariant, error) {
	s, _, err := h.summary(ctx, path)
	if err != nil {
		return runVariant{}, &RunVariantUnverifiedError{Path: path, Why: err.Error()}
	}
	sum := sha256.Sum256(data)
	wantCaps := append([]string(nil), spec.caps...)
	sort.Strings(wantCaps)
	gotCaps := append([]string(nil), s.Caps...)
	sort.Strings(gotCaps)
	sandbox, sandboxErr := filepath.EvalSymlinks(s.FSSandbox)
	switch {
	case s.SecurityMode != spec.mode:
		return runVariant{}, &RunVariantUnverifiedError{Path: path, Why: fmt.Sprintf("security_mode %q, rendered %q", s.SecurityMode, spec.mode)}
	case !slices.Equal(gotCaps, wantCaps):
		return runVariant{}, &RunVariantUnverifiedError{Path: path, Why: fmt.Sprintf("caps %v, rendered %v", s.Caps, wantCaps)}
	case !slices.Equal(s.NetAllow, spec.netAllow) && !(len(s.NetAllow) == 0 && len(spec.netAllow) == 0):
		return runVariant{}, &RunVariantUnverifiedError{Path: path, Why: fmt.Sprintf("net_allow %v, rendered %v", s.NetAllow, spec.netAllow)}
	case sandboxErr != nil || sandbox != h.root:
		return runVariant{}, &RunVariantUnverifiedError{Path: path, Why: fmt.Sprintf("fs_sandbox %q is not the workspace root %q", s.FSSandbox, h.root)}
	case s.PolicyDigest != hex.EncodeToString(sum[:]):
		return runVariant{}, &RunVariantUnverifiedError{Path: path, Why: fmt.Sprintf("policy_digest %q is not the rendered bytes' sha256", s.PolicyDigest)}
	}
	netAllow := s.NetAllow
	if netAllow == nil {
		netAllow = []string{}
	}
	return runVariant{path: path, info: runPolicyInfo{Digest: s.PolicyDigest, SecurityMode: s.SecurityMode,
		Caps: gotCaps, NetAllow: netAllow}}, nil
}

// VerifyRunCaps is the startup half of §4.3 gate 3: it renders the base
// policy at cfg.PolicyPath for cfg.Root, then every variant class the
// operator's allowlist can produce (Declassify, Net, Env), and has the tool
// binary's own `policy-tool summary` verify each. cfg.Root is any directory
// (the daemon passes the workspace root); the files are the caller's to
// remove.
func VerifyRunCaps(ctx context.Context, cfg AilangToolConfig) error {
	base, err := RenderEpisodePolicy(cfg.Root)
	if err != nil {
		return err
	}
	if err := writeFileAtomic(cfg.PolicyPath, base); err != nil {
		return err
	}
	h, err := NewAilangToolHandler(ctx, cfg)
	if err != nil {
		return err
	}
	classes := map[string][]string{
		RunCapDeclassify: {RunCapDeclassify, RunCapIO},
		RunCapNet:        {RunCapIO, RunCapNet},
		RunCapEnv:        {RunCapEnv, RunCapFS, RunCapIO},
	}
	for _, name := range runOperatorCaps {
		if !cfg.RunCaps.allows(name) {
			continue
		}
		if _, err := h.runVariantFor(ctx, classes[name]); err != nil {
			return fmt.Errorf("the %s run variant: %w", name, err)
		}
	}
	return nil
}

// writeFileAtomic replaces path with data at mode 0600.
func writeFileAtomic(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".policy-*.tmp")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// runRequest is a decoded, checked run payload.
type runRequest struct {
	path     string
	argsJSON string
	stdin    *string
	argv     []string
	caps     []string // canonical; runBaseCaps when absent
}

var runPayloadKeys = []string{"path", "args_json", "stdin", "argv", "caps"}

// parseRunPayload is the handler's defence-in-depth twin of run.ail's plan
// checks (§4.5): exact keys (a case-aliased key is named), every bound, the
// cap names, and the effect the caps select must be the effect requested.
func parseRunPayload(effect string, payload []byte) (runRequest, error) {
	refuse := func(why string) (runRequest, error) {
		return runRequest{}, &AilangToolRefusalError{Effect: effect, Why: why}
	}
	fields, err := decodePayloadObject(payload)
	if err != nil {
		return refuse(err.Error())
	}
	req := runRequest{}
	capsGiven := false
	for key, value := range fields {
		switch key {
		case "path":
			if json.Unmarshal(value, &req.path) != nil || req.path == "" {
				return refuse("path must be a non-empty string")
			}
		case "args_json":
			if json.Unmarshal(value, &req.argsJSON) != nil || !json.Valid([]byte(req.argsJSON)) {
				return refuse("args_json must be a string holding JSON")
			}
		case "stdin":
			var s string
			if json.Unmarshal(value, &s) != nil {
				return refuse("stdin must be a string")
			}
			if len(s) > runMaxStdinBytes {
				return refuse(fmt.Sprintf("stdin is %d bytes; the limit is %d", len(s), runMaxStdinBytes))
			}
			req.stdin = &s
		case "argv":
			var argv []string
			if json.Unmarshal(value, &argv) != nil || argv == nil {
				return refuse("argv must be an array of strings")
			}
			if len(argv) > runMaxArgv {
				return refuse(fmt.Sprintf("argv has %d items; the limit is %d", len(argv), runMaxArgv))
			}
			for _, a := range argv {
				if len(a) > runMaxArgBytes {
					return refuse(fmt.Sprintf("an argv item is %d bytes; the limit is %d", len(a), runMaxArgBytes))
				}
				if strings.ContainsRune(a, 0) {
					return refuse("an argv item contains a NUL byte")
				}
			}
			req.argv = argv
		case "caps":
			var caps []string
			if json.Unmarshal(value, &caps) != nil || caps == nil {
				return refuse("caps must be an array of capability names")
			}
			sorted := append([]string(nil), caps...)
			sort.Strings(sorted)
			if err := validateRunCaps(sorted); err != nil {
				return refuse(err.Error())
			}
			req.caps, capsGiven = sorted, true
		default:
			for _, name := range runPayloadKeys {
				if strings.EqualFold(key, name) {
					return refuse(fmt.Sprintf("key %q aliases %q", key, name))
				}
			}
			return refuse(fmt.Sprintf("unknown payload key %q", key))
		}
	}
	if req.path == "" {
		return refuse("path is required")
	}
	if !capsGiven {
		req.caps = append([]string(nil), runBaseCaps...)
	}
	if want := runEffectFor(req.caps); want != effect {
		return refuse(fmt.Sprintf("caps %v run as %s, not %s", req.caps, want, effect))
	}
	return req, nil
}

// operatorRefusal is §4.3 gate 2 inside the handler: "" when every cap
// beyond IO/FS is enabled by the operator (and Net has somewhere to go).
func (h *AilangToolHandler) operatorRefusal(caps []string) string {
	for _, c := range caps {
		if c == RunCapIO || c == RunCapFS {
			continue
		}
		if !h.runCaps.allows(c) {
			return fmt.Sprintf("capability %q is not enabled by the operator (serve --run-allow-caps)", c)
		}
	}
	if slices.Contains(caps, RunCapNet) && len(h.runCaps.NetAllow) == 0 {
		return "Net names no reachable host:port (serve --run-net-allow)"
	}
	return ""
}

// runArgv is the exact command line: `run --policy P [--args-json J] --
// <path> [-- <argv…>]`. The first `--` makes a hyphen-leading path a file
// (V58); the inner `--` keeps every argv item away from the flag parser
// (V13, V39). No widening flag is ever passed (V41).
func runArgv(policyPath string, req runRequest) []string {
	args := []string{"run", "--policy", policyPath}
	if req.argsJSON != "" {
		args = append(args, "--args-json", req.argsJSON)
	}
	args = append(args, "--", req.path)
	if len(req.argv) > 0 {
		args = append(args, "--")
		args = append(args, req.argv...)
	}
	return args
}
