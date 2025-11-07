/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 23:18:44
 * @FilePath: \go-logger\logger.go
 * @Description: 统一的日志工具包，支持 emoji 和结构化日志
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// defaultLogger 默认日志记录器
var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(DefaultConfig())
}

// NewLogger 创建新的日志记录器
func NewLogger(config *LogConfig) *Logger {
	if config == nil {
		config = DefaultConfig()
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		config = DefaultConfig()
	}

	prefix := config.Prefix
	if prefix != "" && !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}

	return &Logger{
		level:      config.Level,
		showCaller: config.ShowCaller,
		logger:     log.New(config.Output, prefix, log.LstdFlags),
		config:     config.Clone(),
	}
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() LogLevel {
	return l.level
}

// SetShowCaller 设置是否显示调用者信息
func (l *Logger) SetShowCaller(show bool) {
	l.showCaller = show
}

// IsShowCaller 检查是否显示调用者信息
func (l *Logger) IsShowCaller() bool {
	return l.showCaller
}

// IsLevelEnabled 检查给定级别是否启用
func (l *Logger) IsLevelEnabled(level LogLevel) bool {
	return level >= l.level
}

// GetConfig 获取日志配置的副本
func (l *Logger) GetConfig() *LogConfig {
	return l.config.Clone()
}

// UpdateConfig 更新日志配置
func (l *Logger) UpdateConfig(config *LogConfig) {
	if config == nil {
		return
	}
	
	l.config = config.Clone()
	l.level = config.Level
	l.showCaller = config.ShowCaller
	
	// 更新内部logger
	prefix := config.Prefix
	if prefix != "" && !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}
	l.logger = log.New(config.Output, prefix, log.LstdFlags)
}

// formatMessage 格式化消息
func (l *Logger) formatMessage(level LogLevel, format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	
	// 构建级别信息
	levelStr := fmt.Sprintf("%s [%s]", level.Emoji(), level.String())
	if l.config.Colorful {
		levelStr = level.Color() + levelStr + "\033[0m"
	}
	
	var callerInfo string
	if l.showCaller {
		if pc, file, line, ok := runtime.Caller(3); ok {
			funcName := runtime.FuncForPC(pc).Name()
			if idx := strings.LastIndex(funcName, "."); idx != -1 {
				funcName = funcName[idx+1:]
			}
			if idx := strings.LastIndex(file, "/"); idx != -1 {
				file = file[idx+1:]
			}
			callerInfo = fmt.Sprintf(" [%s:%d:%s]", file, line, funcName)
		}
	}

	return fmt.Sprintf("%s%s %s", levelStr, callerInfo, message)
}

// log 记录日志
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	message := l.formatMessage(level, format, args...)
	l.logger.Print(message)

	if level == FATAL {
		os.Exit(1)
	}
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// WithField 添加字段信息（结构化日志）
func (l *Logger) WithField(key string, value interface{}) *Logger {
	prefix := fmt.Sprintf("%s%s=%v ", l.config.Prefix, key, value)
	config := l.config.Clone()
	config.Prefix = prefix
	
	return NewLogger(config)
}

// WithFields 添加多个字段信息（结构化日志）
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	if len(fields) == 0 {
		return l
	}
	
	var prefix strings.Builder
	prefix.WriteString(l.config.Prefix)
	
	for key, value := range fields {
		prefix.WriteString(fmt.Sprintf("%s=%v ", key, value))
	}
	
	config := l.config.Clone()
	config.Prefix = prefix.String()
	
	return NewLogger(config)
}

// WithError 添加错误信息
func (l *Logger) WithError(err error) *Logger {
	return l.WithField("error", err.Error())
}

// Clone 克隆当前Logger
func (l *Logger) Clone() *Logger {
	// 创建新的配置反映当前logger状态
	newConfig := l.config.Clone()
	newConfig.Level = l.level
	newConfig.ShowCaller = l.showCaller
	return NewLogger(newConfig)
}

// 全局方法
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// SetGlobalLevel 设置全局日志级别
func SetGlobalLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// SetGlobalShowCaller 设置全局是否显示调用者信息
func SetGlobalShowCaller(show bool) {
	defaultLogger.SetShowCaller(show)
}

// GetGlobalLogger 获取全局Logger
func GetGlobalLogger() *Logger {
	return defaultLogger
}

// WithField 全局添加字段信息
func WithField(key string, value interface{}) *Logger {
	return defaultLogger.WithField(key, value)
}

// WithFields 全局添加多个字段信息
func WithFields(fields map[string]interface{}) *Logger {
	return defaultLogger.WithFields(fields)
}

// WithError 全局添加错误信息
func WithError(err error) *Logger {
	return defaultLogger.WithError(err)
}

// SetGlobalConfig 设置全局配置
func SetGlobalConfig(config *LogConfig) {
	defaultLogger.UpdateConfig(config)
}

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *LogConfig {
	return defaultLogger.GetConfig()
}