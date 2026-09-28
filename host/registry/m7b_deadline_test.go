package registry

import (
	"context"
	"testing"
	"time"
)

// boundedTestContext gives test store calls a finite caller budget.
func boundedTestContext(t testing.TB) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}
