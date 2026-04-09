package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/project47/cmd/mywebapp/models"
)

// StatusHandler 返回所有监控目标的状态
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	results := make([]models.Status, 0, len(models.Targets))
	for _, target := range models.Targets {
		results = append(results, models.GetStatus(target))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// RestartHandler 重启指定服务，接收 JSON: {"target": "nginx"}
func RestartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Target string `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Target == "" {
		http.Error(w, "Missing target field", http.StatusBadRequest)
		return
	}
	err := models.RestartService(req.Target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "restart command sent"})
}
