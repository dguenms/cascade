package main

import (
    "testing"
    "github.com/stretchr/testify/require"
)

var (
    stepDef  = BuildStepDef{"some-repo", "some-path", "some command", []string(nil), []string{"some-path"}}
    stepsDef = BuildStepsDef{"build_step": stepDef}

    expectedBuildStep = BuildStep{"some-repo", "some-path", "some command", []string(nil), []string{"some-path"}, []string(nil)}
    expectedBuildSteps = BuildSteps{"build_step": expectedBuildStep}
)

func TestCreatesBuildStepFromDefinition(t *testing.T) {
    buildStep := stepDef.create([]string(nil))

    require.Equal(t, expectedBuildStep, buildStep)
}

func TestCreatesBuildStepFromDefinitionWithNextSteps(t *testing.T) {
    buildStep := stepDef.create([]string{"first", "second"})

    expectedBuildStep := expectedBuildStep
    expectedBuildStep.NextSteps = []string{"first", "second"}

    require.Equal(t, expectedBuildStep, buildStep)
}

func TestCreatesSingleBuildStepFromBuildStepsDefinition(t *testing.T) {
    buildSteps := stepsDef.create()

    require.Equal(t, expectedBuildSteps, buildSteps)
}

func TestCreatesMultipleBuildStepsFromBuildStepsDefinition(t *testing.T) {
    stepDef2 := BuildStepDef{"second-repo", "second-path", "second command", []string(nil), []string{"second-path"}}
    stepsDef := BuildStepsDef{"first_step": stepDef, "second_step": stepDef2}

    buildSteps := stepsDef.create()

    expectedBuildStep2 := BuildStep{"second-repo", "second-path", "second command", []string(nil), []string{"second-path"}, []string(nil)}
    expectedBuildSteps := BuildSteps{"first_step": expectedBuildStep, "second_step": expectedBuildStep2}

    require.Equal(t, expectedBuildSteps, buildSteps)
}
