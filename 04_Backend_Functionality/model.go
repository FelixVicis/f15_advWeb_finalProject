package main

/*
model.go by Allen J. Mills
    CREATION: 11.17.15
    COMPLETION: mm.dd.yy

    This will serve as the file holding
    the structures for our application.
*/

type User struct {
	Email    string
	Username string
	Password string // hash string
}

type Session struct {
	User
	LoggedIn bool
	viewing  string
}

type blobbedImage struct {
	BlobId   string // blobid of image
	URL      string // ulr of image
	UsrEmail string // email of uploaded user
	Uploaded string // time
}
