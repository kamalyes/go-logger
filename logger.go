/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-27 09:55:17
 * @FilePath: \go-logger\logger.go
 * @Description: 统一的高性能日志实现 - 整合所有功能
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"google.golang.org/grpc/metadata"
)

// ============================================================================
// 性能优化：对象池和预计算常量
// ============================================================================

const (
	maxLogMessageSize    = 1024 // 单条日志消息的最大预分配大小
	estimatedContextSize = 100  // 预估的上下文信息大小（TraceID/RequestID 等）

	// Metadata keys - 用于从 gRPC metadata 获取
	MetadataKeyTraceID   = "x-trace-id"
	MetadataKeyRequestID = "x-request-id"
)

// 字节池 - 用于日志消息构建
var bytePool = sync.Pool{
	New: func() any {
		return make([]byte, 0, maxLogMessageSize)
	},
}

// 上下文信息池 - 用于构建上下文字符串
var contextPool = sync.Pool{
	New: func() any {
		return make([]byte, 0, estimatedContextSize)
	},
}

// fieldMap 池 - 用于 fieldLogger 的 map 复用
var fieldMapPool = sync.Pool{
	New: func() any {
		return make(map[string]any, 8) // 预分配常见大小
	},
}

// 预计算的常量字节切片
var (
	debugPrefix = []byte("🐛 [DEBUG] ")
	infoPrefix  = []byte("ℹ️ [INFO] ")
	warnPrefix  = []byte("⚠️ [WARN] ")
	errorPrefix = []byte("❌ [ERROR] ")
	fatalPrefix = []byte("💀 [FATAL] ")

	debugPrefixColor = []byte("\033[36m🐛 [DEBUG]\033[0m ")
	infoPrefixColor  = []byte("\033[32mℹ️ [INFO]\033[0m ")
	warnPrefixColor  = []byte("\033[33m⚠️ [WARN]\033[0m ")
	errorPrefixColor = []byte("\033[31m❌ [ERROR]\033[0m ")
	fatalPrefixColor = []byte("\033[35m💀 [FATAL]\033[0m ")

	newline = []byte("\n")

	// context 提取的常量前缀（优化字符串拼接）
	traceIDPrefix   = []byte("TraceID=")
	requestIDPrefix = []byte(" RequestID=")
	bracketOpen     = []byte("[")
	bracketClose    = []byte("] ")

	// 键值对日志的常量字符串
	kvSeparator  = []byte(": ")
	kvDelimiter  = []byte(", ")
	kvBraceOpen  = []byte(" {")
	kvBraceClose = []byte("}")
	kvMissing    = []byte("<missing>")
)

var (
	levelPrefixes = map[LogLevel][]byte{
		DEBUG: debugPrefix,
		INFO:  infoPrefix,
		WARN:  warnPrefix,
		ERROR: errorPrefix,
		FATAL: fatalPrefix,
	}

	levelPrefixesColor = map[LogLevel][]byte{
		DEBUG: debugPrefixColor,
		INFO:  infoPrefixColor,
		WARN:  warnPrefixColor,
		ERROR: errorPrefixColor,
		FATAL: fatalPrefixColor,
	}
)

// ============================================================================
// 上下文提取器
// ============================================================================

// ContextExtractor 上下文信息提取器函数类型
type ContextExtractor func(ctx context.Context) string

// defaultContextExtractor 默认的上下文信息提取器
// 从 context.Context 中提取 TraceID 和 RequestID
func defaultContextExtractor(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	var traceID, requestID string

	// 1. 尝试从 context.Value 获取
	if tid, ok := ctx.Value(KeyTraceID).(string); ok && tid != "" {
		traceID = tid
	}
	if rid, ok := ctx.Value(KeyRequestID).(string); ok && rid != "" {
		requestID = rid
	}

	// 2. 如果还没找到，尝试从 gRPC metadata 获取
	if traceID == "" || requestID == "" {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if traceID == "" {
				if values := md.Get(MetadataKeyTraceID); len(values) > 0 {
					traceID = values[0]
				}
			}
			if requestID == "" {
				if values := md.Get(MetadataKeyRequestID); len(values) > 0 {
					requestID = values[0]
				}
			}
		}
	}

	// 3. 构建前缀（使用专用的上下文池和预分配的常量）
	if traceID != "" || requestID != "" {
		buf := contextPool.Get().([]byte)
		buf = buf[:0]
		defer contextPool.Put(buf)

		buf = append(buf, bracketOpen...)
		if traceID != "" {
			buf = append(buf, traceIDPrefix...)
			buf = append(buf, convert.S2B(traceID)...)
		}
		if requestID != "" {
			buf = append(buf, requestIDPrefix...)
			buf = append(buf, convert.S2B(requestID)...)
		}
		buf = append(buf, bracketClose...)
		return string(buf)
	}

	return ""
}

// ============================================================================
// Logger 结构体和初始化
// ============================================================================

// defaultLogger 默认日志记录器
var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger()
}

// New 创建新的日志记录器（简化版本）
func New() *Logger {
	return NewLogger()
}

