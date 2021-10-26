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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	schedulingv1alpha1 "github.com/moirai-io/moirai/api/v1alpha1"
	"github.com/moirai-io/moirai/pkg/queue"
)

// QueueReconciler reconciles a Queue object
type QueueReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
}

//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queues,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queues/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queues/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Queue object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *QueueReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var original schedulingv1alpha1.Queue
	if err := r.Get(ctx, req.NamespacedName, &original); err != nil {
		log.Error(err, "unable to fetch Queue")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	instance := original.DeepCopy()

	reconciler, err := queue.NewReconciler(r.Client, log, r.Recorder, r.Scheme, instance)
	if err != nil {
		return ctrl.Result{}, err
	}
	if err := reconciler.Reconcile(); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Status().Update(ctx, instance); err != nil {
		log.Error(err, "unable to update CronJob status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *QueueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&schedulingv1alpha1.Queue{}).
		Complete(r)
}
