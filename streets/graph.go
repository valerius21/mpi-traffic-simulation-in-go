package streets

import (
	"github.com/dominikbraun/graph"
)

// StreetGraph is a graph of streets with vertices of type int and edges of type JVertex
type StreetGraph struct {
	// ID is the ID of the graph. Root graph has ID 0
	ID string

	// RootGraph is the root graph of the graph, nil if the graph is the root graph
	RootGraph *StreetGraph

	// graph is the graph
	Graph graph.Graph[int, JVertex]
}

// vertexExists checks if a vertex exists in a graph
func (g *StreetGraph) vertexExists(v int) bool {
	_, err := g.Graph.Vertex(v)
	return err == nil
}
