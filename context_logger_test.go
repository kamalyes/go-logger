/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-22 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 10:25:40
 * @FilePath: \go-logger\context_logger_test.go
 * @Description: 上下文日志集成测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

// createTestLogger 创建测试用 logger
func createTestLogger() *Logger {
	config := &LogConfig{
		Level: DEBUG,
	}
	return NewLogger(config)
}

// TestWithLogger 测试创建带上下文字段的 logger
func TestWithLogger(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background()

	// 添加上下文字段
	ctx = WithTraceID(ctx, "test-trace-123")
	ctx = WithUserID(ctx, "test-user-456")

	// 创建带字段的 logger
	contextLogger := WithLogger(ctx, baseLogger)

	assert.NotNil(t, contextLogger)
	assert.NotEqual(t, baseLogger, contextLogger, "应该返回新的 logger 实例")
}

// TestWithLoggerNilInputs 测试 nil 输入
func TestWithLoggerNilInputs(t *testing.T) {
	baseLogger := createTestLogger()

	// nil context
	result := WithLogger(nil, baseLogger)
	assert.Equal(t, baseLogger, result)

	// nil logger
	ctx := context.Background()
	result = WithLogger(ctx, nil)
	assert.Nil(t, result)
}

// TestWithLoggerEmptyContext 测试空上下文
func TestWithLoggerEmptyContext(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background()

	// 空上下文应该返回原 logger
	contextLogger := WithLogger(ctx, baseLogger)
	assert.Equal(t, baseLogger, contextLogger)
}

// TestLogWithContext 测试带上下文的日志记录
func TestLogWithContext(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background()

	// 添加上下文字段
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithRequestID(ctx, "req-456")

	// 使用 LogWithContext
	LogWithContext(ctx, baseLogger, INFO, "test message")

	// 基本验证：没有 panic 即成功
	assert.NotNil(t, ctx)
}

// TestTraceWithContext 测试 Trace 级别日志
func TestTraceWithContext(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 使用 TraceWithContext
	TraceWithContext(ctx, baseLogger, "trace message")

	// 基本验证：没有 panic 即成功
	assert.NotNil(t, ctx)
}

// TestDebugWithContext 测试 Debug 级别日志
func TestDebugWithContext(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 使用 DebugWithContext
	DebugWithContext(ctx, baseLogger, "debug message")

	// 基本验证：没有 panic 即成功
	assert.NotNil(t, ctx)
}

// TestInfoWithContext 测试 Info 级别日志
func TestInfoWithContext(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 使用 InfoWithContext
	InfoWithContext(ctx, baseLogger, "info message")

	// 基本验证：没有 panic 即成功
	assert.NotNil(t, ctx)
}

// TestWarnWithContext 测试 Warn 级别日志
func TestWarnWithContext(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 使用 WarnWithContext
	WarnWithContext(ctx, baseLogger, "warn message")

	// 基本验证：没有 panic 即成功
	assert.NotNil(t, ctx)
}

// TestErrorWithContext 测试 Error 级别日志
func TestErrorWithContext(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 使用 ErrorWithContext
	ErrorWithContext(ctx, baseLogger, "error message")

	// 基本验证：没有 panic 即成功
	assert.NotNil(t, ctx)
}

// TestFatalWithContext 测试 Fatal 级别日志
func TestFatalWithContext(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")

	// 注意：FatalWithContext 会调用 os.Exit(1)，所以在测试中跳过实际调用
	// 只验证函数存在且不会 panic（如果 logger 为 nil）
	assert.NotPanics(t, func() {
		FatalWithContext(ctx, nil, "fatal message")
	})

	// 基本验证：没有 panic 即成功
	assert.NotNil(t, ctx)
}

// TestOperationLogger 测试操作日志记录器
func TestOperationLogger(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 创建操作日志记录器
	ctx, chain := StartOperation(ctx, baseLogger, "test-operation")

	assert.NotNil(t, ctx)
	assert.NotNil(t, chain)

	// 结束操作
	EndOperation(ctx, baseLogger, chain)
}

// TestOperationLoggerSuccess 测试成功操作日志
func TestOperationLoggerSuccess(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 模拟成功操作
	ctx, chain := StartOperation(ctx, baseLogger, "test-operation")
	InfoWithContext(ctx, baseLogger, "operation completed successfully")
	EndOperation(ctx, baseLogger, chain)

	assert.NotNil(t, chain)
}

