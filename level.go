/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 12:33:32
 * @FilePath: \go-logger\level.go
 * @Description: 统一的日志级别定义和管理
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package logger

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// LogLevel 日志级别类型 (保持向后兼容)
type LogLevel int

// 基础日志级别常量 (保持现有API兼容性)
const (
	DEBUG LogLevel = iota // 调试级别 - 最详细的信息
	INFO                  // 信息级别 - 一般信息
	WARN                  // 警告级别 - 警告信息
	ERROR                 // 错误级别 - 错误信息
	FATAL                 // 致命级别 - 致命错误，程序将退出
)

// 扩展级别常量 (新增的高级功能)
const (
	TRACE LogLevel = -1 // 跟踪级别 - 更详细的调试信息
	OFF   LogLevel = 99 // 关闭级别 - 禁用所有日志
)

// 系统级扩展级别
const (
	SYSTEM      LogLevel = 100 + iota // 系统级信息
	KERNEL                            // 内核级信息
	DRIVER                            // 驱动级信息
	APPLICATION                       // 应用级信息
	SERVICE                           // 服务级信息
	COMPONENT                         // 组件级信息
	MODULE                            // 模块级信息
)

// 业务级扩展级别
const (
	BUSINESS    LogLevel = 200 + iota // 业务级信息
	TRANSACTION                       // 事务级信息
	WORKFLOW                          // 工作流信息
	PROCESS                           // 流程级信息
)

// 安全级扩展级别
const (
	SECURITY   LogLevel = 300 + iota // 安全级信息
	AUDIT                            // 审计级信息
	COMPLIANCE                       // 合规级信息
	THREAT                           // 威胁级信息
)

// 性能级扩展级别
const (
	PERFORMANCE LogLevel = 400 + iota // 性能级信息
	METRIC                            // 指标级信息
	BENCHMARK                         // 基准测试信息
	PROFILING                         // 性能分析信息
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
var levelInfoMap = map[LogLevel]LevelInfo{
	// 基础级别
	TRACE: {"TRACE", "TRC", "🔍", "\033[90m", -1, "详细跟踪信息", "basic"},
	DEBUG: {"DEBUG", "DBG", "🐛", "\033[36m", 0, "调试信息", "basic"},
	INFO:  {"INFO", "INF", "ℹ️", "\033[32m", 1, "一般信息", "basic"},
	WARN:  {"WARN", "WRN", "⚠️", "\033[33m", 2, "警告信息", "basic"},
	ERROR: {"ERROR", "ERR", "❌", "\033[31m", 3, "错误信息", "basic"},
	FATAL: {"FATAL", "FTL", "💀", "\033[35m", 4, "致命错误", "basic"},
	OFF:   {"OFF", "OFF", "🚫", "\033[0m", 99, "关闭日志", "basic"},

	// 系统级别
	SYSTEM:      {"SYSTEM", "SYS", "🖥️", "\033[94m", 100, "系统级信息", "system"},
	KERNEL:      {"KERNEL", "KRN", "⚙️", "\033[95m", 101, "内核级信息", "system"},
	DRIVER:      {"DRIVER", "DRV", "🔌", "\033[96m", 102, "驱动级信息", "system"},
	APPLICATION: {"APPLICATION", "APP", "📱", "\033[92m", 103, "应用级信息", "application"},
	SERVICE:     {"SERVICE", "SVC", "🔧", "\033[93m", 104, "服务级信息", "application"},
	COMPONENT:   {"COMPONENT", "CMP", "🧩", "\033[94m", 105, "组件级信息", "application"},
	MODULE:      {"MODULE", "MOD", "📦", "\033[95m", 106, "模块级信息", "application"},

	// 业务级别
	BUSINESS:    {"BUSINESS", "BIZ", "💼", "\033[38;5;214m", 200, "业务级信息", "business"},
	TRANSACTION: {"TRANSACTION", "TXN", "💳", "\033[38;5;215m", 201, "事务级信息", "business"},
	WORKFLOW:    {"WORKFLOW", "WFL", "🔄", "\033[38;5;216m", 202, "工作流信息", "business"},
	PROCESS:     {"PROCESS", "PRC", "⚡", "\033[38;5;217m", 203, "流程级信息", "business"},

	// 安全级别
	SECURITY:   {"SECURITY", "SEC", "🔒", "\033[38;5;196m", 300, "安全级信息", "security"},
	AUDIT:      {"AUDIT", "ADT", "📋", "\033[38;5;197m", 301, "审计级信息", "security"},
	COMPLIANCE: {"COMPLIANCE", "CMP", "✅", "\033[38;5;198m", 302, "合规级信息", "security"},
	THREAT:     {"THREAT", "THR", "🛡️", "\033[38;5;199m", 303, "威胁级信息", "security"},

	// 性能级别
	PERFORMANCE: {"PERFORMANCE", "PRF", "📊", "\033[38;5;81m", 400, "性能级信息", "performance"},
	METRIC:      {"METRIC", "MTC", "📈", "\033[38;5;82m", 401, "指标级信息", "performance"},
	BENCHMARK:   {"BENCHMARK", "BMK", "⏱️", "\033[38;5;83m", 402, "基准测试信息", "performance"},
	PROFILING:   {"PROFILING", "PRO", "🔬", "\033[38;5;84m", 403, "性能分析信息", "performance"},
}

// 级别名称映射
var levelNameMap = func() map[string]LogLevel {
	m := make(map[string]LogLevel)
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

// =============================================================================
// LogLevel 基础方法 (保持向后兼容)
// =============================================================================

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Name
	}
	return fmt.Sprintf("UNKNOWN(%d)", int(l))
}

