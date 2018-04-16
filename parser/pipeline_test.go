package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTwoStepPipeline(t *testing.T) {
	stepA := Step{"a", "a", "a", "a", []string{}, []int{}, []int{1}}
	stepB := Step{"b", "b", "b", "b", []string{}, []int{0}, []int{}}
	pipeline := Pipeline{stepA, stepB}

	require.True(t, pipeline.Connected())
	require.Len(t, pipeline.Components(), 1)
	require.True(t, pipeline.Acyclic())
	require.True(t, pipeline.Valid())
}

func TestTwoStepCycle(t *testing.T) {
	stepA := Step{"a", "a", "a", "a", []string{}, []int{1}, []int{1}}
	stepB := Step{"b", "b", "b", "b", []string{}, []int{0}, []int{0}}
	pipeline := Pipeline{stepA, stepB}

	require.True(t, pipeline.Connected())
	require.Len(t, pipeline.Components(), 1)
	require.False(t, pipeline.Acyclic())
	require.False(t, pipeline.Valid())
}

func TestUnconnectedSteps(t *testing.T) {
	stepA := Step{"a", "a", "a", "a", []string{}, []int{}, []int{}}
	stepB := Step{"b", "b", "b", "b", []string{}, []int{}, []int{}}
	pipeline := Pipeline{stepA, stepB}

	require.False(t, pipeline.Connected())
	require.Len(t, pipeline.Components(), 2)
	require.True(t, pipeline.Acyclic())
	require.False(t, pipeline.Valid())
}

func TestBifurcatedPipeline(t *testing.T) {
	stepA := Step{"a", "a", "a", "a", []string{}, []int{}, []int{1, 2}}
	stepB := Step{"b", "b", "b", "b", []string{}, []int{0}, []int{}}
	stepC := Step{"c", "c", "c", "c", []string{}, []int{0}, []int{}}
	pipeline := Pipeline{stepA, stepB, stepC}

	require.True(t, pipeline.Connected())
	require.Len(t, pipeline.Components(), 1)
	require.True(t, pipeline.Acyclic())
	require.True(t, pipeline.Valid())
}

func TestMergingPipeline(t *testing.T) {
	stepA := Step{"a", "a", "a", "a", []string{}, []int{}, []int{2}}
	stepB := Step{"b", "b", "b", "b", []string{}, []int{}, []int{2}}
	stepC := Step{"c", "c", "c", "c", []string{}, []int{0, 1}, []int{}}
	pipeline := Pipeline{stepA, stepB, stepC}

	require.True(t, pipeline.Connected())
	require.Len(t, pipeline.Components(), 1)
	require.True(t, pipeline.Acyclic())
	require.True(t, pipeline.Valid())
}

func TestTriangularPipeline(t *testing.T) {
	stepA := Step{"a", "a", "a", "a", []string{}, []int{}, []int{1, 2}}
	stepB := Step{"b", "b", "b", "b", []string{}, []int{1}, []int{2}}
	stepC := Step{"c", "c", "c", "c", []string{}, []int{0, 1}, []int{}}
	pipeline := Pipeline{stepA, stepB, stepC}

	require.True(t, pipeline.Connected())
	require.Len(t, pipeline.Components(), 1)
	require.True(t, pipeline.Acyclic())
	require.True(t, pipeline.Valid())
}
