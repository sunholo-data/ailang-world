package main

import (
	"context"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/pkgproj"
)

// reconcileReadOnlyProbe performs only remote metadata reads; it makes no Store call.
func reconcileReadOnlyProbe(opts options, packet pkgproj.ReadyPacket) (broker.ReconcileReceipt, error) {
	return broker.ReconcileRegistryPublish(context.Background(), reconcileConfigFor(opts, packet))
}
