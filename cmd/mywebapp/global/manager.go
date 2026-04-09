package global

import (
	"fmt"
	"sync"

	"github.com/project47/cmd/mywebapp/data"
)

// GlobalManager 全局管理器
type GlobalManager struct {
	dataManager *data.DataManager
	initialized bool
	mutex       sync.RWMutex
}

var (
	instance *GlobalManager
	once     sync.Once
)

// GetInstance 获取全局管理器单例
func GetInstance() *GlobalManager {
	once.Do(func() {
		instance = &GlobalManager{}
	})
	return instance
}

// Initialize 初始化全局管理器
func (gm *GlobalManager) Initialize(config *data.Config) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if gm.initialized {
		return nil
	}

	if config == nil {
		config = data.DefaultConfig()
	}

	// 创建数据管理器
	dm, err := data.NewDataManager(config)
	if err != nil {
		return fmt.Errorf("初始化数据管理器失败: %v", err)
	}

	gm.dataManager = dm
	gm.initialized = true

	return nil
}

// GetDataManager 获取数据管理器
func (gm *GlobalManager) GetDataManager() (*data.DataManager, error) {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()

	if !gm.initialized {
		return nil, fmt.Errorf("全局管理器未初始化")
	}

	return gm.dataManager, nil
}

// IsInitialized 检查是否已初始化
func (gm *GlobalManager) IsInitialized() bool {
	gm.mutex.RLock()
	defer gm.mutex.RUnlock()
	return gm.initialized
}

// Shutdown 关闭全局管理器
func (gm *GlobalManager) Shutdown() {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	if gm.dataManager != nil {
		gm.dataManager.Stop()
		gm.dataManager = nil
	}

	gm.initialized = false
}

// Reload 重新加载配置
func (gm *GlobalManager) Reload(config *data.Config) error {
	gm.Shutdown()
	return gm.Initialize(config)
}