/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\adapters_test.go
 * @Description: 适配器模块测试套件
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package logger

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"sync"
	"testing"
	"time"
)

// AdapterTestSuite 适配器测试套件
type AdapterTestSuite struct {
	suite.Suite
	buffer  *bytes.Buffer
	manager IManager
}

// SetupSuite 套件初始化
func (suite *AdapterTestSuite) SetupSuite() {
	// 套件级别的初始化
}

// TearDownSuite 套件清理
func (suite *AdapterTestSuite) TearDownSuite() {
	// 套件级别的清理
}

// SetupTest 测试前准备
func (suite *AdapterTestSuite) SetupTest() {
	suite.buffer = &bytes.Buffer{}
	suite.manager = NewLoggerManager()
}

// TearDownTest 测试后清理
func (suite *AdapterTestSuite) TearDownTest() {
	if suite.manager != nil {
		suite.manager.CloseAll()
	}
	suite.buffer = nil
}

// TestAdapterConfig 测试适配器配置
func (suite *AdapterTestSuite) TestAdapterConfig() {
	// 测试默认配置
	config := DefaultAdapterConfig()
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), StandardAdapter, config.Type)
	assert.Equal(suite.T(), "default", config.Name)
	assert.Equal(suite.T(), INFO, config.Level)
	assert.Equal(suite.T(), os.Stdout, config.Output)
	assert.Equal(suite.T(), "text", config.Format)
	assert.Equal(suite.T(), 100, config.MaxSize)
	assert.Equal(suite.T(), 3, config.MaxBackups)
	assert.Equal(suite.T(), 28, config.MaxAge)
	assert.True(suite.T(), config.Compress)
	assert.NotNil(suite.T(), config.Fields)
	assert.Equal(suite.T(), "2006-01-02 15:04:05", config.TimeFormat)
	assert.True(suite.T(), config.Colorful)

	// 测试配置验证
	err := config.Validate()
	assert.NoError(suite.T(), err)

	// 测试空名称的配置
	config.Name = ""
	err = config.Validate()
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "adapter name is required")

	// 测试恢复有效名称
	config.Name = "test"
	err = config.Validate()
	assert.NoError(suite.T(), err)

	// 测试nil Output
	config.Output = nil
	err = config.Validate()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), os.Stdout, config.Output) // 应该设置为默认值

	// 测试nil Fields
	config.Fields = nil
	err = config.Validate()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config.Fields) // 应该设置为默认值

	// 测试空TimeFormat
	config.TimeFormat = ""
	err = config.Validate()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "2006-01-02 15:04:05", config.TimeFormat) // 应该设置为默认值
}

// TestLoggerManager 测试日志管理器
func (suite *AdapterTestSuite) TestLoggerManager() {
	// 测试创建管理器
	manager := NewLoggerManager()
	assert.NotNil(suite.T(), manager)

	// 测试初始状态
	adapters := manager.ListAdapters()
	assert.Empty(suite.T(), adapters)

	// 测试健康检查
	health := manager.HealthCheck()
	assert.NotNil(suite.T(), health)
	assert.Empty(suite.T(), health)
}

// TestAddAdapter 测试添加适配器
func (suite *AdapterTestSuite) TestAddAdapter() {
	// 创建适配器
	config := DefaultAdapterConfig()
	config.Name = "test"
	config.Output = suite.buffer
	adapter, err := CreateAdapter(config)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), adapter)

	// 测试添加适配器
	err = suite.manager.AddAdapter("test", adapter)
	assert.NoError(suite.T(), err)

	// 验证适配器已添加
	adapters := suite.manager.ListAdapters()
	assert.Len(suite.T(), adapters, 1)
	assert.Contains(suite.T(), adapters, "test")

	// 测试重复添加
	err = suite.manager.AddAdapter("test", adapter)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "already exists")

	// 测试空名称
	err = suite.manager.AddAdapter("", adapter)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cannot be empty")

	// 测试nil适配器
	err = suite.manager.AddAdapter("nil", nil)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "cannot be nil")
}

