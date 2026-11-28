package main

import (
	"context"
	"encoding/json"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Read the configured presets on runner ticks, independent of HTTP input.
func (e serviceEnv) readOperatorCatalogue(c runnerConfig, src operatorStartSource) []byte {
	items := make(map[string]any)
	r, err := readServiceRecord(e.paths)
	if err != nil {
		return []byte("{}")
	}
	h, err := e.health(c.Address)
	if err != nil || h.Status != "ok" || h.Build.ExecutableSHA256 != c.Executables[e.paths.Binary] {
		return []byte("{}")
	}
	ctx, cancel := e.operatorContext(3 * time.Second)
	defer cancel()
	for _, name := range c.PresetNames {
		p, err := loadLaunchPreset(ctx, name, src.presets, src.profiles, r.Workspace)
		if err == nil {
			if entry, ok := operatorCatalogueEntry(ctx, c, r.Workspace, p, h); ok {
				items[name] = entry
			}
		}
	}
	if len(c.FinishRemotes) != 0 {
		src.workspace = r.Workspace
		finishCtx, stop := e.operatorContext(4 * time.Second)
		items[store.OperatorFinishCatalogueKey] = e.finishCatalogue(finishCtx, c, src)
		stop()
	}
	raw, _ := json.Marshal(items)
	if len(raw) > 16*approval.MaxOperatorActionBytes {
		delete(items, store.OperatorFinishCatalogueKey)
		raw, _ = json.Marshal(items)
	}
	if len(raw) > 16*approval.MaxOperatorActionBytes {
		return []byte("{}")
	}
	return raw
}

func operatorCatalogueEntry(ctx context.Context, c runnerConfig, root string, p *resolvedPreset, h daemon.HealthResponse) (store.OperatorStartReview, bool) {
	world, worldErr := store.SameDatabase(c.World, p.Daemon.DB)
	creds, credErr := credentialsFor(p.Daemon.DB, p.Daemon.CredentialsDB)
	workspace, workspaceErr := filepath.EvalSymlinks(p.Daemon.WorkspaceRoot)
	physical, physicalErr := filepath.EvalSymlinks(root)
	if ctx.Err() != nil || !world || worldErr != nil || credErr != nil || creds != c.Credentials || p.Daemon.Addr != c.Address ||
		workspaceErr != nil || physicalErr != nil || physical != workspace || p.Preset.Agent != "codex" || p.ProfileDigest == "" {
		return store.OperatorStartReview{}, false
	}
	// Git remote configuration can contain credentials. Never publish these
	// forms to a browser, even when accepted at the terminal.
	if strings.Contains(p.RemoteURL, "://") {
		u, err := url.Parse(p.RemoteURL)
		if err != nil || u.RawQuery != "" || u.Fragment != "" || u.User != nil && (u.Scheme != "ssh" || u.User.String() != "git") {
			return store.OperatorStartReview{}, false
		}
	}
	_, wall, err := launchMode(p.Preset, true, false)
	binary, binErr := launchBinary(p.Preset.Agent)
	sessions, pathErr := operatorPhysicalTarget(filepath.Join(p.Daemon.SessionsDir, "snapshot.session"))
	sum := c.Executables[binary]
	if err != nil || binErr != nil || pathErr != nil || sum == "" || serviceFileHash(binary) != sum {
		return store.OperatorStartReview{}, false
	}
	profile := approval.ProfileRef{Project: p.Preset.Profile, Digest: p.ProfileDigest}
	req := &mintRequest{credPath: creds, prof: &broker.ExecProfile{Project: profile.Project, Digest: profile.Digest}}
	if checkPresetDaemon(h, req, p, c.Address) != nil {
		return store.OperatorStartReview{}, false
	}
	out := store.OperatorStartReview{Start: *p.startRecord("", "", p.Preset.Agent, p.Preset.Model, p.Preset.TTLS, false, []string{}, nil),
		Profile: profile, Store: c.Credentials, Workspace: p.Daemon.WorkspaceRoot, Sessions: filepath.Dir(sessions),
		Address: c.Address, Interpreter: h.ToolRef, Binary: binary, BinaryDigest: sum, Wall: wall.String(), Committer: displayCommitter(p.Preset.Committer)}
	build, _ := json.Marshal(h.Build)
	out.Build = string(build)
	for _, g := range p.Grants {
		out.Grants = append(out.Grants, approval.Grant{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget})
	}
	return out, true
}
