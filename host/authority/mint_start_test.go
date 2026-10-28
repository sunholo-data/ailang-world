package authority

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestMintStartedKeepsHistoryAndRollsBackTogether(t *testing.T) {
	c, err := store.OpenCredentials(filepath.Join(t.TempDir(), "credentials.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	start := &store.EpisodeStart{Schema: "world/episode-start/v1", Episode: "ep"}
	grants := []broker.Capability{{Effect: broker.EffectWorkspaceRead, Scope: "worktree", Budget: 1}}
	_, id, _, err := MintStarted(boundedTestContext(t), c, "ep", ProfileBinding{}, grants, 3600, 1000, nil, start)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := MintStarted(boundedTestContext(t), c, "ep", ProfileBinding{}, grants, 3600, 1000, nil, start); err == nil {
		t.Fatal("duplicate start accepted")
	}
	rows, err := c.ListSessions(boundedTestContext(t), "ep", 0)
	if err != nil || len(rows) != 1 {
		t.Fatalf("%+v %v", rows, err)
	}
	if err := Revoke(boundedTestContext(t), c, id); err != nil {
		t.Fatal(err)
	}
	got, err := c.EpisodeStart(boundedTestContext(t), "ep")
	if err != nil || !reflect.DeepEqual(got, start) {
		t.Fatalf("%+v %v", got, err)
	}
}
