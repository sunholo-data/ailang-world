package broker

import (
	"context"
	"fmt"
	"os"
	"time"
)

const EffectGitCommit = "Git.Commit"

// GitHandlerConfig pins the executable used for Git.Commit. Bounds are kept
// internal and injectable by package tests so production and tests take the
// identical runBounded branch.
type GitHandlerConfig struct {
	GitPath        string
	ExecTimeout    time.Duration
	MaxOutputBytes int64
}

type GitHandler struct {
	gitPath string
	bounds  handlerBounds
}

func NewGitHandler(cfg GitHandlerConfig) (*GitHandler, error) {
	if cfg.GitPath == "" {
		return nil, fmt.Errorf("broker: Git.Commit requires an explicit git executable")
	}
	return &GitHandler{
		gitPath: cfg.GitPath,
		bounds: handlerBounds{
			execTimeout: cfg.ExecTimeout, maxOutputBytes: cfg.MaxOutputBytes,
		}.normalized(),
	}, nil
}

func (h *GitHandler) Execute(ctx context.Context, req EffectRequest, payload []byte) ([]byte, error) {
	if req.Effect != EffectGitCommit {
		return nil, fmt.Errorf("broker: Git handler does not implement %q", req.Effect)
	}
	home, err := os.MkdirTemp("", "broker-git-home-*")
	if err != nil {
		return nil, fmt.Errorf("broker: create empty git HOME: %w", err)
	}
	defer func() { _ = os.RemoveAll(home) }()

	message := string(payload)
	if message == "" {
		message = "AILANG World effect"
	}
	return runBounded(ctx, h.bounds, handlerCommand{
		path: h.gitPath,
		args: []string{"commit", "--no-edit", "-m", message},
		dir:  req.Scope,
		env: []string{
			"HOME=" + home,
			"PATH=/usr/bin:/bin",
			"LANG=C",
			"LC_ALL=C",
			// The empty HOME intentionally removes ambient identity. These four
			// deterministic constants make commit usable without inheriting or
			// leaking any caller-provided GIT_* values.
			"GIT_AUTHOR_NAME=AILANG World",
			"GIT_AUTHOR_EMAIL=ailang-world@invalid",
			"GIT_COMMITTER_NAME=AILANG World",
			"GIT_COMMITTER_EMAIL=ailang-world@invalid",
		},
	})
}
