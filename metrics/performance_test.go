/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-09 10:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 10:31:25
 * @FilePath: \go-logger\metrics\performance_test.go
 * @Description: 监控系统性能基准测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package metrics

import (
	"runtime"
	"testing"
	"time"
)

// BenchmarkMemoryInfoCollection 基准测试：内存信息收集性能
func BenchmarkMemoryInfoCollection(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			monitor.GetMemoryInfo()
		}
	})
}

// BenchmarkGCInfoCollection 基准测试：GC信息收集性能
func BenchmarkGCInfoCollection(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			monitor.GetGCInfo()
		}
	})
}

// BenchmarkHeapInfoCollection 基准测试：堆信息收集性能
func BenchmarkHeapInfoCollection(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			monitor.GetHeapInfo()
		}
	})
}

// BenchmarkMemoryStatsCollection 基准测试：内存统计收集性能
func BenchmarkMemoryStatsCollection(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.Start()
	defer monitor.Stop()
	
	// 等待收集一些历史数据
	time.Sleep(50 * time.Millisecond)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			monitor.GetMemoryStats()
		}
	})
}

// BenchmarkHeapSnapshot 基准测试：堆快照性能
func BenchmarkHeapSnapshot(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitor.TakeHeapSnapshot()
	}
}

// BenchmarkMemoryLeakAnalysis 基准测试：内存泄漏分析性能
func BenchmarkMemoryLeakAnalysis(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.Start()
	defer monitor.Stop()
	
	// 生成一些历史数据
	time.Sleep(100 * time.Millisecond)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitor.AnalyzeMemoryLeaks()
	}
}

// BenchmarkConcurrentMemoryAccess 基准测试：并发内存访问性能
func BenchmarkConcurrentMemoryAccess(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.Start()
	defer monitor.Stop()
	
	// 等待收集一些数据
	time.Sleep(50 * time.Millisecond)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 模拟并发访问多种监控功能
			switch b.N % 4 {
			case 0:
				monitor.GetMemoryInfo()
			case 1:
				monitor.GetGCInfo()
			case 2:
				monitor.GetHeapInfo()
			case 3:
				monitor.GetMemoryStats()
			}
		}
	})
}

// BenchmarkMonitorStartStop 基准测试：监控启停性能
func BenchmarkMonitorStartStop(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitor := NewDefaultMemoryMonitor()
		monitor.Start()
		monitor.Stop()
	}
}

// BenchmarkThresholdCallbacks 基准测试：阈值回调性能
func BenchmarkThresholdCallbacks(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.SetMemoryThreshold(1.0) // 设置很低的阈值确保触发回调
	
	callbackCalled := 0
	monitor.OnMemoryThresholdExceeded(func(info *MemoryInfo) {
		callbackCalled++
	})
	
	monitor.Start()
	defer monitor.Stop()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitor.GetMemoryInfo() // 每次调用可能触发回调
	}
}

// BenchmarkMemoryPressureCalculation 基准测试：内存压力计算性能
func BenchmarkMemoryPressureCalculation(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 直接调用collectMemoryInfo方法来测试内存信息收集性能
			monitor.collectMemoryInfo()
		}
	})
}

// BenchmarkGCTuning 基准测试：GC调优性能
func BenchmarkGCTuning(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.enableGCTuning = true
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitor.SetGCPercent(100 + i%100) // 测试不同的GC百分比设置
	}
}

// BenchmarkStatsCollectorRecord 基准测试：统计收集器记录性能
func BenchmarkStatsCollectorRecord(b *testing.B) {
	stats := NewDefaultStatsCollector()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			duration := time.Millisecond * time.Duration(b.N%100)
			stats.Record("test_operation", duration, 1024) 
		}
	})
}

// BenchmarkPerformanceMonitorRecord 基准测试：性能监控器记录性能
func BenchmarkPerformanceMonitorRecord(b *testing.B) {
	perfMonitor := NewDefaultPerformanceMonitor()
	perfMonitor.Start()
	defer perfMonitor.Stop()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			duration := time.Millisecond * time.Duration(b.N%50)
			perfMonitor.RecordLatency("test_operation", duration) 
		}
	})
}

