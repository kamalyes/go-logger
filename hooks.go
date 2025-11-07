/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\hooks.go
 * @Description: 日志钩子系统实现
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// HookLevel 钩子触发级别
type HookLevel []LogLevel

// AllLevels 所有级别
var AllLevels = HookLevel{DEBUG, INFO, WARN, ERROR, FATAL}

// ErrorLevels 错误级别
var ErrorLevels = HookLevel{ERROR, FATAL}

// WarnLevels 警告级别
var WarnLevels = HookLevel{WARN, ERROR, FATAL}

// BaseHook 基础钩子
type BaseHook struct {
	Name     string    `json:"name"`
	Enabled  bool      `json:"enabled"`
	levels   HookLevel
	Async    bool      `json:"async"`
	Timeout  time.Duration `json:"timeout"`
	mutex    sync.RWMutex
}

// NewBaseHook 创建基础钩子
func NewBaseHook(name string, levels HookLevel) *BaseHook {
	return &BaseHook{
		Name:    name,
		Enabled: true,
		levels:  levels,
		Async:   false,
		Timeout: 30 * time.Second,
	}
}

// IsEnabled 检查钩子是否启用
func (h *BaseHook) IsEnabled() bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.Enabled
}

// SetEnabled 设置钩子启用状态
func (h *BaseHook) SetEnabled(enabled bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Enabled = enabled
}

// Levels 返回支持的日志级别
func (h *BaseHook) Levels() []LogLevel {
	return h.levels
}

// SupportsLevel 检查是否支持指定级别
func (h *BaseHook) SupportsLevel(level LogLevel) bool {
	for _, l := range h.levels {
		if l == level {
			return true
		}
	}
	return false
}

// ConsoleHook 控制台钩子（用于特殊输出）
type ConsoleHook struct {
	*BaseHook
	Prefix string `json:"prefix"`
	Color  bool   `json:"color"`
}

// NewConsoleHook 创建控制台钩子
func NewConsoleHook(levels HookLevel) IHook {
	return &ConsoleHook{
		BaseHook: NewBaseHook("console", levels),
		Prefix:   "[HOOK]",
		Color:    true,
	}
}

// Fire 执行钩子
func (h *ConsoleHook) Fire(entry *LogEntry) error {
	if entry == nil {
		return fmt.Errorf("log entry cannot be nil")
	}
	
	if !h.IsEnabled() || !h.SupportsLevel(entry.Level) {
		return nil
	}
	
	timestamp := time.Unix(0, entry.Timestamp).Format("2006-01-02 15:04:05")
	levelStr := entry.Level.String()
	
	if h.Color {
		levelStr = entry.Level.Color() + levelStr + "\033[0m"
	}
	
	output := fmt.Sprintf("%s %s [%s] %s\n", 
		timestamp, h.Prefix, levelStr, entry.Message)
	
	if len(entry.Fields) > 0 {
		fieldsJson, _ := json.Marshal(entry.Fields)
		output += fmt.Sprintf("%s Fields: %s\n", h.Prefix, string(fieldsJson))
	}
	
	fmt.Print(output)
	return nil
}

// FileHook 文件钩子
type FileHook struct {
	*BaseHook
	FilePath   string      `json:"file_path"`
	MaxSize    int64       `json:"max_size"`
	MaxBackups int         `json:"max_backups"`
	file       *os.File
	currentSize int64
	mutex      sync.Mutex
}

// NewFileHook 创建文件钩子
func NewFileHook(filePath string, levels HookLevel) IHook {
	return &FileHook{
		BaseHook:   NewBaseHook("file", levels),
		FilePath:   filePath,
		MaxSize:    100 * 1024 * 1024, // 100MB
		MaxBackups: 5,
	}
}

