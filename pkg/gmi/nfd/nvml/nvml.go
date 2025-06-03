/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/06/01
 * @Desc    : 内置NVML GPU检测器
 */

package nvml

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/karmada-io/karmada/pkg/gmi/base/types"
	"github.com/sirupsen/logrus"
)

// NVMLDetector nvml detector
type NVMLDetector struct {
	initialized     bool
	processDetector *NVMLProcessDetector
}

// NewNVMLDetector create nvml detector
func NewNVMLDetector() *NVMLDetector {
	return &NVMLDetector{
		initialized:     false,
		processDetector: NewNVMLProcessDetector(),
	}
}

// DiscoverGPUs discover gpus
func (nd *NVMLDetector) DiscoverGPUs() []types.Gpu {
	var gpus []types.Gpu

	logrus.Debug("Starting built-in NVML GPU detection...")

	// init nvml
	if err := nd.initNVML(); err != nil {
		logrus.Errorf("NVML initialization failed: %v", err)
		return gpus
	}
	defer nd.shutdownNVML()

	// get gpu count
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		logrus.Errorf("Failed to get GPU count: %v", nvml.ErrorString(ret))
		return gpus
	}

	if count == 0 {
		logrus.Info("No GPUs found via NVML")
		return gpus
	}

	logrus.Debugf("Found %d GPU(s) via NVML", count)

	// get gpu info
	for i := range count {
		if gpu := nd.getGPUInfo(i); gpu != nil {
			gpus = append(gpus, *gpu)
		}
	}

	logrus.Debugf("Successfully detected %d GPU(s) via built-in NVML", len(gpus))
	return gpus
}

// initNVML
func (nd *NVMLDetector) initNVML() error {
	if nd.initialized {
		return nil
	}

	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		return fmt.Errorf("NVML init failed: %v", nvml.ErrorString(ret))
	}

	nd.initialized = true
	logrus.Debug("NVML initialized successfully")
	return nil
}

// shutdownNVML shutdown nvml
func (nd *NVMLDetector) shutdownNVML() {
	if !nd.initialized {
		return
	}

	ret := nvml.Shutdown()
	if ret != nvml.SUCCESS {
		logrus.Warnf("NVML shutdown warning: %v", nvml.ErrorString(ret))
	} else {
		logrus.Debug("NVML shutdown successfully")
	}

	nd.initialized = false
}

