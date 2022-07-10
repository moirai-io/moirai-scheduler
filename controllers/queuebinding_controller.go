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
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"

	moirai "github.com/rudeigerc/moirai/apis/scheduling/v1alpha1"
)

// QueueBindingReconciler reconciles a QueueBinding object
type QueueBindingReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Handlers []QueueBindingEventHandler
}

// NewQueueBindingReconciler returns a new QueueBindingReconciler.
func NewQueueBindingReconciler(
	client client.Client,
	scheme *runtime.Scheme,
	handlers ...QueueBindingEventHandler,
) *QueueBindingReconciler {
	return &QueueBindingReconciler{
		Client:   client,
		Scheme:   scheme,
		Handlers: handlers,
	}
}

//+kubebuilder:rbac:groups="",resources=events,verbs=create;watch;update
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queuebindings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queuebindings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queuebindings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *QueueBindingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var queueBindingObj moirai.QueueBinding
	if err := r.Get(ctx, req.NamespacedName, &queueBindingObj); err != nil {
		log.Error(err, "unable to fetch QueueBinding")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log = ctrl.LoggerFrom(ctx).WithValues("queuebinding", klog.KObj(&queueBindingObj))
	ctx = ctrl.LoggerInto(ctx, log)

	log.V(2).Info("Reconciling QueueBinding")

	if err := r.Status().Update(ctx, &queueBindingObj); err != nil {
		log.Error(err, "unable to update QueueBinding status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// QueueBindingEventHandler is an interface that handles QueueBinding events.
type QueueBindingEventHandler interface {
	// HandleQueueBindingUpdateEvent handles the QueueBinding update events.
	HandleQueueBindingUpdateEvent(*moirai.QueueBinding)
}

// handleEvents handles QueueBinding events for handlers.
func (r *QueueBindingReconciler) handleEvents(queueBinding *moirai.QueueBinding) {
	for _, handler := range r.Handlers {
		handler.HandleQueueBindingUpdateEvent(queueBinding)
	}
}

// Create returns true if the Create event should be processed
func (r *QueueBindingReconciler) Create(e event.CreateEvent) bool {
	queueBinding := e.Object.(*moirai.QueueBinding)
	defer r.handleEvents(queueBinding)
	return true
}

// Delete returns true if the Delete event should be processed
func (r *QueueBindingReconciler) Delete(e event.DeleteEvent) bool {
	queueBinding := e.Object.(*moirai.QueueBinding)
	defer r.handleEvents(queueBinding)
	return true
}

// Update returns true if the Update event should be processed
func (r *QueueBindingReconciler) Update(e event.UpdateEvent) bool {
	queueBindingOld := e.ObjectOld.(*moirai.QueueBinding)
	queueBindingNew := e.ObjectNew.(*moirai.QueueBinding)

	defer r.handleEvents(queueBindingOld)
	defer r.handleEvents(queueBindingNew)

	return true
}

// Generic returns true if the Generic event should be processed
func (r *QueueBindingReconciler) Generic(e event.GenericEvent) bool {
	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *QueueBindingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&moirai.QueueBinding{}).
		WithEventFilter(r).
		Complete(r)
}
