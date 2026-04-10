package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/project47/cmd/mywebapp/types"
	_ "github.com/mattn/go-sqlite3" // SQLite驱动
)

// Store 数据库存储接口
type Store interface {
	// 设备管理
	SaveDevice(device types.Device) error
	GetDevice(deviceID string) (*types.Device, error)
	GetAllDevices() ([]types.Device, error)
	UpdateDevice(device types.Device) error
	DeleteDevice(deviceID string) error
	CountDevices() (int, error)

	// 试剂管理
	SaveReagents(deviceID string, reagents []types.Reagent) error
	GetReagents(deviceID string) ([]types.Reagent, error)
	UpdateReagent(deviceID string, reagent types.Reagent) error

	// 历史记录
	SaveDeviceStatusHistory(deviceID string, status types.DeviceStatus) error
	GetDeviceStatusHistory(deviceID string, limit int, offset int) ([]types.DeviceStatus, error)
	GetDeviceStatusHistoryByTimeRange(deviceID string, startTime, endTime time.Time) ([]types.DeviceStatus, error)
	SaveReagentConsumptionHistory(deviceID string, reagent types.Reagent, consumptionRate, estimatedHours *float64) error
	GetReagentConsumptionHistory(deviceID, reagentName string, limit int) ([]ReagentConsumptionHistoryModel, error)

	// API日志
	LogMiddlewareAPI(endpoint, method string, statusCode, responseTimeMs int, success bool, errorMessage string) error
	GetAPILogs(limit int) ([]MiddlewareAPILogModel, error)
	GetAPISuccessRate(duration time.Duration) (float64, error)

	// 数据源切换
	LogDataSourceSwitch(source, reason string, deviceCount int) error
	GetDataSourceHistory(limit int) ([]DataSourceHistoryModel, error)

	// 缓存统计
	UpdateCacheStatistics(cacheType string, hit bool) error
	GetCacheStatistics(cacheType string, period time.Duration) (*CacheStatisticsModel, error)

	// 统计查询
	GetDeviceStatusSummary() ([]DeviceStatusSummary, error)
	GetReagentConsumptionTrend(deviceID, reagentName string, hours int) ([]ReagentConsumptionTrend, error)
	GetOnlineDeviceCount() (int, error)
	GetDeviceWithLowReagentCount(threshold float64) (int, error)
	GetAverageAPIResponseTime(duration time.Duration) (float64, error)

	// 维护
	CleanupOldData(retentionDays int) error
	GetDatabaseStats() (map[string]interface{}, error)

	// 事务支持
	BeginTx() (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	RollbackTx(tx *sql.Tx) error

	// 关闭连接
	Close() error
}

// SQLStore SQL数据库存储实现
type SQLStore struct {
	db     *sql.DB
	config *DBConfig
	logger *log.Logger
}

// NewSQLStore 创建新的SQL存储
func NewSQLStore(config *DBConfig) (*SQLStore, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	dsn, err := config.GetDSN()
	if err != nil {
		return nil, fmt.Errorf("获取DSN失败: %v", err)
	}

	driverName := config.GetDriverName()
	if driverName == "" {
		return nil, fmt.Errorf("不支持的数据库类型: %s", config.Type)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接失败: %v", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("数据库连接测试失败: %v", err)
	}

	store := &SQLStore{
		db:     db,
		config: config,
		logger: log.New(log.Writer(), "[DB] ", log.LstdFlags),
	}

	// 执行迁移
	if config.AutoMigrate {
		if err := store.migrate(); err != nil {
			store.Close()
			return nil, fmt.Errorf("数据库迁移失败: %v", err)
		}
	}

	store.logger.Printf("数据库连接成功: %s", config.Type)
	return store, nil
}

// migrate 执行数据库迁移
func (s *SQLStore) migrate() error {
	// 这里简化处理，实际应该从文件读取schema.sql
	// 为了简单起见，我们直接执行一些基本的创建语句
	queries := []string{
		`CREATE TABLE IF NOT EXISTS devices (
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
		)`,
		`CREATE TABLE IF NOT EXISTS reagents (
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
		)`,
		`CREATE TABLE IF NOT EXISTS device_status_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id TEXT NOT NULL,
			status TEXT NOT NULL,
			is_online BOOLEAN NOT NULL,
			has_fault BOOLEAN NOT NULL,
			reagent_status TEXT CHECK (reagent_status IN ('normal', 'warning', 'low', 'empty')),
			uptime TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id)`,
		`CREATE INDEX IF NOT EXISTS idx_device_status_history_device_id_created_at ON device_status_history(device_id, created_at DESC)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("执行迁移语句失败: %v\n查询: %s", err, query)
		}
	}

	s.logger.Println("数据库迁移完成")
	return nil
}

// ==================== 设备管理 ====================

// SaveDevice 保存设备
func (s *SQLStore) SaveDevice(device types.Device) error {
	model := ToDeviceModel(device)

	query := `
		INSERT OR REPLACE INTO devices
		(id, device_id, name, status, ip, port, process_name, is_online, has_fault, last_check, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		model.ID, model.DeviceID, model.Name, model.Status, model.IP, model.Port,
		model.ProcessName, model.IsOnline, model.HasFault, model.LastCheck,
		model.CreatedAt, model.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("保存设备失败: %v", err)
	}

	// 保存试剂信息
	if len(device.Reagents) > 0 {
		if err := s.SaveReagents(device.DeviceID, device.Reagents); err != nil {
			return fmt.Errorf("保存试剂信息失败: %v", err)
		}
	}

	return nil
}

// GetDevice 获取设备
func (s *SQLStore) GetDevice(deviceID string) (*types.Device, error) {
	query := `SELECT * FROM devices WHERE device_id = ?`
	row := s.db.QueryRow(query, deviceID)

	var model DeviceModel
	err := row.Scan(
		&model.ID, &model.DeviceID, &model.Name, &model.Status, &model.IP, &model.Port,
		&model.ProcessName, &model.IsOnline, &model.HasFault, &model.LastCheck,
		&model.CreatedAt, &model.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("查询设备失败: %v", err)
	}

	// 获取试剂信息
	reagents, err := s.GetReagents(deviceID)
	if err != nil {
		return nil, fmt.Errorf("获取试剂信息失败: %v", err)
	}

	device := FromDeviceModel(model, reagents)
	return &device, nil
}

// GetAllDevices 获取所有设备
func (s *SQLStore) GetAllDevices() ([]types.Device, error) {
	query := `SELECT * FROM devices ORDER BY name`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询设备列表失败: %v", err)
	}
	defer rows.Close()

	var devices []types.Device
	for rows.Next() {
		var model DeviceModel
		err := rows.Scan(
			&model.ID, &model.DeviceID, &model.Name, &model.Status, &model.IP, &model.Port,
			&model.ProcessName, &model.IsOnline, &model.HasFault, &model.LastCheck,
			&model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描设备行失败: %v", err)
		}

		// 获取试剂信息
		reagents, err := s.GetReagents(model.DeviceID)
		if err != nil {
			s.logger.Printf("警告: 获取设备%s的试剂信息失败: %v", model.DeviceID, err)
			reagents = []types.Reagent{}
		}

		device := FromDeviceModel(model, reagents)
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历设备行失败: %v", err)
	}

	return devices, nil
}

// UpdateDevice 更新设备
func (s *SQLStore) UpdateDevice(device types.Device) error {
	model := ToDeviceModel(device)

	query := `
		UPDATE devices SET
			name = ?, status = ?, ip = ?, port = ?, process_name = ?,
			is_online = ?, has_fault = ?, last_check = ?, updated_at = ?
		WHERE device_id = ?
	`

	result, err := s.db.Exec(query,
		model.Name, model.Status, model.IP, model.Port, model.ProcessName,
		model.IsOnline, model.HasFault, model.LastCheck, model.UpdatedAt,
		model.DeviceID,
	)

	if err != nil {
		return fmt.Errorf("更新设备失败: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取受影响行数失败: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("设备不存在: %s", device.DeviceID)
	}

	// 更新试剂信息
	if len(device.Reagents) > 0 {
		if err := s.SaveReagents(device.DeviceID, device.Reagents); err != nil {
			return fmt.Errorf("更新试剂信息失败: %v", err)
		}
	}

	return nil
}

// DeleteDevice 删除设备
func (s *SQLStore) DeleteDevice(deviceID string) error {
	query := `DELETE FROM devices WHERE device_id = ?`
	result, err := s.db.Exec(query, deviceID)
	if err != nil {
		return fmt.Errorf("删除设备失败: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取受影响行数失败: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("设备不存在: %s", deviceID)
	}

	return nil
}

// CountDevices 统计设备数量
func (s *SQLStore) CountDevices() (int, error) {
	query := `SELECT COUNT(*) FROM devices`
	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("统计设备数量失败: %v", err)
	}
	return count, nil
}

// ==================== 试剂管理 ====================

// SaveReagents 保存试剂列表
func (s *SQLStore) SaveReagents(deviceID string, reagents []types.Reagent) error {
	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 删除旧的试剂记录
	deleteQuery := `DELETE FROM reagents WHERE device_id = ?`
	if _, err := tx.Exec(deleteQuery, deviceID); err != nil {
		return fmt.Errorf("删除旧试剂记录失败: %v", err)
	}

	// 插入新的试剂记录
	insertQuery := `
		INSERT INTO reagents (device_id, name, current, capacity, unit, percent, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, reagent := range reagents {
		model := ReagentModel{
			DeviceID:  deviceID,
			Name:      reagent.Name,
			Current:   reagent.Current,
			Capacity:  reagent.Capacity,
			Unit:      reagent.Unit,
			Percent:   reagent.Percent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if _, err := tx.Exec(insertQuery,
			model.DeviceID, model.Name, model.Current, model.Capacity,
			model.Unit, model.Percent, model.CreatedAt, model.UpdatedAt,
		); err != nil {
			return fmt.Errorf("插入试剂记录失败: %v", err)
		}
	}

	return tx.Commit()
}

// GetReagents 获取试剂列表
func (s *SQLStore) GetReagents(deviceID string) ([]types.Reagent, error) {
	query := `SELECT * FROM reagents WHERE device_id = ? ORDER BY name`
	rows, err := s.db.Query(query, deviceID)
	if err != nil {
		return nil, fmt.Errorf("查询试剂列表失败: %v", err)
	}
	defer rows.Close()

	var models []ReagentModel
	for rows.Next() {
		var model ReagentModel
		err := rows.Scan(
			&model.ID, &model.DeviceID, &model.Name, &model.Current, &model.Capacity,
			&model.Unit, &model.Percent, &model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描试剂行失败: %v", err)
		}
		models = append(models, model)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历试剂行失败: %v", err)
	}

	return FromReagentModels(models), nil
}

// UpdateReagent 更新试剂
func (s *SQLStore) UpdateReagent(deviceID string, reagent types.Reagent) error {
	query := `
		UPDATE reagents SET
			current = ?, capacity = ?, unit = ?, percent = ?, updated_at = ?
		WHERE device_id = ? AND name = ?
	`

	result, err := s.db.Exec(query,
		reagent.Current, reagent.Capacity, reagent.Unit, reagent.Percent, time.Now(),
		deviceID, reagent.Name,
	)

	if err != nil {
		return fmt.Errorf("更新试剂失败: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取受影响行数失败: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("试剂不存在: %s/%s", deviceID, reagent.Name)
	}

	return nil
}

// ==================== 历史记录 ====================

// SaveDeviceStatusHistory 保存设备状态历史
func (s *SQLStore) SaveDeviceStatusHistory(deviceID string, status types.DeviceStatus) error {
	model := ToDeviceStatusHistoryModel(deviceID, status)

	query := `
		INSERT INTO device_status_history
		(device_id, status, is_online, has_fault, reagent_status, uptime, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		model.DeviceID, model.Status, model.IsOnline, model.HasFault,
		model.ReagentStatus, model.Uptime, model.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("保存设备状态历史失败: %v", err)
	}

	return nil
}

// GetDeviceStatusHistory 获取设备状态历史
func (s *SQLStore) GetDeviceStatusHistory(deviceID string, limit int, offset int) ([]types.DeviceStatus, error) {
	if limit <= 0 {
		limit = 100
	}

	query := `
		SELECT * FROM device_status_history
		WHERE device_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.Query(query, deviceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询设备状态历史失败: %v", err)
	}
	defer rows.Close()

	var statuses []types.DeviceStatus
	for rows.Next() {
		var model DeviceStatusHistoryModel
		err := rows.Scan(
			&model.ID, &model.DeviceID, &model.Status, &model.IsOnline,
			&model.HasFault, &model.ReagentStatus, &model.Uptime, &model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描状态历史行失败: %v", err)
		}

		status := types.DeviceStatus{
			DeviceID:      model.DeviceID,
			DeviceName:    "", // 需要从设备表获取
			Status:        model.Status,
			IsOnline:      model.IsOnline,
			HasFault:      model.HasFault,
			ReagentStatus: model.ReagentStatus,
			Reagents:      []types.Reagent{}, // 需要从试剂表获取
			LastCheck:     model.CreatedAt.Format(time.RFC3339),
			Uptime:        model.Uptime,
		}
		statuses = append(statuses, status)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历状态历史行失败: %v", err)
	}

	return statuses, nil
}

// SaveReagentConsumptionHistory 保存试剂消耗历史
func (s *SQLStore) SaveReagentConsumptionHistory(deviceID string, reagent types.Reagent, consumptionRate, estimatedHours *float64) error {
	model := ToReagentConsumptionHistoryModel(deviceID, reagent, consumptionRate, estimatedHours)

	query := `
		INSERT INTO reagent_consumption_history
		(device_id, reagent_name, current, capacity, percent, consumption_rate, estimated_remaining_hours, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		model.DeviceID, model.ReagentName, model.Current, model.Capacity,
		model.Percent, model.ConsumptionRate, model.EstimatedRemainingHours, model.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("保存试剂消耗历史失败: %v", err)
	}

	return nil
}

// ==================== 其他方法实现 ====================

// LogMiddlewareAPI 记录中间件API调用日志
func (s *SQLStore) LogMiddlewareAPI(endpoint, method string, statusCode, responseTimeMs int, success bool, errorMessage string) error {
	query := `
		INSERT INTO middleware_api_logs
		(endpoint, method, status_code, response_time_ms, success, error_message, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		endpoint, method, statusCode, responseTimeMs, success, errorMessage, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("记录API日志失败: %v", err)
	}

	return nil
}

// LogDataSourceSwitch 记录数据源切换
func (s *SQLStore) LogDataSourceSwitch(source, reason string, deviceCount int) error {
	query := `
		INSERT INTO data_source_history
		(source, reason, device_count, created_at)
		VALUES (?, ?, ?, ?)
	`

	_, err := s.db.Exec(query, source, reason, deviceCount, time.Now())
	if err != nil {
		return fmt.Errorf("记录数据源切换失败: %v", err)
	}

	return nil
}

// UpdateCacheStatistics 更新缓存统计
func (s *SQLStore) UpdateCacheStatistics(cacheType string, hit bool) error {
	// 这里简化实现，实际应该使用更复杂的统计逻辑
	s.logger.Printf("缓存统计: type=%s, hit=%v", cacheType, hit)
	return nil
}

// Close 关闭数据库连接
func (s *SQLStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// 其他方法的简化实现（为了保持代码简洁）
func (s *SQLStore) GetDeviceStatusHistoryByTimeRange(deviceID string, startTime, endTime time.Time) ([]types.DeviceStatus, error) {
	return nil, fmt.Errorf("未实现")
}

func (s *SQLStore) GetReagentConsumptionHistory(deviceID, reagentName string, limit int) ([]ReagentConsumptionHistoryModel, error) {
	return nil, fmt.Errorf("未实现")
}

func (s *SQLStore) GetAPILogs(limit int) ([]MiddlewareAPILogModel, error) {
	return nil, fmt.Errorf("未实现")
}

func (s *SQLStore) GetAPISuccessRate(duration time.Duration) (float64, error) {
	return 0, fmt.Errorf("未实现")
}

func (s *SQLStore) GetDataSourceHistory(limit int) ([]DataSourceHistoryModel, error) {
	return nil, fmt.Errorf("未实现")
}

func (s *SQLStore) GetCacheStatistics(cacheType string, period time.Duration) (*CacheStatisticsModel, error) {
	return nil, fmt.Errorf("未实现")
}

func (s *SQLStore) GetDeviceStatusSummary() ([]DeviceStatusSummary, error) {
	return nil, fmt.Errorf("未实现")
}

func (s *SQLStore) GetReagentConsumptionTrend(deviceID, reagentName string, hours int) ([]ReagentConsumptionTrend, error) {
	return nil, fmt.Errorf("未实现")
}

func (s *SQLStore) GetOnlineDeviceCount() (int, error) {
	query := `SELECT COUNT(*) FROM devices WHERE is_online = 1`
	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("查询在线设备数量失败: %v", err)
	}
	return count, nil
}

func (s *SQLStore) GetDeviceWithLowReagentCount(threshold float64) (int, error) {
	// 这里简化实现
	return 0, fmt.Errorf("未实现")
}

func (s *SQLStore) GetAverageAPIResponseTime(duration time.Duration) (float64, error) {
	return 0, fmt.Errorf("未实现")
}

func (s *SQLStore) CleanupOldData(retentionDays int) error {
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	// 清理设备状态历史
	query1 := `DELETE FROM device_status_history WHERE created_at < ?`
	if _, err := s.db.Exec(query1, cutoffTime); err != nil {
		return fmt.Errorf("清理设备状态历史失败: %v", err)
	}

	// 清理试剂消耗历史
	query2 := `DELETE FROM reagent_consumption_history WHERE created_at < ?`
	if _, err := s.db.Exec(query2, cutoffTime); err != nil {
		return fmt.Errorf("清理试剂消耗历史失败: %v", err)
	}

	// 清理API日志
	query3 := `DELETE FROM middleware_api_logs WHERE created_at < ?`
	if _, err := s.db.Exec(query3, cutoffTime); err != nil {
		return fmt.Errorf("清理API日志失败: %v", err)
	}

	s.logger.Printf("清理了%d天前的历史数据", retentionDays)
	return nil
}

func (s *SQLStore) GetDatabaseStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取表行数
	tables := []string{"devices", "reagents", "device_status_history", "reagent_consumption_history", "middleware_api_logs"}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		var count int
		if err := s.db.QueryRow(query).Scan(&count); err == nil {
			stats[table+"_count"] = count
		}
	}

	// 获取数据库大小（SQLite）
	if s.config.Type == "sqlite" {
		var pageSize, pageCount, freelistCount int
		query := `PRAGMA page_size; PRAGMA page_count; PRAGMA freelist_count;`
		rows, err := s.db.Query(query)
		if err == nil {
			defer rows.Close()
			if rows.Next() {
				rows.Scan(&pageSize)
			}
			if rows.Next() {
				rows.Scan(&pageCount)
			}
			if rows.Next() {
				rows.Scan(&freelistCount)
			}
			stats["database_size_mb"] = float64(pageSize*pageCount) / (1024 * 1024)
			stats["free_pages"] = freelistCount
		}
	}

	stats["database_type"] = s.config.Type
	stats["connected"] = true

	return stats, nil
}

func (s *SQLStore) BeginTx() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *SQLStore) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (s *SQLStore) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}