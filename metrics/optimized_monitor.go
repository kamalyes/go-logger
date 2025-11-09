/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-09 11:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 11:00:00
 * @FilePath: \go-logger\metrics\optimized_monitor.go
 * @Description: 极致性能优化的监控器实现
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package metrics

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// OptimizedMonitor 极致性能优化的监控器
type OptimizedMonitor struct {
	// 原子操作字段 - 避免锁竞争
	enabled       int64 // 0=禁用, 1=启用
	lastCheck     int64 // 上次检查时间 (Unix nano)
	checkInterval int64 // 检查间隔 (nanoseconds)
	
	// 缓存的内存信息 - 减少系统调用
	cachedMemInfo unsafe.Pointer // *MemoryInfo
	cacheValid    int64           // 缓存是否有效
	cacheExpiry   int64           // 缓存过期时间
	
	// 轻量级统计
	totalCalls    uint64
	totalErrors   uint64
	totalLatency  uint64 // 纳秒
	
	// 高性能环形缓冲区 - 避免动态分配
	history      []fastMemInfo
	historyIndex uint32
	historySize  uint32
	
	// 最小化的锁和同步
	mu sync.RWMutex
	
	// 配置
	config OptimizedConfig
}

// OptimizedConfig 优化配置
type OptimizedConfig struct {
	CacheExpiry     time.Duration `json:"cache_expiry"`     // 缓存过期时间
	CheckInterval   time.Duration `json:"check_interval"`   // 检查间隔
	HistorySize     int           `json:"history_size"`     // 历史大小
	EnableCaching   bool          `json:"enable_caching"`   // 启用缓存
	EnableHistory   bool          `json:"enable_history"`   // 启用历史
	LightweightMode bool          `json:"lightweight_mode"` // 轻量级模式
}

// fastMemInfo 轻量级内存信息 - 仅包含关键字段
type fastMemInfo struct {
	timestamp uint64 // Unix nano
	heap      uint64 // 堆内存
	stack     uint64 // 栈内存
	used      uint64 // 使用内存
	numGC     uint32 // GC次数
}

// NewOptimizedMonitor 创建优化监控器
func NewOptimizedMonitor(config OptimizedConfig) *OptimizedMonitor {
	if config.CacheExpiry == 0 {
		config.CacheExpiry = 100 * time.Millisecond
	}
	if config.CheckInterval == 0 {
		config.CheckInterval = 500 * time.Millisecond
	}
	if config.HistorySize == 0 {
		config.HistorySize = 64 // 小的环形缓冲区
	}
	
	return &OptimizedMonitor{
		checkInterval: int64(config.CheckInterval),
		history:       make([]fastMemInfo, config.HistorySize),
		historySize:   uint32(config.HistorySize),
		config:        config,
	}
}

// Start 启动监控 - 超轻量级
func (m *OptimizedMonitor) Start() error {
	if atomic.CompareAndSwapInt64(&m.enabled, 0, 1) {
		atomic.StoreInt64(&m.lastCheck, time.Now().UnixNano())
		return nil
	}
	return nil
}

// Stop 停止监控
func (m *OptimizedMonitor) Stop() error {
	atomic.StoreInt64(&m.enabled, 0)
	return nil
}

