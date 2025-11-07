/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 23:11:53
 * @FilePath: \go-logger\standard_adapter.go
 * @Description: 标准库日志适配器实现
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"fmt"
	"strings"
	"time"
)

// StandardLoggerAdapter 标准库适配器实现
type StandardLoggerAdapter struct {
	logger   *Logger
	config   *AdapterConfig
	name     string
	version  string
	healthy  bool
	initTime time.Time
}

// NewStandardAdapter 创建标准库适配器
func NewStandardAdapter(config *AdapterConfig) (IAdapter, error) {
	if config == nil {
		config = DefaultAdapterConfig()
	}
	
	// 获取前缀
	prefix := ""
	if prefixVal, ok := config.Fields["prefix"]; ok {
		if prefixStr, ok := prefixVal.(string); ok {
			prefix = prefixStr
		}
	}
	
	// 创建Logger配置
	logConfig := &LogConfig{
		Level:      config.Level,
		ShowCaller: true,
		Prefix:     prefix,
		Output:     config.Output,
		Colorful:   config.Colorful,
		TimeFormat: config.TimeFormat,
	}
	
	// 添加字段前缀
	if len(config.Fields) > 0 {
		var prefixParts []string
		for key, value := range config.Fields {
			if key != "prefix" {
				prefixParts = append(prefixParts, fmt.Sprintf("%s=%v", key, value))
			}
		}
		if len(prefixParts) > 0 {
			existingPrefix := logConfig.Prefix
			if existingPrefix != "" {
				existingPrefix += " "
			}
			logConfig.Prefix = existingPrefix + strings.Join(prefixParts, " ") + " "
		}
	}
	
	adapter := &StandardLoggerAdapter{
		logger:   NewLogger(logConfig),
		config:   config,
		name:     config.Name,
		version:  "1.0.0",
		healthy:  true,
		initTime: time.Now(),
	}
	
	return adapter, nil
}

// Initialize 初始化适配器
func (s *StandardLoggerAdapter) Initialize() error {
	s.healthy = true
	s.initTime = time.Now()
	return nil
}

// Close 关闭适配器
func (s *StandardLoggerAdapter) Close() error {
	s.healthy = false
	return nil
}

// Flush 刷新缓冲区
func (s *StandardLoggerAdapter) Flush() error {
	// 标准库的log包没有缓冲区，所以直接返回
	return nil
}

// GetAdapterName 获取适配器名称
func (s *StandardLoggerAdapter) GetAdapterName() string {
	return s.name
}

// GetAdapterVersion 获取适配器版本
func (s *StandardLoggerAdapter) GetAdapterVersion() string {
	return s.version
}

// IsHealthy 检查适配器是否健康
func (s *StandardLoggerAdapter) IsHealthy() bool {
	return s.healthy
}

// Debug 调试日志
func (s *StandardLoggerAdapter) Debug(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Debug(format, args...)
	}
}

// Info 信息日志
func (s *StandardLoggerAdapter) Info(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Info(format, args...)
	}
}

// Warn 警告日志
func (s *StandardLoggerAdapter) Warn(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Warn(format, args...)
	}
}

// Error 错误日志
func (s *StandardLoggerAdapter) Error(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Error(format, args...)
	}
}

// Fatal 致命错误日志
func (s *StandardLoggerAdapter) Fatal(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Fatal(format, args...)
	}
}

// SetLevel 设置日志级别
func (s *StandardLoggerAdapter) SetLevel(level LogLevel) {
	if s.healthy {
		s.logger.SetLevel(level)
	}
}

// GetLevel 获取日志级别
func (s *StandardLoggerAdapter) GetLevel() LogLevel {
	if s.healthy {
		return s.logger.GetLevel()
	}
	return ERROR
}

// SetShowCaller 设置是否显示调用者信息
func (s *StandardLoggerAdapter) SetShowCaller(show bool) {
	if s.healthy {
		s.logger.SetShowCaller(show)
	}
}

// IsShowCaller 是否显示调用者信息
func (s *StandardLoggerAdapter) IsShowCaller() bool {
	if s.healthy {
		return s.logger.IsShowCaller()
	}
	return false
}

// IsLevelEnabled 检查日志级别是否启用
func (s *StandardLoggerAdapter) IsLevelEnabled(level LogLevel) bool {
	if !s.healthy {
		return false
	}
	return s.logger.IsLevelEnabled(level)
}

// WithField 添加字段
func (s *StandardLoggerAdapter) WithField(key string, value interface{}) ILogger {
	if s.healthy {
		return &StandardLoggerAdapter{
			logger:   s.logger.WithField(key, value),
			config:   s.config,
			name:     s.name,
			version:  s.version,
			healthy:  s.healthy,
			initTime: s.initTime,
		}
	}
	return s
}

// WithFields 添加多个字段
func (s *StandardLoggerAdapter) WithFields(fields map[string]interface{}) ILogger {
	if s.healthy {
		return &StandardLoggerAdapter{
			logger:   s.logger.WithFields(fields),
			config:   s.config,
			name:     s.name,
			version:  s.version,
			healthy:  s.healthy,
			initTime: s.initTime,
		}
	}
	return s
}

// WithError 添加错误信息
func (s *StandardLoggerAdapter) WithError(err error) ILogger {
	if s.healthy {
		return &StandardLoggerAdapter{
			logger:   s.logger.WithError(err),
			config:   s.config,
			name:     s.name,
			version:  s.version,
			healthy:  s.healthy,
			initTime: s.initTime,
		}
	}
	return s
}

// Clone 克隆适配器
func (s *StandardLoggerAdapter) Clone() ILogger {
	return &StandardLoggerAdapter{
		logger:   s.logger.Clone(),
		config:   s.config,
		name:     s.name,
		version:  s.version,
		healthy:  s.healthy,
		initTime: s.initTime,
	}
}