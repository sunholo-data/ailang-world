package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorServiceUpgrade(t *testing.T) {
	t.Run("active_child", testOperatorUpgradeActiveChild)
	t.Run("stopped_companion", testOperatorUpgradeStoppedCompanion)
	t.Run("rollback_schema_drift", testOperatorUpgradeRollbackSchema)
	for _, tc := range []struct {
		name         string
		before, next int
		tamper       string
	}{
		{"old reader", 9, 10, ""},
		{"missing previous schema", 0, 10, ""},
		{"missing candidate schema", 10, 0, ""},
		{"substituted previous binary", 10, 10, "different bytes"},
		{"linked previous binary", 10, 10, "link"},
		{"missing previous binary", 10, 10, "missing"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, candidate, _ := upgradeFixture(t)
			original := serviceFileHash(f.env.paths.Binary)
			r, err := readServiceRecord(f.env.paths)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := store.UpgradeCredentials(r.Credentials); err != nil {
				t.Fatal(err)
			}
			db, record := serviceFileHash(r.Credentials), serviceFileHash(f.env.paths.Record)
			probe := f.env.probe
			f.env.probe = func(path, home string) ([]byte, error) {
				b, err := probe(path, home)
				if err != nil {
					return nil, err
				}
				var v serviceVersion
				if err := json.Unmarshal(b, &v); err != nil {
					return nil, err
				}
				v.CredentialsSchema = tc.before
				if path == candidate {
					v.CredentialsSchema = tc.next
				}
				return json.Marshal(v)
			}
			f.ctl.onKick = func() error {
				if tc.tamper != "" {
					prev := f.env.paths.Binary + ".prev"
					if err := os.Remove(prev); err != nil {
						t.Fatal(err)
					}
					switch tc.tamper {
					case "link":
						if err := os.Symlink(f.source, prev); err != nil {
							t.Fatal(err)
						}
					case "different bytes":
						if err := os.WriteFile(prev, []byte("substitution"), 0700); err != nil {
							t.Fatal(err)
						}
					}
				}
				return fmt.Errorf("fixture restart failed")
			}
			code, msg := f.run("upgrade", "--binary", candidate)
			if code != 3 || !strings.Contains(msg, "automatic rollback refused") ||
				!strings.Contains(msg, "verified schema-compatible binary") || f.ctl.kicks != 1 || f.ctl.stops != 0 ||
				serviceFileHash(f.env.paths.Binary) != serviceFileHash(candidate) ||
				serviceFileHash(f.env.paths.Record) != record || serviceFileHash(r.Credentials) != db ||
				serviceFileHash(f.env.paths.Binary+".failed") != "" ||
				(tc.tamper == "" && serviceFileHash(f.env.paths.Binary+".prev") != original) {
				t.Fatalf("upgrade %d, kicks %d, stops %d: %s", code, f.ctl.kicks, f.ctl.stops, msg)
			}
			if strings.Contains(msg, wantStoppedRunnerRemedy) {
				t.Fatalf("no-companion install received a companion restart instruction: %s", msg)
			}
		})
	}
}
