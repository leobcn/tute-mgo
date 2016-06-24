// Copyright 2016 Dorival Pedroso. All rights reserved.
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
	"gopkg.in/mgo.v2/bson"

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

	// dummy http variables
	writer, request := &web.MockResponseWriter{}, &http.Request{}

	// create user
	_, err = control.Add(writer, request, &User{Name: "dorival", Email: "dorival@test.com"})
	if err != nil {
		tst.Errorf("Add failed:\n%v", err)
	}

	// create test server
	server := httptest.NewServer(web.MakeHandler(validpath, Json2dat, control.Get, true))
	defer server.Close()

	// send request and get response
	dat := strings.NewReader(`{"name":"dorival"}`)
	response, err := http.Post(server.URL, "application/json", dat)
	defer response.Body.Close()

	// check
	if err != nil {
		tst.Errorf("http.Post failed:\n%s", err)
	}
	if response.StatusCode != 200 {
		tst.Errorf("http.Post failed with Status = %v", response.Status)
		return
	}

	// check results
	results, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tst.Errorf("ReadAll failed:\n%v", err)
		return
	}
	io.Pfblue2("results = %v\n", string(results))
	var res interface{}
	err = json.Unmarshal(results, &res)
	if err != nil {
		tst.Errorf("Unmarshal failed:\n%v", err)
		return
	}
	r := res.(map[string]interface{})
	users := r["users"].([]interface{})
	chk.IntAssert(len(users), 1)
	user := users[0].(map[string]interface{})
	chk.String(tst, user["name"].(string), "dorival")

	// delete users
	_, err = control.DeleteMany(writer, request, &User{Name: "dorival"})
	if err != nil {
		tst.Errorf("Delete failed:\n%v", err)
	}
}

func Test_control02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("control02. get all users")

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

	// dummy http variables
	writer, request := &web.MockResponseWriter{}, &http.Request{}

	// create users
	_, err = control.Add(writer, request, &User{Name: "dorival", Email: "dorival@test.com"})
	if err != nil {
		tst.Errorf("Add failed:\n%v", err)
	}
	_, err = control.Add(writer, request, &User{Name: "bender", Email: "bender@test.org"})
	if err != nil {
		tst.Errorf("Add failed:\n%v", err)
	}
	_, err = control.Add(writer, request, &User{Name: "leela", Email: "leela@futurama.biz"})
	if err != nil {
		tst.Errorf("Add failed:\n%v", err)
	}
	_, err = control.Add(writer, request, &User{Name: "hermes", Email: "hermes@here.ca"})
	if err != nil {
		tst.Errorf("Add failed:\n%v", err)
	}

	// create test server
	server := httptest.NewServer(web.MakeHandler(validpath, Json2dat, control.Get, true))
	defer server.Close()

	// send request and get response
	dat := strings.NewReader(`{}`)
	response, err := http.Post(server.URL, "application/json", dat)
	defer response.Body.Close()

	// check
	if err != nil {
		tst.Errorf("http.Post failed:\n%s", err)
	}
	if response.StatusCode != 200 {
		tst.Errorf("http.Post failed with Status = %v", response.Status)
		return
	}

	// check results
	results, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tst.Errorf("ReadAll failed:\n%v", err)
		return
	}
	io.Pfblue2("results = %v\n", string(results))
	var res interface{}
	err = json.Unmarshal(results, &res)
	if err != nil {
		tst.Errorf("Unmarshal failed:\n%v", err)
		return
	}
	r := res.(map[string]interface{})
	users := r["users"].([]interface{})
	chk.IntAssert(len(users), 4)

	// delete users
	_, err = control.DeleteMany(writer, request, &User{Name: ""})
	if err != nil {
		tst.Errorf("Delete failed:\n%v", err)
	}
}

func Test_control03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("control03. delete one user")

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
	writer, request := &web.MockResponseWriter{}, &http.Request{}

	// create user
	_, err = control.Add(writer, request, &User{Name: "dorival", Email: "dorival@test.com"})
	if err != nil {
		tst.Errorf("Add failed:\n%v", err)
	}

	// create test server
	server := httptest.NewServer(web.MakeHandler(validpath, Json2dat, control.Get, true))
	defer server.Close()

	// send request and get response
	dat := strings.NewReader(`{"name":"dorival"}`)
	response, err := http.Post(server.URL, "application/json", dat)
	defer response.Body.Close()

	// check
	if err != nil {
		tst.Errorf("http.Post failed:\n%s", err)
	}
	if response.StatusCode != 200 {
		tst.Errorf("http.Post failed with Status = %v", response.Status)
		return
	}

	// check results
	results, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tst.Errorf("ReadAll failed:\n%v", err)
		return
	}
	io.Pfblue2("results = %v\n", string(results))
	var res interface{}
	err = json.Unmarshal(results, &res)
	if err != nil {
		tst.Errorf("Unmarshal failed:\n%v", err)
		return
	}
	r := res.(map[string]interface{})
	users := r["users"].([]interface{})
	chk.IntAssert(len(users), 1)
	user := users[0].(map[string]interface{})
	id := user["id"].(string)
	io.Pforan("id = %v\n", id)

	// delete user
	_, err = control.Delete(writer, request, &User{Id: bson.ObjectIdHex(id)})
	if err != nil {
		tst.Errorf("Delete failed:\n%v", err)
	}
}
