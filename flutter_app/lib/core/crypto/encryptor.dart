// 加密器接口
abstract class Encryptor {
  /// 加密数据
  Future<List<int>> encrypt(List<int> plaintext);

  /// 解密数据
  Future<List<int>> decrypt(List<int> ciphertext);

  /// 获取元数据
  Map<String, String> getMetadata();

  /// 算法名称
  String get algorithmName;
}
