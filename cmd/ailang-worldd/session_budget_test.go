package main

import (
	"os"
	"strings"
	"testing"
)

func TestSessionRootBudgets(t *testing.T) {
	b, err := os.ReadFile("session.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(b)
	for _, spec := range []struct{ name, root, call string }{
		{"mint", "mintCtx, cancelMint := context.WithTimeout(context.Background(), 3*time.Second)", "authority.Mint(mintCtx, st,"},
		{"revoke", "revokeCtx, cancelRevoke := context.WithTimeout(context.Background(), 3*time.Second)", "authority.Revoke(revokeCtx, st,"},
	} {
		a := strings.Index(src, spec.root)
		b := strings.Index(src, spec.call)
		if a < 0 || b < a {
			t.Errorf("%s lacks bounded authority call", spec.name)
		}
	}
}
