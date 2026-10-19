package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

// MaxInvocationIntentsPage is the kernel-owned upper bound for one
// InvocationIntents page.
const MaxInvocationIntentsPage = 1000

// InvocationIntent is one invocation intent and, once its commit landed, the
// outcome the same commit transaction appended (Commit step 6). Outcome nil
// means the invocation never committed (or has not yet).
type InvocationIntent struct {
	Seq       int64
	ObjectRef hashref.HashRef
	Intent    JournalIntent
	Outcome   *JournalOutcome
}

// InvalidIntentPrefixError refuses a prefix InvocationIntents cannot bound.
type InvalidIntentPrefixError struct{ Prefix string }

func (e *InvalidIntentPrefixError) Error() string {
	return fmt.Sprintf("store: invocation intent prefix %q must end with \":\" and be outside the effect: namespace", e.Prefix)
}

// invocationIntentsSQL reads one key range of the journal's UNIQUE
// (invocation_id, kind) index — the access path EffectSpend uses — so its
// cost is the range's rows, never the whole log or journal.
const invocationIntentsSQL = `SELECT j.seq, j.object_ref, o.payload, oc.payload
FROM journal j JOIN objects o ON o.hash_ref = j.object_ref
LEFT JOIN journal x ON x.invocation_id = j.invocation_id AND x.kind = 'outcome'
LEFT JOIN objects oc ON oc.hash_ref = x.object_ref
WHERE j.kind = 'intent' AND j.invocation_id >= ? AND j.invocation_id < ?
AND j.seq > ? AND o.semantic_id = ?
ORDER BY j.seq LIMIT ?`

// InvocationIntents returns up to limit invocation intents whose id starts
// with prefix and whose journal seq is after afterSeq, ascending by seq, each
// with its outcome when one exists (row 198: one episode's calls are the
// coordinator ids "<surface>:<episode>:<task>"). prefix must end with ":", so
// the range [prefix, prefix with ";" for its last byte) holds exactly the ids
// that start with it; callers that need an exact episode still check the
// record, since an episode id may itself contain ":".
func (s *Store) InvocationIntents(ctx context.Context, prefix string, afterSeq int64, limit int) ([]InvocationIntent, error) {
	if err := s.checkQuarantine(); err != nil {
		return nil, err
	}
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	const op = "InvocationIntents"
	if limit < 1 || limit > MaxInvocationIntentsPage {
		return nil, &InvalidLimitError{Op: op, Limit: limit, Max: MaxInvocationIntentsPage}
	}
	if !strings.HasSuffix(prefix, ":") || strings.HasPrefix(prefix, "effect:") {
		return nil, &InvalidIntentPrefixError{Prefix: prefix}
	}
	upper := strings.TrimSuffix(prefix, ":") + ";"
	rows, err := s.db.QueryContext(ctx, invocationIntentsSQL, prefix, upper, afterSeq, JournalIntentV1, limit)
	if err != nil {
		return nil, fmt.Errorf("store: invocation intents %q: %w", prefix, err)
	}
	defer rows.Close()
	var out []InvocationIntent
	for rows.Next() {
		var item InvocationIntent
		var refText string
		var intentPayload, outcomePayload []byte
		if err := rows.Scan(&item.Seq, &refText, &intentPayload, &outcomePayload); err != nil {
			return nil, fmt.Errorf("store: invocation intents %q: scan: %w", prefix, err)
		}
		if item.ObjectRef, err = hashref.Parse(refText); err != nil {
			return nil, fmt.Errorf("store: invocation intent %q ref: %w", refText, err)
		}
		if item.Intent, err = decodeJournalIntent(intentPayload); err != nil {
			return nil, fmt.Errorf("store: decode intent %s: %w", refText, err)
		}
		if outcomePayload != nil {
			outcome, err := decodeJournalOutcome(outcomePayload)
			if err != nil {
				return nil, fmt.Errorf("store: decode outcome of %q: %w", item.Intent.InvocationID, err)
			}
			item.Outcome = &outcome
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: invocation intents %q: %w", prefix, err)
	}
	return out, nil
}
