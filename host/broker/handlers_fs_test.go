package broker

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFSHandlerReadWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.txt")
	h := FSHandler{}
	writeReq := EffectRequest{Effect: EffectFSWrite, Scope: path}
	if _, err := h.Execute(context.Background(), writeReq, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	readReq := EffectRequest{Effect: EffectFSRead, Scope: path}
	got, err := h.Execute(context.Background(), readReq, nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "hello" {
		t.Fatalf("FS.Read = %q, want hello", got)
	}
}

func TestFSHandlerRefusesSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	scopeParent := filepath.Join(root, "scope")
	outside := filepath.Join(root, "outside")
	if err := os.MkdirAll(scopeParent, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outside, 0o700); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(outside, "secret.txt")
	if err := os.WriteFile(target, []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(scopeParent, "link.txt")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	_, err := (FSHandler{}).Execute(
		context.Background(),
		EffectRequest{Effect: EffectFSRead, Scope: link},
		nil,
	)
	var pathErr *FSPathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("FS.Read error = %v, want FSPathError", err)
	}
}

func TestFSHandlerRequiresCanonicalAbsoluteScope(t *testing.T) {
	_, err := (FSHandler{}).Execute(
		context.Background(),
		EffectRequest{Effect: EffectFSRead, Scope: "relative.txt"},
		nil,
	)
	var pathErr *FSPathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("FS.Read error = %v, want FSPathError", err)
	}
}
