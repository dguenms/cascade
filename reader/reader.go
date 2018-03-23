package main

import (
    "os"
    "fmt"

    "github.com/kylelemons/go-gypsy/yaml"
)

type Steps struct {
    steps map[string]*Step
}

func (steps *Steps) addStep(step *Step) {
    if steps.steps == nil {
        steps.steps = make(map[string]*Step)
    }

    steps.steps[step.name] = step
}

func (steps *Steps) String() string {
    return fmt.Sprintf("%s", steps.steps)
}

type Step struct {
    name string

    sourceRepo string
    sourcePath string
    command string

    dependencies []string
    paths []string
}

func (step *Step) String() string {
    return fmt.Sprintf("{name: %s, source_repo: %s, source_path: %s, command: %s, dependencies: %s, paths: %s}", step.name, step.sourceRepo, step.sourcePath, step.command, step.dependencies, step.paths)
}

func parseSteps(file *yaml.File) *Steps {
    steps := new(Steps)
    root := file.Root

    switch node := root.(type) {
    case yaml.Map:
        for k, v := range node {
            steps.addStep(parseStep(k, v))
        }
    default:
        panic("yaml file needs to be a map of build steps")
    }

    return steps
}

func parseStep(name string, node yaml.Node) *Step {
    step := new(Step)

    var sourceRepo string
    var sourcePath string
    var command string

    var dependencies []string
    var paths []string

    switch n := node.(type) {
    case yaml.Map:
        for k, v := range n {
            switch k {
            case "source_repo":
                sourceRepo = parseScalar(v)
            case "source_path":
                sourcePath = parseScalar(v)
            case "command":
                command = parseScalar(v)
            case "dependencies":
                dependencies = parseList(v)
            case "paths":
                paths = parseList(v)
            }
        }
    }

    if sourceRepo == "" {
        panic("source_repo needs to be defined")
    }

    if command == "" {
        panic("command needs to be defined")
    }

    step.name = name
    step.sourceRepo = sourceRepo
    step.sourcePath = sourcePath
    step.command = command
    step.dependencies = dependencies
    step.paths = paths

    return step
}

func parseScalar(node yaml.Node) string {
    switch n := node.(type) {
    case yaml.Scalar:
        return n.String()
    default:
        panic("node expected to be a scalar")
    }
}

func parseList(node yaml.Node) []string {
    var list []string

    switch n := node.(type) {
    case yaml.List:
        for _, e := range n {
            list = append(list, parseScalar(e))
        }
    default:
        panic("node expected to be a list")
    }

    return list
}

func main() {
    if len(os.Args) < 2 {
        panic("Usage: cascade <steps.yaml>")
    }

    filename := os.Args[1]
    file, err := yaml.ReadFile(filename)
    if err != nil {
        panic(err)
    }

    steps := parseSteps(file)

    fmt.Print(steps)
}
