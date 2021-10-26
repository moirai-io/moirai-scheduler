package queue

import (
	schedulingv1alpha1 "github.com/moirai-io/moirai/api/v1alpha1"
)

type generator struct {
	queue *schedulingv1alpha1.Queue
}

func newGenerator(queue *schedulingv1alpha1.Queue) *generator {
	return &generator{
		queue: queue,
	}
}
