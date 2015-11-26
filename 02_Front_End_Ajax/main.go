package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"encoding/json"

	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

var tpl *template.Template

func init() {
	r := httprouter.New()
	r.GET("/", home)                            // root
	r.GET("/login", login)                      // public user has requested a login.
	r.GET("/logout", logout)                    // signed in user has requested a log out
	r.GET("/signup", signup)                    // public user has requested a new user
	r.POST("/api/checkusername", checkUserName) // form has posted to api, check username
	r.POST("/api/createuser", createUser)       // signup has posted to api
	r.POST("/api/login", loginProcess)          // login has posted to api
	r.GET("/api/logout", logout)                // logout has posted to api
	r.GET("/failure", failure)                  // a step has gone awry

	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))
	http.Handle("/", r)

	tpl = template.Must(tpl.ParseGlob("public/templates/*.html"))
}

// ROOT ===================================================================================================

func home(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//ctx := appengine.NewContext(req)
	// get session
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
	tpl.ExecuteTemplate(res, "login.html", nil)
}

func loginProcess(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "Users", req.FormValue("userName"), 0, nil)
	var user User
	err := datastore.Get(ctx, key, &user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.FormValue("password"))) != nil {
		// failure logging in
		http.Redirect(res, req, "/failure", 302)
		// var sd Session
		// sd.LoginFail = true
		// tpl.ExecuteTemplate(res, "login.html", sd)
		return
	}
	user.UserName = req.FormValue("userName")
	// success logging in
	createSession(res, req, user)
	// redirect
	http.Redirect(res, req, "/", 302)
}

func logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

	cookie, err := req.Cookie("session")
	// cookie is not set
	if err != nil {
		http.Redirect(res, req, "/", 302)
		return
	}

	// clear memcache
	sd := memcache.Item{
		Key:        cookie.Value,
		Value:      []byte(""),
		Expiration: time.Duration(1 * time.Microsecond),
	}
	memcache.Set(ctx, &sd)

	// clear the cookie
	cookie.MaxAge = -1
	http.SetCookie(res, cookie)

	// redirect
	http.Redirect(res, req, "/", 302)
}

// NEW USER ===================================================================================================

func signup(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	tpl.ExecuteTemplate(res, "signup.html", nil)
}

func createUser(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	user := User{
		Email:    req.FormValue("email"),
		UserName: req.FormValue("userName"),
		Password: string(hashedPass),
	}
	key := datastore.NewKey(ctx, "Users", user.UserName, 0, nil)
	key, err = datastore.Put(ctx, key, &user)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}

	createSession(res, req, user)
	// redirect
	http.Redirect(res, req, "/", 302)
}

// HELPERS ===================================================================================================

func failure(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	tpl.ExecuteTemplate(res, "failure.html", nil)
}

func checkUserName(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	bs, err := ioutil.ReadAll(req.Body)
	sbs := string(bs)
	var user User
	key := datastore.NewKey(ctx, "Users", sbs, 0, nil)
	err = datastore.Get(ctx, key, &user)
	// if there is an err, there is NO user
	if err != nil {
		// there is an err, there is a NO user
		fmt.Fprint(res, "false")
		return
	}
	fmt.Fprint(res, "true")
}

func createSession(res http.ResponseWriter, req *http.Request, user User) {
	ctx := appengine.NewContext(req)
	// SET COOKIE
	id, _ := uuid.NewV4()
	cookie := &http.Cookie{
		Name:  "session",
		Value: id.String(),
		Path:  "/",
		// twenty minute session:
		MaxAge: 60 * 20,
		//		UNCOMMENT WHEN DEPLOYED:
		//		Secure: true,
		//		HttpOnly: true,
	}
	http.SetCookie(res, cookie)

	// SET MEMCACHE session data (sd)
	json, err := json.Marshal(user)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	sd := memcache.Item{
		Key:   id.String(),
		Value: json,
	}
	memcache.Set(ctx, &sd)
}
