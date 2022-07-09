package queue

import (
	"github.com/spf13/cobra"

	"github.com/rudeigerc/moirai/pkg/cli/options"
)

// NewCmdQueue returns the Queue command for moiraictl.
func NewCmdQueue(globalOpts *options.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "queue <subcommand> [flags]",
		Short:     "Queue",
		Long:      `Queue`,
		ValidArgs: []string{"create", "describe", "list"},
	}

	cmd.AddCommand(newCmdCreate(globalOpts))
	cmd.AddCommand(newCmdDescribe(globalOpts))
	cmd.AddCommand(newCmdList(globalOpts))
	return cmd
}
