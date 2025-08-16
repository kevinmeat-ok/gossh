# 测试配置说明

## 单元测试

项目包含完整的单元测试，覆盖了主要的功能模块：

### 测试文件结构
```
gossh/
├── internal/config/config_test.go      # 配置模块测试
├── internal/sshclient/client_test.go   # SSH 客户端测试
├── pkg/ui/ui_test.go                   # UI 模块测试
└── integration_test.go                 # 集成测试
```

## 运行测试

### 运行所有单元测试
```bash
# 运行所有测试
go test ./...

# 运行测试并显示详细输出
go test -v ./...

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 运行特定模块的测试
```bash
# 只测试配置模块
go test ./internal/config

# 只测试 SSH 客户端模块
go test ./internal/sshclient

# 只测试 UI 模块
go test ./pkg/ui
```

### 运行性能测试
```bash
# 运行基准测试
go test -bench=. ./...

# 运行特定的基准测试
go test -bench=BenchmarkConfigValidation ./internal/config
```

## 集成测试

集成测试需要真实的 SSH 服务器环境。

### 设置环境变量
```bash
# 设置测试服务器信息
export RUN_INTEGRATION_TESTS=1
export TEST_SSH_HOST=192.168.1.100
export TEST_SSH_USER=root
export TEST_SSH_PASS=your_password

# 运行集成测试
go test -v integration_test.go
```

### 使用 Docker 搭建测试环境
```bash
# 启动 SSH 测试服务器
docker run -d --name ssh-test-server \
  -p 2222:22 \
  -e SSH_ENABLE_PASSWORD_AUTH=true \
  panubo/sshd:latest

# 设置环境变量
export TEST_SSH_HOST=localhost:2222
export TEST_SSH_USER=root
export TEST_SSH_PASS=root
```

## 测试覆盖的功能

### 配置模块 (config_test.go)
- ✅ 配置验证功能
- ✅ 地址格式化
- ✅ 认证方式检查
- ✅ 密钥文件验证
- ✅ 错误处理

### SSH 客户端模块 (client_test.go)
- ✅ 客户端创建
- ✅ 配置验证
- ✅ 认证方式配置
- ✅ 连接管理
- ✅ 错误处理

### UI 模块 (ui_test.go)
- ✅ 文件上传下载参数验证
- ✅ 帮助信息显示
- ✅ 字符串处理功能
- ✅ 文件操作辅助功能

### 集成测试 (integration_test.go)
- ✅ 完整的连接流程
- ✅ 命令执行功能
- ✅ 错误处理场景
- ✅ 配置到客户端的集成

## 测试最佳实践

### 1. 测试命名规范
- 测试函数以 `Test` 开头
- 基准测试以 `Benchmark` 开头
- 使用描述性的测试名称

### 2. 测试结构
- 使用表驱动测试 (table-driven tests)
- 每个测试用例包含名称、输入、期望输出
- 使用 `t.Run()` 创建子测试

### 3. 模拟对象
- 对于网络操作使用模拟对象
- 避免在单元测试中进行真实的网络连接
- 使用依赖注入便于测试

### 4. 错误测试
- 测试正常情况和异常情况
- 验证错误信息的准确性
- 测试边界条件

## 持续集成

可以在 CI/CD 流水线中使用以下命令：

```bash
# 安装依赖
go mod download

# 运行测试
go test -v -race -coverprofile=coverage.out ./...

# 检查覆盖率
go tool cover -func=coverage.out

# 生成覆盖率报告
go tool cover -html=coverage.out -o coverage.html
```

## 测试数据

测试使用的示例数据：
- 主机地址: `192.168.1.100`
- 端口: `22`
- 用户名: `root`
- 密码: `123456`
- 测试文件内容: 包含中文和英文的混合内容

注意：这些都是测试数据，不要在生产环境中使用。