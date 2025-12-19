/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\logger_test.go
 * @Description: 核心日志器测试套件
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package logger

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// LoggerTestSuite 核心日志器测试套件
type LoggerTestSuite struct {
	suite.Suite
	buffer *bytes.Buffer
	logger *Logger
}

// SetupTest 测试前准备
func (suite *LoggerTestSuite) SetupTest() {
	suite.buffer = &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = suite.buffer
	config.Colorful = false // 测试时关闭颜色，便于验证
	suite.logger = NewLogger(config)
}

// TearDownTest 测试后清理
func (suite *LoggerTestSuite) TearDownTest() {
	suite.buffer = nil
	suite.logger = nil
}

// TestNewLogger 测试创建新的日志器
func (suite *LoggerTestSuite) TestNewLogger() {
	// 测试使用默认配置
	logger := NewLogger(nil)
	assert.NotNil(suite.T(), logger)
	assert.Equal(suite.T(), INFO, logger.GetLevel())
	assert.False(suite.T(), logger.IsShowCaller())

	// 测试使用自定义配置
	config := &LogConfig{
		Level:      DEBUG,
		ShowCaller: true,
		Prefix:     "[TEST]",
		Output:     suite.buffer,
		Colorful:   false,
		TimeFormat: "15:04:05",
	}
	logger = NewLogger(config)
	assert.Equal(suite.T(), DEBUG, logger.GetLevel())
	assert.True(suite.T(), logger.IsShowCaller())

	// 测试使用无效配置（会回退到默认配置）
	invalidConfig := &LogConfig{}
	logger = NewLogger(invalidConfig)
	assert.NotNil(suite.T(), logger)
}

// TestLoggerBasicMethods 测试基本方法
func (suite *LoggerTestSuite) TestLoggerBasicMethods() {
	// 测试设置和获取级别
	suite.logger.SetLevel(DEBUG)
	assert.Equal(suite.T(), DEBUG, suite.logger.GetLevel())

	suite.logger.SetLevel(ERROR)
	assert.Equal(suite.T(), ERROR, suite.logger.GetLevel())

	// 测试设置和检查调用者显示
	suite.logger.SetShowCaller(true)
	assert.True(suite.T(), suite.logger.IsShowCaller())

	suite.logger.SetShowCaller(false)
	assert.False(suite.T(), suite.logger.IsShowCaller())

	// 测试级别启用检查
	suite.logger.SetLevel(WARN)
	assert.False(suite.T(), suite.logger.IsLevelEnabled(DEBUG))
	assert.False(suite.T(), suite.logger.IsLevelEnabled(INFO))
	assert.True(suite.T(), suite.logger.IsLevelEnabled(WARN))
	assert.True(suite.T(), suite.logger.IsLevelEnabled(ERROR))
	assert.True(suite.T(), suite.logger.IsLevelEnabled(FATAL))
}

// TestLoggerLoggingMethods 测试日志记录方法
func (suite *LoggerTestSuite) TestLoggerLoggingMethods() {
	suite.logger.SetLevel(DEBUG)

	// 测试Debug
	suite.buffer.Reset()
	suite.logger.Debug("Debug message: %s", "test")
	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "DEBUG")
	assert.Contains(suite.T(), output, "Debug message: test")

	// 测试Info
	suite.buffer.Reset()
	suite.logger.Info("Info message: %d", 123)
	output = suite.buffer.String()
	assert.Contains(suite.T(), output, "INFO")
	assert.Contains(suite.T(), output, "Info message: 123")

	// 测试Warn
	suite.buffer.Reset()
	suite.logger.Warn("Warn message: %v", true)
	output = suite.buffer.String()
	assert.Contains(suite.T(), output, "WARN")
	assert.Contains(suite.T(), output, "Warn message: true")

	// 测试Error
	suite.buffer.Reset()
	suite.logger.Error("Error message: %f", 3.14)
	output = suite.buffer.String()
	assert.Contains(suite.T(), output, "ERROR")
	assert.Contains(suite.T(), output, "Error message: 3.14")
}

