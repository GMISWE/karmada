/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/06/01
 * @Desc    : NVML进程检测器 - 获取GPU进程和Pod映射关系
 */

package nvml

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
	"github.com/karmada-io/karmada/pkg/gmi/base/types"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NVMLProcessDetector nvml process detector
type NVMLProcessDetector struct {
	k8sClient kubernetes.Interface
	procRoot  string // /proc root path, support container runtime
}

// NewNVMLProcessDetector create nvml process detector
func NewNVMLProcessDetector() *NVMLProcessDetector {
	detector := &NVMLProcessDetector{
		procRoot: "/proc", // default /proc path
	}

	// check if in container
	if _, err := os.Stat("/host/proc"); err == nil {
		detector.procRoot = "/host/proc"
		logrus.Debug("检测到容器环境，使用 /host/proc")
	}

	// try to init kubernetes client
	if k8sClient, err := detector.initKubernetesClient(); err == nil {
		detector.k8sClient = k8sClient
		logrus.Debug("Kubernetes client initialized for process detection")
	} else {
		logrus.Debugf("Kubernetes client not available: %v", err)
	}

	return detector
}

// GetAllGPUProcesses get all gpu processes
func (npd *NVMLProcessDetector) GetAllGPUProcesses() ([]types.GPUProcess, error) {
	var result []types.GPUProcess

	// get gpu count
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("failed to get GPU count: %v", nvml.ErrorString(ret))
	}

	logrus.Infof("Scanning processes on %d GPU(s)...", count)

	// scan each gpu
	for i := range count {
		gpuProcesses, err := npd.GetGPUProcesses(i)
		if err != nil {
			logrus.Errorf("Failed to get processes for GPU %d: %v", i, err)
			continue
		}

		if len(gpuProcesses) > 0 {
			result = append(result, gpuProcesses...)
			logrus.Infof("GPU %d (%s): %d process(es)", i, gpuProcesses[0].ProcessName, len(gpuProcesses))
		}
	}

	return result, nil
}

// getGPUProcesses get gpu processes
func (npd *NVMLProcessDetector) GetGPUProcesses(gpuID int) ([]types.GPUProcess, error) {
	// get gpu handle
	device, ret := nvml.DeviceGetHandleByIndex(gpuID)
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("failed to get GPU %d handle: %v", gpuID, nvml.ErrorString(ret))
	}

	gpuProcesses := []types.GPUProcess{}
	// get compute processes
	if computeProcesses, ret := device.GetComputeRunningProcesses(); ret == nvml.SUCCESS {
		for _, proc := range computeProcesses {
			gpuProcess := npd.buildGPUProcess(proc.Pid, proc.UsedGpuMemory, "compute")
			gpuProcesses = append(gpuProcesses, gpuProcess)
		}
	}

	// get graphics processes
	if graphicsProcesses, ret := device.GetGraphicsRunningProcesses(); ret == nvml.SUCCESS {
		for _, proc := range graphicsProcesses {
			// check if already in compute processes (avoid duplicate)
			found := false
			for _, existing := range gpuProcesses {
				if existing.PID == proc.Pid {
					found = true
					break
				}
			}
			if !found {
				gpuProcess := npd.buildGPUProcess(proc.Pid, proc.UsedGpuMemory, "graphics")
				gpuProcesses = append(gpuProcesses, gpuProcess)
			}
		}
	}

	logrus.Debugf("GPU %d: Found %d process(es)", gpuID, len(gpuProcesses))

	return gpuProcesses, nil
}

