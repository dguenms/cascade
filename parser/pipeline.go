package parser

import (
	"fmt"

	"github.com/yourbasic/graph"
)

type Pipeline []Step

func (pipeline Pipeline) Order() int {
	return len(pipeline)
}

func (pipeline Pipeline) Visit(v int, do func(w int, c int64) (skip bool)) (aborted bool) {
	node := pipeline[v]

	for _, successor := range node.Successors {
		skip := do(successor, 0)

		if skip {
			return true
		}
	}

	return false
}

func (pipeline Pipeline) Connected() bool {
	return graph.Connected(pipeline)
}

func (pipeline Pipeline) Components() [][]int {
	return graph.Components(pipeline)
}

func (pipeline Pipeline) Acyclic() bool {
	return graph.Acyclic(pipeline)
}

func (pipeline Pipeline) Valid() bool {
	if !pipeline.Connected() {
		return false
	}

	if len(pipeline.Components()) > 1 {
		return false
	}

	if !pipeline.Acyclic() {
		return false
	}

	return true
}

type Step struct {
	Name string

	SourceRepo string
	SourcePath string
	Command    string

	Paths []string

	Predecessors []int
	Successors   []int
}

func (step Step) String() string {
	return fmt.Sprintf("{ name: '%s', source_repo: '%s', source_path: '%s', command: '%s', paths: %v, predecessors: %v, successors: %v }",
		step.Name,
		step.SourceRepo,
		step.SourcePath,
		step.Command,
		step.Paths,
		step.Predecessors,
		step.Successors)
}
