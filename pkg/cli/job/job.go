package job

import (
	"github.com/spf13/cobra"

	"github.com/rudeigerc/moirai/pkg/cli/options"
)

// NewCmdJob returns the Job command for moiraictl.
func NewCmdJob(globalOpts *options.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "job <subcommand> [flags]",
		Short:     "job",
		Long:      `job`,
		Aliases:   []string{"j"},
		ValidArgs: []string{"describe"},
	}

	cmd.AddCommand(newCmdDescribe(globalOpts))
	return cmd
}
