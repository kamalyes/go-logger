/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\config_test.go
 * @Description: 配置模块测试套件
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite 配置测试套件
type ConfigTestSuite struct {
	suite.Suite
	buffer *bytes.Buffer
}

// SetupTest 测试前准备
func (suite *ConfigTestSuite) SetupTest() {
	suite.buffer = &bytes.Buffer{}
}

// TearDownTest 测试后清理
func (suite *ConfigTestSuite) TearDownTest() {
	suite.buffer = nil
}

// TestDefaultConfig 测试默认配置
func (suite *ConfigTestSuite) TestDefaultConfig() {
	config := DefaultConfig()
	
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), INFO, config.Level)
	assert.False(suite.T(), config.ShowCaller)
	assert.Equal(suite.T(), "", config.Prefix)
	assert.Equal(suite.T(), os.Stdout, config.Output)
	assert.True(suite.T(), config.Colorful)
	assert.Equal(suite.T(), "2006-01-02 15:04:05", config.TimeFormat)
}

// TestNewConfig 测试新配置创建
func (suite *ConfigTestSuite) TestNewConfig() {
	config := NewConfig()
	
	assert.NotNil(suite.T(), config)
	assert.Equal(suite.T(), INFO, config.Level)
	assert.False(suite.T(), config.ShowCaller)
	assert.Equal(suite.T(), "", config.Prefix)
	assert.Equal(suite.T(), os.Stdout, config.Output)
	assert.True(suite.T(), config.Colorful)
	assert.Equal(suite.T(), "2006-01-02 15:04:05", config.TimeFormat)
}

// TestWithLevel 测试设置日志级别
func (suite *ConfigTestSuite) TestWithLevel() {
	config := NewConfig()
	
	// 测试设置不同级别
	result := config.WithLevel(DEBUG)
	assert.Equal(suite.T(), config, result) // 应该返回同一个实例
	assert.Equal(suite.T(), DEBUG, config.Level)
	
	config.WithLevel(WARN)
	assert.Equal(suite.T(), WARN, config.Level)
	
	config.WithLevel(ERROR)
	assert.Equal(suite.T(), ERROR, config.Level)
	
	config.WithLevel(FATAL)
	assert.Equal(suite.T(), FATAL, config.Level)
}

// TestWithShowCaller 测试设置显示调用者
func (suite *ConfigTestSuite) TestWithShowCaller() {
	config := NewConfig()
	
	// 默认为false
	assert.False(suite.T(), config.ShowCaller)
	
	// 设置为true
	result := config.WithShowCaller(true)
	assert.Equal(suite.T(), config, result)
	assert.True(suite.T(), config.ShowCaller)
	
	// 设置为false
	config.WithShowCaller(false)
	assert.False(suite.T(), config.ShowCaller)
}

// TestWithPrefix 测试设置前缀
func (suite *ConfigTestSuite) TestWithPrefix() {
	config := NewConfig()
	
	// 测试空前缀
	result := config.WithPrefix("")
	assert.Equal(suite.T(), config, result)
	assert.Equal(suite.T(), "", config.Prefix)
	
	// 测试不以空格结尾的前缀
	config.WithPrefix("TEST")
	assert.Equal(suite.T(), "TEST ", config.Prefix)
	
	// 测试已经以空格结尾的前缀
	config.WithPrefix("DEBUG ")
	assert.Equal(suite.T(), "DEBUG ", config.Prefix)
	
	// 测试多个单词的前缀
	config.WithPrefix("SERVICE API")
	assert.Equal(suite.T(), "SERVICE API ", config.Prefix)
}

// TestWithOutput 测试设置输出目标
func (suite *ConfigTestSuite) TestWithOutput() {
	config := NewConfig()
	
	// 默认输出
	assert.Equal(suite.T(), os.Stdout, config.Output)
	
	// 设置为缓冲区
	result := config.WithOutput(suite.buffer)
	assert.Equal(suite.T(), config, result)
	assert.Equal(suite.T(), suite.buffer, config.Output)
	
	// 设置为标准错误
	config.WithOutput(os.Stderr)
	assert.Equal(suite.T(), os.Stderr, config.Output)
}

// TestWithColorful 测试设置彩色输出
func (suite *ConfigTestSuite) TestWithColorful() {
	config := NewConfig()
	
	// 默认为true
	assert.True(suite.T(), config.Colorful)
	
	// 设置为false
	result := config.WithColorful(false)
	assert.Equal(suite.T(), config, result)
	assert.False(suite.T(), config.Colorful)
	
	// 设置为true
	config.WithColorful(true)
	assert.True(suite.T(), config.Colorful)
}

