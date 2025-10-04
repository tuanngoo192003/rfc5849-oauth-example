// printer-service.go
// Printer Service (OAuth Client) - Port 8080
package main

import (
	"fmt"
	"go-oauth1/printer/authorization"
	"go-oauth1/printer/handlers"
	"net/http"
)

func main() {
	authorization.NewOAuthCredential()

	http.HandleFunc("/", handlers.HandlerHome)
	http.HandleFunc("/start-oauth", handlers.StartAuth)
	http.HandleFunc("/callback", handlers.Callback)
	http.HandleFunc("/fetch-photos", handlers.FetchPhotos)
	http.HandleFunc("/print", handlers.Print)

	fmt.Println("🖨️  Printer Service running on http://localhost:8080")
	fmt.Println("   This is the Printer Service (OAuth Client)")
	fmt.Println("\n   Client Credentials:")
	fmt.Println("   - Client Key:", authorization.GetCredential().GetClientKey())
	fmt.Println("   - Client Secret:", authorization.GetCredential().GetClientSecret())
	fmt.Println("\n   🌐 Open http://localhost:8080 in your browser to start")
	http.ListenAndServe(":8080", nil)
}
