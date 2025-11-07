/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-07 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-07 23:35:55
 * @FilePath: \go-logger\context\correlation.go
 * @Description: 相关性管理和ID生成
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package context

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// IDGenerator ID生成器接口
type IDGenerator interface {
	GenerateTraceID() string
	GenerateSpanID() string
	GenerateRequestID() string
	GenerateCorrelationID() string
}

// DefaultIDGenerator 默认ID生成器
type DefaultIDGenerator struct {
	counter uint64
	prefix  string
	mu      sync.Mutex
}

// NewDefaultIDGenerator 创建默认ID生成器
func NewDefaultIDGenerator() *DefaultIDGenerator {
	return &DefaultIDGenerator{
		prefix: "logger",
	}
}

// SetPrefix 设置ID前缀
func (g *DefaultIDGenerator) SetPrefix(prefix string) *DefaultIDGenerator {
	g.mu.Lock()
	defer g.mu.Unlock()
	
	g.prefix = prefix
	return g
}

// GenerateTraceID 生成跟踪ID
func (g *DefaultIDGenerator) GenerateTraceID() string {
	// 使用时间戳 + 随机数的方式生成32字符的跟踪ID
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 8)
	_, _ = rand.Read(randomBytes)
	
	return fmt.Sprintf("%016x%s", timestamp, hex.EncodeToString(randomBytes))
}

// GenerateSpanID 生成跨度ID
func (g *DefaultIDGenerator) GenerateSpanID() string {
	// 生成16字符的跨度ID
	randomBytes := make([]byte, 8)
	_, _ = rand.Read(randomBytes)
	
	return hex.EncodeToString(randomBytes)
}

// GenerateRequestID 生成请求ID
func (g *DefaultIDGenerator) GenerateRequestID() string {
	// 生成带前缀的请求ID
	counter := atomic.AddUint64(&g.counter, 1)
	timestamp := time.Now().Unix()
	
	return fmt.Sprintf("%s-%d-%d", g.prefix, timestamp, counter)
}

// GenerateCorrelationID 生成关联ID
func (g *DefaultIDGenerator) GenerateCorrelationID() string {
	// 生成UUID风格的关联ID
	randomBytes := make([]byte, 16)
	_, _ = rand.Read(randomBytes)
	
	// 设置版本(4)和变体位
	randomBytes[6] = (randomBytes[6] & 0x0f) | 0x40 // Version 4
	randomBytes[8] = (randomBytes[8] & 0x3f) | 0x80 // Variant 10
	
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		randomBytes[0:4], randomBytes[4:6], randomBytes[6:8], randomBytes[8:10], randomBytes[10:16])
}

// CorrelationManager 相关性管理器
type CorrelationManager struct {
	generator    IDGenerator
	correlations map[string]*CorrelationChain
	mu           sync.RWMutex
	
	// 配置
	maxChainLength int
	maxChainAge    time.Duration
	cleanupInterval time.Duration
}

// CorrelationChain 相关性链
type CorrelationChain struct {
	ID        string                 `json:"id"`
	TraceID   string                 `json:"trace_id"`
	ParentID  string                 `json:"parent_id,omitempty"`
	Children  []string               `json:"children,omitempty"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time,omitempty"`
	Tags      map[string]string      `json:"tags,omitempty"`
	Metrics   map[string]interface{} `json:"metrics,omitempty"`
	
	mu sync.RWMutex
}

// NewCorrelationManager 创建相关性管理器
func NewCorrelationManager(generator IDGenerator) *CorrelationManager {
	if generator == nil {
		generator = NewDefaultIDGenerator()
	}
	
	cm := &CorrelationManager{
		generator:       generator,
		correlations:    make(map[string]*CorrelationChain),
		maxChainLength:  100,
		maxChainAge:     time.Hour,
		cleanupInterval: time.Minute * 10,
	}
	
	// 启动清理协程
	go cm.startCleanup()
	
	return cm
}

// CreateChain 创建相关性链
func (cm *CorrelationManager) CreateChain(traceID string) *CorrelationChain {
	chain := &CorrelationChain{
		ID:        cm.generator.GenerateCorrelationID(),
		TraceID:   traceID,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
		Metrics:   make(map[string]interface{}),
		Children:  make([]string, 0),
	}
	
	cm.mu.Lock()
	cm.correlations[chain.ID] = chain
	cm.mu.Unlock()
	
	return chain
}

// GetChain 获取相关性链
func (cm *CorrelationManager) GetChain(id string) (*CorrelationChain, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	chain, exists := cm.correlations[id]
	if !exists {
		return nil, false
	}
	
	return chain.Clone(), true
}