// UnmarshalYAML 实现 YAML 反序列化，支持字符串和整数值
func (l *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// 优先尝试字符串解析
	var s string
	if err := unmarshal(&s); err == nil {
		level, parseErr := ParseLevel(s)
		if parseErr != nil {
			return parseErr
		}
		*l = level
		return nil
	}

	// 降级为整数解析（向后兼容）
	var i int
	if err := unmarshal(&i); err != nil {
		return err
	}
	*l = LogLevel(i)
	return nil
}

// UnmarshalText 实现文本反序列化（用于 JSON 等）
func (l *LogLevel) UnmarshalText(text []byte) error {
	level, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*l = level
	return nil
}

// MarshalText 实现文本序列化（用于 JSON 等）
func (l LogLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// ShortString 返回级别的短字符串表示
func (l LogLevel) ShortString() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.ShortName
	}
	return "UNK"
}

// Emoji 返回日志级别的emoji
func (l LogLevel) Emoji() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Emoji
	}
	return "❓"
}

// Color 返回日志级别的颜色代码
func (l LogLevel) Color() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Color
	}
	return "\033[0m" // 重置颜色
}

// Priority 返回级别的优先级
func (l LogLevel) Priority() int {
	if info, ok := levelInfoMap[l]; ok {
		return info.Priority
	}
	return int(l)
}

// Description 返回级别的描述
func (l LogLevel) Description() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Description
	}
	return "Unknown level"
}

// Category 返回级别的类别
func (l LogLevel) Category() string {
	if info, ok := levelInfoMap[l]; ok {
		return info.Category
	}
	return "unknown"
}

// Info 返回级别的完整信息
func (l LogLevel) Info() LevelInfo {
	if info, ok := levelInfoMap[l]; ok {
		return info
	}
	return LevelInfo{
		Name:        fmt.Sprintf("UNKNOWN(%d)", int(l)),
		ShortName:   "UNK",
		Emoji:       "❓",
		Color:       "\033[0m",
		Priority:    int(l),
		Description: "Unknown level",
		Category:    "unknown",
	}
}

// IsValid 检查级别是否有效
func (l LogLevel) IsValid() bool {
	_, ok := levelInfoMap[l]
	return ok
}

// IsEnabled 检查当前级别是否启用目标级别
func (l LogLevel) IsEnabled(target LogLevel) bool {
	return target.Priority() >= l.Priority()
}

