package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	r := httprouter.New()
	r.GET("/", home) // root

	r.GET("/login", login)                      // public user has requested a login.
	r.GET("/logout", logout)                    // signed in user has requested a log out
	r.GET("/signup", signup)                    // public user has requested a new user
	r.POST("/api/checkusername", checkUserName) // form has posted to api, check username
	r.POST("/api/createuser", createUser)       // signup has posted to api
	r.POST("/api/login", loginProcess)          // login has posted to api
	r.GET("/api/logout", logout)                // logout has posted to api

	r.GET("/failure", failure) // a step has gone awry

	r.POST("/upload", uploadToBlob) // User has submitted a file to the blob
	r.GET("/upload", uploadForm)    // User has requested a submission

	r.GET("/image/:blobKey", getImage)
	r.POST("/api/image", apiGetImageURL)

	r.GET("/viewAll", requestAllImage)
	r.GET("/view/:key", requestImage)

	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))
	http.Handle("/", r)

	tpl = template.Must(tpl.ParseGlob("public/templates/*.html"))
}

// ROOT ===================================================================================================

func home(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	memItem, err := getSession(req)
	var sd Session
	if err == nil {
		// logged in
		json.Unmarshal(memItem.Value, &sd)
		sd.LoggedIn = true
	}
	tpl.ExecuteTemplate(res, "home.html", &sd)
}

// LOGIN/LOGOUT ==============================================================================================

func login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	serveTemplate(res, req, "login.html")
}

// NEW USER ===================================================================================================

func signup(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	serveTemplate(res, req, "signup.html")
}

// HELPERS ===================================================================================================

func failure(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	serveTemplateWithParams(res, req, "falure.html", "NO MESSAGE AVAILABLE")
}
