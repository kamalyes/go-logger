/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-19 00:00:00
 * @FilePath: \go-logger\return_error_test.go
 * @Description: 返回错误的日志方法测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDebugReturn 测试 DebugReturn 方法
func TestDebugReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	config.Level = DEBUG
	log := NewLogger(config)

	err := log.DebugReturn("测试错误: %s", "调试信息")

	assert.NotNil(t, err, "DebugReturn 应该返回错误")
	assert.Equal(t, "测试错误: 调试信息", err.Error(), "错误信息不匹配")

	output := buf.String()
	assert.Contains(t, output, "DEBUG", "输出应该包含 DEBUG 级别")
	assert.Contains(t, output, "测试错误: 调试信息", "输出应该包含错误信息")
}

// TestInfoReturn 测试 InfoReturn 方法
func TestInfoReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	err := log.InfoReturn("信息提示: %s", "操作完成")

	assert.NotNil(t, err, "InfoReturn 应该返回错误")
	assert.Equal(t, "信息提示: 操作完成", err.Error(), "错误信息不匹配")

	output := buf.String()
	assert.Contains(t, output, "INFO", "输出应该包含 INFO 级别")
}

// TestWarnReturn 测试 WarnReturn 方法
func TestWarnReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	err := log.WarnReturn("警告: 磁盘使用率 %d%%", 85)

	assert.NotNil(t, err, "WarnReturn 应该返回错误")
	assert.Equal(t, "警告: 磁盘使用率 85%", err.Error(), "错误信息不匹配")

	output := buf.String()
	assert.Contains(t, output, "WARN", "输出应该包含 WARN 级别")
}

// TestErrorReturn 测试 ErrorReturn 方法
func TestErrorReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	err := log.ErrorReturn("数据库连接失败: %s", "timeout")

	assert.NotNil(t, err, "ErrorReturn 应该返回错误")
	assert.Equal(t, "数据库连接失败: timeout", err.Error(), "错误信息不匹配")

	output := buf.String()
	assert.Contains(t, output, "ERROR", "输出应该包含 ERROR 级别")
}

// TestDebugCtxReturn 测试带上下文的 DebugCtxReturn 方法
func TestDebugCtxReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	config.Level = DEBUG
	log := NewLogger(config)

	ctx := context.Background()
	err := log.DebugCtxReturn(ctx, "上下文调试: %s", "测试")

	if err == nil {
		t.Error("DebugCtxReturn 应该返回错误")
	}

	if err.Error() != "上下文调试: 测试" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestInfoCtxReturn 测试带上下文的 InfoCtxReturn 方法
func TestInfoCtxReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	ctx := context.Background()
	err := log.InfoCtxReturn(ctx, "请求处理: %s", "成功")

	if err == nil {
		t.Error("InfoCtxReturn 应该返回错误")
	}

	if err.Error() != "请求处理: 成功" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestWarnCtxReturn 测试带上下文的 WarnCtxReturn 方法
func TestWarnCtxReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	ctx := context.Background()
	err := log.WarnCtxReturn(ctx, "缓存未命中: key=%s", "user:123")

	if err == nil {
		t.Error("WarnCtxReturn 应该返回错误")
	}

	if err.Error() != "缓存未命中: key=user:123" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestErrorCtxReturn 测试带上下文的 ErrorCtxReturn 方法
func TestErrorCtxReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	ctx := context.Background()
	err := log.ErrorCtxReturn(ctx, "认证失败: user=%s", "admin")

	if err == nil {
		t.Error("ErrorCtxReturn 应该返回错误")
	}

	if err.Error() != "认证失败: user=admin" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestDebugKVReturn 测试带键值对的 DebugKVReturn 方法
func TestDebugKVReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	config.Level = DEBUG
	log := NewLogger(config)

	err := log.DebugKVReturn("调试操作", "action", "test", "status", "pending")

	if err == nil {
		t.Error("DebugKVReturn 应该返回错误")
	}

	if err.Error() != "调试操作" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}

	output := buf.String()
	if !strings.Contains(output, "DEBUG") {
		t.Error("输出应该包含 DEBUG 级别")
	}
}

