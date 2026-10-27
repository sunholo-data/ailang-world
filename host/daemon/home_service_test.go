package daemon

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/authority"
)

func TestHomeEagerIdentityAC15(t *testing.T) {
	path := filepath.Join(t.TempDir(), "binary")
	original := []byte("original running binary")
	if err := os.WriteFile(path, original, 0700); err != nil {
		t.Fatal(err)
	}
	old := currentBuild
	currentBuild = sync.OnceValue(func() BuildIdentity {
		info, _ := debug.ReadBuildInfo()
		return buildIdentityOf(info, executableDigestAt(path))
	})
	defer func() { currentBuild = old }()
	// Exercise Run's first startup action without opening any store or socket.
	_ = Run(context.Background(), Config{BindHost: "0.0.0.0"}, io.Discard)
	other := path + ".new"
	os.WriteFile(other, []byte("replacement binary"), 0700)
	if err := os.Rename(other, path); err != nil {
		t.Fatal(err)
	}
	d := &Daemon{workspace: &workspaceTools{}}
	w := httptest.NewRecorder()
	d.handleHealth(w, httptest.NewRequest("GET", "/v1/health", nil))
	var h HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &h); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(original)
	want := "sha256:" + hex.EncodeToString(sum[:])
	if h.Build.ExecutableSHA256 != want {
		t.Fatalf("build %s, want original %s", h.Build.ExecutableSHA256, want)
	}
}

type homeResolver struct{ expires int64 }

func (r homeResolver) ResolveContext(context.Context, string, int64) (authority.ResolveOutcome, error) {
	return authority.ResolveOutcome{Success: &authority.SessionBinding{EpisodeID: "working-episode", ExpiresAt: r.expires}}, nil
}
func TestHomeHealthActivity(t *testing.T) {
	t.Setenv("XPC_SERVICE_NAME", serviceLabel)
	now := time.Now().Unix()
	d := &Daemon{startedAt: now, workspace: &workspaceTools{}}
	expired := d.activity.begin("expired", now-1, now)
	expired()
	old := d.activity.begin("idle", now+1000, now-301)
	old()
	resolver := activityResolver{inner: homeResolver{now + 1000}, activity: &d.activity}
	handler := trackActivityRequests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := resolver.ResolveContext(r.Context(), "Bearer secret-token", now); err != nil {
			t.Fatal(err)
		}
		h := httptest.NewRecorder()
		d.handleHealth(h, r)
		var health HealthResponse
		if err := json.Unmarshal(h.Body.Bytes(), &health); err != nil {
			t.Fatal(err)
		}
		if health.StartedAt != now || health.Supervisor != serviceLabel || health.Active.InFlight != 1 || len(health.Active.Sessions) != 1 || health.Active.Sessions[0].Episode != "working-episode" || health.Active.Sessions[0].LastCallUnix != now {
			t.Fatalf("%+v", health)
		}
		for _, secret := range []string{"secret-token", "credential_id", "session_id", "token"} {
			if strings.Contains(h.Body.String(), secret) {
				t.Fatal(h.Body.String())
			}
		}
	}))
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/mcp/", nil))
	if d.activity.snapshot(now).InFlight != 0 {
		t.Fatal("flight did not drain")
	}
	t.Setenv("XPC_SERVICE_NAME", "dev.ailang.other")
	if serviceSupervisor() != "none" {
		t.Fatal(serviceSupervisor())
	}
}
