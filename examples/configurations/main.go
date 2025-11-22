/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 13:32:32
 * @FilePath: \go-logger\examples\configurations\main.go
 * @Description: 配置功能演示
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package main

import (
	"fmt"
	logger "github.com/kamalyes/go-logger"
	"time"
)

func main() {
	fmt.Println("=== 配置功能演示 ===")

	// 1. 基础配置演示
	basicConfigDemo()

	// 2. 性能配置演示
	performanceConfigDemo()

	// 3. 企业功能配置演示
	enterpriseConfigDemo()

	// 4. 环境配置演示
	environmentConfigDemo()

	// 5. 高级配置演示
	advancedConfigDemo()

	// 6. 配置文件示例
	configFileExample()
}

// 基础配置演示
func basicConfigDemo() {
	fmt.Println("\n--- 基础配置演示 ---")

	// 使用默认配置创建适配器
	defaultConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     nil, // 使用默认输出
		TimeFormat: "2006-01-02 15:04:05",
		Colorful:   true,
	}
	fmt.Printf("默认配置: Level=%s, TimeFormat=%s, Colorful=%t\n",
		defaultConfig.Level, defaultConfig.TimeFormat, defaultConfig.Colorful)

	defaultAdapter, _ := logger.NewStandardAdapter(defaultConfig)
	defaultAdapter.Initialize()
	defaultAdapter.Info("使用默认配置的日志")
	defaultAdapter.Debug("这条调试信息不会显示（级别是INFO）")

	// 链式配置创建适配器
	chainConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.DEBUG,
		Output:     nil,
		TimeFormat: time.RFC3339,
		Colorful:   true,
		ShowCaller: true,
		Fields: map[string]interface{}{
			"service": "demo",
			"version": "1.0.0",
		},
	}

	fmt.Printf("链式配置: Level=%s, ShowCaller=%t, Fields=%v\n",
		chainConfig.Level, chainConfig.ShowCaller, chainConfig.Fields)

	chainAdapter, _ := logger.NewStandardAdapter(chainConfig)
	chainAdapter.Initialize()
	chainAdapter.Debug("调试信息 - 现在可以看到了")
	chainAdapter.Info("带有预设字段的信息日志")
	chainAdapter.WithField("user_id", 12345).Info("动态添加字段的日志")

	defer defaultAdapter.Close()
	defer chainAdapter.Close()
}

// 性能配置演示
func performanceConfigDemo() {
	fmt.Println("\n--- 性能配置演示 ---")

	// 高性能配置
	highPerfConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     nil,
		TimeFormat: "15:04:05", // 简短时间格式提高性能
		Colorful:   false,      // 禁用颜色提高性能
		ShowCaller: false,      // 禁用调用者信息提高性能
		Fields: map[string]interface{}{
			"mode": "high-performance",
		},
	}

	fmt.Printf("高性能配置: TimeFormat=%s, Colorful=%t, ShowCaller=%t\n",
		highPerfConfig.TimeFormat, highPerfConfig.Colorful, highPerfConfig.ShowCaller)

	perfAdapter, _ := logger.NewStandardAdapter(highPerfConfig)
	perfAdapter.Initialize()

	// 性能测试 - 快速写入大量日志
	start := time.Now()
	for i := 0; i < 1000; i++ {
		perfAdapter.Info("高性能日志 #%d", i)
	}
	duration := time.Since(start)
	fmt.Printf("写入1000条日志耗时: %v (%.2f μs/log)\n",
		duration, float64(duration.Nanoseconds())/1000.0/1000.0)

	defer perfAdapter.Close()
}

// 企业功能配置演示
func enterpriseConfigDemo() {
	fmt.Println("\n--- 企业功能配置演示 ---")

	// JSON格式配置，适合企业日志收集
	enterpriseConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     nil,
		Format:     "json",
		TimeFormat: time.RFC3339,
		Colorful:   false,
		ShowCaller: true,
		Fields: map[string]interface{}{
			"service":     "enterprise-app",
			"version":     "2.0.0",
			"environment": "production",
			"datacenter":  "us-east-1",
		},
	}

	fmt.Printf("企业配置: Format=%s, Fields=%v\n",
		enterpriseConfig.Format, enterpriseConfig.Fields)

	entAdapter, _ := logger.NewStandardAdapter(enterpriseConfig)
	entAdapter.Initialize()

	entAdapter.Info("企业应用启动")
	entAdapter.WithFields(map[string]interface{}{
		"user_id":    "user_12345",
		"request_id": "req_abc123",
		"action":     "user_login",
	}).Info("用户登录事件")

	entAdapter.WithFields(map[string]interface{}{
		"error_code": "DB_TIMEOUT",
		"table":      "users",
		"duration":   "5000ms",
	}).Error("数据库超时错误")

	defer entAdapter.Close()
}

