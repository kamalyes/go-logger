/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 13:22:51
 * @FilePath: \go-logger\metrics\memory.go
 * @Description: 内存监控和管理模块
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package metrics

import (
	"fmt"
	"math"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// MemoryMonitor 内存监控器接口
type MemoryMonitor interface {
	// 开始监控
	Start() error
	Stop() error
	
	// 获取内存信息
	GetMemoryInfo() *MemoryInfo
	GetMemoryStats() *MemoryStats
	GetGCInfo() *GCInfo
	GetHeapInfo() *HeapInfo
	
	// 内存管理
	ForceGC()
	SetGCPercent(percent int) int
	SetMaxMemory(maxBytes uint64)
	
	// 内存警报
	SetMemoryThreshold(threshold float64)
	OnMemoryThresholdExceeded(callback func(info *MemoryInfo))
	
	// 内存分析
	TakeHeapSnapshot() (*HeapSnapshot, error)
	AnalyzeMemoryLeaks() *MemoryLeakReport
	
	// 清理和优化
	Cleanup()
	Optimize()
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Timestamp     time.Time `json:"timestamp"`
	
	// 总体内存
	TotalMemory   uint64    `json:"total_memory"`   // 总内存
	UsedMemory    uint64    `json:"used_memory"`    // 已使用内存
	FreeMemory    uint64    `json:"free_memory"`    // 空闲内存
	MemoryUsage   float64   `json:"memory_usage"`   // 内存使用率 (%)
	
	// 进程内存
	ProcessRSS    uint64    `json:"process_rss"`    // 进程物理内存
	ProcessVSS    uint64    `json:"process_vss"`    // 进程虚拟内存
	ProcessHeap   uint64    `json:"process_heap"`   // 进程堆内存
	ProcessStack  uint64    `json:"process_stack"`  // 进程栈内存
	
	// Go运行时内存
	GoAllocated   uint64    `json:"go_allocated"`   // Go已分配内存
	GoSys         uint64    `json:"go_sys"`         // Go系统内存
	GoHeap        uint64    `json:"go_heap"`        // Go堆内存
	GoStack       uint64    `json:"go_stack"`       // Go栈内存
	GoMSpan       uint64    `json:"go_mspan"`       // Go MSpan内存
	GoMCache      uint64    `json:"go_mcache"`      // Go MCache内存
	GoBuckHash    uint64    `json:"go_buckhash"`    // Go哈希表内存
	GoGC          uint64    `json:"go_gc"`          // Go GC内存
	GoOther       uint64    `json:"go_other"`       // Go其他内存
	
	// 内存压力
	MemoryPressure string   `json:"memory_pressure"` // low/medium/high/critical
	
	// 阈值状态
	ThresholdExceeded bool   `json:"threshold_exceeded"`
	Threshold         float64 `json:"threshold"`
}

// GCInfo GC信息
type GCInfo struct {
	Timestamp      time.Time     `json:"timestamp"`
	
	// GC统计
	NumGC          uint32        `json:"num_gc"`           // GC次数
	NumForcedGC    uint32        `json:"num_forced_gc"`    // 强制GC次数
	GCCPUFraction  float64       `json:"gc_cpu_fraction"`  // GC CPU占用比例
	
	// GC时间
	TotalPauseNs   uint64        `json:"total_pause_ns"`   // 总暂停时间(纳秒)
	TotalPause     time.Duration `json:"total_pause"`      // 总暂停时间
	LastPause      time.Duration `json:"last_pause"`       // 最后暂停时间
	AvgPause       time.Duration `json:"avg_pause"`        // 平均暂停时间
	MinPause       time.Duration `json:"min_pause"`        // 最小暂停时间
	MaxPause       time.Duration `json:"max_pause"`        // 最大暂停时间
	
	// GC触发
	NextGC         uint64        `json:"next_gc"`          // 下次GC阈值
	LastGC         time.Time     `json:"last_gc"`          // 最后GC时间
	GCTrigger      string        `json:"gc_trigger"`       // GC触发原因
	
	// GC详细统计
	PauseHistory   []time.Duration `json:"pause_history"`  // 暂停历史
	GCFrequency    float64         `json:"gc_frequency"`   // GC频率 (次/秒)
	
	// GC效率
	GCEfficiency   float64       `json:"gc_efficiency"`    // GC效率 (回收字节/暂停时间)
	HeapReleased   uint64        `json:"heap_released"`    // 释放的堆内存
}

// HeapInfo 堆信息
type HeapInfo struct {
	Timestamp      time.Time `json:"timestamp"`
	
	// 堆大小
	HeapAlloc      uint64    `json:"heap_alloc"`       // 堆已分配
	HeapSys        uint64    `json:"heap_sys"`         // 堆系统内存
	HeapIdle       uint64    `json:"heap_idle"`        // 堆空闲内存
	HeapInuse      uint64    `json:"heap_inuse"`       // 堆使用中内存
	HeapReleased   uint64    `json:"heap_released"`    // 堆释放内存
	HeapObjects    uint64    `json:"heap_objects"`     // 堆对象数量
	
	// 堆统计
	HeapUsage      float64   `json:"heap_usage"`       // 堆使用率
	HeapFragmentation float64 `json:"heap_fragmentation"` // 堆碎片率
	
	// 分配统计
	TotalAlloc     uint64    `json:"total_alloc"`      // 总分配字节数
	AllocRate      float64   `json:"alloc_rate"`       // 分配速率 (字节/秒)
	Mallocs        uint64    `json:"mallocs"`          // 分配次数
	Frees          uint64    `json:"frees"`            // 释放次数
	
	// 大对象
	LargeObjects   uint64    `json:"large_objects"`    // 大对象数量
	LargeObjectBytes uint64  `json:"large_object_bytes"` // 大对象字节数
}

// HeapSnapshot 堆快照
type HeapSnapshot struct {
	Timestamp      time.Time             `json:"timestamp"`
	TotalSize      uint64                `json:"total_size"`
	ObjectCount    uint64                `json:"object_count"`
	TypeStats      map[string]*TypeStat  `json:"type_stats"`
	SizeHistogram  map[string]uint64     `json:"size_histogram"`
	AgeHistogram   map[string]uint64     `json:"age_histogram"`
}

// TypeStat 类型统计
type TypeStat struct {
	TypeName     string  `json:"type_name"`
	Count        uint64  `json:"count"`
	TotalSize    uint64  `json:"total_size"`
	AvgSize      float64 `json:"avg_size"`
	Percentage   float64 `json:"percentage"`
}

// MemoryLeakReport 内存泄漏报告
type MemoryLeakReport struct {
	Timestamp          time.Time                `json:"timestamp"`
	SuspiciousTypes    []*SuspiciousType        `json:"suspicious_types"`
	GrowthTrend        string                   `json:"growth_trend"` // stable/growing/leaking
	RecommendedActions []string                 `json:"recommended_actions"`
	MemoryGrowthRate   float64                  `json:"memory_growth_rate"` // 字节/秒
}

// SuspiciousType 可疑类型
type SuspiciousType struct {
	TypeName       string  `json:"type_name"`
	GrowthRate     float64 `json:"growth_rate"`     // 增长率
	CurrentCount   uint64  `json:"current_count"`   // 当前数量
	CurrentSize    uint64  `json:"current_size"`    // 当前大小
	SuspicionLevel string  `json:"suspicion_level"` // low/medium/high
	Description    string  `json:"description"`
}

// DefaultMemoryMonitor 默认内存监控器实现
type DefaultMemoryMonitor struct {
	// 基础字段
	running       bool
	startTime     time.Time
	stopChan      chan struct{}
	
	// 配置
	sampleInterval   time.Duration
	threshold        float64
	maxMemory        uint64
	enableGCTuning   bool
	gcPercent        int
	
	// 历史数据
	memoryHistory    []MemoryInfo
	gcHistory        []GCInfo
	heapHistory      []HeapInfo
	maxHistorySize   int
	
	// 快照数据
	snapshots        []*HeapSnapshot
	maxSnapshots     int
	
	// 回调函数
	thresholdCallback func(info *MemoryInfo)
	
	// 内存泄漏检测
	leakDetectionEnabled bool
	baselineSnapshot     *HeapSnapshot
	lastLeakCheck        time.Time
	
	mu sync.RWMutex
}

// Start 开始监控
func (mm *DefaultMemoryMonitor) Start() error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	if mm.running {
		return nil
	}
	
	mm.running = true
	mm.startTime = time.Now()
	
	// 设置GC百分比
	if mm.enableGCTuning {
		debug.SetGCPercent(mm.gcPercent)
	}
	
	// 创建基准快照
	if mm.leakDetectionEnabled {
		mm.takeBaselineSnapshot()
	}
	
	// 启动监控协程
	go mm.monitorLoop()
	
	return nil
}

