package appregistry

// Go mirrors of design_docs/sketches/appregistry.ail (same names, same bounds). Every
// pure decision the registry makes goes through one of these; drift_test.go parses the
// sketch and replays each of its inline rows through the mirror of the same name.

func revisionParentOk(revision int64, hasParent bool) bool {
	return revision >= 1 && ((revision == 1) == !hasParent)
}

func nextRevisionOk(captured, next int64) bool { return next == captured+1 }

func nextVersion(previous int64) int64 { return previous + 1 }

func versionOk(v int64) bool { return v >= 1 }

func publishChanges(samePinnedManifest bool) bool { return !samePinnedManifest }

func retireOk(present bool) bool { return present }

func appCountOk(n int) bool { return n >= 0 && n <= maxApps }

func transitionCountOk(n int) bool { return n >= 0 && n <= maxTransitions }

func manifestPinned(entriesNamingIt int) bool { return entriesNamingIt >= 1 }

// listingState and launchable are the M1b-3 listing laws.
const (
	StatePublished    = "published"
	StateLocalDiffers = "published-local-differs"
	StateLocal        = "local"
)

func listingState(published, local, sameBytes bool) string {
	if published {
		if local && !sameBytes {
			return StateLocalDiffers
		}
		return StatePublished
	}
	return StateLocal
}

// Launchable is D-WORLD-129 = A: only a published app (in either published state) may be
// launched, and it launches its pinned manifest.
func Launchable(state string) bool { return state == StatePublished || state == StateLocalDiffers }
