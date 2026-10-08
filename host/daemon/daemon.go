// Package daemon implements the `ailang-worldd` local daemon's HTTP transport
// shell (w-worldd-m2, Decisions 1, 4, 5 and 7).
//
// The daemon contains NO semantics. Every request terminates in an existing M1
// host call (host/store, host/archive, host/registry, host/hashref) or in a
// pure predicate already frozen and Z3-proven in AILANG
// (design_docs/sketches/worlddapi.ail). Decision 1's S3 answer: this is the
// OS-process boundary that future AILANG packages are served *through*, which
// is precisely the thing that cannot itself be a package.
//
// Three structural properties are load-bearing and are each proved by a test in
// daemon_test.go rather than asserted in prose:
//
//   - LOCAL-FIRST IS STRUCTURAL (Decision 4). Startup validates the bind host
//     with isLoopbackHost, a byte-for-byte mirror of the Z3-proven predicate in
//     the sketch, and REFUSES to start otherwise. M2 ships no override flag.
//
//   - EVERY WAIT AND ALLOCATION IS BOUNDED (Decision 7). The named constant
//     block below is the single source of the daemon's bounds: all four
//     http.Server timeouts are set at construction (a zero value is a test
//     failure, not a default), the commit body cap mirrors the Z3-proven
//     withinCommitBytes bound, the client deadline bounds every CLI call, and
//     the shutdown drain is deadline-bounded then hard-closed — never
//     unbounded.
//
//   - WRITER AUTHORITY IS FAIL-CLOSED (Decision 2, ratified arm A). The daemon
//     takes its write handle through store.Open, which acquires the non-waiting
//     cross-process writer lock. A second serve on the same database fails
//     immediately with *store.WriterAlreadyActive, surfaced here as a fatal
//     structured StartupError.
//
// Scope boundary: M2.B serves exactly the frozen worldd-native /v1/* route
// table (see the sketch's routes()). CLI client verbs beyond health/head belong
// to M2.C. There is no effect broker, no capability/budget authority and no
// MCP/A2A projection here — those are clause-3 and clause-6 respectively.
package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sunholo-data/ailang/serveapi/protocol"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/procbound"
	"github.com/sunholo-data/ailang-world/host/projection"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Version is the daemon's own release string, reported by GET /v1/health. It
// versions the transport shell, not the world semantics (those are versioned by
// the epoch registry and by the archived interpreter's HashRef).
const Version = "0.1.0"

// ---------------------------------------------------------------------------
// Decision 7 — Bounded Waits & Allocations (Standing Rule 6)
// ---------------------------------------------------------------------------
//
// design_docs/planned/w-worldd-m2.md, "Decision 7 — Bounded Waits &
// Allocations": loopback exposure is a trust statement, not a resource-safety
// statement. A trusted operator can still wedge an unbounded server with one
// hung connection or one giant payload, and the daemon is designed to run
// unattended for days. Therefore EVERY wait and EVERY request-driven allocation
// in the daemon is bounded by one of the named constants below — no unbounded
// http.Server, no unbounded client call, no unbounded drain, no uncapped body
// read.
//
// TestBoundedWaitsAndBodyLimit is the non-vacuity gate on this block: it asserts
// the constructed server's four timeout fields EQUAL these constants and are
// non-zero, and that maxCommitBytes equals the Z3-proven sketch bound.
const (
	// readHeaderTimeout bounds http.Server.ReadHeaderTimeout (D7 table: 5 s) —
	// slow-header connections cannot hold a server goroutine indefinitely.
	readHeaderTimeout = 5 * time.Second

	// readTimeout bounds http.Server.ReadTimeout (D7 table: 30 s) — the full
	// request read, including POST /v1/commit bodies (M2.B).
	readTimeout = 30 * time.Second

	// writeTimeout bounds http.Server.WriteTimeout (D7 table: 30 s) — the
	// response write, including the ?payload=true object reads (M2.B).
	writeTimeout   = 30 * time.Second
	invokeDeadline = 20 * time.Second

	// idleTimeout bounds http.Server.IdleTimeout (D7 table: 120 s) — the
	// keep-alive connection lifetime.
	idleTimeout = 120 * time.Second

	// maxCommitBytes bounds http.MaxBytesReader on every body-reading route
	// (D7 table: 8 MiB; in v1 that is POST /v1/commit alone, added in M2.B;
	// oversized -> 413 PayloadTooLarge).
	//
	// This bound is NOT a tunable Go constant: it is frozen semantically by the
	// Z3-proven withinCommitBytes in design_docs/sketches/worlddapi.ail —
	//
	//	ensures { result == (n >= 0 && n <= 8388608) }
	//
	// verified on the pinned AILANG v0.30.0. TestBoundedWaitsAndBodyLimit part
	// (c) asserts the equality so the Go constant and the proven bound cannot
	// drift silently. Raising it is a doc + sketch change, not a tweak.
	maxCommitBytes = 8388608

	// defaultClientTimeout bounds context.WithTimeout on EVERY CLI REST client
	// call (D7 table: 30 s), so no client call can hang past the deadline.
	// Exported to cmd/ailang-worldd as DefaultClientTimeout.
	defaultClientTimeout = 30 * time.Second

	// shutdownTimeout bounds http.Server.Shutdown (D7 table: 10 s). On expiry
	// the daemon hard-Closes and reports a non-nil error: an incomplete drain is
	// REPORTED, never waited out forever.
	shutdownTimeout = 10 * time.Second

	// readDeadline bounds the ELAPSED TIME of every store read a GET handler
	// performs (D7 addendum, w-daemon-read-cancellation). The four http.Server
	// timeouts above bound the transport; none of them bounds the wait that
	// happens BELOW the transport, inside database/sql. This constant does,
	// EXCEPT for a read blocked on a SQLite lock: SQLite's busy-retry sleep does
	// not stop on the driver's interrupt, so busy_timeout bounds that wait
	// instead. New only validates that the store's CONFIGURED busy_timeout is
	// below this value (checkReadOrdering) — configuration validation, not
	// deadline enforcement (w-daemon-lock-wait-not-deadline-bound).
	//
	// It must stay well below writeTimeout: the 503 must be writable inside the
	// connection's remaining write window. At 10 s against a 30 s writeTimeout,
	// a read that consumes the whole deadline still leaves 20 s to write the
	// error. TestBoundedWaitsAndBodyLimit pins the literal.
	readDeadline = 10 * time.Second

	// commitBudget bounds POST /v1/commit's store work: connection acquisition
	// and every statement up to the cancellation cutoff (row 23 bound table B6,
	// ratified D-WORLD-40: busy_timeout + 1 s). Expiry before the cutoff is 503
	// Timeout (not committed); after it, 503 CommitUncertain (reconcile).
	commitBudget = 3 * time.Second

	// credentialBudget bounds the session middleware's credential lookup
	// (row 23 bound table B4, ratified D-WORLD-40). Expiry is 503 Timeout,
	// never 401: a timeout is not an authentication failure.
	credentialBudget = 3 * time.Second
)

