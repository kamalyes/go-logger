/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-09 16:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 10:19:58
 * @FilePath: \go-logger\ultra_fast_logger.go
 * @Description: 极致性能优化的日志实现
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

// 极致优化常量
const (
	maxLogMessageSize = 512
)

// 预分配的字节池
var (
	bytePool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, maxLogMessageSize)
		},
	}
)

// 预计算的常量
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
	space   = []byte(" ")
)

// levelPrefixes 预计算的级别前缀
var levelPrefixes = map[LogLevel][]byte{
	DEBUG: debugPrefix,
	INFO:  infoPrefix,
	WARN:  warnPrefix,
	ERROR: errorPrefix,
	FATAL: fatalPrefix,
}

var levelPrefixesColor = map[LogLevel][]byte{
	DEBUG: debugPrefixColor,
	INFO:  infoPrefixColor,
	WARN:  warnPrefixColor,
	ERROR: errorPrefixColor,
	FATAL: fatalPrefixColor,
}

// UltraFastLogger 极致性能的日志器
type UltraFastLogger struct {
	level    LogLevel
	colorful bool
	output   io.Writer
	mu       sync.Mutex // 保护并发写入
	
	// 优化选项
	skipTimestamp bool // 跳过时间戳以获得极致性能
	skipCaller    bool // 跳过调用者信息
}

// NewUltraFastLogger 创建极致性能日志器
func NewUltraFastLogger(config *LogConfig) *UltraFastLogger {
	if config == nil {
		config = DefaultConfig()
	}

	return &UltraFastLogger{
		level:         config.Level,
		colorful:      config.Colorful,
		output:        config.Output,
		skipTimestamp: false, // 可配置
		skipCaller:    !config.ShowCaller,
	}
}

// NewUltraFastLoggerNoTime 创建不包含时间戳的极致性能日志器
func NewUltraFastLoggerNoTime(output io.Writer, level LogLevel) *UltraFastLogger {
	return &UltraFastLogger{
		level:         level,
		colorful:      false,
		output:        output,
		skipTimestamp: true,
		skipCaller:    true,
	}
}

// unsafeStringToBytes 零拷贝字符串到字节转换
func unsafeStringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		int
	}{s, len(s)}))
}

// unsafeBytesToString 零拷贝字节到字符串转换
func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// fastAppendInt 快速整数追加，避免 strconv.Itoa 分配
func fastAppendInt(buf []byte, val int) []byte {
	if val == 0 {
		return append(buf, '0')
	}
	
	// 快速路径：小数字
	if val < 10 {
		return append(buf, byte('0'+val))
	}
	if val < 100 {
		return append(buf, byte('0'+val/10), byte('0'+val%10))
	}
	if val < 1000 {
		return append(buf, byte('0'+val/100), byte('0'+(val/10)%10), byte('0'+val%10))
	}
	
	// 通用路径
	return strconv.AppendInt(buf, int64(val), 10)
}

// fastFormatTime 快速时间格式化，避免 time.Format 的开销
func fastFormatTime(buf []byte, t time.Time) []byte {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	
	// 手动格式化 "2006/01/02 15:04:05 "
	buf = fastAppendInt(buf, year)
	buf = append(buf, '/')
	buf = fastAppendInt(buf, int(month))
	buf = append(buf, '/')
	buf = fastAppendInt(buf, day)
	buf = append(buf, ' ')
	buf = fastAppendInt(buf, hour)
	buf = append(buf, ':')
	if min < 10 {
		buf = append(buf, '0')
	}
	buf = fastAppendInt(buf, min)
	buf = append(buf, ':')
	if sec < 10 {
		buf = append(buf, '0')
	}
	buf = fastAppendInt(buf, sec)
	buf = append(buf, ' ')
	
	return buf
}

