package main

/*
filename.go by Allen J. Mills
    mm.d.yy
    Description
*/

import (
	"encoding/json"
	"net/http"
)

func serveTemplate(res http.ResponseWriter, req *http.Request, templateName string) {
	memItem, err := getSession(req)
	if err != nil {
		// not logged in
		tpl.ExecuteTemplate(res, templateName, Session{})
		return
	}
	// logged in
	var sd Session
	json.Unmarshal(memItem.Value, &sd)
	sd.LoggedIn = true
	tpl.ExecuteTemplate(res, templateName, &sd)
}
