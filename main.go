package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(c *cli.Context) error {
			getContainers()
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getContainers() {
	out, err := exec.Command("ls", "-la").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
