/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : 描述信息
 */

package cloud

import (
	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	"github.com/karmada-io/karmada/pkg/cloud/constants"
	"github.com/karmada-io/karmada/pkg/cloud/providers"
)

type CloudConfig struct {
	Aws      *providers.AwsProviderConfig
	Gcp      *providers.GcpProviderConfig
	Azure    *providers.AzureProviderConfig
	Huawei   *providers.HuaweiProviderConfig
	Aliyun   *providers.AliyunProviderConfig
	Tencent  *providers.TencentProviderConfig
	Baidu    *providers.BaiduProviderConfig
	Kingsoft *providers.KingsoftProviderConfig
	Gmi      *providers.GmiProviderConfig
}

type Cloud struct {
	providers map[v1alpha1.CloudProvider]providers.Provider
}

var cloud *Cloud

func NewCloud(config *CloudConfig) *Cloud {
	if cloud == nil {
		cloud = &Cloud{
			providers: map[v1alpha1.CloudProvider]providers.Provider{
				constants.CloudProviderAws:      providers.NewAwsProvider(config.Aws),
				constants.CloudProviderGcp:      providers.NewGcpProvider(config.Gcp),
				constants.CloudProviderAzure:    providers.NewAzureProvider(config.Azure),
				constants.CloudProviderHuawei:   providers.NewHuaweiProvider(config.Huawei),
				constants.CloudProviderAliyun:   providers.NewAliyunProvider(config.Aliyun),
				constants.CloudProviderTencent:  providers.NewTencentProvider(config.Tencent),
				constants.CloudProviderBaidu:    providers.NewBaiduProvider(config.Baidu),
				constants.CloudProviderKingsoft: providers.NewKingsoftProvider(config.Kingsoft),
				constants.CloudProviderGmi:      providers.NewGmiProvider(config.Gmi),
			},
		}
	}
	return cloud
}

func (c *Cloud) GetProvider(cloudProvider v1alpha1.CloudProvider) providers.Provider {
	return c.providers[cloudProvider]
}
