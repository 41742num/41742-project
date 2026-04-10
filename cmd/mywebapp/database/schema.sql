-- 项目47数据库架构
-- 支持俄罗斯中间件API > 数据库(历史数据) > 内存缓存API响应 架构

-- 设备表 - 存储设备基本信息
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    device_id TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('enabled', 'disabled')),
    ip TEXT NOT NULL,
    port INTEGER NOT NULL,
    process_name TEXT,
    is_online BOOLEAN DEFAULT FALSE,
    has_fault BOOLEAN DEFAULT FALSE,
    last_check TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 试剂表 - 存储设备试剂信息
CREATE TABLE IF NOT EXISTS reagents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id TEXT NOT NULL,
    name TEXT NOT NULL,
    current REAL NOT NULL,
    capacity REAL NOT NULL,
    unit TEXT NOT NULL,
    percent REAL NOT NULL CHECK (percent >= 0 AND percent <= 100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE,
    UNIQUE(device_id, name)
);

-- 设备状态历史表 - 记录设备状态变化历史
CREATE TABLE IF NOT EXISTS device_status_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id TEXT NOT NULL,
    status TEXT NOT NULL,
    is_online BOOLEAN NOT NULL,
    has_fault BOOLEAN NOT NULL,
    reagent_status TEXT CHECK (reagent_status IN ('normal', 'warning', 'low', 'empty')),
    uptime TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
);

-- 试剂消耗历史表 - 记录试剂消耗变化
CREATE TABLE IF NOT EXISTS reagent_consumption_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id TEXT NOT NULL,
    reagent_name TEXT NOT NULL,
    current REAL NOT NULL,
    capacity REAL NOT NULL,
    percent REAL NOT NULL CHECK (percent >= 0 AND percent <= 100),
    consumption_rate REAL, -- 消耗速率 (ml/hour)
    estimated_remaining_hours REAL, -- 预计剩余小时数
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
);

-- 中间件API调用日志表 - 记录API调用情况
CREATE TABLE IF NOT EXISTS middleware_api_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    endpoint TEXT NOT NULL,
    method TEXT NOT NULL,
    status_code INTEGER,
    response_time_ms INTEGER,
    success BOOLEAN NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 数据源切换历史表 - 记录数据源切换情况
CREATE TABLE IF NOT EXISTS data_source_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source TEXT NOT NULL CHECK (source IN ('middleware', 'database', 'fallback')),
    reason TEXT,
    device_count INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 缓存命中率统计表
CREATE TABLE IF NOT EXISTS cache_statistics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cache_type TEXT NOT NULL CHECK (cache_type IN ('device_list', 'device_status', 'device_stats')),
    hits INTEGER DEFAULT 0,
    misses INTEGER DEFAULT 0,
    total_requests INTEGER DEFAULT 0,
    hit_rate REAL DEFAULT 0,
    period_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    period_end TIMESTAMP
);

-- 索引优化
CREATE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id);
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_is_online ON devices(is_online);
CREATE INDEX IF NOT EXISTS idx_reagents_device_id ON reagents(device_id);
CREATE INDEX IF NOT EXISTS idx_device_status_history_device_id_created_at ON device_status_history(device_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_device_status_history_created_at ON device_status_history(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reagent_consumption_history_device_id_created_at ON reagent_consumption_history(device_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_middleware_api_logs_created_at ON middleware_api_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_data_source_history_created_at ON data_source_history(created_at DESC);

-- 触发器：更新设备时自动更新updated_at
CREATE TRIGGER IF NOT EXISTS update_devices_timestamp
AFTER UPDATE ON devices
BEGIN
    UPDATE devices SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- 触发器：更新试剂时自动更新updated_at
CREATE TRIGGER IF NOT EXISTS update_reagents_timestamp
AFTER UPDATE ON reagents
BEGIN
    UPDATE reagents SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- 视图：设备状态汇总视图
CREATE VIEW IF NOT EXISTS device_status_summary AS
SELECT
    d.device_id,
    d.name,
    d.status,
    d.is_online,
    d.has_fault,
    d.last_check,
    COALESCE(r.min_percent, 100) as min_reagent_percent,
    CASE
        WHEN COALESCE(r.min_percent, 100) = 0 THEN 'empty'
        WHEN COALESCE(r.min_percent, 100) < 30 THEN 'low'
        WHEN COALESCE(r.min_percent, 100) < 70 THEN 'warning'
        ELSE 'normal'
    END as reagent_status,
    (SELECT COUNT(*) FROM device_status_history WHERE device_id = d.device_id) as history_count,
    (SELECT created_at FROM device_status_history WHERE device_id = d.device_id ORDER BY created_at DESC LIMIT 1) as last_status_update
FROM devices d
LEFT JOIN (
    SELECT device_id, MIN(percent) as min_percent
    FROM reagents
    GROUP BY device_id
) r ON d.device_id = r.device_id;

-- 视图：试剂消耗趋势视图
CREATE VIEW IF NOT EXISTS reagent_consumption_trend AS
SELECT
    rch.device_id,
    rch.reagent_name,
    rch.created_at,
    rch.current,
    rch.capacity,
    rch.percent,
    rch.consumption_rate,
    rch.estimated_remaining_hours,
    LAG(rch.percent) OVER (PARTITION BY rch.device_id, rch.reagent_name ORDER BY rch.created_at) as previous_percent,
    rch.percent - LAG(rch.percent) OVER (PARTITION BY rch.device_id, rch.reagent_name ORDER BY rch.created_at) as percent_change
FROM reagent_consumption_history rch
ORDER BY rch.created_at DESC;

-- 初始化数据（可选）
-- INSERT OR IGNORE INTO devices (id, device_id, name, status, ip, port, process_name, is_online, has_fault)
-- VALUES
-- ('1', 'MIDDLEWARE_001', '俄罗斯中间件服务器', 'enabled', '172.19.14.202', 10001, 'middleware', 1, 0),
-- ('2', 'ANALYZER_001', '分析仪器1', 'enabled', '172.19.14.203', 10002, 'analyzer', 1, 0),
-- ('3', 'ANALYZER_002', '分析仪器2', 'disabled', '172.19.14.204', 10003, 'analyzer', 0, 0);