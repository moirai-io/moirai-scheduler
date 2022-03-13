package client

import (
	schedulingv1alpha1 "github.com/moirai-io/moirai-scheduler/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// NewWatchClient returns a new client for watching objects.
func NewWatchClient() (client.WithWatch, error) {
	scheme := runtime.NewScheme()
	if err := schedulingv1alpha1.AddToScheme(scheme); err != nil {
		klog.Error(err)
		return nil, err
	}

	moiraiClient, err := client.NewWithWatch(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return moiraiClient, nil
}
