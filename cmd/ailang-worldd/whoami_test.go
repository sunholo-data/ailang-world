package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type selfTransport func(*http.Request) (*http.Response, error)

func (f selfTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestWhoami(t *testing.T) {
	token := strings.Repeat("a", 64)
	t.Setenv("WORLD_SESSION", token)
	old := whoamiTransport
	t.Cleanup(func() { whoamiTransport = old })
	for _, tc := range []struct {
		name   string
		args   []string
		status int
		code   int
		body   string
	}{
		{"argv", []string{token}, 0, 1, ""},
		{"down", nil, 0, 10, ""},
		{"expired", nil, 401, 11, ""},
		{"json", []string{"--json"}, 200, 0, `{"episode":"mine","secondsLeft":61,"grants":[],"git":{"base":"abcdef0123","behind":3,"ahead":0,"dirty":false}}`},
		{"base", nil, 200, 0, `{"episode":"mine","grants":[],"git":{"base":"abcdef0123","base_date":"date","base_subject":"subject","behind":3}}`},
		{"git-error", nil, 200, 0, `{"episode":"mine","grants":[],"git":{"error":"timed out"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			whoamiTransport = selfTransport(func(r *http.Request) (*http.Response, error) {
				if r.Header.Get("Authorization") != "Bearer "+token || r.URL.Path != "/v1/session/self" {
					t.Fatal("incorrect authenticated route")
				}
				if tc.status == 0 {
					return nil, errors.New("unavailable")
				}
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body)), Header: http.Header{}}, nil
			})
			var out, err bytes.Buffer
			code := runWhoami("http://127.0.0.1:1", tc.args, &out, &err)
			if code != tc.code {
				t.Fatalf("code %d: %s", code, err.String())
			}
			all := out.String() + err.String()
			if strings.Contains(all, token) {
				t.Fatal("credential leaked")
			}
			if code == 10 || code == 11 {
				for _, part := range []string{"what:", "why:", "next:"} {
					if !strings.Contains(all, part) {
						t.Fatal("missing three-part diagnostic")
					}
				}
			}
			if tc.name == "base" && !strings.Contains(out.String(), "Base abcdef0 (date, \"subject\") · dev is 3 commits ahead of it") {
				t.Fatal(out.String())
			}
			if tc.name == "git-error" && !strings.Contains(out.String(), "Base unknown: timed out") {
				t.Fatal(out.String())
			}
			if tc.name == "json" && (!strings.Contains(out.String(), `"episode":"mine"`) || !strings.Contains(out.String(), `"behind":3`) || !strings.Contains(out.String(), `"ahead":0`)) {
				t.Fatal(out.String())
			}
		})
	}
}
