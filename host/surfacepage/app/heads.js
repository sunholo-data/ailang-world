// End-on view of the recorded branch lanes. The World node view owns state inspection.
(() => {
  'use strict';
  const PAGE_SIZE = 16;
  const HISTORY_PAGE = 40;
  const STORAGE_KEY = 'world-heads-worktable';
  const QUIET_AGE = 7 * 24 * 60 * 60 * 1000;
  const defaultColour = '#657b91';
  const viewByContainer = new WeakMap();
  const timeFormat = new Intl.DateTimeFormat(undefined, {
    year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit', timeZone: 'UTC'
  });
  const rearTimeFormat = new Intl.DateTimeFormat(undefined, {
    month: 'short', day: 'numeric', timeZone: 'UTC'
  });
  const glyphs = { circle: '○', square: '▪', triangle: '▲', diamond: '◆', hexagon: '⬢', pentagon: '⬟', star: '★' };
  const esc = value => String(value ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;', '<':'&lt;', '>':'&gt;', '"':'&quot;', "'":'&#39;'}[c]));
  const title = state => state.title || state.summary || state.label || state.id;
  const stateTime = state => Date.parse(state?.at) || 0;
  const colour = head => /^#[0-9a-f]{3}([0-9a-f]{3})?$/i.test(head.colour) ? head.colour : defaultColour;
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
    let stale = false, quietBefore = 0, heads = [], headByKey = new Map(), entries = new Map();
    let drag, pending, historyKey, historyLimit = HISTORY_PAGE;
    let firstRender = true;
    let width = 280, boardHeight = 0;

    container.classList.add('world-heads');
    container.innerHTML = `
      <header class="wh-bar">
        <div><h1>Branch heads</h1><p>Drag a head to move it. Select any pane to open its recorded state.</p></div>
        <div class="wh-readout"><p class="wh-counts"></p><p class="wh-status" role="status"></p></div>
      </header>
      <div class="wh-board" role="region" aria-label="Branch head worktable" tabindex="0">
        <div class="wh-plane"></div>
        <p class="wh-empty" hidden></p>
      </div>
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

    function focusHead(key) {
      if (!headByKey.has(key)) return;
      selected = key;
      entries.forEach(({el}, id) => el.classList.toggle('is-selected', id === key));
      tokens.querySelectorAll('button').forEach(button => button.setAttribute('aria-pressed', String(button.dataset.locate === key)));
      picker.value = key;
      save();
    }

    function orderedStates(head) {
      return head.states.filter(state => state.id && state.id !== head.latest?.id);
    }

    function card(head) {
      const state = head.latest;
      const earlier = orderedStates(head);
      const past = earlier.slice(0, 3).map((record, i) => `
        <a class="wh-rear" data-depth="${i + 1}" href="${esc(href(head, record))}" title="${esc(title(record))}">
          ${stamp(record, true)} <span>${esc(title(record))}</span>
        </a>`).reverse().join('');
      const text = state ? title(state) : 'No recorded state';
      const stateContent = state?.id
        ? `<a class="wh-open" href="${esc(href(head, state))}" title="${esc(text)}">${esc(text)}</a>`
        : `<strong class="wh-open">${esc(text)}</strong>`;
      const info = state ? stateLabel(state) : '';
      const [session] = head.live;
      const sessionText = stale ? 'Seen live' : head.live.length > 1 ? head.live.length + ' live' : 'Live';
      const sessionTitle = stale ? 'Last seen running' : 'Running session';
      const sessionControl = !session ? '' : head.live.length === 1
        ? `<a class="wh-live${stale ? ' was-live' : ''}" href="${esc(href(head, session))}" title="${esc(sessionTitle + ' / ' + session.episode)}">${sessionText}</a>`
        : `<button class="wh-live${stale ? ' was-live' : ''}" type="button" data-history>${sessionText}</button>`;
      return `
        ${past ? `<div class="wh-depth" data-count="${Math.min(3, earlier.length)}" aria-hidden="true"></div>` : ''}
        ${past}
        <div class="wh-face">
          <div class="wh-handle">
            <button class="wh-grip" type="button" aria-label="Move ${esc(head.name)}" title="Drag or use the arrow keys to move this head">
              <span aria-hidden="true" class="wh-glyph">${glyphs[head.glyph] || '•'}</span>
              <span class="wh-head-name">${esc(head.name)}</span>
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
            ${state ? `<span class="wh-stamp">${stamp(state)}</span>` : ''}
            ${earlier.length > 3 ? `<button class="wh-all" type="button" data-history>All ${earlier.length + 1} states</button>` : ''}
            <button class="wh-details" type="button" data-history aria-label="Details for ${esc(head.name)}">Details</button>
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
      const usedWidth = columns * width + (columns - 1) * gap;
      const left = Math.max(gutter, (board.clientWidth - usedWidth) / 2);
      boardHeight = Math.max(board.clientHeight, Math.ceil(entries.size / columns) * 236 + 18);
      plane.style.setProperty('--wh-width', width + 'px');
      plane.style.height = boardHeight + 'px';
      let index = 0;
      entries.forEach(({el}, key) => {
        const column = index % columns;
        const startX = left + column * (width + gap);
        const startY = 104 + Math.floor(index / columns) * 236 + (column % 2) * 18;
        const saved = positions.get(key);
        if (saved) applyPosition(el, 8 + saved[0] * dragWidth(), saved[1]);
        else applyPosition(el, startX, startY);
        index++;
      });
    }

    function dragWidth() { return Math.max(1, board.clientWidth - width - 16); }

    function applyPosition(el, x, y) {
      const position = [
        Math.round(Math.max(8, Math.min(board.clientWidth - width - 8, x))),
        Math.round(Math.max(88, Math.min(boardHeight - el.offsetHeight - 10, y)))
      ];
      el.style.left = position[0] + 'px';
      el.style.top = position[1] + 'px';
      return position;
    }

    function move(el, x, y) {
      const [left, top] = applyPosition(el, x, y);
      positions.set(el.dataset.head, [(left - 8) / dragWidth(), top]);
    }

    function drawPage() {
      const shown = heads.slice(page * PAGE_SIZE, (page + 1) * PAGE_SIZE);
      const previous = entries;
      const active = document.activeElement;
      entries = new Map();
      shown.forEach((head, index) => {
        const previousEntry = previous.get(head.key);
        const el = previousEntry?.el || document.createElement('article');
        const markup = card(head);
        const focusedClass = el.contains(active) && active.className;
        el.className = 'wh-card' + (head.key === selected ? ' is-selected' : '') + (head.live.length || stateTime(head.latest) >= quietBefore ? '' : ' is-quiet');
        el.dataset.head = head.key;
        el.style.setProperty('--wh-base', colour(head));
        if (previousEntry?.markup !== markup) {
          el.innerHTML = markup;
          if (focusedClass) {
            const replacement = [...el.querySelectorAll('a,button')].find(child =>
              child.className === focusedClass && child.getAttribute('href') === active.getAttribute('href'));
            (replacement || el.querySelector('.wh-grip'))?.focus({preventScroll:true});
          }
        }
        if (plane.children[index] !== el) plane.insertBefore(el, plane.children[index] || null);
        entries.set(head.key, {el, markup});
      });
      previous.forEach(({el}, key) => { if (!entries.has(key)) el.remove(); });
      const tokenFocus = tokens.contains(active) ? active.dataset.locate : '';
      const tokenMarkup = shown.slice(0, 8).map(head => `
        <button type="button" data-locate="${esc(head.key)}" title="${esc(head.name)}"
          aria-label="Locate ${esc(head.name)}" style="--wh-token:${colour(head)}">
          ${glyphs[head.glyph] || '•'}
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
    }

    function renderHistory(resetScroll = false) {
      const head = headByKey.get(historyKey);
      if (!head?.latest?.id) { if (history.open) history.close(); return; }
      const all = [head.latest, ...orderedStates(head)];
      history.querySelector('h2').textContent = head.name;
      const liveMarkup = head.live.length ? `
        <h3>${stale ? 'Last seen running' : 'Running sessions'}</h3>
        <ul>${head.live.map(session => `
          <li><a href="${esc(href(head, session))}">${esc(session.episode)}</a>
          <span>${esc(session.label || '')}</span></li>`).join('')}</ul>` : '';
      if (liveList.dataset.markup !== liveMarkup) {
        const link = liveList.contains(document.activeElement) && document.activeElement.getAttribute('href');
        liveList.innerHTML = liveMarkup;
        liveList.dataset.markup = liveMarkup;
        if (link) [...liveList.querySelectorAll('a')].find(el => el.getAttribute('href') === link)?.focus({preventScroll:true});
      }
      const markup = all.slice(0, historyLimit).map((state, index) => `
        <li><div><a href="${esc(href(head, state))}">${esc(title(state))}</a>
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
      quietBefore = (Date.parse(model?.now) || Date.now()) - QUIET_AGE;
      const selectedWasVisible = entries.has(selected);
      heads = [...projection.heads].sort((a, b) =>
        Number(b.live.length > 0) - Number(a.live.length > 0) || stateTime(b.latest) - stateTime(a.latest) || (b.latest?.rank ?? 0) - (a.latest?.rank ?? 0) || a.name.localeCompare(b.name));
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
        `<option value="${esc(head.key)}">${esc(head.name)}${head.repo ? ' / ' + esc(head.repo) : ''}</option>`).join('');
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