// ultraLog 极致优化的日志方法
func (l *UltraFastLogger) ultraLog(level LogLevel, msg string) {
	// 快速级别检查
	if level < l.level {
		return
	}

	// 获取字节缓冲区
	buf := bytePool.Get().([]byte)
	buf = buf[:0] // 重置长度但保留容量
	
	defer bytePool.Put(buf)

	// 添加时间戳（如果需要）
	if !l.skipTimestamp {
		buf = fastFormatTime(buf, time.Now())
	}

	// 添加级别前缀
	var prefix []byte
	if l.colorful {
		prefix = levelPrefixesColor[level]
	} else {
		prefix = levelPrefixes[level]
	}
	buf = append(buf, prefix...)

	// 添加消息
	buf = append(buf, unsafeStringToBytes(msg)...)
	buf = append(buf, newline...)

	// 写入输出
	l.mu.Lock()
	l.output.Write(buf)
	l.mu.Unlock()

	if level == FATAL {
		os.Exit(1)
	}
}

// ultraLogf 极致优化的格式化日志方法（有限支持格式化以保持性能）
func (l *UltraFastLogger) ultraLogf(level LogLevel, format string, args ...interface{}) {
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

// 实现 ILogger 接口
func (l *UltraFastLogger) Debug(format string, args ...interface{}) {
	l.ultraLogf(DEBUG, format, args...)
}

func (l *UltraFastLogger) Info(format string, args ...interface{}) {
	l.ultraLogf(INFO, format, args...)
}

func (l *UltraFastLogger) Warn(format string, args ...interface{}) {
	l.ultraLogf(WARN, format, args...)
}

func (l *UltraFastLogger) Error(format string, args ...interface{}) {
	l.ultraLogf(ERROR, format, args...)
}

func (l *UltraFastLogger) Fatal(format string, args ...interface{}) {
	l.ultraLogf(FATAL, format, args...)
}

// 纯文本日志方法
func (l *UltraFastLogger) DebugMsg(msg string) {
	l.ultraLog(DEBUG, msg)
}

func (l *UltraFastLogger) InfoMsg(msg string) {
	l.ultraLog(INFO, msg)
}

func (l *UltraFastLogger) WarnMsg(msg string) {
	l.ultraLog(WARN, msg)
}

func (l *UltraFastLogger) ErrorMsg(msg string) {
	l.ultraLog(ERROR, msg)
}

func (l *UltraFastLogger) FatalMsg(msg string) {
	l.ultraLog(FATAL, msg)
}

// 配置方法
func (l *UltraFastLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *UltraFastLogger) GetLevel() LogLevel {
	return l.level
}

func (l *UltraFastLogger) SetShowCaller(show bool) {
	l.skipCaller = !show
}

func (l *UltraFastLogger) IsShowCaller() bool {
	return !l.skipCaller
}

func (l *UltraFastLogger) IsLevelEnabled(level LogLevel) bool {
	return level >= l.level
}

// 日志条目方法
func (l *UltraFastLogger) Log(level LogLevel, msg string) {
	l.ultraLog(level, msg)
}

func (l *UltraFastLogger) LogContext(ctx context.Context, level LogLevel, msg string) {
	l.ultraLog(level, msg) // 简化版本忽略上下文
}

func (l *UltraFastLogger) LogWithFields(level LogLevel, msg string, fields map[string]interface{}) {
	if level < l.level {
		return
	}
	
	// 快速构建字段消息
	if len(fields) == 0 {
		l.ultraLog(level, msg)
		return
	}
	
	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	buf = append(buf, unsafeStringToBytes(msg)...)
	buf = append(buf, " {"...)

	first := true
	for k, v := range fields {
		if !first {
			buf = append(buf, ", "...)
		}
		buf = append(buf, unsafeStringToBytes(k)...)
		buf = append(buf, ": "...)
		val := fmt.Sprint(v)
		buf = append(buf, unsafeStringToBytes(val)...)
		first = false
	}
	
	buf = append(buf, '}')
	l.ultraLog(level, unsafeBytesToString(buf))
}

// 简化的 Print 方法
func (l *UltraFastLogger) Print(v ...interface{}) {
	if len(v) == 1 {
		if s, ok := v[0].(string); ok {
			l.ultraLog(INFO, s)
			return
		}
	}
	msg := fmt.Sprint(v...)
	l.ultraLog(INFO, msg)
}

func (l *UltraFastLogger) Printf(format string, v ...interface{}) {
	l.ultraLogf(INFO, format, v...)
}

func (l *UltraFastLogger) Println(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	l.ultraLog(INFO, msg[:len(msg)-1]) // 移除额外的换行符
}

// 上下文方法（简化版，忽略上下文以保持性能）
func (l *UltraFastLogger) DebugContext(ctx context.Context, format string, args ...interface{}) {
	l.ultraLogf(DEBUG, format, args...)
}

func (l *UltraFastLogger) InfoContext(ctx context.Context, format string, args ...interface{}) {
	l.ultraLogf(INFO, format, args...)
}

func (l *UltraFastLogger) WarnContext(ctx context.Context, format string, args ...interface{}) {
	l.ultraLogf(WARN, format, args...)
}

func (l *UltraFastLogger) ErrorContext(ctx context.Context, format string, args ...interface{}) {
	l.ultraLogf(ERROR, format, args...)
}

func (l *UltraFastLogger) FatalContext(ctx context.Context, format string, args ...interface{}) {
	l.ultraLogf(FATAL, format, args...)
}

// 键值对方法（极简版）
func (l *UltraFastLogger) DebugKV(msg string, keysAndValues ...interface{}) {
	if l.level > DEBUG {
		return
	}
	l.logWithKV(DEBUG, msg, keysAndValues...)
}

func (l *UltraFastLogger) InfoKV(msg string, keysAndValues ...interface{}) {
	if l.level > INFO {
		return
	}
	l.logWithKV(INFO, msg, keysAndValues...)
}

func (l *UltraFastLogger) WarnKV(msg string, keysAndValues ...interface{}) {
	if l.level > WARN {
		return
	}
	l.logWithKV(WARN, msg, keysAndValues...)
}

func (l *UltraFastLogger) ErrorKV(msg string, keysAndValues ...interface{}) {
	if l.level > ERROR {
		return
	}
	l.logWithKV(ERROR, msg, keysAndValues...)
}

func (l *UltraFastLogger) FatalKV(msg string, keysAndValues ...interface{}) {
	l.logWithKV(FATAL, msg, keysAndValues...)
}

func (l *UltraFastLogger) LogKV(level LogLevel, msg string, keysAndValues ...interface{}) {
	if level < l.level {
		return
	}
	l.logWithKV(level, msg, keysAndValues...)
}

// logWithKV 极简键值对实现
func (l *UltraFastLogger) logWithKV(level LogLevel, msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) == 0 {
		l.ultraLog(level, msg)
		return
	}

	// 快速构建带键值对的消息
	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	buf = append(buf, unsafeStringToBytes(msg)...)
	buf = append(buf, " {"...)

	for i := 0; i < len(keysAndValues); i += 2 {
		if i > 0 {
			buf = append(buf, ", "...)
		}
		
		// 键
		key := fmt.Sprint(keysAndValues[i])
		buf = append(buf, unsafeStringToBytes(key)...)
		buf = append(buf, ": "...)
		
		// 值
		if i+1 < len(keysAndValues) {
			val := fmt.Sprint(keysAndValues[i+1])
			buf = append(buf, unsafeStringToBytes(val)...)
		}
	}
	
	buf = append(buf, '}')
	
	l.ultraLog(level, unsafeBytesToString(buf))
}

