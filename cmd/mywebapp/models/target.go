package models

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Target 监控目标
type Target struct {
	Name        string   `json:"name"`         // 服务名称，如 "middleware"
	ProcessName string   `json:"process_name"` // 进程名，如 "nginx"
	Port        int      `json:"port"`         // 监控端口，如 8080
	LogPath     string   `json:"log_path"`     // 日志文件路径，如 "/var/log/app/error.log"
	Keywords    []string `json:"keywords"`     // 故障日志关键字，如 ["ERROR", "FATAL"]
}

// Status 状态报告
type Status struct {
	TargetName   string   `json:"target_name"`
	ProcessAlive bool     `json:"process_alive"`
	PortOpen     bool     `json:"port_open"`
	RecentErrors []string `json:"recent_errors"` // 最近10条错误日志
	LastCheck    string   `json:"last_check"`
}

// 全局监控目标列表（可以从配置文件加载）
var Targets = []Target{
	{
		Name:        "nginx",
		ProcessName: "nginx",
		Port:        80,
		LogPath:     "/var/log/nginx/error.log",
		Keywords:    []string{"error", "crit", "alert"},
	},
	{
		Name:        "redis",
		ProcessName: "redis-server",
		Port:        6379,
		LogPath:     "/var/log/redis/redis-server.log",
		Keywords:    []string{"ERR", "FATAL"},
	},
}

// 真实环境下通过 SSH 执行远程命令（这里使用模拟或本地命令）
// 为了演示，我们使用本地命令模拟，实际可替换为 SSH 客户端
func runRemoteCommand(host, command string) (string, error) {
	// TODO: 使用 golang.org/x/crypto/ssh 连接远程服务器
	// 示例：模拟返回固定输出
	if strings.Contains(command, "pgrep") {
		// 模拟进程存在
		if strings.Contains(command, "nginx") {
			return "12345\n", nil
		}
		return "", fmt.Errorf("process not found")
	}
	if strings.Contains(command, "netstat") || strings.Contains(command, "ss") {
		// 模拟端口监听
		return "LISTEN", nil
	}
	if strings.Contains(command, "tail") {
		// 模拟日志输出
		return "2025-01-15 10:00:01 [ERROR] connection refused\n2025-01-15 09:59:00 [INFO] started\n", nil
	}
	return "", nil
}

// CheckProcess 检查进程是否存活
func CheckProcess(target Target) bool {
	// 实际命令: pgrep -f <process_name>
	cmd := exec.Command("pgrep", "-f", target.ProcessName)
	err := cmd.Run()
	return err == nil
}

// CheckPort 检查端口是否开放（通过 netstat 或 ss）
func CheckPort(port int) bool {
	cmd := exec.Command("ss", "-tln", "sport", "=", fmt.Sprintf(":%d", port))
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), fmt.Sprintf(":%d", port))
}

// CheckLogErrors 检查日志文件中的错误关键字，返回最近若干行
func CheckLogErrors(target Target) []string {
	// tail -n 100 日志文件，过滤关键字
	cmd := exec.Command("tail", "-n", "100", target.LogPath)
	output, err := cmd.Output()
	if err != nil {
		return []string{fmt.Sprintf("无法读取日志: %v", err)}
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	var errors []string
	for scanner.Scan() {
		line := scanner.Text()
		for _, kw := range target.Keywords {
			if strings.Contains(strings.ToLower(line), strings.ToLower(kw)) {
				errors = append(errors, line)
				if len(errors) >= 10 {
					break
				}
			}
		}
	}
	return errors
}

// GetStatus 获取单个目标的状态
func GetStatus(target Target) Status {
	processAlive := CheckProcess(target)
	portOpen := CheckPort(target.Port)
	recentErrors := CheckLogErrors(target)
	return Status{
		TargetName:   target.Name,
		ProcessAlive: processAlive,
		PortOpen:     portOpen,
		RecentErrors: recentErrors,
		LastCheck:    time.Now().Format(time.RFC3339),
	}
}

// RestartService 重启服务（例如通过 systemctl 或 supervisor）
func RestartService(targetName string) error {
	// 根据 targetName 找到对应的 Target
	var target *Target
	for i, t := range Targets {
		if t.Name == targetName {
			target = &Targets[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("unknown target: %s", targetName)
	}
	// 实际重启命令，如 systemctl restart <service_name>
	cmd := exec.Command("systemctl", "restart", target.Name)
	err := cmd.Run()
	if err != nil {
		// 尝试使用 service 命令
		cmd = exec.Command("service", target.Name, "restart")
		err = cmd.Run()
	}
	return err
}
