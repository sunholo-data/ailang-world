package surfacefollow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/surface"
)

const StatusSchema = "world/surface-follow-status/v1"

// RepoStatus contains only browser-safe source facts.
type RepoStatus struct {
	ID           string                 `json:"id"`
	State        string                 `json:"state"`
	Head         string                 `json:"head"`
	AcceptedHead string                 `json:"accepted_head"`
	BuiltAt      string                 `json:"built_at"`
	CheckedAt    string                 `json:"checked_at"`
	LastFailure  *surface.FollowFailure `json:"last_failure,omitempty"`
	RetryAt      string                 `json:"retry_at,omitempty"`
}

type StatusDocument struct {
	Schema     string       `json:"schema"`
	Mode       string       `json:"mode"`
	Generation string       `json:"generation"`
	Repos      []RepoStatus `json:"repos"`
}

// Snapshot and everything reachable from it are immutable after publication.
// Readers capture this pointer once for the entire request.
type Snapshot struct {
	Bytes    []byte
	Envelope surface.FollowIndex
	Status   StatusDocument
}

type cacheFile interface {
	io.Writer
	Sync() error
	Close() error
	Name() string
}

type cacheOps struct {
	create  func(string, string) (cacheFile, error)
	rename  func(string, string) error
	syncDir func(string) error
}

func diskOps() cacheOps {
	return cacheOps{
		create: func(dir, pattern string) (cacheFile, error) { return os.CreateTemp(dir, pattern) },
		rename: os.Rename,
		syncDir: func(dir string) error {
			f, err := os.Open(dir)
			if err != nil {
				return err
			}
			defer f.Close()
			return f.Sync()
		},
	}
}

// Publisher's Sink only records the latest status and signals bounded work.
// Run is the sole writer; no cache work happens on the follower's event loop.
type Publisher struct {
	config  Config
	path    string
	log     io.Writer
	ops     cacheOps
	current atomic.Pointer[Snapshot]
	mu      sync.Mutex
	latest  map[string]Status
	wake    chan struct{}
	cleaned bool
}

func cacheRefusal(reason string) error {
	return fmt.Errorf("surface follow cache refused: %s; next: ailang-worldd serve --help", reason)
}

// NewPublisher establishes cache authority and recovers only this config's cache.
// It starts no goroutine. configBytes must be the bytes of the loaded config.
func NewPublisher(cfg Config, configBytes []byte, stateDir string, forbiddenRoots []string, log io.Writer) (*Publisher, error) {
	dir := filepath.Join(stateDir, "surface-follow")
	for _, root := range forbiddenRoots {
		if broker.CheckPolicyOutsideRoot(dir, root) != nil {
			return nil, cacheRefusal("directory must be outside workspace roots")
		}
	}
	for _, r := range cfg.Repos {
		if broker.CheckPolicyOutsideRoot(dir, r.Repo) != nil {
			return nil, cacheRefusal("directory must be outside source repos")
		}
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, cacheRefusal("cannot create directory")
	}
	info, err := os.Lstat(dir)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, cacheRefusal("directory must be a non-symlink directory")
	}
	if err := os.Chmod(dir, 0700); err != nil {
		return nil, cacheRefusal("cannot secure directory")
	}
	key := sha256.Sum256(configBytes)
	if log == nil {
		log = io.Discard
	}
	p := &Publisher{config: cfg, path: filepath.Join(dir, "index-"+hex.EncodeToString(key[:8])+".json"),
		log: log, ops: diskOps(), latest: map[string]Status{}, wake: make(chan struct{}, 1)}
	for _, r := range cfg.Repos {
		p.latest[r.ID] = Status{ID: r.ID}
	}
	p.recover()
	initial := p.makeSnapshot(p.latest)
	_ = encodeSnapshot(initial) // recovered partitions were strictly decoded
	p.current.Store(initial)
	return p, nil
}

func publishFailure(class string) *surface.FollowFailure {
	text := map[string]string{
		"restart":     "not checked since restart",
		"cache-write": "Cache publication failed; next: ailang-worldd service status",
		"candidate":   "Candidate refused; lower a repo limit; next: ailang-worldd service status",
		"oversize":    "Candidate exceeds byte bound; lower a repo limit; next: ailang-worldd service status",
	}[class]
	return &surface.FollowFailure{Class: class, Text: text}
}

func (p *Publisher) recover() {
	info, err := os.Lstat(p.path)
	if os.IsNotExist(err) {
		return
	}
	if err != nil || !info.Mode().IsRegular() {
		fmt.Fprintln(p.log, "surface follow cache: cache-read")
		return
	}
	f, err := os.Open(p.path)
	if err != nil {
		fmt.Fprintln(p.log, "surface follow cache: cache-read")
		return
	}
	data, err := io.ReadAll(io.LimitReader(f, surface.MaxBytes+1))
	_ = f.Close()
	if err != nil {
		fmt.Fprintln(p.log, "surface follow cache: cache-read")
		return
	}
	ix, err := surface.DecodeFollow(data)
	if err != nil {
		fmt.Fprintln(p.log, "surface follow cache: candidate")
		return
	}
	for _, r := range p.config.Repos {
		for _, part := range ix.Repos {
			if part.ID == r.ID && part.Index.Source.Rev == r.Branch {
				part.LastFailure = publishFailure("restart")
				p.latest[r.ID] = Status{ID: r.ID, Head: part.Head, BuiltAt: part.BuiltAt,
					CheckedAt: part.CheckedAt, LastFailure: part.LastFailure, Partition: &part}
			}
		}
	}
}

