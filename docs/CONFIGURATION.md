# Go Logger 配置指南

## 目录

- [🔧 基础配置](#-基础配置)
- [⚙️ 配置结构](#️-配置结构)
- [🌍 环境配置](#-环境配置)
- [📁 配置文件](#-配置文件)
- [🔄 动态配置](#-动态配置)
- [📊 监控配置](#-监控配置)
- [🎯 最佳实践](#-最佳实践)

## 🔧 基础配置

### 快速配置

```go
// 最简配置
logger := logger.New()

// 基础配置
config := &Config{
    Level:      INFO,
    Output:     os.Stdout,
    TimeFormat: TimeFormatStandard,
    Colorful:   true,
}
logger := logger.NewWithConfig(config)
```

### 性能配置

```go
// 高性能配置
config := &Config{
    Level:      INFO,
    TimeFormat: TimeFormatDisabled, // 禁用时间戳获得最高性能
    AsyncWrite: true,               // 异步写入
    BufferSize: 8192,              // 缓冲区大小
    PoolSize:   10,                // 对象池大小
}

// 极致性能配置
ultraConfig := &UltraFastConfig{
    Level:      INFO,
    TimeFormat: TimeFormatDisabled,
    Colorful:   false,
    SyncMode:   true,
}
```

### 调用者信息配置

`ShowCaller` 功能允许您在日志中显示调用者信息，包括文件名和行号，这对调试非常有用。

```go
// 启用调用者信息
adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
    Type:       logger.StandardAdapter,
    Level:      logger.DEBUG,
    Output:     os.Stdout,
    ShowCaller: true,  // 启用调用者信息显示
    Colorful:   true,
})

// 禁用调用者信息
adapter, _ := logger.NewStandardAdapter(&logger.AdapterConfig{
    Type:       logger.StandardAdapter,
    Level:      logger.INFO,
    Output:     os.Stdout,
    ShowCaller: false, // 禁用调用者信息显示
    Colorful:   true,
})
```

**输出效果对比：**

启用 `ShowCaller: true`：
```
2025/11/22 13:19:09 🐛 [DEBUG] [standard_adapter.go:120:Debug] 显示调用者信息的日志
2025/11/22 13:19:09 ℹ️ [INFO] [standard_adapter.go:127:Info] 文件名和行号会显示在日志中
```

禁用 `ShowCaller: false`：
```
2025/11/22 13:19:09 🐛 [DEBUG] 不显示调用者信息的日志
2025/11/22 13:19:09 ℹ️ [INFO] 没有文件名和行号信息
```

**性能建议：**
- 开发环境：建议启用 `ShowCaller: true` 以便调试
- 生产环境：建议禁用 `ShowCaller: false` 以获得更好的性能

## ⚙️ 配置结构

### 完整配置结构

```go
type Config struct {
    // 基础设置
    Level       LogLevel      `json:"level" yaml:"level"`
    Output      io.Writer     `json:"-" yaml:"-"`
    TimeFormat  TimeFormat    `json:"time_format" yaml:"time_format"`
    Colorful    bool          `json:"colorful" yaml:"colorful"`
    
    // 性能设置
    BufferSize     int  `json:"buffer_size" yaml:"buffer_size"`
    AsyncWrite     bool `json:"async_write" yaml:"async_write"`
    PoolSize       int  `json:"pool_size" yaml:"pool_size"`
    BatchSize      int  `json:"batch_size" yaml:"batch_size"`
    BatchTimeout   time.Duration `json:"batch_timeout" yaml:"batch_timeout"`
    
    // 企业功能
    EnableMemoryStats  bool `json:"enable_memory_stats" yaml:"enable_memory_stats"`
    EnableDistributed  bool `json:"enable_distributed" yaml:"enable_distributed"`
    EnableMetrics      bool `json:"enable_metrics" yaml:"enable_metrics"`
    EnableHooks        bool `json:"enable_hooks" yaml:"enable_hooks"`
    
    // 输出格式
    Format         FormatType `json:"format" yaml:"format"`
    TimestampKey   string     `json:"timestamp_key" yaml:"timestamp_key"`
    LevelKey       string     `json:"level_key" yaml:"level_key"`
    MessageKey     string     `json:"message_key" yaml:"message_key"`
    CallerKey      string     `json:"caller_key" yaml:"caller_key"`
    StacktraceKey  string     `json:"stacktrace_key" yaml:"stacktrace_key"`
    
    // 字段设置
    Fields         map[string]interface{} `json:"fields" yaml:"fields"`
    ContextFields  []string              `json:"context_fields" yaml:"context_fields"`
    
    // 调用者信息配置
    ShowCaller     bool `json:"show_caller" yaml:"show_caller"`         // 显示调用者信息（文件名和行号）
    CallerDepth    int  `json:"caller_depth" yaml:"caller_depth"`       // 调用者深度（默认：2）
    ShowStacktrace bool `json:"show_stacktrace" yaml:"show_stacktrace"` // 显示堆栈跟踪（仅错误日志）
    
    // 组件配置
    Adapters    []AdapterConfig    `json:"adapters" yaml:"adapters"`
    Hooks       []HookConfig       `json:"hooks" yaml:"hooks"`
    Middlewares []MiddlewareConfig `json:"middlewares" yaml:"middlewares"`
    
    // 监控配置
    Monitoring MonitoringConfig `json:"monitoring" yaml:"monitoring"`
}
```

### 适配器配置

```go
type AdapterConfig struct {
    Name     string      `json:"name" yaml:"name"`
    Type     string      `json:"type" yaml:"type"`
    Level    LogLevel    `json:"level" yaml:"level"`
    Enabled  bool        `json:"enabled" yaml:"enabled"`
    Config   interface{} `json:"config" yaml:"config"`
}

// 控制台适配器配置
type ConsoleAdapterConfig struct {
    Level         LogLevel    `json:"level" yaml:"level"`
    Colorful      bool        `json:"colorful" yaml:"colorful"`
    Format        FormatType  `json:"format" yaml:"format"`
    TimeFormat    TimeFormat  `json:"time_format" yaml:"time_format"`
    ShowCaller    bool        `json:"show_caller" yaml:"show_caller"` // 显示调用者信息（文件名:.行号）
    CallerDepth   int         `json:"caller_depth" yaml:"caller_depth"`
}

// 文件适配器配置
type FileAdapterConfig struct {
    Level           LogLevel      `json:"level" yaml:"level"`
    Path            string        `json:"path" yaml:"path"`
    MaxSize         int64         `json:"max_size" yaml:"max_size"`
    MaxFiles        int           `json:"max_files" yaml:"max_files"`
    MaxAge          time.Duration `json:"max_age" yaml:"max_age"`
    Compress        bool          `json:"compress" yaml:"compress"`
    Format          FormatType    `json:"format" yaml:"format"`
    AsyncWrite      bool          `json:"async_write" yaml:"async_write"`
    BufferSize      int           `json:"buffer_size" yaml:"buffer_size"`
    FlushInterval   time.Duration `json:"flush_interval" yaml:"flush_interval"`
    FlushThreshold  int           `json:"flush_threshold" yaml:"flush_threshold"`
}

// 网络适配器配置
type NetworkAdapterConfig struct {
    Level            LogLevel      `json:"level" yaml:"level"`
    Protocol         string        `json:"protocol" yaml:"protocol"` // tcp, udp, http
    Address          string        `json:"address" yaml:"address"`
    Timeout          time.Duration `json:"timeout" yaml:"timeout"`
    MaxConnections   int           `json:"max_connections" yaml:"max_connections"`
    MaxIdleTime      time.Duration `json:"max_idle_time" yaml:"max_idle_time"`
    BatchSize        int           `json:"batch_size" yaml:"batch_size"`
    BatchTimeout     time.Duration `json:"batch_timeout" yaml:"batch_timeout"`
    Compression      string        `json:"compression" yaml:"compression"`
    CompressionLevel int           `json:"compression_level" yaml:"compression_level"`
    MaxRetries       int           `json:"max_retries" yaml:"max_retries"`
    RetryDelay       time.Duration `json:"retry_delay" yaml:"retry_delay"`
    BackoffFactor    float64       `json:"backoff_factor" yaml:"backoff_factor"`
}
```

### 监控配置

```go
type MonitoringConfig struct {
    Memory      MemoryMonitoringConfig      `json:"memory" yaml:"memory"`
    Performance PerformanceMonitoringConfig `json:"performance" yaml:"performance"`
    IO          IOMonitoringConfig          `json:"io" yaml:"io"`
}

type MemoryMonitoringConfig struct {
    Enabled         bool          `json:"enabled" yaml:"enabled"`
    Threshold       float64       `json:"threshold" yaml:"threshold"`
    SampleInterval  time.Duration `json:"sample_interval" yaml:"sample_interval"`
    LeakDetection   bool          `json:"leak_detection" yaml:"leak_detection"`
    MaxHistorySize  int           `json:"max_history_size" yaml:"max_history_size"`
    GCPercent       int           `json:"gc_percent" yaml:"gc_percent"`
    MaxMemory       int64         `json:"max_memory" yaml:"max_memory"`
}

type PerformanceMonitoringConfig struct {
    Enabled             bool          `json:"enabled" yaml:"enabled"`
    LatencyThreshold    time.Duration `json:"latency_threshold" yaml:"latency_threshold"`
    ThroughputThreshold float64       `json:"throughput_threshold" yaml:"throughput_threshold"`
    SampleRate          float64       `json:"sample_rate" yaml:"sample_rate"`
}
```

## 🌍 环境配置

### 开发环境

```go
func NewDevelopmentConfig() *Config {
    return &Config{
        Level:      DEBUG,
        Output:     os.Stdout,
        Colorful:   true,
        TimeFormat: TimeFormatShort,
        Format:     FormatText,
        ShowCaller: true,
        
        // 开发环境不需要高性能优化
        AsyncWrite: false,
        BufferSize: 0,
        
        // 启用详细监控
        EnableMemoryStats: true,
        EnableMetrics:     true,
        
        Fields: map[string]interface{}{
            "env":     "development",
            "service": "my-app",
            "version": "dev",
        },
        
        Adapters: []AdapterConfig{
            {
                Name:    "console",
                Type:    "console",
                Level:   DEBUG,
                Enabled: true,
                Config: &ConsoleAdapterConfig{
                    Colorful: true,
                    Format:   FormatText,
                },
            },
        },
        
        Monitoring: MonitoringConfig{
            Memory: MemoryMonitoringConfig{
                Enabled:        true,
                Threshold:      90.0,
                SampleInterval: time.Second * 5,
            },
        },
    }
}
```

### 测试环境

```go
func NewTestingConfig() *Config {
    return &Config{
        Level:      INFO,
        Output:     os.Stdout,
        Colorful:   false,
        TimeFormat: TimeFormatRFC3339,
        Format:     FormatJSON,
        ShowCaller: false,
        
        // 测试环境使用内存适配器
        Adapters: []AdapterConfig{
            {
                Name:    "memory",
                Type:    "memory",
                Level:   DEBUG,
                Enabled: true,
                Config: &MemoryAdapterConfig{
                    MaxSize: 1000,
                    Format:  FormatJSON,
                },
            },
        },
        
        // 禁用一些监控功能
        EnableMemoryStats: false,
        EnableMetrics:     false,
        
        Fields: map[string]interface{}{
            "env":     "testing",
            "service": "my-app",
            "version": "test",
        },
    }
}
```

### 生产环境

```go
func NewProductionConfig() *Config {
    return &Config{
        Level:      INFO,
        Output:     os.Stdout,
        Colorful:   false,
        TimeFormat: TimeFormatRFC3339,
        Format:     FormatJSON,
        ShowCaller: false,
        
        // 生产环境高性能配置
        AsyncWrite:   true,
        BufferSize:   8192,
        PoolSize:     10,
        BatchSize:    100,
        BatchTimeout: time.Millisecond * 100,
        
        // 启用企业功能
        EnableMemoryStats: true,
        EnableDistributed: true,
        EnableMetrics:     true,
        EnableHooks:       true,
        
        Fields: map[string]interface{}{
            "env":     "production",
            "service": "my-app",
            "version": "1.0.0",
        },
        
        ContextFields: []string{
            "trace_id",
            "user_id",
            "session_id",
            "tenant_id",
        },
        
        Adapters: []AdapterConfig{
            // 文件适配器
            {
                Name:    "file",
                Type:    "file",
                Level:   INFO,
                Enabled: true,
                Config: &FileAdapterConfig{
                    Path:           "/var/log/app.log",
                    MaxSize:        100 * 1024 * 1024, // 100MB
                    MaxFiles:       10,
                    MaxAge:         30 * 24 * time.Hour, // 30天
                    Compress:       true,
                    AsyncWrite:     true,
                    BufferSize:     4096,
                    FlushInterval:  time.Second * 5,
                    FlushThreshold: 1000,
                },
            },
            
            // Elasticsearch适配器
            {
                Name:    "elasticsearch",
                Type:    "elasticsearch",
                Level:   WARN,
                Enabled: true,
                Config: &ElasticsearchAdapterConfig{
                    URLs:          []string{"http://es:9200"},
                    Index:         "logs-2024",
                    BufferSize:    1000,
                    FlushInterval: time.Second * 30,
                },
            },
        },
        
        Hooks: []HookConfig{
            // Prometheus指标钩子
            {
                Name:    "metrics",
                Type:    "prometheus",
                Enabled: true,
                Config: &PrometheusHookConfig{
                    Endpoint: "/metrics",
                    Namespace: "app",
                },
            },
            
            // 告警钩子
            {
                Name:    "alert",
                Type:    "webhook",
                Enabled: true,
                Config: &WebhookHookConfig{
                    URL:    "http://alert:8080/webhook",
                    Levels: []LogLevel{ERROR, FATAL},
                },
            },
        },
        
        Monitoring: MonitoringConfig{
            Memory: MemoryMonitoringConfig{
                Enabled:         true,
                Threshold:       85.0,
                SampleInterval:  time.Second * 5,
                LeakDetection:   true,
                MaxHistorySize:  100,
                GCPercent:       75,
                MaxMemory:       4 * 1024 * 1024 * 1024, // 4GB
            },
            Performance: PerformanceMonitoringConfig{
                Enabled:             true,
                LatencyThreshold:    time.Millisecond * 100,
                ThroughputThreshold: 1000.0,
                SampleRate:          0.1, // 10%采样
            },
        },
    }
}
```

## 📁 配置文件

### YAML 配置文件

```yaml
# config/logger.yaml
logger:
  # 基础设置
  level: info
  format: json
  time_format: rfc3339
  colorful: false
  show_caller: false
  
  # 性能设置
  async_write: true
  buffer_size: 8192
  pool_size: 10
  batch_size: 100
  batch_timeout: 100ms
  
  # 企业功能
  enable_memory_stats: true
  enable_distributed: true
  enable_metrics: true
  enable_hooks: true
  
  # 全局字段
  fields:
    service: "my-app"
    version: "1.2.0"
    environment: "production"
    datacenter: "us-west-2"
  
  # 上下文字段
  context_fields:
    - trace_id
    - span_id
    - user_id
    - session_id
    - tenant_id
    - correlation_id
  
  # 适配器配置
  adapters:
    - name: console
      type: console
      level: debug
      enabled: false  # 生产环境禁用控制台输出
      config:
        colorful: true
        format: text
        show_caller: true
        
    - name: file
      type: file
      level: info
      enabled: true
      config:
        path: "/var/log/app.log"
        max_size: 100MB
        max_files: 10
        max_age: 720h  # 30天
        compress: true
        async_write: true
        buffer_size: 4096
        flush_interval: 5s
        flush_threshold: 1000
        
    - name: elasticsearch
      type: elasticsearch
      level: warn
      enabled: true
      config:
        urls: 
          - "http://es1:9200"
          - "http://es2:9200"
        index: "logs-2024"
        type: "_doc"
        buffer_size: 1000
        flush_interval: 30s
        username: "elastic"
        password: "password"
        
    - name: kafka
      type: kafka
      level: error
      enabled: true
      config:
        brokers:
          - "kafka1:9092"
          - "kafka2:9092"
        topic: "error-logs"
        partition: -1
        compression: "gzip"
  
  # 钩子配置
  hooks:
    - name: metrics
      type: prometheus
      enabled: true
      config:
        endpoint: "/metrics"
        namespace: "app"
        subsystem: "logger"
        
    - name: alert
      type: webhook
      enabled: true
      config:
        url: "http://alert:8080/webhook"
        timeout: 10s
        levels: [error, fatal]
        batch_size: 10
        batch_timeout: 30s
        
    - name: audit
      type: audit
      enabled: true
      config:
        output: "/var/log/audit.log"
        levels: [info, warn, error, fatal]
  
  # 监控配置
  monitoring:
    memory:
      enabled: true
      threshold: 85.0
      sample_interval: 5s
      leak_detection: true
      max_history_size: 100
      gc_percent: 75
      max_memory: 4GB
      
    performance:
      enabled: true
      latency_threshold: 100ms
      throughput_threshold: 1000.0
      sample_rate: 0.1
      
    io:
      enabled: true
      disk_usage_threshold: 80.0
      iops_threshold: 1000
      latency_threshold: 100ms
```

### JSON 配置文件

```json
{
  "logger": {
    "level": "info",
    "format": "json",
    "time_format": "rfc3339",
    "colorful": false,
    "async_write": true,
    "buffer_size": 8192,
    "pool_size": 10,
    "enable_memory_stats": true,
    "enable_distributed": true,
    "enable_metrics": true,
    "fields": {
      "service": "my-app",
      "version": "1.2.0",
      "environment": "production"
    },
    "adapters": [
      {
        "name": "file",
        "type": "file",
        "level": "info",
        "enabled": true,
        "config": {
          "path": "/var/log/app.log",
          "max_size": 104857600,
          "max_files": 10,
          "compress": true
        }
      }
    ],
    "monitoring": {
      "memory": {
        "enabled": true,
        "threshold": 85.0,
        "sample_interval": "5s"
      }
    }
  }
}
```

### 加载配置文件

```go
// 从YAML文件加载
config, err := logger.LoadConfigFromYAML("config/logger.yaml")
if err != nil {
    log.Fatal("加载YAML配置失败:", err)
}

// 从JSON文件加载
config, err := logger.LoadConfigFromJSON("config/logger.json")
if err != nil {
    log.Fatal("加载JSON配置失败:", err)
}

// 自动检测格式
config, err := logger.LoadConfigFromFile("config/logger.yaml")
if err != nil {
    log.Fatal("加载配置失败:", err)
}

// 创建logger
log, err := logger.NewWithConfig(config)
if err != nil {
    log.Fatal("创建logger失败:", err)
}
```

## 🔄 动态配置

### 配置热重载

```go
// 创建配置管理器
configManager := logger.NewConfigManager()

// 监听配置文件变化
configManager.WatchFile("config/logger.yaml", func(newConfig *Config) {
    log.Info("检测到配置变化，正在重新加载...")
    
    // 验证新配置
    if err := newConfig.Validate(); err != nil {
        log.Error("新配置验证失败:", err)
        return
    }
    
    // 应用新配置
    if err := log.UpdateConfig(newConfig); err != nil {
        log.Error("配置更新失败:", err)
    } else {
        log.Info("配置更新成功")
    }
})

// 启动配置监听
if err := configManager.Start(); err != nil {
    log.Fatal("启动配置监听失败:", err)
}
defer configManager.Stop()
```

### HTTP API 配置更新

```go
// 配置更新API
http.HandleFunc("/admin/logger/config", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        // 获取当前配置
        currentConfig := log.GetConfig()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(currentConfig)
        
    case "PUT":
        // 更新配置
        var newConfig Config
        if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
            http.Error(w, "Invalid JSON: "+err.Error(), 400)
            return
        }
        
        // 验证配置
        if err := newConfig.Validate(); err != nil {
            http.Error(w, "Invalid config: "+err.Error(), 400)
            return
        }
        
        // 应用配置
        if err := log.UpdateConfig(&newConfig); err != nil {
            http.Error(w, "Update failed: "+err.Error(), 500)
            return
        }
        
        w.WriteHeader(200)
        json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
        
    case "PATCH":
        // 部分更新配置
        var updates map[string]interface{}
        if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
            http.Error(w, "Invalid JSON: "+err.Error(), 400)
            return
        }
        
        if err := log.UpdatePartialConfig(updates); err != nil {
            http.Error(w, "Update failed: "+err.Error(), 500)
            return
        }
        
        w.WriteHeader(200)
        json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
    }
})
```

### 运行时配置修改

```go
// 修改日志级别
log.SetLevel(DEBUG)

// 修改适配器配置
log.UpdateAdapterConfig("file", &FileAdapterConfig{
    Path:     "/tmp/debug.log",
    Level:    DEBUG,
    MaxSize:  50 * 1024 * 1024, // 50MB
    MaxFiles: 5,
})

// 添加新适配器
newAdapter, err := logger.CreateAdapter("tcp", &TCPAdapterConfig{
    Level:   WARN,
    Address: "log-server:514",
    Format:  FormatJSON,
})
if err == nil {
    log.AddAdapter("tcp", newAdapter)
}

// 移除适配器
log.RemoveAdapter("console")

// 启用/禁用钩子
log.EnableHook("alert")
log.DisableHook("metrics")

// 修改监控配置
log.UpdateMonitoringConfig(&MonitoringConfig{
    Memory: MemoryMonitoringConfig{
        Enabled:   true,
        Threshold: 90.0, // 提高内存阈值
    },
})
```

## 📊 监控配置

### 内存监控配置

```go
memoryConfig := MemoryMonitoringConfig{
    Enabled:         true,
    Threshold:       85.0,                    // 内存使用率阈值
    SampleInterval:  time.Second * 5,         // 采样间隔
    LeakDetection:   true,                    // 启用泄漏检测
    MaxHistorySize:  100,                     // 历史记录数量
    GCPercent:       75,                      // GC百分比
    MaxMemory:       4 * 1024 * 1024 * 1024, // 最大内存限制
    
    // 回调配置
    OnThresholdExceeded: func(info *MemoryInfo) {
        // 内存阈值超出时的处理
        log.Warn("内存使用率超出阈值", 
            "usage", info.MemoryUsage,
            "used", info.UsedMemory,
            "threshold", 85.0)
            
        // 可以触发告警或清理操作
        if info.MemoryUsage > 90.0 {
            runtime.GC() // 强制GC
            log.Info("强制执行GC")
        }
    },
    
    OnLeakDetected: func(report *LeakReport) {
        // 内存泄漏检测到时的处理
        log.Error("检测到内存泄漏",
            "trend", report.GrowthTrend,
            "rate", report.MemoryGrowthRate,
            "risk", report.RiskLevel)
            
        // 发送告警
        alertManager.SendAlert("memory_leak", report)
    },
}
```

### 性能监控配置

```go
performanceConfig := PerformanceMonitoringConfig{
    Enabled:             true,
    LatencyThreshold:    time.Millisecond * 100, // 延迟阈值
    ThroughputThreshold: 1000.0,                 // 吞吐量阈值
    SampleRate:          0.1,                    // 10%采样率
    
    // 回调配置
    OnLatencyExceeded: func(operation string, latency time.Duration) {
        log.Warn("操作延迟超标",
            "operation", operation,
            "latency", latency,
            "threshold", time.Millisecond*100)
    },
    
    OnThroughputExceeded: func(operation string, throughput float64) {
        log.Info("操作吞吐量超标",
            "operation", operation,
            "throughput", throughput,
            "threshold", 1000.0)
    },
}
```

## 🎯 最佳实践

### 配置分层

```go
// 基础配置
baseConfig := &Config{
    Level:      INFO,
    TimeFormat: TimeFormatRFC3339,
    Format:     FormatJSON,
}

// 环境特定配置
var envConfig *Config
switch os.Getenv("ENVIRONMENT") {
case "development":
    envConfig = NewDevelopmentConfig()
case "staging":
    envConfig = NewStagingConfig()
case "production":
    envConfig = NewProductionConfig()
default:
    envConfig = NewDevelopmentConfig()
}

// 合并配置
finalConfig := MergeConfigs(baseConfig, envConfig)
```

### 配置验证

```go
func (c *Config) Validate() error {
    // 验证基础设置
    if c.Level < 0 || c.Level > FATAL {
        return fmt.Errorf("invalid log level: %d", c.Level)
    }
    
    // 验证性能设置
    if c.BufferSize < 0 {
        return fmt.Errorf("buffer size cannot be negative: %d", c.BufferSize)
    }
    
    if c.PoolSize < 0 {
        return fmt.Errorf("pool size cannot be negative: %d", c.PoolSize)
    }
    
    // 验证适配器配置
    for i, adapter := range c.Adapters {
        if adapter.Name == "" {
            return fmt.Errorf("adapter[%d] name cannot be empty", i)
        }
        
        if adapter.Type == "" {
            return fmt.Errorf("adapter[%d] type cannot be empty", i)
        }
        
        // 验证适配器特定配置
        if err := ValidateAdapterConfig(adapter.Type, adapter.Config); err != nil {
            return fmt.Errorf("adapter[%d] config invalid: %v", i, err)
        }
    }
    
    return nil
}
```

### 配置安全

```go
// 敏感信息处理
type SecureConfig struct {
    *Config
    
    // 加密字段
    DatabasePassword string `json:"-" yaml:"-"` // 不序列化
    APIKey          string `json:"-" yaml:"-"`
}

func (c *SecureConfig) LoadFromEnv() {
    // 从环境变量加载敏感信息
    c.DatabasePassword = os.Getenv("DB_PASSWORD")
    c.APIKey = os.Getenv("API_KEY")
}

func (c *SecureConfig) Sanitize() *Config {
    // 返回不包含敏感信息的配置副本
    sanitized := *c.Config
    
    // 清理敏感字段
    if sanitized.Fields == nil {
        sanitized.Fields = make(map[string]interface{})
    }
    
    // 移除或脱敏敏感字段
    for key, value := range sanitized.Fields {
        if isSensitiveField(key) {
            sanitized.Fields[key] = "***" // 脱敏处理
        }
    }
    
    return &sanitized
}
```

### 配置优化

```go
// 性能优化配置
func OptimizeForPerformance(config *Config) *Config {
    optimized := *config
    
    // 根据系统资源调整配置
    numCPU := runtime.NumCPU()
    totalMemory := getTotalMemory()
    
    // 调整缓冲区大小
    if totalMemory > 8*1024*1024*1024 { // 8GB以上
        optimized.BufferSize = 16384
        optimized.PoolSize = numCPU * 2
    } else if totalMemory > 4*1024*1024*1024 { // 4GB以上
        optimized.BufferSize = 8192
        optimized.PoolSize = numCPU
    } else {
        optimized.BufferSize = 4096
        optimized.PoolSize = numCPU / 2
        if optimized.PoolSize < 2 {
            optimized.PoolSize = 2
        }
    }
    
    // 启用异步写入（如果支持）
    optimized.AsyncWrite = true
    
    // 调整批处理设置
    optimized.BatchSize = 100
    optimized.BatchTimeout = time.Millisecond * 100
    
    return &optimized
}

// 内存优化配置
func OptimizeForMemory(config *Config) *Config {
    optimized := *config
    
    // 减少缓冲区大小
    optimized.BufferSize = 1024
    optimized.PoolSize = 2
    
    // 禁用一些内存密集型功能
    optimized.EnableMemoryStats = false
    optimized.EnableMetrics = false
    
    // 减少历史记录
    if optimized.Monitoring.Memory.MaxHistorySize > 10 {
        optimized.Monitoring.Memory.MaxHistorySize = 10
    }
    
    return &optimized
}
```

---

更多配置相关信息请参考：

- [📚 使用指南](USAGE.md) - 完整使用指南
- [📊 性能详解](PERFORMANCE.md) - 性能配置优化
- [🔄 迁移指南](MIGRATION.md) - 配置迁移指南
- [🎯 Context使用指南](CONTEXT_USAGE.md) - 分布式配置