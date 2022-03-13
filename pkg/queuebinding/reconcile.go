package queuebinding

import (
	"github.com/go-logr/logr"
	schedulingv1alpha1 "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	client    client.Client
	log       logr.Logger
	recorder  record.EventRecorder
	scheme    *runtime.Scheme
	instance  *schedulingv1alpha1.QueueBinding
	generator *generator
}

func NewReconciler(
	client client.Client,
	log logr.Logger,
	recorder record.EventRecorder,
	scheme *runtime.Scheme,
	instance *schedulingv1alpha1.QueueBinding,
) (*Reconciler, error) {
	generator := newGenerator(instance)

	return &Reconciler{
		client:    client,
		log:       log,
		recorder:  recorder,
		scheme:    scheme,
		instance:  instance,
		generator: generator,
	}, nil
}

func (r *Reconciler) Reconcile() error {
	r.log.Info("Reconciling QueueBinding")

	return nil
}
