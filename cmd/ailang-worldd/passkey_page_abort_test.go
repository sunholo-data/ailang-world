package main

// Row 201 A2 sprint 4 M7, XB2 (design §2.14 [r2], §6): the page's own 240000 ms timer, in virtual
// time, against an injected credentials operation.

import (
	"encoding/base64"
	"math"
	"testing"
)

// XB2 [r2]: a ceremony the authenticator never answers is aborted by the page's own 240000 ms
// timer, driven in virtual time (paused from before the click, so the timer and the operation start at
// one instant): the injected operation saw the signal abort at exactly 240000 ms of page time
// (Date.now(); performance.now() agrees within one 0.1 ms quantum of its coarsening), the page shows the timeout,
// nothing can be saved, a result arriving after the timer (at 245000 ms) is ignored, and the button
// was disabled while the operation was pending.
func TestPasskeyPageAbortsAPendingCeremony(t *testing.T) {
	b := startPasskeyBrowser(t)
	origin := passkeyPageOrigin(t)
	b64 := base64.RawURLEncoding.EncodeToString
	links := map[string]string{
		"create": ceremonyLink{V: 1, Op: "create", RP: "localhost", Origin: origin, Name: "mark-mac", Challenge: b64(make([]byte, 32))}.link(),
		"get": ceremonyLink{V: 1, Op: "get", RP: "localhost", Origin: origin, Request: "mr-0123456789abcdef", Expires: approveT0 + 600,
			Hash16: "0123456789abcdef", Challenge: b64(make([]byte, 32)), Allow: []string{b64(make([]byte, 32))}}.link(),
	}
	for _, op := range []string{"create", "get"} {
		for _, inject := range []string{"abort", "late"} {
			t.Run(op+"-"+inject, func(t *testing.T) {
				b.do(t, nil, "tab", map[string]any{"inject": inject})
				var view pageView
				b.do(t, &view, "open", map[string]any{"url": links[op]})
				if view.State != "ready" {
					t.Fatalf("the page is %q before the click", view.State)
				}
				dl := t.TempDir()
				var c pageCeremony
				b.do(t, &c, "ceremony", map[string]any{"virtualBudget": 250000, "download": dl})
				t.Logf("%s-%s: 250 s of page time in %d ms", op, inject, c.RealMs)
				if c.Pending.State != "pending" || !c.Pending.ButtonDisabled {
					t.Errorf("while pending: state %q, button disabled %v", c.Pending.State, c.Pending.ButtonDisabled)
				}
				o := c.After.Op
				if o == nil || o.Method != op || !o.Signal || !o.Aborted || o.AbortedAfterDate == nil || *o.AbortedAfterDate != 240000 || o.AbortedAfter == nil || math.Abs(*o.AbortedAfter-240000) > 0.15 {
					t.Fatalf("the operation: %s; want %s with a signal aborted at exactly 240000 ms", opString(o), op)
				}
				if inject == "late" && (o.ResolvedAfterDate == nil || *o.ResolvedAfterDate != 245000) {
					t.Fatalf("the late result was not delivered at 245000 ms: %s", opString(o))
				}
				a := c.After
				if a.State != "expired" || !a.StatusError || a.Status != "Timed out: no passkey answer within 4 minutes. Nothing can be saved from this attempt; try again." {
					t.Errorf("after the timer: state %q, error %v, status %q", a.State, a.StatusError, a.Status)
				}
				if files := worldFiles(t, dl); !a.SaveHidden || a.SaveHref != "" || a.Result != "" || a.Code != "" || c.SaveClicked || len(c.Downloads) != 0 || len(files) != 0 {
					t.Errorf("something can be saved after expiry: %+v, downloads %+v, files %q", a, c.Downloads, files)
				}
				if a.ButtonDisabled || len(a.Violations) != 0 || len(c.Errors) != 0 {
					t.Errorf("button disabled %v after the operation settled; CSP violations %q; page errors %q", a.ButtonDisabled, a.Violations, c.Errors)
				}
			})
		}
	}
}
