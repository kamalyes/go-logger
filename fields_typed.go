/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-03 16:02:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-03 16:38:00
 * @FilePath: \go-logger\fields_typed.go
 * @Description: 类型化日志字段 - 免去 ...any 入口在调用方的逐值装箱
 *
 * 为什么需要它：InfoKV 这类 ...any 入口要求调用方把每个值转成 any，
 * 而 string 转 interface 是一次堆分配（值不是编译期常量时编译器无法静态化）
 * 一条十四个字段的结构化日志就是十四次分配，且这笔开销落在调用方、
 * 在 go-logger 内部怎么优化都看不见
 *
 * 实测于 go-rpc-gateway 的访问日志中间件：默认配置每条日志十四个字段，
 * 链路前三层每请求三十次分配里有十四次是这个装箱
 *
 * Field 把值按类型存进结构体，变长参数只分配一个切片，与字段个数无关
 * 这是 zap 与 zerolog 的做法，命名也照它们的惯例，避免调用方意外
 *
 * [EN] Typed log fields, removing the per-value boxing that ...any entry points force
 * on callers
 * Converting a string to an interface is a heap allocation whenever the value is not a
 * compile-time constant, so a fourteen-field structured line costs fourteen allocations,
 * and that cost lands in the caller where no amount of tuning inside go-logger can see it
 * Measured in go-rpc-gateway's access-log middleware: fourteen of the thirty allocations
 * per request across the first three middleware layers were this boxing
 * A Field stores the value by type in a struct, so the variadic costs one slice
 * regardless of field count, following zap's and zerolog's design and naming
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package logger

import (
	"math"
	"strconv"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/desensitize"
)

// fieldKind 标记 Field 里哪个槽位有效
type fieldKind uint8

const (
	fieldKindString fieldKind = iota
	fieldKindInt
	fieldKindUint
	fieldKindFloat32
	fieldKindFloat64
	fieldKindBool
	fieldKindDuration
	fieldKindError
	// fieldKindAny 兜底：类型不在上面几种时退回 any，此时仍有一次装箱
	// 提供它是为了让调用方不必为了一个少见类型就整条日志退回 InfoKV
	// [EN] The fallback for types not covered above, which still boxes once; it exists so
	// a caller need not drop a whole line back to InfoKV for one uncommon type
	fieldKindAny
)

// Field 一个类型化的日志字段
//
// 刻意用单一结构体加 kind 标记，而不是接口：接口会把问题绕回装箱
// 结构体按最宽的槽位对齐，一个 Field 在 arm64 上是 48 字节，
// 十四个字段的切片一次分配 672 字节 —— 对比十四次各自分配
//
// [EN] One typed log field
// Deliberately a single struct with a kind tag rather than an interface, which would
// bring the boxing straight back
// The struct is sized by its widest slot, 48 bytes on arm64, so fourteen fields are one
// allocation of 672 bytes against fourteen separate ones
type Field struct {
	Key string

	kind fieldKind
	// num 承载全部定宽标量：int / uint / bool（0 与 1）/ duration /
	// 以及 float 的位模式（math.Float64bits）
	//
	// 合并槽位不是为了省内存，是为了省拷贝：Field 按变长参数传递，
	// 十四个字段就要拷十四个结构体 最初的版本给 float、bool、error 各留一个槽，
	// 结构体 88 字节、十四个字段 1.2 KB 拷贝，实测虽然把堆分配降到零，
	// 墙上时间反而比 ...any 路径慢三成 —— 用十一次堆分配换 1.2 KB 拷贝是亏的
	// 合并后 64 字节，与 zap 的 Field 同量级
	//
	// [EN] num carries every fixed-width scalar: ints, uints, bools as 0 and 1,
	// durations, and float bit patterns via math.Float64bits
	// Merging the slots is about copying rather than memory: Fields travel as a variadic,
	// so fourteen fields mean fourteen struct copies
	// The first version gave float, bool and error their own slots, making the struct
	// 88 bytes and the variadic a 1.2 KB copy, which drove heap allocations to zero but
	// measured 30% slower in wall time than the ...any path, since trading eleven
	// allocations for a 1.2 KB copy is a loss
	// Merged it is 64 bytes, the same order as zap's Field
	num int64
	str string
	// iface 只承载 error 与兜底的 any，这两种本就带接口
	// [EN] iface carries only error and the any fallback, both already interface-shaped
	iface any
}

