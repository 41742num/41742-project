# 项目架构升级文档

## 升级目标
实现 **俄罗斯中间件API > 数据库(历史数据) > 内存缓存API响应** 的三层架构

## 架构概述

### 三层数据流
1. **第一层：俄罗斯中间件API** - 实时数据源
2. **第二层：数据库(历史数据)** - 历史数据存储和回退
3. **第三层：内存缓存** - API响应加速

### 数据优先级
1. 优先从俄罗斯中间件API获取实时数据
2. 中间件不可用时从数据库读取历史数据
3. 数据库不可用时使用本地回退数据
4. 所有层级都支持内存缓存

## 新增组件

### 1. 数据库层 (`database/`)
- `config.go` - 数据库配置
- `models.go` - 数据库模型定义
- `store.go` - 数据库存储接口和SQLite实现
- `manager.go` - 数据库管理器
- `schema.sql` - 数据库表结构定义

### 2. 历史数据API (`handlers/history.go`)
- 设备状态历史查询
- 试剂消耗趋势分析
- 数据库统计信息
- 数据源切换历史

### 3. 前端历史数据页面 (`static/history.html`)
- 设备历史状态查询
- 试剂消耗趋势图表
- 统计分析展示
- 响应式设计

## 数据库设计

### 主要表结构
1. **devices** - 设备基本信息
2. **reagents** - 设备试剂信息
3. **device_status_history** - 设备状态历史记录
4. **reagent_consumption_history** - 试剂消耗历史记录
5. **middleware_api_logs** - API调用日志
6. **data_source_history** - 数据源切换历史

### 索引优化
- 设备ID索引
- 时间范围索引
- 状态索引

## API端点扩展

### 历史数据API
```
GET  /api/history/devices/{deviceID}          # 获取设备历史状态
GET  /api/history/devices                     # 获取所有设备历史摘要
GET  /api/history/reagents/{device}/{reagent} # 获取试剂消耗历史
GET  /api/history/statistics/database         # 获取数据库统计
GET  /api/history/statistics/data-source      # 获取数据源历史
```

### 管理API增强
```
GET  /api/admin/status        # 包含数据库状态信息
POST /api/admin/refresh       # 支持数据库同步
POST /api/admin/data-source   # 支持切换到数据库源
```

## 缓存策略

### 三级缓存
1. **内存缓存** - 设备状态、设备列表 (TTL: 30秒)
2. **数据库缓存** - 历史数据 (长期存储)
3. **中间件API缓存** - 实时数据 (TTL: 30秒)

### 缓存更新机制
- 定时更新: 每30秒从中间件同步数据到数据库
- 事件驱动: 设备状态变化时更新缓存
- 手动刷新: 支持管理员手动刷新

## 错误处理

### 故障转移策略
1. **中间件API失败** → 回退到数据库历史数据
2. **数据库失败** → 回退到本地硬编码数据
3. **所有数据源失败** → 显示维护页面

### 重试机制
- 最大重试次数: 3次
- 重试延迟: 1秒
- 超时时间: 10秒

## 数据同步

### 同步流程
1. 从中间件API获取实时设备数据
2. 保存设备信息到数据库
3. 记录设备状态历史
4. 记录试剂消耗历史
5. 更新内存缓存

### 同步频率
- 实时同步: 设备状态变化时立即同步
- 定时同步: 每30秒全量同步
- 手动同步: 管理员触发同步

## 前端增强

### 新页面
- `history.html` - 历史数据查询页面
  - 设备历史状态查询
  - 试剂消耗趋势图表
  - 统计分析展示
  - 数据筛选和分页

### 功能增强
- 实时监控页面添加历史数据链接
- 设备管理页面优化
- 响应式设计改进

## 配置说明

### 数据库配置
```go
type DBConfig struct {
    Type           string        // sqlite, postgres, mysql
    SQLitePath     string        // SQLite数据库路径
    AutoMigrate    bool          // 自动迁移
    EnableQueryLog bool          // 查询日志
}
```

### 数据管理器配置
```go
type Config struct {
    DataSource        DataSource    // middleware, database, mock, fallback
    MiddlewareURL     string        // 俄罗斯中间件地址
    EnableDatabase    bool          // 启用数据库功能
    DatabaseType      string        // 数据库类型
    DataRetentionDays int           // 历史数据保留天数
    CacheTTL          time.Duration // 缓存TTL
    UpdateInterval    time.Duration // 更新间隔
    EnableFallback    bool          // 启用回退
}
```

## 部署说明

### 环境要求
- Go 1.25.5+
- SQLite3 (嵌入式，无需单独安装)
- 网络访问俄罗斯中间件API

### 启动命令
```bash
go run main.go
```

### 访问地址
- 实时监控: http://localhost:8083/
- 历史数据: http://localhost:8083/history.html
- 设备管理: http://localhost:8083/devices.html
- 模拟数据: http://localhost:8083/simulated.html

## 监控和维护

### 监控指标
- 数据库连接状态
- 中间件API响应时间
- 缓存命中率
- 数据同步状态
- 错误率统计

### 维护任务
- 定期清理旧数据 (默认保留30天)
- 数据库备份
- 缓存清理
- 日志轮转

## 性能优化

### 数据库优化
- 索引优化
- 查询优化
- 连接池配置
- 批量操作

### 缓存优化
- 分级缓存策略
- 智能缓存失效
- 缓存预热
- 内存管理

### API优化
- 响应压缩
- 请求合并
- 异步处理
- 限流保护

## 安全考虑

### 数据安全
- 数据库加密 (SQLite加密扩展)
- 访问控制
- 审计日志
- 数据备份

### API安全
- 请求验证
- 频率限制
- 错误信息隐藏
- 日志脱敏

## 扩展性

### 未来扩展
1. **支持更多数据库** - PostgreSQL, MySQL
2. **分布式缓存** - Redis集群
3. **实时通知** - WebSocket推送
4. **数据分析** - 大数据分析平台
5. **移动端** - 响应式移动应用

## 故障排除

### 常见问题
1. **数据库连接失败** - 检查数据库文件权限
2. **中间件API超时** - 检查网络连接和中间件状态
3. **缓存不一致** - 手动刷新缓存
4. **历史数据缺失** - 检查数据同步状态

### 日志查看
- 应用日志: 控制台输出
- 数据库日志: 启用查询日志
- API日志: 数据库记录
- 错误日志: 应用日志文件

## 总结

本次架构升级实现了完整的三层数据架构，提供了：
- ✅ 实时数据监控
- ✅ 历史数据存储
- ✅ 高性能缓存
- ✅ 故障容错
- ✅ 扩展性设计
- ✅ 完整的管理界面

系统现在能够可靠地从俄罗斯中间件获取数据，存储历史记录，并通过缓存提供快速API响应。