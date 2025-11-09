/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-09 19:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 19:00:00
 * @FilePath: \go-logger\examples\compatibility\main.go
 * @Description: 多框架兼容性示例，展示如何使用 go-logger 模拟不同日志框架的API风格
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
	fmt.Println("=== Go Logger 多框架兼容性演示 ===")

	// 创建基础日志器
	config := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithShowCaller(true).
		WithColorful(true)

	baseLogger := logger.NewLogger(config)

	// 1. Zap 风格演示
	demonstrateZapStyle(baseLogger)

	// 2. Logrus 风格演示
	demonstrateLogrusStyle(baseLogger)

	// 3. slog 风格演示
	demonstrateSlogStyle(baseLogger)

	// 4. Zerolog 风格演示
	demonstrateZerologStyle(baseLogger)

	// 5. 标准库 log 风格演示
	demonstrateStdLogStyle(baseLogger)

	// 6. 混合使用演示
	demonstrateMixedUsage(baseLogger)

	fmt.Println("\n=== 兼容性演示完成 ===")
}

// demonstrateZapStyle 演示 Zap 风格的日志记录
func demonstrateZapStyle(log logger.ILogger) {
	fmt.Println("1. === Zap 风格演示 ===")
	fmt.Println("特点：结构化日志、强类型、高性能、键值对格式")

	zapLogger := log.WithField("framework", "zap-style")

	// Zap 风格的基础日志记录
	zapLogger.InfoKV("服务启动",
		"service", "user-api",
		"version", "v1.2.3",
		"port", 8080,
		"env", "production",
	)

	// 错误日志记录（Zap 风格）
	err := errors.New("database connection failed")
	zapLogger.ErrorKV("数据库连接失败",
		"error", err.Error(),
		"database", "postgres",
		"host", "db.example.com",
		"port", 5432,
		"timeout_seconds", 30,
		"retry_count", 3,
	)

	// 性能监控日志
	start := time.Now()
	time.Sleep(time.Millisecond * 50) // 模拟操作
	duration := time.Since(start)

	zapLogger.InfoKV("API请求处理完成",
		"method", "POST",
		"path", "/api/users",
		"status_code", 201,
		"duration_ms", duration.Milliseconds(),
		"request_size_bytes", 1024,
		"response_size_bytes", 512,
		"user_id", 12345,
	)

	// 调试信息
	zapLogger.DebugKV("缓存操作",
		"operation", "SET",
		"key", "user:12345:profile",
		"ttl_seconds", 3600,
		"size_bytes", 256,
		"cache_hit_rate", 0.85,
	)

	fmt.Println()
}

// demonstrateLogrusStyle 演示 Logrus 风格的日志记录
func demonstrateLogrusStyle(log logger.ILogger) {
	fmt.Println("2. === Logrus 风格演示 ===")
	fmt.Println("特点：字段化日志、链式调用、插件系统、钩子支持")

	// Logrus 风格的字段日志
	logrusLogger := log.WithField("framework", "logrus-style")

	// 基础字段日志
	logrusLogger.WithField("component", "auth").
		WithField("action", "login").
		WithField("user_id", 67890).
		WithField("ip_address", "192.168.1.100").
		Info("用户登录成功")

	// 多字段日志
	logrusLogger.WithFields(map[string]interface{}{
		"service":     "payment-processor",
		"transaction": "txn_abc123",
		"amount":      299.99,
		"currency":    "USD",
		"gateway":     "stripe",
		"merchant_id": "merchant_xyz",
	}).Info("支付处理开始")

	// 错误处理（Logrus 风格）
	err := errors.New("payment gateway timeout")
	logrusLogger.WithError(err).
		WithField("transaction_id", "txn_def456").
		WithField("amount", 149.99).
		WithField("retry_attempt", 2).
		Error("支付处理失败")

	// 警告日志
	logrusLogger.WithField("metric", "response_time").
		WithField("value_ms", 1500).
		WithField("threshold_ms", 1000).
		WithField("endpoint", "/api/orders").
		Warn("API响应时间超过阈值")

	// 链式字段构建
	contextLogger := logrusLogger.
		WithField("request_id", "req-12345").
		WithField("session_id", "sess-67890").
		WithField("trace_id", "trace-abcdef")

	contextLogger.Debug("请求上下文信息")
	contextLogger.Info("处理业务逻辑")
	contextLogger.Warn("资源使用率较高")

	fmt.Println()
}

