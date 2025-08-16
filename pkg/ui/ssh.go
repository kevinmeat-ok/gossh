// Package ui 提供了用户界面相关的功能
// 这个包处理用户交互，包括 SSH 会话和 SFTP 操作
// 让用户能够方便地使用 SSH 功能
package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"

	"gossh/internal/sshclient"
)

// StartSSHSession 启动交互式 SSH 会话
// 用户可以在远程服务器上执行命令，就像本地终端一样
// 参数:
//   client: SSH 客户端对象
// 返回值:
//   error: 如果会话启动失败则返回错误信息
func StartSSHSession(client *sshclient.Client) error {
	// 获取底层的 SSH 连接
	conn := client.GetConnection()

	// 创建一个新的 SSH 会话
	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %w", err)
	}
	defer session.Close() // 会话结束时关闭

	// 获取终端的大小信息
	// 这样远程终端的显示效果会更好
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		// 获取当前终端的宽度和高度
		width, height, err := term.GetSize(fd)
		if err == nil {
			// 请求一个伪终端，设置终端类型和大小
			if err := session.RequestPty("xterm", height, width, ssh.TerminalModes{}); err != nil {
				return fmt.Errorf("请求伪终端失败: %w", err)
			}
		}
	}

	// 连接标准输入输出
	// 让用户的输入能够发送到远程服务器，远程的输出能够显示在本地
	session.Stdin = os.Stdin   // 用户输入发送到远程
	session.Stdout = os.Stdout // 远程输出显示在本地
	session.Stderr = os.Stderr // 远程错误信息显示在本地

	// 启动远程 shell
	// 这会在远程服务器上启动一个交互式 shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("启动远程 shell 失败: %w", err)
	}

	// 等待会话结束
	// 当用户输入 exit 或者连接断开时，会话会结束
	if err := session.Wait(); err != nil {
		return fmt.Errorf("SSH 会话异常结束: %w", err)
	}

	return nil
}

// ExecuteInteractiveCommand 执行交互式命令
// 允许用户输入命令并查看结果，支持多次命令执行
// 参数:
//   client: SSH 客户端对象
// 返回值:
//   error: 如果执行过程中出现错误则返回错误信息
func ExecuteInteractiveCommand(client *sshclient.Client) error {
	// 创建标准输入的读取器
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("进入交互式命令模式，输入 'exit' 退出")
	fmt.Printf("连接到: %s@%s\n", client.GetConfig().Username, client.GetConfig().Host)
	fmt.Println("----------------------------------------")

	// 循环接收用户输入的命令
	for {
		// 显示命令提示符
		fmt.Print("$ ")

		// 读取用户输入的命令
		command, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// 用户按了 Ctrl+D，正常退出
				fmt.Println("\n再见!")
				break
			}
			return fmt.Errorf("读取用户输入失败: %w", err)
		}

		// 去除命令字符串两端的空白字符
		command = strings.TrimSpace(command)

		// 检查是否是退出命令
		if command == "exit" || command == "quit" {
			fmt.Println("再见!")
			break
		}

		// 跳过空命令
		if command == "" {
			continue
		}

		// 在远程服务器上执行命令
		output, err := client.ExecuteCommand(command)
		if err != nil {
			// 显示错误信息，但不退出程序
			fmt.Printf("命令执行失败: %v\n", err)
			continue
		}

		// 显示命令执行结果
		fmt.Print(output)
	}

	return nil
}