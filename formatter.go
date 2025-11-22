/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\formatter.go
 * @Description: 日志格式化器实现
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// FormatterType 格式化器类型
type FormatterType string

const (
	TextFormatter FormatterType = "text"
	JSONFormatter FormatterType = "json"
	XMLFormatter  FormatterType = "xml"
	CSVFormatter  FormatterType = "csv"
)

// BaseFormatter 基础格式化器
type BaseFormatter struct {
	TimeFormat   string            `json:"time_format"`
	ShowCaller   bool              `json:"show_caller"`
	ShowLevel    bool              `json:"show_level"`
	ShowEmoji    bool              `json:"show_emoji"`
	Colorful     bool              `json:"colorful"`
	PrettyPrint  bool              `json:"pretty_print"`
	CustomFields map[string]string `json:"custom_fields"`
}

// NewBaseFormatter 创建基础格式化器
func NewBaseFormatter() *BaseFormatter {
	return &BaseFormatter{
		TimeFormat:   "2006-01-02 15:04:05.000",
		ShowCaller:   true,
		ShowLevel:    true,
		ShowEmoji:    true,
		Colorful:     true,
		PrettyPrint:  false,
		CustomFields: make(map[string]string),
	}
}

// TextLogFormatter 文本格式化器
type TextLogFormatter struct {
	*BaseFormatter
	Template string `json:"template"`
}

// NewTextFormatter 创建文本格式化器
func NewTextFormatter() IFormatter {
	return &TextLogFormatter{
		BaseFormatter: NewBaseFormatter(),
		Template:      "${time} ${level} ${caller} ${message} ${fields}",
	}
}

// GetName 获取格式化器名称
func (f *TextLogFormatter) GetName() string {
	return "text"
}

// Format 格式化日志条目
func (f *TextLogFormatter) Format(entry *LogEntry) ([]byte, error) {
	var parts []string

	// 时间
	if f.TimeFormat != "" {
		timeStr := time.Unix(entry.Timestamp, 0).Format(f.TimeFormat)
		parts = append(parts, timeStr)
	}

	// 级别
	if f.ShowLevel {
		levelStr := entry.Level.String()
		if f.ShowEmoji {
			levelStr = fmt.Sprintf("%s [%s]", entry.Level.Emoji(), levelStr)
		} else {
			levelStr = fmt.Sprintf("[%s]", levelStr)
		}

		if f.Colorful {
			levelStr = entry.Level.Color() + levelStr + "\033[0m"
		}

		parts = append(parts, levelStr)
	}

	// 调用者信息
	if f.ShowCaller && entry.Caller != nil {
		callerStr := fmt.Sprintf("[%s:%d:%s]", entry.Caller.File, entry.Caller.Line, entry.Caller.Function)
		parts = append(parts, callerStr)
	}

	// 消息
	parts = append(parts, entry.Message)

	// 字段
	if len(entry.Fields) > 0 {
		var fieldParts []string
		for key, value := range entry.Fields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", key, value))
		}
		if len(fieldParts) > 0 {
			parts = append(parts, fmt.Sprintf("{%s}", strings.Join(fieldParts, ", ")))
		}
	}

	result := strings.Join(parts, " ") + "\n"
	return []byte(result), nil
}

// JSONLogFormatter JSON格式化器
type JSONLogFormatter struct {
	*BaseFormatter
	FieldMap map[string]string `json:"field_map"`
}

// NewJSONFormatter 创建JSON格式化器
func NewJSONFormatter() IFormatter {
	return &JSONLogFormatter{
		BaseFormatter: NewBaseFormatter(),
		FieldMap: map[string]string{
			"time":    "timestamp",
			"level":   "level",
			"message": "message",
			"caller":  "caller",
			"fields":  "fields",
		},
	}
}

// GetName 获取格式化器名称
func (f *JSONLogFormatter) GetName() string {
	return "json"
}

// Format 格式化日志条目
func (f *JSONLogFormatter) Format(entry *LogEntry) ([]byte, error) {
	data := make(map[string]interface{})

	// 基础字段
	data[f.getFieldName("time")] = time.Unix(entry.Timestamp, 0).Format(f.TimeFormat)
	data[f.getFieldName("level")] = entry.Level.String()
	data[f.getFieldName("message")] = entry.Message

	// 调用者信息
	if f.ShowCaller && entry.Caller != nil {
		data[f.getFieldName("caller")] = map[string]interface{}{
			"file":     entry.Caller.File,
			"line":     entry.Caller.Line,
			"function": entry.Caller.Function,
		}
	}

	// 自定义字段
	if len(entry.Fields) > 0 {
		if f.getFieldName("fields") == "fields" {
			data["fields"] = entry.Fields
		} else {
			// 展开字段到顶层
			for key, value := range entry.Fields {
				data[key] = value
			}
		}
	}

	// 添加自定义字段
	for key, value := range f.CustomFields {
		data[key] = value
	}

	var result []byte
	var err error

	if f.PrettyPrint {
		result, err = json.MarshalIndent(data, "", "  ")
	} else {
		result, err = json.Marshal(data)
	}

	if err != nil {
		return nil, err
	}

	return append(result, '\n'), nil
}

// getFieldName 获取字段映射名称
func (f *JSONLogFormatter) getFieldName(field string) string {
	if name, exists := f.FieldMap[field]; exists {
		return name
	}
	return field
}

// XMLLogFormatter XML格式化器
type XMLLogFormatter struct {
	*BaseFormatter
	RootElement string `json:"root_element"`
	LogElement  string `json:"log_element"`
}

