package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage 本地存储实现（用于测试或本地备份）
type LocalStorage struct {
	basePath string // 本地存储根目录
}

// NewLocalStorage 创建本地存储
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if basePath == "" {
		return nil, fmt.Errorf("base path cannot be empty")
	}

	// 确保目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	return &LocalStorage{basePath: basePath}, nil
}

// Upload 上传文件到本地存储
func (s *LocalStorage) Upload(ctx context.Context, remotePath string, data io.Reader, metadata map[string]string) error {
	fullPath := filepath.Join(s.basePath, remotePath)

	// 创建目录
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 写入数据
	if _, err := io.Copy(file, data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	// 保存元数据
	if metadata != nil && len(metadata) > 0 {
		if err := s.saveMetadata(fullPath, metadata); err != nil {
			return fmt.Errorf("failed to save metadata: %w", err)
		}
	}

	return nil
}

// Download 从本地存储下载文件
func (s *LocalStorage) Download(ctx context.Context, remotePath string, dst io.Writer) error {
	fullPath := filepath.Join(s.basePath, remotePath)

	file, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	return nil
}

// Delete 删除本地文件
func (s *LocalStorage) Delete(ctx context.Context, remotePath string) error {
	fullPath := filepath.Join(s.basePath, remotePath)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// 同时删除元数据文件
	metaPath := fullPath + ".meta"
	os.Remove(metaPath) // 忽略错误

	return nil
}

// List 列出目录下的文件
func (s *LocalStorage) List(ctx context.Context, remotePath string) ([]FileInfo, error) {
	fullPath := filepath.Join(s.basePath, remotePath)

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []FileInfo
	for _, entry := range entries {
		// 跳过元数据文件
		if filepath.Ext(entry.Name()) == ".meta" {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfo := FileInfo{
			Path:    filepath.Join(remotePath, entry.Name()),
			Size:    info.Size(),
			IsDir:   entry.IsDir(),
			ModTime: info.ModTime().Unix(),
		}

		// 读取元数据
		metadata, _ := s.loadMetadata(filepath.Join(fullPath, entry.Name()))
		fileInfo.Metadata = metadata

		files = append(files, fileInfo)
	}

	return files, nil
}

// Exists 检查文件是否存在
func (s *LocalStorage) Exists(ctx context.Context, remotePath string) (bool, error) {
	fullPath := filepath.Join(s.basePath, remotePath)

	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetMetadata 获取文件元数据
func (s *LocalStorage) GetMetadata(ctx context.Context, remotePath string) (map[string]string, error) {
	fullPath := filepath.Join(s.basePath, remotePath)

	exists, err := s.Exists(ctx, remotePath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("file not found: %s", remotePath)
	}

	metadata, err := s.loadMetadata(fullPath)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// saveMetadata 保存元数据到.meta文件
func (s *LocalStorage) saveMetadata(filePath string, metadata map[string]string) error {
	metaPath := filePath + ".meta"

	data, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	return os.WriteFile(metaPath, data, 0644)
}

// loadMetadata 从.meta文件加载元数据
func (s *LocalStorage) loadMetadata(filePath string) (map[string]string, error) {
	metaPath := filePath + ".meta"

	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	var metadata map[string]string
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

// GetBasePath 获取基础路径
func (s *LocalStorage) GetBasePath() string {
	return s.basePath
}
