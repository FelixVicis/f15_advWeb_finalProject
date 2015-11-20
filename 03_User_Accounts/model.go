package main

/*
model.go by Allen J. Mills
    CREATION: 11.17.15
    COMPLETION: mm.dd.yy

    This will serve as the file holding
    the structures for our application.
*/
import "time"

type User struct {
	Email    string
	UserName string `datastore:"-"`
	Password string `json:"-"` // hash string
}

type Session struct {
	User
	LoggedIn bool
	viewing  string
}

type blobbedImage struct {
	BlobId   string    // blobid of image
	URL      string    // ulr of image
	UsrEmail string    // email of uploaded user
	Uploaded time.Time // time
}
