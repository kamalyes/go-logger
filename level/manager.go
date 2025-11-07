/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 00:00:00
 * @FilePath: \go-logger\level\manager.go
 * @Description: 级别管理器
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package level

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// LevelManager 级别管理器
type LevelManager struct {
	// 全局级别设置
	globalLevel Level
	
	// 组件级别映射
	componentLevels map[string]Level
	
	// 包级别映射
	packageLevels map[string]Level
	
	// 模式级别映射 (支持正则表达式)
	patternLevels map[*regexp.Regexp]Level
	
	// 动态级别控制
	enableDynamicLevel bool
	dynamicLevelTTL    time.Duration
	dynamicLevels      map[string]dynamicLevelEntry
	
	// 级别统计
	levelStats map[Level]*LevelStats
	
	// 配置
	enableStats    bool
	enableCaching  bool
	cacheSize      int
	levelCache     map[string]levelCacheEntry
	
	mu sync.RWMutex
}

// dynamicLevelEntry 动态级别条目
type dynamicLevelEntry struct {
	Level     Level
	ExpiresAt time.Time
}

// levelCacheEntry 级别缓存条目
type levelCacheEntry struct {
	Level     Level
	UpdatedAt time.Time
}

// LevelStats 级别统计
type LevelStats struct {
	Level       Level     `json:"level"`
	Count       uint64    `json:"count"`
	LastUsed    time.Time `json:"last_used"`
	FirstUsed   time.Time `json:"first_used"`
	TotalBytes  uint64    `json:"total_bytes"`
	AvgBytes    float64   `json:"avg_bytes"`
}

// NewLevelManager 创建级别管理器
func NewLevelManager() *LevelManager {
	lm := &LevelManager{
		globalLevel:        INFO,
		componentLevels:    make(map[string]Level),
		packageLevels:      make(map[string]Level),
		patternLevels:      make(map[*regexp.Regexp]Level),
		enableDynamicLevel: false,
		dynamicLevelTTL:    time.Hour,
		dynamicLevels:      make(map[string]dynamicLevelEntry),
		levelStats:         make(map[Level]*LevelStats),
		enableStats:        true,
		enableCaching:      true,
		cacheSize:          1000,
		levelCache:         make(map[string]levelCacheEntry),
	}
	
	// 初始化级别统计
	for _, level := range GetAllLevels() {
		lm.levelStats[level] = &LevelStats{
			Level: level,
		}
	}
	
	return lm
}

// SetGlobalLevel 设置全局级别
func (lm *LevelManager) SetGlobalLevel(level Level) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	lm.globalLevel = level
	lm.clearCache()
	return lm
}

// GetGlobalLevel 获取全局级别
func (lm *LevelManager) GetGlobalLevel() Level {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	return lm.globalLevel
}

// SetComponentLevel 设置组件级别
func (lm *LevelManager) SetComponentLevel(component string, level Level) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	lm.componentLevels[component] = level
	lm.clearCache()
	return lm
}

// GetComponentLevel 获取组件级别
func (lm *LevelManager) GetComponentLevel(component string) Level {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	if level, ok := lm.componentLevels[component]; ok {
		return level
	}
	return lm.globalLevel
}

// SetPackageLevel 设置包级别
func (lm *LevelManager) SetPackageLevel(pkg string, level Level) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	lm.packageLevels[pkg] = level
	lm.clearCache()
	return lm
}

// GetPackageLevel 获取包级别
func (lm *LevelManager) GetPackageLevel(pkg string) Level {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	// 检查缓存
	if lm.enableCaching {
		if entry, ok := lm.levelCache[pkg]; ok {
			if time.Since(entry.UpdatedAt) < time.Minute {
				return entry.Level
			}
		}
	}
	
	level := lm.resolvePackageLevel(pkg)
	
	// 更新缓存
	if lm.enableCaching && len(lm.levelCache) < lm.cacheSize {
		lm.levelCache[pkg] = levelCacheEntry{
			Level:     level,
			UpdatedAt: time.Now(),
		}
	}
	
	return level
}

// resolvePackageLevel 解析包级别
func (lm *LevelManager) resolvePackageLevel(pkg string) Level {
	// 精确匹配
	if level, ok := lm.packageLevels[pkg]; ok {
		return level
	}
	
	// 前缀匹配
	maxLength := 0
	var matchedLevel Level
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

// SetPatternLevel 设置模式级别
func (lm *LevelManager) SetPatternLevel(pattern string, level Level) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %v", err)
	}
	
	lm.patternLevels[regex] = level
	lm.clearCache()
	return nil
}

// SetDynamicLevel 设置动态级别
func (lm *LevelManager) SetDynamicLevel(key string, level Level, ttl time.Duration) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	if !lm.enableDynamicLevel {
		return lm
	}
	
	lm.dynamicLevels[key] = dynamicLevelEntry{
		Level:     level,
		ExpiresAt: time.Now().Add(ttl),
	}
	
	lm.clearCache()
	return lm
}

// GetDynamicLevel 获取动态级别
func (lm *LevelManager) GetDynamicLevel(key string) (Level, bool) {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	if !lm.enableDynamicLevel {
		return INFO, false
	}
	
	entry, ok := lm.dynamicLevels[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return INFO, false
	}
	
	return entry.Level, true
}

