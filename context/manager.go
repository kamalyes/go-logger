/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\context\manager.go
 * @Description: 上下文管理器
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package context

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ContextKey 上下文键类型
type ContextKey string

// 预定义的上下文键
const (
	KeyTraceID     ContextKey = "trace_id"
	KeySpanID      ContextKey = "span_id"
	KeySessionID   ContextKey = "session_id"
	KeyUserID      ContextKey = "user_id"
	KeyRequestID   ContextKey = "request_id"
	KeyCorrelationID ContextKey = "correlation_id"
	KeyVersion     ContextKey = "version"
	KeyEnvironment ContextKey = "environment"
	KeyComponent   ContextKey = "component"
	KeyOperation   ContextKey = "operation"
	KeyStartTime   ContextKey = "start_time"
	KeyTags        ContextKey = "tags"
	KeyMetrics     ContextKey = "metrics"
)

// LogContext 日志上下文
type LogContext struct {
	// 基础信息
	TraceID       string            `json:"trace_id"`
	SpanID        string            `json:"span_id"`
	SessionID     string            `json:"session_id"`
	UserID        string            `json:"user_id"`
	RequestID     string            `json:"request_id"`
	CorrelationID string            `json:"correlation_id"`
	
	// 应用信息
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Component   string `json:"component"`
	Operation   string `json:"operation"`
	
	// 时间信息
	StartTime time.Time `json:"start_time"`
	Duration  time.Duration `json:"duration,omitempty"`
	
	// 标签和指标
	Tags    map[string]string      `json:"tags,omitempty"`
	Metrics map[string]interface{} `json:"metrics,omitempty"`
	
	// 自定义字段
	Fields map[string]interface{} `json:"fields,omitempty"`
	
	mu sync.RWMutex
}

// ContextManager 上下文管理器
type ContextManager struct {
	// 全局上下文
	globalContext *LogContext
	
	// 默认值
	defaultVersion     string
	defaultEnvironment string
	defaultComponent   string
	
	// 上下文池
	contextPool sync.Pool
	
	// 配置
	enableTracing   bool
	enableMetrics   bool
	maxContextAge   time.Duration
	cleanupInterval time.Duration
	
	// 上下文存储
	contexts map[string]*LogContext
	
	mu sync.RWMutex
}

// NewContextManager 创建上下文管理器
func NewContextManager() *ContextManager {
	cm := &ContextManager{
		globalContext:      NewLogContext(),
		defaultVersion:     "1.0.0",
		defaultEnvironment: "dev",
		defaultComponent:   "unknown",
		enableTracing:      true,
		enableMetrics:      true,
		maxContextAge:      time.Hour,
		cleanupInterval:    time.Minute * 10,
		contexts:           make(map[string]*LogContext),
	}
	
	// 初始化上下文池
	cm.contextPool = sync.Pool{
		New: func() interface{} {
			return NewLogContext()
		},
	}
	
	// 启动清理协程
	go cm.startCleanup()
	
	return cm
}

// NewLogContext 创建新的日志上下文
func NewLogContext() *LogContext {
	return &LogContext{
		TraceID:     generateTraceID(),
		SpanID:      generateSpanID(),
		RequestID:   generateRequestID(),
		StartTime:   time.Now(),
		Tags:        make(map[string]string),
		Metrics:     make(map[string]interface{}),
		Fields:      make(map[string]interface{}),
	}
}

// SetDefaults 设置默认值
func (cm *ContextManager) SetDefaults(version, environment, component string) *ContextManager {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.defaultVersion = version
	cm.defaultEnvironment = environment
	cm.defaultComponent = component
	
	// 更新全局上下文
	cm.globalContext.Version = version
	cm.globalContext.Environment = environment
	cm.globalContext.Component = component
	
	return cm
}

// SetGlobalContext 设置全局上下文
func (cm *ContextManager) SetGlobalContext(ctx *LogContext) *ContextManager {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.globalContext = ctx
	return cm
}

// GetGlobalContext 获取全局上下文
func (cm *ContextManager) GetGlobalContext() *LogContext {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	return cm.globalContext.Clone()
}

// GetContext 从池中获取上下文
func (cm *ContextManager) GetContext() *LogContext {
	ctx := cm.contextPool.Get().(*LogContext)
	
	// 重置上下文
	ctx.Reset()
	
	// 设置默认值
	ctx.Version = cm.defaultVersion
	ctx.Environment = cm.defaultEnvironment
	ctx.Component = cm.defaultComponent
	
	return ctx
}

