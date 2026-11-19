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
const {layout, turn, project, block, paneSize} = require(process.argv[2]);
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
function checkLaneSilhouettes(scene, graph) {
  const g = scene.ground, heads = new Map(scene.heads.map(h=>[h.lane,h]));
  const drawn = [...scene.points].map(([id,p]) => {
    const lane = graph.N[id].lane, head = heads.get(lane);
    const front = head.id === id;
    const {width,height} = front ? head : paneSize(g,head,p);
    assert(width>0 && height>0, 'even the oldest record has a face');
    assert(p.blockDepth>0, 'even the oldest record has positive extrusion');
    const {front:near,back} = block(g,p,width,height,p.blockDepth);
    const corners = [...near,...back];
    const xs = corners.map(v=>v.x), ys = corners.map(v=>v.y);
    return {id,corners,l:Math.min(...xs),r:Math.max(...xs),t:Math.min(...ys),b:Math.max(...ys)};
  });
  for (const [i,a] of drawn.entries()) for (const b of drawn.slice(i+1)) {
    if (a.r<=b.l || b.r<=a.l || a.b<=b.t || b.b<=a.t) continue;
    // In one-point perspective the remaining silhouette edges all run
    // toward the vanishing point. Check those separating axes as well.
    const separate = [...a.corners,...b.corners].some(corner => {
      const dx = corner.x-g.vx, dy = corner.y-g.horizon;
      const normal = p => p.x*dy-p.y*dx;
      const left = a.corners.map(normal), right = b.corners.map(normal);
      return Math.max(...left)<=Math.min(...right) || Math.max(...right)<=Math.min(...left);
    });
    assert(separate,'record silhouettes intersect: '+a.id+' and '+b.id);
  }
}
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
  assert.equal(head.recent.id, 'r1', 'use the last timeline instant for the head');
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
  const {ground:g,time} = scene;
  assert.equal(time.referenceIndex,timeline.presentIndex);
  assert.equal(time.referenceX,timeline.LX[1][timeline.presentIndex]);
  assert.equal(scene.travelZ,0);
  assert.equal(time.depthAt(time.referenceIndex),1);
  const tail = scene.points.get('s0');
  assert(tail.depth>12 && tail.y-g.horizon<32,'long Timeline history tapers toward the horizon');

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
    assert.equal(shape.depth,face.depth);
    assert(nearEqual(face.scale*face.depth,1));
    assert(nearEqual(shape.backDepth-shape.depth,face.blockDepth));
    assert(face.blockDepth>0 && face.blockDepth<=ground.blockDepth);
    const groundCenter = project(ground,(face.x-center)*face.depth,0,face.depth);
    assert(nearEqual(groundCenter.x,face.x) && nearEqual(groundCenter.y,face.y+ground.eventElevation/face.depth));
    if (face.lane !== ground.centerLane) {
      assert(shape.front.concat(shape.back).every(p=>Math.sign(p.x-center)===Math.sign(face.x-center)),
        'branch silhouette stays on its bank');
    }
    for (const [i,p] of shape.front.entries()) checkRay(p,shape.back[i],ground);
  }
  for (const [id,point] of scene.points) {
    assert(nearEqual(point.scale*point.depth,1),'inverse depth controls scale');
    const head = fronts.get(timeline.N[id].lane);
    if (head.worldX===point.worldX) checkRay(head,point,ground);
    const expected = project(ground,point.worldX,ground.eventElevation,point.depth);
    assert(nearEqual(expected.x,point.x) && nearEqual(expected.y,point.y),'record retains fixed elevation');
    const index = timeline.timeIndex.get(id);
    assert.equal(point.timeIndex,index);
    assert(nearEqual(point.depth,1+(time.referenceX-timeline.LX[1][index])/time.unitsPerZ));
    assert.equal(point.y,time.yAt(index));
  }
  assert.equal(ground.centerLane,'shared');
  assert.equal(head.shared.x,center,'shared world sits on the optical axis');
  assert([...scene.points].filter(([id,p])=>timeline.N[id].lane==='shared' && p.worldX===0).every(([,p])=>p.x===center));
  for (const face of [head.shared,head.recent]) {
    assert(face.width / face.height >= 1.9,'foreground faces are landscape');
  }
  assert.notEqual(head.shared.y,head.recent.y,'foreground faces stagger in recorded depth');
  assert(ground.streetY>Math.max(...[head.shared,head.recent].map(h=>h.y+h.height/2)), 'the present plane is now, after latest activity');
  for (const face of scene.heads) {
    assert(face.height>0 && face.height*face.depth<=ground.worldHeight+1e-8,'fitted facades retain a positive height within the world height');
    const relativeWidth = (face.width*face.depth)/(head.shared.width*head.shared.depth);
    assert(relativeWidth>0.99 && relativeWidth<=1.10000001,'far facades retain plausible world widths');
    if (face.lane !== 'shared') {
      assert((Math.abs(face.x-center)-face.width/2)*face.depth>=ground.street/2-1e-8,
        'branch faces stay beyond the curb');
    }
  }
  checkLaneSilhouettes(scene,timeline);
  for (const important of ['selected', 'live']) {
    const a = head[important];
    assert(Math.abs(a.x-center) + a.width/2 <= width/2, 'important heads visible on centered load');
  }
}

