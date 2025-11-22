# go-logger Context 包使用文档

## 简介

`github.com/kamalyes/go-logger` 的 context 包提供了分布式系统中的上下文管理和链路追踪功能，帮助您轻松实现日志的关联和追踪。

## 核心概念

### 追踪 ID 类型

| ID 类型 | 用途 | 示例值 | 生命周期 |
|---------|------|--------|----------|
| **TraceID** | 分布式请求链路标识 | `trace-abc123` | 整个请求调用链 |
| **SpanID** | 单个操作标识 | `span-def456` | 单个操作内 |
| **RequestID** | HTTP/RPC 请求标识 | `req-ghi789` | 单次网络请求 |
| **UserID** | 用户标识 | `user-12345` | 用户会话期间 |
| **SessionID** | 用户会话标识 | `session-67890` | 登录到登出 |
| **CorrelationID** | 业务关联标识 | `corr-xyz999` | 相关业务流程 |
| **TenantID** | 租户标识 | `tenant-001` | 多租户隔离 |

## 基础用法

### 设置和获取 ID

```go
package main

import (
    "context"
    "fmt"
    "github.com/kamalyes/go-logger"
)

func main() {
    ctx := context.Background()
    
    // 设置各种 ID
    ctx = logger.WithTraceID(ctx, "trace-123")
    ctx = logger.WithUserID(ctx, "user-456")
    ctx = logger.WithTenantID(ctx, "tenant-001")
    
    // 获取 ID
    traceID := logger.GetTraceID(ctx)
    userID := logger.GetUserID(ctx)
    
    fmt.Printf("TraceID: %s, UserID: %s\n", traceID, userID)
}
```

### 自动生成 ID

```go
// 获取或自动生成 TraceID
ctx, traceID := logger.GetOrGenerateTraceID(ctx)

// 直接生成新 ID
newTraceID := logger.GenerateTraceID()
newSpanID := logger.GenerateSpanID()
newRequestID := logger.GenerateRequestID()
newCorrelationID := logger.GenerateCorrelationID()
```

### 日志记录（自动包含上下文字段）

```go
// 方式1：直接使用带上下文的日志方法
logger.InfoWithContext(ctx, myLogger, "用户登录成功", "username", "john")

// 方式2：创建包含上下文字段的 logger
contextLogger := logger.WithLogger(ctx, myLogger)
contextLogger.InfoMsg("用户登录成功")
contextLogger.ErrorKV("登录失败", "reason", "密码错误")
```

## 典型使用场景

### 1. HTTP API 服务

#### 设置追踪中间件

```go
// middleware.go
package main

import (
    "net/http"
    "github.com/kamalyes/go-logger"
)

func TraceMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 从请求头提取 TraceID，或者生成新的
        traceID := r.Header.Get("X-Trace-ID")
        if traceID != "" {
            ctx = logger.WithTraceID(ctx, traceID)
        } else {
            ctx, traceID = logger.GetOrGenerateTraceID(ctx)
            w.Header().Set("X-Trace-ID", traceID)
        }
        
        // 生成请求ID
        requestID := logger.GenerateRequestID()
        ctx = logger.WithRequestID(ctx, requestID)
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 提取用户信息中间件
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 从 JWT 或其他方式提取用户信息
        userID := extractUserID(r) // 您的实现
        if userID != "" {
            ctx = logger.WithUserID(ctx, userID)
        }
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

#### 业务处理器

```go
// handler.go
func CreateUserHandler(log logger.ILogger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 创建操作日志记录器
        opLogger := logger.NewOperationLogger(ctx, log, "create_user")
        defer opLogger.End() // 确保记录操作完成时间
        
        // 解析请求
        var req CreateUserRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            opLogger.EndWithError(err, "step", "parse_request")
            http.Error(w, "Invalid request", http.StatusBadRequest)
            return
        }
        
        // 设置业务标签
        opLogger.SetTag("username", req.Username)
        opLogger.SetTag("user_type", req.Type)
        
        // 调用业务服务
        user, err := userService.CreateUser(opLogger.GetContext(), &req)
        if err != nil {
            opLogger.EndWithError(err, "step", "create_user")
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        // 记录成功指标
        opLogger.SetMetric("user_id", user.ID)
        opLogger.Info("用户创建成功")
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(user)
    }
}
```

### 2. 微服务间调用

#### 客户端传递上下文

```go
// client.go
package client

import (
    "context"
    "net/http"
    "github.com/kamalyes/go-logger"
)

type UserClient struct {
    baseURL string
    client  *http.Client
    logger  logger.ILogger
}

