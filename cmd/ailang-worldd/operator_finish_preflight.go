package main

import (
	"bytes"
	"context"
	"crypto/sha256"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

// operatorFinishPrep binds a current durable claim to independently resolved
// inputs. Only the proof returned by the runner's recheck can seal it.
type operatorFinishPrep struct {
	action    approval.OperatorAction
	row       store.OperatorActionRow
	authority *authority.VerifiedOperatorAction
}

func verifyOperatorFinishPrep(ctx context.Context, ro *store.Credentials, frozen approval.OperatorAction, now int64,
	cas store.OperatorCAS, owner string) (operatorFinishPrep, error) {
	if frozen.Kind != "finish" || !frozen.Options.ConfirmOnSurface || frozen.ID != cas.ID {
		return operatorFinishPrep{}, approval.ErrNotVerified
	}
	row, err := ro.OperatorAction(ctx, cas.ID, now)
	if err != nil {
		return operatorFinishPrep{}, err
	}
	phase := row.State == "confirmed"
	if owner != "" {
		phase = row.State == "claimed" || row.State == "running"
	}
	if !phase || row.Owner != owner || row.Generation != cas.Generation || row.Kind != frozen.Kind ||
		row.Episode != frozen.Episode || len(row.Challenge) != 32 || !bytes.Equal(row.Digest, cas.Digest[:]) {
		return operatorFinishPrep{}, approval.ErrNotVerified
	}
	if now < row.CreatedAt || now >= row.ExpiresAt {
		return operatorFinishPrep{}, store.ErrOperatorExpired
	}
	raw, err := approval.CanonicalOperatorAction(frozen)
	if err != nil || sha256.Sum256(raw) != cas.Digest || !bytes.Equal(raw, row.Action) {
		return operatorFinishPrep{}, approval.ErrNotVerified
	}
	proof, err := authority.RecheckClaimedOperatorAction(ctx, ro, frozen, [32]byte(row.Challenge), row.ExpiresAt, now, row)
	if err != nil {
		return operatorFinishPrep{}, approval.ErrNotVerified
	}
	return operatorFinishPrep{frozen, row, proof}, nil
}

func (p operatorFinishPrep) claimed() bool {
	return p.row.Owner != "" && len(p.row.Challenge) == 32 && p.authority.Rechecked() &&
		p.authority.Binds(p.action, [32]byte(p.row.Challenge), p.row.ExpiresAt) == nil
}
