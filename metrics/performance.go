/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\metrics\performance.go
 * @Description: 性能监控模块 - 监控日志系统性能
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package metrics

import (
	"runtime"
	"sync"
	"time"
)

// PerformanceMonitor 性能监控器接口
type PerformanceMonitor interface {
	// 开始监控
	Start() error
	Stop() error
	
	// 获取性能数据
	GetPerformanceData() *PerformanceData
	GetLatencyStats() *LatencyStats
	GetThroughputStats() *ThroughputStats
	GetResourceStats() *ResourceStats
	
	// 记录性能数据
	RecordLatency(operation string, latency time.Duration)
	RecordThroughput(operation string, count uint64)
	RecordResourceUsage()
	
	// 性能警报
	SetLatencyThreshold(operation string, threshold time.Duration)
	SetThroughputThreshold(operation string, threshold float64)
	SetResourceThreshold(cpu, memory float64)
	
	// 回调函数
	OnLatencyThresholdExceeded(callback func(operation string, latency time.Duration))
	OnThroughputThresholdExceeded(callback func(operation string, throughput float64))
	OnResourceThresholdExceeded(callback func(usage *ResourceUsage))
}

// PerformanceData 性能数据
type PerformanceData struct {
	Timestamp time.Time      `json:"timestamp"`
	Latency   *LatencyStats  `json:"latency"`
	Throughput *ThroughputStats `json:"throughput"`
	Resource  *ResourceStats `json:"resource"`
}

// LatencyStats 延迟统计
type LatencyStats struct {
	Operations map[string]*OperationLatency `json:"operations"`
	Overall    *LatencyMetrics              `json:"overall"`
}

// OperationLatency 操作延迟
type OperationLatency struct {
	Operation     string           `json:"operation"`
	Count         uint64           `json:"count"`
	TotalLatency  time.Duration    `json:"total_latency"`
	AvgLatency    time.Duration    `json:"avg_latency"`
	MinLatency    time.Duration    `json:"min_latency"`
	MaxLatency    time.Duration    `json:"max_latency"`
	P50Latency    time.Duration    `json:"p50_latency"`
	P90Latency    time.Duration    `json:"p90_latency"`
	P95Latency    time.Duration    `json:"p95_latency"`
	P99Latency    time.Duration    `json:"p99_latency"`
	Threshold     time.Duration    `json:"threshold"`
	ThresholdExceeded uint64       `json:"threshold_exceeded"`
	RecentLatencies []time.Duration `json:"recent_latencies"`
}

// LatencyMetrics 延迟指标
type LatencyMetrics struct {
	AvgLatency    time.Duration `json:"avg_latency"`
	MedianLatency time.Duration `json:"median_latency"`
	P90Latency    time.Duration `json:"p90_latency"`
	P95Latency    time.Duration `json:"p95_latency"`
	P99Latency    time.Duration `json:"p99_latency"`
}

// ThroughputStats 吞吐量统计
type ThroughputStats struct {
	Operations map[string]*OperationThroughput `json:"operations"`
	Overall    *ThroughputMetrics              `json:"overall"`
}

// OperationThroughput 操作吞吐量
type OperationThroughput struct {
	Operation         string    `json:"operation"`
	Count             uint64    `json:"count"`
	StartTime         time.Time `json:"start_time"`
	LastTime          time.Time `json:"last_time"`
	CurrentThroughput float64   `json:"current_throughput"` // ops/sec
	AvgThroughput     float64   `json:"avg_throughput"`     // ops/sec
	MaxThroughput     float64   `json:"max_throughput"`     // ops/sec
	Threshold         float64   `json:"threshold"`          // ops/sec
	ThresholdExceeded uint64    `json:"threshold_exceeded"`
	RecentCounts      []uint64  `json:"recent_counts"`
}

// ThroughputMetrics 吞吐量指标
type ThroughputMetrics struct {
	CurrentThroughput float64 `json:"current_throughput"` // total ops/sec
	AvgThroughput     float64 `json:"avg_throughput"`     // total ops/sec
	PeakThroughput    float64 `json:"peak_throughput"`    // total ops/sec
}

