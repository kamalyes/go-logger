# Go Logger 使用指南

## 目录

- [🏁 快速开始](#-快速开始)
- [🔧 基础配置](#-基础配置)
- [⚡ 性能层级](#-性能层级)
- [🎯 日志接口](#-日志接口)
- [📊 监控系统](#-监控系统)
- [🔍 分布式追踪](#-分布式追踪)
- [🏭 工厂模式](#-工厂模式)
- [🔌 适配器系统](#-适配器系统)
- [⚙️ 高级配置](#️-高级配置)
- [📈 性能调优](#-性能调优)

## 🏁 快速开始

### 环境要求

- Go 1.20+
- 建议使用 Go 1.21+ 获得最佳性能

### 安装

```bash
go get -u github.com/kamalyes/go-logger
```

### 简单使用

```go
package main

import (
    "github.com/kamalyes/go-logger"
)

func main() {
    // 创建默认日志器
    log := logger.New()
    
    // 基础日志记录
    log.Info("应用程序启动")
    log.Debug("调试信息: %s", "详细参数")
    log.Error("错误发生: %v", err)
    
    // 结构化日志
    log.InfoKV("用户登录", "user_id", 123, "username", "alice")
    
    // 链式调用
    log.WithField("component", "auth").
        WithField("action", "login").
        Info("用户认证成功")
}
```

## 🔧 基础配置

### 配置结构

```go
type Config struct {
    // 基础设置
    Level       LogLevel      `json:"level" yaml:"level"`
    Output      io.Writer     `json:"-" yaml:"-"`
    TimeFormat  TimeFormat    `json:"time_format" yaml:"time_format"`
    Colorful    bool          `json:"colorful" yaml:"colorful"`
    
    // 性能设置
    BufferSize     int  `json:"buffer_size" yaml:"buffer_size"`
    AsyncWrite     bool `json:"async_write" yaml:"async_write"`
    PoolSize       int  `json:"pool_size" yaml:"pool_size"`
    
    // 企业功能
    EnableMemoryStats  bool `json:"enable_memory_stats" yaml:"enable_memory_stats"`
    EnableDistributed  bool `json:"enable_distributed" yaml:"enable_distributed"`
    EnableMetrics      bool `json:"enable_metrics" yaml:"enable_metrics"`
    
    // 输出格式
    Format         FormatType `json:"format" yaml:"format"`
    TimestampKey   string     `json:"timestamp_key" yaml:"timestamp_key"`
    LevelKey       string     `json:"level_key" yaml:"level_key"`
    MessageKey     string     `json:"message_key" yaml:"message_key"`
    
    // 字段设置
    Fields         map[string]interface{} `json:"fields" yaml:"fields"`
    ContextFields  []string              `json:"context_fields" yaml:"context_fields"`
    
    // 钩子和中间件
    Hooks         []Hook        `json:"-" yaml:"-"`
    Middlewares   []Middleware  `json:"-" yaml:"-"`
    Writers       []Writer      `json:"-" yaml:"-"`
}
```

### 配置示例

```go
// 开发环境配置
devConfig := Config{
    Level:      DEBUG,
    Output:     os.Stdout,
    Colorful:   true,
    TimeFormat: TimeFormatStandard,
    Format:     FormatText,
    
    Fields: map[string]interface{}{
        "env":     "development",
        "service": "my-app",
        "version": "1.0.0",
    },
}

// 生产环境配置  
prodConfig := Config{
    Level:      INFO,
    Output:     os.Stdout,
    Colorful:   false,
    TimeFormat: TimeFormatRFC3339,
    Format:     FormatJSON,
    
    // 高性能设置
    AsyncWrite:     true,
    BufferSize:     8192,
    PoolSize:       10,
    
    // 企业功能
    EnableMemoryStats: true,
    EnableDistributed: true,
    EnableMetrics:     true,
    
    Fields: map[string]interface{}{
        "env":     "production",
        "service": "my-app",
        "version": "1.2.0",
    },
}

// 创建日志器
devLogger := logger.NewWithConfig(devConfig)
prodLogger := logger.NewWithConfig(prodConfig)
```

## ⚡ 性能层级

### 三层性能架构

go-logger 提供三种性能层级，满足不同场景需求：

#### 1. UltraFast Logger - 极致性能

适用场景：高并发、性能敏感、实时系统

```go
// 创建极速日志器（便利函数）
logger := logger.NewUltraFast()

// 或者使用配置
config := logger.DefaultConfig()
config.Level = logger.INFO
config.Colorful = false
config.ShowCaller = false
logger := logger.NewUltraFastLogger(config)

// 极致性能版本 - 无时间戳
logger := logger.NewUltraFastLoggerNoTime(os.Stdout, logger.INFO)

// 性能特点：
// - 7.56 ns/op 延迟
// - 0 分配
// - 零锁设计
// - 原子操作
```

#### 2. Optimized Logger - 平衡性能

适用场景：一般应用、开发调试、功能完整

```go
// 创建优化日志器（便利函数）
logger := logger.NewOptimized()

// 或者使用配置
config := logger.DefaultConfig()
config.Level = logger.INFO
config.ShowCaller = false
config.Colorful = true
logger := logger.NewLogger(config)

// 性能特点：
// - 22.85 ns/op 延迟  
// - 1 分配
// - 智能缓存
// - 对象池
```

#### 3. Full Logger - 企业级功能

适用场景：企业应用、监控需求、分布式系统

```go
// 创建全功能日志器（便利函数）
logger := logger.New()

// 或者使用配置
config := logger.DefaultConfig()
config.Level = logger.INFO
config.ShowCaller = true
config.Colorful = true
logger := logger.NewLogger(config)

// 功能特点：
// - 完整功能
// - 字段支持
// - 链式调用
// - 调用者信息
// - 彩色输出
```

### 性能对比

| 日志器类型 | 延迟 | 分配 | 功能完整度 | 适用场景 | 创建方式 |
|-----------|------|------|-----------|----------|----------|
| UltraFast | 7.56ns | 0 | ⭐⭐ | 高并发系统 | `logger.NewUltraFast()` |
| Optimized | 22.85ns | 1 | ⭐⭐⭐⭐ | 普通应用 | `logger.NewOptimized()` |
| Full | 130.1ns | 2 | ⭐⭐⭐⭐⭐ | 企业应用 | `logger.New()` |

## 🎯 日志接口

### 基础日志方法

```go
logger := logger.New()

// Printf 风格 
logger.Debug("调试信息: %s", variable)
logger.Info("信息: %d", count)
logger.Warn("警告: %v", warning)
logger.Error("错误: %v", err)
logger.Fatal("致命错误: %v", fatalErr)

// 纯文本方法
logger.DebugMsg("简单调试信息")
logger.InfoMsg("简单信息")
logger.WarnMsg("简单警告")
logger.ErrorMsg("简单错误")
logger.FatalMsg("简单致命错误")

// 结构化日志
logger.DebugKV("用户操作", "action", "login", "user_id", 123)
logger.InfoKV("请求处理", "method", "POST", "path", "/api/users", "status", 200)
logger.ErrorKV("数据库错误", "error", err, "table", "users", "operation", "insert")
```

### Context 感知日志

```go
import "context"

ctx := context.Background()
ctx = logger.WithTraceID(ctx, "trace-123")
ctx = logger.WithUserID(ctx, "user-456")

// Context 日志方法
logger.DebugContext(ctx, "Context调试: %s", info)
logger.InfoContext(ctx, "Context信息: %v", data)
logger.ErrorContext(ctx, "Context错误: %v", err)

// 或使用简化方法
logger.DebugWithContext(ctx, logger, "调试信息")
logger.InfoWithContext(ctx, logger, "信息内容")
logger.ErrorWithContext(ctx, logger, "错误信息")
```

### 链式日志

```go
// 添加字段
logger.WithField("component", "auth").
       WithField("user_id", 123).
       Info("用户登录成功")

// 添加多个字段
logger.WithFields(map[string]interface{}{
    "component": "database",
    "table": "users", 
    "operation": "select",
}).Debug("执行数据库查询")

// 错误链
logger.WithError(err).
       WithField("function", "processUser").
       Error("处理用户数据失败")
```

### 日志级别管理

```go
import "github.com/kamalyes/go-logger/level"

// 基础级别
logger.SetLevel(level.INFO)

// 24种详细级别支持
levels := []level.Level{
    level.TRACE,      // 最详细追踪
    level.DEBUG,      // 调试信息
    level.INFO,       // 一般信息  
    level.NOTICE,     // 重要信息
    level.WARN,       // 警告
    level.ERROR,      // 错误
    level.CRITICAL,   // 严重错误
    level.ALERT,      // 告警
    level.EMERGENCY,  // 紧急情况
    level.FATAL,      // 致命错误
    
    // 专用级别
    level.AUDIT,      // 审计日志
    level.SECURITY,   // 安全日志
    level.ACCESS,     // 访问日志
    level.PERFORMANCE,// 性能日志
    level.BUSINESS,   // 业务日志
    level.SYSTEM,     // 系统日志
    level.NETWORK,    // 网络日志
    level.DATABASE,   // 数据库日志
    level.CACHE,      // 缓存日志
    level.QUEUE,      // 队列日志
    level.SCHEDULE,   // 调度日志
    level.MONITOR,    // 监控日志
    level.METRIC,     // 指标日志
    level.PROFILING,  // 性能分析
}

// 级别管理器
manager := level.NewManager()
manager.SetLevel(level.INFO)
manager.SetPattern("auth.*", level.DEBUG)     // auth模块使用DEBUG
manager.SetPattern("db.*", level.WARN)       // 数据库模块使用WARN
manager.SetPattern("*.critical", level.ALERT) // 所有critical使用ALERT
```

## 📊 监控系统

### 内存监控

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

// 事件回调
monitor.OnMemoryThresholdExceeded(func(info *metrics.MemoryInfo) {
    fmt.Printf("⚠️ 内存使用率: %.2f%%\n", info.MemoryUsage)
    fmt.Printf("已用内存: %.2f MB\n", float64(info.UsedMemory)/1024/1024)
    fmt.Printf("堆内存: %.2f MB\n", float64(info.GoHeap)/1024/1024)
    
    // 可以触发告警或清理操作
    if info.MemoryUsage > 90.0 {
        runtime.GC() // 手动触发GC
    }
})

monitor.OnMemoryLeakDetected(func(report *metrics.LeakReport) {
    fmt.Printf("🚨 检测到内存泄漏: %s\n", report.GrowthTrend)
    fmt.Printf("增长率: %.2f bytes/s\n", report.MemoryGrowthRate)
    fmt.Printf("建议: %s\n", report.Recommendation)
})

// 启动监控
if err := monitor.Start(); err != nil {
    log.Fatal("启动内存监控失败:", err)
}
defer monitor.Stop()
```

### 实时监控数据

```go
// 获取实时内存信息
memInfo := monitor.GetMemoryInfo()
fmt.Printf("内存监控报告:\n")
fmt.Printf("  使用率: %.2f%%\n", memInfo.MemoryUsage)
fmt.Printf("  总内存: %.2f GB\n", float64(memInfo.TotalMemory)/1024/1024/1024)
fmt.Printf("  已用内存: %.2f MB\n", float64(memInfo.UsedMemory)/1024/1024)
fmt.Printf("  Go堆内存: %.2f MB\n", float64(memInfo.GoHeap)/1024/1024)
fmt.Printf("  Go栈内存: %.2f MB\n", float64(memInfo.GoStack)/1024/1024)

// 获取GC信息
gcInfo := monitor.GetGCInfo()
fmt.Printf("GC信息:\n")
fmt.Printf("  GC次数: %d\n", gcInfo.NumGC)
fmt.Printf("  总GC时间: %v\n", gcInfo.PauseTotalNs)
fmt.Printf("  平均GC时间: %.2f ms\n", float64(gcInfo.PauseTotalNs)/float64(gcInfo.NumGC)/1e6)

// 内存快照
snapshot, err := monitor.TakeHeapSnapshot()
if err == nil {
    fmt.Printf("内存快照:\n")
    fmt.Printf("  时间: %s\n", snapshot.Timestamp)
    fmt.Printf("  对象数量: %d\n", snapshot.ObjectCount)
    fmt.Printf("  内存大小: %.2f MB\n", float64(snapshot.MemorySize)/1024/1024)
}

// 内存泄漏分析
report := monitor.AnalyzeMemoryLeaks()
fmt.Printf("泄漏分析:\n")
fmt.Printf("  增长趋势: %s\n", report.GrowthTrend)
fmt.Printf("  增长率: %.2f bytes/s\n", report.MemoryGrowthRate)
fmt.Printf("  风险级别: %s\n", report.RiskLevel)
fmt.Printf("  建议: %s\n", report.Recommendation)
```

### 性能监控

```go
// 创建性能监控器
perfMonitor := metrics.NewDefaultPerformanceMonitor()

// 配置性能监控
perfMonitor.SetLatencyThreshold("api", time.Millisecond*100)    // API延迟阈值
perfMonitor.SetThroughputThreshold("requests", 1000.0)          // 请求吞吐量阈值

// 设置回调
perfMonitor.OnLatencyThresholdExceeded(func(operation string, latency time.Duration) {
    fmt.Printf("⚠️ %s 延迟超标: %v\n", operation, latency)
})

perfMonitor.OnThroughputThresholdExceeded(func(operation string, throughput float64) {
    fmt.Printf("📈 %s 吞吐量: %.2f ops/s\n", operation, throughput)
})

// 启动性能监控
perfMonitor.Start()
defer perfMonitor.Stop()

// 记录操作性能
start := time.Now()
// ... 执行业务操作 ...
duration := time.Since(start)
perfMonitor.RecordLatency("api_call", duration)
perfMonitor.RecordThroughput("requests", 1)

// 获取性能统计
stats := perfMonitor.GetPerformanceStats()
fmt.Printf("性能统计:\n")
fmt.Printf("  总操作数: %d\n", stats.TotalOperations)
fmt.Printf("  平均延迟: %v\n", stats.AvgLatency)
fmt.Printf("  吞吐量: %.2f ops/s\n", stats.Throughput)
fmt.Printf("  错误率: %.2f%%\n", stats.ErrorRate*100)
```

### 多级监控架构

```go
// 超轻量级监控 - 适用于高频操作
ultraMonitor := metrics.NewUltraLightMonitor()
ultraMonitor.Enable()

// 在关键路径中使用
func criticalPath() {
    done := ultraMonitor.Track()  // 开始追踪
    defer done(nil)               // 结束追踪，3.134ns开销
    
    // 业务逻辑...
}

// 优化监控 - 智能缓存
optimizedConfig := metrics.OptimizedConfig{
    CacheExpiry:     100 * time.Millisecond,
    EnableCaching:   true,
    LightweightMode: true,
}
optimizedMonitor := metrics.NewOptimizedMonitor(optimizedConfig)

// 使用优化监控
optimizedMonitor.Start()
heap, stack, used, numGC := optimizedMonitor.FastMemoryInfo()
fmt.Printf("快速内存信息: 堆=%d, 栈=%d, 已用=%d, GC=%d\n", heap, stack, used, numGC)

// 内存追踪器 - 阈值检测
tracker := metrics.NewMemoryTracker(512 * 1024 * 1024) // 512MB阈值
exceeded := tracker.Update(getCurrentMemory())
if exceeded {
    log.Warn("内存使用超过阈值")
}

// 智能健康检查
healthy, pressure := optimizedMonitor.QuickCheck()
fmt.Printf("系统健康: %v, 内存压力: %s\n", healthy, pressure)
```

## 🔍 分布式追踪

### Context ID 管理

```go
import "context"

ctx := context.Background()

// 设置各种ID
ctx = logger.WithTraceID(ctx, "trace-abc123")        // 分布式请求链路ID
ctx = logger.WithSpanID(ctx, "span-def456")          // 单个操作ID  
ctx = logger.WithRequestID(ctx, "req-ghi789")        // HTTP请求ID
ctx = logger.WithUserID(ctx, "user-12345")           // 用户ID
ctx = logger.WithSessionID(ctx, "session-67890")     // 会话ID
ctx = logger.WithCorrelationID(ctx, "corr-xyz999")   // 业务关联ID
ctx = logger.WithTenantID(ctx, "tenant-001")         // 租户ID

// 获取ID
traceID := logger.GetTraceID(ctx)           // "trace-abc123"
spanID := logger.GetSpanID(ctx)             // "span-def456"
userID := logger.GetUserID(ctx)             // "user-12345"

// 自动生成ID（如果不存在）
ctx, newTraceID := logger.GetOrGenerateTraceID(ctx)
newSpanID := logger.GenerateSpanID()
newRequestID := logger.GenerateRequestID()

// 批量提取所有字段
fields := logger.ExtractFields(ctx)
// fields = {
//   "trace_id": "trace-abc123",
//   "span_id": "span-def456",
//   "request_id": "req-ghi789",
//   "user_id": "user-12345",
//   "session_id": "session-67890",
//   "correlation_id": "corr-xyz999",
//   "tenant_id": "tenant-001"
// }
```

### Span 操作

```go
// 创建子Span（继承TraceID）
spanCtx := logger.CreateSpan(ctx, "database_query")
spanID := logger.GetSpanID(spanCtx) // 新生成的SpanID

// 在Span中记录日志
logger.DebugWithContext(spanCtx, log, "执行数据库查询", "table", "users")
logger.InfoWithContext(spanCtx, log, "查询完成", "rows", 42)

// 嵌套Span
subSpanCtx := logger.CreateSpan(spanCtx, "cache_lookup")  
logger.DebugWithContext(subSpanCtx, log, "检查缓存")

// 并行Span
go func() {
    parallelSpanCtx := logger.CreateSpan(ctx, "async_operation")
    logger.InfoWithContext(parallelSpanCtx, log, "异步操作开始")
    // 异步处理...
    logger.InfoWithContext(parallelSpanCtx, log, "异步操作完成")
}()
```

### 相关性链追踪

相关性链用于关联业务相关的多个操作，即使它们有不同的TraceID：

```go
// 创建相关性链
chain, chainCtx := logger.CreateCorrelationChain(ctx)
defer logger.EndCorrelationChain(chain) // 确保链结束

// 设置链属性
chain.SetTag("workflow", "user_registration")
chain.SetTag("business_type", "premium_user")
chain.SetMetric("expected_duration_ms", 5000)
chain.SetMetric("retry_count", 0)

// 步骤1：用户验证（独立TraceID）
validateCtx, _ := logger.GetOrGenerateTraceID(chainCtx)
if err := validateUser(validateCtx, userData); err != nil {
    chain.SetTag("failure_step", "validation")
    chain.SetMetric("retry_count", chain.GetMetric("retry_count").(int)+1)
    logger.ErrorWithContext(validateCtx, log, "用户验证失败", "error", err)
    return err
}

// 步骤2：创建账户（独立TraceID）
createCtx, _ := logger.GetOrGenerateTraceID(chainCtx)
account, err := createAccount(createCtx, userData)
if err != nil {
    chain.SetTag("failure_step", "account_creation")
    logger.ErrorWithContext(createCtx, log, "账户创建失败", "error", err)
    return err
}

// 步骤3：发送欢迎邮件（独立TraceID）
emailCtx, _ := logger.GetOrGenerateTraceID(chainCtx)
if err := sendWelcomeEmail(emailCtx, account.Email); err != nil {
    // 非关键操作，记录但不中断流程
    logger.WarnWithContext(emailCtx, log, "欢迎邮件发送失败", "error", err)
}

// 设置成功指标
chain.SetTag("status", "completed")
chain.SetMetric("account_id", account.ID)
chain.SetMetric("actual_duration_ms", chain.GetDuration().Milliseconds())

// 链自动结束时会记录完整的业务流程日志
```

### 操作日志记录器

操作日志记录器简化了复杂业务操作的日志记录：

```go
// 创建操作记录器
opLogger := logger.NewOperationLogger(ctx, log, "process_order")
defer func() {
    if r := recover(); r != nil {
        opLogger.EndWithError(fmt.Errorf("panic: %v", r))
        panic(r) // 重新抛出panic
    }
}()

// 设置操作属性
opLogger.SetTag("order_type", "premium")
opLogger.SetTag("customer_type", "enterprise")  
opLogger.SetTag("region", "us-west-2")

// 记录操作过程
opLogger.Info("订单处理开始")

// 步骤1：验证订单
opLogger.Debug("验证订单信息")
if err := validateOrder(opLogger.GetContext(), order); err != nil {
    opLogger.EndWithError(err, "step", "validation", "order_id", order.ID)
    return err
}

// 步骤2：库存检查
opLogger.Debug("检查库存")
available, err := checkInventory(opLogger.GetContext(), order.Items)
if err != nil {
    opLogger.EndWithError(err, "step", "inventory_check")
    return err
}
opLogger.SetMetric("items_checked", len(order.Items))
opLogger.SetMetric("items_available", available)

// 步骤3：处理支付
opLogger.Debug("处理支付")
payment, err := processPayment(opLogger.GetContext(), order.Payment)
if err != nil {
    opLogger.EndWithError(err, "step", "payment", "amount", order.Payment.Amount)
    return err
}
opLogger.SetMetric("payment_amount", payment.Amount)
opLogger.SetMetric("payment_method", payment.Method)

// 步骤4：创建订单
opLogger.Debug("创建订单记录")
createdOrder, err := createOrder(opLogger.GetContext(), order)
if err != nil {
    opLogger.EndWithError(err, "step", "order_creation")
    return err
}

// 成功完成
opLogger.SetMetric("order_id", createdOrder.ID)
opLogger.SetTag("final_status", "completed")
opLogger.End("total_amount", payment.Amount, "processing_time_ms", time.Since(startTime).Milliseconds())

// opLogger 自动记录：
// - 操作开始时间
// - 操作结束时间  
// - 操作总耗时
// - 所有设置的标签和指标
// - 成功/失败状态
```

### HTTP 服务集成

```go
// HTTP 中间件
func TracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 提取或生成TraceID
        traceID := r.Header.Get("X-Trace-ID")
        if traceID != "" {
            ctx = logger.WithTraceID(ctx, traceID)
        } else {
            ctx, traceID = logger.GetOrGenerateTraceID(ctx)
            w.Header().Set("X-Trace-ID", traceID)
        }
        
        // 生成RequestID  
        requestID := logger.GenerateRequestID()
        ctx = logger.WithRequestID(ctx, requestID)
        w.Header().Set("X-Request-ID", requestID)
        
        // 提取用户信息
        if userID := getUserIDFromAuth(r); userID != "" {
            ctx = logger.WithUserID(ctx, userID)
        }
        
        // 传递到下一个处理器
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// API 处理器
func CreateUserHandler(log logger.ILogger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 创建API操作记录器
        apiLogger := logger.NewOperationLogger(ctx, log, "create_user_api")
        defer apiLogger.End()
        
        // 设置API信息
        apiLogger.SetTag("method", r.Method)
        apiLogger.SetTag("path", r.URL.Path)
        apiLogger.SetTag("user_agent", r.Header.Get("User-Agent"))
        
        // 解析请求
        var req CreateUserRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            apiLogger.EndWithError(err, "step", "parse_request")
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        
        // 设置业务标签
        apiLogger.SetTag("username", req.Username)
        apiLogger.SetTag("user_type", req.Type)
        
        // 调用业务服务
        user, err := userService.CreateUser(apiLogger.GetContext(), &req)
        if err != nil {
            apiLogger.EndWithError(err, "step", "create_user", "username", req.Username)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // 记录成功信息
        apiLogger.SetMetric("user_id", user.ID)
        apiLogger.SetTag("status", "created")
        
        // 返回响应
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
        
        apiLogger.Info("用户创建成功")
    }
}
```

### 微服务间调用

```go
// 客户端传递追踪信息
func CallUserService(ctx context.Context, userID string) (*User, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", "/users/"+userID, nil)
    if err != nil {
        return nil, err
    }
    
    // 传递所有追踪信息
    if traceID := logger.GetTraceID(ctx); traceID != "" {
        req.Header.Set("X-Trace-ID", traceID)
    }
    if requestID := logger.GetRequestID(ctx); requestID != "" {
        req.Header.Set("X-Request-ID", requestID)
    }
    if userID := logger.GetUserID(ctx); userID != "" {
        req.Header.Set("X-User-ID", userID)
    }
    if correlationID := logger.GetCorrelationID(ctx); correlationID != "" {
        req.Header.Set("X-Correlation-ID", correlationID)
    }
    
    // 为此次调用创建新的SpanID
    spanCtx := logger.CreateSpan(ctx, "call_user_service")
    req.Header.Set("X-Span-ID", logger.GetSpanID(spanCtx))
    
    // 记录调用开始
    logger.InfoWithContext(spanCtx, log, "调用用户服务开始", "user_id", userID)
    
    // 执行请求
    start := time.Now()
    resp, err := httpClient.Do(req)
    duration := time.Since(start)
    
    if err != nil {
        logger.ErrorWithContext(spanCtx, log, "用户服务调用失败", 
            "error", err, "duration_ms", duration.Milliseconds())
        return nil, err
    }
    defer resp.Body.Close()
    
    // 记录调用结果
    logger.InfoWithContext(spanCtx, log, "用户服务调用完成",
        "status_code", resp.StatusCode, "duration_ms", duration.Milliseconds())
    
    if resp.StatusCode != http.StatusOK {
        err := fmt.Errorf("unexpected status: %d", resp.StatusCode)
        logger.ErrorWithContext(spanCtx, log, "用户服务返回错误", 
            "status_code", resp.StatusCode)
        return nil, err
    }
    
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        logger.ErrorWithContext(spanCtx, log, "解析用户响应失败", "error", err)
        return nil, err
    }
    
    return &user, nil
}

// 服务端提取追踪信息
func ExtractTracingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 提取上游传递的所有追踪信息
        headers := map[string]func(context.Context, string) context.Context{
            "X-Trace-ID":       logger.WithTraceID,
            "X-Span-ID":        logger.WithSpanID,
            "X-Request-ID":     logger.WithRequestID,
            "X-User-ID":        logger.WithUserID,
            "X-Session-ID":     logger.WithSessionID,
            "X-Correlation-ID": logger.WithCorrelationID,
            "X-Tenant-ID":      logger.WithTenantID,
        }
        
        for header, setterFunc := range headers {
            if value := r.Header.Get(header); value != "" {
                ctx = setterFunc(ctx, value)
            }
        }
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## 🏭 工厂模式

### 组件工厂

```go
import "github.com/kamalyes/go-logger/factory"

// 创建日志工厂
loggerFactory := logger.NewLoggerFactory()

// 注册自定义formatter
loggerFactory.RegisterFormatter("custom", func(config interface{}) (logger.IFormatter, error) {
    return &CustomFormatter{}, nil
})

// 注册自定义writer
loggerFactory.RegisterWriter("file", func(config interface{}) (logger.IWriter, error) {
    fileConfig := config.(*FileWriterConfig)
    return NewFileWriter(fileConfig.Path, fileConfig.MaxSize), nil
})

// 注册自定义hook
loggerFactory.RegisterHook("alert", func(config interface{}) (logger.IHook, error) {
    alertConfig := config.(*AlertConfig)
    return NewAlertHook(alertConfig.WebhookURL), nil
})

// 使用工厂创建Logger
config := &FactoryConfig{
    Level:     INFO,
    Formatter: "custom",
    Writers:   []string{"console", "file"},
    Hooks:     []string{"alert"},
}

createdLogger, err := loggerFactory.CreateLogger(config)
if err != nil {
    log.Fatal("创建logger失败:", err)
}
```

### 预设模板

```go
// 开发环境模板
devTemplate := logger.DevelopmentTemplate()
devLogger := loggerFactory.CreateFromTemplate(devTemplate)

// 生产环境模板  
prodTemplate := logger.ProductionTemplate()
prodLogger := loggerFactory.CreateFromTemplate(prodTemplate)

// 高性能模板
perfTemplate := logger.HighPerformanceTemplate()
perfLogger := loggerFactory.CreateFromTemplate(perfTemplate)

// 调试模板
debugTemplate := logger.DebugTemplate()
debugLogger := loggerFactory.CreateFromTemplate(debugTemplate)

// 自定义模板
customTemplate := &Template{
    Name: "api-server",
    Config: Config{
        Level:     INFO,
        Format:    FormatJSON,
        TimeFormat: TimeFormatRFC3339,
        
        Writers: []WriterConfig{
            {Type: "console", Config: ConsoleConfig{Colorful: false}},
            {Type: "file", Config: FileConfig{Path: "/var/log/api.log"}},
            {Type: "elasticsearch", Config: ESConfig{URL: "http://es:9200"}},
        },
        
        Hooks: []HookConfig{
            {Type: "metrics", Config: MetricsConfig{Endpoint: "/metrics"}},
            {Type: "alert", Config: AlertConfig{WebhookURL: "http://alert:8080"}},
        },
        
        Fields: map[string]interface{}{
            "service": "api-server",
            "version": "1.2.0",
            "env":     "production",
        },
    },
}

apiLogger := loggerFactory.CreateFromTemplate(customTemplate)
```

## 🔌 适配器系统

### 多适配器管理

```go
// 创建适配器管理器
manager := logger.NewAdapterManager()

// 添加控制台适配器
consoleConfig := &ConsoleAdapterConfig{
    Level:    DEBUG,
    Colorful: true,
    Format:   FormatText,
}
consoleAdapter, err := logger.CreateAdapter("console", consoleConfig)
if err != nil {
    log.Fatal(err)
}
manager.AddAdapter("console", consoleAdapter)

// 添加文件适配器
fileConfig := &FileAdapterConfig{
    Level:    INFO,
    Path:     "/var/log/app.log",
    MaxSize:  100 * 1024 * 1024, // 100MB
    MaxFiles: 10,
    Format:   FormatJSON,
}
fileAdapter, err := logger.CreateAdapter("file", fileConfig)
if err != nil {
    log.Fatal(err)
}
manager.AddAdapter("file", fileAdapter)

// 添加远程适配器
remoteConfig := &RemoteAdapterConfig{
    Level:    WARN,
    Endpoint: "http://log-server:8080/logs",
    Format:   FormatJSON,
    BufferSize: 1000,
    FlushInterval: time.Second * 30,
}
remoteAdapter, err := logger.CreateAdapter("remote", remoteConfig)
if err != nil {
    log.Fatal(err)
}
manager.AddAdapter("remote", remoteAdapter)

// 使用管理器记录日志
manager.Debug("调试信息")     // 只发送到console
manager.Info("信息内容")      // 发送到console和file  
manager.Error("错误信息")     // 发送到所有适配器

// 广播到所有适配器
manager.Broadcast(INFO, "重要信息")

// 获取适配器健康状态
health := manager.HealthCheck()
for name, healthy := range health {
    fmt.Printf("适配器 %s 健康状态: %v\n", name, healthy)
}

// 移除适配器
manager.RemoveAdapter("remote")

// 关闭所有适配器
manager.CloseAll()
```

### 适配器类型

#### 1. 控制台适配器

```go
config := &ConsoleAdapterConfig{
    Level:         DEBUG,
    Colorful:      true,
    Format:        FormatText,
    TimeFormat:    TimeFormatShort,
    ShowCaller:    true,
    CallerDepth:   4,
}
adapter := logger.CreateConsoleAdapter(config)
```

#### 2. 文件适配器

```go
config := &FileAdapterConfig{
    Level:           INFO,
    Path:            "/var/log/app.log",
    MaxSize:         100 * 1024 * 1024,  // 100MB
    MaxFiles:        10,
    MaxAge:          30 * 24 * time.Hour, // 30天
    Compress:        true,
    Format:          FormatJSON,
    AsyncWrite:      true,
    BufferSize:      4096,
    FlushInterval:   time.Second * 5,
}
adapter := logger.CreateFileAdapter(config)
```

#### 3. 网络适配器

```go
// TCP适配器
tcpConfig := &TCPAdapterConfig{
    Level:     WARN,
    Address:   "log-server:514",
    Network:   "tcp",
    Timeout:   time.Second * 10,
    Format:    FormatJSON,
}
tcpAdapter := logger.CreateTCPAdapter(tcpConfig)

// UDP适配器
udpConfig := &UDPAdapterConfig{
    Level:     INFO,
    Address:   "log-server:514", 
    MaxPacketSize: 1024,
    Format:    FormatJSON,
}
udpAdapter := logger.CreateUDPAdapter(udpConfig)

// HTTP适配器
httpConfig := &HTTPAdapterConfig{
    Level:         WARN,
    URL:           "http://log-server:8080/logs",
    Method:        "POST",
    Headers:       map[string]string{"Authorization": "Bearer token"},
    Timeout:       time.Second * 30,
    BufferSize:    1000,
    FlushInterval: time.Second * 60,
    Format:        FormatJSON,
}
httpAdapter := logger.CreateHTTPAdapter(httpConfig)
```

#### 4. 第三方集成适配器

```go
// Elasticsearch适配器
esConfig := &ElasticsearchAdapterConfig{
    Level:         INFO,
    URLs:          []string{"http://es1:9200", "http://es2:9200"},
    Index:         "logs-2024",
    Type:          "_doc",
    BufferSize:    1000,
    FlushInterval: time.Second * 30,
    Username:      "elastic",
    Password:      "password",
}
esAdapter := logger.CreateElasticsearchAdapter(esConfig)

// Redis适配器
redisConfig := &RedisAdapterConfig{
    Level:     DEBUG,
    Addr:      "redis:6379",
    Password:  "",
    DB:        0,
    Key:       "logs",
    MaxLength: 10000,
}
redisAdapter := logger.CreateRedisAdapter(redisConfig)

// Kafka适配器
kafkaConfig := &KafkaAdapterConfig{
    Level:   INFO,
    Brokers: []string{"kafka1:9092", "kafka2:9092"},
    Topic:   "logs",
    Partition: -1, // 自动分区
}
kafkaAdapter := logger.CreateKafkaAdapter(kafkaConfig)
```

### 自定义适配器

```go
// 实现IAdapter接口
type CustomAdapter struct {
    level  LogLevel
    config *CustomConfig
    client *CustomClient
}

func (a *CustomAdapter) Log(level LogLevel, message string, fields map[string]interface{}) error {
    if !a.IsLevelEnabled(level) {
        return nil
    }
    
    // 自定义日志处理逻辑
    logEntry := &CustomLogEntry{
        Timestamp: time.Now(),
        Level:     level.String(),
        Message:   message,
        Fields:    fields,
        Source:    a.config.Source,
    }
    
    return a.client.Send(logEntry)
}

func (a *CustomAdapter) IsLevelEnabled(level LogLevel) bool {
    return level >= a.level
}

func (a *CustomAdapter) SetLevel(level LogLevel) {
    a.level = level
}

func (a *CustomAdapter) GetLevel() LogLevel {
    return a.level
}

func (a *CustomAdapter) Close() error {
    return a.client.Close()
}

func (a *CustomAdapter) Flush() error {
    return a.client.Flush()
}

func (a *CustomAdapter) IsHealthy() bool {
    return a.client.IsConnected()
}

// 注册自定义适配器
logger.RegisterAdapter("custom", func(config interface{}) (logger.IAdapter, error) {
    customConfig := config.(*CustomConfig)
    client, err := NewCustomClient(customConfig)
    if err != nil {
        return nil, err
    }
    
    return &CustomAdapter{
        level:  customConfig.Level,
        config: customConfig,
        client: client,
    }, nil
})

// 使用自定义适配器
config := &CustomConfig{
    Level:    INFO,
    Endpoint: "http://custom-log-server:8080",
    Source:   "my-app",
}
adapter, err := logger.CreateAdapter("custom", config)
```

## ⚙️ 高级配置

### 配置文件管理

#### YAML 配置

```yaml
# config/logger.yaml
logger:
  # 基础设置
  level: info
  format: json
  time_format: rfc3339
  colorful: false
  
  # 性能设置
  async_write: true
  buffer_size: 8192
  pool_size: 10
  
  # 企业功能
  enable_memory_stats: true
  enable_distributed: true
  enable_metrics: true
  
  # 全局字段
  fields:
    service: "my-app"
    version: "1.2.0"
    environment: "production"
  
  # 上下文字段
  context_fields:
    - trace_id
    - user_id
    - session_id
    - tenant_id
  
  # 适配器配置
  adapters:
    - name: console
      type: console
      level: debug
      config:
        colorful: true
        format: text
        
    - name: file
      type: file
      level: info
      config:
        path: "/var/log/app.log"
        max_size: 100MB
        max_files: 10
        compress: true
        
    - name: elasticsearch
      type: elasticsearch
      level: warn
      config:
        urls: ["http://es:9200"]
        index: "logs-2024"
        buffer_size: 1000
        flush_interval: 30s
  
  # 钩子配置
  hooks:
    - name: metrics
      type: prometheus
      config:
        endpoint: "/metrics"
        
    - name: alert
      type: webhook
      config:
        url: "http://alert:8080/webhook"
        levels: [error, fatal]
  
  # 监控配置
  monitoring:
    memory:
      enabled: true
      threshold: 85.0
      sample_interval: 5s
      leak_detection: true
      
    performance:
      enabled: true
      latency_threshold: 100ms
      throughput_threshold: 1000.0
```

#### 加载配置

```go
// 从文件加载
config, err := logger.LoadConfigFromFile("config/logger.yaml")
if err != nil {
    log.Fatal("加载配置失败:", err)
}

// 从环境变量覆盖
config.OverrideFromEnv()

// 从命令行参数覆盖  
config.OverrideFromFlags()

// 创建logger
log, err := logger.NewWithConfig(config)
if err != nil {
    log.Fatal("创建logger失败:", err)
}
```

### 动态配置更新

```go
// 创建配置管理器
configManager := logger.NewConfigManager()

// 监听配置文件变化
configManager.WatchFile("config/logger.yaml", func(newConfig *Config) {
    log.Info("检测到配置变化，正在重新加载...")
    
    if err := log.UpdateConfig(newConfig); err != nil {
        log.Error("配置更新失败:", err)
    } else {
        log.Info("配置更新成功")
    }
})

// 通过API动态更新
http.HandleFunc("/admin/logger/config", func(w http.ResponseWriter, r *http.Request) {
    var newConfig Config
    if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    
    if err := log.UpdateConfig(&newConfig); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    w.WriteHeader(200)
    json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
})

// 运行时修改级别
log.SetLevel(DEBUG)

// 运行时添加适配器
newAdapter, err := logger.CreateAdapter("new-file", &FileAdapterConfig{
    Path: "/tmp/new.log",
    Level: INFO,
})
if err == nil {
    log.AddAdapter("new-file", newAdapter)
}

// 运行时移除适配器
log.RemoveAdapter("console")
```

### 环境特定配置

```go
// 根据环境加载不同配置
env := os.Getenv("GO_ENV")
if env == "" {
    env = "development"
}

var config *Config
switch env {
case "development":
    config = &Config{
        Level:      DEBUG,
        Format:     FormatText,
        Colorful:   true,
        TimeFormat: TimeFormatShort,
        
        Adapters: []AdapterConfig{
            {Type: "console", Level: DEBUG},
        },
    }
    
case "testing":
    config = &Config{
        Level:      INFO,
        Format:     FormatJSON,
        Colorful:   false,
        TimeFormat: TimeFormatRFC3339,
        
        Adapters: []AdapterConfig{
            {Type: "memory", Level: INFO}, // 内存适配器用于测试
        },
    }
    
case "production":
    config = &Config{
        Level:              INFO,
        Format:             FormatJSON,
        Colorful:           false,
        TimeFormat:         TimeFormatRFC3339,
        AsyncWrite:         true,
        BufferSize:         8192,
        EnableMemoryStats:  true,
        EnableDistributed:  true,
        
        Adapters: []AdapterConfig{
            {Type: "file", Level: INFO, Config: &FileConfig{Path: "/var/log/app.log"}},
            {Type: "elasticsearch", Level: WARN, Config: &ESConfig{URL: "http://es:9200"}},
        },
        
        Hooks: []HookConfig{
            {Type: "metrics", Config: &MetricsConfig{Endpoint: "/metrics"}},
            {Type: "alert", Config: &AlertConfig{WebhookURL: "http://alert:8080"}},
        },
    }
    
default:
    log.Fatal("未知环境:", env)
}

logger := logger.NewWithConfig(config)
```

## 📈 性能调优

### 性能分析工具

```go
import _ "net/http/pprof"

// 启用性能分析
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

// 性能分析示例
func analyzeLoggerPerformance() {
    // 创建不同配置的logger进行对比测试
    loggers := map[string]logger.ILogger{
        "ultra":     logger.NewUltraFast(),
        "optimized": logger.NewOptimized(), 
        "standard":  logger.New(),
    }
    
    for name, log := range loggers {
        // 预热
        for i := 0; i < 1000; i++ {
            log.Info("warm up")
        }
        
        // 性能测试
        start := time.Now()
        for i := 0; i < 100000; i++ {
            log.Info("performance test message")
        }
        duration := time.Since(start)
        
        fmt.Printf("%s logger: %v for 100k logs (%.2f ns/op)\n", 
            name, duration, float64(duration.Nanoseconds())/100000)
    }
}
```

### 内存优化

```go
// 内存池配置
config := &Config{
    // 对象池大小
    PoolSize: 50,
    
    // 缓冲区大小
    BufferSize: 8192,
    
    // 异步写入
    AsyncWrite: true,
    
    // 批量写入
    BatchSize: 100,
    BatchTimeout: time.Millisecond * 100,
}

// 内存监控与优化
monitor := metrics.NewMemoryOptimizer()
monitor.SetThreshold(80.0) // 80%内存使用率触发优化

monitor.OnOptimizationNeeded(func(usage float64) {
    fmt.Printf("内存使用率 %.2f%%，执行优化...\n", usage)
    
    // 强制GC
    runtime.GC()
    
    // 清理缓存
    logger.ClearBuffers()
    
    // 减少池大小
    logger.ShrinkPools()
})

// 启动优化监控
monitor.Start()
defer monitor.Stop()
```

### 并发优化

```go
// 并发安全配置
config := &Config{
    // 使用读写锁而非互斥锁
    UseMutex: false,
    
    // 分片锁减少竞争
    LockShards: 16,
    
    // 每个goroutine独立的缓冲区
    PerGoroutineBuffer: true,
    
    // 原子操作计数器
    AtomicCounters: true,
}

// 高并发场景优化
func optimizeForHighConcurrency() {
    // 使用本地缓冲区
    type localBuffer struct {
        buf    bytes.Buffer
        logger logger.ILogger
    }
    
    localPool := &sync.Pool{
        New: func() interface{} {
            return &localBuffer{
                logger: logger.NewUltraFast(),
            }
        },
    }
    
    // 工作goroutine
    worker := func(id int) {
        local := localPool.Get().(*localBuffer)
        defer localPool.Put(local)
        
        // 重置缓冲区
        local.buf.Reset()
        
        // 批量处理日志
        for i := 0; i < 1000; i++ {
            local.logger.Info("worker %d message %d", id, i)
        }
    }
    
    // 启动多个工作goroutine
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(id)
        }(i)
    }
    wg.Wait()
}
```

### 网络优化

```go
// 网络适配器优化配置
networkConfig := &NetworkAdapterConfig{
    // 连接池
    MaxConnections: 10,
    MaxIdleTime:    time.Minute * 5,
    
    // 批量发送
    BatchSize:      100,
    BatchTimeout:   time.Millisecond * 100,
    
    // 压缩
    Compression:    "gzip",
    CompressionLevel: 6,
    
    // 重试机制
    MaxRetries:     3,
    RetryDelay:     time.Second,
    BackoffFactor:  2.0,
    
    // 缓冲区
    SendBufferSize: 64 * 1024,
    RecvBufferSize: 64 * 1024,
    
    // 超时设置
    ConnectTimeout: time.Second * 10,
    WriteTimeout:   time.Second * 5,
    ReadTimeout:    time.Second * 5,
}

// 创建优化的网络适配器
adapter := logger.CreateOptimizedNetworkAdapter(networkConfig)
```

### 磁盘I/O优化

```go
// 文件适配器优化
fileConfig := &FileAdapterConfig{
    // 异步写入
    AsyncWrite: true,
    
    // 大缓冲区
    BufferSize: 256 * 1024, // 256KB
    
    // 批量刷新
    FlushInterval: time.Second * 5,
    FlushThreshold: 1000, // 1000条日志或5秒，先到者触发
    
    // 预分配文件
    PreallocSize: 100 * 1024 * 1024, // 100MB
    
    // 直接I/O（跳过系统缓存）
    DirectIO: true,
    
    // 文件同步策略
    SyncStrategy: "batch", // none, immediate, batch
    
    // 日志轮转
    RotateSize: 1024 * 1024 * 1024, // 1GB
    RotateTime: time.Hour * 24,      // 24小时
    MaxFiles:   30,                   // 保留30个文件
    
    // 压缩旧文件
    CompressRotated: true,
    CompressionType: "gzip",
}

// I/O监控
ioMonitor := metrics.NewIOMonitor()
ioMonitor.SetThresholds(
    80.0,  // 磁盘使用率阈值
    1000,  // IOPS阈值
    100,   // 延迟阈值(ms)
)

ioMonitor.OnThresholdExceeded(func(metric string, value float64) {
    switch metric {
    case "disk_usage":
        logger.Warn("磁盘使用率过高", "usage", value)
        // 清理旧日志文件
        logger.CleanupOldLogs()
        
    case "iops":
        logger.Warn("磁盘IOPS过高", "iops", value)
        // 增加批量大小，减少写入频率
        logger.AdjustBatchSize(2.0)
        
    case "latency":
        logger.Warn("磁盘延迟过高", "latency_ms", value)
        // 启用压缩，减少I/O量
        logger.EnableCompression()
    }
})
```

---

更多详细信息和高级用法，请参考：

- [📊 性能详解](PERFORMANCE.md) - 深入性能分析和优化技术
- [🔄 迁移指南](MIGRATION.md) - 从其他日志库迁移
- [🎯 Context使用指南](CONTEXT_USAGE.md) - 分布式追踪完整指南
- [📝 更新日志](CHANGELOG.md) - 版本更新记录