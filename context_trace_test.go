/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-22 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 09:53:31
 * @FilePath: \go-logger\context_trace_test.go
 * @Description: Trace 上下文管理器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetTraceManager 测试获取 TraceManager
// 已移除 TraceManager, 改为测试默认服务基本功能
func TestContextServiceBasic(t *testing.T) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-abc")
	ctx = WithSpanID(ctx, "span-xyz")
	fields := ExtractFields(ctx)
	assert.Equal(t, "trace-abc", fields["trace_id"])
	assert.Equal(t, "span-xyz", fields["span_id"])
}

// TestWithValue 测试添加值到 context
func TestWithValue(t *testing.T) {
	ctx := context.Background()
	testValue := "test-value"

	ctx = WithValue(ctx, KeyTraceID, testValue)

	value := ctx.Value(KeyTraceID)
	assert.Equal(t, testValue, value)
}

// TestGetValue 测试从 context 获取值
func TestGetValue(t *testing.T) {
	ctx := context.Background()
	testValue := "test-trace-id"

	ctx = WithValue(ctx, KeyTraceID, testValue)

	value := GetValue(ctx, KeyTraceID)
	assert.Equal(t, testValue, value)
}

// TestGetString 测试获取字符串值
func TestGetString(t *testing.T) {
	ctx := context.Background()

	// 测试字符串值
	ctx = WithValue(ctx, KeyUserID, "user-123")
	assert.Equal(t, "user-123", GetString(ctx, KeyUserID))

	// 测试非字符串值
	ctx = WithValue(ctx, KeyUserID, 12345)
	assert.Equal(t, "", GetString(ctx, KeyUserID))

	// 测试不存在的键
	assert.Equal(t, "", GetString(ctx, ContextKey("nonexistent")))
}

// TestGetOrGenerate 测试获取或生成值
func TestGetOrGenerateTraceID_New(t *testing.T) {
	ctx := context.Background()
	newCtx, id := GetOrGenerateTraceID(ctx)
	assert.NotEmpty(t, id)
	assert.NotNil(t, newCtx)
	// 再次调用应返回同一值
	newCtx2, id2 := GetOrGenerateTraceID(newCtx)
	assert.Equal(t, id, id2)
	assert.Equal(t, newCtx, newCtx2)
}

// TestExtractFields 测试提取所有字段
func TestExtractFields(t *testing.T) {
	ctx := context.Background()

	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithRequestID(ctx, "request-456")
	ctx = WithUserID(ctx, "user-789")

	fields := ExtractFields(ctx)

	assert.Contains(t, fields, "trace_id")
	assert.Contains(t, fields, "request_id")
	assert.Contains(t, fields, "user_id")
	assert.Equal(t, "trace-123", fields["trace_id"])
	assert.Equal(t, "request-456", fields["request_id"])
	assert.Equal(t, "user-789", fields["user_id"])
}

// TestWithTraceID 测试 TraceID 操作
func TestWithTraceID(t *testing.T) {
	ctx := context.Background()
	testTraceID := "trace-test-123"

	ctx = WithTraceID(ctx, testTraceID)
	traceID := GetTraceID(ctx)

	assert.Equal(t, testTraceID, traceID)
}

// TestGetOrGenerateTraceID 测试获取或生成 TraceID
func TestGetOrGenerateTraceID(t *testing.T) {
	// 测试生成新的 TraceID
	ctx := context.Background()
	newCtx, traceID := GetOrGenerateTraceID(ctx)

	assert.NotEmpty(t, traceID)
	assert.NotNil(t, newCtx)

	// 测试获取已存在的 TraceID
	existingID := "existing-trace-id"
	ctx = WithTraceID(context.Background(), existingID)
	newCtx, traceID = GetOrGenerateTraceID(ctx)

	assert.Equal(t, existingID, traceID)
}

// TestGenerateTraceID 测试生成 TraceID
func TestGenerateTraceID(t *testing.T) {
	traceID1 := GenerateTraceID()
	traceID2 := GenerateTraceID()

	assert.NotEmpty(t, traceID1)
	assert.NotEmpty(t, traceID2)
	assert.NotEqual(t, traceID1, traceID2, "每次生成的 TraceID 应该不同")
}

// TestSpanIDMethods 测试 SpanID 方法
func TestSpanIDMethods(t *testing.T) {
	ctx := context.Background()
	testSpanID := "span-test-456"

	ctx = WithSpanID(ctx, testSpanID)
	spanID := GetSpanID(ctx)

	assert.Equal(t, testSpanID, spanID)

	// 测试生成
	generatedID := GenerateSpanID()
	assert.NotEmpty(t, generatedID)
}

// TestRequestIDMethods 测试 RequestID 方法
func TestRequestIDMethods(t *testing.T) {
	ctx := context.Background()
	testRequestID := "request-test-789"

	ctx = WithRequestID(ctx, testRequestID)
	requestID := GetRequestID(ctx)

	assert.Equal(t, testRequestID, requestID)

	// 测试生成
	generatedID := GenerateRequestID()
	assert.NotEmpty(t, generatedID)
}

