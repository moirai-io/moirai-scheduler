package queue

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
		Use:   "describe [flags]",
		Short: "describe",
		Long:  `describe`,
		Example: heredoc.Doc(`
			moiraictl queue describe my-queue
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

			queue := &unstructured.Unstructured{}
			queue.SetGroupVersionKind(moirai.GroupVersion.WithKind("Queue"))

			err = c.Get(context.Background(), client.ObjectKey{
				Namespace: globalOpts.Namespace,
				Name:      opts.Name,
			}, queue)
			if err != nil {
				return err
			}

			return printer.PrintObject(cmd.OutOrStdout(), queue, f)
		},
	}

	f.AddFlags(cmd)

	return cmd
}
