package main

// Row 274 M1 (design_docs/planned/w-world-applications.md §2): the read-only app
// catalog verbs. An app is an extension package with an app.json manifest; M1 reads a
// local catalog (~/.ailang-world/apps/<name>.json) and validates every manifest against
// the law in design_docs/sketches/appmanifest.ail (Go mirror: host/appmanifest). Nothing
// here launches anything or writes anywhere.

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/sunholo-data/ailang-world/host/appmanifest"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

const appHelp = `usage: ailang-worldd app list|show|check [name|file.json] [--apps-dir D] [--json]
  list                       every app in the catalog, valid or not
  show <name>                one app's manifest, digest and validation
  check <name|file.json>     validate one manifest; exit 0 only when it is valid

An app is an extension package with an app.json manifest (schema world/app/v1): a
preset template of existing effects only, declared outputs and catalog-only views.
The catalog defaults to ~/.ailang-world/apps, one <name>.json per app. These verbs
only read it; installing is copying a checked manifest there.
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
		return appList(*dir, *asJSON, out, errOut)
	case "show", "check":
		if len(pos) != 1 {
			fmt.Fprint(errOut, appHelp)
			return exitUsage
		}
		e, ok := appEntry(*dir, pos[0], verb == "check", errOut)
		if !ok {
			return exitUsage
		}
		if verb == "show" {
			return appShow(e, *asJSON, out)
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

func appList(dir string, asJSON bool, out, errOut io.Writer) int {
	entries, err := appmanifest.ReadCatalog(dir)
	if err != nil {
		fmt.Fprintf(errOut, "cannot read the app catalog %s: %v\n  why: app list reads one <name>.json per app there\n  next: check the directory's permissions, or pass --apps-dir\n", dir, err)
		return exitFatal
	}
	if asJSON {
		return writeAppJSON(out, entries)
	}
	if len(entries) == 0 {
		fmt.Fprintf(out, "no apps in %s\n", dir)
		return exitOK
	}
	for _, e := range entries {
		switch {
		case e.Error != "":
			fmt.Fprintf(out, "%-24s ERROR %s\n", e.Name, e.Error)
		case !e.Valid():
			fmt.Fprintf(out, "%-24s INVALID %d problem(s): ailang-worldd app check %s\n", e.Name, len(e.Problems), e.Name)
		default:
			m := e.Manifest
			fmt.Fprintf(out, "%-24s %-6s %-11s %s %s\n", e.Name, m.Preset.Agent, m.Preset.Launch.Mode, e.Digest[7:19], m.Title)
		}
	}
	return exitOK
}

func appShow(e appmanifest.Entry, asJSON bool, out io.Writer) int {
	if asJSON {
		return writeAppJSON(out, e)
	}
	fmt.Fprintf(out, "app %s\n  file    %s\n  digest  %s\n", e.Name, e.File, e.Digest)
	if e.Error != "" {
		fmt.Fprintf(out, "  error   %s\n", e.Error)
		return exitOK
	}
	m := e.Manifest
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
	writeAppVerdict(out, e)
	return exitOK
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