// ResourceStats 资源统计
type ResourceStats struct {
	CPU    *CPUStats    `json:"cpu"`
	Memory *MemoryStats `json:"memory"`
	GC     *GCStats     `json:"gc"`
}

// CPUStats CPU统计
type CPUStats struct {
	Usage         float64   `json:"usage"`          // CPU使用率 (%)
	UserTime      float64   `json:"user_time"`      // 用户CPU时间 (秒)
	SystemTime    float64   `json:"system_time"`    // 系统CPU时间 (秒)
	IdleTime      float64   `json:"idle_time"`      // 空闲CPU时间 (秒)
	LoadAverage   float64   `json:"load_average"`   // 负载平均值
	Cores         int       `json:"cores"`          // CPU核心数
	Threshold     float64   `json:"threshold"`      // CPU阈值 (%)
	ThresholdExceeded uint64 `json:"threshold_exceeded"`
	History       []float64 `json:"history"`        // 历史数据
}

// MemoryStats 内存统计
type MemoryStats struct {
	Used          uint64    `json:"used"`           // 已使用内存 (字节)
	Available     uint64    `json:"available"`      // 可用内存 (字节)
	Total         uint64    `json:"total"`          // 总内存 (字节)
	Usage         float64   `json:"usage"`          // 内存使用率 (%)
	Heap          uint64    `json:"heap"`           // 堆内存 (字节)
	Stack         uint64    `json:"stack"`          // 栈内存 (字节)
	RSS           uint64    `json:"rss"`            // 物理内存 (字节)
	VMS           uint64    `json:"vms"`            // 虚拟内存 (字节)
	Threshold     float64   `json:"threshold"`      // 内存阈值 (%)
	ThresholdExceeded uint64 `json:"threshold_exceeded"`
	History       []float64 `json:"history"`        // 历史数据
}

// GCStats 垃圾回收统计
type GCStats struct {
	NumGC          uint32        `json:"num_gc"`           // GC次数
	TotalPause     time.Duration `json:"total_pause"`      // 总暂停时间
	LastPause      time.Duration `json:"last_pause"`       // 最后一次暂停时间
	AvgPause       time.Duration `json:"avg_pause"`        // 平均暂停时间
	MaxPause       time.Duration `json:"max_pause"`        // 最大暂停时间
	GCCPUFraction  float64       `json:"gc_cpu_fraction"`  // GC CPU占用比例
	NextGC         uint64        `json:"next_gc"`          // 下次GC触发内存量
	LastGC         time.Time     `json:"last_gc"`          // 最后一次GC时间
	RecentPauses   []time.Duration `json:"recent_pauses"`  // 最近的暂停时间
}

// ResourceUsage 资源使用情况
type ResourceUsage struct {
	CPU    float64   `json:"cpu"`
	Memory float64   `json:"memory"`
	Time   time.Time `json:"time"`
}

// DefaultPerformanceMonitor 默认性能监控器实现
type DefaultPerformanceMonitor struct {
	// 基础字段
	running   bool
	startTime time.Time
	stopChan  chan struct{}
	
	// 延迟统计
	latencyOps map[string]*latencyOperation
	latencyThresholds map[string]time.Duration
	
	// 吞吐量统计
	throughputOps map[string]*throughputOperation
	throughputThresholds map[string]float64
	
	// 资源统计
	resourceHistory []ResourceUsage
	cpuThreshold    float64
	memoryThreshold float64
	maxHistorySize  int
	
	// 配置
	sampleInterval time.Duration
	windowSize     int
	
	// 回调函数
	latencyCallback    func(operation string, latency time.Duration)
	throughputCallback func(operation string, throughput float64)
	resourceCallback   func(usage *ResourceUsage)
	
	mu sync.RWMutex
}

// 内部操作结构
type latencyOperation struct {
	count          uint64
	totalLatency   int64 // 纳秒
	minLatency     int64
	maxLatency     int64
	recentLatencies []time.Duration
	thresholdExceeded uint64
}

type throughputOperation struct {
	count         uint64
	startTime     time.Time
	lastTime      time.Time
	recentCounts  []uint64
	thresholdExceeded uint64
}

