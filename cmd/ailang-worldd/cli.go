package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

const (
	maxClientResponseBytes = 8388608
	maxClientCommitBytes   = maxClientResponseBytes
)

type client struct {
	base    string
	timeout time.Duration
	http    *http.Client
}

func newClient(base string) *client {
	return &client{base: strings.TrimSuffix(base, "/"), timeout: daemon.DefaultClientTimeout, http: &http.Client{}}
}

// do is the one transport path for every CLI call. Its caller supplies a
// context so cancellation propagates, and every request additionally receives
// the client's injectable D7 deadline.
func (c *client) do(ctx context.Context, method, path string, body io.Reader) (int, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, body)
	if err != nil {
		return 0, nil, fmt.Errorf("build %s %s%s: %w", method, c.base, path, err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("%s %s%s: %w", method, c.base, path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	limited := io.LimitReader(resp.Body, maxClientResponseBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("read %s%s response: %w", c.base, path, err)
	}
	if len(data) > maxClientResponseBytes {
		return resp.StatusCode, nil, fmt.Errorf("response from %s%s exceeds %d bytes", c.base, path, maxClientResponseBytes)
	}
	return resp.StatusCode, data, nil
}

func (c *client) get(path string) (int, string, error) {
	status, body, err := c.do(context.Background(), http.MethodGet, path, nil)
	return status, string(body), err
}

func reportResponse(path string, status int, body []byte, stdout, stderr io.Writer) int {
	if status < 200 || status >= 300 {
		var apiErr daemon.APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error.Class != "" {
			fmt.Fprintf(stderr, "ailang-worldd: %s returned HTTP %d %s: %s", path, status, apiErr.Error.Class, apiErr.Error.Message)
			if apiErr.Error.ObservedHead != "" || apiErr.Error.SelectedHead != "" {
				fmt.Fprintf(stderr, " (observedHead=%s selectedHead=%s)", apiErr.Error.ObservedHead, apiErr.Error.SelectedHead)
			}
			fmt.Fprintln(stderr)
		} else {
			fmt.Fprintf(stderr, "ailang-worldd: %s returned HTTP %d: %s\n", path, status, strings.TrimSpace(string(body)))
		}
		return exitUsage
	}
	fmt.Fprintln(stdout, strings.TrimSpace(string(body)))
	return exitOK
}

func execute(addr, method, path string, body io.Reader, stdout, stderr io.Writer) int {
	status, data, err := newClient(addr).do(context.Background(), method, path, body)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd: %v\n", err)
		return exitUsage
	}
	return reportResponse(path, status, data, stdout, stderr)
}

func runClientGet(addr, path string, args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintf(stderr, "ailang-worldd: unexpected argument %q\n", args[0])
		return exitUsage
	}
	return execute(addr, http.MethodGet, path, nil, stdout, stderr)
}

func requireGet(args []string, noun string, stdout, stderr io.Writer, build func(string) string, addr string) int {
	if len(args) != 2 || args[0] != "get" {
		fmt.Fprintf(stderr, "ailang-worldd %s: usage: %s get <%s>\n", noun, noun, noun)
		return exitUsage
	}
	return execute(addr, http.MethodGet, build(args[1]), nil, stdout, stderr)
}

func runWorld(addr string, args []string, stdout, stderr io.Writer) int {
	return requireGet(args, "world", stdout, stderr, func(ref string) string {
		return "/v1/worlds/" + url.PathEscape(ref)
	}, addr)
}

func runObject(addr string, args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 || args[0] != "get" {
		fmt.Fprintln(stderr, "ailang-worldd object: usage: object get <ref> [--payload]")
		return exitUsage
	}
	fs := flag.NewFlagSet("ailang-worldd object get", flag.ContinueOnError)
	fs.SetOutput(stderr)
	payload := fs.Bool("payload", false, "include base64 payload")
	if err := fs.Parse(args[2:]); err != nil || len(fs.Args()) != 0 {
		return exitUsage
	}
	path := "/v1/objects/" + url.PathEscape(args[1])
	if *payload {
		path += "?payload=true"
	}
	return execute(addr, http.MethodGet, path, nil, stdout, stderr)
}

func runLog(addr string, args []string, stdout, stderr io.Writer) int {
	if len(args) >= 1 && args[0] == "get" {
		return requireGet(args, "log", stdout, stderr, func(index string) string {
			return "/v1/log/" + url.PathEscape(index)
		}, addr)
	}
	if len(args) < 1 || args[0] != "range" {
		fmt.Fprintln(stderr, "ailang-worldd log: usage: log get <index> | log range --from N [--limit M]")
		return exitUsage
	}
	fs := flag.NewFlagSet("ailang-worldd log range", flag.ContinueOnError)
	fs.SetOutput(stderr)
	from := fs.Int64("from", -1, "first log index")
	limit := fs.Int("limit", 0, "maximum entries")
	if err := fs.Parse(args[1:]); err != nil || len(fs.Args()) != 0 || *from < 0 {
		if *from < 0 {
			fmt.Fprintln(stderr, "ailang-worldd log range: --from must be a non-negative integer")
		}
		return exitUsage
	}
	path := "/v1/log?from=" + strconv.FormatInt(*from, 10)
	if *limit != 0 {
		path += "&limit=" + strconv.Itoa(*limit)
	}
	return execute(addr, http.MethodGet, path, nil, stdout, stderr)
}

func runRegistry(addr string, args []string, stdout, stderr io.Writer) int {
	return requireGet(args, "registry", stdout, stderr, func(name string) string {
		// The route wildcard intentionally spans segments. Escape each segment,
		// preserving slashes rather than turning the semantic ID into one segment.
		parts := strings.Split(name, "/")
		for i := range parts {
			parts[i] = url.PathEscape(parts[i])
		}
		return "/v1/registry/" + strings.Join(parts, "/")
	}, addr)
}

func runCommit(addr string, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ailang-worldd commit", flag.ContinueOnError)
	fs.SetOutput(stderr)
	file := fs.String("file", "", "commit JSON file")
	if err := fs.Parse(args); err != nil || len(fs.Args()) != 0 || *file == "" {
		if *file == "" {
			fmt.Fprintln(stderr, "ailang-worldd commit: --file is required")
		}
		return exitUsage
	}
	f, err := os.Open(*file)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd commit: %v\n", err)
		return exitUsage
	}
	defer func() { _ = f.Close() }()
	data, err := io.ReadAll(io.LimitReader(f, maxClientCommitBytes+1))
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd commit: read %s: %v\n", *file, err)
		return exitUsage
	}
	if len(data) > maxClientCommitBytes {
		fmt.Fprintf(stderr, "ailang-worldd commit: %s exceeds %d bytes\n", *file, maxClientCommitBytes)
		return exitUsage
	}
	if len(bytes.TrimSpace(data)) == 0 {
		fmt.Fprintln(stderr, "ailang-worldd commit: commit file is empty")
		return exitUsage
	}
	return execute(addr, http.MethodPost, "/v1/commit", bytes.NewReader(data), stdout, stderr)
}