// IsBasic 检查是否为基础级别
func (l LogLevel) IsBasic() bool {
	return l.Category() == "basic"
}

// IsSystem 检查是否为系统级别
func (l LogLevel) IsSystem() bool {
	return l.Category() == "system"
}

// IsApplication 检查是否为应用级别
func (l LogLevel) IsApplication() bool {
	return l.Category() == "application"
}

// IsBusiness 检查是否为业务级别
func (l LogLevel) IsBusiness() bool {
	return l.Category() == "business"
}

// IsSecurity 检查是否为安全级别
func (l LogLevel) IsSecurity() bool {
	return l.Category() == "security"
}

// IsPerformance 检查是否为性能级别
func (l LogLevel) IsPerformance() bool {
	return l.Category() == "performance"
}

// =============================================================================
// 全局级别函数
// =============================================================================

// ParseLevel 从字符串解析日志级别（支持大小写）
func ParseLevel(level string) (LogLevel, error) {
	level = strings.ToUpper(strings.TrimSpace(level))
	if l, ok := levelNameMap[level]; ok {
		return l, nil
	}
	return DEBUG, fmt.Errorf("invalid log level: %s", level)
}

// GetAllLevels 获取所有基础级别 (保持向后兼容)
func GetAllLevels() []LogLevel {
	return []LogLevel{DEBUG, INFO, WARN, ERROR, FATAL}
}

// GetAllExtendedLevels 获取所有级别（包括扩展级别）
func GetAllExtendedLevels() []LogLevel {
	var levels []LogLevel
	for level := range levelInfoMap {
		levels = append(levels, level)
	}
	// 按优先级排序
	sort.Slice(levels, func(i, j int) bool {
		return levels[i].Priority() < levels[j].Priority()
	})
	return levels
}

// GetExtendedLevels 获取扩展级别（不包括基础级别）
func GetExtendedLevels() []LogLevel {
	var levels []LogLevel
	basicLevels := map[LogLevel]bool{TRACE: true, DEBUG: true, INFO: true, WARN: true, ERROR: true, FATAL: true, OFF: true}

	for level := range levelInfoMap {
		if !basicLevels[level] {
			levels = append(levels, level)
		}
	}
	// 按优先级排序
	sort.Slice(levels, func(i, j int) bool {
		return levels[i].Priority() < levels[j].Priority()
	})
	return levels
}

// GetBasicLevels 获取基础级别
func GetBasicLevels() []LogLevel {
	return []LogLevel{TRACE, DEBUG, INFO, WARN, ERROR, FATAL, OFF}
}

// GetLevelsByCategory 根据类别获取级别
func GetLevelsByCategory(category string) []LogLevel {
	var levels []LogLevel
	for level, info := range levelInfoMap {
		if info.Category == category {
			levels = append(levels, level)
		}
	}
	// 按优先级排序
	sort.Slice(levels, func(i, j int) bool {
		return levels[i].Priority() < levels[j].Priority()
	})
	return levels
}

