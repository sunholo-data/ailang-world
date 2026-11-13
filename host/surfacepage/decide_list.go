package surfacepage

// Row 201 A4 sprint 4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decisions
// GET /surface/requests.json lists under the mint requests, and approve --list prints. Pure: the
// daemon reads the approvals chain and one bounded packet scan, computes each row's display verdict
// with the decision law, and this fold only selects and shapes the rows. Display only: the POST
// re-runs the law with the complete discovery.

// decisionEffect is the only effect a decision card decides.
const decisionEffect = "Human.Approve"

// PendingDecision is one line of requests.json's decisions: a pending Human.Approve request and the
// packet carrying it. Packet is "" and Deadline 0 unless exactly one scanned packet carries it;
// Reserved names a typed owner ("publish", "finish"); Refusal is the law's reason at now, "" when
// the request can be decided.
type PendingDecision struct {
	Request   string `json:"request"`
	Packet    string `json:"packet"`
	Effect    string `json:"effect"`
	Scope     string `json:"scope"`
	Cost      int64  `json:"cost"`
	Requester string `json:"requester"`
	Deadline  int64  `json:"deadline"`
	Reserved  string `json:"reserved"`
	Refusal   string `json:"refusal"`
}

// DecisionCandidate is one request of the approvals chain as the daemon read it: its summary, the
// refs of the scanned packets carrying it (scan order), and, for those it computed them for, the one
// packet's deadline, the reserved owner and the law's display reason.
type DecisionCandidate struct {
	Request, Effect, Scope, Requester string
	Cost                              int64
	Decided                           bool
	Packets                           []string
	Deadline                          int64
	Reserved, Refusal                 string
}

// FoldDecisions is the decisions list (design §2.8): pending Human.Approve requests only; a request
// no scanned packet carries is omitted; one carried by two or more is listed without a packet, with
// its refusal. truncated is set when the approvals read or the packet scan hit its bound.
func FoldDecisions(cands []DecisionCandidate, chainTruncated, scanTruncated bool) (rows []PendingDecision, truncated bool) {
	for _, c := range cands {
		if c.Decided || c.Effect != decisionEffect || len(c.Packets) == 0 {
			continue
		}
		row := PendingDecision{Request: c.Request, Effect: c.Effect, Scope: c.Scope, Cost: c.Cost, Requester: c.Requester,
			Reserved: c.Reserved, Refusal: c.Refusal}
		if len(c.Packets) == 1 {
			row.Packet, row.Deadline = c.Packets[0], c.Deadline
		}
		rows = append(rows, row)
	}
	return rows, chainTruncated || scanTruncated
}
