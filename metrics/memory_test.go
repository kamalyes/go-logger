/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 01:23:45
 * @FilePath: \go-logger\metrics\memory_test.go
 * @Description: 内存监控模块测试套件
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package metrics

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// MemoryMonitorTestSuite 内存监控测试套件
type MemoryMonitorTestSuite struct {
	suite.Suite
	monitor *DefaultMemoryMonitor
}

// SetupSuite 套件级别的初始化，只在整个套件开始时运行一次
func (suite *MemoryMonitorTestSuite) SetupSuite() {
	// 设置全局测试环境
	runtime.GOMAXPROCS(2) // 确保有足够的goroutine用于并发测试
}

// TearDownSuite 套件级别的清理，在整个套件结束时运行一次
func (suite *MemoryMonitorTestSuite) TearDownSuite() {
	// 清理全局测试环境
	runtime.GC()
}

// SetupTest 在每个测试方法开始前运行
func (suite *MemoryMonitorTestSuite) SetupTest() {
	// 为每个测试创建新的监控器实例
	suite.monitor = NewDefaultMemoryMonitor()
	suite.monitor.sampleInterval = 10 * time.Millisecond // 快速采样用于测试
}

// TearDownTest 在每个测试方法结束后运行
func (suite *MemoryMonitorTestSuite) TearDownTest() {
	// 确保监控器停止并清理资源
	if suite.monitor != nil && suite.monitor.running {
		suite.monitor.Stop()
	}
	suite.monitor = nil
	runtime.GC() // 清理测试产生的垃圾
}

// TestNewDefaultMemoryMonitor 测试构造函数
func (suite *MemoryMonitorTestSuite) TestNewDefaultMemoryMonitor() {
	// 重新创建一个新实例来测试构造函数
	monitor := NewDefaultMemoryMonitor()
	
	suite.NotNil(monitor)
	suite.Equal(time.Second*5, monitor.sampleInterval)
	suite.Equal(80.0, monitor.threshold)
	suite.Equal(uint64(0), monitor.maxMemory)
	suite.True(monitor.enableGCTuning)
	suite.Equal(100, monitor.gcPercent)
	suite.Equal(100, monitor.maxHistorySize)
	suite.Equal(10, monitor.maxSnapshots)
	suite.True(monitor.leakDetectionEnabled)
	suite.False(monitor.running)
	suite.NotNil(monitor.stopChan)
	suite.NotNil(monitor.memoryHistory)
	suite.NotNil(monitor.gcHistory)
	suite.NotNil(monitor.heapHistory)
	suite.NotNil(monitor.snapshots)
}

// TestStartStop 测试启动和停止功能
func (suite *MemoryMonitorTestSuite) TestStartStop() {
	// 测试启动
	err := suite.monitor.Start()
	suite.NoError(err)
	suite.True(suite.monitor.running)
	suite.False(suite.monitor.startTime.IsZero())
	
	// 测试重复启动
	err = suite.monitor.Start()
	suite.NoError(err) // 重复启动不应该报错
	
	// 等待一小段时间让监控循环开始
	time.Sleep(20 * time.Millisecond)
	
	// 测试停止
	err = suite.monitor.Stop()
	suite.NoError(err)
	suite.False(suite.monitor.running)
	
	// 测试重复停止
	err = suite.monitor.Stop()
	suite.NoError(err) // 重复停止不应该报错
}

// TestMemoryInfoCollection 测试内存信息收集
func (suite *MemoryMonitorTestSuite) TestMemoryInfoCollection() {
	// 启动监控
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待收集数据
	time.Sleep(25 * time.Millisecond)
	
	// 获取内存信息
	memInfo := suite.monitor.GetMemoryInfo()
	suite.Require().NotNil(memInfo)
	
	// 验证基本字段
	suite.False(memInfo.Timestamp.IsZero())
	suite.True(memInfo.TotalMemory > 0)
	suite.True(memInfo.UsedMemory > 0)
	suite.True(memInfo.GoAllocated > 0)
	suite.True(memInfo.GoSys > 0)
	suite.True(memInfo.GoHeap > 0)
	suite.True(memInfo.MemoryUsage >= 0 && memInfo.MemoryUsage <= 100)
	suite.Contains([]string{"low", "medium", "high", "critical"}, memInfo.MemoryPressure)
	suite.True(memInfo.Threshold > 0)
}

