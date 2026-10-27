package surfacepage

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSurfaceModelLive(t *testing.T) {
	head := "sha256:abc"
	for _, c := range []*LiveCursor{{LastIndex: 1234, LogHead: &head}, {LastIndex: -1}} {
		now := modelNow
		now.Live = c
		m := BuildModel(rigIndex(t), now)
		if m.NowEdge.Live != c {
			t.Fatalf("cursor not carried: %+v", m.NowEdge.Live)
		}
		b, err := json.Marshal(m.NowEdge)
		if err != nil {
			t.Fatal(err)
		}
		if c.LastIndex == -1 && !strings.Contains(string(b), `"live":{"lastIndex":-1,"logHead":null}`) {
			t.Fatalf("empty cursor: %s", b)
		}
	}
	now := modelNow
	now.LiveUnavailable = "the log could not be read"
	m := BuildModel(rigIndex(t), now)
	if m.NowEdge.Live != nil || m.NowEdge.LiveUnavailable != now.LiveUnavailable {
		t.Fatalf("unavailable: %+v", m.NowEdge)
	}
}
