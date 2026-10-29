package broker

// Row 219 M3: accepted tool identities are shared across episodes. This is
// host code because it checks executable placement and runs confinement probes.
import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type ExecToolProbe func(context.Context, ExecProbeConfig, []*ExecProfile) ([]error, error)

type ExecToolStateConfig struct {
	Node        string // configured path, before EvalSymlinks
	Sandbox     *ExecSandbox
	Profiles    map[string]*ExecProfile
	ProbeConfig ExecProbeConfig
	Probe       ExecToolProbe // nil means all arms of ProbeExecProfilesEach
}

type execToolIdentity struct {
	path string
	info os.FileInfo
	libs []execToolIdentity
}

func toolIdentity(path string, node bool) (execToolIdentity, error) {
	configured := path
	if node {
		var err error
		path, err = execNodePath(path)
		if err != nil {
			return execToolIdentity{path: configured}, err
		}
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return execToolIdentity{path: configured}, err
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return execToolIdentity{path: path}, err
	}
	info, err := os.Stat(real)
	id := execToolIdentity{path: real, info: info}
	if err != nil {
		return id, err
	}
	if node {
		matches, err := filepath.Glob(filepath.Join(filepath.Dir(filepath.Dir(real)), "lib", "libnode.*"))
		if err != nil {
			return id, err
		}
		sort.Strings(matches)
		for _, p := range matches {
			lib, err := toolIdentity(p, false)
			if err != nil {
				return id, err
			}
			id.libs = append(id.libs, lib)
		}
	}
	return id, nil
}

func sameTool(a, b execToolIdentity) bool {
	if a.path != b.path || a.info == nil || b.info == nil || !os.SameFile(a.info, b.info) || len(a.libs) != len(b.libs) {
		return false
	}
	for i := range a.libs {
		if !sameTool(a.libs[i], b.libs[i]) {
			return false
		}
	}
	return true
}

type execAcceptedTool struct {
	configured, name, project, command string
	identity                           execToolIdentity
	sum                                string
	failed                             *execToolIdentity
	observed                           *execToolIdentity
	refused                            string
	why, next, checked                 string
}

type ExecToolStatus struct {
	Configured string `json:"configured"`
	Target     string `json:"target"`
	Accepted   string `json:"accepted"`
	Probe      string `json:"probe"`
	CheckedAt  string `json:"checked_at"`
	Next       string `json:"next,omitempty"`
}

// ExecToolState publishes immutable sandbox/profile snapshots. gate serializes
// re-checks while mu protects only publication and health reads, never a probe.
type ExecToolState struct {
	gate      chan struct{}
	mu        sync.RWMutex
	tools     map[string]execAcceptedTool
	sandbox   *ExecSandbox
	profiles  map[string]*ExecProfile
	cfg       ExecToolStateConfig
	forbidden []namedRoot
}

func execCommandKey(project, command string) string { return project + "/" + command }

func NewExecToolState(cfg ExecToolStateConfig) (*ExecToolState, error) {
	if cfg.Sandbox == nil {
		return nil, errors.New("exec tool state needs a sandbox")
	}
	if cfg.Node == "" {
		cfg.Node = cfg.Sandbox.Node
	}
	s := &ExecToolState{gate: make(chan struct{}, 1), tools: map[string]execAcceptedTool{},
		sandbox: cfg.Sandbox, profiles: cfg.Profiles, cfg: cfg,
		forbidden: []namedRoot{{"an exec cache", filepath.Join(cfg.ProbeConfig.Paths.StateDir, "exec-cache")},
			{"the state dir", cfg.ProbeConfig.Paths.StateDir}, {"the workspace root", cfg.ProbeConfig.Paths.WorkspaceRoot}}}
	if s.cfg.Probe == nil {
		s.cfg.Probe = ProbeExecProfilesEach
	}
	add := func(key, name, configured, target, sum, project, command string, node bool) error {
		id, err := toolIdentity(target, node)
		if err != nil {
			return err
		}
		s.tools[key] = execAcceptedTool{configured: configured, name: name, project: project, command: command,
			identity: id, sum: sum, checked: time.Now().UTC().Format(time.RFC3339Nano)}
		return nil
	}
	if err := add("node", "node", cfg.Node, cfg.Sandbox.Node, cfg.Sandbox.NodeSHA256, "", "", true); err != nil {
		return nil, err
	}
	for project, p := range cfg.Profiles {
		for name, c := range p.Commands {
			if err := add(execCommandKey(project, name), c.Argv[0], c.Argv0, c.Argv0, c.Argv0SHA256, project, name, false); err != nil {
				return nil, err
			}
		}
		if p.Probe.Argv0 != "" {
			c := p.Probe
			if err := add(execCommandKey(project, "@probe"), c.Argv[0], c.Argv0, c.Argv0, c.Argv0SHA256, project, "@probe", false); err != nil {
				return nil, err
			}
		}
	}
	s.sandbox, s.profiles = s.snapshots(s.tools)
	return s, nil
}

