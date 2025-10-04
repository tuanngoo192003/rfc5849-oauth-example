package handlers

import (
	"fmt"
	"go-oauth1/printer/authorization"
	"io"
	"net/http"
	"net/url"
)

func StartAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authorization.GetSession().GetMu().Lock()
	authorization.GetSession().SetLogs([]string{"Starting OAuth 1.0 flow...", "Step 1: Requesting temporary credentials from Photo Service"})
	authorization.GetSession().GetMu().Unlock()

	// Step 1: Request temporary credentials
	data := url.Values{}
	data.Set("oauth_consumer_key", authorization.GetCredential().GetClientKey())
	data.Set("oauth_callback", authorization.GetCredential().GetCallbackURL())

	resp, err := http.PostForm(authorization.GetCredential().GetPhotoServiceURL()+"/oauth/initiate", data)
	if err != nil {
		authorization.GetSession().GetMu().Lock()
		authorization.GetSession().SetLogs([]string{fmt.Sprintf("Error: %v", err)})
		authorization.GetSession().GetMu().Unlock()

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	params, _ := url.ParseQuery(string(body))

	tempToken := params.Get("oauth_token")
	tempSecret := params.Get("oauth_token_secret")

	authorization.GetSession().GetMu().Lock()
	authorization.GetSession().SetTempToken(tempToken)
	authorization.GetSession().SetTempSecret(tempSecret)
	authorization.GetSession().SetLogs([]string{fmt.Sprintf("Received temporary token: %s", tempToken),
		"Step 2: Redirecting user to Photo Service for authorization..."})
	authorization.GetSession().GetMu().Unlock()

	// Step 2: Redirect to authorization page
	authURL := fmt.Sprintf("%s/oauth/authorize?oauth_token=%s", authorization.GetCredential().GetPhotoServiceURL(), tempToken)
	http.Redirect(w, r, authURL, http.StatusFound)
}
