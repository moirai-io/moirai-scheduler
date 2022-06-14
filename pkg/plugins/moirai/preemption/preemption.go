package preemption

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

type Interface interface {
}

// Preemptor is an implementation of Interface.
type Preemptor struct {
	PluginName  string
	Handler     framework.Handle
	PodLister   corelisters.PodLister
	MoiraiCache cache.Cache
}

// NewPreemptor
func NewPreemptor(pluginName string, handler framework.Handle, podLister corelisters.PodLister, moiraiCache cache.Cache) *Preemptor {
	return &Preemptor{
		PluginName:  pluginName,
		Handler:     handler,
		PodLister:   podLister,
		MoiraiCache: moiraiCache,
	}
}

func (p *Preemptor) Preempt(ctx context.Context, pod *corev1.Pod, m framework.NodeToStatusMap) (*framework.PostFilterResult, *framework.Status) {
	return framework.NewPostFilterResultWithNominatedNode(p.PluginName), framework.NewStatus(framework.Success)
}
