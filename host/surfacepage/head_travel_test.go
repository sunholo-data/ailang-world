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
	out, err := exec.CommandContext(ctx, node, "-e", headTravelHarness, script("app.js"), script("head-camera.js"), script("head-travel.js")).CombinedOutput()
	if err != nil {
		t.Fatalf("historical head slice: %v\n%s", err, out)
	}
}

const headTravelHarness = `
const assert = require('node:assert/strict');
const {adaptModel} = require(process.argv[1]);
const camera = require(process.argv[2]);
const travel = require(process.argv[3]);
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
for (const width of [320, 1280]) {
  const base = camera.layout(model,{width,height:640,selected:'a',attention:new Set(['b'])});
  const initial = [...base.points];
  const slice = travel.project(model,base,'a','a1');
  assert.deepEqual([slice.width,slice.height],[base.width,base.height],'retains scene bounds');
  assert.deepEqual([...base.points],initial,'whole-scene geometry remains reusable');
  assert.deepEqual(slice.heads.map(h=>h.id).sort(),['a1','b-tied-latest','s1']);
  assert.equal(slice.heads[0].id,'a1','the chosen exact stop is foreground');
  assert.deepEqual([...slice.points.keys()].sort(),['a0','a1','b-tied-latest','b0','b1','s0','s1']);
  for (const head of slice.heads) {
    const original = base.heads.find(h=>h.lane===head.lane);
    const point = slice.points.get(head.id);
    assert.deepEqual([head.x,head.y,point.x,point.y],[original.x,original.y,original.x,original.y],'lane stays anchored');
    assert.equal(head.attention,false,'present-only attention is cleared');
    assert.equal(head.live,false,'present-only running status is cleared');
  }
  const prior = slice.points.get('a0'), head = slice.heads[0];
  assert(prior.y < head.y - head.height/2,'older state stays just behind the head');
  assert.equal(slice.edges.length,4);
  assert(slice.edges.some(e=>e.from==='s0'&&e.to==='s1'),'shared chronological line survives');
  assert(slice.edges.every(e=>base.edges.includes(e) && slice.points.has(e.from) && slice.points.has(e.to)));
  assert.equal(travel.project(model,base,'a','b0'),base,'a foreign stop leaves the view intact');
  assert.equal(travel.project(model,base,'draft','d0'),base,'local drafts are not historical stops');
}
`
