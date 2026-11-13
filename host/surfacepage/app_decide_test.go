package surfacepage

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestNeedsYouListsDecisions (row 201 A4 XD3, design §2.8) runs app.js's pure part under node: Needs
// you lists GET /surface/requests.json's decisions after the mint requests, a decidable one as a link
// to its card, a reserved one naming its typed owner and a refused one its reason, neither with a
// link; every value escaped; a notification for each new decidable decision. It skips only when no
// node exists; with WORLD_EXEC_NODE set, a harness that cannot run fails.
func TestNeedsYouListsDecisions(t *testing.T) {
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	dir := t.TempDir()
	for name, data := range map[string][]byte{"app.js": AppScript, "harness.js": []byte(needsYouDecisionsHarness)} {
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
	if !strings.Contains(string(out), "NEEDS YOU DECISIONS PASS") {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

const needsYouDecisionsHarness = `
const assert = require('node:assert/strict');
const app = require('./app.js');
const Q = 'ab'.repeat(32), R = 'cd'.repeat(32), T = 'ef'.repeat(32), U = '01'.repeat(32);
const ok = {request: 'sha256:' + Q, packet: 'sha256:' + R, effect: 'Human.Approve', scope: 'worktree:1', cost: 1, requester: 'agent-1', deadline: 50520, reserved: '', refusal: ''};
const pub = {request: 'sha256:' + R, packet: 'sha256:' + T, effect: 'Human.Approve', scope: 'registry:pkg', cost: 1, requester: 'agent-2', deadline: 50520, reserved: 'publish', refusal: 'publish approvals are minted by world-publish approve'};
const fin = {request: 'sha256:' + U, packet: '', effect: 'Human.Approve', scope: 'github-pr:1', cost: 1, requester: 'agent-4', deadline: 0, reserved: 'finish', refusal: 'x'};
const two = {request: 'sha256:' + T, packet: '', effect: 'Human.Approve', scope: 'worktree:2', cost: 1, requester: 'agent-3', deadline: 0, reserved: '', refusal: 'carried by 2 decision packets'};
const line = '<div class="nl"><a id="decide-' + Q + '" href="/surface/decide/sha256:' + Q + '">Decide · worktree:1 · requested by agent-1 · deadline 14:02</a></div>';
// A decidable decision: one line, a link to its card; the count includes it; no mint request needed.
assert.equal(app.needsYouHTML({now: 1000, pending: [], decisions: [ok], decisionsTruncated: false}, 1000),
  '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div>' + line + '</div>');
// After the mint requests, in order; reserved names its owner with no link; refused shows its reason with no link.
{ const h = app.needsYouHTML({now: 0, pending: [{id: 'mr-1', episode: 'e', agent: 'a', grants: 1, expires: 60}], decisions: [ok, pub, fin, two]}, 0);
  assert(h.startsWith('<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 5</div><div class="nl"><a id="needs-mr-1"'), h);
  assert(h.includes(' min</a></div>' + line + '<div class="nl sub">publish approval — use world-publish approve at a terminal</div>'
    + '<div class="nl sub">finish approval — use ailang-worldd episode finish at a terminal</div>'
    + '<div class="nl sub">Decide · worktree:2 · refused: carried by 2 decision packets</div></div>'), h);
  for (const ref of [R, U, T]) assert(!h.includes('/surface/decide/sha256:' + ref) && !h.includes('decide-' + ref), h);
  assert.equal(h.split('<a ').length, 3, h); }
// An unknown reserved owner still renders no link.
assert.equal(app.needsYouHTML({pending: [], decisions: [{...ok, reserved: 'x<y'}]}, 0),
  '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div><div class="nl sub">x&lt;y approval — use its typed owner at a terminal</div></div>');
// Truncation is said once, after the lines.
assert.equal(app.needsYouHTML({pending: [], decisions: [ok], decisionsTruncated: true}, 0, 'granted'),
  '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div>' + line
  + '<div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div><div class="nl sub" data-notify>Notifications on</div></div>');
// A truncated read with nothing listed still says so.
assert.equal(app.needsYouHTML({pending: [], decisions: [], decisionsTruncated: true}, 0),
  '<div class="nsec" data-needs><div class="nh" style="color:var(--sub)">Needs you · 0</div><div class="nl sub">Nothing needs you</div><div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div></div>');
// Every value escaped where it lands: the ref in the id and href, the scope, the requester, the refusal.
{ const bad = '"><img src=x onerror=1>';
  assert.equal(app.needsYouHTML({pending: [], decisions: [{...ok, request: 'sha256:' + bad, scope: '<img src=x onerror=2>', requester: '<b>3</b>'}]}, 0),
    '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div><div class="nl"><a id="decide-&quot;&gt;&lt;img src=x onerror=1&gt;" href="/surface/decide/sha256:&quot;&gt;&lt;img src=x onerror=1&gt;">'
    + 'Decide · &lt;img src=x onerror=2&gt; · requested by &lt;b&gt;3&lt;/b&gt; · deadline 14:02</a></div></div>');
  assert.equal(app.needsYouHTML({pending: [], decisions: [{...two, scope: '<s>', refusal: '<i>4</i>'}]}, 0),
    '<div class="nsec" data-needs><div class="nh" style="color:var(--youShell)">Needs you · 1</div><div class="nl sub">Decide · &lt;s&gt; · refused: &lt;i&gt;4&lt;/i&gt;</div></div>'); }
// Without decisions (or with an empty list), Needs you is exactly as before.
assert.equal(app.needsYouHTML({pending: [], decisions: []}, 0), '<div class="nsec" data-needs><div class="nh" style="color:var(--sub)">Needs you · 0</div><div class="nl sub">Nothing needs you</div></div>');
{ const r = {now: 1000, pending: [{id: 'mr-0123456789abcdef', episode: 'w<1>', agent: 'claude', grants: 2, expires: 1601, recorded: 'webauthn:mark'}]};
  assert.equal(app.needsYouHTML({...r, decisions: [], decisionsTruncated: false}, 1000, 'default'), app.needsYouHTML(r, 1000, 'default')); }
// Notifications: one per new decidable decision while granted; none on the first read, none twice, none for a reserved or refused one.
{ const r1 = {pending: [], decisions: [ok]}, r2 = {pending: [{id: 'mr-a', episode: 'e1', grants: 1, expires: 50520}], decisions: [ok, pub, two, {...ok, request: 'sha256:' + U, scope: 's<2>', requester: 'agent-9', deadline: 50580}]};
  const first = app.notifications(r1, null, 'granted');
  assert.deepEqual(first.notes, []); assert.deepEqual([...first.seen], ['sha256:' + Q]);
  const second = app.notifications(r2, first.seen, 'granted');
  assert.deepEqual(second.notes, [{title: 'World needs you', tag: 'mr-a', body: 'Approve mint request for e1 (1 grant), expires 14:02', href: '/surface/approve/mr-a'},
    {title: 'World needs you', tag: 'sha256:' + U, body: 'Decide s<2> requested by agent-9, deadline 14:03', href: '/surface/decide/sha256:' + U}]);
  assert.deepEqual([...second.seen], ['sha256:' + Q, 'mr-a', 'sha256:' + R, 'sha256:' + T, 'sha256:' + U]);
  assert.deepEqual(app.notifications(r2, second.seen, 'granted').notes, []);
  assert.deepEqual(app.notifications(r2, first.seen, 'default').notes, []);
  assert.deepEqual(app.notifications({pending: [], decisions: [ok]}, new Set(), 'granted').notes, [{title: 'World needs you', tag: 'sha256:' + Q, body: 'Decide worktree:1 requested by agent-1, deadline 14:02', href: '/surface/decide/sha256:' + Q}]);
  assert.deepEqual(app.notifications({pending: [], decisions: [pub, two]}, new Set(), 'granted').notes, []);
  assert.deepEqual(app.notifications({pending: [], decisions: [{...ok, reserved: 'publish'}]}, new Set(), 'granted').notes, []); // reserved alone, refusal empty
  assert.deepEqual(app.notifications({pending: []}, new Set(), 'granted'), {notes: [], seen: new Set()}); }
console.log('NEEDS YOU DECISIONS PASS');
`
