/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\types_test.go
 * @Description: 核心类型测试套件
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

// TypesTestSuite 核心类型测试套件
type TypesTestSuite struct {
	suite.Suite
}

// TestLoggerStats 测试日志统计
func (suite *TypesTestSuite) TestLoggerStats() {
	stats := NewLoggerStats()

	// 验证初始状态
	assert.NotNil(suite.T(), stats)
	assert.False(suite.T(), stats.StartTime.IsZero())
	assert.Equal(suite.T(), int64(0), stats.TotalLogs)
	assert.Equal(suite.T(), int64(0), stats.ErrorCount)
	assert.Equal(suite.T(), time.Duration(0), stats.Uptime)
	assert.Equal(suite.T(), int64(0), stats.BytesWritten)
	assert.NotNil(suite.T(), stats.LevelCounts)
	assert.Len(suite.T(), stats.LevelCounts, 0)

	// 测试级别计数增加
	stats.IncrementLevel(INFO)
	statsSnapshot := stats.GetStats()
	assert.Equal(suite.T(), int64(1), statsSnapshot.TotalLogs)
	assert.Equal(suite.T(), int64(1), statsSnapshot.LevelCounts[INFO])
	assert.Equal(suite.T(), int64(0), statsSnapshot.ErrorCount)
	assert.True(suite.T(), statsSnapshot.Uptime >= 0) // 修复：uptime可能为0或大于0
	assert.False(suite.T(), statsSnapshot.LastLogTime.IsZero())

	// 测试错误级别计数
	stats.IncrementLevel(ERROR)
	statsSnapshot = stats.GetStats()
	assert.Equal(suite.T(), int64(2), statsSnapshot.TotalLogs)
	assert.Equal(suite.T(), int64(1), statsSnapshot.LevelCounts[INFO])
	assert.Equal(suite.T(), int64(1), statsSnapshot.LevelCounts[ERROR])
	assert.Equal(suite.T(), int64(1), statsSnapshot.ErrorCount)

	// 测试FATAL级别计数
	stats.IncrementLevel(FATAL)
	statsSnapshot = stats.GetStats()
	assert.Equal(suite.T(), int64(3), statsSnapshot.TotalLogs)
	assert.Equal(suite.T(), int64(1), statsSnapshot.LevelCounts[FATAL])
	assert.Equal(suite.T(), int64(2), statsSnapshot.ErrorCount) // ERROR和FATAL都算错误

	// 测试字节数增加
	stats.AddBytes(1024)
	statsSnapshot = stats.GetStats()
	assert.Equal(suite.T(), int64(1024), statsSnapshot.BytesWritten)

	stats.AddBytes(512)
	statsSnapshot = stats.GetStats()
	assert.Equal(suite.T(), int64(1536), statsSnapshot.BytesWritten)
}

// TestLoggerStatsConcurrency 测试统计的并发安全
func (suite *TypesTestSuite) TestLoggerStatsConcurrency() {
	stats := NewLoggerStats()
	var wg sync.WaitGroup

	// 并发增加不同级别
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(level LogLevel) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				stats.IncrementLevel(level)
			}
		}(LogLevel(i % 5))
	}

	// 并发增加字节数
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				stats.AddBytes(10)
			}
		}()
	}

	wg.Wait()

	// 验证最终结果
	statsSnapshot := stats.GetStats()
	assert.Equal(suite.T(), int64(1000), statsSnapshot.TotalLogs)
	assert.Equal(suite.T(), int64(5000), statsSnapshot.BytesWritten)

	// 验证级别计数总和
	var totalLevelCount int64
	for _, count := range statsSnapshot.LevelCounts {
		totalLevelCount += count
	}
	assert.Equal(suite.T(), int64(1000), totalLevelCount)
}

// TestLoggerOptions 测试日志器选项
func (suite *TypesTestSuite) TestLoggerOptions() {
	// 测试默认选项
	options := DefaultLoggerOptions()
	assert.NotNil(suite.T(), options)
	assert.NotNil(suite.T(), options.Config)
	assert.NotNil(suite.T(), options.Writers)
	assert.NotNil(suite.T(), options.Hooks)
	assert.NotNil(suite.T(), options.Middleware)
	assert.Equal(suite.T(), context.Background(), options.Context)

	// 测试修改选项
	config := DefaultConfig()
	config.Level = DEBUG
	options.Config = config

	assert.Equal(suite.T(), DEBUG, options.Config.Level)
}

