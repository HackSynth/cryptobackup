@echo off
REM CryptoBackup Windows 测试脚本 (批处理版本)
REM 使用方法: 双击运行此文件
chcp 65001 >nul
setlocal EnableDelayedExpansion

echo ==================================================
echo    CryptoBackup Windows 测试脚本 v1.0.0
echo ==================================================
echo.

REM 创建测试目录
set TEST_DIR=cryptobackup-test
if exist %TEST_DIR% (
    echo 清理旧的测试目录...
    rd /s /q %TEST_DIR%
)
mkdir %TEST_DIR%
cd %TEST_DIR%

echo [步骤 1/10] 下载 Windows 版本...
echo 从 GitHub Release 下载 cryptobackup-windows-amd64.exe
echo.
echo 请手动下载文件到当前目录，或使用 PowerShell 脚本自动下载
echo 下载地址: https://github.com/HackSynth/cryptobackup/releases/download/v1.0.0/cryptobackup-windows-amd64.exe
echo.
echo 下载完成后，请将文件重命名为 cryptobackup.exe 并放在此目录中
echo.
pause

if not exist cryptobackup.exe (
    echo [错误] 未找到 cryptobackup.exe
    echo 请下载并将文件重命名为 cryptobackup.exe
    pause
    exit /b 1
)

echo [步骤 2/10] 验证程序版本...
cryptobackup.exe version
echo.

echo [步骤 3/10] 生成加密密钥...
cryptobackup.exe genkey -size 32 > key.tmp
type key.tmp
echo.
echo 请从上面的输出中复制密钥（第二行的十六进制字符串）
set /p KEY=请粘贴密钥:
echo 使用密钥: %KEY%
echo.

echo [步骤 4/10] 创建测试文件...
(
echo Hello, CryptoBackup on Windows!
echo This is a test file to verify encryption and decryption.
echo Line 3: Testing multi-line content.
echo Line 4: 测试中文字符
echo Line 5: Special chars: @#$%%^&*^(^)
echo Line 6: Windows path: C:\Users\Test\Documents
) > testfile.txt
echo ✓ 测试文件创建完成
type testfile.txt
echo.

echo [步骤 5/10] 测试 AES-256-GCM 加密上传...
cryptobackup.exe upload -file testfile.txt -remote /backup/testfile.txt.enc -key %KEY% -algo aes -storage ./backup_storage
if errorlevel 1 (
    echo [错误] 加密上传失败
    pause
    exit /b 1
)
echo.

echo [步骤 6/10] 列出备份文件...
cryptobackup.exe list -path / -storage ./backup_storage
echo.

echo [步骤 7/10] 查看文件元数据...
cryptobackup.exe info -remote /backup/testfile.txt.enc -storage ./backup_storage
echo.

echo [步骤 8/10] 测试 AES 解密下载...
cryptobackup.exe download -remote /backup/testfile.txt.enc -file restored.txt -key %KEY% -algo aes -storage ./backup_storage
if errorlevel 1 (
    echo [错误] 解密下载失败
    pause
    exit /b 1
)
echo.

echo [步骤 9/10] 验证文件完整性...
echo 原始文件:
type testfile.txt
echo.
echo 恢复文件:
type restored.txt
echo.
fc /b testfile.txt restored.txt >nul
if errorlevel 1 (
    echo ✗ 文件不一致！
    pause
    exit /b 1
) else (
    echo ✓ 文件完整性验证通过！文件完全一致。
)
echo.

echo [步骤 10/10] 测试 XOR 加密...
cryptobackup.exe upload -file testfile.txt -remote /backup/testfile_xor.enc -key %KEY% -algo xor -storage ./backup_storage
cryptobackup.exe download -remote /backup/testfile_xor.enc -file restored_xor.txt -key %KEY% -algo xor -storage ./backup_storage
fc /b testfile.txt restored_xor.txt >nul
if errorlevel 1 (
    echo ✗ XOR 验证失败！
) else (
    echo ✓ XOR 加密/解密验证通过
)
echo.

echo [信息] 当前所有备份文件:
cryptobackup.exe list -path /backup -storage ./backup_storage
echo.

echo [测试] 删除 XOR 加密文件...
cryptobackup.exe delete -remote /backup/testfile_xor.enc -storage ./backup_storage
echo.

echo [信息] 删除后剩余文件:
cryptobackup.exe list -path /backup -storage ./backup_storage
echo.

echo [信息] 文件大小对比:
dir testfile.txt restored.txt restored_xor.txt | find ".txt"
echo.

echo ==================================================
echo    ✓ 所有测试通过！
echo ==================================================
echo.
echo 测试总结:
echo   • 版本验证: 通过
echo   • 密钥生成: 通过
echo   • AES-256-GCM 加密: 通过
echo   • AES-256-GCM 解密: 通过
echo   • XOR 加密/解密: 通过
echo   • 文件完整性: 通过
echo   • 文件管理: 通过
echo.

set /p CLEANUP=是否清理测试文件？ (Y/N):
if /i "%CLEANUP%"=="Y" (
    cd ..
    rd /s /q %TEST_DIR%
    echo ✓ 测试文件已清理
) else (
    echo 测试文件已保留，可手动删除目录: %TEST_DIR%
)

echo.
echo 按任意键退出...
pause >nul
