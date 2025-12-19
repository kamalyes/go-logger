/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 16:30:00
 * @FilePath: \go-logger\standard_adapter.go
 * @Description: 标准库日志适配器实现
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// StandardLoggerAdapter 标准库适配器实现
type StandardLoggerAdapter struct {
	logger   ILogger
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
		ShowCaller: config.ShowCaller,
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

// Printf风格方法（与上面相同，但命名更明确）
func (s *StandardLoggerAdapter) Debugf(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Debug(format, args...)
	}
}

func (s *StandardLoggerAdapter) Infof(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Info(format, args...)
	}
}

func (s *StandardLoggerAdapter) Warnf(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Warn(format, args...)
	}
}

func (s *StandardLoggerAdapter) Errorf(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Error(format, args...)
	}
}

func (s *StandardLoggerAdapter) Fatalf(format string, args ...interface{}) {
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

// 为 StandardLoggerAdapter 添加新接口方法的实现

// 纯文本日志方法
func (s *StandardLoggerAdapter) DebugMsg(msg string) {
	if s.healthy {
		s.logger.DebugMsg(msg)
	}
}

func (s *StandardLoggerAdapter) InfoMsg(msg string) {
	if s.healthy {
		s.logger.InfoMsg(msg)
	}
}

func (s *StandardLoggerAdapter) WarnMsg(msg string) {
	if s.healthy {
		s.logger.WarnMsg(msg)
	}
}

func (s *StandardLoggerAdapter) ErrorMsg(msg string) {
	if s.healthy {
		s.logger.ErrorMsg(msg)
	}
}

func (s *StandardLoggerAdapter) FatalMsg(msg string) {
	if s.healthy {
		s.logger.FatalMsg(msg)
	}
}

// 多行日志方法
func (s *StandardLoggerAdapter) InfoLines(lines ...string) {
	if s.healthy {
		s.logger.InfoLines(lines...)
	}
}

func (s *StandardLoggerAdapter) ErrorLines(lines ...string) {
	if s.healthy {
		s.logger.ErrorLines(lines...)
	}
}

func (s *StandardLoggerAdapter) WarnLines(lines ...string) {
	if s.healthy {
		s.logger.WarnLines(lines...)
	}
}

func (s *StandardLoggerAdapter) DebugLines(lines ...string) {
	if s.healthy {
		s.logger.DebugLines(lines...)
	}
}

// 带上下文的日志方法
func (s *StandardLoggerAdapter) DebugContext(ctx context.Context, format string, args ...interface{}) {
	if s.healthy {
		s.logger.DebugContext(ctx, format, args...)
	}
}

func (s *StandardLoggerAdapter) InfoContext(ctx context.Context, format string, args ...interface{}) {
	if s.healthy {
		s.logger.InfoContext(ctx, format, args...)
	}
}

func (s *StandardLoggerAdapter) WarnContext(ctx context.Context, format string, args ...interface{}) {
	if s.healthy {
		s.logger.WarnContext(ctx, format, args...)
	}
}

func (s *StandardLoggerAdapter) ErrorContext(ctx context.Context, format string, args ...interface{}) {
	if s.healthy {
		s.logger.ErrorContext(ctx, format, args...)
	}
}

func (s *StandardLoggerAdapter) FatalContext(ctx context.Context, format string, args ...interface{}) {
	if s.healthy {
		s.logger.FatalContext(ctx, format, args...)
	}
}

// 结构化日志方法（键值对）
func (s *StandardLoggerAdapter) DebugKV(msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.DebugKV(msg, keysAndValues...)
	}
}

func (s *StandardLoggerAdapter) InfoKV(msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.InfoKV(msg, keysAndValues...)
	}
}

func (s *StandardLoggerAdapter) WarnKV(msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.WarnKV(msg, keysAndValues...)
	}
}

func (s *StandardLoggerAdapter) ErrorKV(msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.ErrorKV(msg, keysAndValues...)
	}
}

func (s *StandardLoggerAdapter) FatalKV(msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.FatalKV(msg, keysAndValues...)
	}
}

// 带上下文的结构化日志方法
func (s *StandardLoggerAdapter) DebugContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.DebugContextKV(ctx, msg, keysAndValues...)
	}
}

func (s *StandardLoggerAdapter) InfoContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.InfoContextKV(ctx, msg, keysAndValues...)
	}
}

func (s *StandardLoggerAdapter) WarnContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.WarnContextKV(ctx, msg, keysAndValues...)
	}
}

func (s *StandardLoggerAdapter) ErrorContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.ErrorContextKV(ctx, msg, keysAndValues...)
	}
}

func (s *StandardLoggerAdapter) FatalContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.FatalContextKV(ctx, msg, keysAndValues...)
	}
}

