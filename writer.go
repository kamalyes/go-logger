/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 23:19:49
 * @FilePath: \go-logger\writer.go
 * @Description: 日志输出器实现
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// WriterType 输出器类型
type WriterType string

const (
	ConsoleWriter WriterType = "console"
	FileWriter    WriterType = "file"
	RotateWriter  WriterType = "rotate"
	BufferWriter  WriterType = "buffer"
	MultiWriter   WriterType = "multi"
	NetworkWriter WriterType = "network"
)

// BaseWriter 基础输出器
type BaseWriter struct {
	Level    LogLevel `json:"level"`
	Healthy  bool     `json:"healthy"`
	Stats    *WriterStats `json:"-"`
	mutex    sync.RWMutex
}

// WriterStats 输出器统计信息
type WriterStats struct {
	BytesWritten int64         `json:"bytes_written"`
	LinesWritten int64         `json:"lines_written"`
	ErrorCount   int64         `json:"error_count"`
	LastWrite    time.Time     `json:"last_write"`
	StartTime    time.Time     `json:"start_time"`
	Uptime       time.Duration `json:"uptime"`
}

// NewWriterStats 创建输出器统计信息
func NewWriterStats() *WriterStats {
	return &WriterStats{
		StartTime: time.Now(),
	}
}

// AddBytes 增加字节统计
func (ws *WriterStats) AddBytes(bytes int64) {
	ws.BytesWritten += bytes
	ws.LinesWritten++
	ws.LastWrite = time.Now()
	ws.Uptime = time.Since(ws.StartTime)
}

// AddError 增加错误统计
func (ws *WriterStats) AddError() {
	ws.ErrorCount++
}

// ConsoleLogWriter 控制台输出器
type ConsoleLogWriter struct {
	BaseWriter
	Output io.Writer `json:"-"`
	Color  bool      `json:"color"`
}

// NewConsoleWriter 创建控制台输出器
func NewConsoleWriter(output io.Writer) IWriter {
	if output == nil {
		output = os.Stdout
	}
	
	return &ConsoleLogWriter{
		BaseWriter: BaseWriter{
			Level:   DEBUG,
			Healthy: true,
			Stats:   NewWriterStats(),
		},
		Output: output,
		Color:  true,
	}
}

// Write 实现io.Writer接口
func (w *ConsoleLogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if !w.Healthy {
		return 0, fmt.Errorf("console writer is not healthy")
	}
	
	n, err = w.Output.Write(p)
	if err != nil {
		w.Stats.AddError()
		return n, err
	}
	
	w.Stats.AddBytes(int64(n))
	return n, nil
}

// WriteLevel 按级别写入
func (w *ConsoleLogWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.Level {
		return len(data), nil // 跳过低级别日志
	}
	return w.Write(data)
}

// Flush 刷新缓冲区
func (w *ConsoleLogWriter) Flush() error {
	if flusher, ok := w.Output.(interface{ Flush() error }); ok {
		return flusher.Flush()
	}
	return nil
}

// Close 关闭输出器
func (w *ConsoleLogWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.Healthy = false
	if closer, ok := w.Output.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// IsHealthy 检查健康状态
func (w *ConsoleLogWriter) IsHealthy() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.Healthy
}

// GetStats 获取统计信息
func (w *ConsoleLogWriter) GetStats() interface{} {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return *w.Stats
}

// FileLogWriter 文件输出器
type FileLogWriter struct {
	BaseWriter
	FilePath   string   `json:"file_path"`
	file       *os.File
	Permission os.FileMode `json:"permission"`
}

// NewFileWriter 创建文件输出器
func NewFileWriter(filePath string) IWriter {
	return &FileLogWriter{
		BaseWriter: BaseWriter{
			Level:   DEBUG,
			Healthy: false, // 需要先打开文件
			Stats:   NewWriterStats(),
		},
		FilePath:   filePath,
		Permission: 0644,
	}
}

