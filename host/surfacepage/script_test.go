package surfacepage

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestSurfaceScriptAssetSeams(t *testing.T) {
	if len(Script) == 0 {
		t.Fatal("surface.js asset empty")
	}
	for _, name := range []string{"laneStates", "parseAddress", "keyAction", "shouldHandle", "buildModel", "loadLook", "saveLook", "resolveLook"} {
		if !strings.Contains(string(Script), name) {
			t.Fatalf("missing seam %s", name)
		}
	}
	if !strings.Contains(string(Script), `if (typeof module !== "undefined") module.exports = {`) {
		t.Fatal("seams are not exported for node")
	}
}

type card struct {
	ID   string `json:"id"`
	Href string `json:"href"`
}

type column struct {
	Brane string `json:"brane"`
	Cards []card `json:"cards"`
}

// stackModel is the model buildModel must read off a rendered stack page,
// built here from the layout the page was rendered from.
func stackModel(s Stack) map[string]any {
	cards := func(cs []StackCard) []card {
		out := []card{}
		for _, c := range cs {
			out = append(out, card{c.ID, c.Href})
		}
		return out
	}
	cols := func(cs []StackColumn) []column {
		out := []column{}
		for _, c := range cs {
			out = append(out, column{c.Brane, cards(c.Cards)})
		}
		return out
	}
	return map[string]any{"view": "stack", "brane": s.Front.Brane, "focus": s.Front.ID,
		"rail": cards(s.Rail), "earlier": cols(s.Earlier), "later": cols(s.Later)}
}

func timelineModel(tl Timeline) map[string]any {
	type state struct {
		ID   string `json:"id"`
		Href string `json:"href"`
		Rank int    `json:"rank"`
	}
	type lane struct {
		Brane  string  `json:"brane"`
		States []state `json:"states"`
	}
	lanes := []lane{}
	for _, l := range tl.Lanes {
		nodes := append([]Node(nil), l.Nodes...)
		sort.Slice(nodes, func(i, j int) bool { return nodes[i].Rank < nodes[j].Rank })
		ln := lane{Brane: l.Key, States: []state{}}
		for _, n := range nodes {
			ln.States = append(ln.States, state{n.ID, n.Href, n.Rank})
		}
		lanes = append(lanes, ln)
	}
	return map[string]any{"view": "timeline", "lanes": lanes}
}

