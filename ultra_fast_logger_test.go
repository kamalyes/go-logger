/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-24 10:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-24 10:30:00
 * @FilePath: \go-logger\ultra_fast_logger_test.go
 * @Description: 极致性能日志器测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

// TestUltraFastLogger_BasicLogging 测试基本日志功能
func TestUltraFastLogger_BasicLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	tests := []struct {
		name     string
		logFunc  func(string, ...interface{})
		format   string
		args     []interface{}
		contains string
	}{
		{"Debug", logger.Debug, "Debug message: %s", []interface{}{"test"}, "Debug message: test"},
		{"Info", logger.Info, "Info message: %d", []interface{}{42}, "Info message: 42"},
		{"Warn", logger.Warn, "Warn message", nil, "Warn message"},
		{"Error", logger.Error, "Error: %v", []interface{}{"error"}, "Error: error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.format, tt.args...)
			output := buf.String()
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Expected output to contain %q, got %q", tt.contains, output)
			}
		})
	}
}

// TestUltraFastLogger_LevelFiltering 测试日志级别过滤
func TestUltraFastLogger_LevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, WARN)

	logger.Debug("debug message")
	if buf.Len() != 0 {
		t.Error("Debug message should be filtered")
	}

	logger.Info("info message")
	if buf.Len() != 0 {
		t.Error("Info message should be filtered")
	}

	logger.Warn("warn message")
	if buf.Len() == 0 {
		t.Error("Warn message should be logged")
	}
}

// TestUltraFastLogger_MessageMethods 测试纯文本日志方法
func TestUltraFastLogger_MessageMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	tests := []struct {
		name     string
		logFunc  func(string)
		msg      string
		contains string
	}{
		{"DebugMsg", logger.DebugMsg, "Debug message", "Debug message"},
		{"InfoMsg", logger.InfoMsg, "Info message", "Info message"},
		{"WarnMsg", logger.WarnMsg, "Warn message", "Warn message"},
		{"ErrorMsg", logger.ErrorMsg, "Error message", "Error message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.msg)
			output := buf.String()
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Expected output to contain %q, got %q", tt.contains, output)
			}
		})
	}
}

// TestUltraFastLogger_ContextLogging 测试上下文日志
func TestUltraFastLogger_ContextLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-123")
	ctx = context.WithValue(ctx, "request_id", "req-456")

	logger.InfoContext(ctx, "Context message")
	output := buf.String()

	if !strings.Contains(output, "TraceID=trace-123") {
		t.Errorf("Expected output to contain trace ID, got %q", output)
	}
	if !strings.Contains(output, "RequestID=req-456") {
		t.Errorf("Expected output to contain request ID, got %q", output)
	}
	if !strings.Contains(output, "Context message") {
		t.Errorf("Expected output to contain message, got %q", output)
	}
}

// TestUltraFastLogger_KeyValueLogging 测试键值对日志
func TestUltraFastLogger_KeyValueLogging(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	logger.InfoKV("User action", "user", "john", "action", "login", "status", "success")
	output := buf.String()

	if !strings.Contains(output, "User action") {
		t.Errorf("Expected output to contain message, got %q", output)
	}
	if !strings.Contains(output, "user: john") {
		t.Errorf("Expected output to contain user field, got %q", output)
	}
	if !strings.Contains(output, "action: login") {
		t.Errorf("Expected output to contain action field, got %q", output)
	}
	if !strings.Contains(output, "status: success") {
		t.Errorf("Expected output to contain status field, got %q", output)
	}
}

// TestUltraFastLogger_WithFields 测试字段链式调用
func TestUltraFastLogger_WithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	fields := map[string]interface{}{
		"user":   "alice",
		"age":    30,
		"active": true,
	}

	logger.WithFields(fields).Info("User info")
	output := buf.String()

	if !strings.Contains(output, "User info") {
		t.Errorf("Expected output to contain message, got %q", output)
	}
	if !strings.Contains(output, "user: alice") {
		t.Errorf("Expected output to contain user field, got %q", output)
	}
	if !strings.Contains(output, "age: 30") {
		t.Errorf("Expected output to contain age field, got %q", output)
	}
	if !strings.Contains(output, "active: true") {
		t.Errorf("Expected output to contain active field, got %q", output)
	}
}

// TestUltraFastLogger_WithError 测试错误字段
func TestUltraFastLogger_WithError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	err := errors.New("test error")
	logger.WithError(err).Error("Operation failed")
	output := buf.String()

	if !strings.Contains(output, "Operation failed") {
		t.Errorf("Expected output to contain message, got %q", output)
	}
	if !strings.Contains(output, "error: test error") {
		t.Errorf("Expected output to contain error field, got %q", output)
	}
}

// TestUltraFastLogger_Clone 测试克隆功能
func TestUltraFastLogger_Clone(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)

	cloned := logger.Clone()
	clonedLogger, ok := cloned.(*UltraFastLogger)
	if !ok {
		t.Fatal("Clone should return *UltraFastLogger")
	}

	if clonedLogger.GetLevel() != logger.GetLevel() {
		t.Error("Cloned logger should have same level")
	}

	// 修改克隆不应影响原始
	clonedLogger.SetLevel(DEBUG)
	if logger.GetLevel() == DEBUG {
		t.Error("Original logger should not be affected by clone modification")
	}
}

// TestUltraFastLogger_PrintMethods 测试 Print 兼容方法
func TestUltraFastLogger_PrintMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	tests := []struct {
		name     string
		logFunc  func()
		contains string
	}{
		{"Print", func() { logger.Print("test message") }, "test message"},
		{"Printf", func() { logger.Printf("test %s", "message") }, "test message"},
		{"Println", func() { logger.Println("test", "message") }, "test message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc()
			output := buf.String()
			if !strings.Contains(output, tt.contains) {
				t.Errorf("Expected output to contain %q, got %q", tt.contains, output)
			}
		})
	}
}

