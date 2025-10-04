package authorization

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
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

	var req struct {
		OAuthToken    string `json:"oauth_token"`
		OAuthVerifier string `json:"oauth_verifier"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	store.mu.Lock()
	defer store.mu.Unlock()

	temp, exists := store.tempCredentials[req.OAuthToken]
	if !exists || !temp.Authorized {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid or unauthorized token"})
		return
	}

	verifier, exists := store.authorizedCredentials[req.OAuthToken]
	if !exists || verifier != req.OAuthVerifier {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid verifier"})
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
	delete(store.tempCredentials, req.OAuthToken)
	delete(store.authorizedCredentials, req.OAuthToken)

	json.NewEncoder(w).Encode(map[string]string{
		"oauth_token":        accessToken,
		"oauth_token_secret": accessSecret,
	})
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
