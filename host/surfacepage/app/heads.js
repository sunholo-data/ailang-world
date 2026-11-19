// End-on view of the recorded branch lanes. The World node view owns state inspection.
(() => {
  'use strict';
  const PAGE_SIZE = 16;
  const HISTORY_PAGE = 40;
  const STORAGE_KEY = 'world-heads-worktable';
  const DAY = 24 * 60 * 60 * 1000;
  const defaultColour = '#657b91';
  const viewByContainer = new WeakMap();
  const timeFormat = new Intl.DateTimeFormat(undefined, {
    year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit', timeZone: 'UTC'
  });
  const rearTimeFormat = new Intl.DateTimeFormat(undefined, {
    month: 'short', day: 'numeric', timeZone: 'UTC'
  });
  const {glyph, label: braneLabel, sharedInk} = WorldBrane;
  const esc = value => String(value ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;', '<':'&lt;', '>':'&gt;', '"':'&quot;', "'":'&#39;'}[c]));
  const title = state => state.title || state.summary || state.label || state.id;
  const stateAttributes = state => `data-head-state="${esc(state.id)}" data-state="${esc(state.id)}"`;
  const colour = head => /^#[0-9a-f]{3}([0-9a-f]{3})?$/i.test(head.colour) ? head.colour : defaultColour;
  const headGlyph = head => glyph(head.glyph, head.glyph === 'circle' ? sharedInk : colour(head), 12);
  const headLabel = head => braneLabel(head.name, head.row);
  function recency(head) {
    const {last, days} = head.trace;
    const label = esc(head.ageLabel);
    const activeDays = days.filter(Boolean).length;
    const explanation = last === null ? 'No dated state in this snapshot'
      : `Last recorded state ${timeFormat.format(last)} UTC. Activity on ${activeDays} of the last 30 days.`;
    const age = last === null ? `<span class="wh-age">${label}</span>`
      : `<time class="wh-age" datetime="${new Date(last).toISOString()}">${label}</time>`;
    const bars = days.map((count, i) => {
      if (!count) return '';
      const x = 4 + i * 2;
      const height = Math.min(11, 5 + Math.log2(count) * 2);
      return `M${x} 12v-${height.toFixed(1)}`;
    }).join('');
    const beforeWindow = !activeDays && last !== null ? '<path class="wh-before" d="m4 8-3 3 3 3"/>' : '';
    return `<span class="wh-recency" title="${esc(explanation)}">
      <svg class="wh-trace" viewBox="0 0 66 16" preserveAspectRatio="none" role="img" aria-label="${esc(explanation)}"><path class="wh-axis" d="M2 12h61m0-10v12"/>
        <path class="wh-bars" d="${bars}"/>${beforeWindow}</svg>${age}</span>`;
  }

  function elapsedLabel(elapsed) {
    if (!Number.isFinite(elapsed)) return 'Undated';
    if (elapsed < 60000) return 'Just now';
    if (elapsed < 3600000) return Math.floor(elapsed / 60000) + 'm ago';
    if (elapsed < DAY) return Math.floor(elapsed / 3600000) + 'h ago';
    return Math.floor(elapsed / DAY) + 'd ago';
  }

  const stamp = (state, compact = false) => state.at && Number.isFinite(Date.parse(state.at))
    ? `<time datetime="${esc(state.at)}" title="${esc(state.at)}">${esc((compact ? rearTimeFormat : timeFormat).format(Date.parse(state.at)))}${compact ? '' : ' UTC'}</time>`
    : '<span>Time unavailable</span>';
  const href = (head, state) => '/surface?' + new URLSearchParams({brane: head.key, at: state.id}).toString() + '#heads';
  const stateLabel = state => [state.open?.label, state.agent, state.episode, state.pr ? '#' + state.pr : ''].filter(Boolean).join(' / ');
  const getPersistence = () => {
    try { return JSON.parse(sessionStorage.getItem(STORAGE_KEY) || '{}'); }
    catch (_) { return {}; }
  };

  function createView(container) {
    const persisted = getPersistence();
    let requested = '';
    try { requested = decodeURIComponent(location.hash.slice(1)); }
    catch (_) { /* A malformed location is equivalent to opening the worktable. */ }
    const positions = new Map(Array.isArray(persisted?.positions) ? persisted.positions : []);
    let selected = typeof persisted?.selected === 'string' ? persisted.selected : '';
    let page = Number.isInteger(persisted?.page) ? Math.max(0, persisted.page) : 0;
    let stale = false, currentModel, heads = [], headByKey = new Map(), entries = new Map();
    let drag, pending, historyKey, historyLimit = HISTORY_PAGE;
    let firstRender = true;
    let width = 280, boardHeight = 0;

    container.classList.add('world-heads');
    container.innerHTML = `
      <header class="wh-bar">
        <div><h1>Branch heads</h1><p>Each strip shows 30 days of recorded activity. Today is at the right edge.</p></div>
        <div class="wh-readout"><p class="wh-counts"></p><p class="wh-status" role="status"></p></div>
      </header>
      <div class="wh-board" role="region" aria-label="Branch head worktable" tabindex="0">
        <div class="wh-plane">
          <svg class="wh-links" aria-hidden="true"></svg>
          <div class="wh-cards"></div>
        </div>
        <p class="wh-empty" hidden></p>
      </div>
      <div class="wh-link-details" hidden></div>
      <footer class="wh-nav" aria-label="Worktable navigation">
        <details class="wh-sessions" hidden><summary></summary><ul></ul></details>
        <div class="wh-tokens" role="group" aria-label="Heads on this page"></div>
        <label class="wh-picker">Find brane <select aria-label="Find brane"></select></label>
        <div class="wh-pages" role="group" aria-label="Worktable pages">
          <button type="button" data-page="-1" aria-label="Previous page">‹</button>
          <span></span>
          <button type="button" data-page="1" aria-label="Next page">›</button>
        </div>
        <button class="wh-reset" type="button">Reset layout</button>
      </footer>
      <dialog class="wh-history" aria-label="Recorded states">
        <header><h2></h2><form method="dialog"><button>Close</button></form></header>
        <div class="wh-history-live"></div><ol></ol><button class="wh-more" type="button">Show earlier states</button>
      </dialog>`;
    const find = q => container.querySelector(q);
    const board = find('.wh-board'), plane = find('.wh-plane'), empty = find('.wh-empty');
    const cards = find('.wh-cards'), svg = find('.wh-links'), linkDetails = find('.wh-link-details');
    const status = find('.wh-status'), counts = find('.wh-counts'), nav = find('.wh-nav'), picker = find('.wh-picker select');
    const tokens = find('.wh-tokens'), pages = find('.wh-pages'), reset = find('.wh-reset');
    const sessions = find('.wh-sessions'), history = find('.wh-history'), list = history.querySelector('ol');
    const more = find('.wh-more'), liveList = find('.wh-history-live');
    let sessionContent = '';

    function save() {
      try { sessionStorage.setItem(STORAGE_KEY, JSON.stringify({selected, page, scrollTop: board.scrollTop, positions: [...positions]})); }
      catch (_) { /* The view still works without browser storage. */ }
    }

    addEventListener('pagehide', save);

    function renderLinks() {
      window.WorldHeadLinks?.render(svg, linkDetails, currentModel, entries, selected);
    }

    function focusHead(key) {
      if (!headByKey.has(key) || selected === key) return;
      const previous = headByKey.get(selected);
      const next = headByKey.get(key);
      selected = key;
      entries.forEach(({el}, id) => el.classList.toggle('is-selected', id === key));
      tokens.querySelectorAll('button').forEach(button => button.setAttribute('aria-pressed', String(button.dataset.locate === key)));
      picker.value = key;
      if ((previous && previous.depth !== 'front') || next.depth !== 'front') layout();
      else renderLinks();
      save();
    }

    function orderedStates(head) {
      return head.states.filter(state => state.id && state.id !== head.latest?.id);
    }

    function card(head) {
      const state = head.latest;
      const earlier = orderedStates(head);
      const past = earlier.slice(0, 3).map((record, i) => `
        <a class="wh-rear" data-depth="${i + 1}" ${stateAttributes(record)} href="${esc(href(head, record))}" title="${esc(title(record))}">
          ${stamp(record, true)} <span>${esc(title(record))}</span>
        </a>`).reverse().join('');
      const activityStrip = recency(head);
      const text = state ? title(state) : 'No recorded state';
      const stateContent = state?.id
        ? `<a class="wh-open" ${stateAttributes(state)} href="${esc(href(head, state))}" title="${esc(text)}">${esc(text)}</a>`
        : `<strong class="wh-open">${esc(text)}</strong>`;
      const info = state ? stateLabel(state) : '';
      const [session] = head.live;
      const sessionText = stale ? 'Seen live' : head.live.length > 1 ? head.live.length + ' live' : 'Live';
      const sessionTitle = stale ? 'Last seen running' : 'Running session';
      const sessionControl = !session ? '' : head.live.length === 1
        ? `<a class="wh-live${stale ? ' was-live' : ''}" ${stateAttributes(session)} href="${esc(href(head, session))}" title="${esc(sessionTitle + ' / ' + session.episode)}">${sessionText}</a>`
        : `<button class="wh-live${stale ? ' was-live' : ''}" type="button" data-history>${sessionText}</button>`;
      return `
        <button class="wh-restore" type="button" aria-label="Bring ${esc(headLabel(head))} forward. ${esc(head.ageLabel)}" title="${esc(headLabel(head))}">
          <span class="wh-glyph" aria-hidden="true">${headGlyph(head)}</span>
          <span class="wh-head-name">${esc(headLabel(head))}</span>
          ${activityStrip}
        </button>
        ${past ? `<div class="wh-depth" data-count="${Math.min(3, earlier.length)}" aria-hidden="true"></div>` : ''}
        ${past}
        <div class="wh-face">
          <div class="wh-handle">
            <button class="wh-grip" type="button" aria-label="Move ${esc(headLabel(head))}" title="Drag or use the arrow keys to move this head">
              <span aria-hidden="true" class="wh-glyph">${headGlyph(head)}</span>
              <span class="wh-head-name">${esc(headLabel(head))}</span>
              <i aria-hidden="true"></i>
            </button>
            ${sessionControl}
          </div>
          <div class="wh-content">
            ${head.repo ? `<span class="wh-repo" title="${esc(head.repo)}">${esc(head.repo)}</span>` : ''}
            ${stateContent}
            ${info && info !== text ? `<span class="wh-record" title="${esc(info)}">${esc(info)}</span>` : ''}
          </div>
          <div class="wh-meta">
            ${activityStrip}
            ${earlier.length > 3 ? `<button class="wh-all" type="button" data-history>All ${earlier.length + 1} states</button>` : ''}
            <button class="wh-details" type="button" data-history aria-label="Details for ${esc(headLabel(head))}">Details</button>
          </div>
        </div>`;
    }

    function layout() {
      if (!board.clientWidth) return;
      const overview = board.clientWidth <= 540;
      const gutter = overview ? 10 : 30;
      const gap = overview ? 12 : 34;
      const available = Math.max(120, board.clientWidth - 2 * gutter);
      const cell = overview ? 127 : 214;
      const columns = Math.min(4, Math.max(1, Math.floor((available + gap) / (cell + gap))));
      width = Math.min(320, Math.floor((available - gap * (columns - 1)) / columns));

      const active = document.activeElement;
      let replacementFocus;
      const pageEntries = [...entries.values()];
      const fallback = pageEntries.find(entry => entry.head.depth === 'front') ? null : pageEntries[0];
      const groups = {front: [], near: [], far: []};
      for (const entry of pageEntries) {
        const {el, head} = entry;
        const home = entry === fallback ? 'front' : head.depth;
        const depth = head.key === selected || positions.has(head.key) ? 'front' : home;
        if (el.dataset.zone !== depth && el.contains(active)) {
          if (depth === 'far') replacementFocus = el.querySelector('.wh-restore');
          else if (active.classList.contains('wh-restore')) replacementFocus = el.querySelector('.wh-grip');
        }
        el.dataset.zone = depth;
        groups[home].push(entry);
      }
      let nextTop = 12;
      const slots = [];
      Object.entries(groups).forEach(([group, members]) => {
        const front = group === 'front', far = group === 'far';
        if (!members.length) return;

        const count = far ? Math.min(6, Math.max(columns, Math.floor((available + 8) / 166))) : columns;
        const cardGap = far ? 8 : gap;
        const cardWidth = front ? width : Math.min(far ? 206 : 270, Math.floor((available - cardGap * (count - 1)) / count), far ? width : Math.floor(width * .9));
        const left = Math.max(gutter, (board.clientWidth - count * cardWidth - (count - 1) * cardGap) / 2);
        const rowHeight = front ? 239 : far ? 60 : 192;
        const firstFace = nextTop + (front ? 90 : far ? 8 : 73);
        members.forEach((entry, index) => {
          const column = index % count;
          entry.el.style.setProperty('--wh-width', (entry.el.dataset.zone === 'front' ? width : cardWidth) + 'px');
          slots.push({
            entry,
            x: left + column * (cardWidth + cardGap),
            y: firstFace + Math.floor(index / count) * rowHeight + (front ? (column % 2) * 14 : 0),
          });
        });
        nextTop = firstFace + (Math.ceil(members.length / count) - 1) * rowHeight + (front ? 154 : far ? 60 : 133);
      });
      boardHeight = Math.max(board.clientHeight, nextTop, ...slots.map(({entry, y}) =>
        Math.max(y, positions.get(entry.head.key)?.[1] || 0) + entry.el.offsetHeight + 12));
      plane.style.height = boardHeight + 'px';
      slots.forEach(({entry: {el, head}, x, y}, index) => {
        if (cards.children[index] !== el) cards.insertBefore(el, cards.children[index] || null);
        const saved = positions.get(head.key);
        if (saved) applyPosition(el, 8 + saved[0] * dragWidth(el), saved[1]);
        else applyPosition(el, x, y);
      });
      if (replacementFocus) replacementFocus.focus({preventScroll:true});
      else if (cards.contains(active) && document.activeElement !== active) active.focus({preventScroll:true});
      renderLinks();
    }

    function dragWidth(el) { return Math.max(1, board.clientWidth - el.offsetWidth - 16); }

    function applyPosition(el, x, y) {
      const position = [
        Math.round(Math.max(8, Math.min(board.clientWidth - el.offsetWidth - 8, x))),
        Math.round(Math.max(el.dataset.zone === 'far' ? 8 : 88, Math.min(boardHeight - el.offsetHeight - 10, y)))
      ];
      el.style.left = position[0] + 'px';
      el.style.top = position[1] + 'px';
      return position;
    }

    function move(el, x, y) {
      const [left, top] = applyPosition(el, x, y);
      positions.set(el.dataset.head, [(left - 8) / dragWidth(el), top]);
      renderLinks();
    }

    function drawPage() {
      const shown = heads.slice(page * PAGE_SIZE, (page + 1) * PAGE_SIZE);
      const previous = entries;
      const active = document.activeElement;
      entries = new Map();
      shown.forEach(head => {
        const previousEntry = previous.get(head.key);
        const el = previousEntry?.el || document.createElement('article');
        const markup = card(head);
        const focusedClass = el.contains(active) && active.className;
        el.className = 'wh-card' + (head.key === selected ? ' is-selected' : '');
        el.dataset.head = head.key;
        el.dataset.recency = head.recency;
        el.style.setProperty('--wh-base', colour(head));
        if (previousEntry?.markup !== markup) {
          el.innerHTML = markup;
          if (focusedClass) {
            const replacement = [...el.querySelectorAll('a,button')].find(child =>
              child.className === focusedClass && child.getAttribute('href') === active.getAttribute('href'));
            (replacement || el.querySelector('.wh-grip'))?.focus({preventScroll:true});
          }
        }
        if (el.parentNode !== cards) cards.append(el);
        entries.set(head.key, {el, markup, head});
      });
      previous.forEach(({el}, key) => { if (!entries.has(key)) el.remove(); });
      const tokenFocus = tokens.contains(active) ? active.dataset.locate : '';
      const tokenMarkup = shown.slice(0, 8).map(head => `
        <button type="button" data-locate="${esc(head.key)}" title="${esc(headLabel(head))}"
          aria-label="Locate ${esc(headLabel(head))}" style="--wh-token:${colour(head)}">
          ${headGlyph(head)}
        </button>`).join('');
      if (tokens.dataset.markup !== tokenMarkup) {
        tokens.innerHTML = tokenMarkup;
        tokens.dataset.markup = tokenMarkup;
        [...tokens.children].find(button => button.dataset.locate === tokenFocus)?.focus({preventScroll:true});
      }
      tokens.querySelectorAll('button').forEach(button => button.setAttribute('aria-pressed', String(button.dataset.locate === selected)));
      const count = Math.ceil(heads.length / PAGE_SIZE);
      pages.hidden = count <= 1;
      pages.querySelector('span').textContent = `${page + 1} / ${count}`;
      pages.firstElementChild.disabled = page === 0;
      pages.lastElementChild.disabled = page + 1 >= count;
      layout();
    }

    function locate(key, focus) {
      const index = heads.findIndex(head => head.key === key);
      if (index < 0) return;
      const targetPage = Math.floor(index / PAGE_SIZE);
      focusHead(key);
      if (targetPage !== page) {
        page = targetPage;
        drawPage();
      }
      const el = entries.get(key).el;
      el.scrollIntoView({block:'center', inline:'nearest'});
      if (focus) el.querySelector('.wh-grip').focus({preventScroll:true});
      save();
    }

    function renderHistory(resetScroll = false) {
      const head = headByKey.get(historyKey);
      if (!head?.latest?.id) { if (history.open) history.close(); return; }
      const all = [head.latest, ...orderedStates(head)];
      history.querySelector('h2').textContent = headLabel(head);
      const liveMarkup = head.live.length ? `
        <h3>${stale ? 'Last seen running' : 'Running sessions'}</h3>
        <ul>${head.live.map(session => `
          <li><a ${stateAttributes(session)} href="${esc(href(head, session))}">${esc(session.episode)}</a>
          <span>${esc(session.label || '')}</span></li>`).join('')}</ul>` : '';
      if (liveList.dataset.markup !== liveMarkup) {
        const link = liveList.contains(document.activeElement) && document.activeElement.getAttribute('href');
        liveList.innerHTML = liveMarkup;
        liveList.dataset.markup = liveMarkup;
        if (link) [...liveList.querySelectorAll('a')].find(el => el.getAttribute('href') === link)?.focus({preventScroll:true});
      }
      const markup = all.slice(0, historyLimit).map((state, index) => `
        <li><div><a ${stateAttributes(state)} href="${esc(href(head, state))}">${esc(title(state))}</a>
          ${index === 0 && state.summary && state.summary !== title(state) ? `<p>${esc(state.summary)}</p>` : ''}
          ${index === 0 && stateLabel(state) ? `<span>${esc(stateLabel(state))}</span>` : ''}
        </div> ${stamp(state)}</li>`).join('');
      if (list.dataset.markup !== markup) {
        const active = list.contains(document.activeElement) && document.activeElement.getAttribute('href');
        list.innerHTML = markup;
        list.dataset.markup = markup;
        if (active) [...list.querySelectorAll('a')].find(link => link.getAttribute('href') === active)?.focus({preventScroll:true});
      }
      more.hidden = all.length <= historyLimit;
      if (resetScroll) history.scrollTop = 0;
    }

    container.addEventListener('click', event => {
      const button = event.target.closest('button');
      const el = event.target.closest('.wh-card');
      if (button?.classList.contains('wh-restore') && el) {
        locate(el.dataset.head, true);
        return;
      }
      if (el) focusHead(el.dataset.head);
      if (!button) return;
      if (button.dataset.locate) locate(button.dataset.locate, true);
      if (button.dataset.page) {
        page += Number(button.dataset.page);
        drawPage();
        board.scrollTop = 0;
      }
      if (button === reset) {
        positions.clear();
        board.scrollTop = 0;
        layout();
        save();
      }
      if (button.dataset.history != null && el) {
        historyKey = el.dataset.head;
        historyLimit = HISTORY_PAGE;
        renderHistory(true);
        history.showModal();
      }
      if (button === more) { historyLimit += HISTORY_PAGE; renderHistory(); }
    });
    picker.addEventListener('change', () => locate(picker.value, true));

    board.addEventListener('keydown', event => {
      if (!event.target.closest('.wh-grip')) return;
      const delta = {ArrowLeft:[-12, 0], ArrowRight:[12, 0], ArrowUp:[0, -12], ArrowDown:[0, 12]}[event.key];
      if (!delta) return;
      event.preventDefault();
      event.stopPropagation();
      const el = event.target.closest('.wh-card');
      focusHead(el.dataset.head);
      move(el, el.offsetLeft + delta[0], el.offsetTop + delta[1]);
      save();
    });
    board.addEventListener('pointerdown', event => {
      const grip = event.target.closest('.wh-grip');
      if (!grip || drag || event.button !== 0) return;
      const el = grip.closest('.wh-card');
      focusHead(el.dataset.head);
      drag = {el, pointer: event.pointerId, x: event.clientX, y: event.clientY,
        left: el.offsetLeft, top: el.offsetTop, scroll: board.scrollTop};
      grip.setPointerCapture(event.pointerId);
    });
    board.addEventListener('pointermove', event => {
      if (drag?.pointer !== event.pointerId) return;
      const dx = event.clientX - drag.x, dy = event.clientY - drag.y;
      if (!drag.moved && Math.hypot(dx, dy) < 4) return;
      drag.moved = true;
      drag.el.classList.add('is-dragging');
      move(drag.el, drag.left + dx, drag.top + dy + board.scrollTop - drag.scroll);
    });
    function finishDrag(event) {
      if (drag?.pointer !== event.pointerId) return;
      drag.el.classList.remove('is-dragging');
      drag = null;
      save();
      if (pending) { const [model, options] = pending; pending = null; update(model, options); }
    }
    board.addEventListener('pointerup', finishDrag);
    board.addEventListener('pointercancel', finishDrag);
    board.addEventListener('lostpointercapture', finishDrag);
    new ResizeObserver(layout).observe(board);

    function update(model, options) {
      if (drag) { pending = [model, options]; return; }
      const projection = model ? window.WorldHeadsData.project(model) : {heads: [], unlinked: [], status: ''};
      stale = Boolean(options.stale);
      currentModel = model;
      const modelTime = Date.parse(model?.now);
      const now = Number.isFinite(modelTime) ? modelTime : Date.now();
      const selectedWasVisible = entries.has(selected);
      heads = projection.heads.map(head => {
        const trace = WorldHeadActivity.trace(head.states, now);
        const activity = trace.last ?? -Infinity;
        const elapsed = now - activity;
        const live = !stale && head.live.length > 0;
        const depth = live || elapsed < 2 * DAY ? 'front' : elapsed < 7 * DAY ? 'near' : 'far';
        const recency = trace.last === null ? 'unknown' : elapsed < 2 * DAY ? 'recent' : elapsed < 7 * DAY ? 'week' : 'older';
        return {...head, trace, activity, liveNow: live, depth, recency, ageLabel: elapsedLabel(elapsed)};
      }).sort((a, b) =>
        Number(b.liveNow) - Number(a.liveNow) || b.activity - a.activity || (b.latest?.rank ?? 0) - (a.latest?.rank ?? 0) || a.name.localeCompare(b.name));
      headByKey = new Map(heads.map(head => [head.key, head]));
      const restore = firstRender && heads.length;
      if (restore) {
        firstRender = false;
        if (headByKey.has(requested)) selected = requested;
      }
      if ((selectedWasVisible || (restore && selected === requested)) && headByKey.has(selected)) {
        page = Math.floor(heads.findIndex(head => head.key === selected) / PAGE_SIZE);
      }
      page = Math.min(page, Math.max(0, Math.ceil(heads.length / PAGE_SIZE) - 1));
      const notice = stale && model ? 'Live updates disconnected. Showing the last recorded branch heads.' : '';
      const message = [notice, projection.status].filter(Boolean).join(' ');
      if (status.textContent !== message) status.textContent = message;
      const edge = model?.nowEdge;
      const liveCount = heads.reduce((total, head) => total + head.live.length, projection.unlinked.length);
      const liveLabel = liveCount === 1 ? 'live session' : 'live sessions';
      const availability = stale || !edge ? 'unavailable'
        : model.schema === 'world/surface-model/v2' ? edge.open?.status || 'unavailable'
        : edge.sessionsUnavailable ? 'unavailable' : 'complete';
      const liveReadout = availability === 'complete' ? `${liveCount} ${liveLabel}`
        : liveCount > 0 ? `${liveCount} ${liveLabel} ${stale ? 'last seen' : 'recorded'}`
        : availability === 'unavailable' ? 'Live sessions unavailable' : 'Live sessions incomplete';
      counts.textContent = model ? `${heads.length} ${heads.length === 1 ? 'head' : 'heads'} / ${liveReadout}` : '';
      empty.hidden = heads.length > 0;
      empty.textContent = model?.reason || (model ? 'No recorded branch heads are available.' : stale
        ? [options.reason || 'World’s branch history is unavailable.', 'This view will retry automatically.'].join(' ')
        : 'Waiting for World’s branch history.');
      const optionMarkup = '<option value="">Jump to…</option>' + heads.map(head =>
        `<option value="${esc(head.key)}">${esc(headLabel(head))}${head.repo ? ' / ' + esc(head.repo) : ''}</option>`).join('');
      if (picker.dataset.markup !== optionMarkup) {
        picker.innerHTML = optionMarkup;
        picker.dataset.markup = optionMarkup;
      }
      picker.value = headByKey.has(selected) ? selected : '';
      picker.disabled = !heads.length;
      reset.disabled = !heads.length;
      nav.classList.toggle('wh-no-cards', !heads.length);
      sessions.hidden = !projection.unlinked.length;
      sessions.querySelector('summary').textContent = `${projection.unlinked.length} unlinked ${stale ? 'last seen' : 'live'} ${projection.unlinked.length === 1 ? 'session' : 'sessions'}`;
      const sessionMarkup = projection.unlinked.map(({episode, reason}) =>
        `<li><strong>${esc(episode)}</strong><span>${esc(reason)}</span></li>`).join('');
      if (sessionMarkup !== sessionContent) {
        sessions.querySelector('ul').innerHTML = sessionMarkup;
        sessionContent = sessionMarkup;
      }
      drawPage();
      if (restore) {
        if (requested && requested !== persisted?.selected && entries.has(requested)) locate(requested, false);
        else if (Number.isFinite(persisted?.scrollTop)) board.scrollTop = persisted.scrollTop;
      }
      if (history.open) renderHistory();
    }
    return update;
  }

  window.WorldHeads = {
    render(container, model, options = {}) {
      if (!viewByContainer.has(container)) viewByContainer.set(container, createView(container));
      viewByContainer.get(container)(model, options);
    }
  };
})();