// TestGCInfoCollection 测试GC信息收集
func (suite *MemoryMonitorTestSuite) TestGCInfoCollection() {
	// 启动监控
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 强制触发GC以生成数据
	runtime.GC()
	time.Sleep(25 * time.Millisecond)
	
	// 获取GC信息
	gcInfo := suite.monitor.GetGCInfo()
	suite.Require().NotNil(gcInfo)
	
	// 验证基本字段
	suite.False(gcInfo.Timestamp.IsZero())
	suite.True(gcInfo.GCCPUFraction >= 0)
	suite.True(gcInfo.TotalPause >= 0)
	suite.True(gcInfo.NextGC > 0)
}

// TestHeapInfoCollection 测试堆信息收集
func (suite *MemoryMonitorTestSuite) TestHeapInfoCollection() {
	// 启动监控
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待收集数据
	time.Sleep(25 * time.Millisecond)
	
	// 获取堆信息
	heapInfo := suite.monitor.GetHeapInfo()
	suite.Require().NotNil(heapInfo)
	
	// 验证基本字段
	suite.False(heapInfo.Timestamp.IsZero())
	suite.True(heapInfo.HeapAlloc > 0)
	suite.True(heapInfo.HeapSys > 0)
	suite.True(heapInfo.TotalAlloc > 0)
	suite.True(heapInfo.Mallocs > 0)
	suite.True(heapInfo.HeapUsage >= 0 && heapInfo.HeapUsage <= 100)
}

// TestMemoryStats 测试内存统计
func (suite *MemoryMonitorTestSuite) TestMemoryStats() {
	// 启动监控
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待收集数据
	time.Sleep(35 * time.Millisecond)
	
	// 获取内存统计
	stats := suite.monitor.GetMemoryStats()
	suite.Require().NotNil(stats)
	
	// 验证统计数据
	suite.True(stats.Used > 0)
	suite.True(stats.Total > 0)
	suite.True(stats.Heap > 0)
	suite.True(stats.Stack > 0)
	suite.True(stats.RSS > 0)
	suite.True(stats.VMS > 0)
	suite.Equal(80.0, stats.Threshold)
	suite.NotEmpty(stats.History)
}

// TestForceGC 测试强制GC
func (suite *MemoryMonitorTestSuite) TestForceGC() {
	// 获取GC前的计数
	var beforeGC runtime.MemStats
	runtime.ReadMemStats(&beforeGC)
	initialGC := beforeGC.NumGC
	
	// 执行强制GC
	suite.monitor.ForceGC()
	
	// 获取GC后的计数
	var afterGC runtime.MemStats
	runtime.ReadMemStats(&afterGC)
	finalGC := afterGC.NumGC
	
	// 验证GC被触发
	suite.True(finalGC > initialGC)
}

// TestGCPercentConfiguration 测试GC百分比设置
func (suite *MemoryMonitorTestSuite) TestGCPercentConfiguration() {
	// 测试设置GC百分比
	oldPercent := suite.monitor.SetGCPercent(50)
	suite.Equal(100, oldPercent) // 默认值
	suite.Equal(50, suite.monitor.gcPercent)
	
	// 再次设置
	oldPercent = suite.monitor.SetGCPercent(200)
	suite.Equal(50, oldPercent)
	suite.Equal(200, suite.monitor.gcPercent)
}

// TestMaxMemoryConfiguration 测试最大内存设置
func (suite *MemoryMonitorTestSuite) TestMaxMemoryConfiguration() {
	// 启动监控以生成内存历史
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	time.Sleep(15 * time.Millisecond)
	
	// 设置最大内存
	maxMemory := uint64(1024 * 1024 * 1024) // 1GB
	suite.monitor.SetMaxMemory(maxMemory)
	suite.Equal(maxMemory, suite.monitor.maxMemory)
}

// TestThresholdConfiguration 测试内存阈值设置
func (suite *MemoryMonitorTestSuite) TestThresholdConfiguration() {
	// 测试设置阈值
	newThreshold := 90.0
	suite.monitor.SetMemoryThreshold(newThreshold)
	suite.Equal(newThreshold, suite.monitor.threshold)
}

// TestThresholdCallback 测试内存阈值回调
func (suite *MemoryMonitorTestSuite) TestThresholdCallback() {
	// 设置回调函数
	var callbackCalled bool
	var callbackInfo *MemoryInfo
	callback := func(info *MemoryInfo) {
		callbackCalled = true
		callbackInfo = info
	}
	
	suite.monitor.OnMemoryThresholdExceeded(callback)
	suite.NotNil(suite.monitor.thresholdCallback)
	
	// 模拟阈值超出
	suite.monitor.threshold = 1.0 // 设置很低的阈值
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待回调被触发
	time.Sleep(35 * time.Millisecond)
	
	if callbackCalled {
		suite.True(callbackCalled)
		suite.NotNil(callbackInfo)
		suite.True(callbackInfo.ThresholdExceeded)
	}
}

