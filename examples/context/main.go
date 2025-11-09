/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-09 18:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-09 18:30:00
 * @FilePath: \go-logger\examples\context\main.go
 * @Description: 上下文感知日志示例，演示如何在微服务中使用上下文进行分布式追踪
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/kamalyes/go-logger"
)

// 模拟的微服务组件
type UserService struct {
	logger logger.ILogger
}

type OrderService struct {
	logger     logger.ILogger
	userSvc    *UserService
	paymentSvc *PaymentService
}

type PaymentService struct {
	logger logger.ILogger
}

// 模拟的请求和响应结构
type CreateOrderRequest struct {
	UserID    int     `json:"user_id"`
	ProductID int     `json:"product_id"`
	Amount    float64 `json:"amount"`
}

type CreateOrderResponse struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

func main() {
		fmt.Println("=== 微服务分布式追踪演示结束 ===")

	// 创建日志器
	config := logger.DefaultConfig().
		WithLevel(logger.DEBUG).
		WithShowCaller(true).
		WithColorful(true).
		WithPrefix("[Context-Demo] ")

	baseLogger := logger.NewLogger(config)

	// 创建服务实例
	userSvc := &UserService{logger: baseLogger.WithField("service", "user")}
	paymentSvc := &PaymentService{logger: baseLogger.WithField("service", "payment")}
	orderSvc := &OrderService{
		logger:     baseLogger.WithField("service", "order"),
		userSvc:    userSvc,
		paymentSvc: paymentSvc,
	}

	// 模拟多个并发请求
	fmt.Println("1. === 模拟HTTP请求处理 ===")
	simulateHTTPRequests(orderSvc)

	// 模拟长时间运行的任务
	fmt.Println("\n2. === 模拟后台任务处理 ===")
	simulateBackgroundTask(baseLogger)

	// 演示错误传播和追踪
	fmt.Println("\n3. === 模拟错误传播追踪 ===")
	simulateErrorPropagation(orderSvc)

	fmt.Println("\n=== 上下文演示完成 ===")
}

// simulateHTTPRequests 模拟HTTP请求处理
func simulateHTTPRequests(orderSvc *OrderService) {
	requests := []CreateOrderRequest{
		{UserID: 1001, ProductID: 2001, Amount: 99.99},
		{UserID: 1002, ProductID: 2002, Amount: 149.99},
		{UserID: 1003, ProductID: 2003, Amount: 299.99},
	}

	for i, req := range requests {
		// 为每个请求创建独立的上下文
		ctx := createRequestContext(fmt.Sprintf("req-%d", i+1), fmt.Sprintf("user-%d", req.UserID))
		
		// 处理请求
		response := orderSvc.CreateOrder(ctx, req)
		
		fmt.Printf("请求处理结果: %+v\n", response)
		
		// 模拟请求间隔
		time.Sleep(time.Millisecond * 100)
	}
}

// simulateBackgroundTask 模拟后台任务
func simulateBackgroundTask(baseLogger logger.ILogger) {
	taskID := "task-cleanup-001"
	ctx := context.Background()
	ctx = context.WithValue(ctx, "task_id", taskID)
	ctx = context.WithValue(ctx, "task_type", "cleanup")
	
	// 创建任务专用的日志器
	taskLogger := baseLogger.WithContext(ctx).WithField("component", "background-task")
	
	taskLogger.InfoKV("后台任务开始",
		"task_id", taskID,
		"task_type", "cleanup",
		"scheduled_time", time.Now().Format(time.RFC3339),
	)
	
	// 模拟任务处理步骤
	steps := []string{"扫描过期数据", "备份重要数据", "删除过期记录", "更新统计信息", "清理临时文件"}
	
	for i, step := range steps {
		stepCtx := context.WithValue(ctx, "step", i+1)
		stepLogger := taskLogger.WithContext(stepCtx)
		
		stepLogger.DebugKV("执行清理步骤",
			"step_number", i+1,
			"step_name", step,
			"total_steps", len(steps),
		)
		
		// 模拟处理时间
		processingTime := time.Duration(rand.Intn(200)+50) * time.Millisecond
		time.Sleep(processingTime)
		
		stepLogger.InfoKV("清理步骤完成",
			"step_number", i+1,
			"step_name", step,
			"duration_ms", processingTime.Milliseconds(),
		)
	}
	
	taskLogger.InfoMsg("后台任务完成")
}

// simulateErrorPropagation 模拟错误传播和追踪
func simulateErrorPropagation(orderSvc *OrderService) {
	// 创建一个会导致错误的请求
	ctx := createRequestContext("req-error-001", "user-9999")
	
	req := CreateOrderRequest{
		UserID:    9999, // 不存在的用户
		ProductID: 2001,
		Amount:    99.99,
	}
	
	response := orderSvc.CreateOrder(ctx, req)
	fmt.Printf("错误请求处理结果: %+v\n", response)
}

// createRequestContext 创建请求上下文
func createRequestContext(requestID, userID string) context.Context {
	ctx := context.Background()
	
	// 添加分布式追踪信息
	traceID := fmt.Sprintf("trace-%d", time.Now().UnixNano())
	spanID := fmt.Sprintf("span-%d", time.Now().UnixNano()%1000000)
	
	ctx = context.WithValue(ctx, "request_id", requestID)
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "trace_id", traceID)
	ctx = context.WithValue(ctx, "span_id", spanID)
	ctx = context.WithValue(ctx, "start_time", time.Now())
	ctx = context.WithValue(ctx, "source_ip", "192.168.1.100")
	ctx = context.WithValue(ctx, "user_agent", "Go-HTTP-Client/1.1")
	
	return ctx
}