// buildGPUProcess build gpu process
func (npd *NVMLProcessDetector) buildGPUProcess(pid uint32, memoryUsed uint64, processType string) types.GPUProcess {
	process := types.GPUProcess{
		PID:           pid,
		GPUMemoryUsed: memoryUsed / (1024 * 1024), // convert to MB
	}

	// get process name
	if processName := npd.getProcessName(pid); processName != "" {
		process.ProcessName = processName
	} else {
		process.ProcessName = fmt.Sprintf("pid-%d", pid)
	}

	// get container id
	if containerID := npd.getContainerID(pid); containerID != "" {
		process.ContainerID = containerID

		// 获取Pod信息
		if podInfo := npd.getPodInfoByContainerID(containerID); podInfo != nil {
			process.PodName = podInfo.Name
			process.PodNamespace = podInfo.Namespace
			process.PodUID = podInfo.UID
		}
	}

	logrus.Debugf("Process: PID=%d, Name=%s, Memory=%dMB, Pod=%s/%s",
		process.PID, process.ProcessName, process.GPUMemoryUsed,
		process.PodNamespace, process.PodName)

	return process
}

// getProcessName get process name
func (npd *NVMLProcessDetector) getProcessName(pid uint32) string {
	cmdlinePath := fmt.Sprintf("%s/%d/cmdline", npd.procRoot, pid)
	if content, err := os.ReadFile(cmdlinePath); err == nil {
		// 处理cmdline中的null字符
		cmdline := strings.ReplaceAll(string(content), "\x00", " ")
		cmdline = strings.TrimSpace(cmdline)
		if cmdline != "" {
			// 只取第一个命令
			if parts := strings.Fields(cmdline); len(parts) > 0 {
				return filepath.Base(parts[0])
			}
		}
	}

	// 尝试读取comm文件
	commPath := fmt.Sprintf("%s/%d/comm", npd.procRoot, pid)
	if content, err := os.ReadFile(commPath); err == nil {
		return strings.TrimSpace(string(content))
	}

	return ""
}

