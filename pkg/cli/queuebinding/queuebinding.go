package queuebinding

import (
	"github.com/spf13/cobra"

	"github.com/rudeigerc/moirai/pkg/cli/options"
)

// NewCmdQueueBinding returns the QueueBinding command for moiraictl.
func NewCmdQueueBinding(globalOpts *options.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "queuebinding <subcommand> [flags]",
		Short:     "queuebinding",
		Long:      `queuebinding`,
		ValidArgs: []string{"create", "describe", "list"},
	}

	cmd.AddCommand(newCmdCreate(globalOpts))
	cmd.AddCommand(newCmdDescribe(globalOpts))
	cmd.AddCommand(newCmdList(globalOpts))
	return cmd
}
