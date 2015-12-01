package main

/*
api.go by Allen J. Mills
    CREATION: 11.17.15
    COMPLETION: 12.1.15
*/
import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
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
		serveTemplateWithParams(res, req, "falure.html", "YOUR PASSWORD DOES NOT MATCH")
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

/* Blobstore Section ----------------------------------------- */
// BLOB ======================================================================================================

func uploadForm(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// serving image upload form
	ctx := appengine.NewContext(req)
	uploadURL, err := blobstore.UploadURL(ctx, "/upload", nil) // setting up the post call for blobstore.
	if err != nil {
		serveTemplateWithParams(res, req, "falure.html", "BLOB URL ERROR") // there was an issue making the post call
		return
	}

	_, err = getSession(req)
	if err != nil {
		serveTemplateWithParams(res, req, "falure.html", "YOU MAY NOT SUBMIT FILES WITHOUT BEING LOGGED IN") // there was an issue making the post call
		return
	}

	serveTemplateWithParams(res, req, "upload.html", uploadURL)
}

func uploadToBlob(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// posting image to blob, posting blobbedImage to datastore
	blobs, _, err := blobstore.ParseUpload(req) // from request, take and parse the blobs from the form.
	if err != nil {
		serveTemplateWithParams(res, req, "falure.html", "BLOB PARSE ERROR") // something went wrong with the parse
		return
	}

	file := blobs["file"] // okay. have blobs now.
	if len(file) == 0 {   // are there any files?
		serveTemplateWithParams(res, req, "falure.html", "NO BLOBS FOUND")
		return
	}

	sess, good := getSession(req) // ensure that there is a good session.
	if good != nil {
		serveTemplateWithParams(res, req, "falure.html", "LOGIN TOKEN HAS EXPIRED, CANNOT PROCESS IMAGE")
		return
	}
	var sd Session
	json.Unmarshal(sess.Value, &sd)

	blobImage := blobbedImage{ // make the blob based on what we took in.
		BlobSRC:  makeImageURL(req, string(file[0].BlobKey)),
		URL:      sd.UserName,
		UsrEmail: sd.Email,
		Uploaded: file[0].CreationTime,
	}

	ctx := appengine.NewContext(req) // prep and submit to datastore
	key := datastore.NewKey(ctx, "Images", blobImage.URL, 0, nil)
	key, err = datastore.Put(ctx, key, &blobImage)
	if err != nil {
		serveTemplateWithParams(res, req, "falure.html", "INTERNAL DATASTORE ERROR")
		return
	}

	http.Redirect(res, req, "/view/"+blobImage.URL, http.StatusFound)
	//http.Redirect(res, req, "/", 302)
}

func makeImageURL(req *http.Request, blob string) string {
	// helper to turn blob into image request string
	return "https://" + req.URL.Host + "/image/" + blob
}

func getImage(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// requesting an image based on blob key
	blobstore.Send(res, appengine.BlobKey(ps.ByName("blobKey")))
}

func apiGetImageURL(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// change blob to url as API call
	bs, _ := ioutil.ReadAll(req.Body)
	sbs := string(bs)

	fmt.Fprint(res, makeImageURL(req, sbs))
}

func requestImage(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//user has requested to see :key image
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "Images", ps.ByName("key"), 0, nil) // request datastore info on image
	var bi blobbedImage
	err := datastore.Get(ctx, key, &bi)
	if err != nil {
		serveTemplateWithParams(res, req, "falure.html", "INTERNAL DATASTORE ERROR, IMAGE REQUEST FAILED\nERROR: "+err.Error()+"\nrequesting key: "+ps.ByName("key"))
		return
	}
	serveTemplateWithParams(res, req, "image.html", bi) // got it? good. send image out.
}

func requestAllImage(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//user has requested to see all images
	sess, _ := getSession(req) // get login info if exists.
	var sd Session
	json.Unmarshal(sess.Value, &sd)

	ctx := appengine.NewContext(req)
	var links []blobbedImage
	for t := datastore.NewQuery("Images").Run(ctx); ; { /// get all images in Images from datastore
		var x blobbedImage
		_, err := t.Next(&x)
		if err == datastore.Done {
			break
		}
		if err != nil {
			serveTemplateWithParams(res, req, "falure.html", "INTERNAL DATASTORE ERROR, IMAGE REQUEST FAILED\nERROR: "+err.Error())
			return
		}
		links = append(links, x)
	}

	sd.Viewing = links  // add the images to our session data.
	if sd.Email != "" { // do a quick logged in check.
		sd.LoggedIn = true
	} else {
		sd.LoggedIn = false
	}

	serveTemplateWithParams(res, req, "imageMulti.html", sd) // serve.
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
