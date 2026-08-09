/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 09:42:37
 * @FilePath: \go-logger\async_batch_test.go
 * @Description: 测试 AsyncBatchWriter
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestAsyncBatchWriter_BasicWrite 测试基本写入功能
func TestAsyncBatchWriter_BasicWrite(t *testing.T) {
	var buf bytes.Buffer
	underlying := NewConsoleWriter(WithConsoleOutput(&buf))
	w := NewAsyncBatchWriter(
		WithAsyncUnderlying(underlying),
		WithAsyncBatchSize(10),
		WithAsyncFlushInterval(50*time.Millisecond),
	)

	w.Write([]byte("hello world\n"))

	// Close 会 drain 所有条目并停止 goroutine，确保后续读取安全
	w.Close()

	output := buf.String()
	if !strings.Contains(output, "hello world") {
		t.Errorf("expected output to contain 'hello world', got: %s", output)
	}
}

// TestAsyncBatchWriter_Flush 测试手动 flush
func TestAsyncBatchWriter_Flush(t *testing.T) {
	var buf bytes.Buffer
	underlying := NewConsoleWriter(WithConsoleOutput(&buf))
	w := NewAsyncBatchWriter(
		WithAsyncUnderlying(underlying),
		WithAsyncBatchSize(1000),               // 大阈值，确保不会因 batch 满而 flush
		WithAsyncFlushInterval(10*time.Second), // 长间隔，确保只能手动 flush
	)

	w.Write([]byte("flush test\n"))

	// 立即 flush，不等定时器
	w.Flush()
	// Close 停止 goroutine 后再读取 buf
	w.Close()

	output := buf.String()
	if !strings.Contains(output, "flush test") {
		t.Errorf("expected output to contain 'flush test' after Flush(), got: %s", output)
	}
}

// TestAsyncBatchWriter_Close 测试关闭时 flush 剩余日志
func TestAsyncBatchWriter_Close(t *testing.T) {
	var buf bytes.Buffer
	underlying := NewConsoleWriter(WithConsoleOutput(&buf))
	w := NewAsyncBatchWriter(
		WithAsyncUnderlying(underlying),
		WithAsyncBatchSize(1000),
		WithAsyncFlushInterval(10*time.Second),
	)

	w.Write([]byte("close test\n"))
	w.Close()

	output := buf.String()
	if !strings.Contains(output, "close test") {
		t.Errorf("expected output to contain 'close test' after Close(), got: %s", output)
	}
}

// TestAsyncBatchWriter_ChannelFull 测试 channel 满时降级为同步写入
func TestAsyncBatchWriter_ChannelFull(t *testing.T) {
	var buf bytes.Buffer
	underlying := NewConsoleWriter(WithConsoleOutput(&buf))
	w := NewAsyncBatchWriter(
		WithAsyncUnderlying(underlying),
		WithAsyncBatchSize(1000),
		WithAsyncFlushInterval(10*time.Second),
		WithAsyncChannelSize(1), // 极小 channel，确保快速填满
	)

	// 写入大量数据，迫使 channel 填满后降级同步写入
	for i := 0; i < 100; i++ {
		w.Write([]byte("overflow test\n"))
	}

	w.Close()

	output := buf.String()
	lines := strings.Count(output, "overflow test")
	if lines != 100 {
		t.Errorf("expected 100 lines, got %d", lines)
	}
}

// TestAsyncBatchWriter_ConcurrentWrite 测试并发写入
func TestAsyncBatchWriter_ConcurrentWrite(t *testing.T) {
	var buf bytes.Buffer
	underlying := NewConsoleWriter(WithConsoleOutput(&buf))
	w := NewAsyncBatchWriter(
		WithAsyncUnderlying(underlying),
		WithAsyncBatchSize(50),
		WithAsyncFlushInterval(50*time.Millisecond),
	)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				w.Write([]byte("concurrent write\n"))
			}
		}(i)
	}
	wg.Wait()

	// Close 会 drain 所有条目
	w.Close()

	output := buf.String()
	lines := strings.Count(output, "concurrent write")
	if lines != 1000 {
		t.Errorf("expected 1000 lines, got %d", lines)
	}
}

// TestAsyncBatchWriter_WithLogger 测试与 Logger 集成
func TestAsyncBatchWriter_WithLogger(t *testing.T) {
	var buf bytes.Buffer
	underlying := NewConsoleWriter(WithConsoleOutput(&buf))
	asyncW := NewAsyncBatchWriter(
		WithAsyncUnderlying(underlying),
		WithAsyncBatchSize(10),
		WithAsyncFlushInterval(50*time.Millisecond),
	)

	logger := NewLogger().
		WithOutput(asyncW).
		WithFormat(FormatJSON).
		WithShowCaller(false)

	logger.Info("async logger test")

	asyncW.Close()

	output := buf.String()
	if !strings.Contains(output, "async logger test") {
		t.Errorf("expected output to contain 'async logger test', got: %s", output)
	}
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("expected JSON output with level INFO, got: %s", output)
	}
}

