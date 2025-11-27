/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 12:44:39
 * @FilePath: \go-logger\examples\adapters\main.go
 * @Description: 适配器系统示例 - 演示实际可用的适配器功能
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package main

import (
	"context"
	"fmt"
	"github.com/kamalyes/go-logger"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("🔌 Go Logger - 适配器系统示例演示")
	fmt.Println(strings.Repeat("=", 50))

	// 1. 标准适配器演示
	demonstrateStandardAdapter()

	fmt.Println()

	// 2. 多适配器配置
	demonstrateMultipleAdapters()

	fmt.Println()

	// 3. 适配器配置和管理
	demonstrateAdapterConfiguration()

	fmt.Println()

	// 4. 自定义适配器扩展
	demonstrateCustomAdapterExtension()

	fmt.Println()

	// 5. 实际应用示例
	demonstrateRealWorldExample()
}

// 标准适配器演示
func demonstrateStandardAdapter() {
	fmt.Println("📋 1. 标准适配器演示")
	fmt.Println(strings.Repeat("-", 30))

	// 1.1 基础标准适配器
	fmt.Println("\n🔹 基础标准适配器:")

	// 创建标准适配器配置
	config := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "console-adapter",
		Level:      logger.INFO,
		Output:     os.Stdout,
		Colorful:   true,
		TimeFormat: "15:04:05",
		Fields: map[string]interface{}{
			"service": "demo-service",
			"version": "1.0.0",
		},
	}

	// 创建标准适配器
	adapter, err := logger.NewStandardAdapter(config)
	if err != nil {
		fmt.Printf("❌ 创建适配器失败: %v\n", err)
		return
	}

	// 初始化适配器
	if err := adapter.Initialize(); err != nil {
		fmt.Printf("❌ 初始化适配器失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 适配器创建成功: %s (版本: %s)\n",
		adapter.GetAdapterName(), adapter.GetAdapterVersion())
	fmt.Printf("✅ 适配器健康状态: %v\n", adapter.IsHealthy())

	// 使用适配器记录日志
	adapter.Info("这是通过标准适配器输出的信息日志")
	adapter.WithField("user_id", 12345).Info("带字段的日志")
	adapter.WithFields(map[string]interface{}{
		"action":    "login",
		"timestamp": time.Now().Unix(),
	}).Info("带多个字段的日志")

	// 测试不同级别的日志
	adapter.Debug("调试信息 (可能不会显示，取决于级别设置)")
	adapter.Warn("警告信息")
	adapter.Error("错误信息")

	// 清理
	defer adapter.Close()
}

// 多适配器配置
func demonstrateMultipleAdapters() {
	fmt.Println("🔀 2. 多适配器配置")
	fmt.Println(strings.Repeat("-", 30))

	// 创建日志目录
	logDir := "./logs"
	os.MkdirAll(logDir, 0755)

	// 2.1 控制台适配器
	fmt.Println("\n🔹 控制台适配器:")
	consoleConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "console",
		Level:      logger.INFO,
		Output:     os.Stdout,
		Colorful:   true,
		TimeFormat: "15:04:05",
		Format:     "text",
	}

	consoleAdapter, err := logger.NewStandardAdapter(consoleConfig)
	if err != nil {
		fmt.Printf("❌ 创建控制台适配器失败: %v\n", err)
		return
	}
	consoleAdapter.Initialize()

	// 2.2 文件适配器 (通过重定向输出)
	fmt.Println("\n🔹 文件适配器:")
	logFile, err := os.OpenFile(filepath.Join(logDir, "adapter.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("❌ 创建日志文件失败: %v\n", err)
		return
	}

	fileConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "file",
		Level:      logger.DEBUG,
		Output:     logFile,
		Colorful:   false, // 文件中不需要颜色
		TimeFormat: "2006-01-02 15:04:05",
		Format:     "json",
		Fields: map[string]interface{}{
			"component": "file-logger",
		},
	}

	fileAdapter, err := logger.NewStandardAdapter(fileConfig)
	if err != nil {
		fmt.Printf("❌ 创建文件适配器失败: %v\n", err)
		return
	}
	fileAdapter.Initialize()

	// 测试两个适配器
	fmt.Println("\n🔹 多适配器测试:")

	// 同时写入控制台和文件
	testMessage := "这是一条测试消息"
	consoleAdapter.Info("控制台: %s", testMessage)
	fileAdapter.Info("文件: %s", testMessage)

	// 带字段的日志
	fields := map[string]interface{}{
		"user_id":  12345,
		"action":   "test",
		"duration": "150ms",
	}

	consoleAdapter.WithFields(fields).Info("控制台带字段日志")
	fileAdapter.WithFields(fields).Info("文件带字段日志")

	fmt.Printf("✅ 日志已写入文件: %s\n", filepath.Join(logDir, "adapter.log"))

	// 清理
	defer consoleAdapter.Close()
	defer fileAdapter.Close()
	defer logFile.Close()
}

