package models

import (
	//"fmt"
	"net/http"
	"time"
)

// ServerStatus 服务器状态
type ServerStatus struct {
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	IP           string    `json:"ip"`
	Port         int       `json:"port"`
	WebService   bool      `json:"web_service"`   // Web服务是否可用
	APIService   bool      `json:"api_service"`   // API服务是否正常
	ResponseTime int       `json:"response_time"` // 响应时间(ms)
	LastCheck    time.Time `json:"last_check"`
	Status       string    `json:"status"` // online, warning, offline
	Uptime       string    `json:"uptime"` // 运行时间
}

// 俄罗斯中间件服务器配置
var MiddlewareServer = ServerStatus{
	Name:   "俄罗斯中间件服务器",
	URL:    "http://172.19.14.202:10001",
	IP:     "172.19.14.202",
	Port:   10001,
	Status: "online",
	Uptime: "24h",
}

// CheckWebService 检查Web服务是否可用
func CheckWebService() bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(MiddlewareServer.URL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// CheckAPIService 检查API服务（模拟）
func CheckAPIService() bool {
	// 模拟API检查
	// 实际应调用中间件的健康检查API
	return true
}

// CheckResponseTime 检查响应时间（模拟）
func CheckResponseTime() int {
	// 模拟响应时间检查
	start := time.Now()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	_, err := client.Get(MiddlewareServer.URL)
	if err != nil {
		return -1
	}

	elapsed := time.Since(start)
	return int(elapsed.Milliseconds())
}

// GetServerStatus 获取服务器状态
func GetServerStatus() ServerStatus {
	// 检查Web服务
	webService := CheckWebService()

	// 检查API服务
	apiService := CheckAPIService()

	// 检查响应时间
	responseTime := CheckResponseTime()

	// 确定状态
	status := "online"
	if !webService {
		status = "offline"
	} else if responseTime > 1000 { // 响应时间超过1秒为警告
		status = "warning"
	} else if !apiService {
		status = "warning"
	}

	// 更新服务器状态
	MiddlewareServer.WebService = webService
	MiddlewareServer.APIService = apiService
	MiddlewareServer.ResponseTime = responseTime
	MiddlewareServer.LastCheck = time.Now()
	MiddlewareServer.Status = status

	return MiddlewareServer
}

// GetServerStats 获取服务器统计信息
func GetServerStats() map[string]interface{} {
	status := GetServerStatus()

	return map[string]interface{}{
		"name":          status.Name,
		"url":           status.URL,
		"status":        status.Status,
		"web_service":   status.WebService,
		"api_service":   status.APIService,
		"response_time": status.ResponseTime,
		"last_check":    status.LastCheck.Format(time.RFC3339),
		"uptime":        status.Uptime,
		"check_time":    time.Now().Format(time.RFC3339),
	}
}