// Operational defaults (Decision 4 / Decision 5). Exported because
// cmd/ailang-worldd's flag defaults must be the same values the daemon
// validates against.
const (
	// DefaultBindHost is the default listen host: loopback, per Decision 4.
	DefaultBindHost = "127.0.0.1"
	// DefaultBindPort is the daemon's operational default port.
	DefaultBindPort = 7644
	// DefaultBind is the default `serve --bind` value.
	DefaultBind = "127.0.0.1:7644"
	// DefaultAddr is the default `--addr` client base URL (Decision 5).
	DefaultAddr = "http://127.0.0.1:7644"
	// DefaultClientTimeout is the exported view of the D7 client deadline for
	// the CLI client in cmd/ailang-worldd. Kept as one constant so the CLI
	// cannot invent a second, unbounded deadline.
	DefaultClientTimeout = defaultClientTimeout
)

// ListenAnnouncePrefix is the STABLE stdout prefix of the listen announcement
// written by Run once the socket is bound. It is a committed interface: M2.C's
// end-to-end test starts the daemon with `--bind 127.0.0.1:0` and recovers the
// kernel-assigned port by trimming this prefix. Changing it breaks that test.
const ListenAnnouncePrefix = "ailang-worldd listening on "

const (
	integrityScanPageSize   = 64
	integrityScanRowBudget  = 20000
	integrityScanTimeBudget = 2 * time.Second
)

// unpinnedRelease is the epoch-registry candidate recorded when serve is started
// WITHOUT --ailang-bin, i.e. when there is no archived interpreter to nominate.
//
// Honest consequence, stated rather than hidden: epoch-1 revisions are
// content-addressed, so a database bootstrapped with an archived interpreter and
// one bootstrapped unpinned are genuinely DIFFERENT revisions. Starting the same
// database both ways is a real registry divergence and is reported as a fatal
// structured StartupError by registry.Bootstrap — never silently rewritten.
const unpinnedRelease = "unpinned"

// isLoopbackHost reports whether host is a loopback bind target.
//
// This mirrors EXACTLY the frozen, Z3-proven predicate of the same name in
// design_docs/sketches/worlddapi.ail:
//
//	ensures { result == (host == "127.0.0.1" || host == "::1" || host == "localhost") }
//
// The sketch's own test vectors ("127.0.0.1" -> true, "::1" -> true,
// "0.0.0.0" -> false, "example.com" -> false) are replayed in
// TestIsLoopbackHostMirrorsSketchPredicate, so the Go mirror and the proven
// predicate cannot drift. Decision 4: a host failing this predicate REFUSES
// startup and M2 ships no override flag, which makes local-first structural
// rather than advisory.
func isLoopbackHost(host string) bool {
	return host == "127.0.0.1" || host == "::1" || host == "localhost"
}

// Config is the resolved `ailang-worldd serve` configuration. It mirrors the
// sketch's DaemonConfig plus the optional interpreter to archive at startup.
type Config struct {
	// DBPath is the world store database. The daemon takes SOLE writer
	// authority over it for the lifetime of the process (Decision 2).
	DBPath string
	// BindHost must satisfy isLoopbackHost (Decision 4).
	BindHost string
	// BindPort is the TCP port; 0 asks the kernel for an ephemeral port, whose
	// resolved value is announced on stdout (see ListenAnnouncePrefix).
	BindPort int
	// AilangBin, when non-empty, is the interpreter archived at startup
	// (Decision 6 pinning) and reported by GET /v1/health.
	AilangBin string
	// WorkspaceRoot (`--workspace-root`, row 134 D-SE-4 A) is the directory
	// whose child <root>/<episode> is that episode's worktree. It must be a
	// directory that contains none of the daemon's state (store, archive,
	// rendered policies, tool cache) — startup refuses otherwise (AC4.5).
	WorkspaceRoot string
	// ToolAilangBin (`--tool-ailang-bin`, D-SE-1 A) is the AILANG binary the
	// workspace tools run. It is archived and hash-verified like AilangBin
	// (independent of the compiler pin) and must be ToolBinaryRelease. The
	// workspace tools are served only when WorkspaceRoot is also set.
	ToolAilangBin string
	// ExamplesDir (`--examples-dir`, row 134 break-3 fix) is the AILANG examples
	// corpus examples-search reads, passed to the tool as AILANG_EXAMPLES. The
	// corpus was not built into the v0.51.0 tool binary (V65); v0.52.1 embeds
	// one, but World serves only the operator's corpus. It must be a directory
	// outside WorkspaceRoot — startup refuses otherwise. Empty means no corpus:
	// examples-search answers NoExamplesCorpusRefusal. (The CLI defaults it to
	// the operator's ~/.ailang/examples when that exists.)
	ExamplesDir string
	// RunAllowCaps, RunNetAllow and RunNetAllowHTTP (`--run-allow-caps`,
	// `--run-net-allow`, `--run-net-allow-http`; row 135 §4.3/§4.6) are the
	// operator's ailang-run allowlist: the extra capabilities (Declassify,
	// Env, Net) a run may request, and the port-qualified loopback literals a
	// Net run may reach. Empty admits only row 134's IO/FS runs. Startup
	// refuses an invalid allowlist, one set without the workspace tools, and
	// one whose rendered variant AILANG's own policy summary refuses.
	RunAllowCaps    []string
	RunNetAllow     []string
	RunNetAllowHTTP bool
	// ExecProfiles, ExecEpisodeProjects, ExecSandbox, ExecNode and
	// ExecMaxOutputBytes are the `--exec-*` flags (row 140 §4.6): operator
	// exec profiles, the episode->project map, the pinned srt install, the
	// node that runs it ("" = node on PATH) and the per-stream output kill
	// (0 = 64 MiB). Startup refuses an exec flag without the workspace tools
	// and any §4.3/§4.6 violation, and runs the startup probe for every
	// profile; with no profile, workspace-exec refuses every call.
	ExecProfiles        []string
	ExecEpisodeProjects []string
	ExecSandbox         string
	ExecNode            string
	ExecMaxOutputBytes  int64
	// WorkspaceModuleRoot and WorkspaceEpisodeModuleRoots
	// (`--workspace-module-root REL`, `--workspace-episode-module-root EP=REL`;
	// row 141 M1) make <workspace-root>/<episode>/REL the AILANG sandbox and
	// module root instead of the worktree itself. REL is clean and relative
	// (the exec profile's path grammar); the episode override beats the
	// default. Both need WorkspaceRoot; a REL that is not an existing real
	// directory refuses that episode's tools (R8), never creates it.
	WorkspaceModuleRoot         string
	WorkspaceEpisodeModuleRoots []string
	// WorkspacePackageCache (`--workspace-package-cache DIR`; row 141 M2) is
	// a read-only operator snapshot of registry packages, linked as every
	// episode's package cache. It must lie outside WorkspaceRoot and the state
	// directory and hold no symlink and nothing writable; startup refuses
	// otherwise, and as uid 0. Needs WorkspaceRoot.
	WorkspacePackageCache string
	// ErrorLog receives the operator-facing detail of every sanitized 500: one
	// line per error, carrying the route and the VERBATIM store error that the
	// response body no longer echoes (Decision: sanitize-vs-expose).
	//
	// nil means os.Stderr — resolved ONCE in New, never at each write, so the
	// wiring is a constructed field a test can both assert on and replace. The
	// party entitled to the daemon's internals is the process owner, who owns
	// stderr; a loopback HTTP client is not that party.
	//
	// It is deliberately NOT the announce writer. ListenAnnouncePrefix is a
	// one-line protocol whose consumers read exactly one line, and extra lines
	// on that stream were measured deadlocking Run against an io.Pipe.
	ErrorLog io.Writer

	// templateRebuildGate is a test seam (row 153 AC3.7): when non-nil, the
	// background template rebuild awaits a receive on it before each rebuild.
	// Production leaves it nil.
	templateRebuildGate <-chan struct{}
}

