package checkrun

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestRunEnvProfile(t *testing.T) {
	t.Setenv("WORLD_CHECK_VALUE", "inherited")
	t.Setenv("WORLD_SESSION", "secret")
	t.Setenv("PATH", "/usr/bin:/bin")
	r, err := RunEnv(t.Context(), []string{"/bin/sh", "-c", `printf '%s|%s|%s' "$WORLD_CHECK_VALUE" "$PATH" "${WORLD_SESSION-unset}"`}, t.TempDir(), time.Second,
		map[string]string{"WORLD_CHECK_VALUE": "profile", "WORLD_SESSION": "also-secret"}, []string{"/first", "/second"})
	want := "profile|" + strings.Join([]string{"/first", "/second", "/usr/bin", "/bin"}, string(os.PathListSeparator)) + "|unset"
	if err != nil || r.ExitCode != 0 || r.Stdout != want {
		t.Fatalf("%+v %v, want %q", r, err, want)
	}
}
