package registry

import (
	"github.com/moirai-io/moirai-scheduler/pkg/scheduler"

	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

// NewRegistryOptions returns the default options for the registry.
func NewRegistryOptions() app.Option {
	return app.WithPlugin(scheduler.Name, scheduler.New)
}
