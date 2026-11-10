// Package sshsign is the one bounded, environment-scrubbed subprocess site of `ailang-worldd
// approve --sign-with` (row 201 A1b, design_docs/planned/w-world-approvals-by-key-a1b.md §2.9):
// `ssh-keygen -Y sign -n world-approval@sunholo -f <key>` with the message on stdin, so the
// signature comes back on stdout and nothing is written beside the key.
//
// Why a package: the AC10 census (host/broker registry_publish_test.go) requires every
// production subprocess site to be DRIVEN and to scrub the registry credential, and a site in
// package main cannot be driven from there, so the site lives here (as host/toolprobe's does)
// with an injectable binary path the census points at its env-dumping probe.
//
// The child environment is the process environment minus every registry variable
// (childenv.Scrubbed): SSH_AUTH_SOCK, HOME and PATH, which ssh-keygen needs, pass through.
// PATH is not a trust boundary: whatever runs, the caller verifies its output against the
// enrolled key over World's own message.
package sshsign

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/childenv"
)

// MaxSignatureBytes bounds the signer's output (G2): an Ed25519 or P-256 SSHSIG is under 1 KiB.
const MaxSignatureBytes = 16 << 10

// waitDelay bounds Wait after the deadline's kill, even if a child of the signer keeps a pipe open.
const waitDelay = time.Second

var (
	ErrSignatureTooLarge = fmt.Errorf("ssh-keygen -Y sign printed over %d bytes", MaxSignatureBytes)
	ErrEmptySignature    = errors.New("ssh-keygen -Y sign printed nothing")
)

// Sign runs `<bin> -Y sign -n world-approval@sunholo -f <key>` under ctx with msg on stdin and
// returns its standard output, the armored signature. stderr (where ssh-keygen reports, and
// from where a passphrase prompt goes to /dev/tty) is the operator's. A failed run (G1), over
// MaxSignatureBytes of output (G2) and no output (G3) are errors.
func Sign(ctx context.Context, bin, key string, msg []byte, stderr io.Writer) ([]byte, error) {
	cmd := exec.CommandContext(ctx, bin, "-Y", "sign", "-n", approval.Namespace, "-f", key)
	cmd.Env = childenv.Scrubbed(os.Environ())
	cmd.Stdin = bytes.NewReader(msg)
	out := &cappedBuffer{max: MaxSignatureBytes}
	cmd.Stdout, cmd.Stderr = out, stderr
	cmd.WaitDelay = waitDelay
	if err := cmd.Run(); err != nil { // G1
		return nil, fmt.Errorf("ssh-keygen -Y sign: %w", err)
	}
	if out.over { // G2
		return nil, ErrSignatureTooLarge
	}
	if out.buf.Len() == 0 { // G3
		return nil, ErrEmptySignature
	}
	return out.buf.Bytes(), nil
}

// cappedBuffer keeps the first max bytes written to it and records whether more came.
type cappedBuffer struct {
	buf  bytes.Buffer
	max  int
	over bool
}

func (c *cappedBuffer) Write(p []byte) (int, error) {
	n := min(len(p), c.max-c.buf.Len())
	c.buf.Write(p[:n])
	c.over = c.over || n < len(p)
	return len(p), nil
}
