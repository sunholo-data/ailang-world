package daemon

// Capsule compile-cache template status, rebuild and GC (row 153 M3b).
//
// Publication builds one template per (interpreter, source) pair (host/archive
// CheckSource). A store published before that, or a deleted cache directory,
// leaves a descriptor without one: its first run compiles cold, which is
// correct but slow. So the daemon (a) says at startup, per descriptor, whether
// the template is ready or missing, (b) rebuilds the missing ones AFTER the
// listener is ready, at concurrency <= 2, each bounded by the check's own 10 s,
// without ever delaying serving, and (c) prunes templates of interpreters no
// descriptor pins any more.
//
// This file execs nothing (the subprocess-site gate in host/broker): the
// rebuild calls archive.CheckSource, the one place a template is built.

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// Template states reported by TemplateStatus.State.
const (
	TemplateReady   = "ready"
	TemplateMissing = "missing"
	TemplateRebuilt = "rebuilt"
	TemplateFailed  = "failed"
)

// templateRebuildConcurrency bounds how many missing templates rebuild at once
// (measured: 21.8 s serial, 10.4 s at 2, for the 9 se-tools under load).
const templateRebuildConcurrency = 2

// templateStopWait bounds how long Close/Run wait for the maintenance
// goroutine after cancelling it (an in-flight check is killed by its context).
const templateStopWait = 12 * time.Second

// TemplateStatus is one descriptor's compile-cache template state.
type TemplateStatus struct {
	ID                  string
	Interpreter, Source hashref.HashRef
	State               string // ready | missing | rebuilt | failed
}

type templateEntry struct {
	TemplateStatus
	source []byte
}

// templateBook is the daemon's template bookkeeping, behind one mutex.
type templateBook struct {
	mu      sync.Mutex
	entries []templateEntry
	keep    map[string]bool // interpreter digests GC must keep; nil = GC disabled
	cancel  context.CancelFunc
	done    chan struct{}
}

// CapsuleTemplates returns a copy of the per-descriptor template status read
// at startup and updated by the rebuild.
func (d *Daemon) CapsuleTemplates() []TemplateStatus {
	d.templates.mu.Lock()
	defer d.templates.mu.Unlock()
	out := make([]TemplateStatus, len(d.templates.entries))
	for i, e := range d.templates.entries {
		out[i] = e.TemplateStatus
	}
	return out
}

func (d *Daemon) setTemplateState(id, state string) {
	d.templates.mu.Lock()
	defer d.templates.mu.Unlock()
	for i := range d.templates.entries {
		if d.templates.entries[i].ID == id {
			d.templates.entries[i].State = state
		}
	}
}

func oneLine(s string) string {
	return strings.NewReplacer("\r", "\\r", "\n", "\\n").Replace(s)
}

// reportTemplates is the startup half: one stat per descriptor and one
// operator line each. It never builds anything and never fails startup; a head
// it cannot read is reported and leaves GC disabled (it cannot know what to
// keep).
func (d *Daemon) reportTemplates(ctx context.Context) {
	if d.arch == nil {
		return
	}
	rd := transitionreg.NewReader(d.store)
	_, rev, found, err := rd.CurrentRevision(ctx)
	if err != nil {
		fmt.Fprintf(d.errLog, "ailang-worldd: capsule template status unavailable: %s\n", oneLine(err.Error()))
		return
	}
	keep := map[string]bool{d.interpreterHash.Digest(): true}
	gcSafe := true
	var entries []templateEntry
	if found {
		for _, desc := range rev.Entries {
			keep[desc.Interpreter.Digest()] = true
			obj, ok, err := d.store.GetObject(ctx, desc.TransitionFn)
			if err != nil || !ok {
				fmt.Fprintf(d.errLog, "ailang-worldd: capsule template %s: source unreadable; status unavailable\n", desc.ID)
				gcSafe = false // an unreadable source: GC cannot know what to keep
				continue
			}
			state := TemplateMissing
			if d.arch.CapsuleTemplateReady(desc.Interpreter, obj.Payload) {
				state = TemplateReady
			}
			entries = append(entries, templateEntry{
				TemplateStatus: TemplateStatus{ID: desc.ID, Interpreter: desc.Interpreter, Source: hashref.SumSHA256(obj.Payload), State: state},
				source:         obj.Payload})
			fmt.Fprintf(d.errLog, "ailang-worldd: capsule template %s: %s\n", desc.ID, state)
		}
	}
	d.templates.mu.Lock()
	d.templates.entries = entries
	if gcSafe {
		d.templates.keep = keep
	}
	d.templates.mu.Unlock()
}

