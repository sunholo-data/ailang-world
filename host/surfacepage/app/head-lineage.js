// The node stack's forward strips, attached to the inspected face in a historical slice.
(() => {
  'use strict';
  function forwardPaths(model, id, cut) {
    const seen = new Set([id]);
    const index = id => model.timeIndex.get(id);
    function later(parent) {
      return [...new Set(model.kids[parent] || [])].filter(key =>
        model.N[key] && !model.N[key].draft && key !== parent && index(key) > cut &&
        index(key) >= index(parent) && index(key) <= model.presentIndex
      ).sort((a,b) => index(a) - index(b) || model.N[a].rank - model.N[b].rank);
    }
    // Separate immediate forks as in the node stack. Deeper strips follow only
    // an unambiguous continuation. Reconverging paths share one destination.
    const paths = later(id).slice(0, 3).map(key => [key]);
    for (const [key] of paths) seen.add(key);
    for (const path of paths) {
      while (path.length < 3) {
        const children = later(path.at(-1));
        if (children.length !== 1 || seen.has(children[0])) break;
        path.push(children[0]);
        seen.add(children[0]);
      }
    }
    return paths;
  }

  function create(space, {title, tag, path, accent}) {
    const stack = document.createElement('div');
    stack.className = 'hc-forward';
    stack.setAttribute('role', 'group');
    stack.setAttribute('aria-label', 'Recorded descendants after this time');
    space.append(stack);
    let shownModel, shownLayout, origin;
    const links = new Map();
    function update(model, layout, id, bounds) {
      if (model === shownModel && layout === shownLayout && id === origin) return;
      shownModel = model;
      shownLayout = layout;
      origin = id;
      links.clear();
      stack.replaceChildren();
      stack.hidden = true;
      if (!bounds || !layout?.time.stop || !model.N[id]) return;
      const paths = forwardPaths(model, id, layout.time.referenceIndex);
      if (!paths.length) return;
      const lanePosition = keys => layout.laneOrder.indexOf(model.N[keys[0]].lane);
      paths.sort((a,b) => lanePosition(a) - lanePosition(b));
      stack.hidden = false;
      stack.dataset.origin = id;
      stack.style.left = bounds.x + 'px';
      stack.style.top = (bounds.y + bounds.ry) + 'px';
      const maxWidth = Math.max(110, space.parentElement.clientWidth - 28);
      stack.style.width = Math.min(maxWidth, Math.max(150, (bounds.rx * 2 + 8), 116 * paths.length)) + 'px';
      for (const keys of paths) {
        const column = document.createElement('div');
        column.className = 'hc-forward-path';
        let parent = id;
        for (const key of keys) {
          const node = model.N[key], el = document.createElement('a');
          const label = title(node), identity = tag(node);
          el.href = path(node);
          el.dataset.lineageState = key;
          el.dataset.lineageFrom = parent;
          el.setAttribute('aria-label', 'Ahead · ' + [model.laneTitle(node.lane), identity, label].filter(Boolean).join(' · '));
          el.textContent = identity && identity !== label ? identity + ' · ' + label : label;
          el.style.setProperty('--hc-lane', accent(model.LANES[node.lane]));
          column.append(el);
          links.set(key, el);
          parent = key;
        }
        stack.append(column);
      }
    }
    return {update, link:id => links.get(id), get origin() { return origin; }};
  }
  if (typeof module !== 'undefined') module.exports = {forwardPaths};
  if (typeof window !== 'undefined') window.WorldHeadLineage = {create};
})();
