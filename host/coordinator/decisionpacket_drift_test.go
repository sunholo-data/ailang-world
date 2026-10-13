package coordinator

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// The decision-packet drift test binds decisionpacket.go to
// design_docs/sketches/decisionpacket.ail (row 172 AC9): it READS the sketch,
// parses every inline `tests [...]` row and evaluates each through the Go mirror
// of the same name. A sketch law with no mirror, a mirror with no law, a changed
// row or a changed expected value turns this test red.

// packetSketchMirrors maps each sketch law to its Go mirror.
var packetSketchMirrors = map[string]func(t *testing.T, a []any) any{
	"policyCode": func(t *testing.T, a []any) any { return int64(policyCode(argString(t, a, 0))) },
	"scheduleWellFormed": func(t *testing.T, a []any) any {
		return scheduleWellFormed(argInt(t, a, 0), argInt(t, a, 1), argInt(t, a, 2))
	},
	"mintable": func(t *testing.T, a []any) any { return mintable(argInt(t, a, 0), argInt(t, a, 1)) },
	"mintedScheduleWellFormed": func(t *testing.T, a []any) any {
		return mintedScheduleWellFormed(argInt(t, a, 0), argInt(t, a, 1))
	},
	"producerFieldsOk": func(t *testing.T, a []any) any {
		return producerFieldsOk(int(argInt(t, a, 0)), argInt(t, a, 1))
	},
	"sameRef": func(t *testing.T, a []any) any { return sameRef(argLawRef(t, a, 0), argLawRef(t, a, 1)) },
	"packetRefusal": func(t *testing.T, a []any) any {
		return int64(packetRefusal(argBool(t, a, 0), argBool(t, a, 1), argBool(t, a, 2), argBool(t, a, 3),
			argBool(t, a, 4), argBool(t, a, 5), argBool(t, a, 6)))
	},
	"samePair": func(t *testing.T, a []any) any {
		return samePair(argLawRef(t, a, 0), argLawRef(t, a, 1), argLawRef(t, a, 2), argLawRef(t, a, 3))
	},
	"lifecycleOk": func(t *testing.T, a []any) any {
		return lifecycleOk(argInt(t, a, 0), argInt(t, a, 1), argInt(t, a, 2), argInt(t, a, 3), int(argInt(t, a, 4)))
	},
}

// implicationLaws are bool laws whose sketch rows are all true by construction
// (an implication: it holds vacuously outside the mint window). They need >= 2 true rows.
var implicationLaws = map[string]bool{"mintedScheduleWellFormed": true}

// argLawRef reads a sketch Ref record {algo, digest}.
func argLawRef(t *testing.T, a []any, i int) lawRef {
	t.Helper()
	rec, ok := argAt(t, a, i).(map[string]any)
	if !ok || len(rec) != 2 {
		t.Fatalf("arg %d is %T, want a {algo, digest} record", i, argAt(t, a, i))
	}
	algo, okA := rec["algo"].(string)
	digest, okD := rec["digest"].(string)
	if !okA || !okD {
		t.Fatalf("arg %d record %v: algo and digest must be strings", i, rec)
	}
	return lawRef{algo, digest}
}

func TestDecisionPacketSketchRows(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "design_docs", "sketches", "decisionpacket.ail"))
	if err != nil {
		t.Fatalf("read sketch: %v", err)
	}
	if !strings.HasPrefix(string(src), "module sketches/decisionpacket\n") {
		t.Fatal("sketch module header moved")
	}
	laws := parseSketchTests(t, string(src))

	var sketchNames, goNames []string
	for name := range laws {
		sketchNames = append(sketchNames, name)
	}
	for name := range packetSketchMirrors {
		goNames = append(goNames, name)
	}
	sort.Strings(sketchNames)
	sort.Strings(goNames)
	if !reflect.DeepEqual(sketchNames, goNames) {
		t.Fatalf("sketch laws = %v\nGo mirrors = %v\nevery sketch law needs a Go mirror and vice versa", sketchNames, goNames)
	}

	total := 0
	for _, name := range sketchNames {
		rows := laws[name]
		total += len(rows)
		distinct := map[any]bool{}
		sawTrue, sawFalse := false, false
		for _, row := range rows {
			distinct[row.want] = true
			if b, ok := row.want.(bool); ok {
				if b {
					sawTrue = true
				} else {
					sawFalse = true
				}
			}
		}
		switch {
		case name == "policyCode" || name == "packetRefusal":
			if len(distinct) < 2 {
				t.Errorf("int law %s has %d distinct expected values, want >= 2", name, len(distinct))
			}
		case implicationLaws[name]:
			if len(rows) < 2 || sawFalse {
				t.Errorf("implication law %s needs >= 2 rows, all true (have %d, false=%v)", name, len(rows), sawFalse)
			}
		default:
			if !sawTrue || !sawFalse {
				t.Errorf("bool law %s needs a true and a false row (true=%v false=%v)", name, sawTrue, sawFalse)
			}
		}
		for _, row := range rows {
			got := packetSketchMirrors[name](t, row.args)
			if got != row.want {
				t.Errorf("sketch line %d %s%v: Go mirror = %v, sketch expects %v", row.line, name, row.args, got, row.want)
			}
		}
	}
	if total != 37 {
		t.Fatalf("parsed %d sketch rows; the frozen sketch carries exactly 37", total)
	}
	t.Logf("evaluated %d sketch rows over %d laws", total, len(sketchNames))
}
