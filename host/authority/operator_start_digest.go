package authority

import (
	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

// OperatorStartLaunchDigest reads the immutable signed binding. Earlier
// envelopes require a new review before an executable may be selected.
func OperatorStartLaunchDigest(a approval.OperatorAction) (string, error) {
	for _, fact := range a.Resolution {
		if fact.Name == "launch_binary_sha256" {
			hash, err := hashref.Parse(fact.Value)
			if err == nil && hash.Algo() == "sha256" {
				return fact.Value, nil
			}
			break
		}
	}
	return "", &approval.BindingError{Field: "launch_binary_sha256"}
}
