package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const (
	maxSurfaceUsageBytes    = 8 << 20
	maxSurfaceUsageLaunches = 16
	maxSurfaceUsageEntries  = 32
)

type surfaceUsageMemo struct {
	mu      sync.Mutex
	parsing bool
	entries map[surfaceUsageEpisode]surfaceUsageEntry
}

type surfaceUsageEpisode struct {
	repo, episode, agent string
}

type surfaceUsageFileVersion struct {
	info      os.FileInfo
	changed   int64
	cacheable bool
}

func (v surfaceUsageFileVersion) same(other surfaceUsageFileVersion) bool {
	return v.cacheable && other.cacheable && v.changed == other.changed &&
		os.SameFile(v.info, other.info) && v.info.Size() == other.info.Size() &&
		v.info.ModTime().Equal(other.info.ModTime())
}

func surfaceUsageVersion(root *os.Root, name string) surfaceUsageFileVersion {
	info, err := root.Lstat(name)
	if err != nil || !info.Mode().IsRegular() || info.Size() <= 0 || info.Size() > maxSurfaceUsageBytes {
		return surfaceUsageFileVersion{}
	}
	changed, known := surfaceUsageChangeTime(info)
	return surfaceUsageFileVersion{info: info, changed: changed, cacheable: known}
}

type surfaceUsageEntry struct {
	file    surfaceUsageFileVersion
	digests [maxSurfaceUsageLaunches]hashref.HashRef
	turns   *int
	costUSD *float64
	tokens  tokenTotals
}

func (m *surfaceUsageMemo) begin(key surfaceUsageEpisode, file surfaceUsageFileVersion, digests [maxSurfaceUsageLaunches]hashref.HashRef) (entry surfaceUsageEntry, parse bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if cached, ok := m.entries[key]; ok && cached.digests == digests && cached.file.same(file) {
		return cached, false
	}
	if m.parsing {
		return surfaceUsageEntry{}, false
	}
	m.parsing = true
	return surfaceUsageEntry{}, true
}

func (m *surfaceUsageMemo) finish(key surfaceUsageEpisode, entry surfaceUsageEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.parsing = false
	if !entry.file.cacheable || (entry.turns == nil && entry.costUSD == nil && entry.tokens == (tokenTotals{})) {
		return
	}
	if m.entries == nil {
		m.entries = make(map[surfaceUsageEpisode]surfaceUsageEntry)
	}
	if _, exists := m.entries[key]; !exists && len(m.entries) == maxSurfaceUsageEntries {
		for k := range m.entries {
			delete(m.entries, k)
			break
		}
	}
	m.entries[key] = entry
}

func (e surfaceUsageEntry) usage() (turns *int, costUSD *float64, tokens *surfacepage.NodeTokens) {
	// Keep the cache's values private from callers.
	if e.turns != nil {
		n := *e.turns
		turns = &n
	}
	if e.costUSD != nil {
		c := *e.costUSD
		costUSD = &c
	}
	return turns, costUSD, e.tokens.project()
}

// surfaceRunUsage derives numbers from the service CLI's retained headless stream.
// The fixed leaf and each cumulative launch digest bind the bytes to the operator's
// launch records. Only metrics common to all launches are returned.
func (d *Daemon) surfaceRunUsage(start *store.EpisodeStart) (*int, *float64, *surfacepage.NodeTokens) {
	if start == nil || !surface.EpisodeGrammar.MatchString(start.Episode) || d.cfg.DBPath == "" ||
		(start.Agent != "claude" && start.Agent != "codex" && start.Agent != "agy") ||
		len(start.Launches) == 0 || len(start.Launches) > maxSurfaceUsageLaunches {
		return nil, nil, nil
	}
	var digests [maxSurfaceUsageLaunches]hashref.HashRef
	for i, launch := range start.Launches {
		digest, err := hashref.Parse(launch.TranscriptSHA256)
		if launch.Mode != "headless" || launch.Exit != 0 || launch.FinishReason != "exit" ||
			launch.TranscriptState != "kept" || err != nil {
			return nil, nil, nil
		}
		digests[i] = digest
	}
	root := openSurfaceDir(filepath.Join(filepath.Dir(d.cfg.DBPath), "sessions"))
	if root == nil {
		return nil, nil, nil
	}
	defer root.Close()
	name := start.Episode + ".transcript.jsonl"
	file := surfaceUsageVersion(root, name)
	if file.info == nil {
		return nil, nil, nil
	}
	key := surfaceUsageEpisode{repo: start.Repo, episode: start.Episode, agent: start.Agent}
	entry, parse := d.surfaceUsage.begin(key, file, digests)
	if !parse {
		return entry.usage()
	}

	// Metering is optional, so a cache miss gets one nonblocking parser slot.
	// Cache only after the retained file and launch record have both agreed.
	defer func() { d.surfaceUsage.finish(key, entry) }()
	data := readSurfaceFile(root, name, maxSurfaceUsageBytes)
	count := len(start.Launches)
	if len(data) == 0 || data[len(data)-1] != '\n' || !utf8.Valid(data) ||
		hashref.SumSHA256(data) != digests[count-1] ||
		(file.cacheable && !file.same(surfaceUsageVersion(root, name))) {
		return nil, nil, nil
	}
	turns, costUSD, tokens := surfaceLaunchTotals(start.Agent, data, digests[:count])
	entry = surfaceUsageEntry{file: file, digests: digests, turns: turns, costUSD: costUSD, tokens: tokens}
	return entry.usage()
}

