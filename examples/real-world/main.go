/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-22 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-22 12:15:22
 * @FilePath: \go-logger\examples\real-world\main.go
 * @Description: 实际应用场景示例 - 展示在真实项目中如何使用logger.New()
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package main

import (
	"context"
	"fmt"
	"github.com/kamalyes/go-logger"
	"math/rand"
	"strings"
	"time"
)

// 模拟用户服务
type UserService struct {
	logger *logger.Logger
}

func NewUserService() *UserService {
	return &UserService{
		logger: logger.New().WithPrefix("[UserService] ").WithLevel(logger.DEBUG),
	}
}

func (s *UserService) Login(userID string, password string) error {
	// 记录登录尝试
	s.logger.WithFields(map[string]interface{}{
		"user_id":   userID,
		"action":    "login_attempt",
		"timestamp": time.Now(),
	}).Info("用户尝试登录")

	// 模拟验证过程
	if password == "wrong" {
		s.logger.WithField("user_id", userID).Error("登录失败：密码错误")
		return fmt.Errorf("密码错误")
	}

	// 登录成功
	s.logger.WithFields(map[string]interface{}{
		"user_id":    userID,
		"action":     "login_success",
		"session_id": generateSessionID(),
	}).Info("用户登录成功")

	return nil
}

func (s *UserService) GetProfile(userID string) (map[string]interface{}, error) {
	requestLogger := s.logger.WithField("user_id", userID).WithField("action", "get_profile")

	requestLogger.Debug("开始获取用户资料")

	// 模拟数据库查询
	time.Sleep(50 * time.Millisecond)

	profile := map[string]interface{}{
		"user_id": userID,
		"name":    "张三",
		"email":   "zhangsan@example.com",
		"role":    "user",
	}

	requestLogger.WithField("profile_loaded", true).Info("用户资料获取成功")
	return profile, nil
}

// 模拟订单服务
type OrderService struct {
	logger *logger.Logger
}

func NewOrderService() *OrderService {
	return &OrderService{
		logger: logger.New().WithPrefix("[OrderService] ").WithShowCaller(true),
	}
}

func (s *OrderService) CreateOrder(userID string, amount float64) (string, error) {
	orderID := generateOrderID()

	orderLogger := s.logger.WithFields(map[string]interface{}{
		"order_id": orderID,
		"user_id":  userID,
		"amount":   amount,
	})

	orderLogger.Info("开始创建订单")

	// 模拟库存检查
	if amount > 1000 {
		orderLogger.Warn("订单金额过大，需要审核")
	}

	// 模拟订单创建过程
	orderLogger.Debug("验证用户权限")
	time.Sleep(30 * time.Millisecond)

	orderLogger.Debug("检查库存")
	time.Sleep(20 * time.Millisecond)

	orderLogger.Debug("计算价格")
	time.Sleep(10 * time.Millisecond)

	orderLogger.WithField("status", "created").Info("订单创建成功")
	return orderID, nil
}

// 模拟支付服务
type PaymentService struct {
	logger *logger.Logger
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		logger: logger.New().
			WithPrefix("[PaymentService] ").
			WithLevel(logger.INFO).
			WithColorful(true),
	}
}

func (s *PaymentService) ProcessPayment(orderID string, amount float64) error {
	paymentLogger := s.logger.WithFields(map[string]interface{}{
		"order_id":   orderID,
		"amount":     amount,
		"payment_id": generatePaymentID(),
	})

	paymentLogger.Info("开始处理支付")

	// 模拟支付过程
	if rand.Float32() < 0.1 { // 10% 概率失败
		err := fmt.Errorf("支付网关连接超时")
		paymentLogger.WithError(err).Error("支付处理失败")
		return err
	}

	time.Sleep(100 * time.Millisecond)
	paymentLogger.WithField("status", "completed").Info("支付处理成功")
	return nil
}

// Web API Handler
type APIHandler struct {
	userService    *UserService
	orderService   *OrderService
	paymentService *PaymentService
	logger         *logger.Logger
}

func NewAPIHandler() *APIHandler {
	return &APIHandler{
		userService:    NewUserService(),
		orderService:   NewOrderService(),
		paymentService: NewPaymentService(),
		logger:         logger.New().WithPrefix("[API] "),
	}
}

func (h *APIHandler) HandleCreateOrder(userID string, amount float64) {
	// 创建请求级别的日志器
	requestID := generateRequestID()
	requestLogger := h.logger.WithFields(map[string]interface{}{
		"request_id": requestID,
		"endpoint":   "/api/orders",
		"method":     "POST",
		"user_id":    userID,
	})

	requestLogger.Info("接收创建订单请求")

	start := time.Now()

	// 1. 获取用户资料
	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		requestLogger.WithError(err).Error("获取用户资料失败")
		return
	}
	requestLogger.WithField("user_profile", profile["name"]).Debug("用户资料获取成功")

	// 2. 创建订单
	orderID, err := h.orderService.CreateOrder(userID, amount)
	if err != nil {
		requestLogger.WithError(err).Error("订单创建失败")
		return
	}
	requestLogger.WithField("order_id", orderID).Info("订单创建成功")

	// 3. 处理支付
	err = h.paymentService.ProcessPayment(orderID, amount)
	if err != nil {
		requestLogger.WithError(err).WithField("order_id", orderID).Error("支付处理失败")
		return
	}

	duration := time.Since(start)
	requestLogger.WithFields(map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
		"duration": duration.String(),
		"status":   "completed",
	}).Info("订单创建流程完成")
}

// 应用程序主类
type Application struct {
	logger     *logger.Logger
	apiHandler *APIHandler
}

