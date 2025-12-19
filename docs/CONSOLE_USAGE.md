# Console 风格日志功能

类似 JavaScript `console` 的日志分组、表格和计时器功能。

## 功能特性

### 1. 日志分组 (Console Group)

类似 JavaScript 的 `console.group()` 和 `console.groupCollapsed()`，支持嵌套分组。

#### 基本使用

```go
import "github.com/kamalyes/go-logger"

logger := logger.NewLogger(logger.DefaultConfig())
cg := logger.NewConsoleGroup()

// 开始分组
cg.Group("用户登录流程")
cg.Info("接收登录请求")
cg.Debug("验证用户名和密码")
cg.Info("登录成功，生成 Token")
cg.GroupEnd() // 结束分组
```

输出：
```
2025/12/19 ℹ️ [INFO] ▼ 用户登录流程
2025/12/19 ℹ️ [INFO]   接收登录请求
2025/12/19 🐛 [DEBUG]   验证用户名和密码
2025/12/19 ℹ️ [INFO]   登录成功，生成 Token
```

#### 嵌套分组

```go
cg.Group("订单处理系统")
cg.Info("开始处理订单批次")

  cg.Group("订单验证")
  cg.Info("检查库存")
  cg.Info("验证用户积分")
  cg.GroupEnd()

  cg.Group("支付处理")
  cg.Info("调用支付网关")
  cg.Info("支付成功")
  cg.GroupEnd()

cg.Info("订单处理完成")
cg.GroupEnd()
```

输出：
```
2025/12/19 ℹ️ [INFO] ▼ 订单处理系统
2025/12/19 ℹ️ [INFO]   开始处理订单批次
2025/12/19 ℹ️ [INFO]   ▼ 订单验证
2025/12/19 ℹ️ [INFO]     检查库存
2025/12/19 ℹ️ [INFO]     验证用户积分
2025/12/19 ℹ️ [INFO]   ▼ 支付处理
2025/12/19 ℹ️ [INFO]     调用支付网关
2025/12/19 ℹ️ [INFO]     支付成功
2025/12/19 ℹ️ [INFO]   订单处理完成
```

#### 带上下文的日志方法

```go
ctx := context.Background()

cg.InfoContext(ctx, "处理用户请求")
cg.DebugContext(ctx, "调试信息: %v", debugData)
cg.WarnContext(ctx, "警告: %s", warningMsg)
cg.ErrorContext(ctx, "错误: %v", err)
```

### 2. 表格展示 (Table)

类似 JavaScript 的 `console.table()`，支持多种数据格式。

#### 从 Map 切片创建表格

```go
users := []map[string]interface{}{
    {"ID": 1, "姓名": "张三", "年龄": 25, "部门": "技术部"},
    {"ID": 2, "姓名": "李四", "年龄": 30, "部门": "产品部"},
    {"ID": 3, "姓名": "王五", "年龄": 28, "部门": "技术部"},
}

cg.Group("用户列表")
cg.Table(users)
cg.GroupEnd()
```

输出：
```
2025/12/19 ℹ️ [INFO] ▼ 用户列表
2025/12/19 ℹ️ [INFO]
  ┌────┬──────┬──────┬────────┐
  │ ID │ 姓名   │ 年龄   │ 部门     │
  ├────┼──────┼──────┼────────┤
  │ 1  │ 张三   │ 25   │ 技术部   │
  │ 2  │ 李四   │ 30   │ 产品部   │
  │ 3  │ 王五   │ 28   │ 技术部   │
  └────┴──────┴──────┴────────┘
```

#### 从单个 Map 创建表格

```go
config := map[string]interface{}{
    "数据库类型":   "MySQL",
    "主机地址":    "localhost",
    "端口":      3306,
    "连接池大小":   100,
}

cg.Table(config)
```

输出：
```
  ┌───────────┬───────────┐
  │ Key       │ Value     │
  ├───────────┼───────────┤
  │ 数据库类型   │ MySQL     │
  │ 主机地址     │ localhost │
  │ 端口        │ 3306      │
  │ 连接池大小   │ 100       │
  └───────────┴───────────┘
```

#### 从字符串二维数组创建表格

```go
data := [][]string{
    {"服务名称", "状态", "响应时间", "错误率"},
    {"API Gateway", "运行中", "45ms", "0.01%"},
    {"Auth Service", "运行中", "23ms", "0.00%"},
    {"Database", "运行中", "12ms", "0.00%"},
}

cg.Table(data)
```

输出：
```
  ┌──────────────┬───────┬──────────┬─────────┐
  │ 服务名称       │ 状态    │ 响应时间   │ 错误率   │
  ├──────────────┼───────┼──────────┼─────────┤
  │ API Gateway  │ 运行中  │ 45ms     │ 0.01%   │
  │ Auth Service │ 运行中  │ 23ms     │ 0.00%   │
  │ Database     │ 运行中  │ 12ms     │ 0.00%   │
  └──────────────┴───────┴──────────┴─────────┘
```

### 3. 计时器 (Timer)

