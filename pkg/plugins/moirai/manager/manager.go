package manager

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	moirai "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
)

type Manager interface {
	GetQueueBinding(pod *corev1.Pod) (*moirai.QueueBinding, error)
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
