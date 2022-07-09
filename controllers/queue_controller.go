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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	moirai "github.com/rudeigerc/moirai/apis/scheduling/v1alpha1"
)

// QueueReconciler reconciles a Queue object
type QueueReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
	sourceCh chan event.GenericEvent
}

// NewQueueReconciler returns a new QueueReconciler.
func NewQueueReconciler(
	client client.Client,
	scheme *runtime.Scheme,
	log logr.Logger,
	recorder record.EventRecorder,
) *QueueReconciler {
	return &QueueReconciler{
		Client:   client,
		Log:      log,
		Recorder: recorder,
		Scheme:   scheme,
		sourceCh: make(chan event.GenericEvent, 10),
	}
}

//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queues,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queues/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=scheduling.moirai.io,resources=queues/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *QueueReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var queueObj moirai.Queue
	if err := r.Get(ctx, req.NamespacedName, &queueObj); err != nil {
		log.Error(err, "unable to fetch Queue")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.V(2).Info("Reconciling Queue")

	if err := r.Status().Update(ctx, &queueObj); err != nil {
		log.Error(err, "unable to update Queue status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *QueueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&moirai.Queue{}).
		Watches(&source.Channel{Source: r.sourceCh}, &queueBindingEventHandler{}).
		Complete(r)
}

type queueBindingEventHandler struct{}

func (h *queueBindingEventHandler) Create(event.CreateEvent, workqueue.RateLimitingInterface) {
}

func (h *queueBindingEventHandler) Update(event.UpdateEvent, workqueue.RateLimitingInterface) {
}

func (h *queueBindingEventHandler) Delete(event.DeleteEvent, workqueue.RateLimitingInterface) {
}

func (h *queueBindingEventHandler) Generic(e event.GenericEvent, q workqueue.RateLimitingInterface) {
	queueBinding := e.Object.(*moirai.QueueBinding)
	if queueBinding.Name == "" {
		return
	}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      queueBinding.Spec.Queue,
			Namespace: queueBinding.Namespace,
		},
	}
	q.AddAfter(req, time.Second)
}
