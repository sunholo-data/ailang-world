package evidence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// ErrInvalidProducerConfig refuses a producer whose configuration cannot make
// a bounded, pinned proof. Each refusal names only its own field.
var ErrInvalidProducerConfig = errors.New("evidence: invalid producer configuration")

// ObjectWriter is the write seam of the producer. BusyTimeout reports the
// store's lock-retry window so the write timeout can be ordered above it.
type ObjectWriter interface {
	PutObject(context.Context, store.Object) error
	BusyTimeout() time.Duration
}

// CheckerConfig pins the external checker by its bytes and its version line.
type CheckerConfig struct {
	ExecPath        string
	Compiler        hashref.HashRef
	CompilerVersion string
}

// ProducerConfig bounds every leg of a proof: the two object calls and the
// checker's wall clock and captured output.
type ProducerConfig struct {
	ObjectReadTimeout  time.Duration
	ObjectWriteTimeout time.Duration
	ExecTimeout        time.Duration
	MaxOutputBytes     int64
}

// Producer mints authenticated proof envelopes from a pinned checker run.
type Producer struct {
	key     *[32]byte
	reader  ObjectReader
	writer  ObjectWriter
	checker CheckerConfig
	cfg     ProducerConfig
}

// NewProducer refuses a configuration it cannot bound. Each refusal is a
// separate branch so each is separately observable and separately killable.
func NewProducer(key [32]byte, reader ObjectReader, writer ObjectWriter, checker CheckerConfig, cfg ProducerConfig) (*Producer, error) {
	if reader == nil {
		return nil, fmt.Errorf("%w: reader is nil", ErrInvalidProducerConfig)
	}
	if writer == nil {
		return nil, fmt.Errorf("%w: writer is nil", ErrInvalidProducerConfig)
	}
	if checker.ExecPath == "" {
		return nil, fmt.Errorf("%w: ExecPath is empty", ErrInvalidProducerConfig)
	}
	if checker.Compiler.IsZero() {
		return nil, fmt.Errorf("%w: Compiler identity is unset", ErrInvalidProducerConfig)
	}
	if checker.CompilerVersion == "" {
		return nil, fmt.Errorf("%w: CompilerVersion is empty", ErrInvalidProducerConfig)
	}
	if cfg.ObjectReadTimeout <= 0 {
		return nil, fmt.Errorf("%w: ObjectReadTimeout must be positive, got %s", ErrInvalidProducerConfig, cfg.ObjectReadTimeout)
	}
	if cfg.ObjectWriteTimeout <= 0 {
		return nil, fmt.Errorf("%w: ObjectWriteTimeout must be positive, got %s", ErrInvalidProducerConfig, cfg.ObjectWriteTimeout)
	}
	if cfg.ExecTimeout <= 0 {
		return nil, fmt.Errorf("%w: ExecTimeout must be positive, got %s", ErrInvalidProducerConfig, cfg.ExecTimeout)
	}
	if cfg.MaxOutputBytes <= 0 {
		return nil, fmt.Errorf("%w: MaxOutputBytes must be positive, got %d", ErrInvalidProducerConfig, cfg.MaxOutputBytes)
	}
	seams := []struct {
		name string
		busy time.Duration
	}{{"reader", reader.BusyTimeout()}, {"writer", writer.BusyTimeout()}}
	for _, s := range seams {
		if s.busy < 0 {
			return nil, fmt.Errorf("%w: %s BusyTimeout is unknown (%s)", ErrInvalidProducerConfig, s.name, s.busy)
		}
		if s.busy > 0 && cfg.ObjectReadTimeout <= s.busy {
			return nil, fmt.Errorf("%w: ObjectReadTimeout %s must exceed %s BusyTimeout %s", ErrUnorderedTimeouts, cfg.ObjectReadTimeout, s.name, s.busy)
		}
		if s.busy > 0 && cfg.ObjectWriteTimeout <= s.busy {
			return nil, fmt.Errorf("%w: ObjectWriteTimeout %s must exceed %s BusyTimeout %s", ErrUnorderedTimeouts, cfg.ObjectWriteTimeout, s.name, s.busy)
		}
	}
	keyCopy := new([32]byte)
	*keyCopy = key
	return &Producer{key: keyCopy, reader: reader, writer: writer, checker: checker, cfg: cfg}, nil
}
