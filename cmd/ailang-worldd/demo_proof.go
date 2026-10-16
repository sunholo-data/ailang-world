package main

// M3a (row 150): the first production caller of evidence.NewProducer. The
// adapters below speak only to the running daemon's REST API: GET
// /v1/objects/{ref}?payload=true to read and POST /v1/commit to write. They
// hold no store handle; S4 uses them after the recorded tool stages.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Bounds from the sprint plan. Each REST call is cut at 10 s, below the 12 s
// object read and write bounds, so the HTTP client always answers first.
const (
	demoProofHTTPTimeout = 10 * time.Second
	demoProofObjectBound = 12 * time.Second
	demoProofExecBound   = 30 * time.Second
	demoProofMaxOutput   = 1 << 20
	demoProofWholeBound  = 90 * time.Second
	demoProofVersion     = "AILANG v0.41.0"
)

// demoProofClient is the production REST client used by seed and both proof
// adapters. doRequest derives a per-request deadline from timeout.
func demoProofClient(addr string) *client {
	c := newClient(addr)
	c.timeout = demoProofHTTPTimeout
	c.http.Timeout = demoProofHTTPTimeout
	return c
}

// demoLoadProofKey mirrors daemon.loadProofKey: stat the open file, require a
// regular file, read at most 33 bytes, and never include key bytes in errors.
func demoLoadProofKey(path string) ([32]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return [32]byte{}, fmt.Errorf("proof key: open: %w", err)
	}
	defer f.Close()
	return demoReadProofKey(f)
}

// The narrow file interface also lets error controls exercise stat/read errors
// without races or a second path lookup. Production always passes the open file.
func demoReadProofKey(f interface {
	io.Reader
	Stat() (os.FileInfo, error)
}) ([32]byte, error) {
	var key [32]byte
	info, err := f.Stat()
	if err != nil {
		return key, fmt.Errorf("proof key: stat: %w", err)
	}
	if !info.Mode().IsRegular() {
		return key, fmt.Errorf("proof key: not a regular file")
	}
	raw, err := io.ReadAll(io.LimitReader(f, 33))
	if err != nil {
		return key, fmt.Errorf("proof key: read: %w", err)
	}
	if len(raw) > len(key) {
		return key, fmt.Errorf("proof key: holds more than 32 bytes; expected exactly 32")
	}
	if len(raw) < len(key) {
		return key, fmt.Errorf("proof key: holds %d bytes; expected exactly 32", len(raw))
	}
	copy(key[:], raw)
	return key, nil
}

func demoProofContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, demoProofWholeBound)
}

// demoProofBounds is the producer configuration every proof runs under. The
// adapters have no store lock-retry window, so BusyTimeout is zero.
func demoProofBounds() evidence.ProducerConfig {
	return evidence.ProducerConfig{
		ObjectReadTimeout:  demoProofObjectBound,
		ObjectWriteTimeout: demoProofObjectBound,
		ExecTimeout:        demoProofExecBound,
		MaxOutputBytes:     demoProofMaxOutput,
	}
}

