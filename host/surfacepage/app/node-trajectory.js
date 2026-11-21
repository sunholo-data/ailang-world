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

  function group(body, heading, source, items, line, empty) {
    const section = append(body, 'section', undefined, 'trajectory-source');
    section.setAttribute('aria-label', heading);
    append(section, 'h3', heading);
    if (!source?.available) {
      append(section, 'p', heading + ' unavailable.', 'muted');
      return;
    }
    if (!Array.isArray(items) || !items.length) {
      append(section, 'p', empty, 'muted');
    } else {
      const list = append(section, 'ol');
      for (const row of items) append(list, 'li', line(row));
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

  function render(state) {
    const body = state.body;
    body.replaceChildren();
    if (state.failure) {
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
    group(body, 'World tool calls', world, world?.calls,
      row => [row.tool, row.surface, 'entry ' + row.entry, row.at].filter(Boolean).join(' · '),
      'No recorded World tool calls.');
    group(body, 'Harness events', harness, harness?.events,
      row => ['Launch ' + row.launch, labels[row.kind] || row.kind, row.at].filter(Boolean).join(' · '),
      'No recorded harness events.');
    if (rawReceipt(rawTranscript, state.episode)) rawViewer(body, state, rawTranscript);
  }

  async function load(state) {
    if (state.request || current !== state) return;
    const request = new AbortController();
    state.request = request;
    state.failure = '';
    render(state);
    try {
      const response = await state.getJSON(trajectoryPath(state.brane, state.id), request.signal);
      if (current !== state || request.signal.aborted) return;
      if (response.ok && response.body?.schema === schema && response.body.episode === state.episode) {
        state.result = response.body;
      } else {
        state.failure = response.status === 404 ? 'No verified episode trajectory is available.' : 'The recorded trajectory could not be loaded.';
      }
    } catch (_) {
      if (current !== state || request.signal.aborted) return;
      state.failure = 'The recorded trajectory could not be loaded.';
    } finally {
      state.request = null;
    }
    render(state);
  }

  function mount(front, node, getJSON, visible) {
    const body = front?.querySelector('[data-trajectory-body]');
    const episode = node?.trajectoryEpisode || node?.run?.episode;
    if (!body || !episode || !node.id || !node.brane) return clear();
    if (!current || current.id !== node.id || current.brane !== node.brane || current.episode !== episode) {
      clear();
      current = { id: node.id, brane: node.brane, episode, rawSelection: 0 };
    }
    const state = current;
    if (state.body !== body) {
      body.addEventListener('keydown', event => {
        if (event.target.closest?.('input, button, pre')) event.stopPropagation();
      });
    }
    state.body = body;
    state.getJSON = getJSON;
    if (!visible) return;
    if (state.result || state.failure || state.request) render(state);
    else load(state);
  }

  root.WorldNodeTrajectory = { mount, clear };
  if (typeof module === 'object' && module.exports) {
    module.exports = { ...module.exports, trajectoryPath, rawReceipt, verifyTranscript, transcriptPages };
  }
})(globalThis);
