package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreatesBuildStepFromDefinition(t *testing.T) {
    stepDef := StepDef{"repo", "path", "command", []string{}, []string{"path"}}

    step := stepDef.create()

	expectedStep := newStep("repo", "path", "command", []string{"path"})

	require.Equal(t, expectedStep, step)
}

func TestCreatesSingleBuildStepFromBuildStepsDefinition(t *testing.T) {
    stepDef := StepDef{"repo", "path", "command", []string{}, []string{"path"}}
    pipelineDef := PipelineDef{"first": stepDef}

	pipeline := pipelineDef.create()

	expectedSteps := Steps{"first": newStep("repo", "path", "command", []string{"path"})}
	expectedDependencies := Dependencies{"first": []string{}}
	expectedPipeline := newPipeline(expectedSteps, expectedDependencies)

	require.Equal(t, expectedPipeline, pipeline)
}

func TestCreatesMultipleBuildStepsFromBuildStepsDefinition(t *testing.T) {
    firstDef := StepDef{"first-repo", "first-path", "first command", []string{}, []string{"first-path"}}
    secondDef := StepDef{"second-repo", "second-path", "second command", []string{}, []string{"second-path"}}
    pipelineDef := PipelineDef{"first": firstDef, "second": secondDef}

    pipeline := pipelineDef.create()

    expectedSteps := Steps{
        "first": newStep("first-repo", "first-path", "first command", []string{"first-path"}),
        "second": newStep("second-repo", "second-path", "second command", []string{"second-path"}),
    }
    expectedDependencies := Dependencies{"first": []string{}, "second": []string{}}
    expectedPipeline := newPipeline(expectedSteps, expectedDependencies)

    require.Equal(t, expectedPipeline, pipeline)
}

func TestResolvesDependencies(t *testing.T) {
    firstDef := StepDef{"first-repo", "first-path", "first command", []string{}, []string{"first-path"}}
    secondDef := StepDef{"second-repo", "second-path", "second command", []string{"first"}, []string{"second-path"}}
    pipelineDef := PipelineDef{"first": firstDef, "second": secondDef}

    pipeline := pipelineDef.create()

    expectedSteps := Steps{
        "first": newStep("first-repo", "first-path", "first command", []string{"first-path"}),
        "second": newStep("second-repo", "second-path", "second command", []string{"second-path"}),
    }
    expectedDependencies := Dependencies{"first": []string{}, "second": []string{"first"}}
    expectedPipeline := newPipeline(expectedSteps, expectedDependencies)

    require.Equal(t, expectedPipeline, pipeline)
}

func TestResolvesMultipleStepsDependOnStep(t *testing.T) {
    firstDef := StepDef{"first-repo", "first-path", "first command", []string{}, []string{"first-path"}}
    secondDef := StepDef{"second-repo", "second-path", "second command", []string{"first"}, []string{"second-path"}}
    thirdDef := StepDef{"third-repo", "third-path", "third command", []string{"first"}, []string{"third-path"}}
    pipelineDef := PipelineDef{"first": firstDef, "second": secondDef, "third": thirdDef}

    pipeline := pipelineDef.create()

    expectedSteps := Steps{
        "first": newStep("first-repo", "first-path", "first command", []string{"first-path"}),
        "second": newStep("second-repo", "second-path", "second command", []string{"second-path"}),
        "third": newStep("third-repo", "third-path", "third command", []string{"third-path"}),
    }
    expectedDependencies := Dependencies{"first": []string{}, "second": []string{"first"}, "third": []string{"first"}}
    expectedPipeline := newPipeline(expectedSteps, expectedDependencies)

    require.Equal(t, expectedPipeline, pipeline)
}

func TestResolvesMultipleDependencies(t *testing.T) {
    firstDef := StepDef{"first-repo", "first-path", "first command", []string{}, []string{"first-path"}}
    secondDef := StepDef{"second-repo", "second-path", "second command", []string{}, []string{"second-path"}}
    thirdDef := StepDef{"third-repo", "third-path", "third command", []string{"first", "second"}, []string{"third-path"}}
    pipelineDef := PipelineDef{"first": firstDef, "second": secondDef, "third": thirdDef}

    pipeline := pipelineDef.create()

    expectedSteps := Steps{
        "first": newStep("first-repo", "first-path", "first command", []string{"first-path"}),
        "second": newStep("second-repo", "second-path", "second command", []string{"second-path"}),
        "third": newStep("third-repo", "third-path", "third command", []string{"third-path"}),
    }
    expectedDependencies := Dependencies{"first": []string{}, "second": []string{}, "third": []string{"first", "second"}}
    expectedPipeline := newPipeline(expectedSteps, expectedDependencies)

    require.Equal(t, expectedPipeline, pipeline)
}
