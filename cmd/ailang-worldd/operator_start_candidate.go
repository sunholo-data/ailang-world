package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

// withOperatorStartCandidate owns the episode descriptor for the entire
// visit. The runner may wire a writer here only once the process and recovery
// gates are met. The visit must use a fresh clock and the locked row for CAS.
func withOperatorStartCandidate(ctx context.Context, credentials string, now func() int64,
	acquire func(string) (*os.File, error), visit func(*os.File, store.OperatorActionRow) error) error {
	canonical, err := store.CanonicalPath(credentials)
	if err != nil {
		return err
	}
	db, err := store.OpenCredentialsReadOnly(runnerDB(canonical))
	if err != nil {
		return err
	}
	defer db.Close()
	initial, err := db.NextOperatorStart(ctx, now())
	if err != nil {
		return err
	}
	if !verbs.EpisodeIDPattern.MatchString(initial.Episode) {
		return store.ErrOperatorConflict
	}
	lock, err := acquire(operatorEpisodeLockPath(canonical, initial.Episode))
	if err != nil {
		return fmt.Errorf("%w: episode lock unavailable: %w", store.ErrOperatorConflict, err)
	}
	defer lock.Close()
	current, err := db.NextOperatorStart(ctx, now())
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(current, initial) {
		return store.ErrOperatorConflict
	}
	return visit(lock, current)
}
