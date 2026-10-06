package main

import (
	"io"
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
// The third line is why R-TTY-SAMEFILE existed. /dev/null IS a character
// device, so a naive isatty check admits `world-publish publish --live <
// /dev/null` — the mission's recurring class, a check that looks like
// verification. The fourth line is the positive control: the stdin path is
// reachable in a unit test with no pty, because /dev/null is the same file as
// itself. In PRODUCTION the ctty FileInfo comes only from os.Open("/dev/tty"),
// which never resolves to /dev/null.
//
// WHERE THE TYPED LINE COMES FROM (2026-10-06, attended-steps ergonomics).
// R-TTY-CHARDEV and R-TTY-SAMEFILE used to REFUSE any process whose stdin was
// not the very file /dev/tty. That stopped every human too: in an IDE terminal
// pane — and in an ordinary terminal window, whose stdin is /dev/ttysNNN, a
// different file from /dev/tty — the operator was told
// `reason=stdin-is-not-the-controlling-terminal`, and the documented way past it
// was to re-run with `< /dev/tty`. The fence now does that itself: when stdin is
// not the controlling terminal, the prompt is written to, and the typed line is
// read from, the /dev/tty handle the probe opened.
//
// THIS IS EXACTLY AS STRONG AS THE `< /dev/tty` ROUTE QUICKSTART ALREADY
// SANCTIONED. That route reads the line from /dev/tty too; here the command
// opens the same device instead of the shell. Either way the line can only
// come from the controlling terminal, and a process that has none cannot open
// /dev/tty at all — R-TTY-OPEN, unchanged, is still the refusal an agent
// harness, setsid, cron or CI meets. What a redirected stdin carries (`echo
// phrase |`, `< /dev/null`) is never read as the confirmation: the stdin
// identity checks did not stop being enforced, they stopped deciding WHETHER to
// refuse and now decide WHERE to read.
//
// Stdlib only. No build tags. Works on darwin and linux.
// ---------------------------------------------------------------------------

// devTTY is the controlling-terminal device. It is the kernel's answer to "is
// there a human at this process", not a heuristic about file descriptors: a
// process with no controlling terminal cannot open it at all.
const devTTY = "/dev/tty"

// ttyFix is the exact fix every terminal refusal ends with. An operator who
// reads a STOP line needs the next command, not the theory.
const ttyFix = "fix: run it yourself from a terminal window (an IDE terminal pane works); an agent cannot run this step"

// ttyProbe is one observation of this process's terminal situation. It is a
// VALUE so the refusal and the choice of input below can be driven from a unit
// test without a pty, a subprocess or a new dependency — and so the production
// path can be asserted to build it from nothing but the two syscalls named here.
type ttyProbe struct {
	// stdin is os.Stdin's FileInfo, or nil if it could not be stat'ed.
	stdin fs.FileInfo
	// ctty is the FileInfo of an opened /dev/tty, or nil.
	ctty fs.FileInfo
	// cttyErr is the error from opening /dev/tty. It is carried rather than
	// collapsed into "ctty == nil" because "there is no controlling terminal"
	// is the single most informative thing this fence can tell an operator.
	cttyErr error
	// term is the opened /dev/tty itself, read-write. When stdin is not the
	// controlling terminal, the prompt is written here and the typed line is
	// read from here. nil whenever cttyErr is set.
	term io.ReadWriter
}

// probeControllingTerminal performs the two syscalls. It is the ONLY place the
// production path touches the terminal. It opens /dev/tty read-write and keeps
// the handle in the probe, because the confirmation may be read from it; the
// fence stack closes it once the line has been read.
func probeControllingTerminal() ttyProbe {
	var p ttyProbe
	if info, err := os.Stdin.Stat(); err == nil {
		p.stdin = info
	}
	tty, err := os.OpenFile(devTTY, os.O_RDWR, 0)
	if err != nil {
		p.cttyErr = err
		return p
	}
	info, statErr := tty.Stat()
	if statErr != nil {
		_ = tty.Close()
		p.cttyErr = statErr
		return p
	}
	p.ctty = info
	p.term = tty
	return p
}

// requireControllingTerminal is the fence: a process with no controlling
// terminal is refused. It is the only terminal refusal; which device the typed
// line is then read from is confirmationSource's decision.
func requireControllingTerminal(p ttyProbe) *stopError {
	// R-TTY-OPEN. Measured to fire in this loop's own shell today.
	if p.cttyErr != nil {
		return noControllingTerminal("opening " + devTTY + " failed: " + p.cttyErr.Error())
	}
	return nil
}

func noControllingTerminal(cause string) *stopError {
	return &stopError{
		Fence:  fenceTTY,
		Reason: "no-controlling-terminal",
		Detail: cause + "; there is no controlling terminal here (an agent harness, setsid, cron or CI). " + ttyFix,
	}
}

// confirmationSource picks where the prompt goes and the typed line comes
// from. When stdin IS the controlling terminal (`< /dev/tty`, or a process
// whose stdin was opened from it) that is stdin and stdout, as before. In every
// other case it is the opened /dev/tty — never the redirected stdin.
//
// It fails CLOSED: a stdin that is not the terminal with no terminal handle to
// read from is refused as no-controlling-terminal, not read. In production
// that state cannot follow a passing requireControllingTerminal; the arm is
// here so that neutering R-TTY-OPEN alone cannot turn a redirected stdin back
// into the confirmation.
func confirmationSource(p ttyProbe, in io.Reader, out io.Writer) (io.Reader, io.Writer, *stopError) {
	if sameFile(p.stdin, p.ctty) {
		return in, out, nil
	}
	if p.term == nil {
		return nil, nil, noControllingTerminal(devTTY + " is not open for reading")
	}
	return p.term, p.term, nil
}

// closeTerminal releases the probe's /dev/tty handle, if it holds one.
func closeTerminal(p ttyProbe) {
	if c, ok := p.term.(io.Closer); ok {
		_ = c.Close()
	}
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
