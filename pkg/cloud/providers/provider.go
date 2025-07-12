/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : 云提供商
 */

package providers

import (
	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
)

type Provider interface {
	GetName() string
	GetLocations() ([]v1alpha1.CloudLocation, error)
	GetInstanceList() ([]string, error)
	GetInstancePrice(instanceType string) (float64, error)
}
