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
	HardwareNodeKindHardwareNode            = "HardwareNode"
	HardwareNodePluralHardwareNode          = "hardwarenodes"
	HardwareNodeSingularHardwareNode        = "HardwareNode"
	HardwareNodeNamespaceScopedHardwareNode = false
	HardwareNodeShortName                   = "hd"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:storageversion
// +kubebuilder:resource:scope=Cluster,shortName=hd
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.cluster"
// +kubebuilder:printcolumn:name="GpuType",type="string",JSONPath=".spec.host.gpuType"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// HardwareNode represents the hardware information of a specific node in a cluster.
type HardwareNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HardwareNodeSpec   `json:"spec,omitempty"`
	Status HardwareNodeStatus `json:"status,omitempty"`
}

// HardwareNodeSpec represents the specification of the HardwareNode.
type HardwareNodeSpec struct {
	// +required
	Host *HostInfo `json:"host"`
	// +optional
	Cluster string `json:"cluster,omitempty"`
}

// HostInfo represents the detailed hardware information of a host.
type HostInfo struct {
	// +required
	Name string `json:"name"`
	// +optional
	GpuType string `json:"gpuType,omitempty"`
	// +optional
	CudaVersion string `json:"cudaVersion,omitempty"`
	// +optional
	GpuFailure []string `json:"gpuFailure,omitempty"`
	// +optional
	Gpu *GpuInfo `json:"gpu,omitempty"`
	// +required
	Cpu *CpuInfo `json:"cpu"`
	// +required
	Mem *MemInfo `json:"mem"`
	// +optional
	Gpus []GpuDetail `json:"gpus"`
}

// GpuDetail represents the information of a GPU device.
type GpuDetail struct {
	// +required
	Index int64 `json:"index"`
	// +required
	Idle bool `json:"idle"`
	// +optional
	Process []string `json:"process,omitempty"`
}

type GpuInfo struct {
	Total int64 `json:"total"`
	Idle  int64 `json:"idle"`
}

// CpuInfo represents the CPU information of a host.
type CpuInfo struct {
	// +required
	Total int64 `json:"total"`
	// +required
	Idle int64 `json:"idle"`
}

// MemInfo represents the memory information of a host.
type MemInfo struct {
	// +required
	Total int64 `json:"total"`
	// +required
	Idle int64 `json:"idle"`
}

// HardwareNodeStatus represents the status of HardwareNode.
type HardwareNodeStatus struct {
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty"`
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// HardwareNodeList contains a list of HardwareNode
type HardwareNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items holds a list of HardwareNode.
	Items []HardwareNode `json:"items"`
}
