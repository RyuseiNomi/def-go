package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

type workFlowManager struct {
	reflesh chan bool
	exit    chan bool
}

const (
	containerIDIndex = 0
	namesIndex       = 1
	statusIndex      = 2
)

func newWorkflowManager() *workFlowManager {
	return &workFlowManager{
		reflesh: make(chan bool, 1),
		exit:    make(chan bool, 1),
	}
}

func main() {
	app := &cli.App{
		Name:  "def",
		Usage: "manipulate docker containers",
		Action: func(c *cli.Context) error {
			wfm := newWorkflowManager()
			wfm.Handle()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
