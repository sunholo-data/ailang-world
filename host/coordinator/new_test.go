package coordinator

import (
	"context"
	"errors"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/capsule"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

type constructionRunner struct{}

func (constructionRunner) RunContext(_ context.Context, _ capsule.Entry) (capsule.Result, error) {
	return capsule.Result{}, nil
}

func TestNewRefusesMissingSeamsAndCaps(t *testing.T) {
	valid := Config{Store: (*store.Store)(nil), Runner: constructionRunner{}, Binder: func(string, []broker.Capability) transitionreg.Binder { return nil }, Now: func() int64 { return 1 }, MaxInput: 10, MaxOutput: 10}
	if _, err := New(Config{}); err == nil {
		t.Fatal("accepted missing seams")
	}
	checks := []struct {
		name   string
		change func(*Config)
	}{
		{"store", func(c *Config) { c.Store = nil }}, {"runner", func(c *Config) { c.Runner = nil }},
		{"binder", func(c *Config) { c.Binder = nil }}, {"clock", func(c *Config) { c.Now = nil }},
		{"input cap", func(c *Config) { c.MaxInput = 0 }}, {"output cap", func(c *Config) { c.MaxOutput = -1 }},
	}
	for _, tc := range checks {
		t.Run(tc.name, func(t *testing.T) {
			c := valid
			tc.change(&c)
			if _, err := New(c); err == nil {
				t.Fatal("accepted invalid config")
			}
		})
	}
	c, err := New(valid)
	if err != nil {
		t.Fatal(err)
	}
	if InvocationID("ep", "task") != "a2a:ep:task" {
		t.Fatal("wrong invocation namespace")
	}
	if validTaskID("bad id") || !validTaskID("a-1") {
		t.Fatal("task validation")
	}
	if _, _, err := c.parseOutput([]byte("[]")); err == nil {
		t.Fatal("accepted array output")
	}
	if _, _, err := c.parseOutput([]byte(`{"a":1}`)); err != nil {
		t.Fatal(err)
	}
	execErr := &capsule.ExecError{Stderr: []byte("ARG_DECODE_MISMATCH")}
	var incompatible *IncompatibleError
	if !errors.As(classifyExec(execErr), &incompatible) {
		t.Fatal("wrong execution classification")
	}
}
