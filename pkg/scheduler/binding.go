package scheduler

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// PostBind Plugin

// PostBind...
func (p *Plugin) PostBind(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) {
	klog.V(5).InfoS("PostBind extension point", "pod", klog.KObj(pod))

	_, err := p.manager.GetQueueBinding(ctx, pod)
	if err != nil {
		klog.Errorf("unable to get QueueBinding: %v", err)
	}
}