// 适配器配置和管理
func demonstrateAdapterConfiguration() {
	fmt.Println("⚙️ 3. 适配器配置和管理")
	fmt.Println(strings.Repeat("-", 30))

	// 3.1 不同级别配置
	fmt.Println("\n🔹 不同级别配置:")

	// DEBUG级别适配器
	debugConfig := &logger.AdapterConfig{
		Type:   logger.StandardAdapter,
		Name:   "debug-adapter",
		Level:  logger.DEBUG,
		Output: os.Stdout,
		Fields: map[string]interface{}{"level": "debug"},
	}

	debugAdapter, _ := logger.NewStandardAdapter(debugConfig)
	debugAdapter.Initialize()

	// INFO级别适配器
	infoConfig := &logger.AdapterConfig{
		Type:   logger.StandardAdapter,
		Name:   "info-adapter",
		Level:  logger.INFO,
		Output: os.Stdout,
		Fields: map[string]interface{}{"level": "info"},
	}

	infoAdapter, _ := logger.NewStandardAdapter(infoConfig)
	infoAdapter.Initialize()

	// 测试级别过滤
	fmt.Println("测试级别过滤:")
	debugAdapter.Debug("DEBUG适配器: 调试信息") // 会显示
	debugAdapter.Info("DEBUG适配器: 普通信息")  // 会显示

	infoAdapter.Debug("INFO适配器: 调试信息") // 不会显示
	infoAdapter.Info("INFO适配器: 普通信息")  // 会显示

	// 3.2 动态配置修改
	fmt.Println("\n🔹 动态配置修改:")

	// 创建可配置的适配器
	dynamicConfig := &logger.AdapterConfig{
		Type:   logger.StandardAdapter,
		Name:   "dynamic-adapter",
		Level:  logger.WARN,
		Output: os.Stdout,
	}

	dynamicAdapter, _ := logger.NewStandardAdapter(dynamicConfig)
	dynamicAdapter.Initialize()

	fmt.Println("初始配置 (WARN级别):")
	dynamicAdapter.Debug("调试信息 (不显示)")
	dynamicAdapter.Info("普通信息 (不显示)")
	dynamicAdapter.Warn("警告信息 (显示)")

	// 动态修改级别
	dynamicAdapter.SetLevel(logger.DEBUG)
	fmt.Println("修改为DEBUG级别后:")
	dynamicAdapter.Debug("调试信息 (现在显示)")
	dynamicAdapter.Info("普通信息 (现在显示)")

	// 3.3 适配器健康检查
	fmt.Println("\n🔹 适配器健康检查:")
	adapters := []logger.IAdapter{debugAdapter, infoAdapter, dynamicAdapter}

	for i, adapter := range adapters {
		fmt.Printf("适配器 %d [%s]: 健康状态 = %v\n",
			i+1, adapter.GetAdapterName(), adapter.IsHealthy())
	}

	// 清理
	defer debugAdapter.Close()
	defer infoAdapter.Close()
	defer dynamicAdapter.Close()
}

// 自定义适配器扩展
func demonstrateCustomAdapterExtension() {
	fmt.Println("🛠️ 4. 自定义适配器扩展")
	fmt.Println(strings.Repeat("-", 30))

	// 4.1 内存缓存适配器
	fmt.Println("\n🔹 内存缓存适配器:")
	memoryAdapter := NewMemoryAdapter(100)

	// 写入一些测试日志
	memoryAdapter.Info("内存日志 1")
	memoryAdapter.WithField("test", true).Info("内存日志 2")
	memoryAdapter.Error("内存错误日志")

	logs := memoryAdapter.GetLogs()
	fmt.Printf("内存中缓存的日志数量: %d\n", len(logs))
	for i, logMsg := range logs {
		fmt.Printf("  [%d] %s\n", i+1, logMsg)
	}

	// 4.2 过滤适配器
	fmt.Println("\n🔹 过滤适配器:")

	baseAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:   logger.StandardAdapter,
		Name:   "base",
		Level:  logger.DEBUG,
		Output: os.Stdout,
	})
	baseAdapter.Initialize()

	// 创建只允许错误级别通过的过滤适配器
	filterAdapter := NewFilterAdapter(baseAdapter, func(level logger.LogLevel, msg string) bool {
		return level >= logger.ERROR
	})

	fmt.Println("过滤器测试 (只允许ERROR级别通过):")
	filterAdapter.Debug("调试信息 (被过滤)")
	filterAdapter.Info("普通信息 (被过滤)")
	filterAdapter.Warn("警告信息 (被过滤)")
	filterAdapter.Error("错误信息 (通过过滤)")

	// 4.3 统计适配器
	fmt.Println("\n🔹 统计适配器:")

	statsAdapter := NewStatsAdapter(baseAdapter)

	// 写入各种级别的日志
	statsAdapter.Debug("调试日志")
	statsAdapter.Info("信息日志 1")
	statsAdapter.Info("信息日志 2")
	statsAdapter.Warn("警告日志")
	statsAdapter.Error("错误日志 1")
	statsAdapter.Error("错误日志 2")
	statsAdapter.Error("错误日志 3")

	stats := statsAdapter.GetStats()
	fmt.Printf("日志统计:\n")
	fmt.Printf("  DEBUG: %d\n", stats.DebugCount)
	fmt.Printf("  INFO:  %d\n", stats.InfoCount)
	fmt.Printf("  WARN:  %d\n", stats.WarnCount)
	fmt.Printf("  ERROR: %d\n", stats.ErrorCount)
	fmt.Printf("  总计:   %d\n", stats.TotalCount)

	defer memoryAdapter.Close()
	defer baseAdapter.Close()
}