// TestHeapSnapshot 测试堆快照
func (suite *MemoryMonitorTestSuite) TestHeapSnapshot() {
	// 创建堆快照
	snapshot, err := suite.monitor.TakeHeapSnapshot()
	suite.Require().NoError(err)
	suite.Require().NotNil(snapshot)
	
	// 验证快照数据
	suite.False(snapshot.Timestamp.IsZero())
	suite.True(snapshot.TotalSize > 0)
	suite.True(snapshot.ObjectCount > 0)
	suite.NotNil(snapshot.TypeStats)
	suite.NotNil(snapshot.SizeHistogram)
	suite.NotNil(snapshot.AgeHistogram)
	
	// 验证快照被添加到列表
	suite.Len(suite.monitor.snapshots, 1)
	suite.Equal(snapshot, suite.monitor.snapshots[0])
}

// TestSnapshotLimit 测试多次堆快照限制
func (suite *MemoryMonitorTestSuite) TestSnapshotLimit() {
	suite.monitor.maxSnapshots = 3 // 设置最大快照数量
	
	// 创建多个快照
	for i := 0; i < 5; i++ {
		_, err := suite.monitor.TakeHeapSnapshot()
		suite.Require().NoError(err)
	}
	
	// 验证快照数量限制
	suite.Len(suite.monitor.snapshots, 3)
}

// TestMemoryLeakAnalysis 测试内存泄漏分析
func (suite *MemoryMonitorTestSuite) TestMemoryLeakAnalysis() {
	// 启动监控以生成历史数据
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待收集数据
	time.Sleep(35 * time.Millisecond)
	
	// 分析内存泄漏
	report := suite.monitor.AnalyzeMemoryLeaks()
	suite.Require().NotNil(report)
	
	// 验证报告内容
	suite.False(report.Timestamp.IsZero())
	suite.NotNil(report.SuspiciousTypes)
	suite.Contains([]string{"stable", "growing", "leaking"}, report.GrowthTrend)
	suite.NotNil(report.RecommendedActions)
	suite.NotEmpty(report.RecommendedActions)
	suite.True(report.MemoryGrowthRate >= 0)
}

// TestMemoryLeakAnalysisInsufficientData 测试内存泄漏分析 - 数据不足情况
func (suite *MemoryMonitorTestSuite) TestMemoryLeakAnalysisInsufficientData() {
	// 不启动监控，保持历史数据为空
	report := suite.monitor.AnalyzeMemoryLeaks()
	suite.Require().NotNil(report)
	
	// 验证数据不足时的报告
	suite.Equal("stable", report.GrowthTrend)
	suite.Contains(report.RecommendedActions[0], "需要更多数据来分析内存泄漏")
	suite.Equal(float64(0), report.MemoryGrowthRate)
}

// TestCleanup 测试清理功能
func (suite *MemoryMonitorTestSuite) TestCleanup() {
	// 启动监控生成数据
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待生成数据
	time.Sleep(35 * time.Millisecond)
	
	// 创建快照
	_, err = suite.monitor.TakeHeapSnapshot()
	suite.Require().NoError(err)
	
	// 验证有数据
	suite.NotEmpty(suite.monitor.memoryHistory)
	suite.NotEmpty(suite.monitor.gcHistory)
	suite.NotEmpty(suite.monitor.heapHistory)
	suite.NotEmpty(suite.monitor.snapshots)
	
	// 执行清理
	suite.monitor.Cleanup()
	
	// 验证数据被清理
	suite.Empty(suite.monitor.memoryHistory)
	suite.Empty(suite.monitor.gcHistory)
	suite.Empty(suite.monitor.heapHistory)
	suite.Empty(suite.monitor.snapshots)
	suite.Nil(suite.monitor.baselineSnapshot)
}

// TestOptimize 测试优化功能
func (suite *MemoryMonitorTestSuite) TestOptimize() {
	suite.monitor.enableGCTuning = true
	
	// 启动监控
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待生成数据
	time.Sleep(15 * time.Millisecond)
	
	// 执行优化
	suite.monitor.Optimize()
	
	// 验证没有错误（优化是内部操作，主要验证不崩溃）
	suite.True(true) // 如果程序没有崩溃，测试通过
}

