/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-09 18:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 18:00:00
 * @FilePath: \go-logger\examples\enhanced\main.go
 * @Description: 增强接口功能示例，展示新增的上下文感知、键值对日志、多框架兼容等特性
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kamalyes/go-logger"
)

func main() {
	fmt.Println("=== Go Logger 增强接口功能演示 ===")

	// 创建日志器
	config := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithShowCaller(true).
		WithColorful(true).
		WithPrefix("[Enhanced] ")

	log := logger.NewLogger(config)

	// 1. 纯文本消息日志
	demonstrateTextMessages(log)

	// 2. 上下文感知日志
	demonstrateContextAware(log)

	// 3. 结构化键值对日志
	demonstrateStructuredLogging(log)

	// 4. 原始日志条目方法
	demonstrateRawLogging(log)

	// 5. 多框架兼容性
	demonstrateFrameworkCompatibility(log)

	// 6. 标准库兼容性
	demonstrateStandardLibCompatibility(log)

	// 7. 高级用法示例
	demonstrateAdvancedUsage(log)

	fmt.Println("=== 演示完成 ===")
}

// demonstrateTextMessages 演示纯文本消息日志
func demonstrateTextMessages(log logger.ILogger) {
	fmt.Println("1. === 纯文本消息日志演示 ===")
	
	log.DebugMsg("这是一条调试消息")
	log.InfoMsg("这是一条信息消息")
	log.WarnMsg("这是一条警告消息")
	log.ErrorMsg("这是一条错误消息")
	
	fmt.Println()
}

// demonstrateContextAware 演示上下文感知日志
func demonstrateContextAware(log logger.ILogger) {
	fmt.Println("2. === 上下文感知日志演示 ===")
	
	// 创建带值的上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-12345")
	ctx = context.WithValue(ctx, "user_id", 67890)
	ctx = context.WithValue(ctx, "trace_id", "trace-abc123")
	
	// 使用带上下文的日志方法
	log.DebugContext(ctx, "开始处理用户请求")
	log.InfoContext(ctx, "用户 %d 执行了操作: %s", 67890, "login")
	log.WarnContext(ctx, "检测到可疑活动，IP: %s", "192.168.1.100")
	log.ErrorContext(ctx, "处理请求时发生错误: %v", errors.New("数据库连接超时"))
	
	// 创建带上下文的日志器
	ctxLogger := log.WithContext(ctx)
	ctxLogger.Info("这个日志器携带了请求上下文")
	ctxLogger.Error("上下文感知的错误日志")
	
	fmt.Println()
}

// demonstrateStructuredLogging 演示结构化键值对日志
func demonstrateStructuredLogging(log logger.ILogger) {
	fmt.Println("3. === 结构化键值对日志演示 ===")
	
	// Zap风格的键值对日志
	log.InfoKV("用户登录成功",
		"user_id", 12345,
		"username", "john_doe",
		"ip_address", "192.168.1.100",
		"login_time", time.Now().Format(time.RFC3339),
		"user_agent", "Mozilla/5.0...",
	)
	
	log.WarnKV("数据库连接性能警告",
		"database", "user_db",
		"host", "db.example.com",
		"latency_ms", 450,
		"threshold_ms", 200,
		"query", "SELECT * FROM users WHERE id = ?",
	)
	
	log.ErrorKV("支付处理失败",
		"order_id", 98765,
		"amount", 299.99,
		"currency", "USD",
		"payment_method", "credit_card",
		"error_code", "CARD_DECLINED",
		"retry_count", 3,
	)
	
	// 调试级别的键值对日志
	log.DebugKV("缓存操作",
		"operation", "GET",
		"key", "user:12345:profile",
		"cache_hit", true,
		"ttl_seconds", 3600,
	)
	
	fmt.Println()
}

// demonstrateRawLogging 演示原始日志条目方法
func demonstrateRawLogging(log logger.ILogger) {
	fmt.Println("4. === 原始日志条目方法演示 ===")
	
	// 基础日志方法
	log.Log(logger.DEBUG, "原始调试日志")
	log.Log(logger.INFO, "原始信息日志")
	log.Log(logger.WARN, "原始警告日志")
	log.Log(logger.ERROR, "原始错误日志")
	
	// 带上下文的原始日志
	ctx := context.Background()
	ctx = context.WithValue(ctx, "component", "auth-service")
	
	log.LogContext(ctx, logger.INFO, "带上下文的原始日志")
	log.LogContext(ctx, logger.ERROR, "带上下文的错误日志")
	
	// 带键值对的原始日志
	log.LogKV(logger.INFO, "系统状态检查",
		"status", "healthy",
		"uptime_seconds", 86400,
		"memory_usage_mb", 256,
		"cpu_usage_percent", 15.5,
	)
	
	// 带字段映射的原始日志
	fields := map[string]interface{}{
		"service": "order-service",
		"version": "v2.1.0",
		"environment": "production",
		"region": "us-east-1",
	}
	log.LogWithFields(logger.INFO, "服务启动完成", fields)
	
	fmt.Println()
}

