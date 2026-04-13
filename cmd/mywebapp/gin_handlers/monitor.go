package gin_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/project47/cmd/mywebapp/models"
)

// GinStatusHandler 返回所有监控目标的状态 (Gin版本)
func GinStatusHandler(c *gin.Context) {
	results := make([]models.Status, 0, len(models.Targets))
	for _, target := range models.Targets {
		results = append(results, models.GetStatus(target))
	}
	c.JSON(http.StatusOK, results)
}

// GinRestartHandler 重启指定服务 (Gin版本)
func GinRestartHandler(c *gin.Context) {
	var req struct {
		Target string `json:"target" binding:"required"`
	}

	// 使用 Gin 的绑定功能解析 JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON or missing target field",
		})
		return
	}

	// 调用原有的重启逻辑
	err := models.RestartService(req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "restart command sent",
	})
}

// GinStatusHandlerWithAdapter 使用适配器模式的版本（临时）
// 这个版本保持与原有处理函数完全相同的响应格式
func GinStatusHandlerWithAdapter(c *gin.Context) {
	// 这里可以调用原有的 StatusHandler 逻辑
	// 为了保持响应格式一致，我们复制原有逻辑
	results := make([]models.Status, 0, len(models.Targets))
	for _, target := range models.Targets {
		results = append(results, models.GetStatus(target))
	}

	// 设置 Content-Type 为 application/json
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, results)
}

// GinRestartHandlerWithAdapter 使用适配器模式的版本（临时）
func GinRestartHandlerWithAdapter(c *gin.Context) {
	var req struct {
		Target string `json:"target"`
	}

	// 手动解析 JSON，保持与原有逻辑一致
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	if req.Target == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing target field",
		})
		return
	}

	err := models.RestartService(req.Target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "restart command sent",
	})
}