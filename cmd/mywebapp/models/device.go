package models

import (
	"time"
)

// Device 设备/仪器结构体
type Device struct {
	ID          string    `json:"id"`           // 仪器编号no
	DeviceID    string    `json:"device_id"`    // 仪器ID
	Name        string    `json:"name"`         // 仪器名称
	Status      string    `json:"status"`       // enabled/disabled
	IP          string    `json:"ip"`           // 设备IP地址
	Port        int       `json:"port"`         // 监控端口
	ProcessName string    `json:"process_name"` // 进程名
	Reagents    []Reagent `json:"reagents"`     // 试剂信息
	LastCheck   time.Time `json:"last_check"`   // 最后检查时间
	IsOnline    bool      `json:"is_online"`    // 是否在线
	HasFault    bool      `json:"has_fault"`    // 是否有故障
}

// Reagent 试剂信息
type Reagent struct {
	Name     string  `json:"name"`     // 试剂名称
	Current  float64 `json:"current"`  // 当前余量
	Capacity float64 `json:"capacity"` // 总容量
	Unit     string  `json:"unit"`     // 单位
	Percent  float64 `json:"percent"`  // 百分比 (Current/Capacity * 100)
}

// DeviceStatus 设备状态报告
type DeviceStatus struct {
	DeviceID      string    `json:"device_id"`
	DeviceName    string    `json:"device_name"`
	Status        string    `json:"status"`         // enabled/disabled
	IsOnline      bool      `json:"is_online"`      // 是否在线
	HasFault      bool      `json:"has_fault"`      // 是否有故障
	ReagentStatus string    `json:"reagent_status"` // 试剂状态: normal/low/empty
	Reagents      []Reagent `json:"reagents"`       // 试剂详细信息
	LastCheck     string    `json:"last_check"`     // 最后检查时间
	Uptime        string    `json:"uptime"`         // 运行时间
}

// 全局设备列表（可以从配置文件或数据库加载）
var Devices = []Device{
	{
		ID:          "DEV001",
		DeviceID:    "MIDDLEWARE_001",
		Name:        "俄罗斯中间件服务器",
		Status:      "enabled",
		IP:          "172.19.14.202",
		Port:        10001,
		ProcessName: "middleware",
		Reagents: []Reagent{
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
	},
	{
		ID:          "DEV002",
		DeviceID:    "ANALYZER_001",
		Name:        "分析仪器1",
		Status:      "enabled",
		IP:          "172.19.14.203",
		Port:        10002,
		ProcessName: "analyzer",
		Reagents: []Reagent{
			{
				Name:     "标准液",
				Current:  70.0,
				Capacity: 100.0,
				Unit:     "ml",
				Percent:  70.0,
			},
		},
	},
	{
		ID:          "DEV003",
		DeviceID:    "ANALYZER_002",
		Name:        "分析仪器2",
		Status:      "disabled",
		IP:          "172.19.14.204",
		Port:        10003,
		ProcessName: "analyzer",
		Reagents: []Reagent{
			{
				Name:     "标准液",
				Current:  0.0,
				Capacity: 100.0,
				Unit:     "ml",
				Percent:  0.0,
			},
		},
	},
}

// GetDeviceByID 根据设备ID获取设备
func GetDeviceByID(deviceID string) *Device {
	for i, device := range Devices {
		if device.DeviceID == deviceID {
			return &Devices[i]
		}
	}
	return nil
}

// UpdateDeviceStatus 更新设备状态
func UpdateDeviceStatus(deviceID string, isOnline bool, hasFault bool) {
	for i, device := range Devices {
		if device.DeviceID == deviceID {
			Devices[i].IsOnline = isOnline
			Devices[i].HasFault = hasFault
			Devices[i].LastCheck = time.Now()
			break
		}
	}
}

// UpdateReagentLevel 更新试剂余量
func UpdateReagentLevel(deviceID string, reagentName string, currentLevel float64) {
	for i, device := range Devices {
		if device.DeviceID == deviceID {
			for j, reagent := range device.Reagents {
				if reagent.Name == reagentName {
					Devices[i].Reagents[j].Current = currentLevel
					Devices[i].Reagents[j].Percent = (currentLevel / Devices[i].Reagents[j].Capacity) * 100
					break
				}
			}
			Devices[i].LastCheck = time.Now()
			break
		}
	}
}

// GetReagentStatus 获取试剂状态描述
func GetReagentStatus(percent float64) string {
	if percent == 0 {
		return "empty"
	} else if percent < 30 {
		return "low"
	} else if percent < 70 {
		return "warning"
	}
	return "normal"
}

// GetDeviceStatus 获取设备整体状态
func GetDeviceStatus(device Device) DeviceStatus {
	// 检查试剂状态（取最低的试剂百分比）
	minPercent := 100.0
	for _, reagent := range device.Reagents {
		if reagent.Percent < minPercent {
			minPercent = reagent.Percent
		}
	}

	reagentStatus := GetReagentStatus(minPercent)

	return DeviceStatus{
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