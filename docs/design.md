# Design

## Motivation

1. Multi-tenancy
2. Fair-sharing
3. Resource utilization

## APIs

`apis/scheduling/v1alpha1`

### Queue

```go
type QueueSpec struct {
	Capacity corev1.ResourceList `json:"capacity,omitempty"`
}
```

### QueueBinding

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

## Extension Points

### PreFilter

> These plugins are used to pre-process info about the Pod, or to check certain conditions that the cluster or the Pod must meet. If a PreFilter plugin returns an error, the scheduling cycle is aborted.

### Filter

> These plugins are used to filter out nodes that cannot run the Pod. For each node, the scheduler will call filter plugins in their configured order. If any filter plugin marks the node as infeasible, the remaining plugins will not be called for that node. Nodes may be evaluated concurrently.

### PostFilter

> These plugins are called after Filter phase, but only when no feasible nodes were found for the pod. Plugins are called in their configured order. If any postFilter plugin marks the node as Schedulable, the remaining plugins will not be called. A typical PostFilter implementation is preemption, which tries to make the pod schedulable by preempting other Pods.

### PreScore

> These plugins are used to perform "pre-scoring" work, which generates a sharable state for Score plugins to use. If a PreScore plugin returns an error, the scheduling cycle is aborted.

### Score

> These plugins are used to rank nodes that have passed the filtering phase. The scheduler will call each scoring plugin for each node. There will be a well defined range of integers representing the minimum and maximum scores. After the NormalizeScore phase, the scheduler will combine node scores from all plugins according to the configured plugin weights.

### Reserve

> A plugin that implements the Reserve extension has two methods, namely `Reserve` and `Unreserve`, that back two informational scheduling phases called Reserve and Unreserve, respectively. Plugins which maintain runtime state (aka "stateful plugins") should use these phases to be notified by the scheduler when resources on a node are being reserved and unreserved for a given Pod.

#### `Reserve`

nil

#### `Unreserve`

1. Deny all Pods belonging to the same QueueBinding in the scheduling queue.

## Permit

> Permit plugins are invoked at the end of the scheduling cycle for each Pod, to prevent or delay the binding to the candidate node.

## Limitation

## Referenced Projects

- [kubernetes/kube-scheduler](https://github.com/kubernetes/kube-scheduler)
- [kubernetes-sigs/kueue](https://github.com/kubernetes-sigs/kueue)
- [kube-queue/kube-queue](https://github.com/kube-queue/kube-queue)
- [koordinator-sh/koordinator](https://github.com/koordinator-sh/koordinator)
- [volcano-sh/volcano](https://github.com/volcano-sh/volcano)
- [apache/incubator-yunikorn-core](https://github.com/apache/incubator-yunikorn-core)
