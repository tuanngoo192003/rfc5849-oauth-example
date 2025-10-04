package main

import (
	"fmt"
	"go-oauth1/photos/authorization"
	"go-oauth1/photos/handlers"
	"net/http"
)

func main() {
	// Info page
	http.HandleFunc("/", handlers.HandleHome)

	http.HandleFunc("/api/photos", handlers.HandlePhotos)

	// OAuth endpoints
	http.HandleFunc("/oauth/initiate", authorization.HandleInitiate)
	http.HandleFunc("/oauth/authorize", authorization.HandleAuthorize)
	http.HandleFunc("/oauth/authorize-submit", authorization.HandleAuthorizeSubmit)
	http.HandleFunc("/oauth/token", authorization.HandleToken)

	// Protected resources
	http.HandleFunc("/photos", handlers.HandlePhotos)

	http.ListenAndServe(":8081", nil)

	fmt.Println("🖼️  Photo Service running on http://localhost:8081")
	fmt.Println("   This is Jane's Photo Service (OAuth Provider)")
	fmt.Println("\n   Registered Client:")
	fmt.Println("   - Client Key:", authorization.GetClientKey())
	fmt.Println("   - Client Secret:", authorization.GetClientSecret())
	fmt.Println("\n   User Credentials:")
	fmt.Println("   - Username: jane")
	fmt.Println("   - Password: password123")
	http.ListenAndServe(":8081", nil)
}
