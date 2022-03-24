package moirai

import (
	corev1helpers "k8s.io/component-helpers/scheduling/corev1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// Less are used to sort pods in the scheduling queue.
func (p *Plugin) Less(podInfo1, podInfo2 *framework.QueuedPodInfo) bool {
	p1 := corev1helpers.PodPriority(podInfo1.Pod)
	p2 := corev1helpers.PodPriority(podInfo2.Pod)
	return p1 > p2 || (p1 == p2 && podInfo1.Timestamp.Before(podInfo2.Timestamp))
}
