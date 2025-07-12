/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : aws云提供商
 */

package providers

import (
	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	"github.com/karmada-io/karmada/pkg/cloud/constants"
)

type AwsProvider struct {
	Provider
	name      string
	locations []v1alpha1.CloudLocation
	config    *AwsProviderConfig
}

type AwsProviderConfig struct {
	AccessKey string
	SecretKey string
}

func NewAwsProvider(config *AwsProviderConfig) *AwsProvider {
	return &AwsProvider{
		name:      string(constants.CloudProviderAws),
		locations: []v1alpha1.CloudLocation{},
		config:    config,
	}
}

func (p *AwsProvider) GetName() string {
	return p.name
}

func (p *AwsProvider) GetLocations() ([]v1alpha1.CloudLocation, error) {
	return p.locations, nil
}

func (p *AwsProvider) GetInstanceList() ([]string, error) {
	return []string{}, nil
}

func (p *AwsProvider) GetInstancePrice(instanceType string) (float64, error) {
	return 0, nil
}

var _ Provider = &AwsProvider{}