// 字段方法返回简化的包装器
func (l *UltraFastLogger) WithField(key string, value interface{}) ILogger {
	return &ultraFieldLogger{logger: l, key: key, value: value}
}

func (l *UltraFastLogger) WithFields(fields map[string]interface{}) ILogger {
	return &ultraFieldLogger{logger: l, fields: fields}
}

func (l *UltraFastLogger) WithError(err error) ILogger {
	return &ultraFieldLogger{logger: l, key: "error", value: err.Error()}
}

func (l *UltraFastLogger) WithContext(ctx context.Context) ILogger {
	return l // 简化版本忽略上下文
}

func (l *UltraFastLogger) Clone() ILogger {
	return &UltraFastLogger{
		level:         l.level,
		colorful:      l.colorful,
		output:        l.output,
		skipTimestamp: l.skipTimestamp,
		skipCaller:    l.skipCaller,
	}
}

// ultraFieldLogger 超轻量级字段日志器
type ultraFieldLogger struct {
	logger ILogger
	key    string
	value  interface{}
	fields map[string]interface{}
}

func (f *ultraFieldLogger) Debug(format string, args ...interface{}) {
	f.logWithFields(DEBUG, format, args...)
}

func (f *ultraFieldLogger) Info(format string, args ...interface{}) {
	f.logWithFields(INFO, format, args...)
}

