/*
Copyright 2021 The Karmada Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindMesh            = "Mesh"
	ResourcePluralMesh          = "meshs"
	ResourceSingularMesh        = "mesh"
	ResourceNamespaceScopedMesh = false
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Cluster,shortName=mesh,categories={karmada-io}
// +kubebuilder:subresource:status

type Mesh struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// Spec represents the desired behavior of Mesh.
	Spec MeshSpec `json:"spec"`
}

type MeshSpec struct {
	// +required
	Plugins []Plugin `json:"plugins"`
}

type Plugin struct {
	// +required
	Name string `json:"name"`
	// +required
	Version string `json:"version"`
	// +required
	Path string `json:"path"`
	// +required
	Symbol string `json:"symbol"`
	// +optional
	Description string `json:"description"`
	// +required
	Enabled bool `json:"enabled"`
	// +required
	Graceful int `json:"graceful"`
	// +required
	Md5 string `json:"md5"`
	// +required
	Config string `json:"config"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type MeshList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Mesh `json:"items"`
}
