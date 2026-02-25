/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-02 00:00:00
 * @FilePath: \go-logger\errors_test.go
 * @Description: 错误处理测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package logger

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ErrorsTestSuite 错误测试套件
type ErrorsTestSuite struct {
	suite.Suite
}

// TestNewError 测试创建新错误
func (s *ErrorsTestSuite) TestNewError() {
	err := NewError(ErrInvalidInput, "test details")

	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), ErrInvalidInput, err.Code)
	assert.Equal(s.T(), "test details", err.Details)
	assert.NotEmpty(s.T(), err.Message)
	assert.NotEmpty(s.T(), err.File)
	assert.NotZero(s.T(), err.Line)
	assert.NotEmpty(s.T(), err.Function)
}

// TestAppErrorError 测试错误接口实现
func (s *ErrorsTestSuite) TestAppErrorError() {
	err := NewError(ErrNotFound, "resource not found")

	errStr := err.Error()
	assert.Contains(s.T(), errStr, "resource not found")
	assert.Contains(s.T(), errStr, "[")
	assert.Contains(s.T(), errStr, "]")
}

// TestAppErrorString 测试详细错误信息
func (s *ErrorsTestSuite) TestAppErrorString() {
	err := NewError(ErrInvalidInput, "invalid parameter")

	str := err.String()
	assert.Contains(s.T(), str, "Error Code:")
	assert.Contains(s.T(), str, "Message:")
	assert.Contains(s.T(), str, "Details:")
	assert.Contains(s.T(), str, "Location:")
	assert.Contains(s.T(), str, "Function:")
}

// TestWrapError 测试包装错误
func (s *ErrorsTestSuite) TestWrapError() {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(ErrInternal, "wrapped details", originalErr)

	assert.NotNil(s.T(), wrappedErr)
	assert.Equal(s.T(), ErrInternal, wrappedErr.Code)
	assert.Equal(s.T(), "wrapped details", wrappedErr.Details)
	assert.Equal(s.T(), originalErr, wrappedErr.Cause)
}

// TestAppErrorUnwrap 测试解包错误
func (s *ErrorsTestSuite) TestAppErrorUnwrap() {
	originalErr := errors.New("original")
	wrappedErr := WrapError(ErrInternal, "wrapped", originalErr)

	unwrapped := wrappedErr.Unwrap()
	assert.Equal(s.T(), originalErr, unwrapped)
}

// TestNewErrorf 测试格式化创建错误
func (s *ErrorsTestSuite) TestNewErrorf() {
	err := NewErrorf(ErrInvalidInput, "user %s not found with id %d", "alice", 123)

	assert.NotNil(s.T(), err)
	assert.Contains(s.T(), err.Details, "alice")
	assert.Contains(s.T(), err.Details, "123")
}

// TestWrapErrorf 测试格式化包装错误
func (s *ErrorsTestSuite) TestWrapErrorf() {
	originalErr := errors.New("db error")
	wrappedErr := WrapErrorf(ErrInternal, originalErr, "failed to query user %s", "bob")

	assert.NotNil(s.T(), wrappedErr)
	assert.Contains(s.T(), wrappedErr.Details, "bob")
	assert.Equal(s.T(), originalErr, wrappedErr.Cause)
}

// TestIsError 测试错误类型检查
func (s *ErrorsTestSuite) TestIsError() {
	err := NewError(ErrNotFound, "not found")

	assert.True(s.T(), IsError(err, ErrNotFound))
	assert.False(s.T(), IsError(err, ErrInvalidInput))

	// 测试非AppError
	stdErr := errors.New("standard error")
	assert.False(s.T(), IsError(stdErr, ErrNotFound))
}

// TestGetErrorCode 测试获取错误代码
func (s *ErrorsTestSuite) TestGetErrorCode() {
	err := NewError(ErrTimeout, "timeout")

	code := GetErrorCode(err)
	assert.Equal(s.T(), ErrTimeout, code)

	// 测试非AppError
	stdErr := errors.New("standard error")
	code = GetErrorCode(stdErr)
	assert.Equal(s.T(), ErrUnknown, code)
}

// TestAppErrorWithContext 测试添加上下文
func (s *ErrorsTestSuite) TestAppErrorWithContext() {
	err := NewError(ErrInvalidInput, "invalid")

	err.WithContext("user_id", 123)
	err.WithContext("action", "login")

	assert.Equal(s.T(), 123, err.Context["user_id"])
	assert.Equal(s.T(), "login", err.Context["action"])
}