// TestOperationLoggerFailure 测试失败操作日志
func TestOperationLoggerFailure(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 模拟失败操作
	ctx, chain := StartOperation(ctx, baseLogger, "test-operation")
	ErrorWithContext(ctx, baseLogger, "operation failed")
	EndOperation(ctx, baseLogger, chain, "error", "some error")

	assert.NotNil(t, chain)
}

// TestLogWithCorrelation 测试关联链日志记录
func TestLogWithCorrelation(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")

	// 使用 LogWithCorrelation
	chain := LogWithCorrelation(ctx, baseLogger, INFO, "test message", "key", "value")

	// 基本验证：没有 panic 即成功
	assert.NotNil(t, chain)
	assert.Equal(t, "trace-123", chain.TraceID)
}

// TestConcurrentLogging 测试并发日志记录
func TestConcurrentLogging(t *testing.T) {
	baseLogger := createTestLogger()

	var wg sync.WaitGroup
	concurrency := 100

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			ctx := context.Background()
			ctx = WithTraceID(ctx, "trace-"+string(rune(id)))
			ctx = WithUserID(ctx, "user-"+string(rune(id)))

			InfoWithContext(ctx, baseLogger, "concurrent log message")
		}(i)
	}

	wg.Wait()

	// 基本验证：没有 panic 即成功
	assert.NotNil(t, baseLogger)
}

// BenchmarkWithLogger 性能测试：创建带上下文的 logger
func BenchmarkWithLogger(b *testing.B) {
	baseLogger := createTestLogger()
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithUserID(ctx, "user-456")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithLogger(ctx, baseLogger)
	}
}

// BenchmarkLogWithContext 性能测试:带上下文的日志记录
func BenchmarkLogWithContext(b *testing.B) {
	baseLogger := createTestLogger()
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithUserID(ctx, "user-456")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		InfoWithContext(ctx, baseLogger, "benchmark message")
	}
}

// ========== 完整覆盖测试套件 ==========

// TestWithLogger_AllFields 测试所有上下文字段类型
func TestWithLogger_AllFields(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background()

	// 添加所有7种上下文字段
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithSpanID(ctx, "span-456")
	ctx = WithRequestID(ctx, "req-789")
	ctx = WithUserID(ctx, "user-abc")
	ctx = WithSessionID(ctx, "session-def")
	ctx = WithTenantID(ctx, "tenant-ghi")
	ctx = WithCorrelationID(ctx, "corr-jkl")

	contextLogger := WithLogger(ctx, baseLogger)

	assert.NotNil(t, contextLogger)
	assert.NotEqual(t, baseLogger, contextLogger)
}

// TestWithLogger_SingleField 测试单个字段
func TestWithLogger_SingleField(t *testing.T) {
	baseLogger := createTestLogger()

	tests := []struct {
		name     string
		setupCtx func(context.Context) context.Context
	}{
		{"TraceID", func(ctx context.Context) context.Context { return WithTraceID(ctx, "t1") }},
		{"SpanID", func(ctx context.Context) context.Context { return WithSpanID(ctx, "s1") }},
		{"RequestID", func(ctx context.Context) context.Context { return WithRequestID(ctx, "r1") }},
		{"UserID", func(ctx context.Context) context.Context { return WithUserID(ctx, "u1") }},
		{"SessionID", func(ctx context.Context) context.Context { return WithSessionID(ctx, "ss1") }},
		{"TenantID", func(ctx context.Context) context.Context { return WithTenantID(ctx, "tn1") }},
		{"CorrelationID", func(ctx context.Context) context.Context { return WithCorrelationID(ctx, "c1") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupCtx(context.Background())
			contextLogger := WithLogger(ctx, baseLogger)
			assert.NotNil(t, contextLogger)
			assert.NotEqual(t, baseLogger, contextLogger)
		})
	}
}

// TestLogWithContext_WithKeysAndValues 测试带额外键值对
func TestLogWithContext_WithKeysAndValues(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 单个键值对
	LogWithContext(ctx, baseLogger, INFO, "message1", "key1", "value1")

	// 多个键值对
	LogWithContext(ctx, baseLogger, INFO, "message2", "key1", "value1", "key2", "value2")

	// 奇数个参数(最后一个会被忽略)
	LogWithContext(ctx, baseLogger, INFO, "message3", "key1", "value1", "odd")

	assert.NotNil(t, ctx)
}

