/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-19 23:52:45
 * @FilePath: \go-logger\timer.go
 * @Description: 计时器功能，类似 JavaScript console.time()
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Timer 计时器
type Timer struct {
	logger      ILogger
	label       string
	startTime   time.Time
	indentLevel int
	mutex       sync.Mutex
}

// 全局计时器管理
var (
	timers      = make(map[string]*Timer)
	timersMutex sync.RWMutex
)

// NewTimer 创建新的计时器
func NewTimer(logger ILogger, label string, indentLevel int) *Timer {
	timer := &Timer{
		logger:      logger,
		label:       label,
		startTime:   time.Now(),
		indentLevel: indentLevel,
	}

	// 记录开始信息
	indent := strings.Repeat("  ", indentLevel)
	logger.InfoMsg(fmt.Sprintf("%s⏱️  %s: 计时开始", indent, label))

	// 存储到全局映射
	timersMutex.Lock()
	timers[label] = timer
	timersMutex.Unlock()

	return timer
}

// End 结束计时并输出耗时
// 类似 JavaScript console.timeEnd()
func (t *Timer) End() time.Duration {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	elapsed := time.Since(t.startTime)
	indent := strings.Repeat("  ", t.indentLevel)

	// 格式化耗时
	timeStr := formatDuration(elapsed)
	t.logger.InfoMsg(fmt.Sprintf("%s⏱️  %s: %s", indent, t.label, timeStr))

	// 从全局映射中移除
	timersMutex.Lock()
	delete(timers, t.label)
	timersMutex.Unlock()

	return elapsed
}

// Log 输出当前耗时（不结束计时）
// 类似 JavaScript console.timeLog()
func (t *Timer) Log(msg string, args ...interface{}) time.Duration {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	elapsed := time.Since(t.startTime)
	indent := strings.Repeat("  ", t.indentLevel)

	timeStr := formatDuration(elapsed)
	message := fmt.Sprintf(msg, args...)

	if message != "" {
		t.logger.InfoMsg(fmt.Sprintf("%s⏱️  %s: %s - %s", indent, t.label, timeStr, message))
	} else {
		t.logger.InfoMsg(fmt.Sprintf("%s⏱️  %s: %s", indent, t.label, timeStr))
	}

	return elapsed
}

// Elapsed 获取已经过的时间（不输出日志）
func (t *Timer) Elapsed() time.Duration {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return time.Since(t.startTime)
}

// formatDuration 格式化时间间隔
func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	} else if d < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(d.Nanoseconds())/1000.0)
	} else if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000.0)
	} else if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	} else {
		return d.String()
	}
}

// ============================================================================
// 全局计时器方法
// ============================================================================

// Time 使用默认日志器创建计时器
func Time(label string) *Timer {
	return NewTimer(defaultLogger, label, 0)
}

// TimeEnd 结束指定标签的计时器
func TimeEnd(label string) time.Duration {
	timersMutex.RLock()
	timer, exists := timers[label]
	timersMutex.RUnlock()

	if !exists {
		defaultLogger.WarnMsg(fmt.Sprintf("⚠️  计时器 '%s' 不存在", label))
		return 0
	}

	return timer.End()
}

// TimeLog 输出指定标签计时器的当前耗时
func TimeLog(label string, msg string, args ...interface{}) time.Duration {
	timersMutex.RLock()
	timer, exists := timers[label]
	timersMutex.RUnlock()

	if !exists {
		defaultLogger.WarnMsg(fmt.Sprintf("⚠️  计时器 '%s' 不存在", label))
		return 0
	}

	return timer.Log(msg, args...)
}
