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
	path   string
	source string // absolute spelling used to match libnode digests by path
	info   os.FileInfo
	libs   []execToolIdentity
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
	id := execToolIdentity{path: real, source: abs, info: info}
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
	revision                           uint64 // status publication order, independent of wall-clock time
}

type ExecToolStatus struct {
	Configured string `json:"configured"`
	Target     string `json:"target"`
	Accepted   string `json:"accepted"`
	Probe      string `json:"probe"`
	CheckedAt  string `json:"checked_at"`
	Next       string `json:"next,omitempty"`
}

// ExecToolState publishes immutable sandbox/profile snapshots. A shared flight
// serializes re-checks; mu protects publication, never hashing or a probe.
type ExecToolState struct {
	flight    *execToolFlight
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
	s := &ExecToolState{tools: map[string]execAcceptedTool{},
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
	t.revision++
	if probe {
		t.failed = &id
	} else {
		t.failed = nil
	}
	s.tools[key] = t
}

type execToolFlight struct {
	done                  chan struct{}
	project, tool, target string
	sandbox               *ExecSandbox
	profile               *ExecProfile
	why                   string
	err                   error
}

// execRecheckBudget is the sum of the profiles' arm timeouts plus ten
// seconds for setup, capped at 120 seconds. It belongs to the shared flight,
// never to the first caller's deadline.
func execRecheckBudget(profiles map[string]*ExecProfile) time.Duration {
	budget := 10 * time.Second
	for _, p := range profiles {
		budget += time.Duration(len(execProbeOrder(p))-1)*execProbeRunBudget +
			time.Duration(p.TimeoutMS)*time.Millisecond
	}
	if budget > 120*time.Second {
		budget = 120 * time.Second
	}
	return budget
}

func (s *ExecToolState) resolve(ctx context.Context, project, repo string) (*ExecSandbox, *ExecProfile, string, error) {
	for {
		s.mu.RLock()
		sb, profile := s.sandbox, s.profiles[project]
		tools := make(map[string]execAcceptedTool, len(s.tools))
		for key, t := range s.tools {
			tools[key] = t
		}
		flight := s.flight
		s.mu.RUnlock()
		if profile == nil {
			return nil, nil, fmt.Sprintf("unknown exec project %q", project), nil
		}
		// The common path does all filesystem work outside the publication
		// gate. Revalidate the snapshot under mu before returning it.
		pendingTool, pendingTarget := "node", tools["node"].configured
		clean := flight == nil
		if clean {
			if err := sb.verifyArchive(); err != nil {
				return nil, nil, "", err
			}
			for key, t := range tools {
				if key != "node" && t.project != project {
					continue
				}
				id, err := toolIdentity(t.configured, key == "node")
				if err != nil || !sameTool(id, t.identity) || t.why != "" {
					clean = false
					pendingTool, pendingTarget = t.name, t.configured
					break
				}
			}
			if clean {
				if err := sb.Verify(); err != nil {
					return nil, nil, "", err
				}
				for key, t := range tools {
					if key != "node" && t.project == project {
						if err := verifyFileDigest(t.identity.path, t.sum); err != nil {
							return nil, nil, "", err
						}
					}
				}
			}
		}
		s.mu.Lock()
		if s.sandbox != sb {
			s.mu.Unlock()
			continue
		}
		for key, t := range tools {
			if s.tools[key].revision != t.revision {
				clean = false
			}
		}
		if clean && s.flight == nil {
			s.mu.Unlock()
			return sb, profile, "", nil
		}
		flight = s.flight
		if flight == nil {
			flight = &execToolFlight{done: make(chan struct{}), project: project, tool: pendingTool, target: pendingTarget}
			s.flight = flight
			// Bound the detached root as well as the probe below: at most a
			// node version check plus the 120-second probe ceiling.
			workCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), execNodeVersionBudget+120*time.Second)
			go func(f *execToolFlight) {
				defer cancel()
				f.sandbox, f.profile, f.why, f.err = s.resolveOnce(workCtx, project, repo)
				s.mu.Lock()
				s.flight = nil
				close(f.done)
				s.mu.Unlock()
			}(flight)
		}
		s.mu.Unlock()
		select {
		case <-ctx.Done():
			return nil, nil, fmt.Sprintf("exec %s: re-checking the changed target\n  why: the startup probe is still running after %s changed\n  next: retry in a few seconds (no restart needed)", flight.tool, flight.target), nil
		case <-flight.done:
			if flight.project == project {
				return flight.sandbox, flight.profile, flight.why, flight.err
			}
		}
	}
}

