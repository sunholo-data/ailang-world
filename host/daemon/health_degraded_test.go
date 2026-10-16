package daemon

// Row 185 (w-episode-requires-git): /v1/health names the degraded episodes —
// episode directories under --workspace-root that are not git worktrees
// (`session new --allow-non-git`) — and adds nothing when there is none.

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func healthBody(t *testing.T, d *Daemon) (HealthResponse, string) {
	t.Helper()
	rec := requestRecorder(t, d, http.MethodGet, "/v1/health", nil)
	code, raw := rec.Code, rec.Body.String()
	if code != http.StatusOK {
		t.Fatalf("GET /v1/health = %d (body %q), want 200", code, raw)
	}
	var h HealthResponse
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		t.Fatalf("decode health %q: %v", raw, err)
	}
	return h, raw
}

// The fixture root holds ep1 and ep2 (plain directories), EP1 (outside the
// grammar), a plain file and three symlinks. ep1 is made a linked worktree;
// only ep2 is degraded. Kills MUT-HEALTH-DEGRADED-SILENT (an empty list) and
// MUT-HEALTH-DEGRADED-ALL (listing every directory, git or not).
func TestHealthNamesDegradedEpisodes(t *testing.T) {
	f := newWSFixture(t)
	writeFile(t, filepath.Join(f.root, "ep1", ".git"), "gitdir: /nowhere/.git/worktrees/ep1\n")
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root})
	h, _ := healthBody(t, d)
	if want := []string{"ep2"}; !reflect.DeepEqual(h.DegradedEpisodes, want) {
		t.Fatalf("degraded_episodes = %q, want %q", h.DegradedEpisodes, want)
	}

	// Once ep2 is a worktree too, the key disappears.
	writeFile(t, filepath.Join(f.root, "ep2", ".git"), "gitdir: /nowhere/.git/worktrees/ep2\n")
	if _, raw := healthBody(t, d); strings.Contains(raw, "degraded_episodes") {
		t.Fatalf("health with no degraded episode carries the key: %s", raw)
	}
}

// Without --workspace-root the body is unchanged: no degraded_episodes key.
func TestHealthWithoutAWorkspaceRootHasNoDegradedKey(t *testing.T) {
	db := filepath.Join(canonicalTemp(t), "world.db")
	d := mustWSDaemon(t, Config{DBPath: db})
	if _, raw := healthBody(t, d); strings.Contains(raw, "degraded_episodes") {
		t.Fatalf("health without a workspace root carries degraded_episodes: %s", raw)
	}
	if _, err := os.Stat(db); err != nil {
		t.Fatalf("store: %v", err)
	}
}