// TestUltraFastLogger_LogWithFields 测试带字段的日志条目
func TestUltraFastLogger_LogWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	fields := map[string]interface{}{
		"component": "database",
		"duration":  42,
	}

	logger.LogWithFields(INFO, "Query executed", fields)
	output := buf.String()

	if !strings.Contains(output, "Query executed") {
		t.Errorf("Expected output to contain message, got %q", output)
	}
	if !strings.Contains(output, "component: database") {
		t.Errorf("Expected output to contain component field, got %q", output)
	}
	if !strings.Contains(output, "duration: 42") {
		t.Errorf("Expected output to contain duration field, got %q", output)
	}
}

// TestUltraFastLogger_IsLevelEnabled 测试级别检查
func TestUltraFastLogger_IsLevelEnabled(t *testing.T) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, WARN)

	tests := []struct {
		level   LogLevel
		enabled bool
	}{
		{DEBUG, false},
		{INFO, false},
		{WARN, true},
		{ERROR, true},
		{FATAL, true},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			if logger.IsLevelEnabled(tt.level) != tt.enabled {
				t.Errorf("IsLevelEnabled(%s) = %v, want %v", tt.level, !tt.enabled, tt.enabled)
			}
		})
	}
}

// TestUltraFastLogger_ShowCaller 测试调用者信息配置
func TestUltraFastLogger_ShowCaller(t *testing.T) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, DEBUG)

	// 默认不显示调用者
	if logger.IsShowCaller() {
		t.Error("ShowCaller should be false by default")
	}

	logger.SetShowCaller(true)
	if !logger.IsShowCaller() {
		t.Error("ShowCaller should be true after SetShowCaller(true)")
	}

	logger.SetShowCaller(false)
	if logger.IsShowCaller() {
		t.Error("ShowCaller should be false after SetShowCaller(false)")
	}
}

// TestUltraFastLogger_Colorful 测试彩色输出
func TestUltraFastLogger_Colorful(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Colorful = true
	config.Output = buf
	logger := NewUltraFastLogger(config)

	logger.Info("Colorful message")
	output := buf.String()

	// 彩色输出应包含 ANSI 转义码
	if !strings.Contains(output, "\033[") {
		t.Errorf("Expected colorful output to contain ANSI codes, got %q", output)
	}
}

// TestUltraFieldLogger_ChainedFields 测试链式字段调用
func TestUltraFieldLogger_ChainedFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	logger.WithField("user", "bob").
		WithField("age", 25).
		WithField("role", "admin").
		Info("User details")

	output := buf.String()

	if !strings.Contains(output, "user: bob") {
		t.Errorf("Expected output to contain user field, got %q", output)
	}
	if !strings.Contains(output, "age: 25") {
		t.Errorf("Expected output to contain age field, got %q", output)
	}
	if !strings.Contains(output, "role: admin") {
		t.Errorf("Expected output to contain role field, got %q", output)
	}
}

// TestUltraFastLogger_EmptyContext 测试空上下文
func TestUltraFastLogger_EmptyContext(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	// nil context
	logger.InfoContext(nil, "Message with nil context")
	output := buf.String()
	if !strings.Contains(output, "Message with nil context") {
		t.Errorf("Expected output to contain message, got %q", output)
	}

	// empty context
	buf.Reset()
	ctx := context.Background()
	logger.InfoContext(ctx, "Message with empty context")
	output = buf.String()
	if !strings.Contains(output, "Message with empty context") {
		t.Errorf("Expected output to contain message, got %q", output)
	}
}

// TestUltraFastLogger_OddKeyValues 测试奇数个键值对
func TestUltraFastLogger_OddKeyValues(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, DEBUG)

	// 奇数个参数,最后一个键应该有默认值
	logger.InfoKV("Odd key-values", "key1", "value1", "key2")
	output := buf.String()

	if !strings.Contains(output, "key1: value1") {
		t.Errorf("Expected output to contain key1 field, got %q", output)
	}
	if !strings.Contains(output, "key2:") {
		t.Errorf("Expected output to contain key2 field, got %q", output)
	}
}

// BenchmarkUltraFastLogger_SimpleLog 基准测试 - 简单日志
func BenchmarkUltraFastLogger_SimpleLog(b *testing.B) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, INFO)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Benchmark message")
		}
	})
}

// BenchmarkUltraFastLogger_FormattedLog 基准测试 - 格式化日志
func BenchmarkUltraFastLogger_FormattedLog(b *testing.B) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, INFO)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("User: %s, Age: %d, Active: %v", "john", 30, true)
		}
	})
}

// BenchmarkUltraFastLogger_WithFields 基准测试 - 字段日志
func BenchmarkUltraFastLogger_WithFields(b *testing.B) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, INFO)
	fields := map[string]interface{}{
		"user":   "john",
		"age":    30,
		"active": true,
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.LogWithFields(INFO, "User action", fields)
		}
	})
}

// BenchmarkUltraFastLogger_KeyValue 基准测试 - 键值对日志
func BenchmarkUltraFastLogger_KeyValue(b *testing.B) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, INFO)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.InfoKV("User action", "user", "john", "age", 30, "active", true)
		}
	})
}

// BenchmarkUltraFastLogger_Context 基准测试 - 上下文日志
func BenchmarkUltraFastLogger_Context(b *testing.B) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, INFO)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-123")
	ctx = context.WithValue(ctx, "request_id", "req-456")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.InfoContext(ctx, "Context message")
		}
	})
}
