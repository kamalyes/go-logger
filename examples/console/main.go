/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-19 00:00:00
 * @FilePath: \go-logger\examples\console\main.go
 * @Description: Console 风格日志完整示例 - 分组、表格、计时器、折叠功能
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package main

import (
	"time"

	"github.com/kamalyes/go-logger"
)

func main() {
	// 创建日志器
	log := logger.NewLogger(logger.DefaultConfig())

	println("╔════════════════════════════════════════════════════════════════╗")
	println("║          Console 风格日志功能完整演示                            ║")
	println("╚════════════════════════════════════════════════════════════════╝\n")

	// ============================================================================
	// 示例 1: 基本分组
	// ============================================================================
	basicGroupExample(log)

	// ============================================================================
	// 示例 2: 嵌套分组
	// ============================================================================
	nestedGroupExample(log)

	// ============================================================================
	// 示例 3: 折叠分组功能
	// ============================================================================
	collapsedGroupExample(log)

	// ============================================================================
	// 示例 4: 表格展示
	// ============================================================================
	tableExample(log)

	// ============================================================================
	// 示例 5: 计时器
	// ============================================================================
	timerExample(log)

	// ============================================================================
	// 示例 6: 复杂场景 - API 请求处理
	// ============================================================================
	apiRequestExample(log)

	// ============================================================================
	// 示例 7: 折叠在实际场景中的应用
	// ============================================================================
	collapsedPracticalExample(log)

	// ============================================================================
	// 示例 8: 使用全局方法
	// ============================================================================
	globalMethodsExample()

	println("\n╔════════════════════════════════════════════════════════════════╗")
	println("║                   演示完成                                      ║")
	println("╚════════════════════════════════════════════════════════════════╝")
}

// basicGroupExample 基本分组示例
func basicGroupExample(log *logger.Logger) {
	println("【示例 1: 基本分组】")
	cg := log.NewConsoleGroup()

	cg.Group("用户登录流程")
	cg.Info("接收登录请求")
	cg.Debug("验证用户名和密码")
	cg.Info("登录成功，生成 Token")
	cg.GroupEnd()

	println() // 空行分隔
}

// nestedGroupExample 嵌套分组示例
func nestedGroupExample(log *logger.Logger) {
	println("【示例 2: 嵌套分组】")
	cg := log.NewConsoleGroup()

	cg.Group("订单处理系统")
	cg.Info("开始处理订单批次")

	cg.Group("订单验证")
	cg.Info("检查库存")
	cg.Info("验证用户积分")
	cg.Info("计算优惠金额")
	cg.GroupEnd()

	cg.Group("支付处理")
	cg.Info("调用支付网关")
	cg.Debug("支付参数: amount=299.00, currency=CNY")
	cg.Info("支付成功")
	cg.GroupEnd()

	cg.Group("订单确认")
	cg.Info("更新订单状态")
	cg.Info("发送确认邮件")
	cg.GroupEnd()

	cg.Info("订单处理完成")
	cg.GroupEnd()

	println()
}

// collapsedGroupExample 折叠分组示例
func collapsedGroupExample(log *logger.Logger) {
	println("【示例 3: 折叠分组功能】")
	cg := log.NewConsoleGroup()

	println("→ 正常分组 - 所有日志都会显示")
	cg.Group("📦 正常分组示例")
	cg.Info("这是普通的 Info 日志")
	cg.Debug("这是普通的 Debug 日志")
	cg.Warn("这是普通的 Warn 日志")
	cg.Error("这是普通的 Error 日志")
	cg.GroupEnd()

	println("\n→ 折叠分组 - 只有 Error 和 Fatal 级别会显示")
	cg.GroupCollapsed("📦 折叠分组示例")
	cg.Info("这条 Info 日志不会显示（已折叠）")
	cg.Debug("这条 Debug 日志不会显示（已折叠）")
	cg.Warn("这条 Warn 日志不会显示（已折叠）")
	cg.Error("❌ 这条 Error 日志会显示（即使在折叠状态）")
	cg.GroupEnd()

	println()
}