func (f *ultraFieldLogger) Warn(format string, args ...interface{}) {
	f.logWithFields(WARN, format, args...)
}

func (f *ultraFieldLogger) Error(format string, args ...interface{}) {
	f.logWithFields(ERROR, format, args...)
}

func (f *ultraFieldLogger) Fatal(format string, args ...interface{}) {
	f.logWithFields(FATAL, format, args...)
}

func (f *ultraFieldLogger) logWithFields(level LogLevel, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	
	// 快速构建字段字符串
	if f.key != "" {
		msg = fmt.Sprintf("%s {%s: %v}", msg, f.key, f.value)
	} else if len(f.fields) > 0 {
		fieldsStr := ""
		first := true
		for k, v := range f.fields {
			if !first {
				fieldsStr += ", "
			}
			fieldsStr += fmt.Sprintf("%s: %v", k, v)
			first = false
		}
		msg = fmt.Sprintf("%s {%s}", msg, fieldsStr)
	}

	switch level {
	case DEBUG:
		f.logger.Debug(msg)
	case INFO:
		f.logger.Info(msg)
	case WARN:
		f.logger.Warn(msg)
	case ERROR:
		f.logger.Error(msg)
	case FATAL:
		f.logger.Fatal(msg)
	}
}

// 添加缺失的纯文本日志方法
func (f *ultraFieldLogger) DebugMsg(msg string) {
	f.logWithFieldsMsg(DEBUG, msg)
}

func (f *ultraFieldLogger) InfoMsg(msg string) {
	f.logWithFieldsMsg(INFO, msg)
}

func (f *ultraFieldLogger) WarnMsg(msg string) {
	f.logWithFieldsMsg(WARN, msg)
}

func (f *ultraFieldLogger) ErrorMsg(msg string) {
	f.logWithFieldsMsg(ERROR, msg)
}

func (f *ultraFieldLogger) FatalMsg(msg string) {
	f.logWithFieldsMsg(FATAL, msg)
}

func (f *ultraFieldLogger) logWithFieldsMsg(level LogLevel, msg string) {
	// 快速构建字段字符串
	if f.key != "" {
		msg = fmt.Sprintf("%s {%s: %v}", msg, f.key, f.value)
	} else if len(f.fields) > 0 {
		fieldsStr := ""
		first := true
		for k, v := range f.fields {
			if !first {
				fieldsStr += ", "
			}
			fieldsStr += fmt.Sprintf("%s: %v", k, v)
			first = false
		}
		msg = fmt.Sprintf("%s {%s}", msg, fieldsStr)
	}

	switch level {
	case DEBUG:
		f.logger.DebugMsg(msg)
	case INFO:
		f.logger.InfoMsg(msg)
	case WARN:
		f.logger.WarnMsg(msg)
	case ERROR:
		f.logger.ErrorMsg(msg)
	case FATAL:
		f.logger.FatalMsg(msg)
	}
}

// 添加缺失的配置方法
func (f *ultraFieldLogger) SetLevel(level LogLevel) {
	f.logger.SetLevel(level)
}

func (f *ultraFieldLogger) GetLevel() LogLevel {
	return f.logger.GetLevel()
}

func (f *ultraFieldLogger) SetShowCaller(show bool) {
	f.logger.SetShowCaller(show)
}

func (f *ultraFieldLogger) IsShowCaller() bool {
	return f.logger.IsShowCaller()
}

func (f *ultraFieldLogger) IsLevelEnabled(level LogLevel) bool {
	return f.logger.IsLevelEnabled(level)
}

// 添加缺失的日志条目方法
func (f *ultraFieldLogger) Log(level LogLevel, msg string) {
	f.logWithFieldsMsg(level, msg)
}

