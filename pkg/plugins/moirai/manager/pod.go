package manager

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AnnotatePod sets the annotation to the specified pod
func (m *MoiraiManager) AnnotatePod(ctx context.Context, pod *corev1.Pod) {
	annotations := map[string]string{}
	if pod.Annotations != nil {
		annotations = pod.Annotations
	}
	// FIXME:
	annotations["test"] = "test"
	pod.Annotations = annotations
	m.KubeClient.CoreV1().Pods(pod.Namespace).Update(ctx, pod, metav1.UpdateOptions{})
}
