package manager

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// GetAvailableNodeResources returns the available resources of the node
func (m *MoiraiManager) GetAvailableNodeResources(nodeInfo *framework.NodeInfo, queueBindingName string) {
	_ = nodeInfo.Clone()
}

// CheckResources checks the resource of the node
func (m *MoiraiManager) CheckResources(nodeList []*framework.NodeInfo, resource corev1.ResourceList) {
	for _, nodeInfo := range nodeList {
		if nodeInfo == nil {
			continue
		}
		node := nodeInfo.Node()
		if node == nil {
			continue
		}

	}
}

// GetNodeAvaliableResource returns the available resource of the node
func (m *MoiraiManager) GetNodeAvaliableResource() *framework.Resource {
	return nil
}
