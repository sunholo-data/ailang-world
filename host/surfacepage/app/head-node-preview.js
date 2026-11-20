// A lightweight preview of the same records that the Heads links open.
(() => {
  'use strict';
  if (typeof window === 'undefined') return;

  const exactTime = new Intl.DateTimeFormat(undefined, {
    dateStyle:'medium', timeStyle:'long', timeZone:'UTC',
  });
  let nextID = 0;

  function create(scene, options) {
    const tip = scene.querySelector('.hc-inspect-tip');
    tip.id = 'hc-node-preview-' + (++nextID);
    tip.removeAttribute('aria-hidden');
    tip.setAttribute('role', 'tooltip');

    let source, trigger, shownModel, shownID, dismissed, lastPointer;
    const events = new AbortController();
    const title = node => options.title?.(node) || node.title || node.episode || node.label || node.id;
    function text(tag, className, value) {
      const el = document.createElement(tag);
      el.className = className;
      el.textContent = value;
      return el;
    }

    function describe(nodes, ids, label) {
      const titles = [...new Set(ids || [])].map(id => nodes[id]).filter(Boolean).map(title);
      if (!titles.length) return;
      const list = document.createElement('div');
      list.className = 'hc-preview-relationship';
      list.append(text('b', '', label + ' '), document.createTextNode(titles.slice(0, 2).join(' · ')));
      if (titles.length > 2) list.append(document.createTextNode(' · +' + (titles.length - 2) + ' more'));
      tip.append(list);
    }

    function render(model, id) {
      const node = model.N[id];
      const lane = model.laneTitle(node.lane);
      const laneColour = model.LANES[node.lane]?.colour;
      if (/^#[0-9a-f]{3}([0-9a-f]{3})?$/i.test(laneColour)) tip.style.setProperty('--hc-lane', laneColour);
      else tip.style.removeProperty('--hc-lane');
      const status = node.draft ? 'Draft' : node.kind === 'shared' ? 'Merged into the shared world'
        : node.open?.label || (node.sessionLive ? 'Session live' : node.status || (node.kind === 'step' ? 'Step' : ''));
      const context = [lane, status !== lane && status].filter(Boolean).join(' · ');
      const meta = text('div', 'hc-preview-meta', context);
      const stamp = Date.parse(node.at);
      if (Number.isFinite(stamp)) {
        if (context) meta.append(document.createTextNode(' · '));
        const time = text('time', '', exactTime.format(stamp));
        time.dateTime = node.at;
        meta.append(time);
      }
      const heading = title(node);
      tip.replaceChildren(text('strong', 'hc-preview-title', heading), meta);
      const episode = node.episode && !heading.toLowerCase().includes(node.episode.toLowerCase()) ? node.episode : '';
      const agent = options.agent?.(node) || node.agent;
      const identity = [episode, node.pr && 'PR #' + node.pr, agent].filter(Boolean).join(' · ');
      if (identity) tip.append(text('div', 'hc-preview-meta', identity));
      const summary = options.summary?.(node) ?? (node.summary || '').trim();
      if (summary && summary.toLowerCase() !== heading.toLowerCase()) tip.append(text('p', 'hc-preview-summary', summary));
      describe(model.N, model.parents[id], 'From');
      describe(model.N, model.kids[id], 'Leads to');
    }

    function setTrigger(next) {
      if (trigger === next) return;
      if (trigger) {
        const ids = (trigger.getAttribute('aria-describedby') || '').split(/\s+/).filter(id => id && id !== tip.id);
        if (ids.length) trigger.setAttribute('aria-describedby', ids.join(' '));
        else trigger.removeAttribute('aria-describedby');
      }
      trigger = next;
      if (trigger) {
        const ids = (trigger.getAttribute('aria-describedby') || '').split(/\s+/).filter(Boolean);
        trigger.setAttribute('aria-describedby', [...new Set([...ids, tip.id])].join(' '));
      }
    }

    function hide() {
      tip.hidden = true;
      setTrigger(null);
      source = null;
      dismissed = null;
    }

    function place() {
      const box = scene.getBoundingClientRect();
      const rect = source.getBoundingClientRect();
      const viewport = document.documentElement;
      const margin = 8;
      const left = Math.max(0, box.left), top = Math.max(0, box.top);
      const right = Math.min(viewport.clientWidth, box.right), bottom = Math.min(viewport.clientHeight, box.bottom);
      if (right - left < margin * 2 || bottom - top < margin * 2) return false;
      if (!lastPointer && (rect.right < left || rect.left > right || rect.bottom < top || rect.top > bottom)) return false;
      tip.style.maxWidth = Math.min(340, right - left - margin * 2) + 'px';
      tip.style.maxHeight = bottom - top - margin * 2 + 'px';
      const x = lastPointer ? lastPointer.clientX + 12 : rect.left + rect.width / 2;
      const below = lastPointer ? lastPointer.clientY + 14 : rect.bottom + margin;
      const above = (lastPointer ? lastPointer.clientY : rect.top) - tip.offsetHeight - margin;
      const y = below + tip.offsetHeight <= bottom - margin || above < top + margin ? below : above;
      tip.style.left = Math.max(left + margin, Math.min(x, right - tip.offsetWidth - margin)) - box.left + 'px';
      tip.style.top = Math.max(top + margin, Math.min(y, bottom - tip.offsetHeight - margin)) - box.top + 'px';
      return true;
    }

    function show(model, id, anchor, pointer) {
      if (!model?.N[id] || !anchor?.isConnected) return hide();
      if (anchor === dismissed && id === shownID) return;
      dismissed = null;
      source = anchor;
      lastPointer = pointer;
      if (model !== shownModel || id !== shownID) render(model, id);
      shownModel = model;
      shownID = id;
      tip.hidden = false;
      if (!place()) return hide();
      setTrigger(anchor);
    }

    window.addEventListener('keydown', event => {
      if (event.key !== 'Escape' || tip.hidden || event.defaultPrevented ||
          event.target.closest?.('dialog, [role="dialog"], input, textarea, [contenteditable="true"]')) return;
      dismissed = source;
      tip.hidden = true;
      setTrigger(null);
      event.preventDefault();
      event.stopPropagation();
    }, {capture:true, signal:events.signal});

    return {element:tip, show, hide, destroy() { hide(); events.abort(); }};
  }

  window.WorldHeadNodePreview = {create};
})();
