# Go Logger Examples Makefile

.PHONY: all basic adapters factory benchmark configuration clean help

# 默认目标
all: basic adapters factory benchmark configuration

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

# 检查所有示例的语法
check:
	@echo "=== 检查示例语法 ==="
	@for dir in basic adapters factory benchmark configuration; do \
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
	@echo "  basic         - 运行基础使用示例"
	@echo "  adapters      - 运行适配器使用示例"
	@echo "  factory       - 运行工厂模式示例"
	@echo "  benchmark     - 运行性能基准测试示例"
	@echo "  configuration - 运行配置使用示例"
	@echo "  check         - 检查所有示例的语法"
	@echo "  clean         - 清理构建文件"
	@echo "  help          - 显示此帮助信息"