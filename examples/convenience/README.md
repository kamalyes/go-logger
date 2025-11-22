# 便利函数示例

本示例演示了 go-logger 提供的三个便利函数的使用方法。

## 函数说明

### NewUltraFast()
返回极致性能的 `*UltraFastLogger`，适用于高并发场景：
```go
logger := logger.NewUltraFast()
logger.Info("极致性能日志")
```

### NewOptimized()
返回优化配置的 `*Logger`，平衡性能与功能：
```go
logger := logger.NewOptimized()
logger.Info("优化性能日志")
```

### New()
返回标准配置的 `*Logger`，提供完整功能：
```go
logger := logger.New()
logger.WithField("key", "value").Info("完整功能日志")
```

## 运行示例

```bash
go run main.go
```