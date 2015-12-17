package main

/*
model.go by Allen J. Mills
    CREATION: 11.17.15
    COMPLETION: 12.1.15

    This will serve as the file holding
    the structures for our application.
*/
import "time"

type User struct {
	Email    string
	UserName string `datastore:"-"`
	Password string `json:"-"` // hash string
}

type blobbedImage struct {
	BlobSRC  string    // blobid of image
	URL      string    // ulr of image
	UsrEmail string    // email of uploaded user
	Uploaded time.Time // time
}
