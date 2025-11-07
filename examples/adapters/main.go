/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\examples\adapters\main.go
 * @Description: 适配器模式使用示例
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package main

import (
	"fmt"
	"os"
	"time"

	logger "github.com/kamalyes/go-logger"
)

func main() {
	fmt.Println("=== Go Logger 适配器模式示例 ===")
	fmt.Println()

	// 创建日志管理器
	manager := logger.NewLoggerManager()
	defer manager.CloseAll()

	// 1. 创建控制台适配器
	consoleConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "console",
		Level:      logger.DEBUG,
		Output:     os.Stdout,
		Format:     "text",
		TimeFormat: "2006-01-02 15:04:05",
		Colorful:   true,
		Fields: map[string]interface{}{
			"app":     "adapter-demo",
			"version": "1.0.0",
		},
	}

	consoleAdapter, err := logger.CreateAdapter(consoleConfig)
	if err != nil {
		fmt.Printf("创建控制台适配器失败: %v\n", err)
		return
	}

	err = manager.AddAdapter("console", consoleAdapter)
	if err != nil {
		fmt.Printf("添加控制台适配器失败: %v\n", err)
		return
	}

	// 2. 创建文件适配器
	fileConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "file",
		Level:      logger.INFO,
		File:       "app.log",
		MaxSize:    10, // 10MB
		MaxBackups: 3,
		MaxAge:     7, // 7天
		Compress:   true,
		Format:     "json",
		TimeFormat: time.RFC3339,
		Fields: map[string]interface{}{
			"component": "file-logger",
			"env":       "development",
		},
	}

	// 创建文件输出
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("创建日志文件失败: %v\n", err)
		return
	}
	defer logFile.Close()

	fileConfig.Output = logFile
	fileAdapter, err := logger.CreateAdapter(fileConfig)
	if err != nil {
		fmt.Printf("创建文件适配器失败: %v\n", err)
		return
	}

	err = manager.AddAdapter("file", fileAdapter)
	if err != nil {
		fmt.Printf("添加文件适配器失败: %v\n", err)
		return
	}

	// 3. 创建错误专用适配器
	errorConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Name:       "error",
		Level:      logger.ERROR,
		Output:     os.Stderr,
		Format:     "text",
		TimeFormat: "15:04:05",
		Colorful:   true,
		Fields: map[string]interface{}{
			"severity": "high",
			"alert":    true,
		},
	}

	errorAdapter, err := logger.CreateAdapter(errorConfig)
	if err != nil {
		fmt.Printf("创建错误适配器失败: %v\n", err)
		return
	}

	err = manager.AddAdapter("error", errorAdapter)
	if err != nil {
		fmt.Printf("添加错误适配器失败: %v\n", err)
		return
	}

	// 4. 展示所有适配器
	fmt.Println("已注册的适配器:")
	for _, name := range manager.ListAdapters() {
		adapter, exists := manager.GetAdapter(name)
		if exists {
			fmt.Printf("- %s: %s (v%s) - %s\n", 
				name, 
				adapter.GetAdapterName(), 
				adapter.GetAdapterVersion(),
				func() string {
					if adapter.IsHealthy() {
						return "健康"
					}
					return "异常"
				}())
		}
	}
	fmt.Println()

	// 5. 通过管理器广播日志
	fmt.Println("=== 通过管理器广播日志 ===")
	manager.Broadcast(logger.DEBUG, "这是调试信息，只有控制台适配器会输出")
	manager.Broadcast(logger.INFO, "这是信息日志，控制台和文件适配器都会输出")
	manager.Broadcast(logger.WARN, "这是警告信息，所有适配器都会输出")
	manager.Broadcast(logger.ERROR, "这是错误信息，所有适配器都会输出")
	fmt.Println()

	// 6. 获取特定适配器并单独使用
	fmt.Println("=== 使用特定适配器 ===")
	
	if consoleAdapter, exists := manager.GetAdapter("console"); exists {
		console := consoleAdapter.WithField("module", "user-service")
		console.Info("用户登录成功")
		console.WithField("user_id", 12345).Info("用户操作记录")
		console.WithFields(map[string]interface{}{
			"action": "purchase",
			"amount": 99.99,
			"currency": "USD",
		}).Info("用户购买商品")
	}

	if fileAdapter, exists := manager.GetAdapter("file"); exists {
		file := fileAdapter.WithField("module", "order-service")
		file.Info("订单创建成功")
		file.WithError(fmt.Errorf("库存不足")).Error("订单处理失败")
	}

	if errorAdapter, exists := manager.GetAdapter("error"); exists {
		error := errorAdapter.WithField("module", "payment-service")
		error.Error("支付网关连接失败")
		error.Error("严重错误：数据库连接丢失（演示用，已改为ERROR级别避免程序退出）")
	}

	time.Sleep(100 * time.Millisecond) // 确保日志写入完成
	fmt.Println()

	// 7. 动态管理适配器
	fmt.Println("=== 动态管理适配器 ===")
	
	// 检查健康状态
	fmt.Println("适配器健康状态:")
	health := manager.HealthCheck()
	for name, healthy := range health {
		status := "正常"
		if !healthy {
			status = "异常"
		}
		fmt.Printf("- %s: %s\n", name, status)
	}

	// 调整日志级别
	fmt.Println("\n调整所有适配器日志级别为 WARN...")
	manager.SetLevelAll(logger.WARN)
	
	fmt.Println("测试日志级别调整效果:")
	manager.Broadcast(logger.DEBUG, "DEBUG级别 - 不应该输出")
	manager.Broadcast(logger.INFO, "INFO级别 - 不应该输出") 
	manager.Broadcast(logger.WARN, "WARN级别 - 应该输出")
	manager.Broadcast(logger.ERROR, "ERROR级别 - 应该输出")

	time.Sleep(100 * time.Millisecond)
	fmt.Println()

	// 8. 移除适配器
	fmt.Println("=== 移除适配器示例 ===")
	fmt.Printf("移除前适配器数量: %d\n", len(manager.ListAdapters()))
	
	err = manager.RemoveAdapter("error")
	if err != nil {
		fmt.Printf("移除错误适配器失败: %v\n", err)
	} else {
		fmt.Println("成功移除错误适配器")
	}
	
	fmt.Printf("移除后适配器数量: %d\n", len(manager.ListAdapters()))
	fmt.Println("剩余适配器:", manager.ListAdapters())

	// 9. 刷新所有适配器
	fmt.Println("\n刷新所有适配器缓冲区...")
	err = manager.FlushAll()
	if err != nil {
		fmt.Printf("刷新失败: %v\n", err)
	} else {
		fmt.Println("刷新成功")
	}

	fmt.Println("\n=== 适配器模式示例完成 ===")
}