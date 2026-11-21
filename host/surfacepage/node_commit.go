package surfacepage

// NodeCommit is a recorded Git message on the selected node.
type NodeCommit struct {
	ID       string   `json:"id"`
	Kinds    []string `json:"kinds"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body,omitempty"`
	Redacted bool     `json:"redacted,omitempty"`
}
