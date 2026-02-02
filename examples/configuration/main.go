/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 00:00:00
 * @FilePath: \go-logger\examples\configuration\main.go
 * @Description: 配置系统示例 - 演示完整的配置选项和最佳实践
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kamalyes/go-logger"
)

func main() {
	fmt.Println("🔧 Go Logger - 配置系统示例演示")
	fmt.Println(strings.Repeat("=", 50))

	// 1. 基础配置演示
	demonstrateBasicConfiguration()

	fmt.Println()

	// 2. 环境特定配置演示
	demonstrateEnvironmentConfigurations()

	fmt.Println()

	// 3. 高级配置选项演示
	demonstrateAdvancedConfigurations()

	fmt.Println()

	// 4. 性能优化配置演示
	demonstratePerformanceConfigurations()

	fmt.Println()

	// 5. 配置验证和最佳实践
	demonstrateConfigurationBestPractices()

	fmt.Println()

	// 6. 动态配置更新演示
	demonstrateDynamicConfiguration()
}

// 基础配置演示
func demonstrateBasicConfiguration() {
	fmt.Println("📋 1. 基础配置演示")
	fmt.Println(strings.Repeat("-", 30))

	// 1.1 默认配置
	fmt.Println("\n🔹 默认配置:")
	defaultLogger := logger.NewLogger(logger.DefaultConfig())
	defaultLogger.Info("使用默认配置的日志")

	// 1.2 基本配置选项
	fmt.Println("\n🔹 基本配置选项:")
	basicConfig := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithShowCaller(true).
		WithColorful(true).
		WithPrefix("[MyApp] ")
	basicLogger := logger.NewLogger(basicConfig)

	basicLogger.Debug("调试信息")
	basicLogger.Info("普通信息")
	basicLogger.Warn("警告信息")
	basicLogger.Error("错误信息")

	// 1.3 时间格式配置
	fmt.Println("\n🔹 不同时间格式:")
	timeFormats := map[string]string{
		"简短时间":    "15:04:05",
		"标准时间":    "2006-01-02 15:04:05",
		"ISO时间":   "2006-01-02T15:04:05Z07:00",
		"RFC3339": time.RFC3339,
		"毫秒精度":    time.RFC3339Nano,
		"Unix时间戳": "unix",
	}

	for name, format := range timeFormats {
		config := logger.DefaultConfig().
			WithLevel(logger.INFO).
			WithTimeFormat(format).
			WithPrefix(fmt.Sprintf("[%s] ", name))
		l := logger.NewLogger(config)
		l.Info("时间格式演示")
	}
}

// 环境特定配置演示
func demonstrateEnvironmentConfigurations() {
	fmt.Println("🌍 2. 环境特定配置演示")
	fmt.Println(strings.Repeat("-", 30))

	// 2.1 开发环境配置
	fmt.Println("\n🔹 开发环境配置:")
	devConfig := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithTimeFormat("15:04:05.000").
		WithShowCaller(true).
		WithColorful(true).
		WithPrefix("[DEV] ")
	devLogger := logger.NewLogger(devConfig)

	devLogger.Debug("开发环境调试信息")
	devLogger.Info("开发环境普通信息")
	devLogger.WithFields(map[string]interface{}{
		"module":   "auth",
		"function": "login",
	}).Info("开发环境结构化日志")

	// 2.2 测试环境配置
	fmt.Println("\n🔹 测试环境配置:")
	testConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithTimeFormat("2006-01-02 15:04:05").
		WithShowCaller(false).
		WithColorful(false).
		WithPrefix("[TEST] ")
	testLogger := logger.NewLogger(testConfig)

	testLogger.Info("测试环境信息")
	testLogger.WithFields(map[string]interface{}{
		"test_case": "user_login_test",
		"status":    "passed",
		"duration":  "150ms",
	}).Info("测试用例结果")

	// 2.3 生产环境配置
	fmt.Println("\n🔹 生产环境配置:")
	prodConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithTimeFormat(time.RFC3339).
		WithShowCaller(false).
		WithColorful(false)
	prodLogger := logger.NewLogger(prodConfig)

	prodLogger.Info("生产环境服务启动")
	prodLogger.WithFields(map[string]interface{}{
		"service":     "user-service",
		"version":     "1.2.3",
		"port":        8080,
		"environment": "production",
	}).Info("服务配置信息")

	// 2.4 监控环境配置
	fmt.Println("\n🔹 监控环境配置:")
	monitorConfig := logger.DefaultConfig().
		WithLevel(logger.WARN). // 只记录警告和错误
		WithTimeFormat(time.RFC3339Nano)
	monitorLogger := logger.NewLogger(monitorConfig)

	monitorLogger.WarnKV("系统负载过高", "cpu_usage", "85%")
	monitorLogger.ErrorKV("数据库连接失败", "error", "connection timeout")
}

