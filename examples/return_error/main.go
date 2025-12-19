/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-19 00:00:00
 * @FilePath: \go-logger\examples\return_error\main.go
 * @Description: 演示返回错误的日志方法使用示例
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"context"
	"fmt"

	"github.com/kamalyes/go-logger"
)

func main() {
	// 创建一个logger实例
	log := logger.New()

	fmt.Println("=== 基本返回错误示例 ===")
	basicReturnError(log)

	fmt.Println("\n=== 上下文返回错误示例 ===")
	contextReturnError(log)

	fmt.Println("\n=== 键值对返回错误示例 ===")
	kvReturnError(log)

	fmt.Println("\n=== 实际业务场景示例 ===")
	businessExample(log)

	fmt.Println("\n=== 全局方法示例 ===")
	globalMethodExample()
}

// basicReturnError 演示基本的返回错误方法
func basicReturnError(log *logger.Logger) {
	// 使用 ErrorReturn 记录错误并返回
	if err := processData(""); err != nil {
		fmt.Printf("收到错误: %v\n", err)
	}

	// 使用 WarnReturn 记录警告并返回
	if err := validateInput(5); err != nil {
		fmt.Printf("收到警告: %v\n", err)
	}

	// 使用 InfoReturn 记录信息并返回
	if err := notifyUser("user123"); err != nil {
		fmt.Printf("收到通知: %v\n", err)
	}
}

// contextReturnError 演示带上下文的返回错误方法
func contextReturnError(log *logger.Logger) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-12345")

	// 使用 ErrorCtxReturn
	if err := fetchData(ctx, "user-id-999"); err != nil {
		fmt.Printf("获取数据失败: %v\n", err)
	}

	// 使用 WarnCtxReturn
	if err := cacheData(ctx, "cache-key", "value"); err != nil {
		fmt.Printf("缓存警告: %v\n", err)
	}
}

// kvReturnError 演示键值对返回错误方法
func kvReturnError(log *logger.Logger) {
	// 使用 ErrorKVReturn 带结构化字段
	if err := dbOperation("users", 100); err != nil {
		fmt.Printf("数据库操作错误: %v\n", err)
	}

	// 使用 WarnKVReturn 带结构化字段
	if err := apiCall("https://api.example.com", 429); err != nil {
		fmt.Printf("API调用警告: %v\n", err)
	}
}

// businessExample 实际业务场景示例
func businessExample(log *logger.Logger) {
	// 场景1: 数据验证失败
	user := map[string]interface{}{
		"id":   "",
		"name": "John",
	}

	if err := validateUser(user); err != nil {
		// 错误已经被记录，直接返回给调用者
		fmt.Printf("用户验证失败: %v\n", err)
		return
	}

	// 场景2: 业务逻辑错误
	if err := transferMoney("user1", "user2", -100); err != nil {
		fmt.Printf("转账失败: %v\n", err)
		return
	}

	fmt.Println("业务操作成功")
}

// globalMethodExample 演示全局方法
func globalMethodExample() {
	// 使用全局的返回错误方法
	if err := logger.ErrorReturn("全局错误: %s", "系统繁忙"); err != nil {
		fmt.Printf("全局错误捕获: %v\n", err)
	}

	if err := logger.WarnReturn("全局警告: %s, 重试次数: %d", "连接超时", 3); err != nil {
		fmt.Printf("全局警告捕获: %v\n", err)
	}

	ctx := context.Background()
	if err := logger.ErrorCtxReturn(ctx, "全局上下文错误: %s", "请求失败"); err != nil {
		fmt.Printf("全局上下文错误捕获: %v\n", err)
	}

	if err := logger.ErrorKVReturn("数据库连接失败", "host", "localhost", "port", 5432); err != nil {
		fmt.Printf("全局KV错误捕获: %v\n", err)
	}
}

// ========== 辅助函数 ==========

func processData(data string) error {
	if data == "" {
		return logger.ErrorReturn("数据为空，无法处理")
	}
	return nil
}

func validateInput(value int) error {
	if value < 10 {
		return logger.WarnReturn("输入值 %d 小于推荐值 10", value)
	}
	return nil
}

func notifyUser(userID string) error {
	return logger.InfoReturn("已发送通知给用户: %s", userID)
}

func fetchData(ctx context.Context, userID string) error {
	// 模拟数据获取失败
	return logger.ErrorCtxReturn(ctx, "无法获取用户数据: userID=%s", userID)
}

func cacheData(ctx context.Context, key, value string) error {
	// 模拟缓存警告
	return logger.WarnCtxReturn(ctx, "缓存服务响应缓慢: key=%s", key)
}

func dbOperation(table string, recordCount int) error {
	return logger.ErrorKVReturn(
		"数据库查询超时",
		"table", table,
		"record_count", recordCount,
		"timeout", "5s",
	)
}

func apiCall(url string, statusCode int) error {
	return logger.WarnKVReturn(
		"API限流警告",
		"url", url,
		"status_code", statusCode,
		"retry_after", "60s",
	)
}

func validateUser(user map[string]interface{}) error {
	if user["id"] == "" {
		return logger.ErrorReturn("用户ID不能为空")
	}
	return nil
}

func transferMoney(from, to string, amount float64) error {
	if amount <= 0 {
		return logger.ErrorReturn("转账金额必须大于0，当前金额: %.2f", amount)
	}
	return nil
}
