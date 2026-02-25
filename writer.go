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
	"sync/atomic"
	"time"
)

// baseWriter 基础输出器（所有输出器的通用字段）
type baseWriter struct {
	level      LogLevel      // 日志级别过滤
	healthy    bool          // 健康状态标识
	stats      *writerStats  // 统计信息
	permission os.FileMode   // 文件权限（适用于文件类输出器）
	maxAge     time.Duration // 最大保留时间（适用于轮转输出器）
	compress   bool          // 是否压缩旧文件（适用于轮转输出器）
	mutex      sync.RWMutex  // 读写锁保护并发访问
}

// writerStats 输出器统计信息（使用 atomic 优化并发性能）
type writerStats struct {
	bytesWritten int64        // 已写入字节数（atomic 计数器）
	linesWritten int64        // 已写入行数（atomic 计数器）
	errorCount   int64        // 错误次数（atomic 计数器）
	lastWrite    int64        // 最后写入时间（atomic unix nano）
	startTime    time.Time    // 启动时间（不可变）
	mu           sync.RWMutex // 保护同步操作的锁
}

// newWriterStats 创建输出器统计信息
func newWriterStats() *writerStats {
	return &writerStats{
		startTime: time.Now(),
	}
}

// addBytes 增加字节统计（使用 atomic 快速更新）
func (ws *writerStats) addBytes(bytes int64) {
	atomic.AddInt64(&ws.bytesWritten, bytes)
	atomic.AddInt64(&ws.linesWritten, 1)
	atomic.StoreInt64(&ws.lastWrite, time.Now().UnixNano())
}

// addError 增加错误统计（使用 atomic 快速更新）
func (ws *writerStats) addError() {
	atomic.AddInt64(&ws.errorCount, 1)
}

// WriterStatsSnapshot 统计信息快照（用于外部访问的只读数据）
type WriterStatsSnapshot struct {
	BytesWritten int64         `json:"bytes_written"` // 已写入字节总数
	LinesWritten int64         `json:"lines_written"` // 已写入行数总数
	ErrorCount   int64         `json:"error_count"`   // 错误次数总数
	LastWrite    time.Time     `json:"last_write"`    // 最后一次写入时间
	StartTime    time.Time     `json:"start_time"`    // 输出器启动时间
	Uptime       time.Duration `json:"uptime"`        // 运行时长
}

// getSnapshot 获取统计信息快照
func (ws *writerStats) getSnapshot() WriterStatsSnapshot {
	bytesWritten := atomic.LoadInt64(&ws.bytesWritten)
	linesWritten := atomic.LoadInt64(&ws.linesWritten)
	errorCount := atomic.LoadInt64(&ws.errorCount)
	lastWriteNano := atomic.LoadInt64(&ws.lastWrite)

	var lastWrite time.Time
	if lastWriteNano > 0 {
		lastWrite = time.Unix(0, lastWriteNano)
	}

	return WriterStatsSnapshot{
		BytesWritten: bytesWritten,
		LinesWritten: linesWritten,
		ErrorCount:   errorCount,
		LastWrite:    lastWrite,
		StartTime:    ws.startTime,
		Uptime:       time.Since(ws.startTime),
	}
}

// consoleLogWriter 控制台输出器（输出到标准输出或标准错误）
type consoleLogWriter struct {
	baseWriter              // 继承基础输出器字段
	output        io.Writer // 输出目标（如 os.Stdout, os.Stderr）
	color         bool      // 是否启用颜色输出
	healthyAtomic int32     // 健康状态（atomic bool: 0=false, 1=true）
}

// ConsoleWriterOption 控制台输出器配置选项
type ConsoleWriterOption func(*consoleLogWriter)

// WithConsoleOutput 设置输出目标
func WithConsoleOutput(output io.Writer) ConsoleWriterOption {
	return func(w *consoleLogWriter) {
		w.output = output
	}
}

// WithConsoleColor 设置是否启用颜色
func WithConsoleColor(color bool) ConsoleWriterOption {
	return func(w *consoleLogWriter) {
		w.color = color
	}
}

