package authorization

import (
	"encoding/json"
	"net/http"
)

func HandleAuthorizeSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OAuthToken string `json:"oauth_token"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		Authorize  bool   `json:"authorize"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	store.mu.Lock()
	temp, exists := store.tempCredentials[req.OAuthToken]
	if !exists {
		store.mu.Unlock()
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
		return
	}

	if !req.Authorize {
		store.mu.Unlock()
		json.NewEncoder(w).Encode(map[string]string{"error": "Authorization denied"})
		return
	}

	// Simple authentication check
	if req.Username == "" {
		store.mu.Unlock()
		json.NewEncoder(w).Encode(map[string]string{"error": "Username required"})
		return
	}

	verifier := generateToken()
	temp.Authorized = true
	temp.Username = req.Username
	store.authorizedCredentials[req.OAuthToken] = verifier
	store.mu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"oauth_token":    req.OAuthToken,
		"oauth_verifier": verifier,
	})
}