// BenchmarkPerformanceMonitorGetData 基准测试：性能监控器数据获取性能
func BenchmarkPerformanceMonitorGetData(b *testing.B) {
	perfMonitor := NewDefaultPerformanceMonitor()
	perfMonitor.Start()
	defer perfMonitor.Stop()
	
	// 生成一些操作数据
	for i := 0; i < 100; i++ {
		perfMonitor.RecordLatency("test_operation", time.Millisecond*time.Duration(i%50))
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			perfMonitor.GetPerformanceData() 
		}
	})
}

// BenchmarkMemoryMonitorOptimize 基准测试：内存监控器优化性能
func BenchmarkMemoryMonitorOptimize(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.enableGCTuning = true
	monitor.Start()
	defer monitor.Stop()
	
	// 生成一些历史数据
	time.Sleep(50 * time.Millisecond)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitor.Optimize()
	}
}

// BenchmarkMemoryMonitorCleanup 基准测试：内存监控器清理性能
func BenchmarkMemoryMonitorCleanup(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		monitor := NewDefaultMemoryMonitor()
		monitor.Start()
		
		// 生成一些数据
		time.Sleep(10 * time.Millisecond)
		for j := 0; j < 10; j++ {
			monitor.TakeHeapSnapshot()
		}
		
		b.StartTimer()
		monitor.Cleanup()
		b.StopTimer()
		
		monitor.Stop()
		b.StartTimer()
	}
}

// BenchmarkConcurrentMonitoring 基准测试：完整并发监控性能
func BenchmarkConcurrentMonitoring(b *testing.B) {
	numMonitors := 4
	monitors := make([]*DefaultMemoryMonitor, numMonitors)
	
	// 启动多个监控器
	for i := 0; i < numMonitors; i++ {
		monitors[i] = NewDefaultMemoryMonitor()
		monitors[i].sampleInterval = 5 * time.Millisecond // 快速采样
		monitors[i].Start()
	}
	
	defer func() {
		for _, monitor := range monitors {
			monitor.Stop()
		}
	}()
	
	// 等待收集数据
	time.Sleep(25 * time.Millisecond)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		monitorIndex := 0
		for pb.Next() {
			monitor := monitors[monitorIndex%numMonitors]
			monitorIndex++
			
			// 随机执行不同的监控操作
			switch b.N % 6 {
			case 0:
				monitor.GetMemoryInfo()
			case 1:
				monitor.GetGCInfo()
			case 2:
				monitor.GetHeapInfo()
			case 3:
				monitor.GetMemoryStats()
			case 4:
				monitor.TakeHeapSnapshot()
			case 5:
				monitor.AnalyzeMemoryLeaks()
			}
		}
	})
}

// BenchmarkMemoryMonitorString 基准测试：字符串表示性能
func BenchmarkMemoryMonitorString(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.Start()
	defer monitor.Stop()
	
	// 生成一些数据
	time.Sleep(25 * time.Millisecond)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = monitor.String()
		}
	})
}

// BenchmarkHistoryManagement 基准测试：历史数据管理性能
func BenchmarkHistoryManagement(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.maxHistorySize = 1000 // 设置较大的历史大小
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memInfo := MemoryInfo{Timestamp: time.Now()}
		gcInfo := GCInfo{Timestamp: time.Now()}
		heapInfo := HeapInfo{Timestamp: time.Now()}
		
		monitor.mu.Lock()
		monitor.addToHistory(memInfo, gcInfo, heapInfo)
		monitor.mu.Unlock()
	}
}

// BenchmarkRiskScoreCalculation 基准测试：风险评分计算性能
func BenchmarkRiskScoreCalculation(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	
	leak := &leakAnalysis{SeverityLevel: 2}
	trend := &trendAnalysis{TrendQuality: "moderate"}
	heap := &heapAnalysis{ObjectGrowthRate: 500}
	gc := &gcAnalysis{CPUFraction: 0.15}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			monitor.calculateRiskScore(leak, trend, heap, gc)
		}
	})
}

// 内存分配测试辅助函数
func createMemoryLoad() {
	data := make([][]byte, 100)
	for i := range data {
		data[i] = make([]byte, 1024) // 1KB per slice
	}
	runtime.KeepAlive(data)
}

// BenchmarkMonitoringUnderLoad 基准测试：负载下的监控性能
func BenchmarkMonitoringUnderLoad(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.Start()
	defer monitor.Stop()
	
	// 在后台创建内存负载
	go func() {
		for i := 0; i < 1000; i++ {
			createMemoryLoad()
			time.Sleep(time.Microsecond)
		}
	}()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			monitor.GetMemoryInfo()
		}
	})
}

