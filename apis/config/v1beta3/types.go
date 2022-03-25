package v1beta3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MoiraiArgs defines the scheduling parameters for Moirai plugin
type MoiraiArgs struct {
	metav1.TypeMeta `json:",inline"`

	// FIXME:
	Name string `json:"name,omitempty"`
}
