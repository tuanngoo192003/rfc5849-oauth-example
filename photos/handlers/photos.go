package handlers

import (
	"encoding/json"
	"fmt"
	"go-oauth1/photos/authorization"
	"net/http"
)

func HandlePhotos(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	fmt.Printf("[PHOTOS] Access request with auth: %s\n", authHeader)

	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Extract token from "Bearer <token>"
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authorization header"})
		return
	}

	authorization.GetStore().GetMu().RLock()
	access, exists := authorization.GetStore().GetAccessToken()[token]
	authorization.GetStore().GetMu().RUnlock()

	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid access token"})
		fmt.Printf("[PHOTOS] Invalid access token: %s\n", token)
		return
	}

	fmt.Printf("[PHOTOS] Access granted for user: %s\n", access.Username)

	// Return Jane's vacation photos
	photos := []map[string]interface{}{
		{
			"id":    "1",
			"title": "Beach Sunset",
			"url":   "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=400",
			"owner": access.Username,
		},
		{
			"id":    "2",
			"title": "Mountain View",
			"url":   "https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=400",
			"owner": access.Username,
		},
		{
			"id":    "3",
			"title": "Tropical Paradise",
			"url":   "https://images.unsplash.com/photo-1559827260-dc66d52bef19?w=400",
			"owner": access.Username,
		},
		{
			"id":    "4",
			"title": "Desert Dunes",
			"url":   "https://images.unsplash.com/photo-1509316785289-025f5b846b35?w=400",
			"owner": access.Username,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"photos": photos,
		"user":   access.Username,
	})
}
