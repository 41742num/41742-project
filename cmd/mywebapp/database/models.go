package database

import (
	"time"

	"github.com/project47/cmd/mywebapp/types"
)

// DeviceModel 数据库设备模型
type DeviceModel struct {
	ID          string    `db:"id"`
	DeviceID    string    `db:"device_id"`
	Name        string    `db:"name"`
	Status      string    `db:"status"`
	IP          string    `db:"ip"`
	Port        int       `db:"port"`
	ProcessName string    `db:"process_name"`
	IsOnline    bool      `db:"is_online"`
	HasFault    bool      `db:"has_fault"`
	LastCheck   time.Time `db:"last_check"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// ReagentModel 数据库试剂模型
type ReagentModel struct {
	ID        int       `db:"id"`
	DeviceID  string    `db:"device_id"`
	Name      string    `db:"name"`
	Current   float64   `db:"current"`
	Capacity  float64   `db:"capacity"`
	Unit      string    `db:"unit"`
	Percent   float64   `db:"percent"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// DeviceStatusHistoryModel 设备状态历史模型
type DeviceStatusHistoryModel struct {
	ID           int       `db:"id"`
	DeviceID     string    `db:"device_id"`
	Status       string    `db:"status"`
	IsOnline     bool      `db:"is_online"`
	HasFault     bool      `db:"has_fault"`
	ReagentStatus string   `db:"reagent_status"`
	Uptime       string    `db:"uptime"`
	CreatedAt    time.Time `db:"created_at"`
}

// ReagentConsumptionHistoryModel 试剂消耗历史模型
type ReagentConsumptionHistoryModel struct {
	ID                     int       `db:"id"`
	DeviceID               string    `db:"device_id"`
	ReagentName            string    `db:"reagent_name"`
	Current                float64   `db:"current"`
	Capacity               float64   `db:"capacity"`
	Percent                float64   `db:"percent"`
	ConsumptionRate        *float64  `db:"consumption_rate"`
	EstimatedRemainingHours *float64 `db:"estimated_remaining_hours"`
	CreatedAt              time.Time `db:"created_at"`
}

// MiddlewareAPILogModel 中间件API日志模型
type MiddlewareAPILogModel struct {
	ID             int       `db:"id"`
	Endpoint       string    `db:"endpoint"`
	Method         string    `db:"method"`
	StatusCode     *int      `db:"status_code"`
	ResponseTimeMs *int      `db:"response_time_ms"`
	Success        bool      `db:"success"`
	ErrorMessage   *string   `db:"error_message"`
	CreatedAt      time.Time `db:"created_at"`
}

// DataSourceHistoryModel 数据源切换历史模型
type DataSourceHistoryModel struct {
	ID         int       `db:"id"`
	Source     string    `db:"source"`
	Reason     *string   `db:"reason"`
	DeviceCount *int     `db:"device_count"`
	CreatedAt  time.Time `db:"created_at"`
}

// CacheStatisticsModel 缓存统计模型
type CacheStatisticsModel struct {
	ID            int       `db:"id"`
	CacheType     string    `db:"cache_type"`
	Hits          int       `db:"hits"`
	Misses        int       `db:"misses"`
	TotalRequests int       `db:"total_requests"`
	HitRate       float64   `db:"hit_rate"`
	PeriodStart   time.Time `db:"period_start"`
	PeriodEnd     *time.Time `db:"period_end"`
}

// ==================== 转换函数 ====================

// ToDeviceModel 将types.Device转换为DeviceModel
func ToDeviceModel(device types.Device) DeviceModel {
	return DeviceModel{
		ID:          device.ID,
		DeviceID:    device.DeviceID,
		Name:        device.Name,
		Status:      device.Status,
		IP:          device.IP,
		Port:        device.Port,
		ProcessName: device.ProcessName,
		IsOnline:    device.IsOnline,
		HasFault:    device.HasFault,
		LastCheck:   device.LastCheck,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// FromDeviceModel 将DeviceModel转换为types.Device
func FromDeviceModel(model DeviceModel, reagents []types.Reagent) types.Device {
	return types.Device{
		ID:          model.ID,
		DeviceID:    model.DeviceID,
		Name:        model.Name,
		Status:      model.Status,
		IP:          model.IP,
		Port:        model.Port,
		ProcessName: model.ProcessName,
		Reagents:    reagents,
		LastCheck:   model.LastCheck,
		IsOnline:    model.IsOnline,
		HasFault:    model.HasFault,
	}
}

// ToReagentModels 将types.Reagent列表转换为ReagentModel列表
func ToReagentModels(deviceID string, reagents []types.Reagent) []ReagentModel {
	models := make([]ReagentModel, len(reagents))
	for i, reagent := range reagents {
		models[i] = ReagentModel{
			DeviceID:  deviceID,
			Name:      reagent.Name,
			Current:   reagent.Current,
			Capacity:  reagent.Capacity,
			Unit:      reagent.Unit,
			Percent:   reagent.Percent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	return models
}

// FromReagentModels 将ReagentModel列表转换为types.Reagent列表
func FromReagentModels(models []ReagentModel) []types.Reagent {
	reagents := make([]types.Reagent, len(models))
	for i, model := range models {
		reagents[i] = types.Reagent{
			Name:     model.Name,
			Current:  model.Current,
			Capacity: model.Capacity,
			Unit:     model.Unit,
			Percent:  model.Percent,
		}
	}
	return reagents
}

// ToDeviceStatusHistoryModel 将设备状态转换为历史记录
func ToDeviceStatusHistoryModel(deviceID string, status types.DeviceStatus) DeviceStatusHistoryModel {
	return DeviceStatusHistoryModel{
		DeviceID:     deviceID,
		Status:       status.Status,
		IsOnline:     status.IsOnline,
		HasFault:     status.HasFault,
		ReagentStatus: status.ReagentStatus,
		Uptime:       status.Uptime,
		CreatedAt:    time.Now(),
	}
}

// ToReagentConsumptionHistoryModel 将试剂信息转换为消耗历史记录
func ToReagentConsumptionHistoryModel(deviceID string, reagent types.Reagent, consumptionRate, estimatedHours *float64) ReagentConsumptionHistoryModel {
	return ReagentConsumptionHistoryModel{
		DeviceID:               deviceID,
		ReagentName:            reagent.Name,
		Current:                reagent.Current,
		Capacity:               reagent.Capacity,
		Percent:                reagent.Percent,
		ConsumptionRate:        consumptionRate,
		EstimatedRemainingHours: estimatedHours,
		CreatedAt:              time.Now(),
	}
}

// DeviceStatusSummary 设备状态汇总视图
type DeviceStatusSummary struct {
	DeviceID         string    `db:"device_id"`
	Name             string    `db:"name"`
	Status           string    `db:"status"`
	IsOnline         bool      `db:"is_online"`
	HasFault         bool      `db:"has_fault"`
	LastCheck        time.Time `db:"last_check"`
	MinReagentPercent float64  `db:"min_reagent_percent"`
	ReagentStatus    string    `db:"reagent_status"`
	HistoryCount     int       `db:"history_count"`
	LastStatusUpdate *time.Time `db:"last_status_update"`
}

// ReagentConsumptionTrend 试剂消耗趋势视图
type ReagentConsumptionTrend struct {
	DeviceID               string    `db:"device_id"`
	ReagentName            string    `db:"reagent_name"`
	CreatedAt              time.Time `db:"created_at"`
	Current                float64   `db:"current"`
	Capacity               float64   `db:"capacity"`
	Percent                float64   `db:"percent"`
	ConsumptionRate        *float64  `db:"consumption_rate"`
	EstimatedRemainingHours *float64 `db:"estimated_remaining_hours"`
	PreviousPercent        *float64  `db:"previous_percent"`
	PercentChange          *float64  `db:"percent_change"`
}