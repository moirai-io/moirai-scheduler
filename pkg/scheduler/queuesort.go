package scheduler

import (
	corev1helpers "k8s.io/component-helpers/scheduling/corev1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// Less are used to sort pods in the scheduling queue.
func (p *Plugin) Less(podInfo1, podInfo2 *framework.QueuedPodInfo) bool {
	priority1 := corev1helpers.PodPriority(podInfo1.Pod)
	priority2 := corev1helpers.PodPriority(podInfo2.Pod)
	if priority1 != priority2 {
		return priority1 > priority2
	}

	timestamp1 := podInfo1.Timestamp
	timestamp2 := podInfo2.Timestamp
	return timestamp1.Before(timestamp2)
}
