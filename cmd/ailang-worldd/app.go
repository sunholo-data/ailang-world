package main

// Row 274 M1 (design_docs/planned/w-world-applications.md §2): the read-only app
// catalog verbs. An app is an extension package with an app.json manifest; M1 reads a
// local catalog (~/.ailang-world/apps/<name>.json) and validates every manifest against
// the law in design_docs/sketches/appmanifest.ail (Go mirror: host/appmanifest). Nothing
// here launches anything or writes anywhere.

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/appmanifest"
	"github.com/sunholo-data/ailang-world/host/appregistry"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

const appHelp = `usage: ailang-worldd app list|show|check [name|file.json] [--db W] [--apps-dir D] [--json]
  list                       every app: published (the app registry in --db) and local drafts
  show <name>                one app's state, pinned manifest, grants and validation
  check <name|file.json>     validate one local manifest; exit 0 only when it is valid

An app is an extension package with an app.json manifest (schema world/app/v1): a
preset template of existing effects only, declared outputs and catalog-only views.
The app registry (world/app-registry/v1, read-only from --db, default the service store)
holds published apps; only an attended world-publish apps writes it. The local catalog
(~/.ailang-world/apps, one <name>.json per app) holds drafts. Only published apps launch.
Example:
  ailang-worldd app check docs/examples/apps/docparse-convert.json
  mkdir -p ~/.ailang-world/apps
  cp docs/examples/apps/docparse-convert.json ~/.ailang-world/apps/
  ailang-worldd app list
`

func runApp(args []string, out, errOut io.Writer) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(out, appHelp)
		return exitOK
	}
	verb := args[0]
	fs := flag.NewFlagSet("app "+verb, flag.ContinueOnError)
	fs.SetOutput(errOut)
	dir := fs.String("apps-dir", appmanifest.DefaultDir(), "app catalog directory")
	db := fs.String("db", defaultWorldPaths().db, "the world store whose app registry to read, read-only (default ~/.ailang/world/world.db)")
	asJSON := fs.Bool("json", false, "JSON output")
	pos, err := verbs.ParseInterleaved(fs, args[1:])
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(out, appHelp)
		return exitOK
	}
	if err != nil {
		fmt.Fprintln(errOut, err)
		return exitUsage
	}
	switch verb {
	case "list":
		if len(pos) != 0 {
			fmt.Fprint(errOut, appHelp)
			return exitUsage
		}
		return appList(*db, *dir, *asJSON, out, errOut)
	case "show":
		if len(pos) != 1 {
			fmt.Fprint(errOut, appHelp)
			return exitUsage
		}
		return appShow(*db, *dir, pos[0], *asJSON, out, errOut)
	case "check":
		if len(pos) != 1 {
			fmt.Fprint(errOut, appHelp)
			return exitUsage
		}
		e, ok := appEntry(*dir, pos[0], true, errOut)
		if !ok {
			return exitUsage
		}
		return appCheck(e, *asJSON, out)
	}
	fmt.Fprintf(errOut, "unknown app verb %q\n  why: app verbs are list, show and check\n  next: ailang-worldd app --help\n", verb)
	return exitUsage
}

// appEntry finds one catalog entry by name, or (for check) a manifest file by path.
func appEntry(dir, arg string, allowFile bool, errOut io.Writer) (appmanifest.Entry, bool) {
	if allowFile && strings.HasSuffix(arg, ".json") {
		return appmanifest.ReadEntry(arg), true
	}
	if !appmanifest.NameOk(arg) {
		fmt.Fprintf(errOut, "invalid app name %q\n  why: app names use the preset grammar ^[a-z0-9][a-z0-9-]*$, 1 to 51 characters\n  next: ailang-worldd app list --apps-dir %s\n", arg, dir)
		return appmanifest.Entry{}, false
	}
	entries, err := appmanifest.ReadCatalog(dir)
	if err != nil {
		fmt.Fprintf(errOut, "cannot read the app catalog %s: %v\n  why: app verbs read one <name>.json per app there\n  next: check the directory's permissions, or pass --apps-dir\n", dir, err)
		return appmanifest.Entry{}, false
	}
	var names []string
	for _, e := range entries {
		if e.Name == arg {
			return e, true
		}
		names = append(names, e.Name)
	}
	have := "none"
	if len(names) > 0 {
		have = strings.Join(names, ", ")
	}
	fmt.Fprintf(errOut, "no app %s in %s (have: %s)\n  why: an app is installed by copying its checked app.json there as %s.json\n  next: ailang-worldd app check <file.json>, then copy it into %s\n", arg, dir, have, arg, dir)
	return appmanifest.Entry{}, false
}

