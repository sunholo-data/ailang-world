package daemon

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

const operatorPath = "/surface/operator-actions"
const operatorUnavailable = "runner or reviewed preset unavailable; inspect ailang-worldd service runner status or start at a terminal"

var operatorPolicy = verbs.OperatorPolicy{Start: "A", Finish: "A"} // Headless start is the only installed S7 backend.

func (d *Daemon) operatorVerb(r *http.Request, argv []string, out *verbAnswer) bool {
	a := verbs.AdmitOperator(argv, operatorPolicy)
	switch a.Class {
	case verbs.OperatorNever:
		return out.Class == string(verbs.Never) // Lookup already supplied the refusal.
	case verbs.OperatorReview:
		if a.Intent.Kind == "start" && a.Intent.Headless {
			out.Outcome, out.Reason, out.CLI = "terminal", operatorUnavailable, a.CLI
			ctx, cancel := d.readCtx(r)
			defer cancel()
			ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
			if err != nil {
				return true
			}
			defer ro.Close()
			h, entry, err := d.operatorReview(ctx, ro, a.Intent.Preset, d.approveNow())
			if err == nil {
				if action, err := entry.Action(a.Intent.Episode, "", h.CatalogueGeneration); err == nil {
					out.Outcome, out.Reason, out.Argv, out.Text = "review", "", action.Argv, action.Review
					out.CatalogueGeneration = h.CatalogueGeneration
				}
			}
			return true
		}
	}
	return false
}

func (d *Daemon) operatorReview(ctx context.Context, db *store.Credentials, preset string, now int64) (store.OperatorHeartbeat, store.OperatorStartReview, error) {
	h, err := db.OperatorHeartbeat(ctx)
	if err != nil {
		return h, store.OperatorStartReview{}, err
	}
	build := CurrentBuild().ExecutableSHA256
	if build == "" || !h.Available(build, "start", now) || h.CatalogueGeneration <= 0 {
		return h, store.OperatorStartReview{}, store.ErrOperatorRunner
	}
	entry, err := h.StartReview(preset)
	served, pathErr := store.CanonicalPath(d.credentialsPath)
	if err == nil && (pathErr != nil || entry.Store != served) {
		err = store.ErrOperatorRunner
	}
	return h, entry, err
}

// Handle the whole route family before ServeMux path cleaning. Malformed
// action paths receive a refusal and cannot become redirect destinations.
func (d *Daemon) operatorRoutes(next http.Handler) http.Handler {
	family := func(p string) bool { return p == operatorPath || strings.HasPrefix(p, operatorPath+"/") }
	operator := d.approvalRoute(d.handleOperatorRoutes)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if family(r.URL.Path) || family(path.Clean(r.URL.Path)) {
			operator(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (d *Daemon) handleOperatorRoutes(w http.ResponseWriter, r *http.Request, origins []string) {
	if r.URL.Path != r.URL.EscapedPath() || path.Clean(r.URL.Path) != r.URL.Path || r.URL.RawQuery != "" {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "use the exact operator action path")
		return
	}
	if r.URL.Path == operatorPath {
		if r.Method == http.MethodPost {
			d.handleOperatorRequest(w, r, origins)
			return
		}
	} else {
		id, suffix, hasSuffix := strings.Cut(strings.TrimPrefix(r.URL.Path, operatorPath+"/"), "/")
		if !verbs.ValidOperatorActionID(id) {
			writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "an operator action id is oa- and 32 lowercase hex digits")
			return
		}
		switch {
		case !hasSuffix || suffix == "status.json":
			if r.Method == http.MethodGet {
				d.handleOperatorRead(w, r, id, hasSuffix)
				return
			}
		case suffix == "confirm":
			if r.Method == http.MethodPost {
				d.handleOperatorConfirm(w, r, id, origins)
				return
			}
		default:
			writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no operator action route")
			return
		}
	}
	writeApprovalRefusal(w, http.StatusMethodNotAllowed, "MethodNotAllowed", "use GET for a card/status or POST to request an action")
}

// readOperatorPOST is reached after the Host and method checks, before any body read.
func readOperatorPOST(w http.ResponseWriter, r *http.Request, origins []string, cap int64) ([]byte, bool) {
	origin, fetch, media := r.Header.Values("Origin"), r.Header.Values("Sec-Fetch-Site"), r.Header.Values("Content-Type")
	switch {
	case len(r.Header.Values("Authorization")) != 0 || len(r.Header.Values("Cookie")) != 0:
		writeApprovalRefusal(w, http.StatusForbidden, "Forbidden", "use the operator card without a session or cookie")
	case len(origin) != 1 || !allowedApprovalOrigin(r.Header.Get("Origin"), origins):
		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", "use an exact approval origin")
	case len(fetch) > 1 || !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")):
		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", "Sec-Fetch-Site must be same-origin when sent")
	case len(media) != 1 || !jsonMediaType(r.Header.Get("Content-Type")):
		writeApprovalRefusal(w, http.StatusUnsupportedMediaType, "UnsupportedMediaType", "use application/json")
	default:
		raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, cap))
		var tooLarge *http.MaxBytesError
		switch {
		case errors.As(err, &tooLarge):
			writeApprovalRefusal(w, http.StatusRequestEntityTooLarge, "TooLarge", "the operator action body exceeds its byte limit")
		case err != nil:
			writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "the operator action body could not be read")
		default:
			return raw, true
		}
	}
	return nil, false
}

