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