func (p *Publisher) Seed() map[string]surface.FollowPartition {
	out := map[string]surface.FollowPartition{}
	for _, part := range p.Snapshot().Envelope.Repos {
		out[part.ID] = part
	}
	return out
}

func (p *Publisher) Snapshot() *Snapshot { return p.current.Load() }

func (p *Publisher) Sink(s Status) {
	p.mu.Lock()
	if _, ok := p.latest[s.ID]; ok {
		p.latest[s.ID] = s
	}
	p.mu.Unlock()
	select {
	case p.wake <- struct{}{}:
	default:
	}
}

func (p *Publisher) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.wake:
			if ctx.Err() != nil {
				return
			}
			p.publish()
		}
	}
}

func repoStatus(s Status) RepoStatus {
	r := RepoStatus{ID: s.ID, Head: s.Head, CheckedAt: s.CheckedAt, LastFailure: s.LastFailure}
	if !s.RetryAt.IsZero() {
		r.RetryAt = s.RetryAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	if s.Partition == nil {
		r.State = "building"
		if s.LastFailure != nil {
			r.State = "refused"
		}
		return r
	}
	r.AcceptedHead, r.BuiltAt = s.Partition.Head, s.Partition.BuiltAt
	switch {
	case s.LastFailure != nil:
		r.State = "stale"
	case s.Head != s.Partition.Head:
		r.State = "updating"
	default:
		r.State = "fresh"
	}
	return r
}

func (p *Publisher) makeSnapshot(states map[string]Status) *Snapshot {
	s := &Snapshot{Envelope: surface.FollowIndex{Schema: surface.FollowSchema},
		Status: StatusDocument{Schema: StatusSchema, Mode: "follow", Repos: []RepoStatus{}}}
	keys := make([][3]string, 0, len(p.config.Repos))
	for _, r := range p.config.Repos {
		st := states[r.ID]
		s.Status.Repos = append(s.Status.Repos, repoStatus(st))
		if st.Partition != nil {
			part := *st.Partition
			part.CheckedAt, part.LastFailure = st.CheckedAt, st.LastFailure
			s.Envelope.Repos = append(s.Envelope.Repos, part)
			keys = append(keys, [3]string{part.ID, part.Head, part.BuiltAt})
		}
	}
	b, _ := json.Marshal(keys)
	h := sha256.Sum256(b)
	s.Envelope.Generation = "g" + hex.EncodeToString(h[:8])
	s.Status.Generation = s.Envelope.Generation
	return s
}

// encodeSnapshot validates before serialization and classifies publication
// failures without exposing decoder details or source paths.
func encodeSnapshot(s *Snapshot) string {
	if len(s.Envelope.Repos) > 0 && surface.ValidateFollow(s.Envelope) != nil {
		return "candidate"
	}
	data, err := json.Marshal(s.Envelope)
	if err != nil {
		return "candidate"
	}
	if len(data) > surface.MaxBytes {
		return "oversize"
	}
	s.Bytes = data
	return ""
}

func (p *Publisher) publish() {
	p.mu.Lock()
	states := make(map[string]Status, len(p.latest))
	for id, s := range p.latest {
		states[id] = s
	}
	p.mu.Unlock()
	candidate := p.makeSnapshot(states)
	old := p.Snapshot()
	class := encodeSnapshot(candidate)
	if class == "" && len(candidate.Envelope.Repos) > 0 && candidate.Status.Generation != old.Status.Generation {
		if p.write(candidate.Bytes) != nil {
			class = "cache-write"
		}
	}
	if class != "" {
		// Roll back data only. Observations still describe the latest probe.
		accepted := map[string]*surface.FollowPartition{}
		for i := range old.Envelope.Repos {
			part := old.Envelope.Repos[i]
			accepted[part.ID] = &part
		}
		for id, st := range states {
			st.Partition = accepted[id]
			if st.Partition != nil {
				st.CheckedAt = st.Partition.CheckedAt
			}
			st.LastFailure = publishFailure(class)
			states[id] = st
		}
		candidate = p.makeSnapshot(states)
		_ = encodeSnapshot(candidate) // last-good data plus fixed failure text
	}
	p.current.Store(candidate)
}

func (p *Publisher) write(data []byte) error {
	dir := filepath.Dir(p.path)
	f, err := p.ops.create(dir, ".index-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	n, err := f.Write(data)
	if err == nil && n != len(data) {
		err = io.ErrShortWrite
	}
	if err != nil {
		_ = f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	// Preserve the previous inode until directory durability is established,
	// allowing a failed post-rename sync to restore the previous cache bytes.
	backup := name + ".previous"
	hadOld := false
	if err = os.Link(p.path, backup); err == nil {
		hadOld = true
	} else if !os.IsNotExist(err) {
		return err
	}
	defer os.Remove(backup)
	if err = p.ops.rename(name, p.path); err != nil {
		return err
	}
	if err = p.ops.syncDir(dir); err != nil {
		if hadOld {
			_ = os.Rename(backup, p.path)
		} else {
			_ = os.Remove(p.path)
		}
		_ = p.ops.syncDir(dir)
		return err
	}
	if !p.cleaned {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintln(p.log, "surface follow cache cleanup: cache-read")
		}
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "index-") && strings.HasSuffix(e.Name(), ".json") && e.Name() != filepath.Base(p.path) {
				if err := os.Remove(filepath.Join(dir, e.Name())); err != nil {
					fmt.Fprintln(p.log, "surface follow cache cleanup: cache-write")
				}
			}
		}
		p.cleaned = true
	}
	return nil
}