// tableExample 表格展示示例
func tableExample(log *logger.Logger) {
	println("【示例 4: 表格展示】")
	cg := log.NewConsoleGroup()

	// 示例 1: 用户列表
	println("→ Map 切片表格")
	cg.Group("用户管理")
	users := []map[string]interface{}{
		{"ID": 1, "姓名": "张三", "年龄": 25, "部门": "技术部", "职位": "工程师"},
		{"ID": 2, "姓名": "李四", "年龄": 30, "部门": "产品部", "职位": "产品经理"},
		{"ID": 3, "姓名": "王五", "年龄": 28, "部门": "技术部", "职位": "架构师"},
		{"ID": 4, "姓名": "赵六", "年龄": 26, "部门": "运营部", "职位": "运营专员"},
	}
	cg.Table(users)
	cg.GroupEnd()

	// 示例 2: 配置信息
	println("\n→ Map 键值对表格")
	cg.Group("系统配置")
	config := map[string]interface{}{
		"数据库类型":   "MySQL",
		"主机地址":    "localhost",
		"端口":      3306,
		"数据库名":    "production",
		"连接池大小":   100,
		"超时时间(秒)": 30,
		"是否启用SSL": true,
	}
	cg.Table(config)
	cg.GroupEnd()

	// 示例 3: 服务状态
	println("\n→ 字符串二维数组表格")
	cg.Group("微服务健康检查")
	services := [][]string{
		{"服务名称", "实例", "状态", "CPU", "内存", "响应时间"},
		{"user-service", "192.168.1.10:8080", "✅ 健康", "15%", "512MB", "23ms"},
		{"order-service", "192.168.1.11:8081", "✅ 健康", "25%", "768MB", "45ms"},
		{"payment-service", "192.168.1.12:8082", "⚠️  警告", "85%", "1.2GB", "156ms"},
		{"notification-service", "192.168.1.13:8083", "❌ 异常", "5%", "256MB", "超时"},
	}
	cg.Table(services)
	cg.GroupEnd()

	println()
}

// timerExample 计时器示例
func timerExample(log *logger.Logger) {
	println("【示例 5: 计时器】")
	cg := log.NewConsoleGroup()

	cg.Group("性能测试")

	// 基本计时
	println("→ 基本计时")
	timer1 := cg.Time("数据库查询")
	time.Sleep(120 * time.Millisecond)
	timer1.End()

	// 带中间日志的计时
	println("\n→ 带中间检查点的计时")
	timer2 := cg.Time("文件处理")
	time.Sleep(50 * time.Millisecond)
	timer2.Log("已处理 1000 条记录")
	time.Sleep(50 * time.Millisecond)
	timer2.Log("已处理 2000 条记录")
	time.Sleep(50 * time.Millisecond)
	timer2.End()

	// 嵌套计时
	println("\n→ 嵌套计时")
	totalTimer := cg.Time("总耗时")

	cg.Group("子任务")
	subTimer1 := cg.Time("子任务1")
	time.Sleep(80 * time.Millisecond)
	subTimer1.End()

	subTimer2 := cg.Time("子任务2")
	time.Sleep(60 * time.Millisecond)
	subTimer2.End()
	cg.GroupEnd()

	totalTimer.End()
	cg.GroupEnd()

	println()
}

// apiRequestExample API 请求处理示例
func apiRequestExample(log *logger.Logger) {
	println("【示例 6: 复杂场景 - API 请求处理】")
	cg := log.NewConsoleGroup()

	cg.Group("🌐 API 请求: GET /api/users")
	requestTimer := cg.Time("请求总耗时")

	// 请求信息
	cg.Group("📋 请求信息")
	requestInfo := map[string]interface{}{
		"Method":     "GET",
		"Path":       "/api/users",
		"Query":      "page=1&limit=10",
		"User-Agent": "Mozilla/5.0",
		"IP":         "192.168.1.100",
	}
	cg.Table(requestInfo)
	cg.GroupEnd()

	// 中间件处理
	cg.Group("🔧 中间件处理")
	cg.Info("✅ 认证中间件通过")
	cg.Info("✅ 权限验证通过")
	cg.Info("✅ 限流检查通过")
	cg.GroupEnd()

	// 业务处理
	cg.Group("💼 业务处理")
	dbTimer := cg.Time("数据库查询")
	time.Sleep(85 * time.Millisecond)
	dbTimer.End()

	// 查询结果
	users := []map[string]interface{}{
		{"ID": 1, "Name": "张三", "Email": "zhangsan@example.com", "Status": "Active"},
		{"ID": 2, "Name": "李四", "Email": "lisi@example.com", "Status": "Active"},
		{"ID": 3, "Name": "王五", "Email": "wangwu@example.com", "Status": "Inactive"},
	}
	cg.Table(users)
	cg.GroupEnd()

	// 响应信息
	cg.Group("📤 响应信息")
	responseInfo := map[string]interface{}{
		"Status Code":  200,
		"Content-Type": "application/json",
		"Records":      3,
		"Cache-Hit":    false,
	}
	cg.Table(responseInfo)
	cg.GroupEnd()

	requestTimer.End()
	cg.Info("✅ 请求处理完成")
	cg.GroupEnd()

	println()
}

