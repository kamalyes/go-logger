/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\config.go
 * @Description: 日志配置
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"io"
	"os"
	"strings"
	"time"
)

// FormatType 输出格式类型
type FormatType string

const (
	FormatText FormatType = "text"
	FormatJSON FormatType = "json"
	FormatXML  FormatType = "xml"
	FormatCSV  FormatType = "csv"
)

// TimeFormat 时间格式类型
type TimeFormat string

const (
	TimeFormatStandard TimeFormat = "2006-01-02 15:04:05"
	TimeFormatISO8601  TimeFormat = "2006-01-02T15:04:05Z"
	TimeFormatRFC3339  TimeFormat = "2006-01-02T15:04:05Z07:00"
	TimeFormatUnix     TimeFormat = "unix"
	TimeFormatDisabled TimeFormat = "disabled"
)

// Config 主配置结构
type Config struct {
	// 基础设置
	Level      LogLevel   `json:"level" yaml:"level"`
	Output     io.Writer  `json:"-" yaml:"-"`
	TimeFormat TimeFormat `json:"time_format" yaml:"time_format"`
	Colorful   bool       `json:"colorful" yaml:"colorful"`

	// 性能设置
	BufferSize   int           `json:"buffer_size" yaml:"buffer_size"`
	AsyncWrite   bool          `json:"async_write" yaml:"async_write"`
	PoolSize     int           `json:"pool_size" yaml:"pool_size"`
	BatchSize    int           `json:"batch_size" yaml:"batch_size"`
	BatchTimeout time.Duration `json:"batch_timeout" yaml:"batch_timeout"`

	// 企业功能
	EnableMemoryStats bool `json:"enable_memory_stats" yaml:"enable_memory_stats"`
	EnableDistributed bool `json:"enable_distributed" yaml:"enable_distributed"`
	EnableMetrics     bool `json:"enable_metrics" yaml:"enable_metrics"`
	EnableHooks       bool `json:"enable_hooks" yaml:"enable_hooks"`

	// 输出格式
	Format        FormatType `json:"format" yaml:"format"`
	TimestampKey  string     `json:"timestamp_key" yaml:"timestamp_key"`
	LevelKey      string     `json:"level_key" yaml:"level_key"`
	MessageKey    string     `json:"message_key" yaml:"message_key"`
	CallerKey     string     `json:"caller_key" yaml:"caller_key"`
	StacktraceKey string     `json:"stacktrace_key" yaml:"stacktrace_key"`

	// 字段设置
	Fields        map[string]interface{} `json:"fields" yaml:"fields"`
	ContextFields []string               `json:"context_fields" yaml:"context_fields"`

	// 调用者信息配置
	ShowCaller     bool `json:"show_caller" yaml:"show_caller"`
	CallerDepth    int  `json:"caller_depth" yaml:"caller_depth"`
	ShowStacktrace bool `json:"show_stacktrace" yaml:"show_stacktrace"`

	// 组件配置
	Adapters []AdapterConfig `json:"adapters" yaml:"adapters"`

	// 监控配置
	Monitoring MonitoringConfig `json:"monitoring" yaml:"monitoring"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled     bool                        `json:"enabled" yaml:"enabled"`
	Memory      MemoryMonitoringConfig      `json:"memory" yaml:"memory"`
	Performance PerformanceMonitoringConfig `json:"performance" yaml:"performance"`
	IO          IOMonitoringConfig          `json:"io" yaml:"io"`
}

// MemoryMonitoringConfig 内存监控配置
type MemoryMonitoringConfig struct {
	Enabled   bool          `json:"enabled" yaml:"enabled"`
	Interval  time.Duration `json:"interval" yaml:"interval"`
	Threshold int64         `json:"threshold" yaml:"threshold"`
}

// PerformanceMonitoringConfig 性能监控配置
type PerformanceMonitoringConfig struct {
	Enabled         bool    `json:"enabled" yaml:"enabled"`
	TrackLatency    bool    `json:"track_latency" yaml:"track_latency"`
	TrackThroughput bool    `json:"track_throughput" yaml:"track_throughput"`
	SampleRate      float64 `json:"sample_rate" yaml:"sample_rate"`
}

// IOMonitoringConfig IO监控配置
type IOMonitoringConfig struct {
	Enabled     bool `json:"enabled" yaml:"enabled"`
	TrackWrites bool `json:"track_writes" yaml:"track_writes"`
	TrackReads  bool `json:"track_reads" yaml:"track_reads"`
	TrackErrors bool `json:"track_errors" yaml:"track_errors"`
}

// UltraFastConfig 极致性能配置
type UltraFastConfig struct {
	Level      LogLevel   `json:"level" yaml:"level"`
	TimeFormat TimeFormat `json:"time_format" yaml:"time_format"`
	Colorful   bool       `json:"colorful" yaml:"colorful"`
	SyncMode   bool       `json:"sync_mode" yaml:"sync_mode"`
	Output     io.Writer  `json:"-" yaml:"-"`
}

// NewUltraFastConfig 创建极致性能配置
func NewUltraFastConfig() *UltraFastConfig {
	return &UltraFastConfig{
		Level:      INFO,
		TimeFormat: TimeFormatDisabled,
		Colorful:   false,
		SyncMode:   true,
		Output:     os.Stdout,
	}
}

// ToConfig 转换为标准配置
func (c *UltraFastConfig) ToConfig() *LogConfig {
	config := DefaultConfig()
	config.Level = c.Level
	config.TimeFormat = string(c.TimeFormat)
	config.Colorful = c.Colorful
	config.Output = c.Output
	config.ShowCaller = false
	return config
}

// FileAdapterConfig 文件适配器配置
type FileAdapterConfig struct {
	Path         string        `json:"path" yaml:"path"`
	MaxSize      int64         `json:"max_size" yaml:"max_size"` // 单位：字节
	MaxBackups   int           `json:"max_backups" yaml:"max_backups"`
	MaxAge       time.Duration `json:"max_age" yaml:"max_age"` // 单位：天
	Compress     bool          `json:"compress" yaml:"compress"`
	LocalTime    bool          `json:"local_time" yaml:"local_time"`
	BufferSize   int           `json:"buffer_size" yaml:"buffer_size"`
	SyncInterval time.Duration `json:"sync_interval" yaml:"sync_interval"`
}

// NetworkAdapterConfig 网络适配器配置
type NetworkAdapterConfig struct {
	Protocol        string            `json:"protocol" yaml:"protocol"` // tcp, udp, http, https
	Address         string            `json:"address" yaml:"address"`
	Port            int               `json:"port" yaml:"port"`
	Timeout         time.Duration     `json:"timeout" yaml:"timeout"`
	RetryAttempts   int               `json:"retry_attempts" yaml:"retry_attempts"`
	RetryInterval   time.Duration     `json:"retry_interval" yaml:"retry_interval"`
	TLS             TLSConfig         `json:"tls" yaml:"tls"`
	Authentication  AuthConfig        `json:"authentication" yaml:"authentication"`
	Headers         map[string]string `json:"headers" yaml:"headers"`
	CompressionType string            `json:"compression_type" yaml:"compression_type"` // gzip, lz4, none
}

// TLSConfig TLS配置
type TLSConfig struct {
	Enabled            bool   `json:"enabled" yaml:"enabled"`
	CertFile           string `json:"cert_file" yaml:"cert_file"`
	KeyFile            string `json:"key_file" yaml:"key_file"`
	CAFile             string `json:"ca_file" yaml:"ca_file"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify" yaml:"insecure_skip_verify"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Type     string            `json:"type" yaml:"type"` // basic, bearer, api_key
	Username string            `json:"username" yaml:"username"`
	Password string            `json:"password" yaml:"password"`
	Token    string            `json:"token" yaml:"token"`
	APIKey   string            `json:"api_key" yaml:"api_key"`
	Headers  map[string]string `json:"headers" yaml:"headers"`
}

