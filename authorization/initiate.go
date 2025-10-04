package authorization

import (
	"encoding/json"
	"net/http"
	"time"
)

func HandleInitiate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OAuthCallback string `json:"oauth_callback"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Generate temporary credentials
	tempToken := generateToken()
	tempSecret := generateToken()

	store.mu.Lock()
	store.tempCredentials[tempToken] = &TempCredential{
		Token:     tempToken,
		Secret:    tempSecret,
		Callback:  req.OAuthCallback,
		Timestamp: time.Now(),
	}
	store.mu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"oauth_token":              tempToken,
		"oauth_token_secret":       tempSecret,
		"oauth_callback_confirmed": "true",
	})
}
