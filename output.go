/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-02-02 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-02-02 00:00:00
 * @FilePath: \go-logger\output.go
 * @Description: 输出类型和配置
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"io"

	"github.com/kamalyes/go-toolbox/pkg/mathx"
)

// OutputType 输出类型
type OutputType string

const (
	OutputConsole OutputType = "console" // 控制台输出
	OutputFile    OutputType = "file"    // 文件输出
	OutputRotate  OutputType = "rotate"  // 轮转文件输出
	OutputStdout  OutputType = "stdout"  // 标准输出
	OutputStderr  OutputType = "stderr"  // 标准错误输出
)

// 默认配置常量
const (
	DefaultMaxSize  = 100 * 1024 * 1024 // 默认最大文件大小 100MB
	DefaultMaxFiles = 5                 // 默认最大文件数
)

// 错误消息常量
const (
	ErrMsgFilePathEmpty   = "file_path is required for file output"
	ErrMsgRotatePathEmpty = "file_path is required for rotate output"
)

// WriterConfig Writer 配置
type WriterConfig struct {
	Type     OutputType `json:"type" yaml:"type"`           // 输出类型
	FilePath string     `json:"file_path" yaml:"file_path"` // 文件路径
	MaxSize  int64      `json:"max_size" yaml:"max_size"`   // 最大文件大小(字节)
	MaxFiles int        `json:"max_files" yaml:"max_files"` // 最大文件数
	MaxAge   int        `json:"max_age" yaml:"max_age"`     // 最大保存天数
	Compress bool       `json:"compress" yaml:"compress"`   // 是否压缩
	Output   io.Writer  `json:"-" yaml:"-"`                 // 自定义输出
}

// CreateWriter 根据配置创建 Writer
func CreateWriter(config *WriterConfig) (IWriter, error) {
	if config == nil {
		return NewConsoleWriter(nil), nil
	}

	switch config.Type {
	case OutputConsole, OutputStdout:
		return createConsoleWriter(config), nil

	case OutputFile:
		return createFileWriter(config)

	case OutputRotate:
		return createRotateWriter(config)

	default:
		return NewConsoleWriter(nil), nil
	}
}

// createConsoleWriter 创建控制台 Writer
func createConsoleWriter(config *WriterConfig) IWriter {
	if config.Output != nil {
		return NewConsoleWriter(config.Output)
	}
	return NewConsoleWriter(nil)
}

// createFileWriter 创建文件 Writer
func createFileWriter(config *WriterConfig) (IWriter, error) {
	if config.FilePath == "" {
		return nil, NewConfigError(ErrInvalidInput, ErrMsgFilePathEmpty)
	}
	return NewFileWriter(config.FilePath), nil
}

// createRotateWriter 创建轮转文件 Writer
func createRotateWriter(config *WriterConfig) (IWriter, error) {
	if config.FilePath == "" {
		return nil, NewConfigError(ErrInvalidInput, ErrMsgRotatePathEmpty)
	}
	maxSize := mathx.IfNotZero(config.MaxSize, DefaultMaxSize)
	maxFiles := mathx.IfNotZero(config.MaxFiles, DefaultMaxFiles)
	return NewRotateWriter(config.FilePath, maxSize, maxFiles), nil
}

// AddWriterFromConfig 根据配置添加 Writer 到 Builder
func AddWriterFromConfig(builder *LoggerBuilder, config *WriterConfig) error {
	if config == nil {
		addConsoleWriter(builder)
		return nil
	}

	switch config.Type {
	case OutputConsole, OutputStdout:
		addConsoleWriter(builder)

	case OutputFile:
		return addFileWriter(builder, config)

	case OutputRotate:
		return addRotateWriter(builder, config)

	default:
		addConsoleWriter(builder)
	}

	return nil
}

// addConsoleWriter 添加控制台 Writer 到 Builder
func addConsoleWriter(builder *LoggerBuilder) {
	builder.WithWriter("console", nil)
}

// addFileWriter 添加文件 Writer 到 Builder
func addFileWriter(builder *LoggerBuilder, config *WriterConfig) error {
	if config.FilePath == "" {
		return NewConfigError(ErrInvalidInput, ErrMsgFilePathEmpty)
	}
	writerConfig := map[string]any{
		"file_path": config.FilePath,
	}
	builder.WithWriter("file", writerConfig)
	return nil
}

// addRotateWriter 添加轮转文件 Writer 到 Builder
func addRotateWriter(builder *LoggerBuilder, config *WriterConfig) error {
	if config.FilePath == "" {
		return NewConfigError(ErrInvalidInput, ErrMsgRotatePathEmpty)
	}
	maxSize := mathx.IfNotZero(config.MaxSize, DefaultMaxSize)
	maxFiles := mathx.IfNotZero(config.MaxFiles, DefaultMaxFiles)
	writerConfig := map[string]any{
		"file_path": config.FilePath,
		"max_size":  maxSize,
		"max_files": maxFiles,
	}
	builder.WithWriter("rotate", writerConfig)
	return nil
}
