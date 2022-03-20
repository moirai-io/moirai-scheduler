package cli

import (
	"github.com/spf13/cobra"
)

var (
	Namespace  string
	Kubeconfig string
)

// NewCmdRoot returns the root command for moiraictl.
func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "moiraictl <command> <subcommand> [flags]",
		Short:     "Moirai CLI",
		Long:      `Moirai command line tool.`,
		ValidArgs: []string{"queue", "queuebinding", "version"},
	}

	cmd.PersistentFlags().StringVarP(&Namespace, "namespace", "n", "default", "If present, the namespace scope for this CLI request")
	cmd.PersistentFlags().StringVarP(&Kubeconfig, "kubeconfig", "", "", "Path to the kubeconfig file to use for CLI requests")

	cmd.AddCommand(NewCmdQueue())
	cmd.AddCommand(NewCmdQueueBinding())
	cmd.AddCommand(NewCmdVersion())

	return cmd
}
