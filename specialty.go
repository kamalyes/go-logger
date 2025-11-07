/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\specialty.go
 * @Description: 特殊场景的日志方法
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"fmt"
	"time"
)

// SpecialLogType 特殊日志类型
type SpecialLogType struct {
	emoji string
	name  string
}

// 特殊日志类型定义
var (
	SuccessType     = SpecialLogType{"✅", "SUCCESS"}
	LoadingType     = SpecialLogType{"⏳", "LOADING"}
	ConfigType      = SpecialLogType{"⚙️", "CONFIG"}
	StartType       = SpecialLogType{"🚀", "START"}
	StopType        = SpecialLogType{"🛑", "STOP"}
	DatabaseType    = SpecialLogType{"💾", "DATABASE"}
	NetworkType     = SpecialLogType{"🌐", "NETWORK"}
	SecurityType    = SpecialLogType{"🔒", "SECURITY"}
	CacheType       = SpecialLogType{"🗄️", "CACHE"}
	EnvironmentType = SpecialLogType{"🌍", "ENV"}
)

// logSpecial 记录特殊类型的日志
func logSpecial(logType SpecialLogType, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	defaultLogger.logger.Printf("%s [%s] %s", logType.emoji, logType.name, message)
}

// Success 成功日志
func Success(format string, args ...interface{}) {
	logSpecial(SuccessType, format, args...)
}

// Loading 加载日志
func Loading(format string, args ...interface{}) {
	logSpecial(LoadingType, format, args...)
}

// ConfigLogger 配置日志
func ConfigLogger(format string, args ...interface{}) {
	logSpecial(ConfigType, format, args...)
}

// Start 启动日志
func Start(format string, args ...interface{}) {
	logSpecial(StartType, format, args...)
}

// Stop 停止日志
func Stop(format string, args ...interface{}) {
	logSpecial(StopType, format, args...)
}

// Database 数据库日志
func Database(format string, args ...interface{}) {
	logSpecial(DatabaseType, format, args...)
}

// Network 网络日志
func Network(format string, args ...interface{}) {
	logSpecial(NetworkType, format, args...)
}

// Security 安全日志
func Security(format string, args ...interface{}) {
	logSpecial(SecurityType, format, args...)
}

// Cache 缓存日志
func Cache(format string, args ...interface{}) {
	logSpecial(CacheType, format, args...)
}

// Environment 环境日志
func Environment(format string, args ...interface{}) {
	logSpecial(EnvironmentType, format, args...)
}

// Performance 性能日志
func Performance(operation string, duration time.Duration) {
	var emoji string
	var level string
	
	switch {
	case duration < 50*time.Millisecond:
		emoji = "⚡"
		level = "EXCELLENT"
	case duration < 100*time.Millisecond:
		emoji = "🏃"
		level = "FAST"
	case duration < 500*time.Millisecond:
		emoji = "🚶"
		level = "NORMAL"
	case duration < 2*time.Second:
		emoji = "🐢"
		level = "SLOW"
	default:
		emoji = "🐌"
		level = "VERY_SLOW"
	}
	
	defaultLogger.logger.Printf("%s [PERF-%s] %s completed in %v", emoji, level, operation, duration)
}

// PerformanceWithDetails 带详细信息的性能日志
func PerformanceWithDetails(operation string, duration time.Duration, details map[string]interface{}) {
	var emoji string
	var level string
	
	switch {
	case duration < 50*time.Millisecond:
		emoji = "⚡"
		level = "EXCELLENT"
	case duration < 100*time.Millisecond:
		emoji = "🏃"
		level = "FAST"
	case duration < 500*time.Millisecond:
		emoji = "🚶"
		level = "NORMAL"
	case duration < 2*time.Second:
		emoji = "🐢"
		level = "SLOW"
	default:
		emoji = "🐌"
		level = "VERY_SLOW"
	}
	
	detailStr := ""
	if len(details) > 0 {
		detailStr = fmt.Sprintf(" | Details: %+v", details)
	}
	
	defaultLogger.logger.Printf("%s [PERF-%s] %s completed in %v%s", 
		emoji, level, operation, duration, detailStr)
}

// Timing 计时器辅助结构
type Timing struct {
	operation string
	startTime time.Time
	details   map[string]interface{}
}

// StartTiming 开始计时
func StartTiming(operation string) *Timing {
	return &Timing{
		operation: operation,
		startTime: time.Now(),
		details:   make(map[string]interface{}),
	}
}

// AddDetail 添加详细信息
func (t *Timing) AddDetail(key string, value interface{}) *Timing {
	t.details[key] = value
	return t
}

// End 结束计时并记录性能日志
func (t *Timing) End() time.Duration {
	duration := time.Since(t.startTime)
	PerformanceWithDetails(t.operation, duration, t.details)
	return duration
}

// EndSimple 简单结束计时
func (t *Timing) EndSimple() time.Duration {
	duration := time.Since(t.startTime)
	Performance(t.operation, duration)
	return duration
}

// Progress 进度日志
func Progress(current, total int, operation string) {
	percentage := float64(current) / float64(total) * 100
	var emoji string
	
	switch {
	case percentage == 100:
		emoji = "✅"
	case percentage >= 75:
		emoji = "🔵"
	case percentage >= 50:
		emoji = "🟡"
	case percentage >= 25:
		emoji = "🟠"
	default:
		emoji = "🔴"
	}
	
	defaultLogger.logger.Printf("%s [PROGRESS] %s: %d/%d (%.1f%%)", 
		emoji, operation, current, total, percentage)
}

// Milestone 里程碑日志
func Milestone(message string) {
	defaultLogger.logger.Printf("🎯 [MILESTONE] %s", message)
}

// Health 健康检查日志
func Health(service string, status bool, details string) {
	emoji := "❌"
	statusStr := "UNHEALTHY"
	if status {
		emoji = "✅"
		statusStr = "HEALTHY"
	}
	
	detailStr := ""
	if details != "" {
		detailStr = fmt.Sprintf(" | %s", details)
	}
	
	defaultLogger.logger.Printf("%s [HEALTH] %s: %s%s", emoji, service, statusStr, detailStr)
}

// Audit 审计日志
func Audit(action string, user string, resource string, result string) {
	defaultLogger.logger.Printf("📋 [AUDIT] User: %s | Action: %s | Resource: %s | Result: %s", 
		user, action, resource, result)
}