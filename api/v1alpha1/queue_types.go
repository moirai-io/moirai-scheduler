/*
Copyright 2021 Yuchen Cheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SuspendAnnotation string = "moirai.io/suspend"
	QueueLabel        string = "moirai.io/queue"
)

// QueueSpec defines the desired state of Queue
type QueueSpec struct {
	Capacity corev1.ResourceList `json:"capacity,omitempty"`
}

// QueueConditionType defines the condition type of Queue
type QueueConditionType string

// QueueCondition defines the observed state of Queue
type QueueCondition struct {
	Type               QueueConditionType     `json:"type,omitempty"`
	Status             corev1.ConditionStatus `json:"status,omitempty"`
	LastTransitionTime metav1.Time            `json:"lastTransitionTime,omitempty"`
	Reason             string                 `json:"reason,omitempty"`
	Message            string                 `json:"message,omitempty"`
}

// QueueState defines the state of Queue
type QueueState string

const (
	// QueueStateReady represents the state that the Queue is ready
	QueueStateReady QueueState = "Ready"
	// QueueStatePending represents the state that that the Queue is unavailable
	QueueStateUnavailable QueueState = "Unavailable"
)

// QueueStatus defines the observed state of Queue
type QueueStatus struct {
	State      QueueState           `json:"state,omitempty"`
	Used       corev1.ResourceList  `json:"used,omitempty"`
	Conditions []QueueConditionType `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="CPU",type="string",JSONPath=".spec.resources.cpu"
//+kubebuilder:printcolumn:name="Memory",type="string",JSONPath=".spec.resources.memory"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Queue is the Schema for the queues API
type Queue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QueueSpec   `json:"spec,omitempty"`
	Status QueueStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QueueList contains a list of Queue
type QueueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Queue `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Queue{}, &QueueList{})
}
