/*
 @Version : 1.0
 @Author  : wangxiaokang
 @Email   : xiaokang.w@gmicloud.ai
 @Time    : 2025/05/30 11:03:14
 Desc     :
*/

package nfd

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/karmada-io/karmada/pkg/gmi/base/types"
	"github.com/karmada-io/karmada/pkg/util"
	"github.com/sirupsen/logrus"
)

type HostDetector struct {
	timeout time.Duration
}

// OSInfo holds basic operating system information
type OSInfo struct {
	OS              string
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	KernelVersion   string
}

func NewHostDetector() *HostDetector {
	timeout := util.GetEnv("HOST_DETECTOR_TIMEOUT", "100ms")
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		logrus.Fatalf("Failed to parse HOST_DETECTOR_TIMEOUT: %v", err)
	}
	return &HostDetector{
		timeout: timeoutDuration,
	}
}

func (hd *HostDetector) DiscoverHost() types.Host {
	hostname := hd.fetchHostname()

	// Get basic system information from /proc and /sys
	osInfo := hd.getOSInfo()
	cpuInfo := hd.getCPUInfo()
	memInfo := hd.getMemoryInfo()

	return types.Host{
		Hostname:        hostname,
		Os:              osInfo.OS,
		Platform:        osInfo.Platform,
		PlatformFamily:  osInfo.PlatformFamily,
		PlatformVersion: osInfo.PlatformVersion,
		KernalVersion:   osInfo.KernelVersion,
		KernelArch:      runtime.GOARCH,
		Cpu:             cpuInfo,
		Mem:             memInfo,
	}
}

// getOSInfo retrieves operating system information from /proc and /etc files
func (hd *HostDetector) getOSInfo() OSInfo {
	osInfo := OSInfo{
		OS:             runtime.GOOS,
		Platform:       "unknown",
		PlatformFamily: "unknown",
		KernelVersion:  "unknown",
	}

	// Get kernel version from /proc/version
	if version, err := util.FileToString("/proc/version", true); err == nil && version != "" {
		// Extract kernel version from the version string
		parts := strings.Fields(version)
		if len(parts) >= 3 {
			osInfo.KernelVersion = parts[2]
		}
	}

	// Try to get OS release information from /etc/os-release
	if osRelease, err := util.FileToString("/etc/os-release", true); err == nil && osRelease != "" {
		osInfo.Platform, osInfo.PlatformFamily, osInfo.PlatformVersion = hd.parseOSRelease(osRelease)
	}

	return osInfo
}

// FetchNodeUUID retrieves the host machine UUID from various sources
// Priority order: /etc/machine-id -> /var/lib/dbus/machine-id -> /sys/class/dmi/id/product_uuid -> /proc/sys/kernel/random/uuid
func (hd *HostDetector) FetchNodeUUID() string {
	// Try multiple sources for machine UUID
	sources := []string{
		"/etc/machine-id",
		"/var/lib/dbus/machine-id",
		"/sys/class/dmi/id/product_uuid",
		"/proc/sys/kernel/random/uuid",
	}

	for _, source := range sources {
		if uuid, err := util.FileToString(source, true); err == nil && uuid != "" {
			logrus.Debugf("Successfully read UUID from %s: %s", source, uuid)
			return uuid
		}
	}

	// Fallback: generate a UUID based on hostname and current time
	hostname := hd.fetchHostname()
	fallbackUUID := hd.generateFallbackUUID(hostname)
	logrus.Warnf("Could not read machine UUID from any source, using fallback: %s", fallbackUUID)
	return fallbackUUID
}

// generateFallbackUUID creates a deterministic UUID-like string based on hostname
func (hd *HostDetector) generateFallbackUUID(hostname string) string {
	// Simple fallback UUID generation based on hostname and timestamp
	// In production, you might want to use a proper UUID library
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%s-%d", strings.ToLower(hostname), util.GetEnv("FALLBACK_UUID_SUFFIX", "fallback"), timestamp%10000)
}

// FetchHostname retrieves the hostname from various sources
// Priority order: os.Hostname() -> /etc/hostname -> /proc/sys/kernel/hostname
func (hd *HostDetector) fetchHostname() string {
	// Try os.Hostname() first
	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		logrus.Debugf("Successfully got hostname from os.Hostname(): %s", hostname)
		return hostname
	}

	// Try reading from files
	sources := []string{
		"/etc/hostname",
		"/proc/sys/kernel/hostname",
	}

	for _, source := range sources {
		if hostname, err := util.FileToString(source, true); err == nil && hostname != "" {
			logrus.Debugf("Successfully read hostname from %s: %s", source, hostname)
			return hostname
		}
	}

	// Fallback to "unknown"
	logrus.Warn("Could not determine hostname from any source, using 'unknown'")
	return "unknown"
}

