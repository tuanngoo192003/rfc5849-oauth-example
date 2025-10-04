package handlers

import (
	"encoding/json"
	"go-oauth1/authorization"
	"net/http"
	"strings"
)

func HandlePhotos(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
		return
	}

	token := parts[1]

	authorization.GetStore().GetMu().RLock()
	access, exists := authorization.GetStore().GetAccessToken()[token]
	authorization.GetStore().GetMu().RUnlock()

	if !exists {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

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
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"photos": photos,
		"user":   access.Username,
	})
}