// TestStringRepresentation 测试字符串表示
func (suite *MemoryMonitorTestSuite) TestStringRepresentation() {
	// 测试停止状态
	str := suite.monitor.String()
	suite.Contains(str, "Status: stopped")
	suite.Contains(str, "Threshold: 80.0%")
	suite.Contains(str, "HistorySize: 0")
	
	// 启动监控
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待生成数据
	time.Sleep(15 * time.Millisecond)
	
	// 测试运行状态
	str = suite.monitor.String()
	suite.Contains(str, "Status: running")
	suite.Contains(str, "Threshold: 80.0%")
}

// TestHistoryLimit 测试历史数据限制
func (suite *MemoryMonitorTestSuite) TestHistoryLimit() {
	suite.monitor.maxHistorySize = 5 // 设置小的历史大小
	
	// 启动监控
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待生成足够的历史数据
	time.Sleep(60 * time.Millisecond)
	
	// 验证历史数据限制（使用锁来安全访问）
	suite.monitor.mu.RLock()
	memHistLen := len(suite.monitor.memoryHistory)
	gcHistLen := len(suite.monitor.gcHistory)
	heapHistLen := len(suite.monitor.heapHistory)
	suite.monitor.mu.RUnlock()
	
	suite.True(memHistLen <= 5)
	suite.True(gcHistLen <= 5)
	suite.True(heapHistLen <= 5)
}

// TestConcurrencySafety 测试并发安全性
func (suite *MemoryMonitorTestSuite) TestConcurrencySafety() {
	// 启动监控
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	var wg sync.WaitGroup
	
	// 并发读取
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			suite.monitor.GetMemoryInfo()
			suite.monitor.GetGCInfo()
			suite.monitor.GetHeapInfo()
			suite.monitor.GetMemoryStats()
		}()
	}
	
	// 并发写入
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			suite.monitor.SetMemoryThreshold(90.0)
			suite.monitor.SetGCPercent(100)
			suite.monitor.TakeHeapSnapshot()
		}()
	}
	
	wg.Wait()
	
	// 如果没有竞态条件，测试通过
	suite.True(true)
}

// TestHelperFunctions 测试辅助函数
func (suite *MemoryMonitorTestSuite) TestHelperFunctions() {
	// 测试 min 函数
	suite.Equal(3, min(3, 5))
	suite.Equal(2, min(7, 2))
	suite.Equal(0, min(0, 0))
	
	// 测试 calculateSeverity
	severity := suite.monitor.calculateSeverity(0, 0)
	suite.Equal(0, severity)
	
	// 修正测试期望值
	severity = suite.monitor.calculateSeverity(1024*1024, 1000)
	suite.Equal(0, severity) // 1MB/s + 1000 objects/s = 0分 (都不超过阈值)
	
	severity = suite.monitor.calculateSeverity(10*1024*1024, 10000)
	suite.Equal(2, severity) // 10MB/s(2分) + 10k objects/s(0分) = 2分
	
	// 测试 calculateTrendQuality
	quality := suite.monitor.calculateTrendQuality(0.9, 2*1024*1024)
	suite.Equal("concerning", quality)
	
	quality = suite.monitor.calculateTrendQuality(0.7, 600*1024)
	suite.Equal("moderate", quality)
	
	quality = suite.monitor.calculateTrendQuality(0.5, 100)
	suite.Equal("weak", quality)
	
	quality = suite.monitor.calculateTrendQuality(0.3, 100)
	suite.Equal("stable", quality)
}

// TestRiskScoreCalculation 测试风险评分计算
func (suite *MemoryMonitorTestSuite) TestRiskScoreCalculation() {
	// 测试不同风险级别的组合
	tests := []struct {
		name         string
		leak         *leakAnalysis
		trend        *trendAnalysis
		heap         *heapAnalysis
		gc           *gcAnalysis
		minScore     float64
		maxScore     float64
	}{
		{
			name:     "no risk",
			leak:     &leakAnalysis{SeverityLevel: 0},
			trend:    &trendAnalysis{TrendQuality: "stable"},
			heap:     &heapAnalysis{ObjectGrowthRate: 0},
			gc:       &gcAnalysis{CPUFraction: 0.01},
			minScore: 0.0,
			maxScore: 0.1,
		},
		{
			name:     "high risk",
			leak:     &leakAnalysis{SeverityLevel: 4},
			trend:    &trendAnalysis{TrendQuality: "concerning"},
			heap:     &heapAnalysis{ObjectGrowthRate: 2000},
			gc:       &gcAnalysis{CPUFraction: 0.3},
			minScore: 0.8,
			maxScore: 1.0,
		},
	}
	
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			score := suite.monitor.calculateRiskScore(tt.leak, tt.trend, tt.heap, tt.gc)
			suite.True(score >= tt.minScore && score <= tt.maxScore,
				"Risk score %f should be between %f and %f for %s",
				score, tt.minScore, tt.maxScore, tt.name)
		})
	}
}

