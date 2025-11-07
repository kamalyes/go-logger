/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\interfaces.go
 * @Description: 日志接口定义
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import "io"

// ILogger 日志记录器接口
type ILogger interface {
	// 基本日志方法
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	
	// 配置方法
	SetLevel(level LogLevel)
	GetLevel() LogLevel
	SetShowCaller(show bool)
	IsShowCaller() bool
	IsLevelEnabled(level LogLevel) bool
	
	// 结构化日志
	WithField(key string, value interface{}) ILogger
	WithFields(fields map[string]interface{}) ILogger
	WithError(err error) ILogger
	
	// 实用方法
	Clone() ILogger
}

// IAdapter 日志适配器接口
type IAdapter interface {
	ILogger
	
	// 生命周期管理
	Initialize() error
	Close() error
	Flush() error
	
	// 适配器特定功能
	GetAdapterName() string
	GetAdapterVersion() string
	IsHealthy() bool
}

// IManager 日志管理器接口  
type IManager interface {
	// 适配器管理
	AddAdapter(name string, adapter IAdapter) error
	GetAdapter(name string) (IAdapter, bool)
	RemoveAdapter(name string) error
	ListAdapters() []string
	
	// 生命周期管理
	CloseAll() error
	FlushAll() error
	
	// 全局设置
	SetLevelAll(level LogLevel)
	Broadcast(level LogLevel, format string, args ...interface{})
	
	// 健康检查
	HealthCheck() map[string]bool
}

// IFileWriter 文件写入器接口
type IFileWriter interface {
	io.WriteCloser
	
	// 文件管理
	Rotate() error
	GetCurrentFile() string
	GetFileSize() int64
	
	// 配置
	SetMaxSize(size int64)
	SetMaxBackups(backups int)
	SetMaxAge(days int)
	SetCompress(compress bool)
}

// IFormatter 日志格式化器接口
type IFormatter interface {
	Format(entry *LogEntry) ([]byte, error)
	GetName() string
}

// LogEntry 日志条目
type LogEntry struct {
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Timestamp int64                  `json:"timestamp"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Caller    *CallerInfo            `json:"caller,omitempty"`
}

// CallerInfo 调用者信息
type CallerInfo struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

// IHook 日志钩子接口
type IHook interface {
	Fire(entry *LogEntry) error
	Levels() []LogLevel
}

// IWriter 日志写入器接口
type IWriter interface {
	io.Writer
	
	// 写入控制
	WriteLevel(level LogLevel, data []byte) (n int, err error)
	Flush() error
	Close() error
	
	// 状态查询
	IsHealthy() bool
	GetStats() interface{}
}

// IMiddleware 日志中间件接口
type IMiddleware interface {
	Process(entry *LogEntry, next func(*LogEntry) error) error
	GetName() string
	GetPriority() int
}