// appListing reads the merged view: the app registry in db (read-only, no writer lock, so it
// works beside a running daemon; a missing store file is an empty registry) and the local
// catalog in dir.
func appListing(db, dir string, errOut io.Writer) ([]appregistry.Listing, string, bool) {
	local, err := appmanifest.ReadCatalog(dir)
	if err != nil {
		fmt.Fprintf(errOut, "cannot read the app catalog %s: %v\n  why: app verbs read one <name>.json per app there\n  next: check the directory's permissions, or pass --apps-dir\n", dir, err)
		return nil, "", false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := os.Stat(db); errors.Is(err, os.ErrNotExist) {
		return appregistry.List(ctx, nil, appregistry.Snapshot{}, local), fmt.Sprintf("no world store at %s, so no published apps", db), true
	}
	ro, err := store.OpenReadOnly(db)
	if err != nil {
		fmt.Fprintf(errOut, "cannot read the world store %s: %v\n  why: published apps live in its app registry\n  next: pass --db <world.db>, or ailang-worldd doctor\n", db, err)
		return nil, "", false
	}
	defer func() { _ = ro.Close() }()
	snap, err := appregistry.Read(ctx, ro)
	if err != nil {
		fmt.Fprintf(errOut, "cannot read the app registry in %s: %v\n  next: ailang-worldd doctor\n", db, err)
		return nil, "", false
	}
	note := fmt.Sprintf("app registry: none yet in %s", db)
	if snap.Present {
		note = fmt.Sprintf("app registry: revision %d in %s", snap.Revision.Revision, db)
	}
	return appregistry.List(ctx, ro, snap, local), note, true
}

// appStateLabel is the listing state in words.
func appStateLabel(l appregistry.Listing) string {
	switch l.State {
	case appregistry.StatePublished:
		return fmt.Sprintf("published v%d", l.Version)
	case appregistry.StateLocalDiffers:
		return fmt.Sprintf("published v%d, local copy differs", l.Version)
	}
	return "local (not published)"
}

// appRow is the JSON shape of one listing row (no file paths beyond the local file's).
type appRow struct {
	Name        string                `json:"name"`
	State       string                `json:"state"`
	Launchable  bool                  `json:"launchable"`
	Valid       bool                  `json:"valid"`
	Version     int64                 `json:"version,omitempty"`
	Manifest    string                `json:"manifest,omitempty"`
	PublishedAt string                `json:"published_at,omitempty"`
	Title       string                `json:"title,omitempty"`
	Agent       string                `json:"agent,omitempty"`
	Mode        string                `json:"mode,omitempty"`
	LocalFile   string                `json:"local_file,omitempty"`
	LocalDigest string                `json:"local_digest,omitempty"`
	Problems    []appmanifest.Problem `json:"problems,omitempty"`
	Error       string                `json:"error,omitempty"`
}

func toAppRow(l appregistry.Listing) appRow {
	r := appRow{Name: l.Name, State: l.State, Launchable: l.Launchable(), Valid: l.Valid(), Version: l.Version,
		Manifest: l.ManifestRef, PublishedAt: l.PublishedAt, Problems: l.Problems, Error: l.Error}
	if m := l.Manifest; m != nil {
		r.Title, r.Agent, r.Mode = m.Title, m.Preset.Agent, m.Preset.Launch.Mode
	}
	if l.Local != nil {
		r.LocalFile, r.LocalDigest = l.Local.File, l.Local.Digest
	}
	return r
}

func appList(db, dir string, asJSON bool, out, errOut io.Writer) int {
	rows, note, ok := appListing(db, dir, errOut)
	if !ok {
		return exitFatal
	}
	if asJSON {
		js := make([]appRow, 0, len(rows))
		for _, l := range rows {
			js = append(js, toAppRow(l))
		}
		return writeAppJSON(out, js)
	}
	fmt.Fprintln(out, note)
	if len(rows) == 0 {
		fmt.Fprintf(out, "no apps (catalog %s)\n", dir)
		return exitOK
	}
	for _, l := range rows {
		switch {
		case l.Error != "":
			fmt.Fprintf(out, "%-24s %-34s ERROR %s\n", l.Name, appStateLabel(l), l.Error)
		case !l.Valid():
			fmt.Fprintf(out, "%-24s %-34s INVALID %d problem(s): ailang-worldd app check %s\n", l.Name, appStateLabel(l), len(l.Problems), l.Name)
		default:
			m := l.Manifest
			fmt.Fprintf(out, "%-24s %-34s %-6s %-11s %s\n", l.Name, appStateLabel(l), m.Preset.Agent, m.Preset.Launch.Mode, m.Title)
		}
	}
	return exitOK
}

func appShow(db, dir, name string, asJSON bool, out, errOut io.Writer) int {
	if !appmanifest.NameOk(name) {
		fmt.Fprintf(errOut, "invalid app name %q\n  why: app names use the preset grammar ^[a-z0-9][a-z0-9-]*$, 1 to 51 characters\n  next: ailang-worldd app list\n", name)
		return exitUsage
	}
	rows, _, ok := appListing(db, dir, errOut)
	if !ok {
		return exitFatal
	}
	var names []string
	for _, l := range rows {
		names = append(names, l.Name)
		if l.Name != name {
			continue
		}
		if asJSON {
			return writeAppJSON(out, toAppRow(l))
		}
		fmt.Fprintf(out, "app %s\n  state   %s\n", l.Name, appStateLabel(l))
		if l.ManifestRef != "" {
			fmt.Fprintf(out, "  pinned  %s (published %s)\n", l.ManifestRef, l.PublishedAt)
		}
		if l.Local != nil {
			fmt.Fprintf(out, "  local   %s (%s)\n", l.Local.File, l.Local.Digest)
		}
		if l.Error != "" {
			fmt.Fprintf(out, "  error   %s\n", l.Error)
			return exitOK
		}
		writeManifestDetails(out, l.Manifest)
		if len(l.Problems) > 0 {
			fmt.Fprintf(out, "INVALID %s: %d problem(s)\n", l.Name, len(l.Problems))
			for _, p := range l.Problems {
				fmt.Fprintf(out, "  %s\n", p)
			}
		} else if appregistry.Launchable(l.State) {
			fmt.Fprintf(out, "LAUNCHABLE %s v%d\n", l.Name, l.Version)
		} else {
			fmt.Fprintf(out, "NOT PUBLISHED %s: only published apps launch (D-WORLD-129)\n  next: world-publish apps --store %s --manifest %s\n", l.Name, db, l.Local.File)
		}
		return exitOK
	}
	have := "none"
	if len(names) > 0 {
		have = strings.Join(names, ", ")
	}
	fmt.Fprintf(errOut, "no app %s in the registry or %s (have: %s)\n  next: ailang-worldd app check <file.json>, then copy it into %s or publish it with world-publish apps\n", name, dir, have, dir)
	return exitUsage
}

func writeManifestDetails(out io.Writer, m *appmanifest.Manifest) {
	if m == nil {
		return
	}
	fmt.Fprintf(out, "  title   %s\n  summary %s\n  agent   %s, %s, deadline %ds, ttl %ds\n", m.Title, m.Summary, m.Preset.Agent, m.Preset.Launch.Mode, m.Preset.Launch.DeadlineS, m.Preset.TTLS)
	for _, g := range m.Preset.Grants {
		fmt.Fprintf(out, "  grant   %s=%s:%d\n", g.Effect, g.Scope, g.Budget)
	}
	if m.Preset.ExecProfile != "" {
		fmt.Fprintf(out, "  profile %s\n", m.Preset.ExecProfile)
	}
	for _, c := range m.Preset.Checks {
		fmt.Fprintf(out, "  check   %s: %s\n", c.Name, strings.Join(c.Argv, " "))
	}
	for _, o := range m.Outputs {
		fmt.Fprintf(out, "  output  %s %s", o.Kind, o.Glob)
		if o.Check != "" {
			fmt.Fprintf(out, " (checked by %s)", o.Check)
		}
		fmt.Fprintln(out)
	}
	if len(m.Views) > 0 {
		fmt.Fprintf(out, "  views   %s\n", strings.Join(m.Views, ", "))
	}
}

func appCheck(e appmanifest.Entry, asJSON bool, out io.Writer) int {
	if asJSON {
		writeAppJSON(out, map[string]any{"name": e.Name, "file": e.File, "digest": e.Digest, "valid": e.Valid(), "problems": e.Problems, "error": e.Error})
	} else {
		writeAppVerdict(out, e)
	}
	if !e.Valid() {
		return exitUsage
	}
	return exitOK
}

func writeAppVerdict(out io.Writer, e appmanifest.Entry) {
	switch {
	case e.Error != "":
		fmt.Fprintf(out, "INVALID %s: %s\n  next: fix the JSON, then ailang-worldd app check %s\n", e.Name, e.Error, e.File)
	case len(e.Problems) > 0:
		fmt.Fprintf(out, "INVALID %s: %d problem(s)\n", e.Name, len(e.Problems))
		for _, p := range e.Problems {
			fmt.Fprintf(out, "  %s\n", p)
		}
		fmt.Fprintf(out, "  next: fix each field, then ailang-worldd app check %s\n", e.File)
	default:
		fmt.Fprintf(out, "VALID %s (%s)\n", e.Name, e.Digest)
	}
}

func writeAppJSON(out io.Writer, v any) int {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return exitFatal
	}
	fmt.Fprintln(out, string(raw))
	return exitOK
}
