package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
)

const (
	conversationSchema     = "world/surface-conversation/v1"
	maxConversationBytes   = 16 << 20
	maxConversationTurns   = 80
	maxMessageTextRunes    = 8000
	maxToolPillsPerMessage = 8
)

type surfaceConversation struct {
	Schema      string               `json:"schema"`
	Episode     string               `json:"episode"`
	SourceState string               `json:"sourceState"`
	Agent       string               `json:"agent,omitempty"`
	Status      string               `json:"status"` // "verified", "live", "unavailable"
	SHA256      string               `json:"sha256,omitempty"`
	Reason      string               `json:"reason,omitempty"`
	Total       int                  `json:"total"`
	Truncated   bool                 `json:"truncated,omitempty"`
	Messages    []conversationTurn   `json:"messages"`
}

type conversationMedia struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type conversationTurn struct {
	Role  string              `json:"role"` // "user" or "assistant"
	Text  string              `json:"text,omitempty"`
	Tools []string            `json:"tools,omitempty"`
	Media []conversationMedia `json:"media,omitempty"`
	At    string              `json:"at,omitempty"`
}

var (
	tokenScrubPattern  = regexp.MustCompile(`\b(?:wld_[A-Za-z0-9_-]{16,}|gh[pousr]_[A-Za-z0-9]{20,}|sk-[A-Za-z0-9_-]{20,})\b`)
	userUploadedFileRe = regexp.MustCompile(`\.user_uploaded/([A-Za-z0-9._-]+\.(?:png|jpe?g|gif|webp))`)
	safeMediaFileRe    = regexp.MustCompile(`^[A-Za-z0-9._-]+\.(?:png|jpe?g|gif|webp)$`)
)

func extractAgyMedia(raw, episode string) []conversationMedia {
	if raw == "" || !surface.EpisodeGrammar.MatchString(episode) {
		return nil
	}
	matches := userUploadedFileRe.FindAllStringSubmatch(raw, 8)
	if len(matches) == 0 {
		return nil
	}
	seen := map[string]bool{}
	var out []conversationMedia
	for _, m := range matches {
		name := m[1]
		if seen[name] {
			continue
		}
		seen[name] = true
		out = append(out, conversationMedia{
			Name: name,
			URL:  "/surface/conversation-media?ep=" + url.QueryEscape(episode) + "&file=" + url.QueryEscape(name),
		})
	}
	return out
}

func unpackAgyToolLabel(name string, rawArgs, rawAlt json.RawMessage) string {
	raw := rawArgs
	if len(bytes.TrimSpace(raw)) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		raw = rawAlt
	}
	var obj map[string]json.RawMessage
	if len(raw) > 0 {
		if raw[0] == '"' {
			var s string
			if json.Unmarshal(raw, &s) == nil {
				raw = []byte(s)
			}
		}
		_ = json.Unmarshal(raw, &obj)
	}
	if (name == "call_mcp_tool" || strings.HasPrefix(name, "mcp__")) && obj != nil {
		for _, k := range []string{"ToolName", "toolName", "tool_name"} {
			if v, ok := obj[k]; ok {
				var inner string
				if json.Unmarshal(v, &inner) == nil && strings.TrimSpace(inner) != "" {
					name = strings.TrimSpace(inner)
					break
				}
			}
		}
		for _, k := range []string{"Arguments", "arguments", "args"} {
			if v, ok := obj[k]; ok {
				if len(v) > 0 && v[0] == '"' {
					var s string
					if json.Unmarshal(v, &s) == nil {
						v = []byte(s)
					}
				}
				var innerObj map[string]json.RawMessage
				if json.Unmarshal(v, &innerObj) == nil && innerObj != nil {
					obj = innerObj
					break
				}
			}
		}
	}
	if i := strings.Index(name, "__"); strings.HasPrefix(name, "mcp__") && i > 0 {
		if rest := name[5:]; strings.Contains(rest, "__") {
			name = rest[strings.Index(rest, "__")+2:]
		}
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if obj != nil {
		var p, cmd string
		for _, k := range []string{"path", "Path"} {
			if v, ok := obj[k]; ok && json.Unmarshal(v, &p) == nil && strings.TrimSpace(p) != "" {
				p = strings.TrimPrefix(strings.TrimSpace(p), "./")
				if !strings.HasPrefix(p, "/") {
					return name + " " + p
				}
			}
		}
		for _, k := range []string{"command", "CommandLine"} {
			if v, ok := obj[k]; ok && json.Unmarshal(v, &cmd) == nil && strings.TrimSpace(cmd) != "" {
				return name + " " + truncateRunes(strings.TrimSpace(cmd), 28)
			}
		}
	}
	return name
}

