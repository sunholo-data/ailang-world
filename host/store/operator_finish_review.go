package store

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sunholo-data/ailang-world/host/approval"
)

type OperatorFinishPath struct {
	Path, Mode, Digest string
}

type OperatorFinishTree struct {
	Head, Branch, Index, Config string
	Paths                       []OperatorFinishPath
}

type OperatorLocalCommit struct {
	Snapshot           OperatorFinishTree
	Message, Committer string
}

// OperatorFinishReview is the bounded manifest published by the runner. The
// daemon fills only the draft decision and signed catalogue generation.
type OperatorFinishReview struct {
	Store, Episode, Preset, Address, Worktree string
	Report, PresetDigest, Profile, Sessions   string
	Launch                                    string
	Remote                                    StartRemote
	Base                                      StartBase
	Checks                                    []StartCheck
	Protected, Unattributed                   []string
	Commit                                    OperatorLocalCommit
}

func (v OperatorFinishReview) Action(generation int64, ready bool) (approval.OperatorAction, error) {
	manifest, err := json.Marshal(v)
	if err != nil {
		return approval.OperatorAction{}, err
	}
	intent := "draft"
	if ready {
		intent = "ready"
	}
	o := approval.OperatorOptions{Preset: v.Preset, ConfirmOnSurface: true, Ready: ready}
	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, Kind: "finish", Episode: v.Episode, Store: v.Store,
		Argv: approval.OperatorArgv("finish", v.Episode, o), Options: o,
		Resolution: []approval.OperatorResolved{{Name: "finish_catalogue", Value: strconv.FormatInt(generation, 10)},
			{Name: "finish_manifest", Value: string(manifest)}},
		Review: fmt.Sprintf("Prepare local commit for %s (%s pull request review follows)\n%s", v.Episode, intent, manifest),
	}
	_, err = approval.OperatorReviewDigest(a)
	return a, err
}

const OperatorFinishCatalogueKey = "@finish"

func (h OperatorHeartbeat) FinishReview(episode, preset string) (OperatorFinishReview, error) {
	var catalogue map[string]json.RawMessage
	var entries map[string]OperatorFinishReview
	if len(h.Catalogue) > 16*approval.MaxOperatorActionBytes || json.Unmarshal(h.Catalogue, &catalogue) != nil ||
		json.Unmarshal(catalogue[OperatorFinishCatalogueKey], &entries) != nil {
		return OperatorFinishReview{}, ErrOperatorRunner
	}
	v, ok := entries[episode]
	if !ok || v.Episode != episode || v.Preset != preset {
		return OperatorFinishReview{}, ErrOperatorRunner
	}
	return v, nil
}
