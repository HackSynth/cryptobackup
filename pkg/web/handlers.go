package web

import (
	"bytes"
	"context"
	"crypto/rand"
	"cryptobackup/pkg/crypto"
	"cryptobackup/pkg/uploader"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Handler holds the server configuration and handles HTTP requests
type Handler struct {
	Config *ServerConfig
}

// NewHandler creates a new handler instance
func NewHandler(config *ServerConfig) *Handler {
	return &Handler{
		Config: config,
	}
}

// LoginPage displays the login page
func (h *Handler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Error": c.Query("error"),
	})
}

// LoginPost handles login form submission
func (h *Handler) LoginPost(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Verify credentials
	if username != h.Config.Username {
		c.Redirect(http.StatusFound, "/login?error=invalid")
		return
	}

	// Compare password with hash
	err := bcrypt.CompareHashAndPassword([]byte(h.Config.PasswordHash), []byte(password))
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=invalid")
		return
	}

	// Create session token (simple implementation)
	token := generateSessionToken()
	c.SetCookie("session_token", token, 3600*24, "/", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

// Logout handles user logout
func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// Dashboard displays the file list dashboard
func (h *Handler) Dashboard(c *gin.Context) {
	ctx := context.Background()

	// List all files
	files, err := h.Config.Storage.List(ctx, "/")
	if err != nil {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"Error": fmt.Sprintf("Failed to list files: %v", err),
			"Files": []interface{}{},
		})
		return
	}

	// Filter out directories and format file list
	var fileList []gin.H
	for _, file := range files {
		if !file.IsDir {
			fileList = append(fileList, gin.H{
				"Path":     file.Path,
				"Name":     filepath.Base(file.Path),
				"Size":     formatSize(file.Size),
				"ModTime":  file.ModTime,
				"Metadata": file.Metadata,
			})
		}
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Files":   fileList,
		"Success": c.Query("success"),
	})
}

// UploadPage displays the upload form
func (h *Handler) UploadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{})
}

// UploadPost handles file upload
func (h *Handler) UploadPost(c *gin.Context) {
	// Get form data
	file, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"Error": "Please select a file",
		})
		return
	}

	remotePath := c.PostForm("remote_path")
	algorithm := c.PostForm("algorithm")
	keyHex := c.PostForm("key")

	if remotePath == "" {
		remotePath = "/" + file.Filename
	}

	// Validate inputs
	if keyHex == "" {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"Error": "Encryption key is required",
		})
		return
	}

	// Create encryptor
	encryptor, err := createEncryptor(algorithm, keyHex)
	if err != nil {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"Error": fmt.Sprintf("Failed to create encryptor: %v", err),
		})
		return
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"Error": fmt.Sprintf("Failed to open file: %v", err),
		})
		return
	}
	defer src.Close()

	// Create uploader
	ul := uploader.NewUploader(encryptor, h.Config.Storage)

	// Read file to buffer
	var buf bytes.Buffer
	_, err = buf.ReadFrom(src)
	if err != nil {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"Error": fmt.Sprintf("Failed to read file: %v", err),
		})
		return
	}

	// Upload file
	ctx := context.Background()
	metadata := map[string]string{
		"original_name": file.Filename,
		"original_size": fmt.Sprintf("%d", file.Size),
	}
	err = ul.UploadStream(ctx, &buf, remotePath, metadata)
	if err != nil {
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"Error": fmt.Sprintf("Failed to upload file: %v", err),
		})
		return
	}

	c.Redirect(http.StatusFound, "/?success=File uploaded successfully")
}

// Download handles file download
func (h *Handler) Download(c *gin.Context) {
	// Get encrypted file path
	path := c.Param("path")
	if path == "" {
		c.String(http.StatusBadRequest, "File path is required")
		return
	}

	// Get encryption key from form
	keyHex := c.Query("key")
	algorithm := c.Query("algorithm")

	if keyHex == "" {
		c.String(http.StatusBadRequest, "Encryption key is required")
		return
	}

	// Create encryptor
	encryptor, err := createEncryptor(algorithm, keyHex)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid encryption algorithm")
		return
	}

	// Create uploader
	ul := uploader.NewUploader(encryptor, h.Config.Storage)

	// Download and decrypt file
	ctx := context.Background()
	var buf bytes.Buffer
	err = ul.DownloadStream(ctx, path, &buf)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to download file: %v", err)
		return
	}

	// Get original filename from metadata
	metadata, _ := h.Config.Storage.GetMetadata(ctx, path)
	filename := metadata["original_name"]
	if filename == "" {
		filename = filepath.Base(path)
		// Remove .enc extension if present
		filename = strings.TrimSuffix(filename, ".enc")
	}

	// Send file to browser
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
}

// Delete handles file deletion
func (h *Handler) Delete(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.Redirect(http.StatusFound, "/?error=File path is required")
		return
	}

	ctx := context.Background()
	err := h.Config.Storage.Delete(ctx, path)
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/?error=Failed to delete file: %v", err))
		return
	}

	c.Redirect(http.StatusFound, "/?success=File deleted successfully")
}

// Info displays file information
func (h *Handler) Info(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	ctx := context.Background()
	metadata, err := h.Config.Storage.GetMetadata(ctx, path)
	if err != nil {
		c.HTML(http.StatusOK, "info.html", gin.H{
			"Error": fmt.Sprintf("Failed to get file info: %v", err),
		})
		return
	}

	c.HTML(http.StatusOK, "info.html", gin.H{
		"Path":     path,
		"Name":     filepath.Base(path),
		"Metadata": metadata,
	})
}

// GenKeyPage displays the key generation page
func (h *Handler) GenKeyPage(c *gin.Context) {
	c.HTML(http.StatusOK, "genkey.html", gin.H{})
}

// GenKeyPost handles key generation
func (h *Handler) GenKeyPost(c *gin.Context) {
	sizeStr := c.PostForm("size")
	var size int
	fmt.Sscanf(sizeStr, "%d", &size)

	if size <= 0 || size > 64 {
		c.HTML(http.StatusOK, "genkey.html", gin.H{
			"Error": "Key size must be between 1 and 64 bytes",
		})
		return
	}

	// Generate random key
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		c.HTML(http.StatusOK, "genkey.html", gin.H{
			"Error": fmt.Sprintf("Failed to generate key: %v", err),
		})
		return
	}

	keyHex := hex.EncodeToString(key)

	c.HTML(http.StatusOK, "genkey.html", gin.H{
		"GeneratedKey": keyHex,
		"Size":         size,
	})
}

// Helper functions

// createEncryptor creates an encryptor based on algorithm and key
func createEncryptor(algo string, keyHex string) (crypto.Encryptor, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid key format: %w", err)
	}

	switch algo {
	case "aes":
		return crypto.NewAESEncryptor(key)
	case "xor":
		return crypto.NewXOREncryptor(key)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algo)
	}
}

// generateSessionToken generates a random session token
func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// formatSize formats file size in human-readable format
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