// NewDefaultPerformanceMonitor 创建默认性能监控器
func NewDefaultPerformanceMonitor() *DefaultPerformanceMonitor {
	return &DefaultPerformanceMonitor{
		latencyOps:           make(map[string]*latencyOperation),
		latencyThresholds:    make(map[string]time.Duration),
		throughputOps:        make(map[string]*throughputOperation),
		throughputThresholds: make(map[string]float64),
		resourceHistory:      make([]ResourceUsage, 0),
		cpuThreshold:         80.0,     // 80% CPU
		memoryThreshold:      85.0,     // 85% Memory
		maxHistorySize:       100,      // 保存100个历史记录
		sampleInterval:       time.Second, // 1秒采样
		windowSize:           60,        // 60秒窗口
		stopChan:            make(chan struct{}),
	}
}

// Start 开始监控
func (pm *DefaultPerformanceMonitor) Start() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if pm.running {
		return nil
	}
	
	pm.running = true
	pm.startTime = time.Now()
	
	// 启动监控协程
	go pm.monitorLoop()
	
	return nil
}

// Stop 停止监控
func (pm *DefaultPerformanceMonitor) Stop() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	if !pm.running {
		return nil
	}
	
	pm.running = false
	close(pm.stopChan)
	
	return nil
}

// monitorLoop 监控循环
func (pm *DefaultPerformanceMonitor) monitorLoop() {
	ticker := time.NewTicker(pm.sampleInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			pm.RecordResourceUsage()
		case <-pm.stopChan:
			return
		}
	}
}

// RecordLatency 记录延迟
func (pm *DefaultPerformanceMonitor) RecordLatency(operation string, latency time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	op, ok := pm.latencyOps[operation]
	if !ok {
		op = &latencyOperation{
			recentLatencies: make([]time.Duration, 0),
			minLatency:      latency.Nanoseconds(),
			maxLatency:      latency.Nanoseconds(),
		}
		pm.latencyOps[operation] = op
	}
	
	op.count++
	op.totalLatency += latency.Nanoseconds()
	
	// 更新最小最大延迟
	if latency.Nanoseconds() < op.minLatency {
		op.minLatency = latency.Nanoseconds()
	}
	if latency.Nanoseconds() > op.maxLatency {
		op.maxLatency = latency.Nanoseconds()
	}
	
	// 记录最近延迟
	op.recentLatencies = append(op.recentLatencies, latency)
	if len(op.recentLatencies) > pm.windowSize {
		op.recentLatencies = op.recentLatencies[1:]
	}
	
	// 检查阈值
	if threshold, ok := pm.latencyThresholds[operation]; ok && latency > threshold {
		op.thresholdExceeded++
		if pm.latencyCallback != nil {
			go pm.latencyCallback(operation, latency)
		}
	}
}

// RecordThroughput 记录吞吐量
func (pm *DefaultPerformanceMonitor) RecordThroughput(operation string, count uint64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	now := time.Now()
	op, ok := pm.throughputOps[operation]
	if !ok {
		op = &throughputOperation{
			startTime:    now,
			recentCounts: make([]uint64, 0),
		}
		pm.throughputOps[operation] = op
	}
	
	op.count += count
	op.lastTime = now
	
	// 记录最近计数
	op.recentCounts = append(op.recentCounts, count)
	if len(op.recentCounts) > pm.windowSize {
		op.recentCounts = op.recentCounts[1:]
	}
	
	// 计算当前吞吐量
	if len(op.recentCounts) > 0 {
		var total uint64
		for _, c := range op.recentCounts {
			total += c
		}
		duration := time.Duration(len(op.recentCounts)) * pm.sampleInterval
		currentThroughput := float64(total) / duration.Seconds()
		
		// 检查阈值
		if threshold, ok := pm.throughputThresholds[operation]; ok && currentThroughput > threshold {
			op.thresholdExceeded++
			if pm.throughputCallback != nil {
				go pm.throughputCallback(operation, currentThroughput)
			}
		}
	}
}

