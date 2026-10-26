/* Server HTML owns the static pages; this script adds keys, the address bar and Appearance.
   The app shell (data-app) loads it for its seams only: window.WorldSurface. */
(() => {
  "use strict";
  const SHARED = "shared";
  const APP = "/surface", STATIC = "/surface/static";
  const LOOK_KEY = "world.appearance";
  const LOOKS = {theme: ["paper", "slate", "night", "system"], marks: ["colour", "shape"], motion: ["full", "reduced", "system"]};
  const DEFAULT_LOOK = {theme: "paper", marks: "colour", motion: "system"};
  const ARROWS = ["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight"];
  const ADDRESS_HELP = 'try "world", "world NAME", "world NAME@N", "world NAME@HASH" or "back"';

  // The lane a state is on, as M2 and M3 place it.
  function stateLane(s) { return s.kind === "step" && s.brane ? s.brane : SHARED; }
  function laneStates(index, brane) {
    const states = index && Array.isArray(index.states) ? index.states : [];
    return states.filter(s => stateLane(s) === brane)
      .map(s => ({s, t: Date.parse(s.at)}))
      .sort((a, b) => a.t !== b.t ? a.t - b.t : a.s.id < b.s.id ? -1 : a.s.id > b.s.id ? 1 : 0)
      .map(x => x.s);
  }
  function stateHref(brane, id, base) { return base + "?brane=" + encodeURIComponent(brane) + "&at=" + encodeURIComponent(id); }
  // A brane answers to its key, its title, and each hyphen-joined suffix of the title.
  function braneNames(b) {
    const names = [String(b.key || "").toLowerCase()];
    const parts = String(b.title || "").toLowerCase().split("-");
    for (let i = 0; i < parts.length; i++) {
      const name = parts.slice(i).join("-");
      if (name) names.push(name);
    }
    return names;
  }
  // base is the page the hrefs land on: the app (the default) or the static pages.
  function parseAddress(text, index, base) {
    base = base || APP;
    const words = String(text == null ? "" : text).trim().toLowerCase().split(/\s+/).filter(Boolean);
    if (words.length === 1 && (words[0] === "world" || words[0] === "back")) return {href: base};
    const m = words.length === 2 && words[0] === "world" ? /^([^@]+)(?:@(.*))?$/.exec(words[1]) : null;
    if (!m) return {error: "unknown address: " + ADDRESS_HELP};
    const name = m[1], at = m[2];
    const keys = [];
    if (name === SHARED || name === "dev") keys.push(SHARED);
    for (const b of index && Array.isArray(index.branes) ? index.branes : []) {
      if (braneNames(b).includes(name) && !keys.includes(b.key)) keys.push(b.key);
    }
    if (keys.length === 0) return {error: 'no brane is named "' + name + '"'};
    if (keys.length > 1) return {error: 'ambiguous name "' + name + '": ' + keys.join(", ")};
    const key = keys[0], lane = laneStates(index, key);
    if (lane.length === 0) return {error: key + " has no states"};
    if (at === undefined) return {href: stateHref(key, lane[lane.length - 1].id, base)};
    if (/^[0-9]+$/.test(at)) {
      const n = Number(at);
      if (n >= 1 && n <= lane.length) return {href: stateHref(key, lane[n - 1].id, base)};
      if (at.length < 4) return {error: key + " has states 1 to " + lane.length + ", not " + at};
    }
    if (!/^[0-9a-f]{4,}$/.test(at)) return {error: '"@' + at + '" is neither a state number nor 4 or more hex characters'};
    const hits = lane.filter(s => String(s.id).toLowerCase().startsWith(at));
    if (hits.length === 1) return {href: stateHref(key, hits[0].id, base)};
    if (hits.length > 1) return {error: '"' + at + '" matches ' + hits.length + " states on " + key + ": type more of the hash"};
    return {error: "no state on " + key + ' starts with "' + at + '"'};
  }

  function timelineKey(key, model, focus) {
    const lanes = Array.isArray(model.lanes) ? model.lanes : [];
    const all = lanes.flatMap(l => l.states || []);
    if (all.length === 0) return null;
    const newest = all.reduce((a, b) => b.rank > a.rank ? b : a), oldest = all.reduce((a, b) => b.rank < a.rank ? b : a);
    if (key === "End") return {focus: newest.id};
    if (key === "Home") return {focus: oldest.id};
    let li = -1, si = -1;
    lanes.forEach((l, i) => (l.states || []).forEach((s, j) => { if (focus != null && s.id === focus) { li = i; si = j; } }));
    if (li < 0) return ARROWS.includes(key) ? {focus: newest.id} : null;
    const lane = lanes[li].states, here = lane[si];
    switch (key) {
      case "Enter": return here.href ? {go: here.href} : null;
      case "ArrowLeft": return si > 0 ? {focus: lane[si - 1].id} : null;
      case "ArrowRight": return si + 1 < lane.length ? {focus: lane[si + 1].id} : null;
      case "ArrowUp": case "ArrowDown": {
        const step = key === "ArrowUp" ? -1 : 1;
        for (let i = li + step; i >= 0 && i < lanes.length; i += step) {
          const states = lanes[i].states || [];
          if (states.length === 0) continue;
          // Nearest rank; on a tie the earlier state, which comes first.
          let best = states[0];
          for (const s of states) if (Math.abs(s.rank - here.rank) < Math.abs(best.rank - here.rank)) best = s;
          return {focus: best.id};
        }
        return null;
      }
    }
    return null;
  }
  function stackKey(key, model) {
    const rail = Array.isArray(model.rail) ? model.rail : [];
    const pos = rail.findIndex(s => s.id === model.focus);
    const go = s => s && s.href ? {go: s.href} : null;
    const other = cols => (Array.isArray(cols) ? cols : []).find(c => c.brane !== model.brane && c.cards && c.cards.length > 0);
    switch (key) {
      case "Escape": return {go: STATIC};
      case "End": return go(rail[rail.length - 1]);
      case "Home": return go(rail[0]);
      case "ArrowUp": return pos > 0 ? go(rail[pos - 1]) : null;
      case "ArrowDown": {
        if (pos < 0) return null;
        if (pos + 1 < rail.length) return go(rail[pos + 1]);
        const c = other(model.later); // the end of the brane follows on into what it led to
        return c ? go(c.cards[0]) : null;
      }
      case "ArrowLeft": { const c = other(model.earlier); return c ? go(c.cards[c.cards.length - 1]) : null; }
      case "ArrowRight": { const c = other(model.later); return c ? go(c.cards[0]) : null; }
    }
    return null;
  }
  function keyAction(key, model, focus) {
    switch (key) {
      case "/": return {ui: "address"};
      case "?": return {ui: "keys"};
      case ",": return {ui: "appearance"};
      case "n": case "N": return {ui: "needs"};
    }
    if (!model) return null;
    if (model.view === "timeline") return timelineKey(key, model, focus);
    if (model.view === "stack") return stackKey(key, model);
    return null;
  }

  function shouldHandle(target) {
    if (!target) return true;
    const tag = String(target.tagName || "").toUpperCase();
    if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return false;
    if (target.isContentEditable === true) return false;
    if (target.contentEditable === "true" || target.contentEditable === "plaintext-only") return false;
    const editable = typeof target.getAttribute === "function" ? target.getAttribute("contenteditable") : null;
    return editable == null || String(editable).toLowerCase() === "false";
  }

  function attr(node, name) { const v = node.getAttribute(name); return v == null ? "" : String(v); }
  // A state's x: a circle's cx, or the middle of a polygon's x extent (every glyph is symmetric).
  function centreX(a) {
    const circle = a.querySelector("circle");
    if (circle) return Number(circle.getAttribute("cx"));
    const polygon = a.querySelector("polygon");
    if (!polygon) return NaN;
    const xs = attr(polygon, "points").trim().split(/\s+/).map(p => Number(p.split(",")[0]));
    return (Math.min(...xs) + Math.max(...xs)) / 2;
  }
  function buildModel(doc) {
    const front = doc.querySelector("article.front");
    if (front || doc.querySelector("nav.rail") || doc.querySelector("main a.back")) {
      const card = a => ({id: attr(a, "data-state"), href: attr(a, "href")});
      const columns = section => Array.from(doc.querySelectorAll(section + " div.column"))
        .map(c => ({brane: attr(c, "data-brane"), cards: Array.from(c.querySelectorAll("a.card")).map(card)}));
      return {view: "stack", brane: front ? attr(front, "data-brane") : "", focus: front ? attr(front, "data-state") : "",
        rail: Array.from(doc.querySelectorAll("nav.rail a")).map(card),
        earlier: columns("section.earlier"), later: columns("section.later")};
    }
    const lanes = Array.from(doc.querySelectorAll("g.lane")).map(g => ({brane: attr(g, "data-brane"), states: []}));
    const byBrane = new Map(lanes.map(l => [l.brane, l]));
    const states = Array.from(doc.querySelectorAll("a.state")).map((a, order) => {
      const x = centreX(a);
      return {a, order, x: Number.isFinite(x) ? x : Infinity};
    });
    states.slice().sort((p, q) => p.x !== q.x ? p.x - q.x : p.order - q.order).forEach((s, rank) => { s.rank = rank; });
    for (const s of states) {
      const lane = byBrane.get(attr(s.a, "data-brane"));
      if (lane) lane.states.push({id: attr(s.a, "data-state"), href: attr(s.a, "href"), rank: s.rank});
    }
    for (const lane of lanes) lane.states.sort((p, q) => p.rank - q.rank);
    return {view: "timeline", lanes};
  }

  function cleanLook(raw) {
    const look = Object.assign({}, DEFAULT_LOOK);
    if (raw && typeof raw === "object") {
      for (const k of Object.keys(LOOKS)) if (LOOKS[k].includes(raw[k])) look[k] = raw[k];
    }
    return look;
  }
  function loadLook(storage) {
    try {
      const raw = storage ? storage.getItem(LOOK_KEY) : null;
      return cleanLook(typeof raw === "string" ? JSON.parse(raw) : null);
    } catch (error) { return cleanLook(null); }
  }
  function saveLook(storage, look) {
    try { if (storage) storage.setItem(LOOK_KEY, JSON.stringify(cleanLook(look))); } catch (error) { /* this browser keeps no look */ }
  }
  function resolveLook(look, prefersDark, prefersReduced) {
    const l = cleanLook(look);
    return {
      "data-theme": l.theme === "system" ? (prefersDark ? "slate" : "paper") : l.theme,
      "data-marks": l.marks,
      "data-motion": l.motion === "system" ? (prefersReduced ? "reduced" : "full") : l.motion,
    };
  }

  if (typeof module !== "undefined") module.exports = {laneStates, parseAddress, keyAction, shouldHandle, buildModel, loadLook, saveLook, resolveLook};
  if (typeof window !== "undefined") window.WorldSurface = {laneStates, parseAddress, keyAction, shouldHandle, buildModel, loadLook, saveLook, resolveLook};
  if (typeof document === "undefined") return;
  // The app shell has its own keys, address bar and Appearance (app.js); it borrows the seams only.
  if (document.documentElement && document.documentElement.getAttribute("data-app") != null) return;

  // Appearance first, so a stored look replaces the default as early as defer allows.
  const root = document.documentElement;
  const media = q => typeof matchMedia === "function" ? matchMedia(q) : null;
  const dark = media("(prefers-color-scheme: dark)"), reduced = media("(prefers-reduced-motion: reduce)");
  let storage = null;
  try { storage = window.localStorage; } catch (error) { storage = null; }
  let look = loadLook(storage);
  function applyLook() {
    const attrs = resolveLook(look, !!(dark && dark.matches), !!(reduced && reduced.matches));
    for (const name of Object.keys(attrs)) root.setAttribute(name, attrs[name]);
  }
  applyLook();
  for (const m of [dark, reduced]) if (m && typeof m.addEventListener === "function") m.addEventListener("change", applyLook);

  const model = buildModel(document);
  const KEYS = {
    timeline: [["↑ ↓", "lane above · below"], ["← →", "earlier · later on the lane"], ["End · Home", "most recent · earliest"],
      ["Enter", "open the brane"], ["N", "jump to what needs you"], ["/", "the address bar"], [",", "appearance"], ["?", "show these keys"]],
    stack: [["↑ ↓", "back in time · toward now"], ["← →", "the brane behind · ahead"], ["End · Home", "most recent · earliest"],
      ["Esc", "back to the timeline"], ["N", "jump to what needs you"], ["/", "the address bar"], [",", "appearance"], ["?", "show these keys"]],
  };
  const CHOICES = [
    ["theme", "Theme", [["paper", "Paper"], ["slate", "Slate"], ["night", "Night"], ["system", "Match system"]]],
    ["marks", "Brane marks", [["colour", "Shape and colour"], ["shape", "Shape"]]],
    ["motion", "Motion", [["full", "Full"], ["reduced", "Reduced"], ["system", "System setting"]]],
  ];
  function el(tag, attrs, content) {
    const node = document.createElement(tag);
    for (const name of Object.keys(attrs || {})) node.setAttribute(name, attrs[name]);
    if (typeof content === "string") node.textContent = content;
    else for (const child of content || []) node.appendChild(child);
    return node;
  }

  let panel = null, panelKind = "", returnTo = null;
  function closePanel() {
    if (!panel) return;
    panel.remove(); panel = null; panelKind = "";
    if (returnTo && typeof returnTo.focus === "function") returnTo.focus();
    returnTo = null;
  }
  function openPanel(kind) {
    const was = panelKind;
    closePanel();
    if (was === kind) return;
    returnTo = document.activeElement;
    const close = el("button", {type: "button", class: "close", "aria-label": "Close"}, "×");
    close.addEventListener("click", closePanel);
    let body = [];
    if (kind === "keys") {
      body = [el("h2", {}, "Keys"), el("dl", {}, (KEYS[model.view] || KEYS.timeline).flatMap(([k, what]) => [el("dt", {}, k), el("dd", {}, what)]))];
    } else if (kind === "needs") {
      body = [el("h2", {}, "Needs you"), el("p", {}, "Nothing needs you")];
    } else {
      body = [el("h2", {}, "Appearance")];
      for (const [name, legend, options] of CHOICES) {
        body.push(el("fieldset", {}, [el("legend", {}, legend)].concat(options.map(([value, label]) => {
          const input = el("input", {type: "radio", name: "look-" + name, value});
          input.checked = look[name] === value;
          input.addEventListener("change", () => {
            look = cleanLook(Object.assign({}, look, {[name]: value}));
            saveLook(storage, look); applyLook();
          });
          return el("label", {}, [input, el("span", {}, " " + label)]);
        }))));
      }
      body.push(el("p", {class: "sunholo"}, "Sunholo"));
    }
    const titles = {keys: "Keys", needs: "Needs you", appearance: "Appearance"};
    panel = el("div", {class: "panel", role: "dialog", "aria-label": titles[kind] || kind, "data-panel": kind}, [close].concat(body));
    panelKind = kind;
    document.body.appendChild(panel);
    const first = panel.querySelector("input:checked") || close;
    first.focus();
  }

  const bar = document.querySelector("input.address"), note = document.querySelector(".address-note");
  let index = null, lastState = null;
  function loadIndex() {
    if (!index) {
      index = fetch("/surface/index.json", {headers: {Accept: "application/json"}})
        .then(r => { if (!r.ok) throw new Error("HTTP " + r.status); return r.json(); });
      index.catch(() => { index = null; });
    }
    return index;
  }
  function say(text) {
    if (note) note.textContent = text;
    if (bar) { if (text) bar.setAttribute("aria-invalid", "true"); else bar.removeAttribute("aria-invalid"); }
  }
  async function runAddress() {
    let result = parseAddress(bar.value, null, STATIC);
    if (!result.href) {
      try { result = parseAddress(bar.value, await loadIndex(), STATIC); } catch (error) { result = {error: "the surface index could not be read: " + error.message}; }
    }
    if (result.href) location.assign(result.href);
    else say(result.error); // an error never navigates
  }
  if (bar) {
    bar.addEventListener("keydown", e => {
      if (e.key === "Enter") { e.preventDefault(); runAddress(); }
      else if (e.key === "Escape") {
        e.preventDefault(); bar.blur();
        if (lastState && typeof lastState.focus === "function") lastState.focus();
      }
    });
    bar.addEventListener("input", () => say(""));
  }
  const gear = document.querySelector("button.gear");
  if (gear) { gear.hidden = false; gear.addEventListener("click", () => openPanel("appearance")); }

  const isState = node => !!node && typeof node.matches === "function" && node.matches("a.state");
  document.addEventListener("focusin", e => { if (isState(e.target)) lastState = e.target; });
  document.addEventListener("keydown", e => {
    if (e.defaultPrevented || e.ctrlKey || e.metaKey || e.altKey) return;
    if (panel && e.key === "Escape") { e.preventDefault(); closePanel(); return; }
    if (!shouldHandle(e.target)) return;
    const active = document.activeElement;
    const action = keyAction(e.key, model, isState(active) ? active.getAttribute("data-state") : null);
    if (!action) return;
    e.preventDefault();
    if (action.focus) {
      const link = Array.from(document.querySelectorAll("a.state")).find(a => a.getAttribute("data-state") === action.focus);
      if (link) { lastState = link; link.focus(); }
    } else if (action.go) location.assign(action.go);
    else if (action.ui === "address") { if (bar) { bar.focus(); bar.select(); } }
    else openPanel(action.ui);
  });
})();
