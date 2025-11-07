/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\level.go
 * @Description: 日志级别定义和配置
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"fmt"
	"strings"
)

// LogLevel 日志级别
type LogLevel int

const (
	// DEBUG 调试级别 - 最详细的信息
	DEBUG LogLevel = iota
	// INFO 信息级别 - 一般信息
	INFO
	// WARN 警告级别 - 警告信息
	WARN
	// ERROR 错误级别 - 错误信息
	ERROR
	// FATAL 致命级别 - 致命错误，程序将退出
	FATAL
)

// levelInfo 日志级别信息
type levelInfo struct {
	emoji string
	name  string
	color string
}

// 日志级别对应的emoji、名称和颜色
var levelConfig = map[LogLevel]levelInfo{
	DEBUG: {"🐛", "DEBUG", "\033[36m"},   // 青色
	INFO:  {"ℹ️", "INFO", "\033[32m"},    // 绿色
	WARN:  {"⚠️", "WARN", "\033[33m"},    // 黄色
	ERROR: {"❌", "ERROR", "\033[31m"},   // 红色
	FATAL: {"💀", "FATAL", "\033[35m"},   // 紫色
}

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	if info, ok := levelConfig[l]; ok {
		return info.name
	}
	return fmt.Sprintf("UNKNOWN(%d)", int(l))
}

// Emoji 返回日志级别的emoji
func (l LogLevel) Emoji() string {
	if info, ok := levelConfig[l]; ok {
		return info.emoji
	}
	return "❓"
}

// Color 返回日志级别的颜色代码
func (l LogLevel) Color() string {
	if info, ok := levelConfig[l]; ok {
		return info.color
	}
	return "\033[0m" // 重置颜色
}

// ParseLevel 从字符串解析日志级别
func ParseLevel(level string) (LogLevel, error) {
	level = strings.ToUpper(strings.TrimSpace(level))
	switch level {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN", "WARNING":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		return INFO, fmt.Errorf("invalid log level: %s", level)
	}
}

// IsEnabled 检查给定级别是否在当前级别下启用
func (l LogLevel) IsEnabled(targetLevel LogLevel) bool {
	return targetLevel >= l
}

// GetAllLevels 获取所有可用的日志级别
func GetAllLevels() []LogLevel {
	return []LogLevel{DEBUG, INFO, WARN, ERROR, FATAL}
}

// GetLevelNames 获取所有可用的日志级别名称
func GetLevelNames() []string {
	levels := GetAllLevels()
	names := make([]string, len(levels))
	for i, level := range levels {
		names[i] = level.String()
	}
	return names
}