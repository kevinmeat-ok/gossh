// Package main 提供集成测试
// 这些测试验证各个模块之间的协作是否正常
// 注意：这些测试需要真实的 SSH 服务器环境
package main

import (
	"os"
	"testing"

	"gossh/internal/config"
	"gossh/internal/sshclient"
)

// TestIntegration_ConfigToClient 测试配置到客户端的完整流程
// 这是一个集成测试，验证从配置创建到客户端初始化的完整过程
func TestIntegration_ConfigToClient(t *testing.T) {
	// 跳过集成测试，除非设置了环境变量
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.Skip("跳过集成测试，设置 RUN_INTEGRATION_TESTS 环境变量来运行")
	}

	// 从环境变量获取测试服务器信息
	testHost := os.Getenv("TEST_SSH_HOST")
	testUser := os.Getenv("TEST_SSH_USER")
	testPass := os.Getenv("TEST_SSH_PASS")

	if testHost == "" || testUser == "" || testPass == "" {
		t.Skip("跳过集成测试：缺少必要的环境变量 TEST_SSH_HOST, TEST_SSH_USER, TEST_SSH_PASS")
	}

	// 创建测试配置
	cfg := &config.SSHConfig{
		Host:     testHost,
		Port:     22,
		Username: testUser,
		Password: testPass,
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		t.Fatalf("配置验证失败: %v", err)
	}

	// 创建客户端（这会尝试实际连接）
	client, err := sshclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("创建 SSH 客户端失败: %v", err)
	}
	defer client.Close()

	// 测试基本功能
	t.Run("执行简单命令", func(t *testing.T) {
		output, err := client.ExecuteCommand("echo 'Hello, World!'")
		if err != nil {
			t.Errorf("执行命令失败: %v", err)
		}
		
		expected := "Hello, World!\n"
		if output != expected {
			t.Errorf("命令输出不匹配，got = %v, want = %v", output, expected)
		}
	})

	t.Run("获取当前目录", func(t *testing.T) {
		output, err := client.ExecuteCommand("pwd")
		if err != nil {
			t.Errorf("执行 pwd 命令失败: %v", err)
		}
		
		if len(output) == 0 {
			t.Error("pwd 命令没有返回输出")
		}
		
		t.Logf("当前目录: %s", output)
	})

	t.Run("检查用户身份", func(t *testing.T) {
		output, err := client.ExecuteCommand("whoami")
		if err != nil {
			t.Errorf("执行 whoami 命令失败: %v", err)
		}
		
		if len(output) == 0 {
			t.Error("whoami 命令没有返回输出")
		}
		
		t.Logf("当前用户: %s", output)
	})
}

// TestIntegration_FileTransfer 测试文件传输集成功能
func TestIntegration_FileTransfer(t *testing.T) {
	// 跳过集成测试，除非设置了环境变量
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.Skip("跳过集成测试，设置 RUN_INTEGRATION_TESTS 环境变量来运行")
	}

	// 从环境变量获取测试服务器信息
	testHost := os.Getenv("TEST_SSH_HOST")
	testUser := os.Getenv("TEST_SSH_USER")
	testPass := os.Getenv("TEST_SSH_PASS")

	if testHost == "" || testUser == "" || testPass == "" {
		t.Skip("跳过集成测试：缺少必要的环境变量")
	}

	// 创建临时文件用于测试
	tmpFile, err := os.CreateTemp("", "integration_test_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试内容
	testContent := "这是集成测试文件\n包含中文内容\nIntegration test file"
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("写入测试内容失败: %v", err)
	}
	tmpFile.Close()

	// 创建配置和客户端
	cfg := &config.SSHConfig{
		Host:     testHost,
		Port:     22,
		Username: testUser,
		Password: testPass,
	}

	client, err := sshclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("创建 SSH 客户端失败: %v", err)
	}
	defer client.Close()

	// 注意：由于 UI 包中的文件传输函数需要实际的 SFTP 连接
	// 这里我们主要测试客户端是否能够成功创建
	t.Run("客户端连接成功", func(t *testing.T) {
		if client.GetConnection() == nil {
			t.Error("SSH 连接为空")
		}
		
		if client.GetConfig() == nil {
			t.Error("客户端配置为空")
		}
		
		// 验证配置信息
		gotConfig := client.GetConfig()
		if gotConfig.Host != testHost {
			t.Errorf("主机地址不匹配，got = %v, want = %v", gotConfig.Host, testHost)
		}
		
		if gotConfig.Username != testUser {
			t.Errorf("用户名不匹配，got = %v, want = %v", gotConfig.Username, testUser)
		}
	})
}

// TestIntegration_ErrorHandling 测试错误处理的集成场景
func TestIntegration_ErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		config *config.SSHConfig
		errMsg string
	}{
		{
			name: "连接不存在的主机",
			config: &config.SSHConfig{
				Host:     "192.168.255.255", // 不存在的 IP
				Port:     22,
				Username: "root",
				Password: "123456",
			},
			errMsg: "SSH 连接失败",
		},
		{
			name: "连接错误的端口",
			config: &config.SSHConfig{
				Host:     "127.0.0.1",
				Port:     12345, // 不存在的端口
				Username: "root",
				Password: "123456",
			},
			errMsg: "SSH 连接失败",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 这些测试应该快速失败，不需要长时间等待
			_, err := sshclient.NewClient(tt.config)
			
			if err == nil {
				t.Error("期望连接失败，但连接成功了")
				return
			}
			
			if !contains(err.Error(), tt.errMsg) {
				t.Errorf("错误信息不匹配，got = %v, want containing = %v", err.Error(), tt.errMsg)
			}
			
			t.Logf("正确捕获错误: %v", err)
		})
	}
}

// contains 辅助函数，检查字符串包含关系
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}