// AddChild 添加子链
func (cm *CorrelationManager) AddChild(parentID, childID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	parent, exists := cm.correlations[parentID]
	if !exists {
		return fmt.Errorf("parent correlation not found: %s", parentID)
	}
	
	child, exists := cm.correlations[childID]
	if !exists {
		return fmt.Errorf("child correlation not found: %s", childID)
	}
	
	// 设置父子关系
	parent.Children = append(parent.Children, childID)
	child.ParentID = parentID
	
	return nil
}

// EndChain 结束相关性链
func (cm *CorrelationManager) EndChain(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	chain, exists := cm.correlations[id]
	if !exists {
		return fmt.Errorf("correlation not found: %s", id)
	}
	
	chain.EndTime = time.Now()
	return nil
}

// GetChainsByTrace 根据跟踪ID获取相关性链
func (cm *CorrelationManager) GetChainsByTrace(traceID string) []*CorrelationChain {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	var chains []*CorrelationChain
	for _, chain := range cm.correlations {
		if chain.TraceID == traceID {
			chains = append(chains, chain.Clone())
		}
	}
	
	return chains
}

// GetActiveChains 获取活跃的相关性链
func (cm *CorrelationManager) GetActiveChains() []*CorrelationChain {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	var chains []*CorrelationChain
	for _, chain := range cm.correlations {
		if chain.EndTime.IsZero() {
			chains = append(chains, chain.Clone())
		}
	}
	
	return chains
}

// SetMaxChainLength 设置最大链长度
func (cm *CorrelationManager) SetMaxChainLength(length int) *CorrelationManager {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.maxChainLength = length
	return cm
}

// SetCleanupConfig 设置清理配置
func (cm *CorrelationManager) SetCleanupConfig(maxAge, interval time.Duration) *CorrelationManager {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.maxChainAge = maxAge
	cm.cleanupInterval = interval
	return cm
}

// startCleanup 启动清理协程
func (cm *CorrelationManager) startCleanup() {
	ticker := time.NewTicker(cm.cleanupInterval)
	defer ticker.Stop()
	
	for range ticker.C {
		cm.cleanup()
	}
}

// cleanup 清理过期的相关性链
func (cm *CorrelationManager) cleanup() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	now := time.Now()
	for id, chain := range cm.correlations {
		// 清理超时的链
		if now.Sub(chain.StartTime) > cm.maxChainAge {
			delete(cm.correlations, id)
			continue
		}
		
		// 清理过长的链
		if len(chain.Children) > cm.maxChainLength {
			// 只保留最新的子链
			chain.Children = chain.Children[len(chain.Children)-cm.maxChainLength:]
		}
	}
}

// Clone 克隆相关性链
func (c *CorrelationChain) Clone() *CorrelationChain {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	tags := make(map[string]string)
	for k, v := range c.Tags {
		tags[k] = v
	}
	
	metrics := make(map[string]interface{})
	for k, v := range c.Metrics {
		metrics[k] = v
	}
	
	children := make([]string, len(c.Children))
	copy(children, c.Children)
	
	return &CorrelationChain{
		ID:        c.ID,
		TraceID:   c.TraceID,
		ParentID:  c.ParentID,
		Children:  children,
		StartTime: c.StartTime,
		EndTime:   c.EndTime,
		Tags:      tags,
		Metrics:   metrics,
	}
}

// SetTag 设置标签
func (c *CorrelationChain) SetTag(key, value string) *CorrelationChain {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.Tags[key] = value
	return c
}

// SetMetric 设置指标
func (c *CorrelationChain) SetMetric(key string, value interface{}) *CorrelationChain {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.Metrics[key] = value
	return c
}

// GetDuration 获取持续时间
func (c *CorrelationChain) GetDuration() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if c.EndTime.IsZero() {
		return time.Since(c.StartTime)
	}
	return c.EndTime.Sub(c.StartTime)
}

// IsActive 检查是否活跃
func (c *CorrelationChain) IsActive() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.EndTime.IsZero()
}

// 全局ID生成器实例
var globalIDGenerator IDGenerator = NewDefaultIDGenerator()

// SetGlobalIDGenerator 设置全局ID生成器
func SetGlobalIDGenerator(generator IDGenerator) {
	globalIDGenerator = generator
}

// 全局ID生成函数
func generateTraceID() string {
	return globalIDGenerator.GenerateTraceID()
}

func generateSpanID() string {
	return globalIDGenerator.GenerateSpanID()
}

func generateRequestID() string {
	return globalIDGenerator.GenerateRequestID()
}

// String 字符串表示
func (c *CorrelationChain) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return fmt.Sprintf("CorrelationChain{ID: %s, TraceID: %s, ParentID: %s, Children: %d, Duration: %v}",
		c.ID, c.TraceID, c.ParentID, len(c.Children), c.GetDuration())
}