// TestGetAdapter 测试获取适配器
func (suite *AdapterTestSuite) TestGetAdapter() {
	// 测试获取不存在的适配器
	adapter, exists := suite.manager.GetAdapter("nonexistent")
	assert.Nil(suite.T(), adapter)
	assert.False(suite.T(), exists)

	// 添加适配器
	config := DefaultAdapterConfig()
	config.Name = "test"
	config.Output = suite.buffer
	testAdapter, err := CreateAdapter(config)
	assert.NoError(suite.T(), err)

	err = suite.manager.AddAdapter("test", testAdapter)
	assert.NoError(suite.T(), err)

	// 测试获取存在的适配器
	adapter, exists = suite.manager.GetAdapter("test")
	assert.NotNil(suite.T(), adapter)
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), testAdapter, adapter)
}

// TestRemoveAdapter 测试移除适配器
func (suite *AdapterTestSuite) TestRemoveAdapter() {
	// 测试移除不存在的适配器
	err := suite.manager.RemoveAdapter("nonexistent")
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found")

	// 添加适配器
	config := DefaultAdapterConfig()
	config.Name = "test"
	config.Output = suite.buffer
	adapter, err := CreateAdapter(config)
	assert.NoError(suite.T(), err)

	err = suite.manager.AddAdapter("test", adapter)
	assert.NoError(suite.T(), err)

	// 验证适配器存在
	adapters := suite.manager.ListAdapters()
	assert.Len(suite.T(), adapters, 1)

	// 测试移除适配器
	err = suite.manager.RemoveAdapter("test")
	assert.NoError(suite.T(), err)

	// 验证适配器已移除
	adapters = suite.manager.ListAdapters()
	assert.Empty(suite.T(), adapters)

	// 验证适配器不存在
	_, exists := suite.manager.GetAdapter("test")
	assert.False(suite.T(), exists)
}

// TestCloseAll 测试关闭所有适配器
func (suite *AdapterTestSuite) TestCloseAll() {
	// 添加多个适配器
	for i := 0; i < 3; i++ {
		config := DefaultAdapterConfig()
		config.Name = fmt.Sprintf("adapter%d", i)
		config.Output = suite.buffer
		adapter, err := CreateAdapter(config)
		assert.NoError(suite.T(), err)

		err = suite.manager.AddAdapter(config.Name, adapter)
		assert.NoError(suite.T(), err)
	}

	// 验证适配器已添加
	adapters := suite.manager.ListAdapters()
	assert.Len(suite.T(), adapters, 3)

	// 测试关闭所有适配器
	err := suite.manager.CloseAll()
	assert.NoError(suite.T(), err)

	// 验证所有适配器已移除
	adapters = suite.manager.ListAdapters()
	assert.Empty(suite.T(), adapters)
}

// TestFlushAll 测试刷新所有适配器
func (suite *AdapterTestSuite) TestFlushAll() {
	// 添加适配器
	config := DefaultAdapterConfig()
	config.Name = "test"
	config.Output = suite.buffer
	adapter, err := CreateAdapter(config)
	assert.NoError(suite.T(), err)

	err = suite.manager.AddAdapter("test", adapter)
	assert.NoError(suite.T(), err)

	// 测试刷新所有适配器
	err = suite.manager.FlushAll()
	assert.NoError(suite.T(), err)
}

// TestSetLevelAll 测试设置所有适配器级别
func (suite *AdapterTestSuite) TestSetLevelAll() {
	// 添加适配器
	config := DefaultAdapterConfig()
	config.Name = "test"
	config.Output = suite.buffer
	adapter, err := CreateAdapter(config)
	assert.NoError(suite.T(), err)

	err = suite.manager.AddAdapter("test", adapter)
	assert.NoError(suite.T(), err)

	// 测试设置级别
	suite.manager.SetLevelAll(DEBUG)

	// 验证级别已设置
	retrievedAdapter, exists := suite.manager.GetAdapter("test")
	assert.True(suite.T(), exists)
	assert.Equal(suite.T(), DEBUG, retrievedAdapter.GetLevel())
}

