package authorization

import (
	"fmt"
	"net/http"
	"time"
)

func HandleInitiate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	consumerKey := r.FormValue("oauth_consumer_key")
	callback := r.FormValue("oauth_callback")

	fmt.Printf("[INITIATE] Received request from client: %s, callback: %s\n", consumerKey, callback)

	if consumerKey != clientKey {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "oauth_problem=consumer_key_rejected")
		fmt.Printf("[INITIATE] Invalid consumer key: %s\n", consumerKey)
		return
	}

	// Generate temporary credentials
	tempToken := generateToken()
	tempSecret := generateToken()

	store.mu.Lock()
	store.tempCredentials[tempToken] = &TempCredential{
		Token:     tempToken,
		Secret:    tempSecret,
		Callback:  callback,
		Timestamp: time.Now(),
	}
	store.mu.Unlock()

	fmt.Printf("[INITIATE] Issued temp credentials: %s\n", tempToken)

	// Return as application/x-www-form-urlencoded
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Fprintf(w, "oauth_token=%s&oauth_token_secret=%s&oauth_callback_confirmed=true",
		tempToken, tempSecret)
}