// TestLogWithContext_NilLogger 测试nil logger
func TestLogWithContext_NilLogger(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")

	// nil logger应该不panic
	assert.NotPanics(t, func() {
		LogWithContext(ctx, nil, INFO, "message")
	})
}

// TestLogWithContext_EmptyFields 测试空字段
func TestLogWithContext_EmptyFields(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background() // 空context,没有字段

	// 空字段应该正常工作
	LogWithContext(ctx, baseLogger, INFO, "message")

	assert.NotNil(t, ctx)
}

// TestAllLogLevels_WithKeysAndValues 测试所有日志级别带键值对
func TestAllLogLevels_WithKeysAndValues(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	tests := []struct {
		name    string
		logFunc func(context.Context, ILogger, string, ...interface{})
		message string
	}{
		{"Trace", TraceWithContext, "trace msg"},
		{"Debug", DebugWithContext, "debug msg"},
		{"Info", InfoWithContext, "info msg"},
		{"Warn", WarnWithContext, "warn msg"},
		{"Error", ErrorWithContext, "error msg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 无额外参数
			tt.logFunc(ctx, baseLogger, tt.message)

			// 带键值对
			tt.logFunc(ctx, baseLogger, tt.message, "key1", "val1")
			tt.logFunc(ctx, baseLogger, tt.message, "key1", "val1", "key2", "val2")
		})
	}

	// Fatal 需要单独测试（使用 nil logger 避免退出）
	t.Run("Fatal", func(t *testing.T) {
		assert.NotPanics(t, func() {
			FatalWithContext(ctx, nil, "fatal msg")
			FatalWithContext(ctx, nil, "fatal msg", "key1", "val1")
		})
	})
}

// TestAllLogLevels_NilLogger 测试所有日志级别nil logger
func TestAllLogLevels_NilLogger(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")

	tests := []struct {
		name    string
		logFunc func(context.Context, ILogger, string, ...interface{})
	}{
		{"Trace", TraceWithContext},
		{"Debug", DebugWithContext},
		{"Info", InfoWithContext},
		{"Warn", WarnWithContext},
		{"Error", ErrorWithContext},
		{"Fatal", FatalWithContext},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				tt.logFunc(ctx, nil, "message")
			})
		})
	}
}

// TestWithCorrelation_NilLogger 测试nil logger
func TestWithCorrelation_NilLogger(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")

	contextLogger, chain := WithCorrelation(ctx, nil)

	assert.Nil(t, contextLogger)
	assert.NotNil(t, chain)
	assert.Equal(t, "trace-123", chain.TraceID)
}

// TestWithCorrelation_EmptyContext 测试空上下文
func TestWithCorrelation_EmptyContext(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background()

	contextLogger, chain := WithCorrelation(ctx, baseLogger)

	assert.NotNil(t, chain)
	assert.NotEmpty(t, chain.TraceID) // CreateCorrelationChain会自动生成TraceID
	assert.NotNil(t, contextLogger)
}

// TestWithCorrelation_WithParentChain 测试父链继承
func TestWithCorrelation_WithParentChain(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")

	// 创建父链
	parentLogger, parentChain := WithCorrelation(ctx, baseLogger)
	assert.NotNil(t, parentLogger)
	assert.NotNil(t, parentChain)

	// 在相同trace下创建子链
	childLogger, childChain := WithCorrelation(ctx, baseLogger)

	assert.NotNil(t, childChain)
	assert.Equal(t, "trace-123", childChain.TraceID)
	assert.NotNil(t, childLogger)
}

// TestLogWithCorrelation_AllParameters 测试所有参数组合
func TestLogWithCorrelation_AllParameters(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	// 无额外参数
	chain1 := LogWithCorrelation(ctx, baseLogger, INFO, "msg1")
	assert.NotNil(t, chain1)

	// 带键值对
	chain2 := LogWithCorrelation(ctx, baseLogger, INFO, "msg2", "k1", "v1")
	assert.NotNil(t, chain2)

	// 多个键值对
	chain3 := LogWithCorrelation(ctx, baseLogger, INFO, "msg3", "k1", "v1", "k2", "v2")
	assert.NotNil(t, chain3)
}

