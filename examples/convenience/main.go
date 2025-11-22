/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-22 14:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 14:00:00
 * @FilePath: \go-logger\examples\convenience\main.go
 * @Description: 便利函数使用示例 - 演示NewUltraFast()、NewOptimized()和New()
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"fmt"
	"github.com/kamalyes/go-logger"
	"strings"
)

func main() {
	fmt.Println("🚀 Go Logger - 便利函数示例")
	fmt.Println(strings.Repeat("=", 40))

	// 演示三个便利函数的使用
	demonstrateConvenienceFunctions()

	fmt.Println()

	// 性能对比演示
	demonstratePerformanceComparison()

	fmt.Println()

	// 功能对比演示
	demonstrateFunctionComparison()
}

// 演示便利函数的基本使用
func demonstrateConvenienceFunctions() {
	fmt.Println("📋 1. 便利函数基本使用")
	fmt.Println(strings.Repeat("-", 30))

	fmt.Println("\n🔹 NewUltraFast() - 极致性能:")
	ultraLogger := logger.NewUltraFast()
	ultraLogger.Info("这是极致性能日志器 - 适用于高并发场景")
	ultraLogger.InfoKV("带键值的极速日志", "performance", "ultra")

	fmt.Println("\n🔹 NewOptimized() - 平衡性能:")
	optimizedLogger := logger.NewOptimized()
	optimizedLogger.Info("这是优化性能日志器 - 平衡性能与功能")
	optimizedLogger.WithField("type", "optimized").Info("带字段的优化日志")

	fmt.Println("\n🔹 New() - 完整功能:")
	standardLogger := logger.New()
	standardLogger.Info("这是标准功能日志器 - 提供完整企业级功能")
	standardLogger.WithField("feature", "complete").
		WithField("level", "enterprise").
		Info("带多字段的标准日志")
}

// 演示性能对比
func demonstratePerformanceComparison() {
	fmt.Println("📊 2. 性能特点对比")
	fmt.Println(strings.Repeat("-", 30))

	fmt.Println("\n性能特点:")
	fmt.Println("┌─────────────┬──────────┬────────┬─────────────┐")
	fmt.Println("│ 函数名称    │ 延迟     │ 分配   │ 适用场景    │")
	fmt.Println("├─────────────┼──────────┼────────┼─────────────┤")
	fmt.Println("│ UltraFast() │ ~7.56ns  │ 0      │ 高并发系统  │")
	fmt.Println("│ Optimized() │ ~22.85ns │ 1      │ 普通应用    │")
	fmt.Println("│ New()       │ ~130ns   │ 2      │ 企业应用    │")
	fmt.Println("└─────────────┴──────────┴────────┴─────────────┘")

	fmt.Println("\n💡 性能提示:")
	fmt.Println("  • UltraFast: 零分配设计，适合性能敏感场景")
	fmt.Println("  • Optimized: 智能缓存，平衡性能与功能")
	fmt.Println("  • New:       完整功能，支持所有企业级特性")
}

// 演示功能对比
func demonstrateFunctionComparison() {
	fmt.Println("🔧 3. 功能特性对比")
	fmt.Println(strings.Repeat("-", 30))

	// 创建三种日志器
	ultraLogger := logger.NewUltraFast()
	optimizedLogger := logger.NewOptimized()
	standardLogger := logger.New()

	fmt.Println("\n🔹 基础日志功能 (所有都支持):")
	ultraLogger.Info("UltraFast: 支持基础日志")
	optimizedLogger.Info("Optimized: 支持基础日志")
	standardLogger.Info("Standard: 支持基础日志")

	fmt.Println("\n🔹 结构化日志功能:")
	fmt.Println("UltraFast: 使用 InfoKV 方法")
	ultraLogger.InfoKV("键值对日志", "method", "InfoKV")

	fmt.Println("Optimized & Standard: 使用 WithField 方法")
	optimizedLogger.WithField("method", "WithField").Info("结构化日志")
	standardLogger.WithField("method", "WithField").
		WithField("feature", "rich").Info("多字段结构化日志")

	fmt.Println("\n🔹 链式配置 (运行时修改):")
	fmt.Println("Optimized & Standard: 支持链式配置")
	optimizedLogger.WithLevel(logger.DEBUG).Debug("运行时修改的调试日志")
	standardLogger.WithPrefix("[Runtime] ").Info("运行时添加前缀")

	fmt.Println("\n🔹 高级功能 (仅 Standard):")
	standardLogger.WithShowCaller(true).Info("显示调用者信息的日志")

	fmt.Println("\n✅ 所有便利函数功能演示完成!")
}
