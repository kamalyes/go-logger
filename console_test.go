/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-19 23:02:35
 * @FilePath: \go-logger\console_test.go
 * @Description: Console 分组和表格功能测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"testing"
	"time"
)

func TestConsoleGroup(t *testing.T) {
	logger := NewLogger(DefaultConfig())
	cg := logger.NewConsoleGroup()

	// 测试基本分组
	cg.Group("用户信息处理")
	cg.Info("开始处理用户 ID: %d", 1001)
	cg.Debug("验证用户权限")
	cg.Info("权限验证通过")
	cg.GroupEnd()

	// 测试嵌套分组
	cg.Group("数据库操作")
	cg.Info("连接数据库")

	cg.Group("查询操作")
	cg.Info("执行查询: SELECT * FROM users")
	cg.Debug("返回 10 条记录")
	cg.GroupEnd()

	cg.Group("更新操作")
	cg.Info("执行更新: UPDATE users SET status=1")
	cg.Info("影响 5 条记录")
	cg.GroupEnd()

	cg.GroupEnd()
}

func TestConsoleGroupCollapsed(t *testing.T) {
	logger := NewLogger(DefaultConfig())
	cg := logger.NewConsoleGroup()

	// 测试折叠分组
	cg.GroupCollapsed("详细日志")
	cg.Debug("这是一些详细的调试信息")
	cg.Debug("通常不需要查看")
	cg.GroupEnd()

	cg.Info("主要流程继续...")
}

func TestTable(t *testing.T) {
	logger := NewLogger(DefaultConfig())
	cg := logger.NewConsoleGroup()

	// 测试 map 切片表格
	t.Run("MapSliceTable", func(t *testing.T) {
		users := []map[string]interface{}{
			{"ID": 1, "Name": "张三", "Age": 25, "Role": "Admin"},
			{"ID": 2, "Name": "李四", "Age": 30, "Role": "User"},
			{"ID": 3, "Name": "王五", "Age": 28, "Role": "User"},
		}

		cg.Group("用户列表")
		cg.Table(users)
		cg.GroupEnd()
	})

	// 测试单个 map 表格
	t.Run("MapTable", func(t *testing.T) {
		config := map[string]interface{}{
			"database":  "mysql",
			"host":      "localhost",
			"port":      3306,
			"username":  "root",
			"pool_size": 10,
		}

		cg.Group("配置信息")
		cg.Table(config)
		cg.GroupEnd()
	})

	// 测试字符串二维数组表格
	t.Run("StringSliceTable", func(t *testing.T) {
		data := [][]string{
			{"服务名称", "状态", "响应时间", "错误率"},
			{"API Gateway", "运行中", "45ms", "0.01%"},
			{"Auth Service", "运行中", "23ms", "0.00%"},
			{"Database", "运行中", "12ms", "0.00%"},
			{"Redis", "警告", "156ms", "0.05%"},
		}

		cg.Group("服务监控")
		cg.Table(data)
		cg.GroupEnd()
	})
}

func TestGlobalMethods(t *testing.T) {
	// 测试全局 Group
	cg := Group("全局分组测试")
	cg.Info("这是全局方法创建的分组")
	cg.Debug("支持各种日志级别")
	cg.GroupEnd()

	// 测试全局 Table
	data := map[string]interface{}{
		"version": "1.0.0",
		"env":     "production",
		"region":  "cn-east-1",
	}
	Table(data)
}

