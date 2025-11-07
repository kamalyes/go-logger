/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 23:36:45
 * @FilePath: \go-logger\level\constants.go
 * @Description: 日志级别常量定义
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package level

import (
	"fmt"
	"strings"
)

// Level 日志级别类型
type Level int

// 基础日志级别常量
const (
	TRACE Level = iota - 2 // 跟踪级别 (-2)
	DEBUG                  // 调试级别 (-1)
	INFO                   // 信息级别 (0)
	WARN                   // 警告级别 (1)
	ERROR                  // 错误级别 (2)
	FATAL                  // 致命级别 (3)
	OFF                    // 关闭级别 (4)
)

// 扩展日志级别
const (
	// 系统级别
	SYSTEM Level = 100 + iota
	KERNEL
	DRIVER
	
	// 应用级别
	APPLICATION Level = 200 + iota
	SERVICE
	COMPONENT
	MODULE
	
	// 业务级别
	BUSINESS Level = 300 + iota
	TRANSACTION
	WORKFLOW
	PROCESS
	
	// 安全级别
	SECURITY Level = 400 + iota
	AUDIT
	COMPLIANCE
	THREAT
	
	// 性能级别
	PERFORMANCE Level = 500 + iota
	METRIC
	BENCHMARK
	PROFILING
)

// LevelInfo 级别信息结构
type LevelInfo struct {
	Name        string `json:"name"`        // 级别名称
	ShortName   string `json:"short_name"`  // 短名称
	Emoji       string `json:"emoji"`       // 表情符号
	Color       string `json:"color"`       // 颜色代码
	Priority    int    `json:"priority"`    // 优先级
	Description string `json:"description"` // 描述
	Category    string `json:"category"`    // 类别
}

