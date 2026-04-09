package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/project47/cmd/mywebapp/mock"
)

// SimulatedDevicesHandler 获取模拟的设备列表
func SimulatedDevicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	count := 10
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		if n, err := strconv.Atoi(countStr); err == nil && n > 0 && n <= 100 {
			count = n
		}
	}

	// 生成模拟设备
	simulator := mock.GetDynamicSimulator() //获取动态模拟器实例
	devices := simulator.GetDevices()

	// 如果请求了特定数量，只返回前N个
	if len(devices) > count {
		devices = devices[:count]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// SimulatedDeviceStatusHandler 获取模拟的设备状态
func SimulatedDeviceStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径提取设备ID
	path := strings.TrimPrefix(r.URL.Path, "/api/simulated/devices/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] != "status" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	deviceID := parts[0]
	status, err := mock.GetSimulatedDeviceStatus(deviceID) //获取单个设备的模拟状态
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// SimulatedAllDevicesStatusHandler 获取所有模拟设备状态
func SimulatedAllDevicesStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	statuses := mock.GetSimulatedAllDevicesStatus() //获取所有设备模拟状态
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

// SimulatedDeviceStatsHandler 获取模拟设备统计信息
func SimulatedDeviceStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := mock.GetSimulatedDeviceStats() //获取设备模拟统计
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// SimulatedServerStatusHandler 获取模拟服务器状态
func SimulatedServerStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := mock.GetSimulatedServerStatus() //获取服务器模拟状态
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// SimulatedServerStatsHandler 获取模拟服务器统计信息
func SimulatedServerStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := mock.GetSimulatedServerStats() //获取服务器模拟统计
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// SimulatedDeviceRestartHandler 重启模拟设备
func SimulatedDeviceRestartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径提取设备ID
	path := strings.TrimPrefix(r.URL.Path, "/api/simulated/devices/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] != "restart" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	deviceID := parts[0]
	err := mock.RestartDeviceWithSimulation(deviceID) //模拟重启设备
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":   "Device restart command sent (simulated)",
		"device_id": deviceID,
	})
}

// SimulatedTestHandler 测试模拟数据生成
func SimulatedTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	testData := mock.TestGenerator() //生成测试数据
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(testData)
}

// SimulatedSampleHandler 生成示例数据
func SimulatedSampleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sampleData := mock.GenerateSampleData() //生成示例数据
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sampleData)
}

// OverrideDevicesHandler 用模拟数据覆盖现有设备
func OverrideDevicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取设备数量参数
	count := 10
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		if n, err := strconv.Atoi(countStr); err == nil && n > 0 && n <= 50 {
			count = n
		}
	}

	// 初始化模拟器并覆盖设备
	mock.InitDynamicSimulator(count)
	mock.OverrideDevices()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Devices overridden with simulated data",
		"device_count":  count,
		"original_api":  "Original APIs still available at /api/*",
		"simulated_api": "Simulated APIs available at /api/simulated/*",
	})
}