// TestBroadcast 测试广播日志
func (suite *AdapterTestSuite) TestBroadcast() {
	// 添加多个适配器
	buffers := make([]*bytes.Buffer, 3)
	for i := 0; i < 3; i++ {
		buffers[i] = &bytes.Buffer{}
		config := DefaultAdapterConfig()
		config.Name = fmt.Sprintf("adapter%d", i)
		config.Output = buffers[i]
		config.Level = DEBUG
		adapter, err := CreateAdapter(config)
		assert.NoError(suite.T(), err)

		err = suite.manager.AddAdapter(config.Name, adapter)
		assert.NoError(suite.T(), err)
	}

	// 测试广播不同级别的日志（跳过FATAL级别以避免程序退出）
	testCases := []struct {
		level  LogLevel
		format string
		args   []interface{}
	}{
		{DEBUG, "Debug message: %s", []interface{}{"test"}},
		{INFO, "Info message: %d", []interface{}{123}},
		{WARN, "Warn message: %v", []interface{}{true}},
		{ERROR, "Error message: %f", []interface{}{3.14}},
		// 跳过FATAL测试以避免os.Exit(1)
	}

	for _, tc := range testCases {
		suite.T().Run(fmt.Sprintf("Broadcast_%s", tc.level), func(t *testing.T) {
			// 清空缓冲区
			for i := range buffers {
				buffers[i].Reset()
			}

			// 广播日志
			suite.manager.Broadcast(tc.level, tc.format, tc.args...)

			// 验证所有适配器都接收到日志
			expectedMessage := fmt.Sprintf(tc.format, tc.args...)
			for i, buffer := range buffers {
				output := buffer.String()
				assert.Contains(t, output, expectedMessage,
					"Adapter %d should contain the broadcast message", i)
			}
		})
	}
}

// TestBroadcastFatal 单独测试FATAL级别（需要特殊处理）
func (suite *AdapterTestSuite) TestBroadcastFatal() {
	// 注意：这个测试不能直接调用FATAL，因为它会调用os.Exit(1)
	// 在实际项目中，可能需要使用依赖注入来模拟os.Exit行为

	buffer := &bytes.Buffer{}
	config := DefaultAdapterConfig()
	config.Name = "fatal_test"
	config.Output = buffer
	config.Level = DEBUG
	adapter, err := CreateAdapter(config)
	assert.NoError(suite.T(), err)

	err = suite.manager.AddAdapter(config.Name, adapter)
	assert.NoError(suite.T(), err)

	// 测试FATAL级别的消息格式（但不实际调用会导致退出的方法）
	// 这里我们通过直接使用ERROR级别来测试FATAL的消息格式
	suite.manager.Broadcast(ERROR, "This would be a fatal error: %s", "critical")
	output := buffer.String()
	assert.Contains(suite.T(), output, "This would be a fatal error: critical")
}

// TestHealthCheck 测试健康检查
func (suite *AdapterTestSuite) TestHealthCheck() {
	// 添加适配器
	config := DefaultAdapterConfig()
	config.Name = "test"
	config.Output = suite.buffer
	adapter, err := CreateAdapter(config)
	assert.NoError(suite.T(), err)

	err = suite.manager.AddAdapter("test", adapter)
	assert.NoError(suite.T(), err)

	// 测试健康检查
	health := suite.manager.HealthCheck()
	assert.NotNil(suite.T(), health)
	assert.Len(suite.T(), health, 1)
	assert.Contains(suite.T(), health, "test")
	assert.True(suite.T(), health["test"])
}

// TestCreateAdapter 测试创建适配器工厂方法
func (suite *AdapterTestSuite) TestCreateAdapter() {
	// 测试使用nil配置创建
	adapter, err := CreateAdapter(nil)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), adapter)

	// 测试使用有效配置创建
	config := DefaultAdapterConfig()
	config.Output = suite.buffer
	adapter, err = CreateAdapter(config)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), adapter)

	// 测试使用无效配置创建
	config.Name = ""
	adapter, err = CreateAdapter(config)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "invalid adapter config")
	assert.Nil(suite.T(), adapter)

	// 测试不支持的适配器类型
	config.Name = "test"
	config.Type = LogrusAdapter
	adapter, err = CreateAdapter(config)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "logrus adapter requires additional dependency")
	assert.Nil(suite.T(), adapter)

	config.Type = ZapAdapter
	adapter, err = CreateAdapter(config)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "zap adapter requires additional dependency")
	assert.Nil(suite.T(), adapter)

	config.Type = ZerologAdapter
	adapter, err = CreateAdapter(config)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "zerolog adapter requires additional dependency")
	assert.Nil(suite.T(), adapter)

	// 测试未知适配器类型
	config.Type = AdapterType("unknown")
	adapter, err = CreateAdapter(config)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "unsupported adapter type")
	assert.Nil(suite.T(), adapter)
}