// getContainerID 通过PID获取容器ID
func (npd *NVMLProcessDetector) getContainerID(pid uint32) string {
	cgroupPath := fmt.Sprintf("%s/%d/cgroup", npd.procRoot, pid)
	file, err := os.Open(cgroupPath)
	if err != nil {
		logrus.Debugf("无法读取cgroup文件 %s: %v", cgroupPath, err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 解析cgroup路径中的容器ID
		// 支持多种容器运行时格式
		if containerID := npd.extractContainerIDFromCgroup(line); containerID != "" {
			logrus.Debugf("Found container ID: %s from cgroup: %s", containerID, line)
			return containerID
		}
	}

	return ""
}

// extractContainerIDFromCgroup extract container id from cgroup path
func (npd *NVMLProcessDetector) extractContainerIDFromCgroup(cgroupLine string) string {
	// Docker格式: /docker/<container_id>
	if strings.Contains(cgroupLine, "/docker/") {
		parts := strings.Split(cgroupLine, "/docker/")
		if len(parts) > 1 {
			containerID := strings.Split(parts[1], "/")[0]
			if len(containerID) >= 12 { // 容器ID至少12位
				return containerID
			}
		}
	}

	// Containerd格式: /containerd/<container_id>
	if strings.Contains(cgroupLine, "/containerd/") {
		parts := strings.Split(cgroupLine, "/containerd/")
		if len(parts) > 1 {
			containerID := strings.Split(parts[1], "/")[0]
			if len(containerID) >= 12 {
				return containerID
			}
		}
	}

	// CRI-O格式: /crio-<container_id>
	if strings.Contains(cgroupLine, "/crio-") {
		parts := strings.Split(cgroupLine, "/crio-")
		if len(parts) > 1 {
			containerID := strings.Split(parts[1], ".")[0]
			if len(containerID) >= 12 {
				return containerID
			}
		}
	}

	// Kubernetes Pod格式: 提取Pod相关信息
	if strings.Contains(cgroupLine, "/kubepods/") {
		// 尝试从路径中提取容器ID
		if idx := strings.LastIndex(cgroupLine, "/"); idx != -1 {
			containerID := cgroupLine[idx+1:]
			// 移除可能的scope后缀
			if idx := strings.Index(containerID, ".scope"); idx != -1 {
				containerID = containerID[:idx]
			}
			if len(containerID) >= 12 && !strings.Contains(containerID, ".") {
				return containerID
			}
		}
	}

	// systemd格式: docker-<container_id>.scope
	if strings.Contains(cgroupLine, "docker-") && strings.Contains(cgroupLine, ".scope") {
		parts := strings.Split(cgroupLine, "docker-")
		if len(parts) > 1 {
			containerID := strings.Split(parts[1], ".scope")[0]
			if len(containerID) >= 12 {
				return containerID
			}
		}
	}

	return ""
}

// PodInfo pod basic info
type PodInfo struct {
	Name      string
	Namespace string
	UID       string
}

// getPodInfoByContainerID 通过容器ID获取Pod信息
func (npd *NVMLProcessDetector) getPodInfoByContainerID(containerID string) *PodInfo {
	// 优先尝试从环境变量获取当前Pod信息
	if podInfo := npd.getPodInfoFromEnv(); podInfo != nil {
		logrus.Debugf("Found current Pod info from environment: %s/%s", podInfo.Namespace, podInfo.Name)
		return podInfo
	}

	// 尝试从Downward API挂载的文件获取
	if podInfo := npd.getPodInfoFromDownwardAPI(); podInfo != nil {
		logrus.Debugf("Found current Pod info from Downward API: %s/%s", podInfo.Namespace, podInfo.Name)
		return podInfo
	}

	// 如果Kubernetes客户端可用，作为后备方案
	if npd.k8sClient != nil {
		return npd.getPodInfoFromAPIServer(containerID)
	}

	logrus.Debugf("No method available to get Pod info for container: %s", containerID[:12])
	return nil
}

// getPodInfoFromEnv 从环境变量获取Pod信息
func (npd *NVMLProcessDetector) getPodInfoFromEnv() *PodInfo {
	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	podUID := os.Getenv("POD_UID")

	// 检查是否有Kubernetes注入的标准环境变量
	if podName == "" {
		podName = os.Getenv("HOSTNAME") // Pod名称通常等于hostname
	}
	if podNamespace == "" {
		podNamespace = os.Getenv("NAMESPACE")
	}

	if podName != "" && podNamespace != "" {
		return &PodInfo{
			Name:      podName,
			Namespace: podNamespace,
			UID:       podUID,
		}
	}

	return nil
}

// getPodInfoFromDownwardAPI 从Downward API文件获取Pod信息
func (npd *NVMLProcessDetector) getPodInfoFromDownwardAPI() *PodInfo {
	podInfo := &PodInfo{}

	// 常见的Downward API挂载路径
	paths := []string{
		"/etc/podinfo",
		"/var/run/secrets/kubernetes.io/podinfo",
		"/pod-info",
	}

	for _, basePath := range paths {
		if npd.readPodInfoFromPath(basePath, podInfo) {
			if podInfo.Name != "" && podInfo.Namespace != "" {
				return podInfo
			}
		}
	}

	return nil
}

// readPodInfoFromPath 从指定路径读取Pod信息
func (npd *NVMLProcessDetector) readPodInfoFromPath(basePath string, podInfo *PodInfo) bool {
	found := false

	// 读取Pod名称
	if content, err := os.ReadFile(filepath.Join(basePath, "name")); err == nil {
		podInfo.Name = strings.TrimSpace(string(content))
		found = true
	}

	// 读取命名空间
	if content, err := os.ReadFile(filepath.Join(basePath, "namespace")); err == nil {
		podInfo.Namespace = strings.TrimSpace(string(content))
		found = true
	}

	// 读取UID
	if content, err := os.ReadFile(filepath.Join(basePath, "uid")); err == nil {
		podInfo.UID = strings.TrimSpace(string(content))
		found = true
	}

	// 读取标签文件（可能包含Pod信息）
	if content, err := os.ReadFile(filepath.Join(basePath, "labels")); err == nil {
		npd.parsePodLabels(string(content), podInfo)
		found = true
	}

	// 读取注解文件
	if content, err := os.ReadFile(filepath.Join(basePath, "annotations")); err == nil {
		npd.parsePodAnnotations(string(content), podInfo)
		found = true
	}

	return found
}

// parsePodLabels 解析Pod标签
func (npd *NVMLProcessDetector) parsePodLabels(content string, podInfo *PodInfo) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

				// 从标签中提取有用信息
				switch key {
				case "app", "app.kubernetes.io/name":
					if podInfo.Name == "" {
						podInfo.Name = value
					}
				}
			}
		}
	}
}

