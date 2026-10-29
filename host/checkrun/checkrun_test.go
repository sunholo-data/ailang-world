package checkrun

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCheckrunHeadsAndSessionScrub(t *testing.T) {
	t.Setenv("WORLD_SESSION", "sentinel")
	r, err := Run(context.Background(), []string{"/bin/sh", "-c", "test -z \"$WORLD_SESSION\" || exit 9; i=0; while test $i -lt 500; do printf x; printf y >&2; i=$((i+1)); done"}, t.TempDir(), time.Second)
	if err != nil || r.ExitCode != 0 || r.Stdout != strings.Repeat("x", HeadBytes) || r.Stderr != strings.Repeat("y", HeadBytes) {
		t.Fatalf("result %+v err %v", r, err)
	}
}

func TestCheckrunWallCapKillsProcessGroup(t *testing.T) {
	dir := t.TempDir()
	start := time.Now()
	r, err := Run(context.Background(), []string{"/bin/sh", "-c", "(sleep .4; echo escaped > marker) & wait"}, dir, 50*time.Millisecond)
	if err == nil || time.Since(start) > time.Second {
		t.Fatalf("wall cap: %+v %v in %s", r, err, time.Since(start))
	}
	time.Sleep(450 * time.Millisecond)
	if _, err := os.Stat(filepath.Join(dir, "marker")); !os.IsNotExist(err) {
		t.Fatal("descendant survived wall cap")
	}
}

func TestCheckrunDeclaresInvalidInputs(t *testing.T) {
	if _, err := Run(context.Background(), nil, t.TempDir(), time.Second); err == nil {
		t.Fatal("empty argv accepted")
	}
}
