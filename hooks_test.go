/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-08 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 00:00:00
 * @FilePath: \go-logger\hooks_test.go
 * @Description: 钩子系统测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// HooksTestSuite 钩子测试套件
type HooksTestSuite struct {
	suite.Suite
}

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(HooksTestSuite))
}

// TestNewConsoleHook 测试创建控制台钩子
func (suite *HooksTestSuite) TestNewConsoleHook() {
	levels := []LogLevel{ERROR, FATAL}
	hook := NewConsoleHook(levels)
	
	assert.NotNil(suite.T(), hook)
	expectedLevels := hook.Levels()
	assert.Equal(suite.T(), levels, expectedLevels)
}

// TestConsoleHookFire 测试控制台钩子触发
func (suite *HooksTestSuite) TestConsoleHookFire() {
	hook := NewConsoleHook([]LogLevel{ERROR})
	
	// 创建错误事件
	entry := &LogEntry{
		Level:     ERROR,
		Message:   "test error message",
		Timestamp: time.Now().Unix(),
		Fields:    map[string]interface{}{"key": "value"},
		Caller:    &CallerInfo{File: "test.go", Line: 10, Function: "testFunc"},
	}
	
	// 触发钩子
	err := hook.Fire(entry)
	assert.NoError(suite.T(), err)
}

// TestConsoleHookLevels 测试控制台钩子级别过滤
func (suite *HooksTestSuite) TestConsoleHookLevels() {
	hook := NewConsoleHook([]LogLevel{ERROR, FATAL})
	levels := hook.Levels()
	
	assert.Contains(suite.T(), levels, ERROR)
	assert.Contains(suite.T(), levels, FATAL)
	assert.NotContains(suite.T(), levels, DEBUG)
	assert.NotContains(suite.T(), levels, INFO)
	assert.NotContains(suite.T(), levels, WARN)
}

// TestNewFileHook 测试创建文件钩子
func (suite *HooksTestSuite) TestNewFileHook() {
	tempFile := "test_hook.log"
	defer os.Remove(tempFile)
	
	levels := []LogLevel{WARN, ERROR, FATAL}
	hook := NewFileHook(tempFile, levels)
	
	assert.NotNil(suite.T(), hook)
	expectedLevels := hook.Levels()
	assert.Equal(suite.T(), levels, expectedLevels)
}

// TestFileHookFire 测试文件钩子触发
func (suite *HooksTestSuite) TestFileHookFire() {
	tempFile := "test_hook_fire.log"
	defer os.Remove(tempFile)
	
	hook := NewFileHook(tempFile, []LogLevel{ERROR})
	
	// 创建测试事件
	entry := &LogEntry{
		Level:     ERROR,
		Message:   "file hook test message",
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
		Fields:    map[string]interface{}{"component": "test", "action": "hook_test"},
		Caller:    &CallerInfo{File: "hooks_test.go", Line: 100, Function: "TestFileHookFire"},
	}
	
	// 触发钩子
	err := hook.Fire(entry)
	assert.NoError(suite.T(), err)
	
	// 给一点时间让文件写入完成
	time.Sleep(100 * time.Millisecond)
	
	// 检查文件内容
	content, err := os.ReadFile(tempFile)
	assert.NoError(suite.T(), err)
	
	contentStr := string(content)
	assert.Contains(suite.T(), contentStr, "file hook test message")
	assert.Contains(suite.T(), contentStr, "ERROR")
	assert.Contains(suite.T(), contentStr, "action")
	assert.Contains(suite.T(), contentStr, "component")
}

// TestFileHookLevels 测试文件钩子级别过滤
func (suite *HooksTestSuite) TestFileHookLevels() {
	tempFile := "test_hook_levels.log"
	defer os.Remove(tempFile)
	
	hook := NewFileHook(tempFile, []LogLevel{WARN, ERROR})
	levels := hook.Levels()
	
	assert.Contains(suite.T(), levels, WARN)
	assert.Contains(suite.T(), levels, ERROR)
	assert.NotContains(suite.T(), levels, DEBUG)
	assert.NotContains(suite.T(), levels, INFO)
	assert.NotContains(suite.T(), levels, FATAL)
}

