package surfacepage

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestHeadDraftsProjectExactTabLocalBranches(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for the head draft projection test")
	}
	script, err := filepath.Abs(filepath.Join("app", "head-drafts.js"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, node, "-e", headDraftHarness, script).CombinedOutput()
	if err != nil {
		t.Fatalf("head drafts: %v\n%s", err, out)
	}
}

const headDraftHarness = `
const assert = require('node:assert/strict');
const {project, restore} = require(process.argv[1]);
const at = '2026-11-19T12:00:00.000Z';
const draft = {id:'draft:active-01', parent:'commit-1', at, title:'Check   the rig'};
const stored = restore([
  draft,
  {...draft, title:'Duplicate'},
  {...draft, id:'draft:active-02', parent:'draft:active-01'},
  {...draft, id:'draft:active-03', at:'not a date'},
  {...draft, id:'draft:active-04', title:'    '},
  {...draft, id:'draft:active-05', parent:'missing'},
  {...draft, id:'wrong-prefix'},
  null
]);
assert.deepEqual([...stored.keys()], ['draft:active-01', 'draft:active-05']);

const parent = Object.freeze({id:'commit-1', kind:'step', at, rank:19});
const newer = Object.freeze({id:'commit-2', kind:'step', at, rank:22});
const source = Object.freeze({key:'repo:row-2', repo:'repo', colour:'#497a82', glyph:'star', states:Object.freeze([parent,newer])});
const link = Object.freeze({from:'commit-1', to:'commit-2', colour:source.colour});
const raw = Object.freeze({lanes:Object.freeze([source]), links:Object.freeze([link])});
const {model, parents} = project(raw, stored);
assert.equal(parents.get(parent.id).state, parent);
assert.equal(model.lanes[0], source);
assert.equal(model.links[0], link);
assert.equal(raw.lanes.length, 1);
assert.equal(model.lanes.length, 2);
const lane = model.lanes[1], state = lane.states[0], edge = model.links[1];
assert.deepEqual([lane.name,lane.repo,lane.colour,lane.glyph], ['Draft · this tab','repo','#497a82','star']);
assert.deepEqual([state.kind,state.draft,state.parent,state.title,state.at,state.rank],
  ['draft',true,'commit-1','Check the rig',at,23]);
assert.deepEqual(edge,{from:'commit-1',to:draft.id,colour:'#497a82'});
const refreshed = project(model, stored);
assert.deepEqual(refreshed.model, model, 'a second projection does not duplicate overlays');
assert(!refreshed.parents.has(draft.id), 'drafts are not eligible parents');
assert.equal(project({lanes:[],links:[]}, stored).model.links.length, 0, 'missing parent does not create a dangling edge');
assert.deepEqual(project(model, new Map()).model, raw, 'discard removes the local lane and its edge');
`
