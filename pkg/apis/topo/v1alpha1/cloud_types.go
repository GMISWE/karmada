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
	CloudKindCloud            = "Cloud"
	CloudPluralCloud          = "clouds"
	CloudSingularCloud        = "Cloud"
	CloudNamespaceScopedCloud = false
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:storageversion

// Cloud represents the desired state and status of a member cluster.
type Cloud struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   *CloudSpec   `json:"spec,omitempty"`
	Status *CloudStatus `json:"status,omitempty"`
}

type CloudSpec struct {
	// +required
	CloudProvider CloudProvider `json:"cloudProvider"`
	// +required
	Location *CloudLocation `json:"location,omitempty"`
	// +optional
	CloudAws *CloudAws `json:"cloudAws,omitempty"`
	// +optional
	CloudGcp *CloudGcp `json:"cloudGcp,omitempty"`
	// +optional
	CloudAzure *CloudAzure `json:"cloudAzure,omitempty"`
	// +optional
	CloudHuawei *CloudHuawei `json:"cloudHuawei,omitempty"`
	// +optional
	CloudAliyun *CloudAliyun `json:"cloudAliyun,omitempty"`
	// +optional
	CloudTencent *CloudTencent `json:"cloudTencent,omitempty"`
}

type CloudAws struct {
	// +required
	GpuPrices []GpuPrice `json:"gpuPrices"`
}

type CloudGcp struct {
	// +required
	GpuPrices []GpuPrice `json:"gpuPrices"`
}

type CloudAzure struct {
	// +required
	GpuPrices []GpuPrice `json:"gpuPrices"`
}

type CloudHuawei struct {
	// +required
	GpuPrices []GpuPrice `json:"gpuPrices"`
}

type CloudAliyun struct {
	GpuPrices []GpuPrice `json:"gpuPrices"`
}

type CloudTencent struct {
	GpuPrices []GpuPrice `json:"gpuPrices"`
}

type GpuPrice struct {
	InstanceType string `json:"instanceType"`
	Price        string `json:"price"`
	PriceUnit    string `json:"priceUnit"`
	Stock        bool   `json:"stock"`
}

type CloudLocation struct {
	// +required
	Region CloudRegion `json:"region"`
	// +required
	Zones []string `json:"zones"`
}

type CloudStatus struct {
	Phase      string             `json:"phase,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// CloudList contains a list of Cloud
type CloudList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items holds a list of Cloud.
	Items []Cloud `json:"items"`
}

type CloudProvider string
type CloudRegion string

// const (
// 	CloudProviderAws      CloudProvider = "aws"
// 	CloudProviderGcp      CloudProvider = "gcp"
// 	CloudProviderAzure    CloudProvider = "azure"
// 	CloudProviderHuawei   CloudProvider = "huawei"
// 	CloudProviderAliyun   CloudProvider = "aliyun"
// 	CloudProviderTencent  CloudProvider = "tencent"
// 	CloudProviderBaidu    CloudProvider = "baidu"
// 	CloudProviderKingsoft CloudProvider = "kingsoft"
// 	CloudProviderGmi      CloudProvider = "gmi"
// 	CloudProviderOthers   CloudProvider = "others"
// )

// const (
// 	AwsRegionUsEast1  CloudRegion = "us-east-1"
// 	AwsRegionUsEast2  CloudRegion = "us-east-2"
// 	AwsRegionUsEast3  CloudRegion = "us-east-3"
// 	AwsRegionUsEast4  CloudRegion = "us-east-4"
// 	AwsRegionUsEast5  CloudRegion = "us-east-5"
// 	AwsRegionUsEast6  CloudRegion = "us-east-6"
// 	AwsRegionUsEast7  CloudRegion = "us-east-7"
// 	AwsRegionUsEast8  CloudRegion = "us-east-8"
// 	AwsRegionUsEast9  CloudRegion = "us-east-9"
// 	AwsRegionUsEast10 CloudRegion = "us-east-10"
// )

// var (
// 	AwsLocations = AwsLocation{
// 		Region: "us-east-1",
// 		Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
// 	}
// 	GcpLocations = GcpLocation{
// 		Region: "us-east-1",
// 		Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
// 	}
// )

// type AwsLocation struct {
// 	Region string   `json:"region"`
// 	Zones  []string `json:"zones"`
// }

// type GcpLocation struct {
// 	Region string   `json:"region"`
// 	Zones  []string `json:"zones"`
// }

// type AzureLocation struct {
// }
