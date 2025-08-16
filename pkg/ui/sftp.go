// Package ui 的 SFTP 功能模块
// 提供文件传输的用户界面功能
// 支持上传、下载和交互式文件管理
package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"

	"gossh/internal/sshclient"
)

// StartSFTPSession 启动交互式 SFTP 会话
// 用户可以通过命令行进行文件操作
// 参数:
//   client: SSH 客户端对象
// 返回值:
//   error: 如果会话启动失败则返回错误信息
func StartSFTPSession(client *sshclient.Client) error {
	// 基于 SSH 连接创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(client.GetConnection())
	if err != nil {
		return fmt.Errorf("创建 SFTP 客户端失败: %w", err)
	}
	defer sftpClient.Close() // 会话结束时关闭 SFTP 连接

	// 获取当前远程工作目录
	pwd, err := sftpClient.Getwd()
	if err != nil {
		pwd = "/" // 如果获取失败，默认为根目录
	}

	// 创建标准输入读取器
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("进入 SFTP 交互模式，输入 'help' 查看可用命令")
	fmt.Printf("连接到: %s@%s\n", client.GetConfig().Username, client.GetConfig().Host)
	fmt.Printf("当前远程目录: %s\n", pwd)
	fmt.Println("----------------------------------------")

	// 主命令循环
	for {
		// 显示 SFTP 提示符
		fmt.Print("sftp> ")

		// 读取用户输入
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\n再见!")
				break
			}
			return fmt.Errorf("读取用户输入失败: %w", err)
		}

		// 解析命令和参数
		parts := strings.Fields(strings.TrimSpace(input))
		if len(parts) == 0 {
			continue // 跳过空命令
		}

		command := parts[0]
		args := parts[1:]

		// 执行相应的 SFTP 命令
		if err := executeSFTPCommand(sftpClient, command, args); err != nil {
			fmt.Printf("错误: %v\n", err)
		}

		// 如果是退出命令，跳出循环
		if command == "exit" || command == "quit" {
			break
		}
	}

	return nil
}

// executeSFTPCommand 执行具体的 SFTP 命令
// 根据用户输入的命令执行相应的文件操作
// 参数:
//   client: SFTP 客户端对象
//   command: 用户输入的命令
//   args: 命令参数
// 返回值:
//   error: 如果命令执行失败则返回错误信息
func executeSFTPCommand(client *sftp.Client, command string, args []string) error {
	switch command {
	case "help":
		// 显示帮助信息
		showSFTPHelp()
	case "ls", "dir":
		// 列出远程目录内容
		return listRemoteDirectory(client, args)
	case "pwd":
		// 显示当前远程工作目录
		return showRemotePwd(client)
	case "cd":
		// 切换远程工作目录
		return changeRemoteDirectory(client, args)
	case "get":
		// 下载文件
		return downloadFileCommand(client, args)
	case "put":
		// 上传文件
		return uploadFileCommand(client, args)
	case "mkdir":
		// 创建远程目录
		return createRemoteDirectory(client, args)
	case "rm":
		// 删除远程文件
		return removeRemoteFile(client, args)
	case "exit", "quit":
		// 退出命令
		fmt.Println("再见!")
		return nil
	default:
		return fmt.Errorf("未知命令: %s，输入 'help' 查看可用命令", command)
	}
	return nil
}

// showSFTPHelp 显示 SFTP 命令帮助信息
func showSFTPHelp() {
	fmt.Println("可用的 SFTP 命令:")
	fmt.Println("  ls [目录]     - 列出远程目录内容")
	fmt.Println("  pwd          - 显示当前远程工作目录")
	fmt.Println("  cd <目录>     - 切换远程工作目录")
	fmt.Println("  get <远程文件> [本地文件] - 下载文件")
	fmt.Println("  put <本地文件> [远程文件] - 上传文件")
	fmt.Println("  mkdir <目录>  - 创建远程目录")
	fmt.Println("  rm <文件>     - 删除远程文件")
	fmt.Println("  help         - 显示此帮助信息")
	fmt.Println("  exit/quit    - 退出 SFTP 会话")
}

// listRemoteDirectory 列出远程目录内容
func listRemoteDirectory(client *sftp.Client, args []string) error {
	// 确定要列出的目录
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// 读取目录内容
	files, err := client.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	// 显示文件列表
	fmt.Printf("目录 %s 的内容:\n", dir)
	for _, file := range files {
		// 显示文件类型标识
		fileType := "-"
		if file.IsDir() {
			fileType = "d"
		}
		
		// 显示文件信息：类型、大小、名称
		fmt.Printf("%s %8d %s\n", fileType, file.Size(), file.Name())
	}

	return nil
}

