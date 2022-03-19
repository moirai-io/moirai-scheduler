package internal

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	schedulingv1alpha1 "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(schedulingv1alpha1.AddToScheme(scheme))
}
