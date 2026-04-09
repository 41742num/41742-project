# 俄罗斯中间件监控系统 - 数据源管理演示

## 系统架构

```
俄罗斯中间件API → 内存缓存 → API响应
         ↓
      (失败时)
         ↓
Mock数据/回退数据 → API响应
```

## 启动系统

```bash
cd cmd/mywebapp
go run main.go
```

## API端点

### 1. 设备管理API
- `GET /api/devices` - 获取设备列表
- `GET /api/devices/status` - 获取所有设备状态
- `GET /api/devices/{id}/status` - 获取单个设备状态
- `POST /api/devices/{id}/restart` - 重启设备

### 2. 数据管理API
- `GET /api/admin/data-source` - 获取当前数据源
- `POST /api/admin/data-source` - 切换数据源
- `POST /api/admin/refresh` - 手动刷新数据
- `GET /api/admin/status` - 获取数据管理器状态

### 3. 模拟数据API（兼容旧版本）
- `GET /api/simulated/devices` - 获取模拟设备
- `POST /api/simulated/override` - 用模拟数据覆盖

## 使用示例

### 1. 检查当前数据源
```bash
curl http://localhost:8083/api/admin/data-source
```

响应：
```json
{
  "current_source": "middleware",
  "status": {
    "data_source": "middleware",
    "device_count": 3,
    "last_update": "2026-04-09T10:30:00Z",
    "cache_info": {...}
  }
}
```

### 2. 切换到Mock数据
```bash
curl -X POST http://localhost:8083/api/admin/data-source \
  -H "Content-Type: application/json" \
  -d '{"source": "mock"}'
```

### 3. 切换到回退数据（本地硬编码）
```bash
curl -X POST http://localhost:8083/api/admin/data-source \
  -H "Content-Type: application/json" \
  -d '{"source": "fallback"}'
```

### 4. 切回中间件数据
```bash
curl -X POST http://localhost:8083/api/admin/data-source \
  -H "Content-Type: application/json" \
  -d '{"source": "middleware"}'
```

### 5. 手动刷新数据
```bash
curl -X POST http://localhost:8083/api/admin/refresh
```

### 6. 获取设备列表（自动使用当前数据源）
```bash
curl http://localhost:8083/api/devices
```

## 自动切换逻辑

系统启动时：
1. 尝试连接俄罗斯中间件（`http://localhost:8080`）
2. 如果连接成功，使用中间件数据
3. 如果连接失败，自动切换到回退数据
4. 后台每30秒检查中间件连接状态

手动切换：
- 可以通过API随时切换数据源
- 切换后立即生效，所有API使用新数据源

## 配置说明

在 `main.go` 的 `initDataManager()` 函数中配置：

```go
config := &data.Config{
    DataSource:      data.SourceMiddleware, // 默认数据源
    MiddlewareURL:   "http://localhost:8080", // 中间件地址
    CacheTTL:        30 * time.Second, // 缓存时间
    UpdateInterval:  30 * time.Second, // 更新间隔
    EnableFallback:  true, // 启用自动回退
    FallbackTimeout: 5 * time.Second, // 回退超时
}
```

## 故障排除

### 中间件连接失败
1. 检查中间件服务是否运行在 `http://localhost:8080`
2. 检查网络连接
3. 查看日志中的错误信息

### 数据不更新
1. 检查 `GET /api/admin/status` 查看最后更新时间
2. 使用 `POST /api/admin/refresh` 手动刷新
3. 检查缓存配置

### 切换数据源无效
1. 确认POST请求的JSON格式正确
2. 检查数据源名称：`middleware`、`mock`、`fallback`
3. 查看响应中的错误信息

## 扩展功能

### 添加新的Mock数据源
1. 在 `data/manager.go` 的 `loadFromMock()` 中添加mock数据生成逻辑
2. 可以集成现有的 `mock/generator.go`

### 添加其他数据源
1. 在 `data.DataSource` 中添加新的数据源类型
2. 在 `DataManager` 中添加对应的加载逻辑
3. 更新API处理器支持新数据源

### 添加数据库支持
1. 创建数据库客户端
2. 添加新的数据源类型 `database`
3. 实现数据库查询和缓存逻辑