/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-09-18 11:15:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 09:29:33
 * @FilePath: \go-logger\empty.go
 * @Description: 空日志实现，用于禁用日志输出的场景
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"sync"
)

// EmptyLogger 是 ILogger 接口的空实现，所有方法都不执行任何操作
// 常用于测试环境或需要禁用日志输出的场景
type EmptyLogger struct {
	level      LogLevel
	showCaller bool
	stats      *LoggerStats
}

// NewEmptyLogger 创建一个新的 EmptyLogger 实例
func NewEmptyLogger() *EmptyLogger {
	return &EmptyLogger{
		level:      INFO,
		showCaller: false,
		stats:      NewLoggerStats(),
	}
}

// NewEmptyLoggerWithLevel 创建一个指定级别的空日志实例
func NewEmptyLoggerWithLevel(level LogLevel) *EmptyLogger {
	logger := NewEmptyLogger()
	logger.level = level
	return logger
}

// 基本日志方法 - 所有方法都是空实现
func (e *EmptyLogger) Debug(format string, args ...interface{}) {}
func (e *EmptyLogger) Info(format string, args ...interface{})  {}
func (e *EmptyLogger) Warn(format string, args ...interface{})  {}
func (e *EmptyLogger) Error(format string, args ...interface{}) {}
func (e *EmptyLogger) Fatal(format string, args ...interface{}) {}

// Printf风格方法（与上面相同，但命名更明确）
func (e *EmptyLogger) Debugf(format string, args ...interface{}) {}
func (e *EmptyLogger) Infof(format string, args ...interface{})  {}
func (e *EmptyLogger) Warnf(format string, args ...interface{})  {}
func (e *EmptyLogger) Errorf(format string, args ...interface{}) {}
func (e *EmptyLogger) Fatalf(format string, args ...interface{}) {}

// 纯文本日志方法
func (e *EmptyLogger) DebugMsg(msg string) {}
func (e *EmptyLogger) InfoMsg(msg string)  {}
func (e *EmptyLogger) WarnMsg(msg string)  {}
func (e *EmptyLogger) ErrorMsg(msg string) {}
func (e *EmptyLogger) FatalMsg(msg string) {}

// 多行日志方法
func (e *EmptyLogger) InfoLines(lines ...string)  {}
func (e *EmptyLogger) ErrorLines(lines ...string) {}
func (e *EmptyLogger) WarnLines(lines ...string)  {}
func (e *EmptyLogger) DebugLines(lines ...string) {}

// 带上下文的日志方法
func (e *EmptyLogger) DebugContext(ctx context.Context, format string, args ...interface{}) {}
func (e *EmptyLogger) InfoContext(ctx context.Context, format string, args ...interface{})  {}
func (e *EmptyLogger) WarnContext(ctx context.Context, format string, args ...interface{})  {}
func (e *EmptyLogger) ErrorContext(ctx context.Context, format string, args ...interface{}) {}
func (e *EmptyLogger) FatalContext(ctx context.Context, format string, args ...interface{}) {}

// 结构化日志方法（键值对）
func (e *EmptyLogger) DebugKV(msg string, keysAndValues ...interface{}) {}
func (e *EmptyLogger) InfoKV(msg string, keysAndValues ...interface{})  {}
func (e *EmptyLogger) WarnKV(msg string, keysAndValues ...interface{})  {}
func (e *EmptyLogger) ErrorKV(msg string, keysAndValues ...interface{}) {}
func (e *EmptyLogger) FatalKV(msg string, keysAndValues ...interface{}) {}