// ensureFile 确保文件已打开
func (w *FileLogWriter) ensureFile() error {
	if w.file != nil {
		return nil
	}
	
	// 创建目录
	dir := filepath.Dir(w.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// 打开文件
	file, err := os.OpenFile(w.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, w.Permission)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	
	w.file = file
	w.Healthy = true
	return nil
}

// Write 实现io.Writer接口
func (w *FileLogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if err := w.ensureFile(); err != nil {
		w.Stats.AddError()
		return 0, err
	}
	
	n, err = w.file.Write(p)
	if err != nil {
		w.Stats.AddError()
		w.Healthy = false
		return n, err
	}
	
	w.Stats.AddBytes(int64(n))
	return n, nil
}

// WriteLevel 按级别写入
func (w *FileLogWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.Level {
		return len(data), nil
	}
	return w.Write(data)
}

// Flush 刷新文件缓冲区
func (w *FileLogWriter) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// Close 关闭文件
func (w *FileLogWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.Healthy = false
	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}

// IsHealthy 检查健康状态
func (w *FileLogWriter) IsHealthy() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.Healthy
}

// GetStats 获取统计信息
func (w *FileLogWriter) GetStats() interface{} {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return *w.Stats
}

// RotateLogWriter 轮转文件输出器
type RotateLogWriter struct {
	BaseWriter
	FilePath    string      `json:"file_path"`
	MaxSize     int64       `json:"max_size"`     // 字节
	MaxFiles    int         `json:"max_files"`
	MaxAge      time.Duration `json:"max_age"`
	Compress    bool        `json:"compress"`
	currentFile *os.File
	currentSize int64
}

// NewRotateWriter 创建轮转文件输出器
func NewRotateWriter(filePath string, maxSize int64, maxFiles int) IWriter {
	return &RotateLogWriter{
		BaseWriter: BaseWriter{
			Level:   DEBUG,
			Healthy: false,
			Stats:   NewWriterStats(),
		},
		FilePath: filePath,
		MaxSize:  maxSize,
		MaxFiles: maxFiles,
		MaxAge:   30 * 24 * time.Hour, // 默认30天
		Compress: false,
	}
}

// shouldRotate 检查是否需要轮转
func (w *RotateLogWriter) shouldRotate(dataSize int) bool {
	return w.currentSize+int64(dataSize) > w.MaxSize
}

// rotate 执行文件轮转
func (w *RotateLogWriter) rotate() error {
	// 关闭当前文件
	if w.currentFile != nil {
		w.currentFile.Close()
		w.currentFile = nil
	}
	
	// 重命名现有文件
	for i := w.MaxFiles - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", w.FilePath, i)
		newPath := fmt.Sprintf("%s.%d", w.FilePath, i+1)
		
		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}
	
	// 移动当前文件
	if _, err := os.Stat(w.FilePath); err == nil {
		os.Rename(w.FilePath, w.FilePath+".1")
	}
	
	// 重置大小
	w.currentSize = 0
	
	return w.ensureFile()
}

// ensureFile 确保文件已打开
func (w *RotateLogWriter) ensureFile() error {
	if w.currentFile != nil {
		return nil
	}
	
	dir := filepath.Dir(w.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	file, err := os.OpenFile(w.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	
	// 获取当前文件大小
	if stat, err := file.Stat(); err == nil {
		w.currentSize = stat.Size()
	}
	
	w.currentFile = file
	w.Healthy = true
	return nil
}

// Write 实现io.Writer接口
func (w *RotateLogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	// 检查是否需要轮转
	if w.shouldRotate(len(p)) {
		if err := w.rotate(); err != nil {
			w.Stats.AddError()
			return 0, err
		}
	}
	
	if err := w.ensureFile(); err != nil {
		w.Stats.AddError()
		return 0, err
	}
	
	n, err = w.currentFile.Write(p)
	if err != nil {
		w.Stats.AddError()
		w.Healthy = false
		return n, err
	}
	
	w.currentSize += int64(n)
	w.Stats.AddBytes(int64(n))
	return n, nil
}

// WriteLevel 按级别写入
func (w *RotateLogWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.Level {
		return len(data), nil
	}
	return w.Write(data)
}

// Flush 刷新缓冲区
func (w *RotateLogWriter) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if w.currentFile != nil {
		return w.currentFile.Sync()
	}
	return nil
}

// Close 关闭输出器
func (w *RotateLogWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.Healthy = false
	if w.currentFile != nil {
		err := w.currentFile.Close()
		w.currentFile = nil
		return err
	}
	return nil
}

// IsHealthy 检查健康状态
func (w *RotateLogWriter) IsHealthy() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.Healthy
}

