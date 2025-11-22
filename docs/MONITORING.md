# Go Logger 监控系统指南

## 目录

- [📊 监控概述](#-监控概述)
- [🧠 内存监控](#-内存监控)
- [⚡ 性能监控](#-性能监控)
- [💾 I/O 监控](#-io-监控)
- [🎯 指标收集](#-指标收集)
- [📈 监控配置](#-监控配置)
- [🚨 告警系统](#-告警系统)

## 📊 监控概述

go-logger 提供了全面的监控系统，可以实时监控日志系统的性能、内存使用、I/O 状态等关键指标。监控系统采用分层架构设计，提供三种不同级别的监控：

### 监控架构

```
┌─────────────────────────────────────────────┐
│              监控层 (Monitoring Layer)        │
├─────────────────┬─────────────────┬─────────┤
│   内存监控       │    性能监控       │  I/O监控 │
│ MemoryMonitor   │ PerformanceMonitor│IOMonitor│
└─────────────────┴─────────────────┴─────────┘
┌─────────────────────────────────────────────┐
│             指标层 (Metrics Layer)            │
├─────────────────┬─────────────────┬─────────┤
│   统计收集       │    告警管理       │  数据存储│
│ StatsCollector  │  AlertManager    │DataStore│
└─────────────────┴─────────────────┴─────────┘
┌─────────────────────────────────────────────┐
│            应用层 (Application Layer)         │
├─────────────────┬─────────────────┬─────────┤
│    日志器        │     适配器       │   钩子   │
│   Logger        │    Adapters      │ Hooks   │
└─────────────────┴─────────────────┴─────────┘
```

### 监控级别

| 级别 | 性能开销 | 功能完整度 | 适用场景 |
|------|---------|-----------|----------|
| UltraLight | 3.134ns | ⭐⭐ | 极高性能要求 |
| Optimized | 3.094ns | ⭐⭐⭐⭐ | 一般生产环境 |
| Full | 122.4ns | ⭐⭐⭐⭐⭐ | 企业级应用 |

## 🧠 内存监控

### 内存监控器

```go
import "github.com/kamalyes/go-logger/metrics"

// 创建内存监控器
monitor := metrics.NewDefaultMemoryMonitor()

// 基础配置
monitor.SetSampleInterval(time.Second * 5)      // 采样间隔
monitor.SetMemoryThreshold(85.0)                // 内存阈值85%
monitor.SetMaxMemory(4 * 1024 * 1024 * 1024)   // 最大内存4GB

// 高级配置
monitor.EnableLeakDetection(true)               // 启用泄漏检测
monitor.SetMaxHistorySize(100)                 // 历史记录数量
monitor.SetGCPercent(75)                       // GC百分比

// 启动监控
if err := monitor.Start(); err != nil {
    log.Fatal("启动内存监控失败:", err)
}
defer monitor.Stop()
```

### 内存信息获取

```go
// 获取实时内存信息
memInfo := monitor.GetMemoryInfo()
fmt.Printf("内存监控报告:\n")
fmt.Printf("  使用率: %.2f%%\n", memInfo.MemoryUsage)
fmt.Printf("  总内存: %.2f GB\n", float64(memInfo.TotalMemory)/1024/1024/1024)
fmt.Printf("  已用内存: %.2f MB\n", float64(memInfo.UsedMemory)/1024/1024)
fmt.Printf("  Go堆内存: %.2f MB\n", float64(memInfo.GoHeap)/1024/1024)
fmt.Printf("  Go栈内存: %.2f MB\n", float64(memInfo.GoStack)/1024/1024)
fmt.Printf("  Go系统内存: %.2f MB\n", float64(memInfo.GoSys)/1024/1024)

// 获取详细的内存统计
detailInfo := monitor.GetDetailedMemoryInfo()
fmt.Printf("详细内存统计:\n")
fmt.Printf("  堆对象数: %d\n", detailInfo.HeapObjects)
fmt.Printf("  堆分配数: %d\n", detailInfo.HeapAlloc)
fmt.Printf("  堆空闲数: %d\n", detailInfo.HeapIdle)
fmt.Printf("  堆已释放: %d\n", detailInfo.HeapReleased)
fmt.Printf("  下次GC触发: %d bytes\n", detailInfo.NextGC)

// 获取GC信息
gcInfo := monitor.GetGCInfo()
fmt.Printf("GC信息:\n")
fmt.Printf("  GC次数: %d\n", gcInfo.NumGC)
fmt.Printf("  总GC时间: %v\n", gcInfo.PauseTotalNs)
fmt.Printf("  平均GC时间: %.2f ms\n", float64(gcInfo.PauseTotalNs)/float64(gcInfo.NumGC)/1e6)
fmt.Printf("  最后GC时间: %s\n", time.Unix(0, int64(gcInfo.LastGC)).Format("2006-01-02 15:04:05"))
```

### 内存快照和分析

```go
// 创建内存快照
snapshot, err := monitor.TakeHeapSnapshot()
if err != nil {
    log.Printf("创建内存快照失败: %v", err)
} else {
    fmt.Printf("内存快照:\n")
    fmt.Printf("  时间: %s\n", snapshot.Timestamp.Format("2006-01-02 15:04:05"))
    fmt.Printf("  对象数量: %d\n", snapshot.ObjectCount)
    fmt.Printf("  内存大小: %.2f MB\n", float64(snapshot.MemorySize)/1024/1024)
    fmt.Printf("  分配数量: %d\n", snapshot.AllocCount)
    fmt.Printf("  栈信息: %s\n", snapshot.StackTrace)
}

// 分析内存历史趋势
history := monitor.GetMemoryHistory(10) // 获取最近10个快照
fmt.Printf("内存历史趋势:\n")
for i, h := range history {
    growth := ""
    if i > 0 {
        prevUsage := history[i-1].MemoryUsage
        if h.MemoryUsage > prevUsage {
            growth = fmt.Sprintf(" (+%.2f%%)", h.MemoryUsage-prevUsage)
        } else {
            growth = fmt.Sprintf(" (%.2f%%)", h.MemoryUsage-prevUsage)
        }
    }
    fmt.Printf("  %s: %.2f%%%s\n", 
        h.Timestamp.Format("15:04:05"), h.MemoryUsage, growth)
}

// 内存泄漏分析
report := monitor.AnalyzeMemoryLeaks()
fmt.Printf("泄漏分析:\n")
fmt.Printf("  增长趋势: %s\n", report.GrowthTrend)
fmt.Printf("  增长率: %.2f bytes/s\n", report.MemoryGrowthRate)
fmt.Printf("  风险级别: %s\n", report.RiskLevel)
fmt.Printf("  建议: %s\n", report.Recommendation)

if report.RiskLevel == "HIGH" {
    fmt.Printf("  🚨 高风险泄漏检测到!\n")
    fmt.Printf("  可疑分配点:\n")
    for _, point := range report.SuspiciousPoints {
        fmt.Printf("    - %s: %d bytes\n", point.Function, point.Bytes)
    }
}
```

### 内存事件回调

```go
// 设置内存阈值超出回调
monitor.OnMemoryThresholdExceeded(func(info *metrics.MemoryInfo) {
    log.Warn("⚠️ 内存使用率超出阈值",
        "usage_percent", info.MemoryUsage,
        "used_mb", info.UsedMemory/1024/1024,
        "heap_mb", info.GoHeap/1024/1024,
        "threshold", 85.0)
    
    // 可以触发告警或清理操作
    if info.MemoryUsage > 90.0 {
        log.Warn("内存使用率过高，强制执行GC")
        runtime.GC()
        runtime.GC() // 连续两次GC确保充分回收
        
        // 发送告警
        sendAlert("high_memory_usage", map[string]interface{}{
            "usage":   info.MemoryUsage,
            "used_mb": info.UsedMemory / 1024 / 1024,
            "action":  "force_gc",
        })
    }
    
    if info.MemoryUsage > 95.0 {
        log.Error("🚨 内存使用率极高，可能需要重启应用")
        sendCriticalAlert("critical_memory_usage", info)
    }
})

// 设置内存泄漏检测回调
monitor.OnMemoryLeakDetected(func(report *metrics.LeakReport) {
    log.Error("🚨 检测到内存泄漏",
        "trend", report.GrowthTrend,
        "rate_bytes_per_sec", report.MemoryGrowthRate,
        "risk_level", report.RiskLevel,
        "duration", report.Duration)
    
    // 记录详细信息
    log.Debug("内存泄漏详细信息",
        "start_memory_mb", report.StartMemory/1024/1024,
        "current_memory_mb", report.CurrentMemory/1024/1024,
        "growth_mb", (report.CurrentMemory-report.StartMemory)/1024/1024,
        "growth_rate_mb_per_min", report.MemoryGrowthRate*60/1024/1024)
    
    // 发送告警和建议
    alertData := map[string]interface{}{
        "risk_level":    report.RiskLevel,
        "growth_trend":  report.GrowthTrend,
        "growth_rate":   report.MemoryGrowthRate,
        "recommendation": report.Recommendation,
    }
    
    switch report.RiskLevel {
    case "LOW":
        sendInfoAlert("memory_leak_detected", alertData)
    case "MEDIUM":
        sendWarningAlert("memory_leak_detected", alertData)
    case "HIGH":
        sendCriticalAlert("memory_leak_detected", alertData)
        
        // 高风险时自动执行一些清理操作
        log.Info("执行内存清理操作")
        runtime.GC()
        runtime.FreeOSMemory()
        
        // 如果有可疑分配点，记录详细信息
        for _, point := range report.SuspiciousPoints {
            log.Warn("可疑内存分配点",
                "function", point.Function,
                "bytes", point.Bytes,
                "count", point.Count)
        }
    }
})

// 设置GC事件回调
monitor.OnGCCompleted(func(gcInfo *metrics.GCInfo) {
    // 计算GC效率
    avgPause := float64(gcInfo.PauseTotalNs) / float64(gcInfo.NumGC) / 1e6
    
    log.Debug("GC完成",
        "gc_count", gcInfo.NumGC,
        "avg_pause_ms", avgPause,
        "last_pause_ms", float64(gcInfo.LastPause)/1e6)
    
    // 如果GC暂停时间过长，发出警告
    if avgPause > 10.0 { // 10ms
        log.Warn("GC暂停时间过长",
            "avg_pause_ms", avgPause,
            "threshold_ms", 10.0,
            "suggestion", "考虑调整GOGC参数或优化内存分配")
    }
    
    // 记录GC统计到指标系统
    metricsCollector.RecordGCMetrics(gcInfo)
})
```

### 超轻量级内存监控

适用于高性能场景的极简内存监控：

```go
// 创建超轻量级监控器
ultraMonitor := metrics.NewUltraLightMonitor()
ultraMonitor.Enable()

// 在关键路径中使用
func criticalPathFunction() {
    done := ultraMonitor.Track()  // 开始追踪，仅3.134ns开销
    defer done(nil)               // 结束追踪
    
    // 执行业务逻辑...
    processBusinessLogic()
}

// 获取快速内存状态
func quickMemoryCheck() {
    heap, stack, used, numGC := ultraMonitor.FastMemoryInfo()
    
    // 简单的内存压力检测
    if used > 1024*1024*1024 { // 1GB
        log.Warn("内存使用量较高", "used_bytes", used)
    }
    
    // 统计信息（极低开销）
    stats := ultraMonitor.GetStats()
    if stats.TotalOperations > 0 {
        avgMemory := stats.TotalMemoryUsed / stats.TotalOperations
        log.Debug("平均内存使用", "avg_bytes", avgMemory)
    }
}
```

### 优化内存监控

智能缓存和批处理的优化监控：

```go
// 创建优化监控器
optimizedConfig := metrics.OptimizedConfig{
    CacheExpiry:     100 * time.Millisecond, // 缓存100ms
    EnableCaching:   true,
    LightweightMode: true,
    BatchInterval:   time.Second,             // 批处理间隔
    BatchSize:       100,                     // 批处理大小
}
optimizedMonitor := metrics.NewOptimizedMonitor(optimizedConfig)

// 启动优化监控
optimizedMonitor.Start()
defer optimizedMonitor.Stop()

// 快速获取缓存的内存信息
heap, stack, used, numGC := optimizedMonitor.FastMemoryInfo()
fmt.Printf("快速内存信息: 堆=%d, 栈=%d, 已用=%d, GC=%d\n", 
    heap, stack, used, numGC)

// 内存追踪器 - 阈值检测
tracker := metrics.NewMemoryTracker(512 * 1024 * 1024) // 512MB阈值
exceeded := tracker.Update(used)
if exceeded {
    log.Warn("内存使用超过阈值", "threshold_mb", 512, "used_mb", used/1024/1024)
}

// 智能健康检查
healthy, pressure := optimizedMonitor.QuickCheck()
fmt.Printf("系统健康: %v, 内存压力: %s\n", healthy, pressure)

// 根据内存压力调整监控频率
switch pressure {
case "LOW":
    optimizedMonitor.SetSampleInterval(time.Second * 10)
case "MEDIUM":
    optimizedMonitor.SetSampleInterval(time.Second * 5)
case "HIGH":
    optimizedMonitor.SetSampleInterval(time.Second * 1)
case "CRITICAL":
    // 切换到全功能监控
    fullMonitor := metrics.NewDefaultMemoryMonitor()
    fullMonitor.Start()
    optimizedMonitor.Stop()
}
```

## ⚡ 性能监控

### 性能监控器

```go
// 创建性能监控器
perfMonitor := metrics.NewDefaultPerformanceMonitor()

// 配置性能监控
perfMonitor.SetLatencyThreshold("api", time.Millisecond*100)    // API延迟阈值
perfMonitor.SetThroughputThreshold("requests", 1000.0)          // 请求吞吐量阈值
perfMonitor.SetResourceThreshold(80.0, 85.0)                   // CPU和内存阈值

// 启动性能监控
perfMonitor.Start()
defer perfMonitor.Stop()
```

### 延迟监控

```go
// 记录操作延迟
func monitoredOperation(name string) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        perfMonitor.RecordLatency(name, duration)
    }()
    
    // 执行业务操作...
    performOperation()
}

// 批量延迟记录
func batchOperations() {
    operations := []string{"user_query", "data_process", "cache_update"}
    durations := make([]time.Duration, len(operations))
    
    for i, op := range operations {
        start := time.Now()
        executeOperation(op)
        durations[i] = time.Since(start)
    }
    
    // 批量记录
    perfMonitor.RecordLatencies(operations, durations)
}

// 获取延迟统计
latencyStats := perfMonitor.GetLatencyStats()
for operation, stats := range latencyStats {
    fmt.Printf("操作 %s 延迟统计:\n", operation)
    fmt.Printf("  平均延迟: %v\n", stats.AvgLatency)
    fmt.Printf("  最小延迟: %v\n", stats.MinLatency)
    fmt.Printf("  最大延迟: %v\n", stats.MaxLatency)
    fmt.Printf("  P50延迟: %v\n", stats.P50Latency)
    fmt.Printf("  P95延迟: %v\n", stats.P95Latency)
    fmt.Printf("  P99延迟: %v\n", stats.P99Latency)
    fmt.Printf("  请求总数: %d\n", stats.TotalRequests)
}
```

### 吞吐量监控

```go
// 记录吞吐量
func recordThroughput() {
    // 记录单个操作
    perfMonitor.RecordThroughput("requests", 1)
    perfMonitor.RecordThroughput("messages", 10)
    perfMonitor.RecordThroughput("bytes", 1024)
    
    // 记录批量操作
    perfMonitor.RecordBatchThroughput("batch_process", 100)
}

// 获取吞吐量统计
throughputStats := perfMonitor.GetThroughputStats()
for operation, stats := range throughputStats {
    fmt.Printf("操作 %s 吞吐量统计:\n", operation)
    fmt.Printf("  当前吞吐量: %.2f ops/s\n", stats.CurrentThroughput)
    fmt.Printf("  平均吞吐量: %.2f ops/s\n", stats.AvgThroughput)
    fmt.Printf("  峰值吞吐量: %.2f ops/s\n", stats.PeakThroughput)
    fmt.Printf("  总操作数: %d\n", stats.TotalOperations)
}

// 实时吞吐量监控
go func() {
    ticker := time.NewTicker(time.Second * 5)
    defer ticker.Stop()
    
    for range ticker.C {
        currentStats := perfMonitor.GetCurrentThroughput()
        for operation, throughput := range currentStats {
            if throughput > 1000 {
                log.Info("高吞吐量检测", "operation", operation, "tps", throughput)
            }
        }
    }
}()
```

### 资源使用监控

```go
// 记录资源使用
func monitorResourceUsage() {
    // 手动记录
    perfMonitor.RecordResourceUsage()
    
    // 获取资源统计
    resourceStats := perfMonitor.GetResourceStats()
    fmt.Printf("资源使用统计:\n")
    fmt.Printf("  CPU使用率: %.2f%%\n", resourceStats.CPUUsage)
    fmt.Printf("  内存使用率: %.2f%%\n", resourceStats.MemoryUsage)
    fmt.Printf("  磁盘使用率: %.2f%%\n", resourceStats.DiskUsage)
    fmt.Printf("  网络入带宽: %.2f MB/s\n", resourceStats.NetworkIn/1024/1024)
    fmt.Printf("  网络出带宽: %.2f MB/s\n", resourceStats.NetworkOut/1024/1024)
    fmt.Printf("  文件描述符: %d/%d\n", resourceStats.OpenFiles, resourceStats.MaxFiles)
    fmt.Printf("  线程数: %d\n", resourceStats.ThreadCount)
}

// 自动资源监控
perfMonitor.EnableAutoResourceMonitoring(time.Second * 10) // 每10秒采样一次
```

### 性能事件回调

```go
// 延迟阈值超出回调
perfMonitor.OnLatencyThresholdExceeded(func(operation string, latency time.Duration) {
    log.Warn("⚠️ 操作延迟超标",
        "operation", operation,
        "latency_ms", float64(latency.Nanoseconds())/1e6,
        "threshold_ms", 100)
    
    // 可以触发告警或降级措施
    if latency > time.Millisecond*500 { // 500ms
        log.Error("🚨 严重延迟，考虑降级服务")
        triggerCircuitBreaker(operation)
    }
})

// 吞吐量阈值超出回调
perfMonitor.OnThroughputThresholdExceeded(func(operation string, throughput float64) {
    log.Info("📈 高吞吐量",
        "operation", operation,
        "throughput_ops", throughput,
        "threshold", 1000.0)
    
    // 高吞吐量时可能需要扩容
    if throughput > 5000 {
        log.Warn("🚀 超高吞吐量，建议扩容")
        triggerAutoScaling(operation, throughput)
    }
})

// 资源阈值超出回调
perfMonitor.OnResourceThresholdExceeded(func(usage *metrics.ResourceUsage) {
    if usage.CPUUsage > 80.0 {
        log.Warn("⚠️ CPU使用率过高", "usage", usage.CPUUsage)
    }
    
    if usage.MemoryUsage > 85.0 {
        log.Warn("⚠️ 内存使用率过高", "usage", usage.MemoryUsage)
        runtime.GC() // 尝试释放内存
    }
    
    if usage.DiskUsage > 90.0 {
        log.Error("🚨 磁盘使用率过高", "usage", usage.DiskUsage)
        startLogRotation() // 触发日志清理
    }
    
    if usage.OpenFiles > float64(usage.MaxFiles)*0.9 {
        log.Warn("⚠️ 文件描述符使用率过高",
            "open", usage.OpenFiles, "max", usage.MaxFiles)
    }
})
```

### 性能基准测试

```go
// 性能基准测试
func runPerformanceBenchmark() {
    benchmark := metrics.NewPerformanceBenchmark()
    
    // 设置基准测试参数
    benchmark.SetDuration(time.Minute)      // 测试1分钟
    benchmark.SetConcurrency(100)           // 100个并发
    benchmark.SetOperations([]string{
        "log_info", "log_error", "log_debug",
    })
    
    // 运行基准测试
    results, err := benchmark.Run(func(operation string) error {
        switch operation {
        case "log_info":
            logger.Info("benchmark test message")
        case "log_error":
            logger.Error("benchmark error message")
        case "log_debug":
            logger.Debug("benchmark debug message")
        }
        return nil
    })
    
    if err != nil {
        log.Error("基准测试失败", "error", err)
        return
    }
    
    // 输出测试结果
    fmt.Printf("性能基准测试结果:\n")
    for _, result := range results {
        fmt.Printf("操作 %s:\n", result.Operation)
        fmt.Printf("  总请求数: %d\n", result.TotalRequests)
        fmt.Printf("  成功请求数: %d\n", result.SuccessfulRequests)
        fmt.Printf("  失败请求数: %d\n", result.FailedRequests)
        fmt.Printf("  平均延迟: %v\n", result.AvgLatency)
        fmt.Printf("  P95延迟: %v\n", result.P95Latency)
        fmt.Printf("  P99延迟: %v\n", result.P99Latency)
        fmt.Printf("  吞吐量: %.2f ops/s\n", result.Throughput)
        fmt.Printf("  错误率: %.2f%%\n", result.ErrorRate*100)
    }
}
```

## 💾 I/O 监控

### I/O 监控器

```go
// 创建I/O监控器
ioMonitor := metrics.NewIOMonitor()

// 设置阈值
ioMonitor.SetThresholds(
    80.0,  // 磁盘使用率阈值
    1000,  // IOPS阈值
    100,   // 延迟阈值(ms)
)

// 启动I/O监控
ioMonitor.Start()
defer ioMonitor.Stop()
```

### 磁盘I/O监控

```go
// 获取磁盘I/O统计
diskStats := ioMonitor.GetDiskIOStats()
fmt.Printf("磁盘I/O统计:\n")
fmt.Printf("  读取字节: %.2f MB\n", float64(diskStats.ReadBytes)/1024/1024)
fmt.Printf("  写入字节: %.2f MB\n", float64(diskStats.WriteBytes)/1024/1024)
fmt.Printf("  读取次数: %d\n", diskStats.ReadOps)
fmt.Printf("  写入次数: %d\n", diskStats.WriteOps)
fmt.Printf("  读取延迟: %v\n", diskStats.ReadLatency)
fmt.Printf("  写入延迟: %v\n", diskStats.WriteLatency)
fmt.Printf("  读取IOPS: %.2f\n", diskStats.ReadIOPS)
fmt.Printf("  写入IOPS: %.2f\n", diskStats.WriteIOPS)
fmt.Printf("  磁盘使用率: %.2f%%\n", diskStats.DiskUsage)

// 监控特定文件的I/O
fileIOStats := ioMonitor.GetFileIOStats("/var/log/app.log")
if fileIOStats != nil {
    fmt.Printf("文件I/O统计 (/var/log/app.log):\n")
    fmt.Printf("  写入字节: %.2f MB\n", float64(fileIOStats.WriteBytes)/1024/1024)
    fmt.Printf("  写入次数: %d\n", fileIOStats.WriteOps)
    fmt.Printf("  平均写入延迟: %v\n", fileIOStats.AvgWriteLatency)
}
```

### 网络I/O监控

```go
// 获取网络I/O统计
networkStats := ioMonitor.GetNetworkIOStats()
fmt.Printf("网络I/O统计:\n")
fmt.Printf("  接收字节: %.2f MB\n", float64(networkStats.RxBytes)/1024/1024)
fmt.Printf("  发送字节: %.2f MB\n", float64(networkStats.TxBytes)/1024/1024)
fmt.Printf("  接收包数: %d\n", networkStats.RxPackets)
fmt.Printf("  发送包数: %d\n", networkStats.TxPackets)
fmt.Printf("  接收错误: %d\n", networkStats.RxErrors)
fmt.Printf("  发送错误: %d\n", networkStats.TxErrors)
fmt.Printf("  网络延迟: %v\n", networkStats.Latency)

// 监控特定连接的I/O
connStats := ioMonitor.GetConnectionIOStats("tcp", "elasticsearch:9200")
if connStats != nil {
    fmt.Printf("连接I/O统计 (elasticsearch:9200):\n")
    fmt.Printf("  连接状态: %s\n", connStats.State)
    fmt.Printf("  发送字节: %.2f KB\n", float64(connStats.TxBytes)/1024)
    fmt.Printf("  接收字节: %.2f KB\n", float64(connStats.RxBytes)/1024)
    fmt.Printf("  连接延迟: %v\n", connStats.Latency)
}
```

### I/O 事件回调

```go
// I/O阈值超出回调
ioMonitor.OnThresholdExceeded(func(metric string, value float64) {
    switch metric {
    case "disk_usage":
        log.Warn("磁盘使用率过高", "usage_percent", value)
        if value > 95.0 {
            log.Error("🚨 磁盘空间严重不足")
            // 清理旧日志文件
            cleanupOldLogs()
            // 发送紧急告警
            sendCriticalAlert("disk_full", map[string]interface{}{
                "usage": value,
                "action": "log_cleanup",
            })
        }
        
    case "iops":
        log.Warn("磁盘IOPS过高", "iops", value)
        // 增加批量大小，减少写入频率
        adjustBatchSize(2.0)
        
    case "latency":
        log.Warn("I/O延迟过高", "latency_ms", value)
        if value > 1000 { // 1秒
            log.Error("🚨 I/O延迟极高，可能影响性能")
            // 启用压缩，减少I/O量
            enableCompression()
        }
        
    case "network_error_rate":
        log.Warn("网络错误率过高", "error_rate_percent", value)
        if value > 10.0 {
            log.Error("🚨 网络连接不稳定")
            // 重试连接
            retryNetworkConnections()
        }
    }
})

// I/O性能异常回调
ioMonitor.OnPerformanceAnomaly(func(anomaly *metrics.IOAnomaly) {
    log.Warn("I/O性能异常检测",
        "type", anomaly.Type,
        "severity", anomaly.Severity,
        "description", anomaly.Description,
        "metric_name", anomaly.MetricName,
        "current_value", anomaly.CurrentValue,
        "expected_value", anomaly.ExpectedValue)
    
    switch anomaly.Type {
    case "SUDDEN_LATENCY_SPIKE":
        log.Error("🚨 I/O延迟突然增加")
        // 可能的磁盘问题，需要检查
        checkDiskHealth()
        
    case "THROUGHPUT_DROP":
        log.Warn("⚠️ I/O吞吐量下降")
        // 可能的网络问题或磁盘问题
        diagnosePerfIssues()
        
    case "ERROR_RATE_INCREASE":
        log.Error("🚨 I/O错误率增加")
        // 连接或硬件问题
        escalateToOpsTeam(anomaly)
    }
})
```

## 🎯 指标收集

### 统计收集器

```go
// 创建统计收集器
statsCollector := metrics.NewDefaultStatsCollector()

// 记录操作统计
func recordOperationStats() {
    start := time.Now()
    
    // 执行操作
    err := performOperation()
    
    duration := time.Since(start)
    
    // 记录统计信息
    statsCollector.RecordOperation("user_query", duration, err)
    
    // 记录自定义指标
    statsCollector.RecordCustomMetric("custom.operation.size", 1024)
    statsCollector.RecordCustomMetric("custom.operation.complexity", 5)
}

// 批量记录统计
func batchRecordStats() {
    operations := []metrics.OperationRecord{
        {Name: "db_query", Duration: time.Millisecond * 50, Error: nil},
        {Name: "cache_lookup", Duration: time.Millisecond * 5, Error: nil},
        {Name: "api_call", Duration: time.Millisecond * 200, Error: fmt.Errorf("timeout")},
    }
    
    statsCollector.RecordOperations(operations)
}
```

### 获取统计信息

```go
// 获取所有统计信息
allStats := statsCollector.GetAllStats()
fmt.Printf("所有操作统计:\n")
for operation, stats := range allStats {
    fmt.Printf("操作 %s:\n", operation)
    fmt.Printf("  总数: %d\n", stats.Count)
    fmt.Printf("  成功数: %d\n", stats.SuccessCount)
    fmt.Printf("  失败数: %d\n", stats.ErrorCount)
    fmt.Printf("  成功率: %.2f%%\n", stats.SuccessRate*100)
    fmt.Printf("  平均耗时: %v\n", stats.AvgDuration)
    fmt.Printf("  最小耗时: %v\n", stats.MinDuration)
    fmt.Printf("  最大耗时: %v\n", stats.MaxDuration)
}

// 获取特定操作统计
userQueryStats := statsCollector.GetOperationStats("user_query")
if userQueryStats != nil {
    fmt.Printf("用户查询统计:\n")
    fmt.Printf("  总查询数: %d\n", userQueryStats.Count)
    fmt.Printf("  平均耗时: %v\n", userQueryStats.AvgDuration)
    fmt.Printf("  成功率: %.2f%%\n", userQueryStats.SuccessRate*100)
    
    // 获取最近的错误
    recentErrors := userQueryStats.GetRecentErrors(10)
    if len(recentErrors) > 0 {
        fmt.Printf("  最近错误:\n")
        for _, err := range recentErrors {
            fmt.Printf("    - %s: %v\n", err.Timestamp.Format("15:04:05"), err.Error)
        }
    }
}

// 获取自定义指标
customMetrics := statsCollector.GetCustomMetrics()
for name, value := range customMetrics {
    fmt.Printf("自定义指标 %s: %v\n", name, value)
}
```

### 指标聚合和分析

```go
// 获取时间窗口内的统计
windowStats := statsCollector.GetStatsInTimeWindow(
    time.Now().Add(-time.Hour), // 1小时前
    time.Now(),                 // 现在
)

for operation, stats := range windowStats {
    fmt.Printf("近1小时操作 %s 统计:\n", operation)
    fmt.Printf("  请求量: %d\n", stats.Count)
    fmt.Printf("  QPS: %.2f\n", float64(stats.Count)/3600.0) // 每秒请求数
    fmt.Printf("  平均延迟: %v\n", stats.AvgDuration)
    fmt.Printf("  错误率: %.2f%%\n", stats.ErrorRate*100)
}

// 获取趋势分析
trendAnalysis := statsCollector.GetTrendAnalysis("user_query", time.Hour)
fmt.Printf("用户查询趋势分析 (1小时):\n")
fmt.Printf("  请求量趋势: %s\n", trendAnalysis.RequestTrend)      // INCREASING, DECREASING, STABLE
fmt.Printf("  延迟趋势: %s\n", trendAnalysis.LatencyTrend)        // INCREASING, DECREASING, STABLE
fmt.Printf("  错误率趋势: %s\n", trendAnalysis.ErrorRateTrend)    // INCREASING, DECREASING, STABLE
fmt.Printf("  预测下小时请求量: %d\n", trendAnalysis.PredictedNextHourRequests)

// 异常检测
anomalies := statsCollector.DetectAnomalies("user_query")
for _, anomaly := range anomalies {
    fmt.Printf("检测到异常:\n")
    fmt.Printf("  类型: %s\n", anomaly.Type)
    fmt.Printf("  严重程度: %s\n", anomaly.Severity)
    fmt.Printf("  描述: %s\n", anomaly.Description)
    fmt.Printf("  时间: %s\n", anomaly.Timestamp.Format("2006-01-02 15:04:05"))
}
```

## 📈 监控配置

### 综合监控配置

```yaml
# config/monitoring.yaml
monitoring:
  # 全局设置
  enabled: true
  sampling_rate: 1.0        # 100% 采样
  metrics_interval: 30s     # 指标收集间隔
  retention_period: 24h     # 数据保留期
  
  # 内存监控
  memory:
    enabled: true
    threshold: 85.0           # 内存阈值 85%
    sample_interval: 5s       # 采样间隔
    leak_detection: true      # 启用泄漏检测
    max_history_size: 100     # 历史记录数量
    gc_percent: 75           # GC 百分比
    max_memory: 4GB          # 最大内存限制
    
    # 告警配置
    alerts:
      threshold_exceeded:
        enabled: true
        webhook_url: "http://alert:8080/memory"
      leak_detected:
        enabled: true
        webhook_url: "http://alert:8080/leak"
        
  # 性能监控
  performance:
    enabled: true
    latency_threshold: 100ms  # 延迟阈值
    throughput_threshold: 1000.0  # 吞吐量阈值
    sample_rate: 0.1         # 10% 采样
    enable_profiling: true   # 启用性能分析
    
    # 资源监控
    resource_monitoring:
      enabled: true
      cpu_threshold: 80.0    # CPU 阈值
      memory_threshold: 85.0 # 内存阈值
      disk_threshold: 90.0   # 磁盘阈值
      sample_interval: 10s   # 资源采样间隔
      
    # 告警配置
    alerts:
      latency_exceeded:
        enabled: true
        webhook_url: "http://alert:8080/latency"
      resource_exceeded:
        enabled: true
        webhook_url: "http://alert:8080/resource"
        
  # I/O 监控
  io:
    enabled: true
    disk_usage_threshold: 80.0  # 磁盘使用率阈值
    iops_threshold: 1000        # IOPS 阈值
    latency_threshold: 100ms    # I/O 延迟阈值
    sample_interval: 10s        # 采样间隔
    
    # 文件监控
    file_monitoring:
      enabled: true
      watch_files:
        - "/var/log/app.log"
        - "/var/log/error.log"
        
    # 网络监控
    network_monitoring:
      enabled: true
      connections:
        - "tcp:elasticsearch:9200"
        - "tcp:redis:6379"
        - "tcp:kafka:9092"
        
  # 指标收集
  metrics:
    enabled: true
    collection_interval: 30s    # 收集间隔
    retention_size: 10000       # 保留记录数
    enable_custom_metrics: true # 启用自定义指标
    
    # 导出配置
    exporters:
      prometheus:
        enabled: true
        endpoint: "/metrics"
        namespace: "app"
        
      influxdb:
        enabled: false
        url: "http://influxdb:8086"
        database: "metrics"
        username: "admin"
        password: "password"
        
  # 告警管理
  alerting:
    enabled: true
    default_webhook: "http://alert:8080/webhook"
    retry_attempts: 3
    retry_delay: 5s
    
    # 告警规则
    rules:
      - name: "high_memory_usage"
        condition: "memory.usage > 90"
        severity: "critical"
        message: "内存使用率超过90%"
        
      - name: "high_latency"
        condition: "performance.avg_latency > 500ms"
        severity: "warning"
        message: "平均延迟超过500ms"
        
      - name: "high_error_rate"
        condition: "performance.error_rate > 5"
        severity: "warning"
        message: "错误率超过5%"
        
      - name: "disk_space_low"
        condition: "io.disk_usage > 95"
        severity: "critical"
        message: "磁盘空间不足5%"
```

### 编程方式配置

```go
// 创建监控配置
config := &MonitoringConfig{
    Enabled:         true,
    SamplingRate:    1.0,
    MetricsInterval: time.Second * 30,
    RetentionPeriod: time.Hour * 24,
    
    Memory: MemoryMonitoringConfig{
        Enabled:         true,
        Threshold:       85.0,
        SampleInterval:  time.Second * 5,
        LeakDetection:   true,
        MaxHistorySize:  100,
        GCPercent:       75,
        MaxMemory:       4 * 1024 * 1024 * 1024, // 4GB
    },
    
    Performance: PerformanceMonitoringConfig{
        Enabled:             true,
        LatencyThreshold:    time.Millisecond * 100,
        ThroughputThreshold: 1000.0,
        SampleRate:          0.1,
        EnableProfiling:     true,
        
        ResourceMonitoring: ResourceMonitoringConfig{
            Enabled:          true,
            CPUThreshold:     80.0,
            MemoryThreshold:  85.0,
            DiskThreshold:    90.0,
            SampleInterval:   time.Second * 10,
        },
    },
    
    IO: IOMonitoringConfig{
        Enabled:              true,
        DiskUsageThreshold:   80.0,
        IOPSThreshold:        1000,
        LatencyThreshold:     time.Millisecond * 100,
        SampleInterval:       time.Second * 10,
        
        FileMonitoring: FileMonitoringConfig{
            Enabled: true,
            WatchFiles: []string{
                "/var/log/app.log",
                "/var/log/error.log",
            },
        },
    },
}

// 应用配置
monitoringManager := metrics.NewMonitoringManager(config)
if err := monitoringManager.Start(); err != nil {
    log.Fatal("启动监控失败:", err)
}
defer monitoringManager.Stop()
```

## 🚨 告警系统

### 告警管理器

```go
// 创建告警管理器
alertManager := metrics.NewAlertManager()

// 配置告警规则
rules := []AlertRule{
    {
        Name:      "high_memory_usage",
        Condition: func(metrics map[string]interface{}) bool {
            if usage, ok := metrics["memory.usage"].(float64); ok {
                return usage > 90.0
            }
            return false
        },
        Severity: "critical",
        Message:  "内存使用率超过90%",
        Actions: []AlertAction{
            &WebhookAction{URL: "http://alert:8080/webhook"},
            &EmailAction{To: []string{"admin@company.com"}},
        },
    },
    
    {
        Name:      "high_latency",
        Condition: func(metrics map[string]interface{}) bool {
            if latency, ok := metrics["performance.avg_latency"].(time.Duration); ok {
                return latency > time.Millisecond*500
            }
            return false
        },
        Severity: "warning",
        Message:  "平均延迟超过500ms",
        Actions: []AlertAction{
            &SlackAction{Channel: "#ops", Message: "高延迟告警"},
        },
    },
}

// 添加告警规则
for _, rule := range rules {
    alertManager.AddRule(rule)
}

// 启动告警管理器
alertManager.Start()
defer alertManager.Stop()
```

### 自定义告警动作

```go
// 实现自定义告警动作
type CustomAlertAction struct {
    APIEndpoint string
    APIKey      string
}

func (a *CustomAlertAction) Execute(alert *Alert) error {
    payload := map[string]interface{}{
        "rule_name": alert.RuleName,
        "severity":  alert.Severity,
        "message":   alert.Message,
        "timestamp": alert.Timestamp,
        "metrics":   alert.Metrics,
    }
    
    jsonData, _ := json.Marshal(payload)
    
    req, _ := http.NewRequest("POST", a.APIEndpoint, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+a.APIKey)
    
    client := &http.Client{Timeout: time.Second * 30}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("alert API returned status %d", resp.StatusCode)
    }
    
    return nil
}

// 使用自定义告警动作
customAction := &CustomAlertAction{
    APIEndpoint: "http://custom-alert-system:8080/alerts",
    APIKey:      "your-api-key",
}

rule := AlertRule{
    Name:      "custom_alert",
    Condition: customCondition,
    Severity:  "warning",
    Message:   "自定义告警",
    Actions:   []AlertAction{customAction},
}

alertManager.AddRule(rule)
```

### 告警抑制和静默

```go
// 配置告警抑制
alertManager.AddSuppressionRule(&SuppressionRule{
    Name: "maintenance_window",
    Condition: func(alert *Alert) bool {
        // 在维护窗口期间抑制所有告警
        now := time.Now()
        maintenanceStart := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
        maintenanceEnd := maintenanceStart.Add(time.Hour * 2)
        return now.After(maintenanceStart) && now.Before(maintenanceEnd)
    },
    Duration: time.Hour * 2,
})

// 临时静默特定告警
alertManager.SilenceAlert("high_memory_usage", time.Minute*30) // 静默30分钟

// 取消静默
alertManager.UnsilenceAlert("high_memory_usage")

// 获取告警状态
alertStatus := alertManager.GetAlertStatus()
for ruleName, status := range alertStatus {
    fmt.Printf("告警规则 %s: 状态=%s, 触发次数=%d, 最后触发=%s\n",
        ruleName, status.State, status.TriggerCount, 
        status.LastTrigger.Format("2006-01-02 15:04:05"))
}
```

---

更多监控相关信息请参考：

- [📊 性能详解](PERFORMANCE.md) - 详细性能分析和优化
- [🔧 配置指南](CONFIGURATION.md) - 监控配置详解
- [📚 使用指南](USAGE.md) - 完整使用指南
- [🎯 Context使用指南](CONTEXT_USAGE.md) - 分布式监控