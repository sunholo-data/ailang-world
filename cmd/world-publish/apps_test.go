package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/appregistry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// Row 274 M1b-2: `world-publish apps`, the attended app-registry writer.

func appsExample(t *testing.T, dir string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(commandRepoRoot(t), "docs", "examples", "apps", "docparse-convert.json"))
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(dir, "docparse-convert.json")
	if err := os.WriteFile(p, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func driveApps(t *testing.T, flags map[string]string, stdin string, getenv func(string) string) armResult {
	t.Helper()
	return drive(t, invocation{verb: "apps", flags: flags}, stdin, getenv, satisfiedProbe(t))
}

func readApps(t *testing.T, storePath string) appregistry.Snapshot {
	t.Helper()
	db, err := store.Open(storePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	snap, err := appregistry.Read(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	return snap
}

func TestAppsVerbPublishVersionAndRetire(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "world.db")
	manifest := appsExample(t, dir)
	pub := map[string]string{"store": storePath, "manifest": manifest}

	res := driveApps(t, pub, appsPublishPhrase("docparse-convert", 1, storePath)+"\n", noEnv)
	if res.code != exitOK || !strings.Contains(res.stdout, "published app docparse-convert version 1 (manifest sha256:") ||
		!strings.Contains(res.stdout, "app registry revision 1") {
		t.Fatalf("first publish = (%d, %q, %q)", res.code, res.stdout, res.stderr)
	}
	if !strings.Contains(res.stdout, "to publish app docparse-convert version 1 (sha256:") {
		t.Fatalf("the prompt must name the app, the version and the manifest ref: %q", res.stdout)
	}
	snap := readApps(t, storePath)
	e, ok := snap.Revision.Lookup("docparse-convert")
	if !ok || e.Version != 1 || snap.Revision.Revision != 1 {
		t.Fatalf("after publish: %+v", snap)
	}

	// The same bytes again are refused BEFORE the gate: no prompt, nothing written.
	again := driveApps(t, pub, appsPublishPhrase("docparse-convert", 2, storePath)+"\n", noEnv)
	if again.code != exitError || !strings.Contains(again.stderr, "already pins") || strings.Contains(again.stdout, "Type exactly") {
		t.Fatalf("republish = (%d, %q, %q)", again.code, again.stdout, again.stderr)
	}

	// Edited bytes become version 2 over revision 1.
	raw, _ := os.ReadFile(manifest)
	if err := os.WriteFile(manifest, bytes.Replace(raw, []byte(`"budget": 2000`), []byte(`"budget": 1500`), 1), 0o644); err != nil {
		t.Fatal(err)
	}
	v2 := driveApps(t, pub, appsPublishPhrase("docparse-convert", 2, storePath)+"\n", noEnv)
	if v2.code != exitOK || !strings.Contains(v2.stdout, "version 2") || !strings.Contains(v2.stdout, "revision 2") {
		t.Fatalf("second version = (%d, %q, %q)", v2.code, v2.stdout, v2.stderr)
	}
	snap2 := readApps(t, storePath)
	if e2, _ := snap2.Revision.Lookup("docparse-convert"); e2.Version != 2 || e2.Manifest == e.Manifest || snap2.Revision.Parent != snap.Head.String() {
		t.Fatalf("version 2 = %+v (parent %s, want %s)", snap2.Revision, snap2.Revision.Parent, snap.Head)
	}

	ret := map[string]string{"store": storePath, "retire": "docparse-convert"}
	r := driveApps(t, ret, appsRetirePhrase("docparse-convert", storePath)+"\n", noEnv)
	if r.code != exitOK || !strings.Contains(r.stdout, "retired app docparse-convert: app registry revision 3") {
		t.Fatalf("retire = (%d, %q, %q)", r.code, r.stdout, r.stderr)
	}
	if snap3 := readApps(t, storePath); len(snap3.Revision.Apps) != 0 || snap3.Revision.Revision != 3 {
		t.Fatalf("after retire: %+v", snap3.Revision)
	}
	if r2 := driveApps(t, ret, appsRetirePhrase("docparse-convert", storePath)+"\n", noEnv); r2.code != exitError || !strings.Contains(r2.stderr, "no app docparse-convert to retire") {
		t.Fatalf("retire again = (%d, %q)", r2.code, r2.stderr)
	}
}

func TestAppsVerbRefusesBeforeTheGate(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "world.db")
	manifest := appsExample(t, dir)
	raw, _ := os.ReadFile(manifest)

	invalid := filepath.Join(dir, "wide.json")
	_ = os.WriteFile(invalid, bytes.Replace(bytes.Replace(raw, []byte(`"docparse-convert"`), []byte(`"wide"`), 1), []byte(`"Workspace.Read"`), []byte(`"Model.Infer"`), 1), 0o644)
	renamed := filepath.Join(dir, "other-name.json")
	_ = os.WriteFile(renamed, raw, 0o644)

	cases := map[string]struct {
		flags map[string]string
		code  int
		want  string
	}{
		"invalid manifest":  {map[string]string{"store": storePath, "manifest": invalid}, exitError, "breaks the law"},
		"basename mismatch": {map[string]string{"store": storePath, "manifest": renamed}, exitError, "but the file is other-name.json"},
		"absent file":       {map[string]string{"store": storePath, "manifest": filepath.Join(dir, "nope.json")}, exitError, "read "},
		"neither":           {map[string]string{"store": storePath}, exitUsage, "exactly one of --manifest"},
		"both":              {map[string]string{"store": storePath, "manifest": manifest, "retire": "x"}, exitUsage, "exactly one of --manifest"},
		"retire on empty":   {map[string]string{"store": storePath, "retire": "docparse-convert"}, exitError, "registry is empty"},
	}
	for name, c := range cases {
		res := driveApps(t, c.flags, "anything\n", noEnv)
		if res.code != c.code || !strings.Contains(res.stderr, c.want) || strings.Contains(res.stdout, "Type exactly") {
			t.Errorf("%s = (%d, %q, %q), want exit %d naming %q with no prompt", name, res.code, res.stdout, res.stderr, c.code, c.want)
		}
	}
	if snap := readApps(t, storePath); snap.Present {
		t.Fatalf("a refusal wrote a head: %+v", snap)
	}
}

func TestAppsVerbFences(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "world.db")
	manifest := appsExample(t, dir)
	base := map[string]string{"store": storePath, "manifest": manifest}
	right := appsPublishPhrase("docparse-convert", 1, storePath) + "\n"

	if res := drive(t, invocation{verb: "apps", flags: base, bools: []string{"live"}}, right, noEnv, satisfiedProbe(t)); res.code != exitStop || !strings.Contains(res.stderr, "--live") {
		t.Fatalf("--live = (%d, %q), want STOP naming --live", res.code, res.stderr)
	}
	if res := driveApps(t, base, right, ciEnv); res.code != exitStop {
		t.Fatalf("ci = (%d, %q), want STOP", res.code, res.stderr)
	}
	for name, line := range map[string]string{
		"package-publish phrase": attendedPhrase + "\n",
		"wrong version":          appsPublishPhrase("docparse-convert", 2, storePath) + "\n",
		"wrong store":            appsPublishPhrase("docparse-convert", 1, storePath+".other") + "\n",
		"retire phrase":          appsRetirePhrase("docparse-convert", storePath) + "\n",
		"yes":                    "yes\n",
	} {
		res := driveApps(t, base, line, noEnv)
		if res.code != exitStop || !strings.Contains(res.stderr, "fence=confirmation reason=mismatch") {
			t.Errorf("%s = (%d, %q), want STOP fence=confirmation reason=mismatch", name, res.code, res.stderr)
		}
	}
	if snap := readApps(t, storePath); snap.Present {
		t.Fatalf("a stopped write moved the head: %+v", snap)
	}
	if res := driveApps(t, base, right, noEnv); res.code != exitOK {
		t.Fatalf("the right phrase = (%d, %q)", res.code, res.stderr)
	}
	if got := appsPublishPhrase("docparse-convert", 1, storePath); got != "publish app docparse-convert version 1 to "+storePath {
		t.Fatalf("phrase oracle = %q", got)
	}
}