// demonstrateSlogStyle 演示 slog 风格的日志记录
func demonstrateSlogStyle(log logger.ILogger) {
	fmt.Println("3. === slog 风格演示 ===")
	fmt.Println("特点：上下文感知、层次化属性、Go 1.21+ 标准库")

	slogLogger := log.WithField("framework", "slog-style")

	// 上下文感知日志
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-slog-001")
	ctx = context.WithValue(ctx, "user_id", "user-123")
	ctx = context.WithValue(ctx, "trace_id", "trace-slog-abc")

	// slog 风格的上下文日志
	slogLogger.InfoContext(ctx, "处理用户请求: %s", "GET /profile")
	slogLogger.DebugContext(ctx, "查询数据库: %s 表", "users")
	slogLogger.WarnContext(ctx, "数据库查询耗时: %dms", 800)

	// 带参数的上下文日志
	slogLogger.InfoContext(ctx, "用户操作记录: 用户 %d 执行了 %s 操作", 123, "profile_update")
	slogLogger.ErrorContext(ctx, "操作失败: %s", "profile update validation failed")

	// 层次化属性（使用字段模拟）
	slogLogger.WithField("module", "authentication").
		WithField("sub_module", "oauth2").
		InfoContext(ctx, "OAuth2 令牌验证: %s", "成功")

	// 结合上下文和键值对
	ctxLogger := slogLogger.WithContext(ctx)
	ctxLogger.InfoKV("上下文感知的结构化日志",
		"operation", "data_processing",
		"records_processed", 1500,
		"processing_time_ms", 2500,
		"memory_usage_mb", 128,
	)

	fmt.Println()
}

// demonstrateZerologStyle 演示 Zerolog 风格的日志记录
func demonstrateZerologStyle(log logger.ILogger) {
	fmt.Println("4. === Zerolog 风格演示 ===")
	fmt.Println("特点：零内存分配、JSON输出、链式API、高性能")

	zerologLogger := log.WithField("framework", "zerolog-style")

	// Zerolog 风格的链式调用
	zerologLogger.WithField("service", "notification").
		WithField("channel", "email").
		WithField("template", "welcome").
		WithField("recipient", "user@example.com").
		WithField("attempt", 1).
		Info("发送欢迎邮件")

	// 事件驱动的日志记录
	zerologLogger.WithField("event", "order_created").
		WithField("order_id", "ord_123456").
		WithField("customer_id", 78901).
		WithField("total_amount", 459.99).
		WithField("currency", "USD").
		WithField("payment_method", "credit_card").
		WithField("shipping_method", "express").
		Info("订单创建事件")

	// 性能指标记录
	zerologLogger.WithField("metric_type", "performance").
		WithField("operation", "image_processing").
		WithField("image_size_mb", 15.6).
		WithField("processing_time_ms", 3200).
		WithField("output_format", "webp").
		WithField("quality", 85).
		WithField("compression_ratio", 0.75).
		Debug("图片处理性能指标")

	// 错误事件记录
	zerologLogger.WithField("event", "error_occurred").
		WithField("error_code", "CONN_TIMEOUT").
		WithField("error_message", "connection timeout after 30s").
		WithField("service", "external_api").
		WithField("endpoint", "https://api.partner.com/v1/data").
		WithField("timeout_seconds", 30).
		WithField("retry_count", 3).
		WithField("circuit_breaker_state", "open").
		Error("外部服务连接超时")

	// 带时间戳的事件
	zerologLogger.WithField("event_time", time.Now().Format(time.RFC3339)).
		WithField("event_type", "system_health_check").
		WithField("cpu_usage_percent", 25.5).
		WithField("memory_usage_percent", 68.2).
		WithField("disk_usage_percent", 45.8).
		WithField("network_latency_ms", 12).
		WithField("health_status", "healthy").
		Info("系统健康检查完成")

	fmt.Println()
}

