package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DBConfig 数据库配置
type DBConfig struct {
	// 数据库类型
	Type string `json:"type"` // sqlite, postgres, mysql

	// SQLite配置
	SQLitePath string `json:"sqlite_path"`

	// PostgreSQL配置
	PostgresHost     string `json:"postgres_host"`
	PostgresPort     int    `json:"postgres_port"`
	PostgresUser     string `json:"postgres_user"`
	PostgresPassword string `json:"postgres_password"`
	PostgresDBName   string `json:"postgres_dbname"`
	PostgresSSLMode  string `json:"postgres_sslmode"`

	// 连接池配置
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`

	// 迁移配置
	AutoMigrate bool `json:"auto_migrate"`

	// 日志配置
	EnableQueryLog bool `json:"enable_query_log"`
}

// DefaultConfig 默认配置
func DefaultConfig() *DBConfig {
	// 获取可执行文件所在目录
	exeDir, err := os.Executable()
	var dbPath string
	if err == nil {
		exeDir = filepath.Dir(exeDir)
		dbPath = filepath.Join(exeDir, "data", "project47.db")
	} else {
		// 回退到当前目录
		dbPath = "./data/project47.db"
	}

	// 确保目录存在
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	return &DBConfig{
		Type:             "sqlite",
		SQLitePath:       dbPath,
		PostgresHost:     "localhost",
		PostgresPort:     5432,
		PostgresUser:     "postgres",
		PostgresPassword: "",
		PostgresDBName:   "project47",
		PostgresSSLMode:  "disable",
		MaxOpenConns:     10,
		MaxIdleConns:     5,
		ConnMaxLifetime:  time.Hour,
		ConnMaxIdleTime:  30 * time.Minute,
		AutoMigrate:      true,
		EnableQueryLog:   false,
	}
}

// Validate 验证配置
func (c *DBConfig) Validate() error {
	if c.Type == "" {
		return fmt.Errorf("数据库类型不能为空")
	}

	switch c.Type {
	case "sqlite":
		if c.SQLitePath == "" {
			return fmt.Errorf("SQLite路径不能为空")
		}
	case "postgres":
		if c.PostgresHost == "" {
			return fmt.Errorf("PostgreSQL主机不能为空")
		}
		if c.PostgresPort <= 0 {
			return fmt.Errorf("PostgreSQL端口无效")
		}
		if c.PostgresUser == "" {
			return fmt.Errorf("PostgreSQL用户不能为空")
		}
		if c.PostgresDBName == "" {
			return fmt.Errorf("PostgreSQL数据库名不能为空")
		}
	default:
		return fmt.Errorf("不支持的数据库类型: %s", c.Type)
	}

	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DBConfig) GetDSN() (string, error) {
	switch c.Type {
	case "sqlite":
		return c.SQLitePath, nil
	case "postgres":
		sslMode := c.PostgresSSLMode
		if sslMode == "" {
			sslMode = "disable"
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword,
			c.PostgresDBName, sslMode), nil
	default:
		return "", fmt.Errorf("不支持的数据库类型: %s", c.Type)
	}
}

// GetDriverName 获取数据库驱动名称
func (c *DBConfig) GetDriverName() string {
	switch c.Type {
	case "sqlite":
		return "sqlite3"
	case "postgres":
		return "postgres"
	default:
		return ""
	}
}