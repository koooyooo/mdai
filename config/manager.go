package config

import "sync"

type Manager struct {
	config *Config
	mu     sync.RWMutex
}

var (
	instance *Manager
	once     sync.Once
)

// GetInstance returns the singleton instance of the config manager
func GetInstance() *Manager {
	once.Do(func() {
		instance = &Manager{}
		// 初期化時に設定を読み込む
		if cfg, err := LoadConfig(); err == nil {
			instance.config = cfg
		} else {
			instance.config = GetDefaultConfig()
		}
	})
	return instance
}

// GetConfig returns a copy of the current configuration
func (m *Manager) GetConfig() Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return *m.config
}

// ReloadConfig reloads the configuration from disk
func (m *Manager) ReloadConfig() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.config = cfg
	m.mu.Unlock()
	return nil
}
