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
	// The one departure from v3 is appended, marked: the brand as the way home (Mark, 26 Oct).
	style := string(AppStyle)
	cut := strings.Index(style, "\n/* ---------- After v3: the brand is the way home")
	if cut < 0 || style[:cut] != css {
		t.Error("app.css differs from the working copy beyond its font URLs and the brand link")
	} else {
		for _, rule := range []string{"a.brand { color: inherit; text-decoration: none;", "a.brand:hover { background: rgba(255,255,255,0.06); }",
			"a.brand:focus-visible { outline: 2px solid var(--hCtl);"} {
			if !strings.Contains(style[cut:], rule) {
				t.Errorf("the brand link lacks %q", rule)
			}
		}
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
	// The working copy's brand is a span; the app's is the link home, with the same children.
	brand := `<a class="brand" id="home" href="/surface" aria-label="World: back to the timeline">` + "\n      " + `<span class="logo-tile"><svg `
	if !strings.Contains(index, `<span class="brand">`+"\n      "+`<span class="logo-tile">`) {
		t.Error("the working copy's brand is no longer a span")
	}
	if strings.Count(shell, `class="brand"`) != 1 || !strings.Contains(shell, brand) || strings.Contains(shell, `<span class="brand">`) {
		t.Errorf("the shell's brand is not %q", brand)
	}
	if !strings.Contains(shell, `</svg></span>`+"\n      "+`<span class="brand-name">World</span>`+"\n    </a>") {
		t.Error("the brand link does not hold the logo and \"World\"")
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
		"golden":      filepath.Join("..", "agui", "testdata", "stream_fixture.golden"),
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
assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'applyDelta', 'barHTML', 'blendXS', 'cleanTitle', 'crumbsHTML', 'decidedHTML', 'esc', 'fetchCoherent', 'filePath', 'fileEvidenceHTML', 'filesTabHTML', 'frontHTML',
  'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'notifications', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict',
  'VERB_ROUTE_BUDGET_MS', 'askVerb', 'verbAct', 'verbOutcome', 'verbPanel', 'goWith', 'addressGo'].sort());

assert.equal(app.serviceLineHTML(undefined), '');
// Row 201 A3: Needs you is the pending requests, each a link to its card; a refused row shows why, with no link.
assert.equal(app.needsYouHTML(null, 0), '<div class="nsec" data-needs><div class="nh" style="color:var(--sub)">Needs you · 0</div><div class="nl sub">Nothing needs you</div></div>');
{ const h = app.needsYouHTML({now: 1000, pending: [{id: 'mr-0123456789abcdef', episode: 'w<1>', agent: 'claude', grants: 2, expires: 1601, recorded: 'webauthn:mark'}, {id: 'mr-00000000000000ff', episode: 'x', refused: 'digest <bad>', recorded: null}]}, 1000);
  assert(h.includes('Needs you · 2') && h.includes('<a id="needs-mr-0123456789abcdef" href="/surface/approve/mr-0123456789abcdef"><b>w&lt;1&gt;</b> · claude · 2 grants · expires in 11 min</a> · approved by webauthn:mark, waiting for the minter</div>'), h);
  assert(h.includes('mr-00000000000000ff · refused: digest &lt;bad&gt;') && !h.includes('/surface/approve/mr-00000000000000ff'), h); }
assert(app.needsYouHTML({pending: [], unavailable: 'the mint requests could not be read'}, 0).includes('>the mint requests could not be read<'));
// Every value is escaped where it lands: hostile text in each field and in the link's id and href (sprint-4 judge A02/A03/A05).
assert.equal(app.needsYouHTML({pending: [{id: 'mr-"><img src=x onerror=1>', episode: 'e', agent: '<img src=x onerror=2>', grants: 1, expires: 60, recorded: '"><script>3</script>'},
  {id: 'mr-<b>"4</b>', episode: 'f', refused: '<i>5</i>'}]}, 0), '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 2</div>'
  + '<div class="nl"><a id="needs-mr-&quot;&gt;&lt;img src=x onerror=1&gt;" href="/surface/approve/mr-&quot;&gt;&lt;img src=x onerror=1&gt;"><b>e</b> · &lt;img src=x onerror=2&gt; · 1 grant · expires in 1 min</a>'
  + ' · approved by &quot;&gt;&lt;script&gt;3&lt;/script&gt;, waiting for the minter</div><div class="nl sub">mr-&lt;b&gt;&quot;4&lt;/b&gt; · refused: &lt;i&gt;5&lt;/i&gt;</div></div>');