// surfaceFixtures renders the rig's timeline and every state's stack, with
// the models and lanes Go computes for them.
func surfaceFixtures(t *testing.T) []byte {
	t.Helper()
	ix := rigIndex(t)
	raw, err := os.ReadFile(filepath.Join("..", "surface", "testdata", "index-rig.json"))
	if err != nil {
		t.Fatal(err)
	}
	tl := LayoutTimeline(ix, fixedNow)
	page, err := RenderTimeline(tl, rigNow)
	if err != nil {
		t.Fatal(err)
	}
	type stackFixture struct {
		HTML    string         `json:"html"`
		Model   map[string]any `json:"model"`
		Address string         `json:"address"`
		Href    string         `json:"href"`
	}
	var stacks []stackFixture
	lanes := map[string][]string{}
	for _, x := range ix.States {
		s, err := LayoutStack(ix, stateBrane(x), x.ID, fixedNow)
		if err != nil {
			t.Fatal(err)
		}
		body, err := RenderStack(s)
		if err != nil {
			t.Fatal(err)
		}
		stacks = append(stacks, stackFixture{string(body), stackModel(s), StackAddress(s), s.Front.Href})
		var ids []string
		for _, c := range s.Rail {
			ids = append(ids, c.ID)
		}
		lanes[s.Lane.Key] = ids
	}
	out, err := json.Marshal(map[string]any{
		"index":    json.RawMessage(raw),
		"timeline": map[string]any{"html": string(page), "model": timelineModel(tl)},
		"stacks":   stacks,
		"lanes":    lanes,
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TestSurfaceScriptPure runs surface.js's seams and its DOM wiring under
// node. It skips only when no node exists at all; with WORLD_EXEC_NODE set,
// a harness that cannot run fails.
func TestSurfaceScriptPure(t *testing.T) {
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	dir := t.TempDir()
	for name, data := range map[string][]byte{"surface.js": Script, "fixtures.json": surfaceFixtures(t), "harness.js": []byte(surfaceHarness)} {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, filepath.Join(dir, "harness.js"))
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("node harness: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "SURFACE HARNESS PASS") {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

const surfaceHarness = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const api = require('./surface.js');
const fx = JSON.parse(fs.readFileSync('fixtures.json', 'utf8'));
const index = fx.index;
assert.deepEqual(Object.keys(api).sort(), ['buildModel','keyAction','laneStates','loadLook','parseAddress','resolveLook','saveLook','shouldHandle']);

// ---- A tiny DOM: enough of the server's own HTML for buildModel and the wiring.
const VOID = new Set(['area','base','br','col','embed','hr','img','input','link','meta','source','track','wbr']);
function decode(s) {
  return s.replace(/&(#x[0-9a-f]+|#[0-9]+|amp|lt|gt|quot|apos);/gi, (m, e) => {
    const k = e.toLowerCase();
    if (k === 'amp') return '&'; if (k === 'lt') return '<'; if (k === 'gt') return '>'; if (k === 'quot') return '"'; if (k === 'apos') return "'";
    return String.fromCodePoint(k[1] === 'x' ? parseInt(k.slice(2), 16) : parseInt(k.slice(1), 10));
  });
}
function parseSelector(sel) {
  return sel.split(',').map(part => part.trim().split(/\s+/).map(compound => {
    const m = /^([a-z0-9]+)?((?:\.[a-z0-9_-]+)*)((?:\[[a-z-]+(?:="[^"]*")?\])*)(:checked)?$/i.exec(compound);
    if (!m) throw new Error('unsupported selector ' + sel);
    const attrs = [...(m[3] || '').matchAll(/\[([a-z-]+)(?:="([^"]*)")?\]/gi)].map(a => [a[1], a[2]]);
    return {tag: m[1] && m[1].toLowerCase(), classes: (m[2] || '').split('.').filter(Boolean), attrs, checked: !!m[4]};
  }));
}
function matchCompound(n, c) {
  if (!n.attrs) return false;
  if (c.tag && n.localName !== c.tag) return false;
  const cls = (n.attrs.class || '').split(/\s+/);
  if (!c.classes.every(k => cls.includes(k))) return false;
  if (!c.attrs.every(([k, v]) => v === undefined ? k in n.attrs : n.attrs[k] === v)) return false;
  return !c.checked || n.checked === true;
}
function matchChain(n, chain) {
  if (!matchCompound(n, chain[chain.length - 1])) return false;
  let i = chain.length - 2, p = n.parent;
  while (i >= 0 && p) { if (matchCompound(p, chain[i])) i--; p = p.parent; }
  return i < 0;
}
class El {
  constructor(tag, attrs, doc) { this.localName = tag; this.tagName = tag.toUpperCase(); this.attrs = attrs || {}; this.children = []; this.parent = null; this.listeners = {}; this.ownerDocument = doc; this.textContent = ''; }
  getAttribute(n) { return Object.prototype.hasOwnProperty.call(this.attrs, n) ? this.attrs[n] : null; }
  setAttribute(n, v) { this.attrs[n] = String(v); }
  removeAttribute(n) { delete this.attrs[n]; }
  get hidden() { return 'hidden' in this.attrs; }
  set hidden(v) { if (v) this.attrs.hidden = ''; else delete this.attrs.hidden; }
  get value() { return this._value !== undefined ? this._value : (this.attrs.value || ''); }
  set value(v) { this._value = String(v); }
  appendChild(c) { c.parent = this; this.children.push(c); return c; }
  remove() { if (this.parent) { this.parent.children = this.parent.children.filter(c => c !== this); this.parent = null; } }
  addEventListener(type, fn) { (this.listeners[type] = this.listeners[type] || []).push(fn); }
  fire(type, event) { for (const fn of this.listeners[type] || []) fn(event); }
  all() { const out = []; const walk = n => { for (const c of n.children) { out.push(c); walk(c); } }; walk(this); return out; }
  matches(sel) { return parseSelector(sel).some(chain => matchChain(this, chain)); }
  querySelectorAll(sel) { const chains = parseSelector(sel); return this.all().filter(n => chains.some(c => matchChain(n, c))); }
  querySelector(sel) { return this.querySelectorAll(sel)[0] || null; }
  focus() { this.ownerDocument.activeElement = this; this.ownerDocument.fire('focusin', {target: this}); }
  blur() { if (this.ownerDocument.activeElement === this) this.ownerDocument.activeElement = this.ownerDocument.body; }
  select() {}
  text() { return this.textContent + this.children.map(c => c.text()).join(''); }
}
function parseHTML(html) {
  const doc = new El('#document', null, null); doc.attrs = null; doc.ownerDocument = doc;
  doc.createElement = tag => new El(tag, {}, doc);
  let cur = doc, m;
  const re = /<!--[\s\S]*?-->|<!doctype[^>]*>|<\/([a-zA-Z0-9]+)\s*>|<([a-zA-Z][a-zA-Z0-9]*)((?:\s+[^\s"'>\/=]+(?:\s*=\s*(?:"[^"]*"|'[^']*'|[^\s"'>]+))?)*)\s*(\/?)>/gi;
  while ((m = re.exec(html))) {
    if (m[1]) { const tag = m[1].toLowerCase(); let n = cur; while (n && n.localName !== tag) n = n.parent; if (n) cur = n.parent; continue; }
    if (!m[2]) continue;
    const tag = m[2].toLowerCase(), attrs = {};
    for (const a of (m[3] || '').matchAll(/([^\s"'>\/=]+)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s"'>]+)))?/g)) attrs[a[1].toLowerCase()] = decode(a[2] ?? a[3] ?? a[4] ?? '');
    const node = cur.appendChild(new El(tag, attrs, doc));
    if (m[4] === '/' || VOID.has(tag)) continue;
    if (tag === 'style' || tag === 'script') { const end = html.toLowerCase().indexOf('</' + tag, re.lastIndex); re.lastIndex = end < 0 ? html.length : end; }
    cur = node;
  }
  doc.documentElement = doc.querySelector('html'); doc.body = doc.querySelector('body'); doc.activeElement = doc.body;
  return doc;
}
const plain = v => JSON.parse(JSON.stringify(v));

// ---- buildModel on the rendered fixture pages equals the model Go built from their layouts.
const timelineDoc = parseHTML(fx.timeline.html);
const timeline = api.buildModel(timelineDoc);
assert.deepEqual(plain(timeline), fx.timeline.model, 'timeline model');
assert.equal(timeline.lanes.flatMap(l => l.states).length, 48);
for (const s of fx.stacks) assert.deepEqual(plain(api.buildModel(parseHTML(s.html))), s.model, 'stack model ' + s.model.focus);
assert(fx.stacks.some(s => s.model.rail.some(c => c.href.includes('&at='))), 'hrefs are decoded');

// ---- laneStates: the same lanes LayoutStack builds, for every brane.
for (const [key, ids] of Object.entries(fx.lanes)) assert.deepEqual(api.laneStates(index, key).map(s => s.id), ids, 'lane ' + key);
const tiny = {branes: [{key: 'b'}], states: [
  {id: 'c2', kind: 'shared', at: '2026-10-20T09:00:00Z', brane: ''},
  {id: 'c1', kind: 'shared', at: '2026-10-20T10:00:00+02:00', brane: ''},
  {id: 'n0', kind: 'step', at: '2026-10-20T09:00:00Z', brane: ''},
  {id: 'b1', kind: 'step', at: '2026-10-20T07:00:00Z', brane: 'b'},
  {id: 'x0', kind: 'shared', at: '2026-10-20T06:00:00Z', brane: 'b'}]};
assert.deepEqual(api.laneStates(tiny, 'shared').map(s => s.id), ['x0', 'c1', 'c2', 'n0'], 'by at, then id; shared states and brane-less steps');
assert.deepEqual(api.laneStates(tiny, 'b').map(s => s.id), ['b1']);
assert.deepEqual(api.laneStates(tiny, 'none'), []);
assert.deepEqual(api.laneStates(null, 'shared'), []);

// ---- parseAddress.
const lane189 = fx.lanes['row-189'], shared = fx.lanes['shared'];
const href = (k, id) => '/surface?brane=' + k + '&at=' + id;
const ok = (text, want, ix = index) => assert.deepEqual(api.parseAddress(text, ix), {href: want}, text);
const bad = (text, pattern, ix = index) => { const r = api.parseAddress(text, ix); assert.equal(r.href, undefined, text); assert.match(r.error, pattern, text); };
ok('world', '/surface'); ok('  WORLD  ', '/surface'); ok('back', '/surface'); ok('world', '/surface', null);
for (const name of ['row-189', 'w-world-onboarding', 'world-onboarding', 'onboarding', 'ONBOARDING'])
  ok('world ' + name, href('row-189', lane189[lane189.length - 1]));
ok('world row-139', href('row-139', fx.lanes['row-139'][fx.lanes['row-139'].length - 1]));
ok('world report', href('row-198', fx.lanes['row-198'][fx.lanes['row-198'].length - 1]));
ok('world shared', href('shared', shared[shared.length - 1])); ok('world dev', href('shared', shared[shared.length - 1]));
ok('world onboarding@1', href('row-189', lane189[0]));
ok('world onboarding@' + lane189.length, href('row-189', lane189[lane189.length - 1]));
ok('world row-189@f36428a', href('row-189', 'f36428abd9a72e5093b61cba14e716c37019f1ef'));
ok('world row-189@F36428ABD9A72E5093B61CBA14E716C37019F1EF', href('row-189', 'f36428abd9a72e5093b61cba14e716c37019f1ef'));
ok('world shared@0ee1945', href('shared', '0ee1945ffba795812f06dce95998b646adfac343'));
bad('world onboarding@0', /states 1 to/);
bad('world onboarding@' + (lane189.length + 1), /states 1 to/);
bad('world onboarding@1000', /no state on row-189/);
bad('world row-189@0ee1945', /no state on row-189/, index);
bad('world row-189@f36', /neither/);
bad('world row-189@zzzz', /neither/);
bad('world nosuch', /no brane/); bad('world onboarding', /no brane/, null); bad('world onboarding', /no brane/, {schema: 'world/surface-index/v1', empty: true});
bad('hello', /unknown address/); bad('', /unknown address/); bad('world a b', /unknown address/); bad('back there', /unknown address/); bad('world @1', /unknown address/);
const twin = {branes: [{key: 'row-1', title: 'w-a-report'}, {key: 'row-2', title: 'w-b-report'}, {key: 'row-3', title: 'w-shared'}], states: [
  {id: 'abcd1111', kind: 'step', brane: 'row-1', at: '2026-10-20T01:00:00Z'},
  {id: 'abcd2222', kind: 'step', brane: 'row-1', at: '2026-10-20T02:00:00Z'},
  {id: 'eeee0000', kind: 'shared', brane: '', at: '2026-10-20T03:00:00Z'}]};
bad('world report', /ambiguous name "report": row-1, row-2/, twin);
bad('world shared', /ambiguous/, twin);
bad('world row-1@abcd', /matches 2 states/, twin);
ok('world row-1@abcd2', href('row-1', 'abcd2222'), twin);
ok('world a-report', href('row-1', 'abcd2222'), twin);
bad('world row-2', /row-2 has no states/, twin);
// Every stack's own address bar text opens that stack.
for (const s of fx.stacks) ok(s.address, s.href);

// ---- keyAction, timeline.
const S = (id, rank) => ({id, href: '/h/' + id, rank});
const tm = {view: 'timeline', lanes: [
  {brane: 'shared', states: [S('s0', 0), S('s1', 2), S('s2', 6)]},
  {brane: 'x', states: [S('x0', 1), S('x1', 4)]},
  {brane: 'y', states: []},
  {brane: 'z', states: [S('z0', 3), S('z1', 5)]}]};
const k = (key, model, focus) => api.keyAction(key, model, focus);
assert.deepEqual(k('ArrowUp', tm, 'x0'), {focus: 's0'}, 'tie goes to the earlier state');
assert.deepEqual(k('ArrowUp', tm, 'x1'), {focus: 's1'}, 'tie goes to the earlier state');
assert.deepEqual(k('ArrowDown', tm, 'x1'), {focus: 'z0'}, 'skips an empty lane; tie earlier');
assert.deepEqual(k('ArrowDown', tm, 's2'), {focus: 'x1'}, 'nearest rank');
assert.deepEqual(k('ArrowUp', tm, 'z1'), {focus: 'x1'});
assert.equal(k('ArrowUp', tm, 's0'), null); assert.equal(k('ArrowDown', tm, 'z1'), null);
assert.deepEqual(k('ArrowRight', tm, 'x0'), {focus: 'x1'}); assert.equal(k('ArrowRight', tm, 'x1'), null);
assert.deepEqual(k('ArrowLeft', tm, 'x1'), {focus: 'x0'}); assert.equal(k('ArrowLeft', tm, 'x0'), null);
assert.deepEqual(k('End', tm, 'x0'), {focus: 's2'}); assert.deepEqual(k('Home', tm, null), {focus: 's0'});
for (const a of ['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight']) {
  assert.deepEqual(k(a, tm, null), {focus: 's2'}, 'no focus: ' + a); assert.deepEqual(k(a, tm, 'gone'), {focus: 's2'}, 'unknown focus: ' + a);
}
assert.deepEqual(k('Enter', tm, 'z1'), {go: '/h/z1'}); assert.equal(k('Enter', tm, null), null);
assert.equal(k('Escape', tm, 'z1'), null); assert.equal(k('a', tm, 'z1'), null); assert.equal(k('Tab', tm, null), null);
assert.equal(k('ArrowUp', {view: 'timeline', lanes: []}, null), null);
// On the real page: End is the newest state, and every step along a lane is its neighbour.
const ranked = timeline.lanes.flatMap(l => l.states).sort((a, b) => a.rank - b.rank);
assert.deepEqual(k('End', timeline, null), {focus: ranked[ranked.length - 1].id}); assert.deepEqual(k('Home', timeline, null), {focus: ranked[0].id});
for (const l of timeline.lanes) for (let i = 1; i < l.states.length; i++) assert.deepEqual(k('ArrowLeft', timeline, l.states[i].id), {focus: l.states[i - 1].id});
assert.deepEqual(k('Enter', timeline, 'f36428abd9a72e5093b61cba14e716c37019f1ef'), {go: href('row-189', 'f36428abd9a72e5093b61cba14e716c37019f1ef')});

// ---- keyAction, stack.
const C = id => ({id, href: '/h/' + id});
const sm = {view: 'stack', brane: 'b', focus: 'b1', rail: [C('b0'), C('b1'), C('b2')],
  earlier: [{brane: 'b', cards: [C('b0')]}, {brane: 'shared', cards: [C('e0'), C('e1')]}, {brane: 'c', cards: [C('e2')]}],
  later: [{brane: 'b', cards: [C('b2')]}, {brane: 'shared', cards: [C('l0'), C('l1')]}, {brane: 'c', cards: [C('c0')]}]};
const at = focus => Object.assign({}, sm, {focus});
assert.deepEqual(k('ArrowUp', sm, null), {go: '/h/b0'}); assert.deepEqual(k('ArrowDown', sm, null), {go: '/h/b2'});
assert.deepEqual(k('ArrowUp', sm, 'b2'), {go: '/h/b0'}, 'the focus argument is ignored on the stack');
assert.equal(k('ArrowUp', at('b0'), null), null, 'the start of the rail');
assert.deepEqual(k('ArrowDown', at('b2'), null), {go: '/h/l0'}, 'the end of the rail follows on');
assert.deepEqual(k('ArrowLeft', sm, null), {go: '/h/e1'}); assert.deepEqual(k('ArrowRight', sm, null), {go: '/h/l0'});
assert.deepEqual(k('End', sm, null), {go: '/h/b2'}); assert.deepEqual(k('Home', sm, null), {go: '/h/b0'});
assert.deepEqual(k('Escape', sm, null), {go: '/surface'}); assert.equal(k('Enter', sm, null), null); assert.equal(k('x', sm, null), null);
const lone = {view: 'stack', brane: 'b', focus: 'b1', rail: [C('b0'), C('b1')], earlier: [{brane: 'b', cards: [C('b0')]}], later: []};
assert.equal(k('ArrowDown', lone, null), null, 'no other column: the end of the rail stays');
assert.equal(k('ArrowLeft', lone, null), null); assert.equal(k('ArrowRight', lone, null), null);
assert.deepEqual(k('ArrowUp', lone, null), {go: '/h/b0'});
// On the real pages: w189d's ↓ is w189e, and w189e (the brane's newest) follows on into the shared world.
const real = id => fx.stacks.find(s => s.model.focus.startsWith(id)).model;
assert.deepEqual(k('ArrowDown', real('f36428a'), null), {go: href('row-189', fx.stacks.find(s => s.model.focus.startsWith('652c341')).model.focus)});
const e = real('652c341');
assert.equal(e.rail[e.rail.length - 1].id, e.focus);
assert.match(k('ArrowDown', e, null).go, /^\/surface\?brane=shared&at=0ffbb9c/);
assert.match(k('ArrowLeft', e, null).go, /^\/surface\?brane=row-156&at=/);
// ---- Both views.
for (const m of [tm, sm, null]) {
  assert.deepEqual(k('/', m, null), {ui: 'address'}); assert.deepEqual(k('?', m, null), {ui: 'keys'});
  assert.deepEqual(k(',', m, null), {ui: 'appearance'}); assert.deepEqual(k('n', m, null), {ui: 'needs'}); assert.deepEqual(k('N', m, null), {ui: 'needs'});
}

// ---- shouldHandle.
const attrOf = v => ({getAttribute: n => n === 'contenteditable' ? v : null});
for (const tag of ['INPUT', 'TEXTAREA', 'SELECT', 'input']) assert.equal(api.shouldHandle({tagName: tag}), false, tag);
assert.equal(api.shouldHandle({tagName: 'DIV', isContentEditable: true}), false);
assert.equal(api.shouldHandle({tagName: 'DIV', contentEditable: 'true'}), false);
assert.equal(api.shouldHandle(Object.assign({tagName: 'DIV'}, attrOf(''))), false);
assert.equal(api.shouldHandle(Object.assign({tagName: 'DIV'}, attrOf('plaintext-only'))), false);
assert.equal(api.shouldHandle(Object.assign({tagName: 'DIV', contentEditable: 'inherit', isContentEditable: false}, attrOf('false'))), true);
assert.equal(api.shouldHandle(Object.assign({tagName: 'A'}, attrOf(null))), true);
assert.equal(api.shouldHandle({tagName: 'BODY'}), true); assert.equal(api.shouldHandle(null), true);
assert.equal(api.shouldHandle(timelineDoc.querySelector('input.address')), false);

// ---- loadLook, saveLook, resolveLook.
const DEF = {theme: 'paper', marks: 'colour', motion: 'system'};
const store = (v) => { const m = new Map(v === undefined ? [] : [['world.appearance', v]]); return {getItem: n => m.has(n) ? m.get(n) : null, setItem: (n, x) => m.set(n, String(x)), m}; };
const throwing = {getItem() { throw new Error('denied'); }, setItem() { throw new Error('denied'); }};
assert.deepEqual(api.loadLook(store()), DEF); assert.deepEqual(api.loadLook(null), DEF); assert.deepEqual(api.loadLook(throwing), DEF);
for (const garbage of ['{{', '42', 'null', '"night"', '[]', JSON.stringify({theme: 'neon', marks: 1, motion: null})]) assert.deepEqual(api.loadLook(store(garbage)), DEF, garbage);
assert.deepEqual(api.loadLook(store(JSON.stringify({theme: 'night'}))), {theme: 'night', marks: 'colour', motion: 'system'});
assert.doesNotThrow(() => api.saveLook(throwing, {theme: 'night', marks: 'shape', motion: 'full'}));
assert.doesNotThrow(() => api.saveLook(null, DEF));
const kept = store();
for (const look of [{theme: 'night', marks: 'shape', motion: 'reduced'}, {theme: 'system', marks: 'colour', motion: 'full'}, {theme: 'slate', marks: 'shape', motion: 'system'}]) {
  api.saveLook(kept, look); assert.deepEqual(api.loadLook(kept), look);
}
const R = (theme, marks, motion) => ({'data-theme': theme, 'data-marks': marks, 'data-motion': motion});
assert.deepEqual(api.resolveLook({theme: 'system', marks: 'colour', motion: 'system'}, true, false), R('slate', 'colour', 'full'));
assert.deepEqual(api.resolveLook({theme: 'system', marks: 'shape', motion: 'system'}, false, true), R('paper', 'shape', 'reduced'));
assert.deepEqual(api.resolveLook({theme: 'night', marks: 'colour', motion: 'full'}, false, true), R('night', 'colour', 'full'));
assert.deepEqual(api.resolveLook({theme: 'paper', marks: 'colour', motion: 'reduced'}, true, false), R('paper', 'colour', 'reduced'));
assert.deepEqual(api.resolveLook({theme: 'slate', marks: 'shape', motion: 'full'}, false, false), R('slate', 'shape', 'full'));
assert.deepEqual(api.resolveLook(null, false, false), R('paper', 'colour', 'full'));

// ---- The DOM wiring, in a fresh context over a parsed page.
function boot(html, stored) {
  const doc = parseHTML(html), gone = [], fetched = [], storage = store(stored);
  const context = {document: doc, window: {localStorage: storage}, location: {assign(h) { gone.push(h); }},
    matchMedia: q => ({matches: q.includes('dark'), addEventListener() {}}),
    fetch(url) { fetched.push(url); return Promise.resolve({ok: true, status: 200, json: () => Promise.resolve(index)}); }};
  vm.runInNewContext(fs.readFileSync('surface.js', 'utf8'), context);
  const press = (key, target = doc.activeElement) => {
    const ev = {key, target, defaultPrevented: false, ctrlKey: false, metaKey: false, altKey: false, preventDefault() { this.defaultPrevented = true; }};
    if (target && target !== doc.body && target.fire) target.fire('keydown', ev);
    doc.fire('keydown', ev);
    return ev;
  };
  return {doc, gone, fetched, storage, press, html: doc.documentElement};
}
const flush = async () => { for (let i = 0; i < 5; i++) await new Promise(r => setImmediate(r)); };
(async () => {
  let w = boot(fx.timeline.html, JSON.stringify({theme: 'system', marks: 'shape', motion: 'reduced'}));
  assert.equal(w.html.getAttribute('data-theme'), 'slate'); assert.equal(w.html.getAttribute('data-marks'), 'shape'); assert.equal(w.html.getAttribute('data-motion'), 'reduced');
  assert.equal(w.doc.querySelector('button.gear').hidden, false);
  const id = () => w.doc.activeElement.getAttribute('data-state');
  assert(w.press('End').defaultPrevented); assert.equal(id(), ranked[ranked.length - 1].id);
  for (const key of ['ArrowLeft', 'ArrowDown', 'ArrowDown', 'ArrowUp', 'ArrowRight', 'ArrowLeft']) {
    const before = id(), want = k(key, timeline, before);
    w.press(key); assert.equal(id(), want ? want.focus : before, key + ' from ' + before);
  }
  assert.notEqual(id(), ranked[ranked.length - 1].id, 'the keys moved');
  const focused = id(); w.press('Enter'); assert.deepEqual(w.gone, [timeline.lanes.flatMap(l => l.states).find(s => s.id === focused).href]);
  assert.equal(w.press('q').defaultPrevented, false);
  const bar = w.doc.querySelector('input.address'), note = w.doc.querySelector('.address-note');
  assert.equal(bar.value, 'world');
  w.press('/'); assert.equal(w.doc.activeElement, bar);
  for (const key of ['n', '?', ',', 'End']) w.press(key);
  assert.equal(w.doc.querySelector('div.panel'), null, 'no shortcut fires in the address bar'); assert.equal(w.doc.activeElement, bar);
  bar.value = 'world nosuch'; w.press('Enter'); await flush();
  assert.match(note.textContent, /no brane/); assert.equal(bar.getAttribute('aria-invalid'), 'true'); assert.equal(w.gone.length, 1, 'an error never navigates');
  bar.value = 'world onboarding@0'; w.press('Enter'); await flush(); assert.match(note.textContent, /states 1 to/); assert.equal(w.gone.length, 1);
  bar.value = 'world onboarding'; w.press('Enter'); await flush(); assert.equal(w.gone[1], href('row-189', lane189[lane189.length - 1]));
  bar.value = 'back'; w.press('Enter'); await flush(); assert.equal(w.gone[2], '/surface');
  assert.deepEqual(w.fetched, ['/surface/index.json'], 'the index is fetched once');
  w.press('Escape'); assert.equal(id(), focused, 'Esc in the bar returns to the page');
  w.press('?'); const keys = w.doc.querySelector('div.panel'); assert.equal(keys.getAttribute('data-panel'), 'keys'); assert.match(keys.text(), /lane above · below/);
  w.press('Escape'); assert.equal(w.doc.querySelector('div.panel'), null);
  w.press('N'); assert.match(w.doc.querySelector('div.panel').text(), /Nothing needs you/); assert.equal(w.gone.length, 3, 'needs goes nowhere'); w.press('Escape');
  w.press(','); const panel = w.doc.querySelector('div.panel'); assert.equal(panel.getAttribute('data-panel'), 'appearance'); assert.match(panel.text(), /Sunholo/);
  const labels = panel.querySelectorAll('label').map(l => l.text().trim());
  assert.deepEqual(labels, ['Paper', 'Slate', 'Night', 'Match system', 'Shape and colour', 'Shape', 'Full', 'Reduced', 'System setting']);
  const night = panel.querySelectorAll('input').find(i => i.getAttribute('name') === 'look-theme' && i.getAttribute('value') === 'night');
  night.checked = true; night.fire('change', {target: night});
  assert.equal(w.html.getAttribute('data-theme'), 'night'); assert.equal(JSON.parse(w.storage.m.get('world.appearance')).theme, 'night');
  assert.equal(w.press(',').defaultPrevented, false, 'no shortcut fires on a radio in the panel');
  w.press('Escape'); assert.equal(w.doc.querySelector('div.panel'), null, 'Esc closes the panel from its inputs');
  w.press(','); w.doc.querySelector('div.panel .close').focus(); w.press(','); assert.equal(w.doc.querySelector('div.panel'), null, ', toggles the panel');
  w.doc.querySelector('button.gear').fire('click', {}); assert.equal(w.doc.querySelector('div.panel').getAttribute('data-panel'), 'appearance');

  w = boot(fx.timeline.html, '{{');
  assert.equal(w.html.getAttribute('data-theme'), 'paper'); assert.equal(w.html.getAttribute('data-marks'), 'colour'); assert.equal(w.html.getAttribute('data-motion'), 'full');
  w.press('ArrowRight'); assert.equal(w.doc.activeElement.getAttribute('data-state'), ranked[ranked.length - 1].id, 'no focus: an arrow focuses the newest');

  const front = fx.stacks.find(s => s.model.focus.startsWith('f36428a'));
  w = boot(front.html);
  assert.equal(w.doc.querySelector('input.address').value, front.address); assert.equal(front.address, 'world row-189@f36428a');
  w.press('ArrowUp'); assert.deepEqual(w.gone, [k('ArrowUp', front.model, null).go]);
  w.press('?'); assert.match(w.doc.querySelector('div.panel').text(), /back in time · toward now/);
  w.press('Escape'); assert.equal(w.gone.length, 1, 'Esc closes a panel first');
  w.press('Escape'); assert.equal(w.gone[1], '/surface');
  console.log('SURFACE HARNESS PASS');
})().catch(e => { console.error(e); process.exitCode = 1; });
`
