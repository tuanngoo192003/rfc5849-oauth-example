package authorization

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

func HandleToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	token := r.FormValue("oauth_token")
	verifier := r.FormValue("oauth_verifier")

	fmt.Printf("📥 [TOKEN] Token exchange request: %s, verifier: %s\n", token, verifier)

	store.mu.Lock()
	defer store.mu.Unlock()

	temp, exists := store.tempCredentials[token]
	if !exists || !temp.Authorized {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "oauth_problem=token_rejected")
		fmt.Printf("[TOKEN] Invalid or unauthorized token: %s\n", token)
		return
	}

	storedVerifier, exists := store.authorizedCredentials[token]
	if !exists || storedVerifier != verifier {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "oauth_problem=verifier_invalid")
		fmt.Printf("[TOKEN] Invalid verifier\n")
		return
	}

	// Generate access token
	accessToken := generateToken()
	accessSecret := generateToken()

	store.accessTokens[accessToken] = &AccessToken{
		Token:    accessToken,
		Secret:   accessSecret,
		Username: temp.Username,
	}

	// Clean up temporary credentials
	delete(store.tempCredentials, token)
	delete(store.authorizedCredentials, token)

	fmt.Printf("✅ [TOKEN] Access token issued: %s for user: %s\n", accessToken, temp.Username)

	// Return access token
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Fprintf(w, "oauth_token=%s&oauth_token_secret=%s", accessToken, accessSecret)
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// OAuth 1.0 signature generation (for reference - simplified version)
func generateSignature(method, urlString string, params map[string]string, consumerSecret, tokenSecret string) string {
	// Normalize parameters
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(params[k])))
	}
	paramString := strings.Join(pairs, "&")

	// Create signature base string
	baseString := fmt.Sprintf("%s&%s&%s",
		strings.ToUpper(method),
		url.QueryEscape(urlString),
		url.QueryEscape(paramString))

	// Create signing key
	signingKey := fmt.Sprintf("%s&%s",
		url.QueryEscape(consumerSecret),
		url.QueryEscape(tokenSecret))

	// Generate HMAC-SHA1 signature
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(baseString))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signature
}
