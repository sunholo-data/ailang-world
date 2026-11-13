package main

// Row 201 A4 sprint 3 M6a (the sprint-2 judge's finding 4, JM46): with no --daemon, the decision
// mode reads the installed service's address from $HOME/.ailang-world/service.json. HOME is a temp
// dir whose record names the fake daemon on the loopback (D-WORLD-109: no real service, no real port).

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestApproveDecideDefaultDaemonIsTheService(t *testing.T) {
	f := newFakeDecider(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	paths := homeServicePaths(home)
	rec, err := json.Marshal(serviceRecord{Schema: "world/service/v1", Label: worldServiceLabel, Plist: paths.Plist, Binary: paths.Binary, Address: f.addr()})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(home, ".ailang-world"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".ailang-world", "service.json"), rec, 0o600); err != nil {
		t.Fatal(err)
	}
	if got := defaultServiceAddr(); got != "http://"+f.addr() { // never fall through to a real daemon's port
		t.Fatalf("the temp service record is not read: %q", got)
	}
	var stdout, stderr bytes.Buffer
	env := approveEnv(nil, nil, atTime(1710))
	env.signDecision = (&fakeSigner{ret: func([]byte) ([]byte, error) { return []byte("ARMORED"), nil }}).sign
	if code := runApprove([]string{f.ref, "--db", f.db, "--sign-with", "/k/id"}, &stdout, &stderr, env); code != exitOK || f.gets != 1 || len(f.posts) != 1 {
		t.Fatalf("no --daemon: rc %d, %d GETs, %d POSTs\nstderr %s", code, f.gets, len(f.posts), stderr.String())
	}
}
