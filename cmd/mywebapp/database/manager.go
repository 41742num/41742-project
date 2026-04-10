package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/project47/cmd/mywebapp/types"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	store    Store
	config   *DBConfig
	logger   *log.Logger
	mu       sync.RWMutex
	isClosed bool
}

// NewDatabaseManager 创建新的数据库管理器
func NewDatabaseManager(config *DBConfig) (*DatabaseManager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	store, err := NewSQLStore(config)
	if err != nil {
		return nil, fmt.Errorf("创建数据库存储失败: %v", err)
	}

	return &DatabaseManager{
		store:  store,
		config: config,
		logger: log.New(log.Writer(), "[DB-Manager] ", log.LstdFlags),
	}, nil
}

// ==================== 设备数据同步 ====================

// SyncDevicesFromMiddleware 从中间件同步设备数据到数据库
func (dm *DatabaseManager) SyncDevicesFromMiddleware(devices []types.Device) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.isClosed {
		return fmt.Errorf("数据库管理器已关闭")
	}

	startTime := time.Now()
	successCount := 0
	errorCount := 0

	// 记录数据源切换
	dm.store.LogDataSourceSwitch("middleware", "定期同步", len(devices))

	for _, device := range devices {
		// 保存设备到数据库
		if err := dm.store.SaveDevice(device); err != nil {
			dm.logger.Printf("保存设备%s到数据库失败: %v", device.DeviceID, err)
			errorCount++
			continue
		}

		// 保存设备状态历史
		status := dm.calculateDeviceStatus(device)
		if err := dm.store.SaveDeviceStatusHistory(device.DeviceID, status); err != nil {
			dm.logger.Printf("保存设备%s状态历史失败: %v", device.DeviceID, err)
			// 继续处理，不中断
		}

		// 保存试剂消耗历史
		for _, reagent := range device.Reagents {
			// 计算消耗速率（简化实现）
			var consumptionRate *float64
			var estimatedHours *float64

			if reagent.Percent < 100 && reagent.Percent > 0 {
				// 这里可以添加更复杂的消耗速率计算逻辑
				rate := 0.5 // 示例值，实际应该基于历史数据计算
				consumptionRate = &rate

				// 计算预计剩余小时数
				remainingHours := (reagent.Current / rate) * 24
				estimatedHours = &remainingHours
			}

			if err := dm.store.SaveReagentConsumptionHistory(device.DeviceID, reagent, consumptionRate, estimatedHours); err != nil {
				dm.logger.Printf("保存设备%s试剂%s消耗历史失败: %v", device.DeviceID, reagent.Name, err)
			}
		}

		successCount++
	}

	duration := time.Since(startTime)
	dm.logger.Printf("设备同步完成: 成功%d个, 失败%d个, 耗时%s", successCount, errorCount, duration)

	// 记录API日志
	dm.store.LogMiddlewareAPI("/sync/devices", "POST", 200, int(duration.Milliseconds()), true, "")

	return nil
}

// GetDevicesFromDatabase 从数据库获取设备数据
func (dm *DatabaseManager) GetDevicesFromDatabase() ([]types.Device, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if dm.isClosed {
		return nil, fmt.Errorf("数据库管理器已关闭")
	}

	devices, err := dm.store.GetAllDevices()
	if err != nil {
		return nil, fmt.Errorf("从数据库获取设备失败: %v", err)
	}

	// 记录缓存统计
	dm.store.UpdateCacheStatistics("device_list", len(devices) > 0)

	return devices, nil
}

// GetDeviceFromDatabase 从数据库获取单个设备
func (dm *DatabaseManager) GetDeviceFromDatabase(deviceID string) (*types.Device, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if dm.isClosed {
		return nil, fmt.Errorf("数据库管理器已关闭")
	}

	device, err := dm.store.GetDevice(deviceID)
	if err != nil {
		return nil, fmt.Errorf("从数据库获取设备%s失败: %v", deviceID, err)
	}

	// 记录缓存统计
	dm.store.UpdateCacheStatistics("device_status", device != nil)

	return device, nil
}

// ==================== 历史数据查询 ====================

// GetDeviceStatusHistory 获取设备状态历史
func (dm *DatabaseManager) GetDeviceStatusHistory(deviceID string, hours int) ([]types.DeviceStatus, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if dm.isClosed {
		return nil, fmt.Errorf("数据库管理器已关闭")
	}

	// 计算时间范围
	endTime := time.Now()
	startTime := endTime.Add(-time.Duration(hours) * time.Hour)

	// 获取历史数据
	history, err := dm.store.GetDeviceStatusHistoryByTimeRange(deviceID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("获取设备%s状态历史失败: %v", deviceID, err)
	}

	return history, nil
}

// GetReagentConsumptionTrend 获取试剂消耗趋势
func (dm *DatabaseManager) GetReagentConsumptionTrend(deviceID, reagentName string, hours int) ([]ReagentConsumptionTrend, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if dm.isClosed {
		return nil, fmt.Errorf("数据库管理器已关闭")
	}

	trend, err := dm.store.GetReagentConsumptionTrend(deviceID, reagentName, hours)
	if err != nil {
		return nil, fmt.Errorf("获取试剂消耗趋势失败: %v", err)
	}

	return trend, nil
}