func (c *UserClient) GetUser(ctx context.Context, userID string) (*User, error) {
    // 创建请求
    req, _ := http.NewRequestWithContext(ctx, "GET", 
        c.baseURL+"/users/"+userID, nil)
    
    // 传递追踪信息到下游服务
    req.Header.Set("X-Trace-ID", logger.GetTraceID(ctx))
    req.Header.Set("X-User-ID", logger.GetUserID(ctx))
    req.Header.Set("X-Correlation-ID", logger.GetCorrelationID(ctx))
    
    // 为这次调用创建新的 Span
    spanCtx := logger.CreateSpan(ctx, "call_user_service")
    req.Header.Set("X-Span-ID", logger.GetSpanID(spanCtx))
    
    // 记录调用日志
    logger.InfoWithContext(spanCtx, c.logger, "调用用户服务", 
        "user_id", userID, "service", "user-service")
    
    // 执行请求
    resp, err := c.client.Do(req)
    if err != nil {
        logger.ErrorWithContext(spanCtx, c.logger, "用户服务调用失败", "error", err)
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        err := fmt.Errorf("用户服务返回错误状态: %d", resp.StatusCode)
        logger.ErrorWithContext(spanCtx, c.logger, "用户服务响应异常", "status_code", resp.StatusCode)
        return nil, err
    }
    
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        logger.ErrorWithContext(spanCtx, c.logger, "解析响应失败", "error", err)
        return nil, err
    }
    
    logger.InfoWithContext(spanCtx, c.logger, "用户服务调用成功")
    return &user, nil
}
```

#### 服务端提取上下文

```go
// server.go
func ExtractTraceMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 提取上游传递的追踪信息
        if traceID := r.Header.Get("X-Trace-ID"); traceID != "" {
            ctx = logger.WithTraceID(ctx, traceID)
        }
        if spanID := r.Header.Get("X-Span-ID"); spanID != "" {
            ctx = logger.WithSpanID(ctx, spanID)
        }
        if userID := r.Header.Get("X-User-ID"); userID != "" {
            ctx = logger.WithUserID(ctx, userID)
        }
        if correlationID := r.Header.Get("X-Correlation-ID"); correlationID != "" {
            ctx = logger.WithCorrelationID(ctx, correlationID)
        }
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 3. 数据库操作

```go
// repository.go
package repository

import (
    "context"
    "database/sql"
    "github.com/kamalyes/go-logger"
)

type UserRepository struct {
    db     *sql.DB
    logger logger.ILogger
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
    // 创建数据库查询的 Span
    spanCtx := logger.CreateSpan(ctx, "db_query_user")
    
    // 记录查询日志
    logger.DebugWithContext(spanCtx, r.logger, "执行用户查询", 
        "user_id", userID, "table", "users")
    
    query := "SELECT id, username, email, created_at FROM users WHERE id = ?"
    row := r.db.QueryRowContext(spanCtx, query, userID)
    
    var user User
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            logger.WarnWithContext(spanCtx, r.logger, "用户不存在", "user_id", userID)
            return nil, ErrUserNotFound
        }
        logger.ErrorWithContext(spanCtx, r.logger, "数据库查询失败", 
            "error", err, "user_id", userID)
        return nil, err
    }
    
    logger.InfoWithContext(spanCtx, r.logger, "用户查询成功", "user_id", userID)
    return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
    spanCtx := logger.CreateSpan(ctx, "db_insert_user")
    
    logger.DebugWithContext(spanCtx, r.logger, "插入用户记录", 
        "username", user.Username)
    
    query := `INSERT INTO users (id, username, email, password_hash, created_at) 
              VALUES (?, ?, ?, ?, ?)`
    
    _, err := r.db.ExecContext(spanCtx, query, 
        user.ID, user.Username, user.Email, user.PasswordHash, time.Now())
    
    if err != nil {
        logger.ErrorWithContext(spanCtx, r.logger, "用户创建失败", 
            "error", err, "username", user.Username)
        return err
    }
    
    logger.InfoWithContext(spanCtx, r.logger, "用户创建成功", 
        "user_id", user.ID, "username", user.Username)
    return nil
}
```

### 4. 异步任务处理

#### 任务发布

