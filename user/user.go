// Copyright 2016 Dorival Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"encoding/json"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Phone     string        `json:"phone"`
	Tag       int           `json:"tag"`
	Timestamp time.Time     `json:"timestamp"`
}

func Json2dat(in []byte) interface{} {
	var dat User
	err := json.Unmarshal(in, &dat)
	if err != nil {
		return nil
	}
	return &dat
}
