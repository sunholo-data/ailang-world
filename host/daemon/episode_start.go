package daemon

import (
	"context"
	"errors"
	"os"

	"github.com/sunholo-data/ailang-world/host/store"
)

// Each projection opens only a read-only credentials handle. No session route
// has a writer for attended start history.
func (d *Daemon) readEpisodeStart(ctx context.Context, episode string) (*store.EpisodeStart, error) {
	if d.credentialsPath == "" {
		return nil, nil
	}
	c, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = c.Close() }()
	return c.EpisodeStart(ctx, episode)
}