```go
// publisher.go
package task

import (
    "context"
    "encoding/json"
    "github.com/kamalyes/go-logger"
)

type TaskMessage struct {
    TaskID   string      `json:"task_id"`
    TaskType string      `json:"task_type"`
    Data     interface{} `json:"data"`
    
    // 追踪上下文
    TraceID   string `json:"trace_id"`
    UserID    string `json:"user_id"`
    TenantID  string `json:"tenant_id"`
}

type TaskPublisher struct {
    queue  MessageQueue // 您的消息队列实现
    logger logger.ILogger
}

func (p *TaskPublisher) PublishEmailTask(ctx context.Context, email EmailData) error {
    taskID := generateTaskID()
    
    msg := &TaskMessage{
        TaskID:   taskID,
        TaskType: "send_email",
        Data:     email,
        TraceID:  logger.GetTraceID(ctx),
        UserID:   logger.GetUserID(ctx),
        TenantID: logger.GetTenantID(ctx),
    }
    
    logger.InfoWithContext(ctx, p.logger, "发布邮件任务", 
        "task_id", taskID,
        "email_to", email.To,
        "trace_id", msg.TraceID)
    
    msgBytes, _ := json.Marshal(msg)
    return p.queue.Publish("email_tasks", msgBytes)
}
```

#### 任务消费

```go
// consumer.go
func (c *TaskConsumer) ConsumeEmailTask(msgData []byte) {
    var msg TaskMessage
    if err := json.Unmarshal(msgData, &msg); err != nil {
        c.logger.ErrorMsg("解析任务消息失败", "error", err)
        return
    }
    
    // 恢复上下文信息
    ctx := context.Background()
    ctx = logger.WithTraceID(ctx, msg.TraceID)
    ctx = logger.WithUserID(ctx, msg.UserID)
    ctx = logger.WithTenantID(ctx, msg.TenantID)
    
    // 为任务处理创建新的 Span
    spanCtx := logger.CreateSpan(ctx, "process_email_task")
    
    // 创建操作日志记录器
    opLogger := logger.NewOperationLogger(spanCtx, c.logger, "email_task")
    opLogger.SetTag("task_id", msg.TaskID)
    opLogger.SetTag("task_type", msg.TaskType)
    
    // 处理任务
    var emailData EmailData
    json.Unmarshal(msg.Data.([]byte), &emailData)
    
    if err := c.emailService.SendEmail(opLogger.GetContext(), &emailData); err != nil {
        opLogger.EndWithError(err, "email_to", emailData.To)
        return
    }
    
    opLogger.End("status", "sent", "email_to", emailData.To)
}
```

### 5. 批量处理

```go
// batch.go
func ProcessUserBatch(ctx context.Context, userIDs []string) error {
    // 为整个批次创建相关性ID
    correlationID := logger.GenerateCorrelationID()
    ctx = logger.WithCorrelationID(ctx, correlationID)
    
    // 创建批次操作记录器
    batchLogger := logger.NewOperationLogger(ctx, myLogger, "process_user_batch")
    batchLogger.SetMetric("total_users", len(userIDs))
    defer batchLogger.End()
    
    successCount := 0
    failedCount := 0
    
    for i, userID := range userIDs {
        // 每个用户处理有独立的 TraceID，但共享 CorrelationID
        userCtx, _ := logger.GetOrGenerateTraceID(ctx)
        userCtx = logger.CreateSpan(userCtx, "process_user")
        
        logger.InfoWithContext(userCtx, myLogger, "处理用户", 
            "index", i+1,
            "total", len(userIDs),
            "user_id", userID,
            "batch_correlation", correlationID)
        
        if err := processUser(userCtx, userID); err != nil {
            failedCount++
            logger.ErrorWithContext(userCtx, myLogger, "用户处理失败", 
                "user_id", userID, "error", err)
        } else {
            successCount++
        }
    }
    
    // 记录批次处理结果
    batchLogger.SetMetric("success_count", successCount)
    batchLogger.SetMetric("failed_count", failedCount)
    batchLogger.Info("批次处理完成", 
        "success", successCount, "failed", failedCount)
    
    return nil
}
```

## 进阶功能

### 操作日志记录器

操作日志记录器可以自动追踪操作的开始、结束时间和相关指标：

```go
func ComplexBusinessOperation(ctx context.Context) error {
    // 创建操作记录器
    opLogger := logger.NewOperationLogger(ctx, myLogger, "complex_operation")
    
    // 使用 defer 确保记录结束时间
    defer func() {
        if err := recover(); err != nil {
            opLogger.EndWithError(fmt.Errorf("panic: %v", err))
            panic(err)
        }
    }()
    
    // 设置操作标签
    opLogger.SetTag("operation_type", "data_migration")
    opLogger.SetTag("environment", "production")
    
    // 步骤1：验证数据
    opLogger.Info("开始数据验证")
    if err := validateData(opLogger.GetContext()); err != nil {
        opLogger.EndWithError(err, "step", "validation")
        return err
    }
    opLogger.SetMetric("validation_duration_ms", 100)
    
    // 步骤2：处理数据
    opLogger.Info("开始数据处理")
    processedCount, err := processData(opLogger.GetContext())
    if err != nil {
        opLogger.EndWithError(err, "step", "processing")
        return err
    }
    opLogger.SetMetric("processed_count", processedCount)
    
    // 步骤3：保存结果
    opLogger.Info("保存处理结果")
    if err := saveResults(opLogger.GetContext()); err != nil {
        opLogger.EndWithError(err, "step", "saving")
        return err
    }
    
    // 正常结束，记录成功指标
    opLogger.SetMetric("total_processed", processedCount)
    opLogger.End("status", "completed")
    
    return nil
}
```

