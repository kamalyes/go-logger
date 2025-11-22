/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 01:54:27
 * @FilePath: \go-logger\formatter_test.go
 * @Description: 格式化器测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

// FormatterTestSuite 格式化器测试套件
type FormatterTestSuite struct {
	suite.Suite
	textFormatter IFormatter
	jsonFormatter IFormatter
	registry      *FormatRegistry
}

func (suite *FormatterTestSuite) SetupTest() {
	suite.textFormatter = NewTextFormatter()
	suite.jsonFormatter = NewJSONFormatter()
	suite.registry = NewFormatRegistry()
}

func TestFormatterTestSuite(t *testing.T) {
	suite.Run(t, new(FormatterTestSuite))
}

// TestNewTextFormatter 测试创建文本格式化器
func (suite *FormatterTestSuite) TestNewTextFormatter() {
	formatter := NewTextFormatter()
	assert.NotNil(suite.T(), formatter)
	assert.Equal(suite.T(), "text", formatter.GetName())
}

// TestNewJSONFormatter 测试创建JSON格式化器
func (suite *FormatterTestSuite) TestNewJSONFormatter() {
	formatter := NewJSONFormatter()
	assert.NotNil(suite.T(), formatter)
	assert.Equal(suite.T(), "json", formatter.GetName())
}

// TestTextFormatterFormat 测试文本格式化器格式化
func (suite *FormatterTestSuite) TestTextFormatterFormat() {
	entry := &LogEntry{
		Level:     INFO,
		Message:   "test message",
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
		Fields:    map[string]interface{}{},
		Caller:    &CallerInfo{File: "main.go", Line: 10, Function: "main"},
	}

	result, err := suite.textFormatter.Format(entry)
	assert.NoError(suite.T(), err)

	resultStr := string(result)
	assert.Contains(suite.T(), resultStr, "test message")
	assert.Contains(suite.T(), resultStr, "[INFO]")
	assert.Contains(suite.T(), resultStr, "2024-01-01")
}

// TestTextFormatterWithFields 测试文本格式化器处理字段
func (suite *FormatterTestSuite) TestTextFormatterWithFields() {
	event := &LogEntry{
		Level:     ERROR,
		Message:   "error occurred",
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
		Fields: map[string]interface{}{
			"user_id": 123,
			"action":  "login",
			"error":   "invalid password",
		},
		Caller: &CallerInfo{File: "auth.go", Line: 25, Function: "authenticate"},
	}

	result, err := suite.textFormatter.Format(event)
	assert.NoError(suite.T(), err)

	resultStr := string(result)
	assert.Contains(suite.T(), resultStr, "error occurred")
	assert.Contains(suite.T(), resultStr, "[ERROR]")
	assert.Contains(suite.T(), resultStr, "user_id=123")
	assert.Contains(suite.T(), resultStr, "action=login")
	assert.Contains(suite.T(), resultStr, "error=invalid password")
	assert.Contains(suite.T(), resultStr, "auth.go:25:authenticate")
}

// TestJSONFormatterFormat 测试JSON格式化器格式化
func (suite *FormatterTestSuite) TestJSONFormatterFormat() {
	entry := &LogEntry{
		Level:     WARN,
		Message:   "warning message",
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
		Fields:    map[string]interface{}{},
		Caller:    &CallerInfo{File: "service.go", Line: 50, Function: "processRequest"},
	}

	result, err := suite.jsonFormatter.Format(entry)
	assert.NoError(suite.T(), err)

	// 解析JSON以验证格式
	var logData map[string]interface{}
	err = json.Unmarshal(result, &logData)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "WARN", logData["level"])
	assert.Equal(suite.T(), "warning message", logData["message"])
	assert.NotEmpty(suite.T(), logData["timestamp"])
}