// Stop 停止监控
func (mm *DefaultMemoryMonitor) Stop() error {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	if !mm.running {
		return nil
	}
	
	mm.running = false

	if mm.stopChan != nil {
    	close(mm.stopChan)
    	mm.stopChan = nil
	}
	
	return nil
}

// monitorLoop 监控循环
func (mm *DefaultMemoryMonitor) monitorLoop() {
	ticker := time.NewTicker(mm.sampleInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			mm.collectMemoryData()
		case <-mm.stopChan:
			return
		}
	}
}

// collectMemoryData 收集内存数据
func (mm *DefaultMemoryMonitor) collectMemoryData() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	// 收集内存信息
	memInfo := mm.collectMemoryInfo()
	gcInfo := mm.collectGCInfo()
	heapInfo := mm.collectHeapInfo()
	
	// 添加到历史记录
	mm.addToHistory(memInfo, gcInfo, heapInfo)
	
	// 检查阈值
	if memInfo.ThresholdExceeded && mm.thresholdCallback != nil {
		go mm.thresholdCallback(&memInfo)
	}
	
	// 内存泄漏检测
	if mm.leakDetectionEnabled && time.Since(mm.lastLeakCheck) > time.Minute*10 {
		go mm.checkMemoryLeaks()
		mm.lastLeakCheck = time.Now()
	}
}

