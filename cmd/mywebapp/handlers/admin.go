package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/project47/cmd/mywebapp/data"
	"github.com/project47/cmd/mywebapp/global"
)

// DataSourceHandler 管理数据源
func DataSourceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getDataSource(w, r)
	case http.MethodPost:
		setDataSource(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getDataSource 获取当前数据源
func getDataSource(w http.ResponseWriter, r *http.Request) {
	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		http.Error(w, "数据管理器未初始化", http.StatusServiceUnavailable)
		return
	}

	status := dm.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"current_source": status["data_source"],
		"status":         status,
	})
}

// setDataSource 设置数据源
func setDataSource(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Source data.DataSource `json:"source"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 验证数据源
	switch request.Source {
	case data.SourceMiddleware, data.SourceMock, data.SourceFallback:
		// 有效数据源
	default:
		http.Error(w, "Invalid data source", http.StatusBadRequest)
		return
	}

	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		http.Error(w, "数据管理器未初始化", http.StatusServiceUnavailable)
		return
	}

	// 切换数据源
	if err := dm.SwitchDataSource(request.Source); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := dm.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "数据源切换成功",
		"current_source": status["data_source"],
		"status":         status,
	})
}

// RefreshDataHandler 刷新数据
func RefreshDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		http.Error(w, "数据管理器未初始化", http.StatusServiceUnavailable)
		return
	}

	// 手动刷新数据
	if err := dm.Refresh(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 清除缓存
	dm.ClearCache()

	status := dm.GetStatus()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "数据刷新成功",
		"current_source": status["data_source"],
		"device_count":   status["device_count"],
		"last_update":    status["last_update"],
	})
}

// DataManagerStatusHandler 获取数据管理器状态
func DataManagerStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		// 返回基本状态
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"initialized": false,
			"error":       err.Error(),
		})
		return
	}

	status := dm.GetStatus()
	status["initialized"] = true

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// AutoSwitchHandler 自动切换数据源（内部使用）
func AutoSwitchHandler() {
	dm, err := global.GetInstance().GetDataManager()
	if err != nil {
		return
	}

	currentSource := dm.GetDataSource()

	// 如果当前是中间件源，检查连接状态
	if currentSource == data.SourceMiddleware {
		status := dm.GetStatus()
		cacheInfo, ok := status["cache_info"].(map[string]interface{})
		if ok {
			// 检查中间件连接状态
			if middlewareStatus, ok := cacheInfo["middleware_status"]; ok {
				if middlewareStatus == "disconnected" || middlewareStatus == "error" {
					// 自动切换到回退数据
					dm.SwitchDataSource(data.SourceFallback)
				}
			}
		}
	}
}