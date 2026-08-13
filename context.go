/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-08 02:07:30
 * @FilePath: \go-logger\context.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// ContextKeyTraceID 是 traceId 在日志输出的字段名，同时作为 ctx.Value 的 fallback key
// traceId 的唯一真相源是 OTel span（gRPC 的 otelgrpc StatsHandler / HTTP 的 Tracing 中间件创建）
// extract 优先从 OTel span 提取；当 ctx 不在 span 内时，fallback 到 ctx.Value("trace_id")
const ContextKeyTraceID = "trace_id"

type compiledContextKey struct {
	key      string
	keyBytes []byte
}

var defaultContextKeys = []string{
	ContextKeyTraceID,
}

var defaultCompiledContextKeys = compileContextKeys(defaultContextKeys)

// extractOTelTraceID 从 OTel span 提取 traceId
// OTel span 是 traceId 的单一真相源：gRPC（otelgrpc StatsHandler）和 HTTP（Tracing 中间件）均由 OTel 创建 span
// 当 ctx 不在 span 内时返回空字符串，调用方 fallback 到 ctx.Value / metadata
func extractOTelTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	sc := trace.SpanContextFromContext(ctx)
	if !sc.HasTraceID() {
		return ""
	}
	return sc.TraceID().String()
}

// DefaultContextKeys 返回默认上下文提取 key
func DefaultContextKeys() []compiledContextKey {
	return compileContextKeys(defaultContextKeys)
}

func compileContextKeys(keys []string) []compiledContextKey {
	if len(keys) == 0 {
		return nil
	}

	compiled := make([]compiledContextKey, 0, len(keys))
	for _, key := range keys {
		compiled = append(compiled, compiledContextKey{
			key:      key,
			keyBytes: []byte(key),
		})
	}

	return compiled
}

// mdCache 延迟加载 gRPC incoming metadata，避免多次重复解析
type mdCache struct {
	md     metadata.MD
	loaded bool
	has    bool
}

// get 从 metadata 获取指定 key 的首个非空值（首次调用时延迟加载）
func (c *mdCache) get(ctx context.Context, key string) string {
	if !c.loaded {
		c.md, c.has = metadata.FromIncomingContext(ctx)
		c.loaded = true
	}
	if !c.has {
		return ""
	}
	if values := c.md.Get(key); len(values) > 0 && values[0] != "" {
		return values[0]
	}
	return ""
}

// extractKeyValue 从 context 提取单个 key 的值
// 优先级：OTel traceId（仅 trace_id key）> ctx.Value > gRPC metadata
func extractKeyValue(ctx context.Context, key compiledContextKey, otelTraceID string, md *mdCache) string {
	if key.key == ContextKeyTraceID && otelTraceID != "" {
		return otelTraceID
	}
	if raw := ctx.Value(key.key); raw != nil {
		if text, ok := raw.(string); ok && text != "" {
			return text
		}
	}
	return md.get(ctx, key.key)
}

func extractContextWithCompiledKeys(ctx context.Context, keys []compiledContextKey) string {
	if ctx == nil {
		return ""
	}

	buf := contextPool.Get().([]byte)
	buf = buf[:0]
	defer contextPool.Put(buf)

	buf = append(buf, '[')

	var (
		md         mdCache
		wroteField bool
	)

	// traceId 始终从 OTel span 提取（单一真相源），独立于 keys 配置
	// 即使 keys 为空或不含 trace_id，也输出 traceId，保证全链路日志打通
	otelTraceID := extractOTelTraceID(ctx)
	traceIDWritten := false
	if otelTraceID != "" {
		buf = append(buf, ContextKeyTraceID...)
		buf = append(buf, '=')
		buf = append(buf, otelTraceID...)
		wroteField = true
		traceIDWritten = true
	}

	for _, key := range keys {
		// trace_id 已由 OTel 写入，跳过
		if key.key == ContextKeyTraceID && traceIDWritten {
			continue
		}
		value := extractKeyValue(ctx, key, otelTraceID, &md)
		if value == "" {
			continue
		}
		if wroteField {
			buf = append(buf, ' ')
		}
		buf = append(buf, key.keyBytes...)
		buf = append(buf, '=')
		buf = append(buf, value...)
		wroteField = true
	}

	if !wroteField {
		return ""
	}

	buf = append(buf, ']', ' ')
	return string(buf)
}