// demoProofChecker pins the checker to the health-archived interpreter. The
// executable is the archived copy of the ref that health reports, and its bytes
// must hash to that ref. The producer hashes the bytes again before any run.
func demoProofChecker(h daemon.HealthResponse) (evidence.CheckerConfig, error) {
	line := strings.TrimRight(strings.SplitN(h.InterpreterVersion, "\n", 2)[0], "\r")
	if line != demoProofVersion {
		return evidence.CheckerConfig{}, fmt.Errorf("proof checker: interpreter version %q is not %s", line, demoProofVersion)
	}
	ref, err := hashref.Parse(h.InterpreterRef)
	if err != nil {
		return evidence.CheckerConfig{}, fmt.Errorf("proof checker: interpreter ref: %w", err)
	}
	path, err := archive.New(h.DBPath).Resolve(ref)
	if err != nil {
		return evidence.CheckerConfig{}, fmt.Errorf("proof checker: archived interpreter: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return evidence.CheckerConfig{}, fmt.Errorf("proof checker: read archived interpreter: %w", err)
	}
	if got := hashref.SumSHA256(data); got != ref {
		return evidence.CheckerConfig{}, fmt.Errorf("proof checker: archived pin bytes hash to %s; health pins %s", got, ref)
	}
	return evidence.CheckerConfig{ExecPath: path, Compiler: ref, CompilerVersion: line}, nil
}

// demoProofReader is the producer's ObjectReader over GET /v1/objects/{ref}.
type demoProofReader struct{ c *client }

// BusyTimeout is zero: a REST read has no store lock-retry window to order.
func (demoProofReader) BusyTimeout() time.Duration { return 0 }

// ReadObject returns the payload of ref. An absent object is (zero, nil, nil),
// which the producer names itself; the reader never invents bytes.
func (r demoProofReader) ReadObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (store.ObjectMeta, []byte, error) {
	status, body, err := r.c.do(ctx, http.MethodGet, "/v1/objects/"+url.PathEscape(ref.String())+"?payload=true", nil)
	if err != nil {
		return store.ObjectMeta{}, nil, fmt.Errorf("proof reader: GET object %s: %w", ref, err)
	}
	if status == http.StatusNotFound {
		return store.ObjectMeta{}, nil, nil
	}
	if status != http.StatusOK {
		return store.ObjectMeta{}, nil, fmt.Errorf("proof reader: GET object %s returned HTTP %d", ref, status)
	}
	var wire struct {
		Hash          string  `json:"hash"`
		InterfaceHash string  `json:"interfaceHash"`
		SemanticID    string  `json:"semanticId"`
		Provenance    string  `json:"provenance"`
		Payload       *[]byte `json:"payload"`
	}
	if err := json.Unmarshal(body, &wire); err != nil {
		return store.ObjectMeta{}, nil, fmt.Errorf("proof reader: object %s: malformed body: %w", ref, err)
	}
	if wire.Hash != ref.String() {
		return store.ObjectMeta{}, nil, fmt.Errorf("proof reader: object %s answered as %q", ref, wire.Hash)
	}
	if wire.Payload == nil {
		return store.ObjectMeta{}, nil, fmt.Errorf("proof reader: object %s response carries no payload", ref)
	}
	if n := int64(len(*wire.Payload)); n > maxBytes {
		return store.ObjectMeta{}, nil, &store.ObjectTooLargeError{Size: n, MaxBytes: maxBytes}
	}
	iface, err := hashref.Parse(wire.InterfaceHash)
	if err != nil {
		return store.ObjectMeta{}, nil, fmt.Errorf("proof reader: object %s interface hash: %w", ref, err)
	}
	meta := store.ObjectMeta{InterfaceHash: iface, SemanticID: wire.SemanticID, Provenance: wire.Provenance, PayloadLength: int64(len(*wire.Payload))}
	return meta, *wire.Payload, nil
}

// demoProofWriter is the producer's ObjectWriter: one object per POST
// /v1/commit, appended as the next entry of the selected world. A write whose
// outcome is unknown is reported as store.UncertainError and never retried.
type demoProofWriter struct {
	c           *client
	token       string
	interpreter string // the pinned interpreter ref; the head entry must carry it
}

// BusyTimeout is zero: the REST commit has no store lock-retry window to order.
func (demoProofWriter) BusyTimeout() time.Duration { return 0 }

func (w demoProofWriter) PutObject(ctx context.Context, o store.Object) error {
	obj := demoObject{Hash: o.Hash.String(), InterfaceHash: o.InterfaceHash.String(), SemanticID: o.SemanticID, Provenance: o.Provenance, Payload: o.Payload}
	_, err := commitOneObject(ctx, w.c, w.token, obj, "demo", w.interpreter, "proof")
	return err
}

// commitOneObject appends obj as the next log entry over POST /v1/commit,
// linked to the selected head's entry and written by writtenBy. When
// interpreter is non-empty the head entry must carry it (the proof's pin);
// otherwise the new entry repeats the head entry's interpreter. It returns the
// committed entry index. A transport error, 503 CommitUncertain, or a 200 that
// names another head is a *store.UncertainError (never retried); any other
// non-200 is a plain error. Shared by the row-150 proof writer and row 182's
// `outcome record`; op names the caller in errors.
func commitOneObject(ctx context.Context, c *client, token string, obj demoObject, writtenBy, interpreter, op string) (int64, error) {
	head, has, err := headState(ctx, c)
	if err != nil {
		return 0, fmt.Errorf("%s writer: read head: %w", op, err)
	}
	if !has {
		return 0, fmt.Errorf("%s writer: no selected head to commit onto", op)
	}
	r := worldReader{ctx: ctx, c: c}
	world, found, err := r.world(head)
	if err != nil {
		return 0, fmt.Errorf("%s writer: read world %s: %w", op, head, err)
	}
	if !found {
		return 0, fmt.Errorf("%s writer: selected head %s has no world row", op, head)
	}
	last, found, err := r.entry(world.Revision)
	if err != nil {
		return 0, fmt.Errorf("%s writer: read head entry: %w", op, err)
	}
	if !found {
		return 0, fmt.Errorf("%s writer: head entry %d is absent", op, world.Revision)
	}
	if interpreter == "" {
		interpreter = last.Header.Interpreter
	} else if last.Header.Interpreter != interpreter {
		return 0, fmt.Errorf("%s writer: head entry interpreter %s is not the pinned %s", op, last.Header.Interpreter, interpreter)
	}
	commit := makeOneObjectCommit(obj, world.Revision+1, last.Header.SemanticsEpoch, interpreter, head, world.LogHead, writtenBy)
	data, err := json.Marshal(commit)
	if err != nil {
		return 0, fmt.Errorf("%s writer: encode commit: %w", op, err)
	}
	status, body, err := c.doAuth(ctx, http.MethodPost, "/v1/commit", bytes.NewReader(data), "Bearer "+token)
	if err != nil {
		// The request may already have reached the daemon: the outcome is unknown.
		return 0, &store.UncertainError{Op: op + " commit", Cause: err}
	}
	var failure daemon.APIError
	_ = json.Unmarshal(body, &failure)
	if status == http.StatusServiceUnavailable && failure.Error.Class == "CommitUncertain" {
		return 0, &store.UncertainError{Op: op + " commit", Cause: fmt.Errorf("HTTP 503 CommitUncertain")}
	}
	if status != http.StatusOK {
		return 0, fmt.Errorf("%s writer: commit refused HTTP %d %s; not retried", op, status, failure.Error.Class)
	}
	var result struct {
		SelectedHead string `json:"selectedHead"`
	}
	if json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref {
		// A 200 that does not name the committed head may still have landed.
		return 0, &store.UncertainError{Op: op + " commit", Cause: fmt.Errorf("commit answered %q, not the committed %s", result.SelectedHead, commit.NextWorld.Ref)}
	}
	return commit.Entry.Header.EntryIndex, nil
}

// demoProve mints the proof of the committed source srcRef through the
// production producer and returns the envelope ref. The reader and writer are
// the caller's; the production wiring passes demoProofReader and demoProofWriter.
func demoProve(ctx context.Context, reader evidence.ObjectReader, writer evidence.ObjectWriter, h daemon.HealthResponse, key [32]byte, srcRef hashref.HashRef, required []string) (hashref.HashRef, error) {
	ctx, cancel := demoProofContext(ctx)
	defer cancel()
	checker, err := demoProofChecker(h)
	if err != nil {
		return hashref.HashRef{}, err
	}
	p, err := evidence.NewProducer(key, reader, writer, checker, demoProofBounds())
	if err != nil {
		return hashref.HashRef{}, err
	}
	return p.GenerateProof(ctx, srcRef, required)
}
