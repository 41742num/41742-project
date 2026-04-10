package data

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/project47/cmd/mywebapp/database"
	"github.com/project47/cmd/mywebapp/middleware"
	"github.com/project47/cmd/mywebapp/types"
)

// DataSource 数据源类型
type DataSource string

const (
	// SourceMiddleware 从俄罗斯中间件获取数据
	SourceMiddleware DataSource = "middleware"
	// SourceDatabase 从数据库获取历史数据
	SourceDatabase DataSource = "database"
	// SourceMock 使用模拟数据
	SourceMock DataSource = "mock"
	// SourceFallback 回退到本地数据
	SourceFallback DataSource = "fallback"
)

// DataManager 数据管理器
type DataManager struct {
	mutex          sync.RWMutex
	config         *Config
	middleware     *middleware.MiddlewareClient
	dbManager      *database.DatabaseManager
	currentSource  DataSource
	devices        []types.Device
	deviceStatuses map[string]types.DeviceStatus
	lastUpdate     time.Time
	updateInterval time.Duration
	stopChan       chan struct{}
}

// Config 数据管理器配置
type Config struct {
	// 数据源配置
	DataSource      DataSource `json:"data_source"`
	MiddlewareURL   string     `json:"middleware_url"`
	MiddlewareAPIKey string    `json:"middleware_api_key"`

	// 数据库配置
	DatabaseType    string        `json:"database_type"`
	SQLitePath      string        `json:"sqlite_path"`
	EnableDatabase  bool          `json:"enable_database"`

	// 缓存配置
	CacheTTL        time.Duration `json:"cache_ttl"`
	UpdateInterval  time.Duration `json:"update_interval"`

	// 重试配置
	MaxRetries      int           `json:"max_retries"`
	RetryDelay      time.Duration `json:"retry_delay"`

	// 回退配置
	EnableFallback  bool          `json:"enable_fallback"`
	FallbackTimeout time.Duration `json:"fallback_timeout"`

	// 历史数据保留天数
	DataRetentionDays int `json:"data_retention_days"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		DataSource:        SourceMiddleware,
		MiddlewareURL:     "http://localhost:8080",
		DatabaseType:      "sqlite",
		SQLitePath:        "./data/project47.db",
		EnableDatabase:    true,
		CacheTTL:          30 * time.Second,
		UpdateInterval:    30 * time.Second,
		MaxRetries:        3,
		RetryDelay:        1 * time.Second,
		EnableFallback:    true,
		FallbackTimeout:   5 * time.Second,
		DataRetentionDays: 30,
	}
}

// NewDataManager 创建新的数据管理器
func NewDataManager(config *Config) (*DataManager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	dm := &DataManager{
		config:         config,
		currentSource:  config.DataSource,
		deviceStatuses: make(map[string]types.DeviceStatus),
		updateInterval: config.UpdateInterval,
		stopChan:       make(chan struct{}),
	}

	// 初始化中间件客户端
	if config.DataSource == SourceMiddleware || config.EnableFallback {
		middlewareConfig := &middleware.MiddlewareConfig{
			BaseURL:    config.MiddlewareURL,
			APIKey:     config.MiddlewareAPIKey,
			Timeout:    10 * time.Second,
			CacheTTL:   config.CacheTTL,
			MaxRetries: config.MaxRetries,
			RetryDelay: config.RetryDelay,
		}
		dm.middleware = middleware.NewMiddlewareClient(middlewareConfig)
	}

	// 初始化数据库管理器
	if config.EnableDatabase {
		dbConfig := &database.DBConfig{
			Type:             config.DatabaseType,
			SQLitePath:       config.SQLitePath,
			AutoMigrate:      true,
			EnableQueryLog:   false,
		}
		dbManager, err := database.NewDatabaseManager(dbConfig)
		if err != nil {
			log.Printf("警告: 初始化数据库管理器失败: %v", err)
			log.Println("将禁用数据库功能")
		} else {
			dm.dbManager = dbManager
		}
	}

	// 初始化设备数据
	if err := dm.initializeDevices(); err != nil {
		return nil, fmt.Errorf("初始化设备数据失败: %v", err)
	}

	// 启动后台更新任务
	go dm.startBackgroundUpdate()

	// 启动数据库清理任务
	if dm.dbManager != nil {
		go dm.startDatabaseCleanup()
	}

	return dm, nil
}

// initializeDevices 初始化设备数据
func (dm *DataManager) initializeDevices() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	switch dm.currentSource {
	case SourceMiddleware:
		return dm.loadFromMiddleware()
	case SourceDatabase:
		return dm.loadFromDatabase()
	case SourceMock:
		return dm.loadFromMock()
	case SourceFallback:
		return dm.loadFromFallback()
	default:
		return dm.loadFromFallback()
	}
}

// loadFromMiddleware 从中间件加载数据
func (dm *DataManager) loadFromMiddleware() error {
	if dm.middleware == nil {
		return fmt.Errorf("中间件客户端未初始化")
	}

	devices, err := dm.middleware.GetDevices()
	if err != nil {
		if dm.config.EnableFallback {
			dm.currentSource = SourceFallback
			return dm.loadFromFallback()
		}
		return fmt.Errorf("从中间件加载设备失败: %v", err)
	}

	dm.devices = devices
	dm.lastUpdate = time.Now()
	return nil
}

// loadFromDatabase 从数据库加载数据
func (dm *DataManager) loadFromDatabase() error {
	if dm.dbManager == nil {
		if dm.config.EnableFallback {
			dm.currentSource = SourceFallback
			return dm.loadFromFallback()
		}
		return fmt.Errorf("数据库管理器未初始化")
	}

	devices, err := dm.dbManager.GetDevicesFromDatabase()
	if err != nil {
		if dm.config.EnableFallback {
			dm.currentSource = SourceFallback
			return dm.loadFromFallback()
		}
		return fmt.Errorf("从数据库加载设备失败: %v", err)
	}

	dm.devices = devices
	dm.lastUpdate = time.Now()

	// 记录数据源切换
	// 注意：这里应该通过dbManager的store来记录，但为了简化先跳过

	return nil
}

// loadFromMock 加载模拟数据
func (dm *DataManager) loadFromMock() error {
	// 这里可以调用现有的mock生成器
	// 暂时使用回退数据
	return dm.loadFromFallback()
}

// loadFromFallback 加载回退数据（本地硬编码数据）
func (dm *DataManager) loadFromFallback() error {
	// 创建一些默认的回退数据
	dm.devices = []types.Device{
		{
			ID:          "DEV001",
			DeviceID:    "MIDDLEWARE_001",
			Name:        "俄罗斯中间件服务器",
			Status:      "enabled",
			IP:          "172.19.14.202",
			Port:        10001,
			ProcessName: "middleware",
			Reagents: []types.Reagent{
				{
					Name:     "试剂A",
					Current:  85.0,
					Capacity: 100.0,
					Unit:     "ml",
					Percent:  85.0,
				},
				{
					Name:     "试剂B",
					Current:  45.0,
					Capacity: 100.0,
					Unit:     "ml",
					Percent:  45.0,
				},
			},
			LastCheck: time.Now(),
			IsOnline:  true,
			HasFault:  false,
		},
		{
			ID:          "DEV002",
			DeviceID:    "ANALYZER_001",
			Name:        "分析仪器1",
			Status:      "enabled",
			IP:          "172.19.14.203",
			Port:        10002,
			ProcessName: "analyzer",
			Reagents: []types.Reagent{
				{
					Name:     "标准液",
					Current:  70.0,
					Capacity: 100.0,
					Unit:     "ml",
					Percent:  70.0,
				},
			},
			LastCheck: time.Now(),
			IsOnline:  true,
			HasFault:  false,
		},
		{
			ID:          "DEV003",
			DeviceID:    "ANALYZER_002",
			Name:        "分析仪器2",
			Status:      "disabled",
			IP:          "172.19.14.204",
			Port:        10003,
			ProcessName: "analyzer",
			Reagents: []types.Reagent{
				{
					Name:     "标准液",
					Current:  0.0,
					Capacity: 100.0,
					Unit:     "ml",
					Percent:  0.0,
				},
			},
			LastCheck: time.Now(),
			IsOnline:  false,
			HasFault:  false,
		},
	}
	dm.lastUpdate = time.Now()
	return nil
}

// startBackgroundUpdate 启动后台更新任务
func (dm *DataManager) startBackgroundUpdate() {
	ticker := time.NewTicker(dm.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dm.updateDevices()
		case <-dm.stopChan:
			return
		}
	}
}

// updateDevices 更新设备数据
func (dm *DataManager) updateDevices() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// 检查是否需要更新
	if time.Since(dm.lastUpdate) < dm.config.CacheTTL {
		return
	}

	switch dm.currentSource {
	case SourceMiddleware:
		if devices, err := dm.middleware.GetDevices(); err == nil {
			dm.devices = devices
			dm.lastUpdate = time.Now()
			// 清除状态缓存，因为设备列表可能已更新
			dm.deviceStatuses = make(map[string]types.DeviceStatus)

			// 同步到数据库
			if dm.dbManager != nil {
				go func() {
					if err := dm.dbManager.SyncDevicesFromMiddleware(devices); err != nil {
						log.Printf("同步到数据库失败: %v", err)
					}
				}()
			}
		} else if dm.config.EnableFallback {
			dm.currentSource = SourceFallback
			dm.loadFromFallback()
		}
	case SourceDatabase:
		// 从数据库更新
		if dm.dbManager != nil {
			if devices, err := dm.dbManager.GetDevicesFromDatabase(); err == nil {
				dm.devices = devices
				dm.lastUpdate = time.Now()
				dm.deviceStatuses = make(map[string]types.DeviceStatus)
			} else if dm.config.EnableFallback {
				dm.currentSource = SourceFallback
				dm.loadFromFallback()
			}
		}
	case SourceMock, SourceFallback:
		// 对于模拟和回退数据，只需更新最后更新时间
		dm.lastUpdate = time.Now()
	}
}

// ==================== 公共接口 ====================

// GetDevices 获取设备列表
func (dm *DataManager) GetDevices() []types.Device {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	// 返回副本，避免外部修改
	devices := make([]types.Device, len(dm.devices))
	copy(devices, dm.devices)
	return devices
}

// GetDeviceByID 根据ID获取设备
func (dm *DataManager) GetDeviceByID(deviceID string) *types.Device {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	for _, device := range dm.devices {
		if device.DeviceID == deviceID {
			// 返回指针副本
			deviceCopy := device
			return &deviceCopy
		}
	}
	return nil
}

// GetDeviceStatus 获取设备状态
func (dm *DataManager) GetDeviceStatus(deviceID string) (types.DeviceStatus, error) {
	// 首先检查缓存
	dm.mutex.RLock()
	if status, ok := dm.deviceStatuses[deviceID]; ok {
		dm.mutex.RUnlock()
		return status, nil
	}
	dm.mutex.RUnlock()

	// 获取设备
	device := dm.GetDeviceByID(deviceID)
	if device == nil {
		return types.DeviceStatus{}, fmt.Errorf("设备不存在: %s", deviceID)
	}

	var status types.DeviceStatus
	var err error

	switch dm.currentSource {
	case SourceMiddleware:
		status, err = dm.middleware.GetDeviceStatus(deviceID)
		if err != nil && dm.config.EnableFallback {
			// 回退到本地计算
			status = dm.calculateDeviceStatus(*device)
			err = nil
		}
	case SourceDatabase:
		// 从数据库获取最新状态
		if dm.dbManager != nil {
			// 尝试从数据库获取设备
			if dbDevice, err := dm.dbManager.GetDeviceFromDatabase(deviceID); err == nil && dbDevice != nil {
				status = dm.calculateDeviceStatus(*dbDevice)
			} else {
				// 回退到当前设备数据
				status = dm.calculateDeviceStatus(*device)
			}
		} else {
			status = dm.calculateDeviceStatus(*device)
		}
	case SourceMock, SourceFallback:
		status = dm.calculateDeviceStatus(*device)
	default:
		status = dm.calculateDeviceStatus(*device)
	}

	if err == nil {
		dm.mutex.Lock()
		dm.deviceStatuses[deviceID] = status
		dm.mutex.Unlock()
	}

	return status, err
}

// GetAllDevicesStatus 获取所有设备状态
func (dm *DataManager) GetAllDevicesStatus() []types.DeviceStatus {
	dm.mutex.RLock()
	devices := dm.devices
	dm.mutex.RUnlock()

	var statuses []types.DeviceStatus
	for _, device := range devices {
		status, err := dm.GetDeviceStatus(device.DeviceID)
		if err == nil {
			statuses = append(statuses, status)
		}
	}

	return statuses
}

// GetDeviceStats 获取设备统计信息
func (dm *DataManager) GetDeviceStats() map[string]interface{} {
	dm.mutex.RLock()
	devices := dm.devices
	dm.mutex.RUnlock()

	total := len(devices)
	enabled := 0
	online := 0
	withFault := 0
	lowReagent := 0
	emptyReagent := 0

	for _, device := range devices {
		if device.Status == "enabled" {
			enabled++
		}
		if device.IsOnline {
			online++
		}
		if device.HasFault {
			withFault++
		}

		// 检查试剂状态
		if device.Status == "enabled" {
			minPercent := 100.0
			for _, reagent := range device.Reagents {
				if reagent.Percent < minPercent {
					minPercent = reagent.Percent
				}
			}

			if minPercent == 0 {
				emptyReagent++
			} else if minPercent < 30 {
				lowReagent++
			}
		}
	}

	return map[string]interface{}{
		"total_devices":         total,
		"enabled_devices":       enabled,
		"disabled_devices":      total - enabled,
		"online_devices":        online,
		"offline_devices":       total - online,
		"devices_with_fault":    withFault,
		"devices_low_reagent":   lowReagent,
		"devices_empty_reagent": emptyReagent,
		"data_source":           string(dm.currentSource),
		"last_update":           dm.lastUpdate.Format(time.RFC3339),
		"update_interval":       dm.updateInterval.Seconds(),
	}
}

// RestartDevice 重启设备
func (dm *DataManager) RestartDevice(deviceID string) error {
	if dm.currentSource == SourceMiddleware && dm.middleware != nil {
		err := dm.middleware.RestartDevice(deviceID)
		if err == nil {
			// 清除该设备的状态缓存
			dm.mutex.Lock()
			delete(dm.deviceStatuses, deviceID)
			dm.mutex.Unlock()
		}
		return err
	}

	// 对于模拟和回退数据，模拟重启
	device := dm.GetDeviceByID(deviceID)
	if device == nil {
		return fmt.Errorf("设备不存在: %s", deviceID)
	}

	if device.Status != "enabled" {
		return fmt.Errorf("设备已禁用，无法重启")
	}

	// 模拟重启：短暂离线后恢复
	dm.mutex.Lock()
	for idx, d := range dm.devices {
		if d.DeviceID == deviceID {
			dm.devices[idx].IsOnline = false
			dm.devices[idx].HasFault = false
			dm.devices[idx].LastCheck = time.Now()

			// 重置试剂余量
			for j := range dm.devices[idx].Reagents {
				dm.devices[idx].Reagents[j].Current = dm.devices[idx].Reagents[j].Capacity
				dm.devices[idx].Reagents[j].Percent = 100.0
			}
			break
		}
	}
	dm.mutex.Unlock()

	// 清除状态缓存
	dm.mutex.Lock()
	delete(dm.deviceStatuses, deviceID)
	dm.mutex.Unlock()

	return nil
}

// calculateDeviceStatus 计算设备状态
func (dm *DataManager) calculateDeviceStatus(device types.Device) types.DeviceStatus {
	// 检查试剂状态（取最低的试剂百分比）
	minPercent := 100.0
	for _, reagent := range device.Reagents {
		if reagent.Percent < minPercent {
			minPercent = reagent.Percent
		}
	}

	reagentStatus := "normal"
	if minPercent == 0 {
		reagentStatus = "empty"
	} else if minPercent < 30 {
		reagentStatus = "low"
	} else if minPercent < 70 {
		reagentStatus = "warning"
	}

	return types.DeviceStatus{
		DeviceID:      device.DeviceID,
		DeviceName:    device.Name,
		Status:        device.Status,
		IsOnline:      device.IsOnline,
		HasFault:      device.HasFault,
		ReagentStatus: reagentStatus,
		Reagents:      device.Reagents,
		LastCheck:     device.LastCheck.Format(time.RFC3339),
		Uptime:        "24h", // 模拟数据，实际应从设备获取
	}
}

// UpdateDevice 更新设备信息
func (dm *DataManager) UpdateDevice(deviceID string, updates map[string]interface{}) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	for i, device := range dm.devices {
		if device.DeviceID == deviceID {
			// 更新设备信息
			if name, ok := updates["name"].(string); ok && name != "" {
				dm.devices[i].Name = name
			}
			if status, ok := updates["status"].(string); ok && (status == "enabled" || status == "disabled") {
				dm.devices[i].Status = status
			}
			if ip, ok := updates["ip"].(string); ok && ip != "" {
				dm.devices[i].IP = ip
			}
			if port, ok := updates["port"].(float64); ok && port > 0 {
				dm.devices[i].Port = int(port)
			}

			dm.devices[i].LastCheck = time.Now()

			// 清除状态缓存
			delete(dm.deviceStatuses, deviceID)

			return nil
		}
	}

	return fmt.Errorf("设备不存在: %s", deviceID)
}

// ==================== 管理接口 ====================

// SwitchDataSource 切换数据源
func (dm *DataManager) SwitchDataSource(source DataSource) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.currentSource = source

	// 清除所有缓存
	dm.deviceStatuses = make(map[string]types.DeviceStatus)

	// 重新加载数据
	switch source {
	case SourceMiddleware:
		return dm.loadFromMiddleware()
	case SourceDatabase:
		return dm.loadFromDatabase()
	case SourceMock:
		return dm.loadFromMock()
	case SourceFallback:
		return dm.loadFromFallback()
	default:
		return fmt.Errorf("未知的数据源: %s", source)
	}
}

// GetDataSource 获取当前数据源
func (dm *DataManager) GetDataSource() DataSource {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()
	return dm.currentSource
}

// GetStatus 获取数据管理器状态
func (dm *DataManager) GetStatus() map[string]interface{} {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	var cacheInfo map[string]interface{}
	if dm.middleware != nil {
		cacheInfo = dm.middleware.GetCacheInfo()
	} else {
		cacheInfo = map[string]interface{}{
			"middleware_client": "not_initialized",
		}
	}

	// 数据库状态
	var dbStatus map[string]interface{}
	if dm.dbManager != nil {
		dbStatus = dm.dbManager.GetStatus()
	} else {
		dbStatus = map[string]interface{}{
			"enabled": false,
		}
	}

	return map[string]interface{}{
		"data_source":       string(dm.currentSource),
		"device_count":      len(dm.devices),
		"status_cache_size": len(dm.deviceStatuses),
		"last_update":       dm.lastUpdate.Format(time.RFC3339),
		"update_interval":   dm.updateInterval.Seconds(),
		"cache_info":        cacheInfo,
		"database_status":   dbStatus,
		"config": map[string]interface{}{
			"middleware_url":      dm.config.MiddlewareURL,
			"enable_database":     dm.config.EnableDatabase,
			"database_type":       dm.config.DatabaseType,
			"data_retention_days": dm.config.DataRetentionDays,
			"cache_ttl":           dm.config.CacheTTL.Seconds(),
			"enable_fallback":     dm.config.EnableFallback,
			"fallback_timeout":    dm.config.FallbackTimeout.Seconds(),
		},
	}
}

// ClearCache 清除缓存
func (dm *DataManager) ClearCache() {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.deviceStatuses = make(map[string]types.DeviceStatus)
	if dm.middleware != nil {
		dm.middleware.ClearCache()
	}
}

// Stop 停止数据管理器
func (dm *DataManager) Stop() {
	close(dm.stopChan)

	// 关闭数据库管理器
	if dm.dbManager != nil {
		if err := dm.dbManager.Close(); err != nil {
			log.Printf("关闭数据库管理器失败: %v", err)
		}
	}
}

// Refresh 手动刷新数据
func (dm *DataManager) Refresh() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// 清除缓存
	dm.deviceStatuses = make(map[string]types.DeviceStatus)

	// 重新加载数据
	switch dm.currentSource {
	case SourceMiddleware:
		return dm.loadFromMiddleware()
	case SourceDatabase:
		return dm.loadFromDatabase()
	case SourceMock:
		return dm.loadFromMock()
	case SourceFallback:
		return dm.loadFromFallback()
	default:
		return dm.loadFromFallback()
	}
}

// ==================== 数据库相关方法 ====================

// startDatabaseCleanup 启动数据库清理任务
func (dm *DataManager) startDatabaseCleanup() {
	if dm.dbManager == nil {
		return
	}

	// 每天清理一次旧数据
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			retentionDays := dm.config.DataRetentionDays
			if retentionDays <= 0 {
				retentionDays = 30 // 默认保留30天
			}

			if err := dm.dbManager.CleanupOldData(retentionDays); err != nil {
				log.Printf("数据库清理失败: %v", err)
			}
		case <-dm.stopChan:
			return
		}
	}
}

// SyncToDatabase 同步数据到数据库
func (dm *DataManager) SyncToDatabase() error {
	if dm.dbManager == nil {
		return fmt.Errorf("数据库管理器未初始化")
	}

	dm.mutex.RLock()
	devices := dm.devices
	dm.mutex.RUnlock()

	if err := dm.dbManager.SyncDevicesFromMiddleware(devices); err != nil {
		return fmt.Errorf("同步数据到数据库失败: %v", err)
	}

	return nil
}

// GetDatabaseStatistics 获取数据库统计信息
func (dm *DataManager) GetDatabaseStatistics() (map[string]interface{}, error) {
	if dm.dbManager == nil {
		return nil, fmt.Errorf("数据库管理器未初始化")
	}

	return dm.dbManager.GetStatistics()
}

// GetDeviceHistory 获取设备历史数据
func (dm *DataManager) GetDeviceHistory(deviceID string, hours int) ([]types.DeviceStatus, error) {
	if dm.dbManager == nil {
		return nil, fmt.Errorf("数据库管理器未初始化")
	}

	return dm.dbManager.GetDeviceStatusHistory(deviceID, hours)
}

// GetReagentConsumptionTrend 获取试剂消耗趋势
func (dm *DataManager) GetReagentConsumptionTrend(deviceID, reagentName string, hours int) ([]database.ReagentConsumptionTrend, error) {
	if dm.dbManager == nil {
		return nil, fmt.Errorf("数据库管理器未初始化")
	}

	return dm.dbManager.GetReagentConsumptionTrend(deviceID, reagentName, hours)
}