package authorization

import (
	"fmt"
	"net/http"
)

func HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("oauth_token")

	store.mu.RLock()
	_, exists := store.tempCredentials[token]
	store.mu.RUnlock()

	if !exists {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	// In a real implementation, this would be a proper login page
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
		<h1>Authorize Application</h1>
		<p>The printer service wants to access your photos.</p>
		<p>Token: %s</p>
		<form method="POST" action="/authorize-submit">
			<input type="hidden" name="oauth_token" value="%s">
			<input type="text" name="username" placeholder="Username"><br>
			<button type="submit" name="authorize" value="true">Authorize</button>
			<button type="submit" name="authorize" value="false">Deny</button>
		</form>
	`, token, token)
}