// 实际应用示例
func demonstrateRealWorldExample() {
	fmt.Println("🌍 5. 实际应用示例")
	fmt.Println(strings.Repeat("-", 30))

	// 5.1 Web应用日志配置
	fmt.Println("\n🔹 Web应用日志配置:")

	// 创建日志目录
	logDir := "./logs"
	os.MkdirAll(logDir, 0755)

	// 访问日志 (控制台 + 文件)
	accessFile, _ := os.OpenFile(filepath.Join(logDir, "access.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	accessAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "access-log",
		Level:      logger.INFO,
		Output:     accessFile,
		TimeFormat: "2006-01-02 15:04:05",
		Fields: map[string]interface{}{
			"service": "web-server",
			"type":    "access",
		},
	})
	accessAdapter.Initialize()

	// 错误日志 (控制台 + 错误文件)
	errorFile, _ := os.OpenFile(filepath.Join(logDir, "error.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	errorAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "error-log",
		Level:      logger.ERROR,
		Output:     errorFile,
		TimeFormat: "2006-01-02 15:04:05",
		Fields: map[string]interface{}{
			"service": "web-server",
			"type":    "error",
		},
	})
	errorAdapter.Initialize()

	// 控制台日志
	consoleAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "console-log",
		Level:      logger.INFO,
		Output:     os.Stdout,
		Colorful:   true,
		TimeFormat: "15:04:05",
	})
	consoleAdapter.Initialize()

	// 模拟Web请求处理
	fmt.Println("模拟Web请求处理:")

	// 正常访问
	requestFields := map[string]interface{}{
		"method":   "GET",
		"url":      "/api/users",
		"ip":       "192.168.1.100",
		"status":   200,
		"duration": "150ms",
	}

	consoleAdapter.WithFields(requestFields).Info("API请求")
	accessAdapter.WithFields(requestFields).Info("API访问记录")

	// 错误处理
	errorFields := map[string]interface{}{
		"method":   "POST",
		"url":      "/api/orders",
		"ip":       "192.168.1.100",
		"status":   500,
		"error":    "数据库连接超时",
		"duration": "5000ms",
	}

	consoleAdapter.WithFields(errorFields).Error("API错误")
	accessAdapter.WithFields(errorFields).Error("API错误访问")
	errorAdapter.WithFields(errorFields).Error("服务器错误")

	// 5.2 高并发日志处理
	fmt.Println("\n🔹 高并发日志处理:")

	// 创建高性能日志适配器
	perfFile, _ := os.OpenFile(filepath.Join(logDir, "performance.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	perfAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "perf-log",
		Level:      logger.INFO,
		Output:     perfFile,
		TimeFormat: time.RFC3339Nano,
	})
	perfAdapter.Initialize()

	// 并发测试
	var wg sync.WaitGroup
	numWorkers := 10
	logsPerWorker := 100

	start := time.Now()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < logsPerWorker; j++ {
				perfAdapter.WithFields(map[string]interface{}{
					"worker_id": workerID,
					"task_id":   j,
					"timestamp": time.Now().UnixNano(),
				}).Info("并发处理任务")
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	totalLogs := numWorkers * logsPerWorker
	fmt.Printf("✅ 并发日志测试完成:\n")
	fmt.Printf("  工作协程: %d\n", numWorkers)
	fmt.Printf("  总日志数: %d\n", totalLogs)
	fmt.Printf("  总耗时: %v\n", duration)
	fmt.Printf("  平均耗时: %v/log\n", duration/time.Duration(totalLogs))

	// 5.3 日志轮转模拟
	fmt.Println("\n🔹 日志轮转模拟:")

	// 模拟大量日志写入
	rotateFile, _ := os.OpenFile(filepath.Join(logDir, "rotate.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	rotateAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "rotate-log",
		Level:      logger.DEBUG,
		Output:     rotateFile,
		TimeFormat: time.RFC3339,
	})
	rotateAdapter.Initialize()

	fmt.Println("模拟大量日志写入...")
	for i := 0; i < 1000; i++ {
		rotateAdapter.WithFields(map[string]interface{}{
			"sequence": i,
			"batch":    i / 100,
		}).Info("批量日志数据 %d", i)
	}

	// 检查文件大小
	if stat, err := rotateFile.Stat(); err == nil {
		fmt.Printf("✅ 日志文件大小: %d bytes\n", stat.Size())
	}

	fmt.Printf("✅ 所有日志文件已保存到: %s\n", logDir)

	// 清理资源
	defer accessAdapter.Close()
	defer errorAdapter.Close()
	defer consoleAdapter.Close()
	defer perfAdapter.Close()
	defer rotateAdapter.Close()
	defer accessFile.Close()
	defer errorFile.Close()
	defer perfFile.Close()
	defer rotateFile.Close()
}