// resolve returns a snapshot for this call. Waiting calls obey their own
// contexts. Failed probes are cached by identity; start failures are retried.
func (s *ExecToolState) resolveOnce(ctx context.Context, project, repo string) (*ExecSandbox, *ExecProfile, string, error) {
	if repo == "" {
		repo = s.cfg.ProbeConfig.Paths.WorkspaceRoot
	}
	interrupted := func(tool string) (*ExecSandbox, *ExecProfile, string, error) {
		return nil, nil, "the re-check after " + tool + " moved was interrupted; retry", nil
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
	excluded := map[string]string{}
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
			refused := execToolF15(t, target, err.Error(), why, next)
			s.remember(key, id, "does-not-start: "+why, next, false)
			if key != "node" && nodeChanged {
				excluded[t.project] = refused
				delete(current, key)
				continue
			}
			return nil, nil, refused, nil
		}
		if t.failed != nil && sameTool(*t.failed, id) {
			if key != "node" && nodeChanged {
				excluded[t.project] = t.refused
				delete(current, key)
				continue
			}
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
		if key != "node" && excluded[t.project] != "" {
			delete(current, key)
			continue
		}
		// An unchanged node (including a start-failure retry) is checked
		// by sb.Verify against its retained digests below. This loop protects
		// unchanged components when accepting a different node identity.
		if key == "node" && !changed[key] {
			continue
		}
		if t.identity.path == id.path && os.SameFile(t.identity.info, id.info) {
			if err := verifyFileDigest(id.path, t.sum); err != nil {
				return nil, nil, "", err
			}
		}
		if key == "node" {
			for _, oldLib := range t.identity.libs {
				for _, lib := range id.libs {
					if sameTool(oldLib, lib) {
						for _, digest := range oldSB.Libnode {
							if digest.Path == oldLib.source {
								if err := verifyFileDigest(lib.path, digest.SHA256); err != nil {
									return nil, nil, "", err
								}
							}
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
		if key != "node" && excluded[t.project] != "" {
			continue
		}
		if !changed[key] && !(nodeChanged && key != "node") && !(key == "node" && strings.HasPrefix(t.why, "does-not-start: ")) {
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
				err = checkExecChainPlacement("argv[0]", t.configured, t.configured, s.forbidden)
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
			refused := execToolF15(t, id.path, err.Error(), why, next)
			s.remember(key, id, "does-not-start: "+why+": "+err.Error(), next, false)
			if key != "node" && nodeChanged {
				excluded[t.project] = refused
				delete(current, key)
				continue
			}
			return nil, nil, refused, nil
		}
		// An unchanged identity retains its accepted digest even after a
		// start failure. Verify below refuses rewritten bytes without a probe.
		if changed[key] {
			// A libnode change alone does not change the node file identity.
			// Preserve that file's digest while updating the composite identity.
			sum := t.sum
			var err error
			if t.identity.path != id.path || !os.SameFile(t.identity.info, id.info) {
				sum, err = sha256File(id.path)
			}
			if err != nil {
				return nil, nil, "", err
			}
			t.identity, t.sum = id, sum
		}
		tools[key] = t
	}
	// Excluding a profile must not publish any of its unprobed candidates.
	// Its original accepted tools remain available for a later repair/retry.
	s.mu.RLock()
	for key, t := range tools {
		if key != "node" && excluded[t.project] != "" {
			tools[key] = s.tools[key]
			delete(current, key)
		}
	}
	s.mu.RUnlock()
	if nodeChanged {
		moved, movedKey = tools["node"].name, "node"
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
			if excluded[p] == "" && (nodeChanged || affected[p]) {
				order = append(order, profiles[p])
			}
		}
		cfg := s.cfg.ProbeConfig
		cfg.Sandbox = sb
		probed := make(map[string]*ExecProfile, len(order))
		for _, p := range order {
			probed[p.Project] = p
		}
		probeCtx, cancel := context.WithTimeout(ctx, execRecheckBudget(probed))
		errs, err := s.cfg.Probe(probeCtx, cfg, order)
		if probeCtx.Err() != nil {
			err = errors.Join(err, probeCtx.Err())
		}
		cancel()
		if ctx.Err() != nil {
			return interrupted(moved)
		}
		// Like startup's applyProbe, accept passing profiles and exclude a
		// failing profile with its own tool's refusal. A global probe error
		// (or no passing profile) still refuses the changed stack.
		passing := 0
		for i := range order {
			if i < len(errs) && errs[i] == nil {
				passing++
			}
		}
		if err == nil && passing > 0 {
			for i, p := range order {
				if i >= len(errs) || errs[i] == nil {
					continue
				}
				s.mu.RLock()
				profiles[p.Project] = s.profiles[p.Project]
				s.mu.RUnlock()
				for _, key := range sortedKeys(current) {
					t := tools[key]
					if key == "node" || t.project != p.Project {
						continue
					}
					why := "the startup probe failed after the target changed: " + errs[i].Error()
					next := execToolNext(t.configured, t.identity.path, repo, t.project, false)
					refused := execToolF15(t, t.identity.path, errs[i].Error(), why, next)
					excluded[p.Project] = refused
					s.remember(key, current[key], "failed: "+why, next, true)
					s.mu.Lock()
					old := s.tools[key]
					old.refused = refused
					s.tools[key] = old
					s.mu.Unlock()
					delete(current, key)
				}
			}
		} else {
			err = errors.Join(append(errs, err)...)
		}
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
					old.revision++
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
		// A failure recorded after our snapshot wins over this late call.
		latest := s.tools[key]
		if latest.revision != tools[key].revision && sameTool(latest.identity, t.identity) {
			t.failed, t.observed, t.why, t.next, t.refused = latest.failed, latest.observed, latest.why, latest.next, latest.refused
		}
		t.checked = time.Now().UTC().Format(time.RFC3339Nano)
		t.revision = latest.revision + 1
		s.tools[key] = t
	}
	s.sandbox, s.profiles = sb, profiles
	s.mu.Unlock()
	return sb, profiles[project], excluded[project], nil
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
		t.revision++
		s.tools["node"] = t
	}
	return execToolF15(t, sb.Node, why, "--version failed", next)
}
