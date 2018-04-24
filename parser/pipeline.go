package parser

import "fmt"

type Steps map[string]Step
type Dependencies map[string][]string

type Pipeline interface {
    Execute()
}

type ConcretePipeline struct {
    Steps Steps
    Dependencies Dependencies

    completeSteps map[string]bool
}

func newPipeline(steps Steps, dependencies Dependencies) Pipeline {
    return &ConcretePipeline{Steps: steps, Dependencies: dependencies, completeSteps: make(map[string]bool)}
}

func (pipeline *ConcretePipeline) Execute() {
    ready := make(chan string, len(pipeline.Steps))
    complete := make(chan string, len(pipeline.Steps))

    fmt.Println("Execute()")

    pipeline.startReadySteps(ready)

    for {
        select {
        case name := <- ready:
            fmt.Printf("Execute step: %s\n", name)
            go pipeline.executeStep(name, complete)
        case name := <- complete:
            fmt.Printf("Step complete: %s\n", name)
            pipeline.completeSteps[name] = true
            pipeline.updateReadySteps(name, ready)
        }

        if pipeline.complete(ready, complete) {
            break
        }
    }

    fmt.Println("complete")
}

func (pipeline *ConcretePipeline) startReadySteps(ready chan string) {
    startedSteps := []string{}

    for name := range pipeline.Dependencies {
        if pipeline.ready(name) {
            startedSteps = append(startedSteps, name)
            ready <- name
        }
    }

    for _, name := range startedSteps {
        pipeline.removeStep(name)
    }
}

func (pipeline *ConcretePipeline) ready(stepName string) bool {
    return pipeline.Steps[stepName].Ready() && len(pipeline.Dependencies[stepName]) == 0
}

func (pipeline *ConcretePipeline) executeStep(stepName string, complete chan string) {
    pipeline.Steps[stepName].Execute()
    complete <- stepName
}

func (pipeline *ConcretePipeline) updateReadySteps(stepName string, ready chan string) {
    pipeline.removeDependency(stepName)
    pipeline.startReadySteps(ready)
}

func (pipeline *ConcretePipeline) complete(ready chan string, complete chan string) bool {
    for name := range pipeline.Steps {
        if _, ok := pipeline.completeSteps[name]; !ok {
            return false
        }
    }

    return true
}

func (pipeline *ConcretePipeline) removeStep(stepName string) {
    delete(pipeline.Dependencies, stepName)
}

func (pipeline *ConcretePipeline) removeDependency(stepName string) {
    for entryName, entryDependencies := range pipeline.Dependencies {
        updatedDependencies := []string{}
        for _, dependency := range entryDependencies {
            if stepName != dependency {
                updatedDependencies = append(updatedDependencies, dependency)
            }
        }

        pipeline.Dependencies[entryName] = updatedDependencies
    }
}

type Step interface {
    Ready() bool
    Execute()
}

type ConcreteStep struct {
	SourceRepo string
	SourcePath string
	Command    string

	Paths []string

	executed bool
}

func newStep(sourceRepo string, sourcePath string, command string, paths []string) Step {
    return &ConcreteStep{SourceRepo: sourceRepo, SourcePath: sourcePath, Command: command, Paths: paths}
}

func (step *ConcreteStep) Ready() bool {
    return !step.executed
}

func (step *ConcreteStep) Execute() {
    step.executed = true
}

func Main() {
    step := newStep("repo", "path", "command", []string{})
    steps := Steps{"first": step}
    dependencies := Dependencies{"first": []string{}}
    pipeline := newPipeline(steps, dependencies)

    fmt.Println("Main()")

    pipeline.Execute()
}