assert(app.needsYouHTML({pending: [], unavailable: '<i>"down"</i>'}, 0).includes('<div class="nl sub">&lt;i&gt;&quot;down&quot;&lt;/i&gt;</div>'));
// The notification opt-in follows the browser's permission; nothing without the API (the calls above).
assert(app.needsYouHTML(null, 0, 'default').endsWith('<div class="nl sub">Nothing needs you</div><button type="button" class="nl run" id="notify-me">Notify me</button></div>'));
assert(app.needsYouHTML({pending: [{id: 'mr-1', episode: 'e', agent: 'a', grants: 1, expires: 60}]}, 0, 'granted').endsWith(' min</a></div><div class="nl sub" data-notify>Notifications on</div></div>'));
assert(app.needsYouHTML(null, 0, 'denied').endsWith('<div class="nl sub" data-notify>Notifications are blocked for this page</div></div>'));
// A notification per request this tab has not seen, while granted: none on the page's first read, none twice, none for a refused or recorded one.
{ const r1 = {pending: [{id: 'mr-a', episode: 'e1', grants: 2, expires: 50520}]}, r2 = {pending: r1.pending.concat([{id: 'mr-b', episode: 'e<2>', grants: 1, expires: 50580},
    {id: 'mr-c', episode: 'e3', refused: 'bad'}, {id: 'mr-d', episode: 'e4', grants: 3, expires: 60, recorded: 'webauthn:m'}])};
  assert.deepEqual(app.notifications(null, null, 'granted'), {notes: [], seen: null});
  const first = app.notifications(r1, null, 'granted');
  assert.deepEqual(first.notes, []); assert.deepEqual([...first.seen], ['mr-a']);
  assert.equal(app.notifications({}, first.seen, 'granted').seen, first.seen);
  const second = app.notifications(r2, first.seen, 'granted');
  assert.deepEqual(second.notes, [{title: 'World needs you', tag: 'mr-b', body: 'Approve mint request for e<2> (1 grant), expires 14:03', href: '/surface/approve/mr-b'}]);
  assert.deepEqual([...second.seen], ['mr-a', 'mr-b', 'mr-c', 'mr-d']);
  assert.deepEqual(app.notifications(r2, second.seen, 'granted').notes, []);
  const asked = app.notifications(r2, first.seen, 'default');
  assert.deepEqual(asked.notes, []); assert.deepEqual([...asked.seen], ['mr-a', 'mr-b', 'mr-c', 'mr-d']);
  assert.equal(app.notifications({pending: [{id: 'mr-e', episode: 'e5', grants: 2, expires: 50520}]}, new Set(), 'granted').notes[0].body, 'Approve mint request for e5 (2 grants), expires 14:02'); }
// The decision list: the last three decided requests, newest first, and a toggle for all of them; every value escaped.
assert.equal(app.decidedHTML(null, false), '');
assert.equal(app.decidedHTML({decided: []}, true), '');
{ const d = [{episode: 'w<4>', approvedBy: 'webauthn:"m"', at: 50520}, {episode: 'w3', approvedBy: 'ssh:<k>', at: 50460}, {episode: 'w2', approvedBy: 'tty', at: 3600}, {episode: 'w1', approvedBy: 'tty', at: 0}];
  const few = '<div class="nsec" data-decided><div class="nh">Decided</div><div class="nl sub">You · Minted · w&lt;4&gt; · webauthn:&quot;m&quot; · 14:02</div>'
    + '<div class="nl sub">You · Minted · w3 · ssh:&lt;k&gt; · 14:01</div><div class="nl sub">You · Minted · w2 · tty · 01:00</div>';
  assert.equal(app.decidedHTML({decided: d}, false), few + '<button type="button" class="nl run" id="decided-all" aria-expanded="false">Show all 4</button></div>');
  assert.equal(app.decidedHTML({decided: d}, true), few + '<div class="nl sub">You · Minted · w1 · tty · 00:00</div><button type="button" class="nl run" id="decided-all" aria-expanded="true">Show fewer</button></div>');
  assert.equal(app.decidedHTML({decided: d.slice(0, 3)}, false), few + '</div>'); }