// =============================================================================
// 自定义适配器实现
// =============================================================================

// MemoryAdapter - 内存适配器，缓存最近的日志
type MemoryAdapter struct {
	logs    []string
	maxSize int
	mu      sync.RWMutex
	level   logger.LogLevel
	name    string
	healthy bool
}

// NewMemoryAdapter 创建内存适配器
func NewMemoryAdapter(maxSize int) *MemoryAdapter {
	return &MemoryAdapter{
		logs:    make([]string, 0, maxSize),
		maxSize: maxSize,
		level:   logger.INFO,
		name:    "memory-adapter",
		healthy: true,
	}
}

// 实现 IAdapter 接口
func (a *MemoryAdapter) Initialize() error {
	a.healthy = true
	return nil
}

func (a *MemoryAdapter) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.logs = nil
	a.healthy = false
	return nil
}

func (a *MemoryAdapter) Flush() error {
	return nil
}

func (a *MemoryAdapter) GetAdapterName() string {
	return a.name
}

func (a *MemoryAdapter) GetAdapterVersion() string {
	return "1.0.0"
}

func (a *MemoryAdapter) IsHealthy() bool {
	return a.healthy
}

// 实现 ILogger 接口
func (a *MemoryAdapter) Debug(format string, args ...interface{}) {
	a.logMessage(logger.DEBUG, format, args...)
}

func (a *MemoryAdapter) Info(format string, args ...interface{}) {
	a.logMessage(logger.INFO, format, args...)
}

func (a *MemoryAdapter) Warn(format string, args ...interface{}) {
	a.logMessage(logger.WARN, format, args...)
}

func (a *MemoryAdapter) Error(format string, args ...interface{}) {
	a.logMessage(logger.ERROR, format, args...)
}

func (a *MemoryAdapter) Fatal(format string, args ...interface{}) {
	a.logMessage(logger.FATAL, format, args...)
}

func (a *MemoryAdapter) logMessage(level logger.LogLevel, format string, args ...interface{}) {
	if !a.healthy || level < a.level {
		return
	}

	message := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("%s [%s] %s", timestamp, level.String(), message)

	a.mu.Lock()
	defer a.mu.Unlock()

	// 如果超过最大容量，移除最老的日志
	if len(a.logs) >= a.maxSize {
		a.logs = a.logs[1:]
	}

	a.logs = append(a.logs, logEntry)
}

// 其他必需的方法实现（简化版本）
func (a *MemoryAdapter) SetLevel(level logger.LogLevel)                          { a.level = level }
func (a *MemoryAdapter) GetLevel() logger.LogLevel                               { return a.level }
func (a *MemoryAdapter) SetShowCaller(show bool)                                 {}
func (a *MemoryAdapter) IsShowCaller() bool                                      { return false }
func (a *MemoryAdapter) IsLevelEnabled(level logger.LogLevel) bool               { return level >= a.level }
func (a *MemoryAdapter) WithField(key string, value interface{}) logger.ILogger  { return a }
func (a *MemoryAdapter) WithFields(fields map[string]interface{}) logger.ILogger { return a }
func (a *MemoryAdapter) WithError(err error) logger.ILogger                      { return a }
func (a *MemoryAdapter) Clone() logger.ILogger                                   { return a }