func surfaceLaunchTotals(agent string, data []byte, digests []hashref.HashRef) (turns *int, costUSD *float64, tokens tokenTotals) {
	var totalTurns int
	var totalCost float64
	allTurns, allCosts := true, true
	begin, offset, launch := 0, 0, 0
	addRun := func(run []byte) bool {
		n, dollars, counters, complete := surfaceHeadlessUsage(agent, run)
		tokens.accumulate(counters, launch == 0)
		if n == nil || *n < 0 || *n > math.MaxInt-totalTurns {
			allTurns = false
		} else {
			totalTurns += *n
		}
		if dollars == nil || *dollars < 0 {
			allCosts = false
		} else {
			totalCost += *dollars
		}
		return complete
	}
	for offset < len(data) {
		line, _, _ := bytes.Cut(data[offset:], []byte{'\n'})
		next := offset + len(line) + 1
		if bytes.HasPrefix(line, []byte(`{"world_launch":`)) {
			if launch+1 >= len(digests) || string(line) != fmt.Sprintf(`{"world_launch": %d}`, launch+2) ||
				hashref.SumSHA256(data[:offset]) != digests[launch] || !addRun(data[begin:offset]) {
				return nil, nil, tokenTotals{}
			}
			launch++
			begin = next
		}
		offset = next
	}
	if launch+1 != len(digests) || !addRun(data[begin:]) {
		return nil, nil, tokenTotals{}
	}
	if allTurns {
		turns = &totalTurns
	}
	if allCosts && !math.IsInf(totalCost, 0) {
		costUSD = &totalCost
	}
	return turns, costUSD, tokens
}

// surfaceHeadlessUsage validates the host's launch boundaries while decoding only
// the event discriminator and final usage fields. Message and tool data stay private.
func surfaceHeadlessUsage(agent string, stream []byte) (turns *int, costUSD *float64, tokens tokenTotals, complete bool) {
	var started, inTurn, finished bool
	var completedTurns int
	for len(stream) > 0 {
		line, rest, hasLine := bytes.Cut(stream, []byte{'\n'})
		stream = rest
		var event struct {
			Type    string          `json:"type"`
			Subtype string          `json:"subtype"`
			Turns   *int            `json:"num_turns"`
			Cost    *float64        `json:"total_cost_usd"`
			IsError *bool           `json:"is_error"`
			Usage   json.RawMessage `json:"usage"`
		}
		if !hasLine || finished || json.Unmarshal(line, &event) != nil {
			return nil, nil, tokenTotals{}, false
		}
		if agent == "codex" {
			switch event.Type {
			case "thread.started":
				if started {
					return nil, nil, tokenTotals{}, false
				}
				started = true
			case "turn.started":
				if !started || inTurn {
					return nil, nil, tokenTotals{}, false
				}
				inTurn = true
			case "turn.completed":
				if !inTurn {
					return nil, nil, tokenTotals{}, false
				}
				// The exec event carries launch-local cumulative counters.
				// Use its latest values, including fields it omits.
				tokens = decodeRunTokens(agent, event.Usage)
				completedTurns++
				inTurn = false
			case "item.started", "item.updated", "item.completed":
				if !inTurn {
					return nil, nil, tokenTotals{}, false
				}
			default:
				return nil, nil, tokenTotals{}, false
			}
			continue
		}
		switch event.Type {
		case "system":
			if event.Subtype == "init" {
				if started {
					return nil, nil, tokenTotals{}, false
				}
				started = true
			}
		case "assistant", "user", "stream_event", "tool_progress", "tool_use_summary", "rate_limit_event", "prompt_suggestion":
			// These events can carry content, but do not contribute a usage number.
		case "result":
			if !started || event.Subtype != "success" || event.IsError == nil || *event.IsError {
				return nil, nil, tokenTotals{}, false
			}
			turns, costUSD, finished = event.Turns, event.Cost, true
			tokens = decodeRunTokens(agent, event.Usage)
		default:
			return nil, nil, tokenTotals{}, false
		}
	}
	if agent == "codex" && started && !inTurn && completedTurns > 0 {
		return &completedTurns, nil, tokens, true
	}
	return turns, costUSD, tokens, finished
}

