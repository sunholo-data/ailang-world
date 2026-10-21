package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/surface"
)

// Check both the file and its directory, including symlink targets, before
// taking writer authority. The builder may create the file after startup.
func checkSurfaceIndexOutsideRoot(path, root string) error {
	for _, p := range []string{filepath.Dir(path), path} {
		if err := broker.CheckPolicyOutsideRoot(p, root); err != nil {
			return fmt.Errorf("the configured surface index must be outside --workspace-root: %w", err)
		}
	}
	return nil
}

func (d *Daemon) handleSurfaceIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	data, _, empty, err := d.readSurfaceIndex()
	switch {
	case err != nil:
		writeAPIError(w, "SurfaceIndexRefused", d.surfaceIndexRefusal(err), http.StatusServiceUnavailable)
	case empty != "":
		// These fixed reasons contain no filesystem path.
		_, _ = io.WriteString(w, `{"schema":"`+surface.Schema+`","empty":true,"reason":"`+empty+`"}`)
	default:
		_, _ = w.Write(data)
	}
}

// readSurfaceIndex is the one read of the configured index that both
// surface routes share: re-read on every request, bounded by
// surface.MaxBytes, decoded and validated. empty is a positive empty state's
// fixed reason (no index configured, or not built yet); err is a refusal, to
// be reported through surfaceIndexRefusal.
func (d *Daemon) readSurfaceIndex() (data []byte, ix surface.Index, empty string, err error) {
	if d.cfg.SurfaceIndex == "" {
		return nil, ix, "no surface index is configured: start serve with --surface-index <file>", nil
	}
	f, err := os.Open(d.cfg.SurfaceIndex)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ix, "the surface index has not been built yet: run ailang-worldd surface-index", nil
	}
	if err != nil {
		return nil, ix, "", err
	}
	defer f.Close()
	data, err = io.ReadAll(io.LimitReader(f, surface.MaxBytes+1))
	if err != nil {
		return nil, ix, "", err
	}
	if len(data) > surface.MaxBytes {
		return nil, ix, "", fmt.Errorf("over the 4 MiB (%d bytes) bound", surface.MaxBytes)
	}
	if err := json.Unmarshal(data, &ix); err != nil {
		return nil, surface.Index{}, "", err
	}
	if err := surface.Validate(ix); err != nil {
		return nil, surface.Index{}, "", err
	}
	return data, ix, "", nil
}

// surfaceIndexRefusal is a refusal's text, without the configured path.
func (d *Daemon) surfaceIndexRefusal(err error) string {
	// os.PathError includes its path. Unwrap it before reporting detail,
	// and redact a configured path echoed by a decoder or validator.
	var pe *os.PathError
	if errors.As(err, &pe) {
		err = pe.Err
	}
	reason := strings.ReplaceAll(err.Error(), d.cfg.SurfaceIndex, "the configured surface index")
	return "the configured surface index: " + reason
}
