/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : 描述信息
 */

package providers

import (
	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	"github.com/karmada-io/karmada/pkg/cloud/constants"
)

type GmiProviderConfig struct {
	AccessKey string
	SecretKey string
}

type GmiProvider struct {
	Provider
	name      string
	locations []v1alpha1.CloudLocation
	config    *GmiProviderConfig
}

func NewGmiProvider(config *GmiProviderConfig) *GmiProvider {
	return &GmiProvider{
		name:      string(constants.CloudProviderGmi),
		locations: []v1alpha1.CloudLocation{},
		config:    config,
	}
}

func (p *GmiProvider) GetName() string {
	return p.name
}

func (p *GmiProvider) GetLocations() ([]v1alpha1.CloudLocation, error) {
	return p.locations, nil
}

func (p *GmiProvider) GetInstanceList() ([]string, error) {
	return []string{}, nil
}

func (p *GmiProvider) GetInstancePrice(instanceType string) (float64, error) {
	return 0, nil
}

var _ Provider = &GmiProvider{}