// EnableDynamicLevel 启用动态级别
func (lm *LevelManager) EnableDynamicLevel(enable bool, ttl time.Duration) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	lm.enableDynamicLevel = enable
	if ttl > 0 {
		lm.dynamicLevelTTL = ttl
	}
	return lm
}

// IsLevelEnabled 检查级别是否启用
func (lm *LevelManager) IsLevelEnabled(level Level, component, pkg string) bool {
	// 获取有效级别
	effectiveLevel := lm.GetEffectiveLevel(component, pkg)
	
	return level.IsEnabled(effectiveLevel)
}

// GetEffectiveLevel 获取有效级别
func (lm *LevelManager) GetEffectiveLevel(component, pkg string) Level {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	// 优先级：动态级别 > 组件级别 > 包级别 > 全局级别
	
	// 检查动态级别
	if lm.enableDynamicLevel {
		// 尝试组件动态级别
		if level, ok := lm.getDynamicLevel(fmt.Sprintf("component:%s", component)); ok {
			return level
		}
		
		// 尝试包动态级别
		if level, ok := lm.getDynamicLevel(fmt.Sprintf("package:%s", pkg)); ok {
			return level
		}
	}
	
	// 检查组件级别
	if level, ok := lm.componentLevels[component]; ok {
		return level
	}
	
	// 检查包级别
	pkgLevel := lm.resolvePackageLevel(pkg)
	if pkgLevel != lm.globalLevel {
		return pkgLevel
	}
	
	return lm.globalLevel
}

// getDynamicLevel 获取动态级别（内部方法，无锁）
func (lm *LevelManager) getDynamicLevel(key string) (Level, bool) {
	entry, ok := lm.dynamicLevels[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return INFO, false
	}
	return entry.Level, true
}

// RecordLevelUsage 记录级别使用
func (lm *LevelManager) RecordLevelUsage(level Level, messageSize uint64) {
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
	stats.TotalBytes += messageSize
	if stats.Count > 0 {
		stats.AvgBytes = float64(stats.TotalBytes) / float64(stats.Count)
	}
}

// GetLevelStats 获取级别统计
func (lm *LevelManager) GetLevelStats(level Level) *LevelStats {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	if stats, ok := lm.levelStats[level]; ok {
		// 返回副本
		return &LevelStats{
			Level:      stats.Level,
			Count:      stats.Count,
			LastUsed:   stats.LastUsed,
			FirstUsed:  stats.FirstUsed,
			TotalBytes: stats.TotalBytes,
			AvgBytes:   stats.AvgBytes,
		}
	}
	return nil
}

// GetAllLevelStats 获取所有级别统计
func (lm *LevelManager) GetAllLevelStats() map[Level]*LevelStats {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	result := make(map[Level]*LevelStats)
	for level, stats := range lm.levelStats {
		result[level] = &LevelStats{
			Level:      stats.Level,
			Count:      stats.Count,
			LastUsed:   stats.LastUsed,
			FirstUsed:  stats.FirstUsed,
			TotalBytes: stats.TotalBytes,
			AvgBytes:   stats.AvgBytes,
		}
	}
	return result
}

// GetTopLevels 获取使用最多的级别
func (lm *LevelManager) GetTopLevels(limit int) []*LevelStats {
	allStats := lm.GetAllLevelStats()
	
	var stats []*LevelStats
	for _, stat := range allStats {
		stats = append(stats, stat)
	}
	
	// 按使用次数排序
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})
	
	if limit > 0 && len(stats) > limit {
		stats = stats[:limit]
	}
	
	return stats
}

// EnableStats 启用统计
func (lm *LevelManager) EnableStats(enable bool) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	lm.enableStats = enable
	return lm
}

// EnableCaching 启用缓存
func (lm *LevelManager) EnableCaching(enable bool, size int) *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	lm.enableCaching = enable
	if size > 0 {
		lm.cacheSize = size
	}
	if !enable || size != lm.cacheSize {
		lm.clearCache()
	}
	return lm
}

// clearCache 清理缓存（内部方法，无锁）
func (lm *LevelManager) clearCache() {
	lm.levelCache = make(map[string]levelCacheEntry)
}

// ClearCache 清理缓存
func (lm *LevelManager) ClearCache() *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	lm.clearCache()
	return lm
}

// CleanupExpiredDynamicLevels 清理过期的动态级别
func (lm *LevelManager) CleanupExpiredDynamicLevels() {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	now := time.Now()
	for key, entry := range lm.dynamicLevels {
		if now.After(entry.ExpiresAt) {
			delete(lm.dynamicLevels, key)
		}
	}
}

// Reset 重置管理器
func (lm *LevelManager) Reset() *LevelManager {
	lm.mu.Lock()
	defer lm.mu.Unlock()
	
	lm.globalLevel = INFO
	lm.componentLevels = make(map[string]Level)
	lm.packageLevels = make(map[string]Level)
	lm.patternLevels = make(map[*regexp.Regexp]Level)
	lm.dynamicLevels = make(map[string]dynamicLevelEntry)
	lm.clearCache()
	
	return lm
}

// String 字符串表示
func (lm *LevelManager) String() string {
	lm.mu.RLock()
	defer lm.mu.RUnlock()
	
	return fmt.Sprintf("LevelManager{Global: %s, Components: %d, Packages: %d, Patterns: %d, Dynamic: %v}",
		lm.globalLevel, len(lm.componentLevels), len(lm.packageLevels), len(lm.patternLevels), lm.enableDynamicLevel)
}