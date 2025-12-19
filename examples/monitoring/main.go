/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 13:11:00
 * @FilePath: \go-logger\examples\monitoring\main.go
 * @Description: 监控示例 - 演示日志监控、指标收集和健康检查
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kamalyes/go-logger"
)

// 监控指标
type Metrics struct {
	TotalLogs   int64
	ErrorLogs   int64
	WarningLogs int64
	InfoLogs    int64
	DebugLogs   int64
	StartTime   time.Time
	LastLogTime time.Time
	mu          sync.RWMutex
}

func NewMetrics() *Metrics {
	return &Metrics{
		StartTime: time.Now(),
	}
}

func (m *Metrics) IncrementLog(level logger.LogLevel) {
	atomic.AddInt64(&m.TotalLogs, 1)
	m.mu.Lock()
	m.LastLogTime = time.Now()
	m.mu.Unlock()

	switch level {
	case logger.DEBUG:
		atomic.AddInt64(&m.DebugLogs, 1)
	case logger.INFO:
		atomic.AddInt64(&m.InfoLogs, 1)
	case logger.WARN:
		atomic.AddInt64(&m.WarningLogs, 1)
	case logger.ERROR:
		atomic.AddInt64(&m.ErrorLogs, 1)
	}
}

func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	uptime := time.Since(m.StartTime)
	logsPerSecond := float64(atomic.LoadInt64(&m.TotalLogs)) / uptime.Seconds()

	return map[string]interface{}{
		"total_logs":     atomic.LoadInt64(&m.TotalLogs),
		"error_logs":     atomic.LoadInt64(&m.ErrorLogs),
		"warning_logs":   atomic.LoadInt64(&m.WarningLogs),
		"info_logs":      atomic.LoadInt64(&m.InfoLogs),
		"debug_logs":     atomic.LoadInt64(&m.DebugLogs),
		"uptime_seconds": int64(uptime.Seconds()),
		"logs_per_sec":   logsPerSecond,
		"last_log_time":  m.LastLogTime.Format(time.RFC3339),
	}
}

// 监控适配器
type MonitoringAdapter struct {
	adapter logger.IAdapter
	metrics *Metrics
}

func NewMonitoringAdapter(adapter logger.IAdapter, metrics *Metrics) *MonitoringAdapter {
	return &MonitoringAdapter{
		adapter: adapter,
		metrics: metrics,
	}
}

func (m *MonitoringAdapter) Initialize() error {
	return m.adapter.Initialize()
}

func (m *MonitoringAdapter) Close() error {
	return m.adapter.Close()
}

func (m *MonitoringAdapter) Flush() error {
	return m.adapter.Flush()
}

func (m *MonitoringAdapter) GetAdapterName() string {
	return "monitoring-" + m.adapter.GetAdapterName()
}

func (m *MonitoringAdapter) GetAdapterVersion() string {
	return m.adapter.GetAdapterVersion()
}

func (m *MonitoringAdapter) IsHealthy() bool {
	return m.adapter.IsHealthy()
}

// 实现ILogger接口
func (m *MonitoringAdapter) Debug(format string, args ...interface{}) {
	m.metrics.IncrementLog(logger.DEBUG)
	m.adapter.Debug(format, args...)
}

func (m *MonitoringAdapter) Info(format string, args ...interface{}) {
	m.metrics.IncrementLog(logger.INFO)
	m.adapter.Info(format, args...)
}

func (m *MonitoringAdapter) Warn(format string, args ...interface{}) {
	m.metrics.IncrementLog(logger.WARN)
	m.adapter.Warn(format, args...)
}

func (m *MonitoringAdapter) Error(format string, args ...interface{}) {
	m.metrics.IncrementLog(logger.ERROR)
	m.adapter.Error(format, args...)
}

func (m *MonitoringAdapter) Fatal(format string, args ...interface{}) {
	m.metrics.IncrementLog(logger.ERROR)
	m.adapter.Fatal(format, args...)
}