// ==================== 统计信息 ====================

// GetStatistics 获取统计信息
func (dm *DatabaseManager) GetStatistics() (map[string]interface{}, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if dm.isClosed {
		return nil, fmt.Errorf("数据库管理器已关闭")
	}

	stats := make(map[string]interface{})

	// 设备统计
	totalDevices, err := dm.store.CountDevices()
	if err != nil {
		dm.logger.Printf("获取设备总数失败: %v", err)
	} else {
		stats["total_devices"] = totalDevices
	}

	// 在线设备统计
	onlineDevices, err := dm.store.GetOnlineDeviceCount()
	if err != nil {
		dm.logger.Printf("获取在线设备数失败: %v", err)
	} else {
		stats["online_devices"] = onlineDevices
		stats["offline_devices"] = totalDevices - onlineDevices
	}

	// 低试剂设备统计
	lowReagentThreshold := 30.0
	lowReagentDevices, err := dm.store.GetDeviceWithLowReagentCount(lowReagentThreshold)
	if err != nil {
		dm.logger.Printf("获取低试剂设备数失败: %v", err)
	} else {
		stats["low_reagent_devices"] = lowReagentDevices
	}

	// 数据库统计
	dbStats, err := dm.store.GetDatabaseStats()
	if err != nil {
		dm.logger.Printf("获取数据库统计失败: %v", err)
	} else {
		stats["database"] = dbStats
	}

	// API成功率
	apiSuccessRate, err := dm.store.GetAPISuccessRate(24 * time.Hour)
	if err != nil {
		dm.logger.Printf("获取API成功率失败: %v", err)
	} else {
		stats["api_success_rate_24h"] = apiSuccessRate
	}

	// 平均响应时间
	avgResponseTime, err := dm.store.GetAverageAPIResponseTime(24 * time.Hour)
	if err != nil {
		dm.logger.Printf("获取平均响应时间失败: %v", err)
	} else {
		stats["avg_api_response_time_ms_24h"] = avgResponseTime
	}

	stats["last_update"] = time.Now().Format(time.RFC3339)

	return stats, nil
}

// GetDeviceStatusSummary 获取设备状态汇总
func (dm *DatabaseManager) GetDeviceStatusSummary() ([]DeviceStatusSummary, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if dm.isClosed {
		return nil, fmt.Errorf("数据库管理器已关闭")
	}

	summary, err := dm.store.GetDeviceStatusSummary()
	if err != nil {
		return nil, fmt.Errorf("获取设备状态汇总失败: %v", err)
	}

	return summary, nil
}

// ==================== 维护操作 ====================

// CleanupOldData 清理旧数据
func (dm *DatabaseManager) CleanupOldData(retentionDays int) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.isClosed {
		return fmt.Errorf("数据库管理器已关闭")
	}

	if err := dm.store.CleanupOldData(retentionDays); err != nil {
		return fmt.Errorf("清理旧数据失败: %v", err)
	}

	dm.logger.Printf("成功清理了%d天前的历史数据", retentionDays)
	return nil
}

// BackupDatabase 备份数据库
func (dm *DatabaseManager) BackupDatabase(backupPath string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.isClosed {
		return fmt.Errorf("数据库管理器已关闭")
	}

	// 这里简化实现，实际应该使用数据库的备份功能
	dm.logger.Printf("数据库备份到: %s", backupPath)
	return nil
}

// GetStatus 获取数据库管理器状态
func (dm *DatabaseManager) GetStatus() map[string]interface{} {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	status := map[string]interface{}{
		"is_closed":     dm.isClosed,
		"database_type": dm.config.Type,
		"auto_migrate":  dm.config.AutoMigrate,
	}

	if !dm.isClosed {
		// 获取数据库统计
		if dbStats, err := dm.store.GetDatabaseStats(); err == nil {
			status["database_stats"] = dbStats
		}

		// 获取设备数量
		if count, err := dm.store.CountDevices(); err == nil {
			status["device_count"] = count
		}
	}

	return status
}

// Close 关闭数据库管理器
func (dm *DatabaseManager) Close() error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.isClosed {
		return nil
	}

	if err := dm.store.Close(); err != nil {
		return fmt.Errorf("关闭数据库存储失败: %v", err)
	}

	dm.isClosed = true
	dm.logger.Println("数据库管理器已关闭")
	return nil
}

// ==================== 辅助函数 ====================

// calculateDeviceStatus 计算设备状态
func (dm *DatabaseManager) calculateDeviceStatus(device types.Device) types.DeviceStatus {
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

// LogAPIRequest 记录API请求
func (dm *DatabaseManager) LogAPIRequest(endpoint, method string, statusCode, responseTimeMs int, success bool, errorMessage string) error {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if dm.isClosed {
		return fmt.Errorf("数据库管理器已关闭")
	}

	return dm.store.LogMiddlewareAPI(endpoint, method, statusCode, responseTimeMs, success, errorMessage)
}