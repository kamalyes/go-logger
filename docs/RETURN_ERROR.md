# 返回错误的日志方法 (Return Error Logging)

## 📖 概述

`go-logger` 提供了一套强大的返回错误日志方法，允许你在记录日志的同时返回格式化的错误对象。这个特性简化了错误处理流程，让代码更加简洁优雅。

## ✨ 核心优势

- 🎯 **简化错误处理**: 一行代码同时完成日志记录和错误返回
- 🔄 **保持错误链**: 返回的错误可以继续在调用链中传递
- 📝 **统一格式**: 错误信息与日志信息保持一致
- ⚡ **零性能开销**: 基于已有的日志方法，无额外性能损失
- 🎨 **多种场景**: 支持基本日志、上下文日志、键值对日志

## 🚀 快速开始

### 基本用法

```go
package main

import (
    "github.com/kamalyes/go-logger"
)

func processData(data string) error {
    if data == "" {
        // 记录错误日志并返回错误
        return logger.ErrorReturn("数据为空，无法处理")
    }
    
    // 业务逻辑...
    return nil
}

func connectDatabase(host string, port int) error {
    if host == "" {
        return logger.ErrorReturn("数据库连接失败: host=%s, port=%d", host, port)
    }
    
    // 连接逻辑...
    return nil
}
```

### 对比传统方式

**传统方式** ❌
```go
func oldWay(data string) error {
    if data == "" {
        err := fmt.Errorf("数据为空")
        logger.Error("数据为空")  // 重复的信息
        return err
    }
    return nil
}
```

**使用 Return 方法** ✅
```go
func newWay(data string) error {
    if data == "" {
        return logger.ErrorReturn("数据为空")  // 一行搞定！
    }
    return nil
}
```

## 📚 API 参考

### 基础返回错误方法

所有日志级别都支持返回错误：

```go
// 调试级别
err := log.DebugReturn("调试信息: %s", detail)

// 信息级别
err := log.InfoReturn("操作完成: %s", operation)

// 警告级别
err := log.WarnReturn("警告: 磁盘使用率 %d%%", usage)

// 错误级别
err := log.ErrorReturn("错误: %s", message)
```

### 带上下文的返回错误方法

支持在分布式系统中传递上下文信息：

```go
import "context"

func handleRequest(ctx context.Context, userID string) error {
    // 自动提取 TraceID、RequestID 等信息
    if userID == "" {
        return logger.ErrorCtxReturn(ctx, "用户ID不能为空")
    }
    
    // 带格式化参数
    if err := validateUser(userID); err != nil {
        return logger.ErrorCtxReturn(ctx, "用户验证失败: %v", err)
    }
    
    return nil
}
```

**支持的上下文方法:**
- `DebugCtxReturn(ctx, format, args...) error`
- `InfoCtxReturn(ctx, format, args...) error`
- `WarnCtxReturn(ctx, format, args...) error`
- `ErrorCtxReturn(ctx, format, args...) error`

### 键值对返回错误方法

适合结构化日志场景：

```go
func updateUser(userID string, name string, age int) error {
    if age < 0 {
        return logger.ErrorKVReturn(
            "用户年龄无效",
            "user_id", userID,
            "name", name,
            "age", age,
        )
    }
    
    return nil
}

func queryDatabase(table string, limit int) error {
    return logger.WarnKVReturn(
        "查询超时",
        "table", table,
        "limit", limit,
        "timeout", "5s",
    )
}
```

**支持的键值对方法:**
- `DebugKVReturn(msg, keysAndValues...) error`
- `InfoKVReturn(msg, keysAndValues...) error`
- `WarnKVReturn(msg, keysAndValues...) error`
- `ErrorKVReturn(msg, keysAndValues...) error`

## 💡 使用场景

### 1. 数据验证

```go
func validateOrder(order *Order) error {
    if order == nil {
        return logger.ErrorReturn("订单对象为空")
    }
    
    if order.Amount <= 0 {
        return logger.ErrorReturn("订单金额必须大于0: %.2f", order.Amount)
    }
    
    if order.UserID == "" {
        return logger.ErrorKVReturn(
            "订单用户ID为空",
            "order_id", order.ID,
            "amount", order.Amount,
        )
    }
    
    return nil
}
```

### 2. API 请求处理

```go
func handleAPIRequest(ctx context.Context, req *Request) error {
    // 验证请求
    if req.Token == "" {
        return logger.ErrorCtxReturn(ctx, "认证令牌缺失")
    }
    
    // 业务逻辑
    if err := processRequest(req); err != nil {
        return logger.ErrorCtxReturn(ctx, "请求处理失败: %v", err)
    }
    
    logger.InfoCtxReturn(ctx, "请求处理成功")
    return nil
}
```

### 3. 数据库操作

```go
func getUserByID(ctx context.Context, userID string) (*User, error) {
    if userID == "" {
        return nil, logger.ErrorKVReturn(
            "用户ID不能为空",
            "operation", "getUserByID",
        )
    }
    
    user, err := db.Query(ctx, userID)
    if err != nil {
        return nil, logger.ErrorCtxReturn(ctx, 
            "数据库查询失败: user_id=%s, error=%v", userID, err)
    }
    
    if user == nil {
        return nil, logger.WarnKVReturn(
            "用户不存在",
            "user_id", userID,
            "operation", "getUserByID",
        )
    }
    
    return user, nil
}
```

### 4. 业务流程控制

