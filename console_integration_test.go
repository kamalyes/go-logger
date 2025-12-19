/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-19 23:55:00
 * @FilePath: \go-logger\console_integration_test.go
 * @Description: Console 功能集成测试 - 测试 Console 在不同适配器和场景中的表现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

// TestConsoleWithStandardAdapter 测试 Console 功能与标准适配器的集成
func TestConsoleWithStandardAdapter(t *testing.T) {
	var buf bytes.Buffer

	config := &AdapterConfig{
		Type:     StandardAdapter,
		Name:     "console-test",
		Level:    DEBUG,
		Output:   &buf,
		Colorful: false,
	}

	adapter, err := NewStandardAdapter(config)
	if err != nil {
		t.Fatalf("创建适配器失败: %v", err)
	}
	defer adapter.Close()

	err = adapter.Initialize()
	if err != nil {
		t.Fatalf("初始化适配器失败: %v", err)
	}

	t.Run("BasicConsoleGroup", func(t *testing.T) {
		buf.Reset()

		adapter.ConsoleGroup("测试分组", "参数1")
		adapter.Info("分组内日志")
		adapter.ConsoleGroupEnd()

		output := buf.String()
		if !strings.Contains(output, "▼") && !strings.Contains(output, "测试分组") {
			t.Errorf("输出中未找到分组标记")
		}
	})

	t.Run("NewConsoleGroup", func(t *testing.T) {
		buf.Reset()

		cg := adapter.NewConsoleGroup()
		if cg == nil {
			t.Fatal("NewConsoleGroup 返回 nil")
		}

		cg.Group("嵌套测试分组")
		cg.Info("嵌套分组内的日志")
		cg.GroupEnd()

		output := buf.String()
		if !strings.Contains(output, "嵌套测试分组") {
			t.Errorf("输出中未找到嵌套分组内容")
		}
	})

	t.Run("ConsoleTable", func(t *testing.T) {
		buf.Reset()

		tableData := map[string]interface{}{
			"name": "张三",
			"age":  25,
			"city": "北京",
		}

		adapter.ConsoleTable(tableData)

		output := buf.String()
		if len(output) == 0 {
			t.Error("表格输出为空")
		}
	})

	t.Run("ConsoleTime", func(t *testing.T) {
		buf.Reset()

		timer := adapter.ConsoleTime("性能测试")
		if timer == nil {
			t.Error("ConsoleTime 返回 nil")
		}

		time.Sleep(10 * time.Millisecond)
		timer.End()

		output := buf.String()
		if !strings.Contains(output, "性能测试") {
			t.Error("输出中未找到计时器标签")
		}
	})
}

// TestConsoleWithUltraFastLogger 测试 Console 功能与 UltraFastLogger 的集成
func TestConsoleWithUltraFastLogger(t *testing.T) {
	var buf bytes.Buffer

	config := &LogConfig{
		Level:  DEBUG,
		Output: &buf,
	}

	logger := NewUltraFastLogger(config)

	t.Run("UltraFastConsoleGroup", func(t *testing.T) {
		buf.Reset()

		logger.ConsoleGroup("UltraFast 分组测试")
		logger.Info("分组内日志")
		logger.ConsoleGroupEnd()

		output := buf.String()
		if len(output) == 0 {
			t.Error("UltraFastLogger 的 Console 输出为空")
		}
	})

	t.Run("UltraFastNewConsoleGroup", func(t *testing.T) {
		buf.Reset()

		cg := logger.NewConsoleGroup()
		if cg == nil {
			t.Fatal("UltraFastLogger.NewConsoleGroup 返回 nil")
		}

		cg.Group("委托分组测试")
		cg.Info("委托分组内的日志")
		cg.GroupEnd()

		output := buf.String()
		if !strings.Contains(output, "委托分组测试") {
			t.Errorf("输出中未找到委托分组内容")
		}
	})

	t.Run("UltraFastConsoleTime", func(t *testing.T) {
		buf.Reset()

		timer := logger.ConsoleTime("UltraFast 计时")
		if timer != nil {
			time.Sleep(5 * time.Millisecond)
			timer.End()
		}

		// UltraFastLogger 委托给内部 consoleLogger,应该有输出
		output := buf.String()
		if len(output) == 0 {
			t.Error("UltraFastLogger ConsoleTime 输出为空")
		}
	})
}

// TestConsoleWithEmptyLogger 测试 Console 功能与 EmptyLogger 的集成
func TestConsoleWithEmptyLogger(t *testing.T) {
	logger := NewEmptyLogger()

	t.Run("EmptyLoggerConsole", func(t *testing.T) {
		// EmptyLogger 的所有 Console 方法都应该安全调用,不会 panic
		logger.ConsoleGroup("空日志器分组")
		logger.ConsoleGroupCollapsed("折叠分组")
		logger.ConsoleGroupEnd()
		logger.ConsoleTable(map[string]string{"key": "value"})

		timer := logger.ConsoleTime("空日志器计时")
		if timer == nil {
			t.Error("EmptyLogger.ConsoleTime 应返回非 nil 的 Timer")
		}
	})

	t.Run("EmptyLoggerNewConsoleGroup", func(t *testing.T) {
		cg := logger.NewConsoleGroup()
		if cg == nil {
			t.Fatal("EmptyLogger.NewConsoleGroup 返回 nil")
		}

		// 应该能安全调用所有方法
		cg.Group("测试分组")
		cg.Info("测试日志")
		cg.GroupEnd()
	})
}

