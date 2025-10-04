package handlers

import (
	"fmt"
	"go-oauth1/printer/authorization"
	"io"
	"net/http"
	"net/url"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("oauth_token")
	verifier := r.URL.Query().Get("oauth_verifier")
	problem := r.URL.Query().Get("oauth_problem")

	if problem != "" {
		authorization.GetSession().GetMu().Lock()
		authorization.GetSession().SetLogs([]string{fmt.Sprintf("Authorization failed: %s", problem)})
		authorization.GetSession().GetMu().Unlock()
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	authorization.GetSession().GetMu().Lock()
	authorization.GetSession().SetLogs([]string{"User authorized! Received verifier",
		"Step 3: Exchanging temporary credentials for access token...",
	})
	authorization.GetSession().GetMu().Unlock()

	// Step 3: Exchange for access token
	data := url.Values{}
	data.Set("oauth_token", token)
	data.Set("oauth_verifier", verifier)

	resp, err := http.PostForm(authorization.GetCredential().GetPhotoServiceURL()+"/oauth/token", data)
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

	accessToken := params.Get("oauth_token")
	accessSecret := params.Get("oauth_token_secret")

	authorization.GetSession().GetMu().Lock()
	authorization.GetSession().SetAccessToken(accessToken)
	authorization.GetSession().SetAccessToken(accessSecret)
	authorization.GetSession().SetLogs([]string{fmt.Sprintf("Access token received: %s", accessToken),
		"OAuth flow complete! Ready to access photos",
	})
	authorization.GetSession().GetMu().Unlock()

	http.Redirect(w, r, "/", http.StatusFound)
}
