package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger middleware with detailed error handling
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Try executing the handler and catch panic
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)

		log.Printf("Request: %s %s | Duration: %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// package middleware

// import (
// 	"log"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// // Logger middleware for Gin
// func Logger(c *gin.Context) {
// 	start := time.Now()

// 	// Before request
// 	c.Next()

// 	// After request
// 	duration := time.Since(start)
// 	log.Printf("Request: %s %s | Duration: %v", c.Request.Method, c.Request.URL.Path, duration)
// }
