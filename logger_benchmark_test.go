/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-07-29 10:38:52
 * @FilePath: \go-logger\logger_benchmark_test.go
 * @Description: 日志性能基准测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"io"
	"testing"
)

// 创建输出到 io.Discard 的 logger，消除 I/O 开销，专注测量 logger 自身开销
func benchLogger(format FormatType, showCaller bool) *Logger {
	l := &Logger{
		colorful:     false,
		format:       format,
		timeFormat:   "2006/01/02 15:04:05",
		output:       io.Discard,
		timestampKey: "timestamp",
		levelKey:     "level",
		messageKey:   "message",
		callerKey:    "caller",
		contextKeys:  append([]compiledContextKey(nil), defaultCompiledContextKeys...),
	}
	l.level.Store(int32(DEBUG))
	l.showCaller.Store(showCaller)
	return l
}

// ==================== JSON 模式（生产默认） ====================

func BenchmarkJSON_Info_NoCaller(b *testing.B) {
	l := benchLogger(FormatJSON, false)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("用户登录成功 user_id=%d ip=%s", 12345, "192.168.1.1")
	}
}

func BenchmarkJSON_Info_WithCaller(b *testing.B) {
	l := benchLogger(FormatJSON, true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("用户登录成功 user_id=%d ip=%s", 12345, "192.168.1.1")
	}
}

func BenchmarkJSON_InfoMsg_WithCaller(b *testing.B) {
	l := benchLogger(FormatJSON, true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.InfoMsg("服务启动完成")
	}
}

func BenchmarkJSON_InfoKV_WithCaller(b *testing.B) {
	l := benchLogger(FormatJSON, true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.InfoKV("数据库查询", "table", "users", "rows", 42, "cost_ms", 3)
	}
}

func BenchmarkJSON_InfoContext_WithCaller(b *testing.B) {
	l := benchLogger(FormatJSON, true)
	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.InfoContext(ctx, "处理请求 method=%s path=%s", "GET", "/api/v1/users")
	}
}

func BenchmarkJSON_Debug_Disabled(b *testing.B) {
	l := benchLogger(FormatJSON, true)
	l.SetLevel(INFO) // DEBUG 被禁用，测试快速路径
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Debug("这条日志不会输出 %d", 42)
	}
}

// ==================== Text 模式 ====================

func BenchmarkText_Info_WithCaller(b *testing.B) {
	l := benchLogger(FormatText, true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.Info("用户登录成功 user_id=%d ip=%s", 12345, "192.168.1.1")
	}
}

func BenchmarkText_InfoKV_WithCaller(b *testing.B) {
	l := benchLogger(FormatText, true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.InfoKV("数据库查询", "table", "users", "rows", 42, "cost_ms", 3)
	}
}

// ==================== 并发场景 ====================

func BenchmarkJSON_Info_Parallel(b *testing.B) {
	l := benchLogger(FormatJSON, true)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("并发日志 user_id=%d", 12345)
		}
	})
}

// ==================== Caller 开销隔离 ====================

func BenchmarkFindExternalCaller(b *testing.B) {
	l := benchLogger(FormatJSON, true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.findExternalCaller()
	}
}

// ==================== WithFields 场景 ====================

func BenchmarkJSON_WithFields(b *testing.B) {
	l := benchLogger(FormatJSON, true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l.WithField("request_id", "abc-123").WithField("user_id", 42).Info("处理完成")
	}
}
