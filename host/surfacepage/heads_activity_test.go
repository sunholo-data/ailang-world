package surfacepage

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestHeadsActivityTrace(t *testing.T) {
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	script, err := filepath.Abs(filepath.Join("app", "heads-activity.js"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, "-e", headsActivityHarness, script)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("heads activity: %v\n%s", err, out)
	}
}

const headsActivityHarness = `
const assert = require('node:assert/strict');
const {trace} = require(process.argv[1]);
const now = Date.parse('2026-11-19T15:00:00Z');
const states = [
  {rank: 100},
  {rank: 99, at: 'invalid'},
  {rank: 3, at: '2026-10-21T00:00:00Z'},
  {rank: 4, at: '2026-10-20T23:59:59.999Z'},
  {rank: 2, at: '2026-11-19T01:00:00.123Z'},
  {rank: 5, at: '2026-11-18T19:00:00-05:00'},
  {rank: 6, at: '2026-11-19T00:20:00+01:00'},
  {rank: 7, at: '2026-11-19T15:00:00.001Z'},
  {rank: 8, at: '2026-11-20T00:00:00Z'},
];
const days = Array(30).fill(0);
days[0] = 1;
days[28] = 1;
days[29] = 2;
const expected = {days, last: Date.parse('2026-11-19T01:00:00.123Z')};
assert.deepEqual(trace(states, now), expected);
assert.deepEqual(trace(states.slice().reverse(), now), expected);
assert.deepEqual(trace([states[3]], now), {days: Array(30).fill(0), last: Date.parse(states[3].at)});
assert.deepEqual(trace(states.slice(0, 2), now), {days: Array(30).fill(0), last: null});
assert.deepEqual(trace([], now), {days: Array(30).fill(0), last: null});
`
