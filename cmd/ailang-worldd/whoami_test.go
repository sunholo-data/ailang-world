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
		{"json", []string{"--json"}, 200, 0, `{"episode":"mine","secondsLeft":61,"grants":[]}`},
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
			if tc.name == "json" && !strings.Contains(out.String(), `"episode":"mine"`) {
				t.Fatal(out.String())
			}
		})
	}
}
