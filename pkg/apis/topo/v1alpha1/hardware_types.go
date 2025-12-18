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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	HardwareKindHardware            = "Hardware"
	HardwarePluralHardware          = "hardwares"
	HardwareSingularHardware        = "Hardware"
	HardwareNamespaceScopedHardware = false
	HardwareShortName               = "hw"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Cluster,shortName=hw
// +kubebuilder:subresource:status

// Hardware represents the hardware information of a cluster.
type Hardware struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HardwareSpec   `json:"spec,omitempty"`
	Status HardwareStatus `json:"status,omitempty"`
}

// HardwareSpec represents the specification of the Hardware.
type HardwareSpec struct {
	// +required
	ClusterName string `json:"clusterName"`
	// +optional
	Region string `json:"region,omitempty"`
	// +optional
	Provider string `json:"provider,omitempty"`
	// +required
	Hardwares *Hardwares `json:"hardwares"`
}

// Hardwares represents the hardware information of a cluster.
type Hardwares struct {
	// +optional
	Entropy *EntropyInfo `json:"entropy"`
	// +optional
	GpuType []GpuTypeInfo `json:"gpuType"`
}

// EntropyInfo represents the entropy information of a cluster.
type EntropyInfo struct {
	// +required
	Mismatch int64 `json:"mismatch"`
	// +required
	Sparsity int64 `json:"sparsity"`
	// +required
	ZoneMismatch int64 `json:"zone_mismatch"`
}

// GpuTypeInfo represents the GPU type information.
type GpuTypeInfo struct {
	// +optional
	Name string `json:"name"`
	// +optional
	Allocatable *corev1.ResourceList `json:"allocatable,omitempty"`
	// +optional
	Requests *corev1.ResourceList `json:"requests,omitempty"`
}

// HardwareStatus represents the status of Hardware.
type HardwareStatus struct {
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty"`
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// HardwareList contains a list of Hardware
type HardwareList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items holds a list of Hardware.
	Items []Hardware `json:"items"`
}