// collapsedPracticalExample 折叠功能在实际场景中的应用
func collapsedPracticalExample(log *logger.Logger) {
	println("【示例 7: 折叠功能实际应用 - 应用启动流程】")
	cg := log.NewConsoleGroup()

	cg.Group("🚀 应用启动流程")
	cg.Info("开始启动应用...")

	// 详细的初始化日志可以折叠
	cg.GroupCollapsed("🔧 配置加载（详细日志已折叠）")
	cg.Info("加载配置文件 config.yaml")
	cg.Debug("解析配置项: database")
	cg.Debug("解析配置项: redis")
	cg.Debug("解析配置项: logging")
	cg.Info("配置加载完成")
	cg.GroupEnd()

	cg.Info("配置验证通过")

	// 数据库连接日志保持展开
	cg.Group("🗄️  数据库连接")
	cg.Info("连接 MySQL: localhost:3306")
	cg.Info("连接池初始化完成")
	cg.GroupEnd()

	// 详细的健康检查可以折叠
	cg.GroupCollapsed("🏥 健康检查（详细日志已折叠）")
	cg.Info("检查数据库连接...")
	cg.Debug("数据库响应时间: 5ms")
	cg.Info("检查 Redis 连接...")
	cg.Debug("Redis 响应时间: 2ms")
	cg.Info("检查磁盘空间...")
	cg.Debug("可用空间: 50GB")
	cg.Info("所有检查通过")
	cg.GroupEnd()

	// 表格在折叠分组中的应用
	println("\n→ 折叠分组中的表格")
	cg.Group("📊 用户统计报告")

	// 详细数据表格可以折叠
	cg.GroupCollapsed("📋 详细用户列表（已折叠，不显示表格）")
	detailUsers := []map[string]interface{}{
		{"ID": 1, "Name": "张三", "Age": 25, "Department": "技术部"},
		{"ID": 2, "Name": "李四", "Age": 30, "Department": "产品部"},
		{"ID": 3, "Name": "王五", "Age": 28, "Department": "运营部"},
	}
	cg.Table(detailUsers) // 这个表格不会显示（在折叠分组中）
	cg.Info("总计 %d 个用户", len(detailUsers))
	cg.GroupEnd()

	// 摘要信息保持展开
	cg.Group("📈 统计摘要（展开，显示表格）")
	summary := map[string]interface{}{
		"总用户数": 3,
		"活跃用户": 2,
		"新增用户": 1,
	}
	cg.Table(summary) // 这个表格会显示
	cg.GroupEnd()

	cg.GroupEnd()

	// 错误处理演示
	println("\n→ 折叠分组中的错误依然可见")
	cg.Group("🔄 数据处理任务")

	cg.GroupCollapsed("🔍 数据验证（已折叠）")
	cg.Info("验证字段1")
	cg.Info("验证字段2")
	cg.Error("❌ 字段3验证失败：格式不正确") // Error 会显示
	cg.Info("验证字段4")
	cg.GroupEnd()

	cg.Info("⚠️  发现 1 个错误，请检查日志")
	cg.GroupEnd()

	cg.Info("✅ 应用启动成功！")
	cg.GroupEnd()

	println("\n💡 折叠功能使用建议:")
	println("  1. 使用 GroupCollapsed() 隐藏详细的调试信息")
	println("  2. 重要的流程信息使用 Group() 保持可见")
	println("  3. Error 和 Fatal 级别的日志即使在折叠状态也会显示")
	println("  4. 可以减少日志噪音，专注于关键信息")
	println()
}

// globalMethodsExample 全局方法示例
func globalMethodsExample() {
	println("【示例 8: 全局便捷方法】")

	// 使用全局 Group
	println("→ 全局分组")
	cg := logger.Group("全局分组示例")
	cg.Info("这是使用全局方法创建的分组")
	cg.Debug("可以直接使用，无需创建 logger 实例")
	cg.GroupEnd()

	// 使用全局 Table
	println("\n→ 全局表格")
	logger.Table(map[string]interface{}{
		"功能":   "全局表格",
		"便捷性": "⭐⭐⭐⭐⭐",
		"性能":   "优秀",
	})

	// 使用全局 Timer
	println("\n→ 全局计时器")
	logger.Time("全局计时器")
	time.Sleep(100 * time.Millisecond)
	logger.TimeLog("全局计时器", "中间检查点")
	time.Sleep(100 * time.Millisecond)
	logger.TimeEnd("全局计时器")

	println()
}