// collectMemoryInfo 收集内存信息
func (mm *DefaultMemoryMonitor) collectMemoryInfo() MemoryInfo {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	// 估算总内存（简化版本）
	totalMemory := uint64(8 * 1024 * 1024 * 1024) // 假设8GB总内存
	usedMemory := mem.Sys
	freeMemory := totalMemory - usedMemory
	memoryUsage := float64(usedMemory) / float64(totalMemory) * 100.0
	
	// 确定内存压力
	var memoryPressure string
	switch {
	case memoryUsage < 50:
		memoryPressure = "low"
	case memoryUsage < 70:
		memoryPressure = "medium"
	case memoryUsage < 90:
		memoryPressure = "high"
	default:
		memoryPressure = "critical"
	}
	
	return MemoryInfo{
		Timestamp:         time.Now(),
		TotalMemory:       totalMemory,
		UsedMemory:        usedMemory,
		FreeMemory:        freeMemory,
		MemoryUsage:       memoryUsage,
		ProcessRSS:        mem.Sys, // 简化
		ProcessVSS:        mem.Sys, // 简化
		ProcessHeap:       mem.HeapSys,
		ProcessStack:      mem.StackSys,
		GoAllocated:       mem.Alloc,
		GoSys:            mem.Sys,
		GoHeap:           mem.HeapSys,
		GoStack:          mem.StackSys,
		GoMSpan:          mem.MSpanSys,
		GoMCache:         mem.MCacheSys,
		GoBuckHash:       mem.BuckHashSys,
		GoGC:             mem.GCSys,
		GoOther:          mem.OtherSys,
		MemoryPressure:   memoryPressure,
		ThresholdExceeded: memoryUsage > mm.threshold,
		Threshold:        mm.threshold,
	}
}

// collectGCInfo 收集GC信息
func (mm *DefaultMemoryMonitor) collectGCInfo() GCInfo {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	gcInfo := GCInfo{
		Timestamp:     time.Now(),
		NumGC:         mem.NumGC,
		NumForcedGC:   mem.NumForcedGC,
		GCCPUFraction: mem.GCCPUFraction,
		TotalPauseNs:  mem.PauseTotalNs,
		TotalPause:    time.Duration(mem.PauseTotalNs),
		NextGC:        mem.NextGC,
	}
	
	// 计算平均暂停时间
	if mem.NumGC > 0 {
		gcInfo.AvgPause = time.Duration(mem.PauseTotalNs / uint64(mem.NumGC))
	}
	
	// 收集暂停历史
	gcInfo.PauseHistory = make([]time.Duration, 0)
	minPause := time.Duration(^uint64(0) >> 1) // 最大值
	maxPause := time.Duration(0)
	
	for _, pause := range mem.PauseNs {
		if pause > 0 {
			pauseDuration := time.Duration(pause)
			gcInfo.PauseHistory = append(gcInfo.PauseHistory, pauseDuration)
			
			if pauseDuration < minPause {
				minPause = pauseDuration
			}
			if pauseDuration > maxPause {
				maxPause = pauseDuration
			}
		}
	}
	
	gcInfo.MinPause = minPause
	gcInfo.MaxPause = maxPause
	
	// 获取最后一次GC暂停时间
	if len(gcInfo.PauseHistory) > 0 {
		gcInfo.LastPause = gcInfo.PauseHistory[len(gcInfo.PauseHistory)-1]
	}
	
	// 估算最后GC时间
	if mem.LastGC > 0 {
		gcInfo.LastGC = time.Unix(0, int64(mem.LastGC))
	}
	
	// 计算GC频率
	if len(mm.gcHistory) > 0 {
		timeDiff := gcInfo.Timestamp.Sub(mm.gcHistory[0].Timestamp).Seconds()
		if timeDiff > 0 {
			gcDiff := float64(gcInfo.NumGC - mm.gcHistory[0].NumGC)
			gcInfo.GCFrequency = gcDiff / timeDiff
		}
	}
	
	// 计算GC效率
	if gcInfo.TotalPause > 0 {
		if len(mm.gcHistory) > 0 {
			heapReleased := mm.gcHistory[len(mm.gcHistory)-1].HeapReleased
			gcInfo.GCEfficiency = float64(heapReleased) / gcInfo.TotalPause.Seconds()
		}
	}
	
	return gcInfo
}

