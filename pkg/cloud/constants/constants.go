/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/11
 * @Desc    : 描述信息
 */

package constants

import (
	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
)

// cloud provider type
const (
	CloudProviderAws      v1alpha1.CloudProvider = "aws"
	CloudProviderGcp      v1alpha1.CloudProvider = "gcp"
	CloudProviderAzure    v1alpha1.CloudProvider = "azure"
	CloudProviderHuawei   v1alpha1.CloudProvider = "huawei"
	CloudProviderAliyun   v1alpha1.CloudProvider = "aliyun"
	CloudProviderTencent  v1alpha1.CloudProvider = "tencent"
	CloudProviderBaidu    v1alpha1.CloudProvider = "baidu"
	CloudProviderKingsoft v1alpha1.CloudProvider = "kingsoft"
	CloudProviderGmi      v1alpha1.CloudProvider = "gmi"
	CloudProviderOthers   v1alpha1.CloudProvider = "others"
)

// aws region
const (
	AwsRegionUsEast1  v1alpha1.CloudRegion = "us-east-1"
	AwsRegionUsEast2  v1alpha1.CloudRegion = "us-east-2"
	AwsRegionUsEast3  v1alpha1.CloudRegion = "us-east-3"
	AwsRegionUsEast4  v1alpha1.CloudRegion = "us-east-4"
	AwsRegionUsEast5  v1alpha1.CloudRegion = "us-east-5"
	AwsRegionUsEast6  v1alpha1.CloudRegion = "us-east-6"
	AwsRegionUsEast7  v1alpha1.CloudRegion = "us-east-7"
	AwsRegionUsEast8  v1alpha1.CloudRegion = "us-east-8"
	AwsRegionUsEast9  v1alpha1.CloudRegion = "us-east-9"
	AwsRegionUsEast10 v1alpha1.CloudRegion = "us-east-10"
)

// gcp region
const (
	GcpRegionUsEast1  v1alpha1.CloudRegion = "us-east-1"
	GcpRegionUsEast2  v1alpha1.CloudRegion = "us-east-2"
	GcpRegionUsEast3  v1alpha1.CloudRegion = "us-east-3"
	GcpRegionUsEast4  v1alpha1.CloudRegion = "us-east-4"
	GcpRegionUsEast5  v1alpha1.CloudRegion = "us-east-5"
	GcpRegionUsEast6  v1alpha1.CloudRegion = "us-east-6"
	GcpRegionUsEast7  v1alpha1.CloudRegion = "us-east-7"
	GcpRegionUsEast8  v1alpha1.CloudRegion = "us-east-8"
	GcpRegionUsEast9  v1alpha1.CloudRegion = "us-east-9"
	GcpRegionUsEast10 v1alpha1.CloudRegion = "us-east-10"
)

// azure region
const (
	AzureRegionUsEast1  v1alpha1.CloudRegion = "us-east-1"
	AzureRegionUsEast2  v1alpha1.CloudRegion = "us-east-2"
	AzureRegionUsEast3  v1alpha1.CloudRegion = "us-east-3"
	AzureRegionUsEast4  v1alpha1.CloudRegion = "us-east-4"
	AzureRegionUsEast5  v1alpha1.CloudRegion = "us-east-5"
	AzureRegionUsEast6  v1alpha1.CloudRegion = "us-east-6"
	AzureRegionUsEast7  v1alpha1.CloudRegion = "us-east-7"
	AzureRegionUsEast8  v1alpha1.CloudRegion = "us-east-8"
	AzureRegionUsEast9  v1alpha1.CloudRegion = "us-east-9"
	AzureRegionUsEast10 v1alpha1.CloudRegion = "us-east-10"
)

// huawei region
const (
	HuaweiRegionUsEast1 v1alpha1.CloudRegion = "us-east-1"
	HuaweiRegionUsEast2 v1alpha1.CloudRegion = "us-east-2"
	HuaweiRegionUsEast3 v1alpha1.CloudRegion = "us-east-3"
	HuaweiRegionUsEast4 v1alpha1.CloudRegion = "us-east-4"
	HuaweiRegionUsEast5 v1alpha1.CloudRegion = "us-east-5"
)