// 扩展的方法实现（基础版本）
func (a *MemoryAdapter) Debugf(format string, args ...interface{}) { a.Debug(format, args...) }
func (a *MemoryAdapter) Infof(format string, args ...interface{})  { a.Info(format, args...) }
func (a *MemoryAdapter) Warnf(format string, args ...interface{})  { a.Warn(format, args...) }
func (a *MemoryAdapter) Errorf(format string, args ...interface{}) { a.Error(format, args...) }
func (a *MemoryAdapter) Fatalf(format string, args ...interface{}) { a.Fatal(format, args...) }

// 其他必需方法的空实现
func (a *MemoryAdapter) DebugMsg(msg string)                       { a.Debug("%s", msg) }
func (a *MemoryAdapter) InfoMsg(msg string)                        { a.Info("%s", msg) }
func (a *MemoryAdapter) WarnMsg(msg string)                        { a.Warn("%s", msg) }
func (a *MemoryAdapter) ErrorMsg(msg string)                       { a.Error("%s", msg) }
func (a *MemoryAdapter) FatalMsg(msg string)                       { a.Fatal("%s", msg) }
func (a *MemoryAdapter) Print(args ...interface{})                 {}
func (a *MemoryAdapter) Printf(format string, args ...interface{}) {}
func (a *MemoryAdapter) Println(args ...interface{})               {}

// 空实现的上下文方法
func (a *MemoryAdapter) DebugContext(ctx context.Context, format string, args ...interface{}) {
	a.Debug(format, args...)
}
func (a *MemoryAdapter) InfoContext(ctx context.Context, format string, args ...interface{}) {
	a.Info(format, args...)
}
func (a *MemoryAdapter) WarnContext(ctx context.Context, format string, args ...interface{}) {
	a.Warn(format, args...)
}
func (a *MemoryAdapter) ErrorContext(ctx context.Context, format string, args ...interface{}) {
	a.Error(format, args...)
}
func (a *MemoryAdapter) FatalContext(ctx context.Context, format string, args ...interface{}) {
	a.Fatal(format, args...)
}
func (a *MemoryAdapter) WithContext(ctx context.Context) logger.ILogger { return a }

// 空实现的KV方法
func (a *MemoryAdapter) DebugKV(msg string, keysAndValues ...interface{}) { a.Debug("%s", msg) }
func (a *MemoryAdapter) InfoKV(msg string, keysAndValues ...interface{})  { a.Info("%s", msg) }
func (a *MemoryAdapter) WarnKV(msg string, keysAndValues ...interface{})  { a.Warn("%s", msg) }
func (a *MemoryAdapter) ErrorKV(msg string, keysAndValues ...interface{}) { a.Error("%s", msg) }
func (a *MemoryAdapter) FatalKV(msg string, keysAndValues ...interface{}) { a.Fatal("%s", msg) }

// 空实现的原始日志方法
func (a *MemoryAdapter) Log(level logger.LogLevel, msg string) { a.logMessage(level, "%s", msg) }
func (a *MemoryAdapter) LogContext(ctx context.Context, level logger.LogLevel, msg string) {
	a.Log(level, msg)
}
func (a *MemoryAdapter) LogKV(level logger.LogLevel, msg string, keysAndValues ...interface{}) {
	a.Log(level, msg)
}
func (a *MemoryAdapter) LogWithFields(level logger.LogLevel, msg string, fields map[string]interface{}) {
	a.Log(level, msg)
}

// 多行日志方法实现
func (a *MemoryAdapter) DebugLines(lines ...string) {
	for _, line := range lines {
		a.Debug("%s", line)
	}
}

func (a *MemoryAdapter) InfoLines(lines ...string) {
	for _, line := range lines {
		a.Info("%s", line)
	}
}

func (a *MemoryAdapter) WarnLines(lines ...string) {
	for _, line := range lines {
		a.Warn("%s", line)
	}
}

func (a *MemoryAdapter) ErrorLines(lines ...string) {
	for _, line := range lines {
		a.Error("%s", line)
	}
}

// GetLogs 获取缓存的日志
func (a *MemoryAdapter) GetLogs() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]string, len(a.logs))
	copy(result, a.logs)
	return result
}

// FilterAdapter - 过滤适配器
type FilterAdapter struct {
	adapter logger.IAdapter
	filter  func(logger.LogLevel, string) bool
	mu      sync.RWMutex
}

// NewFilterAdapter 创建过滤适配器
func NewFilterAdapter(adapter logger.IAdapter, filter func(logger.LogLevel, string) bool) *FilterAdapter {
	return &FilterAdapter{
		adapter: adapter,
		filter:  filter,
	}
}

