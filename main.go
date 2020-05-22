package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

var ui cli.Ui

func main() {
	os.Exit(realMain())
}

func realMain() int {
	ui = &cli.BasicUi{Writer: os.Stdout}
	cli := &cli.CLI{
		Args: os.Args[1:],
		Commands: map[string]cli.CommandFactory{
			"create": func() (cli.Command, error) {
				return &GenCommand{}, nil
			},
		},
		HelpFunc: cli.BasicHelpFunc("ssl-create"),
		Version:  "1.0.0",
	}

	exitCode, err := cli.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
