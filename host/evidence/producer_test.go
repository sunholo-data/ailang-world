package evidence_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// fakeWriter is a constructor-only ObjectWriter: it reports a busy window and
// never stores anything.
type fakeWriter struct{ busy time.Duration }

func (f *fakeWriter) PutObject(context.Context, store.Object) error { return nil }
func (f *fakeWriter) BusyTimeout() time.Duration                    { return f.busy }

func validProducerArgs() (evidence.ObjectReader, evidence.ObjectWriter, evidence.CheckerConfig, evidence.ProducerConfig) {
	return &fakeReader{}, &fakeWriter{}, evidence.CheckerConfig{
			ExecPath:        "/nonexistent/ailang",
			Compiler:        testCompiler,
			CompilerVersion: "AILANG v0.41.0",
		}, evidence.ProducerConfig{
			ObjectReadTimeout:  3 * time.Second,
			ObjectWriteTimeout: 3 * time.Second,
			ExecTimeout:        time.Minute,
			MaxOutputBytes:     1 << 20,
		}
}

// requireConfigRefusal demands that build refuse with ErrInvalidProducerConfig
// and a detail naming only the field under test. A panic is a failure, not a pass.
func requireConfigRefusal(t *testing.T, field string, build func() error) {
	t.Helper()
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("NewProducer panicked: %v; want ErrInvalidProducerConfig naming %s", r, field)
			}
		}()
		err = build()
	}()
	if !errors.Is(err, evidence.ErrInvalidProducerConfig) || !strings.Contains(err.Error(), field) {
		t.Fatalf("want ErrInvalidProducerConfig naming %s; got %v", field, err)
	}
}

func requireValidProducerControl(t *testing.T) {
	t.Helper()
	reader, writer, checker, cfg := validProducerArgs()
	p, err := evidence.NewProducer(testKey, reader, writer, checker, cfg)
	if err != nil || p == nil {
		t.Fatalf("valid control must construct: p=%v err=%v", p, err)
	}
}

func TestProducerConstructorRefusesNonPositiveObjectReadTimeout(t *testing.T) {
	requireValidProducerControl(t)
	for _, d := range []time.Duration{0, -time.Second} {
		requireConfigRefusal(t, "ObjectReadTimeout", func() error {
			reader, writer, checker, cfg := validProducerArgs()
			cfg.ObjectReadTimeout = d
			_, err := evidence.NewProducer(testKey, reader, writer, checker, cfg)
			return err
		})
	}
}

func TestProducerConstructorRefusesNonPositiveObjectWriteTimeout(t *testing.T) {
	requireValidProducerControl(t)
	for _, d := range []time.Duration{0, -time.Second} {
		requireConfigRefusal(t, "ObjectWriteTimeout", func() error {
			reader, writer, checker, cfg := validProducerArgs()
			cfg.ObjectWriteTimeout = d
			_, err := evidence.NewProducer(testKey, reader, writer, checker, cfg)
			return err
		})
	}
}

func TestProducerConstructorRefusesNonPositiveExecTimeout(t *testing.T) {
	requireValidProducerControl(t)
	for _, d := range []time.Duration{0, -time.Second} {
		requireConfigRefusal(t, "ExecTimeout", func() error {
			reader, writer, checker, cfg := validProducerArgs()
			cfg.ExecTimeout = d
			_, err := evidence.NewProducer(testKey, reader, writer, checker, cfg)
			return err
		})
	}
}

func TestProducerConstructorRefusesNonPositiveMaxOutputBytes(t *testing.T) {
	requireValidProducerControl(t)
	for _, n := range []int64{0, -1} {
		requireConfigRefusal(t, "MaxOutputBytes", func() error {
			reader, writer, checker, cfg := validProducerArgs()
			cfg.MaxOutputBytes = n
			_, err := evidence.NewProducer(testKey, reader, writer, checker, cfg)
			return err
		})
	}
}

func TestProducerConstructorRefusesNilReader(t *testing.T) {
	requireValidProducerControl(t)
	requireConfigRefusal(t, "reader", func() error {
		_, writer, checker, cfg := validProducerArgs()
		_, err := evidence.NewProducer(testKey, nil, writer, checker, cfg)
		return err
	})
}

func TestProducerConstructorRefusesNilWriter(t *testing.T) {
	requireValidProducerControl(t)
	requireConfigRefusal(t, "writer", func() error {
		reader, _, checker, cfg := validProducerArgs()
		_, err := evidence.NewProducer(testKey, reader, nil, checker, cfg)
		return err
	})
}

func TestProducerConstructorRefusesEmptyExecPath(t *testing.T) {
	requireValidProducerControl(t)
	requireConfigRefusal(t, "ExecPath", func() error {
		reader, writer, checker, cfg := validProducerArgs()
		checker.ExecPath = ""
		_, err := evidence.NewProducer(testKey, reader, writer, checker, cfg)
		return err
	})
}

func TestProducerConstructorRefusesUnsetCompilerIdentity(t *testing.T) {
	requireValidProducerControl(t)
	requireConfigRefusal(t, "Compiler", func() error {
		reader, writer, checker, cfg := validProducerArgs()
		checker.Compiler = hashref.HashRef{}
		_, err := evidence.NewProducer(testKey, reader, writer, checker, cfg)
		return err
	})
}

func TestProducerConstructorRefusesEmptyCompilerVersion(t *testing.T) {
	requireValidProducerControl(t)
	requireConfigRefusal(t, "CompilerVersion", func() error {
		reader, writer, checker, cfg := validProducerArgs()
		checker.CompilerVersion = ""
		_, err := evidence.NewProducer(testKey, reader, writer, checker, cfg)
		return err
	})
}
