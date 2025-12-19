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

### 📖 官方文档
- [🏠 项目主页](https://github.com/kamalyes/go-logger)
- [📖 API 文档](https://pkg.go.dev/github.com/kamalyes/go-logger)
- [📊 代码覆盖率](https://codecov.io/gh/kamalyes/go-logger)

### 📋 技术文档
- 📊 **[性能详解](docs/PERFORMANCE.md)** - 深入了解性能优化技术和基准测试结果
- 🔄 **[迁移指南](docs/MIGRATION.md)** - 从其他日志库迁移的完整指南
- 🎯 **[Context使用指南](docs/CONTEXT_USAGE.md)** - 分布式系统上下文管理和链路追踪
- 🔌 **[自定义上下文提取器](docs/CUSTOM_CONTEXT_EXTRACTOR.md)** - 灵活提取和自定义上下文信息
- ↩️ **[返回错误日志](docs/RETURN_ERROR.md)** - 简化错误处理的日志方法
- 🎨 **[Console 风格日志](docs/CONSOLE_USAGE.md)** - JavaScript Console 风格的分组、表格和计时器功能
- 📝 **[更新日志](./CHANGELOG.md)** - 版本更新和功能变更记录
- 🔧 **[配置指南](docs/CONFIGURATION.md)** - 完整配置选项和最佳实践
- 🧩 **[适配器系统](docs/ADAPTERS.md)** - 适配器完整指南和自定义开发
- 📊 **[监控系统](docs/MONITORING.md)** - 内存监控、性能分析和告警系统
- 🎨 **[格式化器](docs/FORMATTERS.md)** - 日志格式化器详解和自定义开发

### 🔗 代码资源
- 📋 **[示例代码](examples/README.md)** - 丰富的使用示例和最佳实践
- 🧪 **[基准测试](benchmark_test.go)** - 性能测试和对比分析
- ⚡ **[极速日志器](ultra_fast_logger.go)** - 极致性能实现源码

### 💬 社区支持
- [🐛 问题反馈](https://github.com/kamalyes/go-logger/issues)
- [💬 讨论区](https://github.com/kamalyes/go-logger/discussions)

## 🚀 为什么选择 go-logger？

### ⚡ 极致性能 
- **🏆 业界领先**: 相比标准库 slog **快 7.7倍** (75.8ns vs 585.2ns)
- **💾 内存优化**: **83% 内存减少** (144B → 24B)，**50% 分配减少** (2 → 1 allocs)
- **🔧 分层设计**: 三层性能架构满足不同性能需求
- **📊 零开销**: 级别过滤接近零性能开销

### 核心功能
- **📊 内存监控系统**：实时监控内存使用、GC性能、堆分析，支持内存泄漏检测
- **🔍 分布式追踪**：统一的Context服务架构，支持TraceID、SpanID、CorrelationID等多维度追踪
- **🔌 自定义上下文提取器**：灵活的上下文信息提取机制，支持完全自定义链路追踪字段
- **🎯 多级日志系统**：支持24种日志级别，从TRACE到PROFILING，满足不同场景需求
- **📈 性能监控**：实时统计操作性能、延迟分析、吞吐量监控
- **⚡ 架构重构**：Context管理代码减少88%，从1059行优化到128行，性能显著提升

### 企业级功能
- **🛡️ 内存安全**：智能内存管理、GC优化、内存压力检测与自动释放
- **📊 统计分析**：详细的运行时统计、性能指标收集、趋势分析
- **🔧 配置管理**：细粒度配置系统，支持动态配置更新
- **⚙️ 适配器模式**：支持多种输出适配器，灵活扩展输出目标
- **🧪 完善测试**：基于测试套件的全面测试，覆盖率90%+

### 🔌 自定义上下文提取器

支持灵活提取和自定义上下文信息，满足不同场景需求：

**核心能力**：
- 🎯 **预定义提取器**: SimpleTraceIDExtractor、SimpleRequestIDExtractor、NoOpContextExtractor
- 🔧 **自定义字段**: CustomFieldExtractor - 从 context 或 gRPC metadata 提取任意字段
- 🔗 **链式组合**: ChainContextExtractors - 组合多个提取器
- 🏗️ **构建器模式**: ContextExtractorBuilder - 流式 API 构建复杂提取器
- ⚡ **完全自定义**: 支持自定义 ContextExtractor 函数

**性能表现**: NoOp (137ns) | 默认 (466ns) | 链式 (430ns)

**适用场景**: 微服务追踪 | 多租户系统 | API 网关 | 分布式链路追踪

📖 **[查看完整文档和示例 →](docs/CUSTOM_CONTEXT_EXTRACTOR.md)**

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
ultraLogger := logger.NewUltraFast()

// 或使用完整配置
config := logger.DefaultConfig()
config.Level = logger.INFO
config.Colorful = false
config.ShowCaller = false
ultraLogger = logger.NewUltraFastLogger(config)

// ⚡ 高性能 - 优化版标准Logger  
optimizedLogger := logger.NewOptimized()

// 🛡️ 全功能 - 企业级Logger (默认)
fullLogger := logger.New()

// 或使用完整配置
enterpriseConfig := logger.DefaultConfig()
enterpriseConfig.Level = logger.INFO
enterpriseConfig.ShowCaller = true
enterpriseConfig.Colorful = true
fullLogger = logger.NewLogger(enterpriseConfig)
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

📖 **[查看详细性能分析 →](docs/PERFORMANCE.md)**

## 🏗️ 模块化架构
```
go-logger/
├── config/              # 配置管理模块
│   ├── base.go          # 基础配置
│   ├── adapter.go       # 适配器配置
│   ├── output.go        # 输出配置
│   └── level.go         # 日志级别配置
├── context_service.go   # 统一上下文服务（新架构核心）
├── level/               # 日志级别管理
│   ├── constants.go     # 级别常量定义
│   └── manager.go       # 级别管理器
├── metrics/             # 监控指标模块
│   ├── stats.go         # 统计收集
│   ├── performance.go   # 性能监控
│   └── memory.go        # 内存监控
├── docs/                # 文档目录
│   ├── CONTEXT_USAGE.md # Context使用指南
│   ├── PERFORMANCE.md   # 性能详解
│   └── MIGRATION.md     # 迁移指南
└── examples/            # 示例代码
```

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
)

func main() {
    // 🏆 极致性能版本 (推荐高并发场景)
    ultraLogger := logger.NewUltraFast()
    ultraLogger.Info("High performance logging")
    
    // 结构化日志 - 键值对方式
    ultraLogger.InfoKV("High performance with fields", "key", "value")
    
    // 🎯 结构化日志 - 对象方式 (自动解析)
    type User struct {
        ID    int    `json:"id"`
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    user := User{ID: 1001, Name: "张三", Email: "user@example.com"}
    
    // 直接传递对象，自动解析为键值对
    ultraLogger.InfoKV("用户登录", user)
    
    // 也支持 map
    data := map[string]interface{}{
        "request_id": "req-123",
        "method":     "POST",
        "status":     200,
    }
    ultraLogger.InfoKV("API 请求", data)
    
    // ⚡ 优化版标准Logger
    optimizedLogger := logger.NewOptimized()
    optimizedLogger.Info("Optimized logging with features")
    
    // 🛡️ 全功能企业版 (默认)
    fullLogger := logger.New()
    fullLogger.Info("Full featured logging")
    
    // 🎨 Console 风格日志 - JavaScript 风格的分组、表格和计时器
    cg := fullLogger.NewConsoleGroup()
    
    // 📊 分组日志 - 组织相关日志输出
    cg.Group("🚀 应用启动流程")
    cg.Info("开始初始化...")
    
    // 📋 表格展示 - 美化配置信息
    config := map[string]interface{}{
        "环境":   "生产环境",
        "端口":   8080,
        "调试模式": false,
    }
    cg.Table(config)
    // 输出:
    //   ┌──────────┬────────────┐
    //   │ Key      │ Value      │
    //   ├──────────┼────────────┤
    //   │ 环境      │ 生产环境   │
    //   │ 端口      │ 8080       │
    //   │ 调试模式   │ false     │
    //   └──────────┴────────────┘
    
    // ⏱️  计时器 - 测量操作耗时
    timer := cg.Time("数据库连接")
    // ... 执行数据库连接 ...
    timer.End() // 输出: ⏱️  数据库连接: 123.45ms
    
    cg.Info("✅ 启动完成")
    cg.GroupEnd()
    
    // 📖 查看完整 Console 功能文档: docs/CONSOLE_USAGE.md
    
    // 🎯 使用现有的Context ID管理
    ctx := context.Background()
    
    // 直接使用日志记录（结构化字段通过WithField添加）
    fullLogger.WithField("trace_id", "trace-123").
               WithField("user_id", "user-456").
               Info("带上下文的日志")
    
    // 🔌 自定义上下文提取器 (灵活提取链路信息)
    ctx = context.WithValue(ctx, "trace_id", "trace-12345")
    ctx = context.WithValue(ctx, "request_id", "req-67890")
    
    // 使用默认提取器
    ultraLogger.InfoContext(ctx, "用户登录成功")
    // 输出: [TraceID=trace-12345 RequestID=req-67890] 用户登录成功
    
    // 自定义提取器（详见文档）
    extractor := logger.NewContextExtractorBuilder().
        AddTraceID().
        AddRequestID().
        AddContextValue("user_id", "User").
        Build()
    ultraLogger.SetContextExtractor(extractor)
}
```

### 🎨 Console 风格日志功能

类似 JavaScript `console` 的日志分组、表格和计时器功能,让日志输出更加结构化和易读。

```go
log := logger.NewLogger(logger.DefaultConfig())
cg := log.NewConsoleGroup()

// 📊 日志分组 - 组织相关日志
cg.Group("🌐 API 请求处理")
cg.Info("接收到请求: GET /api/users")

    // 嵌套分组
    cg.Group("参数验证")
    cg.Info("验证通过")
    cg.GroupEnd()

// 📋 表格展示 - 结构化数据可视化
users := []map[string]interface{}{
    {"ID": 1, "姓名": "张三", "年龄": 25, "状态": "Active"},
    {"ID": 2, "姓名": "李四", "年龄": 30, "状态": "Active"},
}
cg.Table(users)
// 输出美观的表格:
//   ┌────┬──────┬──────┬────────┐
//   │ ID │ 姓名  │ 年龄  │ 状态    │
//   ├────┼──────┼──────┼────────┤
//   │ 1  │ 张三  │ 25   │ Active │
//   │ 2  │ 李四  │ 30   │ Active │
//   └────┴──────┴──────┴────────┘

// ⏱️  性能计时 - 测量操作耗时
timer := cg.Time("数据库查询")
// ... 执行数据库操作 ...
timer.End()  // 输出: ⏱️  数据库查询: 123.45ms

// 中间检查点
timer2 := cg.Time("文件处理")
// ... 执行部分操作 ...
timer2.Log("处理 50%")  // 输出: ⏱️  文件处理: 50.12ms - 处理 50%
// ... 继续操作 ...
timer2.End()  // 输出: ⏱️  文件处理: 102.34ms

cg.Info("✅ 请求处理完成")
cg.GroupEnd()

// 🎯 折叠分组 - 隐藏详细日志（仅显示 ERROR/FATAL）
cg.GroupCollapsed("调试信息")
cg.Debug("这条不会显示")
cg.Info("这条也不会显示")
cg.Error("但错误日志会显示")  // ❌ 会显示
cg.GroupEnd()

// 🌐 全局便捷方法 - 不需要 ConsoleGroup
logger.Group("全局分组")
logger.Info("这是全局分组内的日志")
logger.Table(map[string]string{"key": "value"})
logger.GroupEnd()

timer := logger.Time("全局计时器")
// ... 操作 ...
timer.End()
```

**主要特性**：
- 🎯 **日志分组**: 
  - `Group(label, ...args)` - 开始新分组
  - `GroupCollapsed(label, ...args)` - 开始折叠分组（仅显示 ERROR/FATAL）
  - `GroupEnd()` - 结束当前分组
  - 支持无限层级嵌套，自动缩进

- 📊 **表格渲染**: 
  - `Table(data)` - 智能表格渲染
  - 支持格式: `[]map[string]interface{}`, `map[string]interface{}`, `[][]string`, `[]string`
  - 自动对齐、美化边框、智能列宽

- ⏱️  **计时器**: 
  - `Time(label)` - 开始计时，返回 Timer 对象
  - `Timer.End()` - 结束计时并输出总耗时
  - `Timer.Log(message)` - 输出中间检查点
  - 智能时间格式化 (ms/s/m)

- 🔄 **Context 集成**: 
  - `InfoContext(ctx, ...)` - 带上下文的 Info 日志
  - `DebugContext(ctx, ...)` - 带上下文的 Debug 日志
  - `WarnContext(ctx, ...)` - 带上下文的 Warn 日志
  - `ErrorContext(ctx, ...)` - 带上下文的 Error 日志
  - 在分组内使用，自动继承缩进

- 🌐 **全局方法**: 
  - `logger.Group()` / `logger.GroupEnd()`
  - `logger.Table(data)`
  - `logger.Time(label)`
  - 无需创建 ConsoleGroup，直接使用

**适用场景**：
- 🚀 应用启动流程展示
- 📊 批量数据处理进度
- 🔍 复杂业务流程追踪
- ⚡ 性能瓶颈分析
- 🐛 调试信息结构化输出

**兼容性**：
- ✅ 所有日志器均支持: `Logger`, `UltraFastLogger`, `StandardAdapter`
- ✅ `UltraFastLogger` 通过内部委托实现完整 Console 功能
- ✅ `EmptyLogger` 提供空实现，不影响生产环境性能

📖 **[查看完整 Console 使用文档和高级示例 →](docs/CONSOLE_USAGE.md)**


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

最新性能测试结果：

```
BenchmarkUltraFastLogger-8       157894737     7.56 ns/op     0 B/op     0 allocs/op
BenchmarkStandardLogger-8         52631578    22.85 ns/op     8 B/op     1 allocs/op
BenchmarkMemoryMonitor-8           9803921   122.4 ns/op    48 B/op     2 allocs/op
```

详细性能分析请参考 [性能文档](docs/PERFORMANCE.md)。

## ⭐ Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=kamalyes/go-logger&type=Date)](https://star-history.com/#kamalyes/go-logger&Date)

## 许可证

该项目使用 MIT 许可证，详见 [LICENSE](LICENSE) 文件