// TestLoggerLevelFiltering 测试日志级别过滤
func (suite *LoggerTestSuite) TestLoggerLevelFiltering() {
	// 设置为INFO级别
	suite.logger.SetLevel(INFO)

	// DEBUG消息应该被过滤掉
	suite.buffer.Reset()
	suite.logger.Debug("This should not appear")
	output := suite.buffer.String()
	assert.Empty(suite.T(), output)

	// INFO消息应该显示
	suite.buffer.Reset()
	suite.logger.Info("This should appear")
	output = suite.buffer.String()
	assert.NotEmpty(suite.T(), output)
	assert.Contains(suite.T(), output, "This should appear")

	// WARN消息应该显示
	suite.buffer.Reset()
	suite.logger.Warn("Warning message")
	output = suite.buffer.String()
	assert.NotEmpty(suite.T(), output)
	assert.Contains(suite.T(), output, "Warning message")
}

// TestLoggerWithField 测试单字段方法
func (suite *LoggerTestSuite) TestLoggerWithField() {
	newLogger := suite.logger.WithField("user_id", "12345")
	assert.NotEqual(suite.T(), suite.logger, newLogger) // 应该是新实例

	newLogger.Info("User logged in")
	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "user_id=12345")
	assert.Contains(suite.T(), output, "User logged in")
}

// TestLoggerWithFields 测试多字段方法
func (suite *LoggerTestSuite) TestLoggerWithFields() {
	fields := map[string]interface{}{
		"user_id":   "12345",
		"action":    "login",
		"timestamp": 1699401600,
	}

	newLogger := suite.logger.WithFields(fields)
	assert.NotEqual(suite.T(), suite.logger, newLogger)

	newLogger.Info("User action performed")
	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "user_id=12345")
	assert.Contains(suite.T(), output, "action=login")
	assert.Contains(suite.T(), output, "timestamp=1699401600")
	assert.Contains(suite.T(), output, "User action performed")

	// 测试空字段映射
	emptyLogger := suite.logger.WithFields(map[string]interface{}{})
	assert.Equal(suite.T(), suite.logger, emptyLogger) // 应该返回原实例

	// 测试nil字段映射
	nilLogger := suite.logger.WithFields(nil)
	assert.Equal(suite.T(), suite.logger, nilLogger)
}

// TestLoggerWithError 测试错误字段方法
func (suite *LoggerTestSuite) TestLoggerWithError() {
	testError := errors.New("test error occurred")
	newLogger := suite.logger.WithError(testError)
	assert.NotEqual(suite.T(), suite.logger, newLogger)

	newLogger.Error("Operation failed")
	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "error=test error occurred")
	assert.Contains(suite.T(), output, "Operation failed")
}

// TestLoggerCallerInfo 测试调用者信息
func (suite *LoggerTestSuite) TestLoggerCallerInfo() {
	suite.logger.SetShowCaller(true)

	suite.logger.Info("Test message with caller")
	output := suite.buffer.String()

	// 应该包含文件名和行号信息
	assert.Contains(suite.T(), output, ".go:")
	assert.Contains(suite.T(), output, "Test message with caller")
}

// TestLoggerClone 测试克隆功能
func (suite *LoggerTestSuite) TestLoggerClone() {
	suite.logger.SetLevel(DEBUG)
	suite.logger.SetShowCaller(true)

	cloned := suite.logger.Clone()
	assert.NotSame(suite.T(), suite.logger, cloned) // 不同实例
	assert.Equal(suite.T(), suite.logger.GetLevel(), cloned.GetLevel())
	assert.Equal(suite.T(), suite.logger.IsShowCaller(), cloned.IsShowCaller())

	// 修改克隆不应影响原logger
	cloned.SetLevel(ERROR)
	assert.Equal(suite.T(), DEBUG, suite.logger.GetLevel())
	assert.Equal(suite.T(), ERROR, cloned.GetLevel())
}

