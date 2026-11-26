package authority

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorPreDigestSignedStart(t *testing.T) {
	for _, modern := range []bool{false, true} {
		digest := hashref.SumSHA256([]byte("fixture executable")).String()
		r := newOperatorStoreRigAction(t, false, func(a *approval.OperatorAction) {
			if modern {
				a.Resolution = append(a.Resolution, approval.OperatorResolved{Name: "launch_binary_sha256", Value: digest})
			}
		})
		before := r.confirm(t)
		ctx := boundedTestContext(t)
		action, err := approval.DecodeOperatorAction(before.Action)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := RecheckClaimedOperatorAction(ctx, r.Credentials, action, [32]byte(before.Challenge), before.ExpiresAt, 1011, before); err != nil {
			t.Fatalf("old resolution must still verify: %v", err)
		}
		got, err := OperatorStartLaunchDigest(action)
		var stale *approval.BindingError
		if modern && (err != nil || got != digest) {
			t.Fatalf("modern digest = %q %v", got, err)
		}
		if !modern && (!errors.As(err, &stale) || stale.Field != "launch_binary_sha256" || got != "") {
			t.Fatalf("pre-digest signature needs stale binding: %q %v", got, err)
		}
		after, err := r.OperatorAction(ctx, before.ID, 1011)
		sessions, listErr := r.ListSessions(ctx, action.Episode, 1)
		start, startErr := r.EpisodeStart(ctx, action.Episode)
		_, stepErr := r.OperatorStep(ctx, store.OperatorStepKey{Action: r.cas, Kind: "mint", Identity: before.ID})
		if err != nil || listErr != nil || startErr != nil || start != nil || len(sessions) != 0 ||
			!errors.Is(stepErr, store.ErrOperatorMissing) || !reflect.DeepEqual(before, after) {
			t.Fatalf("digest admission changed the signed row or created effects: %v %v %v %v", err, listErr, startErr, stepErr)
		}
	}
}