类似 JavaScript 的 `console.time()` 和 `console.timeEnd()`。

#### 基本计时

```go
timer := cg.Time("数据库查询")
// ... 执行操作 ...
timer.End() // 输出: ⏱️  数据库查询: 123.45ms
```

#### 中间检查点

```go
timer := cg.Time("文件处理")
time.Sleep(50 * time.Millisecond)
timer.Log("已处理 50%%") // 输出: ⏱️  文件处理: 50.00ms - 已处理 50%
time.Sleep(50 * time.Millisecond)
timer.Log("已处理 100%%") // 输出: ⏱️  文件处理: 100.00ms - 已处理 100%
timer.End() // 输出: ⏱️  文件处理: 100.00ms
```

#### 嵌套计时

```go
cg.Group("API 请求处理")
totalTimer := cg.Time("总耗时")

dbTimer := cg.Time("数据库查询")
time.Sleep(80 * time.Millisecond)
dbTimer.End()

cacheTimer := cg.Time("缓存更新")
time.Sleep(30 * time.Millisecond)
cacheTimer.End()

totalTimer.End()
cg.GroupEnd()
```

输出：
```
2025/12/19 ℹ️ [INFO] ▼ API 请求处理
2025/12/19 ℹ️ [INFO]   ⏱️  总耗时: 计时开始
2025/12/19 ℹ️ [INFO]   ⏱️  数据库查询: 计时开始
2025/12/19 ℹ️ [INFO]   ⏱️  数据库查询: 80.12ms
2025/12/19 ℹ️ [INFO]   ⏱️  缓存更新: 计时开始
2025/12/19 ℹ️ [INFO]   ⏱️  缓存更新: 30.45ms
2025/12/19 ℹ️ [INFO]   ⏱️  总耗时: 110.57ms
```

### 4. 全局便捷方法

无需创建 logger 实例，直接使用全局方法：

```go
import "github.com/kamalyes/go-logger"

// 全局分组
cg := logger.Group("全局分组")
cg.Info("这是全局方法")
cg.GroupEnd()

// 全局表格
logger.Table(map[string]interface{}{
    "功能":   "全局表格",
    "便捷性": "⭐⭐⭐⭐⭐",
})

// 全局计时器
logger.Time("全局任务")
// ... 执行操作 ...
logger.TimeLog("全局任务", "中间检查点")
// ... 继续操作 ...
logger.TimeEnd("全局任务")
```

## 完整示例

```go
package main

import (
    "time"
    "github.com/kamalyes/go-logger"
)

func main() {
    log := logger.NewLogger(logger.DefaultConfig())
    cg := log.NewConsoleGroup()

    cg.Group("🌐 API 请求: GET /api/users")
    requestTimer := cg.Time("请求总耗时")

    // 请求信息
    cg.Group("📋 请求信息")
    requestInfo := map[string]interface{}{
        "Method":     "GET",
        "Path":       "/api/users",
        "IP":         "192.168.1.100",
    }
    cg.Table(requestInfo)
    cg.GroupEnd()

    // 业务处理
    cg.Group("💼 业务处理")
    dbTimer := cg.Time("数据库查询")
    time.Sleep(85 * time.Millisecond)
    dbTimer.End()

    // 查询结果
    users := []map[string]interface{}{
        {"ID": 1, "Name": "张三", "Status": "Active"},
        {"ID": 2, "Name": "李四", "Status": "Active"},
    }
    cg.Table(users)
    cg.GroupEnd()

    requestTimer.End()
    cg.Info("✅ 请求处理完成")
    cg.GroupEnd()
}
```

## API 参考

### ConsoleGroup 方法

- `Group(label string, args ...interface{})` - 开始分组
- `GroupCollapsed(label string, args ...interface{})` - 开始折叠分组
- `GroupEnd()` - 结束当前分组
- `Info/Debug/Warn/Error(format string, args ...interface{})` - 记录日志
- `InfoContext/DebugContext/WarnContext/ErrorContext(ctx context.Context, format string, args ...interface{})` - 带上下文的日志
- `Table(data interface{})` - 显示表格
- `Time(label string) *Timer` - 创建计时器

### Timer 方法

- `End() time.Duration` - 结束计时并输出
- `Log(msg string, args ...interface{}) time.Duration` - 输出当前耗时
- `Elapsed() time.Duration` - 获取已过时间（不输出）

### 全局方法

- `Group(label string, args ...interface{}) *ConsoleGroup`
- `GroupCollapsed(label string, args ...interface{}) *ConsoleGroup`
- `Table(data interface{})`
- `Time(label string) *Timer`
- `TimeLog(label string, msg string, args ...interface{}) time.Duration`
- `TimeEnd(label string) time.Duration`

## 时间格式化

计时器会自动选择合适的时间单位：

- < 1μs: 纳秒 (ns)
- < 1ms: 微秒 (μs)
- < 1s: 毫秒 (ms)
- < 1m: 秒 (s)
- ≥ 1m: 分钟格式 (1m30s)

## 许可证

Copyright (c) 2025 by kamalyes, All Rights Reserved.
