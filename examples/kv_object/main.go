/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-19 00:00:00
 * @FilePath: \go-logger\examples\kv_object\main.go
 * @Description: KV 方法对象参数示例 - 展示如何使用对象自动解析为键值对
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"time"

	"github.com/kamalyes/go-logger"
)

// User 用户信息结构体
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// APIRequest API 请求信息
type APIRequest struct {
	RequestID string            `json:"request_id"`
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	Headers   map[string]string `json:"headers"`
	Duration  time.Duration     `json:"duration"`
	StatusCode int              `json:"status_code"`
}

// DatabaseError 数据库错误信息
type DatabaseError struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Query     string `json:"query"`
	Table     string `json:"table"`
	Duration  int64  `json:"duration_ms"`
}

func main() {
	// 创建 logger
	log := logger.New()
	log.SetLevel(logger.DEBUG)

	log.InfoMsg("=== KV 方法对象参数示例 ===\n")

	// 示例 1: 使用结构体对象记录用户信息
	log.InfoMsg("【示例 1】使用结构体对象")
	user := User{
		ID:        1001,
		Username:  "zhangsan",
		Email:     "zhangsan@example.com",
		Age:       28,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	
	// 直接传递结构体对象，自动解析为键值对
	log.InfoKV("用户登录成功", user)
	log.DebugKV("用户详细信息", user)
	log.InfoMsg("")

	// 示例 2: 使用结构体指针
	log.InfoMsg("【示例 2】使用结构体指针")
	userPtr := &user
	log.InfoKV("用户信息查询", userPtr)
	log.InfoMsg("")

	// 示例 3: 使用 map 记录 API 请求
	log.InfoMsg("【示例 3】使用 map 对象")
	requestData := map[string]interface{}{
		"request_id": "req-12345-abcde",
		"method":     "POST",
		"path":       "/api/v1/users",
		"client_ip":  "192.168.1.100",
		"user_agent": "Mozilla/5.0",
		"status":     200,
		"duration":   "125ms",
	}
	log.InfoKV("API 请求处理完成", requestData)
	log.InfoMsg("")

	// 示例 4: 使用复杂结构体
	log.InfoMsg("【示例 4】使用复杂结构体")
	apiReq := APIRequest{
		RequestID: "req-67890-fghij",
		Method:    "GET",
		Path:      "/api/v1/products/123",
		Headers: map[string]string{
			"Authorization": "Bearer token123",
			"Content-Type":  "application/json",
		},
		Duration:   150 * time.Millisecond,
		StatusCode: 200,
	}
	log.InfoKV("API 请求详情", apiReq)
	log.InfoMsg("")

	// 示例 5: 错误日志使用对象
	log.InfoMsg("【示例 5】错误日志使用对象")
	dbError := DatabaseError{
		Code:     1062,
		Message:  "Duplicate entry 'user@example.com' for key 'email'",
		Query:    "INSERT INTO users (email, username) VALUES (?, ?)",
		Table:    "users",
		Duration: 45,
	}
	log.ErrorKV("数据库操作失败", dbError)
	log.WarnKV("数据库性能警告", dbError)
	log.InfoMsg("")

	// 示例 6: 传统方式仍然支持（向后兼容）
	log.InfoMsg("【示例 6】传统键值对方式（向后兼容）")
	log.InfoKV("传统方式",
		"key1", "value1",
		"key2", 123,
		"key3", true,
		"key4", time.Now(),
	)
	log.InfoMsg("")

	// 示例 7: 返回错误日志使用对象
	log.InfoMsg("【示例 7】返回错误日志使用对象")
	validationError := map[string]interface{}{
		"field":   "email",
		"value":   "invalid-email",
		"reason":  "格式不正确",
		"pattern": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
	}
	
	err := log.ErrorKVReturn("参数验证失败", validationError)
	if err != nil {
		log.Errorf("返回的错误: %v", err)
	}
	log.InfoMsg("")

	// 示例 8: 性能监控数据
	log.InfoMsg("【示例 8】性能监控数据")
	metrics := map[string]interface{}{
		"cpu_usage":       85.5,
		"memory_usage":    2048.75,
		"disk_usage":      65.3,
		"network_in_mb":   150.2,
		"network_out_mb":  89.7,
		"active_sessions": 1250,
		"requests_per_sec": 450,
	}
	log.InfoKV("系统性能指标", metrics)
	log.WarnKV("资源使用率较高", metrics)
	log.InfoMsg("")

	// 示例 9: 嵌套结构（map 包含结构体）
	log.InfoMsg("【示例 9】混合使用")
	mixedData := map[string]interface{}{
		"operation": "user_update",
		"user_id":   user.ID,
		"timestamp": time.Now().Unix(),
		"changes": map[string]interface{}{
			"old_email": "old@example.com",
			"new_email": user.Email,
		},
		"operator": "admin",
	}
	log.InfoKV("用户信息更新", mixedData)
	log.InfoMsg("")

	log.InfoMsg("=== 示例完成 ===")
}
