package main

import (
    "bytes"
    "fmt"

    "github.com/toolkits/slice"
)

type BuildStepsDef map[string]BuildStepDef

func (def BuildStepsDef) validate() error {
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

func (def BuildStepsDef) create() BuildSteps {
    var steps BuildSteps = make(BuildSteps)

    for name, step := range def {
        steps[name] = step.create(def.determineNextSteps(name))
    }

    return steps
}

func (def BuildStepsDef) determineNextSteps(stepName string) []string {
    var nextSteps []string

    for name, step := range def {
        if slice.ContainsString(step.Dependencies, stepName) {
            nextSteps = append(nextSteps, name)
        }
    }

    return nextSteps
}

func (def BuildStepsDef) String() string {
    var buffer bytes.Buffer

    for k, v := range def {
        buffer.WriteString(fmt.Sprintf("%s: %s\n", k, v))
    }

    return buffer.String()
}

type BuildStepDef struct {
    SourceRepo string `yaml:"source_repo"`
    SourcePath string `yaml:"source_path"`
    Command string

    Dependencies []string
    Paths []string
}

func (def BuildStepDef) validate() error {
    var missingEntries []string

    if def.SourceRepo == "" {
        missingEntries = append(missingEntries, "source_repo")
    }

    if def.Command == "" {
        missingEntries = append(missingEntries, "command")
    }

    if missingEntries != nil {
        return fmt.Errorf("Build def has missing entries: %s", missingEntries)
    }

    return nil
}

func (def BuildStepDef) create(nextSteps []string) BuildStep {
    return BuildStep{def.SourceRepo, def.SourcePath, def.Command, def.Dependencies, def.Paths, nextSteps}
}

func (def BuildStepDef) String() string {
    return fmt.Sprintf("{ source_repo: %s, source_path: %s, command: %s, dependencies: %s, paths: %s }",
        def.SourceRepo,
        def.SourcePath,
        def.Command,
        def.Dependencies,
        def.Paths)
}

type BuildSteps map[string]BuildStep

type BuildStep struct {
    SourceRepo string
    SourcePath string
    Command string

    Dependencies []string
    Paths []string

    NextSteps []string
}
