package queuebinding

import (
	"github.com/spf13/cobra"

	"github.com/moirai-io/moirai-scheduler/pkg/cli/options"
)

// NewCmdQueueBinding returns the QueueBinding command for moiraictl.
func NewCmdQueueBinding(globalOpts *options.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "queuebinding <subcommand> [flags]",
		Short:     "queuebinding",
		Long:      `queuebinding`,
		ValidArgs: []string{"create", "get"},
	}

	return cmd
}
