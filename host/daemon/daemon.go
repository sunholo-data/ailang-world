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
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/registry"
	"github.com/sunholo-data/ailang-world/host/store"
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
	writeTimeout = 30 * time.Second

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
	cfg   Config
	store *store.Store
	srv   *http.Server
	ln    net.Listener

	// drainTimeout bounds the graceful shutdown. New always sets it to the D7
	// shutdownTimeout; it is a field rather than a direct constant reference so
	// the bound is (a) assertable as wiring and (b) shrinkable in tests, which
	// is the only way the expiry branch of the SHIPPED Shutdown path can be
	// exercised without a ten-second test.
	drainTimeout   time.Duration
	scanPageSize   int
	scanRowBudget  int
	scanTimeBudget time.Duration
	integrity      IntegrityReport

	// Health facts resolved once at startup and served verbatim.
	interpreterRef     string
	interpreterVersion string
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
	// InterpreterVersion is the verbatim `ailang --version` output captured in
	// the archive manifest, or "" when no interpreter was archived.
	InterpreterVersion string `json:"interpreter_version"`
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
func New(cfg Config) (*Daemon, error) {
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

	s, err := store.Open(cfg.DBPath)
	if err != nil {
		detail := "cannot open the world store for writing"
		if store.IsWriterAlreadyActive(err) {
			detail = "another process already holds writer authority for this database " +
				"(single-writer is enforced, not conventional)"
		}
		return nil, &StartupError{Stage: StageStoreOpen, Detail: detail, Err: err}
	}

	d := &Daemon{
		cfg: cfg, store: s, drainTimeout: shutdownTimeout,
		scanPageSize: integrityScanPageSize, scanRowBudget: integrityScanRowBudget,
		scanTimeBudget: integrityScanTimeBudget,
	}
	release := unpinnedRelease

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
		d.interpreterVersion = m.Version
		release = releaseFromVersion(m.Version)
	}

	// Idempotent by construction; a head naming different bytes is a genuine
	// divergence and registry.Bootstrap returns an error for it. Surfacing that
	// as a fatal StartupError is the point: a divergent registry head must never
	// be silently accepted or rewritten.
	if _, _, err := registry.Bootstrap(s, release); err != nil {
		return nil, d.abort(StageRegistry,
			fmt.Sprintf("cannot bootstrap %s with release %q", registry.SemanticID, release), err)
	}

	d.integrity = d.scanIntegrity()
	d.srv = newServer(d.Handler())
	return d, nil
}

func (d *Daemon) scanIntegrity() IntegrityReport {
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
			page, err := d.store.ScanUnreadableLog(report.ResumeLogIndex, limit)
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
			page, err := d.store.ScanUnreadableWorlds(report.ResumeWorldRef, limit)
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
// The seven patterns below are the complete frozen v1 table. The registry
// pattern deliberately uses a multi-segment wildcard: registry semantic IDs
// such as "world/epoch-registry/v1" contain slashes.
func (d *Daemon) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", d.handleHealth)
	mux.HandleFunc("GET /v1/head", d.handleHead)
	mux.HandleFunc("GET /v1/worlds/{ref}", d.handleWorld)
	mux.HandleFunc("GET /v1/objects/{ref}", d.handleObject)
	mux.HandleFunc("GET /v1/log/{index}", d.handleLogEntry)
	mux.HandleFunc("GET /v1/log", d.handleLogRange)
	mux.HandleFunc("GET /v1/registry/{name...}", d.handleRegistry)
	mux.HandleFunc("POST /v1/commit", d.handleCommit)
	return mux
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
// store failure is Internal (500).
func (d *Daemon) handleHead(w http.ResponseWriter, _ *http.Request) {
	ref, ok, err := d.store.SelectedHead()
	if err != nil {
		writeAPIError(w, "Internal", err.Error(), http.StatusInternalServerError)
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
func (d *Daemon) Close() error { return d.store.Close() }

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
func Run(ctx context.Context, cfg Config, announce io.Writer) error {
	d, err := New(cfg)
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
	// This is load-bearing, not a style choice. `announce` is frequently an
	// io.Pipe (daemon_test.go:585, and the CLI/main subprocess tests), which is
	// SYNCHRONOUS and unbuffered: every Write blocks until a reader consumes it.
	// Those callers read exactly ONE line — the listen announcement — and then stop
	// reading. Writing any further line here therefore blocks Run forever, so it
	// never reaches Serve() below and the socket it just announced never serves.
	// Emitting only warnings keeps the healthy path write-free and non-blocking;
	// TestRunAnnouncesResolvedListenAddress is the standing guard (it deadlocked,
	// and timed out on GET /v1/health, when this loop ran unconditionally).
	//
	// It is NOT a silent skip: holes and truncation — the two states an operator
	// must act on — are still reported in full, and a truncated scan still says so
	// explicitly rather than reading as a clean bill of health.
	if report := d.IntegrityReport(); !report.Complete || len(report.Holes) > 0 {
		for _, line := range report.Lines() {
			if _, err := fmt.Fprintln(announce, line); err != nil {
				return &StartupError{Stage: StageListen, Detail: "cannot announce the integrity scan", Err: err}
			}
		}
	}

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
