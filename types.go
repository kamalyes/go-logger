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
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// FormatType 定义
// ============================================================================

// FormatType 输出格式类型
type FormatType string

const (
	FormatText FormatType = "text"
	FormatJSON FormatType = "json"
	FormatXML  FormatType = "xml"
	FormatCSV  FormatType = "csv"
)

// Logger 主要的日志记录器结构体
type Logger struct {
	// 基本配置（level 和 showCaller 使用原子操作，允许运行时并发修改）
	level          atomic.Int32
	showCaller     atomic.Bool
	colorful       bool
	prefix         string
	timeFormat     string
	format         FormatType
	callerDepth    int
	showStacktrace bool

	// 字段名配置：热路径 4 键直接存预构建 "key": JSON 片段（设置时构建一次，单次 append 拷贝），
	// 默认值共享包级片段（defaultXxxKey，见 logger.go）；stacktrace 为冷路径保留原字符串
	timestampKey  []byte
	levelKey      []byte
	messageKey    []byte
	callerKey     []byte
	stacktraceKey string

	// 异步写入配置
	asyncWrite   bool
	bufferSize   int
	batchSize    int
	batchTimeout time.Duration

	// 输出和同步
	output      io.Writer
	outLockFree bool       // 输出目标自带并发安全（内置 writer / io.Discard）时热路径免互斥
	mu          sync.Mutex // 保护非并发安全输出的写入

	// 内部组件
	logger     *log.Logger
	formatter  IFormatter
	writers    []IWriter
	hooks      []IHook
	middleware []IMiddleware

	// 上下文支持
	context          context.Context
	cancel           context.CancelFunc
	contextKeys      []compiledContextKey
	contextExtractor ContextExtractor

	// 统计信息
	stats *LoggerStats

	// Console 功能
	consoleGroup     *ConsoleGroup
	consoleGroupOnce sync.Once
}

// LoggerStats 日志统计信息
type LoggerStats struct {
	StartTime    time.Time          `json:"start_time"`
	TotalLogs    int64              `json:"total_logs"`
	LevelCounts  map[LogLevel]int64 `json:"level_counts"`
	ErrorCount   int64              `json:"error_count"`
	LastLogTime  time.Time          `json:"last_log_time"`
	Uptime       time.Duration      `json:"uptime"`
	BytesWritten int64              `json:"bytes_written"`
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

// GetStats 获取统计信息快照
func (s *LoggerStats) GetStats() *LoggerStats {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 创建一个新的快照避免复制 mutex
	clone := &LoggerStats{
		LevelCounts: make(map[LogLevel]int64),
	}

	// 手动复制值字段（跳过 mutex 与 LevelCounts map）
	clone.StartTime = s.StartTime
	clone.TotalLogs = s.TotalLogs
	clone.ErrorCount = s.ErrorCount
	clone.LastLogTime = s.LastLogTime
	clone.Uptime = s.Uptime
	clone.BytesWritten = s.BytesWritten

	// 手动复制 map
	for k, v := range s.LevelCounts {
		clone.LevelCounts[k] = v
	}

	return clone
}

// ============================================================================
// Logger 构造函数
// ============================================================================

// concurrentSafeWriter 标记 Write 已自带并发安全（内部持锁或无状态丢弃），
// Logger 输出路径遇到此标记时跳过外层互斥，消除并行写日志的锁竞争
type concurrentSafeWriter interface {
	concurrentSafeMarker()
}

// isConcurrentSafeOutput 判断输出目标是否无需 Logger 外层互斥：
// 内置 writer（自带锁）与 io.Discard（无状态）为安全；其余（os.Stdout、用户自定义 writer）保守视为不安全，维持加锁语义
func isConcurrentSafeOutput(w io.Writer) bool {
	if w == io.Discard {
		return true
	}
	_, ok := w.(concurrentSafeWriter)
	return ok
}

// NewLogger 创建新的日志记录器（默认配置）
func NewLogger() *Logger {
	l := &Logger{
		colorful:       true,
		prefix:         "",
		timeFormat:     time.RFC3339Nano,
		format:         FormatJSON,
		callerDepth:    2,
		showStacktrace: false,
		timestampKey:   defaultTimestampKey,
		levelKey:       defaultLevelKey,
		messageKey:     defaultMessageKey,
		callerKey:      defaultCallerKey,
		stacktraceKey:  "stacktrace",
		asyncWrite:     false,
		bufferSize:     0,
		batchSize:      100,
		batchTimeout:   100 * time.Millisecond,
		output:         os.Stdout,
		outLockFree:    false, // os.Stdout 并发安全未受保证，保守加锁
		logger:         log.New(os.Stdout, "", log.LstdFlags),
		contextKeys:    append([]compiledContextKey(nil), defaultCompiledContextKeys...),
		stats:          NewLoggerStats(),
		mu:             sync.Mutex{},
	}
	l.level.Store(int32(DEBUG))
	l.showCaller.Store(true)
	return l
}

// ============================================================================
// Builder 模式方法（链式调用）
// ============================================================================

// WithLevel 设置日志级别
func (l *Logger) WithLevel(level LogLevel) *Logger {
	l.level.Store(int32(level))
	return l
}

// WithShowCaller 设置是否显示调用者信息
func (l *Logger) WithShowCaller(show bool) *Logger {
	l.showCaller.Store(show)
	return l
}

// WithPrefix 设置日志前缀
func (l *Logger) WithPrefix(prefix string) *Logger {
	if prefix != "" && !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}
	l.prefix = prefix
	l.logger = log.New(l.output, prefix, log.LstdFlags)
	return l
}

