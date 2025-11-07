/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 02:03:28
 * @FilePath: \go-logger\factory.go
 * @Description: 工厂模式 - 统一管理各种组件的创建
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"fmt"
	"os"
	"sync"
)

// ComponentFactory 组件工厂接口
type ComponentFactory interface {
	GetName() string
	GetVersion() string
	CreateComponent(config interface{}) (interface{}, error)
	ValidateConfig(config interface{}) error
}

// LoggerFactory 日志器工厂
type LoggerFactory struct {
	formatters  *FormatRegistry
	adapters    *AdapterRegistry
	writers     map[string]func(config interface{}) (IWriter, error)
	hooks       map[string]func(config interface{}) (IHook, error)
	middlewares map[string]func(config interface{}) (IMiddleware, error)
	mutex       sync.RWMutex
}

// NewLoggerFactory 创建日志器工厂
func NewLoggerFactory() *LoggerFactory {
	factory := &LoggerFactory{
		formatters:  NewFormatRegistry(),
		adapters:    NewAdapterRegistry(),
		writers:     make(map[string]func(config interface{}) (IWriter, error)),
		hooks:       make(map[string]func(config interface{}) (IHook, error)),
		middlewares: make(map[string]func(config interface{}) (IMiddleware, error)),
	}
	
	// 注册默认组件
	factory.registerDefaultComponents()
	
	return factory
}

// registerDefaultComponents 注册默认组件
func (f *LoggerFactory) registerDefaultComponents() {
	// 注册写入器
	f.RegisterWriter("console", func(config interface{}) (IWriter, error) {
		if config == nil {
			return NewConsoleWriter(os.Stdout), nil
		}
		// 这里可以添加配置解析逻辑
		return NewConsoleWriter(os.Stdout), nil
	})
	
	f.RegisterWriter("file", func(config interface{}) (IWriter, error) {
		configMap, ok := config.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid file writer config")
		}
		
		filePath, ok := configMap["file_path"].(string)
		if !ok {
			return nil, fmt.Errorf("file_path is required for file writer")
		}
		
		return NewFileWriter(filePath), nil
	})
	
	f.RegisterWriter("rotate", func(config interface{}) (IWriter, error) {
		configMap, ok := config.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid rotate writer config")
		}
		
		filePath, ok := configMap["file_path"].(string)
		if !ok {
			return nil, fmt.Errorf("file_path is required for rotate writer")
		}
		
		maxSize := int64(100 * 1024 * 1024) // 默认100MB
		if size, exists := configMap["max_size"]; exists {
			if sizeVal, ok := size.(int64); ok {
				maxSize = sizeVal
			}
		}
		
		maxFiles := 5 // 默认5个文件
		if files, exists := configMap["max_files"]; exists {
			if filesVal, ok := files.(int); ok {
				maxFiles = filesVal
			}
		}
		
		return NewRotateWriter(filePath, maxSize, maxFiles), nil
	})
	
	// 注册钩子
	f.RegisterHook("console", func(config interface{}) (IHook, error) {
		levels := AllLevels
		if config != nil {
			if configMap, ok := config.(map[string]interface{}); ok {
				if levelStr, exists := configMap["levels"]; exists {
					// 解析级别配置
					_ = levelStr
				}
			}
		}
		return NewConsoleHook(levels), nil
	})
	
	f.RegisterHook("file", func(config interface{}) (IHook, error) {
		configMap, ok := config.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid file hook config")
		}
		
		filePath, ok := configMap["file_path"].(string)
		if !ok {
			return nil, fmt.Errorf("file_path is required for file hook")
		}
		
		levels := AllLevels
		return NewFileHook(filePath, levels), nil
	})
	
	f.RegisterHook("webhook", func(config interface{}) (IHook, error) {
		configMap, ok := config.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid webhook hook config")
		}
		
		url, ok := configMap["url"].(string)
		if !ok {
			return nil, fmt.Errorf("url is required for webhook hook")
		}
		
		levels := ErrorLevels // 默认只发送错误级别
		return NewWebhookHook(url, levels), nil
	})
}

// RegisterWriter 注册写入器工厂
func (f *LoggerFactory) RegisterWriter(name string, factory func(config interface{}) (IWriter, error)) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.writers[name] = factory
}

// RegisterHook 注册钩子工厂
func (f *LoggerFactory) RegisterHook(name string, factory func(config interface{}) (IHook, error)) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.hooks[name] = factory
}

// RegisterMiddleware 注册中间件工厂
func (f *LoggerFactory) RegisterMiddleware(name string, factory func(config interface{}) (IMiddleware, error)) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.middlewares[name] = factory
}

// CreateFormatter 创建格式化器
func (f *LoggerFactory) CreateFormatter(formatterType FormatterType) (IFormatter, error) {
	return f.formatters.Create(formatterType)
}

// CreateAdapter 创建适配器
func (f *LoggerFactory) CreateAdapter(name string, config *AdapterConfig) (IAdapter, error) {
	return f.adapters.Create(name, config)
}

// CreateWriter 创建写入器
func (f *LoggerFactory) CreateWriter(name string, config interface{}) (IWriter, error) {
	f.mutex.RLock()
	factory, exists := f.writers[name]
	f.mutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("unknown writer type: %s", name)
	}
	
	return factory(config)
}