// TestWithTimeFormat 测试设置时间格式
func (suite *ConfigTestSuite) TestWithTimeFormat() {
	config := NewConfig()
	
	// 默认格式
	assert.Equal(suite.T(), "2006-01-02 15:04:05", config.TimeFormat)
	
	// 设置自定义格式
	customFormat := "15:04:05.000"
	result := config.WithTimeFormat(customFormat)
	assert.Equal(suite.T(), config, result)
	assert.Equal(suite.T(), customFormat, config.TimeFormat)
	
	// 设置ISO格式
	isoFormat := "2006-01-02T15:04:05Z07:00"
	config.WithTimeFormat(isoFormat)
	assert.Equal(suite.T(), isoFormat, config.TimeFormat)
}

// TestClone 测试配置克隆
func (suite *ConfigTestSuite) TestClone() {
	original := NewConfig()
	original.WithLevel(DEBUG)
	original.WithShowCaller(true)
	original.WithPrefix("TEST")
	original.WithOutput(suite.buffer)
	original.WithColorful(false)
	original.WithTimeFormat("15:04:05")
	
	// 克隆配置
	cloned := original.Clone()
	
	// 验证克隆的配置与原配置相同
	assert.NotSame(suite.T(), original, cloned) // 不同的实例但内容相同
	assert.Equal(suite.T(), original.Level, cloned.Level)
	assert.Equal(suite.T(), original.ShowCaller, cloned.ShowCaller)
	assert.Equal(suite.T(), original.Prefix, cloned.Prefix)
	assert.Equal(suite.T(), original.Output, cloned.Output)
	assert.Equal(suite.T(), original.Colorful, cloned.Colorful)
	assert.Equal(suite.T(), original.TimeFormat, cloned.TimeFormat)
	
	// 修改克隆的配置不应影响原配置
	cloned.WithLevel(ERROR)
	assert.Equal(suite.T(), DEBUG, original.Level)
	assert.Equal(suite.T(), ERROR, cloned.Level)
}

// TestValidate 测试配置验证
func (suite *ConfigTestSuite) TestValidate() {
	config := NewConfig()
	
	// 有效配置
	err := config.Validate()
	assert.NoError(suite.T(), err)
	
	// 测试nil输出
	config.Output = nil
	err = config.Validate()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), os.Stdout, config.Output) // 应该设置为默认值
	
	// 测试空时间格式
	config.TimeFormat = ""
	err = config.Validate()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "2006-01-02 15:04:05", config.TimeFormat) // 应该设置为默认值
}

// TestConfigChaining 测试配置链式调用
func (suite *ConfigTestSuite) TestConfigChaining() {
	config := NewConfig().
		WithLevel(WARN).
		WithShowCaller(true).
		WithPrefix("CHAIN").
		WithOutput(suite.buffer).
		WithColorful(false).
		WithTimeFormat("15:04:05.000")
	
	assert.Equal(suite.T(), WARN, config.Level)
	assert.True(suite.T(), config.ShowCaller)
	assert.Equal(suite.T(), "CHAIN ", config.Prefix)
	assert.Equal(suite.T(), suite.buffer, config.Output)
	assert.False(suite.T(), config.Colorful)
	assert.Equal(suite.T(), "15:04:05.000", config.TimeFormat)
}

// TestConfigEquality 测试配置相等性
func (suite *ConfigTestSuite) TestConfigEquality() {
	config1 := NewConfig().WithLevel(DEBUG).WithShowCaller(true)
	config2 := NewConfig().WithLevel(DEBUG).WithShowCaller(true)
	
	// 同样的设置应该产生相同的值
	assert.Equal(suite.T(), config1.Level, config2.Level)
	assert.Equal(suite.T(), config1.ShowCaller, config2.ShowCaller)
	assert.Equal(suite.T(), config1.Prefix, config2.Prefix)
	assert.Equal(suite.T(), config1.Output, config2.Output)
	assert.Equal(suite.T(), config1.Colorful, config2.Colorful)
	assert.Equal(suite.T(), config1.TimeFormat, config2.TimeFormat)
}