// demonstrateFrameworkCompatibility 演示多框架兼容性
func demonstrateFrameworkCompatibility(log logger.ILogger) {
	fmt.Println("5. === 多框架兼容性演示 ===")
	
	// === Zap 风格 ===
	fmt.Println("Zap 风格键值对日志:")
	log.InfoKV("订单创建",
		"action", "create_order",
		"user_id", 12345,
		"order_id", 67890,
		"amount", 299.99,
		"currency", "USD",
		"timestamp", time.Now().Unix(),
	)
	
	// === Logrus 风格 ===
	fmt.Println("\nLogrus 风格字段日志:")
	log.WithFields(map[string]interface{}{
		"component": "payment-processor",
		"method": "stripe",
		"transaction_id": "txn_abc123",
		"amount": 299.99,
		"currency": "USD",
	}).Info("支付处理开始")
	
	log.WithField("user_id", 12345).
		WithField("action", "password_reset").
		WithField("ip", "192.168.1.50").
		Warn("密码重置请求")
	
	// === slog 风格 ===
	fmt.Println("\nslog 风格上下文日志:")
	ctx := context.Background()
	log.InfoContext(ctx, "处理HTTP请求: %s %s", "POST", "/api/v1/orders")
	log.WarnContext(ctx, "请求处理缓慢: %dms", 1200)
	
	// === Zerolog 风格 ===
	fmt.Println("\nZerolog 风格链式调用:")
	log.WithField("service", "notification").
		WithField("channel", "email").
		WithField("recipient", "user@example.com").
		Info("发送通知邮件")
	
	fmt.Println()
}

// demonstrateStandardLibCompatibility 演示标准库兼容性
func demonstrateStandardLibCompatibility(log logger.ILogger) {
	fmt.Println("6. === 标准库兼容性演示 ===")
	
	// 兼容标准库 log 包的方法
	log.Print("标准库风格的 Print 方法")
	log.Printf("标准库风格的 Printf 方法: 用户 %d 订单 %d", 12345, 67890)
	log.Println("标准库风格的 Println 方法")
	
	fmt.Println()
}

// demonstrateAdvancedUsage 演示高级用法
func demonstrateAdvancedUsage(log logger.ILogger) {
	fmt.Println("7. === 高级用法演示 ===")
	
	// 模拟一个复杂的业务场景
	simulateOrderProcessing(log)
	
	// 演示错误处理
	simulateErrorHandling(log)
	
	// 演示性能监控
	simulatePerformanceLogging(log)
	
	fmt.Println()
}

// simulateOrderProcessing 模拟订单处理流程
func simulateOrderProcessing(log logger.ILogger) {
	fmt.Println("模拟订单处理流程:")
	
	// 创建请求上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-order-001")
	ctx = context.WithValue(ctx, "user_id", 12345)
	
	// 创建带上下文的日志器
	orderLogger := log.WithContext(ctx).WithField("component", "order-service")
	
	// 1. 接收订单
	orderLogger.InfoKV("订单接收",
		"order_id", "order-12345",
		"user_id", 12345,
		"total_amount", 599.99,
		"item_count", 3,
	)
	
	// 2. 验证库存
	orderLogger.DebugKV("检查库存",
		"action", "inventory_check",
		"items", []string{"item-001", "item-002", "item-003"},
		"check_duration_ms", 45,
	)
	
	// 3. 处理支付
	orderLogger.WithField("payment_method", "credit_card").
		InfoKV("处理支付",
			"amount", 599.99,
			"currency", "USD",
			"gateway", "stripe",
		)
	
	// 4. 订单完成
	orderLogger.InfoMsg("订单处理完成")
}

// simulateErrorHandling 模拟错误处理
func simulateErrorHandling(log logger.ILogger) {
	fmt.Println("\n模拟错误处理:")
	
	// 模拟数据库错误
	dbError := errors.New("connection timeout after 30s")
	log.WithError(dbError).
		ErrorKV("数据库操作失败",
			"operation", "INSERT",
			"table", "orders",
			"timeout_seconds", 30,
			"retry_count", 3,
		)
	
	// 模拟API错误
	log.ErrorContext(context.Background(), "外部API调用失败: %s", "payment gateway timeout")
	
	// 模拟业务逻辑错误
	log.WarnKV("业务规则违规",
		"rule", "max_order_amount",
		"user_id", 12345,
		"attempted_amount", 10000.00,
		"limit", 5000.00,
		"action", "order_rejected",
	)
}

// simulatePerformanceLogging 模拟性能监控日志
func simulatePerformanceLogging(log logger.ILogger) {
	fmt.Println("\n模拟性能监控:")
	
	start := time.Now()
	
	// 模拟一些操作
	time.Sleep(time.Millisecond * 50)
	
	duration := time.Since(start)
	
	log.InfoKV("操作性能统计",
		"operation", "order_processing",
		"duration_ms", duration.Milliseconds(),
		"threshold_ms", 100,
		"status", func() string {
			if duration.Milliseconds() > 100 {
				return "slow"
			}
			return "normal"
		}(),
		"memory_mb", 128,
		"cpu_percent", 15.5,
	)
	
	// API响应时间监控
	log.DebugKV("API响应时间",
		"endpoint", "/api/v1/orders",
		"method", "POST",
		"response_time_ms", 85,
		"status_code", 201,
		"response_size_bytes", 1024,
	)
}