func scrubConversationText(s string) string {
	return tokenScrubPattern.ReplaceAllString(s, "[redacted]")
}

func truncateRunes(s string, limit int) string {
	if utf8.RuneCountInString(s) <= limit {
		return s
	}
	var b strings.Builder
	count := 0
	for _, r := range s {
		if count >= limit {
			b.WriteString("\n…")
			break
		}
		b.WriteRune(r)
		count++
	}
	return b.String()
}

// cleanAgyUserContent extracts the human's actual <USER_REQUEST> text when present,
// stripping harness context summaries and system reminders.
func cleanAgyUserContent(raw string) string {
	if start := strings.Index(raw, "<USER_REQUEST>"); start >= 0 {
		rest := raw[start+len("<USER_REQUEST>"):]
		if end := strings.Index(rest, "</USER_REQUEST>"); end >= 0 {
			return strings.TrimSpace(rest[:end])
		}
		return strings.TrimSpace(rest)
	}
	if strings.Contains(raw, "<CONTEXT_SUMMARY>") {
		if end := strings.Index(raw, "</CONTEXT_SUMMARY>"); end >= 0 {
			after := strings.TrimSpace(raw[end+len("</CONTEXT_SUMMARY>"):])
			if after != "" {
				return after
			}
			return ""
		}
	}
	if cut := strings.Index(raw, "The following is a <SYSTEM_MESSAGE>"); cut >= 0 {
		raw = raw[:cut]
	}
	return strings.TrimSpace(raw)
}

