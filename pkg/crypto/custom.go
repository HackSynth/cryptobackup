package crypto

import (
	"fmt"
	"io"
)

// XOREncryptor 简单的XOR加密器（仅作为自定义加密的示例）
// 注意：XOR加密安全性较低，实际生产环境请使用更强的加密算法
type XOREncryptor struct {
	key []byte
}

// NewXOREncryptor 创建XOR加密器
func NewXOREncryptor(key []byte) (*XOREncryptor, error) {
	if len(key) == 0 {
		return nil, fmt.Errorf("key cannot be empty")
	}
	return &XOREncryptor{key: key}, nil
}

// Encrypt 使用XOR加密数据
func (e *XOREncryptor) Encrypt(src io.Reader, dst io.Writer) error {
	return e.xorData(src, dst)
}

// Decrypt 使用XOR解密数据（XOR加密和解密操作相同）
func (e *XOREncryptor) Decrypt(src io.Reader, dst io.Writer) error {
	return e.xorData(src, dst)
}

// xorData XOR数据处理
func (e *XOREncryptor) xorData(src io.Reader, dst io.Writer) error {
	buf := make([]byte, 4096)
	keyLen := len(e.key)
	offset := 0

	for {
		n, err := src.Read(buf)
		if n > 0 {
			// XOR操作
			for i := 0; i < n; i++ {
				buf[i] ^= e.key[offset%keyLen]
				offset++
			}

			if _, writeErr := dst.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("failed to write data: %w", writeErr)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read data: %w", err)
		}
	}

	return nil
}

// GetMetadata 获取加密元数据
func (e *XOREncryptor) GetMetadata() map[string]string {
	return map[string]string{
		"algorithm": "XOR",
		"key_size":  fmt.Sprintf("%d", len(e.key)),
	}
}

// CustomEncryptor 自定义加密器的基础结构
// 用户可以继承这个结构并实现自己的加密逻辑
type CustomEncryptor struct {
	name         string
	encryptFunc  func(io.Reader, io.Writer) error
	decryptFunc  func(io.Reader, io.Writer) error
	metadataFunc func() map[string]string
}

// NewCustomEncryptor 创建自定义加密器
func NewCustomEncryptor(
	name string,
	encrypt func(io.Reader, io.Writer) error,
	decrypt func(io.Reader, io.Writer) error,
	metadata func() map[string]string,
) *CustomEncryptor {
	return &CustomEncryptor{
		name:         name,
		encryptFunc:  encrypt,
		decryptFunc:  decrypt,
		metadataFunc: metadata,
	}
}

// Encrypt 执行加密
func (c *CustomEncryptor) Encrypt(src io.Reader, dst io.Writer) error {
	if c.encryptFunc == nil {
		return fmt.Errorf("encrypt function not implemented")
	}
	return c.encryptFunc(src, dst)
}

// Decrypt 执行解密
func (c *CustomEncryptor) Decrypt(src io.Reader, dst io.Writer) error {
	if c.decryptFunc == nil {
		return fmt.Errorf("decrypt function not implemented")
	}
	return c.decryptFunc(src, dst)
}

// GetMetadata 获取元数据
func (c *CustomEncryptor) GetMetadata() map[string]string {
	if c.metadataFunc == nil {
		return map[string]string{"algorithm": c.name}
	}
	return c.metadataFunc()
}
