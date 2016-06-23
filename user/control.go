// Copyright 2016 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/web"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Control defines a controller for User data
type Control struct {
	database   *mgo.Database
	collection *mgo.Collection
}

// NewControl returns a new User Controller
//  Note: the collection is assumed to be "users"
func NewControl(database *mgo.Database) *Control {
	return &Control{database, database.C("users")}
}

// Create creates new user in database
func (o *Control) Create(w http.ResponseWriter, r *http.Request, datUser interface{}) (err error) {

	// response
	res := web.Response{"OK": false}
	defer func() { fmt.Fprintf(w, res.String()) }()

	// parse data
	if datUser == nil {
		err = chk.Err("user.Control.Create: invalid datUser = %v", datUser)
		return
	}
	user := datUser.(*User)
	user.Id = bson.NewObjectId()
	user.Timestamp = time.Now()

	// write to database
	err = o.collection.Insert(user)
	if err != nil {
		return
	}

	// success
	res["OK"] = true
	return
}

// Get retrieves Users for given Id, Name, Email, Phone and Tag (in this order of preference)
func (o *Control) Get(w http.ResponseWriter, r *http.Request, datUser interface{}) (err error) {

	// respose
	res := web.Response{"OK": false}
	defer func() { fmt.Fprintf(w, res.String()) }()

	// parse data
	if datUser == nil {
		err = chk.Err("user.Control.Get: invalid datUser = %v", datUser)
		return
	}
	user := datUser.(*User)

	// find users
	users, err := o.find(user, true)
	if err != nil {
		err = chk.Err("user.Control.Get: cannot find users with datUser = %v", datUser)
		return
	}

	// success
	res["users"] = users
	res["OK"] = true
	return
}

// Delete deletes one user with given Id, Name, Email, Phone and Tag (in this order of preference)
func (o *Control) Delete(w http.ResponseWriter, r *http.Request, datUser interface{}) (err error) {

	// response
	res := web.Response{"OK": false}
	defer func() { fmt.Fprintf(w, res.String()) }()

	// parse data
	if datUser == nil {
		err = chk.Err("user.Control.Delete: invalid datUser = %v", datUser)
		return
	}
	user := datUser.(*User)

	// find users
	users, err := o.find(user, false)
	if err != nil {
		err = chk.Err("user.Control.Delete: cannot find user with datUser = %v", datUser)
		return
	}

	// delete
	err = o.collection.RemoveId(users[0].Id)
	if err != nil {
		chk.Err("user.Control.Delete: cannot remove user = %v; err = %v", users[0], err)
		return
	}

	// success
	res["OK"] = true
	return
}

// DeleteMany deletes users with given Id, Name, Email, Phone and Tag (in this order of preference)
func (o *Control) DeleteMany(w http.ResponseWriter, r *http.Request, datUser interface{}) (err error) {

	// response
	res := web.Response{"OK": false}
	defer func() { fmt.Fprintf(w, res.String()) }()

	// parse data
	if datUser == nil {
		err = chk.Err("user.Control.DeleteMany: invalid datUser = %v", datUser)
		return
	}
	user := datUser.(*User)

	// find users
	users, err := o.find(user, true)
	if err != nil {
		err = chk.Err("user.Control.DeleteMany: cannot find users with datUser = %v", datUser)
		return
	}

	// delete
	for _, u := range users {
		err = o.collection.RemoveId(u.Id)
		if err != nil {
			chk.Err("user.Control.DeleteMany: cannot remove user = %v; err = %v", u, err)
			return
		}
	}

	// success
	res["OK"] = true
	return
}

// find finds users with given Id, Name, Email, Phone and Tag (in this order of preference)
func (o *Control) find(user *User, many bool) (users []*User, err error) {

	// search key
	var key bson.M
	switch {
	case user.Id != "":
		key = bson.M{"_id": user.Id}
	case user.Name != "":
		key = bson.M{"name": user.Name}
	case user.Email != "":
		key = bson.M{"email": user.Email}
	case user.Phone != "":
		key = bson.M{"phone": user.Phone}
	default:
		key = bson.M{"tag": user.Tag}
	}

	// search
	query := o.collection.Find(key).Sort("-timestamp")
	n, err := query.Count()
	if n < 1 || err != nil {
		return
	}
	if many {
		query.All(&users)
		return
	}
	users = make([]*User, 1)
	query.One(&users)
	return
}
