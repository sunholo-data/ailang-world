package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sunholo-data/ailang-world/host/agui"
	"github.com/sunholo-data/ailang-world/host/store"
)

const (
	aguiWriteMargin = 2 * time.Second
	aguiRunBudget   = writeTimeout - readDeadline - aguiWriteMargin
	aguiBodyBound   = 2 * time.Second
	aguiTick        = 250 * time.Millisecond
	aguiPage        = 100
	aguiCap         = 16
	aguiMaxBody     = 64 << 10
)

var errAGUISlowBody = errors.New("AGUI body exceeded its time bound")

// aguiBodyReader checks every return, including EOF/error. It does not relax
// the server's frozen transport deadlines or start a goroutine around Read.
type aguiBodyReader struct {
	io.Reader
	until time.Time
}

func (b aguiBodyReader) Read(p []byte) (int, error) {
	n, err := b.Reader.Read(p)
	if time.Now().After(b.until) {
		return 0, errAGUISlowBody
	}
	return n, err
}

type aguiInput struct {
	ThreadID *string         `json:"threadId"`
	RunID    *string         `json:"runId"`
	Messages json.RawMessage `json:"messages"`
	State    json.RawMessage `json:"state"`
}

func aguiCursor(state json.RawMessage, header string) (int64, error) {
	after := int64(-1)
	explicit := false
	if len(state) > 0 {
		var s struct {
			Schema    string          `json:"schema"`
			LastIndex json.RawMessage `json:"lastIndex"`
		}
		// AG-UI state is arbitrary JSON. Only our recognized object schema
		// supplies an explicit cursor; all other initial state is genesis.
		err := json.Unmarshal(state, &s)
		if err == nil && s.Schema == agui.StateSchema {
			explicit = true
			if err := json.Unmarshal(s.LastIndex, &after); err != nil || string(s.LastIndex) == "null" || after < -1 {
				return 0, errors.New("invalid state.lastIndex")
			}
		}
	}
	if header != "" {
		h, err := strconv.ParseInt(header, 10, 64)
		if err != nil || h < -1 {
			return 0, errors.New("invalid Last-Event-ID")
		}
		if explicit && after != h {
			return 0, errors.New("conflicting resume cursors")
		}
		after = h
	}
	return after, nil
}
func aguiEntry(e store.LogEntry) agui.Entry {
	b, err := json.Marshal(logJSON(e))
	if err != nil {
		panic(err)
	}
	return agui.Entry{Index: e.Header.EntryIndex, EntryHash: e.EntryHash.String(), Value: b}
}

// handleAGUI serves one bounded, resumable read-only AG-UI 1.0 run.
func (d *Daemon) handleAGUI(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	body, err := io.ReadAll(aguiBodyReader{http.MaxBytesReader(w, r.Body, aguiMaxBody), t0.Add(d.aguiBodyBound)})
	if err != nil {
		if errors.Is(err, errAGUISlowBody) {
			if time.Since(t0) <= 25*time.Second {
				writeAPIError(w, "SlowBody", "request body was too slow", http.StatusRequestTimeout)
			}
			return
		}
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeAPIError(w, "PayloadTooLarge", "request body exceeds 64 KiB", http.StatusRequestEntityTooLarge)
		} else {
			writeAPIError(w, "BadRequest", "invalid request body", http.StatusBadRequest)
		}
		return
	}
	var input aguiInput
	if err := json.Unmarshal(body, &input); err != nil || input.ThreadID == nil || input.RunID == nil || len(input.Messages) == 0 || input.Messages[0] != '[' {
		writeAPIError(w, "BadRequest", "threadId, runId and messages array are required", http.StatusBadRequest)
		return
	}
	after, err := aguiCursor(input.State, r.Header.Get("Last-Event-ID"))
	if err != nil {
		writeAPIError(w, "BadRequest", err.Error(), http.StatusBadRequest)
		return
	}
	select {
	case d.aguiSlots <- struct{}{}:
		defer func() { <-d.aguiSlots }()
	default:
		w.Header().Set("Retry-After", "1")
		writeAPIError(w, "StreamLimit", "too many live streams", http.StatusServiceUnavailable)
		return
	}
	// All reads inherit both disconnect and shutdown. Join the stop watcher on
	// every exit so neither completed nor refused runs leave a goroutine behind.
	runCtx, cancel := context.WithCancel(r.Context())
	watchDone := make(chan struct{})
	go func() {
		defer close(watchDone)
		select {
		case <-d.aguiStop:
			cancel()
		case <-runCtx.Done():
		}
	}()
	defer func() { cancel(); <-watchDone }()
	var snap *agui.Entry
	if after >= 0 {
		ctx, stop := context.WithTimeout(runCtx, d.readDeadline)
		e, ok, err := d.reads.GetLogEntry(ctx, after)
		stop()
		if err != nil {
			d.writeInternalError(w, r, err)
			return
		}
		if !ok {
			writeAPIError(w, "NotFound", "unknown cursor", http.StatusNotFound)
			return
		}
		entry := aguiEntry(e)
		snap = &entry
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeAPIError(w, "Internal", internalErrorMessage, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	emit := func(b []byte) bool {
		_, err := w.Write(b)
		if err != nil {
			return false
		}
		flusher.Flush()
		return true
	}
	run := agui.Run{ThreadID: *input.ThreadID, RunID: *input.RunID}
	if !emit(agui.Start(run, after, snap)) {
		return
	}
	end := t0.Add(d.aguiBudget)
	for time.Now().Before(end) && runCtx.Err() == nil {
		ctx, stop := context.WithTimeout(runCtx, d.readDeadline)
		page, err := d.reads.LogEntriesAfter(ctx, after, aguiPage)
		stop()
		if err != nil {
			if runCtx.Err() != nil {
				break
			}
			code := "Internal"
			if errors.Is(err, context.DeadlineExceeded) {
				code = "Timeout"
			}
			// Match the existing internal log path without writing a second HTTP body.
			fmt.Fprintf(d.errLog, "ailang-worldd: internal error: %s %s: %v\n", r.Method, r.URL.Path, err)
			emit(agui.Error(code, internalErrorMessage))
			return
		}
		for _, e := range page {
			if !emit(agui.EntryFrames(aguiEntry(e))) {
				return
			}
			after = e.Header.EntryIndex
		}
		if len(page) == aguiPage {
			continue
		}
		wait := d.aguiTick
		if remaining := time.Until(end); remaining < wait {
			wait = remaining
		}
		if wait <= 0 {
			break
		}
		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
		case <-runCtx.Done():
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
	}
	if r.Context().Err() == nil {
		emit(agui.Finished(run, after))
	}
}
