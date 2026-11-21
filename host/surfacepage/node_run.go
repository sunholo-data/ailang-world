package surfacepage

// NodeRun holds evidence for the selected episode only. Durations use seconds.
type NodeRun struct {
	Episode     string           `json:"episode"`
	Harness     string           `json:"harness"`
	Model       string           `json:"model,omitempty"`
	Launches    int              `json:"launches"`
	Runtime     *int64           `json:"runtime,omitempty"`
	Elapsed     *int64           `json:"elapsed,omitempty"`
	Turns       *int             `json:"turns,omitempty"`
	Tokens      *NodeTokens      `json:"tokens,omitempty"`
	CostUSD     *float64         `json:"costUsd,omitempty"`
	Decisions   []NodeDecision   `json:"decisions,omitempty"`
	Checks      []NodeCheck      `json:"checks,omitempty"`
	Permissions *NodePermissions `json:"permissions,omitempty"`
}

// NodeTokens holds totals normalized to the provider's input and output
// semantics. A nil field means the retained evidence did not report a total.
type NodeTokens struct {
	Input      *int64 `json:"input,omitempty"`
	Output     *int64 `json:"output,omitempty"`
	CacheRead  *int64 `json:"cacheRead,omitempty"`
	CacheWrite *int64 `json:"cacheWrite,omitempty"`
	Reasoning  *int64 `json:"reasoning,omitempty"`
}

type NodeDecision struct {
	Action string `json:"action"`
	At     string `json:"at"`
}

type NodeCheck struct {
	Name   string `json:"name"`
	Result string `json:"result"`
}

// NodePermissions separates retained credential grants from historical debit
// evidence. Neither view combines overlapping credential allowances.
type NodePermissions struct {
	Sessions  []NodePermissionSession `json:"sessions"`
	Observed  []NodeGrant             `json:"observed,omitempty"`
	Truncated bool                    `json:"truncated,omitempty"`
}

type NodePermissionSession struct {
	State     string      `json:"state"`
	Current   bool        `json:"current,omitempty"`
	CreatedAt string      `json:"createdAt,omitempty"`
	ExpiresAt string      `json:"expiresAt,omitempty"`
	Grants    []NodeGrant `json:"grants"`
}

type NodeGrant struct {
	Effect          string `json:"effect"`
	Scope           string `json:"scope"`
	ScopeKey        int    `json:"scopeKey"`
	State           string `json:"state,omitempty"`
	ExpiresAt       string `json:"expiresAt,omitempty"`
	Budget          *int64 `json:"budget,omitempty"`
	Used            *int64 `json:"used,omitempty"`
	Reserved        *int64 `json:"reserved,omitempty"`
	Remaining       *int64 `json:"remaining,omitempty"`
	LedgerRemaining *int64 `json:"ledgerRemaining,omitempty"`
}