func (d *Daemon) handleSurfaceConversation(w http.ResponseWriter, r *http.Request, _ []string) {
	refuse := func(status int, reason string) {
		d.writeSurfaceJSON(w, r, status, conversationSchema, surfaceReason{Schema: conversationSchema, Reason: reason})
	}
	if len(r.URL.RawQuery) > 2048 {
		refuse(http.StatusBadRequest, "expected a brane and state")
		return
	}
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(q) != 2 || len(q["brane"]) != 1 || len(q["at"]) != 1 {
		refuse(http.StatusBadRequest, "expected a brane and state")
		return
	}
	if d.follow == nil {
		refuse(http.StatusNotFound, "no verified episode for this state")
		return
	}
	snapshot := d.followSnapshot()
	env := snapshot.Envelope
	selection, through, err := selectSurfaceEpisode(env, q.Get("brane"), q.Get("at"), time.Now())
	if errors.Is(err, errSurfaceEpisodeAddress) {
		refuse(http.StatusBadRequest, errSurfaceEpisodeAddress.Error())
		return
	}
	if err != nil {
		refuse(http.StatusNotFound, "no verified episode for this state")
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	start := d.readFollowEpisodeStart(ctx, env, selection)
	if start == nil {
		refuse(http.StatusNotFound, "no verified episode for this state")
		return
	}
	out := d.readEpisodeConversation(ctx, start, selection.Front.ID, through)
	d.writeSurfaceJSON(w, r, http.StatusOK, conversationSchema, out)
}

func (d *Daemon) readEpisodeConversation(ctx context.Context, start *store.EpisodeStart, sourceState string, through time.Time) surfaceConversation {
	out := surfaceConversation{
		Schema:      conversationSchema,
		Episode:     start.Episode,
		SourceState: sourceState,
		Agent:       start.Agent,
		Status:      "unavailable",
		Messages:    []conversationTurn{},
	}
	if !surface.EpisodeGrammar.MatchString(start.Episode) || d.cfg.DBPath == "" {
		out.Reason = "episode transcript path unavailable"
		return out
	}
	worldDir := filepath.Dir(d.cfg.DBPath)

	// 1. Prefer the sealed session transcript when available.
	if len(start.Launches) > 0 {
		last := start.Launches[len(start.Launches)-1]
		if finalDigest, err := hashref.Parse(last.TranscriptSHA256); err == nil && (last.TranscriptState == "kept" || last.TranscriptState == "scrubbed") {
			if root := openSurfaceDir(filepath.Join(worldDir, "sessions")); root != nil {
				data := readSurfaceFile(root, start.Episode+".transcript.jsonl", maxConversationBytes)
				_ = root.Close()
				if len(data) > 0 && hashref.SumSHA256(data) == finalDigest && utf8.Valid(data) {
					turns := parseTranscriptTurns(ctx, data, start.Episode, through)
					return finalizeConversation(out, "verified", finalDigest.Digest(), turns)
				}
			}
		}
	}

	// 2. Otherwise read the live host transcript while the episode is active or not yet sealed.
	liveData := readLiveHostTranscript(worldDir, start)
	if len(liveData) > 0 && utf8.Valid(liveData) {
		turns := parseTranscriptTurns(ctx, liveData, start.Episode, through)
		return finalizeConversation(out, "live", "", turns)
	}
	out.Reason = "no transcript lines recorded yet for this episode"
	return out
}

func finalizeConversation(out surfaceConversation, status, digest string, turns []conversationTurn) surfaceConversation {
	merged := make([]conversationTurn, 0, len(turns))
	for _, t := range turns {
		if t.Role == "assistant" && len(merged) > 0 && merged[len(merged)-1].Role == "assistant" && merged[len(merged)-1].Text == "" {
			prev := &merged[len(merged)-1]
			prev.Tools = append(prev.Tools, t.Tools...)
			if t.Text != "" {
				prev.Text = t.Text
			}
			if t.At != "" {
				prev.At = t.At
			}
			continue
		}
		merged = append(merged, t)
	}
	out.Status = status
	out.SHA256 = digest
	out.Total = len(merged)
	if len(merged) > maxConversationTurns {
		out.Truncated = true
		merged = merged[len(merged)-maxConversationTurns:]
	}
	out.Messages = merged
	return out
}

func readRegularNoFollow(path string, maxBytes int64) []byte {
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() <= 0 || info.Size() > maxBytes {
		return nil
	}
	buf := make([]byte, info.Size())
	n, _ := f.Read(buf)
	return buf[:n]
}

func pathWithinRoot(candidate, root string) bool {
	cleanRoot := filepath.Clean(root)
	cleanCand := filepath.Clean(candidate)
	rel, err := filepath.Rel(cleanRoot, cleanCand)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func readLiveHostTranscript(worldDir string, start *store.EpisodeStart) []byte {
	ailangDir := filepath.Dir(worldDir)
	homeDir := filepath.Dir(ailangDir)
	if envHome, err := os.UserHomeDir(); err == nil && envHome != "" {
		// Allow tests whose worldDir is inside t.TempDir() to override via start worktree Layout,
		// falling back to envHome when worldDir is under ~/.ailang/world.
		if filepath.Base(ailangDir) == ".ailang" {
			homeDir = filepath.Dir(ailangDir)
		} else {
			homeDir = envHome
		}
	}
	candidates := []string{
		filepath.Join(worldDir, "launch", start.Episode),
		filepath.Join(ailangDir, "world-ws", start.Episode),
	}
	switch start.Agent {
	case "agy":
		brainRoot := filepath.Join(homeDir, ".gemini", "antigravity-cli", "brain")
		readAgy := func(txPath string) []byte {
			clean := filepath.Clean(txPath)
			if filepath.Base(clean) == "transcript.jsonl" {
				full := filepath.Join(filepath.Dir(clean), "transcript_full.jsonl")
				if data := readRegularNoFollow(full, maxConversationBytes); len(data) > 0 {
					return data
				}
			}
			return readRegularNoFollow(clean, maxConversationBytes)
		}
		for _, dir := range candidates {
			metaBytes := readRegularNoFollow(filepath.Join(dir, ".agents", ".world-session-meta.json"), 64<<10)
			if len(metaBytes) == 0 {
				continue
			}
			var meta struct {
				TranscriptPath string `json:"transcriptPath"`
			}
			if json.Unmarshal(metaBytes, &meta) != nil || meta.TranscriptPath == "" || !filepath.IsAbs(meta.TranscriptPath) {
				continue
			}
			if !pathWithinRoot(meta.TranscriptPath, brainRoot) {
				continue
			}
			if data := readAgy(meta.TranscriptPath); len(data) > 0 {
				return data
			}
		}
		if entries, err := os.ReadDir(brainRoot); err == nil {
			type brainCand struct {
				path string
				mod  time.Time
			}
			var cands []brainCand
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				txPath := filepath.Join(brainRoot, e.Name(), ".system_generated", "logs", "transcript.jsonl")
				if info, err := os.Lstat(txPath); err == nil && info.Mode().IsRegular() {
					cands = append(cands, brainCand{path: txPath, mod: info.ModTime()})
				}
			}
			sort.Slice(cands, func(i, j int) bool { return cands[i].mod.After(cands[j].mod) })
			if len(cands) > 20 {
				cands = cands[:20]
			}
			needle := []byte("/" + start.Episode)
			for _, c := range cands {
				if data := readRegularNoFollow(c.path, maxConversationBytes); len(data) > 0 && bytes.Contains(data, needle) {
					return readAgy(c.path)
				}
			}
		}
	case "claude":
		projectsRoot := filepath.Join(homeDir, ".claude", "projects")
		for _, dir := range candidates {
			real, err := filepath.EvalSymlinks(dir)
			if err != nil {
				continue
			}
			slugBytes := []byte(real)
			for i, c := range slugBytes {
				if !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9') {
					slugBytes[i] = '-'
				}
			}
			projDir := filepath.Join(projectsRoot, string(slugBytes))
			entries, err := os.ReadDir(projDir)
			if err != nil {
				continue
			}
			type candidateFile struct {
				path string
				mod  time.Time
			}
			var files []candidateFile
			for _, e := range entries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
					continue
				}
				info, statErr := e.Info()
				if statErr == nil && info.Mode().IsRegular() {
					files = append(files, candidateFile{path: filepath.Join(projDir, e.Name()), mod: info.ModTime()})
				}
			}
			sort.Slice(files, func(i, j int) bool { return files[i].mod.After(files[j].mod) })
			if len(files) > 0 {
				if data := readRegularNoFollow(files[0].path, maxConversationBytes); len(data) > 0 {
					return data
				}
			}
		}
	}
	return nil
}

