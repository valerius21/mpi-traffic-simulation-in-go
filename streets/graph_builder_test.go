package streets

import (
	"github.com/aidarkhanov/nanoid"
	"github.com/stretchr/testify/assert"

	"os"
	"reflect"
	"testing"
)

func TestNewGraphBuilder(t *testing.T) {
	builder := NewGraphBuilder()
	if builder == nil {
		t.Error("Builder is nil")
	}
	testType := reflect.TypeOf(builder)
	if testType.String() != "*streets.GraphBuilder" {
		t.Error("Builder is not of type *streets.GraphBuilder")
	}
}

func TestGraphBuilder_WithVertices(t *testing.T) {
	vertices := make([]JVertex, 0)
	for i := 0; i < 10; i++ {
		vertices = append(vertices, JVertex{
			X:  1.,
			Y:  1.,
			ID: i,
		})
	}
	b := NewGraphBuilder().WithVertices(vertices)

	if len(b.vertices) != 10 {
		t.Errorf("Expected 10 vertices, got %d", len(b.vertices))
	}

	for i := 0; i < 10; i++ {
		if b.vertices[i].ID != vertices[i].ID {
			t.Errorf("Expected vertex to have ID %d, got %d", vertices[i].ID, b.vertices[i].ID)
		}
	}
}

func TestGraphBuilder_WithEdges(t *testing.T) {
	edges := make([]JEdge, 0)
	for i := 0; i < 10; i++ {
		id := nanoid.New()
		edges = append(edges, JEdge{
			From:     i,
			To:       i + 1,
			Length:   10,
			MaxSpeed: "10",
			Name:     "10 st.",
			ID:       id,
			Data:     Data{},
		})
	}
	b := NewGraphBuilder().WithEdges(edges)

	if len(b.edges) != 10 {
		t.Errorf("Expected 10 edges, got %d", len(b.edges))
	}

	for i := 0; i < 10; i++ {
		if b.edges[i].ID == "" {
			t.Errorf("Expected edge to have ID, got empty string")
		}
		sameIds := b.edges[i].ID == edges[i].ID
		if !sameIds {
			t.Errorf("Expected edge to have ID %s, got %s", edges[i].ID, b.edges[i].ID)
		}
	}
}

func TestGraphBuilder_WithRectangleParts(t *testing.T) {
	builder := NewGraphBuilder()

	builder.WithRectangleParts(5)

	assert.Equal(t, 5, builder.rectangleParts)
}

func TestGraphBuilder_FromJsonBytes(t *testing.T) {
	jBytes, err := os.ReadFile("../assets/out.json")
	if err != nil {
		t.Error(err)
	}
	if len(jBytes) == 0 {
		t.Error("JBytes is empty")
	}

	jGraph, err := UnmarshalGraphJSON(jBytes)
	if err != nil {
		t.Error(err)
	}

	b := NewGraphBuilder().FromJsonBytes(jBytes)
	if b == nil {
		t.Error("Builder is nil")
	}

	if len(b.vertices) != len(jGraph.Graph.Vertices) {
		t.Errorf("Expected %d vertices, got %d", len(jGraph.Graph.Vertices), len(b.vertices))
	}

	if len(b.edges) != len(jGraph.Graph.Edges) {
		t.Errorf("Expected %d edges, got %d", len(jGraph.Graph.Edges), len(b.edges))
	}

	for i := 0; i < len(b.vertices); i++ {
		if b.vertices[i].ID != jGraph.Graph.Vertices[i].ID {
			t.Errorf("Expected vertex to have ID %d, got %d", jGraph.Graph.Vertices[i].ID, b.vertices[i].ID)
		}
	}
}

func TestGraphBuilder_FromJsonFile(t *testing.T) {
	b := NewGraphBuilder().FromJsonFile("../assets/out.json")
	if b == nil {
		t.Error("Builder is nil")
	}
	if len(b.vertices) == 0 {
		t.Error("Vertices are empty")
	}
	if len(b.edges) == 0 {
		t.Error("Edges are empty")
	}
}

func TestGraphBuilder_SetTopRightBottomLeftVertices(t *testing.T) {
	vertices := make([]JVertex, 0)
	bottomVertex := JVertex{
		X:  0.,
		Y:  0.,
		ID: 0,
	}
	topVertex := JVertex{
		X:  1.,
		Y:  1.,
		ID: 1,
	}

	vertices = append(vertices, bottomVertex)
	vertices = append(vertices, topVertex)

	b := NewGraphBuilder().WithVertices(vertices).SetTopRightBottomLeftVertices()
	if b == nil {
		t.Error("Builder is nil")
	}

	if b.top.Y != topVertex.Y {
		t.Errorf("Expected top vertex to have Y %f, got %f", topVertex.Y, b.top.Y)
	}
	if b.top.X != topVertex.X {
		t.Errorf("Expected top vertex to have X %f, got %f", topVertex.X, b.top.X)
	}
	if b.bot.Y != bottomVertex.Y {
		t.Errorf("Expected bottom vertex to have Y %f, got %f", bottomVertex.Y, b.bot.Y)
	}
	if b.bot.X != bottomVertex.X {
		t.Errorf("Expected bottom vertex to have X %f, got %f", bottomVertex.X, b.bot.X)
	}
}