### 相关性链追踪

用于关联业务相关的多个操作：

```go
func OrderWorkflow(ctx context.Context, orderData OrderData) error {
    // 创建工作流的相关性链
    chain, ctx := logger.CreateCorrelationChain(ctx)
    chain.SetTag("workflow", "order_processing")
    chain.SetTag("order_type", orderData.Type)
    defer logger.EndCorrelationChain(chain)
    
    // 步骤1：验证订单 - 独立的 Trace
    validateCtx, _ := logger.GetOrGenerateTraceID(ctx)
    if err := validateOrder(validateCtx, orderData); err != nil {
        logger.ErrorWithContext(validateCtx, myLogger, "订单验证失败", "error", err)
        return err
    }
    
    // 步骤2：库存检查 - 独立的 Trace
    inventoryCtx, _ := logger.GetOrGenerateTraceID(ctx)
    if err := checkInventory(inventoryCtx, orderData.Items); err != nil {
        logger.ErrorWithContext(inventoryCtx, myLogger, "库存检查失败", "error", err)
        return err
    }
    
    // 步骤3：支付处理 - 独立的 Trace
    paymentCtx, _ := logger.GetOrGenerateTraceID(ctx)
    if err := processPayment(paymentCtx, orderData.Payment); err != nil {
        logger.ErrorWithContext(paymentCtx, myLogger, "支付处理失败", "error", err)
        return err
    }
    
    chain.SetMetric("total_amount", orderData.TotalAmount)
    chain.SetTag("status", "completed")
    
    return nil
}
```

## 最佳实践

### 1. 性能优化

```go
// ✅ 推荐：批量日志记录时先创建带字段的 logger
func ProcessManyItems(ctx context.Context, items []Item) {
    contextLogger := logger.WithLogger(ctx, myLogger)
    
    for _, item := range items {
        // 避免重复提取上下文字段
        contextLogger.InfoKV("处理项目", "item_id", item.ID)
    }
}

// ❌ 避免：每次都提取字段
func ProcessManyItemsSlowly(ctx context.Context, items []Item) {
    for _, item := range items {
        // 每次调用都会提取上下文字段，性能较差
        logger.InfoWithContext(ctx, myLogger, "处理项目", "item_id", item.ID)
    }
}
```

### 2. 错误处理

```go
func RobustBusinessFunction(ctx context.Context) (err error) {
    opLogger := logger.NewOperationLogger(ctx, myLogger, "business_function")
    
    // 确保总是记录操作结果
    defer func() {
        if err != nil {
            opLogger.EndWithError(err)
        } else {
            opLogger.End("status", "success")
        }
    }()
    
    // 业务逻辑
    return doBusinessLogic(opLogger.GetContext())
}
```

### 3. goroutine 中的上下文传递

```go
func ProcessAsync(ctx context.Context, data ProcessData) {
    // 复制上下文信息，但创建独立的 context
    asyncCtx := context.Background()
    asyncCtx = logger.WithTraceID(asyncCtx, logger.GetTraceID(ctx))
    asyncCtx = logger.WithUserID(asyncCtx, logger.GetUserID(ctx))
    asyncCtx = logger.CreateSpan(asyncCtx, "async_processing")
    
    go func() {
        logger.InfoWithContext(asyncCtx, myLogger, "开始异步处理")
        
        if err := processDataAsync(asyncCtx, data); err != nil {
            logger.ErrorWithContext(asyncCtx, myLogger, "异步处理失败", "error", err)
            return
        }
        
        logger.InfoWithContext(asyncCtx, myLogger, "异步处理完成")
    }()
}
```

## API 参考

### ID 操作方法