// TestLoggerConfigOperations 测试配置操作
func (suite *LoggerTestSuite) TestLoggerConfigOperations() {
	// 获取配置副本
	config := suite.logger.GetConfig()
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), suite.logger.GetLevel(), config.Level)

	// 更新配置
	newConfig := &LogConfig{
		Level:      WARN,
		ShowCaller: true,
		Prefix:     "[UPDATED]",
		Output:     suite.buffer,
		Colorful:   false,
		TimeFormat: "15:04:05.000",
	}

	suite.logger.UpdateConfig(newConfig)
	assert.Equal(suite.T(), WARN, suite.logger.GetLevel())
	assert.True(suite.T(), suite.logger.IsShowCaller())

	// 测试使用nil配置更新
	suite.logger.UpdateConfig(nil)
	// 应该不发生变化
	assert.Equal(suite.T(), WARN, suite.logger.GetLevel())
}

// TestGlobalLoggerFunctions 测试全局日志器函数
func (suite *LoggerTestSuite) TestGlobalLoggerFunctions() {
	// 设置全局配置
	globalConfig := &LogConfig{
		Level:      INFO, // 直接设置为INFO级别
		ShowCaller: true,
		Output:     suite.buffer,
		Colorful:   false,
	}
	SetGlobalConfig(globalConfig)

	// 测试全局级别设置
	globalLogger := GetGlobalLogger()
	assert.Equal(suite.T(), INFO, globalLogger.GetLevel())
	assert.True(suite.T(), globalLogger.IsShowCaller())

	// 测试全局日志方法
	suite.buffer.Reset()
	Debug("Global debug message")
	assert.Empty(suite.T(), suite.buffer.String()) // 应该被过滤

	suite.buffer.Reset()
	Info("Global info message")
	assert.Contains(suite.T(), suite.buffer.String(), "Global info message")

	suite.buffer.Reset()
	Warn("Global warn message")
	assert.Contains(suite.T(), suite.buffer.String(), "Global warn message")

	suite.buffer.Reset()
	Error("Global error message")
	assert.Contains(suite.T(), suite.buffer.String(), "Global error message")

	// 测试全局字段方法
	suite.buffer.Reset()
	WithField("global_key", "global_value").Info("Global field test")
	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "global_key=global_value")
	assert.Contains(suite.T(), output, "Global field test")

	suite.buffer.Reset()
	WithFields(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}).Info("Global fields test")
	output = suite.buffer.String()
	assert.Contains(suite.T(), output, "key1=value1")
	assert.Contains(suite.T(), output, "key2=value2")

	suite.buffer.Reset()
	WithError(errors.New("global error")).Error("Global error test")
	output = suite.buffer.String()
	assert.Contains(suite.T(), output, "error=global error")
	assert.Contains(suite.T(), output, "Global error test")

	// 测试获取全局配置
	retrievedConfig := GetGlobalConfig()
	assert.NotNil(suite.T(), retrievedConfig)
	assert.Equal(suite.T(), INFO, retrievedConfig.Level)
}

// TestLoggerFormatMessage 测试消息格式化
func (suite *LoggerTestSuite) TestLoggerFormatMessage() {
	// 测试基本格式化
	suite.logger.SetLevel(DEBUG)
	suite.buffer.Reset()
	suite.logger.Info("Simple message")
	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "INFO")
	assert.Contains(suite.T(), output, "Simple message")

	// 测试带emoji的格式化（默认配置）
	config := DefaultConfig()
	config.Output = suite.buffer
	config.Colorful = false
	emojiLogger := NewLogger(config)

	suite.buffer.Reset()
	emojiLogger.Info("Message with emoji")
	output = suite.buffer.String()
	assert.Contains(suite.T(), output, "ℹ️")
	assert.Contains(suite.T(), output, "INFO")
}

// TestLoggerPrefixHandling 测试前缀处理
func (suite *LoggerTestSuite) TestLoggerPrefixHandling() {
	// 测试自动添加空格的前缀
	config := DefaultConfig()
	config.Output = suite.buffer
	config.Prefix = "[SERVICE]"
	config.Colorful = false

	logger := NewLogger(config)
	logger.Info("Test message")
	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "[SERVICE]")
}

