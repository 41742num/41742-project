package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/project47/cmd/mywebapp/types"
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	BaseURL     string        `json:"base_url"`
	APIKey      string        `json:"api_key"`
	Timeout     time.Duration `json:"timeout"`
	CacheTTL    time.Duration `json:"cache_ttl"`
	MaxRetries  int           `json:"max_retries"`
	RetryDelay  time.Duration `json:"retry_delay"`
}

// MiddlewareClient 中间件API客户端
type MiddlewareClient struct {
	config     *MiddlewareConfig
	httpClient *http.Client
	cache      *Cache
}

// Cache 内存缓存
type Cache struct {
	devices        []types.Device
	devicesExpiry  time.Time
	deviceStatuses map[string]types.DeviceStatus
	statusExpiry   time.Time
}

// NewMiddlewareClient 创建新的中间件客户端
func NewMiddlewareClient(config *MiddlewareConfig) *MiddlewareClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 1 * time.Second
	}

	return &MiddlewareClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		cache: &Cache{
			deviceStatuses: make(map[string]types.DeviceStatus),
		},
	}
}

// DefaultConfig 默认配置
func DefaultConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		BaseURL:    "http://localhost:8080", // 俄罗斯中间件默认地址
		Timeout:    30 * time.Second,
		CacheTTL:   30 * time.Second,
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
	}
}

// ==================== 设备相关API ====================