// RecordResourceUsage 记录资源使用情况
func (pm *DefaultPerformanceMonitor) RecordResourceUsage() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	// 获取CPU使用率（简化版本）
	cpuUsage := pm.getCPUUsage()
	
	// 获取内存使用率
	memoryUsage := pm.getMemoryUsage(&mem)
	
	usage := ResourceUsage{
		CPU:    cpuUsage,
		Memory: memoryUsage,
		Time:   time.Now(),
	}
	
	// 添加到历史记录
	pm.resourceHistory = append(pm.resourceHistory, usage)
	if len(pm.resourceHistory) > pm.maxHistorySize {
		pm.resourceHistory = pm.resourceHistory[1:]
	}
	
	// 检查阈值
	if (cpuUsage > pm.cpuThreshold || memoryUsage > pm.memoryThreshold) && pm.resourceCallback != nil {
		go pm.resourceCallback(&usage)
	}
}

// getCPUUsage 获取CPU使用率（简化实现）
func (pm *DefaultPerformanceMonitor) getCPUUsage() float64 {
	// 这里使用简化的CPU使用率计算
	// 在实际应用中，应该使用更精确的系统调用
	return float64(runtime.NumGoroutine()) / float64(runtime.NumCPU()) * 10.0
}

// getMemoryUsage 获取内存使用率
func (pm *DefaultPerformanceMonitor) getMemoryUsage(mem *runtime.MemStats) float64 {
	// 使用已分配的堆内存计算使用率
	// 这里使用一个估算值，实际应该获取系统总内存
	const estimatedTotalMemory = 8 * 1024 * 1024 * 1024 // 8GB
	return float64(mem.Sys) / estimatedTotalMemory * 100.0
}

// SetLatencyThreshold 设置延迟阈值
func (pm *DefaultPerformanceMonitor) SetLatencyThreshold(operation string, threshold time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.latencyThresholds[operation] = threshold
}

// SetThroughputThreshold 设置吞吐量阈值
func (pm *DefaultPerformanceMonitor) SetThroughputThreshold(operation string, threshold float64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.throughputThresholds[operation] = threshold
}

// SetResourceThreshold 设置资源阈值
func (pm *DefaultPerformanceMonitor) SetResourceThreshold(cpu, memory float64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.cpuThreshold = cpu
	pm.memoryThreshold = memory
}

// OnLatencyThresholdExceeded 设置延迟阈值超出回调
func (pm *DefaultPerformanceMonitor) OnLatencyThresholdExceeded(callback func(operation string, latency time.Duration)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.latencyCallback = callback
}

// OnThroughputThresholdExceeded 设置吞吐量阈值超出回调
func (pm *DefaultPerformanceMonitor) OnThroughputThresholdExceeded(callback func(operation string, throughput float64)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.throughputCallback = callback
}

// OnResourceThresholdExceeded 设置资源阈值超出回调
func (pm *DefaultPerformanceMonitor) OnResourceThresholdExceeded(callback func(usage *ResourceUsage)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	pm.resourceCallback = callback
}

// GetPerformanceData 获取性能数据
func (pm *DefaultPerformanceMonitor) GetPerformanceData() *PerformanceData {
	return &PerformanceData{
		Timestamp:  time.Now(),
		Latency:    pm.GetLatencyStats(),
		Throughput: pm.GetThroughputStats(),
		Resource:   pm.GetResourceStats(),
	}
}

