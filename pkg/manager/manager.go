package manager

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	schedulingv1alpha1 "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
)

type Manager interface {
	GetQueueBinding(pod *corev1.Pod) (*schedulingv1alpha1.QueueBinding, error)
	AnnotatePod(ctx context.Context, pod *corev1.Pod)
}

type MoiraiManager struct {
	client      kubernetes.Interface
	moiraiCache cache.Cache
}

func NewMoiraiManager(client kubernetes.Interface, moiraiCache cache.Cache) *MoiraiManager {
	return &MoiraiManager{
		client:      client,
		moiraiCache: moiraiCache,
	}
}

// GetQueueBinding returns the queue binding of the specified pod
func (m *MoiraiManager) GetQueueBinding(ctx context.Context, pod *corev1.Pod) (*schedulingv1alpha1.QueueBinding, error) {
	queueBindingLabel := pod.Labels[schedulingv1alpha1.QueueBindingLabel]
	if len(queueBindingLabel) == 0 {
		return nil, fmt.Errorf("unable to fetch the QueueBinding label of the pod")
	}

	var queueBinding schedulingv1alpha1.QueueBinding
	err := m.moiraiCache.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: queueBindingLabel}, &queueBinding)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch QueueBinding %s", queueBindingLabel)
	}

	return &queueBinding, nil
}

// AnnotatePod sets the annotation to the specified pod
func (m *MoiraiManager) AnnotatePod(ctx context.Context, pod *corev1.Pod) {
	annotations := map[string]string{}
	if pod.Annotations != nil {
		annotations = pod.Annotations
	}
	// TODO:
	annotations["test"] = "test"
	pod.Annotations = annotations
	m.client.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
}

// GetNodeAvaliableResource returns the available resource of the node
func (m *MoiraiManager) GetNodeAvaliableResource() *framework.Resource {
	return nil
}