// GetStats 获取统计信息
func (w *RotateLogWriter) GetStats() interface{} {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return *w.Stats
}

// BufferedWriter 缓冲输出器
type BufferedWriter struct {
	BaseWriter
	underlying IWriter
	buffer     *bufio.Writer
	bufferSize int
}

// NewBufferedWriter 创建缓冲输出器
func NewBufferedWriter(underlying IWriter, bufferSize int) IWriter {
	return &BufferedWriter{
		BaseWriter: BaseWriter{
			Level:   DEBUG,
			Healthy: true,
			Stats:   NewWriterStats(),
		},
		underlying: underlying,
		buffer:     bufio.NewWriterSize(underlying, bufferSize),
		bufferSize: bufferSize,
	}
}

// Write 实现io.Writer接口
func (w *BufferedWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if !w.Healthy {
		return 0, fmt.Errorf("buffered writer is not healthy")
	}
	
	n, err = w.buffer.Write(p)
	if err != nil {
		w.Stats.AddError()
		w.Healthy = false
		return n, err
	}
	
	w.Stats.AddBytes(int64(n))
	return n, nil
}

// WriteLevel 按级别写入
func (w *BufferedWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.Level {
		return len(data), nil
	}
	return w.Write(data)
}

// Flush 刷新缓冲区
func (w *BufferedWriter) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	if w.buffer != nil {
		return w.buffer.Flush()
	}
	return nil
}

// Close 关闭输出器
func (w *BufferedWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.Healthy = false
	if w.buffer != nil {
		w.buffer.Flush()
	}
	if w.underlying != nil {
		return w.underlying.Close()
	}
	return nil
}

// IsHealthy 检查健康状态
func (w *BufferedWriter) IsHealthy() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.Healthy && w.underlying.IsHealthy()
}

// GetStats 获取统计信息
func (w *BufferedWriter) GetStats() interface{} {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return *w.Stats
}

// MultiLogWriter 多输出器
type MultiLogWriter struct {
	BaseWriter
	writers []IWriter
}

// NewMultiWriter 创建多输出器
func NewMultiWriter(writers ...IWriter) IWriter {
	return &MultiLogWriter{
		BaseWriter: BaseWriter{
			Level:   DEBUG,
			Healthy: true,
			Stats:   NewWriterStats(),
		},
		writers: writers,
	}
}

// Write 实现io.Writer接口
func (w *MultiLogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	var lastErr error
	for _, writer := range w.writers {
		if !writer.IsHealthy() {
			continue
		}
		
		if _, werr := writer.Write(p); werr != nil {
			lastErr = werr
			w.Stats.AddError()
		}
	}
	
	if lastErr != nil {
		return 0, lastErr
	}
	
	w.Stats.AddBytes(int64(len(p)))
	return len(p), nil
}

// WriteLevel 按级别写入
func (w *MultiLogWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.Level {
		return len(data), nil
	}
	return w.Write(data)
}

// Flush 刷新所有输出器
func (w *MultiLogWriter) Flush() error {
	var lastErr error
	for _, writer := range w.writers {
		if err := writer.Flush(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Close 关闭所有输出器
func (w *MultiLogWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	w.Healthy = false
	var lastErr error
	for _, writer := range w.writers {
		if err := writer.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// IsHealthy 检查是否至少有一个输出器健康
func (w *MultiLogWriter) IsHealthy() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	
	if !w.Healthy {
		return false
	}
	
	for _, writer := range w.writers {
		if writer.IsHealthy() {
			return true
		}
	}
	return false
}

// GetStats 获取统计信息
func (w *MultiLogWriter) GetStats() interface{} {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return *w.Stats
}