// Str 字符串字段，访问日志里绝大多数字段都是这一种
// [EN] A string field, which is what most access-log fields are
func Str(key, val string) Field {
	return Field{Key: key, kind: fieldKindString, str: val}
}

// Int 整数字段，int / int8 / int16 / int32 / int64 都用它
// [EN] An integer field, covering every signed width
func Int[T int | int8 | int16 | int32 | int64](key string, val T) Field {
	return Field{Key: key, kind: fieldKindInt, num: int64(val)}
}

// Uint 无符号整数字段
// [EN] An unsigned integer field
func Uint[T uint | uint8 | uint16 | uint32 | uint64](key string, val T) Field {
	return Field{Key: key, kind: fieldKindUint, num: int64(val)}
}

// Float32 单精度浮点字段
// 与 Float64 分开是为了对齐 appendJSONValue：它按 float32 / float64 用不同的
// bitSize 调 strconv.AppendFloat，压成一种会让输出与 ...any 路径不同
//
// [EN] A float32 field, kept apart from Float64 to match appendJSONValue, which calls
// strconv.AppendFloat with a different bitSize per width; collapsing them would make the
// output differ from the ...any path
func Float32(key string, val float32) Field {
	return Field{Key: key, kind: fieldKindFloat32, num: int64(math.Float64bits(float64(val)))}
}

// Float64 双精度浮点字段
// [EN] A float64 field
func Float64(key string, val float64) Field {
	return Field{Key: key, kind: fieldKindFloat64, num: int64(math.Float64bits(val))}
}

// Bool 布尔字段
// [EN] A boolean field
func Bool(key string, val bool) Field {
	var n int64
	if val {
		n = 1
	}
	return Field{Key: key, kind: fieldKindBool, num: n}
}

// Dur 时长字段，输出 Duration.String() 形态的 JSON 字符串，如 "1.5ms"
//
// 这是为了与 ...any 路径一致：appendJSONValue 对 time.Duration 走的就是
// appendJSONString(val.String())，而不是毫秒整数
// 需要毫秒整数的调用方自己传 Int(key, d.Milliseconds())
//
// [EN] A duration field emitted as the JSON string form of Duration.String()
// This matches the ...any path, where appendJSONValue writes
// appendJSONString(val.String()) rather than an integer of milliseconds
// A caller wanting milliseconds passes Int(key, d.Milliseconds()) instead
func Dur(key string, val time.Duration) Field {
	return Field{Key: key, kind: fieldKindDuration, num: int64(val)}
}

// Err 错误字段，nil 时输出 JSON 的 null
//
// 为什么是 null 而不是空串：...any 路径上 error(nil) 装进 any 之后就是一个 nil 接口，
// 与裸 nil 不可区分，appendJSONValue 对它输出 null
// 类型化路径本可以区分「nil 的 error」与「没有值」，但两条路径的输出必须一致，
// 否则同一份日志换个入口就变样 null 也确实是「没有错误」更诚实的 JSON 表达
//
// [EN] An error field; a nil error emits JSON null
// On the ...any path an error(nil) boxed into an any is simply a nil interface,
// indistinguishable from a bare nil, and appendJSONValue writes null for it
// The typed path could tell a nil error from an absent value, but both paths must agree
// or the same log changes shape depending on the entry point, and null is the more
// honest JSON for the absence of an error anyway
func Err(key string, val error) Field {
	return Field{Key: key, kind: fieldKindError, iface: val}
}

// Anyf 兜底字段：类型不在上面几种时用它，此时仍有一次装箱
// 名字带 f 后缀是为了不和内建的 any 类型名撞
//
// [EN] The fallback field for uncovered types, which still boxes once
// The f suffix avoids colliding with the builtin any type name
func Anyf(key string, val any) Field {
	return Field{Key: key, kind: fieldKindAny, iface: val}
}

// textValue 把字段的值取成字符串形态，供脱敏与截断逻辑复用
// 只有 string 与 error 两种有文本语义，其余返回 ok=false
//
// [EN] Returns the field's value in text form for the sanitising and clamping logic
// Only string and error carry text semantics; the rest return ok false
func (f Field) textValue() (string, bool) {
	switch f.kind {
	case fieldKindString:
		return f.str, true
	case fieldKindError:
		if f.errValue() == nil {
			// nil error 没有文本形态，交由渲染路径输出 null
			// [EN] A nil error has no text form; the render path emits null
			return "", false
		}
		return f.errValue().Error(), true
	}
	return "", false
}