// 级别信息映射
var levelInfoMap = map[Level]LevelInfo{
	// 基础级别
	TRACE: {
		Name: "TRACE", ShortName: "TRC", Emoji: "🔍",
		Color: "\033[90m", Priority: -2, Description: "详细跟踪信息", Category: "basic",
	},
	DEBUG: {
		Name: "DEBUG", ShortName: "DBG", Emoji: "🐛",
		Color: "\033[36m", Priority: -1, Description: "调试信息", Category: "basic",
	},
	INFO: {
		Name: "INFO", ShortName: "INF", Emoji: "ℹ️",
		Color: "\033[32m", Priority: 0, Description: "一般信息", Category: "basic",
	},
	WARN: {
		Name: "WARN", ShortName: "WRN", Emoji: "⚠️",
		Color: "\033[33m", Priority: 1, Description: "警告信息", Category: "basic",
	},
	ERROR: {
		Name: "ERROR", ShortName: "ERR", Emoji: "❌",
		Color: "\033[31m", Priority: 2, Description: "错误信息", Category: "basic",
	},
	FATAL: {
		Name: "FATAL", ShortName: "FTL", Emoji: "💀",
		Color: "\033[35m", Priority: 3, Description: "致命错误", Category: "basic",
	},
	OFF: {
		Name: "OFF", ShortName: "OFF", Emoji: "🚫",
		Color: "\033[0m", Priority: 4, Description: "关闭日志", Category: "basic",
	},
	
	// 系统级别
	SYSTEM: {
		Name: "SYSTEM", ShortName: "SYS", Emoji: "🖥️",
		Color: "\033[94m", Priority: 100, Description: "系统级信息", Category: "system",
	},
	KERNEL: {
		Name: "KERNEL", ShortName: "KRN", Emoji: "⚙️",
		Color: "\033[95m", Priority: 101, Description: "内核级信息", Category: "system",
	},
	DRIVER: {
		Name: "DRIVER", ShortName: "DRV", Emoji: "🔌",
		Color: "\033[96m", Priority: 102, Description: "驱动级信息", Category: "system",
	},
	
	// 应用级别
	APPLICATION: {
		Name: "APPLICATION", ShortName: "APP", Emoji: "📱",
		Color: "\033[92m", Priority: 200, Description: "应用级信息", Category: "application",
	},
	SERVICE: {
		Name: "SERVICE", ShortName: "SVC", Emoji: "🔧",
		Color: "\033[93m", Priority: 201, Description: "服务级信息", Category: "application",
	},
	COMPONENT: {
		Name: "COMPONENT", ShortName: "CMP", Emoji: "🧩",
		Color: "\033[94m", Priority: 202, Description: "组件级信息", Category: "application",
	},
	MODULE: {
		Name: "MODULE", ShortName: "MOD", Emoji: "📦",
		Color: "\033[95m", Priority: 203, Description: "模块级信息", Category: "application",
	},
	
	// 业务级别
	BUSINESS: {
		Name: "BUSINESS", ShortName: "BIZ", Emoji: "💼",
		Color: "\033[38;5;214m", Priority: 300, Description: "业务级信息", Category: "business",
	},
	TRANSACTION: {
		Name: "TRANSACTION", ShortName: "TXN", Emoji: "💳",
		Color: "\033[38;5;215m", Priority: 301, Description: "事务级信息", Category: "business",
	},
	WORKFLOW: {
		Name: "WORKFLOW", ShortName: "WFL", Emoji: "🔄",
		Color: "\033[38;5;216m", Priority: 302, Description: "工作流信息", Category: "business",
	},
	PROCESS: {
		Name: "PROCESS", ShortName: "PRC", Emoji: "⚡",
		Color: "\033[38;5;217m", Priority: 303, Description: "流程级信息", Category: "business",
	},
	
	// 安全级别
	SECURITY: {
		Name: "SECURITY", ShortName: "SEC", Emoji: "🔒",
		Color: "\033[38;5;196m", Priority: 400, Description: "安全级信息", Category: "security",
	},
	AUDIT: {
		Name: "AUDIT", ShortName: "ADT", Emoji: "📋",
		Color: "\033[38;5;197m", Priority: 401, Description: "审计级信息", Category: "security",
	},
	COMPLIANCE: {
		Name: "COMPLIANCE", ShortName: "CMP", Emoji: "✅",
		Color: "\033[38;5;198m", Priority: 402, Description: "合规级信息", Category: "security",
	},
	THREAT: {
		Name: "THREAT", ShortName: "THR", Emoji: "🛡️",
		Color: "\033[38;5;199m", Priority: 403, Description: "威胁级信息", Category: "security",
	},
	
	// 性能级别
	PERFORMANCE: {
		Name: "PERFORMANCE", ShortName: "PRF", Emoji: "📊",
		Color: "\033[38;5;81m", Priority: 500, Description: "性能级信息", Category: "performance",
	},
	METRIC: {
		Name: "METRIC", ShortName: "MTC", Emoji: "📈",
		Color: "\033[38;5;82m", Priority: 501, Description: "指标级信息", Category: "performance",
	},
	BENCHMARK: {
		Name: "BENCHMARK", ShortName: "BMK", Emoji: "⏱️",
		Color: "\033[38;5;83m", Priority: 502, Description: "基准测试信息", Category: "performance",
	},
	PROFILING: {
		Name: "PROFILING", ShortName: "PRO", Emoji: "🔬",
		Color: "\033[38;5;84m", Priority: 503, Description: "性能分析信息", Category: "performance",
	},
}

// 级别名称映射
var levelNameMap = func() map[string]Level {
	m := make(map[string]Level)
	for level, info := range levelInfoMap {
		m[info.Name] = level
		m[info.ShortName] = level
	}
	// 添加别名
	m["WARNING"] = WARN
	m["CRITICAL"] = FATAL
	m["EMERGENCY"] = FATAL
	return m
}()

// 类别级别映射
var categoryLevelMap = func() map[string][]Level {
	m := make(map[string][]Level)
	for level, info := range levelInfoMap {
		m[info.Category] = append(m[info.Category], level)
	}
	return m
}()

