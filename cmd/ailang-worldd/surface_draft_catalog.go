package main

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Only expose the operator-selected fields needed for drafting. In particular,
// the raw preset includes local service paths and may contain a remote URL.
type draftPreset struct {
	Name     string              `json:"name"`
	Repo     string              `json:"repo"`
	Base     string              `json:"base"`
	Commit   string              `json:"commit"`
	Host     string              `json:"host"`
	Model    string              `json:"model"`
	Lifetime int64               `json:"seconds"`
	Profile  string              `json:"profile"`
	Digest   string              `json:"digest"`
	Grants   []broker.Capability `json:"grants"`
	Checks   []store.StartCheck  `json:"checks"`
}

func surfaceDraftCatalog(ctx context.Context, paths []string, profiles, root string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	var catalog []draftPreset
	for _, path := range paths {
		if filepath.Ext(path) != ".json" {
			return nil, fmt.Errorf("expected a .json preset file")
		}
		name := strings.TrimSuffix(filepath.Base(path), ".json")
		r, err := loadLaunchPreset(ctx, name, filepath.Dir(path), profiles, root)
		if err != nil {
			return nil, err
		}
		p := r.Preset
		catalog = append(catalog, draftPreset{Name: p.Name, Repo: p.Repo, Base: p.Base, Commit: r.BaseSHA, Host: p.Agent, Model: p.Model, Lifetime: p.TTLS, Profile: p.Profile, Digest: r.SHA256, Grants: r.Grants, Checks: r.Checks})
	}
	return json.Marshal(catalog)
}
