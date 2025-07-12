/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : gcp云提供商
 */

package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	"github.com/karmada-io/karmada/pkg/cloud/constants"
)

type GcpProviderConfig struct {
	ProjectID   string
	AccessToken string // OAuth 2.0 access token or API key
}

type GcpProvider struct {
	Provider
	name      string
	locations []v1alpha1.CloudLocation
	config    *GcpProviderConfig
	client    *http.Client
}

// GCP API 响应结构
type GcpRegionsResponse struct {
	Items []struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Status      string   `json:"status"`
		Zones       []string `json:"zones"`
	} `json:"items"`
}

type GcpZonesResponse struct {
	Items []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Region      string `json:"region"`
	} `json:"items"`
}

type GcpInstancesResponse struct {
	Items map[string]struct {
		Instances []struct {
			Id          string `json:"id"`
			Name        string `json:"name"`
			MachineType string `json:"machineType"`
			Status      string `json:"status"`
			Zone        string `json:"zone"`
		} `json:"instances"`
	} `json:"items"`
}

type GcpMachineTypesResponse struct {
	Items []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		GuestCpus   int    `json:"guestCpus"`
		MemoryMb    int    `json:"memoryMb"`
		Zone        string `json:"zone"`
	} `json:"items"`
}

// 简化的价格响应结构（实际GCP计费API更复杂）
type GcpPriceInfo struct {
	MachineType  string  `json:"machineType"`
	PricePerHour float64 `json:"pricePerHour"`
	Region       string  `json:"region"`
}

func NewGcpProvider(config *GcpProviderConfig) *GcpProvider {
	return &GcpProvider{
		name:      string(constants.CloudProviderGcp),
		locations: []v1alpha1.CloudLocation{},
		config:    config,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *GcpProvider) GetName() string {
	return p.name
}

// 调用 GCP API
func (p *GcpProvider) callAPI(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// 设置认证头
	req.Header.Set("Authorization", "Bearer "+p.config.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API call failed with status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (p *GcpProvider) GetLocations() ([]v1alpha1.CloudLocation, error) {
	// 获取所有区域
	regionsURL := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/regions", p.config.ProjectID)
	regionsData, err := p.callAPI(regionsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get regions: %v", err)
	}

	var regionsResp GcpRegionsResponse
	if err := json.Unmarshal(regionsData, &regionsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal regions response: %v", err)
	}

	// 获取所有可用区
	zonesURL := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones", p.config.ProjectID)
	zonesData, err := p.callAPI(zonesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get zones: %v", err)
	}

	var zonesResp GcpZonesResponse
	if err := json.Unmarshal(zonesData, &zonesResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal zones response: %v", err)
	}

	// 按区域分组可用区
	regionZones := make(map[string][]string)
	for _, zone := range zonesResp.Items {
		if zone.Status == "UP" {
			// 从zone的region URL中提取region名称
			regionName := extractRegionFromURL(zone.Region)
			regionZones[regionName] = append(regionZones[regionName], zone.Name)
		}
	}

	var locations []v1alpha1.CloudLocation
	for _, region := range regionsResp.Items {
		if region.Status == "UP" {
			zones := regionZones[region.Name]
			location := v1alpha1.CloudLocation{
				Region: v1alpha1.CloudRegion(region.Name),
				Zones:  zones,
			}
			locations = append(locations, location)
		}
	}

	p.locations = locations
	return locations, nil
}

func (p *GcpProvider) GetInstanceList() ([]string, error) {
	// 获取所有实例（聚合接口）
	instancesURL := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/aggregated/instances", p.config.ProjectID)
	data, err := p.callAPI(instancesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances: %v", err)
	}

	var instancesResp GcpInstancesResponse
	if err := json.Unmarshal(data, &instancesResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal instances response: %v", err)
	}

	var allInstances []string
	for _, zoneData := range instancesResp.Items {
		for _, instance := range zoneData.Instances {
			allInstances = append(allInstances, instance.Name)
		}
	}

	return allInstances, nil
}

func (p *GcpProvider) GetInstancePrice(instanceType string) (float64, error) {
	// 由于GCP的计费API比较复杂，这里提供一个简化版本
	// 实际应该使用Cloud Billing API或者预定义的价格表

	// 获取机器类型信息（用于计算基础价格）
	machineTypesURL := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/aggregated/machineTypes", p.config.ProjectID)
	data, err := p.callAPI(machineTypesURL)
	if err != nil {
		return 0, fmt.Errorf("failed to get machine types: %v", err)
	}

	var machineTypesResp map[string]struct {
		MachineTypes []struct {
			Name      string `json:"name"`
			GuestCpus int    `json:"guestCpus"`
			MemoryMb  int    `json:"memoryMb"`
		} `json:"machineTypes"`
	}

	if err := json.Unmarshal(data, &machineTypesResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal machine types response: %v", err)
	}

	// 查找指定的机器类型
	for _, zoneData := range machineTypesResp {
		for _, machineType := range zoneData.MachineTypes {
			if machineType.Name == instanceType {
				// 简化的价格计算：基于CPU和内存
				// 实际价格应该从GCP官方价格API获取
				cpuPrice := float64(machineType.GuestCpus) * 0.0475          // 每CPU每小时约$0.0475
				memoryPrice := float64(machineType.MemoryMb) / 1024 * 0.0063 // 每GB内存每小时约$0.0063
				totalPrice := cpuPrice + memoryPrice
				return totalPrice, nil
			}
		}
	}

	return 0, fmt.Errorf("machine type %s not found", instanceType)
}

// 辅助函数：从URL中提取区域名称
func extractRegionFromURL(regionURL string) string {
	parts := strings.Split(regionURL, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

var _ Provider = &GcpProvider{}
