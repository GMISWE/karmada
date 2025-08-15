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
	HardwareNamespaceScopedHardware = true
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:storageversion

// Hardware represents the desired state and status of a member cluster.
type Hardware struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HardwareSpec   `json:"spec,omitempty"`
	Status HardwareStatus `json:"status,omitempty"`
}

// HardwareSpec represents the specification of the desired behavior of Hardware.
type HardwareSpec struct {
	// +required
	ClusterName string `json:"clusterName"`
	// +required
	Timestamp int64 `json:"timestamp"`
	// +optional
	Provider string `json:"provider,omitempty"`
	// +optional
	Region string `json:"region,omitempty"`
	// +optional
	Zone string `json:"zone,omitempty"`
	// +optional
	Zones []string `json:"zones,omitempty"`
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty"`
	// +optional
	Hardwares *Hardwares `json:"hardwares,omitempty"`
}

type Hardwares struct {
	// +required
	Nodes *TopoNodes `json:"nodes"`
	// +required
	Entropy *TopoEntropy `json:"entropy"`
}

type TopoEntropy struct {
	// +required
	Sparsity int64 `json:"sparsity"` // 资源稀疏度
	// +required
	MisMatch int64 `json:"mismatch"` // 资源错配度
	// +required
	ZoneMisMatch int64 `json:"zone_mismatch"` // 区域错配度
}

type TopoNodes struct {
	// +required
	Num int64 `json:"num"`
	// +required
	Gpus map[GpuType]NodeGpu `json:"gpus"`
	// +required
	Cpu *NodeCpu `json:"cpu"`
	// +required
	Mem *NodeMem `json:"mem"`
}

type NodeGpu struct {
	// +required
	Mem int64 `json:"mem"`
	// +required
	Total int64 `json:"total"`
	// +required
	Idle int64 `json:"idle"`
}

type NodeCpu struct {
	// +required
	Total int64 `json:"total"`
	// +required
	Usage int64 `json:"usage"`
}

type NodeMem struct {
	// +required
	Total int64 `json:"total"`
	// +required
	Usage int64 `json:"usage"`
}

