import 'package:flutter/material.dart';
import 'ui/theme/app_theme.dart';

void main() {
  runApp(const CryptoBackupApp());
}

class CryptoBackupApp extends StatelessWidget {
  const CryptoBackupApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'CryptoBackup',
      theme: AppTheme.lightTheme,
      home: const HomePage(),
      debugShowCheckedModeBanner: false,
    );
  }
}

/// 主页 - 文件列表
class HomePage extends StatelessWidget {
  const HomePage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: AppTheme.primaryGradient,
        ),
        child: SafeArea(
          child: Column(
            children: [
              // 毛玻璃 AppBar
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.1),
                  border: Border(
                    bottom: BorderSide(
                      color: Colors.white.withOpacity(0.2),
                    ),
                  ),
                ),
                child: Row(
                  children: [
                    const Icon(Icons.shield, color: Colors.white),
                    const SizedBox(width: 8),
                    const Text(
                      'CryptoBackup',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const Spacer(),
                    IconButton(
                      icon: const Icon(Icons.settings, color: Colors.white),
                      onPressed: () {},
                    ),
                  ],
                ),
              ),

              // 内容区域
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      // 欢迎卡片
                      Card(
                        elevation: 8,
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(16),
                        ),
                        child: Container(
                          padding: const EdgeInsets.all(24),
                          decoration: BoxDecoration(
                            gradient: AppTheme.primaryGradient,
                            borderRadius: BorderRadius.circular(16),
                          ),
                          child: Column(
                            children: [
                              const Icon(
                                Icons.lock,
                                size: 64,
                                color: Colors.white,
                              ),
                              const SizedBox(height: 16),
                              const Text(
                                '欢迎使用 CryptoBackup',
                                style: TextStyle(
                                  color: Colors.white,
                                  fontSize: 24,
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                              const SizedBox(height: 8),
                              Text(
                                '安全的文件加密备份应用',
                                style: TextStyle(
                                  color: Colors.white.withOpacity(0.9),
                                  fontSize: 16,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),

                      const SizedBox(height: 24),

                      // 功能按钮
                      Expanded(
                        child: GridView.count(
                          crossAxisCount: 2,
                          crossAxisSpacing: 16,
                          mainAxisSpacing: 16,
                          children: [
                            _buildFeatureCard(
                              context,
                              icon: Icons.upload_file,
                              title: '加密文件',
                              gradient: AppTheme.primaryGradient,
                              onTap: () {
                                // TODO: 导航到加密页面
                              },
                            ),
                            _buildFeatureCard(
                              context,
                              icon: Icons.folder,
                              title: '文件列表',
                              gradient: AppTheme.secondaryGradient,
                              onTap: () {
                                // TODO: 导航到文件列表
                              },
                            ),
                            _buildFeatureCard(
                              context,
                              icon: Icons.key,
                              title: '生成密钥',
                              gradient: AppTheme.successGradient,
                              onTap: () {
                                // TODO: 导航到密钥生成
                              },
                            ),
                            _buildFeatureCard(
                              context,
                              icon: Icons.lock_open,
                              title: '解密文件',
                              gradient: const LinearGradient(
                                colors: [Color(0xFFFA709A), Color(0xFFFEE140)],
                              ),
                              onTap: () {
                                // TODO: 导航到解密页面
                              },
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildFeatureCard(
    BuildContext context, {
    required IconData icon,
    required String title,
    required Gradient gradient,
    required VoidCallback onTap,
  }) {
    return Card(
      elevation: 4,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Container(
          decoration: BoxDecoration(
            gradient: gradient,
            borderRadius: BorderRadius.circular(16),
          ),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                icon,
                size: 48,
                color: Colors.white,
              ),
              const SizedBox(height: 12),
              Text(
                title,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