// 适配器接口实现
func (f *FilterAdapter) Initialize() error         { return f.adapter.Initialize() }
func (f *FilterAdapter) Close() error              { return f.adapter.Close() }
func (f *FilterAdapter) Flush() error              { return f.adapter.Flush() }
func (f *FilterAdapter) GetAdapterName() string    { return "filter-" + f.adapter.GetAdapterName() }
func (f *FilterAdapter) GetAdapterVersion() string { return f.adapter.GetAdapterVersion() }
func (f *FilterAdapter) IsHealthy() bool           { return f.adapter.IsHealthy() }

// 日志方法实现（带过滤）
func (f *FilterAdapter) Debug(format string, args ...interface{}) {
	if f.shouldLog(logger.DEBUG, format) {
		f.adapter.Debug(format, args...)
	}
}

func (f *FilterAdapter) Info(format string, args ...interface{}) {
	if f.shouldLog(logger.INFO, format) {
		f.adapter.Info(format, args...)
	}
}

func (f *FilterAdapter) Warn(format string, args ...interface{}) {
	if f.shouldLog(logger.WARN, format) {
		f.adapter.Warn(format, args...)
	}
}

func (f *FilterAdapter) Error(format string, args ...interface{}) {
	if f.shouldLog(logger.ERROR, format) {
		f.adapter.Error(format, args...)
	}
}

func (f *FilterAdapter) Fatal(format string, args ...interface{}) {
	if f.shouldLog(logger.FATAL, format) {
		f.adapter.Fatal(format, args...)
	}
}

func (f *FilterAdapter) shouldLog(level logger.LogLevel, format string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.filter(level, format)
}

// 委托给基础适配器的其他方法
func (f *FilterAdapter) SetLevel(level logger.LogLevel) { f.adapter.SetLevel(level) }
func (f *FilterAdapter) GetLevel() logger.LogLevel      { return f.adapter.GetLevel() }
func (f *FilterAdapter) SetShowCaller(show bool)        { f.adapter.SetShowCaller(show) }
func (f *FilterAdapter) IsShowCaller() bool             { return f.adapter.IsShowCaller() }
func (f *FilterAdapter) IsLevelEnabled(level logger.LogLevel) bool {
	return f.adapter.IsLevelEnabled(level)
}
func (f *FilterAdapter) WithField(key string, value interface{}) logger.ILogger  { return f }
func (f *FilterAdapter) WithFields(fields map[string]interface{}) logger.ILogger { return f }
func (f *FilterAdapter) WithError(err error) logger.ILogger                      { return f }
func (f *FilterAdapter) Clone() logger.ILogger                                   { return f }

// 其他方法的简单委托
func (f *FilterAdapter) Debugf(format string, args ...interface{}) { f.Debug(format, args...) }
func (f *FilterAdapter) Infof(format string, args ...interface{})  { f.Info(format, args...) }
func (f *FilterAdapter) Warnf(format string, args ...interface{})  { f.Warn(format, args...) }
func (f *FilterAdapter) Errorf(format string, args ...interface{}) { f.Error(format, args...) }
func (f *FilterAdapter) Fatalf(format string, args ...interface{}) { f.Fatal(format, args...) }
func (f *FilterAdapter) DebugMsg(msg string)                       { f.Debug(msg) }
func (f *FilterAdapter) InfoMsg(msg string)                        { f.Info(msg) }
func (f *FilterAdapter) WarnMsg(msg string)                        { f.Warn(msg) }
func (f *FilterAdapter) ErrorMsg(msg string)                       { f.Error(msg) }
func (f *FilterAdapter) FatalMsg(msg string)                       { f.Fatal(msg) }
func (f *FilterAdapter) Print(args ...interface{})                 { f.adapter.Print(args...) }
func (f *FilterAdapter) Printf(format string, args ...interface{}) { f.adapter.Printf(format, args...) }
func (f *FilterAdapter) Println(args ...interface{})               { f.adapter.Println(args...) }
func (f *FilterAdapter) DebugContext(ctx context.Context, format string, args ...interface{}) {
	f.Debug(format, args...)
}
func (f *FilterAdapter) InfoContext(ctx context.Context, format string, args ...interface{}) {
	f.Info(format, args...)
}
func (f *FilterAdapter) WarnContext(ctx context.Context, format string, args ...interface{}) {
	f.Warn(format, args...)
}
func (f *FilterAdapter) ErrorContext(ctx context.Context, format string, args ...interface{}) {
	f.Error(format, args...)
}
func (f *FilterAdapter) FatalContext(ctx context.Context, format string, args ...interface{}) {
	f.Fatal(format, args...)
}
func (f *FilterAdapter) WithContext(ctx context.Context) logger.ILogger   { return f }
func (f *FilterAdapter) DebugKV(msg string, keysAndValues ...interface{}) { f.Debug(msg) }
func (f *FilterAdapter) InfoKV(msg string, keysAndValues ...interface{})  { f.Info(msg) }
func (f *FilterAdapter) WarnKV(msg string, keysAndValues ...interface{})  { f.Warn(msg) }
func (f *FilterAdapter) ErrorKV(msg string, keysAndValues ...interface{}) { f.Error(msg) }
func (f *FilterAdapter) FatalKV(msg string, keysAndValues ...interface{}) { f.Fatal(msg) }
func (f *FilterAdapter) Log(level logger.LogLevel, msg string) {
	switch level {
	case logger.DEBUG:
		f.Debug(msg)
	case logger.INFO:
		f.Info(msg)
	case logger.WARN:
		f.Warn(msg)
	case logger.ERROR:
		f.Error(msg)
	case logger.FATAL:
		f.Fatal(msg)
	}
}
func (f *FilterAdapter) LogContext(ctx context.Context, level logger.LogLevel, msg string) {
	f.Log(level, msg)
}
func (f *FilterAdapter) LogKV(level logger.LogLevel, msg string, keysAndValues ...interface{}) {
	f.Log(level, msg)
}
func (f *FilterAdapter) LogWithFields(level logger.LogLevel, msg string, fields map[string]interface{}) {
	f.Log(level, msg)
}