// HardwareStatus represents the status of Hardware.
type HardwareStatus struct {
	Phase      string             `json:"phase,omitempty"`
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

type GpuType string

const (
	// NVIDIA 数据中心GPU
	GpuNvGB300 GpuType = "NVIDIA GB300"
	GpuNvGB200 GpuType = "NVIDIA GB200"
	GpuNvH200  GpuType = "NVIDIA H200"
	GpuNvH100  GpuType = "NVIDIA H100"
	GpuNvH20   GpuType = "NVIDIA H20"
	GpuNvA100  GpuType = "NVIDIA A100"
	GpuNvA800  GpuType = "NVIDIA A800"
	GpuNvA40   GpuType = "NVIDIA A40"
	GpuNvA30   GpuType = "NVIDIA A30"
	GpuNvA10   GpuType = "NVIDIA A10"
	GpuNvL40   GpuType = "NVIDIA L40"
	GpuNvL40S  GpuType = "NVIDIA L40S"
	GpuNvL4    GpuType = "NVIDIA L4"
	GpuNvV100  GpuType = "NVIDIA V100"
	GpuNvT4    GpuType = "NVIDIA T4"
	GpuNvT4G   GpuType = "NVIDIA T4G"

	// NVIDIA RTX 消费级GPU
	GpuNvRTX4090 GpuType = "NVIDIA RTX 4090"
	GpuNvRTX4080 GpuType = "NVIDIA RTX 4080"
	GpuNvRTX4070 GpuType = "NVIDIA RTX 4070"
	GpuNvRTX4060 GpuType = "NVIDIA RTX 4060"
	GpuNvRTX3090 GpuType = "NVIDIA RTX 3090"
	GpuNvRTX3080 GpuType = "NVIDIA RTX 3080"
	GpuNvRTX3070 GpuType = "NVIDIA RTX 3070"
	GpuNvRTX3060 GpuType = "NVIDIA RTX 3060"
	GpuNvRTX2080 GpuType = "NVIDIA RTX 2080"
	GpuNvRTX2070 GpuType = "NVIDIA RTX 2070"
	GpuNvRTX2060 GpuType = "NVIDIA RTX 2060"

	// NVIDIA GTX 系列
	GpuNvGTX1080 GpuType = "NVIDIA GTX 1080"
	GpuNvGTX1070 GpuType = "NVIDIA GTX 1070"
	GpuNvGTX1060 GpuType = "NVIDIA GTX 1060"

	// NVIDIA Titan 系列
	GpuNvTitanRTX GpuType = "NVIDIA Titan RTX"
	GpuNvTitanV   GpuType = "NVIDIA Titan V"
	GpuNvTitanXp  GpuType = "NVIDIA Titan Xp"

	// NVIDIA Quadro 系列
	GpuNvQuadroRTX8000 GpuType = "NVIDIA Quadro RTX 8000"
	GpuNvQuadroRTX6000 GpuType = "NVIDIA Quadro RTX 6000"
	GpuNvQuadroRTX5000 GpuType = "NVIDIA Quadro RTX 5000"
	GpuNvQuadroP6000   GpuType = "NVIDIA Quadro P6000"
	GpuNvQuadroP5000   GpuType = "NVIDIA Quadro P5000"

	// NVIDIA Tesla 系列
	GpuNvTeslaV100 GpuType = "NVIDIA Tesla V100"
	GpuNvTeslaP100 GpuType = "NVIDIA Tesla P100"
	GpuNvTeslaP40  GpuType = "NVIDIA Tesla P40"
	GpuNvTeslaP4   GpuType = "NVIDIA Tesla P4"
	GpuNvTeslaK80  GpuType = "NVIDIA Tesla K80"

	// AMD Radeon RX 系列
	GpuAmdRX7900XTX GpuType = "AMD Radeon RX 7900 XTX"
	GpuAmdRX7900XT  GpuType = "AMD Radeon RX 7900 XT"
	GpuAmdRX7800XT  GpuType = "AMD Radeon RX 7800 XT"
	GpuAmdRX7700XT  GpuType = "AMD Radeon RX 7700 XT"
	GpuAmdRX6950XT  GpuType = "AMD Radeon RX 6950 XT"
	GpuAmdRX6900XT  GpuType = "AMD Radeon RX 6900 XT"
	GpuAmdRX6800XT  GpuType = "AMD Radeon RX 6800 XT"
	GpuAmdRX6700XT  GpuType = "AMD Radeon RX 6700 XT"
	GpuAmdRX6600XT  GpuType = "AMD Radeon RX 6600 XT"
	GpuAmdRX6500XT  GpuType = "AMD Radeon RX 6500 XT"

	// AMD Radeon Pro 系列
	GpuAmdRadeonProW6800 GpuType = "AMD Radeon Pro W6800"
	GpuAmdRadeonProW6600 GpuType = "AMD Radeon Pro W6600"
	GpuAmdRadeonProW5700 GpuType = "AMD Radeon Pro W5700"
	GpuAmdRadeonProW5500 GpuType = "AMD Radeon Pro W5500"

	// AMD Instinct 系列
	GpuAmdInstinctMI250X GpuType = "AMD Instinct MI250X"
	GpuAmdInstinctMI210  GpuType = "AMD Instinct MI210"
	GpuAmdInstinctMI100  GpuType = "AMD Instinct MI100"
	GpuAmdInstinctMI50   GpuType = "AMD Instinct MI50"
	GpuAmdInstinctMI25   GpuType = "AMD Instinct MI25"

	// Intel Arc 系列
	GpuIntelArcA770 GpuType = "Intel Arc A770"
	GpuIntelArcA750 GpuType = "Intel Arc A750"
	GpuIntelArcA580 GpuType = "Intel Arc A580"
	GpuIntelArcA380 GpuType = "Intel Arc A380"
	GpuIntelArcA310 GpuType = "Intel Arc A310"

	// Intel Xe 系列
	GpuIntelXeHP  GpuType = "Intel Xe-HP"
	GpuIntelXeHPG GpuType = "Intel Xe-HPG"
	GpuIntelXeHPC GpuType = "Intel Xe-HPC"

	// Intel Data Center GPU
	GpuIntelDataCenterGPUMax  GpuType = "Intel Data Center GPU Max"
	GpuIntelDataCenterGPUFlex GpuType = "Intel Data Center GPU Flex"
)