assert.equal(app.serviceLineHTML({line:'<build & file>',warn:false}), '<div class="nl sub" data-service>&lt;build &amp; file&gt;</div>');
assert.equal(app.serviceLineHTML({line:'restart',warn:true}), '<div class="nl sub" data-service style="color:var(--you)">restart</div>');

const bytes = fs.readFileSync('golden');
const expected = bytes.toString().split('\n\n').slice(0,-1).flatMap(f => f.split('\n').filter(l => l.startsWith('data: ')).map(l => JSON.parse(l.slice(6))));
assert(expected.length > 0);
function chunked(size) {
 const decoder = new TextDecoder(); let tail = '', events = [];
 for (let i = 0; i < bytes.length;) {
   const n = size(); const p = app.splitFrames(tail, decoder.decode(bytes.subarray(i, i+n), {stream:true}));
   i += n; tail = p.tail; events.push(...p.events);
 }
 const p = app.splitFrames(tail, decoder.decode()); events.push(...p.events);
 assert.equal(p.tail, ''); return events;
}
let seed = 188;
for (const size of [() => bytes.length, () => 1, () => {seed = (seed * 1664525 + 1013904223) >>> 0; return 1 + seed % 47;}]) assert.deepEqual(chunked(size), expected);
let cursor = {lastIndex:-1};
for (const e of expected) {
 if (e.type === 'STATE_SNAPSHOT') cursor = {...e.snapshot};
 if (e.type === 'STATE_DELTA') app.applyDelta(cursor, e.delta);
}
assert.equal(cursor.lastIndex, 6);
assert.equal(cursor.logHead, 'sha256:' + '0'.repeat(63) + '7');
assert.equal(typeof cursor.logHead, 'string');
assert.throws(() => app.applyDelta({}, [{op:'replace',path:'/other',value:1}]));
assert.throws(() => app.applyDelta({}, [{op:'add',path:'/lastIndex',value:1}]));
let back = 0;
for (const want of [1,2,4,8,16,30,30]) {back = app.nextBackoff(null,back); assert.equal(back,want);}
assert.equal(app.nextBackoff('4',0),4);
assert.deepEqual(app.liveRequest(1234,7), {threadId:'surface',runId:'surface-7',messages:[],state:{schema:'world/agui-state/v1',lastIndex:1234}});
for (const [s,want] of [
 [{mode:'live',cursor:1234,checked:'14:05',count:0},'Live · up to entry 1234 · checked 14:05'],
 [{mode:'live',count:3},'Live · 3 new entries · refreshing'],
 [{mode:'live',count:1},'Live · 1 new entry · refreshing'],
 [{mode:'paused',reason:'stream limit',wait:4},'Paused: stream limit · retrying in 4 s'],
 [{mode:'hidden'},'Paused while the tab is hidden'],
 [{mode:'unavailable',reason:'the log could not be read'},'Live updates unavailable: the log could not be read'],
 [{mode:'stopped',reason:'unknown cursor'},'Live stopped: unknown cursor']
]) assert.equal(app.liveStatus(s),want);

const hollow = s => { for (const w of ['undefined', 'NaN', 'null', '[object']) assert(!s.includes(w), 'rendered ' + w + ': ' + s.slice(Math.max(0, s.indexOf(w) - 80), s.indexOf(w) + 40)); };

