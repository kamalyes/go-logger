/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-24 10:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-24 11:00:00
 * @FilePath: \go-logger\context_extractors.go
 * @Description: 预定义的上下文提取器和辅助函数
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
)

// NoOpContextExtractor 空操作提取器，不提取任何上下文信息
func NoOpContextExtractor(ctx context.Context) string {
	return ""
}

// SimpleTraceIDExtractor 只提取 TraceID
func SimpleTraceIDExtractor(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// 从 context.Value 获取
	if tid, ok := ctx.Value("trace_id").(string); ok && tid != "" {
		return "[TraceID=" + tid + "] "
	}

	// 从 gRPC metadata 获取
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get("x-trace-id"); len(values) > 0 && values[0] != "" {
			return "[TraceID=" + values[0] + "] "
		}
	}

	return ""
}

// SimpleRequestIDExtractor 只提取 RequestID
func SimpleRequestIDExtractor(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// 从 context.Value 获取
	if rid, ok := ctx.Value("request_id").(string); ok && rid != "" {
		return "[RequestID=" + rid + "] "
	}

	// 从 gRPC metadata 获取
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get("x-request-id"); len(values) > 0 && values[0] != "" {
			return "[RequestID=" + values[0] + "] "
		}
	}

	return ""
}

// CustomFieldExtractor 自定义字段提取器生成器
// 用于从 context 中提取指定的字段
func CustomFieldExtractor(contextKeys []string, metadataKeys []string) ContextExtractor {
	return func(ctx context.Context) string {
		if ctx == nil {
			return ""
		}

		fields := make(map[string]string)

		// 从 context.Value 提取
		for _, key := range contextKeys {
			if val, ok := ctx.Value(key).(string); ok && val != "" {
				fields[key] = val
			}
		}

		// 从 gRPC metadata 提取
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for _, key := range metadataKeys {
				if values := md.Get(key); len(values) > 0 && values[0] != "" {
					fields[key] = values[0]
				}
			}
		}

		if len(fields) == 0 {
			return ""
		}

		// 构建前缀
		prefix := "["
		first := true
		for k, v := range fields {
			if !first {
				prefix += " "
			}
			prefix += fmt.Sprintf("%s=%s", k, v)
			first = false
		}
		prefix += "] "

		return prefix
	}
}

// ChainContextExtractors 链接多个提取器
// 按顺序调用所有提取器，并合并结果
func ChainContextExtractors(extractors ...ContextExtractor) ContextExtractor {
	return func(ctx context.Context) string {
		var result string
		for _, extractor := range extractors {
			if extractor != nil {
				if info := extractor(ctx); info != "" {
					result += info
				}
			}
		}
		return result
	}
}

// ConditionalContextExtractor 条件提取器
// 只有当条件函数返回 true 时才调用提取器
func ConditionalContextExtractor(condition func(context.Context) bool, extractor ContextExtractor) ContextExtractor {
	return func(ctx context.Context) string {
		if condition != nil && condition(ctx) && extractor != nil {
			return extractor(ctx)
		}
		return ""
	}
}

// PrefixedContextExtractor 带前缀的提取器
// 为提取的信息添加自定义前缀
func PrefixedContextExtractor(prefix string, extractor ContextExtractor) ContextExtractor {
	return func(ctx context.Context) string {
		if extractor == nil {
			return ""
		}
		if info := extractor(ctx); info != "" {
			return prefix + info
		}
		return ""
	}
}

// CachedContextExtractor 缓存提取器
// 在同一个 context 中只提取一次，结果缓存在 context 中
func CachedContextExtractor(cacheKey string, extractor ContextExtractor) ContextExtractor {
	return func(ctx context.Context) string {
		if ctx == nil {
			return ""
		}

		// 尝试从缓存获取
		if cached, ok := ctx.Value(cacheKey).(string); ok {
			return cached
		}

		// 提取并缓存（注意：这里无法修改原 context，只能在当前调用返回）
		if extractor != nil {
			return extractor(ctx)
		}
		return ""
	}
}

// ExtractFromContextValue 从 context.Value 提取指定键的值
func ExtractFromContextValue(key string, label string) ContextExtractor {
	return func(ctx context.Context) string {
		if ctx == nil {
			return ""
		}
		if val, ok := ctx.Value(key).(string); ok && val != "" {
			return fmt.Sprintf("[%s=%s] ", label, val)
		}
		return ""
	}
}

// ExtractFromGRPCMetadata 从 gRPC metadata 提取指定键的值
func ExtractFromGRPCMetadata(key string, label string) ContextExtractor {
	return func(ctx context.Context) string {
		if ctx == nil {
			return ""
		}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get(key); len(values) > 0 && values[0] != "" {
				return fmt.Sprintf("[%s=%s] ", label, values[0])
			}
		}
		return ""
	}
}

// CreateContextExtractor 创建自定义上下文提取器的构建器
type ContextExtractorBuilder struct {
	extractors []ContextExtractor
}

// NewContextExtractorBuilder 创建新的提取器构建器
func NewContextExtractorBuilder() *ContextExtractorBuilder {
	return &ContextExtractorBuilder{
		extractors: make([]ContextExtractor, 0),
	}
}

// AddExtractor 添加提取器
func (b *ContextExtractorBuilder) AddExtractor(extractor ContextExtractor) *ContextExtractorBuilder {
	if extractor != nil {
		b.extractors = append(b.extractors, extractor)
	}
	return b
}

// AddTraceID 添加 TraceID 提取器
func (b *ContextExtractorBuilder) AddTraceID() *ContextExtractorBuilder {
	return b.AddExtractor(SimpleTraceIDExtractor)
}

// AddRequestID 添加 RequestID 提取器
func (b *ContextExtractorBuilder) AddRequestID() *ContextExtractorBuilder {
	return b.AddExtractor(SimpleRequestIDExtractor)
}

// AddContextValue 添加从 context.Value 提取的字段
func (b *ContextExtractorBuilder) AddContextValue(key, label string) *ContextExtractorBuilder {
	return b.AddExtractor(ExtractFromContextValue(key, label))
}

// AddGRPCMetadata 添加从 gRPC metadata 提取的字段
func (b *ContextExtractorBuilder) AddGRPCMetadata(key, label string) *ContextExtractorBuilder {
	return b.AddExtractor(ExtractFromGRPCMetadata(key, label))
}

// Build 构建最终的提取器
func (b *ContextExtractorBuilder) Build() ContextExtractor {
	if len(b.extractors) == 0 {
		return NoOpContextExtractor
	}
	if len(b.extractors) == 1 {
		return b.extractors[0]
	}
	return ChainContextExtractors(b.extractors...)
}
