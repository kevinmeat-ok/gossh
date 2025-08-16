// Package main 提供了 SFTP 子命令的功能
// 这个程序专门用于文件传输操作
// 可以作为独立的 SFTP 客户端使用
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

// main 是 SFTP 子命令的入口函数
// 专门处理文件传输相关的操作
func main() {
	// 定义 SFTP 专用的命令行参数
	var (
		host     = flag.String("host", "", "SFTP 服务器地址 (必填)")
		port     = flag.Int("port", 22, "SFTP 服务器端口 (默认: 22)")
		username = flag.String("user", "", "用户名 (必填)")
		password = flag.String("pass", "", "密码")
		keyFile  = flag.String("key", "", "私钥文件路径")
		upload   = flag.String("upload", "", "上传文件路径")
		download = flag.String("download", "", "下载文件路径")
		remote   = flag.String("remote", "", "远程文件路径")
	)

	// 解析命令行参数
	flag.Parse()

	// 检查必填参数
	if *host == "" || *username == "" {
		fmt.Println("错误: 必须提供主机地址和用户名")
		fmt.Println("\n使用示例:")
		fmt.Println("  sftp -host=192.168.1.100 -user=root -pass=123456")
		fmt.Println("  sftp -host=192.168.1.100 -user=root -key=/path/to/key -upload=/local/file -remote=/remote/path")
		flag.Usage()
		os.Exit(1)
	}

	// 创建 SSH 配置
	cfg := &config.SSHConfig{
		Host:     *host,
		Port:     *port,
		Username: *username,
		Password: *password,
		KeyFile:  *keyFile,
	}

	// 创建 SSH 客户端
	client, err := sshclient.NewClient(cfg)
	if err != nil {
		log.Fatalf("创建 SSH 客户端失败: %v", err)
	}
	defer client.Close()

	// 根据参数决定操作模式
	if *upload != "" && *remote != "" {
		// 上传文件模式
		fmt.Printf("正在上传文件 %s 到 %s...\n", *upload, *remote)
		if err := ui.UploadFile(client, *upload, *remote); err != nil {
			log.Fatalf("文件上传失败: %v", err)
		}
		fmt.Println("文件上传成功!")
	} else if *download != "" && *remote != "" {
		// 下载文件模式
		fmt.Printf("正在下载文件 %s 到 %s...\n", *remote, *download)
		if err := ui.DownloadFile(client, *remote, *download); err != nil {
			log.Fatalf("文件下载失败: %v", err)
		}
		fmt.Println("文件下载成功!")
	} else {
		// 交互式 SFTP 模式
		fmt.Printf("正在启动 SFTP 会话到 %s@%s:%d...\n", *username, *host, *port)
		if err := ui.StartSFTPSession(client); err != nil {
			log.Fatalf("SFTP 会话启动失败: %v", err)
		}
	}
}