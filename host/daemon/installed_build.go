package daemon

import (
	"os"
	"sync"
	"time"
)

// Keep one entry: production reads its own path. Tests can exercise other paths.
var installedBuildCache struct {
	sync.Mutex
	path   string
	size   int64
	mtime  time.Time
	digest string
}

// installedDigestAt reads the path's current bytes, independently of CurrentBuild.
func installedDigestAt(path string) string {
	installedBuildCache.Lock()
	defer installedBuildCache.Unlock()
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}
	c := &installedBuildCache
	if c.path == path && c.size == info.Size() && c.mtime.Equal(info.ModTime()) {
		return c.digest
	}
	digest := executableDigestAt(path)
	c.path, c.size, c.mtime, c.digest = path, info.Size(), info.ModTime(), digest
	return digest
}
