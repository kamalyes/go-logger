/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 13:02:26
 * @FilePath: \go-logger\examples\formatters\main.go
 * @Description: 格式化器示例 - 演示不同的日志输出格式
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kamalyes/go-logger"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("🎨 Go Logger - 格式化器示例演示")
	fmt.Println(strings.Repeat("=", 50))

	// 1. 基础格式演示
	demonstrateBasicFormats()

	fmt.Println()

	// 2. 时间格式演示
	demonstrateTimeFormats()

	fmt.Println()

	// 3. 级别显示演示
	demonstrateLevelFormats()

	fmt.Println()

	// 4. 字段格式演示
	demonstrateFieldFormats()

	fmt.Println()

	// 5. 颜色格式演示
	demonstrateColorFormats()

	fmt.Println()

	// 6. 结构化日志演示
	demonstrateStructuredFormats()
}

// 基础格式演示
func demonstrateBasicFormats() {
	fmt.Println("📄 1. 基础格式演示")
	fmt.Println(strings.Repeat("-", 30))

	// 1.1 标准文本格式
	fmt.Println("\n🔹 标准文本格式:")
	textLogger, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     os.Stdout,
		Format:     "text",
		TimeFormat: "15:04:05",
		Colorful:   false,
	})
	textLogger.Initialize()

	textLogger.Info("这是标准文本格式的日志")
	textLogger.WithField("component", "formatter").Info("带字段的文本日志")

	// 1.2 JSON格式
	fmt.Println("\n🔹 JSON格式:")
	jsonBuffer := &bytes.Buffer{}
	jsonLogger, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     jsonBuffer,
		Format:     "json",
		TimeFormat: time.RFC3339,
		Colorful:   false,
	})
	jsonLogger.Initialize()

	jsonLogger.Info("这是JSON格式的日志")
	jsonLogger.WithField("component", "formatter").WithField("version", "1.0.0").Info("带字段的JSON日志")

	fmt.Print(jsonBuffer.String())

	defer textLogger.Close()
	defer jsonLogger.Close()
}

// 时间格式演示
func demonstrateTimeFormats() {
	fmt.Println("⏰ 2. 时间格式演示")
	fmt.Println(strings.Repeat("-", 30))

	timeFormats := map[string]string{
		"标准时间":    "15:04:05",
		"完整日期":    "2006-01-02 15:04:05",
		"RFC3339": time.RFC3339,
		"RFC822":  time.RFC822,
		"自定义格式":   "2006年01月02日 15:04:05",
	}

	for name, format := range timeFormats {
		fmt.Printf("\n🔹 %s (%s):\n", name, format)

		adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
			Type:       logger.StandardAdapter,
			Level:      logger.INFO,
			Output:     os.Stdout,
			Format:     "text",
			TimeFormat: format,
			Colorful:   false,
		})
		adapter.Initialize()

		adapter.Info("时间格式示例")
		adapter.Close()
	}
}

// 级别显示演示
func demonstrateLevelFormats() {
	fmt.Println("📊 3. 级别显示演示")
	fmt.Println(strings.Repeat("-", 30))

	adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.DEBUG,
		Output:     os.Stdout,
		Format:     "text",
		TimeFormat: "15:04:05",
		Colorful:   false,
	})
	adapter.Initialize()

	fmt.Println("\n🔹 所有级别展示:")
	adapter.Debug("这是调试信息 - 用于开发阶段")
	adapter.Info("这是普通信息 - 一般运行信息")
	adapter.Warn("这是警告信息 - 需要注意")
	adapter.Error("这是错误信息 - 发生了错误")

	adapter.Close()
}

// 字段格式演示
func demonstrateFieldFormats() {
	fmt.Println("🏷️ 4. 字段格式演示")
	fmt.Println(strings.Repeat("-", 30))

	adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     os.Stdout,
		Format:     "text",
		TimeFormat: "15:04:05",
		Colorful:   false,
	})
	adapter.Initialize()

	fmt.Println("\n🔹 不同字段类型:")

	// 单个字段
	adapter.WithField("user_id", 12345).Info("单个字段示例")

	// 多个字段
	adapter.WithFields(map[string]interface{}{
		"user_id":   12345,
		"action":    "login",
		"ip":        "192.168.1.100",
		"timestamp": time.Now().Unix(),
		"success":   true,
	}).Info("多字段示例")

	// 错误字段
	err := fmt.Errorf("数据库连接失败")
	adapter.WithError(err).Error("错误日志示例")

	adapter.Close()
}

