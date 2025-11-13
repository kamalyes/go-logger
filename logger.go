/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 09:42:37
 * @FilePath: \go-logger\logger.go
 * @Description: 统一的日志工具包，支持 emoji 和结构化日志
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

// 性能优化: 缓冲池
var stringBuilderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// 性能优化: 预计算的级别格式
var (
	levelFormatsCache      = make(map[LogLevel]string)
	colorLevelFormatsCache = make(map[LogLevel]string)
	initCacheOnce          sync.Once
)

func initLevelFormatsCache() {
	levels := []LogLevel{DEBUG, INFO, WARN, ERROR, FATAL}

	for _, level := range levels {
		// 普通格式
		levelFormatsCache[level] = fmt.Sprintf("%s [%s]", level.Emoji(), level.String())

		// 彩色格式
		colorLevelFormatsCache[level] = level.Color() + levelFormatsCache[level] + "\033[0m"
	}
}

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

// formatMessage 格式化消息 - 优化版本
func (l *Logger) formatMessage(level LogLevel, format string, args ...interface{}) string {
	// 初始化缓存
	initCacheOnce.Do(initLevelFormatsCache)

	// 提前检查级别，避免不必要的计算
	if level < l.level {
		return ""
	}

	// 使用 strings.Builder 减少内存分配
	sb := stringBuilderPool.Get().(*strings.Builder)
	defer func() {
		sb.Reset()
		stringBuilderPool.Put(sb)
	}()

	// 预估容量
	estimatedSize := len(format) + 100
	if l.showCaller {
		estimatedSize += 50
	}
	sb.Grow(estimatedSize)

	// 使用预计算的级别格式
	if l.config.Colorful {
		sb.WriteString(colorLevelFormatsCache[level])
	} else {
		sb.WriteString(levelFormatsCache[level])
	}

	// 添加调用者信息（如果需要）
	if l.showCaller {
		if pc, file, line, ok := runtime.Caller(3); ok {
			funcName := runtime.FuncForPC(pc).Name()
			if idx := strings.LastIndex(funcName, "."); idx != -1 {
				funcName = funcName[idx+1:]
			}
			if idx := strings.LastIndex(file, "/"); idx != -1 {
				file = file[idx+1:]
			}
			sb.WriteString(fmt.Sprintf(" [%s:%d:%s]", file, line, funcName))
		}
	}

	// 添加消息
	sb.WriteByte(' ')
	if len(args) == 0 {
		sb.WriteString(format)
	} else {
		// 只在需要时格式化
		sb.WriteString(fmt.Sprintf(format, args...))
	}

	return sb.String()
}

// log 记录日志 - 优化版本
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	// 提前检查级别，避免不必要的计算
	if level < l.level {
		return
	}

	message := l.formatMessage(level, format, args...)
	if message != "" { // 只有非空消息才输出
		l.logger.Print(message)
	}

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

// Printf风格方法（与上面相同，但命名更明确）
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// WithField 添加字段信息（结构化日志）
func (l *Logger) WithField(key string, value interface{}) ILogger {
	prefix := fmt.Sprintf("%s%s=%v ", l.config.Prefix, key, value)
	config := l.config.Clone()
	config.Prefix = prefix

	return NewLogger(config)
}

