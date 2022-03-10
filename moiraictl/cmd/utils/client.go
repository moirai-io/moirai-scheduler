package utils

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	schedulingv1alpha1 "github.com/moirai-io/moirai-operator/api/v1alpha1"
)

var (
	scheme = runtime.NewScheme()
)

func NewClient() (client.Client, error) {
	err := schedulingv1alpha1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	c, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return c, nil
}
