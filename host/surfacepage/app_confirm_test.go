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

// TestNeedsYouListsConfirmations (row 243 V3 S6, design §2.5, §6 CP7) runs app.js's pure part under node
// on the daemon's own GET /surface/requests.json answers (testdata/requests-confirmations.golden.json,
// written by host/daemon TestRequestsJSONConfirmationsGolden): Needs you lists each pending finish
// confirmation after the decisions as a link to its card; a refused row shows its reason and no link, as
// does a row whose id is not a confirmation id; confirmationsUnavailable is said; every value is escaped.
// It skips only when no node exists; with WORLD_EXEC_NODE set, a harness that cannot run fails.
func TestNeedsYouListsConfirmations(t *testing.T) {
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	golden, err := os.ReadFile(filepath.Join("testdata", "requests-confirmations.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name, data := range map[string][]byte{"app.js": AppScript, "harness.js": []byte(needsYouConfirmationsHarness), "golden.json": golden} {
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
	if !strings.Contains(string(out), "NEEDS YOU CONFIRMATIONS PASS") {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

const needsYouConfirmationsHarness = `
const assert = require('node:assert/strict');
const app = require('./app.js');
const G = JSON.parse(require('node:fs').readFileSync('golden.json', 'utf8'));
const HEAD = (n) => '<div class="nsec" data-needs><div class="nh" style="color:var(--' + (n ? 'youShell' : 'sub') + ')">Needs you · ' + n + '</div>';
const refused = '<div class="nl sub">Confirm · vr-00000000000000c6 · refused: the stored request does not match its digest</div>';
const link = (id, ep, word) => '<div class="nl"><a id="confirm-' + id + '" href="/surface/confirm/' + id + '#activity">Confirm · episode finish ' + ep
  + ' · sunholo-data/ailang-world world/' + ep + '@0a0a0a0a0a0a · ' + word + ' · expires 08:09</a></div>';
const lines = refused + link('vr-00000000000000c4', 'w1', 'draft') + link('vr-00000000000000c2', 'w2', 'ready');
const request = G.v9.confirmations[1], heads = app.needsYouHTML({confirmations:[request]}, 0, null, true);
assert(heads.includes('/surface/confirm/' + request.id + '#heads/activity'));
const node = {kind:'open', id:'world:episode:w1', brane:'world:brane:shared', episode:'w1'};
const href = '/surface?at=world%3Aepisode%3Aw1&brane=world%3Abrane%3Ashared';
assert(app.finishMatchesNode(href, node));
assert(!app.finishMatchesNode(href, {...node, id:'other:episode:w1', brane:'other:brane:shared'}));
assert(!app.finishMatchesNode(href, {...node, kind:'shared'}));
assert(!app.finishMatchesNode(href.replace('shared','row-200'), node));
assert(!app.finishMatchesNode(href + '&at=world%3Aepisode%3Aw1', node));
assert(!app.finishMatchesNode('https://other.invalid' + href, node));
const step = {...node, id:'world:commit:abc', kind:'step'};
assert(app.finishMatchesNode(href, step, {'world:episode:w1':step.id}));
assert(!app.finishMatchesNode(href, step, {'world:episode:w1':'world:commit:def'}));
// The handler's answer: one line each, newest first, after the decisions; the refused one has no link.
const h = app.needsYouHTML(G.v9, G.v9.now);
assert.equal(h, HEAD(3) + lines + '</div>');
assert(!h.includes('/surface/confirm/vr-00000000000000c6') && !h.includes('confirm-vr-00000000000000c6') && h.split('<a ').length === 3, h);
// Below schema 9: the reason is said, not hidden — alone, and after the lines (escaped).
assert.equal(app.needsYouHTML(G.v8, G.v8.now), HEAD(0) + '<div class="nl sub">Nothing needs you</div><div class="nl sub">Confirm · unavailable: '
  + 'the credentials store predates verb requests: restart ` + "`ailang-worldd serve`" + ` with this build once</div></div>');
assert.equal(app.needsYouHTML({...G.v9, confirmationsUnavailable: '<i>u</i>', decisionsTruncated: true}, 0, 'granted'), HEAD(3) + lines
  + '<div class="nl sub">Confirm · unavailable: &lt;i&gt;u&lt;/i&gt;</div><div class="nl sub">Older approval requests are not listed; decide them at a terminal with ailang-worldd approve &lt;request-ref&gt;</div>'
  + '<div class="nl sub" data-notify>Notifications on</div></div>');
// None listed: exactly the empty Needs you.
assert.equal(app.needsYouHTML({...G.v9, confirmations: []}, 0), HEAD(0) + '<div class="nl sub">Nothing needs you</div></div>');
// After the mint requests and the decisions, counted with them.
{ const Q = 'ab'.repeat(32), dec = {request: 'sha256:' + Q, scope: 'worktree:1', requester: 'agent-1', deadline: 50520, reserved: '', refusal: ''};
  const m = app.needsYouHTML({...G.v9, pending: [{id: 'mr-1', episode: 'e', agent: 'a', grants: 1, expires: 60}], decisions: [dec]}, 0);
  assert(m.startsWith(HEAD(5) + '<div class="nl"><a id="needs-mr-1"'), m);
  assert(m.endsWith('<a id="decide-' + Q + '" href="/surface/decide/sha256:' + Q + '">Decide · worktree:1 · requested by agent-1 · deadline 14:02</a></div>' + lines + '</div>'), m); }
// Every value escaped where it lands (the sha cut to its first 12 characters, then escaped).
{ const real = G.v9.confirmations[1];
  assert.equal(app.needsYouHTML({confirmations: [{...real, episode: '<e>', repository: 'a&b', head: '"h"', sha: '<sha>0123456789'}]}, 0), HEAD(1)
    + '<div class="nl"><a id="confirm-vr-00000000000000c4" href="/surface/confirm/vr-00000000000000c4#activity">Confirm · episode finish &lt;e&gt; · a&amp;b &quot;h&quot;@&lt;sha&gt;0123456 · draft · expires 08:09</a></div></div>');
  assert.equal(app.needsYouHTML({confirmations: [{...G.v9.confirmations[0], id: '<x>', refused: '<i>r</i>'}]}, 0), HEAD(1)
    + '<div class="nl sub">Confirm · &lt;x&gt; · refused: &lt;i&gt;r&lt;/i&gt;</div></div>');
  // A summarised row whose id is not exactly vr- and 16 lower-case hex: its id, escaped, and no link.
  for (const id of ['vr-' + 'a'.repeat(15), 'vr-' + 'a'.repeat(17), 'vr-' + 'A'.repeat(16), 'xvr-' + 'a'.repeat(16), 'vr' + 'a'.repeat(16), 'vr_' + 'a'.repeat(16),
    'vr-' + 'g'.repeat(16), 'vr-' + 'a'.repeat(16) + '"><b>']) {
    const x = app.needsYouHTML({confirmations: [{...real, id}]}, 0);
    assert.equal(x, HEAD(1) + '<div class="nl sub">Confirm · ' + app.esc(id) + ' · refused: not a confirmation id</div></div>', id); } }
// A confirmation raises no notification (design §2.5 lists it; it does not notify).
assert.deepEqual(app.notifications(G.v9, new Set(), 'granted').notes, []);
console.log('NEEDS YOU CONFIRMATIONS PASS');
`
