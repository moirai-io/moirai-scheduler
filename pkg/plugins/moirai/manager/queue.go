package manager

import (
	"context"

	moirai "github.com/rudeigerc/moirai/apis/scheduling/v1alpha1"
	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetQueueBindingListFromQueue returns the list of queue bindings of the specified queue
func (m *MoiraiManager) GetQueueBindingListFromQueue(ctx context.Context, queue string) (*moirai.QueueBindingList, error) {
	var queueBindingList moirai.QueueBindingList
	opts := &client.ListOptions{
		FieldSelector: fields.SelectorFromSet(
			fields.Set{
				"spec.queue": queue,
			},
		),
	}

	err := m.MoiraiCache.List(ctx, &queueBindingList, opts)
	if err != nil {
		return nil, err
	}
	return &queueBindingList, nil
}
