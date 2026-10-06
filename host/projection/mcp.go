package projection

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang/serveapi/protocol"
	"github.com/sunholo-data/ailang/serveapi/protocol/hostcall"
	"github.com/sunholo-data/ailang/serveapi/protocol/mcphttp"
	"net/http"
	"time"
)

func newMCPHandler(h *Handler, cfg Config) (http.Handler, error) {
	runner, err := hostcall.New(cfg.CallbackTimeout, cfg.MaxCallbacks)
	if err != nil {
		return nil, err
	}
	adapter := mcpAdapter{h}
	return mcphttp.NewHandler(mcphttp.Config{Agent: cfg.Agent, Resolver: adapter, Tools: adapter, Invoker: adapter, Runner: runner})
}

// MCP supplies the single aggregate request deadline; the released handler owns
// request parsing, batch iteration, all wire bytes and the response envelope.
func (h *Handler) MCP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.invokeWait)
	defer cancel()
	h.mcp.ServeHTTP(w, r.WithContext(ctx))
}

type mcpAdapter struct{ h *Handler }

func (a mcpAdapter) ResolveSession(ctx context.Context, r *http.Request) (protocol.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, a.h.credentialWait)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	out, err := a.h.resolver.ResolveContext(ctx, r.Header.Get("Authorization"), time.Now().Unix())
	if err != nil {
		return nil, err
	}
	if out.Denied != nil {
		message := msgUnknown
		switch *out.Denied {
		case authority.DenialAbsent:
			message = msgAbsent
		case authority.DenialMalformed:
			message = msgMalformed
		case authority.DenialExpired:
			message = msgExpired
		}
		return nil, &protocol.AuthorizationError{Status: http.StatusUnauthorized, Err: errors.New(message)}
	}
	if out.Success == nil {
		return nil, errors.New("projection: missing session binding")
	}
	return out.Success, nil
}
func (a mcpAdapter) Tools(ctx context.Context, session protocol.Session) ([]protocol.ToolDescriptor, error) {
	binding, ok := session.(*authority.SessionBinding)
	if !ok || binding == nil {
		return nil, errors.New("projection: invalid session binding")
	}
	ctx, cancel := context.WithTimeout(ctx, a.h.maxWait)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	_, ds, err := a.h.allowedDescriptors(ctx, binding)
	if err != nil {
		a.h.logRefusal("mcp", "tools/list", "-", err)
		return nil, err
	}
	tools, err := mcpDescriptors(ds)
	if err != nil {
		a.h.logRefusal("mcp", "tools/list", "-", err)
	}
	return tools, err
}
func (a mcpAdapter) Invoke(ctx context.Context, session protocol.Session, inv protocol.Invocation) (protocol.InvocationResult, error) {
	binding, ok := session.(*authority.SessionBinding)
	if !ok || binding == nil {
		return protocol.InvocationResult{}, errors.New("projection: invalid session binding")
	}
	ctx, cancel := context.WithTimeout(ctx, a.h.invokeWait)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return protocol.InvocationResult{}, err
	}
	id, err := decodeMCPName(inv.Name)
	if err != nil {
		return protocol.InvocationResult{}, err
	}
	var input map[string]any
	if err := json.Unmarshal(inv.Arguments, &input); err != nil || input == nil {
		return protocol.InvocationResult{}, errors.New("projection: arguments must be a JSON object")
	}
	admissionCtx, stop := context.WithTimeout(ctx, a.h.maxWait)
	request, ds, err := a.h.allowedDescriptors(admissionCtx, binding)
	stop()
	if err != nil {
		a.h.logRefusal("mcp", "tools/call", "-", err)
		return protocol.InvocationResult{}, err
	}
	allowed := false
	for _, d := range ds {
		if d.ID == id {
			allowed = true
			break
		}
	}
	if !allowed {
		return protocol.InvocationResult{}, errors.New("projection: transition is not authorized")
	}
	if err := ctx.Err(); err != nil {
		return protocol.InvocationResult{}, err
	}
	if a.h.coord == nil {
		return protocol.InvocationResult{}, errors.New("projection: invocation coordinator is unavailable")
	}
	task, err := a.h.mintTask()
	if err != nil {
		return protocol.InvocationResult{}, err
	}
	result, err := a.h.coord.Dispatch(ctx, coordinator.Call{Request: request, EpisodeID: binding.EpisodeID, Grants: binding.Caps, SkillID: id, TaskID: task, Input: input})
	if err != nil {
		a.h.logRefusal("mcp", "tools/call", coordinator.InvocationID(binding.EpisodeID, task), err)
		return protocol.InvocationResult{}, err
	}
	return protocol.InvocationResult{Value: result.OutputBytes}, nil
}
func mintMCPTask() (string, error) {
	var bytes [32]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes[:]), nil
}
