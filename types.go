/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 02:07:30
 * @FilePath: \go-logger\types.go
 * @Description: 核心类型定义
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Logger 主要的日志记录器结构体
type Logger struct {
	// 基本配置
	level      LogLevel
	showCaller bool
	config     *LogConfig
	
	// 内部组件
	logger     *log.Logger
	formatter  IFormatter
	writers    []IWriter
	hooks      []IHook
	middleware []IMiddleware
	
	// 同步和上下文
	context context.Context
	cancel  context.CancelFunc
	
	// 统计信息
	stats *LoggerStats
}

// LoggerStats 日志统计信息
type LoggerStats struct {
	StartTime    time.Time            `json:"start_time"`
	TotalLogs    int64                `json:"total_logs"`
	LevelCounts  map[LogLevel]int64   `json:"level_counts"`
	ErrorCount   int64                `json:"error_count"`
	LastLogTime  time.Time            `json:"last_log_time"`
	Uptime       time.Duration        `json:"uptime"`
	BytesWritten int64                `json:"bytes_written"`
	mutex        sync.RWMutex
}

// NewLoggerStats 创建新的统计信息
func NewLoggerStats() *LoggerStats {
	return &LoggerStats{
		StartTime:   time.Now(),
		LevelCounts: make(map[LogLevel]int64),
	}
}

// IncrementLevel 增加级别计数
func (s *LoggerStats) IncrementLevel(level LogLevel) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.TotalLogs++
	s.LevelCounts[level]++
	s.LastLogTime = time.Now()
	s.Uptime = time.Since(s.StartTime)
	
	if level >= ERROR {
		s.ErrorCount++
	}
}

// AddBytes 增加字节数
func (s *LoggerStats) AddBytes(bytes int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.BytesWritten += bytes
}

// LoggerStatsSnapshot 统计信息快照
type LoggerStatsSnapshot struct {
	StartTime    time.Time            `json:"start_time"`
	TotalLogs    int64                `json:"total_logs"`
	LevelCounts  map[LogLevel]int64   `json:"level_counts"`
	ErrorCount   int64                `json:"error_count"`
	LastLogTime  time.Time            `json:"last_log_time"`
	Uptime       time.Duration        `json:"uptime"`
	BytesWritten int64                `json:"bytes_written"`
}

// GetStats 获取统计信息快照
func (s *LoggerStats) GetStats() LoggerStatsSnapshot {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	// 创建一个新的快照避免复制 mutex
	stats := LoggerStatsSnapshot{
		StartTime:    s.StartTime,
		TotalLogs:    s.TotalLogs,
		ErrorCount:   s.ErrorCount,
		LastLogTime:  s.LastLogTime,
		Uptime:       s.Uptime,
		BytesWritten: s.BytesWritten,
		LevelCounts:  make(map[LogLevel]int64),
	}
	
	for k, v := range s.LevelCounts {
		stats.LevelCounts[k] = v
	}
	
	return stats
}

// LoggerOptions 创建Logger时的选项
type LoggerOptions struct {
	Config     *LogConfig    `json:"config"`
	Formatter  IFormatter    `json:"-"`
	Writers    []IWriter     `json:"-"`
	Hooks      []IHook       `json:"-"`
	Middleware []IMiddleware `json:"-"`
	Context    context.Context `json:"-"`
}

// DefaultLoggerOptions 默认Logger选项
func DefaultLoggerOptions() *LoggerOptions {
	return &LoggerOptions{
		Config:     DefaultConfig(),
		Writers:    []IWriter{},
		Hooks:      []IHook{},
		Middleware: []IMiddleware{},
		Context:    context.Background(),
	}
}

// FieldMap 字段映射类型
type FieldMap map[string]interface{}

// LogContext 日志上下文
type LogContext struct {
	RequestID   string    `json:"request_id,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	SessionID   string    `json:"session_id,omitempty"`
	Operation   string    `json:"operation,omitempty"`
	Component   string    `json:"component,omitempty"`
	Version     string    `json:"version,omitempty"`
	Environment string    `json:"environment,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Fields      FieldMap  `json:"fields,omitempty"`
}

// NewLogContext 创建新的日志上下文
func NewLogContext() *LogContext {
	return &LogContext{
		Timestamp: time.Now(),
		Fields:    make(FieldMap),
	}
}

