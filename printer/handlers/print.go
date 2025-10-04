package handlers

import (
	"fmt"
	"go-oauth1/printer/authorization"
	"net/http"
)

func Print(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	photoTitle := r.FormValue("photo_title")

	authorization.GetSession().GetMu().Lock()
	authorization.GetSession().SetLogs([]string{fmt.Sprintf("PRINTING: %s", photoTitle)})
	authorization.GetSession().GetMu().Unlock()

	http.Redirect(w, r, "/", http.StatusFound)
}
