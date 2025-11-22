# Go Logger Examples Makefile

.PHONY: all basic adapters factory benchmark configuration enhanced context compatibility formatters monitoring clean help run-all

# 默认目标
all: basic adapters factory benchmark configuration enhanced context compatibility formatters monitoring

# 运行所有示例（新增别名）
run-all: all

# 运行基础示例
basic:
	@echo "=== 运行基础使用示例 ==="
	@cd basic && go run main.go
	@echo ""

# 运行适配器示例
adapters:
	@echo "=== 运行适配器使用示例 ==="
	@cd adapters && go run main.go
	@echo ""

# 运行工厂模式示例
factory:
	@echo "=== 运行工厂模式示例 ==="
	@cd factory && go run main.go
	@echo ""

# 运行性能基准测试示例
benchmark:
	@echo "=== 运行性能基准测试示例 ==="
	@cd benchmark && go run main.go
	@echo ""

# 运行配置使用示例
configuration:
	@echo "=== 运行配置使用示例 ==="
	@cd configuration && go run main.go
	@echo ""

# 运行增强接口功能示例
enhanced:
	@echo "=== 运行增强接口功能示例 ==="
	@cd enhanced && go run main.go
	@echo ""

# 运行上下文感知日志示例
context:
	@echo "=== 运行上下文感知日志示例 ==="
	@cd context && go run main.go
	@echo ""

# 运行多框架兼容性示例
compatibility:
	@echo "=== 运行多框架兼容性示例 ==="
	@cd compatibility && go run main.go
	@echo ""

# 运行格式化器示例（新增）
formatters:
	@echo "=== 运行格式化器示例 ==="
	@cd formatters && go run main.go
	@echo ""

# 运行监控系统示例（新增）
monitoring:
	@echo "=== 运行监控系统示例 ==="
	@cd monitoring && go run main.go
	@echo ""

# 检查所有示例的语法
check:
	@echo "=== 检查示例语法 ==="
	@for dir in basic adapters factory benchmark configuration enhanced context compatibility formatters monitoring; do \
		echo "检查 $$dir..."; \
		cd $$dir && go build -o /dev/null . || exit 1; \
		cd ..; \
	done
	@echo "所有示例语法检查通过"

# 清理构建文件
clean:
	@echo "=== 清理构建文件 ==="
	@find . -name "*.log" -delete
	@find . -name "*.out" -delete
	@echo "清理完成"

# 显示帮助信息
help:
	@echo "Go Logger Examples"
	@echo ""
	@echo "可用目标:"
	@echo "  all           - 运行所有示例"
	@echo "  run-all       - 运行所有示例 (别名)"
	@echo "  basic         - 运行基础使用示例"
	@echo "  adapters      - 运行适配器使用示例"
	@echo "  factory       - 运行工厂模式示例"
	@echo "  benchmark     - 运行性能基准测试示例"
	@echo "  configuration - 运行配置使用示例"
	@echo "  enhanced      - 运行增强接口功能示例"
	@echo "  context       - 运行上下文感知日志示例"
	@echo "  compatibility - 运行多框架兼容性示例"
	@echo "  formatters    - 运行格式化器示例 (新增)"
	@echo "  monitoring    - 运行监控系统示例 (新增)"
	@echo "  check         - 检查所有示例的语法"
	@echo "  clean         - 清理构建文件"
	@echo "  help          - 显示此帮助信息"
	@echo ""
	@echo "核心功能示例:"
	@echo "  formatters    - 演示JSON、Text、Structured、CSV、XML等格式化器"
	@echo "  monitoring    - 演示内存监控、性能分析、I/O监控等功能"
	@echo ""
	@echo "高级功能示例:"
	@echo "  enhanced      - 演示纯文本日志、键值对日志、原始日志方法等新功能"
	@echo "  context       - 演示微服务上下文追踪、分布式日志关联等功能"
	@echo "  compatibility - 演示与Zap、Logrus、slog、Zerolog等框架的兼容性"