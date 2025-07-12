/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : baidu云提供商
 */

package providers

import (
	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	"github.com/karmada-io/karmada/pkg/cloud/constants"
)

type BaiduProviderConfig struct {
	AccessKey string
	SecretKey string
}

type BaiduProvider struct {
	Provider
	name      string
	locations []v1alpha1.CloudLocation
	config    *BaiduProviderConfig
}

func NewBaiduProvider(config *BaiduProviderConfig) *BaiduProvider {
	return &BaiduProvider{
		name:      string(constants.CloudProviderBaidu),
		locations: []v1alpha1.CloudLocation{},
		config:    config,
	}
}

func (p *BaiduProvider) GetName() string {
	return p.name
}

func (p *BaiduProvider) GetLocations() ([]v1alpha1.CloudLocation, error) {
	return p.locations, nil
}

func (p *BaiduProvider) GetInstanceList() ([]string, error) {
	return []string{}, nil
}

func (p *BaiduProvider) GetInstancePrice(instanceType string) (float64, error) {
	return 0, nil
}

var _ Provider = &BaiduProvider{}
