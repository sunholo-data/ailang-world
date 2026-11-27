package store

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// OperatorStartReview is a runner-published snapshot. The daemon fills only
// the episode and literal task, using the same mint resolution as the runner.
type OperatorStartReview struct {
	Start                                 EpisodeStart
	Grants                                []approval.Grant
	Profile                               approval.ProfileRef
	Store, Workspace, Sessions            string
	Address, Build, Interpreter           string
	Binary, BinaryDigest, Wall, Committer string
}

func (h OperatorHeartbeat) StartReview(preset string) (OperatorStartReview, error) {
	var catalogue map[string]OperatorStartReview
	if len(h.Catalogue) > 16*approval.MaxOperatorActionBytes || json.Unmarshal(h.Catalogue, &catalogue) != nil {
		return OperatorStartReview{}, ErrOperatorRunner
	}
	v, ok := catalogue[preset]
	if !ok || v.Start.Preset.Name != preset || v.Start.Agent != "codex" {
		return OperatorStartReview{}, ErrOperatorRunner
	}
	return v, nil
}

// Action derives the only browser-start envelope from the published inputs.
// The initially empty ID is assigned only by the request writer.
func (v OperatorStartReview) Action(episode, task string, generation int64) (approval.OperatorAction, error) {
	start := v.Start
	start.Episode, start.Branch = episode, "world/"+episode
	start.Task = &StartTask{
		Path:   filepath.Join(v.Sessions, episode+".task.md"),
		SHA256: fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(task))), Title: StartTaskTitle([]byte(task)),
	}
	in := OperatorMintInput{Attempt: OperatorMintAttempt{Catalogue: generation}, Start: start, Grants: v.Grants,
		Profile: v.Profile, TokenPath: filepath.Join(v.Sessions, episode+".session")}
	mint, err := in.Resolution()
	if err != nil {
		return approval.OperatorAction{}, err
	}
	worktree := filepath.Join(v.Workspace, episode)
	facts := []approval.OperatorResolved{{Name: "launch_binary", Value: v.Binary},
		{Name: "launch_binary_sha256", Value: v.BinaryDigest}, {Name: "launch_wall", Value: v.Wall}, {Name: "mcp_path"}}
	facts = append(facts, mint...)
	facts = append(facts, approval.OperatorResolved{Name: "served_address", Value: v.Address},
		approval.OperatorResolved{Name: "served_build", Value: v.Build}, approval.OperatorResolved{Name: "served_interpreter", Value: v.Interpreter},
		approval.OperatorResolved{Name: "worktree_path", Value: worktree})
	var review strings.Builder
	fmt.Fprintf(&review, "Episode %s from %s (%s)\nRepository %s\nRemote %s = %s\nBase %s @ %s\nNew worktree %s on %s\nHeadless %s %s, cap %s\nCommitter %s\nProfile %s %s\nSession %ds\n",
		episode, start.Preset.Name, start.Preset.SHA256, start.Repo, start.Remote.Name, start.Remote.URL, start.Base.Branch, start.Base.SHA, worktree, start.Branch,
		start.Agent, start.Model, v.Wall, v.Committer, v.Profile.Project, v.Profile.Digest, start.TTLS)
	for _, check := range start.Checks {
		fmt.Fprintf(&review, "Check %s (%s)\n", check.Name, check.Sandbox)
	}
	for _, grant := range v.Grants {
		fmt.Fprintf(&review, "Grant %s %s budget %d\n", grant.Effect, grant.Scope, grant.Budget)
	}
	o := approval.OperatorOptions{Preset: start.Preset.Name, Headless: true}
	return approval.OperatorAction{Schema: approval.OperatorActionSchema, Kind: "start", Episode: episode, Store: v.Store,
		Argv: approval.OperatorArgv("start", episode, o), Options: o, Task: task, Review: review.String(), Resolution: facts}, nil
}