for (const width of [320, 390, 1280]) {
  const single = layout(adaptModel({...model,lanes:[model.lanes[1]],links:[]}),{width,height:620});
  const only = single.heads[0], center = single.width/2;
  assert(Math.abs(only.x-center)+only.width/2 <= width/2,'a single foreground face is readable on arrival');
  assert.equal(single.ground.centerLane,only.lane,'branch-only worlds still have a camera reference');
  assert.equal(only.x,center);
}

const recentLanes = Array.from({length:12},(_, i)=>lane('recent-'+i,[state('new'+i,200+i,.5)]));
const threeNear = layout(adaptModel({...model,lanes:recentLanes.slice(0,3),links:[]}),{width:1280,height:780});
assert.equal(threeNear.width,1280,'all three near faces fit in the desktop frame');
assert(threeNear.heads.every(h=>h.width*h.depth>=240 && h.width*h.depth<=250));
const busy = layout(adaptModel({...model,lanes:recentLanes,links:[]}),{width:1280,height:780});
assert(busy.width>1280 && busy.heads.every(h=>h.width*h.depth>=171.999999),'busy desktop can pan at a readable width');
const mobile = layout(adaptModel({...model,lanes:recentLanes.slice(0,3),links:[]}),{width:320,height:640});
assert(mobile.heads.every(h=>nearEqual(h.width*h.depth,122) && nearEqual(h.height*h.depth,61)),
  'desktop fitting leaves the reference facade size on phones alone');

const longLane = lane('long',Array.from({length:420},(_,i)=>state('long'+i,i,420-i)));
const long = layout(adaptModel({...model,lanes:[longLane],links:[]}),{width:320,height:640});
assert.equal(new Set([...long.points.values()].map(p=>p.depth)).size,420);
assert(long.points.get('long0').depth>60 && long.points.get('long0').y-long.ground.horizon<8);

const exact = adaptModel({...model,lanes:[
  lane('shared',[state('old',1,3,{kind:'shared'}),state('tie1',4,1,{kind:'shared'}),state('tie2',5,1,{kind:'shared'}),
    state('latest',8,.2,{kind:'shared'})]),
  lane('b',[state('b-old',0,10),state('b-at-tie',2,1),state('b-after',3,.5)]),
  lane('unknown',[state('undated',12,null)]),
  lane('future',[state('future',99,-1)]),
  lane('draft',[state('local',100,99,{draft:true,kind:'draft'})])
],links:[{from:'b-at-tie',to:'tie1'},{from:'b-at-tie',to:'b-after'}]});
for (const width of [320,1280]) {
  for (const timeXs of [exact.LX[1],exact.LX[3]]) {
    const base = layout(exact,{width,height:650,timeXs,attention:new Set(['b'])});
    assert(base.points.has('local') && !base.points.has('future'));
    assert.equal(base.points.get('local').depth,1,'draft fixes to the real present');
    const first = base.points.get('tie1'), second = base.points.get('tie2'), other = base.points.get('b-at-tie');
    assert.deepEqual([first.depth,first.y],[second.depth,second.y]);
    assert.deepEqual([second.depth,second.y],[other.depth,other.y]);
    assert.notEqual(first.x,second.x,'same-lane tied records use lateral placements');
    const cut = layout(exact,{width,height:650,timeXs,base,stop:'tie1',attention:new Set(['b'])});
    assert.equal(cut.heads.find(h=>h.lane==='shared').id,'tie1');
    assert.equal(cut.heads.find(h=>h.lane==='b').id,'b-at-tie');
    assert.deepEqual([...cut.points.keys()].sort(),['b-at-tie','b-old','old','tie1','tie2','undated']);
    assert(cut.heads.every(h=>!h.live && !h.attention));
    assert.equal(cut.points.get('tie1').depth,1);
    assert(nearEqual(cut.travelZ,base.points.get('tie1').depth-1));
    assert(cut.edges.some(e=>e.from==='b-at-tie'&&e.to==='tie1'));
    for (const [id,p] of cut.points) assert(nearEqual(base.points.get(id).depth-p.depth,cut.travelZ));
    checkLaneSilhouettes(base,exact);
    checkLaneSilhouettes(cut,exact);

    const shift = layout(exact,{width,height:650,timeXs,base,viewLane:'b',selected:'b'});
    assert.equal(shift.ground.centerLane,'b');
    assert.deepEqual(shift.laneOrder,base.laneOrder);
    assert.equal(shift.heads.find(h=>h.lane==='b').x,shift.ground.vx);
    assert.deepEqual([...shift.points].map(([id,p])=>[id,p.depth,p.y]).sort(),
      [...base.points].map(([id,p])=>[id,p.depth,p.y]).sort());
    for (let i=1;i<shift.laneOrder.length;i++) {
      const prev=shift.heads.find(h=>h.lane===shift.laneOrder[i-1]);
      const next=shift.heads.find(h=>h.lane===shift.laneOrder[i]);
      assert(prev.worldX<=next.worldX,'recenter keeps the left-to-right order');
    }
    checkLaneSilhouettes(shift,exact);
  }
}
const branchOnly = layout(adaptModel({...model,lanes:[
  lane('oldbranch',[state('oldbranch',1,20)]),lane('mainbranch',[state('mainbranch',3,.1)]),
  lane('otherbranch',[state('otherbranch',2,7)])
],links:[]}),{width:320,height:640});
assert.equal(branchOnly.ground.centerLane,'mainbranch');
assert.equal(branchOnly.laneOrder.indexOf('mainbranch'),1,'a branch-only default permits lane switching on either side');