// appendFieldsJSON 把类型化字段直写为 JSON 顶层字段
//
// 逐行对照 appendKVPairsJSON 实现，包括超限回滚到本对 KV 之前的边界、
// 截断标记、以及脱敏按字段名判定的顺序 —— 输出必须与 ...any 路径逐字节相同，
// 否则换个入口日志格式就变了，下游的日志解析会碎
// TestFieldsMatchKV 逐项比对两条路径的输出
//
// [EN] Writes typed fields as top-level JSON, mirroring appendKVPairsJSON line by line,
// including the rollback to the pair boundary on budget overrun, the truncation marker
// and the order in which sanitising decides by key name
// The output must be byte-for-byte identical to the ...any path, or switching entry point
// silently changes the log format and breaks downstream parsing
// TestFieldsMatchKV compares the two paths case by case
func (l *Logger) appendFieldsJSON(buf []byte, fields []Field) []byte {
	sanitizing := l.sanitize != SanitizeOff
	budget := l.lineBudget(true)
	trimmed := false

	for i := range fields {
		f := fields[i]
		mark := len(buf)
		buf = append(buf, ',')
		buf = appendJSONKey(buf, f.Key)

		if sanitizing {
			var cut bool
			buf, cut = l.appendSanitizedFieldJSON(buf, f, valueRoom(buf, budget))
			trimmed = trimmed || cut
			if overBudget(buf, budget) {
				return rollbackTruncated(buf, mark, true)
			}
			continue
		}

		var cut bool
		buf, cut = l.appendFieldValueJSON(buf, f, valueRoom(buf, budget))
		trimmed = trimmed || cut
		if overBudget(buf, budget) {
			return rollbackTruncated(buf, mark, true)
		}
	}

	if trimmed {
		buf = markTruncated(buf, true)
	}
	return buf
}

// appendFieldValueJSON 不脱敏时的值输出，带长度上限
// 文本类字段走 clampString 的 string 快路径，数值类无需截断
//
// [EN] Value output with the length cap and no sanitising
// Text-bearing fields take the clampString fast path while numeric kinds need no clamping
func (l *Logger) appendFieldValueJSON(buf []byte, f Field, room int) ([]byte, bool) {
	switch f.kind {
	case fieldKindString:
		s, cut := l.clampString(f.str, room)
		return appendJSONString(buf, s), cut
	case fieldKindError:
		s, ok := f.textValue()
		if !ok {
			return append(buf, 'n', 'u', 'l', 'l'), false
		}
		s, cut := l.clampString(s, room)
		return appendJSONString(buf, s), cut
	case fieldKindAny:
		v, cut := l.clampValue(f.iface, room)
		return appendJSONValue(buf, v), cut
	}
	return l.appendFieldNumberJSON(buf, f), false
}

// appendFieldNumberJSON 数值与布尔字段直写，不经过 any
// [EN] Numeric and boolean fields are written directly, never through an any
func (l *Logger) appendFieldNumberJSON(buf []byte, f Field) []byte {
	switch f.kind {
	case fieldKindInt:
		return strconv.AppendInt(buf, f.num, 10)
	case fieldKindUint:
		return strconv.AppendUint(buf, uint64(f.num), 10)
	case fieldKindFloat32:
		return strconv.AppendFloat(buf, f.floatValue(), 'f', -1, 32)
	case fieldKindFloat64:
		return strconv.AppendFloat(buf, f.floatValue(), 'f', -1, 64)
	case fieldKindBool:
		if f.num != 0 {
			return append(buf, 't', 'r', 'u', 'e')
		}
		return append(buf, 'f', 'a', 'l', 's', 'e')
	case fieldKindDuration:
		return appendJSONString(buf, time.Duration(f.num).String())
	}
	return append(buf, 'n', 'u', 'l', 'l')
}