// parseOSRelease parses /etc/os-release content to extract platform information
func (hd *HostDetector) parseOSRelease(content string) (platform, family, version string) {
	platform = "unknown"
	family = "unknown"
	version = "unknown"

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ID=") {
			platform = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		} else if strings.HasPrefix(line, "ID_LIKE=") {
			family = strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), "\"")
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		}
	}

	if family == "unknown" {
		family = platform
	}

	return platform, family, version
}

// getCPUInfo retrieves CPU information from /proc/cpuinfo
func (hd *HostDetector) getCPUInfo() *types.Cpu {
	cpuInfo := &types.Cpu{
		Cores: int32(runtime.NumCPU()),
	}

	// Read CPU information from /proc/cpuinfo
	if cpuData, err := util.FileToString("/proc/cpuinfo", true); err == nil && cpuData != "" {
		cpuInfo.ModelName, cpuInfo.Family, cpuInfo.Mhz = hd.parseCPUInfo(cpuData)
	}

	// Get CPU usage from /proc/stat
	cpuInfo.Used = hd.getCPUUsage()

	return cpuInfo
}

// parseCPUInfo parses /proc/cpuinfo to extract CPU details
func (hd *HostDetector) parseCPUInfo(content string) (modelName, family string, mhz float64) {
	modelName = "unknown"
	family = "unknown"
	mhz = 0.0

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				switch key {
				case "model name":
					if modelName == "unknown" { // Only take the first CPU's info
						modelName = value
					}
				case "cpu family":
					if family == "unknown" {
						family = value
					}
				case "cpu MHz":
					if mhz == 0.0 {
						if parsedMhz, err := strconv.ParseFloat(value, 64); err == nil {
							mhz = parsedMhz
						}
					}
				}
			}
		}
	}

	return modelName, family, mhz
}

// getCPUUsage calculates CPU usage percentage from /proc/stat
func (hd *HostDetector) getCPUUsage() float64 {
	// Read /proc/stat twice with a small interval to calculate usage
	stat1 := hd.readCPUStat()
	time.Sleep(100 * time.Millisecond)
	stat2 := hd.readCPUStat()

	if stat1 == nil || stat2 == nil {
		return 0.0
	}

	// Calculate CPU usage percentage
	idle1 := stat1[3] + stat1[4] // idle + iowait
	idle2 := stat2[3] + stat2[4]

	total1 := int64(0)
	total2 := int64(0)
	for _, v := range stat1 {
		total1 += v
	}
	for _, v := range stat2 {
		total2 += v
	}

	totalDiff := total2 - total1
	idleDiff := idle2 - idle1

	if totalDiff == 0 {
		return 0.0
	}

	return float64(totalDiff-idleDiff) / float64(totalDiff) * 100.0
}

// readCPUStat reads CPU statistics from /proc/stat
func (hd *HostDetector) readCPUStat() []int64 {
	content, err := util.FileToString("/proc/stat", true)
	if err != nil {
		return nil
	}
	if content == "" {
		return nil
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			if len(fields) >= 8 {
				stats := make([]int64, len(fields)-1)
				for i := 1; i < len(fields); i++ {
					if val, err := strconv.ParseInt(fields[i], 10, 64); err == nil {
						stats[i-1] = val
					}
				}
				return stats
			}
		}
	}

	return nil
}

// getMemoryInfo retrieves memory information from /proc/meminfo
func (hd *HostDetector) getMemoryInfo() *types.Mem {
	memInfo := &types.Mem{}

	content, err := util.FileToString("/proc/meminfo", true)
	if err != nil {
		return memInfo
	}
	if content == "" {
		return memInfo
	}

	var memTotal, memFree, memActive uint64

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// Remove "kB" suffix and parse the number
				value = strings.TrimSuffix(value, " kB")
				if val, err := strconv.ParseUint(value, 10, 64); err == nil {
					val *= 1024 // Convert from kB to bytes

					switch key {
					case "MemTotal":
						memTotal = val
					case "MemFree":
						memFree = val
					case "Active":
						memActive = val
					}
				}
			}
		}
	}

	memInfo.Size = hd.formatBytes(memTotal)
	memInfo.Free = hd.formatBytes(memFree)
	memInfo.Active = hd.formatBytes(memActive)

	// Calculate used percentage
	if memTotal > 0 {
		used := memTotal - memFree
		memInfo.Used = float64(used) / float64(memTotal) * 100.0
	}

	return memInfo
}

// formatBytes converts bytes to human readable format
func (hd *HostDetector) formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
