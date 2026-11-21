//go:build !darwin && !linux

package daemon

import "os"

// A verified read still works when change time is unavailable for cache invalidation.
func surfaceUsageChangeTime(os.FileInfo) (int64, bool) { return 0, false }
