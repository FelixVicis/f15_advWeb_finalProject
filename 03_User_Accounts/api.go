package main

/*
api.go by Allen J. Mills
    CREATION: 11.17.15
    COMPLETION: mm.dd.yy
*/
import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"io/ioutil"
	"net/http"
	"time"
)

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
	} else {
		user.UserName = req.FormValue("userName")
		// success logging in
		createSession(res, req, user)
		// redirect
		http.Redirect(res, req, "/", 302)
	}
}

func logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

	cookie, err := req.Cookie("session")
	// cookie is not set
	if err != nil {
		http.Redirect(res, req, "/failure", 302)
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
	} else {
		fmt.Fprint(res, "true")
	}
}

/*
   Our goal here will be to make all of the back end functionality

   Our goals:
        1) Allow users to sign in
            1.1) Store sign in passwords securely
        2) Allow users to make new account
            2.1) see 1.1
        3) Allow users, given url, to go to an image refrenced by it.
            3.1) retrive the blob id of an image by given url quickly
            3.2) if that url is not valid, post an error
        4) IF logged in, Allow user to upload a new image
            4.1) this image must against the blobbedImage structure in model.
                4.1.1) MUST have URL given by uploading user

    Our structure
        GET - / - Root
        GET - /view/[username]/[imageURL] - An image by username with url of imageurl
        GET - /view/[username] - all images by username
        GET - /view - All uploaded images
        GET - /upload - if logged in, give user the uploading form
        POST - /upload - add blobbedImage to datastore
            -redirect to /veiw/[username]
        GET - /login - give user the login form
        POST - /login - hold session data for logged in user
            -redirect to /view/[username]
        POST - /login/newuser - add user to datastore
            -redirect to /

*/