// ensureFile 确保文件已打开
func (h *FileHook) ensureFile() error {
	if h.file != nil {
		return nil
	}
	
	dir := filepath.Dir(h.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	file, err := os.OpenFile(h.FilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	
	if stat, err := file.Stat(); err == nil {
		h.currentSize = stat.Size()
	}
	
	h.file = file
	return nil
}

// rotate 轮转文件
func (h *FileHook) rotate() error {
	if h.file != nil {
		h.file.Close()
		h.file = nil
	}
	
	// 轮转文件
	for i := h.MaxBackups - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", h.FilePath, i)
		newPath := fmt.Sprintf("%s.%d", h.FilePath, i+1)
		
		if _, err := os.Stat(oldPath); err == nil {
			os.Rename(oldPath, newPath)
		}
	}
	
	if _, err := os.Stat(h.FilePath); err == nil {
		os.Rename(h.FilePath, h.FilePath+".1")
	}
	
	h.currentSize = 0
	return h.ensureFile()
}

// Fire 执行钩子
func (h *FileHook) Fire(entry *LogEntry) error {
	if entry == nil {
		return fmt.Errorf("log entry cannot be nil")
	}
	
	if !h.IsEnabled() || !h.SupportsLevel(entry.Level) {
		return nil
	}
	
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	// 构建日志内容
	timestamp := time.Unix(0, entry.Timestamp).Format("2006-01-02 15:04:05.000")
	content := fmt.Sprintf("%s [%s] %s", timestamp, entry.Level.String(), entry.Message)
	
	if len(entry.Fields) > 0 {
		fieldsJson, _ := json.Marshal(entry.Fields)
		content += " " + string(fieldsJson)
	}
	
	if entry.Caller != nil {
		content += fmt.Sprintf(" [%s:%d]", entry.Caller.File, entry.Caller.Line)
	}
	
	content += "\n"
	
	// 检查是否需要轮转
	if h.currentSize+int64(len(content)) > h.MaxSize {
		if err := h.rotate(); err != nil {
			return err
		}
	}
	
	if err := h.ensureFile(); err != nil {
		return err
	}
	
	n, err := h.file.WriteString(content)
	if err != nil {
		return err
	}
	
	h.currentSize += int64(n)
	return h.file.Sync()
}

// EmailHook 邮件钩子
type EmailHook struct {
	*BaseHook
	SMTPHost     string   `json:"smtp_host"`
	SMTPPort     string   `json:"smtp_port"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	From         string   `json:"from"`
	To           []string `json:"to"`
	Subject      string   `json:"subject"`
	BatchSize    int      `json:"batch_size"`
	FlushTimeout time.Duration `json:"flush_timeout"`
	buffer       []*LogEntry
	lastFlush    time.Time
	bufferMutex  sync.Mutex
}

// NewEmailHook 创建邮件钩子
func NewEmailHook(host, port, username, password, from string, to []string) IHook {
	hook := &EmailHook{
		BaseHook:     NewBaseHook("email", ErrorLevels),
		SMTPHost:     host,
		SMTPPort:     port,
		Username:     username,
		Password:     password,
		From:         from,
		To:           to,
		Subject:      "Application Log Alert",
		BatchSize:    10,
		FlushTimeout: 5 * time.Minute,
		buffer:       make([]*LogEntry, 0),
		lastFlush:    time.Now(),
	}
	
	// 启动定时刷新
	go hook.flushPeriodically()
	
	return hook
}

// Fire 执行钩子
func (h *EmailHook) Fire(entry *LogEntry) error {
	if !h.IsEnabled() || !h.SupportsLevel(entry.Level) {
		return nil
	}
	
	h.bufferMutex.Lock()
	defer h.bufferMutex.Unlock()
	
	h.buffer = append(h.buffer, entry)
	
	// 检查是否需要立即发送
	if len(h.buffer) >= h.BatchSize {
		return h.flush()
	}
	
	return nil
}

// flushPeriodically 定期刷新缓冲区
func (h *EmailHook) flushPeriodically() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		h.bufferMutex.Lock()
		if len(h.buffer) > 0 && time.Since(h.lastFlush) > h.FlushTimeout {
			h.flush()
		}
		h.bufferMutex.Unlock()
	}
}

// flush 发送邮件
func (h *EmailHook) flush() error {
	if len(h.buffer) == 0 {
		return nil
	}
	
	// 构建邮件内容
	var content strings.Builder
	content.WriteString("Application Log Report\n")
	content.WriteString("======================\n\n")
	
	for _, entry := range h.buffer {
		timestamp := time.Unix(0, entry.Timestamp).Format("2006-01-02 15:04:05")
		content.WriteString(fmt.Sprintf("%s [%s] %s\n", 
			timestamp, entry.Level.String(), entry.Message))
		
		if len(entry.Fields) > 0 {
			fieldsJson, _ := json.MarshalIndent(entry.Fields, "  ", "  ")
			content.WriteString(fmt.Sprintf("  Fields: %s\n", string(fieldsJson)))
		}
		
		if entry.Caller != nil {
			content.WriteString(fmt.Sprintf("  Caller: %s:%d:%s\n", 
				entry.Caller.File, entry.Caller.Line, entry.Caller.Function))
		}
		
		content.WriteString("\n")
	}
	
	// 发送邮件
	err := h.sendEmail(h.Subject, content.String())
	if err != nil {
		return err
	}
	
	// 清空缓冲区
	h.buffer = h.buffer[:0]
	h.lastFlush = time.Now()
	
	return nil
}

// sendEmail 发送邮件
func (h *EmailHook) sendEmail(subject, body string) error {
	auth := smtp.PlainAuth("", h.Username, h.Password, h.SMTPHost)
	
	msg := []byte("To: " + strings.Join(h.To, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body)
	
	return smtp.SendMail(h.SMTPHost+":"+h.SMTPPort, auth, h.From, h.To, msg)
}

// WebhookHook HTTP Webhook钩子
type WebhookHook struct {
	*BaseHook
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers"`
	Timeout     time.Duration     `json:"timeout"`
	MaxRetries  int               `json:"max_retries"`
	RetryDelay  time.Duration     `json:"retry_delay"`
	client      *http.Client
}

// NewWebhookHook 创建Webhook钩子
func NewWebhookHook(url string, levels HookLevel) IHook {
	return &WebhookHook{
		BaseHook:   NewBaseHook("webhook", levels),
		URL:        url,
		Method:     "POST",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

// Fire 执行钩子
func (h *WebhookHook) Fire(entry *LogEntry) error {
	if entry == nil {
		return fmt.Errorf("log entry cannot be nil")
	}
	
	if !h.IsEnabled() || !h.SupportsLevel(entry.Level) {
		return nil
	}
	
	// 构建请求数据
	data := map[string]interface{}{
		"timestamp": time.Unix(0, entry.Timestamp).Format(time.RFC3339),
		"level":     entry.Level.String(),
		"message":   entry.Message,
		"fields":    entry.Fields,
	}
	
	if entry.Caller != nil {
		data["caller"] = map[string]interface{}{
			"file":     entry.Caller.File,
			"line":     entry.Caller.Line,
			"function": entry.Caller.Function,
		}
	}
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	// 发送请求（带重试）
	for retry := 0; retry <= h.MaxRetries; retry++ {
		if retry > 0 {
			time.Sleep(h.RetryDelay * time.Duration(retry))
		}
		
		req, err := http.NewRequest(h.Method, h.URL, bytes.NewBuffer(jsonData))
		if err != nil {
			continue
		}
		
		for key, value := range h.Headers {
			req.Header.Set(key, value)
		}
		
		resp, err := h.client.Do(req)
		if err != nil {
			if retry == h.MaxRetries {
				return err
			}
			continue
		}
		
		resp.Body.Close()
		
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		
		if retry == h.MaxRetries {
			return fmt.Errorf("webhook failed with status: %d", resp.StatusCode)
		}
	}
	
	return nil
}

// HookManager 钩子管理器
type HookManager struct {
	hooks map[string]IHook
	mutex sync.RWMutex
}

// NewHookManager 创建钩子管理器
func NewHookManager() *HookManager {
	return &HookManager{
		hooks: make(map[string]IHook),
	}
}

// AddHook 添加钩子
func (hm *HookManager) AddHook(name string, hook IHook) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	hm.hooks[name] = hook
}

// RemoveHook 移除钩子
func (hm *HookManager) RemoveHook(name string) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	delete(hm.hooks, name)
}

// GetHook 获取钩子
func (hm *HookManager) GetHook(name string) (IHook, bool) {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	hook, exists := hm.hooks[name]
	return hook, exists
}

// FireHooks 触发所有适用的钩子
func (hm *HookManager) FireHooks(entry *LogEntry) {
	hm.mutex.RLock()
	hooks := make([]IHook, 0, len(hm.hooks))
	for _, hook := range hm.hooks {
		hooks = append(hooks, hook)
	}
	hm.mutex.RUnlock()
	
	for _, hook := range hooks {
		if hook != nil {
			go func(h IHook) {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Hook panic: %v\n", r)
					}
				}()
				h.Fire(entry)
			}(hook)
		}
	}
}

// ListHooks 列出所有钩子名称
func (hm *HookManager) ListHooks() []string {
	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	
	names := make([]string, 0, len(hm.hooks))
	for name := range hm.hooks {
		names = append(names, name)
	}
	return names
}