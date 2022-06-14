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

// Package v1beta3 is the v1beta3 version of the API
package v1beta3

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kube-scheduler/config/v1beta3"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = v1beta3.SchemeGroupVersion

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &v1beta3.SchemeBuilder

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// addKnownTypes registers known types to the given scheme
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&MoiraiArgs{},
	)
	return nil
}

func init() {
	SchemeBuilder.Register(addKnownTypes)
}