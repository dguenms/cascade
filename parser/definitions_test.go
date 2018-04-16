package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	stepDef     = StepDef{"some-repo", "some-path", "some command", []string{}, []string{"some-path"}}
	pipelineDef = PipelineDef{"first_step": stepDef}

	expectedStep     = Step{"first_step", "some-repo", "some-path", "some command", []string{"some-path"}, []int{}, []int{}}
	expectedPipeline = Pipeline{expectedStep}
)

func TestCreatesBuildStepFromDefinition(t *testing.T) {
	step := stepDef.create("first_step", []int{}, []int{})

	require.Equal(t, expectedStep, step)
}

func TestCreatesBuildStepFromDefinitionWithSuccessors(t *testing.T) {
	step := stepDef.create("first_step", []int{}, []int{1, 2})

	expectedStep := expectedStep
	expectedStep.Successors = []int{1, 2}

	require.Equal(t, expectedStep, step)
}

func TestCreatesSingleBuildStepFromBuildStepsDefinition(t *testing.T) {
	pipeline := pipelineDef.create()

	require.Equal(t, expectedPipeline, pipeline)
}

func TestCreatesMultipleBuildStepsFromBuildStepsDefinition(t *testing.T) {
	stepDef2 := StepDef{"second-repo", "second-path", "second command", []string{}, []string{"second-path"}}
	pipelineDef := PipelineDef{"first_step": stepDef, "second_step": stepDef2}

	pipeline := pipelineDef.create()

	expectedStep2 := Step{"second_step", "second-repo", "second-path", "second command", []string{"second-path"}, []int{}, []int{}}

	require.Len(t, pipeline, 2)
	require.Contains(t, pipeline, expectedStep)
	require.Contains(t, pipeline, expectedStep2)
}

func TestResolvesSuccessors(t *testing.T) {
	stepDef2 := StepDef{"second-repo", "second-path", "second command", []string{"first_step"}, []string{"second-path"}}
	pipelineDef := PipelineDef{"first_step": stepDef, "second_step": stepDef2}

	pipeline := pipelineDef.create()

	expectedStep := Step{"first_step", "some-repo", "some-path", "some command", []string{"some-path"}, []int{}, []int{1}}
	expectedStep2 := Step{"second_step", "second-repo", "second-path", "second command", []string{"second-path"}, []int{0}, []int{}}

	require.Len(t, pipeline, 2)
	require.Contains(t, pipeline, expectedStep)
	require.Contains(t, pipeline, expectedStep2)
}

func TestResolvesMultipleSuccessors(t *testing.T) {
	stepDef2 := StepDef{"second-repo", "second-path", "second command", []string{"first_step"}, []string{}}
	stepDef3 := StepDef{"third-repo", "third-path", "third command", []string{"first_step"}, []string{}}
	pipelineDef := PipelineDef{"first_step": stepDef, "second_step": stepDef2, "third_step": stepDef3}

	pipeline := pipelineDef.create()

	expectedStep := Step{"first_step", "some-repo", "some-path", "some command", []string{"some-path"}, []int{}, []int{1, 2}}
	expectedStep2 := Step{"second_step", "second-repo", "second-path", "second command", []string{}, []int{0}, []int{}}
	expectedStep3 := Step{"third_step", "third-repo", "third-path", "third command", []string{}, []int{0}, []int{}}

	require.Len(t, pipeline, 3)
	require.Contains(t, pipeline, expectedStep)
	require.Contains(t, pipeline, expectedStep2)
	require.Contains(t, pipeline, expectedStep3)
}
