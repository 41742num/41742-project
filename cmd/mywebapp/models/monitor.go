package models

import (
	"fmt"
	"net"
	"os/exec"
	"time"
)

// DeviceMonitor 设备监控器
type DeviceMonitor struct {
	CheckInterval time.Duration // 检查间隔
}

// NewDeviceMonitor 创建设备监控器
func NewDeviceMonitor(interval time.Duration) *DeviceMonitor {
	return &DeviceMonitor{
		CheckInterval: interval,
	}
}

// CheckDeviceOnline 检查设备是否在线（通过端口检测）
func CheckDeviceOnline(device Device) bool {
	address := fmt.Sprintf("%s:%d", device.IP, device.Port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// CheckDeviceProcess 检查设备进程是否运行（通过SSH或本地命令）
func CheckDeviceProcess(device Device) bool {
	// 模拟进程检查，实际应通过SSH执行远程命令
	// 这里使用本地pgrep命令作为示例
	cmd := exec.Command("pgrep", "-f", device.ProcessName)
	err := cmd.Run()

	// 如果本地检查失败，尝试模拟远程检查
	if err != nil {
		// 模拟远程设备检查逻辑
		// 实际应使用SSH客户端执行远程命令
		return simulateRemoteProcessCheck(device)
	}

	return true
}

// simulateRemoteProcessCheck 模拟远程进程检查
func simulateRemoteProcessCheck(device Device) bool {
	// 模拟逻辑：根据设备状态返回
	// 实际应通过SSH连接到设备执行命令

	// 模拟俄罗斯中间件服务器总是在线
	if device.IP == "172.19.14.202" {
		return true
	}

	// 模拟其他设备：enabled的设备在线，disabled的设备离线
	return device.Status == "enabled"
}

// CheckDeviceFault 检查设备是否有故障
func CheckDeviceFault(device Device) bool {
	// 模拟故障检查逻辑
	// 实际应通过SSH检查设备日志或状态接口

	// 检查试剂状态：如果有试剂为0，则认为有故障
	for _, reagent := range device.Reagents {
		if reagent.Current == 0 {
			return true
		}
	}

	// 模拟随机故障（用于测试）
	// 在实际应用中，应通过设备API或日志检查
	return false
}

// CheckReagentLevels 检查试剂余量（模拟）
func CheckReagentLevels(device Device) []Reagent {
	// 模拟试剂余量检查
	// 实际应通过设备API获取实时数据

	updatedReagents := make([]Reagent, len(device.Reagents))

	for i, reagent := range device.Reagents {
		// 模拟试剂消耗：每次检查减少0.1-0.5%
		consumption := 0.1 + (float64(i) * 0.1)
		newCurrent := reagent.Current - consumption
		if newCurrent < 0 {
			newCurrent = 0
		}

		updatedReagents[i] = Reagent{
			Name:     reagent.Name,
			Current:  newCurrent,
			Capacity: reagent.Capacity,
			Unit:     reagent.Unit,
			Percent:  (newCurrent / reagent.Capacity) * 100,
		}
	}

	return updatedReagents
}

// MonitorAllDevices 监控所有设备
func MonitorAllDevices() []DeviceStatus {
	var statuses []DeviceStatus

	for i, device := range Devices {
		// 检查设备在线状态
		isOnline := CheckDeviceOnline(device)

		// 检查设备故障
		hasFault := CheckDeviceFault(device)

		// 更新设备状态
		Devices[i].IsOnline = isOnline
		Devices[i].HasFault = hasFault
		Devices[i].LastCheck = time.Now()

		// 更新试剂余量（模拟）
		if isOnline && device.Status == "enabled" {
			Devices[i].Reagents = CheckReagentLevels(device)
		}

		// 获取设备状态报告
		status := GetDeviceStatus(Devices[i])
		statuses = append(statuses, status)
	}

	return statuses
}

// GetDeviceDetailedStatus 获取设备详细状态
func GetDeviceDetailedStatus(deviceID string) (DeviceStatus, error) {
	device := GetDeviceByID(deviceID)
	if device == nil {
		return DeviceStatus{}, fmt.Errorf("device not found: %s", deviceID)
	}

	// 执行实时检查
	isOnline := CheckDeviceOnline(*device)
	hasFault := CheckDeviceFault(*device)

	// 更新设备状态
	UpdateDeviceStatus(deviceID, isOnline, hasFault)

	// 如果在线且启用，更新试剂余量
	if isOnline && device.Status == "enabled" {
		updatedReagents := CheckReagentLevels(*device)
		for _, reagent := range updatedReagents {
			UpdateReagentLevel(deviceID, reagent.Name, reagent.Current)
		}
	}

	// 获取最新状态
	device = GetDeviceByID(deviceID)
	return GetDeviceStatus(*device), nil
}

// RestartDevice 重启设备（模拟）
func RestartDevice(deviceID string) error {
	device := GetDeviceByID(deviceID)
	if device == nil {
		return fmt.Errorf("device not found: %s", deviceID)
	}

	if device.Status != "enabled" {
		return fmt.Errorf("device is disabled, cannot restart")
	}

	// 模拟重启命令
	// 实际应通过SSH执行重启命令
	fmt.Printf("模拟重启设备: %s (%s)\n", device.Name, device.IP)

	// 模拟重启过程：设备会短暂离线然后恢复
	UpdateDeviceStatus(deviceID, false, false)

	// 模拟重启后状态恢复（在实际应用中应有延迟）
	time.Sleep(100 * time.Millisecond) // 模拟重启时间
	UpdateDeviceStatus(deviceID, true, false)

	// 重置试剂余量为满（模拟重启后补充试剂）
	for _, reagent := range device.Reagents {
		UpdateReagentLevel(deviceID, reagent.Name, reagent.Capacity)
	}

	return nil
}

// GetEnabledDevices 获取所有启用的设备
func GetEnabledDevices() []Device {
	var enabledDevices []Device
	for _, device := range Devices {
		if device.Status == "enabled" {
			enabledDevices = append(enabledDevices, device)
		}
	}
	return enabledDevices
}

// GetDeviceStats 获取设备统计信息
func GetDeviceStats() map[string]interface{} {
	total := len(Devices)
	enabled := 0
	online := 0
	withFault := 0

	for _, device := range Devices {
		if device.Status == "enabled" {
			enabled++
		}
		if device.IsOnline {
			online++
		}
		if device.HasFault {
			withFault++
		}
	}

	// 检查试剂状态
	lowReagentDevices := 0
	emptyReagentDevices := 0

	for _, device := range Devices {
		if device.Status == "enabled" {
			minPercent := 100.0
			for _, reagent := range device.Reagents {
				if reagent.Percent < minPercent {
					minPercent = reagent.Percent
				}
			}

			if minPercent == 0 {
				emptyReagentDevices++
			} else if minPercent < 30 {
				lowReagentDevices++
			}
		}
	}

	return map[string]interface{}{
		"total_devices":          total,
		"enabled_devices":        enabled,
		"disabled_devices":       total - enabled,
		"online_devices":         online,
		"offline_devices":        total - online,
		"devices_with_fault":     withFault,
		"devices_low_reagent":    lowReagentDevices,
		"devices_empty_reagent":  emptyReagentDevices,
		"last_update":            time.Now().Format(time.RFC3339),
	}
}