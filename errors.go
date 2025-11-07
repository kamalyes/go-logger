/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\errors.go
 * @Description: 统一的错误处理包
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"fmt"
	"runtime"
	"strings"
)

// ErrorCode 错误代码
type ErrorCode int

const (
	// 通用错误
	ErrUnknown ErrorCode = iota + 1000
	ErrInvalidInput
	ErrNotFound
	ErrAlreadyExists
	ErrPermissionDenied
	ErrTimeout
	ErrInternal

	// 日志相关错误
	ErrLoggerNotInitialized ErrorCode = iota + 2000
	ErrFormatterNotFound
	ErrWriterNotFound
	ErrHookNotFound
	ErrFilterNotFound
	ErrMiddlewareNotFound

	// 配置相关错误
	ErrConfigInvalid ErrorCode = iota + 3000
	ErrConfigNotFound
	ErrConfigLoadFailed
	ErrConfigSaveFailed

	// IO相关错误
	ErrFileNotFound ErrorCode = iota + 4000
	ErrFilePermission
	ErrFileWrite
	ErrFileRead
	ErrDiskFull
)

// 错误代码到消息的映射
var errorMessages = map[ErrorCode]string{
	ErrUnknown:          "未知错误",
	ErrInvalidInput:     "无效的输入参数",
	ErrNotFound:         "资源未找到",
	ErrAlreadyExists:    "资源已存在",
	ErrPermissionDenied: "权限被拒绝",
	ErrTimeout:          "操作超时",
	ErrInternal:         "内部错误",

	ErrLoggerNotInitialized: "日志器未初始化",
	ErrFormatterNotFound:    "格式化器未找到",
	ErrWriterNotFound:       "写入器未找到",
	ErrHookNotFound:         "钩子未找到",
	ErrFilterNotFound:       "过滤器未找到",
	ErrMiddlewareNotFound:   "中间件未找到",

	ErrConfigInvalid:    "无效的配置",
	ErrConfigNotFound:   "配置文件未找到",
	ErrConfigLoadFailed: "加载配置失败",
	ErrConfigSaveFailed: "保存配置失败",

	ErrFileNotFound:   "文件未找到",
	ErrFilePermission: "文件权限不足",
	ErrFileWrite:      "文件写入失败",
	ErrFileRead:       "文件读取失败",
	ErrDiskFull:       "磁盘空间不足",
}

