// Branch Heads projects the timeline's recorded state-to-state edges.
(() => {
  'use strict';

  const projectionByModel = new WeakMap();
  const viewBySVG = new WeakMap();
  const esc = value => String(value ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;', '<':'&lt;', '>':'&gt;', '"':'&quot;', "'":'&#39;'}[c]));
  const laneName = lane => [lane.name || lane.key, lane.repo].filter(Boolean).join(' / ');
  const stateTitle = state => state.title || state.summary || state.label || state.id;
  const stateHref = ({lane, state}) => '/surface?' + new URLSearchParams({brane: lane.key, at: state.id}) + '#heads';
  const relation = ({from, to}) => from.state.kind === 'step' && to.state.kind === 'shared' ? 'Merge' : 'Base';

  function project(model) {
    if (model && projectionByModel.has(model)) return projectionByModel.get(model);
    const states = new Map(), tips = new Map(), pairs = new Map(), byLane = new Map();
    for (const lane of model?.lanes || []) {
      let tip;
      for (const state of lane.states) {
        states.set(state.id, {lane, state});
        if (!tip || state.rank > tip.rank) tip = state;
      }
      if (tip) tips.set(lane.key, tip.id);
    }
    for (const {from: source, to: target} of model?.links || []) {
      const from = states.get(source), to = states.get(target);
      if (!from || !to || from.lane.key === to.lane.key) continue;
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
      pair.edges.push({from, to});
    }
    const projection = {pairs: [...pairs.values()], byLane, tips};
    if (model) projectionByModel.set(model, projection);
    return projection;
  }

  function stateLink(endpoint, tips) {
    const {lane, state} = endpoint;
    const position = state.id === tips.get(lane.key) ? 'Head' : 'Earlier state';
    return `<a href="${esc(stateHref(endpoint))}" title="${esc(state.id)}">${esc(stateTitle(state))}</a>
      <span>${esc(laneName(lane))} · ${position}${state.at ? ' · <time datetime="' + esc(state.at) + '">' + esc(state.at.replace('T', ' ')) + '</time>' : ''}</span>`;
  }

  function detailsMarkup(projection, entries, selected) {
    const {byLane, tips} = projection;
    const pairs = byLane.get(selected) || [];
    const anyOnPage = [...entries.keys()].some(key => byLane.has(key));
    if (!anyOnPage && !pairs.length) return '';
    const legend = '<span class="wh-link-key">Lines connect branches with recorded links. Dashed lines include earlier states or multiple links.</span>';
    if (!pairs.length) {
      return legend + '<span class="wh-link-note">' + (selected && entries.has(selected)
        ? 'No recorded links to other branches for this head.'
        : 'Select a head to inspect its links.') + '</span>';
    }
    let count = 0, offPage = 0;
    const items = [];
    for (const pair of pairs) {
      const peer = pair.lanes.find(key => key !== selected);
      const away = !entries.has(peer);
      count += pair.edges.length;
      if (away) offPage += pair.edges.length;
      for (const edge of pair.edges) {
        items.push(`<li>
          <strong>${relation(edge)}<span>${away ? ' · Other page' : ''}</span></strong>
          <div class="wh-link-states">
            <div><small>From</small>${stateLink(edge.from, tips)}</div>
            <span class="wh-link-arrow" aria-hidden="true">→</span>
            <div><small>To</small>${stateLink(edge.to, tips)}</div>
          </div>
        </li>`);
      }
    }
    const label = laneName(pairs[0].edges[0].from.lane.key === selected ? pairs[0].edges[0].from.lane : pairs[0].edges[0].to.lane);
    return legend + `<details><summary>${count} recorded ${count === 1 ? 'link' : 'links'} · ${esc(label)}${offPage ? ' · ' + offPage + ' off page' : ''}</summary>
      <ul>${items.join('')}</ul></details>`;
  }

  function drawDetails(view, details, projection, entries, selected) {
    const page = [...entries.keys()].join('\n');
    if (view.detailProjection === projection && view.detailPage === page && view.detailSelected === selected) return;
    const markup = detailsMarkup(projection, entries, selected);
    if (view.markup !== markup) {
      const open = details.querySelector('details')?.open;
      const active = document.activeElement;
      const focused = details.contains(active);
      details.innerHTML = markup;
      const panel = details.querySelector('details');
      if (panel && view.detailSelected === selected) panel.open = Boolean(open);
      if (focused) {
        const replacement = [...details.querySelectorAll('a')].find(link => link.href === active.href);
        (replacement || panel?.querySelector('summary'))?.focus({preventScroll: true});
      }
      details.hidden = !markup;
      view.markup = markup;
    }
    Object.assign(view, {detailProjection: projection, detailPage: page, detailSelected: selected});
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
    const lines = [], highlighted = [];
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
      const active = lanes.includes(selected);
      const current = a.exact && b.exact && endpoints.every(({lane, state}) => projection.tips.get(lane.key) === state.id);
      const classes = 'wh-link' + (active ? ' is-active' : entries.has(selected) ? ' is-muted' : '') + (current ? '' : ' is-history');
      const port = ({x, y}, exact) => `<circle cx="${x}" cy="${y}" r="2.5" class="${exact ? 'is-exact' : ''}"/>`;
      const line = `<g class="${classes}"><path d="${d}"/>${port(source, a.exact)}${port(target, b.exact)}</g>`;
      (active ? highlighted : lines).push(line);
    }
    const markup = lines.concat(highlighted).join('');
    if (view.svgMarkup !== markup) {
      svg.innerHTML = markup;
      view.svgMarkup = markup;
    }
  }

  function render(svg, details, model, entries, selected) {
    let view = viewBySVG.get(svg);
    if (!view) {
      view = {};
      viewBySVG.set(svg, view);
    }
    const projection = project(model);
    drawDetails(view, details, projection, entries, selected);
    if (view.geometry?.projection !== projection || view.geometry.entries !== entries) {
      const visiblePairs = new Set(), visibleKeys = new Set();
      for (const key of entries.keys()) {
        for (const pair of projection.byLane.get(key) || []) {
          if (!pair.lanes.every(lane => entries.has(lane))) continue;
          visiblePairs.add(pair);
          pair.lanes.forEach(lane => visibleKeys.add(lane));
        }
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
  if (typeof module !== 'undefined') module.exports = {project, detailsMarkup, geometry};
})();