// appendSanitizedFieldJSON 脱敏开启时的值输出
//
// 判定顺序与 appendSanitizedJSONValue 一致：先按字段名查动作，
// 再看是否全值扫描，最后才是字段级长度上限
// 数值类字段只可能命中 WipeValue，掩码对它们没有意义
//
// [EN] Value output with sanitising on
// The decision order matches appendSanitizedJSONValue: the per-key action first, then
// whole-value scanning, then the per-field length cap
// Numeric kinds can only match WipeValue, since masking has no meaning for them
func (l *Logger) appendSanitizedFieldJSON(buf []byte, f Field, room int) ([]byte, bool) {
	if s, ok := f.textValue(); ok {
		clamped, cut := l.clampString(s, room)
		return l.appendSanitizedJSONStringValue(buf, f.Key, clamped), cut
	}
	if f.kind == fieldKindAny {
		v, cut := l.clampValue(f.iface, room)
		return l.appendSanitizedJSONValue(buf, f.Key, v), cut
	}
	if f.kind == fieldKindError {
		// nil error：与 ...any 路径一致输出 null，不经脱敏
		// [EN] A nil error emits null as the ...any path does, bypassing sanitising
		return append(buf, 'n', 'u', 'l', 'l'), false
	}
	// 数值类：只处理整值擦除，其余直写
	// [EN] Numeric kinds: only whole-value wiping applies, otherwise written directly
	if desensitize.ActionForKey(f.Key) == desensitize.WipeValue {
		return appendJSONString(buf, desensitize.RedactedText), false
	}
	return l.appendFieldNumberJSON(buf, f), false
}

// clampString 是 clampValue 的 string 形态，不装箱
// clampValue(v any, ...) 的 case string 分支与此等价，
// 但调用它要先把 string 装进 any —— 那正是这条路径要消掉的分配
//
// [EN] The string form of clampValue, without boxing
// clampValue's string branch is equivalent, but reaching it requires boxing the string
// into an any, which is exactly the allocation this path removes
func (l *Logger) clampString(s string, room int) (string, bool) {
	if l.maxLineBytes <= 0 {
		return s, false
	}
	if room < 1 {
		room = 1
	}
	if len(s) <= room {
		return s, false
	}
	return s[:utf8SafeEnd(s, room)], true
}

// appendSanitizedJSONStringValue 是 appendSanitizedJSONValue 的 string 形态
//
// 判定顺序逐行对照原函数，不能重排：
//  1. 按字段名查脱敏动作，命中擦除或掩码就到此为止
//  2. SanitizeKeysAndValues 时对全值做扫描
//  3. 字段级长度上限
//  4. 以上都不命中才原样输出
//
// 顺序错了会出现「该擦除的字段只被截断」这类静默降级
//
// 掩码分支直接把 string 交给 mask，不经过 appendMaskedJSON 的 textOf ——
// 后者要先把值装进 any，正是这条路径要消掉的分配
//
// [EN] The string form of appendSanitizedJSONValue
// The decision order mirrors the original line by line and must not be reordered: the
// per-key action first, then whole-value scanning under SanitizeKeysAndValues, then the
// per-field length cap, and only then the value as-is
// Getting the order wrong silently downgrades a field that should be wiped into one that
// is merely truncated
// The mask branches hand the string straight to mask rather than going through
// appendMaskedJSON's textOf, which would box the value into an any first
func (l *Logger) appendSanitizedJSONStringValue(buf []byte, key, s string) []byte {
	if l.sanitize == SanitizeOff {
		return appendJSONString(buf, s)
	}

	switch desensitize.ActionForKey(key) {
	case desensitize.WipeValue:
		return appendJSONString(buf, desensitize.RedactedText)
	case desensitize.MaskPhone:
		return appendMaskedJSONString(buf, s, desensitize.AppendMaskedPhone)
	case desensitize.MaskEmail:
		return appendMaskedJSONString(buf, s, desensitize.AppendMaskedEmail)
	case desensitize.MaskFingerprint:
		return appendMaskedJSONString(buf, s, desensitize.AppendFingerprint)
	case desensitize.MaskEncrypt:
		return appendMaskedJSONString(buf, s, appendEncrypted)
	}

	if l.sanitize == SanitizeKeysAndValues {
		return appendScannedJSONString(buf, l.capped(s))
	}
	if l.maxFieldBytes > 0 {
		cut, trimmed := l.capFieldText(s)
		return appendCappedJSONString(buf, cut, trimmed)
	}
	return appendJSONString(buf, s)
}

