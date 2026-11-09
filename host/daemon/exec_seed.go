package daemon

import "github.com/sunholo-data/ailang-world/host/broker"

// Row 257 (w-exec-cache-seed-at-bind r3 §3.6, §4): the bound exec handler's
// live seed state for the episode report. Kept out of workspace.go on purpose.

// execEntryLocked is the one read of w.execH; the caller holds w.execMu.
func (w *workspaceTools) execEntryLocked(episodeID string) (episodeTool, bool) {
	t, ok := w.execH[episodeID]
	return t, ok
}

// boundExecHandler is the mutex-protected accessor: the episode's bound
// *broker.ExecHandler, or nil. It holds execMu only for the map read.
func (w *workspaceTools) boundExecHandler(episodeID string) *broker.ExecHandler {
	if w == nil {
		return nil
	}
	w.execMu.Lock()
	t, ok := w.execEntryLocked(episodeID)
	w.execMu.Unlock()
	if !ok {
		return nil
	}
	eh, _ := t.h.(*broker.ExecHandler)
	return eh
}

// execCacheStatus is the report's exec_cache: live daemon state, empty after
// a restart until the next bind. execMu and the seed's own mutex are never
// held together.
func (w *workspaceTools) execCacheStatus(episodeID string) []ReportExecCache {
	eh := w.boundExecHandler(episodeID)
	if eh == nil {
		return nil
	}
	var out []ReportExecCache
	for _, s := range eh.SeedStatus() {
		out = append(out, ReportExecCache{Project: s.Project, Cache: s.Cache, State: s.State, Attempts: s.Attempts,
			By: s.By, TookMS: s.Took.Milliseconds(), Failure: s.Failure})
	}
	return out
}
