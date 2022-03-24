package moirai

import (
	"context"
	"fmt"
	"math"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"

	moirai "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
)

// PreFilterExtensions is an interface that is included in plugins that allow specifying
// callbacks to make incremental updates to its supposedly pre-calculated
// state.
func (p *Plugin) PreFilterExtensions() framework.PreFilterExtensions {
	return nil
}

// PreFilterPlugin

// PreFilter is called at the beginning of the scheduling cycle. All PreFilter
// plugins must return success or the pod will be rejected.
func (p *Plugin) PreFilter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod) *framework.Status {
	klog.V(5).InfoS("PreFilter extension point", "pod", klog.KObj(pod))

	queueBinding, err := p.manager.GetQueueBinding(ctx, pod)
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("unable to get QueueBinding: %v", err))
	}
	// fetch pods according to the queue binding
	_, err = p.frameworkHandler.SharedInformerFactory().Core().V1().Pods().Lister().List(
		labels.SelectorFromSet(labels.Set{moirai.QueueBindingLabel: queueBinding.Name}),
	)
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("unable to list pods in QueueBinding %s: %v", queueBinding.Name, err))
	}

	nodeInfoList, err := p.frameworkHandler.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("unable to list nodes: %v", err))
	}
	nodeList := make([]*corev1.Node, 0, len(nodeInfoList))
	for _, nodeInfo := range nodeInfoList {
		nodeList = append(nodeList, nodeInfo.Node())
	}

	// resources := queueBinding.Spec.Resources.DeepCopy()

	// FIXME:
	state.Write("", NewStateData(queueBinding.Name))
	return framework.NewStatus(framework.Success, "")
}

// Filter Plugin

// Filter is called by the scheduling framework.
func (p *Plugin) Filter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	klog.V(5).InfoS("Filter extension point", "pod", klog.KObj(pod))

	node := nodeInfo.Node()
	if node == nil {
		return framework.AsStatus(fmt.Errorf("node not found"))
	}

	_, err := p.manager.GetQueueBinding(ctx, pod)
	if err != nil {
		return framework.NewStatus(framework.Error, fmt.Sprintf("unable to get QueueBinding: %v", err))
	}

	return framework.NewStatus(framework.Success, "")
}

// PostFilter Plugin

// PostFilter is called by the scheduling framework.
// Preemption
func (p *Plugin) PostFilter(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, filteredNodeStatusMap framework.NodeToStatusMap) (*framework.PostFilterResult, *framework.Status) {
	klog.V(5).InfoS("PostFilter extension point", "pod", klog.KObj(pod))

	return nil, framework.NewStatus(framework.Success, "")
}

// PreScore Plugin

// PreScore is called by the scheduling framework after a list of nodes
// passed the filtering phase. All prescore plugins must return success or
// the pod will be rejected
func (p *Plugin) PreScore(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodes []*corev1.Node) *framework.Status {
	klog.V(5).InfoS("PreScore extension point", "pod", klog.KObj(pod))

	return framework.NewStatus(framework.Success, "")
}

// Score Plugin

// Score is called on each filtered node. It must return success and an integer
// indicating the rank of the node. All scoring plugins must return success or
// the pod will be rejected.
func (p *Plugin) Score(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) (int64, *framework.Status) {
	klog.V(5).InfoS("Score extension point", "pod", klog.KObj(pod))

	nodeInfo, err := p.frameworkHandler.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.AsStatus(fmt.Errorf("unable to get node %s from snapshot: %v", nodeName, err))
	}
	_ = nodeInfo.Node()
	return 0, framework.NewStatus(framework.Success, "")
}

// Score Extension

// ScoreExtensions is an interface for Score extended functionality.
func (p *Plugin) ScoreExtensions() framework.ScoreExtensions {
	return p
}

// NormalizeScore is called for each scored node.
// From https://github.com/Azure/placement-policy-scheduler-plugins/blob/main/pkg/plugins/placementpolicy/placementpolicy.go
func (p *Plugin) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
	// Find highest and lowest scores.
	var highest int64 = -math.MaxInt64
	var lowest int64 = math.MaxInt64
	for _, nodeScore := range scores {
		if nodeScore.Score > highest {
			highest = nodeScore.Score
		}
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
	}

	// Transform the highest to lowest score range to fit the framework's min to max node score range.
	oldRange := highest - lowest
	newRange := framework.MaxNodeScore - framework.MinNodeScore
	for i, nodeScore := range scores {
		if oldRange == 0 {
			scores[i].Score = framework.MinNodeScore
		} else {
			scores[i].Score = ((nodeScore.Score - lowest) * newRange / oldRange) + framework.MinNodeScore
		}
	}

	klog.InfoS("normalized scores", "pod", pod.Name, "scores", scores)
	return nil
}

// Reserve Plugin

// Reserve is called by the scheduling framework when the scheduler cache is
// updated.
func (p *Plugin) Reserve(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) *framework.Status {
	klog.V(5).InfoS("Reserve extension point", "pod", klog.KObj(pod))

	return framework.NewStatus(framework.Success, "")
}

// Unreserve is called by the scheduling framework when a reserved pod was
// rejected, an error occurred during reservation of subsequent plugins, or
// in a later phase.
func (p *Plugin) Unreserve(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) {
	klog.V(5).InfoS("Unreserve extension point", "pod", klog.KObj(pod))

	queueBinding, err := p.manager.GetQueueBinding(ctx, pod)
	if err != nil {
		klog.Errorf("unable to get QueueBinding: %v", err)
		return
	}

	// Iterate waiting pods to reject pods belonging to the same queue binding to be scheduled
	p.frameworkHandler.IterateOverWaitingPods(func(waitingPod framework.WaitingPod) {
		waitingPodRef := waitingPod.GetPod()
		if waitingPodRef.Namespace == pod.Namespace && waitingPodRef.Labels[moirai.QueueBindingLabel] == queueBinding.Name {
			waitingPod.Reject(p.Name(), fmt.Sprintf("Pod %s is rejected in unreserve phase", pod.Name))
		}
	})
}

// Permit Plugin

// Permit is called before binding a pod (and before prebind plugins). Permit
// plugins are used to prevent or delay the binding of a Pod.
func (p *Plugin) Permit(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) (*framework.Status, time.Duration) {
	klog.V(5).InfoS("Permit extension point", "pod", klog.KObj(pod))

	queueBinding, err := p.manager.GetQueueBinding(ctx, pod)
	if err != nil {
		klog.Errorf("unable to get QueueBinding: %v", err)
		return framework.NewStatus(framework.UnschedulableAndUnresolvable, fmt.Sprintf("unable to get QueueBinding: %v", err)), 0
	}
	if queueBinding.Name == "" {
		return framework.NewStatus(framework.Success, ""), 0
	}

	if queueBinding.Status.Pending > 0 {
		// FIXME: Use configuration instead of hardcoded value
		return framework.NewStatus(framework.Wait, ""), time.Second * 30
	}

	// Iterate waiting pods to allow pods belonging to the same queue binding to be scheduled
	p.frameworkHandler.IterateOverWaitingPods(func(waitingPod framework.WaitingPod) {
		waitingPodRef := waitingPod.GetPod()
		if waitingPodRef.Namespace == pod.Namespace && waitingPodRef.Labels[moirai.QueueBindingLabel] == queueBinding.Name {
			klog.V(3).InfoS("Allowed in permit phase", "pod", klog.KObj(waitingPodRef))
			waitingPod.Allow(p.Name())
		}
	})

	return framework.NewStatus(framework.Success, ""), 0
}