// ElasticsearchAdapterConfig Elasticsearch适配器配置
type ElasticsearchAdapterConfig struct {
	Addresses       []string      `json:"addresses" yaml:"addresses"`
	Index           string        `json:"index" yaml:"index"`
	DocType         string        `json:"doc_type" yaml:"doc_type"`
	Username        string        `json:"username" yaml:"username"`
	Password        string        `json:"password" yaml:"password"`
	CloudID         string        `json:"cloud_id" yaml:"cloud_id"`
	APIKey          string        `json:"api_key" yaml:"api_key"`
	Timeout         time.Duration `json:"timeout" yaml:"timeout"`
	RetryOnConflict int           `json:"retry_on_conflict" yaml:"retry_on_conflict"`
	BufferSize      int           `json:"buffer_size" yaml:"buffer_size"`
	FlushInterval   time.Duration `json:"flush_interval" yaml:"flush_interval"`
	TLS             TLSConfig     `json:"tls" yaml:"tls"`
}

// HookConfig 钩子配置
type HookConfig struct {
	Name    string                 `json:"name" yaml:"name"`
	Type    string                 `json:"type" yaml:"type"` // prometheus, webhook, slack, email
	Level   LogLevel               `json:"level" yaml:"level"`
	Enabled bool                   `json:"enabled" yaml:"enabled"`
	Config  map[string]interface{} `json:"config" yaml:"config"`
}