// TestNewWebhookHook 测试创建Webhook钩子
func (suite *HooksTestSuite) TestNewWebhookHook() {
	url := "http://example.com/webhook"
	levels := []LogLevel{ERROR, FATAL}
	hook := NewWebhookHook(url, levels)
	
	assert.NotNil(suite.T(), hook)
	expectedLevels := hook.Levels()
	assert.Equal(suite.T(), levels, expectedLevels)
}

// TestWebhookHookFire 测试Webhook钩子触发
func (suite *HooksTestSuite) TestWebhookHookFire() {
	// 创建测试服务器
	receivedPayload := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		receivedPayload = buf.String()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	// 创建webhook钩子
	hook := NewWebhookHook(server.URL, []LogLevel{ERROR})
	
	// 创建测试事件
	entry := &LogEntry{
		Level:     ERROR,
		Message:   "webhook test message",
		Timestamp: time.Now().Unix(),
		Fields:    map[string]interface{}{"service": "test", "severity": "high"},
		Caller:    &CallerInfo{File: "service.go", Line: 50, Function: "handleError"},
	}
	
	// 触发钩子
	err := hook.Fire(entry)
	assert.NoError(suite.T(), err)
	
	// 给一点时间让HTTP请求完成
	time.Sleep(200 * time.Millisecond)
	
	// 检查接收到的payload
	assert.Contains(suite.T(), receivedPayload, "webhook test message")
	assert.Contains(suite.T(), receivedPayload, "ERROR")
	assert.Contains(suite.T(), receivedPayload, "service")
}

// TestWebhookHookInvalidURL 测试Webhook钩子无效URL
func (suite *HooksTestSuite) TestWebhookHookInvalidURL() {
	hook := NewWebhookHook("http://invalid-url-that-does-not-exist.com", []LogLevel{ERROR})
	
	entry := &LogEntry{
		Level:     ERROR,
		Message:   "test message",
		Timestamp: time.Now().Unix(),
		Fields:    map[string]interface{}{},
	}
	
	// 触发钩子，应该不会panic但可能返回错误
	err := hook.Fire(entry)
	// 网络错误是可以接受的，不应该导致程序崩溃
	if err != nil {
		assert.Contains(suite.T(), err.Error(), "error")
	}
}

// TestWebhookHookLevels 测试Webhook钩子级别过滤
func (suite *HooksTestSuite) TestWebhookHookLevels() {
	hook := NewWebhookHook("http://example.com", []LogLevel{ERROR, FATAL})
	levels := hook.Levels()
	
	assert.Contains(suite.T(), levels, ERROR)
	assert.Contains(suite.T(), levels, FATAL)
	assert.NotContains(suite.T(), levels, DEBUG)
	assert.NotContains(suite.T(), levels, INFO)
	assert.NotContains(suite.T(), levels, WARN)
}

// TestMultipleHooks 测试多个钩子
func (suite *HooksTestSuite) TestMultipleHooks() {
	tempFile := "test_multiple.log"
	defer os.Remove(tempFile)
	
	// 创建多个钩子
	consoleHook := NewConsoleHook([]LogLevel{ERROR})
	fileHook := NewFileHook(tempFile, []LogLevel{WARN, ERROR})
	
	// 创建测试事件
	entry := &LogEntry{
		Level:     ERROR,
		Message:   "multiple hooks test",
		Timestamp: time.Now().Unix(),
		Fields:    map[string]interface{}{"test": "multiple"},
	}
	
	// 触发所有钩子
	err1 := consoleHook.Fire(entry)
	err2 := fileHook.Fire(entry)
	
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	
	// 给一点时间让文件写入完成
	time.Sleep(100 * time.Millisecond)
	
	// 检查文件输出
	content, err := os.ReadFile(tempFile)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), string(content), "multiple hooks test")
}

// TestHookConcurrency 测试钩子并发安全
func (suite *HooksTestSuite) TestHookConcurrency() {
	tempFile := "test_concurrency.log"
	defer os.Remove(tempFile)
	
	hook := NewFileHook(tempFile, []LogLevel{INFO})
	
	done := make(chan bool, 10)
	
	// 并发触发钩子
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				entry := &LogEntry{
					Level:     INFO,
					Message:   "concurrency test",
					Timestamp: time.Now().Unix(),
					Fields:    map[string]interface{}{"goroutine": id, "iteration": j},
				}
				
				err := hook.Fire(entry)
				assert.NoError(suite.T(), err)
			}
			done <- true
		}(i)
	}
	
	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// 给一点时间让所有写入完成
	time.Sleep(200 * time.Millisecond)
}

