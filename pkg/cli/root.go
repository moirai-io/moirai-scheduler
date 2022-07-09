package cli

import (
	"github.com/spf13/cobra"

	"github.com/rudeigerc/moirai/pkg/cli/job"
	"github.com/rudeigerc/moirai/pkg/cli/options"
	"github.com/rudeigerc/moirai/pkg/cli/queue"
	"github.com/rudeigerc/moirai/pkg/cli/queuebinding"
	"github.com/rudeigerc/moirai/pkg/cli/version"
)

// NewCmdRoot returns the root command for moiraictl.
func NewCmdRoot() *cobra.Command {
	opts := &options.GlobalOptions{}

	cmd := &cobra.Command{
		Use:       "moiraictl <command> <subcommand> [flags]",
		Short:     "Moirai CLI",
		Long:      `Moirai command line tool.`,
		ValidArgs: []string{"queue", "queuebinding", "version"},
	}

	cmd.PersistentFlags().StringVarP(&opts.Namespace, "namespace", "n", "default", "If present, the namespace scope for this CLI request")
	cmd.PersistentFlags().StringVarP(&opts.Kubeconfig, "kubeconfig", "", "", "Path to the kubeconfig file to use for CLI requests")

	cmd.AddCommand(queue.NewCmdQueue(opts))
	cmd.AddCommand(queuebinding.NewCmdQueueBinding(opts))
	cmd.AddCommand(version.NewCmdVersion(opts))
	cmd.AddCommand(job.NewCmdJob(opts))

	return cmd
}
