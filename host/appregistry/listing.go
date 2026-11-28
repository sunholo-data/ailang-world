package appregistry

import (
	"context"
	"sort"

	"github.com/sunholo-data/ailang-world/host/appmanifest"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

func parseRef(s string) (hashref.HashRef, error) { return hashref.Parse(s) }

// Listing is one app name in the merged view of the registry and the local catalog
// (row 274 M1b-3, w-app-registry §4). State is listingState's verdict.
type Listing struct {
	Name  string `json:"name"`
	State string `json:"state"`
	// Published: the registry entry, its pinned manifest, and the manifest's validity.
	Version     int64                 `json:"version,omitempty"`
	ManifestRef string                `json:"manifest,omitempty"`
	PublishedAt string                `json:"published_at,omitempty"`
	Manifest    *appmanifest.Manifest `json:"-"`
	// Local: the catalog entry (a draft when the app is published).
	Local *appmanifest.Entry `json:"-"`
	// Problems and Error describe the row's manifest: the pinned one when published, else the
	// local file's.
	Problems []appmanifest.Problem `json:"problems,omitempty"`
	Error    string                `json:"error,omitempty"`
}

// Valid reports whether the row's manifest parsed and passed the law.
func (l Listing) Valid() bool { return l.Error == "" && len(l.Problems) == 0 && l.Manifest != nil }

// Launchable reports whether M2 may launch this row (D-WORLD-129 = A).
func (l Listing) Launchable() bool { return Launchable(l.State) && l.Valid() }

// List merges snap (the registry, possibly absent) with the local catalog entries, one row
// per name, sorted by name. A published row carries its pinned manifest, read through s and
// re-validated; a pinned object that is missing or no longer valid becomes the row's Error.
func List(ctx context.Context, s ObjectStore, snap Snapshot, local []appmanifest.Entry) []Listing {
	byName := map[string]*appmanifest.Entry{}
	for i := range local {
		byName[local[i].Name] = &local[i]
	}
	var rows []Listing
	seen := map[string]bool{}
	if snap.Present {
		for _, e := range snap.Revision.Apps {
			seen[e.Name] = true
			loc := byName[e.Name]
			row := Listing{Name: e.Name, Version: e.Version, ManifestRef: e.Manifest, PublishedAt: e.PublishedAt, Local: loc}
			row.State = listingState(true, loc != nil, loc != nil && loc.Digest == e.Manifest)
			if ref, err := parseRef(e.Manifest); err != nil {
				row.Error = err.Error()
			} else if m, _, err := ReadManifest(ctx, s, ref); err != nil {
				row.Error = err.Error()
			} else {
				row.Manifest = &m
			}
			rows = append(rows, row)
		}
	}
	for i := range local {
		e := &local[i]
		if seen[e.Name] {
			continue
		}
		rows = append(rows, Listing{Name: e.Name, State: listingState(false, true, false), Local: e,
			Manifest: e.Manifest, Problems: e.Problems, Error: e.Error})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
	return rows
}