// TestLogWithCorrelation_NilLogger 测试nil logger
func TestLogWithCorrelation_NilLogger(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")

	chain := LogWithCorrelation(ctx, nil, INFO, "message")

	assert.NotNil(t, chain)
	assert.Equal(t, "trace-123", chain.TraceID)
}

// TestStartOperation_Basic 测试基本操作启动
func TestStartOperation_Basic(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	newCtx, chain := StartOperation(ctx, baseLogger, "test-op")

	assert.NotNil(t, newCtx)
	assert.NotNil(t, chain)
	assert.Equal(t, "test-op", chain.Tags["operation"])
	assert.Equal(t, "trace-123", chain.TraceID)
}

// TestStartOperation_NilLogger 测试nil logger
func TestStartOperation_NilLogger(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")

	newCtx, chain := StartOperation(ctx, nil, "test-op")

	assert.NotNil(t, newCtx)
	assert.NotNil(t, chain)
}

// TestEndOperation_Basic 测试基本操作结束
func TestEndOperation_Basic(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")

	// 无额外参数
	EndOperation(ctx, baseLogger, chain)

	// 带键值对
	EndOperation(ctx, baseLogger, chain, "key", "value")
}

// TestEndOperation_NilLogger 测试nil logger
func TestEndOperation_NilLogger(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, nil, "test-op")

	assert.NotPanics(t, func() {
		EndOperation(ctx, nil, chain)
	})
}

// TestEndOperation_NilChain 测试nil chain
func TestEndOperation_NilChain(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background()

	assert.NotPanics(t, func() {
		EndOperation(ctx, baseLogger, nil)
	})
}

// ========== OperationLogger 完整测试 ==========

// TestOperationLogger_AllMethods 测试所有方法
func TestOperationLogger_AllMethods(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")
	opLogger := &OperationLogger{
		ctx:    ctx,
		logger: baseLogger,
		chain:  chain,
	}

	// SetTag
	opLogger.SetTag("tag1", "value1")
	opLogger.SetTag("tag2", "123")

	// SetMetric
	opLogger.SetMetric("metric1", 100)
	opLogger.SetMetric("metric2", 200.5)

	// 日志方法
	opLogger.Trace("trace msg")
	opLogger.Debug("debug msg")
	opLogger.Info("info msg")
	opLogger.Warn("warn msg")
	opLogger.Error("error msg")

	// End
	opLogger.End()
}

// TestOperationLogger_EndWithError 测试带错误结束
func TestOperationLogger_EndWithError(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")
	opLogger := &OperationLogger{
		ctx:    ctx,
		logger: baseLogger,
		chain:  chain,
	}

	// EndWithError
	opLogger.EndWithError(assert.AnError)
}

// TestOperationLogger_NilLogger 测试nil logger
func TestOperationLogger_NilLogger(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-123")
	chain := &CorrelationChain{
		TraceID: "trace-123",
		Tags:    make(map[string]string),
		Metrics: make(map[string]interface{}),
	}

	opLogger := &OperationLogger{
		ctx:    ctx,
		logger: nil,
		chain:  chain,
	}

	// 所有方法都应该不panic
	assert.NotPanics(t, func() {
		opLogger.SetTag("tag", "value")
		opLogger.SetMetric("metric", 100)
		opLogger.Trace("trace")
		opLogger.Debug("debug")
		opLogger.Info("info")
		opLogger.Warn("warn")
		opLogger.Error("error")
		opLogger.End()
		opLogger.EndWithError(assert.AnError)
	})
}

// TestOperationLogger_WithKeysAndValues 测试所有日志方法带键值对
func TestOperationLogger_WithKeysAndValues(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")
	opLogger := &OperationLogger{
		ctx:    ctx,
		logger: baseLogger,
		chain:  chain,
	}

	// 测试所有日志级别带键值对
	opLogger.Trace("trace msg", "k1", "v1")
	opLogger.Debug("debug msg", "k1", "v1", "k2", "v2")
	opLogger.Info("info msg", "k1", "v1")
	opLogger.Warn("warn msg", "k1", "v1", "k2", "v2")
	opLogger.Error("error msg", "k1", "v1")

	opLogger.End("result", "success")
}

