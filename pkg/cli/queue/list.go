package queue

import (
	"context"
	"fmt"

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

type listOptions struct {
}

func newCmdList(globalOpts *options.GlobalOptions) *cobra.Command {
	f := genericclioptions.NewPrintFlags("")
	_ = listOptions{}

	cmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "list",
		Long:  `list`,
		Example: heredoc.Doc(`
			moiraictl queue list
		`),
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := internal.NewClient()
			if err != nil {
				return err
			}

			queueList := &unstructured.UnstructuredList{}
			queueList.SetGroupVersionKind(moirai.GroupVersion.WithKind("QueueList"))

			err = c.List(context.Background(), queueList, &client.ListOptions{})
			if err != nil {
				return err
			}

			if len(queueList.Items) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No Queue found.")
				return nil
			}

			return printer.PrintObject(cmd.OutOrStdout(), queueList, f)
		},
	}

	f.AddFlags(cmd)

	return cmd
}
