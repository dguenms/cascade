package main

import (
    "bytes"
    "fmt"
)

type BuildSteps map[string]BuildStep

func (steps BuildSteps) validate() error {
    var buffer bytes.Buffer

    for k, v := range steps {
        err := v.validate()
        if err != nil {
            buffer.WriteString(fmt.Sprintf("%s: %s\n", k, err))
        }
    }

    message := buffer.String()
    if message != "" {
        return fmt.Errorf("Errors in build steps:\n%s", message)
    }

    return nil
}

func (steps BuildSteps) String() string {
    var buffer bytes.Buffer

    for k, v := range steps {
        buffer.WriteString(fmt.Sprintf("%s: %s\n", k, v))
    }

    return buffer.String()
}

type BuildStep struct {
    SourceRepo string `yaml:"source_repo"`
    SourcePath string `yaml:"source_path"`
    Command string

    Dependencies []string
    Paths []string
}

func (step BuildStep) validate() error {
    var missingEntries []string

    if step.SourceRepo == "" {
        missingEntries = append(missingEntries, "source_repo")
    }

    if step.Command == "" {
        missingEntries = append(missingEntries, "command")
    }

    if missingEntries != nil {
        return fmt.Errorf("Build step has missing entries: %s", missingEntries)
    }

    return nil
}

func (step BuildStep) String() string {
    return fmt.Sprintf("{ source_repo: %s, source_path: %s, command: %s, dependencies: %s, paths: %s }",
        step.SourceRepo,
        step.SourcePath,
        step.Command,
        step.Dependencies,
        step.Paths)
}
