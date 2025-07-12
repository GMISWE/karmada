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

type KingsoftProviderConfig struct {
	AccessKey string
	SecretKey string
}

type KingsoftProvider struct {
	Provider
	name      string
	locations []v1alpha1.CloudLocation
	config    *KingsoftProviderConfig
}

func NewKingsoftProvider(config *KingsoftProviderConfig) *KingsoftProvider {
	return &KingsoftProvider{
		name:      string(constants.CloudProviderKingsoft),
		locations: []v1alpha1.CloudLocation{},
		config:    config,
	}
}

func (p *KingsoftProvider) GetName() string {
	return p.name
}

func (p *KingsoftProvider) GetLocations() ([]v1alpha1.CloudLocation, error) {
	return p.locations, nil
}

func (p *KingsoftProvider) GetInstanceList() ([]string, error) {
	return []string{}, nil
}

func (p *KingsoftProvider) GetInstancePrice(instanceType string) (float64, error) {
	return 0, nil
}

var _ Provider = &KingsoftProvider{}
