# CryptoBackup Flutter APP

CryptoBackup 的 Flutter 跨平台应用，支持 Android 和 iOS。

## 功能特性

- ✅ **AES-256-GCM 加密**：与 Go 后端完全兼容
- ✅ **XOR 加密**：简单加密算法（测试用）
- ✅ **本地独立运行**：无需服务器，完全离线工作
- ✅ **现代化 UI**：渐变设计，与 Web UI 保持一致
- 📱 **跨平台支持**：Android + iOS

## 项目结构

```
flutter_app/
├── lib/
│   ├── main.dart                    # 应用入口
│   ├── core/                        # 核心功能
│   │   └── crypto/                  # 加密模块
│   │       ├── encryptor.dart       # 加密接口
│   │       ├── aes_encryptor.dart   # AES-GCM实现 ✅
│   │       └── xor_encryptor.dart   # XOR实现 ✅
│   └── ui/                          # 界面层
│       └── theme/
│           └── app_theme.dart       # 主题配置 ✅
├── pubspec.yaml                     # 依赖配置 ✅
└── README.md                        # 本文件
```

## 核心依赖

```yaml
dependencies:
  pointycastle: ^3.7.3    # AES-GCM加密
  hex: ^0.2.0             # 十六进制转换
  file_picker: ^6.1.1     # 文件选择
  path_provider: ^2.1.1   # 路径获取
  sqflite: ^2.3.0         # SQLite数据库
  flutter_animate: ^4.3.0 # 动画效果
```

## 快速开始

### 1. 安装依赖

```bash
cd flutter_app
flutter pub get
```

### 2. 运行应用

```bash
# Android
flutter run

# iOS（需要 macOS）
flutter run -d ios
```

### 3. 构建发布版

```bash
# Android APK
flutter build apk --release

# Android App Bundle
flutter build appbundle --release

# iOS
flutter build ios --release
```

## AES-GCM 加密实现

### 与 Go 兼容性说明

Flutter 的 AES-GCM 实现与 Go 后端完全兼容：

**数据格式**：
```
[nonce(12字节)] + [ciphertext + tag(16字节)]
```

**关键参数**：
- Nonce Size: 12 字节
- Tag Length: 128 位（16 字节）
- Additional Data: 空

### 使用示例

```dart
import 'dart:typed_data';
import 'package:hex/hex.dart';
import 'core/crypto/aes_encryptor.dart';

// 1. 创建加密器
final keyHex = '0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef';
final key = Uint8List.fromList(HEX.decode(keyHex));
final encryptor = AESEncryptor(key);

// 2. 加密
final plaintext = Uint8List.fromList('Hello CryptoBackup!'.codeUnits);
final ciphertext = await encryptor.encrypt(plaintext);

// 3. 解密
final decrypted = await encryptor.decrypt(ciphertext);
print(String.fromCharCodes(decrypted)); // Hello CryptoBackup!
```

## 兼容性测试

### 测试步骤

1. **Dart 加密 → Go 解密**：
```bash
# Dart加密文件
flutter test test/crypto_test.dart

# Go解密文件
cryptobackup download -remote /tmp/dart_encrypted.bin -file /tmp/decrypted.txt -key <hex> -algo aes
```

2. **Go 加密 → Dart 解密**：
```bash
# Go加密文件
cryptobackup upload -file test.txt -remote /tmp/go_encrypted.bin -key <hex> -algo aes

# Dart解密文件
flutter test test/crypto_test.dart
```

## UI 主题

### 渐变色方案（与 Web UI 一致）

```dart
// 主渐变
primaryGradient: [#667EEA, #764BA2]

// 次级渐变
secondaryGradient: [#F093FB, #F5576C]

// 成功渐变
successGradient: [#11998E, #38EF7D]
```

## 待实现功能

- [ ] 文件列表页面
- [ ] 加密文件页面
- [ ] 解密文件页面
- [ ] 密钥生成器页面
- [ ] 文件存储管理（SQLite）
- [ ] 文件分享功能
- [ ] 设置页面

## 开发注意事项

1. **权限配置**：
   - Android: `WRITE_EXTERNAL_STORAGE`, `READ_EXTERNAL_STORAGE`
   - iOS: `NSPhotoLibraryUsageDescription`

2. **大文件处理**：
   - 限制文件大小（建议 ≤ 100MB）
   - 显示加密/解密进度

3. **密钥管理**：
   - 使用 `flutter_secure_storage` 安全存储密钥
   - 或每次操作时输入密钥

## 许可证

MIT License

## 作者

HackSynth
