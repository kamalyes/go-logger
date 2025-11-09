package logger

import (
	"io"
	"log"
	"log/slog"
	"testing"
	"time"
)

// BenchmarkGoLogger 测试我们的 logger 性能
func BenchmarkGoLogger(b *testing.B) {
	config := DefaultConfig().
		WithLevel(INFO).
		WithShowCaller(false).
		WithColorful(false).
		WithOutput(io.Discard) // 输出到空设备，避免 I/O 影响性能
	
	logger := NewLogger(config)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		logger.Info("This is a test message")
	}
}

// BenchmarkGoLoggerKV 测试键值对方式的性能
func BenchmarkGoLoggerKV(b *testing.B) {
	config := DefaultConfig().
		WithLevel(INFO).
		WithShowCaller(false).
		WithColorful(false).
		WithOutput(io.Discard)
	
	logger := NewLogger(config)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		logger.InfoKV("Test message with KV",
			"component", "benchmark",
			"iteration", i,
			"timestamp", time.Now(),
		)
	}
}

// BenchmarkSlog 对比标准库 slog 性能
func BenchmarkSlog(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		logger.Info("This is a test message")
	}
}

// BenchmarkStdLog 对比标准库 log 性能
func BenchmarkStdLog(b *testing.B) {
	logger := log.New(io.Discard, "", log.LstdFlags)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Print("This is a test message")
	}
}

// BenchmarkUltraFastLogger 测试极致优化版本的性能
func BenchmarkUltraFastLogger(b *testing.B) {
	config := DefaultConfig().
		WithLevel(INFO).
		WithShowCaller(false).
		WithColorful(false).
		WithOutput(io.Discard)
	
	logger := NewUltraFastLogger(config)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		logger.Info("This is a test message")
	}
}

// BenchmarkUltraFastLoggerNoTime 测试无时间戳的极致优化版本
func BenchmarkUltraFastLoggerNoTime(b *testing.B) {
	logger := NewUltraFastLoggerNoTime(io.Discard, INFO)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		logger.Info("This is a test message")
	}
}