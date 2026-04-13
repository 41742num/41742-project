package main

import (
	"fmt"
	"log"
	"time"

	"github.com/project47/cmd/mywebapp/data"
	"github.com/project47/cmd/mywebapp/global"
	"github.com/project47/cmd/mywebapp/router"
)

func main() {
	// 初始化全局数据管理器
	initDataManager()

	// 使用Gin路由
	router := router.SetupRouter()

	fmt.Println("监控服务启动，监听 8083 端口")
	fmt.Println("使用 Gin 框架")
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

	// 启动Gin服务器
	if err := router.Run(":8083"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
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