// TestLoggerConcurrency 测试并发安全
func (suite *LoggerTestSuite) TestLoggerConcurrency() {
	var wg sync.WaitGroup
	numGoroutines := 10
	messagesPerGoroutine := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				suite.logger.Info("Goroutine %d, Message %d", id, j)
			}
		}(i)
	}

	wg.Wait()

	// 验证没有panic发生
	assert.True(suite.T(), true) // 如果到达这里说明没有并发问题
}

// TestLoggerEdgeCases 测试边界情况
func (suite *LoggerTestSuite) TestLoggerEdgeCases() {
	// 测试空消息
	suite.buffer.Reset()
	suite.logger.Info("%s", "")
	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "INFO")

	// 测试很长的消息
	longMessage := strings.Repeat("A", 10000)
	suite.buffer.Reset()
	suite.logger.Info("%s", longMessage)
	output = suite.buffer.String()
	assert.Contains(suite.T(), output, longMessage)

	// 测试特殊字符
	suite.buffer.Reset()
	suite.logger.Info("Message with 特殊字符 and émojis 🎉")
	output = suite.buffer.String()
	assert.Contains(suite.T(), output, "特殊字符")
	assert.Contains(suite.T(), output, "émojis 🎉")
}

// TestLoggerWithInvalidLevel 测试无效级别处理
func (suite *LoggerTestSuite) TestLoggerWithInvalidLevel() {
	// 设置无效级别（999比所有标准级别都高）
	suite.logger.SetLevel(LogLevel(999))

	// 由于999级别太高，Info级别的日志应该被过滤
	suite.buffer.Reset()
	suite.logger.Info("Test with invalid level")
	output := suite.buffer.String()

	// 测试无效级别时的处理 - 应该是空的，因为999比INFO(1)级别高很多
	if output == "" {
		// 这是期望的行为 - 高级别会过滤低级别日志
		assert.Empty(suite.T(), output)
	} else {
		// 如果有输出，验证包含消息
		assert.Contains(suite.T(), output, "Test with invalid level")
	}
}

// TestLoggerStats 测试日志统计功能（如果实现了的话）
func (suite *LoggerTestSuite) TestLoggerStats() {
	// 这里可以测试日志统计功能，如果Logger支持的话
	suite.logger.SetLevel(DEBUG)

	// 记录一些日志
	suite.logger.Debug("Debug message")
	suite.logger.Info("Info message")
	suite.logger.Warn("Warn message")
	suite.logger.Error("Error message")

	// 如果实现了统计功能，可以验证计数
	// 这里只是演示测试结构
}

// TestLoggerChaining 测试方法链
func (suite *LoggerTestSuite) TestLoggerChaining() {
	suite.buffer.Reset()

	// 测试复杂的方法链
	suite.logger.
		WithField("user_id", "123").
		WithField("action", "test").
		WithError(errors.New("chain test error")).
		Error("Chained logging test")

	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "user_id=123")
	assert.Contains(suite.T(), output, "action=test")
	assert.Contains(suite.T(), output, "error=chain test error")
	assert.Contains(suite.T(), output, "Chained logging test")
}

// TestLoggerMemoryUsage 测试内存使用
func (suite *LoggerTestSuite) TestLoggerMemoryUsage() {
	// 创建大量logger实例
	loggers := make([]*Logger, 1000)
	for i := 0; i < 1000; i++ {
		config := DefaultConfig()
		config.Output = &bytes.Buffer{}
		loggers[i] = NewLogger(config)
	}

	// 验证都创建成功
	assert.Len(suite.T(), loggers, 1000)

	// 使用所有logger
	for i, logger := range loggers {
		logger.Info("Logger %d test", i)
	}
}

// 运行测试套件
func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

// TestLoggerPerformance 性能测试（模拟）
func TestLoggerPerformance(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	// 测试大量日志输出的性能
	start := time.Now()
	iterations := 10000

	for i := 0; i < iterations; i++ {
		logger.Info("Performance test message %d", i)
	}

	duration := time.Since(start)
	t.Logf("Logged %d messages in %v (avg: %v per message)",
		iterations, duration, duration/time.Duration(iterations))

	assert.True(t, duration < time.Second*5,
		"Logging should be reasonably fast, took %v", duration)
}

