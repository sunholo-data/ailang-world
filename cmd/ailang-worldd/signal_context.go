package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// serveSignalContext owns the daemon lifetime, separate from its startup budget.
// This file has no Store call.
func serveSignalContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