// CreateOrder 创建订单的业务逻辑
func (os *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) CreateOrderResponse {
	// 创建订单服务的日志器，携带请求上下文
	orderLogger := os.logger.WithContext(ctx)
	
	orderLogger.InfoKV("开始处理订单创建请求",
		"user_id", req.UserID,
		"product_id", req.ProductID,
		"amount", req.Amount,
	)
	
	// 1. 验证用户
	user, err := os.userSvc.ValidateUser(ctx, req.UserID)
	if err != nil {
		orderLogger.ErrorKV("用户验证失败",
			"user_id", req.UserID,
			"error", err.Error(),
		)
		return CreateOrderResponse{
			Status:  "failed",
			Message: "用户验证失败",
		}
	}
	
	orderLogger.DebugKV("用户验证成功",
		"user_id", req.UserID,
		"username", user["username"],
	)
	
	// 2. 处理支付
	paymentResult, err := os.paymentSvc.ProcessPayment(ctx, req.UserID, req.Amount)
	if err != nil {
		orderLogger.ErrorKV("支付处理失败",
			"user_id", req.UserID,
			"amount", req.Amount,
			"error", err.Error(),
		)
		return CreateOrderResponse{
			Status:  "failed",
			Message: "支付处理失败",
		}
	}
	
	// 3. 创建订单
	orderID := fmt.Sprintf("order-%d", time.Now().UnixNano())
	orderLogger.InfoKV("订单创建成功",
		"order_id", orderID,
		"user_id", req.UserID,
		"amount", req.Amount,
		"payment_id", paymentResult["payment_id"],
	)
	
	return CreateOrderResponse{
		OrderID: orderID,
		Status:  "success",
		Message: "订单创建成功",
	}
}

// ValidateUser 验证用户
func (us *UserService) ValidateUser(ctx context.Context, userID int) (map[string]interface{}, error) {
	userLogger := us.logger.WithContext(ctx)
	
	userLogger.DebugKV("开始用户验证",
		"user_id", userID,
		"validation_type", "identity_check",
	)
	
	// 模拟数据库查询
	time.Sleep(time.Millisecond * 50)
	
	// 模拟用户不存在的情况
	if userID == 9999 {
		userLogger.WarnKV("用户不存在",
			"user_id", userID,
			"check_result", "not_found",
		)
		return nil, fmt.Errorf("user %d not found", userID)
	}
	
	// 模拟成功的用户验证
	user := map[string]interface{}{
		"user_id":  userID,
		"username": fmt.Sprintf("user%d", userID),
		"status":   "active",
		"level":    "premium",
	}
	
	userLogger.InfoKV("用户验证通过",
		"user_id", userID,
		"username", user["username"],
		"status", user["status"],
		"level", user["level"],
	)
	
	return user, nil
}

// ProcessPayment 处理支付
func (ps *PaymentService) ProcessPayment(ctx context.Context, userID int, amount float64) (map[string]interface{}, error) {
	paymentLogger := ps.logger.WithContext(ctx)
	
	paymentLogger.InfoKV("开始处理支付",
		"user_id", userID,
		"amount", amount,
		"currency", "USD",
		"payment_method", "credit_card",
	)
	
	// 模拟支付处理时间
	processingTime := time.Duration(rand.Intn(200)+100) * time.Millisecond
	time.Sleep(processingTime)
	
	// 模拟支付ID生成
	paymentID := fmt.Sprintf("pay-%d", time.Now().UnixNano())
	
	paymentLogger.InfoKV("支付处理完成",
		"user_id", userID,
		"amount", amount,
		"payment_id", paymentID,
		"processing_time_ms", processingTime.Milliseconds(),
		"status", "completed",
	)
	
	return map[string]interface{}{
		"payment_id": paymentID,
		"status":     "completed",
		"amount":     amount,
		"currency":   "USD",
	}, nil
}

// 模拟HTTP中间件，展示如何在HTTP处理器中使用上下文日志
func LoggingMiddleware(logger logger.ILogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 为每个HTTP请求创建上下文
			ctx := r.Context()
			requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
			ctx = context.WithValue(ctx, "request_id", requestID)
			ctx = context.WithValue(ctx, "method", r.Method)
			ctx = context.WithValue(ctx, "path", r.URL.Path)
			ctx = context.WithValue(ctx, "remote_addr", r.RemoteAddr)
			ctx = context.WithValue(ctx, "user_agent", r.UserAgent())
			
			// 创建带上下文的日志器
			reqLogger := logger.WithContext(ctx)
			
			// 记录请求开始
			start := time.Now()
			reqLogger.InfoKV("HTTP请求开始",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
			
			// 调用下一个处理器
			next.ServeHTTP(w, r.WithContext(ctx))
			
			// 记录请求完成
			duration := time.Since(start)
			reqLogger.InfoKV("HTTP请求完成",
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", duration.Milliseconds(),
				"status", "200", // 简化示例，实际应该从ResponseWriter获取
			)
		})
	}
}