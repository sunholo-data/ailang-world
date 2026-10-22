// World — Tauri spike v3 frontend. The v2 artifact-faithful port, plus a navigable timeline
// (pan, zoom from steps to weeks, an overview strip) and human-first reading (titles, the agent's
// own summary, lane goals and progress). Renders ONLY spike-data.json; nothing is invented.
(() => {
  'use strict';
  const $ = id => document.getElementById(id);
  const esc = s => String(s == null ? '' : s).replace(/[&<>"]/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c]));
  const MON = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  const DOW = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  const pad = n => String(n).padStart(2, '0');
  const hm = iso => { const d = new Date(iso); return pad(d.getUTCHours()) + ':' + pad(d.getUTCMinutes()); };
  const dayOf = iso => { const d = new Date(iso); return d.getUTCDate() + ' ' + MON[d.getUTCMonth()]; };
  const exactOf = iso => dayOf(iso) + ' ' + hm(iso) + ' UTC';
  const prettyModel = m => (m || '').split('-').reduce((a, w, i, all) => {
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

  // ---------- Model ----------
  const D = window.SPIKE_DATA;
  const S = D.stack, F = S.front, RAIL = S.rail;
  const LANES = {}, N = {};
  D.lanes.forEach(l => { LANES[l.key] = l; l.states.forEach(s => { N[s.id] = { ...s, lane: l.key }; }); });
  const parents = {}, kids = {};
  const addEdge = (p, c) => { (kids[p] = kids[p] || []).push(c); (parents[c] = parents[c] || []).push(p); };
  D.links.forEach(k => { if (N[k.from] && N[k.to]) addEdge(k.from, k.to); }); // links run parent (from) → child (to)
  const SH = LANES.shared.states.filter(s => s.kind === 'shared').sort((a, b) => a.rank - b.rank).map(s => s.id);
  for (let i = 1; i < SH.length; i++) addEdge(SH[i - 1], SH[i]);  // the shared world's own first-parent line
  const nameOf = key => LANES[key].name || LANES[key].title;
  const laneTitle = key => { const l = LANES[key]; return l.row ? `${nameOf(key)} · row ${l.row}` : nameOf(key); };
  const glyphOf = (key, size) => { const l = LANES[key]; return gl(l.glyph, l.glyph === 'circle' ? SHARED_INK : l.colour, size); };
  const glyphOfNode = (n, size) => n.kind === 'shared' ? gl('circle', SHARED_INK, size) : n.lane === 'shared' ? gl(null, null, size) : gl(LANES[n.lane].glyph, LANES[n.lane].colour, size);
  const READS = new Map(F.reads.map(r => [r.fromId, r]));
  const frontIdx = Math.max(0, RAIL.findIndex(r => r.id === F.id));
  const mergeOf = id => (kids[id] || []).find(c => N[c].kind === 'shared');
  const STEPS = Object.values(N).filter(n => n.kind === 'step');
  const TNOW = Date.parse(D.now);
  const NOW_D = new Date(D.now);

  // ---------- Time axis: each shared state is a point; zoom levels blend "one slot each" into real time ----------
  const NS = SH.length;                         // index NS is "now"
  const SIDX = {}; SH.forEach((id, i) => { SIDX[id] = i; });
  const TT = SH.map(id => Date.parse(N[id].at)).concat([TNOW]);
  const BREAK_H = 72;                           // longer quiet gaps are drawn as an axis break, not as empty weeks
  const isBreak = i => i > 0 && (TT[i] - TT[i - 1]) / 3.6e6 > BREAK_H;
  const LEVELS = [{ name: 'Close', G: 84, P: 0 }, { name: 'Steps', G: 50, P: 0 }, { name: 'Hours', G: 14, P: 16 }, { name: 'Days', G: 3, P: 6 }, { name: 'Weeks', G: 0, P: 2.4 }];
  const LX = LEVELS.map(L => {
    const xs = [220];
    for (let i = 1; i <= NS; i++) {
      const dth = (TT[i] - TT[i - 1]) / 3.6e6;
      const gap = isBreak(i) ? Math.max(L.G, 64) : i === NS ? Math.max(L.G * 1.44, dth * L.P, 28) : Math.max(L.G, dth * L.P);
      xs.push(xs[i - 1] + gap);
    }
    return xs;
  });
  const MAXU = LEVELS.length - 1;
  let XS = LX[1].slice();
  const computeXS = u => { const k = Math.min(MAXU - 1, Math.floor(u)), f = u - k; XS = LX[k].map((x, i) => x + (LX[k + 1][i] - x) * f); };
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

  // ---------- Canvas geometry (1600 × 936 design canvas, scaled to the window) ----------
  const PL = 262, PR = 1262, PW = PR - PL, TICK_T = 768, MINI_T = 818, MINI_H = 36, LEGEND_T = 892;
  const recent = {};
  STEPS.forEach(n => { if (n.lane !== 'shared') recent[n.lane] = Math.max(recent[n.lane] || -1, n.rank); });
  const byRecent = Object.keys(recent).sort((a, b) => recent[b] - recent[a]);
  const above = byRecent.slice(0, 1), below = byRecent.slice(1, 4);
  const quiet = Object.keys(LANES).filter(k => k !== 'shared' && !above.includes(k) && !below.includes(k));
  const LY = {};
  LY[above[0]] = 182; LY.shared = 320; below.forEach((k, i) => { LY[k] = 440 + i * 98; }); LY.other = 440 + below.length * 98 - 6;
  const ROWS = [above[0], 'shared', ...below, 'other'];
  const rowOf = id => { const n = N[id]; return n.kind === 'shared' ? 'shared' : n.lane === 'shared' ? 'dip' : LY[n.lane] != null && n.lane !== 'other' ? n.lane : 'other'; };
  const rowY = r => r === 'dip' ? LY.shared + 44 : LY[r];
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
  const st = { d: 0, focus: F.id, ri: frontIdx };
  const view = { u: 1, L: 0 };
  const reduced = () => $('root').classList.contains('reduced');

  // ---------- Layout of the visible graph at the current zoom ----------
  let POS = {}, REP = {}, DISP = {}, LABELED = new Set();
  const isCluster = id => typeof id === 'string' && id.startsWith('c:');
  function layout() {
    POS = {}; REP = {}; DISP = {};
    Object.keys(N).forEach(id => {
      const n = N[id], r = rowOf(id);
      if (n.kind === 'shared') { POS[id] = [XS[SIDX[id]], LY.shared]; return; }
      const m = mergeOf(id);
      let x;
      if (m != null && SIDX[m] != null) { const i = SIDX[m], xm = XS[i], xp = i > 0 ? XS[i - 1] : xm - 52; x = xm - Math.min(26, Math.max(0, (xm - xp) / 2)); }
      else x = XS[NS] - 26;
      POS[id] = [x, rowY(r)];
    });
    // Collapse dense runs on each row into a count bubble.
    const byRow = {};
    Object.keys(POS).forEach(id => { (byRow[rowOf(id)] = byRow[rowOf(id)] || []).push(id); });
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
    LABELED = new Set();
    const lastX = {};
    Object.keys(DISP).filter(id => !isCluster(id)).sort((a, b) => DISP[b].x - DISP[a].x).forEach(id => {
      const d = DISP[id], n = N[id];
      if (n.kind === 'shared' || d.row === 'dip') return;
      if (lastX[d.row] == null || lastX[d.row] - d.x >= 112) { LABELED.add(id); lastX[d.row] = d.x; }
    });
  }

  // ---------- Timeline (static frame) ----------
  const linksOf = id => {
    const r = new Map();
    if (isCluster(id)) return r;
    (parents[id] || []).forEach(p => r.set(p, 'built on'));
    (kids[id] || []).forEach(c => r.set(c, 'built on'));
    if (id === F.id) READS.forEach((_, src) => { if (!r.has(src)) r.set(src, 'read'); });
    if (READS.has(id) && !r.has(F.id)) r.set(F.id, 'read');
    return r;
  };
  const dayDiff = iso => { const a = new Date(iso), b = NOW_D; return Math.round((Date.UTC(b.getUTCFullYear(), b.getUTCMonth(), b.getUTCDate()) - Date.UTC(a.getUTCFullYear(), a.getUTCMonth(), a.getUTCDate())) / 864e5); };
  const lastStep = key => STEPS.filter(n => n.lane === key).sort((a, b) => b.rank - a.rank)[0];
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
    return `<button type="button" class="lhead" data-lane="${key}" style="top:${y - 32}px">
      <span class="l1">${glyphOf(key, 11)}<span class="nm">${esc(nameOf(key))}</span></span>
      ${l.goal ? `<span class="goal" data-tip="goal:${key}">${esc(l.goal)}</span>` : ''}
      ${progressOf(key).replace('</div>', statusOf(key) + '</div>')}</button>`;
  }
  function otherHead(y) {
    const steps = quiet.reduce((a, k) => a + (LANES[k].stepsLanded != null ? LANES[k].stepsLanded : LANES[k].states.length), 0);
    const last = quiet.map(lastStep).filter(Boolean).sort((a, b) => b.rank - a.rank)[0];
    return `<button type="button" class="lhead" id="quiet" aria-expanded="false" style="top:${y - 24}px">
      <span class="l1"><span class="nm">${quiet.length} quieter branes</span><span class="gs">${quiet.map(k => glyphOf(k, 9)).join('')}</span></span>
      <span class="prog">${steps} steps landed${last ? ' · last ' + esc(last.rel) : ''}</span></button>`;
  }

  function renderTimeline() {
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
      return `<div>${glyphOf(k, 11)}<span class="qn"><b>${esc(nameOf(k))}</b><span class="qg">${esc(LANES[k].goal || ('row ' + LANES[k].row))}</span></span><span class="w">${esc(s ? s.rel : '')}</span></div>`;
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
    goNow(false);
  }

  function nowCard() {
    const ne = D.nowEdge, now = NOW_D, sess = ne.sessions || [], needs = ne.needsYou || [], td = ne.today;
    const healthy = /healthy/i.test(ne.health);
    let today = '';
    if (td) {
      const by = Object.entries(td.byAgent || {}).map(([a, n]) => `${cap(a)} ${n}`).join(', ');
      today = `<div class="nsec"><div class="nh">Today</div>
        <div class="nl"><b>${td.landed} step${td.landed === 1 ? '' : 's'} landed</b>${by ? ` · ${esc(by)}` : ''}</div>
        <div class="nl sub">${td.firstRound} first-round pass${td.firstRound === 1 ? '' : 'es'} · ${td.fixRounds} fix round${td.fixRounds === 1 ? '' : 's'} · ${Math.round(td.agentMinutes)} agent-minutes</div></div>`;
    }
    const run = sess[0], runLane = run && Object.keys(LANES).find(k => LANES[k].row && new RegExp('^w' + LANES[k].row + '[a-z]?$').test(run.episode));
    return `<div class="now" style="left:${PR + 22}px;top:56px">
      <div><div style="font-size:12px;font-weight:600;color:var(--sub)">Now</div><div style="font-size:15px;font-weight:700">${DOW[now.getUTCDay()]} ${dayOf(D.now)} · ${hm(D.now)}</div></div>
      <div style="display:flex;gap:9px;align-items:center"><span style="width:9px;height:9px;border-radius:5px;background:${healthy ? 'var(--okDot)' : 'var(--you)'};flex:none"></span><span style="font-weight:700;color:var(--ok)">${healthy ? 'Healthy' : esc(ne.health)}</span><span style="font-size:12.5px;color:var(--sub)">checked ${esc(ne.checked)}</span></div>
      ${today}
      <div class="nsec"><div class="nh">Running · ${sess.length}</div>
        ${run ? `<div class="nl"><b>${esc(run.episode)}</b>${runLane ? ` · ${esc(nameOf(runLane))}` : ''}</div><div class="nl sub">started ${esc(run.started)}</div>` : '<div class="nl sub">Nothing running</div>'}</div>
      <div class="nsec"><div class="nh" style="color:${needs.length ? 'var(--youShell)' : 'var(--sub)'}">Needs you · ${needs.length}</div>
        ${needs.length ? '' : `<div class="nl sub">Nothing needs you</div>`}</div>
    </div>`;
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
        h += `<button type="button" class="node bub-n" data-id="${id}" data-tip="node:${id}" style="left:${x}px;top:${y}px" aria-label="${ns.length} ${word}, ${esc(exactOf(ns[0].at))} to ${esc(exactOf(ns[ns.length - 1].at))}"><span class="halo"><span class="bub" style="border-color:${col}">${ns.length}</span></span></button>`;
        if (badge) h += `<div class="nb" style="left:${x + 8}px;top:${y - 21}px">${badge}</div>`;
        return;
      }
      const n = N[id], shared = n.kind === 'shared', unfiled = !shared && n.lane === 'shared';
      const l = LANES[n.lane];
      const size = shared ? 9 : unfiled ? 9 : l.glyph === 'square' ? 15 : 18;
      const inner = shared ? `<span class="g ring" style="width:9px;height:9px;border-radius:50%;background:var(--shellBg);border:2px solid var(--edge)"></span>`
        : unfiled ? gl(null, null, 9) : gl(l.glyph, l.colour, size);
      const what = shared ? `Merged ${n.pr ? '#' + n.pr : n.label}` : (cleanTitle(n.title) || n.episode || n.label);
      h += `<button type="button" class="node${shared ? ' sh' : ''}" data-id="${id}" data-tip="node:${id}" style="left:${x}px;top:${y}px" aria-label="${esc((shared ? 'Shared world' : nameOf(n.lane)) + ': ' + what + ', ' + n.rel + (n.badge ? ', ' + n.badge + ' of your decisions' : ''))}"><span class="halo">${inner}</span></button>`;
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
      if (n.kind === 'step' && !ps.length && rowOf(c) !== 'dip') {  // its base is further back than this data: it arrives from the left
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
  function applyView() {
    view.L = clampL(view.L);
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
  const MXS = LX[3].map(x => 8 + (x - LX[3][0]) / (LX[3][NS] - LX[3][0]) * (PW - 16));
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
      const r = rowOf(n.id) === 'dip' ? 'shared' : rowOf(n.id), k = ROWS.indexOf(r);
      h += `<div class="mk" style="left:${MX(SIDX[m] - 0.15) - 1}px;top:${top + k * (rowH + 1)}px;height:${rowH}px;width:3px;background:${n.lane === 'shared' ? 'var(--edge)' : LANES[n.lane].colour}"></div>`;
    });
    h += `<div class="mnow" style="left:${MXS[NS]}px"></div><div class="mwin" id="mwin"></div>`;
    $('mini').innerHTML = h;
  }

  // ---------- Brane view ----------
  const tagOf = c => c.kind === 'shared' ? `Merged ${c.pr ? '#' + c.pr : c.label}` : (c.episode || c.label);
  function columns() {
    const ri = st.ri;
    if (ri === frontIdx) return { earlier: S.earlier, later: S.later };
    const c = RAIL[ri], lane = c ? N[c.id].lane : F.brane, l = LANES[lane];
    const col = cards => ({ brane: lane, glyph: l.glyph, colour: l.colour, title: l.title, cards });
    const later = RAIL.slice(ri + 1, ri + 4);
    const m = mergeOf(RAIL[ri].id);
    return { earlier: ri > 0 ? [col(RAIL.slice(Math.max(0, ri - 3), ri))] : [],
      later: [later.length ? col(later) : null, m ? { brane: 'shared', glyph: 'circle', colour: SHARED_INK, title: 'Shared world', cards: [N[m]] } : null].filter(Boolean) };
  }
  function stackBars(cols, dir) {
    const nc = Math.max(1, cols.length), cw = (1180 - (nc - 1) * 12) / nc, out = [];
    cols.slice(0, 3).forEach((col, c) => {
      const near = dir < 0 ? col.cards.slice().reverse() : col.cards;
      near.slice(0, 3).forEach((card, i) => {
        const dd = i + 1, shrink = (dd - 1) * 24;
        out.push({ key: (dir < 0 ? 'u:' : 'd:') + card.id, card, col, dd,
          width: cw - shrink, ml: -590 + c * (cw + 12) + shrink / 2, top: dir < 0 ? 106 - (dd - 1) * 46 : 796 + (dd - 1) * 46,
          op: Math.max(0.35, 1 - (dd - 1) * 0.28), z: 39 - dd, dir, laneLab: nc > 1 && dd === 1 ? (col.brane === 'shared' ? 'Shared world' : nameOf(col.brane)) : '' });
      });
    });
    return out;
  }
  function barHTML(b) {
    const c = b.card, g = c.kind === 'shared' ? gl('circle', SHARED_INK, 11) : gl(b.col.glyph, b.col.colour, 11);
    const t = cleanTitle(N[c.id] ? N[c.id].title : c.title);
    return `<span class="a">${b.dir < 0 ? '↑' : '↓'}</span>${g}<span class="tag">${esc(tagOf(c))}</span>${t ? `<span class="ttl">${esc(t)}</span>` : ''}<span style="flex:1"></span><span class="w">${esc(c.rel)}</span>${c.badge > 0 ? `<span class="nb inl">${c.badge}</span>` : ''}${b.laneLab ? `<span class="ln">${esc(b.laneLab)}</span>` : ''}`;
  }
  function placeBars(animate) {
    const { earlier, later } = columns();
    const bars = stackBars(earlier, -1).concat(stackBars(later, 1));
    const host = $('bstage'), have = new Map([...host.querySelectorAll('.strip')].map(el => [el.dataset.key, el]));
    const keep = new Set();
    bars.forEach(b => {
      keep.add(b.key);
      let el = have.get(b.key);
      const fresh = !el;
      if (fresh) {
        el = document.createElement('button'); el.type = 'button'; el.className = 'strip'; el.dataset.key = b.key; el.dataset.id = b.card.id;
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

  const chip = (id, label, glyphHtml, via, mono, tip) => `<button type="button" class="chip${mono ? ' mono-l' : ''}" data-go="${id}"${tip ? ` data-tip="${tip}"` : ''}>${glyphHtml}<span class="lb">${esc(label)}</span>${via ? `<span class="via">${esc(via)}</span>` : ''}</button>`;
  const stateChip = (id, via) => {
    const n = N[id]; if (!n) return '';
    return chip(id, n.kind === 'shared' ? (n.pr ? `#${n.pr}` : n.label) : (n.episode || n.label), glyphOfNode(n, 10), via, false, 'node:' + id);
  };
  const ACT = { you: ['var(--you)', 'var(--onYou)', 'var(--you)'], agent: ['#eaf1f7', '#2f6690', '#b9cfe1'], world: ['#f6f8fa', '#526677', '#dce3e8'], review: ['#f1edf6', '#6b4f8a', '#d6cbe4'], loop: ['#eaf3ee', '#3f7a5e', '#c3dccf'] };
  const WHO = { you: 'You', agent: 'Agent', world: 'World', review: 'Review', loop: 'Loop' };

  function verdict(graders, rounds, firstPass, minutes) {
    if (!graders) return '';
    const [a, b] = String(graders).split('/').map(Number);
    const all = a === b;
    let s = all ? `Passed all ${b} checks` : `${a} of ${b} checks passed`;
    s += firstPass ? ' on the first round' : rounds > 1 ? ` after ${rounds} rounds` : '';
    return s + (minutes != null ? ` · ${minutes} agent-minutes` : '');
  }

  function cardHTML() {
    const r = RAIL[st.ri], isF = r.id === F.id, n = N[r.id], lane = n.lane;
    const bar = `<div class="fbar"><span id="cb-tag" style="font-weight:700">${esc(isF ? `${F.episode} · ${prettyModel(F.model)}` : (n.episode ? `${n.episode}${n.model ? ' · ' + prettyModel(n.model) : ''}` : r.label))}</span><span class="sub" id="cb-when">${esc(isF ? `${F.rel} · ${F.exact}` : `${r.rel} · ${exactOf(r.at)}`)}</span><span style="flex:1"></span><span class="sub" style="display:inline-flex;align-items:center;gap:7px">${glyphOf(lane, 11)}${esc(laneTitle(lane))}</span><span class="sub" id="cb-pos">${st.ri + 1} / ${RAIL.length}</span></div>`;
    const summ = summaryOf(isF ? F : n);
    const head = `<div class="head"><div class="h">${esc(cleanTitle(r.title) || r.label)}</div>${summ ? `<div class="outcome">${esc(summ)}</div>` : ''}</div>`;
    const view = `<section class="view" aria-label="From the agent"><b>From the agent</b><span>No view from the agent yet</span><span style="flex:1"></span><span class="pill">A2UI · World catalog</span><span class="pill" style="border-style:dashed;border-color:var(--ctlCard)">row 151</span></section>`;
    const bases = (parents[r.id] || []).filter(p => N[p]), merge = mergeOf(r.id);
    if (!isF) {
      const none = t => `<div class="muted">${t}</div>`;
      const v = verdict(n.graders, n.rounds, n.firstRoundPass, n.minutes);
      const tiles = n.graders ? `<div class="tiles">${[[n.graders, 'graders'], [String(n.rounds), n.rounds === 1 ? 'round' : 'rounds'], [String(n.minutes), 'agent min']].map(([vv, k]) => `<div class="ctile"><div class="v">${esc(vv)}</div><div class="k">${k}</div></div>`).join('')}</div>` : '';
      return bar + `<div class="fbody">${head}${view}<div class="cols">
        <div class="col"><div class="colh">Inputs</div>
          ${n.model ? `<div class="ingrid"><span class="k">Agent</span><div class="seg"><button type="button" class="sbtn" aria-pressed="true" disabled style="cursor:default;color:var(--ink)">${esc(prettyModel(n.model))}</button></div></div>` : ''}
          ${none(`Task, grants and reads: not in this data (only ${esc(F.episode)}’s card was captured).`)}
          ${bases.length ? `<div class="sect"><div class="l">↑ Came from</div><div class="chips">${bases.map(p => stateChip(p, 'built on')).join('')}</div></div>` : ''}</div>
        <div class="col g6"><div class="colh">Outputs</div>${r.pr ? `<div class="prpill">#${r.pr}</div>` : ''}
          ${merge ? `<div class="sect" style="margin-top:auto"><div class="l">↓ Led to</div><div class="chips">${stateChip(merge, '')}</div></div>` : ''}</div>
        <div class="col g9"><div class="colh">Decisions</div>${r.badge ? none(`${r.badge} of yours, not itemised in this data.`) : none('None recorded.')}</div>
        <div class="col g9"><div class="colh">Checks</div>${v ? `<div class="verdict">${esc(v)}</div>${tiles}` : none('No episode record in this data.')}</div>
      </div></div>`;
    }
    const agents = [[prettyModel(F.model), true], ['Codex', false]].map(([a, on]) => `<button type="button" class="sbtn" aria-pressed="${on}" disabled style="cursor:default;color:var(--ink)">${esc(a)}</button>`).join('');
    const grants = F.grants.map(g => `<button type="button" class="gbtn" aria-pressed="${!!g.on}"${g.locked ? ' disabled' : ''} style="cursor:default">${g.locked ? '<svg width="10" height="11" viewBox="0 0 10 11" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="1.3"><rect x="1" y="4.6" width="8" height="5.6" rx="1"/><path d="M3 4.6V3a2 2 0 0 1 4 0v1.6"/></svg>' : ''}${esc(g.name)}</button>`).join('');
    const came = (F.builtOn ? stateChip(F.builtOn.id, 'built on') : '') + F.reads.map(rd => {
      const g = rd.via === 'commit' || !rd.glyph ? gl(null, null, 10) : gl(rd.glyph, rd.colour, 10);
      return chip(rd.fromId, rd.label, g, 'read', rd.via === 'commit', 'reads:' + rd.fromId);
    }).join('');
    const uniq = p => { const parts = p.split('/'); for (let k = 1; k <= parts.length; k++) { const suf = parts.slice(-k).join('/'); if (F.files.filter(g => g.path === suf || g.path.endsWith('/' + suf)).length === 1) return suf; } return p; };
    const files = F.files.map(f => {
      const tag = f.status === 'A' ? 'new' : f.status === 'D' ? 'del' : 'edit';
      const tc = tag === 'new' ? 'var(--ink)' : tag === 'del' ? 'var(--youInk)' : 'var(--inkSub)';
      return `<button type="button" class="file" data-tip="path:${esc(f.path)}" aria-label="${esc(f.path + ', ' + tag)}"><span class="nm">${esc(uniq(f.path))}</span><span class="tg" style="color:${tc}">${tag}</span></button>`;
    }).join('');
    const acts = F.decisions.map(a => {
      const [bg, fg, bd] = ACT[a.actor] || ACT.world;
      return `<div class="act"${a.source ? ` data-tip="src:${esc(a.source)}"` : ''}><span class="who" style="background:${bg};color:${fg};border-color:${bd}">${WHO[a.actor] || esc(a.actor)}</span><span class="tx">${esc(a.text)}<span class="tm">${esc(a.rel)}</span></span></div>`;
    }).join('');
    const ck = F.checks;
    const tiles = [[ck.graders, 'graders'], [String(ck.rounds), ck.rounds === 1 ? 'round' : 'rounds'], [String(ck.minutes), 'agent min'], [String(ck.calls), 'World calls']]
      .map(([v, k]) => `<div class="ctile"><div class="v">${esc(v)}</div><div class="k">${k}</div></div>`).join('');
    return bar + `<div class="fbody">${head}${view}<div class="cols">
      <div class="col"><div class="colh">Inputs</div>
        <div class="ingrid">
          <span class="k">Agent</span><div class="seg" role="group" aria-label="Agent">${agents}</div>
          <span class="k top">Task</span><div class="task"><div class="tl1">${esc(F.taskFile)}<span class="sub">minted for ${esc(F.ttl)}</span></div><div class="tid">${esc(F.task)}</div></div>
          <span class="k top">Grants</span><div class="seg" role="group" aria-label="Grants">${grants}</div>
        </div>
        <div class="sect"><div class="l">↑ Came from · ${F.readCalls} reads over ${F.readFiles} files</div><div class="chips">${came}</div></div>
      </div>
      <div class="col g6"><div class="colh">Outputs</div>${files}<div class="prpill">#${F.pr} · ${F.files.length} files</div>
        <div class="sect" style="margin-top:auto"><div class="l">↓ Led to</div><div class="chips">${F.mergedInto ? stateChip(F.mergedInto.id, '') : ''}</div></div>
      </div>
      <div class="col g9"><div class="colh">Decisions</div>${acts}</div>
      <div class="col g9"><div class="colh">Checks</div><div class="verdict">${esc(verdict(ck.graders, ck.rounds, N[F.id].firstRoundPass, ck.minutes))}</div><div class="tiles">${tiles}</div></div>
    </div></div>`;
  }

  function renderBrane(animate) {
    if (!$('bstage')) {
      $('br').innerHTML = `<div class="bstage" id="bstage"><article class="front" id="front"></article></div><nav class="rail" id="rail" aria-label="Brane timeline"></nav>`;
    }
    $('front').innerHTML = cardHTML();
    $('front').setAttribute('aria-label', cleanTitle(RAIL[st.ri].title) || RAIL[st.ri].label);
    placeBars(animate);
    const lane = N[RAIL[0].id].lane;
    $('rail').innerHTML = `<div class="rail-h">${glyphOf(lane, 11)}${esc(laneTitle(lane))}</div>` + RAIL.map((r, k) => ({ r, k })).reverse().map(({ r, k }) => {
      const n = N[r.id];
      return `<button type="button" class="rrow" data-k="${k}" data-tip="node:${r.id}" aria-current="${k === st.ri}"><span class="lb"><span class="rt">${esc(cleanTitle(r.title) || r.label)}</span><span class="rid">${esc(n.episode || r.label)}${r.pr && '#' + r.pr !== (n.episode || r.label) ? ' · #' + r.pr : ''}</span></span>${r.badge > 0 ? `<span class="nb inl">${r.badge}</span>` : ''}<span class="tm">${esc(r.rel)}</span></button>`;
    }).join('');
  }

  // ---------- Header ----------
  const addrOf = () => `${F.brane}@${st.ri + 1}`;
  function renderCrumbs() {
    const cs = [['world', st.d === 0]];
    if (st.d === 1) cs.push([addrOf(), true]);
    $('crumbs').innerHTML = cs.map(([t, cur], k) => (k ? '<span class="crumb-sep" aria-hidden="true">›</span>' : '') +
      `<button type="button" class="crumb" data-crumb="${k}" aria-current="${cur && k === cs.length - 1}">${esc(t)}</button>`).join('');
    $('addr').placeholder = st.d === 0 ? 'go to a place — row-189@6 · shared@0ffbb9c' : '';
  }
  let msgT;
  function say(t) { const m = $('hmsg'); m.textContent = t; m.style.opacity = '1'; clearTimeout(msgT); msgT = setTimeout(() => { m.style.opacity = '0'; }, 2600); }

  // ---------- Navigation between timeline and brane ----------
  function setOrigin(x, y) {
    const o = `${x * scale}px ${y * scale}px`;
    $('L0').style.transformOrigin = o; $('L1').style.transformOrigin = o;
  }
  function open(id) {
    if (isCluster(id)) { const d = DISP[id]; return zoomTo(Math.ceil(view.u - 0.01) - 1, PL + d.x - view.L); }
    const k = RAIL.findIndex(r => r.id === id);
    if (k < 0) { say(`Only ${nameOf(F.brane)}’s brane is in this data`); return; }
    const sx = screenX(id), d = DISP[REP[id]];
    if (st.d === 0) setOrigin(sx != null ? Math.max(PL, Math.min(PR, sx)) : PL, d ? d.y : LY[F.brane]);
    st.d = 1; st.ri = k; st.focus = id; $('tip').hidden = true;
    renderBrane(false);
    $('root').classList.add('zoomed');
    renderCrumbs(); $('root').focus({ preventScroll: true });
  }
  function out(focusId) {
    if (st.d === 0) return;
    const id = focusId || RAIL[st.ri].id;
    if (POS[id]) {
      st.focus = id;
      const x = screenX(id); if (x == null || x < PL + 40 || x > PR - 40) panTo(view.L + (x - (PL + PW / 2)), false);
      applyFocus(); setOrigin(screenX(id), DISP[REP[id]].y);
    }
    st.d = 0; $('root').classList.remove('zoomed'); $('tip').hidden = true;
    renderCrumbs(); $('root').focus({ preventScroll: true });
  }
  function setRail(k) {
    k = Math.max(0, Math.min(RAIL.length - 1, k));
    if (k === st.ri) return;
    st.ri = k; st.focus = RAIL[k].id;
    renderBrane(true); renderCrumbs();
  }
  function go(raw) {
    const t = raw.trim().replace(/^world\b\s*/, '');
    if (t === '' || t === 'back') return out();
    if (t === 'now') { out(); return goNow(); }
    let m = t.match(/^([a-z0-9-]+)@(\d{1,2})$/);
    if (m) {
      if (m[1] !== F.brane) return say(LANES[m[1]] ? `Only ${nameOf(F.brane)}’s brane is in this data` : `no brane called ${m[1]}`);
      const r = RAIL[+m[2] - 1]; return r ? (st.d === 1 ? setRail(+m[2] - 1) : open(r.id)) : say(`no state ${m[2]} on ${m[1]}`);
    }
    m = t.match(/^(?:shared@)?([0-9a-f]{4,40})$/);
    if (m) { const id = Object.keys(N).find(x => x.startsWith(m[1])); if (id && POS[id]) { if (st.d === 1) out(id); else { st.focus = id; applyFocus(); ensureVisible(id); } return; } }
    say('no such place: ' + t);
  }

  // ---------- Popovers ----------
  const KEYROWS = [['← →', 'back · forward a page', '—'], ['⇧← ⇧→', 'previous · next state', '—'], ['↑ ↓', 'lane above · below', 'back in time · toward now'], ['+ −', 'zoom in · out', '—'], ['End', 'back to now', 'most recent'], ['Home', 'the oldest state', 'earliest'], ['Enter', 'open brane', '—'], ['Esc', '—', 'back to timeline'], ['/', 'address bar', 'address bar'], [',', 'appearance', 'appearance'], ['?', 'these keys', 'these keys']];
  function renderKeys() {
    $('keys').innerHTML = `<div class="kgrid"><span></span><span style="font-size:12px;font-weight:700;color:${st.d === 0 ? 'var(--text)' : 'var(--faint)'}">Timeline</span><span style="font-size:12px;font-weight:700;color:${st.d === 1 ? 'var(--text)' : 'var(--faint)'}">Brane</span>` +
      KEYROWS.map(([k, a, b]) => `<span><span class="kbd">${k}</span></span><span>${a}</span><span>${b}</span>`).join('') +
      `</div><div style="margin-top:10px;font-size:12px;color:var(--sub)">Drag or scroll sideways to pan. ⌘-scroll or pinch to zoom around the pointer.</div>`;
  }
  const look = { theme: 'paper', motion: 'full' };
  try { Object.assign(look, JSON.parse(localStorage.getItem('world.appearance') || '{}')); } catch (_) {}
  function applyLook() {
    document.documentElement.dataset.theme = look.theme;
    $('root').classList.toggle('reduced', look.motion === 'reduced');
    try { localStorage.setItem('world.appearance', JSON.stringify(look)); } catch (_) {}
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
    const shared = n.kind === 'shared', l = LANES[n.lane];
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
      ${s ? `<div class="hs">${esc(s)}</div>` : `<div class="hs none">No agent episode in this data.</div>`}
      ${n.model ? `<div class="hm"><b>${esc(agentLine(n))}</b>${g ? ' · ' + g : ''}${n.minutes != null ? ` · ${n.minutes} agent-min` : ''}</div>` : ''}
      ${n.badge ? `<div class="hm"><span class="nb inl">${n.badge}</span> ${n.badge === 1 ? 'decision' : 'decisions'} of yours</div>` : ''}
      <div class="hw">${esc(n.rel)} · ${esc(exactOf(n.at))}</div><div class="hid">${esc(ids)}</div></div>`;
  }
  function tipHTML(spec) {
    const [kind, ...rest] = spec.split(':'), v = rest.join(':');
    if (kind === 'node') return nodeCard(v);
    if (kind === 'reads') {
      const r = READS.get(v), n = r.files.reduce((a, f) => a + f.count, 0), s = N[r.fromId];
      const head = s && s.title ? `<div class="th">${esc(cleanTitle(s.title))}</div><div style="color:var(--sub);margin-bottom:6px">${esc(r.label)} · ${esc(s.rel)} · read ${r.files.length} file${r.files.length === 1 ? '' : 's'}, ${n} time${n === 1 ? '' : 's'}</div>`
        : `<div class="th">${esc(r.label)} <span style="font-weight:500;color:var(--sub)">· ${r.files.length} file${r.files.length === 1 ? '' : 's'}, ${n} read${n === 1 ? '' : 's'}</span></div>`;
      return head + r.files.map(f => `<div class="tf"><span>${esc(f.path)}</span><span class="c">×${f.count}</span></div>`).join('');
    }
    if (kind === 'goal') return `<div class="hc"><div class="hk">${esc(laneTitle(v))} · why it exists</div><div class="ht" style="font-weight:500">${esc(LANES[v].goal)}</div></div>`;
    if (kind === 'ms') {
      const ms = LANES[v].milestones, rest2 = ms.names.filter(x => !ms.doneNames.includes(x) && x !== ms.inReview && x !== ms.next);
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
    if (n && st.d === 0 && !dragging && st.focus !== n.dataset.id) { st.focus = n.dataset.id; applyFocus(); }
  });
  let suppressClick = false;
  $('stage').addEventListener('click', e => {
    if (suppressClick) { suppressClick = false; e.stopPropagation(); return; }
    const n = e.target.closest('#tl .node');
    if (n) return open(n.dataset.id);
    if (e.target.closest('#quiet')) { const l = $('quiet-list'); l.hidden = !l.hidden; $('quiet').setAttribute('aria-expanded', String(!l.hidden)); return; }
    const lab = e.target.closest('.lhead');
    if (lab) return lab.dataset.lane === F.brane ? open(RAIL[RAIL.length - 1].id) : say(`Only ${nameOf(F.brane)}’s brane is in this data`);
    if (e.target.closest('#now-b')) return goNow();
    if (e.target.closest('#pg-back')) return panTo(view.L - PW * 0.8);
    if (e.target.closest('#pg-fwd')) return panTo(view.L + PW * 0.8);
    if (e.target.closest('#z-in')) return zoomStep(-1);
    if (e.target.closest('#z-out')) return zoomStep(1);
    const rr = e.target.closest('.rrow'); if (rr) return setRail(+rr.dataset.k);
    const bar = e.target.closest('.strip');
    if (bar) { const k = RAIL.findIndex(r => r.id === bar.dataset.id); return k >= 0 ? setRail(k) : POS[bar.dataset.id] ? out(bar.dataset.id) : say('That state is outside the timeline'); }
    const ch = e.target.closest('[data-go]');
    if (ch) { const id = ch.dataset.go, k = RAIL.findIndex(r => r.id === id); return k >= 0 ? setRail(k) : POS[id] ? out(id) : say(`${N[id] ? N[id].label : id.slice(0, 7)} is outside the timeline`); }
  }, true);

  // Drag to pan (the plot), drag/click the overview window.
  let drag = null;
  $('stage').addEventListener('pointerdown', e => {
    if (st.d !== 0 || e.button !== 0) return;
    const plot = e.target.closest('#plot'), mini = e.target.closest('#mini');
    if (!plot && !mini) return;
    cancelAnimationFrame(anim);
    if (mini) {
      const r = $('mini').getBoundingClientRect(), mx = (e.clientX - r.left) / scale, mw = $('mwin');
      const wl = parseFloat(mw.style.left), ww = parseFloat(mw.style.width);
      if (mx < wl || mx > wl + ww) {   // click outside the window: jump there, then keep dragging
        const f = fOfMX(mx); panTo(posF(f) - PW / 2);
        drag = { kind: 'mini', x0: e.clientX, f0: fOfMX(mx) , grab: 0.5, moved: true, jumped: true };
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
    if (st.d !== 0) return;
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

  $('crumbs').addEventListener('click', e => { const c = e.target.closest('[data-crumb]'); if (c && c.dataset.crumb === '0') out(); });
  $('keys-btn').addEventListener('click', () => toggle('keys'));
  $('look-btn').addEventListener('click', () => toggle('look'));
  $('look').addEventListener('click', e => { const b = e.target.closest('[data-look]'); if (!b) return; const [k, v] = b.dataset.look.split(':'); look[k] = v; applyLook(); renderLook(); });
  $('addr').addEventListener('keydown', e => {
    e.stopPropagation();
    if (e.key === 'Enter') { go(e.target.value); e.target.value = ''; $('root').focus({ preventScroll: true }); }
    if (e.key === 'Escape') { e.target.value = ''; $('root').focus({ preventScroll: true }); }
  });
  const onRow = row => Object.keys(DISP).filter(id => DISP[id].row === row).sort((a, b) => DISP[a].x - DISP[b].x);
  document.addEventListener('keydown', e => {
    if (e.target && e.target.tagName === 'INPUT') return;
    const k = e.key; let done = true;
    if (k === '?') toggle('keys');
    else if (k === ',') toggle('look');
    else if (k === '/') $('addr').focus();
    else if (k === 'Escape') { if (!$('keys').hidden || !$('look').hidden) { $('keys').hidden = true; $('look').hidden = true; } else out(); }
    else if (st.d === 1 && k === 'ArrowUp') setRail(st.ri - 1);
    else if (st.d === 1 && k === 'ArrowDown') setRail(st.ri + 1);
    else if (st.d === 1 && k === 'Home') setRail(0);
    else if (st.d === 1 && k === 'End') setRail(RAIL.length - 1);
    else if (st.d === 0 && k === 'Enter') open(REP[st.focus] && isCluster(REP[st.focus]) ? REP[st.focus] : st.focus);
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
  window.World = { open: label => { const r = RAIL.find(x => x.label === label); if (r) open(r.id); return !!r; }, out: () => out(), rail: k => setRail(k),
    focus: id => { st.focus = id; applyFocus(); }, zoom: (u, ax) => zoomTo(u, ax == null ? PL + PW / 2 : ax, false), pan: L => panTo(L, false), now: () => goNow(false),
    view: () => ({ u: view.u, L: view.L, level: levelName(), atNow: atNow() }), screenX: id => screenX(id) };

  async function boot() {
    applyLook();
    try { await document.fonts.ready; } catch (_) {}
    renderTimeline(); renderBrane(false); renderCrumbs(); fit();
    document.documentElement.dataset.ready = '1';
    if (location.hash === '#brane') setTimeout(() => window.World.open(F.episode), 300);
    $('root').focus({ preventScroll: true });
  }
  boot();
})();