// 颜色格式演示
func demonstrateColorFormats() {
	fmt.Println("🌈 5. 颜色格式演示")
	fmt.Println(strings.Repeat("-", 30))

	fmt.Println("\n🔹 带颜色的日志:")
	colorAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.DEBUG,
		Output:     os.Stdout,
		Format:     "text",
		TimeFormat: "15:04:05",
		Colorful:   true,
	})
	colorAdapter.Initialize()

	colorAdapter.Debug("彩色调试信息")
	colorAdapter.Info("彩色普通信息")
	colorAdapter.Warn("彩色警告信息")
	colorAdapter.Error("彩色错误信息")

	fmt.Println("\n🔹 不带颜色的日志:")
	noColorAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.DEBUG,
		Output:     os.Stdout,
		Format:     "text",
		TimeFormat: "15:04:05",
		Colorful:   false,
	})
	noColorAdapter.Initialize()

	noColorAdapter.Debug("无色调试信息")
	noColorAdapter.Info("无色普通信息")
	noColorAdapter.Warn("无色警告信息")
	noColorAdapter.Error("无色错误信息")

	colorAdapter.Close()
	noColorAdapter.Close()
}

// 结构化日志演示
func demonstrateStructuredFormats() {
	fmt.Println("🏗️ 6. 结构化日志演示")
	fmt.Println(strings.Repeat("-", 30))

	// 创建JSON格式的适配器
	jsonBuffer := &bytes.Buffer{}
	jsonAdapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     jsonBuffer,
		Format:     "json",
		TimeFormat: time.RFC3339Nano,
		Colorful:   false,
		Fields: map[string]interface{}{
			"service": "demo-app",
			"version": "1.0.0",
		},
	})
	jsonAdapter.Initialize()

	fmt.Println("\n🔹 结构化JSON日志:")

	// API请求日志
	jsonAdapter.WithFields(map[string]interface{}{
		"method":    "GET",
		"endpoint":  "/api/users/123",
		"status":    200,
		"duration":  "45ms",
		"client_ip": "192.168.1.100",
	}).Info("API请求处理完成")

	// 数据库查询日志
	jsonAdapter.WithFields(map[string]interface{}{
		"operation": "SELECT",
		"table":     "users",
		"duration":  "12ms",
		"rows":      1,
		"query_id":  "q_12345",
	}).Info("数据库查询执行")

	// 错误日志
	jsonAdapter.WithFields(map[string]interface{}{
		"error_code":  "DB_CONNECTION_FAILED",
		"retry_count": 3,
		"last_error":  "connection timeout",
		"component":   "database",
	}).Error("数据库连接失败")

	// 输出JSON结果
	jsonOutput := jsonBuffer.String()
	lines := strings.Split(strings.TrimSpace(jsonOutput), "\n")

	for i, line := range lines {
		if line != "" {
			// 美化JSON输出
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(line), &jsonData); err == nil {
				prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
				fmt.Printf("  日志 %d:\n%s\n\n", i+1, string(prettyJSON))
			}
		}
	}

	// 性能统计
	fmt.Println("🔹 性能统计示例:")
	perfFields := map[string]interface{}{
		"request_count":      1000,
		"avg_response_time":  "25ms",
		"error_rate":         "0.5%",
		"memory_usage":       "256MB",
		"cpu_usage":          "15%",
		"active_connections": 45,
	}

	jsonAdapter.WithFields(perfFields).Info("系统性能统计")

	// 输出最后一条日志
	lastLog := jsonBuffer.String()
	lastLines := strings.Split(strings.TrimSpace(lastLog), "\n")
	if len(lastLines) > 0 {
		lastLine := lastLines[len(lastLines)-1]
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(lastLine), &jsonData); err == nil {
			prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
			fmt.Printf("  性能统计:\n%s\n", string(prettyJSON))
		}
	}

	jsonAdapter.Close()
}
