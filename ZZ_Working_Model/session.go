package main

/*
filename.go by Allen J. Mills
    mm.d.yy

    Description
*/

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"net/http"
)

type Session struct {
	User
	LoggedIn bool
	Viewing  []blobbedImage // list of urls for the images
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

func getSession(req *http.Request) (*memcache.Item, error) {
	ctx := appengine.NewContext(req)

	cookie, err := req.Cookie("session")
	if err != nil {
		return &memcache.Item{}, err
	}

	item, err := memcache.Get(ctx, cookie.Value)
	if err != nil {
		return &memcache.Item{}, err
	}
	return item, nil
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