// ultraLog 极致优化的日志方法（使用字节池和零拷贝）
func (l *Logger) ultraLog(level LogLevel, msg string) {
	if level < l.level {
		return
	}

	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	// 添加时间戳
	buf = convert.FastFormatTime(buf, time.Now())

	// 添加前缀（如果有）
	if l.prefix != "" {
		buf = append(buf, convert.S2B(l.prefix)...)
	}

	// 添加级别前缀
	prefix := mathx.IF(l.colorful, levelPrefixesColor[level], levelPrefixes[level])
	buf = append(buf, prefix...)

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
			buf = append(buf, '[')
			buf = append(buf, convert.S2B(file)...)
			buf = append(buf, ':')
			buf = convert.FastAppendInt(buf, line)
			buf = append(buf, ':')
			buf = append(buf, convert.S2B(funcName)...)
			buf = append(buf, ']', ' ')
		}
	}

	// 添加消息
	buf = append(buf, convert.S2B(msg)...)
	buf = append(buf, newline...)

	// 写入输出
	l.mu.Lock()
	l.output.Write(buf)
	l.mu.Unlock()

	if level == FATAL {
		os.Exit(1)
	}
}

// ultraLogf 极致优化的格式化日志方法
func (l *Logger) ultraLogf(level LogLevel, format string, args ...any) {
	if level < l.level {
		return
	}

	// 快速路径：无参数格式化
	if len(args) == 0 {
		l.ultraLog(level, format)
		return
	}

	// 有参数时才进行格式化
	msg := fmt.Sprintf(format, args...)
	l.ultraLog(level, msg)
}