// Startup stages, used as the Stage field of StartupError so an operator (and a
// test) can tell WHICH lifecycle step refused rather than reading prose.
const (
	StageBindPolicy = "bind-policy"
	StageConfig     = "config"
	StageStoreOpen  = "store-open"
	StageArchive    = "archive"
	StageRegistry   = "registry-bootstrap"
	StageListen     = "listen"
)

// ErrNonLoopbackBind is the named sentinel for a refused non-loopback bind. It
// is what makes Decision 4's refusal assertable (errors.Is) instead of a string
// match on a message.
var ErrNonLoopbackBind = errors.New("daemon: bind host is not loopback")

// ErrUnorderedTimeouts is the named sentinel for a store whose CONFIGURED SQLite
// lock-retry window (busy_timeout) is not numerically below the CONFIGURED read
// deadline. This is configuration validation only: the deadline does not govern
// a read blocked on a lock, and a window below the deadline does not guarantee
// such a read completes before it (SQLite's retry granularity can exceed the
// gap). Real deadline enforcement is deferred.
var ErrUnorderedTimeouts = errors.New("daemon: store busy_timeout is not below the read deadline")

// checkReadOrdering refuses a configured lock-retry window at or above the
// configured read deadline. A window <= 0 disables SQLite's busy handler (a
// lock conflict fails immediately) and is accepted.
func checkReadOrdering(window, deadline time.Duration) error {
	if window >= deadline {
		return fmt.Errorf("%w: busy_timeout %s must be below read deadline %s", ErrUnorderedTimeouts, window, deadline)
	}
	return nil
}

// StartupError is the structured fatal error of the serve lifecycle. Every
// startup refusal is one of these: nothing in the lifecycle degrades silently,
// which is the explicit requirement for a divergent registry head.
type StartupError struct {
	// Stage is one of the Stage* constants.
	Stage string
	// Detail is the operator-facing explanation.
	Detail string
	// Err is the wrapped cause, if any.
	Err error
	// Settled closes after retained cleanup releases writer authority.
	Settled <-chan struct{}
}

func (e *StartupError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("daemon startup failed at %s: %s: %v", e.Stage, e.Detail, e.Err)
	}
	return fmt.Sprintf("daemon startup failed at %s: %s", e.Stage, e.Detail)
}

func (e *StartupError) Unwrap() error { return e.Err }

