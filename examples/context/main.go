/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 13:10:00
 * @FilePath: \go-logger\examples\context\main.go
 * @Description: 上下文示例 - 演示上下文相关的日志功能
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package main

import (
	"context"
	"fmt"
	"github.com/kamalyes/go-logger"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("🎯 Go Logger - 上下文示例演示")
	fmt.Println(strings.Repeat("=", 50))

	// 1. 基础上下文日志
	demonstrateBasicContext()

	fmt.Println()

	// 2. 上下文传递演示
	demonstrateContextPropagation()

	fmt.Println()

	// 3. 上下文取消演示
	demonstrateContextCancellation()

	fmt.Println()

	// 4. 上下文超时演示
	demonstrateContextTimeout()

	fmt.Println()

	// 5. 上下文值传递
	demonstrateContextValues()

	fmt.Println()

	// 6. 实际应用场景
	demonstrateRealWorldScenarios()
}

// 基础上下文日志
func demonstrateBasicContext() {
	fmt.Println("📝 1. 基础上下文日志")
	fmt.Println(strings.Repeat("-", 30))

	// 创建适配器
	adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.DEBUG,
		Output:     os.Stdout,
		TimeFormat: "15:04:05",
		Colorful:   true,
	})
	adapter.Initialize()

	// 创建上下文
	ctx := context.Background()

	fmt.Println("\n🔹 基础上下文方法:")
	adapter.DebugContext(ctx, "这是带上下文的调试信息")
	adapter.InfoContext(ctx, "这是带上下文的普通信息")
	adapter.WarnContext(ctx, "这是带上下文的警告信息")
	adapter.ErrorContext(ctx, "这是带上下文的错误信息")

	fmt.Println("\n🔹 带上下文的格式化日志:")
	adapter.DebugContext(ctx, "用户 %d 执行了 %s 操作", 12345, "登录")
	adapter.InfoContext(ctx, "处理请求耗时 %v", 150*time.Millisecond)

	// 与非上下文方法对比
	fmt.Println("\n🔹 对比非上下文方法:")
	adapter.Info("普通日志方法")
	adapter.InfoContext(ctx, "上下文日志方法")

	defer adapter.Close()
}

// 上下文传递演示
func demonstrateContextPropagation() {
	fmt.Println("🔄 2. 上下文传递演示")
	fmt.Println(strings.Repeat("-", 30))

	adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     os.Stdout,
		TimeFormat: "15:04:05",
	})
	adapter.Initialize()

	// 创建带值的上下文
	ctx := context.WithValue(context.Background(), "requestID", "req-12345")
	ctx = context.WithValue(ctx, "userID", "user-67890")

	fmt.Println("\n🔹 模拟请求处理链:")

	// 模拟HTTP处理器
	handleRequest(ctx, adapter)

	defer adapter.Close()
}

func handleRequest(ctx context.Context, logger logger.IAdapter) {
	logger.InfoContext(ctx, "开始处理HTTP请求")

	// 调用业务逻辑
	processBusinessLogic(ctx, logger)

	logger.InfoContext(ctx, "HTTP请求处理完成")
}

func processBusinessLogic(ctx context.Context, logger logger.IAdapter) {
	logger.InfoContext(ctx, "执行业务逻辑")

	// 调用数据库操作
	queryDatabase(ctx, logger)

	// 调用外部API
	callExternalAPI(ctx, logger)

	logger.InfoContext(ctx, "业务逻辑执行完成")
}

func queryDatabase(ctx context.Context, logger logger.IAdapter) {
	logger.InfoContext(ctx, "查询数据库")
	time.Sleep(50 * time.Millisecond) // 模拟数据库查询
	logger.InfoContext(ctx, "数据库查询完成")
}

func callExternalAPI(ctx context.Context, logger logger.IAdapter) {
	logger.InfoContext(ctx, "调用外部API")
	time.Sleep(100 * time.Millisecond) // 模拟API调用
	logger.InfoContext(ctx, "外部API调用完成")
}

// 上下文取消演示
func demonstrateContextCancellation() {
	fmt.Println("❌ 3. 上下文取消演示")
	fmt.Println(strings.Repeat("-", 30))

	adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     os.Stdout,
		TimeFormat: "15:04:05",
	})
	adapter.Initialize()

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())

	fmt.Println("\n🔹 启动可取消的任务:")

	// 启动长时间运行的任务
	go longRunningTask(ctx, adapter)

	// 等待一段时间后取消
	time.Sleep(2 * time.Second)
	fmt.Println("\n🔹 取消任务:")
	cancel()

	// 等待任务完成
	time.Sleep(500 * time.Millisecond)

	defer adapter.Close()
}