func (s *ExecToolState) snapshots(tools map[string]execAcceptedTool) (*ExecSandbox, map[string]*ExecProfile) {
	sb := *s.sandbox
	sb.Node = tools["node"].identity.path
	sb.NodeConfigured = tools["node"].configured
	sb.NodeSHA256 = tools["node"].sum
	nodeIdentity := tools["node"].identity
	sb.nodeIdentity = &nodeIdentity
	profiles := make(map[string]*ExecProfile, len(s.profiles))
	for project, old := range s.profiles {
		p := *old
		p.Commands = make(map[string]ExecProfileCommand, len(old.Commands))
		patch := func(c ExecProfileCommand, key string) ExecProfileCommand {
			t := tools[key]
			c.ConfiguredArgv0, c.Argv0, c.Argv0SHA256 = t.configured, t.identity.path, t.sum
			return c
		}
		for name, c := range old.Commands {
			p.Commands[name] = patch(c, execCommandKey(project, name))
		}
		if old.Probe.Argv0 != "" {
			p.Probe = patch(old.Probe, execCommandKey(project, "@probe"))
		}
		profiles[project] = &p
	}
	return &sb, profiles
}

func execToolNext(configured, target, repo, project string, gone bool) string {
	if gone && strings.Contains(filepath.ToSlash(configured), "/Cellar/") {
		return fmt.Sprintf("ailang-worldd profile init --repo %s --project %s", repo, project)
	}
	parts := strings.Split(filepath.ToSlash(target), "/Cellar/")
	if len(parts) == 2 && strings.Split(parts[1], "/")[0] != "" {
		return "brew reinstall " + strings.Split(parts[1], "/")[0] + ", then retry (no restart needed)"
	}
	return "reinstall the package that provides " + configured + ", then retry (no restart needed)"
}

func execToolF15(t execAcceptedTool, target, failure, why, next string) string {
	return fmt.Sprintf("exec %s %s → %s does not start: %s\n  why: %s\n  next: %s",
		t.name, t.configured, target, firstExecFailure(failure), why, next)
}

func firstExecFailure(s string) string {
	return strings.Split(strings.TrimSpace(s), "\n")[0]
}

func (s *ExecToolState) remember(key string, id execToolIdentity, why, next string, probe bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := s.tools[key]
	t.why, t.next, t.checked = why, next, time.Now().UTC().Format(time.RFC3339Nano)
	t.observed = &id
	if probe {
		t.failed = &id
	} else {
		t.failed = nil
	}
	s.tools[key] = t
}

