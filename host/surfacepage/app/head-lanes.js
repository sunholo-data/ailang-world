// Reuse Timeline's lane identities in the Heads projection.
(() => {
  'use strict';
  if (typeof window === 'undefined') return;

  function create(scene, chooseLane) {
    const strip = document.createElement('nav');
    strip.className = 'hc-lanes';
    strip.setAttribute('aria-label', 'Lanes');
    strip.hidden = true;
    scene.append(strip);

    let priorGraph, priorOrder = '', priorLane = '', priorWidth = 0;
    const buttons = new Map();
    strip.addEventListener('click', event => {
      const button = event.target.closest('button[data-camera-lane]');
      if (button) chooseLane(button.dataset.cameraLane);
    });
    strip.addEventListener('wheel', event => event.stopPropagation());

    function populate(graph, order) {
      const focused = document.activeElement?.dataset.cameraLane;
      const row = document.createDocumentFragment();
      buttons.clear();
      let section;
      for (const lane of order) {
        const group = graph.groupOf(lane);
        const button = document.createElement('button');
        button.type = 'button';
        button.dataset.cameraLane = lane;
        button.classList.toggle('hc-lane-group-start', section !== undefined && section !== group.id);
        button.classList.toggle('hc-lane-shared', lane === group.shared);
        section = group.id;
        const glyph = document.createElement('span');
        glyph.setAttribute('aria-hidden', 'true');
        glyph.innerHTML = graph.glyphOf(lane, 11);
        const name = document.createElement('span');
        name.className = 'hc-lane-name';
        name.textContent = graph.nameOf(lane);
        button.append(glyph, name);
        row.append(button);
        buttons.set(lane, button);
      }
      strip.replaceChildren(row);
      buttons.get(focused)?.focus({preventScroll:true});
    }

    function update({active, graph, current, heads, historical, order = []}) {
      const wasHidden = strip.hidden;
      strip.hidden = !active || !heads.length;
      if (strip.hidden) return;

      const byLane = new Map(heads.map(head => [head.lane, head]));
      const sidebarOrder = graph.ROWS.flatMap(row => graph.QUIET[row] || [row]);
      const lanes = [...new Set([...order, ...sidebarOrder, ...byLane.keys()])].filter(lane => byLane.has(lane));
      const signature = JSON.stringify(lanes);
      if (priorGraph !== graph || priorOrder !== signature) populate(graph, lanes);

      for (const [lane, button] of buttons) {
        const head = byLane.get(lane);
        const needs = !historical && head.attention;
        const live = !historical && head.live;
        const label = [graph.laneTitle(lane), live && 'Running', needs && 'Needs you'].filter(Boolean).join(', ');
        if (button.title !== label) {
          button.title = label;
          button.setAttribute('aria-label', label);
        }
        button.toggleAttribute('data-running', live);
        button.toggleAttribute('data-needs-attention', needs);
        if (lane === current) button.setAttribute('aria-current', 'true');
        else button.removeAttribute('aria-current');
      }

      const width = strip.clientWidth;
      if (priorOrder !== signature || wasHidden || priorLane !== current || priorWidth !== width) {
        buttons.get(current)?.scrollIntoView({block:'nearest', inline:'nearest', behavior:'instant'});
      }
      priorGraph = graph;
      priorOrder = signature;
      priorLane = current;
      priorWidth = width;
    }

    return {update};
  }

  window.WorldHeadLanes = {create};
})();
