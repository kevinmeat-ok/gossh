// Package config 提供了 SSH 连接的配置管理功能
// 这个包定义了连接参数的结构体和验证方法
// 让配置信息更容易管理和传递
package config

import (
	"errors"
	"fmt"
	"os"
)

// SSHConfig 定义了 SSH 连接所需的所有配置信息
// 这个结构体包含了连接远程服务器需要的所有参数
type SSHConfig struct {
	Host     string // 服务器地址，如 "192.168.1.100" 或 "example.com"
	Port     int    // 服务器端口，通常是 22
	Username string // 登录用户名
	Password string // 登录密码（可选，也可以使用密钥）
	KeyFile  string // 私钥文件路径（可选，用于密钥认证）
}

// Validate 验证配置信息是否完整和有效
// 检查必填字段是否为空，以及文件是否存在
// 返回值:
//   error: 如果配置无效则返回错误信息，否则返回 nil
func (c *SSHConfig) Validate() error {
	// 检查主机地址是否为空
	if c.Host == "" {
		return errors.New("主机地址不能为空")
	}

	// 检查用户名是否为空
	if c.Username == "" {
		return errors.New("用户名不能为空")
	}

	// 检查端口是否在有效范围内
	if c.Port <= 0 || c.Port > 65535 {
		return errors.New("端口必须在 1-65535 范围内")
	}

	// 检查认证方式：必须提供密码或密钥文件
	if c.Password == "" && c.KeyFile == "" {
		return errors.New("必须提供密码或私钥文件")
	}

	// 如果指定了密钥文件，检查文件是否存在
	if c.KeyFile != "" {
		if _, err := os.Stat(c.KeyFile); os.IsNotExist(err) {
			return errors.New("指定的私钥文件不存在: " + c.KeyFile)
		}
	}

	return nil
}

// GetAddress 返回完整的服务器地址
// 将主机和端口组合成 "host:port" 格式
// 返回值:
//   string: 格式化的地址字符串
func (c *SSHConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// HasKeyAuth 检查是否使用密钥认证
// 返回值:
//   bool: 如果配置了密钥文件则返回 true，否则返回 false
func (c *SSHConfig) HasKeyAuth() bool {
	return c.KeyFile != ""
}

// HasPasswordAuth 检查是否使用密码认证
// 返回值:
//   bool: 如果配置了密码则返回 true，否则返回 false
func (c *SSHConfig) HasPasswordAuth() bool {
	return c.Password != ""
}