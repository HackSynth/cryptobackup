package crypto

import "io"

// Encryptor 定义加密接口，支持自定义加密实现
type Encryptor interface {
	// Encrypt 加密数据流
	// src: 源数据读取器
	// dst: 加密后数据写入器
	Encrypt(src io.Reader, dst io.Writer) error

	// Decrypt 解密数据流
	// src: 加密数据读取器
	// dst: 解密后数据写入器
	Decrypt(src io.Reader, dst io.Writer) error

	// GetMetadata 获取加密元数据（如加密算法、密钥信息等）
	GetMetadata() map[string]string
}

// Config 加密配置
type Config struct {
	Algorithm string                 // 加密算法名称
	Key       []byte                 // 加密密钥
	Options   map[string]interface{} // 额外选项
}