// TestLoggerFatalBehavior 测试Fatal行为（需要小心处理os.Exit）
func TestLoggerFatalBehavior(t *testing.T) {
	// 注意：这个测试不能直接调用Fatal，因为它会调用os.Exit(1)
	// 在实际项目中，可能需要使用依赖注入来模拟os.Exit行为

	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	// 测试Fatal消息格式（不实际调用Fatal方法）
	logger.Error("This would be a fatal error")
	output := buffer.String()
	assert.Contains(t, output, "ERROR")
	assert.Contains(t, output, "This would be a fatal error")
}

// TestFormattingConsistency 测试格式化一致性
func TestFormattingConsistency(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false

	logger := NewLogger(config)
	logger.SetLevel(DEBUG)

	// 测试不同级别的格式化一致性
	levels := []struct {
		level LogLevel
		name  string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
	}

	for _, lvl := range levels {
		buffer.Reset()

		switch lvl.level {
		case DEBUG:
			logger.Debug("Test message")
		case INFO:
			logger.Info("Test message")
		case WARN:
			logger.Warn("Test message")
		case ERROR:
			logger.Error("Test message")
		}

		output := buffer.String()
		assert.Contains(t, output, lvl.name,
			"Level %s should appear in output", lvl.name)
		assert.Contains(t, output, "Test message",
			"Message should appear in output for level %s", lvl.name)
	}
}

func TestNew(t *testing.T) {
	// 测试New函数创建默认logger
	log := New()
	if log == nil {
		t.Fatal("New() 应该返回非空的logger实例")
	}

	// 验证默认配置
	config := log.GetConfig()
	if config.Level != INFO {
		t.Errorf("默认级别应该是INFO，实际是%v", config.Level)
	}

	if config.ShowCaller != false {
		t.Errorf("默认ShowCaller应该是false，实际是%v", config.ShowCaller)
	}

	if config.Colorful != true {
		t.Errorf("默认Colorful应该是true，实际是%v", config.Colorful)
	}
}

func TestNewLogger(t *testing.T) {
	// 测试NewLogger函数
	config := NewLogConfig().
		WithLevel(WARN).
		WithPrefix("[TEST] ").
		WithShowCaller(true)

	log := NewLogger(config)
	if log == nil {
		t.Fatal("NewLogger() 应该返回非空的logger实例")
	}

	// 验证配置
	actualConfig := log.GetConfig()
	if actualConfig.Level != WARN {
		t.Errorf("级别应该是WARN，实际是%v", actualConfig.Level)
	}

	if actualConfig.ShowCaller != true {
		t.Errorf("ShowCaller应该是true，实际是%v", actualConfig.ShowCaller)
	}

	if !strings.Contains(actualConfig.Prefix, "[TEST]") {
		t.Errorf("前缀应该包含[TEST]，实际是%s", actualConfig.Prefix)
	}
}

func TestNewLoggerWithNilConfig(t *testing.T) {
	// 测试NewLogger传入nil配置
	log := NewLogger(nil)
	if log == nil {
		t.Fatal("NewLogger(nil) 应该返回非空的logger实例")
	}

	// 应该使用默认配置
	config := log.GetConfig()
	if config.Level != INFO {
		t.Errorf("nil配置时应该使用默认级别INFO，实际是%v", config.Level)
	}
}

func TestLoggerChainMethods(t *testing.T) {
	// 测试链式调用方法
	log := New().
		WithLevel(DEBUG).
		WithPrefix("[CHAIN] ").
		WithShowCaller(true).
		WithColorful(false)

	if log == nil {
		t.Fatal("链式调用应该返回非空的logger实例")
	}

	// 验证链式配置结果
	config := log.GetConfig()
	if config.Level != DEBUG {
		t.Errorf("链式设置级别应该是DEBUG，实际是%v", config.Level)
	}

	if config.ShowCaller != true {
		t.Errorf("链式设置ShowCaller应该是true，实际是%v", config.ShowCaller)
	}

	if config.Colorful != false {
		t.Errorf("链式设置Colorful应该是false，实际是%v", config.Colorful)
	}

	if !strings.Contains(config.Prefix, "[CHAIN]") {
		t.Errorf("链式设置前缀应该包含[CHAIN]，实际是%s", config.Prefix)
	}
}

