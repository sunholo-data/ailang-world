package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestAllowanceListAndRevokeCLI(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "credentials.db")
	st, err := store.OpenCredentials(credPath)
	if err != nil {
		t.Fatalf("OpenCredentials: %v", err)
	}
	const now int64 = 1_776_000_000
	al := approval.StandingAllowance{
		Schema:       approval.AllowanceSchema,
		PresetName:   "world-agy",
		PresetSHA256: strings.Repeat("a", 64),
		Profile:      approval.ProfileRef{Project: "ailang-world", Digest: "sha256:" + strings.Repeat("b", 64)},
		Repos:        []string{"/repo/ailang-world"},
		Grants: []approval.Grant{
			{Effect: "Workspace.Read", Scope: "*", Budget: 50},
		},
		TTLSeconds:         7200,
		MaxPerDay:          10,
		ValidDays:          30,
		WorldDB:            "/state/world.db",
		CredentialsDB:      "/state/credentials.db",
		WorkspaceRoot:      "/state/workspaces",
		ExcludeMissionLoop: true,
	}
	payload, err := approval.CanonicalAllowance(al)
	if err != nil {
		t.Fatalf("CanonicalAllowance: %v", err)
	}
	digest, err := approval.AllowanceDigest(al)
	if err != nil {
		t.Fatalf("AllowanceDigest: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := st.GrantAllowance(ctx, store.StandingAllowanceRow{
		AllowanceID:    "sa-1111222233334444",
		PresetName:     "world-agy",
		PresetSHA256:   strings.Repeat("a", 64),
		ProfileProject: "ailang-world",
		ProfileDigest:  "sha256:" + strings.Repeat("b", 64),
		Payload:        payload,
		Digest:         hex.EncodeToString(digest[:]),
		Challenge:      strings.Repeat("c", 64),
		MaxPerDay:      10,
		TTLS:           7200,
		ApprovedAt:     now - 10,
		ExpiresAt:      (now - 10) + 30*86400,
		ApprovedBy:     "tty",
	}); err != nil {
		t.Fatalf("GrantAllowance: %v", err)
	}
	_ = st.Close()

	env := sessionEnv{now: func() int64 { return now }}

	// 1. allowance list --json returns the active allowance.
	var out, errOut bytes.Buffer
	if code := runAllowance([]string{"list", "--json", "--credentials-db", credPath}, &out, &errOut, env); code != exitOK {
		t.Fatalf("allowance list exit=%d stderr=%s", code, errOut.String())
	}
	var items []allowanceListItem
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("unmarshal list: %v (%s)", err, out.String())
	}
	if len(items) != 1 || items[0].AllowanceID != "sa-1111222233334444" || items[0].Status != "active" || items[0].DrawsToday != 0 {
		t.Fatalf("unexpected list items: %+v", items)
	}

	// 2. allowance revoke by preset name revokes it.
	out.Reset()
	errOut.Reset()
	if code := runAllowance([]string{"revoke", "world-agy", "--reason", "rotating", "--credentials-db", credPath}, &out, &errOut, env); code != exitOK {
		t.Fatalf("allowance revoke exit=%d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "revoked standing allowance sa-1111222233334444") {
		t.Fatalf("unexpected revoke output: %s", out.String())
	}

	// 3. allowance list (default active only) is now empty; --all shows revoked.
	out.Reset()
	errOut.Reset()
	if code := runAllowance([]string{"list", "--json", "--credentials-db", credPath}, &out, &errOut, env); code != exitOK {
		t.Fatalf("allowance list after revoke exit=%d stderr=%s", code, errOut.String())
	}
	items = nil
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("unmarshal empty list: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 active allowances after revoke, got %d", len(items))
	}

	out.Reset()
	errOut.Reset()
	if code := runAllowance([]string{"list", "--all", "--json", "--credentials-db", credPath}, &out, &errOut, env); code != exitOK {
		t.Fatalf("allowance list --all exit=%d stderr=%s", code, errOut.String())
	}
	items = nil
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("unmarshal --all list: %v", err)
	}
	if len(items) != 1 || items[0].Status != "revoked" || items[0].RevokedReason != "rotating" {
		t.Fatalf("unexpected --all list items: %+v", items)
	}

	// 4. --from-allowance without --preset is refused.
	out.Reset()
	errOut.Reset()
	if code := runEpisodeStart([]string{"--from-allowance", "ep-1"}, &out, &errOut, env); code == exitOK {
		t.Fatalf("expected --from-allowance without --preset to be refused")
	}
	if !strings.Contains(errOut.String(), "--from-allowance requires a launch preset") {
		t.Fatalf("unexpected error for --from-allowance without --preset: %s", errOut.String())
	}
}
