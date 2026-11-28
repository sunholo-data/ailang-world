package appmanifest

import (
	"reflect"
	"sort"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// TestAllowedEffectsAreBrokerEffects binds effectAllowed's set to the broker's own effect
// constants: an app may only grant an effect a broker handler serves, and a renamed
// constant turns this red instead of leaving apps naming an effect nothing handles.
func TestAllowedEffectsAreBrokerEffects(t *testing.T) {
	want := []string{
		broker.EffectWorkspaceRead, broker.EffectWorkspaceWrite, broker.EffectWorkspaceDelete,
		broker.EffectWorkspaceExec, broker.EffectAilangCheck, broker.EffectAilangRun,
		broker.EffectAilangRunEnv, broker.EffectAilangRunNet, broker.EffectAilangDiscover,
		broker.EffectAilangCLI,
	}
	var got []string
	for e := range allowedEffects {
		got = append(got, e)
	}
	sort.Strings(want)
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("allowed app effects %v, broker episode effects %v", got, want)
	}
	for _, refused := range []string{broker.EffectModelInfer, "Human.Approve", "FS.Write"} {
		if effectAllowed(refused) {
			t.Errorf("%s must not be grantable by an app", refused)
		}
	}
}