// collectHeapInfo 收集堆信息
func (mm *DefaultMemoryMonitor) collectHeapInfo() HeapInfo {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	heapInfo := HeapInfo{
		Timestamp:    time.Now(),
		HeapAlloc:    mem.HeapAlloc,
		HeapSys:      mem.HeapSys,
		HeapIdle:     mem.HeapIdle,
		HeapInuse:    mem.HeapInuse,
		HeapReleased: mem.HeapReleased,
		HeapObjects:  mem.HeapObjects,
		TotalAlloc:   mem.TotalAlloc,
		Mallocs:      mem.Mallocs,
		Frees:        mem.Frees,
	}
	
	// 计算堆使用率
	if heapInfo.HeapSys > 0 {
		heapInfo.HeapUsage = float64(heapInfo.HeapInuse) / float64(heapInfo.HeapSys) * 100.0
	}
	
	// 计算堆碎片率
	if heapInfo.HeapSys > 0 {
		heapInfo.HeapFragmentation = float64(heapInfo.HeapIdle) / float64(heapInfo.HeapSys) * 100.0
	}
	
	// 计算分配速率
	if len(mm.heapHistory) > 0 {
		lastHeap := mm.heapHistory[len(mm.heapHistory)-1]
		timeDiff := heapInfo.Timestamp.Sub(lastHeap.Timestamp).Seconds()
		if timeDiff > 0 {
			allocDiff := heapInfo.TotalAlloc - lastHeap.TotalAlloc
			heapInfo.AllocRate = float64(allocDiff) / timeDiff
		}
	}
	
	return heapInfo
}

// addToHistory 添加到历史记录
func (mm *DefaultMemoryMonitor) addToHistory(memInfo MemoryInfo, gcInfo GCInfo, heapInfo HeapInfo) {
	// 添加内存历史
	mm.memoryHistory = append(mm.memoryHistory, memInfo)
	if len(mm.memoryHistory) > mm.maxHistorySize {
		mm.memoryHistory = mm.memoryHistory[1:]
	}
	
	// 添加GC历史
	mm.gcHistory = append(mm.gcHistory, gcInfo)
	if len(mm.gcHistory) > mm.maxHistorySize {
		mm.gcHistory = mm.gcHistory[1:]
	}
	
	// 添加堆历史
	mm.heapHistory = append(mm.heapHistory, heapInfo)
	if len(mm.heapHistory) > mm.maxHistorySize {
		mm.heapHistory = mm.heapHistory[1:]
	}
}

// GetMemoryInfo 获取当前内存信息
func (mm *DefaultMemoryMonitor) GetMemoryInfo() *MemoryInfo {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.memoryHistory) > 0 {
		// 返回最新的内存信息副本
		latest := mm.memoryHistory[len(mm.memoryHistory)-1]
		return &latest
	}
	
	// 如果没有历史记录，则实时收集
	info := mm.collectMemoryInfo()
	return &info
}

// GetMemoryStats 获取内存统计
func (mm *DefaultMemoryMonitor) GetMemoryStats() *MemoryStats {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.memoryHistory) == 0 {
		return nil
	}
	
	// 计算统计信息
	var totalUsage float64
	var maxUsage float64
	var minUsage float64 = 100.0
	
	for _, info := range mm.memoryHistory {
		totalUsage += info.MemoryUsage
		if info.MemoryUsage > maxUsage {
			maxUsage = info.MemoryUsage
		}
		if info.MemoryUsage < minUsage {
			minUsage = info.MemoryUsage
		}
	}
	
	latest := mm.memoryHistory[len(mm.memoryHistory)-1]
	
	var history []float64
	for _, info := range mm.memoryHistory {
		history = append(history, info.MemoryUsage)
	}
	
	return &MemoryStats{
		Used:      latest.UsedMemory,
		Available: latest.FreeMemory,
		Total:     latest.TotalMemory,
		Usage:     latest.MemoryUsage,
		Heap:      latest.GoHeap,
		Stack:     latest.GoStack,
		RSS:       latest.ProcessRSS,
		VMS:       latest.ProcessVSS,
		Threshold: mm.threshold,
		History:   history,
	}
}

// GetGCInfo 获取GC信息
func (mm *DefaultMemoryMonitor) GetGCInfo() *GCInfo {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.gcHistory) > 0 {
		// 返回最新的GC信息副本
		latest := mm.gcHistory[len(mm.gcHistory)-1]
		return &latest
	}
	
	// 如果没有历史记录，则实时收集
	info := mm.collectGCInfo()
	return &info
}

// GetHeapInfo 获取堆信息
func (mm *DefaultMemoryMonitor) GetHeapInfo() *HeapInfo {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.heapHistory) > 0 {
		// 返回最新的堆信息副本
		latest := mm.heapHistory[len(mm.heapHistory)-1]
		return &latest
	}
	
	// 如果没有历史记录，则实时收集
	info := mm.collectHeapInfo()
	return &info
}

// ForceGC 强制执行GC
func (mm *DefaultMemoryMonitor) ForceGC() {
	runtime.GC()
}