// WithFields 添加多个字段信息（结构化日志）
func (l *Logger) WithFields(fields map[string]interface{}) ILogger {
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
func (l *Logger) WithError(err error) ILogger {
	return l.WithField("error", err.Error())
}

// Clone 克隆当前Logger
func (l *Logger) Clone() ILogger {
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

// WithField 全局添加字段
func WithField(key string, value interface{}) ILogger {
	return defaultLogger.WithField(key, value)
}

// WithFields 全局添加多个字段
func WithFields(fields map[string]interface{}) ILogger {
	return defaultLogger.WithFields(fields)
}

// WithError 全局添加错误信息
func WithError(err error) ILogger {
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

// 为 Logger 添加新接口方法的实现

// 纯文本日志方法
func (l *Logger) DebugMsg(msg string) {
	l.Debug("%s", msg)
}

func (l *Logger) InfoMsg(msg string) {
	l.Info("%s", msg)
}

func (l *Logger) WarnMsg(msg string) {
	l.Warn("%s", msg)
}

func (l *Logger) ErrorMsg(msg string) {
	l.Error("%s", msg)
}

func (l *Logger) FatalMsg(msg string) {
	l.Fatal("%s", msg)
}

// 带上下文的日志方法
func (l *Logger) DebugContext(ctx context.Context, format string, args ...interface{}) {
	// 目前忽略context，委托给基础方法
	l.Debug(format, args...)
}

func (l *Logger) InfoContext(ctx context.Context, format string, args ...interface{}) {
	l.Info(format, args...)
}

func (l *Logger) WarnContext(ctx context.Context, format string, args ...interface{}) {
	l.Warn(format, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, format string, args ...interface{}) {
	l.Error(format, args...)
}

func (l *Logger) FatalContext(ctx context.Context, format string, args ...interface{}) {
	l.Fatal(format, args...)
}

// 结构化日志方法（键值对）
func (l *Logger) DebugKV(msg string, keysAndValues ...interface{}) {
	fields := l.parseKeysAndValues(keysAndValues...)
	if len(fields) > 0 {
		logger := l.WithFields(fields).(*Logger) // 类型转换
		logger.Debug("%s", msg)
	} else {
		l.Debug("%s", msg)
	}
}

func (l *Logger) InfoKV(msg string, keysAndValues ...interface{}) {
	fields := l.parseKeysAndValues(keysAndValues...)
	if len(fields) > 0 {
		logger := l.WithFields(fields).(*Logger) // 类型转换
		logger.Info("%s", msg)
	} else {
		l.Info("%s", msg)
	}
}

func (l *Logger) WarnKV(msg string, keysAndValues ...interface{}) {
	fields := l.parseKeysAndValues(keysAndValues...)
	if len(fields) > 0 {
		logger := l.WithFields(fields).(*Logger) // 类型转换
		logger.Warn("%s", msg)
	} else {
		l.Warn("%s", msg)
	}
}

func (l *Logger) ErrorKV(msg string, keysAndValues ...interface{}) {
	fields := l.parseKeysAndValues(keysAndValues...)
	if len(fields) > 0 {
		logger := l.WithFields(fields).(*Logger) // 类型转换
		logger.Error("%s", msg)
	} else {
		l.Error("%s", msg)
	}
}

func (l *Logger) FatalKV(msg string, keysAndValues ...interface{}) {
	fields := l.parseKeysAndValues(keysAndValues...)
	if len(fields) > 0 {
		logger := l.WithFields(fields).(*Logger) // 类型转换
		logger.Fatal("%s", msg)
	} else {
		l.Fatal("%s", msg)
	}
}

// 原始日志条目方法
func (l *Logger) Log(level LogLevel, msg string) {
	switch level {
	case DEBUG:
		l.Debug("%s", msg)
	case INFO:
		l.Info("%s", msg)
	case WARN:
		l.Warn("%s", msg)
	case ERROR:
		l.Error("%s", msg)
	case FATAL:
		l.Fatal("%s", msg)
	}
}

func (l *Logger) LogContext(ctx context.Context, level LogLevel, msg string) {
	// 默认实现忽略context
	l.Log(level, msg)
}

func (l *Logger) LogKV(level LogLevel, msg string, keysAndValues ...interface{}) {
	fields := l.parseKeysAndValues(keysAndValues...)
	logger := l
	if len(fields) > 0 {
		logger = logger.WithFields(fields).(*Logger) // 这里需要类型转换，因为我们知道返回的是 *Logger
	}

	switch level {
	case DEBUG:
		logger.Debug("%s", msg)
	case INFO:
		logger.Info("%s", msg)
	case WARN:
		logger.Warn("%s", msg)
	case ERROR:
		logger.Error("%s", msg)
	case FATAL:
		logger.Fatal("%s", msg)
	}
}

func (l *Logger) LogWithFields(level LogLevel, msg string, fields map[string]interface{}) {
	logger := l
	if len(fields) > 0 {
		logger = logger.WithFields(fields).(*Logger) // 类型转换
	}

	switch level {
	case DEBUG:
		logger.Debug("%s", msg)
	case INFO:
		logger.Info("%s", msg)
	case WARN:
		logger.Warn("%s", msg)
	case ERROR:
		logger.Error("%s", msg)
	case FATAL:
		logger.Fatal("%s", msg)
	}
}

// WithContext 带上下文的logger（当前实现返回自身）
func (l *Logger) WithContext(ctx context.Context) ILogger {
	// 创建一个新的logger实例并设置context
	newLogger := l.Clone()
	if loggerPtr, ok := newLogger.(*Logger); ok {
		loggerPtr.context = ctx
		return loggerPtr
	}
	return newLogger
}

// 兼容标准log包的方法
func (l *Logger) Print(args ...interface{}) {
	l.Info("%s", fmt.Sprint(args...))
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.Info(format, args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.Info("%s", fmt.Sprintln(args...))
}

// parseKeysAndValues 解析键值对参数 - 优化版本
func (l *Logger) parseKeysAndValues(keysAndValues ...interface{}) map[string]interface{} {
	if len(keysAndValues) == 0 {
		return nil
	}

	// 预分配合适大小的map
	fields := make(map[string]interface{}, len(keysAndValues)/2+1)

	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			// 优化字符串转换
			key := toString(keysAndValues[i])
			fields[key] = keysAndValues[i+1]
		} else {
			// 奇数个参数，最后一个作为无值key
			key := toString(keysAndValues[i])
			fields[key] = ""
		}
	}
	return fields
}

// toString 高效的字符串转换
func toString(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case fmt.Stringer:
		return s.String()
	default:
		return fmt.Sprint(v)
	}
}
