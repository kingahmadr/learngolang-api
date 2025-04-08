// package utils

// import "fmt"

// // Example utility function
// func PrintMessage(msg string) {
// 	fmt.Println("Message:", msg)
// }

package utils

import "fmt"

// LogInfo logs an informational message
func LogInfo(message string) {
    fmt.Println("[INFO] " + message)
}
