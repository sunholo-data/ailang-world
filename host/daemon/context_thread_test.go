package daemon

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"testing"
)

func TestStartupCallerCancellation(t *testing.T) {
	for _, viaRun := range []bool{false, true} {
		t.Run(map[bool]string{false: "New", true: "Run"}[viaRun], func(t *testing.T) {
			cfg := Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			var err error
			if viaRun {
				err = Run(ctx, cfg, io.Discard)
			} else {
				var d *Daemon
				d, err = New(ctx, cfg)
				if d != nil {
					d.Close()
				}
			}
			var se *StartupError
			if !errors.Is(err, context.Canceled) || !errors.As(err, &se) || se.Stage != StageRegistry {
				t.Fatalf("startup must observe cancelled ctx at registry: %v", err)
			}
			d, err := New(context.Background(), cfg)
			if err != nil {
				t.Fatalf("live control/released writer: %v", err)
			}
			d.Close()
		})
	}
}
