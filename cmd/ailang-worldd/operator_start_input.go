package main

import (
	"context"
	"encoding/json"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

// operatorStartFacts projects one live read of a headless start into the
// effective mint input and the complete, name-ordered resolution a reviewed
// action signs. It derives the task only from the signed task bytes.
func operatorStartFacts(ctx context.Context, signed approval.OperatorAction, attempt store.OperatorMintAttempt,
	presets, profiles, workspace, hosts string, clock func() int64) (store.OperatorMintInput, []approval.OperatorResolved, error) {
	if signed.Kind != "start" || !signed.Options.Headless {
		return store.OperatorMintInput{}, nil, store.ErrOperatorConflict
	}
	read, err := resolvePresetStart(ctx, signed.Options.Preset, presets, profiles, presetStartInput{
		episodes: []string{signed.Episode}, workspaceRoot: workspace, hostsDir: hosts, launching: true, headless: true,
		readTask: func() (string, []byte, error) { return "inline", []byte(signed.Task), nil },
	}, clock)
	if err != nil {
		return store.OperatorMintInput{}, nil, err
	}
	req := read.mint
	if req.prof == nil {
		return store.OperatorMintInput{}, nil, store.ErrOperatorInput
	}
	served, err := store.CanonicalPath(req.credPath)
	if err != nil {
		return store.OperatorMintInput{}, nil, err
	}
	if served != signed.Store {
		return store.OperatorMintInput{}, nil, &approval.BindingError{Field: "store"}
	}
	start, task := *req.start, *req.start.Task
	task.Path = read.taskCopy
	start.Task = &task
	in := store.OperatorMintInput{Attempt: attempt, Start: start, TokenPath: req.outPath,
		Profile: approval.ProfileRef{Project: req.profile, Digest: req.prof.Digest}}
	for _, g := range req.grants {
		in.Grants = append(in.Grants, approval.Grant{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget})
	}
	mint, err := in.Resolution()
	if err != nil {
		return store.OperatorMintInput{}, nil, err
	}
	build, err := json.Marshal(read.health.Build)
	if err != nil {
		return store.OperatorMintInput{}, nil, err
	}
	facts := []approval.OperatorResolved{{Name: "launch_binary", Value: read.launchBinary},
		{Name: "launch_wall", Value: read.launchWall.String()}, {Name: "mcp_path", Value: read.mcp}}
	facts = append(facts, mint...)
	facts = append(facts, approval.OperatorResolved{Name: "served_build", Value: string(build)},
		approval.OperatorResolved{Name: "served_interpreter", Value: read.health.ToolRef},
		approval.OperatorResolved{Name: "worktree_path", Value: req.dir})
	return in, facts, nil
}

// reconstructOperatorStart requires every signed resolution fact to equal its
// live value before any effect. A change, even one lowering authority, needs a
// fresh review. It locks, writes, mints and launches nothing, and it has no
// production caller until the paired-process consumer lands.
func reconstructOperatorStart(ctx context.Context, signed approval.OperatorAction, attempt store.OperatorMintAttempt,
	presets, profiles, workspace, hosts string, clock func() int64) (store.OperatorMintInput, error) {
	in, facts, err := operatorStartFacts(ctx, signed, attempt, presets, profiles, workspace, hosts, clock)
	if err != nil {
		return store.OperatorMintInput{}, err
	}
	if len(facts) != len(signed.Resolution) {
		return store.OperatorMintInput{}, &approval.BindingError{Field: "resolution"}
	}
	for i, fact := range facts {
		if signed.Resolution[i] != fact {
			return store.OperatorMintInput{}, &approval.BindingError{Field: fact.Name}
		}
	}
	return in, nil
}