// extractContextFieldsWithCompiledKeys 从 context 中提取键值对，返回 map 形式
// 用于 JSON 格式输出，使 traceId 等信息成为 JSON 字段而非 prepend 文本
func extractContextFieldsWithCompiledKeys(ctx context.Context, keys []compiledContextKey) map[string]any {
	if ctx == nil {
		return nil
	}

	var md mdCache
	fields := make(map[string]any, len(keys)+1)

	// traceId 始终从 OTel span 提取（单一真相源），独立于 keys 配置
	// 即使 keys 为空或不含 trace_id，也输出 traceId，保证全链路日志打通
	otelTraceID := extractOTelTraceID(ctx)
	traceIDWritten := false
	if otelTraceID != "" {
		fields[ContextKeyTraceID] = otelTraceID
		traceIDWritten = true
	}

	for _, key := range keys {
		// trace_id 已由 OTel 写入，跳过
		if key.key == ContextKeyTraceID && traceIDWritten {
			continue
		}
		if value := extractKeyValue(ctx, key, otelTraceID, &md); value != "" {
			fields[key.key] = value
		}
	}

	if len(fields) == 0 {
		return nil
	}
	return fields
}

// WithContextKeys 配置 Logger 在记录 Context 日志时提取哪些 key
func (l *Logger) WithContextKeys(keys ...string) *Logger {
	l.contextKeys = compileContextKeys(keys)
	l.contextExtractor = nil
	return l
}

// ExtractTraceID 从 ctx 提取 trace_id
// 优先级：OTel span > ctx.Value(ContextKeyTraceID) > 空
// 直接复用 context.go 的 extractOTelTraceID，保证与日志输出使用同一真相源
func ExtractTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	// 1. OTel span（单一真相源）- 复用 context.go 内部函数
	if id := extractOTelTraceID(ctx); id != "" {
		return id
	}
	// 2. ctx.Value fallback（ContextKeyTraceID = "trace_id"）
	if v := ctx.Value(ContextKeyTraceID); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// ContextWithTraceID 创建携带 trace_id 的 context
// 如果 ctx 已有 trace_id 则不覆盖（尊重上游传入的值）
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}
	if ExtractTraceID(ctx) != "" {
		return ctx // 已有则不覆盖
	}
	return context.WithValue(ctx, ContextKeyTraceID, traceID)
}

// InjectTraceToOutgoing 将 trace_id 注入 gRPC outgoing metadata
// 用于 gRPC 客户端调用远程节点时传播 trace
// 确保 OTel 未启用时仍可传播
func InjectTraceToOutgoing(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.MD{}
	} else {
		md = md.Copy() // 避免修改共享的 metadata
	}
	md.Set(ContextKeyTraceID, traceID)
	return metadata.NewOutgoingContext(ctx, md)
}

// ExtractTraceFromIncoming 从 gRPC incoming metadata 提取 trace_id
// 优先级：OTel span > metadata > ctx.Value
// 用于 gRPC 服务端接收到远程节点请求时恢复 trace
func ExtractTraceFromIncoming(ctx context.Context) string {
	// 1. OTel span - 复用 context.go 内部函数
	if id := extractOTelTraceID(ctx); id != "" {
		return id
	}
	// 2. metadata fallback
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get(ContextKeyTraceID); len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
	}
	// 3. ctx.Value fallback
	if v := ctx.Value(ContextKeyTraceID); v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// RestoreTraceFromIncoming 从 gRPC incoming 恢复 trace_id 到 ctx
// 如果提取到 trace_id 但 ctx 中没有，则注入到 ctx.Value
// 返回增强后的 ctx，供下游 logger 自动输出 trace_id
func RestoreTraceFromIncoming(ctx context.Context) context.Context {
	traceID := ExtractTraceFromIncoming(ctx)
	if traceID == "" {
		return ctx
	}
	// 只在 ctx 中缺失时注入，不覆盖已有值
	return ContextWithTraceID(ctx, traceID)
}