// 委托其他方法
func (m *MonitoringAdapter) SetLevel(level logger.LogLevel) { m.adapter.SetLevel(level) }
func (m *MonitoringAdapter) GetLevel() logger.LogLevel      { return m.adapter.GetLevel() }
func (m *MonitoringAdapter) SetShowCaller(show bool)        { m.adapter.SetShowCaller(show) }
func (m *MonitoringAdapter) IsShowCaller() bool             { return m.adapter.IsShowCaller() }
func (m *MonitoringAdapter) IsLevelEnabled(level logger.LogLevel) bool {
	return m.adapter.IsLevelEnabled(level)
}
func (m *MonitoringAdapter) WithField(key string, value interface{}) logger.ILogger  { return m }
func (m *MonitoringAdapter) WithFields(fields map[string]interface{}) logger.ILogger { return m }
func (m *MonitoringAdapter) WithError(err error) logger.ILogger                      { return m }
func (m *MonitoringAdapter) Clone() logger.ILogger                                   { return m }
func (m *MonitoringAdapter) Debugf(format string, args ...interface{})               { m.Debug(format, args...) }
func (m *MonitoringAdapter) Infof(format string, args ...interface{})                { m.Info(format, args...) }
func (m *MonitoringAdapter) Warnf(format string, args ...interface{})                { m.Warn(format, args...) }
func (m *MonitoringAdapter) Errorf(format string, args ...interface{})               { m.Error(format, args...) }
func (m *MonitoringAdapter) Fatalf(format string, args ...interface{})               { m.Fatal(format, args...) }
func (m *MonitoringAdapter) DebugMsg(msg string)                                     { m.Debug(msg) }
func (m *MonitoringAdapter) InfoMsg(msg string)                                      { m.Info(msg) }
func (m *MonitoringAdapter) WarnMsg(msg string)                                      { m.Warn(msg) }
func (m *MonitoringAdapter) ErrorMsg(msg string)                                     { m.Error(msg) }
func (m *MonitoringAdapter) FatalMsg(msg string)                                     { m.Fatal(msg) }
func (m *MonitoringAdapter) Print(args ...interface{})                               { m.adapter.Print(args...) }
func (m *MonitoringAdapter) Printf(format string, args ...interface{}) {
	m.adapter.Printf(format, args...)
}
func (m *MonitoringAdapter) Println(args ...interface{}) { m.adapter.Println(args...) }
func (m *MonitoringAdapter) DebugContext(ctx context.Context, format string, args ...interface{}) {
	m.Debug(format, args...)
}
func (m *MonitoringAdapter) InfoContext(ctx context.Context, format string, args ...interface{}) {
	m.Info(format, args...)
}
func (m *MonitoringAdapter) WarnContext(ctx context.Context, format string, args ...interface{}) {
	m.Warn(format, args...)
}
func (m *MonitoringAdapter) ErrorContext(ctx context.Context, format string, args ...interface{}) {
	m.Error(format, args...)
}
func (m *MonitoringAdapter) FatalContext(ctx context.Context, format string, args ...interface{}) {
	m.Fatal(format, args...)
}
func (m *MonitoringAdapter) WithContext(ctx context.Context) logger.ILogger   { return m }
func (m *MonitoringAdapter) DebugKV(msg string, keysAndValues ...interface{}) { m.Debug(msg) }
func (m *MonitoringAdapter) InfoKV(msg string, keysAndValues ...interface{})  { m.Info(msg) }
func (m *MonitoringAdapter) WarnKV(msg string, keysAndValues ...interface{})  { m.Warn(msg) }
func (m *MonitoringAdapter) ErrorKV(msg string, keysAndValues ...interface{}) { m.Error(msg) }
func (m *MonitoringAdapter) FatalKV(msg string, keysAndValues ...interface{}) { m.Fatal(msg) }
func (m *MonitoringAdapter) DebugContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.Debug(msg)
}
func (m *MonitoringAdapter) InfoContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.Info(msg)
}
func (m *MonitoringAdapter) WarnContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.Warn(msg)
}
func (m *MonitoringAdapter) ErrorContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.Error(msg)
}
func (m *MonitoringAdapter) FatalContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {
	m.Fatal(msg)
}
func (m *MonitoringAdapter) Log(level logger.LogLevel, msg string) {
	switch level {
	case logger.DEBUG:
		m.Debug(msg)
	case logger.INFO:
		m.Info(msg)
	case logger.WARN:
		m.Warn(msg)
	case logger.ERROR:
		m.Error(msg)
	case logger.FATAL:
		m.Fatal(msg)
	}
}
func (m *MonitoringAdapter) LogContext(ctx context.Context, level logger.LogLevel, msg string) {
	m.Log(level, msg)
}
func (m *MonitoringAdapter) LogKV(level logger.LogLevel, msg string, keysAndValues ...interface{}) {
	m.Log(level, msg)
}
func (m *MonitoringAdapter) LogWithFields(level logger.LogLevel, msg string, fields map[string]interface{}) {
	m.Log(level, msg)
}