// AppError 应用错误
type AppError struct {
	Code     ErrorCode              `json:"code"`
	Message  string                 `json:"message"`
	Details  string                 `json:"details"`
	Cause    error                  `json:"-"`
	File     string                 `json:"file"`
	Line     int                    `json:"line"`
	Function string                 `json:"function"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Cause
}

// String 返回详细的错误信息
func (e *AppError) String() string {
	var parts []string
	
	parts = append(parts, fmt.Sprintf("Error Code: %d", e.Code))
	parts = append(parts, fmt.Sprintf("Message: %s", e.Message))
	
	if e.Details != "" {
		parts = append(parts, fmt.Sprintf("Details: %s", e.Details))
	}
	
	if e.Cause != nil {
		parts = append(parts, fmt.Sprintf("Cause: %v", e.Cause))
	}
	
	if e.File != "" {
		parts = append(parts, fmt.Sprintf("Location: %s:%d", e.File, e.Line))
	}
	
	if e.Function != "" {
		parts = append(parts, fmt.Sprintf("Function: %s", e.Function))
	}
	
	if len(e.Context) > 0 {
		parts = append(parts, fmt.Sprintf("Context: %+v", e.Context))
	}
	
	return strings.Join(parts, "\n")
}

// WithContext 添加上下文信息
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// NewError 创建新的应用错误
func NewError(code ErrorCode, details string) *AppError {
	message := errorMessages[code]
	if message == "" {
		message = "未定义的错误"
	}

	// 获取调用者信息
	var file, function string
	var line int
	if pc, f, l, ok := runtime.Caller(1); ok {
		file = f
		line = l
		if fn := runtime.FuncForPC(pc); fn != nil {
			function = fn.Name()
			if idx := strings.LastIndex(function, "."); idx != -1 {
				function = function[idx+1:]
			}
		}
		if idx := strings.LastIndex(file, "/"); idx != -1 {
			file = file[idx+1:]
		}
	}

	return &AppError{
		Code:     code,
		Message:  message,
		Details:  details,
		File:     file,
		Line:     line,
		Function: function,
		Context:  make(map[string]interface{}),
	}
}

// WrapError 包装已有错误
func WrapError(code ErrorCode, details string, cause error) *AppError {
	err := NewError(code, details)
	err.Cause = cause
	return err
}

// NewErrorf 创建带格式化详情的应用错误
func NewErrorf(code ErrorCode, format string, args ...interface{}) *AppError {
	return NewError(code, fmt.Sprintf(format, args...))
}

// WrapErrorf 包装已有错误，带格式化详情
func WrapErrorf(code ErrorCode, cause error, format string, args ...interface{}) *AppError {
	return WrapError(code, fmt.Sprintf(format, args...), cause)
}

// IsError 检查错误是否为指定的错误代码
func IsError(err error, code ErrorCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// GetErrorCode 获取错误代码
func GetErrorCode(err error) ErrorCode {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ErrUnknown
}

// 便利方法

// NewInvalidInput 创建无效输入错误
func NewInvalidInput(details string) *AppError {
	return NewError(ErrInvalidInput, details)
}

// NewNotFound 创建未找到错误
func NewNotFound(resource string) *AppError {
	return NewError(ErrNotFound, fmt.Sprintf("资源 '%s' 未找到", resource))
}

// NewAlreadyExists 创建已存在错误
func NewAlreadyExists(resource string) *AppError {
	return NewError(ErrAlreadyExists, fmt.Sprintf("资源 '%s' 已存在", resource))
}

// NewPermissionDenied 创建权限拒绝错误
func NewPermissionDenied(action string) *AppError {
	return NewError(ErrPermissionDenied, fmt.Sprintf("无权限执行操作: %s", action))
}

// NewTimeout 创建超时错误
func NewTimeout(operation string, duration string) *AppError {
	return NewError(ErrTimeout, fmt.Sprintf("操作 '%s' 在 %s 内未完成", operation, duration))
}

// NewLoggerError 创建日志器相关错误
func NewLoggerError(code ErrorCode, component string) *AppError {
	return NewError(code, fmt.Sprintf("日志组件: %s", component))
}

// NewConfigError 创建配置相关错误
func NewConfigError(code ErrorCode, path string) *AppError {
	return NewError(code, fmt.Sprintf("配置路径: %s", path))
}

// NewFileError 创建文件相关错误
func NewFileError(code ErrorCode, filepath string) *AppError {
	return NewError(code, fmt.Sprintf("文件路径: %s", filepath))
}

// ErrorHandler 错误处理器
type ErrorHandler struct {
	logger ILogger
	hooks  []func(*AppError)
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(logger ILogger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
		hooks:  make([]func(*AppError), 0),
	}
}

// AddHook 添加错误处理钩子
func (eh *ErrorHandler) AddHook(hook func(*AppError)) {
	eh.hooks = append(eh.hooks, hook)
}

// Handle 处理错误
func (eh *ErrorHandler) Handle(err error) {
	if err == nil {
		return
	}
	
	var appErr *AppError
	if ae, ok := err.(*AppError); ok {
		appErr = ae
	} else {
		appErr = WrapError(ErrUnknown, "未知错误", err)
	}
	
	// 记录错误日志
	if eh.logger != nil {
		eh.logger.WithFields(map[string]interface{}{
			"error_code": appErr.Code,
			"file":       appErr.File,
			"line":       appErr.Line,
			"function":   appErr.Function,
			"context":    appErr.Context,
		}).Error(appErr.Error())
	}
	
	// 执行钩子
	for _, hook := range eh.hooks {
		go func(h func(*AppError)) {
			defer func() {
				if r := recover(); r != nil {
					if eh.logger != nil {
						eh.logger.Error("错误处理钩子发生panic: %v", r)
					}
				}
			}()
			h(appErr)
		}(hook)
	}
}

// 全局错误处理器
var globalErrorHandler *ErrorHandler

// InitErrorHandler 初始化全局错误处理器
func InitErrorHandler(logger ILogger) {
	globalErrorHandler = NewErrorHandler(logger)
}

// HandleError 使用全局错误处理器处理错误
func HandleError(err error) {
	if globalErrorHandler != nil {
		globalErrorHandler.Handle(err)
	}
}

// AddErrorHook 添加全局错误处理钩子
func AddErrorHook(hook func(*AppError)) {
	if globalErrorHandler != nil {
		globalErrorHandler.AddHook(hook)
	}
}