```go
// TraceID
WithTraceID(ctx, id) context.Context
GetTraceID(ctx) string
GetOrGenerateTraceID(ctx) (context.Context, string)
GenerateTraceID() string

// SpanID
WithSpanID(ctx, id) context.Context
GetSpanID(ctx) string
GenerateSpanID() string
CreateSpan(ctx, operation) context.Context

// RequestID
WithRequestID(ctx, id) context.Context
GetRequestID(ctx) string
GenerateRequestID() string

// UserID
WithUserID(ctx, id) context.Context
GetUserID(ctx) string

// SessionID
WithSessionID(ctx, id) context.Context
GetSessionID(ctx) string

// CorrelationID
WithCorrelationID(ctx, id) context.Context
GetCorrelationID(ctx) string
GenerateCorrelationID() string

// TenantID
WithTenantID(ctx, id) context.Context
GetTenantID(ctx) string
```

### 日志方法

```go
// 带上下文的日志记录
InfoWithContext(ctx, logger, msg, kv...)
DebugWithContext(ctx, logger, msg, kv...)
WarnWithContext(ctx, logger, msg, kv...)
ErrorWithContext(ctx, logger, msg, kv...)

// 创建包含上下文字段的 logger
WithLogger(ctx, baseLogger) ILogger

// 提取上下文字段
ExtractFields(ctx) map[string]interface{}
```

### 相关性链

```go
// 创建相关性链
CreateCorrelationChain(ctx) (*CorrelationChain, context.Context)

// 结束相关性链
EndCorrelationChain(chain)

// 链操作
chain.SetTag(key, value)        // 设置字符串标签
chain.SetMetric(key, value)     // 设置任意类型指标
chain.GetDuration()             // 获取持续时间
```

### 操作日志记录器

```go
// 创建操作记录器
NewOperationLogger(ctx, logger, operation) *OperationLogger

// 记录器方法
opLogger.SetTag(key, value) *OperationLogger
opLogger.SetMetric(key, value) *OperationLogger
opLogger.Info(msg, kv...)
opLogger.Debug(msg, kv...)
opLogger.Warn(msg, kv...)
opLogger.Error(msg, kv...)
opLogger.End(kv...)
opLogger.EndWithError(err, kv...)
opLogger.GetContext() context.Context
```

## 常见问题

### Q: 如何集成到现有项目？

A: 逐步集成，先在入口处设置 TraceID，然后逐步添加其他字段：

```go
// 第一步：在 HTTP 入口添加 TraceID
func addTraceID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, traceID := logger.GetOrGenerateTraceID(r.Context())
        w.Header().Set("X-Trace-ID", traceID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 第二步：在需要的地方使用带上下文的日志
func businessHandler(w http.ResponseWriter, r *http.Request) {
    logger.InfoWithContext(r.Context(), myLogger, "处理业务请求")
    // 您的业务逻辑
}
```

### Q: 如何处理性能敏感的场景？

A: 使用字段提取一次，多次使用：

```go
func highFrequencyOperation(ctx context.Context) {
    // 一次性创建带字段的 logger
    contextLogger := logger.WithLogger(ctx, myLogger)
    
    // 在循环中重复使用，避免重复提取字段
    for i := 0; i < 10000; i++ {
        contextLogger.DebugKV("处理项目", "index", i)
    }
}
```

### Q: 如何在测试中使用？

A: 创建测试专用的上下文：

```go
func TestBusinessLogic(t *testing.T) {
    // 创建测试上下文
    ctx := context.Background()
    ctx = logger.WithTraceID(ctx, "test-trace-123")
    ctx = logger.WithUserID(ctx, "test-user")
    
    // 使用 mock logger
    mockLogger := &MockLogger{}
    
    // 测试业务逻辑
    err := BusinessFunction(ctx, mockLogger)
    
    assert.NoError(t, err)
    assert.Contains(t, mockLogger.Messages, "test-trace-123")
}
```

### ID 类型说明

#### TraceID (追踪ID)
- **用途**: 标识一次完整的分布式请求调用链
- **范围**: 跨服务、跨进程
- **生命周期**: 从请求开始到结束
- **特点**: 在整个调用链中保持不变

#### SpanID (操作ID)  
- **用途**: 标识单个操作或服务调用
- **范围**: 单个操作内
- **生命周期**: 从操作开始到结束
- **特点**: 一个 TraceID 包含多个 SpanID

#### RequestID (请求ID)
- **用途**: 标识单个 HTTP 请求或 RPC 调用
- **范围**: 单次网络请求
- **生命周期**: 请求开始到响应结束
- **特点**: 可返回给客户端用于问题排查

#### UserID (用户ID)
- **用途**: 标识操作的用户
- **范围**: 用户相关的所有操作
- **生命周期**: 用户会话期间
- **特点**: 用于权限控制和行为分析

