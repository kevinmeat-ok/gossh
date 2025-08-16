// Package main 是 SSH 工具的主程序入口
// 这个程序提供了 SSH 连接和 SFTP 文件传输功能
// 支持命令行参数，方便用户使用
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gossh/internal/config"
	"gossh/internal/sshclient"
	"gossh/pkg/ui"
)

// main 是程序的入口函数
// 负责解析命令行参数，初始化配置，启动相应的功能模块
func main() {
	// 定义命令行参数
	// 这些参数让用户可以通过命令行指定连接信息
	var (
		host     = flag.String("host", "", "SSH 服务器地址 (必填)")
		port     = flag.Int("port", 22, "SSH 服务器端口 (默认: 22)")
		username = flag.String("user", "", "用户名 (必填)")
		password = flag.String("pass", "", "密码")
		keyFile  = flag.String("key", "", "私钥文件路径")
		mode     = flag.String("mode", "ssh", "运行模式: ssh 或 sftp (默认: ssh)")
	)

	// 解析命令行参数
	flag.Parse()

	// 检查必填参数
	// 如果用户没有提供必要的连接信息，显示帮助信息并退出
	if *host == "" || *username == "" {
		fmt.Println("错误: 必须提供主机地址和用户名")
		fmt.Println("\n使用示例:")
		fmt.Println("  ssh-tool -host=192.168.1.100 -user=root -pass=123456")
		fmt.Println("  ssh-tool -host=192.168.1.100 -user=root -key=/path/to/key -mode=sftp")
		flag.Usage()
		os.Exit(1)
	}

	// 创建 SSH 配置对象
	// 将用户输入的参数封装成配置结构体
	cfg := &config.SSHConfig{
		Host:     *host,
		Port:     *port,
		Username: *username,
		Password: *password,
		KeyFile:  *keyFile,
	}

	// 创建 SSH 客户端
	// 这个客户端负责实际的 SSH 连接和操作
	client, err := sshclient.NewClient(cfg)
	if err != nil {
		log.Fatalf("创建 SSH 客户端失败: %v", err)
	}
	defer client.Close() // 程序结束时关闭连接

	// 根据用户选择的模式启动相应功能
	switch *mode {
	case "ssh":
		// 启动 SSH 交互模式
		// 用户可以在远程服务器上执行命令
		fmt.Printf("正在连接到 %s@%s:%d...\n", *username, *host, *port)
		if err := ui.StartSSHSession(client); err != nil {
			log.Fatalf("SSH 会话启动失败: %v", err)
		}
	case "sftp":
		// 启动 SFTP 文件传输模式
		// 用户可以上传下载文件
		fmt.Printf("正在启动 SFTP 会话到 %s@%s:%d...\n", *username, *host, *port)
		if err := ui.StartSFTPSession(client); err != nil {
			log.Fatalf("SFTP 会话启动失败: %v", err)
		}
	default:
		// 用户输入了不支持的模式
		fmt.Printf("错误: 不支持的模式 '%s'，请使用 'ssh' 或 'sftp'\n", *mode)
		os.Exit(1)
	}
}