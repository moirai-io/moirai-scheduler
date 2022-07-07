package queuebinding

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
			moiraictl queuebinding list
		`),
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := internal.NewClient()
			if err != nil {
				return err
			}

			queueBindingList := &unstructured.UnstructuredList{}
			queueBindingList.SetGroupVersionKind(moirai.GroupVersion.WithKind("QueueBindingList"))

			err = c.List(context.Background(), queueBindingList, &client.ListOptions{})
			if err != nil {
				return err
			}

			if len(queueBindingList.Items) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No QueueBinding found.")
				return nil
			}

			return printer.PrintObject(cmd.OutOrStdout(), queueBindingList, f)
		},
	}

	f.AddFlags(cmd)

	return cmd
}