```go
func transferMoney(from, to string, amount float64) error {
    // 步骤1: 验证
    if amount <= 0 {
        return logger.ErrorReturn("转账金额必须大于0: %.2f", amount)
    }
    
    // 步骤2: 检查余额
    balance, err := getBalance(from)
    if err != nil {
        return logger.ErrorReturn("获取余额失败: %v", err)
    }
    
    if balance < amount {
        return logger.WarnKVReturn(
            "余额不足",
            "from", from,
            "balance", balance,
            "amount", amount,
        )
    }
    
    // 步骤3: 执行转账
    if err := executeTransfer(from, to, amount); err != nil {
        return logger.ErrorKVReturn(
            "转账执行失败",
            "from", from,
            "to", to,
            "amount", amount,
            "error", err.Error(),
        )
    }
    
    logger.InfoKVReturn("转账成功", "from", from, "to", to, "amount", amount)
    return nil
}
```

### 5. 错误链传递

```go
func processOrder(orderID string) error {
    // 第一层
    if err := validateOrderID(orderID); err != nil {
        return logger.ErrorReturn("订单验证失败: %v", err)
    }
    
    // 第二层
    if err := checkInventory(orderID); err != nil {
        return logger.ErrorReturn("库存检查失败: %v", err)
    }
    
    // 第三层
    if err := createShipment(orderID); err != nil {
        return logger.ErrorReturn("创建发货单失败: %v", err)
    }
    
    return nil
}

// 错误会在调用链中层层传递，每一层都会记录日志
```

## 🎯 全局方法

除了实例方法，还提供了全局便捷方法：

```go
import "github.com/kamalyes/go-logger"

func main() {
    // 全局方法
    if err := logger.ErrorReturn("全局错误: %s", "系统繁忙"); err != nil {
        // 处理错误
    }
    
    // 全局上下文方法
    ctx := context.Background()
    if err := logger.ErrorCtxReturn(ctx, "请求失败"); err != nil {
        // 处理错误
    }
    
    // 全局键值对方法
    if err := logger.ErrorKVReturn("操作失败", "code", 500); err != nil {
        // 处理错误
    }
}
```

**可用的全局方法:**
```go
// 基础方法
logger.DebugReturn(format, args...) error
logger.InfoReturn(format, args...) error
logger.WarnReturn(format, args...) error
logger.ErrorReturn(format, args...) error

// 上下文方法
logger.DebugCtxReturn(ctx, format, args...) error
logger.InfoCtxReturn(ctx, format, args...) error
logger.WarnCtxReturn(ctx, format, args...) error
logger.ErrorCtxReturn(ctx, format, args...) error

// 键值对方法
logger.DebugKVReturn(msg, keysAndValues...) error
logger.InfoKVReturn(msg, keysAndValues...) error
logger.WarnKVReturn(msg, keysAndValues...) error
logger.ErrorKVReturn(msg, keysAndValues...) error
```

## 🔧 配置和定制

### 使用自定义 Logger

```go
// 创建自定义 logger
log := logger.New().
    WithLevel(logger.DEBUG).
    WithColorful(true).
    WithShowCaller(true)

// 使用返回错误方法
if err := log.ErrorReturn("自定义日志错误"); err != nil {
    // 处理
}
```

### 在适配器中使用

```go
// StandardAdapter
adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
    Type:  logger.StandardAdapter,
    Level: logger.INFO,
})

err := adapter.ErrorReturn("适配器错误: %s", "连接失败")
```

### UltraFastLogger 支持

```go
// 极速日志器也支持返回错误
fastLog := logger.NewUltraFast()

err := fastLog.ErrorReturn("高性能错误日志: %d", 500)
```

## 📊 性能特点

返回错误的日志方法基于已有的日志方法实现，具有以下特点：

- ✅ **零额外开销**: 不增加额外的性能损失
- ✅ **相同的级别检查**: 继承原有的级别过滤机制
- ✅ **格式化复用**: 使用相同的格式化逻辑
- ✅ **内存优化**: 与普通日志方法相同的内存表现

## 🎓 最佳实践

### 1. 选择合适的日志级别

```go
// ✅ 使用 Error 记录真正的错误
return logger.ErrorReturn("数据库连接失败: %v", err)

// ✅ 使用 Warn 记录警告但不阻断流程
logger.WarnReturn("缓存未命中: key=%s", key)

// ✅ 使用 Info 记录重要信息
logger.InfoReturn("用户登录成功: user_id=%s", userID)

// ⚠️ 避免滥用 Debug
// logger.DebugReturn(...) 应该用于开发调试
```

### 2. 提供足够的上下文信息

```go
// ❌ 信息不足
return logger.ErrorReturn("操作失败")

// ✅ 提供详细信息
return logger.ErrorKVReturn(
    "用户更新操作失败",
    "user_id", userID,
    "operation", "update_profile",
    "error", err.Error(),
    "timestamp", time.Now(),
)
```

### 3. 使用上下文方法追踪请求

```go
// ✅ 在处理请求时使用 Ctx 方法
func handleRequest(ctx context.Context, req *Request) error {
    // 自动包含 TraceID 和 RequestID
    return logger.ErrorCtxReturn(ctx, "请求处理失败: %v", err)
}
```

### 4. 避免重复日志

```go
// ❌ 重复记录
func bad(data string) error {
    if data == "" {
        logger.Error("数据为空")
        return fmt.Errorf("数据为空")  // 重复了
    }
    return nil
}

// ✅ 只记录一次
func good(data string) error {
    if data == "" {
        return logger.ErrorReturn("数据为空")
    }
    return nil
}
```

## 🔗 相关文档

- [基础用法指南](USAGE.md)
- [上下文使用指南](CONTEXT_USAGE.md)
- [配置指南](CONFIGURATION.md)
- [示例代码](../examples/return_error/main.go)

## 📝 完整示例

查看 [examples/return_error/main.go](../examples/return_error/main.go) 获取完整的使用示例，包括：

- 基本返回错误示例
- 上下文返回错误示例
- 键值对返回错误示例
- 实际业务场景示例
- 全局方法示例

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个功能！
