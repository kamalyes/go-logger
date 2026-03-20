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

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"google.golang.org/grpc/metadata"
)

const (
	ContextKeyTraceID = "trace_id"
)

// Metadata keys - 用于从 gRPC metadata 获取
const (
	MetadataKeyTraceID = "x-trace-id"
)

type compiledContextKey struct {
	key      string
	keyBytes []byte
}

var defaultContextKeys = []string{
	ContextKeyTraceID,
	MetadataKeyTraceID,
}

var defaultCompiledContextKeys = compileContextKeys(defaultContextKeys)

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
			keyBytes: convert.S2B(key),
		})
	}

	return compiled
}

func extractContextWithCompiledKeys(ctx context.Context, keys []compiledContextKey) string {
	if ctx == nil || len(keys) == 0 {
		return ""
	}

	buf := contextPool.Get().([]byte)
	buf = buf[:0]
	defer contextPool.Put(buf)

	buf = append(buf, '[')

	var (
		incomingMD metadata.MD
		mdLoaded   bool
		hasMD      bool
		wroteField bool
	)

	for _, key := range keys {
		value := ""
		if raw := ctx.Value(key.key); raw != nil {
			if text, ok := raw.(string); ok && text != "" {
				value = text
			}
		}

		if value == "" {
			if !mdLoaded {
				incomingMD, hasMD = metadata.FromIncomingContext(ctx)
				mdLoaded = true
			}
			if hasMD {
				if values := incomingMD.Get(key.key); len(values) > 0 && values[0] != "" {
					value = values[0]
				}
			}
		}

		if value == "" {
			continue
		}

		if wroteField {
			buf = append(buf, ' ')
		}
		buf = append(buf, key.keyBytes...)
		buf = append(buf, '=')
		buf = append(buf, convert.S2B(value)...)
		wroteField = true
	}

	if !wroteField {
		return ""
	}

	buf = append(buf, ']', ' ')
	return string(buf)
}

// WithContextKeys 配置 Logger 在记录 Context 日志时提取哪些 key
func (l *Logger) WithContextKeys(keys ...string) *Logger {
	l.contextKeys = compileContextKeys(keys)
	l.contextExtractor = nil
	return l
}