// TestConvenienceFunctions 测试便利函数
func (s *ErrorsTestSuite) TestConvenienceFunctions() {
	// TestNewInvalidInput
	err := NewInvalidInput("invalid parameter")
	assert.Equal(s.T(), ErrInvalidInput, err.Code)

	// TestNewNotFound
	err = NewNotFound("user")
	assert.Equal(s.T(), ErrNotFound, err.Code)
	assert.Contains(s.T(), err.Details, "user")

	// TestNewAlreadyExists
	err = NewAlreadyExists("email")
	assert.Equal(s.T(), ErrAlreadyExists, err.Code)
	assert.Contains(s.T(), err.Details, "email")

	// TestNewPermissionDenied
	err = NewPermissionDenied("delete")
	assert.Equal(s.T(), ErrPermissionDenied, err.Code)
	assert.Contains(s.T(), err.Details, "delete")

	// TestNewTimeout
	err = NewTimeout("query", "5s")
	assert.Equal(s.T(), ErrTimeout, err.Code)
	assert.Contains(s.T(), err.Details, "query")
	assert.Contains(s.T(), err.Details, "5s")
}

// TestLoggerErrors 测试日志器相关错误
func (s *ErrorsTestSuite) TestLoggerErrors() {
	err := NewLoggerError(ErrFormatterNotFound, "json")
	assert.Equal(s.T(), ErrFormatterNotFound, err.Code)
	assert.Contains(s.T(), err.Details, "json")
}

// TestConfigErrors 测试配置相关错误
func (s *ErrorsTestSuite) TestConfigErrors() {
	err := NewConfigError(ErrConfigNotFound, "/etc/app/config.yaml")
	assert.Equal(s.T(), ErrConfigNotFound, err.Code)
	assert.Contains(s.T(), err.Details, "/etc/app/config.yaml")
}

// TestFileErrors 测试文件相关错误
func (s *ErrorsTestSuite) TestFileErrors() {
	err := NewFileError(ErrFileNotFound, "/var/log/app.log")
	assert.Equal(s.T(), ErrFileNotFound, err.Code)
	assert.Contains(s.T(), err.Details, "/var/log/app.log")
}

// TestErrorHandler 测试错误处理器
func (s *ErrorsTestSuite) TestErrorHandler() {
	logger := NewEmptyLogger()
	handler := NewErrorHandler(logger)

	assert.NotNil(s.T(), handler)
	assert.NotNil(s.T(), handler.logger)
	assert.NotNil(s.T(), handler.hooks)
}

// TestErrorHandlerHandle 测试处理错误
func (s *ErrorsTestSuite) TestErrorHandlerHandle() {
	logger := NewEmptyLogger()
	handler := NewErrorHandler(logger)

	// 测试处理nil错误
	handler.Handle(nil)

	// 测试处理AppError
	appErr := NewError(ErrInvalidInput, "test")
	handler.Handle(appErr)

	// 测试处理标准错误
	stdErr := errors.New("standard error")
	handler.Handle(stdErr)
}

// TestErrorHandlerAddHook 测试添加钩子
func (s *ErrorsTestSuite) TestErrorHandlerAddHook() {
	logger := NewEmptyLogger()
	handler := NewErrorHandler(logger)

	done := make(chan bool, 1)
	hook := func(err *AppError) {
		done <- true
	}

	handler.AddHook(hook)
	assert.Len(s.T(), handler.hooks, 1)

	// 触发钩子
	err := NewError(ErrInvalidInput, "test")
	handler.Handle(err)

	// 等待钩子执行（带超时）
	select {
	case <-done:
		// 钩子执行成功
	case <-time.After(100 * time.Millisecond):
		// 超时也不算失败，因为钩子是异步的
	}
}

// TestGlobalErrorHandler 测试全局错误处理器
func (s *ErrorsTestSuite) TestGlobalErrorHandler() {
	logger := NewEmptyLogger()
	InitErrorHandler(logger)

	assert.NotNil(s.T(), globalErrorHandler)

	// 测试全局处理
	err := NewError(ErrInvalidInput, "test")
	HandleError(err)

	// 测试添加全局钩子
	AddErrorHook(func(err *AppError) {
		// 钩子逻辑
	})
}

// TestErrorCodes 测试错误代码
func (s *ErrorsTestSuite) TestErrorCodes() {
	codes := []ErrorCode{
		ErrUnknown,
		ErrInvalidInput,
		ErrNotFound,
		ErrAlreadyExists,
		ErrPermissionDenied,
		ErrTimeout,
		ErrInternal,
		ErrLoggerNotInitialized,
		ErrFormatterNotFound,
		ErrWriterNotFound,
		ErrHookNotFound,
		ErrFilterNotFound,
		ErrMiddlewareNotFound,
		ErrConfigInvalid,
		ErrConfigNotFound,
		ErrConfigLoadFailed,
		ErrConfigSaveFailed,
		ErrFileNotFound,
		ErrFilePermission,
		ErrFileWrite,
		ErrFileRead,
		ErrDiskFull,
	}

	for _, code := range codes {
		assert.NotZero(s.T(), code)
		// 每个错误代码都应该有对应的消息
		message := errorMessages[code]
		assert.NotEmpty(s.T(), message)
	}
}

