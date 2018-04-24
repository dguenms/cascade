package parser

import (
	"bytes"
	"fmt"

	"github.com/toolkits/slice"
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
	steps := make(map[string]Step)
	dependencies := make(map[string][]string)

	for name, step := range def {
	    steps[name] = step.create()
	    dependencies[name] = step.Dependencies
    }

    return newPipeline(steps, dependencies)
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

func (def StepDef) create() Step {
    return newStep(def.SourceRepo, def.SourcePath, def.Command, def.Paths)
}

func (def StepDef) String() string {
	return fmt.Sprintf("{ source_repo: %s, source_path: %s, command: %s, dependencies: %s, paths: %s }",
		def.SourceRepo,
		def.SourcePath,
		def.Command,
		def.Dependencies,
		def.Paths)
}
