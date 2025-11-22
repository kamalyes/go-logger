/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-22 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 09:21:28
 * @FilePath: \go-logger\context_logger.go
 * @Description: 上下文日志集成 - 统一入口设计
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
)

// ========== 统一日志入口 ==========

// ContextLogger 上下文日志门面 - 统一入口
type ContextLogger struct {
	ctx    context.Context
	logger ILogger
}

// NewContextLogger 创建上下文日志记录器（统一入口）
func NewContextLogger(ctx context.Context, logger ILogger) *ContextLogger {
	return &ContextLogger{ctx: ctx, logger: logger}
}

// Log 统一的日志记录入口
func (cl *ContextLogger) Log(level string, msg string, keysAndValues ...interface{}) {
	if cl.logger == nil {
		return
	}

	fields := ExtractFields(cl.ctx)
	logger := mathx.IF(len(fields) > 0, cl.logger.WithFields(fields), cl.logger)
	hasKV := len(keysAndValues) > 0

	switch level {
	case "trace", "debug":
		mathx.IfLazy(hasKV,
			func() any { logger.DebugKV(msg, keysAndValues...); return nil },
			func() any { logger.DebugMsg(msg); return nil },
		)
	case "info":
		mathx.IfLazy(hasKV,
			func() any { logger.InfoKV(msg, keysAndValues...); return nil },
			func() any { logger.InfoMsg(msg); return nil },
		)
	case "warn":
		mathx.IfLazy(hasKV,
			func() any { logger.WarnKV(msg, keysAndValues...); return nil },
			func() any { logger.WarnMsg(msg); return nil },
		)
	case "error":
		mathx.IfLazy(hasKV,
			func() any { logger.ErrorKV(msg, keysAndValues...); return nil },
			func() any { logger.ErrorMsg(msg); return nil },
		)
	case "fatal":
		mathx.IfLazy(hasKV,
			func() any { logger.FatalKV(msg, keysAndValues...); return nil },
			func() any { logger.FatalMsg(msg); return nil },
		)
	}
}

// Trace/Debug/Info/Warn/Error/Fatal 便捷方法
func (cl *ContextLogger) Trace(msg string, kvs ...interface{}) { cl.Log("trace", msg, kvs...) }
func (cl *ContextLogger) Debug(msg string, kvs ...interface{}) { cl.Log("debug", msg, kvs...) }
func (cl *ContextLogger) Info(msg string, kvs ...interface{})  { cl.Log("info", msg, kvs...) }
func (cl *ContextLogger) Warn(msg string, kvs ...interface{})  { cl.Log("warn", msg, kvs...) }
func (cl *ContextLogger) Error(msg string, kvs ...interface{}) { cl.Log("error", msg, kvs...) }
func (cl *ContextLogger) Fatal(msg string, kvs ...interface{}) { cl.Log("fatal", msg, kvs...) }

// ========== 向后兼容的函数式API ==========

// WithLogger 从 context 创建带所有注册字段的 logger
func WithLogger(ctx context.Context, baseLogger ILogger) ILogger {
	if ctx == nil || baseLogger == nil {
		return baseLogger
	}
	fields := ExtractFields(ctx)
	return mathx.IF(len(fields) > 0, baseLogger.WithFields(fields), baseLogger)
}

// LogWithContext 统一的日志记录入口
func LogWithContext(ctx context.Context, baseLogger ILogger, level interface{}, msg string, keysAndValues ...interface{}) {
	// 支持 string 和 LogLevel 两种类型
	var levelStr string
	switch v := level.(type) {
	case string:
		levelStr = v
	case LogLevel:
		levelStr = v.String()
	default:
		levelStr = "info"
	}
	NewContextLogger(ctx, baseLogger).Log(levelStr, msg, keysAndValues...)
}

// 便捷函数 - 所有日志级别都通过统一入口
func TraceWithContext(ctx context.Context, logger ILogger, msg string, kvs ...interface{}) {
	LogWithContext(ctx, logger, "trace", msg, kvs...)
}

func DebugWithContext(ctx context.Context, logger ILogger, msg string, kvs ...interface{}) {
	LogWithContext(ctx, logger, "debug", msg, kvs...)
}

func InfoWithContext(ctx context.Context, logger ILogger, msg string, kvs ...interface{}) {
	LogWithContext(ctx, logger, "info", msg, kvs...)
}

func WarnWithContext(ctx context.Context, logger ILogger, msg string, kvs ...interface{}) {
	LogWithContext(ctx, logger, "warn", msg, kvs...)
}

func ErrorWithContext(ctx context.Context, logger ILogger, msg string, kvs ...interface{}) {
	LogWithContext(ctx, logger, "error", msg, kvs...)
}

func FatalWithContext(ctx context.Context, logger ILogger, msg string, kvs ...interface{}) {
	LogWithContext(ctx, logger, "fatal", msg, kvs...)
}

// ========== 高级日志功能 ==========