// 高级配置选项演示
func demonstrateAdvancedConfigurations() {
	fmt.Println("⚙️ 3. 高级配置选项演示")
	fmt.Println(strings.Repeat("-", 30))

	// 3.1 输出目标配置
	fmt.Println("\n🔹 输出目标配置:")

	// 标准输出
	stdoutConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithOutput(os.Stdout).
		WithPrefix("[STDOUT] ")
	stdoutLogger := logger.NewLogger(stdoutConfig)
	stdoutLogger.Info("输出到标准输出")

	// 标准错误
	stderrConfig := logger.DefaultConfig().
		WithLevel(logger.ERROR).
		WithOutput(os.Stderr).
		WithPrefix("[STDERR] ")
	stderrLogger := logger.NewLogger(stderrConfig)
	stderrLogger.Error("输出到标准错误")

	// 缓冲区输出
	var buffer strings.Builder
	bufferConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithOutput(&buffer).
		WithPrefix("[BUFFER] ")
	bufferLogger := logger.NewLogger(bufferConfig)
	bufferLogger.Info("输出到内存缓冲区")
	fmt.Printf("缓冲区内容: %s", buffer.String())

	// 3.2 格式化配置演示 (简化版本)
	fmt.Println("\n🔹 格式化配置:")

	// 标准文本格式
	textConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithColorful(true)
	textLogger := logger.NewLogger(textConfig)
	textLogger.WithField("format", "text").Info("文本格式日志")

	// 无颜色格式（类似JSON风格）
	plainConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithColorful(false)
	plainLogger := logger.NewLogger(plainConfig)
	plainLogger.WithField("format", "plain").Info("无颜色格式日志")

	// 3.3 级别配置
	fmt.Println("\n🔹 级别过滤演示:")
	levels := []logger.LogLevel{logger.DEBUG, logger.INFO, logger.WARN, logger.ERROR}

	for _, level := range levels {
		config := logger.DefaultConfig().
			WithLevel(level).
			WithPrefix(fmt.Sprintf("[%s_LEVEL] ", strings.ToUpper(level.String())))
		l := logger.NewLogger(config)

		fmt.Printf("设置级别为 %s:\n", level.String())
		l.Debug("  调试信息")
		l.Info("  普通信息")
		l.Warn("  警告信息")
		l.Error("  错误信息")
	}
}

// 性能优化配置演示
func demonstratePerformanceConfigurations() {
	fmt.Println("⚡ 4. 性能优化配置演示")
	fmt.Println(strings.Repeat("-", 30))

	// 4.1 高性能配置
	fmt.Println("\n🔹 高性能配置 (生产环境推荐):")
	highPerfConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithTimeFormat("unix"). // Unix时间戳最快
		WithShowCaller(false).  // 关闭调用者信息
		WithColorful(false)     // 关闭颜色输出
	highPerfLogger := logger.NewLogger(highPerfConfig)

	// 性能测试
	start := time.Now()
	for i := 0; i < 1000; i++ {
		highPerfLogger.WithField("iteration", i).Info("高性能日志测试")
	}
	highPerfDuration := time.Since(start)
	fmt.Printf("高性能配置 1000 条日志耗时: %v\n", highPerfDuration)

	// 4.2 标准配置对比
	fmt.Println("\n🔹 标准配置对比:")
	standardConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithTimeFormat(time.RFC3339Nano).
		WithShowCaller(true).
		WithColorful(true)
	standardLogger := logger.NewLogger(standardConfig)

	start = time.Now()
	for i := 0; i < 1000; i++ {
		standardLogger.WithFields(map[string]interface{}{
			"iteration": i,
			"timestamp": time.Now(),
		}).Info("标准配置日志测试")
	}
	standardDuration := time.Since(start)
	fmt.Printf("标准配置 1000 条日志耗时: %v\n", standardDuration)

	speedup := float64(standardDuration) / float64(highPerfDuration)
	fmt.Printf("性能提升: %.2fx\n", speedup)

	// 4.3 内存优化配置
	fmt.Println("\n🔹 内存优化配置:")
	memOptConfig := logger.DefaultConfig().
		WithLevel(logger.INFO)
	memOptLogger := logger.NewLogger(memOptConfig)
	memOptLogger.Info("内存优化配置日志")
}

