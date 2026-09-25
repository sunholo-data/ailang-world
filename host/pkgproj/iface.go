package pkgproj

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

// InterfaceIdentity is the interface section of `ailang pkg quality --json`.
// V1 is the MANIFEST-coverage hash (same bytes as InterfaceHash); V2 is
// upstream's signature-set identity, the one that moves when the exported
// interface moves.
type InterfaceIdentity struct {
	V1         string
	V2         string
	Signatures int
}

var (
	v1Shape = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	v2Shape = regexp.MustCompile(`^sha256:ifacev2:[0-9a-f]{64}$`)
)

// ParseQualityInterface extracts the interface section. It refuses, never
// defaults: an absent hash_v2 (upstream omits it when v2 cannot be built)
// is an error, not an empty identity.
func ParseQualityInterface(out []byte) (InterfaceIdentity, error) {
	var doc struct {
		Schema    string `json:"schema"`
		Interface *struct {
			HashV1     *string `json:"hash_v1"`
			HashV2     *string `json:"hash_v2"`
			Signatures int     `json:"signatures"`
		} `json:"interface"`
	}
	if err := json.Unmarshal(out, &doc); err != nil {
		return InterfaceIdentity{}, fmt.Errorf("pkg quality: not JSON: %w", err)
	}
	if doc.Schema != "ailang.package-quality/v1" {
		return InterfaceIdentity{}, fmt.Errorf("pkg quality: schema %q", doc.Schema)
	}
	if doc.Interface == nil {
		return InterfaceIdentity{}, errors.New("pkg quality: no interface section")
	}
	in := doc.Interface
	if in.HashV2 == nil || !v2Shape.MatchString(*in.HashV2) {
		return InterfaceIdentity{}, errors.New("pkg quality: interface identity v2 absent or malformed")
	}
	if in.HashV1 == nil || !v1Shape.MatchString(*in.HashV1) {
		return InterfaceIdentity{}, errors.New("pkg quality: interface hash v1 absent or malformed")
	}
	return InterfaceIdentity{V1: *in.HashV1, V2: *in.HashV2, Signatures: in.Signatures}, nil
}

// queryInterfaceTimeout bounds one `ailang pkg quality --json --no-run`.
// Measured 0.28-0.85 s per call on the pinned v0.41.0 (go test timings, V17/V22),
// so 30 s is ~35x headroom, and it sits below step 7's outer run_bounded 120,
// so the named QueryTimeoutError fires before the shell's exit 124.
const queryInterfaceTimeout = 30 * time.Second

// queryInterfaceWaitDelay bounds Wait AFTER the deadline kills the direct
// child: a descendant that inherited stdout would otherwise hold the pipe, and
// therefore Wait, open for its own lifetime. It does NOT kill that descendant
// (no process group here; that lifecycle is queue row 24).
const queryInterfaceWaitDelay = 2 * time.Second

// maxQualityOutputBytes caps collected stdout. The real --no-run document is
// 3267 bytes (V24); 1 MiB is ~320x headroom and a hard allocation bound.
const maxQualityOutputBytes = 1 << 20

type queryBounds struct {
	timeout, waitDelay time.Duration
	maxOutput          int
}

var defaultQueryBounds = queryBounds{queryInterfaceTimeout, queryInterfaceWaitDelay, maxQualityOutputBytes}

// QueryTimeoutError is the named result of a quality run that exceeded its bound.
type QueryTimeoutError struct{ Timeout time.Duration }

func (e *QueryTimeoutError) Error() string {
	return fmt.Sprintf("ailang pkg quality exceeded its %s bound", e.Timeout)
}

// ErrQualityOutputOverflow: stdout exceeded maxQualityOutputBytes.
var ErrQualityOutputOverflow = errors.New("ailang pkg quality: stdout or stderr exceeded its output cap")

// cappedBuffer collects at most limit bytes. Crossing the cap records the
// overflow AND calls onOverflow, which cancels the subprocess context: a
// writer error alone does not make cmd.Run return (Run still waits for the
// child), so cancellation SIGKILLs the direct child and WaitDelay releases
// any inherited pipe.
type cappedBuffer struct {
	buf        bytes.Buffer
	limit      int
	over       bool
	onOverflow func()
}

func (c *cappedBuffer) Write(p []byte) (int, error) {
	if c.buf.Len()+len(p) > c.limit {
		if !c.over && c.onOverflow != nil {
			c.onOverflow()
		}
		c.over = true
		return 0, ErrQualityOutputOverflow
	}
	return c.buf.Write(p)
}