// 多行日志方法
func (m *MonitoringAdapter) DebugLines(lines ...string) {
	for _, line := range lines {
		m.Debug("%s", line)
	}
}

func (m *MonitoringAdapter) InfoLines(lines ...string) {
	for _, line := range lines {
		m.Info("%s", line)
	}
}

func (m *MonitoringAdapter) WarnLines(lines ...string) {
	for _, line := range lines {
		m.Warn("%s", line)
	}
}

func (m *MonitoringAdapter) ErrorLines(lines ...string) {
	for _, line := range lines {
		m.Error("%s", line)
	}
}

// 返回错误的日志方法
func (m *MonitoringAdapter) DebugReturn(format string, args ...interface{}) error {
	m.Debug(format, args...)
	return fmt.Errorf(format, args...)
}

func (m *MonitoringAdapter) InfoReturn(format string, args ...interface{}) error {
	m.Info(format, args...)
	return fmt.Errorf(format, args...)
}

func (m *MonitoringAdapter) WarnReturn(format string, args ...interface{}) error {
	m.Warn(format, args...)
	return fmt.Errorf(format, args...)
}

func (m *MonitoringAdapter) ErrorReturn(format string, args ...interface{}) error {
	m.Error(format, args...)
	return fmt.Errorf(format, args...)
}