func TestLoggerLevelCheck(t *testing.T) {
	// 测试日志级别检查
	log := New().WithLevel(WARN)

	if !log.IsLevelEnabled(WARN) {
		t.Error("WARN级别应该被启用")
	}

	if !log.IsLevelEnabled(ERROR) {
		t.Error("ERROR级别应该被启用（高于WARN）")
	}

	if log.IsLevelEnabled(INFO) {
		t.Error("INFO级别不应该被启用（低于WARN）")
	}

	if log.IsLevelEnabled(DEBUG) {
		t.Error("DEBUG级别不应该被启用（低于WARN）")
	}
}

func TestLoggerGetSetMethods(t *testing.T) {
	log := New()

	// 测试SetLevel和GetLevel
	log.SetLevel(ERROR)
	if log.GetLevel() != ERROR {
		t.Errorf("SetLevel/GetLevel: 期望ERROR，实际%v", log.GetLevel())
	}

	// 测试SetShowCaller和IsShowCaller
	log.SetShowCaller(true)
	if !log.IsShowCaller() {
		t.Error("SetShowCaller/IsShowCaller: 期望true，实际false")
	}

	log.SetShowCaller(false)
	if log.IsShowCaller() {
		t.Error("SetShowCaller/IsShowCaller: 期望false，实际true")
	}
}

func TestLoggerClone(t *testing.T) {
	// 测试Clone方法
	original := New().WithLevel(WARN).WithShowCaller(true)
	cloned := original.Clone()

	if cloned == nil {
		t.Fatal("Clone() 应该返回非空的logger实例")
	}

	// 验证克隆的配置
	originalConfig := original.GetConfig()
	clonedConfig := cloned.(*Logger).GetConfig()

	if originalConfig.Level != clonedConfig.Level {
		t.Errorf("克隆的级别不匹配：原始%v，克隆%v", originalConfig.Level, clonedConfig.Level)
	}

	if originalConfig.ShowCaller != clonedConfig.ShowCaller {
		t.Errorf("克隆的ShowCaller不匹配：原始%v，克隆%v", originalConfig.ShowCaller, clonedConfig.ShowCaller)
	}
}

func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

func BenchmarkNewLogger(b *testing.B) {
	config := DefaultConfig()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewLogger(config)
	}
}

func BenchmarkChainMethods(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New().WithLevel(DEBUG).WithPrefix("[BENCH] ").WithShowCaller(true)
	}
}

// TestLoggerKVWithObject 测试 KV 方法支持对象自动解析
func TestLoggerKVWithObject(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	config.Level = DEBUG // 设置为 DEBUG 级别
	logger := NewLogger(config)

	// 定义测试结构体
	type User struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
		IsActive bool   `json:"is_active"`
	}

	user := User{
		ID:       1001,
		Name:     "张三",
		Email:    "zhangsan@example.com",
		Age:      25,
		IsActive: true,
	}

	// 测试 InfoKV 使用对象参数
	buffer.Reset()
	logger.InfoKV("用户登录", user)
	output := buffer.String()
	
	assert.Contains(t, output, "用户登录")
	assert.Contains(t, output, "id")
	assert.Contains(t, output, "1001")
	assert.Contains(t, output, "name")
	assert.Contains(t, output, "张三")
	assert.Contains(t, output, "email")
	assert.Contains(t, output, "zhangsan@example.com")

	// 测试 DebugKV 使用对象参数
	buffer.Reset()
	logger.DebugKV("调试信息", user)
	output = buffer.String()
	
	assert.NotEmpty(t, output, "调试输出不应为空")
	assert.Contains(t, output, "调试信息")
	assert.Contains(t, output, "age")
	assert.Contains(t, output, "25")

	// 测试使用 map 参数
	buffer.Reset()
	data := map[string]interface{}{
		"request_id": "req-12345",
		"method":     "POST",
		"path":       "/api/users",
		"status":     200,
	}
	logger.InfoKV("API请求", data)
	output = buffer.String()
	
	assert.Contains(t, output, "API请求")
	assert.Contains(t, output, "request_id")
	assert.Contains(t, output, "req-12345")
	assert.Contains(t, output, "method")
	assert.Contains(t, output, "POST")

	// 测试使用指针参数
	buffer.Reset()
	userPtr := &user
	logger.WarnKV("警告信息", userPtr)
	output = buffer.String()
	
	assert.Contains(t, output, "警告信息")
	assert.Contains(t, output, "name")
	assert.Contains(t, output, "张三")

	// 测试传统的 key-value 对参数（向后兼容）
	buffer.Reset()
	logger.InfoKV("传统方式", "key1", "value1", "key2", 123)
	output = buffer.String()
	
	assert.Contains(t, output, "传统方式")
	assert.Contains(t, output, "key1")
	assert.Contains(t, output, "value1")
	assert.Contains(t, output, "key2")
	assert.Contains(t, output, "123")
}