// resolve returns a snapshot for this call. Waiting calls obey their own
// contexts. Failed probes are cached by identity; start failures are retried.
func (s *ExecToolState) resolve(ctx context.Context, project, repo string) (*ExecSandbox, *ExecProfile, string, error) {
	if repo == "" {
		repo = s.cfg.ProbeConfig.Paths.WorkspaceRoot
	}
	interrupted := func(tool string) (*ExecSandbox, *ExecProfile, string, error) {
		return nil, nil, "the re-check after " + tool + " moved was interrupted; retry", nil
	}
	select {
	case s.gate <- struct{}{}:
		defer func() { <-s.gate }()
	case <-ctx.Done():
		return interrupted("exec tool")
	}
	if ctx.Err() != nil {
		return interrupted("exec tool")
	}
	s.mu.RLock()
	tools := make(map[string]execAcceptedTool, len(s.tools))
	for k, t := range s.tools {
		tools[k] = t
	}
	oldSB := s.sandbox
	s.mu.RUnlock()
	if err := oldSB.verifyArchive(); err != nil {
		return nil, nil, "", err
	}
	current := map[string]execToolIdentity{}
	changed := map[string]bool{}
	nodeChanged := false
	affected := map[string]bool{}
	for _, key := range append([]string{"node"}, sortedKeys(toolsWithoutNode(tools))...) {
		t := tools[key]
		if key != "node" && !nodeChanged && t.project != project {
			continue
		}
		id, err := toolIdentity(t.configured, key == "node")
		current[key] = id
		if err != nil {
			target := id.path
			if target == t.configured {
				target = t.identity.path
			}
			next := execToolNext(t.configured, target, repo, execToolProject(t, project), true)
			why := "the target is gone"
			s.remember(key, id, "does-not-start: "+why, next, false)
			return nil, nil, execToolF15(t, target, err.Error(), why, next), nil
		}
		if t.failed != nil && sameTool(*t.failed, id) {
			return nil, nil, t.refused, nil
		}
		changed[key] = !sameTool(t.identity, id)
		if key == "node" {
			nodeChanged = changed[key]
		}
		if changed[key] {
			affected[t.project] = true
		}
	}
	// Never turn an in-place rewrite into an upgrade, including unchanged
	// libnode files when some other member of the node identity changed.
	for key, id := range current {
		t := tools[key]
		if t.identity.path == id.path && os.SameFile(t.identity.info, id.info) {
			if err := verifyFileDigest(id.path, t.sum); err != nil {
				return nil, nil, "", err
			}
		}
		if key == "node" {
			for i, oldLib := range t.identity.libs {
				for _, lib := range id.libs {
					if sameTool(oldLib, lib) && i < len(oldSB.Libnode) {
						if err := verifyFileDigest(lib.path, oldSB.Libnode[i].SHA256); err != nil {
							return nil, nil, "", err
						}
					}
				}
			}
		}
	}
	var moved string
	var movedKey string
	for _, key := range sortedKeys(current) {
		t, id := tools[key], current[key]
		if !changed[key] && !(key == "node" && strings.HasPrefix(t.why, "does-not-start: ")) {
			continue
		}
		if changed[key] {
			moved, movedKey = t.name, key
		}
		var err error
		if key == "node" {
			var target string
			target, err = resolveExecNode(ctx, t.configured, s.forbidden)
			if err == nil && target != id.path {
				return interrupted(t.name)
			}
		} else {
			// Reuse the profile loader's executable and '=' rules on the
			// accepted target, then enforce startup's placement.
			_, _, err = resolveArgv0(id.path, nil)
			if err == nil {
				err = checkPlacement("argv[0]", t.configured, id.path, s.forbidden)
			}
		}
		if ctx.Err() != nil {
			return interrupted(t.name)
		}
		if err != nil {
			next := execToolNext(t.configured, id.path, repo, execToolProject(t, project), false)
			why := "--version failed"
			if key != "node" {
				why = "the executable failed startup checks"
			}
			s.remember(key, id, "does-not-start: "+why+": "+err.Error(), next, false)
			return nil, nil, execToolF15(t, id.path, err.Error(), why, next), nil
		}
		sum, err := sha256File(id.path)
		if err != nil {
			return nil, nil, "", err
		}
		t.identity, t.sum = id, sum
		tools[key] = t
	}
	sb, profiles := s.snapshots(tools)
	if nodeChanged {
		libs, err := libnodeDigests(sb.Node)
		if err != nil {
			return nil, nil, "", err
		}
		sb.Libnode = libs
	}
	// Verify the archive before any probe can spawn with it.
	if err := sb.Verify(); err != nil {
		return nil, nil, "", err
	}
	if moved != "" {
		var order []*ExecProfile
		for _, p := range sortedKeys(profiles) {
			if nodeChanged || affected[p] {
				order = append(order, profiles[p])
			}
		}
		cfg := s.cfg.ProbeConfig
		cfg.Sandbox = sb
		errs, err := s.cfg.Probe(ctx, cfg, order)
		if ctx.Err() != nil {
			return interrupted(moved)
		}
		err = errors.Join(append(errs, err)...)
		if err != nil {
			t := tools[movedKey]
			why := "the startup probe failed after the target changed: " + err.Error()
			next := execToolNext(t.configured, t.identity.path, repo, execToolProject(t, project), false)
			refused := execToolF15(t, t.identity.path, err.Error(), why, next)
			s.mu.Lock()
			for key := range current {
				if changed[key] {
					old := s.tools[key]
					id := current[key]
					old.failed, old.observed = &id, &id
					old.why, old.next, old.refused = "failed: "+why, next, refused
					old.checked = time.Now().UTC().Format(time.RFC3339Nano)
					s.tools[key] = old
				}
			}
			s.mu.Unlock()
			return nil, nil, refused, nil
		}
		// A target moving during the probe is never published as accepted.
		for key, before := range current {
			after, err := toolIdentity(tools[key].configured, key == "node")
			if err != nil || !sameTool(before, after) {
				return interrupted(tools[key].name)
			}
		}
		if err := sb.Verify(); err != nil {
			return nil, nil, "", err
		}
		for key := range current {
			if key != "node" {
				if err := verifyFileDigest(tools[key].identity.path, tools[key].sum); err != nil {
					return nil, nil, "", err
				}
			}
		}
	}
	if ctx.Err() != nil {
		return interrupted(moved)
	}
	s.mu.Lock()
	for key := range current {
		t := tools[key]
		t.failed, t.observed, t.why, t.next, t.refused = nil, nil, "", "", ""
		t.checked = time.Now().UTC().Format(time.RFC3339Nano)
		s.tools[key] = t
	}
	s.sandbox, s.profiles = sb, profiles
	s.mu.Unlock()
	return sb, profiles[project], "", nil
}

