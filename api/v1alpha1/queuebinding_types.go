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
	QueueBindingLabel string = "moirai.io/queue-binding"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// QueueBindingSpec defines the desired state of QueueBinding
type QueueBindingSpec struct {
	// Queue is the name of the queue to bind to
	Queue string `json:"queue,omitempty"`
	// PriorityClassName is the name of the priority class to bind to
	PriorityClassName string                 `json:"priorityClassName,omitempty"`
	Resources         corev1.ResourceList    `json:"resources,omitempty"`
	JobRef            corev1.ObjectReference `json:"jobRef,omitempty"`
}

type QueueBindingConditionType string

// QueueBindingStatus defines the observed state of QueueBinding
type QueueBindingCondition struct {
	Type               QueueBindingConditionType `json:"type,omitempty"`
	Status             corev1.ConditionStatus    `json:"status,omitempty"`
	LastTransitionTime metav1.Time               `json:"lastTransitionTime,omitempty"`
	Reason             string                    `json:"reason,omitempty"`
	Message            string                    `json:"message,omitempty"`
}

// QueueBindingStatus defines the observed phase of QueueBinding
type QueueBindingPhaseType string

const (
	Pending   QueueBindingPhaseType = "Pending"
	Scheduled QueueBindingPhaseType = "Scheduled"
	Failed    QueueBindingPhaseType = "Failed"
)

// QueueBindingStatus defines the observed state of QueueBinding
type QueueBindingStatus struct {
	Conditions []QueueBindingCondition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:subresource:status

// QueueBinding is the Schema for the queuebindings API
type QueueBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QueueBindingSpec   `json:"spec,omitempty"`
	Status QueueBindingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// QueueBindingList contains a list of QueueBinding
type QueueBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []QueueBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&QueueBinding{}, &QueueBindingList{})
}
