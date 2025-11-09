# Go Logger 示例

本目录包含了 go-logger 库的各种使用示例。

## 目录结构

```
examples/
├── basic/          # 基础使用示例
├── adapters/       # 适配器使用示例
├── factory/        # 工厂模式示例
├── benchmark/      # 性能基准测试示例
├── configuration/  # 配置使用示例
├── enhanced/       # 增强接口功能演示（新增）
├── context/        # 上下文感知日志演示（新增）
├── compatibility/  # 多框架兼容性演示（新增）
└── README.md       # 本文件
```

## 快速开始

### 1. 基础使用 (basic/)

演示最基本的日志库使用方法：

```bash
cd basic
go run main.go
```

**特性演示：**
- 基本日志级别使用
- 结构化日志
- 错误日志
- 链式调用
- 日志器克隆

### 2. 适配器使用 (adapters/)

演示如何使用日志适配器和管理器：

```bash
cd adapters
go run main.go
```

**特性演示：**
- 创建和管理适配器
- 适配器级别设置
- 管理器广播
- 健康检查
- 资源管理

### 3. 工厂模式 (factory/)

演示使用工厂模式创建日志器：

```bash
cd factory
go run main.go
```

**特性演示：**
- 日志器构建器
- 链式配置
- 钩子集成
- 中间件使用
- 复杂配置组合

### 4. 性能基准测试 (benchmark/)

演示性能测试和优化：

```bash
cd benchmark
go run main.go
```

**特性演示：**
- 基本性能测试
- 结构化日志性能
- 不同级别性能对比
- 并发日志性能

### 5. 配置使用 (configuration/)

演示各种配置选项：

```bash
cd configuration
go run main.go
```

**特性演示：**
- 默认配置
- 时间格式配置
- 输出目标配置
- 级别配置
- 调用者信息配置
- 彩色输出配置
- 配置验证和克隆

### 6. 增强接口功能 (enhanced/)

演示新增的增强接口功能：

```bash
cd enhanced
go run main.go
```

**特性演示：**
- 纯文本消息日志（DebugMsg, InfoMsg 等）
- 上下文感知日志（WithContext, DebugContext 等）
- 结构化键值对日志（InfoKV, ErrorKV 等）
- 原始日志条目方法（Log, LogKV, LogWithFields）
- 多框架兼容性（Zap, Logrus, slog 风格）
- 标准库兼容性（Print, Printf, Println）

### 7. 上下文感知日志 (context/)

演示微服务中的上下文日志追踪：

```bash
cd context
go run main.go
```

**特性演示：**
- 分布式追踪上下文传递
- 微服务间日志关联
- HTTP请求生命周期跟踪
- 错误传播和追踪
- 后台任务上下文管理

### 8. 多框架兼容性 (compatibility/)

演示与主流日志框架的兼容性：

```bash
cd compatibility
go run main.go
```

**特性演示：**
- Zap 风格的键值对日志
- Logrus 风格的字段日志
- slog 风格的上下文日志
- Zerolog 风格的事件日志
- 标准库 log 的兼容性
- 混合使用多种风格

## 核心概念

### 日志级别

```go
logger.DEBUG   // 调试级别
logger.INFO    // 信息级别  
logger.WARN    // 警告级别
logger.ERROR   // 错误级别
logger.FATAL   // 致命级别
```

### 基本使用

```go
// 创建日志器
log := logger.NewLogger(logger.DefaultConfig())

// 记录不同级别的日志
log.Debug("调试信息")
log.Info("普通信息") 
log.Warn("警告信息")
log.Error("错误信息")
```

### 增强功能（新）

#### 纯文本消息日志
```go
log.InfoMsg("纯文本信息")
log.ErrorMsg("纯文本错误")
```

#### 上下文感知日志
```go
ctx := context.Background()
log.InfoContext(ctx, "处理请求: %s", "login")
log.ErrorContext(ctx, "处理失败: %v", err)

// 创建带上下文的日志器
ctxLogger := log.WithContext(ctx)
ctxLogger.Info("携带上下文的日志")
```

#### 结构化键值对日志（类似Zap）
```go
log.InfoKV("用户登录",
    "user_id", 12345,
    "username", "john_doe",
    "ip_address", "192.168.1.100",
)
```

#### 原始日志条目方法
```go
log.Log(logger.INFO, "原始信息日志")
log.LogKV(logger.ERROR, "错误日志", "error_code", "E001")
log.LogWithFields(logger.DEBUG, "调试信息", map[string]interface{}{
    "component": "auth",
    "action": "validate",
})
```

