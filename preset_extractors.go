/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-24 12:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-24 23:29:14
 * @FilePath: \engine-im-service\go-logger\preset_extractors.go
 * @Description: 预设的上下文提取器 - 为常见场景提供开箱即用的提取器
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */

package logger

// PresetExtractorType 预设提取器类型
type PresetExtractorType string

const (
	PresetExtractorDefault   PresetExtractorType = "default"   // 默认：TraceID + RequestID
	PresetExtractorGateway   PresetExtractorType = "gateway"   // 网关：TraceID + RequestID + User + Tenant
	PresetExtractorService   PresetExtractorType = "service"   // 服务：TraceID + RequestID + User + Tenant
	PresetExtractorWebSocket PresetExtractorType = "websocket" // WebSocket：TraceID + User + Session
	PresetExtractorGRPC      PresetExtractorType = "grpc"      // gRPC：从 metadata 提取
	PresetExtractorFull      PresetExtractorType = "full"      // 完整：所有字段
)

// GetPresetExtractor 获取预设的上下文提取器
func GetPresetExtractor(presetType PresetExtractorType) ContextExtractor {
	switch presetType {
	case PresetExtractorDefault:
		return GetDefaultPresetExtractor()
	case PresetExtractorGateway:
		return GetGatewayPresetExtractor()
	case PresetExtractorService:
		return GetServicePresetExtractor()
	case PresetExtractorWebSocket:
		return GetWebSocketPresetExtractor()
	case PresetExtractorGRPC:
		return GetGRPCPresetExtractor()
	case PresetExtractorFull:
		return GetFullPresetExtractor()
	default:
		return GetDefaultPresetExtractor()
	}
}

// GetDefaultPresetExtractor 获取默认预设提取器
// 提取: TraceID, RequestID
func GetDefaultPresetExtractor() ContextExtractor {
	return NewContextExtractorBuilder().
		AddTraceID().
		AddRequestID().
		Build()
}

// GetGatewayPresetExtractor 获取网关预设提取器
// 提取: TraceID, RequestID, UserID, TenantID
func GetGatewayPresetExtractor() ContextExtractor {
	return NewContextExtractorBuilder().
		AddTraceID().
		AddRequestID().
		AddContextValue(string(KeyUserID), "User").
		AddContextValue(string(KeyTenantID), "Tenant").
		Build()
}

// GetServicePresetExtractor 获取服务预设提取器
// 提取: TraceID, RequestID, UserID, TenantID
func GetServicePresetExtractor() ContextExtractor {
	return NewContextExtractorBuilder().
		AddTraceID().
		AddRequestID().
		AddContextValue(string(KeyUserID), "User").
		AddContextValue(string(KeyTenantID), "Tenant").
		Build()
}

// GetWebSocketPresetExtractor 获取 WebSocket 预设提取器
// 提取: TraceID, UserID, SessionID
func GetWebSocketPresetExtractor() ContextExtractor {
	return NewContextExtractorBuilder().
		AddTraceID().
		AddContextValue(string(KeyUserID), "User").
		AddContextValue(string(KeySessionID), "Session").
		Build()
}

// GetGRPCPresetExtractor 获取 gRPC 预设提取器
// 从 gRPC metadata 提取: TraceID, RequestID, UserID, TenantID
func GetGRPCPresetExtractor() ContextExtractor {
	return CustomFieldExtractor(
		[]string{
			string(KeyTraceID),
			string(KeyRequestID),
			string(KeyUserID),
			string(KeyTenantID),
		},
		[]string{
			"x-trace-id",
			"x-request-id",
			"x-user-id",
			"x-tenant-id",
		},
	)
}

// GetFullPresetExtractor 获取完整预设提取器
// 提取所有标准字段
func GetFullPresetExtractor() ContextExtractor {
	return NewContextExtractorBuilder().
		AddTraceID().
		AddRequestID().
		AddContextValue(string(KeyUserID), "User").
		AddContextValue(string(KeyTenantID), "Tenant").
		AddContextValue(string(KeySessionID), "Session").
		AddContextValue(string(KeyOperation), "Op").
		Build()
}
