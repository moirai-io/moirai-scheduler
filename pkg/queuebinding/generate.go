package queuebinding

import (
	schedulingv1alpha1 "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
)

type generator struct {
	queueBinding *schedulingv1alpha1.QueueBinding
}

func newGenerator(queueBinding *schedulingv1alpha1.QueueBinding) *generator {
	return &generator{
		queueBinding: queueBinding,
	}
}
