package surfacepage

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestSurfaceAppFollow runs app.js's and surface.js's pure parts under node
// on the row 220 M4 follow documents Go builds here: the grouped layout (AC9),
// repo-qualified addresses and their ambiguity, the model/stack generation
// handshake (AC13), and the Running/Sources sections (AC12). It skips only
// when no node exists; with WORLD_EXEC_NODE set, a harness that cannot run
// fails.
func TestSurfaceAppFollow(t *testing.T) {
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	env := followEnv(t)
	model := BuildFollowModel(env, followSources(), Now{At: fixedNow})
	step := stateWithPrefix(t, env.Repos[0].Index, "f36428a")
	commit, err := BuildFollowStack(env, "tools:brane:row-189", "tools:commit:"+step.ID, fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	open, err := BuildFollowStack(env, "world:brane:shared", "world:episode:w189z", fixedNow)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	files := map[string]any{"model.json": model, "envelope.json": env, "commit.json": commit, "open.json": open}
	for name, v := range files {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), b, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for name, data := range map[string][]byte{"app.js": AppScript, "surface.js": Script, "harness.js": []byte(appFollowHarness)} {
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
	if !strings.Contains(string(out), "APP FOLLOW HARNESS PASS") {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

const appFollowHarness = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const app = require('./app.js');
const seams = require('./surface.js');
const read = f => JSON.parse(fs.readFileSync(f, 'utf8'));
const model = read('model.json'), env = read('envelope.json'), commit = read('commit.json'), open = read('open.json');
const hollow = s => { for (const w of ['undefined', 'NaN', 'null', '[object']) assert(!s.includes(w), 'rendered ' + w + ': ' + s.slice(Math.max(0, s.indexOf(w) - 80), s.indexOf(w) + 40)); };
const repo = id => String(id).split(':')[0];

// ---- Grouped layout (AC9): one surface, each repo its own shared world, rows in config order, no edge across repos.
assert(app.isSharedKey('shared') && app.isSharedKey('world:brane:shared') && !app.isSharedKey('world:brane:row-189') && !app.isSharedKey('shared:x'));
const T = app.adaptModel(model);
assert.deepEqual(T.GROUPS.map(g => g.id), ['world', 'tools']); assert(T.multi);
assert.deepEqual(T.ROWS.filter(app.isSharedKey), ['world:brane:shared', 'tools:brane:shared']);
const rowsOf = r => T.ROWS.filter(k => repo(k.replace(/^(dip|other)\|/, '')) === r).map(k => T.ROWS.indexOf(k));
assert(Math.max(...rowsOf('world')) < Math.min(...rowsOf('tools')), 'world rows all above tools rows');
for (const r of T.ROWS) assert(Number.isFinite(T.LY[r]), 'row ' + r);
for (let i = 1; i < T.ROWS.length; i++) assert(T.LY[T.ROWS[i]] > T.LY[T.ROWS[i - 1]], 'rows run down the page');
for (const [c, ps] of Object.entries(T.parents)) for (const p of ps) assert.equal(repo(p), repo(c), 'edge ' + p + ' → ' + c + ' joins repos');
const sharedLine = id => T.N[id].kind === 'shared';
for (const g of T.GROUPS) {
  const line = T.SH.filter(id => T.N[id].lane === g.shared);
  for (let i = 1; i < line.length; i++) assert(T.parents[line[i]].includes(line[i - 1]), g.id + ' shared line');
}
const all = model.lanes.flatMap(l => l.states.map(s => s.id));
assert.equal(new Set(all).size, all.length, 'every id distinct across repos');
for (const u of [0, 1, 2, 4]) {
  const g = app.layoutGraph(T, app.blendXS(T.LX, u));
  assert.deepEqual(Object.values(g.DISP).flatMap(d => d.ids).sort(), Object.keys(T.N).sort(), 'zoom ' + u);
  for (const id of Object.keys(T.N)) assert(Number.isFinite(g.POS[id][0]) && Number.isFinite(g.POS[id][1]), id);
}
const ep = T.N['world:episode:w189z'];
assert.equal(ep.kind, 'open'); assert.equal(T.rowOf(ep.id), 'world:brane:row-189', 'the running episode sits on its branch row');
assert(T.parents[ep.id].includes(model.links.find(k => k.to === ep.id).from), 'anchored at its base');
assert.equal(T.aliasOf['tools:episode:e1'], model.aliases[0].target, 'an alias resolves to its landed step');
// The legacy single-world path is unchanged by grouping (covered by TestSurfaceAppPure); one repo is not "multi".
const one = JSON.parse(JSON.stringify(model)); one.repos = one.repos.slice(0, 1); one.lanes = one.lanes.filter(l => l.repo === 'world');
const T1 = app.adaptModel(one);
assert(!T1.multi); assert.equal(T1.LY['world:brane:shared'], 320); assert(T1.ROWS.includes('world:brane:shared'));

// ---- The open card and the landed alias (minimal; N2d owns the richer card).
const oh = app.frontHTML(open, T);
hollow(oh);
for (const want of ['w189z', 'session live', 'Polish the plugin docs', 'Nothing has landed yet.', 'No checks until it lands.', 'via">built on<', 'trial/row189-m7']) assert(oh.includes(want), want);
assert(!oh.includes('A merge has no checks'), 'an open episode is not a merge');
const ch = app.frontHTML(commit, T);
hollow(ch); assert(ch.includes('data-go="tools:commit:'), 'commit stack chips stay in their repo');
assert(!ch.includes('data-go="world:'), 'a tools stack links no world state');
hollow(app.railHTML(commit)); hollow(app.railHTML(open));

// ---- Running and Sources (AC12): status per repo; an incomplete or unavailable read never says "Nothing running".
const run = app.openRunningHTML(model.nowEdge.open);
hollow(run);
assert(run.includes('Running · 3') && run.includes('data-go="world:episode:w189z"') && run.includes('start record unavailable') && run.includes('1 PR open'));
assert(app.openRunningHTML({status: 'complete', live: []}).includes('Nothing running'));
for (const o of [{status: 'unavailable', note: 'running status unavailable: the session read failed', live: []}, {status: 'incomplete', note: 'session read bound reached', live: []}]) {
  const h = app.openRunningHTML(o);
  assert(!h.includes('Nothing running'), o.status + ' claimed zero'); assert(h.includes(o.note));
}
const src = app.sourcesHTML(model.repos);
hollow(src);
assert(src.includes('data-repo="world" data-state="fresh">world · dev · local branch checked 08:59<'));
assert(src.includes('data-repo="tools" data-state="stale" style="color:var(--you)">tools · dev · stale · last successful local check 08:40'));
assert(src.includes('next: ailang-worldd service status'));
assert.equal(app.sourcesHTML(undefined), ''); assert.equal(app.sourcesHTML([]), '');

// ---- Addresses: short names only when unique; ambiguity lists repo-qualified choices.
const P = (t, base) => seams.parseAddress(t, env, base);
assert.match(P('world row-189').error, /ambiguous name "row-189": choose world:row-189, tools:row-189/);
assert.match(P('world onboarding').error, /choose world:row-189, tools:row-189/);
assert.match(P('world shared').error, /choose world:shared, tools:shared/);
const lane = seams.laneStates(env.repos[1].index, 'row-189');
assert.deepEqual(P('world tools:row-189@1'), {href: '/surface?brane=' + encodeURIComponent('tools:brane:row-189') + '&at=' + encodeURIComponent('tools:commit:' + lane[0].id)});
assert.deepEqual(P('world tools:row-189@' + lane[0].id.slice(0, 7), '/surface/static'), {href: '/surface/static?brane=' + encodeURIComponent('tools:brane:row-189') + '&at=' + encodeURIComponent('tools:commit:' + lane[0].id)});
const sharedTools = seams.laneStates(env.repos[1].index, 'shared');
assert.deepEqual(P('world tools'), {href: '/surface?brane=' + encodeURIComponent('tools:brane:shared') + '&at=' + encodeURIComponent('tools:commit:' + sharedTools[sharedTools.length - 1].id)});
assert.deepEqual(P('world w189z'), {href: '/surface?brane=' + encodeURIComponent('world:brane:shared') + '&at=' + encodeURIComponent('world:episode:w189z')}, 'a unique episode id');
assert.deepEqual(P('world e1'), {href: '/surface?brane=' + encodeURIComponent('world:brane:shared').replace('world', 'tools') + '&at=' + encodeURIComponent('tools:episode:e1')}, 'an alias');
assert.match(P('world tools:row-189@0').error, /tools:row-189 has states 1 to/);
assert.match(P('world nosuch').error, /no brane or episode/);
assert.deepEqual(app.addressTarget('tools:row-189@1', env, seams.parseAddress), {open: {brane: 'tools:brane:row-189', at: 'tools:commit:' + lane[0].id}});

// ---- The generation handshake (AC13): same → no re-render; mismatch → retry once; twice → keep the view.
const reply = (status, body) => ({ok: status === 200, status, body});
const fake = script => { const calls = []; return {calls, get: async url => { calls.push(url); const r = script.shift(); assert(r, 'unexpected fetch ' + url); return r; }}; };
const M = gen => reply(200, {schema: 'world/surface-model/v2', generation: gen, lanes: []});
const S = gen => reply(200, {generation: gen, front: {id: 'x'}});
const sel = {brane: 'world:brane:row-189', at: 'world:commit:x'};
(async () => {
  let f = fake([M('g1')]);
  assert.equal((await app.fetchCoherent(f.get, sel, 'g1')).kind, 'same'); assert.equal(f.calls.length, 1, 'same generation fetches no stack');
  f = fake([M('g2')]);
  let r = await app.fetchCoherent(f.get, null, 'g1'); assert.equal(r.kind, 'apply'); assert.equal(r.stack, null);
  f = fake([M('g2'), S('g2')]);
  r = await app.fetchCoherent(f.get, sel, 'g1'); assert.equal(r.kind, 'apply'); assert.equal(r.stack.generation, 'g2');
  assert.equal(f.calls[1], app.stackPath(sel.brane, sel.at));
  f = fake([M('g2'), S('g3'), M('g3'), S('g3')]);
  r = await app.fetchCoherent(f.get, sel, 'g1'); assert.equal(r.kind, 'apply'); assert.equal(r.model.generation, 'g3'); assert.equal(f.calls.length, 4, 'retried once');
  f = fake([M('g2'), S('g3'), M('g3'), S('g4')]);
  assert.equal((await app.fetchCoherent(f.get, sel, 'g1')).kind, 'keep', 'mismatch twice keeps the previous view');
  f = fake([reply(503, {schema: 'world/surface-model/v2', refused: true, reason: 'no repo has an accepted graph yet', repos: []})]);
  r = await app.fetchCoherent(f.get, sel, 'g1'); assert.equal(r.kind, 'fail'); assert.equal(r.status, 503);
  f = fake([M('g2'), reply(404, {reason: 'episode gone'}), M('g2'), reply(404, {reason: 'episode gone'})]);
  r = await app.fetchCoherent(f.get, sel, 'g1'); assert.equal(r.kind, 'apply'); assert.equal(r.missing, 'episode gone');
  // Legacy (no generation): every refresh applies, a stack failure is a failure, as in N2c.
  f = fake([reply(200, {lanes: []}), reply(200, {front: {id: 'x'}})]);
  r = await app.fetchCoherent(f.get, sel, null); assert.equal(r.kind, 'apply'); assert(r.stack);
  f = fake([reply(200, {lanes: []}), reply(404, {reason: 'gone'})]);
  assert.equal((await app.fetchCoherent(f.get, sel, null)).kind, 'fail');
  console.log('APP FOLLOW HARNESS PASS');
})().catch(e => { console.error(e); process.exit(1); });
`
