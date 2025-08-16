// go.mod 文件定义了 Go 模块的基本信息
// 这个文件告诉 Go 编译器这是一个名为 gossh 的模块
module gossh

// 指定使用的 Go 版本，确保兼容性
go 1.21

// 项目依赖的第三方包
require (
	github.com/pkg/sftp v1.13.6 // SFTP 客户端功能
	golang.org/x/crypto v0.17.0 // SSH 加密相关功能
	golang.org/x/term v0.15.0 // 终端控制功能
)

require (
	github.com/kr/fs v0.1.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
