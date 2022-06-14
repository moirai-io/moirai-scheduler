package moirai

import "k8s.io/kubernetes/pkg/scheduler/framework"

const PreFilterKey framework.StateKey = "moirai.io/prefilter"

type stateData struct {
	name string
	node string
}

// NewStateData returns a new state data object.
func NewStateData(name string, node string) framework.StateData {
	return &stateData{
		name: name,
		node: node,
	}
}

func (d *stateData) Clone() framework.StateData {
	return d
}
