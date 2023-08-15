package streets

import (
	"github.com/dominikbraun/graph"
	"golang.org/x/exp/slices"
	"math/rand"
	"pchpc_next/utils"
)

// StreetGraph is a graph of streets with vertices of type int and edges of type JVertex
type StreetGraph struct {
	// ID is the ID of the graph. Root graph has ID 0, leaf graphs have IDs 1, 2, 3, ...
	ID int

	// RootGraph is the root graph of the graph, nil if the graph is the root graph
	RootGraph *StreetGraph

	// graph is the graph
	Graph graph.Graph[int, JVertex]

	// vertex IDs
	vertexIDs []int
}

// VertexExists checks if a vertex exists in a graph
func (g *StreetGraph) VertexExists(v int) bool {
	_, err := g.Graph.Vertex(v)
	return err == nil
}

// GetVertices gets all vertex ids in a graph
func (g *StreetGraph) GetVertices() ([]int, error) {
	// Look for cached vertex IDs
	if g.vertexIDs != nil {
		return g.vertexIDs, nil
	}

	edges, err := g.Graph.Edges()
	if err != nil {
		return nil, err
	}

	vertices := make([]int, 0)
	for _, edge := range edges {
		dstID := edge.Target
		srcID := edge.Source

		if !slices.Contains(vertices, dstID) {
			vertices = append(vertices, dstID)
		}
		if !slices.Contains(vertices, srcID) {
			vertices = append(vertices, dstID, srcID)
		}
	}

	g.vertexIDs = vertices
	return g.vertexIDs, nil
}

// AddVehicle adds a vehicle to a graph
func (g *StreetGraph) AddVehicle(minSpeed, maxSpeed float64) (*Vehicle, error) {
	vertices, err := g.GetVertices()
	if err != nil {
		return nil, err
	}

	// Calculate the path
	var path []int
	for len(path) < 2 {
		srcIdx := rand.Intn(len(vertices))
		src := vertices[srcIdx]
		destIdx := rand.Intn(len(vertices))
		dest := vertices[destIdx]
		if src == dest {
			continue
		}
		path, err = graph.ShortestPath(g.Graph, src, dest)
		if err == nil {
			break
		}
	}

	speed := utils.RandomFloat64(minSpeed, maxSpeed)

	// Create the vehicle
	vb := NewVehicleBuilder().WithGraph(g).WithPathIDs(path).WithDelta(0.0).WithIsParked(false)
	vb = vb.WithSpeed(speed).WithLastID(path[0]).WithNextID(path[1])

	v, err := vb.Build()
	if err != nil {
		return nil, err
	}

	return &v, nil
}

// AddVehicleFromJson adds a vehicle to a graph from jsonBytes
func (g *StreetGraph) AddVehicleFromJson(jsonBytes []byte) (*Vehicle, error) {
	// Create the vehicle
	vb, err := NewVehicleBuilder().FromJsonBytes(jsonBytes)
	if err != nil {
		return nil, err
	}

	vb = vb.WithGraph(g)
	v, err := vb.Build()

	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (g *StreetGraph) GetRectFromVertexID(vertexID int, leafs []*StreetGraph) (int, error) {
	for _, leaf := range leafs {
		if leaf.VertexExists(vertexID) {
			return leaf.ID, nil
		}
	}
	return -1, nil
}