// Daemon is a constructed, not-yet-listening daemon: the store write handle is
// held, the interpreter (if any) is archived, the epoch registry is bootstrapped
// and the HTTP server is built with its D7 timeouts. Listen/Serve/Shutdown/Close
// drive the rest of the lifecycle.
type Daemon struct {
	cfg       Config
	store     *store.Store
	bootstrap func(context.Context, *store.Store, string) (registry.Registry, hashref.HashRef, error)
	srv       *http.Server
	ln        net.Listener

	// resolver maps an inbound Authorization: Bearer session credential to a
	// binding (w-session-authority D5/D6). New wires it to authority.New over
	// the SAME store the daemon serves, so the session_credentials table and
	// the world tables live in one database. It is consulted by the session
	// middleware before any protected handler runs.
	resolver authority.Resolver

	// reads is the read seam every /v1 GET route goes through. New always wires
	// it to the SAME *store.Store held in `store`, so the production path is
	// unchanged and passes through one extra interface dispatch — the seam
	// exists so a test can WRAP (embed-and-override) the real store to produce
	// a stimulus a real store cannot produce on demand: a getter that blocks.
	//
	// It is deliberately NOT the whole store. Writes (Commit) stay on `store`:
	// putting the write path behind a type named "reads" would widen the seam
	// past this item's scope. The integrity sweep's ScanUnreadableLog /
	// ScanUnreadableWorlds also stay on `store` — they are startup reads, not
	// request reads, and no /v1 route reaches them.
	reads readStore

	// drainTimeout bounds the graceful shutdown. New always sets it to the D7
	// shutdownTimeout; it is a field rather than a direct constant reference so
	// the bound is (a) assertable as wiring and (b) shrinkable in tests, which
	// is the only way the expiry branch of the SHIPPED Shutdown path can be
	// exercised without a ten-second test.
	drainTimeout time.Duration

	// readDeadline bounds every store read a GET handler performs. New always
	// sets it from the readDeadline constant; it is a field rather than a
	// direct constant reference for the same two reasons drainTimeout is —
	// the wiring is assertable (MU8: wire it to 0 and TestBoundedWaitsAndBodyLimit
	// reds) and the value is shrinkable in tests, which is the only way the
	// SHIPPED timeout branch can be exercised without a ten-second test.
	readDeadline time.Duration

	// commits is the durable write seam POST /v1/commit goes through (row 23
	// M2). New wires it to the SAME *store.Store as `store`; a test wraps it
	// to produce outcomes a real store cannot produce on demand (a held
	// connection, a post-cutoff expiry). commitBudget is New's commitBudget,
	// a field so the timeout branches are testable without a 3 s test.
	commits      durableStore
	commitBudget time.Duration

	// credentialBudget is New's credentialBudget (B4); a field for the same
	// reason commitBudget is.
	credentialBudget time.Duration

	// errLog is the RESOLVED destination of every sanitized 500's detail line.
	// New resolves Config.ErrorLog's nil to os.Stderr here, so this field is
	// never nil on a constructed daemon and writeInternalError needs no nil check
	// on a path that only runs when something has already gone wrong.
	errLog io.Writer

	// capsule is the runner behind d.coord (nil without --ailang-bin); kept so
	// the daemon can report its cold-run count (row 153).
	capsule *capsule.Runner

	// arch is the archive behind the capsule (nil without --ailang-bin) and
	// interpreterHash the pinned interpreter's ref; templates is the row-153
	// per-descriptor compile-cache status (templates.go).
	arch            *archive.Archive
	interpreterHash hashref.HashRef
	templates       templateBook

	scanPageSize   int
	scanRowBudget  int
	scanTimeBudget time.Duration
	integrity      IntegrityReport

	// projection serves the two additive A2A surface routes
	// (GET /.well-known/agent.json, POST /a2a/; w-a2a-session-projection
	// P6.B-A2A-CARD) over THIS daemon's one resolver, one store handle and one
	// read seam — no second store, credential path or policy engine. The
	// routes are NOT in isProtected: the projection resolves the session
	// itself so /a2a/ can answer in JSON-RPC form.
	projection *projection.Handler

	// coord is the invocation coordinator (nil without --ailang-bin); its
	// Binder is d.binder, over workspace.registry(episode).
	coord *coordinator.Coordinator

	// workspace is the resolved --workspace-root/--tool-ailang-bin pair
	// (row 134 §4.3). Its registry is empty unless both are set.
	workspace *workspaceTools

	// Per-instance stream bounds and lifetime signals.
	aguiBudget    time.Duration
	aguiBodyBound time.Duration
	aguiTick      time.Duration
	aguiStop      chan struct{}
	aguiStopOnce  sync.Once
	aguiSlots     chan struct{}

	// Health facts resolved once at startup and served verbatim.
	interpreterRef     string
	interpreterVersion string
}

// readStore is the daemon's request-read surface: EXACTLY the six store
// getters the seven /v1 GET routes reach, each context-first.
//
// SIX methods. The route count and the method count are different
// numbers and the design doc conflates them in two places: GetLogEntry serves
// BOTH GET /v1/log/{index} and the bounded loop of GET /v1/log, so seven routes
// reach the store through six distinct getters. Commit and GetVerifyResult
// are outside the daemon read path.
//
// *store.Store satisfies this by construction — New assigns the same handle to
// both d.store and d.reads — so the seam adds no production behaviour. Its only
// job is to let a test wrap the real store rather than replace it, which keeps
// every handler-side mutation non-vacuous: the fake substitutes store BODIES,
// and all handler code still runs the production path.
type readStore interface {
	GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error)
	GetWorld(ctx context.Context, ref hashref.HashRef) (store.World, bool, error)
	GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error)
	LogEntriesAfter(ctx context.Context, after int64, limit int) ([]store.LogEntry, error)
	LogEntriesLatest(ctx context.Context, limit int) ([]store.LogEntry, error)
	GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error)
	SelectedHead(ctx context.Context) (hashref.HashRef, bool, error)
	ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]store.Object, error)
	ObjectReferences(ctx context.Context, ref hashref.HashRef, after *store.ObjectReferenceCursor, limit int) ([]store.ObjectReference, error)
	ObjectCommits(ctx context.Context, ref hashref.HashRef, afterEntry int64, limit int) ([]int64, error)
	GetReceipt(ctx context.Context, id string) (store.Receipt, bool, error)
}

// durableStore is the daemon's durable-write surface: the context-bounded
// commit of row 23's M1. *store.Store satisfies it by construction.
type durableStore interface {
	Commit(ctx context.Context, c store.Commit) error
	AppendIntent(ctx context.Context, id string, intent store.JournalIntent) (int64, hashref.HashRef, error)
}

// IntegrityReport is the bounded startup sweep result.
type IntegrityReport struct {
	LogRowsScanned   int
	WorldRowsScanned int
	Holes            []store.UnreadableRow
	Complete         bool
	ResumeLogIndex   int64
	ResumeWorldRef   string
}

func (d *Daemon) IntegrityReport() IntegrityReport {
	r := d.integrity
	r.Holes = append([]store.UnreadableRow(nil), r.Holes...)
	return r
}

func (r IntegrityReport) Lines() []string {
	lines := make([]string, 0, len(r.Holes)+1)
	for _, hole := range r.Holes {
		if hole.Table == "log_entries" {
			lines = append(lines, fmt.Sprintf("integrity_hole table=log_entries index=%d field=%s",
				hole.Index, hole.Field))
		} else {
			lines = append(lines, fmt.Sprintf("integrity_hole table=worlds ref=%s field=%s",
				hole.Ref, hole.Field))
		}
	}
	if r.Complete {
		lines = append(lines, fmt.Sprintf("integrity_scan_complete log_rows=%d world_rows=%d holes=%d",
			r.LogRowsScanned, r.WorldRowsScanned, len(r.Holes)))
	} else {
		lines = append(lines, fmt.Sprintf("integrity_scan_incomplete log_rows=%d world_rows=%d holes=%d resume_log_index=%d resume_world_ref=%s",
			r.LogRowsScanned, r.WorldRowsScanned, len(r.Holes), r.ResumeLogIndex, r.ResumeWorldRef))
	}
	return lines
}

