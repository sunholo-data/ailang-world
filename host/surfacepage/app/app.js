// World: the surface app (row 188 N2b). The approved v3 working copy, ported onto live data: it
// renders GET /surface/model.json and GET /surface/stack.json (N1's world/surface-model/v1 and
// world/surface-stack/v1), following POST /agui/ from the model's cursor. A field the records do not hold stays absent and
// shows as an honest gap. The data adapters, the layout maths and the brane card are pure and
// exported for node, as /surface/surface.js's seams are; the address grammar and the stored
// look are those seams (window.WorldSurface).
(() => {
  'use strict';
  const esc = s => String(s == null ? '' : s).replace(/[&<>"]/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c]));
  const MON = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  const DOW = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  const pad = n => String(n).padStart(2, '0');
  const hm = iso => { if (iso == null) return ''; const d = new Date(iso); return Number.isFinite(+d) ? pad(d.getUTCHours()) + ':' + pad(d.getUTCMinutes()) : ''; };
  const dayOf = iso => { if (iso == null) return ''; const d = new Date(iso); return Number.isFinite(+d) ? d.getUTCDate() + ' ' + MON[d.getUTCMonth()] : ''; };
  const exactOf = iso => dayOf(iso) ? dayOf(iso) + ' ' + hm(iso) + ' UTC' : 'Time unrecorded';
  const lastRunLabel = run => run ? (run.failed ? 'Last run failed' : 'Last run finished') : '';
  const lastRunDetail = run => run ? lastRunLabel(run) + ' · ' + exactOf(run.endedAt) : '';
  function mergeFileSummary(node) {
    const n = node.landedFileCount;
    if (!Number.isSafeInteger(n) || n < 0) return '';
    if (n === 0) return 'No file changes in merge';
    return n + (n === 1 ? ' file changed in merge' : ' files changed in merge');
  }
  const prettyModel = m => (m || '').split('-').filter(Boolean).reduce((a, w, i, all) => {
    if (/^\d+$/.test(w) && /^\d+$/.test(all[i - 1] || '')) return a + '.' + w;
    const word = /^gpt$/i.test(w) ? 'GPT' : /^\d/.test(w) ? w : w[0].toUpperCase() + w.slice(1);
    return a + (a ? ' ' : '') + word;
  }, '');
  const cap = s => s ? s[0].toUpperCase() + s.slice(1) : '';
  const agentLine = n => { if (!n.model) return ''; const m = prettyModel(n.model); return m.toLowerCase().startsWith(String(n.agent || '').toLowerCase()) ? m : `${cap(n.agent)} · ${m}`; };
  // PR titles carry rig bookkeeping ("(through World, w189c)", "— made through World (w190)"): keep what was done.
  const cleanTitle = t => String(t || '').replace(/`/g, '')
    .replace(/\s*[—–-]\s*(?:made|built) (?:through|inside) World\b.*$/i, '')
    .replace(/\s*\([^()]*\b(?:through World|made through World|built inside World|attended|trial w\d)[^()]*\)/gi, '')
    .replace(/\s+:/g, ':').replace(/\s{2,}/g, ' ').trim();
  const shortTitle = (t, max = 44) => { const c = cleanTitle(t); if (c.length <= max) return c; const cut = c.slice(0, max); return cut.slice(0, Math.max(cut.lastIndexOf(' '), max - 10)).replace(/[\s,;:—–-]+$/, '') + '…'; };
  const summaryOf = n => String(n.summary || '').replace(/\s+-\s+/g, ' ').replace(/\s{2,}/g, ' ').trim();
  // A lane's label: its name and, for a row whose name is not already "row N", the row.
  const rowLabel = (name, row) => row && name !== 'row ' + row ? `${name} · row ${row}` : String(name || '');
  // A shared world's lane: "shared" in v1, "<repo>:brane:shared" in follow mode (row 220 M4).
  const isSharedKey = k => k === 'shared' || /^[a-z][a-z0-9-]*:brane:shared$/.test(String(k == null ? '' : k));

  // ---------- Glyphs (one shape per brane; colour from the data) ----------
  const CLIP = { square: 'none', triangle: 'polygon(50% 4%, 100% 96%, 0 96%)', diamond: 'polygon(50% 0, 100% 50%, 50% 100%, 0 50%)',
    hexagon: 'polygon(25% 4%, 75% 4%, 100% 50%, 75% 96%, 25% 96%, 0 50%)', pentagon: 'polygon(50% 0, 100% 38%, 81% 100%, 19% 100%, 0 38%)',
    star: 'polygon(50% 0, 62% 35%, 99% 36%, 69% 58%, 80% 95%, 50% 73%, 20% 95%, 31% 58%, 1% 36%, 38% 35%)' };
  const SHARED_INK = '#526677';
  function gl(shape, colour, size, extra = '') {
    if (!shape) return `<span class="g" style="width:${size - 2}px;height:${size - 2}px;border-radius:50%;background:#8496a5;margin:1px;${extra}"></span>`;
    if (shape === 'circle') return `<span class="g" style="width:${size}px;height:${size}px;border-radius:50%;border:2px solid ${colour || SHARED_INK};${extra}"></span>`;
    return `<span class="g" style="width:${size}px;height:${size}px;background:${colour};clip-path:${CLIP[shape] || 'none'};border-radius:${shape === 'square' ? '2px' : '0'};${extra}"></span>`;
  }

  const RECORD_LABEL = {
    'github-merge': 'GitHub merge', 'shared-commit': 'Shared commit',
    episode: 'World episode', draft: 'World draft', step: 'Work step',
  };
  const mergeGlyph = (size, timeline = false) => `<svg class="g merge-glyph" viewBox="0 0 16 16" width="${size}" height="${size}" aria-hidden="true" focusable="false"${timeline ? ' style="color:var(--text)"' : ''}>
      <path d="M2 4h3c3 0 3 4 6 4h3M2 12h3c3 0 3-4 6-4" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
    </svg>`;

  // ---------- Time axis: one slot per instant, shared by both views ----------
  const BREAK_H = 72;                           // longer quiet gaps are drawn as an axis break, not as empty weeks
  const LEVELS = [{ name: 'Close', G: 84, P: 0 }, { name: 'Steps', G: 50, P: 0 }, { name: 'Hours', G: 14, P: 16 }, { name: 'Days', G: 3, P: 6 }, { name: 'Weeks', G: 0, P: 2.4 }];
  const breakAt = (TT, i) => i > 0 && (TT[i] - TT[i - 1]) / 3.6e6 > BREAK_H;
  // levelXs blends evenly spaced instants into elapsed time as you zoom out.
  function levelXs(TT, presentIndex = TT.length - 1, undated = false) {
    const NS = TT.length - 1;
    return LEVELS.map(L => {
      const xs = [220];
      for (let i = 1; i <= NS; i++) {
        const dth = (TT[i] - TT[i - 1]) / 3.6e6;
        const gap = (undated && i === 1) || breakAt(TT, i) ? Math.max(L.G, 64) : i === presentIndex ? Math.max(L.G * 1.44, dth * L.P, 28) : Math.max(L.G, dth * L.P);
        xs.push(xs[i - 1] + gap);
      }
      return xs;
    });
  }
  // blendXS is the x of every point at zoom u, between two levels.
  function blendXS(LX, u) { const k = Math.min(LX.length - 2, Math.floor(u)), f = u - k; return LX[k].map((x, i) => x + (LX[k + 1][i] - x) * f); }

  function timeAxis(nodes, now) {
    const entries = nodes.map(n => [n, n.draft ? NaN : Date.parse(n.at)]);
    const times = entries.map(([, t]) => t).filter(Number.isFinite);
    const measured = Date.parse(now);
    const nowKnown = Number.isFinite(measured);
    const TNOW = nowKnown ? measured : times.reduce((latest, time) => Math.max(latest, time), -Infinity);
    const present = Number.isFinite(TNOW) ? TNOW : 0;
    const TT = [...new Set([...times, present])].sort((a, b) => a - b);
    const undated = entries.some(([n, t]) => !n.draft && !Number.isFinite(t));
    // Undated states get one slot outside recorded history. Its ordering key
    // is synthetic and never shown as a timestamp. Local drafts sit at now.
    if (undated) TT.unshift(TT[0] - 1);
    const indexes = new Map(TT.map((t, i) => [t, i]));
    const presentIndex = indexes.get(present);
    const timeIndex = new Map(entries.map(([n, t]) => [n.id, n.draft ? presentIndex : indexes.get(t) ?? 0]));
    const LX = levelXs(TT, presentIndex, undated);
    return { TT, LX, timeIndex, undated, nowKnown, dated: nowKnown || times.length > 0, presentIndex, TNOW:present, NOW_D:new Date(present), NS:TT.length - 1 };
  }

  // ---------- Canvas geometry (1600 × 936 design canvas, scaled to the window) ----------
  let PL = 262, PR = 1262, PW = PR - PL;
  const TICK_T = 768, MINI_T = 818, MINI_H = 36, LEGEND_T = 892;

  // ---------- Model: world/surface-model/v1 adapted to the timeline's lookups ----------
  function adaptModel(D) {
    const LANES = {}, N = {};
    (D.lanes || []).forEach(l => { LANES[l.key] = l; (l.states || []).forEach(s => { N[s.id] = Object.assign({}, s, { lane: l.key }); }); });
    const parents = {}, kids = {};
    const addEdge = (p, c) => { (kids[p] = kids[p] || []).push(c); (parents[c] = parents[c] || []).push(p); };
    (D.links || []).forEach(k => { if (N[k.from] && N[k.to]) addEdge(k.from, k.to); }); // links run parent (from) → child (to)
    // Groups: follow mode's repos in config order, each with its own shared world and branes; else the one world.
    const repos = Array.isArray(D.repos) ? D.repos.filter(r => r.shared && LANES[r.shared]) : [];
    const GROUPS = repos.length ? repos.map(r => ({ id: r.id, shared: r.shared, repo: r })) : [{ id: '', shared: 'shared', repo: null }];
    const multi = GROUPS.length > 1;
    const groupOf = key => { const l = LANES[key]; return (l && l.repo && GROUPS.find(g => g.id === l.repo)) || GROUPS[0]; };
    const dipOf = g => multi ? 'dip|' + g.id : 'dip', otherOf = g => multi ? 'other|' + g.id : 'other';
    const SH = Object.values(N).filter(n => n.kind === 'shared').sort((a, b) => a.rank - b.rank).map(n => n.id);
    GROUPS.forEach(g => {   // each shared world's own first-parent line; no edge joins two repos
      const line = SH.filter(id => N[id].lane === g.shared);
      for (let i = 1; i < line.length; i++) addEdge(line[i - 1], line[i]);
    });
    const nameOf = key => { const l = LANES[key]; return l ? l.name || l.title || key : String(key); };
    const laneTitle = key => rowLabel(nameOf(key), LANES[key] && LANES[key].row);
    const glyphOf = (key, size) => { const l = LANES[key] || {}; return gl(l.glyph, l.glyph === 'circle' ? SHARED_INK : l.colour, size); };
    const recordKind = n => {
      if (n.kind === 'shared') return n.pr && parents[n.id]?.some(id => N[id]?.kind === 'step')
        ? 'github-merge' : 'shared-commit';
      if (n.draft) return 'draft';
      return n.kind === 'open' || n.episode ? 'episode' : 'step';
    };
    const recordLabel = n => RECORD_LABEL[recordKind(n)];
    const activityLabel = n => D.nowEdge?.activityUntil?.[n.id] ? 'Tool activity observed at ' + hm(D.now) + ' UTC' : '';
    const glyphOfNode = (n, size) => {
      if (n.kind === 'shared') return recordKind(n) === 'github-merge' ? mergeGlyph(size) : gl('circle', SHARED_INK, size);
      if (n.kind === 'open') {
        const colour = LANES[n.lane]?.colour || SHARED_INK;
        return gl('circle', colour, size, 'border-style:dashed;background:' + (n.open?.sessionLive ? colour + '55' : 'var(--shellBg)') + ';');
      }
      return isSharedKey(n.lane) ? gl(null, null, size) : glyphOf(n.lane, size);
    };
    const mergeOf = id => (kids[id] || []).find(c => N[c].kind === 'shared');
    const nodes = Object.values(N);
    const STEPS = nodes.filter(n => n.kind === 'step');
    const OPEN = nodes.filter(n => n.kind === 'open');
    const axis = timeAxis(nodes, D.now);
    // Rows: the most recent brane above the shared world, the next three below, the rest share one quiet lane.
    // Several repos: each repo's shared world, then its newest branes, then its quiet lane, in config order.
    const recent = {};
    nodes.filter(n => n.kind !== 'shared').forEach(n => { if (!isSharedKey(n.lane)) recent[n.lane] = Math.max(recent[n.lane] == null ? -1 : recent[n.lane], n.rank); });
    const byRecent = Object.keys(recent).sort((a, b) => recent[b] - recent[a]);
    const LY = {}, QUIET = {}, HOME = {};
    let above = [], below = [], quiet = [], ROWS = [], GAP = 98;
    if (!multi) {
      const g = GROUPS[0];
      above = byRecent.slice(0, 1); below = byRecent.slice(1, 4);
      quiet = Object.keys(LANES).filter(k => !isSharedKey(k) && !above.includes(k) && !below.includes(k));
      above.forEach(k => { LY[k] = 182; }); LY[g.shared] = 320; below.forEach((k, i) => { LY[k] = 440 + i * 98; }); LY[otherOf(g)] = 440 + below.length * 98 - 6;
      ROWS = above.concat(LANES[g.shared] ? [g.shared] : [], below, quiet.length ? [otherOf(g)] : []);
      QUIET[otherOf(g)] = quiet; HOME[dipOf(g)] = g.shared;
    } else {
      const per = GROUPS.length === 2 ? 2 : GROUPS.length === 3 ? 1 : 0;
      GROUPS.forEach(g => {
        const shown = byRecent.filter(k => groupOf(k) === g).slice(0, per);
        const q = Object.keys(LANES).filter(k => groupOf(k) === g && !isSharedKey(k) && !shown.includes(k));
        ROWS.push(g.shared, ...shown); below.push(...shown); quiet.push(...q);
        if (q.length) { ROWS.push(otherOf(g)); QUIET[otherOf(g)] = q; }
        HOME[dipOf(g)] = g.shared;
      });
      GAP = Math.min(98, (716 - 160) / Math.max(1, ROWS.length - 1));
      ROWS.forEach((r, i) => { LY[r] = 160 + i * GAP; });
    }
    const rowOf = id => { const n = N[id], g = groupOf(n.lane); return n.kind === 'shared' ? g.shared : isSharedKey(n.lane) ? dipOf(g) : LY[n.lane] != null ? n.lane : otherOf(g); };
    const rowY = r => HOME[r] ? LY[HOME[r]] + (multi ? Math.round(GAP * 0.45) : 44) : LY[r];
    const isDip = r => !!HOME[r], isOther = r => !!QUIET[r] || r === 'other';
    const homeRow = r => HOME[r] || r;
    const lastStep = key => STEPS.filter(n => n.lane === key).sort((a, b) => b.rank - a.rank)[0];
    const lastState = key => ((LANES[key] && LANES[key].states) || []).reduce((a, s) => !a || s.rank > a.rank ? s : a, null);
    const aliasOf = {}; (D.aliases || []).forEach(a => { if (N[a.target]) aliasOf[a.id] = a.target; });
    return { D, LANES, N, parents, kids, SH, ...axis, MAXU: axis.LX.length - 1, STEPS, OPEN, above, below, quiet, LY, ROWS,
      GROUPS, multi, GAP, QUIET, groupOf, isDip, isOther, homeRow, aliasOf,
      nameOf, laneTitle, glyphOf, glyphOfNode, recordKind, recordLabel, activityLabel, mergeOf, rowOf, rowY, lastStep, lastState };
  }

  // ---------- Layout of the visible graph at one zoom (XS) ----------
  const isCluster = id => typeof id === 'string' && id.startsWith('c:');
  function layoutGraph(T, XS, selectedID) {
    const { N, timeIndex } = T, POS = {}, REP = {}, DISP = {};
    Object.keys(N).forEach(id => {
      POS[id] = [XS[timeIndex.get(id)], T.rowY(T.rowOf(id))];
    });
    // Collapse dense runs on each row into a count bubble.
    const byRow = {};
    Object.keys(POS).forEach(id => { (byRow[T.rowOf(id)] = byRow[T.rowOf(id)] || []).push(id); });
    Object.entries(byRow).forEach(([row, ids]) => {
      ids.sort((a, b) => POS[a][0] - POS[b][0] || N[a].rank - N[b].rank);
      const SEP = isSharedKey(row) ? 15 : 19;
      let g = [];
      const flush = () => {
        if (!g.length) return;
        const separate = g.length > 1 && g.includes(selectedID);
        if (separate) {
          REP[selectedID] = selectedID;
          DISP[selectedID] = {x:POS[selectedID][0], y:POS[selectedID][1], ids:[selectedID], row};
          g = g.filter(id => id !== selectedID);
        }
        const rep = g.length === 1 ? g[0] : 'c:' + g[0];
        const x = g.reduce((sum, id) => sum + POS[id][0], 0) / g.length;
        const y = POS[g[0]][1] + (separate && Math.abs(x - POS[selectedID][0]) < 24 ? 24 : 0);
        g.forEach(id => { REP[id] = rep; });
        DISP[rep] = {x, y, ids:g, row};
        g = [];
      };
      ids.forEach(id => { if (g.length && POS[id][0] - POS[g[g.length - 1]][0] >= SEP) flush(); g.push(id); });
      flush();
    });
    // Labels thin out: newest first, one per ~112 px on a row.
    const LABELED = new Set();
    const lastX = {};
    Object.keys(DISP).filter(id => !isCluster(id)).sort((a, b) => DISP[b].x - DISP[a].x).forEach(id => {
      const d = DISP[id], n = N[id];
      if (n.kind === 'shared' || T.isDip(d.row)) return;
      if (lastX[d.row] == null || lastX[d.row] - d.x >= 112) { LABELED.add(id); lastX[d.row] = d.x; }
    });
    return { POS, REP, DISP, LABELED };
  }

  // ---------- Brane view: world/surface-stack/v1 as bars, the front card and the rail ----------
  const tagOf = c => c.kind === 'shared' ? `Shared ${c.pr ? '#' + c.pr : c.label}` : (c.episode || c.label);
  function stackBars(cols, dir) {
    const shown = (cols || []).slice(0, 3), nc = Math.max(1, shown.length), cw = (1180 - (nc - 1) * 12) / nc, out = [];
    shown.forEach((col, c) => {
      const near = dir < 0 ? col.cards.slice().reverse() : col.cards;
      near.slice(0, 3).forEach((card, i) => {
        const dd = i + 1, shrink = (dd - 1) * 24;
        out.push({ key: (dir < 0 ? 'u:' : 'd:') + card.id, card, col, dd,
          width: cw - shrink, ml: -590 + c * (cw + 12) + shrink / 2, top: dir < 0 ? 106 - (dd - 1) * 46 : 796 + (dd - 1) * 46,
          op: Math.max(0.35, 1 - (dd - 1) * 0.28), z: 39 - dd, dir, laneLab: nc > 1 && dd === 1 ? (col.brane === 'shared' ? 'Shared world' : col.name || col.brane) : '' });
      });
    });
    return out;
  }
  function barHTML(b) {
    const c = b.card, g = c.kind === 'shared' ? gl('circle', SHARED_INK, 11) : isSharedKey(b.col.brane) ? gl(null, null, 11) : gl(b.col.glyph, b.col.colour, 11);
    const t = cleanTitle(c.title);
    return `<span class="a">${b.dir < 0 ? '↑' : '↓'}</span>${g}<span class="tag">${esc(tagOf(c))}</span>${t ? `<span class="ttl">${esc(t)}</span>` : ''}<span style="flex:1"></span><span class="w">${esc(c.rel)}</span>${c.badge > 0 ? `<span class="nb inl">${c.badge}</span>` : ''}${b.laneLab ? `<span class="ln">${esc(b.laneLab)}</span>` : ''}`;
  }

  const chip = (id, label, glyphHtml, via, mono, tip, extra = '') => `<button type="button" class="chip${mono ? ' mono-l' : ''}" data-go="${esc(id)}"${tip ? ` data-tip="${esc(tip)}"` : ''}${extra}>${glyphHtml}<span class="lb">${esc(label)}</span>${via ? `<span class="via">${esc(via)}</span>` : ''}</button>`;
  // stateChip names a state from the stack, glyphed by its lane in the model when the model has it.
  const stateChip = (T, s, via) => {
    if (!s) return '';
    const n = (T && T.N[s.id]) || s;
    const g = n.kind === 'shared' ? gl('circle', SHARED_INK, 10) : T && n.lane && !isSharedKey(n.lane) ? T.glyphOf(n.lane, 10) : gl(null, null, 10);
    return chip(n.id, n.kind === 'shared' ? (n.pr ? `#${n.pr}` : n.label) : (n.episode || n.label), g, via, false, 'node:' + n.id);
  };
  const ACT = { you: ['var(--you)', 'var(--onYou)', 'var(--you)'], agent: ['#eaf1f7', '#2f6690', '#b9cfe1'], world: ['#f6f8fa', '#526677', '#dce3e8'], review: ['#f1edf6', '#6b4f8a', '#d6cbe4'], loop: ['#eaf3ee', '#3f7a5e', '#c3dccf'] };
  const WHO = { you: 'You', agent: 'Agent', world: 'World', review: 'Review', loop: 'Loop' };

  function verdict(graders, rounds, firstPass, minutes) {
    if (!graders) return '';
    const [a, b] = String(graders).split('/').map(Number);
    let s = a === b ? `Passed all ${b} checks` : `${a} of ${b} checks passed`;
    s += firstPass ? ' on the first round' : rounds > 1 ? ` after ${rounds} rounds` : rounds === 1 ? ' in one round' : '';
    return s + (minutes != null ? ` · ${minutes} agent-minutes` : '');
  }

  const FILE_STATUS = new Map([['A', 'Added'], ['M', 'Modified'], ['D', 'Deleted']]);

  function pathAttribute(path) {
    // HTML normalizes literal CR before tokenizing attributes.
    return esc(path).replace(/\r/g, '&#13;');
  }

  function filePathHTML(path, source, origin = '') {
    const name = `<code>${esc(path)}</code>`;
    if (!source) return name;
    const version = source === 'episode-read' ? 'recorded read' : source === 'episode-write' ? 'recorded write' : '';
    return `<button type="button" class="file-path" data-file-path="${pathAttribute(path)}" data-file-source="${source}"${origin ? ` data-file-origin="${esc(origin)}"` : ''}${version ? ` aria-label="Preview ${version} of ${pathAttribute(path)}"` : ''}>${name}</button>`;
  }

  function filesTabHTML(F, activity) {
    const groups = [], options = [];
    const episode = F.trajectoryEpisode || F.run?.episode;
    const addGroup = (kind, title, rows, describe, source = '') => {
      if (!rows.length) return;
      const list = rows.map((row, index) => {
        const metadata = describe(row);
        return `<li class="file-change file-activity-row" data-file-row="${pathAttribute(row.path)}">
          <span class="file-activity-meta">${esc(metadata)}</span>${filePathHTML(row.path, source, kind + ':' + index)}</li>`;
      }).join('');
      options.push(`<option value="${kind}">${title}</option>`);
      groups.push(`<section class="file-activity-group" data-file-group="${kind}"><h3>${title} <span>${rows.length}</span></h3><ul class="file-list">${list}</ul></section>`);
    };
    const mentions = (row, word) => Number.isSafeInteger(row.count) && row.count > 0 ? `${row.count} ${word}${row.count === 1 ? '' : 's'}` : '';
    const writeActions = new Map([['workspace-write', 'write'], ['workspace-edit', 'edit'], ['write', 'write'], ['edit', 'edit'], ['patch', 'patch']]);
    const writeDescription = row => {
      const action = writeActions.get(row.op);
      return [mentions(row, 'write'), action ? 'Last: ' + action : ''].filter(Boolean).join(' · ');
    };
    let message = '';
    if (episode) {
      if (activity == null) {
        message = '<p class="muted" role="status">Loading episode file activity…</p>';
      } else if (activity.status === 'complete' || activity.status === 'partial') {
        const reads = activity.reads || [], writes = activity.writes || [], produced = activity.produced || [];
        addGroup('reads', 'Read', reads, row => mentions(row, 'read'), 'episode-read');
        addGroup('writes', 'Written', writes, writeDescription, 'episode-write');
        addGroup('produced', 'Changes observed after commands', produced, row => mentions(row, 'observation'));
        if (activity.status === 'partial') {
          const omitted = Object.values(activity.omitted || {}).reduce((a,n) => a + (Number(n) || 0), 0);
          message = `<p class="muted">Episode file activity is partial.${omitted ? ' Some path records were omitted.' : ''}${activity.deltaGaps ? ' Some command changes could not be enumerated.' : ''}</p>`;
        } else if (!reads.length && !writes.length && !produced.length) {
          message = '<p class="muted">No file activity was recorded for this episode.</p>';
        } else if (reads.length && !writes.length && !produced.length) {
          message = '<p class="muted">No writes recorded.</p>';
        }
      } else {
        message = '<p class="muted">Episode file activity is unavailable.</p>';
      }
    }

    const landed = F.landedFiles || [];
    const merge = F.kind === 'step' || F.parents?.length > 1 || F.commits?.some(commit => commit.kinds?.includes('merge'));
    addGroup('committed', merge ? 'In the merge' : 'In the commit', landed, row => FILE_STATUS.get(row.status) || 'Changed', 'committed');
    const historical = !episode || activity?.status === 'unavailable';
    const recordedReads = historical ? (F.reads || []).flatMap(group => (group.files || []).map(file => ({...file, from: group.label}))) : [];
    addGroup('recorded-reads', 'Recorded step reads', recordedReads, row => [mentions(row, 'read'), row.from].filter(Boolean).join(' · '));
    const recordedWrites = historical && Array.isArray(F.files);
    if (recordedWrites) addGroup('recorded', 'Recorded step changes', F.files, row => FILE_STATUS.get(row.status) || 'Changed');
    const missing = Math.max(0, (F.landedFileCount || 0) - landed.length);
    const landedGap = missing ? `<p class="muted">${missing} committed file path${missing === 1 ? '' : 's'} unavailable.</p>`
      : F.landedFileCount === 0 ? '<p class="muted">No committed file changes.</p>' : '';

    const recordedGap = recordedWrites
      ? F.filesTruncated > 0 ? `<p class="muted">${F.filesTruncated} recorded step file path${F.filesTruncated === 1 ? '' : 's'} unavailable.</p>`
        : !F.files.length ? '<p class="muted">No recorded step file changes.</p>' : ''
      : '';
    const filters = groups.length ? `<div class="file-activity-filters">
        <label>Filter paths <input type="search" data-file-filter spellcheck="false" autocomplete="off"></label>
        ${options.length > 1 ? `<label>Activity <select data-file-kind><option value="">All activity</option>${options.join('')}</select></label>` : ''}
      </div><p class="muted" data-file-empty role="status" hidden>No paths match these filters.</p>` : '';
    return `<div class="node-files-tab">${filters}${message}${landedGap}${recordedGap}${groups.join('')}${!groups.length && !message && !landedGap && !recordedGap ? '<p class="muted">No file records are available for this node.</p>' : ''}</div>`;
  }

  // Render the recorded request as a document. Text is escaped before any markup is added.
  // Links, images and embedded HTML remain inert in an agent-supplied request.
  function taskDocument(text) {
    const inline = line => line.split(/(`[^`\n]+`|\*\*[^*\n]+\*\*)/g).map(part => {
      if (part.startsWith('`') && part.endsWith('`') && part.length > 2) return '<code>' + esc(part.slice(1,-1)) + '</code>';
      if (part.startsWith('**') && part.endsWith('**') && part.length > 4) return '<strong>' + esc(part.slice(2,-2)) + '</strong>';
      return esc(part);
    }).join('');
    const lines = text.replace(/\r\n?/g, '\n').split('\n'), blocks = [];
    const fence = /^ {0,3}(`{3,}|~{3,})(.*)$/;
    const heading = /^ {0,3}(#{1,6})\s+(.*)$/;
    const list = /^\s*(?:(\d+)\.|[-*+])\s+(.*)$/;
    for (let i = 0; i < lines.length;) {
      const line = lines[i++];
      if (!line.trim()) continue;
      const marker = line.match(fence);
      if (marker) {
        const code = [];
        while (i < lines.length) {
          const next = lines[i++], closing = next.match(fence);
          if (closing && closing[1][0] === marker[1][0] && closing[1].length >= marker[1].length && !closing[2].trim()) break;
          code.push(next);
        }
        blocks.push('<pre><code>' + esc(code.join('\n')) + '</code></pre>');
        continue;
      }
      const h = line.match(heading);
      if (h) {
        const level = Math.min(6, 3 + h[1].length);
        blocks.push('<h' + level + '>' + inline(h[2]) + '</h' + level + '>');
        continue;
      }
      const item = line.match(list);
      if (item) {
        const tag = item[1] ? 'ol' : 'ul';
        const items = [item[2]];
        while (i < lines.length) {
          const next = lines[i].match(list);
          if (!next || Boolean(next[1]) !== Boolean(item[1])) break;
          items.push(next[2]); i++;
        }
        const start = item[1] && item[1] !== '1' ? ' start="' + Number(item[1]) + '"' : '';
        blocks.push('<' + tag + start + '>' + items.map(v => '<li>' + inline(v) + '</li>').join('') + '</' + tag + '>');
        continue;
      }
      const paragraph = [line];
      while (i < lines.length && lines[i].trim() && !heading.test(lines[i]) && !list.test(lines[i]) && !fence.test(lines[i])) paragraph.push(lines[i++]);
      blocks.push('<p>' + paragraph.map(inline).join('\n') + '</p>');
    }
    return blocks.join('');
  }

  function taskHTML(F) {
    const text = F.taskBody || F.task;
    if (!text) return '';
    return `<section class="task-request"><h3>${F.taskBody ? 'Task request' : 'Task title'}</h3>
      <div class="task-document">${taskDocument(text)}</div>
      ${F.taskRedacted ? '<p class="muted">Sensitive task text was redacted.</p>' : ''}</section>`;
  }

  function commitHTML(F) {
    if (!F.commits?.length) return F.kind === 'open' && !F.pr ? '<div class="muted">No finish commit recorded for this episode.</div>' : '';
    return F.commits.map(c => {
      const kind = c.kinds.includes('episode-finish') ? 'Episode finish commit' : c.kinds.includes('branch-tip') ? 'Branch tip commit' : 'Merge commit';
      return `<section class="commit-message"><h3>${kind} <code>${esc(c.id.split(':').pop().slice(0, 8))}</code></h3>
        <p class="commit-subject">${esc(c.subject)}</p>${c.body ? `<div class="commit-body">${esc(c.body)}</div>` : ''}
        ${c.redacted ? '<p class="muted">Sensitive commit text was redacted.</p>' : ''}</section>`;
    }).join('');
  }

  function recordedHeadingHTML(F, title, summary) {
    const rawTitle = F.fullTitle || F.title || '';
    const rawSummary = F.fullSummary || F.summary || '';
    const fullTitle = rawTitle && (rawTitle !== title || rawTitle.length > 60);
    const fullSummary = rawSummary && (rawSummary !== summary || rawSummary.length > 100);
    const titleLabel = F.kind === 'shared' ? (F.pr ? 'Recorded merge title' : 'Recorded commit title') : 'Recorded title';
    if (!fullTitle && !fullSummary) return '';
    return `<details class="stack-disclosure recorded-heading"><summary>${fullTitle && fullSummary ? titleLabel + ' and summary' : fullTitle ? titleLabel : 'Recorded summary'}</summary>
      ${fullTitle ? `<div class="recorded-text">${esc(rawTitle)}</div>` : ''}
      ${fullSummary ? `<div class="recorded-text">${esc(rawSummary)}</div>` : ''}</details>`;
  }

  function runTime(seconds) {
    if (seconds < 60) return seconds + ' s';
    const h = Math.floor(seconds / 3600), m = Math.floor(seconds / 60) % 60;
    return h ? h + ' h ' + m + ' min' : m + ' min ' + (seconds % 60) + ' s';
  }

  function permissionsHTML(permissions) {
    if (!permissions) return '';
    const sessions = permissions.sessions || [];
    const byScope = new Map();
    let overlapping = false;
    const names = { 'Workspace.Read': 'Read files', 'Workspace.Write': 'Write files', 'Workspace.Exec': 'Run commands' };
    const label = effect => Object.hasOwn(names, effect) ? names[effect] : effect;
    const number = n => n.toLocaleString('en-US');
    const key = g => JSON.stringify([g.effect, g.scopeKey ?? g.scope]);
    // One episode spends against a shared counter for each permission and scope.
    // Simultaneous credentials offer alternative ceilings, never additive credit.
    for (const session of sessions) {
      for (const grant of session.grants || []) {
        const id = key(grant), prev = byScope.get(id), active = (grant.state || session.state) === 'active';
        if (prev?.active && active) overlapping = true;
        if (!prev || (active && !prev.active) || (active === prev.active &&
            (active ? grant.budget > prev.grant.budget : (session.createdAt || '') > (prev.session.createdAt || '')))) {
          byScope.set(id, {grant, session, active});
        }
      }
    }
    for (const grant of permissions.observed || []) {
      if (!byScope.has(key(grant))) byScope.set(key(grant), {grant, session: {}, active:false});
    }
    const list = [...byScope.values()].sort((a,b) => a.grant.effect.localeCompare(b.grant.effect) || a.grant.scope.localeCompare(b.grant.scope));
    const stateLabel = state => ({active:'Active',expired:'Expired',revoked:'Revoked',unavailable:'Grant details unavailable',observed:'Historical use','revoked-or-unavailable':'Credential unavailable'})[state] || 'Recorded';
    const grantRows = list.map(({grant:g,session,active}) => {
      const metrics = [];
      metrics.push(g.used != null ? number(g.used) + ' used' : 'Usage unavailable');
      if (g.budget != null) metrics.push('limit ' + number(g.budget));
      if (g.reserved) metrics.push(number(g.reserved) + ' reserved');
      if (active && g.remaining != null) metrics.push(number(g.remaining) + ' left');
      if (!active) {
        if (g.ledgerRemaining != null) metrics.push(number(g.ledgerRemaining) + ' unused');
        metrics.push(stateLabel(g.state || session.state));
      }
      return `<li><span><b>${esc(label(g.effect))}</b> <small class="muted">${esc(g.scope)}</small></span><span>${esc(metrics.join(' · '))}</span></li>`;
    }).join('');
    const history = sessions.map(session => {
      const when = [];
      if (session.createdAt) when.push('Issued ' + exactOf(session.createdAt));
      if (session.expiresAt) when.push((session.state === 'expired' ? 'expired ' : 'expires ') + exactOf(session.expiresAt));
      const grants = (session.grants || []).map(g => {
        const expires = g.expiresAt && g.expiresAt !== session.expiresAt ? ' · expires ' + exactOf(g.expiresAt) : '';
        return `<li><span>${esc(g.effect)} · ${esc(g.scope)}</span><span>${g.budget != null ? 'Limit ' + esc(number(g.budget)) : 'Limit unavailable'}${esc(expires)}</span></li>`;
      }).join('');
      return `<section class="sect"><b>${esc(stateLabel(session.state))}</b>${when.length ? `<span class="muted">${esc(when.join(' · '))}</span>` : ''}${grants ? `<ol>${grants}</ol>` : ''}</section>`;
    }).join('');
    return `<div class="node-decisions node-permissions">
      ${grantRows ? `<ol class="permission-totals">${grantRows}</ol>` : '<p class="muted">Grant details unavailable.</p>'}
      ${history ? `<details class="stack-disclosure" data-detail="grant-history"><summary>Grant history · ${sessions.length}</summary>${overlapping ? '<p class="muted">Active grants for the same permission and scope share one usage counter. The summary shows the highest active limit. Each credential enforces its own limit.</p>' : ''}${history}</details>` : ''}
      ${permissions.truncated ? '<p class="muted">Only part of the grant history is available.</p>' : ''}</div>`;
  }

  function inspectionHTML(F) {
    const run = F.run, github = F.githubChecks;
    const section = (label, body) => `<section class="col node-run-section"><h2 class="colh">${esc(label)}</h2>${body}</section>`;
    const rows = items => `<dl class="node-facts">${items.map(([k,v,note]) => `<dt>${esc(k)}</dt><dd>${esc(v)}${note ? `<div class="muted">${esc(note)}</div>` : ''}</dd>`).join('')}</dl>`;
    let html = '';
    if (run) {
      const harness = [prettyModel(run.harness), prettyModel(run.model)].filter(Boolean).join(' · ');
      const metrics = [['Harness', harness || 'Unrecorded'], ['Recorded launches', run.launches]];
      if (run.launches) {
        metrics.push(['Completed turns', run.turns ?? 'Unreported']);
        if (run.tokens) {
          const tokens = run.tokens, count = n => n.toLocaleString('en-US');
          const cache = [];
          if (tokens.cacheRead != null) cache.push(count(tokens.cacheRead) + ' read from cache');
          if (tokens.cacheWrite != null) cache.push(count(tokens.cacheWrite) + ' written to cache');
          if (tokens.input != null) metrics.push(['Input tokens', count(tokens.input), cache.join(' · ')]);
          else if (cache.length) metrics.push(['Input cache', cache.join(' · ')]);
          if (tokens.output != null) metrics.push(['Output tokens', count(tokens.output), tokens.reasoning != null ? count(tokens.reasoning) + ' reasoning' : '']);
          else if (tokens.reasoning != null) metrics.push(['Reasoning tokens', count(tokens.reasoning)]);
        }
        metrics.push(['Provider cost', run.costUsd == null ? 'Unreported' : new Intl.NumberFormat('en-US', {style:'currency',currency:'USD',maximumFractionDigits:6}).format(run.costUsd)]);
      }
      if (run.runtime != null) metrics.push(['Harness runtime',runTime(run.runtime)]);
      if (run.elapsed != null) metrics.push(['Start to finish',runTime(run.elapsed)]);
      let usageNote = '';
      if (!run.launches) usageNote = 'World has no recorded harness launch for this episode.';
      else if (run.turns == null && run.costUsd == null && !run.tokens) usageNote = 'World has no verified usage totals for this run.';
      else if (run.costUsd == null) usageNote = 'The retained run record supplies no USD amount.';
      html += section('Run details', rows(metrics) + (usageNote ? `<p class="muted">${usageNote}</p>` : ''));
      if (run.permissions) html += section('Permissions & budgets', permissionsHTML(run.permissions));
      if (run.checks?.length) {
        html += section('World finish checks · ' + run.checks.length, rows(run.checks.map(c => [c.name, c.result === 'PASS' ? 'Passed' : 'Failed'])));
      }
    }
    if (github) {
      const checks = github.checks || [];
      const passed = checks.filter(c => (c.conclusion || c.status || '').toUpperCase() === 'SUCCESS').length;
      const label = checks.length ? `GitHub checks · ${passed}/${checks.length} passed` : 'GitHub checks';
      const url = /^https:\/\/github\.com\/[\w.-]+\/[\w.-]+\/pull\/[0-9]+(?:\/checks)?$/.test(github.url) ? github.url.replace(/(?:\/checks)?$/, '/checks') : '';
      const detail = checks.length ? rows(checks.map(c => [c.name, cap(String(c.conclusion || c.status || 'Unknown').toLowerCase().replace(/_/g,' '))])) : '<p>No checks in this snapshot.</p>';
      html += section(label, detail + `<p>Captured ${esc(exactOf(github.capturedAt))}${url ? ` · <a href="${esc(url)}" target="_blank" rel="noopener noreferrer">Open on GitHub</a>` : ''}</p>`);
    }
    return html;
  }

  function humanDecisionsHTML(run) {
    if (!run?.decisions?.length) return '';
    return '<section class="node-decisions"><h3>Human decisions</h3><ol>' + run.decisions.map(d =>
      `<li><span>${esc(d.action)}</span><time datetime="${esc(d.at)}">${esc(exactOf(d.at))}</time></li>`).join('') + '</ol></section>';
  }

  const CONFIRM_ID = /^vr-[0-9a-f]{16}$/;
  const FINISH_NODE_ID = /^([a-z][a-z0-9-]{0,31}):(?:episode|commit):/;
  const finishEpisode = F => F && FINISH_NODE_ID.test(F.id) && (F.kind === 'open' || F.kind === 'step') ? (F.run?.episode || F.episode || '') : '';

  function finishMatchesNode(href, F, aliases) {
    const episode = finishEpisode(F);
    if (!episode || typeof href !== 'string' || !href.startsWith('/surface?')) return false;
    const place = parseAppQuery(href.slice('/surface'.length));
    const repo = FINISH_NODE_ID.exec(F.id)?.[1];
    return !!repo && !!place && place.brane === F.brane && place.brane.startsWith(repo + ':brane:') &&
      place.at === repo + ':episode:' + episode &&
      (place.at === F.id || aliases?.[place.at] === F.id);
  }

  function finishPlaceholder(F) {
    return finishEpisode(F) ? '<section class="node-decisions" data-node-confirmations hidden><h3>Awaiting approval</h3><div data-finish-links></div></section>' : '';
  }

  function nodeLayout(F, bar, head, overview, tab, activity) {
    const run = inspectionHTML(F);
    const sections = [
      ['overview', 'Overview', overview],
      ['files', 'Files', filesTabHTML(F, activity)],
    ];
    if (run) sections.push(['run', F.run ? 'Run' : 'Checks', `<div class="node-run-grid">${run}</div>`]);
    if (F.trajectoryEpisode || F.run?.episode) {
      const label = F.kind === 'shared' ? 'Linked episode trajectory' : 'Agent trajectory';
      sections.push(['trajectory', 'Trajectory', `<h2 class="node-panel-title">${label}</h2><div data-trajectory-body></div>`]);
    }
    if (!sections.some(([id]) => id === tab)) tab = 'overview';
    const buttons = sections.map(([id,label]) =>
      `<button type="button" role="tab" data-node-tab="${id}" id="node-tab-${id}" aria-controls="node-panel-${id}" aria-selected="${tab === id}" tabindex="${tab === id ? 0 : -1}">${label}</button>`).join('');
    const panels = sections.map(([id,,content]) =>
      `<section role="tabpanel" tabindex="0" data-node-panel="${id}" id="node-panel-${id}" aria-labelledby="node-tab-${id}"${tab === id ? '' : ' hidden'}>${content}</section>`).join('');
    return bar + `<nav class="node-tabs" role="tablist" aria-label="Node sections">${buttons}</nav><div class="fbody stack-body" tabindex="0" role="region" aria-label="Node summary">${head}${panels}</div>`;
  }

  function overviewHTML(inputs, history, outputs) {
    return `<div class="node-overview">${inputs}${history || outputs ? `<div class="node-outcome">${history}${outputs}</div>` : ''}</div>`;
  }

  // frontHTML is the front card for the stack's front state. T (the adapted model) only names
  // and glyphs other states; every value on the card is the stack's own, or an honest gap.
  function frontHTML(S, T, tab = 'overview', activity = null) {
    const F = S.front, step = F.kind === 'step';
    const none = t => `<div class="muted">${esc(t)}</div>`;
    if (F.kind === 'open') return openFrontHTML(S, T, none, tab, activity);
    const tag = F.episode ? F.episode + (F.model ? ' · ' + prettyModel(F.model) : '') : F.label;
    const bar = `<div class="fbar"><span id="cb-tag" style="font-weight:700">${esc(tag)}</span><span class="sub" id="cb-when">${esc(`${F.rel} · ${F.exact}`)}</span><span style="flex:1"></span><span class="sub" style="display:inline-flex;align-items:center;gap:7px">${gl(F.glyph, F.glyph === 'circle' ? SHARED_INK : F.colour, 11)}${esc(rowLabel(F.name, F.row))}</span><span class="sub" id="cb-pos">${esc(F.position)}</span></div>`;
    const summ = summaryOf(F);
    const title = cleanTitle(F.title) || F.label;
    const head = `<div class="head"><div class="h">${esc(title)}</div>${summ ? `<div class="outcome">${esc(summ)}</div>` : ''}${recordedHeadingHTML(F, title, summ)}</div>`;
    // Inputs: the agent, the task, then what it came from and read.
    let grid = '';
    if (F.model) grid += `<span class="k">Agent</span><div>${esc(agentLine(F))}</div>`;
    const gaps = !step || F.run ? [] : F.episode ? ['Grants: not recorded yet.'] : ['No episode record for this step: its agent, task, grants and reads are not recorded.'];
    const came = [];
    if (F.builtOn) came.push(stateChip(T, F.builtOn, 'built on'));
    (F.parents || []).forEach(p => came.push(stateChip(T, p, 'parent')));
    (F.reads || []).forEach(rd => {
      const g = rd.via === 'commit' || !rd.glyph ? gl(null, null, 10) : gl(rd.glyph, rd.colour, 10);
      came.push(chip(rd.fromId, rd.label, g, 'read', rd.via === 'commit', 'reads:' + rd.via + ':' + rd.fromId, ` data-read="${esc(rd.fromId)}"`));
    });
    const reads = F.readCalls != null ? ` · ${F.readCalls} read${F.readCalls === 1 ? '' : 's'} over ${F.readFiles} file${F.readFiles === 1 ? '' : 's'}` : '';
    const readGap = step && F.episode && !F.reads && !F.run ? none('Reads: not recorded for this step.') : '';
    const inputs = F.taskBody || F.task || grid || gaps.length || came.length || reads || readGap ? `<section class="col" aria-labelledby="front-inputs"><h2 class="colh" id="front-inputs">Inputs</h2>${taskHTML(F)}${grid ? `<div class="ingrid">${grid}</div>` : ''}${gaps.map(none).join('')}
      ${came.length || reads ? `<div class="sect"><div class="l">↑ Came from${reads}</div><div class="chips">${came.join('')}</div></div>` : ''}${readGap}</section>` : '';

    const pr = F.pr ? `<div class="prpill">PR #${F.pr}</div>` : '';
    const led = F.mergedInto ? `<div class="sect"><div class="l">↓ Led to</div><div class="chips">${stateChip(T, F.mergedInto, '')}</div></div>` : '';
    const outputBody = commitHTML(F) + pr + led;
    const outputs = outputBody ? `<section class="col" aria-labelledby="front-outputs"><h2 class="colh" id="front-outputs">Outputs</h2>${outputBody}</section>` : '';

    const acts = (F.decisions || []).map(a => {
      const [bg, fg, bd] = ACT[a.actor] || ACT.world;
      return `<li class="act"${a.source ? ` data-tip="src:${esc(a.source)}"` : ''}><span class="who" style="background:${bg};color:${fg};border-color:${bd}">${WHO[a.actor] || esc(a.actor)}</span><span class="tx">${esc(a.text)}<span class="tm">${esc(a.rel)}</span></span></li>`;
    }).join('');
    const ck = F.checks || {};
    const totals = [[ck.graders, 'graders'], [ck.rounds, ck.rounds === 1 ? 'round' : 'rounds'], [ck.minutes, 'agent-minutes'], [ck.calls, 'World calls']]
      .filter(([value]) => value != null && value !== '').map(([value, label]) => `${esc(value)} ${label}`).join(' · ');
    const decisions = humanDecisionsHTML(F.run);
    const finish = finishPlaceholder(F);
    const hasHistory = acts || totals || decisions;
    const happened = hasHistory || finish
      ? `<section class="col" aria-labelledby="front-history"${hasHistory ? '' : ' data-finish-history-only hidden'}><h2 class="colh" id="front-history">What happened</h2>${finish}${acts ? `<ol class="front-history">${acts}</ol>` : ''}${totals ? `<div class="check-totals">${totals}</div>` : ''}${decisions}</section>`
      : '';

    return nodeLayout(F, bar, head, overviewHTML(inputs, happened, outputs), tab, activity);
  }

  // openFrontHTML is an open episode's minimal card (row 220 M4): what it is, what it was built on,
  // its lifecycle label. Nothing has landed, so outputs and checks are honest gaps; N2d owns the richer card.
  function openFrontHTML(S, T, none, tab, activity) {
    const F = S.front, o = F.open || {};
    const bar = `<div class="fbar"><span id="cb-tag" style="font-weight:700">${esc(F.episode + (F.model ? ' · ' + prettyModel(F.model) : ''))}</span><span class="sub" id="cb-when">${esc(`started ${F.rel} · ${F.exact}`)}</span><span style="flex:1"></span><span class="sub" style="display:inline-flex;align-items:center;gap:7px">${gl(F.glyph, F.colour, 11)}${esc(rowLabel(F.name, F.row))}</span><span class="sub" id="cb-pos" data-open-state="${esc(o.state)}">${esc(o.label)}</span></div>`;
    const head = `<div class="head"><div class="h">${esc(F.title || F.episode)}</div><div class="outcome">${esc(o.label + (o.ended ? ' · ended ' + o.ended : ''))}</div></div>`;
    let grid = '';
    if (o.lastRun) grid += `<span class="k">Agent run</span><div>${esc(lastRunDetail(o.lastRun))}</div>`;
    if (F.model) grid += `<span class="k">Agent</span><div>${esc(agentLine(F))}</div>`;
    if (o.branch) grid += `<span class="k">Branch</span><div>${esc(o.branch)}</div>`;
    const came = F.builtOn ? `<div class="sect"><div class="l">↑ Came from</div><div class="chips">${stateChip(T, F.builtOn, 'built on')}</div></div>`
      : none(o.baseOutsideWindow ? `Built on ${o.base || 'an unrecorded base'}, outside the retained window.` : 'Base not recorded.');
    const inputs = `<div class="col"><div class="colh">Inputs</div>${taskHTML(F)}${grid ? `<div class="ingrid">${grid}</div>` : ''}${came}</div>`;
    const outputs = `<div class="col g6"><div class="colh">Outputs</div>${commitHTML(F)}${F.pr ? `<div class="prpill">#${F.pr} · waiting for merge</div>` : (F.commits?.length ? none('Nothing has landed yet.') : '')}</div>`;
    const decisions = humanDecisionsHTML(F.run);
    const finish = finishPlaceholder(F);
    const history = decisions || finish ? `<section class="col"${decisions ? '' : ' data-finish-history-only hidden'}><h2 class="colh">What happened</h2>${finish}${decisions}</section>` : '';
    return nodeLayout(F, bar, head, overviewHTML(inputs, history, outputs), tab, activity);
  }

  // railHTML is the whole brane, newest at the top; a long lane shows a window around the front.
  const RAIL_ROWS = 15;
  function railHTML(S) {
    const F = S.front, rail = S.rail || [], fi = Math.max(0, rail.findIndex(r => r.id === F.id));
    let lo = 0, hi = rail.length;
    if (rail.length > RAIL_ROWS) { lo = Math.max(0, Math.min(fi - Math.floor(RAIL_ROWS / 2), rail.length - RAIL_ROWS)); hi = lo + RAIL_ROWS; }
    const rows = rail.slice(lo, hi).map((r, j) => ({ r, k: lo + j })).reverse().map(({ r, k }) => {
      const handle = r.episode || r.label;
      return `<button type="button" class="rrow" data-k="${k}" data-tip="node:${esc(r.id)}" aria-current="${r.id === F.id}"><span class="lb"><span class="rt">${esc(cleanTitle(r.title) || r.label)}</span><span class="rid">${esc(handle)}${r.pr && '#' + r.pr !== handle ? ' · #' + r.pr : ''}</span></span>${r.badge > 0 ? `<span class="nb inl">${r.badge}</span>` : ''}<span class="tm">${esc(r.rel)}</span></button>`;
    }).join('');
    const foot = hi - lo < rail.length ? `<div class="rail-foot">${hi - lo} of ${rail.length} states on this brane · ↑ ↓ to move</div>` : '';
    return `<div class="rail-h">${gl(F.glyph, F.glyph === 'circle' ? SHARED_INK : F.colour, 11)}${esc(rowLabel(F.name, F.row))}</div>` + rows + foot;
  }

  // ---------- Places: the URL follows the view; the address bar is M4a's grammar ----------
  const appPath = (brane, at) => brane == null ? '/surface' : '/surface?brane=' + encodeURIComponent(brane) + '&at=' + encodeURIComponent(at);
  const stackPath = (brane, at) => '/surface/stack.json?brane=' + encodeURIComponent(brane) + '&at=' + encodeURIComponent(at);
  const filePath = (brane, at, path, source = 'committed') => '/surface/file.json?brane=' + encodeURIComponent(brane) + '&at=' + encodeURIComponent(at) + '&path=' + encodeURIComponent(path) +
    (source === 'committed' ? '' : '&source=' + encodeURIComponent(source));

  function patchLines(source) {
    const lines = source.split('\n');
    if (lines.at(-1) === '') lines.pop();
    let old = 0, next = 0, oldLeft = 0, newLeft = 0;
    return lines.map(text => {
      const hunk = /^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@/.exec(text);
      if (hunk) {
        old = Number(hunk[1]); oldLeft = Number(hunk[2] ?? 1);
        next = Number(hunk[3]); newLeft = Number(hunk[4] ?? 1);
        return { text, type: 'hunk' };
      }
      const marker = text[0];
      if (oldLeft > 0 && marker === '-') {
        oldLeft--;
        return { text, type: 'removed', before: old++ };
      }
      if (newLeft > 0 && marker === '+') {
        newLeft--;
        return { text, type: 'added', after: next++ };
      }
      if (oldLeft > 0 && newLeft > 0 && marker === ' ') {
        oldLeft--; newLeft--;
        return { text, type: 'context', before: old++, after: next++ };
      }
      return { text, type: 'meta' };
    });
  }

  function fileSizeLabel(bytes) {
    if (!Number.isSafeInteger(bytes) || bytes < 0) return '';
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KiB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MiB';
  }

  function filePreviewHTML(evidence) {
    const patch = typeof evidence.diff === 'string' ? patchLines(evidence.diff) : [];
    const first = patch.findIndex(line => line.type === 'added' || line.type === 'removed');
    const status = evidence.sourceLabel ? '' : FILE_STATUS.get(evidence.status);
    const size = fileSizeLabel(evidence.status === 'D' ? evidence.size?.before : evidence.size?.after);
    const meta = [status, size].filter(Boolean).join(' · ');
    const sample = first >= 0 ? patch.slice(Math.max(0, first - 1), first + 6)
      : evidence.text && typeof evidence.text === 'string' ? evidence.text.split('\n', 6).map(text => ({text, type: 'context'})) : [];
    const lines = sample.map(({text, type}) => `<span class="patch-${type}">${esc(text) || ' '}</span>`).join('\n');
    return `<div class="file-preview-label">${esc(evidence.sourceLabel || 'Committed file preview')}${meta ? ' · ' + meta : ''}</div>
      ${lines ? `<pre>${lines}</pre>` : evidence.text === '' ? '<div>Empty file.</div>' : evidence.kind === 'unavailable' ? '<div>Retained file contents unavailable.</div>' : '<div>Open for file details.</div>'}`;
  }

  function fileEvidenceHTML(evidence, mode) {
    const messages = {
      binary: 'Binary content has no text preview.',
      symlink: 'This path is a symbolic link.',
      submodule: 'This path is a Git submodule.',
      'too-large': 'This file exceeds the preview limit.',
      unavailable: 'Retained file contents unavailable.',
    };
    const hasText = evidence.kind === 'text' && typeof evidence.text === 'string';
    const hasDiff = typeof evidence.diff === 'string' && evidence.diff.length > 0;
    const diff = hasDiff && (mode !== 'file' || !hasText);
    const controls = hasText && hasDiff ? `<div class="file-modes" role="group" aria-label="File view">
      <button type="button" data-file-mode="diff" aria-pressed="${diff}">Diff</button>
      <button type="button" data-file-mode="file" aria-pressed="${!diff}">File</button>
      </div>` : '';
    const rev = evidence.commit || evidence.merge;
    const origin = evidence.sourceLabel || (evidence.merge ? 'Merge' : '');
    const commit = origin || rev ? `<span>${esc(origin)}${rev ? ` <code>${esc(rev)}</code>` : ''}</span>` : '';
    const before = fileSizeLabel(evidence.size?.before), after = fileSizeLabel(evidence.size?.after);
    const size = before && after ? `${before} → ${after}` : before || after;
    const patch = diff ? patchLines(evidence.diff) : [];
    const visiblePatch = patch.slice(0, 5000);
    const clippedPatch = evidence.truncated || visiblePatch.length < patch.length;
    const added = visiblePatch.filter(line => line.type === 'added').length;
    const removed = visiblePatch.filter(line => line.type === 'removed').length;
    const changes = added || removed ? `<span class="patch-count"><span class="patch-added">+${added}</span> <span class="patch-removed">−${removed}</span>${clippedPatch ? ' shown' : ''}</span>` : '';
    const facts = commit || size || changes ? `<div class="file-evidence-facts muted">${commit}${size ? `<span>${size}</span>` : ''}${changes}</div>` : '';
    const content = diff ? evidence.diff : hasText ? evidence.text : null;
    const hint = diff ? (clippedPatch ? 'Only part of the diff is shown.' : '') :
      evidence.status === 'D' && hasText ? 'Deleted in this merge. Showing the first-parent version.' : hasText && !evidence.sourceLabel ? 'Version at the merge.' : '';
    const details = hint ? `<div class="muted">${hint}</div>` : '';
    const preview = content === null ? `<div class="muted">${Object.hasOwn(messages, evidence.kind) ? messages[evidence.kind] : 'File preview unavailable.'}</div>`
      : content === '' ? '<div class="muted">Empty file.</div>'
      : diff ? `<pre class="file-evidence-code rich-patch" aria-label="Git diff"><code>${visiblePatch.map(line =>
        `<span class="patch-row patch-${line.type}" data-before="${line.before ?? ''}" data-after="${line.after ?? ''}">${esc(line.text) || ' '}</span>`
      ).join('\n')}</code></pre>`
      : `<pre class="file-evidence-code" aria-label="File contents"><code>${esc(content)}</code></pre>`;
    return facts + controls + details + preview;
  }

  // parseAppQuery is /surface's query grammar: none ({}), or exactly one brane and one at; else null.
  function parseAppQuery(search) {
    const p = new URLSearchParams(String(search || '').replace(/^\?/, ''));
    const keys = [...p.keys()];
    if (keys.length === 0) return {};
    if (keys.length !== 2 || p.getAll('brane').length !== 1 || p.getAll('at').length !== 1) return null;
    return { brane: p.get('brane'), at: p.get('at') };
  }
  // addressTarget resolves what was typed after "world>" through parse (surface.js's parseAddress).
  function addressTarget(raw, index, parse) {
    const t = String(raw == null ? '' : raw).trim();
    if (t === '') return { home: true };
    if (typeof parse !== 'function') return { error: 'the address bar is not loaded (/surface/surface.js)' };
    const r = parse(/^(world|back)(\s|$)/i.test(t) ? t : 'world ' + t, index) || {};
    if (!r.href) return { error: r.error || 'unknown address' };
    const i = r.href.indexOf('?'), q = parseAppQuery(i < 0 ? '' : r.href.slice(i));
    if (!q) return { error: 'unknown address' };
    return q.brane == null ? { home: true } : { open: q };
  }
  // lookOf is v3's two choices from M4a's stored look: Paper or Slate, full or reduced motion.
  function lookOf(look, prefersDark, prefersReduced) {
    const l = look && typeof look === 'object' ? look : {};
    const theme = l.theme === 'system' ? (prefersDark ? 'slate' : 'paper') : l.theme === 'slate' || l.theme === 'night' ? 'slate' : 'paper';
    const motion = l.motion === 'reduced' || (l.motion === 'system' && prefersReduced) ? 'reduced' : 'full';
    return { theme, motion };
  }

  // crumbsHTML is the header's crumbs: none on the timeline (the brand is the way home), and in a
  // brane only its address (row-189@6), current.
  const crumbsHTML = addr => addr == null ? '' : `<button type="button" class="crumb" aria-current="true">${esc(addr)}</button>`;

  function splitFrames(tail, chunk) {
    const frames = (tail + chunk).split("\n\n");
    const rest = frames.pop();
    const events = frames.flatMap(frame => frame.split("\n")
      .filter(line => line.startsWith("data: ")).map(line => JSON.parse(line.slice(6))));
    return {events, tail: rest};
  }
  function applyDelta(state, delta) {
    for (const op of delta) {
      if (op.op !== "replace") throw new Error("unsupported state delta");
      if (op.path === "/lastIndex") state.lastIndex = op.value;
      else if (op.path === "/logHead") state.logHead = op.value;
      else throw new Error("unknown state path");
    }
    return state;
  }
  function nextBackoff(retryAfter, previous = 0) {
    const parsed = Number(retryAfter);
    const first = Number.isFinite(parsed) && parsed > 0 ? parsed : 1;
    return Math.min(30, previous > 0 ? previous * 2 : first);
  }
  function liveRequest(cursor, run) {
    return {threadId: 'surface', runId: 'surface-' + run, messages: [], state: {schema: 'world/agui-state/v1', lastIndex: cursor}};
  }
  function liveStatus(s) {
    if (s.mode === 'unavailable') return 'Live updates unavailable: ' + s.reason;
    if (s.mode === 'stopped') return 'Live stopped: ' + s.reason;
    if (s.mode === 'hidden') return 'Paused while the tab is hidden';
    if (s.mode === 'paused') return 'Paused: ' + s.reason + ' · retrying in ' + s.wait + ' s';
    if (s.count) return 'Live · ' + s.count + ' new entr' + (s.count === 1 ? 'y' : 'ies') + ' · refreshing';
    return 'Live · up to entry ' + s.cursor + ' · checked ' + s.checked;
  }

  // ---------- Follow mode (row 220 M4): sources, running, polling ----------
  // sourcesHTML is each repo's freshness, as the daemon worded it ("local branch checked": remote freshness is unknown).
  function sourcesHTML(repos) {
    if (!Array.isArray(repos) || !repos.length) return '';
    return `<div class="nsec" data-sources><div class="nh">Sources</div>${repos.map(r => `<div class="nl sub" data-repo="${esc(r.id)}" data-state="${esc(r.state)}"${r.warn ? ' style="color:var(--you)"' : ''}>${esc(r.line)}</div>`).join('')}</div>`;
  }
  // Show observed launch history beside the credential count. Only a complete read can confirm zero live sessions.
  function openRunningHTML(o) {
    const live = o.live || [], known = o.status === 'complete' || o.status === 'incomplete';
    const rows = live.slice(0, 4).map(r => r.id
      ? `<button type="button" class="nl run" data-go="${esc(r.id)}"><b>${esc(r.episode)}</b>${r.repo ? ' · ' + esc(r.repo) : ''} · ${esc(r.label)}</button>`
      : `<div class="nl sub"><b>${esc(r.episode)}</b> · ${esc(r.reason)}</div>`).join('');
    const counts = [o.prOpen ? `${o.prOpen} PR open` : '', o.ended ? `${o.ended} sessions ended` : '', o.failed ? `${o.failed} finish not landed` : ''].filter(Boolean).join(' · ');
    const head = known ? `Live sessions · ${live.length}${o.status === 'incomplete' ? '+' : ''}` : 'Live sessions';
    const run = o.lastRun;
    const latest = run ? `<button type="button" class="nl run" data-go="${esc(run.id)}"><span><b>${esc(run.episode)}</b> · ${esc(lastRunDetail(run))}</span></button>` : '';
    return `<div class="nsec" data-running><div class="nh">${head}</div>${latest}${rows}${live.length > 4 ? `<div class="nl sub">+${live.length - 4} more</div>` : ''}` +
      `${!live.length && o.status === 'complete' ? '<div class="nl sub">No live sessions</div>' : ''}${counts ? `<div class="nl sub">${esc(counts)}</div>` : ''}${o.note ? `<div class="nl sub">${esc(o.note)}</div>` : ''}</div>`;
  }
  // fetchCoherent reads the model and, with a brane open, its stack from ONE generation. "same": the generation on
  // screen, nothing to re-render; "apply": a coherent pair (stack null with missing when the selection is gone);
  // "keep": the generations disagreed twice, keep the previous view; "fail": the model answered an error.
  async function fetchCoherent(get, selection, shown) {
    let missing = null;
    for (let attempt = 0; attempt < 2; attempt++) {
      const model = await get('/surface/model.json');
      if (!model.ok || !model.body || !Array.isArray(model.body.lanes)) return { kind: 'fail', status: model.status, body: model.body };
      const gen = model.body.generation;
      if (gen != null && gen === shown) return { kind: 'same', model: model.body };
      if (!selection) return { kind: 'apply', model: model.body, stack: null };
      const stack = await get(stackPath(selection.brane, selection.at));
      if (stack.ok && stack.body && stack.body.front) {
        if (stack.body.generation === gen) return { kind: 'apply', model: model.body, stack: stack.body };
        missing = null;
      } else {
        if (gen == null) return { kind: 'fail', status: stack.status, body: stack.body };
        missing = (stack.body && stack.body.reason) || 'the selection could not be read (HTTP ' + stack.status + ')';
      }
    }
    return missing ? { kind: 'apply', missing } : { kind: 'keep' };
  }
  function serviceLineHTML(service) {
    if (!service) return '';
    return '<div class="nl sub" data-service' + (service.warn ? ' style="color:var(--you)"' : '') + '>' + esc(service.line) + '</div>';
  }

  // needsYouHTML is the now card's Needs you (row 201 A3, design §2.8): the pending mint requests of
  // GET /surface/requests.json, each a link to its card; a request the daemon could not read shows why;
  // then the notification opt-in for this browser's permission (design §2.9).
  // Row 201 A4 (design §2.8): then the pending Human.Approve decisions of the same read, each a link to
  // its decision card; a reserved one names its typed owner and a refused one its reason, with no link.
  // Row 243 V3 (design §2.5): then the pending finish confirmations, each a link to its card; a refused one
  // shows its reason with no link; confirmationsUnavailable is said, never hidden.
  function needsYouHTML(r, now, permission, fromHeads = false) {
    const p = r && Array.isArray(r.pending) ? r.pending : [];
    const ds = r && Array.isArray(r.decisions) ? r.decisions : [];
    const cs = r && Array.isArray(r.confirmations) ? r.confirmations : [];
    const n = p.length + ds.length + cs.length;
    const head = `<div class="nh" style="color:${n ? 'var(--youShell)' : 'var(--sub)'}">Needs you · ${n}</div>`;
    const more = (r && r.confirmationsUnavailable ? `<div class="nl sub">Confirm · unavailable: ${esc(r.confirmationsUnavailable)}</div>` : '')
      + (r && r.decisionsTruncated ? '<div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div>' : '');
    if (!n) return `<div class="nsec" data-needs>${head}<div class="nl sub">${esc(r && r.unavailable || 'Nothing needs you')}</div>${more}${notifyHTML(permission)}</div>`;
    return `<div class="nsec" data-needs>${head}` + p.map(q => q.refused ? `<div class="nl sub">${esc(q.id)} · refused: ${esc(q.refused)}</div>`
      : `<div class="nl"><a id="needs-${esc(q.id)}" href="/surface/approve/${esc(q.id)}"><b>${esc(q.episode)}</b> · ${esc(q.agent)} · ${q.grants} grant${q.grants === 1 ? '' : 's'} · expires in ${Math.max(0, Math.ceil((q.expires - now) / 60))} min</a>`
        + (q.recorded ? ` · approved by ${esc(q.recorded)}, waiting for the minter` : '') + '</div>').join('') + ds.map(decisionLineHTML).join('') + cs.map(c => confirmLineHTML(c, fromHeads)).join('')
      + more + notifyHTML(permission) + '</div>';
  }
  // RESERVED_OWNERS is the typed command that owns a reserved decision (design §2.6 L1).
  const RESERVED_OWNERS = { publish: 'world-publish approve', finish: 'ailang-worldd episode finish' };
  // decisionLineHTML is one decision's Needs-you line (design §2.8).
  function decisionLineHTML(x) {
    if (x.reserved) return `<div class="nl sub">${esc(x.reserved)} approval — use ${esc(RESERVED_OWNERS[x.reserved] || 'its typed owner')} at a terminal</div>`;
    if (x.refusal) return `<div class="nl sub">Decide · ${esc(x.scope)} · refused: ${esc(x.refusal)}</div>`;
    const ref = String(x.request || '').replace(/^sha256:/, '');
    return `<div class="nl"><a id="decide-${esc(ref)}" href="/surface/decide/sha256:${esc(ref)}">Decide · ${esc(x.scope)} · requested by ${esc(x.requester)} · deadline ${hm(new Date(x.deadline * 1000))}</a></div>`;
  }
  // confirmLineHTML is one finish confirmation's Needs-you line (row 243 V3 §2.5): a link to its card only for
  // a summarised row whose id is a confirmation id (store.VerbRequestIDPattern); otherwise its reason, no link.
  function confirmLineHTML(c, fromHeads) {
    if (c.refused || !CONFIRM_ID.test(c.id)) return `<div class="nl sub">Confirm · ${esc(c.id)} · refused: ${esc(c.refused || 'not a confirmation id')}</div>`;
    return `<div class="nl"><a id="confirm-${esc(c.id)}" href="/surface/confirm/${esc(c.id)}#${fromHeads ? 'heads/activity' : 'activity'}">Confirm · episode finish ${esc(c.episode)} · ${esc(c.repository)} ${esc(c.head)}@${esc(String(c.sha || '').slice(0, 12))}`
      + ` · ${c.draft ? 'draft' : 'ready'} · expires ${hm(new Date(c.expires * 1000))}</a></div>`;
  }
  // notifyHTML is the opt-in: "Notify me" while this browser has not been asked (permission "default"),
  // its answer once it has, and nothing where the browser has no notifications.
  function notifyHTML(permission) {
    if (permission === 'default') return '<button type="button" class="nl run" id="notify-me">Notify me</button>';
    if (permission === 'granted') return '<div class="nl sub" data-notify>Notifications on</div>';
    return permission === 'denied' ? '<div class="nl sub" data-notify>Notifications are blocked for this page</div>' : '';
  }
  // notifications is what one read of GET /surface/requests.json raises (design §2.9): a "World needs you"
  // for each pending request this tab has not seen, while notifications are granted; each opens its card.
  // seen is null before the page's first read, which raises nothing (the page itself shows those); a read
  // that failed changes nothing. Refused and already-recorded requests are seen but raise nothing.
  // A decision (row 201 A4) is keyed by its request ref; a reserved or refused one is seen but raises nothing.
  function notifications(r, seen, permission) {
    if (!r || !Array.isArray(r.pending)) return { notes: [], seen };
    const ds = Array.isArray(r.decisions) ? r.decisions : [];
    const fresh = seen ? r.pending.filter(q => !seen.has(q.id) && !q.refused && !q.recorded) : [];
    const freshD = seen ? ds.filter(x => !seen.has(x.request) && !x.reserved && !x.refusal) : [];
    return { seen: new Set([...(seen || []), ...r.pending.map(q => q.id), ...ds.map(x => x.request)]), notes: permission !== 'granted' ? [] : fresh.map(q => ({ title: 'World needs you', tag: q.id,
      body: `Approve mint request for ${q.episode} (${q.grants} grant${q.grants === 1 ? '' : 's'}), expires ${hm(new Date(q.expires * 1000))}`, href: '/surface/approve/' + q.id }))
      .concat(freshD.map(x => ({ title: 'World needs you', tag: x.request,
        body: `Decide ${x.scope} requested by ${x.requester}, deadline ${hm(new Date(x.deadline * 1000))}`, href: '/surface/decide/' + x.request }))) };
  }
  // decidedHTML is the decision list under Needs you (design §2.8): the last few of GET
  // /surface/requests.json's decided requests (newest first; the daemon sends 20), and all of them on request.
  const DECIDED_FEW = 3;
  function decidedHTML(r, all) {
    const d = r && Array.isArray(r.decided) ? r.decided : [];
    if (!d.length) return '';
    return '<div class="nsec" data-decided><div class="nh">Decided</div>' + (all ? d : d.slice(0, DECIDED_FEW)).map(x =>
      `<div class="nl sub">You · Minted · ${esc(x.episode)} · ${esc(x.approvedBy)} · ${hm(new Date(x.at * 1000))}</div>`).join('')
      + (d.length > DECIDED_FEW ? `<button type="button" class="nl run" id="decided-all" aria-expanded="${!!all}">${all ? 'Show fewer' : 'Show all ' + d.length}</button>` : '') + '</div>';
  }

  // ---------- The bar's verbs (row 243 V1): GET /surface/verb.json decides; only a schema-valid address navigates ----------
  const VERB_ROUTE_BUDGET_MS = 15000;   // above the daemon's 10 s read deadline; a named constant, not a setting
  const VERB_NEEDS = { address: [], result: ['cli'], navigate: ['href'], terminal: ['cli', 'reason'], refused: ['reason'] };
  const UNREADABLE = { say: 'the verb route answered something this page cannot read' };
  // verbOutcome maps the route's answer (HTTP status, body text) to one action; anything else is a refusal.
  function verbOutcome(status, text) {
    if (status !== 200) return { say: 'the verb route did not answer: HTTP ' + status };
    let b = null;
    try { b = JSON.parse(text); } catch (_) { return UNREADABLE; }
    const need = b && typeof b === 'object' && b.schema === 'world/surface-verb/v1' && b.v === 1 && Object.prototype.hasOwnProperty.call(VERB_NEEDS, b.outcome) ? VERB_NEEDS[b.outcome] : null;
    if (!need || !need.every(k => typeof b[k] === 'string' && b[k] !== '')) return UNREADABLE;
    if (b.outcome === 'address') return { address: true };
    if (b.outcome === 'refused') return typeof b.cli === 'string' && b.cli !== '' ? { panel: { head: '$ ' + b.cli, body: b.reason, copy: b.cli } } : { say: b.reason };
    if (b.outcome === 'navigate') return b.href.startsWith('/surface/approve/mr-') || b.href.startsWith('/surface/decide/sha256:') || b.href.startsWith('/surface/confirm/vr-') ? { assign: b.href } : { say: 'refused: unexpected target' };
    if (b.outcome === 'terminal') return { panel: { head: '$ ' + b.cli, body: b.reason + '\nAt a terminal: ' + b.cli, copy: b.cli } };
    const shown = [b.text, b.json, b.stderr].filter(x => typeof x === 'string' && x !== '');
    if (shown.length !== 1 || (b.exit !== undefined && b.exit !== 1 && b.exit !== 3) || (b.stderr !== undefined) !== (b.exit === 1)) return UNREADABLE;
    return { panel: { head: '$ ' + b.cli, body: shown[0] + (b.exit ? 'exit ' + b.exit : ''), copy: '' } };
  }
  // askVerb asks the route under one budget, across the fetch and the body (io carries fetch, the
  // clock and AbortController). A newer ask supersedes an older one: the older is aborted and its
  // answer dropped, never acted on; the newer answers for the bar.
  let verbJob = null;
  function dropVerb() { if (verbJob) { verbJob.superseded = true; verbJob.ctl.abort(); } }
  async function askVerb(raw, io) {
    dropVerb();
    const job = { ctl: new io.AbortController(), superseded: false, late: false };
    verbJob = job;
    const stop = new Promise((_, reject) => job.ctl.signal.addEventListener('abort', () => reject(new Error('aborted'))));
    stop.catch(() => {});
    const timer = io.setTimeout(() => { job.late = true; job.ctl.abort(); }, VERB_ROUTE_BUDGET_MS);
    try {
      const res = await Promise.race([io.fetch('/surface/verb.json?q=' + encodeURIComponent(raw), { signal: job.ctl.signal, headers: { Accept: 'application/json' } }), stop]);
      const text = await Promise.race([res.text(), stop]);
      return job.superseded ? { superseded: true } : verbOutcome(res.status, text);   // an answer that beat the abort is still dropped
    } catch (e) {
      if (job.superseded) return { superseded: true };
      if (job.late) return { say: 'the verb route did not answer in time' };
      return { say: 'the verb route is unreachable: ' + (e && e.message ? e.message : String(e)) };
    } finally {
      io.clearTimeout(timer);
      if (verbJob === job) verbJob = null;
    }
  }
  // verbAct runs one action through fx: only an address reaches the address path; only an address,
  // a navigation or a panel clears the bar; a refusal says why and leaves the view and the input.
  function verbAct(act, raw, fx) {
    if (act.superseded) return;
    if (act.address) { fx.clear(); return fx.address(raw); }
    if (act.assign) { fx.clear(); return fx.assign(act.assign); }
    if (act.panel) { fx.clear(); return fx.panel(act.panel); }
    fx.say(act.say);
  }
  // goWith is the bar's Enter over fx (design A4): every Enter owns the bar. An empty line supersedes a pending ask and goes
  // home; an answer that arrives after a newer Enter is dropped, and the address path re-checks after its own await.
  let barTurn = 0;
  async function goWith(raw, fx) {
    const turn = ++barTurn, current = () => turn === barTurn;
    if (String(raw == null ? '' : raw).trim() === '') { dropVerb(); return fx.out(); }
    const act = await fx.ask(raw);
    if (!current()) return;
    return verbAct(act, raw, { ...fx.act, address: r => fx.address(r, current) });
  }
  // addressGo is the address bar's own path (M4a's grammar), reached only through the route's address answer; after the
  // index loads it acts only for the current Enter.
  async function addressGo(raw, current, fx) {
    let r = addressTarget(raw, null, fx.parse);
    if (r.error && typeof fx.parse === 'function') {
      try { r = addressTarget(raw, await fx.loadIndex(), fx.parse); } catch (e) { r = { error: 'the surface index could not be read: ' + e.message }; }
      if (!current()) return;
    }
    if (r.home) return fx.out();
    if (r.open) return fx.openState(r.open.brane, r.open.at);
    fx.say(r.error);
  }
  // verbPanel builds the panel's nodes with textContent only: a report or a reason is never markup.
  function verbPanel(doc, p) {
    const el = doc.createElement('div'), head = doc.createElement('div'), body = doc.createElement('pre');
    el.id = 'verb'; el.className = 'pop'; el.setAttribute('role', 'dialog'); el.setAttribute('aria-label', 'Verb');
    head.className = 'mono'; head.textContent = p.head;
    body.style.cssText = 'margin:10px 0 0;white-space:pre-wrap;max-height:60vh;overflow:auto;font-size:12px'; body.textContent = p.body;
    el.append(head, body);
    if (p.copy) { const b = doc.createElement('button'); b.type = 'button'; b.className = 'hbtn'; b.textContent = 'Copy the line'; b.dataset.copy = p.copy; el.append(b); }
    return el;
  }

  if (typeof module !== 'undefined') module.exports = { needsYouHTML, decidedHTML, notifications, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, finishMatchesNode, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
    adaptModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, filePath, fileEvidenceHTML, filesTabHTML, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW },
    VERB_ROUTE_BUDGET_MS, verbOutcome, askVerb, verbAct, verbPanel, goWith, addressGo };
  if (typeof document === 'undefined') return;

  // ================= The page =================
  const $ = id => document.getElementById(id);
  const SEAMS = window.WorldSurface || {};
  let T = null, D = null, LANES = {}, N = {}, parents = {}, kids = {}, SH = [], NS = 0, TT = [], LX = [], MAXU = LEVELS.length - 1;
  let STEPS = [], TNOW = 0, NOW_D = new Date(), quiet = [], LY = {}, ROWS = [], MXS = [], XS = [];
  let nameOf, laneTitle, glyphOf, glyphOfNode, lastStep, lastState;
  let recordedModel;
  const arrivals = new Map();
  const arrivalButton = document.querySelector('#episode-arrivals button');
  function updateArrivals(model) {
    if (T && recordedModel && model.generation && model.generation !== recordedModel.generation) {
      const previous = new Set(T.OPEN.map(n => n.id));
      const open = model.lanes.flatMap(lane => lane.states).filter(n => n.kind === 'open');
      const current = new Set(open.map(n => n.id));
      for (const id of arrivals.keys()) if (!current.has(id)) arrivals.delete(id);
      for (const node of open) {
        if (!previous.has(node.id) && node.open?.sessionLive) arrivals.set(node.id, node.at);
      }
    }
    arrivalButton.hidden = !arrivals.size;
    arrivalButton.textContent = arrivals.size === 1 ? '1 new episode · Show' : arrivals.size + ' new episodes · Show latest';
  }
  arrivalButton.addEventListener('click', () => {
    const latest = [...arrivals].sort((a, b) => b[1].localeCompare(a[1]))[0]?.[0];
    if (!latest || !N[latest]) return;
    arrivals.clear();
    arrivalButton.hidden = true;
    if (camera.active) {
      camera.followLane(N[latest].lane);
    } else {
      pinnedTimelineState = latest;
      relayout();
      chooseTimelineNode(latest);
      goNow(false);
      const button = document.querySelector('#tl .node[data-state="' + CSS.escape(latest) + '"]');
      button?.focus({preventScroll:true});
      if (narrow.matches && button) button.scrollIntoView({block:'nearest', inline:'nearest'});
    }
  });
  const activityStyle = document.createElement('style');
  document.head.append(activityStyle);
  let activityTimer;
  function showRecentCalls(model) {
    clearTimeout(activityTimer);
    // Translate server time to elapsed browser time so device clock skew cannot keep a pulse alive.
    const serverNow = Date.parse(model.now), received = performance.now();
    const signals = Object.entries(model.nowEdge?.activityUntil || {});
    function paint() {
      const now = serverNow + performance.now() - received;
      const active = signals.filter(([, until]) => until > now);
      const selectors = active.flatMap(([id]) => {
        const value = CSS.escape(id);
        return [`.node[data-state="${value}"]`, `.node.bub-n:has(>[data-state="${value}"])`,
          `.hc-face[data-head-state="${value}"]`, `.hc-rear[data-rear-state="${value}"]`];
      });
      const rule = selectors.length ? selectors.join(',') + '{--world-call-ring:block;--world-mini-pulse:world-call-mini}' : '';
      if (activityStyle.textContent !== rule) activityStyle.textContent = rule;
      if (active.length) activityTimer = setTimeout(paint, Math.max(1, Math.min(...active.map(([, until]) => until)) - now));
    }
    paint();
  }
  const drafts = WorldHeadDrafts.create(({id, parent, discarded}) => {
    if (!recordedModel) return;
    setModel(recordedModel);
    st.focus = discarded ? parent : id;
    renderTimeline(true);
    applyView();
    camera.present();
    chooseAngle(true, {instant:true, keyboard:true, anchor:st.focus});
  });
  function setModel(m) {
    updateArrivals(m);
    showRecentCalls(m);
    recordedModel = m;
    T = adaptModel(drafts.project(m));
    ({ D, LANES, N, parents, kids, SH, NS, TT, LX, MAXU, STEPS, TNOW, NOW_D, quiet, LY, ROWS, nameOf, laneTitle, glyphOf, glyphOfNode, lastStep, lastState } = T);
    measureStops();
    computeXS(view.u);
  }
  function measureMini(width) {
    const a = LX[3][0], z = LX[3][NS];
    MXS = LX[3].map(x => 8 + (x - a) / ((z - a) || 1) * Math.max(1, width - 16));
  }
  // The open stack (world/surface-stack/v1), when a brane is open.
  let S = null, F = null, RAIL = [], READS = new Map();
  function setStack(s) { S = s; F = s.front; RAIL = s.rail || []; READS = new Map((F.reads || []).map(r => [r.via + ':' + r.fromId, r])); }
  const railIdx = () => F ? Math.max(0, RAIL.findIndex(r => r.id === F.id)) : 0;
  const setView = v => {
    $('stage').dataset.view = v;
    if (v === 'brane' || !T) { showHeadRibbon(false); placeNow(false); }
  };

  const isBreak = i => breakAt(TT, i);
  const breaksBetween = (a, b) => {
    for (let i = Math.min(a, b) + 1; i <= Math.max(a, b); i++) if (isBreak(i)) return true;
    return false;
  };
  // Enumerate ticks within drawn time spans to skip over unbounded quiet gaps.
  function* tickTimes(step) {
    if (!T.dated) return;
    const first = T.undated ? 1 : 0;
    for (let i = first; i <= NS; i++) {
      if (i > first && !isBreak(i)) {
        for (let t = (Math.floor(TT[i - 1] / step) + 1) * step; t < TT[i]; t += step) {
          yield [t, i - 1 + (t - TT[i - 1]) / (TT[i] - TT[i - 1])];
        }
      }
      if (TT[i] % step === 0) yield [TT[i], i];
    }
  }
  const computeXS = u => { XS = blendXS(LX, u); };
  const posF = f => { f = Math.max(0, Math.min(NS, f)); const i = Math.floor(f), j = Math.min(NS, i + 1); return XS[i] + (XS[j] - XS[i]) * (f - i); };
  const fAt = x => {
    if (x <= XS[0]) return 0; if (x >= XS[NS]) return NS;
    let lo = 0, hi = NS; while (hi - lo > 1) { const m = (lo + hi) >> 1; if (XS[m] <= x) lo = m; else hi = m; }
    return lo + (XS[hi] === XS[lo] ? 0 : (x - XS[lo]) / (XS[hi] - XS[lo]));
  };
  const level = () => view.u < 1.5 ? 'steps' : view.u < 2.5 ? 'hours' : view.u < 3.5 ? 'days' : 'weeks';
  const levelName = () => LEVELS[Math.round(view.u)].name;

  // ---------- Zoom/fit ----------
  const narrow = window.matchMedia('(max-width: 700px)');
  let scale = 1;
  function fit() {
    const stg = $('stage'), w = stg.clientWidth, h = stg.clientHeight;
    const right = narrow.matches ? w : 1262;
    const left = narrow.matches ? Math.min(128, Math.round(w * 0.32)) : 262;
    if (right !== PR || left !== PL) {
      const pinned = T && atNow();
      PL = left; PR = right; PW = Math.max(1, PR - PL);
      if (T) {
        renderTimeline(true);
        if (pinned) goNow(false); else applyView();
      }
    }
    scale = narrow.matches ? 1 : Math.min(w / 1600, h / 936, 1.25);
    const bx = Math.max(0, (w - 1600 * scale) / 2);
    $('tl').style.transform = `scale(${scale})`;
    $('br').style.transform = `translate(${bx}px, 0) scale(${scale})`;
    const edge = bx + 1290 * scale;
    $('L1').style.background = `linear-gradient(90deg, var(--stageBg) ${edge}px, var(--line) ${edge}px, var(--line) ${edge + 1}px, var(--bg) ${edge + 1}px)`;
  }

  // ---------- State ----------
  const st = { d: 0, focus: null };
  const view = { u: 1, L: 0 };
  const reduced = () => $('root').classList.contains('reduced');

  let POS = {}, REP = {}, DISP = {}, LABELED = new Set();
  let pinnedTimelineState = null;
  function layout() { ({ POS, REP, DISP, LABELED } = layoutGraph(T, XS, pinnedTimelineState)); }

  let fromHeads = false, timelineAnchor = '', stackFlight;
  function chooseTimelineNode(id) {
    timelineAnchor = N[id]?.id || DISP[id]?.ids?.at(-1) || '';
    if (st.focus === id) return;
    st.focus = id;
    applyFocus();
  }
  function sidePoints() {
    const points = new Map(), stageBox = $('stage').getBoundingClientRect();
    const world = $('world')?.getBoundingClientRect();
    if (world) {
      for (const [id, rep] of Object.entries(REP)) {
        const d = DISP[rep];
        points.set(id, {x:world.left - stageBox.left + d.x * scale, y:world.top - stageBox.top - toolbar.offsetHeight + d.y * scale});
      }
    }
    return points;
  }

  // Only attribute a pending request to a lane when its episode has one owner.
  function attentionLanes() {
    const requests = [...(REQS?.pending || []), ...(REQS?.decisions || []), ...(REQS?.confirmations || [])];
    const pending = new Set(requests.filter(r => !r.refused && !r.recorded && (!r.expires || r.expires > REQS.now)).map(r => r.episode).filter(Boolean));
    const episodes = new Map(), lanes = new Set();
    for (const node of Object.values(N)) {
      if (!pending.has(node.episode)) continue;
      if (!episodes.has(node.episode)) episodes.set(node.episode, new Set());
      episodes.get(node.episode).add(node.lane);
    }
    for (const owners of episodes.values()) if (owners.size === 1) owners.forEach(lane => lanes.add(lane));
    return lanes;
  }
  const camera = WorldCameraController.create($('stage'), $('head-space'), {
    model:() => T, timeXs:() => XS, sidePoints, attention:attentionLanes, reduced,
    glyph:(n, size) => glyphOfNode(n, size), title:n => cleanTitle(n.title) || n.episode || n.label || n.id,
    summary:summaryOf, agent:agentLine, tag:tagOf, lastRunLabel, lastRunDetail, mergeFileSummary, stops:() => miniStops,
    open:id => open(id),
    select:id => { if (N[id]) { st.focus = id; camera.select(id); } },
    placeStatus:placeNow,
    timeChanged: (index, active) => updateHeadRibbon(index, active && st.d === 0),
    sideVisible:keyboard => {
      $('L0').inert = false;
      if (keyboard) $('timeline-angle').focus({preventScroll:true});
    }
  });

  function markAngle(which) {
    for (const [id, name] of [['timeline-angle','timeline'], ['head-angle','heads']]) {
      if (name === which) $(id).setAttribute('aria-current', 'page');
      else $(id).removeAttribute('aria-current');
    }
  }
  function cameraURL(url, push) {
    if (location.pathname + location.search + location.hash === url) return;
    history[push ? 'pushState' : 'replaceState'](null, '', url);
  }
  function faceFor(id) {
    const space = $('head-space');
    const exact = space.querySelector(`[data-head-state="${CSS.escape(id)}"],[data-rear-state="${CSS.escape(id)}"],[data-lineage-state="${CSS.escape(id)}"]`);
    if (exact) return exact;
    const lane = N[id]?.lane;
    return lane && space.querySelector(`[data-head-lane="${CSS.escape(lane)}"]`);
  }
  function animateStack(face, back) {
    stackFlight?.cancel();
    if (reduced() || !face || !face.width) return;
    const front = $('front').getBoundingClientRect();
    if (!front.width || !front.height) return;
    const zoom = face.width / front.width;
    const x = face.left + face.width / 2 - (front.left + front.width / 2), y = face.top - front.top;
    const crop = Math.max(0, front.height - face.height / zoom);
    const small = {transform:`translate(${x}px,${y}px) scale(${zoom})`, transformOrigin:'50% 0', clipPath:`inset(0 0 ${crop}px)`,opacity:.4};
    const full = {transform:'none', transformOrigin:'50% 0', clipPath:'inset(0 0 0px)',opacity:1};
    stackFlight = $('front').animate(back ? [full,small] : [small,full], {duration:640, easing:'cubic-bezier(.34,.04,.21,1)'});
    return stackFlight;
  }
  function chooseAngle(heads, {navigate=true, instant=false, keyboard=false, anchor=null} = {}) {
    if (!T) return;
    openSeq++;
    const wasStack = st.d === 1;
    $('stage').classList.remove('head-approach','head-return');
    stackFlight?.cancel();
    if (wasStack && heads && !instant && !reduced()) $('stage').classList.add('head-return');
    if (wasStack) {
      resetFinishLinks();
      dismissFilePreview();
      closeFile(false);
      timelineAnchor = N[F.id]?.id || T.aliasOf[F.id] || '';
      st.d = 0;
      $('root').classList.remove('zoomed');
      $('L1').inert = true;
      renderCrumbs();
    }
    const chosen = anchor || timelineAnchor;
    const node = N[chosen] || N[T.aliasOf[chosen]];
    if (node) st.focus = node.id;
    if (wasStack && fromHeads && !heads && !camera.active) camera.turn(true, st.focus, {instant:true});
    if (!heads && st.focus && REP[st.focus]) {
      pinnedTimelineState = st.focus;
      relayout();
      ensureVisible(st.focus, false);
      if (narrow.matches) {
        const y = DISP[REP[st.focus]].y;
        const layer = $('L0');
        if (y < layer.scrollTop + 40 || y > layer.scrollTop + layer.clientHeight - 40) layer.scrollTop = y - layer.clientHeight / 2;
      }
      applyFocus();
    }
    $('L0').inert = heads || camera.active;
    $('tip').hidden = true;
    markAngle(heads ? 'heads' : 'timeline');
    if (navigate) cameraURL(heads ? '/surface#heads' : '/surface', true);
    const id = N[st.focus] ? st.focus : DISP[st.focus]?.ids?.at(-1);
    const frameNode = heads && Boolean(node);
    if (heads) timelineAnchor = '';
    camera.turn(heads, id, {instant:instant || wasStack && heads, keyboard, frameNode});
    if (!heads && !camera.active) { setView('timeline'); $('L0').inert = false; }
    if ($('stage').classList.contains('head-return')) {
      const face = faceFor(st.focus)?.getBoundingClientRect();
      const flight = animateStack(face, true);
      const finish = () => $('stage').classList.remove('head-return');
      if (flight) flight.finished.then(finish, () => {}); else finish();
    }
    fromHeads = false;
  }

  // ---------- Timeline (static frame) ----------
  const readFrom = id => [...READS.values()].some(rd => rd.fromId === id);
  const linksOf = id => {
    const r = new Map();
    if (isCluster(id)) return r;
    (parents[id] || []).forEach(p => r.set(p, 'built on'));
    (kids[id] || []).forEach(c => r.set(c, 'built on'));
    if (F && id === F.id) READS.forEach(rd => { if (N[rd.fromId] && !r.has(rd.fromId)) r.set(rd.fromId, 'read'); });
    if (F && readFrom(id) && !r.has(F.id)) r.set(F.id, 'read');
    return r;
  };
  const dayDiff = iso => { const a = new Date(iso), b = NOW_D; return Math.round((Date.UTC(b.getUTCFullYear(), b.getUTCMonth(), b.getUTCDate()) - Date.UTC(a.getUTCFullYear(), a.getUTCMonth(), a.getUTCDate())) / 864e5); };
  const lastShared = key => ((LANES[key] && LANES[key].states) || []).filter(s => s.kind === 'shared').reduce((a, s) => !a || s.rank > a.rank ? s : a, null);
  function statusOf(key) {
    const run = T.OPEN.find(n => n.lane === key && n.open && n.open.sessionLive);
    if (run) return `<span class="stat on" data-tip="node:${esc(run.id)}">session live</span>`;
    const s = isSharedKey(key) ? (key === 'shared' ? N[SH.at(-1)] : lastShared(key)) : lastStep(key);
    const on = s && dayDiff(s.at) <= 1;
    return `<span class="stat${on ? ' on' : ''}" data-tip="${s ? 'node:' + s.id : ''}">${on ? 'active' : 'quiet'}</span>`;
  }
  function progressOf(key) {
    const l = LANES[key];
    if (isSharedKey(key)) { const g = T.groupOf(key), n = (l.states || []).filter(s => s.kind === 'shared').length; return `<div class="prog">${n} shared record${n === 1 ? '' : 's'} on ${esc(g.repo ? g.repo.branch : 'dev')}</div>`; }
    const ms = l.milestones;
    if (ms && ms.names) {
      const dots = ms.names.map(nm => {
        const done = ms.doneNames.includes(nm), rev = ms.inReview === nm, nx = ms.next === nm;
        const style = done ? `background:${l.colour};border-color:${l.colour}` : rev ? `background:linear-gradient(90deg, ${l.colour} 50%, transparent 50%);border-color:${l.colour}` : nx ? `border-color:${l.colour}` : 'border-color:var(--line)';
        return `<span class="msd" style="${style}"></span>`;
      }).join('');
      return `<div class="prog" data-tip="ms:${key}"><span class="mst">${dots}</span><b>${ms.done} of ${ms.total}</b></div>
        <div class="prog sub">${esc(ms.inReview)} in review · next ${esc(ms.next)}</div>`;
    }
    const n = l.stepsLanded != null ? l.stepsLanded : l.states.length;
    return `<div class="prog">${n} step${n === 1 ? '' : 's'} landed</div>`;
  }
  function laneHead(key, y) {
    const l = LANES[key];
    return `<button type="button" class="lhead" data-lane="${esc(key)}" data-brane="${esc(key)}" style="top:${y - 32}px">
      <span class="l1">${glyphOf(key, 11)}<span class="nm"${T.multi && l.goal ? ` data-tip="goal:${esc(key)}"` : ''}>${esc(nameOf(key))}</span></span>
      ${l.goal && !T.multi ? `<span class="goal" data-tip="goal:${esc(key)}">${esc(l.goal)}</span>` : ''}
      ${progressOf(key).replace('</div>', statusOf(key) + '</div>')}</button>`;
  }
  function otherHead(y, row) {
    const quiet = T.QUIET[row] || [];
    const steps = quiet.reduce((a, k) => a + (LANES[k].stepsLanded != null ? LANES[k].stepsLanded : LANES[k].states.length), 0);
    const last = quiet.map(lastStep).filter(Boolean).sort((a, b) => b.rank - a.rank)[0];
    return `<button type="button" class="lhead" data-quiet="${esc(row)}" aria-expanded="false" style="top:${y - 24}px">
      <span class="l1"><span class="nm">${quiet.length} quieter branes</span><span class="gs">${quiet.map(k => glyphOf(k, 9)).join('')}</span></span>
      <span class="prog">${steps} steps landed${last ? ' · last ' + esc(last.rel) : ''}</span></button>`;
  }

  function renderTimeline(keepView = false) {
    let h = '';
    ROWS.forEach((key, k) => {
      const y = LY[key], tall = T.multi ? T.GAP : T.isOther(key) ? 72 : 98;
      h += `<div class="band" style="top:${y - tall / 2 - 4}px;height:${tall}px;width:${PR}px;background:${k % 2 === 0 ? 'var(--laneBand)' : 'transparent'}"></div>`;
      h += `<div class="lrule" style="top:${y}px;left:${PL}px;width:${PW}px"></div>`;
      h += T.isOther(key) ? otherHead(y, key) : laneHead(key, y);
    });
    h += `<div class="pedge" style="left:${PL - 1}px;top:96px;height:${TICK_T + 36 - 96}px"></div><div class="pedge" style="left:${PR}px;top:96px;height:${TICK_T + 36 - 96}px"></div>`;
    h += `<div class="plot" id="plot" style="left:${PL}px;top:0;width:${PW}px;height:${TICK_T + 40}px"><div class="world" id="world"></div></div>`;
    h += `<div class="tlbar" style="left:${PL}px;top:62px;width:${PW}px">
      <button type="button" class="phone-view" id="phone-activity">Activity</button><span class="range" id="range"></span><span style="flex:1"></span>
      <span class="ctl">
        <button type="button" class="cb" id="pg-back" aria-label="Back in time" data-tip="k:← back a page · drag or scroll to pan">‹</button>
        <button type="button" class="cb" id="pg-fwd" aria-label="Forward in time" data-tip="k:→ forward a page">›</button>
        <button type="button" class="cb now-b" id="now-b" aria-label="Back to now" data-tip="k:End · back to the newest state">Now</button>
      </span>
      <span id="timeline-zoom"></span></div>`;
    h += '<div id="timeline-now"></div>';
    h += `<div class="mini-lab" style="top:${MINI_T + 2}px">Overview<span>drag the window · click to jump</span></div>`;
    h += '<div id="timeline-mini"></div>';
    ROWS.filter(T.isOther).forEach(row => {
      h += `<div class="quiet-list" data-quiet-list="${esc(row)}" style="left:16px;bottom:${936 - (LY[row] - 28)}px" hidden>${(T.QUIET[row] || []).map(k => {
        const s = lastStep(k);
        return `<div data-brane="${esc(k)}">${glyphOf(k, 11)}<span class="qn"><b>${esc(nameOf(k))}</b><span class="qg">${esc(LANES[k].goal || (LANES[k].row ? 'row ' + LANES[k].row : ''))}</span></span><span class="w">${esc(s ? s.rel : '')}</span></div>`;
      }).join('')}</div>`;
    });
    h += `<div class="legend" aria-label="Legend" style="top:${LEGEND_T}px">
      <span>${mergeGlyph(11, true)}GitHub merge</span>
      <span>${gl('square', 'var(--sub)', 10)}World episode or work step</span>
      <span><span class="nb inl">1</span>your decisions</span>
      <span><span style="width:14px;height:14px;border-radius:7px;box-shadow:0 0 0 4px var(--halo)"></span>connected to the step you point at</span>
      <span><span style="width:18px;border-top:2px dotted var(--read)"></span>read</span>
      <span><span class="bub inl">4</span>several records · click to open up</span>
      <span class="off"><span style="width:18px;border-top:2px dashed var(--ran)"></span>ran on <i>· not recorded yet</i></span>
    </div>`;
    $('tl').innerHTML = h;
    repaintNow();
    // The plot clips; it must never scroll itself (focus() or scrollIntoView on a node would).
    $('plot').addEventListener('scroll', e => { e.target.scrollTop = 0; e.target.scrollLeft = 0; });
    mini.style.cssText = `left:${PL}px;top:${MINI_T}px;width:${PW}px;height:${MINI_H}px`;
    placeRibbon();
    drawMini();
    relayout();
    if (!keepView) goNow(false);
  }

  function nowCard() {
    const ne = D.nowEdge || {}, now = NOW_D, sess = ne.sessions || [], td = ne.today;
    const healthy = /\bhealthy\b/i.test(ne.health);
    let today = '';
    if (td) {
      const by = Object.entries(td.byAgent || {}).map(([a, n]) => `${cap(a)} ${n}`).join(', ');
      today = `<div class="nsec" data-today><div class="nh">Today</div>
        <div class="nl"><b>${td.landed} step${td.landed === 1 ? '' : 's'} landed</b>${by ? ` · ${esc(by)}` : ''}</div>
        <div class="nl sub">${td.firstRound} first-round pass${td.firstRound === 1 ? '' : 'es'} · ${td.fixRounds} fix round${td.fixRounds === 1 ? '' : 's'} · ${Math.round(td.agentMinutes)} agent-minutes</div></div>`;
    }
    const run = sess[0], runLane = run && Object.keys(LANES).find(k => LANES[k].row && new RegExp('^w' + LANES[k].row + '[a-z]?$').test(run.episode));
    const running = ne.open ? openRunningHTML(ne.open) : ne.sessionsUnavailable ? `<div class="nsec" data-running><div class="nh">Live sessions</div><div class="nl sub">${esc(ne.sessionsUnavailable)}</div></div>`
      : `<div class="nsec" data-running><div class="nh">Live sessions · ${sess.length}</div>
        ${run ? `<div class="nl"><b>${esc(run.episode)}</b>${runLane ? ` · ${esc(nameOf(runLane))}` : ''}</div><div class="nl sub">started ${esc(run.started)}</div>` : '<div class="nl sub">No live sessions</div>'}</div>`;
    const activity = running + needsYouHTML(REQS, REQS ? REQS.now : 0, notifyPermission(), !toolbar.hidden);
    return `<div class="now" data-now style="left:${PR + 22}px;top:56px">
      <button type="button" class="phone-view" id="activity-back">← Timeline</button>
      <div class="now-overview"><div class="now-clock"><div class="nh">Now</div><div style="font-size:15px;font-weight:700">${T.nowKnown ? DOW[now.getUTCDay()] + ' ' + dayOf(now) + ' · ' + hm(now) : 'Time unavailable'}</div></div>
      <div class="now-health" style="display:flex;gap:9px;align-items:center"><span style="width:9px;height:9px;border-radius:5px;background:${healthy ? 'var(--okDot)' : 'var(--you)'};flex:none"></span><span data-health style="font-weight:700;color:var(--ok)">${healthy ? 'Healthy' : esc(ne.health)}</span><span style="font-size:12.5px;color:var(--sub)">checked ${esc(ne.checked)}</span></div>
      ${serviceLineHTML(ne.service)}
      <div class="nl sub" data-live>${esc(liveStatus(live))}</div>
      ${poll.note ? `<div class="nl sub" data-poll style="color:var(--you)">${esc(poll.note)}</div>` : ''}
      </div>
      ${sourcesHTML(D.repos)}
      ${today}
      ${activity}
      ${decidedHTML(REQS, decidedAll)}
    </div>`;
  }

  // The empty and refused states: the reason, quietly, where the timeline would be.
  function renderVoid(kind, reason, repos) {
    T = null;
    $('tl').innerHTML = `<div data-void style="position:absolute;left:${PL}px;top:360px;width:${PW}px;display:flex;flex-direction:column;align-items:center;gap:6px;text-align:center">
      <div style="font-size:15px;font-weight:700;color:var(--sub)">${kind === 'empty' ? 'Nothing to show yet' : 'The world could not be read'}</div>
      <div style="font-size:13px;color:var(--sub);max-width:640px">${esc(reason)}</div>
      ${(repos || []).map(r => `<div style="font-size:12.5px;color:${r.warn ? 'var(--you)' : 'var(--sub)'};max-width:640px" data-repo="${esc(r.id)}">${esc(r.line)}</div>`).join('')}</div>`;
    setView(kind);
    return false;
  }

  // ---------- The moving part: nodes, labels, axis ----------
  function relayout() {
    computeXS(view.u); layout();
    const W = Math.ceil(XS[NS] + 260);
    let h = `<svg class="edges" id="edges" width="${W}" height="936" style="position:absolute;left:0;top:0;pointer-events:none"></svg>`;
    // Axis breaks (long quiet gaps).
    for (let i = 1; i <= NS; i++) if (isBreak(i)) {
      const days = Math.round((TT[i] - TT[i - 1]) / 864e5), mx = (XS[i - 1] + XS[i]) / 2;
      const sy = LY[T.GROUPS[0].shared];
      h += `<div class="brk" style="left:${mx - 40}px;top:${sy - 30}px">${days} quiet days</div><div class="brk-z" style="left:${mx - 4}px;top:${sy - 9}px"></div>`;
    }
    Object.entries(DISP).forEach(([id, d]) => {
      const { x, y } = d;
      if (isCluster(id)) {
        const ns = d.ids.map(k => N[k]), shared = isSharedKey(d.row);
        const word = shared ? 'shared records' : 'work records';
        const col = shared || T.isDip(d.row) || T.isOther(d.row) ? 'var(--edge)' : LANES[d.row].colour;
        const badge = ns.reduce((a, n) => a + (n.badge || 0), 0);
        h += `<button type="button" class="node bub-n" data-id="${id}" data-tip="node:${id}" style="left:${x}px;top:${y}px" aria-label="${ns.length} ${word}, ${esc(exactOf(ns[0].at))} to ${esc(exactOf(ns[ns.length - 1].at))}"><span class="halo"><span class="bub" style="border-color:${col}">${ns.length}</span></span>${d.ids.map(k => `<span data-state="${esc(k)}" hidden></span>`).join('')}</button>`;
        if (badge) h += `<div class="nb" style="left:${x + 8}px;top:${y - 21}px">${badge}</div>`;
        return;
      }
      const n = N[id], shared = n.kind === 'shared', unfiled = !shared && isSharedKey(n.lane), opened = n.kind === 'open';
      const l = LANES[n.lane];
      const size = shared ? 9 : unfiled ? 9 : l.glyph === 'square' ? 15 : 18;
      // An open episode: a dashed ring in its lane's colour, filled while its session is live. Never a commit glyph.
      const record = T.recordKind(n);
      const inner = record === 'github-merge' ? mergeGlyph(13, true)
        : shared ? gl('circle', SHARED_INK, size, 'background:var(--shellBg);')
        : opened ? glyphOfNode(n, 14)
        : unfiled ? gl(null, null, 9) : gl(l.glyph, l.colour, size);
      const title = cleanTitle(n.title) || n.episode || n.label;
      const pr = n.pr && title !== '#' + n.pr && 'PR #' + n.pr;
      const what = [T.recordLabel(n), T.activityLabel(n), opened && n.open?.label, lastRunDetail(n.open?.lastRun), title, pr].filter(Boolean).join(', ');
      h += `<button type="button" class="node${shared ? ' sh' : ''}" data-id="${id}" data-state="${esc(id)}" data-record-kind="${record}"${opened ? ` data-open="${esc(n.open?.state || '')}"` : ''} data-tip="node:${id}" style="left:${x}px;top:${y}px" aria-label="${esc(nameOf(n.lane) + ', ' + what + (n.rel ? ', ' + n.rel : '') + (n.badge ? ', ' + n.badge + ' of your decisions' : ''))}"><span class="halo">${inner}</span></button>`;
      if (n.badge > 0) h += `<div class="nb" style="left:${x + 6}px;top:${y - 19}px">${n.badge}</div>`;
      if (LABELED.has(id)) h += `<div class="nlab" data-l="${id}" style="left:${x - 56}px;bottom:${936 - (y - 18)}px"><span class="tx">${esc(opened ? n.episode + ' · ' + (lastRunLabel(n.open.lastRun) || n.open.label) : shortTitle(n.title || '') || n.episode || n.label)}</span></div>`;
      h += `<div class="nlab" data-r="${id}" style="left:${x - 56}px;top:${y + (shared ? 9 : 14)}px"><div class="rel"></div></div>`;
    });
    h += `<div class="nowline" style="left:${XS[T.presentIndex]}px;top:92px;height:${TICK_T + 30 - 92}px"></div>`;
    h += ticksHTML();
    $('world').style.width = W + 'px';
    $('world').innerHTML = h;
    $('zl').textContent = levelName();
    $('z-in').disabled = view.u <= 0.001; $('z-out').disabled = view.u >= MAXU - 0.001;
    applyFocus();
  }

  function ticksHTML() {
    const lv = level(), out = [];
    const place = (x, text, cls, sep, placed) => { if (placed.some(p => Math.abs(p - x) < sep)) return false; placed.push(x); out.push(`<div class="${cls}" style="left:${x - 40}px;top:${cls === 'tick-d' ? TICK_T + 16 : TICK_T}px">${text}</div>`); return true; };
    const pt = [], pd = [];
    if (lv === 'steps') {
      let lastDay = '';
      const times = T.dated ? TT.map((t, i) => ({t, i})).slice(T.undated ? 1 : 0) : [];
      if (T.undated) place(XS[0], 'Undated', 'tick-d', 76, pd);
      times.forEach(({t, i}) => {
        const d = dayOf(t), x = XS[i];
        if (d !== lastDay) { place(x, d, 'tick-d', 76, pd); lastDay = d; place(x, hm(t), 'tick-t', 44, pt); }
      });
      times.forEach(({t, i}) => place(XS[i], hm(t), 'tick-t', 44, pt));
    } else {
      if (T.undated) place(XS[0], 'Undated', 'tick-d', 84, pd);
      const cands = [];
      for (const [t, f] of tickTimes(lv === 'weeks' ? 864e5 : lv === 'days' ? 12 * 3.6e6 : 3.6e6)) {
        const d = new Date(t), hh = d.getUTCHours(), mid = hh === 0, mon = mid && d.getUTCDay() === 1;
        const pri = lv === 'weeks' ? (mon ? 0 : mid ? 1 : 9) : mid ? 0 : hh % 12 === 0 ? 1 : hh % 6 === 0 ? 2 : hh % 3 === 0 ? 3 : 4;
        if (pri < 9) cands.push({ t, x: posF(f), pri, d, mid, mon });
      }
      cands.sort((a, b) => a.pri - b.pri || a.t - b.t).forEach(c => {
        const ds = `${DOW[c.d.getUTCDay()]} ${c.d.getUTCDate()} ${MON[c.d.getUTCMonth()]}`;
        if (lv === 'weeks') {
          if (c.mon) { place(c.x, 'w/c ' + c.d.getUTCDate() + ' ' + MON[c.d.getUTCMonth()], 'tick-d', 90, pd); }
          place(c.x, `${DOW[c.d.getUTCDay()][0]} ${c.d.getUTCDate()}`, 'tick-t', 30, pt);
        } else if (c.mid) { place(c.x, ds, 'tick-d', 84, pd); if (lv === 'hours') place(c.x, '00:00', 'tick-t', 46, pt); else pt.push(c.x); }
        else if (lv === 'hours' || c.pri <= 1) place(c.x, pad(c.d.getUTCHours()) + ':00', 'tick-t', lv === 'days' ? 60 : 46, pt);
      });
      // Tick marks on the axis rule.
      cands.filter(c => c.mid).forEach(c => out.push(`<div class="tmark" style="left:${c.x}px;top:${TICK_T - 8}px"></div>`));
    }
    return out.join('');
  }

  function applyFocus() {
    const f = st.focus, fr = REP[f] || f, linked = linksOf(f);
    const linkedR = new Set([...linked.keys()].map(k => REP[k]).filter(Boolean));
    const lit = (pr, cr) => (pr === fr && linkedR.has(cr)) || (cr === fr && linkedR.has(pr));
    let e = '';
    const W = XS[NS] + 260;
    T.GROUPS.forEach(g => {
      if (LANES[g.shared]) e += `<path d="M 0 ${LY[g.shared]} L ${XS[0]} ${LY[g.shared]}" fill="none" stroke="var(--edge)" stroke-width="1.4" stroke-dasharray="3 5"></path>`;
    });
    const seen = new Set();
    Object.keys(POS).forEach(c => {
      const n = N[c], rc = REP[c], [x2, y2] = [DISP[rc].x, DISP[rc].y];
      const ps = parents[c] || [];
      if ((n.kind === 'step' || n.kind === 'open') && !ps.length && !T.isDip(T.rowOf(c))) {  // its base is further back than this data: it arrives from the left
        e += `<path d="M ${x2 - 70} ${y2} L ${x2} ${y2}" fill="none" stroke="${isSharedKey(n.lane) ? 'var(--edge)' : LANES[n.lane].colour}" stroke-width="1.4" stroke-dasharray="3 5" opacity="0.6"></path>`;
      }
      ps.forEach(p => {
        const rp = REP[p]; if (!rp || rp === rc) return;
        const key = rp + '>' + rc; if (seen.has(key)) return; seen.add(key);
        const [x1, y1] = [DISP[rp].x, DISP[rp].y], L = lit(rp, rc);
        const stepish = k => k === 'step' || k === 'open';
        const ownKey = stepish(n.kind) ? n.lane : stepish(N[p].kind) ? N[p].lane : 'shared';
        const col = L ? 'var(--lit)' : isSharedKey(ownKey) ? 'var(--edge)' : LANES[ownKey].colour;
        const w = L ? 2.6 : 1.6, op = L || isSharedKey(ownKey) ? 1 : 0.6;
        let d;
        const brk = (n.kind === 'shared' && N[p].kind === 'shared' && breaksBetween(T.timeIndex.get(p), T.timeIndex.get(c))) || n.kind === 'open';
        if (x1 === x2 || y1 === y2) d = `M ${x1} ${y1} L ${x2} ${y2}`;
        else if (n.kind === 'shared') d = `M ${x1} ${y1} C ${x2} ${y1}, ${x2} ${y2}, ${x2} ${y2}`;
        else { const k = Math.sign(x2 - x1) * Math.min(30, Math.abs(x2 - x1) / 2); d = `M ${x1} ${y1} C ${x1 + k} ${y1}, ${x1} ${y2}, ${x1 + k} ${y2} L ${x2} ${y2}`; }
        e += `<path d="${d}" fill="none" stroke="${col}" stroke-width="${w}" opacity="${op}"${brk ? ' stroke-dasharray="3 5"' : ''}></path>`;
      });
    });
    // Read arcs: dotted purple, lifted above both ends (the artifact's Q curve).
    linked.forEach((via, o) => {
      if (via !== 'read' || !REP[o] || !REP[f]) return;
      const a = DISP[REP[o]], b = DISP[REP[f]];
      if (a === b) return;
      e += `<path d="M ${a.x} ${a.y} Q ${(a.x + b.x) / 2} ${Math.min(a.y, b.y) - 70}, ${b.x} ${b.y}" fill="none" stroke="var(--read)" stroke-width="2" stroke-dasharray="2 4"></path>`;
    });
    $('edges').setAttribute('width', W);
    $('edges').innerHTML = e;
    document.querySelectorAll('#tl .node').forEach(b => {
      const id = b.dataset.id, isF = id === fr, isL = linkedR.has(id);
      b.querySelector('.halo').style.boxShadow = isF ? '0 0 0 3px var(--shellBg), 0 0 0 5px var(--lit)' : isL ? '0 0 0 4px var(--halo)' : 'none';
      const ring = b.querySelector('.ring'); if (ring) ring.style.borderColor = isF || isL ? 'var(--lit)' : 'var(--edge)';
    });
    document.querySelectorAll('#tl .nlab[data-l]').forEach(el => {
      const id = el.dataset.l, isF = id === fr, isL = linkedR.has(id);
      el.style.color = isF || isL ? 'var(--text)' : 'var(--sub)';
      el.style.fontWeight = isF ? '700' : '500';
    });
    document.querySelectorAll('#tl .nlab[data-r]').forEach(el => {
      const id = el.dataset.r; el.querySelector('.rel').textContent = id === fr ? (N[id].kind === 'shared' ? `${dayOf(N[id].at)} · ${hm(N[id].at)}` : N[id].rel) : '';
    });
  }

  // ---------- Panning, zooming, the overview strip ----------
  const minL = () => XS[0] - 150, maxL = () => XS[NS] + 4 - PW;
  const clampL = L => { const lo = minL(), hi = maxL(); return hi < lo ? hi : Math.max(lo, Math.min(hi, L)); };
  const nowL = () => XS[T.presentIndex] + 4 - PW;
  const atNow = () => T != null && Math.abs(view.L - clampL(nowL())) < 1;
  let anim = 0;
  const ease = t => 1 - Math.pow(1 - t, 3);
  function applyView(keepPosition = false) {
    if (!keepPosition) view.L = clampL(view.L);
    $('world').style.transform = `translateX(${-view.L}px)`;
    const fL = fAt(view.L), fR = fAt(view.L + PW), mw = $('mwin');
    if (mw && headRibbon.hidden) updateMiniWindow(fL, fR);
    const tL = timeOfF(fL), tR = timeOfF(fR), nowIn = T.nowKnown && (atNow() || tR === TNOW);
    $('range').innerHTML = `<b>${fmtRange(tL)}</b> → <b>${nowIn ? 'now' : fmtRange(tR)}</b>${atNow() ? '' : `<span class="behind">${behindText(tR)}</span>`}`;
    $('now-b').classList.toggle('live', !atNow());
    $('pg-fwd').disabled = view.L >= clampL(maxL()) - 1; $('pg-back').disabled = view.L <= minL() + 1;
  }
  const timeOfF = f => { const i = Math.floor(Math.max(0, Math.min(NS, f))), j = Math.min(NS, i + 1); return !T.dated || (T.undated && i === 0) ? null : TT[i] + (TT[j] - TT[i]) * (f - i); };
  const fmtRange = t => { if (t == null) return 'Undated'; const d = new Date(t); return `${DOW[d.getUTCDay()]} ${d.getUTCDate()} ${MON[d.getUTCMonth()]} ${pad(d.getUTCHours())}:${pad(d.getUTCMinutes())}`; };
  const behindText = tR => { if (tR == null || !T.nowKnown) return ''; const h = (TNOW - tR) / 3.6e6; return h < 1 ? '' : h < 36 ? `${Math.round(h)} h before now` : `${Math.round(h / 24)} days before now`; };
  function panTo(L, animate = true) {
    cancelAnimationFrame(anim);
    L = clampL(L);
    if (!animate || reduced()) { view.L = L; applyView(); return; }
    const L0 = view.L, t0 = performance.now(), dur = 420;
    const step = now => { const k = Math.min(1, (now - t0) / dur); view.L = L0 + (L - L0) * ease(k); applyView(); if (k < 1) anim = requestAnimationFrame(step); };
    anim = requestAnimationFrame(step);
  }
  const goNow = (animate = true) => panTo(nowL(), animate);
  function setZoom(u, ax) {   // ax: the design-x on screen that stays put
    const f = fAt(view.L + (ax - PL));
    view.u = Math.max(0, Math.min(MAXU, u));
    relayout();
    view.L = clampL(posF(f) - (ax - PL));
    applyView();
    if (!headRibbon.hidden) camera.update();
  }
  let zanim = 0;
  function zoomTo(u, ax, animate = true) {
    cancelAnimationFrame(zanim); cancelAnimationFrame(anim);
    u = Math.max(0, Math.min(MAXU, u));
    if (!animate || reduced()) return setZoom(u, ax);
    const u0 = view.u, t0 = performance.now(), dur = 260;
    const step = now => { const k = Math.min(1, (now - t0) / dur); setZoom(u0 + (u - u0) * ease(k), ax); if (k < 1) zanim = requestAnimationFrame(step); };
    zanim = requestAnimationFrame(step);
  }
  const screenX = id => { const d = DISP[REP[id] || id]; return d ? PL + d.x - view.L : null; };
  const keyAnchor = () => { const x = screenX(st.focus); return x != null && x >= PL && x <= PR ? x : PL + PW / 2; };
  const zoomStep = dir => {
    const target = dir < 0 ? Math.ceil(view.u - 0.01) - 1 : Math.floor(view.u + 0.01) + 1;
    zoomTo(target, keyAnchor(), headRibbon.hidden);
  };
  function ensureVisible(id, animate = true) {
    const x = screenX(id); if (x == null) return;
    if (x < PL + 60) panTo(view.L - (PL + 60 - x) - PW * 0.25, animate);
    else if (x > PR - 40) panTo(view.L + (x - (PR - 40)) + PW * 0.25, animate);
  }

  // Overview: a fixed real-time layout of the whole history (the Days level), normalised to the strip.
  const mini = document.createElement('div'), zoom = document.createElement('span'), headRibbon = document.createElement('div');
  mini.id = 'mini';
  mini.className = 'mini';
  zoom.className = 'ctl zoom';
  zoom.setAttribute('role', 'group');
  zoom.setAttribute('aria-label', 'Zoom');
  zoom.innerHTML = `
    <button type="button" class="cb" id="z-out" aria-label="Zoom out" data-tip="k:− · ⌘-scroll or pinch to zoom">−</button>
    <span class="zl" id="zl">Steps</span>
    <button type="button" class="cb" id="z-in" aria-label="Zoom in" data-tip="k:+ · ⌘-scroll or pinch to zoom">+</button>`;
  headRibbon.id = 'head-ribbon';
  headRibbon.hidden = true;
  $('keys').before(headRibbon);

  let miniWidth = 0, headIndex, miniStops = [];
  function measureStops() {
    const at = new Map();
    for (const n of Object.values(N)) {
      const i = T.timeIndex.get(n.id);
      if (n.draft || (T.undated && i === 0) || i >= T.presentIndex || at.has(i)) continue;
      at.set(i, n.id);
    }
    at.set(T.presentIndex, '');
    miniStops = [...at].map(([index, id]) => ({index, id})).sort((a, b) => a.index - b.index);
  }
  function placeRibbon() {
    const inHeads = !headRibbon.hidden;
    const miniSlot = inHeads ? headRibbon : $('timeline-mini');
    const zoomSlot = inHeads ? headRibbon : $('timeline-zoom');
    if (miniSlot && mini.parentElement !== miniSlot) miniSlot.append(mini);
    if (zoomSlot && zoom.parentElement !== zoomSlot) zoomSlot.append(zoom);
  }
  function showHeadRibbon(show) {
    if (show === !headRibbon.hidden) return;
    headRibbon.hidden = !show;
    placeRibbon();
    if (show) {
      mini.tabIndex = 0;
      mini.setAttribute('role', 'slider');
      mini.setAttribute('aria-label', 'Recorded time');
    } else {
      for (const attr of ['tabindex', 'role', 'aria-label', 'aria-valuemin', 'aria-valuemax', 'aria-valuenow', 'aria-valuetext']) mini.removeAttribute(attr);
    }
    if (!T) return;
    for (const [id, key] of [['z-out', '−'], ['z-in', '+']]) {
      $(id).dataset.tip = show ? '' : `k:${key} · ⌘-scroll or pinch to zoom`;
    }
    drawMini();
  }
  function drawMini() {
    miniWidth = headRibbon.hidden ? PW : mini.clientWidth;
    measureMini(miniWidth);
    renderMini();
    if (headRibbon.hidden) updateMiniWindow(fAt(view.L), fAt(view.L + PW));
    else updateMiniMarker();
  }
  function updateMiniWindow(fL, fR) {
    const a = MX(fL), b = MX(fR), mw = $('mwin');
    mw.style.left = a + 'px';
    mw.style.width = Math.max(8, b - a) + 'px';
  }
  function updateHeadRibbon(index, active) {
    headIndex = index;
    if (!T || !mini.isConnected) return;
    const changed = active === headRibbon.hidden;
    showHeadRibbon(active);
    if (active && !changed) {
      if (miniWidth !== mini.clientWidth) drawMini();
      else updateMiniMarker();
    }
  }
  function updateMiniMarker() {
    const x = MXS[headIndex];
    if (x == null) return;
    const text = headIndex === T.presentIndex ? 'Present' : fmtRange(timeOfF(headIndex));
    mini.style.setProperty('--mini-cursor', x + 'px');
    $('mhead-date').textContent = text;
    mini.setAttribute('aria-valuemin', String(miniStops[0].index));
    mini.setAttribute('aria-valuemax', String(T.presentIndex));
    mini.setAttribute('aria-valuenow', String(Math.max(headIndex, miniStops[0].index)));
    mini.setAttribute('aria-valuetext', headIndex === T.presentIndex || timeOfF(headIndex) == null ? text : text + ' UTC');
  }
  let miniDragId = null, miniDragIndex = null;
  function scrubHeadTime(event) {
    const x = event.clientX - mini.getBoundingClientRect().left - mini.clientLeft;
    const stop = miniStops.reduce((best, next) =>
      Math.abs(MXS[next.index] - x) <= Math.abs(MXS[best.index] - x) ? next : best);
    if (stop.index === miniDragIndex) return;
    miniDragIndex = stop.index;
    camera.seek(stop.id);
  }
  mini.addEventListener('pointerdown', event => {
    if (headRibbon.hidden || event.button !== 0 || !event.isPrimary) return;
    event.stopPropagation();
    event.preventDefault();
    mini.focus({preventScroll:true});
    miniDragId = event.pointerId;
    miniDragIndex = null;
    mini.setPointerCapture(event.pointerId);
    scrubHeadTime(event);
  });
  mini.addEventListener('pointermove', event => {
    if (event.pointerId === miniDragId && !headRibbon.hidden) scrubHeadTime(event);
  });
  mini.addEventListener('pointerup', event => {
    if (event.pointerId === miniDragId && !headRibbon.hidden) scrubHeadTime(event);
  });
  mini.addEventListener('lostpointercapture', event => {
    if (event.pointerId === miniDragId) miniDragId = null;
  });
  mini.addEventListener('keydown', event => {
    if (headRibbon.hidden || event.ctrlKey || event.altKey || event.metaKey) return;
    const key = event.key;
    let stop;
    if (key === 'Home') stop = miniStops[0];
    else if (key === 'End') stop = miniStops.at(-1);
    else if (key === 'ArrowRight' || key === 'ArrowUp') stop = miniStops.find(s => s.index > headIndex);
    else if (key === 'ArrowLeft' || key === 'ArrowDown') stop = miniStops.findLast(s => s.index < headIndex);
    else return;
    event.preventDefault();
    event.stopPropagation();
    if (stop) camera.seek(stop.id);
  });
  const MX = f => { f = Math.max(0, Math.min(NS, f)); const i = Math.floor(f), j = Math.min(NS, i + 1); return MXS[i] + (MXS[j] - MXS[i]) * (f - i); };
  const fOfMX = mx => { if (mx <= MXS[0]) return 0; if (mx >= MXS[NS]) return NS; let i = 0; while (i < NS && MXS[i + 1] < mx) i++; return i + (MXS[i + 1] === MXS[i] ? 0 : (mx - MXS[i]) / (MXS[i + 1] - MXS[i])); };
  function renderMini() {
    const rowH = T.multi ? 3 : 4, top = 5;
    let h = '', lastDateX = -Infinity;
    for (const [t, f] of tickTimes(864e5)) {
      const x = MX(f), d = new Date(t);
      h += `<div class="mday" style="left:${x}px"></div>`;
      if (headRibbon.hidden || (x - lastDateX >= 36 && x + 36 <= miniWidth)) {
        h += `<div class="mdl" style="left:${x + 3}px">${d.getUTCDate()} ${MON[d.getUTCMonth()]}</div>`;
        lastDateX = x;
      }
    }
    for (let i = 1; i <= NS; i++) if (isBreak(i)) h += `<div class="mbrk" style="left:${(MXS[i - 1] + MXS[i]) / 2 - 6}px">≈</div>`;
    ROWS.forEach((r, k) => { h += `<div class="mrow" style="top:${top + k * (rowH + 1)}px"></div>`; });
    for (const n of Object.values(N)) {
      const x = MXS[T.timeIndex.get(n.id)], k = ROWS.indexOf(T.homeRow(T.rowOf(n.id)));
      const width = n.kind === 'shared' ? 1 : 3, colour = isSharedKey(n.lane) ? 'var(--edge)' : LANES[n.lane].colour;
      h += `<div class="mk" style="left:${x - (width - 1) / 2}px;top:${top + k * (rowH + 1)}px;height:${rowH}px;width:${width}px;background:${colour}"></div>`;
    }
    h += `<div class="mnow" style="left:${MXS[T.presentIndex]}px"></div><div class="mwin" id="mwin"></div>
      <div class="mhead" aria-hidden="true"><span id="mhead-date"></span></div>`;
    mini.innerHTML = h;
  }

  // ---------- Brane view ----------
  function placeBars(animate) {
    const bars = stackBars(S.earlier, -1).concat(stackBars(S.later, 1));
    const host = $('bstage'), have = new Map([...host.querySelectorAll('.strip')].map(el => [el.dataset.key, el]));
    const keep = new Set();
    bars.forEach(b => {
      keep.add(b.key);
      let el = have.get(b.key);
      const fresh = !el;
      if (fresh) {
        el = document.createElement('button'); el.type = 'button'; el.className = 'strip'; el.dataset.key = b.key; el.dataset.id = b.card.id; el.dataset.lane = b.col.brane;
        el.style.top = (b.top + (b.dir < 0 ? -46 : 46)) + 'px'; el.style.opacity = '0'; el.style.width = b.width + 'px'; el.style.marginLeft = b.ml + 'px';
        host.appendChild(el);
      }
      el.innerHTML = barHTML(b);
      el.setAttribute('aria-label', (b.dir < 0 ? 'Back to ' : 'Ahead to ') + tagOf(b.card) + (b.card.title ? ': ' + cleanTitle(b.card.title) : ''));
      el.dataset.tip = N[b.card.id] ? 'node:' + b.card.id : '';
      el.style.zIndex = b.z;
      const apply = () => { el.style.top = b.top + 'px'; el.style.width = b.width + 'px'; el.style.marginLeft = b.ml + 'px'; el.style.opacity = b.op; };
      if (fresh && animate) requestAnimationFrame(() => requestAnimationFrame(apply)); else apply();
    });
    have.forEach((el, k) => { if (!keep.has(k)) { if (animate) { el.style.opacity = '0'; setTimeout(() => el.remove(), 300); } else el.remove(); } });
  }

  let selectedFile = null, filePreview = null, restoringFileFocus = false;

  function dismissFilePreview() {
    if (!filePreview) return;
    const { button, timer, leave, request, bubble } = filePreview;
    filePreview = null;
    clearTimeout(timer);
    clearTimeout(leave);
    request.abort();
    button.removeAttribute('aria-describedby');
    bubble?.remove();
  }

  function leaveFilePreview(target) {
    const preview = filePreview;
    if (!preview || preview.button === document.activeElement || preview.button.contains(target) ||
        preview.bubble?.contains(target)) return;
    clearTimeout(preview.leave);
    preview.leave = setTimeout(() => { if (filePreview === preview) dismissFilePreview(); }, 120);
  }

  function placeFilePreview() {
    if (!filePreview?.bubble) return;
    const { button, bubble } = filePreview;
    const rect = button.getBoundingClientRect();
    const visible = button.closest('.stack-body')?.getBoundingClientRect();
    if (!visible || rect.bottom <= Math.max(0, visible.top) ||
        rect.top >= Math.min(window.innerHeight, visible.bottom)) {
      dismissFilePreview();
      return;
    }
    const width = bubble.offsetWidth, height = bubble.offsetHeight;
    const x = rect.left - width >= 16 ? rect.left - width - 8 :
      Math.max(8, Math.min(rect.left, window.innerWidth - width - 8));
    const y = rect.left - width >= 16 ? rect.top : rect.bottom + 8;
    bubble.style.left = x + 'px';
    bubble.style.top = Math.max(8, Math.min(y, window.innerHeight - height - 8)) + 'px';
  }

  function fileSelection(button) {
    const { filePath: path, fileSource: source = 'committed', fileOrigin: origin = '' } = button.dataset;
    const tab = button.closest('[data-node-panel]')?.dataset.nodePanel || '';
    if (source === 'committed' && F.landedFiles?.some(file => file.path === path)) return {path, source, tab, origin};
    if (tab === 'files' && (source === 'episode-read' || source === 'episode-write')) return {path, source, tab, origin};
    return null;
  }

  function applyFileFilters(root) {
    const query = (root.querySelector('[data-file-filter]')?.value || '').trim().toLowerCase();
    const kind = root.querySelector('[data-file-kind]')?.value || '';
    let visible = 0;
    for (const group of root.querySelectorAll('[data-file-group]')) {
      let matches = 0;
      for (const row of group.querySelectorAll('[data-file-row]')) {
        row.hidden = (kind && kind !== group.dataset.fileGroup) || !row.dataset.fileRow.toLowerCase().includes(query);
        if (!row.hidden) matches++;
      }
      group.hidden = !matches;
      visible += matches;
    }
    const empty = root.querySelector('[data-file-empty]');
    if (empty) empty.hidden = visible !== 0;
    if (filePreview?.button.closest('[hidden]')) dismissFilePreview();
  }

  function startFilePreview(button) {
    if (!button || restoringFileFocus || selectedFile || st.d !== 1 || !$('front')?.contains(button)) return;
    if (filePreview?.button === button) {
      clearTimeout(filePreview.leave);
      return;
    }
    dismissFilePreview();
    const selection = fileSelection(button);
    if (!selection) return;
    const {path, source} = selection;
    const request = new AbortController();
    const preview = { button, request };
    filePreview = preview;
    preview.timer = setTimeout(async () => {
      if (filePreview !== preview || !button.isConnected) return;
      const bubble = document.createElement('aside');
      bubble.id = 'file-hover-preview';
      bubble.className = 'file-hover-preview';
      bubble.setAttribute('role', 'tooltip');
      bubble.addEventListener('pointerenter', () => clearTimeout(preview.leave));
      bubble.addEventListener('pointerleave', event => leaveFilePreview(event.relatedTarget));
      bubble.textContent = 'Loading file preview…';
      preview.bubble = bubble;
      document.body.append(bubble);
      button.setAttribute('aria-describedby', bubble.id);
      placeFilePreview();
      if (filePreview !== preview) return;
      try {
        const result = await getJSON(filePath(F.brane, F.id, path, source), request.signal);
        if (filePreview !== preview) return;
        if (!result.ok || result.body?.path !== path) {
          bubble.textContent = 'File preview unavailable.';
        } else {
          preview.evidence = result.body;
          bubble.innerHTML = filePreviewHTML(result.body);
        }
      } catch (_) {
        if (filePreview === preview) bubble.textContent = 'File preview unavailable.';
      }
      if (filePreview === preview) placeFilePreview();
    }, 220);
  }

  function fileButton(front, selection) {
    const buttons = Array.from(front.querySelectorAll('[data-file-path]')).filter(button =>
      button.dataset.filePath === selection.path && (button.dataset.fileSource || 'committed') === selection.source &&
      (button.closest('[data-node-panel]')?.dataset.nodePanel || '') === selection.tab);
    return buttons.find(button => (button.dataset.fileOrigin || '') === selection.origin) || buttons[0];
  }

  function focusFileSource(element) {
    restoringFileFocus = true;
    element.focus({preventScroll:true});
    restoringFileFocus = false;
  }

  function closeFile(restoreFocus = true) {
    if (!selectedFile) return;
    const closing = selectedFile;
    const { panel, sectionTop, layerTop, request } = closing;
    selectedFile = null;
    request.abort();
    panel.remove();
    if ($('phone-back')) $('phone-back').textContent = fromHeads ? '← Heads' : '← Timeline';
    const front = $('front'), summary = front.querySelector('.stack-body');
    if (summary) {
      summary.hidden = false;
      if (restoreFocus) {
        const button = fileButton(front, closing);
        focusFileSource(button || summary);
        summary.scrollTop = sectionTop;
        if (narrow.matches) $('L1').scrollTop = layerTop;
      }
    }
  }

  async function openFile(button) {
    if (st.d !== 1 || !F || !$('front').contains(button) || selectedFile) return;
    const selection = fileSelection(button);
    if (!selection) return;
    const {path, source, tab} = selection;
    const previewed = filePreview?.button === button && filePreview.evidence;
    dismissFilePreview();

    const front = $('front'), summary = front.querySelector('.stack-body');
    const panel = document.createElement('section');
    panel.className = 'fbody file-evidence';
    panel.setAttribute('aria-labelledby', 'file-evidence-path');
    panel.tabIndex = 0;
    panel.innerHTML = '<div class="file-evidence-head"><button type="button" data-file-summary>Back to summary</button><h2 id="file-evidence-path"></h2></div>' +
      '<div role="status" class="muted">Loading file…</div><div class="file-evidence-body"></div>';
    panel.querySelector('#file-evidence-path').textContent = path;
    if (tab === 'files') panel.querySelector('[data-file-summary]').textContent = 'Back to files';
    const request = new AbortController();
    const selected = { ...selection, node: F.id, brane: F.brane, panel, request, sectionTop: summary.scrollTop, layerTop: $('L1').scrollTop };
    selectedFile = selected;
    summary.hidden = true;
    front.append(panel);
    $('phone-back').textContent = tab === 'files' ? '← Files' : '← Summary';
    if (narrow.matches) $('L1').scrollTop = 0;
    panel.querySelector('[data-file-summary]').focus({preventScroll:true});

    const loading = panel.querySelector('[role="status"]'), body = panel.querySelector('.file-evidence-body');
    try {
      const result = previewed ? {ok: true, body: previewed} :
        await getJSON(filePath(selected.brane, selected.node, path, source), request.signal);
      if (selectedFile !== selected) return;
      if (!result.ok || result.body?.path !== path) {
        loading.textContent = 'File preview unavailable.';
        return;
      }
      loading.textContent = '';
      body.innerHTML = fileEvidenceHTML(result.body, 'diff');
      body.addEventListener('click', event => {
        const mode = event.target.closest('[data-file-mode]')?.dataset.fileMode;
        if (mode !== 'file' && mode !== 'diff') return;
        body.innerHTML = fileEvidenceHTML(result.body, mode);
        body.querySelector(`[data-file-mode="${mode}"]`)?.focus({preventScroll:true});
      });
    } catch (_) {
      if (selectedFile === selected) loading.textContent = 'File preview unavailable.';
    }
  }

  let selectedNodeTab = 'overview';
  let episodeFiles = null;

  function tabLocation() {
    const parts = location.hash.slice(1).split('/');
    return {heads: parts[0] === 'heads', tab: (parts[0] === 'heads' ? parts[1] : parts[0]) || 'overview'};
  }

  function nodeURL() {
    const fragment = [fromHeads ? 'heads' : '', selectedNodeTab === 'overview' ? '' : selectedNodeTab].filter(Boolean).join('/');
    return appPath(F.brane, F.id) + (fragment ? '#' + fragment : '');
  }

  function finishReturnFragment(F) {
    const place = new URLSearchParams({brane:F.brane, at:F.id});
    const view = [fromHeads ? 'heads' : '', selectedNodeTab === 'overview' ? '' : selectedNodeTab].filter(Boolean).join('/');
    if (view) place.set('view', view);
    return '#' + place;
  }

  let finishCheck = null;
  function resetFinishLinks() {
    finishCheck?.request?.abort();
    finishCheck = null;
  }

  function paintFinishLinks(rows) {
    const section = $('front')?.querySelector('[data-node-confirmations]');
    if (!section || st.d !== 1) return;
    const links = section.querySelector('[data-finish-links]');
    const href = finishReturnFragment(F);
    const html = rows.map(c =>
      `<a class="finish-link" id="node-confirm-${esc(c.id)}" href="/surface/confirm/${esc(c.id)}${esc(href)}">Review ${c.draft ? 'draft pull request' : 'pull request'}<small>${esc(c.repository)} · ${esc(c.head)} @ ${esc(String(c.sha || '').slice(0, 8))}</small></a>`
    ).join('');
    if (links.dataset.markup !== html) {
      links.innerHTML = html;
      links.dataset.markup = html;
    }
    section.hidden = !rows.length;
    const history = section.closest('[data-finish-history-only]');
    if (history) history.hidden = section.hidden;
  }

  async function showFinishRequests(retry = false) {
    const episode = st.d === 1 && finishEpisode(F);
    const queue = Array.isArray(REQS?.confirmations) ? REQS.confirmations : [];
    const rows = episode ? queue.filter(row => row.episode === episode && !row.refused && CONFIRM_ID.test(row.id)).slice(0, 3) : [];
    if (!rows.length) {
      resetFinishLinks();
      paintFinishLinks([]);
      return;
    }

    const node = F, key = JSON.stringify([node.brane, node.id, rows.map(row => row.id)]);
    if (finishCheck?.key === key && (finishCheck.request || finishCheck.rows.length || !retry)) {
      paintFinishLinks(finishCheck.rows);
      return;
    }
    resetFinishLinks();
    paintFinishLinks([]);
    const request = new AbortController();
    const check = {key, request, rows: []};
    finishCheck = check;
    const timer = setTimeout(() => request.abort(), 5000);
    const results = await Promise.all(rows.map(async row => {
      try {
        const res = await getJSON('/surface/confirmations/' + row.id, request.signal);
        return res.ok && res.body?.id === row.id && res.body.state === 'pending' &&
          finishMatchesNode(res.body.episodeHref, node, T?.aliasOf) ? row : null;
      } catch (_) { return null; }
    }));
    clearTimeout(timer);
    if (finishCheck !== check) return;
    check.request = null;
    check.rows = results.filter(Boolean);
    paintFinishLinks(check.rows);
  }

  function resetEpisodeFiles() {
    episodeFiles?.request?.abort();
    episodeFiles = null;
  }

  function fileFilters(panel) {
    return {path: panel?.querySelector('[data-file-filter]')?.value || '', kind: panel?.querySelector('[data-file-kind]')?.value || ''};
  }

  function restoreFilters(panel, saved) {
    const root = panel?.querySelector('.node-files-tab');
    if (!root) return;
    const path = root.querySelector('[data-file-filter]'), kind = root.querySelector('[data-file-kind]');
    if (path) path.value = saved.path;
    if (kind) kind.value = saved.kind;
    applyFileFilters(root);
  }

  async function loadEpisodeFiles(front) {
    const episode = F.trajectoryEpisode || F.run?.episode;
    if (!episode || selectedNodeTab !== 'files') return;
    if (!episodeFiles) episodeFiles = {id:F.id, brane:F.brane, episode, checked:0};
    const state = episodeFiles;
    if (state.request || Date.now() - state.checked < 10000) return;
    const request = new AbortController();
    state.request = request;
    let data;
    try {
      const query = new URLSearchParams({brane: state.brane, at: state.id});
      const response = await getJSON('/surface/episode-files.json?' + query, request.signal);
      if (response.ok && response.body?.schema === 'world/surface-episode-files/v1' && response.body.episode === episode) data = response.body;
    } catch (_) { /* The manifest can still supply committed paths. */ }
    if (episodeFiles !== state || request.signal.aborted) return;
    state.request = null;
    state.checked = Date.now();
    data ||= {status:'unavailable'};
    if (JSON.stringify(state.data) === JSON.stringify(data)) return;
    state.data = data;
    const panel = front.querySelector('[data-node-panel=files]');
    if (!panel || F.id !== state.id || F.brane !== state.brane) return;
    const focused = document.activeElement;
    const pathControl = focused.matches?.('[data-file-filter]'), kindControl = focused.matches?.('[data-file-kind]');
    const focusedFile = panel.contains(focused) && focused.matches?.('[data-file-path]') ? fileSelection(focused) : null;
    const filters = fileFilters(panel);
    if (filePreview && panel.contains(filePreview.button)) dismissFilePreview();
    panel.innerHTML = filesTabHTML(F, data);
    restoreFilters(panel, filters);
    if (focusedFile) focusFileSource(fileButton(front, focusedFile) || panel);
    else if (pathControl || kindControl) panel.querySelector(pathControl ? '[data-file-filter]' : '[data-file-kind]')?.focus({preventScroll:true});
  }

  function selectNodeTab(tab) {
    const front = $('front'), chosen = front?.querySelector(`[data-node-tab="${tab}"]`);
    if (!chosen || st.d !== 1) return;
    const wasFile = Boolean(selectedFile);
    if (wasFile) closeFile(false);
    dismissFilePreview();
    for (const button of front.querySelectorAll('[data-node-tab]')) {
      const selected = button === chosen;
      button.setAttribute('aria-selected', String(selected));
      button.tabIndex = selected ? 0 : -1;
    }
    for (const panel of front.querySelectorAll('[data-node-panel]')) panel.hidden = panel.dataset.nodePanel !== tab;
    selectedNodeTab = tab;
    if (finishCheck?.rows) paintFinishLinks(finishCheck.rows);
    const summary = front.querySelector('.stack-body');
    summary.scrollTop = 0;
    if (narrow.matches) $('L1').scrollTop = 0;
    if (wasFile && narrow.matches) $('phone-back').textContent = fromHeads ? '← Heads' : '← Timeline';
    cameraURL(nodeURL(), false);
    globalThis.WorldNodeTrajectory?.mount(front, F, getJSON, tab === 'trajectory');
    loadEpisodeFiles(front);
  }

  function renderBrane(animate) {
    dismissFilePreview();
    if (!$('bstage')) {
      $('br').innerHTML = `<div class="bstage" id="bstage"><article class="front" id="front"></article></div><nav class="rail" id="rail" aria-label="Brane timeline"></nav>`;
    }
    const front = $('front');
    if (selectedFile && (selectedFile.node !== F.id || selectedFile.brane !== F.brane)) closeFile(false);
    const activeFileControl = selectedFile?.panel.contains(document.activeElement) ? document.activeElement : null;
    const detailKey = el => el.dataset.detail || el.querySelector('summary')?.textContent;
    const sameNode = front.dataset.front === F.id;
    const openDetails = sameNode ? new Set([...front.querySelectorAll('details[open]')].map(detailKey)) : new Set();
    const focusedDetail = sameNode && document.activeElement?.matches('#front summary') ? detailKey(document.activeElement.parentElement) : null;
    const activeFileButton = sameNode && document.activeElement?.closest('#front [data-file-path]');
    const focusedFile = activeFileButton ? fileSelection(activeFileButton) : null;
    const readingPosition = sameNode ? front.querySelector('.stack-body')?.scrollTop : 0;
    if (episodeFiles && (episodeFiles.id !== F.id || episodeFiles.brane !== F.brane)) resetEpisodeFiles();
    const savedFilters = fileFilters(sameNode && front.querySelector('[data-node-panel=files]') || null);
    const focusedTab = sameNode && document.activeElement?.dataset.nodeTab;
    const oldFilter = sameNode && document.activeElement?.matches?.('[data-file-filter], [data-file-kind]') ? document.activeElement : null;
    front.innerHTML = frontHTML(S, T, selectedNodeTab, episodeFiles?.data);
    selectedNodeTab = front.querySelector('[data-node-tab][aria-selected=true]').dataset.nodeTab;
    restoreFilters(front.querySelector('[data-node-panel=files]'), savedFilters);
    if (focusedTab) front.querySelector(`[data-node-tab="${focusedTab}"]`)?.focus({preventScroll:true});
    if (oldFilter) front.querySelector(oldFilter.matches('[data-file-kind]') ? '[data-file-kind]' : '[data-file-filter]')?.focus({preventScroll:true});
    for (const detail of front.querySelectorAll('details')) {
      const key = detailKey(detail);
      if (openDetails.has(key)) detail.open = true;
      if (key === focusedDetail) detail.querySelector('summary')?.focus({preventScroll:true});
    }
    globalThis.WorldNodeTrajectory?.mount(front, F, getJSON, selectedNodeTab === 'trajectory');
    loadEpisodeFiles(front);
    if (readingPosition) front.querySelector('.stack-body').scrollTop = readingPosition;
    if (focusedFile && !selectedFile) focusFileSource(fileButton(front, focusedFile) || front.querySelector('.stack-body'));
    if (selectedFile) {
      if (fileButton(front, selectedFile) || (selectedFile.tab === 'files' && front.querySelector('[data-node-panel=files]'))) {
        front.querySelector('.stack-body').hidden = true;
        front.append(selectedFile.panel);
        activeFileControl?.focus({preventScroll:true});
      } else closeFile();
    }
    const fork = document.createElement('button');
    fork.type = 'button'; fork.className = 'state-branch'; fork.textContent = 'Branch…';
    fork.setAttribute('aria-label', 'Draft a branch from this state');
    fork.setAttribute('aria-haspopup', 'dialog');
    fork.addEventListener('click', () => drafts.edit(N[F.id]));
    if (N[F.id]) $('front').querySelector('.fbar')?.append(fork);
    $('front').dataset.front = F.id;
    showFinishRequests();
    $('front').setAttribute('aria-label', cleanTitle(F.title) || F.label);
    const pos = railIdx();
    if (!$('phone-step-nav')) {
      $('br').insertAdjacentHTML('afterbegin', '<nav id="phone-step-nav" class="phone-view" aria-label="State navigation"><button type="button" id="phone-back">← Timeline</button><button type="button" id="phone-earlier">Earlier</button><button type="button" id="phone-later">Later</button></nav>');
    }
    $('phone-back').textContent = selectedFile ? (selectedFile.tab === 'files' ? '← Files' : '← Summary') : fromHeads ? '← Heads' : '← Timeline';
    $('phone-earlier').disabled = pos === 0;
    $('phone-later').disabled = pos >= RAIL.length - 1;
    placeBars(animate);
    $('rail').innerHTML = railHTML(S);
  }

  // ---------- Header ----------
  const addrOf = () => F.kind === 'open' ? String(F.id).replace(/:episode:/, ':') : `${String(F.brane).replace(/^([a-z][a-z0-9-]*):brane:/, '$1:')}@${railIdx() + 1}`;
  function renderCrumbs() {
    $('crumbs').innerHTML = crumbsHTML(st.d === 1 && F ? addrOf() : null);
    $('addr').placeholder = st.d !== 0 ? '' : (T && T.GROUPS[0].id ? 'go to a place — ' + T.GROUPS.map(g => g.id).join(' · ') + ' · ' + T.GROUPS[0].id + ':row-N@1' : 'go to a place — row-189@6 · shared@0ffbb9c') + ' · episode report <episode> · why <target>';
  }
  let msgT;
  function say(t) { const m = $('hmsg'); m.textContent = t; m.style.opacity = '1'; clearTimeout(msgT); msgT = setTimeout(() => { m.style.opacity = '0'; }, 2600); }
  function follow(path) { try { if (location.pathname + location.search !== path) history.replaceState(null, '', path); } catch (_) { /* the view still moves */ } }

  // ---------- Live data ----------
  async function getJSON(url, signal) {
    const r = await fetch(url, { headers: { Accept: 'application/json' }, cache: 'no-store', signal });
    let body = null;
    try { body = await r.json(); } catch (_) { body = null; }
    return { ok: r.ok, status: r.status, body };
  }
  // REQS is GET /surface/requests.json, read with the model and on every refresh (no timer of its own);
  // decidedAll is the decision list's "show all".
  let REQS = null, decidedAll = false, seenIds = null;
  const notifyPermission = () => typeof Notification === 'undefined' ? '' : Notification.permission;
  async function loadRequests() {
    try { REQS = (await getJSON('/surface/requests.json')).body; } catch (_) { REQS = null; }
    showFinishRequests(true);
    const n = notifications(REQS, seenIds, notifyPermission());
    seenIds = n.seen;
    n.notes.forEach(x => { const note = new Notification(x.title, { body: x.body, tag: x.tag }); note.onclick = () => { window.focus(); location.assign(x.href); }; });
  }
  // askNotify asks this browser for notifications, only from the "Notify me" click (a user gesture).
  function askNotify() { Promise.resolve(Notification.requestPermission()).then(repaintNow, repaintNow); }
  const FOLLOW_MODEL = 'world/surface-model/v2', POLL_MS = 15000;
  let followMode = false, shownGen = null, shownAt = '', pollTimer = null, liveStarted = false;
  const poll = { note: '' };
  async function loadModel() {
    let res;
    await loadRequests();
    try { res = await getJSON('/surface/model.json'); } catch (e) { return renderVoid('error', 'the surface model could not be fetched: ' + e.message); }
    const b = res.body || {};
    followMode = b.schema === FOLLOW_MODEL;
    if (!res.ok) return renderVoid('error', b.reason || `the surface model answered HTTP ${res.status}`, b.repos);
    if (b.empty) return renderVoid('empty', b.reason || 'the surface index is empty');
    if (!Array.isArray(b.lanes) || !b.lanes.some(l => l.states && l.states.length)) return renderVoid('empty', 'The surface index holds no states yet.', b.repos);
    setModel(b);
    shownGen = b.generation == null ? null : b.generation; shownAt = hm(new Date());
    st.focus = defaultFocus();
    renderTimeline(); renderCrumbs();
    setView('timeline');
    return true;
  }

  const live = {mode: 'unavailable', reason: 'the log could not be read', count: 0};
  let liveActive = null, liveRetry = null, liveRefresh = null, refreshing = false, queued = false, backoff = 0, liveRun = 0;
  const liveVisible = () => document.visibilityState !== 'hidden';
  function livePaint() {
    const el = document.querySelector('[data-live]');
    if (el) el.textContent = liveStatus(live);
  }
  function livePause(reason, retryAfter) {
    if (live.mode === 'stopped') return;
    if (liveRetry) clearTimeout(liveRetry);
    if (liveActive) liveActive.abort();
    backoff = nextBackoff(retryAfter, backoff);
    Object.assign(live, {mode: liveVisible() ? 'paused' : 'hidden', reason, wait: backoff});
    livePaint();
    if (liveVisible()) liveRetry = setTimeout(() => { liveRetry = null; if (live.count) queueRefresh(); liveFollow(); }, backoff * 1000);
  }
  function queueRefresh() {
    if (refreshing) { queued = true; return; }
    if (liveRefresh) return;
    liveRefresh = setTimeout(refreshLive, 300);
  }
  const refreshLive = () => refresh('live');
  // One live Now card moves between the Timeline sidebar and the Heads header.
  const toolbar = $('head-toolbar'), statusToggle = $('head-status-toggle'), statusDetail = $('head-status-detail');
  let statusPreference;
  try { statusPreference = localStorage.getItem('world.headStatus'); } catch (_) {}
  function statusExpanded() {
    return statusPreference == null ? !narrow.matches : statusPreference === 'expanded';
  }
  function setStatusExpanded(expanded) {
    statusPreference = expanded ? 'expanded' : 'compact';
    try { localStorage.setItem('world.headStatus', statusPreference); } catch (_) {}
    applyStatusSize();
  }
  function applyStatusSize() {
    const expanded = statusExpanded();
    statusToggle.setAttribute('aria-expanded', String(expanded));
    statusToggle.title = expanded ? 'Minimise world status' : 'Expand world status';
    statusDetail.hidden = !expanded;
  }
  statusToggle.addEventListener('click', () => setStatusExpanded(!statusExpanded()));
  const toolbarResize = new ResizeObserver(() => {
    if (!toolbar.hidden && !toolbar.closest('[inert]')) camera.update();
  });
  toolbarResize.observe(toolbar);

  function focusActivity() {
    if (!toolbar.hidden && statusDetail.hidden) setStatusExpanded(true);
    const card = document.querySelector('[data-now]'), needs = card?.querySelector('[data-needs]');
    if (!needs) return false;
    needs.tabIndex = -1;
    needs.focus({preventScroll:true});
    if (narrow.matches && toolbar.hidden) card.scrollIntoView({block:'start'});
    else needs.scrollIntoView({block:'nearest', inline:'nearest'});
    return true;
  }

  function placeNow(inHeads) {
    if (toolbar.hidden === !inHeads) return;
    toolbar.hidden = !inHeads;
    const card = document.querySelector('[data-now]');
    const slot = inHeads ? statusDetail : $('timeline-now');
    if (card && slot) {
      slot.append(card);
      for (const link of card.querySelectorAll('a[id^=confirm-]')) link.hash = inHeads ? 'heads/activity' : 'activity';
    }
    applyStatusSize();
  }
  function repaintNow() {
    if (!T) return;
    const slot = toolbar.hidden ? $('timeline-now') : statusDetail;
    if (!slot) return;
    const focused = slot.contains(document.activeElement) && document.activeElement;
    slot.innerHTML = nowCard();
    if (focused) {
      const action = focused.id ? $(focused.id) : focused.dataset.go &&
        slot.querySelector('[data-go="' + CSS.escape(focused.dataset.go) + '"]');
      action?.focus({preventScroll:true});
    }
    const card = slot.firstElementChild;
    $('head-status-health').textContent = card.querySelector('[data-health]')?.textContent || '';
    const needs = card.querySelector('[data-needs] .nh')?.textContent;
    const running = card.querySelector('[data-running] .nh')?.textContent;
    const active = [needs, lastRunLabel(D.nowEdge?.open?.lastRun), running].filter(label => label && !/ · 0$/.test(label));
    $('head-status-activity').textContent = active.join(' · ');
    applyStatusSize();
  }
  // schedulePoll: follow mode polls the model every 15 s while the page is visible, independently of AG-UI.
  function schedulePoll(delay) {
    clearTimeout(pollTimer); pollTimer = null;
    if (!followMode || !liveVisible()) return;
    pollTimer = setTimeout(() => { pollTimer = null; refresh('poll'); }, delay == null ? POLL_MS : delay);
  }
  // applyModel replaces the render with a coherent model (and stack), keeping pan, zoom, focus and address.
  function defaultFocus() {
    const work = STEPS.length ? STEPS : T.OPEN;
    const newest = work.reduce((a, n) => !a || n.rank > a.rank ? n : a, null);
    return newest?.id || SH.at(-1) || null;
  }
  function applyModel(r, selection) {
    if (!T) {
      setModel(r.model);
      st.focus = defaultFocus();
      renderTimeline(); renderCrumbs(); setView('timeline');
      if (location.hash === '#heads') chooseAngle(true, {navigate:false, instant:true});
      if (!liveStarted) startLive();
      return;
    }
    const active = document.activeElement, activeId = active && active.id;
    const activeKey = active && ['data-id', 'data-state', 'data-go', 'data-rail', 'data-lane'].find(k => active.hasAttribute(k));
    const activeValue = activeKey && active.getAttribute(activeKey);
    const focus = st.focus, u = view.u, L = view.L, pinned = atNow(), depth = st.d;
    const railTop = $('rail') ? $('rail').scrollTop : 0, address = $('addr').value;
    setModel(r.model);
    st.focus = N[focus] ? focus : T.aliasOf[focus] || defaultFocus();
    view.u = u; view.L = L;
    renderTimeline(true);
    if (pinned) goNow(false);
    else applyView(true);
    if (depth === 1 && selection && selection.seq === openSeq) {
      if (!r.stack) { say(r.missing || 'the selection is no longer on the surface'); out(); return; }
      setStack(r.stack); renderBrane(false); $('rail').scrollTop = railTop;
      if (N[F.id]) st.focus = F.id;
      cameraURL(nodeURL(), false); renderCrumbs();
    }
    st.d = depth;
    setView(depth === 1 ? 'brane' : camera.active ? 'heads' : 'timeline');
    $('addr').value = address;
    if (active && active.isConnected) active.focus({preventScroll: true});
    else {
      const replacement = activeId ? $(activeId) : activeKey ? Array.from(document.querySelectorAll('[' + activeKey + ']')).find(e => e.getAttribute(activeKey) === activeValue && e.tabIndex >= 0) : null;
      if (replacement) replacement.focus({preventScroll: true});
    }
  }
  // refresh is N2c's refresh on a committed log entry ("live") and follow mode's model poll ("poll"): fetch
  // everything before touching the old render; an unchanged generation keeps the stack.
  async function refresh(trigger) {
    if (trigger === 'live') liveRefresh = null;
    if (refreshing) { if (trigger === 'live') queued = true; return; }
    refreshing = true;
    const count = live.count;
    try {
      const selection = T && st.d === 1 && F ? {brane: F.brane, at: F.id, seq: openSeq} : null;
      await loadRequests();
      const r = await fetchCoherent(getJSON, selection, T ? shownGen : null);
      if (r.kind === 'fail') {
        if (!followMode) throw new Error('refresh failed (HTTP ' + r.status + ')');
        const reason = (r.body && r.body.reason) || 'HTTP ' + r.status;
        if (!T) renderVoid('error', reason, r.body && r.body.repos);
        poll.note = 'Model refresh failed: ' + reason + (T ? ' · showing the view from ' + shownAt : '');
      } else if (r.kind === 'keep') {
        poll.note = 'The graph changed during the refresh · showing the view from ' + shownAt;
      } else {
        if (r.kind === 'same') applyModel(r, null);
        else { applyModel(r, selection); indexP = null; }
        if (r.model) { shownGen = r.model.generation == null ? null : r.model.generation; shownAt = hm(new Date()); }
        poll.note = '';
        const cursor = r.model && r.model.nowEdge && r.model.nowEdge.live;
        if (cursor && live.cursor != null) live.cursor = Math.max(live.cursor, cursor.lastIndex);
        if (trigger === 'live') { live.count = Math.max(0, live.count - count); live.checked = hm(new Date()); }
      }
    } catch (e) {
      queued = false;
      if (followMode) poll.note = 'Model refresh failed: the daemon could not be reached · showing the view from ' + shownAt;
      else livePause(/^refresh failed/.test(e.message) ? e.message : 'refresh failed (HTTP 0)');
    } finally {
      refreshing = false;
      if (queued) { queued = false; queueRefresh(); }
      if (camera.active) camera.update();
      repaintNow();
      livePaint();
      if (followMode) schedulePoll();
    }
  }
  async function liveFollow() {
    if (!liveVisible() || liveActive || liveRetry || live.mode === 'stopped' || live.mode === 'unavailable') return;
    const controller = new AbortController(); liveActive = controller;
    let terminal = false, failure = null, retryAfter;
    try {
      const response = await fetch('/agui/', {method: 'POST', headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(liveRequest(live.cursor, liveRun++)), signal: controller.signal});
      if (response.status === 404) { live.mode = 'stopped'; live.reason = 'unknown cursor'; livePaint(); return; }
      if (!response.ok) {
        retryAfter = response.headers.get('Retry-After');
        throw new Error(response.status === 503 ? 'stream limit' : 'stream HTTP ' + response.status);
      }
      live.mode = 'live'; livePaint();
      const reader = response.body.getReader(), decoder = new TextDecoder();
      let tail = '';
      while (!controller.signal.aborted) {
        const chunk = await reader.read();
        if (controller.signal.aborted) break;
        const parsed = splitFrames(tail, decoder.decode(chunk.value || new Uint8Array(), {stream: !chunk.done}));
        tail = parsed.tail;
        for (const event of parsed.events) {
          if (event.type === 'STATE_DELTA') {
            const next = applyDelta({lastIndex: live.cursor}, event.delta);
            live.cursor = Math.max(live.cursor, next.lastIndex);
          } else if (event.type === 'CUSTOM' && event.name === 'world.entry.committed') {
            live.count++; queueRefresh();
          } else if (event.type === 'RUN_FINISHED') {
            live.cursor = Math.max(live.cursor, event.result.lastIndex);
            live.checked = hm(new Date()); terminal = true;
          } else if (event.type === 'RUN_ERROR') throw new Error(event.code || 'stream error');
        }
        livePaint();
        if (terminal || chunk.done) break;
      }
      if (!terminal && !controller.signal.aborted) throw new Error('connection ended');
    } catch (e) { if (!controller.signal.aborted) failure = e.message; }
    finally { controller.abort(); liveActive = null; }
    if (!liveVisible()) return;
    if (failure) livePause(failure, retryAfter);
    else if (terminal) { backoff = 0; liveFollow(); }
    else if (live.mode === 'live') liveFollow();
  }
  function startLive() {
    liveStarted = true;
    const ne = D.nowEdge || {};
    if (!ne.live) {
      Object.assign(live, {mode: 'unavailable', reason: ne.liveUnavailable || 'the log could not be read'});
    } else {
      Object.assign(live, {mode: liveVisible() ? 'live' : 'hidden', cursor: ne.live.lastIndex, checked: hm(new Date())});
      liveFollow();
    }
    livePaint();
  }
  document.addEventListener('visibilitychange', () => {
    if (live.mode === 'unavailable' || live.mode === 'stopped') return;
    if (!liveVisible()) {
      if (liveActive) liveActive.abort();
      if (liveRetry) clearTimeout(liveRetry);
      liveRetry = null; live.mode = 'hidden'; livePaint();
    } else { live.mode = 'live'; livePaint(); if (live.count) queueRefresh(); liveFollow(); }
  });

  // Follow mode's poll is independent of the stream: hidden pauses it, visible fetches at once.
  document.addEventListener('visibilitychange', () => { if (liveVisible()) schedulePoll(0); else { clearTimeout(pollTimer); pollTimer = null; } });

  let indexP = null;
  function loadIndex() {
    if (!indexP) {
      indexP = getJSON('/surface/index.json').then(r => { if (!r.ok || !r.body) throw new Error('HTTP ' + r.status); return r.body; });
      indexP.catch(() => { indexP = null; });
    }
    return indexP;
  }

  // ---------- Navigation between timeline and brane ----------
  function setOrigin(x, y) {
    const o = `${x * scale}px ${y * scale}px`;
    $('L0').style.transformOrigin = o; $('L1').style.transformOrigin = o;
  }
  function openCluster(cluster) {
    if (view.u > 0.001) return zoomTo(Math.ceil(view.u - 0.01) - 1, PL + cluster.x - view.L);
    const dialog = document.createElement('dialog'), heading = document.createElement('h2');
    heading.id = 'timeline-cluster-title';
    heading.textContent = 'States at this position';
    dialog.setAttribute('aria-labelledby', heading.id);
    dialog.append(heading);
    const list = document.createElement('ul');
    for (const id of cluster.ids) {
      const item = document.createElement('li'), button = document.createElement('button'), state = N[id];
      button.textContent = nameOf(state.lane) + ' · ' + (cleanTitle(state.title) || tagOf(state)) + ' · ' + exactOf(state.at);
      button.addEventListener('click', () => { dialog.close(); open(id); });
      item.append(button);
      list.append(item);
    }
    const close = document.createElement('button');
    close.textContent = 'Close';
    close.addEventListener('click', () => dialog.close());
    dialog.append(list, close);
    dialog.addEventListener('close', () => dialog.remove(), {once: true});
    $('stage').append(dialog);
    dialog.showModal();
  }
  function open(id) {
    if (!T) return;
    if (isCluster(id)) { if (DISP[id]) openCluster(DISP[id]); return; }
    const n = N[id];
    if (!n) { say('That state is not on the timeline'); return; }
    if (n.draft) return drafts.edit(N[n.parent], id);
    return openState(n.lane, id);
  }
  // openState fetches the state's stack and zooms into it; a newer open supersedes an older one.
  let openSeq = 0;
  async function openState(brane, at, navigate = true, tab) {
    if (N[at]?.draft) return drafts.edit(N[N[at].parent], at);
    const seq = ++openSeq;
    let res;
    try { res = await getJSON(stackPath(brane, at)); } catch (e) { if (seq === openSeq) say('the stack could not be fetched: ' + e.message); return false; }
    if (seq !== openSeq) return false;
    const b = res.body || {};
    if (!res.ok || !b.front) { say(b.reason || `no such state (HTTP ${res.status})`); if (st.d === 0 && !camera.active) follow(appPath()); return false; }
    const was = st.d === 1;
    $('stage').classList.remove('head-return');
    const face = camera.active && faceFor(b.front.id)?.getBoundingClientRect();
    if (!was) fromHeads = camera.heads && camera.active || tabLocation().heads;
    if (tab != null || !was) selectedNodeTab = tab || 'overview';
    setStack(b);
    if (!was && T) { const sx = screenX(F.id), d = DISP[REP[F.id]]; setOrigin(sx != null ? Math.max(PL, Math.min(PR, sx)) : PL, d ? d.y : LY[T.GROUPS[0].shared]); }
    st.d = 1; st.focus = F.id; $('tip').hidden = true;
    renderBrane(was);
    if (narrow.matches) $('L1').scrollTop = 0;
    $('root').classList.add('zoomed'); setView('brane');
    $('L0').inert = true;
    $('L1').inert = false;
    markAngle(null);
    if (camera.active) {
      $('head-space').inert = true;
      $('stage').classList.add('head-approach');
      const flight = animateStack(face, false);
      const finish = () => {
        if (st.d !== 1) return;
        camera.hide();
        $('stage').classList.remove('head-approach');
      };
      if (flight) flight.finished.then(finish, () => {}); else finish();
    }
    renderCrumbs();
    if (fromHeads) {
      if (navigate) cameraURL(nodeURL(), !was);
    } else cameraURL(nodeURL(), false);
    $('root').focus({ preventScroll: true });
    return true;
  }
  function out(focusId) {
    resetFinishLinks();
    resetEpisodeFiles();
    selectedNodeTab = 'overview';
    globalThis.WorldNodeTrajectory?.clear();
    if (st.d === 0) {
      if (camera.active) chooseAngle(false, {keyboard:headRibbon.contains(document.activeElement)});
      return;
    }
    if (fromHeads) { chooseAngle(true, {keyboard:true}); return; }
    dismissFilePreview();
    closeFile(false);
    $('L0').inert = false; $('L1').inert = true; markAngle('timeline');
    openSeq++;
    const id = focusId || (F && F.id);
    if (T && id && POS[id]) {
      st.focus = timelineAnchor = id;
      const x = screenX(id); if (x == null || x < PL + 40 || x > PR - 40) panTo(view.L + (x - (PL + PW / 2)), false);
      applyFocus(); setOrigin(screenX(id), DISP[REP[id]].y);
    }
    st.d = 0; $('root').classList.remove('zoomed'); $('tip').hidden = true;
    setView(T ? 'timeline' : $('stage').dataset.view);
    renderCrumbs(); follow(appPath()); $('root').focus({ preventScroll: true });
  }
  function openRail(k) {
    k = Math.max(0, Math.min(RAIL.length - 1, k));
    if (!F || !RAIL[k] || RAIL[k].id === F.id) return;
    openState(F.brane, RAIL[k].id);
  }
  // go is the bar's Enter: empty is home (no fetch); anything else asks the verb route first.
  function go(raw) {
    return goWith(raw, { out, address: goAddress,
      ask: q => askVerb(q, { fetch: (u, o) => fetch(u, o), AbortController, setTimeout: (f, ms) => setTimeout(f, ms), clearTimeout: t => clearTimeout(t) }),
      act: { assign: h => {
        const suffix = !/^\/surface\/confirm\/vr-[0-9a-f]{16}$/.test(h) ? ''
          : st.d === 1 && F ? finishReturnFragment(F)
          : camera.active ? (location.hash === '#heads/activity' ? '#heads/activity' : '#heads')
          : location.hash === '#activity' ? '#activity' : '';
        location.assign(h + suffix);
      }, panel: showVerbPanel, say, clear: () => { if ($('addr').value === raw) $('addr').value = ''; } } });
  }
  function showVerbPanel(p) {
    const old = $('verb'), el = verbPanel(document, p);
    if (old) old.replaceWith(el); else $('stage').appendChild(el);
    const b = el.querySelector('button');
    if (b) b.addEventListener('click', () => { try { navigator.clipboard.writeText(b.dataset.copy).then(() => say('copied'), () => say('copy it from the panel')); } catch (_) { say('copy it from the panel'); } });
  }
  function goAddress(raw, current) { return addressGo(raw, current, { parse: SEAMS.parseAddress, loadIndex, out, openState, say }); }

  // ---------- Popovers ----------
  const KEYROWS = [['← →', 'back · forward a page', '—'], ['⇧← ⇧→', 'previous · next state', '—'], ['↑ ↓', 'lane above · below', 'back in time · toward now'], ['+ −', 'zoom in · out', '—'], ['End', 'back to now', 'most recent'], ['Home', 'the oldest state', 'earliest'], ['Enter', 'open brane', '—'], ['Esc', '—', 'close the file view or return to the timeline'], ['N', 'what needs you', 'what needs you'], ['/', 'address bar', 'address bar'], ['Enter', 'in the bar: a place, episode report E, why T, approve ID', 'the same'], [',', 'appearance', 'appearance'], ['?', 'these keys', 'these keys']];
  function renderKeys() {
    if (camera.heads && camera.active && st.d === 0) {
      const rows = [['← →', 'Bring the neighboring lane’s latest entry forward'], ['↑ / +', 'Move into history'],
        ['↓ / −', 'Move toward the present'],
        ['Enter', 'Open the focused state'], ['Esc', 'Return to Timeline'], ['/', 'Address bar'],
        ['N', 'What needs you'], [',', 'Appearance'], ['?', 'These keys']];
      $('keys').innerHTML = '<div class="kgrid heads-keys"><strong>Heads</strong>' +
        rows.map(([key, action]) => `<span><span class="kbd">${key}</span></span><span>${action}</span>`).join('') +
        '</div><p>Hold an arrow to keep moving. Drag sideways to switch threads or vertically to travel in time. Double-click the grid to hop there.</p>';
      return;
    }
    $('keys').innerHTML = `<div class="kgrid"><span></span><span style="font-size:12px;font-weight:700;color:${st.d === 0 ? 'var(--text)' : 'var(--faint)'}">Timeline</span><span style="font-size:12px;font-weight:700;color:${st.d === 1 ? 'var(--text)' : 'var(--faint)'}">Brane</span>` +
      KEYROWS.map(([k, a, b]) => `<span><span class="kbd">${k}</span></span><span>${a}</span><span>${b}</span>`).join('') +
      `</div><div style="margin-top:10px;font-size:12px;color:var(--sub)">Drag or scroll sideways to pan. ⌘-scroll or pinch to zoom around the pointer.</div>`;
  }
  // Appearance is stored under M4a's key (world.appearance) through its seams; failures are swallowed.
  let storage = null;
  try { storage = window.localStorage; } catch (_) { storage = null; }
  const media = q => { try { return typeof matchMedia === 'function' ? matchMedia(q) : null; } catch (_) { return null; } };
  const darkQ = media('(prefers-color-scheme: dark)'), reducedQ = media('(prefers-reduced-motion: reduce)');
  function readLook() {
    if (typeof SEAMS.loadLook === 'function') return SEAMS.loadLook(storage);
    try { return JSON.parse(storage.getItem('world.appearance') || '{}') || {}; } catch (_) { return {}; }
  }
  let stored = readLook();
  const look = lookOf(stored, !!(darkQ && darkQ.matches), !!(reducedQ && reducedQ.matches));
  function applyLook() {
    document.documentElement.dataset.theme = look.theme;
    $('root').classList.toggle('reduced', look.motion === 'reduced');
  }
  function saveLook() {
    stored = Object.assign({}, stored, look);
    if (typeof SEAMS.saveLook === 'function') SEAMS.saveLook(storage, stored);
    else try { storage.setItem('world.appearance', JSON.stringify(stored)); } catch (_) { /* this browser keeps no look */ }
  }
  function renderLook() {
    const opt = (key, v, label) => { const on = look[key] === v; return `<button type="button" data-look="${key}:${v}" aria-pressed="${on}" style="font:inherit;flex:1;min-height:44px;border-radius:var(--r);cursor:pointer;border:1px solid ${on ? 'var(--you)' : 'var(--ctl)'};background:${on ? 'var(--you)' : 'var(--btn)'};color:${on ? 'var(--onYou)' : 'var(--text)'};font-weight:${on ? 700 : 500};font-size:13px">${label}</button>`; };
    $('look').innerHTML = `<div style="font-size:15px;font-weight:700">Appearance</div>
      <div style="display:flex;flex-direction:column;gap:6px"><div style="font-size:12.5px;color:var(--sub)">Theme</div><div style="display:flex;gap:6px">${opt('theme', 'paper', 'Paper')}${opt('theme', 'slate', 'Slate')}</div></div>
      <div style="display:flex;flex-direction:column;gap:6px"><div style="font-size:12.5px;color:var(--sub)">Motion</div><div style="display:flex;gap:6px">${opt('motion', 'full', 'Full')}${opt('motion', 'reduced', 'Reduced')}</div></div>
      <div style="display:flex;align-items:center;gap:10px;border-top:1px solid var(--line);padding-top:12px;font-size:12px;color:var(--sub)">Paper and Slate follow the Sunholo brand guide. Your choices stay on this device.</div>`;
  }
  function toggle(which) {
    const other = which === 'keys' ? 'look' : 'keys';
    $(other).hidden = true;
    $(which).hidden = !$(which).hidden;
    if (!$(which).hidden) (which === 'keys' ? renderKeys : renderLook)();
    $('keys-btn').setAttribute('aria-expanded', String(!$('keys').hidden));
    $('look-btn').setAttribute('aria-expanded', String(!$('look').hidden));
  }

  // ---------- Hover cards: the main way a human reads the timeline ----------
  function nodeCard(id) {
    if (isCluster(id)) {
      const d = DISP[id]; if (!d) return '';
      const ns = d.ids.map(k => N[k]).sort((a, b) => b.rank - a.rank), shared = ns.every(n => n.kind === 'shared');
      const t0 = ns[ns.length - 1].at, t1 = ns[0].at;
      const span = dayOf(t0) === dayOf(t1) ? `${dayOf(t0)} · ${hm(t0)}–${hm(t1)}` : `${exactOf(t0)} → ${exactOf(t1)}`;
      return `<div class="hc"><div class="hk">${ns.length} ${shared ? 'shared records' : 'work records'} · ${esc(span)}</div>` +
        ns.slice(0, 6).map(n => `<div class="hli">${glyphOfNode(n, 9)}<span>${esc(shortTitle(n.title, 64) || n.episode || n.label)}</span></div>`).join('') +
        `<div class="hid">${ns.length > 6 ? `+${ns.length - 6} more · ` : ''}click or zoom in to open them up</div></div>`;
    }
    const n = N[id]; if (!n) return '';
    const shared = n.kind === 'shared';
    const fileSummary = mergeFileSummary(n);
    const files = fileSummary ? `<div class="hm">${esc(fileSummary)}</div>` : '';
    const ids = [n.episode, n.pr ? '#' + n.pr : '', String(id).replace(/^.*:commit:/, '').slice(0, 7)].filter(Boolean).join(' · ');
    if (shared) {
      const from = (parents[id] || []).map(p => N[p]).find(p => p.kind === 'step');
      return `<div class="hc"><div class="hk">${glyphOfNode(n, 11)}${T.recordLabel(n)} · ${esc(nameOf(n.lane))} · ${esc(n.rel)}</div>
        ${n.title ? `<div class="ht">${esc(cleanTitle(n.title))}</div>` : ''}
        ${files}
        ${from && !isSharedKey(from.lane) ? `<div class="hm">from ${esc(nameOf(from.lane))}${from.episode ? ' · ' + esc(from.episode) : ''}</div>` : ''}
        <div class="hw">${esc(exactOf(n.at))}</div><div class="hid">${esc(ids)}</div></div>`;
    }
    if (n.kind === 'open') {
      const o = n.open || {}, base = o.baseOutsideWindow ? `base ${o.base || 'unrecorded'} is outside the retained window` : o.base ? 'built on ' + o.base : '';
      return `<div class="hc"><div class="hk">${glyphOfNode(n, 9)}${T.recordLabel(n)} · ${esc(nameOf(n.lane))} · ${esc(o.label)}</div>
        <div class="ht">${esc(n.title || n.episode)}</div>
        ${o.lastRun ? `<div class="hm">${esc(lastRunDetail(o.lastRun))}</div>` : ''}
        ${n.model ? `<div class="hm"><b>${esc(agentLine(n))}</b></div>` : ''}
        <div class="hm">${esc([o.branch, base, n.pr ? 'PR #' + n.pr : ''].filter(Boolean).join(' · '))}</div>
        <div class="hw">started ${esc(n.rel)}${o.ended ? ' · ended ' + esc(o.ended) : ''}</div><div class="hid">${esc(n.episode)}</div></div>`;
    }
    const s = summaryOf(n);
    const g = n.graders ? `${esc(n.graders)} graders · ${n.firstRoundPass ? 'first round ✓' : `${n.rounds} round${n.rounds === 1 ? '' : 's'}`}` : '';
    return `<div class="hc"><div class="hk">${glyphOfNode(n, 9)}${T.recordLabel(n)} · ${esc(nameOf(n.lane))}</div>
      <div class="ht">${esc(cleanTitle(n.title) || n.episode || n.label)}</div>
      ${s ? `<div class="hs">${esc(s)}</div>` : ''}
      ${files}
      ${n.model ? `<div class="hm"><b>${esc(agentLine(n))}</b>${g ? ' · ' + g : ''}${n.minutes != null ? ` · ${n.minutes} agent-min` : ''}</div>` : ''}
      ${n.badge ? `<div class="hm"><span class="nb inl">${n.badge}</span> ${n.badge === 1 ? 'decision' : 'decisions'} of yours</div>` : ''}
      ${n.sessionLive ? '<div class="hm">its episode\'s session is still live</div>' : ''}
      <div class="hw">${esc(n.rel)} · ${esc(exactOf(n.at))}</div><div class="hid">${esc(ids)}</div></div>`;
  }
  function tipHTML(spec) {
    const [kind, ...rest] = spec.split(':'), v = rest.join(':');
    if (kind === 'node') return T ? nodeCard(v) : '';
    if (kind === 'reads') {
      const r = READS.get(v); if (!r) return '';
      const n = r.files.reduce((a, f) => a + f.count, 0), s = N[r.fromId], title = (s && s.title) || r.title;
      const head = title ? `<div class="th">${esc(cleanTitle(title))}</div><div style="color:var(--sub);margin-bottom:6px">${esc(r.label)}${s ? ' · ' + esc(s.rel) : ''} · read ${r.files.length} file${r.files.length === 1 ? '' : 's'}, ${n} time${n === 1 ? '' : 's'}</div>`
        : `<div class="th">${esc(r.label)} <span style="font-weight:500;color:var(--sub)">· ${r.files.length} file${r.files.length === 1 ? '' : 's'}, ${n} read${n === 1 ? '' : 's'}</span></div>`;
      return head + r.files.map(f => `<div class="tf"><span>${esc(f.path)}</span><span class="c">×${f.count}</span></div>`).join('');
    }
    if (kind === 'goal') return LANES[v] ? `<div class="hc"><div class="hk">${esc(laneTitle(v))} · why it exists</div><div class="ht" style="font-weight:500">${esc(LANES[v].goal)}</div></div>` : '';
    if (kind === 'ms') {
      const ms = LANES[v] && LANES[v].milestones; if (!ms) return '';
      const rest2 = ms.names.filter(x => !ms.doneNames.includes(x) && x !== ms.inReview && x !== ms.next);
      return `<div class="hc"><div class="hk">${esc(nameOf(v))} · milestones</div><div class="hm"><b>Done</b> ${esc(ms.doneNames.join(', '))}</div><div class="hm"><b>In review</b> ${esc(ms.inReview)}</div><div class="hm"><b>Next</b> ${esc(ms.next)}${rest2.length ? ` · then ${esc(rest2.join(', '))}` : ''}</div></div>`;
    }
    if (kind === 'k') return `<div>${esc(v)}</div>`;
    if (kind === 'path') return `<div class="tf"><span>${esc(v)}</span></div>`;
    if (kind === 'src') return `<div>Source: ${esc(v)}</div>`;
    return '';
  }
  let dragging = false;
  document.addEventListener('pointerover', e => {
    const t = e.target.closest && e.target.closest('[data-tip]'), tip = $('tip');
    if (dragging || !t || !t.dataset.tip) { tip.hidden = true; return; }
    const html = tipHTML(t.dataset.tip); if (!html) { tip.hidden = true; return; }
    tip.innerHTML = html; tip.hidden = false;
    tip.classList.toggle('wide', /^(node|goal|ms):/.test(t.dataset.tip));
    const sr = $('stage').getBoundingClientRect(), r = t.getBoundingClientRect();
    const isNode = t.classList.contains('node');
    let x = isNode ? r.right - sr.left + 8 : r.left - sr.left, y = isNode ? r.top - sr.top - 6 : r.bottom - sr.top + 6;
    const tw = tip.offsetWidth, th = tip.offsetHeight;
    if (x + tw > sr.width - 8) x = isNode ? r.left - sr.left - tw - 8 : sr.width - 8 - tw;
    if (y + th > sr.height - 8) y = isNode ? sr.height - 8 - th : r.top - sr.top - th - 6;
    tip.style.left = Math.max(8, x) + 'px'; tip.style.top = Math.max(8, y) + 'px';
  });

  // ---------- Events ----------
  function inspectTimelineNode(e) {
    const n = e.target.closest('#tl .node');
    if (n && T && st.d === 0 && !camera.active && !dragging) chooseTimelineNode(n.dataset.id);
  }
  $('stage').addEventListener('pointerover', inspectTimelineNode);
  $('stage').addEventListener('focusin', inspectTimelineNode);
  $('stage').addEventListener('pointerover', event => {
    if (event.pointerType !== 'touch') startFilePreview(event.target.closest('[data-file-path]'));
  });
  $('stage').addEventListener('pointerout', event => {
    if (filePreview?.button.contains(event.target)) leaveFilePreview(event.relatedTarget);
  });
  $('stage').addEventListener('focusin', event => startFilePreview(event.target.closest('[data-file-path]')));
  function filterNodeFiles(event) {
    const filter = event.target.closest('[data-file-filter], [data-file-kind]');
    if (filter) applyFileFilters(filter.closest('.node-files-tab'));
  }
  $('stage').addEventListener('input', filterNodeFiles);
  $('stage').addEventListener('change', filterNodeFiles);
  $('stage').addEventListener('focusout', event => {
    const button = filePreview?.button;
    if (button?.contains(event.target) && !button.contains(event.relatedTarget)) dismissFilePreview();
  });
  $('stage').addEventListener('scroll', placeFilePreview, true);
  window.addEventListener('resize', placeFilePreview);
  let suppressClick = false;
  $('stage').addEventListener('click', e => {
    if (suppressClick) { suppressClick = false; e.stopPropagation(); return; }
    if (e.target.closest('#phone-back')) return selectedFile ? closeFile() : out();
    if (e.target.closest('#phone-earlier')) return openRail(railIdx() - 1);
    if (e.target.closest('#phone-later')) return openRail(railIdx() + 1);
    if (e.target.closest('#phone-activity')) {
      cameraURL('/surface#activity', false);
      return focusActivity();
    }
    if (e.target.closest('#activity-back')) {
      cameraURL('/surface', false);
      $('L0').scrollTop = 0;
      return $('phone-activity').focus({preventScroll:true});
    }
    const n = e.target.closest('#tl .node');
    if (n) return open(n.dataset.id);
    const qb = e.target.closest('[data-quiet]');
    if (qb) { const l = [...document.querySelectorAll('[data-quiet-list]')].find(x => x.dataset.quietList === qb.dataset.quiet); if (l) { l.hidden = !l.hidden; qb.setAttribute('aria-expanded', String(!l.hidden)); } return; }
    if (e.target.closest('#decided-all')) { decidedAll = !decidedAll; return repaintNow(); }
    if (e.target.closest('#notify-me')) return askNotify();
    const lab = e.target.closest('.lhead');
    if (lab) { const key = lab.dataset.lane, s = lastState(key); return s ? openState(key, s.id) : say(`${nameOf(key)} has no states yet`); }
    if (e.target.closest('#now-b')) return goNow();
    if (e.target.closest('#pg-back')) return panTo(view.L - PW * 0.8);
    if (e.target.closest('#pg-fwd')) return panTo(view.L + PW * 0.8);
    if (e.target.closest('#z-in')) return zoomStep(-1);
    if (e.target.closest('#z-out')) return zoomStep(1);
    const rr = e.target.closest('.rrow'); if (rr) return openRail(+rr.dataset.k);
    const bar = e.target.closest('.strip');
    if (bar) return openState(bar.dataset.lane, bar.dataset.id);
    if (e.target.closest('[data-file-summary]')) return closeFile();
    const file = e.target.closest('[data-file-path]');
    if (file) return openFile(file);
    const ch = e.target.closest('[data-go]');
    if (ch) {
      const id = ch.dataset.go, k = RAIL.findIndex(r => r.id === id);
      if (st.d === 1 && F && id === F.id) return;
      if (ch.closest('#head-status-detail') && N[id]) return open(id);
      const to = T && T.aliasOf[id] ? T.aliasOf[id] : id;
      return k >= 0 ? openRail(k) : N[to] ? openState(N[to].lane, to) : say(`${ch.textContent.trim() || id.slice(0, 7)} is outside the timeline`);
    }
  }, true);

  // Drag to pan (the plot), drag/click the overview window.
  let drag = null;
  $('stage').addEventListener('pointerdown', e => {
    if (!T || st.d !== 0 || camera.active || e.button !== 0) return;
    const plot = e.target.closest('#plot'), mini = e.target.closest('#mini');
    if (!plot && !mini) return;
    cancelAnimationFrame(anim);
    if (mini) {
      const r = $('mini').getBoundingClientRect(), mx = (e.clientX - r.left) / scale, mw = $('mwin');
      const wl = parseFloat(mw.style.left), ww = parseFloat(mw.style.width);
      if (mx < wl || mx > wl + ww) {   // click outside the window: jump there, then keep dragging
        const f = fOfMX(mx); panTo(posF(f) - PW / 2);
        drag = { kind: 'mini', x0: e.clientX, f0: fOfMX(mx), grab: 0.5, moved: true, jumped: true };
      } else drag = { kind: 'mini', x0: e.clientX, f0: null, grabL: wl, moved: false };
    } else drag = { kind: 'plot', x0: e.clientX, L0: view.L, moved: false };
    drag.pid = e.pointerId;
    if (mini) $('stage').setPointerCapture(e.pointerId);
  });
  $('stage').addEventListener('pointermove', e => {
    if (!drag) return;
    const dx = (e.clientX - drag.x0) / scale;
    if (!drag.moved && Math.abs(dx) < 4) return;
    if (!drag.moved) { try { $('stage').setPointerCapture(drag.pid); } catch (_) {} }
    drag.moved = true; dragging = true; $('tip').hidden = true;
    if (drag.kind === 'plot') { view.L = drag.L0 - dx; applyView(); $('plot').classList.add('grab'); }
    else if (drag.jumped) { const f = fOfMX(MX(drag.f0) + dx); cancelAnimationFrame(anim); view.L = posF(f) - PW / 2; applyView(); }
    else { const f = fOfMX(drag.grabL + dx); cancelAnimationFrame(anim); view.L = posF(f); applyView(); }
  });
  const endDrag = () => { if (drag && drag.moved && drag.kind === 'plot') suppressClick = true; drag = null; dragging = false; const p = $('plot'); if (p) p.classList.remove('grab'); setTimeout(() => { suppressClick = false; }, 0); };
  $('stage').addEventListener('pointerup', endDrag);
  $('stage').addEventListener('pointercancel', endDrag);

  // Scroll / trackpad: sideways (or plain wheel) pans; ⌘/Ctrl-wheel and pinch zoom around the pointer.
  let wq = null;
  $('stage').addEventListener('wheel', e => {
    if (!T || camera.active || st.d !== 0 || narrow.matches) return;
    e.preventDefault();
    const r = $('tl').getBoundingClientRect(), ax = Math.max(PL, Math.min(PR, (e.clientX - r.left) / scale));
    const unit = e.deltaMode === 1 ? 16 : 1;
    if (e.ctrlKey || e.metaKey) {
      const du = e.deltaY * unit * (e.ctrlKey && !e.metaKey && Math.abs(e.deltaY) < 50 ? 0.012 : 0.004);
      wq = wq || { du: 0, ax }; wq.du += du; wq.ax = ax;
    } else {
      const d = (Math.abs(e.deltaX) > Math.abs(e.deltaY) ? e.deltaX : e.deltaY) * unit / scale;
      wq = wq || { du: 0, ax }; wq.dL = (wq.dL || 0) + d;
    }
    if (!wq.raf) wq.raf = requestAnimationFrame(() => {
      const q = wq; wq = null;
      cancelAnimationFrame(anim); cancelAnimationFrame(zanim);
      if (q.du) setZoom(view.u + q.du, q.ax);
      if (q.dL) { view.L += q.dL; applyView(); }
    });
  }, { passive: false });

  // The brand is the way home: in a brane it does what Esc does, on the timeline it goes to Now.
  // A modified or middle click keeps the link's own default (/surface in a new tab or window).
  $('home').addEventListener('click', e => {
    if (e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return;
    e.preventDefault();
    if (st.d === 1) out(); else if (camera.active) chooseAngle(false); else if (T) goNow();
  });
  for (const [id, heads] of [['head-angle',true],['timeline-angle',false]]) {
    $(id).addEventListener('click', event => {
      if (event.button || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) return;
      event.preventDefault();
      chooseAngle(heads, {keyboard:event.detail === 0});
    });
  }
  $('br').addEventListener('click', event => {
    const tab = event.target.closest('[data-node-tab]');
    if (tab) selectNodeTab(tab.dataset.nodeTab);
  });
  $('br').addEventListener('keydown', event => {
    const tab = event.target.closest('[data-node-tab]');
    if (!tab) return;
    if (event.key === 'Enter' || event.key === ' ') { event.stopPropagation(); return; }
    const tabs = Array.from(tab.parentElement.querySelectorAll('[data-node-tab]'));
    const index = tabs.indexOf(tab);
    const next = {Home:0, End:tabs.length - 1, ArrowLeft:(index + tabs.length - 1) % tabs.length, ArrowRight:(index + 1) % tabs.length}[event.key];
    if (next == null) return;
    event.stopPropagation();
    event.preventDefault();
    selectNodeTab(tabs[next].dataset.nodeTab);
    tabs[next].focus({preventScroll:true});
  });
  window.addEventListener('popstate', () => {
    const q = parseAppQuery(location.search);
    if (q?.brane) { const place = tabLocation(); fromHeads = place.heads; openState(q.brane, q.at, false, place.tab); }
    else {
      const activity = location.hash === '#activity' || location.hash === '#heads/activity';
      chooseAngle(location.hash === '#heads' || location.hash === '#heads/activity', {navigate:false, instant:activity});
      if (activity) focusActivity();
    }
  });
  $('keys-btn').addEventListener('click', () => toggle('keys'));
  $('look-btn').addEventListener('click', () => toggle('look'));
  $('look').addEventListener('click', e => { const b = e.target.closest('[data-look]'); if (!b) return; const [k, v] = b.dataset.look.split(':'); look[k] = v; saveLook(); applyLook(); renderLook(); });
  $('addr').addEventListener('keydown', e => {
    e.stopPropagation();
    if (e.key === 'Enter') { go(e.target.value); $('root').focus({ preventScroll: true }); }
    if (e.key === 'Escape') { e.target.value = ''; $('root').focus({ preventScroll: true }); }
  });
  const onRow = row => Object.keys(DISP).filter(id => DISP[id].row === row).sort((a, b) => DISP[a].x - DISP[b].x);
  document.addEventListener('keydown', e => {
    if (e.target && (e.target.tagName === 'SELECT' || e.target.tagName === 'INPUT' || e.target.closest?.('dialog, #compose-btn, .state-branch'))) return;
    const k = e.key; let done = true;
    if (e.target.closest?.('.camera-angles') && k === 'Enter') return;
    if (e.target.closest?.('.stack-disclosure summary, .file-path, .file-evidence button') && k === 'Enter') return;
    if (e.target.closest?.('.file-evidence') && ['ArrowUp', 'ArrowDown', 'Home', 'End'].includes(k)) return;
    if (e.target.matches?.('.stack-body, [role=tabpanel]') && ['ArrowUp', 'ArrowDown', 'Home', 'End'].includes(k)) return;
    if (camera.active && st.d === 0 && !['Escape','?',',','/','n','N'].includes(k)) return;
    if (k === 'Enter' && e.target && e.target.id === 'home') return;   // the link's own activation: its click goes home
    if (k === 'Enter' && e.target && e.target.closest && e.target.closest('[data-now]')) return;   // the now card's links and buttons activate themselves
    if (k === '?') toggle('keys');
    else if (k === ',') toggle('look');
    else if (k === '/') $('addr').focus();
    else if (k === 'Escape') { if ($('verb')) $('verb').remove(); else if (!$('keys').hidden || !$('look').hidden) { $('keys').hidden = true; $('look').hidden = true; } else if (selectedFile) closeFile(); else if (filePreview) dismissFilePreview(); else out(); }
    else if (k === 'n' || k === 'N') {
      const a = document.querySelector('[data-needs] a');
      if (a) {
        if (st.d === 1) out();
        if (!toolbar.hidden && statusDetail.hidden) setStatusExpanded(true);
        a.focus();
      }
    }   // design §2.8: N jumps to what needs you
    else if (st.d === 1 && k === 'ArrowUp') openRail(railIdx() - 1);
    else if (st.d === 1 && k === 'ArrowDown') openRail(railIdx() + 1);
    else if (st.d === 1 && k === 'Home') openRail(0);
    else if (st.d === 1 && k === 'End') openRail(RAIL.length - 1);
    else if (!T) done = false;
    else if (k === 'Enter' && e.target.closest?.('.phone-view, .lhead, .rrow, .node, .cb')) done = false;
    else if (st.d === 0 && k === 'Enter') { if (st.focus) open(REP[st.focus] && isCluster(REP[st.focus]) ? REP[st.focus] : st.focus); }
    else if (st.d === 0 && (k === '+' || k === '=')) zoomStep(-1);
    else if (st.d === 0 && (k === '-' || k === '_')) zoomStep(1);
    else if (st.d === 0 && k === 'End') goNow();
    else if (st.d === 0 && k === 'Home') panTo(minL());
    else if (st.d === 0 && (k === 'ArrowLeft' || k === 'ArrowRight') && !e.shiftKey) panTo(view.L + (k === 'ArrowLeft' ? -1 : 1) * PW * 0.8);
    else if (st.d === 0 && ['ArrowLeft', 'ArrowRight', 'ArrowUp', 'ArrowDown'].includes(k)) {
      const cur = REP[st.focus] || st.focus, d = DISP[cur]; if (!d) return;
      const list = onRow(d.row), i = list.indexOf(cur);
      let nx = cur;
      if (k === 'ArrowLeft') nx = list[Math.max(0, i - 1)];
      else if (k === 'ArrowRight') nx = list[Math.min(list.length - 1, i + 1)];
      else {
        const j = ROWS.indexOf(T.homeRow(d.row)) + (k === 'ArrowUp' ? -1 : 1);
        if (ROWS[j]) { const cand = onRow(ROWS[j]); if (cand.length) nx = cand.reduce((b, id) => Math.abs(DISP[id].x - d.x) < Math.abs(DISP[b].x - d.x) ? id : b); }
      }
      if (nx && nx !== cur) { chooseTimelineNode(nx); ensureVisible(nx); }
    } else done = false;
    if (done) e.preventDefault();
  });
  window.addEventListener('resize', () => { fit(); camera.update(); });

  // Exposed for the screenshot driver.
  window.World = { open: label => { const n = Object.values(N).find(x => x.id === label || x.episode === label || x.label === label) || (T && N[T.aliasOf[label]]); if (n) open(n.id); return !!n; }, out: () => out(), rail: k => openRail(k),
    poll: () => refresh('poll'), generation: () => shownGen,
    focus: chooseTimelineNode, zoom: (u, ax) => zoomTo(u, ax == null ? PL + PW / 2 : ax, false), pan: L => panTo(L, false), now: () => goNow(false),
    view: () => ({ u: view.u, L: view.L, level: levelName(), atNow: atNow() }), screenX: id => screenX(id) };

  async function boot() {
    applyLook();
    $('L1').inert = true;
    const heads = location.hash === '#heads' || location.hash === '#heads/activity';
    const activity = location.hash === '#activity' || location.hash === '#heads/activity';
    if (heads || location.hash.startsWith('#draft=')) $('stage').classList.add('camera-loading');
    try { await document.fonts.ready; } catch (_) {}
    renderCrumbs(); fit();
    const q = parseAppQuery(location.search);
    const ok = await loadModel();
    if (ok) startLive();
    if (followMode) schedulePoll();
    if (ok && q && q.brane != null) await openState(q.brane, q.at, true, tabLocation().tab);   // a deep link, or a reload inside a brane
    if (ok && !q?.brane && (heads || location.hash.startsWith('#draft='))) {
      chooseAngle(true, {navigate:false, instant:true});
      if (location.hash.startsWith('#draft=')) {
        const id = new URLSearchParams(location.hash.slice(1)).get('draft');
        cameraURL('/surface#heads', false);
        if (N[id]?.draft) drafts.edit(N[N[id].parent], id); else say('This draft belongs to its original tab');
      }
    }
    $('stage').classList.remove('camera-loading');
    document.documentElement.dataset.ready = '1';
    if (!(ok && !q?.brane && activity && focusActivity()) && !document.querySelector('dialog[open]')) $('root').focus({ preventScroll: true });
  }
  boot();
})();
