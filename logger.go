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
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/stringx"
)

// ============================================================================
// 性能优化：对象池和预计算常量
// ============================================================================

const (
	maxLogMessageSize    = 4096 // 单条日志消息的最大预分配大小（适应带 caller 的 JSON 日志）
	estimatedContextSize = 128  // 预估的上下文信息大小（TraceID 等）
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

	// 键值对日志的常量字符串
	kvSeparator  = []byte(": ")
	kvDelimiter  = []byte(", ")
	kvBraceOpen  = []byte(" {")
	kvBraceClose = []byte("}")
	kvMissing    = []byte("<missing>")
)

// 使用数组替代 map，O(1) 直接索引（避免哈希计算）
// 仅覆盖基础级别 0-4（DEBUG/INFO/WARN/ERROR/FATAL），其他级别前缀为 nil
var (
	levelPrefixesArr = [5][]byte{
		DEBUG: debugPrefix,
		INFO:  infoPrefix,
		WARN:  warnPrefix,
		ERROR: errorPrefix,
		FATAL: fatalPrefix,
	}

	levelPrefixesColorArr = [5][]byte{
		DEBUG: debugPrefixColor,
		INFO:  infoPrefixColor,
		WARN:  warnPrefixColor,
		ERROR: errorPrefixColor,
		FATAL: fatalPrefixColor,
	}
)

// getLevelPrefix 根据级别和颜色配置获取前缀
// 内联友好：避免 mathx.IF 的双 map 查找，用数组直接索引
func (l *Logger) getLevelPrefix(level LogLevel) []byte {
	idx := int(level)
	if idx < 0 || idx >= len(levelPrefixesArr) {
		return nil
	}
	if l.colorful {
		return levelPrefixesColorArr[idx]
	}
	return levelPrefixesArr[idx]
}

// levelNames 基础级别名称数组，用于 JSON 路径快速查找（避免 level.String() 的 map 查找）
var levelNames = [5]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// getLevelName 快速获取级别名称
func getLevelName(level LogLevel) string {
	idx := int(level)
	if idx >= 0 && idx < len(levelNames) {
		return levelNames[idx]
	}
	return level.String()
}

// ============================================================================
// 上下文提取器
// ============================================================================

// ContextExtractor 上下文信息提取器函数类型
type ContextExtractor func(ctx context.Context) string

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

// appendTextHeader 将时间戳、前缀、级别前缀、调用者信息追加到 buf
// 供带 fields/KV 的文本模式条目构建器共享，避免重复代码
// 注意：此函数成本较高无法内联，故 ultraLog 的纯文本路径保留内联实现以减少调用开销
func (l *Logger) appendTextHeader(buf []byte, level LogLevel) []byte {
	// 添加时间戳
	buf = time.Now().AppendFormat(buf, l.timeFormat)

	// 添加前缀（如果有）
	if l.prefix != "" {
		buf = append(buf, l.prefix...)
	}

	// 添加级别前缀
	buf = append(buf, l.getLevelPrefix(level)...)

	// 添加调用者信息（如果需要）
	if l.showCaller.Load() {
		if file, line, funcName := l.findExternalCaller(); file != "" {
			buf = append(buf, '[')
			buf = append(buf, file...)
			buf = append(buf, ':')
			buf = stringx.FastAppendInt(buf, line)
			buf = append(buf, ':')
			buf = append(buf, funcName...)
			buf = append(buf, ']', ' ')
		}
	}
	return buf
}

// writeTextBufLocked 将已构建好的文本条目缓冲区（含 header 与消息内容，不含换行）
// 追加换行后加锁写入并处理 FATAL 退出。所有文本模式日志方法共享此逻辑。
// 设计为可内联（成本远低于内联预算），避免调用开销。
func (l *Logger) writeTextBufLocked(level LogLevel, buf []byte) {
	buf = append(buf, newline...)
	l.mu.Lock()
	l.output.Write(buf)
	l.mu.Unlock()
	if level == FATAL {
		os.Exit(1)
	}
}

// ultraLog 极致优化的日志方法（使用字节池和零拷贝）
// 纯文本路径保留内联实现以避免 appendTextHeader 无法内联带来的额外调用开销
func (l *Logger) ultraLog(level LogLevel, msg string) {
	if level < LogLevel(l.level.Load()) {
		return
	}

	// JSON 格式走专用路径，输出结构化 JSON（无额外 fields）
	if l.format == FormatJSON {
		l.writeJSONEntry(level, msg, nil)
		return
	}

	// 文本路径：内联构建完整条目
	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	// 添加时间戳
	buf = time.Now().AppendFormat(buf, l.timeFormat)

	// 添加前缀（如果有）
	if l.prefix != "" {
		buf = append(buf, l.prefix...)
	}

	// 添加级别前缀
	buf = append(buf, l.getLevelPrefix(level)...)

	// 添加调用者信息（如果需要）
	if l.showCaller.Load() {
		if file, line, funcName := l.findExternalCaller(); file != "" {
			buf = append(buf, '[')
			buf = append(buf, file...)
			buf = append(buf, ':')
			buf = stringx.FastAppendInt(buf, line)
			buf = append(buf, ':')
			buf = append(buf, funcName...)
			buf = append(buf, ']', ' ')
		}
	}

	// 添加消息
	buf = append(buf, msg...)
	l.writeTextBufLocked(level, buf)
}

