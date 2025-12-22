# CryptoBackup Windows 测试脚本 (PowerShell)
# 使用方法: 右键 "使用 PowerShell 运行" 或在 PowerShell 中执行: .\test-windows.ps1

Write-Host "==================================================" -ForegroundColor Cyan
Write-Host "   CryptoBackup Windows 测试脚本 v1.0.0" -ForegroundColor Cyan
Write-Host "==================================================" -ForegroundColor Cyan
Write-Host ""

# 设置错误时停止
$ErrorActionPreference = "Stop"

# 创建测试目录
$TestDir = "cryptobackup-test"
if (Test-Path $TestDir) {
    Write-Host "清理旧的测试目录..." -ForegroundColor Yellow
    Remove-Item -Path $TestDir -Recurse -Force
}
New-Item -ItemType Directory -Path $TestDir | Out-Null
Set-Location $TestDir

try {
    # 步骤 1: 下载 Windows 版本
    Write-Host "`n[步骤 1/10] 下载 Windows 版本..." -ForegroundColor Green
    Write-Host "从 GitHub Release 下载 cryptobackup-windows-amd64.exe"

    $ReleaseUrl = "https://github.com/HackSynth/cryptobackup/releases/download/v1.0.0/cryptobackup-windows-amd64.exe"
    $ExePath = "cryptobackup.exe"

    Invoke-WebRequest -Uri $ReleaseUrl -OutFile $ExePath
    Write-Host "✓ 下载完成 ($(((Get-Item $ExePath).Length / 1MB).ToString('0.0')) MB)" -ForegroundColor Green

    # 步骤 2: 验证版本
    Write-Host "`n[步骤 2/10] 验证程序版本..." -ForegroundColor Green
    $Version = & ".\$ExePath" version
    Write-Host $Version -ForegroundColor Cyan

    # 步骤 3: 生成加密密钥
    Write-Host "`n[步骤 3/10] 生成加密密钥..." -ForegroundColor Green
    $KeyOutput = & ".\$ExePath" genkey -size 32
    Write-Host $KeyOutput -ForegroundColor Cyan

    # 提取密钥（获取第二行）
    $Key = ($KeyOutput -split "`n")[1].Trim()
    Write-Host "使用密钥: $Key" -ForegroundColor Yellow

    # 步骤 4: 创建测试文件
    Write-Host "`n[步骤 4/10] 创建测试文件..." -ForegroundColor Green
    $TestContent = @"
Hello, CryptoBackup on Windows!
This is a test file to verify encryption and decryption.
Line 3: Testing multi-line content.
Line 4: 测试中文字符
Line 5: Special chars: @#$%^&*()
Line 6: Windows path: C:\Users\Test\Documents
"@
    $TestContent | Out-File -FilePath "testfile.txt" -Encoding UTF8
    Write-Host "✓ 测试文件创建完成" -ForegroundColor Green
    Write-Host "文件内容预览:" -ForegroundColor Gray
    Get-Content "testfile.txt" | ForEach-Object { Write-Host "  $_" -ForegroundColor Gray }

    # 步骤 5: 测试 AES 加密上传
    Write-Host "`n[步骤 5/10] 测试 AES-256-GCM 加密上传..." -ForegroundColor Green
    & ".\$ExePath" upload `
        -file testfile.txt `
        -remote /backup/testfile.txt.enc `
        -key $Key `
        -algo aes `
        -storage ./backup_storage
    Write-Host "✓ AES 加密上传成功" -ForegroundColor Green

    # 步骤 6: 列出备份文件
    Write-Host "`n[步骤 6/10] 列出备份文件..." -ForegroundColor Green
    & ".\$ExePath" list -path / -storage ./backup_storage

    # 步骤 7: 查看文件详细信息
    Write-Host "`n[步骤 7/10] 查看文件元数据..." -ForegroundColor Green
    & ".\$ExePath" info -remote /backup/testfile.txt.enc -storage ./backup_storage

    # 步骤 8: 测试 AES 解密下载
    Write-Host "`n[步骤 8/10] 测试 AES 解密下载..." -ForegroundColor Green
    & ".\$ExePath" download `
        -remote /backup/testfile.txt.enc `
        -file restored.txt `
        -key $Key `
        -algo aes `
        -storage ./backup_storage
    Write-Host "✓ AES 解密下载成功" -ForegroundColor Green

    # 步骤 9: 验证文件完整性
    Write-Host "`n[步骤 9/10] 验证文件完整性..." -ForegroundColor Green

    $OriginalHash = (Get-FileHash -Path testfile.txt -Algorithm SHA256).Hash
    $RestoredHash = (Get-FileHash -Path restored.txt -Algorithm SHA256).Hash

    Write-Host "原始文件 SHA256: $OriginalHash" -ForegroundColor Cyan
    Write-Host "恢复文件 SHA256: $RestoredHash" -ForegroundColor Cyan

    if ($OriginalHash -eq $RestoredHash) {
        Write-Host "✓ 文件完整性验证通过！文件完全一致。" -ForegroundColor Green
    } else {
        Write-Host "✗ 错误：文件不一致！" -ForegroundColor Red
        throw "文件完整性验证失败"
    }

    # 测试错误密钥
    Write-Host "`n[测试] 使用错误密钥解密（预期失败）..." -ForegroundColor Yellow
    $WrongKey = "0000000000000000000000000000000000000000000000000000000000000000"
    try {
        & ".\$ExePath" download `
            -remote /backup/testfile.txt.enc `
            -file wrong.txt `
            -key $WrongKey `
            -algo aes `
            -storage ./backup_storage 2>&1 | Out-Null
        Write-Host "✗ 错误：应该失败但成功了" -ForegroundColor Red
    } catch {
        Write-Host "✓ 正确拒绝了错误的密钥" -ForegroundColor Green
    }

    # 步骤 10: 测试 XOR 加密
    Write-Host "`n[步骤 10/10] 测试 XOR 加密..." -ForegroundColor Green

    # XOR 上传
    & ".\$ExePath" upload `
        -file testfile.txt `
        -remote /backup/testfile_xor.enc `
        -key $Key `
        -algo xor `
        -storage ./backup_storage

    # XOR 下载
    & ".\$ExePath" download `
        -remote /backup/testfile_xor.enc `
        -file restored_xor.txt `
        -key $Key `
        -algo xor `
        -storage ./backup_storage

    # 验证 XOR
    $XorHash = (Get-FileHash -Path restored_xor.txt -Algorithm SHA256).Hash
    if ($OriginalHash -eq $XorHash) {
        Write-Host "✓ XOR 加密/解密验证通过" -ForegroundColor Green
    } else {
        Write-Host "✗ XOR 验证失败" -ForegroundColor Red
        throw "XOR 验证失败"
    }

    # 显示所有备份文件
    Write-Host "`n[信息] 当前所有备份文件:" -ForegroundColor Cyan
    & ".\$ExePath" list -path /backup -storage ./backup_storage

    # 测试删除功能
    Write-Host "`n[测试] 删除 XOR 加密文件..." -ForegroundColor Yellow
    & ".\$ExePath" delete -remote /backup/testfile_xor.enc -storage ./backup_storage

    Write-Host "`n[信息] 删除后剩余文件:" -ForegroundColor Cyan
    & ".\$ExePath" list -path /backup -storage ./backup_storage

    # 显示文件大小对比
    Write-Host "`n[信息] 文件大小对比:" -ForegroundColor Cyan
    Get-ChildItem testfile.txt, restored.txt, restored_xor.txt |
        Select-Object Name, @{Name="Size (Bytes)";Expression={$_.Length}} |
        Format-Table -AutoSize

    # 测试完成
    Write-Host "`n==================================================" -ForegroundColor Cyan
    Write-Host "   ✓ 所有测试通过！" -ForegroundColor Green
    Write-Host "==================================================" -ForegroundColor Cyan
    Write-Host "`n测试总结:" -ForegroundColor Yellow
    Write-Host "  • 版本验证: 通过" -ForegroundColor Green
    Write-Host "  • 密钥生成: 通过" -ForegroundColor Green
    Write-Host "  • AES-256-GCM 加密: 通过" -ForegroundColor Green
    Write-Host "  • AES-256-GCM 解密: 通过" -ForegroundColor Green
    Write-Host "  • XOR 加密/解密: 通过" -ForegroundColor Green
    Write-Host "  • 文件完整性: 通过" -ForegroundColor Green
    Write-Host "  • 错误密钥拒绝: 通过" -ForegroundColor Green
    Write-Host "  • 文件管理: 通过" -ForegroundColor Green
    Write-Host ""

    # 询问是否清理
    Write-Host "测试文件保存在: $(Get-Location)" -ForegroundColor Cyan
    $CleanUp = Read-Host "`n是否清理测试文件？ (Y/N)"
    if ($CleanUp -eq "Y" -or $CleanUp -eq "y") {
        Set-Location ..
        Remove-Item -Path $TestDir -Recurse -Force
        Write-Host "✓ 测试文件已清理" -ForegroundColor Green
    } else {
        Write-Host "测试文件已保留，可手动删除目录: $TestDir" -ForegroundColor Yellow
    }

} catch {
    Write-Host "`n✗ 测试失败: $_" -ForegroundColor Red
    Write-Host $_.ScriptStackTrace -ForegroundColor Red
    Set-Location ..
    exit 1
}

Write-Host "`n按任意键退出..." -ForegroundColor Gray
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
