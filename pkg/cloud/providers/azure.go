/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : azure云提供商
 */

package providers

import (
	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	"github.com/karmada-io/karmada/pkg/cloud/constants"
)

type AzureProviderConfig struct {
	SubscriptionID string
}

type AzureProvider struct {
	Provider
	name      string
	locations []v1alpha1.CloudLocation
	config    *AzureProviderConfig
}

func NewAzureProvider(config *AzureProviderConfig) *AzureProvider {
	return &AzureProvider{
		name:      string(constants.CloudProviderAzure),
		locations: []v1alpha1.CloudLocation{},
		config:    config,
	}
}

func (p *AzureProvider) GetName() string {
	return p.name
}

func (p *AzureProvider) GetLocations() ([]v1alpha1.CloudLocation, error) {
	return p.locations, nil
}

func (p *AzureProvider) GetInstanceList() ([]string, error) {
	return []string{}, nil
}

func (p *AzureProvider) GetInstancePrice(instanceType string) (float64, error) {
	return 0, nil
}

var _ Provider = &AzureProvider{}