func TestNestedGroupsWithTable(t *testing.T) {
	logger := NewLogger(DefaultConfig())
	cg := logger.NewConsoleGroup()

	cg.Group("电商订单处理")
	cg.Info("开始处理订单批次")

	cg.Group("订单验证")
	orders := []map[string]interface{}{
		{"OrderID": "ORD001", "Status": "待支付", "Amount": 299.00},
		{"OrderID": "ORD002", "Status": "已支付", "Amount": 599.00},
		{"OrderID": "ORD003", "Status": "已发货", "Amount": 199.00},
	}
	cg.Table(orders)
	cg.Info("验证完成，共 %d 个订单", len(orders))
	cg.GroupEnd()

	cg.Group("支付处理")
	cg.Info("处理待支付订单")
	cg.Debug("调用支付网关 API")
	cg.Info("支付成功")
	cg.GroupEnd()

	cg.Group("物流处理")
	logistics := []map[string]interface{}{
		{"TrackingNo": "SF123456", "Carrier": "顺丰", "Status": "运输中"},
		{"TrackingNo": "YTO789012", "Carrier": "圆通", "Status": "已送达"},
	}
	cg.Table(logistics)
	cg.GroupEnd()

	cg.Info("订单处理完成")
	cg.GroupEnd()
}

func TestTimer(t *testing.T) {
	logger := NewLogger(DefaultConfig())
	cg := logger.NewConsoleGroup()

	// 测试基本计时
	t.Run("BasicTimer", func(t *testing.T) {
		timer := cg.Time("数据库查询")
		time.Sleep(100 * time.Millisecond)
		timer.End()
	})

	// 测试 TimeLog
	t.Run("TimeLog", func(t *testing.T) {
		timer := cg.Time("文件处理")
		time.Sleep(50 * time.Millisecond)
		timer.Log("已处理 50%%")
		time.Sleep(50 * time.Millisecond)
		timer.Log("已处理 100%%")
		timer.End()
	})

	// 测试嵌套计时
	t.Run("NestedTimer", func(t *testing.T) {
		cg.Group("API 请求处理")
		totalTimer := cg.Time("总耗时")

		cg.Info("验证请求参数")
		time.Sleep(20 * time.Millisecond)

		dbTimer := cg.Time("数据库操作")
		time.Sleep(80 * time.Millisecond)
		dbTimer.End()

		cacheTimer := cg.Time("缓存更新")
		time.Sleep(30 * time.Millisecond)
		cacheTimer.End()

		totalTimer.End()
		cg.GroupEnd()
	})
}

func TestGlobalTimer(t *testing.T) {
	// 测试全局计时器
	Time("全局任务")
	time.Sleep(100 * time.Millisecond)
	TimeLog("全局任务", "中间检查点")
	time.Sleep(100 * time.Millisecond)
	TimeEnd("全局任务")
}

func TestComplexScenario(t *testing.T) {
	logger := NewLogger(DefaultConfig())
	cg := logger.NewConsoleGroup()

	cg.Group("🚀 微服务启动流程")
	startTimer := cg.Time("启动总耗时")

	// 配置加载
	cg.Group("📋 配置加载")
	configTimer := cg.Time("配置加载")
	config := map[string]interface{}{
		"service_name": "user-service",
		"port":         8080,
		"environment":  "production",
		"log_level":    "info",
	}
	cg.Table(config)
	time.Sleep(50 * time.Millisecond)
	configTimer.End()
	cg.GroupEnd()

	// 数据库连接
	cg.Group("🗄️  数据库初始化")
	dbTimer := cg.Time("数据库连接")
	cg.Info("连接到 MySQL: localhost:3306")
	time.Sleep(100 * time.Millisecond)
	dbTimer.End()

	// 显示连接池状态
	poolStats := []map[string]interface{}{
		{"连接池": "主库", "最大连接数": 100, "当前连接": 5, "空闲连接": 95},
		{"连接池": "从库1", "最大连接数": 50, "当前连接": 2, "空闲连接": 48},
		{"连接池": "从库2", "最大连接数": 50, "当前连接": 3, "空闲连接": 47},
	}
	cg.Table(poolStats)
	cg.GroupEnd()

	// Redis 连接
	cg.Group("🔴 Redis 初始化")
	redisTimer := cg.Time("Redis 连接")
	cg.Info("连接到 Redis: localhost:6379")
	time.Sleep(30 * time.Millisecond)
	redisTimer.End()
	cg.GroupEnd()

	// 服务注册
	cg.Group("📡 服务注册")
	cg.Info("注册到 Consul")
	services := [][]string{
		{"服务名称", "地址", "健康检查", "状态"},
		{"user-service", "192.168.1.10:8080", "HTTP /health", "✅ 健康"},
		{"order-service", "192.168.1.11:8081", "HTTP /health", "✅ 健康"},
		{"payment-service", "192.168.1.12:8082", "HTTP /health", "⚠️  警告"},
	}
	cg.Table(services)
	cg.GroupEnd()

	startTimer.End()
	cg.Info("✅ 服务启动完成")
	cg.GroupEnd()
}