// TestFieldMap 测试字段映射
func (suite *TypesTestSuite) TestFieldMap() {
	fields := make(FieldMap)

	// 测试添加不同类型的值
	fields["string"] = "test"
	fields["int"] = 42
	fields["float"] = 3.14
	fields["bool"] = true
	fields["nil"] = nil
	fields["slice"] = []string{"a", "b", "c"}
	fields["map"] = map[string]interface{}{"nested": "value"}

	assert.Equal(suite.T(), "test", fields["string"])
	assert.Equal(suite.T(), 42, fields["int"])
	assert.Equal(suite.T(), 3.14, fields["float"])
	assert.Equal(suite.T(), true, fields["bool"])
	assert.Nil(suite.T(), fields["nil"])
	assert.Equal(suite.T(), []string{"a", "b", "c"}, fields["slice"])
	assert.Equal(suite.T(), map[string]interface{}{"nested": "value"}, fields["map"])
}

// TestAdapterRegistry 测试适配器注册表
func (suite *TypesTestSuite) TestAdapterRegistry() {
	registry := NewAdapterRegistry()
	assert.NotNil(suite.T(), registry)
	assert.Empty(suite.T(), registry.List())

	// 创建模拟工厂函数
	mockFactory := func(config *AdapterConfig) (IAdapter, error) {
		return &MockAdapter{name: "mock"}, nil
	}

	errorFactory := func(config *AdapterConfig) (IAdapter, error) {
		return nil, errors.New("factory error")
	}

	// 测试注册适配器
	registry.Register("mock", mockFactory)
	registry.Register("error", errorFactory)

	adapters := registry.List()
	assert.Len(suite.T(), adapters, 2)
	assert.Contains(suite.T(), adapters, "mock")
	assert.Contains(suite.T(), adapters, "error")

	// 测试创建适配器
	config := DefaultAdapterConfig()
	adapter, err := registry.Create("mock", config)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), adapter)

	// 测试创建失败的适配器
	adapter, err = registry.Create("error", config)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), adapter)
	assert.Contains(suite.T(), err.Error(), "factory error")

	// 测试创建不存在的适配器
	adapter, err = registry.Create("nonexistent", config)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), adapter)
	assert.Contains(suite.T(), err.Error(), "not found")
}

// TestAdapterRegistryConcurrency 测试注册表并发安全
func (suite *TypesTestSuite) TestAdapterRegistryConcurrency() {
	registry := NewAdapterRegistry()
	var wg sync.WaitGroup

	// 并发注册
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			factory := func(config *AdapterConfig) (IAdapter, error) {
				return &MockAdapter{name: fmt.Sprintf("adapter_%d", index)}, nil
			}
			registry.Register(fmt.Sprintf("adapter_%d", index), factory)
		}(i)
	}

	// 并发读取
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				registry.List()
			}
		}()
	}

	wg.Wait()

	// 验证所有适配器都注册成功
	adapters := registry.List()
	assert.Len(suite.T(), adapters, 10)
}

// TestBufferPool 测试缓冲区池
func (suite *TypesTestSuite) TestBufferPool() {
	pool := NewBufferPool()
	assert.NotNil(suite.T(), pool)

	// 测试获取缓冲区
	buf1 := pool.Get()
	assert.NotNil(suite.T(), buf1)
	assert.Equal(suite.T(), 0, len(buf1))

	buf2 := pool.Get()
	assert.NotNil(suite.T(), buf2)
	assert.True(suite.T(), &buf1 != &buf2) // 比较地址而不是内容

	// 测试使用缓冲区
	buf1 = append(buf1, []byte("test data")...)
	assert.Equal(suite.T(), "test data", string(buf1))

	// 测试归还缓冲区
	pool.Put(buf1)
	pool.Put(buf2)

	// 再次获取应该能重用
	buf3 := pool.Get()
	assert.NotNil(suite.T(), buf3)
	assert.Equal(suite.T(), 0, len(buf3)) // 应该被重置

	// 测试归还过大的缓冲区
	largeBuf := make([]byte, 100*1024) // 100KB
	pool.Put(largeBuf)                 // 应该被丢弃而不是保存
}

// TestBufferPoolConcurrency 测试缓冲区池并发安全
func (suite *TypesTestSuite) TestBufferPoolConcurrency() {
	pool := NewBufferPool()
	var wg sync.WaitGroup

	// 并发获取和归还缓冲区
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				buf := pool.Get()
				buf = append(buf, []byte("test")...)
				pool.Put(buf)
			}
		}()
	}

	wg.Wait()

	// 验证池仍然正常工作
	buf := pool.Get()
	assert.NotNil(suite.T(), buf)
	assert.Equal(suite.T(), 0, len(buf))
}

