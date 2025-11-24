# 自定义上下文提取器指南

## 概述

`UltraFastLogger` 支持自定义上下文提取器（Context Extractor），允许你从 `context.Context` 中提取任意信息并附加到日志中。这对于分布式追踪、请求链路跟踪、多租户系统等场景非常有用。

## 快速开始

### 使用默认提取器

默认情况下，`UltraFastLogger` 会自动提取 `TraceID` 和 `RequestID`：

```go
logger := logger.NewUltraFastLogger(logger.DefaultConfig())

ctx := context.Background()
ctx = context.WithValue(ctx, "trace_id", "trace-12345")
ctx = context.WithValue(ctx, "request_id", "req-67890")

logger.InfoContext(ctx, "用户登录成功")
// 输出: [TraceID=trace-12345 RequestID=req-67890] 用户登录成功
```

### 禁用上下文提取

```go
logger.SetContextExtractor(logger.NoOpContextExtractor)

logger.InfoContext(ctx, "这条日志不包含上下文信息")
// 输出: 这条日志不包含上下文信息
```

## 预定义提取器

### 1. NoOpContextExtractor

空操作提取器，不提取任何上下文信息。

```go
logger.SetContextExtractor(logger.NoOpContextExtractor)
```

### 2. SimpleTraceIDExtractor

只提取 `TraceID`，忽略其他字段。

```go
logger.SetContextExtractor(logger.SimpleTraceIDExtractor)

ctx = context.WithValue(ctx, "trace_id", "trace-12345")
logger.InfoContext(ctx, "消息")
// 输出: [TraceID=trace-12345] 消息
```

### 3. SimpleRequestIDExtractor

只提取 `RequestID`。

```go
logger.SetContextExtractor(logger.SimpleRequestIDExtractor)

ctx = context.WithValue(ctx, "request_id", "req-67890")
logger.InfoContext(ctx, "消息")
// 输出: [RequestID=req-67890] 消息
```

## 自定义字段提取

### CustomFieldExtractor

提取指定的自定义字段：

```go
extractor := logger.CustomFieldExtractor(
    []string{"user_id", "session_id"}, // 从 context.Value 提取
    []string{"x-tenant-id"},            // 从 gRPC metadata 提取
)
logger.SetContextExtractor(extractor)

ctx = context.WithValue(ctx, "user_id", "user-123")
ctx = context.WithValue(ctx, "session_id", "sess-456")
logger.InfoContext(ctx, "用户操作")
// 输出: [user_id=user-123 session_id=sess-456] 用户操作
```

## 组合提取器

### ChainContextExtractors

链接多个提取器，合并它们的输出：

```go
extractor := logger.ChainContextExtractors(
    logger.SimpleTraceIDExtractor,
    logger.ExtractFromContextValue("user_id", "User"),
    logger.ExtractFromContextValue("ip", "IP"),
)
logger.SetContextExtractor(extractor)

ctx = context.WithValue(ctx, "trace_id", "trace-12345")
ctx = context.WithValue(ctx, "user_id", "alice")
ctx = context.WithValue(ctx, "ip", "192.168.1.1")
logger.InfoContext(ctx, "API 请求")
// 输出: [TraceID=trace-12345] [User=alice] [IP=192.168.1.1] API 请求
```

## 使用构建器

`ContextExtractorBuilder` 提供了一种流式 API 来构建复杂的提取器：

```go
extractor := logger.NewContextExtractorBuilder().
    AddTraceID().
    AddRequestID().
    AddContextValue("tenant_id", "Tenant").
    AddContextValue("env", "Env").
    AddGRPCMetadata("x-api-key", "ApiKey").
    Build()

logger.SetContextExtractor(extractor)
```

### 构建器方法

- `AddTraceID()` - 添加 TraceID 提取器
- `AddRequestID()` - 添加 RequestID 提取器
- `AddContextValue(key, label)` - 从 context.Value 提取
- `AddGRPCMetadata(key, label)` - 从 gRPC metadata 提取
- `AddExtractor(extractor)` - 添加自定义提取器
- `Build()` - 构建最终提取器

## 高级用法

### 条件提取器

根据条件决定是否提取信息：

