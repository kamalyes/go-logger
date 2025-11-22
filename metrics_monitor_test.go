/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-09 10:50:45
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 12:28:52
 * @FilePath: \go-logger\metrics_monitor_test.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"testing"
)

func TestMonitor(t *testing.T) {
	m := NewMonitor()

	// 测试基础功能
	done := m.Track()
	done(nil)

	heap := m.FastMemory()
	if heap == 0 {
		t.Error("Heap should not be zero")
	}

	healthy, pressure := m.QuickCheck()
	if !healthy {
		t.Error("Should be healthy")
	}
	if pressure == "" {
		t.Error("Pressure should not be empty")
	}

	// 测试统计
	ops, errors, current, peak, warnings, rate := m.Stats()
	if ops == 0 {
		t.Error("Should have operations")
	}

	t.Logf("Stats: ops=%d errors=%d current=%d peak=%d warnings=%d rate=%f",
		ops, errors, current, peak, warnings, rate)
	t.Logf("Monitor: %s", m.String())
}

func BenchmarkTrack(b *testing.B) {
	m := NewMonitor()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			done := m.Track()
			done(nil)
		}
	})
}

func BenchmarkFastMemory(b *testing.B) {
	m := NewMonitor()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.FastMemory()
		}
	})
}

func BenchmarkUpdate(b *testing.B) {
	m := NewMonitor()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Update(uint64(b.N * 1024))
		}
	})
}
