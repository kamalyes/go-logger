/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\config.go
 * @Description: 日志配置
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"io"
	"os"
	"strings"
)

// LogConfig 日志配置
type LogConfig struct {
	Level      LogLevel `json:"level"`       // 日志级别
	ShowCaller bool     `json:"show_caller"` // 是否显示调用者信息
	Prefix     string   `json:"prefix"`      // 日志前缀
	Output     io.Writer `json:"-"`          // 输出目标
	Colorful   bool     `json:"colorful"`    // 是否使用彩色输出
	TimeFormat string   `json:"time_format"` // 时间格式
}

// DefaultConfig 默认配置
func DefaultConfig() *LogConfig {
	return &LogConfig{
		Level:      INFO,
		ShowCaller: false,
		Prefix:     "",
		Output:     os.Stdout,
		Colorful:   true,
		TimeFormat: "2006-01-02 15:04:05",
	}
}

// NewConfig 创建新的配置
func NewConfig() *LogConfig {
	return DefaultConfig()
}

// WithLevel 设置日志级别
func (c *LogConfig) WithLevel(level LogLevel) *LogConfig {
	c.Level = level
	return c
}

// WithShowCaller 设置是否显示调用者信息
func (c *LogConfig) WithShowCaller(show bool) *LogConfig {
	c.ShowCaller = show
	return c
}

// WithPrefix 设置日志前缀
func (c *LogConfig) WithPrefix(prefix string) *LogConfig {
	if prefix != "" && !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}
	c.Prefix = prefix
	return c
}

// WithOutput 设置输出目标
func (c *LogConfig) WithOutput(output io.Writer) *LogConfig {
	c.Output = output
	return c
}

// WithColorful 设置是否使用彩色输出
func (c *LogConfig) WithColorful(colorful bool) *LogConfig {
	c.Colorful = colorful
	return c
}

// WithTimeFormat 设置时间格式
func (c *LogConfig) WithTimeFormat(format string) *LogConfig {
	c.TimeFormat = format
	return c
}

// Clone 克隆配置
func (c *LogConfig) Clone() *LogConfig {
	return &LogConfig{
		Level:      c.Level,
		ShowCaller: c.ShowCaller,
		Prefix:     c.Prefix,
		Output:     c.Output,
		Colorful:   c.Colorful,
		TimeFormat: c.TimeFormat,
	}
}

// Validate 验证配置
func (c *LogConfig) Validate() error {
	if c.Output == nil {
		c.Output = os.Stdout
	}
	if c.TimeFormat == "" {
		c.TimeFormat = "2006-01-02 15:04:05"
	}
	return nil
}