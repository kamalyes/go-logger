/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\metrics\stats.go
 * @Description: 统计模块 - 收集和管理日志统计信息
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package metrics

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// StatsCollector 统计收集器接口
type StatsCollector interface {
	// 记录操作
	Record(operation string, duration time.Duration, size uint64)
	RecordError(operation string, err error)
	RecordCount(metric string, value uint64)
	RecordGauge(metric string, value float64)
	RecordHistogram(metric string, value float64)
	
	// 获取统计
	GetStats() *Stats
	GetOperationStats(operation string) *OperationStats
	GetMetricStats(metric string) *MetricStats
	GetErrorStats() *ErrorStats
	
	// 重置和清理
	Reset()
	Cleanup(age time.Duration)
}

// Stats 综合统计信息
type Stats struct {
	// 基础计数
	TotalOperations   uint64    `json:"total_operations"`
	TotalErrors       uint64    `json:"total_errors"`
	TotalBytes        uint64    `json:"total_bytes"`
	ErrorRate         float64   `json:"error_rate"`
	
	// 时间统计
	StartTime         time.Time `json:"start_time"`
	LastOperationTime time.Time `json:"last_operation_time"`
	Uptime            time.Duration `json:"uptime"`
	
	// 性能统计
	OperationsPerSecond float64   `json:"operations_per_second"`
	BytesPerSecond      float64   `json:"bytes_per_second"`
	AvgOperationTime    time.Duration `json:"avg_operation_time"`
	
	// 详细统计
	Operations map[string]*OperationStats `json:"operations"`
	Metrics    map[string]*MetricStats    `json:"metrics"`
	Errors     *ErrorStats                `json:"errors"`
}

// OperationStats 操作统计
type OperationStats struct {
	Name         string        `json:"name"`
	Count        uint64        `json:"count"`
	TotalTime    time.Duration `json:"total_time"`
	AvgTime      time.Duration `json:"avg_time"`
	MinTime      time.Duration `json:"min_time"`
	MaxTime      time.Duration `json:"max_time"`
	TotalBytes   uint64        `json:"total_bytes"`
	AvgBytes     float64       `json:"avg_bytes"`
	ErrorCount   uint64        `json:"error_count"`
	ErrorRate    float64       `json:"error_rate"`
	LastExecuted time.Time     `json:"last_executed"`
	FirstExecuted time.Time    `json:"first_executed"`
}

// MetricStats 指标统计
type MetricStats struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"` // counter, gauge, histogram
	Count       uint64    `json:"count"`
	Sum         float64   `json:"sum"`
	Min         float64   `json:"min"`
	Max         float64   `json:"max"`
	Mean        float64   `json:"mean"`
	StdDev      float64   `json:"std_dev"`
	Percentiles map[int]float64 `json:"percentiles"` // P50, P90, P95, P99
	LastUpdate  time.Time `json:"last_update"`
}

// ErrorStats 错误统计
type ErrorStats struct {
	TotalErrors    uint64                    `json:"total_errors"`
	ErrorsByType   map[string]uint64         `json:"errors_by_type"`
	ErrorsByOperation map[string]uint64      `json:"errors_by_operation"`
	RecentErrors   []ErrorRecord             `json:"recent_errors"`
	FirstError     time.Time                 `json:"first_error"`
	LastError      time.Time                 `json:"last_error"`
}