// 原始日志条目方法
func (s *StandardLoggerAdapter) Log(level LogLevel, msg string) {
	if s.healthy {
		s.logger.Log(level, msg)
	}
}

func (s *StandardLoggerAdapter) LogContext(ctx context.Context, level LogLevel, msg string) {
	if s.healthy {
		s.logger.LogContext(ctx, level, msg)
	}
}

func (s *StandardLoggerAdapter) LogKV(level LogLevel, msg string, keysAndValues ...interface{}) {
	if s.healthy {
		s.logger.LogKV(level, msg, keysAndValues...)
	}
}

func (s *StandardLoggerAdapter) LogWithFields(level LogLevel, msg string, fields map[string]interface{}) {
	if s.healthy {
		s.logger.LogWithFields(level, msg, fields)
	}
}

// WithContext 的实现
func (s *StandardLoggerAdapter) WithContext(ctx context.Context) ILogger {
	if s.healthy {
		return &StandardLoggerAdapter{
			logger:   s.logger.WithContext(ctx),
			config:   s.config,
			name:     s.name,
			version:  s.version,
			healthy:  s.healthy,
			initTime: s.initTime,
		}
	}
	return s
}

// 兼容标准log包的方法
func (s *StandardLoggerAdapter) Print(args ...interface{}) {
	if s.healthy {
		s.logger.Print(args...)
	}
}

func (s *StandardLoggerAdapter) Printf(format string, args ...interface{}) {
	if s.healthy {
		s.logger.Printf(format, args...)
	}
}

func (s *StandardLoggerAdapter) Println(args ...interface{}) {
	if s.healthy {
		s.logger.Println(args...)
	}
}

// 返回错误的日志方法
func (s *StandardLoggerAdapter) DebugReturn(format string, args ...interface{}) error {
	if s.healthy {
		return s.logger.DebugReturn(format, args...)
	}
	return fmt.Errorf(format, args...)
}

func (s *StandardLoggerAdapter) InfoReturn(format string, args ...interface{}) error {
	if s.healthy {
		return s.logger.InfoReturn(format, args...)
	}
	return fmt.Errorf(format, args...)
}

func (s *StandardLoggerAdapter) WarnReturn(format string, args ...interface{}) error {
	if s.healthy {
		return s.logger.WarnReturn(format, args...)
	}
	return fmt.Errorf(format, args...)
}

func (s *StandardLoggerAdapter) ErrorReturn(format string, args ...interface{}) error {
	if s.healthy {
		return s.logger.ErrorReturn(format, args...)
	}
	return fmt.Errorf(format, args...)
}

// 返回错误的上下文日志方法
func (s *StandardLoggerAdapter) DebugCtxReturn(ctx context.Context, format string, args ...interface{}) error {
	if s.healthy {
		return s.logger.DebugCtxReturn(ctx, format, args...)
	}
	return fmt.Errorf(format, args...)
}

func (s *StandardLoggerAdapter) InfoCtxReturn(ctx context.Context, format string, args ...interface{}) error {
	if s.healthy {
		return s.logger.InfoCtxReturn(ctx, format, args...)
	}
	return fmt.Errorf(format, args...)
}

func (s *StandardLoggerAdapter) WarnCtxReturn(ctx context.Context, format string, args ...interface{}) error {
	if s.healthy {
		return s.logger.WarnCtxReturn(ctx, format, args...)
	}
	return fmt.Errorf(format, args...)
}

func (s *StandardLoggerAdapter) ErrorCtxReturn(ctx context.Context, format string, args ...interface{}) error {
	if s.healthy {
		return s.logger.ErrorCtxReturn(ctx, format, args...)
	}
	return fmt.Errorf(format, args...)
}

// 返回错误的键值对日志方法
func (s *StandardLoggerAdapter) DebugKVReturn(msg string, keysAndValues ...interface{}) error {
	if s.healthy {
		return s.logger.DebugKVReturn(msg, keysAndValues...)
	}
	return fmt.Errorf("%s", msg)
}

func (s *StandardLoggerAdapter) InfoKVReturn(msg string, keysAndValues ...interface{}) error {
	if s.healthy {
		return s.logger.InfoKVReturn(msg, keysAndValues...)
	}
	return fmt.Errorf("%s", msg)
}

func (s *StandardLoggerAdapter) WarnKVReturn(msg string, keysAndValues ...interface{}) error {
	if s.healthy {
		return s.logger.WarnKVReturn(msg, keysAndValues...)
	}
	return fmt.Errorf("%s", msg)
}

func (s *StandardLoggerAdapter) ErrorKVReturn(msg string, keysAndValues ...interface{}) error {
	if s.healthy {
		return s.logger.ErrorKVReturn(msg, keysAndValues...)
	}
	return fmt.Errorf("%s", msg)
}
