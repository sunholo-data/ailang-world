package store

import (
	"database/sql"
	"testing"
)

func TestOperatorMintStoreFilesFence(t *testing.T) {
	for _, name := range []string{"read", "abort", "abort_without_files"} {
		t.Run(name, func(t *testing.T) {
			f := newOperatorPinFixture(t, false)
			f.mint(t)
			id := f.in.Attempt.CAS.ID
			reviewed, err := f.c.ReadOperatorMint(f.ctx, id)
			if err != nil || reviewed.CredentialID != f.in.CredentialID {
				t.Fatalf("read before later reservation: %v", err)
			}
			if name == "abort_without_files" {
				if err := f.c.AbortOperatorMint(f.ctx, reviewed, pinNow); err != nil {
					t.Fatal(err)
				}
				after := f.snapshot(t)
				if after.action.State != "failed" || len(after.sessions) != 0 || after.start == nil {
					t.Fatalf("control failed to revoke and stop while retaining the start: %s, %d sessions, start present %t",
						after.action.State, len(after.sessions), after.start != nil)
				}
				return
			}
			files, _, err := f.c.ReadOperatorFiles(f.ctx, id)
			if err != nil {
				t.Fatal(err)
			}
			if err := f.c.RecordOperatorFiles(f.ctx, files, "files-identity", f.proof, func() int64 { return pinNow }); err != nil {
				t.Fatal(err)
			}
			key := files.Step
			key.Identity, key.Version = "files-identity", 1
			assertReservedOperatorStep(t, f.c, key)
			before := f.snapshot(t)
			if name == "abort" {
				err := f.c.AbortOperatorMint(f.ctx, reviewed, pinNow)
				f.assertRefused(t, before, err, sql.ErrNoRows)
				return
			}
			got, err := f.c.ReadOperatorMint(f.ctx, id)
			f.assertRefused(t, before, err, sql.ErrNoRows)
			if got != (OperatorMintReceipt{}) {
				t.Fatalf("refusal returned mint evidence: %+v", got)
			}
		})
	}
}