func execToolProject(t execAcceptedTool, project string) string {
	if t.project != "" {
		return t.project
	}
	return project
}

func toolsWithoutNode(in map[string]execAcceptedTool) map[string]execAcceptedTool {
	out := make(map[string]execAcceptedTool, len(in))
	for key, t := range in {
		if key != "node" {
			out[key] = t
		}
	}
	return out
}

// Accepted reports only published targets, without looking up configured paths.
func (s *ExecToolState) Accepted() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := map[string]string{}
	for _, key := range append(sortedKeys(toolsWithoutNode(s.tools)), "node") {
		t := s.tools[key]
		if t.command != "@probe" {
			out[t.name] = t.identity.path
		}
	}
	return out
}

// Health resolves fresh identities without probing or changing accepted state.
func (s *ExecToolState) Health() (map[string]string, map[string]ExecToolStatus) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	accepted := map[string]string{}
	status := map[string]ExecToolStatus{}
	for _, key := range append(sortedKeys(toolsWithoutNode(s.tools)), "node") {
		t := s.tools[key]
		if t.command == "@probe" {
			continue
		}
		id, err := toolIdentity(t.configured, key == "node")
		probe := "ok"
		if err != nil {
			probe = "does-not-start: the target is gone"
		} else if t.failed != nil && sameTool(*t.failed, id) {
			probe = t.why
		} else if t.observed != nil && sameTool(*t.observed, id) && strings.HasPrefix(t.why, "does-not-start: ") {
			probe = t.why
		} else if !sameTool(t.identity, id) {
			probe = "pending"
		} else if strings.HasPrefix(t.why, "does-not-start: ") {
			probe = t.why
		}
		next := t.next
		if err != nil {
			next = execToolNext(t.configured, t.identity.path, s.cfg.ProbeConfig.Paths.WorkspaceRoot, t.project, true)
		}
		accepted[t.name] = t.identity.path
		status[t.name] = ExecToolStatus{Configured: t.configured, Target: id.path, Accepted: t.identity.path,
			Probe: probe, CheckedAt: t.checked, Next: next}
	}
	return accepted, status
}

func (s *ExecToolState) nodeFailed(sb *ExecSandbox, why string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := s.tools["node"]
	next := execToolNext(t.configured, sb.Node, s.cfg.ProbeConfig.Paths.WorkspaceRoot, "", false)
	// A late result for an older snapshot must not poison the current node.
	if sb.nodeIdentity != nil && sameTool(t.identity, *sb.nodeIdentity) {
		t.why, t.next = "does-not-start: --version failed: "+why, next
		t.checked = time.Now().UTC().Format(time.RFC3339Nano)
		id := t.identity
		t.observed = &id
		s.tools["node"] = t
	}
	return execToolF15(t, sb.Node, why, "--version failed", next)
}