#### SessionID (会话ID)
- **用途**: 标识用户登录会话
- **范围**: 从登录到登出
- **生命周期**: 会话期间
- **特点**: 一个用户可能有多个会话

#### CorrelationID (关联ID)
- **用途**: 关联业务相关的多个操作
- **范围**: 相关业务流程
- **生命周期**: 业务流程完成
- **特点**: 可跨越多个 TraceID

#### TenantID (租户ID)
- **用途**: 多租户系统中的租户标识
- **范围**: 租户的所有数据和操作
- **生命周期**: 租户存续期间
- **特点**: 用于数据隔离

## 基础 API

### ID 设置和获取

```go
// TraceID
ctx = WithTraceID(ctx, "trace-123")
traceID := GetTraceID(ctx)
ctx, traceID := GetOrGenerateTraceID(ctx)
newTraceID := GenerateTraceID()

// SpanID
ctx = WithSpanID(ctx, "span-456") 
spanID := GetSpanID(ctx)
newSpanID := GenerateSpanID()

// RequestID
ctx = WithRequestID(ctx, "req-789")
requestID := GetRequestID(ctx)
newRequestID := GenerateRequestID()

// UserID
ctx = WithUserID(ctx, "user-abc")
userID := GetUserID(ctx)

// SessionID
ctx = WithSessionID(ctx, "session-def")
sessionID := GetSessionID(ctx)

// CorrelationID
ctx = WithCorrelationID(ctx, "corr-ghi")
correlationID := GetCorrelationID(ctx)
newCorrelationID := GenerateCorrelationID()

// TenantID
ctx = WithTenantID(ctx, "tenant-jkl")
tenantID := GetTenantID(ctx)
```

### 字段提取

```go
// 提取所有上下文字段
fields := ExtractFields(ctx)
// 返回 map[string]interface{} 包含所有设置的字段

// 创建 Span
spanCtx := CreateSpan(ctx, "database_query")
```

### 日志记录

```go
// 直接记录日志（自动包含上下文字段）
InfoWithContext(ctx, logger, "操作完成", "result", "success")
DebugWithContext(ctx, logger, "调试信息", "key", "value")
ErrorWithContext(ctx, logger, "错误发生", "error", err)

// 创建带上下文字段的 logger
contextLogger := WithLogger(ctx, baseLogger)
contextLogger.InfoMsg("操作完成")
```

### 相关性链

```go
// 创建相关性链
chain, ctx := CreateCorrelationChain(ctx)

// 设置标签（字符串类型）
chain.SetTag("operation", "user_login")
chain.SetTag("service", "auth")

// 设置指标（任意类型）
chain.SetMetric("duration_ms", 150)
chain.SetMetric("success", true)

// 获取持续时间
duration := chain.GetDuration()

// 结束相关性链
EndCorrelationChain(chain)
```

### 操作日志记录器

```go
// 创建操作记录器
opLogger := NewOperationLogger(ctx, logger, "create_order")

// 设置标签和指标
opLogger.SetTag("order_type", "premium")
opLogger.SetMetric("amount", 99.99)

// 记录日志
opLogger.Info("订单创建开始")
opLogger.Debug("验证用户信息")

// 结束操作
opLogger.End()
// 或者带错误结束
opLogger.EndWithError(err)
```

## 使用场景

### 场景1: HTTP API 服务

#### 中间件设置

