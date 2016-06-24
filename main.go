// Copyright 2016 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"regexp"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/tute-mgo/user"
	"github.com/cpmech/web"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

/* test with

curl -X POST -d "{\"name\":\"dorival\"}" http://localhost:8080/addUser

curl -X POST -d "{\"name\":\"dorival\"}" http://localhost:8080/getUsers

*/

// entry point
func main() {

	// get database session
	session, err := mgo.Dial("localhost")
	if err != nil {
		chk.Panic("Dial failed:\n%v", err)
		return
	}
	defer session.Close()

	// create User control
	uc := user.NewControl(session.DB("tute-mgo-01"))

	// routes
	r := mux.NewRouter()
	r.HandleFunc("/", web.MakeErrorHandler(root))
	r.HandleFunc("/addUser", web.MakeHandler(regexp.MustCompile("^/(addUser)"), user.Json2dat, uc.Add, true))
	r.HandleFunc("/getUsers", web.MakeHandler(regexp.MustCompile("^/(getUsers)"), user.Json2dat, uc.Get, true))
	http.Handle("/", r)

	// listen
	io.Pf("starting server on :8080\n")
	http.ListenAndServe(":8080", r)
}

// root handles requests to '/'
func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("restricted access\n"))
}
