package main

// Row 201 A3 sprint 4 (design §2.9, §6 XA4): the opt-in notification in a live Chrome, through the
// driver's notify recorder (nothing is shown) and a CDP permission grant.

import (
	"reflect"
	"testing"
	"time"
)

// notifyView is what the notify recorder saw, the permission, and the opt-in's text.
type notifyView struct {
	Notes      []struct{ Title, Body, Tag string }
	Asks       []struct{ Active bool }
	Permission string
	Control    string
	Violations []string
}

// XA4: never asked, the page raises nothing and never requests permission, even for a new request.
// The browser grants notifications and "Notify me" asks once, from the click (a user gesture). A
// request the tab has already seen raises nothing; a new one raises exactly one "World needs you"
// (tag = its id, the body names the episode, its grants and its expiry), and the next refresh none.
// Reloaded with notifications granted, the page raises nothing for what it shows; a new request's
// notification, clicked, focuses the tab and opens its card.
func TestNeedsYouNotifiesOnlyWhenAskedOnce(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"inject": "notify"})
	now := time.Now().Unix()
	a := newApproveRig(t, "ep1", now)
	origin := rigSurface(t, a.credPath)
	seen := func(notes, asks int, permission, control string) notifyView {
		t.Helper()
		var v notifyView
		b.do(t, &v, "waitFor", map[string]any{"expression": `(() => { const c = document.querySelector('[data-needs] #notify-me, [data-needs] [data-notify]');
			return { notes: window.__notes, asks: window.__asks, permission: Notification.permission, control: c ? c.textContent : '', violations: window.__cspv }; })()`})
		if len(v.Notes) != notes || len(v.Asks) != asks || v.Permission != permission || v.Control != control || len(v.Violations) != 0 {
			t.Fatalf("the page made %d notifications %+v and %d permission requests, permission %q, control %q, CSP violations %q; want %d, %d, %q, %q, none",
				len(v.Notes), v.Notes, len(v.Asks), v.Permission, v.Control, v.Violations, notes, asks, permission, control)
		}
		return v
	}
	poll := func(prefix string) {
		t.Helper()
		pollSurface(t, b)
		sectionText(t, b, "data-needs", prefix)
	}
	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
	sectionText(t, b, "data-needs", "Needs you · 1")
	seen(0, 0, "default", "Notify me")
	newApproveRigIn(t, a.r, "ep2", now)
	poll("Needs you · 2")
	seen(0, 0, "default", "Notify me")

	var permission string
	if b.do(t, &permission, "grantNotifications", map[string]any{"origin": origin}); permission != "granted" {
		t.Fatalf("after the grant the permission is %q", permission)
	}
	b.do(t, nil, "tap", map[string]any{"id": "notify-me"})
	b.do(t, nil, "waitFor", map[string]any{"expression": `document.querySelector('[data-needs] [data-notify]') !== null`})
	if v := seen(0, 1, "granted", "Notifications on"); !v.Asks[0].Active {
		t.Fatal("permission was requested without a user gesture")
	}
	poll("Needs you · 2")
	seen(0, 1, "granted", "Notifications on")
	third := newApproveRigIn(t, a.r, "ep3", now)
	poll("Needs you · 3")
	v := seen(1, 1, "granted", "Notifications on")
	want := struct{ Title, Body, Tag string }{"World needs you", "Approve mint request for ep3 (10 grants), expires " + time.Unix(third.row(t).ExpiresAt, 0).UTC().Format("15:04"), third.id}
	if !reflect.DeepEqual(v.Notes[0], want) {
		t.Fatalf("the notification is %+v, want %+v", v.Notes[0], want)
	}
	poll("Needs you · 3")
	seen(1, 1, "granted", "Notifications on")

	b.do(t, nil, "openCard", map[string]any{"url": origin + "/surface", "raw": true})
	sectionText(t, b, "data-needs", "Needs you · 3")
	seen(0, 0, "granted", "Notifications on")
	fourth := newApproveRigIn(t, a.r, "ep4", now)
	poll("Needs you · 4")
	if v := seen(1, 0, "granted", "Notifications on"); v.Notes[0].Tag != fourth.id {
		t.Fatalf("the notification is %+v, want tag %s", v.Notes[0], fourth.id)
	}
	var clicked struct {
		Location string
		Focused  int
	}
	if b.do(t, &clicked, "notification", map[string]any{"index": 0}); clicked.Location != origin+"/surface/approve/"+fourth.id || clicked.Focused != 1 {
		t.Fatalf("the notification's click: %+v", clicked)
	}
	b.do(t, nil, "waitFor", map[string]any{"expression": `document.body.dataset.state === 'ready'`})
}
