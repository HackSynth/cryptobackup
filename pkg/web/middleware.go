package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("session_token")
		if err != nil || session == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Store session in context for handlers to use
		c.Set("session_token", session)
		c.Next()
	}
}

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		fmt.Printf("[%s] %s %s - %d (%v)\n",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			duration)
	}
}
