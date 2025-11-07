/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\factory_test.go
 * @Description: 工厂模式测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// FactoryTestSuite 工厂测试套件
type FactoryTestSuite struct {
	suite.Suite
	factory *LoggerFactory
}

func (suite *FactoryTestSuite) SetupTest() {
	suite.factory = NewLoggerFactory()
}

func TestFactoryTestSuite(t *testing.T) {
	suite.Run(t, new(FactoryTestSuite))
}

// TestNewLoggerFactory 测试创建工厂
func (suite *FactoryTestSuite) TestNewLoggerFactory() {
	factory := NewLoggerFactory()
	
	assert.NotNil(suite.T(), factory)
	assert.NotNil(suite.T(), factory.formatters)
	assert.NotNil(suite.T(), factory.adapters)
	assert.NotNil(suite.T(), factory.writers)
	assert.NotNil(suite.T(), factory.hooks)
	assert.NotNil(suite.T(), factory.middlewares)
}

// TestCreateFormatter 测试创建格式化器
func (suite *FactoryTestSuite) TestCreateFormatter() {
	// 测试文本格式化器
	textFormatter, err := suite.factory.CreateFormatter(TextFormatter)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), textFormatter)
	
	// 测试JSON格式化器
	jsonFormatter, err := suite.factory.CreateFormatter(JSONFormatter)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), jsonFormatter)
	
	// 测试无效格式化器类型
	_, err = suite.factory.CreateFormatter("invalid")
	assert.Error(suite.T(), err)
}

// TestCreateWriter 测试创建写入器
func (suite *FactoryTestSuite) TestCreateWriter() {
	// 测试控制台写入器
	consoleWriter, err := suite.factory.CreateWriter("console", nil)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), consoleWriter)
	
	// 测试文件写入器
	fileConfig := map[string]interface{}{
		"file_path": "test.log",
	}
	fileWriter, err := suite.factory.CreateWriter("file", fileConfig)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), fileWriter)
	
	// 测试文件写入器无效配置
	_, err = suite.factory.CreateWriter("file", "invalid")
	assert.Error(suite.T(), err)
	
	// 测试轮转写入器
	rotateConfig := map[string]interface{}{
		"file_path":  "rotate.log",
		"max_size":   int64(10 * 1024 * 1024),
		"max_files":  5,
	}
	rotateWriter, err := suite.factory.CreateWriter("rotate", rotateConfig)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), rotateWriter)
	
	// 测试不存在的写入器类型
	_, err = suite.factory.CreateWriter("nonexistent", nil)
	assert.Error(suite.T(), err)
}

// TestCreateHook 测试创建钩子
func (suite *FactoryTestSuite) TestCreateHook() {
	// 测试控制台钩子
	consoleHook, err := suite.factory.CreateHook("console", nil)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), consoleHook)
	
	// 测试文件钩子
	fileConfig := map[string]interface{}{
		"file_path": "hook.log",
	}
	fileHook, err := suite.factory.CreateHook("file", fileConfig)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), fileHook)
	
	// 测试文件钩子无效配置
	_, err = suite.factory.CreateHook("file", "invalid")
	assert.Error(suite.T(), err)
	
	// 测试webhook钩子
	webhookConfig := map[string]interface{}{
		"url": "http://example.com/webhook",
	}
	webhookHook, err := suite.factory.CreateHook("webhook", webhookConfig)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), webhookHook)
	
	// 测试webhook钩子无效配置
	_, err = suite.factory.CreateHook("webhook", map[string]interface{}{})
	assert.Error(suite.T(), err)
	
	// 测试不存在的钩子类型
	_, err = suite.factory.CreateHook("nonexistent", nil)
	assert.Error(suite.T(), err)
}

// TestCreateMiddleware 测试创建中间件
func (suite *FactoryTestSuite) TestCreateMiddleware() {
	// 测试认证中间件（应该失败，因为不支持）
	authConfig := map[string]interface{}{
		"required_role": "admin",
		"allowed_users": []interface{}{"user1", "user2"},
	}
	_, err := suite.factory.CreateMiddleware("auth", authConfig)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "unknown middleware type")
	
	// 测试认证中间件无效配置
	_, err = suite.factory.CreateMiddleware("auth", "invalid")
	assert.Error(suite.T(), err)
	
	// 测试丰富中间件（应该失败，因为不支持）
	_, err = suite.factory.CreateMiddleware("enrich", nil)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "unknown middleware type")
	
	// 测试指标中间件（应该失败，因为不支持）
	_, err = suite.factory.CreateMiddleware("metrics", nil)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "unknown middleware type")
	
	// 测试速率限制中间件（应该失败，因为不支持）
	ratelimitConfig := map[string]interface{}{
		"max_rate":    100,
		"time_window": "1s",
		"burst_size":  10,
	}
	_, err = suite.factory.CreateMiddleware("ratelimit", ratelimitConfig)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "unknown middleware type")
	
	// 测试不存在的中间件类型
	_, err = suite.factory.CreateMiddleware("nonexistent", nil)
	assert.Error(suite.T(), err)
}