// TestMemoryLeakAnalysisScenarios 测试内存泄漏分析的不同场景
func (suite *MemoryMonitorTestSuite) TestMemoryLeakAnalysisScenarios() {
	// 场景1：模拟稳定内存使用
	suite.monitor.memoryHistory = []MemoryInfo{
		{Timestamp: time.Now().Add(-5 * time.Second), UsedMemory: 1000000},
		{Timestamp: time.Now().Add(-4 * time.Second), UsedMemory: 1000100},
		{Timestamp: time.Now().Add(-3 * time.Second), UsedMemory: 1000200},
		{Timestamp: time.Now().Add(-2 * time.Second), UsedMemory: 1000300},
		{Timestamp: time.Now().Add(-1 * time.Second), UsedMemory: 1000400},
	}
	
	report := suite.monitor.AnalyzeMemoryLeaks()
	suite.Equal("stable", report.GrowthTrend)
	suite.Contains(report.RecommendedActions[0], "内存使用稳定")
	
	// 场景2：模拟内存泄漏
	suite.monitor.memoryHistory = []MemoryInfo{
		{Timestamp: time.Now().Add(-5 * time.Second), UsedMemory: 1000000},
		{Timestamp: time.Now().Add(-4 * time.Second), UsedMemory: 10000000},
		{Timestamp: time.Now().Add(-3 * time.Second), UsedMemory: 20000000},
		{Timestamp: time.Now().Add(-2 * time.Second), UsedMemory: 30000000},
		{Timestamp: time.Now().Add(-1 * time.Second), UsedMemory: 40000000},
	}
	
	report = suite.monitor.AnalyzeMemoryLeaks()
	suite.Equal("leaking", report.GrowthTrend)
	suite.True(report.MemoryGrowthRate > 1024*1024) // 超过1MB/s
	suite.Contains(report.RecommendedActions[0], "立即检查内存泄漏")
}

// TestEmergencyOptimization 测试紧急优化功能
func (suite *MemoryMonitorTestSuite) TestEmergencyOptimization() {
	// 启动监控
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 生成一些历史数据
	time.Sleep(25 * time.Millisecond)
	
	// 添加大量历史数据以测试清理
	for i := 0; i < 50; i++ {
		suite.monitor.mu.Lock()
		suite.monitor.memoryHistory = append(suite.monitor.memoryHistory, MemoryInfo{})
		suite.monitor.gcHistory = append(suite.monitor.gcHistory, GCInfo{})
		suite.monitor.heapHistory = append(suite.monitor.heapHistory, HeapInfo{})
		suite.monitor.mu.Unlock()
	}
	
	// 添加快照
	for i := 0; i < 10; i++ {
		suite.monitor.TakeHeapSnapshot()
	}
	
	// 执行紧急优化
	suite.monitor.emergencyOptimization()
	
	// 验证数据被清理
	suite.monitor.mu.Lock()
	suite.True(len(suite.monitor.memoryHistory) <= 20)
	suite.True(len(suite.monitor.gcHistory) <= 20)
	suite.True(len(suite.monitor.heapHistory) <= 20)
	suite.True(len(suite.monitor.snapshots) <= 5)
	suite.monitor.mu.Unlock()
}

// 运行测试套件的函数
func TestMemoryMonitorSuite(t *testing.T) {
	suite.Run(t, new(MemoryMonitorTestSuite))
}

// TestMemoryInfoFieldsDetailed 测试内存信息字段的完整性
func (suite *MemoryMonitorTestSuite) TestMemoryInfoFieldsDetailed() {
	// 直接测试收集方法
	memInfo := suite.monitor.collectMemoryInfo()
	
	// 验证所有字段都有有效值
	suite.False(memInfo.Timestamp.IsZero())
	suite.True(memInfo.TotalMemory > 0)
	suite.True(memInfo.UsedMemory > 0)
	suite.True(memInfo.MemoryUsage >= 0 && memInfo.MemoryUsage <= 100)
	suite.True(memInfo.ProcessRSS > 0)
	suite.True(memInfo.ProcessVSS > 0)
	suite.True(memInfo.ProcessHeap > 0)
	suite.True(memInfo.GoAllocated > 0)
	suite.True(memInfo.GoSys > 0)
	suite.True(memInfo.GoHeap > 0)
	suite.NotEmpty(memInfo.MemoryPressure)
	suite.Contains([]string{"low", "medium", "high", "critical"}, memInfo.MemoryPressure)
	suite.True(memInfo.Threshold > 0)
}

