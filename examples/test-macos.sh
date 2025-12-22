#!/bin/bash

# CryptoBackup macOS 测试脚本
# 使用方法: chmod +x test-macos.sh && ./test-macos.sh

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 辅助函数
print_header() {
    echo -e "${CYAN}==================================================${NC}"
    echo -e "${CYAN}   CryptoBackup macOS 测试脚本 v1.0.0${NC}"
    echo -e "${CYAN}==================================================${NC}"
    echo ""
}

print_step() {
    echo -e "${GREEN}[步骤 $1] $2${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${CYAN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}$1${NC}"
}

# 检测 macOS 架构
detect_arch() {
    local arch=$(uname -m)
    case $arch in
        x86_64)
            echo "darwin-amd64"
            ;;
        arm64)
            echo "darwin-arm64"
            ;;
        *)
            print_error "不支持的架构: $arch"
            exit 1
            ;;
    esac
}

# 主测试函数
main() {
    print_header

    # 创建测试目录
    TEST_DIR="cryptobackup-test"
    if [ -d "$TEST_DIR" ]; then
        print_warning "清理旧的测试目录..."
        rm -rf "$TEST_DIR"
    fi
    mkdir -p "$TEST_DIR"
    cd "$TEST_DIR"

    # 步骤 1: 下载 macOS 版本
    print_step "1/11" "检测系统并下载对应版本..."
    ARCH=$(detect_arch)
    print_info "检测到架构: $ARCH ($([ "$ARCH" = "darwin-arm64" ] && echo "Apple Silicon" || echo "Intel"))"

    BINARY_NAME="cryptobackup-${ARCH}"
    DOWNLOAD_URL="https://github.com/HackSynth/cryptobackup/releases/download/v1.0.0/${BINARY_NAME}"

    print_info "从 GitHub Release 下载: $BINARY_NAME"

    if command -v curl &> /dev/null; then
        curl -L -# "$DOWNLOAD_URL" -o cryptobackup
    elif command -v wget &> /dev/null; then
        wget -q --show-progress "$DOWNLOAD_URL" -O cryptobackup
    else
        print_error "未找到 curl 或 wget，请安装其中一个"
        exit 1
    fi

    chmod +x cryptobackup
    FILE_SIZE=$(du -h cryptobackup | awk '{print $1}')
    print_success "下载完成 ($FILE_SIZE)"
    echo ""

    # 步骤 2: 验证版本
    print_step "2/11" "验证程序版本..."
    VERSION=$(./cryptobackup version)
    print_info "$VERSION"
    echo ""

    # 步骤 3: 生成加密密钥
    print_step "3/11" "生成加密密钥..."
    KEY_OUTPUT=$(./cryptobackup genkey -size 32)
    echo "$KEY_OUTPUT"
    KEY=$(echo "$KEY_OUTPUT" | grep -v "生成的密钥" | grep -v "请妥善保管" | tr -d '\n' | xargs)
    print_warning "使用密钥: $KEY"
    echo ""

    # 步骤 4: 创建测试文件
    print_step "4/11" "创建测试文件..."
    cat > testfile.txt << 'EOF'
Hello, CryptoBackup on macOS!
This is a test file to verify encryption and decryption.
Line 3: Testing multi-line content.
Line 4: 测试中文字符
Line 5: Special chars: @#$%^&*()
Line 6: macOS path: /Users/user/Documents
Line 7: Emoji test: 🔐 🍎 🚀 ✅
EOF
    print_success "测试文件创建完成"
    print_info "文件内容预览:"
    sed 's/^/  /' testfile.txt
    echo ""

    # 步骤 5: 测试 AES 加密上传
    print_step "5/11" "测试 AES-256-GCM 加密上传..."
    ./cryptobackup upload \
        -file testfile.txt \
        -remote /backup/testfile.txt.enc \
        -key "$KEY" \
        -algo aes \
        -storage ./backup_storage
    print_success "AES 加密上传成功"
    echo ""

    # 步骤 6: 列出备份文件
    print_step "6/11" "列出备份文件..."
    ./cryptobackup list -path / -storage ./backup_storage
    echo ""

    # 步骤 7: 查看文件详细信息
    print_step "7/11" "查看文件元数据..."
    ./cryptobackup info -remote /backup/testfile.txt.enc -storage ./backup_storage
    echo ""

    # 步骤 8: 检查加密文件内容
    print_step "8/11" "检查加密文件内容（应该是乱码）..."
    print_info "加密文件前 100 字节（十六进制）:"
    hexdump -C backup_storage/backup/testfile.txt.enc | head -10
    echo ""

    # 步骤 9: 测试 AES 解密下载
    print_step "9/11" "测试 AES 解密下载..."
    ./cryptobackup download \
        -remote /backup/testfile.txt.enc \
        -file restored.txt \
        -key "$KEY" \
        -algo aes \
        -storage ./backup_storage
    print_success "AES 解密下载成功"
    echo ""

    # 步骤 10: 验证文件完整性
    print_step "10/11" "验证文件完整性..."

    ORIGINAL_HASH=$(shasum -a 256 testfile.txt | cut -d' ' -f1)
    RESTORED_HASH=$(shasum -a 256 restored.txt | cut -d' ' -f1)

    print_info "原始文件 SHA256: $ORIGINAL_HASH"
    print_info "恢复文件 SHA256: $RESTORED_HASH"

    if [ "$ORIGINAL_HASH" = "$RESTORED_HASH" ]; then
        print_success "文件完整性验证通过！文件完全一致。"
    else
        print_error "文件不一致！"
        exit 1
    fi
    echo ""

    # 测试错误密钥
    print_warning "[测试] 使用错误密钥解密（预期失败）..."
    WRONG_KEY="0000000000000000000000000000000000000000000000000000000000000000"
    if ./cryptobackup download \
        -remote /backup/testfile.txt.enc \
        -file wrong.txt \
        -key "$WRONG_KEY" \
        -algo aes \
        -storage ./backup_storage 2>&1 | grep -q "failed"; then
        print_success "正确拒绝了错误的密钥"
    else
        print_error "应该失败但成功了"
        exit 1
    fi
    echo ""

    # 步骤 11: 测试 XOR 加密
    print_step "11/11" "测试 XOR 加密..."

    # XOR 上传
    ./cryptobackup upload \
        -file testfile.txt \
        -remote /backup/testfile_xor.enc \
        -key "$KEY" \
        -algo xor \
        -storage ./backup_storage

    # XOR 下载
    ./cryptobackup download \
        -remote /backup/testfile_xor.enc \
        -file restored_xor.txt \
        -key "$KEY" \
        -algo xor \
        -storage ./backup_storage

    # 验证 XOR
    XOR_HASH=$(shasum -a 256 restored_xor.txt | cut -d' ' -f1)
    if [ "$ORIGINAL_HASH" = "$XOR_HASH" ]; then
        print_success "XOR 加密/解密验证通过"
    else
        print_error "XOR 验证失败"
        exit 1
    fi
    echo ""

    # 显示所有备份文件
    print_info "[信息] 当前所有备份文件:"
    ./cryptobackup list -path /backup -storage ./backup_storage
    echo ""

    # 测试删除功能
    print_warning "[测试] 删除 XOR 加密文件..."
    ./cryptobackup delete -remote /backup/testfile_xor.enc -storage ./backup_storage
    echo ""

    print_info "[信息] 删除后剩余文件:"
    ./cryptobackup list -path /backup -storage ./backup_storage
    echo ""

    # 显示文件大小对比
    print_info "[信息] 文件大小对比:"
    ls -lh testfile.txt restored.txt restored_xor.txt | awk '{print $9, $5}'
    echo ""

    # 显示加密文件对比
    print_info "[信息] 加密文件大小:"
    ls -lh backup_storage/backup/*.enc 2>/dev/null | awk '{print $9, $5}' || echo "无加密文件"
    echo ""

    # 测试完成
    echo -e "${CYAN}==================================================${NC}"
    echo -e "${GREEN}   ✓ 所有测试通过！${NC}"
    echo -e "${CYAN}==================================================${NC}"
    echo ""

    echo -e "${YELLOW}测试总结:${NC}"
    echo -e "  ${GREEN}• 版本验证: 通过${NC}"
    echo -e "  ${GREEN}• 密钥生成: 通过${NC}"
    echo -e "  ${GREEN}• AES-256-GCM 加密: 通过${NC}"
    echo -e "  ${GREEN}• AES-256-GCM 解密: 通过${NC}"
    echo -e "  ${GREEN}• XOR 加密/解密: 通过${NC}"
    echo -e "  ${GREEN}• 文件完整性: 通过${NC}"
    echo -e "  ${GREEN}• 错误密钥拒绝: 通过${NC}"
    echo -e "  ${GREEN}• 文件管理: 通过${NC}"
    echo ""

    # 询问是否清理
    print_info "测试文件保存在: $(pwd)"
    read -p "是否清理测试文件？ (y/N): " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cd ..
        rm -rf "$TEST_DIR"
        print_success "测试文件已清理"
    else
        print_warning "测试文件已保留，可手动删除目录: $TEST_DIR"
    fi
    echo ""
}

# 错误处理
trap 'print_error "测试失败于第 $LINENO 行"; cd ..; exit 1' ERR

# 运行主函数
main

print_info "测试脚本执行完毕"