// Counters are serialized as JSON numbers and must remain exact in the browser.
const maxTokenCount = 1<<53 - 1

const (
	tokenInput = iota
	tokenOutput
	tokenCacheRead
	tokenCacheWrite
	tokenReasoning
	tokenFieldCount
)

type tokenCounter struct {
	count    int64
	reported bool
}

func (c tokenCounter) plus(other tokenCounter) tokenCounter {
	if c.reported && other.reported && other.count <= maxTokenCount-c.count {
		return tokenCounter{count: c.count + other.count, reported: true}
	}
	return tokenCounter{}
}

func (c tokenCounter) number() *int64 {
	if c.reported {
		return &c.count
	}
	return nil
}

type tokenTotals [tokenFieldCount]tokenCounter

func (t *tokenTotals) accumulate(next tokenTotals, first bool) {
	if first {
		*t = next
		return
	}
	for i := range t {
		t[i] = t[i].plus(next[i])
	}
}

func (t tokenTotals) project() *surfacepage.NodeTokens {
	if t == (tokenTotals{}) {
		return nil
	}
	return &surfacepage.NodeTokens{Input: t[tokenInput].number(), Output: t[tokenOutput].number(),
		CacheRead: t[tokenCacheRead].number(), CacheWrite: t[tokenCacheWrite].number(), Reasoning: t[tokenReasoning].number()}
}

func decodeRunTokens(agent string, raw json.RawMessage) tokenTotals {
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil || fields == nil {
		return tokenTotals{}
	}
	counter := func(key string) tokenCounter {
		var n *int64
		if json.Unmarshal(fields[key], &n) != nil || n == nil || *n < 0 || *n > maxTokenCount {
			return tokenCounter{}
		}
		return tokenCounter{count: *n, reported: true}
	}
	input, output := counter("input_tokens"), counter("output_tokens")
	switch agent {
	case "codex":
		cacheRead := tokenSubset(counter("cached_input_tokens"), input)
		cacheWrite := tokenSubset(counter("cache_write_input_tokens"), input)
		reasoning := tokenSubset(counter("reasoning_output_tokens"), output)
		return tokenTotals{input, output, cacheRead, cacheWrite, reasoning}
	case "claude":
		// Claude's input_tokens excludes cache reads and cache creation.
		cacheRead, cacheWrite := counter("cache_read_input_tokens"), counter("cache_creation_input_tokens")
		return tokenTotals{input.plus(cacheRead).plus(cacheWrite), output, cacheRead, cacheWrite, {}}
	case "agy":
		cacheRead := counter("cache_read_input_tokens")
		if !cacheRead.reported {
			cacheRead = counter("cache_read_tokens")
		}
		cacheWrite := counter("cache_creation_input_tokens")
		if cacheWrite.reported {
			return tokenTotals{input.plus(cacheRead).plus(cacheWrite), output, cacheRead, cacheWrite, {}}
		}
		if cacheRead.reported {
			return tokenTotals{input.plus(cacheRead), output, cacheRead, {}, {}}
		}
		return tokenTotals{input, output, {}, {}, {}}
	default:
		return tokenTotals{}
	}
}

func tokenSubset(part, whole tokenCounter) tokenCounter {
	if whole.reported && part.reported && part.count > whole.count {
		return tokenCounter{}
	}
	return part
}