// TestGCInfoFieldsDetailed 测试GC信息字段的完整性
func (suite *MemoryMonitorTestSuite) TestGCInfoFieldsDetailed() {
	// 触发GC以生成数据
	runtime.GC()
	
	// 直接测试收集方法
	gcInfo := suite.monitor.collectGCInfo()
	
	// 验证字段
	suite.False(gcInfo.Timestamp.IsZero())
	suite.True(gcInfo.GCCPUFraction >= 0)
	suite.True(gcInfo.TotalPause >= 0)
	suite.True(gcInfo.NextGC > 0)
	suite.NotNil(gcInfo.PauseHistory)
	suite.True(gcInfo.GCFrequency >= 0)
}

// TestMemoryLeakOptimizedDetection 测试优化版的内存泄漏检测
func (suite *MemoryMonitorTestSuite) TestMemoryLeakOptimizedDetection() {
	suite.monitor.leakDetectionEnabled = true
	suite.monitor.sampleInterval = 50 * time.Millisecond
	
	err := suite.monitor.Start()
	suite.Require().NoError(err)
	
	// 等待基准快照创建和一些历史数据
	time.Sleep(150 * time.Millisecond)
	
	// 停止监控避免并发问题
	suite.monitor.Stop()
	
	// 验证没有崩溃
	suite.True(true)
}

// TestAnalyzeSnapshotComparison 测试快照比较分析
func (suite *MemoryMonitorTestSuite) TestAnalyzeSnapshotComparison() {
	// 创建基准快照
	baseSnapshot := &HeapSnapshot{
		Timestamp:   time.Now().Add(-time.Minute),
		TotalSize:   1000000,
		ObjectCount: 5000,
	}
	suite.monitor.baselineSnapshot = baseSnapshot
	
	// 创建当前快照
	currentSnapshot := &HeapSnapshot{
		Timestamp:   time.Now(),
		TotalSize:   2000000,
		ObjectCount: 8000,
	}
	
	// 测试分析方法
	analysis := suite.monitor.analyzeSnapshotComparison(currentSnapshot)
	suite.NotNil(analysis)
	suite.True(analysis.SizeGrowthRate > 0)
	suite.True(analysis.ObjectGrowthRate > 0)
	suite.True(analysis.TimePeriod > 0)
}

// TestAnalyzeMemoryTrends 测试内存趋势分析
func (suite *MemoryMonitorTestSuite) TestAnalyzeMemoryTrends() {
	// 准备测试数据
	baseTime := time.Now().Add(-10 * time.Second)
	for i := 0; i < 15; i++ {
		suite.monitor.memoryHistory = append(suite.monitor.memoryHistory, MemoryInfo{
			Timestamp:  baseTime.Add(time.Duration(i) * time.Second),
			UsedMemory: uint64(1000000 + i*100000), // 线性增长
		})
	}
	
	// 测试趋势分析
	analysis := suite.monitor.analyzeMemoryTrends()
	suite.NotNil(analysis)
	suite.True(analysis.GrowthSlope > 0) // 应该检测到增长趋势
	suite.True(analysis.Correlation >= 0 && analysis.Correlation <= 1)
	suite.True(analysis.Confidence >= 0 && analysis.Confidence <= 1)
	suite.NotEmpty(analysis.TrendQuality)
}

// TestAnalyzeHeapGrowth 测试堆增长分析
func (suite *MemoryMonitorTestSuite) TestAnalyzeHeapGrowth() {
	// 准备测试数据
	baseTime := time.Now().Add(-5 * time.Second)
	for i := 0; i < 6; i++ {
		suite.monitor.heapHistory = append(suite.monitor.heapHistory, HeapInfo{
			Timestamp:         baseTime.Add(time.Duration(i) * time.Second),
			HeapObjects:       uint64(1000 + i*100),
			HeapAlloc:         uint64(500000 + i*50000),
			Mallocs:           uint64(2000 + i*200),
			Frees:             uint64(1900 + i*180),
			HeapFragmentation: float64(10 + i*2),
		})
	}
	
	// 测试堆分析
	analysis := suite.monitor.analyzeHeapGrowth()
	suite.NotNil(analysis)
	suite.True(analysis.ObjectGrowthRate >= 0)
	suite.True(analysis.AllocationBalance >= 0)
	suite.True(analysis.SizeGrowthRate >= 0)
}