// startTemplateMaintenance starts the single GC + rebuild goroutine. Run calls
// it after the listener is bound and announced, so nothing here can delay
// serving. It is idempotent; stopTemplateMaintenance cancels and joins it.
func (d *Daemon) startTemplateMaintenance(ctx context.Context) {
	if d.arch == nil {
		return
	}
	d.templates.mu.Lock()
	defer d.templates.mu.Unlock()
	if d.templates.cancel != nil {
		return
	}
	mctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	d.templates.cancel, d.templates.done = cancel, done
	go func() {
		defer close(done)
		d.maintainTemplates(mctx)
	}()
}

// stopTemplateMaintenance cancels the goroutine and waits (bounded) for it.
func (d *Daemon) stopTemplateMaintenance() {
	d.templates.mu.Lock()
	cancel, done := d.templates.cancel, d.templates.done
	d.templates.mu.Unlock()
	if cancel == nil {
		return
	}
	cancel()
	select {
	case <-done:
	case <-time.After(templateStopWait):
		fmt.Fprintln(d.errLog, "ailang-worldd: capsule template maintenance did not stop within "+templateStopWait.String())
	}
}

func (d *Daemon) maintainTemplates(ctx context.Context) {
	d.templates.mu.Lock()
	keep := d.templates.keep
	var missing []templateEntry
	for _, e := range d.templates.entries {
		if e.State == TemplateMissing {
			missing = append(missing, e)
		}
	}
	d.templates.mu.Unlock()

	if keep != nil {
		removed, err := d.arch.PruneCapsuleCache(keep, archive.StaleTemplateTmpAge, time.Now())
		if err != nil {
			fmt.Fprintf(d.errLog, "ailang-worldd: capsule template gc: %s\n", oneLine(err.Error()))
		}
		if len(removed) > 0 {
			fmt.Fprintf(d.errLog, "ailang-worldd: capsule template gc: pruned %d: %s\n", len(removed), strings.Join(removed, " "))
		}
	}

	sem := make(chan struct{}, templateRebuildConcurrency)
	var wg sync.WaitGroup
	defer wg.Wait()
	for _, e := range missing {
		if gate := d.cfg.templateRebuildGate; gate != nil {
			select {
			case <-gate:
			case <-ctx.Done():
				return
			}
		}
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			return
		}
		wg.Add(1)
		go func(e templateEntry) {
			defer wg.Done()
			defer func() { <-sem }()
			d.rebuildTemplate(ctx, e)
		}(e)
	}
}

// rebuildTemplate builds one missing template with archive.CheckSource, which
// bounds the check at its own 10 s.
func (d *Daemon) rebuildTemplate(ctx context.Context, e templateEntry) {
	start := time.Now()
	res, err := d.arch.CheckSource(ctx, e.Interpreter, e.source)
	if err == nil && !res.Passed {
		err = fmt.Errorf("the interpreter refused the source: %s", res.Output)
	}
	if err != nil {
		if ctx.Err() != nil {
			return // shutting down: not a rebuild failure
		}
		d.setTemplateState(e.ID, TemplateFailed)
		fmt.Fprintf(d.errLog, "ailang-worldd: capsule template %s: rebuild failed: %s\n", e.ID, oneLine(err.Error()))
		return
	}
	d.setTemplateState(e.ID, TemplateRebuilt)
	fmt.Fprintf(d.errLog, "ailang-worldd: capsule template %s: rebuilt in %s\n", e.ID, time.Since(start).Round(time.Millisecond))
}
