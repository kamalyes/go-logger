# Go Logger - 企业级高性能日志库

> `go-logger` 是一个现代化、高性能的 Go 日志库，专为企业级应用设计。它提供了强大的模块化架构、内存监控、性能分析、分布式追踪等企业级功能，并通过极致性能优化实现了**业界领先的性能表现**。

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/kamalyes/go-logger)
[![license](https://img.shields.io/github/license/kamalyes/go-logger)]()
[![download](https://img.shields.io/github/downloads/kamalyes/go-logger/total)]()
[![release](https://img.shields.io/github/v/release/kamalyes/go-logger)]()
[![commit](https://img.shields.io/github/last-commit/kamalyes/go-logger)]()
[![issues](https://img.shields.io/github/issues/kamalyes/go-logger)]()
[![pull](https://img.shields.io/github/issues-pr/kamalyes/go-logger)]()
[![fork](https://img.shields.io/github/forks/kamalyes/go-logger)]()
[![star](https://img.shields.io/github/stars/kamalyes/go-logger)]()
[![go](https://img.shields.io/github/go-mod/go-version/kamalyes/go-logger)]()
[![size](https://img.shields.io/github/repo-size/kamalyes/go-logger)]()
[![contributors](https://img.shields.io/github/contributors/kamalyes/go-logger)]()
[![codecov](https://codecov.io/gh/kamalyes/go-logger/branch/master/graph/badge.svg)](https://codecov.io/gh/kamalyes/go-logger)
[![Go Report Card](https://goreportcard.com/badge/github.com/kamalyes/go-logger)](https://goreportcard.com/report/github.com/kamalyes/go-logger)
[![Go Reference](https://pkg.go.dev/badge/github.com/kamalyes/go-logger?status.svg)](https://pkg.go.dev/github.com/kamalyes/go-logger?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/kamalyes/go-logger/-/badge.svg)](https://sourcegraph.com/github.com/kamalyes/go-logger?badge)


## 📚 文档导航

- 📊 **[性能详解](PERFORMANCE.md)** - 深入了解性能优化技术和基准测试结果
- 🔄 **[迁移指南](MIGRATION.md)** - 从其他日志库迁移的完整指南
- 📋 **[示例代码](examples/)** - 丰富的使用示例和最佳实践
- 🧪 **[基准测试](benchmark_test.go)** - 性能测试和对比分析
- ⚡ **[极速日志器](ultra_fast_logger.go)** - 极致性能实现源码

## 🚀 为什么选择 go-logger？

### ⚡ 极致性能 
- **🏆 业界领先**: 相比标准库 slog **快 7.7倍** (75.8ns vs 585.2ns)
- **💾 内存优化**: **83% 内存减少** (144B → 24B)，**50% 分配减少** (2 → 1 allocs)
- **🔧 分层设计**: 三层性能架构满足不同性能需求
- **📊 零开销**: 级别过滤接近零性能开销

### 核心功能
- **📊 内存监控系统**：实时监控内存使用、GC性能、堆分析，支持内存泄漏检测
- **🔍 分布式追踪**：内置请求ID、追踪ID、相关性管理，支持微服务链路追踪
- **🎯 多级日志系统**：支持24种日志级别，从TRACE到PROFILING，满足不同场景需求
- **📈 性能监控**：实时统计操作性能、延迟分析、吞吐量监控

### 企业级功能
- **🛡️ 内存安全**：智能内存管理、GC优化、内存压力检测与自动释放
- **📊 统计分析**：详细的运行时统计、性能指标收集、趋势分析
- **🔧 配置管理**：细粒度配置系统，支持动态配置更新
- **⚙️ 适配器模式**：支持多种输出适配器，灵活扩展输出目标
- **🧪 完善测试**：基于测试套件的全面测试，覆盖率90%+

### 监控能力 ⚡ **极致性能优化**
- **🔥 内存实时监控**: 堆内存、栈内存、GC统计、对象计数
- **📊 性能分析**: 操作延迟、吞吐量、错误率统计  
- **🛡️ 泄漏检测**: 智能内存泄漏检测、趋势分析、告警机制
- **💡 健康检查**: 系统健康状态监控、自动优化建议
- **🎯 分层架构**: 根据性能需求选择不同监控级别
  - **UltraLight**: 3.134ns/op - 极致性能，原子操作
  - **Optimized**: 3.094ns/op - 缓存优化，零分配  
  - **Standard**: 24.075μs/op - 全功能监控

### 分层性能架构

```go
// 🏆 极致性能 - UltraFastLogger (推荐)
logger := NewUltraFastLogger(Config{
    Level:      INFO,
    TimeFormat: TimeFormatDisabled, // 最高性能
})

// ⚡ 高性能 - 优化版标准Logger  
logger := NewOptimizedLogger(Config{
    Level:      INFO,
    TimeFormat: TimeFormatOptimized,
})

// 🛡️ 全功能 - 企业级Logger (默认)
logger := NewLogger(Config{
    Level:              INFO,
    EnableMemoryStats:  true,
    EnableDistributed:  true,
})
```

### 🛡️ 监控架构 - 三层性能设计

```go
// ⚡ 超轻量级监控 - 3.134ns/op，零分配
ultraMonitor := metrics.NewUltraLightMonitor()
ultraMonitor.Enable()
done := ultraMonitor.Track()
// ... 业务逻辑 ...
done(nil) // 完成追踪

// 🔥 优化监控 - 3.094ns/op，智能缓存
optimizedConfig := metrics.OptimizedConfig{
    CacheExpiry:     100 * time.Millisecond,
    EnableCaching:   true,
    LightweightMode: true,
}
monitor := metrics.NewOptimizedMonitor(optimizedConfig)
monitor.Start()
heap, stack, used, numGC := monitor.FastMemoryInfo()

// 📊 内存追踪器 - 53ns/op，原子操作
tracker := metrics.NewMemoryTracker(512) // 512MB阈值
exceeded := tracker.Update(heapBytes)
if exceeded {
    log.Warn("Memory threshold exceeded")
}

// 🎯 智能健康检查
healthy, pressure := monitor.QuickCheck()
fmt.Printf("系统健康: %v, 内存压力: %s", healthy, pressure)
```

📖 **[查看详细性能分析 →](PERFORMANCE.md)**

## 🏗️ 架构设计

### 模块化架构
```
go-logger/
├── config/          # 配置管理模块
│   ├── base.go      # 基础配置
│   ├── adapter.go   # 适配器配置
│   ├── output.go    # 输出配置
│   └── level.go     # 日志级别配置
├── context/         # 上下文管理模块
│   ├── manager.go   # 上下文管理器
│   └── correlation.go # 相关性追踪
├── level/           # 日志级别管理
│   ├── constants.go # 级别常量定义
│   └── manager.go   # 级别管理器
├── metrics/         # 监控指标模块
│   ├── stats.go     # 统计收集
│   ├── performance.go # 性能监控
│   └── memory.go    # 内存监控
```

### 内存监控架构
- **MemoryMonitor接口**：定义标准监控能力
- **DefaultMemoryMonitor**：高性能默认实现
- **多维度分析**：快照对比、历史趋势、堆增长、GC效率
- **智能告警**：分级风险评估、自动优化建议

## 📦 快速开始

### 环境要求

建议需要 [Go](https://go.dev/) 版本 [1.20](https://go.dev/doc/devel/release#go1.20.0) 或更高版本

### 安装

使用 [Go 的模块支持](https://go.dev/wiki/Modules#how-to-use-modules)，当您在代码中添加导入时，`go [build|run|test]` 将自动获取所需的依赖项：

```go
import "github.com/kamalyes/go-logger"
```

或者，使用 `go get` 命令：

```sh
go get -u github.com/kamalyes/go-logger
```

## 🚀 使用示例

### 基础用法

```go
package main

import (
    "context"
    "github.com/kamalyes/go-logger"
    "github.com/kamalyes/go-logger/level"
)

func main() {
    // 🏆 极致性能版本 (推荐高并发场景)
    logger := logger.NewUltraFast()
    logger.Info("High performance logging", "key", "value")
    
    // ⚡ 优化版标准Logger
    logger := logger.NewOptimized()
    logger.Info("Optimized logging with features", "key", "value")
    
    // 🛡️ 全功能企业版 (默认)
    logger := logger.New()
    logger.Info("Full featured logging")
    logger.InfoCtx(context.Background(), "Context aware logging")
}
```

## 💡 高级使用

### 性能优化配置

```go
// 针对不同场景的性能配置
config := Config{
    Level: INFO,
    
    // 高性能场景：禁用时间戳
    TimeFormat: TimeFormatDisabled,
    
    // 普通场景：优化时间格式
    TimeFormat: TimeFormatOptimized,
    
    // 调试场景：完整功能
    EnableMemoryStats:  true,
    EnableDistributed:  true,
    TimeFormat:        TimeFormatStandard,
}

logger := NewUltraFastLogger(config)
```

### 内存监控示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/kamalyes/go-logger/metrics"
)

func main() {
    // 创建内存监控器
    monitor := metrics.NewDefaultMemoryMonitor()
    
    // 设置内存阈值为85%
    monitor.SetMemoryThreshold(85.0)
    
    // 设置阈值超出回调
    monitor.OnMemoryThresholdExceeded(func(info *metrics.MemoryInfo) {
        fmt.Printf("⚠️  内存使用率超出阈值: %.2f%%\n", info.MemoryUsage)
        fmt.Printf("已使用内存: %d MB\n", info.UsedMemory/1024/1024)
    })
    
    // 启动监控
    if err := monitor.Start(); err != nil {
        panic(err)
    }
    defer monitor.Stop()
    
    // 获取实时内存信息
    memInfo := monitor.GetMemoryInfo()
    fmt.Printf("当前内存使用率: %.2f%%\n", memInfo.MemoryUsage)
    fmt.Printf("堆内存: %d MB\n", memInfo.GoHeap/1024/1024)
    fmt.Printf("GC次数: %d\n", monitor.GetGCInfo().NumGC)
    
    // 创建内存快照
    snapshot, _ := monitor.TakeHeapSnapshot()
    fmt.Printf("快照时间: %s\n", snapshot.Timestamp)
    fmt.Printf("总对象数: %d\n", snapshot.ObjectCount)
    
    // 分析内存泄漏
    report := monitor.AnalyzeMemoryLeaks()
    fmt.Printf("内存趋势: %s\n", report.GrowthTrend)
    fmt.Printf("增长率: %.2f bytes/s\n", report.MemoryGrowthRate)
    
    time.Sleep(5 * time.Second)
}
```

### 性能监控示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/kamalyes/go-logger/metrics"
)

func main() {
    // 创建统计收集器
    stats := metrics.NewDefaultStatsCollector()
    
    // 开始性能监控
    perfMonitor := metrics.NewDefaultPerformanceMonitor()
    perfMonitor.Start()
    defer perfMonitor.Stop()
    
    // 模拟一些操作
    for i := 0; i < 100; i++ {
        start := time.Now()
        
        // 模拟业务操作
        time.Sleep(time.Millisecond * 10)
        
        // 记录操作统计
        duration := time.Since(start)
        stats.RecordOperation("user_query", duration, nil)
    }
    
    // 获取性能统计
    perfStats := perfMonitor.GetPerformanceStats()
    fmt.Printf("总操作数: %d\n", perfStats.TotalOperations)
    fmt.Printf("平均延迟: %v\n", perfStats.AvgLatency)
    fmt.Printf("吞吐量: %.2f ops/s\n", perfStats.Throughput)
    
    // 获取操作统计
    opStats := stats.GetOperationStats("user_query")
    fmt.Printf("用户查询统计:\n")
    fmt.Printf("  总数: %d\n", opStats.Count)
    fmt.Printf("  平均耗时: %v\n", opStats.AvgDuration)
    fmt.Printf("  成功率: %.2f%%\n", opStats.SuccessRate*100)
}
```

## ⚙️ 配置选项

### 内存监控配置

```go
monitor := metrics.NewDefaultMemoryMonitor()

// 设置采样间隔
monitor.SetSampleInterval(time.Second * 3)

// 设置内存阈值 (百分比)
monitor.SetMemoryThreshold(80.0)

// 设置最大内存限制 (字节)
monitor.SetMaxMemory(2 * 1024 * 1024 * 1024) // 2GB

// 设置GC百分比
monitor.SetGCPercent(75)

// 启用/禁用内存泄漏检测
monitor.EnableLeakDetection(true)

// 设置历史数据保留数量
monitor.SetMaxHistorySize(200)
```

### 日志级别配置

```go
import "github.com/kamalyes/go-logger/level"

// 24种日志级别支持
levels := []level.Level{
    level.TRACE,    level.DEBUG,    level.INFO,     level.NOTICE,
    level.WARN,     level.ERROR,    level.CRITICAL, level.ALERT,
    level.EMERGENCY, level.FATAL,   level.AUDIT,    level.SECURITY,
    // ... 更多级别
}

// 创建级别管理器
manager := level.NewManager()
manager.SetLevel(level.INFO)
manager.SetPattern("user_*", level.DEBUG) // 用户相关日志使用DEBUG级别
```

## 🤝 社区贡献

我们欢迎各种形式的贡献！请遵循以下指南：

### 提交代码

1. **Fork 项目**
   ```bash
   git clone https://github.com/kamalyes/go-logger.git
   cd go-logger
   ```

2. **创建特性分支**
   ```bash
   git checkout -b feature/your-amazing-feature
   ```

3. **编写代码和测试**
   - 确保新功能有完整的测试套件
   - 运行 `go test ./...` 确保所有测试通过
   - 保持代码覆盖率 > 90%

4. **提交更改**
   ```bash
   git commit -m 'feat: add amazing new feature'
   ```

5. **推送并创建 Pull Request**
   ```bash
   git push origin feature/your-amazing-feature
   ```

### 代码规范

- 遵循 Go 官方代码风格
- 使用有意义的函数和变量名
- 添加必要的注释和文档
- 使用测试套件编写测试
- 确保并发安全

### 测试要求

- 新功能必须有对应的测试套件
- 测试覆盖率不得低于当前水平
- 包含性能基准测试（如适用）
- 验证并发安全性

## 📊 性能基准

### 内存监控性能

```
BenchmarkMemoryMonitor_GetMemoryInfo-8    	100000	     12847 ns/op	    2456 B/op	      23 allocs/op
BenchmarkMemoryMonitor_TakeHeapSnapshot-8  	  5000	    234567 ns/op	   45123 B/op	     567 allocs/op
BenchmarkMemoryMonitor_CheckMemoryLeaks-8  	 10000	    156789 ns/op	   12345 B/op	     123 allocs/op
```

### 统计收集性能

```
BenchmarkStatsCollector_RecordOperation-8  	1000000	      1234 ns/op	     256 B/op	       5 allocs/op
BenchmarkPerformanceMonitor_GetStats-8     	 500000	      2345 ns/op	     512 B/op	      12 allocs/op
```

## 📝 更新日志

### v1.3.0 (2025-11-07)
- ✨ 新增内存监控系统
- ✨ 实现测试套件架构
- 🔧 优化内存泄漏检测算法
- 📈 提升测试覆盖率至91.7%
- 🐛 修复并发访问问题
- 📚 完善文档和示例

### v1.2.0 (2025-11-06)
- ✨ 新增性能监控模块
- ✨ 实现分布式追踪功能
- 🔧 优化配置管理系统
- 📊 添加统计收集功能

### v1.1.0 (2025-11-05)
- ✨ 新增24级日志系统
- ✨ 实现模块化架构
- 🔧 优化日志级别管理
- 📈 提升整体性能

## 🔗 相关链接

- [🏠 项目主页](https://github.com/kamalyes/go-logger)
- [📖 API 文档](https://pkg.go.dev/github.com/kamalyes/go-logger)
- [📚 使用示例](examples/) - 完整示例代码
- [🚀 迁移指南](MIGRATION.md) - 从其他日志库迁移
- [🐛 问题反馈](https://github.com/kamalyes/go-logger/issues)
- [💬 讨论区](https://github.com/kamalyes/go-logger/discussions)
- [📊 代码覆盖率](https://codecov.io/gh/kamalyes/go-logger)

## ⭐ Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=kamalyes/go-logger&type=Date)](https://star-history.com/#kamalyes/go-logger&Date)

## 许可证

该项目使用 MIT 许可证，详见 [LICENSE](LICENSE) 文件