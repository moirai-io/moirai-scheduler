package moirai

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	moirai "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
)

// PostBind Plugin

// PostBind is called after a pod is successfully bound.
func (p *Plugin) PostBind(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) {
	klog.V(5).InfoS("PostBind extension point", "pod", klog.KObj(pod))

	queueBinding, err := p.manager.GetQueueBinding(ctx, pod)
	if err != nil {
		klog.Errorf("unable to get QueueBinding: %v", err)
		return
	}
	queueBindingCopy := queueBinding.DeepCopy()

	queueBindingCopy.Status.Scheduled++
	queueBindingCopy.Status.Pending--

	if queueBindingCopy.Status.Phase == moirai.QueueBindingPhaseTypeReady && queueBindingCopy.Status.Pending == 0 {
		queueBindingCopy.Status.Phase = moirai.QueueBindingPhaseTypePending
	}

	if queueBindingCopy.Status.Pending == 0 {
		queueBindingCopy.Status.Phase = moirai.QueueBindingPhaseTypeScheduled
	}

	if !equality.Semantic.DeepEqual(queueBindingCopy.Status, queueBindingCopy.Status) {
		err := p.client.Status().Update(ctx, queueBindingCopy)
		klog.Errorf("unable to update the status of QueueBinding: %v", err)
	}
}
