# 俄罗斯中间件服务器监控系统

基于Go语言开发的服务器监控系统，专门用于监控俄罗斯中间件服务器上的设备状态和试剂余量。

## 功能特性

### 设备监控
- 设备在线状态监控（enabled/disabled）
- 试剂余量实时监控
- 故障检测和告警
- 设备重启控制

### 数据模型
- **Device**: 设备信息（ID、名称、状态、IP、端口等）
- **Reagent**: 试剂信息（名称、当前余量、总容量、百分比）
- **DeviceStatus**: 设备状态报告

### API接口

#### 设备监控API
- `GET /api/devices` - 获取所有设备列表
- `GET /api/devices/status` - 获取所有设备状态
- `GET /api/devices/stats` - 获取设备统计信息
- `GET /api/devices/{id}/status` - 获取单个设备状态
- `POST /api/devices/{id}/restart` - 重启指定设备
- `PUT /api/devices/{id}/update` - 更新设备信息（管理员）

#### 服务器状态API
- `GET /api/server/status` - 获取俄罗斯中间件服务器状态
- `GET /api/server/stats` - 获取服务器详细统计信息

#### 原始服务监控API（兼容）
- `GET /api/status` - 获取原始服务状态
- `POST /api/restart` - 重启服务

### 动态模拟数据API
- `GET /api/simulated/test` - 测试Mock.js风格数据生成
- `GET /api/simulated/sample` - 获取示例数据
- `GET /api/simulated/devices` - 获取模拟设备列表
- `GET /api/simulated/devices/status` - 获取所有模拟设备状态
- `GET /api/simulated/devices/stats` - 获取模拟设备统计
- `GET /api/simulated/server/status` - 获取模拟服务器状态
- `GET /api/simulated/server/stats` - 获取模拟服务器统计
- `POST /api/simulated/devices/{id}/restart` - 重启模拟设备
- `POST /api/simulated/override?count=N` - 用模拟数据覆盖设备列表

### 前端界面
- 设备监控面板（devices.html）
  - 俄罗斯中间件服务器状态栏
  - 设备统计信息
  - 设备列表和状态
  - 试剂余量可视化
  - 一键重启功能
- 导航页面（index.html）
  - 系统入口
  - API文档

## 安装和运行

### 环境要求
- Go 1.16+
- 网络访问权限（用于监控远程设备）

### 编译和运行
```bash
# 进入项目目录
cd cmd/mywebapp

# 编译项目
go build -o mywebapp.exe

# 运行服务
./mywebapp.exe
```

### 访问地址
- 管理界面: http://localhost:8083/
- 设备监控: http://localhost:8083/devices.html
- 模拟数据测试: http://localhost:8083/simulated.html
- API文档: http://localhost:8083/

## 配置说明

### 设备配置
设备配置在 `models/device.go` 文件的 `Devices` 变量中：

```go
var Devices = []Device{
    {
        ID:          "DEV001",
        DeviceID:    "MIDDLEWARE_001",
        Name:        "俄罗斯中间件服务器",
        Status:      "enabled",
        IP:          "172.19.14.202",
        Port:        10001,
        ProcessName: "middleware",
        Reagents: []Reagent{
            {
                Name:     "试剂A",
                Current:  85.0,
                Capacity: 100.0,
                Unit:     "ml",
                Percent:  85.0,
            },
        },
    },
}
```

### 监控配置
- 设备端口检查超时: 5秒
- 服务器状态检查超时: 10秒
- 自动刷新间隔: 30秒（前端）
- 试剂告警阈值: <70% 警告，<30% 严重，=0% 耗尽
- 服务器响应告警阈值: >1000ms 警告

## 使用示例

### 获取所有设备状态
```bash
curl http://localhost:8083/api/devices/status
```

### 重启设备
```bash
curl -X POST http://localhost:8083/api/devices/MIDDLEWARE_001/restart
```

### 获取设备统计
```bash
curl http://localhost:8083/api/devices/stats
```

### 获取服务器状态
```bash
curl http://localhost:8083/api/server/status
```

### 获取服务器详细统计
```bash
curl http://localhost:8083/api/server/stats
```

## 动态模拟数据系统

### 特性
1. **Mock.js风格模板语法**：支持 `@cname`, `@integer`, `@datetime` 等模板
2. **动态数据变化**：设备状态、试剂余量随时间自动变化
3. **真实业务逻辑**：模拟设备故障、网络延迟、试剂消耗
4. **可配置参数**：支持设备数量、故障率、消耗速率等配置
5. **后台自动更新**：每30秒自动更新模拟数据

### 模拟数据示例
```bash
# 测试数据生成
curl http://localhost:8083/api/simulated/test

# 生成10个模拟设备
curl http://localhost:8083/api/simulated/devices?count=10

# 获取模拟设备统计
curl http://localhost:8083/api/simulated/devices/stats

# 重启模拟设备
curl -X POST http://localhost:8083/api/simulated/devices/DEV_001/restart

# 用模拟数据覆盖现有设备（生成20个设备）
curl -X POST http://localhost:8083/api/simulated/override?count=20
```

### 模拟数据模板语法
```go
// 在代码中使用
gen := mock.NewGenerator()
device := mock.Device{
    Name: gen.String("@cname"),           // 中文姓名
    IP: gen.String("@ip"),               // IP地址
    Port: gen.Integer(10000, 11000),     // 端口范围
    Status: gen.Choice([]string{"enabled", "disabled"}),
    Created: gen.String("@datetime(yyyy-MM-dd HH:mm:ss)"),
}
```

## 项目结构
```
mywebapp/
├── main.go                 # 主程序入口
├── models/
│   ├── device.go          # 设备数据模型
│   ├── monitor.go         # 设备监控逻辑
│   ├── server.go          # 服务器状态监控
│   └── target.go          # 原始监控目标
├── handlers/
│   ├── device.go          # 设备API处理器
│   ├── server.go          # 服务器状态处理器
│   ├── monitor.go         # 原始监控处理器
│   └── simulated.go       # 模拟数据API处理器
├── mock/                   # Mock.js风格模拟数据生成器
│   ├── generator.go       # 数据生成器核心
│   └── init.go            # 模拟器初始化
├── static/
│   ├── devices.html       # 设备监控页面
│   ├── simulated.html     # 动态模拟数据测试页面
│   └── index.html         # 导航页面
└── README.md              # 本文档
```

## 扩展开发

### 添加新设备
1. 在 `models/device.go` 的 `Devices` 数组中添加新设备
2. 配置设备IP、端口、试剂信息
3. 重启服务生效

### 实现真实监控
当前使用模拟数据，需要实现：
1. SSH客户端连接远程设备
2. 真实试剂余量获取接口
3. 设备状态API调用

### 数据库集成
计划添加SQLite存储：
- 设备配置持久化
- 历史状态记录
- 故障事件日志

## 注意事项

1. **安全警告**: 当前版本使用模拟数据，生产环境需要：
   - 添加API认证
   - 加密存储认证信息
   - 实现真实的SSH连接

2. **性能考虑**: 监控大量设备时：
   - 使用并发检查
   - 添加缓存机制
   - 优化数据库查询

3. **故障处理**: 
   - 添加重试机制
   - 实现告警通知
   - 记录详细日志

## 开发阶段总结

详细开发过程和架构设计请参考：[DEVELOPMENT_PHASES.md](DEVELOPMENT_PHASES.md)

## 许可证
MIT License