```go
// 追踪中间件
func TraceMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 从请求头提取或生成 TraceID
        traceID := r.Header.Get("X-Trace-ID")
        if traceID != "" {
            ctx = WithTraceID(ctx, traceID)
        } else {
            ctx, traceID = GetOrGenerateTraceID(ctx)
        }
        
        // 生成 RequestID
        requestID := GenerateRequestID()
        ctx = WithRequestID(ctx, requestID)
        
        // 设置响应头
        w.Header().Set("X-Trace-ID", traceID)
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 从 JWT Token 提取用户信息
        token := r.Header.Get("Authorization")
        if token != "" {
            claims := parseJWT(token)
            ctx = WithUserID(ctx, claims.UserID)
            ctx = WithTenantID(ctx, claims.TenantID)
            ctx = WithSessionID(ctx, claims.SessionID)
        }
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

#### 业务处理器

```go
func CreateOrderHandler(logger ILogger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 创建操作记录器
        opLogger := NewOperationLogger(ctx, logger, "create_order")
        defer opLogger.End()
        
        // 设置请求信息
        opLogger.SetTag("method", r.Method)
        opLogger.SetTag("path", r.URL.Path)
        
        // 解析请求
        var req CreateOrderRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            opLogger.EndWithError(err, "step", "parse_request")
            http.Error(w, "Invalid request", 400)
            return
        }
        
        // 调用服务
        order, err := orderService.Create(opLogger.GetContext(), &req)
        if err != nil {
            opLogger.EndWithError(err, "step", "service_call")
            http.Error(w, err.Error(), 500)
            return
        }
        
        opLogger.SetMetric("order_id", order.ID)
        opLogger.Info("订单创建成功")
        
        json.NewEncoder(w).Encode(order)
    }
}
```

### 场景2: 服务间调用

#### 客户端传递上下文

```go
func CallUserService(ctx context.Context, userID string) (*User, error) {
    req, _ := http.NewRequest("GET", "/users/"+userID, nil)
    
    // 传递追踪信息
    req.Header.Set("X-Trace-ID", GetTraceID(ctx))
    req.Header.Set("X-User-ID", GetUserID(ctx))
    req.Header.Set("X-Correlation-ID", GetCorrelationID(ctx))
    
    // 创建新的 SpanID 用于此次调用
    spanCtx := CreateSpan(ctx, "call_user_service")
    req.Header.Set("X-Span-ID", GetSpanID(spanCtx))
    
    InfoWithContext(spanCtx, logger, "调用用户服务", "user_id", userID)
    
    resp, err := client.Do(req)
    if err != nil {
        ErrorWithContext(spanCtx, logger, "用户服务调用失败", "error", err)
        return nil, err
    }
    
    return parseUserResponse(resp)
}
```

#### 服务端提取上下文

```go
func ExtractContextMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        // 提取上游传递的追踪信息
        if traceID := r.Header.Get("X-Trace-ID"); traceID != "" {
            ctx = WithTraceID(ctx, traceID)
        }
        if spanID := r.Header.Get("X-Span-ID"); spanID != "" {
            ctx = WithSpanID(ctx, spanID)
        }
        if userID := r.Header.Get("X-User-ID"); userID != "" {
            ctx = WithUserID(ctx, userID)
        }
        if correlationID := r.Header.Get("X-Correlation-ID"); correlationID != "" {
            ctx = WithCorrelationID(ctx, correlationID)
        }
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 场景3: 数据库操作

```go
type UserRepository struct {
    db     *sql.DB
    logger ILogger
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (*User, error) {
    // 创建数据库查询的 Span
    spanCtx := CreateSpan(ctx, "db_query_user")
    
    DebugWithContext(spanCtx, r.logger, "执行用户查询", 
        "user_id", userID,
        "table", "users")
    
    row := r.db.QueryRowContext(spanCtx, 
        "SELECT id, name, email FROM users WHERE id = ?", userID)
    
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        ErrorWithContext(spanCtx, r.logger, "数据库查询失败", 
            "error", err,
            "user_id", userID)
        return nil, err
    }
    
    InfoWithContext(spanCtx, r.logger, "用户查询成功", "user_id", userID)
    return &user, nil
}
```

### 场景4: 异步任务

#### 消息发布

```go
type TaskMessage struct {
    TaskID    string      `json:"task_id"`
    Data      interface{} `json:"data"`
    
    // 保存上下文信息
    TraceID   string `json:"trace_id"`
    UserID    string `json:"user_id"`
    TenantID  string `json:"tenant_id"`
}

func PublishTask(ctx context.Context, taskData interface{}) error {
    taskID := generateTaskID()
    
    msg := &TaskMessage{
        TaskID:  taskID,
        Data:    taskData,
        TraceID: GetTraceID(ctx),
        UserID:  GetUserID(ctx),
        TenantID: GetTenantID(ctx),
    }
    
    InfoWithContext(ctx, logger, "发布任务", 
        "task_id", taskID,
        "trace_id", msg.TraceID)
    
    return queue.Publish(msg)
}
```

#### 消息消费

```go
func ConsumeTask(msgData []byte) {
    var msg TaskMessage
    json.Unmarshal(msgData, &msg)
    
    // 恢复上下文
    ctx := context.Background()
    ctx = WithTraceID(ctx, msg.TraceID)
    ctx = WithUserID(ctx, msg.UserID)
    ctx = WithTenantID(ctx, msg.TenantID)
    
    // 创建新的 Span 用于任务处理
    spanCtx := CreateSpan(ctx, "process_task")
    
    opLogger := NewOperationLogger(spanCtx, logger, "async_task")
    opLogger.SetTag("task_id", msg.TaskID)
    
    if err := processTaskData(opLogger.GetContext(), msg.Data); err != nil {
        opLogger.EndWithError(err)
    } else {
        opLogger.End("status", "completed")
    }
}
```