// GetAllCategories 获取所有类别
func GetAllCategories() []string {
	categoryMap := make(map[string]bool)
	for _, info := range levelInfoMap {
		categoryMap[info.Category] = true
	}

	var categories []string
	for category := range categoryMap {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// GetLevelNames 获取基础级别名称 (保持向后兼容)
func GetLevelNames() []string {
	levels := GetAllLevels()
	names := make([]string, len(levels))
	for i, level := range levels {
		names[i] = level.String()
	}
	return names
}

// GetAllLevelNames 获取所有级别名称（包括扩展级别）
func GetAllLevelNames() []string {
	var names []string
	for _, info := range levelInfoMap {
		names = append(names, info.Name)
	}
	sort.Strings(names)
	return names
}

// GetLevelShortNames 获取所有级别短名称
func GetLevelShortNames() []string {
	var names []string
	for _, info := range levelInfoMap {
		names = append(names, info.ShortName)
	}
	sort.Strings(names)
	return names
}

// =============================================================================
// 颜色常量
// =============================================================================

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

// =============================================================================
// 全局级别管理器 (自动注入)
// =============================================================================

// LevelManager 级别管理器
type LevelManager struct {
	globalLevel     LogLevel
	componentLevels map[string]LogLevel
	packageLevels   map[string]LogLevel
	patternLevels   map[*regexp.Regexp]LogLevel
	enableStats     bool
	levelStats      map[LogLevel]*LevelStats
	mu              sync.RWMutex
}

// LevelStats 级别统计
type LevelStats struct {
	Level     LogLevel  `json:"level"`
	Count     uint64    `json:"count"`
	LastUsed  time.Time `json:"last_used"`
	FirstUsed time.Time `json:"first_used"`
}

// 全局级别管理器实例 (自动注入)
var (
	globalLevelManager *LevelManager
	managerOnce        sync.Once
)

// GetLevelManager 获取全局级别管理器实例 (单例模式)
func GetLevelManager() *LevelManager {
	managerOnce.Do(func() {
		globalLevelManager = &LevelManager{
			globalLevel:     INFO,
			componentLevels: make(map[string]LogLevel),
			packageLevels:   make(map[string]LogLevel),
			patternLevels:   make(map[*regexp.Regexp]LogLevel),
			enableStats:     true,
			levelStats:      make(map[LogLevel]*LevelStats),
		}

		// 初始化级别统计
		for level := range levelInfoMap {
			globalLevelManager.levelStats[level] = &LevelStats{Level: level}
		}
	})

	return globalLevelManager
}

// NewLevelManager 创建新的级别管理器实例 (可选)
func NewLevelManager() *LevelManager {
	lm := &LevelManager{
		globalLevel:     INFO,
		componentLevels: make(map[string]LogLevel),
		packageLevels:   make(map[string]LogLevel),
		patternLevels:   make(map[*regexp.Regexp]LogLevel),
		enableStats:     true,
		levelStats:      make(map[LogLevel]*LevelStats),
	}

	// 初始化级别统计
	for level := range levelInfoMap {
		lm.levelStats[level] = &LevelStats{Level: level}
	}

	return lm
}

// SetGlobalLevel 设置全局级别
func (lm *LevelManager) SetGlobalLevel(level LogLevel) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.globalLevel = level
	return lm
}

// GetGlobalLevel 获取全局级别
func (lm *LevelManager) GetGlobalLevel() LogLevel {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	return lm.globalLevel
}

// SetComponentLevel 设置组件级别
func (lm *LevelManager) SetComponentLevel(component string, level LogLevel) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.componentLevels[component] = level
	return lm
}

// GetComponentLevel 获取组件级别
func (lm *LevelManager) GetComponentLevel(component string) LogLevel {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if level, ok := lm.componentLevels[component]; ok {
		return level
	}
	return lm.globalLevel
}

// SetPackageLevel 设置包级别
func (lm *LevelManager) SetPackageLevel(pkg string, level LogLevel) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	lm.packageLevels[pkg] = level
	return lm
}

// GetPackageLevel 获取包级别
func (lm *LevelManager) GetPackageLevel(pkg string) LogLevel {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	// 精确匹配
	if level, ok := lm.packageLevels[pkg]; ok {
		return level
	}

	// 前缀匹配
	maxLength := 0
	var matchedLevel LogLevel
	found := false

	for prefix, level := range lm.packageLevels {
		if strings.HasPrefix(pkg, prefix) && len(prefix) > maxLength {
			maxLength = len(prefix)
			matchedLevel = level
			found = true
		}
	}

	if found {
		return matchedLevel
	}

	// 正则模式匹配
	for pattern, level := range lm.patternLevels {
		if pattern.MatchString(pkg) {
			return level
		}
	}

	return lm.globalLevel
}

// IsLevelEnabled 检查级别是否启用
func (lm *LevelManager) IsLevelEnabled(level LogLevel, component, pkg string) bool {
	effectiveLevel := lm.GetEffectiveLevel(component, pkg)
	return level.IsEnabled(effectiveLevel)
}