// TestLoggerBuilder 测试日志器构建器
func (suite *FactoryTestSuite) TestLoggerBuilder() {
	builder := NewLoggerBuilder()
	assert.NotNil(suite.T(), builder)
	assert.NotNil(suite.T(), builder.factory)
	assert.NotNil(suite.T(), builder.config)
	assert.Empty(suite.T(), builder.writers)
	assert.Empty(suite.T(), builder.hooks)
	assert.Empty(suite.T(), builder.middlewares)
}

// TestLoggerBuilderChaining 测试构建器链式调用
func (suite *FactoryTestSuite) TestLoggerBuilderChaining() {
	buffer := &bytes.Buffer{}
	config := DefaultConfig().
		WithLevel(DEBUG).
		WithShowCaller(true).
		WithPrefix("[Test] ")
	config.Output = buffer
	
	logger := NewLoggerBuilder().
		WithConfig(config).
		WithFormatter(JSONFormatter).
		WithWriter("console", nil).
		WithHook("console", nil).
		WithMiddleware("metrics", nil).
		Build()
	
	assert.NotNil(suite.T(), logger)
	assert.Equal(suite.T(), DEBUG, logger.GetLevel())
	assert.True(suite.T(), logger.IsShowCaller())
}

// TestLoggerBuilderBuild 测试构建器构建日志器
func (suite *FactoryTestSuite) TestLoggerBuilderBuild() {
	buffer := &bytes.Buffer{}
	config := DefaultConfig()
	config.Output = buffer
	
	logger := NewLoggerBuilder().
		WithConfig(config).
		Build()
	
	assert.NotNil(suite.T(), logger)
	
	// 测试日志输出
	logger.Info("test message")
	assert.Contains(suite.T(), buffer.String(), "test message")
}

// TestCreateSimpleLogger 测试创建简单日志器
func (suite *FactoryTestSuite) TestCreateSimpleLogger() {
	logger := CreateSimpleLogger(INFO)
	assert.NotNil(suite.T(), logger)
	assert.Equal(suite.T(), INFO, logger.GetLevel())
}

// TestCreateFileLogger 测试创建文件日志器
func (suite *FactoryTestSuite) TestCreateFileLogger() {
	tempFile := "test_file_logger.log"
	defer os.Remove(tempFile)
	
	logger := CreateFileLogger(tempFile, DEBUG)
	assert.NotNil(suite.T(), logger)
	assert.Equal(suite.T(), DEBUG, logger.GetLevel())
}

// TestCreateProductionLogger 测试创建生产环境日志器
func (suite *FactoryTestSuite) TestCreateProductionLogger() {
	tempDir := "test_logs"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)
	
	logger := CreateProductionLogger(tempDir)
	assert.NotNil(suite.T(), logger)
	assert.Equal(suite.T(), INFO, logger.GetLevel())
	assert.True(suite.T(), logger.IsShowCaller())
}

// TestGlobalFactory 测试全局工厂
func (suite *FactoryTestSuite) TestGlobalFactory() {
	globalFactory := GetGlobalFactory()
	assert.NotNil(suite.T(), globalFactory)
	
	// 测试设置新的全局工厂
	newFactory := NewLoggerFactory()
	SetGlobalFactory(newFactory)
	
	retrievedFactory := GetGlobalFactory()
	assert.Equal(suite.T(), newFactory, retrievedFactory)
}

// TestBuilderErrorHandling 测试构建器错误处理
func (suite *FactoryTestSuite) TestBuilderErrorHandling() {
	builder := NewLoggerBuilder()
	
	// 测试无效的写入器不会导致panic
	builder.WithWriter("invalid_writer", nil)
	logger := builder.Build()
	assert.NotNil(suite.T(), logger)
	
	// 测试无效的钩子不会导致panic
	builder.WithHook("invalid_hook", nil)
	logger = builder.Build()
	assert.NotNil(suite.T(), logger)
	
	// 测试无效的中间件不会导致panic
	builder.WithMiddleware("invalid_middleware", nil)
	logger = builder.Build()
	assert.NotNil(suite.T(), logger)
}

// TestBuilderWithDefaults 测试构建器默认值
func (suite *FactoryTestSuite) TestBuilderWithDefaults() {
	logger := NewLoggerBuilder().Build()
	
	assert.NotNil(suite.T(), logger)
	assert.Equal(suite.T(), INFO, logger.GetLevel()) // 默认级别
	assert.False(suite.T(), logger.IsShowCaller())   // 默认不显示调用者
}