// FastMemoryInfo 快速内存信息获取 - 零分配版本
func (m *OptimizedMonitor) FastMemoryInfo() (heap, stack, used uint64, numGC uint32) {
	if atomic.LoadInt64(&m.enabled) == 0 {
		return 0, 0, 0, 0
	}
	
	now := time.Now().UnixNano()
	
	// 检查缓存
	if m.config.EnableCaching {
		if atomic.LoadInt64(&m.cacheValid) == 1 && 
		   now < atomic.LoadInt64(&m.cacheExpiry) {
			cached := (*MemoryInfo)(atomic.LoadPointer(&m.cachedMemInfo))
			if cached != nil {
				return cached.GoHeap, 0, cached.UsedMemory, 0
			}
		}
	}
	
	// 限制检查频率
	lastCheck := atomic.LoadInt64(&m.lastCheck)
	if now-lastCheck < atomic.LoadInt64(&m.checkInterval) {
		return 0, 0, 0, 0
	}
	
	// 原子更新检查时间
	if !atomic.CompareAndSwapInt64(&m.lastCheck, lastCheck, now) {
		return 0, 0, 0, 0 // 其他goroutine已经在检查
	}
	
	// 快速内存信息收集 - 最小化系统调用
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	
	heap = ms.HeapAlloc
	stack = ms.StackInuse  
	used = ms.Sys
	numGC = ms.NumGC
	
	// 更新环形缓冲区（如果启用）
	if m.config.EnableHistory {
		index := atomic.AddUint32(&m.historyIndex, 1) % m.historySize
		m.history[index] = fastMemInfo{
			timestamp: uint64(now),
			heap:     heap,
			stack:    stack, 
			used:     used,
			numGC:    numGC,
		}
	}
	
	// 更新缓存
	if m.config.EnableCaching {
		cached := &MemoryInfo{
			Timestamp:   time.Unix(0, now),
			GoHeap:      heap,
			GoStack:     stack, // 修正字段名
			UsedMemory:  used,
			MemoryUsage: float64(used) / float64(8*1024*1024*1024) * 100, // 假设8GB系统
		}
		atomic.StorePointer(&m.cachedMemInfo, unsafe.Pointer(cached))
		atomic.StoreInt64(&m.cacheExpiry, now+int64(m.config.CacheExpiry))
		atomic.StoreInt64(&m.cacheValid, 1)
	}
	
	return heap, stack, used, numGC
}

// QuickCheck 超快速健康检查 - 仅检查关键指标
func (m *OptimizedMonitor) QuickCheck() (healthy bool, pressure string) {
	heap, _, used, _ := m.FastMemoryInfo()
	
	if heap == 0 && used == 0 {
		return true, "unknown" // 监控未启用或间隔限制
	}
	
	// 简单的压力评估
	if heap > 512*1024*1024 { // 512MB
		return false, "high"
	} else if heap > 256*1024*1024 { // 256MB
		return true, "medium"
	} else {
		return true, "low"
	}
}

// GetStats 获取轻量级统计信息
func (m *OptimizedMonitor) GetStats() OptimizedStats {
	return OptimizedStats{
		TotalCalls:   atomic.LoadUint64(&m.totalCalls),
		TotalErrors:  atomic.LoadUint64(&m.totalErrors),
		TotalLatency: atomic.LoadUint64(&m.totalLatency),
		Enabled:      atomic.LoadInt64(&m.enabled) == 1,
		CacheHits:    atomic.LoadInt64(&m.cacheValid),
	}
}

// RecordCall 记录调用 - 超轻量级
func (m *OptimizedMonitor) RecordCall(latencyNs uint64, isError bool) {
	atomic.AddUint64(&m.totalCalls, 1)
	atomic.AddUint64(&m.totalLatency, latencyNs)
	if isError {
		atomic.AddUint64(&m.totalErrors, 1)
	}
}

// OptimizedStats 优化统计信息
type OptimizedStats struct {
	TotalCalls   uint64 `json:"total_calls"`
	TotalErrors  uint64 `json:"total_errors"`
	TotalLatency uint64 `json:"total_latency_ns"`
	Enabled      bool   `json:"enabled"`
	CacheHits    int64  `json:"cache_hits"`
}

// GetAverageLatency 获取平均延迟（纳秒）
func (s OptimizedStats) GetAverageLatency() uint64 {
	if s.TotalCalls == 0 {
		return 0
	}
	return s.TotalLatency / s.TotalCalls
}

// GetErrorRate 获取错误率
func (s OptimizedStats) GetErrorRate() float64 {
	if s.TotalCalls == 0 {
		return 0
	}
	return float64(s.TotalErrors) / float64(s.TotalCalls)
}

// UltraLightMonitor 超轻量级监控器 - 仅原子操作
type UltraLightMonitor struct {
	// 所有字段都是原子操作，无锁设计
	enabled        int64  // 是否启用
	totalOps       uint64 // 总操作数
	totalErrors    uint64 // 总错误数
	lastHeapCheck  uint64 // 上次堆检查值
	lastCheckTime  int64  // 上次检查时间
	checkInterval  int64  // 检查间隔（纳秒）
}

// NewUltraLightMonitor 创建超轻量级监控器
func NewUltraLightMonitor() *UltraLightMonitor {
	return &UltraLightMonitor{
		checkInterval: int64(time.Second), // 1秒检查一次
	}
}

