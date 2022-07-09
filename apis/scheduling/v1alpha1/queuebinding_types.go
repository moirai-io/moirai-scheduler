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

// QueueBindingSpec defines the desired state of QueueBinding
type QueueBindingSpec struct {
	// Queue is the name of the queue to bind to
	Queue string `json:"queue,omitempty"`
	// PriorityClassName is the name of the priority class to bind to
	PriorityClassName string                 `json:"priorityClassName,omitempty"`
	Resources         corev1.ResourceList    `json:"resources,omitempty"`
	JobRef            corev1.ObjectReference `json:"jobRef,omitempty"`
}

// QueueBindingConditionType defines the condition type of QueueBinding
type QueueBindingConditionType string

// QueueBindingCondition defines the observed conditions of QueueBinding
type QueueBindingCondition struct {
	Type               QueueBindingConditionType `json:"type,omitempty"`
	Status             corev1.ConditionStatus    `json:"status,omitempty"`
	LastTransitionTime metav1.Time               `json:"lastTransitionTime,omitempty"`
	Reason             string                    `json:"reason,omitempty"`
	Message            string                    `json:"message,omitempty"`
}

// QueueBindingPhaseType defines the observed phase of QueueBinding
type QueueBindingPhaseType string

const (
	// QueueBindingPhaseTypeReady represents the phase that the QueueBinding is ready
	QueueBindingPhaseTypeReady QueueBindingPhaseType = "Ready"
	// QueueBindingPhaseTypePending represents the phase that one of the Pod belonging to the QueueBinding is scheduled
	QueueBindingPhaseTypePending QueueBindingPhaseType = "Pending"
	// QueueBindingPhaseTypeScheduled represents the phase that all the Pods belonging to the QueueBinding is scheduled
	QueueBindingPhaseTypeScheduled QueueBindingPhaseType = "Scheduled"
	// QueueBindingPhaseTypeFailed represents the phase that one of the Pod belonging to the QueueBinding is failed to schedule
	QueueBindingPhaseTypeFailed QueueBindingPhaseType = "Failed"
)

// QueueBindingStatus defines the observed state of QueueBinding
type QueueBindingStatus struct {
	// Phase represents the current phase of the QueueBinding
	Phase QueueBindingPhaseType `json:"phase,omitempty"`
	// Pending represents the number of pending Pods in the QueueBinding
	Pending int32 `json:"pending,omitempty"`
	// Scheduled represents the number of scheduled Pods in the QueueBinding
	Scheduled int32 `json:"scheduled,omitempty"`
	// Conditions represents the current conditions of the QueueBinding
	Conditions []QueueBindingCondition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Job",type="string",JSONPath=".spec.jobRef.name"

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
