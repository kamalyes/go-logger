/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\examples\basic\main.go
 * @Description: 基本使用示例
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"os"

	"github.com/kamalyes/go-logger"
)

func main() {
	// 1. 使用默认配置创建日志器
	log := logger.NewLogger(logger.DefaultConfig())
	log.Info("欢迎使用 go-logger!")

	// 2. 创建自定义配置的日志器
	config := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithShowCaller(true).
		WithColorful(true).
		WithPrefix("[MyApp] ")

	customLogger := logger.NewLogger(config)
	
	// 3. 基本日志级别使用
	customLogger.Debug("这是调试信息")
	customLogger.Info("这是普通信息")
	customLogger.Warn("这是警告信息")
	customLogger.Error("这是错误信息")

	// 4. 结构化日志
	customLogger.WithField("user_id", 12345).
		WithField("action", "login").
		Info("用户登录")

	// 5. 多字段日志
	customLogger.WithFields(map[string]interface{}{
		"method":     "POST",
		"url":        "/api/users",
		"status":     200,
		"duration":   "45ms",
		"request_id": "abc123",
	}).Info("API 请求处理完成")

	// 6. 错误日志
	err := someFunction()
	if err != nil {
		customLogger.WithError(err).Error("函数执行失败")
	}

	// 7. 链式调用
	customLogger.WithField("component", "database").
		WithField("operation", "query").
		WithField("table", "users").
		Debug("执行数据库查询")

	// 8. 克隆日志器
	dbLogger := customLogger.WithField("module", "database")
	dbLogger.Info("数据库连接建立")
	dbLogger.Info("数据库查询完成")
}

func someFunction() error {
	// 模拟一个可能失败的函数
	if _, err := os.Stat("nonexistent.txt"); err != nil {
		return err
	}
	return nil
}