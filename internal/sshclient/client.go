// Package sshclient 提供了 SSH 连接的核心功能
// 这个包封装了 SSH 和 SFTP 的底层操作
// 让上层应用更容易使用 SSH 功能
package sshclient

import (
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"

	"gossh/internal/config"
)

// Client 表示一个 SSH 客户端连接
// 这个结构体包含了 SSH 连接和相关的配置信息
type Client struct {
	config *config.SSHConfig // SSH 连接配置
	conn   *ssh.Client       // SSH 连接对象
}

// NewClient 创建一个新的 SSH 客户端
// 根据提供的配置信息建立 SSH 连接
// 参数:
//   cfg: SSH 连接配置信息
// 返回值:
//   *Client: 创建的客户端对象
//   error: 如果连接失败则返回错误信息
func NewClient(cfg *config.SSHConfig) (*Client, error) {
	// 验证配置信息是否有效
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	// 创建 SSH 客户端配置
	sshConfig := &ssh.ClientConfig{
		User:            cfg.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：生产环境应该验证主机密钥
		Timeout:         30 * time.Second,            // 连接超时时间
	}

	// 根据配置添加认证方式
	if err := addAuthMethods(sshConfig, cfg); err != nil {
		return nil, fmt.Errorf("配置认证方式失败: %w", err)
	}

	// 建立 SSH 连接
	conn, err := ssh.Dial("tcp", cfg.GetAddress(), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("SSH 连接失败: %w", err)
	}

	// 创建客户端对象
	client := &Client{
		config: cfg,
		conn:   conn,
	}

	return client, nil
}

// addAuthMethods 为 SSH 配置添加认证方式
// 支持密码认证和密钥认证
// 参数:
//   sshConfig: SSH 客户端配置对象
//   cfg: 用户提供的配置信息
// 返回值:
//   error: 如果配置认证方式失败则返回错误
func addAuthMethods(sshConfig *ssh.ClientConfig, cfg *config.SSHConfig) error {
	var authMethods []ssh.AuthMethod

	// 如果配置了密码，添加密码认证
	if cfg.HasPasswordAuth() {
		authMethods = append(authMethods, ssh.Password(cfg.Password))
	}

	// 如果配置了密钥文件，添加密钥认证
	if cfg.HasKeyAuth() {
		// 读取私钥文件内容
		keyData, err := ioutil.ReadFile(cfg.KeyFile)
		if err != nil {
			return fmt.Errorf("读取私钥文件失败: %w", err)
		}

		// 解析私钥
		signer, err := ssh.ParsePrivateKey(keyData)
		if err != nil {
			return fmt.Errorf("解析私钥失败: %w", err)
		}

		// 添加公钥认证方式
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// 将认证方式设置到 SSH 配置中
	sshConfig.Auth = authMethods
	return nil
}

// GetConnection 返回底层的 SSH 连接对象
// 供其他模块使用原始的 SSH 连接
// 返回值:
//   *ssh.Client: SSH 连接对象
func (c *Client) GetConnection() *ssh.Client {
	return c.conn
}

// GetConfig 返回客户端的配置信息
// 返回值:
//   *config.SSHConfig: 配置对象
func (c *Client) GetConfig() *config.SSHConfig {
	return c.config
}

// ExecuteCommand 在远程服务器上执行单个命令
// 执行命令并返回输出结果
// 参数:
//   command: 要执行的命令字符串
// 返回值:
//   string: 命令的输出结果
//   error: 如果执行失败则返回错误信息
func (c *Client) ExecuteCommand(command string) (string, error) {
	// 创建一个新的会话
	session, err := c.conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建会话失败: %w", err)
	}
	defer session.Close() // 使用完毕后关闭会话

	// 执行命令并获取输出
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("执行命令失败: %w", err)
	}

	return string(output), nil
}

// Close 关闭 SSH 连接
// 释放网络资源，程序结束前应该调用此方法
// 返回值:
//   error: 如果关闭失败则返回错误信息
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}