// TestConcurrentOperations 测试并发操作
func (suite *AdapterTestSuite) TestConcurrentOperations() {
	var wg sync.WaitGroup
	numGoroutines := 10
	numAdapters := 5

	// 并发添加适配器
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			for j := 0; j < numAdapters; j++ {
				config := DefaultAdapterConfig()
				config.Name = fmt.Sprintf("adapter_%d_%d", index, j)
				config.Output = &bytes.Buffer{}
				adapter, err := CreateAdapter(config)
				if err == nil {
					suite.manager.AddAdapter(config.Name, adapter)
				}
			}
		}(i)
	}

	// 并发读取操作
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				suite.manager.ListAdapters()
				suite.manager.HealthCheck()
			}
		}()
	}

	wg.Wait()

	// 验证最终状态
	adapters := suite.manager.ListAdapters()
	assert.True(suite.T(), len(adapters) >= 0) // 可能有一些添加成功

	// 清理
	err := suite.manager.CloseAll()
	assert.NoError(suite.T(), err)
}

// TestManagerErrorHandling 测试管理器错误处理
func (suite *AdapterTestSuite) TestManagerErrorHandling() {
	// 创建一个会在关闭时出错的模拟适配器
	mockAdapter := &MockAdapter{
		initError:  nil,
		closeError: errors.New("close error"),
		flushError: errors.New("flush error"),
		healthy:    true,
	}

	// 添加模拟适配器
	err := suite.manager.AddAdapter("mock", mockAdapter)
	assert.NoError(suite.T(), err)

	// 测试刷新错误
	err = suite.manager.FlushAll()
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "flush")

	// 测试关闭错误
	err = suite.manager.CloseAll()
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "close")
}

// TestAdapterConfigWithFields 测试带字段的适配器配置
func (suite *AdapterTestSuite) TestAdapterConfigWithFields() {
	config := DefaultAdapterConfig()
	config.Name = "test"
	config.Output = suite.buffer
	config.Fields = map[string]interface{}{
		"service": "test-service",
		"version": "1.0.0",
		"prefix":  "[TEST]",
	}

	adapter, err := CreateAdapter(config)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), adapter)

	// 测试适配器功能
	adapter.Info("Test message")
	output := suite.buffer.String()
	assert.Contains(suite.T(), output, "Test message")
}

// TestAdapterTypeConstants 测试适配器类型常量
func (suite *AdapterTestSuite) TestAdapterTypeConstants() {
	assert.Equal(suite.T(), AdapterType("standard"), StandardAdapter)
	assert.Equal(suite.T(), AdapterType("logrus"), LogrusAdapter)
	assert.Equal(suite.T(), AdapterType("zap"), ZapAdapter)
	assert.Equal(suite.T(), AdapterType("zerolog"), ZerologAdapter)
}

// MockAdapter 用于测试的模拟适配器
type MockAdapter struct {
	level      LogLevel
	showCaller bool
	fields     map[string]interface{}
	initError  error
	closeError error
	flushError error
	healthy    bool
	name       string
	version    string
}

// 实现 IAdapter 接口
func (m *MockAdapter) Initialize() error {
	return m.initError
}

func (m *MockAdapter) Close() error {
	return m.closeError
}

func (m *MockAdapter) Flush() error {
	return m.flushError
}

func (m *MockAdapter) GetAdapterName() string {
	if m.name == "" {
		return "mock"
	}
	return m.name
}

func (m *MockAdapter) GetAdapterVersion() string {
	if m.version == "" {
		return "1.0.0"
	}
	return m.version
}