// 多行日志方法实现
func (f *FilterAdapter) DebugLines(lines ...string) {
	for _, line := range lines {
		f.Debug("%s", line)
	}
}

func (f *FilterAdapter) InfoLines(lines ...string) {
	for _, line := range lines {
		f.Info("%s", line)
	}
}

func (f *FilterAdapter) WarnLines(lines ...string) {
	for _, line := range lines {
		f.Warn("%s", line)
	}
}

func (f *FilterAdapter) ErrorLines(lines ...string) {
	for _, line := range lines {
		f.Error("%s", line)
	}
}

// StatsAdapter - 统计适配器
type StatsAdapter struct {
	adapter logger.IAdapter
	stats   LogStats
	mu      sync.RWMutex
}

type LogStats struct {
	DebugCount int64
	InfoCount  int64
	WarnCount  int64
	ErrorCount int64
	TotalCount int64
}

// NewStatsAdapter 创建统计适配器
func NewStatsAdapter(adapter logger.IAdapter) *StatsAdapter {
	return &StatsAdapter{
		adapter: adapter,
	}
}

// 适配器接口实现
func (s *StatsAdapter) Initialize() error         { return s.adapter.Initialize() }
func (s *StatsAdapter) Close() error              { return s.adapter.Close() }
func (s *StatsAdapter) Flush() error              { return s.adapter.Flush() }
func (s *StatsAdapter) GetAdapterName() string    { return "stats-" + s.adapter.GetAdapterName() }
func (s *StatsAdapter) GetAdapterVersion() string { return s.adapter.GetAdapterVersion() }
func (s *StatsAdapter) IsHealthy() bool           { return s.adapter.IsHealthy() }

// 日志方法实现（带统计）
func (s *StatsAdapter) Debug(format string, args ...interface{}) {
	s.incrementCount(logger.DEBUG)
	s.adapter.Debug(format, args...)
}

func (s *StatsAdapter) Info(format string, args ...interface{}) {
	s.incrementCount(logger.INFO)
	s.adapter.Info(format, args...)
}

func (s *StatsAdapter) Warn(format string, args ...interface{}) {
	s.incrementCount(logger.WARN)
	s.adapter.Warn(format, args...)
}

func (s *StatsAdapter) Error(format string, args ...interface{}) {
	s.incrementCount(logger.ERROR)
	s.adapter.Error(format, args...)
}

func (s *StatsAdapter) Fatal(format string, args ...interface{}) {
	s.incrementCount(logger.FATAL)
	s.adapter.Fatal(format, args...)
}

func (s *StatsAdapter) incrementCount(level logger.LogLevel) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats.TotalCount++
	switch level {
	case logger.DEBUG:
		s.stats.DebugCount++
	case logger.INFO:
		s.stats.InfoCount++
	case logger.WARN:
		s.stats.WarnCount++
	case logger.ERROR:
		s.stats.ErrorCount++
	}
}

// GetStats 获取统计信息
func (s *StatsAdapter) GetStats() LogStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.stats
}