func parseTranscriptTurns(ctx context.Context, data []byte, episode string, through time.Time) []conversationTurn {
	var turns []conversationTurn
	appendOrMergeAssistant := func(text string, tools []string, at string) {
		if text == "" && len(tools) == 0 {
			return
		}
		if len(turns) > 0 && turns[len(turns)-1].Role == "assistant" && turns[len(turns)-1].Text == "" && text == "" {
			prev := &turns[len(turns)-1]
			for _, t := range tools {
				if len(prev.Tools) < maxToolPillsPerMessage {
					prev.Tools = append(prev.Tools, t)
				}
			}
			if at != "" {
				prev.At = at
			}
			return
		}
		if len(turns) > 0 && turns[len(turns)-1].Role == "assistant" && turns[len(turns)-1].Text == "" && text != "" {
			prev := &turns[len(turns)-1]
			prev.Text = text
			for _, t := range tools {
				if len(prev.Tools) < maxToolPillsPerMessage {
					prev.Tools = append(prev.Tools, t)
				}
			}
			if at != "" {
				prev.At = at
			}
			return
		}
		if len(tools) > maxToolPillsPerMessage {
			tools = tools[:maxToolPillsPerMessage]
		}
		turns = append(turns, conversationTurn{Role: "assistant", Text: text, Tools: tools, At: at})
	}

	for offset := 0; offset < len(data); {
		if ctx.Err() != nil {
			break
		}
		line, rest, _ := bytes.Cut(data[offset:], []byte{'\n'})
		offset = len(data) - len(rest)
		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.HasPrefix(line, []byte(`{"world_launch":`)) {
			continue
		}
		var env struct {
			Type      string `json:"type"`
			Timestamp string `json:"timestamp"`
			CreatedAt string `json:"created_at"`
			Content   string `json:"content"`
			ToolCalls []struct {
				Name      string          `json:"name"`
				Arguments json.RawMessage `json:"arguments"`
				Args      json.RawMessage `json:"args"`
			} `json:"tool_calls"`
			Message struct {
				Role    string          `json:"role"`
				Content json.RawMessage `json:"content"`
			} `json:"message"`
			Payload struct {
				Type    string `json:"type"`
				Role    string `json:"role"`
				Message string `json:"message"`
				Name    string `json:"name"`
				Content []struct {
					Type string `json:"type"`
					Text string `json:"text"`
				} `json:"content"`
			} `json:"payload"`
		}
		if json.Unmarshal(line, &env) != nil {
			continue
		}
		at := env.Timestamp
		if at == "" {
			at = env.CreatedAt
		}
		if !through.IsZero() && at != "" {
			if stamp, err := time.Parse(time.RFC3339, at); err == nil && stamp.After(through) {
				continue
			}
		}

		switch env.Type {
		case "USER_INPUT": // Antigravity (agy)
			media := extractAgyMedia(env.Content, episode)
			cleaned := cleanAgyUserContent(env.Content)
			if cleaned != "" || len(media) > 0 {
				if cleaned == "" {
					cleaned = "(Uploaded image)"
				}
				turns = append(turns, conversationTurn{
					Role:  "user",
					Text:  truncateRunes(scrubConversationText(cleaned), maxMessageTextRunes),
					Media: media,
					At:    at,
				})
			}
		case "PLANNER_RESPONSE": // Antigravity (agy)
			var tools []string
			for _, tc := range env.ToolCalls {
				if label := unpackAgyToolLabel(tc.Name, tc.Arguments, tc.Args); label != "" {
					tools = append(tools, label)
				}
			}
			txt := strings.TrimSpace(env.Content)
			if txt != "" {
				txt = truncateRunes(scrubConversationText(txt), maxMessageTextRunes)
			}
			appendOrMergeAssistant(txt, tools, at)
		case "user", "assistant": // Claude Code
			text, tools := extractClaudeMessage(env.Message.Content)
			if env.Type == "user" {
				if text != "" {
					turns = append(turns, conversationTurn{
						Role: "user",
						Text: truncateRunes(scrubConversationText(text), maxMessageTextRunes),
						At:   at,
					})
				}
			} else {
				if text != "" {
					text = truncateRunes(scrubConversationText(text), maxMessageTextRunes)
				}
				appendOrMergeAssistant(text, tools, at)
			}
		case "event_msg": // Codex
			if env.Payload.Type == "user_message" && strings.TrimSpace(env.Payload.Message) != "" {
				turns = append(turns, conversationTurn{
					Role: "user",
					Text: truncateRunes(scrubConversationText(strings.TrimSpace(env.Payload.Message)), maxMessageTextRunes),
					At:   at,
				})
			} else if env.Payload.Type == "agent_message" && strings.TrimSpace(env.Payload.Message) != "" {
				appendOrMergeAssistant(truncateRunes(scrubConversationText(strings.TrimSpace(env.Payload.Message)), maxMessageTextRunes), nil, at)
			}
		case "response_item": // Codex
			if env.Payload.Type == "function_call" && env.Payload.Name != "" {
				appendOrMergeAssistant("", []string{env.Payload.Name}, at)
			} else if env.Payload.Role == "user" || env.Payload.Role == "assistant" {
				var parts []string
				for _, c := range env.Payload.Content {
					if (c.Type == "input_text" || c.Type == "output_text" || c.Type == "text") && strings.TrimSpace(c.Text) != "" {
						parts = append(parts, strings.TrimSpace(c.Text))
					}
				}
				joined := strings.Join(parts, "\n\n")
				if joined != "" {
					joined = truncateRunes(scrubConversationText(joined), maxMessageTextRunes)
					if env.Payload.Role == "user" {
						turns = append(turns, conversationTurn{Role: "user", Text: joined, At: at})
					} else {
						appendOrMergeAssistant(joined, nil, at)
					}
				}
			}
		}
	}
	return turns
}

