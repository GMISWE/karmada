package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindJuicefs            = "Juicefs"
	ResourcePluralJuicefs          = "juicefs"
	ResourceSingularJuicefs        = "juicefs"
	ResourceNamespaceScopedJuicefs = false
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Juicefs represents a JuiceFS storage resource
type Juicefs struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JuicefsSpec   `json:"spec,omitempty"`
	Status JuicefsStatus `json:"status,omitempty"`
}

// JuicefsSpec defines the specification for a JuiceFS storage resource
type JuicefsSpec struct {
	Description string        `json:"description,omitempty"`
	Location    Location      `json:"location,omitempty"`
	Provider    Provider      `json:"provider,omitempty"`
	Public      string        `json:"public,omitempty"`
	Client      JuicefsClient `json:"client,omitempty"`
}

// Location represents the storage location
type Location struct {
	Region string `json:"region,omitempty"`
	AZ     string `json:"az,omitempty"`
}

// Provider represents the storage provider information
type Provider struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// JuicefsClient defines the configuration for JuiceFS client
type JuicefsClient struct {
	Image        string                      `json:"image,omitempty"`
	CacheDir     string                      `json:"cache-dir,omitempty"`
	MountOptions []string                    `json:"mount-options,omitempty"`
	Envs         []corev1.EnvVar             `json:"envs,omitempty"`
	Volumes      []corev1.Volume             `json:"volumes,omitempty"`
	VolumeMounts []corev1.VolumeMount        `json:"volumeMounts,omitempty"`
	Resources    corev1.ResourceRequirements `json:"resources,omitempty"`
	EE           *EnterpriseEdition          `json:"ee,omitempty"`
	CE           *CommunityEdition           `json:"ce,omitempty"`
}

// EnterpriseEdition defines configuration for JuiceFS enterprise edition
type EnterpriseEdition struct {
	Version string         `json:"version,omitempty"`
	Auth    EnterpriseAuth `json:"auth,omitempty"`
}

// EnterpriseAuth defines authentication details for JuiceFS enterprise edition
type EnterpriseAuth struct {
	Token      string `json:"token,omitempty"`
	Name       string `json:"name,omitempty"`
	BaseURL    string `json:"base-url,omitempty"`
	AccessKey  string `json:"access-key,omitempty"`
	SecretKey  string `json:"secret-key,omitempty"`
	AccessKey2 string `json:"access-key2,omitempty"`
	SecretKey2 string `json:"secret-key2,omitempty"`
}

// CommunityEdition defines configuration for JuiceFS community edition
type CommunityEdition struct {
	Version string        `json:"version,omitempty"`
	MetaURL string        `json:"meta-url,omitempty"`
	Backend BackendConfig `json:"backend,omitempty"`
}

// BackendConfig defines backend storage configuration for JuiceFS community edition
type BackendConfig struct {
	Storage   string `json:"storage,omitempty"`
	AccessKey string `json:"access-key,omitempty"`
	SecretKey string `json:"secret-key,omitempty"`
}

// JuicefsStatus defines the observed state of JuiceFS resource
type JuicefsStatus struct {
	Phase      string             `json:"phase,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	MountInfo  MountInfo          `json:"mountInfo,omitempty"`
}

// MountInfo provides information about mounted JuiceFS
type MountInfo struct {
	MountedNodes []string    `json:"mountedNodes,omitempty"`
	LastMounted  metav1.Time `json:"lastMounted,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// JuicefsList contains a list of Juicefs resources
type JuicefsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Juicefs `json:"items"`
}