// HealthResponse is the GET /v1/health body: daemon version, db path and the
// archived-interpreter HashRef + verbatim `ailang --version` recorded in the
// archive manifest (frozen route table, Decision 3).
type HealthResponse struct {
	Status string `json:"status"`
	// DaemonVersion is this package's Version.
	DaemonVersion string `json:"daemon_version"`
	// DBPath is the configured store database path.
	DBPath string `json:"db_path"`
	// InterpreterRef is the archived interpreter's canonical "algo:digest", or
	// "" when serve was started without --ailang-bin.
	InterpreterRef string `json:"interpreter_ref"`
	// InterpreterVersion is the verbatim `ailang --version` STDOUT captured in
	// the archive manifest, or "" when no interpreter was archived. STDOUT
	// specifically: the archive probe captures the two streams separately so
	// interpreter chatter on stderr never becomes part of a content-addressed
	// artifact's recorded identity, nor of the epoch-registry candidate
	// releaseFromVersion derives from this same string.
	InterpreterVersion string `json:"interpreter_version"`
}

// resolveErrorLog turns Config.ErrorLog's optional nil into the concrete
// default ONCE, at construction. Resolving at construction rather than at each
// write is what makes the default assertable: a test reads d.errLog and sees
// os.Stderr, instead of having to prove a negative about a nil branch it cannot
// observe.
func resolveErrorLog(w io.Writer) io.Writer {
	if w == nil {
		return os.Stderr
	}
	return w
}