// 委托给基础适配器的其他方法（简化实现）
func (s *StatsAdapter) SetLevel(level logger.LogLevel) { s.adapter.SetLevel(level) }
func (s *StatsAdapter) GetLevel() logger.LogLevel      { return s.adapter.GetLevel() }
func (s *StatsAdapter) SetShowCaller(show bool)        { s.adapter.SetShowCaller(show) }
func (s *StatsAdapter) IsShowCaller() bool             { return s.adapter.IsShowCaller() }
func (s *StatsAdapter) IsLevelEnabled(level logger.LogLevel) bool {
	return s.adapter.IsLevelEnabled(level)
}
func (s *StatsAdapter) WithField(key string, value interface{}) logger.ILogger  { return s }
func (s *StatsAdapter) WithFields(fields map[string]interface{}) logger.ILogger { return s }
func (s *StatsAdapter) WithError(err error) logger.ILogger                      { return s }
func (s *StatsAdapter) Clone() logger.ILogger                                   { return s }
func (s *StatsAdapter) Debugf(format string, args ...interface{})               { s.Debug(format, args...) }
func (s *StatsAdapter) Infof(format string, args ...interface{})                { s.Info(format, args...) }
func (s *StatsAdapter) Warnf(format string, args ...interface{})                { s.Warn(format, args...) }
func (s *StatsAdapter) Errorf(format string, args ...interface{})               { s.Error(format, args...) }
func (s *StatsAdapter) Fatalf(format string, args ...interface{})               { s.Fatal(format, args...) }
func (s *StatsAdapter) DebugMsg(msg string)                                     { s.Debug(msg) }
func (s *StatsAdapter) InfoMsg(msg string)                                      { s.Info(msg) }
func (s *StatsAdapter) WarnMsg(msg string)                                      { s.Warn(msg) }
func (s *StatsAdapter) ErrorMsg(msg string)                                     { s.Error(msg) }
func (s *StatsAdapter) FatalMsg(msg string)                                     { s.Fatal(msg) }
func (s *StatsAdapter) Print(args ...interface{})                               { s.adapter.Print(args...) }
func (s *StatsAdapter) Printf(format string, args ...interface{})               { s.adapter.Printf(format, args...) }
func (s *StatsAdapter) Println(args ...interface{})                             { s.adapter.Println(args...) }
func (s *StatsAdapter) DebugContext(ctx context.Context, format string, args ...interface{}) {
	s.Debug(format, args...)
}
func (s *StatsAdapter) InfoContext(ctx context.Context, format string, args ...interface{}) {
	s.Info(format, args...)
}
func (s *StatsAdapter) WarnContext(ctx context.Context, format string, args ...interface{}) {
	s.Warn(format, args...)
}
func (s *StatsAdapter) ErrorContext(ctx context.Context, format string, args ...interface{}) {
	s.Error(format, args...)
}
func (s *StatsAdapter) FatalContext(ctx context.Context, format string, args ...interface{}) {
	s.Fatal(format, args...)
}
func (s *StatsAdapter) WithContext(ctx context.Context) logger.ILogger   { return s }
func (s *StatsAdapter) DebugKV(msg string, keysAndValues ...interface{}) { s.Debug(msg) }
func (s *StatsAdapter) InfoKV(msg string, keysAndValues ...interface{})  { s.Info(msg) }
func (s *StatsAdapter) WarnKV(msg string, keysAndValues ...interface{})  { s.Warn(msg) }
func (s *StatsAdapter) ErrorKV(msg string, keysAndValues ...interface{}) { s.Error(msg) }
func (s *StatsAdapter) FatalKV(msg string, keysAndValues ...interface{}) { s.Fatal(msg) }
func (s *StatsAdapter) Log(level logger.LogLevel, msg string) {
	switch level {
	case logger.DEBUG:
		s.Debug(msg)
	case logger.INFO:
		s.Info(msg)
	case logger.WARN:
		s.Warn(msg)
	case logger.ERROR:
		s.Error(msg)
	case logger.FATAL:
		s.Fatal(msg)
	}
}
func (s *StatsAdapter) LogContext(ctx context.Context, level logger.LogLevel, msg string) {
	s.Log(level, msg)
}
func (s *StatsAdapter) LogKV(level logger.LogLevel, msg string, keysAndValues ...interface{}) {
	s.Log(level, msg)
}
func (s *StatsAdapter) LogWithFields(level logger.LogLevel, msg string, fields map[string]interface{}) {
	s.Log(level, msg)
}

// 多行日志方法实现
func (s *StatsAdapter) DebugLines(lines ...string) {
	for _, line := range lines {
		s.Debug("%s", line)
	}
}

func (s *StatsAdapter) InfoLines(lines ...string) {
	for _, line := range lines {
		s.Info("%s", line)
	}
}

func (s *StatsAdapter) WarnLines(lines ...string) {
	for _, line := range lines {
		s.Warn("%s", line)
	}
}

func (s *StatsAdapter) ErrorLines(lines ...string) {
	for _, line := range lines {
		s.Error("%s", line)
	}
}
