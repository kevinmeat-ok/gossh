// Package sshclient_test 提供 SSH 客户端的单元测试
// 测试客户端创建、连接管理等功能
// 由于涉及网络连接，部分测试使用模拟对象
package sshclient

import (
	"os"
	"testing"

	"golang.org/x/crypto/ssh"

	"gossh/internal/config"
)

// TestNewClient_ConfigValidation 测试客户端创建时的配置验证
// 验证无效配置是否被正确拒绝
func TestNewClient_ConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.SSHConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "无效配置 - 缺少主机",
			config: &config.SSHConfig{
				Port:     22,
				Username: "root",
				Password: "123456",
			},
			wantErr: true,
			errMsg:  "配置验证失败",
		},
		{
			name: "无效配置 - 缺少用户名",
			config: &config.SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Password: "123456",
			},
			wantErr: true,
			errMsg:  "配置验证失败",
		},
		{
			name: "无效配置 - 缺少认证信息",
			config: &config.SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Username: "root",
			},
			wantErr: true,
			errMsg:  "配置验证失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 尝试创建客户端，应该在配置验证阶段失败
			_, err := NewClient(tt.config)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("NewClient() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

// TestAddAuthMethods 测试认证方式配置
// 这个测试不需要实际的网络连接
func TestAddAuthMethods(t *testing.T) {
	// 创建临时密钥文件用于测试
	tmpFile, err := os.CreateTemp("", "test_key")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	// 写入一个简单的测试密钥内容（不是真正的密钥，只是为了测试文件读取）
	testKeyContent := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEA1234567890abcdef
-----END OPENSSH PRIVATE KEY-----`
	
	if _, err := tmpFile.WriteString(testKeyContent); err != nil {
		t.Fatalf("写入测试密钥失败: %v", err)
	}
	tmpFile.Close()

	tests := []struct {
		name    string
		config  *config.SSHConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "密码认证配置",
			config: &config.SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Username: "root",
				Password: "123456",
			},
			wantErr: false,
		},
		{
			name: "密钥文件不存在",
			config: &config.SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Username: "root",
				KeyFile:  "/path/to/nonexistent/key",
			},
			wantErr: true,
			errMsg:  "读取私钥文件失败",
		},
		{
			name: "密钥文件存在但格式错误",
			config: &config.SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Username: "root",
				KeyFile:  tmpFile.Name(),
			},
			wantErr: true,
			errMsg:  "解析私钥失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个空的 SSH 配置用于测试
			sshConfig := &ssh.ClientConfig{
				User: tt.config.Username,
			}
			
			// 测试认证方式配置
			err := addAuthMethods(sshConfig, tt.config)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("addAuthMethods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("addAuthMethods() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
			
			// 如果没有错误，检查认证方式是否被正确添加
			if !tt.wantErr {
				if len(sshConfig.Auth) == 0 {
					t.Error("addAuthMethods() 没有添加任何认证方式")
				}
			}
		})
	}
}

// TestClient_GetConfig 测试获取配置功能
// 使用模拟客户端对象进行测试
func TestClient_GetConfig(t *testing.T) {
	// 创建测试配置
	testConfig := &config.SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
		Password: "123456",
	}
	
	// 创建模拟客户端（不进行实际连接）
	client := &Client{
		config: testConfig,
		conn:   nil, // 在单元测试中不需要真实连接
	}
	
	// 测试获取配置
	got := client.GetConfig()
	if got != testConfig {
		t.Errorf("Client.GetConfig() = %v, want %v", got, testConfig)
	}
	
	// 验证配置内容
	if got.Host != testConfig.Host {
		t.Errorf("Client.GetConfig().Host = %v, want %v", got.Host, testConfig.Host)
	}
	if got.Port != testConfig.Port {
		t.Errorf("Client.GetConfig().Port = %v, want %v", got.Port, testConfig.Port)
	}
	if got.Username != testConfig.Username {
		t.Errorf("Client.GetConfig().Username = %v, want %v", got.Username, testConfig.Username)
	}
}

// TestClient_Close 测试连接关闭功能
func TestClient_Close(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		wantErr bool
	}{
		{
			name: "关闭空连接",
			client: &Client{
				config: &config.SSHConfig{},
				conn:   nil,
			},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// BenchmarkConfigValidation 性能测试 - 配置验证
// 测试配置验证的性能表现
func BenchmarkConfigValidation(b *testing.B) {
	config := &config.SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
		Password: "123456",
	}
	
	// 运行基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}

// contains 辅助函数，检查字符串包含关系
func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsSubstring(s, substr)
}

// containsSubstring 检查字符串中是否包含子字符串
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}