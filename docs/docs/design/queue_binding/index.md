# QueueBinding

```go
type QueueBindingSpec struct {
	// Queue is the name of the queue to bind to
	Queue string `json:"queue,omitempty"`
	// PriorityClassName is the name of the priority class to bind to
	PriorityClassName string                 `json:"priorityClassName,omitempty"`
	Resources         corev1.ResourceList    `json:"resources,omitempty"`
	JobRef            corev1.ObjectReference `json:"jobRef,omitempty"`
}
```
