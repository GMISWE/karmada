/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/06/02
 * @Desc    : topo types
 */

package types

import (
	"time"
)

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

type GpuType string

const (
	GpuTypeA100 GpuType = "A100"
	GpuTypeH100 GpuType = "H100"
	GpuTypeV100 GpuType = "V100"
	GpuTypeT4   GpuType = "T4"
)

type TopoGpu struct {
	ID   string `json:"id"`
	Node string `json:"node"`
}

// NodeTaint 节点污点
type NodeTaint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"` // NoSchedule, PreferNoSchedule, NoExecute
}

type NodeLabels map[string]string

type Topo struct {
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
	// Resources ResourceMetrics `json:"resources"`
	RawNode Node `json:"raw_node"` // 原始节点数据

	// 调度相关
	Taints      []NodeTaint       `json:"taints"`      // 污点
	Annotations map[string]string `json:"annotations"` // 注解
	Priority    int               `json:"priority"`    // 优先级
}