const history = adaptModel({...model,lanes:[
  lane('shared',[state('main',20,2,{kind:'shared'})]),
  lane('left',[state('left-now',21,1),state('left-gap',14,7),state('left-old',3,15)]),
  lane('road',[state('road-now',22,.5),state('road-tie-2',11,9),state('road-tie-1',10,9),state('road-old',2,16)]),
  lane('right',[state('right-new',16,5),state('right-old',1,18)]),
  lane('late',[state('late',19,3)])
],links:[]});
const historyOptions = {width:390,height:640};
const present = layout(history,historyOptions);
assert.deepEqual(present.laneOrder,[...present.plan.lanes.keys()].sort((a,b)=>
  present.plan.lanes.get(a).origin-present.plan.lanes.get(b).origin));
for (const center of ['left','shared','road','late']) {
  for (const stop of ['right-new','left-gap','road-tie-2','road-tie-1','left-old']) {
    const options = {...historyOptions,stop,viewLane:center};
    const past = layout(history,{...options,base:present});
    const automatic = layout(history,options);
    assert.equal(past.ground.centerLane,center,'the camera keeps its center through gaps');
    assert.deepEqual(past.laneOrder,present.laneOrder,'cuts retain the present lane order');
    const positions = scene => [...scene.heads].sort((a,b) => a.worldX - b.worldX).map(h=>[h.lane,h.worldX]);
    assert.deepEqual(positions(automatic),positions(past),'an omitted base uses the same lane positions');
    if (center === 'shared') {
      const defaultView = layout(history,{...historyOptions,base:present,stop});
      assert.equal(defaultView.ground.centerLane,'shared','the default reference survives cuts without shared records');
      assert.deepEqual(positions(defaultView),positions(past));
    }
    const visible = new Set(past.heads.map(h=>h.lane));
    assert.deepEqual(positions(past).map(([lane])=>lane),present.laneOrder.filter(lane=>visible.has(lane)));
    const centered = layout(history,{...historyOptions,base:present,viewLane:center});
    for (const [id,point] of past.points) {
      assert.equal(point.worldX,centered.points.get(id).worldX,id+' retains its track at '+stop);
    }
    checkLaneSilhouettes(past,history);
  }
}

const clustered = adaptModel({...model,lanes:[lane('shared',[
  state('newest',9,1),state('newer-2',8,4),state('newer-1',7,4),state('newer-0',6,4),
  ...Array.from({length:6},(_,i)=>state('cluster-'+i,i,5))
])],links:[]});
const clusterBase = layout(clustered,historyOptions);
let companions = 0;
for (const stop of [undefined,'newer-2','newer-1','newer-0','cluster-5','cluster-3','cluster-0']) {
  const past = layout(clustered,{...historyOptions,base:clusterBase,stop});
  for (const [id,point] of past.points) {
    assert.equal(point.worldX,clusterBase.points.get(id).worldX,id+' keeps its exact track at '+stop);
    if (!stop && point.worldX !== 0) companions++;
  }
  checkLaneSilhouettes(past,clustered);
}
assert(companions >= 5,'the test exercises several companions');