// WithField 添加字段
func (lc *LogContext) WithField(key string, value interface{}) *LogContext {
	if lc.Fields == nil {
		lc.Fields = make(FieldMap)
	}
	lc.Fields[key] = value
	return lc
}

// WithFields 添加多个字段
func (lc *LogContext) WithFields(fields FieldMap) *LogContext {
	if lc.Fields == nil {
		lc.Fields = make(FieldMap)
	}
	for k, v := range fields {
		lc.Fields[k] = v
	}
	return lc
}

// Clone 克隆上下文
func (lc *LogContext) Clone() *LogContext {
	clone := &LogContext{
		RequestID:   lc.RequestID,
		UserID:      lc.UserID,
		SessionID:   lc.SessionID,
		Operation:   lc.Operation,
		Component:   lc.Component,
		Version:     lc.Version,
		Environment: lc.Environment,
		Timestamp:   lc.Timestamp,
		Fields:      make(FieldMap),
	}
	
	for k, v := range lc.Fields {
		clone.Fields[k] = v
	}
	
	return clone
}

// AdapterRegistry 适配器注册表
type AdapterRegistry struct {
	adapters map[string]AdapterFactory
	mutex    sync.RWMutex
}

// AdapterFactory 适配器工厂函数
type AdapterFactory func(config *AdapterConfig) (IAdapter, error)

// NewAdapterRegistry 创建适配器注册表
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: make(map[string]AdapterFactory),
	}
}

// Register 注册适配器工厂
func (r *AdapterRegistry) Register(name string, factory AdapterFactory) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.adapters[name] = factory
}

// Create 创建适配器
func (r *AdapterRegistry) Create(name string, config *AdapterConfig) (IAdapter, error) {
	r.mutex.RLock()
	factory, exists := r.adapters[name]
	r.mutex.RUnlock()
	
	if !exists {
				return nil, fmt.Errorf("adapter factory '%s' not found", name)
	}
	
	return factory(config)
}

// List 列出所有注册的适配器
func (r *AdapterRegistry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	return names
}

// BufferPool 缓冲区池
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool 创建缓冲区池
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, 1024)
			},
		},
	}
}

// Get 获取缓冲区
func (bp *BufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}

// Put 归还缓冲区
func (bp *BufferPool) Put(buf []byte) {
	if cap(buf) > 64*1024 {
		return // 防止缓冲区过大
	}
	bp.pool.Put(buf[:0])
}

// EventType 事件类型
type EventType int

const (
	EventLogCreated EventType = iota
	EventLogProcessed
	EventLogWritten
	EventLogError
	EventLoggerStarted
	EventLoggerStopped
	EventLoggerConfigChanged
)

// LogEvent 日志事件
type LogEvent struct {
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Level     LogLevel    `json:"level"`
	Message   string      `json:"message"`
	Fields    FieldMap    `json:"fields,omitempty"`
	Error     error       `json:"error,omitempty"`
	Context   *LogContext `json:"context,omitempty"`
}

// NewLogEvent 创建日志事件
func NewLogEvent(eventType EventType, level LogLevel, message string) *LogEvent {
	event := &LogEvent{
		Type:      eventType,
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    make(FieldMap),
	}
	return event
}

// NewLoggerWithOptions 使用选项创建Logger
func NewLoggerWithOptions(options *LoggerOptions) *Logger {
	if options == nil {
		options = DefaultLoggerOptions()
	}
	
	// 创建上下文
	ctx := options.Context
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithCancel(ctx)
	
	logger := &Logger{
		level:      options.Config.Level,
		showCaller: options.Config.ShowCaller,
		config:     options.Config,
		formatter:  options.Formatter,
		writers:    options.Writers,
		hooks:      options.Hooks,
		middleware: options.Middleware,
		context:    ctx,
		cancel:     cancel,
		stats:      NewLoggerStats(),
		logger:     log.New(options.Config.Output, options.Config.Prefix, log.LstdFlags),
	}
	
	// 如果没有指定格式化器，使用默认的
	if logger.formatter == nil {
		logger.formatter = NewTextFormatter()
	}
	
	// 如果没有写入器，添加默认控制台写入器
	if len(logger.writers) == 0 {
		logger.writers = append(logger.writers, NewConsoleWriter(options.Config.Output))
	}
	
	return logger
}