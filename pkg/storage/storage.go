package storage

import (
	"context"
	"io"
)

// Storage 定义网盘存储接口，方便后续接入不同网盘API
type Storage interface {
	// Upload 上传文件
	// ctx: 上下文，用于超时控制和取消
	// remotePath: 远程路径
	// data: 文件数据
	// metadata: 文件元数据
	Upload(ctx context.Context, remotePath string, data io.Reader, metadata map[string]string) error

	// Download 下载文件
	// ctx: 上下文
	// remotePath: 远程路径
	// dst: 目标写入器
	Download(ctx context.Context, remotePath string, dst io.Writer) error

	// Delete 删除文件
	// ctx: 上下文
	// remotePath: 远程路径
	Delete(ctx context.Context, remotePath string) error

	// List 列出文件
	// ctx: 上下文
	// remotePath: 远程目录路径
	List(ctx context.Context, remotePath string) ([]FileInfo, error)

	// Exists 检查文件是否存在
	// ctx: 上下文
	// remotePath: 远程路径
	Exists(ctx context.Context, remotePath string) (bool, error)

	// GetMetadata 获取文件元数据
	// ctx: 上下文
	// remotePath: 远程路径
	GetMetadata(ctx context.Context, remotePath string) (map[string]string, error)
}

// FileInfo 文件信息
type FileInfo struct {
	Path         string            // 文件路径
	Size         int64             // 文件大小（字节）
	IsDir        bool              // 是否是目录
	ModTime      int64             // 修改时间（Unix时间戳）
	Metadata     map[string]string // 元数据
}

// Config 存储配置
type Config struct {
	Type    string                 // 存储类型（如 "local", "baidu", "aliyun" 等）
	Options map[string]interface{} // 配置选项
}
