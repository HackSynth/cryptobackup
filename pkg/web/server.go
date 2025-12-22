package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Server represents the web server
type Server struct {
	Config  *ServerConfig
	Router  *gin.Engine
	Handler *Handler
}

// NewServer creates a new web server instance
func NewServer(config *ServerConfig) *Server {
	// Create handler
	handler := NewHandler(config)

	// Create router
	router := setupRouter(handler)

	return &Server{
		Config:  config,
		Router:  router,
		Handler: handler,
	}
}

// setupRouter configures the Gin router with all routes and middleware
func setupRouter(handler *Handler) *gin.Engine {
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(LoggerMiddleware())

	// Load HTML templates
	router.LoadHTMLGlob("pkg/web/templates/*.html")

	// Static files (for future CSS/JS if needed)
	// router.Static("/static", "./pkg/web/static")

	// Public routes (no authentication required)
	router.GET("/login", handler.LoginPage)
	router.POST("/login", handler.LoginPost)

	// Protected routes (authentication required)
	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/", handler.Dashboard)
		protected.GET("/upload", handler.UploadPage)
		protected.POST("/upload", handler.UploadPost)
		protected.GET("/download/*path", handler.Download)
		protected.POST("/delete/*path", handler.Delete)
		protected.GET("/info/*path", handler.Info)
		protected.GET("/genkey", handler.GenKeyPage)
		protected.POST("/genkey", handler.GenKeyPost)
		protected.GET("/logout", handler.Logout)
	}

	return router
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
	fmt.Printf("Starting CryptoBackup web server at http://%s\n", addr)
	fmt.Printf("Login with username: %s\n", s.Config.Username)
	return s.Router.Run(addr)
}

// StartServer is a convenience function to create and start the server
func StartServer(config *ServerConfig) error {
	server := NewServer(config)
	return server.Start()
}