```go
extractor := logger.ConditionalContextExtractor(
    func(ctx context.Context) bool {
        env, ok := ctx.Value("env").(string)
        return ok && env == "production"
    },
    logger.ChainContextExtractors(
        logger.SimpleTraceIDExtractor,
        logger.SimpleRequestIDExtractor,
    ),
)
logger.SetContextExtractor(extractor)

// 只在生产环境提取详细信息
```

### 完全自定义提取器

实现 `ContextExtractor` 函数类型：

```go
customExtractor := func(ctx context.Context) string {
    if ctx == nil {
        return ""
    }
    
    userId, _ := ctx.Value("user_id").(string)
    userName, _ := ctx.Value("user_name").(string)
    
    if userId != "" || userName != "" {
        return fmt.Sprintf("[👤 %s (%s)] ", userId, userName)
    }
    
    return ""
}

logger.SetContextExtractor(customExtractor)

ctx = context.WithValue(ctx, "user_id", "12345")
ctx = context.WithValue(ctx, "user_name", "张三")
logger.InfoContext(ctx, "用户订单")
// 输出: [👤 12345 (张三)] 用户订单
```

### 带前缀的提取器

为提取的信息添加自定义前缀：

```go
extractor := logger.PrefixedContextExtractor(
    "🔍 ",
    logger.SimpleTraceIDExtractor,
)
logger.SetContextExtractor(extractor)
// 输出: 🔍 [TraceID=xxx] 消息
```

## 性能考虑

1. **NoOpContextExtractor** - 最快，完全跳过上下文提取
2. **SimpleTraceIDExtractor** - 很快，只提取一个字段
3. **默认提取器** - 平衡性能与功能
4. **ChainContextExtractors** - 性能随提取器数量增加而降低
5. **CustomFieldExtractor** - 性能取决于提取的字段数量

### 基准测试结果示例

```
BenchmarkNoOpContextExtractor        - 最快（零开销）
BenchmarkDefaultContextExtractor     - 稍慢（提取 2 个字段）
BenchmarkChainedContextExtractor     - 更慢（多个提取器）
```

## 最佳实践

1. **只提取需要的字段** - 避免提取过多信息影响性能
2. **使用构建器** - 代码更清晰易维护
3. **条件提取** - 在不同环境使用不同的提取策略
4. **缓存结果** - 对于昂贵的提取操作考虑缓存
5. **统一命名** - 团队内部统一 context key 的命名

## 示例场景

### 微服务追踪

```go
extractor := logger.NewContextExtractorBuilder().
    AddTraceID().
    AddContextValue("service_name", "Service").
    AddContextValue("span_id", "Span").
    Build()
```

### 多租户系统

```go
extractor := logger.NewContextExtractorBuilder().
    AddTraceID().
    AddRequestID().
    AddContextValue("tenant_id", "Tenant").
    AddContextValue("org_id", "Org").
    Build()
```

### API 网关

```go
extractor := logger.ChainContextExtractors(
    logger.SimpleTraceIDExtractor,
    logger.ExtractFromContextValue("api_key", "ApiKey"),
    logger.ExtractFromContextValue("client_ip", "IP"),
    logger.ExtractFromContextValue("user_agent", "UA"),
)
```

## 获取和恢复提取器

```go
// 保存当前提取器
originalExtractor := logger.GetContextExtractor()

// 临时更换
logger.SetContextExtractor(logger.NoOpContextExtractor)
// ... 执行一些操作 ...

// 恢复原始提取器
logger.SetContextExtractor(originalExtractor)
```

## 与其他功能配合

自定义上下文提取器可以与以下功能无缝配合：

- `WithField()` / `WithFields()` - 字段日志器
- `Clone()` - 克隆日志器时会复制提取器配置
- 所有日志级别的 `*Context()` 方法

```go
logger.WithField("component", "auth").
    InfoContext(ctx, "用户认证成功")
// 输出: [TraceID=xxx] 用户认证成功 {component: auth}
```

## 注意事项

1. 提取器函数应该是**无副作用**的
2. 提取器应该**快速返回**，避免阻塞
3. 返回的字符串会直接拼接到日志消息前，注意格式
4. 设置 `nil` 提取器会自动回退到默认提取器
5. 提取器在克隆日志器时会被复制

## 完整示例

查看 `examples/custom_context_extractor/main.go` 获取完整的使用示例。