// spuriousNoRunGate is the ONE gate `pkg quality --no-run` emits on a package
// whose smoke passes (v0.41.0, V8; upstream report §11 O1). "Not run" is
// reported as "failed". It is the only gate an exit 2 may carry.
const spuriousNoRunGateCode, spuriousNoRunGateMsg = "PUB015", "_smoke.ail failed"

type qualityGate struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// checkQualityGates is the ONE exit-code/gate-payload rule: only exit 0 or 2
// can carry an identity; exit 0 requires gates == []; exit 2 requires gates ==
// exactly the spurious PUB015. Any other combination is refused, naming the
// gate codes.
func checkQualityGates(out []byte, exitCode int) error {
	if exitCode != 0 && exitCode != 2 {
		return fmt.Errorf("pkg quality: exit %d (only 0 or 2 carry an identity)", exitCode)
	}
	var doc struct {
		Gates *[]qualityGate `json:"gates"`
	}
	if err := json.Unmarshal(out, &doc); err != nil {
		return fmt.Errorf("pkg quality: not JSON: %w", err)
	}
	if doc.Gates == nil {
		return errors.New("pkg quality: no gates array")
	}
	gates := *doc.Gates
	codes := make([]string, len(gates))
	for i, g := range gates {
		codes[i] = g.Code
	}
	if exitCode == 0 {
		if len(gates) != 0 {
			return fmt.Errorf("pkg quality: exit 0 with gates %v, want none", codes)
		}
		return nil
	}
	if len(gates) != 1 || gates[0].Code != spuriousNoRunGateCode || gates[0].Msg != spuriousNoRunGateMsg {
		return fmt.Errorf("pkg quality: exit 2 with gates %v; only the spurious %s %q is accepted", codes, spuriousNoRunGateCode, spuriousNoRunGateMsg)
	}
	return nil
}

// QueryInterface runs the pinned binary under a finite deadline.
func QueryInterface(ctx context.Context, packageDir string, manifest Manifest, ailangBin string) (InterfaceIdentity, error) {
	return queryInterface(ctx, packageDir, manifest, ailangBin, defaultQueryBounds)
}

func queryInterface(ctx context.Context, packageDir string, manifest Manifest, ailangBin string, b queryBounds) (InterfaceIdentity, error) {
	runCtx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, ailangBin, "pkg", "quality", "--json", "--no-run", ".")
	cmd.Dir = packageDir
	cmd.Env = childenv.Scrubbed(os.Environ())
	cmd.WaitDelay = b.waitDelay
	stdout := &cappedBuffer{limit: b.maxOutput, onOverflow: cancel}
	stderr := &cappedBuffer{limit: 64 << 10, onOverflow: cancel}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	err := cmd.Run()
	// Overflow first: it is definitive, and a child blocked on a pipe we
	// stopped draining may also run into the deadline.
	if stdout.over || stderr.over {
		return InterfaceIdentity{}, ErrQualityOutputOverflow
	}
	// The CALLER's context is checked first: a caller deadline also marks
	// runCtx DeadlineExceeded, and must not be reported as our bound.
	if cerr := ctx.Err(); cerr != nil {
		return InterfaceIdentity{}, fmt.Errorf("ailang pkg quality: cancelled by caller: %w", cerr)
	}
	if runCtx.Err() == context.DeadlineExceeded {
		return InterfaceIdentity{}, &QueryTimeoutError{Timeout: b.timeout}
	}
	exitCode := 0
	if err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			return InterfaceIdentity{}, fmt.Errorf("ailang pkg quality: %w", err)
		}
		exitCode = ee.ExitCode()
	}
	if gerr := checkQualityGates(stdout.buf.Bytes(), exitCode); gerr != nil {
		return InterfaceIdentity{}, fmt.Errorf("%w: stderr: %s", gerr, bytes.TrimSpace(stderr.buf.Bytes()))
	}
	id, perr := ParseQualityInterface(stdout.buf.Bytes())
	if perr != nil {
		return InterfaceIdentity{}, perr
	}
	if local := InterfaceHash(manifest); id.V1 != local {
		return InterfaceIdentity{}, fmt.Errorf("pkg quality hash_v1 %s disagrees with pkgproj.InterfaceHash %s", id.V1, local)
	}
	return id, nil
}
