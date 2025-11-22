# Go Logger 迁移指南

本指南帮助您从其他日志框架迁移到 go-logger，或者了解如何使用 go-logger 的多框架兼容接口。

## 📋 目录

- [从 Logrus 迁移](#从-logrus-迁移)
- [从 Zap 迁移](#从-zap-迁移)
- [从 slog 迁移](#从-slog-迁移)
- [从 Zerolog 迁移](#从-zerolog-迁移)
- [从标准库 log 迁移](#从标准库-log-迁移)
- [混合使用策略](#混合使用策略)
- [性能对比](#性能对比)
- [最佳实践](#最佳实践)

## 从 Logrus 迁移

### 原始 Logrus 代码
```go
import "github.com/sirupsen/logrus"

func example() {
    log := logrus.New()
    log.SetLevel(logrus.DebugLevel)
    
    log.Info("服务启动")
    log.WithField("user_id", 12345).Info("用户登录")
    log.WithFields(logrus.Fields{
        "component": "auth",
        "action": "login",
        "ip": "192.168.1.1",
    }).Info("认证成功")
    
    err := errors.New("连接失败")
    log.WithError(err).Error("数据库连接错误")
}
```

### 迁移到 go-logger
```go
import "github.com/kamalyes/go-logger"

func example() {
    config := logger.DefaultConfig().
        WithLevel(logger.DEBUG).
        WithShowCaller(true).
        WithColorful(true)
    log := logger.NewLogger(config)
    
    log.Info("服务启动")
    log.WithField("user_id", 12345).Info("用户登录")
    log.WithFields(map[string]interface{}{
        "component": "auth",
        "action": "login",
        "ip": "192.168.1.1",
    }).Info("认证成功")
    
    err := errors.New("连接失败")
    log.WithError(err).Error("数据库连接错误")
}
```

### 迁移要点
- ✅ **API 完全兼容**：`WithField`、`WithFields`、`WithError` 等方法完全一样
- ✅ **级别映射简单**：`logrus.InfoLevel` → `logger.INFO`
- ✅ **零学习成本**：保持原有的编程习惯

## 从 Zap 迁移

### 原始 Zap 代码
```go
import "go.uber.org/zap"

func example() {
    logger, _ := zap.NewProduction()
    defer logger.Sync()
    
    logger.Info("服务启动",
        zap.String("version", "1.0.0"),
        zap.Int("port", 8080),
    )
    
    logger.Error("数据库错误",
        zap.String("database", "postgres"),
        zap.String("host", "localhost"),
        zap.Error(err),
        zap.Duration("timeout", 30*time.Second),
    )
}
```

### 迁移到 go-logger
```go
import "github.com/kamalyes/go-logger"

func example() {
    config := logger.DefaultConfig().
        WithLevel(logger.INFO).
        WithShowCaller(false)
    log := logger.NewLogger(config)
    
    // 使用键值对方式（推荐）
    log.InfoKV("服务启动",
        "version", "1.0.0",
        "port", 8080,
    )
    
    log.ErrorKV("数据库错误",
        "database", "postgres", 
        "host", "localhost",
        "error", err.Error(),
        "timeout", 30*time.Second,
    )
    
    // 或使用字段方式
    log.WithField("version", "1.0.0").
        WithField("port", 8080).
        Info("服务启动")
}
```

### 迁移要点
- ✅ **结构化日志**：使用 `InfoKV` 系列方法实现键值对日志
- ✅ **类型灵活**：支持任意类型的值，自动序列化
- ✅ **性能优化**：无需手动 Sync，自动管理资源

## 从 slog 迁移

### 原始 slog 代码
```go
import "log/slog"

func example(ctx context.Context) {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    
    logger.Info("处理请求", "method", "GET", "path", "/api/users")
    logger.InfoContext(ctx, "用户查询", "user_id", 12345)
    
    logger.Error("查询失败", "error", err, "query", "SELECT * FROM users")
}
```

### 迁移到 go-logger
```go
import "github.com/kamalyes/go-logger"

func example(ctx context.Context) {
    config := logger.DefaultConfig()
    log := logger.NewLogger(config)
    
    // 键值对方式
    log.InfoKV("处理请求", "method", "GET", "path", "/api/users")
    
    // 上下文感知日志（完全兼容）
    log.InfoContext(ctx, "用户查询，user_id: %d", 12345)
    
    // 或使用键值对 + 上下文
    log.LogKV(logger.ERROR, "查询失败", 
        "error", err.Error(), 
        "query", "SELECT * FROM users")
}
```

### 迁移要点
- ✅ **上下文兼容**：`InfoContext` 等方法完全兼容
- ✅ **结构化支持**：支持键值对和字段两种方式
- ✅ **格式灵活**：支持格式化字符串和纯键值对

## 从 Zerolog 迁移

### 原始 Zerolog 代码
```go
import "github.com/rs/zerolog/log"

func example() {
    log.Info().
        Str("service", "api").
        Int("port", 8080).
        Msg("服务启动")
    
    log.Error().
        Err(err).
        Str("component", "database").
        Msg("连接失败")
}
```

### 迁移到 go-logger
```go
import "github.com/kamalyes/go-logger"

func example() {
    log := logger.NewLogger(logger.DefaultConfig())
    
    // 使用链式调用方式
    log.WithField("service", "api").
        WithField("port", 8080).
        Info("服务启动")
    
    // 使用键值对方式
    log.InfoKV("服务启动",
        "service", "api",
        "port", 8080,
    )
    
    log.WithError(err).
        WithField("component", "database").
        Error("连接失败")
}
```

### 迁移要点
- ✅ **链式调用**：支持 `WithField` 链式调用
- ✅ **事件驱动**：支持基于事件的日志记录
- ✅ **零分配**：在可能的情况下优化内存分配

## 从标准库 log 迁移

### 原始标准库代码
```go
import "log"

func example() {
    log.SetPrefix("[APP] ")
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    log.Print("服务启动")
    log.Printf("监听端口: %d", 8080)
    log.Println("准备接受连接")
}
```

### 迁移到 go-logger
```go
import "github.com/kamalyes/go-logger"

func example() {
    config := logger.DefaultConfig().
        WithPrefix("[APP] ").
        WithShowCaller(true).
        WithTimeFormat("2006/01/02 15:04:05")
    log := logger.NewLogger(config)
    
    // 完全兼容的方法
    log.Print("服务启动")
    log.Printf("监听端口: %d", 8080)
    log.Println("准备接受连接")
    
    // 或使用增强方法
    log.Info("服务启动")
    log.InfoKV("服务配置", "port", 8080)
}
```

### 迁移要点
- ✅ **API 兼容**：完全支持 `Print`、`Printf`、`Println` 方法
- ✅ **配置映射**：前缀、时间格式等配置完全对应
- ✅ **逐步升级**：可以渐进式地使用新功能

## 混合使用策略

### 渐进式迁移
```go
// 步骤1：替换日志器创建
func step1() {
    // 原: log := logrus.New()
    log := logger.NewLogger(logger.DefaultConfig())
    
    // 保持原有调用不变
    log.WithField("user_id", 123).Info("用户操作")
}

// 步骤2：引入新功能
func step2() {
    log := logger.NewLogger(logger.DefaultConfig())
    
    // 混合使用
    log.WithField("component", "auth").Info("传统方式")
    log.InfoKV("新方式", "component", "auth", "action", "login")
}

// 步骤3：统一风格
func step3() {
    log := logger.NewLogger(logger.DefaultConfig())
    
    // 统一使用键值对方式（推荐）
    log.InfoKV("用户操作",
        "user_id", 123,
        "action", "login",
        "timestamp", time.Now(),
    )
}
```

### 团队协作策略
```go
// 定义团队标准的日志器创建函数
func NewAppLogger(component string) logger.ILogger {
    config := logger.DefaultConfig().
        WithLevel(logger.INFO).
        WithShowCaller(true).
        WithPrefix(fmt.Sprintf("[%s] ", component))
    
    return logger.NewLogger(config)
}

// 在各个模块中使用
func userService() {
    log := NewAppLogger("UserService")
    
    log.InfoKV("服务启动",
        "component", "user-service",
        "version", "v1.2.0",
    )
}

func authService() {
    log := NewAppLogger("AuthService")
    
    log.InfoKV("服务启动", 
        "component", "auth-service",
        "version", "v1.1.0",
    )
}
```

## 性能对比

### 基准测试结果
```
BenchmarkGoLogger-8                      8867611     130.1 ns/op   144 B/op    2 allocs/op
BenchmarkUltraFastLoggerNoTime-8        15794086      75.8 ns/op    24 B/op    1 allocs/op
BenchmarkSlog-8                          2085189     585.2 ns/op     0 B/op    0 allocs/op
BenchmarkStdLog-8                      305145283       3.9 ns/op     0 B/op    0 allocs/op
```

**性能分析**：
- ✅ **极致优化版本**: UltraFastLoggerNoTime **75.8ns/op, 24B/op, 1 alloc/op**
- ✅ **相比 slog**: go-logger 极致版快 **7.7倍** (75.8ns vs 585.2ns)
- ✅ **内存效率**: 相比标准版本减少 **83%** 内存使用 (144B → 24B)
- ✅ **功能完整**: 支持级别、字段、颜色、emoji 等丰富功能
- ⚠️ **vs 标准库**: std log 虽然极快，但功能有限（无级别、无字段支持）

**选择建议**：
- **极致性能场景**: 使用 `NewUltraFastLoggerNoTime()` - 最快速度
- **平衡性能与功能**: 使用 `NewUltraFastLogger()` - 完整功能 + 高性能  
- **标准使用**: 使用 `NewLogger()` - 完整功能 + 良好性能

### 性能优化建议
1. **关闭不必要功能**
   ```go
   config := logger.DefaultConfig().
       WithShowCaller(false).  // 生产环境关闭
       WithColorful(false)     // 非终端输出关闭
   ```

2. **选择合适的日志级别**
   ```go
   // 生产环境
   config.WithLevel(logger.INFO)
   
   // 开发环境
   config.WithLevel(logger.DEBUG)
   ```

3. **使用结构化日志**
   ```go
   // 高效：预分配字段
   log.InfoKV("操作完成",
       "duration", duration,
       "status", "success",
   )
   
   // 低效：字符串拼接
   log.Info("操作完成，耗时%v，状态%s", duration, "success")
   ```

## 最佳实践

### 1. 统一日志格式
```go
// 定义标准字段
type LogFields struct {
    Component   string `json:"component"`
    RequestID   string `json:"request_id"`
    UserID      int64  `json:"user_id"`
    Action      string `json:"action"`
    Duration    int64  `json:"duration_ms"`
    Error       string `json:"error,omitempty"`
}

func logUserAction(log logger.ILogger, fields LogFields) {
    log.InfoKV("用户操作",
        "component", fields.Component,
        "request_id", fields.RequestID,
        "user_id", fields.UserID,
        "action", fields.Action,
        "duration_ms", fields.Duration,
        "error", fields.Error,
    )
}
```

### 2. 上下文传递
```go
func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
    // 创建请求上下文
    ctx := context.WithValue(r.Context(), "request_id", generateID())
    log := logger.NewLogger(logger.DefaultConfig()).WithContext(ctx)
    
    // 在整个请求处理过程中使用
    processRequest(ctx, log, r)
}

func processRequest(ctx context.Context, log logger.ILogger, r *http.Request) {
    log.InfoContext(ctx, "开始处理请求: %s %s", r.Method, r.URL.Path)
    
    // 传递到下级函数
    result, err := businessLogic(ctx, log)
    if err != nil {
        log.ErrorContext(ctx, "处理失败: %v", err)
        return
    }
    
    log.InfoContext(ctx, "处理完成")
}
```

### 3. 错误处理标准化
```go
func standardErrorLogging(log logger.ILogger, err error, context string) {
    log.ErrorKV("操作失败",
        "context", context,
        "error", err.Error(),
        "error_type", fmt.Sprintf("%T", err),
        "timestamp", time.Now().Format(time.RFC3339),
    )
}

// 使用示例
func businessOperation(log logger.ILogger) error {
    err := someOperation()
    if err != nil {
        standardErrorLogging(log, err, "business_operation")
        return err
    }
    return nil
}
```

### 4. 配置管理
```go
type LoggerConfig struct {
    Level       string `yaml:"level"`
    Format      string `yaml:"format"`
    Output      string `yaml:"output"`
    ShowCaller  bool   `yaml:"show_caller"`
    Colorful    bool   `yaml:"colorful"`
}

func createLoggerFromConfig(cfg LoggerConfig) logger.ILogger {
    level, _ := logger.ParseLevel(cfg.Level)
    
    config := logger.DefaultConfig().
        WithLevel(level).
        WithShowCaller(cfg.ShowCaller).
        WithColorful(cfg.Colorful)
    
    return logger.NewLogger(config)
}
```

## 常见问题

### Q: 迁移过程中如何保证兼容性？
A: go-logger 提供了多种兼容接口，可以渐进式迁移，不需要一次性修改所有代码。

### Q: 性能会有影响吗？
A: 在正确配置的情况下，go-logger 的性能与主流日志库相当，某些场景下更优。

### Q: 如何处理现有的日志分析工具？
A: 保持日志格式不变，或者使用自定义格式化器来匹配现有工具的要求。

### Q: 是否需要修改现有的监控告警？
A: 如果保持相同的日志级别和关键字段，通常不需要修改监控配置。

## 总结

go-logger 提供了强大的多框架兼容性，让您可以：

- 🔄 **无痛迁移**：保持现有代码风格不变
- 🚀 **渐进升级**：逐步引入新功能
- 🎯 **团队协作**：支持多种编程习惯
- ⚡ **性能优化**：在兼容的基础上提供更好的性能
- 🛠️ **功能增强**：获得额外的功能支持

选择适合您团队和项目的迁移策略，享受现代化日志系统带来的便利！