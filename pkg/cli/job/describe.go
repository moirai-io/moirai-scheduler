package job

import (
	"context"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/moirai-io/moirai-scheduler/pkg/cli/options"
	"github.com/moirai-io/moirai-scheduler/pkg/cli/printer"
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
			moiraictl job describe my-queue
		`),
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Name = args[0]
			}

			c, err := client.New(config.GetConfigOrDie(), client.Options{})
			if err != nil {
				return err
			}

			job := &unstructured.Unstructured{}
			job.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   "batch",
				Kind:    "Job",
				Version: "v1",
			})

			err = c.Get(context.Background(), client.ObjectKey{
				Namespace: globalOpts.Namespace,
				Name:      opts.Name,
			}, job)

			if err != nil {
				return err
			}

			return printer.PrintObject(cmd.OutOrStdout(), job, f)
		},
	}

	f.AddFlags(cmd)

	return cmd
}