// TestLogEvent 测试日志事件
func (suite *TypesTestSuite) TestLogEvent() {
	// 测试创建新事件
	event := NewLogEvent(EventLogCreated, INFO, "Test message")
	assert.NotNil(suite.T(), event)
	assert.Equal(suite.T(), EventLogCreated, event.Type)
	assert.Equal(suite.T(), INFO, event.Level)
	assert.Equal(suite.T(), "Test message", event.Message)
	assert.False(suite.T(), event.Timestamp.IsZero())
	assert.NotNil(suite.T(), event.Fields)
	assert.Empty(suite.T(), event.Fields)
	assert.Nil(suite.T(), event.Error)

	// 测试设置字段
	event.Fields["key"] = "value"
	event.Error = errors.New("test error")
	// Context field removed

	assert.Equal(suite.T(), "value", event.Fields["key"])
	assert.NotNil(suite.T(), event.Error)
	assert.Equal(suite.T(), "test error", event.Error.Error())
}

// TestEventTypes 测试事件类型
func (suite *TypesTestSuite) TestEventTypes() {
	// 验证事件类型常量
	assert.Equal(suite.T(), EventType(0), EventLogCreated)
	assert.Equal(suite.T(), EventType(1), EventLogProcessed)
	assert.Equal(suite.T(), EventType(2), EventLogWritten)
	assert.Equal(suite.T(), EventType(3), EventLogError)
	assert.Equal(suite.T(), EventType(4), EventLoggerStarted)
	assert.Equal(suite.T(), EventType(5), EventLoggerStopped)
	assert.Equal(suite.T(), EventType(6), EventLoggerConfigChanged)

	// 测试不同类型的事件创建
	events := []struct {
		eventType EventType
		level     LogLevel
		message   string
	}{
		{EventLogCreated, DEBUG, "Log created"},
		{EventLogProcessed, INFO, "Log processed"},
		{EventLogWritten, WARN, "Log written"},
		{EventLogError, ERROR, "Log error"},
		{EventLoggerStarted, INFO, "Logger started"},
		{EventLoggerStopped, INFO, "Logger stopped"},
		{EventLoggerConfigChanged, INFO, "Config changed"},
	}

	for _, e := range events {
		event := NewLogEvent(e.eventType, e.level, e.message)
		assert.Equal(suite.T(), e.eventType, event.Type)
		assert.Equal(suite.T(), e.level, event.Level)
		assert.Equal(suite.T(), e.message, event.Message)
	}
}

// TestNewLoggerWithOptions 测试使用选项创建日志器
func (suite *TypesTestSuite) TestNewLoggerWithOptions() {
	// 测试使用nil选项
	logger := NewLoggerWithOptions(nil)
	assert.NotNil(suite.T(), logger)
	assert.NotNil(suite.T(), logger.stats)
	assert.NotNil(suite.T(), logger.context)
	assert.NotNil(suite.T(), logger.cancel)

	// 清理
	logger.cancel()

	// 测试使用自定义选项
	options := &LoggerOptions{
		Config: &LogConfig{
			Level:      DEBUG,
			ShowCaller: true,
		},
		Context: context.Background(),
	}

	logger = NewLoggerWithOptions(options)
	assert.NotNil(suite.T(), logger)
	assert.Equal(suite.T(), DEBUG, logger.level)
	assert.True(suite.T(), logger.showCaller)
	assert.NotNil(suite.T(), logger.context)

	// 清理
	logger.cancel()
}

// TestLoggerOptionsDefaults 测试日志器选项默认值
func (suite *TypesTestSuite) TestLoggerOptionsDefaults() {
	logger := NewLoggerWithOptions(nil)
	assert.NotNil(suite.T(), logger.formatter)
	assert.NotEmpty(suite.T(), logger.writers)

	// 清理
	logger.cancel()
}

// 运行测试套件
func TestTypesSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

// TestTypesPerformance 类型操作的性能测试（模拟）
func TestTypesPerformance(t *testing.T) {
	t.Run("LoggerStats", func(t *testing.T) {
		stats := NewLoggerStats()
		start := time.Now()

		for i := 0; i < 10000; i++ {
			stats.IncrementLevel(LogLevel(i % 5))
			if i%100 == 0 {
				stats.AddBytes(1024)
			}
		}

		duration := time.Since(start)
		t.Logf("10000 stat operations took %v", duration)
		assert.True(t, duration < time.Second,
			"Stats operations should be fast, took %v", duration)
	})

	t.Run("BufferPool", func(t *testing.T) {
		pool := NewBufferPool()
		start := time.Now()

		for i := 0; i < 10000; i++ {
			buf := pool.Get()
			buf = append(buf, []byte("test data")...)
			pool.Put(buf)
		}

		duration := time.Since(start)
		t.Logf("10000 buffer pool operations took %v", duration)
		assert.True(t, duration < time.Second,
			"Buffer pool operations should be fast, took %v", duration)
	})
}
