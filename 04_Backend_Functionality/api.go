package main

/*
api.go by Allen J. Mills
    CREATION: 11.17.15
    COMPLETION: mm.dd.yy


*/

import (
	"fmt"
)

func main() {
	fmt.Println("")
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
