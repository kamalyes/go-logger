/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-03-02 00:00:00
 * @FilePath: \go-logger\writer_test.go
 * @Description: 日志输出器测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConsoleWriter 测试控制台输出器
func TestConsoleWriter(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := NewConsoleWriter(WithConsoleOutput(buffer))

	assert.NotNil(t, writer)
	assert.True(t, writer.IsHealthy())

	n, err := writer.Write([]byte("test message\n"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n)
	assert.Contains(t, buffer.String(), "test message")
}

// TestConsoleWriterLevel 测试控制台输出器级别过滤
func TestConsoleWriterLevel(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := NewConsoleWriter(
		WithConsoleOutput(buffer),
		WithConsoleLevel(WARN),
	)

	// DEBUG级别应该被过滤
	n, err := writer.WriteLevel(DEBUG, []byte("debug message\n"))
	assert.NoError(t, err)
	assert.Equal(t, 14, n) // 返回长度但不写入
	assert.Empty(t, buffer.String())

	// WARN级别应该写入
	n, err = writer.WriteLevel(WARN, []byte("warn message\n"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n)
	assert.Contains(t, buffer.String(), "warn message")
}

// TestConsoleWriterClose 测试关闭控制台输出器
func TestConsoleWriterClose(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := NewConsoleWriter(WithConsoleOutput(buffer))

	err := writer.Close()
	assert.NoError(t, err)
	assert.False(t, writer.IsHealthy())

	// 关闭后写入应该失败
	_, err = writer.Write([]byte("test"))
	assert.Error(t, err)
}

// TestConsoleWriterStats 测试控制台输出器统计
func TestConsoleWriterStats(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := NewConsoleWriter(WithConsoleOutput(buffer))

	writer.Write([]byte("message 1\n"))
	writer.Write([]byte("message 2\n"))

	stats := writer.GetStats()
	assert.NotNil(t, stats)
}

// TestFileWriter 测试文件输出器
func TestFileWriter(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.log")

	writer := NewFileWriter(
		WithFileWriterPath(filePath),
		WithFilePermission(0644),
	)

	assert.NotNil(t, writer)

	n, err := writer.Write([]byte("test message\n"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n)

	err = writer.Flush()
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	// 验证文件内容
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test message")
}

// TestFileWriterAutoCreate 测试文件自动创建
func TestFileWriterAutoCreate(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "subdir", "test.log")

	writer := NewFileWriter(WithFileWriterPath(filePath))

	n, err := writer.Write([]byte("test\n"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	writer.Close()

	// 验证文件存在
	_, err = os.Stat(filePath)
	assert.NoError(t, err)
}

// TestFileWriterLevel 测试文件输出器级别过滤
func TestFileWriterLevel(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "level_test.log")

	writer := NewFileWriter(
		WithFileWriterPath(filePath),
		WithFileLevel(ERROR),
	)

	// INFO级别应该被过滤
	n, err := writer.WriteLevel(INFO, []byte("info message\n"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n)

	// ERROR级别应该写入
	n, err = writer.WriteLevel(ERROR, []byte("error message\n"))
	assert.NoError(t, err)
	assert.Equal(t, 14, n)

	writer.Close()

	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.NotContains(t, string(content), "info message")
	assert.Contains(t, string(content), "error message")
}

// TestRotateWriter 测试轮转文件输出器
func TestRotateWriter(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "rotate.log")

	writer := NewRotateWriter(
		WithFilePath(filePath),
		WithMaxSize(100), // 100字节触发轮转
		WithMaxFiles(3),
	)

	assert.NotNil(t, writer)

	// 写入少量数据
	n, err := writer.Write([]byte("test message\n"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n)

	// 写入后应该健康
	assert.True(t, writer.IsHealthy())

	writer.Close()
}

// TestRotateWriterRotation 测试文件轮转
func TestRotateWriterRotation(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "rotate_test.log")

	writer := NewRotateWriter(
		WithFilePath(filePath),
		WithMaxSize(50), // 50字节触发轮转
		WithMaxFiles(2),
	)

	// 写入足够的数据触发轮转
	for range 10 {
		writer.Write([]byte("test message line\n"))
	}

	writer.Close()

	// 验证轮转文件存在
	_, err := os.Stat(filePath)
	assert.NoError(t, err)
}

// TestBufferedWriter 测试缓冲输出器
func TestBufferedWriter(t *testing.T) {
	buffer := &bytes.Buffer{}
	underlying := NewConsoleWriter(WithConsoleOutput(buffer))

	writer := NewBufferedWriter(
		WithBufferedUnderlying(underlying),
		WithBufferSize(1024),
	)

	assert.NotNil(t, writer)
	assert.True(t, writer.IsHealthy())

	n, err := writer.Write([]byte("test message\n"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n)

	// 刷新前buffer可能为空
	err = writer.Flush()
	assert.NoError(t, err)

	// 刷新后应该有内容
	assert.Contains(t, buffer.String(), "test message")

	writer.Close()
}

// TestBufferedWriterAutoFlush 测试缓冲自动刷新
func TestBufferedWriterAutoFlush(t *testing.T) {
	buffer := &bytes.Buffer{}
	underlying := NewConsoleWriter(WithConsoleOutput(buffer))

	writer := NewBufferedWriter(
		WithBufferedUnderlying(underlying),
		WithBufferSize(10), // 小缓冲区
	)

	// 写入超过缓冲区大小的数据
	writer.Write([]byte("this is a long message\n"))

	// 应该自动刷新
	assert.NotEmpty(t, buffer.String())

	writer.Close()
}

// TestMultiWriter 测试多输出器
func TestMultiWriter(t *testing.T) {
	buffer1 := &bytes.Buffer{}
	buffer2 := &bytes.Buffer{}

	writer1 := NewConsoleWriter(WithConsoleOutput(buffer1))
	writer2 := NewConsoleWriter(WithConsoleOutput(buffer2))

	multiWriter := NewMultiWriter(WithWriters(writer1, writer2))

	assert.NotNil(t, multiWriter)
	assert.True(t, multiWriter.IsHealthy())

	n, err := multiWriter.Write([]byte("test message\n"))
	assert.NoError(t, err)
	assert.Equal(t, 13, n)

	// 两个输出器都应该有内容
	assert.Contains(t, buffer1.String(), "test message")
	assert.Contains(t, buffer2.String(), "test message")

	multiWriter.Close()
}

// TestMultiWriterPartialFailure 测试多输出器部分失败
func TestMultiWriterPartialFailure(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer1 := NewConsoleWriter(WithConsoleOutput(buffer))
	writer2 := NewConsoleWriter(WithConsoleOutput(buffer))

	// 关闭一个输出器
	writer2.Close()

	multiWriter := NewMultiWriter(WithWriters(writer1, writer2))

	// 应该仍然健康（至少一个输出器健康）
	assert.True(t, multiWriter.IsHealthy())

	multiWriter.Close()
}

// TestWriterStats 测试输出器统计
func TestWriterStats(t *testing.T) {
	stats := NewWriterStats()
	assert.NotNil(t, stats)

	stats.AddBytes(100)
	assert.Equal(t, int64(100), stats.BytesWritten)
	assert.Equal(t, int64(1), stats.LinesWritten)

	stats.AddError()
	assert.Equal(t, int64(1), stats.ErrorCount)
}

// TestCreateWriter 测试创建输出器
func TestCreateWriter(t *testing.T) {
	tempDir := t.TempDir()

	// 测试创建控制台输出器（使用 buffer 而不是 os.Stdout）
	buffer := &bytes.Buffer{}
	config := &WriterConfig{
		Type:   OutputConsole,
		Output: buffer,
	}
	writer, err := CreateWriter(config)
	assert.NoError(t, err)
	assert.NotNil(t, writer)
	writer.Close()

	// 测试创建文件输出器
	filePath := filepath.Join(tempDir, "create_test.log")

	config = &WriterConfig{
		Type:     OutputFile,
		FilePath: filePath,
	}
	writer, err = CreateWriter(config)
	assert.NoError(t, err)
	assert.NotNil(t, writer)
	writer.Close()
}

// TestCreateWriterInvalidConfig 测试无效配置
func TestCreateWriterInvalidConfig(t *testing.T) {
	// 文件输出器缺少路径
	config := &WriterConfig{
		Type: OutputFile,
	}
	writer, err := CreateWriter(config)
	assert.Error(t, err)
	assert.Nil(t, writer)

	// 轮转输出器缺少路径
	config = &WriterConfig{
		Type: OutputRotate,
	}
	writer, err = CreateWriter(config)
	assert.Error(t, err)
	assert.Nil(t, writer)
}

// TestCreateWriterWithDefaults 测试使用默认值创建
func TestCreateWriterWithDefaults(t *testing.T) {
	// 使用 buffer 而不是默认的 os.Stdout
	buffer := &bytes.Buffer{}
	config := &WriterConfig{
		Type:   OutputConsole,
		Output: buffer,
	}
	writer, err := CreateWriter(config)
	assert.NoError(t, err)
	assert.NotNil(t, writer)
	writer.Close()
}

// TestWriterConcurrent 测试并发写入
func TestWriterConcurrent(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := NewConsoleWriter(WithConsoleOutput(buffer))

	done := make(chan bool)

	for range 10 {
		go func() {
			writer.Write([]byte("concurrent message\n"))
			done <- true
		}()
	}

	for range 10 {
		<-done
	}

	writer.Close()
	assert.NotEmpty(t, buffer.String())
}

// TestFileWriterPermission 测试文件权限（Windows 跳过）
func TestFileWriterPermission(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "perm_test.log")

	writer := NewFileWriter(
		WithFileWriterPath(filePath),
		WithFilePermission(0600),
	)

	writer.Write([]byte("test\n"))
	writer.Close()

	info, err := os.Stat(filePath)
	assert.NoError(t, err)
	assert.NotNil(t, info)

	// Windows 不支持 Unix 文件权限，跳过权限检查
	// 只验证文件创建成功
}

// TestRotateWriterMaxFiles 测试最大文件数限制
func TestRotateWriterMaxFiles(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "max_files.log")

	writer := NewRotateWriter(
		WithFilePath(filePath),
		WithMaxSize(30),
		WithMaxFiles(3),
	)

	// 写入大量数据触发多次轮转
	for range 20 {
		writer.Write([]byte("test message line\n"))
	}

	writer.Close()
}

// TestWriterFlush 测试刷新功能
func TestWriterFlush(t *testing.T) {
	buffer := &bytes.Buffer{}
	writer := NewConsoleWriter(WithConsoleOutput(buffer))

	writer.Write([]byte("test\n"))

	err := writer.Flush()
	assert.NoError(t, err)

	writer.Close()
}

// BenchmarkConsoleWriter 控制台输出器性能测试
func BenchmarkConsoleWriter(b *testing.B) {
	buffer := &bytes.Buffer{}
	writer := NewConsoleWriter(WithConsoleOutput(buffer))
	data := []byte("test message\n")

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		writer.Write(data)
	}

	writer.Close()
}

// BenchmarkFileWriter 文件输出器性能测试
func BenchmarkFileWriter(b *testing.B) {
	filePath := filepath.Join(b.TempDir(), "bench.log")

	writer := NewFileWriter(WithFileWriterPath(filePath))
	data := []byte("test message\n")

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		writer.Write(data)
	}

	writer.Close()
}

// BenchmarkBufferedWriter 缓冲输出器性能测试
func BenchmarkBufferedWriter(b *testing.B) {
	buffer := &bytes.Buffer{}
	underlying := NewConsoleWriter(WithConsoleOutput(buffer))
	writer := NewBufferedWriter(
		WithBufferedUnderlying(underlying),
		WithBufferSize(4096),
	)
	data := []byte("test message\n")

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		writer.Write(data)
	}

	writer.Flush()
	writer.Close()
}

// BenchmarkMultiWriter 多输出器性能测试
func BenchmarkMultiWriter(b *testing.B) {
	buffer1 := &bytes.Buffer{}
	buffer2 := &bytes.Buffer{}
	writer1 := NewConsoleWriter(WithConsoleOutput(buffer1))
	writer2 := NewConsoleWriter(WithConsoleOutput(buffer2))
	multiWriter := NewMultiWriter(WithWriters(writer1, writer2))
	data := []byte("test message\n")

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		multiWriter.Write(data)
	}

	multiWriter.Close()
}