func TestGraphBuilder_NumberOfRects(t *testing.T) {
	n := 3
	b := NewGraphBuilder().NumberOfRects(n)
	if b == nil {
		t.Error("Builder is nil")
	}
	if b.rectangleParts != n {
		t.Errorf("Expected number of rects to be %d, got %d", n, b.rectangleParts)
	}
}

func TestGraphBuilder_DivideGraphsIntoRects(t *testing.T) {
	builder := NewGraphBuilder()
	builder.WithRectangleParts(2).WithVertices([]JVertex{
		{X: 0, Y: 0},
		{X: 4, Y: 4},
		{X: 3, Y: 3},
	}).SetTopRightBottomLeftVertices()

	builder.DivideGraphsIntoRects()

	assert.Equal(t, 2, len(builder.rects))
	assert.Equal(t, float64(2), builder.rects[0].TopRight.X)
	assert.Equal(t, float64(0), builder.rects[0].BotLeft.X)
}

func TestGraphBuilder_PickRect(t *testing.T) {
	builder := NewGraphBuilder()
	builder.WithRectangleParts(2).WithVertices([]JVertex{
		{X: 0, Y: 0},
		{X: 4, Y: 4},
		{X: 3, Y: 3},
	}).SetTopRightBottomLeftVertices().DivideGraphsIntoRects()

	builder.PickRect(1)

	assert.Equal(t, float64(4), builder.pickedRect.TopRight.X)
	assert.Equal(t, float64(2), builder.pickedRect.BotLeft.X)
}

func TestGraphBuilder_FilterForRect(t *testing.T) {
	gb := NewGraphBuilder()

	// Predefined data for the test
	vertices := []JVertex{
		{ID: 1, X: 1, Y: 1},
		{ID: 2, X: 2, Y: 2},
		{ID: 3, X: 4, Y: 4},
		{ID: 4, X: 5, Y: 5},
	}

	edges := []JEdge{
		{From: 1, To: 2},
		{From: 3, To: 4},
	}

	gb.WithVertices(vertices).WithEdges(edges).WithRectangleParts(2).SetTopRightBottomLeftVertices().DivideGraphsIntoRects().PickRect(0).FilterForRect()

	// Expectations
	expectedVertices := 2 // As both {1,1} and {2,2} lie within the picked rectangle
	expectedEdges := 1    // Only the edge {From: 1, To: 2} lies within the picked rectangle

	// Assertions
	if len(gb.vertices) != expectedVertices {
		t.Errorf("Expected %d vertices in the filtered rectangle, but got %d", expectedVertices, len(gb.vertices))
	}

	if len(gb.edges) != expectedEdges {
		t.Errorf("Expected %d edges in the filtered rectangle, but got %d", expectedEdges, len(gb.edges))
	}
}

func TestGraphBuilder_IsRoot(t *testing.T) {
	b := NewGraphBuilder().IsRoot()
	if b == nil {
		t.Error("Builder is nil")
	}
	if b.id < 0 {
		t.Error("Expected builder to be root")
	}
	if b.root != nil {
		t.Error("Expected root to be nil")
	}
}

func TestGraphBuilder_IsLeaf(t *testing.T) {
	root := StreetGraph{}

	b := NewGraphBuilder().IsLeaf(&root, 0)

	if b == nil {
		t.Error("Builder is nil")
	}

	if b.id < 0 {
		t.Error("Expected builder to have ID")
	}

	if b.root != &root {
		t.Error("Expected root to be nil")
	}
}

func TestGraphBuilder_Build(t *testing.T) {
	builder := NewGraphBuilder()
	builder.WithRectangleParts(2).WithVertices([]JVertex{
		{ID: 1, X: 1, Y: 1},
		{ID: 2, X: 3, Y: 3},
	}).WithEdges([]JEdge{
		{From: 1, To: 2},
	}).SetTopRightBottomLeftVertices().DivideGraphsIntoRects().PickRect(1).FilterForRect().IsRoot()

	g, err := builder.Build()

	assert.NoError(t, err)
	assert.NotNil(t, g)
}
