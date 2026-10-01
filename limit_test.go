/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-26 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-01-26 00:00:00
 * @FilePath: \go-logger\limit_test.go
 * @Description: 单条日志字节上限与超限裁剪测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestLimit_DefaultCapsHugeLine 默认上限必须挡住失控的大 payload
// 回归点：无上限时一个 200KB 字段原样写出 204894 字节，且缓冲超过池阈值后
// 每条大日志都是一次大块分配加一次丢弃
func TestLimit_DefaultCapsHugeLine(t *testing.T) {
	var buf bytes.Buffer
	l := jsonFieldLogger(&buf)

	l.InfoKV("dump", "payload", strings.Repeat("x", 200*1024))

	if got := buf.Len(); got > DefaultMaxLineBytes+8 {
		t.Fatalf("默认上限未生效：写出 %d 字节，上限 %d", got, DefaultMaxLineBytes)
	}
	if !strings.Contains(buf.String(), `"truncated":true`) {
		t.Errorf("超限条目应补 truncated 标记\n原文前 200 字节: %.200s", buf.String())
	}
}

// TestLimit_JSONStaysValid 超限裁剪后 JSON 必须仍然合法可解析
// 这是裁剪最关键的约束：从字段中间切断会产出无法解析的行，
// 下游解析器可能因此丢掉整个批次
func TestLimit_JSONStaysValid(t *testing.T) {
	for _, limit := range []int{64, 128, 256, 1024, 4096} {
		var buf bytes.Buffer
		l := jsonFieldLogger(&buf).WithMaxLineBytes(limit)

		l.InfoKV("dump",
			"a", strings.Repeat("1", 100),
			"b", strings.Repeat("2", 100),
			"c", strings.Repeat("3", 100),
			"d", strings.Repeat("4", 100))

		line := lastLine(t, &buf)
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("上限 %d 时输出非合法 JSON: %v\n原文: %s", limit, err, line)
		}
	}
}

// TestLimit_HeaderSurvives 骨架字段（级别、消息）最先写出，故最后被牺牲
func TestLimit_HeaderSurvives(t *testing.T) {
	var buf bytes.Buffer
	l := jsonFieldLogger(&buf).WithMaxLineBytes(200)

	l.ErrorKV("db failed", "stack", strings.Repeat("frame ", 200), "code", "E500")

	var m map[string]any
	if err := json.Unmarshal([]byte(lastLine(t, &buf)), &m); err != nil {
		t.Fatalf("输出非合法 JSON: %v", err)
	}
	if m["level"] != "ERROR" || m["message"] != "db failed" {
		t.Errorf("骨架字段应优先保住，实得 level=%v message=%v", m["level"], m["message"])
	}
	if m["truncated"] != true {
		t.Errorf("应补 truncated 标记，实得 %#v", m["truncated"])
	}
}

// TestLimit_FieldOrderPreserved 先写的字段先保住，顺序可预期
func TestLimit_FieldOrderPreserved(t *testing.T) {
	var buf bytes.Buffer
	l := jsonFieldLogger(&buf).WithMaxLineBytes(220)

	l.InfoKV("msg", "keep", "v1", "bulk", strings.Repeat("x", 500), "drop", "v2")

	var m map[string]any
	if err := json.Unmarshal([]byte(lastLine(t, &buf)), &m); err != nil {
		t.Fatalf("输出非合法 JSON: %v", err)
	}
	if m["keep"] != "v1" {
		t.Errorf("先写的字段应保住，实得 %#v", m["keep"])
	}
	if _, ok := m["drop"]; ok {
		t.Errorf("超限后的字段应被丢弃，实得 %#v", m["drop"])
	}
}

// TestLimit_MaxFieldBytes 单字段上限就地截断值本体并留标记
func TestLimit_MaxFieldBytes(t *testing.T) {
	var buf bytes.Buffer
	l := jsonFieldLogger(&buf).WithMaxFieldBytes(32)

	l.InfoKV("msg", "small", "ok", "big", strings.Repeat("y", 500))

	var m map[string]any
	if err := json.Unmarshal([]byte(lastLine(t, &buf)), &m); err != nil {
		t.Fatalf("输出非合法 JSON: %v", err)
	}
	if m["small"] != "ok" {
		t.Errorf("未超限字段应原样保留，实得 %#v", m["small"])
	}
	big, _ := m["big"].(string)
	if !strings.HasSuffix(big, truncatedValueMark) {
		t.Errorf("超限字段应带 truncated 标记，实得 %q", big)
	}
	if len(big) > 32+len(truncatedValueMark) {
		t.Errorf("截断后仍超限：%d 字节", len(big))
	}
}

// TestLimit_UTF8Boundary 截断必须落在 UTF-8 边界，不得产出非法字节序列
func TestLimit_UTF8Boundary(t *testing.T) {
	var buf bytes.Buffer
	// 上限取 10：中文 3 字节一个，10 不是 3 的整数倍，必然触发边界回退
	l := jsonFieldLogger(&buf).WithMaxFieldBytes(10)

	l.InfoKV("msg", "text", strings.Repeat("中", 50))

	line := lastLine(t, &buf)
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("输出非合法 JSON（截断可能切断了多字节字符）: %v\n原文: %s", err, line)
	}
	text, _ := m["text"].(string)
	body := strings.TrimSuffix(text, truncatedValueMark)
	if strings.ContainsRune(body, '�') {
		t.Errorf("截断产出了非法 UTF-8，实得 %q", body)
	}
	if len(body)%3 != 0 {
		t.Errorf("截断未落在字符边界：%d 字节", len(body))
	}
}

// TestLimit_LongMessage 超长 msg 单独也能撑爆条目，须一并受限
func TestLimit_LongMessage(t *testing.T) {
	var buf bytes.Buffer
	l := jsonFieldLogger(&buf).WithMaxLineBytes(512)

	l.InfoMsg(strings.Repeat("m", 100*1024))

	if got := buf.Len(); got > 512+8 {
		t.Fatalf("超长 msg 未受限：写出 %d 字节", got)
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(lastLine(t, &buf)), &m); err != nil {
		t.Fatalf("输出非合法 JSON: %v", err)
	}
}

// TestLimit_Disabled 上限设为 0 时不做任何裁剪
func TestLimit_Disabled(t *testing.T) {
	var buf bytes.Buffer
	l := jsonFieldLogger(&buf).WithMaxLineBytes(0)

	payload := strings.Repeat("z", 200*1024)
	l.InfoKV("dump", "payload", payload)

	if !strings.Contains(buf.String(), payload) {
		t.Error("上限关闭后应原样写出")
	}
}

// TestLimit_TextMode text 模式同样受限且补标记
func TestLimit_TextMode(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger().WithOutput(&buf).WithFormat(FormatText).
		WithColorful(false).WithShowCaller(false).WithMaxLineBytes(200)

	l.InfoKV("dump", "a", "keep", "b", strings.Repeat("x", 500))

	out := buf.String()
	if len(out) > 200+len(truncatedFieldText)+2 {
		t.Fatalf("text 模式上限未生效：%d 字节", len(out))
	}
	if !strings.Contains(out, truncatedFieldText) {
		t.Errorf("text 模式超限应补标记\n原文: %s", out)
	}
	if strings.Count(out, "\n") != 1 {
		t.Errorf("裁剪后仍须恰为一行，实得 %d 个换行", strings.Count(out, "\n"))
	}
}
