package scheduler

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// PreScore ...
func (p *Plugin) PreScore(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodes []*corev1.Node) *framework.Status {
	return framework.NewStatus(framework.Success, "")
}

// Score ...
func (p *Plugin) Score(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, nodeName string) (int64, *framework.Status) {
	return 0, framework.NewStatus(framework.Success, "")
}

// ScoreExtensions is an interface for Score extended functionality.
func (p *Plugin) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

// NormalizeScore is called for each scored node.
func (p *Plugin) NormalizeScore(ctx context.Context, state *framework.CycleState, pod *corev1.Pod, scores framework.NodeScoreList) *framework.Status {
	return framework.NewStatus(framework.Success, "")
}
