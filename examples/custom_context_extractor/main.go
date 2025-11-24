/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-24 11:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-24 11:30:00
 * @FilePath: \go-logger\examples\custom_context_extractor\main.go
 * @Description: 自定义上下文提取器使用示例
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"context"
	logger "github.com/kamalyes/go-logger"
	"os"
)

func main() {
	// 示例 1: 使用默认提取器
	example1DefaultExtractor()

	// 示例 2: 禁用上下文提取
	example2NoExtractor()

	// 示例 3: 只提取 TraceID
	example3TraceIDOnly()

	// 示例 4: 自定义字段提取
	example4CustomFields()

	// 示例 5: 链式提取器
	example5ChainedExtractors()

	// 示例 6: 使用构建器
	example6Builder()

	// 示例 7: 条件提取器
	example7Conditional()
}

// 示例 1: 使用默认提取器（TraceID + RequestID）
func example1DefaultExtractor() {
	println("=== 示例 1: 默认提取器 ===")

	config := logger.DefaultConfig()
	config.Output = os.Stdout
	log := logger.NewUltraFastLogger(config)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-12345")
	ctx = context.WithValue(ctx, "request_id", "req-67890")

	log.InfoContext(ctx, "用户登录成功")
	// 输出: [TraceID=trace-12345 RequestID=req-67890] 用户登录成功
}

// 示例 2: 禁用上下文提取
func example2NoExtractor() {
	println("\n=== 示例 2: 禁用上下文提取 ===")

	config := logger.DefaultConfig()
	config.Output = os.Stdout
	log := logger.NewUltraFastLogger(config)

	// 设置为空操作提取器
	log.SetContextExtractor(logger.NoOpContextExtractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-12345")

	log.InfoContext(ctx, "这条日志不会包含上下文信息")
	// 输出: 这条日志不会包含上下文信息
}

// 示例 3: 只提取 TraceID
func example3TraceIDOnly() {
	println("\n=== 示例 3: 只提取 TraceID ===")

	config := logger.DefaultConfig()
	config.Output = os.Stdout
	log := logger.NewUltraFastLogger(config)

	// 只提取 TraceID
	log.SetContextExtractor(logger.SimpleTraceIDExtractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-12345")
	ctx = context.WithValue(ctx, "request_id", "req-67890")

	log.InfoContext(ctx, "只显示 TraceID")
	// 输出: [TraceID=trace-12345] 只显示 TraceID
}

// 示例 4: 自定义字段提取
func example4CustomFields() {
	println("\n=== 示例 4: 自定义字段提取 ===")

	config := logger.DefaultConfig()
	config.Output = os.Stdout
	log := logger.NewUltraFastLogger(config)

	// 提取自定义字段：user_id 和 session_id
	extractor := logger.CustomFieldExtractor(
		[]string{"user_id", "session_id"}, // context keys
		[]string{},                        // metadata keys
	)
	log.SetContextExtractor(extractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", "user-123")
	ctx = context.WithValue(ctx, "session_id", "sess-456")

	log.InfoContext(ctx, "用户操作")
	// 输出: [user_id=user-123 session_id=sess-456] 用户操作
}

// 示例 5: 链式提取器
func example5ChainedExtractors() {
	println("\n=== 示例 5: 链式提取器 ===")

	config := logger.DefaultConfig()
	config.Output = os.Stdout
	log := logger.NewUltraFastLogger(config)

	// 组合多个提取器
	extractor := logger.ChainContextExtractors(
		logger.SimpleTraceIDExtractor,
		logger.ExtractFromContextValue("user_id", "User"),
		logger.ExtractFromContextValue("ip", "IP"),
	)
	log.SetContextExtractor(extractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-12345")
	ctx = context.WithValue(ctx, "user_id", "alice")
	ctx = context.WithValue(ctx, "ip", "192.168.1.1")

	log.InfoContext(ctx, "API 请求")
	// 输出: [TraceID=trace-12345] [User=alice] [IP=192.168.1.1] API 请求
}

// 示例 6: 使用构建器创建提取器
func example6Builder() {
	println("\n=== 示例 6: 使用构建器 ===")

	config := logger.DefaultConfig()
	config.Output = os.Stdout
	log := logger.NewUltraFastLogger(config)

	// 使用构建器创建提取器
	extractor := logger.NewContextExtractorBuilder().
		AddTraceID().
		AddRequestID().
		AddContextValue("tenant_id", "Tenant").
		AddContextValue("env", "Env").
		Build()

	log.SetContextExtractor(extractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-99999")
	ctx = context.WithValue(ctx, "request_id", "req-88888")
	ctx = context.WithValue(ctx, "tenant_id", "tenant-A")
	ctx = context.WithValue(ctx, "env", "production")

	log.InfoContext(ctx, "多租户请求")
	// 输出: [TraceID=trace-99999] [RequestID=req-88888] [Tenant=tenant-A] [Env=production] 多租户请求
}

// 示例 7: 条件提取器
func example7Conditional() {
	println("\n=== 示例 7: 条件提取器 ===")

	config := logger.DefaultConfig()
	config.Output = os.Stdout
	log := logger.NewUltraFastLogger(config)

	// 只在生产环境才提取详细信息
	extractor := logger.ConditionalContextExtractor(
		func(ctx context.Context) bool {
			env, ok := ctx.Value("env").(string)
			return ok && env == "production"
		},
		logger.ChainContextExtractors(
			logger.SimpleTraceIDExtractor,
			logger.SimpleRequestIDExtractor,
		),
	)
	log.SetContextExtractor(extractor)

	// 生产环境
	prodCtx := context.Background()
	prodCtx = context.WithValue(prodCtx, "env", "production")
	prodCtx = context.WithValue(prodCtx, "trace_id", "trace-prod")
	prodCtx = context.WithValue(prodCtx, "request_id", "req-prod")
	log.InfoContext(prodCtx, "生产环境日志")
	// 输出: [TraceID=trace-prod] [RequestID=req-prod] 生产环境日志

	// 开发环境
	devCtx := context.Background()
	devCtx = context.WithValue(devCtx, "env", "development")
	devCtx = context.WithValue(devCtx, "trace_id", "trace-dev")
	log.InfoContext(devCtx, "开发环境日志")
	// 输出: 开发环境日志
}

// 示例 8: 自定义提取器函数
func example8CustomFunction() {
	println("\n=== 示例 8: 自定义提取器函数 ===")

	config := logger.DefaultConfig()
	config.Output = os.Stdout
	log := logger.NewUltraFastLogger(config)

	// 完全自定义的提取器
	customExtractor := func(ctx context.Context) string {
		if ctx == nil {
			return ""
		}

		// 提取用户信息并格式化
		userId, _ := ctx.Value("user_id").(string)
		userName, _ := ctx.Value("user_name").(string)

		if userId != "" || userName != "" {
			return "[👤 User: " + userId + " (" + userName + ")] "
		}

		return ""
	}

	log.SetContextExtractor(customExtractor)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", "12345")
	ctx = context.WithValue(ctx, "user_name", "张三")

	log.InfoContext(ctx, "用户订单")
	// 输出: [👤 User: 12345 (张三)] 用户订单
}
