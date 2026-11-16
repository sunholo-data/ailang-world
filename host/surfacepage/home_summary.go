package surfacepage

import (
	"strconv"
	"strings"
)

// HomeSummarySchema is the Studio start page's card contract (the sunholo-platform homepage's
// `summary` registry field): at most four short labelled values, read over loopback every 30 s.
const HomeSummarySchema = "sunholo/home-summary/v1"

// HomeItem is one value on the start page's card. Tone is "attention", "ok" or empty.
type HomeItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Tone  string `json:"tone,omitempty"`
}

// HomeSummary is GET /surface/home-summary.json.
type HomeSummary struct {
	Schema string     `json:"schema"`
	Items  []HomeItem `json:"items"`
}

// BuildHomeSummary gives the start page the app's own Now-card counts: Needs you (pending mint
// requests plus pending decisions, counted as the card counts them; needs < 0 when they could not
// be read), then, when a model was built, Running, Landed today and Health. It never carries an
// episode name, path or id: the card links to /surface for those.
func BuildHomeSummary(needs int, edge *NowEdge) HomeSummary {
	items := []HomeItem{}
	switch {
	case needs < 0:
		items = append(items, HomeItem{Label: "Needs you", Value: "unknown", Tone: "attention"})
	case needs > 0:
		items = append(items, HomeItem{Label: "Needs you", Value: strconv.Itoa(needs), Tone: "attention"})
	default:
		items = append(items, HomeItem{Label: "Needs you", Value: "0"})
	}
	if edge == nil {
		return HomeSummary{Schema: HomeSummarySchema, Items: items}
	}
	running := strconv.Itoa(len(edge.Sessions))
	if edge.Open != nil {
		running = strconv.Itoa(len(edge.Open.Live))
	} else if edge.SessionsUnavailable != "" {
		running = "unknown"
	}
	items = append(items,
		HomeItem{Label: "Running", Value: running},
		HomeItem{Label: "Landed today", Value: strconv.Itoa(edge.Today.Landed)})
	if strings.Contains(strings.ToLower(edge.Health), "healthy") && (edge.Service == nil || !edge.Service.Warn) {
		items = append(items, HomeItem{Label: "Health", Value: "healthy", Tone: "ok"})
	} else {
		items = append(items, HomeItem{Label: "Health", Value: "check", Tone: "attention"})
	}
	return HomeSummary{Schema: HomeSummarySchema, Items: items}
}
