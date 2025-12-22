# CryptoBackup

一个安全、简单的加密文件备份工具，支持多种加密算法和本地存储。

## 特性

- **多种加密算法支持**: AES-256-GCM、自定义XOR加密
- **流式加密**: 支持大文件的流式加密处理
- **本地存储**: 将加密文件安全存储在本地
- **文件管理**: 上传、下载、列表、删除、查看文件信息
- **密钥生成**: 内置安全的随机密钥生成器
- **跨平台**: 支持 Linux、macOS、Windows

## 安装

### 从Release下载

从 [Releases](https://github.com/HackSynth/cryptobackup/releases) 页面下载适合你系统的预编译二进制文件。

### 从源码构建

```bash
git clone https://github.com/HackSynth/cryptobackup.git
cd cryptobackup
go build -o cryptobackup ./cmd/cryptobackup
```

## 快速开始

### 1. 生成加密密钥

```bash
cryptobackup genkey -size 32
```

这将生成一个32字节（256位）的随机密钥，输出为十六进制字符串。**请妥善保管此密钥，丢失后将无法解密文件！**

### 2. 加密并备份文件

```bash
cryptobackup upload \
  -file ./mydata.txt \
  -remote /backup/mydata.txt.enc \
  -key <your-hex-key> \
  -algo aes \
  -storage ./backup
```

参数说明：
- `-file`: 要加密的本地文件路径
- `-remote`: 加密后的文件在存储中的路径
- `-key`: 加密密钥（十六进制字符串）
- `-algo`: 加密算法（`aes` 或 `xor`）
- `-storage`: 本地存储目录（默认: `./backup`）

### 3. 恢复文件

```bash
cryptobackup download \
  -remote /backup/mydata.txt.enc \
  -file ./restored.txt \
  -key <your-hex-key> \
  -algo aes \
  -storage ./backup
```

### 4. 列出备份文件

```bash
cryptobackup list -path / -storage ./backup
```

### 5. 查看文件信息

```bash
cryptobackup info -remote /backup/mydata.txt.enc -storage ./backup
```

### 6. 删除备份文件

```bash
cryptobackup delete -remote /backup/mydata.txt.enc -storage ./backup
```

## 命令详解

### `genkey` - 生成密钥

```bash
cryptobackup genkey -size <bytes>
```

- `-size`: 密钥大小（字节），AES推荐使用 16、24 或 32

### `upload` - 上传文件

```bash
cryptobackup upload -file <local> -remote <remote> -key <key> [-algo <algorithm>] [-storage <path>]
```

### `download` - 下载文件

```bash
cryptobackup download -remote <remote> -file <local> -key <key> [-algo <algorithm>] [-storage <path>]
```

### `list` - 列出文件

```bash
cryptobackup list [-path <path>] [-storage <path>]
```

### `info` - 查看文件信息

```bash
cryptobackup info -remote <remote> [-storage <path>]
```

### `delete` - 删除文件

```bash
cryptobackup delete -remote <remote> [-storage <path>]
```

### `version` - 显示版本

```bash
cryptobackup version
```

## 加密算法

### AES-256-GCM（推荐）

- 使用 AES-256-GCM 对称加密
- 提供认证加密（AEAD）
- 密钥大小: 32字节（256位）
- 高度安全，适合生产环境

使用示例：
```bash
cryptobackup upload -file data.txt -remote /data.enc -key <32-byte-hex> -algo aes
```

### XOR 加密

- 简单的异或加密
- 仅用于测试或非敏感数据
- 不建议用于生产环境

使用示例：
```bash
cryptobackup upload -file data.txt -remote /data.enc -key <hex-key> -algo xor
```

## 安全建议

1. **密钥管理**
   - 使用 `genkey` 命令生成强随机密钥
   - 将密钥存储在安全的位置（密码管理器、密钥保险库等）
   - 不要将密钥提交到版本控制系统
   - 定期轮换密钥

2. **备份存储**
   - 将备份存储目录保存在安全位置
   - 考虑使用额外的物理隔离（如外部硬盘）
   - 定期验证备份完整性

3. **加密算法选择**
   - 生产环境使用 AES 算法
   - 避免在敏感数据上使用 XOR 加密

## 项目结构

```
cryptobackup/
├── cmd/
│   └── cryptobackup/     # 主程序入口
│       └── main.go
├── pkg/
│   ├── crypto/           # 加密模块
│   │   ├── crypto.go     # 加密接口定义
│   │   ├── aes.go        # AES实现
│   │   └── custom.go     # XOR实现
│   ├── storage/          # 存储模块
│   │   ├── storage.go    # 存储接口定义
│   │   └── local.go      # 本地存储实现
│   └── uploader/         # 上传下载模块
│       └── uploader.go
├── go.mod
└── README.md
```

## 开发

### 环境要求

- Go 1.21 或更高版本

### 构建

```bash
go build -o cryptobackup ./cmd/cryptobackup
```

### 运行测试

```bash
go test ./...
```

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 作者

HackSynth