// demonstrateStdLogStyle 演示标准库 log 风格
func demonstrateStdLogStyle(log logger.ILogger) {
	fmt.Println("5. === 标准库 log 风格演示 ===")
	fmt.Println("特点：简单直接、格式化输出、兼容性好")

	stdLogger := log.WithField("framework", "stdlib-style")

	// 标准库风格的基础日志
	stdLogger.Print("服务启动完成")
	stdLogger.Printf("监听端口: %d", 8080)
	stdLogger.Println("准备接受请求")

	// 格式化日志
	userID := 12345
	action := "login"
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	stdLogger.Printf("用户活动: [%s] 用户 %d 执行了 %s 操作", timestamp, userID, action)
	stdLogger.Printf("数据库查询: 查询用户表，返回 %d 条记录，耗时 %dms", 1, 45)
	stdLogger.Printf("缓存操作: 设置键 user:%d:session，TTL=%d秒", userID, 3600)

	// 错误和状态报告
	stdLogger.Printf("错误统计: 过去24小时内发生 %d 个错误，成功率 %.2f%%", 23, 99.85)
	stdLogger.Printf("性能指标: 平均响应时间 %dms，QPS %d", 120, 850)

	// 系统事件日志
	stdLogger.Printf("系统事件: [%s] %s", "STARTUP", "所有服务启动完成")
	stdLogger.Printf("系统事件: [%s] %s", "CONFIG_RELOAD", "配置文件重新加载")
	stdLogger.Printf("系统事件: [%s] %s", "HEALTH_CHECK", "健康检查通过")

	fmt.Println()
}

// demonstrateMixedUsage 演示混合使用不同风格
func demonstrateMixedUsage(log logger.ILogger) {
	fmt.Println("6. === 混合使用演示 ===")
	fmt.Println("特点：在同一应用中根据场景选择最适合的日志风格")

	// 创建带基础字段的日志器
	appLogger := log.WithFields(map[string]interface{}{
		"app_name":    "e-commerce-api",
		"app_version": "v2.1.0",
		"environment": "production",
		"region":      "us-east-1",
	})

	// 场景1：HTTP请求处理（使用 slog 风格的上下文日志）
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-mixed-001")
	ctx = context.WithValue(ctx, "user_id", 54321)

	appLogger.InfoContext(ctx, "收到API请求: %s %s", "POST", "/api/checkout")

	// 场景2：业务逻辑处理（使用 Zap 风格的结构化日志）
	appLogger.InfoKV("开始结账流程",
		"cart_id", "cart-789012",
		"item_count", 3,
		"total_amount", 237.50,
		"coupon_code", "SAVE10",
		"discount_amount", 23.75,
	)

	// 场景3：外部服务调用（使用 Logrus 风格的字段日志）
	appLogger.WithField("external_service", "payment-gateway").
		WithField("provider", "stripe").
		WithField("amount", 213.75).
		WithField("currency", "USD").
		Info("调用外部支付服务")

	// 场景4：错误处理（使用 Zerolog 风格的事件日志）
	appLogger.WithField("event", "payment_failed").
		WithField("error_code", "CARD_DECLINED").
		WithField("transaction_id", "txn_failed_001").
		WithField("amount", 213.75).
		WithField("card_last_four", "1234").
		WithField("decline_reason", "insufficient_funds").
		Error("支付被拒绝")

	// 场景5：性能监控（使用标准库风格的简单日志）
	appLogger.Printf("性能报告: 结账流程总耗时 %dms，平均响应时间 %dms", 1850, 1200)

	// 场景6：系统状态（混合使用）
	appLogger.WithField("component", "health_monitor").
		InfoKV("系统健康状态更新",
			"cpu_usage", 45.2,
			"memory_usage", 72.8,
			"active_connections", 1247,
			"queue_length", 23,
			"error_rate", 0.02,
		)

	// 场景7：审计日志（使用完整的上下文和字段）
	auditLogger := appLogger.WithContext(ctx).
		WithField("audit_type", "security").
		WithField("action", "sensitive_data_access").
		WithField("resource", "customer_payment_info")

	auditLogger.InfoKV("敏感数据访问审计",
		"accessed_by", "user_54321",
		"access_method", "api_call",
		"data_type", "encrypted_card_number",
		"access_time", time.Now().Format(time.RFC3339),
		"ip_address", "192.168.1.200",
		"user_agent", "Mobile-App/2.1.0",
		"permission_level", "customer_data_read",
	)

	fmt.Printf("\n混合使用演示了如何在同一应用中灵活选择最适合的日志风格：\n")
	fmt.Printf("- HTTP请求: slog风格的上下文日志\n")
	fmt.Printf("- 业务逻辑: Zap风格的键值对日志\n")
	fmt.Printf("- 外部调用: Logrus风格的字段日志\n")
	fmt.Printf("- 错误处理: Zerolog风格的事件日志\n")
	fmt.Printf("- 性能监控: 标准库风格的简单日志\n")
	fmt.Printf("- 审计记录: 组合多种风格的完整日志\n")

	fmt.Println()
}

