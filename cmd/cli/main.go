package main

import (
	"os"

	moiraictl "github.com/moirai-io/moirai-scheduler/pkg/cli"
)

func main() {
	cli := moiraictl.NewCmdRoot()
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