// 带上下文的结构化日志方法
func (e *EmptyLogger) DebugContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {}
func (e *EmptyLogger) InfoContextKV(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (e *EmptyLogger) WarnContextKV(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (e *EmptyLogger) ErrorContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {}
func (e *EmptyLogger) FatalContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {}

// 原始日志条目方法
func (e *EmptyLogger) Log(level LogLevel, msg string)                                          {}
func (e *EmptyLogger) LogContext(ctx context.Context, level LogLevel, msg string)              {}
func (e *EmptyLogger) LogKV(level LogLevel, msg string, keysAndValues ...interface{})          {}
func (e *EmptyLogger) LogWithFields(level LogLevel, msg string, fields map[string]interface{}) {}

// 配置方法
func (e *EmptyLogger) SetLevel(level LogLevel) {
	e.level = level
}

func (e *EmptyLogger) GetLevel() LogLevel {
	return e.level
}

func (e *EmptyLogger) SetShowCaller(show bool) {
	e.showCaller = show
}

func (e *EmptyLogger) IsShowCaller() bool {
	return e.showCaller
}

func (e *EmptyLogger) IsLevelEnabled(level LogLevel) bool {
	return level >= e.level
}

// 结构化日志构建器方法 - 返回自身以保持链式调用
func (e *EmptyLogger) WithField(key string, value interface{}) ILogger {
	return e
}

func (e *EmptyLogger) WithFields(fields map[string]interface{}) ILogger {
	return e
}

func (e *EmptyLogger) WithError(err error) ILogger {
	return e
}

func (e *EmptyLogger) WithContext(ctx context.Context) ILogger {
	return e
}

// 实用方法
func (e *EmptyLogger) Clone() ILogger {
	return &EmptyLogger{
		level:      e.level,
		showCaller: e.showCaller,
		stats:      NewLoggerStats(),
	}
}

// 兼容标准log包的方法
func (e *EmptyLogger) Print(args ...interface{})                 {}
func (e *EmptyLogger) Printf(format string, args ...interface{}) {}
func (e *EmptyLogger) Println(args ...interface{})               {}

// GetStats 返回统计信息
func (e *EmptyLogger) GetStats() *LoggerStats {
	return e.stats
}

// EmptyAdapter 是 IAdapter 接口的空实现
type EmptyAdapter struct {
	*EmptyLogger
	name    string
	version string
	healthy bool
}

// NewEmptyAdapter 创建一个新的空适配器
func NewEmptyAdapter(name string) *EmptyAdapter {
	return &EmptyAdapter{
		EmptyLogger: NewEmptyLogger(),
		name:        name,
		version:     "1.0.0",
		healthy:     true,
	}
}

// Initialize 初始化适配器（空实现）
func (e *EmptyAdapter) Initialize() error {
	return nil
}

// Close 关闭适配器（空实现）
func (e *EmptyAdapter) Close() error {
	return nil
}

// Flush 刷新缓冲区（空实现）
func (e *EmptyAdapter) Flush() error {
	return nil
}

// GetAdapterName 获取适配器名称
func (e *EmptyAdapter) GetAdapterName() string {
	return e.name
}

// GetAdapterVersion 获取适配器版本
func (e *EmptyAdapter) GetAdapterVersion() string {
	return e.version
}

// IsHealthy 检查适配器健康状态
func (e *EmptyAdapter) IsHealthy() bool {
	return e.healthy
}

// SetHealthy 设置适配器健康状态
func (e *EmptyAdapter) SetHealthy(healthy bool) {
	e.healthy = healthy
}

// EmptyWriter 是 IWriter 接口的空实现
type EmptyWriter struct {
	mutex sync.Mutex
	stats map[string]interface{}
}

// NewEmptyWriter 创建一个新的空写入器
func NewEmptyWriter() *EmptyWriter {
	return &EmptyWriter{
		stats: make(map[string]interface{}),
	}
}

// Write 写入数据（空实现）
func (e *EmptyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// WriteLevel 按级别写入数据（空实现）
func (e *EmptyWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	return len(data), nil
}

// Flush 刷新缓冲区（空实现）
func (e *EmptyWriter) Flush() error {
	return nil
}

// Close 关闭写入器（空实现）
func (e *EmptyWriter) Close() error {
	return nil
}

// IsHealthy 检查写入器健康状态
func (e *EmptyWriter) IsHealthy() bool {
	return true
}

// GetStats 获取统计信息
func (e *EmptyWriter) GetStats() interface{} {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	result := make(map[string]interface{})
	for k, v := range e.stats {
		result[k] = v
	}
	return result
}

// EmptyHook 是 IHook 接口的空实现
type EmptyHook struct {
	levels []LogLevel
}

// NewEmptyHook 创建一个新的空钩子
func NewEmptyHook(levels []LogLevel) *EmptyHook {
	if levels == nil {
		levels = []LogLevel{DEBUG, INFO, WARN, ERROR, FATAL}
	}
	return &EmptyHook{
		levels: levels,
	}
}

// Fire 执行钩子（空实现）
func (e *EmptyHook) Fire(entry *LogEntry) error {
	return nil
}

// Levels 获取支持的级别
func (e *EmptyHook) Levels() []LogLevel {
	return e.levels
}

// EmptyMiddleware 是 IMiddleware 接口的空实现
type EmptyMiddleware struct {
	name     string
	priority int
}

// NewEmptyMiddleware 创建一个新的空中间件
func NewEmptyMiddleware(name string, priority int) *EmptyMiddleware {
	return &EmptyMiddleware{
		name:     name,
		priority: priority,
	}
}

// Process 处理日志条目（空实现，直接调用下一个处理器）
func (e *EmptyMiddleware) Process(entry *LogEntry, next func(*LogEntry) error) error {
	if next != nil {
		return next(entry)
	}
	return nil
}

// GetName 获取中间件名称
func (e *EmptyMiddleware) GetName() string {
	return e.name
}

// GetPriority 获取中间件优先级
func (e *EmptyMiddleware) GetPriority() int {
	return e.priority
}

// EmptyFormatter 是 IFormatter 接口的空实现
type EmptyFormatter struct {
	name string
}

// NewEmptyFormatter 创建一个新的空格式化器
func NewEmptyFormatter() *EmptyFormatter {
	return &EmptyFormatter{
		name: "empty",
	}
}

// Format 格式化日志条目（返回空字节）
func (e *EmptyFormatter) Format(entry *LogEntry) ([]byte, error) {
	return []byte{}, nil
}

// GetName 获取格式化器名称
func (e *EmptyFormatter) GetName() string {
	return e.name
}

// 全局空日志实例
var (
	// NoLogger 全局的空日志实例，用于禁用所有日志输出
	NoLogger ILogger = NewEmptyLogger()

	// DiscardLogger 废弃所有日志的实例，与NoLogger功能相同
	DiscardLogger ILogger = NewEmptyLogger()

	// NullLogger 空日志实例，用于测试场景
	NullLogger ILogger = NewEmptyLogger()
)

// DisableLogging 创建一个完全禁用日志的配置
func DisableLogging() *LogConfig {
	return &LogConfig{
		Level:      FATAL + 1, // 设置一个比FATAL更高的级别，禁用所有日志
		Output:     NewEmptyWriter(),
		Colorful:   false,
		ShowCaller: false,
		TimeFormat: "",
		Prefix:     "",
	}
}

// IsEmptyLogger 检查给定的logger是否是空日志实现
func IsEmptyLogger(logger ILogger) bool {
	_, ok := logger.(*EmptyLogger)
	return ok
}

// WrapWithEmpty 将任何日志器包装为空实现（用于临时禁用日志）
func WrapWithEmpty(original ILogger) ILogger {
	if IsEmptyLogger(original) {
		return original
	}

	empty := NewEmptyLogger()
	if original != nil {
		empty.SetLevel(original.GetLevel())
		empty.SetShowCaller(original.IsShowCaller())
	}
	return empty
}