// Enable 启用监控
func (u *UltraLightMonitor) Enable() {
	atomic.StoreInt64(&u.enabled, 1)
	atomic.StoreInt64(&u.lastCheckTime, time.Now().UnixNano())
}

// Disable 禁用监控
func (u *UltraLightMonitor) Disable() {
	atomic.StoreInt64(&u.enabled, 0)
}

// Track 跟踪操作 - 零开销版本
func (u *UltraLightMonitor) Track() func(error) {
	if atomic.LoadInt64(&u.enabled) == 0 {
		return func(error) {} // 返回空函数，零开销
	}
	
	atomic.AddUint64(&u.totalOps, 1)
	
	return func(err error) {
		if err != nil {
			atomic.AddUint64(&u.totalErrors, 1)
		}
	}
}

// CheckMemory 检查内存 - 限频版本
func (u *UltraLightMonitor) CheckMemory() (heap uint64, changed bool) {
	if atomic.LoadInt64(&u.enabled) == 0 {
		return 0, false
	}
	
	now := time.Now().UnixNano()
	lastCheck := atomic.LoadInt64(&u.lastCheckTime)
	
	// 限制检查频率
	if now-lastCheck < atomic.LoadInt64(&u.checkInterval) {
		return atomic.LoadUint64(&u.lastHeapCheck), false
	}
	
	// 尝试原子更新检查时间
	if !atomic.CompareAndSwapInt64(&u.lastCheckTime, lastCheck, now) {
		return atomic.LoadUint64(&u.lastHeapCheck), false
	}
	
	// 执行检查
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	
	oldHeap := atomic.SwapUint64(&u.lastHeapCheck, ms.HeapAlloc)
	
	return ms.HeapAlloc, ms.HeapAlloc != oldHeap
}

// Stats 获取统计信息
func (u *UltraLightMonitor) Stats() (ops, errors uint64, errorRate float64) {
	ops = atomic.LoadUint64(&u.totalOps)
	errors = atomic.LoadUint64(&u.totalErrors)
	
	if ops == 0 {
		errorRate = 0
	} else {
		errorRate = float64(errors) / float64(ops)
	}
	
	return ops, errors, errorRate
}

// String 字符串表示
func (u *UltraLightMonitor) String() string {
	ops, errors, errorRate := u.Stats()
	heap, _ := u.CheckMemory()
	enabled := atomic.LoadInt64(&u.enabled) == 1
	
	if enabled {
		return fmt.Sprintf("UltraLight[enabled] ops=%d errors=%d rate=%.2f%% heap=%dKB", 
			ops, errors, errorRate*100, heap/1024)
	} else {
		return "UltraLight[disabled]"
	}
}

// MemoryTracker 专门用于内存跟踪的超轻量级实现
type MemoryTracker struct {
	threshold uint64 // 内存阈值
	current   uint64 // 当前内存使用（原子）
	peak      uint64 // 峰值内存使用（原子）
	warnings  uint64 // 警告计数（原子）
}

// NewMemoryTracker 创建内存追踪器
func NewMemoryTracker(thresholdMB uint64) *MemoryTracker {
	return &MemoryTracker{
		threshold: thresholdMB * 1024 * 1024, // 转换为字节
	}
}

// Update 更新内存使用 - 超快速
func (mt *MemoryTracker) Update(heapBytes uint64) bool {
	atomic.StoreUint64(&mt.current, heapBytes)
	
	// 更新峰值
	for {
		peak := atomic.LoadUint64(&mt.peak)
		if heapBytes <= peak {
			break
		}
		if atomic.CompareAndSwapUint64(&mt.peak, peak, heapBytes) {
			break
		}
	}
	
	// 检查阈值
	if heapBytes > mt.threshold {
		atomic.AddUint64(&mt.warnings, 1)
		return true // 超过阈值
	}
	
	return false
}

// Status 获取状态
func (mt *MemoryTracker) Status() (current, peak, threshold uint64, warnings uint64) {
	return atomic.LoadUint64(&mt.current),
		   atomic.LoadUint64(&mt.peak),
		   mt.threshold,
		   atomic.LoadUint64(&mt.warnings)
}