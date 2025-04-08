// package middleware

// import (
// 	"fmt"
// 	"net/http"
// )

// // Logger logs incoming HTTP requests
//
//	func Logger(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)
//			next.ServeHTTP(w, r)
//		})
//	}
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