// New runs the pre-listen half of the serve lifecycle (Decision 5):
//
//	loopback bind check -> store.Open (fail-closed writer) -> optional
//	archive.Archive + ReadManifest -> registry.Bootstrap (idempotent;
//	divergent head = FATAL) -> build the D7-bounded http.Server
//
// The bind policy is checked FIRST, before any writer lock is taken, so a
// misconfigured bind never disturbs a database another process is serving.
// Every failure after store.Open closes the store, releasing the writer lock:
// a refused startup must not strand writer authority.
// ctx is threaded to bootstrap reads; it does not bound every startup operation.
// Pass the caller's intended bootstrap-read ctx; this API supplies no default timeout.
func New(ctx context.Context, cfg Config) (*Daemon, error) {
	if !isLoopbackHost(cfg.BindHost) {
		return nil, &StartupError{
			Stage: StageBindPolicy,
			Detail: fmt.Sprintf(
				"bind host %q is not loopback (allowed: 127.0.0.1, ::1, localhost); "+
					"M2 is local-first by construction and ships no override flag", cfg.BindHost),
			Err: ErrNonLoopbackBind,
		}
	}
	if cfg.DBPath == "" {
		return nil, &StartupError{Stage: StageConfig, Detail: "no database path configured (--db is required)"}
	}
	if cfg.WorkspaceRoot == "" && (cfg.WorkspaceModuleRoot != "" || len(cfg.WorkspaceEpisodeModuleRoots) > 0 || cfg.WorkspacePackageCache != "") {
		return nil, &StartupError{Stage: StageConfig, Detail: "--workspace-module-root, --workspace-episode-module-root and --workspace-package-cache need --workspace-root"}
	}
	// Row 134 AC4.5: refuse a workspace root that contains the daemon's own
	// state BEFORE taking writer authority, like the bind policy.
	workspace := &workspaceTools{errLog: resolveErrorLog(cfg.ErrorLog)}
	if cfg.WorkspaceRoot != "" {
		root, stateDir, err := resolveWorkspaceRoot(cfg.WorkspaceRoot, cfg.DBPath)
		if err != nil {
			return nil, &StartupError{Stage: StageConfig, Detail: "the workspace root is refused", Err: err}
		}
		workspace.root, workspace.stateDir = root, stateDir
		if cfg.ExamplesDir != "" {
			examples, err := resolveExamplesDir(cfg.ExamplesDir, root)
			if err != nil {
				return nil, &StartupError{Stage: StageConfig, Detail: "the examples corpus directory is refused", Err: err}
			}
			workspace.examplesDir = examples
		}
		def, byEp, err := resolveModuleRoots(cfg.WorkspaceModuleRoot, cfg.WorkspaceEpisodeModuleRoots, root)
		if err != nil {
			return nil, &StartupError{Stage: StageConfig, Detail: "the module root is refused", Err: err}
		}
		workspace.moduleRoot, workspace.episodeModuleRoot = def, byEp
		if cfg.WorkspacePackageCache != "" {
			cache, digest, n, err := resolvePackageCache(cfg.WorkspacePackageCache, root, stateDir)
			if err != nil {
				return nil, &StartupError{Stage: StageConfig, Detail: "the --workspace-package-cache directory is refused", Err: err}
			}
			workspace.packageCache, workspace.packageCacheDigest = cache, digest
			fmt.Fprintf(workspace.errLog, "ailang-worldd: workspace package cache %s: %d packages, read-only, %s\n", cache, n, digest)
		}
	}

	s, err := store.Open(cfg.DBPath)
	if err != nil {
		detail := "cannot open the world store for writing"
		if store.IsWriterAlreadyActive(err) {
			detail = "another process already holds writer authority for this database " +
				"(single-writer is enforced, not conventional)"
		}
		return nil, &StartupError{Stage: StageStoreOpen, Detail: detail, Err: err}
	}
	if err := checkReadOrdering(s.BusyTimeout(), readDeadline); err != nil {
		_ = s.Close()
		return nil, &StartupError{Stage: StageStoreOpen,
			Detail: "the store's lock-retry window is not below the read deadline", Err: err}
	}

	d := &Daemon{
		cfg: cfg, store: s, bootstrap: registry.Bootstrap, reads: s, commits: s, commitBudget: commitBudget, credentialBudget: credentialBudget, drainTimeout: shutdownTimeout,
		aguiBudget: aguiRunBudget, aguiBodyBound: aguiBodyBound, aguiTick: aguiTick,
		aguiStop: make(chan struct{}), aguiSlots: make(chan struct{}, aguiCap),
		readDeadline: readDeadline, errLog: resolveErrorLog(cfg.ErrorLog),
		scanPageSize: integrityScanPageSize, scanRowBudget: integrityScanRowBudget,
		scanTimeBudget: integrityScanTimeBudget, resolver: authority.New(s), workspace: workspace,
	}
	release := unpinnedRelease
	if cfg.ToolAilangBin != "" {
		bin, ref, err := archiveToolBinary(archive.New(cfg.DBPath), cfg.ToolAilangBin)
		if err != nil {
			return nil, d.abort(StageArchive, "cannot archive the configured workspace tool binary", err)
		}
		workspace.bin, workspace.binRef = bin, ref
	}
	if err := workspace.configureRunCaps(ctx, cfg); err != nil {
		return nil, d.abort(StageConfig, "the run capabilities are refused", err)
	}
	if err := workspace.configureExec(ctx, cfg); err != nil {
		return nil, d.abort(StageConfig, "the exec configuration is refused", err)
	}
	if (cfg.WorkspaceRoot == "") != (cfg.ToolAilangBin == "") {
		fmt.Fprintln(d.errLog, "ailang-worldd: workspace tools disabled: --workspace-root and --tool-ailang-bin must both be set; "+
			"every effect-declaring transition is refused")
	}

	if cfg.AilangBin != "" {
		a := archive.New(cfg.DBPath)
		ref, err := a.Archive(cfg.AilangBin)
		if err != nil {
			return nil, d.abort(StageArchive, "cannot archive the configured interpreter", err)
		}
		m, err := a.ReadManifest(ref)
		if err != nil {
			return nil, d.abort(StageArchive, "cannot read the archived interpreter manifest", err)
		}
		d.interpreterRef = ref.String()
		d.arch, d.interpreterHash = a, ref
		d.interpreterVersion = m.Version
		release = releaseFromVersion(m.Version)
		// The capsule's log is the daemon's operator log: a cold (template-less)
		// run writes its `capsule: cold compile` line there (row 153).
		d.capsule = capsule.New(a, capsule.Config{Log: d.errLog})
		d.coord, err = coordinator.New(coordinator.Config{Store: d.store, Runner: d.capsule,
			Binder: d.binder, Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
		if err != nil {
			return nil, d.abort(StageConfig, "cannot construct invocation coordinator", err)
		}
	}
	if invokeDeadline <= 0 || invokeDeadline >= writeTimeout {
		return nil, d.abort(StageConfig, "invocation deadline must be shorter than write timeout", nil)
	}

	// Idempotent by construction; a head naming different bytes is a genuine
	// divergence and registry.Bootstrap returns an error for it. Surfacing that
	// as a fatal StartupError is the point: a divergent registry head must never
	// be silently accepted or rewritten.
	if err := d.bootstrapRegistry(ctx, release); err != nil {
		return nil, err
	}

	// Row 153 M3b: say, per descriptor, whether its compile-cache template is
	// ready. A stat each, no build; the rebuild runs from Run after Listen.
	d.reportTemplates(ctx)

	d.integrity = d.scanIntegrity(ctx)

	// The A2A projection (w-a2a-session-projection) mounts two additive routes
	// over the handles New already owns: the ONE resolver instance (F3), the
	// registry reader over the ONE open store (the projection holds no store
	// handle at all — P5 made structural), and the read seam for the B3
	// absent-head pre-check. The card-route denial writer is the
	// middleware's own writeSessionDenial so a card denial is byte-identical
	// to the /v1/commit denial (B1); the card route's 503/504 failures reuse
	// the daemon's writeAPIError envelope — the envelope stays daemon-owned
	// and projection never formats it (AC1). MaxWait is the D7 readDeadline:
	// a projection request IS store reads below the transport, so the same
	// proven bound applies.
	proj, err := projection.New(projection.Config{
		Resolver: d.resolver,
		Reader:   transitionreg.NewReader(d.store),
		Heads:    d.reads,
		Deny:     writeSessionDenial,
		Fail:     writeAPIError,
		ErrorLog: d.errLog,
		Agent: protocol.AgentInfo{
			Name:        "ailang-worldd",
			Description: "AILANG World daemon: a published skill accepts one JSON-object data part and produces a JSON-object data artifact; invocation requires an archived interpreter; schemas are carried but not validated.",
			Version:     Version,
		},
		MaxWait:        readDeadline,
		InvokeWait:     invokeDeadline,
		CredentialWait: credentialBudget, CallbackTimeout: invokeDeadline, MaxCallbacks: procbound.MaxOutstanding, WriteWait: writeTimeout,
		Coordinator: d.coord,
	})
	if err != nil {
		return nil, d.abort(StageConfig, "cannot construct the A2A projection", err)
	}
	d.projection = proj

	d.srv = newServer(d.Handler())
	d.srv.RegisterOnShutdown(func() { d.aguiStopOnce.Do(func() { close(d.aguiStop) }) })
	return d, nil
}

// binder is the coordinator's production BinderFor: a live binder over the
// episode's workspace registry (row 134 §4.3) — empty, so R8 refuses every
// declared effect, unless the episode resolves to a served worktree; then
// Workspace.Exec is bound and the eight AILANG names only when that
// episode's AILANG handler builds (row 152).
func (d *Daemon) binder(episodeID string, grants []broker.Capability) transitionreg.Binder {
	return broker.OpenBinder(d.store, episodeID, grants, d.workspace.registry(episodeID))
}

func (d *Daemon) bootstrapRegistry(ctx context.Context, release string) error {
	if _, _, err := d.bootstrap(ctx, d.store, release); err != nil {
		detail := fmt.Sprintf("cannot bootstrap %s with release %q", registry.SemanticID, release)
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return d.abortBudget(StageRegistry, detail, ctx.Err())
		}
		return d.abort(StageRegistry, detail, err)
	}
	return nil
}

func (d *Daemon) scanIntegrity(ctx context.Context) IntegrityReport {
	start := time.Now()
	report := IntegrityReport{}
	logDone, worldDone := false, false
	for !logDone || !worldDone {
		remaining := d.scanRowBudget - report.LogRowsScanned - report.WorldRowsScanned
		if remaining <= 0 || time.Since(start) >= d.scanTimeBudget {
			return report
		}
		limit := d.scanPageSize
		if limit > remaining {
			limit = remaining
		}
		if !logDone {
			page, err := d.store.ScanUnreadableLog(ctx, report.ResumeLogIndex, limit)
			if err != nil {
				report.Holes = append(report.Holes, store.UnreadableRow{
					Table: "log_entries", Index: report.ResumeLogIndex,
					Field: "scan", Reason: err.Error(),
				})
				return report
			}
			report.LogRowsScanned += page.Scanned
			report.ResumeLogIndex = page.NextIndex
			report.Holes = append(report.Holes, page.Rows...)
			logDone = page.Done
		}
		remaining = d.scanRowBudget - report.LogRowsScanned - report.WorldRowsScanned
		if remaining <= 0 || time.Since(start) >= d.scanTimeBudget {
			return report
		}
		limit = d.scanPageSize
		if limit > remaining {
			limit = remaining
		}
		if !worldDone {
			page, err := d.store.ScanUnreadableWorlds(ctx, report.ResumeWorldRef, limit)
			if err != nil {
				report.Holes = append(report.Holes, store.UnreadableRow{
					Table: "worlds", Ref: report.ResumeWorldRef,
					Field: "scan", Reason: err.Error(),
				})
				return report
			}
			report.WorldRowsScanned += page.Scanned
			report.ResumeWorldRef = page.NextRef
			report.Holes = append(report.Holes, page.Rows...)
			worldDone = page.Done
		}
	}
	report.Complete = true
	return report
}

// abort releases writer authority and wraps err as a fatal StartupError. Used
// for every failure that happens after store.Open has succeeded.
func (d *Daemon) abort(stage, detail string, err error) error {
	_ = d.store.Close()
	return &StartupError{Stage: stage, Detail: detail, Err: err}
}

// abortBudget hands cleanup to a retained owner so the startup error can
// return while a post-cutoff COMMIT still owns the sole connection.
func (d *Daemon) abortBudget(stage, detail string, err error) error {
	store.Quarantine(d.store)
	settled := make(chan struct{})
	go func() {
		defer close(settled)
		_ = d.store.Close()
	}()
	return &StartupError{Stage: stage, Detail: detail, Err: err, Settled: settled}
}

// releaseFromVersion reduces a verbatim `ailang --version` output (which is
// multi-line: version, commit, build date, banner) to its first non-empty line,
// e.g. "AILANG v0.30.0". That line is the epoch-registry candidate string,
// matching the "AILANG v0.30.0 (commit e37b370)" shape M1's registry tests use.
func releaseFromVersion(version string) string {
	for _, line := range strings.Split(version, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			return trimmed
		}
	}
	return unpinnedRelease
}

