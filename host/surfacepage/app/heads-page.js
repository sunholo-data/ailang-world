// Refresh even when the graph generation stays put: session liveness and source
// freshness can change independently of commits. Only one read is in flight.
(() => {
  'use strict';
  const container = document.getElementById('heads-workspace');
  const dark = matchMedia('(prefers-color-scheme: dark)');
  const reduced = matchMedia('(prefers-reduced-motion: reduce)');
  function applyLook() {
    let storage;
    try { storage = localStorage; } catch (_) { /* browser preferences may be disabled */ }
    const shared = window.WorldSurface;
    const attrs = shared.resolveLook(shared.loadLook(storage), dark.matches, reduced.matches);
    document.documentElement.dataset.theme = attrs['data-theme'] === 'night' ? 'slate' : attrs['data-theme'];
    document.documentElement.dataset.motion = attrs['data-motion'];
  }
  applyLook();
  dark.addEventListener('change', applyLook);
  reduced.addEventListener('change', applyLook);
  window.addEventListener('storage', applyLook);

  let lastModel = null, timer, loading = false;
  async function refresh() {
    clearTimeout(timer);
    if (loading || document.hidden) return;
    loading = true;
    const controller = new AbortController();
    const deadline = setTimeout(() => controller.abort(), 12000);
    try {
      const response = await fetch('/surface/model.json', {cache:'no-store', signal:controller.signal});
      const data = await response.json();
      if (!response.ok || (!Array.isArray(data.lanes) && !data.empty)) throw new Error(data.reason || 'Model unavailable');
      window.WorldHeads.render(container, data);
      lastModel = data;
    } catch (error) {
      window.WorldHeads.render(container, lastModel, {stale:true, reason:error.message});
    } finally {
      clearTimeout(deadline);
      loading = false;
      if (!document.hidden) timer = setTimeout(refresh, 15000);
    }
  }
  document.addEventListener('visibilitychange', refresh);
  window.addEventListener('pageshow', event => { if (event.persisted) refresh(); });
  refresh();
})();
