package main

import (
	"os"

	"github.com/moirai-io/moirai-scheduler/moiraictl/cmd"
)

func main() {
	cli := cmd.NewCmdRoot()
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}