package broker

import (
	"context"
	"fmt"
	"strings"
)

// protectedGitPreflight refuses config-defined host commands before status.
func protectedGitPreflight(ctx context.Context, git, root string) error {
	filters, err := configuredGitFilters(ctx, git, root)
	if err != nil {
		return fmt.Errorf("protected Git filter query: %w", err)
	}
	if len(filters) > 0 {
		return fmt.Errorf("%s: %s (git status would run a configured command)",
			WorktreeFilterRefused, strings.Join(filters, ", "))
	}
	return nil
}