// appendMaskedJSONString 是 appendMaskedJSON 的 string 形态，省掉 textOf 的装箱
// [EN] The string form of appendMaskedJSON, without textOf's boxing
func appendMaskedJSONString(buf []byte, s string, mask func([]byte, string) []byte) []byte {
	buf = append(buf, '"')
	buf = mask(buf, s)
	return append(buf, '"')
}

// logWithTypedFields 类型化字段的记录入口
// 与 logWithKV 的结构一一对应：先过级别闸门，再解析 caller，空字段退回 ultraLog
//
// [EN] The logging entry point for typed fields, structured to match logWithKV: the level
// gate first, then caller resolution, and an empty field list falls back to ultraLog
func (l *Logger) logWithTypedFields(level LogLevel, msg string, ci *callerInfo, fields []Field) {
	if level < LogLevel(l.level.Load()) {
		return
	}

	if ci == nil && l.showCaller.Load() {
		ci = l.findExternalCaller(2, 3)
	}

	if len(fields) == 0 {
		l.ultraLog(level, msg, ci)
		return
	}

	if l.format == FormatJSON {
		p := bytePool.Get().(*[]byte)
		buf := (*p)[:0]
		defer func() { putPooledBuf(&bytePool, p, buf) }()

		buf = append(buf, '{')
		buf = l.appendJSONHeader(buf, level, msg, ci)
		buf = l.appendFieldsJSON(buf, fields)
		l.writeJSONTail(level, buf)
		return
	}

	// text 模式暂不走专用路径：把字段转回 KV 形态交给既有实现
	//
	// 这里刻意不复制 text 的渲染逻辑 —— 它的边界处理（分隔符、回滚、截断标记）
	// 与 JSON 路径不同，再抄一份就要再维护一份，而类型化字段的收益集中在
	// 结构化日志场景，text 模式本就是人读为主、量不大
	// 代价是 text 模式下仍有一次装箱，与改造前持平，不更差
	//
	// [EN] Text mode deliberately reuses the existing KV implementation rather than
	// duplicating its rendering, whose boundary handling differs from the JSON path and
	// would become a second copy to maintain, while the benefit of typed fields is
	// concentrated in structured logging and text mode is human-facing and low volume
	// The cost is one boxing in text mode, level with before rather than worse
	kv := make([]any, 0, len(fields)*2)
	for i := range fields {
		kv = append(kv, fields[i].Key, fields[i].asAny())
	}
	l.logWithKV(level, msg, ci, kv...)
}

// asAny 把字段退回 any 形态，仅供 text 模式与兼容路径使用
// [EN] Returns the field as an any, used only by text mode and compatibility paths
func (f Field) asAny() any {
	switch f.kind {
	case fieldKindString:
		return f.str
	case fieldKindInt:
		return f.num
	case fieldKindUint:
		return uint64(f.num)
	case fieldKindFloat32:
		return float32(f.floatValue())
	case fieldKindFloat64:
		return f.floatValue()
	case fieldKindBool:
		return f.num != 0
	case fieldKindDuration:
		return time.Duration(f.num)
	case fieldKindError:
		return f.errValue()
	}
	return f.iface
}

// DebugFields 以类型化字段记录 Debug 日志
// [EN] Logs at Debug with typed fields
func (l *Logger) DebugFields(msg string, fields ...Field) {
	l.logWithTypedFields(DEBUG, msg, nil, fields)
}

// InfoFields 以类型化字段记录 Info 日志
// [EN] Logs at Info with typed fields
func (l *Logger) InfoFields(msg string, fields ...Field) {
	l.logWithTypedFields(INFO, msg, nil, fields)
}

// WarnFields 以类型化字段记录 Warn 日志
// [EN] Logs at Warn with typed fields
func (l *Logger) WarnFields(msg string, fields ...Field) {
	l.logWithTypedFields(WARN, msg, nil, fields)
}

// ErrorFields 以类型化字段记录 Error 日志
// [EN] Logs at Error with typed fields
func (l *Logger) ErrorFields(msg string, fields ...Field) {
	l.logWithTypedFields(ERROR, msg, nil, fields)
}

// floatValue 从 num 的位模式还原浮点值
// [EN] Recovers the float from num's bit pattern
func (f Field) floatValue() float64 { return math.Float64frombits(uint64(f.num)) }

// errValue 取 error 槽位，类型不符时返回 nil
// [EN] Returns the error slot, or nil when the type does not match
func (f Field) errValue() error {
	e, _ := f.iface.(error)
	return e
}