// BenchmarkLightweightMonitoring 基准测试：轻量级监控模式性能
func BenchmarkLightweightMonitoring(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	
	// 配置为轻量级模式
	monitor.leakDetectionEnabled = false
	monitor.enableGCTuning = false
	monitor.maxHistorySize = 10
	monitor.maxSnapshots = 2
	monitor.sampleInterval = 100 * time.Millisecond
	
	monitor.Start()
	defer monitor.Stop()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			monitor.GetMemoryInfo()
		}
	})
}

// BenchmarkHighFrequencyMonitoring 基准测试：高频监控性能
func BenchmarkHighFrequencyMonitoring(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	monitor.sampleInterval = time.Millisecond // 非常高频
	monitor.maxHistorySize = 5000 // 大历史缓存
	
	monitor.Start()
	defer monitor.Stop()
	
	// 等待收集大量数据
	time.Sleep(50 * time.Millisecond)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		monitor.GetMemoryStats()
	}
}

// BenchmarkEmergencyOptimization 基准测试：紧急优化性能
func BenchmarkEmergencyOptimization(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		monitor := NewDefaultMemoryMonitor()
		monitor.Start()
		
		// 生成大量数据触发优化需求
		for j := 0; j < 100; j++ {
			monitor.memoryHistory = append(monitor.memoryHistory, MemoryInfo{})
			monitor.gcHistory = append(monitor.gcHistory, GCInfo{})
			monitor.heapHistory = append(monitor.heapHistory, HeapInfo{})
		}
		for j := 0; j < 20; j++ {
			monitor.TakeHeapSnapshot()
		}
		
		b.StartTimer()
		monitor.emergencyOptimization()
		b.StopTimer()
		
		monitor.Stop()
		b.StartTimer()
	}
}

// BenchmarkMemoryInfoCaching 基准测试：内存信息缓存性能
func BenchmarkMemoryInfoCaching(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	
	// 预热缓存
	for i := 0; i < 10; i++ {
		monitor.GetMemoryInfo()
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 连续快速调用，测试缓存效果
			monitor.GetMemoryInfo()
			monitor.GetMemoryInfo()
			monitor.GetMemoryInfo()
		}
	})
}

// BenchmarkAnalyticsCalculation 基准测试：分析计算性能
func BenchmarkAnalyticsCalculation(b *testing.B) {
	monitor := NewDefaultMemoryMonitor()
	
	// 准备测试数据
	baseTime := time.Now().Add(-time.Minute)
	for i := 0; i < 100; i++ {
		monitor.memoryHistory = append(monitor.memoryHistory, MemoryInfo{
			Timestamp:  baseTime.Add(time.Duration(i) * time.Second),
			UsedMemory: uint64(1000000 + i*10000),
		})
		monitor.gcHistory = append(monitor.gcHistory, GCInfo{
			Timestamp:     baseTime.Add(time.Duration(i) * time.Second),
			NumGC:         uint32(i),
			GCCPUFraction: 0.1,
		})
		monitor.heapHistory = append(monitor.heapHistory, HeapInfo{
			Timestamp:   baseTime.Add(time.Duration(i) * time.Second),
			HeapObjects: uint64(1000 + i*10),
		})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 测试各种分析计算
		monitor.analyzeMemoryTrends()
		monitor.analyzeHeapGrowth()
		monitor.analyzeGCEfficiency()
	}
}

// BenchmarkMonitoringOverhead 基准测试：监控开销对应用性能的影响
func BenchmarkMonitoringOverhead(b *testing.B) {
	// 测试没有监控的基准性能
	b.Run("without-monitoring", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// 模拟应用工作负载
				data := make([]byte, 1024)
				_ = data
				runtime.KeepAlive(data)
			}
		})
	})
	
	// 测试有监控的性能
	b.Run("with-monitoring", func(b *testing.B) {
		monitor := NewDefaultMemoryMonitor()
		monitor.Start()
		defer monitor.Stop()
		
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// 模拟应用工作负载
				data := make([]byte, 1024)
				_ = data
				runtime.KeepAlive(data)
				
				// 偶尔触发监控
				if b.N%100 == 0 {
					monitor.GetMemoryInfo()
				}
			}
		})
	})
}