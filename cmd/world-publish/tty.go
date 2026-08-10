package main

import (
	"io/fs"
	"os"
)

// ---------------------------------------------------------------------------
// THE LOAD-BEARING FENCE: a controlling terminal
//
// The strongest fence already landed in this repository is the APPROVAL OBJECT:
// validatePublishApproval refuses before the credential is loaded and before
// any POST, and approval_claims' PRIMARY KEY makes the stamp single-use durably
// across a process restart. But it is not a HEADLESS/ATTENDED fence — a headless
// loop can mint one by calling broker.DecideApproval. It is a SINGLE-USE fence.
// That distinction is why a new layer is needed.
//
// Every candidate for the new layer was rejected by asking one question: can
// THIS loop satisfy it?
//
//	an environment variable  -> defeated by `env -u`
//	a typed phrase alone     -> defeated by `echo ... |`
//	a sentinel file          -> this loop can create files
//	isatty via ModeCharDevice alone -> REFUTED BY MEASUREMENT (see below)
//
// A controlling terminal is the one candidate this loop is STRUCTURALLY unable
// to satisfy. Measured first-party in the executing agent's own shell at
// 6d1dce0:
//
//	os.Stdin.Stat()      -> mode Srw-rw-rw-, chardev=false   (a SOCKET)
//	os.Open("/dev/tty")  -> "device not configured" (ENXIO)
//	os.Stat("/dev/null") -> mode Dcrw-rw-rw-, chardev=TRUE
//	os.SameFile(/dev/null, /dev/null) -> true
//
// The third line is why R-TTY-SAMEFILE exists. /dev/null IS a character device,
// so a naive isatty check admits `world-publish publish --live < /dev/null` —
// the mission's recurring class, a check that looks like verification. The
// fourth line is the positive control: the passing branch is reachable in a unit
// test with no pty and no new dependency, because /dev/null is the same file as
// itself. In PRODUCTION the ctty FileInfo comes only from os.Open("/dev/tty"),
// which never resolves to /dev/null.
//
// Stdlib only. No build tags. Works on darwin and linux.
// ---------------------------------------------------------------------------

// devTTY is the controlling-terminal device. It is the kernel's answer to "is
// there a human at this process", not a heuristic about file descriptors: a
// process with no controlling terminal cannot open it at all.
const devTTY = "/dev/tty"

// ttyProbe is one observation of this process's terminal situation. It is a
// VALUE so the three refusals below can be driven from a unit test without a
// pty, a subprocess or a new dependency — and so the production path can be
// asserted to build it from nothing but the two syscalls named here.
type ttyProbe struct {
	// stdin is os.Stdin's FileInfo, or nil if it could not be stat'ed.
	stdin fs.FileInfo
	// ctty is the FileInfo of an opened /dev/tty, or nil.
	ctty fs.FileInfo
	// cttyErr is the error from opening /dev/tty. It is carried rather than
	// collapsed into "ctty == nil" because "there is no controlling terminal"
	// is the single most informative thing this fence can tell an operator.
	cttyErr error
}

// probeControllingTerminal performs the two syscalls. It is the ONLY place the
// production path touches the terminal, and it opens /dev/tty read-only and
// closes it immediately: the fence needs the file's identity, not a handle.
func probeControllingTerminal() ttyProbe {
	var p ttyProbe
	if info, err := os.Stdin.Stat(); err == nil {
		p.stdin = info
	}
	tty, err := os.OpenFile(devTTY, os.O_RDONLY, 0)
	if err != nil {
		p.cttyErr = err
		return p
	}
	defer func() { _ = tty.Close() }()
	if info, statErr := tty.Stat(); statErr == nil {
		p.ctty = info
	} else {
		p.cttyErr = statErr
	}
	return p
}

// requireControllingTerminal is the fence. Its three refusals are independent
// and are ordered so that each one's precondition is established by the one
// before it.
func requireControllingTerminal(p ttyProbe) *stopError {
	// R-TTY-OPEN. Measured to fire in this loop's own shell today.
	if p.cttyErr != nil {
		return &stopError{
			Fence:  fenceTTY,
			Reason: "no-controlling-terminal",
			Detail: "opening " + devTTY + " failed: " + p.cttyErr.Error(),
		}
	}
	// R-TTY-CHARDEV. Measured to fire in this loop's own shell today: stdin is
	// a socket.
	if p.stdin == nil || p.stdin.Mode()&os.ModeCharDevice == 0 {
		return &stopError{
			Fence:  fenceTTY,
			Reason: "stdin-not-a-terminal",
			Detail: describeStdin(p.stdin),
		}
	}
	// R-TTY-SAMEFILE. THE REPAIR: /dev/null is a character device, so the check
	// above alone admits `--live < /dev/null`. Redirected stdin is not the
	// controlling terminal, whatever kind of device it happens to be.
	if !sameFile(p.stdin, p.ctty) {
		return &stopError{
			Fence:  fenceTTY,
			Reason: "stdin-is-not-the-controlling-terminal",
			Detail: "stdin is a character device but not " + devTTY +
				"; a redirect such as `< /dev/null` is not an attended operator",
		}
	}
	return nil
}

// sameFile is os.SameFile with a nil-tolerant front. It fails CLOSED: an
// observation this fence could not make is not evidence that the two are the
// same file.
//
// The nil case is reachable — that is the point. Neutering R-TTY-OPEN (mutation
// MUT-D0-04) leaves ctty nil while execution continues, and os.SameFile would
// panic on a nil FileInfo's Sys(). A mutant that panics still reds, but it reds
// as a crash rather than as the refusal the row names, which is a worse signal.
func sameFile(a, b fs.FileInfo) bool {
	if a == nil || b == nil {
		return false
	}
	return os.SameFile(a, b)
}

// describeStdin names what stdin actually was, so an operator who gets this
// refusal in a terminal knows to look at their redirect rather than their tty.
func describeStdin(info fs.FileInfo) string {
	if info == nil {
		return "stdin could not be stat'ed"
	}
	return "stdin mode is " + info.Mode().String() + ", which carries no character-device bit"
}
