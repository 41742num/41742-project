package gin_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/project47/cmd/mywebapp/models"
)

// GinServerStatusHandler 获取服务器状态 (Gin版本)
func GinServerStatusHandler(c *gin.Context) {
	status := models.GetServerStatus()
	c.JSON(http.StatusOK, status)
}

// GinServerStatsHandler 获取服务器统计信息 (Gin版本)
func GinServerStatsHandler(c *gin.Context) {
	stats := models.GetServerStats()
	c.JSON(http.StatusOK, stats)
}

// GinServerStatusHandlerWithQuery 支持查询参数的版本
func GinServerStatusHandlerWithQuery(c *gin.Context) {
	// 可以添加查询参数支持，例如：
	// detailed := c.Query("detailed") == "true"
	// format := c.Query("format") // json, xml, etc.

	status := models.GetServerStatus()
	c.JSON(http.StatusOK, status)
}

// GinServerStatsHandlerWithQuery 支持查询参数的统计信息版本
func GinServerStatsHandlerWithQuery(c *gin.Context) {
	// 可以添加查询参数支持，例如时间范围
	// startTime := c.Query("start_time")
	// endTime := c.Query("end_time")

	stats := models.GetServerStats()
	c.JSON(http.StatusOK, stats)
}