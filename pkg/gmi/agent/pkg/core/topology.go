/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/06/01
 * @Desc    : 资源拓扑数据结构，用于任务调度
 */

package core

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/karmada-io/karmada/pkg/gmi/base/types"
)

// ====================
// 资源度量和状态定义
// ====================

// ResourceMetrics 资源度量
type ResourceMetrics struct {
	CPU     CPUMetrics     `json:"cpu"`
	Memory  MemoryMetrics  `json:"memory"`
	GPU     GPUMetrics     `json:"gpu"`
	Disk    DiskMetrics    `json:"disk"`
	Network NetworkMetrics `json:"network"`
}

// CPUMetrics CPU资源度量
type CPUMetrics struct {
	Total     int32   `json:"total"`     // 总核数
	Available int32   `json:"available"` // 可用核数
	Used      int32   `json:"used"`      // 已使用核数
	Usage     float64 `json:"usage"`     // 使用率 0-1
}

// MemoryMetrics 内存资源度量
type MemoryMetrics struct {
	TotalGB     float64 `json:"total_gb"`     // 总内存GB
	AvailableGB float64 `json:"available_gb"` // 可用内存GB
	UsedGB      float64 `json:"used_gb"`      // 已使用内存GB
	Usage       float64 `json:"usage"`        // 使用率 0-1
}

// GPUMetrics GPU资源度量
type GPUMetrics struct {
	Total             int      `json:"total"`               // GPU总数
	Available         int      `json:"available"`           // 可用GPU数
	Used              int      `json:"used"`                // 已使用GPU数
	Models            []string `json:"models"`              // GPU型号列表
	TotalMemoryGB     float64  `json:"total_memory_gb"`     // 总显存GB
	AvailableMemoryGB float64  `json:"available_memory_gb"` // 可用显存GB
}

// DiskMetrics 磁盘资源度量
type DiskMetrics struct {
	TotalGB     float64 `json:"total_gb"`     // 总空间GB
	AvailableGB float64 `json:"available_gb"` // 可用空间GB
	UsedGB      float64 `json:"used_gb"`      // 已使用空间GB
	Usage       float64 `json:"usage"`        // 使用率 0-1
	IOPSLimit   int     `json:"iops_limit"`   // IOPS限制
}

// NetworkMetrics 网络资源度量
type NetworkMetrics struct {
	Bandwidth   float64 `json:"bandwidth"`   // 带宽 Mbps
	Latency     float64 `json:"latency"`     // 延迟 ms
	PacketLoss  float64 `json:"packet_loss"` // 丢包率 0-1
	Connections int     `json:"connections"` // 当前连接数
}

// NodeStatus 节点状态
type NodeStatus int

const (
	NodeStatusUnknown     NodeStatus = iota
	NodeStatusReady                  // 就绪
	NodeStatusBusy                   // 繁忙
	NodeStatusMaintenance            // 维护中
	NodeStatusOffline                // 离线
	NodeStatusError                  // 错误状态
)

func (s NodeStatus) String() string {
	switch s {
	case NodeStatusReady:
		return "Ready"
	case NodeStatusBusy:
		return "Busy"
	case NodeStatusMaintenance:
		return "Maintenance"
	case NodeStatusOffline:
		return "Offline"
	case NodeStatusError:
		return "Error"
	default:
		return "Unknown"
	}
}

// ====================
// 节点和拓扑结构定义
// ====================

// NodeLabels 节点标签，用于调度策略
type NodeLabels map[string]string

// TopologyNode 拓扑节点
type TopologyNode struct {
	// 基础信息
	UUID    string     `json:"uuid"`
	Name    string     `json:"name"`
	Labels  NodeLabels `json:"labels"`
	Zone    string     `json:"zone"`    // 可用区
	Region  string     `json:"region"`  // 地区
	Cluster string     `json:"cluster"` // 集群

	// 节点状态
	Status      NodeStatus `json:"status"`
	LastUpdated time.Time  `json:"last_updated"`
	Health      float64    `json:"health"` // 健康度 0-1

	// 资源信息
	Resources ResourceMetrics `json:"resources"`
	RawNode   types.Node      `json:"raw_node"` // 原始节点数据

	// 调度相关
	Taints      []NodeTaint       `json:"taints"`      // 污点
	Annotations map[string]string `json:"annotations"` // 注解
	Priority    int               `json:"priority"`    // 优先级
}

// NodeTaint 节点污点
type NodeTaint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"` // NoSchedule, PreferNoSchedule, NoExecute
}

// ====================
// 资源拓扑管理器
// ====================

