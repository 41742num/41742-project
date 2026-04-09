package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/project47/cmd/mywebapp/handlers"
)

func main() {
	// 使用自定义多路复用器
	mux := http.NewServeMux()

	// API路由
	mux.HandleFunc("/api/status", handlers.StatusHandler)
	mux.HandleFunc("/api/restart", handlers.RestartHandler)

	// 设备管理API
	mux.HandleFunc("/api/devices", handlers.DevicesHandler)
	mux.HandleFunc("/api/devices/status", handlers.AllDevicesStatusHandler)
	mux.HandleFunc("/api/devices/stats", handlers.DeviceStatsHandler)

	// 模拟数据API
	mux.HandleFunc("/api/simulated/devices", handlers.SimulatedDevicesHandler)
	mux.HandleFunc("/api/simulated/devices/status", handlers.SimulatedAllDevicesStatusHandler)
	mux.HandleFunc("/api/simulated/devices/stats", handlers.SimulatedDeviceStatsHandler)
	mux.HandleFunc("/api/simulated/server/status", handlers.SimulatedServerStatusHandler)
	mux.HandleFunc("/api/simulated/server/stats", handlers.SimulatedServerStatsHandler)
	mux.HandleFunc("/api/simulated/test", handlers.SimulatedTestHandler)
	mux.HandleFunc("/api/simulated/sample", handlers.SimulatedSampleHandler)
	mux.HandleFunc("/api/simulated/override", handlers.OverrideDevicesHandler)

	// 服务器状态API
	mux.HandleFunc("/api/server/status", handlers.ServerStatusHandler)
	mux.HandleFunc("/api/server/stats", handlers.ServerStatsHandler)

	// 模拟设备动态路由
	mux.HandleFunc("/api/simulated/devices/", func(w http.ResponseWriter, r *http.Request) {
		// 动态路由处理
		path := r.URL.Path
		if strings.HasSuffix(path, "/status") {
			handlers.SimulatedDeviceStatusHandler(w, r)
		} else if strings.HasSuffix(path, "/restart") {
			handlers.SimulatedDeviceRestartHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/api/devices/", func(w http.ResponseWriter, r *http.Request) {
		// 动态路由处理
		path := r.URL.Path
		if strings.HasSuffix(path, "/status") {
			handlers.DeviceStatusHandler(w, r)
		} else if strings.HasSuffix(path, "/restart") {
			handlers.DeviceRestartHandler(w, r)
		} else if strings.HasSuffix(path, "/update") {
			handlers.UpdateDeviceHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// 静态文件服务 - 使用绝对路径
	fs := http.FileServer(http.Dir("./static"))

	// 处理根路径
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// 根路径，直接提供index.html
			http.ServeFile(w, r, "./static/index.html")
		} else {
			// 其他路径使用文件服务器
			fs.ServeHTTP(w, r)
		}
	})

	fmt.Println("监控服务启动，监听 8083 端口")
	fmt.Println("管理页面: http://localhost:8083/")
	fmt.Println("状态API: http://localhost:8083/api/status")
	fmt.Println("重启API: http://localhost:8083/api/restart")
	fmt.Println("设备列表: http://localhost:8083/api/devices")
	fmt.Println("设备状态: http://localhost:8083/api/devices/status")
	fmt.Println("设备统计: http://localhost:8083/api/devices/stats")
	fmt.Println("服务器状态: http://localhost:8083/api/server/status")
	fmt.Println("服务器统计: http://localhost:8083/api/server/stats")
	log.Fatal(http.ListenAndServe(":8083", mux))
}