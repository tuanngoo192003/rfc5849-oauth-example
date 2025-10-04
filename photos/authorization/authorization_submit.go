package authorization

import (
	"fmt"
	"net/http"
)

func HandleAuthorizeSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	token := r.FormValue("oauth_token")
	username := r.FormValue("username")
	password := r.FormValue("password")
	authorize := r.FormValue("authorize")

	fmt.Printf("[AUTHORIZE-SUBMIT] User: %s, Token: %s, Action: %s\n", username, token, authorize)

	store.mu.Lock()
	temp, exists := store.tempCredentials[token]
	store.mu.Unlock()

	if !exists {
		http.Error(w, "Invalid token", http.StatusBadRequest)

		return
	}

	if authorize != "true" {
		fmt.Printf("[AUTHORIZE-SUBMIT] User denied authorization\n")
		http.Redirect(w, r, temp.Callback+"?oauth_problem=user_denied", http.StatusFound)

		return
	}

	// Simple authentication check
	expectedPassword, userExists := users[username]
	if !userExists || expectedPassword != password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		fmt.Printf("[AUTHORIZE-SUBMIT] Invalid credentials for user: %s\n", username)

		return
	}

	verifier := generateToken()

	store.mu.Lock()
	temp.Authorized = true
	temp.Username = username
	store.authorizedCredentials[token] = verifier
	store.mu.Unlock()

	fmt.Printf("[AUTHORIZE-SUBMIT] Authorization successful for user: %s, verifier: %s\n", username, verifier)

	// Redirect back to client with verifier
	redirectURL := fmt.Sprintf("%s?oauth_token=%s&oauth_verifier=%s", temp.Callback, token, verifier)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