// ResourceTopology 资源拓扑管理器
type ResourceTopology struct {
	mu    sync.RWMutex
	nodes map[string]*TopologyNode // UUID -> Node

	// 索引
	zoneNodes    map[string][]*TopologyNode // Zone -> Nodes
	regionNodes  map[string][]*TopologyNode // Region -> Nodes
	clusterNodes map[string][]*TopologyNode // Cluster -> Nodes
	labelIndex   map[string][]*TopologyNode // Label -> Nodes

	// 统计信息
	totalResources     ResourceMetrics `json:"total_resources"`
	availableResources ResourceMetrics `json:"available_resources"`
	lastUpdated        time.Time       `json:"last_updated"`
}

// NewResourceTopology 创建资源拓扑管理器
func NewResourceTopology() *ResourceTopology {
	return &ResourceTopology{
		nodes:        make(map[string]*TopologyNode),
		zoneNodes:    make(map[string][]*TopologyNode),
		regionNodes:  make(map[string][]*TopologyNode),
		clusterNodes: make(map[string][]*TopologyNode),
		labelIndex:   make(map[string][]*TopologyNode),
	}
}

// AddNode 添加节点到拓扑
func (rt *ResourceTopology) AddNode(node types.Node) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	topologyNode := &TopologyNode{
		UUID:        node.UUID,
		Name:        node.Host.Hostname,
		Zone:        extractZone(node),
		Region:      extractRegion(node),
		Cluster:     extractCluster(node),
		Status:      NodeStatusReady,
		LastUpdated: time.Now(),
		Health:      1.0,
		Resources:   convertToResourceMetrics(node),
		RawNode:     node,
		Labels:      generateNodeLabels(node),
		Annotations: make(map[string]string),
		Priority:    calculateNodePriority(node),
	}

	// 添加到主索引
	rt.nodes[node.UUID] = topologyNode

	// 更新各种索引
	rt.updateIndexes(topologyNode)

	// 重新计算统计信息
	rt.updateStatistics()

	return nil
}

// UpdateNode 更新节点信息
func (rt *ResourceTopology) UpdateNode(node types.Node) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	existing, exists := rt.nodes[node.UUID]
	if !exists {
		return rt.AddNode(node)
	}

	// 更新节点信息
	existing.Resources = convertToResourceMetrics(node)
	existing.RawNode = node
	existing.LastUpdated = time.Now()
	existing.Health = calculateNodeHealth(node)

	// 重新计算统计信息
	rt.updateStatistics()

	return nil
}

// RemoveNode 移除节点
func (rt *ResourceTopology) RemoveNode(uuid string) {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	node, exists := rt.nodes[uuid]
	if !exists {
		return
	}

	// 从主索引删除
	delete(rt.nodes, uuid)

	// 从其他索引删除
	rt.removeFromIndexes(node)

	// 重新计算统计信息
	rt.updateStatistics()
}

// GetNode 获取指定节点
func (rt *ResourceTopology) GetNode(uuid string) (*TopologyNode, bool) {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	node, exists := rt.nodes[uuid]
	return node, exists
}

