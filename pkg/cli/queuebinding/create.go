package queuebinding

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/moirai-io/moirai-scheduler/pkg/cli/options"
)

type createOptions struct {
}

func newCmdCreate(globalOpts *options.GlobalOptions) *cobra.Command {
	_ = &createOptions{}

	cmd := &cobra.Command{
		Use:   "create <name> [flags]",
		Short: "create",
		Long:  `create`,
		Example: heredoc.Doc(`
			moiraictl queuebinding create my-queue
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO:

			return nil
		},
	}

	return cmd
}