// aliyun region
const (
	AliyunRegionUsEast1 v1alpha1.CloudRegion = "us-east-1"
	AliyunRegionUsEast2 v1alpha1.CloudRegion = "us-east-2"
	AliyunRegionUsEast3 v1alpha1.CloudRegion = "us-east-3"
	AliyunRegionUsEast4 v1alpha1.CloudRegion = "us-east-4"
	AliyunRegionUsEast5 v1alpha1.CloudRegion = "us-east-5"
)

// tencent region
const (
	TencentRegionUsEast1 v1alpha1.CloudRegion = "us-east-1"
	TencentRegionUsEast2 v1alpha1.CloudRegion = "us-east-2"
	TencentRegionUsEast3 v1alpha1.CloudRegion = "us-east-3"
	TencentRegionUsEast4 v1alpha1.CloudRegion = "us-east-4"
	TencentRegionUsEast5 v1alpha1.CloudRegion = "us-east-5"
)

// baidu region
const (
	BaiduRegionUsEast1 v1alpha1.CloudRegion = "us-east-1"
	BaiduRegionUsEast2 v1alpha1.CloudRegion = "us-east-2"
	BaiduRegionUsEast3 v1alpha1.CloudRegion = "us-east-3"
	BaiduRegionUsEast4 v1alpha1.CloudRegion = "us-east-4"
	BaiduRegionUsEast5 v1alpha1.CloudRegion = "us-east-5"
)

// kingsoft region
const (
	KingsoftRegionUsEast1 v1alpha1.CloudRegion = "us-east-1"
	KingsoftRegionUsEast2 v1alpha1.CloudRegion = "us-east-2"
	KingsoftRegionUsEast3 v1alpha1.CloudRegion = "us-east-3"
	KingsoftRegionUsEast4 v1alpha1.CloudRegion = "us-east-4"
	KingsoftRegionUsEast5 v1alpha1.CloudRegion = "us-east-5"
)

// gmi region
const (
	GmiRegionUsEast1 v1alpha1.CloudRegion = "us-east-1"
	GmiRegionUsEast2 v1alpha1.CloudRegion = "us-east-2"
	GmiRegionUsEast3 v1alpha1.CloudRegion = "us-east-3"
	GmiRegionUsEast4 v1alpha1.CloudRegion = "us-east-4"
	GmiRegionUsEast5 v1alpha1.CloudRegion = "us-east-5"
)

// others region
const (
	OthersRegionUsEast1 v1alpha1.CloudRegion = "us-east-1"
	OthersRegionUsEast2 v1alpha1.CloudRegion = "us-east-2"
	OthersRegionUsEast3 v1alpha1.CloudRegion = "us-east-3"
	OthersRegionUsEast4 v1alpha1.CloudRegion = "us-east-4"
	OthersRegionUsEast5 v1alpha1.CloudRegion = "us-east-5"
)

var (
	// aws locations
	AwsLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		AwsRegionUsEast1: {
			Region: AwsRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
	// gcp locations
	GcpLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		GcpRegionUsEast1: {
			Region: GcpRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
	// azure locations
	AzureLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		AzureRegionUsEast1: {
			Region: AzureRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
	// huawei locations
	HuaweiLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		HuaweiRegionUsEast1: {
			Region: HuaweiRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
	// aliyun locations
	AliyunLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		AliyunRegionUsEast1: {
			Region: AliyunRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
	// tencent locations
	TencentLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		TencentRegionUsEast1: {
			Region: TencentRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
	// baidu locations
	BaiduLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		BaiduRegionUsEast1: {
			Region: BaiduRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
	// kingsoft locations
	KingsoftLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		KingsoftRegionUsEast1: {
			Region: KingsoftRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
	// gmi locations
	GmiLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		GmiRegionUsEast1: {
			Region: GmiRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
	// others locations
	OthersLocations = map[v1alpha1.CloudRegion]v1alpha1.CloudLocation{
		OthersRegionUsEast1: {
			Region: OthersRegionUsEast1,
			Zones:  []string{"us-east-1a", "us-east-1b", "us-east-1c"},
		},
	}
)
