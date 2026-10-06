package daemon

// Row 141 (w-workspace-project-layouts): project layouts the one-directory-
// per-episode default does not fit.
//
// M1 — a MODULE ROOT: `--workspace-module-root REL` (and the per-episode
// `--workspace-episode-module-root EP=REL`) makes <root>/<episode>/REL the
// AILANG sandbox, because policy-tool runs every CLI child with cwd = the
// sandbox and AILANG resolves bare imports against the cwd. The grammar is
// row 140's exec-profile path grammar; resolution is projectDir's rule
// (symlink-free, an existing directory) and NOTHING is ever created.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// checkModuleRoot refuses a REL outside the shared exec path grammar, an
// unclean one, and one renderPolicy would refuse. Every error names the flag
// and the raw value.
func checkModuleRoot(flag, value, rel, root string) error {
	switch {
	case !broker.ExecPathOk(rel):
		return fmt.Errorf("%s %q: %q must be relative, with no leading / or -, and no .. segment", flag, value, rel)
	case filepath.Clean(rel) != rel:
		return fmt.Errorf("%s %q: %q is not clean", flag, value, rel)
	}
	if _, err := broker.RenderEpisodePolicy(filepath.Join(root, rel)); err != nil {
		return fmt.Errorf("%s %q: %v", flag, value, err)
	}
	return nil
}

// resolveModuleRoots validates the default and the per-episode overrides
// without touching the filesystem.
func resolveModuleRoots(def string, pairs []string, root string) (string, map[string]string, error) {
	if def != "" {
		if err := checkModuleRoot("--workspace-module-root", def, def, root); err != nil {
			return "", nil, err
		}
	}
	byEp := map[string]string{}
	for _, pair := range pairs {
		const flag = "--workspace-episode-module-root"
		ep, rel, ok := strings.Cut(pair, "=")
		switch {
		case !ok:
			return "", nil, fmt.Errorf("%s %q: is not EP=REL", flag, pair)
		case !episodeIDPattern.MatchString(ep):
			return "", nil, fmt.Errorf("%s %q: episode %q is outside the episode grammar %s", flag, pair, ep, episodeIDPattern)
		}
		if _, dup := byEp[ep]; dup {
			return "", nil, fmt.Errorf("%s %q: maps episode %q a second time", flag, pair, ep)
		}
		if err := checkModuleRoot(flag, pair, rel, root); err != nil {
			return "", nil, err
		}
		byEp[ep] = rel
	}
	return def, byEp, nil
}

// sandboxRoot is the episode's AILANG sandbox: <epRoot>/REL, REL being the
// episode's override, else the default, else ".". It must be a real directory
// reached through no symlink (projectDir's rule). Never MkdirAll.
func (w *workspaceTools) sandboxRoot(episodeID, epRoot string) (string, error) {
	rel := w.moduleRoot
	if o, ok := w.episodeModuleRoot[episodeID]; ok {
		rel = o
	}
	if rel == "" {
		rel = "."
	}
	want := filepath.Join(epRoot, rel)
	got, err := filepath.EvalSymlinks(want)
	if err != nil {
		return "", fmt.Errorf("module root %q: %v", rel, err)
	}
	if got != want {
		return "", fmt.Errorf("module root %q: resolves through a symlink to %q", rel, got)
	}
	if info, err := os.Stat(want); err != nil || !info.IsDir() {
		return "", fmt.Errorf("module root %q: not a directory", rel)
	}
	return want, nil
}