// WithCorrelation 创建带相关性链的日志记录器
func WithCorrelation(ctx context.Context, baseLogger ILogger) (ILogger, *CorrelationChain) {
	chain, newCtx := CreateCorrelationChain(ctx)
	contextLogger := WithLogger(newCtx, baseLogger)
	return contextLogger, chain
}

// LogWithCorrelation 使用相关性链记录日志
func LogWithCorrelation(ctx context.Context, baseLogger ILogger, lvl LogLevel, msg string, keysAndValues ...interface{}) *CorrelationChain {
	chain, newCtx := CreateCorrelationChain(ctx)
	LogWithContext(newCtx, baseLogger, lvl, msg, keysAndValues...)
	return chain
}

// StartOperation 开始一个操作并记录
func StartOperation(ctx context.Context, baseLogger ILogger, operation string, keysAndValues ...interface{}) (context.Context, *CorrelationChain) {
	spanCtx := CreateSpan(ctx, operation)
	chain, chainCtx := CreateCorrelationChain(spanCtx)
	chain.SetTag("operation", operation)
	InfoWithContext(chainCtx, baseLogger, "Operation started", append(keysAndValues, "operation", operation)...)
	return chainCtx, chain
}

// EndOperation 结束一个操作并记录
func EndOperation(ctx context.Context, baseLogger ILogger, chain *CorrelationChain, keysAndValues ...interface{}) {
	if chain != nil {
		EndCorrelationChain(chain)
		keysAndValues = append(keysAndValues, "duration_ms", chain.GetDuration().Milliseconds())
	}
	InfoWithContext(ctx, baseLogger, "Operation completed", keysAndValues...)
}

// OperationLogger 操作日志记录器
type OperationLogger struct {
	ctx    context.Context
	logger ILogger
	chain  *CorrelationChain
}

// NewOperationLogger 创建操作日志记录器
func NewOperationLogger(ctx context.Context, baseLogger ILogger, operation string) *OperationLogger {
	spanCtx := CreateSpan(ctx, operation)
	chain, chainCtx := CreateCorrelationChain(spanCtx)
	chain.SetTag("operation", operation)

	ol := &OperationLogger{ctx: chainCtx, logger: baseLogger, chain: chain}
	ol.Info("Operation started")
	return ol
}

// GetContext 获取上下文
func (ol *OperationLogger) GetContext() context.Context { return ol.ctx }

// GetChain 获取相关性链
func (ol *OperationLogger) GetChain() *CorrelationChain { return ol.chain }

// SetTag 设置标签
func (ol *OperationLogger) SetTag(key, value string) *OperationLogger {
	if ol.chain != nil {
		ol.chain.SetTag(key, value)
	}
	return ol
}

// SetMetric 设置指标
func (ol *OperationLogger) SetMetric(key string, value interface{}) *OperationLogger {
	if ol.chain != nil {
		ol.chain.SetMetric(key, value)
	}
	return ol
}

// log 统一的日志记录方法
func (ol *OperationLogger) log(level string, msg string, keysAndValues ...interface{}) {
	NewContextLogger(ol.ctx, ol.logger).Log(level, msg, keysAndValues...)
}

// Trace 记录 Trace 日志
func (ol *OperationLogger) Trace(msg string, kvs ...interface{}) { ol.log("trace", msg, kvs...) }

// Debug 记录 Debug 日志
func (ol *OperationLogger) Debug(msg string, kvs ...interface{}) { ol.log("debug", msg, kvs...) }

// Info 记录 Info 日志
func (ol *OperationLogger) Info(msg string, kvs ...interface{}) { ol.log("info", msg, kvs...) }

// Warn 记录 Warn 日志
func (ol *OperationLogger) Warn(msg string, kvs ...interface{}) { ol.log("warn", msg, kvs...) }

// Error 记录 Error 日志
func (ol *OperationLogger) Error(msg string, kvs ...interface{}) { ol.log("error", msg, kvs...) }

// endOperation 统一的结束操作辅助方法
func (ol *OperationLogger) endOperation(err error, keysAndValues ...interface{}) {
	if ol.chain != nil {
		EndCorrelationChain(ol.chain)
		keysAndValues = append(keysAndValues, "duration_ms", ol.chain.GetDuration().Milliseconds())
		if err != nil {
			ol.chain.SetTag("error", "true")
			ol.chain.SetTag("error_message", err.Error())
		}
	}

	if err != nil {
		keysAndValues = append(keysAndValues, "error", err.Error())
		ol.Error("Operation failed", keysAndValues...)
	} else {
		ol.Info("Operation completed", keysAndValues...)
	}
}

// End 结束操作
func (ol *OperationLogger) End(keysAndValues ...interface{}) {
	ol.endOperation(nil, keysAndValues...)
}

// EndWithError 结束操作（带错误）
func (ol *OperationLogger) EndWithError(err error, keysAndValues ...interface{}) {
	ol.endOperation(err, keysAndValues...)
}
