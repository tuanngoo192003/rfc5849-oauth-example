package handlers

import (
	"encoding/json"
	"fmt"
	"go-oauth1/printer/authorization"
	"net/http"
)

func FetchPhotos(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	authorization.GetSession().GetMu().Lock()
	accessToken := authorization.GetSession().GetAccessToken()
	authorization.GetSession().GetMu().Unlock()

	if accessToken == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	authorization.GetSession().GetMu().Lock()
	authorization.GetSession().SetLogs([]string{"Step 4: Fetching photos from Photo Service..."})
	authorization.GetSession().GetMu().Unlock()

	// Step 4: Access protected resource
	req, _ := http.NewRequest("GET", authorization.GetCredential().GetPhotoServiceURL()+"/api/photos", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		authorization.GetSession().GetMu().Lock()
		authorization.GetSession().SetLogs([]string{fmt.Sprintf("Error: %v", err)})
		authorization.GetSession().GetMu().Unlock()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Photos []authorization.Photo `json:"photos"`
		User   string                `json:"user"`
	}

	json.NewDecoder(resp.Body).Decode(&result)

	authorization.GetSession().GetMu().Lock()
	authorization.GetSession().SetPhotos(result.Photos)
	authorization.GetSession().SetLogs([]string{fmt.Sprintf("Successfully fetched %d photos for user: %s", len(result.Photos), result.User)})
	authorization.GetSession().GetMu().Unlock()

	http.Redirect(w, r, "/", http.StatusFound)
}
