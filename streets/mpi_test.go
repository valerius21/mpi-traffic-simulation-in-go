package streets

import (
	mpi2 "github.com/sbromberger/gompi"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewMPI(t *testing.T) {
	t.Skip("Skipping test that requires MPI")

	mpi2.Start(false)
	defer mpi2.Stop()
	comm := mpi2.NewCommunicator(nil)
	taskID := comm.Rank()
	mpi := NewMPI(taskID, *comm, nil)
	if mpi == nil {
		t.Error("MPI is nil")
	}
	testType := reflect.TypeOf(mpi)
	if testType.String() != "*streets.MPI" {
		t.Error("MPI is not of type *streets.MPI")
	}
}

func TestEdgeRequest(t *testing.T) {
	t.Skip("Skipping test that requires MPI")

	mpi2.Start(false)
	defer mpi2.Stop()
	edges := make([]JEdge, 0)
	edge := JEdge{
		From:     2,
		To:       4,
		Length:   10,
		MaxSpeed: "10",
		Name:     "test",
		ID:       "asdf",
		Data:     Data{},
	}
	edges = append(edges, edge)
	leftVertex := JVertex{
		X:  1.,
		Y:  1.,
		ID: 2,
	}
	rightVertex := JVertex{
		X:  2.,
		Y:  1.,
		ID: 4,
	}
	vertices := make([]JVertex, 0)
	vertices = append(vertices, leftVertex)
	vertices = append(vertices, rightVertex)

	bRoot := NewGraphBuilder().WithEdges(edges).WithVertices(vertices).SetTopRightBottomLeftVertices()
	bRoot = bRoot.NumberOfRects(1).DivideGraphsIntoRects().PickRect(0).IsRoot()
	rootGraph, _ := bRoot.Build()

	leftB := NewGraphBuilder().WithEdges(edges).WithVertices(vertices).SetTopRightBottomLeftVertices()
	leftB = leftB.NumberOfRects(2).DivideGraphsIntoRects().PickRect(0).FilterForRect().IsLeaf(rootGraph, 0)

	rightB := NewGraphBuilder().WithEdges(edges).WithVertices(vertices).SetTopRightBottomLeftVertices()
	rightB = rightB.NumberOfRects(2).DivideGraphsIntoRects().PickRect(1).FilterForRect().IsLeaf(rootGraph, 0)

	leftGraph, _ := leftB.Build()
	_, _ = rightB.Build()

	if mpi2.WorldRank() == 0 {
		comm := mpi2.NewCommunicator(nil)
		mpi := NewMPI(0, *comm, rootGraph)
		err := mpi.RespondToEdgeLengthRequest()
		if err != nil {
			t.Error(err)
		}
	} else {
		taskID := mpi2.WorldRank()
		comm := mpi2.NewCommunicator([]int{taskID})
		mpi := NewMPI(taskID, *comm, leftGraph)
		length, err := mpi.AskRootForEdgeLength(2, 4)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, edge.Length, length)
	}
}
