package surfacepage

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestHeadCameraUsesTimelineGeometry(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for the head camera geometry test")
	}
	app, err := filepath.Abs(filepath.Join("app", "app.js"))
	if err != nil {
		t.Fatal(err)
	}
	camera, err := filepath.Abs(filepath.Join("app", "head-camera.js"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, node, "-e", headCameraHarness, app, camera).CombinedOutput()
	if err != nil {
		t.Fatalf("camera geometry: %v\n%s", err, out)
	}
}

const headCameraHarness = `
const assert = require('node:assert/strict');
const {adaptModel} = require(process.argv[1]);
const {layout, turn} = require(process.argv[2]);
const now = Date.parse('2026-11-19T12:00:00Z');
const state = (id, rank, days, extra = {}) => ({
  id, rank, kind: 'step', at: days == null ? '' : new Date(now - days * 86400000).toISOString(), ...extra
});
const lane = (key, states) => ({key, states});
const model = {
  now: new Date(now).toISOString(),
  lanes: [
    lane('shared', Array.from({length: 90}, (_, i) => state('s' + i, i, 90 - i, {kind:'shared'}))),
    lane('recent', [state('r1', 105, .5), state('r2', 106, 11)]),
    lane('live', [state('a1', 92, 40, {sessionLive:true})]),
    lane('attention', [state('q1', 93, 9)]),
    lane('selected', [state('e1', 94, 28), state('e2', 95, 15)]),
    lane('unknown', [state('u1', 96, null)]),
    lane('middle', [state('m1', 97, 5)]),
    ...Array.from({length: 12}, (_, i) => lane('old' + i, [state('o' + i, 110 + i, 13 + i)]))
  ],
  links: [{from:'s89',to:'r1'}, {from:'r1',to:'r2'}, {from:'r1',to:'r2'}]
};
const timeline = adaptModel(model);
const options = {selected:'selected', attention:new Set(['attention'])};
for (const width of [320, 390, 1300]) {
  const scene = layout(timeline, {width, height:640, ...options});
  assert.equal(scene.points.size, Object.keys(timeline.N).length, 'every timeline state has a point');
  assert.equal(scene.heads.length, model.lanes.length, 'shared and all thread lanes have heads');
  assert.equal(new Set([...scene.points.values()].map(p => p.x + ',' + p.y)).size, scene.points.size);
  assert.equal(new Set(scene.edges.map(e => e.from + '>' + e.to)).size, 91, 'deduped recorded and shared first-parent edges');
  assert(!scene.edges.some(e => e.from === 'e1' && e.to === 'e2'), 'positions never invent links');

  const head = Object.fromEntries(scene.heads.map(h => [h.lane, h]));
  const unselected = layout(timeline, {width, height:640});
  assert.equal(head.recent.id, 'r2', 'use rank for identity');
  assert.equal(head.recent.priority, 1, 'use newest recorded time for lane relevance');
  assert.deepEqual([head.selected.priority, head.live.priority, head.attention.priority], [4,3,2]);
  assert(head.recent.y > head.middle.y && head.middle.y > head.selected.y && head.selected.y > head.unknown.y,
    'quiet lanes recede up toward the horizon even when selected');
  assert(head.unknown.width < head.selected.width && head.unknown.height < head.recent.height);
  for (const id of ['e2', 'q1']) {
    const picked = scene.points.get(id), regular = unselected.points.get(id);
    assert.deepEqual([picked.y, picked.scale, picked.depth], [regular.y, regular.scale, regular.depth],
      'selection and attention change lateral prominence, not projected depth or scale');
  }
  assert(scene.points.get('u1').depth > scene.points.get('r2').depth, 'unknown dates do not become fresh');
  const older = scene.points.get('r1');
  const lip = head.recent.y - head.recent.height / 2;
  assert(older.y < lip && lip - older.y < 9 && older.depth > scene.points.get('r2').depth,
    'the first past state sits immediately behind the face');
  assert(scene.points.get('s0').y < scene.points.get('s80').y, 'older shared history recedes towards the horizon');

  for (const [i, a] of scene.heads.entries()) {
    assert(a.x-a.width/2 >= 0 && a.x+a.width/2 <= scene.width);
    assert(a.y-a.height/2 >= 0 && a.y+a.height/2 <= scene.height);
    const {x, y} = scene.points.get(a.id);
    assert.deepEqual([x,y], [a.x,a.y]);
    for (const b of scene.heads.slice(i+1)) {
      assert(Math.abs(a.x-b.x) >= (a.width+b.width)/2 || Math.abs(a.y-b.y) >= (a.height+b.height)/2,
        'head faces overlap: ' + a.id + ' and ' + b.id);
    }
  }
  const center = scene.width/2;
  for (const important of ['selected', 'live']) {
    const a = head[important];
    assert(Math.abs(a.x-center) + a.width/2 <= width/2, 'important heads visible on centered load');
  }
}

const allOld = layout(adaptModel({...model, lanes:model.lanes.slice(7), links:[]}), {width:320, height:400});
assert.equal(allOld.heads.filter(h => h.priority).length, 1);
assert.equal(allOld.heads[0].id, 'o0', 'most recent quiet head is the default priority in an otherwise distant scene');

// Polling can replace the Running summary without rebuilding the graph.
const follow = adaptModel({...model, nowEdge:{open:{status:'complete', live:[{id:'e1', lane:'selected'}]}}});
const running = () => layout(follow).heads.filter(h => h.live).map(h => h.lane);
assert.deepEqual(running(), ['selected'], 'the latest summary supersedes stored session flags');
follow.D.nowEdge.open.live = [];
assert.deepEqual(running(), [], 'an ended session loses its live mark on the next poll');

const side = {x:10, y:40, scale:1, depth:0}, end = {x:100, y:80, scale:.5, depth:2};
assert.deepEqual(turn(side,end,0), side);
assert.deepEqual(turn(side,end,1), end);
const middle = turn(side,end,.5);
assert(middle.x > 55 && middle.x < 100, 'farther points show perspective during the same camera turn');
assert(middle.scale > .5 && middle.scale < 1);
`
