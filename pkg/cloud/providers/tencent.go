/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : tencent云提供商
 */

package providers

import (
	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	"github.com/karmada-io/karmada/pkg/cloud/constants"
)

type TencentProviderConfig struct {
	SecretID  string
	SecretKey string
}

type TencentProvider struct {
	Provider
	name      string
	locations []v1alpha1.CloudLocation
	config    *TencentProviderConfig
}

func NewTencentProvider(config *TencentProviderConfig) *TencentProvider {
	return &TencentProvider{
		name:      string(constants.CloudProviderTencent),
		locations: []v1alpha1.CloudLocation{},
		config:    config,
	}
}

func (p *TencentProvider) GetName() string {
	return p.name
}

func (p *TencentProvider) GetLocations() ([]v1alpha1.CloudLocation, error) {
	return p.locations, nil
}

func (p *TencentProvider) GetInstanceList() ([]string, error) {
	return []string{}, nil
}

func (p *TencentProvider) GetInstancePrice(instanceType string) (float64, error) {
	return 0, nil
}

var _ Provider = &TencentProvider{}
