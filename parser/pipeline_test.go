package parser

import (
    "testing"
    "github.com/stretchr/testify/require"
)

type MockedStep struct {
    step Step

    readyCount int
    executeCount int
}

func mockStep(step Step) *MockedStep {
    return &MockedStep{step: step}
}

func (mockedStep *MockedStep) Ready() bool {
    mockedStep.readyCount += 1
    return mockedStep.step.Ready()
}

func (mockedStep *MockedStep) Execute() {
    mockedStep.executeCount += 1
    mockedStep.step.Execute()
}

func TestExecutesOneStep(t *testing.T) {
    step := mockStep(newStep("repo", "path", "command", []string{}))
    steps := Steps{"first": step}
    dependencies := Dependencies{"first": []string{}}
    pipeline := newPipeline(steps, dependencies)

    pipeline.Execute()

    require.Equal(t, 1, step.executeCount)
}

func TestExecutesIndependentSteps(t *testing.T) {
    firstStep := mockStep(newStep("repo", "path", "command", []string{}))
    secondStep := mockStep(newStep("repo", "path", "command", []string{}))
    steps := Steps{"first": firstStep, "second": secondStep}
    dependencies := Dependencies{"first": []string{}, "second": []string{}}
    pipeline := newPipeline(steps, dependencies)

    pipeline.Execute()

    require.Equal(t, 1, firstStep.executeCount)
    require.Equal(t, 1, secondStep.executeCount)
}

func TestExecutesSequentialSteps(t *testing.T) {
    firstStep := mockStep(newStep("repo", "path", "command", []string{}))
    secondStep := mockStep(newStep("repo", "path", "command", []string{}))
    steps := Steps{"first": firstStep, "second": secondStep}
    dependencies := Dependencies{"first": []string{}, "second": []string{"first"}}
    pipeline := newPipeline(steps, dependencies)

    pipeline.Execute()

    require.Equal(t, 1, firstStep.executeCount)
    require.Equal(t, 1, secondStep.executeCount)
}

func TestExecutesDiamondSteps(t *testing.T) {
    firstStep := mockStep(newStep("repo", "path", "command", []string{}))
    secondStep := mockStep(newStep("repo", "path", "command", []string{}))
    thirdStep := mockStep(newStep("repo", "path", "command", []string{}))
    fourthStep := mockStep(newStep("repo", "path", "command", []string{}))
    steps := Steps{"first": firstStep, "second": secondStep, "third": thirdStep, "fourth": fourthStep}
    dependencies := Dependencies{
        "first": []string{},
        "second": []string{"first"},
        "third": []string{"first"},
        "fourth": []string{"second", "third"},
    }

    pipeline := newPipeline(steps, dependencies)
    pipeline.Execute()

    require.Equal(t, 1, firstStep.executeCount)
    require.Equal(t, 1, secondStep.executeCount)
    require.Equal(t, 1, thirdStep.executeCount)
    require.Equal(t, 1, fourthStep.executeCount)
}
