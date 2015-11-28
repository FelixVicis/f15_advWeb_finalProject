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
	sess, err := getSession(req)
	if err != nil {
		// not logged in
		tpl.ExecuteTemplate(res, templateName, Session{})
		return
	}
	// logged in
	var sd Session
	json.Unmarshal(sess.Value, &sd)
	sd.LoggedIn = true
	tpl.ExecuteTemplate(res, templateName, &sd)
}
