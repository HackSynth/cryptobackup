#!/bin/bash

# CryptoBackup 使用示例脚本

echo "=== CryptoBackup 使用示例 ==="
echo ""

# 1. 生成密钥
echo "步骤 1: 生成加密密钥"
echo "运行命令: cryptobackup genkey -size 32"
echo ""
KEY=$(cryptobackup genkey -size 32 | grep -v "生成的密钥" | grep -v "请妥善保管" | tr -d '\n')
echo "生成的密钥: $KEY"
echo ""

# 2. 创建测试文件
echo "步骤 2: 创建测试文件"
echo "Hello, CryptoBackup!" > test.txt
echo "This is a test file for encryption." >> test.txt
echo "创建文件: test.txt"
cat test.txt
echo ""

# 3. 上传并加密文件
echo "步骤 3: 加密并备份文件"
echo "运行命令: cryptobackup upload -file test.txt -remote /backup/test.txt.enc -key \$KEY -algo aes"
cryptobackup upload -file test.txt -remote /backup/test.txt.enc -key "$KEY" -algo aes -storage ./backup
echo ""

# 4. 列出备份文件
echo "步骤 4: 列出备份文件"
echo "运行命令: cryptobackup list -path / -storage ./backup"
cryptobackup list -path / -storage ./backup
echo ""

# 5. 查看文件信息
echo "步骤 5: 查看文件信息"
echo "运行命令: cryptobackup info -remote /backup/test.txt.enc -storage ./backup"
cryptobackup info -remote /backup/test.txt.enc -storage ./backup
echo ""

# 6. 下载并解密文件
echo "步骤 6: 下载并解密文件"
echo "运行命令: cryptobackup download -remote /backup/test.txt.enc -file restored.txt -key \$KEY -algo aes"
cryptobackup download -remote /backup/test.txt.enc -file restored.txt -key "$KEY" -algo aes -storage ./backup
echo ""

# 7. 验证恢复的文件
echo "步骤 7: 验证恢复的文件"
echo "原始文件:"
cat test.txt
echo ""
echo "恢复的文件:"
cat restored.txt
echo ""

# 8. 清理
echo "步骤 8: 清理测试文件"
rm -f test.txt restored.txt
rm -rf ./backup
echo "清理完成！"
echo ""

echo "=== 示例完成 ==="