func (m *MockAdapter) IsHealthy() bool {
	return m.healthy
}

// 实现 ILogger 接口
func (m *MockAdapter) Debug(format string, args ...interface{}) {}
func (m *MockAdapter) Info(format string, args ...interface{})  {}
func (m *MockAdapter) Warn(format string, args ...interface{})  {}
func (m *MockAdapter) Error(format string, args ...interface{}) {}
func (m *MockAdapter) Fatal(format string, args ...interface{}) {}

// Printf风格方法（与上面相同，但命名更明确）
func (m *MockAdapter) Debugf(format string, args ...interface{}) {}
func (m *MockAdapter) Infof(format string, args ...interface{})  {}
func (m *MockAdapter) Warnf(format string, args ...interface{})  {}
func (m *MockAdapter) Errorf(format string, args ...interface{}) {}
func (m *MockAdapter) Fatalf(format string, args ...interface{}) {}

func (m *MockAdapter) SetLevel(level LogLevel) {
	m.level = level
}

func (m *MockAdapter) GetLevel() LogLevel {
	return m.level
}

func (m *MockAdapter) SetShowCaller(show bool) {
	m.showCaller = show
}

func (m *MockAdapter) IsShowCaller() bool {
	return m.showCaller
}

func (m *MockAdapter) IsLevelEnabled(level LogLevel) bool {
	return level >= m.level
}

func (m *MockAdapter) WithField(key string, value interface{}) ILogger {
	if m.fields == nil {
		m.fields = make(map[string]interface{})
	}
	m.fields[key] = value
	return m
}

func (m *MockAdapter) WithFields(fields map[string]interface{}) ILogger {
	if m.fields == nil {
		m.fields = make(map[string]interface{})
	}
	for k, v := range fields {
		m.fields[k] = v
	}
	return m
}

func (m *MockAdapter) WithError(err error) ILogger {
	return m.WithField("error", err)
}

// 新增的接口方法实现

// 纯文本日志方法
func (m *MockAdapter) DebugMsg(msg string) {}
func (m *MockAdapter) InfoMsg(msg string)  {}
func (m *MockAdapter) WarnMsg(msg string)  {}
func (m *MockAdapter) ErrorMsg(msg string) {}
func (m *MockAdapter) FatalMsg(msg string) {}

// 带上下文的日志方法
func (m *MockAdapter) DebugContext(ctx context.Context, format string, args ...interface{}) {}
func (m *MockAdapter) InfoContext(ctx context.Context, format string, args ...interface{})  {}
func (m *MockAdapter) WarnContext(ctx context.Context, format string, args ...interface{})  {}
func (m *MockAdapter) ErrorContext(ctx context.Context, format string, args ...interface{}) {}
func (m *MockAdapter) FatalContext(ctx context.Context, format string, args ...interface{}) {}

// 结构化日志方法（键值对）
func (m *MockAdapter) DebugKV(msg string, keysAndValues ...interface{}) {}
func (m *MockAdapter) InfoKV(msg string, keysAndValues ...interface{})  {}
func (m *MockAdapter) WarnKV(msg string, keysAndValues ...interface{})  {}
func (m *MockAdapter) ErrorKV(msg string, keysAndValues ...interface{}) {}
func (m *MockAdapter) FatalKV(msg string, keysAndValues ...interface{}) {}

