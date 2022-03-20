package cli

import (
	"github.com/spf13/cobra"
)

// NewCmdQueueBinding returns the QueueBinding command for moiraictl.
func NewCmdQueueBinding() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "queuebinding <subcommand> [flags]",
		Short:     "queuebinding",
		Long:      `queuebinding`,
		ValidArgs: []string{"create", "get"},
	}

	return cmd
}