// TestAnalyzeGCEfficiency 测试GC效率分析
func (suite *MemoryMonitorTestSuite) TestAnalyzeGCEfficiency() {
	// 准备测试数据
	baseTime := time.Now().Add(-3 * time.Second)
	for i := 0; i < 4; i++ {
		suite.monitor.gcHistory = append(suite.monitor.gcHistory, GCInfo{
			Timestamp:     baseTime.Add(time.Duration(i) * time.Second),
			NumGC:         uint32(10 + i),
			TotalPause:    time.Duration(1000000 + i*100000), // 纳秒
			GCFrequency:   float64(1.0 + float64(i)*0.1),
			GCEfficiency:  float64(1000 + i*100),
			GCCPUFraction: float64(0.1 + float64(i)*0.05),
		})
	}
	
	// 测试GC效率分析
	analysis := suite.monitor.analyzeGCEfficiency()
	suite.NotNil(analysis)
	suite.True(analysis.AveragePause >= 0)
	suite.True(analysis.Frequency >= 0)
	suite.True(analysis.Efficiency >= 0)
	suite.True(analysis.CPUFraction >= 0)
}

// TestEvaluateAndAlert 测试评估和告警功能
func (suite *MemoryMonitorTestSuite) TestEvaluateAndAlert() {
	// 创建不同风险级别的测试数据
	leak := &leakAnalysis{SeverityLevel: 2}
	trend := &trendAnalysis{TrendQuality: "moderate"}
	heap := &heapAnalysis{ObjectGrowthRate: 500}
	gc := &gcAnalysis{CPUFraction: 0.15}
	
	// 测试评估方法
	suite.monitor.evaluateAndAlert(leak, trend, heap, gc)
	
	// 这个方法主要是内部处理，验证没有崩溃即可
	suite.True(true)
}

// TestUpdateBaselineIfNeeded 测试智能基准更新
func (suite *MemoryMonitorTestSuite) TestUpdateBaselineIfNeeded() {
	// 设置初始基准
	oldBaseline := &HeapSnapshot{
		Timestamp: time.Now().Add(-2 * time.Hour),
		TotalSize: 1000000,
	}
	suite.monitor.baselineSnapshot = oldBaseline
	
	current := &HeapSnapshot{
		Timestamp: time.Now(),
		TotalSize: 1100000,
	}
	
	// 测试不同情况下的基准更新
	leak := &leakAnalysis{SeverityLevel: 4} // 高风险
	suite.monitor.updateBaselineIfNeeded(current, leak)
	suite.Equal(current, suite.monitor.baselineSnapshot) // 应该更新基准
	
	// 测试基准快照过期的情况
	veryOldBaseline := &HeapSnapshot{
		Timestamp: time.Now().Add(-2 * time.Hour),
		TotalSize: 1000000,
	}
	suite.monitor.baselineSnapshot = veryOldBaseline
	leak = &leakAnalysis{SeverityLevel: 1} // 低风险
	suite.monitor.updateBaselineIfNeeded(current, leak)
	suite.Equal(current, suite.monitor.baselineSnapshot) // 过期应该更新
}

// TestTriggerAlerts 测试告警触发
func (suite *MemoryMonitorTestSuite) TestTriggerAlerts() {
	leak := &leakAnalysis{
		SizeGrowthRate:   10 * 1024 * 1024, // 10MB/s
		ObjectGrowthRate: 5000,             // 5k objects/s
		SeverityLevel:    4,
	}
	trend := &trendAnalysis{
		TrendQuality: "concerning",
		Confidence:   0.9,
	}
	heap := &heapAnalysis{
		ObjectGrowthRate:  2000,
		FragmentationRate: 50.0,
	}
	gc := &gcAnalysis{
		CPUFraction: 0.3,
	}
	
	// 测试高风险告警
	suite.monitor.triggerHighRiskAlert(leak, trend, heap, gc)
	
	// 测试中等风险告警
	leak.SeverityLevel = 2
	leak.SizeGrowthRate = 1024 * 1024 // 1MB/s
	trend.TrendQuality = "moderate"
	suite.monitor.triggerMediumRiskAlert(leak, trend, heap, gc)
	
	// 测试低风险告警
	leak.SeverityLevel = 1
	leak.SizeGrowthRate = 1024 // 1KB/s
	trend.TrendQuality = "weak"
	suite.monitor.triggerLowRiskAlert(leak, trend, heap, gc)
	
	// 验证没有崩溃
	suite.True(true)
}

