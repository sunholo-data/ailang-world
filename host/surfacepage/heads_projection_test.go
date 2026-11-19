package surfacepage

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

func TestHeadsProjection(t *testing.T) {
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	env := followEnv(t)
	second := env.Open.Episodes[0]
	second.ID, second.Episode = "world:episode:w189y", "w189y"
	second.StartedAt = "2026-10-01T08:30:00Z"
	env.Open.Episodes[0].StartedAt = "2026-10-01T08:00:00Z"
	env.Open.Episodes = append(env.Open.Episodes, second)
	env.Open.Status, env.Open.Reason = surface.OpenIncomplete, "some sessions could not be read"
	if err := surface.ValidateFollow(env); err != nil {
		t.Fatal(err)
	}

	ix := env.Repos[0].Index
	documents := map[string]any{
		"follow": BuildFollowModel(env, followSources(), Now{At: fixedNow}),
		"single": BuildModel(ix, Now{At: fixedNow, Sessions: []Session{
			{Episode: "e1", Started: fixedNow}, {Episode: "w999", Started: fixedNow},
		}}),
		"unavailable": BuildModel(ix, Now{At: fixedNow, SessionsUnavailable: "credentials could not be read"}),
	}
	b, err := json.Marshal(documents)
	if err != nil {
		t.Fatal(err)
	}
	script, err := filepath.Abs(filepath.Join("app", "heads-data.js"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, "-e", headsProjectionHarness, script)
	cmd.Stdin = bytes.NewReader(b)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("heads projection: %v\n%s", err, out)
	}
}

const headsProjectionHarness = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const {project} = require(process.argv[1]);
const {follow, single, unavailable} = JSON.parse(fs.readFileSync(0));

function checkLanes(model) {
  const before = JSON.stringify(model);
  const projected = project(model);
  const lanes = model.lanes.filter(lane => lane.states.length);
  assert.deepEqual(projected.heads.map(head => head.key), lanes.map(lane => lane.key));
  projected.heads.forEach((head, i) => {
    const lane = lanes[i];
    assert.equal(head.name, lane.name);
    assert.equal(head.repo, lane.repo);
    assert.equal(head.glyph, lane.glyph);
    assert.equal(head.colour, lane.colour);
    assert.equal(head.latest.rank, Math.max(...lane.states.map(state => state.rank)));
    assert.equal(head.states.length, lane.states.length);
    for (const state of head.states) assert(lane.states.includes(state));
    for (let j = 1; j < head.states.length; j++) {
      const a = head.states[j - 1], b = head.states[j];
      assert(a.at > b.at || (a.at === b.at && a.rank > b.rank));
    }
  });
  assert.equal(JSON.stringify(model), before, 'project keeps the shared model intact');
  return projected;
}

const v2 = checkLanes(follow);
const world = v2.heads.find(head => head.key === 'world:brane:row-189');
const tools = v2.heads.find(head => head.key === 'tools:brane:row-189');
assert.equal(world.latest.id, 'world:episode:w189y');
assert.equal(world.states[0].kind, 'step', 'open started before the newest landed step');
assert.deepEqual(world.live.map(session => session.id), ['world:episode:w189z', 'world:episode:w189y']);
assert.equal(tools.live.length, 1);
assert.equal(tools.live[0].episode, 'e1');
assert.match(tools.live[0].id, /^tools:commit:/);
assert.equal(tools.live[0].label, 'landed; session live');
assert.equal(tools.source.state, 'stale');
assert.equal(v2.heads.find(head => head.latest.id === 'world:episode:w190q').live.length, 0);
assert.equal(v2.heads.find(head => head.latest.id === 'tools:episode:t1').live.length, 0);
assert.deepEqual(v2.unlinked, [{episode: 'w999', reason: 'start record unavailable'}]);
assert.match(v2.status, /some sessions could not be read/);
assert.match(v2.status, /tools · dev · stale/);

const v1 = checkLanes(single);
assert(v1.heads.some(head => head.key === 'shared'));
const linked = v1.heads.flatMap(head => head.live);
assert.equal(linked.length, 1);
assert.equal(linked[0].episode, 'e1');
assert.equal(v1.heads.find(head => head.key === linked[0].lane).states.find(state => state.id === linked[0].id).episode, 'e1');
assert.deepEqual(v1.unlinked, [{episode: 'w999', reason: 'No indexed state for this session'}]);
assert.equal(v1.status, '');
assert.equal(project(unavailable).status, 'credentials could not be read');
`
