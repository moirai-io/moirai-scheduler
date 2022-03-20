package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "moiraictl version 0.0.1"
)

// NewCmdQueue returns the Queue command for moiraictl.
func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "version",
		Long:  `Version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	return cmd
}