// TestConfigImmutability 测试配置不可变性
func (suite *ConfigTestSuite) TestConfigImmutability() {
	original := NewConfig()
	originalLevel := original.Level
	originalShowCaller := original.ShowCaller
	originalPrefix := original.Prefix
	
	// 创建新配置应该不影响原配置
	modified := original.Clone()
	modified.WithLevel(ERROR)
	modified.WithShowCaller(true)
	modified.WithPrefix("MODIFIED")
	
	assert.Equal(suite.T(), originalLevel, original.Level)
	assert.Equal(suite.T(), originalShowCaller, original.ShowCaller)
	assert.Equal(suite.T(), originalPrefix, original.Prefix)
}

// TestConfigWithDifferentLevels 测试不同级别的配置
func (suite *ConfigTestSuite) TestConfigWithDifferentLevels() {
	levels := []LogLevel{DEBUG, INFO, WARN, ERROR, FATAL}
	levelNames := []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
	
	for i, level := range levels {
		suite.T().Run(levelNames[i], func(t *testing.T) {
			config := NewConfig().WithLevel(level)
			assert.Equal(t, level, config.Level)
		})
	}
}

// TestConfigWithDifferentOutputs 测试不同输出目标的配置
func (suite *ConfigTestSuite) TestConfigWithDifferentOutputs() {
	outputs := map[string]interface{}{
		"stdout": os.Stdout,
		"stderr": os.Stderr,
		"buffer": suite.buffer,
	}
	
	for name, output := range outputs {
		suite.T().Run(name, func(t *testing.T) {
			config := NewConfig().WithOutput(output.(interface{ Write([]byte) (int, error) }))
			assert.Equal(t, output, config.Output)
		})
	}
}

// TestConfigTimeFormats 测试不同时间格式
func (suite *ConfigTestSuite) TestConfigTimeFormats() {
	timeFormats := []string{
		"15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"Jan 2 15:04:05",
		"2006/01/02 15:04:05.000",
	}
	
	for _, format := range timeFormats {
		suite.T().Run(strings.ReplaceAll(format, ":", "_"), func(t *testing.T) {
			config := NewConfig().WithTimeFormat(format)
			assert.Equal(t, format, config.TimeFormat)
		})
	}
}

// TestConfigPrefixes 测试不同前缀格式
func (suite *ConfigTestSuite) TestConfigPrefixes() {
	prefixes := map[string]string{
		"simple":      "LOG",
		"bracketed":   "[APP]",
		"with_space":  "DEBUG ",
		"complex":     "[2023-11-08] [SERVICE]",
		"empty":       "",
	}
	
	for name, prefix := range prefixes {
		suite.T().Run(name, func(t *testing.T) {
			config := NewConfig().WithPrefix(prefix)
			
			if prefix == "" {
				assert.Equal(t, "", config.Prefix)
			} else if strings.HasSuffix(prefix, " ") {
				assert.Equal(t, prefix, config.Prefix)
			} else {
				assert.Equal(t, prefix+" ", config.Prefix)
			}
		})
	}
}

// TestConfigValidationEdgeCases 测试配置验证边界情况
func (suite *ConfigTestSuite) TestConfigValidationEdgeCases() {
	// 测试多次验证
	config := NewConfig()
	
	err := config.Validate()
	assert.NoError(suite.T(), err)
	
	err = config.Validate()
	assert.NoError(suite.T(), err)
	
	// 测试验证后的状态
	assert.NotNil(suite.T(), config.Output)
	assert.NotEmpty(suite.T(), config.TimeFormat)
}

// TestConfigConcurrentAccess 测试并发访问配置
func (suite *ConfigTestSuite) TestConfigConcurrentAccess() {
	config := NewConfig()
	done := make(chan bool, 2)
	
	// 并发读取配置
	go func() {
		for i := 0; i < 100; i++ {
			_ = config.Level
			_ = config.ShowCaller
			_ = config.Prefix
			_ = config.Output
			_ = config.Colorful
			_ = config.TimeFormat
		}
		done <- true
	}()
	
	// 并发修改配置（通过方法链）
	go func() {
		for i := 0; i < 100; i++ {
			tempConfig := config.Clone()
			tempConfig.WithLevel(DEBUG)
			tempConfig.WithShowCaller(true)
			tempConfig.WithPrefix("TEMP")
		}
		done <- true
	}()
	
	// 等待完成
	<-done
	<-done
	
	// 验证原配置未受影响
	assert.Equal(suite.T(), INFO, config.Level) // 默认级别
}

