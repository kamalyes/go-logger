# Go Logger 示例集合

本目录包含了 go-logger 库的完整使用示例，涵盖从基础使用到企业级功能的各种场景。

## 📚 示例导航

### 🚀 快速开始
- **[基础使用](basic/)** - 最简单的日志使用方法，包括便利函数验证
- **[便利函数](convenience/)** - NewUltraFast()、NewOptimized()、New() 使用示例
- **[性能测试](benchmark/)** - 性能基准测试和优化

### 🔧 核心功能  
- **[配置系统](configuration/)** - 完整的配置选项演示
- **[适配器系统](adapters/)** - 多种适配器使用和管理
- **[格式化器](formatters/)** - 日志格式化器的使用
- **[监控系统](monitoring/)** - 内存和性能监控功能

### 🎯 高级功能
- **[Context追踪](context/)** - 分布式系统上下文管理
- **[工厂模式](factory/)** - 高级日志器构建
- **[增强功能](enhanced/)** - 企业级增强功能
- **[兼容性](compatibility/)** - 多框架兼容性演示

## 🏃‍♂️ 快速运行

### 运行所有示例
```bash
make run-all
```

### 运行特定示例
```bash
cd <example-directory>
go run main.go
```

## 📖 示例说明

### 1. 基础使用 (basic/)
演示最基本的日志功能：
- 基本日志级别
- 结构化日志
- 错误处理
- 链式调用

```bash
cd basic && go run main.go
```

### 2. 配置系统 (configuration/)
演示完整的配置功能：
- 环境特定配置
- 动态配置更新
- 配置文件支持
- 最佳实践

```bash
cd configuration && go run main.go
```

### 3. 适配器系统 (adapters/)
演示各种适配器的使用：
- Console、File、Network 适配器
- 企业级适配器 (Elasticsearch, Redis, Kafka)
- 自定义适配器开发
- 适配器管理和路由

```bash
cd adapters && go run main.go
```

### 4. 格式化器 (formatters/)
演示日志格式化功能：
- 内置格式化器 (JSON, Text, CSV, XML)
- 自定义格式化器
- 模板引擎
- 条件格式化

```bash
cd formatters && go run main.go
```

### 5. 监控系统 (monitoring/)
演示监控和性能分析：
- 内存监控
- 性能监控
- I/O 监控
- 告警系统

```bash
cd monitoring && go run main.go
```

### 6. Context追踪 (context/)
演示分布式系统上下文管理：
- TraceID 和 SpanID
- 微服务间追踪
- HTTP 请求追踪
- 错误传播

```bash
cd context && go run main.go
```

### 7. 性能测试 (benchmark/)
演示性能测试和优化：
- 基准测试
- 并发性能测试
- 内存使用分析
- 性能对比

```bash
cd benchmark && go run main.go
```

### 8. 工厂模式 (factory/)
演示高级日志器构建：
- 构建器模式
- 复杂配置组合
- 中间件集成
- 插件系统

```bash
cd factory && go run main.go
```

### 9. 增强功能 (enhanced/)
演示企业级增强功能：
- 纯文本消息日志
- 上下文感知日志
- 结构化键值对日志
- 多框架兼容性

```bash
cd enhanced && go run main.go
```

### 10. 兼容性 (compatibility/)
演示多框架兼容性：
- Zap 风格日志
- Logrus 风格日志
- slog 风格日志
- 标准库兼容

```bash
cd compatibility && go run main.go
```

## 🎯 最佳实践指南

### 选择合适的日志方法

```go
// 1. 简单消息 - 使用 *Msg 方法
logger.InfoMsg("操作完成")
logger.ErrorMsg("操作失败")

// 2. 格式化消息 - 使用标准方法
logger.Info("处理用户 %s 的请求", username)
logger.Error("连接数据库失败: %v", err)

// 3. 结构化日志 - 使用 *KV 方法
logger.InfoKV("用户登录",
    "user_id", 12345,
    "username", "john",
    "ip", "192.168.1.100",
)

// 4. 上下文追踪 - 使用 *Context 方法
logger.InfoContext(ctx, "处理请求")
logger.ErrorContext(ctx, "请求失败: %v", err)
```

### 性能优化建议

```go
// 生产环境配置
config := logger.Config{
    Level:      logger.INFO,          // 合适的日志级别
    ShowCaller: false,                // 关闭调用者信息
    Colorful:   false,                // 关闭颜色输出
    TimeFormat: "2006-01-02T15:04:05Z", // 标准时间格式
}
```

### 错误处理模式

```go
// 方式1: 使用 WithError
if err != nil {
    logger.WithError(err).Error("操作失败")
    return err
}

// 方式2: 使用结构化日志
if err != nil {
    logger.ErrorKV("数据库操作失败",
        "operation", "INSERT",
        "table", "users",
        "error", err.Error(),
    )
    return err
}
```

### 上下文使用模式

```go
// HTTP 处理器模式
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    reqLogger := logger.WithContext(ctx)
    
    reqLogger.Info("开始处理请求")
    defer reqLogger.Info("请求处理完成")
    
    // 使用带上下文的日志器
    processRequest(reqLogger)
}

// 服务层模式
func processRequest(log logger.ILogger) {
    log.Info("执行业务逻辑")
    // ...
}
```

## 📊 性能参考

基于最新的基准测试结果：

```
BenchmarkUltraFastLogger-8       157894737     7.56 ns/op     0 B/op     0 allocs/op
BenchmarkStandardLogger-8         52631578    22.85 ns/op     8 B/op     1 allocs/op
BenchmarkStructuredLogging-8      15789473    75.8 ns/op    24 B/op     1 allocs/op
```

详细性能分析请查看 [性能文档](../docs/PERFORMANCE.md)。

## 🔗 相关文档

- [📊 性能详解](../docs/PERFORMANCE.md)
- [🔧 配置指南](../docs/CONFIGURATION.md)  
- [🧩 适配器系统](../docs/ADAPTERS.md)
- [📊 监控系统](../docs/MONITORING.md)
- [🎨 格式化器](../docs/FORMATTERS.md)
- [🎯 Context使用指南](../docs/CONTEXT_USAGE.md)

## 🤝 贡献指南

欢迎提交新的示例！请确保：

1. 添加详细的注释说明
2. 包含错误处理
3. 遵循最佳实践
4. 更新相应的文档

## 📞 获取帮助

- [GitHub Issues](https://github.com/kamalyes/go-logger/issues)
- [API 文档](https://pkg.go.dev/github.com/kamalyes/go-logger)
- [项目主页](https://github.com/kamalyes/go-logger)