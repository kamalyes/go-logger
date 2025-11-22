/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 13:43:27
 * @FilePath: \go-logger\examples\basic\main.go
 * @Description: 基本使用示例 - 演示logger.New()和基础功能
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
	fmt.Println("🚀 Go Logger - 基础使用示例")
	fmt.Println(strings.Repeat("=", 40))

	// 1. 最简单的使用方式 - New()函数
	demonstrateBasicUsage()

	fmt.Println()

	// 2. 链式配置使用
	demonstrateChainConfiguration()

	fmt.Println()

	// 3. 结构化日志
	demonstrateStructuredLogging()

	fmt.Println()

	// 4. 不同日志级别演示
	demonstrateLogLevels()

	fmt.Println()

	// 5. 全局日志器使用
	demonstrateGlobalLogger()

	fmt.Println()

	// 6. 不同日志方法演示
	demonstrateLogMethods()
}

// 演示基础使用
func demonstrateBasicUsage() {
	fmt.Println("📋 1. 基础使用演示")
	fmt.Println(strings.Repeat("-", 25))

	// 使用New()创建默认配置的日志器
	fmt.Println("\n🔹 使用logger.New():")
	log := logger.New()
	log.Info("欢迎使用 go-logger! 这是使用New()创建的日志器")
	log.Debug("这是调试信息（默认不显示，因为默认级别是INFO）")

	// 使用NewLogger()创建自定义配置
	fmt.Println("\n🔹 使用logger.NewLogger():")
	config := logger.NewLogConfig().
		WithLevel(logger.DEBUG).
		WithPrefix("[Custom] ")
	customLog := logger.NewLogger(config)
	customLog.Debug("现在可以看到调试信息了")
	customLog.Info("自定义配置的普通信息")
}

// 演示链式配置
func demonstrateChainConfiguration() {
	fmt.Println("🔧 2. 链式配置演示")
	fmt.Println(strings.Repeat("-", 25))

	// 方式1: 创建时链式配置
	fmt.Println("\n🔹 创建时链式配置:")
	log1 := logger.New().
		WithLevel(logger.DEBUG).
		WithPrefix("[Chain1] ").
		WithShowCaller(true).
		WithColorful(true)

	log1.Debug("链式配置的调试信息")
	log1.Info("带调用者信息的日志")

	// 方式2: 分步配置
	fmt.Println("\n🔹 分步配置:")
	log2 := logger.New()
	log2.WithLevel(logger.WARN)
	log2.WithPrefix("[Chain2] ")
	log2.WithShowCaller(false)

	log2.Debug("不会显示（级别低于WARN）")
	log2.Info("不会显示（级别低于WARN）")
	log2.Warn("这条警告会显示")
	log2.Error("这条错误会显示")
}

// 演示结构化日志
func demonstrateStructuredLogging() {
	fmt.Println("🏗️ 3. 结构化日志演示")
	fmt.Println(strings.Repeat("-", 25))

	log := logger.New().WithPrefix("[Struct] ")

	// 单个字段
	fmt.Println("\n🔹 单个字段:")
	log.WithField("user_id", 12345).Info("用户操作")
	log.WithField("operation", "login").WithField("ip", "192.168.1.1").Info("用户登录")

	// 多个字段
	fmt.Println("\n🔹 多个字段:")
	log.WithFields(map[string]interface{}{
		"service":    "user-api",
		"method":     "POST",
		"endpoint":   "/api/users/login",
		"status":     200,
		"duration":   "150ms",
		"request_id": "req-abc123",
		"user_agent": "Mozilla/5.0",
	}).Info("API请求处理完成")

	// 错误信息
	fmt.Println("\n🔹 错误信息:")
	err := fmt.Errorf("数据库连接失败")
	log.WithError(err).WithFields(map[string]interface{}{
		"database": "mysql",
		"host":     "localhost:3306",
		"retry":    3,
	}).Error("数据库操作失败")
}

