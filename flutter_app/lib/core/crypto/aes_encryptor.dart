import 'dart:typed_data';
import 'dart:math';
import 'package:pointycastle/export.dart';
import 'encryptor.dart';

/// AES-GCM 加密器
/// 与 Go 实现完全兼容
class AESEncryptor implements Encryptor {
  final Uint8List key;

  AESEncryptor(this.key) {
    if (key.length != 16 && key.length != 24 && key.length != 32) {
      throw ArgumentError('密钥必须是 16、24 或 32 字节（AES-128/192/256）');
    }
  }

  @override
  Future<List<int>> encrypt(List<int> plaintext) async {
    // 1. 生成随机 nonce（12字节，与 Go 的 gcm.NonceSize() 一致）
    final nonce = _generateNonce(12);

    // 2. 创建 AES-GCM cipher
    final cipher = GCMBlockCipher(AESEngine());

    // 3. 初始化加密模式
    // tag length = 128 bits (16 bytes)，与 Go 默认一致
    final params = AEADParameters(
      KeyParameter(key),
      128,  // tag length in bits
      nonce,
      Uint8List(0),  // additional data 为空，与 Go 实现一致
    );

    cipher.init(true, params);  // true = 加密模式

    // 4. 加密数据
    final ciphertext = cipher.process(Uint8List.fromList(plaintext));

    // 5. 组合输出：nonce(12) + ciphertext + tag(16)
    // 这与 Go 的格式完全一致
    final result = Uint8List(nonce.length + ciphertext.length);
    result.setRange(0, nonce.length, nonce);
    result.setRange(nonce.length, result.length, ciphertext);

    return result;
  }

  @override
  Future<List<int>> decrypt(List<int> ciphertext) async {
    if (ciphertext.length < 12) {
      throw ArgumentError('密文太短，无效的加密数据');
    }

    // 1. 提取 nonce（前12字节）
    final nonce = Uint8List.fromList(ciphertext.sublist(0, 12));

    // 2. 提取加密数据（nonce之后的所有数据，包含tag）
    final encryptedData = Uint8List.fromList(ciphertext.sublist(12));

    // 3. 创建 AES-GCM cipher
    final cipher = GCMBlockCipher(AESEngine());

    // 4. 初始化解密模式
    final params = AEADParameters(
      KeyParameter(key),
      128,
      nonce,
      Uint8List(0),
    );

    cipher.init(false, params);  // false = 解密模式

    // 5. 解密数据
    try {
      final plaintext = cipher.process(encryptedData);
      return plaintext;
    } catch (e) {
      throw Exception('解密失败：密钥错误或数据已损坏 - $e');
    }
  }

  /// 生成随机 nonce
  Uint8List _generateNonce(int length) {
    final random = Random.secure();
    return Uint8List.fromList(
      List.generate(length, (_) => random.nextInt(256))
    );
  }

  @override
  String get algorithmName => 'AES-GCM';

  @override
  Map<String, String> getMetadata() {
    return {
      'algorithm': 'AES-GCM',
      'key_size': '${key.length * 8}',
    };
  }
}
