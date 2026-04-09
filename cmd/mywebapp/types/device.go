package types

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

// ServerStatus 服务器状态
type ServerStatus struct {
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	IP           string    `json:"ip"`
	Port         int       `json:"port"`
	WebService   bool      `json:"web_service"`
	APIService   bool      `json:"api_service"`
	ResponseTime int       `json:"response_time"`
	LastCheck    time.Time `json:"last_check"`
	Status       string    `json:"status"` // online/warning/offline
	Uptime       string    `json:"uptime"`
}