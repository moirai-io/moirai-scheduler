package batch

import (
	batchv1 "k8s.io/api/batch/v1"

	moirai "github.com/rudeigerc/moirai/apis/scheduling/v1alpha1"
)

func jobSuspended(job *batchv1.Job) bool {
	return job.Spec.Suspend != nil && *job.Spec.Suspend
}

func queueName(job *batchv1.Job) string {
	return job.Labels[moirai.QueueLabel]
}