// GetEffectiveLevel 获取有效级别
func (lm *LevelManager) GetEffectiveLevel(component, pkg string) LogLevel {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	// 优先级：组件级别 > 包级别 > 全局级别

	// 检查组件级别
	if level, ok := lm.componentLevels[component]; ok {
		return level
	}

	// 检查包级别
	pkgLevel := lm.GetPackageLevel(pkg)
	if pkgLevel != lm.globalLevel {
		return pkgLevel
	}

	return lm.globalLevel
}

// RecordLevelUsage 记录级别使用
func (lm *LevelManager) RecordLevelUsage(level LogLevel) {
	if !lm.enableStats {
		return
	}

	lm.mu.Lock()
	defer lm.mu.Unlock()

	stats, ok := lm.levelStats[level]
	if !ok {
		stats = &LevelStats{Level: level}
		lm.levelStats[level] = stats
	}

	now := time.Now()
	stats.Count++
	stats.LastUsed = now
	if stats.FirstUsed.IsZero() {
		stats.FirstUsed = now
	}
}

// GetLevelStats 获取级别统计
func (lm *LevelManager) GetLevelStats(level LogLevel) *LevelStats {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	if stats, ok := lm.levelStats[level]; ok {
		// 返回副本
		return &LevelStats{
			Level:     stats.Level,
			Count:     stats.Count,
			LastUsed:  stats.LastUsed,
			FirstUsed: stats.FirstUsed,
		}
	}
	return nil
}

// String 字符串表示
func (lm *LevelManager) String() string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()

	return fmt.Sprintf("LevelManager{Global: %s, Components: %d, Packages: %d}",
		lm.globalLevel, len(lm.componentLevels), len(lm.packageLevels))
}

// =============================================================================
// 全局级别管理 (自动注入便捷接口)
// =============================================================================

// GlobalLevel 全局级别管理接口
var GlobalLevel = struct {
	// Set 设置全局级别
	Set func(level LogLevel)
	// Get 获取全局级别
	Get func() LogLevel
	// SetComponent 设置组件级别
	SetComponent func(component string, level LogLevel)
	// GetComponent 获取组件级别
	GetComponent func(component string) LogLevel
	// SetPackage 设置包级别
	SetPackage func(pkg string, level LogLevel)
	// GetPackage 获取包级别
	GetPackage func(pkg string) LogLevel
	// IsEnabled 检查级别是否启用
	IsEnabled func(level LogLevel, component, pkg string) bool
	// GetEffective 获取有效级别
	GetEffective func(component, pkg string) LogLevel
	// RecordUsage 记录级别使用
	RecordUsage func(level LogLevel)
	// GetStats 获取级别统计
	GetStats func(level LogLevel) *LevelStats
}{
	Set: func(level LogLevel) {
		GetLevelManager().SetGlobalLevel(level)
	},
	Get: func() LogLevel {
		return GetLevelManager().GetGlobalLevel()
	},
	SetComponent: func(component string, level LogLevel) {
		GetLevelManager().SetComponentLevel(component, level)
	},
	GetComponent: func(component string) LogLevel {
		return GetLevelManager().GetComponentLevel(component)
	},
	SetPackage: func(pkg string, level LogLevel) {
		GetLevelManager().SetPackageLevel(pkg, level)
	},
	GetPackage: func(pkg string) LogLevel {
		return GetLevelManager().GetPackageLevel(pkg)
	},
	IsEnabled: func(level LogLevel, component, pkg string) bool {
		return GetLevelManager().IsLevelEnabled(level, component, pkg)
	},
	GetEffective: func(component, pkg string) LogLevel {
		return GetLevelManager().GetEffectiveLevel(component, pkg)
	},
	RecordUsage: func(level LogLevel) {
		GetLevelManager().RecordLevelUsage(level)
	},
	GetStats: func(level LogLevel) *LevelStats {
		return GetLevelManager().GetLevelStats(level)
	},
}