// NewXMLFormatter 创建XML格式化器
func NewXMLFormatter() IFormatter {
	return &XMLLogFormatter{
		BaseFormatter: NewBaseFormatter(),
		RootElement:   "log",
		LogElement:    "entry",
	}
}

// GetName 获取格式化器名称
func (f *XMLLogFormatter) GetName() string {
	return "xml"
}

// Format 格式化日志条目
func (f *XMLLogFormatter) Format(entry *LogEntry) ([]byte, error) {
	var parts []string

	parts = append(parts, fmt.Sprintf("<%s>", f.LogElement))

	// 时间
	if f.TimeFormat != "" {
		timeStr := time.Unix(entry.Timestamp, 0).Format(f.TimeFormat)
		parts = append(parts, fmt.Sprintf("  <time>%s</time>", timeStr))
	}

	// 级别
	if f.ShowLevel {
		parts = append(parts, fmt.Sprintf("  <level>%s</level>", entry.Level.String()))
	}

	// 消息
	parts = append(parts, fmt.Sprintf("  <message><![CDATA[%s]]></message>", entry.Message))

	// 调用者信息
	if f.ShowCaller && entry.Caller != nil {
		parts = append(parts, "  <caller>")
		parts = append(parts, fmt.Sprintf("    <file>%s</file>", entry.Caller.File))
		parts = append(parts, fmt.Sprintf("    <line>%d</line>", entry.Caller.Line))
		parts = append(parts, fmt.Sprintf("    <function>%s</function>", entry.Caller.Function))
		parts = append(parts, "  </caller>")
	}

	// 字段
	if len(entry.Fields) > 0 {
		parts = append(parts, "  <fields>")
		for key, value := range entry.Fields {
			parts = append(parts, fmt.Sprintf("    <%s><![CDATA[%v]]></%s>", key, value, key))
		}
		parts = append(parts, "  </fields>")
	}

	parts = append(parts, fmt.Sprintf("</%s>", f.LogElement))

	result := strings.Join(parts, "\n") + "\n"
	return []byte(result), nil
}

// CSVLogFormatter CSV格式化器
type CSVLogFormatter struct {
	*BaseFormatter
	Headers   []string `json:"headers"`
	Delimiter string   `json:"delimiter"`
}

// NewCSVFormatter 创建CSV格式化器
func NewCSVFormatter() IFormatter {
	return &CSVLogFormatter{
		BaseFormatter: NewBaseFormatter(),
		Headers:       []string{"time", "level", "message", "file", "line", "function"},
		Delimiter:     ",",
	}
}

// GetName 获取格式化器名称
func (f *CSVLogFormatter) GetName() string {
	return "csv"
}

// Format 格式化日志条目
func (f *CSVLogFormatter) Format(entry *LogEntry) ([]byte, error) {
	var values []string

	for _, header := range f.Headers {
		switch header {
		case "time":
			values = append(values, time.Unix(entry.Timestamp, 0).Format(f.TimeFormat))
		case "level":
			values = append(values, entry.Level.String())
		case "message":
			values = append(values, fmt.Sprintf("\"%s\"", strings.ReplaceAll(entry.Message, "\"", "\"\"")))
		case "file":
			if entry.Caller != nil {
				values = append(values, entry.Caller.File)
			} else {
				values = append(values, "")
			}
		case "line":
			if entry.Caller != nil {
				values = append(values, fmt.Sprintf("%d", entry.Caller.Line))
			} else {
				values = append(values, "")
			}
		case "function":
			if entry.Caller != nil {
				values = append(values, entry.Caller.Function)
			} else {
				values = append(values, "")
			}
		default:
			// 尝试从字段中获取
			if value, exists := entry.Fields[header]; exists {
				values = append(values, fmt.Sprintf("%v", value))
			} else {
				values = append(values, "")
			}
		}
	}

	result := strings.Join(values, f.Delimiter) + "\n"
	return []byte(result), nil
}

// FormatRegistry 格式化器注册表
type FormatRegistry struct {
	formatters map[FormatterType]func() IFormatter
}

// NewFormatRegistry 创建格式化器注册表
func NewFormatRegistry() *FormatRegistry {
	registry := &FormatRegistry{
		formatters: make(map[FormatterType]func() IFormatter),
	}

	// 注册默认格式化器
	registry.Register(TextFormatter, NewTextFormatter)
	registry.Register(JSONFormatter, NewJSONFormatter)
	registry.Register(XMLFormatter, NewXMLFormatter)
	registry.Register(CSVFormatter, NewCSVFormatter)

	return registry
}

// Register 注册格式化器
func (r *FormatRegistry) Register(formatterType FormatterType, factory func() IFormatter) {
	r.formatters[formatterType] = factory
}

// Create 创建格式化器
func (r *FormatRegistry) Create(formatterType FormatterType) (IFormatter, error) {
	factory, exists := r.formatters[formatterType]
	if !exists {
		return nil, fmt.Errorf("unknown formatter type: %s", formatterType)
	}
	return factory(), nil
}

// List 列出所有注册的格式化器
func (r *FormatRegistry) List() []FormatterType {
	var types []FormatterType
	for t := range r.formatters {
		types = append(types, t)
	}
	return types
}

// GetCallerInfo 获取调用者信息
func GetCallerInfo(skip int) *CallerInfo {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil
	}

	funcName := runtime.FuncForPC(pc).Name()
	if idx := strings.LastIndex(funcName, "."); idx != -1 {
		funcName = funcName[idx+1:]
	}
	if idx := strings.LastIndex(file, "/"); idx != -1 {
		file = file[idx+1:]
	}

	return &CallerInfo{
		File:     file,
		Line:     line,
		Function: funcName,
	}
}