// WithConsoleLevel 设置日志级别
func WithConsoleLevel(level LogLevel) ConsoleWriterOption {
	return func(w *consoleLogWriter) {
		w.level = level
	}
}

// NewConsoleWriter 创建控制台输出器
func NewConsoleWriter(opts ...ConsoleWriterOption) IWriter {
	w := &consoleLogWriter{
		baseWriter: baseWriter{
			level:   DEBUG,
			healthy: true,
			stats:   newWriterStats(),
		},
		output:        os.Stdout,
		color:         true,
		healthyAtomic: 1,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// Write 实现io.Writer接口（优化：减少函数调用开销）
func (w *consoleLogWriter) Write(p []byte) (n int, err error) {
	// 快速健康检查（无锁）
	if atomic.LoadInt32(&w.healthyAtomic) == 0 {
		return 0, fmt.Errorf("console writer is not healthy")
	}

	// 只锁写入操作
	w.mutex.Lock()
	n, err = w.output.Write(p)
	w.mutex.Unlock()

	if err != nil {
		w.stats.addError()
		return n, err
	}

	w.stats.addBytes(int64(n))
	return n, nil
}

// WriteLevel 按级别写入
func (w *consoleLogWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.level {
		return len(data), nil // 跳过低级别日志
	}
	return w.Write(data)
}

// Flush 刷新缓冲区
func (w *consoleLogWriter) Flush() error {
	if flusher, ok := w.output.(interface{ Flush() error }); ok {
		return flusher.Flush()
	}
	return nil
}

// Close 关闭输出器
func (w *consoleLogWriter) Close() error {
	atomic.StoreInt32(&w.healthyAtomic, 0)
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.healthy = false
	if closer, ok := w.output.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// IsHealthy 检查健康状态（使用 atomic 快速检查）
func (w *consoleLogWriter) IsHealthy() bool {
	return atomic.LoadInt32(&w.healthyAtomic) == 1
}

// GetStats 获取统计信息
func (w *consoleLogWriter) GetStats() WriterStatsSnapshot {
	return w.stats.getSnapshot()
}

// FileLogWriter 文件输出器（支持可配置缓冲区大小）
type FileLogWriter struct {
	baseWriter                  // 继承基础输出器字段
	filePath      string        // 日志文件路径
	file          *os.File      // 文件句柄
	buffer        *bufio.Writer // 内部创建的缓冲区
	bufferSize    int           // 缓冲区大小（字节，默认 64KB）
	healthyAtomic int32         // 健康状态（atomic bool: 0=false, 1=true）
}

// FileWriterOption 文件输出器配置选项
type FileWriterOption func(*FileLogWriter)

// WithFileLevel 设置日志级别
func WithFileLevel(level LogLevel) FileWriterOption {
	return func(w *FileLogWriter) {
		w.level = level
	}
}

// WithFileWriterPath 设置文件路径
func WithFileWriterPath(filePath string) FileWriterOption {
	return func(w *FileLogWriter) {
		w.filePath = filePath
	}
}

// WithFilePermission 设置文件权限
func WithFilePermission(permission os.FileMode) FileWriterOption {
	return func(w *FileLogWriter) {
		w.permission = permission
	}
}

// WithFileBufferSize 设置缓冲区大小（字节）
func WithFileBufferSize(size int) FileWriterOption {
	return func(w *FileLogWriter) {
		if size > 0 {
			w.bufferSize = size
		}
	}
}

// NewFileWriter 创建文件输出器
func NewFileWriter(opts ...FileWriterOption) IWriter {
	w := &FileLogWriter{
		baseWriter: baseWriter{
			level:      DEBUG,
			healthy:    false, // 需要先打开文件
			stats:      newWriterStats(),
			permission: DefaultFilePermission,
		},
		bufferSize: 64 * 1024, // 默认 64KB
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// ensureFile 确保文件已打开（添加缓冲层）
func (w *FileLogWriter) ensureFile() error {
	if w.file != nil && w.buffer != nil {
		return nil
	}

	// 创建目录
	dir := filepath.Dir(w.filePath)
	if err := os.MkdirAll(dir, DefaultDirPermission); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 打开文件
	file, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, w.permission)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	w.file = file
	w.buffer = bufio.NewWriterSize(file, w.bufferSize)
	w.healthy = true
	atomic.StoreInt32(&w.healthyAtomic, 1)
	return nil
}

// Write 实现io.Writer接口（写入缓冲区）
func (w *FileLogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if err := w.ensureFile(); err != nil {
		w.stats.addError()
		return 0, err
	}

	n, err = w.buffer.Write(p)
	if err != nil {
		w.stats.addError()
		w.healthy = false
		atomic.StoreInt32(&w.healthyAtomic, 0)
		return n, err
	}

	w.stats.addBytes(int64(n))
	return n, nil
}

// WriteLevel 按级别写入
func (w *FileLogWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.level {
		return len(data), nil
	}
	return w.Write(data)
}

// Flush 刷新文件缓冲区（先刷新缓冲再同步文件）
func (w *FileLogWriter) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.buffer != nil {
		if err := w.buffer.Flush(); err != nil {
			return err
		}
	}
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

// Close 关闭文件（确保缓冲刷新）
func (w *FileLogWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.healthy = false
	atomic.StoreInt32(&w.healthyAtomic, 0)

	if w.buffer != nil {
		w.buffer.Flush()
		w.buffer = nil
	}

	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}
	return nil
}

// IsHealthy 检查健康状态（使用 atomic 快速检查）
func (w *FileLogWriter) IsHealthy() bool {
	return atomic.LoadInt32(&w.healthyAtomic) == 1
}

// GetStats 获取统计信息
func (w *FileLogWriter) GetStats() WriterStatsSnapshot {
	return w.stats.getSnapshot()
}

// RotateLogWriter 轮转文件输出器（支持按大小自动轮转，支持可配置缓冲区）
type RotateLogWriter struct {
	baseWriter                  // 继承基础输出器字段
	filePath      string        // 日志文件路径
	maxSize       int64         // 单个文件最大字节数（超过后轮转）
	maxFiles      int           // 最大保留文件数（旧文件会被删除）
	currentFile   *os.File      // 当前文件句柄
	currentSize   int64         // 当前文件已写入字节数
	buffer        *bufio.Writer // 内部创建的缓冲区
	bufferSize    int           // 缓冲区大小（字节，默认 64KB）
	healthyAtomic int32         // 健康状态（atomic bool: 0=false, 1=true）
}

// RotateWriterOption 轮转文件输出器配置选项
type RotateWriterOption func(*RotateLogWriter)

// WithRotateLevel 设置日志级别
func WithRotateLevel(level LogLevel) RotateWriterOption {
	return func(w *RotateLogWriter) {
		w.level = level
	}
}

// WithFilePath 设置文件路径
func WithFilePath(filePath string) RotateWriterOption {
	return func(w *RotateLogWriter) {
		w.filePath = filePath
	}
}

// WithMaxSize 设置最大文件大小（字节）
func WithMaxSize(maxSize int64) RotateWriterOption {
	return func(w *RotateLogWriter) {
		w.maxSize = maxSize
	}
}

// WithMaxFiles 设置最大文件数
func WithMaxFiles(maxFiles int) RotateWriterOption {
	return func(w *RotateLogWriter) {
		w.maxFiles = maxFiles
	}
}

// WithMaxAge 设置最大保留时间
func WithMaxAge(maxAge time.Duration) RotateWriterOption {
	return func(w *RotateLogWriter) {
		w.maxAge = maxAge
	}
}

// WithCompress 设置是否压缩旧文件
func WithCompress(compress bool) RotateWriterOption {
	return func(w *RotateLogWriter) {
		w.compress = compress
	}
}

// WithRotatePermission 设置文件权限
func WithRotatePermission(permission os.FileMode) RotateWriterOption {
	return func(w *RotateLogWriter) {
		w.permission = permission
	}
}

// WithRotateBufferSize 设置缓冲区大小（字节）
func WithRotateBufferSize(size int) RotateWriterOption {
	return func(w *RotateLogWriter) {
		if size > 0 {
			w.bufferSize = size
		}
	}
}

// NewRotateWriter 创建轮转文件输出器
func NewRotateWriter(opts ...RotateWriterOption) IWriter {
	w := &RotateLogWriter{
		baseWriter: baseWriter{
			level:      DEBUG,
			healthy:    false,
			stats:      newWriterStats(),
			permission: DefaultFilePermission,
			maxAge:     DefaultMaxAge,
			compress:   false,
		},
		maxSize:    DefaultMaxSize,
		maxFiles:   DefaultMaxFiles,
		bufferSize: 64 * 1024, // 默认 64KB
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// shouldRotate 检查是否需要轮转
func (w *RotateLogWriter) shouldRotate(dataSize int) bool {
	return w.currentSize+int64(dataSize) > w.maxSize
}

// rotate 执行文件轮转（刷新缓冲后轮转）
func (w *RotateLogWriter) rotate() error {
	// 刷新并关闭缓冲
	if w.buffer != nil {
		w.buffer.Flush()
		w.buffer = nil
	}
	if w.currentFile != nil {
		w.currentFile.Close()
		w.currentFile = nil
	}

	// 重命名现有文件
	for i := w.maxFiles - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", w.filePath, i)
		newPath := fmt.Sprintf("%s.%d", w.filePath, i+1)

		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}

	// 移动当前文件
	if _, err := os.Stat(w.filePath); err == nil {
		os.Rename(w.filePath, w.filePath+".1")
	}

	// 重置大小
	w.currentSize = 0

	return w.ensureFile()
}

// ensureFile 确保文件已打开（添加缓冲层）
func (w *RotateLogWriter) ensureFile() error {
	if w.currentFile != nil && w.buffer != nil {
		return nil
	}

	dir := filepath.Dir(w.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.OpenFile(w.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, w.permission)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	// 获取当前文件大小
	if stat, err := file.Stat(); err == nil {
		w.currentSize = stat.Size()
	}

	w.currentFile = file
	w.buffer = bufio.NewWriterSize(file, w.bufferSize)
	w.healthy = true
	atomic.StoreInt32(&w.healthyAtomic, 1)
	return nil
}

// Write 实现io.Writer接口（写入缓冲区）
func (w *RotateLogWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 检查是否需要轮转
	if w.shouldRotate(len(p)) {
		if err := w.rotate(); err != nil {
			w.stats.addError()
			return 0, err
		}
	}

	if err := w.ensureFile(); err != nil {
		w.stats.addError()
		return 0, err
	}

	n, err = w.buffer.Write(p)
	if err != nil {
		w.stats.addError()
		w.healthy = false
		atomic.StoreInt32(&w.healthyAtomic, 0)
		return n, err
	}

	w.currentSize += int64(n)
	w.stats.addBytes(int64(n))
	return n, nil
}

// WriteLevel 按级别写入
func (w *RotateLogWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.level {
		return len(data), nil
	}
	return w.Write(data)
}

// Flush 刷新缓冲区（刷新缓冲和文件）
func (w *RotateLogWriter) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.buffer != nil {
		if err := w.buffer.Flush(); err != nil {
			return err
		}
	}
	if w.currentFile != nil {
		return w.currentFile.Sync()
	}
	return nil
}

// Close 关闭输出器（确保缓冲刷新）
func (w *RotateLogWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.healthy = false
	atomic.StoreInt32(&w.healthyAtomic, 0)

	if w.buffer != nil {
		w.buffer.Flush()
		w.buffer = nil
	}

	if w.currentFile != nil {
		err := w.currentFile.Close()
		w.currentFile = nil
		return err
	}
	return nil
}

// IsHealthy 检查健康状态（使用 atomic 快速检查）
func (w *RotateLogWriter) IsHealthy() bool {
	return atomic.LoadInt32(&w.healthyAtomic) == 1
}

// GetStats 获取统计信息
func (w *RotateLogWriter) GetStats() WriterStatsSnapshot {
	return w.stats.getSnapshot()
}

// BufferedWriter 缓冲输出器（为底层输出器添加可配置的缓冲层）
type BufferedWriter struct {
	baseWriter               // 继承基础输出器字段
	underlying IWriter       // 底层输出器（实际写入目标）
	buffer     *bufio.Writer // 基于 underlying 创建的缓冲区
	bufferSize int           // 缓冲区大小（字节，可配置）
}

// BufferedWriterOption 缓冲输出器配置选项
type BufferedWriterOption func(*BufferedWriter)

// WithBufferedUnderlying 设置底层输出器
func WithBufferedUnderlying(underlying IWriter) BufferedWriterOption {
	return func(w *BufferedWriter) {
		w.underlying = underlying
		w.buffer = bufio.NewWriterSize(underlying, w.bufferSize)
	}
}

// WithBufferSize 设置缓冲区大小
func WithBufferSize(bufferSize int) BufferedWriterOption {
	return func(w *BufferedWriter) {
		w.bufferSize = bufferSize
		if w.underlying != nil {
			w.buffer = bufio.NewWriterSize(w.underlying, bufferSize)
		}
	}
}

// WithBufferedLevel 设置日志级别
func WithBufferedLevel(level LogLevel) BufferedWriterOption {
	return func(w *BufferedWriter) {
		w.level = level
	}
}

// NewBufferedWriter 创建缓冲输出器
func NewBufferedWriter(opts ...BufferedWriterOption) IWriter {
	w := &BufferedWriter{
		baseWriter: baseWriter{
			level:   DEBUG,
			healthy: true,
			stats:   newWriterStats(),
		},
		bufferSize: DefaultBufferSize,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// Write 实现io.Writer接口
func (w *BufferedWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.healthy {
		return 0, fmt.Errorf("buffered writer is not healthy")
	}

	n, err = w.buffer.Write(p)
	if err != nil {
		w.stats.addError()
		w.healthy = false
		return n, err
	}

	w.stats.addBytes(int64(n))
	return n, nil
}

// WriteLevel 按级别写入
func (w *BufferedWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.level {
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

	w.healthy = false
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
	return w.healthy && w.underlying.IsHealthy()
}

// GetStats 获取统计信息
func (w *BufferedWriter) GetStats() WriterStatsSnapshot {
	return w.stats.getSnapshot()
}

// MultiLogWriter 多输出器（同时写入多个输出器，实现日志分发）
type MultiLogWriter struct {
	baseWriter           // 继承基础输出器字段
	writers    []IWriter // 输出器列表（日志会写入所有健康的输出器）
}

// MultiWriterOption 多输出器配置选项
type MultiWriterOption func(*MultiLogWriter)

// WithWriters 设置输出器列表
func WithWriters(writers ...IWriter) MultiWriterOption {
	return func(w *MultiLogWriter) {
		w.writers = writers
	}
}

// WithMultiLevel 设置日志级别
func WithMultiLevel(level LogLevel) MultiWriterOption {
	return func(w *MultiLogWriter) {
		w.level = level
	}
}

// NewMultiWriter 创建多输出器
func NewMultiWriter(opts ...MultiWriterOption) IWriter {
	w := &MultiLogWriter{
		baseWriter: baseWriter{
			level:   DEBUG,
			healthy: true,
			stats:   newWriterStats(),
		},
		writers: []IWriter{},
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
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
			w.stats.addError()
		}
	}

	if lastErr != nil {
		return 0, lastErr
	}

	w.stats.addBytes(int64(len(p)))
	return len(p), nil
}

// WriteLevel 按级别写入
func (w *MultiLogWriter) WriteLevel(level LogLevel, data []byte) (n int, err error) {
	if level < w.level {
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

	w.healthy = false
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

	if !w.healthy {
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
func (w *MultiLogWriter) GetStats() WriterStatsSnapshot {
	return w.stats.getSnapshot()
}
