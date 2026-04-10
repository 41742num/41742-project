package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/project47/cmd/mywebapp/database"
	"github.com/project47/cmd/mywebapp/global"
	"github.com/project47/cmd/mywebapp/types"
)

// HistoryAPIHandler 历史数据API处理器
type HistoryAPIHandler struct{}

// NewHistoryAPIHandler 创建新的历史数据API处理器
func NewHistoryAPIHandler() *HistoryAPIHandler {
	return &HistoryAPIHandler{}
}

// ==================== 设备历史数据API ====================

// DeviceHistoryHandler 获取设备历史数据
func DeviceHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径提取设备ID
	// 路径格式: /api/history/devices/{deviceID}
	path := r.URL.Path
	deviceID := extractDeviceIDFromPath(path, "/api/history/devices/")
	if deviceID == "" {
		http.Error(w, "设备ID不能为空", http.StatusBadRequest)
		return
	}

	// 获取查询参数
	hoursStr := r.URL.Query().Get("hours")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	// 解析参数
	hours := 24 // 默认24小时
	if hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil && h > 0 {
			hours = h
		}
	}

	limit := 100 // 默认100条
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 获取数据管理器
	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		http.Error(w, fmt.Sprintf("获取数据管理器失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 获取设备历史数据
	history, err := dm.GetDeviceHistory(deviceID, hours)
	if err != nil {
		// 如果历史数据不可用，返回空数组而不是错误
		log.Printf("获取设备%s历史数据失败: %v", deviceID, err)
		history = []types.DeviceStatus{}
	}

	// 应用分页
	start := offset
	end := offset + limit
	if start > len(history) {
		start = len(history)
	}
	if end > len(history) {
		end = len(history)
	}

	pagedHistory := history[start:end]

	response := map[string]interface{}{
		"device_id": deviceID,
		"hours":     hours,
		"total":     len(history),
		"offset":    offset,
		"limit":     limit,
		"count":     len(pagedHistory),
		"history":   pagedHistory,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("编码响应失败: %v", err), http.StatusInternalServerError)
		return
	}
}

// AllDevicesHistoryHandler 获取所有设备历史数据摘要
func AllDevicesHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	hoursStr := r.URL.Query().Get("hours")
	hours := 24 // 默认24小时
	if hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil && h > 0 {
			hours = h
		}
	}

	// 获取数据管理器
	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		http.Error(w, fmt.Sprintf("获取数据管理器失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 获取所有设备
	devices := dm.GetDevices()

	// 为每个设备获取历史数据摘要
	var summaries []map[string]interface{}
	for _, device := range devices {
		history, err := dm.GetDeviceHistory(device.DeviceID, hours)
		if err != nil {
			log.Printf("获取设备%s历史数据失败: %v", device.DeviceID, err)
			continue
		}

		summary := map[string]interface{}{
			"device_id":   device.DeviceID,
			"device_name": device.Name,
			"total_records": len(history),
			"online_rate":  calculateOnlineRate(history),
			"last_status":  getLastStatus(history),
			"has_history":  len(history) > 0,
		}

		if len(history) > 0 {
			summary["first_record"] = history[0].LastCheck
			summary["last_record"] = history[len(history)-1].LastCheck
		}

		summaries = append(summaries, summary)
	}

	response := map[string]interface{}{
		"hours":      hours,
		"total_devices": len(devices),
		"devices_with_history": len(summaries),
		"summaries":  summaries,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("编码响应失败: %v", err), http.StatusInternalServerError)
		return
	}
}

// ==================== 试剂消耗历史API ====================

// ReagentConsumptionHistoryHandler 获取试剂消耗历史
func ReagentConsumptionHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径提取设备ID和试剂名称
	// 路径格式: /api/history/reagents/{deviceID}/{reagentName}
	path := r.URL.Path
	deviceID, reagentName := extractDeviceIDAndReagentNameFromPath(path, "/api/history/reagents/")
	if deviceID == "" || reagentName == "" {
		http.Error(w, "设备ID和试剂名称不能为空", http.StatusBadRequest)
		return
	}

	// 获取查询参数
	hoursStr := r.URL.Query().Get("hours")
	hours := 24 // 默认24小时
	if hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil && h > 0 {
			hours = h
		}
	}

	// 获取数据管理器
	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		http.Error(w, fmt.Sprintf("获取数据管理器失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 获取试剂消耗趋势
	trend, err := dm.GetReagentConsumptionTrend(deviceID, reagentName, hours)
	if err != nil {
		log.Printf("获取试剂消耗趋势失败: %v", err)
		// 返回空数据
		// 创建空趋势数据
		trend = []database.ReagentConsumptionTrend{}
	}

	response := map[string]interface{}{
		"device_id":    deviceID,
		"reagent_name": reagentName,
		"hours":        hours,
		"trend":        trend,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("编码响应失败: %v", err), http.StatusInternalServerError)
		return
	}
}

