// Package config_test 提供配置模块的单元测试
// 测试配置验证、地址格式化等功能
// 确保配置管理功能的正确性
package config

import (
	"os"
	"testing"
)

// TestSSHConfig_Validate 测试配置验证功能
// 验证各种配置情况下的验证结果
func TestSSHConfig_Validate(t *testing.T) {
	// 定义测试用例
	// 每个测试用例包含配置、期望结果和描述
	tests := []struct {
		name    string      // 测试用例名称
		config  *SSHConfig  // 测试的配置
		wantErr bool        // 是否期望出现错误
		errMsg  string      // 期望的错误信息（部分匹配）
	}{
		{
			name: "有效的密码配置",
			config: &SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Username: "root",
				Password: "123456",
			},
			wantErr: false,
		},
		{
			name: "缺少主机地址",
			config: &SSHConfig{
				Port:     22,
				Username: "root",
				Password: "123456",
			},
			wantErr: true,
			errMsg:  "主机地址不能为空",
		},
		{
			name: "缺少用户名",
			config: &SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Password: "123456",
			},
			wantErr: true,
			errMsg:  "用户名不能为空",
		},
		{
			name: "端口号无效（太小）",
			config: &SSHConfig{
				Host:     "192.168.1.100",
				Port:     0,
				Username: "root",
				Password: "123456",
			},
			wantErr: true,
			errMsg:  "端口必须在 1-65535 范围内",
		},
		{
			name: "端口号无效（太大）",
			config: &SSHConfig{
				Host:     "192.168.1.100",
				Port:     70000,
				Username: "root",
				Password: "123456",
			},
			wantErr: true,
			errMsg:  "端口必须在 1-65535 范围内",
		},
		{
			name: "缺少认证信息",
			config: &SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Username: "root",
			},
			wantErr: true,
			errMsg:  "必须提供密码或私钥文件",
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		// 使用 t.Run 创建子测试，便于识别失败的测试用例
		t.Run(tt.name, func(t *testing.T) {
			// 调用被测试的方法
			err := tt.config.Validate()
			
			// 检查错误结果是否符合期望
			if (err != nil) != tt.wantErr {
				t.Errorf("SSHConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// 如果期望有错误，检查错误信息是否包含期望的内容
			if tt.wantErr && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("SSHConfig.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

// TestSSHConfig_GetAddress 测试地址格式化功能
// 验证主机和端口是否正确组合
func TestSSHConfig_GetAddress(t *testing.T) {
	tests := []struct {
		name   string
		config *SSHConfig
		want   string
	}{
		{
			name: "标准端口",
			config: &SSHConfig{
				Host: "192.168.1.100",
				Port: 22,
			},
			want: "192.168.1.100:22",
		},
		{
			name: "自定义端口",
			config: &SSHConfig{
				Host: "example.com",
				Port: 2222,
			},
			want: "example.com:2222",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetAddress()
			if got != tt.want {
				t.Errorf("SSHConfig.GetAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSSHConfig_HasKeyAuth 测试密钥认证检查功能
func TestSSHConfig_HasKeyAuth(t *testing.T) {
	tests := []struct {
		name   string
		config *SSHConfig
		want   bool
	}{
		{
			name: "有密钥文件",
			config: &SSHConfig{
				KeyFile: "/path/to/key",
			},
			want: true,
		},
		{
			name: "无密钥文件",
			config: &SSHConfig{
				KeyFile: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.HasKeyAuth()
			if got != tt.want {
				t.Errorf("SSHConfig.HasKeyAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSSHConfig_HasPasswordAuth 测试密码认证检查功能
func TestSSHConfig_HasPasswordAuth(t *testing.T) {
	tests := []struct {
		name   string
		config *SSHConfig
		want   bool
	}{
		{
			name: "有密码",
			config: &SSHConfig{
				Password: "123456",
			},
			want: true,
		},
		{
			name: "无密码",
			config: &SSHConfig{
				Password: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.HasPasswordAuth()
			if got != tt.want {
				t.Errorf("SSHConfig.HasPasswordAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSSHConfig_ValidateWithKeyFile 测试密钥文件验证
// 这个测试需要创建临时文件
func TestSSHConfig_ValidateWithKeyFile(t *testing.T) {
	// 创建临时密钥文件用于测试
	tmpFile, err := os.CreateTemp("", "test_key")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // 测试结束后删除临时文件
	tmpFile.Close()

	tests := []struct {
		name    string
		config  *SSHConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "存在的密钥文件",
			config: &SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Username: "root",
				KeyFile:  tmpFile.Name(),
			},
			wantErr: false,
		},
		{
			name: "不存在的密钥文件",
			config: &SSHConfig{
				Host:     "192.168.1.100",
				Port:     22,
				Username: "root",
				KeyFile:  "/path/to/nonexistent/key",
			},
			wantErr: true,
			errMsg:  "指定的私钥文件不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SSHConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("SSHConfig.Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

// contains 检查字符串是否包含子字符串
// 这是一个辅助函数，用于错误信息的部分匹配
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
		 (s[:len(substr)] == substr || 
		  s[len(s)-len(substr):] == substr || 
		  containsSubstring(s, substr))))
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