// TestOperationLogger_SetTagMultipleTypes 测试不同类型的tag
func TestOperationLogger_SetTagMultipleTypes(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")
	opLogger := &OperationLogger{
		ctx:    ctx,
		logger: baseLogger,
		chain:  chain,
	}

	// 不同类型的值
	opLogger.SetTag("string", "value")
	opLogger.SetTag("int", "123")
	opLogger.SetTag("float", "123.45")
	opLogger.SetTag("bool", "true")
	opLogger.SetTag("nil", "")

	opLogger.End()
}

// TestOperationLogger_SetMetricMultipleTypes 测试不同类型的metric
func TestOperationLogger_SetMetricMultipleTypes(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")
	opLogger := &OperationLogger{
		ctx:    ctx,
		logger: baseLogger,
		chain:  chain,
	}

	// 不同类型的值
	opLogger.SetMetric("int", 100)
	opLogger.SetMetric("float", 100.5)
	opLogger.SetMetric("negative", -50)
	opLogger.SetMetric("zero", 0)

	opLogger.End()
}

// TestOperationLogger_ChainModification 测试链的修改
func TestOperationLogger_ChainModification(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")
	opLogger := &OperationLogger{
		ctx:    ctx,
		logger: baseLogger,
		chain:  chain,
	}

	// 验证chain初始状态（StartOperation会自动添加operation tag）
	assert.Len(t, chain.Tags, 1) // operation tag
	assert.Empty(t, chain.Metrics)

	// 添加数据
	opLogger.SetTag("tag1", "value1")
	assert.Len(t, chain.Tags, 2) // operation + tag1

	opLogger.SetMetric("metric1", 100)
	assert.Len(t, chain.Metrics, 1)

	opLogger.End()
}

// TestOperationLogger_ContextPropagation 测试context传播
func TestOperationLogger_ContextPropagation(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")
	opLogger := &OperationLogger{
		ctx:    ctx,
		logger: baseLogger,
		chain:  chain,
	}

	// 验证context传播
	assert.Equal(t, ctx, opLogger.ctx)
	assert.NotNil(t, opLogger.ctx)
	assert.Equal(t, chain, opLogger.chain)
}

// TestOperationLogger_LoggerReplacement 测试logger替换
func TestOperationLogger_LoggerReplacement(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")
	opLogger := &OperationLogger{
		ctx:    ctx,
		logger: baseLogger,
		chain:  chain,
	}

	// 验证logger设置
	assert.Equal(t, baseLogger, opLogger.logger)
	assert.Equal(t, ctx, opLogger.ctx)
	assert.Equal(t, chain, opLogger.chain)
}

// TestWithLogger_MultipleFieldsSameType 测试多个相同类型字段
func TestWithLogger_MultipleFieldsSameType(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := context.Background()

	// 多次设置同一字段(应该使用最后一个)
	ctx = WithTraceID(ctx, "trace-1")
	ctx = WithTraceID(ctx, "trace-2")
	ctx = WithTraceID(ctx, "trace-3")

	contextLogger := WithLogger(ctx, baseLogger)
	assert.NotNil(t, contextLogger)
}

// TestEndOperation_WithMultipleKeysAndValues 测试多种键值对组合
func TestEndOperation_WithMultipleKeysAndValues(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	ctx, chain := StartOperation(ctx, baseLogger, "test-op")

	// 多种键值对
	EndOperation(ctx, baseLogger, chain,
		"key1", "value1",
		"key2", 123,
		"key3", true,
		"key4", 123.45,
	)
}

// TestLogWithCorrelation_AllLogLevels 测试所有日志级别
func TestLogWithCorrelation_AllLogLevels(t *testing.T) {
	baseLogger := createTestLogger()
	ctx := WithTraceID(context.Background(), "trace-123")

	levels := []LogLevel{
		DEBUG,
		INFO,
		WARN,
		ERROR,
	}

	for _, level := range levels {
		t.Run(level.String(), func(t *testing.T) {
			chain := LogWithCorrelation(ctx, baseLogger, level, "test message")
			assert.NotNil(t, chain)
			assert.Equal(t, "trace-123", chain.TraceID)
		})
	}

	// Fatal 需要单独测试（使用 nil logger 避免退出）
	t.Run("FATAL", func(t *testing.T) {
		chain := LogWithCorrelation(ctx, nil, FATAL, "test message")
		assert.NotNil(t, chain)
		assert.Equal(t, "trace-123", chain.TraceID)
	})
}