// log 记录日志 - 使用 ultraLogf 提升性能
func (l *Logger) log(level LogLevel, format string, args ...any) {
	l.ultraLogf(level, format, args...)
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...any) {
	if l.level > DEBUG {
		return
	}
	l.ultraLogf(DEBUG, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...any) {
	if l.level > INFO {
		return
	}
	l.ultraLogf(INFO, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...any) {
	if l.level > WARN {
		return
	}
	l.ultraLogf(WARN, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...any) {
	if l.level > ERROR {
		return
	}
	l.ultraLogf(ERROR, format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(format string, args ...any) {
	l.ultraLogf(FATAL, format, args...)
}

// Printf风格方法（与上面相同，但命名更明确）
func (l *Logger) Debugf(format string, args ...any) {
	if l.level > DEBUG {
		return
	}
	l.ultraLogf(DEBUG, format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	if l.level > INFO {
		return
	}
	l.ultraLogf(INFO, format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	if l.level > WARN {
		return
	}
	l.ultraLogf(WARN, format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	if l.level > ERROR {
		return
	}
	l.ultraLogf(ERROR, format, args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.ultraLogf(FATAL, format, args...)
}

// WithField 添加字段信息（结构化日志）
func (l *Logger) WithField(key string, value any) ILogger {
	return &fieldLogger{
		logger: l,
		fields: map[string]any{key: value},
	}
}

// WithFields 添加多个字段信息（结构化日志）
func (l *Logger) WithFields(fields map[string]any) ILogger {
	if len(fields) == 0 {
		return l
	}

	return &fieldLogger{
		logger: l,
		fields: fields,
	}
}

// WithError 添加错误信息
func (l *Logger) WithError(err error) ILogger {
	return l.WithField("error", err.Error())
}

// ============================================================================
// Logger 实例方法
// ============================================================================

// 纯文本日志方法
func (l *Logger) DebugMsg(msg string) {
	if l.level > DEBUG {
		return
	}
	l.ultraLog(DEBUG, msg)
}

func (l *Logger) InfoMsg(msg string) {
	if l.level > INFO {
		return
	}
	l.ultraLog(INFO, msg)
}

func (l *Logger) WarnMsg(msg string) {
	if l.level > WARN {
		return
	}
	l.ultraLog(WARN, msg)
}

func (l *Logger) ErrorMsg(msg string) {
	if l.level > ERROR {
		return
	}
	l.ultraLog(ERROR, msg)
}

func (l *Logger) FatalMsg(msg string) {
	l.ultraLog(FATAL, msg)
}

// 多行日志方法 - 自动处理换行符
func (l *Logger) InfoLines(lines ...string) {
	if l.level > INFO {
		return
	}
	for _, line := range lines {
		l.ultraLog(INFO, line)
	}
}

func (l *Logger) ErrorLines(lines ...string) {
	if l.level > ERROR {
		return
	}
	for _, line := range lines {
		l.ultraLog(ERROR, line)
	}
}

func (l *Logger) WarnLines(lines ...string) {
	if l.level > WARN {
		return
	}
	for _, line := range lines {
		l.ultraLog(WARN, line)
	}
}

func (l *Logger) DebugLines(lines ...string) {
	if l.level > DEBUG {
		return
	}
	for _, line := range lines {
		l.ultraLog(DEBUG, line)
	}
}

// SetContextExtractor 设置自定义上下文提取器
func (l *Logger) SetContextExtractor(extractor ContextExtractor) {
	if extractor == nil {
		l.contextExtractor = defaultContextExtractor
	} else {
		l.contextExtractor = extractor
	}
}

// GetContextExtractor 获取当前的上下文提取器
func (l *Logger) GetContextExtractor() ContextExtractor {
	if l.contextExtractor == nil {
		return defaultContextExtractor
	}
	return l.contextExtractor
}

// extractContextInfo 从上下文中提取信息
func (l *Logger) extractContextInfo(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if l.contextExtractor == nil {
		return defaultContextExtractor(ctx)
	}
	return l.contextExtractor(ctx)
}

// 带上下文的日志方法
func (l *Logger) DebugContext(ctx context.Context, format string, args ...any) {
	if l.level > DEBUG {
		return
	}
	contextInfo := l.extractContextInfo(ctx)
	if contextInfo != "" {
		format = contextInfo + format
	}
	l.ultraLogf(DEBUG, format, args...)
}

func (l *Logger) InfoContext(ctx context.Context, format string, args ...any) {
	if l.level > INFO {
		return
	}
	contextInfo := l.extractContextInfo(ctx)
	if contextInfo != "" {
		format = contextInfo + format
	}
	l.ultraLogf(INFO, format, args...)
}

func (l *Logger) WarnContext(ctx context.Context, format string, args ...any) {
	if l.level > WARN {
		return
	}
	contextInfo := l.extractContextInfo(ctx)
	if contextInfo != "" {
		format = contextInfo + format
	}
	l.ultraLogf(WARN, format, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, format string, args ...any) {
	if l.level > ERROR {
		return
	}
	contextInfo := l.extractContextInfo(ctx)
	if contextInfo != "" {
		format = contextInfo + format
	}
	l.ultraLogf(ERROR, format, args...)
}

func (l *Logger) FatalContext(ctx context.Context, format string, args ...any) {
	contextInfo := l.extractContextInfo(ctx)
	if contextInfo != "" {
		format = contextInfo + format
	}
	l.ultraLogf(FATAL, format, args...)
}

// ============================================================================
// 键值对和结构化日志辅助方法
// ============================================================================

// logWithKV 极简键值对实现 - 零分配优化
func (l *Logger) logWithKV(level LogLevel, msg string, keysAndValues ...any) {
	if level < l.level {
		return
	}

	if len(keysAndValues) == 0 {
		l.ultraLog(level, msg)
		return
	}

	// 检查是否是单个对象参数
	if len(keysAndValues) == 1 {
		if objFields := convert.ParseObjectToMap(keysAndValues[0]); objFields != nil {
			l.logWithFields(level, msg, objFields)
			return
		}
	}

	// 快速构建带键值对的消息
	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	buf = append(buf, convert.S2B(msg)...)
	buf = append(buf, kvBraceOpen...)

	for i := 0; i < len(keysAndValues); i += 2 {
		if i > 0 {
			buf = append(buf, kvDelimiter...)
		}

		// 键
		buf = convert.AppendValue(buf, keysAndValues[i])
		buf = append(buf, kvSeparator...)

		// 值
		if i+1 < len(keysAndValues) {
			buf = convert.AppendValue(buf, keysAndValues[i+1])
		} else {
			buf = append(buf, kvMissing...)
		}
	}

	buf = append(buf, kvBraceClose...)
	l.ultraLog(level, string(buf))
}

// logWithFields 使用字段映射记录日志
func (l *Logger) logWithFields(level LogLevel, msg string, fields map[string]any) {
	if level < l.level {
		return
	}

	if len(fields) == 0 {
		l.ultraLog(level, msg)
		return
	}

	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	buf = append(buf, convert.S2B(msg)...)
	buf = append(buf, kvBraceOpen...)

	first := true
	for k, v := range fields {
		if !first {
			buf = append(buf, kvDelimiter...)
		}
		buf = append(buf, convert.S2B(k)...)
		buf = append(buf, kvSeparator...)
		buf = convert.AppendValue(buf, v)
		first = false
	}

	buf = append(buf, kvBraceClose...)
	l.ultraLog(level, string(buf))
}

// logWithContextKV 带上下文的键值对日志
func (l *Logger) logWithContextKV(ctx context.Context, level LogLevel, msg string, keysAndValues ...any) {
	if level < l.level {
		return
	}

	// 先从context提取信息
	contextInfo := l.extractContextInfo(ctx)
	if contextInfo != "" {
		msg = contextInfo + msg
	}

	l.logWithKV(level, msg, keysAndValues...)
}

// ============================================================================
// 结构化日志方法（键值对）
// ============================================================================

// 结构化日志方法（键值对）
func (l *Logger) DebugKV(msg string, keysAndValues ...any) {
	if l.level > DEBUG {
		return
	}
	l.logWithKV(DEBUG, msg, keysAndValues...)
}

func (l *Logger) DebugContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if l.level > DEBUG {
		return
	}
	l.logWithContextKV(ctx, DEBUG, msg, keysAndValues...)
}

func (l *Logger) InfoKV(msg string, keysAndValues ...any) {
	if l.level > INFO {
		return
	}
	l.logWithKV(INFO, msg, keysAndValues...)
}

func (l *Logger) InfoContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if l.level > INFO {
		return
	}
	l.logWithContextKV(ctx, INFO, msg, keysAndValues...)
}

func (l *Logger) WarnKV(msg string, keysAndValues ...any) {
	if l.level > WARN {
		return
	}
	l.logWithKV(WARN, msg, keysAndValues...)
}

func (l *Logger) WarnContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if l.level > WARN {
		return
	}
	l.logWithContextKV(ctx, WARN, msg, keysAndValues...)
}

func (l *Logger) ErrorKV(msg string, keysAndValues ...any) {
	if l.level > ERROR {
		return
	}
	l.logWithKV(ERROR, msg, keysAndValues...)
}

func (l *Logger) ErrorContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if l.level > ERROR {
		return
	}
	l.logWithContextKV(ctx, ERROR, msg, keysAndValues...)
}

func (l *Logger) FatalKV(msg string, keysAndValues ...any) {
	l.logWithKV(FATAL, msg, keysAndValues...)
}

func (l *Logger) FatalContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	l.logWithContextKV(ctx, FATAL, msg, keysAndValues...)
}

// 字段映射方法（直接支持 map[string]any）
func (l *Logger) DebugWithFields(msg string, fields map[string]any) {
	if l.level > DEBUG {
		return
	}
	l.logWithFields(DEBUG, msg, fields)
}

func (l *Logger) InfoWithFields(msg string, fields map[string]any) {
	if l.level > INFO {
		return
	}
	l.logWithFields(INFO, msg, fields)
}

func (l *Logger) WarnWithFields(msg string, fields map[string]any) {
	if l.level > WARN {
		return
	}
	l.logWithFields(WARN, msg, fields)
}

func (l *Logger) ErrorWithFields(msg string, fields map[string]any) {
	if l.level > ERROR {
		return
	}
	l.logWithFields(ERROR, msg, fields)
}

func (l *Logger) FatalWithFields(msg string, fields map[string]any) {
	l.logWithFields(FATAL, msg, fields)
}

// 原始日志条目方法
func (l *Logger) Log(level LogLevel, msg string) {
	if level < l.level {
		return
	}
	l.ultraLog(level, msg)
}

func (l *Logger) LogContext(ctx context.Context, level LogLevel, msg string) {
	if level < l.level {
		return
	}
	contextInfo := l.extractContextInfo(ctx)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	l.ultraLog(level, msg)
}

func (l *Logger) LogKV(level LogLevel, msg string, keysAndValues ...any) {
	if level < l.level {
		return
	}
	l.logWithKV(level, msg, keysAndValues...)
}

func (l *Logger) LogWithFields(level LogLevel, msg string, fields map[string]any) {
	if level < l.level {
		return
	}
	l.logWithFields(level, msg, fields)
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
func (l *Logger) Print(args ...any) {
	l.Info("%s", fmt.Sprint(args...))
}

func (l *Logger) Printf(format string, args ...any) {
	l.Info(format, args...)
}

func (l *Logger) Println(args ...any) {
	l.Info("%s", fmt.Sprintln(args...))
}

// ============================================================================
// 返回错误的日志方法
// ============================================================================

// DebugReturn 记录调试日志并返回格式化的错误
func (l *Logger) DebugReturn(format string, args ...any) error {
	l.log(DEBUG, format, args...)
	return fmt.Errorf(format, args...)
}

// InfoReturn 记录信息日志并返回格式化的错误
func (l *Logger) InfoReturn(format string, args ...any) error {
	l.log(INFO, format, args...)
	return fmt.Errorf(format, args...)
}

// WarnReturn 记录警告日志并返回格式化的错误
func (l *Logger) WarnReturn(format string, args ...any) error {
	l.log(WARN, format, args...)
	return fmt.Errorf(format, args...)
}

// ErrorReturn 记录错误日志并返回格式化的错误
func (l *Logger) ErrorReturn(format string, args ...any) error {
	l.log(ERROR, format, args...)
	return fmt.Errorf(format, args...)
}

// DebugCtxReturn 记录带上下文的调试日志并返回格式化的错误
func (l *Logger) DebugCtxReturn(ctx context.Context, format string, args ...any) error {
	l.DebugContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

// InfoCtxReturn 记录带上下文的信息日志并返回格式化的错误
func (l *Logger) InfoCtxReturn(ctx context.Context, format string, args ...any) error {
	l.InfoContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

// WarnCtxReturn 记录带上下文的警告日志并返回格式化的错误
func (l *Logger) WarnCtxReturn(ctx context.Context, format string, args ...any) error {
	l.WarnContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

// ErrorCtxReturn 记录带上下文的错误日志并返回格式化的错误
func (l *Logger) ErrorCtxReturn(ctx context.Context, format string, args ...any) error {
	l.ErrorContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

// DebugKVReturn 记录带键值对的调试日志并返回错误
func (l *Logger) DebugKVReturn(msg string, keysAndValues ...any) error {
	l.DebugKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

// InfoKVReturn 记录带键值对的信息日志并返回错误
func (l *Logger) InfoKVReturn(msg string, keysAndValues ...any) error {
	l.InfoKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

// WarnKVReturn 记录带键值对的警告日志并返回错误
func (l *Logger) WarnKVReturn(msg string, keysAndValues ...any) error {
	l.WarnKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

// ErrorKVReturn 记录带键值对的错误日志并返回错误
func (l *Logger) ErrorKVReturn(msg string, keysAndValues ...any) error {
	l.ErrorKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

// ============================================================================
// Console 风格日志方法实现
// ============================================================================

// getOrCreateConsoleGroup 获取或创建 ConsoleGroup（延迟初始化）
func (l *Logger) getOrCreateConsoleGroup() *ConsoleGroup {
	l.consoleGroupOnce.Do(func() {
		l.consoleGroup = &ConsoleGroup{
			logger:          l,
			indentLevel:     0,
			collapsed:       false,
			collapsedLevels: make([]bool, 0, 16), // 预分配 16 层嵌套容量
		}
	})
	return l.consoleGroup
}

// ConsoleGroup 开始一个新的日志分组
func (l *Logger) ConsoleGroup(label string, args ...any) {
	cg := l.getOrCreateConsoleGroup()
	cg.Group(label, args...)
}

// ConsoleGroupCollapsed 开始一个折叠的日志分组
func (l *Logger) ConsoleGroupCollapsed(label string, args ...any) {
	cg := l.getOrCreateConsoleGroup()
	cg.GroupCollapsed(label, args...)
}

// ConsoleGroupEnd 结束当前分组
func (l *Logger) ConsoleGroupEnd() {
	cg := l.getOrCreateConsoleGroup()
	cg.GroupEnd()
}

// ConsoleTable 显示表格
func (l *Logger) ConsoleTable(data any) {
	cg := l.getOrCreateConsoleGroup()
	cg.Table(data)
}

// ConsoleTime 开始计时
func (l *Logger) ConsoleTime(label string) *Timer {
	cg := l.getOrCreateConsoleGroup()
	return cg.Time(label)
}

// ============================================================================
// 配置方法实现
// ============================================================================

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

// ============================================================================
// fieldLogger - 字段日志包装器（用于 WithField/WithFields）
// ============================================================================

// fieldLogger 轻量级字段日志器包装
type fieldLogger struct {
	logger *Logger
	fields map[string]any
}

// 实现所有 ILogger 接口方法，将字段附加到日志消息

// Debug 调试日志
func (f *fieldLogger) Debug(format string, args ...any) {
	if !f.logger.IsLevelEnabled(DEBUG) {
		return
	}
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	f.logger.logWithFields(DEBUG, msg, f.fields)
}

// Info 信息日志
func (f *fieldLogger) Info(format string, args ...any) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	f.logger.logWithFields(INFO, msg, f.fields)
}

// Warn 警告日志
func (f *fieldLogger) Warn(format string, args ...any) {
	if !f.logger.IsLevelEnabled(WARN) {
		return
	}
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	f.logger.logWithFields(WARN, msg, f.fields)
}

// Error 错误日志
func (f *fieldLogger) Error(format string, args ...any) {
	if !f.logger.IsLevelEnabled(ERROR) {
		return
	}
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	f.logger.logWithFields(ERROR, msg, f.fields)
}

// Fatal 致命错误日志
func (f *fieldLogger) Fatal(format string, args ...any) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	f.logger.logWithFields(FATAL, msg, f.fields)
}

// Printf风格方法
func (f *fieldLogger) Debugf(format string, args ...any) {
	f.Debug(format, args...)
}

func (f *fieldLogger) Infof(format string, args ...any) {
	f.Info(format, args...)
}

func (f *fieldLogger) Warnf(format string, args ...any) {
	f.Warn(format, args...)
}

func (f *fieldLogger) Errorf(format string, args ...any) {
	f.Error(format, args...)
}

func (f *fieldLogger) Fatalf(format string, args ...any) {
	f.Fatal(format, args...)
}

// 纯文本日志方法
func (f *fieldLogger) DebugMsg(msg string) {
	if f.logger.level > DEBUG {
		return
	}
	f.logger.logWithFields(DEBUG, msg, f.fields)
}

func (f *fieldLogger) InfoMsg(msg string) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	f.logger.logWithFields(INFO, msg, f.fields)
}

func (f *fieldLogger) WarnMsg(msg string) {
	if !f.logger.IsLevelEnabled(WARN) {
		return
	}
	f.logger.logWithFields(WARN, msg, f.fields)
}

func (f *fieldLogger) ErrorMsg(msg string) {
	if !f.logger.IsLevelEnabled(ERROR) {
		return
	}
	f.logger.logWithFields(ERROR, msg, f.fields)
}

func (f *fieldLogger) FatalMsg(msg string) {
	f.logger.logWithFields(FATAL, msg, f.fields)
}

// 多行日志方法
func (f *fieldLogger) InfoLines(lines ...string) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	for _, line := range lines {
		f.logger.logWithFields(INFO, line, f.fields)
	}
}

func (f *fieldLogger) ErrorLines(lines ...string) {
	if !f.logger.IsLevelEnabled(ERROR) {
		return
	}
	for _, line := range lines {
		f.logger.logWithFields(ERROR, line, f.fields)
	}
}

func (f *fieldLogger) WarnLines(lines ...string) {
	if !f.logger.IsLevelEnabled(WARN) {
		return
	}
	for _, line := range lines {
		f.logger.logWithFields(WARN, line, f.fields)
	}
}

func (f *fieldLogger) DebugLines(lines ...string) {
	if !f.logger.IsLevelEnabled(DEBUG) {
		return
	}
	for _, line := range lines {
		f.logger.logWithFields(DEBUG, line, f.fields)
	}
}

// 上下文日志方法
func (f *fieldLogger) DebugContext(ctx context.Context, format string, args ...any) {
	if f.logger.level > DEBUG {
		return
	}
	contextInfo := f.logger.extractContextInfo(ctx)
	msg := fmt.Sprintf(format, args...)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	f.logger.logWithFields(DEBUG, msg, f.fields)
}

func (f *fieldLogger) InfoContext(ctx context.Context, format string, args ...any) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	contextInfo := f.logger.extractContextInfo(ctx)
	msg := fmt.Sprintf(format, args...)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	f.logger.logWithFields(INFO, msg, f.fields)
}

func (f *fieldLogger) WarnContext(ctx context.Context, format string, args ...any) {
	if !f.logger.IsLevelEnabled(WARN) {
		return
	}
	contextInfo := f.logger.extractContextInfo(ctx)
	msg := fmt.Sprintf(format, args...)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	f.logger.logWithFields(WARN, msg, f.fields)
}

func (f *fieldLogger) ErrorContext(ctx context.Context, format string, args ...any) {
	if !f.logger.IsLevelEnabled(ERROR) {
		return
	}
	contextInfo := f.logger.extractContextInfo(ctx)
	msg := fmt.Sprintf(format, args...)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	f.logger.logWithFields(ERROR, msg, f.fields)
}

func (f *fieldLogger) FatalContext(ctx context.Context, format string, args ...any) {
	contextInfo := f.logger.extractContextInfo(ctx)
	msg := fmt.Sprintf(format, args...)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	f.logger.logWithFields(FATAL, msg, f.fields)
}

// 键值对日志方法
func (f *fieldLogger) DebugKV(msg string, keysAndValues ...any) {
	if !f.logger.IsLevelEnabled(DEBUG) {
		return
	}
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithKV(DEBUG, msg, allFields...)
}

func (f *fieldLogger) InfoKV(msg string, keysAndValues ...any) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithKV(INFO, msg, allFields...)
}

func (f *fieldLogger) WarnKV(msg string, keysAndValues ...any) {
	if !f.logger.IsLevelEnabled(WARN) {
		return
	}
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithKV(WARN, msg, allFields...)
}

func (f *fieldLogger) ErrorKV(msg string, keysAndValues ...any) {
	if !f.logger.IsLevelEnabled(ERROR) {
		return
	}
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithKV(ERROR, msg, allFields...)
}

func (f *fieldLogger) FatalKV(msg string, keysAndValues ...any) {
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithKV(FATAL, msg, allFields...)
}

// 带上下文的键值对日志方法
func (f *fieldLogger) DebugContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if !f.logger.IsLevelEnabled(DEBUG) {
		return
	}
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithContextKV(ctx, DEBUG, msg, allFields...)
}

func (f *fieldLogger) InfoContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithContextKV(ctx, INFO, msg, allFields...)
}

func (f *fieldLogger) WarnContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if !f.logger.IsLevelEnabled(WARN) {
		return
	}
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithContextKV(ctx, WARN, msg, allFields...)
}

func (f *fieldLogger) ErrorContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if !f.logger.IsLevelEnabled(ERROR) {
		return
	}
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithContextKV(ctx, ERROR, msg, allFields...)
}

func (f *fieldLogger) FatalContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithContextKV(ctx, FATAL, msg, allFields...)
}

// 字段映射方法
func (f *fieldLogger) DebugWithFields(msg string, fields map[string]any) {
	if !f.logger.IsLevelEnabled(DEBUG) {
		return
	}
	mergedFields := f.mergeFieldsMap(fields)
	f.logger.logWithFields(DEBUG, msg, mergedFields)
}

func (f *fieldLogger) InfoWithFields(msg string, fields map[string]any) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	mergedFields := f.mergeFieldsMap(fields)
	f.logger.logWithFields(INFO, msg, mergedFields)
}

func (f *fieldLogger) WarnWithFields(msg string, fields map[string]any) {
	if !f.logger.IsLevelEnabled(WARN) {
		return
	}
	mergedFields := f.mergeFieldsMap(fields)
	f.logger.logWithFields(WARN, msg, mergedFields)
}

func (f *fieldLogger) ErrorWithFields(msg string, fields map[string]any) {
	if !f.logger.IsLevelEnabled(ERROR) {
		return
	}
	mergedFields := f.mergeFieldsMap(fields)
	f.logger.logWithFields(ERROR, msg, mergedFields)
}

func (f *fieldLogger) FatalWithFields(msg string, fields map[string]any) {
	mergedFields := f.mergeFieldsMap(fields)
	f.logger.logWithFields(FATAL, msg, mergedFields)
}

// 原始日志条目方法
func (f *fieldLogger) Log(level LogLevel, msg string) {
	if !f.logger.IsLevelEnabled(level) {
		return
	}
	f.logger.logWithFields(level, msg, f.fields)
}

func (f *fieldLogger) LogContext(ctx context.Context, level LogLevel, msg string) {
	if !f.logger.IsLevelEnabled(level) {
		return
	}
	contextInfo := f.logger.extractContextInfo(ctx)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	f.logger.logWithFields(level, msg, f.fields)
}

func (f *fieldLogger) LogKV(level LogLevel, msg string, keysAndValues ...any) {
	if !f.logger.IsLevelEnabled(level) {
		return
	}
	allFields := f.mergeKV(keysAndValues...)
	f.logger.logWithKV(level, msg, allFields...)
}

func (f *fieldLogger) LogWithFields(level LogLevel, msg string, fields map[string]any) {
	if !f.logger.IsLevelEnabled(level) {
		return
	}
	mergedFields := f.mergeFieldsMap(fields)
	f.logger.logWithFields(level, msg, mergedFields)
}

// 配置方法
func (f *fieldLogger) SetLevel(level LogLevel) {
	f.logger.SetLevel(level)
}

func (f *fieldLogger) GetLevel() LogLevel {
	return f.logger.GetLevel()
}

func (f *fieldLogger) SetShowCaller(show bool) {
	f.logger.SetShowCaller(show)
}

func (f *fieldLogger) IsShowCaller() bool {
	return f.logger.IsShowCaller()
}

func (f *fieldLogger) IsLevelEnabled(level LogLevel) bool {
	return f.logger.IsLevelEnabled(level)
}

// 结构化日志构建器 - 使用 clear() 优化（Go 1.21+）
func (f *fieldLogger) WithField(key string, value any) ILogger {
	// 从对象池获取 map
	newFields := fieldMapPool.Get().(map[string]any)

	clear(newFields)

	// 复制现有字段
	for k, v := range f.fields {
		newFields[k] = v
	}
	newFields[key] = value

	return &fieldLogger{logger: f.logger, fields: newFields}
}

func (f *fieldLogger) WithFields(fields map[string]any) ILogger {
	if len(fields) == 0 {
		return f
	}

	// 从对象池获取 map
	newFields := fieldMapPool.Get().(map[string]any)

	clear(newFields)

	// 复制现有字段
	for k, v := range f.fields {
		newFields[k] = v
	}
	// 添加新字段
	for k, v := range fields {
		newFields[k] = v
	}

	return &fieldLogger{logger: f.logger, fields: newFields}
}

func (f *fieldLogger) WithError(err error) ILogger {
	return f.WithField("error", err.Error())
}

func (f *fieldLogger) WithContext(ctx context.Context) ILogger {
	return f
}

// Clone 克隆当前Logger
func (f *fieldLogger) Clone() ILogger {
	newFields := make(map[string]any, len(f.fields))
	for k, v := range f.fields {
		newFields[k] = v
	}
	return &fieldLogger{logger: f.logger, fields: newFields}
}

// 兼容标准log包的方法
func (f *fieldLogger) Print(args ...any) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	f.logger.logWithFields(INFO, fmt.Sprint(args...), f.fields)
}

func (f *fieldLogger) Printf(format string, args ...any) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	f.logger.logWithFields(INFO, fmt.Sprintf(format, args...), f.fields)
}