// showRemotePwd 显示当前远程工作目录
func showRemotePwd(client *sftp.Client) error {
	pwd, err := client.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}
	fmt.Println(pwd)
	return nil
}

// changeRemoteDirectory 切换远程工作目录
func changeRemoteDirectory(client *sftp.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定要切换到的目录")
	}

	// SFTP 客户端没有 Chdir 方法，我们需要通过其他方式实现
	// 先检查目录是否存在
	_, err := client.Stat(args[0])
	if err != nil {
		return fmt.Errorf("目录不存在或无法访问: %w", err)
	}

	// 显示提示信息（注意：SFTP 协议本身不支持切换工作目录）
	fmt.Printf("注意: SFTP 协议不支持切换工作目录，请在命令中使用完整路径\n")
	fmt.Printf("目录 %s 存在且可访问\n", args[0])
	return nil
}

// UploadFile 上传文件到远程服务器
// 这是一个公共函数，可以被其他模块调用
// 参数:
//   client: SSH 客户端对象
//   localPath: 本地文件路径
//   remotePath: 远程文件路径
// 返回值:
//   error: 如果上传失败则返回错误信息
func UploadFile(client *sshclient.Client, localPath, remotePath string) error {
	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(client.GetConnection())
	if err != nil {
		return fmt.Errorf("创建 SFTP 客户端失败: %w", err)
	}
	defer sftpClient.Close()

	// 打开本地文件
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %w", err)
	}
	defer localFile.Close()

	// 创建远程文件
	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("创建远程文件失败: %w", err)
	}
	defer remoteFile.Close()

	// 复制文件内容
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("文件传输失败: %w", err)
	}

	return nil
}

// DownloadFile 从远程服务器下载文件
// 这是一个公共函数，可以被其他模块调用
// 参数:
//   client: SSH 客户端对象
//   remotePath: 远程文件路径
//   localPath: 本地文件路径
// 返回值:
//   error: 如果下载失败则返回错误信息
func DownloadFile(client *sshclient.Client, remotePath, localPath string) error {
	// 创建 SFTP 客户端
	sftpClient, err := sftp.NewClient(client.GetConnection())
	if err != nil {
		return fmt.Errorf("创建 SFTP 客户端失败: %w", err)
	}
	defer sftpClient.Close()

	// 打开远程文件
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("打开远程文件失败: %w", err)
	}
	defer remoteFile.Close()

	// 创建本地文件
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %w", err)
	}
	defer localFile.Close()

	// 复制文件内容
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("文件传输失败: %w", err)
	}

	return nil
}

// uploadFileCommand 处理上传文件命令
func uploadFileCommand(client *sftp.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定要上传的本地文件")
	}

	localPath := args[0]
	remotePath := filepath.Base(localPath) // 默认使用文件名作为远程路径

	if len(args) > 1 {
		remotePath = args[1] // 用户指定了远程路径
	}

	// 这里需要将 sftp.Client 转换为 sshclient.Client
	// 实际实现中需要保存原始的 sshclient.Client 引用
	fmt.Printf("上传 %s 到 %s...\n", localPath, remotePath)
	return fmt.Errorf("上传功能需要完整的客户端对象")
}

// downloadFileCommand 处理下载文件命令
func downloadFileCommand(client *sftp.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定要下载的远程文件")
	}

	remotePath := args[0]
	localPath := filepath.Base(remotePath) // 默认使用文件名作为本地路径

	if len(args) > 1 {
		localPath = args[1] // 用户指定了本地路径
	}

	fmt.Printf("下载 %s 到 %s...\n", remotePath, localPath)
	return fmt.Errorf("下载功能需要完整的客户端对象")
}

// createRemoteDirectory 创建远程目录
func createRemoteDirectory(client *sftp.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定要创建的目录名")
	}

	if err := client.Mkdir(args[0]); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	fmt.Printf("目录 %s 创建成功\n", args[0])
	return nil
}

// removeRemoteFile 删除远程文件
func removeRemoteFile(client *sftp.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("请指定要删除的文件名")
	}

	if err := client.Remove(args[0]); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	fmt.Printf("文件 %s 删除成功\n", args[0])
	return nil
}