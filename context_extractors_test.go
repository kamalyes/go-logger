/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-24 11:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-24 11:30:00
 * @FilePath: \go-logger\context_extractors_test.go
 * @Description: 上下文提取器测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestNoOpContextExtractor 测试空操作提取器
func TestNoOpContextExtractor(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)
	logger.SetContextExtractor(NoOpContextExtractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-123")

	logger.InfoContext(ctx, "Test message")
	output := buf.String()

	// 不应包含上下文信息
	if strings.Contains(output, "trace-123") {
		t.Errorf("NoOpContextExtractor should not extract context info, got: %s", output)
	}
	if !strings.Contains(output, "Test message") {
		t.Errorf("Output should contain message, got: %s", output)
	}
}

// TestSimpleTraceIDExtractor 测试 TraceID 提取器
func TestSimpleTraceIDExtractor(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)
	logger.SetContextExtractor(SimpleTraceIDExtractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-456")
	ctx = context.WithValue(ctx, "request_id", "req-789")

	logger.InfoContext(ctx, "Test message")
	output := buf.String()

	// 应包含 TraceID
	if !strings.Contains(output, "TraceID=trace-456") {
		t.Errorf("Expected TraceID in output, got: %s", output)
	}
	// 不应包含 RequestID
	if strings.Contains(output, "RequestID") {
		t.Errorf("Should not contain RequestID, got: %s", output)
	}
}

// TestSimpleRequestIDExtractor 测试 RequestID 提取器
func TestSimpleRequestIDExtractor(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)
	logger.SetContextExtractor(SimpleRequestIDExtractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-456")
	ctx = context.WithValue(ctx, "request_id", "req-789")

	logger.InfoContext(ctx, "Test message")
	output := buf.String()

	// 应包含 RequestID
	if !strings.Contains(output, "RequestID=req-789") {
		t.Errorf("Expected RequestID in output, got: %s", output)
	}
	// 不应包含 TraceID
	if strings.Contains(output, "TraceID") {
		t.Errorf("Should not contain TraceID, got: %s", output)
	}
}

// TestCustomFieldExtractor 测试自定义字段提取器
func TestCustomFieldExtractor(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)

	extractor := CustomFieldExtractor(
		[]string{"user_id", "session_id"},
		[]string{},
	)
	logger.SetContextExtractor(extractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", "user-123")
	ctx = context.WithValue(ctx, "session_id", "sess-456")
	ctx = context.WithValue(ctx, "ignored", "should-not-appear")

	logger.InfoContext(ctx, "Test message")
	output := buf.String()

	if !strings.Contains(output, "user_id=user-123") {
		t.Errorf("Expected user_id in output, got: %s", output)
	}
	if !strings.Contains(output, "session_id=sess-456") {
		t.Errorf("Expected session_id in output, got: %s", output)
	}
	if strings.Contains(output, "ignored") {
		t.Errorf("Should not contain ignored field, got: %s", output)
	}
}

// TestChainContextExtractors 测试链式提取器
func TestChainContextExtractors(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)

	extractor := ChainContextExtractors(
		SimpleTraceIDExtractor,
		SimpleRequestIDExtractor,
		ExtractFromContextValue("user_id", "User"),
	)
	logger.SetContextExtractor(extractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-111")
	ctx = context.WithValue(ctx, "request_id", "req-222")
	ctx = context.WithValue(ctx, "user_id", "alice")

	logger.InfoContext(ctx, "Test message")
	output := buf.String()

	if !strings.Contains(output, "TraceID=trace-111") {
		t.Errorf("Expected TraceID in output, got: %s", output)
	}
	if !strings.Contains(output, "RequestID=req-222") {
		t.Errorf("Expected RequestID in output, got: %s", output)
	}
	if !strings.Contains(output, "User=alice") {
		t.Errorf("Expected User in output, got: %s", output)
	}
}

// TestConditionalContextExtractor 测试条件提取器
func TestConditionalContextExtractor(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)

	extractor := ConditionalContextExtractor(
		func(ctx context.Context) bool {
			env, ok := ctx.Value("env").(string)
			return ok && env == "production"
		},
		SimpleTraceIDExtractor,
	)
	logger.SetContextExtractor(extractor)

	// 测试生产环境
	prodCtx := context.Background()
	prodCtx = context.WithValue(prodCtx, "env", "production")
	prodCtx = context.WithValue(prodCtx, "trace_id", "trace-prod")

	buf.Reset()
	logger.InfoContext(prodCtx, "Production message")
	output := buf.String()

	if !strings.Contains(output, "TraceID=trace-prod") {
		t.Errorf("Expected TraceID in production, got: %s", output)
	}

	// 测试开发环境
	devCtx := context.Background()
	devCtx = context.WithValue(devCtx, "env", "development")
	devCtx = context.WithValue(devCtx, "trace_id", "trace-dev")

	buf.Reset()
	logger.InfoContext(devCtx, "Development message")
	output = buf.String()

	if strings.Contains(output, "TraceID") {
		t.Errorf("Should not contain TraceID in development, got: %s", output)
	}
}

