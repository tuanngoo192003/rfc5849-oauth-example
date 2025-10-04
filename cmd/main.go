package main

import (
	"fmt"
	"go-oauth1/authorization"
	"go-oauth1/internal/handlers"
	"net/http"
)

// Client credentials (printer.example.com)
const (
	clientKey    = "dpf43f3p2l4k3l03"
	clientSecret = "kd94hf93k423kf44"
)

func main() {
	// Enable CORS
	http.HandleFunc("/", authorization.CorsMiddleware(handlers.HandleRoot))
	http.HandleFunc("/initiate", authorization.CorsMiddleware(authorization.HandleInitiate))
	http.HandleFunc("/authorize", authorization.CorsMiddleware(authorization.HandleAuthorize))
	http.HandleFunc("/authorize-submit", authorization.CorsMiddleware(authorization.HandleAuthorizeSubmit))
	http.HandleFunc("/token", authorization.CorsMiddleware(authorization.HandleToken))
	http.HandleFunc("/photos", authorization.CorsMiddleware(handlers.HandlePhotos))

	fmt.Println("OAuth 1.0 Server running on :8080")
	fmt.Println("Client Key:", clientKey)
	fmt.Println("Client Secret:", clientSecret)
	http.ListenAndServe(":8080", nil)
}
