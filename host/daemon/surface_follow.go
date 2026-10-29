package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sunholo-data/ailang-world/host/surfacefollow"
)

const groupedFollowReason = "the grouped multi-repo model arrives in row 220 M4; use /surface/index.json and /surface/follow-status.json"

type surfaceFollow struct {
	config        surfacefollow.Config
	publisher     *surfacefollow.Publisher
	collector     *surfacefollow.Collector
	mu            sync.Mutex
	cancel        context.CancelFunc
	followerDone  chan struct{}
	publisherDone chan struct{}
}

func newSurfaceFollow(cfg Config, root string) (*surfaceFollow, error) {
	if cfg.SurfaceFollow == "" {
		return nil, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("surface follow refused: cannot resolve home; next: ailang-worldd serve --help")
	}
	var roots []string
	if root != "" {
		roots = []string{root}
	}
	config, err := surfacefollow.LoadConfig(cfg.SurfaceFollow, surfacefollow.LoadOptions{Home: home, ForbiddenRoots: roots})
	if err != nil {
		return nil, err
	}
	path := cfg.SurfaceFollow
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("surface follow refused: cannot read config bytes; next: ailang-worldd serve --help")
	}
	stateDir, err := canonicalStateDir(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("surface follow refused: cannot resolve state directory; next: ailang-worldd serve --help")
	}
	p, err := surfacefollow.NewPublisher(config, data, stateDir, roots, resolveErrorLog(cfg.ErrorLog))
	if err != nil {
		return nil, err
	}
	return &surfaceFollow{config: config, publisher: p, collector: surfacefollow.NewCollector(nil, nil)}, nil
}

func (d *Daemon) startSurfaceFollow(ctx context.Context) {
	f := d.follow
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.cancel != nil {
		return
	}
	ctx, f.cancel = context.WithCancel(ctx)
	f.followerDone, f.publisherDone = make(chan struct{}), make(chan struct{})
	follower := &surfacefollow.Follower{Config: f.config, Collector: f.collector, Sink: f.publisher.Sink, Seed: f.publisher.Seed()}
	go func() {
		defer close(f.publisherDone)
		f.publisher.Run(ctx)
	}()
	go func() {
		defer close(f.followerDone)
		_ = follower.Run(ctx)
	}()
}

func (d *Daemon) stopSurfaceFollow() {
	f := d.follow
	if f == nil {
		return
	}
	f.mu.Lock()
	cancel, followerDone, publisherDone := f.cancel, f.followerDone, f.publisherDone
	f.mu.Unlock()
	if cancel == nil {
		return
	}
	cancel()
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	for _, done := range []chan struct{}{followerDone, publisherDone} {
		select {
		case <-done:
		case <-timer.C:
			fmt.Fprintln(d.errLog, "ailang-worldd: surface follow did not stop within 2s")
			return
		}
	}
}

func (d *Daemon) handleFollowIndex(w http.ResponseWriter, _ *http.Request) {
	s := d.follow.publisher.Snapshot()
	if len(s.Envelope.Repos) == 0 {
		var states []string
		for _, r := range s.Status.Repos {
			state := r.State
			if r.LastFailure != nil {
				state += ": " + r.LastFailure.Class
			}
			states = append(states, r.ID+": "+state)
		}
		writeAPIError(w, "SurfaceIndexRefused", "the follow index has no accepted partitions: "+strings.Join(states, "; "), http.StatusServiceUnavailable)
		return
	}
	_, _ = w.Write(s.Bytes)
}

func (d *Daemon) handleSurfaceFollowStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	status := surfacefollow.StatusDocument{Schema: surfacefollow.StatusSchema, Mode: "none", Repos: []surfacefollow.RepoStatus{}}
	if d.follow != nil {
		s := d.follow.publisher.Snapshot()
		status = s.Status
	} else if d.cfg.SurfaceIndex != "" {
		status.Mode = "index"
	}
	_ = json.NewEncoder(w).Encode(status)
}
