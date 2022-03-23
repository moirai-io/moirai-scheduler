package registry

import (
	"github.com/moirai-io/moirai-scheduler/pkg/plugins/moirai"

	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

// NewRegistryOptions returns the default options for the registry.
func NewRegistryOptions() app.Option {
	return app.WithPlugin(moirai.Name, moirai.New)
}