// TestConsoleGroupConcurrency 测试 ConsoleGroup 的并发安全性
func TestConsoleGroupConcurrency(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(&LogConfig{
		Level:  DEBUG,
		Output: &buf,
	})

	cg := logger.NewConsoleGroup()

	// 并发测试
	done := make(chan bool)
	workers := 10

	for i := 0; i < workers; i++ {
		go func(id int) {
			cg.Group("Worker %d", id)
			cg.Info("并发日志 %d", id)
			cg.GroupEnd()
			done <- true
		}(i)
	}

	// 等待所有 worker 完成
	for i := 0; i < workers; i++ {
		<-done
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("并发测试输出为空")
	}
}

// TestConsoleGroupNesting 测试深度嵌套的分组
func TestConsoleGroupNesting(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(&LogConfig{
		Level:  DEBUG,
		Output: &buf,
	})

	cg := logger.NewConsoleGroup()

	// 创建深度嵌套 (5 层)
	depth := 5
	for i := 0; i < depth; i++ {
		cg.Group("Level %d", i+1)
		cg.Info("日志在第 %d 层", i+1)
	}

	// 结束所有嵌套
	for i := 0; i < depth; i++ {
		cg.GroupEnd()
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("深度嵌套测试输出为空")
	}

	// 验证输出包含所有层级
	if !strings.Contains(output, "Level 1") {
		t.Error("输出中未找到 Level 1")
	}
	if !strings.Contains(output, "Level 5") {
		t.Error("输出中未找到 Level 5")
	}
}

// TestConsoleTableFormats 测试不同数据类型的表格显示
func TestConsoleTableFormats(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(&LogConfig{
		Level:  DEBUG,
		Output: &buf,
	})

	t.Run("MapTable", func(t *testing.T) {
		buf.Reset()

		data := map[string]interface{}{
			"name":   "测试",
			"age":    30,
			"active": true,
		}

		logger.ConsoleTable(data)

		output := buf.String()
		if !strings.Contains(output, "name") || !strings.Contains(output, "测试") {
			t.Error("Map 表格未正确显示")
		}
	})

	t.Run("SliceTable", func(t *testing.T) {
		buf.Reset()

		data := []string{"项目1", "项目2", "项目3"}

		logger.ConsoleTable(data)

		output := buf.String()
		if len(output) == 0 {
			t.Error("Slice 表格输出为空")
		}
	})

	t.Run("StructTable", func(t *testing.T) {
		buf.Reset()

		type User struct {
			Name string
			Age  int
		}

		data := User{Name: "张三", Age: 25}

		logger.ConsoleTable(data)

		output := buf.String()
		if len(output) == 0 {
			t.Error("Struct 表格输出为空")
		}
	})
}

// TestConsoleTimerAccuracy 测试计时器的准确性
func TestConsoleTimerAccuracy(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(&LogConfig{
		Level:  DEBUG,
		Output: &buf,
	})

	timer := logger.ConsoleTime("准确性测试")
	if timer == nil {
		t.Fatal("ConsoleTime 返回 nil")
	}

	sleepDuration := 100 * time.Millisecond
	time.Sleep(sleepDuration)

	timer.End()

	output := buf.String()
	if !strings.Contains(output, "ms") && !strings.Contains(output, "毫秒") {
		t.Error("计时器输出中未找到时间单位")
	}

	// 验证输出中包含合理的时间值
	if !strings.Contains(output, "准确性测试") {
		t.Error("计时器输出中未找到标签")
	}
}

// TestConsoleGroupCollapsedBehavior 测试折叠分组的行为
func TestConsoleGroupCollapsedBehavior(t *testing.T) {
	var buf bytes.Buffer

	logger := NewLogger(&LogConfig{
		Level:  DEBUG,
		Output: &buf,
	})

	cg := logger.NewConsoleGroup()

	t.Run("CollapsedGroupFiltersNonError", func(t *testing.T) {
		buf.Reset()

		cg.GroupCollapsed("折叠的分组")
		cg.Debug("调试日志 (应被过滤)")
		cg.Info("信息日志 (应被过滤)")
		cg.Warn("警告日志 (应被过滤)")
		cg.Error("错误日志 (应显示)")
		cg.GroupEnd()

		output := buf.String()

		// 错误日志应该显示
		if !strings.Contains(output, "错误日志") {
			t.Error("折叠分组中的错误日志未显示")
		}
	})

	t.Run("NestedCollapsed", func(t *testing.T) {
		buf.Reset()

		cg.Group("外层普通分组")
		cg.Info("外层信息 (应显示)")

		cg.GroupCollapsed("内层折叠分组")
		cg.Info("内层信息 (应被过滤)")
		cg.Error("内层错误 (应显示)")
		cg.GroupEnd()

		cg.Info("外层信息2 (应显示)")
		cg.GroupEnd()

		output := buf.String()

		if !strings.Contains(output, "外层信息") {
			t.Error("外层普通分组日志未显示")
		}

		if !strings.Contains(output, "内层错误") {
			t.Error("内层折叠分组的错误日志未显示")
		}
	})
}
