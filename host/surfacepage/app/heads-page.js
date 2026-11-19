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
    const brane = selectedCard()?.dataset.head;
    timeline.href = '/surface' + (brane ? '#brane=' + encodeURIComponent(brane) : '');
  }
  timeline.addEventListener('pointerdown', updateTimelineLink);
  timeline.addEventListener('focus', updateTimelineLink);
  timeline.addEventListener('click', event => {
    updateTimelineLink();
    const selected = selectedCard();
    if (selected && projection?.plainClick(event)) {
      projection.save(timeline.href, selected.dataset.head, selected.querySelector('.wh-grip .wh-glyph'),
        selected.querySelector('.wh-face'), selected.style.getPropertyValue('--wh-base'), false);
    }
  });
  container.addEventListener('click', event => {
    if (!projection?.plainClick(event)) return;
    const link = event.target.closest('a[data-head-state]');
    if (!link) return;
    const target = new URL(link.href);
    const brane = target.searchParams.get('brane');
    const card = link.closest('.wh-card') || selectedCard();
    const glyph = card?.querySelector('.wh-grip .wh-glyph');
    const face = link.closest('.wh-rear') || link.closest('.wh-face') || link;
    projection.save(link.href, JSON.stringify([brane, link.dataset.headState]), glyph || link, face,
      card?.style.getPropertyValue('--wh-base'), false);
  }, {capture:true});
  container.addEventListener('click', updateTimelineLink);
  container.addEventListener('change', updateTimelineLink);

  function showProjection() {
    updateTimelineLink();
    if (projected || !incoming) return;
    const selected = selectedCard();
    if (!selected) return;
    projected = true;
    projection.animate(incoming, selected.dataset.head, selected.querySelector('.wh-grip .wh-glyph'),
      selected.querySelector('.wh-face'), false);
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
