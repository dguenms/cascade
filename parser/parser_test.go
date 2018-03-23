package main

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestParsesSingleEntry(t *testing.T) {
    steps, err := parse("resources/single_entry.yaml")

    step := BuildStep{"some-repo", "some-path", "some command", []string(nil), []string(nil)}
    expectedSteps := BuildSteps{"build_step_name": step}

    require.NoError(t, err)
    require.Equal(t, expectedSteps, steps)
}

func TestParsesMultipleEntries(t *testing.T) {
    steps, err := parse("resources/multiple_entries.yaml")

    firstStep := BuildStep{"first-repo", "first-path", "first command", []string(nil), []string(nil)}
    secondStep := BuildStep{"second-repo", "second-path", "second command", []string(nil), []string(nil)}
    expectedSteps := BuildSteps{"first_step": firstStep, "second_step": secondStep}

    require.NoError(t, err)
    require.Equal(t, expectedSteps, steps)
}

func TestDoesNotParseNonMap(t *testing.T) {
    _, err := parse("resources/not_map.yaml")

    require.Error(t, err)
}

func TestParsesDepedenciesList(t *testing.T) {
    steps, err := parse("resources/dependencies_list.yaml")

    step := BuildStep{"some-repo", "some-path", "some command", []string{"first-dependency", "second-dependency"}, []string(nil)}
    expectedSteps := BuildSteps{"build_step": step}

    require.NoError(t, err)
    require.Equal(t, expectedSteps, steps)
}

func TestParsesPathsList(t *testing.T) {
    steps, err := parse("resources/paths_list.yaml")

    step := BuildStep{"some-repo", "some-path", "some command", []string(nil), []string{"first-path", "second-path"}}
    expectedSteps := BuildSteps{"build_step": step}

    require.NoError(t, err)
    require.Equal(t, expectedSteps, steps)
}

func TestFailsIfSourceRepoMissing(t *testing.T) {
    _, err := parse("resources/source_repo_missing.yaml")

    require.Error(t, err)
}

func TestParsesIfSourcePathMissing(t *testing.T) {
    steps, err := parse("resources/source_path_missing.yaml")

    step := BuildStep{"some-repo", "", "some command", []string(nil), []string(nil)}
    expectedSteps := BuildSteps{"build_step": step}

    require.NoError(t, err)
    require.Equal(t, expectedSteps, steps)
}

func TestFailsIfCommandMissing(t *testing.T) {
    _, err := parse("resources/command_missing.yaml")

    require.Error(t, err)
}

func TestParsesIfDependenciesMissing(t *testing.T) {
    steps, err := parse("resources/dependencies_missing.yaml")

    step := BuildStep{"some-repo", "some-path", "some command", []string(nil), []string(nil)}
    expectedSteps := BuildSteps{"build_step": step}

    require.NoError(t, err)
    require.Equal(t, expectedSteps, steps)
}

func TestParsesIfPathsMissing(t *testing.T) {
    steps, err := parse("resources/paths_missing.yaml")

    step := BuildStep{"some-repo", "some-path", "some command", []string(nil), []string(nil)}
    expectedSteps := BuildSteps{"build_step": step}

    require.NoError(t, err)
    require.Equal(t, expectedSteps, steps)
}

func TestFailsWithUndefinedAttribute(t *testing.T) {
    _, err := parse("resources/undefined_attribute.yaml")

    require.Error(t, err)
}