// CreateHook 创建钩子
func (f *LoggerFactory) CreateHook(name string, config interface{}) (IHook, error) {
	f.mutex.RLock()
	factory, exists := f.hooks[name]
	f.mutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("unknown hook type: %s", name)
	}
	
	return factory(config)
}

// CreateMiddleware 创建中间件
func (f *LoggerFactory) CreateMiddleware(name string, config interface{}) (IMiddleware, error) {
	f.mutex.RLock()
	factory, exists := f.middlewares[name]
	f.mutex.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("unknown middleware type: %s", name)
	}
	
	return factory(config)
}

// LoggerBuilder 日志器构建器
type LoggerBuilder struct {
	factory     *LoggerFactory
	config      *LogConfig
	formatter   IFormatter
	writers     []IWriter
	hooks       []IHook
	middlewares []IMiddleware
}

// NewLoggerBuilder 创建日志器构建器
func NewLoggerBuilder() *LoggerBuilder {
	return &LoggerBuilder{
		factory:     NewLoggerFactory(),
		config:      DefaultConfig(),
		writers:     make([]IWriter, 0),
		hooks:       make([]IHook, 0),
		middlewares: make([]IMiddleware, 0),
	}
}

// WithConfig 设置配置
func (b *LoggerBuilder) WithConfig(config *LogConfig) *LoggerBuilder {
	b.config = config
	return b
}

// WithFormatter 设置格式化器
func (b *LoggerBuilder) WithFormatter(formatterType FormatterType) *LoggerBuilder {
	if formatter, err := b.factory.CreateFormatter(formatterType); err == nil {
		b.formatter = formatter
	}
	return b
}

// WithWriter 添加写入器
func (b *LoggerBuilder) WithWriter(name string, config interface{}) *LoggerBuilder {
	if writer, err := b.factory.CreateWriter(name, config); err == nil {
		b.writers = append(b.writers, writer)
	}
	return b
}

// WithHook 添加钩子
func (b *LoggerBuilder) WithHook(name string, config interface{}) *LoggerBuilder {
	if hook, err := b.factory.CreateHook(name, config); err == nil {
		b.hooks = append(b.hooks, hook)
	}
	return b
}

// WithMiddleware 添加中间件
func (b *LoggerBuilder) WithMiddleware(name string, config interface{}) *LoggerBuilder {
	if middleware, err := b.factory.CreateMiddleware(name, config); err == nil {
		b.middlewares = append(b.middlewares, middleware)
	}
	return b
}

// Build 构建日志器
func (b *LoggerBuilder) Build() *Logger {
	options := &LoggerOptions{
		Config:     b.config,
		Formatter:  b.formatter,
		Writers:    b.writers,
		Hooks:      b.hooks,
		Middleware: b.middlewares,
	}
	
	return NewLoggerWithOptions(options)
}

// 全局工厂实例
var globalFactory = NewLoggerFactory()

// GetGlobalFactory 获取全局工厂实例
func GetGlobalFactory() *LoggerFactory {
	return globalFactory
}

// SetGlobalFactory 设置全局工厂实例
func SetGlobalFactory(factory *LoggerFactory) {
	globalFactory = factory
}

// 便利函数

// CreateSimpleLogger 创建简单的日志器
func CreateSimpleLogger(level LogLevel) *Logger {
	return NewLoggerBuilder().
		WithConfig(DefaultConfig().WithLevel(level)).
		WithFormatter(TextFormatter).
		WithWriter("console", nil).
		Build()
}

// CreateFileLogger 创建文件日志器
func CreateFileLogger(filePath string, level LogLevel) *Logger {
	writerConfig := map[string]interface{}{
		"file_path": filePath,
	}
	
	return NewLoggerBuilder().
		WithConfig(DefaultConfig().WithLevel(level)).
		WithFormatter(JSONFormatter).
		WithWriter("file", writerConfig).
		Build()
}

// CreateProductionLogger 创建生产环境日志器
func CreateProductionLogger(logDir string) *Logger {
	// 控制台输出配置
	consoleConfig := map[string]interface{}{}
	
	// 文件输出配置
	fileConfig := map[string]interface{}{
		"file_path": logDir + "/app.log",
		"max_size":  100 * 1024 * 1024, // 100MB
		"max_files": 10,
	}
	
	// 错误文件配置
	errorConfig := map[string]interface{}{
		"file_path": logDir + "/error.log",
		"max_size":  50 * 1024 * 1024, // 50MB
		"max_files": 5,
	}
	
	// 速率限制配置
	rateLimitConfig := map[string]interface{}{
		"max_rate":    1000,
		"time_window": "1s",
		"burst_size":  100,
	}
	
	// 指标中间件配置
	metricsConfig := map[string]interface{}{}
	
	return NewLoggerBuilder().
		WithConfig(DefaultConfig().WithLevel(INFO).WithShowCaller(true)).
		WithFormatter(JSONFormatter).
		WithWriter("console", consoleConfig).
		WithWriter("rotate", fileConfig).
		WithWriter("rotate", errorConfig).
		WithMiddleware("metrics", metricsConfig).
		WithMiddleware("ratelimit", rateLimitConfig).
		Build()
}