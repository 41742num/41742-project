package mock

import (
	"fmt"
	"sync"
	"time"

	"github.com/project47/cmd/mywebapp/models"
)

var (
	dynamicSimulator *DynamicSimulator
	simulatorOnce    sync.Once
)

// InitDynamicSimulator 初始化动态模拟器
func InitDynamicSimulator(deviceCount int) {
	simulatorOnce.Do(func() {
		dynamicSimulator = NewDynamicSimulator(deviceCount)

		// 启动后台更新协程
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				dynamicSimulator.Update()
			}
		}()
	})
}

// GetDynamicSimulator 获取动态模拟器实例
func GetDynamicSimulator() *DynamicSimulator {
	if dynamicSimulator == nil {
		InitDynamicSimulator(10) // 默认10个设备
	}
	return dynamicSimulator
}

// ==================== 与现有模型的集成函数 ====================

// OverrideDevices 用模拟数据覆盖现有设备列表
func OverrideDevices() {
	simulator := GetDynamicSimulator()
	devices := simulator.GetDevices()

	// 清空现有设备列表
	models.Devices = devices
}

// GetSimulatedDeviceStatus 获取模拟的设备状态
func GetSimulatedDeviceStatus(deviceID string) (models.DeviceStatus, error) {
	simulator := GetDynamicSimulator()

	for _, device := range simulator.GetDevices() {
		if device.DeviceID == deviceID {
			return simulator.generator.DeviceStatus(device), nil
		}
	}

	// 如果模拟器中找不到，尝试从原始模型获取
	if device := models.GetDeviceByID(deviceID); device != nil {
		return models.GetDeviceStatus(*device), nil
	}

	return models.DeviceStatus{}, fmt.Errorf("device not found: %s", deviceID)
}

// GetSimulatedAllDevicesStatus 获取所有设备的模拟状态
func GetSimulatedAllDevicesStatus() []models.DeviceStatus {
	simulator := GetDynamicSimulator()
	return simulator.GetDeviceStatuses()
}

// GetSimulatedDeviceStats 获取模拟的设备统计
func GetSimulatedDeviceStats() map[string]interface{} {
	simulator := GetDynamicSimulator()
	return simulator.GetDeviceStats()
}

// GetSimulatedServerStatus 获取模拟的服务器状态
func GetSimulatedServerStatus() models.ServerStatus {
	simulator := GetDynamicSimulator()
	return simulator.GetServerStatus()
}

// GetSimulatedServerStats 获取模拟的服务器统计
func GetSimulatedServerStats() map[string]interface{} {
	server := GetSimulatedServerStatus()

	return map[string]interface{}{
		"name":          server.Name,
		"url":           server.URL,
		"status":        server.Status,
		"web_service":   server.WebService,
		"api_service":   server.APIService,
		"response_time": server.ResponseTime,
		"last_check":    server.LastCheck.Format(time.RFC3339),
		"uptime":        server.Uptime,
		"check_time":    time.Now().Format(time.RFC3339),
	}
}

// RestartDeviceWithSimulation 使用模拟器重启设备
func RestartDeviceWithSimulation(deviceID string) error {
	simulator := GetDynamicSimulator()
	return simulator.RestartDevice(deviceID)
}

// ==================== 测试函数 ====================

// TestGenerator 测试生成器
func TestGenerator() map[string]interface{} {
	gen := NewGenerator()

	return map[string]interface{}{
		"string_template": gen.String("@cname"),
		"integer_range":   gen.Integer(1, 100),
		"float_value":     gen.Float(0, 100, 2),
		"boolean":         gen.Boolean(),
		"datetime":        gen.String("@datetime(yyyy-MM-dd HH:mm:ss)"),
		"ip_address":      gen.String("@ip"),
		"url":             gen.String("@url"),
		"device":          gen.Device(),
		"server_status":   gen.ServerStatus(),
	}
}

// GenerateSampleData 生成示例数据
func GenerateSampleData() map[string]interface{} {
	gen := NewGenerator()

	return map[string]interface{}{
		"devices": gen.Devices(5),
		"device_statuses": gen.DeviceStatuses(3),
		"server": gen.ServerStatus(),
		"stats": map[string]interface{}{
			"total_devices": 5,
			"online_devices": 3,
			"offline_devices": 2,
		},
	}
}