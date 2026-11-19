// Refresh even when the graph generation stays put: session liveness and source
// freshness can change independently of commits. Only one read is in flight.
(() => {
  'use strict';
  const container = document.getElementById('heads-workspace');
  const projection = window.WorldProjection;
  const incoming = projection?.take();
  const timeline = document.getElementById('timeline-link');
  let projected = false;

  const selectedCard = () => container.querySelector('.wh-card.is-selected');
  function updateTimelineLink() {
    const selected = selectedCard();
    const state = selected?.querySelector('.wh-rear[aria-current="location"], .wh-open[data-state]');
    timeline.href = selected ? WorldBrane.timelineHref(selected.dataset.head, state?.dataset.state) : '/surface';
  }
  function saveProjection(link, stateLink) {
    const target = new URL(link.href);
    const side = !target.searchParams.has('brane');
    const address = side ? new URLSearchParams(target.hash.slice(1)) : target.searchParams;
    const brane = address.get('brane'), at = address.get('at');
    if (!brane) return;
    const card = [...container.querySelectorAll('.wh-card')].find(el => el.dataset.head === brane);
    const pane = stateLink ? link.closest('.wh-rear, .wh-face') : card?.querySelector('.wh-rear[aria-current="location"]');
    const visible = el => el?.getClientRects().length;
    const anchor = pane || (stateLink ? link : card?.querySelector('.wh-face'));
    const ownGlyph = card?.querySelector('.wh-grip .g');
    const glyph = pane?.querySelector('.g') || (visible(ownGlyph) && (pane || !stateLink) ? ownGlyph : link.querySelector('.g') || link);
    projection.save(link.href, at ? JSON.stringify([brane, at]) : brane, glyph, anchor,
      card?.style.getPropertyValue('--wh-base'), false);
  }
  timeline.addEventListener('pointerdown', updateTimelineLink);
  timeline.addEventListener('focus', updateTimelineLink);
  timeline.addEventListener('click', event => {
    updateTimelineLink();
    if (projection?.plainClick(event)) saveProjection(timeline, false);
  });
  container.addEventListener('click', event => {
    if (!projection?.plainClick(event)) return;
    const link = event.target.closest('a[data-head-state], a[data-head-timeline]');
    if (link) saveProjection(link, true);
  }, {capture:true});
  container.addEventListener('click', updateTimelineLink);
  container.addEventListener('change', updateTimelineLink);

  function showProjection() {
    updateTimelineLink();
    if (projected || !incoming) return;
    const selected = selectedCard();
    if (!selected) return;
    projected = true;
    const pane = selected.querySelector('.wh-rear[aria-current="location"]') || selected.querySelector('.wh-face');
    projection.animate(incoming, selected.dataset.head, pane.querySelector('.g'), pane, false);
  }

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
      showProjection();
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