// TestConfigMemoryUsage 测试配置内存使用
func (suite *ConfigTestSuite) TestConfigMemoryUsage() {
	configs := make([]*LogConfig, 1000)
	
	for i := 0; i < 1000; i++ {
		configs[i] = NewConfig().
			WithLevel(LogLevel(i % 5)).
			WithShowCaller(i%2 == 0).
			WithPrefix("CONFIG_" + string(rune('A'+i%26))).
			WithColorful(i%3 == 0)
	}
	
	// 验证配置创建成功
	assert.Len(suite.T(), configs, 1000)
	
	// 验证每个配置都是独立的
	for i := 0; i < 10; i++ {
		config := configs[i]
		assert.NotNil(suite.T(), config)
		assert.Equal(suite.T(), LogLevel(i%5), config.Level)
		assert.Equal(suite.T(), i%2 == 0, config.ShowCaller)
		expectedPrefix := "CONFIG_" + string(rune('A'+i%26)) + " "
		assert.Equal(suite.T(), expectedPrefix, config.Prefix)
	}
}

// TestConfigSerialization 测试配置序列化兼容性
func (suite *ConfigTestSuite) TestConfigSerialization() {
	// 创建一个复杂的配置
	original := NewConfig().
		WithLevel(DEBUG).
		WithShowCaller(true).
		WithPrefix("[TEST]").
		WithColorful(false).
		WithTimeFormat("2006-01-02T15:04:05.000Z")
	
	// 通过克隆模拟序列化/反序列化
	cloned := original.Clone()
	
	// 验证克隆后的配置
	assert.Equal(suite.T(), original.Level, cloned.Level)
	assert.Equal(suite.T(), original.ShowCaller, cloned.ShowCaller)
	assert.Equal(suite.T(), original.Prefix, cloned.Prefix)
	assert.Equal(suite.T(), original.Colorful, cloned.Colorful)
	assert.Equal(suite.T(), original.TimeFormat, cloned.TimeFormat)
}

// 运行测试套件
func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// TestLogConfigDefaults 测试配置默认值
func TestLogConfigDefaults(t *testing.T) {
	config := &LogConfig{}
	
	// 验证零值状态
	assert.Equal(t, LogLevel(0), config.Level)
	assert.False(t, config.ShowCaller)
	assert.Equal(t, "", config.Prefix)
	assert.Nil(t, config.Output)
	assert.False(t, config.Colorful)
	assert.Equal(t, "", config.TimeFormat)
}

// TestLogConfigValidation 测试独立的配置验证功能
func TestLogConfigValidation(t *testing.T) {
	// 测试完全空的配置
	config := &LogConfig{}
	err := config.Validate()
	assert.NoError(t, err)
	assert.Equal(t, os.Stdout, config.Output)
	assert.Equal(t, "2006-01-02 15:04:05", config.TimeFormat)
	
	// 测试部分填充的配置
	config = &LogConfig{
		Level:      DEBUG,
		ShowCaller: true,
		Prefix:     "[TEST] ",
	}
	err = config.Validate()
	assert.NoError(t, err)
	assert.Equal(t, os.Stdout, config.Output)
	assert.Equal(t, "2006-01-02 15:04:05", config.TimeFormat)
}

// BenchmarkConfigCreation 配置创建的基准测试（模拟）
func TestBenchmarkConfigCreation(t *testing.T) {
	start := time.Now()
	iterations := 10000
	
	for i := 0; i < iterations; i++ {
		config := NewConfig()
		config.WithLevel(DEBUG).
			WithShowCaller(true).
			WithPrefix("BENCH").
			WithColorful(false)
	}
	
	duration := time.Since(start)
	t.Logf("Created %d configs in %v (avg: %v per config)",
		iterations, duration, duration/time.Duration(iterations))
	
	// 验证性能在合理范围内
	assert.True(t, duration < time.Second,
		"Config creation should be fast, took %v", duration)
}

// BenchmarkConfigCloning 配置克隆的基准测试（模拟）
func TestBenchmarkConfigCloning(t *testing.T) {
	original := NewConfig().
		WithLevel(DEBUG).
		WithShowCaller(true).
		WithPrefix("ORIGINAL").
		WithColorful(false).
		WithTimeFormat("2006-01-02T15:04:05.000Z")
	
	start := time.Now()
	iterations := 10000
	
	for i := 0; i < iterations; i++ {
		clone := original.Clone()
		_ = clone
	}
	
	duration := time.Since(start)
	t.Logf("Cloned config %d times in %v (avg: %v per clone)",
		iterations, duration, duration/time.Duration(iterations))
	
	// 验证性能在合理范围内
	assert.True(t, duration < time.Second,
		"Config cloning should be fast, took %v", duration)
}