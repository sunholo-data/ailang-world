package surfacepage

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// fixtureDist is the approved working copy, which the app is ported from.
var fixtureDist = filepath.Join("..", "..", "design_docs", "mockups", "surface-app-v3", "dist")

// TestSurfaceAppMatchesWorkingCopy: app.css is the fixture's but for its
// three font URLs, the logo is the fixture's, and the shell inlines the same
// logo markup and keeps the fixture's header and stage.
func TestSurfaceAppMatchesWorkingCopy(t *testing.T) {
	read := func(name string) []byte {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(fixtureDist, name))
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	css := string(read("app.css"))
	for _, f := range []string{"Montserrat-Variable.ttf", "IBMPlexMono-Regular.ttf", "IBMPlexMono-Medium.ttf"} {
		if _, _, ok := Font(f); !ok {
			t.Fatalf("font %s is not embedded", f)
		}
		css = strings.Replace(css, "url('fonts/"+f+"')", "url('/surface/fonts/"+f+"')", 1)
	}
	if css != string(AppStyle) {
		t.Error("app.css differs from the working copy beyond its font URLs")
	}
	logo := read("ailang-logo.svg")
	if !bytes.Equal(logo, AppLogo) {
		t.Error("ailang-logo.svg differs from the working copy")
	}
	if !strings.HasPrefix(string(AppLogo), "<svg") || bytes.IndexByte(AppLogo, 0) >= 0 {
		t.Error("the logo is not SVG text")
	}
	inner := string(logo)
	inner = inner[strings.Index(inner, ">")+1 : strings.LastIndex(inner, "</svg>")]
	shell := string(AppShell)
	if !strings.Contains(shell, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="32" height="32" role="img" aria-label="AILANG"><title>AILANG</title>`+inner+`</svg>`) {
		t.Error("the shell does not inline the logo's own markup")
	}
	index := string(read("index.html"))
	for _, part := range []string{"<header class=\"hdr\">", `<nav class="crumbs" id="crumbs" aria-label="Where you are"></nav>`, `<span class="prompt">world&gt;</span>`,
		`<span class="hmsg" id="hmsg"></span>`, `<section class="layer" id="L0" aria-label="The world timeline"><div class="canvas" id="tl"></div></section>`,
		`<section class="layer" id="L1" aria-label="One brane"><div class="canvas" id="br"></div></section>`, `<div class="tip" id="tip" hidden></div>`} {
		if !strings.Contains(index, part) || !strings.Contains(shell, part) {
			t.Errorf("the shell lost the working copy's %q", part)
		}
	}
	for _, banned := range []string{"SPIKE_DATA", "spike-data", "data.js"} {
		if strings.Contains(string(AppScript), banned) || strings.Contains(shell, banned) {
			t.Errorf("the app still names %s", banned)
		}
	}
	for _, name := range []string{"data.js", "spike-data.json"} {
		if _, err := os.Stat(filepath.Join("app", name)); err == nil {
			t.Errorf("app/%s is embedded", name)
		}
	}
}

func TestSurfaceAppShellAssembly(t *testing.T) {
	if bytes.Contains(AppShell, []byte(AppStylesheetLink)) || bytes.Count(AppShell, []byte("<style>")) != 1 {
		t.Fatal("the stylesheet is linked, not inlined once")
	}
	got := inlineStyle([]byte("a"+AppStylesheetLink+"b"), []byte("x{}"))
	if string(got) != "a<style>\nx{}</style>b" {
		t.Errorf("inlineStyle = %q", got)
	}
	for _, bad := range []struct{ index, css string }{{"no link", "x{}"}, {AppStylesheetLink + AppStylesheetLink, "x{}"}, {AppStylesheetLink, "</STYLE><script>"}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("inlineStyle(%q, %q) did not refuse", bad.index, bad.css)
				}
			}()
			inlineStyle([]byte(bad.index), []byte(bad.css))
		}()
	}
	for _, seam := range []string{"adaptModel", "layoutGraph", "frontHTML", "parseAppQuery", "addressTarget", "lookOf"} {
		if !strings.Contains(string(AppScript), seam) {
			t.Errorf("app.js lacks %s", seam)
		}
	}
	if !regexp.MustCompile(`if \(typeof module !== 'undefined'\) module\.exports = \{`).Match(AppScript) {
		t.Error("app.js's pure part is not exported for node")
	}
}

// TestSurfaceAppPure runs app.js's pure part under node on N1's committed
// goldens, with surface.js's parseAddress for the address bar. It skips only
// when no node exists; with WORLD_EXEC_NODE set, a harness that cannot run fails.
func TestSurfaceAppPure(t *testing.T) {
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	dir := t.TempDir()
	files := map[string][]byte{"app.js": AppScript, "surface.js": Script, "harness.js": []byte(appHarness)}
	for name, src := range map[string]string{
		"model.json":  filepath.Join("testdata", "model-rig.golden.json"),
		"stacks.json": filepath.Join("testdata", "stack-rig.golden.json"),
		"index.json":  filepath.Join("..", "surface", "testdata", "index-rig.json"),
	} {
		b, err := os.ReadFile(src)
		if err != nil {
			t.Fatal(err)
		}
		files[name] = b
	}
	for name, data := range files {
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
	if !strings.Contains(string(out), "APP HARNESS PASS") {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

const appHarness = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const app = require('./app.js');
const seams = require('./surface.js');
const model = JSON.parse(fs.readFileSync('model.json', 'utf8'));
const stacks = JSON.parse(fs.readFileSync('stacks.json', 'utf8'));
const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'barHTML', 'blendXS', 'cleanTitle', 'esc', 'frontHTML',
  'geometry', 'layoutGraph', 'levelXs', 'lookOf', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'shortTitle', 'stackBars', 'stackPath', 'summaryOf', 'verdict']);
const hollow = s => { for (const w of ['undefined', 'NaN', 'null', '[object']) assert(!s.includes(w), 'rendered ' + w + ': ' + s.slice(Math.max(0, s.indexOf(w) - 80), s.indexOf(w) + 40)); };

// ---- adaptModel: every state of model.json, on its lane, with the links as edges.
const T = app.adaptModel(model);
const all = model.lanes.flatMap(l => l.states.map(s => [l.key, s]));
assert.equal(Object.keys(T.N).length, all.length);
for (const [key, s] of all) { assert.equal(T.N[s.id].lane, key); assert.equal(T.N[s.id].title, s.title); assert.equal(T.N[s.id].summary, s.summary); }
for (const k of model.links) { assert(T.kids[k.from].includes(k.to)); assert(T.parents[k.to].includes(k.from)); }
const shared = model.lanes[0].states.filter(s => s.kind === 'shared').sort((a, b) => a.rank - b.rank).map(s => s.id);
assert.deepEqual(T.SH, shared); assert.equal(T.NS, shared.length); assert.equal(T.TT.length, T.NS + 1); assert.equal(T.TT[T.NS], Date.parse(model.now));
for (let i = 1; i < T.SH.length; i++) assert(T.parents[T.SH[i]].includes(T.SH[i - 1]), 'the shared line');
// Rows: the newest brane above, three below, the rest quiet; every brane lane exactly once.
const branes = model.lanes.slice(1).map(l => l.key);
assert.deepEqual([...T.above, ...T.below, ...T.quiet].sort(), branes.slice().sort());
const newest = model.lanes.slice(1).filter(l => l.states.some(s => s.kind === 'step'))
  .map(l => [l.key, Math.max(...l.states.filter(s => s.kind === 'step').map(s => s.rank))]).sort((a, b) => b[1] - a[1]).map(x => x[0]);
assert(newest.length >= 4 && T.quiet.length >= 1, 'the rig fills every row kind');
assert.deepEqual(T.above, newest.slice(0, 1)); assert.deepEqual(T.below, newest.slice(1, 4));
assert.deepEqual(T.ROWS, [...T.above, 'shared', ...T.below, 'other']);
assert.equal(T.LY.shared, 320); assert.equal(T.LY[T.above[0]], 182);
for (const [key, s] of all) {
  const r = T.rowOf(s.id);
  assert.equal(r, s.kind === 'shared' ? 'shared' : key === 'shared' ? 'dip' : T.quiet.includes(key) ? 'other' : key, s.id);
}
assert.equal(T.nameOf('shared'), 'Shared world'); assert.equal(T.nameOf('row-189'), 'World onboarding'); assert.equal(T.laneTitle('row-189'), 'World onboarding · row 189');
assert.equal(T.laneTitle('row-139'), 'row 139', 'a lane named by its row is not "row 139 · row 139"');
assert.equal(app.rowLabel('x', null), 'x');
assert.equal(T.lastState('row-189').id, model.lanes.find(l => l.key === 'row-189').states.slice(-1)[0].id);

// ---- levelXs and blendXS: five levels, steps get their slot, time never runs backwards.
assert.equal(T.LX.length, 5); assert.deepEqual(app.LEVELS.map(l => l.name), ['Close', 'Steps', 'Hours', 'Days', 'Weeks']);
T.LX.forEach((xs, k) => {
  assert.equal(xs.length, T.NS + 1); assert.equal(xs[0], 220);
  for (let i = 1; i < xs.length; i++) { assert(xs[i] >= xs[i - 1], 'level ' + k + ' runs backwards'); if (k <= 1) assert(xs[i] - xs[i - 1] >= app.LEVELS[k].G - 1e-9); }
});
assert.deepEqual(app.blendXS(T.LX, 1), T.LX[1]);
app.blendXS(T.LX, 4).forEach((x, i) => assert(Math.abs(x - T.LX[4][i]) < 1e-6, 'Weeks is the last level'));
app.blendXS(T.LX, 2.5).forEach((x, i) => assert(x >= Math.min(T.LX[2][i], T.LX[3][i]) && x <= Math.max(T.LX[2][i], T.LX[3][i])));
const gap = [Date.parse('2026-09-01T00:00:00Z'), Date.parse('2026-10-01T00:00:00Z'), Date.parse('2026-10-01T01:00:00Z')];
assert.equal(app.levelXs(gap)[4][1] - 220, 64, 'a quiet month is an axis break, not a month of axis');

// ---- layoutGraph: at every zoom, every state is drawn once, as a node or inside a bubble.
for (const u of [0, 0.5, 1, 2, 3, 4]) {
  const g = app.layoutGraph(T, app.blendXS(T.LX, u));
  const drawn = Object.values(g.DISP).flatMap(d => d.ids).sort();
  assert.deepEqual(drawn, Object.keys(T.N).sort(), 'zoom ' + u);
  for (const id of Object.keys(T.N)) { assert(g.DISP[g.REP[id]].ids.includes(id)); assert(Number.isFinite(g.POS[id][0]) && Number.isFinite(g.POS[id][1])); }
  for (const id of g.LABELED) assert(!id.startsWith('c:') && T.N[id].kind !== 'shared');
  for (const [id, d] of Object.entries(g.DISP)) assert.equal(id.startsWith('c:'), d.ids.length > 1);
}
const steps = app.layoutGraph(T, T.LX[1]);
for (const s of T.STEPS) { const m = T.mergeOf(s.id); if (m) assert(steps.POS[s.id][0] < steps.POS[m][0], s.id + ' sits before its merge'); }

// ---- The brane view, on every committed stack.
const w189d = stacks.find(s => s.front.id.startsWith('f36428a'));
for (const S of stacks) {
  const F = S.front, html = app.frontHTML(S, T);
  hollow(html);
  assert(html.includes(app.esc(app.cleanTitle(F.title) || F.label)), 'headline ' + F.id);
  assert(html.includes(app.esc(F.position)) && html.includes(app.esc(F.exact)));
  assert.equal((html.match(/data-read="/g) || []).length, (F.reads || []).length, 'one chip per read source');
  assert(html.includes('No view from the agent yet') && html.includes('row 151'), 'the agent band is an honest gap');
  if (!F.files) assert(html.includes('Files: not recorded'), 'absent files are a gap: ' + F.id);
  if (F.kind === 'step' && !F.episode) assert(html.includes('No episode record for this step'), F.id);
  if (F.kind === 'shared') { assert(html.includes('A merge has no checks of its own.')); assert(!html.includes('Grants')); assert.equal((html.match(/via">parent</g) || []).length, F.parents.length); }
  if (!F.checks.graders) assert(!/class="ctile"/.test(html), 'no tile invents a check: ' + F.id);
  if (!F.decisions.length) assert(html.includes('None recorded.'));
  const bars = app.stackBars(S.earlier, -1).concat(app.stackBars(S.later, 1));
  const want = [S.earlier, S.later].reduce((a, cols) => a + cols.slice(0, 3).reduce((b, c) => b + Math.min(3, c.cards.length), 0), 0);
  assert.equal(bars.length, want); assert.equal(new Set(bars.map(b => b.key)).size, bars.length);
  bars.forEach(b => { hollow(app.barHTML(b)); assert(b.width > 0 && b.op >= 0.35 && b.op <= 1); });
  const rail = app.railHTML(S);
  hollow(rail);
  assert.equal((rail.match(/aria-current="true"/g) || []).length, 1);
  assert.equal((rail.match(/class="rrow"/g) || []).length, Math.min(15, S.rail.length));
}
const dh = app.frontHTML(w189d, T);
assert(dh.includes('Passed all 4 checks after 2 rounds · 13.5 agent-minutes'), 'w189d verdict');
assert(dh.includes('↑ Came from · 8 reads over 4 files'));
assert.equal((dh.match(/data-read="/g) || []).length, 4);
assert(dh.includes('Task</span><div class="task"><div class="tl1">row-189-M5-plugin</div>'));
assert(dh.includes('Grants: not recorded yet.'), 'grants are not in the records');
assert(dh.includes('data-go="0ee1945ffba795812f06dce95998b646adfac343"') && dh.includes('via">built on<'));
assert(dh.includes('data-go="90035557b2012895551867c0ebe2f36b853aac09"'), 'led to its merge');
assert.equal((dh.match(/class="act"/g) || []).length, w189d.front.decisions.length);
assert.equal((dh.match(/class="file"/g) || []).length, w189d.front.files.length);
// Absent stays absent: a stack without the optional fields renders gaps, not values.
const bare = JSON.parse(JSON.stringify(w189d));
for (const k of ['reads', 'readCalls', 'readFiles', 'files', 'filesTruncated', 'task', 'model', 'summary', 'mergedInto', 'builtOn']) delete bare.front[k];
bare.front.checks = {}; bare.front.decisions = [];
const bh = app.frontHTML(bare, T);
hollow(bh);
for (const gap of ['Reads: not recorded for this step.', 'Files: not recorded for this step.', 'No episode record for this step.', 'None recorded.', 'Grants: not recorded yet.']) assert(bh.includes(gap), gap);
assert(!bh.includes('class="ctile"') && !bh.includes('Agent</span>') && !bh.includes('Came from'));
assert.equal(app.verdict('8/8', 1, undefined, 3.5), 'Passed all 8 checks in one round · 3.5 agent-minutes');
assert.equal(app.verdict('7/8', 3, false, 18.4), '7 of 8 checks passed after 3 rounds · 18.4 agent-minutes');
assert.equal(app.verdict('5/5', 1, true, null), 'Passed all 5 checks on the first round');
assert.equal(app.verdict('', 1, true, 2), '');
// A long rail shows a window around its front.
const long = JSON.parse(JSON.stringify(w189d));
long.rail = Array.from({length: 40}, (_, i) => Object.assign({}, w189d.rail[0], {id: 'r' + i, label: 'r' + i}));
long.front.id = 'r3';
const lr = app.railHTML(long);
assert.equal((lr.match(/class="rrow"/g) || []).length, 15); assert(lr.includes('aria-current="true"')); assert(lr.includes('15 of 40 states'));

// ---- Text helpers (v3's human layer).
assert.equal(app.cleanTitle('Exec profile bound at mint, and ` + "`episode start`" + ` (through World, w189b)'), 'Exec profile bound at mint, and episode start');
assert.equal(app.cleanTitle('Session credentials in their own store — mint while the daemon serves (made through World, w189a)'), 'Session credentials in their own store — mint while the daemon serves');
assert.equal(app.prettyModel('claude-opus-5-5'), 'Claude Opus 5.5'); assert.equal(app.prettyModel('gpt-5-codex'), 'GPT 5 Codex'); assert.equal(app.prettyModel(''), '');
assert(app.shortTitle('a '.repeat(40)).endsWith('…')); assert.equal(app.esc('<a "b">&'), '&lt;a &quot;b&quot;&gt;&amp;');

// ---- Places: /surface's query grammar, the URL the view follows, and the address bar.
assert.deepEqual(app.parseAppQuery(''), {}); assert.deepEqual(app.parseAppQuery('?'), {});
assert.deepEqual(app.parseAppQuery('?brane=row-189&at=abc'), {brane: 'row-189', at: 'abc'});
assert.deepEqual(app.parseAppQuery('at=abc&brane=row%2D189'), {brane: 'row-189', at: 'abc'});
for (const bad of ['?x=1', '?brane=a', '?at=b', '?brane=a&at=b&c=d', '?brane=a&brane=b&at=c', '?brane=a&at=b&at=c']) assert.equal(app.parseAppQuery(bad), null, bad);
assert.equal(app.appPath(), '/surface'); assert.equal(app.appPath('row 1', 'a&b'), '/surface?brane=row%201&at=a%26b');
assert.equal(app.stackPath('row-189', 'f3'), '/surface/stack.json?brane=row-189&at=f3');
for (const s of stacks) assert.deepEqual(app.parseAppQuery(app.appPath(s.front.brane, s.front.id).slice('/surface'.length)), {brane: s.front.brane, at: s.front.id});
const P = seams.parseAddress, lane189 = seams.laneStates(index, 'row-189');
const at = (text, ix = index) => app.addressTarget(text, ix, P);
assert.deepEqual(at('row-189@1'), {open: {brane: 'row-189', at: lane189[0].id}}, 'typed after world>');
assert.deepEqual(at('world row-189@1'), {open: {brane: 'row-189', at: lane189[0].id}});
assert.deepEqual(at('onboarding'), {open: {brane: 'row-189', at: lane189[lane189.length - 1].id}});
assert.deepEqual(at('shared@0ee1945'), {open: {brane: 'shared', at: '0ee1945ffba795812f06dce95998b646adfac343'}});
assert.deepEqual(at('back', null), {home: true}); assert.deepEqual(at('world', null), {home: true}); assert.deepEqual(at('  '), {home: true});
assert.match(at('nosuch').error, /no brane/); assert.match(at('onboarding@0').error, /states 1 to/); assert.match(at('a b c').error, /unknown address/);
assert.match(at('row-189@1', null).error, /no brane/, 'without the index a name cannot resolve');
assert.match(app.addressTarget('row-189@1', index, undefined).error, /not loaded/);
// Every committed stack's own address opens it.
for (const s of stacks) {
  const lane = seams.laneStates(index, s.front.brane), n = lane.findIndex(x => x.id === s.front.id) + 1;
  assert.deepEqual(at(s.front.brane + '@' + n), {open: {brane: s.front.brane, at: s.front.id}});
  assert.deepEqual(at(s.front.brane + '@' + s.front.id.slice(0, 7)), {open: {brane: s.front.brane, at: s.front.id}});
}

// ---- Appearance: v3's two choices from M4a's stored look.
const L = (theme, motion) => ({theme, motion});
assert.deepEqual(app.lookOf(seams.loadLook(null), false, false), L('paper', 'full'));
assert.deepEqual(app.lookOf(seams.loadLook(null), true, true), L('paper', 'reduced'), 'M4a default motion follows the system');
assert.deepEqual(app.lookOf({theme: 'slate', motion: 'full'}, false, true), L('slate', 'full'));
assert.deepEqual(app.lookOf({theme: 'night'}, false, false), L('slate', 'full'));
assert.deepEqual(app.lookOf({theme: 'system', motion: 'reduced'}, true, false), L('slate', 'reduced'));
assert.deepEqual(app.lookOf(null, true, true), L('paper', 'full')); assert.deepEqual(app.lookOf('junk'), L('paper', 'full'));
const kept = new Map(), store = {getItem: k => kept.has(k) ? kept.get(k) : null, setItem: (k, v) => kept.set(k, v)};
seams.saveLook(store, Object.assign(seams.loadLook(store), L('slate', 'reduced')));
assert.deepEqual(app.lookOf(seams.loadLook(store), false, false), L('slate', 'reduced'), 'kept under world.appearance');
assert(kept.has('world.appearance'));
console.log('APP HARNESS PASS');
`