// TestAsyncBatchWriter_HealthCheck 测试健康状态
// 注意：underlying 必须使用 bytes.Buffer 而非默认的 os.Stdout，
// 否则 Close() 会关闭全局 os.Stdout，破坏后续测试和测试框架自身的输出
func TestAsyncBatchWriter_HealthCheck(t *testing.T) {
	underlying := NewConsoleWriter(WithConsoleOutput(&bytes.Buffer{}))
	w := NewAsyncBatchWriter(
		WithAsyncUnderlying(underlying),
	)
	defer w.Close()

	if !w.IsHealthy() {
		t.Error("expected writer to be healthy")
	}

	w.Close()

	if w.IsHealthy() {
		t.Error("expected writer to be unhealthy after Close()")
	}
}

// TestAsyncBatchWriter_Stats 测试统计信息
func TestAsyncBatchWriter_Stats(t *testing.T) {
	underlying := NewConsoleWriter(WithConsoleOutput(&bytes.Buffer{}))
	w := NewAsyncBatchWriter(
		WithAsyncUnderlying(underlying),
		WithAsyncBatchSize(1000),
		WithAsyncFlushInterval(10*time.Second),
	)
	defer w.Close()

	w.Write([]byte("test data"))

	stats := w.GetStats()
	if stats.BytesWritten != 9 { // len("test data") = 9
		t.Errorf("expected 9 bytes written, got %d", stats.BytesWritten)
	}
	if stats.LinesWritten != 1 {
		t.Errorf("expected 1 line written, got %d", stats.LinesWritten)
	}
}

// TestLogSpecialContext 测试带 context 的特殊日志方法
func TestLogSpecialContext(t *testing.T) {
	// text 模式测试
	t.Run("Text", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger().
			WithOutput(&buf).
			WithFormat(FormatText).
			WithShowCaller(false)

		ctx := context.WithValue(context.Background(), ContextKeyTraceID, "trace-123")

		logger.LogSpecialContext(ctx, SuccessType, INFO, "操作成功")
		logger.LogSpecialContext(ctx, LoadingType, INFO, "正在加载")
		logger.LogSpecialContext(ctx, DatabaseType, INFO, "DB查询耗时 %v", 50*time.Millisecond)

		output := buf.String()
		if !strings.Contains(output, "trace-123") {
			t.Error("expected trace ID in text output")
		}
		if !strings.Contains(output, "✅") {
			t.Error("expected success emoji")
		}
		if !strings.Contains(output, "[LOADING]") {
			t.Error("expected LOADING tag")
		}
		if !strings.Contains(output, "[DATABASE]") {
			t.Error("expected DATABASE tag")
		}
	})

	// JSON 模式测试
	t.Run("JSON", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewLogger().
			WithOutput(&buf).
			WithFormat(FormatJSON).
			WithShowCaller(false)

		ctx := context.WithValue(context.Background(), ContextKeyTraceID, "trace-456")

		logger.LogSpecialContext(ctx, SuccessType, INFO, "JSON模式成功")
		logger.LogSpecialContext(ctx, CacheType, INFO, "缓存命中 key=%s", "user:123")

		output := buf.String()
		if !strings.Contains(output, "trace-456") {
			t.Error("expected trace ID in JSON output")
		}
		if !strings.Contains(output, "[SUCCESS]") {
			t.Error("expected SUCCESS tag in JSON")
		}
		if !strings.Contains(output, "[CACHE]") {
			t.Error("expected CACHE tag in JSON")
		}
	})
}

// BenchmarkAsyncBatchWriter 对比异步 vs 同步写入性能
func BenchmarkAsyncBatchWriter(b *testing.B) {
	underlying := NewConsoleWriter(WithConsoleOutput(&bytes.Buffer{}))
	w := NewAsyncBatchWriter(
		WithAsyncUnderlying(underlying),
		WithAsyncBatchSize(100),
		WithAsyncFlushInterval(100*time.Millisecond),
	)
	defer w.Close()

	data := []byte(`{"timestamp":"2026-01-01T00:00:00Z","level":"INFO","message":"benchmark test"}` + "\n")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w.Write(data)
		}
	})
}

// BenchmarkSyncWriter 同步写入对比基准
func BenchmarkSyncWriter(b *testing.B) {
	w := NewConsoleWriter(WithConsoleOutput(&bytes.Buffer{}))

	data := []byte(`{"timestamp":"2026-01-01T00:00:00Z","level":"INFO","message":"benchmark test"}` + "\n")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w.Write(data)
		}
	})
}