type operatorRequestBody struct {
	Argv                []string `json:"argv"`
	OperationID         string   `json:"operationId"`
	CatalogueGeneration int64    `json:"catalogueGeneration"`
	Task                string   `json:"task"`
}

func (d *Daemon) handleOperatorRequest(w http.ResponseWriter, r *http.Request, origins []string) {
	raw, ok := readOperatorPOST(w, r, origins, 64<<10)
	if !ok {
		return
	}
	var body operatorRequestBody
	if approval.DecodeOperatorObject(raw, &body, "argv", "operationId", "catalogueGeneration", "task") != nil ||
		len(body.Task) == 0 || len(body.Task) > approval.MaxOperatorTaskBytes || strings.ContainsRune(body.Task, 0) {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "send argv, operationId, catalogueGeneration and a task of 1–8192 UTF-8 bytes")
		return
	}
	a := verbs.AdmitOperator(body.Argv, operatorPolicy)
	if a.Class != verbs.OperatorReview || a.Intent.Kind != "start" || !a.Intent.Headless {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "review a supported headless start first")
		return
	}
	if body.CatalogueGeneration <= 0 || !store.ValidOperatorOperationID(body.OperationID) {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "use a positive catalogue generation and 32 lowercase hex operationId")
		return
	}

	d.humanMu.Lock()
	defer d.humanMu.Unlock()
	ctx, cancel := d.readCtx(r)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.operatorError(w, err)
		return
	}
	defer ro.Close()
	now := d.approveNow()
	if row, err := ro.OperatorActionForOperation(ctx, body.OperationID, now); !errors.Is(err, store.ErrOperatorMissing) {
		if err != nil {
			d.operatorError(w, err)
			return
		}
		frozen, err := approval.DecodeOperatorAction(row.Action)
		if err != nil || !slices.Equal(frozen.Argv, a.Intent.Argv) || frozen.Task != body.Task ||
			!slices.Contains(frozen.Resolution, approval.OperatorResolved{Name: "mint_catalogue", Value: strconv.FormatInt(body.CatalogueGeneration, 10)}) {
			d.operatorError(w, store.ErrOperatorConflict)
			return
		}
		writeApprovalJSON(w, http.StatusOK, operatorStatusFor(row, now))
		return
	}
	wr, err := d.openConfirmationWriter()
	if err != nil {
		d.operatorError(w, err)
		return
	}
	defer wr.Close()
	now = d.approveNow()
	h, entry, err := d.operatorReview(ctx, wr, a.Intent.Preset, now)
	if err == nil && h.CatalogueGeneration != body.CatalogueGeneration {
		err = store.ErrOperatorRunner
	}
	if err != nil {
		d.operatorError(w, err)
		return
	}
	action, err := entry.Action(a.Intent.Episode, body.Task, h.CatalogueGeneration)
	if err != nil {
		d.operatorError(w, err)
		return
	}
	var id [16]byte
	var challenge [32]byte
	if _, err := rand.Read(id[:]); err != nil {
		d.operatorError(w, err)
		return
	}
	if _, err := rand.Read(challenge[:]); err != nil {
		d.operatorError(w, err)
		return
	}
	action.ID = "oa-" + hex.EncodeToString(id[:])
	row, err := wr.RequestOperatorAction(ctx, store.OperatorRequest{Action: action, OperationID: body.OperationID, Challenge: challenge,
		CreatedAt: now, ExpiresAt: now + int64(approval.OperatorActionLifetime/time.Second)}, now)
	if err != nil {
		d.operatorError(w, err)
		return
	}
	writeApprovalJSON(w, http.StatusOK, operatorStatusFor(row, now))
}

func (d *Daemon) operatorError(w http.ResponseWriter, err error) {
	switch {
	case isOlderStore(err), errors.Is(err, store.ErrOperatorRunner):
		writeApprovalRefusal(w, http.StatusServiceUnavailable, "RunnerUnavailable", operatorUnavailable+"; review again")
	case store.IsWriterAlreadyActive(err):
		w.Header().Set("Retry-After", "1")
		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Busy", "credentials writer busy; retry")
	case errors.Is(err, store.ErrOperatorMissing):
		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no operator action")
	case errors.Is(err, store.ErrOperatorCapacity):
		writeApprovalRefusal(w, http.StatusConflict, "QueueFull", "sixteen actions are already pending")
	case errors.Is(err, store.ErrOperatorConflict):
		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "operator action changed or conflicts with existing work")
	case errors.Is(err, store.ErrOperatorInput), errors.Is(err, approval.ErrOperatorTooLarge):
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "invalid reviewed action")
	default:
		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Disconnected", "operator actions could not be read or written; retry")
	}
}