// getGPUInfo get gpu info
func (nd *NVMLDetector) getGPUInfo(index int) *types.Gpu {
	// get gpu handle
	device, ret := nvml.DeviceGetHandleByIndex(index)
	if ret != nvml.SUCCESS {
		logrus.Errorf("Failed to get GPU %d handle: %v", index, nvml.ErrorString(ret))
		return nil
	}

	gpu := &types.Gpu{
		ID:      fmt.Sprintf("%d", index),
		Cluster: "local",
	}

	// get uuid
	if uuid, ret := device.GetUUID(); ret == nvml.SUCCESS {
		gpu.UUID = uuid
		logrus.Debugf("GPU %d UUID: %s", index, uuid)
	} else {
		logrus.Debugf("Failed to get GPU %d UUID: %v", index, nvml.ErrorString(ret))
		gpu.UUID = fmt.Sprintf("nvml-gpu-%d", index)
	}

	// get device name
	if name, ret := device.GetName(); ret == nvml.SUCCESS {
		gpu.Model = name
		logrus.Debugf("GPU %d Name: %s", index, name)
	} else {
		logrus.Debugf("Failed to get GPU %d name: %v", index, nvml.ErrorString(ret))
		gpu.Model = "Unknown NVIDIA GPU"
	}

	// get memory info
	if memInfo, ret := device.GetMemoryInfo(); ret == nvml.SUCCESS {
		gpu.FreeMem = fmt.Sprintf("%.0f MB", float64(memInfo.Free)/(1024*1024))
		gpu.UseMem = fmt.Sprintf("%.0f MB", float64(memInfo.Used)/(1024*1024))
		logrus.Debugf("GPU %d Memory: %.1f GB total, %.1f GB used, %.1f GB free",
			index,
			float64(memInfo.Total)/(1024*1024*1024),
			float64(memInfo.Used)/(1024*1024*1024),
			float64(memInfo.Free)/(1024*1024*1024))
	} else {
		logrus.Debugf("Failed to get GPU %d memory info: %v", index, nvml.ErrorString(ret))
	}

	// get temperature
	if temp, ret := device.GetTemperature(nvml.TEMPERATURE_GPU); ret == nvml.SUCCESS {
		gpu.Temp = int64(temp)
		logrus.Debugf("GPU %d Temperature: %d°C", index, temp)
	} else {
		logrus.Debugf("Failed to get GPU %d temperature: %v", index, nvml.ErrorString(ret))
	}

	// get power
	if power, ret := device.GetPowerUsage(); ret == nvml.SUCCESS {
		gpu.Power = int64(power / 1000) // convert to W (NVML returns mW)
		logrus.Debugf("GPU %d Power: %d W", index, gpu.Power)
	} else {
		logrus.Debugf("Failed to get GPU %d power: %v", index, nvml.ErrorString(ret))
	}

	// get utilization
	if util, ret := device.GetUtilizationRates(); ret == nvml.SUCCESS {
		gpu.Utilization = fmt.Sprintf("%d%%", util.Gpu)
		gpu.Load = int64(util.Gpu) // GPU load percentage
		logrus.Debugf("GPU %d Utilization: GPU=%d%%, Memory=%d%%", index, util.Gpu, util.Memory)
	} else {
		logrus.Debugf("Failed to get GPU %d utilization: %v", index, nvml.ErrorString(ret))
	}

	// get fan speed
	if fanSpeed, ret := device.GetFanSpeed(); ret == nvml.SUCCESS {
		gpu.Fan = int64(fanSpeed) // fan speed percentage
		logrus.Debugf("GPU %d Fan Speed: %d%%", index, fanSpeed)
	} else {
		logrus.Debugf("Failed to get GPU %d fan speed: %v", index, nvml.ErrorString(ret))
	}

	// get cuda cores count
	if coreCount, ret := device.GetNumGpuCores(); ret == nvml.SUCCESS {
		gpu.Core = int64(coreCount) // CUDA cores count
		logrus.Debugf("GPU %d CUDA Cores: %d", index, coreCount)
	} else {
		logrus.Debugf("Failed to get GPU %d core count: %v", index, nvml.ErrorString(ret))
	}

	// get memory clock
	if clockMem, ret := device.GetClockInfo(nvml.CLOCK_MEM); ret == nvml.SUCCESS {
		logrus.Debugf("GPU %d Memory Clock: %d MHz", index, clockMem)
	}

	// get core clock
	if clockGpu, ret := device.GetClockInfo(nvml.CLOCK_GRAPHICS); ret == nvml.SUCCESS {
		logrus.Debugf("GPU %d Graphics Clock: %d MHz", index, clockGpu)
	}

	// get pci info
	if pciInfo, ret := device.GetPciInfo(); ret == nvml.SUCCESS {
		logrus.Debugf("GPU %d PCI: %d, Bus: %d, Device: %d", index, pciInfo.BusId, pciInfo.Bus, pciInfo.Device)
	}

	// get gpu processes
	if processes, err := nd.processDetector.GetGPUProcesses(index); err == nil {
		gpu.Processes = processes
	} else {
		logrus.Debugf("Failed to get GPU %d processes: %v", index, err)
	}

	return gpu
}

// IsNVMLAvailable check nvml is available
func (nd *NVMLDetector) IsNVMLAvailable() bool {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		return false
	}

	// shutdown nvml
	nvml.Shutdown()
	return true
}

// GetDriverVersion get nvidia driver version
func (nd *NVMLDetector) GetDriverVersion() (string, error) {
	if err := nd.initNVML(); err != nil {
		return "", err
	}
	defer nd.shutdownNVML()

	version, ret := nvml.SystemGetDriverVersion()
	if ret != nvml.SUCCESS {
		return "", fmt.Errorf("failed to get driver version: %v", nvml.ErrorString(ret))
	}

	return version, nil
}

// GetNVMLVersion get nvml version
func (nd *NVMLDetector) GetNVMLVersion() (string, error) {
	if err := nd.initNVML(); err != nil {
		return "", err
	}
	defer nd.shutdownNVML()

	version, ret := nvml.SystemGetNVMLVersion()
	if ret != nvml.SUCCESS {
		return "", fmt.Errorf("failed to get NVML version: %v", nvml.ErrorString(ret))
	}

	return version, nil
}