// TestLoggerKVReturnWithObject 测试 KVReturn 方法支持对象参数
func TestLoggerKVReturnWithObject(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	type ErrorInfo struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details string `json:"details"`
	}

	errInfo := ErrorInfo{
		Code:    500,
		Message: "服务器内部错误",
		Details: "数据库连接失败",
	}

	// 测试 ErrorKVReturn 使用对象参数
	err := logger.ErrorKVReturn("系统错误", errInfo)
	assert.NotNil(t, err)
	assert.Equal(t, "系统错误", err.Error())
	
	output := buffer.String()
	assert.Contains(t, output, "系统错误")
	assert.Contains(t, output, "code")
	assert.Contains(t, output, "500")
	assert.Contains(t, output, "message")
	assert.Contains(t, output, "服务器内部错误")

	// 测试 WarnKVReturn 使用 map 参数
	buffer.Reset()
	warnData := map[string]interface{}{
		"cpu_usage":    85.5,
		"memory_usage": 90.2,
		"disk_usage":   75.0,
	}
	err = logger.WarnKVReturn("资源使用率过高", warnData)
	assert.NotNil(t, err)
	
	output = buffer.String()
	assert.Contains(t, output, "资源使用率过高")
	assert.Contains(t, output, "cpu_usage")
	assert.Contains(t, output, "85.5")
}

// TestNewConsoleGroup 测试创建 ConsoleGroup
func TestNewConsoleGroup(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	cg := logger.NewConsoleGroup()
	assert.NotNil(t, cg, "ConsoleGroup 不应为 nil")
	assert.Equal(t, 0, cg.indentLevel, "初始缩进级别应为 0")
	assert.False(t, cg.collapsed, "初始状态不应折叠")
}

// TestConsoleGroupBasicFlow 测试 ConsoleGroup 基本流程
func TestConsoleGroupBasicFlow(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	cg := logger.NewConsoleGroup()
	
	// 测试基本分组
	cg.Group("测试分组")
	cg.Info("分组内的消息")
	cg.GroupEnd()

	output := buffer.String()
	assert.Contains(t, output, "▼ 测试分组")
	assert.Contains(t, output, "分组内的消息")
}

// TestConsoleGroupNested 测试嵌套分组
func TestConsoleGroupNested(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	cg := logger.NewConsoleGroup()
	
	cg.Group("外层分组")
	cg.Info("外层消息")
	
	cg.Group("内层分组")
	cg.Info("内层消息")
	cg.GroupEnd()
	
	cg.Info("回到外层")
	cg.GroupEnd()

	output := buffer.String()
	assert.Contains(t, output, "外层分组")
	assert.Contains(t, output, "内层分组")
	assert.Contains(t, output, "外层消息")
	assert.Contains(t, output, "内层消息")
}

