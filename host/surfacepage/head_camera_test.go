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
const {layout, turn, project, block} = require(process.argv[2]);
const nearEqual = (a,b) => Math.abs(a-b)<1e-7;
const checkRay = (a,b,g) => assert(nearEqual(
  (a.x-g.vx)*(b.y-g.horizon), (b.x-g.vx)*(a.y-g.horizon)), 'points share the camera vanishing point');
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
  assert(older.y < lip && lip - older.y < 25 && older.depth > scene.points.get('r2').depth,
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
  const ground = scene.ground, center = scene.width/2;
  assert.equal(ground.vx,center);
  assert(nearEqual(ground.focal,640/(2*Math.tan(27*Math.PI/180))));
  const fronts = new Map(scene.heads.map(h=>[h.lane,h]));
  for (const face of scene.heads) {
    const shape = block(ground,face,face.width,face.height);
    const inferred = block(ground,{x:face.x,y:face.y},face.width,face.height);
    assert(nearEqual(shape.depth,inferred.depth) && nearEqual(shape.depth,face.depth));
    assert(nearEqual(face.scale*face.depth,1));
    assert(nearEqual(shape.backDepth-shape.depth,ground.blockDepth));
    const groundCenter = project(ground,(face.x-center)*face.depth,0,face.depth);
    assert(nearEqual(groundCenter.x,face.x) && nearEqual(groundCenter.y,face.y+face.height/2));
    assert(shape.front.concat(shape.back).every(p=>Math.sign(p.x-center)===Math.sign(face.x-center)),
      'the solid building leaves the central view open');
    for (const [i,p] of shape.front.entries()) checkRay(p,shape.back[i],ground);
  }
  for (const [id,point] of scene.points) {
    assert(nearEqual(point.scale*point.depth,1),'inverse depth controls scale');
    const head = fronts.get(timeline.N[id].lane);
    checkRay(head,point,ground);
    const elevation = head.height*head.depth/2;
    const expected = project(ground,(head.x-center)*head.depth,elevation,point.depth);
    assert(nearEqual(expected.x,point.x) && nearEqual(expected.y,point.y),'lane histories retain fixed physical X/elevation');
  }
  const foreground = [head.recent, head.shared].sort((a,b) => a.x-b.x);
  const [left,right] = foreground;
  assert(left.x < center && right.x > center);
  assert(left.width / left.height >= 1.9 && right.width / right.height >= 1.9, 'foreground faces are landscape');
  assert.notEqual(left.y,right.y,'foreground faces stagger in recorded depth');
  assert(nearEqual((center-left.x-left.width/2)*left.depth,ground.street/2));
  assert(nearEqual((right.x-right.width/2-center)*right.depth,ground.street/2));
  assert.equal(ground.streetY, Math.max(...foreground.map(h=>h.y+h.height/2)));
  for (const face of scene.heads) {
    assert(nearEqual(face.height*face.depth,left.height*left.depth),'facades share the same world height');
    const relativeWidth = (face.width*face.depth)/(left.width*left.depth);
    assert(relativeWidth>0.99 && relativeWidth<=1.10000001,'far facades retain plausible world widths');
  }
  for (const important of ['selected', 'live']) {
    const a = head[important];
    assert(Math.abs(a.x-center) + a.width/2 <= width/2, 'important heads visible on centered load');
  }
}

for (const width of [320, 390, 1280]) {
  const single = layout(adaptModel({...model,lanes:[model.lanes[1]],links:[]}),{width,height:620});
  const only = single.heads[0], center = single.width/2;
  assert(Math.abs(only.x-center)+only.width/2 <= width/2,'a single foreground face is readable on arrival');
  assert(Math.abs(only.x-center)-only.width/2 > 0,'a lone lane leaves the through-path open');
}

const recentLanes = Array.from({length:12},(_, i)=>lane('recent-'+i,[state('new'+i,200+i,.5)]));
const threeNear = layout(adaptModel({...model,lanes:recentLanes.slice(0,3),links:[]}),{width:1280,height:780});
assert.equal(threeNear.width,1280,'all three near faces fit in the desktop frame');
assert(threeNear.heads.every(h=>h.width>=226 && h.width<=250 && h.height>=80 && h.height<=87));
const busy = layout(adaptModel({...model,lanes:recentLanes,links:[]}),{width:1280,height:780});
assert(busy.width>1280 && busy.heads.every(h=>h.width*h.depth>=171.999999),'busy desktop can pan at a readable width');
const mobile = layout(adaptModel({...model,lanes:recentLanes.slice(0,3),links:[]}),{width:320,height:640});
assert(nearEqual(Math.max(...mobile.heads.map(h=>h.width)),122) &&
  nearEqual(Math.max(...mobile.heads.map(h=>h.height)),61),'desktop fitting leaves the nearest phone face alone');

const allOld = layout(adaptModel({...model, lanes:model.lanes.slice(7), links:[]}), {width:320, height:400});
assert.equal(allOld.heads.filter(h => h.priority).length, 1);
assert.equal(allOld.heads[0].id, 'o0', 'most recent quiet head is the default priority in an otherwise distant scene');

// Polling can replace the Running summary without rebuilding the graph.
const follow = adaptModel({...model, nowEdge:{open:{status:'complete', live:[{id:'e1', lane:'selected'}]}}});
const running = () => layout(follow).heads.filter(h => h.live).map(h => h.lane);
assert.deepEqual(running(), ['selected'], 'the latest summary supersedes stored session flags');
follow.D.nowEdge.open.live = [];
assert.deepEqual(running(), [], 'an ended session loses its live mark on the next poll');

const testGround = {vx:160,horizon:70,streetY:310,blockDepth:.2};
assert.deepEqual(project(testGround,0,0,1),{x:160,y:310},'Z=1 lands at the street reference');
assert.deepEqual(project(testGround,100,120,2),{x:210,y:130});
const solid = block(testGround,{x:235,y:200},80,60);
const factor = solid.depth/solid.backDepth;
assert(nearEqual(solid.back[1].x-solid.back[0].x,80*factor));
assert(nearEqual(solid.back[2].y-solid.back[1].y,60*factor));
assert(solid.top && solid.side && solid.floor);
assert.equal(solid.side[0],solid.front[0],'right bank exposes the side facing inward');
const opposite = block(testGround,{x:85,y:200},80,60);
assert.equal(opposite.side[0],opposite.front[1],'left bank exposes the opposite face');
const straddled = block(testGround,{x:160,y:80},60,60);
assert.equal(straddled.top,null,'roof above eye height is hidden');
assert.equal(straddled.side,null,'a face spanning the vanishing point hides both sides');
assert.equal(block(testGround,{x:160,y:40},40,40),null,'front feet above the horizon have no inferred ground depth');
const raised = block(testGround,{x:235,y:40,depth:2},60,24);
assert.equal(raised.depth,2,'an elevated pane retains its actual camera Z');
assert.equal(raised.floor[0].y,190,'the elevated pane footprint still lies on the ground');
assert(raised.front[3].y < testGround.horizon && raised.floor[0].y > testGround.horizon);
checkRay(raised.front[0],raised.back[0],testGround);

const side = {x:10, y:40, scale:1, depth:0}, end = {x:100, y:80, scale:.5, depth:2};
assert.deepEqual(turn(side,end,0), side);
assert.deepEqual(turn(side,end,1), end);
const middle = turn(side,end,.5);
assert(middle.x > 55 && middle.x < 100, 'farther points show perspective during the same camera turn');
assert(middle.scale > .5 && middle.scale < 1);
assert(nearEqual(middle.depth*middle.scale,1),'camera turn carries a consistent depth and inverse scale');
`
