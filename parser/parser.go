package parser

import (
    "os"
    "log"

    "github.com/spf13/afero"
    "gopkg.in/yaml.v2"
    "fmt"
)

var fs = afero.NewOsFs()

func parse(filename string) (BuildStepsDef, error) {
    var steps BuildStepsDef

    file, err := afero.ReadFile(fs, filename)
    if err != nil {
        return nil, err
    }

    err = yaml.UnmarshalStrict(file, &steps)
    if err != nil {
        return nil, err
    }

    err = steps.validate()
    if err != nil {
        return nil, err
    }

    return steps, nil
}

func help() string {
    return "Usage: cascade <steps.yaml>\n"
}

func Main() {
    if len(os.Args) < 2 {
        fmt.Printf(help())
        os.Exit(1)
    }

    filename := os.Args[1]

    steps, err := parse(filename)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Print(steps)
}