// MiddlewareDevice 中间件返回的设备结构
type MiddlewareDevice struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	IP           string                 `json:"ip"`
	Port         int                    `json:"port"`
	Status       string                 `json:"status"` // online/offline/maintenance
	IsOnline     bool                   `json:"is_online"`
	HasFault     bool                   `json:"has_fault"`
	ProcessInfo  map[string]interface{} `json:"process_info"`
	ReagentInfo  []MiddlewareReagent    `json:"reagent_info"`
	LastCheck    string                 `json:"last_check"`
	Uptime       string                 `json:"uptime"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

// MiddlewareReagent 中间件返回的试剂结构
type MiddlewareReagent struct {
	Name     string  `json:"name"`
	Current  float64 `json:"current"`
	Capacity float64 `json:"capacity"`
	Unit     string  `json:"unit"`
	Percent  float64 `json:"percent"`
	Warning  bool    `json:"warning"`
	Critical bool    `json:"critical"`
}

// GetDevices 从中间件获取设备列表
func (c *MiddlewareClient) GetDevices() ([]types.Device, error) {
	// 检查缓存
	if !c.cache.devicesExpiry.IsZero() && time.Now().Before(c.cache.devicesExpiry) {
		return c.cache.devices, nil
	}

	// 调用API
	url := fmt.Sprintf("%s/api/devices", c.config.BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	if c.config.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	}

	// 重试逻辑
	var resp *http.Response
	for i := 0; i < c.config.MaxRetries; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		if i < c.config.MaxRetries-1 {
			time.Sleep(c.config.RetryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("API调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误: %s - %s", resp.Status, string(body))
	}

	// 解析响应
	var middlewareDevices []MiddlewareDevice
	if err := json.NewDecoder(resp.Body).Decode(&middlewareDevices); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 转换为模型设备
	devices := make([]types.Device, len(middlewareDevices))
	for i, md := range middlewareDevices {
		devices[i] = c.convertToModelDevice(md)
	}

	// 更新缓存
	c.cache.devices = devices
	c.cache.devicesExpiry = time.Now().Add(c.config.CacheTTL)

	return devices, nil
}

// GetDeviceStatus 从中间件获取设备状态
func (c *MiddlewareClient) GetDeviceStatus(deviceID string) (types.DeviceStatus, error) {
	// 检查缓存
	if status, ok := c.cache.deviceStatuses[deviceID]; ok {
		if !c.cache.statusExpiry.IsZero() && time.Now().Before(c.cache.statusExpiry) {
			return status, nil
		}
	}

	// 调用API
	url := fmt.Sprintf("%s/api/devices/%s/status", c.config.BaseURL, deviceID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.DeviceStatus{}, fmt.Errorf("创建请求失败: %v", err)
	}

	if c.config.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return types.DeviceStatus{}, fmt.Errorf("API调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return types.DeviceStatus{}, fmt.Errorf("设备不存在: %s", deviceID)
		}
		body, _ := io.ReadAll(resp.Body)
		return types.DeviceStatus{}, fmt.Errorf("API返回错误: %s - %s", resp.Status, string(body))
	}

	// 解析响应
	var middlewareStatus map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&middlewareStatus); err != nil {
		return types.DeviceStatus{}, fmt.Errorf("解析响应失败: %v", err)
	}

	// 转换为模型状态
	status := c.convertToModelStatus(deviceID, middlewareStatus)

	// 更新缓存
	c.cache.deviceStatuses[deviceID] = status
	c.cache.statusExpiry = time.Now().Add(c.config.CacheTTL)

	return status, nil
}

// GetDeviceRealtimeData 获取设备实时数据
func (c *MiddlewareClient) GetDeviceRealtimeData(deviceID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/devices/%s/realtime", c.config.BaseURL, deviceID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	if c.config.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误: %s - %s", resp.Status, string(body))
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return data, nil
}

// RestartDevice 重启设备
func (c *MiddlewareClient) RestartDevice(deviceID string) error {
	url := fmt.Sprintf("%s/api/devices/%s/restart", c.config.BaseURL, deviceID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	if c.config.APIKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("API调用失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API返回错误: %s - %s", resp.Status, string(body))
	}

	// 清除缓存
	delete(c.cache.deviceStatuses, deviceID)

	return nil
}

// ==================== 转换函数 ====================

func (c *MiddlewareClient) convertToModelDevice(md MiddlewareDevice) types.Device {
	// 解析最后检查时间
	var lastCheck time.Time
	if md.LastCheck != "" {
		lastCheck, _ = time.Parse(time.RFC3339, md.LastCheck)
	} else {
		lastCheck = time.Now()
	}

	// 转换试剂信息
	reagents := make([]types.Reagent, len(md.ReagentInfo))
	for i, mr := range md.ReagentInfo {
		reagents[i] = types.Reagent{
			Name:     mr.Name,
			Current:  mr.Current,
			Capacity: mr.Capacity,
			Unit:     mr.Unit,
			Percent:  mr.Percent,
		}
	}

	// 确定设备状态
	status := "enabled"
	if md.Status == "maintenance" || md.Status == "offline" {
		status = "disabled"
	}

	return types.Device{
		ID:          md.ID,
		DeviceID:    md.ID,
		Name:        md.Name,
		Status:      status,
		IP:          md.IP,
		Port:        md.Port,
		ProcessName: c.extractProcessName(md.ProcessInfo),
		Reagents:    reagents,
		LastCheck:   lastCheck,
		IsOnline:    md.IsOnline,
		HasFault:    md.HasFault,
	}
}

func (c *MiddlewareClient) convertToModelStatus(deviceID string, data map[string]interface{}) types.DeviceStatus {
	// 提取基本信息
	name, _ := data["name"].(string)
	status, _ := data["status"].(string)
	isOnline, _ := data["is_online"].(bool)
	hasFault, _ := data["has_fault"].(bool)
	lastCheck, _ := data["last_check"].(string)
	uptime, _ := data["uptime"].(string)

	// 转换试剂信息
	var reagents []types.Reagent
	if reagentInfo, ok := data["reagent_info"].([]interface{}); ok {
		for _, ri := range reagentInfo {
			if reagentMap, ok := ri.(map[string]interface{}); ok {
				reagent := types.Reagent{
					Name:     reagentMap["name"].(string),
					Current:  reagentMap["current"].(float64),
					Capacity: reagentMap["capacity"].(float64),
					Unit:     reagentMap["unit"].(string),
					Percent:  reagentMap["percent"].(float64),
				}
				reagents = append(reagents, reagent)
			}
		}
	}

	// 计算试剂状态
	reagentStatus := "normal"
	if len(reagents) > 0 {
		minPercent := 100.0
		for _, reagent := range reagents {
			if reagent.Percent < minPercent {
				minPercent = reagent.Percent
			}
		}

		if minPercent == 0 {
			reagentStatus = "empty"
		} else if minPercent < 30 {
			reagentStatus = "low"
		} else if minPercent < 70 {
			reagentStatus = "warning"
		}
	}

	return types.DeviceStatus{
		DeviceID:      deviceID,
		DeviceName:    name,
		Status:        status,
		IsOnline:      isOnline,
		HasFault:      hasFault,
		ReagentStatus: reagentStatus,
		Reagents:      reagents,
		LastCheck:     lastCheck,
		Uptime:        uptime,
	}
}

func (c *MiddlewareClient) extractProcessName(processInfo map[string]interface{}) string {
	if name, ok := processInfo["name"].(string); ok {
		return name
	}
	if cmd, ok := processInfo["command"].(string); ok {
		return cmd
	}
	return "unknown"
}

// ==================== 缓存管理 ====================

// ClearCache 清除缓存
func (c *MiddlewareClient) ClearCache() {
	c.cache.devices = nil
	c.cache.devicesExpiry = time.Time{}
	c.cache.deviceStatuses = make(map[string]types.DeviceStatus)
	c.cache.statusExpiry = time.Time{}
}

// GetCacheInfo 获取缓存信息
func (c *MiddlewareClient) GetCacheInfo() map[string]interface{} {
	devicesCount := len(c.cache.devices)
	statusCount := len(c.cache.deviceStatuses)

	devicesExpired := false
	if !c.cache.devicesExpiry.IsZero() {
		devicesExpired = time.Now().After(c.cache.devicesExpiry)
	}

	statusExpired := false
	if !c.cache.statusExpiry.IsZero() {
		statusExpired = time.Now().After(c.cache.statusExpiry)
	}

	return map[string]interface{}{
		"devices_count":      devicesCount,
		"status_count":       statusCount,
		"devices_expired":    devicesExpired,
		"status_expired":     statusExpired,
		"devices_expiry":     c.cache.devicesExpiry.Format(time.RFC3339),
		"status_expiry":      c.cache.statusExpiry.Format(time.RFC3339),
		"cache_ttl_seconds":  c.config.CacheTTL.Seconds(),
	}
}