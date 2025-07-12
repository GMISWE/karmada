/*
 * @Version : 1.0
 * @Author  : xiaokang.w
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/07/12
 * @Desc    : aliyun云提供商
 */

package providers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	"github.com/karmada-io/karmada/pkg/cloud/constants"
)

type AliyunProviderConfig struct {
	AccessKey string
	SecretKey string
}

type AliyunProvider struct {
	Provider
	name      string
	locations []v1alpha1.CloudLocation
	config    *AliyunProviderConfig
	client    *http.Client
}

// 阿里云API响应结构
type AliyunRegionsResponse struct {
	Regions struct {
		Region []struct {
			RegionId   string `json:"RegionId"`
			RegionName string `json:"RegionName"`
		} `json:"Region"`
	} `json:"Regions"`
}

type AliyunZonesResponse struct {
	Zones struct {
		Zone []struct {
			ZoneId   string `json:"ZoneId"`
			ZoneName string `json:"ZoneName"`
		} `json:"Zone"`
	} `json:"Zones"`
}

type AliyunInstancesResponse struct {
	Instances struct {
		Instance []struct {
			InstanceId   string `json:"InstanceId"`
			InstanceName string `json:"InstanceName"`
			InstanceType string `json:"InstanceType"`
		} `json:"Instance"`
	} `json:"Instances"`
}

type AliyunPriceResponse struct {
	PriceInfo struct {
		Price struct {
			OriginalPrice float64 `json:"OriginalPrice"`
			DiscountPrice float64 `json:"DiscountPrice"`
		} `json:"Price"`
	} `json:"PriceInfo"`
}

func NewAliyunProvider(config *AliyunProviderConfig) *AliyunProvider {
	return &AliyunProvider{
		name:      string(constants.CloudProviderAliyun),
		locations: []v1alpha1.CloudLocation{},
		config:    config,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *AliyunProvider) GetName() string {
	return p.name
}

// 生成阿里云API签名
func (p *AliyunProvider) generateSignature(method, canonicalizedQuery string) string {
	stringToSign := method + "&" + url.QueryEscape("/") + "&" + url.QueryEscape(canonicalizedQuery)
	h := hmac.New(sha1.New, []byte(p.config.SecretKey+"&"))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 调用阿里云API
func (p *AliyunProvider) callAPI(action, regionId string, params map[string]string) ([]byte, error) {
	// 基本参数
	baseParams := map[string]string{
		"Action":           action,
		"Version":          "2014-05-26",
		"AccessKeyId":      p.config.AccessKey,
		"SignatureMethod":  "HMAC-SHA1",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"SignatureVersion": "1.0",
		"SignatureNonce":   fmt.Sprintf("%d", time.Now().UnixNano()),
		"Format":           "JSON",
	}

	// 合并参数
	allParams := make(map[string]string)
	for k, v := range baseParams {
		allParams[k] = v
	}
	for k, v := range params {
		allParams[k] = v
	}

	// 构建查询字符串
	var keys []string
	for k := range allParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var query []string
	for _, k := range keys {
		query = append(query, url.QueryEscape(k)+"="+url.QueryEscape(allParams[k]))
	}
	canonicalizedQuery := strings.Join(query, "&")

	// 生成签名
	signature := p.generateSignature("GET", canonicalizedQuery)

	// 构建最终URL
	endpoint := fmt.Sprintf("https://ecs.%s.aliyuncs.com", regionId)
	finalURL := endpoint + "/?" + canonicalizedQuery + "&Signature=" + url.QueryEscape(signature)

	// 发送请求
	resp, err := p.client.Get(finalURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (p *AliyunProvider) GetLocations() ([]v1alpha1.CloudLocation, error) {
	// 获取所有区域
	data, err := p.callAPI("DescribeRegions", "cn-hangzhou", map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe regions: %v", err)
	}

	var regionsResp AliyunRegionsResponse
	if err := json.Unmarshal(data, &regionsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal regions response: %v", err)
	}

	var locations []v1alpha1.CloudLocation

	for _, region := range regionsResp.Regions.Region {
		// 为每个区域获取可用区
		zonesData, err := p.callAPI("DescribeZones", region.RegionId, map[string]string{
			"RegionId": region.RegionId,
		})
		if err != nil {
			continue // 跳过出错的区域
		}

		var zonesResp AliyunZonesResponse
		if err := json.Unmarshal(zonesData, &zonesResp); err != nil {
			continue
		}

		var zones []string
		for _, zone := range zonesResp.Zones.Zone {
			zones = append(zones, zone.ZoneId)
		}

		location := v1alpha1.CloudLocation{
			Region: v1alpha1.CloudRegion(region.RegionId),
			Zones:  zones,
		}
		locations = append(locations, location)
	}

	p.locations = locations
	return locations, nil
}

func (p *AliyunProvider) GetInstanceList() ([]string, error) {
	var allInstances []string

	// 遍历所有区域获取实例
	for _, location := range p.locations {
		pageNumber := 1
		pageSize := 100

		for {
			params := map[string]string{
				"RegionId":   string(location.Region),
				"PageNumber": strconv.Itoa(pageNumber),
				"PageSize":   strconv.Itoa(pageSize),
			}

			data, err := p.callAPI("DescribeInstances", string(location.Region), params)
			if err != nil {
				break
			}

			var instancesResp AliyunInstancesResponse
			if err := json.Unmarshal(data, &instancesResp); err != nil {
				break
			}

			for _, instance := range instancesResp.Instances.Instance {
				allInstances = append(allInstances, instance.InstanceId)
			}

			if len(instancesResp.Instances.Instance) < pageSize {
				break
			}
			pageNumber++
		}
	}

	return allInstances, nil
}

func (p *AliyunProvider) GetInstancePrice(instanceType string) (float64, error) {
	params := map[string]string{
		"ResourceType": "instance",
		"InstanceType": instanceType,
		"RegionId":     "cn-hangzhou", // 默认使用杭州区域查询价格
		"Period":       "1",
		"PriceUnit":    "Hour",
	}

	data, err := p.callAPI("DescribePrice", "cn-hangzhou", params)
	if err != nil {
		return 0, fmt.Errorf("failed to get instance price: %v", err)
	}

	var priceResp AliyunPriceResponse
	if err := json.Unmarshal(data, &priceResp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal price response: %v", err)
	}

	return priceResp.PriceInfo.Price.OriginalPrice, nil
}

var _ Provider = &AliyunProvider{}
