package moirai

import "k8s.io/kubernetes/pkg/scheduler/framework"

type stateData struct {
	name string
}

// NewStateData returns a new state data object.
func NewStateData(name string) framework.StateData {
	return &stateData{
		name: name,
	}
}

func (d *stateData) Clone() framework.StateData {
	return d
}
