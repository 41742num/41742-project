package router

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/project47/cmd/mywebapp/gin_handlers"
	"github.com/project47/cmd/mywebapp/handlers"
)

// SetupRouter 配置并返回Gin路由器
func SetupRouter() *gin.Engine {
	// 创建Gin实例，使用默认中间件（Logger和Recovery）
	router := gin.Default()

	// 配置静态文件服务
	setupStaticFiles(router)

	// 配置API路由
	setupAPIRoutes(router)

	return router
}

// setupStaticFiles 配置静态文件服务
func setupStaticFiles(router *gin.Engine) {
	// 智能查找静态文件目录
	staticDir := findStaticDir()

	// 使用Gin的静态文件服务
	router.Static("/static", staticDir)

	// 根路径重定向到index.html
	router.GET("/", func(c *gin.Context) {
		indexPath := filepath.Join(staticDir, "index.html")
		c.File(indexPath)
	})

	// 其他静态文件路由
	router.GET("/devices.html", func(c *gin.Context) {
		filePath := filepath.Join(staticDir, "devices.html")
		c.File(filePath)
	})

	router.GET("/history.html", func(c *gin.Context) {
		filePath := filepath.Join(staticDir, "history.html")
		c.File(filePath)
	})

	router.GET("/simulated.html", func(c *gin.Context) {
		filePath := filepath.Join(staticDir, "simulated.html")
		c.File(filePath)
	})
}

// setupAPIRoutes 配置API路由
func setupAPIRoutes(router *gin.Engine) {
	// API路由组
	api := router.Group("/api")
	{
		// 监控API
		api.GET("/status", gin_handlers.GinStatusHandler)
		api.POST("/restart", gin_handlers.GinRestartHandler)

		// 服务器状态API
		api.GET("/server/status", gin_handlers.GinServerStatusHandler)
		api.GET("/server/stats", gin_handlers.GinServerStatsHandler)

		// 设备管理API路由组
		setupDeviceRoutes(api)

		// 模拟数据API路由组
		setupSimulatedRoutes(api)

		// 管理API路由组
		setupAdminRoutes(api)

		// 历史数据API路由组
		setupHistoryRoutes(api)
	}
}

// setupDeviceRoutes 配置设备管理API路由
func setupDeviceRoutes(api *gin.RouterGroup) {
	devices := api.Group("/devices")
	{
		// 使用适配器暂时保持兼容性
		devices.GET("", handlers.AdaptHandler(handlers.DevicesHandler))
		devices.GET("/status", handlers.AdaptHandler(handlers.AllDevicesStatusHandler))
		devices.GET("/stats", handlers.AdaptHandler(handlers.DeviceStatsHandler))

		// 动态路由 - 使用Gin的参数绑定
		devices.GET("/:id/status", handlers.AdaptHandler(handlers.DeviceStatusHandler))
		devices.POST("/:id/restart", handlers.AdaptHandler(handlers.DeviceRestartHandler))
		devices.PUT("/:id/update", handlers.AdaptHandler(handlers.UpdateDeviceHandler))
	}
}

// setupSimulatedRoutes 配置模拟数据API路由
func setupSimulatedRoutes(api *gin.RouterGroup) {
	simulated := api.Group("/simulated")
	{
		// 使用适配器暂时保持兼容性
		simulated.GET("/devices", handlers.AdaptHandler(handlers.SimulatedDevicesHandler))
		simulated.GET("/devices/status", handlers.AdaptHandler(handlers.SimulatedAllDevicesStatusHandler))
		simulated.GET("/devices/stats", handlers.AdaptHandler(handlers.SimulatedDeviceStatsHandler))
		simulated.GET("/server/status", handlers.AdaptHandler(handlers.SimulatedServerStatusHandler))
		simulated.GET("/server/stats", handlers.AdaptHandler(handlers.SimulatedServerStatsHandler))
		simulated.GET("/test", handlers.AdaptHandler(handlers.SimulatedTestHandler))
		simulated.GET("/sample", handlers.AdaptHandler(handlers.SimulatedSampleHandler))
		simulated.POST("/override", handlers.AdaptHandler(handlers.OverrideDevicesHandler))

		// 动态路由
		simulated.GET("/devices/:id/status", handlers.AdaptHandler(handlers.SimulatedDeviceStatusHandler))
		simulated.POST("/devices/:id/restart", handlers.AdaptHandler(handlers.SimulatedDeviceRestartHandler))
	}
}

// setupAdminRoutes 配置管理API路由
func setupAdminRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin")
	{
		// 使用适配器暂时保持兼容性
		admin.GET("/data-source", handlers.AdaptHandler(handlers.DataSourceHandler))
		admin.POST("/refresh", handlers.AdaptHandler(handlers.RefreshDataHandler))
		admin.GET("/status", handlers.AdaptHandler(handlers.DataManagerStatusHandler))
	}
}

// setupHistoryRoutes 配置历史数据API路由
func setupHistoryRoutes(api *gin.RouterGroup) {
	history := api.Group("/history")
	{
		// 使用适配器暂时保持兼容性
		history.GET("/devices", handlers.AdaptHandler(handlers.AllDevicesHistoryHandler))
		history.GET("/statistics/database", handlers.AdaptHandler(handlers.DatabaseStatisticsHandler))
		history.GET("/statistics/data-source", handlers.AdaptHandler(handlers.DataSourceHistoryHandler))

		// 动态路由
		history.GET("/devices/:deviceID", handlers.AdaptHandler(handlers.DeviceHistoryHandler))
		history.GET("/reagents/:deviceID/:reagentName", handlers.AdaptHandler(handlers.ReagentConsumptionHistoryHandler))
	}
}

// findStaticDir 智能查找静态文件目录（从main.go复制）
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
	panic("无法找到静态文件目录(static/)。请确保static目录存在。")
}

// SetupTestRouter 配置测试用的路由器（不包含静态文件）
func SetupTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	setupAPIRoutes(router)
	return router
}

// HealthCheck 健康检查端点
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "project47-web-monitor",
		"version": "1.0.0",
		"gin":     true,
	})
}

// NotFoundHandler 404处理
func NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error":   "Not Found",
		"message": "The requested resource was not found",
		"path":    c.Request.URL.Path,
	})
}

// MethodNotAllowedHandler 405处理
func MethodNotAllowedHandler(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"error":   "Method Not Allowed",
		"message": "The requested method is not allowed for this resource",
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
	})
}