// SetGCPercent 设置GC百分比
func (mm *DefaultMemoryMonitor) SetGCPercent(percent int) int {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	oldPercent := mm.gcPercent
	mm.gcPercent = percent
	
	if mm.enableGCTuning {
		debug.SetGCPercent(percent)
	}
	
	return oldPercent
}

// SetMaxMemory 设置最大内存限制
func (mm *DefaultMemoryMonitor) SetMaxMemory(maxBytes uint64) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	mm.maxMemory = maxBytes
	
	// 如果当前内存使用超过限制，强制GC
	if len(mm.memoryHistory) > 0 {
		latest := mm.memoryHistory[len(mm.memoryHistory)-1]
		if latest.UsedMemory > maxBytes {
			go mm.ForceGC()
		}
	}
}

// SetMemoryThreshold 设置内存阈值
func (mm *DefaultMemoryMonitor) SetMemoryThreshold(threshold float64) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	mm.threshold = threshold
}

// OnMemoryThresholdExceeded 设置内存阈值超出回调
func (mm *DefaultMemoryMonitor) OnMemoryThresholdExceeded(callback func(info *MemoryInfo)) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	mm.thresholdCallback = callback
}

// TakeHeapSnapshot 创建堆快照
func (mm *DefaultMemoryMonitor) TakeHeapSnapshot() (*HeapSnapshot, error) {
	// 强制GC以获得更准确的快照
	runtime.GC()
	
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	snapshot := &HeapSnapshot{
		Timestamp:     time.Now(),
		TotalSize:     mem.HeapAlloc,
		ObjectCount:   mem.HeapObjects,
		TypeStats:     make(map[string]*TypeStat),
		SizeHistogram: make(map[string]uint64),
		AgeHistogram:  make(map[string]uint64),
	}
	
	// 这里应该使用更详细的堆分析，但由于Go的限制，我们使用简化版本
	// 在实际应用中，可能需要使用pprof或其他工具
	
	// 单独获取锁来添加到快照列表
	mm.mu.Lock()
	// 添加到快照列表
	mm.snapshots = append(mm.snapshots, snapshot)
	if len(mm.snapshots) > mm.maxSnapshots {
		mm.snapshots = mm.snapshots[1:]
	}
	mm.mu.Unlock()
	
	return snapshot, nil
}

// takeBaselineSnapshot 创建基准快照
func (mm *DefaultMemoryMonitor) takeBaselineSnapshot() {
	// 不在这里获取锁，因为调用者已经有锁了
	go func() {
		snapshot, _ := mm.TakeHeapSnapshot()
		mm.mu.Lock()
		mm.baselineSnapshot = snapshot
		mm.mu.Unlock()
	}()
}

// checkMemoryLeaks 检查内存泄漏
func (mm *DefaultMemoryMonitor) checkMemoryLeaks() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	// 检查基础条件
	if mm.baselineSnapshot == nil || len(mm.memoryHistory) < 5 {
		return
	}
	
	// 获取当前快照（不需要锁，因为TakeHeapSnapshot会处理锁）
	mm.mu.Unlock() // 临时释放锁避免死锁
	currentSnapshot, err := mm.TakeHeapSnapshot()
	mm.mu.Lock()   // 重新获取锁
	
	if err != nil {
		return
	}
	
	// 1. 快照比较分析
	leakInfo := mm.analyzeSnapshotComparison(currentSnapshot)
	
	// 2. 历史趋势分析
	trendInfo := mm.analyzeMemoryTrends()
	
	// 3. 堆分析
	heapInfo := mm.analyzeHeapGrowth()
	
	// 4. GC效率分析
	gcInfo := mm.analyzeGCEfficiency()
	
	// 5. 综合评估并触发告警
	mm.evaluateAndAlert(leakInfo, trendInfo, heapInfo, gcInfo)
	
	// 6. 智能基准快照更新
	mm.updateBaselineIfNeeded(currentSnapshot, leakInfo)
}

// analyzeSnapshotComparison 分析快照比较
func (mm *DefaultMemoryMonitor) analyzeSnapshotComparison(current *HeapSnapshot) *leakAnalysis {
	if mm.baselineSnapshot == nil {
		return &leakAnalysis{}
	}
	
	timeDiff := current.Timestamp.Sub(mm.baselineSnapshot.Timestamp).Seconds()
	if timeDiff <= 0 {
		return &leakAnalysis{}
	}
	
	sizeGrowth := float64(current.TotalSize - mm.baselineSnapshot.TotalSize)
	objectGrowth := float64(current.ObjectCount - mm.baselineSnapshot.ObjectCount)
	
	return &leakAnalysis{
		SizeGrowthRate:   sizeGrowth / timeDiff,
		ObjectGrowthRate: objectGrowth / timeDiff,
		TimePeriod:       timeDiff,
		SeverityLevel:    mm.calculateSeverity(sizeGrowth/timeDiff, objectGrowth/timeDiff),
	}
}

