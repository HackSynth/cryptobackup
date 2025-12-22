package uploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cryptobackup/pkg/crypto"
	"cryptobackup/pkg/storage"
)

// Uploader 文件上传器，整合加密和存储功能
type Uploader struct {
	encryptor crypto.Encryptor
	storage   storage.Storage
}

// NewUploader 创建上传器
func NewUploader(encryptor crypto.Encryptor, storage storage.Storage) *Uploader {
	return &Uploader{
		encryptor: encryptor,
		storage:   storage,
	}
}

// UploadFile 加密并上传文件
func (u *Uploader) UploadFile(ctx context.Context, localPath string, remotePath string) error {
	// 打开本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// 加密文件
	var encryptedData bytes.Buffer
	if err := u.encryptor.Encrypt(file, &encryptedData); err != nil {
		return fmt.Errorf("failed to encrypt file: %w", err)
	}

	// 准备元数据
	metadata := u.encryptor.GetMetadata()
	metadata["original_name"] = filepath.Base(localPath)
	metadata["original_size"] = fmt.Sprintf("%d", fileInfo.Size())
	metadata["encrypted_size"] = fmt.Sprintf("%d", encryptedData.Len())
	metadata["upload_time"] = time.Now().Format(time.RFC3339)

	// 上传到存储
	if err := u.storage.Upload(ctx, remotePath, &encryptedData, metadata); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// DownloadFile 下载并解密文件
func (u *Uploader) DownloadFile(ctx context.Context, remotePath string, localPath string) error {
	// 下载文件
	var encryptedData bytes.Buffer
	if err := u.storage.Download(ctx, remotePath, &encryptedData); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	// 解密文件
	var decryptedData bytes.Buffer
	if err := u.encryptor.Decrypt(&encryptedData, &decryptedData); err != nil {
		return fmt.Errorf("failed to decrypt file: %w", err)
	}

	// 创建本地目录
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 保存到本地文件
	if err := os.WriteFile(localPath, decryptedData.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ListFiles 列出远程文件
func (u *Uploader) ListFiles(ctx context.Context, remotePath string) ([]storage.FileInfo, error) {
	return u.storage.List(ctx, remotePath)
}

// DeleteFile 删除远程文件
func (u *Uploader) DeleteFile(ctx context.Context, remotePath string) error {
	return u.storage.Delete(ctx, remotePath)
}

// GetFileInfo 获取文件信息
func (u *Uploader) GetFileInfo(ctx context.Context, remotePath string) (map[string]string, error) {
	return u.storage.GetMetadata(ctx, remotePath)
}

// UploadStream 加密并上传数据流
func (u *Uploader) UploadStream(ctx context.Context, data io.Reader, remotePath string, metadata map[string]string) error {
	// 加密数据流
	var encryptedData bytes.Buffer
	if err := u.encryptor.Encrypt(data, &encryptedData); err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	// 合并元数据
	finalMetadata := u.encryptor.GetMetadata()
	for k, v := range metadata {
		finalMetadata[k] = v
	}
	finalMetadata["encrypted_size"] = fmt.Sprintf("%d", encryptedData.Len())
	finalMetadata["upload_time"] = time.Now().Format(time.RFC3339)

	// 上传到存储
	if err := u.storage.Upload(ctx, remotePath, &encryptedData, finalMetadata); err != nil {
		return fmt.Errorf("failed to upload data: %w", err)
	}

	return nil
}

// DownloadStream 下载并解密数据流
func (u *Uploader) DownloadStream(ctx context.Context, remotePath string, dst io.Writer) error {
	// 下载文件
	var encryptedData bytes.Buffer
	if err := u.storage.Download(ctx, remotePath, &encryptedData); err != nil {
		return fmt.Errorf("failed to download data: %w", err)
	}

	// 解密数据
	if err := u.encryptor.Decrypt(&encryptedData, dst); err != nil {
		return fmt.Errorf("failed to decrypt data: %w", err)
	}

	return nil
}
