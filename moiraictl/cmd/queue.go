package cmd

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	"github.com/moirai-io/moirai-operator/moiraictl/cmd/utils"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	schedulingv1alpha1 "github.com/moirai-io/moirai-operator/api/v1alpha1"
)

// NewCmdQueue returns the Queue command for moiraictl.
func NewCmdQueue() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "queue <subcommand> [flags]",
		Short:     "Queue",
		Long:      `Queue`,
		ValidArgs: []string{"create", "get"},
	}

	cmd.AddCommand(newCmdCreate())
	cmd.AddCommand(newCmdGet())
	return cmd
}

type createOptions struct {
	Name   string
	CPU    string
	Memory string
}

func newCmdCreate() *cobra.Command {
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

			c, err := utils.NewClient()
			if err != nil {
				return err
			}

			queue := &schedulingv1alpha1.Queue{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: Namespace,
					Name:      opts.Name,
				},
				Spec: schedulingv1alpha1.QueueSpec{
					Resources: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse(opts.CPU),
						v1.ResourceMemory: resource.MustParse(opts.Memory),
					},
				},
			}

			err = c.Create(context.Background(), queue)

			return err
		},
	}

	cmd.Flags().StringVar(&opts.CPU, "cpu", "", "")
	cmd.Flags().StringVar(&opts.Memory, "memory", "", "")

	return cmd
}

type getOptions struct {
	Name string
}

func newCmdGet() *cobra.Command {
	opts := getOptions{}

	cmd := &cobra.Command{
		Use:   "get [flags]",
		Short: "get",
		Long:  `get`,
		Example: heredoc.Doc(`
			moiraictl queue create my-queue
		`),
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.Name = args[0]
			}

			c, err := utils.NewClient()
			if err != nil {
				return err
			}

			queue := &schedulingv1alpha1.Queue{}
			err = c.Get(context.Background(), client.ObjectKey{
				Namespace: Namespace,
				Name:      opts.Name,
			}, queue)

			return err
		},
	}

	return cmd
}
