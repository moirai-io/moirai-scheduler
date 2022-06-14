package manager

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	moirai "github.com/moirai-io/moirai-scheduler/apis/scheduling/v1alpha1"
)

// GetQueueBinding returns the queue binding of the specified pod
func (m *MoiraiManager) GetQueueBinding(ctx context.Context, pod *corev1.Pod) (*moirai.QueueBinding, error) {
	queueBindingLabel := pod.Labels[moirai.QueueBindingLabel]
	if len(queueBindingLabel) == 0 {
		return nil, fmt.Errorf("unable to fetch the QueueBinding label of the pod")
	}

	queueBinding, err := m.GetQueueBindingByName(ctx, pod.Namespace, queueBindingLabel)
	if err != nil {
		return nil, err
	}

	return queueBinding, nil
}

// GetQueueBindingByName returns the queue binding of the specified name
func (m *MoiraiManager) GetQueueBindingByName(ctx context.Context, namespace string, name string) (*moirai.QueueBinding, error) {
	var queueBinding moirai.QueueBinding
	err := m.MoiraiCache.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &queueBinding)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch QueueBinding %s", name)
	}

	return &queueBinding, nil
}
