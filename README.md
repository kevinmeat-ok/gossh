# GoSSH - Go 语言 SSH 客户端工具

这是一个用 Go 语言开发的 SSH 客户端工具，支持 SSH 连接和 SFTP 文件传输功能。专为 PC 客户端设计，支持跨平台部署。

## 功能特性

- **SSH 连接**: 支持密码和密钥认证
- **交互式 Shell**: 在远程服务器上执行命令
- **SFTP 文件传输**: 上传下载文件，支持交互式操作
- **跨平台支持**: Windows、macOS、Linux
- **详细注释**: 代码包含详细的中文注释，适合学习

## 项目结构

```
gossh/
├── cmd/
│   ├── ssh-tool/          # 主程序入口
│   └── sftp/              # SFTP 子命令
├── internal/
│   ├── sshclient/         # SSH 核心逻辑
│   └── config/            # 配置管理
├── pkg/
│   └── ui/                # 用户界面
├── go.mod                 # Go 模块定义
└── README.md              # 项目说明
```

## 安装和编译

1. 确保已安装 Go 1.21 或更高版本
2. 克隆项目到本地
3. 编译主程序：

```bash
# 编译主程序
go build -o ssh-tool ./cmd/ssh-tool

# 编译 SFTP 子命令
go build -o sftp ./cmd/sftp
```

## 使用方法

### SSH 连接

```bash
# 使用密码连接
./ssh-tool -host=192.168.1.100 -user=root -pass=123456

# 使用密钥连接
./ssh-tool -host=192.168.1.100 -user=root -key=/path/to/private/key

# 指定端口
./ssh-tool -host=192.168.1.100 -port=2222 -user=root -pass=123456
```

### SFTP 文件传输

```bash
# 启动交互式 SFTP 会话
./ssh-tool -host=192.168.1.100 -user=root -pass=123456 -mode=sftp

# 或者使用专用的 SFTP 命令
./sftp -host=192.168.1.100 -user=root -pass=123456

# 直接上传文件
./sftp -host=192.168.1.100 -user=root -pass=123456 -upload=/local/file -remote=/remote/path

# 直接下载文件
./sftp -host=192.168.1.100 -user=root -pass=123456 -download=/local/file -remote=/remote/path
```

### SFTP 交互命令

在 SFTP 交互模式下，支持以下命令：

- `ls [目录]` - 列出远程目录内容
- `pwd` - 显示当前远程工作目录
- `cd <目录>` - 切换远程工作目录
- `get <远程文件> [本地文件]` - 下载文件
- `put <本地文件> [远程文件]` - 上传文件
- `mkdir <目录>` - 创建远程目录
- `rm <文件>` - 删除远程文件
- `help` - 显示帮助信息
- `exit` 或 `quit` - 退出会话

## 安全注意事项

- 生产环境中应该验证主机密钥，避免中间人攻击
- 建议使用密钥认证而不是密码认证
- 私钥文件应该设置适当的权限（600）

## 开发说明

项目遵循 Go 语言最佳实践：

- 使用标准的项目结构
- 包含详细的中文注释
- 错误处理遵循 Go 惯用方式
- 支持优雅的资源清理

## 依赖包

- `golang.org/x/crypto` - SSH 加密功能
- `golang.org/x/term` - 终端控制
- `github.com/pkg/sftp` - SFTP 客户端

## 许可证

MIT License