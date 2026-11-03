/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-03 17:02:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-03 17:32:00
 * @FilePath: \go-logger\fields_typed_test.go
 * @Description: 类型化字段与 ...any 路径的逐字节等价测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package logger

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// fieldCase 一个字段的两种写法：类型化的 Field 与等价的 KV 对
type fieldCase struct {
	name  string
	field Field
	key   string
	val   any
}

func fieldCases() []fieldCase {
	err := errors.New("boom")
	return []fieldCase{
		{"string", Str("path", "/v1/x"), "path", "/v1/x"},
		{"string空", Str("path", ""), "path", ""},
		{"string带引号", Str("msg", `he said "hi"`), "msg", `he said "hi"`},
		{"string带换行", Str("msg", "a\nb"), "msg", "a\nb"},
		{"string中文", Str("msg", "中文内容"), "msg", "中文内容"},
		{"string带反斜杠", Str("p", `C:\tmp`), "p", `C:\tmp`},
		{"int", Int("status", 200), "status", int64(200)},
		{"int负数", Int("delta", -7), "delta", int64(-7)},
		{"int零", Int("n", 0), "n", int64(0)},
		{"int64大值", Int("big", int64(1)<<62), "big", int64(1) << 62},
		{"uint", Uint("bytes", uint64(4096)), "bytes", uint64(4096)},
		{"uint零", Uint("bytes", uint64(0)), "bytes", uint64(0)},
		{"float64", Float64("ratio", 1.5), "ratio", 1.5},
		{"float64整数值", Float64("ratio", 2), "ratio", float64(2)},
		{"float32", Float32("r32", 1.5), "r32", float32(1.5)},
		{"bool真", Bool("slow", true), "slow", true},
		{"bool假", Bool("slow", false), "slow", false},
		{"duration", Dur("cost", 1500*time.Microsecond), "cost", 1500 * time.Microsecond},
		{"duration零", Dur("cost", 0), "cost", time.Duration(0)},
		{"error", Err("error", err), "error", err},
		{"error为nil", Err("error", nil), "error", error(nil)},
		{"兜底any", Anyf("other", []string{"a"}), "other", []string{"a"}},
		// 脱敏相关：字段名命中擦除与掩码规则
		{"脱敏擦除", Str("password", "s3cr3t"), "password", "s3cr3t"},
		{"脱敏手机", Str("phone", "13800138000"), "phone", "13800138000"},
		{"脱敏邮箱", Str("email", "a@b.com"), "email", "a@b.com"},
	}
}

// captureJSON 用同一个 logger 配置分别走两条路径，返回去掉时间戳后的输出
func captureJSON(t *testing.T, cfg func(*Logger) *Logger, emit func(*Logger)) string {
	t.Helper()
	var buf bytes.Buffer
	l := NewLogger().WithFormat(FormatJSON).WithShowCaller(false).WithOutput(&buf)
	if cfg != nil {
		l = cfg(l)
	}
	emit(l)
	return stripVolatile(buf.String())
}

