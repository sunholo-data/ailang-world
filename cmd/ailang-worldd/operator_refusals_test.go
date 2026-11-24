package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestOperatorRunnerLockFileRefusals(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lock")
	if err := os.WriteFile(path, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	before := serviceFileHash(path)
	alias := path + "-alias"
	if err := os.Symlink(path, alias); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, path string
		mode       os.FileMode
	}{
		{"symlink", alias, 0600},
		{"mode", path, 0644},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.Chmod(path, tc.mode); err != nil {
				t.Fatal(err)
			}
			lock, err := runnerKernelLock(tc.path)
			if lock != nil {
				lock.Close()
			}
			info, statErr := os.Stat(path)
			target, linkErr := os.Readlink(alias)
			if err == nil || lock != nil || statErr != nil || info.Mode().Perm() != tc.mode ||
				serviceFileHash(path) != before || linkErr != nil || target != path {
				t.Fatalf("unsafe lock did not refuse without mutation: %v %v", lock, err)
			}
		})
	}
}

func TestOperatorRunnerSetupRefusals(t *testing.T) {
	t.Run("linked config", func(t *testing.T) {
		f := operatorService(t)
		f.mustInstall(t)
		p := f.env.runnerPaths()
		alias := f.configFile + "-alias"
		if err := os.Symlink(f.configFile, alias); err != nil {
			t.Fatal(err)
		}
		before := []string{serviceFileHash(f.configFile), serviceFileHash(p.Record), serviceFileHash(p.Plist), serviceFileHash(f.config.Credentials)}
		f.writes = nil
		code, msg := f.command("install", "--config", alias)
		f.refused(t, "invalid runner setup file", code, msg)
		after := []string{serviceFileHash(f.configFile), serviceFileHash(p.Record), serviceFileHash(p.Plist), serviceFileHash(f.config.Credentials)}
		if target, err := os.Readlink(alias); err != nil || target != f.configFile || !reflect.DeepEqual(after, before) {
			t.Fatal("refusal modified configuration or credentials")
		}
	})
	t.Run("reserved preset", func(t *testing.T) {
		f := operatorService(t)
		f.config.PresetNames = []string{"se-tools"}
		code, msg := f.install(t)
		f.refused(t, "invalid runner preset name", code, msg)
	})
	t.Run("global address", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)
		t.Setenv("AILANG_REGISTRY_API_KEY", "")
		var out, errOut bytes.Buffer
		if code := run([]string{"--addr", "127.0.0.1:9", "operator-runner"}, &out, &errOut); code != exitUsage || out.Len()+errOut.Len() != 0 {
			t.Fatalf("explicit address reached runner startup: %d %s", code, errOut.String())
		}
		if files, err := os.ReadDir(home); err != nil || len(files) != 0 {
			t.Fatal("refused invocation changed its home")
		}
	})
}
