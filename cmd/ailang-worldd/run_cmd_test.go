package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestRunWrapperBuildsEpisodeStartArgs(t *testing.T) {
	orig := runEpisodeStartFunc
	defer func() { runEpisodeStartFunc = orig }()

	var gotArgs []string
	var gotTaskBody string
	runEpisodeStartFunc = func(args []string, stdout, stderr io.Writer, env sessionEnv) int {
		gotArgs = append([]string(nil), args...)
		for i := 0; i+1 < len(args); i++ {
			if args[i] == "--task" {
				b, err := os.ReadFile(args[i+1])
				if err == nil {
					gotTaskBody = string(b)
				}
			}
		}
		return exitOK
	}

	// 1. Default inline task -> --preset, --launch, --headless, --from-allowance, --task <tmp>
	var out, errOut bytes.Buffer
	if rc := runWrapper([]string{"world-agy", "fix", "the", "bug"}, &out, &errOut, sessionEnv{}); rc != exitOK {
		t.Fatalf("inline run failed rc=%d stderr=%s", rc, errOut.String())
	}
	if gotTaskBody != "fix the bug\n" {
		t.Fatalf("got task body %q, want %q", gotTaskBody, "fix the bug\n")
	}
	wantPrefix := []string{"--preset", "world-agy", "--launch", "--headless", "--from-allowance", "--task"}
	if len(gotArgs) != 7 || !reflect.DeepEqual(gotArgs[:6], wantPrefix) {
		t.Fatalf("got args %q, want prefix %q + tmp path", gotArgs, wantPrefix)
	}
	// Temp file must be removed after runWrapper returns.
	if _, err := os.Stat(gotArgs[6]); !os.IsNotExist(err) {
		t.Fatalf("expected temp task file %q to be cleaned up, err=%v", gotArgs[6], err)
	}

	// 2. Interactive + --prompt + --episode
	gotArgs = nil
	out.Reset()
	errOut.Reset()
	if rc := runWrapper([]string{"daneel", "--interactive", "--prompt", "--episode", "my-ep"}, &out, &errOut, sessionEnv{}); rc != exitOK {
		t.Fatalf("interactive run failed rc=%d stderr=%s", rc, errOut.String())
	}
	wantInteractive := []string{"--preset", "daneel", "--launch", "my-ep"}
	if !reflect.DeepEqual(gotArgs, wantInteractive) {
		t.Fatalf("interactive args = %q, want %q", gotArgs, wantInteractive)
	}

	// 3. Refusals: missing preset, missing task for headless, both --task and inline text.
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"no preset", nil, "<preset> is required"},
		{"no task headless", []string{"world-agy"}, "headless run requires a task"},
		{"both task and inline", []string{"world-agy", "--task", "/tmp/x.md", "inline"}, "cannot pass both --task"},
	} {
		out.Reset()
		errOut.Reset()
		if rc := runWrapper(tc.args, &out, &errOut, sessionEnv{}); rc == exitOK {
			t.Errorf("%s: expected refusal", tc.name)
		} else if !strings.Contains(errOut.String(), tc.want) {
			t.Errorf("%s: stderr %q does not contain %q", tc.name, errOut.String(), tc.want)
		}
	}
}
