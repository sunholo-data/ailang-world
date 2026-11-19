// Branch Heads projects the timeline's recorded state-to-state edges.
(() => {
  'use strict';

  const brane = typeof WorldBrane !== 'undefined' ? WorldBrane : require('./brane.js');
  const projectionByModel = new WeakMap();
  const viewBySVG = new WeakMap();
  const esc = value => String(value ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;', '<':'&lt;', '>':'&gt;', '"':'&quot;', "'":'&#39;'}[c]));
  const laneName = lane => [brane.label(lane.name || lane.key, lane.row) || lane.name || lane.key, lane.repo].filter(Boolean).join(' / ');
  const stateTitle = state => state.title || state.summary || state.label || state.id;
  const stateHref = ({lane, state}) => '/surface?' + new URLSearchParams({brane: lane.key, at: state.id}) + '#heads';
  const timelineHref = ({lane, state}) => brane.timelineHref(lane.key, state.id);
  const relation = ({from, to}) => {
    if (from.lane.key === to.lane.key) return 'Continues';
    if (from.state.kind === 'step' && to.state.kind === 'shared') return 'Merge';
    return from.state.kind === 'shared' ? 'Fork' : 'Base';
  };
  const colour = lane => /^#[0-9a-f]{3}([0-9a-f]{3})?$/i.test(lane?.colour) ? lane.colour : '#657b91';
  const glyph = lane => brane.glyph(lane.glyph, colour(lane), 10);
  const date = new Intl.DateTimeFormat(undefined, {month:'short', day:'numeric', year:'numeric', timeZone:'UTC'});

  function project(model) {
    if (model && projectionByModel.has(model)) return projectionByModel.get(model);
    const lanes = new Map(), states = new Map(), tips = new Map(), pairs = new Map(), byLane = new Map(), byState = new Map();
    for (const lane of model?.lanes || []) {
      lanes.set(lane.key, lane);
      let tip;
      for (const state of lane.states) {
        states.set(state.id, {lane, state});
        if (!tip || state.rank > tip.rank) tip = state;
      }
      if (tip) tips.set(lane.key, tip.id);
    }
    for (const {from: source, to: target} of brane.graphEdges(model || {})) {
      const from = states.get(source), to = states.get(target);
      if (!from || !to) continue;
      const edge = {from, to};
      for (const id of new Set([source, target])) {
        if (!byState.has(id)) byState.set(id, []);
        byState.get(id).push(edge);
      }
      if (from.lane.key === to.lane.key) continue;
      const lanes = [from.lane.key, to.lane.key].sort();
      const key = JSON.stringify(lanes);
      let pair = pairs.get(key);
      if (!pair) {
        pair = {lanes, edges: []};
        pairs.set(key, pair);
        for (const lane of lanes) {
          if (!byLane.has(lane)) byLane.set(lane, []);
          byLane.get(lane).push(pair);
        }
      }
      pair.edges.push(edge);
    }
    const projection = {lanes, states, pairs: [...pairs.values()], byLane, byState, tips};
    if (model) projectionByModel.set(model, projection);
    return projection;
  }

  function stateLink(endpoint, tips) {
    const {lane, state} = endpoint;
    const position = state.id === tips.get(lane.key) ? 'Head' : 'Earlier state';
    return `<div class="wh-evidence-state">
      <div>${esc(laneName(lane))} · ${position}</div>
      <a data-head-state="${esc(state.id)}" href="${esc(stateHref(endpoint))}">${esc(stateTitle(state))}</a>
      <a data-head-timeline="${esc(state.id)}" href="${esc(timelineHref(endpoint))}" aria-label="View ${esc(stateTitle(state))} on the timeline">Timeline</a>
    </div>`;
  }

  function neighbors(projection, selected) {
    const peers = new Map();
    for (const {edges} of projection.byLane.get(selected) || []) {
      const {from, to} = edges[0];
      const lane = from.lane.key === selected ? to.lane : from.lane;
      peers.set(lane.key, {lane, edges});
    }
    return peers;
  }

  function selectedEndpoint(projection, selected, requested) {
    const endpoint = projection.states.get(requested);
    return endpoint?.lane.key === selected ? endpoint : projection.states.get(projection.tips.get(selected));
  }

  function stationsFor(projection, peers, selected, requested) {
    const current = selectedEndpoint(projection, selected, requested);
    if (!current) return [];
    const byState = new Map([[current.state.id, current.state]]);
    for (const {edges} of peers.values()) {
      for (const {from, to} of edges) {
        const {state} = to.lane.key === selected ? to : from;
        byState.set(state.id, state);
      }
    }
    const earlier = [...byState.values()].filter(state => state.id !== current.state.id).sort((a, b) =>
      (Date.parse(b.at) - Date.parse(a.at)) || b.rank - a.rank);
    return [current.state, ...earlier].map(state => {
      const edges = projection.byState.get(state.id) || [];
      return {
        state,
        incoming: edges.filter(edge => edge.to.state.id === state.id),
        outgoing: edges.filter(edge => edge.from.state.id === state.id),
      };
    });
  }

  function topologyEdge(edge, incoming, entries) {
    const endpoint = incoming ? edge.from : edge.to;
    const {lane, state} = endpoint;
    const sameLane = edge.from.lane.key === edge.to.lane.key;
    const name = sameLane ? stateTitle(state) : laneName(lane);
    const type = relation(edge);
    const offPage = entries.has(lane.key) ? '' : ' · Other page';
    const content = `${glyph(lane)}<span>${esc(sameLane ? name : lane.name || lane.key)}</span>`;
    const control = sameLane
      ? `<a href="${esc(timelineHref(endpoint))}" data-head-timeline="${esc(state.id)}"
          aria-label="View ${esc(name)} on the timeline" title="${esc(name)}">${content}</a>`
      : `<button type="button" data-locate="${esc(lane.key)}" title="${esc(name + offPage)}"
          aria-label="${esc('Show branch ' + name + offPage)}">${content}</button>`;
    return `<div class="wh-topology-edge" data-relation="${type.toLowerCase()}">
      ${control}<span class="wh-edge-direction" title="${esc(stateTitle(state))}">${type}<span aria-hidden="true"> →</span></span>
    </div>`;
  }

  function detailsMarkup(projection, entries, selected, requested, peers = neighbors(projection, selected)) {
    const {tips} = projection;
    const lane = projection.lanes.get(selected);
    if (!lane || !tips.has(selected)) return '';
    const stations = stationsFor(projection, peers, selected, requested);
    const rows = stations.map(({state, incoming, outgoing}, index) => {
      const parsed = Date.parse(state.at);
      const at = Number.isFinite(parsed) ? `<time datetime="${esc(state.at)}">${esc(date.format(parsed))}</time>` : '';
      const id = 'wh-evidence-' + index;
      return `<li class="wh-topology-station">
        <div class="wh-topology-in">${incoming.map(edge => topologyEdge(edge, true, entries)).join('')}</div>
        <div class="wh-topology-node">
          <button type="button" data-topology-state="${esc(state.id)}" aria-controls="${id}" aria-expanded="false"
            aria-label="Evidence for ${esc(stateTitle(state))}" title="${esc(stateTitle(state))}">${glyph(lane)}${at}</button>
        </div>
        <div class="wh-topology-out">${outgoing.map(edge => topologyEdge(edge, false, entries)).join('')}</div>
        <div class="wh-topology-evidence" id="${id}" hidden>
          ${incoming.length || outgoing.length ? '' : '<p class="wh-link-note">No recorded links for this state.</p>'}
          ${[...incoming, ...outgoing].map(edge => `<div class="wh-evidence-edge">
            <strong>${relation(edge)}</strong>${stateLink(edge.from, tips)}
            <span aria-hidden="true">→</span>${stateLink(edge.to, tips)}
          </div>`).join('')}
        </div>
      </li>`;
    });
    const offPage = [...peers.keys()].filter(key => !entries.has(key)).length;
    const label = brane.label(lane.name || lane.key, lane.row);
    const pickedTip = stations[0].state.id === tips.get(selected);
    const earlier = rows.length > 1 ? `<details class="wh-topology-earlier">
      <summary>${rows.length - 1} ${pickedTip ? 'earlier' : 'other'} linked ${rows.length === 2 ? 'state' : 'states'}</summary>
      <ol class="wh-topology" start="2">${rows.slice(1).join('')}</ol></details>` : '';
    return `<section class="wh-neighborhood" aria-label="Recorded connections">
      <header><span class="wh-topology-title">${glyph(lane)}<strong>${esc(label)}</strong></span>
        <span>${peers.size ? peers.size + ' linked ' + (peers.size === 1 ? 'branch' : 'branches') : 'No branch links'}${offPage ? ' · ' + offPage + ' off page' : ''}</span></header>
      <div class="wh-topology-panel"><ol class="wh-topology">${rows[0]}</ol>${earlier}</div></section>`;
  }

  function toggleEvidence(control, open) {
    const evidence = control.closest('.wh-topology-station').querySelector('.wh-topology-evidence');
    evidence.hidden = !open;
    control.setAttribute('aria-expanded', String(open));
  }

  function drawDetails(view, details, projection, entries, selected, requested) {
    if (view.detailProjection === projection && view.detailEntries === entries && view.detailSelected === selected && view.detailState === requested) return;
    const peers = neighbors(projection, selected);
    const lane = projection.lanes.get(selected);
    const root = details.closest('.world-heads');
    if (lane) root?.style.setProperty('--wh-selection', colour(lane));
    else root?.style.removeProperty('--wh-selection');
    for (const [key, {el}] of entries) {
      el.classList.toggle('has-wh-link', peers.has(key));
      el.classList.toggle('wh-unrelated', peers.size > 0 && key !== selected && !peers.has(key));
    }
    Object.assign(view, {detailProjection: projection, detailEntries: entries, detailSelected: selected, detailState: requested});
    const markup = detailsMarkup(projection, entries, selected, requested, peers);
    if (view.markup === markup) return;
    const oldPanel = details.querySelector('.wh-neighborhood');
    const scroll = oldPanel?.querySelector('.wh-topology-panel')?.scrollTop;
    const earlierOpen = oldPanel?.querySelector('.wh-topology-earlier')?.open;
    const openStates = new Set([...details.querySelectorAll('[data-topology-state][aria-expanded="true"]')].map(node => node.dataset.topologyState));
    const active = document.activeElement;
    const focused = details.contains(active);
    details.innerHTML = markup;
    const panel = details.querySelector('.wh-neighborhood');
    if (panel && oldPanel && view.neighborhoodSelected === selected && view.neighborhoodState === requested) {
      const earlier = panel.querySelector('.wh-topology-earlier');
      if (earlier) earlier.open = Boolean(earlierOpen);
      for (const control of panel.querySelectorAll('[data-topology-state]')) {
        if (openStates.has(control.dataset.topologyState)) toggleEvidence(control, true);
      }
      panel.querySelector('.wh-topology-panel').scrollTop = scroll || 0;
    }
    if (focused) {
      const replacement = [...details.querySelectorAll('a, button')].find(control =>
        (active.tagName === 'A' && control.href === active.href) ||
        (active.dataset.locate && control.dataset.locate === active.dataset.locate) ||
        (active.dataset.topologyState && control.dataset.topologyState === active.dataset.topologyState));
      (replacement || panel?.querySelector('summary'))?.focus({preventScroll: true});
    }
    details.hidden = !markup;
    Object.assign(view, {markup, neighborhoodSelected: selected, neighborhoodState: requested});
  }

  function positions(entries, keys, bounds) {
    const boxes = new Map();
    for (const key of keys) {
      const {el} = entries.get(key);
      const getBox = node => {
        const rect = node.getBoundingClientRect();
        return {x: rect.left - bounds.left, y: rect.top - bounds.top, width: rect.width, height: rect.height};
      };
      const panes = new Map();
      for (const pane of el.querySelectorAll(':is(.wh-open, .wh-rear)[data-state]')) {
        if (pane.getClientRects().length) panes.set(pane.dataset.state, getBox(pane.closest('.wh-face') || pane));
      }
      const folded = el.querySelector('.wh-restore');
      boxes.set(key, {card: getBox(folded?.getClientRects().length ? folded : el), panes});
    }
    return boxes;
  }

  function port(box, horizontal, sign, exact) {
    const gap = exact ? 2 : 5;
    return horizontal
      ? {x: box.x + (sign > 0 ? box.width + gap : -gap), y: box.y + box.height / 2}
      : {x: box.x + box.width / 2, y: box.y + (sign > 0 ? box.height + gap : -gap)};
  }

  function geometry(a, b) {
    const dx = b.x + b.width / 2 - a.x - a.width / 2;
    const dy = b.y + b.height / 2 - a.y - a.height / 2;
    const horizontal = Math.abs(dx) / Math.max(a.width, b.width, 1) >= Math.abs(dy) / Math.max(a.height, b.height, 1);
    const sign = (horizontal ? dx : dy) >= 0 ? 1 : -1;
    const source = port(a, horizontal, sign, a.exact);
    const target = port(b, horizontal, -sign, b.exact);
    const tension = Math.min(90, Math.max(12, Math.abs(horizontal ? target.x - source.x : target.y - source.y) * .45));
    const x = horizontal ? tension * sign : 0, y = horizontal ? 0 : tension * sign;
    return {
      source, target,
      d: `M${source.x},${source.y} C${source.x + x},${source.y + y} ${target.x - x},${target.y - y} ${target.x},${target.y}`
    };
  }

  function drawSVG(svg, view) {
    const {projection, entries, selected} = view.geometry;
    const bounds = svg.getBoundingClientRect();
    if (!bounds.width || !bounds.height) return;
    const boxes = positions(entries, view.visibleKeys, bounds);
    const lines = [];
    for (const {lanes, edges} of view.visiblePairs) {
      const [first, second] = lanes.map(lane => boxes.get(lane));
      const endpoints = edges.length === 1 ? [edges[0].from, edges[0].to] : [];
      const endpoint = (side, index) => {
        const item = endpoints.find(({lane}) => lane.key === lanes[index]);
        const pane = item && side.panes.get(item.state.id);
        return pane ? {...pane, exact: true} : side.card;
      };
      const a = endpoint(first, 0), b = endpoint(second, 1);
      const {source, target, d} = geometry(a, b);
      const current = a.exact && b.exact && endpoints.every(({lane, state}) => projection.tips.get(lane.key) === state.id);
      const classes = 'wh-link' + (selected ? ' is-active' : '') + (current ? '' : ' is-history');
      const port = ({x, y}, exact) => `<circle cx="${x}" cy="${y}" r="2.5" class="${exact ? 'is-exact' : ''}"/>`;
      const line = `<g class="${classes}"><path d="${d}"/>${port(source, a.exact)}${port(target, b.exact)}</g>`;
      lines.push(line);
    }
    const markup = lines.join('');
    if (view.svgMarkup !== markup) {
      svg.innerHTML = markup;
      view.svgMarkup = markup;
    }
  }

  function render(svg, details, model, entries, selected, requested) {
    let view = viewBySVG.get(svg);
    if (!view) {
      view = {};
      viewBySVG.set(svg, view);
      details.addEventListener('click', event => {
        const control = event.target.closest('[data-topology-state]');
        if (control) toggleEvidence(control, control.getAttribute('aria-expanded') !== 'true');
      });
    }
    const projection = project(model);
    drawDetails(view, details, projection, entries, selected, requested);
    if (view.geometry?.projection !== projection || view.geometry.entries !== entries || view.geometry.selected !== selected) {
      const visiblePairs = new Set(), visibleKeys = new Set();
      for (const pair of (selected ? projection.byLane.get(selected) || [] : projection.pairs)) {
        if (!pair.lanes.every(lane => entries.has(lane))) continue;
        visiblePairs.add(pair);
        pair.lanes.forEach(lane => visibleKeys.add(lane));
      }
      Object.assign(view, {visiblePairs, visibleKeys});
    }
    view.geometry = {projection, entries, selected};
    if (view.frame != null) return;
    view.frame = requestAnimationFrame(() => {
      view.frame = null;
      drawSVG(svg, view);
    });
  }

  if (typeof window !== 'undefined') window.WorldHeadLinks = {render};
  if (typeof module !== 'undefined') module.exports = {project, neighbors, stationsFor, detailsMarkup, geometry};
})();
