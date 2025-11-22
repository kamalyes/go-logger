/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-22 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 09:53:31
 * @FilePath: \go-logger\context_service.go
 * @Description: Trace 上下文管理器测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"github.com/kamalyes/go-toolbox/pkg/idgen"
	"sync"
	"time"
)

// ContextKey 统一键类型
type ContextKey string

// 常用键定义
const (
	KeyTraceID       ContextKey = "trace_id"
	KeySpanID        ContextKey = "span_id"
	KeySessionID     ContextKey = "session_id"
	KeyUserID        ContextKey = "user_id"
	KeyRequestID     ContextKey = "request_id"
	KeyCorrelationID ContextKey = "correlation_id"
	KeyOperation     ContextKey = "operation"
	KeyTenantID      ContextKey = "tenant_id"
)

// Field 定义一个可提取字段
type Field struct {
	Key       ContextKey
	LogName   string
	Generator func() string
}

// CorrelationChain 简化后的相关链
type CorrelationChain struct {
	ID        string
	TraceID   string
	StartTime time.Time
	EndTime   time.Time
	Tags      map[string]string
	Metrics   map[string]interface{}
	mu        sync.RWMutex
}

func (c *CorrelationChain) SetTag(k, v string) { c.mu.Lock(); c.Tags[k] = v; c.mu.Unlock() }
func (c *CorrelationChain) SetMetric(k string, v interface{}) {
	c.mu.Lock()
	c.Metrics[k] = v
	c.mu.Unlock()
}
func (c *CorrelationChain) GetDuration() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.EndTime.IsZero() {
		return time.Since(c.StartTime)
	}
	return c.EndTime.Sub(c.StartTime)
}

// ContextService 单一核心服务
type ContextService struct {
	idGen  idgen.IDGenerator
	fields []Field
	mu     sync.RWMutex
}

// NewContextService 创建服务
func NewContextService(gen idgen.IDGenerator) *ContextService {
	if gen == nil {
		gen = idgen.NewDefaultIDGenerator()
	}
	cs := &ContextService{idGen: gen}
	cs.fields = []Field{
		{KeyTraceID, "trace_id", func() string { return cs.idGen.GenerateTraceID() }},
		{KeySpanID, "span_id", func() string { return cs.idGen.GenerateSpanID() }},
		{KeyRequestID, "request_id", func() string { return cs.idGen.GenerateRequestID() }},
		{KeySessionID, "session_id", func() string { return cs.idGen.GenerateCorrelationID() }},
		{KeyCorrelationID, "correlation_id", func() string { return cs.idGen.GenerateCorrelationID() }},
	}
	return cs
}

// 默认实例
var defaultContextService = NewContextService(nil)