// PrometheusHookConfig Prometheus钩子配置
type PrometheusHookConfig struct {
	Namespace   string            `json:"namespace" yaml:"namespace"`
	Subsystem   string            `json:"subsystem" yaml:"subsystem"`
	MetricName  string            `json:"metric_name" yaml:"metric_name"`
	Help        string            `json:"help" yaml:"help"`
	Labels      map[string]string `json:"labels" yaml:"labels"`
	Registry    string            `json:"registry" yaml:"registry"`
	PushGateway string            `json:"push_gateway" yaml:"push_gateway"`
}

// WebhookHookConfig Webhook钩子配置
type WebhookHookConfig struct {
	URL         string            `json:"url" yaml:"url"`
	Method      string            `json:"method" yaml:"method"`
	Headers     map[string]string `json:"headers" yaml:"headers"`
	Timeout     time.Duration     `json:"timeout" yaml:"timeout"`
	RetryCount  int               `json:"retry_count" yaml:"retry_count"`
	ContentType string            `json:"content_type" yaml:"content_type"`
	Template    string            `json:"template" yaml:"template"`
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	Name    string                 `json:"name" yaml:"name"`
	Type    string                 `json:"type" yaml:"type"`
	Enabled bool                   `json:"enabled" yaml:"enabled"`
	Order   int                    `json:"order" yaml:"order"`
	Config  map[string]interface{} `json:"config" yaml:"config"`
}

// EnvironmentConfig 环境配置
type EnvironmentConfig struct {
	Development *LogConfig `json:"development" yaml:"development"`
	Testing     *LogConfig `json:"testing" yaml:"testing"`
	Production  *LogConfig `json:"production" yaml:"production"`
}

// NewEnvironmentConfig 创建环境配置
func NewEnvironmentConfig() *EnvironmentConfig {
	return &EnvironmentConfig{
		Development: NewDevelopmentConfig(),
		Testing:     NewTestingConfig(),
		Production:  NewProductionConfig(),
	}
}

// NewDevelopmentConfig 创建开发环境配置
func NewDevelopmentConfig() *LogConfig {
	config := DefaultConfig()
	config.Level = DEBUG
	config.ShowCaller = true
	config.Colorful = true
	return config
}

// NewTestingConfig 创建测试环境配置
func NewTestingConfig() *LogConfig {
	config := DefaultConfig()
	config.Level = INFO
	config.ShowCaller = false
	config.Colorful = false
	return config
}

// NewProductionConfig 创建生产环境配置
func NewProductionConfig() *LogConfig {
	config := DefaultConfig()
	config.Level = WARN
	config.ShowCaller = false
	config.Colorful = false
	return config
}

// LogConfig 日志配置（兼容旧版本）
type LogConfig struct {
	Level      LogLevel  `json:"level"`       // 日志级别
	ShowCaller bool      `json:"show_caller"` // 是否显示调用者信息
	Prefix     string    `json:"prefix"`      // 日志前缀
	Output     io.Writer `json:"-"`           // 输出目标
	Colorful   bool      `json:"colorful"`    // 是否使用彩色输出
	TimeFormat string    `json:"time_format"` // 时间格式
}

// WithLevel 设置日志级别
func (c *Config) WithLevel(level LogLevel) *Config {
	c.Level = level
	return c
}

// WithOutput 设置输出目标
func (c *Config) WithOutput(output io.Writer) *Config {
	c.Output = output
	return c
}

// WithTimeFormat 设置时间格式
func (c *Config) WithTimeFormat(format TimeFormat) *Config {
	c.TimeFormat = format
	return c
}

// WithColorful 设置是否使用彩色输出
func (c *Config) WithColorful(colorful bool) *Config {
	c.Colorful = colorful
	return c
}

// WithAsyncWrite 设置是否异步写入
func (c *Config) WithAsyncWrite(async bool) *Config {
	c.AsyncWrite = async
	return c
}

// WithBufferSize 设置缓冲区大小
func (c *Config) WithBufferSize(size int) *Config {
	c.BufferSize = size
	return c
}

// WithPoolSize 设置对象池大小
func (c *Config) WithPoolSize(size int) *Config {
	c.PoolSize = size
	return c
}

// WithShowCaller 设置是否显示调用者信息
func (c *Config) WithShowCaller(show bool) *Config {
	c.ShowCaller = show
	return c
}

// WithCallerDepth 设置调用者深度
func (c *Config) WithCallerDepth(depth int) *Config {
	c.CallerDepth = depth
	return c
}

// WithFields 设置额外字段
func (c *Config) WithFields(fields map[string]interface{}) *Config {
	c.Fields = fields
	return c
}

// WithField 添加单个字段
func (c *Config) WithField(key string, value interface{}) *Config {
	if c.Fields == nil {
		c.Fields = make(map[string]interface{})
	}
	c.Fields[key] = value
	return c
}

