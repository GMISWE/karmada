/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : huawei云提供商
 */

package providers

import (
	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	"github.com/karmada-io/karmada/pkg/cloud/constants"
)

type HuaweiProviderConfig struct {
	AccessKey string
	SecretKey string
}

type HuaweiProvider struct {
	Provider
	name      string
	locations []v1alpha1.CloudLocation
	config    *HuaweiProviderConfig
}

func NewHuaweiProvider(config *HuaweiProviderConfig) *HuaweiProvider {
	return &HuaweiProvider{
		name:      string(constants.CloudProviderHuawei),
		locations: []v1alpha1.CloudLocation{},
		config:    config,
	}
}

func (p *HuaweiProvider) GetName() string {
	return p.name
}

func (p *HuaweiProvider) GetLocations() ([]v1alpha1.CloudLocation, error) {
	return p.locations, nil
}

func (p *HuaweiProvider) GetInstanceList() ([]string, error) {
	return []string{}, nil
}

func (p *HuaweiProvider) GetInstancePrice(instanceType string) (float64, error) {
	return 0, nil
}

var _ Provider = &HuaweiProvider{}