func (f *fieldLogger) Println(args ...any) {
	if !f.logger.IsLevelEnabled(INFO) {
		return
	}
	msg := fmt.Sprintln(args...)
	f.logger.logWithFields(INFO, msg[:len(msg)-1], f.fields)
}

// 返回错误的日志方法
func (f *fieldLogger) DebugReturn(format string, args ...any) error {
	f.Debug(format, args...)
	return fmt.Errorf(format, args...)
}

func (f *fieldLogger) InfoReturn(format string, args ...any) error {
	f.Info(format, args...)
	return fmt.Errorf(format, args...)
}

func (f *fieldLogger) WarnReturn(format string, args ...any) error {
	f.Warn(format, args...)
	return fmt.Errorf(format, args...)
}

func (f *fieldLogger) ErrorReturn(format string, args ...any) error {
	f.Error(format, args...)
	return fmt.Errorf(format, args...)
}

// 返回错误的上下文日志方法
func (f *fieldLogger) DebugCtxReturn(ctx context.Context, format string, args ...any) error {
	f.DebugContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

func (f *fieldLogger) InfoCtxReturn(ctx context.Context, format string, args ...any) error {
	f.InfoContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

func (f *fieldLogger) WarnCtxReturn(ctx context.Context, format string, args ...any) error {
	f.WarnContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

func (f *fieldLogger) ErrorCtxReturn(ctx context.Context, format string, args ...any) error {
	f.ErrorContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

// 返回错误的键值对日志方法
func (f *fieldLogger) DebugKVReturn(msg string, keysAndValues ...any) error {
	f.DebugKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

func (f *fieldLogger) InfoKVReturn(msg string, keysAndValues ...any) error {
	f.InfoKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

func (f *fieldLogger) WarnKVReturn(msg string, keysAndValues ...any) error {
	f.WarnKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

func (f *fieldLogger) ErrorKVReturn(msg string, keysAndValues ...any) error {
	f.ErrorKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

// Console 相关方法
func (f *fieldLogger) ConsoleGroup(label string, args ...any) {
	f.logger.ConsoleGroup(label, args...)
}

func (f *fieldLogger) ConsoleGroupCollapsed(label string, args ...any) {
	f.logger.ConsoleGroupCollapsed(label, args...)
}

func (f *fieldLogger) ConsoleGroupEnd() {
	f.logger.ConsoleGroupEnd()
}

func (f *fieldLogger) ConsoleTable(data any) {
	f.logger.ConsoleTable(data)
}

func (f *fieldLogger) ConsoleTime(label string) *Timer {
	return f.logger.ConsoleTime(label)
}

func (f *fieldLogger) NewConsoleGroup() *ConsoleGroup {
	return f.logger.getOrCreateConsoleGroup()
}

// 辅助方法：合并字段和键值对
func (f *fieldLogger) mergeKV(keysAndValues ...any) []any {
	totalLen := len(f.fields)*2 + len(keysAndValues)
	result := make([]any, 0, totalLen)

	// 添加现有字段
	for k, v := range f.fields {
		result = append(result, k, v)
	}

	// 添加传入的键值对
	result = append(result, keysAndValues...)

	return result
}

// 辅助方法：合并字段映射 - 使用对象池优化
func (f *fieldLogger) mergeFieldsMap(fields map[string]any) map[string]any {
	if len(fields) == 0 {
		return f.fields
	}

	// 从对象池获取 map
	merged := fieldMapPool.Get().(map[string]any)
	clear(merged)

	// 添加现有字段
	for k, v := range f.fields {
		merged[k] = v
	}

	// 添加传入的字段
	for k, v := range fields {
		merged[k] = v
	}

	return merged
}

// ============================================================================
// 特殊场景日志方法（specialty）
// ============================================================================

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

// logSpecial 记录特殊类型的日志（使用 INFO 级别）
func (l *Logger) logSpecial(logType SpecialLogType, level LogLevel, format string, args ...any) {
	if level < l.level {
		return
	}
	message := fmt.Sprintf(format, args...)
	l.ultraLog(level, fmt.Sprintf("%s [%s] %s", logType.emoji, logType.name, message))
}

// Success 成功日志（INFO 级别）
func (l *Logger) Success(format string, args ...any) {
	l.logSpecial(SuccessType, INFO, format, args...)
}

// Loading 加载日志（INFO 级别）
func (l *Logger) Loading(format string, args ...any) {
	l.logSpecial(LoadingType, INFO, format, args...)
}

// ConfigLog 配置日志（INFO 级别）
func (l *Logger) ConfigLog(format string, args ...any) {
	l.logSpecial(ConfigType, INFO, format, args...)
}

// Start 启动日志（INFO 级别）
func (l *Logger) Start(format string, args ...any) {
	l.logSpecial(StartType, INFO, format, args...)
}

// Stop 停止日志（INFO 级别）
func (l *Logger) Stop(format string, args ...any) {
	l.logSpecial(StopType, INFO, format, args...)
}

// Database 数据库日志（INFO 级别）
func (l *Logger) Database(format string, args ...any) {
	l.logSpecial(DatabaseType, INFO, format, args...)
}

// Network 网络日志（INFO 级别）
func (l *Logger) Network(format string, args ...any) {
	l.logSpecial(NetworkType, INFO, format, args...)
}

// Security 安全日志（SECURITY 级别）
func (l *Logger) Security(format string, args ...any) {
	l.logSpecial(SecurityType, SECURITY, format, args...)
}

// Cache 缓存日志（INFO 级别）
func (l *Logger) Cache(format string, args ...any) {
	l.logSpecial(CacheType, INFO, format, args...)
}

// Environment 环境日志（INFO 级别）
func (l *Logger) Environment(format string, args ...any) {
	l.logSpecial(EnvironmentType, INFO, format, args...)
}

// ============================================================================
// 性能日志方法
// ============================================================================

// PerformanceLevel 性能级别定义
type PerformanceLevel struct {
	Threshold time.Duration
	Emoji     string
	Level     string
}

// 性能级别配置（按阈值从小到大排序）
var performanceLevels = []PerformanceLevel{
	{50 * time.Millisecond, "⚡", "EXCELLENT"},
	{100 * time.Millisecond, "🏃", "FAST"},
	{500 * time.Millisecond, "🚶", "NORMAL"},
	{2 * time.Second, "🐢", "SLOW"},
	{0, "🐌", "VERY_SLOW"}, // 0 表示默认值（最后一个）
}

// getPerformanceLevel 获取性能级别和表情符号
func getPerformanceLevel(duration time.Duration) (emoji, level string) {
	for _, pl := range performanceLevels {
		if pl.Threshold == 0 || duration < pl.Threshold {
			return pl.Emoji, pl.Level
		}
	}
	return "🐌", "VERY_SLOW"
}

// Performance 性能日志（PERFORMANCE 级别，支持可选的详细信息）
func (l *Logger) Performance(operation string, duration time.Duration, details ...map[string]any) {
	if PERFORMANCE < l.level {
		return
	}

	emoji, level := getPerformanceLevel(duration)
	msg := fmt.Sprintf("%s [PERF-%s] %s completed in %v", emoji, level, operation, duration)

	if len(details) > 0 && len(details[0]) > 0 {
		msg += fmt.Sprintf(" | Details: %+v", details[0])
	}

	l.ultraLog(PERFORMANCE, msg)
}

// Timing 计时器辅助结构
type Timing struct {
	logger    *Logger
	operation string
	startTime time.Time
	details   map[string]any
}

// StartTiming 开始计时
func (l *Logger) StartTiming(operation string) *Timing {
	return &Timing{
		logger:    l,
		operation: operation,
		startTime: time.Now(),
		details:   make(map[string]any),
	}
}

// AddDetail 添加详细信息
func (t *Timing) AddDetail(key string, value any) *Timing {
	t.details[key] = value
	return t
}

// End 结束计时并记录性能日志
func (t *Timing) End() time.Duration {
	duration := time.Since(t.startTime)
	if len(t.details) > 0 {
		t.logger.Performance(t.operation, duration, t.details)
	} else {
		t.logger.Performance(t.operation, duration)
	}
	return duration
}

// getProgressEmoji 根据进度百分比获取表情符号
func getProgressEmoji(percentage float64) string {
	switch {
	case percentage == 100:
		return "✅"
	case percentage >= 75:
		return "🔵"
	case percentage >= 50:
		return "🟡"
	case percentage >= 25:
		return "🟠"
	default:
		return "🔴"
	}
}

// Progress 进度日志（INFO 级别）
func (l *Logger) Progress(current, total int, operation string) {
	if INFO < l.level {
		return
	}

	percentage := float64(current) / float64(total) * 100
	emoji := getProgressEmoji(percentage)
	l.ultraLog(INFO, fmt.Sprintf("%s [PROGRESS] %s: %d/%d (%.1f%%)", emoji, operation, current, total, percentage))
}

// Milestone 里程碑日志（INFO 级别）
func (l *Logger) Milestone(message string) {
	if INFO < l.level {
		return
	}
	l.ultraLog(INFO, fmt.Sprintf("🎯 [MILESTONE] %s", message))
}

// Health 健康检查日志（WARN 级别用于不健康，INFO 级别用于健康）
func (l *Logger) Health(service string, status bool, details string) {
	level := WARN
	if status {
		level = INFO
	}

	if level < l.level {
		return
	}

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

	l.ultraLog(level, fmt.Sprintf("%s [HEALTH] %s: %s%s", emoji, service, statusStr, detailStr))
}

// Audit 审计日志（AUDIT 级别）
func (l *Logger) Audit(action, user, resource, result string) {
	if AUDIT < l.level {
		return
	}
	l.ultraLog(AUDIT, fmt.Sprintf("📋 [AUDIT] User: %s | Action: %s | Resource: %s | Result: %s", user, action, resource, result))
}
