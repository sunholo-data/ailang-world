package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOperatorActivationDrainRefusals(t *testing.T) {
	t.Run("from_binary", func(t *testing.T) {
		f, _, candidate := pairedUpgradeFixture(t)
		before := pairedUpgradeEvidence(f)
		hash := serviceFileHash(candidate)
		if hash == before[0] || hash == "" {
			t.Fatal("need distinct verified build hashes")
		}
		pair, err := f.env.reviewRunnerUpgrade(f.env.runnerUpgradeReview(), hash, hash)
		if pair != nil || err == nil || !strings.Contains(err.Error(), "changed during confirmation") ||
			pairedUpgradeEvidence(f) != before {
			t.Fatalf("approved pair crossed the from-build binding: %v", err)
		}
	})
	t.Run("budget_expired", func(t *testing.T) {
		f, s, candidate := pairedUpgradeFixture(t)
		before := pairedUpgradeEvidence(f)
		start := time.Unix(1000, 0)
		elapsed := false
		f.env.runnerNow = func() time.Time {
			if elapsed {
				return start.Add(f.env.wait)
			}
			elapsed = true
			return start
		}
		candidateBefore := serviceFileHash(candidate)
		guard, err := f.env.drainRunner(f.config)
		if guard != nil {
			guard.Close()
		}
		_, lockErr := os.Lstat(f.config.Credentials + ".operator-runner.lock")
		if guard != nil || err == nil || !strings.Contains(err.Error(), "companion drain budget expired") || s.runnerLoaded ||
			f.homeJob.kicks != 0 || pairedUpgradeEvidence(f) != before || serviceFileHash(candidate) != candidateBefore ||
			!errors.Is(lockErr, os.ErrNotExist) {
			t.Fatalf("expired drain acquired the lock or swapped: %v (%v)", err, lockErr)
		}
	})
	t.Run("stop_canonical", func(t *testing.T) {
		f, s, _ := pairedUpgradeFixture(t)
		alias := filepath.Join(f.home, "credentials-alias")
		if err := os.Symlink(f.config.Credentials, alias); err != nil {
			t.Fatal(err)
		}
		got, gotErr := filepath.EvalSymlinks(alias)
		want, wantErr := filepath.EvalSymlinks(f.config.Credentials)
		if gotErr != nil || wantErr != nil || got != want || got == alias {
			t.Fatal("need a distinct path to the same credentials")
		}
		r, err := readServiceRecord(f.env.paths)
		if err != nil {
			t.Fatal(err)
		}
		r.Credentials = alias
		b, err := json.Marshal(r)
		if err != nil {
			t.Fatal(err)
		}
		if err := f.env.atomic(f.env.paths.Record, b, 0600); err != nil {
			t.Fatal(err)
		}
		before := pairedUpgradeEvidence(f)
		code, msg := f.command("stop")
		if code != 1 || !strings.Contains(msg, "canonical credentials lock") || !s.runnerLoaded ||
			pairedUpgradeEvidence(f) != before {
			t.Fatalf("stop accepted an alias in the service record: %d %s", code, msg)
		}
	})
	t.Run("stopped_failure", func(t *testing.T) {
		f, s, candidate := pairedUpgradeFixture(t)
		s.runnerLoaded = false
		before := pairedUpgradeEvidence(f)
		read, calls := f.env.runnerProcesses, 0
		f.env.runnerProcesses = func() ([]worldProcess, error) {
			calls++
			if calls == 2 {
				return nil, errors.New("fixture post-drain census denied")
			}
			return read()
		}
		code, msg := f.run("upgrade", "--binary", candidate)
		if code != 1 || !strings.Contains(msg, "terminal drain refused") || strings.Contains(msg, wantStoppedRunnerRemedy) ||
			s.runnerLoaded || f.homeJob.kicks != 0 || pairedUpgradeEvidence(f) != before {
			t.Fatalf("pre-stopped companion failure: %d %s", code, msg)
		}
	})
}

func TestOperatorRevokeStatRefusal(t *testing.T) {
	f := operatorService(t)
	path := filepath.Join(f.home, "stat-loop")
	if err := os.Symlink(path, path); err != nil {
		t.Fatal(err)
	}
	before := pairedUpgradeEvidence(f)
	ctx := activationTestContext(t)
	id := strings.Repeat("a", 64)
	if _, ok, err := readSessionRevokeTarget(ctx, f.config.World, path+".absent", id); ok || err != nil {
		t.Fatalf("missing credentials control: %t %v", ok, err)
	}
	_, ok, err := readSessionRevokeTarget(ctx, f.config.World, path, id)
	target, linkErr := os.Readlink(path)
	if ok || err == nil || errors.Is(err, os.ErrNotExist) || target != path || linkErr != nil || pairedUpgradeEvidence(f) != before {
		t.Fatalf("stat failure was treated as missing credentials: %t %v", ok, err)
	}
}