// GetLatencyStats 获取延迟统计
func (pm *DefaultPerformanceMonitor) GetLatencyStats() *LatencyStats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	stats := &LatencyStats{
		Operations: make(map[string]*OperationLatency),
	}
	
	var allLatencies []time.Duration
	var totalLatency int64
	var totalCount uint64
	
	for operation, op := range pm.latencyOps {
		opLatency := &OperationLatency{
			Operation:         operation,
			Count:             op.count,
			TotalLatency:      time.Duration(op.totalLatency),
			MinLatency:        time.Duration(op.minLatency),
			MaxLatency:        time.Duration(op.maxLatency),
			ThresholdExceeded: op.thresholdExceeded,
			RecentLatencies:   make([]time.Duration, len(op.recentLatencies)),
		}
		
		if op.count > 0 {
			opLatency.AvgLatency = time.Duration(op.totalLatency / int64(op.count))
		}
		
		if threshold, ok := pm.latencyThresholds[operation]; ok {
			opLatency.Threshold = threshold
		}
		
		// 复制最近延迟
		copy(opLatency.RecentLatencies, op.recentLatencies)
		
		// 计算百分位数
		if len(op.recentLatencies) > 0 {
			pm.calculateLatencyPercentiles(opLatency, op.recentLatencies)
		}
		
		stats.Operations[operation] = opLatency
		
		// 收集总体统计数据
		allLatencies = append(allLatencies, op.recentLatencies...)
		totalLatency += op.totalLatency
		totalCount += op.count
	}
	
	// 计算总体延迟指标
	if len(allLatencies) > 0 {
		stats.Overall = &LatencyMetrics{}
		if totalCount > 0 {
			stats.Overall.AvgLatency = time.Duration(totalLatency / int64(totalCount))
		}
		pm.calculateOverallLatencyMetrics(stats.Overall, allLatencies)
	}
	
	return stats
}

// calculateLatencyPercentiles 计算延迟百分位数
func (pm *DefaultPerformanceMonitor) calculateLatencyPercentiles(opLatency *OperationLatency, latencies []time.Duration) {
	if len(latencies) == 0 {
		return
	}
	
	// 排序延迟数据
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	
	// 简单的插入排序
	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && sorted[j] > key {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}
	
	// 计算百分位数
	opLatency.P50Latency = sorted[len(sorted)*50/100]
	opLatency.P90Latency = sorted[len(sorted)*90/100]
	opLatency.P95Latency = sorted[len(sorted)*95/100]
	opLatency.P99Latency = sorted[len(sorted)*99/100]
}

// calculateOverallLatencyMetrics 计算总体延迟指标
func (pm *DefaultPerformanceMonitor) calculateOverallLatencyMetrics(metrics *LatencyMetrics, latencies []time.Duration) {
	if len(latencies) == 0 {
		return
	}
	
	// 排序延迟数据
	sorted := make([]time.Duration, len(latencies))
	copy(sorted, latencies)
	
	// 简单的插入排序
	for i := 1; i < len(sorted); i++ {
		key := sorted[i]
		j := i - 1
		for j >= 0 && sorted[j] > key {
			sorted[j+1] = sorted[j]
			j--
		}
		sorted[j+1] = key
	}
	
	// 计算指标
	metrics.MedianLatency = sorted[len(sorted)/2]
	metrics.P90Latency = sorted[len(sorted)*90/100]
	metrics.P95Latency = sorted[len(sorted)*95/100]
	metrics.P99Latency = sorted[len(sorted)*99/100]
}

// GetThroughputStats 获取吞吐量统计
func (pm *DefaultPerformanceMonitor) GetThroughputStats() *ThroughputStats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	stats := &ThroughputStats{
		Operations: make(map[string]*OperationThroughput),
	}
	
	var totalCurrentThroughput float64
	var totalAvgThroughput float64
	var maxThroughput float64
	operationCount := 0
	
	for operation, op := range pm.throughputOps {
		opThroughput := &OperationThroughput{
			Operation:         operation,
			Count:             op.count,
			StartTime:         op.startTime,
			LastTime:          op.lastTime,
			ThresholdExceeded: op.thresholdExceeded,
			RecentCounts:      make([]uint64, len(op.recentCounts)),
		}
		
		if threshold, ok := pm.throughputThresholds[operation]; ok {
			opThroughput.Threshold = threshold
		}
		
		// 复制最近计数
		copy(opThroughput.RecentCounts, op.recentCounts)
		
		// 计算吞吐量
		if !op.startTime.IsZero() && !op.lastTime.IsZero() {
			duration := op.lastTime.Sub(op.startTime).Seconds()
			if duration > 0 {
				opThroughput.AvgThroughput = float64(op.count) / duration
			}
		}
		
		// 计算当前吞吐量
		if len(op.recentCounts) > 0 {
			var total uint64
			for _, c := range op.recentCounts {
				total += c
			}
			windowDuration := time.Duration(len(op.recentCounts)) * pm.sampleInterval
			if windowDuration.Seconds() > 0 {
				opThroughput.CurrentThroughput = float64(total) / windowDuration.Seconds()
			}
		}
		
		// 设置最大吞吐量
		if opThroughput.CurrentThroughput > opThroughput.AvgThroughput {
			opThroughput.MaxThroughput = opThroughput.CurrentThroughput
		} else {
			opThroughput.MaxThroughput = opThroughput.AvgThroughput
		}
		
		stats.Operations[operation] = opThroughput
		
		// 累计总体统计
		totalCurrentThroughput += opThroughput.CurrentThroughput
		totalAvgThroughput += opThroughput.AvgThroughput
		if opThroughput.MaxThroughput > maxThroughput {
			maxThroughput = opThroughput.MaxThroughput
		}
		operationCount++
	}
	
	// 计算总体吞吐量指标
	stats.Overall = &ThroughputMetrics{
		CurrentThroughput: totalCurrentThroughput,
		PeakThroughput:    maxThroughput,
	}
	
	if operationCount > 0 {
		stats.Overall.AvgThroughput = totalAvgThroughput / float64(operationCount)
	}
	
	return stats
}

