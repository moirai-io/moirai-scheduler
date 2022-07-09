package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rudeigerc/moirai/pkg/cli/options"
)

const (
	Version = "moiraictl 0.0.1"
)

// NewCmdVersion returns the Queue command for moiraictl.
func NewCmdVersion(globalOpts *options.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "version",
		Long:  `Version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), Version)

			return nil
		},
	}

	return cmd
}
