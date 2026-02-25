/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-22 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 09:53:31
 * @FilePath: \go-logger\context_service_test.go
 * @Description: 上下文管理器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"testing"

	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/stretchr/testify/assert"
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
	testValue := random.UUID()

	ctx = WithValue(ctx, KeyTraceID, testValue)

	value := ctx.Value(KeyTraceID)
	assert.Equal(t, testValue, value)
}

// TestGetValue 测试从 context 获取值
func TestGetValue(t *testing.T) {
	ctx := context.Background()
	testValue := random.UUID()

	ctx = WithValue(ctx, KeyTraceID, testValue)

	value := GetValue(ctx, KeyTraceID)
	assert.Equal(t, testValue, value)
}

// TestGetString 测试获取字符串值
func TestGetString(t *testing.T) {
	ctx := context.Background()

	// 测试字符串值
	testUserID := random.UUID()
	ctx = WithValue(ctx, KeyUserID, testUserID)
	assert.Equal(t, testUserID, GetString(ctx, KeyUserID))

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

	traceID := random.UUID()
	requestID := random.UUID()
	userID := random.UUID()

	ctx = WithTraceID(ctx, traceID)
	ctx = WithRequestID(ctx, requestID)
	ctx = WithUserID(ctx, userID)

	fields := ExtractFields(ctx)

	assert.Contains(t, fields, "trace_id")
	assert.Contains(t, fields, "request_id")
	assert.Contains(t, fields, "user_id")
	assert.Equal(t, traceID, fields["trace_id"])
	assert.Equal(t, requestID, fields["request_id"])
	assert.Equal(t, userID, fields["user_id"])
}

// TestWithTraceID 测试 TraceID 操作
func TestWithTraceID(t *testing.T) {
	ctx := context.Background()
	testTraceID := random.UUID()

	ctx = WithTraceID(ctx, testTraceID)
	traceID := GetTraceID(ctx)

	assert.Equal(t, testTraceID, traceID)
}

// TestGenerateTraceID 测试生成 TraceID
func TestGenerateTraceID(t *testing.T) {
	traceID1 := GenerateTraceID()
	traceID2 := GenerateTraceID()

	assert.NotEmpty(t, traceID1)
	assert.NotEmpty(t, traceID2)
	assert.NotEqual(t, traceID1, traceID2, "每次生成的 TraceID 应该不同")
}

// TestIDMethodsWithGenerate 测试带生成函数的 ID 方法（table-driven）
func TestIDMethodsWithGenerate(t *testing.T) {
	tests := map[string]struct {
		withFunc     func(context.Context, string) context.Context
		getFunc      func(context.Context) string
		generateFunc func() string
	}{
		"SpanID":    {WithSpanID, GetSpanID, GenerateSpanID},
		"RequestID": {WithRequestID, GetRequestID, GenerateRequestID},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testID := random.UUID()
			ctx := context.Background()
			ctx = tt.withFunc(ctx, testID)
			result := tt.getFunc(ctx)
			assert.Equal(t, testID, result)

			// 测试生成函数
			generatedID := tt.generateFunc()
			assert.NotEmpty(t, generatedID)
		})
	}
}

// TestContextIDMethods 测试各种 ID 方法（table-driven）
func TestContextIDMethods(t *testing.T) {
	tests := map[string]struct {
		withFunc func(context.Context, string) context.Context
		getFunc  func(context.Context) string
	}{
		"UserID":        {WithUserID, GetUserID},
		"SessionID":     {WithSessionID, GetSessionID},
		"TenantID":      {WithTenantID, GetTenantID},
		"CorrelationID": {WithCorrelationID, GetCorrelationID},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testID := random.UUID()
			ctx := context.Background()
			ctx = tt.withFunc(ctx, testID)
			result := tt.getFunc(ctx)
			assert.Equal(t, testID, result)
		})
	}

	// 单独测试生成函数
	t.Run("GenerateCorrelationID", func(t *testing.T) {
		generatedID := GenerateCorrelationID()
		assert.NotEmpty(t, generatedID)
	})
}

// TestCreateSpan 测试创建 Span
func TestCreateSpan(t *testing.T) {
	parentCtx := context.Background()
	parentTraceID := random.UUID()
	operationName := random.UUID()
	parentCtx = WithTraceID(parentCtx, parentTraceID)

	spanCtx := CreateSpan(parentCtx, operationName)

	// 应该继承 TraceID
	assert.Equal(t, parentTraceID, GetTraceID(spanCtx))

	// 应该有新的 SpanID
	spanID := GetSpanID(spanCtx)
	assert.NotEmpty(t, spanID)

	// 应该有操作名称
	operation := spanCtx.Value(KeyOperation)
	assert.Equal(t, operationName, operation)
}

// TestCreateChildContext 测试创建子上下文
func TestCreateChildContext(t *testing.T) {
	// CreateChildContext removed in simplified service
	// Use CreateSpan instead for similar functionality
	parentCtx := context.Background()
	parentTrace := random.UUID()
	parentUser := random.UUID()
	operationName := random.UUID()
	parentCtx = WithTraceID(parentCtx, parentTrace)
	parentCtx = WithUserID(parentCtx, parentUser)

	childCtx := CreateSpan(parentCtx, operationName)

	// 应该继承 TraceID 和 UserID
	assert.Equal(t, parentTrace, GetTraceID(childCtx))
	assert.Equal(t, parentUser, GetUserID(childCtx))

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

	traceID := random.UUID()
	spanID := random.UUID()
	requestID := random.UUID()
	userID := random.UUID()
	sessionID := random.UUID()
	tenantID := random.UUID()
	correlationID := random.UUID()

	ctx = WithTraceID(ctx, traceID)
	ctx = WithSpanID(ctx, spanID)
	ctx = WithRequestID(ctx, requestID)
	ctx = WithUserID(ctx, userID)
	ctx = WithSessionID(ctx, sessionID)
	ctx = WithTenantID(ctx, tenantID)
	ctx = WithCorrelationID(ctx, correlationID)

	fields := ExtractFields(ctx)

	assert.Len(t, fields, 7)
	assert.Equal(t, traceID, fields["trace_id"])
	assert.Equal(t, spanID, fields["span_id"])
	assert.Equal(t, requestID, fields["request_id"])
	assert.Equal(t, userID, fields["user_id"])
	assert.Equal(t, sessionID, fields["session_id"])
	assert.Equal(t, tenantID, fields["tenant_id"])
	assert.Equal(t, correlationID, fields["correlation_id"])
}

// BenchmarkGetTraceID 性能测试：获取 TraceID
func BenchmarkGetTraceID(b *testing.B) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, random.UUID())

	b.ResetTimer()
	for range b.N {
		_ = GetTraceID(ctx)
	}
}

// BenchmarkExtractFields 性能测试：提取字段
func BenchmarkExtractFields(b *testing.B) {
	ctx := context.Background()
	ctx = WithTraceID(ctx, random.UUID())
	ctx = WithRequestID(ctx, random.UUID())
	ctx = WithUserID(ctx, random.UUID())

	b.ResetTimer()
	for range b.N {
		_ = ExtractFields(ctx)
	}
}

// BenchmarkGetOrGenerateTraceID 性能测试：获取或生成 TraceID
func BenchmarkGetOrGenerateTraceID(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		ctx, _ = GetOrGenerateTraceID(ctx)
	}
}
