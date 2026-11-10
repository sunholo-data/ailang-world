package main

// Row 201 A0 (design_docs/planned/w-world-approvals-by-key-a0a1.md §2.1a, §2.5): the mint
// request `episode start --request` stores and `approve` mints from.

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

// approveReadBudget bounds the store reads of `episode start --request` and `approve` (3 s).
const approveReadBudget = 3 * time.Second

// buildMintRequest is design §2.5 step 3: the MintRequest of a validated mintRequest. Every value
// approve acts on is in it; paths are canonical (stores) or absolute (the rest).
func buildMintRequest(r *mintRequest, build daemon.BuildIdentity, addr, agent, model, task, mcpFile string) (approval.MintRequest, error) {
	for i, g := range r.grants {
		if g.ExpiresAt != 0 {
			return approval.MintRequest{}, fmt.Errorf("grant %d carries an expiry; a request's grants expire at the mint", i)
		}
	}
	worldDB, err := store.CanonicalPath(r.db)
	if err != nil {
		return approval.MintRequest{}, err
	}
	wsRoot, err := filepath.Abs(r.wsRoot)
	if err != nil {
		return approval.MintRequest{}, err
	}
	repo := ""
	if r.repo != "" {
		if repo, err = filepath.Abs(r.repo); err != nil {
			return approval.MintRequest{}, err
		}
	}
	profilesDir, err := filepath.Abs(r.profilesDir)
	if err != nil {
		return approval.MintRequest{}, err
	}
	prof := approval.ProfileRef{}
	if r.prof != nil {
		prof = approval.ProfileRef{Project: r.prof.Project, Digest: r.prof.Digest}
	}
	mr := approval.MintRequest{Schema: approval.RequestSchema, Episode: r.episode, WorldDB: worldDB,
		WorkspaceRoot: wsRoot, Worktree: r.dir, Reuse: r.reuse, Degraded: r.degraded, Repo: repo, Git: r.git,
		Branch: r.branch, Grants: fromCapabilities(r.grants), TTLSeconds: r.ttl, Profile: prof,
		ProfilesDir: profilesDir, CredentialsDB: r.credPath, SessionFile: r.outPath, MCPConfig: mcpFile,
		Agent: agent, Model: model, Task: task,
		Daemon: approval.DaemonRef{Addr: addr, VCSRevision: build.VCSRevision, VCSModified: build.VCSModified,
			ExecutableSHA256: build.ExecutableSHA256, GoVersion: build.GoVersion}}
	return mr, mr.Validate()
}

// fromRequest is the §2.1a map read backwards: the mintOptions approve validates, built from the
// decoded request only. No flag or default reaches it.
func fromRequest(mr approval.MintRequest) mintOptions {
	return mintOptions{verb: "approve", episode: mr.Episode, db: mr.WorldDB, credsDB: mr.CredentialsDB,
		wsRoot: mr.WorkspaceRoot, repo: mr.Repo, branch: mr.Branch, allowNonGit: false,
		grants: toCapabilities(mr.Grants), ttl: mr.TTLSeconds, outPath: mr.SessionFile, outDirCreate: true,
		profile: mr.Profile.Project, profilesDir: mr.ProfilesDir}
}

func fromCapabilities(cs []broker.Capability) []approval.Grant {
	gs := make([]approval.Grant, 0, len(cs))
	for _, c := range cs {
		gs = append(gs, approval.Grant{Effect: c.Effect, Scope: c.Scope, Budget: c.Budget})
	}
	return gs
}

func toCapabilities(gs []approval.Grant) []broker.Capability {
	cs := make([]broker.Capability, 0, len(gs))
	for _, g := range gs {
		cs = append(cs, broker.Capability{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget})
	}
	return cs
}

// mcpAfterToken is the afterToken that writes the episode's MCP config; nil iff mcpFile is "".
func mcpAfterToken(mcpFile, addr string) func(token string) error {
	if mcpFile == "" {
		return nil
	}
	return func(token string) error { return writeMCPConfig(mcpFile, addr, token) }
}