func TestConsoleGroupWithContext(t *testing.T) {
	logger := NewLogger(DefaultConfig())
	cg := logger.NewConsoleGroup()
	ctx := context.Background()

	// 测试带 Context 的日志方法
	t.Run("ContextMethods", func(t *testing.T) {
		cg.Group("带上下文的日志测试")

		cg.InfoContext(ctx, "这是带 Context 的 Info 日志")
		cg.DebugContext(ctx, "这是带 Context 的 Debug 日志")
		cg.WarnContext(ctx, "这是带 Context 的 Warn 日志")
		cg.ErrorContext(ctx, "这是带 Context 的 Error 日志")

		cg.GroupEnd()
	})

	// 测试在折叠分组中使用 Context 方法
	t.Run("ContextInCollapsedGroup", func(t *testing.T) {
		cg.GroupCollapsed("折叠分组中的 Context 方法")

		cg.InfoContext(ctx, "这条不会显示（折叠状态）")
		cg.DebugContext(ctx, "这条不会显示（折叠状态）")
		cg.WarnContext(ctx, "这条不会显示（折叠状态）")
		cg.ErrorContext(ctx, "这条 Error 会显示（即使折叠）")

		cg.GroupEnd()
	})

	// 测试嵌套分组中的 Context 方法
	t.Run("NestedContextGroups", func(t *testing.T) {
		cg.Group("API 请求处理 (带 Context)")
		cg.InfoContext(ctx, "收到请求: GET /api/users")

		cg.Group("参数验证")
		cg.DebugContext(ctx, "验证参数: page=1, limit=10")
		cg.InfoContext(ctx, "参数验证通过")
		cg.GroupEnd()

		cg.Group("业务处理")
		cg.InfoContext(ctx, "查询数据库")
		cg.DebugContext(ctx, "SQL: SELECT * FROM users LIMIT 10")
		cg.InfoContext(ctx, "查询完成，返回 10 条记录")
		cg.GroupEnd()

		cg.InfoContext(ctx, "请求处理完成")
		cg.GroupEnd()
	})
}

