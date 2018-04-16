package parser

import (
	"bytes"
	"fmt"

	"github.com/toolkits/slice"
	"sort"
)

type PipelineDef map[string]StepDef

func (def PipelineDef) validate() error {
	var buffer bytes.Buffer

	for name, step := range def {
		err := step.validate()
		if err != nil {
			buffer.WriteString(fmt.Sprintf("%s: %s\n", name, err))
		}

		if slice.ContainsString(step.Dependencies, name) {
			buffer.WriteString(fmt.Sprintf("%s: has dependency on itself", name))
		}

		for _, dependency := range step.Dependencies {
			if _, found := def[dependency]; !found {
				buffer.WriteString(fmt.Sprintf("%s: declared dependency %s is not defined\n", name, dependency))
			}
		}
	}

	message := buffer.String()
	if message != "" {
		return fmt.Errorf("Errors in build def:\n%s", message)
	}

	return nil
}

func (def PipelineDef) create() Pipeline {
	pipeline := Pipeline{}

	var sortedEntries []string
	for name, _ := range def {
		sortedEntries = append(sortedEntries, name)
	}
	sort.Strings(sortedEntries)

	indices := make(map[string]int)
	for index, name := range sortedEntries {
		indices[name] = index
	}

	for _, name := range sortedEntries {
		step := def[name]
		predecessors := def.predecessors(name, indices)
		successors := def.successors(name, indices)
		pipeline = append(pipeline, step.create(name, predecessors, successors))
	}

	return pipeline
}

func (def PipelineDef) predecessors(stepName string, indices map[string]int) []int {
	step := def[stepName]
	predecessors := []int{}

	for _, dependency := range step.Dependencies {
		predecessors = append(predecessors, indices[dependency])
	}

	sort.Ints(predecessors)

	return predecessors
}

func (def PipelineDef) successors(stepName string, indices map[string]int) []int {
	successors := []int{}

	for name, step := range def {
		if slice.ContainsString(step.Dependencies, stepName) {
			index := indices[name]
			if !slice.ContainsInt(successors, index) {
				successors = append(successors, index)
			}
		}
	}

	sort.Ints(successors)

	return successors
}

func (def PipelineDef) String() string {
	var buffer bytes.Buffer

	for k, v := range def {
		buffer.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	return buffer.String()
}

type StepDef struct {
	SourceRepo string `yaml:"source_repo"`
	SourcePath string `yaml:"source_path"`
	Command    string

	Dependencies []string
	Paths        []string
}

func (def StepDef) validate() error {
	var missingEntries []string

	if def.SourceRepo == "" {
		missingEntries = append(missingEntries, "source_repo")
	}

	if def.Command == "" {
		missingEntries = append(missingEntries, "command")
	}

	if len(missingEntries) > 0 {
		return fmt.Errorf("Build def has missing entries: %s", missingEntries)
	}

	return nil
}

func (def StepDef) create(name string, predecessors []int, successors []int) Step {
	return Step{name, def.SourceRepo, def.SourcePath, def.Command, def.Paths, predecessors, successors}
}

func (def StepDef) String() string {
	return fmt.Sprintf("{ source_repo: %s, source_path: %s, command: %s, dependencies: %s, paths: %s }",
		def.SourceRepo,
		def.SourcePath,
		def.Command,
		def.Dependencies,
		def.Paths)
}
