package mock

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/project47/cmd/mywebapp/models"
)

// Generator Mock.js风格的数据生成器
type Generator struct {
	rand *rand.Rand
}

// NewGenerator 创建新的生成器
func NewGenerator() *Generator {
	return &Generator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// WithSeed 设置随机种子
func (g *Generator) WithSeed(seed int64) *Generator {
	g.rand = rand.New(rand.NewSource(seed))
	return g
}

// ==================== 基础数据类型生成 ====================

// String 生成字符串
func (g *Generator) String(template string) string {
	if strings.HasPrefix(template, "@") {
		return g.parseTemplate(template)
	}
	return template
}

// Integer 生成整数
func (g *Generator) Integer(min, max int) int {
	if min == max {
		return min
	}
	return min + g.rand.Intn(max-min+1)
}

// Float 生成浮点数
func (g *Generator) Float(min, max float64, precision int) float64 {
	value := min + g.rand.Float64()*(max-min)
	format := fmt.Sprintf("%%.%df", precision)
	formatted, _ := strconv.ParseFloat(fmt.Sprintf(format, value), 64)
	return formatted
}

// Boolean 生成布尔值
func (g *Generator) Boolean() bool {
	return g.rand.Intn(2) == 1
}

// Choice 从选项中随机选择
func (g *Generator) Choice(options []interface{}) interface{} {
	if len(options) == 0 {
		return nil
	}
	return options[g.rand.Intn(len(options))]
}

// ==================== Mock.js模板解析 ====================

// parseTemplate 解析Mock.js模板
func (g *Generator) parseTemplate(template string) string {
	// 移除@符号
	template = strings.TrimPrefix(template, "@")

	// 解析带参数的模板
	parts := strings.Split(template, "(")
	funcName := parts[0]

	switch funcName {
	case "cname":
		return g.chineseName()
	case "cfirst":
		return g.chineseFirstName()
	case "clast":
		return g.chineseLastName()
	case "ctitle":
		length := 3
		if len(parts) > 1 {
			params := strings.TrimSuffix(parts[1], ")")
			if strings.Contains(params, ",") {
				rangeParts := strings.Split(params, ",")
				min, _ := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				max, _ := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
				length = g.Integer(min, max)
			} else {
				length, _ = strconv.Atoi(params)
			}
		}
		return g.chineseTitle(length)
	case "integer":
		min, max := 1, 100
		if len(parts) > 1 {
			params := strings.TrimSuffix(parts[1], ")")
			if strings.Contains(params, ",") {
				rangeParts := strings.Split(params, ",")
				min, _ = strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				max, _ = strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			} else {
				max, _ = strconv.Atoi(params)
			}
		}
		return strconv.Itoa(g.Integer(min, max))
	case "float":
		min, max, precision := 0.0, 100.0, 2
		if len(parts) > 1 {
			params := strings.TrimSuffix(parts[1], ")")
			paramParts := strings.Split(params, ",")
			if len(paramParts) >= 2 {
				min, _ = strconv.ParseFloat(strings.TrimSpace(paramParts[0]), 64)
				max, _ = strconv.ParseFloat(strings.TrimSpace(paramParts[1]), 64)
			}
			if len(paramParts) >= 3 {
				precision, _ = strconv.Atoi(strings.TrimSpace(paramParts[2]))
			}
		}
		return fmt.Sprintf("%.*f", precision, g.Float(min, max, precision))
	case "datetime":
		format := "2006-01-02 15:04:05"
		if len(parts) > 1 {
			format = strings.TrimSuffix(parts[1], ")")
			format = strings.ReplaceAll(format, "yyyy", "2006")
			format = strings.ReplaceAll(format, "MM", "01")
			format = strings.ReplaceAll(format, "dd", "02")
			format = strings.ReplaceAll(format, "HH", "15")
			format = strings.ReplaceAll(format, "mm", "04")
			format = strings.ReplaceAll(format, "ss", "05")
		}
		// 生成最近30天内的随机时间
		daysAgo := g.Integer(0, 30)
		hoursAgo := g.Integer(0, 24)
		minutesAgo := g.Integer(0, 60)
		t := time.Now().Add(-time.Duration(daysAgo)*24*time.Hour -
			time.Duration(hoursAgo)*time.Hour -
			time.Duration(minutesAgo)*time.Minute)
		return t.Format(format)
	case "ip":
		return fmt.Sprintf("%d.%d.%d.%d",
			g.Integer(172, 172),
			g.Integer(19, 19),
			g.Integer(14, 20),
			g.Integer(1, 254))
	case "url":
		return fmt.Sprintf("http://%d.%d.%d.%d:%d",
			g.Integer(172, 172),
			g.Integer(19, 19),
			g.Integer(14, 20),
			g.Integer(1, 254),
			g.Integer(10000, 11000))
	default:
		return template
	}
}

// ==================== 中文数据生成 ====================

var (
	chineseFirstNames = []string{"张", "王", "李", "赵", "刘", "陈", "杨", "黄", "周", "吴"}
	chineseLastNames  = []string{"伟", "芳", "娜", "秀英", "敏", "静", "丽", "强", "磊", "军"}
	chineseWords      = []string{"服务器", "中间件", "分析仪", "检测器", "控制器", "监控", "设备", "系统", "平台", "终端"}
)

func (g *Generator) chineseName() string {
	return g.chineseFirstName() + g.chineseLastName()
}

func (g *Generator) chineseFirstName() string {
	return chineseFirstNames[g.rand.Intn(len(chineseFirstNames))]
}

func (g *Generator) chineseLastName() string {
	return chineseLastNames[g.rand.Intn(len(chineseLastNames))]
}

func (g *Generator) chineseTitle(length int) string {
	var words []string
	for i := 0; i < length; i++ {
		words = append(words, chineseWords[g.rand.Intn(len(chineseWords))])
	}
	return strings.Join(words, "")
}

// ==================== 设备数据生成 ====================

// Device 生成设备数据
func (g *Generator) Device() models.Device {
	deviceID := fmt.Sprintf("DEV%03d", g.Integer(1, 999))

	return models.Device{
		ID:          deviceID,
		DeviceID:    fmt.Sprintf("%s_%03d", strings.ToUpper(g.chineseTitle(2)), g.Integer(1, 999)),
		Name:        g.chineseTitle(g.Integer(2, 4)) + "设备",
		Status:      g.Choice([]interface{}{"enabled", "disabled"}).(string),
		IP:          g.String("@ip"),
		Port:        g.Integer(10000, 11000),
		ProcessName: g.Choice([]interface{}{"middleware", "analyzer", "controller", "monitor"}).(string),
		Reagents:    g.Reagents(g.Integer(1, 4)),
		LastCheck:   time.Now().Add(-time.Duration(g.Integer(0, 300)) * time.Second),
		IsOnline:    g.Boolean(),
		HasFault:    g.rand.Float32() < 0.1, // 10%故障率
	}
}

// Reagents 生成试剂数据
func (g *Generator) Reagents(count int) []models.Reagent {
	reagentNames := []string{"试剂A", "试剂B", "标准液", "缓冲液", "清洗液", "校准液"}
	units := []string{"ml", "L", "g", "kg"}

	var reagents []models.Reagent
	for i := 0; i < count && i < len(reagentNames); i++ {
		capacity := g.Float(50, 200, 1)
		current := g.Float(0, capacity, 1)

		reagents = append(reagents, models.Reagent{
			Name:     reagentNames[i],
			Current:  current,
			Capacity: capacity,
			Unit:     units[g.rand.Intn(len(units))],
			Percent:  (current / capacity) * 100,
		})
	}
	return reagents
}

// DeviceStatus 生成设备状态
func (g *Generator) DeviceStatus(device models.Device) models.DeviceStatus {
	// 计算试剂状态
	minPercent := 100.0
	for _, reagent := range device.Reagents {
		if reagent.Percent < minPercent {
			minPercent = reagent.Percent
		}
	}

	reagentStatus := "normal"
	if minPercent == 0 {
		reagentStatus = "empty"
	} else if minPercent < 30 {
		reagentStatus = "low"
	} else if minPercent < 70 {
		reagentStatus = "warning"
	}

	// 模拟运行时间（1小时到30天）
	uptimeHours := g.Integer(1, 720)
	uptime := fmt.Sprintf("%d小时", uptimeHours)
	if uptimeHours > 24 {
		uptime = fmt.Sprintf("%d天", uptimeHours/24)
	}

	return models.DeviceStatus{
		DeviceID:      device.DeviceID,
		DeviceName:    device.Name,
		Status:        device.Status,
		IsOnline:      device.IsOnline,
		HasFault:      device.HasFault,
		ReagentStatus: reagentStatus,
		Reagents:      device.Reagents,
		LastCheck:     device.LastCheck.Format(time.RFC3339),
		Uptime:        uptime,
	}
}

// ServerStatus 生成服务器状态
func (g *Generator) ServerStatus() models.ServerStatus {
	responseTime := g.Integer(10, 500)
	status := "online"
	if responseTime > 1000 {
		status = "warning"
	} else if g.rand.Float32() < 0.05 { // 5%离线率
		status = "offline"
	}

	uptimeHours := g.Integer(24, 720)
	uptime := fmt.Sprintf("%d小时", uptimeHours)
	if uptimeHours > 24 {
		uptime = fmt.Sprintf("%d天", uptimeHours/24)
	}

	return models.ServerStatus{
		Name:         "俄罗斯中间件服务器",
		URL:          g.String("@url"),
		IP:           g.String("@ip"),
		Port:         g.Integer(10000, 11000),
		WebService:   status != "offline",
		APIService:   g.rand.Float32() < 0.95, // 95%正常
		ResponseTime: responseTime,
		LastCheck:    time.Now(),
		Status:       status,
		Uptime:       uptime,
	}
}

// ==================== 批量生成 ====================

// Devices 批量生成设备
func (g *Generator) Devices(count int) []models.Device {
	var devices []models.Device
	for i := 0; i < count; i++ {
		devices = append(devices, g.Device())
	}
	return devices
}

// DeviceStatuses 批量生成设备状态
func (g *Generator) DeviceStatuses(count int) []models.DeviceStatus {
	var statuses []models.DeviceStatus
	devices := g.Devices(count)
	for _, device := range devices {
		statuses = append(statuses, g.DeviceStatus(device))
	}
	return statuses
}

// ==================== 动态模拟器 ====================

// DynamicSimulator 动态模拟器
type DynamicSimulator struct {
	generator *Generator
	devices   []models.Device
	server    models.ServerStatus
	lastUpdate time.Time
}

// NewDynamicSimulator 创建动态模拟器
func NewDynamicSimulator(deviceCount int) *DynamicSimulator {
	gen := NewGenerator()
	return &DynamicSimulator{
		generator: gen,
		devices:   gen.Devices(deviceCount),
		server:    gen.ServerStatus(),
		lastUpdate: time.Now(),
	}
}

// Update 更新模拟数据
func (ds *DynamicSimulator) Update() {
	elapsed := time.Since(ds.lastUpdate)
	ds.lastUpdate = time.Now()

	// 更新设备状态
	for i := range ds.devices {
		ds.updateDevice(&ds.devices[i], elapsed)
	}

	// 更新服务器状态
	ds.updateServer(elapsed)
}

// updateDevice 更新单个设备状态
func (ds *DynamicSimulator) updateDevice(device *models.Device, elapsed time.Duration) {
	// 模拟设备随机离线（0.5%的概率）
	if ds.generator.rand.Float64() < 0.005 {
		device.IsOnline = false
	}

	// 离线设备有2%的概率恢复
	if !device.IsOnline && ds.generator.rand.Float64() < 0.02 {
		device.IsOnline = true
	}

	// 模拟设备随机故障（0.3%的概率）
	if device.IsOnline && ds.generator.rand.Float64() < 0.003 {
		device.HasFault = true
	}

	// 故障设备有1%的概率自动恢复
	if device.HasFault && ds.generator.rand.Float64() < 0.01 {
		device.HasFault = false
	}

	// 更新试剂余量（模拟消耗）
	for j := range device.Reagents {
		if device.IsOnline && device.Status == "enabled" {
			// 模拟消耗：每小时消耗0.1-1%的容量
			hourlyConsumption := device.Reagents[j].Capacity * ds.generator.Float(0.001, 0.01, 3)
			consumed := hourlyConsumption * elapsed.Hours()

			device.Reagents[j].Current -= consumed
			if device.Reagents[j].Current < 0 {
				device.Reagents[j].Current = 0
			}
			device.Reagents[j].Percent = (device.Reagents[j].Current / device.Reagents[j].Capacity) * 100
		}
	}

	device.LastCheck = time.Now()
}

// updateServer 更新服务器状态
func (ds *DynamicSimulator) updateServer(elapsed time.Duration) {
	// 模拟响应时间波动
	baseResponse := 50.0
	jitter := ds.generator.Float(-20, 20, 0)
	ds.server.ResponseTime = int(baseResponse + jitter)

	// 模拟服务状态变化
	if ds.generator.rand.Float64() < 0.001 { // 0.1%概率服务异常
		ds.server.WebService = false
		ds.server.Status = "offline"
	} else if ds.server.ResponseTime > 1000 {
		ds.server.Status = "warning"
	} else {
		ds.server.Status = "online"
		ds.server.WebService = true
	}

	ds.server.LastCheck = time.Now()
}

// GetDevices 获取当前设备列表
func (ds *DynamicSimulator) GetDevices() []models.Device {
	ds.Update() // 每次获取都更新
	return ds.devices
}

// GetDeviceStatuses 获取设备状态列表
func (ds *DynamicSimulator) GetDeviceStatuses() []models.DeviceStatus {
	ds.Update()

	var statuses []models.DeviceStatus
	for _, device := range ds.devices {
		statuses = append(statuses, ds.generator.DeviceStatus(device))
	}
	return statuses
}

// GetServerStatus 获取服务器状态
func (ds *DynamicSimulator) GetServerStatus() models.ServerStatus {
	ds.Update()
	return ds.server
}

// GetDeviceStats 获取设备统计信息
func (ds *DynamicSimulator) GetDeviceStats() map[string]interface{} {
	ds.Update()

	total := len(ds.devices)
	enabled := 0
	online := 0
	withFault := 0
	lowReagent := 0
	emptyReagent := 0

	for _, device := range ds.devices {
		if device.Status == "enabled" {
			enabled++
		}
		if device.IsOnline {
			online++
		}
		if device.HasFault {
			withFault++
		}

		// 检查试剂状态
		if device.Status == "enabled" {
			minPercent := 100.0
			for _, reagent := range device.Reagents {
				if reagent.Percent < minPercent {
					minPercent = reagent.Percent
				}
			}

			if minPercent == 0 {
				emptyReagent++
			} else if minPercent < 30 {
				lowReagent++
			}
		}
	}

	return map[string]interface{}{
		"total_devices":         total,
		"enabled_devices":       enabled,
		"disabled_devices":      total - enabled,
		"online_devices":        online,
		"offline_devices":       total - online,
		"devices_with_fault":    withFault,
		"devices_low_reagent":   lowReagent,
		"devices_empty_reagent": emptyReagent,
		"last_update":           time.Now().Format(time.RFC3339),
	}
}

// RestartDevice 重启设备（模拟）
func (ds *DynamicSimulator) RestartDevice(deviceID string) error {
	for i, device := range ds.devices {
		if device.DeviceID == deviceID {
			if device.Status != "enabled" {
				return fmt.Errorf("device is disabled, cannot restart")
			}

			// 模拟重启过程
			ds.devices[i].IsOnline = false
			ds.devices[i].HasFault = false

			// 重置试剂余量
			for j := range ds.devices[i].Reagents {
				ds.devices[i].Reagents[j].Current = ds.devices[i].Reagents[j].Capacity
				ds.devices[i].Reagents[j].Percent = 100.0
			}

			// 模拟重启延迟后恢复在线
			time.Sleep(100 * time.Millisecond)
			ds.devices[i].IsOnline = true
			ds.devices[i].LastCheck = time.Now()

			return nil
		}
	}

	return fmt.Errorf("device not found: %s", deviceID)
}