// 环境配置演示
func environmentConfigDemo() {
	fmt.Println("\n--- 环境配置演示 ---")

	// 开发环境配置
	devConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.DEBUG,
		Output:     nil,
		Format:     "text",
		TimeFormat: "15:04:05",
		Colorful:   true,
		ShowCaller: true,
		Fields: map[string]interface{}{
			"environment": "development",
			"debug_mode":  true,
		},
	}
	fmt.Printf("开发环境: Level=%s, ShowCaller=%t, Colorful=%t\n",
		devConfig.Level, devConfig.ShowCaller, devConfig.Colorful)

	devAdapter, _ := logger.NewStandardAdapter(devConfig)
	devAdapter.Initialize()
	devAdapter.Debug("开发环境调试信息")
	defer devAdapter.Close()

	// 测试环境配置
	testConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.INFO,
		Output:     nil,
		Format:     "json",
		TimeFormat: time.RFC3339,
		Colorful:   false,
		ShowCaller: true,
		Fields: map[string]interface{}{
			"environment": "testing",
			"test_suite":  "integration",
		},
	}
	fmt.Printf("测试环境: Level=%s, Format=%s\n",
		testConfig.Level, testConfig.Format)

	testAdapter, _ := logger.NewStandardAdapter(testConfig)
	testAdapter.Initialize()
	testAdapter.Info("测试环境信息日志")
	defer testAdapter.Close()

	// 生产环境配置
	prodConfig := &logger.AdapterConfig{
		Type:       logger.StandardAdapter,
		Level:      logger.WARN,
		Output:     nil,
		Format:     "json",
		TimeFormat: time.RFC3339,
		Colorful:   false,
		ShowCaller: false, // 生产环境禁用调用者信息
		Fields: map[string]interface{}{
			"environment": "production",
			"service":     "app",
			"version":     "1.0.0",
		},
	}
	fmt.Printf("生产环境: Level=%s, ShowCaller=%t\n",
		prodConfig.Level, prodConfig.ShowCaller)

	prodAdapter, _ := logger.NewStandardAdapter(prodConfig)
	prodAdapter.Initialize()
	prodAdapter.Warn("生产环境警告")
	prodAdapter.Error("生产环境错误")
	defer prodAdapter.Close()
}

// 高级配置演示
func advancedConfigDemo() {
	fmt.Println("\n--- 高级配置演示 ---")

	// 多适配器配置示例
	configs := []*logger.AdapterConfig{
		{
			Type:       logger.StandardAdapter,
			Level:      logger.DEBUG,
			Output:     nil,
			Format:     "text",
			TimeFormat: "15:04:05",
			Colorful:   true,
			ShowCaller: true,
			Fields: map[string]interface{}{
				"adapter": "console",
			},
		},
		{
			Type:       logger.StandardAdapter,
			Level:      logger.INFO,
			Output:     nil,
			Format:     "json",
			TimeFormat: time.RFC3339,
			Colorful:   false,
			ShowCaller: false,
			Fields: map[string]interface{}{
				"adapter": "structured",
				"service": "advanced-demo",
			},
		},
	}

	var adapters []logger.IAdapter
	for i, config := range configs {
		fmt.Printf("适配器 #%d: Format=%s, Level=%s\n",
			i+1, config.Format, config.Level)

		adapter, _ := logger.NewStandardAdapter(config)
		adapter.Initialize()
		adapters = append(adapters, adapter)

		// 每个适配器记录不同的消息
		adapter.Info("适配器 #%d 已初始化", i+1)
	}

	// 同步写入测试
	fmt.Println("\n同步写入测试:")
	for i, adapter := range adapters {
		adapter.WithField("test_type", "sync").Info("同步测试消息 #%d", i+1)
	}

	// 带字段写入测试
	fmt.Println("\n结构化字段测试:")
	for i, adapter := range adapters {
		adapter.WithFields(map[string]interface{}{
			"test_id":   fmt.Sprintf("test_%d", i+1),
			"timestamp": time.Now().Unix(),
			"thread_id": i + 1,
			"memory_mb": float64(i*10 + 50),
			"status":    "active",
		}).Info("结构化数据测试")
	}

	// 错误级别测试
	fmt.Println("\n错误级别测试:")
	for i, adapter := range adapters {
		adapter.WithField("error_code", "TEST_ERROR").Error("模拟错误 #%d", i+1)
	}

	// 清理所有适配器
	for _, adapter := range adapters {
		defer adapter.Close()
	}
}

// 配置文件示例
func configFileExample() {
	fmt.Println("\n--- 配置文件示例 ---")

	// JSON配置示例
	jsonConfig := `{
  "level": "info",
  "time_format": "iso8601",
  "colorful": false,
  "async_write": true,
  "buffer_size": 8192,
  "pool_size": 20,
  "show_caller": false,
  "fields": {
    "service": "my-service",
    "version": "1.0.0"
  },
  "monitoring": {
    "enabled": true,
    "memory": {
      "enabled": true,
      "interval": "1m",
      "threshold": 104857600
    },
    "performance": {
      "enabled": true,
      "track_latency": true,
      "track_throughput": true,
      "sample_rate": 0.1
    }
  }
}`
	fmt.Println("JSON配置示例:")
	fmt.Println(jsonConfig)

	// YAML配置示例
	yamlConfig := `# 日志配置
level: info
time_format: iso8601
colorful: false
async_write: true
buffer_size: 8192

# 字段配置
fields:
  service: my-service
  version: 1.0.0

# 监控配置
monitoring:
  enabled: true
  memory:
    enabled: true
    interval: 1m
    threshold: 104857600
  performance:
    enabled: true
    track_latency: true
    sample_rate: 0.1

# 适配器配置
adapters:
  - name: console
    type: standard
    level: debug
    colorful: true
    
  - name: file
    type: file
    level: info
    file: ./logs/app.log
    max_size: 104857600
    max_backups: 5
`
	fmt.Println("\nYAML配置示例:")
	fmt.Println(yamlConfig)
}