### 场景5: 批量处理

```go
func ProcessBatch(ctx context.Context, items []Item) error {
    // 为整个批次创建相关性ID
    correlationID := GenerateCorrelationID()
    ctx = WithCorrelationID(ctx, correlationID)
    
    batchLogger := NewOperationLogger(ctx, logger, "batch_process")
    batchLogger.SetMetric("total_items", len(items))
    defer batchLogger.End()
    
    successCount := 0
    failedCount := 0
    
    for i, item := range items {
        // 每个 item 有独立的 TraceID，但共享 CorrelationID
        itemCtx, _ := GetOrGenerateTraceID(ctx)
        itemCtx = CreateSpan(itemCtx, "process_item")
        
        InfoWithContext(itemCtx, logger, "处理项目", 
            "index", i,
            "item_id", item.ID,
            "batch_correlation", correlationID)
        
        if err := processItem(itemCtx, item); err != nil {
            failedCount++
            ErrorWithContext(itemCtx, logger, "项目处理失败", 
                "error", err, "item_id", item.ID)
        } else {
            successCount++
        }
    }
    
    batchLogger.SetMetric("success_count", successCount)
    batchLogger.SetMetric("failed_count", failedCount)
    
    return nil
}
```

## 最佳实践

### 1. ID 传播原则

- **TraceID**: 在整个调用链中保持不变，跨服务传递
- **SpanID**: 每个独立操作生成新的 SpanID
- **CorrelationID**: 用于关联业务相关的多个请求

### 2. 性能优化

```go
// ✅ 推荐：一次性提取字段
contextLogger := WithLogger(ctx, baseLogger)
for i := 0; i < 1000; i++ {
    contextLogger.InfoKV("Processing", "index", i)
}

// ❌ 避免：重复提取字段
for i := 0; i < 1000; i++ {
    InfoWithContext(ctx, baseLogger, "Processing", "index", i)
}
```

### 3. 错误处理

```go
func BusinessFunction(ctx context.Context) error {
    opLogger := NewOperationLogger(ctx, logger, "business_op")
    defer func() {
        if err := recover(); err != nil {
            opLogger.EndWithError(fmt.Errorf("panic: %v", err))
            panic(err)
        }
    }()
    
    // 业务逻辑
    if err := doSomething(); err != nil {
        opLogger.EndWithError(err)
        return err
    }
    
    opLogger.End("status", "success")
    return nil
}
```

### 4. 测试支持

```go
func TestBusinessLogic(t *testing.T) {
    // 创建测试上下文
    ctx := context.Background()
    ctx = WithTraceID(ctx, "test-trace-123")
    ctx = WithUserID(ctx, "test-user")
    
    // 使用 mock logger
    mockLogger := &MockLogger{}
    
    err := BusinessFunction(ctx, mockLogger)
    
    assert.NoError(t, err)
    assert.Contains(t, mockLogger.Logs, "test-trace-123")
}
```

## 常见问题

### Q: 如何在 goroutine 中传递上下文？

```go
func ProcessAsync(ctx context.Context) {
    // 创建独立的上下文（避免被父 goroutine 取消）
    asyncCtx := context.Background()
    asyncCtx = WithTraceID(asyncCtx, GetTraceID(ctx))
    asyncCtx = WithUserID(asyncCtx, GetUserID(ctx))
    asyncCtx = CreateSpan(asyncCtx, "async_work")
    
    go func() {
        doAsyncWork(asyncCtx)
    }()
}
```

### Q: 如何处理没有上下文的情况？

```go
func SafeGetTraceID(ctx context.Context) string {
    if ctx == nil {
        return ""
    }
    return GetTraceID(ctx)
}

func EnsureTraceID(ctx context.Context) (context.Context, string) {
    if ctx == nil {
        ctx = context.Background()
    }
    return GetOrGenerateTraceID(ctx)
}
```

### Q: 如何在数据库事务中保持上下文？

```go
func WithTransaction(ctx context.Context, db *sql.DB, fn func(context.Context, *sql.Tx) error) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // 创建事务专用的 Span
    txCtx := CreateSpan(ctx, "database_transaction")
    
    if err := fn(txCtx, tx); err != nil {
        ErrorWithContext(txCtx, logger, "事务执行失败", "error", err)
        return err
    }
    
    if err := tx.Commit(); err != nil {
        ErrorWithContext(txCtx, logger, "事务提交失败", "error", err)
        return err
    }
    
    InfoWithContext(txCtx, logger, "事务提交成功")
    return nil
}
```