func NewApplication() *Application {
	return &Application{
		logger:     logger.New().WithPrefix("[App] ").WithLevel(logger.DEBUG),
		apiHandler: NewAPIHandler(),
	}
}

func (app *Application) Start() {
	app.logger.Info("🚀 应用程序启动")

	// 设置全局日志级别
	logger.SetGlobalLevel(logger.DEBUG)
	app.logger.Debug("全局日志级别设置为DEBUG")

	// 模拟应用启动过程
	app.logger.WithField("component", "database").Info("连接数据库")
	app.logger.WithField("component", "cache").Info("初始化缓存")
	app.logger.WithField("component", "http_server").Info("启动HTTP服务器")

	app.logger.WithFields(map[string]interface{}{
		"port":    8080,
		"mode":    "development",
		"version": "1.0.0",
	}).Info("应用程序启动完成")
}

func (app *Application) SimulateTraffic() {
	app.logger.Info("📊 开始模拟用户请求")

	users := []string{"user001", "user002", "user003", "user004", "user005"}
	amounts := []float64{99.99, 199.99, 299.99, 599.99, 1299.99}

	for i := 0; i < 10; i++ {
		userID := users[rand.Intn(len(users))]
		amount := amounts[rand.Intn(len(amounts))]

		app.logger.WithFields(map[string]interface{}{
			"simulation_round": i + 1,
			"user_id":          userID,
			"amount":           amount,
		}).Debug("模拟用户请求")

		app.apiHandler.HandleCreateOrder(userID, amount)

		// 随机延迟
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}

	app.logger.Info("📊 用户请求模拟完成")
}

func (app *Application) ShowStatistics() {
	app.logger.Info("📈 应用统计信息")

	stats := map[string]interface{}{
		"requests_processed": 10,
		"average_duration":   "150ms",
		"success_rate":       "90%",
		"memory_usage":       "45MB",
		"goroutines":         15,
		"uptime":             time.Since(time.Now().Add(-5 * time.Minute)).String(),
	}

	app.logger.WithFields(stats).Info("当前系统统计")
}

func (app *Application) Shutdown() {
	app.logger.Info("🛑 应用程序关闭")

	app.logger.WithField("component", "http_server").Info("关闭HTTP服务器")
	app.logger.WithField("component", "database").Info("关闭数据库连接")
	app.logger.WithField("component", "cache").Info("清理缓存")

	app.logger.Info("✅ 应用程序关闭完成")
}

func main() {
	fmt.Println("🌟 Go Logger - 实际应用场景示例")
	fmt.Println(strings.Repeat("=", 50))

	// 创建应用程序
	app := NewApplication()

	// 1. 应用启动
	app.Start()
	fmt.Println()

	// 2. 用户登录测试
	fmt.Println("👤 用户服务测试")
	fmt.Println(strings.Repeat("-", 30))
	userService := NewUserService()

	// 成功登录
	err := userService.Login("user001", "correct")
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
	}

	// 失败登录
	err = userService.Login("user002", "wrong")
	if err != nil {
		fmt.Printf("预期的登录失败: %v\n", err)
	}

	fmt.Println()

	// 3. 模拟业务流量
	app.SimulateTraffic()
	fmt.Println()

	// 4. 显示统计信息
	app.ShowStatistics()
	fmt.Println()

	// 5. 演示全局日志器
	fmt.Println("🌍 全局日志器使用")
	fmt.Println(strings.Repeat("-", 30))

	logger.WithField("module", "global").Info("使用全局日志器记录")
	logger.WithFields(map[string]interface{}{
		"event":     "system_event",
		"level":     "info",
		"timestamp": time.Now(),
	}).Info("全局结构化日志")

	fmt.Println()

	// 6. 上下文日志演示
	fmt.Println("📋 上下文日志演示")
	fmt.Println(strings.Repeat("-", 30))
	demonstrateContextLogging()
	fmt.Println()

	// 7. 应用关闭
	app.Shutdown()
}

// 演示上下文日志
func demonstrateContextLogging() {
	ctx := context.Background()

	// 创建带上下文信息的日志器
	contextLogger := logger.New().WithPrefix("[Context] ")

	// 模拟处理用户请求的上下文
	ctx = context.WithValue(ctx, "trace_id", generateTraceID())
	ctx = context.WithValue(ctx, "user_id", "user001")
	ctx = context.WithValue(ctx, "request_path", "/api/users/profile")

	// 使用上下文信息记录日志
	contextLogger.WithFields(map[string]interface{}{
		"trace_id":     ctx.Value("trace_id"),
		"user_id":      ctx.Value("user_id"),
		"request_path": ctx.Value("request_path"),
	}).Info("处理用户请求")

	// 模拟异步操作
	go func() {
		asyncLogger := contextLogger.WithFields(map[string]interface{}{
			"trace_id": ctx.Value("trace_id"),
			"worker":   "async_worker_1",
		})

		asyncLogger.Debug("异步任务开始")
		time.Sleep(100 * time.Millisecond)
		asyncLogger.Info("异步任务完成")
	}()

	time.Sleep(200 * time.Millisecond) // 等待异步任务完成
}

// 辅助函数
func generateSessionID() string {
	return fmt.Sprintf("sess_%d", rand.Int31())
}

func generateOrderID() string {
	return fmt.Sprintf("order_%d", rand.Int31())
}

func generatePaymentID() string {
	return fmt.Sprintf("pay_%d", rand.Int31())
}

func generateRequestID() string {
	return fmt.Sprintf("req_%d", rand.Int31())
}

func generateTraceID() string {
	return fmt.Sprintf("trace_%d", rand.Int31())
}