// TestInfoKVReturn 测试带键值对的 InfoKVReturn 方法
func TestInfoKVReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	err := log.InfoKVReturn("用户登录", "user_id", "12345", "ip", "192.168.1.1")

	if err == nil {
		t.Error("InfoKVReturn 应该返回错误")
	}

	if err.Error() != "用户登录" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestWarnKVReturn 测试带键值对的 WarnKVReturn 方法
func TestWarnKVReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	err := log.WarnKVReturn("性能警告", "duration", "5s", "threshold", "3s")

	if err == nil {
		t.Error("WarnKVReturn 应该返回错误")
	}

	if err.Error() != "性能警告" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestErrorKVReturn 测试带键值对的 ErrorKVReturn 方法
func TestErrorKVReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	err := log.ErrorKVReturn("数据库错误", "table", "users", "operation", "insert", "error", "duplicate key")

	if err == nil {
		t.Error("ErrorKVReturn 应该返回错误")
	}

	if err.Error() != "数据库错误" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}

	output := buf.String()
	if !strings.Contains(output, "ERROR") {
		t.Error("输出应该包含 ERROR 级别")
	}
}

// TestGlobalDebugReturn 测试全局 DebugReturn 方法
func TestGlobalDebugReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	config.Level = DEBUG
	SetGlobalConfig(config)

	err := DebugReturn("全局调试: %s", "测试")

	if err == nil {
		t.Error("全局 DebugReturn 应该返回错误")
	}

	if err.Error() != "全局调试: 测试" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestGlobalErrorReturn 测试全局 ErrorReturn 方法
func TestGlobalErrorReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	SetGlobalConfig(config)

	err := ErrorReturn("全局错误: code=%d", 500)

	if err == nil {
		t.Error("全局 ErrorReturn 应该返回错误")
	}

	if err.Error() != "全局错误: code=500" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestGlobalErrorCtxReturn 测试全局 ErrorCtxReturn 方法
func TestGlobalErrorCtxReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	SetGlobalConfig(config)

	ctx := context.Background()
	err := ErrorCtxReturn(ctx, "请求失败: %s", "超时")

	if err == nil {
		t.Error("全局 ErrorCtxReturn 应该返回错误")
	}

	if err.Error() != "请求失败: 超时" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestGlobalErrorKVReturn 测试全局 ErrorKVReturn 方法
func TestGlobalErrorKVReturn(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	SetGlobalConfig(config)

	err := ErrorKVReturn("连接失败", "host", "localhost", "port", 3306)

	if err == nil {
		t.Error("全局 ErrorKVReturn 应该返回错误")
	}

	if err.Error() != "连接失败" {
		t.Errorf("错误信息不匹配，得到: %s", err.Error())
	}
}

// TestReturnErrorChaining 测试返回错误的链式调用
func TestReturnErrorChaining(t *testing.T) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	// 模拟业务流程
	if err := step1(log); err != nil {
		if err2 := step2(log, err); err2 != nil {
			t.Logf("捕获到最终错误: %v", err2)
		}
	}
}

func step1(log *Logger) error {
	return log.ErrorReturn("步骤1失败: %s", "数据验证错误")
}

func step2(log *Logger, prevErr error) error {
	return log.ErrorReturn("步骤2失败，原因: %v", prevErr)
}

// BenchmarkErrorReturn 基准测试 ErrorReturn 方法
func BenchmarkErrorReturn(b *testing.B) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = log.ErrorReturn("基准测试错误: %d", i)
	}
}

// BenchmarkErrorCtxReturn 基准测试 ErrorCtxReturn 方法
func BenchmarkErrorCtxReturn(b *testing.B) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = log.ErrorCtxReturn(ctx, "基准测试错误: %d", i)
	}
}

// BenchmarkErrorKVReturn 基准测试 ErrorKVReturn 方法
func BenchmarkErrorKVReturn(b *testing.B) {
	buf := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buf
	log := NewLogger(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = log.ErrorKVReturn("基准测试", "iteration", i, "status", "running")
	}
}