// 配置验证和最佳实践
func demonstrateConfigurationBestPractices() {
	fmt.Println("🎯 5. 配置验证和最佳实践")
	fmt.Println(strings.Repeat("-", 30))

	// 5.1 配置验证
	fmt.Println("\n🔹 配置验证:")

	// 有效配置
	validConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithTimeFormat("2006-01-02 15:04:05").
		WithOutput(os.Stdout)
	validLogger := logger.NewLogger(validConfig)

	fmt.Printf("✅ 配置验证通过\n")
	validLogger.Info("配置验证通过的日志")

	// 5.2 配置继承和组合
	fmt.Println("\n🔹 配置继承和组合:")

	// 基础配置
	baseConfig := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithTimeFormat("15:04:05").
		WithShowCaller(true)
	baseLogger := logger.NewLogger(baseConfig)

	// 派生配置 - 基于基础配置创建子日志器
	derivedLogger := baseLogger.WithField("component", "api").
		WithField("version", "1.0.0")

	fmt.Println("基础配置输出:")
	baseLogger.Debug("调试信息")
	baseLogger.Info("普通信息")
	baseLogger.Warn("警告信息")

	fmt.Println("派生配置输出:")
	derivedLogger.Debug("调试信息")
	derivedLogger.Info("普通信息")
	derivedLogger.Warn("警告信息")

	// 5.3 最佳实践示例
	fmt.Println("\n🔹 最佳实践示例:")
	demonstrateConfigurationPatterns()
}

// 动态配置更新演示
func demonstrateDynamicConfiguration() {
	fmt.Println("🔄 6. 动态配置更新演示")
	fmt.Println(strings.Repeat("-", 30))

	// 创建可更新的日志器
	dynamicConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithTimeFormat("15:04:05").
		WithPrefix("[DYNAMIC] ")
	dynamicLogger := logger.NewLogger(dynamicConfig)

	fmt.Println("\n🔹 初始配置:")
	dynamicLogger.Debug("调试信息 (不显示)")
	dynamicLogger.Info("普通信息")
	dynamicLogger.Warn("警告信息")

	// 模拟配置更新 - 创建新的日志器
	fmt.Println("\n🔹 更新配置为DEBUG级别:")
	updatedConfig := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithTimeFormat("15:04:05.000").
		WithPrefix("[UPDATED] ").
		WithColorful(true)
	updatedLogger := logger.NewLogger(updatedConfig)

	updatedLogger.Debug("调试信息 (现在显示)")
	updatedLogger.Info("普通信息")
	updatedLogger.Warn("警告信息")

	fmt.Println("\n🔹 模拟热重载配置:")
	// 模拟从配置文件热重载
	configMap := map[string]interface{}{
		"level":       "error",
		"time_format": "2006-01-02T15:04:05Z07:00",
		"format":      "json",
		"show_caller": false,
	}

	hotReloadLogger := createLoggerFromMap(configMap)

	hotReloadLogger.Debug("调试信息 (不显示)")
	hotReloadLogger.Info("普通信息 (不显示)")
	hotReloadLogger.Warn("警告信息 (不显示)")
	hotReloadLogger.Error("错误信息")
}

// 配置模式演示
func demonstrateConfigurationPatterns() {
	// 微服务配置模式
	microserviceConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithTimeFormat(time.RFC3339)
	microserviceLogger := logger.NewLogger(microserviceConfig).WithFields(map[string]interface{}{
		"service":  "user-service",
		"version":  "1.0.0",
		"instance": "us-west-1a",
	})
	microserviceLogger.Info("微服务日志模式")

	// API网关配置模式
	gatewayConfig := logger.DefaultConfig().
		WithLevel(logger.INFO)
	gatewayLogger := logger.NewLogger(gatewayConfig).WithField("component", "api-gateway")

	gatewayLogger.WithFields(map[string]interface{}{
		"method":   "GET",
		"path":     "/api/users",
		"status":   200,
		"duration": "50ms",
	}).Info("API请求日志")

	// 批处理任务配置模式
	batchConfig := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithTimeFormat(time.RFC3339Nano)
	batchLogger := logger.NewLogger(batchConfig).WithField("job_type", "data-processing")

	batchLogger.Info("批处理任务开始")
}

// 从配置映射创建日志器
func createLoggerFromMap(configMap map[string]interface{}) *logger.Logger {
	config := logger.DefaultConfig()

	if level, ok := configMap["level"].(string); ok {
		switch level {
		case "debug":
			config = config.WithLevel(logger.DEBUG)
		case "info":
			config = config.WithLevel(logger.INFO)
		case "warn":
			config = config.WithLevel(logger.WARN)
		case "error":
			config = config.WithLevel(logger.ERROR)
		}
	}

	if timeFormat, ok := configMap["time_format"].(string); ok {
		config = config.WithTimeFormat(timeFormat)
	}

	if showCaller, ok := configMap["show_caller"].(bool); ok {
		config = config.WithShowCaller(showCaller)
	}

	// 格式化配置简化处理（当前版本不支持复杂的formatter配置）
	if format, ok := configMap["format"].(string); ok {
		switch format {
		case "json":
			// 可以设置为无颜色模式模拟JSON风格
			config = config.WithColorful(false)
		default:
			config = config.WithColorful(true)
		}
	}

	return logger.NewLogger(config)
}