func (m *MonitoringAdapter) DebugCtxReturn(ctx context.Context, format string, args ...interface{}) error {
	m.DebugContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

func (m *MonitoringAdapter) InfoCtxReturn(ctx context.Context, format string, args ...interface{}) error {
	m.InfoContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

func (m *MonitoringAdapter) WarnCtxReturn(ctx context.Context, format string, args ...interface{}) error {
	m.WarnContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

func (m *MonitoringAdapter) ErrorCtxReturn(ctx context.Context, format string, args ...interface{}) error {
	m.ErrorContext(ctx, format, args...)
	return fmt.Errorf(format, args...)
}

func (m *MonitoringAdapter) DebugKVReturn(msg string, keysAndValues ...interface{}) error {
	m.DebugKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

func (m *MonitoringAdapter) InfoKVReturn(msg string, keysAndValues ...interface{}) error {
	m.InfoKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

func (m *MonitoringAdapter) WarnKVReturn(msg string, keysAndValues ...interface{}) error {
	m.WarnKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

func (m *MonitoringAdapter) ErrorKVReturn(msg string, keysAndValues ...interface{}) error {
	m.ErrorKV(msg, keysAndValues...)
	return fmt.Errorf("%s", msg)
}

// Console 相关方法
func (m *MonitoringAdapter) ConsoleGroup(label string, args ...interface{}) {
	m.adapter.ConsoleGroup(label, args...)
}
func (m *MonitoringAdapter) ConsoleGroupCollapsed(label string, args ...interface{}) {
	m.adapter.ConsoleGroupCollapsed(label, args...)
}
func (m *MonitoringAdapter) ConsoleGroupEnd() {
	m.adapter.ConsoleGroupEnd()
}
func (m *MonitoringAdapter) ConsoleTable(data interface{}) {
	m.adapter.ConsoleTable(data)
}
func (m *MonitoringAdapter) ConsoleTime(label string) *logger.Timer {
	return m.adapter.ConsoleTime(label)
}
func (m *MonitoringAdapter) NewConsoleGroup() *logger.ConsoleGroup {
	return m.adapter.NewConsoleGroup()
}

func main() {
	fmt.Println("📊 Go Logger - 监控示例演示")
	fmt.Println(strings.Repeat("=", 50))

	// 1. 基础监控演示
	demonstrateBasicMonitoring()

	fmt.Println()

	// 2. 实时指标收集
	demonstrateMetricsCollection()

	fmt.Println()

	// 3. 健康检查演示
	demonstrateHealthCheck()

	fmt.Println()

	// 4. 性能监控
	demonstratePerformanceMonitoring()

	fmt.Println()

	// 5. 告警模拟
	demonstrateAlerting()
}

func demonstrateBasicMonitoring() {
	fmt.Println("📈 1. 基础监控演示")
	fmt.Println(strings.Repeat("-", 30))

	// 创建指标收集器
	metrics := NewMetrics()

	// 创建基础适配器
	baseAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.DEBUG,
		Output:     os.Stdout,
		TimeFormat: "15:04:05",
		Colorful:   true,
	})
	baseAdapter.Initialize()

	// 包装为监控适配器
	monitoringAdapter := NewMonitoringAdapter(baseAdapter, metrics)

	fmt.Println("\n🔹 模拟应用日志:")
	monitoringAdapter.Info("应用启动完成")
	monitoringAdapter.Debug("加载配置文件")
	monitoringAdapter.Info("连接数据库成功")
	monitoringAdapter.Warn("内存使用率较高: 85%")
	monitoringAdapter.Error("连接超时")
	monitoringAdapter.Info("重试连接成功")

	// 显示统计信息
	fmt.Println("\n📊 当前监控指标:")
	stats := metrics.GetStats()
	for key, value := range stats {
		fmt.Printf("  %s: %v\n", key, value)
	}

	defer baseAdapter.Close()
}

func demonstrateMetricsCollection() {
	fmt.Println("📋 2. 实时指标收集")
	fmt.Println(strings.Repeat("-", 30))

	metrics := NewMetrics()

	baseAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:   logger.StandardAdapter,
		Level:  logger.DEBUG,
		Output: os.Stdout,
	})
	baseAdapter.Initialize()

	monitoringAdapter := NewMonitoringAdapter(baseAdapter, metrics)

	// 启动指标收集协程
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 模拟并发日志生成
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				select {
				case <-ctx.Done():
					return
				default:
					level := rand.Intn(4)
					switch level {
					case 0:
						monitoringAdapter.Debug("Worker %d: Debug message %d", workerID, j)
					case 1:
						monitoringAdapter.Info("Worker %d: Info message %d", workerID, j)
					case 2:
						monitoringAdapter.Warn("Worker %d: Warning message %d", workerID, j)
					case 3:
						monitoringAdapter.Error("Worker %d: Error message %d", workerID, j)
					}
					time.Sleep(50 * time.Millisecond)
				}
			}
		}(i)
	}

	// 实时显示指标
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stats := metrics.GetStats()
				fmt.Printf("\r实时统计: 总日志=%d, 错误=%d, 警告=%d, 信息=%d, 调试=%d, 速率=%.1f/s",
					stats["total_logs"], stats["error_logs"], stats["warning_logs"],
					stats["info_logs"], stats["debug_logs"], stats["logs_per_sec"])
			}
		}
	}()

	wg.Wait()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n\n📈 最终统计:")
	stats := metrics.GetStats()
	for key, value := range stats {
		fmt.Printf("  %s: %v\n", key, value)
	}

	defer baseAdapter.Close()
}