// TestUserIDMethods 测试 UserID 方法
func TestUserIDMethods(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-test-001"

	ctx = WithUserID(ctx, testUserID)
	userID := GetUserID(ctx)

	assert.Equal(t, testUserID, userID)
}

// TestSessionIDMethods 测试 SessionID 方法
func TestSessionIDMethods(t *testing.T) {
	ctx := context.Background()
	testSessionID := "session-test-002"

	ctx = WithSessionID(ctx, testSessionID)
	sessionID := GetSessionID(ctx)

	assert.Equal(t, testSessionID, sessionID)
}

// TestTenantIDMethods 测试 TenantID 方法
func TestTenantIDMethods(t *testing.T) {
	ctx := context.Background()
	testTenantID := "tenant-test-003"

	ctx = WithTenantID(ctx, testTenantID)
	tenantID := GetTenantID(ctx)

	assert.Equal(t, testTenantID, tenantID)
}

// TestCorrelationIDMethods 测试 CorrelationID 方法
func TestCorrelationIDMethods(t *testing.T) {
	ctx := context.Background()
	testCorrelationID := "correlation-test-004"

	ctx = WithCorrelationID(ctx, testCorrelationID)
	correlationID := GetCorrelationID(ctx)

	assert.Equal(t, testCorrelationID, correlationID)

	// 测试生成
	generatedID := GenerateCorrelationID()
	assert.NotEmpty(t, generatedID)
}

// TestCreateSpan 测试创建 Span
func TestCreateSpan(t *testing.T) {
	parentCtx := context.Background()
	parentCtx = WithTraceID(parentCtx, "parent-trace-id")

	spanCtx := CreateSpan(parentCtx, "test-operation")

	// 应该继承 TraceID
	assert.Equal(t, "parent-trace-id", GetTraceID(spanCtx))

	// 应该有新的 SpanID
	spanID := GetSpanID(spanCtx)
	assert.NotEmpty(t, spanID)

	// 应该有操作名称
	operation := spanCtx.Value(KeyOperation)
	assert.Equal(t, "test-operation", operation)
}

// TestCreateChildContext 测试创建子上下文
func TestCreateChildContext(t *testing.T) {
	// CreateChildContext removed in simplified service
	// Use CreateSpan instead for similar functionality
	parentCtx := context.Background()
	parentCtx = WithTraceID(parentCtx, "parent-trace")
	parentCtx = WithUserID(parentCtx, "parent-user")

	childCtx := CreateSpan(parentCtx, "test-operation")

	// 应该继承 TraceID 和 UserID
	assert.Equal(t, "parent-trace", GetTraceID(childCtx))
	assert.Equal(t, "parent-user", GetUserID(childCtx))

	// 应该有新的 SpanID
	childSpanID := GetSpanID(childCtx)
	assert.NotEmpty(t, childSpanID)
}

// TestCreateCorrelationChain 测试创建相关性链
func TestCreateCorrelationChain_New(t *testing.T) {
	ctx := context.Background()
	chain, newCtx := CreateCorrelationChain(ctx)
	assert.NotNil(t, chain)
	assert.NotEmpty(t, chain.ID)
	assert.NotEmpty(t, chain.TraceID)
	assert.Equal(t, chain.ID, GetCorrelationID(newCtx))
}

// TestMultipleFields 测试多个字段组合
func TestMultipleFields(t *testing.T) {
	ctx := context.Background()

	ctx = WithTraceID(ctx, "trace-multi")
	ctx = WithSpanID(ctx, "span-multi")
	ctx = WithRequestID(ctx, "request-multi")
	ctx = WithUserID(ctx, "user-multi")
	ctx = WithSessionID(ctx, "session-multi")
	ctx = WithTenantID(ctx, "tenant-multi")
	ctx = WithCorrelationID(ctx, "correlation-multi")

	fields := ExtractFields(ctx)

	assert.Len(t, fields, 7)
	assert.Equal(t, "trace-multi", fields["trace_id"])
	assert.Equal(t, "span-multi", fields["span_id"])
	assert.Equal(t, "request-multi", fields["request_id"])
	assert.Equal(t, "user-multi", fields["user_id"])
	assert.Equal(t, "session-multi", fields["session_id"])
	assert.Equal(t, "tenant-multi", fields["tenant_id"])
	assert.Equal(t, "correlation-multi", fields["correlation_id"])
}

// BenchmarkGetTraceID 性能测试：获取 TraceID
func BenchmarkGetTraceID(b *testing.B) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "benchmark-trace-id")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetTraceID(ctx)
	}
}

// BenchmarkExtractFields 性能测试：提取字段
func BenchmarkExtractFields(b *testing.B) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithRequestID(ctx, "request-456")
	ctx = WithUserID(ctx, "user-789")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ExtractFields(ctx)
	}
}

// BenchmarkGetOrGenerateTraceID 性能测试：获取或生成 TraceID
func BenchmarkGetOrGenerateTraceID(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, _ = GetOrGenerateTraceID(ctx)
	}
}