// parsePodAnnotations 解析Pod注解
func (npd *NVMLProcessDetector) parsePodAnnotations(content string, podInfo *PodInfo) {
	// 类似于标签解析，可以从注解中提取额外信息
	// 这里可以根据需要添加特定的注解解析逻辑
}

// getPodInfoFromAPIServer 从API Server获取Pod信息（后备方案）
func (npd *NVMLProcessDetector) getPodInfoFromAPIServer(containerID string) *PodInfo {
	if npd.k8sClient == nil {
		logrus.Debug("Kubernetes client not available for Pod lookup")
		return nil
	}

	// 查询所有Pod
	pods, err := npd.k8sClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logrus.Debugf("Failed to list pods: %v", err)
		return nil
	}

	// 遍历所有Pod查找匹配的容器
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.ContainerID != "" {
				// 提取容器ID（去掉runtime前缀）
				statusContainerID := npd.extractContainerIDFromStatus(containerStatus.ContainerID)
				if npd.containerIDMatches(statusContainerID, containerID) {
					logrus.Debugf("Found Pod via API Server: %s/%s for container %s",
						pod.Namespace, pod.Name, containerID[:12])
					return &PodInfo{
						Name:      pod.Name,
						Namespace: pod.Namespace,
						UID:       string(pod.UID),
					}
				}
			}
		}
	}

	logrus.Debugf("No Pod found via API Server for container ID: %s", containerID[:12])
	return nil
}

// containerIDMatches check if two container ids match
func (npd *NVMLProcessDetector) containerIDMatches(id1, id2 string) bool {
	// exact match
	if id1 == id2 {
		return true
	}

	// 前缀匹配（至少12位）
	minLen := 12
	if len(id1) >= minLen && len(id2) >= minLen {
		if strings.HasPrefix(id1, id2[:minLen]) || strings.HasPrefix(id2, id1[:minLen]) {
			return true
		}
	}

	return false
}

// extractContainerIDFromStatus extract container id from container status
func (npd *NVMLProcessDetector) extractContainerIDFromStatus(containerID string) string {
	// 格式: docker://abc123... 或 containerd://abc123...
	if strings.Contains(containerID, "://") {
		parts := strings.Split(containerID, "://")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return containerID
}

// initKubernetesClient init kubernetes client
func (npd *NVMLProcessDetector) initKubernetesClient() (kubernetes.Interface, error) {
	// try to get in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %v", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	return client, nil
}

// PrintGPUProcesses print gpu processes
func (npd *NVMLProcessDetector) PrintGPUProcesses() error {
	processes, err := npd.GetAllGPUProcesses()
	if err != nil {
		return err
	}

	if len(processes) == 0 {
		logrus.Info("No GPU processes found")
		return nil
	}

	logrus.Info("GPU Process Information:")
	logrus.Info(strings.Repeat("=", 80))

	for _, proc := range processes {
		logrus.Infof("  PID: %d | Process: %s | Memory: %d MB",
			proc.PID, proc.ProcessName, proc.GPUMemoryUsed)

		if proc.PodName != "" {
			logrus.Infof("     Pod: %s/%s (UID: %s)",
				proc.PodNamespace, proc.PodName, proc.PodUID)
		} else if proc.ContainerID != "" {
			logrus.Infof("     Container: %s", proc.ContainerID[:12])
		} else {
			logrus.Info("     Host Process")
		}
		logrus.Info("")
	}

	return nil
}