// stripVolatile 去掉每次都不同的时间戳，其余必须逐字节相同
func stripVolatile(s string) string {
	out := make([]string, 0, 4)
	for _, line := range strings.Split(strings.TrimRight(s, "\n"), "\n") {
		if i := strings.Index(line, `"timestamp":"`); i >= 0 {
			if j := strings.Index(line[i+13:], `"`); j >= 0 {
				line = line[:i+13] + "<ts>" + line[i+13+j:]
			}
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}

// TestFieldsMatchKV 类型化字段的输出必须与等价的 ...any 调用逐字节相同
//
// 这是整个类型化字段改动的正确性闸门：它换的是调用方的分配方式，
// 不该换日志的格式 格式一变，下游的日志解析就会碎，而且是静默的
// 逐个字段单独比一遍，再合起来比一遍 —— 单独比能指出是哪个类型错了
//
// [EN] Typed-field output must be byte-for-byte identical to the equivalent ...any call
// This is the correctness gate for the whole change: it alters how the caller allocates,
// not how the log is formatted, and a format change would silently break downstream
// parsing
// Each field is compared on its own and then all together, so a failure names the type
func TestFieldsMatchKV(t *testing.T) {
	for _, tc := range fieldCases() {
		t.Run(tc.name, func(t *testing.T) {
			got := captureJSON(t, nil, func(l *Logger) { l.InfoFields("m", tc.field) })
			want := captureJSON(t, nil, func(l *Logger) { l.InfoKV("m", tc.key, tc.val) })
			if got != want {
				t.Errorf("输出不一致\n类型化: %s\n...any : %s", got, want)
			}
		})
	}

	t.Run("全部字段合起来", func(t *testing.T) {
		cases := fieldCases()
		fields := make([]Field, 0, len(cases))
		kv := make([]any, 0, len(cases)*2)
		for _, tc := range cases {
			fields = append(fields, tc.field)
			kv = append(kv, tc.key, tc.val)
		}
		got := captureJSON(t, nil, func(l *Logger) { l.InfoFields("m", fields...) })
		want := captureJSON(t, nil, func(l *Logger) { l.InfoKV("m", kv...) })
		if got != want {
			t.Errorf("输出不一致\n类型化: %s\n...any : %s", got, want)
		}
	})
}

// TestFieldsMatchKVUnderSanitizeModes 三种脱敏模式下都必须等价
// 脱敏的判定顺序最容易在复制实现时被重排，重排的表现是「该擦除的只被截断」
//
// [EN] Equivalence must hold under every sanitising mode
// The decision order is the easiest thing to reorder when duplicating an implementation,
// and the symptom is a field that should be wiped merely getting truncated
func TestFieldsMatchKVUnderSanitizeModes(t *testing.T) {
	modes := []struct {
		name string
		mode SanitizeMode
	}{
		{"关闭", SanitizeOff},
		{"按键名", SanitizeKeys},
		{"键名与值", SanitizeKeysAndValues},
	}
	for _, m := range modes {
		for _, tc := range fieldCases() {
			t.Run(fmt.Sprintf("%s/%s", m.name, tc.name), func(t *testing.T) {
				cfg := func(l *Logger) *Logger { return l.WithSanitize(m.mode) }
				got := captureJSON(t, cfg, func(l *Logger) { l.InfoFields("m", tc.field) })
				want := captureJSON(t, cfg, func(l *Logger) { l.InfoKV("m", tc.key, tc.val) })
				if got != want {
					t.Errorf("模式 %s 下输出不一致\n类型化: %s\n...any : %s", m.name, got, want)
				}
			})
		}
	}
}

// TestFieldsMatchKVUnderLineBudget 行长度上限下的截断与回滚也必须等价
// [EN] Truncation and rollback under the line budget must match too
func TestFieldsMatchKVUnderLineBudget(t *testing.T) {
	long := strings.Repeat("x", 500)
	for _, limit := range []int{0, 64, 128, 256, 1024} {
		t.Run(fmt.Sprintf("上限%d", limit), func(t *testing.T) {
			cfg := func(l *Logger) *Logger { return l.WithMaxLineBytes(limit) }
			got := captureJSON(t, cfg, func(l *Logger) {
				l.InfoFields("m", Str("a", long), Int("b", 1), Str("c", long))
			})
			want := captureJSON(t, cfg, func(l *Logger) {
				l.InfoKV("m", "a", long, "b", int64(1), "c", long)
			})
			if got != want {
				t.Errorf("上限 %d 下输出不一致\n类型化: %s\n...any : %s", limit, got, want)
			}
		})
	}
}

// TestFieldsLevelGate 级别闸门要挡住低于当前级别的调用
// [EN] The level gate must drop calls below the current level
func TestFieldsLevelGate(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger().WithFormat(FormatJSON).WithShowCaller(false).WithOutput(&buf).WithLevel(ERROR)
	l.InfoFields("不该出现", Str("k", "v"))
	l.DebugFields("不该出现", Str("k", "v"))
	if buf.Len() != 0 {
		t.Errorf("级别闸门没挡住: %s", buf.String())
	}
	l.ErrorFields("该出现", Str("k", "v"))
	if !strings.Contains(buf.String(), "该出现") {
		t.Errorf("ERROR 级别被误挡: %s", buf.String())
	}
}

// TestFieldsEmpty 没有字段时退回无字段条目，与 InfoKV 空参一致
// [EN] With no fields the entry falls back to the field-less form, matching empty InfoKV
func TestFieldsEmpty(t *testing.T) {
	got := captureJSON(t, nil, func(l *Logger) { l.InfoFields("m") })
	want := captureJSON(t, nil, func(l *Logger) { l.InfoKV("m") })
	if got != want {
		t.Errorf("空字段输出不一致\n类型化: %s\n...any : %s", got, want)
	}
}