// 演示不同日志级别
func demonstrateLogLevels() {
	fmt.Println("📊 4. 日志级别演示")
	fmt.Println(strings.Repeat("-", 25))

	// 设置为DEBUG级别，显示所有日志
	fmt.Println("\n🔹 DEBUG级别（显示所有）:")
	debugLog := logger.New().WithLevel(logger.DEBUG).WithPrefix("[DEBUG] ")
	debugLog.Debug("🐛 调试信息 - 用于开发阶段")
	debugLog.Info("ℹ️ 普通信息 - 一般运行信息")
	debugLog.Warn("⚠️ 警告信息 - 需要注意但不影响运行")
	debugLog.Error("❌ 错误信息 - 发生了错误但程序可以继续")

	// 设置为INFO级别
	fmt.Println("\n🔹 INFO级别（不显示DEBUG）:")
	infoLog := logger.New().WithLevel(logger.INFO).WithPrefix("[INFO] ")
	infoLog.Debug("不会显示这条调试信息")
	infoLog.Info("显示这条普通信息")
	infoLog.Warn("显示这条警告信息")
	infoLog.Error("显示这条错误信息")

	// 设置为ERROR级别
	fmt.Println("\n🔹 ERROR级别（只显示错误）:")
	errorLog := logger.New().WithLevel(logger.ERROR).WithPrefix("[ERROR] ")
	errorLog.Debug("不显示")
	errorLog.Info("不显示")
	errorLog.Warn("不显示")
	errorLog.Error("只显示错误信息")

	// 级别检查
	fmt.Println("\n🔹 级别检查:")
	log := logger.New().WithLevel(logger.WARN)
	fmt.Printf("当前级别: %v\n", log.GetLevel())
	fmt.Printf("DEBUG启用: %v\n", log.IsLevelEnabled(logger.DEBUG))
	fmt.Printf("INFO启用: %v\n", log.IsLevelEnabled(logger.INFO))
	fmt.Printf("WARN启用: %v\n", log.IsLevelEnabled(logger.WARN))
	fmt.Printf("ERROR启用: %v\n", log.IsLevelEnabled(logger.ERROR))
}

// 演示全局日志器
func demonstrateGlobalLogger() {
	fmt.Println("🌍 5. 全局日志器演示")
	fmt.Println(strings.Repeat("-", 25))

	// 使用全局日志方法
	fmt.Println("\n🔹 全局日志方法:")
	logger.Info("使用全局Info方法")
	logger.Warn("使用全局Warn方法")

	// 配置全局日志器
	fmt.Println("\n🔹 配置全局日志器:")
	logger.SetGlobalLevel(logger.DEBUG)
	logger.SetGlobalShowCaller(true)
	logger.Debug("配置后的全局调试信息")

	// 全局结构化日志
	fmt.Println("\n🔹 全局结构化日志:")
	logger.WithField("component", "main").
		WithField("version", "1.0.0").
		Info("应用启动完成")

	logger.WithFields(map[string]interface{}{
		"memory_usage": "45MB",
		"goroutines":   12,
		"uptime":       "30s",
	}).Info("系统状态")

	// 获取全局日志器
	fmt.Println("\n🔹 获取全局日志器:")
	globalLogger := logger.GetGlobalLogger()
	globalLogger.WithPrefix("[Global] ").Info("通过全局日志器记录")
}

// 演示不同的日志方法
func demonstrateLogMethods() {
	fmt.Println("🔧 6. 不同日志方法演示")
	fmt.Println(strings.Repeat("-", 25))

	log := logger.New().WithPrefix("[Methods] ")

	// Printf风格
	fmt.Println("\n🔹 Printf风格:")
	log.Infof("用户%s登录成功，ID: %d", "张三", 12345)
	log.Warnf("磁盘使用率达到%.1f%%", 85.6)

	// 纯文本方法
	fmt.Println("\n🔹 纯文本方法:")
	log.InfoMsg("这是一条纯文本信息")
	log.WarnMsg("这是一条纯文本警告")

	// 兼容标准log
	fmt.Println("\n🔹 兼容标准log:")
	log.Print("兼容Print方法")
	log.Printf("兼容Printf方法：%s", "格式化文本")
	log.Println("兼容Println方法")
}