func longRunningTask(ctx context.Context, logger logger.IAdapter) {
	logger.InfoContext(ctx, "开始长时间运行的任务")

	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			logger.WarnContext(ctx, "任务被取消: %v", ctx.Err())
			return
		default:
			logger.InfoContext(ctx, "任务进度: %d/10", i+1)
			time.Sleep(500 * time.Millisecond)
		}
	}

	logger.InfoContext(ctx, "长时间任务完成")
}

// 上下文超时演示
func demonstrateContextTimeout() {
	fmt.Println("⏰ 4. 上下文超时演示")
	fmt.Println(strings.Repeat("-", 30))

	adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     os.Stdout,
		TimeFormat: "15:04:05",
	})
	adapter.Initialize()

	fmt.Println("\n🔹 设置2秒超时的任务:")

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 启动可能超时的任务
	timeoutTask(ctx, adapter)

	defer adapter.Close()
}

func timeoutTask(ctx context.Context, logger logger.IAdapter) {
	logger.InfoContext(ctx, "开始可能超时的任务")

	// 模拟耗时操作
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				logger.ErrorContext(ctx, "任务超时")
			} else {
				logger.WarnContext(ctx, "任务被取消")
			}
			return
		default:
			logger.InfoContext(ctx, "执行步骤 %d", i+1)
			time.Sleep(800 * time.Millisecond) // 每步800ms，总共需要4秒
		}
	}

	logger.InfoContext(ctx, "任务成功完成")
}

// 上下文值传递
func demonstrateContextValues() {
	fmt.Println("💼 5. 上下文值传递")
	fmt.Println(strings.Repeat("-", 30))

	adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     os.Stdout,
		TimeFormat: "15:04:05",
	})
	adapter.Initialize()

	fmt.Println("\n🔹 在上下文中传递跟踪信息:")

	// 创建带跟踪信息的上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceID", "trace-abc123")
	ctx = context.WithValue(ctx, "spanID", "span-def456")
	ctx = context.WithValue(ctx, "userID", "user-12345")
	ctx = context.WithValue(ctx, "sessionID", "sess-789")

	// 使用WithContext创建带上下文的logger
	contextLogger := adapter.WithContext(ctx)

	fmt.Println("  模拟业务流程:")
	contextLogger.Info("用户认证")
	contextLogger.Info("权限检查")
	contextLogger.Info("数据查询")
	contextLogger.Info("结果返回")

	// 演示从上下文提取值
	fmt.Println("\n🔹 从上下文提取信息:")
	if traceID := ctx.Value("traceID"); traceID != nil {
		adapter.InfoContext(ctx, "当前追踪ID: %s", traceID)
	}
	if userID := ctx.Value("userID"); userID != nil {
		adapter.InfoContext(ctx, "当前用户ID: %s", userID)
	}

	defer adapter.Close()
}

// 实际应用场景
func demonstrateRealWorldScenarios() {
	fmt.Println("🌍 6. 实际应用场景")
	fmt.Println(strings.Repeat("-", 30))

	adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     os.Stdout,
		TimeFormat: "15:04:05",
	})
	adapter.Initialize()

	fmt.Println("\n🔹 Web服务器请求处理:")
	simulateWebServer(adapter)

	fmt.Println("\n🔹 并发任务处理:")
	simulateConcurrentTasks(adapter)

	defer adapter.Close()
}

func simulateWebServer(logger logger.IAdapter) {
	// 模拟3个并发的HTTP请求
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			// 为每个请求创建独立的上下文
			ctx := context.Background()
			ctx = context.WithValue(ctx, "requestID", fmt.Sprintf("req-%d", requestID))
			ctx = context.WithValue(ctx, "startTime", time.Now())

			// 模拟请求处理
			logger.InfoContext(ctx, "收到HTTP请求")

			// 随机处理时间
			processingTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
			time.Sleep(processingTime)

			if rand.Float32() > 0.7 { // 30% 概率出错
				logger.ErrorContext(ctx, "请求处理失败")
			} else {
				logger.InfoContext(ctx, "请求处理成功，耗时: %v", processingTime)
			}
		}(i + 1)
	}

	wg.Wait()
}

func simulateConcurrentTasks(logger logger.IAdapter) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 启动多个并发任务
	var wg sync.WaitGroup
	taskCount := 5

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			// 为每个任务创建子上下文
			taskCtx := context.WithValue(ctx, "taskID", taskID)

			logger.InfoContext(taskCtx, "任务开始")

			// 模拟任务执行
			for step := 0; step < 3; step++ {
				select {
				case <-taskCtx.Done():
					logger.WarnContext(taskCtx, "任务被中断: %v", taskCtx.Err())
					return
				default:
					logger.InfoContext(taskCtx, "执行步骤 %d", step+1)
					time.Sleep(time.Duration(rand.Intn(800)+200) * time.Millisecond)
				}
			}

			logger.InfoContext(taskCtx, "任务完成")
		}(i + 1)
	}

	wg.Wait()
}
