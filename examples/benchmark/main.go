/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\examples\benchmark\main.go
 * @Description: 性能基准测试示例
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kamalyes/go-logger"
)

func main() {
	// 创建一个高性能的日志器配置
	config := logger.DefaultConfig().
		WithLevel(logger.INFO).
		WithShowCaller(false). // 关闭调用者信息以提高性能
		WithColorful(false).   // 关闭彩色输出以提高性能
		WithOutput(os.Stdout)

	log := logger.NewLogger(config)

	// 1. 基本性能测试
	fmt.Println("=== 基本性能测试 ===")
	testBasicPerformance(log)

	// 2. 结构化日志性能测试
	fmt.Println("\n=== 结构化日志性能测试 ===")
	testStructuredLogging(log)

	// 3. 不同级别的性能对比
	fmt.Println("\n=== 不同级别性能对比 ===")
	testLevelPerformance(log)

	// 4. 大量并发日志性能
	fmt.Println("\n=== 并发日志性能测试 ===")
	testConcurrentLogging(log)
}

func testBasicPerformance(log *logger.Logger) {
	iterations := 10000
	
	start := time.Now()
	for i := 0; i < iterations; i++ {
		log.Info("This is a basic performance test message #%d", i)
	}
	duration := time.Since(start)
	
	fmt.Printf("基本日志: %d 条消息耗时 %v\n", iterations, duration)
	fmt.Printf("平均每条消息: %v\n", duration/time.Duration(iterations))
	fmt.Printf("每秒处理消息: %.0f 条/秒\n", float64(iterations)/duration.Seconds())
}

func testStructuredLogging(log *logger.Logger) {
	iterations := 10000
	
	start := time.Now()
	for i := 0; i < iterations; i++ {
		log.WithFields(map[string]interface{}{
			"user_id":    12345,
			"session_id": "sess_abc123",
			"action":     "click",
			"page":       "/dashboard",
			"timestamp":  time.Now().Unix(),
			"iteration":  i,
		}).Info("User action performed")
	}
	duration := time.Since(start)
	
	fmt.Printf("结构化日志: %d 条消息耗时 %v\n", iterations, duration)
	fmt.Printf("平均每条消息: %v\n", duration/time.Duration(iterations))
	fmt.Printf("每秒处理消息: %.0f 条/秒\n", float64(iterations)/duration.Seconds())
}

func testLevelPerformance(log *logger.Logger) {
	iterations := 5000
	
	// 测试不同级别的性能
	levels := []struct {
		name string
		fn   func(string, ...interface{})
	}{
		{"DEBUG", log.Debug},
		{"INFO", log.Info},
		{"WARN", log.Warn},
		{"ERROR", log.Error},
	}
	
	for _, level := range levels {
		start := time.Now()
		for i := 0; i < iterations; i++ {
			level.fn("Level %s test message #%d", level.name, i)
		}
		duration := time.Since(start)
		
		fmt.Printf("%s 级别: %d 条消息耗时 %v (%.0f 条/秒)\n", 
			level.name, iterations, duration, float64(iterations)/duration.Seconds())
	}
}

func testConcurrentLogging(log *logger.Logger) {
	numGoroutines := 10
	messagesPerGoroutine := 1000
	totalMessages := numGoroutines * messagesPerGoroutine
	
	start := time.Now()
	done := make(chan bool, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineId int) {
			for j := 0; j < messagesPerGoroutine; j++ {
				log.WithFields(map[string]interface{}{
					"goroutine_id": goroutineId,
					"message_id":   j,
					"timestamp":    time.Now().UnixNano(),
				}).Info("Concurrent logging test")
			}
			done <- true
		}(i)
	}
	
	// 等待所有 goroutine 完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	duration := time.Since(start)
	
	fmt.Printf("并发日志: %d 个协程, 总计 %d 条消息, 耗时 %v\n", 
		numGoroutines, totalMessages, duration)
	fmt.Printf("平均每条消息: %v\n", duration/time.Duration(totalMessages))
	fmt.Printf("每秒处理消息: %.0f 条/秒\n", float64(totalMessages)/duration.Seconds())
}

// 额外的基准测试函数
func BenchmarkBasicLogging(log *logger.Logger, iterations int) time.Duration {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		log.Info("Benchmark message %d", i)
	}
	return time.Since(start)
}

func BenchmarkStructuredLogging(log *logger.Logger, iterations int) time.Duration {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		log.WithField("iteration", i).
			WithField("benchmark", true).
			Info("Structured benchmark message")
	}
	return time.Since(start)
}