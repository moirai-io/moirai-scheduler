/*
Copyright 2021 Yuchen Cheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package batch

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	moirai "github.com/rudeigerc/moirai/apis/scheduling/v1alpha1"
)

// JobReconciler reconciles a Job object
type JobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch,resources=jobs/finalizers,verbs=update
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queuebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queuebindings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queuebindings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *JobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var jobObj batchv1.Job
	if err := r.Get(ctx, req.NamespacedName, &jobObj); err != nil {
		klog.Error(err, "unable to fetch Job")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log = ctrl.LoggerFrom(ctx).WithValues("job", klog.KObj(&jobObj))
	ctx = ctrl.LoggerInto(ctx, log)

	if queueName(&jobObj) == "" {
		return ctrl.Result{}, nil
	}

	log.V(2).Info("Reconciling Job")

	var queueBindingList moirai.QueueBindingList
	if err := r.List(ctx, &queueBindingList, client.InNamespace(req.Namespace), client.MatchingFields{ownerKey: req.Name}); err != nil {
		klog.Error(err, "unable to list QueueBindings")
		return ctrl.Result{}, err
	}

	// create a new QueueBinding for the Job if it doesn't exist
	if len(queueBindingList.Items) == 0 {
		err := r.createQueueBinding(ctx, &jobObj)
		if err != nil {
			klog.Error(err, "unable to create QueueBinding")

		}
		return ctrl.Result{}, err
	}

	// queueBinding := queueBindingList.Items[0]
	// if jobSuspended(&jobObj) {
	// }

	if err := r.Status().Update(ctx, &jobObj); err != nil {
		log.Error(err, "unable to update Job status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

var (
	ownerKey = ".metadata.controller"
)

// SetupWithManager sets up the controller with the Manager.
func (r *JobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &moirai.QueueBinding{}, ownerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner...
		queueBinding := rawObj.(*moirai.QueueBinding)
		owner := metav1.GetControllerOf(queueBinding)
		if owner == nil {
			return nil
		}
		// ...make sure it's a Job...
		if owner.APIVersion != "batch/v1" || owner.Kind != "Job" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Job{}).
		Owns(&moirai.QueueBinding{}).
		Complete(r)
}

func (r *JobReconciler) createQueueBinding(ctx context.Context, job *batchv1.Job) error {
	log := log.FromContext(ctx)

	queueBinding := &moirai.QueueBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "queuebinding-" + job.Name,
			Namespace: job.Namespace,
		},
		Spec: moirai.QueueBindingSpec{
			Queue: queueName(job),
			JobRef: v1.ObjectReference{
				APIVersion: "batch/v1",
				Kind:       "Job",
				Namespace:  job.Namespace,
				Name:       job.Name,
			},
		},
	}

	if err := r.Create(ctx, queueBinding); err != nil {
		return err
	}
	log.V(3).Info("Create QueueBinding", "QueueBinding", queueBinding.Name)

	if err := ctrl.SetControllerReference(job, queueBinding, r.Scheme); err != nil {
		return err
	}

	// job.Spec.Template.Labels[moirai.QueueBindingLabel] = queueBinding.Name
	// if err := r.Update(ctx, job); err != nil {
	// 	return err
	// }

	return nil
}