// analyzeMemoryTrends 分析内存趋势
func (mm *DefaultMemoryMonitor) analyzeMemoryTrends() *trendAnalysis {
	if len(mm.memoryHistory) < 10 {
		return &trendAnalysis{}
	}
	
	// 使用线性回归分析趋势
	recentHistory := mm.memoryHistory[len(mm.memoryHistory)-min(20, len(mm.memoryHistory)):]
	
	var sumX, sumY, sumXY, sumX2 float64
	n := float64(len(recentHistory))
	
	for i, info := range recentHistory {
		x := float64(i)
		y := float64(info.UsedMemory)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	
	// 计算线性回归斜率 (内存增长率)
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	
	// 计算相关系数 (趋势稳定性)
	meanX := sumX / n
	meanY := sumY / n
	var ssXX, ssYY, ssXY float64
	
	for i, info := range recentHistory {
		x := float64(i)
		y := float64(info.UsedMemory)
		ssXX += (x - meanX) * (x - meanX)
		ssYY += (y - meanY) * (y - meanY)
		ssXY += (x - meanX) * (y - meanY)
	}
	
	correlation := ssXY / (math.Sqrt(ssXX) * math.Sqrt(ssYY))
	
	// 计算变化率
	if len(recentHistory) >= 2 {
		oldest := recentHistory[0]
		newest := recentHistory[len(recentHistory)-1]
		timeDiff := newest.Timestamp.Sub(oldest.Timestamp).Seconds()
		if timeDiff > 0 {
			slope = float64(newest.UsedMemory-oldest.UsedMemory) / timeDiff
		}
	}
	
	return &trendAnalysis{
		GrowthSlope:   slope,
		Correlation:   correlation,
		Confidence:    math.Abs(correlation),
		TrendQuality:  mm.calculateTrendQuality(correlation, slope),
	}
}

// analyzeHeapGrowth 分析堆增长
func (mm *DefaultMemoryMonitor) analyzeHeapGrowth() *heapAnalysis {
	if len(mm.heapHistory) < 5 {
		return &heapAnalysis{}
	}
	
	recentHeap := mm.heapHistory[len(mm.heapHistory)-5:]
	oldest := recentHeap[0]
	newest := recentHeap[len(recentHeap)-1]
	
	timeDiff := newest.Timestamp.Sub(oldest.Timestamp).Seconds()
	if timeDiff <= 0 {
		return &heapAnalysis{}
	}
	
	objectGrowth := newest.HeapObjects - oldest.HeapObjects
	allocGrowth := newest.Mallocs - newest.Frees
	sizeGrowth := newest.HeapAlloc - oldest.HeapAlloc
	
	return &heapAnalysis{
		ObjectGrowthRate:  float64(objectGrowth) / timeDiff,
		AllocationBalance: float64(allocGrowth),
		SizeGrowthRate:    float64(sizeGrowth) / timeDiff,
		FragmentationRate: newest.HeapFragmentation - oldest.HeapFragmentation,
	}
}

// analyzeGCEfficiency 分析GC效率
func (mm *DefaultMemoryMonitor) analyzeGCEfficiency() *gcAnalysis {
	if len(mm.gcHistory) < 3 {
		return &gcAnalysis{}
	}
	
	recent := mm.gcHistory[len(mm.gcHistory)-3:]
	
	var totalPause time.Duration
	var totalGCs uint32
	var avgEfficiency float64
	
	for _, gc := range recent {
		totalPause += gc.TotalPause
		totalGCs += gc.NumGC
		if gc.GCEfficiency > 0 {
			avgEfficiency += gc.GCEfficiency
		}
	}
	
	avgEfficiency /= float64(len(recent))
	
	return &gcAnalysis{
		AveragePause:   totalPause / time.Duration(len(recent)),
		Frequency:      recent[len(recent)-1].GCFrequency,
		Efficiency:     avgEfficiency,
		CPUFraction:    recent[len(recent)-1].GCCPUFraction,
	}
}

// evaluateAndAlert 评估并触发告警
func (mm *DefaultMemoryMonitor) evaluateAndAlert(leak *leakAnalysis, trend *trendAnalysis, heap *heapAnalysis, gc *gcAnalysis) {
	// 计算综合风险分数
	riskScore := mm.calculateRiskScore(leak, trend, heap, gc)
	
	if riskScore > 0.7 { // 高风险
		mm.triggerHighRiskAlert(leak, trend, heap, gc)
	} else if riskScore > 0.4 { // 中等风险
		mm.triggerMediumRiskAlert(leak, trend, heap, gc)
	} else if riskScore > 0.2 { // 低风险
		mm.triggerLowRiskAlert(leak, trend, heap, gc)
	}
}

// updateBaselineIfNeeded 智能更新基准快照
func (mm *DefaultMemoryMonitor) updateBaselineIfNeeded(current *HeapSnapshot, leak *leakAnalysis) {
	// 如果检测到严重泄漏，更新基准以避免重复告警
	if leak.SeverityLevel >= 3 {
		mm.baselineSnapshot = current
		return
	}
	
	// 如果基准快照太旧（超过1小时），更新基准
	if time.Since(mm.baselineSnapshot.Timestamp) > time.Hour {
		mm.baselineSnapshot = current
		return
	}
	
	// 如果内存使用相对稳定，定期更新基准
	if leak.SeverityLevel == 0 && time.Since(mm.baselineSnapshot.Timestamp) > 30*time.Minute {
		mm.baselineSnapshot = current
	}
}

// AnalyzeMemoryLeaks 分析内存泄漏
func (mm *DefaultMemoryMonitor) AnalyzeMemoryLeaks() *MemoryLeakReport {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.memoryHistory) < 5 {
		return &MemoryLeakReport{
			Timestamp:        time.Now(),
			SuspiciousTypes:  []*SuspiciousType{},
			GrowthTrend:      "stable",
			RecommendedActions: []string{"需要更多数据来分析内存泄漏"},
			MemoryGrowthRate: 0,
		}
	}
	
	// 计算内存增长趋势
	oldest := mm.memoryHistory[0]
	latest := mm.memoryHistory[len(mm.memoryHistory)-1]
	
	timeDiff := latest.Timestamp.Sub(oldest.Timestamp).Seconds()
	memoryGrowth := float64(latest.UsedMemory) - float64(oldest.UsedMemory)
	growthRate := memoryGrowth / timeDiff
	
	var growthTrend string
	var recommendedActions []string
	
	switch {
	case growthRate < 1024: // < 1KB/s
		growthTrend = "stable"
		recommendedActions = []string{"内存使用稳定，无需特殊处理"}
	case growthRate < 1024*1024: // < 1MB/s
		growthTrend = "growing"
		recommendedActions = []string{
			"监控内存使用趋势",
			"检查是否有不必要的缓存",
			"考虑优化数据结构",
		}
	default: // >= 1MB/s
		growthTrend = "leaking"
		recommendedActions = []string{
			"立即检查内存泄漏",
			"使用pprof工具进行详细分析",
			"检查goroutine是否有泄漏",
			"检查是否有循环引用",
			"考虑强制GC",
		}
	}
	
	return &MemoryLeakReport{
		Timestamp:          time.Now(),
		SuspiciousTypes:    []*SuspiciousType{}, // 简化版本不包含具体类型
		GrowthTrend:        growthTrend,
		RecommendedActions: recommendedActions,
		MemoryGrowthRate:   growthRate,
	}
}