// GetResourceStats 获取资源统计
func (pm *DefaultPerformanceMonitor) GetResourceStats() *ResourceStats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	// 获取当前资源使用情况
	currentCPU := pm.getCPUUsage()
	currentMemory := pm.getMemoryUsage(&mem)
	
	// 计算历史统计
	var cpuHistory []float64
	var memoryHistory []float64
	var cpuThresholdExceeded, memoryThresholdExceeded uint64
	
	for _, usage := range pm.resourceHistory {
		cpuHistory = append(cpuHistory, usage.CPU)
		memoryHistory = append(memoryHistory, usage.Memory)
		
		if usage.CPU > pm.cpuThreshold {
			cpuThresholdExceeded++
		}
		if usage.Memory > pm.memoryThreshold {
			memoryThresholdExceeded++
		}
	}
	
	// 构建统计信息
	stats := &ResourceStats{
		CPU: &CPUStats{
			Usage:             currentCPU,
			Cores:             runtime.NumCPU(),
			Threshold:         pm.cpuThreshold,
			ThresholdExceeded: cpuThresholdExceeded,
			History:           cpuHistory,
		},
		Memory: &MemoryStats{
			Used:              mem.Alloc,
			Total:             mem.Sys,
			Usage:             currentMemory,
			Heap:              mem.HeapAlloc,
			Stack:             mem.StackInuse,
			Threshold:         pm.memoryThreshold,
			ThresholdExceeded: memoryThresholdExceeded,
			History:           memoryHistory,
		},
		GC: &GCStats{
			NumGC:         mem.NumGC,
			TotalPause:    time.Duration(mem.PauseTotalNs),
			GCCPUFraction: mem.GCCPUFraction,
			NextGC:        mem.NextGC,
		},
	}
	
	// 计算GC统计
	if mem.NumGC > 0 {
		stats.GC.AvgPause = time.Duration(mem.PauseTotalNs / uint64(mem.NumGC))
		if len(mem.PauseNs) > 0 {
			// 获取最近的GC暂停时间
			for i, pause := range mem.PauseNs {
				if pause > 0 {
					pauseDuration := time.Duration(pause)
					stats.GC.LastPause = pauseDuration
					if pauseDuration > stats.GC.MaxPause {
						stats.GC.MaxPause = pauseDuration
					}
					
					// 只保存最近的一些暂停时间
					if len(stats.GC.RecentPauses) < 10 {
						stats.GC.RecentPauses = append(stats.GC.RecentPauses, pauseDuration)
					} else {
						// 替换最旧的记录
						stats.GC.RecentPauses[i%10] = pauseDuration
					}
				}
			}
		}
		
		// 估算最后GC时间
		stats.GC.LastGC = time.Now().Add(-time.Duration(mem.LastGC))
	}
	
	return stats
}