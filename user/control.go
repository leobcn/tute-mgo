// Copyright 2016 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
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

// Add creates new user in database
func (o *Control) Add(w http.ResponseWriter, r *http.Request, datUser interface{}) (res *web.Results, err error) {

	// parse data
	if datUser == nil {
		err = chk.Err("user.Control.Add: invalid datUser = %v", datUser)
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
	return
}

// Get retrieves users
func (o *Control) Get(w http.ResponseWriter, r *http.Request, datUser interface{}) (res *web.Results, err error) {

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

	// results
	res = &web.Results{"users": users}
	return
}

// Delete deletes one user
func (o *Control) Delete(w http.ResponseWriter, r *http.Request, datUser interface{}) (res *web.Results, err error) {

	// parse data
	if datUser == nil {
		err = chk.Err("user.Control.Delete: invalid datUser = %v", datUser)
		return
	}
	user := datUser.(*User)

	// find users
	users, err := o.find(user, false)
	if err != nil {
		err = chk.Err("user.Control.Delete: cannot find user with datUser = %v. err = %v", datUser, err)
		return
	}

	// check
	if users[0] == nil {
		chk.Err("user.Control.Delete: cannot find user to delete. datUser = %v", datUser)
		return
	}

	// delete
	err = o.collection.RemoveId(users[0].Id)
	if err != nil {
		chk.Err("user.Control.Delete: cannot remove user = %v; err = %v", users[0], err)
		return
	}
	return
}

// DeleteMany deletes users
func (o *Control) DeleteMany(w http.ResponseWriter, r *http.Request, datUser interface{}) (res *web.Results, err error) {

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
	for _, c := range users {
		if c == nil {
			continue
		}
		err = o.collection.RemoveId(c.Id)
		if err != nil {
			chk.Err("user.Control.DeleteMany: cannot remove user = %v; err = %v", c, err)
			return
		}
	}
	return
}

// find finds users with given Id, Name, Email, Phone and Tag (in this order of preference)
// Note: empty user will return all users
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
	case user.Tag > 0:
		key = bson.M{"tag": user.Tag}
	}

	// search
	query := o.collection.Find(key).Sort("-timestamp")
	n, err := query.Count()
	if err != nil {
		return
	}
	if n < 1 {
		err = chk.Err("cannot find entry in database")
		return
	}
	if many {
		query.All(&users)
		return
	}
	users = []*User{&User{}}
	err = query.One(users[0])
	return
}
