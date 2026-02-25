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

// 全局计时器管理（使用 sync.Map 优化并发性能）
var timers sync.Map

// 计时器配置
const (
	defaultTimerMaxAge     = 24 * time.Hour  // 默认最大存活时间
	defaultCleanupInterval = 5 * time.Minute // 默认清理间隔
)

var (
	timerCleanupOnce sync.Once
	timerMaxAge      = defaultTimerMaxAge
	cleanupInterval  = defaultCleanupInterval
)

// init 初始化自动清理
func init() {
	startTimerCleanup()
}

// startTimerCleanup 启动自动清理 goroutine
func startTimerCleanup() {
	timerCleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(cleanupInterval)
			defer ticker.Stop()

			for range ticker.C {
				cleanupExpiredTimers()
			}
		}()
	})
}

// cleanupExpiredTimers 清理过期的计时器（内部使用）
func cleanupExpiredTimers() int {
	count := 0
	now := time.Now()

	timers.Range(func(key, value any) bool {
		timer := value.(*Timer)
		timer.mutex.Lock()
		age := now.Sub(timer.startTime)
		timer.mutex.Unlock()

		if age > timerMaxAge {
			timers.Delete(key)
			count++
		}
		return true
	})

	return count
}

// SetTimerMaxAge 设置计时器最大存活时间（可选配置）
func SetTimerMaxAge(maxAge time.Duration) {
	if maxAge > 0 {
		timerMaxAge = maxAge
	}
}

// SetTimerCleanupInterval 设置清理间隔（可选配置）
func SetTimerCleanupInterval(interval time.Duration) {
	if interval > 0 {
		cleanupInterval = interval
	}
}

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

	// 存储到 sync.Map（优化：避免全局锁竞争）
	timers.Store(label, timer)

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

	// 从 sync.Map 中移除（优化：避免全局锁竞争）
	timers.Delete(t.label)

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

// CleanupExpiredTimers 手动清理超过指定时间未结束的计时器
func CleanupExpiredTimers(maxAge time.Duration) int {
	count := 0
	now := time.Now()

	timers.Range(func(key, value any) bool {
		timer := value.(*Timer)
		timer.mutex.Lock()
		age := now.Sub(timer.startTime)
		timer.mutex.Unlock()

		if age > maxAge {
			timers.Delete(key)
			count++
		}
		return true
	})

	return count
}

// GetActiveTimersCount 获取当前活跃的计时器数量
func GetActiveTimersCount() int {
	count := 0
	timers.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}
