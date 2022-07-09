package manager

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	moirai "github.com/rudeigerc/moirai/apis/scheduling/v1alpha1"
)

// Manager is an interface for managing scheduler plugin.
type Manager interface {
	// Pod
	AnnotatePod(ctx context.Context, pod *corev1.Pod)
	// QueueBinding
	GetQueueBinding(pod *corev1.Pod) (*moirai.QueueBinding, error)
	GetQueueBindingListFromQueue(ctx context.Context, queue string) (*moirai.QueueBindingList, error)
}

// MoiraiManager is a concrete implementation of Manager interface for Moirai.
type MoiraiManager struct {
	KubeClient  kubernetes.Interface
	MoiraiCache cache.Cache
}

// NewMoiraiManager returns a new MoiraiManager.
func NewMoiraiManager(client kubernetes.Interface, moiraiCache cache.Cache) *MoiraiManager {
	return &MoiraiManager{
		KubeClient:  client,
		MoiraiCache: moiraiCache,
	}
}
