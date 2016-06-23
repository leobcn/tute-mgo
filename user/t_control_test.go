// Copyright 2015 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"gopkg.in/mgo.v2"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/web"
)

var validpath = regexp.MustCompile("^/")

func Test_control01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("control01. create and delete users")

	// get database session
	session, err := mgo.Dial("localhost")
	if err != nil {
		tst.Errorf("Dial failed:\n%v", err)
		return
	}
	defer session.Close()

	// create control
	database := session.DB("testing_tute-mgo-01")
	control := NewControl(database)

	// http variables
	w, r := &web.MockResponseWriter{}, &http.Request{}

	// create user
	err = control.Create(w, r, &User{Name: "dorival", Email: "secret", Tag: 123})
	if err != nil {
		tst.Errorf("Create failed:\n%v", err)
	}

	// create another user
	err = control.Create(w, r, &User{Name: "bender", Email: "bender@futurama", Tag: 666})
	if err != nil {
		tst.Errorf("Create failed:\n%v", err)
	}

	// create test server
	server := httptest.NewServer(web.MakeHandler(validpath, Json2dat, control.Get, true))
	defer server.Close()

	// send request and get response
	dat := strings.NewReader(`{"name":"bender", "email":"bender@futurama", "tag":666}`)
	res, err := http.Post(server.URL, "application/json", dat)
	defer res.Body.Close()

	// check
	if err != nil {
		tst.Errorf("http.Post failed:\n%s", err)
	}
	if res.StatusCode != 200 {
		tst.Errorf("http.Post failed with Status = %v", res.Status)
		return
	}

	// check response
	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		tst.Errorf("ReadAll failed:\n%v", err)
		return
	}
	type ResponseData struct {
		Ok    bool
		Users []byte
	}
	var rd ResponseData
	err = json.Unmarshal(got, &rd)
	if err != nil {
		tst.Errorf("cannot unmarshal response\n%v", err)
		return
	}
	str := string(rd.Users)
	io.Pforan("users (str) = %v\n", str)

	// check users
	var users []*User
	json.Unmarshal(rd.Users, &users)
	if err != nil {
		tst.Errorf("cannot unmarshal users\n%v", err)
		return
	}
	chk.IntAssert(len(users), 1)
	u := users[0]
	io.Pf("got user = %+#v\n", u)
	chk.String(tst, "bender", u.Name)
	chk.String(tst, "bender@futurama", u.Email)

	// delete users
	err = control.DeleteMany(w, r, &User{Name: "dorival"})
	if err != nil {
		tst.Errorf("Delete failed:\n%v", err)
	}
	err = control.DeleteMany(w, r, &User{Name: "bender"})
	if err != nil {
		tst.Errorf("Delete failed:\n%v", err)
	}
}
