package queue

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	moirai "github.com/moirai-io/moirai-scheduler/apis/scheduling/v1alpha1"
	"github.com/moirai-io/moirai-scheduler/pkg/cli/options"
	"github.com/moirai-io/moirai-scheduler/pkg/internal"
)

type createOptions struct {
	Name   string
	CPU    string
	Memory string
}

func newCmdCreate(globalOpts *options.GlobalOptions) *cobra.Command {
	opts := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create <name> [flags]",
		Short: "create",
		Long:  `create`,
		Example: heredoc.Doc(`
			moiraictl queue create my-queue
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Name = args[0]
			}

			c, err := internal.NewClient()
			if err != nil {
				return err
			}

			queue := &moirai.Queue{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: globalOpts.Namespace,
					Name:      opts.Name,
				},
				Spec: moirai.QueueSpec{
					Capacity: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse(opts.CPU),
						v1.ResourceMemory: resource.MustParse(opts.Memory),
					},
				},
			}

			err = c.Create(context.Background(), queue)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Queue '%s' created.\n", opts.Name)

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.CPU, "cpu", "", "")
	cmd.Flags().StringVar(&opts.Memory, "memory", "", "")

	return cmd
}
