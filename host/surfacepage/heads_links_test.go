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
)

func TestHeadsLinks(t *testing.T) {
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	env := followEnv(t)
	model, err := json.Marshal(BuildFollowModel(env, followSources(), Now{At: fixedNow}))
	if err != nil {
		t.Fatal(err)
	}
	script, err := filepath.Abs(filepath.Join("app", "heads-links.js"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, "-e", headsLinksHarness, script)
	cmd.Stdin = bytes.NewReader(model)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("heads links: %v\n%s", err, out)
	}
}

const headsLinksHarness = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const {project, detailsMarkup, geometry} = require(process.argv[1]);
const model = JSON.parse(fs.readFileSync(0));

const before = JSON.stringify(model);
const projected = project(model);
const endpointLane = new Map(model.lanes.flatMap(lane => lane.states.map(state => [state.id, lane.key])));
const expected = model.links.filter(({from, to}) => endpointLane.get(from) !== endpointLane.get(to));
const actual = projected.pairs.flatMap(pair => pair.edges.map(({from, to}) => [from.state.id, to.state.id]));
assert.deepEqual(actual.sort(), expected.map(({from, to}) => [from, to]).sort(), 'retain all directed timeline links across lanes');
assert(projected.pairs.some(pair => pair.edges.length > 1), 'collapse multiple state links between the same pair of branes');
assert(projected.pairs.every(pair => pair.edges.every(({from, to}) => from.lane.repo === to.lane.repo)), 'keep colliding repo IDs separate');
assert.equal(JSON.stringify(model), before, 'the timeline model is unchanged');

const selected = 'world:brane:row-189';
const shown = new Map([[selected, {}]]);
const details = detailsMarkup(projected, shown, selected);
assert.match(details, /off page/);
assert.match(details, /Earlier state/);
assert.match(details, /<strong>Base/);
assert.match(details, /<strong>Merge/);
assert.match(details, /row-189&(?:amp;)?at=world%3Aepisode%3Aw189z#heads/);
assert.doesNotMatch(details, /tools%3A/);
assert(!actual.some(edge => edge.includes('world:episode:w190q')), 'an open state with an unavailable base gets no inferred link');

const points = geometry({x: 10, y: 80, width: 210, height: 30, exact: true}, {x: 500, y: 100, width: 210, height: 30});
assert.deepEqual(points.source, {x: 222, y: 95});
assert.deepEqual(points.target, {x: 495, y: 115});
`
