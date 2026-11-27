package surfacepage

// PendingOperatorAction is one Needs-you row. Only an anchored, same-origin
// card route can appear in Href.
type PendingOperatorAction struct {
	ID        string `json:"id"`
	Episode   string `json:"episode"`
	Kind      string `json:"kind"`
	ExpiresAt int64  `json:"expiresAt"`
	Href      string `json:"href,omitempty"`
}
