package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"cryptobackup/pkg/crypto"
	"cryptobackup/pkg/storage"
	"cryptobackup/pkg/uploader"
)

const (
	version = "1.0.0"
)

func main() {
	// 定义子命令
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)
	genkeyCmd := flag.NewFlagSet("genkey", flag.ExitOnError)

	// upload 命令参数
	uploadFile := uploadCmd.String("file", "", "要上传的本地文件路径")
	uploadRemote := uploadCmd.String("remote", "", "远程文件路径")
	uploadAlgo := uploadCmd.String("algo", "aes", "加密算法 (aes|xor)")
	uploadKey := uploadCmd.String("key", "", "加密密钥（16进制字符串）")
	uploadStorage := uploadCmd.String("storage", "./backup", "存储路径")

	// download 命令参数
	downloadRemote := downloadCmd.String("remote", "", "远程文件路径")
	downloadFile := downloadCmd.String("file", "", "保存到本地的文件路径")
	downloadAlgo := downloadCmd.String("algo", "aes", "加密算法 (aes|xor)")
	downloadKey := downloadCmd.String("key", "", "解密密钥（16进制字符串）")
	downloadStorage := downloadCmd.String("storage", "./backup", "存储路径")

	// list 命令参数
	listPath := listCmd.String("path", "/", "要列出的远程目录路径")
	listStorage := listCmd.String("storage", "./backup", "存储路径")

	// delete 命令参数
	deleteRemote := deleteCmd.String("remote", "", "要删除的远程文件路径")
	deleteStorage := deleteCmd.String("storage", "./backup", "存储路径")

	// info 命令参数
	infoRemote := infoCmd.String("remote", "", "远程文件路径")
	infoStorage := infoCmd.String("storage", "./backup", "存储路径")

	// genkey 命令参数
	genkeySize := genkeyCmd.Int("size", 32, "密钥大小（字节），AES推荐16/24/32")

	// 检查参数
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// 解析子命令
	switch os.Args[1] {
	case "upload":
		uploadCmd.Parse(os.Args[2:])
		if *uploadFile == "" || *uploadRemote == "" || *uploadKey == "" {
			fmt.Println("错误: upload 命令需要 -file, -remote 和 -key 参数")
			uploadCmd.PrintDefaults()
			os.Exit(1)
		}
		handleUpload(*uploadFile, *uploadRemote, *uploadAlgo, *uploadKey, *uploadStorage)

	case "download":
		downloadCmd.Parse(os.Args[2:])
		if *downloadRemote == "" || *downloadFile == "" || *downloadKey == "" {
			fmt.Println("错误: download 命令需要 -remote, -file 和 -key 参数")
			downloadCmd.PrintDefaults()
			os.Exit(1)
		}
		handleDownload(*downloadRemote, *downloadFile, *downloadAlgo, *downloadKey, *downloadStorage)

	case "list":
		listCmd.Parse(os.Args[2:])
		handleList(*listPath, *listStorage)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteRemote == "" {
			fmt.Println("错误: delete 命令需要 -remote 参数")
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}
		handleDelete(*deleteRemote, *deleteStorage)

	case "info":
		infoCmd.Parse(os.Args[2:])
		if *infoRemote == "" {
			fmt.Println("错误: info 命令需要 -remote 参数")
			infoCmd.PrintDefaults()
			os.Exit(1)
		}
		handleInfo(*infoRemote, *infoStorage)

	case "genkey":
		genkeyCmd.Parse(os.Args[2:])
		handleGenKey(*genkeySize)

	case "version":
		fmt.Printf("cryptobackup version %s\n", version)

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Printf("未知命令: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`cryptobackup - 加密文件备份工具 v%s

用法:
  cryptobackup <command> [options]

命令:
  upload      加密并上传文件
  download    下载并解密文件
  list        列出远程文件
  delete      删除远程文件
  info        查看文件信息
  genkey      生成随机密钥
  version     显示版本信息
  help        显示帮助信息

示例:
  # 生成密钥
  cryptobackup genkey -size 32

  # 上传文件
  cryptobackup upload -file ./test.txt -remote /backup/test.txt.enc -key <your-key> -algo aes

  # 下载文件
  cryptobackup download -remote /backup/test.txt.enc -file ./restored.txt -key <your-key> -algo aes

  # 列出文件
  cryptobackup list -path / -storage ./backup

使用 'cryptobackup <command> -h' 查看各命令的详细帮助
`, version)
}

func createEncryptor(algo string, keyHex string) (crypto.Encryptor, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("无效的密钥格式，必须是16进制字符串: %w", err)
	}

	switch algo {
	case "aes":
		return crypto.NewAESEncryptor(key)
	case "xor":
		return crypto.NewXOREncryptor(key)
	default:
		return nil, fmt.Errorf("不支持的加密算法: %s", algo)
	}
}

func handleUpload(localFile, remotePath, algo, keyHex, storagePath string) {
	// 创建加密器
	encryptor, err := createEncryptor(algo, keyHex)
	if err != nil {
		fmt.Printf("创建加密器失败: %v\n", err)
		os.Exit(1)
	}

	// 创建存储
	store, err := storage.NewLocalStorage(storagePath)
	if err != nil {
		fmt.Printf("创建存储失败: %v\n", err)
		os.Exit(1)
	}

	// 创建上传器
	ul := uploader.NewUploader(encryptor, store)

	// 上传文件
	ctx := context.Background()
	fmt.Printf("正在加密并上传文件: %s -> %s\n", localFile, remotePath)
	if err := ul.UploadFile(ctx, localFile, remotePath); err != nil {
		fmt.Printf("上传失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ 上传成功！")
}

func handleDownload(remotePath, localFile, algo, keyHex, storagePath string) {
	// 创建加密器
	encryptor, err := createEncryptor(algo, keyHex)
	if err != nil {
		fmt.Printf("创建加密器失败: %v\n", err)
		os.Exit(1)
	}

	// 创建存储
	store, err := storage.NewLocalStorage(storagePath)
	if err != nil {
		fmt.Printf("创建存储失败: %v\n", err)
		os.Exit(1)
	}

	// 创建上传器
	ul := uploader.NewUploader(encryptor, store)

	// 下载文件
	ctx := context.Background()
	fmt.Printf("正在下载并解密文件: %s -> %s\n", remotePath, localFile)
	if err := ul.DownloadFile(ctx, remotePath, localFile); err != nil {
		fmt.Printf("下载失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ 下载成功！")
}

func handleList(path, storagePath string) {
	// 创建存储
	store, err := storage.NewLocalStorage(storagePath)
	if err != nil {
		fmt.Printf("创建存储失败: %v\n", err)
		os.Exit(1)
	}

	// 列出文件
	ctx := context.Background()
	files, err := store.List(ctx, path)
	if err != nil {
		fmt.Printf("列出文件失败: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("目录为空")
		return
	}

	fmt.Printf("路径: %s\n", path)
	fmt.Println("----------------------------------------")
	for _, file := range files {
		if file.IsDir {
			fmt.Printf("[DIR]  %s\n", filepath.Base(file.Path))
		} else {
			fmt.Printf("[FILE] %s (%d bytes)\n", filepath.Base(file.Path), file.Size)
		}
	}
}

func handleDelete(remotePath, storagePath string) {
	// 创建存储
	store, err := storage.NewLocalStorage(storagePath)
	if err != nil {
		fmt.Printf("创建存储失败: %v\n", err)
		os.Exit(1)
	}

	// 删除文件
	ctx := context.Background()
	fmt.Printf("正在删除文件: %s\n", remotePath)
	if err := store.Delete(ctx, remotePath); err != nil {
		fmt.Printf("删除失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ 删除成功！")
}

func handleInfo(remotePath, storagePath string) {
	// 创建存储
	store, err := storage.NewLocalStorage(storagePath)
	if err != nil {
		fmt.Printf("创建存储失败: %v\n", err)
		os.Exit(1)
	}

	// 获取文件信息
	ctx := context.Background()
	metadata, err := store.GetMetadata(ctx, remotePath)
	if err != nil {
		fmt.Printf("获取文件信息失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("文件信息: %s\n", remotePath)
	fmt.Println("----------------------------------------")
	for k, v := range metadata {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func handleGenKey(size int) {
	key := make([]byte, size)
	if _, err := rand.Read(key); err != nil {
		fmt.Printf("生成密钥失败: %v\n", err)
		os.Exit(1)
	}

	keyHex := hex.EncodeToString(key)
	fmt.Printf("生成的密钥 (%d 字节):\n%s\n", size, keyHex)
	fmt.Println("\n请妥善保管此密钥，丢失后将无法解密文件！")
}
