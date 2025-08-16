# Makefile for GoSSH project
# 这个文件定义了编译和管理项目的常用命令

# 定义变量
BINARY_NAME=ssh-tool
SFTP_BINARY=sftp
BUILD_DIR=build

# 默认目标：编译所有程序
all: build

# 编译主程序和 SFTP 工具
build:
	@echo "正在编译 SSH 工具..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ssh-tool
	@echo "正在编译 SFTP 工具..."
	go build -o $(BUILD_DIR)/$(SFTP_BINARY) ./cmd/sftp
	@echo "编译完成！可执行文件位于 $(BUILD_DIR)/ 目录"

# 清理编译产物
clean:
	@echo "清理编译文件..."
	rm -rf $(BUILD_DIR)
	@echo "清理完成！"

# 运行测试
test:
	@echo "运行单元测试..."
	go test -v ./...

# 运行测试并显示覆盖率
test-cover:
	@echo "运行测试并生成覆盖率报告..."
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# 生成 HTML 覆盖率报告
test-html:
	@echo "生成 HTML 覆盖率报告..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 运行基准测试
bench:
	@echo "运行性能测试..."
	go test -bench=. ./...

# 运行集成测试
test-integration:
	@echo "运行集成测试..."
	@echo "请确保设置了必要的环境变量："
	@echo "  export RUN_INTEGRATION_TESTS=1"
	@echo "  export TEST_SSH_HOST=your_host"
	@echo "  export TEST_SSH_USER=your_user"
	@echo "  export TEST_SSH_PASS=your_pass"
	go test -v integration_test.go

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 检查代码
vet:
	@echo "检查代码..."
	go vet ./...

# 下载依赖
deps:
	@echo "下载依赖包..."
	go mod download
	go mod tidy

# 安装到系统
install: build
	@echo "安装到系统..."
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	cp $(BUILD_DIR)/$(SFTP_BINARY) /usr/local/bin/
	@echo "安装完成！"

# 创建发布包
release: clean build
	@echo "创建发布包..."
	mkdir -p release
	tar -czf release/gossh-$(shell date +%Y%m%d).tar.gz $(BUILD_DIR)/ README.md
	@echo "发布包创建完成！"

# 显示帮助信息
help:
	@echo "可用的命令："
	@echo "  make build    - 编译程序"
	@echo "  make clean    - 清理编译文件"
	@echo "  make test     - 运行单元测试"
	@echo "  make test-cover - 运行测试并显示覆盖率"
	@echo "  make test-html  - 生成 HTML 覆盖率报告"
	@echo "  make bench      - 运行性能测试"
	@echo "  make test-integration - 运行集成测试"
	@echo "  make fmt      - 格式化代码"
	@echo "  make vet      - 检查代码"
	@echo "  make deps     - 下载依赖"
	@echo "  make install  - 安装到系统"
	@echo "  make release  - 创建发布包"
	@echo "  make help     - 显示此帮助"

# 声明伪目标
.PHONY: all build clean test test-cover test-html bench test-integration fmt vet deps install release help