// TestTableWithMixedContent 测试包含混合内容的表格（中英文、表情、符号、不同长度）
func TestTableWithMixedContent(t *testing.T) {
	logger := NewLogger(DefaultConfig())
	cg := logger.NewConsoleGroup()

	cg.Group("📊 混合内容表格测试")

	// 测试1: 不同长度的中英文混合
	t.Run("中英文混合", func(t *testing.T) {
		data := map[string]interface{}{
			"用户名":          "张三 (Zhang San)",
			"Email":        "zhangsan@example.com",
			"手机号":          "+86 138-1234-5678",
			"地址":           "北京市朝阳区建国路88号SOHO现代城A座2501室",
			"Status":       "Active ✓",
			"会员等级":         "💎 Diamond VIP",
			"积分":           "12,345",
			"注册时间":         "2023-01-15 14:30:25",
			"最后登录":         "2025-12-20 19:45:30",
			"Account Type": "Premium",
		}
		cg.Info("示例1: 用户信息表（长短不一）")
		cg.Table(data)
	})

	// 测试2: 包含表情符号
	t.Run("表情符号", func(t *testing.T) {
		emojiData := map[string]interface{}{
			"🎉 活动名称": "双十二大促销",
			"📅 开始时间": "2025-12-12 00:00:00",
			"⏰ 结束时间": "2025-12-12 23:59:59",
			"💰 优惠金额": "¥500",
			"🛒 订单数":  "8,888",
			"👥 参与人数": "15,234",
			"✅ 状态":   "进行中",
			"🔥 热度":   "⭐⭐⭐⭐⭐",
			"📊 完成率":  "85.6%",
			"🎯 目标":   "10,000单",
		}
		cg.Info("示例2: 活动统计表（含表情）")
		cg.Table(emojiData)
	})

	// 测试3: 特殊符号和长文本
	t.Run("特殊符号", func(t *testing.T) {
		specialData := map[string]interface{}{
			"API接口":        "/api/v1/users/{id}/profile",
			"请求方法":         "POST → PUT → DELETE",
			"状态码":          "200 ✓ | 404 ✗ | 500 ⚠",
			"响应时间":         "≈ 125ms (avg) ± 15ms",
			"Success Rate": "99.99% ≥ 99.9%",
			"QPS":          "10K~50K req/s",
			"数据大小":         "≤ 1MB (max: 5MB)",
			"编码格式":         "UTF-8 / GBK / GB2312",
			"Content-Type": "application/json; charset=utf-8",
			"认证方式":         "Bearer Token (JWT) & API Key",
		}
		cg.Info("示例3: API接口信息（特殊符号）")
		cg.Table(specialData)
	})

	// 测试4: 极短和极长混合
	t.Run("长度差异大", func(t *testing.T) {
		lengthData := map[string]interface{}{
			"ID":           "1",
			"超长字段测试内容":     "这是一个非常非常非常非常非常非常非常非常长的字符串，用来测试表格在处理超长内容时的表现，包含中文、English、数字123、符号!@#$%^&*()以及表情😀😁😂🤣",
			"短":            "A",
			"Description":  "A comprehensive system monitoring and alerting platform with real-time data visualization",
			"中":            "测试",
			"Mixed_测试_123": "Test测试🔥",
			"URL":          "https://www.example.com/path/to/resource?param1=value1&param2=value2#section",
			"简":            "简",
			"版本号":          "v2.15.8-beta.3+build.20251220",
			"S":            "S",
		}
		cg.Info("示例4: 长度差异测试")
		cg.Table(lengthData)
	})

	// 测试5: 数值和单位混合
	t.Run("数值单位", func(t *testing.T) {
		numericData := map[string]interface{}{
			"CPU使用率": "45.8% ↑",
			"内存占用":   "8.5 GB / 16 GB",
			"磁盘空间":   "256 GB (剩余: 128 GB)",
			"网络流量 ↓": "1.25 MB/s",
			"网络流量 ↑": "850 KB/s",
			"温度":     "65°C ~ 75°C",
			"转速":     "2,400 RPM",
			"电压":     "3.3V ± 0.1V",
			"功耗":     "≈ 95W (max: 150W)",
			"运行时长":   "15天 8小时 32分钟",
		}
		cg.Info("示例5: 系统监控数据")
		cg.Table(numericData)
	})

	// 测试6: 多语言混合
	t.Run("多语言", func(t *testing.T) {
		multiLangData := map[string]interface{}{
			"中文":       "你好世界",
			"English":  "Hello World",
			"日本語":      "こんにちは世界",
			"한국어":      "안녕하세요 세계",
			"Français": "Bonjour le monde",
			"Deutsch":  "Hallo Welt",
			"Русский":  "Привет мир",
			"العربية":  "مرحبا بالعالم",
			"emoji":    "👋🌍🌎🌏",
			"混合 Mixed": "你好 Hello 世界 World 🌟",
		}
		cg.Info("示例6: 多语言支持测试")
		cg.Table(multiLangData)
	})

	cg.GroupEnd()
}
