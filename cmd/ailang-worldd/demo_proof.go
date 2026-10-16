package main

// M3a (row 150): the first production caller of evidence.NewProducer. The
// adapters below speak only to the running daemon's REST API: GET
// /v1/objects/{ref}?payload=true to read and POST /v1/commit to write. They
// hold no store handle, and the seed does not call them until M3b.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	head, has, err := headState(ctx, w.c)
	if err != nil {
		return fmt.Errorf("proof writer: read head: %w", err)
	}
	if !has {
		return fmt.Errorf("proof writer: no selected head to commit onto")
	}
	r := worldReader{ctx: ctx, c: w.c}
	world, found, err := r.world(head)
	if err != nil {
		return fmt.Errorf("proof writer: read world %s: %w", head, err)
	}
	if !found {
		return fmt.Errorf("proof writer: selected head %s has no world row", head)
	}
	last, found, err := r.entry(world.Revision)
	if err != nil {
		return fmt.Errorf("proof writer: read head entry: %w", err)
	}
	if !found {
		return fmt.Errorf("proof writer: head entry %d is absent", world.Revision)
	}
	if last.Header.Interpreter != w.interpreter {
		return fmt.Errorf("proof writer: head entry interpreter %s is not the pinned %s", last.Header.Interpreter, w.interpreter)
	}
	obj := demoObject{Hash: o.Hash.String(), InterfaceHash: o.InterfaceHash.String(), SemanticID: o.SemanticID, Provenance: o.Provenance, Payload: o.Payload}
	commit := demoMakeCommit(obj, world.Revision+1, last.Header.SemanticsEpoch, w.interpreter, head, world.LogHead)
	data, err := json.Marshal(commit)
	if err != nil {
		return fmt.Errorf("proof writer: encode commit: %w", err)
	}
	status, body, err := w.c.doAuth(ctx, http.MethodPost, "/v1/commit", bytes.NewReader(data), "Bearer "+w.token)
	if err != nil {
		// The request may already have reached the daemon: the outcome is unknown.
		return &store.UncertainError{Op: "proof commit", Cause: err}
	}
	var failure daemon.APIError
	_ = json.Unmarshal(body, &failure)
	if status == http.StatusServiceUnavailable && failure.Error.Class == "CommitUncertain" {
		return &store.UncertainError{Op: "proof commit", Cause: fmt.Errorf("HTTP 503 CommitUncertain")}
	}
	if status != http.StatusOK {
		return fmt.Errorf("proof writer: commit refused HTTP %d %s; not retried", status, failure.Error.Class)
	}
	var result struct {
		SelectedHead string `json:"selectedHead"`
	}
	if json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref {
		// A 200 that does not name the committed head may still have landed.
		return &store.UncertainError{Op: "proof commit", Cause: fmt.Errorf("commit answered %q, not the committed %s", result.SelectedHead, commit.NextWorld.Ref)}
	}
	return nil
}

// demoProve mints the proof of the committed source srcRef through the
// production producer and returns the envelope ref. The reader and writer are
// the caller's; the production wiring passes demoProofReader and demoProofWriter.
func demoProve(ctx context.Context, reader evidence.ObjectReader, writer evidence.ObjectWriter, h daemon.HealthResponse, key [32]byte, srcRef hashref.HashRef, required []string) (hashref.HashRef, error) {
	ctx, cancel := context.WithTimeout(ctx, demoProofWholeBound)
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