// TestContextExtractorBuilder 测试构建器
func TestContextExtractorBuilder(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)

	extractor := NewContextExtractorBuilder().
		AddTraceID().
		AddRequestID().
		AddContextValue("tenant_id", "Tenant").
		Build()

	logger.SetContextExtractor(extractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-999")
	ctx = context.WithValue(ctx, "request_id", "req-888")
	ctx = context.WithValue(ctx, "tenant_id", "tenant-A")

	logger.InfoContext(ctx, "Test message")
	output := buf.String()

	if !strings.Contains(output, "TraceID=trace-999") {
		t.Errorf("Expected TraceID in output, got: %s", output)
	}
	if !strings.Contains(output, "RequestID=req-888") {
		t.Errorf("Expected RequestID in output, got: %s", output)
	}
	if !strings.Contains(output, "Tenant=tenant-A") {
		t.Errorf("Expected Tenant in output, got: %s", output)
	}
}

// TestExtractFromContextValue 测试从 context.Value 提取
func TestExtractFromContextValue(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)

	extractor := ExtractFromContextValue("custom_key", "CustomLabel")
	logger.SetContextExtractor(extractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "custom_key", "custom_value")

	logger.InfoContext(ctx, "Test message")
	output := buf.String()

	if !strings.Contains(output, "CustomLabel=custom_value") {
		t.Errorf("Expected CustomLabel in output, got: %s", output)
	}
}

// TestSetContextExtractorNil 测试设置 nil 提取器
func TestSetContextExtractorNil(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)

	// 设置 nil 应回退到默认提取器
	logger.SetContextExtractor(nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-default")
	ctx = context.WithValue(ctx, "request_id", "req-default")

	logger.InfoContext(ctx, "Test message")
	output := buf.String()

	// 应使用默认提取器
	if !strings.Contains(output, "TraceID=trace-default") {
		t.Errorf("Expected default extractor to work, got: %s", output)
	}
	if !strings.Contains(output, "RequestID=req-default") {
		t.Errorf("Expected default extractor to work, got: %s", output)
	}
}

// TestGetContextExtractor 测试获取提取器
func TestGetContextExtractor(t *testing.T) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, INFO)

	// 默认提取器
	extractor := logger.GetContextExtractor()
	if extractor == nil {
		t.Error("Default extractor should not be nil")
	}

	// 设置自定义提取器
	customExtractor := SimpleTraceIDExtractor
	logger.SetContextExtractor(customExtractor)

	// 验证获取的是设置的提取器
	retrieved := logger.GetContextExtractor()
	if retrieved == nil {
		t.Error("Retrieved extractor should not be nil")
	}
}

// TestCustomExtractorFunction 测试完全自定义的提取器函数
func TestCustomExtractorFunction(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)

	// 自定义提取器
	customExtractor := func(ctx context.Context) string {
		if ctx == nil {
			return ""
		}
		userId, _ := ctx.Value("user_id").(string)
		if userId != "" {
			return "[👤 " + userId + "] "
		}
		return ""
	}

	logger.SetContextExtractor(customExtractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", "john")

	logger.InfoContext(ctx, "User action")
	output := buf.String()

	if !strings.Contains(output, "👤 john") {
		t.Errorf("Expected custom format in output, got: %s", output)
	}
}

// TestContextExtractorWithFieldLogger 测试字段日志器与上下文提取器配合
func TestContextExtractorWithFieldLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)
	logger.SetContextExtractor(SimpleTraceIDExtractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-field-test")

	// 使用字段日志器
	fieldLogger := logger.WithField("component", "api")

	fieldLogger.InfoContext(ctx, "API call")
	output := buf.String()

	if !strings.Contains(output, "TraceID=trace-field-test") {
		t.Errorf("Expected TraceID in field logger output, got: %s", output)
	}
	if !strings.Contains(output, "component: api") {
		t.Errorf("Expected component field in output, got: %s", output)
	}
}

// TestEmptyContextExtractor 测试空上下文
func TestEmptyContextExtractor(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewUltraFastLoggerNoTime(buf, INFO)

	// nil context
	logger.InfoContext(nil, "Nil context message")
	output := buf.String()
	if !strings.Contains(output, "Nil context message") {
		t.Errorf("Expected message in output, got: %s", output)
	}

	// 空 context
	buf.Reset()
	ctx := context.Background()
	logger.InfoContext(ctx, "Empty context message")
	output = buf.String()
	if !strings.Contains(output, "Empty context message") {
		t.Errorf("Expected message in output, got: %s", output)
	}
}

// BenchmarkDefaultContextExtractor 基准测试 - 默认提取器
func BenchmarkDefaultContextExtractor(b *testing.B) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, INFO)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-123")
	ctx = context.WithValue(ctx, "request_id", "req-456")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.InfoContext(ctx, "Benchmark message")
		}
	})
}

// BenchmarkNoOpContextExtractor 基准测试 - 空操作提取器
func BenchmarkNoOpContextExtractor(b *testing.B) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, INFO)
	logger.SetContextExtractor(NoOpContextExtractor)
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.InfoContext(ctx, "Benchmark message")
		}
	})
}

// BenchmarkChainedContextExtractor 基准测试 - 链式提取器
func BenchmarkChainedContextExtractor(b *testing.B) {
	logger := NewUltraFastLoggerNoTime(&bytes.Buffer{}, INFO)
	extractor := ChainContextExtractors(
		SimpleTraceIDExtractor,
		SimpleRequestIDExtractor,
	)
	logger.SetContextExtractor(extractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-123")
	ctx = context.WithValue(ctx, "request_id", "req-456")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.InfoContext(ctx, "Benchmark message")
		}
	})
}