// writeJSONEntry 输出 JSON 格式日志条目
// traceId 等上下文字段、KV 字段都会作为 JSON 顶层字段输出，便于日志收集系统（如 openobserve）索引
func (l *Logger) writeJSONEntry(level LogLevel, msg string, fields map[string]any) {
	// 使用有序键值对构建，避免 map 遍历顺序不稳定
	// 同时复用 bytePool 减少分配
	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	buf = append(buf, '{')

	// timestamp
	buf = appendJSONKey(buf, l.timestampKey)
	buf = append(buf, '"')
	buf = time.Now().AppendFormat(buf, l.timeFormat)
	buf = append(buf, '"')

	// level
	buf = append(buf, ',')
	buf = appendJSONKey(buf, l.levelKey)
	buf = append(buf, '"')
	buf = append(buf, getLevelName(level)...)
	buf = append(buf, '"')

	// prefix（如果有）
	if l.prefix != "" {
		buf = append(buf, ',')
		buf = appendJSONKey(buf, "prefix")
		buf = appendJSONString(buf, strings.TrimSpace(l.prefix))
	}

	// message
	buf = append(buf, ',')
	buf = appendJSONKey(buf, l.messageKey)
	buf = appendJSONString(buf, msg)

	// caller（如果启用）- 用循环回溯找到第一个 go-logger 包外的调用者
	if l.showCaller.Load() {
		if file, line, funcName := l.findExternalCaller(); file != "" {
			buf = append(buf, ',')
			buf = appendJSONKey(buf, l.callerKey)
			buf = append(buf, '"')
			buf = appendJSONStringContent(buf, file)
			buf = append(buf, ':')
			buf = strconv.AppendInt(buf, int64(line), 10)
			buf = append(buf, '"')
			buf = append(buf, ',')
			buf = appendJSONKey(buf, "callerfunc")
			buf = appendJSONString(buf, funcName)
		}
	}

	// 额外 fields（traceId、KV 等）
	for k, v := range fields {
		buf = append(buf, ',')
		buf = appendJSONKey(buf, k)
		buf = appendJSONValue(buf, v)
	}

	buf = append(buf, '}', '\n')

	l.mu.Lock()
	l.output.Write(buf)
	l.mu.Unlock()

	if level == FATAL {
		os.Exit(1)
	}
}

// ultraLogWithFields 带额外 fields 的日志方法
// JSON 模式下 fields 作为 JSON 顶层字段输出；text 模式下 fields 拼接到 msg 后
// 在单个缓冲区中构建完整条目，避免 string(buf) 分配和二次缓冲
func (l *Logger) ultraLogWithFields(level LogLevel, msg string, fields map[string]any) {
	if level < LogLevel(l.level.Load()) {
		return
	}

	if l.format == FormatJSON {
		l.writeJSONEntry(level, msg, fields)
		return
	}

	// text 模式：无 fields 时直接走纯文本路径
	if len(fields) == 0 {
		l.ultraLog(level, msg)
		return
	}

	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	buf = l.appendTextHeader(buf, level)
	buf = append(buf, msg...)
	buf = append(buf, kvBraceOpen...)

	first := true
	for k, v := range fields {
		if !first {
			buf = append(buf, kvDelimiter...)
		}
		buf = append(buf, k...)
		buf = append(buf, kvSeparator...)
		buf = convert.AppendValue(buf, v)
		first = false
	}

	buf = append(buf, kvBraceClose...)
	l.writeTextBufLocked(level, buf)
}

// extractContextFields 从上下文提取信息为 map（用于 JSON 输出）
// 与 extractContextInfo 对应，但返回 map 而非拼接字符串
// 自定义 extractor 返回 string 无法拆分为 fields，降级放入 "context" 字段保留信息
func (l *Logger) extractContextFields(ctx context.Context) map[string]any {
	if ctx == nil {
		return nil
	}
	if l.contextExtractor != nil {
		info := l.contextExtractor(ctx)
		if info = strings.TrimSpace(info); info != "" {
			return map[string]any{"context": info}
		}
		return nil
	}
	return extractContextFieldsWithCompiledKeys(ctx, l.contextKeys)
}

// appendJSONKey 追加 JSON 键（假设 key 是简单 ASCII 字符串，无需转义）
func appendJSONKey(buf []byte, key string) []byte {
	buf = append(buf, '"')
	buf = append(buf, key...)
	buf = append(buf, '"', ':')
	return buf
}

