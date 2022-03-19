package scheduler

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// PreFilterExtensions is an interface that is included in plugins that allow specifying
// callbacks to make incremental updates to its supposedly pre-calculated
// state.
func (p *Plugin) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// PreFilter ...
func (p *Plugin) PreFilter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod) *framework.Status {
	nodeInfoList, err := p.frameworkHandler.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("failed to list nodes: %v", err))
	}
	nodeList := make([]*corev1.Node, 0, len(nodeInfoList))
	for _, nodeInfo := range nodeInfoList {
		nodeList = append(nodeList, nodeInfo.Node())
	}
	klog.V(3).InfoS("get node list", nodeList)

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