// 演示如何封装不同框架风格的快捷方法
type FrameworkLoggers struct {
	base logger.ILogger
}

func NewFrameworkLoggers(base logger.ILogger) *FrameworkLoggers {
	return &FrameworkLoggers{base: base}
}

// Zap 风格的快捷方法
func (fl *FrameworkLoggers) ZapInfo(msg string, fields ...interface{}) {
	fl.base.InfoKV(msg, fields...)
}

func (fl *FrameworkLoggers) ZapError(msg string, err error, fields ...interface{}) {
	allFields := append([]interface{}{"error", err.Error()}, fields...)
	fl.base.ErrorKV(msg, allFields...)
}

// Logrus 风格的快捷方法
func (fl *FrameworkLoggers) LogrusWithField(key string, value interface{}) logger.ILogger {
	return fl.base.WithField(key, value)
}

func (fl *FrameworkLoggers) LogrusWithError(err error) logger.ILogger {
	return fl.base.WithError(err)
}

// slog 风格的快捷方法
func (fl *FrameworkLoggers) SlogInfoContext(ctx context.Context, msg string, args ...interface{}) {
	fl.base.InfoContext(ctx, msg, args...)
}

// Zerolog 风格的快捷方法
func (fl *FrameworkLoggers) ZerologEvent() logger.ILogger {
	return fl.base
}

// 演示封装使用
func demonstrateFrameworkWrappers() {
	fmt.Println("7. === 框架封装演示 ===")

	config := logger.DefaultConfig().WithPrefix("[Wrapper] ")
	baseLogger := logger.NewLogger(config)
	frameworks := NewFrameworkLoggers(baseLogger)

	// 使用封装的方法
	frameworks.ZapInfo("Zap风格的信息日志",
		"service", "api",
		"action", "startup",
	)

	err := errors.New("示例错误")
	frameworks.ZapError("Zap风格的错误日志", err,
		"component", "database",
		"retry_count", 3,
	)

	frameworks.LogrusWithField("component", "cache").
		WithField("operation", "flush").
		Info("Logrus风格的字段日志")

	ctx := context.Background()
	frameworks.SlogInfoContext(ctx, "slog风格的上下文日志: %s", "示例消息")

	frameworks.ZerologEvent().
		WithField("event", "system_event").
		WithField("type", "maintenance").
		Info("Zerolog风格的事件日志")

	fmt.Println()
}