// WithColorful 设置是否使用彩色输出
func (l *Logger) WithColorful(colorful bool) *Logger {
	l.colorful = colorful
	return l
}

// WithOutput 设置输出目标
func (l *Logger) WithOutput(output io.Writer) *Logger {
	l.output = output
	l.outLockFree = isConcurrentSafeOutput(output)
	l.logger = log.New(output, l.prefix, log.LstdFlags)
	return l
}

// WithTimeFormat 设置时间格式
func (l *Logger) WithTimeFormat(format string) *Logger {
	l.timeFormat = format
	return l
}

// WithFormat 设置输出格式
func (l *Logger) WithFormat(format FormatType) *Logger {
	l.format = format
	return l
}

// WithCallerDepth 设置调用者深度
func (l *Logger) WithCallerDepth(depth int) *Logger {
	l.callerDepth = depth
	return l
}

// WithShowStacktrace 设置是否显示堆栈跟踪
func (l *Logger) WithShowStacktrace(show bool) *Logger {
	l.showStacktrace = show
	return l
}

// WithTimestampKey 设置时间戳字段名（条目首字段，片段不带前导逗号）
func (l *Logger) WithTimestampKey(key string) *Logger {
	l.timestampKey = keyJSONFragment(key, false)
	return l
}

// WithLevelKey 设置日志级别字段名
func (l *Logger) WithLevelKey(key string) *Logger {
	l.levelKey = keyJSONFragment(key, true)
	return l
}

// WithMessageKey 设置消息字段名
func (l *Logger) WithMessageKey(key string) *Logger {
	l.messageKey = keyJSONFragment(key, true)
	return l
}

// WithCallerKey 设置调用者字段名
func (l *Logger) WithCallerKey(key string) *Logger {
	l.callerKey = keyJSONFragment(key, true)
	return l
}

// WithStacktraceKey 设置堆栈跟踪字段名
func (l *Logger) WithStacktraceKey(key string) *Logger {
	l.stacktraceKey = key
	return l
}

// WithAsyncWrite 设置是否异步写入
func (l *Logger) WithAsyncWrite(async bool) *Logger {
	l.asyncWrite = async
	return l
}

// WithBufferSize 设置缓冲区大小
func (l *Logger) WithBufferSize(size int) *Logger {
	l.bufferSize = size
	return l
}

// WithBatchSize 设置批量写入大小
func (l *Logger) WithBatchSize(size int) *Logger {
	l.batchSize = size
	return l
}

// WithBatchTimeout 设置批量写入超时时间
func (l *Logger) WithBatchTimeout(timeout time.Duration) *Logger {
	l.batchTimeout = timeout
	return l
}

// WithFormatter 设置格式化器
func (l *Logger) WithFormatter(formatter IFormatter) *Logger {
	l.formatter = formatter
	return l
}

// WithWriters 设置写入器列表
func (l *Logger) WithWriters(writers []IWriter) *Logger {
	l.writers = writers
	return l
}

// WithHooks 设置钩子列表
func (l *Logger) WithHooks(hooks []IHook) *Logger {
	l.hooks = hooks
	return l
}

// WithMiddleware 设置中间件列表
func (l *Logger) WithMiddleware(middleware []IMiddleware) *Logger {
	l.middleware = middleware
	return l
}

// WithContextExtractor 设置上下文提取器
func (l *Logger) WithContextExtractor(extractor ContextExtractor) *Logger {
	l.contextExtractor = extractor
	return l
}

// IsShowCaller 检查是否显示调用者信息
func (l *Logger) IsShowCaller() bool {
	return l.showCaller.Load()
}

// IsLevelEnabled 检查给定级别是否启用
func (l *Logger) IsLevelEnabled(level LogLevel) bool {
	return level >= LogLevel(l.level.Load())
}

func (l *Logger) Clone() ILogger {
	newLogger := &Logger{
		stats: NewLoggerStats(),
	}

	// 手动复制字段（跳过 mutex、sync.Once、stats 与原子状态）
	newLogger.colorful = l.colorful
	newLogger.prefix = l.prefix
	newLogger.timeFormat = l.timeFormat
	newLogger.format = l.format
	newLogger.callerDepth = l.callerDepth
	newLogger.showStacktrace = l.showStacktrace
	newLogger.timestampKey = l.timestampKey
	newLogger.levelKey = l.levelKey
	newLogger.messageKey = l.messageKey
	newLogger.callerKey = l.callerKey
	newLogger.stacktraceKey = l.stacktraceKey
	newLogger.asyncWrite = l.asyncWrite
	newLogger.bufferSize = l.bufferSize
	newLogger.batchSize = l.batchSize
	newLogger.batchTimeout = l.batchTimeout
	newLogger.output = l.output
	newLogger.outLockFree = l.outLockFree
	newLogger.logger = l.logger
	newLogger.formatter = l.formatter
	newLogger.writers = l.writers
	newLogger.contextKeys = append([]compiledContextKey(nil), l.contextKeys...)

	// 原子字段必须手动拷贝（DeepCopy 无法处理 atomic.Int32/atomic.Bool）
	newLogger.level.Store(l.level.Load())
	newLogger.showCaller.Store(l.showCaller.Load())

	// 确保使用新的统计信息
	newLogger.stats = NewLoggerStats()
	newLogger.contextExtractor = l.contextExtractor

	return newLogger
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