// TestJSONFormatterWithFields 测试JSON格式化器处理字段
func (suite *FormatterTestSuite) TestJSONFormatterWithFields() {
	entry := &LogEntry{
		Level:     DEBUG,
		Message:   "debug info",
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
		Fields: map[string]interface{}{
			"request_id": "abc123",
			"duration":   "150ms",
			"status":     200,
			"nested": map[string]interface{}{
				"key": "value",
			},
		},
		Caller: &CallerInfo{File: "handler.go", Line: 75, Function: "handleRequest"},
	}

	result, err := suite.jsonFormatter.Format(entry)
	assert.NoError(suite.T(), err)

	// 解析JSON以验证格式
	var logData map[string]interface{}
	err = json.Unmarshal(result, &logData)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "DEBUG", logData["level"])
	assert.Equal(suite.T(), "debug info", logData["message"])

	// 检查fields对象
	fields, ok := logData["fields"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "abc123", fields["request_id"])
	assert.Equal(suite.T(), "150ms", fields["duration"])
	assert.Equal(suite.T(), float64(200), fields["status"]) // JSON数字解析为float64

	// 测试嵌套对象
	nested, ok := fields["nested"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "value", nested["key"])
}

// TestTextFormatterColorSupport 测试文本格式化器颜色支持
func (suite *FormatterTestSuite) TestTextFormatterColorSupport() {
	formatter := NewTextFormatter()

	// 测试不同级别的颜色
	levels := []LogLevel{DEBUG, INFO, WARN, ERROR, FATAL}

	for _, level := range levels {
		entry := &LogEntry{
			Level:     level,
			Message:   "test message",
			Timestamp: time.Now().Unix(),
			Fields:    map[string]interface{}{},
		}

		result, err := formatter.Format(entry)
		assert.NoError(suite.T(), err)
		resultStr := string(result)
		assert.Contains(suite.T(), resultStr, "test message")
	}
}

// TestJSONFormatterLevelMapping 测试JSON格式化器级别映射
func (suite *FormatterTestSuite) TestJSONFormatterLevelMapping() {
	levels := map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}

	for level, expectedStr := range levels {
		entry := &LogEntry{
			Level:     level,
			Message:   "test",
			Timestamp: time.Now().Unix(),
			Fields:    map[string]interface{}{},
		}

		result, err := suite.jsonFormatter.Format(entry)
		assert.NoError(suite.T(), err)

		var logData map[string]interface{}
		err = json.Unmarshal(result, &logData)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), expectedStr, logData["level"])
	}
}

// TestFormatterEmptyFields 测试格式化器处理空字段
func (suite *FormatterTestSuite) TestFormatterEmptyFields() {
	entry := &LogEntry{
		Level:     INFO,
		Message:   "test message",
		Timestamp: time.Now().Unix(),
		Fields:    nil, // 空字段
	}

	// 测试文本格式化器
	textResult, err := suite.textFormatter.Format(entry)
	assert.NoError(suite.T(), err)
	textStr := string(textResult)
	assert.Contains(suite.T(), textStr, "test message")
	assert.Contains(suite.T(), textStr, "[INFO]")

	// 测试JSON格式化器
	jsonResult, err := suite.jsonFormatter.Format(entry)
	assert.NoError(suite.T(), err)
	var logData map[string]interface{}
	err = json.Unmarshal(jsonResult, &logData)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test message", logData["message"])
	assert.Equal(suite.T(), "INFO", logData["level"])
}

// TestFormatterSpecialCharacters 测试格式化器处理特殊字符
func (suite *FormatterTestSuite) TestFormatterSpecialCharacters() {
	entry := &LogEntry{
		Level:     INFO,
		Message:   "测试消息 with \"quotes\" and\nnewlines",
		Timestamp: time.Now().Unix(),
		Fields: map[string]interface{}{
			"special": "value with\ttabs and\nnewlines",
			"unicode": "测试🎉",
		},
	}

	// 测试文本格式化器
	textResult, err := suite.textFormatter.Format(entry)
	assert.NoError(suite.T(), err)
	textStr := string(textResult)
	assert.Contains(suite.T(), textStr, "测试消息")
	assert.Contains(suite.T(), textStr, "测试🎉")

	// 测试JSON格式化器
	jsonResult, err := suite.jsonFormatter.Format(entry)
	assert.NoError(suite.T(), err)
	var logData map[string]interface{}
	err = json.Unmarshal(jsonResult, &logData)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), logData["message"].(string), "测试消息")

	// 检查fields对象中的unicode字段
	fields, ok := logData["fields"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "测试🎉", fields["unicode"])
}