// Cleanup 清理内存监控器
func (mm *DefaultMemoryMonitor) Cleanup() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	// 强制GC
	runtime.GC()
	
	// 清理历史数据
	mm.memoryHistory = mm.memoryHistory[:0]
	mm.gcHistory = mm.gcHistory[:0]
	mm.heapHistory = mm.heapHistory[:0]
	
	// 清理快照
	mm.snapshots = mm.snapshots[:0]
	mm.baselineSnapshot = nil
}

// Optimize 优化内存使用
func (mm *DefaultMemoryMonitor) Optimize() {
	// 强制GC
	runtime.GC()
	
	// 调整GC目标百分比
	if mm.enableGCTuning {
		memInfo := mm.GetMemoryInfo()
		if memInfo != nil && memInfo.MemoryUsage > 70 {
			// 内存使用率高，降低GC目标以更频繁回收
			debug.SetGCPercent(50)
		} else {
			// 内存使用率正常，恢复默认设置
			debug.SetGCPercent(100)
		}
	}
	
	// 释放未使用的内存给操作系统
	debug.FreeOSMemory()
}

// String 字符串表示
func (mm *DefaultMemoryMonitor) String() string {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	status := "stopped"
	if mm.running {
		status = "running"
	}
	
	return fmt.Sprintf("MemoryMonitor{Status: %s, Threshold: %.1f%%, HistorySize: %d}",
		status, mm.threshold, len(mm.memoryHistory))
}

// 辅助类型定义
type leakAnalysis struct {
	SizeGrowthRate   float64
	ObjectGrowthRate float64
	TimePeriod       float64
	SeverityLevel    int
}

type trendAnalysis struct {
	GrowthSlope  float64
	Correlation  float64
	Confidence   float64
	TrendQuality string
}

type heapAnalysis struct {
	ObjectGrowthRate  float64
	AllocationBalance float64
	SizeGrowthRate    float64
	FragmentationRate float64
}

