# Go Logger - 企业级高性能日志库

> `go-logger` 是一个现代化、高性能的 Go 日志库，专为企业级应用设计。它提供了强大的结构化日志、基于 OpenTelemetry 的分布式追踪、Console 风格日志等企业级功能，并通过对象池、零拷贝、原子操作等性能优化实现了出色的性能表现。

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

### 💬 社区支持

- [🐛 问题反馈](https://github.com/kamalyes/go-logger/issues)
- [💬 讨论区](https://github.com/kamalyes/go-logger/discussions)

## 🚀 为什么选择 go-logger？

### ⚡ 核心特性

- **📊 结构化日志**: 支持键值对、字段映射、对象自动解析、多行日志等多种结构化方式
- **🔍 分布式追踪**: 基于 OpenTelemetry 的链路追踪，`trace_id` 单一真相源，自动从 OTel span 提取
- **🎨 Console 风格**: JavaScript Console 风格的分组、表格、计时器功能
- **🔌 灵活扩展**: 多种写入器适配器（Console、File、Rotate、Buffered、Multi、AsyncBatch）、自定义上下文提取器、格式化器、钩子、中间件
- **⚡ 高性能**: 对象池、零拷贝、原子操作、预计算常量、数组索引替 map 等性能优化
- **🛡️ 并发安全**: 完善的并发控制，适合高并发场景
- **🎯 多输出格式**: 支持 Text / JSON / XML / CSV 四种输出格式

### 核心功能

#### 📝 多种日志方式

- **Printf 风格**: `Info(format, args...)` / `Infof(format, args...)`
- **纯文本**: `InfoMsg(msg)` - 单个消息，无格式化开销
- **键值对**: `InfoKV(msg, key, value, ...)` - 结构化日志
- **对象解析**: `InfoKV(msg, struct)` - 单个结构体参数自动解析为键值对（通过 `convert.ParseObjectToMap`）
- **字段映射**: `InfoWithFields(msg, fields)` - map 方式
- **多行日志**: `InfoLines(line1, line2, ...)` - 自动处理多行
- **返回错误**: `InfoReturn(format, args...)` - 记录日志并返回 error
- **标准库兼容**: `Print` / `Printf` / `Println` 兼容标准 `log` 包

#### 🔍 分布式追踪

- **OpenTelemetry 单一真相源**: `trace_id` 优先从 OTel span 提取（gRPC otelgrpc StatsHandler / HTTP Tracing 中间件创建 span）
- **多来源 fallback**: OTel span → `ctx.Value` → gRPC metadata
- **gRPC 链路传播**: `InjectTraceToOutgoing` / `ExtractTraceFromIncoming` / `RestoreTraceFromIncoming`
- **便捷辅助函数**: `ExtractTraceID` / `ContextWithTraceID`
- **自定义提取器**: `SetContextExtractor` / `WithContextExtractor` / `WithContextKeys`

#### 🎨 Console 风格日志

- **日志分组**: 支持嵌套分组和折叠分组（折叠时仅显示 ERROR/FATAL）
- **表格渲染**: 自动对齐、美化边框、智能列宽、东亚宽度计算
- **计时器**: 支持中间检查点（`Log`）、自动格式化（ns/μs/ms/s）、自动清理过期计时器
- **Context 集成**: 分组内支持带上下文的日志（`InfoContext` 等）
- **输出隔离**: Console 展示内容直接写入 output，不走 logger 的 JSON/文本格式化，避免污染结构化日志流

#### 🎯 多级日志系统

- **26 种日志级别**（按 6 个类别组织）：
  - 基础级别（7个）：TRACE(-1) → DEBUG(0) → INFO(1) → WARN(2) → ERROR(3) → FATAL(4) → OFF(99)
  - 系统级别（3个）：SYSTEM(100) → KERNEL(101) → DRIVER(102)
  - 应用级别（4个）：APPLICATION(103) → SERVICE(104) → COMPONENT(105) → MODULE(106)
  - 业务级别（4个）：BUSINESS(200) → TRANSACTION(201) → WORKFLOW(202) → PROCESS(203)
  - 安全级别（4个）：SECURITY(300) → AUDIT(301) → COMPLIANCE(302) → THREAT(303)
  - 性能级别（4个）：PERFORMANCE(400) → METRIC(401) → BENCHMARK(402) → PROFILING(403)
- **级别过滤**: 动态调整日志级别（原子操作，运行时并发安全）
- **YAML/JSON 反序列化**: `LogLevel` 实现 `UnmarshalYAML` / `UnmarshalText` / `MarshalText`，支持字符串与整数
- **特殊日志类型**: Success、Loading、Config、Start、Stop、Database、Network、Security、Cache、Environment
- **性能日志**: 自动分级（EXCELLENT/FAST/NORMAL/SLOW/VERY_SLOW）
- **进度日志**: 百分比和进度条显示
- **健康检查**: 服务健康状态监控
- **审计日志**: 用户操作审计追踪

### 企业级功能

- **🔧 灵活配置**: Builder 模式链式调用，支持运行时动态配置
- **⚙️ 多种写入器**: Console、File、Rotate、Buffered、Multi、AsyncBatch，全部采用 Options 模式
- **⚙️ 配置驱动创建**: `WriterConfig` + `CreateWriter` 支持通过配置统一创建 Writer
- **🎯 错误处理**: 返回 error 的日志方法（`*Return` / `*CtxReturn` / `*KVReturn`），简化错误处理流程
- **📊 统计分析**: `LoggerStats` 内置日志统计，`WriterStatsSnapshot` 收集写入器运行时指标
- **🧪 完善测试**: 全面的测试覆盖与基准测试（`logger_benchmark_test.go`、`writer_benchmark_test.go`、`async_batch_test.go` 等）
- **🔌 接口丰富**: 定义 `ILogger` / `IWriter` / `IAdapter` / `IManager` / `IFormatter` / `IHook` / `IMiddleware` 等接口，并兼容 Zap/Logrus/slog/Zerolog 框架接口

## 🏗️ 项目结构

```
go-logger/
├── logger.go                  # 核心日志实现（基础/结构化/上下文/特殊/性能日志）
├── types.go                   # 核心类型定义与 Builder 方法、全局 Logger
├── interfaces.go              # 接口定义（ILogger 及框架兼容接口）
├── level.go                   # 日志级别管理（26 种级别，6 个类别）
├── context.go                 # 上下文服务（OTel 链路追踪提取）
├── console.go                 # Console 风格日志（分组/表格）
├── timer.go                   # 计时器实现（含自动清理）
├── writer.go                  # 写入器实现（Console/File/Rotate/Buffered/Multi/AsyncBatch）
├── output.go                  # 输出类型与配置（WriterConfig + CreateWriter）
├── empty.go                   # 空实现（EmptyLogger/EmptyAdapter/EmptyWriter/EmptyHook）
├── *_test.go                  # 测试文件
├── logger_benchmark_test.go   # 日志基准测试
├── writer_benchmark_test.go   # 写入器基准测试
└── benchmark.sh               # 基准测试脚本
```

## 📦 快速开始

### 环境要求

需要 [Go](https://go.dev/) 版本 [1.25.0](https://go.dev/doc/devel/release) 或更高版本

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
 // 创建日志实例
 log := logger.NewLogger()

 // 基础日志（Printf 风格）
 log.Info("应用启动")
 log.Debug("调试信息: %s", "debug message")
 log.Warn("警告信息")
 log.Error("错误信息")

 // 结构化日志 - 键值对方式
 log.InfoKV("用户登录", "user_id", 1001, "username", "张三")

 // 🎯 结构化日志 - 对象方式 (自动解析)
 type User struct {
  ID    int    `json:"id"`
  Name  string `json:"name"`
  Email string `json:"email"`
 }
 user := User{ID: 1001, Name: "张三", Email: "user@example.com"}

 // 直接传递单个对象，自动解析为键值对
 log.InfoKV("用户信息", user)
 // 输出: 用户信息 {id: 1001, name: 张三, email: user@example.com}

 // 也支持 map
 data := map[string]interface{}{
  "request_id": "req-123",
  "method":     "POST",
  "status":     200,
 }
 log.InfoKV("API 请求", data)

 // 字段映射方式
 log.InfoWithFields("订单创建", map[string]interface{}{
  "order_id": "ORD-12345",
  "amount":   99.99,
  "user_id":  1001,
 })

 // 🎯 带上下文的日志（基于 OTel 链路追踪）
 ctx := logger.ContextWithTraceID(context.Background(), "trace-12345")
 log.InfoContext(ctx, "用户登录成功")
 // 输出: [trace_id=trace-12345] 用户登录成功

 // 链式调用添加字段
 log.WithField("trace_id", "trace-123").
  WithField("user_id", "user-456").
  Info("带字段的日志")

 // 🎨 Console 风格日志
 cg := log.NewConsoleGroup()

 // 📊 分组日志
 cg.Group("🚀 应用启动流程")
 cg.Info("开始初始化...")

 // 📋 表格展示
 config := map[string]interface{}{
  "环境":   "生产环境",
  "端口":   8080,
  "调试模式": false,
 }
 cg.Table(config)

 // ⏱️  计时器
 timer := cg.Time("数据库连接")
 // ... 执行数据库连接 ...
 timer.End() // 输出: ⏱️  数据库连接: 123.45ms

 cg.Info("✅ 启动完成")
 cg.GroupEnd()
}
```

### 链式配置

```go
// Builder 模式链式调用
log := logger.NewLogger().
 WithLevel(logger.INFO).
 WithShowCaller(true).
 WithColorful(true).
 WithPrefix("[MyApp]").
 WithFormat(logger.FormatJSON) // 支持 FormatText/FormatJSON/FormatXML/FormatCSV

log.Info("配置完成")
```

### 分布式追踪（OpenTelemetry 集成）

`trace_id` 的单一真相源是 OTel span。当 gRPC（otelgrpc StatsHandler）或 HTTP（Tracing 中间件）创建 span 后，日志会自动从 span 提取 `trace_id`，无需手动注入。

```go
// 1. 便捷方式：手动注入 trace_id（OTel 未启用时的 fallback）
ctx := logger.ContextWithTraceID(context.Background(), "trace-12345")
log.InfoContext(ctx, "带 trace 的日志")

// 2. 提取 trace_id
traceID := logger.ExtractTraceID(ctx)

// 3. gRPC 跨节点传播
// 客户端：注入到 outgoing metadata
outCtx := logger.InjectTraceToOutgoing(ctx, traceID)
// 服务端：从 incoming metadata 恢复到 ctx
restoredCtx := logger.RestoreTraceFromIncoming(incomingCtx)

// 4. 配置自定义上下文 key（除 trace_id 外额外提取）
log.WithContextKeys("request_id", "user_id")
// 此后 InfoContextKV 等方法会同时提取 trace_id（OTel）和这些 key
```

### 自定义上下文提取器

```go
// 自定义提取器，完全控制上下文信息输出
extractor := func(ctx context.Context) string {
 if ctx == nil {
  return ""
 }
 var parts []string
 if traceID, ok := ctx.Value("trace_id").(string); ok {
  parts = append(parts, "TraceID="+traceID)
 }
 if userID, ok := ctx.Value("user_id").(string); ok {
  parts = append(parts, "UserID="+userID)
 }
 if len(parts) > 0 {
  return "[" + strings.Join(parts, " ") + "] "
 }
 return ""
}

log.SetContextExtractor(extractor)
// 或 Builder 模式：log.WithContextExtractor(extractor)

ctx := context.WithValue(context.Background(), "trace_id", "trace-123")
ctx = context.WithValue(ctx, "user_id", "user-456")
log.InfoContext(ctx, "自定义上下文")
// 输出: [TraceID=trace-123 UserID=user-456] 自定义上下文
```

### 🎨 Console 风格日志功能

类似 JavaScript `console` 的日志分组、表格和计时器功能，让日志输出更加结构化和易读。

```go
log := logger.NewLogger()
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
// 输出美观的表格（自动处理中文字符宽度对齐）

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
```

> **说明**：包级便捷函数（如 `logger.Group()` / `logger.Info()` / `logger.Table()`）未提供。如需全局使用，请通过 `logger.GetGlobalLogger()` 获取默认 Logger，或用 `logger.New()`（`NewLogger` 的简写）创建实例。

**主要特性**：

- 🎯 **日志分组**:
  - `Group(label, ...args)` - 开始新分组
  - `GroupCollapsed(label, ...args)` - 开始折叠分组（仅显示 ERROR/FATAL）
  - `GroupEnd()` - 结束当前分组
  - 支持无限层级嵌套，自动缩进

- 📊 **表格渲染**:
  - `Table(data)` - 智能表格渲染
  - 支持格式: `[]map[string]interface{}`, `map[string]interface{}`, `[][]string`, `[]string`
  - 自动对齐、美化边框、智能列宽、东亚宽度对齐

- ⏱️  **计时器**:
  - `Time(label)` - 开始计时，返回 `*Timer`
  - `Timer.End()` - 结束计时并输出总耗时
  - `Timer.Log(message)` - 输出中间检查点
  - `Timer.Elapsed()` - 获取已过时间（不输出）
  - 智能时间格式化 (ns/μs/ms/s)
  - 自动清理过期计时器（`SetTimerMaxAge` / `SetTimerCleanupInterval` 可调）

- 🔄 **Context 集成**:
  - `InfoContext(ctx, ...)` / `DebugContext(ctx, ...)` / `WarnContext(ctx, ...)` / `ErrorContext(ctx, ...)`
  - 在分组内使用，自动继承缩进

**适用场景**：

- 🚀 应用启动流程展示
- 📊 批量数据处理进度
- 🔍 复杂业务流程追踪
- ⚡ 性能瓶颈分析
- 🐛 调试信息结构化输出

### 特殊日志类型

```go
log := logger.NewLogger()

// 成功日志
log.Success("用户注册成功")

// 加载日志
log.Loading("正在加载配置文件...")

// 配置日志
log.ConfigLog("数据库连接: %s", "localhost:3306")

// 启动/停止日志
log.Start("HTTP 服务器启动在端口 8080")
log.Stop("HTTP 服务器已停止")

// 数据库日志
log.Database("执行查询: SELECT * FROM users")

// 网络日志
log.Network("发送 HTTP 请求: GET /api/users")

// 安全日志（SECURITY 级别）
log.Security("检测到可疑登录尝试")

// 缓存日志
log.Cache("缓存命中: user:1001")

// 环境日志
log.Environment("当前环境: production")

// 性能日志（自动分级）
log.Performance("数据库查询", 50*time.Millisecond)
// 输出: ⚡ [PERF-EXCELLENT] 数据库查询 completed in 50ms

log.Performance("API 调用", 2*time.Second, map[string]any{
 "endpoint": "/api/users",
 "method":   "GET",
})
// 输出: 🐢 [PERF-SLOW] API 调用 completed in 2s | Details: map[endpoint:/api/users method:GET]

// 进度日志
log.Progress(50, 100, "数据导入")
// 输出: 🟡 [PROGRESS] 数据导入: 50/100 (50.0%)

// 里程碑日志
log.Milestone("系统初始化完成")
// 输出: 🎯 [MILESTONE] 系统初始化完成

// 健康检查日志
log.Health("数据库", true, "连接正常")
// 输出: ✅ [HEALTH] 数据库: HEALTHY | 连接正常

log.Health("Redis", false, "连接超时")
// 输出: ❌ [HEALTH] Redis: UNHEALTHY | 连接超时

// 审计日志（AUDIT 级别）
log.Audit("删除用户", "admin", "user:1001", "成功")
// 输出: 📋 [AUDIT] User: admin | Action: 删除用户 | Resource: user:1001 | Result: 成功

// 带上下文的特殊日志（自定义类型与级别）
log.LogSpecialContext(ctx, logger.SuccessType, logger.INFO, "操作成功")
log.LogSpecialContext(ctx, logger.DatabaseType, logger.INFO, "查询 %s 耗时 %v", table, duration)
log.LogSpecialContext(ctx, logger.SecurityType, logger.SECURITY, "检测到异常访问 %s", ip)
```

### 计时器辅助

```go
log := logger.NewLogger()

// 开始计时
timing := log.StartTiming("数据处理")

// 添加详细信息（链式）
timing.AddDetail("records", 1000).
 AddDetail("source", "database")

// 结束计时并自动记录性能日志
timing.End()
// 输出: ⚡ [PERF-FAST] 数据处理 completed in 80ms | Details: map[records:1000 source:database]
```

### 返回错误的日志方法

```go
func processUser(id int) error {
 user, err := findUser(id)
 if err != nil {
  // 记录错误日志并返回错误
  return log.ErrorReturn("查找用户失败: %v", err)
 }

 if user == nil {
  // 记录警告日志并返回错误
  return log.WarnReturn("用户不存在: %d", id)
 }

 return nil
}

// 带上下文的返回错误
func processOrder(ctx context.Context, orderID string) error {
 order, err := findOrder(orderID)
 if err != nil {
  return log.ErrorCtxReturn(ctx, "查找订单失败: %v", err)
 }

 return nil
}

// 键值对方式返回错误
func validateInput(data map[string]interface{}) error {
 if data["email"] == nil {
  return log.WarnKVReturn("缺少必填字段", "field", "email")
 }

 return nil
}
```

### 多输出写入器（Options 模式）

所有写入器均采用 **Options 函数式选项模式**构造：

```go
// 控制台写入器
consoleWriter := logger.NewConsoleWriter(
 logger.WithConsoleOutput(os.Stdout),
 logger.WithConsoleColor(true),
)

// 文件写入器
fileWriter := logger.NewFileWriter(
 logger.WithFileWriterPath("app.log"),
 logger.WithFilePermission(0644),
 logger.WithFileBufferSize(4096),
)

// 轮转文件写入器（按大小/数量/时间轮转，支持压缩）
rotateWriter := logger.NewRotateWriter(
 logger.WithFilePath("app.log"),
 logger.WithMaxSize(100*1024*1024), // 100MB
 logger.WithMaxFiles(10),          // 保留 10 个文件
 logger.WithMaxAge(30*24*time.Hour),
 logger.WithCompress(true),
)

// 缓冲写入器（在底层 writer 之上做缓冲）
bufferedWriter := logger.NewBufferedWriter(
 logger.WithBufferedUnderlying(fileWriter),
 logger.WithBufferSize(64*1024),
)

// 多输出写入器（同时分发到多个 writer）
multiWriter := logger.NewMultiWriter(
 logger.WithWriters(consoleWriter, fileWriter, rotateWriter),
)

// 异步批量写入器（后台 goroutine 批量 flush，高并发场景推荐）
asyncWriter := logger.NewAsyncBatchWriter(
 logger.WithAsyncUnderlying(fileWriter),
 logger.WithAsyncBatchSize(200),
 logger.WithAsyncFlushInterval(50*time.Millisecond),
 logger.WithAsyncChannelSize(4096),
)
defer asyncWriter.Close() // 应用退出时务必 Close

// 创建日志实例
log := logger.NewLogger().WithOutput(multiWriter)
log.Info("这条日志会同时输出到控制台和文件")
```

**通过配置统一创建 Writer**（适合配置驱动的场景）：

```go
config := &logger.WriterConfig{
 Type:     logger.OutputRotate,
 FilePath: "app.log",
 MaxSize:  100 * 1024 * 1024,
 MaxFiles: 10,
 Compress: true,
}
writer, err := logger.CreateWriter(config)
log := logger.NewLogger().WithOutput(writer)
```

### 空实现（禁用日志）

```go
// 创建空日志实例（不输出任何日志）
emptyLog := logger.NewEmptyLogger()

// 所有日志方法都不会产生任何输出
emptyLog.Info("这条日志不会输出")
emptyLog.Error("这条错误也不会输出")

// 带级别的空日志
emptyLogWithLevel := logger.NewEmptyLoggerWithLevel(logger.ERROR)

// 其他空实现
adapter := logger.NewEmptyAdapter("noop")
writer := logger.NewEmptyWriter()
hook := logger.NewEmptyHook([]logger.LogLevel{logger.ERROR})

// 判断是否为空日志 / 包装为空日志
logger.IsEmptyLogger(emptyLog)
logger.WrapWithEmpty(originalLogger)

// 适用场景：测试环境、性能敏感场景、占位实现
```

### 全局 Logger

```go
// 获取默认全局 Logger（包 init 时创建）
glog := logger.GetGlobalLogger()
glog.Info("全局日志")

// 全局设置
logger.SetGlobalLevel(logger.WARN)
logger.SetGlobalShowCaller(false)
```

## 📋 API 速览

### ILogger 接口方法（按类别）

| 类别 | 方法 |
| --- | --- |
| **基础（Printf）** | `Debug` `Info` `Warn` `Error` `Fatal` `Debugf` `Infof` `Warnf` `Errorf` `Fatalf` |
| **纯文本** | `DebugMsg` `InfoMsg` `WarnMsg` `ErrorMsg` `FatalMsg` |
| **结构化 KV** | `DebugKV` `InfoKV` `WarnKV` `ErrorKV` `FatalKV` |
| **字段映射** | `DebugWithFields` `InfoWithFields` `WarnWithFields` `ErrorWithFields` `FatalWithFields` |
| **多行** | `InfoLines` `ErrorLines` `WarnLines` `DebugLines` |
| **上下文** | `DebugContext` `InfoContext` `WarnContext` `ErrorContext` `FatalContext` |
| **上下文 KV** | `DebugContextKV` `InfoContextKV` `WarnContextKV` `ErrorContextKV` `FatalContextKV` |
| **上下文多行** | `DebugContextLines` `InfoContextLines` `WarnContextLines` `ErrorContextLines` |
| **返回错误** | `DebugReturn` `InfoReturn` `WarnReturn` `ErrorReturn` |
| **返回错误（ctx）** | `DebugCtxReturn` `InfoCtxReturn` `WarnCtxReturn` `ErrorCtxReturn` |
| **返回错误（KV）** | `DebugKVReturn` `InfoKVReturn` `WarnKVReturn` `ErrorKVReturn` |
| **原始条目** | `Log` `LogContext` `LogKV` `LogWithFields` `LogSpecialContext` |
| **链式构建** | `WithField` `WithFields` `WithError` `WithContext` |
| **标准库兼容** | `Print` `Printf` `Println` |
| **Console** | `NewConsoleGroup` `ConsoleGroup` `ConsoleGroupCollapsed` `ConsoleGroupEnd` `ConsoleTable` `ConsoleTime` |
| **配置** | `SetLevel` `GetLevel` `SetShowCaller` `IsShowCaller` `IsLevelEnabled` |

### Logger Builder 方法

`WithLevel` `WithShowCaller` `WithPrefix` `WithColorful` `WithOutput` `WithTimeFormat` `WithFormat` `WithCallerDepth` `WithShowStacktrace` `WithTimestampKey` `WithLevelKey` `WithMessageKey` `WithCallerKey` `WithStacktraceKey` `WithAsyncWrite` `WithBufferSize` `WithBatchSize` `WithBatchTimeout` `WithFormatter` `WithWriters` `WithHooks` `WithMiddleware` `WithContextExtractor` `WithContextKeys`

## 🤝 社区贡献

我们欢迎各种形式的贡献！请遵循以下指南：

### 提交代码

- **Fork 项目**

```bash
git clone https://github.com/kamalyes/go-logger.git
cd go-logger
```

- **创建特性分支**

```bash
git checkout -b feature/your-amazing-feature
```

- **编写代码和测试**

- 确保新功能有完整的测试套件
- 运行 `go test ./...` 确保所有测试通过
- 保持代码覆盖率 > 90%

- **提交更改**

```bash
git commit -m 'feat: add amazing new feature'
```

- **推送并创建 Pull Request**

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
- 包含性能基准测试（如适用，可参考 `benchmark.sh`）
- 验证并发安全性

## ⭐ Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=kamalyes/go-logger&type=Date)](https://star-history.com/#kamalyes/go-logger&Date)

## 许可证

该项目使用 MIT 许可证，详见 [LICENSE](LICENSE) 文件