func (f *ultraFieldLogger) LogContext(ctx context.Context, level LogLevel, msg string) {
	f.logWithFieldsMsg(level, msg)
}

func (f *ultraFieldLogger) LogWithFields(level LogLevel, msg string, fields map[string]interface{}) {
	// 合并字段
	allFields := make(map[string]interface{})
	if f.fields != nil {
		for k, v := range f.fields {
			allFields[k] = v
		}
	}
	if f.key != "" {
		allFields[f.key] = f.value
	}
	for k, v := range fields {
		allFields[k] = v
	}
	
	f.logger.LogWithFields(level, msg, allFields)
}

// 委托其他方法
func (f *ultraFieldLogger) Print(v ...interface{}) { f.logger.Print(v...) }
func (f *ultraFieldLogger) Printf(format string, v ...interface{}) { f.logger.Printf(format, v...) }
func (f *ultraFieldLogger) Println(v ...interface{}) { f.logger.Println(v...) }

func (f *ultraFieldLogger) DebugContext(ctx context.Context, format string, args ...interface{}) {
	f.logger.DebugContext(ctx, format, args...)
}
func (f *ultraFieldLogger) InfoContext(ctx context.Context, format string, args ...interface{}) {
	f.logger.InfoContext(ctx, format, args...)
}
func (f *ultraFieldLogger) WarnContext(ctx context.Context, format string, args ...interface{}) {
	f.logger.WarnContext(ctx, format, args...)
}
func (f *ultraFieldLogger) ErrorContext(ctx context.Context, format string, args ...interface{}) {
	f.logger.ErrorContext(ctx, format, args...)
}
func (f *ultraFieldLogger) FatalContext(ctx context.Context, format string, args ...interface{}) {
	f.logger.FatalContext(ctx, format, args...)
}

func (f *ultraFieldLogger) DebugKV(msg string, keysAndValues ...interface{}) { f.logger.DebugKV(msg, keysAndValues...) }
func (f *ultraFieldLogger) InfoKV(msg string, keysAndValues ...interface{}) { f.logger.InfoKV(msg, keysAndValues...) }
func (f *ultraFieldLogger) WarnKV(msg string, keysAndValues ...interface{}) { f.logger.WarnKV(msg, keysAndValues...) }
func (f *ultraFieldLogger) ErrorKV(msg string, keysAndValues ...interface{}) { f.logger.ErrorKV(msg, keysAndValues...) }
func (f *ultraFieldLogger) FatalKV(msg string, keysAndValues ...interface{}) { f.logger.FatalKV(msg, keysAndValues...) }
func (f *ultraFieldLogger) LogKV(level LogLevel, msg string, keysAndValues ...interface{}) { f.logger.LogKV(level, msg, keysAndValues...) }

func (f *ultraFieldLogger) WithField(key string, value interface{}) ILogger {
	newFields := make(map[string]interface{})
	if f.fields != nil {
		for k, v := range f.fields {
			newFields[k] = v
		}
	}
	if f.key != "" {
		newFields[f.key] = f.value
	}
	newFields[key] = value
	return &ultraFieldLogger{logger: f.logger, fields: newFields}
}

func (f *ultraFieldLogger) WithFields(fields map[string]interface{}) ILogger {
	newFields := make(map[string]interface{})
	if f.fields != nil {
		for k, v := range f.fields {
			newFields[k] = v
		}
	}
	if f.key != "" {
		newFields[f.key] = f.value
	}
	for k, v := range fields {
		newFields[k] = v
	}
	return &ultraFieldLogger{logger: f.logger, fields: newFields}
}

func (f *ultraFieldLogger) WithError(err error) ILogger {
	return f.WithField("error", err.Error())
}

func (f *ultraFieldLogger) WithContext(ctx context.Context) ILogger {
	return f
}

func (f *ultraFieldLogger) Clone() ILogger {
	newFields := make(map[string]interface{})
	if f.fields != nil {
		for k, v := range f.fields {
			newFields[k] = v
		}
	}
	if f.key != "" {
		newFields[f.key] = f.value
	}
	return &ultraFieldLogger{logger: f.logger.Clone(), fields: newFields}
}