type gcAnalysis struct {
	AveragePause time.Duration
	Frequency    float64
	Efficiency   float64
	CPUFraction  float64
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (mm *DefaultMemoryMonitor) calculateSeverity(sizeRate, objectRate float64) int {
	// 计算严重程度：0-低，1-中低，2-中，3-中高，4-高
	score := 0
	
	// 基于内存增长率评分
	if sizeRate > 10*1024*1024 { // 10MB/s
		score += 2
	} else if sizeRate > 1024*1024 { // 1MB/s
		score += 1
	}
	
	// 基于对象增长率评分
	if objectRate > 10000 { // 10k objects/s
		score += 2
	} else if objectRate > 1000 { // 1k objects/s
		score += 1
	}
	
	if score > 4 {
		score = 4
	}
	
	return score
}

func (mm *DefaultMemoryMonitor) calculateTrendQuality(correlation, slope float64) string {
	confidence := math.Abs(correlation)
	
	if confidence > 0.8 && slope > 1024*1024 { // 强相关且高增长
		return "concerning"
	} else if confidence > 0.6 && slope > 512*1024 { // 中等相关且中等增长
		return "moderate"
	} else if confidence > 0.4 {
		return "weak"
	}
	return "stable"
}

func (mm *DefaultMemoryMonitor) calculateRiskScore(leak *leakAnalysis, trend *trendAnalysis, heap *heapAnalysis, gc *gcAnalysis) float64 {
	var score float64
	
	// 泄漏分析权重：40%
	switch leak.SeverityLevel {
	case 4:
		score += 0.4
	case 3:
		score += 0.3
	case 2:
		score += 0.2
	case 1:
		score += 0.1
	}
	
	// 趋势分析权重：30%
	switch trend.TrendQuality {
	case "concerning":
		score += 0.3
	case "moderate":
		score += 0.2
	case "weak":
		score += 0.1
	}
	
	// 堆分析权重：20%
	if heap.ObjectGrowthRate > 1000 {
		score += 0.2
	} else if heap.ObjectGrowthRate > 100 {
		score += 0.1
	}
	
	// GC效率权重：10%
	if gc.CPUFraction > 0.25 { // GC占用超过25% CPU
		score += 0.1
	} else if gc.CPUFraction > 0.1 {
		score += 0.05
	}
	
	return score
}

func (mm *DefaultMemoryMonitor) triggerHighRiskAlert(leak *leakAnalysis, trend *trendAnalysis, heap *heapAnalysis, gc *gcAnalysis) {
	fmt.Printf("[CRITICAL] Memory leak detected!\n")
	fmt.Printf("  Memory growth: %.2f MB/s\n", leak.SizeGrowthRate/(1024*1024))
	fmt.Printf("  Object growth: %.2f objects/s\n", leak.ObjectGrowthRate)
	fmt.Printf("  Trend quality: %s\n", trend.TrendQuality)
	fmt.Printf("  Heap fragmentation rate: %.2f%%\n", heap.FragmentationRate)
	fmt.Printf("  GC CPU fraction: %.2f%%\n", gc.CPUFraction*100)
	
	// 执行紧急优化
	go mm.emergencyOptimization()
}

func (mm *DefaultMemoryMonitor) triggerMediumRiskAlert(leak *leakAnalysis, trend *trendAnalysis, heap *heapAnalysis, gc *gcAnalysis) {
	fmt.Printf("[WARNING] Potential memory leak detected!\n")
	fmt.Printf("  Memory growth: %.2f KB/s\n", leak.SizeGrowthRate/1024)
	fmt.Printf("  Object growth: %.2f objects/s\n", leak.ObjectGrowthRate)
	fmt.Printf("  Trend: %s (confidence: %.2f)\n", trend.TrendQuality, trend.Confidence)
	
	// 建议优化
	fmt.Printf("  Recommended: Monitor closely and consider optimization\n")
}

func (mm *DefaultMemoryMonitor) triggerLowRiskAlert(leak *leakAnalysis, trend *trendAnalysis, heap *heapAnalysis, gc *gcAnalysis) {
	fmt.Printf("[INFO] Minor memory growth detected\n")
	fmt.Printf("  Growth rate: %.2f KB/s\n", leak.SizeGrowthRate/1024)
	fmt.Printf("  Trend: %s\n", trend.TrendQuality)
}

func (mm *DefaultMemoryMonitor) emergencyOptimization() {
	// 立即执行GC
	runtime.GC()
	runtime.GC() // 双重GC确保清理
	
	// 释放OS内存
	debug.FreeOSMemory()
	
	// 调整GC参数为更激进的设置
	debug.SetGCPercent(25) // 降低GC触发阈值
	
	// 清理监控器自身的历史数据
	mm.mu.Lock()
	// 保留最近的数据，清理旧数据
	if len(mm.memoryHistory) > 20 {
		mm.memoryHistory = mm.memoryHistory[len(mm.memoryHistory)-20:]
	}
	if len(mm.gcHistory) > 20 {
		mm.gcHistory = mm.gcHistory[len(mm.gcHistory)-20:]
	}
	if len(mm.heapHistory) > 20 {
		mm.heapHistory = mm.heapHistory[len(mm.heapHistory)-20:]
	}
	if len(mm.snapshots) > 5 {
		mm.snapshots = mm.snapshots[len(mm.snapshots)-5:]
	}
	mm.mu.Unlock()
	
	fmt.Printf("[EMERGENCY] Memory optimization completed\n")
}