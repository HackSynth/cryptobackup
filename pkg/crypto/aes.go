package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// AESEncryptor 使用AES-GCM模式的加密器
type AESEncryptor struct {
	key []byte
}

// NewAESEncryptor 创建AES加密器
// key: 密钥，必须是16、24或32字节（对应AES-128、AES-192、AES-256）
func NewAESEncryptor(key []byte) (*AESEncryptor, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: %d, must be 16, 24 or 32 bytes", len(key))
	}
	return &AESEncryptor{key: key}, nil
}

// Encrypt 使用AES-GCM加密数据
func (e *AESEncryptor) Encrypt(src io.Reader, dst io.Writer) error {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 先写入nonce
	if _, err := dst.Write(nonce); err != nil {
		return fmt.Errorf("failed to write nonce: %w", err)
	}

	// 读取所有数据
	plaintext, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("failed to read source data: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// 写入加密后的数据
	if _, err := dst.Write(ciphertext); err != nil {
		return fmt.Errorf("failed to write encrypted data: %w", err)
	}

	return nil
}

// Decrypt 使用AES-GCM解密数据
func (e *AESEncryptor) Decrypt(src io.Reader, dst io.Writer) error {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	// 读取nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(src, nonce); err != nil {
		return fmt.Errorf("failed to read nonce: %w", err)
	}

	// 读取加密数据
	ciphertext, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("failed to read encrypted data: %w", err)
	}

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}

	// 写入解密后的数据
	if _, err := dst.Write(plaintext); err != nil {
		return fmt.Errorf("failed to write decrypted data: %w", err)
	}

	return nil
}

// GetMetadata 获取加密元数据
func (e *AESEncryptor) GetMetadata() map[string]string {
	return map[string]string{
		"algorithm": "AES-GCM",
		"key_size":  fmt.Sprintf("%d", len(e.key)*8),
	}
}