// EnableMonitoring 启用监控
func (c *Config) EnableMonitoring() *Config {
	c.Monitoring.Enabled = true
	c.EnableMetrics = true
	return c
}

// EnableMemoryMonitoring 启用内存监控
func (c *Config) EnableMemoryMonitoring() *Config {
	c.Monitoring.Memory.Enabled = true
	c.EnableMemoryStats = true
	return c
}

// EnablePerformanceMonitoring 启用性能监控
func (c *Config) EnablePerformanceMonitoring() *Config {
	c.Monitoring.Performance.Enabled = true
	c.Monitoring.Performance.TrackLatency = true
	c.Monitoring.Performance.TrackThroughput = true
	return c
}

// Clone 克隆配置
func (c *Config) Clone() *Config {
	clone := *c

	// 深拷贝字段映射
	if c.Fields != nil {
		clone.Fields = make(map[string]interface{}, len(c.Fields))
		for k, v := range c.Fields {
			clone.Fields[k] = v
		}
	}

	// 深拷贝上下文字段切片
	if c.ContextFields != nil {
		clone.ContextFields = make([]string, len(c.ContextFields))
		copy(clone.ContextFields, c.ContextFields)
	}

	// 深拷贝适配器配置
	if c.Adapters != nil {
		clone.Adapters = make([]AdapterConfig, len(c.Adapters))
		copy(clone.Adapters, c.Adapters)
	}

	return &clone
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Output == nil {
		c.Output = os.Stdout
	}
	if c.TimeFormat == "" {
		c.TimeFormat = TimeFormatStandard
	}
	if c.BufferSize <= 0 {
		c.BufferSize = 4096
	}
	if c.PoolSize <= 0 {
		c.PoolSize = 10
	}
	if c.BatchSize <= 0 {
		c.BatchSize = 100
	}
	if c.BatchTimeout <= 0 {
		c.BatchTimeout = time.Millisecond * 100
	}
	if c.CallerDepth <= 0 {
		c.CallerDepth = 2
	}
	if c.Fields == nil {
		c.Fields = make(map[string]interface{})
	}

	return nil
}

// DefaultConfig 默认配置
func DefaultConfig() *LogConfig {
	return &LogConfig{
		Level:      INFO,
		ShowCaller: false,
		Prefix:     "",
		Output:     os.Stdout,
		Colorful:   true,
		TimeFormat: "2006-01-02 15:04:05",
	}
}

// NewLogConfig 创建新的日志配置（兼容旧版本）
func NewLogConfig() *LogConfig {
	return DefaultConfig()
}

// WithLevel 设置日志级别（兼容旧版本）
func (c *LogConfig) WithLevel(level LogLevel) *LogConfig {
	c.Level = level
	return c
}

// WithShowCaller 设置是否显示调用者信息（兼容旧版本）
func (c *LogConfig) WithShowCaller(show bool) *LogConfig {
	c.ShowCaller = show
	return c
}

// WithPrefix 设置日志前缀（兼容旧版本）
func (c *LogConfig) WithPrefix(prefix string) *LogConfig {
	if prefix != "" && !strings.HasSuffix(prefix, " ") {
		prefix += " "
	}
	c.Prefix = prefix
	return c
}

// WithOutput 设置输出目标（兼容旧版本）
func (c *LogConfig) WithOutput(output io.Writer) *LogConfig {
	c.Output = output
	return c
}

// WithColorful 设置是否使用彩色输出（兼容旧版本）
func (c *LogConfig) WithColorful(colorful bool) *LogConfig {
	c.Colorful = colorful
	return c
}

// WithTimeFormat 设置时间格式（兼容旧版本）
func (c *LogConfig) WithTimeFormat(format string) *LogConfig {
	c.TimeFormat = format
	return c
}

// Clone 克隆配置（兼容旧版本）
func (c *LogConfig) Clone() *LogConfig {
	return &LogConfig{
		Level:      c.Level,
		ShowCaller: c.ShowCaller,
		Prefix:     c.Prefix,
		Output:     c.Output,
		Colorful:   c.Colorful,
		TimeFormat: c.TimeFormat,
	}
}

// Validate 验证配置（兼容旧版本）
func (c *LogConfig) Validate() error {
	if c.Output == nil {
		c.Output = os.Stdout
	}
	if c.TimeFormat == "" {
		c.TimeFormat = "2006-01-02 15:04:05"
	}
	return nil
}

// ToConfig 转换为新的配置结构
func (c *LogConfig) ToConfig() *Config {
	return &Config{
		Level:      c.Level,
		Output:     c.Output,
		TimeFormat: TimeFormat(c.TimeFormat),
		Colorful:   c.Colorful,
		ShowCaller: c.ShowCaller,
		Fields: map[string]interface{}{
			"prefix": c.Prefix,
		},
	}
}
