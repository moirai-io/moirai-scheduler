package scheduler

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	schedulingv1alpha1 "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
)

// PreFilterExtensions is an interface that is included in plugins that allow specifying
// callbacks to make incremental updates to its supposedly pre-calculated
// state.
func (p *Plugin) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// PreFilter ...
func (p *Plugin) PreFilter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod) *framework.Status {
	queueBinding, err := p.manager.GetQueueBinding(ctx, pod)
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("unable to get QueueBinding: %v", err))
	}
	// fetch pods according to the queue binding
	_, err = p.frameworkHandler.SharedInformerFactory().Core().V1().Pods().Lister().List(
		labels.SelectorFromSet(labels.Set{schedulingv1alpha1.QueueBindingLabel: queueBinding.Name}),
	)
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("unable to list pod: %v", err))
	}

	nodeInfoList, err := p.frameworkHandler.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("unable to list nodes: %v", err))
	}
	nodeList := make([]*corev1.Node, 0, len(nodeInfoList))
	for _, nodeInfo := range nodeInfoList {
		nodeList = append(nodeList, nodeInfo.Node())
	}
	// FIXME:
	p.manager.AnnotatePod(ctx, pod)

	return framework.NewStatus(framework.Success, "")
}

// Filter ...
func (p *Plugin) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	if nodeInfo.Node() == nil {
		return framework.AsStatus(fmt.Errorf("node not found"))
	}

	return framework.NewStatus(framework.Success, "")
}

// PostFilter ...
// Preemption
func (p *Plugin) PostFilter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, filteredNodeStatusMap framework.NodeToStatusMap) (*framework.PostFilterResult, *framework.Status) {
	return nil, framework.NewStatus(framework.Success, "")
}