// TestFormatRegistry 测试格式化器注册表
func (suite *FormatterTestSuite) TestFormatRegistry() {
	registry := NewFormatRegistry()
	assert.NotNil(suite.T(), registry)
}

// TestFormatRegistryCreate 测试注册表创建格式化器
func (suite *FormatterTestSuite) TestFormatRegistryCreate() {
	// 测试创建文本格式化器
	textFormatter, err := suite.registry.Create(TextFormatter)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), textFormatter)
	assert.Equal(suite.T(), "text", textFormatter.GetName())

	// 测试创建JSON格式化器
	jsonFormatter, err := suite.registry.Create(JSONFormatter)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), jsonFormatter)
	assert.Equal(suite.T(), "json", jsonFormatter.GetName())

	// 测试创建未知类型
	_, err = suite.registry.Create("unknown")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "unknown formatter type")
}

// TestFormatRegistryRegister 测试注册表注册自定义格式化器
func (suite *FormatterTestSuite) TestFormatRegistryRegister() {
	// 创建自定义格式化器工厂
	customFactory := func() IFormatter {
		return NewTextFormatter() // 简单返回文本格式化器作为示例
	}

	// 注册自定义格式化器
	suite.registry.Register("custom", customFactory)

	// 测试创建自定义格式化器
	customFormatter, err := suite.registry.Create("custom")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), customFormatter)
}

// TestFormatterPerformance 测试格式化器性能
func (suite *FormatterTestSuite) TestFormatterPerformance() {
	entry := &LogEntry{
		Level:     INFO,
		Message:   "performance test message",
		Timestamp: time.Now().Unix(),
		Fields: map[string]interface{}{
			"key1": "value1",
			"key2": 12345,
			"key3": true,
		},
		Caller: &CallerInfo{File: "test.go", Line: 100, Function: "testFunc"},
	}

	// 测试文本格式化器性能
	start := time.Now()
	for i := 0; i < 1000; i++ {
		_, err := suite.textFormatter.Format(entry)
		assert.NoError(suite.T(), err)
	}
	textDuration := time.Since(start)

	// 测试JSON格式化器性能
	start = time.Now()
	for i := 0; i < 1000; i++ {
		_, err := suite.jsonFormatter.Format(entry)
		assert.NoError(suite.T(), err)
	}
	jsonDuration := time.Since(start)

	// 性能测试只是确保没有异常，不做严格的时间限制
	assert.True(suite.T(), textDuration > 0)
	assert.True(suite.T(), jsonDuration > 0)
}

// TestFormatterConcurrency 测试格式化器并发安全
func (suite *FormatterTestSuite) TestFormatterConcurrency() {
	entry := &LogEntry{
		Level:     INFO,
		Message:   "concurrency test",
		Timestamp: time.Now().Unix(),
		Fields: map[string]interface{}{
			"thread": "test",
		},
	}

	// 并发测试文本格式化器
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				result, err := suite.textFormatter.Format(entry)
				assert.NoError(suite.T(), err)
				assert.Contains(suite.T(), string(result), "concurrency test")
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 并发测试JSON格式化器
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				result, err := suite.jsonFormatter.Format(entry)
				assert.NoError(suite.T(), err)
				var logData map[string]interface{}
				err = json.Unmarshal(result, &logData)
				assert.NoError(suite.T(), err)
			}
			done <- true
		}()
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestFormatterTimeFormat 测试格式化器时间格式
func (suite *FormatterTestSuite) TestFormatterTimeFormat() {
	timestamp := time.Date(2024, 12, 25, 15, 30, 45, 123456789, time.UTC)
	entry := &LogEntry{
		Level:     INFO,
		Message:   "time format test",
		Timestamp: timestamp.Unix(),
		Fields:    map[string]interface{}{},
	}

	// 测试文本格式化器时间格式
	textResult, err := suite.textFormatter.Format(entry)
	assert.NoError(suite.T(), err)
	textStr := string(textResult)
	assert.Contains(suite.T(), textStr, "time format test")

	// 测试JSON格式化器时间格式
	jsonResult, err := suite.jsonFormatter.Format(entry)
	assert.NoError(suite.T(), err)
	var logData map[string]interface{}
	err = json.Unmarshal(jsonResult, &logData)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), logData["timestamp"])
}
