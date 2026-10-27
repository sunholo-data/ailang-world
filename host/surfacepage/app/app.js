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
  const hm = iso => { const d = new Date(iso); return pad(d.getUTCHours()) + ':' + pad(d.getUTCMinutes()); };
  const dayOf = iso => { const d = new Date(iso); return d.getUTCDate() + ' ' + MON[d.getUTCMonth()]; };
  const exactOf = iso => dayOf(iso) + ' ' + hm(iso) + ' UTC';
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

  // ---------- Time axis: each shared state is a point; zoom levels blend "one slot each" into real time ----------
  const BREAK_H = 72;                           // longer quiet gaps are drawn as an axis break, not as empty weeks
  const LEVELS = [{ name: 'Close', G: 84, P: 0 }, { name: 'Steps', G: 50, P: 0 }, { name: 'Hours', G: 14, P: 16 }, { name: 'Days', G: 3, P: 6 }, { name: 'Weeks', G: 0, P: 2.4 }];
  const breakAt = (TT, i) => i > 0 && (TT[i] - TT[i - 1]) / 3.6e6 > BREAK_H;
  // levelXs is each level's x for every shared state, then now; TT is their times, then now's.
  function levelXs(TT) {
    const NS = TT.length - 1;
    return LEVELS.map(L => {
      const xs = [220];
      for (let i = 1; i <= NS; i++) {
        const dth = (TT[i] - TT[i - 1]) / 3.6e6;
        const gap = breakAt(TT, i) ? Math.max(L.G, 64) : i === NS ? Math.max(L.G * 1.44, dth * L.P, 28) : Math.max(L.G, dth * L.P);
        xs.push(xs[i - 1] + gap);
      }
      return xs;
    });
  }
  // blendXS is the x of every point at zoom u, between two levels.
  function blendXS(LX, u) { const k = Math.min(LX.length - 2, Math.floor(u)), f = u - k; return LX[k].map((x, i) => x + (LX[k + 1][i] - x) * f); }

  // ---------- Canvas geometry (1600 × 936 design canvas, scaled to the window) ----------
  const PL = 262, PR = 1262, PW = PR - PL, TICK_T = 768, MINI_T = 818, MINI_H = 36, LEGEND_T = 892;

  // ---------- Model: world/surface-model/v1 adapted to the timeline's lookups ----------
  function adaptModel(D) {
    const LANES = {}, N = {};
    (D.lanes || []).forEach(l => { LANES[l.key] = l; (l.states || []).forEach(s => { N[s.id] = Object.assign({}, s, { lane: l.key }); }); });
    const parents = {}, kids = {};
    const addEdge = (p, c) => { (kids[p] = kids[p] || []).push(c); (parents[c] = parents[c] || []).push(p); };
    (D.links || []).forEach(k => { if (N[k.from] && N[k.to]) addEdge(k.from, k.to); }); // links run parent (from) → child (to)
    const SH = (LANES.shared ? LANES.shared.states || [] : []).filter(s => s.kind === 'shared').sort((a, b) => a.rank - b.rank).map(s => s.id);
    for (let i = 1; i < SH.length; i++) addEdge(SH[i - 1], SH[i]);  // the shared world's own first-parent line
    const nameOf = key => { const l = LANES[key]; return l ? l.name || l.title || key : String(key); };
    const laneTitle = key => rowLabel(nameOf(key), LANES[key] && LANES[key].row);
    const glyphOf = (key, size) => { const l = LANES[key] || {}; return gl(l.glyph, l.glyph === 'circle' ? SHARED_INK : l.colour, size); };
    const glyphOfNode = (n, size) => n.kind === 'shared' ? gl('circle', SHARED_INK, size) : n.lane === 'shared' ? gl(null, null, size) : glyphOf(n.lane, size);
    const mergeOf = id => (kids[id] || []).find(c => N[c].kind === 'shared');
    const STEPS = Object.values(N).filter(n => n.kind === 'step');
    const TNOW = Date.parse(D.now), NOW_D = new Date(D.now);
    const NS = SH.length;                         // index NS is "now"
    const SIDX = {}; SH.forEach((id, i) => { SIDX[id] = i; });
    const TT = SH.map(id => Date.parse(N[id].at)).concat([TNOW]);
    const LX = levelXs(TT);
    // Rows: the most recent brane above the shared world, the next three below, the rest share one quiet lane.
    const recent = {};
    STEPS.forEach(n => { if (n.lane !== 'shared') recent[n.lane] = Math.max(recent[n.lane] == null ? -1 : recent[n.lane], n.rank); });
    const byRecent = Object.keys(recent).sort((a, b) => recent[b] - recent[a]);
    const above = byRecent.slice(0, 1), below = byRecent.slice(1, 4);
    const quiet = Object.keys(LANES).filter(k => k !== 'shared' && !above.includes(k) && !below.includes(k));
    const LY = {};
    above.forEach(k => { LY[k] = 182; }); LY.shared = 320; below.forEach((k, i) => { LY[k] = 440 + i * 98; }); LY.other = 440 + below.length * 98 - 6;
    const ROWS = above.concat(['shared'], below, quiet.length ? ['other'] : []);
    const rowOf = id => { const n = N[id]; return n.kind === 'shared' ? 'shared' : n.lane === 'shared' ? 'dip' : LY[n.lane] != null && n.lane !== 'other' ? n.lane : 'other'; };
    const rowY = r => r === 'dip' ? LY.shared + 44 : LY[r];
    const isBreak = i => breakAt(TT, i);
    const lastStep = key => STEPS.filter(n => n.lane === key).sort((a, b) => b.rank - a.rank)[0];
    const lastState = key => ((LANES[key] && LANES[key].states) || []).reduce((a, s) => !a || s.rank > a.rank ? s : a, null);
    return { D, LANES, N, parents, kids, SH, NS, SIDX, TT, LX, MAXU: LX.length - 1, STEPS, TNOW, NOW_D, above, below, quiet, LY, ROWS,
      nameOf, laneTitle, glyphOf, glyphOfNode, mergeOf, rowOf, rowY, isBreak, lastStep, lastState };
  }

  // ---------- Layout of the visible graph at one zoom (XS) ----------
  const isCluster = id => typeof id === 'string' && id.startsWith('c:');
  function layoutGraph(T, XS) {
    const { N, SIDX, NS, LY } = T, POS = {}, REP = {}, DISP = {};
    Object.keys(N).forEach(id => {
      const n = N[id], r = T.rowOf(id);
      if (n.kind === 'shared') { POS[id] = [XS[SIDX[id]], LY.shared]; return; }
      const m = T.mergeOf(id);
      let x;
      if (m != null && SIDX[m] != null) { const i = SIDX[m], xm = XS[i], xp = i > 0 ? XS[i - 1] : xm - 52; x = xm - Math.min(26, Math.max(0, (xm - xp) / 2)); }
      else x = XS[NS] - 26;
      POS[id] = [x, T.rowY(r)];
    });
    // Collapse dense runs on each row into a count bubble.
    const byRow = {};
    Object.keys(POS).forEach(id => { (byRow[T.rowOf(id)] = byRow[T.rowOf(id)] || []).push(id); });
    Object.entries(byRow).forEach(([row, ids]) => {
      ids.sort((a, b) => POS[a][0] - POS[b][0] || N[a].rank - N[b].rank);
      const SEP = row === 'shared' ? 15 : 19;
      let g = [];
      const flush = () => {
        if (!g.length) return;
        if (g.length === 1) { REP[g[0]] = g[0]; DISP[g[0]] = { x: POS[g[0]][0], y: POS[g[0]][1], ids: g, row }; }
        else { const rid = 'c:' + g[0]; const x = g.reduce((a, id) => a + POS[id][0], 0) / g.length; g.forEach(id => { REP[id] = rid; }); DISP[rid] = { x, y: POS[g[0]][1], ids: g, row }; }
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
      if (n.kind === 'shared' || d.row === 'dip') return;
      if (lastX[d.row] == null || lastX[d.row] - d.x >= 112) { LABELED.add(id); lastX[d.row] = d.x; }
    });
    return { POS, REP, DISP, LABELED };
  }

  // ---------- Brane view: world/surface-stack/v1 as bars, the front card and the rail ----------
  const tagOf = c => c.kind === 'shared' ? `Merged ${c.pr ? '#' + c.pr : c.label}` : (c.episode || c.label);
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
    const c = b.card, g = c.kind === 'shared' ? gl('circle', SHARED_INK, 11) : b.col.brane === 'shared' ? gl(null, null, 11) : gl(b.col.glyph, b.col.colour, 11);
    const t = cleanTitle(c.title);
    return `<span class="a">${b.dir < 0 ? '↑' : '↓'}</span>${g}<span class="tag">${esc(tagOf(c))}</span>${t ? `<span class="ttl">${esc(t)}</span>` : ''}<span style="flex:1"></span><span class="w">${esc(c.rel)}</span>${c.badge > 0 ? `<span class="nb inl">${c.badge}</span>` : ''}${b.laneLab ? `<span class="ln">${esc(b.laneLab)}</span>` : ''}`;
  }

  const chip = (id, label, glyphHtml, via, mono, tip, extra = '') => `<button type="button" class="chip${mono ? ' mono-l' : ''}" data-go="${esc(id)}"${tip ? ` data-tip="${esc(tip)}"` : ''}${extra}>${glyphHtml}<span class="lb">${esc(label)}</span>${via ? `<span class="via">${esc(via)}</span>` : ''}</button>`;
  // stateChip names a state from the stack, glyphed by its lane in the model when the model has it.
  const stateChip = (T, s, via) => {
    if (!s) return '';
    const n = (T && T.N[s.id]) || s;
    const g = n.kind === 'shared' ? gl('circle', SHARED_INK, 10) : T && n.lane && n.lane !== 'shared' ? T.glyphOf(n.lane, 10) : gl(null, null, 10);
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

  // frontHTML is the front card for the stack's front state. T (the adapted model) only names
  // and glyphs other states; every value on the card is the stack's own, or an honest gap.
  function frontHTML(S, T) {
    const F = S.front, step = F.kind === 'step';
    const none = t => `<div class="muted">${esc(t)}</div>`;
    const tag = F.episode ? F.episode + (F.model ? ' · ' + prettyModel(F.model) : '') : F.label;
    const bar = `<div class="fbar"><span id="cb-tag" style="font-weight:700">${esc(tag)}</span><span class="sub" id="cb-when">${esc(`${F.rel} · ${F.exact}`)}</span><span style="flex:1"></span><span class="sub" style="display:inline-flex;align-items:center;gap:7px">${gl(F.glyph, F.glyph === 'circle' ? SHARED_INK : F.colour, 11)}${esc(rowLabel(F.name, F.row))}</span><span class="sub" id="cb-pos">${esc(F.position)}</span></div>`;
    const summ = summaryOf(F);
    const head = `<div class="head"><div class="h">${esc(cleanTitle(F.title) || F.label)}</div>${summ ? `<div class="outcome">${esc(summ)}</div>` : ''}</div>`;
    const view = `<section class="view" aria-label="From the agent"><b>From the agent</b><span>No view from the agent yet</span><span style="flex:1"></span><span class="pill">A2UI · World catalog</span><span class="pill" style="border-style:dashed;border-color:var(--ctlCard)">row 151</span></section>`;

    // Inputs: the agent, the task, then what it came from and read.
    let grid = '';
    if (F.model) {
      const other = /^codex/i.test(F.agent || '') ? 'Claude' : 'Codex';
      const agents = [[prettyModel(F.model), true], [other, false]].map(([a, on]) => `<button type="button" class="sbtn" aria-pressed="${on}" disabled style="cursor:default;color:var(--ink)">${esc(a)}</button>`).join('');
      grid += `<span class="k">Agent</span><div class="seg" role="group" aria-label="Agent">${agents}</div>`;
    }
    if (F.task) grid += `<span class="k top">Task</span><div class="task"><div class="tl1">${esc(F.task)}</div></div>`;
    const gaps = !step ? [] : F.episode ? ['Grants: not recorded yet.'] : ['No episode record for this step: its agent, task, grants and reads are not recorded.'];
    const came = [];
    if (F.builtOn) came.push(stateChip(T, F.builtOn, 'built on'));
    (F.parents || []).forEach(p => came.push(stateChip(T, p, 'parent')));
    (F.reads || []).forEach(rd => {
      const g = rd.via === 'commit' || !rd.glyph ? gl(null, null, 10) : gl(rd.glyph, rd.colour, 10);
      came.push(chip(rd.fromId, rd.label, g, 'read', rd.via === 'commit', 'reads:' + rd.via + ':' + rd.fromId, ` data-read="${esc(rd.fromId)}"`));
    });
    const reads = F.readCalls != null ? ` · ${F.readCalls} read${F.readCalls === 1 ? '' : 's'} over ${F.readFiles} file${F.readFiles === 1 ? '' : 's'}` : '';
    const readGap = step && F.episode && !F.reads ? none('Reads: not recorded for this step.') : '';
    const inputs = `<div class="col"><div class="colh">Inputs</div>${grid ? `<div class="ingrid">${grid}</div>` : ''}${gaps.map(none).join('')}
      ${came.length || reads ? `<div class="sect"><div class="l">↑ Came from${reads}</div><div class="chips">${came.join('')}</div></div>` : ''}${readGap}</div>`;

    // Outputs: the files changed, the PR, and the merge it led to.
    const fl = F.files || [];
    const uniq = p => { const parts = p.split('/'); for (let k = 1; k <= parts.length; k++) { const suf = parts.slice(-k).join('/'); if (fl.filter(g => g.path === suf || g.path.endsWith('/' + suf)).length === 1) return suf; } return p; };
    const files = fl.map(f => {
      const tag = f.status === 'A' ? 'new' : f.status === 'D' ? 'del' : 'edit';
      const tc = tag === 'new' ? 'var(--ink)' : tag === 'del' ? 'var(--youInk)' : 'var(--inkSub)';
      return `<button type="button" class="file" data-tip="path:${esc(f.path)}" aria-label="${esc(f.path + ', ' + tag)}"><span class="nm">${esc(uniq(f.path))}</span><span class="tg" style="color:${tc}">${tag}</span></button>`;
    }).join('');
    const more = F.filesTruncated > 0 ? none(`and ${F.filesTruncated} more not listed`) : '';
    const fileGap = F.files ? '' : none(`Files: not recorded for this ${step ? 'step' : 'state'}.`);
    const pr = F.pr ? `<div class="prpill">#${F.pr}${F.files ? ` · ${fl.length} file${fl.length === 1 ? '' : 's'}` : ''}</div>` : '';
    const led = F.mergedInto ? `<div class="sect" style="margin-top:auto"><div class="l">↓ Led to</div><div class="chips">${stateChip(T, F.mergedInto, '')}</div></div>` : '';
    const outputs = `<div class="col g6"><div class="colh">Outputs</div>${files}${more}${fileGap}${pr}${led}</div>`;

    const acts = (F.decisions || []).map(a => {
      const [bg, fg, bd] = ACT[a.actor] || ACT.world;
      return `<div class="act"${a.source ? ` data-tip="src:${esc(a.source)}"` : ''}><span class="who" style="background:${bg};color:${fg};border-color:${bd}">${WHO[a.actor] || esc(a.actor)}</span><span class="tx">${esc(a.text)}<span class="tm">${esc(a.rel)}</span></span></div>`;
    }).join('');
    const decisions = `<div class="col g9"><div class="colh">Decisions</div>${acts || none('None recorded.')}</div>`;

    const ck = F.checks || {};
    const v = verdict(ck.graders, ck.rounds, F.firstRoundPass, ck.minutes);
    const tiles = [[ck.graders, 'graders'], [ck.rounds, ck.rounds === 1 ? 'round' : 'rounds'], [ck.minutes, 'agent min'], [ck.calls, 'World calls']]
      .filter(([x]) => x != null && x !== '').map(([x, k]) => `<div class="ctile"><div class="v">${esc(x)}</div><div class="k">${k}</div></div>`).join('');
    const checks = `<div class="col g9"><div class="colh">Checks</div>${v ? `<div class="verdict">${esc(v)}</div>` : none(step ? 'No episode record for this step.' : 'A merge has no checks of its own.')}${tiles ? `<div class="tiles">${tiles}</div>` : ''}</div>`;

    return bar + `<div class="fbody">${head}${view}<div class="cols">${inputs}${outputs}${decisions}${checks}</div></div>`;
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

  if (typeof module !== 'undefined') module.exports = { splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
    adaptModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW } };
  if (typeof document === 'undefined') return;

  // ================= The page =================
  const $ = id => document.getElementById(id);
  const SEAMS = window.WorldSurface || {};
  let T = null, D = null, LANES = {}, N = {}, parents = {}, kids = {}, SH = [], NS = 0, SIDX = {}, TT = [], LX = [], MAXU = LEVELS.length - 1;
  let STEPS = [], TNOW = 0, NOW_D = new Date(), quiet = [], LY = {}, ROWS = [], MXS = [], XS = [];
  let nameOf, laneTitle, glyphOf, glyphOfNode, mergeOf, isBreak, lastStep, lastState;
  function setModel(m) {
    T = adaptModel(m);
    ({ D, LANES, N, parents, kids, SH, NS, SIDX, TT, LX, MAXU, STEPS, TNOW, NOW_D, quiet, LY, ROWS, nameOf, laneTitle, glyphOf, glyphOfNode, mergeOf, isBreak, lastStep, lastState } = T);
    const a = LX[3][0], z = LX[3][NS];
    MXS = LX[3].map(x => 8 + (x - a) / ((z - a) || 1) * (PW - 16));
    XS = LX[1].slice();
  }
  // The open stack (world/surface-stack/v1), when a brane is open.
  let S = null, F = null, RAIL = [], READS = new Map();
  function setStack(s) { S = s; F = s.front; RAIL = s.rail || []; READS = new Map((F.reads || []).map(r => [r.via + ':' + r.fromId, r])); }
  const railIdx = () => F ? Math.max(0, RAIL.findIndex(r => r.id === F.id)) : 0;
  const setView = v => { $('stage').dataset.view = v; };

  const computeXS = u => { XS = blendXS(LX, u); };
  const posF = f => { f = Math.max(0, Math.min(NS, f)); const i = Math.floor(f), j = Math.min(NS, i + 1); return XS[i] + (XS[j] - XS[i]) * (f - i); };
  const fAt = x => {
    if (x <= XS[0]) return 0; if (x >= XS[NS]) return NS;
    let lo = 0, hi = NS; while (hi - lo > 1) { const m = (lo + hi) >> 1; if (XS[m] <= x) lo = m; else hi = m; }
    return lo + (XS[hi] === XS[lo] ? 0 : (x - XS[lo]) / (XS[hi] - XS[lo]));
  };
  const fOfTime = t => {
    if (t < TT[0] || t > TT[NS]) return null;
    let lo = 0, hi = NS; while (hi - lo > 1) { const m = (lo + hi) >> 1; if (TT[m] <= t) lo = m; else hi = m; }
    if (isBreak(hi)) return null;
    return lo + (TT[hi] === TT[lo] ? 0 : (t - TT[lo]) / (TT[hi] - TT[lo]));
  };
  const level = () => view.u < 1.5 ? 'steps' : view.u < 2.5 ? 'hours' : view.u < 3.5 ? 'days' : 'weeks';
  const levelName = () => LEVELS[Math.round(view.u)].name;

  // ---------- Zoom/fit ----------
  let scale = 1;
  function fit() {
    const stg = $('stage'), w = stg.clientWidth, h = stg.clientHeight;
    scale = Math.min(w / 1600, h / 936, 1.25);
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
  function layout() { ({ POS, REP, DISP, LABELED } = layoutGraph(T, XS)); }

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
  function statusOf(key) {
    const s = key === 'shared' ? N[SH[NS - 1]] : lastStep(key);
    const on = s && dayDiff(s.at) <= 1;
    return `<span class="stat${on ? ' on' : ''}" data-tip="${s ? 'node:' + s.id : ''}">${on ? 'active' : 'quiet'}</span>`;
  }
  function progressOf(key) {
    const l = LANES[key];
    if (key === 'shared') return `<div class="prog">${NS} merges into dev</div>`;
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
      <span class="l1">${glyphOf(key, 11)}<span class="nm">${esc(nameOf(key))}</span></span>
      ${l.goal ? `<span class="goal" data-tip="goal:${esc(key)}">${esc(l.goal)}</span>` : ''}
      ${progressOf(key).replace('</div>', statusOf(key) + '</div>')}</button>`;
  }
  function otherHead(y) {
    const steps = quiet.reduce((a, k) => a + (LANES[k].stepsLanded != null ? LANES[k].stepsLanded : LANES[k].states.length), 0);
    const last = quiet.map(lastStep).filter(Boolean).sort((a, b) => b.rank - a.rank)[0];
    return `<button type="button" class="lhead" id="quiet" aria-expanded="false" style="top:${y - 24}px">
      <span class="l1"><span class="nm">${quiet.length} quieter branes</span><span class="gs">${quiet.map(k => glyphOf(k, 9)).join('')}</span></span>
      <span class="prog">${steps} steps landed${last ? ' · last ' + esc(last.rel) : ''}</span></button>`;
  }

  function renderTimeline(keepView = false) {
    let h = '';
    ROWS.forEach((key, k) => {
      const y = LY[key], tall = key === 'other' ? 72 : 98;
      h += `<div class="band" style="top:${y - tall / 2 - 4}px;height:${tall}px;width:${PR}px;background:${k % 2 === 0 ? 'var(--laneBand)' : 'transparent'}"></div>`;
      h += `<div class="lrule" style="top:${y}px;left:${PL}px;width:${PW}px"></div>`;
      h += key === 'other' ? otherHead(y) : laneHead(key, y);
    });
    h += `<div class="pedge" style="left:${PL - 1}px;top:96px;height:${TICK_T + 36 - 96}px"></div><div class="pedge" style="left:${PR}px;top:96px;height:${TICK_T + 36 - 96}px"></div>`;
    h += `<div class="plot" id="plot" style="left:${PL}px;top:0;width:${PW}px;height:${TICK_T + 40}px"><div class="world" id="world"></div></div>`;
    h += `<div class="tlbar" style="left:${PL}px;top:62px;width:${PW}px">
      <span class="range" id="range"></span><span style="flex:1"></span>
      <span class="ctl">
        <button type="button" class="cb" id="pg-back" aria-label="Back in time" data-tip="k:← back a page · drag or scroll to pan">‹</button>
        <button type="button" class="cb" id="pg-fwd" aria-label="Forward in time" data-tip="k:→ forward a page">›</button>
        <button type="button" class="cb now-b" id="now-b" aria-label="Back to now" data-tip="k:End · back to the newest state">Now</button>
      </span>
      <span class="ctl zoom" role="group" aria-label="Zoom">
        <button type="button" class="cb" id="z-out" aria-label="Zoom out" data-tip="k:− · ⌘-scroll or pinch to zoom">−</button>
        <span class="zl" id="zl">Steps</span>
        <button type="button" class="cb" id="z-in" aria-label="Zoom in" data-tip="k:+ · ⌘-scroll or pinch to zoom">+</button>
      </span></div>`;
    h += nowCard();
    h += `<div class="mini-lab" style="top:${MINI_T + 2}px">Overview<span>drag the window · click to jump</span></div>`;
    h += `<div class="mini" id="mini" style="left:${PL}px;top:${MINI_T}px;width:${PW}px;height:${MINI_H}px"></div>`;
    h += `<div class="quiet-list" id="quiet-list" style="left:16px;bottom:${936 - (LY.other - 28)}px" hidden>${quiet.map(k => {
      const s = lastStep(k);
      return `<div data-brane="${esc(k)}">${glyphOf(k, 11)}<span class="qn"><b>${esc(nameOf(k))}</b><span class="qg">${esc(LANES[k].goal || ('row ' + LANES[k].row))}</span></span><span class="w">${esc(s ? s.rel : '')}</span></div>`;
    }).join('')}</div>`;
    h += `<div class="legend" aria-label="Legend" style="top:${LEGEND_T}px">
      <span><span class="nb inl">1</span>your decisions</span>
      <span><span style="width:14px;height:14px;border-radius:7px;box-shadow:0 0 0 4px var(--halo)"></span>connected to the step you point at</span>
      <span><span style="width:18px;border-top:2px dotted var(--read)"></span>read</span>
      <span><span class="bub inl">4</span>several merges or steps · click to open up</span>
      <span class="off"><span style="width:18px;border-top:2px dashed var(--ran)"></span>ran on <i>· not recorded yet</i></span>
    </div>`;
    $('tl').innerHTML = h;
    // The plot clips; it must never scroll itself (focus() or scrollIntoView on a node would).
    $('plot').addEventListener('scroll', e => { e.target.scrollTop = 0; e.target.scrollLeft = 0; });
    renderMini();
    relayout();
    if (!keepView) goNow(false);
  }

  function nowCard() {
    const ne = D.nowEdge || {}, now = NOW_D, sess = ne.sessions || [], needs = ne.needsYou || [], td = ne.today;
    const healthy = /healthy/i.test(ne.health);
    let today = '';
    if (td) {
      const by = Object.entries(td.byAgent || {}).map(([a, n]) => `${cap(a)} ${n}`).join(', ');
      today = `<div class="nsec"><div class="nh">Today</div>
        <div class="nl"><b>${td.landed} step${td.landed === 1 ? '' : 's'} landed</b>${by ? ` · ${esc(by)}` : ''}</div>
        <div class="nl sub">${td.firstRound} first-round pass${td.firstRound === 1 ? '' : 'es'} · ${td.fixRounds} fix round${td.fixRounds === 1 ? '' : 's'} · ${Math.round(td.agentMinutes)} agent-minutes</div></div>`;
    }
    const run = sess[0], runLane = run && Object.keys(LANES).find(k => LANES[k].row && new RegExp('^w' + LANES[k].row + '[a-z]?$').test(run.episode));
    const running = ne.sessionsUnavailable ? `<div class="nsec"><div class="nh">Running</div><div class="nl sub">${esc(ne.sessionsUnavailable)}</div></div>`
      : `<div class="nsec"><div class="nh">Running · ${sess.length}</div>
        ${run ? `<div class="nl"><b>${esc(run.episode)}</b>${runLane ? ` · ${esc(nameOf(runLane))}` : ''}</div><div class="nl sub">started ${esc(run.started)}</div>` : '<div class="nl sub">Nothing running</div>'}</div>`;
    return `<div class="now" data-now style="left:${PR + 22}px;top:56px">
      <div><div style="font-size:12px;font-weight:600;color:var(--sub)">Now</div><div style="font-size:15px;font-weight:700">${DOW[now.getUTCDay()]} ${dayOf(D.now)} · ${hm(D.now)}</div></div>
      <div style="display:flex;gap:9px;align-items:center"><span style="width:9px;height:9px;border-radius:5px;background:${healthy ? 'var(--okDot)' : 'var(--you)'};flex:none"></span><span style="font-weight:700;color:var(--ok)">${healthy ? 'Healthy' : esc(ne.health)}</span><span style="font-size:12.5px;color:var(--sub)">checked ${esc(ne.checked)}</span></div>
      <div class="nl sub" data-live>${esc(liveStatus(live))}</div>
      ${today}
      ${running}
      <div class="nsec"><div class="nh" style="color:${needs.length ? 'var(--youShell)' : 'var(--sub)'}">Needs you · ${needs.length}</div>
        ${needs.length ? '' : `<div class="nl sub">Nothing needs you</div>`}</div>
    </div>`;
  }

  // The empty and refused states: the reason, quietly, where the timeline would be.
  function renderVoid(kind, reason) {
    T = null;
    $('tl').innerHTML = `<div style="position:absolute;left:${PL}px;top:360px;width:${PW}px;display:flex;flex-direction:column;align-items:center;gap:6px;text-align:center">
      <div style="font-size:15px;font-weight:700;color:var(--sub)">${kind === 'empty' ? 'Nothing to show yet' : 'The world could not be read'}</div>
      <div style="font-size:13px;color:var(--sub);max-width:640px">${esc(reason)}</div></div>`;
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
      h += `<div class="brk" style="left:${mx - 40}px;top:${LY.shared - 30}px">${days} quiet days</div><div class="brk-z" style="left:${mx - 4}px;top:${LY.shared - 9}px"></div>`;
    }
    Object.entries(DISP).forEach(([id, d]) => {
      const { x, y } = d;
      if (isCluster(id)) {
        const ns = d.ids.map(k => N[k]), shared = d.row === 'shared';
        const word = shared ? 'merges' : 'steps';
        const col = shared || d.row === 'dip' ? 'var(--edge)' : d.row === 'other' ? 'var(--edge)' : LANES[d.row].colour;
        const badge = ns.reduce((a, n) => a + (n.badge || 0), 0);
        h += `<button type="button" class="node bub-n" data-id="${id}" data-tip="node:${id}" style="left:${x}px;top:${y}px" aria-label="${ns.length} ${word}, ${esc(exactOf(ns[0].at))} to ${esc(exactOf(ns[ns.length - 1].at))}"><span class="halo"><span class="bub" style="border-color:${col}">${ns.length}</span></span>${d.ids.map(k => `<span data-state="${esc(k)}" hidden></span>`).join('')}</button>`;
        if (badge) h += `<div class="nb" style="left:${x + 8}px;top:${y - 21}px">${badge}</div>`;
        return;
      }
      const n = N[id], shared = n.kind === 'shared', unfiled = !shared && n.lane === 'shared';
      const l = LANES[n.lane];
      const size = shared ? 9 : unfiled ? 9 : l.glyph === 'square' ? 15 : 18;
      const inner = shared ? `<span class="g ring" style="width:9px;height:9px;border-radius:50%;background:var(--shellBg);border:2px solid var(--edge)"></span>`
        : unfiled ? gl(null, null, 9) : gl(l.glyph, l.colour, size);
      const what = shared ? `Merged ${n.pr ? '#' + n.pr : n.label}` : (cleanTitle(n.title) || n.episode || n.label);
      h += `<button type="button" class="node${shared ? ' sh' : ''}" data-id="${id}" data-state="${esc(id)}" data-tip="node:${id}" style="left:${x}px;top:${y}px" aria-label="${esc((shared ? 'Shared world' : nameOf(n.lane)) + ': ' + what + ', ' + n.rel + (n.badge ? ', ' + n.badge + ' of your decisions' : ''))}"><span class="halo">${inner}</span></button>`;
      if (n.badge > 0) h += `<div class="nb" style="left:${x + 6}px;top:${y - 19}px">${n.badge}</div>`;
      if (LABELED.has(id)) h += `<div class="nlab" data-l="${id}" style="left:${x - 56}px;bottom:${936 - (y - 18)}px"><span class="tx">${esc(shortTitle(n.title || '') || n.episode || n.label)}</span></div>`;
      h += `<div class="nlab" data-r="${id}" style="left:${x - 56}px;top:${y + (shared ? 9 : 14)}px"><div class="rel"></div></div>`;
    });
    h += `<div class="nowline" style="left:${XS[NS]}px;top:92px;height:${TICK_T + 30 - 92}px"></div>`;
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
      SH.forEach((id, i) => {
        const n = N[id], d = dayOf(n.at), x = XS[i];
        if (d !== lastDay) { place(x, d, 'tick-d', 76, pd); lastDay = d; place(x, hm(n.at), 'tick-t', 44, pt); }
      });
      SH.forEach((id, i) => place(XS[i], hm(N[id].at), 'tick-t', 44, pt));
    } else {
      const H = 3.6e6, start = Math.floor(TT[0] / 864e5) * 864e5;
      const cands = [];
      for (let t = start; t <= TNOW; t += H) {
        const f = fOfTime(t); if (f == null) continue;
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
    e += `<path d="M 0 ${LY.shared} L ${XS[0]} ${LY.shared}" fill="none" stroke="var(--edge)" stroke-width="1.4" stroke-dasharray="3 5"></path>`;
    const seen = new Set();
    Object.keys(POS).forEach(c => {
      const n = N[c], rc = REP[c], [x2, y2] = [DISP[rc].x, DISP[rc].y];
      const ps = parents[c] || [];
      if (n.kind === 'step' && !ps.length && T.rowOf(c) !== 'dip') {  // its base is further back than this data: it arrives from the left
        e += `<path d="M ${x2 - 70} ${y2} L ${x2} ${y2}" fill="none" stroke="${n.lane === 'shared' ? 'var(--edge)' : LANES[n.lane].colour}" stroke-width="1.4" stroke-dasharray="3 5" opacity="0.6"></path>`;
      }
      ps.forEach(p => {
        const rp = REP[p]; if (!rp || rp === rc) return;
        const key = rp + '>' + rc; if (seen.has(key)) return; seen.add(key);
        const [x1, y1] = [DISP[rp].x, DISP[rp].y], L = lit(rp, rc);
        const ownKey = n.kind === 'step' ? n.lane : N[p].kind === 'step' ? N[p].lane : 'shared';
        const col = L ? 'var(--lit)' : ownKey === 'shared' ? 'var(--edge)' : LANES[ownKey].colour;
        const w = L ? 2.6 : 1.6, op = L || ownKey === 'shared' ? 1 : 0.6;
        let d;
        const brk = n.kind === 'shared' && N[p].kind === 'shared' && isBreak(SIDX[c]);
        if (y1 === y2) d = `M ${x1} ${y1} L ${x2} ${y2}`;
        else if (n.kind === 'shared') d = `M ${x1} ${y1} C ${x2} ${y1}, ${x2} ${y2}, ${x2} ${y2}`;
        else { const k = Math.min(30, Math.max(4, (x2 - x1) / 2)); d = `M ${x1} ${y1} C ${x1 + k} ${y1}, ${x1} ${y2}, ${x1 + k} ${y2} L ${x2} ${y2}`; }
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
  const atNow = () => Math.abs(view.L - clampL(maxL())) < 1;
  let anim = 0;
  const ease = t => 1 - Math.pow(1 - t, 3);
  function applyView(keepPosition = false) {
    if (!keepPosition) view.L = clampL(view.L);
    $('world').style.transform = `translateX(${-view.L}px)`;
    const fL = fAt(view.L), fR = fAt(view.L + PW), mw = $('mwin');
    if (mw) { const a = MX(fL), b = MX(fR); mw.style.left = a + 'px'; mw.style.width = Math.max(8, b - a) + 'px'; }
    const tL = timeOfF(fL), tR = timeOfF(fR), nowIn = view.L + PW >= XS[NS] - 1;
    $('range').innerHTML = `<b>${fmtRange(tL)}</b> → <b>${nowIn ? 'now' : fmtRange(tR)}</b>${atNow() ? '' : `<span class="behind">${behindText(tR)}</span>`}`;
    $('now-b').classList.toggle('live', !atNow());
    $('pg-fwd').disabled = atNow(); $('pg-back').disabled = view.L <= minL() + 1;
  }
  const timeOfF = f => { const i = Math.floor(Math.max(0, Math.min(NS, f))), j = Math.min(NS, i + 1); return TT[i] + (TT[j] - TT[i]) * (f - i); };
  const fmtRange = t => { const d = new Date(t); return `${DOW[d.getUTCDay()]} ${d.getUTCDate()} ${MON[d.getUTCMonth()]} ${pad(d.getUTCHours())}:${pad(d.getUTCMinutes())}`; };
  const behindText = tR => { const h = (TNOW - tR) / 3.6e6; return h < 1 ? '' : h < 36 ? `${Math.round(h)} h before now` : `${Math.round(h / 24)} days before now`; };
  function panTo(L, animate = true) {
    cancelAnimationFrame(anim);
    L = clampL(L);
    if (!animate || reduced()) { view.L = L; applyView(); return; }
    const L0 = view.L, t0 = performance.now(), dur = 420;
    const step = now => { const k = Math.min(1, (now - t0) / dur); view.L = L0 + (L - L0) * ease(k); applyView(); if (k < 1) anim = requestAnimationFrame(step); };
    anim = requestAnimationFrame(step);
  }
  const goNow = (animate = true) => panTo(maxL(), animate);
  function setZoom(u, ax) {   // ax: the design-x on screen that stays put
    const f = fAt(view.L + (ax - PL));
    view.u = Math.max(0, Math.min(MAXU, u));
    relayout();
    view.L = clampL(posF(f) - (ax - PL));
    applyView();
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
  const zoomStep = dir => zoomTo(dir < 0 ? Math.ceil(view.u - 0.01) - 1 : Math.floor(view.u + 0.01) + 1, keyAnchor());
  function ensureVisible(id, animate = true) {
    const x = screenX(id); if (x == null) return;
    if (x < PL + 60) panTo(view.L - (PL + 60 - x) - PW * 0.25, animate);
    else if (x > PR - 40) panTo(view.L + (x - (PR - 40)) + PW * 0.25, animate);
  }

  // Overview: a fixed real-time layout of the whole history (the Days level), normalised to the strip.
  const MX = f => { f = Math.max(0, Math.min(NS, f)); const i = Math.floor(f), j = Math.min(NS, i + 1); return MXS[i] + (MXS[j] - MXS[i]) * (f - i); };
  const fOfMX = mx => { if (mx <= MXS[0]) return 0; if (mx >= MXS[NS]) return NS; let i = 0; while (i < NS && MXS[i + 1] < mx) i++; return i + (MXS[i + 1] === MXS[i] ? 0 : (mx - MXS[i]) / (MXS[i + 1] - MXS[i])); };
  function renderMini() {
    const rowH = 4, top = 5;
    let h = '';
    const start = Math.floor(TT[0] / 864e5) * 864e5;
    for (let t = start; t <= TNOW; t += 864e5) {
      const f = fOfTime(t); if (f == null) continue;
      const x = MX(f), d = new Date(t);
      h += `<div class="mday" style="left:${x}px"></div><div class="mdl" style="left:${x + 3}px">${d.getUTCDate()} ${MON[d.getUTCMonth()]}</div>`;
    }
    for (let i = 1; i <= NS; i++) if (isBreak(i)) h += `<div class="mbrk" style="left:${(MXS[i - 1] + MXS[i]) / 2 - 6}px">≈</div>`;
    ROWS.forEach((r, k) => { h += `<div class="mrow" style="top:${top + k * (rowH + 1)}px"></div>`; });
    SH.forEach((id, i) => { h += `<div class="mk" style="left:${MXS[i]}px;top:${top + ROWS.indexOf('shared') * (rowH + 1)}px;height:${rowH}px;width:1px;background:var(--edge)"></div>`; });
    STEPS.forEach(n => {
      const m = mergeOf(n.id); if (m == null) return;
      const r = T.rowOf(n.id) === 'dip' ? 'shared' : T.rowOf(n.id), k = ROWS.indexOf(r);
      h += `<div class="mk" style="left:${MX(SIDX[m] - 0.15) - 1}px;top:${top + k * (rowH + 1)}px;height:${rowH}px;width:3px;background:${n.lane === 'shared' ? 'var(--edge)' : LANES[n.lane].colour}"></div>`;
    });
    h += `<div class="mnow" style="left:${MXS[NS]}px"></div><div class="mwin" id="mwin"></div>`;
    $('mini').innerHTML = h;
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

  function renderBrane(animate) {
    if (!$('bstage')) {
      $('br').innerHTML = `<div class="bstage" id="bstage"><article class="front" id="front"></article></div><nav class="rail" id="rail" aria-label="Brane timeline"></nav>`;
    }
    $('front').innerHTML = frontHTML(S, T);
    $('front').dataset.front = F.id;
    $('front').setAttribute('aria-label', cleanTitle(F.title) || F.label);
    placeBars(animate);
    $('rail').innerHTML = railHTML(S);
  }

  // ---------- Header ----------
  const addrOf = () => `${F.brane}@${railIdx() + 1}`;
  function renderCrumbs() {
    $('crumbs').innerHTML = crumbsHTML(st.d === 1 && F ? addrOf() : null);
    $('addr').placeholder = st.d === 0 ? 'go to a place — row-189@6 · shared@0ffbb9c' : '';
  }
  let msgT;
  function say(t) { const m = $('hmsg'); m.textContent = t; m.style.opacity = '1'; clearTimeout(msgT); msgT = setTimeout(() => { m.style.opacity = '0'; }, 2600); }
  function follow(path) { try { if (location.pathname + location.search !== path) history.replaceState(null, '', path); } catch (_) { /* the view still moves */ } }

  // ---------- Live data ----------
  async function getJSON(url) {
    const r = await fetch(url, { headers: { Accept: 'application/json' }, cache: 'no-store' });
    let body = null;
    try { body = await r.json(); } catch (_) { body = null; }
    return { ok: r.ok, status: r.status, body };
  }
  async function loadModel() {
    let res;
    try { res = await getJSON('/surface/model.json'); } catch (e) { return renderVoid('error', 'the surface model could not be fetched: ' + e.message); }
    const b = res.body || {};
    if (!res.ok) return renderVoid('error', b.reason || `the surface model answered HTTP ${res.status}`);
    if (b.empty) return renderVoid('empty', b.reason || 'the surface index is empty');
    if (!Array.isArray(b.lanes) || !b.lanes.some(l => l.states && l.states.length)) return renderVoid('empty', 'The surface index holds no states yet.');
    setModel(b);
    const newest = STEPS.reduce((a, n) => !a || n.rank > a.rank ? n : a, null);   // the latest work, as v3 opened on its front
    st.focus = newest ? newest.id : SH[NS - 1] || null;
    renderTimeline();
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
  async function refreshLive() {
    liveRefresh = null;
    refreshing = true;
    const count = live.count;
    try {
      const selection = st.d === 1 && F ? {brane: F.brane, at: F.id, seq: openSeq} : null;
      const model = await getJSON('/surface/model.json');
      if (!model.ok || !model.body || !Array.isArray(model.body.lanes)) throw new Error('refresh failed (HTTP ' + model.status + ')');
      const stack = selection ? await getJSON(stackPath(selection.brane, selection.at)) : null;
      if (stack && (!stack.ok || !stack.body || !stack.body.front)) throw new Error('refresh failed (HTTP ' + stack.status + ')');
      // Fetch everything before touching the old render. Navigation during the fetch wins.
      const active = document.activeElement, activeId = active && active.id;
      const activeKey = active && ['data-id', 'data-state', 'data-go', 'data-rail', 'data-lane'].find(k => active.hasAttribute(k));
      const activeValue = activeKey && active.getAttribute(activeKey);
      const focus = st.focus, u = view.u, L = view.L, depth = st.d;
      const railTop = $('rail') ? $('rail').scrollTop : 0, address = $('addr').value;
      setModel(model.body);
      const newest = STEPS.reduce((a, n) => !a || n.rank > a.rank ? n : a, null);
      st.focus = N[focus] ? focus : newest ? newest.id : SH[NS - 1] || null;
      view.u = u; view.L = L;
      renderTimeline(true);
      applyView(true);
      if (stack && selection.seq === openSeq && depth === 1) { setStack(stack.body); renderBrane(false); $('rail').scrollTop = railTop; }
      st.d = depth;
      setView(depth === 1 ? 'brane' : 'timeline');
      $('addr').value = address;
      if (active && active.isConnected) active.focus({preventScroll: true});
      else {
        const replacement = activeId ? $(activeId) : activeKey ? Array.from(document.querySelectorAll('[' + activeKey + ']')).find(e => e.getAttribute(activeKey) === activeValue && e.tabIndex >= 0) : null;
        if (replacement) replacement.focus({preventScroll: true});
      }
      const cursor = model.body.nowEdge && model.body.nowEdge.live;
      if (cursor) live.cursor = Math.max(live.cursor, cursor.lastIndex);
      live.count = Math.max(0, live.count - count);
      live.checked = hm(new Date());
      indexP = null;
    } catch (e) {
      queued = false;
      livePause(/^refresh failed/.test(e.message) ? e.message : 'refresh failed (HTTP 0)');
    } finally {
      refreshing = false;
      if (queued) { queued = false; queueRefresh(); }
      livePaint();
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
  function open(id) {
    if (!T) return;
    if (isCluster(id)) { const d = DISP[id]; return zoomTo(Math.ceil(view.u - 0.01) - 1, PL + d.x - view.L); }
    const n = N[id];
    if (!n) { say('That state is not on the timeline'); return; }
    return openState(n.lane, id);
  }
  // openState fetches the state's stack and zooms into it; a newer open supersedes an older one.
  let openSeq = 0;
  async function openState(brane, at) {
    const seq = ++openSeq;
    let res;
    try { res = await getJSON(stackPath(brane, at)); } catch (e) { if (seq === openSeq) say('the stack could not be fetched: ' + e.message); return false; }
    if (seq !== openSeq) return false;
    const b = res.body || {};
    if (!res.ok || !b.front) { say(b.reason || `no such state (HTTP ${res.status})`); if (st.d === 0) follow(appPath()); return false; }
    const was = st.d === 1;
    setStack(b);
    if (!was && T) { const sx = screenX(F.id), d = DISP[REP[F.id]]; setOrigin(sx != null ? Math.max(PL, Math.min(PR, sx)) : PL, d ? d.y : LY.shared); }
    st.d = 1; st.focus = F.id; $('tip').hidden = true;
    renderBrane(was);
    $('root').classList.add('zoomed'); setView('brane');
    renderCrumbs(); follow(appPath(F.brane, F.id)); $('root').focus({ preventScroll: true });
    return true;
  }
  function out(focusId) {
    if (st.d === 0) return;
    openSeq++;
    const id = focusId || (F && F.id);
    if (T && id && POS[id]) {
      st.focus = id;
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
  async function go(raw) {
    const parse = SEAMS.parseAddress;
    let r = addressTarget(raw, null, parse);
    if (r.error && typeof parse === 'function') {
      try { r = addressTarget(raw, await loadIndex(), parse); } catch (e) { r = { error: 'the surface index could not be read: ' + e.message }; }
    }
    if (r.home) return out();
    if (r.open) return openState(r.open.brane, r.open.at);
    say(r.error);
  }

  // ---------- Popovers ----------
  const KEYROWS = [['← →', 'back · forward a page', '—'], ['⇧← ⇧→', 'previous · next state', '—'], ['↑ ↓', 'lane above · below', 'back in time · toward now'], ['+ −', 'zoom in · out', '—'], ['End', 'back to now', 'most recent'], ['Home', 'the oldest state', 'earliest'], ['Enter', 'open brane', '—'], ['Esc', '—', 'back to timeline'], ['/', 'address bar', 'address bar'], [',', 'appearance', 'appearance'], ['?', 'these keys', 'these keys']];
  function renderKeys() {
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
      const ns = d.ids.map(k => N[k]).sort((a, b) => b.rank - a.rank), shared = d.row === 'shared';
      const t0 = ns[ns.length - 1].at, t1 = ns[0].at;
      const span = dayOf(t0) === dayOf(t1) ? `${dayOf(t0)} · ${hm(t0)}–${hm(t1)}` : `${exactOf(t0)} → ${exactOf(t1)}`;
      return `<div class="hc"><div class="hk">${ns.length} ${shared ? 'merges into the shared world' : 'steps'} · ${esc(span)}</div>` +
        ns.slice(0, 6).map(n => `<div class="hli">${glyphOfNode(n, 9)}<span>${esc(shortTitle(n.title, 64) || n.episode || n.label)}</span></div>`).join('') +
        `<div class="hid">${ns.length > 6 ? `+${ns.length - 6} more · ` : ''}click or zoom in to open them up</div></div>`;
    }
    const n = N[id]; if (!n) return '';
    const shared = n.kind === 'shared';
    const ids = [n.episode, n.pr ? '#' + n.pr : '', id.slice(0, 7)].filter(Boolean).join(' · ');
    if (shared) {
      const from = (parents[id] || []).map(p => N[p]).find(p => p.kind === 'step');
      return `<div class="hc"><div class="hk">${gl('circle', SHARED_INK, 9)}Merged into the shared world · ${esc(n.rel)}</div>
        ${n.title ? `<div class="ht">${esc(cleanTitle(n.title))}</div>` : ''}
        ${from && from.lane !== 'shared' ? `<div class="hm">from ${esc(nameOf(from.lane))}${from.episode ? ' · ' + esc(from.episode) : ''}</div>` : ''}
        <div class="hw">${esc(exactOf(n.at))}</div><div class="hid">${esc(ids)}</div></div>`;
    }
    const s = summaryOf(n);
    const g = n.graders ? `${esc(n.graders)} graders · ${n.firstRoundPass ? 'first round ✓' : `${n.rounds} round${n.rounds === 1 ? '' : 's'}`}` : '';
    return `<div class="hc"><div class="hk">${glyphOfNode(n, 9)}${esc(n.lane === 'shared' ? 'Straight to the shared world' : nameOf(n.lane))}</div>
      <div class="ht">${esc(cleanTitle(n.title) || n.episode || n.label)}</div>
      ${s ? `<div class="hs">${esc(s)}</div>` : `<div class="hs none">No agent episode recorded.</div>`}
      ${n.model ? `<div class="hm"><b>${esc(agentLine(n))}</b>${g ? ' · ' + g : ''}${n.minutes != null ? ` · ${n.minutes} agent-min` : ''}</div>` : ''}
      ${n.badge ? `<div class="hm"><span class="nb inl">${n.badge}</span> ${n.badge === 1 ? 'decision' : 'decisions'} of yours</div>` : ''}
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
  $('stage').addEventListener('pointerover', e => {
    const n = e.target.closest('#tl .node');
    if (n && T && st.d === 0 && !dragging && st.focus !== n.dataset.id) { st.focus = n.dataset.id; applyFocus(); }
  });
  let suppressClick = false;
  $('stage').addEventListener('click', e => {
    if (suppressClick) { suppressClick = false; e.stopPropagation(); return; }
    const n = e.target.closest('#tl .node');
    if (n) return open(n.dataset.id);
    if (e.target.closest('#quiet')) { const l = $('quiet-list'); l.hidden = !l.hidden; $('quiet').setAttribute('aria-expanded', String(!l.hidden)); return; }
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
    const ch = e.target.closest('[data-go]');
    if (ch) {
      const id = ch.dataset.go, k = RAIL.findIndex(r => r.id === id);
      if (F && id === F.id) return;
      return k >= 0 ? openRail(k) : N[id] ? openState(N[id].lane, id) : say(`${ch.textContent.trim() || id.slice(0, 7)} is outside the timeline`);
    }
  }, true);

  // Drag to pan (the plot), drag/click the overview window.
  let drag = null;
  $('stage').addEventListener('pointerdown', e => {
    if (!T || st.d !== 0 || e.button !== 0) return;
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
    if (!T || st.d !== 0) return;
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
    if (st.d === 1) out(); else if (T) goNow();
  });
  $('keys-btn').addEventListener('click', () => toggle('keys'));
  $('look-btn').addEventListener('click', () => toggle('look'));
  $('look').addEventListener('click', e => { const b = e.target.closest('[data-look]'); if (!b) return; const [k, v] = b.dataset.look.split(':'); look[k] = v; saveLook(); applyLook(); renderLook(); });
  $('addr').addEventListener('keydown', e => {
    e.stopPropagation();
    if (e.key === 'Enter') { go(e.target.value); e.target.value = ''; $('root').focus({ preventScroll: true }); }
    if (e.key === 'Escape') { e.target.value = ''; $('root').focus({ preventScroll: true }); }
  });
  const onRow = row => Object.keys(DISP).filter(id => DISP[id].row === row).sort((a, b) => DISP[a].x - DISP[b].x);
  document.addEventListener('keydown', e => {
    if (e.target && e.target.tagName === 'INPUT') return;
    const k = e.key; let done = true;
    if (k === 'Enter' && e.target && e.target.id === 'home') return;   // the link's own activation: its click goes home
    if (k === '?') toggle('keys');
    else if (k === ',') toggle('look');
    else if (k === '/') $('addr').focus();
    else if (k === 'Escape') { if (!$('keys').hidden || !$('look').hidden) { $('keys').hidden = true; $('look').hidden = true; } else out(); }
    else if (st.d === 1 && k === 'ArrowUp') openRail(railIdx() - 1);
    else if (st.d === 1 && k === 'ArrowDown') openRail(railIdx() + 1);
    else if (st.d === 1 && k === 'Home') openRail(0);
    else if (st.d === 1 && k === 'End') openRail(RAIL.length - 1);
    else if (!T) done = false;
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
        const j = ROWS.indexOf(d.row === 'dip' ? 'shared' : d.row) + (k === 'ArrowUp' ? -1 : 1);
        if (ROWS[j]) { const cand = onRow(ROWS[j]); if (cand.length) nx = cand.reduce((b, id) => Math.abs(DISP[id].x - d.x) < Math.abs(DISP[b].x - d.x) ? id : b); }
      }
      if (nx && nx !== cur) { st.focus = nx; applyFocus(); ensureVisible(nx); }
    } else done = false;
    if (done) e.preventDefault();
  });
  window.addEventListener('resize', fit);

  // Exposed for the screenshot driver.
  window.World = { open: label => { const n = Object.values(N).find(x => x.id === label || x.episode === label || x.label === label); if (n) open(n.id); return !!n; }, out: () => out(), rail: k => openRail(k),
    focus: id => { st.focus = id; applyFocus(); }, zoom: (u, ax) => zoomTo(u, ax == null ? PL + PW / 2 : ax, false), pan: L => panTo(L, false), now: () => goNow(false),
    view: () => ({ u: view.u, L: view.L, level: levelName(), atNow: atNow() }), screenX: id => screenX(id) };

  async function boot() {
    applyLook();
    try { await document.fonts.ready; } catch (_) {}
    renderCrumbs(); fit();
    const q = parseAppQuery(location.search);
    const ok = await loadModel();
    if (ok) startLive();
    if (ok && q && q.brane != null) await openState(q.brane, q.at);   // a deep link, or a reload inside a brane
    document.documentElement.dataset.ready = '1';
    $('root').focus({ preventScroll: true });
  }
  boot();
})();
