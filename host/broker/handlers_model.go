package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const EffectModelInfer = "Model.Infer"

const modelProgram = `module host/broker/model_infer

import std/ai (call)

export func main(prompt: string) -> string ! {AI} {
  call(prompt)
}
`

type ModelHandlerConfig struct {
	AILANGPath     string
	Stub           bool
	Model          string
	ExecTimeout    time.Duration
	MaxOutputBytes int64
}

type ModelHandler struct {
	ailangPath string
	stub       bool
	model      string
	bounds     handlerBounds
}

func NewModelHandler(cfg ModelHandlerConfig) (*ModelHandler, error) {
	if cfg.AILANGPath == "" {
		return nil, fmt.Errorf("broker: Model.Infer requires an explicit ailang executable")
	}
	if !cfg.Stub && cfg.Model == "" {
		return nil, fmt.Errorf("broker: live Model.Infer requires an explicit model")
	}
	return &ModelHandler{
		ailangPath: cfg.AILANGPath, stub: cfg.Stub, model: cfg.Model,
		bounds: handlerBounds{
			execTimeout: cfg.ExecTimeout, maxOutputBytes: cfg.MaxOutputBytes,
		}.normalized(),
	}, nil
}

func (h *ModelHandler) Execute(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	if req.Effect != EffectModelInfer {
		return nil, fmt.Errorf("broker: Model handler does not implement %q", req.Effect)
	}
	root, err := os.MkdirTemp("", "broker-model-*")
	if err != nil {
		return nil, fmt.Errorf("broker: create model project: %w", err)
	}
	defer func() { _ = os.RemoveAll(root) }()
	rel := filepath.FromSlash("host/broker/model_infer.ail")
	source := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(source), 0o700); err != nil {
		return nil, fmt.Errorf("broker: create model source directory: %w", err)
	}
	if err := os.WriteFile(source, []byte(modelProgram), 0o600); err != nil {
		return nil, fmt.Errorf("broker: stage model program: %w", err)
	}
	promptJSON, err := json.Marshal(string(payload))
	if err != nil {
		return nil, fmt.Errorf("broker: encode model prompt: %w", err)
	}
	args := []string{"run", "--quiet", "--caps", "AI", "--entry", "main",
		"--args-json", string(promptJSON)}
	if h.stub {
		args = append(args, "--ai-stub")
	} else {
		args = append(args, "--ai", h.model)
	}
	args = append(args, filepath.ToSlash(rel))
	return runBounded(ctx, h.bounds, handlerCommand{
		path: h.ailangPath, args: args, dir: root,
		env: []string{"HOME=" + root, "PATH=/usr/bin:/bin", "LANG=C", "LC_ALL=C"},
	})
}