// PutContext 返回上下文到池
func (cm *ContextManager) PutContext(ctx *LogContext) {
	if ctx != nil {
		cm.contextPool.Put(ctx)
	}
}

// SaveContext 保存上下文
func (cm *ContextManager) SaveContext(key string, ctx *LogContext) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.contexts[key] = ctx.Clone()
}

// LoadContext 加载上下文
func (cm *ContextManager) LoadContext(key string) (*LogContext, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	ctx, exists := cm.contexts[key]
	if !exists {
		return nil, false
	}
	
	return ctx.Clone(), true
}

// CreateChildContext 创建子上下文
func (cm *ContextManager) CreateChildContext(parent *LogContext, operation string) *LogContext {
	child := cm.GetContext()
	
	// 继承父上下文信息
	if parent != nil {
		child.TraceID = parent.TraceID
		child.SessionID = parent.SessionID
		child.UserID = parent.UserID
		child.RequestID = parent.RequestID
		child.CorrelationID = parent.CorrelationID
		child.Version = parent.Version
		child.Environment = parent.Environment
		child.Component = parent.Component
		
		// 复制标签
		for k, v := range parent.Tags {
			child.Tags[k] = v
		}
	}
	
	// 设置新的SpanID和操作
	child.SpanID = generateSpanID()
	child.Operation = operation
	child.StartTime = time.Now()
	
	return child
}

// FromStdContext 从标准上下文中提取日志上下文
func (cm *ContextManager) FromStdContext(ctx context.Context) *LogContext {
	logCtx := cm.GetContext()
	
	// 提取已知字段
	if val := ctx.Value(KeyTraceID); val != nil {
		if traceID, ok := val.(string); ok {
			logCtx.TraceID = traceID
		}
	}
	
	if val := ctx.Value(KeySpanID); val != nil {
		if spanID, ok := val.(string); ok {
			logCtx.SpanID = spanID
		}
	}
	
	if val := ctx.Value(KeySessionID); val != nil {
		if sessionID, ok := val.(string); ok {
			logCtx.SessionID = sessionID
		}
	}
	
	if val := ctx.Value(KeyUserID); val != nil {
		if userID, ok := val.(string); ok {
			logCtx.UserID = userID
		}
	}
	
	if val := ctx.Value(KeyRequestID); val != nil {
		if requestID, ok := val.(string); ok {
			logCtx.RequestID = requestID
		}
	}
	
	if val := ctx.Value(KeyCorrelationID); val != nil {
		if correlationID, ok := val.(string); ok {
			logCtx.CorrelationID = correlationID
		}
	}
	
	// 提取其他字段
	if val := ctx.Value(KeyComponent); val != nil {
		if component, ok := val.(string); ok {
			logCtx.Component = component
		}
	}
	
	if val := ctx.Value(KeyOperation); val != nil {
		if operation, ok := val.(string); ok {
			logCtx.Operation = operation
		}
	}
	
	return logCtx
}

// ToStdContext 将日志上下文添加到标准上下文
func (cm *ContextManager) ToStdContext(ctx context.Context, logCtx *LogContext) context.Context {
	if logCtx == nil {
		return ctx
	}
	
	// 添加基础字段
	if logCtx.TraceID != "" {
		ctx = context.WithValue(ctx, KeyTraceID, logCtx.TraceID)
	}
	if logCtx.SpanID != "" {
		ctx = context.WithValue(ctx, KeySpanID, logCtx.SpanID)
	}
	if logCtx.SessionID != "" {
		ctx = context.WithValue(ctx, KeySessionID, logCtx.SessionID)
	}
	if logCtx.UserID != "" {
		ctx = context.WithValue(ctx, KeyUserID, logCtx.UserID)
	}
	if logCtx.RequestID != "" {
		ctx = context.WithValue(ctx, KeyRequestID, logCtx.RequestID)
	}
	if logCtx.CorrelationID != "" {
		ctx = context.WithValue(ctx, KeyCorrelationID, logCtx.CorrelationID)
	}
	if logCtx.Component != "" {
		ctx = context.WithValue(ctx, KeyComponent, logCtx.Component)
	}
	if logCtx.Operation != "" {
		ctx = context.WithValue(ctx, KeyOperation, logCtx.Operation)
	}
	
	// 添加标签
	if len(logCtx.Tags) > 0 {
		ctx = context.WithValue(ctx, KeyTags, logCtx.Tags)
	}
	
	// 添加指标
	if len(logCtx.Metrics) > 0 {
		ctx = context.WithValue(ctx, KeyMetrics, logCtx.Metrics)
	}
	
	return ctx
}