#### 多框架兼容性
```go
// Zap 风格
log.InfoKV("消息", "key1", "value1", "key2", "value2")

// Logrus 风格  
log.WithField("component", "auth").Info("消息")
log.WithFields(map[string]interface{}{"k1": "v1"}).Info("消息")

// slog 风格
log.InfoContext(ctx, "处理请求: %s", "data")

// 标准库 log 风格
log.Printf("格式化消息: %s", "data")
log.Println("简单消息")
```

### 结构化日志

```go
// 单个字段
log.WithField("user_id", 12345).Info("用户操作")

// 多个字段
log.WithFields(map[string]interface{}{
    "user_id": 12345,
    "action": "login",
    "ip": "192.168.1.1",
}).Info("用户登录")
```

### 错误处理

```go
err := someFunction()
if err != nil {
    log.WithError(err).Error("操作失败")
}
```

### 配置选项

```go
config := logger.DefaultConfig().
    WithLevel(logger.DEBUG).
    WithShowCaller(true).
    WithColorful(true).
    WithPrefix("[MyApp] ").
    WithTimeFormat("2006-01-02 15:04:05")

log := logger.NewLogger(config)
```

## 最佳实践

### 1. 选择合适的日志级别和方法
   - 生产环境使用 INFO 或更高级别
   - 开发环境可以使用 DEBUG 级别
   - 根据场景选择最合适的日志方法：
     * 简单消息：使用 `InfoMsg()`, `ErrorMsg()` 等
     * 格式化消息：使用 `Info()`, `Error()` 等
     * 结构化数据：使用 `InfoKV()`, `ErrorKV()` 等
     * 上下文追踪：使用 `InfoContext()`, `ErrorContext()` 等

### 2. 上下文感知日志使用
   - 在HTTP服务中传递请求上下文
   - 使用 `WithContext()` 创建携带上下文的日志器
   - 在微服务间传递追踪信息
   ```go
   // 推荐模式
   func handleRequest(ctx context.Context, log logger.ILogger) {
       reqLogger := log.WithContext(ctx)
       reqLogger.Info("开始处理请求")
       // ... 处理逻辑
       reqLogger.Info("请求处理完成")
   }
   ```

### 3. 结构化日志使用
   - 优先使用键值对日志而非字符串拼接
   - 保持字段名称的一致性
   - 使用有意义的字段名
   ```go
   // 推荐
   log.InfoKV("用户操作",
       "user_id", 12345,
       "action", "login",
       "ip", "192.168.1.100",
   )
   
   // 不推荐
   log.Info("用户 %d 从 %s 登录", 12345, "192.168.1.100")
   ```

### 4. 错误处理最佳实践
   ```go
   // 使用 WithError 记录错误
   if err != nil {
       log.WithError(err).Error("操作失败")
       return err
   }
   
   // 或使用键值对格式
   if err != nil {
       log.ErrorKV("数据库操作失败",
           "operation", "INSERT",
           "table", "users",
           "error", err.Error(),
       )
       return err
   }
   ```

### 5. 多框架兼容性使用
   - 在迁移项目时可以保持原有的日志调用风格
   - 混合使用不同风格以适应不同场景
   - 团队内保持风格一致性

### 6. 合理使用调用者信息
   - 开发时开启，便于调试
   - 生产环境关闭，提高性能

### 7. 配置验证
   - 在应用启动时验证日志配置
   - 确保配置的正确性

### 8. 资源管理
   - 适当时机刷新和关闭日志器
   - 避免内存泄漏

## 性能提示

1. **关闭不必要的功能**
   ```go
   config := logger.DefaultConfig().
       WithShowCaller(false).  // 关闭调用者信息
       WithColorful(false)     // 关闭彩色输出
   ```

2. **批量日志处理**
   - 在高并发场景下考虑缓冲
   - 合理设置刷新间隔

3. **级别过滤**
   - 设置合适的日志级别
   - 避免记录不必要的日志

## 常见问题

### Q: 如何在生产环境中使用？
A: 建议使用 INFO 级别，关闭彩色输出和调用者信息，设置合适的时间格式。

### Q: 如何处理大量日志？
A: 考虑使用适配器将日志输出到文件或外部系统，设置合适的缓冲和刷新策略。

### Q: 如何集成到现有项目？
A: 可以逐步替换现有日志库，先在新功能中使用，然后逐步迁移旧代码。

### Q: 性能如何？
A: 在关闭不必要功能的情况下，性能表现良好。参考 benchmark 示例了解具体数据。

## 更多信息

- [项目主页](https://github.com/kamalyes/go-logger)
- [API 文档](https://pkg.go.dev/github.com/kamalyes/go-logger)
- [问题反馈](https://github.com/kamalyes/go-logger/issues)