// renderRequest is the request block --request prints and approve shows above its y/N.
func renderRequest(id string, mr approval.MintRequest, digestHex string, expiresAt, now int64) string {
	var b strings.Builder
	fmt.Fprintf(&b, "mint request %s\n", id)
	fmt.Fprintf(&b, "  digest    %s\n", digestHex)
	fmt.Fprintf(&b, "  expires   %s (in %s)\n", time.Unix(expiresAt, 0).UTC().Format(time.RFC3339), time.Duration(expiresAt-now)*time.Second)
	fmt.Fprintf(&b, "  episode   %s\n", mr.Episode)
	where := mr.Worktree + " (new, from " + mr.Repo + " on branch " + mr.Branch + ")"
	if mr.Reuse {
		where = mr.Worktree + " (existing worktree)"
	}
	fmt.Fprintf(&b, "  worktree  %s\n", where)
	fmt.Fprintf(&b, "  grants    %s\n", grantList(toCapabilities(mr.Grants)))
	fmt.Fprintf(&b, "  session   %ds; token file %s\n", mr.TTLSeconds, mr.SessionFile)
	profile := "unbound"
	if mr.Profile.Project != "" {
		profile = mr.Profile.Project + " " + mr.Profile.Digest + " (" + mr.ProfilesDir + ")"
	}
	fmt.Fprintf(&b, "  profile   %s\n", profile)
	mcp, task := "none", "none"
	if mr.MCPConfig != "" {
		mcp = mr.MCPConfig
	}
	if mr.Task != "" {
		task = mr.Task
	}
	fmt.Fprintf(&b, "  agent     %s %s; MCP config %s\n", mr.Agent, mr.Model, mcp)
	fmt.Fprintf(&b, "  task      %s\n", task)
	fmt.Fprintf(&b, "  daemon    %s build %s\n", mr.Daemon.Addr, daemon.BuildIdentity{VCSRevision: mr.Daemon.VCSRevision,
		VCSModified: mr.Daemon.VCSModified, ExecutableSHA256: mr.Daemon.ExecutableSHA256, GoVersion: mr.Daemon.GoVersion})
	fmt.Fprintf(&b, "  stores    %s, %s\n", mr.WorldDB, mr.CredentialsDB)
	return b.String()
}

var mintRequestFieldMap = map[string]string{
	"verb":               "excluded: display only, the message prefix; approve sets \"approve\"",
	"episode":            "MintRequest.Episode",
	"db":                 "MintRequest.WorldDB",
	"credsDB":            "MintRequest.CredentialsDB",
	"wsRoot":             "MintRequest.WorkspaceRoot",
	"repo":               "MintRequest.Repo",
	"branch":             "MintRequest.Branch",
	"allowNonGit":        "fixed: false (episode start never sets it)",
	"grants":             "MintRequest.Grants",
	"grants[].Effect":    "MintRequest.Grants[].Effect",
	"grants[].Scope":     "MintRequest.Grants[].Scope",
	"grants[].Budget":    "MintRequest.Grants[].Budget",
	"grants[].ExpiresAt": "fixed: 0, set at the mint to now + ttl",
	"ttl":                "MintRequest.TTLSeconds",
	"outPath":            "MintRequest.SessionFile",
	"outDirCreate":       "fixed: true (episode start's value)",
	"profile":            "MintRequest.Profile.Project",
	"profilesDir":        "MintRequest.ProfilesDir",
	"credPath":           "MintRequest.CredentialsDB",
	"dir":                "MintRequest.Worktree",
	"reuse":              "MintRequest.Reuse",
	"degraded":           "MintRequest.Degraded",
	"git":                "MintRequest.Git",
	"prof":               "MintRequest.Profile.Digest",
	"profFile":           "excluded: derived as ProfilesDir/<Profile.Project>.json",
	"afterToken":         "excluded: set by approve from MintRequest.MCPConfig and MintRequest.Daemon.Addr (nil iff no MCP config)",
	"resume":             "nil: set only by episode_resume.go; --request refuses --resume",
	"beforeToken":        "nil: set only by episode_resume.go",
	"resumeRollback":     "nil: set only by episode_resume.go",
	"start":              "nil: set only by episode_preset.go; --request refuses launch presets",
	"startPoint":         "fixed: \"\" (set only by episode_preset.go)",
	"existingBranch":     "fixed: false (set only by episode_preset.go)",
	"presetPrompt":       "excluded: display only; approve sets renderRequest(request)",
}

// requestOnlyLeaves are the MintRequest leaves no mintOptions/mintRequest field maps to: approve
// uses them directly.
var requestOnlyLeaves = []string{"Schema", "MCPConfig", "Agent", "Model", "Task",
	"Daemon.Addr", "Daemon.VCSRevision", "Daemon.VCSModified", "Daemon.ExecutableSHA256", "Daemon.GoVersion"}
