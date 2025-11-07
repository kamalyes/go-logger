/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\examples\factory\main.go
 * @Description: 工厂模式使用示例
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"fmt"

	"github.com/kamalyes/go-logger"
)

func main() {
	// 1. 使用工厂创建日志器构建器
	builder := logger.NewLoggerBuilder()

	// 2. 链式配置日志器
	config := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithShowCaller(true).
		WithColorful(true).
		WithPrefix("[Factory] ")

	loggerInstance := builder.
		WithConfig(config).
		Build()

	fmt.Println("=== 工厂构建的日志器 ===")
	loggerInstance.Debug("调试信息 - 工厂模式")
	loggerInstance.Info("信息日志 - 工厂模式")
	loggerInstance.Warn("警告日志 - 工厂模式")
	loggerInstance.Error("错误日志 - 工厂模式")

	// 3. 创建带有钩子的日志器
	builder2 := logger.NewLoggerBuilder()
	
	// 添加控制台钩子
	consoleHookConfig := map[string]interface{}{
		"levels": []string{"error", "fatal"},
	}
	
	config2 := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithPrefix("[WithHooks] ")
	
	loggerWithHooks := builder2.
		WithConfig(config2).
		WithHook("console", consoleHookConfig).
		Build()

	fmt.Println("\n=== 带钩子的日志器 ===")
	loggerWithHooks.Info("这是普通信息，不会触发钩子")
	loggerWithHooks.Error("这是错误信息，会触发钩子")

	// 4. 创建带有中间件的日志器
	builder3 := logger.NewLoggerBuilder()
	
	// 添加认证中间件
	authMiddlewareConfig := map[string]interface{}{
		"required_role": "admin",
		"allowed_users": []string{"user1", "user2"},
	}

	config3 := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithPrefix("[WithMiddleware] ")

	loggerWithMiddleware := builder3.
		WithConfig(config3).
		WithMiddleware("auth", authMiddlewareConfig).
		Build()

	fmt.Println("\n=== 带中间件的日志器 ===")
	loggerWithMiddleware.Info("带有中间件处理的日志")

	// 5. 多个钩子和中间件组合
	complexBuilder := logger.NewLoggerBuilder()
	
	// 文件钩子配置
	fileHookConfig := map[string]interface{}{
		"filename": "app.log",
		"max_size": 10, // MB
		"levels":   []string{"info", "warn", "error", "fatal"},
	}

	// 度量中间件配置
	metricsMiddlewareConfig := map[string]interface{}{
		"enabled":     true,
		"sample_rate": 1.0,
	}

	complexConfig := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithShowCaller(true).
		WithPrefix("[Complex] ")

	complexLogger := complexBuilder.
		WithConfig(complexConfig).
		WithHook("console", consoleHookConfig).
		WithHook("file", fileHookConfig).
		WithMiddleware("auth", authMiddlewareConfig).
		WithMiddleware("metrics", metricsMiddlewareConfig).
		Build()

	fmt.Println("\n=== 复杂配置的日志器 ===")
	complexLogger.WithField("operation", "user_login").
		WithField("user_id", 12345).
		Info("复杂配置日志器测试")

	complexLogger.WithFields(map[string]interface{}{
		"module":   "authentication",
		"action":   "login_attempt",
		"success":  true,
		"duration": "150ms",
	}).Info("用户登录成功")

	// 6. 演示不同配置的区别
	fmt.Println("\n=== 配置对比 ===")

	// 简单配置
	simpleConfig := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithPrefix("[Simple] ")
	simple := logger.NewLoggerBuilder().
		WithConfig(simpleConfig).
		Build()
	simple.Info("简单配置的日志")

	// 详细配置
	detailedConfig := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithShowCaller(true).
		WithColorful(true).
		WithPrefix("[Detailed] ")
	detailed := logger.NewLoggerBuilder().
		WithConfig(detailedConfig).
		Build()
	detailed.Info("详细配置的日志")

	fmt.Println("\n工厂模式示例完成")
}