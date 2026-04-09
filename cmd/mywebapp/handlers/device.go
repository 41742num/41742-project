package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/project47/cmd/mywebapp/models"
)

// DevicesHandler 获取所有设备列表
func DevicesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	statusFilter := r.URL.Query().Get("status")
	onlineFilter := r.URL.Query().Get("online")

	var devices []models.Device
	if statusFilter == "enabled" {
		devices = models.GetEnabledDevices()
	} else {
		devices = models.Devices
	}

	// 过滤在线状态
	if onlineFilter != "" {
		var filtered []models.Device
		isOnline := onlineFilter == "true"
		for _, device := range devices {
			if device.IsOnline == isOnline {
				filtered = append(filtered, device)
			}
		}
		devices = filtered
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// DeviceStatusHandler 获取设备状态
func DeviceStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径提取设备ID
	path := strings.TrimPrefix(r.URL.Path, "/api/devices/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] != "status" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	deviceID := parts[0]
	status, err := models.GetDeviceDetailedStatus(deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// AllDevicesStatusHandler 获取所有设备状态
func AllDevicesStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	statuses := models.MonitorAllDevices()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

// DeviceRestartHandler 重启设备
func DeviceRestartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径提取设备ID
	path := strings.TrimPrefix(r.URL.Path, "/api/devices/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[1] != "restart" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	deviceID := parts[0]
	err := models.RestartDevice(deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Device restart command sent",
		"device_id": deviceID,
	})
}

// DeviceStatsHandler 获取设备统计信息
func DeviceStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := models.GetDeviceStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// UpdateDeviceHandler 更新设备信息（管理员）
func UpdateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径提取设备ID
	path := strings.TrimPrefix(r.URL.Path, "/api/devices/")
	deviceID := strings.TrimSuffix(path, "/update")

	var updateData struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	device := models.GetDeviceByID(deviceID)
	if device == nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	// 更新设备信息
	// 注意：这里简化处理，实际应更新数据库
	if updateData.Name != "" {
		device.Name = updateData.Name
	}
	if updateData.Status != "" && (updateData.Status == "enabled" || updateData.Status == "disabled") {
		device.Status = updateData.Status
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Device updated successfully",
		"device_id": deviceID,
	})
}