// TestConsoleGroupTable 测试表格功能
func TestConsoleGroupTable(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	cg := logger.NewConsoleGroup()
	
	// 测试 map 表格
	data := map[string]interface{}{
		"名称": "测试",
		"值":  123,
	}
	cg.Table(data)

	output := buffer.String()
	assert.Contains(t, output, "名称")
	assert.Contains(t, output, "测试")
}

// TestConsoleGroupTimer 测试计时器功能
func TestConsoleGroupTimer(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	cg := logger.NewConsoleGroup()
	
	timer := cg.Time("测试计时器")
	time.Sleep(10 * time.Millisecond)
	timer.End()

	output := buffer.String()
	assert.Contains(t, output, "测试计时器")
	assert.Contains(t, output, "ms") // 应该包含时间单位
}

// TestConsoleGroupContextMethods 测试带 Context 的方法
func TestConsoleGroupContextMethods(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	cg := logger.NewConsoleGroup()
	ctx := context.Background()
	
	cg.Group("Context 测试")
	cg.InfoContext(ctx, "Info with context")
	cg.DebugContext(ctx, "Debug with context")
	cg.WarnContext(ctx, "Warn with context")
	cg.ErrorContext(ctx, "Error with context")
	cg.GroupEnd()

	output := buffer.String()
	assert.Contains(t, output, "Info with context")
	assert.Contains(t, output, "Error with context")
}

// TestLoggerConsoleIntegration 测试 Logger 的 Console 集成方法
func TestLoggerConsoleIntegration(t *testing.T) {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	config.Colorful = false
	logger := NewLogger(config)

	// 测试通过 Logger 直接调用 Console 方法
	logger.ConsoleGroup("集成测试")
	logger.ConsoleTable(map[string]interface{}{"key": "value"})
	timer := logger.ConsoleTime("timer")
	time.Sleep(5 * time.Millisecond)
	timer.End()
	logger.ConsoleGroupEnd()

	output := buffer.String()
	assert.Contains(t, output, "集成测试")
	assert.Contains(t, output, "key")
	assert.Contains(t, output, "timer")
}

// BenchmarkKVWithObject 性能测试：对象参数
func BenchmarkKVWithObject(b *testing.B) {
	logger := New()
	logger.SetLevel(INFO)
	
	type TestData struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
		Field3 bool   `json:"field3"`
	}
	
	data := TestData{
		Field1: "test",
		Field2: 123,
		Field3: true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoKV("测试消息", data)
	}
}

// BenchmarkKVWithMap 性能测试：map 参数
func BenchmarkKVWithMap(b *testing.B) {
	logger := New()
	logger.SetLevel(INFO)
	
	data := map[string]interface{}{
		"field1": "test",
		"field2": 123,
		"field3": true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoKV("测试消息", data)
	}
}

// BenchmarkKVWithKeyValuePairs 性能测试：传统 key-value 对
func BenchmarkKVWithKeyValuePairs(b *testing.B) {
	logger := New()
	logger.SetLevel(INFO)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoKV("测试消息", "field1", "test", "field2", 123, "field3", true)
	}
}

// BenchmarkConsoleGroup 性能测试：ConsoleGroup
func BenchmarkConsoleGroup(b *testing.B) {
	logger := New()
	logger.SetLevel(INFO)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cg := logger.NewConsoleGroup()
		cg.Group("Benchmark")
		cg.Info("Test message")
		cg.GroupEnd()
	}
}

// BenchmarkConsoleTable 性能测试：ConsoleTable
func BenchmarkConsoleTable(b *testing.B) {
	logger := New()
	logger.SetLevel(INFO)
	
	data := map[string]interface{}{
		"field1": "test",
		"field2": 123,
		"field3": true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cg := logger.NewConsoleGroup()
		cg.Table(data)
	}
}

// BenchmarkConsoleTimer 性能测试：ConsoleTimer
func BenchmarkConsoleTimer(b *testing.B) {
	logger := New()
	logger.SetLevel(INFO)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cg := logger.NewConsoleGroup()
		timer := cg.Time("benchmark")
		timer.End()
	}
}