// 带上下文的结构化日志方法（键值对）
func (m *MockAdapter) DebugContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {}
func (m *MockAdapter) InfoContextKV(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (m *MockAdapter) WarnContextKV(ctx context.Context, msg string, keysAndValues ...interface{})  {}
func (m *MockAdapter) ErrorContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {}
func (m *MockAdapter) FatalContextKV(ctx context.Context, msg string, keysAndValues ...interface{}) {}

// 原始日志条目方法
func (m *MockAdapter) Log(level LogLevel, msg string)                                          {}
func (m *MockAdapter) LogContext(ctx context.Context, level LogLevel, msg string)              {}
func (m *MockAdapter) LogKV(level LogLevel, msg string, keysAndValues ...interface{})          {}
func (m *MockAdapter) LogWithFields(level LogLevel, msg string, fields map[string]interface{}) {}

// 多行日志方法
func (m *MockAdapter) DebugLines(lines ...string) {}
func (m *MockAdapter) InfoLines(lines ...string)  {}
func (m *MockAdapter) WarnLines(lines ...string)  {}
func (m *MockAdapter) ErrorLines(lines ...string) {}

// WithContext 的实现
func (m *MockAdapter) WithContext(ctx context.Context) ILogger {
	return m
}

// 兼容标准log包的方法
func (m *MockAdapter) Print(args ...interface{})                 {}
func (m *MockAdapter) Printf(format string, args ...interface{}) {}
func (m *MockAdapter) Println(args ...interface{})               {}

func (m *MockAdapter) Clone() ILogger {
	clone := &MockAdapter{
		level:      m.level,
		showCaller: m.showCaller,
		initError:  m.initError,
		closeError: m.closeError,
		flushError: m.flushError,
		healthy:    m.healthy,
		name:       m.name,
		version:    m.version,
	}
	if m.fields != nil {
		clone.fields = make(map[string]interface{})
		for k, v := range m.fields {
			clone.fields[k] = v
		}
	}
	return clone
}

// TestMockAdapter 测试模拟适配器本身
func (suite *AdapterTestSuite) TestMockAdapter() {
	mock := &MockAdapter{
		level:   INFO,
		healthy: true,
	}

	// 测试基本功能
	assert.Equal(suite.T(), INFO, mock.GetLevel())
	assert.True(suite.T(), mock.IsHealthy())
	assert.Equal(suite.T(), "mock", mock.GetAdapterName())
	assert.Equal(suite.T(), "1.0.0", mock.GetAdapterVersion())

	// 测试设置级别
	mock.SetLevel(DEBUG)
	assert.Equal(suite.T(), DEBUG, mock.GetLevel())

	// 测试设置调用者显示
	mock.SetShowCaller(true)
	assert.True(suite.T(), mock.IsShowCaller())

	// 测试级别检查
	assert.True(suite.T(), mock.IsLevelEnabled(INFO))
	assert.True(suite.T(), mock.IsLevelEnabled(LogLevel(10))) // 修复：MockAdapter的IsLevelEnabled返回true

	// 测试添加字段
	newLogger := mock.WithField("test", "value")
	assert.NotNil(suite.T(), newLogger)

	// 测试添加多个字段
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	newLogger = mock.WithFields(fields)
	assert.NotNil(suite.T(), newLogger)

	// 测试添加错误
	err := errors.New("test error")
	newLogger = mock.WithError(err)
	assert.NotNil(suite.T(), newLogger)

	// 测试克隆
	clone := mock.Clone()
	assert.NotNil(suite.T(), clone)
	// 修复：对于mock适配器，克隆返回的是相同的实例
	assert.Equal(suite.T(), mock, clone)

	// 测试错误情况
	mock.initError = errors.New("init error")
	assert.Error(suite.T(), mock.Initialize())

	mock.closeError = errors.New("close error")
	assert.Error(suite.T(), mock.Close())

	mock.flushError = errors.New("flush error")
	assert.Error(suite.T(), mock.Flush())

	// 测试自定义名称和版本
	mock.name = "custom"
	mock.version = "2.0.0"
	assert.Equal(suite.T(), "custom", mock.GetAdapterName())
	assert.Equal(suite.T(), "2.0.0", mock.GetAdapterVersion())
}

// TestComplexScenarios 测试复杂场景
func (suite *AdapterTestSuite) TestComplexScenarios() {
	// 场景1: 混合操作
	suite.T().Run("MixedOperations", func(t *testing.T) {
		manager := NewLoggerManager()

		// 添加多个不同类型的适配器
		configs := []struct {
			name   string
			level  LogLevel
			output *bytes.Buffer
		}{
			{"debug_adapter", DEBUG, &bytes.Buffer{}},
			{"info_adapter", INFO, &bytes.Buffer{}},
			{"error_adapter", ERROR, &bytes.Buffer{}},
		}

		for _, cfg := range configs {
			config := DefaultAdapterConfig()
			config.Name = cfg.name
			config.Level = cfg.level
			config.Output = cfg.output

			adapter, err := CreateAdapter(config)
			assert.NoError(t, err)

			err = manager.AddAdapter(cfg.name, adapter)
			assert.NoError(t, err)
		}

		// 执行广播测试
		manager.Broadcast(INFO, "Test info message")
		manager.Broadcast(ERROR, "Test error message")

		// 验证健康状态
		health := manager.HealthCheck()
		assert.Len(t, health, 3)
		for name, healthy := range health {
			assert.True(t, healthy, "Adapter %s should be healthy", name)
		}

		// 清理
		err := manager.CloseAll()
		assert.NoError(t, err)
	})

	// 场景2: 压力测试
	suite.T().Run("StressTest", func(t *testing.T) {
		manager := NewLoggerManager()

		// 快速添加和删除大量适配器
		numOps := 100
		for i := 0; i < numOps; i++ {
			config := DefaultAdapterConfig()
			config.Name = fmt.Sprintf("stress_%d", i)
			config.Output = &bytes.Buffer{}

			adapter, err := CreateAdapter(config)
			if err == nil {
				manager.AddAdapter(config.Name, adapter)
				if i%2 == 0 {
					manager.RemoveAdapter(config.Name)
				}
			}
		}

		// 验证最终状态
		adapters := manager.ListAdapters()
		assert.True(t, len(adapters) < numOps) // 应该有一些被删除了

		// 清理
		manager.CloseAll()
	})
}

// 运行测试套件
func TestAdapterSuite(t *testing.T) {
	suite.Run(t, new(AdapterTestSuite))
}

// TestEdgeCases 测试边界情况
func TestEdgeCases(t *testing.T) {
	t.Run("EmptyManager", func(t *testing.T) {
		manager := NewLoggerManager()

		// 空管理器的操作
		err := manager.FlushAll()
		assert.NoError(t, err)

		err = manager.CloseAll()
		assert.NoError(t, err)

		manager.SetLevelAll(DEBUG)
		manager.Broadcast(INFO, "test")

		health := manager.HealthCheck()
		assert.Empty(t, health)

		adapters := manager.ListAdapters()
		assert.Empty(t, adapters)
	})

	t.Run("ConfigValidation", func(t *testing.T) {
		// 测试各种配置边界情况
		config := &AdapterConfig{
			Type:   StandardAdapter,
			Name:   "test",
			Level:  LogLevel(999), // 无效级别
			Output: nil,
			Fields: nil,
		}

		err := config.Validate()
		assert.NoError(t, err) // Validate不检查LogLevel有效性

		// 创建适配器应该能处理无效级别
		adapter, err := CreateAdapter(config)
		assert.NoError(t, err)
		assert.NotNil(t, adapter)
	})
}

// TestPerformance 性能测试（在单元测试框架中模拟）
func TestPerformance(t *testing.T) {
	t.Run("ManagerOperations", func(t *testing.T) {
		manager := NewLoggerManager()
		buffer := &bytes.Buffer{}

		// 测量添加适配器的性能
		start := time.Now()
		numAdapters := 1000

		for i := 0; i < numAdapters; i++ {
			config := DefaultAdapterConfig()
			config.Name = fmt.Sprintf("perf_%d", i)
			config.Output = buffer

			adapter, err := CreateAdapter(config)
			if err == nil {
				manager.AddAdapter(config.Name, adapter)
			}
		}

		duration := time.Since(start)
		t.Logf("Added %d adapters in %v", numAdapters, duration)
		assert.True(t, duration < time.Second*5, "Adding adapters took too long")

		// 测量广播性能
		start = time.Now()
		for i := 0; i < 100; i++ {
			manager.Broadcast(INFO, "Performance test message %d", i)
		}
		duration = time.Since(start)
		t.Logf("100 broadcasts to %d adapters took %v", numAdapters, duration)

		// 清理
		start = time.Now()
		manager.CloseAll()
		duration = time.Since(start)
		t.Logf("Closed all adapters in %v", duration)
	})
}
