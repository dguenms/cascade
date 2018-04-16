package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsesSingleEntry(t *testing.T) {
	steps, err := parse("resources/single_entry.yaml")

	step := StepDef{"some-repo", "some-path", "some command", []string(nil), []string(nil)}
	expectedSteps := PipelineDef{"build_step_name": step}

	require.NoError(t, err)
	require.Equal(t, expectedSteps, steps)
}

func TestParsesMultipleEntries(t *testing.T) {
	steps, err := parse("resources/multiple_entries.yaml")

	firstStep := StepDef{"first-repo", "first-path", "first command", []string(nil), []string(nil)}
	secondStep := StepDef{"second-repo", "second-path", "second command", []string(nil), []string(nil)}
	expectedSteps := PipelineDef{"first_step": firstStep, "second_step": secondStep}

	require.NoError(t, err)
	require.Equal(t, expectedSteps, steps)
}

func TestDoesNotParseNonMap(t *testing.T) {
	_, err := parse("resources/not_map.yaml")

	require.Error(t, err)
}

func TestParsesDependenciesList(t *testing.T) {
	steps, err := parse("resources/dependencies_list.yaml")

	step := StepDef{"some-repo", "some-path", "some command", []string{"first_dependency", "second_dependency"}, []string(nil)}
	step2 := StepDef{"first", "", "first", []string(nil), []string(nil)}
	step3 := StepDef{"second", "", "second", []string(nil), []string(nil)}
	expectedSteps := PipelineDef{"build_step": step, "first_dependency": step2, "second_dependency": step3}

	require.NoError(t, err)
	require.Equal(t, expectedSteps, steps)
}

func TestFailsIfUnresolvableDependencies(t *testing.T) {
	_, err := parse("resources/unresolvable_dependencies.yaml")

	require.Error(t, err)
}

func TestFailsIfSelfReferentialDependency(t *testing.T) {
	_, err := parse("resources/self_referential_dependency.yaml")

	require.Error(t, err)
}

func TestParsesPathsList(t *testing.T) {
	steps, err := parse("resources/paths_list.yaml")

	step := StepDef{"some-repo", "some-path", "some command", []string(nil), []string{"first-path", "second-path"}}
	expectedSteps := PipelineDef{"build_step": step}

	require.NoError(t, err)
	require.Equal(t, expectedSteps, steps)
}

func TestFailsIfSourceRepoMissing(t *testing.T) {
	_, err := parse("resources/source_repo_missing.yaml")

	require.Error(t, err)
}

func TestParsesIfSourcePathMissing(t *testing.T) {
	steps, err := parse("resources/source_path_missing.yaml")

	step := StepDef{"some-repo", "", "some command", []string(nil), []string(nil)}
	expectedSteps := PipelineDef{"build_step": step}

	require.NoError(t, err)
	require.Equal(t, expectedSteps, steps)
}

func TestFailsIfCommandMissing(t *testing.T) {
	_, err := parse("resources/command_missing.yaml")

	require.Error(t, err)
}

func TestParsesIfDependenciesMissing(t *testing.T) {
	steps, err := parse("resources/dependencies_missing.yaml")

	step := StepDef{"some-repo", "some-path", "some command", []string(nil), []string(nil)}
	expectedSteps := PipelineDef{"build_step": step}

	require.NoError(t, err)
	require.Equal(t, expectedSteps, steps)
}

func TestParsesIfPathsMissing(t *testing.T) {
	steps, err := parse("resources/paths_missing.yaml")

	step := StepDef{"some-repo", "some-path", "some command", []string(nil), []string(nil)}
	expectedSteps := PipelineDef{"build_step": step}

	require.NoError(t, err)
	require.Equal(t, expectedSteps, steps)
}

func TestFailsWithUndefinedAttribute(t *testing.T) {
	_, err := parse("resources/undefined_attribute.yaml")

	require.Error(t, err)
}
