package surfacepage

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestHistoricalHeadSliceRetainsTheScene(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for the historical head slice test")
	}
	script := func(name string) string {
		t.Helper()
		path, err := filepath.Abs(filepath.Join("app", name))
		if err != nil {
			t.Fatal(err)
		}
		return path
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, node, "-e", headTravelHarness, script("app.js"), script("head-camera.js"), script("head-lineage.js")).CombinedOutput()
	if err != nil {
		t.Fatalf("historical head slice: %v\n%s", err, out)
	}
}

const headTravelHarness = `
const assert = require('node:assert/strict');
const {adaptModel} = require(process.argv[1]);
const camera = require(process.argv[2]);
const state = (id,rank,age,extra={}) => ({id,kind:'step',rank,at:new Date(Date.UTC(2026,10,19-age)).toISOString(),...extra});
const model = adaptModel({
  now:'2026-11-19T12:00:00Z',
  lanes:[
    {key:'shared',states:[state('s0',0,24,{kind:'shared'}),state('s1',3,18,{kind:'shared'}),state('s2',5,12,{kind:'shared'})]},
    {key:'a',states:[state('a0',1,24),state('a1',4,14),state('a-tied-later',4,14),state('a2',10,9)]},
    {key:'b',states:[state('b0',2,15),state('b1',4,12),state('b-tied-latest',4,20),state('b2',11,13)]},
    {key:'later',states:[state('c0',7,21)]},
    {key:'draft',states:[state('d0',14,0,{kind:'draft',draft:true})]}
  ],
  links:[
    {from:'s0',to:'a0'},{from:'a0',to:'a1'},{from:'a1',to:'a-tied-later'},
    {from:'a1',to:'b-tied-latest'},{from:'s1',to:'a2'},{from:'b2',to:'c0'}
  ]
});
const {forwardPaths} = require(process.argv[3]);
const cut = model.timeIndex.get('s0');
assert.deepEqual(forwardPaths(model,'s0',cut), [['s1']],
  'only recorded successors beyond the cut are eligible for the forward stack');
assert.deepEqual(forwardPaths(model,'a1',model.timeIndex.get('a1')), [],
  'a simultaneous neighbor and an earlier child belong to the visible side of the cut');
for (const width of [320, 1280]) {
  const options = {width,height:640,selected:'a',attention:new Set(['b'])};
  const base = camera.layout(model,options);
  const slice = camera.layout(model,{...options,base,stop:'a1'});
  assert.equal(slice.travelZ,base.points.get('a1').depth-slice.points.get('a1').depth);
  assert(slice.travelZ>0,'the camera advances by the actual time between now and the stop');
  assert.deepEqual(slice.heads.map(h=>h.id).sort(),['a1','b0','c0','s1']);
  assert.deepEqual([...slice.points.keys()].sort(),
    ['a-tied-later','a0','a1','b-tied-latest','b0','c0','s0','s1'],
    'includes simultaneous and earlier states regardless of their ranks');
  for (const [id,point] of slice.points) {
    const index = model.timeIndex.get(id);
    assert(Math.abs(point.depth - (base.points.get(id).depth - slice.travelZ))<1e-8,
      'every surviving record moves through the same physical distance');
    assert(Math.abs(point.y - slice.time.yAt(index))<1e-8, 'the time ruler marks the actual card center');
  }
  assert.equal(slice.points.get('a1').y,slice.points.get('a-tied-later').y);
  assert.notEqual(slice.points.get('a1').x,slice.points.get('a-tied-later').x);
  for (const head of slice.heads) {
    assert(!head.live && !head.attention, 'a historical slice clears current session/approval signals');
  }
  const edges = slice.edges.map(e=>[e.from,e.to]);
  assert(edges.some(([from,to])=>from==='s0' && to==='s1'), 'shared chronology survives travel');
  assert(edges.every(([from,to])=>slice.points.has(from) && slice.points.has(to)));
  assert(base.points.has('d0') && base.points.has('s2'),'the present scene remains available for return');
}
`