// 简单辅助
func (cs *ContextService) withValue(ctx context.Context, key ContextKey, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}
func (cs *ContextService) get(ctx context.Context, key ContextKey) interface{} {
	if ctx == nil {
		return nil
	}
	return ctx.Value(key)
}
func (cs *ContextService) getString(ctx context.Context, key ContextKey) string {
	v := cs.get(ctx, key)
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// EnsureID 确保某个 ID 存在
func (cs *ContextService) EnsureID(ctx context.Context, key ContextKey) (context.Context, string) {
	if existing := cs.getString(ctx, key); existing != "" {
		return ctx, existing
	}
	for _, f := range cs.fields {
		if f.Key == key && f.Generator != nil {
			id := f.Generator()
			ctx = cs.withValue(ctx, key, id)
			return ctx, id
		}
	}
	return ctx, ""
}

// ExtractFields 提取已注册字段
func (cs *ContextService) ExtractFields(ctx context.Context) map[string]interface{} {
	result := make(map[string]interface{})
	for _, f := range cs.fields {
		if v := cs.getString(ctx, f.Key); v != "" {
			result[f.LogName] = v
		}
	}
	// user/session/tenant/correlation 可能未在 fields 静态数组中显式映射成日志名时补充
	extras := []struct {
		key  ContextKey
		name string
	}{
		{KeyUserID, "user_id"}, {KeyTenantID, "tenant_id"}, {KeyCorrelationID, "correlation_id"}, {KeySessionID, "session_id"},
	}
	for _, e := range extras {
		if _, exists := result[e.name]; !exists {
			if v := cs.getString(ctx, e.key); v != "" {
				result[e.name] = v
			}
		}
	}
	return result
}

// CreateSpan 创建 span (继承 trace, 新 spanID)
func (cs *ContextService) CreateSpan(ctx context.Context, operation string) context.Context {
	ctx, _ = cs.EnsureID(ctx, KeyTraceID)
	ctx, _ = cs.EnsureID(ctx, KeySpanID)
	if operation != "" {
		ctx = cs.withValue(ctx, KeyOperation, operation)
	}
	return ctx
}

// CreateCorrelationChain 创建相关性链
func (cs *ContextService) CreateCorrelationChain(ctx context.Context) (*CorrelationChain, context.Context) {
	ctx, traceID := cs.EnsureID(ctx, KeyTraceID)
	id := cs.idGen.GenerateCorrelationID()
	chain := &CorrelationChain{ID: id, TraceID: traceID, StartTime: time.Now(), Tags: map[string]string{}, Metrics: map[string]interface{}{}}
	ctx = cs.withValue(ctx, KeyCorrelationID, id)
	return chain, ctx
}

// EndChain 结束相关性链
func (cs *ContextService) EndChain(chain *CorrelationChain) {
	if chain != nil {
		chain.EndTime = time.Now()
	}
}

// ========== 全局函数（新的唯一入口） ==========

func WithValue(ctx context.Context, key ContextKey, val interface{}) context.Context {
	return defaultContextService.withValue(ctx, key, val)
}
func GetValue(ctx context.Context, key ContextKey) interface{} {
	return defaultContextService.get(ctx, key)
}
func GetString(ctx context.Context, key ContextKey) string {
	return defaultContextService.getString(ctx, key)
}
func ExtractFields(ctx context.Context) map[string]interface{} {
	return defaultContextService.ExtractFields(ctx)
}
func CreateSpan(ctx context.Context, operation string) context.Context {
	return defaultContextService.CreateSpan(ctx, operation)
}
func GenerateTraceID() string       { return defaultContextService.idGen.GenerateTraceID() }
func GenerateSpanID() string        { return defaultContextService.idGen.GenerateSpanID() }
func GenerateRequestID() string     { return defaultContextService.idGen.GenerateRequestID() }
func GenerateCorrelationID() string { return defaultContextService.idGen.GenerateCorrelationID() }

// ID 操作便捷函数
func WithTraceID(ctx context.Context, id string) context.Context {
	return WithValue(ctx, KeyTraceID, id)
}
func GetTraceID(ctx context.Context) string { return GetString(ctx, KeyTraceID) }
func GetOrGenerateTraceID(ctx context.Context) (context.Context, string) {
	return defaultContextService.EnsureID(ctx, KeyTraceID)
}
func WithSpanID(ctx context.Context, id string) context.Context { return WithValue(ctx, KeySpanID, id) }
func GetSpanID(ctx context.Context) string                      { return GetString(ctx, KeySpanID) }
func WithRequestID(ctx context.Context, id string) context.Context {
	return WithValue(ctx, KeyRequestID, id)
}
func GetRequestID(ctx context.Context) string                   { return GetString(ctx, KeyRequestID) }
func WithUserID(ctx context.Context, id string) context.Context { return WithValue(ctx, KeyUserID, id) }
func GetUserID(ctx context.Context) string                      { return GetString(ctx, KeyUserID) }
func WithSessionID(ctx context.Context, id string) context.Context {
	return WithValue(ctx, KeySessionID, id)
}
func GetSessionID(ctx context.Context) string { return GetString(ctx, KeySessionID) }
func WithTenantID(ctx context.Context, id string) context.Context {
	return WithValue(ctx, KeyTenantID, id)
}
func GetTenantID(ctx context.Context) string { return GetString(ctx, KeyTenantID) }
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return WithValue(ctx, KeyCorrelationID, id)
}
func GetCorrelationID(ctx context.Context) string { return GetString(ctx, KeyCorrelationID) }

// Correlation 操作
func CreateCorrelationChain(ctx context.Context) (*CorrelationChain, context.Context) {
	return defaultContextService.CreateCorrelationChain(ctx)
}
func EndCorrelationChain(chain *CorrelationChain) { defaultContextService.EndChain(chain) }