// GetAllNodes 获取所有节点
func (rt *ResourceTopology) GetAllNodes() []*TopologyNode {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	nodes := make([]*TopologyNode, 0, len(rt.nodes))
	for _, node := range rt.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

// GetNodesByZone 按可用区获取节点
func (rt *ResourceTopology) GetNodesByZone(zone string) []*TopologyNode {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	return rt.zoneNodes[zone]
}

// GetNodesByRegion 按地区获取节点
func (rt *ResourceTopology) GetNodesByRegion(region string) []*TopologyNode {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	return rt.regionNodes[region]
}

// GetNodesByLabel 按标签获取节点
func (rt *ResourceTopology) GetNodesByLabel(key, value string) []*TopologyNode {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	var result []*TopologyNode
	for _, node := range rt.nodes {
		if labelValue, exists := node.Labels[key]; exists && labelValue == value {
			result = append(result, node)
		}
	}
	return result
}

// GetAvailableResources 获取可用资源统计
func (rt *ResourceTopology) GetAvailableResources() ResourceMetrics {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	return rt.availableResources
}

// GetTotalResources 获取总资源统计
func (rt *ResourceTopology) GetTotalResources() ResourceMetrics {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	return rt.totalResources
}

// ====================
// 内部辅助方法
// ====================

// updateIndexes 更新索引
func (rt *ResourceTopology) updateIndexes(node *TopologyNode) {
	// 更新zone索引
	rt.zoneNodes[node.Zone] = append(rt.zoneNodes[node.Zone], node)

	// 更新region索引
	rt.regionNodes[node.Region] = append(rt.regionNodes[node.Region], node)

	// 更新cluster索引
	rt.clusterNodes[node.Cluster] = append(rt.clusterNodes[node.Cluster], node)

	// 更新label索引
	for key, value := range node.Labels {
		labelKey := key + "=" + value
		rt.labelIndex[labelKey] = append(rt.labelIndex[labelKey], node)
	}
}

// removeFromIndexes 从索引中移除
func (rt *ResourceTopology) removeFromIndexes(node *TopologyNode) {
	// 从zone索引移除
	if nodes, exists := rt.zoneNodes[node.Zone]; exists {
		rt.zoneNodes[node.Zone] = rt.removeNodeFromSlice(nodes, node)
	}

	// 从region索引移除
	if nodes, exists := rt.regionNodes[node.Region]; exists {
		rt.regionNodes[node.Region] = rt.removeNodeFromSlice(nodes, node)
	}

	// 从cluster索引移除
	if nodes, exists := rt.clusterNodes[node.Cluster]; exists {
		rt.clusterNodes[node.Cluster] = rt.removeNodeFromSlice(nodes, node)
	}

	// 从label索引移除
	for key, value := range node.Labels {
		labelKey := key + "=" + value
		if nodes, exists := rt.labelIndex[labelKey]; exists {
			rt.labelIndex[labelKey] = rt.removeNodeFromSlice(nodes, node)
		}
	}
}

// removeNodeFromSlice 从切片中移除节点
func (rt *ResourceTopology) removeNodeFromSlice(slice []*TopologyNode, node *TopologyNode) []*TopologyNode {
	for i, n := range slice {
		if n.UUID == node.UUID {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// updateStatistics 更新统计信息
func (rt *ResourceTopology) updateStatistics() {
	var total, available ResourceMetrics

	for _, node := range rt.nodes {
		if node.Status != NodeStatusReady {
			continue
		}

		// 累加总资源
		total.CPU.Total += node.Resources.CPU.Total
		total.Memory.TotalGB += node.Resources.Memory.TotalGB
		total.GPU.Total += node.Resources.GPU.Total
		total.Disk.TotalGB += node.Resources.Disk.TotalGB

		// 累加可用资源
		available.CPU.Total += node.Resources.CPU.Available
		available.Memory.TotalGB += node.Resources.Memory.AvailableGB
		available.GPU.Total += node.Resources.GPU.Available
		available.Disk.TotalGB += node.Resources.Disk.AvailableGB
	}

	rt.totalResources = total
	rt.availableResources = available
	rt.lastUpdated = time.Now()
}

// ====================
// 辅助函数实现
// ====================

// extractZone 从节点信息中提取可用区
func extractZone(node types.Node) string {
	// 可以根据主机名或其他信息推断
	hostname := node.Host.Hostname
	if strings.Contains(hostname, "zone") {
		parts := strings.Split(hostname, "-")
		for _, part := range parts {
			if strings.HasPrefix(part, "zone") {
				return part
			}
		}
	}
	return "zone-default"
}

// extractRegion 从节点信息中提取地区
func extractRegion(node types.Node) string {
	// 可以根据主机名或其他信息推断
	hostname := node.Host.Hostname
	if strings.Contains(hostname, "region") {
		parts := strings.Split(hostname, "-")
		for _, part := range parts {
			if strings.HasPrefix(part, "region") {
				return part
			}
		}
	}
	return "region-default"
}

// extractCluster 从节点信息中提取集群信息
func extractCluster(node types.Node) string {
	// 可以根据主机名或其他信息推断
	hostname := node.Host.Hostname
	if strings.Contains(hostname, "cluster") {
		parts := strings.Split(hostname, "-")
		for _, part := range parts {
			if strings.HasPrefix(part, "cluster") {
				return part
			}
		}
	}
	return "cluster-default"
}

// convertToResourceMetrics 将原始节点数据转换为资源度量
func convertToResourceMetrics(node types.Node) ResourceMetrics {
	metrics := ResourceMetrics{}

	// CPU度量
	if node.Host.Cpu != nil {
		metrics.CPU.Total = node.Host.Cpu.Cores
		metrics.CPU.Usage = node.Host.Cpu.Used
		metrics.CPU.Used = int32(float64(node.Host.Cpu.Cores) * node.Host.Cpu.Used)
		metrics.CPU.Available = metrics.CPU.Total - metrics.CPU.Used
	}

	// 内存度量
	if node.Host.Mem != nil {
		totalGB := parseMemorySize(node.Host.Mem.Size)
		usedGB := totalGB * node.Host.Mem.Used
		metrics.Memory.TotalGB = totalGB
		metrics.Memory.UsedGB = usedGB
		metrics.Memory.AvailableGB = totalGB - usedGB
		metrics.Memory.Usage = node.Host.Mem.Used
	}

	// GPU度量
	metrics.GPU.Total = len(node.GPUs)
	metrics.GPU.Available = len(node.GPUs) // 简化处理，假设都可用
	models := make([]string, 0)
	for _, gpu := range node.GPUs {
		if gpu.Model != "" {
			models = append(models, gpu.Model)
		}
		// 累加GPU显存
		totalMem := parseMemorySize(gpu.UseMem) + parseMemorySize(gpu.FreeMem)
		metrics.GPU.TotalMemoryGB += totalMem
		metrics.GPU.AvailableMemoryGB += parseMemorySize(gpu.FreeMem)
	}
	metrics.GPU.Models = models

	// 磁盘度量
	for _, disk := range node.Host.Disk {
		totalGB := parseMemorySize(disk.Total)
		usedGB := totalGB * disk.Used
		metrics.Disk.TotalGB += totalGB
		metrics.Disk.UsedGB += usedGB
		metrics.Disk.AvailableGB += totalGB - usedGB
	}
	if metrics.Disk.TotalGB > 0 {
		metrics.Disk.Usage = metrics.Disk.UsedGB / metrics.Disk.TotalGB
	}

	return metrics
}

// parseMemorySize 解析内存大小字符串，返回GB数
func parseMemorySize(sizeStr string) float64 {
	if sizeStr == "" {
		return 0
	}

	// 移除空格
	sizeStr = strings.TrimSpace(sizeStr)

	// 提取数字部分
	var numStr strings.Builder
	var unit string

	for i, r := range sizeStr {
		if r >= '0' && r <= '9' || r == '.' {
			numStr.WriteRune(r)
		} else {
			unit = sizeStr[i:]
			break
		}
	}

	num, err := strconv.ParseFloat(numStr.String(), 64)
	if err != nil {
		return 0
	}

	// 转换单位到GB
	unit = strings.ToUpper(strings.TrimSpace(unit))
	switch unit {
	case "B", "BYTES":
		return num / (1024 * 1024 * 1024)
	case "KB", "KIB":
		return num / (1024 * 1024)
	case "MB", "MIB":
		return num / 1024
	case "GB", "GIB":
		return num
	case "TB", "TIB":
		return num * 1024
	default:
		// 默认当作字节处理
		return num / (1024 * 1024 * 1024)
	}
}

// generateNodeLabels 生成节点标签
func generateNodeLabels(node types.Node) NodeLabels {
	labels := make(NodeLabels)

	// 基础标签
	labels["hostname"] = node.Host.Hostname
	labels["os"] = node.Host.Os
	labels["platform"] = node.Host.Platform
	labels["platform_family"] = node.Host.PlatformFamily
	labels["kernel_arch"] = node.Host.KernelArch

	// CPU标签
	if node.Host.Cpu != nil {
		labels["cpu_cores"] = strconv.Itoa(int(node.Host.Cpu.Cores))
		labels["cpu_family"] = node.Host.Cpu.Family
		labels["cpu_model"] = node.Host.Cpu.ModelName
	}

	// GPU标签
	if len(node.GPUs) > 0 {
		labels["gpu_count"] = strconv.Itoa(len(node.GPUs))
		// 取第一个GPU的型号作为代表
		if node.GPUs[0].Model != "" {
			labels["gpu_model"] = node.GPUs[0].Model
		}
		labels["has_gpu"] = "true"
	} else {
		labels["has_gpu"] = "false"
	}

	return labels
}

// calculateNodePriority 计算节点优先级
func calculateNodePriority(node types.Node) int {
	priority := 0

	// 基于CPU核数
	if node.Host.Cpu != nil {
		priority += int(node.Host.Cpu.Cores) * 10
	}

	// 基于GPU数量
	priority += len(node.GPUs) * 100

	// 基于内存大小（简化计算）
	if node.Host.Mem != nil {
		memGB := parseMemorySize(node.Host.Mem.Size)
		priority += int(memGB)
	}

	return priority
}

// calculateNodeHealth 计算节点健康度
func calculateNodeHealth(node types.Node) float64 {
	health := 1.0

	// 基于CPU使用率
	if node.Host.Cpu != nil {
		if node.Host.Cpu.Used > 0.9 {
			health -= 0.3
		} else if node.Host.Cpu.Used > 0.7 {
			health -= 0.1
		}
	}

	// 基于内存使用率
	if node.Host.Mem != nil {
		if node.Host.Mem.Used > 0.9 {
			health -= 0.3
		} else if node.Host.Mem.Used > 0.7 {
			health -= 0.1
		}
	}

	// 确保健康度不低于0
	if health < 0 {
		health = 0
	}

	return health
}
