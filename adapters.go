/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\adapters.go
 * @Description: 第三方日志库适配器和管理器
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

// AdapterType 适配器类型
type AdapterType string

const (
	StandardAdapter AdapterType = "standard" // 标准库适配器
	LogrusAdapter   AdapterType = "logrus"   // Logrus适配器 (需要额外安装)
	ZapAdapter      AdapterType = "zap"      // Zap适配器 (需要额外安装)
	ZerologAdapter  AdapterType = "zerolog"  // Zerolog适配器 (需要额外安装)
)

// AdapterConfig 适配器配置
type AdapterConfig struct {
	Type       AdapterType            `json:"type" yaml:"type"`
	Name       string                 `json:"name" yaml:"name"`
	Level      LogLevel               `json:"level" yaml:"level"`
	Output     io.Writer              `json:"-" yaml:"-"`
	File       string                 `json:"file,omitempty" yaml:"file,omitempty"`
	MaxSize    int                    `json:"max_size,omitempty" yaml:"max_size,omitempty"`       // MB
	MaxBackups int                    `json:"max_backups,omitempty" yaml:"max_backups,omitempty"`
	MaxAge     int                    `json:"max_age,omitempty" yaml:"max_age,omitempty"`          // days
	Compress   bool                   `json:"compress,omitempty" yaml:"compress,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty" yaml:"fields,omitempty"`
	Format     string                 `json:"format,omitempty" yaml:"format,omitempty"`           // json, text, custom
	TimeFormat string                 `json:"time_format,omitempty" yaml:"time_format,omitempty"`
	Colorful   bool                   `json:"colorful,omitempty" yaml:"colorful,omitempty"`
}

// DefaultAdapterConfig 创建默认适配器配置
func DefaultAdapterConfig() *AdapterConfig {
	return &AdapterConfig{
		Type:       StandardAdapter,
		Name:       "default",
		Level:      INFO,
		Output:     os.Stdout,
		Format:     "text",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
		Fields:     make(map[string]interface{}),
		TimeFormat: "2006-01-02 15:04:05",
		Colorful:   true,
	}
}

// Validate 验证适配器配置
func (c *AdapterConfig) Validate() error {
	if c.Name == "" {
		return errors.New("adapter name is required")
	}
	
	if c.Output == nil {
		c.Output = os.Stdout
	}
	
	if c.Fields == nil {
		c.Fields = make(map[string]interface{})
	}
	
	if c.TimeFormat == "" {
		c.TimeFormat = "2006-01-02 15:04:05"
	}
	
	return nil
}

// LoggerManager 日志管理器实现
type LoggerManager struct {
	adapters map[string]IAdapter
	mutex    sync.RWMutex
	config   *LogConfig
}

// NewLoggerManager 创建日志管理器
func NewLoggerManager() IManager {
	return &LoggerManager{
		adapters: make(map[string]IAdapter),
		config:   DefaultConfig(),
	}
}

// AddAdapter 添加适配器
func (lm *LoggerManager) AddAdapter(name string, adapter IAdapter) error {
	if name == "" {
		return errors.New("adapter name cannot be empty")
	}
	
	if adapter == nil {
		return errors.New("adapter cannot be nil")
	}
	
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	
	// 检查是否已存在
	if _, exists := lm.adapters[name]; exists {
		return fmt.Errorf("adapter with name '%s' already exists", name)
	}
	
	// 初始化适配器
	if err := adapter.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize adapter '%s': %w", name, err)
	}
	
	lm.adapters[name] = adapter
	return nil
}

// GetAdapter 获取适配器
func (lm *LoggerManager) GetAdapter(name string) (IAdapter, bool) {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	
	adapter, exists := lm.adapters[name]
	return adapter, exists
}

// RemoveAdapter 移除适配器
func (lm *LoggerManager) RemoveAdapter(name string) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	
	adapter, exists := lm.adapters[name]
	if !exists {
		return fmt.Errorf("adapter '%s' not found", name)
	}
	
	// 关闭适配器
	if err := adapter.Close(); err != nil {
		return fmt.Errorf("failed to close adapter '%s': %w", name, err)
	}
	
	delete(lm.adapters, name)
	return nil
}

// ListAdapters 列出所有适配器名称
func (lm *LoggerManager) ListAdapters() []string {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	
	names := make([]string, 0, len(lm.adapters))
	for name := range lm.adapters {
		names = append(names, name)
	}
	
	return names
}

// CloseAll 关闭所有适配器
func (lm *LoggerManager) CloseAll() error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	
	var errs []error
	for name, adapter := range lm.adapters {
		if err := adapter.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close adapter '%s': %w", name, err))
		}
	}
	
	// 清空适配器
	lm.adapters = make(map[string]IAdapter)
	
	if len(errs) > 0 {
		return fmt.Errorf("multiple errors occurred while closing adapters: %v", errs)
	}
	
	return nil
}

// FlushAll 刷新所有适配器
func (lm *LoggerManager) FlushAll() error {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	
	var errs []error
	for name, adapter := range lm.adapters {
		if err := adapter.Flush(); err != nil {
			errs = append(errs, fmt.Errorf("failed to flush adapter '%s': %w", name, err))
		}
	}
	
	if len(errs) > 0 {
		return fmt.Errorf("multiple errors occurred while flushing adapters: %v", errs)
	}
	
	return nil
}

// SetLevelAll 设置所有适配器的日志级别
func (lm *LoggerManager) SetLevelAll(level LogLevel) {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	
	for _, adapter := range lm.adapters {
		adapter.SetLevel(level)
	}
}

// Broadcast 广播日志到所有适配器
func (lm *LoggerManager) Broadcast(level LogLevel, format string, args ...interface{}) {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	
	for _, adapter := range lm.adapters {
		if !adapter.IsLevelEnabled(level) {
			continue
		}
		
		switch level {
		case DEBUG:
			adapter.Debug(format, args...)
		case INFO:
			adapter.Info(format, args...)
		case WARN:
			adapter.Warn(format, args...)
		case ERROR:
			adapter.Error(format, args...)
		case FATAL:
			adapter.Fatal(format, args...)
		}
	}
}

// HealthCheck 检查所有适配器的健康状态
func (lm *LoggerManager) HealthCheck() map[string]bool {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	
	health := make(map[string]bool)
	for name, adapter := range lm.adapters {
		health[name] = adapter.IsHealthy()
	}
	
	return health
}

// CreateAdapter 创建适配器的工厂方法
func CreateAdapter(config *AdapterConfig) (IAdapter, error) {
	if config == nil {
		config = DefaultAdapterConfig()
	}
	
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid adapter config: %w", err)
	}
	
	switch config.Type {
	case StandardAdapter:
		return NewStandardAdapter(config)
	case LogrusAdapter:
		return nil, fmt.Errorf("logrus adapter requires additional dependency: github.com/sirupsen/logrus")
	case ZapAdapter:
		return nil, fmt.Errorf("zap adapter requires additional dependency: go.uber.org/zap")
	case ZerologAdapter:
		return nil, fmt.Errorf("zerolog adapter requires additional dependency: github.com/rs/zerolog")
	default:
		return nil, fmt.Errorf("unsupported adapter type: %s", config.Type)
	}
}