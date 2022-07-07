package queuebinding

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"

	moirai "github.com/moirai-io/moirai-scheduler/apis/scheduling/v1alpha1"
	"github.com/moirai-io/moirai-scheduler/pkg/cli/options"
	"github.com/moirai-io/moirai-scheduler/pkg/cli/printer"
	"github.com/moirai-io/moirai-scheduler/pkg/internal"
)

type describeOptions struct {
	Name string
}

func newCmdDescribe(globalOpts *options.GlobalOptions) *cobra.Command {
	f := genericclioptions.NewPrintFlags("")
	opts := describeOptions{}

	cmd := &cobra.Command{
		Use:     "describe [flags]",
		Short:   "describe",
		Long:    `describe`,
		Aliases: []string{"get"},
		Example: heredoc.Doc(`
			moiraictl queuebinding describe my-queue
		`),
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Name = args[0]
			}

			c, err := internal.NewClient()
			if err != nil {
				return err
			}

			queueBinding := &unstructured.Unstructured{}
			queueBinding.SetGroupVersionKind(moirai.GroupVersion.WithKind("QueueBinding"))

			err = c.Get(context.Background(), client.ObjectKey{
				Namespace: globalOpts.Namespace,
				Name:      opts.Name,
			}, queueBinding)
			if err != nil {
				return err
			}

			return printer.PrintObject(cmd.OutOrStdout(), queueBinding, f)
		},
	}

	f.AddFlags(cmd)

	return cmd
}