// EnableTracing 启用追踪
func (cm *ContextManager) EnableTracing(enable bool) *ContextManager {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.enableTracing = enable
	return cm
}

// EnableMetrics 启用指标
func (cm *ContextManager) EnableMetrics(enable bool) *ContextManager {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.enableMetrics = enable
	return cm
}

// SetCleanupConfig 设置清理配置
func (cm *ContextManager) SetCleanupConfig(maxAge, interval time.Duration) *ContextManager {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.maxContextAge = maxAge
	cm.cleanupInterval = interval
	return cm
}

// startCleanup 启动清理协程
func (cm *ContextManager) startCleanup() {
	ticker := time.NewTicker(cm.cleanupInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		cm.cleanup()
	}
}

// cleanup 清理过期上下文
func (cm *ContextManager) cleanup() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	now := time.Now()
	for key, ctx := range cm.contexts {
		if now.Sub(ctx.StartTime) > cm.maxContextAge {
			delete(cm.contexts, key)
		}
	}
}

// Reset 重置上下文
func (c *LogContext) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.TraceID = ""
	c.SpanID = ""
	c.SessionID = ""
	c.UserID = ""
	c.RequestID = ""
	c.CorrelationID = ""
	c.Version = ""
	c.Environment = ""
	c.Component = ""
	c.Operation = ""
	c.StartTime = time.Time{}
	c.Duration = 0
	
	// 清空映射
	for k := range c.Tags {
		delete(c.Tags, k)
	}
	for k := range c.Metrics {
		delete(c.Metrics, k)
	}
	for k := range c.Fields {
		delete(c.Fields, k)
	}
}

// Clone 克隆上下文
func (c *LogContext) Clone() *LogContext {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// 克隆映射
	tags := make(map[string]string)
	for k, v := range c.Tags {
		tags[k] = v
	}
	
	metrics := make(map[string]interface{})
	for k, v := range c.Metrics {
		metrics[k] = v
	}
	
	fields := make(map[string]interface{})
	for k, v := range c.Fields {
		fields[k] = v
	}
	
	return &LogContext{
		TraceID:       c.TraceID,
		SpanID:        c.SpanID,
		SessionID:     c.SessionID,
		UserID:        c.UserID,
		RequestID:     c.RequestID,
		CorrelationID: c.CorrelationID,
		Version:       c.Version,
		Environment:   c.Environment,
		Component:     c.Component,
		Operation:     c.Operation,
		StartTime:     c.StartTime,
		Duration:      c.Duration,
		Tags:          tags,
		Metrics:       metrics,
		Fields:        fields,
	}
}

// SetField 设置字段
func (c *LogContext) SetField(key string, value interface{}) *LogContext {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.Fields[key] = value
	return c
}

// GetField 获取字段
func (c *LogContext) GetField(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	value, exists := c.Fields[key]
	return value, exists
}

// SetTag 设置标签
func (c *LogContext) SetTag(key, value string) *LogContext {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.Tags[key] = value
	return c
}

// GetTag 获取标签
func (c *LogContext) GetTag(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	value, exists := c.Tags[key]
	return value, exists
}

// SetMetric 设置指标
func (c *LogContext) SetMetric(key string, value interface{}) *LogContext {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.Metrics[key] = value
	return c
}

// GetMetric 获取指标
func (c *LogContext) GetMetric(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	value, exists := c.Metrics[key]
	return value, exists
}

// UpdateDuration 更新持续时间
func (c *LogContext) UpdateDuration() *LogContext {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if !c.StartTime.IsZero() {
		c.Duration = time.Since(c.StartTime)
	}
	return c
}

// String 字符串表示
func (c *LogContext) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return fmt.Sprintf("LogContext{TraceID: %s, SpanID: %s, RequestID: %s, Component: %s, Operation: %s}",
		c.TraceID, c.SpanID, c.RequestID, c.Component, c.Operation)
}