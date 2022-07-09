package internal

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	moirai "github.com/rudeigerc/moirai/apis/scheduling/v1alpha1"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(moirai.AddToScheme(scheme))
}