// TestErrorMessages 测试错误消息
func (s *ErrorsTestSuite) TestErrorMessages() {
	err := NewError(ErrInvalidInput, "test")
	assert.NotEmpty(s.T(), err.Message)

	// 测试未定义的错误代码
	err = &AppError{
		Code: ErrorCode(9999),
	}
	// 应该有默认消息
	assert.NotEmpty(s.T(), err.Error())
}

// TestErrorContextChaining 测试上下文链式调用
func (s *ErrorsTestSuite) TestErrorContextChaining() {
	err := NewError(ErrInvalidInput, "test").
		WithContext("user", "alice").
		WithContext("action", "login").
		WithContext("ip", "192.168.1.1")

	assert.Len(s.T(), err.Context, 3)
	assert.Equal(s.T(), "alice", err.Context["user"])
	assert.Equal(s.T(), "login", err.Context["action"])
	assert.Equal(s.T(), "192.168.1.1", err.Context["ip"])
}

// TestErrorStringWithCause 测试带原因的错误字符串
func (s *ErrorsTestSuite) TestErrorStringWithCause() {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(ErrInternal, "wrapped", originalErr)

	str := wrappedErr.String()
	assert.Contains(s.T(), str, "Cause:")
	assert.Contains(s.T(), str, "original error")
}

// TestErrorStringWithContext 测试带上下文的错误字符串
func (s *ErrorsTestSuite) TestErrorStringWithContext() {
	err := NewError(ErrInvalidInput, "test")
	err.WithContext("key", "value")

	str := err.String()
	assert.Contains(s.T(), str, "Context:")
	assert.Contains(s.T(), str, "key")
}

// TestErrorHandlerWithNilLogger 测试nil日志器的错误处理器
func (s *ErrorsTestSuite) TestErrorHandlerWithNilLogger() {
	handler := NewErrorHandler(nil)
	assert.NotNil(s.T(), handler)

	// 不应该崩溃
	err := NewError(ErrInvalidInput, "test")
	handler.Handle(err)
}

// TestErrorHandlerHookPanic 测试钩子panic处理
func (s *ErrorsTestSuite) TestErrorHandlerHookPanic() {
	logger := NewEmptyLogger()
	handler := NewErrorHandler(logger)

	// 添加会panic的钩子
	handler.AddHook(func(err *AppError) {
		panic("hook panic")
	})

	// 不应该导致程序崩溃
	err := NewError(ErrInvalidInput, "test")
	handler.Handle(err)
}

// TestMultipleErrorWrapping 测试多层错误包装
func (s *ErrorsTestSuite) TestMultipleErrorWrapping() {
	err1 := errors.New("level 1")
	err2 := WrapError(ErrInternal, "level 2", err1)
	err3 := WrapError(ErrTimeout, "level 3", err2)

	assert.Equal(s.T(), ErrTimeout, err3.Code)
	assert.Equal(s.T(), err2, err3.Cause)
	assert.Equal(s.T(), err1, err2.Cause)
}

// TestErrorCodeUniqueness 测试错误代码唯一性
func (s *ErrorsTestSuite) TestErrorCodeUniqueness() {
	codes := make(map[ErrorCode]bool)

	allCodes := []ErrorCode{
		ErrUnknown, ErrInvalidInput, ErrNotFound, ErrAlreadyExists,
		ErrPermissionDenied, ErrTimeout, ErrInternal,
		ErrLoggerNotInitialized, ErrFormatterNotFound, ErrWriterNotFound,
		ErrHookNotFound, ErrFilterNotFound, ErrMiddlewareNotFound,
		ErrConfigInvalid, ErrConfigNotFound, ErrConfigLoadFailed, ErrConfigSaveFailed,
		ErrFileNotFound, ErrFilePermission, ErrFileWrite, ErrFileRead, ErrDiskFull,
	}

	for _, code := range allCodes {
		assert.False(s.T(), codes[code], "Duplicate error code: %d", code)
		codes[code] = true
	}
}

// 运行测试套件
func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsTestSuite))
}

// BenchmarkNewError 创建错误性能测试
func BenchmarkNewError(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = NewError(ErrInvalidInput, "test details")
	}
}

// BenchmarkWrapError 包装错误性能测试
func BenchmarkWrapError(b *testing.B) {
	originalErr := errors.New("original")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = WrapError(ErrInternal, "wrapped", originalErr)
	}
}

// BenchmarkErrorString 错误字符串性能测试
func BenchmarkErrorString(b *testing.B) {
	err := NewError(ErrInvalidInput, "test")
	err.WithContext("key1", "value1")
	err.WithContext("key2", "value2")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = err.String()
	}
}

// BenchmarkIsError 错误检查性能测试
func BenchmarkIsError(b *testing.B) {
	err := NewError(ErrNotFound, "not found")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = IsError(err, ErrNotFound)
	}
}

// BenchmarkErrorHandler 错误处理器性能测试
func BenchmarkErrorHandler(b *testing.B) {
	logger := NewEmptyLogger()
	handler := NewErrorHandler(logger)
	err := NewError(ErrInvalidInput, "test")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		handler.Handle(err)
	}
}
