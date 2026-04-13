package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/project47/cmd/mywebapp/data"
	"github.com/project47/cmd/mywebapp/global"
	"github.com/project47/cmd/mywebapp/handlers"
)

func main() {
	// 初始化全局数据管理器
	initDataManager()

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

	// 添加管理API
	mux.HandleFunc("/api/admin/data-source", handlers.DataSourceHandler)
	mux.HandleFunc("/api/admin/refresh", handlers.RefreshDataHandler)
	mux.HandleFunc("/api/admin/status", handlers.DataManagerStatusHandler)

	// 历史数据API
	mux.HandleFunc("/api/history/devices/", handlers.DeviceHistoryHandler)
	mux.HandleFunc("/api/history/devices", handlers.AllDevicesHistoryHandler)
	mux.HandleFunc("/api/history/reagents/", handlers.ReagentConsumptionHistoryHandler)
	mux.HandleFunc("/api/history/statistics/database", handlers.DatabaseStatisticsHandler)
	mux.HandleFunc("/api/history/statistics/data-source", handlers.DataSourceHistoryHandler)

	// 智能查找静态文件目录（三重回退机制）
	staticDir := findStaticDir()
	log.Printf("使用的静态文件目录: %s", staticDir)

	fs := http.FileServer(http.Dir(staticDir))

	// 处理根路径
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// 根路径，直接提供index.html
			indexPath := filepath.Join(staticDir, "index.html")
			http.ServeFile(w, r, indexPath)
		} else {
			// 其他路径，去掉开头的"/"再交给文件服务器
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/")
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
	fmt.Println("数据管理API: http://localhost:8083/api/admin/data-source")
	fmt.Println("数据刷新API: http://localhost:8083/api/admin/refresh")
	fmt.Println("管理器状态API: http://localhost:8083/api/admin/status")
	fmt.Println("设备历史API: http://localhost:8083/api/history/devices/{deviceID}")
	fmt.Println("所有设备历史API: http://localhost:8083/api/history/devices")
	fmt.Println("试剂消耗历史API: http://localhost:8083/api/history/reagents/{deviceID}/{reagentName}")
	fmt.Println("数据库统计API: http://localhost:8083/api/history/statistics/database")
	fmt.Println("数据源历史API: http://localhost:8083/api/history/statistics/data-source")

	log.Fatal(http.ListenAndServe(":8083", mux))
}

// findStaticDir 智能查找静态文件目录
func findStaticDir() string {
	// 方法1：尝试当前工作目录
	cwd, err := os.Getwd()
	if err == nil {
		staticDir := filepath.Join(cwd, "static")
		if _, err := os.Stat(staticDir); err == nil {
			return staticDir
		}
	}

	// 方法2：尝试可执行文件目录
	exeDir, err := os.Executable()
	if err == nil {
		exeDir = filepath.Dir(exeDir)
		staticDir := filepath.Join(exeDir, "static")
		if _, err := os.Stat(staticDir); err == nil {
			return staticDir
		}
	}

	// 方法3：尝试源代码目录（通过调用栈）
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		sourceDir := filepath.Dir(filename)
		staticDir := filepath.Join(sourceDir, "static")
		if _, err := os.Stat(staticDir); err == nil {
			return staticDir
		}
	}

	// 方法4：硬编码回退（开发时使用）
	devPaths := []string{
		"./static",                            // 相对路径
		"E:\\FILE\\gostudy\\project47\\cmd\\mywebapp\\static", // 绝对路径
	}

	for _, path := range devPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// 所有方法都失败
	log.Fatal("无法找到静态文件目录(static/)。请确保static目录存在。")
	return ""
}

// initDataManager 初始化数据管理器
func initDataManager() {
	// 配置数据管理器
	config := &data.Config{
		DataSource:        data.SourceMiddleware, // 默认尝试连接中间件
		MiddlewareURL:     "http://localhost:8080", // 俄罗斯中间件地址
		DatabaseType:      "sqlite",
		SQLitePath:        "./data/project47.db",
		EnableDatabase:    true, // 启用数据库功能
		CacheTTL:          30 * time.Second,
		UpdateInterval:    30 * time.Second,
		MaxRetries:        3,
		RetryDelay:        1 * time.Second,
		EnableFallback:    true, // 启用回退
		FallbackTimeout:   5 * time.Second,
		DataRetentionDays: 30,   // 保留30天历史数据
	}

	// 初始化全局管理器
	if err := global.GetInstance().Initialize(config); err != nil {
		log.Printf("警告: 初始化数据管理器失败: %v", err)
		log.Println("将使用回退数据（本地硬编码数据）")

		// 尝试使用回退配置
		fallbackConfig := &data.Config{
			DataSource:     data.SourceFallback,
			EnableFallback: true,
		}

		if err := global.GetInstance().Initialize(fallbackConfig); err != nil {
			log.Fatalf("致命错误: 无法初始化数据管理器: %v", err)
		}
	}

	// 检查数据源状态
	dm, err := global.GetInstance().GetDataManager()
	if err == nil {
		status := dm.GetStatus()
		log.Printf("数据管理器初始化完成")
		log.Printf("当前数据源: %s", status["data_source"])
		log.Printf("设备数量: %v", status["device_count"])

		// 显示数据库状态
		if dbStatus, ok := status["database_status"].(map[string]interface{}); ok {
			if enabled, ok := dbStatus["enabled"].(bool); ok && enabled {
				log.Printf("数据库功能已启用")
			} else {
				log.Printf("数据库功能未启用或初始化失败")
			}
		}
	}
}