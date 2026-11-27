(function (root) {
  'use strict';

  const schema = 'world/surface-trajectory/v1';
  const maxRawBytes = 8 * 1024 * 1024;
  const pageChars = 12000;
  const labels = {
    'session-started': 'Session started', 'turn-started': 'Turn started', 'turn-completed': 'Turn completed',
    'turn-failed': 'Turn failed', 'item-completed': 'Item completed', 'assistant-event': 'Assistant event',
    'run-result': 'Run result', 'tool-called': 'Tool called', 'tool-result': 'Tool result'
  };

  function trajectoryPath(brane, at) {
    return '/surface/trajectory.json?brane=' + encodeURIComponent(brane) + '&at=' + encodeURIComponent(at);
  }

  function rawReceipt(raw, episode) {
    return Boolean(raw && /^[a-z0-9][a-z0-9-]{0,63}$/.test(episode) &&
      raw.filename === episode + '.transcript.jsonl' && /^[a-f0-9]{64}$/.test(raw.sha256));
  }

  async function verifyTranscript(file, receipt, cryptoAPI) {
    if (file.name !== receipt.filename) throw new Error('Choose ' + receipt.filename + '.');
    if (file.size > maxRawBytes) throw new Error('This transcript exceeds the 8 MiB viewer limit.');
    if (!cryptoAPI?.subtle) throw new Error('Open the preview over HTTPS or on localhost to verify a local transcript.');
    const bytes = await file.arrayBuffer();
    if (bytes.byteLength !== file.size || bytes.byteLength > maxRawBytes) {
      throw new Error('The file size changed while reading it.');
    }
    const digest = await cryptoAPI.subtle.digest('SHA-256', bytes);
    const hex = Array.from(new Uint8Array(digest), byte => byte.toString(16).padStart(2, '0')).join('');
    if (hex !== receipt.sha256) throw new Error('The transcript does not match the recorded SHA-256.');
    try {
      return new TextDecoder('utf-8', { fatal: true, ignoreBOM: true }).decode(bytes);
    } catch (_) {
      throw new Error('The verified transcript is not UTF-8 text.');
    }
  }

  // Prefer full lines. A very long JSONL record can span pages.
  function transcriptPages(text) {
    const pages = [];
    for (let start = 0; start < text.length;) {
      let end = Math.min(start + pageChars, text.length);
      if (end < text.length) {
        const line = text.lastIndexOf('\n', end - 1);
        if (line >= start + pageChars / 2) end = line + 1;
        const next = text.charCodeAt(end);
        if (next >= 0xDC00 && next <= 0xDFFF) end--;
      }
      pages.push(text.slice(start, end));
      start = end;
    }
    return pages.length ? pages : [''];
  }

  function append(parent, tag, text, className) {
    const el = parent.ownerDocument.createElement(tag);
    if (className) el.className = className;
    if (text !== undefined) el.textContent = text;
    parent.append(el);
    return el;
  }

  function baseTarget(target) {
    return typeof target === 'string' ? target.replace(/:\d+-\d+$/, '') : '';
  }

  function shortTime(iso) {
    if (typeof iso !== 'string' || !iso) return '';
    const match = /T(\d{2}:\d{2}:\d{2})/.exec(iso);
    return match ? match[1] + ' UTC' : iso;
  }

  // Collapse consecutive calls to the same tool, surface, status, and base file path
  // so ranged file reads or repeated searches scan as one high-signal row.
  function collapseCalls(calls) {
    if (!Array.isArray(calls)) return [];
    const groups = [];
    for (const row of calls) {
      const base = baseTarget(row.target);
      const status = row.status && row.status !== 'ok' ? row.status : '';
      const prev = groups.at(-1);
      if (prev && prev.tool === row.tool && prev.surface === row.surface && prev.status === status && prev.base === base) {
        prev.contiguous = prev.contiguous !== false && row.entry === prev.lastEntry + 1;
        prev.count++;
        prev.lastEntry = row.entry;
        prev.at = row.at || prev.at;
        if (row.target !== prev.target) prev.target = base;
      } else {
        groups.push({
          tool: row.tool, surface: row.surface, status, base,
          target: row.target || '', count: 1, contiguous: true,
          firstEntry: row.entry, lastEntry: row.entry, at: row.at || ''
        });
      }
    }
    return groups;
  }

  function isHighSignalCall(c) {
    return Boolean(c && ((c.status && c.status !== 'ok') || /write|edit|delete|exec|run|cli/i.test(c.tool || '')));
  }

  function looksLikeFilePath(s) {
    return Boolean(s && !s.startsWith('-') && !/\s/.test(s) && (s.includes('/') || /\.[a-z0-9]{1,8}$/i.test(s)));
  }

  function renderCallRow(list, item) {
    const li = append(list, 'li', undefined, 'trajectory-row');
    append(li, 'span', item.tool, 'trajectory-tool');
    if (item.target) {
      li.append(' ');
      const cleanPath = baseTarget(item.target);
      if (looksLikeFilePath(cleanPath)) {
        const btn = append(li, 'button', item.target, 'trajectory-target trajectory-file-btn');
        btn.type = 'button';
        btn.setAttribute('data-file-path', cleanPath);
        btn.setAttribute('data-file-source', /write|edit|delete/i.test(item.tool) ? 'episode-write' : 'episode-read');
      } else {
        append(li, 'code', item.target, 'trajectory-target');
      }
    }
    if (item.count > 1) {
      li.append(' ');
      append(li, 'span', '×' + item.count, 'trajectory-burst');
    }
    if (item.status) {
      li.append(' ');
      append(li, 'span', item.status, 'trajectory-status');
    }
    const entryLabel = item.count > 1 && item.contiguous !== false ? 'entries ' + item.firstEntry + '–' + item.lastEntry : 'entry ' + item.lastEntry;
    li.append(' ');
    const metaWrap = append(li, 'span', undefined, 'trajectory-meta');
    if (item.lastEntry != null) {
      const whyBtn = append(metaWrap, 'button', entryLabel, 'trajectory-why-btn');
      whyBtn.type = 'button';
      whyBtn.setAttribute('data-verb', 'why ' + item.lastEntry);
      whyBtn.setAttribute('title', 'Explain why World allowed ' + entryLabel);
    }
    const tm = shortTime(item.at);
    if (tm) append(metaWrap, 'span', ' · ' + tm);
  }

  function renderWorldCalls(body, state, source) {
    const section = append(body, 'section', undefined, 'trajectory-source');
    const calls = Array.isArray(source?.calls) ? source.calls : [];
    const total = Number.isSafeInteger(source?.total) && source.total > 0 ? source.total : calls.length;
    const heading = total ? 'World tool calls · ' + total : 'World tool calls';
    section.setAttribute('aria-label', 'World tool calls');
    const head = append(section, 'div', undefined, 'trajectory-source-head');
    append(head, 'h3', heading);
    if (!source?.available) {
      append(section, 'p', 'World tool calls unavailable.', 'muted');
      return;
    }
    if (!calls.length) {
      append(section, 'p', 'No recorded World tool calls.', 'muted');
      return;
    }
    const sigCalls = calls.filter(isHighSignalCall);
    const hasFilter = calls.length > 5 && sigCalls.length > 0 && sigCalls.length < calls.length;
    const onlySig = hasFilter && state.showAllCalls !== true;
    const active = onlySig ? sigCalls : calls;
    if (calls.length > 1) {
      const controls = append(head, 'div', undefined, 'trajectory-controls');
      if (hasFilter) {
        const sigBtn = append(controls, 'button', onlySig ? 'Changes & exec · ' + sigCalls.length + ' ✓' : 'All calls · ' + calls.length);
        sigBtn.type = 'button';
        sigBtn.setAttribute('aria-pressed', String(onlySig));
        sigBtn.addEventListener('click', () => { state.showAllCalls = onlySig; render(state); });
      }
      const newestBtn = append(controls, 'button', state.oldestFirst ? 'Newest first' : 'Newest first ✓');
      newestBtn.type = 'button';
      newestBtn.setAttribute('aria-pressed', String(!state.oldestFirst));
      newestBtn.addEventListener('click', () => { state.oldestFirst = !state.oldestFirst; render(state); });
      const groupBtn = append(controls, 'button', state.ungrouped ? 'Group bursts' : 'Grouped ✓');
      groupBtn.type = 'button';
      groupBtn.setAttribute('aria-pressed', String(!state.ungrouped));
      groupBtn.addEventListener('click', () => { state.ungrouped = !state.ungrouped; render(state); });
    }
    const rows = state.ungrouped
      ? active.map(c => ({ tool: c.tool, surface: c.surface, status: c.status && c.status !== 'ok' ? c.status : '', target: c.target || '', count: 1, firstEntry: c.entry, lastEntry: c.entry, at: c.at || '' }))
      : collapseCalls(active);
    const ordered = state.oldestFirst ? rows : rows.slice().reverse();
    const list = append(section, 'ol', undefined, 'trajectory-list');
    for (const item of ordered) renderCallRow(list, item);
    if (source.truncated) {
      const note = source.total > calls.length
        ? 'Showing the latest ' + calls.length + ' of ' + source.total + ' recorded entries.'
        : 'Only the first recorded entries are shown.';
      append(section, 'p', note, 'muted');
    }
  }

  function group(body, heading, source, items, line, empty) {
    const section = append(body, 'section', undefined, 'trajectory-source');
    section.setAttribute('aria-label', heading);
    const head = append(section, 'div', undefined, 'trajectory-source-head');
    append(head, 'h3', items?.length ? heading + ' · ' + items.length : heading);
    if (!source?.available) {
      append(section, 'p', heading + ' unavailable.', 'muted');
      return;
    }
    if (!Array.isArray(items) || !items.length) {
      append(section, 'p', empty, 'muted');
    } else {
      const list = append(section, 'ol', undefined, 'trajectory-list');
      for (const row of items) {
        const li = append(list, 'li', line(row), 'trajectory-row');
        li.setAttribute('data-harness-row', '1');
      }
    }
    if (source.truncated) append(section, 'p', 'Only the first recorded entries are shown.', 'muted');
  }

  let current = null;

  function clear() {
    current?.request?.abort();
    current = null;
  }

  function rawViewer(body, state, receipt) {
    const raw = append(body, 'details', undefined, 'trajectory-raw');
    raw.open = state.rawOpen;
    append(raw, 'summary', 'View raw transcript locally');
    const info = append(raw, 'p');
    info.append('Default rig location: ');
    append(info, 'code', '~/.ailang/world/sessions/' + receipt.filename);
    append(raw, 'p', 'Choose a local copy. The browser checks its SHA-256 and displays it on this device. The file is never uploaded.', 'muted');
    const label = append(raw, 'label', 'Choose transcript ');
    const input = append(label, 'input');
    input.type = 'file';
    input.accept = '.jsonl,application/x-ndjson';
    const status = append(raw, 'p', state.rawStatus || '', 'trajectory-raw-status');
    status.setAttribute('role', 'status');
    raw.addEventListener('toggle', () => {
      if (state.body !== body || !raw.isConnected) return;
      state.rawOpen = raw.open;
    });
    input.addEventListener('change', async () => {
      const file = input.files?.[0];
      if (!file || current !== state) return;
      const selection = ++state.rawSelection;
      state.rawOpen = true;
      state.rawPages = null;
      state.rawPage = 0;
      state.rawStatus = 'Checking SHA-256…';
      const hadFocus = input.ownerDocument.activeElement === input;
      render(state);
      if (hadFocus) state.body.querySelector('.trajectory-raw input')?.focus();
      try {
        const text = await verifyTranscript(file, receipt, root.crypto);
        if (current !== state || selection !== state.rawSelection) return;
        state.rawPages = transcriptPages(text);
        state.rawStatus = 'SHA-256 verified.';
      } catch (error) {
        if (current !== state || selection !== state.rawSelection) return;
        state.rawStatus = error.message;
      }
      const focused = state.body.ownerDocument.activeElement === state.body.querySelector('.trajectory-raw input');
      render(state);
      if (focused) state.body.querySelector(state.rawPages ? '.trajectory-raw pre' : '.trajectory-raw input')?.focus();
    });
    raw.addEventListener('keydown', event => {
      if (event.key === 'Escape' && raw.open) {
        event.stopPropagation();
        raw.open = false;
        raw.querySelector('summary')?.focus();
      } else if (event.target.closest?.('input, button, pre')) {
        event.stopPropagation();
      }
    });
    if (!state.rawPages) return;
    const controls = append(raw, 'div', undefined, 'trajectory-pagination');
    const prev = append(controls, 'button', 'Previous');
    const page = append(controls, 'span');
    page.setAttribute('aria-live', 'polite');
    const next = append(controls, 'button', 'Next');
    prev.type = next.type = 'button';
    const content = append(raw, 'pre');
    content.tabIndex = 0;
    content.setAttribute('aria-label', 'Verified raw transcript');
    function showPage() {
      page.textContent = 'Page ' + (state.rawPage + 1) + ' of ' + state.rawPages.length;
      prev.disabled = state.rawPage === 0;
      next.disabled = state.rawPage === state.rawPages.length - 1;
      content.textContent = state.rawPages[state.rawPage];
      content.scrollTop = 0;
    }
    prev.addEventListener('click', () => { state.rawPage--; showPage(); });
    next.addEventListener('click', () => { state.rawPage++; showPage(); });
    showPage();
  }

  function renderOverview(state) {
    const slot = state.overview;
    if (!slot) return;
    slot.replaceChildren();
    if (state.failure) {
      append(slot, 'p', state.failure, 'muted');
      return;
    }
    if (!state.result) {
      append(slot, 'p', 'Loading live tool activity…', 'muted');
      return;
    }
    const world = state.result.world;
    const calls = Array.isArray(world?.calls) ? world.calls : [];
    const total = Number.isSafeInteger(world?.total) && world.total > 0 ? world.total : calls.length;
    if (!world?.available) {
      append(slot, 'p', 'Tool activity unavailable.', 'muted');
      return;
    }
    if (!calls.length) {
      append(slot, 'p', 'No World tool calls recorded yet.', 'muted');
      return;
    }
    const groups = collapseCalls(calls).slice(-6).reverse();
    const list = append(slot, 'ol', undefined, 'overview-trajectory-list');
    for (const item of groups) renderCallRow(list, item);
    const foot = append(slot, 'div', undefined, 'overview-trajectory-foot');
    append(foot, 'span', total + (total === 1 ? ' World call recorded' : ' World calls recorded'), 'muted');
    const openBtn = append(foot, 'button', 'Open full trajectory →', 'link-btn');
    openBtn.type = 'button';
    openBtn.setAttribute('data-jump-tab', 'trajectory');
  }

  function render(state) {
    renderOverview(state);
    const body = state.body;
    if (!body) return;
    body.replaceChildren();
    if (state.failure && !state.result) {
      append(body, 'p', state.failure, 'muted').setAttribute('role', 'status');
      const retry = append(body, 'button', 'Retry');
      retry.type = 'button';
      retry.addEventListener('click', () => load(state));
      return;
    }
    if (!state.result) {
      append(body, 'p', 'Loading recorded trajectory…', 'muted').setAttribute('role', 'status');
      return;
    }
    const { world, harness, rawTranscript } = state.result;
    renderWorldCalls(body, state, world);
    const harnessNote = !harness?.available && !rawTranscript
      ? 'Available once the interactive session closes and writes its verified transcript.'
      : 'No recorded harness events.';
    if (harness?.available && Array.isArray(harness.events) && harness.events.length) {
      group(body, 'Harness events', harness, harness.events,
        row => ['Launch ' + row.launch, labels[row.kind] || row.kind, shortTime(row.at)].filter(Boolean).join(' · '),
        'No recorded harness events.');
    }
    if (rawReceipt(rawTranscript, state.episode)) rawViewer(body, state, rawTranscript);
  }

  async function load(state, background = false) {
    if (state.request || current !== state) return;
    const request = new AbortController();
    state.request = request;
    if (!background) {
      state.failure = '';
      render(state);
    }
    try {
      const response = await state.getJSON(trajectoryPath(state.brane, state.id), request.signal);
      if (current !== state || request.signal.aborted) return;
      if (response.ok && response.body?.schema === schema && response.body.episode === state.episode) {
        state.result = response.body;
        state.failure = '';
        state.loadedAt = Date.now();
      } else if (!background || !state.result) {
        state.failure = response.status === 404 ? 'No verified episode trajectory is available.' : 'The recorded trajectory could not be loaded.';
      }
    } catch (_) {
      if (current !== state || request.signal.aborted) return;
      if (!background || !state.result) state.failure = 'The recorded trajectory could not be loaded.';
    } finally {
      state.request = null;
    }
    render(state);
  }

  function refresh() {
    if (!current || current.request || !current.getJSON) return;
    if (current.loadedAt && Date.now() - current.loadedAt < 800) return;
    load(current, Boolean(current.result));
  }

  function mount(front, node, getJSON, visible) {
    const body = front?.querySelector('[data-trajectory-body]');
    const overview = front?.querySelector('[data-overview-trajectory]') || null;
    const episode = node?.trajectoryEpisode || node?.run?.episode;
    if ((!body && !overview) || !episode || !node.id || !node.brane) return clear();
    if (!current || current.id !== node.id || current.brane !== node.brane || current.episode !== episode) {
      clear();
      current = { id: node.id, brane: node.brane, episode, rawSelection: 0, loadedAt: 0 };
    }
    const state = current;
    if (body && state.body !== body) {
      body.addEventListener('keydown', event => {
        if (event.target.closest?.('input, button, pre')) event.stopPropagation();
      });
    }
    state.body = body;
    state.overview = overview;
    state.getJSON = getJSON;
    if (!visible) return;
    if (state.result || state.failure || state.request) render(state);
    else load(state);
  }

  root.WorldNodeTrajectory = { mount, clear, refresh };
  if (typeof module === 'object' && module.exports) {
    module.exports = { ...module.exports, trajectoryPath, rawReceipt, verifyTranscript, transcriptPages, collapseCalls };
  }
})(globalThis);