// String 返回级别的字符串表示
func (l Level) String() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Name
	}
	return fmt.Sprintf("UNKNOWN(%d)", int(l))
}

// ShortString 返回级别的短字符串表示
func (l Level) ShortString() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.ShortName
	}
	return "UNK"
}

// Emoji 返回级别的表情符号
func (l Level) Emoji() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Emoji
	}
	return "❓"
}

// Color 返回级别的颜色代码
func (l Level) Color() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Color
	}
	return "\033[0m" // 重置颜色
}

// Priority 返回级别的优先级
func (l Level) Priority() int {
	if info, ok := levelInfoMap[l]; ok {
		return info.Priority
	}
	return 0
}

// Description 返回级别的描述
func (l Level) Description() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Description
	}
	return "Unknown level"
}

// Category 返回级别的类别
func (l Level) Category() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Category
	}
	return "unknown"
}

// Info 返回级别的完整信息
func (l Level) Info() LevelInfo {
	if info, ok := levelInfoMap[l]; ok {
		return info
	}
	return LevelInfo{
		Name: fmt.Sprintf("UNKNOWN(%d)", int(l)),
		ShortName: "UNK",
		Emoji: "❓",
		Color: "\033[0m",
		Priority: int(l),
		Description: "Unknown level",
		Category: "unknown",
	}
}

// IsValid 检查级别是否有效
func (l Level) IsValid() bool {
	_, ok := levelInfoMap[l]
	return ok
}

// IsEnabled 检查当前级别是否启用目标级别
func (l Level) IsEnabled(target Level) bool {
	return target.Priority() >= l.Priority()
}

// IsBasic 检查是否为基础级别
func (l Level) IsBasic() bool {
	return l.Category() == "basic"
}

// IsSystem 检查是否为系统级别
func (l Level) IsSystem() bool {
	return l.Category() == "system"
}

// IsApplication 检查是否为应用级别
func (l Level) IsApplication() bool {
	return l.Category() == "application"
}

// IsBusiness 检查是否为业务级别
func (l Level) IsBusiness() bool {
	return l.Category() == "business"
}

// IsSecurity 检查是否为安全级别
func (l Level) IsSecurity() bool {
	return l.Category() == "security"
}

// IsPerformance 检查是否为性能级别
func (l Level) IsPerformance() bool {
	return l.Category() == "performance"
}

// ParseLevel 解析级别字符串
func ParseLevel(s string) (Level, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	if level, ok := levelNameMap[s]; ok {
		return level, nil
	}
	return INFO, fmt.Errorf("invalid log level: %s", s)
}

// GetAllLevels 获取所有级别
func GetAllLevels() []Level {
	var levels []Level
	for level := range levelInfoMap {
		levels = append(levels, level)
	}
	return levels
}

// GetBasicLevels 获取基础级别
func GetBasicLevels() []Level {
	return []Level{TRACE, DEBUG, INFO, WARN, ERROR, FATAL, OFF}
}

// GetLevelsByCategory 根据类别获取级别
func GetLevelsByCategory(category string) []Level {
	return categoryLevelMap[category]
}

// GetAllCategories 获取所有类别
func GetAllCategories() []string {
	var categories []string
	for category := range categoryLevelMap {
		categories = append(categories, category)
	}
	return categories
}

// GetLevelNames 获取所有级别名称
func GetLevelNames() []string {
	var names []string
	for _, info := range levelInfoMap {
		names = append(names, info.Name)
	}
	return names
}

// GetLevelShortNames 获取所有级别短名称
func GetLevelShortNames() []string {
	var names []string
	for _, info := range levelInfoMap {
		names = append(names, info.ShortName)
	}
	return names
}

// ColorReset 颜色重置代码
const ColorReset = "\033[0m"

// 预定义颜色常量
const (
	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorGray    = "\033[90m"
	
	// 亮色版本
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightWhite   = "\033[97m"
)