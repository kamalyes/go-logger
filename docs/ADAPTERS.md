# Go Logger 适配器系统指南

## 目录

- [🔌 适配器概述](#-适配器概述)
- [🎯 内置适配器](#-内置适配器)
- [📊 适配器管理](#-适配器管理)
- [🏗️ 自定义适配器](#️-自定义适配器)
- [🚀 高级功能](#-高级功能)
- [⚙️ 配置示例](#️-配置示例)

## 🔌 适配器概述

适配器系统是 go-logger 的核心特性之一，它允许日志输出到多种不同的目标。每个适配器都可以独立配置，支持不同的日志级别、格式和特殊功能。

### 适配器接口

```go
type IAdapter interface {
    // 基础日志方法
    Log(level LogLevel, message string, fields map[string]interface{}) error
    
    // 级别管理
    IsLevelEnabled(level LogLevel) bool
    SetLevel(level LogLevel)
    GetLevel() LogLevel
    
    // 生命周期管理
    Initialize() error
    Close() error
    Flush() error
    
    // 健康检查
    IsHealthy() bool
    
    // 元信息
    GetAdapterName() string
    GetAdapterVersion() string
    GetAdapterType() string
    
    // 字段管理
    WithField(key string, value interface{}) IAdapter
    WithFields(fields map[string]interface{}) IAdapter
    
    // 调用者信息
    SetShowCaller(show bool)
    IsShowCaller() bool
}
```

### 适配器类型

| 类型 | 描述 | 适用场景 |
|------|------|----------|
| Console | 控制台输出 | 开发调试 |
| File | 文件输出 | 本地日志存储 |
| TCP/UDP | 网络传输 | 远程日志收集 |
| HTTP | RESTful API | 日志服务集成 |
| Elasticsearch | 搜索引擎 | 日志分析查询 |
| Redis | 内存数据库 | 缓存和队列 |
| Kafka | 消息队列 | 大数据处理 |
| Database | 数据库存储 | 结构化存储 |
| Email | 邮件通知 | 错误告警 |
| Webhook | HTTP回调 | 第三方集成 |

## 🎯 内置适配器

### 1. Console 适配器

控制台适配器用于将日志输出到标准输出或标准错误。

```go
// 基础配置
consoleConfig := &ConsoleAdapterConfig{
    Level:         DEBUG,
    Colorful:      true,
    Format:        FormatText,
    TimeFormat:    TimeFormatShort,
    ShowCaller:    true,
    CallerDepth:   4,
    Output:        os.Stdout, // 或 os.Stderr
}

// 创建适配器
adapter, err := logger.CreateConsoleAdapter(consoleConfig)
if err != nil {
    log.Fatal(err)
}

// 高级配置
advancedConfig := &ConsoleAdapterConfig{
    Level:      INFO,
    Colorful:   true,
    Format:     FormatText,
    TimeFormat: TimeFormatStandard,
    
    // 颜色自定义
    Colors: ColorConfig{
        Debug: "\033[36m", // 青色
        Info:  "\033[32m", // 绿色
        Warn:  "\033[33m", // 黄色
        Error: "\033[31m", // 红色
        Fatal: "\033[35m", // 紫色
        Reset: "\033[0m",  // 重置
    },
    
    // 格式模板
    Template: "{{.Time}} [{{.Level}}] {{if .Caller}}{{.Caller}} {{end}}{{.Message}}{{if .Fields}} {{.Fields}}{{end}}\n",
}
```

### 2. File 适配器

文件适配器支持日志轮转、压缩和异步写入。

```go
// 基础文件配置
fileConfig := &FileAdapterConfig{
    Level:    INFO,
    Path:     "/var/log/app.log",
    Format:   FormatJSON,
    MaxSize:  100 * 1024 * 1024, // 100MB
    MaxFiles: 10,
    MaxAge:   30 * 24 * time.Hour, // 30天
    Compress: true,
}

// 高性能文件配置
performanceConfig := &FileAdapterConfig{
    Level:    INFO,
    Path:     "/var/log/app.log",
    Format:   FormatJSON,
    
    // 文件轮转
    MaxSize:         100 * 1024 * 1024, // 100MB
    MaxFiles:        10,
    MaxAge:          30 * 24 * time.Hour,
    Compress:        true,
    CompressLevel:   6, // gzip压缩级别
    
    // 性能优化
    AsyncWrite:      true,
    BufferSize:      8192,
    FlushInterval:   time.Second * 5,
    FlushThreshold:  1000,
    PreallocSize:    10 * 1024 * 1024, // 预分配10MB
    
    // 高级特性
    DirectIO:        false, // 直接I/O
    SyncStrategy:    "batch", // none, immediate, batch
    Permissions:     0644,
    
    // 文件命名
    RotationPattern: "2006-01-02-15", // 按小时轮转
    LinkName:        "/var/log/current.log", // 软链接
}

// 创建文件适配器
adapter, err := logger.CreateFileAdapter(fileConfig)
if err != nil {
    log.Fatal(err)
}
```

### 3. Network 适配器

#### TCP 适配器

```go
tcpConfig := &TCPAdapterConfig{
    Level:     WARN,
    Address:   "log-server:514",
    Network:   "tcp",
    Timeout:   time.Second * 10,
    Format:    FormatJSON,
    
    // 连接池
    MaxConnections: 5,
    MaxIdleTime:    time.Minute * 5,
    KeepAlive:      true,
    KeepAlivePeriod: time.Second * 30,
    
    // 重连机制
    EnableReconnect: true,
    ReconnectDelay:  time.Second,
    MaxReconnects:   10,
    
    // 缓冲和批处理
    BufferSize:    4096,
    BatchSize:     100,
    BatchTimeout:  time.Second,
    
    // 安全设置
    TLS: &TLSConfig{
        Enabled:            true,
        InsecureSkipVerify: false,
        CertFile:           "/etc/ssl/client.crt",
        KeyFile:            "/etc/ssl/client.key",
        CAFile:             "/etc/ssl/ca.crt",
    },
}

adapter, err := logger.CreateTCPAdapter(tcpConfig)
```

#### UDP 适配器

```go
udpConfig := &UDPAdapterConfig{
    Level:         INFO,
    Address:       "log-server:514",
    MaxPacketSize: 1024,
    Format:        FormatJSON,
    
    // 缓冲设置
    BufferSize:   4096,
    BatchSize:    50,
    BatchTimeout: time.Millisecond * 500,
    
    // 网络设置
    LocalAddr: "0.0.0.0:0",
    TTL:       64,
}

adapter, err := logger.CreateUDPAdapter(udpConfig)
```

#### HTTP 适配器

```go
httpConfig := &HTTPAdapterConfig{
    Level:  WARN,
    URL:    "http://log-server:8080/logs",
    Method: "POST",
    Format: FormatJSON,
    
    // 认证
    Headers: map[string]string{
        "Authorization": "Bearer your-token",
        "Content-Type":  "application/json",
        "User-Agent":    "go-logger/1.0",
    },
    
    // 超时设置
    Timeout:        time.Second * 30,
    ConnectTimeout: time.Second * 10,
    WriteTimeout:   time.Second * 15,
    ReadTimeout:    time.Second * 15,
    
    // 性能设置
    BufferSize:     1000,
    FlushInterval:  time.Second * 60,
    MaxConnections: 10,
    IdleTimeout:    time.Minute * 5,
    
    // 重试机制
    MaxRetries:    3,
    RetryDelay:    time.Second,
    BackoffFactor: 2.0,
    
    // 压缩
    Compression: "gzip",
    CompressionLevel: 6,
}

adapter, err := logger.CreateHTTPAdapter(httpConfig)
```

### 4. Elasticsearch 适配器

```go
esConfig := &ElasticsearchAdapterConfig{
    Level:  INFO,
    URLs:   []string{"http://es1:9200", "http://es2:9200"},
    Index:  "logs-2024",
    Type:   "_doc",
    Format: FormatJSON,
    
    // 认证
    Username: "elastic",
    Password: "password",
    
    // 批处理
    BufferSize:    1000,
    FlushInterval: time.Second * 30,
    FlushTimeout:  time.Second * 10,
    
    // 索引设置
    IndexPattern:   "logs-2006-01-02", // 按日期分割索引
    IndexTemplate:  "log-template",
    Pipeline:       "log-pipeline",
    RoutingField:   "service",
    
    // 文档设置
    DocumentType:   "_doc",
    DocumentID:     "", // 空值表示自动生成
    TimestampField: "@timestamp",
    
    // 性能优化
    Compression:    true,
    MaxRetries:     3,
    RetryDelay:     time.Second,
    HealthCheck:    time.Minute,
    
    // 映射配置
    Mapping: map[string]interface{}{
        "properties": map[string]interface{}{
            "@timestamp": map[string]interface{}{
                "type": "date",
            },
            "level": map[string]interface{}{
                "type": "keyword",
            },
            "message": map[string]interface{}{
                "type": "text",
                "analyzer": "standard",
            },
            "service": map[string]interface{}{
                "type": "keyword",
            },
        },
    },
}

adapter, err := logger.CreateElasticsearchAdapter(esConfig)
```

### 5. Redis 适配器

```go
redisConfig := &RedisAdapterConfig{
    Level:  DEBUG,
    Format: FormatJSON,
    
    // 连接设置
    Addr:     "redis:6379",
    Password: "",
    DB:       0,
    
    // 存储设置
    Key:       "logs",
    KeyType:   "list", // list, stream, pubsub
    MaxLength: 10000,
    
    // Stream 模式设置（当 KeyType = "stream" 时）
    StreamConfig: &RedisStreamConfig{
        StreamKey:    "logs-stream",
        ConsumerGroup: "log-processors",
        MaxLength:     10000,
        Approximate:   true,
    },
    
    // Pub/Sub 模式设置（当 KeyType = "pubsub" 时）
    PubSubConfig: &RedisPubSubConfig{
        Channel: "logs-channel",
        Pattern: "logs.*",
    },
    
    // 连接池设置
    PoolSize:     10,
    MinIdleConns: 5,
    MaxRetries:   3,
    
    // 超时设置
    DialTimeout:  time.Second * 5,
    ReadTimeout:  time.Second * 3,
    WriteTimeout: time.Second * 3,
    PoolTimeout:  time.Second * 4,
    IdleTimeout:  time.Minute * 5,
    
    // 批处理
    BatchSize:    100,
    BatchTimeout: time.Second,
}

adapter, err := logger.CreateRedisAdapter(redisConfig)
```

### 6. Kafka 适配器

```go
kafkaConfig := &KafkaAdapterConfig{
    Level:   INFO,
    Format:  FormatJSON,
    
    // 集群设置
    Brokers: []string{"kafka1:9092", "kafka2:9092", "kafka3:9092"},
    Topic:   "logs",
    
    // 分区设置
    Partition:      -1, // 自动分区
    PartitionKey:   "service", // 分区键字段
    PartitionFunc:  "hash", // hash, random, round-robin
    
    // 生产者设置
    ProducerConfig: &KafkaProducerConfig{
        RequiredAcks:    1, // 0=no acks, 1=leader ack, -1=all acks
        Timeout:         time.Second * 10,
        Compression:     "gzip", // none, gzip, snappy, lz4, zstd
        MaxMessageSize:  1000000, // 1MB
        BatchSize:       16384,
        BatchTimeout:    time.Millisecond * 100,
        RetryMax:        3,
        RetryBackoff:    time.Millisecond * 100,
        
        // 幂等配置
        Idempotent: true,
        
        // 事务配置（可选）
        TransactionID: "log-producer-1",
    },
    
    // 安全设置
    Security: &KafkaSecurityConfig{
        Protocol: "SASL_PLAINTEXT", // PLAINTEXT, SSL, SASL_PLAINTEXT, SASL_SSL
        SASL: &SASLConfig{
            Mechanism: "PLAIN", // PLAIN, SCRAM-SHA-256, SCRAM-SHA-512
            Username:  "user",
            Password:  "password",
        },
        TLS: &TLSConfig{
            Enabled:    true,
            CertFile:   "/etc/ssl/client.crt",
            KeyFile:    "/etc/ssl/client.key",
            CAFile:     "/etc/ssl/ca.crt",
        },
    },
    
    // 消息格式
    MessageFormat: &KafkaMessageFormat{
        KeyField:       "trace_id", // 消息键字段
        TimestampField: "@timestamp",
        HeaderFields:   []string{"service", "version", "env"},
    },
}

adapter, err := logger.CreateKafkaAdapter(kafkaConfig)
```

### 7. Database 适配器

```go
dbConfig := &DatabaseAdapterConfig{
    Level:      INFO,
    Format:     FormatJSON,
    
    // 数据库连接
    Driver:     "postgres", // mysql, postgres, sqlite3
    DSN:        "postgres://user:pass@localhost/logs?sslmode=disable",
    
    // 表设置
    TableName:  "logs",
    AutoCreate: true,
    
    // 字段映射
    FieldMapping: map[string]string{
        "timestamp": "created_at",
        "level":     "log_level",
        "message":   "log_message",
        "trace_id":  "trace_id",
        "user_id":   "user_id",
    },
    
    // 连接池
    MaxOpenConns:    25,
    MaxIdleConns:    5,
    ConnMaxLifetime: time.Hour,
    ConnMaxIdleTime: time.Minute * 10,
    
    // 批处理
    BatchSize:    100,
    BatchTimeout: time.Second * 5,
    
    // 数据保留
    RetentionDays: 90,
    CleanupCron:   "0 2 * * *", // 每天凌晨2点清理
    
    // 表结构（自动创建时使用）
    Schema: &DatabaseSchema{
        Columns: []ColumnDefinition{
            {Name: "id", Type: "SERIAL PRIMARY KEY"},
            {Name: "created_at", Type: "TIMESTAMP"},
            {Name: "log_level", Type: "VARCHAR(10)"},
            {Name: "log_message", Type: "TEXT"},
            {Name: "trace_id", Type: "VARCHAR(64)"},
            {Name: "user_id", Type: "VARCHAR(64)"},
            {Name: "metadata", Type: "JSONB"},
        },
        Indexes: []IndexDefinition{
            {Name: "idx_created_at", Columns: []string{"created_at"}},
            {Name: "idx_log_level", Columns: []string{"log_level"}},
            {Name: "idx_trace_id", Columns: []string{"trace_id"}},
        },
    },
}

adapter, err := logger.CreateDatabaseAdapter(dbConfig)
```

## 📊 适配器管理

### 多适配器管理器

```go
// 创建适配器管理器
manager := logger.NewAdapterManager()

// 添加多个适配器
adapters := []struct {
    name   string
    config AdapterConfig
}{
    {"console", ConsoleAdapterConfig{Level: DEBUG, Colorful: true}},
    {"file", FileAdapterConfig{Level: INFO, Path: "/var/log/app.log"}},
    {"elasticsearch", ElasticsearchAdapterConfig{Level: WARN, URLs: []string{"http://es:9200"}}},
}

for _, a := range adapters {
    adapter, err := logger.CreateAdapter(a.name, a.config)
    if err != nil {
        log.Printf("创建适配器 %s 失败: %v", a.name, err)
        continue
    }
    
    if err := manager.AddAdapter(a.name, adapter); err != nil {
        log.Printf("添加适配器 %s 失败: %v", a.name, err)
        adapter.Close()
        continue
    }
    
    log.Printf("适配器 %s 添加成功", a.name)
}

// 使用管理器记录日志
manager.Debug("调试信息")     // 只发送到 console (DEBUG级别)
manager.Info("普通信息")      // 发送到 console 和 file
manager.Error("错误信息")     // 发送到所有适配器

// 广播到所有适配器
manager.Broadcast(INFO, "重要信息", map[string]interface{}{
    "event_type": "system_notification",
    "severity":   "high",
})
```

### 适配器状态管理

```go
// 健康检查
health := manager.HealthCheck()
for name, healthy := range health {
    if healthy {
        fmt.Printf("✅ 适配器 %s 正常\n", name)
    } else {
        fmt.Printf("❌ 适配器 %s 异常\n", name)
        
        // 尝试重新初始化异常的适配器
        if adapter, exists := manager.GetAdapter(name); exists {
            if err := adapter.Initialize(); err != nil {
                log.Printf("重新初始化适配器 %s 失败: %v", name, err)
                manager.RemoveAdapter(name)
            } else {
                log.Printf("适配器 %s 重新初始化成功", name)
            }
        }
    }
}

// 获取适配器统计信息
stats := manager.GetStatistics()
for name, stat := range stats {
    fmt.Printf("适配器 %s 统计: 消息数=%d, 错误数=%d, 延迟=%.2fms\n",
        name, stat.MessageCount, stat.ErrorCount, stat.AvgLatency.Seconds()*1000)
}

// 动态调整适配器级别
manager.SetAdapterLevel("file", DEBUG)
manager.SetAdapterLevel("elasticsearch", ERROR)

// 刷新所有适配器
manager.FlushAll()

// 优雅关闭
manager.CloseAll()
```

### 适配器路由

```go
// 创建带路由功能的管理器
routingManager := logger.NewRoutingAdapterManager()

// 定义路由规则
rules := []RoutingRule{
    {
        Name:      "debug_to_console",
        Condition: func(level LogLevel, msg string, fields map[string]interface{}) bool {
            return level == DEBUG
        },
        Adapters: []string{"console"},
    },
    {
        Name:      "error_to_alert",
        Condition: func(level LogLevel, msg string, fields map[string]interface{}) bool {
            return level >= ERROR
        },
        Adapters: []string{"elasticsearch", "email", "slack"},
    },
    {
        Name:      "service_specific",
        Condition: func(level LogLevel, msg string, fields map[string]interface{}) bool {
            if service, ok := fields["service"]; ok {
                return service == "payment-service"
            }
            return false
        },
        Adapters: []string{"database", "audit-file"},
    },
    {
        Name:      "default",
        Condition: func(level LogLevel, msg string, fields map[string]interface{}) bool {
            return true // 默认规则
        },
        Adapters: []string{"file"},
    },
}

// 添加路由规则
for _, rule := range rules {
    routingManager.AddRoutingRule(rule)
}

// 使用路由管理器
routingManager.Log(DEBUG, "调试信息", nil)           // 只发送到 console
routingManager.Log(ERROR, "错误信息", nil)           // 发送到 elasticsearch, email, slack, file
routingManager.Log(INFO, "支付完成", map[string]interface{}{
    "service": "payment-service",
})  // 发送到 database, audit-file, file
```

## 🏗️ 自定义适配器

### 基础适配器实现

```go
// 自定义适配器结构
type CustomAdapter struct {
    level       LogLevel
    config      *CustomConfig
    client      *CustomClient
    fields      map[string]interface{}
    showCaller  bool
    
    // 统计信息
    messageCount int64
    errorCount   int64
    lastError    error
    lastMessage  time.Time
    
    // 并发控制
    mutex sync.RWMutex
}

// 自定义配置
type CustomConfig struct {
    Level       LogLevel `json:"level" yaml:"level"`
    Endpoint    string   `json:"endpoint" yaml:"endpoint"`
    APIKey      string   `json:"api_key" yaml:"api_key"`
    Format      string   `json:"format" yaml:"format"`
    BufferSize  int      `json:"buffer_size" yaml:"buffer_size"`
    Timeout     time.Duration `json:"timeout" yaml:"timeout"`
    MaxRetries  int      `json:"max_retries" yaml:"max_retries"`
}

// 实现 IAdapter 接口
func (a *CustomAdapter) Log(level LogLevel, message string, fields map[string]interface{}) error {
    if !a.IsLevelEnabled(level) {
        return nil
    }
    
    a.mutex.Lock()
    defer a.mutex.Unlock()
    
    // 合并字段
    combinedFields := make(map[string]interface{})
    for k, v := range a.fields {
        combinedFields[k] = v
    }
    for k, v := range fields {
        combinedFields[k] = v
    }
    
    // 创建日志条目
    entry := &LogEntry{
        Timestamp: time.Now(),
        Level:     level.String(),
        Message:   message,
        Fields:    combinedFields,
        Source:    a.config.Endpoint,
    }
    
    // 发送日志
    if err := a.client.Send(entry); err != nil {
        a.errorCount++
        a.lastError = err
        return err
    }
    
    a.messageCount++
    a.lastMessage = time.Now()
    return nil
}

func (a *CustomAdapter) IsLevelEnabled(level LogLevel) bool {
    a.mutex.RLock()
    defer a.mutex.RUnlock()
    return level >= a.level
}

func (a *CustomAdapter) SetLevel(level LogLevel) {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    a.level = level
}

func (a *CustomAdapter) GetLevel() LogLevel {
    a.mutex.RLock()
    defer a.mutex.RUnlock()
    return a.level
}

func (a *CustomAdapter) Initialize() error {
    client, err := NewCustomClient(a.config)
    if err != nil {
        return fmt.Errorf("初始化客户端失败: %v", err)
    }
    
    a.client = client
    return nil
}

func (a *CustomAdapter) Close() error {
    if a.client != nil {
        return a.client.Close()
    }
    return nil
}

func (a *CustomAdapter) Flush() error {
    if a.client != nil {
        return a.client.Flush()
    }
    return nil
}

func (a *CustomAdapter) IsHealthy() bool {
    if a.client == nil {
        return false
    }
    return a.client.IsConnected()
}

func (a *CustomAdapter) GetAdapterName() string {
    return "custom"
}

func (a *CustomAdapter) GetAdapterVersion() string {
    return "1.0.0"
}

func (a *CustomAdapter) GetAdapterType() string {
    return "third-party"
}

func (a *CustomAdapter) WithField(key string, value interface{}) IAdapter {
    newAdapter := *a
    newFields := make(map[string]interface{})
    for k, v := range a.fields {
        newFields[k] = v
    }
    newFields[key] = value
    newAdapter.fields = newFields
    return &newAdapter
}

func (a *CustomAdapter) WithFields(fields map[string]interface{}) IAdapter {
    newAdapter := *a
    newFields := make(map[string]interface{})
    for k, v := range a.fields {
        newFields[k] = v
    }
    for k, v := range fields {
        newFields[k] = v
    }
    newAdapter.fields = newFields
    return &newAdapter
}

func (a *CustomAdapter) SetShowCaller(show bool) {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    a.showCaller = show
}

func (a *CustomAdapter) IsShowCaller() bool {
    a.mutex.RLock()
    defer a.mutex.RUnlock()
    return a.showCaller
}
```

### 注册自定义适配器

```go
// 注册适配器工厂函数
func init() {
    logger.RegisterAdapterFactory("custom", func(config interface{}) (logger.IAdapter, error) {
        customConfig, ok := config.(*CustomConfig)
        if !ok {
            return nil, fmt.Errorf("invalid config type for custom adapter")
        }
        
        adapter := &CustomAdapter{
            level:  customConfig.Level,
            config: customConfig,
            fields: make(map[string]interface{}),
        }
        
        if err := adapter.Initialize(); err != nil {
            return nil, err
        }
        
        return adapter, nil
    })
}

// 使用自定义适配器
config := &CustomConfig{
    Level:      INFO,
    Endpoint:   "http://custom-log-server:8080",
    APIKey:     "your-api-key",
    Format:     "json",
    BufferSize: 1000,
    Timeout:    time.Second * 30,
    MaxRetries: 3,
}

adapter, err := logger.CreateAdapter("custom", config)
if err != nil {
    log.Fatal("创建自定义适配器失败:", err)
}

// 添加到管理器
manager.AddAdapter("custom", adapter)
```

## 🚀 高级功能

### 适配器中间件

```go
// 适配器中间件接口
type AdapterMiddleware interface {
    Process(entry *LogEntry, next func(*LogEntry) error) error
}

// 字段增强中间件
type FieldEnhancerMiddleware struct {
    enhancer func(map[string]interface{}) map[string]interface{}
}

func (m *FieldEnhancerMiddleware) Process(entry *LogEntry, next func(*LogEntry) error) error {
    // 增强字段
    entry.Fields = m.enhancer(entry.Fields)
    return next(entry)
}

// 过滤中间件
type FilterMiddleware struct {
    filter func(*LogEntry) bool
}

func (m *FilterMiddleware) Process(entry *LogEntry, next func(*LogEntry) error) error {
    if !m.filter(entry) {
        return nil // 过滤掉这条日志
    }
    return next(entry)
}

// 速率限制中间件
type RateLimitMiddleware struct {
    limiter *rate.Limiter
}

func (m *RateLimitMiddleware) Process(entry *LogEntry, next func(*LogEntry) error) error {
    if !m.limiter.Allow() {
        return fmt.Errorf("rate limit exceeded")
    }
    return next(entry)
}

// 在适配器中使用中间件
type MiddlewareAdapter struct {
    IAdapter
    middlewares []AdapterMiddleware
}

func (a *MiddlewareAdapter) Log(level LogLevel, message string, fields map[string]interface{}) error {
    entry := &LogEntry{
        Timestamp: time.Now(),
        Level:     level.String(),
        Message:   message,
        Fields:    fields,
    }
    
    return a.processMiddlewares(entry, 0)
}

func (a *MiddlewareAdapter) processMiddlewares(entry *LogEntry, index int) error {
    if index >= len(a.middlewares) {
        // 所有中间件处理完毕，调用底层适配器
        return a.IAdapter.Log(
            ParseLogLevel(entry.Level),
            entry.Message,
            entry.Fields,
        )
    }
    
    middleware := a.middlewares[index]
    return middleware.Process(entry, func(e *LogEntry) error {
        return a.processMiddlewares(e, index+1)
    })
}
```

### 适配器插件系统

```go
// 插件接口
type AdapterPlugin interface {
    Name() string
    Version() string
    Init(config map[string]interface{}) error
    CreateAdapter(config interface{}) (IAdapter, error)
    Shutdown() error
}

// 插件管理器
type PluginManager struct {
    plugins map[string]AdapterPlugin
    mutex   sync.RWMutex
}

func NewPluginManager() *PluginManager {
    return &PluginManager{
        plugins: make(map[string]AdapterPlugin),
    }
}

func (pm *PluginManager) LoadPlugin(pluginPath string) error {
    // 动态加载插件 (示例使用 plugin 包)
    p, err := plugin.Open(pluginPath)
    if err != nil {
        return err
    }
    
    // 查找插件符号
    symbol, err := p.Lookup("Plugin")
    if err != nil {
        return err
    }
    
    // 类型断言
    adapterPlugin, ok := symbol.(AdapterPlugin)
    if !ok {
        return fmt.Errorf("插件不实现 AdapterPlugin 接口")
    }
    
    // 初始化插件
    if err := adapterPlugin.Init(nil); err != nil {
        return err
    }
    
    pm.mutex.Lock()
    pm.plugins[adapterPlugin.Name()] = adapterPlugin
    pm.mutex.Unlock()
    
    return nil
}

func (pm *PluginManager) CreateAdapter(pluginName string, config interface{}) (IAdapter, error) {
    pm.mutex.RLock()
    plugin, exists := pm.plugins[pluginName]
    pm.mutex.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("插件 %s 不存在", pluginName)
    }
    
    return plugin.CreateAdapter(config)
}
```

### 适配器性能监控

```go
// 性能监控适配器包装器
type PerformanceMonitoringAdapter struct {
    IAdapter
    metrics *AdapterMetrics
    monitor *PerformanceMonitor
}

type AdapterMetrics struct {
    MessageCount    int64
    ErrorCount      int64
    TotalLatency    time.Duration
    MinLatency      time.Duration
    MaxLatency      time.Duration
    LastMessageTime time.Time
    LastErrorTime   time.Time
    LastError       error
}

func NewPerformanceMonitoringAdapter(adapter IAdapter, monitor *PerformanceMonitor) *PerformanceMonitoringAdapter {
    return &PerformanceMonitoringAdapter{
        IAdapter: adapter,
        metrics:  &AdapterMetrics{},
        monitor:  monitor,
    }
}

func (p *PerformanceMonitoringAdapter) Log(level LogLevel, message string, fields map[string]interface{}) error {
    start := time.Now()
    
    err := p.IAdapter.Log(level, message, fields)
    
    latency := time.Since(start)
    
    // 更新指标
    atomic.AddInt64(&p.metrics.MessageCount, 1)
    p.metrics.TotalLatency += latency
    p.metrics.LastMessageTime = time.Now()
    
    if p.metrics.MinLatency == 0 || latency < p.metrics.MinLatency {
        p.metrics.MinLatency = latency
    }
    if latency > p.metrics.MaxLatency {
        p.metrics.MaxLatency = latency
    }
    
    if err != nil {
        atomic.AddInt64(&p.metrics.ErrorCount, 1)
        p.metrics.LastErrorTime = time.Now()
        p.metrics.LastError = err
    }
    
    // 记录到性能监控器
    p.monitor.RecordLatency(p.GetAdapterName(), latency)
    if err != nil {
        p.monitor.RecordError(p.GetAdapterName(), err)
    }
    
    return err
}

func (p *PerformanceMonitoringAdapter) GetMetrics() *AdapterMetrics {
    return &AdapterMetrics{
        MessageCount:    atomic.LoadInt64(&p.metrics.MessageCount),
        ErrorCount:      atomic.LoadInt64(&p.metrics.ErrorCount),
        TotalLatency:    p.metrics.TotalLatency,
        MinLatency:      p.metrics.MinLatency,
        MaxLatency:      p.metrics.MaxLatency,
        LastMessageTime: p.metrics.LastMessageTime,
        LastErrorTime:   p.metrics.LastErrorTime,
        LastError:       p.metrics.LastError,
    }
}

func (p *PerformanceMonitoringAdapter) GetAvgLatency() time.Duration {
    messageCount := atomic.LoadInt64(&p.metrics.MessageCount)
    if messageCount == 0 {
        return 0
    }
    return p.metrics.TotalLatency / time.Duration(messageCount)
}
```

## ⚙️ 配置示例

### 完整适配器配置

```yaml
# config/adapters.yaml
adapters:
  # 控制台适配器 - 开发环境
  console:
    type: console
    level: debug
    enabled: true
    config:
      colorful: true
      format: text
      time_format: short
      show_caller: true
      
  # 文件适配器 - 应用日志
  app_file:
    type: file
    level: info
    enabled: true
    config:
      path: "/var/log/app/app.log"
      max_size: 100MB
      max_files: 10
      max_age: 720h
      compress: true
      async_write: true
      buffer_size: 8192
      flush_interval: 5s
      
  # 错误文件适配器 - 错误日志
  error_file:
    type: file
    level: error
    enabled: true
    config:
      path: "/var/log/app/error.log"
      max_size: 50MB
      max_files: 5
      max_age: 720h
      compress: true
      
  # Elasticsearch - 日志检索
  elasticsearch:
    type: elasticsearch
    level: warn
    enabled: true
    config:
      urls: ["http://elasticsearch:9200"]
      index: "app-logs-2024"
      buffer_size: 1000
      flush_interval: 30s
      username: "elastic"
      password: "password"
      
  # Kafka - 大数据处理
  kafka:
    type: kafka
    level: info
    enabled: true
    config:
      brokers: ["kafka1:9092", "kafka2:9092"]
      topic: "app-logs"
      batch_size: 100
      batch_timeout: 1s
      compression: "gzip"
      
  # HTTP - 告警系统
  alert_webhook:
    type: http
    level: error
    enabled: true
    config:
      url: "http://alert-system:8080/webhook"
      method: "POST"
      headers:
        Authorization: "Bearer alert-token"
        Content-Type: "application/json"
      timeout: 30s
      buffer_size: 100
      flush_interval: 60s
      
  # Database - 审计日志
  audit_db:
    type: database
    level: info
    enabled: true
    config:
      driver: "postgres"
      dsn: "postgres://user:pass@db:5432/logs?sslmode=disable"
      table_name: "audit_logs"
      batch_size: 50
      batch_timeout: 10s
      
  # Redis - 实时监控
  redis_stream:
    type: redis
    level: debug
    enabled: true
    config:
      addr: "redis:6379"
      key_type: "stream"
      stream_config:
        stream_key: "logs-stream"
        max_length: 10000
      batch_size: 100
      batch_timeout: 1s
```

### 适配器路由配置

```yaml
# config/routing.yaml
routing:
  rules:
    # 调试信息只输出到控制台
    - name: "debug_to_console"
      priority: 1
      condition:
        level: debug
      adapters: ["console"]
      
    # 错误信息发送到多个目标
    - name: "errors_to_multiple"
      priority: 2
      condition:
        level: ["error", "fatal"]
      adapters: ["error_file", "elasticsearch", "alert_webhook"]
      
    # 特定服务的日志
    - name: "payment_service_logs"
      priority: 3
      condition:
        fields:
          service: "payment"
      adapters: ["app_file", "audit_db", "kafka"]
      
    # 安全相关日志
    - name: "security_logs"
      priority: 4
      condition:
        fields:
          category: "security"
      adapters: ["audit_db", "elasticsearch", "alert_webhook"]
      
    # 默认路由
    - name: "default"
      priority: 999
      condition:
        level: ["info", "warn"]
      adapters: ["app_file", "redis_stream"]
```

---

更多适配器相关信息请参考：

- [📚 使用指南](USAGE.md) - 完整使用指南
- [🔧 配置指南](CONFIGURATION.md) - 适配器配置详解
- [📊 性能详解](PERFORMANCE.md) - 适配器性能优化
- [🔄 迁移指南](MIGRATION.md) - 适配器迁移指南