// appendJSONStringContent 追加 JSON 字符串内容（不带引号，带转义处理）
func appendJSONStringContent(buf []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"':
			buf = append(buf, '\\', '"')
		case '\\':
			buf = append(buf, '\\', '\\')
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		case '\b':
			buf = append(buf, '\\', 'b')
		case '\f':
			buf = append(buf, '\\', 'f')
		default:
			if c < 0x20 {
				buf = append(buf, '\\', 'u', '0', '0',
					hexChars[c>>4], hexChars[c&0xf])
			} else {
				buf = append(buf, c)
			}
		}
	}
	return buf
}

// appendJSONString 追加 JSON 字符串值（带引号和转义处理）
func appendJSONString(buf []byte, s string) []byte {
	buf = append(buf, '"')
	buf = appendJSONStringContent(buf, s)
	buf = append(buf, '"')
	return buf
}

// appendJSONValue 追加 JSON 值（根据类型序列化）
func appendJSONValue(buf []byte, v any) []byte {
	switch val := v.(type) {
	case string:
		return appendJSONString(buf, val)
	case bool:
		if val {
			return append(buf, 't', 'r', 'u', 'e')
		}
		return append(buf, 'f', 'a', 'l', 's', 'e')
	case int:
		return strconv.AppendInt(buf, int64(val), 10)
	case int8:
		return strconv.AppendInt(buf, int64(val), 10)
	case int16:
		return strconv.AppendInt(buf, int64(val), 10)
	case int32:
		return strconv.AppendInt(buf, int64(val), 10)
	case int64:
		return strconv.AppendInt(buf, val, 10)
	case uint:
		return strconv.AppendUint(buf, uint64(val), 10)
	case uint8:
		return strconv.AppendUint(buf, uint64(val), 10)
	case uint16:
		return strconv.AppendUint(buf, uint64(val), 10)
	case uint32:
		return strconv.AppendUint(buf, uint64(val), 10)
	case uint64:
		return strconv.AppendUint(buf, val, 10)
	case float32:
		return strconv.AppendFloat(buf, float64(val), 'f', -1, 32)
	case float64:
		return strconv.AppendFloat(buf, val, 'f', -1, 64)
	case nil:
		return append(buf, 'n', 'u', 'l', 'l')
	default:
		// 其他类型降级用 fmt.Sprintf 转字符串后输出
		return appendJSONString(buf, fmt.Sprintf("%v", v))
	}
}

// hexChars 用于 JSON unicode 转义
var hexChars = []byte("0123456789abcdef")

// callerInfo 缓存的调用者信息，按 PC 索引
type callerInfo struct {
	file     string
	line     int
	funcName string
}

// callerCache 按调用者 PC 缓存调用者信息
// 同一调用点的 PC 固定不变，首次遍历栈后缓存，后续调用零分配命中缓存
var callerCache sync.Map

