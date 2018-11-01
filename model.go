package main

import (
	"gopkg.in/mgo.v2/bson"
)

// Item - Our struct for all items
type Item struct {
	ID          bson.ObjectId `bson:"_id" json:"id" `
	Name        *string       `bson:"name" json:"name"`
	Value       *string       `bson:"value" json:"value"`
	Description string        `bson:"description" json:"description"`
	ItemId      string        `bson:"-" json:"-"`
	Count       int           `bson:"count" json:"count"`
}

//Config db connection info struct
type Config struct {
	Hosts    string `json:"hosts"`
	Database string `json:"database"`
	Item     string `json:"item"`
	User     string `json:"user"`
}

//User info
type User struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Username *string       `bson:"username" json:"username"`
	Password *string       `bson:"password" json:"password"`
	Token    string        `json:"token"`
}

// JwtToken struct
type JwtToken struct {
	Token string `json:"token"`
}

// Exception jwt message
type Exception struct {
	Message string `json:"message"`
}
