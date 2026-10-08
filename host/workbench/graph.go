package workbench

import "fmt"

// GraphEntry carries only the entry's already checked outbound edges.
type GraphEntry struct {
	Index int64
	Href  string
	Edges []EdgeView
}
type GraphNode struct {
	X         int
	Y         int
	Label     string
	Available bool
	Href      string
}
type GraphEdge struct {
	X1       int
	Y1       int
	X2       int
	Y2       int
	Relation string
}
type GraphView struct {
	Unavailable string
	Height      int
	Nodes       []GraphNode
	Edges       []GraphEdge
}

// LayoutGraph preserves entry order and target first-appearance order. Maps
// serve only as lookup indexes; output never depends on map iteration.
func LayoutGraph(entries []GraphEntry) GraphView {
	g := GraphView{}
	targets := map[string]int{}
	var objects []GraphNode
	for i, e := range entries {
		g.Nodes = append(g.Nodes, GraphNode{X: 16, Y: 16 + 36*i, Label: fmt.Sprintf("entry %d", e.Index), Available: true, Href: e.Href})
		for _, edge := range e.Edges {
			row, ok := targets[edge.Target]
			if !ok {
				row = len(objects)
				targets[edge.Target] = row
				objects = append(objects, GraphNode{X: 344, Y: 16 + 36*row, Label: edge.Target, Available: edge.Available, Href: edge.Href})
			}
			g.Edges = append(g.Edges, GraphEdge{X1: 296, Y1: 30 + 36*i, X2: 344, Y2: 30 + 36*row, Relation: edge.Relation})
		}
	}
	g.Nodes = append(g.Nodes, objects...)
	g.Height = 32 + 36*max(len(entries), len(objects))
	return g
}