// findExternalCaller 回溯调用栈，找到第一个不在 go-logger/reflect/testing/runtime 内的调用者
// 优化策略：
//  1. 用 runtime.Callers 一次性获取所有 PC（替代循环 runtime.Caller）
//  2. 用文件路径判断内部调用者（避免 runtime.FuncForPC().Name() 的字符串分配）
//  3. 仅对外部调用者调用 .Name()（从 2-3 次分配降至 1 次）
//  4. 按 PC 缓存结果（热路径零分配）
func (l *Logger) findExternalCaller() (file string, line int, funcName string) {
	// 栈上数组，无需 Pool，无堆分配
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])

	for i := 0; i < n; i++ {
		pc := pcs[i]

		// 快速路径：缓存命中（sync.Map.Load 无分配）
		if cached, ok := callerCache.Load(pc); ok {
			ci := cached.(*callerInfo)
			return ci.file, ci.line, ci.funcName
		}

		// 未缓存：获取 *Func（无分配），用 FileLine 拿文件路径（无字符串分配）
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		f, ln := fn.FileLine(pc)

		// 用文件路径判断是否内部调用（无字符串分配）
		if isInternalCallerByFile(f) {
			continue
		}

		// 找到外部调用者：此处是唯一需要分配的地方（首次调用，后续命中缓存）
		name := fn.Name()
		// 统一路径分隔符为正斜杠，保留完整路径以便 IDE 点击跳转
		if strings.ContainsRune(f, '\\') {
			f = strings.ReplaceAll(f, `\`, `/`)
		}

		// 缓存结果（每个调用点只分配一次）
		ci := &callerInfo{file: f, line: ln, funcName: name}
		callerCache.Store(pc, ci)
		return f, ln, name
	}
	return "", 0, ""
}

// isInternalCallerByFile 根据文件路径判断是否为需要跳过的内部调用
// 仅使用文件路径（无需函数名），避免 runtime.FuncForPC().Name() 的字符串分配
// 跳过 go-logger 自身、reflect 反射、testing 框架、runtime 内部
// 注意：_test.go 文件即使属于 go-logger 包也不跳过（测试代码本身就是调用者）
func isInternalCallerByFile(file string) bool {
	if strings.HasSuffix(file, "_test.go") {
		return false
	}
	// 检查 go-logger 包（跨平台：同时检查正斜杠和反斜杠）
	if strings.Contains(file, `/go-logger/`) || strings.Contains(file, `\go-logger\`) {
		return true
	}
	// 检查 Go 标准库 runtime/reflect/testing
	if strings.Contains(file, `/src/runtime/`) || strings.Contains(file, `\src\runtime\`) {
		return true
	}
	if strings.Contains(file, `/src/reflect/`) || strings.Contains(file, `\src\reflect\`) {
		return true
	}
	if strings.Contains(file, `/src/testing/`) || strings.Contains(file, `\src\testing\`) {
		return true
	}
	return false
}

// logWithContextFormat 是 *Context 系列方法的公共实现
// 统一处理 level 检查、JSON/text 模式分支、traceId 提取
func (l *Logger) logWithContextFormat(ctx context.Context, level LogLevel, format string, args ...any) {
	if level < LogLevel(l.level.Load()) {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if l.format == FormatJSON {
		l.ultraLogWithFields(level, msg, l.extractContextFields(ctx))
		return
	}
	contextInfo := l.extractContextInfo(ctx)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	l.ultraLog(level, msg)
}

// ultraLogf 极致优化的格式化日志方法
func (l *Logger) ultraLogf(level LogLevel, format string, args ...any) {
	if level < LogLevel(l.level.Load()) {
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
	if LogLevel(l.level.Load()) > DEBUG {
		return
	}
	l.ultraLogf(DEBUG, format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...any) {
	if LogLevel(l.level.Load()) > INFO {
		return
	}
	l.ultraLogf(INFO, format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...any) {
	if LogLevel(l.level.Load()) > WARN {
		return
	}
	l.ultraLogf(WARN, format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...any) {
	if LogLevel(l.level.Load()) > ERROR {
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
	if LogLevel(l.level.Load()) > DEBUG {
		return
	}
	l.ultraLogf(DEBUG, format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	if LogLevel(l.level.Load()) > INFO {
		return
	}
	l.ultraLogf(INFO, format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	if LogLevel(l.level.Load()) > WARN {
		return
	}
	l.ultraLogf(WARN, format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	if LogLevel(l.level.Load()) > ERROR {
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
	if LogLevel(l.level.Load()) > DEBUG {
		return
	}
	l.ultraLog(DEBUG, msg)
}

func (l *Logger) InfoMsg(msg string) {
	if LogLevel(l.level.Load()) > INFO {
		return
	}
	l.ultraLog(INFO, msg)
}

func (l *Logger) WarnMsg(msg string) {
	if LogLevel(l.level.Load()) > WARN {
		return
	}
	l.ultraLog(WARN, msg)
}

func (l *Logger) ErrorMsg(msg string) {
	if LogLevel(l.level.Load()) > ERROR {
		return
	}
	l.ultraLog(ERROR, msg)
}

func (l *Logger) FatalMsg(msg string) {
	l.ultraLog(FATAL, msg)
}

// 多行日志方法 - 自动处理换行符
func (l *Logger) InfoLines(lines ...string) {
	if LogLevel(l.level.Load()) > INFO {
		return
	}
	for _, line := range lines {
		l.ultraLog(INFO, line)
	}
}

func (l *Logger) ErrorLines(lines ...string) {
	if LogLevel(l.level.Load()) > ERROR {
		return
	}
	for _, line := range lines {
		l.ultraLog(ERROR, line)
	}
}

func (l *Logger) WarnLines(lines ...string) {
	if LogLevel(l.level.Load()) > WARN {
		return
	}
	for _, line := range lines {
		l.ultraLog(WARN, line)
	}
}

func (l *Logger) DebugLines(lines ...string) {
	if LogLevel(l.level.Load()) > DEBUG {
		return
	}
	for _, line := range lines {
		l.ultraLog(DEBUG, line)
	}
}

// logWithContextLines 多行上下文日志内部实现 - 每一行均附加 ctx 信息
func (l *Logger) logWithContextLines(ctx context.Context, level LogLevel, lines ...string) {
	if level < LogLevel(l.level.Load()) {
		return
	}
	if l.format == FormatJSON {
		fields := l.extractContextFields(ctx)
		for _, line := range lines {
			l.ultraLogWithFields(level, line, fields)
		}
		return
	}
	contextInfo := l.extractContextInfo(ctx)
	for _, line := range lines {
		msg := line
		if contextInfo != "" {
			msg = contextInfo + line
		}
		l.ultraLog(level, msg)
	}
}

// 多行日志方法（带上下文）
func (l *Logger) DebugContextLines(ctx context.Context, lines ...string) {
	l.logWithContextLines(ctx, DEBUG, lines...)
}

func (l *Logger) InfoContextLines(ctx context.Context, lines ...string) {
	l.logWithContextLines(ctx, INFO, lines...)
}

func (l *Logger) WarnContextLines(ctx context.Context, lines ...string) {
	l.logWithContextLines(ctx, WARN, lines...)
}

func (l *Logger) ErrorContextLines(ctx context.Context, lines ...string) {
	l.logWithContextLines(ctx, ERROR, lines...)
}

// SetContextExtractor 设置自定义上下文提取器
func (l *Logger) SetContextExtractor(extractor ContextExtractor) {
	l.contextExtractor = extractor
}

// GetContextExtractor 获取当前的上下文提取器
func (l *Logger) GetContextExtractor() ContextExtractor {
	return l.contextExtractor
}

// extractContextInfo 从上下文中提取信息
func (l *Logger) extractContextInfo(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if l.contextExtractor != nil {
		return l.contextExtractor(ctx)
	}
	return extractContextWithCompiledKeys(ctx, l.contextKeys)
}

// 带上下文的日志方法
func (l *Logger) DebugContext(ctx context.Context, format string, args ...any) {
	l.logWithContextFormat(ctx, DEBUG, format, args...)
}

func (l *Logger) InfoContext(ctx context.Context, format string, args ...any) {
	l.logWithContextFormat(ctx, INFO, format, args...)
}

func (l *Logger) WarnContext(ctx context.Context, format string, args ...any) {
	l.logWithContextFormat(ctx, WARN, format, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, format string, args ...any) {
	l.logWithContextFormat(ctx, ERROR, format, args...)
}

func (l *Logger) FatalContext(ctx context.Context, format string, args ...any) {
	l.logWithContextFormat(ctx, FATAL, format, args...)
}

// ============================================================================
// 键值对和结构化日志辅助方法
// ============================================================================

// logWithKV 极简键值对实现 - 零分配优化
func (l *Logger) logWithKV(level LogLevel, msg string, keysAndValues ...any) {
	if level < LogLevel(l.level.Load()) {
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

	// JSON 模式：把 KV 转为 fields map，作为 JSON 顶层字段输出
	if l.format == FormatJSON {
		fields := kvToFields(keysAndValues)
		l.ultraLogWithFields(level, msg, fields)
		return
	}

	// text 模式：在单个缓冲区中构建完整条目，避免 string(buf) 分配和二次缓冲
	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	buf = l.appendTextHeader(buf, level)
	buf = append(buf, msg...)
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
	l.writeTextBufLocked(level, buf)
}

// kvToFields 将键值对切片转为 map[string]any
// 用于 JSON 模式下将 KV 作为 JSON 顶层字段输出
func kvToFields(keysAndValues []any) map[string]any {
	if len(keysAndValues) == 0 {
		return nil
	}
	fields := make(map[string]any, len(keysAndValues)/2+1)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keysAndValues[i])
		}
		if i+1 < len(keysAndValues) {
			fields[key] = keysAndValues[i+1]
		} else {
			fields[key] = "<missing>"
		}
	}
	return fields
}

// logWithFields 使用字段映射记录日志
func (l *Logger) logWithFields(level LogLevel, msg string, fields map[string]any) {
	if level < LogLevel(l.level.Load()) {
		return
	}

	// JSON 模式：fields 作为 JSON 顶层字段输出
	if l.format == FormatJSON {
		l.ultraLogWithFields(level, msg, fields)
		return
	}

	if len(fields) == 0 {
		l.ultraLog(level, msg)
		return
	}

	// text 模式：在单个缓冲区中构建完整条目，避免 string(buf) 分配和二次缓冲
	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	defer bytePool.Put(buf)

	buf = l.appendTextHeader(buf, level)
	buf = append(buf, msg...)
	buf = append(buf, kvBraceOpen...)

	first := true
	for k, v := range fields {
		if !first {
			buf = append(buf, kvDelimiter...)
		}
		buf = append(buf, k...)
		buf = append(buf, kvSeparator...)
		buf = convert.AppendValue(buf, v)
		first = false
	}

	buf = append(buf, kvBraceClose...)
	l.writeTextBufLocked(level, buf)
}

// logWithContextKV 带上下文的键值对日志
func (l *Logger) logWithContextKV(ctx context.Context, level LogLevel, msg string, keysAndValues ...any) {
	if level < LogLevel(l.level.Load()) {
		return
	}

	// JSON 模式：traceId 和 KV 都作为 JSON 顶层字段输出
	if l.format == FormatJSON {
		fields := l.extractContextFields(ctx)
		if len(keysAndValues) > 0 {
			// 合并 KV 到 fields（KV 优先级高于 context）
			if fields == nil {
				fields = kvToFields(keysAndValues)
			} else {
				kvFields := kvToFields(keysAndValues)
				for k, v := range kvFields {
					fields[k] = v
				}
			}
		}
		l.ultraLogWithFields(level, msg, fields)
		return
	}

	// text 模式：先从 context 提取信息 prepend 到 msg
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
	if LogLevel(l.level.Load()) > DEBUG {
		return
	}
	l.logWithKV(DEBUG, msg, keysAndValues...)
}

func (l *Logger) DebugContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if LogLevel(l.level.Load()) > DEBUG {
		return
	}
	l.logWithContextKV(ctx, DEBUG, msg, keysAndValues...)
}

func (l *Logger) InfoKV(msg string, keysAndValues ...any) {
	if LogLevel(l.level.Load()) > INFO {
		return
	}
	l.logWithKV(INFO, msg, keysAndValues...)
}

func (l *Logger) InfoContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if LogLevel(l.level.Load()) > INFO {
		return
	}
	l.logWithContextKV(ctx, INFO, msg, keysAndValues...)
}

func (l *Logger) WarnKV(msg string, keysAndValues ...any) {
	if LogLevel(l.level.Load()) > WARN {
		return
	}
	l.logWithKV(WARN, msg, keysAndValues...)
}

func (l *Logger) WarnContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if LogLevel(l.level.Load()) > WARN {
		return
	}
	l.logWithContextKV(ctx, WARN, msg, keysAndValues...)
}

func (l *Logger) ErrorKV(msg string, keysAndValues ...any) {
	if LogLevel(l.level.Load()) > ERROR {
		return
	}
	l.logWithKV(ERROR, msg, keysAndValues...)
}

func (l *Logger) ErrorContextKV(ctx context.Context, msg string, keysAndValues ...any) {
	if LogLevel(l.level.Load()) > ERROR {
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
	if LogLevel(l.level.Load()) > DEBUG {
		return
	}
	l.logWithFields(DEBUG, msg, fields)
}

func (l *Logger) InfoWithFields(msg string, fields map[string]any) {
	if LogLevel(l.level.Load()) > INFO {
		return
	}
	l.logWithFields(INFO, msg, fields)
}

func (l *Logger) WarnWithFields(msg string, fields map[string]any) {
	if LogLevel(l.level.Load()) > WARN {
		return
	}
	l.logWithFields(WARN, msg, fields)
}

func (l *Logger) ErrorWithFields(msg string, fields map[string]any) {
	if LogLevel(l.level.Load()) > ERROR {
		return
	}
	l.logWithFields(ERROR, msg, fields)
}

func (l *Logger) FatalWithFields(msg string, fields map[string]any) {
	l.logWithFields(FATAL, msg, fields)
}

// 原始日志条目方法
func (l *Logger) Log(level LogLevel, msg string) {
	if level < LogLevel(l.level.Load()) {
		return
	}
	l.ultraLog(level, msg)
}

func (l *Logger) LogContext(ctx context.Context, level LogLevel, msg string) {
	if level < LogLevel(l.level.Load()) {
		return
	}
	// JSON 模式：traceId 作为 JSON 顶层字段
	if l.format == FormatJSON {
		l.ultraLogWithFields(level, msg, l.extractContextFields(ctx))
		return
	}
	contextInfo := l.extractContextInfo(ctx)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	l.ultraLog(level, msg)
}

func (l *Logger) LogKV(level LogLevel, msg string, keysAndValues ...any) {
	if level < LogLevel(l.level.Load()) {
		return
	}
	l.logWithKV(level, msg, keysAndValues...)
}

func (l *Logger) LogWithFields(level LogLevel, msg string, fields map[string]any) {
	if level < LogLevel(l.level.Load()) {
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
			output:          l.output,
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
	l.level.Store(int32(level))
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() LogLevel {
	return LogLevel(l.level.Load())
}

// SetShowCaller 设置是否显示调用者信息
func (l *Logger) SetShowCaller(show bool) {
	l.showCaller.Store(show)
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
	if LogLevel(f.logger.level.Load()) > DEBUG {
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

// logWithContextFieldsLines 是 fieldLogger.*ContextLines 系列方法的公共实现
// 统一处理 level 检查、JSON/text 模式分支、ctx 信息 + f.fields 合并
func (f *fieldLogger) logWithContextFieldsLines(ctx context.Context, level LogLevel, lines ...string) {
	if level < LogLevel(f.logger.level.Load()) {
		return
	}
	if f.logger.format == FormatJSON {
		fields := f.mergeContextFields(ctx)
		for _, line := range lines {
			f.logger.ultraLogWithFields(level, line, fields)
		}
		return
	}
	contextInfo := f.logger.extractContextInfo(ctx)
	for _, line := range lines {
		msg := line
		if contextInfo != "" {
			msg = contextInfo + line
		}
		f.logger.logWithFields(level, msg, f.fields)
	}
}

// 多行日志方法（带上下文）
func (f *fieldLogger) DebugContextLines(ctx context.Context, lines ...string) {
	f.logWithContextFieldsLines(ctx, DEBUG, lines...)
}

func (f *fieldLogger) InfoContextLines(ctx context.Context, lines ...string) {
	f.logWithContextFieldsLines(ctx, INFO, lines...)
}

func (f *fieldLogger) WarnContextLines(ctx context.Context, lines ...string) {
	f.logWithContextFieldsLines(ctx, WARN, lines...)
}

func (f *fieldLogger) ErrorContextLines(ctx context.Context, lines ...string) {
	f.logWithContextFieldsLines(ctx, ERROR, lines...)
}

// mergeContextFields 合并 context 提取的 fields 和 f.fields（用于 JSON 模式）
// traceId 等上下文字段与 f.fields 合并，f.fields 优先级更高（避免被覆盖）
func (f *fieldLogger) mergeContextFields(ctx context.Context) map[string]any {
	ctxFields := f.logger.extractContextFields(ctx)
	if len(f.fields) == 0 {
		return ctxFields
	}
	if ctxFields == nil {
		return f.fields
	}
	merged := make(map[string]any, len(ctxFields)+len(f.fields))
	for k, v := range ctxFields {
		merged[k] = v
	}
	for k, v := range f.fields {
		merged[k] = v
	}
	return merged
}

// logWithContextFieldsFormat 是 fieldLogger.*Context 系列方法的公共实现
// 统一处理 level 检查、JSON/text 模式分支、traceId + f.fields 合并
func (f *fieldLogger) logWithContextFieldsFormat(ctx context.Context, level LogLevel, format string, args ...any) {
	if level < LogLevel(f.logger.level.Load()) {
		return
	}
	msg := fmt.Sprintf(format, args...)
	if f.logger.format == FormatJSON {
		f.logger.ultraLogWithFields(level, msg, f.mergeContextFields(ctx))
		return
	}
	contextInfo := f.logger.extractContextInfo(ctx)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	f.logger.logWithFields(level, msg, f.fields)
}

// 上下文日志方法
func (f *fieldLogger) DebugContext(ctx context.Context, format string, args ...any) {
	f.logWithContextFieldsFormat(ctx, DEBUG, format, args...)
}

func (f *fieldLogger) InfoContext(ctx context.Context, format string, args ...any) {
	f.logWithContextFieldsFormat(ctx, INFO, format, args...)
}

func (f *fieldLogger) WarnContext(ctx context.Context, format string, args ...any) {
	f.logWithContextFieldsFormat(ctx, WARN, format, args...)
}

func (f *fieldLogger) ErrorContext(ctx context.Context, format string, args ...any) {
	f.logWithContextFieldsFormat(ctx, ERROR, format, args...)
}

func (f *fieldLogger) FatalContext(ctx context.Context, format string, args ...any) {
	f.logWithContextFieldsFormat(ctx, FATAL, format, args...)
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
	// JSON 模式：traceId 和 f.fields 都作为 JSON 顶层字段
	if f.logger.format == FormatJSON {
		f.logger.ultraLogWithFields(level, msg, f.mergeContextFields(ctx))
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

// LogSpecialContext 委托给底层 Logger，并合并 fieldLogger 的字段
func (f *fieldLogger) LogSpecialContext(ctx context.Context, logType SpecialLogType, level LogLevel, format string, args ...any) {
	if !f.logger.IsLevelEnabled(level) {
		return
	}
	// 构建消息：emoji [NAME] format
	msg := logType.emoji + " [" + logType.name + "] "
	if len(args) == 0 {
		msg += format
	} else {
		msg += fmt.Sprintf(format, args...)
	}
	// JSON 模式：traceId 和 f.fields 都作为 JSON 顶层字段
	if f.logger.format == FormatJSON {
		f.logger.ultraLogWithFields(level, msg, f.mergeContextFields(ctx))
		return
	}
	// text 模式：ctx 信息前缀
	contextInfo := f.logger.extractContextInfo(ctx)
	if contextInfo != "" {
		msg = contextInfo + msg
	}
	f.logger.logWithFields(level, msg, f.fields)
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

// 结构化日志构建器
func (f *fieldLogger) WithField(key string, value any) ILogger {
	// 按实际所需容量分配 map，避免过度预分配
	newFields := make(map[string]any, len(f.fields)+1)

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

	// 按实际所需容量分配 map
	newFields := make(map[string]any, len(f.fields)+len(fields))

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

// 辅助方法：合并字段映射
func (f *fieldLogger) mergeFieldsMap(fields map[string]any) map[string]any {
	if len(fields) == 0 {
		return f.fields
	}

	// 按实际所需容量分配 map
	merged := make(map[string]any, len(f.fields)+len(fields))

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

// buildSpecialMessage 构建特殊日志消息内容：emoji [NAME] format
// 返回池化 buffer，调用方负责归还 bytePool
func buildSpecialMessage(emoji, name, format string, args ...any) []byte {
	buf := bytePool.Get().([]byte)
	buf = buf[:0]
	buf = append(buf, emoji...)
	buf = append(buf, ' ', '[')
	buf = append(buf, name...)
	buf = append(buf, ']', ' ')
	if len(args) == 0 {
		buf = append(buf, format...)
	} else {
		buf = append(buf, fmt.Sprintf(format, args...)...)
	}
	return buf
}

// logSpecialInternal 特殊日志内部实现（无 context）
// 所有特殊日志方法（logSpecial/Performance/Progress/Milestone/Health/Audit）的公共路径
func (l *Logger) logSpecialInternal(level LogLevel, emoji, name, format string, args ...any) {
	if level < LogLevel(l.level.Load()) {
		return
	}

	msgBuf := buildSpecialMessage(emoji, name, format, args...)
	defer bytePool.Put(msgBuf)

	if l.format == FormatJSON {
		l.writeJSONEntry(level, string(msgBuf), nil)
		return
	}

	buf := bytePool.Get().([]byte)
	defer bytePool.Put(buf)
	buf = buf[:0]
	buf = l.appendTextHeader(buf, level)
	buf = append(buf, msgBuf...)
	l.writeTextBufLocked(level, buf)
}

// logSpecial 记录特殊类型的日志（无 context）
func (l *Logger) logSpecial(logType SpecialLogType, level LogLevel, format string, args ...any) {
	l.logSpecialInternal(level, logType.emoji, logType.name, format, args...)
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
	emoji, level := getPerformanceLevel(duration)
	format := "%s completed in %v"
	args := []any{operation, duration}
	if len(details) > 0 && len(details[0]) > 0 {
		format += " | Details: %+v"
		args = append(args, details[0])
	}
	l.logSpecialInternal(PERFORMANCE, emoji, "PERF-"+level, format, args...)
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
	percentage := float64(current) / float64(total) * 100
	emoji := getProgressEmoji(percentage)
	l.logSpecialInternal(INFO, emoji, "PROGRESS", "%s: %d/%d (%.1f%%)", operation, current, total, percentage)
}

// Milestone 里程碑日志（INFO 级别）
func (l *Logger) Milestone(message string) {
	l.logSpecialInternal(INFO, "🎯", "MILESTONE", "%s", message)
}

// Health 健康检查日志（WARN 级别用于不健康，INFO 级别用于健康）
func (l *Logger) Health(service string, status bool, details string) {
	level := WARN
	emoji := "❌"
	statusStr := "UNHEALTHY"
	if status {
		level = INFO
		emoji = "✅"
		statusStr = "HEALTHY"
	}
	format := "%s: %s"
	args := []any{service, statusStr}
	if details != "" {
		format += " | %s"
		args = append(args, details)
	}
	l.logSpecialInternal(level, emoji, "HEALTH", format, args...)
}

// Audit 审计日志（AUDIT 级别）
func (l *Logger) Audit(action, user, resource, result string) {
	l.logSpecialInternal(AUDIT, "📋", "AUDIT", "User: %s | Action: %s | Resource: %s | Result: %s", user, action, resource, result)
}

// LogSpecialContext 记录带上下文的特殊类型日志（公开方法）
//
// 使用示例：
//
//	logger.LogSpecialContext(ctx, SuccessType, INFO, "操作成功")
//	logger.LogSpecialContext(ctx, DatabaseType, INFO, "查询 %s 耗时 %v", table, duration)
//	logger.LogSpecialContext(ctx, SecurityType, SECURITY, "检测到异常访问 %s", ip)
//
// JSON 模式：ctx 字段提取到 JSON fields；text 模式：ctx 信息前缀到消息
func (l *Logger) LogSpecialContext(ctx context.Context, logType SpecialLogType, level LogLevel, format string, args ...any) {
	if level < LogLevel(l.level.Load()) {
		return
	}

	msgBuf := buildSpecialMessage(logType.emoji, logType.name, format, args...)
	defer bytePool.Put(msgBuf)

	if l.format == FormatJSON {
		l.writeJSONEntry(level, string(msgBuf), l.extractContextFields(ctx))
		return
	}

	// text 模式：ctx 信息前缀到消息
	buf := bytePool.Get().([]byte)
	defer bytePool.Put(buf)
	buf = buf[:0]
	buf = l.appendTextHeader(buf, level)
	if contextInfo := l.extractContextInfo(ctx); contextInfo != "" {
		buf = append(buf, contextInfo...)
	}
	buf = append(buf, msgBuf...)
	l.writeTextBufLocked(level, buf)
}