// ErrorRecord 错误记录
type ErrorRecord struct {
	Operation string    `json:"operation"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
	Count     uint64    `json:"count"`
}

// DefaultStatsCollector 默认统计收集器实现
type DefaultStatsCollector struct {
	// 基础字段
	startTime    time.Time
	totalOps     uint64
	totalErrors  uint64
	totalBytes   uint64
	
	// 操作统计
	operations   map[string]*operationStats
	
	// 指标统计
	counters     map[string]*counterMetric
	gauges       map[string]*gaugeMetric
	histograms   map[string]*histogramMetric
	
	// 错误统计
	errorsByType map[string]uint64
	errorsByOp   map[string]uint64
	recentErrors []ErrorRecord
	firstError   time.Time
	lastError    time.Time
	
	// 配置
	maxRecentErrors int
	enablePercentiles bool
	percentiles     []int
	
	mu sync.RWMutex
}

// 内部统计结构
type operationStats struct {
	count        uint64
	totalTime    int64 // 纳秒
	minTime      int64
	maxTime      int64
	totalBytes   uint64
	errorCount   uint64
	lastExecuted int64
	firstExecuted int64
}

type counterMetric struct {
	value      uint64
	lastUpdate int64
}

type gaugeMetric struct {
	value      float64
	lastUpdate int64
}

type histogramMetric struct {
	values     []float64
	count      uint64
	sum        float64
	min        float64
	max        float64
	lastUpdate int64
	mu         sync.Mutex
}

// NewDefaultStatsCollector 创建默认统计收集器
func NewDefaultStatsCollector() *DefaultStatsCollector {
	return &DefaultStatsCollector{
		startTime:         time.Now(),
		operations:        make(map[string]*operationStats),
		counters:          make(map[string]*counterMetric),
		gauges:            make(map[string]*gaugeMetric),
		histograms:        make(map[string]*histogramMetric),
		errorsByType:      make(map[string]uint64),
		errorsByOp:        make(map[string]uint64),
		recentErrors:      make([]ErrorRecord, 0),
		maxRecentErrors:   100,
		enablePercentiles: true,
		percentiles:       []int{50, 90, 95, 99},
	}
}

// Record 记录操作
func (sc *DefaultStatsCollector) Record(operation string, duration time.Duration, size uint64) {
	now := time.Now().UnixNano()
	durationNanos := duration.Nanoseconds()
	
	// 更新全局计数
	atomic.AddUint64(&sc.totalOps, 1)
	atomic.AddUint64(&sc.totalBytes, size)
	
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	// 更新操作统计
	stats, ok := sc.operations[operation]
	if !ok {
		stats = &operationStats{
			minTime:       math.MaxInt64,
			firstExecuted: now,
		}
		sc.operations[operation] = stats
	}
	
	atomic.AddUint64(&stats.count, 1)
	atomic.AddInt64(&stats.totalTime, durationNanos)
	atomic.AddUint64(&stats.totalBytes, size)
	atomic.StoreInt64(&stats.lastExecuted, now)
	
	// 更新最小最大时间
	for {
		currentMin := atomic.LoadInt64(&stats.minTime)
		if durationNanos >= currentMin || atomic.CompareAndSwapInt64(&stats.minTime, currentMin, durationNanos) {
			break
		}
	}
	
	for {
		currentMax := atomic.LoadInt64(&stats.maxTime)
		if durationNanos <= currentMax || atomic.CompareAndSwapInt64(&stats.maxTime, currentMax, durationNanos) {
			break
		}
	}
}

// RecordError 记录错误
func (sc *DefaultStatsCollector) RecordError(operation string, err error) {
	atomic.AddUint64(&sc.totalErrors, 1)
	
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	errorType := fmt.Sprintf("%T", err)
	sc.errorsByType[errorType]++
	sc.errorsByOp[operation]++
	
	// 更新操作错误计数
	if stats, ok := sc.operations[operation]; ok {
		atomic.AddUint64(&stats.errorCount, 1)
	}
	
	// 记录最近错误
	errorRecord := ErrorRecord{
		Operation: operation,
		Error:     err.Error(),
		Timestamp: time.Now(),
		Count:     1,
	}
	
	// 检查是否已有相同错误
	found := false
	for i := range sc.recentErrors {
		if sc.recentErrors[i].Operation == operation && sc.recentErrors[i].Error == err.Error() {
			sc.recentErrors[i].Count++
			sc.recentErrors[i].Timestamp = errorRecord.Timestamp
			found = true
			break
		}
	}
	
	if !found {
		sc.recentErrors = append(sc.recentErrors, errorRecord)
		if len(sc.recentErrors) > sc.maxRecentErrors {
			sc.recentErrors = sc.recentErrors[1:]
		}
	}
	
	// 更新错误时间
	now := time.Now()
	if sc.firstError.IsZero() {
		sc.firstError = now
	}
	sc.lastError = now
}

// RecordCount 记录计数指标
func (sc *DefaultStatsCollector) RecordCount(metric string, value uint64) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	counter, ok := sc.counters[metric]
	if !ok {
		counter = &counterMetric{}
		sc.counters[metric] = counter
	}
	
	atomic.AddUint64(&counter.value, value)
	atomic.StoreInt64(&counter.lastUpdate, time.Now().UnixNano())
}

// RecordGauge 记录仪表指标
func (sc *DefaultStatsCollector) RecordGauge(metric string, value float64) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	gauge, ok := sc.gauges[metric]
	if !ok {
		gauge = &gaugeMetric{}
		sc.gauges[metric] = gauge
	}
	
	gauge.value = value
	atomic.StoreInt64(&gauge.lastUpdate, time.Now().UnixNano())
}

// RecordHistogram 记录直方图指标
func (sc *DefaultStatsCollector) RecordHistogram(metric string, value float64) {
	sc.mu.Lock()
	histogram, ok := sc.histograms[metric]
	if !ok {
		histogram = &histogramMetric{
			values: make([]float64, 0),
			min:    math.MaxFloat64,
			max:    -math.MaxFloat64,
		}
		sc.histograms[metric] = histogram
	}
	sc.mu.Unlock()
	
	histogram.mu.Lock()
	defer histogram.mu.Unlock()
	
	histogram.values = append(histogram.values, value)
	histogram.count++
	histogram.sum += value
	
	if value < histogram.min {
		histogram.min = value
	}
	if value > histogram.max {
		histogram.max = value
	}
	
	atomic.StoreInt64(&histogram.lastUpdate, time.Now().UnixNano())
}

// GetStats 获取综合统计
func (sc *DefaultStatsCollector) GetStats() *Stats {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	
	now := time.Now()
	uptime := now.Sub(sc.startTime)
	totalOps := atomic.LoadUint64(&sc.totalOps)
	totalErrors := atomic.LoadUint64(&sc.totalErrors)
	totalBytes := atomic.LoadUint64(&sc.totalBytes)
	
	var errorRate float64
	if totalOps > 0 {
		errorRate = float64(totalErrors) / float64(totalOps)
	}
	
	var opsPerSec, bytesPerSec float64
	if uptime.Seconds() > 0 {
		opsPerSec = float64(totalOps) / uptime.Seconds()
		bytesPerSec = float64(totalBytes) / uptime.Seconds()
	}
	
	// 计算平均操作时间
	var avgOpTime time.Duration
	if totalOps > 0 {
		var totalTime int64
		for _, stats := range sc.operations {
			totalTime += atomic.LoadInt64(&stats.totalTime)
		}
		avgOpTime = time.Duration(totalTime / int64(totalOps))
	}
	
	// 构建操作统计
	operations := make(map[string]*OperationStats)
	for name, stats := range sc.operations {
		operations[name] = sc.buildOperationStats(name, stats)
	}
	
	// 构建指标统计
	metrics := make(map[string]*MetricStats)
	for name, counter := range sc.counters {
		metrics[name] = &MetricStats{
			Name:       name,
			Type:       "counter",
			Count:      1,
			Sum:        float64(atomic.LoadUint64(&counter.value)),
			LastUpdate: time.Unix(0, atomic.LoadInt64(&counter.lastUpdate)),
		}
	}
	
	for name, gauge := range sc.gauges {
		metrics[name] = &MetricStats{
			Name:       name,
			Type:       "gauge",
			Count:      1,
			Sum:        gauge.value,
			LastUpdate: time.Unix(0, atomic.LoadInt64(&gauge.lastUpdate)),
		}
	}
	
	for name, histogram := range sc.histograms {
		metrics[name] = sc.buildHistogramStats(name, histogram)
	}
	
	// 构建错误统计
	errorStats := &ErrorStats{
		TotalErrors:       totalErrors,
		ErrorsByType:      make(map[string]uint64),
		ErrorsByOperation: make(map[string]uint64),
		RecentErrors:      make([]ErrorRecord, len(sc.recentErrors)),
		FirstError:        sc.firstError,
		LastError:         sc.lastError,
	}
	
	for k, v := range sc.errorsByType {
		errorStats.ErrorsByType[k] = v
	}
	for k, v := range sc.errorsByOp {
		errorStats.ErrorsByOperation[k] = v
	}
	copy(errorStats.RecentErrors, sc.recentErrors)
	
	return &Stats{
		TotalOperations:     totalOps,
		TotalErrors:         totalErrors,
		TotalBytes:          totalBytes,
		ErrorRate:           errorRate,
		StartTime:           sc.startTime,
		LastOperationTime:   sc.getLastOperationTime(),
		Uptime:              uptime,
		OperationsPerSecond: opsPerSec,
		BytesPerSecond:      bytesPerSec,
		AvgOperationTime:    avgOpTime,
		Operations:          operations,
		Metrics:             metrics,
		Errors:              errorStats,
	}
}

// buildOperationStats 构建操作统计
func (sc *DefaultStatsCollector) buildOperationStats(name string, stats *operationStats) *OperationStats {
	count := atomic.LoadUint64(&stats.count)
	totalTime := atomic.LoadInt64(&stats.totalTime)
	totalBytes := atomic.LoadUint64(&stats.totalBytes)
	errorCount := atomic.LoadUint64(&stats.errorCount)
	
	var avgTime time.Duration
	var avgBytes float64
	var errorRate float64
	
	if count > 0 {
		avgTime = time.Duration(totalTime / int64(count))
		avgBytes = float64(totalBytes) / float64(count)
		errorRate = float64(errorCount) / float64(count)
	}
	
	return &OperationStats{
		Name:          name,
		Count:         count,
		TotalTime:     time.Duration(totalTime),
		AvgTime:       avgTime,
		MinTime:       time.Duration(atomic.LoadInt64(&stats.minTime)),
		MaxTime:       time.Duration(atomic.LoadInt64(&stats.maxTime)),
		TotalBytes:    totalBytes,
		AvgBytes:      avgBytes,
		ErrorCount:    errorCount,
		ErrorRate:     errorRate,
		LastExecuted:  time.Unix(0, atomic.LoadInt64(&stats.lastExecuted)),
		FirstExecuted: time.Unix(0, atomic.LoadInt64(&stats.firstExecuted)),
	}
}

// buildHistogramStats 构建直方图统计
func (sc *DefaultStatsCollector) buildHistogramStats(name string, histogram *histogramMetric) *MetricStats {
	histogram.mu.Lock()
	defer histogram.mu.Unlock()
	
	stats := &MetricStats{
		Name:        name,
		Type:        "histogram",
		Count:       histogram.count,
		Sum:         histogram.sum,
		Min:         histogram.min,
		Max:         histogram.max,
		LastUpdate:  time.Unix(0, atomic.LoadInt64(&histogram.lastUpdate)),
		Percentiles: make(map[int]float64),
	}
	
	if histogram.count > 0 {
		stats.Mean = histogram.sum / float64(histogram.count)
		
		// 计算标准差
		var variance float64
		for _, value := range histogram.values {
			diff := value - stats.Mean
			variance += diff * diff
		}
		variance /= float64(histogram.count)
		stats.StdDev = math.Sqrt(variance)
		
		// 计算百分位数
		if sc.enablePercentiles && len(histogram.values) > 0 {
			sorted := make([]float64, len(histogram.values))
			copy(sorted, histogram.values)
			sort.Float64s(sorted)
			
			for _, p := range sc.percentiles {
				index := int(float64(p)/100.0*float64(len(sorted)-1) + 0.5)
				if index >= len(sorted) {
					index = len(sorted) - 1
				}
				stats.Percentiles[p] = sorted[index]
			}
		}
	}
	
	return stats
}

// getLastOperationTime 获取最后操作时间
func (sc *DefaultStatsCollector) getLastOperationTime() time.Time {
	var lastTime int64
	for _, stats := range sc.operations {
		opLastTime := atomic.LoadInt64(&stats.lastExecuted)
		if opLastTime > lastTime {
			lastTime = opLastTime
		}
	}
	
	if lastTime > 0 {
		return time.Unix(0, lastTime)
	}
	return time.Time{}
}

// GetOperationStats 获取特定操作统计
func (sc *DefaultStatsCollector) GetOperationStats(operation string) *OperationStats {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	
	stats, ok := sc.operations[operation]
	if !ok {
		return nil
	}
	
	return sc.buildOperationStats(operation, stats)
}

// GetMetricStats 获取特定指标统计
func (sc *DefaultStatsCollector) GetMetricStats(metric string) *MetricStats {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	
	// 检查计数器
	if counter, ok := sc.counters[metric]; ok {
		return &MetricStats{
			Name:       metric,
			Type:       "counter",
			Count:      1,
			Sum:        float64(atomic.LoadUint64(&counter.value)),
			LastUpdate: time.Unix(0, atomic.LoadInt64(&counter.lastUpdate)),
		}
	}
	
	// 检查仪表
	if gauge, ok := sc.gauges[metric]; ok {
		return &MetricStats{
			Name:       metric,
			Type:       "gauge",
			Count:      1,
			Sum:        gauge.value,
			LastUpdate: time.Unix(0, atomic.LoadInt64(&gauge.lastUpdate)),
		}
	}
	
	// 检查直方图
	if histogram, ok := sc.histograms[metric]; ok {
		return sc.buildHistogramStats(metric, histogram)
	}
	
	return nil
}

// GetErrorStats 获取错误统计
func (sc *DefaultStatsCollector) GetErrorStats() *ErrorStats {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	
	errorStats := &ErrorStats{
		TotalErrors:       atomic.LoadUint64(&sc.totalErrors),
		ErrorsByType:      make(map[string]uint64),
		ErrorsByOperation: make(map[string]uint64),
		RecentErrors:      make([]ErrorRecord, len(sc.recentErrors)),
		FirstError:        sc.firstError,
		LastError:         sc.lastError,
	}
	
	for k, v := range sc.errorsByType {
		errorStats.ErrorsByType[k] = v
	}
	for k, v := range sc.errorsByOp {
		errorStats.ErrorsByOperation[k] = v
	}
	copy(errorStats.RecentErrors, sc.recentErrors)
	
	return errorStats
}

// Reset 重置统计
func (sc *DefaultStatsCollector) Reset() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	atomic.StoreUint64(&sc.totalOps, 0)
	atomic.StoreUint64(&sc.totalErrors, 0)
	atomic.StoreUint64(&sc.totalBytes, 0)
	
	sc.startTime = time.Now()
	sc.operations = make(map[string]*operationStats)
	sc.counters = make(map[string]*counterMetric)
	sc.gauges = make(map[string]*gaugeMetric)
	sc.histograms = make(map[string]*histogramMetric)
	sc.errorsByType = make(map[string]uint64)
	sc.errorsByOp = make(map[string]uint64)
	sc.recentErrors = make([]ErrorRecord, 0)
	sc.firstError = time.Time{}
	sc.lastError = time.Time{}
}

// Cleanup 清理旧数据
func (sc *DefaultStatsCollector) Cleanup(age time.Duration) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	cutoff := time.Now().Add(-age)
	
	// 清理操作统计
	for name, stats := range sc.operations {
		if time.Unix(0, atomic.LoadInt64(&stats.lastExecuted)).Before(cutoff) {
			delete(sc.operations, name)
		}
	}
	
	// 清理指标统计
	for name, counter := range sc.counters {
		if time.Unix(0, atomic.LoadInt64(&counter.lastUpdate)).Before(cutoff) {
			delete(sc.counters, name)
		}
	}
	
	for name, gauge := range sc.gauges {
		if time.Unix(0, atomic.LoadInt64(&gauge.lastUpdate)).Before(cutoff) {
			delete(sc.gauges, name)
		}
	}
	
	for name, histogram := range sc.histograms {
		if time.Unix(0, atomic.LoadInt64(&histogram.lastUpdate)).Before(cutoff) {
			delete(sc.histograms, name)
		}
	}
	
	// 清理最近错误
	var filteredErrors []ErrorRecord
	for _, errorRecord := range sc.recentErrors {
		if errorRecord.Timestamp.After(cutoff) {
			filteredErrors = append(filteredErrors, errorRecord)
		}
	}
	sc.recentErrors = filteredErrors
}

// SetMaxRecentErrors 设置最大最近错误数
func (sc *DefaultStatsCollector) SetMaxRecentErrors(max int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	sc.maxRecentErrors = max
	if len(sc.recentErrors) > max {
		sc.recentErrors = sc.recentErrors[len(sc.recentErrors)-max:]
	}
}

// SetPercentiles 设置百分位数
func (sc *DefaultStatsCollector) SetPercentiles(percentiles []int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	sc.percentiles = percentiles
	sc.enablePercentiles = len(percentiles) > 0
}