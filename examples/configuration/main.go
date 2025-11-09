/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\examples\configuration\main.go
 * @Description: 配置使用示例
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kamalyes/go-logger"
)

func main() {
	// 1. 默认配置
	fmt.Println("=== 默认配置 ===")
	defaultLogger := logger.NewLogger(logger.DefaultConfig())
	defaultLogger.Info("使用默认配置的日志")

	// 2. 自定义基本配置
	fmt.Println("\n=== 自定义基本配置 ===")
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

	// 3. 时间格式配置
	fmt.Println("\n=== 时间格式配置 ===")
	timeFormats := []string{
		"15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"Jan 2 15:04:05",
		"2006/01/02 15:04:05.000",
	}

	for i, format := range timeFormats {
		config := logger.DefaultConfig().
			WithTimeFormat(format).
			WithPrefix(fmt.Sprintf("[Format%d] ", i+1))

		timeLogger := logger.NewLogger(config)
		timeLogger.Info("时间格式示例: %s", format)
	}

	// 4. 输出目标配置
	fmt.Println("\n=== 输出目标配置 ===")

	// 标准输出
	stdoutConfig := logger.DefaultConfig().
		WithOutput(os.Stdout).
		WithPrefix("[STDOUT] ")
	stdoutLogger := logger.NewLogger(stdoutConfig)
	stdoutLogger.Info("输出到标准输出")

	// 标准错误
	stderrConfig := logger.DefaultConfig().
		WithOutput(os.Stderr).
		WithPrefix("[STDERR] ")
	stderrLogger := logger.NewLogger(stderrConfig)
	stderrLogger.Error("输出到标准错误")

	// 字符串缓冲区
	var buffer strings.Builder
	bufferConfig := logger.DefaultConfig().
		WithOutput(&buffer).
		WithPrefix("[BUFFER] ")
	bufferLogger := logger.NewLogger(bufferConfig)
	bufferLogger.Info("输出到缓冲区")
	fmt.Printf("缓冲区内容: %s", buffer.String())

	// 5. 级别配置演示
	fmt.Println("\n=== 级别配置演示 ===")
	levels := []logger.LogLevel{
		logger.DEBUG,
		logger.INFO,
		logger.WARN,
		logger.ERROR,
		logger.FATAL,
	}

	for _, level := range levels {
		config := logger.DefaultConfig().
			WithLevel(level).
			WithPrefix(fmt.Sprintf("[%s] ", level))

		levelLogger := logger.NewLogger(config)
		fmt.Printf("设置级别为 %s:\n", level)
		levelLogger.Debug("  这是调试信息")
		levelLogger.Info("  这是普通信息")
		levelLogger.Warn("  这是警告信息")
		levelLogger.Error("  这是错误信息")
		fmt.Println()
	}

	// 6. 调用者信息配置
	fmt.Println("=== 调用者信息配置 ===")
	
	// 显示调用者信息
	callerConfig := logger.DefaultConfig().
		WithShowCaller(true).
		WithPrefix("[WithCaller] ")
	callerLogger := logger.NewLogger(callerConfig)
	demonstrateCaller(callerLogger)

	// 不显示调用者信息
	noCallerConfig := logger.DefaultConfig().
		WithShowCaller(false).
		WithPrefix("[NoCaller] ")
	noCallerLogger := logger.NewLogger(noCallerConfig)
	demonstrateCaller(noCallerLogger)

	// 7. 彩色输出配置
	fmt.Println("\n=== 彩色输出配置 ===")
	
	// 彩色输出
	colorConfig := logger.DefaultConfig().
		WithColorful(true).
		WithPrefix("[Color] ")
	colorLogger := logger.NewLogger(colorConfig)
	colorLogger.Debug("彩色调试信息")
	colorLogger.Info("彩色普通信息")
	colorLogger.Warn("彩色警告信息")
	colorLogger.Error("彩色错误信息")

	// 无彩色输出
	noColorConfig := logger.DefaultConfig().
		WithColorful(false).
		WithPrefix("[NoColor] ")
	noColorLogger := logger.NewLogger(noColorConfig)
	noColorLogger.Debug("无彩色调试信息")
	noColorLogger.Info("无彩色普通信息")
	noColorLogger.Warn("无彩色警告信息")
	noColorLogger.Error("无彩色错误信息")

	// 8. 配置验证
	fmt.Println("\n=== 配置验证 ===")
	validConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithTimeFormat("2006-01-02 15:04:05").
		WithPrefix("[Valid] ")

	if err := validConfig.Validate(); err != nil {
		fmt.Printf("配置验证失败: %v\n", err)
	} else {
		fmt.Println("配置验证通过")
		validLogger := logger.NewLogger(validConfig)
		validLogger.Info("验证通过的配置日志")
	}

	// 9. 配置克隆
	fmt.Println("\n=== 配置克隆 ===")
	originalConfig := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithPrefix("[Original] ")

	// 克隆配置并修改
	clonedConfig := originalConfig.Clone().
		WithLevel(logger.WARN).
		WithPrefix("[Cloned] ")

	originalLogger := logger.NewLogger(originalConfig)
	clonedLogger := logger.NewLogger(clonedConfig)

	fmt.Println("原始配置输出:")
	originalLogger.Debug("调试信息")
	originalLogger.Info("普通信息")
	originalLogger.Warn("警告信息")

	fmt.Println("克隆配置输出:")
	clonedLogger.Debug("调试信息 (不会显示)")
	clonedLogger.Info("普通信息 (不会显示)")
	clonedLogger.Warn("警告信息")

	// 8. 演示不同场景的用法
	fmt.Println("\n=== 不同场景配置演示 ===")
	demonstrateUsageScenarios()

	fmt.Println("\n配置示例完成")
}

func demonstrateCaller(log *logger.Logger) {
	log.Info("这条日志用于演示调用者信息")
}

// 演示配置的不同用途
func demonstrateUsageScenarios() {
	// 开发环境配置
	devConfig := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithShowCaller(true).
		WithColorful(true).
		WithPrefix("[DEV] ")

	// 生产环境配置
	prodConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithShowCaller(false).
		WithColorful(false).
		WithTimeFormat("2006-01-02T15:04:05Z07:00").
		WithPrefix("[PROD] ")

	// 测试环境配置
	testConfig := logger.DefaultConfig().
		WithLevel(logger.WARN).
		WithShowCaller(false).
		WithColorful(false).
		WithPrefix("[TEST] ")

	devLogger := logger.NewLogger(devConfig)
	prodLogger := logger.NewLogger(prodConfig)
	testLogger := logger.NewLogger(testConfig)

	message := "配置演示消息"
	devLogger.Info(message)
	prodLogger.Info(message)
	testLogger.Info(message)
}