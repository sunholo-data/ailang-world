package daemon

import (
	"context"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/sunholo-data/ailang-world/host/authority"
)

const serviceLabel = "com.sunholo.ailang-world.worldd"

func serviceSupervisor() string {
	if os.Getenv("XPC_SERVICE_NAME") == serviceLabel {
		return serviceLabel
	}
	return "none"
}

type HealthActive struct {
	InFlight int                   `json:"in_flight"`
	Sessions []HealthActiveSession `json:"sessions"`
}
type HealthActiveSession struct {
	Episode      string `json:"episode"`
	LastCallUnix int64  `json:"last_call_unix"`
}
type activityKey struct {
	episode string
	expires int64
}
type activityTracker struct {
	mu       sync.Mutex
	inFlight int
	episodes map[activityKey]int64
}

func (a *activityTracker) begin(episode string, expires, now int64) func() {
	a.mu.Lock()
	if a.episodes == nil {
		a.episodes = make(map[activityKey]int64)
	}
	for ep, last := range a.episodes {
		if ep.expires < now || last < now-300 {
			delete(a.episodes, ep)
		}
	}
	a.episodes[activityKey{episode, expires}] = now
	a.inFlight++
	a.mu.Unlock()
	return func() { a.mu.Lock(); a.inFlight--; a.mu.Unlock() }
}
func (a *activityTracker) snapshot(now int64) HealthActive {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := HealthActive{InFlight: a.inFlight, Sessions: []HealthActiveSession{}}
	recent := make(map[string]int64)
	for key, last := range a.episodes {
		if key.expires >= now && last >= now-300 && last > recent[key.episode] {
			recent[key.episode] = last
		}
	}
	for episode, last := range recent {
		out.Sessions = append(out.Sessions, HealthActiveSession{episode, last})
	}
	sort.Slice(out.Sessions, func(i, j int) bool { return out.Sessions[i].Episode < out.Sessions[j].Episode })
	return out
}

type activityRequestKey struct{}
type activityRequest struct {
	mu   sync.Mutex
	done func()
}

func trackActivityRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := &activityRequest{}
		defer func() {
			state.mu.Lock()
			defer state.mu.Unlock()
			if state.done != nil {
				state.done()
			}
		}()
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), activityRequestKey{}, state)))
	})
}

type activityResolver struct {
	inner    authority.Resolver
	activity *activityTracker
}

func (r activityResolver) ResolveContext(ctx context.Context, header string, now int64) (authority.ResolveOutcome, error) {
	out, err := r.inner.ResolveContext(ctx, header, now)
	if err == nil && out.Success != nil {
		if state, ok := ctx.Value(activityRequestKey{}).(*activityRequest); ok {
			state.mu.Lock()
			if state.done == nil {
				state.done = r.activity.begin(out.Success.EpisodeID, out.Success.ExpiresAt, now)
			}
			state.mu.Unlock()
		}
	}
	return out, err
}
