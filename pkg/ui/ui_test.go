// Package ui_test 提供用户界面模块的单元测试
// 测试文件传输、会话管理等功能
// 使用模拟对象避免实际的网络操作
package ui

import (
	"os"
	"path/filepath"
	"testing"

	"gossh/internal/config"
	"gossh/internal/sshclient"
)

// MockSSHClient 模拟 SSH 客户端，用于测试
// 避免在测试中进行真实的网络连接
type MockSSHClient struct {
	config *config.SSHConfig
	closed bool
}

// NewMockSSHClient 创建模拟的 SSH 客户端
func NewMockSSHClient(cfg *config.SSHConfig) *MockSSHClient {
	return &MockSSHClient{
		config: cfg,
		closed: false,
	}
}

// GetConfig 返回配置信息
func (m *MockSSHClient) GetConfig() *config.SSHConfig {
	return m.config
}

// Close 模拟关闭连接
func (m *MockSSHClient) Close() error {
	m.closed = true
	return nil
}

// ExecuteCommand 模拟执行命令
func (m *MockSSHClient) ExecuteCommand(command string) (string, error) {
	// 模拟一些常见命令的输出
	switch command {
	case "pwd":
		return "/home/root\n", nil
	case "ls":
		return "file1.txt\nfile2.txt\ndir1/\n", nil
	case "whoami":
		return "root\n", nil
	default:
		return "command output for: " + command + "\n", nil
	}
}

// TestUploadFile 测试文件上传功能
// 由于涉及实际的文件操作和网络传输，这里主要测试参数验证
func TestUploadFile(t *testing.T) {
	// 创建临时文件用于测试
	tmpDir, err := os.MkdirTemp("", "upload_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "这是测试文件内容"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 创建模拟客户端
	cfg := &config.SSHConfig{
		Host:     "192.168.1.100",
		Port:     22,
		Username: "root",
		Password: "123456",
	}

	tests := []struct {
		name       string
		localPath  string
		remotePath string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "本地文件不存在",
			localPath:  "/path/to/nonexistent/file",
			remotePath: "/remote/path",
			wantErr:    true,
			errMsg:     "打开本地文件失败",
		},
		{
			name:       "空的本地路径",
			localPath:  "",
			remotePath: "/remote/path",
			wantErr:    true,
			errMsg:     "打开本地文件失败",
		},
		{
			name:       "空的远程路径",
			localPath:  testFile,
			remotePath: "",
			wantErr:    true,
			errMsg:     "创建 SFTP 客户端失败", // 由于没有真实连接，会在这里失败
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 注意：由于 UploadFile 需要真实的 SSH 连接来创建 SFTP 客户端
			// 在单元测试中会失败，但我们可以测试参数验证部分
			client := &sshclient.Client{} // 空客户端，会在 SFTP 创建时失败
			
			err := UploadFile(client, tt.localPath, tt.remotePath)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("UploadFile() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

// TestDownloadFile 测试文件下载功能
func TestDownloadFile(t *testing.T) {
	// 创建临时目录用于下载测试
	tmpDir, err := os.MkdirTemp("", "download_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	downloadPath := filepath.Join(tmpDir, "downloaded.txt")

	tests := []struct {
		name       string
		remotePath string
		localPath  string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "空的远程路径",
			remotePath: "",
			localPath:  downloadPath,
			wantErr:    true,
			errMsg:     "创建 SFTP 客户端失败", // 由于没有真实连接，会在这里失败
		},
		{
			name:       "空的本地路径",
			remotePath: "/remote/file.txt",
			localPath:  "",
			wantErr:    true,
			errMsg:     "创建 SFTP 客户端失败", // 由于没有真实连接，会在这里失败
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &sshclient.Client{} // 空客户端，会在 SFTP 创建时失败
			
			err := DownloadFile(client, tt.remotePath, tt.localPath)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && err != nil {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("DownloadFile() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

// TestSFTPHelp 测试帮助信息显示
// 这个测试验证帮助信息是否能正常显示（不会崩溃）
func TestSFTPHelp(t *testing.T) {
	// 由于 showSFTPHelp 只是打印信息，我们主要测试它不会崩溃
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("showSFTPHelp() 发生 panic: %v", r)
		}
	}()
	
	// 调用帮助函数
	showSFTPHelp()
	
	// 如果执行到这里没有 panic，测试通过
	t.Log("showSFTPHelp() 执行成功")
}

// TestFileOperations 测试文件操作相关的辅助功能
func TestFileOperations(t *testing.T) {
	// 创建临时目录和文件用于测试
	tmpDir, err := os.MkdirTemp("", "file_ops_test")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := "测试文件内容\n包含中文字符"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 测试文件是否存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("测试文件不存在: %v", testFile)
	}

	// 测试文件内容
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("读取测试文件失败: %v", err)
	}
	
	if string(content) != testContent {
		t.Errorf("文件内容不匹配，got = %v, want = %v", string(content), testContent)
	}
}

// BenchmarkStringContains 性能测试 - 字符串包含检查
func BenchmarkStringContains(b *testing.B) {
	testString := "这是一个用于测试字符串包含功能的长字符串，包含了中文和英文 English characters"
	searchString := "字符串包含"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = contains(testString, searchString)
	}
}

// TestStringContains 测试字符串包含功能
func TestStringContains(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{
			name:   "包含子字符串",
			s:      "hello world",
			substr: "world",
			want:   true,
		},
		{
			name:   "不包含子字符串",
			s:      "hello world",
			substr: "golang",
			want:   false,
		},
		{
			name:   "空子字符串",
			s:      "hello world",
			substr: "",
			want:   true,
		},
		{
			name:   "相同字符串",
			s:      "hello",
			substr: "hello",
			want:   true,
		},
		{
			name:   "中文字符串",
			s:      "你好世界",
			substr: "世界",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contains(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

// contains 辅助函数，检查字符串包含关系
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	return containsSubstring(s, substr)
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