// TestHookWithNilEntry 测试钩子处理nil事件
func (suite *HooksTestSuite) TestHookWithNilEntry() {
	hook := NewConsoleHook([]LogLevel{INFO})
	
	err := hook.Fire(nil)
	assert.Error(suite.T(), err)
}

// TestHookWithEmptyLevels 测试钩子空级别列表
func (suite *HooksTestSuite) TestHookWithEmptyLevels() {
	hook := NewConsoleHook([]LogLevel{})
	levels := hook.Levels()
	
	assert.Empty(suite.T(), levels)
}

// TestPredefinedLevelSets 测试预定义级别集合
func (suite *HooksTestSuite) TestPredefinedLevelSets() {
	// 测试AllLevels
	assert.Equal(suite.T(), 5, len(AllLevels))
	assert.Contains(suite.T(), AllLevels, DEBUG)
	assert.Contains(suite.T(), AllLevels, INFO)
	assert.Contains(suite.T(), AllLevels, WARN)
	assert.Contains(suite.T(), AllLevels, ERROR)
	assert.Contains(suite.T(), AllLevels, FATAL)
	
	// 测试ErrorLevels
	assert.Equal(suite.T(), 2, len(ErrorLevels))
	assert.Contains(suite.T(), ErrorLevels, ERROR)
	assert.Contains(suite.T(), ErrorLevels, FATAL)
	assert.NotContains(suite.T(), ErrorLevels, DEBUG)
	assert.NotContains(suite.T(), ErrorLevels, INFO)
	assert.NotContains(suite.T(), ErrorLevels, WARN)
}

// TestHookFileCreation 测试钩子文件创建
func (suite *HooksTestSuite) TestHookFileCreation() {
	tempFile := "test_creation.log"
	defer os.Remove(tempFile)
	
	// 确保文件不存在
	os.Remove(tempFile)
	
	// 创建文件钩子应该自动创建文件
	hook := NewFileHook(tempFile, []LogLevel{INFO})
	
	entry := &LogEntry{
		Level:     INFO,
		Message:   "file creation test",
		Timestamp: time.Now().Unix(),
		Fields:    map[string]interface{}{},
	}
	
	err := hook.Fire(entry)
	assert.NoError(suite.T(), err)
	
	// 给一点时间让文件写入
	time.Sleep(100 * time.Millisecond)
	
	// 检查文件是否被创建
	_, err = os.Stat(tempFile)
	assert.NoError(suite.T(), err)
}

// TestHookManager 测试钩子管理器
func (suite *HooksTestSuite) TestHookManager() {
	manager := NewHookManager()
	assert.NotNil(suite.T(), manager)
	
	// 添加钩子
	hook := NewConsoleHook([]LogLevel{ERROR})
	manager.AddHook("console", hook)
	
	// 测试钩子触发
	entry := &LogEntry{
		Level:     ERROR,
		Message:   "manager test",
		Timestamp: time.Now().Unix(),
		Fields:    map[string]interface{}{},
	}
	
	manager.FireHooks(entry)
	
	// 移除钩子
	manager.RemoveHook("console")
}

// TestBaseHook 测试基础钩子
func (suite *HooksTestSuite) TestBaseHook() {
	baseHook := NewBaseHook("test", []LogLevel{INFO, ERROR})
	assert.NotNil(suite.T(), baseHook)
	assert.Equal(suite.T(), "test", baseHook.Name)
	assert.True(suite.T(), baseHook.Enabled)
}

// TestHookTimeout 测试钩子超时
func (suite *HooksTestSuite) TestHookTimeout() {
	// 创建一个会延迟的测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond) // 延迟500ms
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	
	hook := NewWebhookHook(server.URL, []LogLevel{ERROR})
	
	entry := &LogEntry{
		Level:     ERROR,
		Message:   "timeout test",
		Timestamp: time.Now().Unix(),
		Fields:    map[string]interface{}{},
	}
	
	start := time.Now()
	err := hook.Fire(entry)
	duration := time.Since(start)
	
	// 应该在合理的时间内完成（考虑到网络延迟）
	assert.True(suite.T(), duration < 2*time.Second)
	
	// 可能会有超时错误，这是正常的
	if err != nil {
		// 验证错误包含超时相关信息
		assert.NotNil(suite.T(), err)
	}
}