func demonstrateHealthCheck() {
	fmt.Println("🏥 3. 健康检查演示")
	fmt.Println(strings.Repeat("-", 30))

	adapters := make([]logger.IAdapter, 0)

	// 创建多个适配器
	for j := 0; j < 3; j++ {
		adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
			Type:   logger.StandardAdapter,
			Name:   fmt.Sprintf("adapter-%d", j+1),
			Level:  logger.INFO,
			Output: os.Stdout,
		})
		adapter.Initialize()
		adapters = append(adapters, adapter)
	}

	fmt.Println("\n🔹 健康检查结果:")
	healthyCount := 0
	for _, adapter := range adapters {
		isHealthy := adapter.IsHealthy()
		status := "❌ 不健康"
		if isHealthy {
			status = "✅ 健康"
			healthyCount++
		}

		fmt.Printf("  %s: %s (版本: %s)\n",
			adapter.GetAdapterName(), status, adapter.GetAdapterVersion())
	}

	fmt.Printf("\n📊 健康状态汇总: %d/%d 适配器健康\n", healthyCount, len(adapters))

	if healthyCount == len(adapters) {
		fmt.Println("🎉 所有适配器运行正常!")
	} else {
		fmt.Println("⚠️ 存在不健康的适配器，需要检查!")
	}

	// 清理
	for _, adapter := range adapters {
		adapter.Close()
	}
}

func demonstratePerformanceMonitoring() {
	fmt.Println("⚡ 4. 性能监控")
	fmt.Println(strings.Repeat("-", 30))

	adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:   logger.StandardAdapter,
		Level:  logger.INFO,
		Output: os.Stdout,
	})
	adapter.Initialize()

	fmt.Println("\n🔹 性能基准测试:")

	// 测试单线程性能
	start := time.Now()
	for i := 0; i < 1000; i++ {
		adapter.Info("Performance test message %d", i)
	}
	singleDuration := time.Since(start)

	// 测试并发性能
	start = time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				adapter.Info("Concurrent test worker %d message %d", workerID, j)
			}
		}(i)
	}
	wg.Wait()
	concurrentDuration := time.Since(start)

	fmt.Printf("单线程 1000 条日志耗时: %v (%.2f μs/条)\n",
		singleDuration, float64(singleDuration.Nanoseconds())/1000.0/1000.0)
	fmt.Printf("并发 1000 条日志耗时: %v (%.2f μs/条)\n",
		concurrentDuration, float64(concurrentDuration.Nanoseconds())/1000.0/1000.0)

	improvement := float64(singleDuration.Nanoseconds()) / float64(concurrentDuration.Nanoseconds())
	fmt.Printf("并发性能提升: %.2fx\n", improvement)

	defer adapter.Close()
}

func demonstrateAlerting() {
	fmt.Println("🚨 5. 告警模拟")
	fmt.Println(strings.Repeat("-", 30))

	metrics := NewMetrics()

	baseAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:   logger.StandardAdapter,
		Level:  logger.DEBUG,
		Output: os.Stdout,
	})
	baseAdapter.Initialize()

	monitoringAdapter := NewMonitoringAdapter(baseAdapter, metrics)

	// 模拟正常运行
	fmt.Println("\n🔹 正常运行期:")
	for i := 0; i < 5; i++ {
		monitoringAdapter.Info("正常业务处理 %d", i+1)
		time.Sleep(100 * time.Millisecond)
	}

	// 模拟异常情况
	fmt.Println("\n🔹 异常检测:")
	errorThreshold := int64(3)

	for i := 0; i < 5; i++ {
		monitoringAdapter.Error("数据库连接失败 %d", i+1)

		currentErrors := atomic.LoadInt64(&metrics.ErrorLogs)
		if currentErrors >= errorThreshold {
			fmt.Printf("\n🚨 告警触发! 错误日志数量达到阈值: %d >= %d\n", currentErrors, errorThreshold)
			fmt.Println("📧 发送告警邮件...")
			fmt.Println("📱 推送告警通知...")
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	// 显示最终指标
	fmt.Println("\n📊 告警期间指标:")
	stats := metrics.GetStats()
	for key, value := range stats {
		fmt.Printf("  %s: %v\n", key, value)
	}

	defer baseAdapter.Close()
}