const consecutive = adaptModel({...model,lanes:[
  lane('shared',[
    state('merge-old',1,3,{kind:'shared'}),
    state('merge-paired',2,2,{kind:'shared'}),
    state('source',10,2),
    state('merge-new',4,1,{kind:'shared'})
  ]),
  lane('branch',[state('step-old',3,3),state('step-paired',5,2),state('step-new',6,1)])
],links:[{from:'source',to:'merge-paired'}]});
const primaries = {
  shared:['merge-old','merge-paired','merge-new'],
  branch:['step-old','step-paired','step-new']
};
for (const width of [320,1280]) {
  const options = {width,height:640};
  const base = layout(consecutive,options);
  for (const viewLane of Object.keys(primaries)) {
    const present = layout(consecutive,{...options,base,viewLane});
    const axis = present.plan.lanes.get(viewLane).origin;
    assert.equal(present.ground.centerLane,viewLane);
    for (const [key,ids] of Object.entries(primaries)) {
      const x = (present.plan.lanes.get(key).origin-axis)*present.plan.pitch;
      for (const id of ids) assert.equal(present.points.get(id).worldX,x,id+' follows its lane axis');
      const primary = present.points.get(ids.at(-1));
      if (key === viewLane) assert.equal(primary.x,present.ground.vx);
    }
    const merge = present.points.get('merge-paired'), source = present.points.get('source');
    assert.equal(source.depth,merge.depth);
    assert.equal(source.y,merge.y);
    assert(Math.abs(source.worldX-merge.worldX)>0 && Math.abs(source.worldX-merge.worldX)<=present.plan.pitch,
      'paired source uses the adjacent track');
    for (const stop of Object.keys(consecutive.N)) {
      const past = layout(consecutive,{...options,base,viewLane,stop});
      const head = past.heads.find(h=>h.id===stop);
      assert(head,stop+' is promoted to a head at its stop');
      assert.equal(head.depth,1);
      assert.equal(head.x,past.ground.vx+present.points.get(stop).worldX);
      assert(head.width>=(stop==='source'?80:122) && head.height>=40,stop+' retains a readable near-plane face');
      if (viewLane==='shared' && stop==='merge-paired') {
        const peer = past.points.get('source');
        assert(Math.abs(peer.x-past.ground.vx)+peer.paneWidth/2 <= width/2,
          'the nearest pair fits together in the viewport');
      }
      const cutoff = consecutive.timeIndex.get(stop);
      assert.deepEqual([...past.points.keys()].sort(),[...present.points.keys()].filter(
        id=>consecutive.timeIndex.get(id)<=cutoff).sort());
      for (const [id,point] of past.points) {
        assert.equal(point.worldX,present.points.get(id).worldX,id+' keeps its track at '+stop+' from '+viewLane);
      }
      if (stop==='merge-paired' || stop==='source') {
        assert(past.edges.some(e=>e.from==='source' && e.to==='merge-paired'));
      }
      checkLaneSilhouettes(past,consecutive);
    }
    checkLaneSilhouettes(present,consecutive);
  }
}


// Four seconds can be a fraction of a Timeline pixel at Weeks zoom.
// Advancing to that instant must still yield a usable foreground card.
const dense = adaptModel({...model,lanes:[lane('shared',[12,8,4].map((seconds,i)=>
  state('fast-'+i,i,seconds/86400,{kind:'shared'})))],links:[]});
const denseOptions = {width:320,height:640,timeXs:dense.LX[4]};
const denseBase = layout(dense,denseOptions);
for (const stop of ['fast-2','fast-1','fast-0']) {
  const scene = layout(dense,{...denseOptions,base:denseBase,stop});
  const head = scene.heads[0];
  assert.equal(head.id,stop);
  assert(head.height>=44 && head.worldX===0);
  checkLaneSilhouettes(scene,dense);
}

const repoLane = (repo,name,kind,age) => ({
  ...lane(repo+':brane:'+name,Array.from({length:65},(_,i)=>state(repo+name+i,i,age+64-i,{kind}))),repo
});
const following = adaptModel({...model,repos:[
  {id:'b',shared:'b:brane:shared'},{id:'a',shared:'a:brane:shared'}
],lanes:[
  repoLane('a','shared','shared',.5),repoLane('b','shared','shared',.5),
  repoLane('a','branch','step',8),repoLane('b','branch','step',10)
],links:[]});
for (const width of [320,1280]) {
  const multi = layout(following,{width,height:640});
  assert.equal(multi.ground.centerLane,'b:brane:shared','first configured repo receives the optical axis');
  assert.equal(multi.heads.find(h=>h.lane==='b:brane:shared').x,multi.ground.vx);
  assert.notEqual(multi.heads.find(h=>h.lane==='a:brane:shared').x,multi.ground.vx);
  assert.equal(multi.points.size,260);
  checkLaneSilhouettes(multi,following);
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
