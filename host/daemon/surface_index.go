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
	empty := func(reason string) {
		// These fixed reasons contain no filesystem path.
		_, _ = io.WriteString(w, `{"schema":"`+surface.Schema+`","empty":true,"reason":"`+reason+`"}`)
	}
	refuse := func(err error) {
		// os.PathError includes its path. Unwrap it before reporting detail,
		// and redact a configured path echoed by a decoder or validator.
		var pe *os.PathError
		if errors.As(err, &pe) {
			err = pe.Err
		}
		reason := strings.ReplaceAll(err.Error(), d.cfg.SurfaceIndex, "the configured surface index")
		writeAPIError(w, "SurfaceIndexRefused", "the configured surface index: "+reason, http.StatusServiceUnavailable)
	}
	if d.cfg.SurfaceIndex == "" {
		empty("no surface index is configured: start serve with --surface-index <file>")
		return
	}
	f, err := os.Open(d.cfg.SurfaceIndex)
	if errors.Is(err, os.ErrNotExist) {
		empty("the surface index has not been built yet: run ailang-worldd surface-index")
		return
	}
	if err != nil {
		refuse(err)
		return
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, surface.MaxBytes+1))
	if err != nil {
		refuse(err)
		return
	}
	if len(data) > surface.MaxBytes {
		refuse(fmt.Errorf("over the 4 MiB (%d bytes) bound", surface.MaxBytes))
		return
	}
	var ix surface.Index
	if err := json.Unmarshal(data, &ix); err != nil {
		refuse(err)
		return
	}
	if err := surface.Validate(ix); err != nil {
		refuse(err)
		return
	}
	_, _ = w.Write(data)
}