// ---- adaptModel: every state of model.json, on its lane, with the links as edges.
const T = app.adaptModel(model);
const all = model.lanes.flatMap(l => l.states.map(s => [l.key, s]));
assert.equal(Object.keys(T.N).length, all.length);
for (const [key, s] of all) { assert.equal(T.N[s.id].lane, key); assert.equal(T.N[s.id].title, s.title); assert.equal(T.N[s.id].summary, s.summary); }
for (const k of model.links) { assert(T.kids[k.from].includes(k.to)); assert(T.parents[k.to].includes(k.from)); }
const shared = model.lanes[0].states.filter(s => s.kind === 'shared').sort((a, b) => a.rank - b.rank).map(s => s.id);
assert.deepEqual(T.SH, shared);
assert.equal(T.TT.length, T.NS + 1);
assert.deepEqual(T.TT, [...new Set([...all.map(([, s]) => Date.parse(s.at)), Date.parse(model.now)])].sort((a, b) => a - b));
assert.equal(T.TT[T.presentIndex], Date.parse(model.now));
for (const [, state] of all) assert.equal(T.TT[T.timeIndex.get(state.id)], Date.parse(state.at), state.id);
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
  const xs = app.blendXS(T.LX, u), g = app.layoutGraph(T, xs);
  const drawn = Object.values(g.DISP).flatMap(d => d.ids).sort();
  assert.deepEqual(drawn, Object.keys(T.N).sort(), 'zoom ' + u);
  for (const id of Object.keys(T.N)) { assert(g.DISP[g.REP[id]].ids.includes(id)); assert.equal(g.POS[id][0], xs[T.timeIndex.get(id)]); assert(Number.isFinite(g.POS[id][1])); }
  for (const id of g.LABELED) assert(!id.startsWith('c:') && T.N[id].kind !== 'shared');
  for (const [id, d] of Object.entries(g.DISP)) assert.equal(id.startsWith('c:'), d.ids.length > 1);
}
const steps = app.layoutGraph(T, T.LX[1]);
for (const s of T.STEPS) {
  const m = T.mergeOf(s.id);
  if (m) assert.equal(Math.sign(steps.POS[s.id][0] - steps.POS[m][0]), Math.sign(Date.parse(s.at) - Date.parse(T.N[m].at)), s.id + ' uses its own time');
}

// A branch need not have a merge to appear on the axis. Missing timestamps
// get one slot before history, simultaneous states share x and drafts use now.
const branch = app.adaptModel({now:'2026-10-01T00:00:00Z', lanes:[
  {key:'trial', states:[
    {id:'a', kind:'step', rank:0, at:'2026-09-01T00:00:00Z'},
    {id:'tie', kind:'open', rank:1, at:'2026-08-31T20:00:00-04:00'},
    {id:'future', kind:'open', rank:2, at:'2026-10-02T00:00:00Z'},
    {id:'unknown', kind:'step', rank:3},
    {id:'invalid', kind:'open', rank:4, at:'unrecorded'},
    {id:'draft', kind:'draft', rank:5, at:'2026-09-23T00:00:00Z', draft:true},
  ]},
], links:[{from:'a', to:'tie'}]});
assert.equal(branch.SH.length, 0);
assert.deepEqual(branch.ROWS, ['trial']);
assert.equal(branch.TT.length, 4);
assert(branch.undated);
assert.equal(branch.timeIndex.get('unknown'), 0);
assert.equal(branch.timeIndex.get('invalid'), 0);
assert.equal(branch.timeIndex.get('draft'), branch.presentIndex);
assert.equal(branch.timeIndex.get('future'), branch.NS);
assert(branch.presentIndex < branch.NS);
assert.deepEqual(branch.parents.tie, ['a']);
for (const xs of branch.LX) {
  const {POS, DISP, REP} = app.layoutGraph(branch, xs, 'tie');
  assert.equal(POS.a[0], POS.tie[0]);
  assert.equal(POS.tie[0], DISP[REP.tie].x);
  assert(POS.a[0] < POS.draft[0] && POS.draft[0] < POS.future[0]);
  assert(Object.values(POS).every(at => at.every(Number.isFinite)));
}
const withoutClock = app.adaptModel({now:'unknown', lanes:[]});
assert(!withoutClock.nowKnown && !withoutClock.dated);
assert(withoutClock.TT.every(Number.isFinite));

