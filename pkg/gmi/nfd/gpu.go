/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/29
 * @Desc    : GPU detection implementation using NVML only
 */

package nfd

import (
	"time"

	"github.com/karmada-io/karmada/pkg/gmi/base/types"
	"github.com/karmada-io/karmada/pkg/gmi/nfd/nvml"
	"github.com/sirupsen/logrus"
)

// GPUDetector provides GPU detection using NVML only
type GPUDetector struct {
	nvmlDetector *nvml.NVMLDetector
	timeout      time.Duration
}

// NewGPUDetector creates a new GPU detector using NVML
func NewGPUDetector() *GPUDetector {
	return &GPUDetector{
		nvmlDetector: nvml.NewNVMLDetector(),
		timeout:      30 * time.Second,
	}
}

// DiscoverGPUs discovers GPUs using NVML only
func (gd *GPUDetector) DiscoverGPUs() []types.Gpu {
	// 使用NVML检测GPU
	gpus := gd.nvmlDetector.DiscoverGPUs()

	if len(gpus) == 0 {
		logrus.Warn("No GPUs detected!")
		logrus.Info("This could be due to:")
		logrus.Info("   - No NVIDIA GPU hardware present")
		logrus.Info("   - NVIDIA driver not installed or not compatible")
		logrus.Info("   - Missing device mounts in container (/dev/nvidia*)")
		logrus.Info("   - Insufficient permissions to access GPU devices")
		logrus.Info("   - NVML library not available")
	}

	return gpus
}
