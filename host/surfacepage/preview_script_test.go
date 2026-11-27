package surfacepage

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// TestPreviewScriptMatchesEmbeddedBundle ensures scripts/surface/preview.mjs
// concatenates AppScript, AppStyle, and AppShell in the exact same order as
// app.go so live disk previews never drift from compiled binaries.
func TestPreviewScriptMatchesEmbeddedBundle(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node unavailable locally")
	}
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	preview := filepath.Join(root, "scripts", "surface", "preview.mjs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	check := `
import assert from 'node:assert/strict';
import { pathToFileURL } from 'node:url';
const mod = await import(pathToFileURL(process.argv[1]).href);
const bundle = mod.assembleBundle(process.argv[2]);
process.stdout.write(JSON.stringify({
  scriptBytes: Buffer.byteLength(bundle.script),
  styleBytes: Buffer.byteLength(bundle.style),
  shellBytes: Buffer.byteLength(bundle.shell),
}));
`
	out, err := exec.CommandContext(ctx, node, "--input-type=module", "-e", check, preview, root).CombinedOutput()
	if err != nil {
		t.Fatalf("preview.mjs failed: %v\n%s", err, out)
	}
	want := `{"scriptBytes":` + itoa(len(AppScript)) + `,"styleBytes":` + itoa(len(AppStyle)) + `,"shellBytes":` + itoa(len(AppShell)) + `}`
	if string(out) != want {
		t.Fatalf("preview bundle sizes = %s, want %s (preview.mjs drifted from app.go)", out, want)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
