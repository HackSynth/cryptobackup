import 'dart:typed_data';
import 'encryptor.dart';

/// XOR 加密器
/// 与 Go 实现完全兼容
class XOREncryptor implements Encryptor {
  final Uint8List key;

  XOREncryptor(this.key) {
    if (key.isEmpty) {
      throw ArgumentError('密钥不能为空');
    }
  }

  @override
  Future<List<int>> encrypt(List<int> plaintext) async {
    return _xorData(plaintext);
  }

  @override
  Future<List<int>> decrypt(List<int> ciphertext) async {
    return _xorData(ciphertext);
  }

  /// XOR 加密/解密（相同操作）
  List<int> _xorData(List<int> data) {
    final result = Uint8List(data.length);
    for (int i = 0; i < data.length; i++) {
      result[i] = data[i] ^ key[i % key.length];
    }
    return result;
  }

  @override
  String get algorithmName => 'XOR';

  @override
  Map<String, String> getMetadata() {
    return {
      'algorithm': 'XOR',
      'key_size': '${key.length}',
    };
  }
}
