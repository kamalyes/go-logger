/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-11-03 17:52:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-11-03 18:18:00
 * @FilePath: \\go-logger\\fields_typed_bench_test.go
 * @Description: 类型化字段与 ...any 路径的基准 - 单线程与并发两档都必须给
 *
 * 两档都给是因为结论在两档之间会反转：单线程下类型化更慢（结构体拷贝主导），
 * 并发下类型化更快（GC 压力主导）
 * 只报一档就是在选一个对自己有利的数字
 *
 * [EN] Both a single-threaded and a parallel variant, because the conclusion reverses
 * between them: typed loses single-threaded where struct copying dominates and wins in
 * parallel where GC pressure does
 * Reporting only one of them is choosing the flattering number
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package logger

import (
	"io"
	"testing"
)

var tfVals = func() []string {
	out := make([]string, 256)
	for i := range out {
		out[i] = "value-" + string(rune('a'+i%26)) + string(rune('0'+i%10))
	}
	return out
}()

func tfLogger() *Logger {
	return NewLogger().WithFormat(FormatJSON).WithShowCaller(false).WithOutput(io.Discard)
}

// 十四个字段，与 go-rpc-gateway 访问日志的真实规模一致，值每次都不同
func BenchmarkTypedFields_KV(b *testing.B) {
	l := tfLogger()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := tfVals[i%len(tfVals)]
		l.InfoKV("HTTP Request",
			"method", v, "path", v, "status", int64(200), "bytes", int64(290),
			"duration_ms", int64(1), "ip", v, "user_agent", v, "trace_id", v,
			"request_id", v, "user_id", v, "tenant_id", v, "domain", v,
			"platform_id", v, "region_id", v)
	}
}

func BenchmarkTypedFields_Typed(b *testing.B) {
	l := tfLogger()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v := tfVals[i%len(tfVals)]
		l.InfoFields("HTTP Request",
			Str("method", v), Str("path", v), Int("status", 200), Int("bytes", 290),
			Int("duration_ms", 1), Str("ip", v), Str("user_agent", v), Str("trace_id", v),
			Str("request_id", v), Str("user_id", v), Str("tenant_id", v), Str("domain", v),
			Str("platform_id", v), Str("region_id", v))
	}
}

// 并发档：单线程微基准看不出「去掉分配」的价值 ——
// 分配的代价是 GC 工作量，它体现在并发下的吞吐与尾延迟，不是单线程的 ns/op
func BenchmarkTypedFields_KV_Parallel(b *testing.B) {
	l := tfLogger()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			v := tfVals[i%len(tfVals)]
			i++
			l.InfoKV("HTTP Request",
				"method", v, "path", v, "status", int64(200), "bytes", int64(290),
				"duration_ms", int64(1), "ip", v, "user_agent", v, "trace_id", v,
				"request_id", v, "user_id", v, "tenant_id", v, "domain", v,
				"platform_id", v, "region_id", v)
		}
	})
}

func BenchmarkTypedFields_Typed_Parallel(b *testing.B) {
	l := tfLogger()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			v := tfVals[i%len(tfVals)]
			i++
			l.InfoFields("HTTP Request",
				Str("method", v), Str("path", v), Int("status", 200), Int("bytes", 290),
				Int("duration_ms", 1), Str("ip", v), Str("user_agent", v), Str("trace_id", v),
				Str("request_id", v), Str("user_id", v), Str("tenant_id", v), Str("domain", v),
				Str("platform_id", v), Str("region_id", v))
		}
	})
}