func extractClaudeMessage(raw json.RawMessage) (string, []string) {
	if len(raw) == 0 {
		return "", nil
	}
	var plain string
	if json.Unmarshal(raw, &plain) == nil {
		return strings.TrimSpace(plain), nil
	}
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
		Name string `json:"name"`
	}
	if json.Unmarshal(raw, &blocks) != nil {
		return "", nil
	}
	var texts []string
	var tools []string
	for _, b := range blocks {
		switch b.Type {
		case "text":
			if t := strings.TrimSpace(b.Text); t != "" {
				texts = append(texts, t)
			}
		case "tool_use":
			if b.Name != "" {
				tools = append(tools, strings.TrimPrefix(b.Name, "mcp__world__"))
			}
		}
	}
	return strings.Join(texts, "\n\n"), tools
}

func (d *Daemon) handleSurfaceConversationMedia(w http.ResponseWriter, r *http.Request, _ []string) {
	q := r.URL.Query()
	ep := q.Get("ep")
	file := q.Get("file")
	if !surface.EpisodeGrammar.MatchString(ep) || !safeMediaFileRe.MatchString(strings.ToLower(file)) || d.cfg.DBPath == "" {
		http.Error(w, "invalid media request", http.StatusBadRequest)
		return
	}
	worldDir := filepath.Dir(d.cfg.DBPath)
	ailangDir := filepath.Dir(worldDir)
	homeDir := filepath.Dir(ailangDir)
	if filepath.Base(ailangDir) != ".ailang" {
		if envHome, err := os.UserHomeDir(); err == nil && envHome != "" {
			homeDir = envHome
		}
	}
	brainRoot := filepath.Join(homeDir, ".gemini", "antigravity-cli", "brain")
	var convDir string
	for _, dir := range []string{filepath.Join(worldDir, "launch", ep), filepath.Join(ailangDir, "world-ws", ep)} {
		metaBytes := readRegularNoFollow(filepath.Join(dir, ".agents", ".world-session-meta.json"), 64<<10)
		var meta struct {
			TranscriptPath string `json:"transcriptPath"`
		}
		if len(metaBytes) > 0 && json.Unmarshal(metaBytes, &meta) == nil && meta.TranscriptPath != "" && pathWithinRoot(meta.TranscriptPath, brainRoot) {
			convDir = filepath.Dir(filepath.Dir(filepath.Dir(filepath.Clean(meta.TranscriptPath))))
			break
		}
	}
	if convDir == "" {
		if entries, err := os.ReadDir(brainRoot); err == nil {
			needle := []byte("/" + ep)
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				txPath := filepath.Join(brainRoot, e.Name(), ".system_generated", "logs", "transcript.jsonl")
				if data := readRegularNoFollow(txPath, maxConversationBytes); len(data) > 0 && bytes.Contains(data, needle) {
					convDir = filepath.Join(brainRoot, e.Name())
					break
				}
			}
		}
	}
	if convDir == "" || !pathWithinRoot(convDir, brainRoot) {
		http.NotFound(w, r)
		return
	}
	imgBytes := readRegularNoFollow(filepath.Join(convDir, ".user_uploaded", file), 10<<20)
	if len(imgBytes) == 0 {
		http.NotFound(w, r)
		return
	}
	ext := strings.ToLower(filepath.Ext(file))
	mime := "image/jpeg"
	switch ext {
	case ".png":
		mime = "image/png"
	case ".gif":
		mime = "image/gif"
	case ".webp":
		mime = "image/webp"
	}
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Security-Policy", "default-src 'none'")
	w.Header().Set("Cache-Control", "private, max-age=300")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(imgBytes)
}
