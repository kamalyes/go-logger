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

1. **选择合适的日志级别**
   - 生产环境使用 INFO 或更高级别
   - 开发环境可以使用 DEBUG 级别

2. **使用结构化日志**
   - 便于日志分析和搜索
   - 提供更多上下文信息

3. **合理使用调用者信息**
   - 开发时开启，便于调试
   - 生产环境关闭，提高性能

4. **配置验证**
   - 在应用启动时验证日志配置
   - 确保配置的正确性

5. **资源管理**
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