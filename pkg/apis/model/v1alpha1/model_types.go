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
	ResourceKindModel            = "Model"
	ResourcePluralModel          = "models"
	ResourceSingularModel        = "model"
	ResourceNamespaceScopedModel = false
)

type ModelType string

const (
	ModelTypeLLM   ModelType = "LLM"
	ModelTypeVideo ModelType = "Video"
	ModelTypeAudio ModelType = "Audio"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:path=models,scope=Namespaced,shortName=m,categories={karmada-io}
// +kubebuilder:modelversion

type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// Spec represents the desired behavior of Model.
	Spec ModelSpec `json:"spec"`
}

type LLMConfig struct {
	// +required
	ParameterSize int64 `json:"parameterSize"` // 模型参数大小，单位：GB
	// +required
	RunArgs map[string]string `json:"runArgs,omitempty"` // 模型运行参数，json格式，例如：{"temperature": 0.5, "max_tokens": 100}
}

type VideoConfig struct {
	// +required
	ParameterSize int64 `json:"parameterSize"` // 模型参数大小，单位：GB
	// +required
	RunArgs map[string]string `json:"runArgs,omitempty"` // 模型运行参数
}

type AudioConfig struct {
	// +required
	ParameterSize int64 `json:"parameterSize"` // 模型参数大小，单位：GB
	// +required
	RunArgs map[string]string `json:"runArgs,omitempty"` // 模型运行参数
}

type ModelSpec struct {
	// +required
	ModelType ModelType `json:"modelType"` // 模型类型：LLM、Video、Audio
	// +required
	ModelName string `json:"modelName"`
	// +required
	ModelVersion string `json:"modelVersion"`
	// +required
	ModelPath string `json:"modelPath"`
	// +required
	ModelImage string `json:"modelImage"` // 模型镜像
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// +optional
	MinReplicas int `json:"minReplicas,omitempty"` // 最小副本数，如果为0，则默认模型可以自动下线
	// +optional
	MaxReplicas int `json:"maxReplicas,omitempty"` // 最大副本数，如果为0，则默认math.MaxInt

	// +optional
	LLMConfig *LLMConfig `json:"llmConfig,omitempty"`
	// +optional
	VideoConfig *VideoConfig `json:"videoConfig,omitempty"`
	// +optional
	AudioConfig *AudioConfig `json:"audioConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
type ModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Model `json:"items"`
}