// ---- The brane view, on every committed stack.
const w189d = stacks.find(s => s.front.id.startsWith('f36428a'));
for (const S of stacks) {
  const F = S.front, html = app.frontHTML(S, T);
  hollow(html);
  assert(html.includes(app.esc(app.cleanTitle(F.title) || F.label)), 'headline ' + F.id);
  assert(html.includes(app.esc(F.position)) && html.includes(app.esc(F.exact)));
  assert.equal((html.match(/data-read="/g) || []).length, (F.reads || []).length, 'one chip per read source');
  assert(!html.includes('From the agent') && !html.includes('row 151'), 'no output advertises an absent view');
  if (!F.files && !F.landedFiles && F.landedFileCount == null) assert(app.filesTabHTML(F).includes('No file records are available'), 'absent files are a gap: ' + F.id);
  if (F.kind === 'step' && !F.episode) assert(html.includes('No episode record for this step'), F.id);
  if (F.kind === 'shared') { assert(!html.includes('Grants')); assert.equal((html.match(/via">parent</g) || []).length, F.parents.length); }
  if (!F.decisions.length && !Object.values(F.checks).some(value => value != null && value !== '')) {
    assert(!html.includes('What happened'), 'no empty history column: ' + F.id);
  }
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
assert(dh.includes('4/4 graders · 2 rounds · 13.5 agent-minutes'), 'w189d measured checks');
assert(dh.indexOf('>Inputs<') < dh.indexOf('>What happened<') && dh.indexOf('>What happened<') < dh.indexOf('>Outputs<'));
const decisionPositions = w189d.front.decisions.map(d => dh.indexOf(app.esc(d.text)));
assert(decisionPositions.every((p, i) => p >= 0 && (i === 0 || p > decisionPositions[i - 1])), 'recorded events retain chronological order');
assert(dh.includes('↑ Came from · 8 reads over 4 files'));
assert.equal((dh.match(/data-read="/g) || []).length, 4);
assert(dh.includes('Task title</h3>') && dh.includes('<p>row-189-M5-plugin</p>'));
assert(dh.includes('Grants: not recorded yet.'), 'grants are not in the records');
assert(dh.includes('data-go="0ee1945ffba795812f06dce95998b646adfac343"') && dh.includes('via">built on<'));
assert(dh.includes('data-go="90035557b2012895551867c0ebe2f36b853aac09"'), 'led to its merge');
assert.equal((dh.match(/class="act"/g) || []).length, w189d.front.decisions.length);
const legacyFiles = app.filesTabHTML(w189d.front);
assert(legacyFiles.includes('Recorded step reads') && legacyFiles.includes('Recorded step changes'));
assert.equal((legacyFiles.match(/data-file-row=/g) || []).length,
  w189d.front.files.length + w189d.front.reads.reduce((n, group) => n + group.files.length, 0));
// Absent stays absent: a stack without the optional fields renders gaps, not values.
const bare = JSON.parse(JSON.stringify(w189d));
for (const k of ['reads', 'readCalls', 'readFiles', 'files', 'filesTruncated', 'landedFiles', 'landedFileCount', 'task', 'model', 'summary', 'mergedInto', 'builtOn']) delete bare.front[k];
bare.front.checks = {}; bare.front.decisions = [];
const bh = app.frontHTML(bare, T);
hollow(bh);
for (const gap of ['Reads: not recorded for this step.', 'Grants: not recorded yet.']) assert(bh.includes(gap), gap);
assert(!bh.includes('What happened') && !bh.includes('Agent</span>') && !bh.includes('Task title</span>') && !bh.includes('Came from'));

// The merge's Git paths provide a fallback, with their own provenance and no invented grading.
const gitFiles = [
  {path: 'asset/scene & <room>.tscn', status: 'A'},
  {path: 'src/a/same.go', status: 'M'},
  {path: 'src/b/same.go', status: 'D'},
  ...Array.from({length: 5}, (_, n) => ({path: 'docs/deep/file-' + n + '.md', status: 'M'})),
];
function detailWith(fields) {
  const html = app.frontHTML({...bare, front: {...bare.front, ...fields}}, T);
  hollow(html);
  return html;
}
// The Files tab keeps historical step changes separate from committed paths and actual episode evidence.
const filesWith = (fields, activity) => app.filesTabHTML({...bare.front, ...fields}, activity);
const merged = filesWith({landedFiles: gitFiles, landedFileCount: 8});
assert(merged.includes('In the merge <span>8</span>'));
assert(merged.includes('>Added<') && merged.includes('>Modified<') && merged.includes('>Deleted<'));
assert(merged.includes('data-file-path="asset/scene &amp; &lt;room&gt;.tscn"') && merged.includes('src/a/same.go') && merged.includes('src/b/same.go'));
const carriageReturn = filesWith({landedFiles: [{path: 'docs/line\rbreak.md', status: 'M'}], landedFileCount: 1});
assert(carriageReturn.includes('data-file-path="docs/line&#13;break.md"'));
assert(app.filePath('world:brane:shared', 'world:commit:abc', 'docs/line\rbreak.md').endsWith('path=docs%2Fline%0Dbreak.md'));
assert(!merged.includes('file-more') && merged.includes('docs/deep/file-4.md'));
assert.equal((merged.match(/data-file-row=/g) || []).length, gitFiles.length);
const sharedFiles = filesWith({kind: 'shared', episode: '', parents: [], landedFiles: gitFiles, landedFileCount: 10});
assert(sharedFiles.includes('In the commit <span>8</span>') && sharedFiles.includes('2 committed file paths unavailable.'));
assert(!sharedFiles.includes('In the merge'));
const emptyStep = filesWith({files: [], filesTruncated: 0, landedFiles: gitFiles, landedFileCount: 8});
assert(emptyStep.includes('No recorded step file changes.') && emptyStep.includes('In the merge') && emptyStep.includes('src/a/same.go'));
const partialStep = filesWith({files: [gitFiles[0]], filesTruncated: 1, landedFiles: gitFiles, landedFileCount: 8});
assert(partialStep.includes('Recorded step changes <span>1</span>') && partialStep.includes('1 recorded step file path unavailable.'));
assert(partialStep.includes('In the merge') && partialStep.includes('src/a/same.go'));
assert(filesWith({landedFileCount: 0, landedFiles: []}).includes('No committed file changes.'));
const evidencePath = 'docs/a & b.md';
const episodeFiles = filesWith({trajectoryEpisode:'e1', landedFiles: [{path:evidencePath,status:'M'}], landedFileCount:1}, {
  status:'partial', reads:[{path:evidencePath,count:2}], writes:[{path:evidencePath,count:1,op:'workspace-edit'}],
  produced:[{path:'reports/log.txt',count:1}], omitted:{writes:1}, deltaGaps:1
});
for (const source of ['episode-read','episode-write','committed']) assert(episodeFiles.includes('data-file-path="docs/a &amp; b.md" data-file-source="' + source + '"'));
assert(episodeFiles.includes('2 reads') && episodeFiles.includes('Last: edit') && episodeFiles.includes('Changes observed after commands'));
assert(!episodeFiles.includes('data-file-path="reports/log.txt"') && episodeFiles.includes('Some path records were omitted.'));
assert(filesWith({trajectoryEpisode:'e1'}, {status:'complete', reads:[],writes:[],produced:[]}).includes('No file activity was recorded'));
assert(filesWith({trajectoryEpisode:'e1'}, {status:'unavailable'}).includes('Episode file activity is unavailable'));
assert(app.filePath('world:brane:shared','world:commit:abc',evidencePath,'episode-write').endsWith('path=docs%2Fa%20%26%20b.md&source=episode-write'));
const fullTask = '# Inspect the run\n\nUse <room> & preserve the next line.\n\n- First file\n- **Second file**\n\nLast line.';
const task = detailWith({task: 'Short recorded task', taskBody: fullTask, taskRedacted: true});
assert(task.includes('<h4>Inspect the run</h4>') && task.includes('<p>Use &lt;room&gt; &amp; preserve the next line.</p>'));
assert(task.includes('<ul><li>First file</li><li><strong>Second file</strong></li></ul>') && task.includes('<p>Last line.</p>'));
assert(!task.includes('Read full task') && !task.includes('Short recorded task') && task.includes('Sensitive task text was redacted.'));
const sameTask = detailWith({task: 'Already short', taskBody: 'Already short'});
assert(sameTask.includes('Task request</h3>') && sameTask.includes('<p>Already short</p>'));
const liveTask = detailWith({kind: 'open', pr: 0, open: {label: 'Session live'}, taskBody: fullTask});
assert(liveTask.includes('<h4>Inspect the run</h4>') && liveTask.includes('<p>Last line.</p>') && liveTask.includes('No finish commit recorded'));
const reviewed = detailWith({commits: [{id: 'world:commit:123456789', kinds: ['episode-finish'],subject:'A <change>', body:'Detailed summary.'}],
  run: {episode:'ep-1',launches:1,turns:2,costUsd:0,decisions:[{action:'Approved pull request creation',at:'2026-11-21T12:00:00Z'}]}});
assert(reviewed.includes('Episode finish commit <code>12345678</code>') && reviewed.includes('A &lt;change&gt;') && reviewed.includes('Detailed summary.'));
assert(reviewed.includes('Human decisions</h3><ol>') && reviewed.includes('Approved pull request creation'));
assert(reviewed.includes('Agent trajectory</h2>') && reviewed.includes('Completed turns</dt><dd>2') && reviewed.includes('Provider cost</dt><dd>$0.00'));

const metered = detailWith({run:{launches:1, tokens:{input:258254,cacheRead:198784,cacheWrite:0,output:1883,reasoning:24}}});
assert(metered.includes('<dt>Input tokens</dt><dd>258,254') && metered.includes('198,784 read from cache · 0 written to cache'));
assert(metered.includes('<dt>Output tokens</dt><dd>1,883') && metered.includes('24 reasoning'));
assert(!metered.includes('no verified usage totals'));
const partialUsage = detailWith({run:{launches:1,tokens:{cacheRead:0,reasoning:0}}});
assert(partialUsage.includes('0 read from cache') && partialUsage.includes('<dt>Reasoning tokens</dt><dd>0</dd>') && !partialUsage.includes('<dt>Input tokens</dt>'));
const permission = (state,budget,used,remaining) => ({state, grants:[{effect:'Workspace.Write',scope:'worktree',budget,used,remaining}]});
const grants = detailWith({run:{launches:0,permissions:{sessions:[permission('active',80,78,2),permission('active',8,78,0),permission('active',86,78,8)]}}});
assert(grants.includes('Permissions &amp; budgets') && grants.includes('78 used · limit 86 · 8 left') && grants.includes('Grant history · 3'));
assert.equal((grants.match(/<b>Write files<\/b>/g)||[]).length,1);
const expired = detailWith({run:{launches:0,permissions:{sessions:[permission('expired',80,17,63)]}}});
assert(expired.includes('17 used · limit 80 · Expired') && !expired.includes('63 left'));
const unused = detailWith({run:{launches:0,permissions:{sessions:[permission('active',8,0,8)]}}});
assert(unused.includes('0 used · limit 8 · 8 left'));
const privateGrants = detailWith({run:{launches:0,permissions:{sessions:[{state:'active',grants:[
  {effect:'Workspace.Read',scope:'[private scope]',scopeKey:1,budget:5,used:2,remaining:3},
  {effect:'Workspace.Read',scope:'[private scope]',scopeKey:2,budget:6,used:4,remaining:2}
]}]}}});
assert(privateGrants.includes('2 used · limit 5 · 3 left') && privateGrants.includes('4 used · limit 6 · 2 left'));
const linkedTrace = detailWith({kind:'shared',trajectoryEpisode:'ep-1'});
assert(linkedTrace.includes('Linked episode trajectory</h2>') && !linkedTrace.includes('Run details</h2>'));

const sections = app.frontHTML({...bare,front:{...bare.front,run:{episode:'scope',launches:0}}},T,'files');
assert(sections.includes('aria-label="Node sections"') && sections.includes('id="node-tab-files" aria-controls="node-panel-files" aria-selected="true"'));
assert(sections.includes('id="node-panel-overview" aria-labelledby="node-tab-overview" hidden') && sections.includes('data-node-tab="run"') && sections.includes('data-node-tab="trajectory"'));
const gitOnly = detailWith({kind:'shared',githubChecks:{checks:[]}});
assert(gitOnly.includes('>Checks</button>') && !gitOnly.includes('data-node-tab="trajectory"'));

const rawMerge = 'An original merge title with a suffix — made through World (w999) & <proof>';
const mergeTitle = detailWith({kind: 'shared', episode: '', pr: 631, fullTitle: rawMerge, title: 'An original merge title'});
assert(mergeTitle.includes('<summary>Recorded merge title</summary>') && mergeTitle.includes(app.esc(rawMerge)));
assert(!mergeTitle.includes('Task title') && !mergeTitle.includes('Read full task'));

// File content is pinned to a node's exact qualified locator and treated as plain text.
assert.equal(app.filePath('world:brane:shared', 'world:commit:abc', 'docs/a+b #.md'),
  '/surface/file.json?brane=world%3Abrane%3Ashared&at=world%3Acommit%3Aabc&path=docs%2Fa%2Bb%20%23.md');
const proof = {kind: 'text', merge: 'world:commit:abc', status: 'D', text: '<script>bad</script>\nBefore', diff: '-<script>bad</script>', truncated: true};
const patchView = app.fileEvidenceHTML(proof, 'diff');
hollow(patchView);
assert(patchView.includes('data-file-mode="diff" aria-pressed="true"') && patchView.includes('Only part of the diff is shown.'));
assert(patchView.includes('-&lt;script&gt;bad&lt;/script&gt;') && !patchView.includes('Before'));
const fileView = app.fileEvidenceHTML(proof, 'file');
hollow(fileView);
assert(fileView.includes('data-file-mode="file" aria-pressed="true"') && fileView.includes('Showing the first-parent version.'));
assert(fileView.includes('&lt;script&gt;bad&lt;/script&gt;\nBefore') && !fileView.includes('Only part of the diff'));
assert(app.fileEvidenceHTML({kind: 'text', status: 'A', text: ''}).includes('Empty file.'));
for (const kind of ['binary', 'symlink', 'submodule', 'too-large', 'unexpected']) {
  const html = app.fileEvidenceHTML({kind, text: 'unexpected text'});
  assert(html.includes('class="muted"') && !html.includes('<pre') && !html.includes('unexpected text'));
}


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
// Crumbs: the brand is the way home, so the timeline has none and a brane has only its address.
assert.equal(app.crumbsHTML(null), ''); assert.equal(app.crumbsHTML(undefined), '');
assert.equal(app.crumbsHTML('row-189@6'), '<button type="button" class="crumb" aria-current="true">row-189@6</button>');
assert(!/world|›|data-crumb/.test(app.crumbsHTML('row-189@6')), 'no root crumb in a brane');
assert.equal(app.crumbsHTML('a<b'), '<button type="button" class="crumb" aria-current="true">a&lt;b</button>');
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

// The operator could not copy a request id or hash from the app (2026-11-16): v3's body rule is
// user-select: none. The marked after-v3 override restores text selection after that rule, and
// keeps buttons unselectable.
func TestSurfaceAppTextIsSelectable(t *testing.T) {
	style := string(AppStyle)
	v3 := strings.Index(style, "user-select: none; -webkit-user-select: none; cursor: default; }")
	marker := strings.Index(style, "/* ---------- After v3: text can be selected and copied")
	override := strings.Index(style, "body { user-select: text; -webkit-user-select: text; }")
	if v3 < 0 || marker < v3 || override < marker {
		t.Fatalf("want v3's body rule, then the marked override restoring text selection (v3 %d, marker %d, override %d)", v3, marker, override)
	}
	if !strings.Contains(style[marker:], "button { user-select: none; -webkit-user-select: none; }") {
		t.Fatal("buttons must stay unselectable")
	}
}