// DeviceReagentsHistoryHandler 获取设备所有试剂消耗历史
func DeviceReagentsHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径提取设备ID
	path := r.URL.Path
	deviceID := extractDeviceIDFromPath(path, "/api/history/devices/")
	if deviceID == "" {
		http.Error(w, "设备ID不能为空", http.StatusBadRequest)
		return
	}

	// 获取查询参数
	hoursStr := r.URL.Query().Get("hours")
	hours := 24 // 默认24小时
	if hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil && h > 0 {
			hours = h
		}
	}

	// 获取数据管理器
	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		http.Error(w, fmt.Sprintf("获取数据管理器失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 获取设备
	device := dm.GetDeviceByID(deviceID)
	if device == nil {
		http.Error(w, fmt.Sprintf("设备不存在: %s", deviceID), http.StatusNotFound)
		return
	}

	// 获取每个试剂的消耗趋势
	reagentsHistory := make(map[string]interface{})
	for _, reagent := range device.Reagents {
		trend, err := dm.GetReagentConsumptionTrend(deviceID, reagent.Name, hours)
		if err != nil {
			log.Printf("获取试剂%s消耗趋势失败: %v", reagent.Name, err)
			reagentsHistory[reagent.Name] = map[string]interface{}{
				"error": "获取数据失败",
			}
			continue
		}

		reagentsHistory[reagent.Name] = map[string]interface{}{
			"current":  reagent.Current,
			"capacity": reagent.Capacity,
			"percent":  reagent.Percent,
			"trend":    trend,
		}
	}

	response := map[string]interface{}{
		"device_id":    deviceID,
		"device_name":  device.Name,
		"hours":        hours,
		"reagents":     reagentsHistory,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("编码响应失败: %v", err), http.StatusInternalServerError)
		return
	}
}

// ==================== 统计API ====================

// DatabaseStatisticsHandler 获取数据库统计信息
func DatabaseStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 获取数据管理器
	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		http.Error(w, fmt.Sprintf("获取数据管理器失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 获取数据库统计
	stats, err := dm.GetDatabaseStatistics()
	if err != nil {
		log.Printf("获取数据库统计失败: %v", err)
		stats = map[string]interface{}{
			"error": "获取统计信息失败",
		}
	}

	response := map[string]interface{}{
		"statistics": stats,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("编码响应失败: %v", err), http.StatusInternalServerError)
		return
	}
}

// DataSourceHistoryHandler 获取数据源切换历史
func DataSourceHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	// 获取查询参数
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // 默认50条
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 这里简化实现，实际应该从数据库获取
	response := map[string]interface{}{
		"limit":     limit,
		"history":   []map[string]interface{}{},
		"timestamp": time.Now().Format(time.RFC3339),
		"note":     "数据源切换历史功能待实现",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("编码响应失败: %v", err), http.StatusInternalServerError)
		return
	}
}

// ==================== 辅助函数 ====================

// extractDeviceIDFromPath 从路径中提取设备ID
func extractDeviceIDFromPath(path, prefix string) string {
	if len(path) <= len(prefix) {
		return ""
	}
	return path[len(prefix):]
}

// extractDeviceIDAndReagentNameFromPath 从路径中提取设备ID和试剂名称
func extractDeviceIDAndReagentNameFromPath(path, prefix string) (string, string) {
	if len(path) <= len(prefix) {
		return "", ""
	}

	rest := path[len(prefix):]
	for i := 0; i < len(rest); i++ {
		if rest[i] == '/' {
			return rest[:i], rest[i+1:]
		}
	}

	return rest, ""
}

// calculateOnlineRate 计算在线率
func calculateOnlineRate(history []types.DeviceStatus) float64 {
	if len(history) == 0 {
		return 0
	}

	onlineCount := 0
	for _, status := range history {
		if status.IsOnline {
			onlineCount++
		}
	}

	return float64(onlineCount) / float64(len(history)) * 100
}

// getLastStatus 获取最后状态
func getLastStatus(history []types.DeviceStatus) map[string]interface{} {
	if len(history) == 0 {
		return map[string]interface{}{
			"has_data": false,
		}
	}

	last := history[len(history)-1]
	return map[string]interface{}{
		"has_data":      true,
		"status":        last.Status,
		"is_online":     last.IsOnline,
		"has_fault":     last.HasFault,
		"reagent_status": last.ReagentStatus,
		"last_check":    last.LastCheck,
	}
}