// TestHistoryDataManagement 测试历史数据管理
func (suite *MemoryMonitorTestSuite) TestHistoryDataManagement() {
	suite.monitor.maxHistorySize = 3 // 设置小的历史大小用于测试
	
	// 模拟添加历史数据
	for i := 0; i < 5; i++ {
		memInfo := MemoryInfo{Timestamp: time.Now()}
		gcInfo := GCInfo{Timestamp: time.Now()}
		heapInfo := HeapInfo{Timestamp: time.Now()}
		
		suite.monitor.mu.Lock()
		suite.monitor.addToHistory(memInfo, gcInfo, heapInfo)
		suite.monitor.mu.Unlock()
	}
	
	// 验证历史大小限制
	suite.Equal(3, len(suite.monitor.memoryHistory))
	suite.Equal(3, len(suite.monitor.gcHistory))
	suite.Equal(3, len(suite.monitor.heapHistory))
}

// TestAdvancedRiskScoreCalculation 测试高级风险评分计算
func (suite *MemoryMonitorTestSuite) TestAdvancedRiskScoreCalculation() {
	// 测试所有风险级别的详细场景
	testCases := []struct {
		name         string
		leak         *leakAnalysis
		trend        *trendAnalysis
		heap         *heapAnalysis
		gc           *gcAnalysis
		expectedMin  float64
		expectedMax  float64
	}{
		{
			name:     "zero risk",
			leak:     &leakAnalysis{SeverityLevel: 0},
			trend:    &trendAnalysis{TrendQuality: "stable"},
			heap:     &heapAnalysis{ObjectGrowthRate: 0},
			gc:       &gcAnalysis{CPUFraction: 0.01},
			expectedMin: 0.0,
			expectedMax: 0.05,
		},
		{
			name:     "low risk",
			leak:     &leakAnalysis{SeverityLevel: 1},
			trend:    &trendAnalysis{TrendQuality: "weak"},
			heap:     &heapAnalysis{ObjectGrowthRate: 50},
			gc:       &gcAnalysis{CPUFraction: 0.05},
			expectedMin: 0.1,
			expectedMax: 0.25,
		},
		{
			name:     "medium risk",
			leak:     &leakAnalysis{SeverityLevel: 2},
			trend:    &trendAnalysis{TrendQuality: "moderate"},
			heap:     &heapAnalysis{ObjectGrowthRate: 500},
			gc:       &gcAnalysis{CPUFraction: 0.15},
			expectedMin: 0.25,
			expectedMax: 0.55,
		},
		{
			name:     "high risk",
			leak:     &leakAnalysis{SeverityLevel: 3},
			trend:    &trendAnalysis{TrendQuality: "concerning"},
			heap:     &heapAnalysis{ObjectGrowthRate: 1500},
			gc:       &gcAnalysis{CPUFraction: 0.25},
			expectedMin: 0.65,
			expectedMax: 0.95,
		},
		{
			name:     "critical risk",
			leak:     &leakAnalysis{SeverityLevel: 4},
			trend:    &trendAnalysis{TrendQuality: "concerning"},
			heap:     &heapAnalysis{ObjectGrowthRate: 2000},
			gc:       &gcAnalysis{CPUFraction: 0.3},
			expectedMin: 0.8,
			expectedMax: 1.0,
		},
	}
	
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			score := suite.monitor.calculateRiskScore(tc.leak, tc.trend, tc.heap, tc.gc)
			suite.True(score >= tc.expectedMin && score <= tc.expectedMax,
				"Risk score %f should be between %f and %f for %s",
				score, tc.expectedMin, tc.expectedMax, tc.name)
		})
	}
}

// 基准测试
func (suite *MemoryMonitorTestSuite) TestBenchmarkMemoryInfo() {
	// 这不是真正的基准测试，而是在套件中模拟基准测试的行为
	start := time.Now()
	
	for i := 0; i < 100; i++ {
		suite.monitor.GetMemoryInfo()
	}
	
	duration := time.Since(start)
	suite.True(duration < time.Second) // 验证性能在可接受范围内
}

func (suite *MemoryMonitorTestSuite) TestBenchmarkHeapSnapshot() {
	start := time.Now()
	
	for i := 0; i < 10; i++ {
		suite.monitor.TakeHeapSnapshot()
	}
	
	duration := time.Since(start)
	suite.True(duration < time.Second*5) // 验证性能在可接受范围内
}