// Handler builds the route table. ServeMux METHOD PATTERNS (go 1.26) make the
// method part of the pattern, so a non-GET on these paths is a 405 from the mux
// rather than a hand-rolled check.
//
// POST /agui/ is an additive read-only AG-UI 1.0 SSE run, bounded to 18 s
// from handler entry. Resume with standard state or Last-Event-ID; see
// docs/QUICKSTART.md, "Watch the world live".
// GET /workbench/live.js serves the embedded same-origin live client; see
// docs/QUICKSTART.md, "Open the live workbench" for the read-only panes.
//
// The ten /v1 patterns below are the complete frozen v1 machine table (nine GET, one POST;
// GET /v1/receipts/{id} added by row 23 M3, D-WORLD-40).
// The tenth registration, GET /workbench, is the unversioned read-only operator renderer: it is
// NOT part of the frozen table, its HTML may evolve, and it changes the semantics of none of the
// ten. The registry pattern deliberately uses a multi-segment wildcard: registry semantic IDs
// such as "world/epoch-registry/v1" contain slashes.
func (d *Daemon) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", d.handleHealth)
	mux.HandleFunc("GET /v1/head", d.handleHead)
	mux.HandleFunc("GET /v1/worlds/{ref}", d.handleWorld)
	mux.HandleFunc("GET /v1/objects/{ref}", d.handleObject)
	mux.HandleFunc("GET /v1/objects/by-semantic-id/{name...}", d.handleObjectsBySemanticID)
	mux.HandleFunc("GET /v1/log/{index}", d.handleLogEntry)
	mux.HandleFunc("GET /v1/log", d.handleLogRange)
	mux.HandleFunc("GET /v1/registry/{name...}", d.handleRegistry)
	mux.HandleFunc("POST /v1/commit", d.handleCommit)
	mux.HandleFunc("GET /v1/receipts/{id}", d.handleReceipt)
	mux.HandleFunc("GET /workbench", d.handleWorkbench)
	mux.HandleFunc("GET /workbench/live.js", d.handleWorkbenchScript)
	// The two A2A projection routes (w-a2a-session-projection P6.B-A2A-CARD)
	// are ADDITIVE: the frozen /v1/ table above is untouched, and the routes
	// are NOT in isProtected — the projection handler resolves the session
	// itself so /a2a/ can answer in JSON-RPC form (B5).
	mux.HandleFunc("GET /.well-known/agent.json", d.projection.AgentCard)
	mux.HandleFunc("POST /a2a/", d.projection.A2A)
	mux.HandleFunc("POST /mcp/", d.projection.MCP)
	mux.HandleFunc("POST /agui/", d.handleAGUI)
	return NewSessionMiddleware(d.resolver, d.credentialBudget, d.writeInternalError).Wrap(d.isProtected, mux)
}

// isProtected reports whether a request must carry a valid session credential
// (w-session-authority D6). This sprint it is ONLY POST /v1/commit — the
// daemon's sole mutation. The eight GET routes pass through unauthenticated as
// declared residual R1 (full /v1/* read enforcement is a follow-up queue row; a
// config flip here later is a rewrite-free change).
func (d *Daemon) isProtected(r *http.Request) bool {
	return r.Method == http.MethodPost && r.URL.Path == "/v1/commit"
}

// handleHealth serves GET /v1/health.
func (d *Daemon) handleHealth(w http.ResponseWriter, _ *http.Request) {
	body := HealthResponse{
		Status:             "ok",
		DaemonVersion:      Version,
		DBPath:             d.cfg.DBPath,
		InterpreterRef:     d.interpreterRef,
		InterpreterVersion: d.interpreterVersion,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		// The status line is already written; there is nothing left to say to
		// the client, and the server must not panic on a dropped connection.
		return
	}
}

