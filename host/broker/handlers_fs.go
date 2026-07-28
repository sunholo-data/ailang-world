package broker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	EffectFSRead  = "FS.Read"
	EffectFSWrite = "FS.Write"
)

// FSPathError is a defense-in-depth handler refusal after authority has
// already been decided.
type FSPathError struct {
	Scope    string
	Resolved string
	Why      string
}

func (e *FSPathError) Error() string {
	return fmt.Sprintf("broker: FS path %q resolved to %q: %s", e.Scope, e.Resolved, e.Why)
}

// FSHandler performs direct file I/O. The exact canonical absolute path is the
// request scope; payload is the complete content for FS.Write.
type FSHandler struct{}

func (FSHandler) Execute(_ context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	path, err := guardedFSPath(req.Scope, req.Effect == EffectFSWrite)
	if err != nil {
		return nil, err
	}
	switch req.Effect {
	case EffectFSRead:
		result, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("broker: FS.Read %q: %w", path, err)
		}
		return result, nil
	case EffectFSWrite:
		if err := os.WriteFile(path, payload, 0o600); err != nil {
			return nil, fmt.Errorf("broker: FS.Write %q: %w", path, err)
		}
		return []byte{}, nil
	default:
		return nil, fmt.Errorf("broker: FS handler does not implement %q", req.Effect)
	}
}

func guardedFSPath(scope string, allowMissing bool) (string, error) {
	if !filepath.IsAbs(scope) || filepath.Clean(scope) != scope {
		return "", &FSPathError{Scope: scope, Why: "scope is not a canonical absolute path"}
	}
	parent := filepath.Dir(scope)
	resolvedParent, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return "", &FSPathError{Scope: scope, Why: "resolve parent: " + err.Error()}
	}
	resolved, err := filepath.EvalSymlinks(scope)
	if err != nil {
		if !allowMissing || !os.IsNotExist(err) {
			return "", &FSPathError{Scope: scope, Why: "resolve path: " + err.Error()}
		}
		resolved = filepath.Join(resolvedParent, filepath.Base(scope))
	}
	rel, err := filepath.Rel(resolvedParent, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", &FSPathError{Scope: scope, Resolved: resolved, Why: "resolution leaves scope parent"}
	}
	return resolved, nil
}

// ProbeHandler is a deterministic echo surface for pipeline and replay tests.
type ProbeHandler struct {
	Dispatches *int
}

func (h ProbeHandler) Execute(_ context.Context, _ EffectRequest, payload []byte) ([]byte, error) {
	if h.Dispatches != nil {
		*h.Dispatches++
	}
	return append([]byte("probe:"), payload...), nil
}