// handleHead serves GET /v1/head: the store's selected head as canonical
// "algo:digest" text (Decision 3 — one hash encoding everywhere; the daemon
// never invents a second one).
//
// Success remains canonical plain text. Errors use the same JSON APIError
// envelope as every other v1 route: no selected head is NotFound (404), and a
// store failure is Internal (500) — sanitized through writeInternalError like
// every other 500, because this route reaches the store through the SAME read
// seam as the five in handlers.go and leaks exactly the same host detail.
func (d *Daemon) handleHead(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := d.readCtx(r)
	defer cancel()
	ref, ok, err := d.reads.SelectedHead(ctx)
	if err != nil {
		if timedOut(ctx, err) {
			writeReadTimeout(w, d.readDeadline)
			return
		}
		d.writeInternalError(w, r, err)
		return
	}
	if !ok {
		writeAPIError(w, "NotFound", "no world head has been selected yet", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.WriteString(w, ref.String())
}

// newServer constructs the daemon's http.Server with ALL FOUR D7 timeouts set
// at construction. It is a named constructor precisely so the timeouts have a
// single definition site that a test can assert against — a zero value in any
// of these fields is a test failure, not a default.
func newServer(h http.Handler) *http.Server {
	return &http.Server{
		Handler:           h,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}

// Listen binds the configured loopback socket. Port 0 yields a kernel-assigned
// ephemeral port, which is why Addr must be read back from the listener rather
// than reconstructed from Config.
func (d *Daemon) Listen() error {
	addr := net.JoinHostPort(d.cfg.BindHost, strconv.Itoa(d.cfg.BindPort))
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return &StartupError{Stage: StageListen, Detail: "cannot bind " + addr, Err: err}
	}
	d.ln = ln
	return nil
}

// Addr returns the RESOLVED listen address ("host:port"), empty before Listen.
func (d *Daemon) Addr() string {
	if d.ln == nil {
		return ""
	}
	return d.ln.Addr().String()
}

// URL returns the resolved base URL a client should pass as --addr.
func (d *Daemon) URL() string {
	if d.ln == nil {
		return ""
	}
	return "http://" + d.Addr()
}

// Serve runs the HTTP server until Shutdown or Close; it returns
// http.ErrServerClosed on an orderly stop, exactly like http.Server.Serve.
func (d *Daemon) Serve() error { return d.srv.Serve(d.ln) }

// Shutdown performs the D7 bounded graceful drain: in-flight requests finish
// under drainTimeout (= shutdownTimeout), and if they do not, connections are
// hard-closed and a non-nil error is returned so the caller can exit non-zero.
// The drain is never unbounded.
func (d *Daemon) Shutdown() error { return drain(d.srv, d.drainTimeout) }

// Close releases writer authority by closing the store (which releases the
// cross-process writer lock).
func (d *Daemon) Close() error {
	d.stopTemplateMaintenance()
	return d.store.Close()
}

// shutdowner is the http.Server subset drain needs. It exists so the
// deadline-expiry branch — the branch that must never be reachable in
// production, and therefore the branch a real server cannot easily exercise —
// is still honestly tested rather than assumed.
type shutdowner interface {
	Shutdown(ctx context.Context) error
	Close() error
}

// drain implements Decision 7's shutdown rule: Shutdown(ctx) under a deadline,
// then a hard Close() if the deadline expires, reporting the incomplete drain
// instead of waiting it out forever.
func drain(srv shutdowner, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		_ = srv.Close()
		return fmt.Errorf(
			"daemon: graceful drain did not finish within %s; connections were hard-closed: %w",
			timeout, err)
	}
	return nil
}

// Run is the whole serve lifecycle (Decision 5):
//
//	New -> Listen -> ANNOUNCE the resolved address on announce -> Serve ->
//	ctx cancelled (SIGINT/SIGTERM at the caller) -> bounded Shutdown -> Close
//
// The announcement is written AFTER the socket is bound and BEFORE serving, so
// a caller that reads the line is guaranteed the port is already accepting.
// Run returns a non-nil error when the drain did not finish, so the process can
// exit non-zero on an incomplete shutdown.
func startBounded(ctx context.Context, cfg Config, budget time.Duration, construct func(context.Context, Config) (*Daemon, error)) (*Daemon, error) {
	startupCtx, cancel := context.WithTimeout(ctx, budget)
	defer cancel()
	return construct(startupCtx, cfg)
}

func Run(ctx context.Context, cfg Config, announce io.Writer) error {
	d, err := startBounded(ctx, cfg, 9*time.Second, New)
	if err != nil {
		return err
	}
	defer func() { _ = d.Close() }()

	if err := d.Listen(); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(announce, "%s%s\n", ListenAnnouncePrefix, d.URL()); err != nil {
		return &StartupError{Stage: StageListen, Detail: "cannot announce the resolved listen address", Err: err}
	}
	// Announce the sweep ONLY when there is something to warn about — a hole, or a
	// scan that could not finish. A clean, complete sweep says nothing, so a healthy
	// store's announce stream stays byte-identical to the pre-sweep contract.
	//
	// It is NOT a silent skip: holes and truncation — the two states an operator
	// must act on — are still reported, and a truncated scan still says so
	// explicitly rather than reading as a clean bill of health. The daemon never
	// emits an all-clear for a store it did not fully verify.
	//
	// Be precise about WHICH guarantee this is, because the two are different and
	// only one of them holds. TRUTHFULNESS is guaranteed: no output ever claims a
	// store is clean when it is not. DELIVERY is best-effort: a consumer that
	// reads to EOF — which is what an operator running the daemon as a subprocess
	// does — receives every line, but a caller that stops reading after the
	// listen line, or a process torn down before the goroutine is scheduled, may
	// never see the warnings. That is the deliberate trade for never letting a
	// diagnostic wedge startup, and it is a weaker DELIVERY promise than the
	// synchronous form had — not a weaker claim about the store.
	//
	// The write happens on its own goroutine, and that is load-bearing rather than
	// stylistic. `announce` is frequently an io.Pipe (daemon_test.go:585 and the
	// CLI/main subprocess tests): SYNCHRONOUS and unbuffered, so every Write blocks
	// until a reader consumes it, and those callers read exactly ONE line — the
	// listen announcement — then stop reading. A synchronous write here therefore
	// blocked Run before it reached Serve() below, so the socket it had just
	// announced never served (observed as GET /v1/health timing out).
	//
	// Emitting only warnings fixed that for a healthy store, but NOT for a store
	// that genuinely has something to say: with a hole present, or a sweep
	// truncated by the row/time budget, the same one-line reader deadlocked again.
	// Startup diagnostics must never be able to wedge startup, whatever the store
	// contains — so Run never waits on the consumer. Ordering is still guaranteed:
	// the listen line is written synchronously above, before this goroutine starts.
	// A consumer that stops reading simply leaves the goroutine parked until the
	// stream is closed; it cannot delay Serve, the health surface, or shutdown.
	if report := d.IntegrityReport(); !report.Complete || len(report.Holes) > 0 {
		lines := report.Lines()
		go func() {
			for _, line := range lines {
				if _, err := fmt.Fprintln(announce, line); err != nil {
					return // the consumer went away; a diagnostic must not outlive it
				}
			}
		}()
	}

	// Row 153 M3b: rebuild missing templates only now that the listener is bound
	// and announced; the goroutine is cancelled and joined by Close.
	d.startTemplateMaintenance(ctx)

	served := make(chan error, 1)
	go func() { served <- d.Serve() }()

	select {
	case err := <-served:
		// The server stopped on its own (a listener failure); ErrServerClosed
		// cannot occur here because nothing has asked it to stop yet.
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
	}

	drainErr := d.Shutdown()
	if err := <-served; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return drainErr
}
